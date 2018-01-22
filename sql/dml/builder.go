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
	"strconv"
	"unicode/utf8"

	"github.com/corestoreio/errors"
	"github.com/corestoreio/log"
)

const (
	dmlSourceSelect = 's'
	dmlSourceInsert = 'i'
	dmlSourceUpdate = 'u'
	dmlSourceDelete = 'd'
	dmlSourceWith   = 'w'
	dmlSourceUnion  = 'n'
	dmlSourceShow   = 'h'
)

type writer interface {
	WriteByte(c byte) error
	WriteRune(r rune) (int, error)
	Write(p []byte) (int, error)
}

// QueryBuilder assembles a query and returns the raw SQL without parameter
// substitution and the arguments.
type QueryBuilder interface {
	ToSQL() (string, []interface{}, error)
}

// QuerySQL a helper type to transform a string into a QueryBuilder compatible
// type.
type QuerySQL string

// ToSQL satisfies interface QueryBuilder and returns always nil arguments and
// nil error.
func (s QuerySQL) ToSQL() (string, []interface{}, error) {
	return string(s), nil, nil
}

type queryBuilder interface {
	toSQL(w *bytes.Buffer, placeHolders []string) ([]string, error)
	writeBuildCache(sql []byte)
	// readBuildCache returns the cached SQL string
	readBuildCache() (sql []byte)
}

// builderCommon
type builderCommon struct {
	// cachedSQL contains the final SQL string which gets send to the server.
	cachedSQL []byte
	// EstimatedCachedSQLSize specifies the estimated size in bytes of the final
	// SQL string. This value gets used during SQL string building process to
	// reduce the allocations and speed up the process. Default Value is xxxx
	// Bytes.
	EstimatedCachedSQLSize uint16
	// source defines with which DML statement the builderCommon struct has been initialized.
	// Constants are `dmlType*`
	source rune
	// ID of a statement. Used in logging. The ID gets generated with function
	// signature `func() string`. This func gets applied to the logger when
	// setting up a logger.
	id  string     // tracing ID
	Log log.Logger // Log optional logger

	// ärgErr represents an argument error caused in one of the three With
	// functions.
	ärgErr error // Sorry Germans for that terrible pun #notSorry

	defaultQualifier string
	// isWithInterfaces will be set to true if the raw interface arguments are
	// getting applied.
	isWithInterfaces bool
	// qualifiedColumns gets collected before calling ToSQL, and clearing the all
	// pointers, to know which columns need values from the QualifiedRecords
	qualifiedColumns []string
	// templateStmtCount only used in case a UNION statement acts as a template.
	// Create one SELECT statement and by setting the data for
	// Union.StringReplace function additional SELECT statements are getting
	// created. Now the arguments must be multiplied by the number of new
	// created SELECT statements. This value  gets stored in templateStmtCount.
	// An example exists in TestUnionTemplate_ReuseArgs.
	templateStmtCount int
	// DB can be either a *sql.DB (connection pool), a *sql.Conn (a single
	// dedicated database session) or a *sql.Tx (an in-progress database
	// transaction).
	DB QueryExecPreparer
}

func (b *builderCommon) prepare(ctx context.Context, db QueryExecPreparer, qb QueryBuilder, source rune) (_ *Stmt, err error) {
	var sqlStr string
	sqlStr, _, err = qb.ToSQL()
	if b.Log != nil && b.Log.IsDebug() {
		defer log.WhenDone(b.Log).Debug("Prepare", log.Err(err), log.String("sql", sqlStr))
	}
	if err != nil {
		return nil, errors.WithStack(err)
	}
	sqlStmt, err := db.PrepareContext(ctx, sqlStr)
	if err != nil {
		return nil, errors.Wrapf(err, "[dml] Prepare.PrepareContext with query %q", sqlStr)
	}

	stmt := &Stmt{
		base: *b,
		Stmt: sqlStmt,
	}
	stmt.base.DB = stmtWrapper{stmt: sqlStmt}
	stmt.base.source = source
	return stmt, nil

}

// estimatedCachedSQLSize 1024 bytes value got retrieved by analyzing and
// reviewing some M2 SQL queries.
const estimatedCachedSQLSize = 1024

// BuilderBase contains fields which all SQL query builder have in common, the
// same base. Exported for documentation reasons.
type BuilderBase struct {
	builderCommon
	RawFullSQL string
	Table      id
	// PropagationStopped set to true if you would like to interrupt the
	// listener chain. Once set to true all sub sequent calls of the next
	// listeners will be suppressed.
	PropagationStopped   bool
	IsBuildCacheDisabled bool // see DisableBuildCache()
	// IsUnsafe if set to true the functions AddColumn* will turn any
	// non valid identifier (not `{a-z}[a-z0-9$_]+`i) into an expression.
	IsUnsafe bool
	// propagationStoppedAt position in the slice where the stopped propagation
	// has been requested. for every new iteration the propagation must stop at
	// this position.
	propagationStoppedAt int
}

// WithArgs sets the optional interfaced arguments for the later execution.
func (bb *BuilderBase) withArgs(qb QueryBuilder, rawArgs ...interface{}) *Arguments {
	sqlStr, argsRaw, err := qb.ToSQL()
	var args [defaultArgumentsCapacity]argument
	a := Arguments{
		base: bb.builderCommon, // might be a source of a possible race condition, fix later
		raw:  append(rawArgs, argsRaw...),
		args: args[:0],
	}
	a.base.cachedSQL = []byte(sqlStr)
	a.base.ärgErr = errors.WithStack(err)
	return &a
}

// hasBuildCache satisfies partially interface queryBuilder
func (bb *BuilderBase) hasBuildCache() bool {
	return !bb.IsBuildCacheDisabled
}

// buildToSQL builds the raw SQL string and caches it as a byte slice. It gets
// called by toSQL.
// buildArgsAndSQL generates the SQL string and its place holders. Takes care of
// caching. It returns the string with placeholders.
func (bb *BuilderBase) buildToSQL(qb queryBuilder) ([]byte, error) {
	if bb.ärgErr != nil {
		return nil, errors.WithStack(bb.ärgErr)
	}
	rawSQL := qb.readBuildCache()
	if rawSQL == nil || bb.IsBuildCacheDisabled {
		bb.qualifiedColumns = bb.qualifiedColumns[:0]
		// Pre allocating that with a decent size, can speed up writing due to
		// less re-slicing / buffer.Grow.
		size := bb.EstimatedCachedSQLSize
		if size == 0 {
			size = estimatedCachedSQLSize
		}
		buf := bytes.NewBuffer(make([]byte, 0, size))
		var err error
		bb.qualifiedColumns, err = qb.toSQL(buf, bb.qualifiedColumns)
		if err != nil {
			return nil, errors.WithStack(err)
		}
		if !bb.IsBuildCacheDisabled {
			qb.writeBuildCache(buf.Bytes())
		}
		rawSQL = buf.Bytes()
	}
	return rawSQL, nil
}

// BuilderConditional defines base fields used in statements which can have
// conditional constraints like WHERE, JOIN, ORDER, etc. Exported for
// documentation reasons.
type BuilderConditional struct {
	Joins      Joins
	Wheres     Conditions
	OrderBys   ids
	LimitCount uint64
	LimitValid bool
}

func (b *BuilderConditional) join(j string, t id, on ...*Condition) {
	jf := &join{
		JoinType: j,
		Table:    t,
	}
	jf.On = append(jf.On, on...)
	b.Joins = append(b.Joins, jf)
}

func sqlObjToString(rawSQL []byte, err error) string {
	if err != nil {
		return fmt.Sprintf("[dml] String Error: %+v", err)
	}
	return string(rawSQL)
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

func sqlWriteLimitOffset(w *bytes.Buffer, limitValid bool, limitCount uint64, offsetValid bool, offsetCount uint64) {
	if limitValid {
		w.WriteString(" LIMIT ")
		writeUint64(w, limitCount)
	}
	if offsetValid {
		w.WriteString(" OFFSET ")
		writeUint64(w, offsetCount)
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
