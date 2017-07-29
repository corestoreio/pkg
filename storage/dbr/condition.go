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

	"github.com/corestoreio/csfw/util/bufferpool"
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

func writePlaceHolderList(w queryWriter, valLen int) {
	w.WriteByte('(')
	for j := 0; j < valLen; j++ {
		if j > 0 {
			w.WriteByte(',')
		}
		w.WriteByte('?')
	}
	w.WriteByte(')')
}

func (o Op) write(w queryWriter, argLen int) (err error) {
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
			err = w.WriteByte('?')
		}
	case NotLike:
		_, err = w.WriteString(" NOT LIKE ")
		if hasArgs {
			err = w.WriteByte('?')
		}
	case Regexp:
		_, err = w.WriteString(" REGEXP ")
		if hasArgs {
			err = w.WriteByte('?')
		}
	case NotRegexp:
		_, err = w.WriteString(" NOT REGEXP ")
		if hasArgs {
			err = w.WriteByte('?')
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
			err = w.WriteByte('?')
		}
	case Exists:
		_, err = w.WriteString(" EXISTS ")
	case NotExists:
		_, err = w.WriteString(" NOT EXISTS ")
	case NotEqual:
		_, err = w.WriteString(" != ")
		if hasArgs {
			err = w.WriteByte('?')
		}
	case Less:
		_, err = w.WriteString(" < ")
		if hasArgs {
			err = w.WriteByte('?')
		}
	case Greater:
		_, err = w.WriteString(" > ")
		if hasArgs {
			err = w.WriteByte('?')
		}
	case LessOrEqual:
		_, err = w.WriteString(" <= ")
		if hasArgs {
			err = w.WriteByte('?')
		}
	case GreaterOrEqual:
		_, err = w.WriteString(" >= ")
		if hasArgs {
			err = w.WriteByte('?')
		}
	case SpaceShip:
		_, err = w.WriteString(" <=> ")
		if hasArgs {
			err = w.WriteByte('?')
		}
	default: // and case Equal
		_, err = w.WriteString(" = ")
		if hasArgs {
			err = w.WriteByte('?')
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

// expression is just some lines to avoid a fully written string created by
// other functions. Each line can contain arbitrary characters. The lines get
// written without any separator into a buffer. Hooks/Event/Observer allow an
// easily modification of different line items. Much better than having a long
// string with a wrapped complex SQL expression. There is no need to export it.
type expressions []string

// write writes the strings into `w` and correctly handles the place holder
// repetition depending on the number of arguments.
func (e expressions) write(w queryWriter, arg ...Argument) (phCount int, err error) {
	eBuf := bufferpool.Get()
	defer bufferpool.Put(eBuf)

	args := Arguments(arg)

	for _, expr := range e {
		phCount += strings.Count(expr, placeHolderStr)
		if _, err = eBuf.WriteString(expr); err != nil {
			return phCount, errors.Wrapf(err, "[dbr] expression.write: failed to write %q", expr)
		}
	}
	if args != nil && phCount != args.len() {
		if err = repeatPlaceHolders(w, eBuf.Bytes(), args...); err != nil {
			return phCount, errors.WithStack(err)
		}
	} else {
		_, err = eBuf.WriteTo(w)
	}
	return phCount, errors.WithStack(err)
}

func (e expressions) isset() bool {
	return len(e) > 0
}

// Aliased appends a quoted alias name to the expression
func (e expressions) Alias(a string) expressions {
	e = append(e, " AS ")
	e = Quoter.appendQuote(e, a)
	return e
}

func (e expressions) String() string {
	buf := bufferpool.Get()
	defer bufferpool.Put(buf)
	e.write(buf)
	return buf.String()
}

// Condition implements a single condition often used in WHERE, ON, SET and ON
// DUPLICATE KEY UPDATE. Please use the helper functions instead of using this
// type directly.
type Condition struct {
	// Left can contain either a valid identifier or an expression. Set field
	// `IsLeftExpression` to true to avoid quoting of the this field. Left can also
	// contain a string in the format `qualifier.identifier`.
	Left string
	// LeftExpression defines multiple strings as an expression. Each string
	// gets written without any separator. If `LeftExpression` has been set, the
	// field `Left` gets ignored.
	LeftExpression expressions
	// Right defines the right hand side for an assignment
	Right struct {
		// Expression can contain multiple entries. Each slice item gets written
		// into the buffer when the SQL string gets build. Usage in SET and ON
		// DUPLICATE KEY.
		Expression expressions
		Argument   Argument  // Either this or the slice is set.
		Arguments  Arguments // Only set in case of an expression.
		// Select adds a sub-select to the where statement. Column must be
		// either a column name or anything else which can handle the result of
		// a sub-select.
		Sub *Select
	}
	// Operator contains the comparison logic like LIKE, IN, GREATER, etc ...
	// defaults to EQUAL.
	Operator Op

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
	return c.LeftExpression.isset() || c.Right.Expression.isset()
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

// Expression adds an unquoted SQL expression to a WHERE, HAVING, SET or ON
// DUPLICATE KEY statement. Each item of an expression gets written into the
// buffer without a separator.
func Expression(expression ...string) *Condition {
	return &Condition{
		LeftExpression: expression,
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

// Cahen's Constant, used as a random identifier that an argument is a place
// holder.
const cahensConstant = -64341

// PlaceHolder sets the database specific place holder character. Mostly used in
// prepared statements and for interpolation.
func (c *Condition) PlaceHolder() *Condition {
	c.Right.Argument = placeHolderOp(cahensConstant)
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

// Expression compares the left hand side with the expression of the right hand
// side.
func (c *Condition) Expression(exp ...string) *Condition {
	c.Right.Expression = exp
	return c
}

func (c *Condition) Int(i int) *Condition {
	if c.isExpression() {
		c.Right.Arguments = append(c.Right.Arguments, Int(i))
		return c
	}
	c.Right.Argument = Int(i)
	return c
}

func (c *Condition) Ints(i ...int) *Condition {
	if c.isExpression() {
		c.Right.Arguments = append(c.Right.Arguments, Ints(i))
		return c
	}
	c.Right.Argument = Ints(i)
	return c
}
func (c *Condition) Int64(i int64) *Condition {
	if c.isExpression() {
		c.Right.Arguments = append(c.Right.Arguments, Int64(i))
		return c
	}
	c.Right.Argument = Int64(i)
	return c
}

func (c *Condition) Int64s(i ...int64) *Condition {
	if c.isExpression() {
		c.Right.Arguments = append(c.Right.Arguments, Int64s(i))
		return c
	}
	c.Right.Argument = Int64s(i)
	return c
}

func (c *Condition) Uint64(i uint64) *Condition {
	if c.isExpression() {
		c.Right.Arguments = append(c.Right.Arguments, Uint64(i))
		return c
	}
	c.Right.Argument = Uint64(i)
	return c
}

func (c *Condition) Float64(f float64) *Condition {
	if c.isExpression() {
		c.Right.Arguments = append(c.Right.Arguments, Float64(f))
		return c
	}
	c.Right.Argument = Float64(f)
	return c
}
func (c *Condition) Float64s(f ...float64) *Condition {
	if c.isExpression() {
		c.Right.Arguments = append(c.Right.Arguments, Float64s(f))
		return c
	}
	c.Right.Argument = Float64s(f)
	return c
}
func (c *Condition) String(s string) *Condition { // TODO rename to Str and Strs and use String() as fmt.Stringer
	if c.isExpression() {
		c.Right.Arguments = append(c.Right.Arguments, String(s))
		return c
	}
	c.Right.Argument = String(s)
	return c
}

func (c *Condition) Strings(s ...string) *Condition {
	if c.isExpression() {
		c.Right.Arguments = append(c.Right.Arguments, Strings(s))
		return c
	}
	c.Right.Argument = Strings(s)
	return c
}

func (c *Condition) Bool(b bool) *Condition {
	if c.isExpression() {
		c.Right.Arguments = append(c.Right.Arguments, Bool(b))
		return c
	}
	c.Right.Argument = Bool(b)
	return c
}

func (c *Condition) Bools(b ...bool) *Condition {
	if c.isExpression() {
		c.Right.Arguments = append(c.Right.Arguments, Bools(b))
		return c
	}
	c.Right.Argument = Bools(b)
	return c
}

// Bytes uses a byte slice for comparison. Providing a nil value returns a
// NULL type. Detects between valid UTF-8 strings and binary data. Later gets
// hex encoded.
func (c *Condition) Bytes(p []byte) *Condition {
	if c.isExpression() {
		c.Right.Arguments = append(c.Right.Arguments, Bytes(p))
		return c
	}
	c.Right.Argument = Bytes(p)
	return c
}

func (c *Condition) BytesSlice(p ...[]byte) *Condition {
	if c.isExpression() {
		c.Right.Arguments = append(c.Right.Arguments, BytesSlice(p))
		return c
	}
	c.Right.Argument = BytesSlice(p)
	return c
}

// Time uses time.Time arguments for comparison.
func (c *Condition) Time(t time.Time) *Condition {
	if c.isExpression() {
		c.Right.Arguments = append(c.Right.Arguments, MakeTime(t))
		return c
	}
	c.Right.Argument = MakeTime(t)
	return c
}

// Times uses time.Time arguments for comparison.
func (c *Condition) Times(t ...time.Time) *Condition {
	if c.isExpression() {
		c.Right.Arguments = append(c.Right.Arguments, Times(t))
		return c
	}
	c.Right.Argument = Times(t)
	return c
}

// NullString uses nullable string arguments for comparison.
func (c *Condition) NullString(nv ...NullString) *Condition {
	switch {
	case c.isExpression():
		c.Right.Arguments = append(c.Right.Arguments, NullStrings(nv))
	case len(nv) == 1:
		c.Right.Argument = nv[0]
	default:
		c.Right.Argument = NullStrings(nv)
	}
	return c
}

// NullFloat64 uses nullable float64 arguments for comparison.
func (c *Condition) NullFloat64(nv ...NullFloat64) *Condition {
	switch {
	case c.isExpression():
		c.Right.Arguments = append(c.Right.Arguments, NullFloat64s(nv))
	case len(nv) == 1:
		c.Right.Argument = nv[0]
	default:
		c.Right.Argument = NullFloat64s(nv)
	}
	return c
}

// NullInt64 uses nullable int64 arguments for comparison.
func (c *Condition) NullInt64(nv ...NullInt64) *Condition {
	switch {
	case c.isExpression():
		c.Right.Arguments = append(c.Right.Arguments, NullInt64s(nv))
	case len(nv) == 1:
		c.Right.Argument = nv[0]
	default:
		c.Right.Argument = NullInt64s(nv)
	}
	return c
}

// NullBool uses nullable bool arguments for comparison.
func (c *Condition) NullBool(nv NullBool) *Condition {
	switch {
	case c.isExpression():
		c.Right.Arguments = append(c.Right.Arguments, nv)
	default:
		c.Right.Argument = nv
	}
	return c
}

// NullTime uses nullable time arguments for comparison.
func (c *Condition) NullTime(nv ...NullTime) *Condition {
	switch {
	case c.isExpression():
		c.Right.Arguments = append(c.Right.Arguments, NullTimes(nv))
	case len(nv) == 1:
		c.Right.Argument = nv[0]
	default:
		c.Right.Argument = NullTimes(nv)
	}
	return c
}

// Values onlny usable in case for ON DUPLCIATE KEY to generate a statement like:
//		column=VALUES(column)
func (c *Condition) Values() *Condition {
	// noop just to lower the cognitive overload.
	return c
}

// DriverValue uses driver.Valuers for comparison.
func (c *Condition) DriverValue(dv ...driver.Valuer) *Condition {
	switch {
	case c.isExpression():
		c.Right.Arguments = append(c.Right.Arguments, DriverValues(dv))
	default:
		c.Right.Argument = DriverValues(dv)
	}
	return c
}

// write writes the conditions for usage as restrictions in WHERE, HAVING or
// JOIN clauses. conditionType enum of j=join, w=where, h=having
func (cs Conditions) write(w queryWriter, conditionType byte) error {
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
		case cnd.LeftExpression.isset():
			phCount, err := cnd.LeftExpression.write(w, cnd.Right.Arguments...)
			if err != nil {
				return errors.WithStack(err)
			}

			// Only write the operator in case there is no place holder and we
			// have one value.
			if phCount == 0 && (len(cnd.Right.Arguments) == 1 || cnd.Right.Argument != nil) && cnd.Operator > 0 {
				eArg := cnd.Right.Argument
				if eArg == nil {
					eArg = cnd.Right.Arguments[0]
				}
				cnd.Operator.write(w, eArg.len())
			}
			// TODO a case where left and right are expressions
			// if cnd.Right.Expression.isset() {
			// }
		case cnd.Right.Expression.isset():
			Quoter.WriteIdentifier(w, cnd.Left)
			cnd.Operator.write(w, 0) // must be zero because place holder get handled via repeatPlaceHolders function
			cnd.Right.Expression.write(w, cnd.Right.Arguments...)

		case cnd.Right.Sub != nil:
			Quoter.WriteIdentifier(w, cnd.Left)
			cnd.Operator.write(w, 0)
			w.WriteByte('(')
			if err := cnd.Right.Sub.toSQL(w); err != nil {
				return errors.Wrapf(err, "[dbr] write failed SubSelect for table: %q", cnd.Right.Sub.Table.String())
			}
			w.WriteByte(')')

		case cnd.Right.Argument != nil && cnd.Right.Arguments == nil:
			Quoter.WriteIdentifier(w, cnd.Left)
			al := cnd.Right.Argument.len()
			if cnd.IsPlaceHolder {
				al = 1
			}
			if al > 1 && cnd.Operator == 0 { // no operator but slice applied, so creating an IN query.
				cnd.Operator = In
			}
			cnd.Operator.write(w, al)

		case cnd.Right.Argument == nil && cnd.Right.Arguments == nil:
			Quoter.WriteIdentifier(w, cnd.Left)
			cOp := cnd.Operator
			if cOp == 0 {
				cOp = Null
			}
			cOp.write(w, 1)

		default:
			panic(errors.NewNotSupportedf("[dbr] Multiple arguments for a column are not supported\nWhereFragment: %#v\n", cnd))
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
			addArg = true
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

		case cnd.Right.Argument != nil:
			// a column only supports one value.
			addArg = cnd.Operator.hasArgs(cnd.Right.Argument.len())
		case cnd.Right.Sub != nil:
			args, err = cnd.Right.Sub.appendArgs(args)
			if err != nil {
				return nil, pendingArgPos, errors.Wrapf(err, "[dbr] write failed SubSelect for table: %q", cnd.Right.Sub.Table.String())
			}
		}

		if addArg {
			if cnd.Right.Argument != nil {
				args = append(args, cnd.Right.Argument)
			}
			args = append(args, cnd.Right.Arguments...)
		}
		pendingArgPosCount++
		i++
	}
	return args, pendingArgPos, nil
}

func (cs Conditions) writeSetClauses(w queryWriter) error {
	for i, cnd := range cs {
		if i > 0 {
			w.WriteString(", ")
		}
		Quoter.quote(w, cnd.Left)
		w.WriteByte('=')

		switch {
		case cnd.Right.Expression.isset(): // maybe that case is superfluous
			cnd.Right.Expression.write(w)
		case cnd.Right.Sub != nil:
			w.WriteByte('(')
			if err := cnd.Right.Sub.toSQL(w); err != nil {
				return errors.WithStack(err)
			}
			w.WriteByte(')')
		default:
			w.WriteByte('?')
		}
	}
	return nil
}

func writeValues(w queryWriter, column string) {
	w.WriteString("VALUES(")
	Quoter.quote(w, column)
	w.WriteByte(')')
}

// writeOnDuplicateKey writes the columns to `w` and appends the arguments to
// `args` and returns `args`.
// https://dev.mysql.com/doc/refman/5.7/en/insert-on-duplicate.html
func (cs Conditions) writeOnDuplicateKey(w queryWriter) error {
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
		case cnd.Right.Expression.isset(): // maybe that case is superfluous
			cnd.Right.Expression.write(w)
		//case cnd.Right.Sub != nil:
		//	w.WriteByte('(')
		//	if err := cnd.Right.Sub.toSQL(w); err != nil {
		//		return errors.WithStack(err)
		//	}
		//	w.WriteByte(')')
		case cnd.Right.Argument == nil:
			writeValues(w, cnd.Left)
		default:
			w.WriteByte('?')
		}
	}
	return nil
}

func appendAssembledArgs(pendingArgPos []int, rec ArgumentsAppender, args Arguments, stmtType int, columns []string) (_ Arguments, err error) {
	if rec == nil {
		return args, nil
	}

	lenBefore := len(args)
	args, err = rec.AppendArguments(stmtType, args, columns)
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
		args = args[:newLen] // remove the appended placeHolderOp types after swapping
	}
	return args, nil
}
