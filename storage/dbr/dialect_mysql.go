package dbr

import (
	"strconv"
	"strings"
	"time"
)

const mysqlTimeFormat = "2006-01-02 15:04:05"

const DriverNameMySQL = "mysql"

type mysqlDialect struct{}

func (mysqlDialect) EscapeIdent(w QueryWriter, ident string) {
	w.WriteRune('`')
	r := strings.NewReplacer("`", "``", ".", "`.`")
	w.WriteString(r.Replace(ident))
	w.WriteRune('`')
}

func (mysqlDialect) EscapeBool(w QueryWriter, b bool) {
	if b {
		w.WriteRune('1')
	} else {
		w.WriteRune('0')
	}
}

// Need to turn \x00, \n, \r, \, ', " and \x1a.
// Returns an escaped, quoted string. eg, "hello 'world'" -> "'hello \'world\''".
func (mysqlDialect) EscapeString(w QueryWriter, s string) {
	w.WriteRune('\'')
	for _, char := range s {
		switch char {
		case '\'':
			w.WriteString(`\'`)
		case '"':
			w.WriteString(`\"`)
		case '\\':
			w.WriteString(`\\`)
		case '\n':
			w.WriteString(`\n`)
		case '\r':
			w.WriteString(`\r`)
		case 0:
			w.WriteString(`\x00`)
		case 0x1a:
			w.WriteString(`\x1a`)
		default:
			w.WriteRune(char)
		}
	}
	w.WriteRune('\'')
}

func (d mysqlDialect) EscapeTime(w QueryWriter, t time.Time) {
	d.EscapeString(w, t.Format(mysqlTimeFormat))
}

func (mysqlDialect) ApplyLimitAndOffset(w QueryWriter, limit, offset uint64) {
	w.WriteString(" LIMIT ")
	if limit == 0 {
		// In MYSQL, OFFSET cannot be used alone. Set the limit to the max possible value.
		w.WriteString("18446744073709551615")
	} else {
		w.WriteString(strconv.FormatUint(limit, 10))
	}
	if offset > 0 {
		w.WriteString(" OFFSET ")
		w.WriteString(strconv.FormatUint(offset, 10))
	}

}
