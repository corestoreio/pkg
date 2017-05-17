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
	DB  struct {
		Preparer
		Execer
	}
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
	// Select used to create an "INSERT INTO `table` SELECT ..." statement.
	Select *Select
	Maps   map[string]Argument

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
	i.DB.Execer = c.DB
	i.DB.Preparer = c.DB
	return i
}

// InsertInto instantiates a Insert for the given table bound to a transaction
func (tx *Tx) InsertInto(into string) *Insert {
	i := &Insert{
		Log:  tx.Logger,
		Into: into,
	}
	i.DB.Execer = tx.Tx
	i.DB.Preparer = tx.Tx
	return i
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
	b.Columns = append(b.Columns, columns...)
	return b
}

// AddValues appends a set of values to the statement. Each call of AddValues
// creates a new set of values.
// TODO(CyS) refactor and maybe use an empty interface slice ...
// Deprecated: Gets refactored
func (b *Insert) AddValues(vals ...Argument) *Insert {
	if lv, mod := len(vals), len(b.Columns); mod > 0 && lv > mod && (lv%mod) == 0 {
		// now we have more arguments than columns and we can assume that more
		// rows gets inserted.
		for i := 0; i < len(vals); i = i + mod {
			b.Values = append(b.Values, vals[i:i+mod])
		}
	} else {
		// each call to AddValues equals one row in a table.
		b.Values = append(b.Values, vals)
	}
	return b
}

// AddRecords pulls in values to match Columns from the record generator.
func (b *Insert) AddRecords(recs ...ArgumentAssembler) *Insert {
	b.Records = append(b.Records, recs...)
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

// Map pulls in values to match Columns from the record. Calling multiple
// times will add new map entries to the Insert map.
func (b *Insert) Map(m map[string]Argument) *Insert {
	if b.Maps == nil {
		b.Maps = make(map[string]Argument)
	}
	for col, val := range m {
		b.Maps[col] = val
	}
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

func (b *Insert) toSQL(buf queryWriter) (Arguments, error) {
	if b.previousError != nil {
		return nil, errors.Wrap(b.previousError, "[dbr] Insert.ToSQL")
	}

	if err := b.Listeners.dispatch(OnBeforeToSQL, b); err != nil {
		return nil, errors.Wrap(err, "[dbr] Insert.Listeners.dispatch")
	}

	if b.RawFullSQL != "" {
		buf.WriteString(b.RawFullSQL)
		return b.RawArguments, nil
	}

	if len(b.Into) == 0 {
		return nil, errors.NewEmptyf(errTableMissing)
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
		sArgs, err := b.Select.toSQL(buf)
		return sArgs, errors.Wrap(err, "[dbr] Insert.FromSelect")
	}

	if len(b.Maps) > 0 {
		args, err := b.mapToSQL(buf)
		return args, errors.Wrap(err, "[dbr] Insert.mapToSQL")
	}

	if lv := len(b.Values); b.Records == nil && (lv == 0 || (lv > 0 && len(b.Values[0]) == 0)) {
		return nil, errors.NewEmptyf("[dbr] Insert.ToSQL cannot find any Values for table %q", b.Into)
	}

	var ph = bufferpool.Get() // Build the ph like "(?,?,?)"

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
	} else {
		// no columns provided so build the place holders.
		ph.WriteByte('(')
		for i := range b.Values[0] {
			if i > 0 {
				ph.WriteByte(',')
			}
			ph.WriteByte('?')
		}
		ph.WriteByte(')')
	}
	buf.WriteString("VALUES ")

	placeholderStr := ph.String()
	bufferpool.Put(ph)

	var argCount0 int
	if len(b.Values) > 0 {
		argCount0 = len(b.Values[0])
	}
	totalArgCount := len(b.Values) * argCount0
	args := make(Arguments, 0, totalArgCount+len(b.Records)+len(b.OnDuplicateKey.Columns)) // sneaky ;-)
	for i, v := range b.Values {
		if i > 0 {
			buf.WriteByte(',')
		}
		buf.WriteString(placeholderStr)
		args = append(args, v...)
	}

	var err error
	if b.Records == nil {
		args, err = b.OnDuplicateKey.writeOnDuplicateKey(buf, args)
		return args, errors.Wrap(err, "[dbr] Insert.OnDuplicateKey.writeOnDuplicateKey")
	}

	for i, rec := range b.Records {
		args, err = rec.AssembleArguments(stmtTypeInsert, args, b.Columns, nil)
		if err != nil {
			return nil, errors.Wrap(err, "[dbr] Insert.ToSQL.Record")
		}
		if i > 0 {
			buf.WriteByte(',')
		}
		buf.WriteString(placeholderStr)
	}

	if args, err = b.OnDuplicateKey.writeOnDuplicateKey(buf, args); err != nil {
		return nil, errors.Wrap(err, "[dbr] Insert.OnDuplicateKey.writeOnDuplicateKey")
	}
	return args, nil
}

// mapToSQL serialized the Insert to a SQL string
// It goes through the Maps param and combined its keys/values into the SQL query string
// It returns the string with placeholders and a slice of query arguments
func (b *Insert) mapToSQL(w queryWriter) (Arguments, error) {
	if b.previousError != nil {
		return nil, errors.Wrap(b.previousError, "[dbr] Insert.ToSQL")
	}

	keys := make([]string, len(b.Maps))
	vals := make(Arguments, len(b.Maps))
	i := 0
	for k, v := range b.Maps {
		keys[i] = k
		vals[i] = v
		i++
	}
	var args Arguments
	var placeholder = bufferpool.Get() // Build the placeholder like "(?,?,?)"
	defer bufferpool.Put(placeholder)

	placeholder.WriteByte('(')
	for i, c := range keys {
		if i > 0 {
			w.WriteByte(',')
			placeholder.WriteByte(',')
		}
		Quoter.FquoteAs(w, c)
		placeholder.WriteByte('?')
	}
	w.WriteString(") VALUES ")
	placeholder.WriteByte(')')
	w.WriteString(placeholder.String())

	args = append(args, vals...)

	return args, nil
}

// Exec executes the statement represented by the Insert
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
