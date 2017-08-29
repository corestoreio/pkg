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
	"bytes"
	"context"
	"database/sql"
	"strings"

	"github.com/corestoreio/errors"
	"github.com/corestoreio/log"
)

// Insert contains the clauses for an INSERT statement
type Insert struct {
	BuilderBase
	// DB can be either a *sql.DB (connection pool), a *sql.Conn (a single
	// dedicated database session) or a *sql.Tx (an in-progress database
	// transaction).
	DB ExecPreparer

	Into    string
	Columns []string
	Values  []Arguments
	// RowCount defines the number of expected rows.
	RowCount int // See SetRowCount()

	Records []ArgumentsAppender
	// RecordValueCount defines the number of place holders for each value
	// within the brackets for each set. Must only be set when Records are set
	// and `Columns` field has been omitted.
	RecordValueCount int
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
	// The function dbr.ExpressionValue is supported and allows SQL
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
	Listeners InsertListeners
}

// NewInsert creates a new Insert object.
func NewInsert(into string) *Insert {
	return &Insert{
		Into: into,
	}
}

func newInsertInto(db ExecPreparer, l log.Logger, into, id string) *Insert {
	return &Insert{
		BuilderBase: BuilderBase{
			id:  id,
			Log: l,
		},
		DB:   db,
		Into: into,
	}
}

// InsertInto instantiates a Insert for the given table
func (c *ConnPool) InsertInto(into string) *Insert {
	l := c.Log
	id := c.makeUniqueID()
	if l != nil {
		l = c.Log.With(log.String("ConnPool", "Insert"), log.String("id", id), log.String("table", into))
	}
	return newInsertInto(c.DB, l, into, id)
}

// InsertInto instantiates a Insert for the given table
func (c *Conn) InsertInto(into string) *Insert {
	l := c.Log
	id := c.makeUniqueID()
	if l != nil {
		l = c.Log.With(log.String("Conn", "Insert"), log.String("id", id), log.String("table", into))
	}
	return newInsertInto(c.Conn, l, into, id)
}

// InsertInto instantiates a Insert for the given table bound to a transaction
func (tx *Tx) InsertInto(into string) *Insert {
	l := tx.Log
	id := tx.makeUniqueID()
	if l != nil {
		l = tx.Log.With(log.String("Tx", "Insert"), log.String("id", id), log.String("table", into))
	}
	return newInsertInto(tx.Tx, l, into, id)
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

// AddColumns appends columns and increases the `RecordValueCount` variable.
func (b *Insert) AddColumns(columns ...string) *Insert {
	b.RecordValueCount += len(columns)
	b.Columns = append(b.Columns, columns...)
	return b
}

// TODO (CyS) maybe write an intermediate type which will be used to get rid of
// AddValues and AddArguments. Maybe same pattern as Column() function.

// AddValues appends a set of values to the statement. Each call of AddValues
// creates a new set of values. Only primitive types are supported. Runtime type
// safety only.
func (b *Insert) AddValues(values ...interface{}) *Insert {
	return b.AddArguments(iFaceToArgs(values...))
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

// AddArguments appends a set of values to the statement. Each call of
// AddArguments creates a new set of values. Only primitive types are supported.
// Runtime type safety only.
func (b *Insert) AddArguments(args Arguments) *Insert {
	if lv, mod := len(args), len(b.Columns); mod > 0 && lv > mod && (lv%mod) == 0 {
		// now we have more arguments than columns and we can assume that more
		// rows gets inserted.
		for i := 0; i < len(args); i = i + mod {
			b.Values = append(b.Values, args[i:i+mod])
		}
	} else {
		// each call to AddValues equals one row in a table.
		b.Values = append(b.Values, args)
	}
	return b
}

// BindRecord appends a new record for each INSERT VALUES (),[(...)...] case. A
// record can also be e.g. a slice which appends all requested arguments at
// once. Using a slice requires to call `SetRowCount` to tell the Insert object
// the number of rows.
func (b *Insert) BindRecord(recs ...ArgumentsAppender) *Insert {
	b.Records = append(b.Records, recs...)
	return b
}

// SetRecordValueCount number of expected values within each set. Must be
// applied if a call to AddColumns has been omitted and BindRecord gets called
// or Records gets set in a different way.
//		INSERT INTO tableX (?,?,?)
// SetRecordValueCount would now be 3 because of the three place holders.
func (b *Insert) SetRecordValueCount(valueCount int) *Insert {
	// maybe we can do better and remove this method ...
	b.RecordValueCount = valueCount
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
			b.Values[0] = append(b.Values[0], cv.Right.Argument)

		} else { // this is not an ELSEIF
			if colPos == 0 { // create new slice
				b.Values = append(b.Values, make(Arguments, len(b.Columns)))
			}
			pos := len(b.Values) - 1
			b.Values[pos][colPos] = cv.Right.Argument
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
	return toSQL(b, b.IsInterpolate, _isNotPrepared)
}

func (b *Insert) writeBuildCache(sql []byte) {
	b.cacheSQL = sql
}

func (b *Insert) readBuildCache() (sql []byte, _ Arguments, err error) {
	if b.cacheSQL == nil {
		return nil, nil, nil
	}
	b.cacheArgs, err = b.appendArgs(b.cacheArgs[:0])
	return b.cacheSQL, b.cacheArgs, err
}

// BuildCache if `true` the final build query including place holders will be
// cached in a private field. Each time a call to function ToSQL happens, the
// arguments will be re-evaluated and returned or interpolated.
func (b *Insert) BuildCache() *Insert {
	b.IsBuildCache = true
	return b
}

func (b *Insert) hasBuildCache() bool {
	return b.IsBuildCache
}

func (b *Insert) toSQL(buf *bytes.Buffer) error {

	if err := b.Listeners.dispatch(OnBeforeToSQL, b); err != nil {
		return errors.WithStack(err)
	}

	if b.RawFullSQL != "" {
		buf.WriteString(b.RawFullSQL)
		return nil
	}

	if len(b.Into) == 0 {
		return errors.NewEmptyf("[dbr] Insert Table is missing")
	}

	ior := "INSERT "
	if b.IsReplace {
		ior = "REPLACE "
	}
	buf.WriteString(ior)
	if b.IsIgnore {
		buf.WriteString("IGNORE ")
	}

	buf.WriteString("INTO ")
	Quoter.quote(buf, b.Into)
	buf.WriteByte(' ')

	if b.Select != nil {
		return errors.WithStack(b.Select.toSQL(buf))
	}

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
		if b.RecordValueCount > 0 {
			argCount0 = b.RecordValueCount
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

	return errors.Wrap(b.OnDuplicateKeys.writeOnDuplicateKey(buf), "[dbr] Insert.toSQL.writeOnDuplicateKey\n")
}

func (b *Insert) appendArgs(args Arguments) (_ Arguments, err error) {

	if b.RawFullSQL != "" {
		return b.RawArguments, nil
	}

	if b.Select != nil {
		args, err = b.Select.appendArgs(args)
		return args, errors.WithStack(err)
	}

	if lv := len(b.Values); b.Records == nil && (lv == 0 || (lv > 0 && len(b.Values[0]) == 0)) {
		return nil, nil
	}

	argCount0 := b.RecordValueCount
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
		args, _, err = b.OnDuplicateKeys.appendArgs(args, appendArgsDUPKEY)
		return args, errors.WithStack(err)
	}

	for _, rec := range b.Records {
		alBefore := len(args)
		args, err = rec.AppendArgs(args, b.Columns) // b.Columns can be nil
		if err != nil {
			return nil, errors.WithStack(err)
		}
		if addedArgs := len(args) - alBefore; addedArgs%argCount0 != 0 {
			return nil, errors.NewMismatchf("[dbr] Insert.appendArgs RecordValueCount(%d) does not match the number of assembled arguments (%d)", b.RecordValueCount, addedArgs)
		}
	}

	if args, _, err = b.OnDuplicateKeys.appendArgs(args, appendArgsDUPKEY); err != nil {
		return nil, errors.WithStack(err)
	}

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
// of the statement.
func (b *Insert) Prepare(ctx context.Context) (*StmtInsert, error) {
	sqlStmt, err := Prepare(ctx, b.DB, b)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	cap := len(b.Columns) * b.RecordValueCount
	return &StmtInsert{
		StmtBase: StmtBase{
			stmt:      sqlStmt,
			argsCache: make(Arguments, 0, cap),
			argsRaw:   make([]interface{}, 0, cap),
		},
		ins: b,
	}, nil
}

// StmtInsert wraps a *sql.Stmt with a specific SQL query. To create a
// StmtInsert call the Prepare function of type Insert. StmtInsert is not safe
// for concurrent use, despite the underlying *sql.Stmt is. Don't forget to call
// Close!
type StmtInsert struct {
	StmtBase
	ins *Insert
}

// WithArguments sets the arguments for the execution with Exec. It internally resets
// previously applied arguments.
func (st *StmtInsert) WithArguments(args Arguments) *StmtInsert {
	st.ins.Records = nil
	st.ins.Values = st.ins.Values[:0]
	st.argsCache = st.argsCache[:0]

	if lv, mod := len(args), len(st.ins.Columns); mod > 0 && lv > mod && (lv%mod) == 0 {
		// now we have more arguments than columns and we can assume that more
		// rows gets inserted.
		for i := 0; i < len(args); i = i + mod {
			st.ins.Values = append(st.ins.Values, args[i:i+mod])
		}
	} else {
		// each call to AddValues equals one row in a table.
		st.ins.Values = append(st.ins.Values, args)
	}

	st.argsCache, st.ärgErr = st.ins.appendArgs(st.argsCache)

	return st
}

// WithRecords sets the records for the execution with Exec. It internally
// resets previously applied arguments.
func (st *StmtInsert) WithRecords(records ...ArgumentsAppender) *StmtInsert {
	st.argsCache = st.argsCache[:0]
	st.ins.Records = nil
	st.ins.BindRecord(records...)
	st.argsCache, st.ärgErr = st.ins.appendArgs(st.argsCache)
	return st
}

// ExecContext executes a query with the previous set arguments or records or
// without arguments. It does not reset the internal arguments, so multiple
// executions with the same arguments/records are possible. Number of previously
// applied arguments or records must be the same as in the defined SQL but
// With*().ExecContext() can be called in a loop, both are not thread safe.
func (st *StmtInsert) ExecContext(ctx context.Context, args ...interface{}) (sql.Result, error) {

	result, err := st.StmtBase.Exec(ctx, args...)
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

func strInSlice(search string, sl []string) bool {
	for _, s := range sl {
		if s == search {
			return true
		}
	}
	return false
}
