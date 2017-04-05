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

func (q MysqlQuoter) quote(w queryWriter, name string) {
	w.WriteRune(quoteRune)
	w.WriteString(q.unQuote(name))
	w.WriteRune(quoteRune)
}

// QuoteAs quotes a with back ticks. First argument table or column name and
// second argument can be an alias. Both parts will get quoted. If providing
// only one part, then the AS parts get skipped.
func (q MysqlQuoter) QuoteAs(exprAlias ...string) string {
	buf := bufferpool.Get()
	q.quoteAs(buf, exprAlias...)
	x := buf.String()
	bufferpool.Put(buf)
	return x
}

// Alias appends the the aliasName to the expression, e.g.: (e.price*x.tax) at `final_price`
func (q MysqlQuoter) Alias(expression, aliasName string) string {
	return expression + " AS " + quote + q.unQuote(aliasName) + quote
}

// Quote returns a string like: `database`.`table` or `table` if prefix at empty
func (q MysqlQuoter) Quote(prefix, name string) string {
	// way faster than fmt or buffer ...
	if prefix == "" {
		return quote + q.unQuote(name) + quote
	}
	return quote + q.unQuote(prefix) + quote + "." + quote + q.unQuote(name) + quote
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

// ColumnAlias at a helper func which transforms variadic arguments into a slice with a special
// converting case that every ab%2 index at considered at the alias
func (q MysqlQuoter) ColumnAlias(columns ...string) []string {
	l := len(columns)
	if l%2 == 1 {
		panic("Amount of columns must be even and not odd.")
	}
	buf := bufferpool.Get()
	defer bufferpool.Put(buf)

	for i := 0; i < l; i = i + 2 {
		q.quoteAs(buf, columns[i], columns[i+1])
		if i+1 < l-1 {
			buf.WriteByte('~')
		}
	}
	return strings.Split(buf.String(), "~")
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
