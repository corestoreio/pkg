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
	"context"
	"database/sql"
	"fmt"
	"io"
	"math"
	"strconv"
	"unicode/utf8"

	"github.com/corestoreio/csfw/util/bufferpool"
	"github.com/corestoreio/csfw/util/byteconv"
	"github.com/corestoreio/errors"
)

// Scanner allows a type to load data from database query. It's used in the
// rows.Next() for-loop.
type Scanner interface {
	// RowScan implementation must use function `Scan` to scan the values of the
	// query into its own type. See database/sql package for examples.
	RowScan(*sql.Rows) error
}

// RowCloser allows to execute special functions after the scanning has
// happened. Should only be implemented in a custom type, if the interface
// Scanner has been implemented too. Not every type might need a RowCloser.
type RowCloser interface {
	// RowClose gets called at the very end and even after rows.Close. Allows
	// to implement special functions, like unlocking a mutex or updating
	// internal structures or resetting internal type containers.
	RowClose() error
}

// Exec executes the statement represented by the QueryBuilder. It returns the
// raw database/sql Result or an error if there was one. Regarding
// LastInsertID(): If you insert multiple rows using a single INSERT statement,
// LAST_INSERT_ID() returns the value generated for the first inserted row only.
// The reason for this is to make it possible to reproduce easily the same
// INSERT statement against some other server.
func Exec(ctx context.Context, db Execer, b QueryBuilder) (sql.Result, error) {
	sqlStr, args, err := b.ToSQL()
	if err != nil {
		return nil, errors.WithStack(err)
	}
	result, err := db.ExecContext(ctx, sqlStr, args...)
	return result, errors.Wrapf(err, "[dbr] Exec.ExecContext with query %q", sqlStr)
}

// Prepare prepares a SQL statement. Sets IsInterpolate to false.
func Prepare(ctx context.Context, db Preparer, b QueryBuilder) (*sql.Stmt, error) {
	var sqlStr string
	var err error
	if qb, ok := b.(queryBuilder); ok { // Interface upgrade
		sqlStr, _, err = toSQL(qb, _isNotInterpolate, _isPrepared)
		if err != nil {
			return nil, errors.WithStack(err)
		}
	} else {
		sqlStr, _, err = b.ToSQL()
		if err != nil {
			return nil, errors.WithStack(err)
		}
	}
	stmt, err := db.PrepareContext(ctx, sqlStr)
	return stmt, errors.Wrapf(err, "[dbr] Prepare.PrepareContext with query %q", sqlStr)
}

// Query executes a query and returns many rows.
func Query(ctx context.Context, db Querier, b QueryBuilder) (*sql.Rows, error) {
	sqlStr, args, err := b.ToSQL()
	if err != nil {
		return nil, errors.WithStack(err)
	}
	rows, err := db.QueryContext(ctx, sqlStr, args...)
	return rows, errors.Wrapf(err, "[dbr] Query.QueryContext with query %q", sqlStr)
}

// Load loads data from a query into `s`. Load supports loading of up to n-rows.
// Load checks if a type implements RowCloser interface.
func Load(ctx context.Context, db Querier, b QueryBuilder, s Scanner) (rowCount int64, err error) {
	sqlStr, args, err := b.ToSQL()
	if err != nil {
		return 0, errors.WithStack(err)
	}

	r, err := db.QueryContext(ctx, sqlStr, args...)
	if err != nil {
		return 0, errors.Wrapf(err, "[dbr] Load.QueryContext with query %q", sqlStr)
	}
	defer func() {
		// Not testable with the sqlmock package :-(
		if err2 := r.Close(); err2 != nil && err == nil {
			err = errors.Wrap(err2, "[dbr] Load.QueryContext.Rows.Close")
		}
		if rc, ok := s.(RowCloser); ok {
			if err2 := rc.RowClose(); err2 != nil && err == nil {
				err = errors.Wrap(err2, "[dbr] Load.QueryContext.Scanner.RowClose")
			}
		}
	}()

	for r.Next() {
		err = s.RowScan(r)
		if err != nil {
			return 0, errors.WithStack(err)
		}
		rowCount++
	}
	if err = r.Err(); err != nil {
		return rowCount, errors.WithStack(err)
	}
	return rowCount, err
}

// RowConvert represents the commonly used fields for each struct for a database
// table or a view to convert the sql.RawBytes into the desired final type. It
// scans a *sql.Rows into a *sql.RawBytes slice. The conversion into the desired
// final type can happen without allocating of memory. RowConvert should be used
// as a composite field in a database table struct.
type RowConvert struct {
	// Count increments on call to Scan.
	Count uint64
	// Columns contains the names of the column returned from the query.
	Columns []string
	// Aliased maps a `key` containing the alias name, used in the query, to the
	// `value`, the original snake case name used in the parent struct.
	Alias map[string]string
	// Initialized gets set to true after the first call to Scan to initialize
	// the internal slices.
	Initialized bool
	// CheckValidUTF8 if enabled checks if strings contains valid UTF-8 characters.
	CheckValidUTF8 bool
	scanArgs       []interface{}
	scanRaw        []*sql.RawBytes
	index          int
	current        []byte
}

// Scan calls rows.Scan and builds an internal stack of sql.RawBytes for further
// processing and type conversion.
//
// Each function for a specific type converts the underlying byte slice at the
// current set index (see function Index) to the appropriate type. You can call
// as many times as you want the specific functions. The underlying byte slice
// value is valid until the next call to rows.Next, rows,Scan or rows.Close. See
// the example for further usages.
//
// sqlRower relates to type *sql.Rows, but kept private to not confuse
// developers with another exported interface. The interface exists mainly for
// testing purposes.
func (b *RowConvert) Scan(r *sql.Rows) error {
	if !b.Initialized {
		var err error
		b.Columns, err = r.Columns()
		if err != nil {
			return errors.WithStack(err)
		}
		lc := len(b.Columns)
		b.scanRaw = make([]*sql.RawBytes, lc)
		b.scanArgs = make([]interface{}, lc)
		for i := range b.Columns {
			rb := new(sql.RawBytes)
			b.scanRaw[i] = rb
			b.scanArgs[i] = rb
		}
		b.Initialized = true
		b.Count = 0
	}
	if err := r.Scan(b.scanArgs...); err != nil {
		return errors.WithStack(err)
	}
	b.Count++
	return nil
}

// Index sets the current column index to read data from. You must call this
// first before calling any other type conversion function or they will return
// empty or NULL values.
func (b *RowConvert) Index(i int) *RowConvert {
	b.index = i
	b.current = *b.scanRaw[i]
	return b
}

// Bool see the documentation for function Scan.
func (b *RowConvert) Bool() (bool, error) {
	return byteconv.ParseBool(b.current)
}

// NullBool see the documentation for function Scan.
func (b *RowConvert) NullBool() (sql.NullBool, error) {
	return byteconv.ParseNullBool(b.current)
}

// Int see the documentation for function Scan.
func (b *RowConvert) Int() (int, error) {
	i, err := byteconv.ParseInt(b.current)
	if err != nil {
		return 0, err
	}
	if strconv.IntSize == 32 && (i < -math.MaxInt32 || i > math.MaxInt32) {
		return 0, rangeError("RowConvert.Int", string(b.current))
	}
	return int(i), nil
}

// Int64 see the documentation for function Scan.
func (b *RowConvert) Int64() (int64, error) {
	return byteconv.ParseInt(b.current)
}

// NullInt64 see the documentation for function Scan.
func (b *RowConvert) NullInt64() (sql.NullInt64, error) {
	return byteconv.ParseNullInt64(b.current)
}

// Float64 see the documentation for function Scan.
func (b *RowConvert) Float64() (float64, error) {
	return byteconv.ParseFloat(b.current)
}

// NullFloat64 see the documentation for function Scan.
func (b *RowConvert) NullFloat64() (sql.NullFloat64, error) {
	return byteconv.ParseNullFloat64(b.current)
}

// Uint see the documentation for function Scan.
func (b *RowConvert) Uint() (uint, error) {
	i, _, err := byteconv.ParseUintSQL(b.current, 10, strconv.IntSize)
	if err != nil {
		return 0, err
	}
	if strconv.IntSize == 32 && i > math.MaxUint32 {
		return 0, rangeError("RowConvert.Uint", string(b.current))
	}
	return uint(i), nil
}

// Uint8 see the documentation for function Scan.
func (b *RowConvert) Uint8() (uint8, error) {
	i, _, err := byteconv.ParseUintSQL(b.current, 10, 8)
	if err != nil {
		return 0, err
	}
	return uint8(i), nil
}

// Uint16 see the documentation for function Scan.
func (b *RowConvert) Uint16() (uint16, error) {
	i, _, err := byteconv.ParseUintSQL(b.current, 10, 16)
	if err != nil {
		return 0, err
	}
	return uint16(i), nil
}

// Uint32 see the documentation for function Scan.
func (b *RowConvert) Uint32() (uint32, error) {
	i, _, err := byteconv.ParseUintSQL(b.current, 10, 32)
	if err != nil {
		return 0, err
	}
	return uint32(i), nil
}

// Uint64 see the documentation for function Scan.
func (b *RowConvert) Uint64() (uint64, error) {
	i, _, err := byteconv.ParseUintSQL(b.current, 10, 64)
	return i, err
}

// String implements fmt.Stringer interface and returns the column names with
// their values. Mostly useful for debugging purposes. The output format might
// change.
func (b *RowConvert) String() string {
	buf := bufferpool.Get()
	defer bufferpool.Put(buf)
	for i, c := range b.Columns {
		if i > 0 {
			buf.WriteByte('\n')
		}
		buf.WriteString(c)
		b := *b.scanRaw[i]
		if b == nil {
			buf.WriteString(": <nil>")
		} else {
			fmt.Fprintf(buf, ": %q", string(b))
		}
	}
	return buf.String()
}

// Byte copies the value byte slice at index `idx` into a new slice. See the
// documentation for function Scan.
func (b *RowConvert) Byte() []byte {
	if b.current == nil {
		return nil
	}
	ret := make([]byte, len(b.current))
	copy(ret, b.current)
	return ret
}

// WriteTo implements interface io.WriterTo. It puts the underlying byte slice
// directly into w. The value is valid until the next call to rows.Next.
// See the documentation for function Scan.
func (b *RowConvert) WriteTo(w io.Writer) (n int64, err error) {
	var n2 int
	n2, err = w.Write(b.current)
	return int64(n2), errors.WithStack(err)
}

// Str see the documentation for function Scan.
func (b *RowConvert) Str() (string, error) {
	if b.CheckValidUTF8 && !utf8.Valid(b.current) {
		return "", errors.NewNotValidf("[dbr] Column Index %d at position %d contains invalid UTF-8 characters", b.index, b.Count)
	}
	return string(b.current), nil
}

// NullString see the documentation for function Scan.
func (b *RowConvert) NullString() (sql.NullString, error) {
	if b.CheckValidUTF8 && !utf8.Valid(b.current) {
		return sql.NullString{}, errors.NewNotValidf("[dbr] Column Index %d at position %d contains invalid UTF-8 characters", b.index, b.Count)
	}
	s := byteconv.ParseNullString(b.current)
	return s, nil
}

func rangeError(fn, str string) *strconv.NumError {
	return &strconv.NumError{Func: fn, Num: str, Err: strconv.ErrRange}
}
