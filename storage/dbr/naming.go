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
	"unicode/utf8"

	"github.com/corestoreio/csfw/util/bufferpool"
	"github.com/corestoreio/errors"
)

const quote string = "`"
const quoteRune = '`'

// Quoter at the quoter to use for quoting text; use Mysql quoting by default.
var Quoter = MysqlQuoter{
	replacer: strings.NewReplacer(quote, ""),
}

type alias struct {
	// Derived Tables (Subqueries in the FROM Clause). A derived table is a
	// subquery in a SELECT statement FROM clause. Derived tables can return a
	// scalar, column, row, or table. Ignored in any other case.
	DerivedTable *Select
	// Name can be any kind of SQL expression or a valid identifier. It gets
	// quoted when `IsExpression` is false.
	Name string
	// Alias must be a valid identifier allowed for alias usage.
	Alias string
	// IsExpression if true the field `Name` will be treated as an expression and
	// won't get quoted when generating the SQL.
	IsExpression bool
	// Sort applies only to GROUP BY and ORDER BY clauses. 'd'=descending,
	// 0=default or nothing; 'a'=ascending.
	Sort byte
}

const (
	sortDescending byte = 'd'
	sortAscending  byte = 'a'
)

// MakeAlias creates a new name with an optional alias. Supports two arguments.
// 1. a qualifier name and 2. an alias.
func MakeAlias(nameAlias ...string) alias {
	a := alias{
		Name: nameAlias[0],
	}
	if len(nameAlias) > 1 {
		a.Alias = nameAlias[1]
	}
	return a
}

// MakeAliasExpr creates a new expression with an optional alias. Supports two
// arguments. 1. an expression and 2. an alias.
func MakeAliasExpr(expressionAlias ...string) alias {
	a := alias{
		IsExpression: true,
		Name:         expressionAlias[0],
	}
	if len(expressionAlias) > 1 {
		a.Alias = expressionAlias[1]
	}
	return a
}

func (a alias) isEmpty() bool { return a.Name == "" && a.DerivedTable == nil }

// String returns the correct stringyfied statement.
func (a alias) String() string {
	if a.IsExpression {
		return Quoter.exprAlias(a.Name, a.Alias)
	}
	return a.QuoteAs()
}

// QuoteAs always quuotes the name and the alias
func (a alias) QuoteAs() string {
	return Quoter.QuoteAs(a.Name, a.Alias)
}

// appendArgs assembles the arguments and appends them to `args`
func (a alias) appendArgs(args Arguments) (_ Arguments, err error) {
	if a.DerivedTable != nil {
		args, err = a.DerivedTable.appendArgs(args)
	}
	return args, errors.WithStack(err)
}

// FquoteAs writes the quoted table and its maybe alias into w.
func (a alias) FquoteAs(w queryWriter) error {
	if a.DerivedTable != nil {
		w.WriteByte('(')
		if err := a.DerivedTable.toSQL(w); err != nil {
			return errors.WithStack(err)
		}
		w.WriteByte(')')
		w.WriteString(" AS ")
		Quoter.quote(w, a.Alias)
		return nil
	}

	qf := Quoter.FquoteAs
	if a.IsExpression {
		qf = Quoter.FquoteExprAs
	}
	qf(w, a.Name, a.Alias)

	if a.Sort == sortAscending {
		w.WriteString(" ASC")
	}
	if a.Sort == sortDescending {
		w.WriteString(" DESC")
	}
	return nil
}

// TODO(CyS) if we need to distinguish between table name and the column or even need
// a sub select in the column list, then we can implement type aliases and replace
// all []string with type aliases. This costs some allocs but for modifying queries
// in dispatched events, it's getting easier ...
type aliases []alias

func (as aliases) FquoteAs(w queryWriter) error {
	for i, a := range as {
		if i > 0 {
			w.WriteString(", ")
		}
		if err := a.FquoteAs(w); err != nil {
			return errors.WithStack(err)
		}
	}
	return nil
}

func (as aliases) appendArgs(args Arguments) (Arguments, error) {
	for _, a := range as {
		var err error
		args, err = a.appendArgs(args)
		if err != nil {
			return nil, errors.WithStack(err)
		}
	}
	return args, nil
}

// setSort applies to last n items the sort order `sort` in reverse iteration.
// Usuallay `lastNindexes` is len(object) because we decrement 1 from
// `lastNindexes`. This function panics when lastNindexes does not match the
// length of `aliases`.
func (as aliases) applySort(lastNindexes int, sort byte) aliases {
	to := len(as) - lastNindexes
	for i := len(as) - 1; i >= to; i-- {
		as[i].Sort = sort
	}
	return as
}

// AddColumns adds more columns to the aliases. Columns get quoted.
func (as aliases) AddColumns(columns ...string) aliases {
	return as.appendColumns(columns, false)
}

func (as aliases) appendColumns(columns []string, isExpression bool) aliases {
	if len(as) == 0 {
		as = make(aliases, 0, len(columns))
	}
	for _, c := range columns {
		as = append(as, alias{Name: c, IsExpression: isExpression})
	}
	return as
}

// columns must be balanced slice. i=column name, i+1=alias name
func (as aliases) appendColumnsAliases(columns []string, isExpression bool) aliases {
	if len(as) == 0 {
		as = make(aliases, 0, len(columns)/2)
	}
	for i := 0; i < len(columns); i = i + 2 {
		as = append(as, alias{Name: columns[i], Alias: columns[i+1], IsExpression: isExpression})
	}
	return as
}

// MysqlQuoter implements Mysql-specific quoting
type MysqlQuoter struct {
	replacer *strings.Replacer
}

func (q MysqlQuoter) unQuote(s string) string {
	if strings.IndexByte(s, quoteRune) == -1 {
		return s
	}
	return q.replacer.Replace(s)
}

func (q MysqlQuoter) quote(w queryWriter, qualifierName ...string) {
	for i, qn := range qualifierName {
		if i > 0 {
			w.WriteByte('.')
		}
		w.WriteByte(quoteRune)
		w.WriteString(q.unQuote(qn))
		w.WriteByte(quoteRune)
	}
}

// exprAlias appends to the provided `expression` the quote alias name, e.g.:
// 		exprAlias("(e.price*x.tax*t.weee)", "final_price") // (e.price*x.tax*t.weee) AS `final_price`
func (q MysqlQuoter) exprAlias(expression, aliasName string) string {
	if aliasName == "" {
		return expression
	}
	return expression + " AS " + quote + q.unQuote(aliasName) + quote
}

// Quote quotes an optional qualifier and its required name. Returns a string
// like: `database`.`table` or `table`, if qualifier has been omitted.
// 		Quote("dbName", "tableName") => `dbName`.`tableName`
// 		Quote("tableName") => `tableName`
// It panics when no arguments have been given.
// https://dev.mysql.com/doc/refman/5.7/en/identifier-qualifiers.html
func (q MysqlQuoter) Quote(qualifierName ...string) string {
	// way faster than fmt or buffer ...
	l := len(qualifierName)
	idx1 := 0
	if l == 2 && qualifierName[idx1] == "" {
		idx1++
	}
	if l == 1 || idx1 == 1 {
		return quote + q.unQuote(qualifierName[idx1]) + quote
	}
	return quote + q.unQuote(qualifierName[0]) + quote + "." + quote + q.unQuote(qualifierName[1]) + quote
}

// QuoteAs quotes with back ticks and splits at a dot in the name. First
// argument table and/or column name (separated by a dot) and second argument
// can be an alias. Both parts will get quoted. If providing only one part, then
// the last `alias` parts gets skipped.
//		QuoteAs("f", "g", "h") 			// "`f` AS `g_h`"
//		QuoteAs("e.entity_id", "ee") 	// `e`.`entity_id` AS `ee`
func (q MysqlQuoter) QuoteAs(expressionAlias ...string) string {
	buf := bufferpool.Get()
	q.FquoteAs(buf, expressionAlias...)
	x := buf.String()
	bufferpool.Put(buf)
	return x
}

// FquoteExprAs quotes an expression with an optional alias into w.
func (q MysqlQuoter) FquoteExprAs(w queryWriter, expressionAlias ...string) {
	w.WriteString(expressionAlias[0])
	if len(expressionAlias) > 1 && expressionAlias[1] != "" {
		w.WriteString(" AS ")
		q.quote(w, expressionAlias[1])
	}
}

// FquoteAs same as QuoteAs but writes into w which is a bytes.Buffer. It quotes always and each part.
func (q MysqlQuoter) FquoteAs(w queryWriter, expressionAlias ...string) {

	lp := len(expressionAlias)
	if lp == 2 && expressionAlias[1] == "" {
		lp = 1
		expressionAlias = expressionAlias[:1]
	}
	expr := expressionAlias[0]

	// checks if there are quotes at the beginning and at the end. no white spaces allowed.
	hasQuote0 := strings.HasPrefix(expr, quote) && strings.HasSuffix(expr, quote)
	hasDot0 := strings.IndexByte(expr, '.') >= 0

	//fmt.Printf("lp %d expr %q hasQuote0 %t hasDot0 %t | %#v\n", lp, expr, hasQuote0, hasDot0, expressionAlias)

	switch {
	case lp == 1 && hasQuote0:
		// already quoted
		w.WriteString(expr)
		return
	case lp > 1 && expressionAlias[1] == "" && !hasQuote0 && !hasDot0:
		// must be quoted
		q.quote(w, expr)
		return
	case lp == 1 && !hasQuote0 && hasDot0:
		q.splitDotAndQuote(w, expr)
		return
	case lp == 1 && expr == "":
		// just an empty string
		return
	}

	q.splitDotAndQuote(w, expr)
	switch lp {
	case 1:
		// do nothing
	case 2:
		w.WriteString(" AS ")
		q.quote(w, expressionAlias[1])
	default:
		w.WriteString(" AS ")
		q.quote(w, strings.Join(expressionAlias[1:], "_"))
	}
	return
}

func (q MysqlQuoter) splitDotAndQuote(w queryWriter, part string) {
	dotIndex := strings.IndexByte(part, '.')
	if dotIndex > 0 { // dot at a beginning of a string at illegal
		q.quote(w, part[:dotIndex])
		w.WriteByte('.')
		if a := part[dotIndex+1:]; a == sqlStar {
			w.WriteByte('*')
		} else {
			q.quote(w, part[dotIndex+1:])
		}
		return
	}
	q.quote(w, part)
}

// TableColumnAlias prefixes all columns with a table name/alias and puts quotes around them.
// If a column name has already been prefixed by a name or an alias it will be ignored.
func (q MysqlQuoter) TableColumnAlias(t string, cols ...string) []string {
	for i, c := range cols {
		switch {
		case strings.ContainsRune(c, quoteRune):
			cols[i] = c
		case strings.ContainsRune(c, '.'):
			cols[i] = q.QuoteAs(c)
		default:
			cols[i] = q.QuoteAs(t + "." + c)
		}
	}
	return cols
}

// maxIdentifierLength see http://dev.mysql.com/doc/refman/5.7/en/identifiers.html
const maxIdentifierLength = 64
const dummyQualifier = "X" // just a dummy value, can be optimized later

// IsValidIdentifier checks the permissible syntax for identifiers. Certain
// objects within MySQL, including database, table, index, column, alias, view,
// stored procedure, partition, tablespace, and other object names are known as
// identifiers. ASCII: [0-9,a-z,A-Z$_] (basic Latin letters, digits 0-9, dollar,
// underscore) Max length 63 characters.
//
// Returns 0 if the identifier is valid.
//
// http://dev.mysql.com/doc/refman/5.7/en/identifiers.html
func isValidIdentifier(objectName string) int8 {
	if objectName == sqlStar {
		return 0
	}
	qualifier := dummyQualifier
	if i := strings.IndexByte(objectName, '.'); i >= 0 {
		qualifier = objectName[:i]
		objectName = objectName[i+1:]
	}

	validQualifier := isNameValid(qualifier)
	if validQualifier == 0 && objectName == sqlStar {
		return 0
	}
	if validQualifier > 0 {
		return validQualifier
	}
	return isNameValid(objectName)
}

// isNameValid returns 0 if the name is valid or an error number identifying
// where the name becomes invalid.
func isNameValid(name string) int8 {
	if name == dummyQualifier {
		return 0
	}

	ln := len(name)
	if ln > maxIdentifierLength || name == "" {
		return 1 //errors.NewNotValidf("[csdb] Incorrect identifier. Too long or empty: %q", name)
	}
	pos := 0
	for pos < ln {
		r, w := utf8.DecodeRuneInString(name[pos:])
		pos += w
		if !mapAlNum(r) {
			return 2 // errors.NewNotValidf("[csdb] Invalid character in name %q", name)
		}
	}
	return 0
}

func mapAlNum(r rune) bool {
	var ok bool
	switch {
	case '0' <= r && r <= '9':
		ok = true
	case 'a' <= r && r <= 'z', 'A' <= r && r <= 'Z':
		ok = true
	case r == '$', r == '_':
		ok = true
	}
	return ok
}
