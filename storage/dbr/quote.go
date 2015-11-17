package dbr

import "strings"

const quote string = "`"
const quoteRune rune = '`'

// Quoter is the quoter to use for quoting text; use Mysql quoting by default
var Quoter = MysqlQuoter{}

// MysqlQuoter implements Mysql-specific quoting
type MysqlQuoter struct{}

func (q MysqlQuoter) writeQuotedColumn(column string, sql QueryWriter) {
	_, _ = sql.WriteRune(quoteRune)
	_, _ = sql.WriteString(column)
	_, _ = sql.WriteRune(quoteRune)
}

func (q MysqlQuoter) unQuote(s string) string {
	if !strings.ContainsRune(s, quoteRune) {
		return s
	}
	return strings.Replace(s, quote, "", -1)
}

// QuoteAs quotes a with back ticks. First argument table or column name and
// second argument can be an alias. Both parts will get quoted.
func (q MysqlQuoter) QuoteAs(exprAlias ...string) string {
	return q.quoteAs(exprAlias...)
}

func (q MysqlQuoter) Alias(expression, as string) string {
	return expression + " AS " + quote + q.unQuote(as) + quote
}

func (q MysqlQuoter) quoteAs(parts ...string) string {

	lp := len(parts)

	hasQuote0 := strings.ContainsRune(parts[0], quoteRune)
	hasDot0 := strings.ContainsRune(parts[0], '.')

	switch {
	case lp == 1 && hasQuote0:
		return parts[0] // already quoted
	case lp > 1 && parts[1] == "" && !hasQuote0 && !hasDot0:
		return quote + q.unQuote(parts[0]) + quote // must be quoted
	case lp == 1 && !hasQuote0 && hasDot0:
		return q.splitDotAndQuote(parts[0])
	}

	n := q.splitDotAndQuote(parts[0])

	switch lp {
	case 1:
		return n
	case 2:
		return n + " AS " + quote + q.unQuote(parts[1]) + quote
	default:
		return n + " AS " + quote + q.unQuote(strings.Join(parts[1:], "_")) + quote
	}
}

func (q MysqlQuoter) splitDotAndQuote(part string) string {
	dotIndex := strings.Index(part, ".")
	if dotIndex > 0 { // dot at a beginning of a string is illegal
		return quote + q.unQuote(part[:dotIndex]) + quote + "." + quote + q.unQuote(part[dotIndex+1:]) + quote
	}
	return quote + q.unQuote(part) + quote
}

// ColumnAlias is a helper func which transforms variadic arguments into a slice with a special
// converting case that every i%2 index is considered as the alias
func (q MysqlQuoter) ColumnAlias(columns ...string) []string {
	l := len(columns)
	if l%2 == 1 {
		panic("Amount of columns must be even and not odd.")
	}
	cols := make([]string, l/2)
	j := 0
	for i := 0; i < l; i = i + 2 {
		cols[j] = q.quoteAs(columns[i], columns[i+1])
		j++
	}
	return cols
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
