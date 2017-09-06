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

package dml

import (
	"bytes"
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
	Coalesce       Op = 'c'          // Returns the first non-NULL value in the list, or NULL if there are no non-NULL arguments.
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

func writePlaceHolderList(w *bytes.Buffer, valLen int) {
	w.WriteByte('(')
	for j := 0; j < valLen; j++ {
		if j > 0 {
			w.WriteByte(',')
		}
		w.WriteByte(placeHolderRune)
	}
	w.WriteByte(')')
}

func (o Op) write(w *bytes.Buffer, argLen int) (err error) {
	// hasArgs value only used in cases where we have in the parent caller
	// function a sub-select. sub-selects do not need a place holder.
	hasArgs := argLen > 0
	switch o {
	case Null:
		_, err = w.WriteString(" IS NULL")
	case NotNull:
		_, err = w.WriteString(" IS NOT NULL")
	case In:
		_, err = w.WriteString(" IN ")
		if hasArgs {
			writePlaceHolderList(w, argLen)
		}
	case NotIn:
		_, err = w.WriteString(" NOT IN ")
		if hasArgs {
			writePlaceHolderList(w, argLen)
		}
	case Like:
		_, err = w.WriteString(" LIKE ")
		if hasArgs {
			err = w.WriteByte(placeHolderRune)
		}
	case NotLike:
		_, err = w.WriteString(" NOT LIKE ")
		if hasArgs {
			err = w.WriteByte(placeHolderRune)
		}
	case Regexp:
		_, err = w.WriteString(" REGEXP ")
		if hasArgs {
			err = w.WriteByte(placeHolderRune)
		}
	case NotRegexp:
		_, err = w.WriteString(" NOT REGEXP ")
		if hasArgs {
			err = w.WriteByte(placeHolderRune)
		}
	case Between:
		_, err = w.WriteString(" BETWEEN ? AND ?")
	case NotBetween:
		_, err = w.WriteString(" NOT BETWEEN ? AND ?")
	case Greatest:
		_, err = w.WriteString(" GREATEST ")
		writePlaceHolderList(w, argLen)
	case Least:
		_, err = w.WriteString(" LEAST ")
		writePlaceHolderList(w, argLen)
	case Coalesce:
		_, err = w.WriteString(" COALESCE ")
		writePlaceHolderList(w, argLen)
	case Xor:
		_, err = w.WriteString(" XOR ")
		if hasArgs {
			err = w.WriteByte(placeHolderRune)
		}
	case Exists:
		_, err = w.WriteString(" EXISTS ")
	case NotExists:
		_, err = w.WriteString(" NOT EXISTS ")
	case NotEqual:
		_, err = w.WriteString(" != ")
		if hasArgs {
			err = w.WriteByte(placeHolderRune)
		}
	case Less:
		_, err = w.WriteString(" < ")
		if hasArgs {
			err = w.WriteByte(placeHolderRune)
		}
	case Greater:
		_, err = w.WriteString(" > ")
		if hasArgs {
			err = w.WriteByte(placeHolderRune)
		}
	case LessOrEqual:
		_, err = w.WriteString(" <= ")
		if hasArgs {
			err = w.WriteByte(placeHolderRune)
		}
	case GreaterOrEqual:
		_, err = w.WriteString(" >= ")
		if hasArgs {
			err = w.WriteByte(placeHolderRune)
		}
	case SpaceShip:
		_, err = w.WriteString(" <=> ")
		if hasArgs {
			err = w.WriteByte(placeHolderRune)
		}
	default: // and case Equal
		_, err = w.WriteString(" = ")
		if hasArgs {
			err = w.WriteByte(placeHolderRune)
		}
	}
	return
}

// hasArgs returns true if the Operator requires arguments to compare with.
func (o Op) hasArgs(argLen int) (addArg bool) {
	// hasArg value only used in cases where we have in the parent caller
	// function a sub-select. sub-selects do not need a place holder.
	switch o {
	case Null, NotNull:
		addArg = false
	case Like, NotLike, Regexp, NotRegexp, Between, NotBetween, Greatest, Least, Coalesce, Xor:
		addArg = true
	default:
		addArg = argLen > 0
	}
	return
}

// Condition implements a single condition often used in WHERE, ON, SET and ON
// DUPLICATE KEY UPDATE. Please use the helper functions instead of using this
// type directly.
type Condition struct {
	Aliased string
	// Left can contain either a valid identifier or an expression. Set field
	// `IsLeftExpression` to true to avoid quoting of the this field. Left can also
	// contain a string in the format `qualifier.identifier`.
	Left string
	// Right defines the right hand side for an assignment which can be either a
	// single argument, multiple arguments in case of an expression, a sub
	// select or a name of a column.
	Right struct {
		// Column defines a column name to compare to. The column, with an
		// optional qualifier, gets quoted, in case IsExpression is false.
		Column    string
		NamedArg  string
		Argument  argument // Only set in case of no expression
		Arguments Arguments
		// Select adds a sub-select to the where statement. Column must be
		// either a column name or anything else which can handle the result of
		// a sub-select.
		Sub *Select
		// IsExpression if true field `Column` gets treated as an expression.
		// Additionally the field Right.Arguments will be read to extract any
		// given args.
		IsExpression bool
	}
	// Operator contains the comparison logic like LIKE, IN, GREATER, etc ...
	// defaults to EQUAL.
	Operator Op
	// IsLeftExpression if set to true, the field Left won't get quoted and
	// treated as an expression. Additionally the field Right.Arguments will be
	// read to extract any given args.
	IsLeftExpression bool
	// Logical states how multiple WHERE statements will be connected.
	// Default to AND. Possible values are a=AND, o=OR, x=XOR, n=NOT
	Logical byte
	// IsPlaceHolder true if the current WHERE condition just acts as a place
	// holder for a prepared statement or an interpolation.
	IsPlaceHolder bool
	// Columns is a list of column names which get quoted during SQL statement
	// creation in the JOIN part for the USING syntax. Additionally used in ON
	// DUPLICATE KEY.
	Columns []string
}

// Conditions provides a list where the left hand side gets an assignment from
// the right hand side. Mostly used in
type Conditions []*Condition

// Joins defines multiple join conditions.
type Joins []*join

type join struct {
	// JoinType can be LEFT, RIGHT, INNER, OUTER, CROSS or another word.
	JoinType string
	// Table name and alias of the table
	Table identifier
	// On join on those conditions
	On Conditions
}

// Alias assigns an alias name to the condition.
func (c *Condition) Alias(a string) *Condition {
	c.Aliased = a
	return c
}

// And sets the logical AND operator
func (c *Condition) And() *Condition {
	c.Logical = logicalAnd
	return c
}

// Or sets the logical OR operator
func (c *Condition) Or() *Condition {
	c.Logical = logicalOr
	return c
}

func (c *Condition) isExpression() bool {
	return c.IsLeftExpression || c.Right.IsExpression
}

// intersectConditions iterates over each WHERE fragment and appends all
// conditions aka column names to the slice c.
func (cs Conditions) intersectConditions(cols []string) []string {
	// this calculates the intersection of the columns in Conditions which
	// already have an value provided/assigned and those where the arguments
	// must be assembled from the interface ArgumentsAppender. If the arguments
	// should be assembled from the interface IsPlaceHolder is true.
	for _, cnd := range cs {
		if cnd.IsPlaceHolder {
			cols = append(cols, cnd.Left)
		}
	}
	return cols
}

// leftHands appends all Left strings to c. Use to get a list of all columns.
func (cs Conditions) leftHands(columns []string) []string {
	for _, cnd := range cs {
		columns = append(columns, cnd.Left)
	}
	return columns
}

// Columns add syntactic sugar to a JOIN or ON DUPLICATE KEY statement: In case
// of JOIN: The USING(column_list) clause names a list of columns that must
// exist in both tables. If tables a and b both contain columns c1, c2, and c3,
// the following join compares corresponding columns from the two tables:
//		a LEFT JOIN b USING (c1, c2, c3)
// The columns list gets quoted while writing the query string. In case of ON
// DUPLICATE KEY each column gets written like: `column`=VALUES(`column`).
// Any other field in *Condition gets ignored once field Columns has been set.
func Columns(columns ...string) *Condition {
	return &Condition{
		Columns: columns,
	}
}

// Column adds a new condition.
func Column(columnName string) *Condition {
	return &Condition{
		Left: columnName,
	}
}

// Expr adds an unquoted SQL expression to a column, WHERE, HAVING, SET or ON DUPLICATE
// KEY statement. Each item of an expression gets written into the buffer
// without a separator.
func Expr(expression string) *Condition {
	return &Condition{
		Left:             expression,
		IsLeftExpression: true,
	}
}

// ParenthesisOpen sets an open parenthesis "(". Mostly used for OR conditions
// in combination with AND conditions.
func ParenthesisOpen() *Condition {
	return &Condition{
		Left: "(",
	}
}

// ParenthesisClose sets a closing parenthesis ")". Mostly used for OR
// conditions in combination with AND conditions.
func ParenthesisClose() *Condition {
	return &Condition{
		Left: ")",
	}
}

///////////////////////////////////////////////////////////////////////////////
// COMPARIONS OPERATOR
///////////////////////////////////////////////////////////////////////////////

func (c *Condition) Null() *Condition {
	c.Operator = Null
	return c
}
func (c *Condition) NotNull() *Condition {
	c.Operator = NotNull
	return c
}
func (c *Condition) In() *Condition {
	c.Operator = In
	return c
}
func (c *Condition) NotIn() *Condition {
	c.Operator = NotIn
	return c
}
func (c *Condition) Between() *Condition {
	c.Operator = Between
	return c
}
func (c *Condition) NotBetween() *Condition {
	c.Operator = NotBetween
	return c
}
func (c *Condition) Like() *Condition {
	c.Operator = Like
	return c
}
func (c *Condition) NotLike() *Condition {
	c.Operator = NotLike
	return c
}
func (c *Condition) Greatest() *Condition {
	c.Operator = Greatest
	return c
}
func (c *Condition) Least() *Condition {
	c.Operator = Least
	return c
}
func (c *Condition) Equal() *Condition {
	c.Operator = Equal
	return c
}
func (c *Condition) NotEqual() *Condition {
	c.Operator = NotEqual
	return c
}
func (c *Condition) Exists() *Condition {
	c.Operator = Exists
	return c
}
func (c *Condition) NotExists() *Condition {
	c.Operator = NotExists
	return c
}
func (c *Condition) Less() *Condition {
	c.Operator = Less
	return c
}
func (c *Condition) Greater() *Condition {
	c.Operator = Greater
	return c
}
func (c *Condition) LessOrEqual() *Condition {
	c.Operator = LessOrEqual
	return c
}
func (c *Condition) GreaterOrEqual() *Condition {
	c.Operator = GreaterOrEqual
	return c
}

func (c *Condition) Regexp() *Condition {
	c.Operator = Regexp
	return c
}

func (c *Condition) NotRegexp() *Condition {
	c.Operator = NotRegexp
	return c
}

func (c *Condition) Xor() *Condition {
	c.Operator = Xor
	return c
}

func (c *Condition) SpaceShip() *Condition {
	c.Operator = SpaceShip
	return c
}

func (c *Condition) Coalesce() *Condition {
	c.Operator = Coalesce
	return c
}

///////////////////////////////////////////////////////////////////////////////
//		TYPES
///////////////////////////////////////////////////////////////////////////////

// Column compares the left hand side with this column name.
func (c *Condition) Column(col string) *Condition {
	c.Right.Column = col
	return c
}

// NamedArg if set the MySQL/MariaDB placeholder `?` will be replaced with a
// name. Records which implement ArgumentAppender must also use this name.
func (c *Condition) NamedArg(n string) *Condition {
	c.Right.NamedArg = n
	c.IsPlaceHolder = true
	return c
}

// PlaceHolder sets the database specific place holder character. Mostly used in
// prepared statements and for interpolation.
func (c *Condition) PlaceHolder() *Condition {
	c.Right.Argument.set(placeHolder(-7)) // value -7 does not matter
	c.IsPlaceHolder = true
	return c
}

// Sub compares the left hand side with the SELECT of the right hand side.
// Choose the appropriate comparison operator, default is IN.
func (c *Condition) Sub(sub *Select) *Condition {
	c.Right.Sub = sub
	if c.Operator == 0 {
		c.Operator = In
	}
	return c
}

// Expr compares the left hand side with the expression of the right hand
// side.
func (c *Condition) Expr(expression string) *Condition {
	c.Right.Column = expression
	c.Right.IsExpression = c.Right.Column != ""
	return c
}

func (c *Condition) Int(i int) *Condition {
	if c.isExpression() {
		c.Right.Arguments = c.Right.Arguments.Int64(int64(i))
		return c
	}
	c.Right.Argument.set(i)
	return c
}

func (c *Condition) Ints(i ...int) *Condition {
	if c.isExpression() {
		c.Right.Arguments = c.Right.Arguments.Ints(i...)
		return c
	}
	c.Right.Argument.set(i)
	return c
}

func (c *Condition) Int64(i int64) *Condition {
	if c.isExpression() {
		c.Right.Arguments = c.Right.Arguments.Int64(i)
		return c
	}
	c.Right.Argument.set(i)
	return c
}

func (c *Condition) Int64s(i ...int64) *Condition {
	if c.isExpression() {
		c.Right.Arguments = c.Right.Arguments.Int64s(i...)
		return c
	}
	c.Right.Argument.set(i)
	return c
}

func (c *Condition) Uint64(i uint64) *Condition {
	if c.isExpression() {
		c.Right.Arguments = c.Right.Arguments.Uint64(i)
		return c
	}
	c.Right.Argument.set(i)
	return c
}

func (c *Condition) Uint64s(i ...uint64) *Condition {
	if c.isExpression() {
		c.Right.Arguments = c.Right.Arguments.Uint64s(i...)
		return c
	}
	c.Right.Argument.set(i)
	return c
}

func (c *Condition) Float64(f float64) *Condition {
	if c.isExpression() {
		c.Right.Arguments = c.Right.Arguments.Float64(f)
		return c
	}
	c.Right.Argument.set(f)
	return c
}
func (c *Condition) Float64s(f ...float64) *Condition {
	if c.isExpression() {
		c.Right.Arguments = c.Right.Arguments.Float64s(f...)
		return c
	}
	c.Right.Argument.set(f)
	return c
}
func (c *Condition) Str(s string) *Condition {
	if c.isExpression() {
		c.Right.Arguments = c.Right.Arguments.String(s)
		return c
	}
	c.Right.Argument.set(s)
	return c
}

func (c *Condition) Strs(s ...string) *Condition {
	if c.isExpression() {
		c.Right.Arguments = c.Right.Arguments.Strings(s...)
		return c
	}
	c.Right.Argument.set(s)
	return c
}

func (c *Condition) Bool(b bool) *Condition {
	if c.isExpression() {
		c.Right.Arguments = c.Right.Arguments.Bool(b)
		return c
	}
	c.Right.Argument.set(b)
	return c
}

func (c *Condition) Bools(b ...bool) *Condition {
	if c.isExpression() {
		c.Right.Arguments = c.Right.Arguments.Bools(b...)
		return c
	}
	c.Right.Argument.set(b)
	return c
}

// Bytes uses a byte slice for comparison. Providing a nil value returns a
// NULL type. Detects between valid UTF-8 strings and binary data. Later gets
// hex encoded.
func (c *Condition) Bytes(p []byte) *Condition {
	if c.isExpression() {
		c.Right.Arguments = c.Right.Arguments.Bytes(p)
		return c
	}
	c.Right.Argument.set(p)
	return c
}

func (c *Condition) BytesSlice(p ...[]byte) *Condition {
	if c.isExpression() {
		c.Right.Arguments = c.Right.Arguments.BytesSlice(p...)
		return c
	}
	c.Right.Argument.set(p)
	return c
}

func (c *Condition) Time(t time.Time) *Condition {
	if c.isExpression() {
		c.Right.Arguments = c.Right.Arguments.Time(t)
		return c
	}
	c.Right.Argument.set(t)
	return c
}

func (c *Condition) Times(t ...time.Time) *Condition {
	if c.isExpression() {
		c.Right.Arguments = c.Right.Arguments.Times(t...)
		return c
	}
	c.Right.Argument.set(t)
	return c
}

func (c *Condition) NullString(nv ...NullString) *Condition {
	if c.isExpression() {
		c.Right.Arguments = c.Right.Arguments.NullString(nv...)
		return c
	}
	c.Right.Argument.set(nv)
	return c
}

func (c *Condition) NullFloat64(nv ...NullFloat64) *Condition {
	if c.isExpression() {
		c.Right.Arguments = c.Right.Arguments.NullFloat64(nv...)
		return c
	}
	c.Right.Argument.set(nv)
	return c
}

func (c *Condition) NullInt64(nv ...NullInt64) *Condition {
	if c.isExpression() {
		c.Right.Arguments = c.Right.Arguments.NullInt64(nv...)
		return c
	}
	c.Right.Argument.set(nv)
	return c
}

func (c *Condition) NullBool(nv ...NullBool) *Condition {
	if c.isExpression() {
		c.Right.Arguments = c.Right.Arguments.NullBool(nv...)
		return c
	}
	c.Right.Argument.set(nv)
	return c
}

func (c *Condition) NullTime(nv ...NullTime) *Condition {
	if c.isExpression() {
		c.Right.Arguments = c.Right.Arguments.NullTime(nv...)
		return c
	}
	c.Right.Argument.set(nv)
	return c
}

// Values onlny usable in case for ON DUPLCIATE KEY to generate a statement like:
//		column=VALUES(column)
func (c *Condition) Values() *Condition {
	// noop just to lower the cognitive overload when reading the code where
	// this function gets used.
	return c
}

// DriverValue uses driver.Valuers for comparison. Named DriverValue to avoid
// confusion with Values() function.
func (c *Condition) DriverValue(dv ...driver.Valuer) *Condition {
	c.Right.Arguments = c.Right.Arguments.DriverValue(dv...)
	return c
}
func (c *Condition) DriverValues(dv ...driver.Valuer) *Condition {
	c.Right.Arguments = c.Right.Arguments.DriverValues(dv...)
	return c
}

///////////////////////////////////////////////////////////////////////////////
//		FUNCTIONS / EXPRESSIONS
///////////////////////////////////////////////////////////////////////////////

// SQLCase see description at function SQLCase.
func (c *Condition) SQLCase(value, defaultValue string, compareResult ...string) *Condition {
	c.Right.Column = sqlCase(value, defaultValue, compareResult...)
	c.Right.IsExpression = c.Right.Column != ""
	return c
}

// SQLIfNull see description at function SQLIfNull.
func (c *Condition) SQLIfNull(expression ...string) *Condition {
	c.Right.Column = sqlIfNull(expression)
	c.Right.IsExpression = c.Right.Column != ""
	return c
}

///////////////////////////////////////////////////////////////////////////////
//		INTERNAL
///////////////////////////////////////////////////////////////////////////////

// write writes the conditions for usage as restrictions in WHERE, HAVING or
// JOIN clauses. conditionType enum of j=join, w=where, h=having
func (cs Conditions) write(w *bytes.Buffer, conditionType byte) error {
	if len(cs) == 0 {
		return nil
	}

	switch conditionType {
	case 'w':
		w.WriteString(" WHERE ")
	case 'h':
		w.WriteString(" HAVING ")
	}

	i := 0
	for _, cnd := range cs {

		if conditionType == 'j' {
			if len(cnd.Columns) > 0 {
				w.WriteString(" USING (")
				for j, c := range cnd.Columns {
					if j > 0 {
						w.WriteByte(',')
					}
					Quoter.quote(w, c)
				}
				w.WriteByte(')')
				return nil // done, only one USING allowed
			}
			if i == 0 {
				w.WriteString(" ON ")
			}
		}

		if cnd.Left == ")" {
			w.WriteString(cnd.Left)
			continue
		}

		if i > 0 {
			// How the WHERE conditions are connected
			switch cnd.Logical {
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

		if cnd.Left == "(" {
			i = 0
			w.WriteString(cnd.Left)
			continue
		}

		w.WriteByte('(')
		// Code is a bit duplicated but can be refactored later.
		switch {
		case cnd.IsLeftExpression:
			phCount, err := writeExpression(w, cnd.Left, cnd.Right.Arguments)
			if err != nil {
				return errors.WithStack(err)
			}

			// Only write the operator in case there is no place holder and we
			// have one value.
			if phCount == 0 && (len(cnd.Right.Arguments) == 1 || cnd.Right.Argument.isSet) && cnd.Operator > 0 {
				eArg := cnd.Right.Argument
				if !eArg.isSet {
					eArg = cnd.Right.Arguments[0]
				}
				cnd.Operator.write(w, eArg.len())
			}
			// TODO a case where left and right are expressions
			// if cnd.Right.Expression.isset() {
			// }
		case cnd.Right.IsExpression:
			Quoter.WriteIdentifier(w, cnd.Left)
			cnd.Operator.write(w, 0) // must be zero because place holder get handled via repeatPlaceHolders function
			writeExpression(w, cnd.Right.Column, cnd.Right.Arguments)

		case cnd.Right.Sub != nil:
			Quoter.WriteIdentifier(w, cnd.Left)
			cnd.Operator.write(w, 0)
			w.WriteByte('(')
			if err := cnd.Right.Sub.toSQL(w); err != nil {
				return errors.Wrapf(err, "[dml] write failed SubSelect for table: %q", cnd.Right.Sub.Table.String())
			}
			w.WriteByte(')')

			// One Argument and no expression
		case cnd.Right.Argument.isSet && cnd.Right.Arguments == nil:
			Quoter.WriteIdentifier(w, cnd.Left)
			al := cnd.Right.Argument.len()
			if cnd.IsPlaceHolder {
				al = 1
			}
			if al > 1 && cnd.Operator == 0 { // no operator but slice applied, so creating an IN query.
				cnd.Operator = In
			}
			cnd.Operator.write(w, al)

			// No Argument and expression arguments
		case !cnd.Right.Argument.isSet && cnd.Right.Arguments != nil:
			Quoter.WriteIdentifier(w, cnd.Left)
			al := cnd.Right.Arguments.Len()

			if cnd.IsPlaceHolder {
				al = 1
			}
			if al > 1 && cnd.Operator == 0 { // no operator but slice applied, so creating an IN query.
				cnd.Operator = In
			}
			cnd.Operator.write(w, al)

			// compares the left column with the right column
		case cnd.Right.Column != "":
			Quoter.WriteIdentifier(w, cnd.Left)
			cnd.Operator.write(w, 0)
			Quoter.WriteIdentifier(w, cnd.Right.Column)

			// No Argument at all, which kinda is the default case
		case !cnd.Right.Argument.isSet && cnd.Right.Arguments == nil:
			Quoter.WriteIdentifier(w, cnd.Left)
			cOp := cnd.Operator
			if cOp == 0 {
				cOp = Null
			}
			cOp.write(w, 1)

		default:
			panic(errors.NewNotSupportedf("[dml] Multiple arguments for a column are not supported\nWhereFragment: %#v\n", cnd))
		}

		w.WriteByte(')')
		i++
	}
	return nil
}

const (
	appendArgsJOIN   = 'j'
	appendArgsWHERE  = 'w'
	appendArgsHAVING = 'h'
	appendArgsSET    = 's'
	appendArgsDUPKEY = 'd'
)

// conditionType enum of: see constants appendArgs
func (cs Conditions) appendArgs(args Arguments, conditionType byte) (_ Arguments, pendingArgPos []int, err error) {
	if len(cs) == 0 {
		return args, pendingArgPos, nil
	}

	pendingArgPosCount := len(args)
	i := 0
	for _, cnd := range cs {

		switch {
		case conditionType == appendArgsJOIN && len(cnd.Columns) > 0:
			return args, pendingArgPos, nil // done, only one USING allowed

		case cnd.Left == ")":
			continue

		case cnd.Left == "(":
			i = 0
			continue
		}

		addArg := false
		switch {
		case cnd.isExpression():
			addArg = cnd.Operator.hasArgs(len(cnd.Right.Arguments))
		case cnd.IsPlaceHolder:
			addArg = cnd.Operator.hasArgs(1) // always a length of one, see the `repeatPlaceHolders()` function
			// By keeping addArg as it is and not setting
			// addArg=false, this []int avoids
			// https://en.wikipedia.org/wiki/Permutation Which would
			// result in a Go Code like
			// https://play.golang.org/p/rZvW0qW1N7 (C) Volker Dobler
			// Because addArg=false does not add below the arguments and we must
			// later swap the positions.
			pendingArgPos = append(pendingArgPos, pendingArgPosCount)
		case cnd.Right.Argument.isSet:
			addArg = cnd.Operator.hasArgs(cnd.Right.Argument.len())
		case cnd.Right.Arguments != nil:
			addArg = cnd.Operator.hasArgs(cnd.Right.Arguments.Len())
		case cnd.Right.Sub != nil:
			args, err = cnd.Right.Sub.appendArgs(args)
			if err != nil {
				return nil, pendingArgPos, errors.Wrapf(err, "[dml] write failed SubSelect for table: %q", cnd.Right.Sub.Table.String())
			}
		}

		if addArg {
			if cnd.Right.Argument.isSet {
				args = append(args, cnd.Right.Argument)
			}
			args = append(args, cnd.Right.Arguments...)
			pendingArgPosCount++
		}
		i++
	}
	return args, pendingArgPos, nil
}

func (cs Conditions) writeSetClauses(w *bytes.Buffer) error {
	for i, cnd := range cs {
		if i > 0 {
			w.WriteString(", ")
		}
		Quoter.quote(w, cnd.Left)
		w.WriteByte('=')

		switch {
		case cnd.Right.IsExpression: // maybe that case is superfluous
			writeExpression(w, cnd.Right.Column, nil)
		case cnd.Right.Sub != nil:
			w.WriteByte('(')
			if err := cnd.Right.Sub.toSQL(w); err != nil {
				return errors.WithStack(err)
			}
			w.WriteByte(')')
		default:
			w.WriteByte(placeHolderRune)
		}
	}
	return nil
}

func writeValues(w *bytes.Buffer, column string) {
	w.WriteString("VALUES(")
	Quoter.quote(w, column)
	w.WriteByte(')')
}

// writeOnDuplicateKey writes the columns to `w` and appends the arguments to
// `args` and returns `args`.
// https://dev.mysql.com/doc/refman/5.7/en/insert-on-duplicate.html
func (cs Conditions) writeOnDuplicateKey(w *bytes.Buffer) error {
	if len(cs) == 0 {
		return nil
	}

	w.WriteString(" ON DUPLICATE KEY UPDATE ")
	for i, cnd := range cs {
		addColon := false
		for j, col := range cnd.Columns {
			if j > 0 {
				w.WriteString(", ")
			}
			Quoter.quote(w, col)
			w.WriteByte('=')
			writeValues(w, col)
			addColon = true
		}
		if cnd.Left == "" {
			continue
		}
		if i > 0 || addColon {
			w.WriteString(", ")
		}
		Quoter.quote(w, cnd.Left)
		w.WriteByte('=')

		switch {
		case cnd.Right.IsExpression: // maybe that case is superfluous
			writeExpression(w, cnd.Right.Column, nil)
		//case cnd.Right.Sub != nil:
		//	w.WriteByte('(')
		//	if err := cnd.Right.Sub.toSQL(w); err != nil {
		//		return errors.WithStack(err)
		//	}
		//	w.WriteByte(')')
		case !cnd.Right.Argument.isSet:
			writeValues(w, cnd.Left)
		default:
			w.WriteByte(placeHolderRune)
		}
	}
	return nil
}

func appendArgs(pendingArgPos []int, records map[string]ArgumentsAppender, args Arguments, defaultQualifier string, columns []string) (_ Arguments, err error) {
	// arguments list above is a bit long, maybe later this function can be
	// integrated into Conditions.appendArgs.
	if records == nil {
		return args, nil
	}

	lenBefore := len(args)
	var argCol [1]string
	for _, identifier := range columns {
		qualifier, column := splitColumn(identifier)
		if qualifier == "" {
			qualifier = defaultQualifier
		}
		if rec, ok := records[qualifier]; ok {
			argCol[0] = column
			args, err = rec.AppendArgs(args, argCol[:])
			if err != nil {
				return nil, errors.WithStack(err)
			}
		}
	}
	lenAfter := len(args)
	if lenAfter > lenBefore {
		j := 0
		newLen := lenAfter - len(pendingArgPos)
		for i := newLen; i < lenAfter; i++ {
			args[pendingArgPos[j]], args[i] = args[i], args[pendingArgPos[j]]
			j++
		}
		args = args[:newLen] // remove the appended placeHolderOp types after swapping
	}
	return args, nil
}

// splitColumn splits a string via its last dot into the qualifier and the
// column name.
func splitColumn(identifier string) (qualifier, column string) {
	// dot at a beginning and end of a string is illegal.
	// Using LastIndexByte allows to retain the database qualifier, so:
	// database.table.column will become in the return "database.table", "column"
	if dotIndex := strings.LastIndexByte(identifier, '.'); dotIndex > 0 && dotIndex+1 < len(identifier) {
		return identifier[:dotIndex], identifier[dotIndex+1:]
	}
	return "", identifier
}
