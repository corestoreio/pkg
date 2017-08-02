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
	"database/sql/driver"
	"time"
	"unicode/utf8"

	"github.com/corestoreio/errors"
)

const (
	argFieldNull uint8 = iota + 1
	argFieldInt64
	argFieldInt64s
	argFieldUint64
	argFieldUint64s
	argFieldFloat64
	argFieldFloat64s
	argFieldBool
	argFieldBools
	argFieldString
	argFieldStrings
	argFieldByte
	argFieldBytes
	argFieldTime
	argFieldTimes
	argFieldNullStrings
	argFieldNullInt64s
	argFieldNullFloat64s
	argFieldNullBool
	argFieldNullTimes
)

// argUnion is union type for different Go primitives and their slice
// representation. argUnion must be used as a pointer because it slows
// everything down. Check the benchmarks.
type argUnion struct {
	field uint8
	bool
	int64
	uint64
	float64
	string
	int64s   []int64
	uint64s  []uint64
	float64s []float64
	bools    []bool
	strings  []string
	bytes    []byte
	bytess   [][]byte
	times    []time.Time
	time     time.Time

	nullStrings  NullStrings
	nullInt64s   NullInt64s
	nullFloat64s NullFloat64s
	nullBool     NullBool
	nullTimes    NullTimes
}

func (arg argUnion) len() (l int) {
	switch arg.field {
	case argFieldNull, argFieldInt64, argFieldUint64, argFieldFloat64, argFieldBool, argFieldString, argFieldByte, argFieldTime, argFieldNullBool:
		l = 1
	case argFieldInt64s:
		l = len(arg.int64s)
	case argFieldUint64s:
		l = len(arg.uint64s)
	case argFieldFloat64s:
		l = len(arg.float64s)
	case argFieldBools:
		l = len(arg.bools)
	case argFieldStrings:
		l = len(arg.strings)
	case argFieldBytes:
		l = len(arg.bytess)
	case argFieldTimes:
		l = len(arg.times)
	case argFieldNullStrings:
		l = len(arg.nullStrings)
	case argFieldNullInt64s:
		l = len(arg.nullInt64s)
	case argFieldNullFloat64s:
		l = len(arg.nullFloat64s)
	case argFieldNullTimes:
		l = len(arg.nullTimes)
	}
	// default is 0
	return
}

func (arg argUnion) writeTo(w *bytes.Buffer, pos int) (err error) {
	switch arg.field {
	case argFieldInt64:
		err = writeInt64(w, arg.int64)
	case argFieldInt64s:
		err = writeInt64(w, arg.int64s[pos])
	case argFieldNullInt64s:
		if s := arg.nullInt64s[pos]; s.Valid {
			return writeInt64(w, s.Int64)
		}
		_, err = w.WriteString(sqlStrNull)

	case argFieldUint64:
		err = writeUint64(w, arg.uint64)
	case argFieldUint64s:
		err = writeUint64(w, arg.uint64s[pos])

	case argFieldFloat64:
		err = writeFloat64(w, arg.float64)
	case argFieldFloat64s:
		err = writeFloat64(w, arg.float64s[pos])
	case argFieldNullFloat64s:
		if s := arg.nullFloat64s[pos]; s.Valid {
			return writeFloat64(w, s.Float64)
		}
		_, err = w.WriteString(sqlStrNull)

	case argFieldBool:
		dialect.EscapeBool(w, arg.bool)
	case argFieldBools:
		dialect.EscapeBool(w, arg.bools[pos])
	case argFieldNullBool:
		if s := arg.nullBool; s.Valid {
			dialect.EscapeBool(w, s.Bool)
			return nil
		}
		_, err = w.WriteString(sqlStrNull)

		// TODO(CyS) Cut the printed string in errors if it's longer than XX chars
	case argFieldString:
		if !utf8.ValidString(arg.string) {
			return errors.NewNotValidf("[dbr] Argument.WriteTo: String is not UTF-8: %q", arg.string)
		}
		dialect.EscapeString(w, arg.string)
	case argFieldStrings:
		if !utf8.ValidString(arg.strings[pos]) {
			return errors.NewNotValidf("[dbr] Argument.WriteTo: String is not UTF-8: %q", arg.strings[pos])
		}
		dialect.EscapeString(w, arg.strings[pos])
	case argFieldNullStrings:
		if s := arg.nullStrings[pos]; s.Valid {
			if !utf8.ValidString(s.String) {
				return errors.NewNotValidf("[dbr] Argument.WriteTo: String is not UTF-8: %q", s.String)
			}
			dialect.EscapeString(w, s.String)
		} else {
			_, err = w.WriteString(sqlStrNull)
		}

	case argFieldByte:
		if !utf8.Valid(arg.bytes) {
			dialect.EscapeBinary(w, arg.bytes)
		} else {
			dialect.EscapeString(w, string(arg.bytes)) // maybe create an EscapeByteString version to avoid one alloc ;-)
		}

	case argFieldBytes:
		if !utf8.Valid(arg.bytess[pos]) {
			dialect.EscapeBinary(w, arg.bytess[pos])
		} else {
			dialect.EscapeString(w, string(arg.bytess[pos]))
		}

	case argFieldTime:
		dialect.EscapeTime(w, arg.time)
	case argFieldTimes:
		dialect.EscapeTime(w, arg.times[pos])
	case argFieldNullTimes:
		if nt := arg.nullTimes[pos]; nt.Valid {
			dialect.EscapeTime(w, nt.Time)
		} else {
			_, err = w.WriteString(sqlStrNull)
		}

	case argFieldNull:
		_, err = w.WriteString(sqlStrNull)

	default:
		panic(errors.NewNotSupportedf("[dbr] Unsupported field type: %d", arg.field))
	}
	return err
}

// ArgUninons a collection of primitive types or slice of primitive types. Using
// pointers in *argUnion would slow down the program.
type ArgUninons []argUnion

func makeArgUninons(cap int) ArgUninons {
	return make(ArgUninons, 0, cap)
}

// Len returns the total length of all arguments.
func (a ArgUninons) Len() int {
	var l int
	for _, arg := range a {
		l += arg.len()
	}
	return l
}

// Write writes all arguments into buf and separated by a colon.
func (a ArgUninons) Write(buf *bytes.Buffer) error {
	buf.WriteByte('(')
	for j, arg := range a {
		l := arg.len()
		for i := 0; i < l; i++ {
			if i > 0 || j > 0 {
				buf.WriteByte(',')
			}
			if err := arg.writeTo(buf, i); err != nil {
				return errors.Wrapf(err, "[dbr] ArgUninons write failed at pos %d with argument %#v", j, arg)
			}
		}
	}
	return buf.WriteByte(')')
}

// Interfaces creates an interface slice with flat values. Each type is one of
// the allowed in driver.Value.
func (a ArgUninons) Interfaces(args ...interface{}) []interface{} {
	if len(a) == 0 {
		return nil
	}
	if args == nil {
		args = make([]interface{}, 0, 2*len(a))
	}
	for _, arg := range a { // run bench between arg and a[i]
		switch arg.field {

		case argFieldInt64:
			args = append(args, arg.int64)
		case argFieldInt64s:
			for _, v := range arg.int64s {
				args = append(args, v)
			}
		case argFieldNullInt64s:
			for _, v := range arg.nullInt64s {
				if v.Valid {
					args = append(args, v.Int64)
				} else {
					args = append(args, nil)
				}
			}

			// TODO check if uint64 overflows int64
		case argFieldUint64:
			args = append(args, int64(arg.uint64))
		case argFieldUint64s:
			for _, v := range arg.uint64s {
				args = append(args, int64(v))
			}

		case argFieldFloat64:
			args = append(args, arg.float64)
		case argFieldFloat64s:
			for _, v := range arg.float64s {
				args = append(args, v)
			}
		case argFieldNullFloat64s:
			for _, v := range arg.nullFloat64s {
				if v.Valid {
					args = append(args, v.Float64)
				} else {
					args = append(args, nil)
				}
			}

		case argFieldBool:
			args = append(args, arg.bool)
		case argFieldBools:
			for _, v := range arg.bools {
				args = append(args, v)
			}
		case argFieldNullBool:
			if arg.nullBool.Valid {
				args = append(args, arg.nullBool.Bool)
			} else {
				args = append(args, nil)
			}

		case argFieldString:
			args = append(args, arg.string)
		case argFieldStrings:
			for _, v := range arg.strings {
				args = append(args, v)
			}
		case argFieldNullStrings:
			for _, v := range arg.nullStrings {
				if v.Valid {
					args = append(args, v.String)
				} else {
					args = append(args, nil)
				}
			}

		case argFieldByte:
			args = append(args, arg.bytes)
		case argFieldBytes:
			for _, v := range arg.bytess {
				args = append(args, v)
			}

		case argFieldTime:
			args = append(args, arg.time)
		case argFieldTimes:
			for _, v := range arg.times {
				args = append(args, v)
			}
		case argFieldNullTimes:
			for _, v := range arg.nullTimes {
				if v.Valid {
					args = append(args, v.Time)
				} else {
					args = append(args, nil)
				}
			}
		case argFieldNull:
			args = append(args, nil)
		}
	}
	return args
}

func (a ArgUninons) Null() ArgUninons {
	return append(a, argUnion{field: argFieldNull})
}
func (a ArgUninons) Int64(i int64) ArgUninons {
	return append(a, argUnion{field: argFieldInt64, int64: int64(i)})
}
func (a ArgUninons) Int64s(i ...int64) ArgUninons {
	return append(a, argUnion{field: argFieldInt64s, int64s: i})
}
func (a ArgUninons) Uint64(i uint64) ArgUninons {
	return append(a, argUnion{field: argFieldUint64, uint64: i})
}
func (a ArgUninons) Uint64s(i ...uint64) ArgUninons {
	return append(a, argUnion{field: argFieldUint64s, uint64s: i})
}
func (a ArgUninons) Float64(f float64) ArgUninons {
	return append(a, argUnion{field: argFieldFloat64, float64: f})
}
func (a ArgUninons) Float64s(f ...float64) ArgUninons {
	return append(a, argUnion{field: argFieldFloat64s, float64s: f})
}
func (a ArgUninons) Bool(f bool) ArgUninons {
	return append(a, argUnion{field: argFieldBool, bool: f})
}
func (a ArgUninons) Bools(f ...bool) ArgUninons {
	return append(a, argUnion{field: argFieldBools, bools: f})
}
func (a ArgUninons) String(f string) ArgUninons {
	return append(a, argUnion{field: argFieldString, string: f})
}
func (a ArgUninons) Strings(f ...string) ArgUninons {
	return append(a, argUnion{field: argFieldStrings, strings: f})
}
func (a ArgUninons) Bytes(b []byte) ArgUninons {
	return append(a, argUnion{field: argFieldByte, bytes: b})
}
func (a ArgUninons) BytesSlice(b ...[]byte) ArgUninons {
	return append(a, argUnion{field: argFieldBytes, bytess: b})
}
func (a ArgUninons) Time(t time.Time) ArgUninons {
	return append(a, argUnion{field: argFieldTime, time: t})
}
func (a ArgUninons) Times(t ...time.Time) ArgUninons {
	return append(a, argUnion{field: argFieldTimes, times: t})
}
func (a ArgUninons) NullString(nv ...NullString) ArgUninons {
	return append(a, argUnion{field: argFieldNullStrings, nullStrings: nv})
}
func (a ArgUninons) NullFloat64(nv ...NullFloat64) ArgUninons {
	return append(a, argUnion{field: argFieldNullFloat64s, nullFloat64s: nv})
}
func (a ArgUninons) NullInt64(nv ...NullInt64) ArgUninons {
	return append(a, argUnion{field: argFieldNullInt64s, nullInt64s: nv})
}
func (a ArgUninons) NullBool(nv NullBool) ArgUninons {
	return append(a, argUnion{field: argFieldNullBool, nullBool: nv})
}
func (a ArgUninons) NullTime(nv ...NullTime) ArgUninons {
	return append(a, argUnion{field: argFieldNullTimes, nullTimes: nv})
}
func (a ArgUninons) DriverValue(dvs ...driver.Valuer) ArgUninons {
	// Value is a value that drivers must be able to handle.
	// It is either nil or an instance of one of these types:
	//
	//   int64
	//   float64
	//   bool
	//   []byte
	//   string
	//   time.Time
	for _, dv := range dvs {
		v, err := dv.Value()
		if err != nil {
			// TODO: Either keep panic or delay the error until another function gets called which also returns an error.
			panic(errors.NewFatal(err, "[dbr] Driver.Value error for %#v", dv))
		}
		switch t := v.(type) {
		case nil:
			a = append(a, argUnion{field: argFieldNull})
		case int64:
			a = append(a, argUnion{field: argFieldInt64, int64: t})
		case float64:
			a = append(a, argUnion{field: argFieldFloat64, float64: t})
		case bool:
			a = append(a, argUnion{field: argFieldBool, bool: t})
		case []byte:
			a = append(a, argUnion{field: argFieldByte, bytes: t})
		case string:
			a = append(a, argUnion{field: argFieldString, string: t})
		case time.Time:
			a = append(a, argUnion{field: argFieldTime, time: t})
		default:
			panic(errors.NewNotSupportedf("[dbr] Type %#v not supported", t))
		}
	}
	return a
}
