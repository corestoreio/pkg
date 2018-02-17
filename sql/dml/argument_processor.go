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
	"database/sql/driver"
	"sync"
	"time"

	"github.com/corestoreio/errors"
	"github.com/corestoreio/log"
	"github.com/corestoreio/pkg/util/bufferpool"
)

// Arguments a collection of primitive types or slices of primitive types.
// It acts as some kind of prepared statement.
type Arguments struct {
	base builderCommon
	// QualifiedColumnsAliases allows to overwrite the internal qualified
	// columns slice with custom names. Only in the use case when records are
	// applied. The list of column names in `QualifiedColumnsAliases` gets
	// passed to the ColumnMapper and back to the provided object. The
	// `QualifiedColumnsAliases` slice must have the same length as the
	// qualified columns slice. The order of the alias names must be in the same
	// order as the qualified columns or as the placeholders occur.
	QualifiedColumnsAliases []string
	// insertCachedSQL contains the final build SQL string with the correct
	// amount of placeholders.
	insertCachedSQL     []byte
	insertColumnCount   uint
	insertRowCount      uint
	insertIsBuildValues bool
	// isPrepared if true the cachedSQL field in base gets ignored
	isPrepared bool
	Options    uint
	// hasNamedArgs checks if the SQL string in the cachedSQL field contains
	// named arguments. 0 not yet checked, 1=does not contain, 2 = yes
	hasNamedArgs      uint8 // 0 not checked, 1=no, 2=yes
	nextUnnamedArgPos int
	raw               []interface{}
	arguments
	recs []QualifiedRecord
}

const (
	argOptionExpandPlaceholder = 1 << iota
	argOptionInterpolate
)

// WithQualifiedColumnsAliases for documentation please see:
// Arguments.QualifiedColumnsAliases.
func (a *Arguments) WithQualifiedColumnsAliases(aliases ...string) *Arguments {
	a.QualifiedColumnsAliases = aliases
	return a
}

// ToSQL the returned interface slice is owned by the callee.
func (a *Arguments) ToSQL() (string, []interface{}, error) {
	return a.prepareArgs()
}

// Interpolate if set stringyfies the arguments into the SQL string and returns
// pre-processed SQL command when calling the function ToSQL. Not suitable for
// prepared statements. ToSQLs second argument `args` will then be nil.
func (a *Arguments) Interpolate() *Arguments {
	a.Options = a.Options | argOptionInterpolate
	return a
}

// ExpandPlaceHolders repeats the place holders with the provided argument
// count. If the amount of arguments does not match the number of place holders,
// a mismatch error gets returned.
//		ExpandPlaceHolders("SELECT * FROM table WHERE id IN (?) AND status IN (?)", Int(myIntSlice...), String(myStrSlice...))
// Gets converted to:
//		SELECT * FROM table WHERE id IN (?,?) AND status IN (?,?,?)
// The place holders are of course depending on the values in the Arg*
// functions. This function should be generally used when dealing with prepared
// statements or interpolation.
func (a *Arguments) ExpandPlaceHolders() *Arguments {
	a.Options = a.Options | argOptionExpandPlaceholder
	return a
}

func (a *Arguments) isEmpty() bool {
	if a == nil {
		return true
	}
	return len(a.raw) == 0 && len(a.arguments) == 0 && len(a.recs) == 0
}

// prepareArgs transforms mainly the Arguments into []interface{}. It appends
// its arguments to the `extArgs` arguments from the Exec+ or Query+ function.
// This allows for a developer to reuse the interface slice and save
// allocations. All method receivers are not thread safe. The returned interface
// slice is the same as `extArgs`.
func (a *Arguments) prepareArgs(extArgs ...interface{}) (_ string, _ []interface{}, err error) {
	if a.base.ärgErr != nil {
		return "", nil, errors.WithStack(a.base.ärgErr)
	}
	if !a.isPrepared && len(a.base.cachedSQL) == 0 {
		return "", nil, errors.Empty.Newf("[dml] Arguments: The SQL string is empty.")
	}

	if a.base.source == dmlSourceInsert {
		return a.prepareArgsInsert(extArgs...)
	}

	if a.isEmpty() {
		a.hasNamedArgs = 1
		if a.isPrepared {
			return "", extArgs, nil
		}
		return string(a.base.cachedSQL), extArgs, nil
	}

	if !a.isPrepared && a.hasNamedArgs == 0 {
		var found bool
		a.hasNamedArgs = 1
		a.base.cachedSQL, a.base.qualifiedColumns, found = extractReplaceNamedArgs(a.base.cachedSQL, a.base.qualifiedColumns)

		switch la := len(a.arguments); true {
		case found:
			a.hasNamedArgs = 2
		case !found && len(a.recs) == 0 && la > 0:
			for _, arg := range a.arguments {
				if arg.name != "" {
					a.hasNamedArgs = 2
					break
				}
			}
		}
	}

	sqlBuf := bufferpool.GetTwin()
	collectedArgs := pooledArgumentsGet()
	defer pooledArgumentsPut(collectedArgs, sqlBuf)
	collectedArgs = append(collectedArgs, a.arguments...)

	extArgs = append(extArgs, a.raw...)
	if collectedArgs, err = a.appendConvertedRecordsToArguments(collectedArgs); err != nil {
		return "", nil, errors.WithStack(err)
	}

	if a.isPrepared {
		return "", collectedArgs.Interfaces(extArgs...), nil
	}

	// Make a copy of the original SQL statement because it gets modified in the
	// worst case. Best case would be no modification and hence we don't need a
	// bytes.Buffer from the pool! TODO(CYS) optimize this and only acquire a
	// buffer from the pool in the worse case.
	if _, err := sqlBuf.First.Write(a.base.cachedSQL); err != nil {
		return "", nil, errors.WithStack(err)
	}

	// `switch` statement no suitable.
	if a.Options > 0 && len(extArgs) > 0 && len(a.recs) == 0 && len(a.arguments) == 0 {
		return "", nil, errors.NotAllowed.Newf("[dml] Interpolation/ExpandPlaceholders supports only Records and Arguments and not yet an interface slice.")
	}

	if a.Options&argOptionExpandPlaceholder != 0 {
		if phCount := bytes.Count(sqlBuf.First.Bytes(), placeHolderByte); phCount < a.Len() {
			if err := expandPlaceHolders(sqlBuf.Second, sqlBuf.First.Bytes(), collectedArgs); err != nil {
				return "", nil, errors.WithStack(err)
			}
			if _, err := sqlBuf.CopySecondToFirst(); err != nil {
				return "", nil, errors.WithStack(err)
			}
		}
	}
	if a.Options&argOptionInterpolate != 0 {
		if err := writeInterpolateBytes(sqlBuf.Second, sqlBuf.First.Bytes(), collectedArgs); err != nil {
			return "", nil, errors.Wrapf(err, "[dml] Interpolation failed: %q", sqlBuf.String())
		}
		return sqlBuf.Second.String(), nil, nil
	}

	return sqlBuf.First.String(), collectedArgs.Interfaces(extArgs...), nil
}

func (a *Arguments) appendConvertedRecordsToArguments(collectedArgs arguments) (arguments, error) {
	if a.base.templateStmtCount == 0 {
		a.base.templateStmtCount = 1
	}
	if len(a.arguments) == 0 && len(a.recs) == 0 {
		return collectedArgs, nil
	}

	if len(a.arguments) > 0 && len(a.recs) == 0 && a.hasNamedArgs < 2 {
		if a.base.templateStmtCount > 1 {
			collectedArgs = a.multiplyArguments(a.base.templateStmtCount)
		}
		// This is also a case where there are no records and only arguments and
		// those arguments do not contain any name. Then we can skip the column
		// mapper and ignore the qualifiedColumns.
		return collectedArgs, nil
	}

	qualifiedColumns := a.base.qualifiedColumns
	if lqca := len(a.QualifiedColumnsAliases); lqca > 0 {
		if lqca != len(a.base.qualifiedColumns) {
			return nil, errors.Mismatch.Newf("[dml] Argument.Record: QualifiedColumnsAliases slice %v and qualifiedColumns slice %v must have the same length", a.QualifiedColumnsAliases, a.base.qualifiedColumns)
		}
		qualifiedColumns = a.QualifiedColumnsAliases
	}

	// TODO refactor prototype and make it performant and beautiful code
	cm := NewColumnMap(len(a.arguments)+len(a.recs), "") // can use an arg pool Arguments sync.Pool, nope.

	for tsc := 0; tsc < a.base.templateStmtCount; tsc++ { // only in case of UNION statements in combination with a template SELECT, can be optimized later

		// `qualifiedColumns` contains the correct order as the place holders
		// appear in the SQL string.
		for _, identifier := range qualifiedColumns {
			// identifier can be either: column or qualifier.column or :column
			qualifier, column := splitColumn(identifier)
			// a.base.defaultQualifier is empty in case of INSERT statements

			column, isNamedArg := cutNamedArgStartStr(column) // removes the colon for named arguments
			cm.columns[0] = column                            // length is always one, as created in NewColumnMap

			if isNamedArg && len(a.arguments) > 0 {
				// if the colon : cannot be found then a simple place holder ? has been detected
				if err := a.MapColumns(cm); err != nil {
					return collectedArgs, errors.WithStack(err)
				}
			} else {
				found := false
				for _, qRec := range a.recs {
					if qRec.Qualifier == "" && qualifier != "" {
						qRec.Qualifier = a.base.defaultQualifier
					}
					if qRec.Qualifier != "" && qualifier == "" {
						qualifier = a.base.defaultQualifier
					}

					if qRec.Qualifier == qualifier {
						if err := qRec.Record.MapColumns(cm); err != nil {
							return collectedArgs, errors.WithStack(err)
						}
						found = true
					}
				}
				if !found {
					// If the argument cannot be found in the records then we assume the argument
					// has a numerical position and we grab just the next unnamed argument.
					if pArg, ok := a.nextUnnamedArg(); ok {
						cm.arguments = append(cm.arguments, pArg)
					}
				}
			}
		}
		a.nextUnnamedArgPos = 0
	}
	if len(cm.arguments) > 0 {
		collectedArgs = cm.arguments
	}

	return collectedArgs, nil
}

// prepareArgsInsert prepares the special arguments for an INSERT statement. The
// returned interface slice is the same as the `extArgs` slice. extArgs =
// external arguments.
func (a *Arguments) prepareArgsInsert(extArgs ...interface{}) (string, []interface{}, error) {

	//cm := pooledColumnMapGet()
	sqlBuf := bufferpool.GetTwin()
	defer bufferpool.PutTwin(sqlBuf)
	//defer pooledBufferColumnMapPut(cm, sqlBuf, nil)

	cm := NewColumnMap(16)
	cm.setColumns(a.base.qualifiedColumns)
	//defer bufferpool.PutTwin(sqlBuf)
	cm.arguments = append(cm.arguments, a.arguments...)
	lenInsertCachedSQL := len(a.insertCachedSQL)
	{
		cachedSQL := a.base.cachedSQL
		if lenInsertCachedSQL > 0 {
			cachedSQL = a.insertCachedSQL
		}
		if _, err := sqlBuf.First.Write(cachedSQL); err != nil {
			return "", nil, errors.WithStack(err)
		}

		for _, qRec := range a.recs {
			if qRec.Qualifier != "" {
				return "", nil, errors.Fatal.Newf("[dml] Qualifier in %T is not supported and not needed.", qRec)
			}

			if err := qRec.Record.MapColumns(cm); err != nil {
				return "", nil, errors.WithStack(err)
			}
		}
	}

	if a.isPrepared {
		// TODO above construct can be more optimized when using prepared statements
		return "", cm.arguments.Interfaces(extArgs...), nil
	}

	extArgs = append(extArgs, a.raw...)
	totalArgLen := uint(len(cm.arguments) + len(extArgs))

	if !a.insertIsBuildValues && lenInsertCachedSQL == 0 { // Write placeholder list e.g. "VALUES (?,?),(?,?)"
		odkPos := bytes.Index(a.base.cachedSQL, onDuplicateKeyPart)
		if odkPos > 0 {
			sqlBuf.First.Reset()
			sqlBuf.First.Write(a.base.cachedSQL[:odkPos])
		}

		if a.insertRowCount > 0 {
			columnCount := totalArgLen / a.insertRowCount
			writeInsertPlaceholders(sqlBuf.First, a.insertRowCount, columnCount)
		} else if a.insertColumnCount > 0 {
			rowCount := totalArgLen / a.insertColumnCount
			if rowCount == 0 {
				rowCount = 1
			}
			writeInsertPlaceholders(sqlBuf.First, rowCount, a.insertColumnCount)
		}
		if odkPos > 0 {
			sqlBuf.First.Write(a.base.cachedSQL[odkPos:])
		}
		a.insertCachedSQL = bufTrySizeByResliceOrNew(a.insertCachedSQL, sqlBuf.First.Len())
		copy(a.insertCachedSQL, sqlBuf.First.Bytes())
	}

	if a.Options > 0 {
		if len(extArgs) > 0 && len(a.recs) == 0 && len(cm.arguments) == 0 {
			return "", nil, errors.NotAllowed.Newf("[dml] Interpolation/ExpandPlaceholders supports only Records and Arguments and not yet an interface slice.")
		}

		if a.Options&argOptionInterpolate != 0 {
			if err := writeInterpolateBytes(sqlBuf.Second, sqlBuf.First.Bytes(), cm.arguments); err != nil {
				return "", nil, errors.Wrapf(err, "[dml] Interpolation failed: %q", sqlBuf.First.String())
			}
			return sqlBuf.Second.String(), nil, nil
		}
	}

	return sqlBuf.First.String(), cm.arguments.Interfaces(extArgs...), nil
}

// nextUnnamedArg returns an unnamed argument by its position.
func (a *Arguments) nextUnnamedArg() (argument, bool) {
	var unnamedCounter int
	lenArg := len(a.arguments)
	for i := 0; i < lenArg && a.nextUnnamedArgPos >= 0; i++ {
		if arg := a.arguments[i]; arg.name == "" {
			if unnamedCounter == a.nextUnnamedArgPos {
				a.nextUnnamedArgPos++
				return arg, true
			}
			unnamedCounter++
		}
	}
	a.nextUnnamedArgPos = -1 // nothing found, so no need to further iterate through the []argument slice.
	return argument{}, false
}

// MapColumns allows to merge one argument slice with another depending on the
// matched columns. Each argument in the slice must be a named argument.
// Implements interface ColumnMapper.
func (a *Arguments) MapColumns(cm *ColumnMap) error {
	if cm.Mode() == ColumnMapEntityReadAll {
		cm.arguments = append(cm.arguments, a.arguments...)
		return cm.Err()
	}
	for cm.Next() {
		// now a bit slow ... but will be refactored later with constant time
		// access, but first benchmark it. This for loop can be the 3rd one in the
		// overall chain.
		c := cm.Column()
		for _, arg := range a.arguments {
			// Case sensitive comparison
			if c != "" && arg.name == c {
				cm.arguments = append(cm.arguments, arg)
				break
			}
		}
	}
	return cm.Err()
}

func (a *Arguments) add(v interface{}) *Arguments {
	if a == nil {
		a = MakeArgs(defaultArgumentsCapacity)
	}
	a.arguments = a.arguments.add(v)
	return a
}

func (a *Arguments) Record(qualifier string, record ColumnMapper) *Arguments {
	a.recs = append(a.recs, Qualify(qualifier, record))
	return a
}

// Arguments sets the internal arguments slice to the provided argument. Those
// are the slices Arguments, records and raw.
func (a *Arguments) Arguments(args *Arguments) *Arguments {
	// maybe deprecated this function.
	a.arguments = args.arguments
	a.recs = args.recs
	a.raw = args.raw
	return a
}

func (a *Arguments) Records(records ...QualifiedRecord) *Arguments { a.recs = records; return a }
func (a *Arguments) Raw(raw ...interface{}) *Arguments             { a.raw = raw; return a }

func (a *Arguments) Null() *Arguments                          { return a.add(nil) }
func (a *Arguments) Unsafe(arg interface{}) *Arguments         { return a.add(arg) }
func (a *Arguments) Int(i int) *Arguments                      { return a.add(i) }
func (a *Arguments) Ints(i ...int) *Arguments                  { return a.add(i) }
func (a *Arguments) Int64(i int64) *Arguments                  { return a.add(i) }
func (a *Arguments) Int64s(i ...int64) *Arguments              { return a.add(i) }
func (a *Arguments) Uint(i uint) *Arguments                    { return a.add(uint64(i)) }
func (a *Arguments) Uints(i ...uint) *Arguments                { return a.add(i) }
func (a *Arguments) Uint64(i uint64) *Arguments                { return a.add(i) }
func (a *Arguments) Uint64s(i ...uint64) *Arguments            { return a.add(i) }
func (a *Arguments) Float64(f float64) *Arguments              { return a.add(f) }
func (a *Arguments) Float64s(f ...float64) *Arguments          { return a.add(f) }
func (a *Arguments) Bool(b bool) *Arguments                    { return a.add(b) }
func (a *Arguments) Bools(b ...bool) *Arguments                { return a.add(b) }
func (a *Arguments) String(s string) *Arguments                { return a.add(s) }
func (a *Arguments) Strings(s ...string) *Arguments            { return a.add(s) }
func (a *Arguments) Time(t time.Time) *Arguments               { return a.add(t) }
func (a *Arguments) Times(t ...time.Time) *Arguments           { return a.add(t) }
func (a *Arguments) Bytes(b []byte) *Arguments                 { return a.add(b) }
func (a *Arguments) BytesSlice(b ...[]byte) *Arguments         { return a.add(b) }
func (a *Arguments) NullString(nv NullString) *Arguments       { return a.add(nv) }
func (a *Arguments) NullStrings(nv ...NullString) *Arguments   { return a.add(nv) }
func (a *Arguments) NullFloat64(nv NullFloat64) *Arguments     { return a.add(nv) }
func (a *Arguments) NullFloat64s(nv ...NullFloat64) *Arguments { return a.add(nv) }
func (a *Arguments) NullInt64(nv NullInt64) *Arguments         { return a.add(nv) }
func (a *Arguments) NullInt64s(nv ...NullInt64) *Arguments     { return a.add(nv) }
func (a *Arguments) NullBool(nv NullBool) *Arguments           { return a.add(nv) }
func (a *Arguments) NullBools(nv ...NullBool) *Arguments       { return a.add(nv) }
func (a *Arguments) NullTime(nv NullTime) *Arguments           { return a.add(nv) }
func (a *Arguments) NullTimes(nv ...NullTime) *Arguments       { return a.add(nv) }

// Name sets the name for the following argument. Calling Name two times after
// each other sets the first call to Name to a NULL value. A call to Name should
// always follow a call to a function type like Int, Float64s or NullTime.
// Name may contain the placeholder prefix colon.
func (a *Arguments) Name(n string) *Arguments {
	a.arguments = append(a.arguments, argument{name: n})
	return a
}

// Reset resets the slice for new usage retaining the already allocated memory.
// It does not reset the Options field.
func (a *Arguments) Reset() *Arguments {
	for i := range a.recs {
		a.recs[i].Qualifier = ""
		a.recs[i].Record = nil
	}
	a.recs = a.recs[:0]
	a.arguments = a.arguments[:0]
	a.raw = a.raw[:0]
	a.nextUnnamedArgPos = 0
	return a
}

// ResetInsert same as Reset but only applicable with INSERT statements and
// triggers a new build of the VALUES part. This function must be called when
// the number of argument changes.
func (a *Arguments) ResetInsert() *Arguments {
	a.insertIsBuildValues = false
	a.insertCachedSQL = a.insertCachedSQL[:0]
	return a.Reset()
}

// DriverValue adds multiple of the same underlying values to the argument
// slice. When using different values, the last applied value wins and gets
// added to the argument slice. For example driver.Values of type `int` will
// result in []int.
func (a *Arguments) DriverValue(dvs ...driver.Valuer) *Arguments {
	if a == nil {
		a = MakeArgs(len(dvs))
	}
	a.arguments, a.base.ärgErr = driverValue(a.arguments, dvs...)
	return a
}

// DriverValues adds each driver.Value as its own argument to the argument
// slice. It panics if the underlying type is not one of the allowed of
// interface driver.Valuer.
func (a *Arguments) DriverValues(dvs ...driver.Valuer) *Arguments {
	if a == nil {
		a = MakeArgs(len(dvs))
	}
	a.arguments, a.base.ärgErr = driverValues(a.arguments, dvs...)
	return a
}

// WithDB sets the database query object.
func (a *Arguments) WithDB(db QueryExecPreparer) *Arguments {
	a.base.DB = db
	return a
}

// WithTx sets the transaction query executor and the logger to run this query
// within a transaction.
func (a *Arguments) WithTx(tx *Tx) *Arguments {
	if a.base.id == "" {
		a.base.id = tx.makeUniqueID()
	}
	a.base.Log = tx.Log
	a.base.DB = tx.DB
	return a
}

/*****************************************************************************************************
	LOAD / QUERY and EXEC functions
*****************************************************************************************************/

var pooledColumnMap = sync.Pool{
	New: func() interface{} {
		return NewColumnMap(20, "")
	},
}

func pooledColumnMapGet() *ColumnMap {
	return pooledColumnMap.Get().(*ColumnMap)
}

func pooledBufferColumnMapPut(cm *ColumnMap, buf *bufferpool.TwinBuffer, fn func()) {
	if buf != nil {
		bufferpool.PutTwin(buf)
	}
	if fn != nil {
		fn()
	}
	cm.reset()
	pooledColumnMap.Put(cm)
}

var pooledArguments = sync.Pool{
	New: func() interface{} {
		var a [16]argument
		return arguments(a[:0])
	},
}

func pooledArgumentsGet() arguments {
	return pooledArguments.Get().(arguments)
}

func pooledArgumentsPut(a arguments, buf *bufferpool.TwinBuffer) {
	a = a[:0]
	pooledArguments.Put(a)
	if buf != nil {
		bufferpool.PutTwin(buf)
	}
}

// Exec executes the statement represented by the Insert object. It returns the
// raw database/sql Result or an error if there was one. Regarding
// LastInsertID(): If you insert multiple rows using a single INSERT statement,
// LAST_INSERT_ID() returns the value generated for the first inserted row only.
// The reason for this at to make it possible to reproduce easily the same
// INSERT statement against some other server. If a record resp. and object
// implements the interface LastInsertIDAssigner then the LastInsertID gets
// assigned incrementally to the objects.
func (a *Arguments) ExecContext(ctx context.Context, args ...interface{}) (sql.Result, error) {
	return a.exec(ctx, args...)
}

// QueryContext traditional way of the databasel/sql package.
func (a *Arguments) QueryContext(ctx context.Context, args ...interface{}) (*sql.Rows, error) {
	return a.query(ctx, args...)
}

// QueryRowContext traditional way of the databasel/sql package.
func (a *Arguments) QueryRowContext(ctx context.Context, args ...interface{}) *sql.Row {
	sqlStr, args, err := a.prepareArgs(args...)
	if a.base.Log != nil && a.base.Log.IsDebug() {
		defer log.WhenDone(a.base.Log).Debug("QueryRowContext", log.String("sql", sqlStr), log.String("source", string(a.base.source)), log.Err(err))
	}
	return a.base.DB.QueryRowContext(ctx, sqlStr, args...)
}

// IterateSerial iterates in serial order over the result set by loading one row each
// iteration and then discarding it. Handles records one by one. The context
// gets only used in the Query function.
func (a *Arguments) IterateSerial(ctx context.Context, callBack func(*ColumnMap) error, args ...interface{}) (err error) {
	if a.base.Log != nil && a.base.Log.IsDebug() {
		defer log.WhenDone(a.base.Log).Debug("IterateSerial", log.String("id", a.base.id), log.Err(err))
	}

	r, err := a.query(ctx, args...)
	if err != nil {
		err = errors.Wrapf(err, "[dml] IterateSerial.Query with query ID %q", a.base.id)
		return
	}
	cmr := pooledColumnMapGet() // this sync.Pool might not work correctly, write a complex test.
	defer pooledBufferColumnMapPut(cmr, nil, func() {
		// Not testable with the sqlmock package :-(
		if err2 := r.Close(); err2 != nil && err == nil {
			err = errors.Wrap(err2, "[dml] IterateSerial.QueryContext.Rows.Close")
		}
	})

	for r.Next() {
		if err = cmr.Scan(r); err != nil {
			err = errors.WithStack(err)
			return
		}
		if err = callBack(cmr); err != nil {
			err = errors.WithStack(err)
			return
		}
	}
	err = errors.WithStack(r.Err())
	return
}

// iterateParallelForNextLoop has been extracted from IterateParallel to not
// mess around with closing channels in different locations of the source code
// when an error occurs.
func iterateParallelForNextLoop(r *sql.Rows, rowChan chan<- *ColumnMap, errChan <-chan error) (err error) {
	defer func() {
		if err2 := r.Err(); err2 != nil && err == nil {
			err = errors.WithStack(err)
		}
		if err2 := r.Close(); err2 != nil && err == nil {
			err = errors.Wrap(err2, "[dml] IterateParallel.QueryContext.Rows.Close")
		}
	}()

	var idx uint64
	for r.Next() {
		var cm ColumnMap // must be empty because we're not collecting data
		if errS := cm.Scan(r); errS != nil {
			err = errors.WithStack(errS)
			return
		}
		cm.Count = idx
		select {
		case rowChan <- &cm:
		case errC := <-errChan:
			if errC != nil {
				err = errC
			}
			return
		}
		idx++
	}
	return
}

// IterateParallel starts a number of workers as defined by variable
// concurrencyLevel and executes the query. Each database row gets evenly
// distributed to the workers. The callback function gets called within a
// worker. concurrencyLevel should be the number of CPUs.
func (a *Arguments) IterateParallel(ctx context.Context, concurrencyLevel int, callBack func(*ColumnMap) error, args ...interface{}) (err error) {
	if a.base.Log != nil && a.base.Log.IsDebug() {
		defer log.WhenDone(a.base.Log).Debug("IterateParallel", log.String("id", a.base.id), log.Err(err))
	}
	if concurrencyLevel < 1 {
		return errors.OutofRange.Newf("[dml] Arguments.IterateParallel concurrencyLevel %d for query ID %q cannot be smaller zero.", concurrencyLevel, a.base.id)
	}

	r, err := a.query(ctx, args...)
	if err != nil {
		err = errors.Wrapf(err, "[dml] IterateParallel.Query with query ID %q", a.base.id)
		return
	}

	// start workers and a channel for communicating
	rowChan := make(chan *ColumnMap)
	errChan := make(chan error, concurrencyLevel)
	var wg sync.WaitGroup
	for i := 0; i < concurrencyLevel; i++ {
		wg.Add(1)
		//i := i
		go func(wg *sync.WaitGroup, rowChan <-chan *ColumnMap, errChan chan<- error) {
			defer wg.Done()
			for cmr := range rowChan {
				if cbErr := callBack(cmr); cbErr != nil {
					errChan <- errors.WithStack(cbErr)
					return
				}
			}
		}(&wg, rowChan, errChan)
	}

	if err2 := iterateParallelForNextLoop(r, rowChan, errChan); err2 != nil {
		err = err2
	}
	close(rowChan)
	wg.Wait()
	close(errChan)

	var multiErr *errors.MultiErr
	i := 0
	for errC2 := range errChan {
		if i == 0 && err != nil {
			multiErr = multiErr.AppendErrors(err)
		}
		multiErr = multiErr.AppendErrors(errC2)
		i++
	}
	if multiErr != nil {
		err = multiErr
	}
	return
}

// Load loads data from a query into an object. Load can load a single row or
// muliple-rows. It checks on top if ColumnMapper `s` implements io.Closer, to
// call the custom close function. This is useful for e.g. unlocking a mutex.
func (a *Arguments) Load(ctx context.Context, s ColumnMapper, args ...interface{}) (rowCount uint64, err error) {
	if a.base.Log != nil && a.base.Log.IsDebug() {
		defer log.WhenDone(a.base.Log).Debug("Load", log.String("id", a.base.id), log.Err(err), log.ObjectTypeOf("ColumnMapper", s), log.Uint64("row_count", rowCount))
	}

	r, err := a.query(ctx, args...)
	if err != nil {
		err = errors.Wrapf(err, "[dml] Arguments.Load.QueryContext failed with queryID %q and ColumnMapper %T", a.base.id, s)
		return
	}
	cm := pooledColumnMapGet()
	defer pooledBufferColumnMapPut(cm, nil, func() {
		// Not testable with the sqlmock package :-(
		if err2 := r.Close(); err2 != nil && err == nil {
			err = errors.Wrap(err2, "[dml] Arguments.Load.Rows.Close")
		}
		if rc, ok := s.(ioCloser); ok {
			if err2 := rc.Close(); err2 != nil && err == nil {
				err = errors.Wrap(err2, "[dml] Arguments.Load.ColumnMapper.Close")
			}
		}
	})

	for r.Next() {
		if err = cm.Scan(r); err != nil {
			return 0, errors.WithStack(err)
		}
		if err = s.MapColumns(cm); err != nil {
			return 0, errors.Wrapf(err, "[dml] Arguments.Load failed with queryID %q and ColumnMapper %T", a.base.id, s)
		}
	}
	if err = r.Err(); err != nil {
		return 0, errors.WithStack(err)
	}
	if cm.HasRows {
		cm.Count++ // because first row is zero but we want the actual row number
	}
	rowCount = cm.Count
	return
}

// LoadInt64 executes the prepared statement and returns the value as an
// int64. It returns a NotFound error if the query returns nothing.
func (a *Arguments) LoadInt64(ctx context.Context, args ...interface{}) (int64, error) {
	if a.base.Log != nil && a.base.Log.IsDebug() {
		defer log.WhenDone(a.base.Log).Debug("LoadInt64")
	}
	return loadInt64(a.query(ctx, args...))
}

// LoadInt64s executes the Select and returns the value as a slice of
// int64s.
func (a *Arguments) LoadInt64s(ctx context.Context, args ...interface{}) (ret []int64, err error) {
	if a.base.Log != nil && a.base.Log.IsDebug() {
		// do not use fullSQL because we might log sensitive data
		defer log.WhenDone(a.base.Log).Debug("LoadInt64s", log.Int("row_count", len(ret)), log.Err(err))
	}
	ret, err = loadInt64s(a.query(ctx, args...))
	// Do not simplify it because we need ret in the defer. we don't log errors
	// because they get handled.
	return
}

// LoadUint64 executes the Select and returns the value at an uint64. It returns
// a NotFound error if the query returns nothing. This function comes in handy
// when performing a COUNT(*) query. See function `Select.Count`.
func (a *Arguments) LoadUint64(ctx context.Context, args ...interface{}) (_ uint64, err error) {
	if a.base.Log != nil && a.base.Log.IsDebug() {
		// do not use fullSQL because we might log sensitive data
		defer log.WhenDone(a.base.Log).Debug("LoadUint64", log.String("id", a.base.id), log.Err(err))
	}

	rows, err := a.query(ctx, args...)
	if err != nil {
		return 0, errors.WithStack(err)
	}
	defer func() {
		if errC := rows.Close(); err == nil && errC != nil {
			err = errors.WithStack(errC)
		}
	}()

	var value uint64
	found := false
	for rows.Next() {
		if err = rows.Scan(&value); err != nil {
			return 0, errors.WithStack(err)
		}
		found = true
	}
	if err = rows.Err(); err != nil {
		return 0, errors.WithStack(err)
	}
	if !found {
		err = errors.NotFound.Newf("[dml] LoadUint64 value not found")
	}
	return value, err
}

// LoadUint64s executes the Select and returns the value at a slice of uint64s.
func (a *Arguments) LoadUint64s(ctx context.Context, args ...interface{}) (values []uint64, err error) {
	if a.base.Log != nil && a.base.Log.IsDebug() {
		// do not use fullSQL because we might log sensitive data
		defer log.WhenDone(a.base.Log).Debug("LoadUint64s", log.Int("row_count", len(values)), log.String("id", a.base.id), log.Err(err))
	}

	rows, err := a.query(ctx, args...)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	defer func() {
		if errC := rows.Close(); err == nil && errC != nil {
			err = errors.WithStack(errC)
		}
	}()

	values = make([]uint64, 0, 10)
	for rows.Next() {
		var value uint64
		if err = rows.Scan(&value); err != nil {
			return nil, errors.WithStack(err)
		}
		values = append(values, value)
	}
	if err = rows.Err(); err != nil {
		return nil, errors.WithStack(err)
	}
	return
}

// LoadFloat64 executes the Select and returns the value at an float64. It
// returns a NotFound error if the query returns nothing.
func (a *Arguments) LoadFloat64(ctx context.Context, args ...interface{}) (_ float64, err error) {
	if a.base.Log != nil && a.base.Log.IsDebug() {
		// do not use fullSQL because we might log sensitive data
		defer log.WhenDone(a.base.Log).Debug("LoadFloat64", log.String("id", a.base.id), log.Err(err))
	}

	rows, err := a.query(ctx, args...)
	if err != nil {
		return 0, errors.WithStack(err)
	}
	defer func() {
		if errC := rows.Close(); err == nil && errC != nil {
			err = errors.WithStack(errC)
		}
	}()

	var value float64
	found := false
	for rows.Next() {
		if err = rows.Scan(&value); err != nil {
			return 0, errors.WithStack(err)
		}
		found = true
	}
	if err = rows.Err(); err != nil {
		return 0, errors.WithStack(err)
	}
	if !found {
		err = errors.NotFound.Newf("[dml] LoadFloat64 value not found")
	}
	return value, err
}

// LoadFloat64s executes the Select and returns the value at a slice of float64s.
func (a *Arguments) LoadFloat64s(ctx context.Context, args ...interface{}) (_ []float64, err error) {
	if a.base.Log != nil && a.base.Log.IsDebug() {
		// do not use fullSQL because we might log sensitive data
		defer log.WhenDone(a.base.Log).Debug("LoadFloat64s", log.String("id", a.base.id), log.Err(err))
	}

	rows, err := a.query(ctx, args...)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	defer func() {
		if errC := rows.Close(); err == nil && errC != nil {
			err = errors.WithStack(errC)
		}
	}()

	values := make([]float64, 0, 10)
	for rows.Next() {
		var value float64
		if err = rows.Scan(&value); err != nil {
			return nil, errors.WithStack(err)
		}
		values = append(values, value)
	}
	if err = rows.Err(); err != nil {
		return nil, errors.WithStack(err)
	}
	return values, err
}

// LoadString executes the Select and returns the value as a string. It
// returns a NotFound error if the row amount is not equal one.
func (a *Arguments) LoadString(ctx context.Context, args ...interface{}) (_ string, err error) {
	if a.base.Log != nil && a.base.Log.IsDebug() {
		// do not use fullSQL because we might log sensitive data
		defer log.WhenDone(a.base.Log).Debug("LoadString", log.String("id", a.base.id), log.Err(err))
	}

	rows, err := a.query(ctx, args...)
	if err != nil {
		return "", errors.WithStack(err)
	}
	defer func() {
		if errC := rows.Close(); err == nil && errC != nil {
			err = errors.WithStack(errC)
		}
	}()

	var value string
	found := false
	for rows.Next() {
		if err = rows.Scan(&value); err != nil {
			return "", errors.WithStack(err)
		}
		found = true
	}
	if err = rows.Err(); err != nil {
		return "", errors.WithStack(err)
	}
	if !found {
		err = errors.NotFound.Newf("[dml] LoadInt64 value not found")
	}
	return value, err
}

// LoadStrings executes the Select and returns a slice of strings.
func (a *Arguments) LoadStrings(ctx context.Context, args ...interface{}) (values []string, err error) {
	if a.base.Log != nil && a.base.Log.IsDebug() {
		// do not use fullSQL because we might log sensitive data
		defer log.WhenDone(a.base.Log).Debug("LoadStrings", log.Int("row_count", len(values)), log.String("id", a.base.id), log.Err(err))
	}

	rows, err := a.query(ctx, args...)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	defer func() {
		if errC := rows.Close(); err == nil && errC != nil {
			err = errors.WithStack(errC)
		}
	}()

	values = make([]string, 0, 10)
	for rows.Next() {
		var value string
		if err = rows.Scan(&value); err != nil {
			return nil, errors.WithStack(err)
		}
		values = append(values, value)
	}
	if err = rows.Err(); err != nil {
		return nil, errors.WithStack(err)
	}
	return values, err
}

func (a *Arguments) query(ctx context.Context, args ...interface{}) (rows *sql.Rows, err error) {
	sqlStr, args, err2 := a.prepareArgs(args...)
	err = err2
	if a.base.Log != nil && a.base.Log.IsDebug() {
		defer log.WhenDone(a.base.Log).Debug("Query", log.String("sql", sqlStr), log.String("source", string(a.base.source)), log.Err(err))
	}
	if err != nil {
		return nil, errors.WithStack(err)
	}

	rows, err = a.base.DB.QueryContext(ctx, sqlStr, args...)
	if err != nil {
		err = errors.Wrapf(err, "[dml] Query.QueryContext with query %q", sqlStr)
	}
	return
}

func loadInt64(rows *sql.Rows, errIn error) (value int64, err error) {
	if errIn != nil {
		return 0, errors.WithStack(errIn)
	}

	defer func() {
		if cErr := rows.Close(); err == nil && cErr != nil {
			err = errors.WithStack(cErr)
		}
	}()

	found := false
	for rows.Next() {
		if err = rows.Scan(&value); err != nil {
			return 0, errors.WithStack(err)
		}
		found = true
	}
	if err = rows.Err(); err != nil {
		return 0, errors.WithStack(err)
	}
	if !found {
		err = errors.NotFound.Newf("[dml] LoadInt64 value not found")
	}
	return value, err
}

func loadInt64s(rows *sql.Rows, errIn error) (_ []int64, err error) {
	if errIn != nil {
		return nil, errors.WithStack(errIn)
	}
	defer func() {
		if cErr := rows.Close(); err == nil && cErr != nil {
			err = errors.WithStack(cErr)
		}
	}()

	values := make([]int64, 0, 16)
	for rows.Next() {
		var value int64
		if err = rows.Scan(&value); err != nil {
			return nil, errors.WithStack(err)
		}
		values = append(values, value)
	}
	if err = rows.Err(); err != nil {
		return nil, errors.WithStack(err)
	}
	return values, nil
}

func (a *Arguments) exec(ctx context.Context, args ...interface{}) (result sql.Result, err error) {
	sqlStr, args, err2 := a.prepareArgs(args...)
	err = err2
	if a.base.Log != nil && a.base.Log.IsDebug() {
		defer log.WhenDone(a.base.Log).Debug("Exec", log.String("sql", sqlStr), log.String("source", string(a.base.source)), log.Err(err))
	}
	if err != nil {
		return nil, errors.WithStack(err)
	}

	result, err = a.base.DB.ExecContext(ctx, sqlStr, args...)
	if err != nil {
		err = errors.Wrapf(err, "[dml] ExecContext with query %q", sqlStr) // err gets catched by the defer
		return
	}

	if a.recs == nil {
		return result, nil
	}
	lID, err := result.LastInsertId()
	if err != nil {
		err = errors.WithStack(err)
		return
	}
	for i, rec := range a.recs {
		if a, ok := rec.Record.(LastInsertIDAssigner); ok {
			a.AssignLastInsertID(lID + int64(i))
		}
	}
	return
}
