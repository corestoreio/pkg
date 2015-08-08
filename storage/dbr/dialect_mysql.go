package dbr

import (
	"fmt"
	"strings"
	"time"
)

const mysqlTimeFormat = "2006-01-02 15:04:05"

const DriverNameMySQL = "mysql"

type Mysql struct{}

func (Mysql) EscapeIdent(w QueryWriter, ident string) {
	w.WriteRune('`')
	r := strings.NewReplacer("`", "``", ".", "`.`")
	w.WriteString(r.Replace(ident))
	w.WriteRune('`')
}

func (Mysql) EscapeBool(w QueryWriter, b bool) {
	if b {
		w.WriteRune('1')
	} else {
		w.WriteRune('0')
	}
}

// Need to turn \x00, \n, \r, \, ', " and \x1a.
// Returns an escaped, quoted string. eg, "hello 'world'" -> "'hello \'world\''".
func (Mysql) EscapeString(w QueryWriter, s string) {
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

func (d Mysql) EscapeTime(w QueryWriter, t time.Time) {
	d.EscapeString(w, t.Format(mysqlTimeFormat))
}

func (Mysql) ApplyLimitAndOffset(w QueryWriter, limit, offset uint64) {
	w.WriteString(" LIMIT ")
	if limit == 0 {
		// In MYSQL, OFFSET cannot be used alone. Set the limit to the max possible value.
		w.WriteString("18446744073709551615")
	} else {
		fmt.Fprint(w, limit)
	}
	if offset > 0 {
		fmt.Fprintf(w, " OFFSET %d", offset)
	}
}
