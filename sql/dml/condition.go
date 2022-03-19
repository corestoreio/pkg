// Copyright 2015-present, Cyrill @ Schumacher.fm and the CoreStore contributors
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
	"fmt"
	"strings"
	"time"

	"github.com/corestoreio/errors"
	"github.com/corestoreio/pkg/storage/null"
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

func (o Op) write(w *bytes.Buffer, args ...any) (err error) {
	var arg any
	if len(args) == 1 {
		arg = args[0]
	}

	switch o {
	case NotIn, NotLike, NotRegexp, NotBetween, NotExists:
		w.WriteString(" NOT")
	}

	switch o {
	case Null:
		_, err = w.WriteString(" IS NULL")
	case NotNull:
		_, err = w.WriteString(" IS NOT NULL")
	case In, NotIn:
		w.WriteString(" IN ")
		err = writeInterfaces(w, args)
	case Like, NotLike:
		w.WriteString(" LIKE ")
		err = writeInterfaceValue(arg, w, 0)
	case Regexp, NotRegexp:
		w.WriteString(" REGEXP ")
		err = writeInterfaceValue(arg, w, 0)
	case Between, NotBetween:
		w.WriteString(" BETWEEN ")
		if arg == nil {
			w.WriteByte(placeHolderRune)
			w.WriteString(" AND ") // don't write the last place holder as it gets written somewhere else
		} else {
			if err = writeInterfaceValue(arg, w, 1); err != nil {
				return errors.WithStack(err)
			}
			w.WriteString(" AND ")
			if err = writeInterfaceValue(arg, w, 2); err != nil {
				return errors.WithStack(err)
			}
		}
	case Greatest:
		w.WriteString(" GREATEST ")
		err = writeInterfaces(w, args)
	case Least:
		w.WriteString(" LEAST ")
		err = writeInterfaces(w, args)
	case Coalesce:
		w.WriteString(" COALESCE ")
		err = writeInterfaces(w, args)
	case Xor:
		w.WriteString(" XOR ")
		err = writeInterfaceValue(arg, w, 0)
	case Exists, NotExists:
		w.WriteString(" EXISTS ")
		err = writeInterfaces(w, args)
	case Less:
		w.WriteString(" < ")
		err = writeInterfaceValue(arg, w, 0)
	case Greater:
		w.WriteString(" > ")
		err = writeInterfaceValue(arg, w, 0)
	case LessOrEqual:
		w.WriteString(" <= ")
		err = writeInterfaceValue(arg, w, 0)
	case GreaterOrEqual:
		w.WriteString(" >= ")
		err = writeInterfaceValue(arg, w, 0)
	case SpaceShip:
		w.WriteString(" <=> ")
		err = writeInterfaceValue(arg, w, 0)
	case NotEqual:
		w.WriteString(" != ")
		err = writeInterfaceValue(arg, w, 0)
	default: // and case Equal
		w.WriteString(" = ")
		err = writeInterfaceValue(arg, w, 0)
	}
	return
}

// Conditions provides a list where the left hand side gets an assignment from
// the right hand side. Mostly used in
type Conditions []*Condition

// Clone creates a clone of the current object.
func (cs Conditions) Clone() Conditions {
	if cs == nil {
		return nil
	}
	cs2 := make(Conditions, len(cs))
	for i, c := range cs {
		cs2[i] = c.Clone()
	}
	return cs2
}

// Reset resets the slice to length zero and retains the allocated memory.
func (cs *Conditions) Reset() Conditions {
	cs2 := *cs
	for i := range cs2 {
		cs2[i] = nil
	}
	*cs = cs2[:0]
	return *cs
}

// Joins defines multiple join conditions.
type Joins []*join

// Clone creates a new clone of the current object.
func (js Joins) Clone() Joins {
	if js == nil {
		return nil
	}
	c := make(Joins, len(js))
	for i, j := range js {
		c[i] = j.Clone()
	}
	return c
}

type join struct {
	// JoinType can be LEFT, RIGHT, INNER, OUTER, CROSS or another word.
	JoinType string
	// Table name and alias of the table
	Table id
	// On join on those conditions
	On Conditions
}

// Clone creates a new clone of the current object.
func (j *join) Clone() *join {
	if j == nil {
		return nil
	}
	c := *j
	c.Table = j.Table.Clone()
	c.On = j.On.Clone()
	return &c
}

// Condition implements a single condition often used in WHERE, ON, SET and ON
// DUPLICATE KEY UPDATE. Please use the helper functions instead of using this
// type directly.
type Condition struct {
	previousErr error
	Aliased     string
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
		Column string
		// PlaceHolder can be a :named or the MySQL/MariaDB place holder
		// character `?`. If set, the current condition just acts as a place
		// holder for a prepared statement or an interpolation. In case of a
		// :named place holder for a prepared statement, the :named string gets
		// replaced with the `?`. The allowed characters are unicode letters and
		// digits.
		PlaceHolder string
		// arg gets written into the SQL string as a persistent argument
		arg any // Only set in case of no expression
		// args same as arg but only used in case of an expression.
		args []any
		// Select adds a sub-select to the where statement. Column must be
		// either a column name or anything else which can handle the result of
		// a sub-select.
		Sub *Select
		// IsExpression if true field `Column` gets treated as an expression.
		// Additionally the field Right.args will be read to extract any
		// given args.
		IsExpression bool
	}
	// Operator contains the comparison logic like LIKE, IN, GREATER, etc ...
	// defaults to EQUAL.
	Operator Op
	// IsLeftExpression if set to true, the field Left won't get quoted and
	// treated as an expression. Additionally the field Right.args will be
	// read to extract any given args.
	IsLeftExpression bool
	// Logical states how multiple WHERE statements will be connected.
	// Default to AND. Possible values are a=AND, o=OR, x=XOR, n=NOT
	Logical byte
	// Columns is a list of column names which get quoted during SQL statement
	// creation in the JOIN part for the USING syntax. Additionally used in ON
	// DUPLICATE KEY.
	Columns []string
}

// Clone creates a new clone of the current object. It resets the internal error
// field.
func (c *Condition) Clone() *Condition {
	if c == nil {
		return nil
	}
	c2 := *c
	c2.previousErr = nil
	if c2.Right.args != nil {
		a2 := make([]any, len(c2.Right.args))
		copy(a2, c2.Right.args)
		c2.Right.args = a2
	}
	c2.Right.Sub = c.Right.Sub.Clone()
	c2.Columns = cloneStringSlice(c.Columns)
	return &c2
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
// COMPARISON OPERATOR
///////////////////////////////////////////////////////////////////////////////

// Op sets a custom operator
func (c *Condition) Op(o Op) *Condition {
	c.Operator = o
	return c
}

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

// NamedArg treats a condition as a place holder. If set the MySQL/MariaDB
// placeholder `?` will be used and the provided name gets replaced. Records
// which implement ColumnMapper must also use this name. A dot in the name (for
// e.g. setting a qualifier) is not allowed.
func (c *Condition) NamedArg(n string) *Condition {
	c.Right.PlaceHolder = n
	return c
}

// PlaceHolder treats a condition as a placeholder. Sets the database specific
// placeholder character "?". Mostly used in prepared statements and for
// interpolation.
func (c *Condition) PlaceHolder() *Condition {
	c.Right.PlaceHolder = placeHolderStr
	return c
}

// Tuples allows to build a query string for tuple comparison.
// 	SELECT * FROM catalog_product_index_decimal_idx WHERE
//	(entity_id,attribute_id,store_id,source_id) IN (
//		(4,4,4,4), (3,3,3,3), (dynamical values)
//	);
// See test ... TBC
func (c *Condition) Tuples() *Condition {
	c.Right.PlaceHolder = placeHolderTuples
	return c
}

// PlaceHolders treats a condition as a string with multiple placeholders. Sets
// the database specific placeholder character "?" as many times as specified in
// variable count. Mostly used in prepared statements and for interpolation and
// when using the IN clause.
func (c *Condition) PlaceHolders(count int) *Condition {
	switch count {
	case 1:
		c.Right.PlaceHolder = "(?)"
	case 2:
		c.Right.PlaceHolder = "(?,?)"
	case 3:
		c.Right.PlaceHolder = "(?,?,?)"
	}
	if c.Right.PlaceHolder != "" {
		return c
	}
	var buf strings.Builder
	buf.WriteByte('(')
	for i := 0; i < count; i++ {
		if i > 0 {
			buf.WriteByte(',')
		}
		buf.WriteByte(placeHolderRune)
	}
	buf.WriteByte(')')
	c.Right.PlaceHolder = buf.String()
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
		c.Right.args = append(c.Right.args, int64(i))
		return c
	}
	c.Right.arg = i
	return c
}

func (c *Condition) Ints(i ...int) *Condition {
	if c.isExpression() {
		c.Right.args = append(c.Right.args, i)
		return c
	}
	c.Right.arg = i
	return c
}

func (c *Condition) Int64(i int64) *Condition {
	if c.isExpression() {
		c.Right.args = append(c.Right.args, i)
		return c
	}
	c.Right.arg = i
	return c
}

func (c *Condition) Int64s(i ...int64) *Condition {
	if c.isExpression() {
		c.Right.args = append(c.Right.args, i)
		return c
	}
	c.Right.arg = i
	return c
}

func (c *Condition) Uint64(i uint64) *Condition {
	if c.isExpression() {
		c.Right.args = append(c.Right.args, i)
		return c
	}
	c.Right.arg = i
	return c
}

func (c *Condition) Uint64s(i ...uint64) *Condition {
	if c.isExpression() {
		c.Right.args = append(c.Right.args, i)
		return c
	}
	c.Right.arg = i
	return c
}

// Add when needed
// func (c *Condition) Decimals(d ...Decimal) *Condition {}

func (c *Condition) Decimal(d null.Decimal) *Condition {
	v := d.String()
	if c.isExpression() {
		if v == sqlStrNullUC {
			c.Right.args = append(c.Right.args, nil)
		} else {
			c.Right.args = append(c.Right.args, v)
		}
		return c
	}
	if v == sqlStrNullUC {
		c.Right.arg = nil
	} else {
		c.Right.arg = v
	}
	return c
}

func (c *Condition) Float64(f float64) *Condition {
	if c.isExpression() {
		c.Right.args = append(c.Right.args, f)
		return c
	}
	c.Right.arg = f
	return c
}

func (c *Condition) Float64s(f ...float64) *Condition {
	if c.isExpression() {
		c.Right.args = append(c.Right.args, f)
		return c
	}
	c.Right.arg = f
	return c
}

func (c *Condition) Str(s string) *Condition {
	if c.isExpression() {
		c.Right.args = append(c.Right.args, s)
		return c
	}
	c.Right.arg = s
	return c
}

func (c *Condition) Strs(s ...string) *Condition {
	if c.isExpression() {
		c.Right.args = append(c.Right.args, s)
		return c
	}
	c.Right.arg = s
	return c
}

func (c *Condition) Bool(b bool) *Condition {
	if c.isExpression() {
		c.Right.args = append(c.Right.args, b)
		return c
	}
	c.Right.arg = b
	return c
}

func (c *Condition) Bools(b ...bool) *Condition {
	if c.isExpression() {
		c.Right.args = append(c.Right.args, b)
		return c
	}
	c.Right.arg = b
	return c
}

// Bytes uses a byte slice for comparison. Providing a nil value returns a
// NULL type. Detects between valid UTF-8 strings and binary data. Later gets
// hex encoded.
func (c *Condition) Bytes(p []byte) *Condition {
	if c.isExpression() {
		c.Right.args = append(c.Right.args, p)
		return c
	}
	c.Right.arg = p
	return c
}

func (c *Condition) BytesSlice(p ...[]byte) *Condition {
	if c.isExpression() {
		c.Right.args = append(c.Right.args, p)
		return c
	}
	c.Right.arg = p
	return c
}

func (c *Condition) Time(t time.Time) *Condition {
	if c.isExpression() {
		c.Right.args = append(c.Right.args, t)
		return c
	}
	c.Right.arg = t
	return c
}

func (c *Condition) Times(t ...time.Time) *Condition {
	if c.isExpression() {
		c.Right.args = append(c.Right.args, t)
		return c
	}
	c.Right.arg = t
	return c
}

func (c *Condition) NullString(nv null.String) *Condition {
	if c.isExpression() {
		c.Right.args = append(c.Right.args, nv)
		return c
	}
	c.Right.arg = nv
	return c
}

func (c *Condition) NullStrings(nv ...null.String) *Condition {
	if c.isExpression() {
		c.Right.args = append(c.Right.args, nv)
		return c
	}
	c.Right.arg = nv
	return c
}

func (c *Condition) NullFloat64(nv null.Float64) *Condition {
	if c.isExpression() {
		c.Right.args = append(c.Right.args, nv)
		return c
	}
	c.Right.arg = nv
	return c
}

func (c *Condition) NullFloat64s(nv ...null.Float64) *Condition {
	if c.isExpression() {
		c.Right.args = append(c.Right.args, nv)
		return c
	}
	c.Right.arg = nv
	return c
}

func (c *Condition) NullInt64(nv null.Int64) *Condition {
	if c.isExpression() {
		c.Right.args = append(c.Right.args, nv)
		return c
	}
	c.Right.arg = nv
	return c
}

func (c *Condition) NullInt64s(nv ...null.Int64) *Condition {
	if c.isExpression() {
		c.Right.args = append(c.Right.args, nv)
		return c
	}
	c.Right.arg = nv
	return c
}

func (c *Condition) NullBool(nv null.Bool) *Condition {
	if c.isExpression() {
		c.Right.args = append(c.Right.args, nv)
		return c
	}
	c.Right.arg = nv
	return c
}

func (c *Condition) NullBools(nv ...null.Bool) *Condition {
	if c.isExpression() {
		c.Right.args = append(c.Right.args, nv)
		return c
	}
	c.Right.arg = nv
	return c
}

func (c *Condition) NullTime(nv null.Time) *Condition {
	if c.isExpression() {
		c.Right.args = append(c.Right.args, nv)
		return c
	}
	c.Right.arg = nv
	return c
}

func (c *Condition) NullTimes(nv ...null.Time) *Condition {
	if c.isExpression() {
		c.Right.args = append(c.Right.args, nv)
		return c
	}
	c.Right.arg = nv
	return c
}

// Values only usable in case for ON DUPLICATE KEY to generate a statement like:
//		column=VALUES(column)
func (c *Condition) Values() *Condition {
	// noop just to lower the cognitive overload when reading the code where
	// this function gets used.
	return c
}

// DriverValue adds multiple of the same underlying values to the argument
// slice. When using different values, the last applied value wins and gets
// added to the argument slice. For example driver.Values of type `int` will
// result in []int.
func (c *Condition) DriverValue(dv ...driver.Valuer) *Condition {
	if c.previousErr != nil {
		return c
	}
	c.Right.args, c.previousErr = driverValue(c.Right.args, dv...)
	return c
}

// DriverValues adds each driver.Value as its own argument to the argument
// slice. It panics if the underlying type is not one of the allowed of
// interface driver.Valuer.
func (c *Condition) DriverValues(dv ...driver.Valuer) *Condition {
	if c.previousErr != nil {
		return c
	}
	c.Right.args, c.previousErr = driverValues(c.Right.args, dv...)
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
func (cs Conditions) write(w *bytes.Buffer, conditionType byte, placeHolders []string, isWithDBR bool) (_placeHolders []string, err error) {
	if len(cs) == 0 {
		return placeHolders, nil
	}

	switch conditionType {
	case 'w':
		w.WriteString(" WHERE ")
	case 'h':
		w.WriteString(" HAVING ")
	}

	i := 0
	for _, cnd := range cs {
		if cnd.previousErr != nil {
			return nil, errors.WithStack(cnd.previousErr)
		}
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
				return placeHolders, nil // done, only one USING allowed
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
		// Code is a bit duplicated but can be refactored later. The order of
		// the `case`s has been carefully implemented.
		switch lenArgs := len(cnd.Right.args); {
		case cnd.IsLeftExpression:
			var phCount int
			phCount, err = writeExpression(w, cnd.Left, cnd.Right.args)
			if err != nil {
				return nil, errors.WithStack(err)
			}

			// Only write the operator in case there is no place holder and we
			// have one value.
			switch {
			case phCount == 0 && (lenArgs == 1 || cnd.Right.arg != nil) && cnd.Operator > 0:
				eArg := cnd.Right.arg
				if eArg == nil {
					eArg = cnd.Right.args[0]
				}
				cnd.Operator.write(w, eArg)

			case cnd.Right.Sub != nil:
				if err = cnd.Operator.write(w); err != nil {
					return nil, errors.WithStack(err)
				}
				w.WriteByte('(')
				placeHolders, err = cnd.Right.Sub.toSQL(w, placeHolders)
				if err != nil {
					return nil, errors.Wrapf(err, "[dml] write failed SubSelect for table: %q", cnd.Right.Sub.Table.String())
				}
				w.WriteByte(')')
			}

		case cnd.Right.IsExpression:
			Quoter.WriteIdentifier(w, cnd.Left)
			if err = cnd.Operator.write(w); err != nil {
				return nil, errors.WithStack(err)
			}
			if _, err = writeExpression(w, cnd.Right.Column, cnd.Right.args); err != nil {
				return nil, errors.WithStack(err)
			}
		case cnd.Right.Sub != nil:
			Quoter.WriteIdentifier(w, cnd.Left)
			if err = cnd.Operator.write(w); err != nil {
				return nil, errors.WithStack(err)
			}
			w.WriteByte('(')
			placeHolders, err = cnd.Right.Sub.toSQL(w, placeHolders)
			if err != nil {
				return nil, errors.Wrapf(err, "[dml] write failed SubSelect for table: %q", cnd.Right.Sub.Table.String())
			}
			w.WriteByte(')')

		case cnd.Right.arg != nil && lenArgs == 0: // One Argument and no expression
			Quoter.WriteIdentifier(w, cnd.Left)
			if al, _ := sliceLen(cnd.Right.arg); al > 1 && cnd.Operator == 0 { // no operator but slice applied, so creating an IN query.
				cnd.Operator = In
			}
			if err = cnd.Operator.write(w, cnd.Right.arg); err != nil {
				return nil, errors.WithStack(err)
			}

		case cnd.Right.arg == nil && lenArgs > 0:
			Quoter.WriteIdentifier(w, cnd.Left)
			if totalSliceLenSimple(cnd.Right.args) > 1 && cnd.Operator == 0 { // no operator but slice applied, so creating an IN query.
				cnd.Operator = In
			}
			if err = cnd.Operator.write(w, cnd.Right.args...); err != nil {
				return nil, errors.WithStack(err)
			}

		case cnd.Right.Column != "": // compares the left column with the right column
			Quoter.WriteIdentifier(w, cnd.Left)
			if err = cnd.Operator.write(w); err != nil {
				return nil, errors.WithStack(err)
			}
			Quoter.WriteIdentifier(w, cnd.Right.Column)

		case cnd.Right.PlaceHolder == placeHolderTuples:
			w.WriteByte('(')
			for j, col := range cnd.Columns {
				if j > 0 {
					w.WriteString(", ")
				}
				Quoter.quote(w, col)
			}
			w.WriteByte(')')
			if err = cnd.Operator.write(w); err != nil {
				return nil, errors.WithStack(err)
			}
			if isWithDBR {
				fmt.Fprintf(w, placeHolderTuples, len(cnd.Columns))
				// BuilderBase.buildToSQL needs this hack to see if we have a tuple. If
				// so, sets containsTuples to true.
				placeHolders = append(placeHolders, placeHolderTuples)
			} else {
				w.WriteByte('(')
				writeTuplePlaceholders(w, 1, uint(len(cnd.Columns)))
				w.WriteByte(')')
			}

		case cnd.Right.PlaceHolder != "":
			Quoter.WriteIdentifier(w, cnd.Left)
			if err = cnd.Operator.write(w); err != nil {
				return nil, errors.WithStack(err)
			}

			switch {
			case cnd.Right.PlaceHolder == placeHolderStr:
				placeHolders = append(placeHolders, cnd.Left)
				w.WriteByte(placeHolderRune)
			case isNamedArg(cnd.Right.PlaceHolder):
				w.WriteByte(placeHolderRune)
				ph := cnd.Right.PlaceHolder
				if !strings.HasPrefix(cnd.Right.PlaceHolder, namedArgStartStr) {
					ph = namedArgStartStr + ph
				}
				placeHolders = append(placeHolders, ph)
			default:
				placeHolders = append(placeHolders, cnd.Left)
				w.WriteString(cnd.Right.PlaceHolder)
			}

		case cnd.Right.arg == nil && lenArgs == 0: // No Argument at all, which kinda is the default case
			Quoter.WriteIdentifier(w, cnd.Left)
			cOp := cnd.Operator
			if cOp == 0 {
				cOp = Null
			}
			if err = cOp.write(w); err != nil {
				return nil, errors.WithStack(err)
			}

		default:
			panic(errors.NotSupported.Newf("[dml] Multiple arguments for a column are not supported\nWhereFragment: %#v\n", cnd))
		}

		w.WriteByte(')')
		i++
	}
	return placeHolders, errors.WithStack(err)
}

func (cs Conditions) writeSetClauses(w *bytes.Buffer, placeHolders []string) ([]string, error) {
	for i, cnd := range cs {
		if i > 0 {
			w.WriteString(", ")
		}
		Quoter.quote(w, cnd.Left)
		w.WriteByte('=')

		switch {
		case cnd.Right.arg != nil && len(cnd.Right.args) == 0: // One Argument and no expression
			if err := writeInterfaceValue(cnd.Right.arg, w, 0); err != nil {
				return nil, errors.WithStack(err)
			}
		case cnd.Right.IsExpression: // maybe that case is superfluous
			if _, err := writeExpression(w, cnd.Right.Column, cnd.Right.args); err != nil {
				return nil, errors.WithStack(err)
			}
			placeHolders = append(placeHolders, cnd.Left)
		case cnd.Right.Sub != nil:
			w.WriteByte('(')
			var err error
			if placeHolders, err = cnd.Right.Sub.toSQL(w, placeHolders); err != nil {
				return nil, errors.WithStack(err)
			}
			w.WriteByte(')')
		default:
			placeHolders = append(placeHolders, cnd.Left)
			w.WriteByte(placeHolderRune)
		}
	}
	return placeHolders, nil
}

func writeSQLValues(w *bytes.Buffer, column string) {
	w.WriteString("VALUES(")
	Quoter.quote(w, column)
	w.WriteByte(')')
}

var onDuplicateKeyPart = []byte(` ON DUPLICATE KEY UPDATE `)

const onDuplicateKeyPartS = ` ON DUPLICATE KEY UPDATE `

// writeOnDuplicateKey writes the columns to `w` and appends the arguments to
// `args` and returns `args`.
// https://dev.mysql.com/doc/refman/5.7/en/insert-on-duplicate.html
func (cs Conditions) writeOnDuplicateKey(w *bytes.Buffer, placeHolders []string) ([]string, error) {
	if len(cs) == 0 {
		return placeHolders, nil
	}

	w.Write(onDuplicateKeyPart)
	for i, cnd := range cs {
		addColon := false
		for j, col := range cnd.Columns {
			if j > 0 {
				w.WriteString(", ")
			}
			Quoter.quote(w, col)
			w.WriteByte('=')
			writeSQLValues(w, col)
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
			if _, err := writeExpression(w, cnd.Right.Column, cnd.Right.args); err != nil {
				return nil, errors.WithStack(err)
			}

		case cnd.Right.PlaceHolder != "":

			switch {
			case cnd.Right.PlaceHolder == placeHolderStr:
				placeHolders = append(placeHolders, cnd.Left)
				w.WriteByte(placeHolderRune)
			case isNamedArg(cnd.Right.PlaceHolder):
				w.WriteByte(placeHolderRune)
				ph := cnd.Right.PlaceHolder
				if !strings.HasPrefix(cnd.Right.PlaceHolder, namedArgStartStr) {
					ph = namedArgStartStr + ph
				}
				placeHolders = append(placeHolders, ph)
			default:
				placeHolders = append(placeHolders, cnd.Left)
				w.WriteString(cnd.Right.PlaceHolder)
			}

		case cnd.Right.arg == nil:
			writeSQLValues(w, cnd.Left)
		case cnd.Right.arg != nil:
			if err := writeInterfaceValue(cnd.Right.arg, w, 0); err != nil {
				return nil, errors.WithStack(err)
			}

		default:
			placeHolders = append(placeHolders, cnd.Left)
			w.WriteByte(placeHolderRune)
		}
	}
	return placeHolders, nil
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
