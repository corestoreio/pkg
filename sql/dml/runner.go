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
	"context"
	"database/sql"
	"encoding"
	"fmt"
	"io"
	"strconv"
	"time"
	"unicode/utf8"

	"github.com/corestoreio/errors"
	"github.com/corestoreio/log"
	"github.com/corestoreio/pkg/util/byteconv"
)

// Exec executes the statement represented by the QueryBuilder. It returns the
// raw database/sql Result or an error if there was one. Regarding
// LastInsertID(): If you insert multiple rows using a single INSERT statement,
// LAST_INSERT_ID() returns the value generated for the first inserted row only.
// The reason for this is to make it possible to reproduce easily the same
// INSERT statement against some other server.
// `db` can be either a *sql.DB (connection pool), a *sql.Conn (a single
// dedicated database session) or a *sql.Tx (an in-progress database
// transaction).
func Exec(ctx context.Context, db Execer, b QueryBuilder) (sql.Result, error) {
	sqlStr, args, err := b.ToSQL()
	if err != nil {
		return nil, errors.WithStack(err)
	}
	result, err := db.ExecContext(ctx, sqlStr, args...)
	return result, errors.Wrapf(err, "[dml] Exec.ExecContext with query %q", sqlStr)
}

// Prepare prepares a SQL statement. Sets IsInterpolate to false.
// `db` can be either a *sql.DB (connection pool), a *sql.Conn (a single
// dedicated database session) or a *sql.Tx (an in-progress database
// transaction).
func Prepare(ctx context.Context, db Preparer, b QueryBuilder) (*sql.Stmt, error) {
	sqlStr, _, err := b.ToSQL()
	if err != nil {
		return nil, errors.WithStack(err)
	}
	stmt, err := db.PrepareContext(ctx, sqlStr)
	return stmt, errors.Wrapf(err, "[dml] Prepare.PrepareContext with query %q", sqlStr)
}

// Query executes a query and returns many rows.
// `db` can be either a *sql.DB (connection pool), a *sql.Conn (a single
// dedicated database session) or a *sql.Tx (an in-progress database
// transaction).
func Query(ctx context.Context, db Querier, b QueryBuilder) (*sql.Rows, error) {
	sqlStr, args, err := b.ToSQL()
	if err != nil {
		return nil, errors.WithStack(err)
	}
	rows, err := db.QueryContext(ctx, sqlStr, args...)
	return rows, errors.Wrapf(err, "[dml] Query.QueryContext with query %q", sqlStr)
}

// Load loads data from a query into `s`. Load supports loading of up to n-rows.
// Load checks if a type implements io.Closer interface. `db` can be either a
// *sql.DB (connection pool), a *sql.Conn (a single dedicated database session)
// or a *sql.Tx (an in-progress database transaction).
//
// If ColumnMapper `s` implements io.Closer, the Close() function gets called in
// the defer function after rows.Close. `s.Close` function allows to implement
// features like unlocking a mutex or updating internal structures or closing a
// connection.
func Load(ctx context.Context, db Querier, b QueryBuilder, s ColumnMapper) (rowCount uint64, err error) {
	sqlStr, args, err := b.ToSQL()
	if err != nil {
		return 0, errors.WithStack(err)
	}
	rows, err := db.QueryContext(ctx, sqlStr, args...)
	rowCount, err = load(rows, err, s, new(ColumnMap))
	if err != nil {
		return 0, errors.Wrapf(err, "[dml] Load.QueryContext with query %q", sqlStr)
	}
	return rowCount, nil
}

func load(r *sql.Rows, errIn error, s ColumnMapper, cm *ColumnMap) (rowCount uint64, err error) {
	if errIn != nil {
		return 0, errors.WithStack(errIn)
	}
	defer func() {
		// Not testable with the sqlmock package :-(
		if err2 := r.Close(); err2 != nil && err == nil {
			err = errors.Wrap(err2, "[dml] Load.QueryContext.Rows.Close")
		}
		if rc, ok := s.(io.Closer); ok {
			if err2 := rc.Close(); err2 != nil && err == nil {
				err = errors.Wrap(err2, "[dml] Load.QueryContext.ColumnMapper.Close")
			}
		}
	}()

	for r.Next() {
		if err = cm.Scan(r); err != nil {
			return 0, errors.WithStack(err)
		}
		if err = s.MapColumns(cm); err != nil {
			return 0, errors.WithStack(err)
		}
	}
	if err = r.Err(); err != nil {
		return 0, errors.WithStack(err)
	}
	if cm.HasRows {
		cm.Count++ // because first row is zero but we want the actual row number
	}
	return cm.Count, err
}

// ColumnMapper allows a type to load data from database query into its fields
// or return the fields values as arguments for a query. It's used in the
// rows.Next() for-loop. A ColumnMapper is usually a single record/row or in
// case of a slice a complete query result.
type ColumnMapper interface {
	// RowScan implementation must use function `Scan` to scan the values of the
	// query into its own type. See database/sql package for examples.
	MapColumns(rc *ColumnMap) error
}

// Maybe add the following functions to ColumnMapper. Mostly useful
// when dealing with INSERT statements.
//FieldCount() int
//Length() int

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
	Args Arguments // in case we collect arguments

	// initialized gets set to true after the first call to Scan to initialize
	// the internal slices.
	initialized bool
	// CheckValidUTF8 if enabled checks if strings contains valid UTF-8 characters.
	CheckValidUTF8 bool
	// HasRows set to true if at least one row has been found.
	HasRows bool
	// Count increments on call to Scan.
	Count    uint64
	scanArgs []interface{} // could be a sync.Pool but check it in benchmarks.
	scanCol  []scannedColumn
	// Columns contains the names of the column returned from the query. One
	// should only read from the slice. Never modify it.
	columns    []string
	columnsLen int
	// scanErr is a delayed error and also used to avoid `if err != nil` in
	// generated code. This reduces the boiler plate code a lot! A trade off
	// between chainable API and too verbose error checking.
	scanErr error
	index   int
}

func newColumnMap(args Arguments, columns ...string) *ColumnMap {
	cm := &ColumnMap{Args: args}
	cm.setColumns(columns...)
	return cm
}

func (b *ColumnMap) setColumns(cols ...string) {
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
		return ColumnMapScan // assign the column values from the DB to the structs and create new structs in a slice.
	}

	// case b.Args != nil
	switch b.columnsLen {
	case 0:
		m = ColumnMapEntityReadAll // Entity: read all mode; Collection jump into loop and pass on to Entity
	case 1:
		m = ColumnMapCollectionReadSet // request certain column values as a slice.
	default:
		m = ColumnMapEntityReadSet // Entity: calls the for cm.Next loop; Collection jump into loop and pass on to Entity
	}
	return m
}

// scannedColumn represents an intermediate type (or DTO) to scan the
// driver.Values into. It supports the private types textRows and binaryRows in
// go-sql-driver/mysql. TextRows gets used during a normal query to read its
// result set and binaryRows gets used when a prepared statement gets executed
// and returns a result set. TextRows contains only byte slices whereas
// binaryRows contains already decoded types as defined in driver.Value. Avoids
// the reflection soup in database/sql.convertAssign.
// The supported data types depend on the this function:
// github.com/go-sql-driver/mysql/packets.go:1133 `func (rows *binaryRows) readRow(dest []driver.Value) error`
// Then all the type functions (Int,String,Uint8, etc) in type ColumnMap can
// support all of the MySQL protocol field types. Hence we can support
// fieldTypeGeometry, fieldTypeJSON with custom decoders.
type scannedColumn struct {
	field   byte      // i,f,b,y,s,t, n == null; nothing equals null of nil/empty
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

func (s *scannedColumn) Scan(src interface{}) (err error) {
	switch val := src.(type) {
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
	case []byte:
		s.field = 'y'
		s.byte = val
	case string:
		s.field = 's'
		s.string = val
	case time.Time:
		s.field = 't'
		s.time = val
	case nil:
		s.field = 'n'
	default:
		err = errors.NotSupported.Newf("[dml] ColumnMap.Scan does not yet support type %T with value: %#v", val, val)
	}
	return err
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
			return errors.WithStack(err)
		}
		b.setColumns(cols...)
		b.scanCol = make([]scannedColumn, b.columnsLen)
		b.scanArgs = make([]interface{}, b.columnsLen)
		for i := 0; i < b.columnsLen; i++ {
			b.scanArgs[i] = &b.scanCol[i]
		}
		b.initialized = true
		b.Count = 0
		b.HasRows = true
	} else {
		b.Count++
	}
	if err := r.Scan(b.scanArgs...); err != nil {
		return errors.WithStack(err)
	}
	return nil
}

// Err returns the delayed error from one of the scans and parsings. Function is
// idempotent.
func (b *ColumnMap) Err() error {
	return b.scanErr
}

// Column returns the current column name after calling `Next`.
func (b *ColumnMap) Column() string {
	return b.columns[b.index]
}

// Next moves the internal index to the next position. It may return false if
// during RawBytes scanning an error has occurred.
func (b *ColumnMap) Next() bool {
	b.index++
	ok := b.index < b.columnsLen && b.scanErr == nil
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
	if b.Args != nil {
		if ptr == nil {
			b.Args = b.Args.Null()
		} else {
			b.Args = b.Args.Bool(*ptr)
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
			*ptr, b.scanErr = byteconv.ParseBool(v.byte)
			if b.scanErr != nil {
				b.scanErr = errors.BadEncoding.New(b.scanErr, "[dml] Column %q", b.Column())
			}
		case 's':
			*ptr, b.scanErr = strconv.ParseBool(v.string)
			if b.scanErr != nil {
				b.scanErr = errors.BadEncoding.New(b.scanErr, "[dml] Column %q", b.Column())
			}
		default:
			b.scanErr = errors.NotSupported.Newf("[dml] Column %q does not support field type: %q", b.Column(), v.field)
		}
	}
	return b
}

// NullBool reads a bool value and appends it to the arguments slice or assigns the
// bool value stored in sql.RawBytes to the pointer. See the documentation for
// function Scan.
func (b *ColumnMap) NullBool(ptr *NullBool) *ColumnMap {
	if b.Args != nil {
		if ptr == nil {
			b.Args = b.Args.Null()
		} else {
			b.Args = b.Args.NullBool(*ptr)
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
			ptr.NullBool, b.scanErr = byteconv.ParseNullBool(v.byte)
			if b.scanErr != nil {
				b.scanErr = errors.BadEncoding.New(b.scanErr, "[dml] Column %q", b.Column())
			}
		case 's':
			if v.string != "" {
				ptr.Bool, b.scanErr = strconv.ParseBool(v.string)
				if b.scanErr != nil {
					b.scanErr = errors.BadEncoding.New(b.scanErr, "[dml] Column %q", b.Column())
				}
				ptr.Valid = b.scanErr == nil
			} else {
				ptr.Bool = false
				ptr.Valid = false
			}
		default:
			b.scanErr = errors.NotSupported.Newf("[dml] Column %q does not support field type: %q", b.Column(), v.field)
		}
	}
	return b
}

// Int reads an int value and appends it to the arguments slice or assigns the
// int value stored in sql.RawBytes to the pointer. See the documentation for
// function Scan.
func (b *ColumnMap) Int(ptr *int) *ColumnMap {
	if b.Args != nil {
		if ptr == nil {
			b.Args = b.Args.Null()
		} else {
			b.Args = b.Args.Int(*ptr)
		}
		return b
	}
	if b.scanErr == nil {
		switch v := b.scanCol[b.index]; v.field {
		case 'i':
			*ptr = int(v.int64)
		case 'y':
			var i64 int64
			i64, b.scanErr = byteconv.ParseInt(v.byte)
			*ptr = int(i64)
			if b.scanErr != nil {
				b.scanErr = errors.BadEncoding.New(b.scanErr, "[dml] Column %q", b.Column())
			}
		default:
			b.scanErr = errors.NotSupported.Newf("[dml] Column %q does not support field type: %q", b.Column(), v.field)
		}
	}
	return b
}

// Int64 reads a int64 value and appends it to the arguments slice or assigns
// the int64 value stored in sql.RawBytes to the pointer. See the documentation
// for function Scan.
func (b *ColumnMap) Int64(ptr *int64) *ColumnMap {
	if b.Args != nil {
		if ptr == nil {
			b.Args = b.Args.Null()
		} else {
			b.Args = b.Args.Int64(*ptr)
		}
		return b
	}
	if b.scanErr == nil {
		switch v := b.scanCol[b.index]; v.field {
		case 'i':
			*ptr = v.int64
		case 'y':
			*ptr, b.scanErr = byteconv.ParseInt(v.byte)
			if b.scanErr != nil {
				b.scanErr = errors.BadEncoding.New(b.scanErr, "[dml] Column %q", b.Column())
			}
		default:
			b.scanErr = errors.NotSupported.Newf("[dml] Column %q does not support field type: %q", b.Column(), v.field)
		}
	}
	return b
}

// NullInt64 reads an int64 value and appends it to the arguments slice or
// assigns the int64 value stored in sql.RawBytes to the pointer. See the
// documentation for function Scan.
func (b *ColumnMap) NullInt64(ptr *NullInt64) *ColumnMap {
	if b.Args != nil {
		if ptr == nil {
			b.Args = b.Args.Null()
		} else {
			b.Args = b.Args.NullInt64(*ptr)
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
			ptr.NullInt64, b.scanErr = byteconv.ParseNullInt64(v.byte)
			if b.scanErr != nil {
				b.scanErr = errors.BadEncoding.New(b.scanErr, "[dml] Column %q", b.Column())
			}
		default:
			b.scanErr = errors.NotSupported.Newf("[dml] Column %q does not support field type: %q", b.Column(), v.field)
		}
	}
	return b
}

// Float64 reads a float64 value and appends it to the arguments slice or
// assigns the float64 value stored in sql.RawBytes to the pointer. See the
// documentation for function Scan.
func (b *ColumnMap) Float64(ptr *float64) *ColumnMap {
	if b.Args != nil {
		if ptr == nil {
			b.Args = b.Args.Null()
		} else {
			b.Args = b.Args.Float64(*ptr)
		}
		return b
	}
	if b.scanErr == nil {
		switch v := b.scanCol[b.index]; v.field {
		case 'f':
			*ptr = v.float64
		case 'y':
			*ptr, b.scanErr = byteconv.ParseFloat(v.byte)
			if b.scanErr != nil {
				b.scanErr = errors.BadEncoding.New(b.scanErr, "[dml] Column %q", b.Column())
			}
		default:
			b.scanErr = errors.NotSupported.Newf("[dml] Column %q does not support field type: %q", b.Column(), v.field)
		}
	}
	return b
}

// Decimal reads a Decimal value and appends it to the arguments slice or
// assigns the numeric value stored in sql.RawBytes to the pointer. See the
// documentation for function Scan.
func (b *ColumnMap) Decimal(ptr *Decimal) *ColumnMap {
	if b.Args != nil {
		if v := ptr.String(); ptr == nil || v == sqlStrNullUC {
			b.Args = b.Args.Null()
		} else {
			b.Args = b.Args.String(v)
		}
		return b
	}
	if b.scanErr == nil {
		switch v := b.scanCol[b.index]; v.field {
		case 'f':
			*ptr, b.scanErr = MakeDecimalFloat64(v.float64)
		case 'y':
			*ptr, b.scanErr = MakeDecimalBytes(v.byte)
		case 's':
			*ptr, b.scanErr = MakeDecimalBytes([]byte(v.string)) // mostly used for testing
		case 'n':
			ptr.Valid = false
		default:
			b.scanErr = errors.NotSupported.Newf("[dml] Column %q does not support field type: %q", b.Column(), v.field)
		}
	}
	return b
}

// NullFloat64 reads a float64 value and appends it to the arguments slice or
// assigns the float64 value stored in sql.RawBytes to the pointer. See the
// documentation for function Scan.
func (b *ColumnMap) NullFloat64(ptr *NullFloat64) *ColumnMap {
	if b.Args != nil {
		if ptr == nil {
			b.Args = b.Args.Null()
		} else {
			b.Args = b.Args.NullFloat64(*ptr)
		}
		return b
	}
	if b.scanErr == nil {
		switch v := b.scanCol[b.index]; v.field {
		case 'f':
			ptr.Float64 = v.float64
			ptr.Valid = true
		case 'y':
			ptr.NullFloat64, b.scanErr = byteconv.ParseNullFloat64(v.byte)
			if b.scanErr != nil {
				b.scanErr = errors.BadEncoding.New(b.scanErr, "[dml] Column %q", b.Column())
			}
		case 'n':
			ptr.Valid = false
		default:
			b.scanErr = errors.NotSupported.Newf("[dml] Column %q does not support field type: %q", b.Column(), v.field)
		}
	}
	return b
}

// Uint reads an uint value and appends it to the arguments slice or assigns the
// uint value stored in sql.RawBytes to the pointer. See the documentation for
// function Scan.
func (b *ColumnMap) Uint(ptr *uint) *ColumnMap {
	if b.Args != nil {
		if ptr == nil {
			b.Args = b.Args.Null()
		} else {
			b.Args = b.Args.Uint(*ptr)
		}
		return b
	}
	if b.scanErr == nil {
		switch v := b.scanCol[b.index]; v.field {
		case 'i':
			*ptr = uint(v.int64)
		case 'y':
			var u64 uint64
			u64, b.scanErr = byteconv.ParseUint(v.byte, 10, strconv.IntSize)
			*ptr = uint(u64)
			if b.scanErr != nil {
				b.scanErr = errors.BadEncoding.New(b.scanErr, "[dml] Column %q", b.Column())
			}
		default:
			b.scanErr = errors.NotSupported.Newf("[dml] Column %q does not support field type: %q", b.Column(), v.field)
		}
	}
	return b
}

// Uint8 reads an uint8 value and appends it to the arguments slice or assigns
// the uint8 value stored in sql.RawBytes to the pointer. See the documentation
// for function Scan.
func (b *ColumnMap) Uint8(ptr *uint8) *ColumnMap {
	if b.Args != nil {
		if ptr == nil {
			b.Args = b.Args.Null()
		} else {
			b.Args = b.Args.Uint(uint(*ptr))
		}
		return b
	}
	if b.scanErr == nil {
		switch v := b.scanCol[b.index]; v.field {
		case 'i':
			*ptr = uint8(v.int64)
		case 'y':
			var u64 uint64
			u64, b.scanErr = byteconv.ParseUint(v.byte, 10, 8)
			*ptr = uint8(u64)
			if b.scanErr != nil {
				b.scanErr = errors.BadEncoding.New(b.scanErr, "[dml] Column %q", b.Column())
			}
		default:
			b.scanErr = errors.NotSupported.Newf("[dml] Column %q does not support field type: %q", b.Column(), v.field)
		}
	}
	return b
}

// Uint16 reads an uint16 value and appends it to the arguments slice or assigns
// the uint16 value stored in sql.RawBytes to the pointer. See the documentation
// for function Scan.
func (b *ColumnMap) Uint16(ptr *uint16) *ColumnMap {
	if b.Args != nil {
		if ptr == nil {
			b.Args = b.Args.Null()
		} else {
			b.Args = b.Args.Uint(uint(*ptr))
		}
		return b
	}
	if b.scanErr == nil {
		switch v := b.scanCol[b.index]; v.field {
		case 'i':
			*ptr = uint16(v.int64)
		case 'y':
			var u64 uint64
			u64, b.scanErr = byteconv.ParseUint(v.byte, 10, 16)
			*ptr = uint16(u64)
			if b.scanErr != nil {
				b.scanErr = errors.BadEncoding.New(b.scanErr, "[dml] Column %q", b.Column())
			}
		default:
			b.scanErr = errors.NotSupported.Newf("[dml] Column %q does not support field type: %q", b.Column(), v.field)
		}
	}
	return b
}

// Uint32 reads an uint32 value and appends it to the arguments slice or assigns
// the uint32 value stored in sql.RawBytes to the pointer. See the documentation
// for function Scan.
func (b *ColumnMap) Uint32(ptr *uint32) *ColumnMap {
	if b.Args != nil {
		if ptr == nil {
			b.Args = b.Args.Null()
		} else {
			b.Args = b.Args.Uint(uint(*ptr))
		}
		return b
	}
	if b.scanErr == nil {
		switch v := b.scanCol[b.index]; v.field {
		case 'i':
			*ptr = uint32(v.int64)
		case 'y':
			var u64 uint64
			u64, b.scanErr = byteconv.ParseUint(v.byte, 10, 32)
			*ptr = uint32(u64)
			if b.scanErr != nil {
				b.scanErr = errors.BadEncoding.New(b.scanErr, "[dml] Column %q", b.Column())
			}
		default:
			b.scanErr = errors.NotSupported.Newf("[dml] Column %q does not support field type: %q", b.Column(), v.field)
		}
	}
	return b
}

// Uint64 reads an uint64 value and appends it to the arguments slice or assigns
// the uint64 value stored in sql.RawBytes to the pointer. See the documentation
// for function Scan.
func (b *ColumnMap) Uint64(ptr *uint64) *ColumnMap {
	if b.Args != nil {
		if ptr == nil {
			b.Args = b.Args.Null()
		} else {
			b.Args = b.Args.Uint64(*ptr)
		}
		return b
	}
	if b.scanErr == nil {
		switch v := b.scanCol[b.index]; v.field {
		case 'i':
			*ptr = uint64(v.int64)
		case 'y':
			*ptr, b.scanErr = byteconv.ParseUint(v.byte, 10, strconv.IntSize)
			if b.scanErr != nil {
				b.scanErr = errors.BadEncoding.New(b.scanErr, "[dml] Column %q", b.Column())
			}
		default:
			b.scanErr = errors.NotSupported.Newf("[dml] Column %q does not support field type: %q", b.Column(), v.field)
		}
	}
	return b
}

// Debug writes the column names with their values into `w`. The output format
// might change.
func (b *ColumnMap) Debug(w io.Writer) (err error) {
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
		} else {
			if _, err = fmt.Fprintf(w, ": %q", b); err != nil {
				return errors.WithStack(err)
			}
		}
	}
	return nil
}

// Byte reads a []byte value and appends it to the arguments slice or assigns
// the []byte value stored in sql.RawBytes to the pointer. See the documentation
// for function Scan.
func (b *ColumnMap) Byte(ptr *[]byte) *ColumnMap {
	if b.Args != nil {
		if ptr == nil {
			b.Args = b.Args.Null()
		} else {
			b.Args = b.Args.Bytes(*ptr)
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
			b.scanErr = errors.NotSupported.Newf("[dml] Column %q does not support field type: %q", b.Column(), v.field)
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
	if b.Args != nil {
		var data []byte
		data, b.scanErr = enc.MarshalText()
		if b.CheckValidUTF8 && !utf8.Valid(data) {
			b.scanErr = errors.NotValid.Newf("[dml] Column Index %d at position %d contains invalid UTF-8 characters", b.index, b.Count)
		} else {
			b.Args = b.Args.Bytes(data)
		}
		return b
	}

	switch v := b.scanCol[b.index]; v.field {
	case 'y', 'n':
		if b.CheckValidUTF8 && !utf8.Valid(v.byte) {
			b.scanErr = errors.NotValid.Newf("[dml] Column %q contains invalid UTF-8 characters", b.Column(), b.Count)
		} else if b.scanErr = enc.UnmarshalText(v.byte); b.scanErr != nil {
			b.scanErr = errors.BadEncoding.New(b.scanErr, "[dml] Column %q", b.Column())
		}
	default:
		b.scanErr = errors.NotSupported.Newf("[dml] Column %q does not support field type: %q", b.Column(), v.field)
	}
	return b
}

// Text allows to encode an object to its binary representation when arguments are
// requested and to decode a byte slice into its object when data is retrieved
// from the server. Use this function for GOB, Protocol Buffers, etc formats.
func (b *ColumnMap) Binary(enc interface {
	encoding.BinaryMarshaler
	encoding.BinaryUnmarshaler
}) *ColumnMap {
	if b.scanErr != nil {
		return b
	}
	if b.Args != nil {
		var data []byte
		data, b.scanErr = enc.MarshalBinary()
		b.Args = b.Args.Bytes(data)
		return b
	}

	switch v := b.scanCol[b.index]; v.field {
	case 'y', 'n':
		if b.scanErr = enc.UnmarshalBinary(v.byte); b.scanErr != nil {
			b.scanErr = errors.BadEncoding.New(b.scanErr, "[dml] Column %q", b.Column())
		}
	default:
		b.scanErr = errors.NotSupported.Newf("[dml] Column %q does not support field type: %q", b.Column(), v.field)
	}
	return b
}

// String reads a string value and appends it to the arguments slice or assigns
// the string value stored in sql.RawBytes to the pointer. See the documentation
// for function Scan.
func (b *ColumnMap) String(ptr *string) *ColumnMap {
	if b.Args != nil {
		if ptr == nil {
			b.Args = b.Args.Null()
		} else {
			b.Args = b.Args.String(*ptr)
		}
		return b
	}

	if b.scanErr == nil {
		switch v := b.scanCol[b.index]; v.field {
		case 's':
			if b.CheckValidUTF8 && !utf8.ValidString(v.string) {
				b.scanErr = errors.NotValid.Newf("[dml] Column Index %d at position %d contains invalid UTF-8 characters", b.index, b.Count)
			} else {
				*ptr = v.string
			}
		case 'y':
			if b.CheckValidUTF8 && !utf8.Valid(v.byte) {
				b.scanErr = errors.NotValid.Newf("[dml] Column Index %d at position %d contains invalid UTF-8 characters", b.index, b.Count)
			} else {
				*ptr = string(v.byte)
			}
		default:
			b.scanErr = errors.NotSupported.Newf("[dml] Column %q does not support field type: %q", b.Column(), v.field)
		}
	}
	return b
}

// NullString reads a string value and appends it to the arguments slice or
// assigns the string value stored in sql.RawBytes to the pointer. See the
// documentation for function Scan.
func (b *ColumnMap) NullString(ptr *NullString) *ColumnMap {
	if b.Args != nil {
		if ptr == nil {
			b.Args = b.Args.Null()
		} else {
			b.Args = b.Args.NullString(*ptr)
		}
		return b
	}

	if b.scanErr == nil {
		switch v := b.scanCol[b.index]; v.field {
		case 's':
			if b.CheckValidUTF8 && !utf8.ValidString(v.string) {
				b.scanErr = errors.NotValid.Newf("[dml] Column Index %d at position %d contains invalid UTF-8 characters", b.index, b.Count)
			} else {
				ptr.String = v.string
				ptr.Valid = true
			}
		case 'y':
			if b.CheckValidUTF8 && !utf8.Valid(v.byte) {
				b.scanErr = errors.NotValid.Newf("[dml] Column Index %d at position %d contains invalid UTF-8 characters", b.index, b.Count)
			} else {
				ptr.String = string(v.byte)
				ptr.Valid = v.byte != nil
			}
		case 'n':
			ptr.String = ""
			ptr.Valid = false
		default:
			b.scanErr = errors.NotSupported.Newf("[dml] Column %q does not support field type: %q", b.Column(), v.field)
		}
	}

	return b
}

// Time reads a time.Time value and appends it to the arguments slice or assigns
// the time.Time value stored in sql.RawBytes to the pointer. See the
// documentation for function Scan. It supports all MySQL/MariaDB date/time types.
func (b *ColumnMap) Time(ptr *time.Time) *ColumnMap {
	if b.Args != nil {
		if ptr == nil {
			b.Args = b.Args.Null()
		} else {
			b.Args = b.Args.Time(*ptr)
		}
		return b
	}
	if b.scanErr == nil {
		switch v := b.scanCol[b.index]; v.field {
		case 't':
			*ptr = v.time
		case 'y':
			*ptr, b.scanErr = parseDateTime(string(v.byte), time.UTC) // time.Location can be merged into ColumnMap but then change NullTime method receiver.
			if b.scanErr != nil {
				b.scanErr = errors.BadEncoding.New(b.scanErr, "[dml] Column %q", b.Column())
			}
		default:
			b.scanErr = errors.NotSupported.Newf("[dml] Column %q does not support field type: %q", b.Column(), v.field)
		}
	}
	return b
}

// NullTime reads a time value and appends it to the arguments slice or assigns
// the NullTime value stored in sql.RawBytes to the pointer. See the
// documentation for function Scan.
func (b *ColumnMap) NullTime(ptr *NullTime) *ColumnMap {
	if b.Args != nil {
		if ptr == nil {
			b.Args = b.Args.Null()
		} else {
			b.Args = b.Args.NullTime(*ptr)
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
				b.scanErr = errors.NotValid.Newf("[dml] ColumnMap NullTime: Invalid time string: %q with error %s", v.string, err)
			}
		case 'y':
			if err := ptr.Scan(v.byte); err != nil {
				b.scanErr = errors.NotValid.Newf("[dml] ColumnMap NullTime: Invalid time string: %q with error %s", v.byte, err)
			}
		case 'n':
			ptr.Time = time.Time{}
			ptr.Valid = false
		default:
			b.scanErr = errors.NotSupported.Newf("[dml] Column %q does not support field type: %q", b.Column(), v.field)
		}
	}
	return b
}

// StmtBase wraps a *sql.Stmt (a prepared statement) with a specific SQL query.
// To create a StmtBase call the Prepare function of type Select. StmtBase is
// not safe for concurrent use, despite the underlying *sql.Stmt is. Don't
// forget to call Close!
type StmtBase struct {
	builderCommon
	stmt *sql.Stmt
}

// Close closes the underlying prepared statement.
func (st *StmtBase) Close() error { return st.stmt.Close() }

func (st *StmtBase) resetArgs() {
	st.argsArgs = st.argsArgs[:0]
	st.argsRaw = st.argsRaw[:0]
	st.argsRecords = st.argsRecords[:0]
}

func (st *StmtBase) withArgs(args []interface{}) {
	st.resetArgs()
	st.argsRaw = args
	st.isWithInterfaces = true
}

func (st *StmtBase) withArguments(args Arguments) {
	st.resetArgs()
	st.argsArgs = args
	st.isWithInterfaces = false
}

// withRecords sets the records for the execution with Query or Exec. It
// internally resets previously applied arguments.
func (st *StmtBase) withRecords(records []QualifiedRecord) {
	st.resetArgs()
	st.argsRecords = records
	st.isWithInterfaces = false
}

// prepareArgs transforms mainly the Arguments into []interface{} but also
// appends the `args` from the Exec+ or Query+ function.
// All method receivers are not thread safe.
func (st *StmtBase) prepareArgs(args ...interface{}) error {
	if st.ärgErr != nil {
		return st.ärgErr
	}

	if !st.isWithInterfaces {
		st.argsRaw = st.argsRaw[:0]
	}

	argsArgs, err := st.convertRecordsToArguments()
	st.argsRaw = append(st.argsRaw, argsArgs.Interfaces()...)
	st.argsRaw = append(st.argsRaw, args...)
	return err
}

// Errors do not get logged in the next functions. Errors are getting handled.

// Exec supports both either the traditional way or passing arguments or
// in combination with the previously called WithArguments, WithRecords or
// WithArgs functions. If you want to call it multiple times with the same
// arguments, do not use the `args` variable, instead use the With+ functions.
// Calling any of the With+ function and additionally setting the `args`, will
// append the `args` at the end to the previously set or generated arguments.
// This function is not thread safe.
func (st *StmtBase) Exec(ctx context.Context, args ...interface{}) (sql.Result, error) {
	if err := st.prepareArgs(args...); err != nil {
		return nil, errors.WithStack(err)
	}
	if st.Log != nil && st.Log.IsDebug() {
		defer log.WhenDone(st.Log).Debug("Exec", log.Int("arg_len", len(st.argsRaw)))
	}
	return st.stmt.ExecContext(ctx, st.argsRaw...)
}

// Query traditional way, allocation heavy.
func (st *StmtBase) Query(ctx context.Context, args ...interface{}) (*sql.Rows, error) {
	if err := st.prepareArgs(args...); err != nil {
		return nil, errors.WithStack(err)
	}
	if st.Log != nil && st.Log.IsDebug() {
		defer log.WhenDone(st.Log).Debug("Query", log.Int("arg_len", len(st.argsRaw)))
	}
	return st.stmt.QueryContext(ctx, st.argsRaw...)
}

// QueryRow traditional way, allocation heavy.
func (st *StmtBase) QueryRow(ctx context.Context, args ...interface{}) *sql.Row {
	if err := st.prepareArgs(args...); err != nil {
		_ = err
		// Hmmm what should happen here?
	}
	if st.Log != nil && st.Log.IsDebug() {
		defer log.WhenDone(st.Log).Debug("QueryRow", log.Int("arg_len", len(st.argsRaw)))
	}
	return st.stmt.QueryRowContext(ctx, st.argsRaw...)
}

// Load loads data from a query into an object. You must set DB.QueryContext on
// the Select object or it just panics. Load can load a single row or n-rows.
func (st *StmtBase) Load(ctx context.Context, s ColumnMapper) (rowCount uint64, err error) {
	if st.Log != nil && st.Log.IsDebug() {
		defer log.WhenDone(st.Log).Debug("Load", log.Uint64("row_count", rowCount), log.String("object_type", fmt.Sprintf("%T", s)))
	}
	r, err := st.Query(ctx)
	rowCount, err = load(r, err, s, &st.ColumnMap)
	return rowCount, errors.WithStack(err)
}

// LoadInt64 executes the prepared statement and returns the value at an int64.
// It returns a NotFound error if the query returns nothing.
func (st *StmtBase) LoadInt64(ctx context.Context) (int64, error) {
	if st.Log != nil && st.Log.IsDebug() {
		defer log.WhenDone(st.Log).Debug("LoadInt64")
	}
	return loadInt64(st.Query(ctx))
}

// LoadInt64s executes the Select and returns the value as a slice of int64s.
func (st *StmtBase) LoadInt64s(ctx context.Context) (ret []int64, err error) {
	if st.Log != nil && st.Log.IsDebug() {
		// do not use fullSQL because we might log sensitive data
		defer log.WhenDone(st.Log).Debug("LoadInt64s", log.Int("row_count", len(ret)))
	}
	ret, err = loadInt64s(st.Query(ctx))
	// Do not simplify it because we need ret in the defer. we don't log errors
	// because they get handled.
	return ret, err
}

// More Load* functions can be added later
