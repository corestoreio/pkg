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
	"strings"

	"github.com/corestoreio/errors"
	"github.com/corestoreio/log"
)

// Insert contains the clauses for an INSERT statement
type Insert struct {
	Log log.Logger // Log optional logger
	DB  ExecPreparer

	// UseBuildCache if `true` the final build query including place holders
	// will be cached in a private field. Each time a call to function ToSQL
	// happens, the arguments will be re-evaluated and returned or interpolated.
	UseBuildCache bool
	cacheSQL      []byte
	cacheArgs     Arguments // like a buffer, gets reused

	RawFullSQL   string
	RawArguments Arguments // Arguments used by RawFullSQL

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

	// OnDuplicateKey updates the referenced columns. See documentation for type
	// `UpdatedColumns`. For more details
	// https://dev.mysql.com/doc/refman/5.7/en/insert-on-duplicate.html
	OnDuplicateKey UpdatedColumns
	// OnDuplicateKeyAvoidUpdate skips adding the mentioned columns to the ON
	// DUPLICATE KEY UPDATE section. Otherwise all columns in the field
	// `Columns` will be added to the ON DUPLICATE KEY UPDATE expression.
	// Usually the slice `OnDuplicateKeyAvoidUpdate` contains the primary key
	// columns. Case-sensitive comparison.
	OnDuplicateKeyAvoidUpdate []string
	// IsReplace uses the REPLACE syntax. See function Replace().
	IsReplace bool
	// IsIgnore ignores error. See function Ignore().
	IsIgnore      bool
	IsInterpolate bool // See Interpolate()
	// PropagationStopped set to true if you would like to interrupt the
	// listener chain. Once set to true all sub sequent calls of the next
	// listeners will be suppressed.
	PropagationStopped bool
	// Listeners allows to dispatch certain functions in different
	// situations.
	Listeners InsertListeners
	// propagationStoppedAt position in the slice where the stopped propagation
	// has been requested. for every new iteration the propagation must stop at
	// this position.
	propagationStoppedAt int
}

// NewInsert creates a new Insert object.
func NewInsert(into string) *Insert {
	return &Insert{
		Into: into,
	}
}

// InsertInto instantiates a Insert for the given table
func (c *Connection) InsertInto(into string) *Insert {
	i := &Insert{
		Log:  c.Log,
		Into: into,
	}
	i.DB = c.DB
	return i
}

// InsertInto instantiates a Insert for the given table bound to a transaction
func (tx *Tx) InsertInto(into string) *Insert {
	i := &Insert{
		Log:  tx.Logger,
		Into: into,
	}
	i.DB = tx.Tx
	return i
}

// WithDB sets the database query object.
func (b *Insert) WithDB(db ExecPreparer) *Insert {
	b.DB = db
	return b
}

// AddOnDuplicateKeyAvoidUpdate skips adding the mentioned columns to the ON
// DUPLICATE KEY UPDATE section. Otherwise all columns in the field
// `Columns` will be added to the ON DUPLICATE KEY UPDATE expression.
// Usually the slice `OnDuplicateKeyAvoidUpdate` contains the primary key
// columns. Case-sensitive comparison.
func (b *Insert) AddOnDuplicateKeyAvoidUpdate(primaryKeyColumnNames ...string) *Insert {
	b.OnDuplicateKeyAvoidUpdate = append(b.OnDuplicateKeyAvoidUpdate, primaryKeyColumnNames...)
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

// AddValues appends a set of values to the statement. Each call of AddValues
// creates a new set of values. Only primitive types are supported. Runtime type
// safety only.
func (b *Insert) AddValues(values ...interface{}) *Insert {
	return b.AddArguments(iFaceToArgs(values...)...)
}

// SetRowCount defines the number of expected rows. Each set of place holders
// within the brackets defines a row. This setting defaults to one. It gets
// applied when fields `Arguments` and `Records` have been left empty. For each
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
func (b *Insert) AddArguments(args ...Argument) *Insert {
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

// AddRecords appends record generators.
func (b *Insert) AddRecords(recs ...ArgumentsAppender) *Insert {
	b.Records = append(b.Records, recs...)
	return b
}

// SetRecordValueCount number of expected values within each set. Must be
// applied if a call to AddColumns has been omitted and AddRecords gets called
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
func (b *Insert) AddOnDuplicateKey(column string, arg Argument) *Insert {
	b.OnDuplicateKey.Columns = append(b.OnDuplicateKey.Columns, column)
	b.OnDuplicateKey.Arguments = append(b.OnDuplicateKey.Arguments, arg)
	return b
}

// Pair appends a key/value (column/value) pair to the statement. Calling this
// function multiple times with the same column name produces invalid SQL.
func (b *Insert) Pair(column string, arg Argument) *Insert {
	colPos := -1
	for i, c := range b.Columns {
		if strings.EqualFold(c, column) {
			colPos = i
			break
		}
	}
	if colPos == -1 {
		b.Columns = append(b.Columns, column)
		if len(b.Values) == 0 {
			b.Values = make([]Arguments, 1, 5)
		}
		b.Values[0] = append(b.Values[0], arg)
		return b
	}

	if colPos == 0 { // create new slice
		b.Values = append(b.Values, make(Arguments, len(b.Columns)))
	}
	pos := len(b.Values) - 1
	b.Values[pos][colPos] = arg
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
// prepared statements. ToSQLs second argument `Arguments` will then be nil.
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

func (b *Insert) hasBuildCache() bool {
	return b.UseBuildCache
}

func (b *Insert) toSQL(buf queryWriter) error {

	if err := b.Listeners.dispatch(OnBeforeToSQL, b); err != nil {
		return errors.WithStack(err)
	}

	if len(b.OnDuplicateKeyAvoidUpdate) > 0 {
		for _, c := range b.Columns {
			// Wow two times a comparison with a slice. That costs a bit
			// performance but a reliable way to avoid writing duplicate ON
			// DUPLICATE KEY UPDATE sets. If there is something faster, write us.
			if !strInSlice(c, b.OnDuplicateKeyAvoidUpdate) && !strInSlice(c, b.OnDuplicateKey.Columns) {
				b.OnDuplicateKey.Columns = append(b.OnDuplicateKey.Columns, c)
			}
		}
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
	Quoter.writeName(buf, b.Into)
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
			Quoter.writeName(buf, c)
		}
		buf.WriteString(") ")
	}
	buf.WriteString("VALUES ")

	rowCount := 1
	if b.RowCount > 0 {
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
		rowCount = len(b.Records)
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
			buf.WriteByte('?')
		}
		buf.WriteByte(')')
	}
	return errors.Wrap(b.OnDuplicateKey.writeOnDuplicateKey(buf), "[dbr] Insert.toSQL.writeOnDuplicateKey\n")
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
		args = make(Arguments, 0, totalArgCount+len(b.Records)+len(b.OnDuplicateKey.Columns)) // sneaky ;-)
	}
	for _, v := range b.Values {
		args = append(args, v...)
	}

	if b.Records == nil {
		args, err = b.OnDuplicateKey.appendArgs(args)
		return args, errors.WithStack(err)
	}

	for _, rec := range b.Records {
		alBefore := len(args)
		args, err = rec.AppendArguments(SQLStmtInsert|SQLPartValues, args, b.Columns) // Columns can be empty
		if err != nil {
			return nil, errors.WithStack(err)
		}
		if addedArgs := len(args) - alBefore; addedArgs != argCount0 {
			return nil, errors.NewMismatchf("[dbr] Insert.appendArgs RecordValueCount(%d) does not match the number of assembled arguments (%d)", b.RecordValueCount, addedArgs)
		}
	}

	if args, err = b.OnDuplicateKey.appendArgs(args); err != nil {
		return nil, errors.WithStack(err)
	}
	return args, nil
}

// Exec executes the statement represented by the Insert object.
// It returns the raw database/sql Result and an error if there was one.
// Regarding LastInsertID(): If you insert multiple rows using a single
// INSERT statement, LAST_INSERT_ID() returns the value generated for
// the first inserted row only. The reason for this at to make it possible to
// reproduce easily the same INSERT statement against some other server.
func (b *Insert) Exec(ctx context.Context) (sql.Result, error) {
	result, err := Exec(ctx, b.DB, b)
	return result, errors.WithStack(err)
}

// Prepare creates a prepared statement
func (b *Insert) Prepare(ctx context.Context) (*sql.Stmt, error) {
	stmt, err := Prepare(ctx, b.DB, b)
	return stmt, errors.WithStack(err)
}
