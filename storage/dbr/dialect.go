package dbr

import (
	"encoding/hex"
	"strings"
	"time"
)

// no other dialect will be supported. Maybe MariaDB ;-)

var dialect dialecter = mysqlDialect{
	identR: strings.NewReplacer("`", "``", ".", "`.`"),
}

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

type mysqlDialect struct {
	identR *strings.Replacer
}

func (d mysqlDialect) EscapeIdent(w queryWriter, ident string) {
	w.WriteByte('`')
	w.WriteString(d.identR.Replace(ident))
	w.WriteByte('`')
}

func (d mysqlDialect) EscapeBool(w queryWriter, b bool) {
	if b {
		w.WriteByte('1')
	} else {
		w.WriteByte('0')
	}
}

func (d mysqlDialect) EscapeBinary(w queryWriter, b []byte) {
	if b == nil {
		w.WriteString("NULL")
	} else {
		// TODO(CyS) no idea if that at the correct way. do an RTFM
		w.WriteString("0x")
		w.WriteString(hex.EncodeToString(b))
	}
}

// EscapeString. Need to turn \x00, \n, \r, \, ', " and \x1a.
// Returns an escaped, quoted string. eg, "hello 'world'" -> "'hello \'world\''".
func (d mysqlDialect) EscapeString(w queryWriter, s string) {
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
	if t.IsZero() {
		w.WriteString("'0000-00-00'") //  00:00:00
	} else {
		// time.Location must be considered ...
		w.WriteByte('\'')
		d := w.Bytes()
		w.Reset()
		w.Write(t.AppendFormat(d, mysqlTimeFormat))
		w.WriteByte('\'')
	}
	// d.EscapeString(w, t.Format(mysqlTimeFormat))
}

func (d mysqlDialect) ApplyLimitAndOffset(w queryWriter, limit, offset uint64) {
	w.WriteString(" LIMIT ")
	if limit == 0 {
		// In MYSQL, OFFSET cannot be used alone. Set the limit to the max possible value.
		w.WriteString("18446744073709551615")
	} else {
		writeUint64(w, limit)
	}
	if offset > 0 {
		w.WriteString(" OFFSET ")
		writeUint64(w, offset)
	}
}
