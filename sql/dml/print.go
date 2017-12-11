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
)

// QueryBuilder assembles a query and returns the raw SQL without parameter
// substitution and the arguments.
type QueryBuilder interface {
	ToSQL() (string, []interface{}, error)
}

type queryBuilder interface {
	toSQL(w *bytes.Buffer, placeHolders []string) ([]string, error)
	writeBuildCache(sql []byte)
	// readBuildCache returns the cached SQL string
	readBuildCache() (sql []byte)
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
		w.WriteString("/*ID:")
		w.WriteString(id)
		w.WriteString("*/ ")
	}
}
