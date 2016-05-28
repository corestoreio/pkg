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
	"bytes"
	"encoding/json"
	"errors"
	"math"
	"net/http"
	"testing"
	"time"

	"github.com/corestoreio/csfw/storage/text"
	"github.com/stretchr/testify/assert"
)

const testKey = "MyTestKey"

func TestField_Bool(t *testing.T) {
	f := Bool(testKey, true)
	assert.Exactly(t, typeBool, f.fieldType)
	assert.Exactly(t, int64(1), f.int64)
	assert.Exactly(t, testKey, f.key)
}

func TestField_Float64(t *testing.T) {
	f := Float64(testKey, math.Pi)
	assert.Exactly(t, typeFloat64, f.fieldType)
	assert.Exactly(t, math.Pi, f.float64)
	assert.Exactly(t, testKey, f.key)
}

func TestField_Int(t *testing.T) {
	f := Int(testKey, math.MaxInt32)
	assert.Exactly(t, typeInt, f.fieldType)
	assert.Exactly(t, int64(math.MaxInt32), f.int64)
	assert.Exactly(t, testKey, f.key)
}

func TestField_Int64(t *testing.T) {
	f := Int64(testKey, math.MaxInt64)
	assert.Exactly(t, typeInt64, f.fieldType)
	assert.Exactly(t, int64(math.MaxInt64), f.int64)
	assert.Exactly(t, testKey, f.key)
}

func TestField_Uint(t *testing.T) {
	f := Uint(testKey, math.MaxUint32)
	assert.Exactly(t, typeInt, f.fieldType)
	assert.Exactly(t, int64(math.MaxUint32), f.int64)
	assert.Exactly(t, testKey, f.key)
}

func TestField_Uint64(t *testing.T) {
	f := Uint64(testKey, math.MaxUint64)
	assert.Exactly(t, typeInt64, f.fieldType)
	assert.Exactly(t, int64(math.MaxInt64), f.int64)
	assert.Exactly(t, testKey, f.key)
}

func TestField_String(t *testing.T) {
	const data = `16. “One is never alone with a rubber duck.” Douglas Adams`
	f := String(testKey, data)
	assert.Exactly(t, typeString, f.fieldType)
	assert.Exactly(t, data, f.string)
	assert.Exactly(t, testKey, f.key)
}

func TestField_Stringer(t *testing.T) {
	const data = `27. “Anything invented after you're thirty-five is against the natural order of things.” Douglas Adams`
	f := Stringer(testKey, bytes.NewBufferString(data))
	assert.Exactly(t, typeString, f.fieldType)
	assert.Exactly(t, data, f.string)
	assert.Exactly(t, testKey, f.key)
}

type gs struct {
	data interface{}
	err  error
}

func (g gs) MarshalText() ([]byte, error) {
	if g.err != nil {
		return nil, g.err
	}
	return g.data.([]byte), nil
}
func (gs) GoString() string { return "gs struct {}" }
func (g gs) MarshalJSON() ([]byte, error) {
	d, err := json.Marshal(g.data)
	if err != nil {
		g.err = err
	}
	return d, g.err
}
func (g gs) MarshalLog() (Field, error) {
	if g.err != nil {
		return Field{}, g.err
	}
	return String("ignored", "Val1x"), nil
}

func TestField_GoStringer(t *testing.T) {
	f := GoStringer(testKey, gs{})
	assert.Exactly(t, typeString, f.fieldType)
	assert.Exactly(t, "gs struct {}", f.string)
	assert.Exactly(t, testKey, f.key)
}

func TestField_Text(t *testing.T) {
	const data = `35. “My universe is my eyes and my ears. Anything else is hearsay.” Douglas Adams`
	f := Text(testKey, text.Chars(data))
	assert.Exactly(t, typeString, f.fieldType)
	assert.Exactly(t, data, f.string)
	assert.Exactly(t, testKey, f.key)
}
func TestField_TextError(t *testing.T) {
	var data = gs{data: nil, err: errors.New("Errr")}
	f := Text(testKey, data)
	assert.Exactly(t, typeString, f.fieldType)
	assert.Exactly(t, "[log] TextMarshaler: Errr", f.string)
	assert.Exactly(t, ErrorKeyName, f.key)
}

func TestField_JSON(t *testing.T) {
	const data = `12. “Reality is frequently inaccurate.” Douglas Adams`
	f := JSON(testKey, gs{data: data})
	assert.Exactly(t, typeString, f.fieldType)
	assert.Exactly(t, `"`+data+`"`, f.string)
	assert.Exactly(t, testKey, f.key)
}

func TestField_JSONError(t *testing.T) {
	f := JSON(testKey, gs{data: make(chan struct{})})
	assert.Exactly(t, typeString, f.fieldType)
	assert.Exactly(t, "[log] MarshalJSON: json: unsupported type: chan struct {}", f.string)
	assert.Exactly(t, ErrorKeyName, f.key)
}

func TestField_Time(t *testing.T) {
	now := time.Now()
	f := Time(testKey, now)
	assert.Exactly(t, typeInt64, f.fieldType)
	assert.Exactly(t, now.UnixNano(), f.int64)
	assert.Exactly(t, testKey, f.key)
}

func TestField_Duration(t *testing.T) {
	now := time.Hour * 2
	f := Duration(testKey, now)
	assert.Exactly(t, typeInt64, f.fieldType)
	assert.Exactly(t, now.Nanoseconds(), f.int64)
	assert.Exactly(t, testKey, f.key)
}

func TestField_Error(t *testing.T) {
	const data = `15. “There is no point in using the word 'impossible' to describe something that has clearly happened.” Douglas Adams`
	err := errors.New(data)
	f := Err(err)
	assert.Exactly(t, typeString, f.fieldType)
	assert.Exactly(t, data, f.string)
	assert.Exactly(t, ErrorKeyName, f.key)
}

func TestField_Object(t *testing.T) {
	req, _ := http.NewRequest("GET", "http://corestore.io", nil)
	req.RemoteAddr = "192.168.0.42"
	f := Object(testKey, req)
	assert.Exactly(t, typeObject, f.fieldType)
	assert.Exactly(t, req, f.obj)
	assert.Exactly(t, testKey, f.key)
}
