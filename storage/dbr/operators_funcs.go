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
	"bytes"
	"strings"
)

// SQLIfNull appends the (optional) alias to the IFNULL expression. Value
// count can be between 1-n. IFNULL(expr1,expr2) If expr1 is not NULL, IFNULL()
// returns expr1; otherwise it returns expr2. IFNULL() returns a numeric or
// string value, depending on the context in which it is used. See the examples.
// Returns a []string.
func SQLIfNull(expressionAlias ...string) expressions {
	ret := make([]string, 0, 12+len(expressionAlias))

	switch len(expressionAlias) {
	case 1:
		ret = sqlIfNullQuote2(ret, expressionAlias[0], "NULL ")
	case 2:
		// Input:  dbr.SQLIfNull(col1,col2,alias)
		// Output: IFNULL(`col1`, `col2`)

		// Input:  dbr.SQLIfNull(expr1,expr2,alias)
		// Output: IFNULL(expr1, expr1)
		ret = sqlIfNullQuote2(ret, expressionAlias...)

	case 3:
		// Input:  dbr.SQLIfNull(table1.col1,table2.col2,alias)
		// Output: IFNULL(`table1`.`col1`, `table2`.`col2`) AS `alias`

		// Input:  dbr.SQLIfNull(col1,col2,alias)
		// Output: IFNULL(`col1`, `col2`) AS `alias`

		// Input:  dbr.SQLIfNull(expr1,expr2,alias)
		// Output: IFNULL(expr1, expr1) AS `alias`
		ret = sqlIfNullQuote2(ret, expressionAlias...)
		ret = append(ret, " AS ")
		ret = Quoter.appendName(ret, expressionAlias[2])

	case 4:
		// Input:  dbr.SQLIfNull(table1,col1,table2,col2)
		// Output: IFNULL(`table1`.`col1`, `table2`.`col2`)
		ret = sqlIfNullQuote4(ret, expressionAlias...)
	case 5:
		// Input:  dbr.SQLIfNull(table1,col1,table2,col2,alias)
		// Output: IFNULL(`table1`.`col1`, `table2`.`col2`) AS `alias`
		ret = sqlIfNullQuote4(ret, expressionAlias[:4]...)
		ret = append(ret, " AS ")
		ret = Quoter.appendName(ret, expressionAlias[4])
	default:
		ret = sqlIfNullQuote4(ret, expressionAlias[:4]...)
		ret = append(ret, " AS ")
		ret = Quoter.appendName(ret, strings.Join(expressionAlias[4:], "_"))
	}
	return ret
}

func sqlIfNullQuote2(ret []string, expressionAlias ...string) []string {
	ret = append(ret, "IFNULL(")
	if isValidIdentifier(expressionAlias[0]) == 0 {
		var buf bytes.Buffer // for now this hacky way until we've create a appendNameAlias function.
		Quoter.WriteIdentifier(&buf, expressionAlias[0])
		ret = append(ret, buf.String())
	} else {
		ret = append(ret, "(", expressionAlias[0], ")")
	}
	ret = append(ret, ",")
	if isValidIdentifier(expressionAlias[1]) == 0 {
		var buf bytes.Buffer // for now this hacky way until we've create a appendNameAlias function.
		Quoter.WriteIdentifier(&buf, expressionAlias[1])
		ret = append(ret, buf.String())
	} else {
		ret = append(ret, "(", expressionAlias[1], ")")
	}
	return append(ret, ")")
}

func sqlIfNullQuote4(ret []string, qualifierName ...string) []string {
	ret = append(ret, "IFNULL(")
	ret = Quoter.appendQualifierName(ret, qualifierName[0], qualifierName[1])
	ret = append(ret, ",")
	ret = Quoter.appendQualifierName(ret, qualifierName[2], qualifierName[3])
	return append(ret, ")")
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
// where index `i` represents the case part and index `i+1` the result. If the
// `compareResult` slice is imbalanced the function assumes that the last item
// of compareResult should be printed as an alias. The returned slice must be
// joined or written into a buffer by the caller of this function. Returns a
// []string.
// https://dev.mysql.com/doc/refman/5.7/en/control-flow-functions.html#operator_case
func SQLCase(value, defaultValue string, compareResult ...string) expressions {
	if len(compareResult) == 1 {
		return []string{"<SQLCase error len(compareResult) == 1>"}
	}
	ret := make([]string, 0, 10+len(compareResult))

	useAlias := len(compareResult)%2 == 1

	lcr := len(compareResult)
	if useAlias {
		lcr--
		ret = append(ret, "(")
	}
	ret = append(ret, "CASE ", value)
	for i := 0; i < lcr; i = i + 2 {
		ret = append(ret, " WHEN ", compareResult[i], " THEN ", compareResult[i+1])
	}
	if defaultValue != "" {
		ret = append(ret, " ELSE ", defaultValue)
	}
	ret = append(ret, " END")
	if useAlias {
		ret = append(ret, ") AS ")
		ret = Quoter.appendName(ret, compareResult[len(compareResult)-1])
	}
	return ret
}
