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
	"bytes"
	"context"
	"database/sql"
	"strings"

	"github.com/corestoreio/errors"
	"github.com/corestoreio/log"
	"github.com/corestoreio/pkg/util/bufferpool"
)

// LastInsertIDAssigner assigns the last insert ID of an auto increment
// column back to the objects.
type LastInsertIDAssigner interface {
	AssignLastInsertID(int64)
}

// Insert contains the clauses for an INSERT statement
type Insert struct {
	BuilderBase
	// DB can be either a *sql.DB (connection pool), a *sql.Conn (a single
	// dedicated database session) or a *sql.Tx (an in-progress database
	// transaction).
	DB ExecPreparer

	Into    string
	Columns []string
	// RowCount defines the number of expected rows.
	RowCount int // See SetRowCount()
	Values   []Arguments
	Records  []ColumnMapper
	// RecordPlaceHolderCount defines the number of place holders for each set
	// within the brackets. Must only be set when Records have been applied
	// and `Columns` field has been omitted.
	RecordPlaceHolderCount int
	// Select used to create an "INSERT INTO `table` SELECT ..." statement.
	Select *Select

	// OnDuplicateKeys updates the referenced columns. See documentation for type
	// `Conditions`. For more details
	// https://dev.mysql.com/doc/refman/5.7/en/insert-on-duplicate.html
	// Conditions contains the column/argument association for either the SET
	// clause in an UPDATE statement or to be used in an INSERT ... ON DUPLICATE KEY
	// statement. For each column there must be one argument which can either be nil
	// or has an actual value.
	//
	// When using the ON DUPLICATE KEY feature in the Insert builder:
	//
	// The function dml.ExpressionValue is supported and allows SQL
	// constructs like (ib == InsertBuilder builds INSERT statements):
	// 		`columnA`=VALUES(`columnB`)+2
	// by writing the Go code:
	//		ib.AddOnDuplicateKey("columnA", ExpressionValue("VALUES(`columnB`)+?", Int(2)))
	// Omitting the argument and using the keyword nil will turn this Go code:
	//		ib.AddOnDuplicateKey("columnA", nil)
	// into that SQL:
	// 		`columnA`=VALUES(`columnA`)
	// Same applies as when the columns gets only assigned without any arguments:
	//		ib.OnDuplicateKeys.Columns = []string{"name","sku"}
	// will turn into:
	// 		`name`=VALUES(`name`), `sku`=VALUES(`sku`)
	// Type `Conditions` gets used in type `Update` with field
	// `SetClauses` and in type `Insert` with field OnDuplicateKeys.
	OnDuplicateKeys Conditions
	// OnDuplicateKeyExclude excludes the mentioned columns to the ON DUPLICATE
	// KEY UPDATE section. Otherwise all columns in the field `Columns` will be
	// added to the ON DUPLICATE KEY UPDATE expression. Usually the slice
	// `OnDuplicateKeyExclude` contains the primary key columns. Case-sensitive
	// comparison.
	OnDuplicateKeyExclude []string
	// IsOnDuplicateKey if enabled adds all columns to the ON DUPLICATE KEY
	// claus. Takes the OnDuplicateKeyExclude field into consideration.
	IsOnDuplicateKey bool
	// IsReplace uses the REPLACE syntax. See function Replace().
	IsReplace bool
	// IsIgnore ignores error. See function Ignore().
	IsIgnore bool
	// Listeners allows to dispatch certain functions in different
	// situations.
	Listeners ListenersInsert
	argsPool  Arguments // only used during ToSQL
}

// NewInsert creates a new Insert object.
func NewInsert(into string) *Insert {
	return &Insert{
		Into: into,
	}
}

func newInsertInto(db ExecPreparer, idFn uniqueIDFn, l log.Logger, into string) *Insert {
	id := idFn()
	if l != nil {
		l = l.With(log.String("insert_id", id), log.String("table", into))
	}
	return &Insert{
		BuilderBase: BuilderBase{
			builderCommon: builderCommon{
				id:  id,
				Log: l,
			},
		},
		DB:   db,
		Into: into,
	}
}

// InsertInto instantiates a Insert for the given table
func (c *ConnPool) InsertInto(into string) *Insert {
	return newInsertInto(c.DB, c.makeUniqueID, c.Log, into)
}

// InsertInto instantiates a Insert for the given table
func (c *Conn) InsertInto(into string) *Insert {
	return newInsertInto(c.DB, c.makeUniqueID, c.Log, into)
}

// InsertInto instantiates a Insert for the given table bound to a transaction
func (tx *Tx) InsertInto(into string) *Insert {
	return newInsertInto(tx.DB, tx.makeUniqueID, tx.Log, into)
}

// WithDB sets the database query object.
func (b *Insert) WithDB(db ExecPreparer) *Insert {
	b.DB = db
	return b
}

// Ignore modifier enables errors that occur while executing the INSERT
// statement are getting ignored. For example, without IGNORE, a row that
// duplicates an existing UNIQUE index or PRIMARY KEY value in the table causes
// a duplicate-key error and the statement is aborted. With IGNORE, the row is
// discarded and no error occurs. Ignored errors generate warnings instead.
// https://dev.mysql.com/doc/refman/5.7/en/insert.html
func (b *Insert) Ignore() *Insert {
	b.IsIgnore = true
	return b
}

// Replace instead of INSERT to overwrite old rows. REPLACE is the counterpart
// to INSERT IGNORE in the treatment of new rows that contain unique key values
// that duplicate old rows: The new rows are used to replace the old rows rather
// than being discarded.
// https://dev.mysql.com/doc/refman/5.7/en/replace.html
func (b *Insert) Replace() *Insert {
	b.IsReplace = true
	return b
}

// AddColumns appends columns and increases the `RecordPlaceHolderCount` variable.
func (b *Insert) AddColumns(columns ...string) *Insert {
	b.RecordPlaceHolderCount += len(columns)
	b.Columns = append(b.Columns, columns...)
	return b
}

// SetRowCount defines the number of expected rows. Each set of place holders
// within the brackets defines a row. This setting defaults to one. It gets
// applied when fields `args` and `Records` have been left empty. For each
// defined column the QueryBuilder creates a place holder. Use when creating a
// prepared statement. See the example for more details.
// 		RowCount = 2 ==> (?,?,?),(?,?,?)
// 		RowCount = 3 ==> (?,?,?),(?,?,?),(?,?,?)
func (b *Insert) SetRowCount(rows int) *Insert {
	b.RowCount = rows
	return b
}

// AddValuesUnsafe appends a set of primitives, packed in interfaces, to the
// statement. Each call of AddValuesUnsafe creates a new set of values. Only
// primitive types are supported. Runtime type safety only. It panics for
// unknown types.
func (b *Insert) AddValuesUnsafe(args ...interface{}) *Insert {
	return b.AddValues(iFaceToArgs(args...))
}

// AddValues appends a set of arguments to the statement. Each call of
// AddValues creates a new set of values. Only primitive types are supported.
func (b *Insert) AddValues(args Arguments) *Insert {
	if lv, mod := len(args), len(b.Columns); mod > 0 && lv > mod && (lv%mod) == 0 {
		// now we have more arguments than columns and we can assume that more
		// rows gets inserted.
		for i := 0; i < len(args); i = i + mod {
			b.Values = append(b.Values, args[i:i+mod])
		}
	} else {
		// each call to AddValuesUnsafe equals one row in a table.
		b.Values = append(b.Values, args)
	}
	return b
}

// AddRecords appends a new record for each INSERT VALUES (),[(...)...] case. A
// record can also be e.g. a slice which appends all requested arguments at
// once. Using a slice requires to call `SetRowCount` to tell the Insert object
// the number of rows.
func (b *Insert) AddRecords(recs ...ColumnMapper) *Insert {
	b.Records = append(b.Records, recs...)
	return b
}

// WithArgs applies only in the case where place holders are used in the ON
// DUPLICATE KEY part. It sets the interfaced arguments for the execution with
// Query+. It internally resets previously applied arguments. This function does
// not support interpolation.
func (b *Insert) WithArgs(args ...interface{}) *Insert {
	b.withArgs(args)
	return b
}

// WithArguments applies only in the case where place holders are used in the ON
// DUPLICATE KEY part. It sets the arguments for the execution with Query+. It
// internally resets previously applied arguments. This function supports
// interpolation.
func (b *Insert) WithArguments(args Arguments) *Insert {
	b.withArguments(args)
	return b
}

// Reset resets the Records and Values slices to be reused. It sets the records
// slice items to nil for the GC.
func (b *Insert) Reset() *Insert {
	for i := 0; i < len(b.Records); i++ {
		b.Records[i] = nil // remove pointer, etc for GC
	}
	b.Records = b.Records[:0]
	b.Values = b.Values[:0]
	return b
}

// SetRecordPlaceHolderCount number of expected place holders within each set.
// Must be applied if a call to AddColumns has been omitted and WithRecords gets
// called or Records gets set in a different way.
//		INSERT INTO tableX (?,?,?)
// SetRecordPlaceHolderCount would now be 3 because of the three place holders.
func (b *Insert) SetRecordPlaceHolderCount(valueCount int) *Insert {
	// maybe we can do better and remove this method ...
	b.RecordPlaceHolderCount = valueCount
	return b
}

// AddOnDuplicateKey has some hidden features for best flexibility. You can only
// set the Columns itself to allow the following SQL construct:
//		`columnA`=VALUES(`columnA`)
// Means columnA gets automatically mapped to the VALUES column name.
func (b *Insert) AddOnDuplicateKey(c ...*Condition) *Insert {
	b.OnDuplicateKeys = append(b.OnDuplicateKeys, c...)
	return b
}

// AddOnDuplicateKeyExclude adds a column to the exclude list. As soon as a
// column gets set with this function the ON DUPLICATE KEY clause gets
// generated. Usually the slice `OnDuplicateKeyExclude` contains the
// primary/unique key columns. Case-sensitive comparison.
func (b *Insert) AddOnDuplicateKeyExclude(columns ...string) *Insert {
	b.OnDuplicateKeyExclude = append(b.OnDuplicateKeyExclude, columns...)
	return b
}

// OnDuplicateKey enables for all columns to be written into the ON DUPLICATE
// KEY claus. Takes the field OnDuplicateKeyExclude into consideration.
func (b *Insert) OnDuplicateKey() *Insert {
	b.IsOnDuplicateKey = true
	return b
}

// Pair appends a column/value pair to the statement. Calling this function
// multiple times with the same column name produces invalid SQL. Slice values
// and right/left side expressions are not supported and ignored.
func (b *Insert) Pair(cvs ...*Condition) *Insert {
	// TODO(CyS) support right side expressions, requires some internal refactoring
	for _, cv := range cvs {
		colPos := -1
		for i, c := range b.Columns {
			if strings.EqualFold(c, cv.Left) {
				colPos = i
				break
			}
		}
		if colPos == -1 {
			b.Columns = append(b.Columns, cv.Left)
			if len(b.Values) == 0 {
				b.Values = make([]Arguments, 1, 5)
			}
			b.Values[0] = append(b.Values[0], cv.Right.arg)

		} else { // this is not an ELSEIF
			if colPos == 0 { // create new slice
				b.Values = append(b.Values, make(Arguments, len(b.Columns)))
			}
			pos := len(b.Values) - 1
			b.Values[pos][colPos] = cv.Right.arg
		}
	}
	return b
}

// FromSelect creates an "INSERT INTO `table` SELECT ..." statement from a
// previously created SELECT statement.
func (b *Insert) FromSelect(s *Select) *Insert {
	b.Select = s
	return b
}

// Interpolate if set stringyfies the arguments into the SQL string and returns
// pre-processed SQL command when calling the function ToSQL. Not suitable for
// prepared statements. ToSQLs second argument `args` will then be nil.
func (b *Insert) Interpolate() *Insert {
	b.IsInterpolate = true
	return b
}

// ToSQL serialized the Insert to a SQL string
// It returns the string with placeholders and a slice of query arguments
func (b *Insert) ToSQL() (string, []interface{}, error) {
	rawSQL, err := b.buildToSQL(b)
	if err != nil {
		return "", nil, errors.WithStack(err)
	}
	if b.argsPool == nil {
		b.argsPool = make(Arguments, 0, b.totalArgCount())
	}
	b.argsPool = b.argsPool[:0]
	b.argsPool, err = b.appendArgs(b.argsPool)
	if err != nil {
		return "", nil, errors.WithStack(err)
	}

	if b.IsInterpolate {
		buf := bufferpool.Get()
		err := writeInterpolate(buf, rawSQL, b.argsPool)
		s := buf.String()
		bufferpool.Put(buf)
		return s, nil, err
	}

	return string(rawSQL), append(b.argsPool.Interfaces(), b.argsRaw...), nil
}

func (b *Insert) writeBuildCache(sql []byte) {
	// think about resetting ...
	b.cachedSQL = sql
}

func (b *Insert) readBuildCache() (sql []byte) {
	return b.cachedSQL
}

// DisableBuildCache if enabled it does not cache the SQL string as a final
// rendered byte slice. Allows you to rebuild the query with different
// statements.
func (b *Insert) DisableBuildCache() *Insert {
	b.IsBuildCacheDisabled = true
	return b
}

func (b *Insert) toSQL(buf *bytes.Buffer, placeHolders []string) ([]string, error) {

	if err := b.Listeners.dispatch(OnBeforeToSQL, b); err != nil {
		return nil, errors.WithStack(err)
	}

	if b.RawFullSQL != "" {
		buf.WriteString(b.RawFullSQL)
		return placeHolders, nil
	}

	if len(b.Into) == 0 {
		return nil, errors.Empty.Newf("[dml] Inserted table is missing")
	}

	ior := "INSERT "
	if b.IsReplace {
		ior = "REPLACE "
	}
	buf.WriteString(ior)
	writeStmtID(buf, b.id)
	if b.IsIgnore {
		buf.WriteString("IGNORE ")
	}

	buf.WriteString("INTO ")
	Quoter.quote(buf, b.Into)
	buf.WriteByte(' ')

	if b.Select != nil {
		if len(b.Columns) > 0 {
			buf.WriteByte('(')
			for i, c := range b.Columns {
				if i > 0 {
					buf.WriteByte(',')
				}
				Quoter.quote(buf, c)
			}
			buf.WriteString(") ")
		}
		ph, err := b.Select.toSQL(buf, placeHolders)
		return ph, errors.WithStack(err)
	}

	if len(b.Columns) > 0 {
		buf.WriteByte('(')
		for i, c := range b.Columns {
			if i > 0 {
				buf.WriteByte(',')
			}
			Quoter.quote(buf, c)
		}
		placeHolders = append(placeHolders, b.Columns...)
		buf.WriteString(") ")
	}
	buf.WriteString("VALUES ")

	rowCount := 1
	if lr := len(b.Records); lr > 0 {
		rowCount = lr
	}
	if b.RowCount > 0 { // no switch statement
		rowCount = b.RowCount
	}

	var argCount0 int
	if b.Records == nil {
		argCount0 = len(b.Columns)
		lv := len(b.Values)
		if lv > 0 {
			rowCount = lv
		}
		if argCount0 == 0 && lv > 0 {
			argCount0 = len(b.Values[0])
		}
	} else {
		argCount0 = len(b.Columns)
		if b.RecordPlaceHolderCount > 0 {
			argCount0 = b.RecordPlaceHolderCount
		}
	}

	// write the place holders: (?,?,?)[,(?,?,?)...]
	for vc := 0; vc < rowCount; vc++ {
		if vc > 0 {
			buf.WriteByte(',')
		}
		buf.WriteByte('(')
		for i := 0; i < argCount0; i++ {
			if i > 0 {
				buf.WriteByte(',')
			}
			buf.WriteByte(placeHolderRune)
		}
		buf.WriteByte(')')
	}
	if len(b.OnDuplicateKeyExclude) > 0 || b.IsOnDuplicateKey {
		if len(b.OnDuplicateKeys) == 0 {
			b.OnDuplicateKeys = append(b.OnDuplicateKeys, &Condition{})
		}
	ColumnsLoop:
		for _, c := range b.Columns {
			// Wow two times a comparison with a slice. That costs a bit
			// performance but a reliable way to avoid writing duplicate ON
			// DUPLICATE KEY UPDATE sets. If there is something faster, write us.
			if strInSlice(c, b.OnDuplicateKeyExclude) {
				continue
			}
			for _, cnd := range b.OnDuplicateKeys {
				if c == cnd.Left || strInSlice(c, cnd.Columns) {
					continue ColumnsLoop
				}
			}
			b.OnDuplicateKeys[0].Columns = append(b.OnDuplicateKeys[0].Columns, c)
		}
	}

	placeHolders, err := b.OnDuplicateKeys.writeOnDuplicateKey(buf, placeHolders)
	return placeHolders, errors.Wrap(err, "[dml] Insert.toSQL.writeOnDuplicateKey\n")
}

func strInSlice(search string, sl []string) bool {
	for _, s := range sl {
		if s == search {
			return true
		}
	}
	return false
}

func (b *Insert) totalArgCount() int {
	argCount0 := b.RecordPlaceHolderCount
	if len(b.Values) > 0 && len(b.Columns) == 0 {
		argCount0 = len(b.Values[0])
	}
	if lc := len(b.Columns); argCount0 < 1 && lc > 0 {
		argCount0 = lc
	}

	return len(b.Values) * argCount0
}

func (b *Insert) appendArgs(args Arguments) (_ Arguments, err error) {

	if b.Select != nil && (b.Select.argsArgs != nil || b.Select.argsRecords != nil) {
		args, err = b.Select.convertRecordsToArguments()
		return args, errors.WithStack(err)
	}

	argCount0 := b.RecordPlaceHolderCount
	if len(b.Values) > 0 && len(b.Columns) == 0 {
		argCount0 = len(b.Values[0])
	}
	if lc := len(b.Columns); argCount0 < 1 && lc > 0 {
		argCount0 = lc
	}
	totalArgCount := len(b.Values) * argCount0
	if cap(args) == 0 {
		args = make(Arguments, 0, totalArgCount+len(b.Records)+len(b.OnDuplicateKeys))
	}
	for _, v := range b.Values {
		args = append(args, v...)
	}

	if b.Records == nil {
		args = append(args, b.argsArgs...)
		return args, errors.WithStack(err)
	}

	cm := newColumnMap(args, b.qualifiedColumns...) // b.Columns can be nil
	for _, rec := range b.Records {
		alBefore := len(cm.Args)
		if err = rec.MapColumns(cm); err != nil {
			return nil, errors.WithStack(err)
		}
		if addedArgs := len(cm.Args) - alBefore; argCount0 > 0 && addedArgs%argCount0 != 0 {
			return nil, errors.Mismatch.Newf("[dml] Insert.appendArgs RecordPlaceHolderCount(%d) does not match the number of assembled arguments (%d)", b.RecordPlaceHolderCount, addedArgs)
		}
	}
	args = cm.Args

	args = append(args, b.argsArgs...)

	return args, nil
}

// Exec executes the statement represented by the Insert object. It returns the
// raw database/sql Result or an error if there was one. Regarding
// LastInsertID(): If you insert multiple rows using a single INSERT statement,
// LAST_INSERT_ID() returns the value generated for the first inserted row only.
// The reason for this at to make it possible to reproduce easily the same
// INSERT statement against some other server. If a record resp. and object
// implements the interface LastInsertIDAssigner then the LastInsertID gets
// assigned incrementally to the objects.
func (b *Insert) Exec(ctx context.Context) (sql.Result, error) {
	if b.Log != nil && b.Log.IsDebug() {
		defer log.WhenDone(b.Log).Debug("Exec", log.Stringer("sql", b))
	}
	result, err := Exec(ctx, b.DB, b)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	if b.Records == nil {
		return result, nil
	}
	lID, err := result.LastInsertId()
	if err != nil {
		return nil, errors.WithStack(err)
	}
	for i, rec := range b.Records {
		if a, ok := rec.(LastInsertIDAssigner); ok {
			a.AssignLastInsertID(lID + int64(i))
		}
	}
	return result, nil
}

// Prepare executes the statement represented by the Insert to create a prepared
// statement. It returns a custom statement type or an error if there was one.
// Provided arguments or records in the Insert are getting ignored. The provided
// context is used for the preparation of the statement, not for the execution
// of the statement. The returned Stmter is not safe for concurrent use, despite
// the underlying *sql.Stmt is.
func (b *Insert) Prepare(ctx context.Context) (Stmter, error) {
	if b.Log != nil && b.Log.IsDebug() {
		defer log.WhenDone(b.Log).Debug("Prepare", log.Stringer("sql", b))
	}
	sqlStmt, err := Prepare(ctx, b.DB, b)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	cap := len(b.Columns) * b.RecordPlaceHolderCount
	return &stmtInsert{
		stmtBase: stmtBase{
			builderCommon: builderCommon{
				id:       b.id,
				argsArgs: make(Arguments, 0, cap),
				argsRaw:  make([]interface{}, 0, cap),
				Log:      b.Log,
			},
			source: dmlTypeInsert,
			stmt:   sqlStmt,
		},
		ins: b,
	}, nil
}

// stmtInsert wraps a *sql.Stmt with a specific SQL query. To create a
// stmtInsert call the Prepare function of type Insert. stmtInsert is not safe
// for concurrent use, despite the underlying *sql.Stmt is. Don't forget to call
// Close!
type stmtInsert struct {
	stmtBase
	ins *Insert
}

// WithArguments sets the arguments for the execution with Exec. It internally resets
// previously applied arguments.
func (st *stmtInsert) WithArguments(args Arguments) Stmter {
	st.ins.Records = nil
	st.ins.Values = st.ins.Values[:0]
	st.argsArgs = st.argsArgs[:0]

	if lv, mod := len(args), len(st.ins.Columns); mod > 0 && lv > mod && (lv%mod) == 0 {
		// now we have more arguments than columns and we can assume that more
		// rows gets inserted.
		for i := 0; i < len(args); i = i + mod {
			st.ins.Values = append(st.ins.Values, args[i:i+mod])
		}
	} else {
		// each call to AddValuesUnsafe equals one row in a table.
		st.ins.Values = append(st.ins.Values, args)
	}

	st.argsArgs, st.ärgErr = st.ins.appendArgs(st.argsArgs)

	return st
}

// WithRecords sets the records for the execution with Exec. It internally
// resets previously applied arguments.
func (st *stmtInsert) WithRecords(records ...QualifiedRecord) Stmter {
	st.argsArgs = st.argsArgs[:0]
	st.ins.Records = nil
	recs := make([]ColumnMapper, len(records))
	for i, r := range records {
		recs[i] = r.Record
	}
	st.ins.AddRecords(recs...)
	st.argsArgs, st.ärgErr = st.ins.appendArgs(st.argsArgs)
	return st
}

// Exec executes a query with the previous set arguments or records or
// without arguments. It does not reset the internal arguments, so multiple
// executions with the same arguments/records are possible. Number of previously
// applied arguments or records must be the same as in the defined SQL but
// With*().ExecContext() can be called in a loop, both are not thread safe.
func (st *stmtInsert) Exec(ctx context.Context, args ...interface{}) (sql.Result, error) {

	result, err := st.stmtBase.Exec(ctx, args...)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	if st.ins.Records == nil {
		return result, nil
	}

	lID, err := result.LastInsertId()
	if err != nil {
		return nil, errors.WithStack(err)
	}
	for i, rec := range st.ins.Records {
		if a, ok := rec.(LastInsertIDAssigner); ok {
			a.AssignLastInsertID(lID + int64(i))
		}
	}
	return result, nil
}
