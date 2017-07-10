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
	"strings"

	"github.com/corestoreio/errors"
)

const (
	logicalAnd byte = 'a'
	logicalOr  byte = 'o'
	logicalXor byte = 'x'
	logicalNot byte = 'n'
)

// WhereFragment implements a single WHERE condition. Please use the helper
// functions instead of using this type directly.
type WhereFragment struct {
	// Condition can contain either a valid identifier or an expression. Set
	// field `IsExpression` to true to avoid quoting of the `Column` field.
	// Condition can also contain `qualifier.identifier`.
	Condition string
	// Operator  Op todo rethink that
	Argument  Argument // Either this or the slice is set.
	Arguments Arguments

	Sub struct {
		// Select adds a sub-select to the where statement. Column must be either
		// a column name or anything else which can handle the result of a
		// sub-select.
		Select   *Select
		Operator Op
	}
	// Logical states how multiple WHERE statements will be connected.
	// Default to AND. Possible values are a=AND, o=OR, x=XOR, n=NOT
	Logical byte
	// IsExpression set to true if the Column contains an expression.
	// Otherwise the `Column` gets always quoted.
	IsExpression bool
	// Using a list of column names which get quoted during SQL statement
	// creation.
	Using []string
}

// WhereFragments provides a list WHERE resp. ON clauses.
type WhereFragments []*WhereFragment

// JoinFragments defines multiple join conditions.
type JoinFragments []*joinFragment

type joinFragment struct {
	// JoinType can be LEFT, RIGHT, INNER, OUTER or CROSS
	JoinType string
	// Table name and alias of the table
	Table alias
	// OnConditions join on those conditions
	OnConditions WhereFragments
}

// And sets the logical AND operator
func (wf *WhereFragment) And() *WhereFragment {
	wf.Logical = logicalAnd
	return wf
}

// Or sets the logical OR operator
func (wf *WhereFragment) Or() *WhereFragment {
	wf.Logical = logicalOr
	return wf
}

// cahensConstant a randomly chosen number aka unique identifier to check if
// some conditions are only placeholder in combination with the interface
// ArgumentAssembler. A bit hacky ...
const cahensConstant = -6434105

// Conditions iterates over each WHERE fragment and assembles all conditions
// into a new slice.
func (wf WhereFragments) Conditions() []string {
	// this calculates the intersection of the columns in WhereFragments which
	// already have an argument provided and those where the arguments must be
	// assembled from the interface ArgumentAssembler. If the arguments should
	// be assembled from the interface Argument.len() returns cahensConstant.
	c := make([]string, 0, len(wf))
	for _, w := range wf {
		if w.Argument != nil && w.Argument.len() == cahensConstant {
			c = append(c, w.Condition)
		}
		if len(w.Arguments) > 0 && w.Arguments[0].len() == cahensConstant {
			c = append(c, w.Condition)
		}
	}
	return c
}

// Using add syntactic sugar to a JOIN statement: The USING(column_list) clause
// names a list of columns that must exist in both tables. If tables a and b
// both contain columns c1, c2, and c3, the following join compares
// corresponding columns from the two tables:
//	a LEFT JOIN b USING (c1, c2, c3)
// The columns list gets quoted while writing the query string.
func Using(columns ...string) *WhereFragment {
	return &WhereFragment{
		Using: columns, // gets quoted during writing the query in ToSQL
	}
}

// SubSelect creates a condition for a WHERE or JOIN statement to compare the
// data in `rawStatementOrColumnName` with the returned value/s of the
// sub-select.
func SubSelect(columnName string, operator Op, s *Select) *WhereFragment {
	wf := &WhereFragment{
		Condition: columnName,
	}
	wf.Sub.Select = s
	wf.Sub.Operator = operator
	return wf
}

// TODO(CyS) implement once needed
//func WithSelect(columnName string, operator Op, wth *With) *WhereFragment {
//	wf := &WhereFragment{
//		Condition: columnName,
//	}
//	wf.Sub.With = wth
//	wf.Sub.Operator = operator
//	return wf
//}

// Column adds a condition to a WHERE or HAVING statement.
func Column(columnName string, arg Argument) *WhereFragment {
	return &WhereFragment{
		Condition: columnName,
		Argument:  arg,
	}
}

// Expression adds an unquoted SQL expression to a WHERE or HAVING statement.
func Expression(expression string, arg ...Argument) *WhereFragment {
	return &WhereFragment{
		IsExpression: true,
		Condition:    expression,
		Arguments:    arg,
	}
}

// ParenthesisOpen sets an open parenthesis "(". Mostly used for OR conditions
// in combination with AND conditions.
func ParenthesisOpen() *WhereFragment {
	return &WhereFragment{
		Condition: "(",
	}
}

// ParenthesisClose sets a closing parenthesis ")". Mostly used for OR
// conditions in combination with AND conditions.
func ParenthesisClose() *WhereFragment {
	return &WhereFragment{
		Condition: ")",
	}
}

// conditionType enum of j=join, w=where, h=having
func (wf WhereFragments) write(w queryWriter, conditionType byte) error {
	if len(wf) == 0 {
		return nil
	}

	switch conditionType {
	case 'w':
		w.WriteString(" WHERE ")
	case 'h':
		w.WriteString(" HAVING ")
	}

	i := 0
	for _, f := range wf {

		if conditionType == 'j' {
			if len(f.Using) > 0 {
				w.WriteString(" USING (")
				for j, c := range f.Using {
					if j > 0 {
						w.WriteByte(',')
					}
					Quoter.writeName(w, c)
				}
				w.WriteByte(')')
				return nil // done, only one USING allowed
			}
			if i == 0 {
				w.WriteString(" ON ")
			}
		}

		if f.Condition == ")" {
			w.WriteString(f.Condition)
			continue
		}

		if i > 0 {
			// How the WHERE conditions are connected
			switch f.Logical {
			case logicalAnd:
				w.WriteString(" AND ")
			case logicalOr:
				w.WriteString(" OR ")
			case logicalXor:
				w.WriteString(" XOR ")
			case logicalNot:
				w.WriteString(" NOT ")
			default:
				w.WriteString(" AND ")
			}
		}

		if f.Condition == "(" {
			i = 0
			w.WriteString(f.Condition)
			continue
		}

		w.WriteByte('(')

		switch {
		case f.IsExpression:
			_, _ = w.WriteString(f.Condition)
			// Only write the operator in case there is no place holder and we have one argument
			if strings.IndexByte(f.Condition, '?') == -1 && len(f.Arguments) == 1 && f.Arguments[0].operator() > 0 {
				writeOperator(w, true, f.Arguments[0])
			}
			// TODO: case f.Sub.Select != nil:
		default:
			Quoter.WriteNameAlias(w, f.Condition, "")

			if f.Sub.Select != nil {
				writeOperator(w, false, argNull(f.Sub.Operator))
				w.WriteByte('(')
				if err := f.Sub.Select.toSQL(w); err != nil {
					return errors.Wrapf(err, "[dbr] write failed SubSelect for table: %q", f.Sub.Select.Table.String())
				}
				w.WriteByte(')')
			} else {
				// a column only supports one argument.
				if f.Argument != nil && f.Arguments == nil {
					writeOperator(w, true, f.Argument)
				}
				if f.Argument == nil && f.Arguments == nil {
					writeOperator(w, true, ArgNull())
				}
			}
		}

		w.WriteByte(')')
		i++
	}
	return nil
}

// conditionType enum of j=join, w=where, h=having
func (wf WhereFragments) appendArgs(args Arguments, conditionType byte) (_ Arguments, pendingArgPos []int, err error) {
	if len(wf) == 0 {
		return args, pendingArgPos, nil
	}

	pendingArgPosCount := len(args)
	i := 0
	for _, f := range wf {

		switch {
		case conditionType == 'j' && len(f.Using) > 0:
			return args, pendingArgPos, nil // done, only one USING allowed

		case f.Condition == ")":
			continue

		case f.Condition == "(":
			i = 0
			continue
		}

		addArg := false
		switch {
		case f.IsExpression:
			addArg = true
		case f.Sub.Select != nil:
			args, err = f.Sub.Select.appendArgs(args)
			if err != nil {
				return nil, pendingArgPos, errors.Wrapf(err, "[dbr] write failed SubSelect for table: %q", f.Sub.Select.Table.String())
			}
		case f.Argument != nil:
			// a column only supports one argument.
			addArg = writeOperator(backHole{}, true, f.Argument)
			if f.Argument.len() == cahensConstant {
				// By keeping addArg as it is and not setting
				// addArg=false, this []int avoids
				// https://en.wikipedia.org/wiki/Permutation Which would
				// result in a Go Code like
				// https://play.golang.org/p/rZvW0qW1N7 (C) Volker Dobler
				// Because addArg=false does not add below the arguments and we must
				// later swap the positions.
				pendingArgPos = append(pendingArgPos, pendingArgPosCount)
			}
		}

		if addArg {
			if f.Argument != nil {
				args = append(args, f.Argument)
			}
			args = append(args, f.Arguments...)
		}
		pendingArgPosCount++
		i++
	}
	return args, pendingArgPos, nil
}

func appendAssembledArgs(pendingArgPos []int, rec ArgumentAssembler, args Arguments, stmtType int, columns []string) (_ Arguments, err error) {
	if rec == nil {
		return args, nil
	}

	lenBefore := len(args)
	args, err = rec.AssembleArguments(stmtType, args, columns)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	lenAfter := len(args)
	if lenAfter > lenBefore {
		j := 0
		newLen := lenAfter - len(pendingArgPos)
		for i := newLen; i < lenAfter; i++ {
			args[pendingArgPos[j]], args[i] = args[i], args[pendingArgPos[j]]
			j++
		}
		args = args[:newLen] // remove the appended argPlaceHolder types after swapping
	}
	return args, nil
}
