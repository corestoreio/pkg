package dbr

import "strings"

const Quote string = "`"

// Quoter is the quoter to use for quoting text; use Mysql quoting by default
var Quoter = MysqlQuoter{}

// Interface for driver-swappable quoting
type quoter interface {
	writeQuotedColumn()
}

// MysqlQuoter implements Mysql-specific quoting
type MysqlQuoter struct{}

func (q MysqlQuoter) writeQuotedColumn(column string, sql QueryWriter) {
	sql.WriteString(Quote + column + Quote)
}

// Table quotes a table with back ticks. First argument table name, second argument
// can be an alias.
func (q MysqlQuoter) Table(tableName ...string) string {
	//	tn := ""
	//	for _, n := range tableName {
	//		if n != '`' { // @todo
	//			tn = tn + n
	//		}
	//	}
	if len(tableName) == 1 {
		return Quote + tableName[0] + Quote
	}
	return quoteAs(tableName...)
}

func quoteAs(parts ...string) string {
	if len(parts) == 1 {
		return parts[0]
	}
	if len(parts) != 2 {
		panic("from can either be a table name or table name and alias")
	}
	n := parts[0]
	dotIndex := strings.Index(n, ".")
	if dotIndex > 0 {
		n = Quote + parts[0][:dotIndex] + Quote
		n = n + "." + Quote + parts[0][dotIndex+1:] + Quote
	} else {
		n = Quote + n + Quote
	}
	return n + " AS " + Quote + parts[1] + Quote
}

// ColumnAlias is a helper func which transforms variadic arguments into a slice with a special
// converting case that every i%2 index is considered as the alias
func ColumnAlias(columns ...string) []string {
	l := len(columns)
	if l%2 == 1 {
		panic("Amount of columns must be even and not odd.")
	}
	cols := make([]string, l/2)
	j := 0
	for i := 0; i < l; i = i + 2 {
		cols[j] = quoteAs(columns[i], columns[i+1])
		j++
	}
	return cols
}
