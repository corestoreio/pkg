package dbr

import "strings"

// Stmt is helper for various method to check statements
var Stmt = stmtChecker{}

// stmtChecker @todo better checking ...
type stmtChecker struct{}

func (stmtChecker) startContain(sql, starts, contains string) bool {
	sql = strings.ToLower(sql)
	return strings.Index(sql, starts) == 0 && strings.Index(sql, contains) > 4
}

// IsSelect checks if string is a SELECT statement
func (sc stmtChecker) IsSelect(sql string) bool {
	return sc.startContain(sql, "select", "from")
}

// IsUpdate checks if string is an UPDATE statement
func (sc stmtChecker) IsUpdate(sql string) bool {
	return sc.startContain(sql, "update", "from")
}

// IsDelete checks if string is a DELETE statement
func (sc stmtChecker) IsDelete(sql string) bool {
	return sc.startContain(sql, "delete", "from")
}

// IsInsert checks if string is an INSERT statement
func (sc stmtChecker) IsInsert(sql string) bool {
	return sc.startContain(sql, "insert", " ")
}

func strInSlice(search string, sl []string) bool {
	for _, s := range sl {
		if s == search {
			return true
		}
	}
	return false
}
