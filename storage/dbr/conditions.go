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

import "github.com/corestoreio/errors"

const (
	logicalAnd byte = 'a'
	logicalOr  byte = 'o'
	logicalXor byte = 'x'
	logicalNot byte = 'n'
)

type whereFragment struct {
	// Condition can contain either a valid identifier or an expression. Set
	// field `IsExpression` to true to avoid quoting of the `Column` field.
	Condition string
	Arguments Arguments
	Sub       struct {
		// Select adds a sub-select to the where statement. Column must be either
		// a column name or anything else which can handle the result of a
		// sub-select.
		Select   *Select
		Operator rune
	}
	// Logical states how multiple where statements will be connected.
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
type WhereFragments []*whereFragment

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

// ConditionArg used at argument in Where()
type ConditionArg interface {
	appendConditions(WhereFragments) WhereFragments
	And() ConditionArg // And connects next condition via AND
	Or() ConditionArg  // Or connects next condition via OR
}

// Eq is a map Expression -> value pairs which must be matched in a query.
// Joined at AND statements to the WHERE clause. Implements ConditionArg
// interface. Eq = EqualityMap.
type Eq map[string]Argument

func (eq Eq) appendConditions(wfs WhereFragments) WhereFragments {
	for c, arg := range eq {
		if arg == nil {
			arg = ArgNull()
		}
		wfs = append(wfs, &whereFragment{
			Condition: c,
			Arguments: Arguments{arg},
		})
	}
	return wfs
}

// And sets the logical AND operator. Default case.
func (eq Eq) And() ConditionArg {
	return eq
}

// Or not supported
func (eq Eq) Or() ConditionArg {
	return eq
}

func (wf *whereFragment) appendConditions(wfs WhereFragments) WhereFragments {
	return append(wfs, wf)
}

// And sets the logical AND operator
func (wf *whereFragment) And() ConditionArg {
	wf.Logical = logicalAnd
	return wf
}

// Or sets the logical OR operator
func (wf *whereFragment) Or() ConditionArg {
	wf.Logical = logicalOr
	return wf
}

// Conditions iterates over each WHERE fragment and assembles all conditions
// into a new slice.
func (wf WhereFragments) Conditions() []string {
	c := make([]string, len(wf))
	for i, w := range wf {
		c[i] = w.Condition
	}
	return c
}

// Using add syntactic sugar to a JOIN statement: The USING(column_list) clause
// names a list of columns that must exist in both tables. If tables a and b
// both contain columns c1, c2, and c3, the following join compares
// corresponding columns from the two tables:
//	a LEFT JOIN b USING (c1, c2, c3)
// The columns list gets quoted while writing the query string.
func Using(columns ...string) ConditionArg {
	return &whereFragment{
		Using: columns, // gets quoted during writing the query in ToSQL
	}
}

// SubSelect creates a condition for a WHERE or JOIN statement to compare the
// data in `rawStatementOrColumnName` with the returned value/s of the
// sub-select.
func SubSelect(columnName string, operator rune, s *Select) ConditionArg {
	wf := &whereFragment{
		Condition: columnName,
	}
	wf.Sub.Select = s
	wf.Sub.Operator = operator
	return wf
}

// Column adds a condition to a WHERE or HAVING statement.
func Column(columnName string, arg ...Argument) ConditionArg {
	return &whereFragment{
		Condition: columnName,
		Arguments: arg,
	}
}

// Expression adds an unquoted SQL expression to a WHERE or HAVING statement.
func Expression(expression string, arg ...Argument) ConditionArg {
	return &whereFragment{
		IsExpression: true,
		Condition:    expression,
		Arguments:    arg,
	}
}

// ParenthesisOpen sets an open parenthesis "(". Mostly used for OR conditions
// in combination with AND conditions.
func ParenthesisOpen() ConditionArg {
	return &whereFragment{
		Condition: "(",
	}
}

// ParenthesisClose sets a closing parenthesis ")". Mostly used for OR
// conditions in combination with AND conditions.
func ParenthesisClose() ConditionArg {
	return &whereFragment{
		Condition: ")",
	}
}

func (wf WhereFragments) append(wargs ...ConditionArg) WhereFragments {
	for _, warg := range wargs {
		wf = warg.appendConditions(wf)
	}
	return wf
}

// stmtType enum of j=join, w=where, h=having
func (wf WhereFragments) write(w queryWriter, args Arguments, stmtType byte) (Arguments, error) {
	if len(wf) == 0 {
		return args, nil
	}

	switch stmtType {
	case 'w':
		w.WriteString(" WHERE ")
	case 'h':
		w.WriteString(" HAVING ")
	}

	i := 0
	for _, f := range wf {

		if stmtType == 'j' {
			if len(f.Using) > 0 {
				w.WriteString(" USING (")
				for j, c := range f.Using {
					if j > 0 {
						w.WriteByte(',')
					}
					Quoter.quote(w, c)
				}
				w.WriteByte(')')
				return args, nil // done, only one using allowed
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
		addArg := false
		if f.IsExpression {
			_, _ = w.WriteString(f.Condition)
			addArg = true
			if len(f.Arguments) == 1 && f.Arguments[0].operator() > 0 {
				writeOperator(w, f.Arguments[0].operator(), true)
			}
		} else {
			Quoter.FquoteAs(w, f.Condition)

			if f.Sub.Select != nil {
				writeOperator(w, f.Sub.Operator, false)
				w.WriteByte('(')
				subArgs, err := f.Sub.Select.toSQL(w)
				w.WriteByte(')')
				if err != nil {
					return nil, errors.Wrapf(err, "[dbr] write failed SubSelect for table: %q", f.Sub.Select.Table.String())
				}
				args = append(args, subArgs...)
			} else {
				// a column only supports one argument.
				if len(f.Arguments) == 1 {
					addArg = writeOperator(w, f.Arguments[0].operator(), true)
				}
			}
		}
		w.WriteByte(')')

		if addArg {
			args = append(args, f.Arguments...)
		}
		i++
	}
	return args, nil
}
