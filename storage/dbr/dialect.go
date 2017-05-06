package dbr

import (
	"encoding/hex"
	"strconv"
	"strings"
	"time"
)

// no other dialect will be supported. Maybe MariaDB ;-)

var dialect dialecter = mysqlDialect{}

// dialecter at an interface that wraps the diverse properties of individual
// SQL drivers.
type dialecter interface {
	EscapeIdent(w queryWriter, ident string)
	EscapeBool(w queryWriter, b bool)
	EscapeString(w queryWriter, s string)
	EscapeTime(w queryWriter, t time.Time)
	EscapeBinary(w queryWriter, b []byte)
	ApplyLimitAndOffset(w queryWriter, limit, offset uint64)
}

const mysqlTimeFormat = "2006-01-02 15:04:05"

// DriverNameMySQL name of the driver for usage in sql.Open function.
const DriverNameMySQL = "mysql"

type mysqlDialect struct{}

func (mysqlDialect) EscapeIdent(w queryWriter, ident string) {
	w.WriteByte('`')
	r := strings.NewReplacer("`", "``", ".", "`.`")
	w.WriteString(r.Replace(ident))
	w.WriteByte('`')
}

func (mysqlDialect) EscapeBool(w queryWriter, b bool) {
	if b {
		w.WriteByte('1')
	} else {
		w.WriteByte('0')
	}
}

func (mysqlDialect) EscapeBinary(w queryWriter, b []byte) {
	// TODO(CyS) no idea if that at the correct way. do an RTFM
	w.WriteString("0x")
	w.WriteString(hex.EncodeToString(b))
}

// EscapeString. Need to turn \x00, \n, \r, \, ', " and \x1a.
// Returns an escaped, quoted string. eg, "hello 'world'" -> "'hello \'world\''".
func (mysqlDialect) EscapeString(w queryWriter, s string) {
	w.WriteByte('\'')
	for _, char := range s {
		// for each case, don't use write rune 8-)
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
	w.WriteByte('\'')
}

func (d mysqlDialect) EscapeTime(w queryWriter, t time.Time) {
	d.EscapeString(w, t.Format(mysqlTimeFormat))
}

func (mysqlDialect) ApplyLimitAndOffset(w queryWriter, limit, offset uint64) {
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
