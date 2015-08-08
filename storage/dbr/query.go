package dbr

import "io"

// Writer is used to write a query.
type QueryWriter interface {
	io.Writer
	WriteString(s string) (n int, err error)
	WriteRune(r rune) (n int, err error)
}
