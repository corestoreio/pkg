// Copyright 2015-present, Cyrill @ Schumacher.fm and the CoreStore contributors
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

// Package null contains types which can be NULL in some storage engines.
//
// The aim is to import database/sql to avoid the dependency on SQL. Null values
// can occur everywhere hence we want to keep the deps minimal.
package null

import (
	"bytes"
	"strconv"
	"time"
)

const (
	sqlStrNullUC = "NULL"
	sqlStrNullLC = "null"
)

var (
	bTextNullUC  = []byte(sqlStrNullUC)
	bTextNullLC  = []byte(sqlStrNullLC)
	bTextFalseLC = []byte("false")
	bTextTrueLC  = []byte("true")
)

// Dialecter at an interface that wraps the diverse properties of individual
// SQL drivers.
type Dialecter interface {
	EscapeIdent(w *bytes.Buffer, ident string)
	EscapeBool(w *bytes.Buffer, b bool)
	EscapeString(w *bytes.Buffer, s string)
	EscapeTime(w *bytes.Buffer, t time.Time)
	EscapeBinary(w *bytes.Buffer, b []byte)
}

func writeFloat64(w *bytes.Buffer, f float64) (err error) {
	d := w.Bytes()
	w.Reset()
	_, err = w.Write(strconv.AppendFloat(d, f, 'g', -1, 64))
	return err
}

func writeInt64(w *bytes.Buffer, i int64) (err error) {
	d := w.Bytes()
	w.Reset()
	_, err = w.Write(strconv.AppendInt(d, i, 10))
	return err
}

func writeUint64(w *bytes.Buffer, i uint64) (err error) {
	d := w.Bytes()
	w.Reset()
	_, err = w.Write(strconv.AppendUint(d, i, 10))
	return err
}
