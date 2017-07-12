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
	"github.com/corestoreio/errors"
	"time"
	"unicode/utf8"
)

// ValueAssembler assembles arguments for CRUD statements. The `stmtType`
// variable contains a bit flag from the constants SQLStmt* and SQLPart* to
// allow the knowledge in which case the function AssembleArguments gets called.
// Any new arguments must be append to variable `args` and then returned.
// Variable `columns` contains the name of the requested columns. E.g. if the
// first requested column names `id` then the first appended argument must be an
// integer. Variable `columns` can additionally contain the names and/or
// expressions used in the WHERE, JOIN or HAVING clauses, if applicable for the
// SQL statement type. In case where stmtType has been set to SQLStmtInsert|SQLPartValues, the
// `columns` slice can be empty which means that all arguments are requested.
type ValueAssembler interface {
	AssembleValues(stmtType int, vals Values, columns []string) (Values, error)
}

type Value struct {
	int
	int64
	float64
	string // can contain also the SQL expression
	bool
	byte []byte
	time time.Time

	hasInt        bool
	hasInt64      bool
	hasFloat64    bool
	hasString     bool
	hasBool       bool
	hasTime       bool
	isPlaceHolder bool // aka cahens constant
	isExpression  bool

	ints     []int
	int64s   []int64
	float64s []float64
	strings  []string
	times    []time.Time
	bytes    [][]byte
	vals     Values // in case of expression
}

func newValue() *Value { return new(Value) }

func (v *Value) reset() *Value {
	v.int = 0
	v.int64 = 0
	v.float64 = 0
	v.string = ""
	v.bool = false
	v.byte = nil
	v.time = time.Time{}

	v.hasInt = false
	v.hasInt64 = false
	v.hasFloat64 = false
	v.hasString = false
	v.hasBool = false
	v.hasTime = false

	v.isPlaceHolder = false
	v.isExpression = false

	v.ints = nil
	v.int64s = nil
	v.float64s = nil
	v.strings = nil
	v.times = nil
	v.bytes = nil
	v.vals = nil
	return v
}

func (v *Value) setInt(i int) *Value {
	v.int = i
	v.hasInt = true
	return v
}
func (v *Value) setInts(i ...int) *Value {
	v.isPlaceHolder = len(i) == 0
	v.ints = i
	return v
}
func (v *Value) setInt64(i int64) *Value {
	v.int64 = i
	v.hasInt64 = true
	return v
}
func (v *Value) setInt64s(i ...int64) *Value {
	v.isPlaceHolder = len(i) == 0
	v.int64s = i
	return v
}
func (v *Value) setFloat64(i float64) *Value {
	v.float64 = i
	v.hasFloat64 = true
	return v
}
func (v *Value) setFloat64s(i ...float64) *Value {
	v.isPlaceHolder = len(i) == 0
	v.float64s = i
	return v
}

func (v *Value) setString(s string) *Value {
	v.string = s
	v.hasString = true
	return v
}
func (v *Value) setStrings(s ...string) *Value {
	v.isPlaceHolder = len(s) == 0
	v.strings = s
	return v
}
func (v *Value) setBool(b bool) *Value {
	v.bool = b
	v.hasBool = true
	return v
}
func (v *Value) setTime(t time.Time) *Value {
	v.time = t
	v.hasTime = true
	return v
}
func (v *Value) setTimes(t ...time.Time) *Value {
	v.isPlaceHolder = len(t) == 0
	v.times = t
	return v
}

func (v *Value) toIFace(args []interface{}) []interface{} {

	switch {
	case v.hasInt:
		args = append(args, v.int)
	case v.hasInt64:
		args = append(args, v.int64)
	case v.hasFloat64:
		args = append(args, v.float64)
	case v.hasString:
		args = append(args, v.string)
	case v.hasBool:
		args = append(args, v.bool)
	case v.byte != nil:
		args = append(args, v.byte)
	case v.hasTime:
		args = append(args, v.time)
	case v.isPlaceHolder:
		// do nothing

	case v.ints != nil:
		for _, i := range v.ints {
			args = append(args, i)
		}
	case v.int64s != nil:
		for _, i := range v.int64s {
			args = append(args, i)
		}
	case v.float64s != nil:
		for _, i := range v.float64s {
			args = append(args, i)
		}
	case v.strings != nil:
		for _, i := range v.strings {
			args = append(args, i)
		}
	case v.times != nil:
		for _, i := range v.times {
			args = append(args, i)
		}
	case v.bytes != nil:
		for _, i := range v.bytes {
			args = append(args, i)
		}
	case v.isExpression:
		for _, val := range v.vals {
			args = val.toIFace(args) // recursion!
		}
	}

	return args
}

func (v *Value) writeTo(w queryWriter, pos int) (err error) {

	switch {
	case v.hasInt:
		err = writeInt64(w, int64(v.int))
	case v.hasInt64:
		err = writeInt64(w, v.int64)
	case v.hasFloat64:
		err = writeFloat64(w, v.float64)
	case v.hasString:
		if !utf8.ValidString(v.string) {
			return errors.NewNotValidf("[dbr] Value.WriteTo: String is not UTF-8: %q", v.string)
		}
		dialect.EscapeString(w, v.string)
	case v.hasBool:
		dialect.EscapeBool(w, v.bool)
	case v.byte != nil:
		if !utf8.Valid(v.byte) {
			dialect.EscapeBinary(w, v.byte)
		} else {
			dialect.EscapeString(w, string(v.byte))
		}
	case v.hasTime:
		dialect.EscapeTime(w, v.time)
	case v.isPlaceHolder:
		_, err = w.WriteString("? /*PLACEHOLDER*/") // maybe remove /*PLACEHOLDER*/ if it's annoying

	case v.ints != nil:
		err = writeInt64(w, int64(v.ints[pos]))
	case v.int64s != nil:
		err = writeInt64(w, v.int64s[pos])
	case v.float64s != nil:
		err = writeFloat64(w, v.float64s[pos])
	case v.strings != nil:
		if !utf8.ValidString(v.strings[pos]) {
			return errors.NewNotValidf("[dbr] Value.WriteTo: String is not UTF-8: %q", v.strings[pos])
		}
		dialect.EscapeString(w, v.strings[pos])
	case v.times != nil:
		dialect.EscapeTime(w, v.times[pos])
	case v.bytes != nil:
		if !utf8.Valid(v.bytes[pos]) {
			dialect.EscapeBinary(w, v.bytes[pos])
		} else {
			dialect.EscapeString(w, string(v.bytes[pos]))
		}
	case v.isExpression:
		_, err = w.WriteString(v.string)
	default:
		_, err = w.WriteString("NULL")
	}

	return err
}

func (v *Value) len() (l int) {

	switch {
	case v.hasInt:
		l = 1
	case v.hasInt64:
		l = 1
	case v.hasFloat64:
		l = 1
	case v.hasString:
		l = 1
	case v.hasBool:
		l = 1
	case v.byte != nil:
		l = 1
	case v.hasTime:
		l = 1
	case v.isPlaceHolder:
		// do nothing

	case v.ints != nil:
		l = len(v.ints)
	case v.int64s != nil:
		l = len(v.int64s)
	case v.float64s != nil:
		l = len(v.float64s)
	case v.strings != nil:
		l = len(v.strings)
	case v.times != nil:
		l = len(v.times)
	case v.bytes != nil:
		l = len(v.bytes)
	case v.isExpression:
		l = len(v.vals)
	}

	return l
}

type Values []*Value

// Interfaces converts the underlying concrete types into an interface slice.
// Each entry in the interface is guaranteed to be one of the following values:
// []byte, bool, float64, int64, string or time.Time. Use driver.IsValue() for a
// check.
func (as Values) Interfaces() []interface{} {
	if len(as) == 0 {
		return nil
	}
	ret := make([]interface{}, 0, len(as))
	for _, a := range as {
		ret = a.toIFace(ret)
	}
	return ret
}

// len calculates the total length of all values
func (as Values) len() (tl int) {
	for _, a := range as {
		l := a.len()
		if l == cahensConstant {
			l = 1
		}
		tl += l
	}
	return
}
