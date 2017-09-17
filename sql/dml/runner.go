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
// Load checks if a type implements RowCloser interface.
// `db` can be either a *sql.DB (connection pool), a *sql.Conn (a single
// dedicated database session) or a *sql.Tx (an in-progress database
// transaction).
func Load(ctx context.Context, db Querier, b QueryBuilder, s Scanner) (rowCount int64, err error) {
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

func load(r *sql.Rows, errIn error, s Scanner) (rowCount int64, err error) {
	if errIn != nil {
		return 0, errors.WithStack(errIn)
	}
	defer func() {
		// Not testable with the sqlmock package :-(
		if err2 := r.Close(); err2 != nil && err == nil {
			err = errors.Wrap(err2, "[dml] Load.QueryContext.Rows.Close")
		}
		if rc, ok := s.(RowCloser); ok {
			if err2 := rc.RowClose(); err2 != nil && err == nil {
				err = errors.Wrap(err2, "[dml] Load.QueryContext.Scanner.RowClose")
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
// as a composite field in a database table struct. Does not support streaming
// because neither database/sql does :-(  The method receiver functions have the
// same names as in type RowConvert.
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
	scanArgs       []interface{} // could be a sync.Pool but check it in benchmarks.
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
func (b *RowConvert) NullBool() (NullBool, error) {
	nv, err := byteconv.ParseNullBool(b.current)
	return NullBool{NullBool: nv}, err
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
func (b *RowConvert) NullInt64() (NullInt64, error) {
	nv, err := byteconv.ParseNullInt64(b.current)
	return NullInt64{NullInt64: nv}, err
}

// Float64 see the documentation for function Scan.
func (b *RowConvert) Float64() (float64, error) {
	return byteconv.ParseFloat(b.current)
}

// NullFloat64 see the documentation for function Scan.
func (b *RowConvert) NullFloat64() (NullFloat64, error) {
	nv, err := byteconv.ParseNullFloat64(b.current)
	return NullFloat64{NullFloat64: nv}, err
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

// Debug writes the column names with their values into `w`. The output format
// might change.
func (b *RowConvert) Debug(w io.Writer) (err error) {
	nl := []byte("\n")
	tNil := []byte(": <nil>")
	for i, c := range b.Columns {
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

// String see the documentation for function Scan.
func (b *RowConvert) String() (string, error) {
	if b.CheckValidUTF8 && !utf8.Valid(b.current) {
		return "", errors.NewNotValidf("[dml] Column Index %d at position %d contains invalid UTF-8 characters", b.index, b.Count)
	}
	return string(b.current), nil
}

// NullString see the documentation for function Scan.
func (b *RowConvert) NullString() (NullString, error) {
	if b.CheckValidUTF8 && !utf8.Valid(b.current) {
		return NullString{}, errors.NewNotValidf("[dml] Column Index %d at position %d contains invalid UTF-8 characters", b.index, b.Count)
	}
	return NullString{NullString: byteconv.ParseNullString(b.current)}, nil
}

// String see the documentation for function Scan.
func (b *RowConvert) Time() (time.Time, error) {
	t, err := time.Parse(time.RFC3339Nano, string(b.current))
	if err != nil {
		return time.Time{}, errors.NewNotValidf("[dml] RowConvert Time: Invalid time string: %q", string(b.current))
	}
	return t, nil
}

// NullString see the documentation for function Scan.
func (b *RowConvert) NullTime() (NullTime, error) {
	var nt NullTime
	if err := nt.Scan(b.current); err != nil {
		return NullTime{}, errors.NewNotValidf("[dml] RowConvert NullTime: Invalid time string: %q", string(b.current))
	}
	return nt, nil
}

func rangeError(fn, str string) *strconv.NumError {
	return &strconv.NumError{Func: fn, Num: str, Err: strconv.ErrRange}
}

// BuilderBase contains fields which all SQL query builder have in common, the
// same base. Exported for documentation reasons.
type BuilderBase struct {
	// ID of a statement. Used in logging. If empty the generated SQL string
	// gets used which can might contain sensitive information which should not
	// get logged. TODO implement
	id           string
	Log          log.Logger // Log optional logger
	RawFullSQL   string
	RawArguments Arguments // args used by RawFullSQL

	Table identifier

	// PropagationStopped set to true if you would like to interrupt the
	// listener chain. Once set to true all sub sequent calls of the next
	// listeners will be suppressed.
	PropagationStopped bool
	IsInterpolate      bool // See Interpolate()
	IsBuildCache       bool // see BuildCache()
	// IsUnsafe if set to true the functions AddColumn* will turn any
	// non valid identifier (not `{a-z}[a-z0-9$_]+`i) into an expression.
	IsUnsafe  bool
	cacheSQL  []byte
	cacheArgs Arguments // like a buffer, gets reused
	// propagationStoppedAt position in the slice where the stopped propagation
	// has been requested. for every new iteration the propagation must stop at
	// this position.
	propagationStoppedAt int
}

// BuilderConditional defines base fields used in statements which can have
// conditional constraints like WHERE, JOIN, ORDER, etc. Exported for
// documentation reasons.
type BuilderConditional struct {
	// ArgumentsAppender a map of ArgumentsAppender to retrieve the necessary
	// arguments from the interface implementations of the objects. The string
	// key (the qualifier) can either be the table or object name or in cases,
	// where an alias gets used, the string key must be the same as the alias.
	// The map get called internally when the arguments are getting assembled.
	ArgumentsAppender map[string]ArgumentsAppender
	Joins             Joins
	Wheres            Conditions
	OrderBys          identifiers
	LimitCount        uint64
	LimitValid        bool
}

func (b *BuilderConditional) join(j string, t identifier, on ...*Condition) {
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
func (st *StmtBase) Load(ctx context.Context, s Scanner) (rowCount int64, err error) {
	if st.log != nil && st.log.IsDebug() {
		defer log.WhenDone(st.log).Debug("Load", log.Int64("row_count", rowCount), log.String("object_type", fmt.Sprintf("%T", s)))
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
