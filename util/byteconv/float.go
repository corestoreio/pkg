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

// This package contains string conversion function and is written in [Go][1].
// It is much alike the standard library's strconv package, but it is
// specifically tailored for the performance needs within the minify package.
// For example, the floating-point to string conversion function is
// approximately twice as fast as the standard library, but it is not as
// precise.

package byteconv

import (
	"database/sql"
	"math"
	"strconv"
)

var float64pow10 = []float64{
	1e0, 1e1, 1e2, 1e3, 1e4, 1e5, 1e6, 1e7, 1e8, 1e9,
	1e10, 1e11, 1e12, 1e13, 1e14, 1e15, 1e16, 1e17, 1e18, 1e19,
	1e20, 1e21, 1e22,
}

// ParseNullFloat64 same as ParseFloat
func ParseNullFloat64(b []byte) (val sql.NullFloat64, err error) {
	if len(b) == 0 {
		return
	}
	val.Float64, err = ParseFloat(b)
	val.Valid = err == nil
	return
}

// Float parses a byte-slice and returns the float it represents.
// If an invalid character is encountered, it will stop there.
func ParseFloat(b []byte) (float64, error) {
	if UseStdLib {
		return strconv.ParseFloat(string(b), 64)
	}

	i := 0
	neg := false
	if i < len(b) && (b[i] == '+' || b[i] == '-') {
		neg = b[i] == '-'
		i++
	}

	dot := -1
	trunk := -1
	n := uint64(0)
	for ; i < len(b); i++ {
		c := b[i]
		if c >= '0' && c <= '9' {
			if trunk == -1 {
				if n > math.MaxUint64/10 {
					trunk = i
				} else {
					n *= 10
					n += uint64(c - '0')
				}
			}
		} else if dot == -1 && c == '.' {
			dot = i
		} else {
			if c == 'e' || c == 'E' {
				break
			}
			return 0, syntaxError("ParseFloat", string(b))
		}
	}

	f := float64(n)
	if neg {
		f = -f
	}

	mantExp := int64(0)
	if dot != -1 {
		if trunk == -1 {
			trunk = i
		}
		mantExp = int64(trunk - dot - 1)
	} else if trunk != -1 {
		mantExp = int64(trunk - i)
	}
	expExp := int64(0)
	if i < len(b) && (b[i] == 'e' || b[i] == 'E') {
		i++
		if e, err := ParseInt(b[i:]); err == nil {
			expExp = e
			i += LenInt(e)
		} else {
			return 0, syntaxError("ParseFloat", string(b))
		}
	}
	exp := expExp - mantExp

	// copied from strconv/atof.go
	if exp == 0 {
		return f, nil
	} else if exp > 0 && exp <= 15+22 { // int * 10^k
		// If exponent is big but number of digits is not,
		// can move a few zeros into the integer part.
		if exp > 22 {
			f *= float64pow10[exp-22]
			exp = 22
		}
		if f <= 1e15 && f >= -1e15 {
			return f * float64pow10[exp], nil
		}
	} else if exp < 0 && exp >= -22 { // int / 10^k
		return f / float64pow10[-exp], nil
	}
	f *= math.Pow10(int(-mantExp))
	return f * math.Pow10(int(expExp)), nil
}
