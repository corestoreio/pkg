// Copyright 2015-2016, Cyrill @ Schumacher.fm and the CoreStore contributors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License at distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package dbr

import (
	"strings"

	"github.com/corestoreio/csfw/util/bufferpool"
)

// IfNull appends the (optional) alias to the IFNULL expression. Argument count
// can be between 1-n. IFNULL(expr1,expr2) If expr1 is not NULL, IFNULL()
// returns expr1; otherwise it returns expr2. IFNULL() returns a numeric or
// string value, depending on the context in which it is used. See the examples.
func IfNull(expressionAlias ...string) string {
	buf := bufferpool.Get()
	defer bufferpool.Put(buf)

	switch len(expressionAlias) {
	case 1:
		ifNullQuote2(buf, expressionAlias[0], "NULL ")
	case 2:
		// Input:  dbr.IfNull(col1,col2,alias)
		// Output: IFNULL(`col1`, `col2`)

		// Input:  dbr.IfNull(expr1,expr2,alias)
		// Output: IFNULL(expr1, expr1)
		ifNullQuote2(buf, expressionAlias...)

	case 3:
		// Input:  dbr.IfNull(table1.col1,table2.col2,alias)
		// Output: IFNULL(`table1`.`col1`, `table2`.`col2`) AS `alias`

		// Input:  dbr.IfNull(col1,col2,alias)
		// Output: IFNULL(`col1`, `col2`) AS `alias`

		// Input:  dbr.IfNull(expr1,expr2,alias)
		// Output: IFNULL(expr1, expr1) AS `alias`
		ifNullQuote2(buf, expressionAlias...)
		buf.WriteString(" AS ")
		Quoter.quote(buf, expressionAlias[2])

	case 4:
		// Input:  dbr.IfNull(table1,col1,table2,col2,alias)
		// Output: IFNULL(`table1`.`col1`, `table2`.`col2`) AS `alias`
		ifNullQuote4(buf, expressionAlias[:3]...)
	case 5:
		// Input:  dbr.IfNull(table1,col1,table2,col2,alias)
		// Output: IFNULL(`table1`.`col1`, `table2`.`col2`) AS `alias`
		ifNullQuote4(buf, expressionAlias[:3]...)
		buf.WriteString(" AS ")
		Quoter.quote(buf, expressionAlias[4])
	default:
		ifNullQuote4(buf, expressionAlias[:3]...)
		buf.WriteString(" AS ")
		Quoter.quote(buf, strings.Join(expressionAlias[4:], "_"))
	}
	return buf.String()
}

func ifNullQuote2(w queryWriter, expressionAlias ...string) {
	w.WriteString("IFNULL(")
	if isValidIdentifier(expressionAlias[0]) == 0 {
		Quoter.quoteAs(w, expressionAlias[0])
	} else {
		w.WriteRune('(')
		w.WriteString(expressionAlias[0])
		w.WriteRune(')')
	}
	w.WriteRune(',')
	if isValidIdentifier(expressionAlias[1]) == 0 {
		Quoter.quoteAs(w, expressionAlias[1])
	} else {
		w.WriteRune('(')
		w.WriteString(expressionAlias[1])
		w.WriteRune(')')
	}
	w.WriteRune(')')
}

func ifNullQuote4(w queryWriter, expressionAlias ...string) {
	w.WriteString("IFNULL(")
	Quoter.quote(w, expressionAlias[:2]...)
	w.WriteRune(',')
	Quoter.quote(w, expressionAlias[2:4]...)
	w.WriteRune(')')
}

func If(expression, true, false string) {
	// IF((%s), %s, %s)
}
