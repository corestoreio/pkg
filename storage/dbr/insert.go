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

	"github.com/corestoreio/csfw/util/bufferpool"
	"github.com/corestoreio/errors"
	"github.com/corestoreio/log"
)

// Insert contains the clauses for an INSERT statement
type Insert struct {
	Log log.Logger // Log optional logger
	DB  Execer

	// UseBuildCache if set to true the final build query will be stored in
	// field private field `buildCache` and the arguments in field `Arguments`
	UseBuildCache bool
	buildCache    []byte

	RawFullSQL   string
	RawArguments Arguments // Arguments used by RawFullSQL or BuildCache

	Into    string
	Columns []string
	Values  []Arguments

	Records []ArgumentAssembler
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
	// previousError any error occurred during construction the SQL statement
	previousError error
}

// NewInsert creates a new object with a black hole logger.
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
func (b *Insert) WithDB(db Execer) *Insert {
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

// AddColumns appends columns to insert in the statement.
func (b *Insert) AddColumns(columns ...string) *Insert {
	b.RecordValueCount += len(columns)
	b.Columns = append(b.Columns, columns...)
	return b
}

// AddValues appends a set of values to the statement. Each call of AddValues
// creates a new set of values. Only primitive types are supported. Runtime type
// safety only.
func (b *Insert) AddValues(values ...interface{}) *Insert {
	args, err := iFaceToArgs(values...)
	if err != nil {
		b.previousError = errors.Wrap(err, "[dbr] Insert.AddValues.iFaceToArgs")
		return b
	}
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

// AddRecords pulls in values to match Columns from the record generator.
func (b *Insert) AddRecords(recs ...ArgumentAssembler) *Insert {
	b.Records = append(b.Records, recs...)
	return b
}

// SetRecordValueCount number of expected values within each set. Must be
// applied if columns have been omitted and AddRecords gets called or Records
// gets set in a different way.
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

// Pair adds a key/value pair to the statement.
func (b *Insert) Pair(column string, arg Argument) *Insert {
	if b.previousError != nil {
		return b
	}
	for _, c := range b.Columns {
		if c == column {
			b.previousError = errors.NewAlreadyExistsf("[dbr] Column %q has already been added", c)
			return b
		}
	}

	b.Columns = append(b.Columns, column)
	if len(b.Values) == 0 {
		b.Values = make([]Arguments, 1, 5)
	}
	b.Values[0] = append(b.Values[0], arg)

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
func (b *Insert) ToSQL() (string, Arguments, error) {
	return toSQL(b, b.IsInterpolate)
}

func (b *Insert) writeBuildCache(sql []byte, arguments Arguments) {
	b.buildCache = sql
	b.RawArguments = arguments
}

func (b *Insert) readBuildCache() (sql []byte, arguments Arguments) {
	return b.buildCache, b.RawArguments
}

func (b *Insert) hasBuildCache() bool {
	return b.UseBuildCache
}

func (b *Insert) toSQL(buf queryWriter) error {
	if b.previousError != nil {
		return errors.Wrap(b.previousError, "[dbr] Insert.ToSQL")
	}

	if err := b.Listeners.dispatch(OnBeforeToSQL, b); err != nil {
		return errors.Wrap(err, "[dbr] Insert.Listeners.dispatch")
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
		return errors.Wrap(b.Select.toSQL(buf), "[dbr] Insert.FromSelect")
	}

	if lv := len(b.Values); b.Records == nil && (lv == 0 || (lv > 0 && len(b.Values[0]) == 0)) {
		return errors.NewEmptyf("[dbr] Insert.ToSQL cannot find any Values for table %q", b.Into)
	}

	var ph = bufferpool.Get() // Build the ph like "(?,?,?)"
	defer bufferpool.Put(ph)

	if len(b.Columns) > 0 {
		ph.WriteByte('(')
		buf.WriteByte('(')
		for i, c := range b.Columns {
			if i > 0 {
				buf.WriteByte(',')
				ph.WriteByte(',')
			}
			Quoter.FquoteAs(buf, c)
			ph.WriteByte('?')
		}
		ph.WriteByte(')')
		buf.WriteByte(')')
		buf.WriteByte(' ')
	}
	buf.WriteString("VALUES ")

	var argCount0 int
	if len(b.Values) > 0 && len(b.Columns) == 0 {
		argCount0 = len(b.Values[0])
		// no columns provided so build the place holders.
		ph.WriteByte('(')
		for i := 0; i < argCount0; i++ {
			if i > 0 {
				ph.WriteByte(',')
			}
			ph.WriteByte('?')
		}
		ph.WriteByte(')')
	}

	placeholderStr := ph.String()

	for i := range b.Values {
		if i > 0 {
			buf.WriteByte(',')
		}
		buf.WriteString(placeholderStr)
	}

	if b.Records == nil {
		return errors.Wrap(b.OnDuplicateKey.writeOnDuplicateKey(buf), "[dbr] Insert.OnDuplicateKey.writeOnDuplicateKey")
	}

	for i := range b.Records {
		if i == 0 && placeholderStr == "" {

			// Build place holder string because here the func knows how many arguments it have.
			ph.WriteByte('(')
			for j := 0; j < b.RecordValueCount; j++ {
				if j > 0 {
					ph.WriteByte(',')
				}
				ph.WriteByte('?')
			}
			ph.WriteByte(')')
			placeholderStr = ph.String()
		}
		if i > 0 {
			buf.WriteByte(',')
		}
		buf.WriteString(placeholderStr)
	}

	if err := b.OnDuplicateKey.writeOnDuplicateKey(buf); err != nil {
		return errors.Wrap(err, "[dbr] Insert.OnDuplicateKey.writeOnDuplicateKey")
	}
	return nil
}

func (b *Insert) appendArgs(args Arguments) (_ Arguments, err error) {
	if b.previousError != nil {
		return nil, errors.Wrap(b.previousError, "[dbr] Insert.ToSQL")
	}

	if b.RawFullSQL != "" {
		return b.RawArguments, nil
	}

	if b.Select != nil {
		args, err = b.Select.appendArgs(args)
		return args, errors.Wrap(err, "[dbr] Insert.FromSelect")
	}

	if lv := len(b.Values); b.Records == nil && (lv == 0 || (lv > 0 && len(b.Values[0]) == 0)) {
		return nil, errors.NewEmptyf("[dbr] Insert.ToSQL cannot find any Values for table %q", b.Into)
	}

	argCount0 := b.RecordValueCount
	if len(b.Values) > 0 && len(b.Columns) == 0 {
		argCount0 = len(b.Values[0])
	}
	if lc := len(b.Columns); argCount0 < 1 && lc > 0 {
		argCount0 = lc
	}

	totalArgCount := len(b.Values) * argCount0
	if args == nil {
		args = make(Arguments, 0, totalArgCount+len(b.Records)+len(b.OnDuplicateKey.Columns)) // sneaky ;-)
	}
	for _, v := range b.Values {
		args = append(args, v...)
	}

	if b.Records == nil {
		args, err = b.OnDuplicateKey.appendArgs(args)
		return args, errors.Wrap(err, "[dbr] Insert.OnDuplicateKey.appendArgs")
	}

	for _, rec := range b.Records {
		alBefore := len(args)
		args, err = rec.AssembleArguments(SQLStmtInsert|SQLPartValues, args, b.Columns) // Columns can be empty
		if err != nil {
			return nil, errors.Wrap(err, "[dbr] Insert.ToSQL.Record")
		}
		if addedArgs := len(args) - alBefore; addedArgs != argCount0 {
			return nil, errors.NewMismatchf("[dbr] Insert.appendArgs RecordValueCount(%d) does not match the number of assembled arguments (%d)", b.RecordValueCount, addedArgs)
		}
	}

	if args, err = b.OnDuplicateKey.appendArgs(args); err != nil {
		return nil, errors.Wrap(err, "[dbr] Insert.OnDuplicateKey.appendArgs")
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
	sqlStr, args, err := b.ToSQL()
	if err != nil {
		return nil, errors.Wrap(err, "[dbr] Insert.Exec.ToSQL")
	}

	if b.Log != nil && b.Log.IsInfo() {
		defer log.WhenDone(b.Log).Info("dbr.Insert.Exec.Timing", log.String("sqlStr", sqlStr))
	}
	result, err := b.DB.ExecContext(ctx, sqlStr, args.Interfaces()...)
	if err != nil {
		return result, errors.Wrap(err, "[dbr] Insert.Exec.Exec")
	}

	return result, nil
}

// Prepare creates a prepared statement
func (b *Insert) Prepare(ctx context.Context) (*sql.Stmt, error) {
	rawSQL, _, err := b.ToSQL() // TODO create a ToSQL version without any arguments
	if err != nil {
		return nil, errors.Wrap(err, "[dbr] Insert.Exec.ToSQL")
	}

	if b.Log != nil && b.Log.IsInfo() {
		defer log.WhenDone(b.Log).Info("dbr.Insert.Prepare.Timing", log.String("sql", rawSQL))
	}

	stmt, err := b.DB.PrepareContext(ctx, rawSQL)
	return stmt, errors.Wrap(err, "[dbr] Insert.Prepare.Prepare")
}
