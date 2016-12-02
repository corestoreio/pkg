package dbr

import "time"

var dialect dialecter = mysqlDialect{}

// dialecter is an interface that wraps the diverse properties of individual
// SQL drivers.
type dialecter interface {
	EscapeIdent(w QueryWriter, ident string)
	EscapeBool(w QueryWriter, b bool)
	EscapeString(w QueryWriter, s string)
	EscapeTime(w QueryWriter, t time.Time)
	ApplyLimitAndOffset(w QueryWriter, limit, offset uint64)
}
