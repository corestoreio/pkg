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

// FquoteAs same as QuoteAs but writes into w which is a bytes.Buffer.
func (q MysqlQuoter) FquoteAs(w queryWriter, expressionAlias ...string) {

	lp := len(expressionAlias)
	if lp == 2 && expressionAlias[1] == "" {
		lp = 1
		expressionAlias = expressionAlias[:1]
	}

	hasQuote0 := strings.ContainsRune(expressionAlias[0], quoteRune)
	hasDot0 := strings.ContainsRune(expressionAlias[0], '.')

	switch {
	case lp == 1 && hasQuote0:
		// already quoted
		w.WriteString(expressionAlias[0])
		return
	case lp > 1 && expressionAlias[1] == "" && !hasQuote0 && !hasDot0:
		// must be quoted
		q.quote(w, expressionAlias[0])
		return
	case lp == 1 && !hasQuote0 && hasDot0:
		q.splitDotAndQuote(w, expressionAlias[0])
		return
	case lp == 1 && expressionAlias[0] == "":
		// just an empty string
		return
	}

	q.splitDotAndQuote(w, expressionAlias[0])
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
