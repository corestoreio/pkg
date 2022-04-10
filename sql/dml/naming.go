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
	"fmt"
	"strings"
	"unicode"
	"unicode/utf8"

	"github.com/corestoreio/pkg/util/bufferpool"
)

const (
	quote     = "`"
	quoteRune = '`'
)

// Quoter at the quoter to use for quoting text; use Mysql quoting by default.
var Quoter = MysqlQuoter{
	replacer: strings.NewReplacer(quote, ""),
}

// id is an identifier for table name or a column name or an alias for a sub
// query.
type id struct {
	// Derived Tables (Subqueries in the FROM Clause). A derived table is a
	// subquery in a SELECT statement FROM clause. Derived tables can return a
	// scalar, column, row, or table. Ignored in any other case.
	DerivedTable *Select
	// Name can be any kind of SQL expression or a valid identifier. It gets
	// quoted when `IsLeftExpression` is false.
	Name string
	// Expression has precedence over the `Name` field. Each line in an expression
	// gets written unchanged to the final SQL string.
	Expression string
	// Aliased must be a valid identifier allowed for alias usage. As soon as the field `Aliased` has been set
	// it gets append to the Name and Expression field: "sql AS Aliased"
	Aliased string
	// Sort applies only to GROUP BY and ORDER BY clauses. 'd'=descending,
	// 0=default or nothing; 'a'=ascending.
	Sort byte
}

const (
	sortDescending byte = 'd'
	sortAscending  byte = 'a'
)

// MakeIdentifier creates a new quoted name with an optional alias `a`, which can be
// empty.
func MakeIdentifier(name string) id { return id{Name: name} }

// Alias sets the aliased name for the `Name` field.
func (a id) Alias(alias string) id { a.Aliased = alias; return a }

// Clone creates a new object and takes care of a cloned DerivedTable field.
func (a id) Clone() id {
	if nil != a.DerivedTable {
		a.DerivedTable = a.DerivedTable.Clone()
	}
	return a
}

// uncomment this functions and its test once needed
// MakeExpressionAlias creates a new unquoted expression with an optional alias
// `a`, which can be empty.
// func MakeExpressionAlias(expression []string, a string) identifier {
//	return identifier{
//		Expression: expression,
//		Aliased:      a,
//	}
//}

func (a id) isEmpty() bool { return a.Name == "" && a.DerivedTable == nil && a.Expression == "" }

// qualifier returns the correct qualifier for an identifier
func (a id) qualifier() string {
	if a.Aliased != "" {
		return a.Aliased
	}
	return a.Name
}

// String returns the correct stringyfied statement.
func (a id) String() string {
	if a.Expression != "" {
		buf := bufferpool.Get()
		defer bufferpool.Put(buf)
		buf.WriteString(a.Expression)
		buf.WriteString(" AS ")
		Quoter.quote(buf, a.Aliased)
		return buf.String()
	}
	return a.QuoteAs()
}

// NameAlias always quuotes the name and the alias
func (a id) QuoteAs() string { return Quoter.NameAlias(a.Name, a.Aliased) }

// writeQuoted writes the quoted table and its maybe alias into w.
func (a id) writeQuoted(w *bytes.Buffer, placeHolders []string) (_ []string, err error) {
	if a.DerivedTable != nil {
		w.WriteByte('(')
		if placeHolders, err = a.DerivedTable.toSQL(w, placeHolders); err != nil {
			return nil, fmt.Errorf("[dml] 1648324312088 writeQuoted failed: %w", err)
		}
		w.WriteByte(')')
		w.WriteString(" AS ")
		Quoter.quote(w, a.Aliased)
		return placeHolders, nil
	}

	if a.Expression != "" {
		writeExpression(w, a.Expression, nil)
	} else {
		Quoter.WriteIdentifier(w, a.Name)
	}
	if a.Aliased != "" {
		w.WriteString(" AS ")
		Quoter.quote(w, a.Aliased)
	}

	if a.Sort == sortAscending {
		w.WriteString(" ASC")
	}
	if a.Sort == sortDescending {
		w.WriteString(" DESC")
	}
	return placeHolders, nil
}

// ids is a slice of identifiers. `idc` in the receiver means id-collection.
type ids []id

func (idc ids) Clone() ids {
	if idc == nil {
		return nil
	}
	c := make(ids, len(idc))
	for idx, ido := range idc {
		c[idx] = ido.Clone()
	}
	return c
}

// writeQuoted writes all identifiers comma separated and quoted into w.
func (idc ids) writeQuoted(w *bytes.Buffer, placeHolders []string) (_ []string, err error) {
	for i, a := range idc {
		if i > 0 {
			w.WriteString(", ")
		}
		if placeHolders, err = a.writeQuoted(w, placeHolders); err != nil {
			return nil, fmt.Errorf("[dml] 1648324339117 %w", err)
		}
	}
	return placeHolders, nil
}

// setSort applies to last n items the sort order `sort` in reverse iteration.
// Usuallay `lastNindexes` is len(object) because we decrement 1 from
// `lastNindexes`. This function panics when lastNindexes does not match the
// length of `identifiers`.
func (idc ids) applySort(lastNindexes int, sort byte) ids {
	to := len(idc) - lastNindexes
	for i := len(idc) - 1; i >= to; i-- {
		idc[i].Sort = sort
	}
	return idc
}

// AppendColumns adds new columns to the identifier slice. If a column name is
// not a valid identifier that column gets switched into an expression. You
// should use this function when no arguments should be attached to an
// expression, otherwise use the function appendConditions. If a column name
// contains " ASC" or " DESC" as a suffix, the internal sorting flag gets set
// and the words ASC or DESC removed.
func (idc ids) AppendColumns(isUnsafe bool, columns ...string) ids {
	if cap(idc) == 0 {
		idc = make(ids, 0, len(columns)*2)
	}
	for _, c := range columns {
		var sorting byte
		switch {
		case strings.HasSuffix(c, " ASC"):
			c = c[:len(c)-4]
			sorting = sortAscending
		case strings.HasSuffix(c, " DESC"):
			c = c[:len(c)-5]
			sorting = sortDescending
		}

		ident := id{Name: c, Sort: sorting}
		if isUnsafe && isValidIdentifier(c) != 0 {
			ident.Expression = ident.Name
			ident.Name = ""
		}
		idc = append(idc, ident)
	}
	return idc
}

// AppendColumnsAliases expects a balanced slice where i=column name and
// i+1=alias name. An imbalanced slice will cause a panic. If a column name is
// not valid identifier that column gets switched into an expression. The alias
// does not change. You should use this function when no arguments should be
// attached to an expression, otherwise use the function appendConditions.
func (idc ids) AppendColumnsAliases(isUnsafe bool, columns ...string) ids {
	if (len(columns) % 2) == 1 {
		// A programmer made an error
		panic(fmt.Errorf("[dml] Expecting a balanced slice! Got: %v", columns))
	}
	if cap(idc) == 0 {
		idc = make(ids, 0, len(columns)/2)
	}

	for i := 0; i < len(columns); i = i + 2 {
		ident := id{Name: columns[i], Aliased: columns[i+1]}
		if isUnsafe && isValidIdentifier(ident.Name) != 0 {
			ident.Expression = ident.Name
			ident.Name = ""
		}
		idc = append(idc, ident)
	}
	return idc
}

// appendConditions adds an expression with arguments. SubSelects are not yet
// supported. You should use this function when arguments should be attached to
// the expression, otherwise use the function AppendColumns*.
func (idc ids) appendConditions(expressions Conditions) (ids, error) {
	buf := bufferpool.Get()
	for _, e := range expressions {
		idf := id{Name: e.Left, Aliased: e.Aliased}
		if e.IsLeftExpression {
			idf.Expression = idf.Name
			idf.Name = ""

			if len(e.Right.args) > 0 {
				if err := writeInterpolate(buf, idf.Expression, e.Right.args); err != nil {
					bufferpool.Put(buf)
					return nil, fmt.Errorf("[dml] 1648324526672 ids.appendConditions with expression: %q error: %w", idf.Expression, err)
				}
				idf.Expression = buf.String()
				buf.Reset()
			}
		}
		idc = append(idc, idf)
	}
	bufferpool.Put(buf)
	return idc, nil
}

// MysqlQuoter implements Mysql-specific quoting
type MysqlQuoter struct {
	replacer *strings.Replacer
}

func (mq MysqlQuoter) unQuote(s string) string {
	if strings.IndexByte(s, quoteRune) == -1 {
		return s
	}
	return mq.replacer.Replace(s)
}

func (mq MysqlQuoter) quote(w *bytes.Buffer, str string) {
	w.WriteByte(quoteRune)
	w.WriteString(mq.unQuote(str))
	w.WriteByte(quoteRune)
}

func (mq MysqlQuoter) writeQualifierName(w *bytes.Buffer, q, n string) {
	mq.quote(w, q)
	w.WriteByte('.')
	mq.quote(w, n)
}

// Name quotes securely a name.
// 		Name("tableName") => `tableName`
// 		Name("table`Name") => `tableName`
// https://dev.mysql.com/doc/refman/5.7/en/identifier-qualifiers.html
func (mq MysqlQuoter) Name(n string) string {
	return quote + mq.unQuote(n) + quote
}

// QualifierName quotes securely a qualifier and its name.
// 		QualifierName("dbName", "tableName") => `dbName`.`tableName`
// 		QualifierName("db`Name", "`tableName`") => `dbName`.`tableName`
// https://dev.mysql.com/doc/refman/5.7/en/identifier-qualifiers.html
func (mq MysqlQuoter) QualifierName(q, n string) string {
	if q == "" {
		return mq.Name(n)
	}
	// return mq.Name(q) + "." + mq.Name(n) <-- too slow, too many allocs
	return quote + mq.unQuote(q) + quote + "." + quote + mq.unQuote(n) + quote
}

// WriteQualifierName same as function QualifierName but writes into w.
func (mq MysqlQuoter) WriteQualifierName(w *bytes.Buffer, qualifier, name string) {
	if qualifier == "" {
		mq.quote(w, name)
		return
	}
	mq.quote(w, qualifier)
	w.WriteByte('.')
	mq.quote(w, name)
}

// NameAlias quotes with back ticks and splits at a dot into the qualified or
// unqualified identifier. First argument table and/or column name (separated by
// a dot) and second argument can be an alias. Both parts will get quoted.
//		NameAlias("f", "g") 			// "`f` AS `g`"
//		NameAlias("e.entity_id", "ee") 	// `e`.`entity_id` AS `ee`
//		NameAlias("e.entity_id", "") 	// `e`.`entity_id`
func (mq MysqlQuoter) NameAlias(name, alias string) string {
	buf := bufferpool.Get()
	mq.WriteIdentifier(buf, name)
	if alias != "" {
		buf.WriteString(" AS ")
		Quoter.quote(buf, alias)
	}
	x := buf.String()
	bufferpool.Put(buf)
	return x
}

// WriteIdentifier quotes with back ticks and splits at a dot into the qualified
// or unqualified identifier. First argument table and/or column name (separated
// by a dot). It quotes always and each part. If a string contains quotes, they
// won't get stripped.
//		WriteIdentifier(&buf,"tableName.ColumnName") -> `tableName`.`ColumnName`
func (mq MysqlQuoter) WriteIdentifier(w *bytes.Buffer, name string) {
	switch {
	case name == "":
		return
	case name == sqlStrNullUC: // see calling func sqlIfNullQuote2
		w.WriteString(name)
		return
	case strings.HasPrefix(name, quote) && strings.HasSuffix(name, quote): // not really secure
		// checks if there are quotes at the beginning and at the end. no white spaces allowed.
		w.WriteString(name) // already quoted
		return
	case strings.IndexByte(name, '.') == -1:
		// must be quoted
		mq.quote(w, name)
		return
	}

	if dotIndex := strings.IndexByte(name, '.'); dotIndex > 0 && dotIndex+1 < len(name) { // dot at a beginning of a string is illegal and at the end
		mq.quote(w, name[:dotIndex])
		w.WriteByte('.')
		if a := name[dotIndex+1:]; a == sqlStar {
			w.WriteByte('*')
		} else {
			mq.quote(w, name[dotIndex+1:])
		}
		return
	}
	mq.quote(w, name)
}

// ColumnsWithQualifier prefixes all columns in the slice `cols` with a qualifier and applies backticks. If a column name has already been
// prefixed with a qualifier or an alias it will be ignored. This functions modifies
// the argument slice `cols`.
func (mq MysqlQuoter) ColumnsWithQualifier(t string, cols ...string) []string {
	for i, c := range cols {
		switch {
		case strings.IndexByte(c, quoteRune) >= 0:
			cols[i] = c
		case strings.IndexByte(c, '.') > 0:
			cols[i] = mq.NameAlias(c, "")
		default:
			cols[i] = mq.QualifierName(t, c)
		}
	}
	return cols
}

// MaxIdentifierLength see http://dev.mysql.com/doc/refman/5.7/en/identifiers.html
const MaxIdentifierLength = 64

const dummyQualifier = "X" // just a dummy value, can be optimized later

// IsValidIdentifier checks the permissible syntax for identifiers. Certain
// objects within MySQL, including database, table, index, column, alias, view,
// stored procedure, partition, tablespace, and other object names are known as
// identifiers. ASCII: [0-9,a-z,A-Z$_] (basic Latin letters, digits 0-9, dollar,
// underscore) Max length 63 characters.
// It is recommended that you do not use names that begin with Me or MeN, where
// M and N are integers. For example, avoid using 1e as an identifier, because
// an expression such as 1e+3 is ambiguous. Depending on context, it might be
// interpreted as the expression 1e + 3 or as the number 1e+3.
//
// Returns 0 if the identifier is valid.
//
// http://dev.mysql.com/doc/refman/5.7/en/identifiers.html
func IsValidIdentifier(objectName string) (err error) {
	if v := isValidIdentifier(objectName); v != 0 {
		err = fmt.Errorf("[dml] 1648324659366 Invalid identifier %q (Case %d)", objectName, v)
	}
	return
}

func isValidIdentifier(s string) int8 {
	if s == sqlStar {
		return 0
	}
	qualifier := dummyQualifier
	if i := strings.IndexByte(s, '.'); i >= 0 {
		qualifier = s[:i]
		s = s[i+1:]
	}

	validQualifier := isNameValid(qualifier)
	if validQualifier == 0 && s == sqlStar {
		return 0
	}
	if validQualifier > 0 {
		return validQualifier
	}
	return isNameValid(s)
}

// isNameValid returns 0 if the name is valid or an error number identifying
// where the name becomes invalid.
func isNameValid(name string) int8 {
	if name == dummyQualifier {
		return 0
	}

	ln := len(name)
	if ln > MaxIdentifierLength || name == "" {
		return 1
	}
	pos := 0
	for pos < ln {
		r, w := utf8.DecodeRuneInString(name[pos:])
		if pos == 0 && unicode.IsDigit(r) {
			return 3
		}
		pos += w
		if !mapAlNum(r) {
			return 2
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

func cloneStringSlice(sl []string) []string {
	if sl == nil {
		return nil
	}
	c := make([]string, len(sl))
	copy(c, sl)
	return c
}
