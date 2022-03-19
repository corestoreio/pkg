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
	"fmt"
	"strconv"
	"unicode/utf8"

	"github.com/corestoreio/errors"
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
	ToSQL() (string, []any, error)
}

// QuerySQL a helper type to transform a string into a QueryBuilder compatible
// type.
type QuerySQLFn func() (string, []any, error)

// ToSQL satisfies interface QueryBuilder and returns always nil arguments and
// nil error.
func (fn QuerySQLFn) ToSQL() (string, []any, error) {
	return fn()
}

// QuerySQL simple type to satisfy the QueryBuilder interface.
type QuerySQL string

// ToSQL satisfies interface QueryBuilder and returns always nil arguments and
// nil error.
func (qs QuerySQL) ToSQL() (string, []any, error) {
	return string(qs), nil, nil
}

// queryBuilder must support thread safety when writing and reading the cache.
type queryBuilder interface {
	toSQL(w *bytes.Buffer, placeHolders []string) ([]string, error)
}

// BuilderBase contains fields which all SQL query builder have in common, the
// same base. Exported for documentation reasons.
type BuilderBase struct {
	Table id
	// IsUnsafe if set to true the functions AddColumn* will turn any
	// non valid identifier (not `{a-z}[a-z0-9$_]+`i) into an expression.
	IsUnsafe         bool
	ärgErr           error
	isWithDBR        bool // tuple handling before building the SQL string
	containsTuples   bool
	qualifiedColumns []string
}

// Clone creates a clone of the current object.
func (bb BuilderBase) Clone() BuilderBase {
	cc := bb
	cc.Table = bb.Table.Clone()
	cc.ärgErr = nil
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

	buf := bufferpool.Get()
	defer bufferpool.Put(buf)
	qualifiedColumns, err := qb.toSQL(buf, []string{})
	if err != nil {
		return "", errors.WithStack(err)
	}

	// the qualifiedColumns might have an entry from Conditions.write to
	// indicate there is a tuple placeholder.
	qualifiedColumns2 := qualifiedColumns[:0]
	for _, pc := range qualifiedColumns {
		if pc != placeHolderTuples {
			qualifiedColumns2 = append(qualifiedColumns2, pc)
		} else {
			bb.containsTuples = true
		}
	}
	bb.qualifiedColumns = qualifiedColumns2

	return buf.String(), nil
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

func sqlWriteUnionAll(w *bytes.Buffer, isAll, isIntersect, isExcept bool) {
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
			writeNumber(w, offsetCount)
			w.WriteByte(',')
		}
		writeNumber(w, limitCount)
	}
}

type writeNumberTypes interface {
	~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64 | ~int64 | ~float64
}

func writeNumber[N writeNumberTypes](w *bytes.Buffer, i N) (err error) {
	d := w.Bytes()
	w.Reset()
	switch it := any(i).(type) { // https://github.com/golang/go/issues/45380#issuecomment-1014950980
	case uint:
		_, err = w.Write(strconv.AppendUint(d, uint64(it), 10))
	case uint8:
		_, err = w.Write(strconv.AppendUint(d, uint64(it), 10))
	case uint16:
		_, err = w.Write(strconv.AppendUint(d, uint64(it), 10))
	case uint32:
		_, err = w.Write(strconv.AppendUint(d, uint64(it), 10))
	case uint64:
		_, err = w.Write(strconv.AppendUint(d, it, 10))
	case int64:
		_, err = w.Write(strconv.AppendInt(d, it, 10))
	case float64:
		_, err = w.Write(strconv.AppendFloat(d, it, 'g', -1, 64))
	default:
		panic(fmt.Sprintf("type not supported: %T", any(i)))
	}
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

const (
	tupleTemplate      = `(?)(?,?)(?,?,?)(?,?,?,?)(?,?,?,?,?)(?,?,?,?,?,?)(?,?,?,?,?,?,?)(?,?,?,?,?,?,?,?)(?,?,?,?,?,?,?,?,?)(?,?,?,?,?,?,?,?,?,?)(?,?,?,?,?,?,?,?,?,?,?)(?,?,?,?,?,?,?,?,?,?,?,?)(?,?,?,?,?,?,?,?,?,?,?,?,?)(?,?,?,?,?,?,?,?,?,?,?,?,?,?)(?,?,?,?,?,?,?,?,?,?,?,?,?,?,?)`
	tupleTemplateCount = 15
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

func writeTuplePlaceholders(buf *bytes.Buffer, rowCount, columnCount uint) {
	start, end := calcInsertTemplatePlaceholderPos(columnCount)
	for r := uint(0); r < rowCount; r++ {
		if r > 0 {
			buf.WriteByte(',')
		}
		if columnCount <= tupleTemplateCount {
			buf.WriteString(tupleTemplate[start:end])
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
