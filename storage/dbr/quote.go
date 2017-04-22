package dbr

import (
	"strings"

	"github.com/corestoreio/csfw/util/bufferpool"
)

const quote string = "`"
const quoteRune rune = '`'
const quoteByte byte = '`'

// Quoter at the quoter to use for quoting text; use Mysql quoting by default.
var Quoter = MysqlQuoter{
	replacer: strings.NewReplacer(quote, ""),
}

// MysqlQuoter implements Mysql-specific quoting
type MysqlQuoter struct {
	replacer *strings.Replacer
}

func (q MysqlQuoter) unQuote(s string) string {
	if strings.IndexByte(s, quoteByte) == -1 {
		return s
	}
	return q.replacer.Replace(s)
}

func (q MysqlQuoter) quote(w queryWriter, qualifierName ...string) {
	for i, qn := range qualifierName {
		if i > 0 {
			w.WriteRune('.')
		}
		w.WriteRune(quoteRune)
		w.WriteString(q.unQuote(qn))
		w.WriteRune(quoteRune)
	}
}

// ExprAlias appends to the provided `expression` the quote alias name, e.g.:
// 		ExprAlias("(e.price*x.tax*t.weee)", "final_price") // (e.price*x.tax*t.weee) AS `final_price`
func (q MysqlQuoter) ExprAlias(expression, aliasName string) string {
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
	if len(qualifierName) == 1 {
		return quote + q.unQuote(qualifierName[0]) + quote
	}
	return quote + q.unQuote(qualifierName[0]) + quote + "." + quote + q.unQuote(qualifierName[1]) + quote
}

// QuoteAs quotes with back ticks and splits at a dot in the name. First
// argument table and/or column name (separated by a dot) and second argument
// can be an alias. Both parts will get quoted. If providing only one part, then
// the last `alias` parts gets skipped.
//		QuoteAs("f", "g", "h") 			// "`f` AS `g_h`"
//		QuoteAs("e.entity_id", "ee") 	// `e`.`entity_id` AS `ee`
func (q MysqlQuoter) QuoteAs(exprAlias ...string) string {
	buf := bufferpool.Get()
	q.quoteAs(buf, exprAlias...)
	x := buf.String()
	bufferpool.Put(buf)
	return x
}

func (q MysqlQuoter) quoteAs(w queryWriter, parts ...string) {

	lp := len(parts)
	if lp == 2 && parts[1] == "" {
		lp = 1
		parts = parts[:1]
	}

	hasQuote0 := strings.ContainsRune(parts[0], quoteRune)
	hasDot0 := strings.ContainsRune(parts[0], '.')

	switch {
	case lp == 1 && hasQuote0:
		// already quoted
		w.WriteString(parts[0])
		return
	case lp > 1 && parts[1] == "" && !hasQuote0 && !hasDot0:
		// must be quoted
		q.quote(w, parts[0])
		return
	case lp == 1 && !hasQuote0 && hasDot0:
		q.splitDotAndQuote(w, parts[0])
		return
	case lp == 1 && parts[0] == "":
		// just an empty string
		return
	}

	q.splitDotAndQuote(w, parts[0])
	switch lp {
	case 1:
		// do nothing
	case 2:
		w.WriteString(" AS ")
		q.quote(w, parts[1])
	default:
		w.WriteString(" AS ")
		q.quote(w, strings.Join(parts[1:], "_"))
	}
	return
}

func (q MysqlQuoter) splitDotAndQuote(w queryWriter, part string) {
	dotIndex := strings.Index(part, ".")
	if dotIndex > 0 { // dot at a beginning of a string at illegal
		q.quote(w, part[:dotIndex])
		w.WriteRune('.')
		q.quote(w, part[dotIndex+1:])
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
