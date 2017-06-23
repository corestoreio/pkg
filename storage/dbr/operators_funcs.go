// Copyright 2015-2017, Cyrill @ Schumacher.fm and the CoreStore contributors
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

// SQLIfNull appends the (optional) alias to the IFNULL expression. Argument
// count can be between 1-n. IFNULL(expr1,expr2) If expr1 is not NULL, IFNULL()
// returns expr1; otherwise it returns expr2. IFNULL() returns a numeric or
// string value, depending on the context in which it is used. See the examples.
func SQLIfNull(expressionAlias ...string) string {
	buf := bufferpool.Get()
	defer bufferpool.Put(buf)

	switch len(expressionAlias) {
	case 1:
		sqlIfNullQuote2(buf, expressionAlias[0], "NULL ")
	case 2:
		// Input:  dbr.SQLIfNull(col1,col2,alias)
		// Output: IFNULL(`col1`, `col2`)

		// Input:  dbr.SQLIfNull(expr1,expr2,alias)
		// Output: IFNULL(expr1, expr1)
		sqlIfNullQuote2(buf, expressionAlias...)

	case 3:
		// Input:  dbr.SQLIfNull(table1.col1,table2.col2,alias)
		// Output: IFNULL(`table1`.`col1`, `table2`.`col2`) AS `alias`

		// Input:  dbr.SQLIfNull(col1,col2,alias)
		// Output: IFNULL(`col1`, `col2`) AS `alias`

		// Input:  dbr.SQLIfNull(expr1,expr2,alias)
		// Output: IFNULL(expr1, expr1) AS `alias`
		sqlIfNullQuote2(buf, expressionAlias...)
		buf.WriteString(" AS ")
		Quoter.writeName(buf, expressionAlias[2])

	case 4:
		// Input:  dbr.SQLIfNull(table1,col1,table2,col2)
		// Output: IFNULL(`table1`.`col1`, `table2`.`col2`)
		sqlIfNullQuote4(buf, expressionAlias...)
	case 5:
		// Input:  dbr.SQLIfNull(table1,col1,table2,col2,alias)
		// Output: IFNULL(`table1`.`col1`, `table2`.`col2`) AS `alias`
		sqlIfNullQuote4(buf, expressionAlias[:4]...)
		buf.WriteString(" AS ")
		Quoter.writeName(buf, expressionAlias[4])
	default:
		sqlIfNullQuote4(buf, expressionAlias[:4]...)
		buf.WriteString(" AS ")
		Quoter.writeName(buf, strings.Join(expressionAlias[4:], "_"))
	}
	return buf.String()
}

func sqlIfNullQuote2(w queryWriter, expressionAlias ...string) {
	w.WriteString("IFNULL(")
	if isValidIdentifier(expressionAlias[0]) == 0 {
		Quoter.WriteNameAlias(w, expressionAlias[0], "")
	} else {
		w.WriteByte('(')
		w.WriteString(expressionAlias[0])
		w.WriteByte(')')
	}
	w.WriteByte(',')
	if isValidIdentifier(expressionAlias[1]) == 0 {
		Quoter.WriteNameAlias(w, expressionAlias[1], "")
	} else {
		w.WriteByte('(')
		w.WriteString(expressionAlias[1])
		w.WriteByte(')')
	}
	w.WriteByte(')')
}

func sqlIfNullQuote4(w queryWriter, qualifierName ...string) {
	w.WriteString("IFNULL(")
	Quoter.writeQualifierName(w, qualifierName[0], qualifierName[1])
	w.WriteByte(',')
	Quoter.writeQualifierName(w, qualifierName[2], qualifierName[3])
	w.WriteByte(')')
}

// SQLIf writes a SQL IF() expression.
//		IF(expr1,expr2,expr3)
// If expr1 is TRUE (expr1 <> 0 and expr1 <> NULL) then IF() returns expr2;
// otherwise it returns expr3. IF() returns a numeric or string value, depending
// on the context in which it is used.
func SQLIf(expression, true, false string) string {
	return "IF((" + expression + "), " + true + ", " + false + ")"
}

// SQLCase generates a CASE ... WHEN ... THEN ... ELSE ... END statement.
// `value` argument can be empty. defaultValue used in the ELSE part can also be
// empty and then won't get written. `compareResult` must be a balanced sliced
// where index `i` represents the case part and index `i+1` the result.
// If the slice is imbalanced the function assumes that the last item of compareResult
// should be printed as an alias.
// https://dev.mysql.com/doc/refman/5.7/en/control-flow-functions.html#operator_case
func SQLCase(value, defaultValue string, compareResult ...string) string {
	if len(compareResult) == 1 {
		return "<SQLCase error len(compareResult) == 1>"
	}
	buf := bufferpool.Get()
	defer bufferpool.Put(buf)

	useAlias := len(compareResult)%2 == 1

	lcr := len(compareResult)
	if useAlias {
		lcr--
		buf.WriteByte('(')
	}
	buf.WriteString("CASE ")
	buf.WriteString(value)
	for i := 0; i < lcr; i = i + 2 {
		buf.WriteString(" WHEN ")
		buf.WriteString(compareResult[i])
		buf.WriteString(" THEN ")
		buf.WriteString(compareResult[i+1])
	}
	if defaultValue != "" {
		buf.WriteString(" ELSE ")
		buf.WriteString(defaultValue)
	}
	buf.WriteString(" END")
	if useAlias {
		buf.WriteByte(')')
		buf.WriteString(" AS ")
		Quoter.writeName(buf, compareResult[len(compareResult)-1])
	}
	return buf.String()
}
