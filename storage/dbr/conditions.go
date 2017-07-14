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
	"database/sql/driver"
	"strings"
	"time"

	"github.com/corestoreio/errors"
)

const (
	logicalAnd byte = 'a'
	logicalOr  byte = 'o'
	logicalXor byte = 'x'
	logicalNot byte = 'n'
)

// whereFragment implements a single WHERE condition. Please use the helper
// functions instead of using this type directly.
type whereFragment struct {
	// Condition can contain either a valid identifier or an expression. Set
	// field `IsExpression` to true to avoid quoting of the `Column` field.
	// Condition can also contain `qualifier.identifier`.
	Condition string
	// Operator contains the comparison logic like LIKE, IN, GREATER, etc ...
	Operator  Op
	Argument  Argument // Either this or the slice is set.
	Arguments Arguments

	Sub struct {
		// Select adds a sub-select to the where statement. Column must be
		// either a column name or anything else which can handle the result of
		// a sub-select.
		Select   *Select
		Operator Op
	}
	// Logical states how multiple WHERE statements will be connected.
	// Default to AND. Possible values are a=AND, o=OR, x=XOR, n=NOT
	Logical byte
	// IsExpression set to true if the Column contains an expression.
	// Otherwise the `Column` gets always quoted.
	IsExpression bool
	// IsPlaceHolder true if the current WHERE condition just acts as a place
	// holder for a prepared statement or an interpolation.
	IsPlaceHolder bool
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

// And sets the logical AND operator
func (wf *whereFragment) And() *whereFragment {
	wf.Logical = logicalAnd
	return wf
}

// Or sets the logical OR operator
func (wf *whereFragment) Or() *whereFragment {
	wf.Logical = logicalOr
	return wf
}

// Conditions iterates over each WHERE fragment and assembles all conditions
// into a new slice.
func (wfs WhereFragments) Conditions() []string {
	// this calculates the intersection of the columns in WhereFragments which
	// already have an argument provided/assigned and those where the arguments
	// must be assembled from the interface ArgumentAssembler. If the arguments
	// should be assembled from the interface IsPlaceHolder is true.
	c := make([]string, 0, len(wfs))
	for _, w := range wfs {
		if w.IsPlaceHolder {
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
func Using(columns ...string) *whereFragment {
	return &whereFragment{
		Using: columns, // gets quoted during writing the query in ToSQL
	}
}

// SubSelect creates a condition for a WHERE or JOIN statement to compare the
// data in `rawStatementOrColumnName` with the returned value/s of the
// sub-select.
func SubSelect(columnName string, operator Op, s *Select) *whereFragment {
	wf := &whereFragment{
		Condition: columnName,
	}
	wf.Sub.Select = s
	wf.Sub.Operator = operator
	return wf
}

// Column adds a condition to a WHERE or HAVING statement.
func Column(columnName string) *whereFragment {
	return &whereFragment{
		Condition: columnName,
	}
}

func (wf *whereFragment) Null() *whereFragment {
	wf.Operator = Null
	return wf
}
func (wf *whereFragment) NotNull() *whereFragment {
	wf.Operator = NotNull
	return wf
}
func (wf *whereFragment) In() *whereFragment {
	wf.Operator = In
	return wf
}
func (wf *whereFragment) NotIn() *whereFragment {
	wf.Operator = NotIn
	return wf
}
func (wf *whereFragment) Between() *whereFragment {
	wf.Operator = Between
	return wf
}
func (wf *whereFragment) NotBetween() *whereFragment {
	wf.Operator = NotBetween
	return wf
}
func (wf *whereFragment) Like() *whereFragment {
	wf.Operator = Like
	return wf
}
func (wf *whereFragment) NotLike() *whereFragment {
	wf.Operator = NotLike
	return wf
}
func (wf *whereFragment) Greatest() *whereFragment {
	wf.Operator = Greatest
	return wf
}
func (wf *whereFragment) Least() *whereFragment {
	wf.Operator = Least
	return wf
}
func (wf *whereFragment) Equal() *whereFragment {
	wf.Operator = Equal
	return wf
}
func (wf *whereFragment) NotEqual() *whereFragment {
	wf.Operator = NotEqual
	return wf
}
func (wf *whereFragment) Exists() *whereFragment {
	wf.Operator = Exists
	return wf
}
func (wf *whereFragment) NotExists() *whereFragment {
	wf.Operator = NotExists
	return wf
}
func (wf *whereFragment) Less() *whereFragment {
	wf.Operator = Less
	return wf
}
func (wf *whereFragment) Greater() *whereFragment {
	wf.Operator = Greater
	return wf
}
func (wf *whereFragment) LessOrEqual() *whereFragment {
	wf.Operator = LessOrEqual
	return wf
}
func (wf *whereFragment) GreaterOrEqual() *whereFragment {
	wf.Operator = GreaterOrEqual
	return wf
}

func (wf *whereFragment) Regexp() *whereFragment {
	wf.Operator = Regexp
	return wf
}

func (wf *whereFragment) NotRegexp() *whereFragment {
	wf.Operator = NotRegexp
	return wf
}

func (wf *whereFragment) Xor() *whereFragment {
	wf.Operator = Xor
	return wf
}

func (wf *whereFragment) SpaceShip() *whereFragment {
	wf.Operator = SpaceShip
	return wf
}

func (wf *whereFragment) Coalesce() *whereFragment {
	wf.Operator = Coalesce
	return wf
}

///////////////////////////////////////////////////////////////////////////////////
//		TYPES
///////////////////////////////////////////////////////////////////////////////////

// Cahen's Constant, used as a random identifier
const cahensConstant = -64341

// PlaceHolder sets the database specific place holder character. Mostly used in
// prepared statements and for interpolation.
func (wf *whereFragment) PlaceHolder() *whereFragment {
	wf.Argument = argPlaceHolder(cahensConstant)
	wf.IsPlaceHolder = true
	return wf
}

func (wf *whereFragment) Int(i int) *whereFragment {
	wf.Argument = ArgInt(i)
	return wf
}

func (wf *whereFragment) Ints(i ...int) *whereFragment {
	wf.Argument = ArgInts(i)
	return wf
}
func (wf *whereFragment) Int64(i int64) *whereFragment {
	wf.Argument = ArgInt64(i)
	return wf
}

func (wf *whereFragment) Int64s(i ...int64) *whereFragment {
	wf.Argument = ArgInt64s(i)
	return wf
}

func (wf *whereFragment) Float64(f float64) *whereFragment {
	wf.Argument = ArgFloat64(f)
	return wf
}
func (wf *whereFragment) Float64s(f ...float64) *whereFragment {
	wf.Argument = ArgFloat64s(f)
	return wf
}
func (wf *whereFragment) String(s string) *whereFragment {
	wf.Argument = ArgString(s)
	return wf
}

func (wf *whereFragment) Strings(s ...string) *whereFragment {
	wf.Argument = ArgStrings(s)
	return wf
}

// Bool uses bool values for comparison.
func (wf *whereFragment) Bool(b bool) *whereFragment {
	wf.Argument = ArgBool(b)
	return wf
}

// Bools uses bool values for comparison.
func (wf *whereFragment) Bools(b ...bool) *whereFragment {
	wf.Argument = ArgBools(b)
	return wf
}

// Bytes uses a byte slice for comparison. Providing a nil argument returns a
// NULL type. Detects between valid UTF-8 strings and binary data. Later gets
// hex encoded.
func (wf *whereFragment) Bytes(p []byte) *whereFragment {
	wf.Argument = ArgBytes(p)
	return wf
}

func (wf *whereFragment) BytesSlice(p ...[]byte) *whereFragment {
	wf.Argument = ArgBytesSlice(p)
	return wf
}

// Time uses time.Time values for comparison.
func (wf *whereFragment) Time(t time.Time) *whereFragment {
	wf.Argument = ArgTime(t)
	return wf
}

// Times uses time.Time values for comparison.
func (wf *whereFragment) Times(t ...time.Time) *whereFragment {
	wf.Argument = ArgTimes(t)
	return wf
}

// NullString uses nullable string values for comparison.
func (wf *whereFragment) NullString(values ...NullString) *whereFragment {
	if len(values) == 1 {
		wf.Argument = values[0]
	} else {
		wf.Argument = ArgNullStrings(values)
	}
	return wf
}

// NullFloat64 uses nullable float64 values for comparison.
func (wf *whereFragment) NullFloat64(values ...NullFloat64) *whereFragment {
	if len(values) == 1 {
		wf.Argument = values[0]
	} else {
		wf.Argument = ArgNullFloat64s(values)
	}
	return wf
}

// NullInt64 uses nullable int64 values for comparison.
func (wf *whereFragment) NullInt64(values ...NullInt64) *whereFragment {
	if len(values) == 1 {
		wf.Argument = values[0]
	} else {
		wf.Argument = ArgNullInt64s(values)
	}
	return wf
}

// NullBool uses nullable bool values for comparison.
func (wf *whereFragment) NullBool(value NullBool) *whereFragment {
	wf.Argument = value
	return wf
}

// NullTime uses nullable time values for comparison.
func (wf *whereFragment) NullTime(values ...NullTime) *whereFragment {
	if len(values) == 1 {
		wf.Argument = values[0]
	} else {
		wf.Argument = ArgNullTimes(values)
	}
	return wf
}

// Value uses driver.Valuers for comparison.
func (wf *whereFragment) Value(values ...driver.Valuer) *whereFragment {
	wf.Argument = ArgValues(values)
	return wf
}

// Expression adds an unquoted SQL expression to a WHERE or HAVING statement.
func Expression(expression string, arg ...Argument) *whereFragment {
	return &whereFragment{
		IsExpression: true,
		Condition:    expression,
		Arguments:    arg,
	}
}

// ParenthesisOpen sets an open parenthesis "(". Mostly used for OR conditions
// in combination with AND conditions.
func ParenthesisOpen() *whereFragment {
	return &whereFragment{
		Condition: "(",
	}
}

// ParenthesisClose sets a closing parenthesis ")". Mostly used for OR
// conditions in combination with AND conditions.
func ParenthesisClose() *whereFragment {
	return &whereFragment{
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
			if strings.IndexByte(f.Condition, '?') == -1 && (len(f.Arguments) == 1 || f.Argument != nil) && f.Operator > 0 {
				eArg := f.Argument
				if eArg == nil {
					eArg = f.Arguments[0]
				}
				writeOperator(w, eArg.len(), f.Operator)
			}

		case f.Sub.Select != nil:
			Quoter.WriteNameAlias(w, f.Condition, "")
			writeOperator(w, 0, f.Sub.Operator)
			w.WriteByte('(')
			if err := f.Sub.Select.toSQL(w); err != nil {
				return errors.Wrapf(err, "[dbr] write failed SubSelect for table: %q", f.Sub.Select.Table.String())
			}
			w.WriteByte(')')

		case f.Argument != nil && f.Arguments == nil:
			Quoter.WriteNameAlias(w, f.Condition, "")
			al := f.Argument.len()
			if f.IsPlaceHolder {
				al = 1
			}
			if al > 1 && f.Operator == 0 { // no operator but slice applied, so creating an IN query.
				f.Operator = In
			}
			writeOperator(w, al, f.Operator)

		case f.Argument == nil && f.Arguments == nil:
			Quoter.WriteNameAlias(w, f.Condition, "")
			c := f.Operator
			if c == 0 {
				c = Null
			}
			writeOperator(w, 1, c)

		default:
			panic(errors.NewNotSupportedf("[dbr] Multiple arguments for a column are not supported\nWhereFragment: %#v\n", f))
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
		case f.IsPlaceHolder:
			addArg = writeOperator(backHole{}, 1, f.Operator) // always a length of one, see the `repeat()` function
			// By keeping addArg as it is and not setting
			// addArg=false, this []int avoids
			// https://en.wikipedia.org/wiki/Permutation Which would
			// result in a Go Code like
			// https://play.golang.org/p/rZvW0qW1N7 (C) Volker Dobler
			// Because addArg=false does not add below the arguments and we must
			// later swap the positions.
			pendingArgPos = append(pendingArgPos, pendingArgPosCount)

		case f.Argument != nil:
			// a column only supports one argument.
			addArg = writeOperator(backHole{}, f.Argument.len(), f.Operator)
		case f.Sub.Select != nil:
			args, err = f.Sub.Select.appendArgs(args)
			if err != nil {
				return nil, pendingArgPos, errors.Wrapf(err, "[dbr] write failed SubSelect for table: %q", f.Sub.Select.Table.String())
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
