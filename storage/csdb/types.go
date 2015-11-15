// Copyright 2015, Cyrill @ Schumacher.fm and the CoreStore contributors
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

package csdb

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"time"

	"fmt"

	"strings"

	"github.com/go-sql-driver/mysql"
	"github.com/ugorji/go/codec"
)

var (
	_          codec.Selfer = (*NullString)(nil)
	nullString              = []byte("null")
)

// NullString is a type that can be null or a string
type NullString struct {
	sql.NullString
}

// NullFloat64 is a type that can be null or a float64
type NullFloat64 struct {
	sql.NullFloat64
}

// NullInt64 is a type that can be null or an int
type NullInt64 struct {
	sql.NullInt64
}

// NullTime is a type that can be null or a time
type NullTime struct {
	mysql.NullTime
}

// NullBool is a type that can be null or a bool
type NullBool struct {
	sql.NullBool
}

// NewNullString creates a new database aware string.
func NewNullString(v interface{}) (n NullString) {
	n.Scan(v)
	return
}

// GoString satisfies the interface fmt.GoStringer when using %#v in Printf methods.
// Returns
// 		csdb.NewNullString(`...`,bool)
func (ns NullString) GoString() string {
	if ns.Valid && strings.ContainsRune(ns.String, '`') {
		// `This is my`string`
		ns.String = strings.Join(strings.Split(ns.String, "`"), "`+\"`\"+`")
		// `This is my`+"`"+`string`
	}

	ns.String = "`" + ns.String + "`"
	if !ns.Valid {
		ns.String = "nil"
	}
	return fmt.Sprintf("csdb.NewNullString(%s)", ns.String)
}

// CodecEncodeSelf for ugorji.go codec package
func (n NullString) CodecEncodeSelf(e *codec.Encoder) {
	if err := e.Encode(n.String); err != nil {
		PkgLog.Debug("csdb.NullString.CodecEncodeSelf", "err", err, "n", n)
	}
}

// CodecDecodeSelf  for ugorji.go codec package @todo write test ... not sure if ok
func (n *NullString) CodecDecodeSelf(d *codec.Decoder) {
	if err := d.Decode(&n.String); err != nil {
		PkgLog.Debug("csdb.NullString.CodecDecodeSelf", "err", err, "n", n)
	}
	// think about empty string and Valid value ...
}

// MarshalJSON correctly serializes a NullString to JSON
func (n NullString) MarshalJSON() ([]byte, error) {
	if n.Valid {
		j, e := json.Marshal(n.String)
		return j, e
	}
	return nullString, nil
}

// UnmarshalJSON correctly deserializes a NullString from JSON
func (n *NullString) UnmarshalJSON(b []byte) error {
	var s interface{}
	if err := json.Unmarshal(b, &s); err != nil {
		return err
	}
	return n.Scan(s)
}

// NewNullInt64 creates a new database aware type.
func NewNullInt64(v interface{}) (n NullInt64) {
	n.Scan(v)
	return
}

// MarshalJSON correctly serializes a NullInt64 to JSON
func (n *NullInt64) MarshalJSON() ([]byte, error) {
	if n.Valid {
		j, e := json.Marshal(n.Int64)
		return j, e
	}
	return nullString, nil
}

// CodecEncodeSelf for ugorji.go codec package
func (n *NullInt64) CodecEncodeSelf(e *codec.Encoder) {
	if err := e.Encode(n.Int64); err != nil {
		PkgLog.Debug("csdb.NullInt64.CodecEncodeSelf", "err", err, "n", n)
	}
}

// CodecDecodeSelf  for ugorji.go codec package @todo write test ... not sure if ok
func (n *NullInt64) CodecDecodeSelf(d *codec.Decoder) {
	if err := d.Decode(&n.Int64); err != nil {
		PkgLog.Debug("csdb.NullInt64.CodecDecodeSelf", "err", err, "n", n)
	}
}

// UnmarshalJSON correctly deserializes a NullInt64 from JSON
func (n *NullInt64) UnmarshalJSON(b []byte) error {
	var s interface{}
	if err := json.Unmarshal(b, &s); err != nil {
		return err
	}
	return n.Scan(s)
}

// NewNullFloat64 creates a new database aware type.
func NewNullFloat64(v interface{}) (n NullFloat64) {
	n.Scan(v)
	return
}

// MarshalJSON correctly serializes a NullFloat64 to JSON
func (n NullFloat64) MarshalJSON() ([]byte, error) {
	if n.Valid {
		j, e := json.Marshal(n.Float64)
		return j, e
	}
	return nullString, nil
}

// CodecEncodeSelf for ugorji.go codec package
func (n NullFloat64) CodecEncodeSelf(e *codec.Encoder) {
	if err := e.Encode(n.Float64); err != nil {
		PkgLog.Debug("csdb.NullFloat64.CodecEncodeSelf", "err", err, "n", n)
	}
}

// CodecDecodeSelf  for ugorji.go codec package @todo write test ... not sure if ok
func (n *NullFloat64) CodecDecodeSelf(d *codec.Decoder) {
	if err := d.Decode(&n.Float64); err != nil {
		PkgLog.Debug("csdb.NullFloat64.CodecDecodeSelf", "err", err, "n", n)
	}
}

// UnmarshalJSON correctly deserializes a NullFloat64 from JSON
func (n *NullFloat64) UnmarshalJSON(b []byte) error {
	var s interface{}
	if err := json.Unmarshal(b, &s); err != nil {
		return err
	}
	return n.Scan(s)
}

// NewNullTime creates a new database aware type.
func NewNullTime(v interface{}) (n NullTime) {
	n.Scan(v)
	return
}

// MarshalJSON correctly serializes a NullTime to JSON
func (n NullTime) MarshalJSON() ([]byte, error) {
	if n.Valid {
		j, e := json.Marshal(n.Time)
		return j, e
	}
	return nullString, nil
}

// CodecEncodeSelf for ugorji.go codec package
func (n NullTime) CodecEncodeSelf(e *codec.Encoder) {
	if err := e.Encode(n.Time); err != nil {
		PkgLog.Debug("csdb.NullTime.CodecEncodeSelf", "err", err, "n", n)
	}
}

// CodecDecodeSelf  for ugorji.go codec package @todo write test ... not sure if ok
func (n *NullTime) CodecDecodeSelf(d *codec.Decoder) {
	if err := d.Decode(&n.Time); err != nil {
		PkgLog.Debug("csdb.NullTime.CodecDecodeSelf", "err", err, "n", n)
	}
}

// UnmarshalJSON correctly deserializes a NullTime from JSON
func (n *NullTime) UnmarshalJSON(b []byte) error {
	// scan for null
	if bytes.Equal(b, nullString) {
		return n.Scan(nil)
	}
	// scan for JSON timestamp
	var t time.Time
	if err := json.Unmarshal(b, &t); err != nil {
		return err
	}
	return n.Scan(t)
}

// NewNullBool creates a new database aware type.
func NewNullBool(v interface{}) (n NullBool) {
	n.Scan(v)
	return
}

// MarshalJSON correctly serializes a NullBool to JSON
func (n *NullBool) MarshalJSON() ([]byte, error) {
	if n.Valid {
		j, e := json.Marshal(n.Bool)
		return j, e
	}
	return nullString, nil
}

// CodecEncodeSelf for ugorji.go codec package
func (n NullBool) CodecEncodeSelf(e *codec.Encoder) {
	if err := e.Encode(n.Bool); err != nil {
		PkgLog.Debug("csdb.NullBool.CodecEncodeSelf", "err", err, "n", n)
	}
}

// CodecDecodeSelf  for ugorji.go codec package @todo write test ... not sure if ok
func (n *NullBool) CodecDecodeSelf(d *codec.Decoder) {
	if err := d.Decode(&n.Bool); err != nil {
		PkgLog.Debug("csdb.NullBool.CodecDecodeSelf", "err", err, "n", n)
	}
}

// UnmarshalJSON correctly deserializes a NullBool from JSON
func (n *NullBool) UnmarshalJSON(b []byte) error {
	var s interface{}
	if err := json.Unmarshal(b, &s); err != nil {
		return err
	}
	return n.Scan(s)
}
