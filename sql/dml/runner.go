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

package dml

import (
	"context"
	"database/sql"
	"fmt"
	"io"
	"math"
	"strconv"
	"time"
	"unicode/utf8"

	"github.com/corestoreio/csfw/util/byteconv"
	"github.com/corestoreio/errors"
	"github.com/corestoreio/log"
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
	rowCount, err = load(rows, err, s)
	if err != nil {
		return 0, errors.Wrapf(err, "[dml] Load.QueryContext with query %q", sqlStr)
	}
	return rowCount, nil
}

func load(r *sql.Rows, errIn error, s ColumnMapper) (rowCount uint64, err error) {
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

	rm := new(ColumnMap) // TODO(CyS) use sync.Pool
	for r.Next() {
		if err = rm.Scan(r); err != nil {
			return 0, errors.WithStack(err)
		}
		if err = s.MapColumns(rm); err != nil {
			return 0, errors.WithStack(err)
		}
	}
	if err = r.Err(); err != nil {
		return 0, errors.WithStack(err)
	}
	if rm.HasRows {
		rm.Count++ // because first row is zero but we want the actual row number
	}
	return rm.Count, err
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

	// HasRows set to true if at least one row has been found.
	HasRows bool
	// Count increments on call to Scan.
	Count uint64
	// Columns contains the names of the column returned from the query. One
	// should only read from the slice. Never modify it.
	columns    []string
	columnsLen int
	// initialized gets set to true after the first call to Scan to initialize
	// the internal slices.
	initialized bool
	// CheckValidUTF8 if enabled checks if strings contains valid UTF-8 characters.
	CheckValidUTF8 bool
	scanArgs       []interface{} // could be a sync.Pool but check it in benchmarks.
	scanRaw        []*sql.RawBytes
	scanErr        error // delayed error and also to avoid `if err != nil`
	index          int
	current        []byte
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
		m = ColumnMapCollectionReadSet // request certain column values as a slice. implemented in func condition.go:appendArgs.
	default:
		m = ColumnMapEntityReadSet // Entity: calls the for cm.Next loop; Collection jump into loop and pass on to Entity
	}
	return m
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
		b.scanRaw = make([]*sql.RawBytes, b.columnsLen)
		b.scanArgs = make([]interface{}, b.columnsLen)
		for i := 0; i < b.columnsLen; i++ {
			rb := new(sql.RawBytes)
			b.scanRaw[i] = rb
			b.scanArgs[i] = rb
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
	if ok && b.scanRaw != nil {
		b.current = *b.scanRaw[b.index]
	}
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
		*ptr, b.scanErr = byteconv.ParseBool(b.current)
		if b.scanErr != nil {
			b.scanErr = errors.Wrapf(b.scanErr, "[dml] Column %q", b.Column())
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
		var nv sql.NullBool
		nv, b.scanErr = byteconv.ParseNullBool(b.current)
		*ptr = NullBool{NullBool: nv}
		if b.scanErr != nil {
			b.scanErr = errors.Wrapf(b.scanErr, "[dml] Column %q", b.Column())
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
		var i64 int64
		i64, b.scanErr = byteconv.ParseInt(b.current)
		if b.scanErr == nil && strconv.IntSize == 32 && (i64 < -math.MaxInt32 || i64 > math.MaxInt32) { // hmm rethink that depending on goarch
			b.scanErr = rangeError("ColumnMap.Int", string(b.current))
		}
		if b.scanErr != nil {
			b.scanErr = errors.Wrapf(b.scanErr, "[dml] Column %q", b.Column())
		}
		*ptr = int(i64)
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
		*ptr, b.scanErr = byteconv.ParseInt(b.current)
		if i64 := *ptr; b.scanErr == nil && strconv.IntSize == 32 && (i64 < -math.MaxInt32 || i64 > math.MaxInt32) { // hmm rethink that depending on goarch
			b.scanErr = rangeError("ColumnMap.Int", string(b.current))
		}
		if b.scanErr != nil {
			b.scanErr = errors.Wrapf(b.scanErr, "[dml] Column %q", b.Column())
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
		var nv sql.NullInt64
		nv, b.scanErr = byteconv.ParseNullInt64(b.current)
		*ptr = NullInt64{NullInt64: nv}
		if b.scanErr != nil {
			b.scanErr = errors.Wrapf(b.scanErr, "[dml] Column %q", b.Column())
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
		*ptr, b.scanErr = byteconv.ParseFloat(b.current)
		if b.scanErr != nil {
			b.scanErr = errors.Wrapf(b.scanErr, "[dml] Column %q", b.Column())
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
		var nv sql.NullFloat64
		nv, b.scanErr = byteconv.ParseNullFloat64(b.current)
		*ptr = NullFloat64{NullFloat64: nv}
		if b.scanErr != nil {
			b.scanErr = errors.Wrapf(b.scanErr, "[dml] Column %q", b.Column())
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
		var u64 uint64
		u64, _, b.scanErr = byteconv.ParseUintSQL(b.current, 10, strconv.IntSize)
		*ptr = uint(u64)
		if b.scanErr != nil {
			b.scanErr = errors.Wrapf(b.scanErr, "[dml] Column %q", b.Column())
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
		var u64 uint64
		u64, _, b.scanErr = byteconv.ParseUintSQL(b.current, 10, 8)
		*ptr = uint8(u64)
		if b.scanErr != nil {
			b.scanErr = errors.Wrapf(b.scanErr, "[dml] Column %q", b.Column())
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
		var u64 uint64
		u64, _, b.scanErr = byteconv.ParseUintSQL(b.current, 10, 16)
		*ptr = uint16(u64)
		if b.scanErr != nil {
			b.scanErr = errors.Wrapf(b.scanErr, "[dml] Column %q", b.Column())
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
		var u64 uint64
		u64, _, b.scanErr = byteconv.ParseUintSQL(b.current, 10, 32)
		*ptr = uint32(u64)
		if b.scanErr != nil {
			b.scanErr = errors.Wrapf(b.scanErr, "[dml] Column %q", b.Column())
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
		*ptr, _, b.scanErr = byteconv.ParseUintSQL(b.current, 10, strconv.IntSize)
		if b.scanErr != nil {
			b.scanErr = errors.Wrapf(b.scanErr, "[dml] Column %q", b.Column())
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
		b := *b.scanRaw[i]
		if b == nil {
			_, _ = w.Write(tNil)
		} else {
			if _, err = fmt.Fprintf(w, ": %q", string(b)); err != nil {
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
		*ptr = append((*ptr)[:0], b.current...)
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
	if b.CheckValidUTF8 && !utf8.Valid(b.current) {
		b.scanErr = errors.NewNotValidf("[dml] Column %q at position %d contains invalid UTF-8 characters", b.Column(), b.Count)
	}
	if b.scanErr == nil {
		*ptr = string(b.current)
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
	if b.CheckValidUTF8 && !utf8.Valid(b.current) {
		b.scanErr = errors.NewNotValidf("[dml] Column Index %d at position %d contains invalid UTF-8 characters", b.index, b.Count)
	}
	if b.scanErr == nil {
		*ptr = NullString{NullString: byteconv.ParseNullString(b.current)}
	} else {
		b.scanErr = errors.Wrapf(b.scanErr, "[dml] Column %q", b.Column())
	}
	return b
}

// Time reads a time.Time value and appends it to the arguments slice or assigns
// the time.Time value stored in sql.RawBytes to the pointer. See the
// documentation for function Scan.
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
		*ptr, b.scanErr = time.Parse(time.RFC3339Nano, string(b.current))
		if b.scanErr != nil {
			b.scanErr = errors.Wrapf(b.scanErr, "[dml] Column %q", b.Column())
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
		if err := ptr.Scan(b.current); err != nil {
			b.scanErr = errors.NewNotValidf("[dml] ColumnMap NullTime: Invalid time string: %q with error %s", string(b.current), err)
		}
	}
	return b
}

func rangeError(fn, str string) *strconv.NumError {
	return &strconv.NumError{Func: fn, Num: str, Err: strconv.ErrRange}
}

// BuilderBase contains fields which all SQL query builder have in common, the
// same base. Exported for documentation reasons.
type BuilderBase struct {
	// ID of a statement. Used in logging. The ID gets generated with function
	// signature `func() string`. This func gets applied to the logger when
	// setting up a logger.
	id           string
	Log          log.Logger // Log optional logger
	RawFullSQL   string
	RawArguments Arguments // args used by RawFullSQL

	Table id

	// PropagationStopped set to true if you would like to interrupt the
	// listener chain. Once set to true all sub sequent calls of the next
	// listeners will be suppressed.
	PropagationStopped bool
	IsInterpolate      bool // See Interpolate()
	IsBuildCache       bool // see BuildCache()
	// IsUnsafe if set to true the functions AddColumn* will turn any
	// non valid identifier (not `{a-z}[a-z0-9$_]+`i) into an expression.
	IsUnsafe bool
	cacheSQL []byte
	argPool  Arguments // like a buffer, gets reused internally, so a pool.
	// propagationStoppedAt position in the slice where the stopped propagation
	// has been requested. for every new iteration the propagation must stop at
	// this position.
	propagationStoppedAt int
}

// BuilderConditional defines base fields used in statements which can have
// conditional constraints like WHERE, JOIN, ORDER, etc. Exported for
// documentation reasons.
type BuilderConditional struct {
	// QualifiedRecords represents a map of ColumnMappers to retrieve the
	// necessary arguments from the interface implementations of the types. The
	// string key (the qualifier) can either be the table or object name or in
	// cases, where an alias gets used, the string key must be the same as the
	// alias. The map get called internally when the arguments are getting
	// assembled.
	QualifiedRecords map[string]ColumnMapper
	Joins            Joins
	Wheres           Conditions
	OrderBys         ids
	LimitCount       uint64
	LimitValid       bool
}

func (b *BuilderConditional) join(j string, t id, on ...*Condition) {
	jf := &join{
		JoinType: j,
		Table:    t,
	}
	jf.On = append(jf.On, on...)
	b.Joins = append(b.Joins, jf)
}

// StmtBase wraps a *sql.Stmt (a prepared statement) with a specific SQL query.
// To create a StmtBase call the Prepare function of type Select. StmtBase is
// not safe for concurrent use, despite the underlying *sql.Stmt is. Don't
// forget to call Close!
type StmtBase struct {
	id   string // tracing ID
	stmt *sql.Stmt
	// argsCache can be a sync.Pool and when calling Close function the
	// interface slice gets returned to the pool
	argsCache Arguments
	// argsRaw can be a sync.Pool and when calling Close function the interface
	// slice gets returned to the pool
	argsRaw []interface{}
	// isWithInterfaces will be set to true if the raw interface arguments are
	// getting applied.
	isWithInterfaces bool
	// ärgErr represents an argument error caused in one of the three With
	// functions.
	ärgErr     error // Sorry Germans for that terrible pun #notSorry
	bindRecord func(records []QualifiedRecord)
	log        log.Logger
}

// Close closes the underlying prepared statement.
func (st *StmtBase) Close() error { return st.stmt.Close() }

func (st *StmtBase) withArgs(args []interface{}) {
	st.argsCache = st.argsCache[:0]
	st.argsRaw = st.argsRaw[:0]
	st.argsRaw = append(st.argsRaw, args...)
	st.isWithInterfaces = true
}

func (st *StmtBase) withArguments(args Arguments) {
	st.argsCache = st.argsCache[:0]
	st.argsCache = append(st.argsCache, args...)
	st.isWithInterfaces = false
}

// withRecords sets the records for the execution with Query or Exec. It
// internally resets previously applied arguments.
func (st *StmtBase) withRecords(appendArgs func(Arguments) (Arguments, error), records ...QualifiedRecord) {
	st.argsCache = st.argsCache[:0]
	st.bindRecord(records)
	st.argsCache, st.ärgErr = appendArgs(st.argsCache)
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
		st.argsRaw = st.argsCache.Interfaces(st.argsRaw...)
	}
	st.argsRaw = append(st.argsRaw, args...)
	return nil
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
	if st.log != nil && st.log.IsDebug() {
		defer log.WhenDone(st.log).Debug("Exec", log.Int("arg_len", len(st.argsRaw)))
	}
	return st.stmt.ExecContext(ctx, st.argsRaw...)
}

// Query traditional way, allocation heavy.
func (st *StmtBase) Query(ctx context.Context, args ...interface{}) (*sql.Rows, error) {
	if err := st.prepareArgs(args...); err != nil {
		return nil, errors.WithStack(err)
	}
	if st.log != nil && st.log.IsDebug() {
		defer log.WhenDone(st.log).Debug("Query", log.Int("arg_len", len(st.argsRaw)))
	}
	return st.stmt.QueryContext(ctx, st.argsRaw...)
}

// QueryRow traditional way, allocation heavy.
func (st *StmtBase) QueryRow(ctx context.Context, args ...interface{}) *sql.Row {
	if err := st.prepareArgs(args...); err != nil {
		_ = err
		// Hmmm what should happen here?
	}
	if st.log != nil && st.log.IsDebug() {
		defer log.WhenDone(st.log).Debug("QueryRow", log.Int("arg_len", len(st.argsRaw)))
	}
	return st.stmt.QueryRowContext(ctx, st.argsRaw...)
}

// Load loads data from a query into an object. You must set DB.QueryContext on
// the Select object or it just panics. Load can load a single row or n-rows.
func (st *StmtBase) Load(ctx context.Context, s ColumnMapper) (rowCount uint64, err error) {
	if st.log != nil && st.log.IsDebug() {
		defer log.WhenDone(st.log).Debug("Load", log.Uint64("row_count", rowCount), log.String("object_type", fmt.Sprintf("%T", s)))
	}
	r, err := st.Query(ctx)
	rowCount, err = load(r, err, s)
	return rowCount, errors.WithStack(err)
}

// LoadInt64 executes the prepared statement and returns the value at an int64.
// It returns a NotFound error if the query returns nothing.
func (st *StmtBase) LoadInt64(ctx context.Context) (int64, error) {
	if st.log != nil && st.log.IsDebug() {
		defer log.WhenDone(st.log).Debug("LoadInt64")
	}
	return loadInt64(st.Query(ctx))
}

// LoadInt64s executes the Select and returns the value as a slice of int64s.
func (st *StmtBase) LoadInt64s(ctx context.Context) (ret []int64, err error) {
	if st.log != nil && st.log.IsDebug() {
		// do not use fullSQL because we might log sensitive data
		defer log.WhenDone(st.log).Debug("LoadInt64s", log.Int("row_count", len(ret)))
	}
	ret, err = loadInt64s(st.Query(ctx))
	// Do not simplify it because we need ret in the defer. we don't log errors
	// because they get handled.
	return ret, err
}

// More Load* functions can be added later
