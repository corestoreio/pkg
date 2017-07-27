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
func Repeat(sql string, args ...Argument) (string, []interface{}, error) {

	phCount := strings.Count(sql, placeHolderStr)
	if want := len(args); phCount != want || want == 0 {
		return "", nil, errors.NewMismatchf("[dbr] Repeat: Number of %s:%d do not match the number of repetitions: %d", placeHolderStr, phCount, want)
	}

	retArgs := make([]interface{}, 0, len(args)*2)

	buf := bufferpool.Get()
	defer bufferpool.Put(buf)

	n := phCount
	i := 0
	for i < n {
		m := strings.IndexByte(sql, placeHolderRune)
		if m < 0 {
			break
		}
		buf.WriteString(sql[:m])

		if i < len(args) {
			prevLen := len(retArgs)
			retArgs = args[i].toIFace(retArgs)
			reps := len(retArgs) - prevLen
			for r := 0; r < reps; r++ {
				buf.WriteByte(placeHolderRune)
				if r < reps-1 {
					buf.WriteByte(',')
				}
			}
		}
		sql = sql[m+len(placeHolderStr):]
		i++
	}
	buf.WriteString(sql)
	return buf.String(), retArgs, nil
}

// repeatPlaceHolders multiplies the place holder with the arguments internal len.
func repeatPlaceHolders(buf queryWriter, sql []byte, args ...Argument) error {

	i := 0
	pos := 0
	for pos < len(sql) {
		r, w := utf8.DecodeRune(sql[pos:])
		pos += w

		switch r {
		case '?':
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

type ipolate struct {
	useCache   bool
	queryCache []byte
	args       Arguments
}

// Interpolate takes a SQL byte slice with placeholders and a list of arguments
// to replace them with. It returns a blank string or an error if the number of
// placeholders does not match the number of arguments. Implements the Repeat
// function.
func Interpolate(sql string) *ipolate {
	return &ipolate{
		queryCache: []byte(sql),
		args:       make(Arguments, 0, 8),
	}
}

// String implements fmt.Stringer and prints errors into the string which will
// maybe generate invalid SQL code.
func (ip *ipolate) String() string {
	sql, _, err := ip.ToSQL()
	if err != nil {
		return err.Error()
	}
	return sql
}

// MustString returns the fully interpolated SQL string or panics on error.
func (ip *ipolate) MustString() string {
	sql, _, err := ip.ToSQL()
	if err != nil {
		panic(err)
	}
	return sql
}

// ToSQL implements dbr.QueryBuilder. The interface slice is always nil.
func (ip *ipolate) ToSQL() (_ string, alwaysNil []interface{}, _ error) {
	buf := bufferpool.Get()
	defer bufferpool.Put(buf)
	if err := writeInterpolate(buf, ip.queryCache, ip.args); err != nil {
		return "", nil, errors.WithStack(err)
	}
	return buf.String(), nil, nil
}

// Reset resets the internal argument cache for reuse. Avoids lots of
// allocations.
func (ip *ipolate) Reset() *ipolate                { ip.args = ip.args[:0]; return ip }
func (ip *ipolate) Null() *ipolate                 { ip.args = append(ip.args, nil); return ip }
func (ip *ipolate) Int(i int) *ipolate             { ip.args = append(ip.args, Int(i)); return ip }
func (ip *ipolate) Ints(i ...int) *ipolate         { ip.args = append(ip.args, Ints(i)); return ip }
func (ip *ipolate) Int64(i int64) *ipolate         { ip.args = append(ip.args, Int64(i)); return ip }
func (ip *ipolate) Int64s(i ...int64) *ipolate     { ip.args = append(ip.args, Int64s(i)); return ip }
func (ip *ipolate) Uint64(i uint64) *ipolate       { ip.args = append(ip.args, Uint64(i)); return ip }
func (ip *ipolate) Float64(f float64) *ipolate     { ip.args = append(ip.args, Float64(f)); return ip }
func (ip *ipolate) Float64s(f ...float64) *ipolate { ip.args = append(ip.args, Float64s(f)); return ip }
func (ip *ipolate) Str(s string) *ipolate          { ip.args = append(ip.args, String(s)); return ip }
func (ip *ipolate) Strs(s ...string) *ipolate      { ip.args = append(ip.args, Strings(s)); return ip }
func (ip *ipolate) Bool(b bool) *ipolate           { ip.args = append(ip.args, Bool(b)); return ip }
func (ip *ipolate) Bools(b ...bool) *ipolate       { ip.args = append(ip.args, Bools(b)); return ip }
func (ip *ipolate) Bytes(p []byte) *ipolate        { ip.args = append(ip.args, Bytes(p)); return ip }
func (ip *ipolate) BytesSlice(p ...[]byte) *ipolate {
	ip.args = append(ip.args, BytesSlice(p))
	return ip
}
func (ip *ipolate) Time(t time.Time) *ipolate     { ip.args = append(ip.args, MakeTime(t)); return ip }
func (ip *ipolate) Times(t ...time.Time) *ipolate { ip.args = append(ip.args, Times(t)); return ip }

func (ip *ipolate) NullString(nv ...NullString) *ipolate {
	if len(nv) == 1 {
		ip.args = append(ip.args, nv[0])
	} else {
		ip.args = append(ip.args, NullStrings(nv))
	}
	return ip
}

func (ip *ipolate) NullFloat64(nv ...NullFloat64) *ipolate {
	if len(nv) == 1 {
		ip.args = append(ip.args, nv[0])
	} else {
		ip.args = append(ip.args, NullFloat64s(nv))
	}
	return ip
}

func (ip *ipolate) NullInt64(nv ...NullInt64) *ipolate {
	if len(nv) == 1 {
		ip.args = append(ip.args, nv[0])
	} else {
		ip.args = append(ip.args, NullInt64s(nv))
	}
	return ip
}

func (ip *ipolate) NullBool(nv NullBool) *ipolate {
	ip.args = append(ip.args, nv)
	return ip
}
func (ip *ipolate) NullTime(nv ...NullTime) *ipolate {
	if len(nv) == 1 {
		ip.args = append(ip.args, nv[0])
	} else {
		ip.args = append(ip.args, NullTimes(nv))
	}
	return ip
}
func (ip *ipolate) Value(dv ...driver.Valuer) *ipolate {
	ip.args = append(ip.args, DriverValues(dv))
	return ip
}
func (ip *ipolate) Arguments(arg ...Argument) *ipolate {
	// todo maybe make this method Arguments private
	ip.args = append(ip.args, arg...)
	return ip
}

// Named uses the NamedArg for string replacement. Replaces the names with place
// holder character. TODO(CyS) Slices in NamedArg.Value are not yet supported.
func (ip *ipolate) Named(nArgs ...sql.NamedArg) *ipolate {
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
func writeInterpolate(buf queryWriter, sql []byte, args Arguments) error {

	// TODO(CyS) due to the type `interpolate`, we can optimize the parsing in
	// the second run with the same SQL slice but different arguments. We know
	// ahead on which position the insertion must happen. Some refactoring needs
	// to be done.

	phCount := bytes.Count(sql, bTextPlaceholder)
	argCount := Arguments(args).len()

	// Repeats the place holders, e.g. IN (?) will become IN (?,?,?)
	if phCount < argCount {
		rBuf := bufferpool.Get()
		defer bufferpool.Put(rBuf)
		if err := repeatPlaceHolders(rBuf, sql, args...); err != nil {
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
		if args[0] != nil {
			argLength = args[0].len()
		}
	}
	pos := 0
	for pos < len(sql) {
		r, w := utf8.DecodeRune(sql[pos:])
		pos += w

		switch {
		case r == '?':
			if phCounter < argLength-1 {
				phCounter++
			} else {
				phCounter = 0 // next argument set starts
				argIndex++
				if argIndex >= len(args) {
					return errors.NewNotValidf("[dbr] Arguments are imbalanced. Argument Index %d but argument count was %d", argIndex, len(args)-1)
				}
				argLength = 1
				if args[argIndex] != nil {
					argLength = args[argIndex].len()
				}
			}
			if args[argIndex] == nil {
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
		return errors.NewNotValidf("[dbr] Arguments are imbalanced. Placeholders: %d Current argument count: %d or %d", phTotals, argCount, len(args))
	}
	return nil
}
