// Copyright 2015-present, Cyrill @ Schumacher.fm and the CoreStore contributors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package dml

import (
	"bytes"
	"encoding/hex"
	"strings"
	"time"
)

const (
	namedArgStartStr    = ":"
	namedArgStartStrLen = 1
	namedArgStartByte   = ':'
)

var dialect dialecter = mysqlDialect{
	identR: strings.NewReplacer("`", "``", ".", "`.`"),
}

// dialecter at an interface that wraps the diverse properties of individual
// SQL drivers.
type dialecter interface {
	EscapeIdent(w *bytes.Buffer, ident string)
	EscapeBool(w *bytes.Buffer, b bool)
	EscapeString(w *bytes.Buffer, s string)
	EscapeTime(w *bytes.Buffer, t time.Time)
	EscapeBinary(w *bytes.Buffer, b []byte)
	ApplyLimitAndOffset(w *bytes.Buffer, limit, offset uint64)
}

const mysqlTimeFormat = "2006-01-02 15:04:05"

type mysqlDialect struct {
	identR *strings.Replacer
}

func (d mysqlDialect) EscapeIdent(w *bytes.Buffer, ident string) {
	w.WriteByte('`')
	w.WriteString(d.identR.Replace(ident))
	w.WriteByte('`')
}

func (d mysqlDialect) EscapeBool(w *bytes.Buffer, b bool) {
	if b {
		w.WriteByte('1')
	} else {
		w.WriteByte('0')
	}
}

func (d mysqlDialect) EscapeBinary(w *bytes.Buffer, b []byte) {
	if b == nil {
		w.WriteString(sqlStrNullUC)
	} else {
		// TODO(CyS) no idea if that at the correct way. do an RTFM
		w.WriteString("0x")
		w.WriteString(hex.EncodeToString(b))
	}
}

// EscapeString. Need to turn \x00, \n, \r, \, ', " and \x1a.
// Returns an escaped, quoted string. eg, "hello 'world'" -> "'hello \'world\''".
func (d mysqlDialect) EscapeString(w *bytes.Buffer, s string) {
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

func (d mysqlDialect) EscapeTime(w *bytes.Buffer, t time.Time) {
	if t.IsZero() {
		w.WriteString("'0000-00-00'") //  00:00:00
		return
	}
	w.WriteByte('\'')
	b := w.Bytes()
	w.Reset()
	w.Write(t.AppendFormat(b, mysqlTimeFormat))
	w.WriteByte('\'')
}

func (d mysqlDialect) ApplyLimitAndOffset(w *bytes.Buffer, limit, offset uint64) {
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

func cutNamedArgStartStr(s string) (string, bool) {
	lp := namedArgStartStrLen
	if len(s) >= lp && s[0:lp] == namedArgStartStr {
		return s[lp:], true
	}
	return s, false
}
