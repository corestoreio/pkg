package dbr

import "time"

var D Dialect = Mysql{}

// Dialect is an interface that wraps the diverse properties of individual
// SQL drivers.
type Dialect interface {
	EscapeIdent(w QueryWriter, ident string)
	EscapeBool(w QueryWriter, b bool)
	EscapeString(w QueryWriter, s string)
	EscapeTime(w QueryWriter, t time.Time)
	ApplyLimitAndOffset(w QueryWriter, limit, offset uint64)
}
