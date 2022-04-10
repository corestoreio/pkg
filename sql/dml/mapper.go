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

package dml

import (
	"database/sql"
	"encoding"
	"fmt"
	"math"
	"strconv"
	"time"
	"unicode/utf8"

	"github.com/corestoreio/pkg/storage/null"
	"github.com/corestoreio/pkg/util/byteconv"
)

// ColumnMapper allows a type to load data from database query into its fields
// or return the fields values as arguments for a query. It's used in the
// rows.Next() for-loop. A ColumnMapper is usually a single record/row or in
// case of a slice a complete query result.
type ColumnMapper interface {
	// RowScan implementation must use function `Scan` to scan the values of the
	// query into its own type. See database/sql package for examples.
	MapColumns(rc *ColumnMap) error
}

// ColumnMap takes care that the table/view/identifiers are getting properly
// mapped to ColumnMapper interface. ColumnMap has two run modes either collect
// arguments from a type for running a SQL query OR to convert the sql.RawBytes
// into the desired final type. ColumnMap scans a *sql.Rows into a *sql.RawBytes
// slice without having a big memory overhead and not a single use of
// reflection. The conversion into the desired final type can happen without
// allocating of memory. It does not support streaming because neither
// database/sql does :-(  The method receiver functions have the same names as
// in type ColumnMap.
type ColumnMap struct {
	args []any

	// initialized gets set to true after the first call to Scan to initialize
	// the internal slices.
	initialized bool
	// CheckValidUTF8 if enabled checks if strings contains valid UTF-8 characters.
	CheckValidUTF8 bool
	// HasRows set to true if at least one row has been found.
	HasRows bool
	// Count increments on call to Scan.
	Count    uint64
	scanArgs []any // could be a sync.Pool but check it in benchmarks.
	scanCol  []scannedColumn
	// Columns contains the names of the column returned from the query or
	// needed to build a query aka reading the arguments in the ColumnMapper
	// interface. One should only read from the slice. Never modify it.
	columns    []string
	columnsLen int
	// scanErr is a delayed error and also used to avoid `if err != nil` in
	// generated code. This reduces the boiler plate code a lot! A trade off
	// between chainable API and too verbose error checking.
	scanErr    error
	index      int // current column index
	fieldCount int
}

// NewColumnMap exported for testing reasons.
func NewColumnMap(cap int, columns ...string) *ColumnMap {
	var a []any
	if cap > 0 {
		// Need the IF in case whe use sync.Pool
		a = make([]any, 0, cap)
	}
	cm := &ColumnMap{args: a}
	cm.setColumns(columns)
	return cm
}

// reset gets called when returning into the pool
func (b *ColumnMap) reset() {
	for i := range b.args {
		b.args[i] = nil
	}
	b.args = b.args[:0]
	b.initialized = false
	b.HasRows = false
	b.Count = 0
	b.scanArgs = b.scanArgs[:0]
	for i := range b.scanCol {
		b.scanCol[i].reset()
	}
	b.scanCol = b.scanCol[:0]
	b.columns = b.columns[:0]
	b.columnsLen = 0
	b.scanErr = nil
	b.index = 0
}

func (b *ColumnMap) setColumns(cols []string) {
	b.columns = cols
	b.columnsLen = len(cols)
	b.index = -1
}

// columnMapMode should be private because no need for a developer to take care
// of this mode in a variable.
type columnMapMode byte

func (m columnMapMode) String() string {
	return string(m)
}

// Those four constants represents the modes for ColumnMap.Mode. An upper case
// letter defines a collection and a lower case letter an entity.
const (
	ColumnMapEntityReadAll     columnMapMode = 'a'
	ColumnMapEntityReadSet     columnMapMode = 'r'
	ColumnMapCollectionReadSet columnMapMode = 'R'
	ColumnMapScan              columnMapMode = 'S' // can be used for both
)

// Mode returns a status byte of four different states. These states are getting
// used in the implementation of ColumnMapper. Each state represents a different
// action while scanning from the query or collecting arguments. ColumnMapper
// can be implemented by either a single type or a slice/map type. Slice or not
// slice requires different states. A primitive type must only handle mode
// ColumnMapEntityReadAll to return all requested fields. A slice type must
// handle additionally the cases ColumnMapEntityReadSet,
// ColumnMapCollectionReadSet and ColumnMapScan. See the examples. Documentation
// needs to be written better.
func (b *ColumnMap) Mode() (m columnMapMode) {
	if b.scanArgs != nil {
		// assign the column values from the DB to the structs and create new
		// structs in a slice.
		return ColumnMapScan
	}

	switch b.columnsLen {
	case 0:
		// Entity: read all mode; Collection jump into loop and pass on to
		// Entity.
		m = ColumnMapEntityReadAll
	case 1:
		// request certain column values as a slice.
		m = ColumnMapCollectionReadSet
	default:
		// Entity: calls the for cm.Next loop; Collection jump into loop and
		// pass on to Entity.
		m = ColumnMapEntityReadSet
	}
	return m
}

// scannedColumn represents an intermediate type (or DTO) to scan the
// driver.Values into. It supports the private types textRows and binaryRows in
// go-sql-driver/mysql. TextRows gets used during a normal query to read its
// result set and binaryRows gets used when a prepared statement gets executed
// and returns a result set. TextRows contains only byte slices whereas
// binaryRows contains already decoded types as defined in driver.Value. Avoids
// the reflection soup in database/sql.convertAssign. The supported data types
// depend on the this function: github.com/go-sql-driver/mysql/packets.go:1120
// `func (rows *binaryRows) readRow(dest []driver.Value) error` Then all the
// type functions (Int,String,Uint8, etc) in type ColumnMap can support all of
// the MySQL protocol field types. Hence we can support fieldTypeGeometry,
// fieldTypeJSON with custom decoders.
// Case for large uint64: if the DB value overflows math.MaxInt64, then it will
// be converted to a []byte slice, otherwise we have to deal with int64
type scannedColumn struct {
	field   byte      // i,j,f,b,y,s,t, n == null; nothing equals null of nil/empty
	bool    bool      // b
	int64   int64     // i
	float64 float64   // f double type
	string  string    // s
	time    time.Time // t
	byte    []byte    // y
}

func (s scannedColumn) String() string {
	switch s.field {
	case 'i':
		return strconv.FormatInt(s.int64, 10)
	case 'f':
		return strconv.FormatFloat(s.float64, 'f', -1, 64)
	case 'b':
		return strconv.FormatBool(s.bool)
	case 'y':
		return string(s.byte)
	case 's':
		return s.string
	case 't':
		return s.time.String()
	case 'n':
		return "<nil>"
	}
	return fmt.Sprintf("Field Type %q not supported", s.field)
}

func (s *scannedColumn) reset() {
	s.field = 0
	s.bool = false
	s.int64 = 0
	s.float64 = 0
	s.string = ""
	s.time = time.Time{}
	s.byte = s.byte[:0]
}

func (s *scannedColumn) Scan(src any) (err error) {
	switch val := src.(type) {
	case []byte: // most important case
		s.field = 'y'
		s.byte = val
	case int64:
		s.field = 'i'
		s.int64 = val
	case int: // sqlmock package requires this
		s.field = 'i'
		s.int64 = int64(val)
	case float32:
		s.field = 'f'
		s.float64 = float64(val)
	case float64:
		s.field = 'f'
		s.float64 = val
	case bool:
		s.field = 'b'
		s.bool = val
	case string:
		s.field = 's'
		s.string = val
	case time.Time:
		s.field = 't'
		s.time = val
	case nil:
		s.field = 'n'
	default:
		err = fmt.Errorf("[dml] 1649532685909 ColumnMap.Scan does not yet support type %T with value: %#v", val, val)
	}
	return err
}

func (b *ColumnMap) shouldCollectArgs() bool {
	return len(b.scanArgs) == 0 && cap(b.args) > 0
}

// Scan calls rows.Scan and builds an internal stack of sql.RawBytes for further
// processing and type conversion.
//
// Each function for a specific type converts the underlying byte slice at the
// current applied index (see function Index) to the appropriate type. You can
// call as many times as you want the specific functions. The underlying byte
// slice value is valid until the next call to rows.Next, rows.Scan or
// rows.Close. See the example for further usages.
func (b *ColumnMap) Scan(r *sql.Rows) error {
	if !b.initialized {
		cols, err := r.Columns()
		if err != nil {
			return err
		}

		b.setColumns(cols)
		if cap(b.scanCol) >= b.columnsLen { // reuse from pool!
			b.scanCol = b.scanCol[:b.columnsLen]
			b.scanArgs = b.scanArgs[:b.columnsLen]
		} else {
			b.scanCol = make([]scannedColumn, b.columnsLen)
			b.scanArgs = make([]any, b.columnsLen)
			for i := 0; i < b.columnsLen; i++ {
				b.scanArgs[i] = &b.scanCol[i]
			}
		}
		b.initialized = true
		b.Count = 0
		b.HasRows = true
	} else {
		b.Count++
	}
	return r.Scan(b.scanArgs...)
}

// Err returns the delayed error from one of the scans and parsings. Function is
// idempotent.
func (b *ColumnMap) Err() error {
	return b.scanErr
}

// Column returns the current column name after calling `Next`.
func (b *ColumnMap) Column() string {
	if b.index == -1 {
		b.index++ // in case of collections
	}
	if b.columnsLen == 0 && b.fieldCount > 0 {
		return strconv.FormatInt(int64(b.index), 10)
	}
	return b.columns[b.index]
}

// Next moves the internal index to the next position. It may return false if
// during RawBytes scanning an error has occurred.
func (b *ColumnMap) Next(fieldCount int) bool {
	b.fieldCount = fieldCount
	b.index++
	columnsLen := b.columnsLen
	if columnsLen == 0 {
		columnsLen = fieldCount
	}
	ok := b.index < columnsLen && b.scanErr == nil
	if !ok && b.scanErr == nil {
		// reset because the next row from the result-set will start or the next
		// Record/ColumnMapper collects the arguments. Only reset the index in
		// case of no-error because with an error you can get the column name
		// where the error has happened.
		b.index = -1
	}
	return ok
}

// Bool reads a bool value and appends it to the arguments slice or assigns the
// bool value stored in sql.RawBytes to the pointer. See the documentation for
// function Scan.
func (b *ColumnMap) Bool(ptr *bool) *ColumnMap {
	if b.shouldCollectArgs() {
		if ptr == nil {
			b.args = append(b.args, internalNULLNIL{})
		} else {
			b.args = append(b.args, *ptr)
		}
		return b
	}
	if b.scanErr == nil {
		// TODO benchmark if b.scanCol[b.index].field is faster than copying the struct in all type functions of ColumnMap
		switch v := b.scanCol[b.index]; v.field {
		case 'b':
			*ptr = v.bool // probably not implemented by go-sql-driver/mysql but to keep to compatibility reason for driver.Value
		case 'i':
			*ptr = v.int64 == 1
		case 'y':
			*ptr, _, b.scanErr = byteconv.ParseBool(v.byte)
			if b.scanErr != nil {
				b.scanErr = fmt.Errorf("[dml] 1649532847084 Column %q with error: %w", b.Column(), b.scanErr)
			}
		case 's':
			*ptr, b.scanErr = strconv.ParseBool(v.string)
			if b.scanErr != nil {
				b.scanErr = fmt.Errorf("[dml] 1649532850027 Column %q with error: %w", b.Column(), b.scanErr)
			}
		default:
			b.scanErr = fmt.Errorf("[dml] Column %q does not support field type: %q", b.Column(), v.field)
		}
	}
	return b
}

// NullBool reads a bool value and appends it to the arguments slice or assigns the
// bool value stored in sql.RawBytes to the pointer. See the documentation for
// function Scan.
func (b *ColumnMap) NullBool(ptr *null.Bool) *ColumnMap {
	if b.shouldCollectArgs() {
		if ptr == nil {
			b.args = append(b.args, internalNULLNIL{})
		} else {
			b.args = append(b.args, *ptr)
		}
		return b
	}
	if b.scanErr == nil {
		switch v := b.scanCol[b.index]; v.field {
		case 'b':
			ptr.Bool = v.bool
			ptr.Valid = true
		case 'n':
			ptr.Bool = false
			ptr.Valid = false
		case 'i':
			ptr.Bool = v.int64 == 1
			ptr.Valid = true
		case 'y':
			*ptr, b.scanErr = null.MakeBoolFromByte(v.byte)
			if b.scanErr != nil {
				b.scanErr = fmt.Errorf("[dml] 1649532889274 Column %q with error: %w", b.Column(), b.scanErr)
			}
		case 's':
			if v.string != "" {
				ptr.Bool, b.scanErr = strconv.ParseBool(v.string)
				if b.scanErr != nil {
					b.scanErr = fmt.Errorf("[dml] 1649532891950 Column %q with error: %w", b.Column(), b.scanErr)
				}
				ptr.Valid = b.scanErr == nil
			} else {
				ptr.Bool = false
				ptr.Valid = false
			}
		default:
			b.scanErr = fmt.Errorf("[dml] Column %q does not support field type: %q", b.Column(), v.field)
		}
	}
	return b
}

// Int reads an int value and appends it to the arguments slice or assigns the
// int value stored in sql.RawBytes to the pointer. See the documentation for
// function Scan.
func (b *ColumnMap) Int(ptr *int) *ColumnMap {
	if b.shouldCollectArgs() {
		if ptr == nil {
			b.args = append(b.args, internalNULLNIL{})
		} else {
			b.args = append(b.args, *ptr)
		}
		return b
	}
	if b.scanErr == nil {
		switch v := b.scanCol[b.index]; v.field {
		case 'i':
			*ptr = int(v.int64)
		case 'y':
			var i64 int64
			i64, _, b.scanErr = byteconv.ParseInt(v.byte)
			*ptr = int(i64)
			if b.scanErr != nil {
				b.scanErr = fmt.Errorf("[dml] 1649532900534 Column %q with error: %w", b.Column(), b.scanErr)
			}
		default:
			b.scanErr = fmt.Errorf("[dml] Column %q does not support field type: %q", b.Column(), v.field)
		}
	}
	return b
}

// Int64 reads a int64 value and appends it to the arguments slice or assigns
// the int64 value stored in sql.RawBytes to the pointer. See the documentation
// for function Scan.
func (b *ColumnMap) Int64(ptr *int64) *ColumnMap {
	if b.shouldCollectArgs() {
		if ptr == nil {
			b.args = append(b.args, internalNULLNIL{})
		} else {
			b.args = append(b.args, *ptr)
		}
		return b
	}
	if b.scanErr != nil {
		return b
	}
	switch v := b.scanCol[b.index]; v.field {
	case 'i':
		*ptr = v.int64
	case 'y':
		*ptr, _, b.scanErr = byteconv.ParseInt(v.byte)
		if b.scanErr != nil {
			b.scanErr = fmt.Errorf("[dml] 1649532905258 Column %q with error: %w", b.Column(), b.scanErr)
		}
	default:
		b.scanErr = fmt.Errorf("[dml] Column %q does not support field type: %q", b.Column(), v.field)
	}
	return b
}

// Int32 reads a int32 value and appends it to the arguments slice or assigns
// the int32 value stored in sql.RawBytes to the pointer. See the documentation
// for function Scan.
func (b *ColumnMap) Int32(ptr *int32) *ColumnMap {
	if b.shouldCollectArgs() {
		if ptr == nil {
			b.args = append(b.args, internalNULLNIL{})
		} else {
			b.args = append(b.args, *ptr)
		}
		return b
	}
	if b.scanErr != nil {
		return b
	}
	switch v := b.scanCol[b.index]; v.field {
	case 'i':
		*ptr = int32(v.int64)
	case 'y':
		var i64 int64
		i64, _, b.scanErr = byteconv.ParseInt(v.byte)
		switch {
		case b.scanErr != nil:
			b.scanErr = fmt.Errorf("[dml] 1649532910346 Column %q with error: %w", b.Column(), b.scanErr)
		case i64 > math.MaxInt32, i64 < -math.MaxInt32:
			b.scanErr = fmt.Errorf("[dml] 1649533361246 Column %q overflows int32: i64: %d", b.Column(), i64)
		default:
			*ptr = int32(i64)
		}
	default:
		b.scanErr = fmt.Errorf("[dml] Column %q does not support field type: %q", b.Column(), v.field)
	}
	return b
}

// Int16 reads a int16 value and appends it to the arguments slice or assigns
// the int16 value stored in sql.RawBytes to the pointer. See the documentation
// for function Scan.
func (b *ColumnMap) Int16(ptr *int16) *ColumnMap {
	if b.shouldCollectArgs() {
		if ptr == nil {
			b.args = append(b.args, internalNULLNIL{})
		} else {
			b.args = append(b.args, *ptr)
		}
		return b
	}
	if b.scanErr != nil {
		return b
	}
	switch v := b.scanCol[b.index]; v.field {
	case 'i':
		*ptr = int16(v.int64)
	case 'y':
		var i64 int64
		i64, _, b.scanErr = byteconv.ParseInt(v.byte)
		switch {
		case b.scanErr != nil:
			b.scanErr = fmt.Errorf("[dml] 1649532915743 Column %q with error: %w", b.Column(), b.scanErr)
		case i64 > math.MaxInt16, i64 < -math.MaxInt16:
			b.scanErr = fmt.Errorf("[dml] 1649533385933 Column %q overflows int16: i64: %d", b.Column(), i64)
		default:
			*ptr = int16(i64)
		}
	default:
		b.scanErr = fmt.Errorf("[dml] 1649533389507 Column %q does not support field type: %q", b.Column(), v.field)
	}
	return b
}

// Int8 reads a int8 value and appends it to the arguments slice or assigns
// the int8 value stored in sql.RawBytes to the pointer. See the documentation
// for function Scan.
func (b *ColumnMap) Int8(ptr *int8) *ColumnMap {
	if b.shouldCollectArgs() {
		if ptr == nil {
			b.args = append(b.args, internalNULLNIL{})
		} else {
			b.args = append(b.args, *ptr)
		}
		return b
	}
	if b.scanErr != nil {
		return b
	}
	switch v := b.scanCol[b.index]; v.field {
	case 'i':
		*ptr = int8(v.int64)
	case 'y':
		var i64 int64
		i64, _, b.scanErr = byteconv.ParseInt(v.byte)
		switch {
		case b.scanErr != nil:
			b.scanErr = fmt.Errorf("[dml] 1649532921154 Column %q with error: %w", b.Column(), b.scanErr)
		case i64 > math.MaxInt8, i64 < -math.MaxInt8:
			b.scanErr = fmt.Errorf("[dml] 1649533410105 Column %q overflows int8: i64: %d", b.Column(), i64)
		default:
			*ptr = int8(i64)
		}
	default:
		b.scanErr = fmt.Errorf("[dml] 1649533413976 Column %q does not support field type: %q", b.Column(), v.field)
	}
	return b
}

// NullInt64 reads an int64 value and appends it to the arguments slice or
// assigns the int64 value stored in sql.RawBytes to the pointer. See the
// documentation for function Scan.
func (b *ColumnMap) NullInt64(ptr *null.Int64) *ColumnMap {
	if b.shouldCollectArgs() {
		if ptr == nil {
			b.args = append(b.args, internalNULLNIL{})
		} else {
			b.args = append(b.args, *ptr)
		}
		return b
	}
	if b.scanErr == nil {
		switch v := b.scanCol[b.index]; v.field {
		case 'i':
			ptr.Int64 = v.int64
			ptr.Valid = true
		case 'n':
			ptr.Int64 = 0
			ptr.Valid = false
		case 'y':
			*ptr, b.scanErr = null.MakeInt64FromByte(v.byte)
			if b.scanErr != nil {
				b.scanErr = fmt.Errorf("[dml] 1649532925895 Column %q with error: %w", b.Column(), b.scanErr)
			}
		default:
			b.scanErr = fmt.Errorf("[dml] Column %q does not support field type: %q", b.Column(), v.field)
		}
	}
	return b
}

// NullInt32 reads an int32 value and appends it to the arguments slice or
// assigns the int32 value stored in sql.RawBytes to the pointer. See the
// documentation for function Scan.
func (b *ColumnMap) NullInt32(ptr *null.Int32) *ColumnMap {
	if b.shouldCollectArgs() {
		if ptr == nil {
			b.args = append(b.args, internalNULLNIL{})
		} else {
			b.args = append(b.args, *ptr)
		}
		return b
	}
	if b.scanErr == nil {
		switch v := b.scanCol[b.index]; v.field {
		case 'i':
			ptr.Int32 = int32(v.int64) // TODO check overflow
			ptr.Valid = true
		case 'n':
			ptr.Int32 = 0
			ptr.Valid = false
		case 'y':
			*ptr, b.scanErr = null.MakeInt32FromByte(v.byte)
			if b.scanErr != nil {
				b.scanErr = fmt.Errorf("[dml] 1649532966947 Column %q with error: %w", b.Column(), b.scanErr)
			}
		default:
			b.scanErr = fmt.Errorf("[dml] Column %q does not support field type: %q", b.Column(), v.field)
		}
	}
	return b
}

// NullInt16 reads an int16 value and appends it to the arguments slice or
// assigns the int16 value stored in sql.RawBytes to the pointer. See the
// documentation for function Scan.
func (b *ColumnMap) NullInt16(ptr *null.Int16) *ColumnMap {
	if b.shouldCollectArgs() {
		if ptr == nil {
			b.args = append(b.args, internalNULLNIL{})
		} else {
			b.args = append(b.args, *ptr)
		}
		return b
	}
	if b.scanErr == nil {
		switch v := b.scanCol[b.index]; v.field {
		case 'i':
			ptr.Int16 = int16(v.int64) // TODO check overflow
			ptr.Valid = true
		case 'n':
			ptr.Int16 = 0
			ptr.Valid = false
		case 'y':
			*ptr, b.scanErr = null.MakeInt16FromByte(v.byte)
			if b.scanErr != nil {
				b.scanErr = fmt.Errorf("[dml] 1649532971518 Column %q with error: %w", b.Column(), b.scanErr)
			}
		default:
			b.scanErr = fmt.Errorf("[dml] Column %q does not support field type: %q", b.Column(), v.field)
		}
	}
	return b
}

// NullInt8 reads an int8 value and appends it to the arguments slice or
// assigns the int8 value stored in sql.RawBytes to the pointer. See the
// documentation for function Scan.
func (b *ColumnMap) NullInt8(ptr *null.Int8) *ColumnMap {
	if b.shouldCollectArgs() {
		if ptr == nil {
			b.args = append(b.args, internalNULLNIL{})
		} else {
			b.args = append(b.args, *ptr)
		}
		return b
	}
	if b.scanErr == nil {
		switch v := b.scanCol[b.index]; v.field {
		case 'i':
			ptr.Int8 = int8(v.int64) // TODO check overflow
			ptr.Valid = true
		case 'n':
			ptr.Int8 = 0
			ptr.Valid = false
		case 'y':
			*ptr, b.scanErr = null.MakeInt8FromByte(v.byte)
			if b.scanErr != nil {
				b.scanErr = fmt.Errorf("[dml] 1649532978741 Column %q with error: %w", b.Column(), b.scanErr)
			}
		default:
			b.scanErr = fmt.Errorf("[dml] Column %q does not support field type: %q", b.Column(), v.field)
		}
	}
	return b
}

// NullInt64 reads an int64 value and appends it to the arguments slice or
// assigns the int64 value stored in sql.RawBytes to the pointer. See the
// documentation for function Scan.
func (b *ColumnMap) NullUint64(ptr *null.Uint64) *ColumnMap {
	if b.shouldCollectArgs() {
		if ptr == nil {
			b.args = append(b.args, internalNULLNIL{})
		} else {
			b.args = append(b.args, *ptr)
		}
		return b
	}
	if b.scanErr == nil {
		switch v := b.scanCol[b.index]; v.field {
		case 'i':
			ptr.Uint64 = uint64(v.int64)
			ptr.Valid = true
		case 'n':
			ptr.Uint64 = 0
			ptr.Valid = false
		case 'y':
			*ptr, b.scanErr = null.MakeUint64FromByte(v.byte)
			if b.scanErr != nil {
				b.scanErr = fmt.Errorf("[dml] 1649532983260 Column %q with error: %w", b.Column(), b.scanErr)
			}
		default:
			b.scanErr = fmt.Errorf("[dml] Column %q does not support field type: %q", b.Column(), v.field)
		}
	}
	return b
}

// NullInt32 reads an int32 value and appends it to the arguments slice or
// assigns the int32 value stored in sql.RawBytes to the pointer. See the
// documentation for function Scan.
func (b *ColumnMap) NullUint32(ptr *null.Uint32) *ColumnMap {
	if b.shouldCollectArgs() {
		if ptr == nil {
			b.args = append(b.args, internalNULLNIL{})
		} else {
			b.args = append(b.args, *ptr)
		}
		return b
	}
	if b.scanErr == nil {
		switch v := b.scanCol[b.index]; v.field {
		case 'i':
			ptr.Uint32 = uint32(v.int64) // TODO check overflow
			ptr.Valid = true
		case 'n':
			ptr.Uint32 = 0
			ptr.Valid = false
		case 'y':
			*ptr, b.scanErr = null.MakeUint32FromByte(v.byte)
			if b.scanErr != nil {
				b.scanErr = fmt.Errorf("[dml] 1649532999016 Column %q with error: %w", b.Column(), b.scanErr)
			}
		default:
			b.scanErr = fmt.Errorf("[dml] Column %q does not support field type: %q", b.Column(), v.field)
		}
	}
	return b
}

// NullInt16 reads an int16 value and appends it to the arguments slice or
// assigns the int16 value stored in sql.RawBytes to the pointer. See the
// documentation for function Scan.
func (b *ColumnMap) NullUint16(ptr *null.Uint16) *ColumnMap {
	if b.shouldCollectArgs() {
		if ptr == nil {
			b.args = append(b.args, internalNULLNIL{})
		} else {
			b.args = append(b.args, *ptr)
		}
		return b
	}
	if b.scanErr == nil {
		switch v := b.scanCol[b.index]; v.field {
		case 'i':
			ptr.Uint16 = uint16(v.int64) // TODO check overflow
			ptr.Valid = true
		case 'n':
			ptr.Uint16 = 0
			ptr.Valid = false
		case 'y':
			*ptr, b.scanErr = null.MakeUint16FromByte(v.byte)
			if b.scanErr != nil {
				b.scanErr = fmt.Errorf("[dml] 1649533005807 Column %q with error: %w", b.Column(), b.scanErr)
			}
		default:
			b.scanErr = fmt.Errorf("[dml] Column %q does not support field type: %q", b.Column(), v.field)
		}
	}
	return b
}

// NullInt8 reads an int8 value and appends it to the arguments slice or
// assigns the int8 value stored in sql.RawBytes to the pointer. See the
// documentation for function Scan.
func (b *ColumnMap) NullUint8(ptr *null.Uint8) *ColumnMap {
	if b.shouldCollectArgs() {
		if ptr == nil {
			b.args = append(b.args, internalNULLNIL{})
		} else {
			b.args = append(b.args, *ptr)
		}
		return b
	}
	if b.scanErr == nil {
		switch v := b.scanCol[b.index]; v.field {
		case 'i':
			ptr.Uint8 = uint8(v.int64) // TODO check overflow
			ptr.Valid = true
		case 'n':
			ptr.Uint8 = 0
			ptr.Valid = false
		case 'y':
			*ptr, b.scanErr = null.MakeUint8FromByte(v.byte)
			if b.scanErr != nil {
				b.scanErr = fmt.Errorf("[dml] 1649533012588 Column %q with error: %w", b.Column(), b.scanErr)
			}
		default:
			b.scanErr = fmt.Errorf("[dml] Column %q does not support field type: %q", b.Column(), v.field)
		}
	}
	return b
}

// Float64 reads a float64 value and appends it to the arguments slice or
// assigns the float64 value stored in sql.RawBytes to the pointer. See the
// documentation for function Scan.
func (b *ColumnMap) Float64(ptr *float64) *ColumnMap {
	if b.shouldCollectArgs() {
		if ptr == nil {
			b.args = append(b.args, internalNULLNIL{})
		} else {
			b.args = append(b.args, *ptr)
		}
		return b
	}
	if b.scanErr == nil {
		switch v := b.scanCol[b.index]; v.field {
		case 'f':
			*ptr = v.float64
		case 'y':
			*ptr, _, b.scanErr = byteconv.ParseFloat(v.byte)
			if b.scanErr != nil {
				b.scanErr = fmt.Errorf("[dml] 1649533032763 Column %q with error: %w", b.Column(), b.scanErr)
			}
		default:
			b.scanErr = fmt.Errorf("[dml] Column %q does not support field type: %q", b.Column(), v.field)
		}
	}
	return b
}

// Decimal reads a Decimal value and appends it to the arguments slice or
// assigns the numeric value stored in sql.RawBytes to the pointer. See the
// documentation for function Scan.
func (b *ColumnMap) Decimal(ptr *null.Decimal) *ColumnMap {
	if b.shouldCollectArgs() {
		if v := ptr.String(); ptr == nil || v == sqlStrNullUC {
			b.args = append(b.args, internalNULLNIL{})
		} else {
			b.args = append(b.args, v)
		}
		return b
	}
	if b.scanErr == nil {
		switch v := b.scanCol[b.index]; v.field {
		case 'f':
			*ptr, b.scanErr = null.MakeDecimalFloat64(v.float64)
		case 'y':
			*ptr, b.scanErr = null.MakeDecimalBytes(v.byte)
		case 's':
			*ptr, b.scanErr = null.MakeDecimalBytes([]byte(v.string)) // mostly used for testing
		case 'n':
			ptr.Valid = false
		default:
			b.scanErr = fmt.Errorf("[dml] Column %q does not support field type: %q", b.Column(), v.field)
		}
	}
	return b
}

// NullFloat64 reads a float64 value and appends it to the arguments slice or
// assigns the float64 value stored in sql.RawBytes to the pointer. See the
// documentation for function Scan.
func (b *ColumnMap) NullFloat64(ptr *null.Float64) *ColumnMap {
	if b.shouldCollectArgs() {
		if ptr == nil {
			b.args = append(b.args, internalNULLNIL{})
		} else {
			b.args = append(b.args, *ptr)
		}
		return b
	}
	if b.scanErr == nil {
		switch v := b.scanCol[b.index]; v.field {
		case 'f':
			ptr.Float64 = v.float64
			ptr.Valid = true
		case 'y':
			*ptr, b.scanErr = null.MakeFloat64FromByte(v.byte)
			if b.scanErr != nil {
				b.scanErr = fmt.Errorf("[dml] 1649533037000 Column %q with error: %w", b.Column(), b.scanErr)
			}
		case 'n':
			ptr.Valid = false
		default:
			b.scanErr = fmt.Errorf("[dml] Column %q does not support field type: %q", b.Column(), v.field)
		}
	}
	return b
}

// Uint reads an uint value and appends it to the arguments slice or assigns the
// uint value stored in sql.RawBytes to the pointer. See the documentation for
// function Scan.
func (b *ColumnMap) Uint(ptr *uint) *ColumnMap {
	if b.shouldCollectArgs() {
		if ptr == nil {
			b.args = append(b.args, internalNULLNIL{})
		} else {
			b.args = append(b.args, *ptr)
		}
		return b
	}
	if b.scanErr == nil {
		switch v := b.scanCol[b.index]; v.field {
		case 'i':
			*ptr = uint(v.int64)
		case 'y':
			var u64 uint64
			u64, _, b.scanErr = byteconv.ParseUint(v.byte, 10, strconv.IntSize)
			*ptr = uint(u64)
			if b.scanErr != nil {
				b.scanErr = fmt.Errorf("[dml] 1649533041002 Column %q with error: %w", b.Column(), b.scanErr)
			}
		default:
			b.scanErr = fmt.Errorf("[dml] Column %q does not support field type: %q", b.Column(), v.field)
		}
	}
	return b
}

// Uint8 reads an uint8 value and appends it to the arguments slice or assigns
// the uint8 value stored in sql.RawBytes to the pointer. See the documentation
// for function Scan.
func (b *ColumnMap) Uint8(ptr *uint8) *ColumnMap {
	if b.shouldCollectArgs() {
		if ptr == nil {
			b.args = append(b.args, internalNULLNIL{})
		} else {
			b.args = append(b.args, *ptr)
		}
		return b
	}
	if b.scanErr == nil {
		switch v := b.scanCol[b.index]; v.field {
		case 'i':
			*ptr = uint8(v.int64)
		case 'y':
			var u64 uint64
			u64, _, b.scanErr = byteconv.ParseUint(v.byte, 10, 8)
			*ptr = uint8(u64)
			if b.scanErr != nil {
				b.scanErr = fmt.Errorf("[dml] 1649533045240 Column %q with error: %w", b.Column(), b.scanErr)
			}
		default:
			b.scanErr = fmt.Errorf("[dml] Column %q does not support field type: %q", b.Column(), v.field)
		}
	}
	return b
}

// Uint16 reads an uint16 value and appends it to the arguments slice or assigns
// the uint16 value stored in sql.RawBytes to the pointer. See the documentation
// for function Scan.
func (b *ColumnMap) Uint16(ptr *uint16) *ColumnMap {
	if b.shouldCollectArgs() {
		if ptr == nil {
			b.args = append(b.args, internalNULLNIL{})
		} else {
			b.args = append(b.args, *ptr)
		}
		return b
	}
	if b.scanErr == nil {
		switch v := b.scanCol[b.index]; v.field {
		case 'i':
			*ptr = uint16(v.int64)
		case 'y':
			var u64 uint64
			u64, _, b.scanErr = byteconv.ParseUint(v.byte, 10, 16)
			*ptr = uint16(u64)
			if b.scanErr != nil {
				b.scanErr = fmt.Errorf("[dml] 1649533050547 Column %q with error: %w", b.Column(), b.scanErr)
			}
		default:
			b.scanErr = fmt.Errorf("[dml] Column %q does not support field type: %q", b.Column(), v.field)
		}
	}
	return b
}

// Uint32 reads an uint32 value and appends it to the arguments slice or assigns
// the uint32 value stored in sql.RawBytes to the pointer. See the documentation
// for function Scan.
func (b *ColumnMap) Uint32(ptr *uint32) *ColumnMap {
	if b.shouldCollectArgs() {
		if ptr == nil {
			b.args = append(b.args, internalNULLNIL{})
		} else {
			b.args = append(b.args, *ptr)
		}
		return b
	}
	if b.scanErr == nil {
		switch v := b.scanCol[b.index]; v.field {
		case 'i':
			*ptr = uint32(v.int64)
		case 'y':
			var u64 uint64
			u64, _, b.scanErr = byteconv.ParseUint(v.byte, 10, 32)
			*ptr = uint32(u64)
			if b.scanErr != nil {
				b.scanErr = fmt.Errorf("[dml] 1649533055112 Column %q with error: %w", b.Column(), b.scanErr)
			}
		default:
			b.scanErr = fmt.Errorf("[dml] Column %q does not support field type: %q", b.Column(), v.field)
		}
	}
	return b
}

// Uint64 reads an uint64 value and appends it to the arguments slice or assigns
// the uint64 value stored in sql.RawBytes to the pointer. See the documentation
// for function Scan.
func (b *ColumnMap) Uint64(ptr *uint64) *ColumnMap {
	if b.shouldCollectArgs() {
		if ptr == nil {
			b.args = append(b.args, internalNULLNIL{})
		} else {
			b.args = append(b.args, *ptr)
		}
		return b
	}
	if b.scanErr == nil {
		switch v := b.scanCol[b.index]; v.field {
		case 'i':
			*ptr = uint64(v.int64)
		case 'y':
			*ptr, _, b.scanErr = byteconv.ParseUint(v.byte, 10, strconv.IntSize)
			if b.scanErr != nil {
				b.scanErr = fmt.Errorf("[dml] 1649533059094 Column %q with error: %w", b.Column(), b.scanErr)
			}
		default:
			b.scanErr = fmt.Errorf("[dml] Column %q does not support field type: %q", b.Column(), v.field)
		}
	}
	return b
}

type ioWriter interface {
	Write(p []byte) (n int, err error)
}

// Debug writes the column names with their values into `w`. The output format
// might change.
func (b *ColumnMap) Debug(w ioWriter) (err error) {
	nl := []byte("\n")
	tNil := []byte(": <nil>")
	for i, c := range b.columns {
		if i > 0 {
			_, _ = w.Write(nl)
		}
		_, _ = w.Write([]byte(c))
		b := b.scanCol[i]
		if b.field == 'n' {
			_, _ = w.Write(tNil)
		} else if _, err = fmt.Fprintf(w, ": %q", b); err != nil {
			return err
		}
	}
	return nil
}

// Byte reads a []byte value and appends it to the arguments slice or assigns
// the []byte value stored in sql.RawBytes to the pointer. See the documentation
// for function Scan.
func (b *ColumnMap) Byte(ptr *[]byte) *ColumnMap {
	if b.shouldCollectArgs() {
		if ptr == nil {
			b.args = append(b.args, internalNULLNIL{})
		} else {
			b.args = append(b.args, *ptr)
		}
		return b
	}
	if b.scanErr == nil {
		switch v := b.scanCol[b.index]; v.field {
		case 's':
			*ptr = append((*ptr)[:0], v.string...)
		case 'y':
			*ptr = append((*ptr)[:0], v.byte...)
		case 'n':
			*ptr = nil
		default:
			b.scanErr = fmt.Errorf("[dml] Column %q does not support field type: %q", b.Column(), v.field)
		}
	}
	return b
}

// Text allows to encode an object to its text representation when arguments are
// requested and to decode a byte slice into its object when data is retrieved
// from the server. Use this function for JSON, XML, YAML, etc formats. This
// function can check for valid UTF8 characters, see field CheckValidUTF8.
func (b *ColumnMap) Text(enc interface {
	encoding.TextMarshaler
	encoding.TextUnmarshaler
}) *ColumnMap {
	if b.scanErr != nil {
		return b
	}
	if b.shouldCollectArgs() {
		var data []byte
		data, b.scanErr = enc.MarshalText()
		if b.CheckValidUTF8 && !utf8.Valid(data) {
			b.scanErr = fmt.Errorf("[dml] 1649533441330 Column Index %d at position %d contains invalid UTF-8 characters", b.index, b.Count)
		} else {
			b.args = append(b.args, data)
		}
		return b
	}

	switch v := b.scanCol[b.index]; v.field {
	case 'y', 'n':
		if b.CheckValidUTF8 && !utf8.Valid(v.byte) {
			b.scanErr = fmt.Errorf("[dml] 1649533448677 Column %q IDX:%d contains invalid UTF-8 characters", b.Column(), b.Count)
		} else if b.scanErr = enc.UnmarshalText(v.byte); b.scanErr != nil {
			b.scanErr = fmt.Errorf("[dml] 1649533063790 Column %q with error: %w", b.Column(), b.scanErr)
		}
	default:
		b.scanErr = fmt.Errorf("[dml] Column %q does not support field type: %q", b.Column(), v.field)
	}
	return b
}

// Binary allows to encode an object to its binary representation when arguments
// are requested and to decode a byte slice into its object when data is
// retrieved from the server. Use this function for GOB, Protocol Buffers, etc
// formats.
func (b *ColumnMap) Binary(enc interface {
	encoding.BinaryMarshaler
	encoding.BinaryUnmarshaler
}) *ColumnMap {
	if b.scanErr != nil {
		return b
	}
	if b.shouldCollectArgs() {
		var data []byte
		data, b.scanErr = enc.MarshalBinary()
		b.args = append(b.args, data)
		return b
	}

	switch v := b.scanCol[b.index]; v.field {
	case 'y', 'n':
		if b.scanErr = enc.UnmarshalBinary(v.byte); b.scanErr != nil {
			b.scanErr = fmt.Errorf("[dml] 1649533068402 Column %q with error: %w", b.Column(), b.scanErr)
		}
	default:
		b.scanErr = fmt.Errorf("[dml] Column %q does not support field type: %q", b.Column(), v.field)
	}
	return b
}

// String reads a string value and appends it to the arguments slice or assigns
// the string value stored in sql.RawBytes to the pointer. See the documentation
// for function Scan.
func (b *ColumnMap) String(ptr *string) *ColumnMap {
	if b.shouldCollectArgs() {
		if ptr == nil {
			b.args = append(b.args, internalNULLNIL{})
		} else {
			b.args = append(b.args, *ptr)
		}
		return b
	}

	if b.scanErr == nil {
		switch v := b.scanCol[b.index]; v.field {
		case 's':
			if b.CheckValidUTF8 && !utf8.ValidString(v.string) {
				b.scanErr = fmt.Errorf("[dml] 1649533459370 Column Index %d at position %d contains invalid UTF-8 characters", b.index, b.Count)
			} else {
				*ptr = v.string
			}
		case 'y':
			if b.CheckValidUTF8 && !utf8.Valid(v.byte) {
				b.scanErr = fmt.Errorf("[dml] 1649533465808 Column Index %d at position %d contains invalid UTF-8 characters", b.index, b.Count)
			} else {
				*ptr = string(v.byte)
			}
		default:
			b.scanErr = fmt.Errorf("[dml] 1649533477981 Column %q does not support field type: %q", b.Column(), v.field)
		}
	}
	return b
}

// NullString reads a string value and appends it to the arguments slice or
// assigns the string value stored in sql.RawBytes to the pointer. See the
// documentation for function Scan.
func (b *ColumnMap) NullString(ptr *null.String) *ColumnMap {
	if b.shouldCollectArgs() {
		if ptr == nil {
			b.args = append(b.args, internalNULLNIL{})
		} else {
			b.args = append(b.args, *ptr)
		}
		return b
	}

	if b.scanErr == nil {
		switch v := b.scanCol[b.index]; v.field {
		case 's':
			if b.CheckValidUTF8 && !utf8.ValidString(v.string) {
				b.scanErr = fmt.Errorf("[dml] 1649533490498 Column Index %d at position %d contains invalid UTF-8 characters", b.index, b.Count)
			} else {
				ptr.Data = v.string
				ptr.Valid = true
			}
		case 'y':
			if b.CheckValidUTF8 && !utf8.Valid(v.byte) {
				b.scanErr = fmt.Errorf("[dml] 1649533498261 Column Index %d at position %d contains invalid UTF-8 characters", b.index, b.Count)
			} else {
				ptr.Data = string(v.byte)
				ptr.Valid = v.byte != nil
			}
		case 'n':
			ptr.Data = ""
			ptr.Valid = false
		default:
			b.scanErr = fmt.Errorf("[dml] 1649533502539 Column %q does not support field type: %q", b.Column(), v.field)
		}
	}

	return b
}

// Time reads a time.Time value and appends it to the arguments slice or assigns
// the time.Time value stored in sql.RawBytes to the pointer. See the
// documentation for function Scan. It supports all MySQL/MariaDB date/time types.
func (b *ColumnMap) Time(ptr *time.Time) *ColumnMap {
	if b.shouldCollectArgs() {
		if ptr == nil {
			b.args = append(b.args, internalNULLNIL{})
		} else {
			b.args = append(b.args, *ptr)
		}
		return b
	}
	if b.scanErr == nil {
		switch v := b.scanCol[b.index]; v.field {
		case 't':
			*ptr = v.time
		case 'y':
			if len(v.byte) == 0 {
				b.scanErr = fmt.Errorf("[dml] 1649533511414 Column %q Time cannot be empty.", b.Column())
			} else {
				var nt null.Time
				nt, b.scanErr = null.ParseDateTime(string(v.byte), time.UTC) // time.Location can be merged into ColumnMap but then change NullTime method receiver.
				*ptr = nt.Time
				if b.scanErr != nil {
					b.scanErr = fmt.Errorf("[dml] 1649533072494 Column %q with error: %w", b.Column(), b.scanErr)
				}
			}
		default:
			b.scanErr = fmt.Errorf("[dml] Column %q does not support field type: %q", b.Column(), v.field)
		}
	}
	return b
}

// NullTime reads a time value and appends it to the arguments slice or assigns
// the NullTime value stored in sql.RawBytes to the pointer. See the
// documentation for function Scan.
func (b *ColumnMap) NullTime(ptr *null.Time) *ColumnMap {
	if b.shouldCollectArgs() {
		if ptr == nil {
			b.args = append(b.args, internalNULLNIL{})
		} else {
			b.args = append(b.args, *ptr)
		}
		return b
	}
	if b.scanErr == nil {
		switch v := b.scanCol[b.index]; v.field {
		case 't':
			ptr.Time = v.time
			ptr.Valid = true
		case 's':
			if err := ptr.Scan(v.string); err != nil {
				b.scanErr = fmt.Errorf("[dml] 1649533520040 ColumnMap NullTime: Invalid time string: %q with error %s", v.string, err)
			}
		case 'y':
			if err := ptr.Scan(v.byte); err != nil {
				b.scanErr = fmt.Errorf("[dml] 1649533526774 ColumnMap NullTime: Invalid time string: %q with error %s", v.byte, err)
			}
		case 'n':
			ptr.Time = time.Time{}
			ptr.Valid = false
		default:
			b.scanErr = fmt.Errorf("[dml] Column %q does not support field type: %q", b.Column(), v.field)
		}
	}
	return b
}

func (b *ColumnMap) addSlice(fnName string, slice any) *ColumnMap {
	if b.shouldCollectArgs() && b.scanErr == nil && b.Mode() == ColumnMapCollectionReadSet {
		b.args = append(b.args, slice)
	} else {
		b.scanErr = fmt.Errorf("[dml] 1649533329737 ColumnMap.%s does only support mode ColumnMapCollectionReadSet", fnName)
	}
	return b
}

func (b *ColumnMap) Uint64s(values ...uint64) *ColumnMap {
	return b.addSlice("Uint64s", values)
}

func (b *ColumnMap) Uint32s(values ...uint32) *ColumnMap {
	return b.addSlice("Uint32s", values)
}

func (b *ColumnMap) Uint16s(values ...uint16) *ColumnMap {
	return b.addSlice("Uint16s", values)
}

func (b *ColumnMap) Uint8s(values ...uint8) *ColumnMap {
	return b.addSlice("Uint8s", values)
}

func (b *ColumnMap) Int64s(values ...int64) *ColumnMap {
	return b.addSlice("Int64s", values)
}

func (b *ColumnMap) Int32s(values ...int32) *ColumnMap {
	return b.addSlice("Int32s", values)
}

func (b *ColumnMap) Int16s(values ...int16) *ColumnMap {
	return b.addSlice("Int16s", values)
}

func (b *ColumnMap) Int8s(values ...int8) *ColumnMap {
	return b.addSlice("Int8s", values)
}

func (b *ColumnMap) Strings(values ...string) *ColumnMap {
	return b.addSlice("Strings", values)
}

func (b *ColumnMap) Times(values ...time.Time) *ColumnMap {
	return b.addSlice("Times", values)
}

func (b *ColumnMap) NullStrings(values ...null.String) *ColumnMap {
	return b.addSlice("NullStrings", values)
}
