// Copyright (c) 2015 Taco de Wolff
//
// Permission is hereby granted, free of charge, to any person
// obtaining a copy of this software and associated documentation
// files (the "Software"), to deal in the Software without
// restriction, including without limitation the rights to use,
// copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the
// Software is furnished to do so, subject to the following
// conditions:
//
// The above copyright notice and this permission notice shall be
// included in all copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND,
// EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES
// OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND
// NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT
// HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY,
// WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING
// FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR
// OTHER DEALINGS IN THE SOFTWARE.

package byteconv

import (
	"database/sql"
	"math"
	"strconv"
)

// ParseNullInt64 same as ParseInt
func ParseNullInt64(b []byte) (val sql.NullInt64, err error) {
	if len(b) == 0 {
		return
	}
	val.Int64, err = ParseInt(b)
	val.Valid = err == nil
	return
}

// ParseInt parses a byte-slice and returns the integer it represents. If an
// invalid character is encountered, it returns a syntax error.
func ParseInt(b []byte) (int64, error) {
	if UseStdLib {
		return strconv.ParseInt(string(b), 10, 64)
	}

	i := 0
	neg := false
	if len(b) > 0 && (b[0] == '+' || b[0] == '-') {
		neg = b[0] == '-'
		i++
	}
	n := uint64(0)
	for i < len(b) {
		c := b[i]
		if n > math.MaxUint64/10 {
			return 0, rangeError("ParseInt", string(b))
		} else if c >= '0' && c <= '9' {
			n *= 10
			n += uint64(c - '0')
		} else {
			break
		}
		i++
	}
	if !neg && n > uint64(math.MaxInt64) || n > uint64(math.MaxInt64)+1 {
		return 0, rangeError("ParseInt", string(b))
	} else if neg {
		return -int64(n), nil
	}
	if len(b) != i {
		return 0, syntaxError("ParseInt", string(b))
	}
	return int64(n), nil
}

func LenInt(i int64) int {
	if i < 0 {
		i = -i
	}
	switch {
	case i < 10:
		return 1
	case i < 100:
		return 2
	case i < 1000:
		return 3
	case i < 10000:
		return 4
	case i < 100000:
		return 5
	case i < 1000000:
		return 6
	case i < 10000000:
		return 7
	case i < 100000000:
		return 8
	case i < 1000000000:
		return 9
	case i < 10000000000:
		return 10
	case i < 100000000000:
		return 11
	case i < 1000000000000:
		return 12
	case i < 10000000000000:
		return 13
	case i < 100000000000000:
		return 14
	case i < 1000000000000000:
		return 15
	case i < 10000000000000000:
		return 16
	case i < 100000000000000000:
		return 17
	case i < 1000000000000000000:
		return 18
	}
	return 19
}
