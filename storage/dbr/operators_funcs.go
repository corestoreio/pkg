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
	"github.com/corestoreio/csfw/util/bufferpool"
	"github.com/corestoreio/errors"
)

// SQLIfNull creates an IFNULL expression. Argument count can be either 1, 2 or
// 4. A single expression can contain a qualified or unqualified identifier. See
// the examples.
//
// IFNULL(expr1,expr2) If expr1 is not NULL, IFNULL() returns expr1; otherwise
// it returns expr2. IFNULL() returns a numeric or string value, depending on
// the context in which it is used.
func SQLIfNull(expression ...string) expressions {
	buf := bufferpool.Get()
	defer bufferpool.Put(buf)

	switch len(expression) {
	case 1:
		sqlIfNullQuote2(buf, expression[0], "NULL ")
	case 2:
		// Input:  dbr.SQLIfNull(col1,col2)
		// Output: IFNULL(`col1`, `col2`)

		// Input:  dbr.SQLIfNull(expr1,expr2)
		// Output: IFNULL(expr1, expr1)
		sqlIfNullQuote2(buf, expression...)
	case 4:
		// Input:  dbr.SQLIfNull(table1,col1,table2,col2)
		// Output: IFNULL(`table1`.`col1`, `table2`.`col2`)
		sqlIfNullQuote4(buf, expression...)
	default:
		panic(errors.NewNotValidf("[dbr] Invalid number of arguments. Max 4 arguments allowed, got: %v", expression))

	}
	return []string{buf.String()}
}

func sqlIfNullQuote2(w queryWriter, expressionAlias ...string) {
	w.WriteString("IFNULL(")
	if isValidIdentifier(expressionAlias[0]) == 0 {
		Quoter.WriteIdentifier(w, expressionAlias[0])
	} else {
		w.WriteByte('(')
		w.WriteString(expressionAlias[0])
		w.WriteByte(')')
	}
	w.WriteByte(',')
	if isValidIdentifier(expressionAlias[1]) == 0 {
		Quoter.WriteIdentifier(w, expressionAlias[1])
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
// Returns a []string.
func SQLIf(expression, true, false string) expressions {
	ret := [...]string{"IF((", expression, "), ", true, ", ", false, ")"}
	return ret[:]
}

// SQLCase generates a CASE ... WHEN ... THEN ... ELSE ... END statement.
// `value` argument can be empty. defaultValue used in the ELSE part can also be
// empty and then won't get written. `compareResult` must be a balanced sliced
// where index `i` represents the case part and index `i+1` the result.
// If the slice is imbalanced the function assumes that the last item of compareResult
// should be printed as an alias.
// https://dev.mysql.com/doc/refman/5.7/en/control-flow-functions.html#operator_case
func SQLCase(value, defaultValue string, compareResult ...string) expressions {
	if len(compareResult) == 1 {
		return []string{"<SQLCase error len(compareResult) == 1>"}
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
		Quoter.quote(buf, compareResult[len(compareResult)-1])
	}
	return []string{buf.String()}
}
