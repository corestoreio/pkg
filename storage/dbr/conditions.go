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

// Comparison functions and operators describe all available possibilities.
const (
	Null           Op = 'n'          // IS NULL
	NotNull        Op = 'N'          // IS NOT NULL
	In             Op = '∈'          // IN ?
	NotIn          Op = '∉'          // NOT IN ?
	Between        Op = 'b'          // BETWEEN ? AND ?
	NotBetween     Op = 'B'          // NOT BETWEEN ? AND ?
	Like           Op = 'l'          // LIKE ?
	NotLike        Op = 'L'          // NOT LIKE ?
	Greatest       Op = '≫'          // GREATEST(?,?,?) returns NULL if any value is NULL.
	Least          Op = '≪'          // LEAST(?,?,?) If any value is NULL, the result is NULL.
	Equal          Op = '='          // = ?
	NotEqual       Op = '≠'          // != ?
	Exists         Op = '∃'          // EXISTS(subquery)
	NotExists      Op = '∄'          // NOT EXISTS(subquery)
	Less           Op = '<'          // <
	Greater        Op = '>'          // >
	LessOrEqual    Op = '≤'          // <=
	GreaterOrEqual Op = '≥'          // >=
	Regexp         Op = 'r'          // REGEXP ?
	NotRegexp      Op = 'R'          // NOT REGEXP ?
	Xor            Op = '⊻'          // XOR ?
	SpaceShip      Op = '\U0001f680' // a <=> b is equivalent to a = b OR (a IS NULL AND b IS NULL) NULL-safe equal to operator
	Coalesce       Op = 'c'          // Returns the first non-NULL value in the list, or NULL if there are no non-NULL values.
)

// Op the Operator, defines comparison and operator functions used in any
// conditions. The upper case letter always negates.
// https://dev.mysql.com/doc/refman/5.7/en/comparison-operators.html
// https://mariadb.com/kb/en/mariadb/comparison-operators/
type Op rune

// String transforms the rune into a string.
func (o Op) String() string {
	if o < 1 {
		o = Equal
	}
	return string(o)
}

// WhereFragment implements a single WHERE condition. Please use the helper
// functions instead of using this type directly.
type WhereFragment struct {
	// Condition can contain either a valid identifier or an expression. Set
	// field `IsExpression` to true to avoid quoting of the `Column` field.
	// Condition can also contain `qualifier.identifier`.
	Condition string
	// Operator contains the comparison logic like LIKE, IN, GREATER, etc ...
	Operator  Op
	Argument  Value // Either this or the slice is set.
	Arguments Values

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

// intersectConditions iterates over each WHERE fragment and appends all
// conditions aka column names to the slice c.
func (wfs WhereFragments) intersectConditions(c []string) []string {
	// this calculates the intersection of the columns in WhereFragments which
	// already have an value provided/assigned and those where the values
	// must be assembled from the interface ValuesAppender. If the values
	// should be assembled from the interface IsPlaceHolder is true.
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
//		a LEFT JOIN b USING (c1, c2, c3)
// The columns list gets quoted while writing the query string.
func Using(columns ...string) *WhereFragment {
	return &WhereFragment{
		Using: columns, // gets quoted during writing the query in ToSQL
	}
}

// SubSelect creates a condition for a WHERE or JOIN statement to compare the
// data in `columnName` with the returned value/s of the sub-select. Choose the
// appropriate comparison operator. Equal comparions operator is the default
// one.
func SubSelect(columnName string, operator Op, s *Select) *WhereFragment {
	wf := &WhereFragment{
		Condition: columnName,
	}
	wf.Sub.Select = s
	wf.Sub.Operator = operator
	return wf
}

// Column adds a condition to a WHERE or HAVING statement.
func Column(columnName string) *WhereFragment {
	return &WhereFragment{
		Condition: columnName,
	}
}

func (wf *WhereFragment) Null() *WhereFragment {
	wf.Operator = Null
	return wf
}
func (wf *WhereFragment) NotNull() *WhereFragment {
	wf.Operator = NotNull
	return wf
}
func (wf *WhereFragment) In() *WhereFragment {
	wf.Operator = In
	return wf
}
func (wf *WhereFragment) NotIn() *WhereFragment {
	wf.Operator = NotIn
	return wf
}
func (wf *WhereFragment) Between() *WhereFragment {
	wf.Operator = Between
	return wf
}
func (wf *WhereFragment) NotBetween() *WhereFragment {
	wf.Operator = NotBetween
	return wf
}
func (wf *WhereFragment) Like() *WhereFragment {
	wf.Operator = Like
	return wf
}
func (wf *WhereFragment) NotLike() *WhereFragment {
	wf.Operator = NotLike
	return wf
}
func (wf *WhereFragment) Greatest() *WhereFragment {
	wf.Operator = Greatest
	return wf
}
func (wf *WhereFragment) Least() *WhereFragment {
	wf.Operator = Least
	return wf
}
func (wf *WhereFragment) Equal() *WhereFragment {
	wf.Operator = Equal
	return wf
}
func (wf *WhereFragment) NotEqual() *WhereFragment {
	wf.Operator = NotEqual
	return wf
}
func (wf *WhereFragment) Exists() *WhereFragment {
	wf.Operator = Exists
	return wf
}
func (wf *WhereFragment) NotExists() *WhereFragment {
	wf.Operator = NotExists
	return wf
}
func (wf *WhereFragment) Less() *WhereFragment {
	wf.Operator = Less
	return wf
}
func (wf *WhereFragment) Greater() *WhereFragment {
	wf.Operator = Greater
	return wf
}
func (wf *WhereFragment) LessOrEqual() *WhereFragment {
	wf.Operator = LessOrEqual
	return wf
}
func (wf *WhereFragment) GreaterOrEqual() *WhereFragment {
	wf.Operator = GreaterOrEqual
	return wf
}

func (wf *WhereFragment) Regexp() *WhereFragment {
	wf.Operator = Regexp
	return wf
}

func (wf *WhereFragment) NotRegexp() *WhereFragment {
	wf.Operator = NotRegexp
	return wf
}

func (wf *WhereFragment) Xor() *WhereFragment {
	wf.Operator = Xor
	return wf
}

func (wf *WhereFragment) SpaceShip() *WhereFragment {
	wf.Operator = SpaceShip
	return wf
}

func (wf *WhereFragment) Coalesce() *WhereFragment {
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
func (wf *WhereFragment) PlaceHolder() *WhereFragment {
	wf.Argument = placeHolderOp(cahensConstant)
	wf.IsPlaceHolder = true
	return wf
}

func (wf *WhereFragment) Int(i int) *WhereFragment {
	wf.Argument = Int(i)
	return wf
}

func (wf *WhereFragment) Ints(i ...int) *WhereFragment {
	wf.Argument = Ints(i)
	return wf
}
func (wf *WhereFragment) Int64(i int64) *WhereFragment {
	wf.Argument = Int64(i)
	return wf
}

func (wf *WhereFragment) Int64s(i ...int64) *WhereFragment {
	wf.Argument = Int64s(i)
	return wf
}

func (wf *WhereFragment) Uint64(i uint64) *WhereFragment {
	wf.Argument = Uint64(i)
	return wf
}

func (wf *WhereFragment) Float64(f float64) *WhereFragment {
	wf.Argument = Float64(f)
	return wf
}
func (wf *WhereFragment) Float64s(f ...float64) *WhereFragment {
	wf.Argument = Float64s(f)
	return wf
}
func (wf *WhereFragment) String(s string) *WhereFragment {
	wf.Argument = String(s)
	return wf
}

func (wf *WhereFragment) Strings(s ...string) *WhereFragment {
	wf.Argument = Strings(s)
	return wf
}

// Bool uses bool values for comparison.
func (wf *WhereFragment) Bool(b bool) *WhereFragment {
	wf.Argument = Bool(b)
	return wf
}

// Bools uses bool values for comparison.
func (wf *WhereFragment) Bools(b ...bool) *WhereFragment {
	wf.Argument = Bools(b)
	return wf
}

// Bytes uses a byte slice for comparison. Providing a nil value returns a
// NULL type. Detects between valid UTF-8 strings and binary data. Later gets
// hex encoded.
func (wf *WhereFragment) Bytes(p []byte) *WhereFragment {
	wf.Argument = Bytes(p)
	return wf
}

func (wf *WhereFragment) BytesSlice(p ...[]byte) *WhereFragment {
	wf.Argument = BytesSlice(p)
	return wf
}

// Time uses time.Time values for comparison.
func (wf *WhereFragment) Time(t time.Time) *WhereFragment {
	wf.Argument = MakeTime(t)
	return wf
}

// Times uses time.Time values for comparison.
func (wf *WhereFragment) Times(t ...time.Time) *WhereFragment {
	wf.Argument = Times(t)
	return wf
}

// NullString uses nullable string values for comparison.
func (wf *WhereFragment) NullString(values ...NullString) *WhereFragment {
	if len(values) == 1 {
		wf.Argument = values[0]
	} else {
		wf.Argument = ArgNullStrings(values)
	}
	return wf
}

// NullFloat64 uses nullable float64 values for comparison.
func (wf *WhereFragment) NullFloat64(values ...NullFloat64) *WhereFragment {
	if len(values) == 1 {
		wf.Argument = values[0]
	} else {
		wf.Argument = ArgNullFloat64s(values)
	}
	return wf
}

// NullInt64 uses nullable int64 values for comparison.
func (wf *WhereFragment) NullInt64(values ...NullInt64) *WhereFragment {
	if len(values) == 1 {
		wf.Argument = values[0]
	} else {
		wf.Argument = ArgNullInt64s(values)
	}
	return wf
}

// NullBool uses nullable bool values for comparison.
func (wf *WhereFragment) NullBool(value NullBool) *WhereFragment {
	wf.Argument = value
	return wf
}

// NullTime uses nullable time values for comparison.
func (wf *WhereFragment) NullTime(values ...NullTime) *WhereFragment {
	if len(values) == 1 {
		wf.Argument = values[0]
	} else {
		wf.Argument = ArgNullTimes(values)
	}
	return wf
}

// Value uses driver.Valuers for comparison.
func (wf *WhereFragment) Value(values ...driver.Valuer) *WhereFragment {
	wf.Argument = DriverValues(values)
	return wf
}

// Expression adds an unquoted SQL expression to a WHERE or HAVING statement.
func Expression(expression string, val ...Value) *WhereFragment {
	return &WhereFragment{
		IsExpression: true,
		Condition:    expression,
		Arguments:    val,
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
func (wfs WhereFragments) write(w queryWriter, conditionType byte) error {
	if len(wfs) == 0 {
		return nil
	}

	switch conditionType {
	case 'w':
		w.WriteString(" WHERE ")
	case 'h':
		w.WriteString(" HAVING ")
	}

	i := 0
	for _, f := range wfs {

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
			// Only write the operator in case there is no place holder and we
			// have one value.
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
			panic(errors.NewNotSupportedf("[dbr] Multiple values for a column are not supported\nWhereFragment: %#v\n", f))
		}

		w.WriteByte(')')
		i++
	}
	return nil
}

// conditionType enum of j=join, w=where, h=having
func (wfs WhereFragments) appendArgs(vals Values, conditionType byte) (_ Values, pendingValPos []int, err error) {
	if len(wfs) == 0 {
		return vals, pendingValPos, nil
	}

	pendingArgPosCount := len(vals)
	i := 0
	for _, f := range wfs {

		switch {
		case conditionType == 'j' && len(f.Using) > 0:
			return vals, pendingValPos, nil // done, only one USING allowed

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
			// Because addArg=false does not add below the values and we must
			// later swap the positions.
			pendingValPos = append(pendingValPos, pendingArgPosCount)

		case f.Argument != nil:
			// a column only supports one value.
			addArg = writeOperator(backHole{}, f.Argument.len(), f.Operator)
		case f.Sub.Select != nil:
			vals, err = f.Sub.Select.appendArgs(vals)
			if err != nil {
				return nil, pendingValPos, errors.Wrapf(err, "[dbr] write failed SubSelect for table: %q", f.Sub.Select.Table.String())
			}
		}

		if addArg {
			if f.Argument != nil {
				vals = append(vals, f.Argument)
			}
			vals = append(vals, f.Arguments...)
		}
		pendingArgPosCount++
		i++
	}
	return vals, pendingValPos, nil
}

func appendAssembledArgs(pendingValPos []int, rec ValuesAppender, vals Values, stmtType int, columns []string) (_ Values, err error) {
	if rec == nil {
		return vals, nil
	}

	lenBefore := len(vals)
	vals, err = rec.AppendValues(stmtType, vals, columns)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	lenAfter := len(vals)
	if lenAfter > lenBefore {
		j := 0
		newLen := lenAfter - len(pendingValPos)
		for i := newLen; i < lenAfter; i++ {
			vals[pendingValPos[j]], vals[i] = vals[i], vals[pendingValPos[j]]
			j++
		}
		vals = vals[:newLen] // remove the appended placeHolderOp types after swapping
	}
	return vals, nil
}
