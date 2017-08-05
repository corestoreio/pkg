// Copyright 2015-2017, Cyrill @ Schumacher.fm and the CoreStore contributors
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

package dbr

import (
	"bytes"
	"database/sql"
	"database/sql/driver"
	"strings"
	"time"
	"unicode/utf8"

	"github.com/corestoreio/csfw/util/bufferpool"
	"github.com/corestoreio/errors"
)

const (
	placeHolderRune = '?'
	placeHolderStr  = `?`
)

// Repeat takes a SQL string and repeats the question marks with the provided
// arguments. If the amount of arguments does not match the number of questions
// marks, a Mismatch error gets returned. The arguments are getting converted to
// an interface slice to easy passing into the db.Query/db.Exec/etc functions at
// an argument.
//		Repeat("SELECT * FROM table WHERE id IN (?) AND status IN (?)", Int(myIntSlice...), String(myStrSlice...))
// Gets converted to:
//		SELECT * FROM table WHERE id IN (?,?) AND status IN (?,?,?)
// The questions marks are of course depending on the values in the Arg*
// functions. This function should be generally used when dealing with prepared
// statements.
func Repeat(sql string, args ArgUnions) (string, error) {

	phCount := strings.Count(sql, placeHolderStr)
	if want := len(args); phCount != want || want == 0 {
		return "", errors.NewMismatchf("[dbr] Repeat: Number of %s:%d do not match the number of repetitions: %d", placeHolderStr, phCount, want)
	}

	buf := bufferpool.Get()
	defer bufferpool.Put(buf)
	err := repeatPlaceHolders(buf, []byte(sql), args)
	return buf.String(), errors.WithStack(err)
}

// repeatPlaceHolders multiplies the place holder with the arguments internal len.
func repeatPlaceHolders(buf *bytes.Buffer, sql []byte, args ArgUnions) error {
	i := 0
	pos := 0
	for pos < len(sql) {
		r, w := utf8.DecodeRune(sql[pos:])
		pos += w

		switch r {
		case placeHolderRune:
			if i < len(args) {
				reps := args[i].len()
				for r := 0; r < reps; r++ {
					buf.WriteByte(placeHolderRune)
					if r < reps-1 {
						buf.WriteByte(',')
					}
				}
			}
			i++
		default:
			buf.WriteRune(r)
		}
	}
	return nil
}

type iPolate struct {
	queryCache []byte
	args       ArgUnions
}

// Interpolate takes a SQL byte slice with placeholders and a list of arguments
// to replace them with. It returns a blank string or an error if the number of
// placeholders does not match the number of arguments. Implements the Repeat
// function.
func Interpolate(sql string) *iPolate {
	return &iPolate{
		queryCache: []byte(sql),
		args:       MakeArgUnions(8),
	}
}

// String implements fmt.Stringer and prints errors into the string which will
// maybe generate invalid SQL code.
func (ip *iPolate) String() string {
	str, _, err := ip.ToSQL()
	if err != nil {
		return err.Error()
	}
	return str
}

// MustString returns the fully interpolated SQL string or panics on error.
func (ip *iPolate) MustString() string {
	str, _, err := ip.ToSQL()
	if err != nil {
		panic(err)
	}
	return str
}

// ToSQL implements dbr.QueryBuilder. The interface slice is always nil.
func (ip *iPolate) ToSQL() (_ string, alwaysNil []interface{}, _ error) {
	buf := bufferpool.Get()
	defer bufferpool.Put(buf)
	if err := writeInterpolate(buf, ip.queryCache, ip.args); err != nil {
		return "", nil, errors.WithStack(err)
	}
	return buf.String(), nil, nil
}

// Reset resets the internal argument cache for reuse. Avoids lots of
// allocations.
func (ip *iPolate) Reset() *iPolate              { ip.args = ip.args[:0]; return ip }
func (ip *iPolate) Null() *iPolate               { ip.args = ip.args.Null(); return ip }
func (ip *iPolate) Int(i int) *iPolate           { ip.args = ip.args.Int(i); return ip }
func (ip *iPolate) Ints(i ...int) *iPolate       { ip.args = ip.args.Ints(i...); return ip }
func (ip *iPolate) Int64(i int64) *iPolate       { ip.args = ip.args.Int64(i); return ip }
func (ip *iPolate) Int64s(i ...int64) *iPolate   { ip.args = ip.args.Int64s(i...); return ip }
func (ip *iPolate) Uint64(i uint64) *iPolate     { ip.args = ip.args.Uint64(i); return ip }
func (ip *iPolate) Uint64s(i ...uint64) *iPolate { ip.args = ip.args.Uint64s(i...); return ip }
func (ip *iPolate) Float64(f float64) *iPolate   { ip.args = ip.args.Float64(f); return ip }
func (ip *iPolate) Float64s(f ...float64) *iPolate {
	ip.args = ip.args.Float64s(f...)
	return ip
}
func (ip *iPolate) Str(s string) *iPolate     { ip.args = ip.args.Str(s); return ip }
func (ip *iPolate) Strs(s ...string) *iPolate { ip.args = ip.args.Strs(s...); return ip }
func (ip *iPolate) Bool(b bool) *iPolate      { ip.args = ip.args.Bool(b); return ip }
func (ip *iPolate) Bools(b ...bool) *iPolate  { ip.args = ip.args.Bools(b...); return ip }

// Bytes uses a byte slice for comparison. Providing a nil value returns a
// NULL type. Detects between valid UTF-8 strings and binary data. Later gets
// hex encoded.
func (ip *iPolate) Bytes(p []byte) *iPolate { ip.args = ip.args.Bytes(p); return ip }
func (ip *iPolate) BytesSlice(p ...[]byte) *iPolate {
	ip.args = ip.args.BytesSlice(p...)
	return ip
}
func (ip *iPolate) Time(t time.Time) *iPolate     { ip.args = ip.args.Time(t); return ip }
func (ip *iPolate) Times(t ...time.Time) *iPolate { ip.args = ip.args.Times(t...); return ip }
func (ip *iPolate) NullString(nv ...NullString) *iPolate {
	ip.args = ip.args.NullString(nv...)
	return ip
}
func (ip *iPolate) NullFloat64(nv ...NullFloat64) *iPolate {
	ip.args = ip.args.NullFloat64(nv...)
	return ip
}
func (ip *iPolate) NullInt64(nv ...NullInt64) *iPolate {
	ip.args = ip.args.NullInt64(nv...)
	return ip
}
func (ip *iPolate) NullBool(nv ...NullBool) *iPolate {
	ip.args = ip.args.NullBool(nv...)
	return ip
}
func (ip *iPolate) NullTime(nv ...NullTime) *iPolate {
	ip.args = ip.args.NullTime(nv...)
	return ip
}
func (ip *iPolate) DriverValue(dvs ...driver.Valuer) *iPolate {
	ip.args = ip.args.DriverValue(dvs...)
	return ip
}
func (ip *iPolate) ArgUnions(args ArgUnions) *iPolate {
	ip.args = args
	return ip
}

// Named uses the NamedArg for string replacement. Replaces the names with place
// holder character. TODO(CyS) Slices in NamedArg.Value are not yet supported.
func (ip *iPolate) Named(nArgs ...sql.NamedArg) *iPolate {
	// for now this unoptimized version with a stupid string replacement and
	// converting between bytes and string.
	sqlStr := string(ip.queryCache)
	for _, na := range nArgs {
		sqlStr = strings.Replace(sqlStr, na.Name, "?", -1)
		ip.args = append(ip.args, iFaceToArgs(na.Value)...)
	}
	ip.queryCache = []byte(sqlStr)
	return ip
}

var bTextPlaceholder = []byte("?")

// writeInterpolate merges `args` into `sql` and writes the result into `buf`. `sql`
// stays unchanged.
func writeInterpolate(buf *bytes.Buffer, sql []byte, args ArgUnions) error {

	// TODO(CyS) due to the type `interpolate`, we can optimize the parsing in
	// the second run with the same SQL slice but different arguments. We know
	// ahead on which position the insertion must happen. Some refactoring needs
	// to be done.

	phCount := bytes.Count(sql, bTextPlaceholder)
	argCount := args.Len()

	// Repeats the place holders, e.g. IN (?) will become IN (?,?,?)
	if phCount < argCount {
		rBuf := bufferpool.Get()
		defer bufferpool.Put(rBuf)
		if err := repeatPlaceHolders(rBuf, sql, args); err != nil {
			return errors.WithStack(err)
		}
		sql = rBuf.Bytes()
	}

	phTotals := 0
	phCounter := -1
	argIndex := 0
	argLength := 0
	if len(args) > 0 {
		argLength = 1
		if args[0].field > 0 {
			argLength = args[0].len()
		}
	}
	pos := 0
	for pos < len(sql) {
		r, w := utf8.DecodeRune(sql[pos:])
		pos += w

		switch {
		case r == placeHolderRune:
			if phCounter < argLength-1 {
				phCounter++
			} else {
				phCounter = 0 // next argument set starts
				argIndex++
				if argIndex >= len(args) {
					return errors.NewNotValidf("[dbr] Arguments are imbalanced. Argument Index %d is greater than argument count %d", argIndex, len(args)-1)
				}
				argLength = 1
				if args[argIndex].field > 0 {
					argLength = args[argIndex].len()
				}
			}
			if args[argIndex].field == 0 {
				buf.WriteString("NULL")
			} else if err := args[argIndex].writeTo(buf, phCounter); err != nil {
				return errors.WithStack(err)
			}

			phTotals++
		case r == '`', r == '\'', r == '"':
			p := bytes.IndexRune(sql[pos:], r)
			if p == -1 {
				return errors.NewNotValidf("[dbr] Interpolate: Invalid syntax")
			}
			if r == '"' {
				r = '\''
			}
			buf.WriteRune(r)
			buf.Write(sql[pos : pos+p])
			buf.WriteRune(r)
			pos += p + 1
		case r == '[':
			w := bytes.IndexRune(sql[pos:], ']')
			col := sql[pos : pos+w]
			dialect.EscapeIdent(buf, string(col))
			pos += w + 1 // size of ']'
		default:
			buf.WriteRune(r)
		}
	}

	if phTotals != argCount {
		return errors.NewNotValidf("[dbr] args are imbalanced. Placeholders: %d Current argument count: %d or %d", phTotals, argCount, len(args))
	}
	return nil
}
