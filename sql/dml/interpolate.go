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

package dml

import (
	"bytes"
	"database/sql"
	"database/sql/driver"
	"strings"
	"time"
	"unicode/utf8"

	"github.com/corestoreio/cspkg/util/bufferpool"
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
func Repeat(sql string, args Arguments) (string, error) {

	phCount := strings.Count(sql, placeHolderStr)
	if want := len(args); phCount != want || want == 0 {
		return "", errors.NewMismatchf("[dml] Repeat: Number of %s:%d do not match the number of repetitions: %d", placeHolderStr, phCount, want)
	}

	buf := bufferpool.Get()
	defer bufferpool.Put(buf)
	err := repeatPlaceHolders(buf, []byte(sql), args)
	return buf.String(), errors.WithStack(err)
}

// repeatPlaceHolders multiplies the place holder with the arguments internal len.
func repeatPlaceHolders(buf *bytes.Buffer, sql []byte, args Arguments) error {
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

// ip handles the interpolation of the SQL string and uses an internal argument
// pool for optimal slice usage.
type ip struct {
	queryCache []byte
	args       Arguments
}

// Interpolate takes a SQL byte slice with placeholders and a list of arguments
// to replace them with. It returns a blank string or an error if the number of
// placeholders does not match the number of arguments. Implements the Repeat
// function.
func Interpolate(sql string) *ip {
	return &ip{
		queryCache: []byte(sql),
		args:       MakeArgs(8),
	}
}

// String implements fmt.Stringer and prints errors into the string which will
// maybe generate invalid SQL code.
func (in *ip) String() string {
	str, _, err := in.ToSQL()
	if err != nil {
		return err.Error()
	}
	return str
}

// MustString returns the fully interpolated SQL string or panics on error.
func (in *ip) MustString() string {
	str, _, err := in.ToSQL()
	if err != nil {
		panic(err)
	}
	return str
}

// ToSQL implements dml.QueryBuilder. The interface slice is always nil.
func (in *ip) ToSQL() (_ string, alwaysNil []interface{}, _ error) {
	buf := bufferpool.Get()
	defer bufferpool.Put(buf)
	if err := writeInterpolate(buf, in.queryCache, in.args); err != nil {
		return "", nil, errors.WithStack(err)
	}
	return buf.String(), nil, nil
}

// Reset resets the internal argument cache for reuse. Avoids lots of
// allocations.
func (in *ip) Reset() *ip              { in.args = in.args[:0]; return in }
func (in *ip) Null() *ip               { in.args = in.args.Null(); return in }
func (in *ip) Int(i int) *ip           { in.args = in.args.Int(i); return in }
func (in *ip) Ints(i ...int) *ip       { in.args = in.args.Ints(i...); return in }
func (in *ip) Int64(i int64) *ip       { in.args = in.args.Int64(i); return in }
func (in *ip) Int64s(i ...int64) *ip   { in.args = in.args.Int64s(i...); return in }
func (in *ip) Uint64(i uint64) *ip     { in.args = in.args.Uint64(i); return in }
func (in *ip) Uint64s(i ...uint64) *ip { in.args = in.args.Uint64s(i...); return in }
func (in *ip) Float64(f float64) *ip   { in.args = in.args.Float64(f); return in }
func (in *ip) Float64s(f ...float64) *ip {
	in.args = in.args.Float64s(f...)
	return in
}
func (in *ip) Str(s string) *ip     { in.args = in.args.String(s); return in }
func (in *ip) Strs(s ...string) *ip { in.args = in.args.Strings(s...); return in }
func (in *ip) Bool(b bool) *ip      { in.args = in.args.Bool(b); return in }
func (in *ip) Bools(b ...bool) *ip  { in.args = in.args.Bools(b...); return in }

// Bytes uses a byte slice for comparison. Providing a nil value returns a
// NULL type. Detects between valid UTF-8 strings and binary data. Later gets
// hex encoded.
func (in *ip) Bytes(p []byte) *ip { in.args = in.args.Bytes(p); return in }
func (in *ip) BytesSlice(p ...[]byte) *ip {
	in.args = in.args.BytesSlice(p...)
	return in
}
func (in *ip) Time(t time.Time) *ip     { in.args = in.args.Time(t); return in }
func (in *ip) Times(t ...time.Time) *ip { in.args = in.args.Times(t...); return in }
func (in *ip) NullString(nv ...NullString) *ip {
	in.args = in.args.NullString(nv...)
	return in
}
func (in *ip) NullFloat64(nv ...NullFloat64) *ip {
	in.args = in.args.NullFloat64(nv...)
	return in
}
func (in *ip) NullInt64(nv ...NullInt64) *ip {
	in.args = in.args.NullInt64(nv...)
	return in
}
func (in *ip) NullBool(nv ...NullBool) *ip {
	in.args = in.args.NullBool(nv...)
	return in
}
func (in *ip) NullTime(nv ...NullTime) *ip {
	in.args = in.args.NullTime(nv...)
	return in
}
func (in *ip) DriverValue(dvs ...driver.Valuer) *ip {
	in.args = in.args.DriverValue(dvs...)
	return in
}
func (in *ip) ArgUnions(args Arguments) *ip {
	in.args = args
	return in
}

// Named uses the NamedArg for string replacement. Replaces the names with place
// holder character. TODO(CyS) Slices in NamedArg.Value are not yet supported.
func (in *ip) Named(nArgs ...sql.NamedArg) *ip {
	// for now this unoptimized version with a stupid string replacement and
	// converting between bytes and string.
	sqlStr := string(in.queryCache)
	for _, na := range nArgs {
		sqlStr = strings.Replace(sqlStr, na.Name, "?", -1)
		in.args = append(in.args, iFaceToArgs(na.Value)...)
	}
	in.queryCache = []byte(sqlStr)
	return in
}

var bTextPlaceholder = []byte("?")

// writeInterpolate merges `args` into `sql` and writes the result into `buf`. `sql`
// stays unchanged.
func writeInterpolate(buf *bytes.Buffer, sql []byte, args Arguments) error {

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
		if args[0].isSet {
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
					return errors.NewNotValidf("[dml] Interpolate: Arguments are imbalanced. Argument Index %d is greater than argument count %d", argIndex, len(args)-1)
				}
				argLength = 1
				if args[argIndex].isSet {
					argLength = args[argIndex].len()
				}
			}
			if !args[argIndex].isSet {
				buf.WriteString(sqlStrNullUC)
			} else if err := args[argIndex].writeTo(buf, phCounter); err != nil {
				return errors.WithStack(err)
			}

			phTotals++
		case r == '`', r == '\'', r == '"':
			p := bytes.IndexRune(sql[pos:], r)
			if p == -1 {
				return errors.NewNotValidf("[dml] Interpolate: Invalid syntax")
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
		return errors.NewNotValidf("[dml] Interpolate: Arguments are imbalanced. Placeholder count: %d. Current argument count: %d or %d", phTotals, argCount, len(args))
	}
	return nil
}
