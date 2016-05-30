// Copyright (c) 2016 Uber Technologies, Inc.
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in
// all copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
// THE SOFTWARE.

// Copyright 2015-2016, Cyrill @ Schumacher.fm and the CoreStore contributors
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

package log

import (
	"encoding"
	"fmt"
	"math"
	"time"

	"github.com/corestoreio/csfw/util/bufferpool"
	"github.com/corestoreio/csfw/util/errors"
)

type fieldType uint8

// Type* constants define all available types which a field can contain.
const (
	typeBool fieldType = iota + 1
	typeInt
	typeInt64
	typeFloat64
	typeString
	typeStringer
	typeGoStringer
	typeObject
	typeMarshaler
	typeTextMarshaler
)

// JSONMarshaler is the interface implemented by types that
// can marshal themselves into valid JSON.
type JSONMarshaler interface {
	MarshalJSON() ([]byte, error)
}

// LogMarshaler allows user-defined types to efficiently add themselves to the
// logging context, and to selectively omit information which shouldn't be
// included in logs (e.g., passwords).
// Compatible to github.com/uber-go/zap
type LogMarshaler interface {
	MarshalLog(KeyValuer) error
}

// KeyValuer is an encoding-agnostic interface to add structured data to the
// logging context. Like maps, KeyValues aren't safe for concurrent use (though
// typical use shouldn't require locks).
//
// Compatible to github.com/uber-go/zap
type KeyValuer interface {
	AddBool(string, bool)
	AddFloat64(string, float64)
	AddInt(string, int)
	AddInt64(string, int64)
	AddMarshaler(string, LogMarshaler) error
	// AddObject uses reflection to serialize arbitrary objects, so it's slow and
	// allocation-heavy. Consider implementing the LogMarshaler interface instead.
	AddObject(string, interface{})
	AddString(string, string)
	Nest(string, func(KeyValuer) error) error
}

// Fields a slice of n Field types
type Fields []Field

// AddTo adds all fields within this slice to a KeyValue encoder.
// Breaks on first error.
func (fs Fields) AddTo(kv KeyValuer) error {
	for _, f := range fs {
		if err := f.AddTo(kv); err != nil {
			return errors.Wrap(err, "[log] Fields.AddTo")
		}
	}
	return nil
}

// MarshalLog satisfies the interface of log.LogMarshaler
func (fs Fields) MarshalLog(kv KeyValuer) error {
	return errors.Wrap(fs.AddTo(kv), "[log] Fields.Marshalog")
}

// ToString transforms multiple fields into a single string using the
// format of the type KVStringify.
func (fs Fields) ToString(msg string) string {
	buf := bufferpool.Get()
	defer bufferpool.Put(buf)
	wt := WriteTypes{W: buf}

	buf.WriteString(msg)
	if err := fs.AddTo(wt); err != nil {
		buf.WriteString(separator)
		buf.WriteString(ErrorKeyName)
		buf.WriteString(assignmentChar)
		buf.WriteString(errors.PrintLoc(err))
	}
	buf.WriteRune('\n')
	return buf.String()
}

// A Field is a deferred marshaling operation used to add a key-value pair to
// a logger's context. Keys and values are appropriately escaped for the current
// encoding scheme (e.g., JSON).
type Field struct {
	key string
	// fieldType specifies the used type. If 0 this struct is empty
	fieldType
	int64
	float64
	string
	obj interface{}
}

// AddTo adds a field to KeyValue encoder
func (f Field) AddTo(kv KeyValuer) error {
	switch f.fieldType {
	case typeBool:
		kv.AddBool(f.key, f.int64 == 1)
	case typeFloat64:
		kv.AddFloat64(f.key, f.float64)
	case typeInt:
		kv.AddInt(f.key, int(f.int64))
	case typeInt64:
		kv.AddInt64(f.key, f.int64)
	case typeString:
		kv.AddString(f.key, f.string)
	case typeStringer:
		kv.AddString(f.key, f.obj.(fmt.Stringer).String())
	case typeGoStringer:
		kv.AddString(f.key, f.obj.(fmt.GoStringer).GoString())
	case typeObject:
		kv.AddObject(f.key, f.obj)
	case typeMarshaler:
		return kv.AddMarshaler(f.key, f.obj.(LogMarshaler))
	case typeTextMarshaler:
		txt, err := f.obj.(encoding.TextMarshaler).MarshalText()
		if err != nil {
			return errors.Wrap(err, "[log] AddTo.TextMarshaler")
		}
		kv.AddString(f.key, string(txt))
	default:
		return errors.NewFatalf("[log] Unknown field type found: %v", f)
	}
	return nil
}

// Bool constructs a Field with the given key and value.
func Bool(key string, value bool) Field {
	var val int64
	if value {
		val = 1
	}
	return Field{key: key, fieldType: typeBool, int64: val}
}

// Float64 constructs a Field with the given key and value.
func Float64(key string, value float64) Field {
	return Field{key: key, fieldType: typeFloat64, float64: value}
}

// Int constructs a Field with the given key and value.
func Int(key string, val int) Field {
	return Field{key: key, fieldType: typeInt, int64: int64(val)}
}

// Int64 constructs a Field with the given key and value.
func Int64(key string, val int64) Field {
	return Field{key: key, fieldType: typeInt64, int64: val}
}

// Uint constructs a Field with the given key and value.
func Uint(key string, val uint) Field {
	return Field{key: key, fieldType: typeInt, int64: int64(val)}
}

// Uint64 constructs a Field with the given key and value.
// If val is bigger than math.MaxInt64 then val gets set to math.MaxInt64.
func Uint64(key string, val uint64) Field {
	if val > math.MaxInt64 {
		val = math.MaxInt64
	}
	return Field{key: key, fieldType: typeInt64, int64: int64(val)}
}

// String constructs a Field with the given key and value.
func String(key string, val string) Field {
	return Field{key: key, fieldType: typeString, string: val}
}

// Stringer constructs a Field with the given key and value. The value
// is the result of the String() method.
func Stringer(key string, val fmt.Stringer) Field {
	return Field{key: key, fieldType: typeStringer, obj: val}
}

// GoStringer constructs a Field with the given key and value. The value
// is the result of the GoString() method.
func GoStringer(key string, val fmt.GoStringer) Field {
	return Field{key: key, fieldType: typeGoStringer, obj: val}
}

// Text constructs a Field with the given key and value. The value
// is the result of the MarshalText() method.
func Text(key string, val encoding.TextMarshaler) Field {
	return Field{key: key, fieldType: typeTextMarshaler, obj: val}
}

// JSON constructs a Field with the given key and value. The value
// is the result of the MarshalJSON() method.
func JSON(key string, val JSONMarshaler) Field {
	j, err := val.MarshalJSON()
	if err != nil {
		return Err(errors.Wrap(err, "[log] MarshalJSON"))
	}
	return Field{key: key, fieldType: typeString, string: string(j)}
}

// Time constructs a Field with the given key and value. It represents a
// time.Time as nanoseconds since epoch.
func Time(key string, val time.Time) Field {
	return Int64(key, val.UnixNano())
}

// Duration constructs a Field with the given key and value. It represents
// durations as an integer number of nanoseconds.
func Duration(key string, val time.Duration) Field {
	return Field{key: key, fieldType: typeInt64, int64: val.Nanoseconds()}
}

// Error constructs a Field that stores err under the key log.ErrorKeyName. Prints
// <nil> if the error is nil.
func Err(err error) Field {
	if err == nil {
		return String(ErrorKeyName, "<nil>")
	}
	return String(ErrorKeyName, err.Error())
}

// ErrWithKey constructs a Field that stores err under a key. Prints
// <nil> if the error is nil.
func ErrWithKey(key string, err error) Field {
	if err == nil {
		return String(key, "<nil>")
	}
	return String(key, err.Error())
}

// Object constructs a field with the given key and an arbitrary object. It uses
// an encoding-appropriate, reflection-based function to serialize nearly any
// object into the logging context, but it's relatively slow and allocation-heavy.
//
// If encoding fails (e.g., trying to serialize a map[int]string to JSON), Object
// includes the error message in the final log output.
func Object(key string, val interface{}) Field {
	return Field{key: key, fieldType: typeObject, obj: val}
}

// Marshaler constructs a field with the given key and zap.LogMarshaler. It
// provides a flexible, but still type-safe and efficient, way to add
// user-defined types to the logging context.
func Marshaler(key string, val LogMarshaler) Field {
	return Field{key: key, fieldType: typeMarshaler, obj: val}
}

// Nest takes a key and a variadic number of Fields and creates a nested
// namespace.
func Nest(key string, fields ...Field) Field {
	return Field{key: key, fieldType: typeMarshaler, obj: Fields(fields)}
}
