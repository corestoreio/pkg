// Copyright 2015-present, Cyrill @ Schumacher.fm and the CoreStore contributors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package dml

import (
	"bytes"
	"context"
	"fmt"
	"sort"
	"strconv"
	"sync"
	"unicode/utf8"

	"github.com/corestoreio/errors"
	"github.com/corestoreio/log"
	"github.com/corestoreio/pkg/util/bufferpool"
)

const (
	dmlSourceSelect       = 's'
	dmlSourceInsert       = 'i'
	dmlSourceInsertSelect = 'I'
	dmlSourceUpdate       = 'u'
	dmlSourceDelete       = 'd'
	dmlSourceWith         = 'w'
	dmlSourceUnion        = 'n'
	dmlSourceShow         = 'h'
)

type writer interface {
	WriteByte(c byte) error
	WriteRune(r rune) (int, error)
	Write(p []byte) (int, error)
}

// QueryBuilder assembles a query and returns the raw SQL without parameter
// substitution and the arguments. The input arguments might be modified and
// returned as plain primitive types.
type QueryBuilder interface {
	ToSQL() (string, []interface{}, error)
}

// QuerySQL a helper type to transform a string into a QueryBuilder compatible
// type.
type QuerySQLFn func() (string, []interface{}, error)

// ToSQL satisfies interface QueryBuilder and returns always nil arguments and
// nil error.
func (fn QuerySQLFn) ToSQL() (string, []interface{}, error) {
	return fn()
}

// queryBuilder must support thread safety when writing and reading the cache.
type queryBuilder interface {
	toSQL(w *bytes.Buffer, placeHolders []string) ([]string, error)
}

// builderCommon
type builderCommon struct {
	defaultQualifier string
	// ID of a statement. Used in logging. The ID gets generated with function
	// signature `func() string`. This func gets applied to the logger when
	// setting up a logger.
	id string // tracing ID
	// ärgErr represents an argument error caused in any of the other functions.
	// A stack has been attached to the error to identify properly the source.
	ärgErr error // Sorry Germans for that terrible pun #notSorry
	// source defines with which DML statement the builderCommon struct has been initialized.
	// Constants are `dmlType*`
	source rune
	Log    log.Logger // Log optional logger
	// templateStmtCount only used in case a UNION statement acts as a template.
	// Create one SELECT statement and by setting the data for
	// Union.StringReplace function additional SELECT statements are getting
	// created. Now the arguments must be multiplied by the number of new
	// created SELECT statements. This value  gets stored in templateStmtCount.
	// An example exists in TestUnionTemplate_ReuseArgs.
	templateStmtCount int
	// EstimatedCachedSQLSize specifies the estimated size in bytes of the final
	// SQL string. This value gets used during SQL string building process to
	// reduce the allocations and speed up the process. Default Value is xxxx
	// Bytes.
	EstimatedCachedSQLSize uint16

	cacheKey string
	// SingleUseCacheKey      bool // TODO implement, should panic when setting the same cahce key the 2nd time
	// cachedSQL contains the final SQL string which gets send to the server.
	// Using the CacheKey allows a dml type (insert,update,select ... ) to build
	// multiple different versions from object.
	cachedSQL map[string]string
	// qualifiedColumns gets collected before calling ToSQL, and clearing the all
	// pointers, to know which columns need values from the QualifiedRecords
	qualifiedColumns []string
	// DB can be either a *sql.DB (connection pool), a *sql.Conn (a single
	// dedicated database session) or a *sql.Tx (an in-progress database
	// transaction).
	db QueryExecPreparer
}

func (bc *builderCommon) withCacheKey(key string, args ...interface{}) {
	if len(args) > 0 {
		key = fmt.Sprintf(key, args...)
	}
	bc.cacheKey = key
}

func (bc *builderCommon) CachedQueries(queries ...string) []string {
	keys := make([]string, 0, len(bc.cachedSQL))
	for key := range bc.cachedSQL {
		keys = append(keys, key)
	}
	sort.Strings(keys)
	for _, k := range keys {
		queries = append(queries, k, bc.cachedSQL[k])
	}
	return queries
}

func (bc *builderCommon) cachedSQLUpsert(key string, sql string) {
	if bc.cachedSQL == nil {
		bc.cachedSQL = make(map[string]string, 32) // 32 is just a guess
	}
	bc.cachedSQL[key] = sql
}

// BuilderBase contains fields which all SQL query builder have in common, the
// same base. Exported for documentation reasons.
type BuilderBase struct {
	Table id
	// PropagationStopped set to true if you would like to interrupt the
	// listener chain. Once set to true all sub sequent calls of the next
	// listeners will be suppressed.
	PropagationStopped bool
	// IsUnsafe if set to true the functions AddColumn* will turn any
	// non valid identifier (not `{a-z}[a-z0-9$_]+`i) into an expression.
	IsUnsafe bool

	rwmu sync.RWMutex // also protects the whole SQL string building process
	builderCommon
}

// Clone creates a clone of the current object.
func (bb BuilderBase) Clone() BuilderBase {
	cc := bb
	cc.Table = bb.Table.Clone()
	cc.rwmu = sync.RWMutex{}
	cc.builderCommon.qualifiedColumns = cloneStringSlice(bb.builderCommon.qualifiedColumns)
	return cc
}

// buildToSQL builds the raw SQL string and caches it as a byte slice. It gets
// called by toSQL.
// buildArgsAndSQL generates the SQL string and its place holders. Takes care of
// caching. It returns the string with placeholders.
func (bb *BuilderBase) buildToSQL(qb queryBuilder) (string, error) {
	if bb.ärgErr != nil {
		return "", errors.WithStack(bb.ärgErr)
	}

	rawSQL, ok := bb.cachedSQL[bb.cacheKey]
	if !ok {
		buf := bufferpool.Get()
		defer bufferpool.Put(buf)
		qualifiedColumns, err := qb.toSQL(buf, []string{})
		if err != nil {
			return "", errors.WithStack(err)
		}
		rawSQL = buf.String()
		bb.qualifiedColumns = qualifiedColumns
		bb.cachedSQLUpsert(bb.cacheKey, rawSQL)
	}
	return rawSQL, nil
}

func (bb *BuilderBase) prepare(ctx context.Context, db Preparer, qb queryBuilder, source rune) (_ *Stmt, err error) {
	if in, ok := qb.(*Insert); ok && in != nil && !in.IsBuildValues {
		return nil, errors.NotAcceptable.Newf("[dml] did you forgot to call .BuildValues()?")
	}

	rawQuery, err := bb.buildToSQL(qb)
	if bb.Log != nil && bb.Log.IsDebug() {
		defer log.WhenDone(bb.Log).Debug("Prepare", log.Err(err), log.String("sql", rawQuery))
	}
	if err != nil {
		return nil, errors.WithStack(err)
	}

	sqlStmt, err := db.PrepareContext(ctx, rawQuery)
	if err != nil {
		return nil, errors.Wrapf(err, "[dml] Prepare.PrepareContext with query %q", rawQuery)
	}

	stmt := &Stmt{
		base: bb.builderCommon,
		Stmt: sqlStmt,
	}
	stmt.base.cacheKey = bb.cacheKey
	stmt.base.cachedSQLUpsert(bb.cacheKey, rawQuery)
	stmt.base.db = stmtWrapper{stmt: sqlStmt}
	stmt.base.source = source
	return stmt, nil
}

// newDBR builds the SQl string and creates a new DBR object for
// collecting arguments and later querying.
func (bb *BuilderBase) newDBR(qb queryBuilder) *DBR {
	bb.rwmu.Lock()
	_, err := bb.buildToSQL(qb)
	a := DBR{
		base: bb.builderCommon,
	}
	if err != nil {
		a.base.ärgErr = errors.WithStack(err)
	}
	bb.rwmu.Unlock()
	return &a
}

// BuilderConditional defines base fields used in statements which can have
// conditional constraints like WHERE, JOIN, ORDER, etc. Exported for
// documentation reasons.
type BuilderConditional struct {
	Joins    Joins
	Wheres   Conditions
	OrderBys ids
	// OrderByRandColumnName defines the column name of the single primary key
	// in a table to build the optimized ORDER BY RAND() JOIN clause.
	OrderByRandColumnName string
	LimitCount            uint64
	LimitValid            bool
}

// Clone creates a new clone of the current object.
func (b BuilderConditional) Clone() BuilderConditional {
	c := b
	c.Joins = b.Joins.Clone()
	c.Wheres = b.Wheres.Clone()
	c.OrderBys = b.OrderBys.Clone()
	return c
}

func (b *BuilderConditional) join(j string, t id, on ...*Condition) {
	jf := &join{
		JoinType: j,
		Table:    t,
	}
	jf.On = append(jf.On, on...)
	b.Joins = append(b.Joins, jf)
}

func sqlObjToString(rawSQL string, err error) string {
	if err != nil {
		return fmt.Sprintf("[dml] String Error: %+v", err)
	}
	return rawSQL
}

// String returns a string representing a preprocessed, interpolated, query.
// On error, the error gets printed. Fulfills interface fmt.Stringer.
func (b *Delete) String() string {
	return sqlObjToString(b.buildToSQL(b))
}

// String returns a string representing a preprocessed, interpolated, query.
// On error, the error gets printed. Fulfills interface fmt.Stringer.
func (b *Insert) String() string {
	return sqlObjToString(b.buildToSQL(b))
}

// String returns a string representing a preprocessed, interpolated, query.
// On error, the error gets printed. Fulfills interface fmt.Stringer.
func (b *Select) String() string {
	return sqlObjToString(b.buildToSQL(b))
}

// String returns a string representing a preprocessed, interpolated, query.
// On error, the error gets printed. Fulfills interface fmt.Stringer.
func (b *Update) String() string {
	return sqlObjToString(b.buildToSQL(b))
}

// String returns a string representing a preprocessed, interpolated, query.
// On error, the error gets printed. Fulfills interface fmt.Stringer.
func (u *Union) String() string {
	return sqlObjToString(u.buildToSQL(u))
}

// String returns a string representing a preprocessed, interpolated, query.
// On error, the error gets printed. Fulfills interface fmt.Stringer.
func (b *With) String() string {
	return sqlObjToString(b.buildToSQL(b))
}

// String returns a string representing a preprocessed, interpolated, query.
// On error, the error gets printed. Fulfills interface fmt.Stringer.
func (b *Show) String() string {
	return sqlObjToString(b.buildToSQL(b))
}

func sqlWriteUnionAll(w *bytes.Buffer, isAll bool, isIntersect bool, isExcept bool) {
	w.WriteByte('\n')
	switch {
	case isIntersect:
		w.WriteString("INTERSECT") // MariaDB >= 10.3
	case isExcept:
		w.WriteString("EXCEPT") // MariaDB >= 10.3
	default:
		w.WriteString("UNION")
		if isAll {
			w.WriteString(" ALL")
		}
	}
	w.WriteByte('\n')
}

func sqlWriteOrderBy(w *bytes.Buffer, orderBys ids, br bool) {
	if len(orderBys) == 0 {
		return
	}
	brS := ' '
	if br {
		brS = '\n'
	}
	w.WriteRune(brS)
	w.WriteString("ORDER BY ")
	orderBys.writeQuoted(w, nil)
}

// LIMIT 0,0 quickly returns an empty set. This can be useful for checking the
// validity of a query. When using one of the MySQL APIs, it can also be
// employed for obtaining the types of the result columns.
func sqlWriteLimitOffset(w *bytes.Buffer, limitValid, offsetValid bool, offsetCount, limitCount uint64) {
	if limitValid {
		w.WriteString(" LIMIT ")
		if offsetValid {
			writeUint64(w, offsetCount)
			w.WriteByte(',')
		}
		writeUint64(w, limitCount)
	}
}

func writeFloat64(w *bytes.Buffer, f float64) (err error) {
	d := w.Bytes()
	w.Reset()
	_, err = w.Write(strconv.AppendFloat(d, f, 'g', -1, 64))
	return err
}

func writeInt64(w *bytes.Buffer, i int64) (err error) {
	d := w.Bytes()
	w.Reset()
	_, err = w.Write(strconv.AppendInt(d, i, 10))
	return err
}

func writeUint64(w *bytes.Buffer, i uint64) (err error) {
	d := w.Bytes()
	w.Reset()
	_, err = w.Write(strconv.AppendUint(d, i, 10))
	return err
}

func writeBytes(w *bytes.Buffer, p []byte) (err error) {
	switch {
	case p == nil:
		_, err = w.WriteString(sqlStrNullUC)
	case !utf8.Valid(p):
		dialect.EscapeBinary(w, p)
	default:
		dialect.EscapeString(w, string(p)) // maybe create an EscapeByteString version to avoid one alloc ;-)
	}
	return
}

func writeStmtID(w *bytes.Buffer, id string) {
	if id != "" {
		w.WriteString("/*ID$") // colon not possible because used for named arguments.
		w.WriteString(id)
		w.WriteString("*/ ")
	}
}

const (
	insertTemplate      = `(?)(?,?)(?,?,?)(?,?,?,?)(?,?,?,?,?)(?,?,?,?,?,?)(?,?,?,?,?,?,?)(?,?,?,?,?,?,?,?)(?,?,?,?,?,?,?,?,?)(?,?,?,?,?,?,?,?,?,?)(?,?,?,?,?,?,?,?,?,?,?)(?,?,?,?,?,?,?,?,?,?,?,?)(?,?,?,?,?,?,?,?,?,?,?,?,?)(?,?,?,?,?,?,?,?,?,?,?,?,?,?)(?,?,?,?,?,?,?,?,?,?,?,?,?,?,?)`
	insertTemplateCount = 15
)

// i = 1 => pos = 0:3 | 1 + 0 + 2
// i = 2 => pos = 3:8 | 2 + 1 + 2
// i = 3 => pos = 8:15 | 3 + 2 + 2
// i = 4 => pos = 15:24 | 4 + 3 + 2
// i = 5 => pos = 24:35 | 5 + 4 + 2 <= 5 = number of placeholders; 4 number of colons; 2 number of brackets

func calcInsertTemplatePlaceholderPos(columnCount uint) (start, end uint) {
	var colons uint
	for i := uint(1); i <= columnCount; i++ {
		colons = i - 1
		start = end
		end = colons + start + i + 2
	}
	if columnCount == 1 {
		start = 0
	}
	return
}

func writeInsertPlaceholders(buf *bytes.Buffer, rowCount, columnCount uint) {
	start, end := calcInsertTemplatePlaceholderPos(columnCount)
	for r := uint(0); r < rowCount; r++ {
		if r > 0 {
			buf.WriteByte(',')
		}
		if columnCount <= insertTemplateCount {
			buf.WriteString(insertTemplate[start:end])
		} else {
			buf.WriteByte('(')
			for c := uint(0); c < columnCount; c++ {
				if c > 0 {
					buf.WriteByte(',')
				}
				buf.WriteByte('?')
			}
			buf.WriteByte(')')
		}
	}
}
