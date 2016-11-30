package dbr

import (
	"fmt"
	"io"
)

// QueryWriter is used to generate a query.
type QueryWriter interface {
	io.Writer
	fmt.Stringer
	WriteString(s string) (n int, err error)
	WriteRune(r rune) (n int, err error)
}
