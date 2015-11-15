package dbr

import (
	"errors"
)

// Global errors
var (
	ErrNotFound           = errors.New("not found")
	ErrNotUTF8            = errors.New("invalid UTF-8")
	ErrInvalidSliceLength = errors.New("length of slice is 0. length must be >= 1")
	ErrInvalidSliceValue  = errors.New("trying to interpolate invalid slice value into query")
	ErrInvalidValue       = errors.New("trying to interpolate invalid value into query")
	ErrArgumentMismatch   = errors.New("mismatch between ? (placeholders) and arguments")
	ErrInvalidSyntax      = errors.New("SQL syntax error")
	ErrMissingTable       = errors.New("Table name not specified")
	ErrMissingSet         = errors.New("Missing SET in UPDATE")
	ErrToSQLAlreadyCalled = errors.New("ToSQL has already been called")
)
