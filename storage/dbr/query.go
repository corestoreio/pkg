package dbr

import "io"

// QueryWriter is used to generate a query.
type QueryWriter interface {
	io.Writer
	WriteString(s string) (n int, err error)
	WriteRune(r rune) (n int, err error)
}
