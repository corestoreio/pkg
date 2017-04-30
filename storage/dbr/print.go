// Copyright 2015-2017, Cyrill @ Schumacher.fm and the CoreStore contributors
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

package dbr

import (
	"fmt"
	"strconv"
)

// QueryBuilder assembles a query and returns the raw SQL without parameter
// substitution and the arguments.
type QueryBuilder interface {
	ToSQL() (string, Arguments, error)
}

// queryWriter at used to generate a query.
type queryWriter interface {
	WriteString(s string) (n int, err error)
	WriteRune(r rune) (n int, err error)
}

func makeSQL(b QueryBuilder) string {
	sRaw, vals, err := b.ToSQL()
	if err != nil {
		return fmt.Sprintf("[dbr] ToSQL Error: %+v", err)
	}
	sql, err := Interpolate(sRaw, vals...)
	if err != nil {
		return fmt.Sprintf("[dbr] Interpolate Error: %+v", err)
	}
	return sql
}

// String returns a string representing a preprocessed, interpolated, query.
// On error, the error gets printed. Fulfills interface fmt.Stringer.
func (b *Delete) String() string {
	return makeSQL(b)
}

// String returns a string representing a preprocessed, interpolated, query.
// On error, the error gets printed. Fulfills interface fmt.Stringer.
func (b *Insert) String() string {
	return makeSQL(b)
}

// String returns a string representing a preprocessed, interpolated, query.
// On error, the error gets printed. Fulfills interface fmt.Stringer.
func (b *Select) String() string {
	return makeSQL(b)
}

// String returns a string representing a preprocessed, interpolated, query.
// On error, the error gets printed. Fulfills interface fmt.Stringer.
func (b *Update) String() string {
	return makeSQL(b)
}

func sqlWriteUnionAll(w queryWriter, isAll bool) {
	w.WriteString("\nUNION")
	if isAll {
		w.WriteString(" ALL")
	}
	w.WriteRune('\n')
}

func sqlWriteOrderBy(w queryWriter, orderBys []string, br bool) {
	if len(orderBys) == 0 {
		return
	}
	brS := ' '
	if br {
		brS = '\n'
	}
	w.WriteRune(brS)
	w.WriteString("ORDER BY ")
	for i, s := range orderBys {
		if i > 0 {
			w.WriteString(", ")
		}
		w.WriteString(s)
	}
}

func sqlWriteLimitOffset(w queryWriter, limitValid bool, limitCount uint64, offsetValid bool, offsetCount uint64) {
	if limitValid {
		w.WriteString(" LIMIT ")
		w.WriteString(strconv.FormatUint(limitCount, 10))
	}
	if offsetValid {
		w.WriteString(" OFFSET ")
		w.WriteString(strconv.FormatUint(offsetCount, 10))
	}
}
