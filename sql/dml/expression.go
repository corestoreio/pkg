// Copyright 2015-present, Cyrill @ Schumacher.fm and the CoreStore contributors
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

package dml

import (
	"bytes"
	"strings"

	"github.com/corestoreio/errors"
	"github.com/corestoreio/pkg/util/bufferpool"
)

// writeExpression writes the strings into `w` and correctly handles the place holder
// repetition depending on the number of arguments.
func writeExpression(w *bytes.Buffer, expression string, args []any) (phCount int, err error) {
	phCount = strings.Count(expression, placeHolderStr)
	if phCount == 0 || len(args) == 0 {
		// fast path
		_, err = w.WriteString(expression)
	} else {
		err = writeInterpolate(w, expression, args)
	}
	return
}

// SQLIfNull creates an IFNULL expression. Argument count can be either 1, 2 or
// 4. A single expression can contain a qualified or unqualified identifier. See
// the examples.
//
// IFNULL(expr1,expr2) If expr1 is not NULL, IFNULL() returns expr1; otherwise
// it returns expr2. IFNULL() returns a numeric or string value, depending on
// the context in which it is used.
func SQLIfNull(expression ...string) *Condition {
	return &Condition{
		Left:             sqlIfNull(expression),
		IsLeftExpression: true,
	}
}

func sqlIfNull(expression []string) string {
	buf := bufferpool.Get() // way faster than strings.Builder

	switch len(expression) {
	case 1:
		sqlIfNullQuote2(buf, expression[0], sqlStrNullUC)
	case 2:
		// Input:  dml.SQLIfNull(col1,col2)
		// Output: IFNULL(`col1`, `col2`)

		// Input:  dml.SQLIfNull(expr1,expr2)
		// Output: IFNULL(expr1, expr1)
		sqlIfNullQuote2(buf, expression...)
	case 4:
		// Input:  dml.SQLIfNull(table1,col1,table2,col2)
		// Output: IFNULL(`table1`.`col1`, `table2`.`col2`)
		sqlIfNullQuote4(buf, expression...)
	default:
		panic(errors.NotValid.Newf("[dml] Invalid number of arguments. Max 4 arguments allowed, got: %v", expression))

	}
	ret := buf.String()
	bufferpool.Put(buf)
	return ret
}

func sqlIfNullQuote2(w *bytes.Buffer, expressionAlias ...string) {
	w.WriteString("IFNULL(")
	if isValidIdentifier(expressionAlias[0]) == 0 {
		Quoter.WriteIdentifier(w, expressionAlias[0])
	} else {
		w.WriteString(expressionAlias[0])
	}
	w.WriteByte(',')
	if isValidIdentifier(expressionAlias[1]) == 0 {
		Quoter.WriteIdentifier(w, expressionAlias[1])
	} else {
		w.WriteString(expressionAlias[1])
	}
	w.WriteByte(')')
}

func sqlIfNullQuote4(w *bytes.Buffer, qualifierName ...string) {
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
func SQLIf(expression, true, false string) *Condition {
	return &Condition{
		Left:             "IF((" + expression + "), " + true + ", " + false + ")",
		IsLeftExpression: 2 > 1,
	}
}

// SQLCase generates a CASE ... WHEN ... THEN ... ELSE ... END statement.
// `value` argument can be empty. defaultValue used in the ELSE part can also be
// empty and then won't get written. `compareResult` must be a balanced sliced
// where index `i` represents the case part and index `i+1` the result.
// If the slice is imbalanced the function assumes that the last item of compareResult
// should be printed as an alias.
// https://dev.mysql.com/doc/refman/5.7/en/control-flow-functions.html#operator_case
func SQLCase(value, defaultValue string, compareResult ...string) *Condition {
	return &Condition{
		Left:             sqlCase(value, defaultValue, compareResult...),
		IsLeftExpression: true,
	}
}

func sqlCase(value, defaultValue string, compareResult ...string) string {
	if len(compareResult) < 2 {
		panic(errors.Fatal.Newf("[dml] SQLCase error incorrect length for compareResult: %v", compareResult))
	}
	buf := bufferpool.Get()

	useAlias := len(compareResult)%2 == 1

	lcr := len(compareResult)
	if useAlias {
		lcr--
		buf.WriteByte('(')
	}
	buf.WriteString("CASE ")
	buf.WriteString(value)
	for i := 0; i < lcr; i += 2 {
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
	e := buf.String()
	bufferpool.Put(buf)
	return e
}
