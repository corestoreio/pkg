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
	"sync"

	"github.com/corestoreio/errors"
	"github.com/corestoreio/log"
	"github.com/corestoreio/pkg/storage/null"
	"github.com/corestoreio/pkg/util/bufferpool"
	"github.com/corestoreio/pkg/util/byteconv"
	"golang.org/x/sync/errgroup"
)

// DBR is a DataBaseRunner which prepares the SQL string from a DML type,
// collects and build a list of arguments for later sending and execution in the
// database server. Arguments are collections of primitive types or slices of
// primitive types. An DBR type acts like a prepared statement. In fact it can
// contain under the hood different connection types. DBR is optimized for reuse
// and allow saving memory allocations. It can't be used in concurrent context.
type DBR struct {
	base builderCommon
	// QualifiedColumnsAliases allows to overwrite the internal qualified
	// columns slice with custom names. Only in the use case when records are
	// applied. The list of column names in `QualifiedColumnsAliases` gets
	// passed to the ColumnMapper and back to the provided object. The
	// `QualifiedColumnsAliases` slice must have the same length as the
	// qualified columns slice. The order of the alias names must be in the same
	// order as the qualified columns or as the placeholders occur.
	QualifiedColumnsAliases []string
	OrderBys                ids
	LimitValid              bool
	OffsetValid             bool
	LimitCount              uint64
	OffsetCount             uint64
	// insertCachedSQL contains the final build SQL string with the correct
	// amount of placeholders.
	insertCachedSQL     string
	insertColumnCount   uint
	tupleRowCount       uint
	insertIsBuildValues bool
	// isPrepared if true the cachedSQL field in base gets ignored
	isPrepared bool
	// Options like enable interpolation or expanding placeholders.
	Options uint
}

const (
	argOptionExpandPlaceholder = 1 << iota
	argOptionInterpolate
)

// DBRFunc defines a call back function used in other packages to allow
// modifications to the DBR object.
type DBRFunc func(*DBR)

// ApplyCallBacks may clone the DBR object and applies various functions to the new
// DBR instance. It only clones if slice `fns` has at least one entry.
func (a *DBR) ApplyCallBacks(fns ...DBRFunc) *DBR {
	if len(fns) == 0 {
		return a
	}

	ac := a.Clone() // locks
	for _, af := range fns {
		af(ac)
	}
	return ac
}

// OrderBy appends columns to the ORDER BY statement for ascending sorting. A
// column gets always quoted if it is a valid identifier otherwise it will be
// treated as an expression. This ORDER BY clause gets appended to the current
// internal cached SQL string independently if the SQL statement supports it or
// not or if there exists already an ORDER BY clause.
// A column name can also contain the suffix words " ASC" or " DESC" to indicate
// the sorting. This avoids using the method OrderByDesc when sorting certain
// columns descending.
func (a *DBR) OrderBy(columns ...string) *DBR {
	a.OrderBys = a.OrderBys.AppendColumns(false, columns...)
	return a
}

// OrderByDesc appends columns to the ORDER BY statement for descending sorting.
// A column gets always quoted if it is a valid identifier otherwise it will be
// treated as an expression. This ORDER BY clause gets appended to the current
// internal cached SQL string independently if the SQL statement supports it or
// not or if there exists already an ORDER BY clause.
func (a *DBR) OrderByDesc(columns ...string) *DBR {
	a.OrderBys = a.OrderBys.AppendColumns(false, columns...).applySort(len(columns), sortDescending)
	return a
}

// Limit sets a LIMIT clause for the statement; overrides any existing LIMIT.
// This LIMIT clause gets appended to the current internal cached SQL string
// independently if the SQL statement supports it or not or if there exists
// already a LIMIT clause.
func (a *DBR) Limit(offset, limit uint64) *DBR {
	a.OffsetCount = offset
	a.LimitCount = limit
	a.OffsetValid = true
	a.LimitValid = true
	return a
}

// Paginate sets LIMIT/OFFSET for the statement based on the given page/perPage
// Assumes page/perPage are valid. Page and perPage must be >= 1
func (a *DBR) Paginate(page, perPage uint64) *DBR {
	a.Limit((page-1)*perPage, perPage)
	return a
}

// WithQualifiedColumnsAliases for documentation please see:
// DBR.QualifiedColumnsAliases.
func (a *DBR) WithQualifiedColumnsAliases(aliases ...string) *DBR {
	a.QualifiedColumnsAliases = aliases
	return a
}

// ToSQL generates the SQL string.
func (a *DBR) ToSQL() (string, []interface{}, error) {
	sqlStr, _, err := a.prepareQueryAndArgs(nil)
	return sqlStr, nil, err
}

// TestWithArgs returns a QueryBuilder with resolved arguments. Mostly used for
// testing and in examples to skip the calls to ExecContext or QueryContext.
// Every 2nd call arguments are getting interpolated.
func (a *DBR) TestWithArgs(args ...interface{}) QueryBuilder {
	var secondCallInterpolates uint
	return QuerySQLFn(func() (string, []interface{}, error) {
		if secondCallInterpolates > 0 && secondCallInterpolates%2 == 1 {
			a.Interpolate()
		} else if a.Options&argOptionInterpolate != 0 {
			a.Options ^= argOptionInterpolate // removes interpolation AND NOT
		}
		secondCallInterpolates++

		sqlStr, args, err := a.prepareQueryAndArgs(args)
		return sqlStr, args, err
	})
}

func (a *DBR) testWithArgs(args ...interface{}) QueryBuilder {
	return QuerySQLFn(func() (string, []interface{}, error) {
		sqlStr, args, err := a.prepareQueryAndArgs(args)
		return sqlStr, args, err
	})
}

// CachedQueries returns a list with all cached SQL queries.
func (a *DBR) CachedQueries(queries ...string) []string {
	return a.base.CachedQueries(queries...)
}

// WithCacheKey sets the currently used cache key when generating a SQL string.
// By setting a different cache key, a previous generated SQL query is
// accessible again. New cache keys allow to change the generated query of the
// current object. E.g. different where clauses or different row counts in
// INSERT ... VALUES statements. The empty string defines the default cache key.
// If the `args` argument contains values, then fmt.Sprintf gets used.
func (a *DBR) WithCacheKey(key string, args ...interface{}) *DBR {
	a.base.withCacheKey(key, args...)
	return a
}

// Interpolate if set stringyfies the arguments into the SQL string and returns
// pre-processed SQL command when calling the function ToSQL. Not suitable for
// prepared statements. ToSQLs second argument `args` will then be nil.
func (a *DBR) Interpolate() *DBR {
	a.Options |= argOptionInterpolate
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
func (a *DBR) ExpandPlaceHolders() *DBR {
	a.Options |= argOptionExpandPlaceholder
	return a
}

// prepareQueryAndArgs transforms mainly the DBR into []interface{}. It appends
// its arguments to the `extArgs` arguments from the Exec+ or Query+ function.
// This allows for a developer to reuse the interface slice and save
// allocations. All method receivers are not thread safe. The returned interface
// slice is the same as `extArgs`.
// The returned []QualifiedRecord slice is needed to use interface LastInsertIDAssigner.
func (a *DBR) prepareQueryAndArgs(extArgs []interface{}) (_ string, _ []interface{}, err error) {
	if a.base.ärgErr != nil {
		return "", nil, errors.WithStack(a.base.ärgErr)
	}
	lenExtArgs := len(extArgs)
	var hasNamedArgs uint8
	var containsQualifiedRecords int
	var primitiveCounts int
	var args []interface{}
	if lenExtArgs > 0 {
		args = pooledInterfacesGet()
		defer pooledInterfacesPut(args)

		for _, ea := range extArgs {
			switch eaTypeValue := ea.(type) {
			case nil:
				args = append(args, internalNULLNIL{})
				primitiveCounts++
			case QualifiedRecord, ColumnMapper:
				containsQualifiedRecords++
				args = append(args, eaTypeValue)
			case sql.NamedArg:
				args = append(args, ea)
				hasNamedArgs = 2
				primitiveCounts++
			case []sql.NamedArg: // insert statement with key/value pairs
				for _, na := range eaTypeValue {
					args = append(args, na.Value)
					primitiveCounts++
				}
			default:
				args = append(args, ea)
				primitiveCounts++
			}
		}
	}
	if a.base.source == dmlSourceInsert {
		return a.prepareQueryAndArgsInsert(args, primitiveCounts)
	}

	cachedSQL, ok := a.base.cachedSQL[a.base.cacheKey]
	if !a.isPrepared && !ok {
		return "", nil, errors.Empty.Newf("[dml] DBR: The SQL string is empty.")
	}

	if a.base.templateStmtCount < 2 && hasNamedArgs == 0 && containsQualifiedRecords == 0 && a.Options == 0 { // no options and qualified records provided
		if a.isPrepared {
			return "", expandInterfaces(args), nil
		}

		if a.Options == 0 && len(a.OrderBys) == 0 && !a.LimitValid {
			return cachedSQL, expandInterfaces(args), nil
		}
		buf := bufferpool.Get()
		defer bufferpool.Put(buf)
		buf.WriteString(cachedSQL)
		sqlWriteOrderBy(buf, a.OrderBys, false)
		sqlWriteLimitOffset(buf, a.LimitValid, a.OffsetValid, a.OffsetCount, a.LimitCount)
		return buf.String(), expandInterfaces(args), nil
	}

	if !a.isPrepared && hasNamedArgs == 0 {
		var found bool
		hasNamedArgs = 1
		cachedSQL, a.base.qualifiedColumns, found = extractReplaceNamedArgs(cachedSQL, a.base.qualifiedColumns)
		if found {
			a.base.cachedSQLUpsert(a.base.cacheKey, cachedSQL)
			hasNamedArgs = 2
		}
	}

	sqlBuf := bufferpool.GetTwin()
	defer bufferpool.PutTwin(sqlBuf)

	if args, err = a.appendConvertedRecordsToArguments(hasNamedArgs, args, containsQualifiedRecords); err != nil {
		return "", nil, errors.WithStack(err)
	}

	if a.isPrepared {
		return "", expandInterfaces(args), nil
	}

	// Make a copy of the original SQL statement because it gets modified in the
	// worst case. Best case would be no modification and hence we don't need a
	// bytes.Buffer from the pool! TODO(CYS) optimize this and only acquire a
	// buffer from the pool in the worse case.
	if _, err := sqlBuf.First.WriteString(cachedSQL); err != nil {
		return "", nil, errors.WithStack(err)
	}

	sqlWriteOrderBy(sqlBuf.First, a.OrderBys, false)
	sqlWriteLimitOffset(sqlBuf.First, a.LimitValid, a.OffsetValid, a.OffsetCount, a.LimitCount)

	// `switch` statement no suitable.
	if a.Options > 0 && lenExtArgs > 0 && containsQualifiedRecords == 0 && len(args) == 0 {
		return "", nil, errors.NotAllowed.Newf("[dml] Interpolation/ExpandPlaceholders supports only Records and Arguments and not yet an interface slice.")
	}

	// TODO more advanced caching of the final non-expanded SQL string

	if a.Options&argOptionExpandPlaceholder != 0 {
		phCount := bytes.Count(sqlBuf.First.Bytes(), placeHolderByte)
		if aLen, hasSlice := totalSliceLen(args); phCount < aLen || hasSlice {
			if a.base.containsTuples {
				if err := expandPlaceHolderTuples(sqlBuf.Second, sqlBuf.First.Bytes(), aLen); err != nil {
					return "", nil, errors.WithStack(err)
				}
				if _, err := sqlBuf.CopySecondToFirst(); err != nil {
					return "", nil, errors.WithStack(err)
				}
			}
			if err := expandPlaceHolders(sqlBuf.Second, sqlBuf.First.Bytes(), args); err != nil {
				return "", nil, errors.WithStack(err)
			}
			if _, err := sqlBuf.CopySecondToFirst(); err != nil {
				return "", nil, errors.WithStack(err)
			}
		}
	}
	if a.Options&argOptionInterpolate != 0 {
		if err := writeInterpolateBytes(sqlBuf.Second, sqlBuf.First.Bytes(), args); err != nil {
			return "", nil, errors.Wrapf(err, "[dml] Interpolation failed: %q", sqlBuf.String())
		}
		return sqlBuf.Second.String(), nil, nil
	}

	return sqlBuf.First.String(), expandInterfaces(args), nil
}

func (a *DBR) appendConvertedRecordsToArguments(hasNamedArgs uint8, collectedArgs []interface{}, containsQualifiedRecords int) ([]interface{}, error) {
	templateStmtCount := a.base.templateStmtCount
	if a.base.templateStmtCount == 0 {
		templateStmtCount = 1
	}

	if containsQualifiedRecords == 0 && hasNamedArgs == 0 {
		return collectedArgs, nil
	}

	if containsQualifiedRecords == 0 && hasNamedArgs < 2 {
		if templateStmtCount > 1 {
			collectedArgs = multiplyInterfaceValues(collectedArgs, templateStmtCount)
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

	var nextUnnamedArgPos int
	// TODO refactor prototype and make it performant and beautiful code
	cm := NewColumnMap(len(collectedArgs)+containsQualifiedRecords, "") // can use an arg pool DBR sync.Pool, nope.
	for tsc := 0; tsc < templateStmtCount; tsc++ {                      // only in case of UNION statements in combination with a template SELECT, can be optimized later

		// `qualifiedColumns` contains the correct order as the place holders
		// appear in the SQL string.
		for _, identifier := range qualifiedColumns {
			// identifier can be either: column or qualifier.column or :column
			qualifier, column := splitColumn(identifier)
			// a.base.defaultQualifier is empty in case of INSERT statements

			column, isNamedArg := cutNamedArgStartStr(column) // removes the colon for named arguments
			cm.columns[0] = column                            // length is always one, as created in NewColumnMap

			if isNamedArg && len(collectedArgs) > 0 {
				// if the colon : cannot be found then a simple place holder ? has been detected
				if err := a.mapColumns(containsQualifiedRecords, collectedArgs, cm); err != nil {
					return collectedArgs, errors.WithStack(err)
				}
			} else {
				found := false
				for _, arg := range collectedArgs {
					switch qRec := arg.(type) {
					case QualifiedRecord:
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

					case ColumnMapper:
						if err := qRec.MapColumns(cm); err != nil {
							return collectedArgs, errors.WithStack(err)
						}

					}
				}
				if !found {
					// If the argument cannot be found in the records then we assume the argument
					// has a numerical position and we grab just the next unnamed argument.
					var ok bool
					var pArg interface{}
					if pArg, nextUnnamedArgPos, ok = a.nextUnnamedArg(nextUnnamedArgPos, collectedArgs); ok {
						cm.args = append(cm.args, pArg)
					}
				}
			}
		}
		nextUnnamedArgPos = 0
	}
	if len(cm.args) > 0 {
		collectedArgs = cm.args
	}

	return collectedArgs, nil
}

// prepareQueryAndArgsInsert prepares the special arguments for an INSERT statement. The
// returned interface slice is the same as the `extArgs` slice. extArgs =
// external arguments.
func (a *DBR) prepareQueryAndArgsInsert(extArgs []interface{}, primitiveCounts int) (string, []interface{}, error) {
	sqlBuf := bufferpool.GetTwin()
	defer bufferpool.PutTwin(sqlBuf)
	cm := NewColumnMap(2*primitiveCounts, a.base.qualifiedColumns...)
	cm.args = extArgs
	lenExtArgsBefore := len(extArgs)
	lenInsertCachedSQL := len(a.insertCachedSQL)
	cachedSQL, _ := a.base.cachedSQL[a.base.cacheKey]
	var containsRecords bool
	{
		if lenInsertCachedSQL > 0 {
			cachedSQL = a.insertCachedSQL
		}
		if _, err := sqlBuf.First.WriteString(cachedSQL); err != nil {
			return "", nil, errors.WithStack(err)
		}

		for _, arg := range extArgs {
			switch qRec := arg.(type) {
			case QualifiedRecord:
				if qRec.Qualifier != "" {
					return "", nil, errors.Fatal.Newf("[dml] Qualifier in %T is not supported and not needed.", qRec)
				}
				if err := qRec.Record.MapColumns(cm); err != nil {
					return "", nil, errors.WithStack(err)
				}
				containsRecords = true
			case ColumnMapper:
				if err := qRec.MapColumns(cm); err != nil {
					return "", nil, errors.WithStack(err)
				}
				containsRecords = true
			}
		}
		primitiveCounts += len(cm.args) - lenExtArgsBefore
	}

	if a.isPrepared {
		// TODO above construct can be more optimized when using prepared statements
		return "", expandInterfaces(cm.args), nil
	}

	if !a.insertIsBuildValues && lenInsertCachedSQL == 0 { // Write placeholder list e.g. "VALUES (?,?),(?,?)"
		odkPos := strings.Index(cachedSQL, onDuplicateKeyPartS)
		if odkPos > 0 {
			sqlBuf.First.Reset()
			sqlBuf.First.WriteString(cachedSQL[:odkPos])
		}

		if a.tupleRowCount > 0 {
			columnCount := uint(primitiveCounts) / a.tupleRowCount
			writeTuplePlaceholders(sqlBuf.First, a.tupleRowCount, columnCount)
		} else if a.insertColumnCount > 0 {
			rowCount := uint(primitiveCounts) / a.insertColumnCount
			if rowCount == 0 {
				rowCount = 1
			}
			writeTuplePlaceholders(sqlBuf.First, rowCount, a.insertColumnCount)
		}
		if odkPos > 0 {
			sqlBuf.First.WriteString(cachedSQL[odkPos:])
		}
	}
	if lenInsertCachedSQL == 0 {
		a.insertCachedSQL = sqlBuf.First.String()
	}

	if a.Options > 0 {
		if primitiveCounts > 0 && !containsRecords && len(cm.args) == 0 {
			return "", nil, errors.NotAllowed.Newf("[dml] Interpolation/ExpandPlaceholders supports only Records and Arguments and not yet an interface slice.")
		}

		if a.Options&argOptionInterpolate != 0 {
			if err := writeInterpolateBytes(sqlBuf.Second, sqlBuf.First.Bytes(), cm.args); err != nil {
				return "", nil, errors.Wrapf(err, "[dml] Interpolation failed: %q", sqlBuf.First.String())
			}
			return sqlBuf.Second.String(), nil, nil
		}
	}

	return a.insertCachedSQL, expandInterfaces(cm.args), nil
}

// nextUnnamedArg returns an unnamed argument by its position.
func (a *DBR) nextUnnamedArg(nextUnnamedArgPos int, args []interface{}) (interface{}, int, bool) {
	var unnamedCounter int
	for _, arg := range args {
		switch arg.(type) {
		case sql.NamedArg, QualifiedRecord, ColumnMapper:
		// skip
		default:
			if unnamedCounter == nextUnnamedArgPos {
				nextUnnamedArgPos++
				return arg, nextUnnamedArgPos, true
			}
			unnamedCounter++
		}
	}
	nextUnnamedArgPos = -1 // nothing found, so no need to further iterate through the []argument slice.
	return nil, nextUnnamedArgPos, false
}

// mapColumns allows to merge one argument slice with another depending on the
// matched columns. Each argument in the slice must be a named argument.
// Implements interface ColumnMapper.
func (a *DBR) mapColumns(containsQualifiedRecords int, args []interface{}, cm *ColumnMap) error {
	if cm.Mode() == ColumnMapEntityReadAll {
		cm.args = append(cm.args, args...)
		return cm.Err()
	}
	var cm2 *ColumnMap
	if containsQualifiedRecords > 0 {
		cm2 = NewColumnMap(1)
	}
	for cm.Next() {
		// now a bit slow ...
		c := cm.Column()
		if containsQualifiedRecords > 0 {
			cm2.setColumns([]string{c})
		}

		keepOnRolling := true
		for i := 0; i < len(args) && keepOnRolling && c != ""; i++ {
			switch at := args[i].(type) {
			case sql.NamedArg:
				if at.Name == c {
					cm.args = append(cm.args, at.Value) // at.Value was previously just at
					keepOnRolling = false
				}
			case QualifiedRecord:
				// ignore the returned error as the mapper might not find the column.
				// in the upper switch case we compare Name==c and so we know which
				// column but with the mapper we don't know.
				_ = at.Record.MapColumns(cm2)
				cm.args = append(cm.args, cm2.args...)
				// do not break the loop like in the upper switch case.
			}
		}
	}
	return cm.Err()
}

// Reset resets the internal slices for new usage retaining the already
// allocated memory. Reset gets called automatically in many Load* functions. In
// case of an INSERT statement, Reset triggers a new build of the VALUES part.
// This function must be called when the number of argument changes.
func (a *DBR) Reset() *DBR {
	a.insertIsBuildValues = false
	a.insertCachedSQL = a.insertCachedSQL[:0]
	return a
}

// WithDB sets the database query object.
func (a *DBR) WithDB(db QueryExecPreparer) *DBR {
	a.base.db = db
	return a
}

// WithPreparedStmt uses a SQL statement as DB connection.
func (a *DBR) WithPreparedStmt(stmt *sql.Stmt) *DBR {
	a.base.db = stmtWrapper{stmt: stmt}
	return a
}

// WithTx sets the transaction query executor and the logger to run this query
// within a transaction.
func (a *DBR) WithTx(tx *Tx) *DBR {
	if a.base.id == "" {
		a.base.id = tx.makeUniqueID()
	}
	a.base.Log = tx.Log
	a.base.db = tx.DB
	return a
}

// Clone creates a shallow clone of the current pointer. The logger gets copied.
// Some underlying slices for the cached SQL statements are still referring to
// the source DBR object.
func (a *DBR) Clone() *DBR {
	c := *a
	if a.QualifiedColumnsAliases != nil {
		c.QualifiedColumnsAliases = make([]string, len(a.QualifiedColumnsAliases))
		copy(c.QualifiedColumnsAliases, a.QualifiedColumnsAliases)
	}
	return &c
}

// Close tries to close the underlying DB connection. Useful in cases of
// prepared statements. If the underlying DB connection does not implement
// io.Closer, nothing will happen.
func (a *DBR) Close() error {
	if a.base.ärgErr != nil {
		return errors.WithStack(a.base.ärgErr)
	}
	if c, ok := a.base.db.(ioCloser); ok {
		return errors.WithStack(c.Close())
	}
	return nil
}

/*****************************************************************************************************
	LOAD / QUERY and EXEC functions
*****************************************************************************************************/

var pooledColumnMap = sync.Pool{
	New: func() interface{} {
		return NewColumnMap(30, "")
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

const argumentPoolMaxSize = 256

// regarding the returned slices in both pools: https://github.com/golang/go/blob/7e394a2/src/net/http/h2_bundle.go#L998-L1043
// they also uses a []byte slice in the pool and not a pointer

var pooledInterfaces = sync.Pool{
	New: func() interface{} {
		return make([]interface{}, 0, argumentPoolMaxSize)
	},
}

func pooledInterfacesGet() []interface{} {
	return pooledInterfaces.Get().([]interface{})
}

func pooledInterfacesPut(args []interface{}) {
	if cap(args) <= argumentPoolMaxSize {
		args = args[:0]
		pooledInterfaces.Put(args)
	}
}

// ExecContext executes the statement represented by the Update/Insert object.
// It returns the raw database/sql Result or an error if there was one.
// Regarding LastInsertID(): If you insert multiple rows using a single INSERT
// statement, LAST_INSERT_ID() returns the value generated for the first
// inserted row only. The reason for this at to make it possible to reproduce
// easily the same INSERT statement against some other server. If a record resp.
// and object implements the interface LastInsertIDAssigner then the
// LastInsertID gets assigned incrementally to the objects. Pro tip: you can use
// function ExecValidateOneAffectedRow to check if the underlying SQL statement
// has affected only one row.
func (a *DBR) ExecContext(ctx context.Context, args ...interface{}) (sql.Result, error) {
	return a.exec(ctx, args)
}

// QueryContext traditional way of the databasel/sql package.
func (a *DBR) QueryContext(ctx context.Context, args ...interface{}) (*sql.Rows, error) {
	return a.query(ctx, args)
}

// QueryRowContext traditional way of the databasel/sql package.
func (a *DBR) QueryRowContext(ctx context.Context, args ...interface{}) *sql.Row {
	sqlStr, args, err := a.prepareQueryAndArgs(args)
	if a.base.Log != nil && a.base.Log.IsDebug() {
		defer log.WhenDone(a.base.Log).Debug(
			"QueryRowContext",
			log.String("sql", sqlStr),
			log.String("source", string(a.base.source)),
			log.Err(err))
	}
	return a.base.db.QueryRowContext(ctx, sqlStr, args...)
}

// IterateSerial iterates in serial order over the result set by loading one row each
// iteration and then discarding it. Handles records one by one. The context
// gets only used in the Query function.
func (a *DBR) IterateSerial(ctx context.Context, callBack func(*ColumnMap) error, args ...interface{}) (err error) {
	if a.base.Log != nil && a.base.Log.IsDebug() {
		defer log.WhenDone(a.base.Log).Debug(
			"IterateSerial",
			log.String("id", a.base.id),
			log.Err(err))
	}

	r, err := a.query(ctx, args)
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
func iterateParallelForNextLoop(ctx context.Context, r *sql.Rows, rowChan chan<- *ColumnMap) (err error) {
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
		case <-ctx.Done():
			return
		}
		idx++
	}
	return
}

// IterateParallel starts a number of workers as defined by variable
// concurrencyLevel and executes the query. Each database row gets evenly
// distributed to the workers. The callback function gets called within a
// worker. concurrencyLevel should be the number of CPUs. You should use this
// function when you expect to process large amount of rows returned from a
// query.
func (a *DBR) IterateParallel(ctx context.Context, concurrencyLevel int, callBack func(*ColumnMap) error, args ...interface{}) (err error) {
	if a.base.Log != nil && a.base.Log.IsDebug() {
		defer log.WhenDone(a.base.Log).Debug("IterateParallel", log.String("id", a.base.id), log.Err(err))
	}
	if concurrencyLevel < 1 {
		return errors.OutOfRange.Newf("[dml] DBR.IterateParallel concurrencyLevel %d for query ID %q cannot be smaller zero.", concurrencyLevel, a.base.id)
	}

	r, err := a.query(ctx, args)
	if err != nil {
		err = errors.Wrapf(err, "[dml] IterateParallel.Query with query ID %q", a.base.id)
		return
	}

	g, ctx := errgroup.WithContext(ctx)

	// start workers and a channel for communicating
	rowChan := make(chan *ColumnMap)
	for i := 0; i < concurrencyLevel; i++ {
		g.Go(func() error {
			for cmr := range rowChan {
				if ctx.Err() != nil {
					return ctx.Err() // terminate this goroutine once the context gets canceled.
				}
				if cbErr := callBack(cmr); cbErr != nil {
					return errors.WithStack(cbErr)
				}
			}
			return nil
		})
	}

	if err2 := iterateParallelForNextLoop(ctx, r, rowChan); err2 != nil {
		err = err2
	}
	close(rowChan)

	return errors.WithStack(g.Wait())
}

// Load loads data from a query into an object. Load can load a single row or
// multiple-rows. It checks on top if ColumnMapper `s` implements io.Closer, to
// call the custom close function. This is useful for e.g. unlocking a mutex.
func (a *DBR) Load(ctx context.Context, s ColumnMapper, args ...interface{}) (rowCount uint64, err error) {
	if a.base.Log != nil && a.base.Log.IsDebug() {
		defer log.WhenDone(a.base.Log).Debug("Load", log.String("id", a.base.id), log.Err(err), log.ObjectTypeOf("ColumnMapper", s), log.Uint64("row_count", rowCount))
	}

	r, err := a.query(ctx, args)
	if err != nil {
		err = errors.Wrapf(err, "[dml] DBR.Load.QueryContext failed with queryID %q and ColumnMapper %T", a.base.id, s)
		return
	}
	cm := pooledColumnMapGet()
	defer pooledBufferColumnMapPut(cm, nil, func() {
		// Not testable with the sqlmock package :-(
		if err2 := r.Close(); err2 != nil && err == nil {
			err = errors.Wrap(err2, "[dml] DBR.Load.Rows.Close")
		}
		if rc, ok := s.(ioCloser); ok {
			if err2 := rc.Close(); err2 != nil && err == nil {
				err = errors.Wrap(err2, "[dml] DBR.Load.ColumnMapper.Close")
			}
		}
	})

	for r.Next() {
		if err = cm.Scan(r); err != nil {
			return 0, errors.WithStack(err)
		}
		if err = s.MapColumns(cm); err != nil {
			return 0, errors.Wrapf(err, "[dml] DBR.Load failed with queryID %q and ColumnMapper %T", a.base.id, s)
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

// LoadNullInt64 executes the query and returns the first row parsed into the
// current type. `Found` might be false if there are no matching rows.
func (a *DBR) LoadNullInt64(ctx context.Context, args ...interface{}) (nv null.Int64, found bool, err error) {
	found, err = a.loadPrimitive(ctx, &nv, args...)
	return
}

// LoadNullUint64 executes the query and returns the first row parsed into the
// current type. `Found` might be false if there are no matching rows.
// This function with ptr type uint64 comes in handy when performing
// a COUNT(*) query. See function `Select.Count`.
func (a *DBR) LoadNullUint64(ctx context.Context, args ...interface{}) (nv null.Uint64, found bool, err error) {
	found, err = a.loadPrimitive(ctx, &nv, args...)
	return
}

// LoadNullFloat64 executes the query and returns the first row parsed into the
// current type. `Found` might be false if there are no matching rows.
func (a *DBR) LoadNullFloat64(ctx context.Context, args ...interface{}) (nv null.Float64, found bool, err error) {
	found, err = a.loadPrimitive(ctx, &nv, args...)
	return
}

// LoadNullString executes the query and returns the first row parsed into the
// current type. `Found` might be false if there are no matching rows.
func (a *DBR) LoadNullString(ctx context.Context, args ...interface{}) (nv null.String, found bool, err error) {
	found, err = a.loadPrimitive(ctx, &nv, args...)
	return
}

// LoadNullTime executes the query and returns the first row parsed into the
// current type. `Found` might be false if there are no matching rows.
func (a *DBR) LoadNullTime(ctx context.Context, args ...interface{}) (nv null.Time, found bool, err error) {
	found, err = a.loadPrimitive(ctx, &nv, args...)
	return
}

// LoadDecimal executes the query and returns the first row parsed into the
// current type. `Found` might be false if there are no matching rows.
func (a *DBR) LoadDecimal(ctx context.Context, args ...interface{}) (nv null.Decimal, found bool, err error) {
	found, err = a.loadPrimitive(ctx, &nv, args...)
	return
}

func (a *DBR) loadPrimitive(ctx context.Context, ptr interface{}, args ...interface{}) (found bool, err error) {
	if a.base.Log != nil && a.base.Log.IsDebug() {
		// do not use fullSQL because we might log sensitive data
		defer log.WhenDone(a.base.Log).Debug("LoadPrimitive", log.String("id", a.base.id), log.Err(err), log.ObjectTypeOf("ptr_type", ptr))
	}
	var rows *sql.Rows
	rows, err = a.query(ctx, args)
	if err != nil {
		return false, errors.WithStack(err)
	}
	defer func() {
		if errC := rows.Close(); err == nil && errC != nil {
			err = errors.WithStack(errC)
		}
	}()

	for rows.Next() && !found {
		if err = rows.Scan(ptr); err != nil {
			return false, errors.WithStack(err)
		}
		found = true
	}
	if err = rows.Err(); err != nil {
		err = errors.WithStack(err)
	}
	return
}

// LoadInt64s executes the query and returns the values appended to slice
// dest. It ignores and skips NULL values.
func (a *DBR) LoadInt64s(ctx context.Context, dest []int64, args ...interface{}) (_ []int64, err error) {
	var rowCount int
	if a.base.Log != nil && a.base.Log.IsDebug() {
		// do not use fullSQL because we might log sensitive data
		defer log.WhenDone(a.base.Log).Debug("LoadInt64s", log.Int("row_count", rowCount), log.Err(err))
	}
	var r *sql.Rows
	r, err = a.query(ctx, args)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	defer func() {
		if cErr := r.Close(); err == nil && cErr != nil {
			err = errors.WithStack(cErr)
		}
	}()
	for r.Next() {
		var nv sql.RawBytes
		if err = r.Scan(&nv); err != nil {
			return nil, errors.WithStack(err)
		}
		if i64, ok, err := byteconv.ParseInt(nv); ok && err == nil {
			dest = append(dest, i64)
		} else if err != nil {
			return nil, errors.WithStack(err)
		}
	}
	if err = r.Err(); err != nil {
		return nil, errors.WithStack(err)
	}

	rowCount = len(dest)
	return dest, err
}

// LoadUint64s executes the query and returns the values appended to slice
// dest. It ignores and skips NULL values.
func (a *DBR) LoadUint64s(ctx context.Context, dest []uint64, args ...interface{}) (_ []uint64, err error) {
	var rowCount int
	if a.base.Log != nil && a.base.Log.IsDebug() {
		// do not use fullSQL because we might log sensitive data
		defer log.WhenDone(a.base.Log).Debug("LoadUint64s", log.Int("row_count", rowCount), log.String("id", a.base.id), log.Err(err))
	}

	rows, err := a.query(ctx, args)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	defer func() {
		if errC := rows.Close(); errC != nil && err == nil {
			err = errors.WithStack(errC)
		}
	}()

	for rows.Next() {
		var nv sql.RawBytes
		if err = rows.Scan(&nv); err != nil {
			return nil, errors.WithStack(err)
		}
		if u64, ok, err := byteconv.ParseUint(nv, 10, 64); ok && err == nil {
			dest = append(dest, u64)
		} else if err != nil {
			return nil, errors.WithStack(err)
		}
	}
	if err = rows.Err(); err != nil {
		return nil, errors.WithStack(err)
	}
	rowCount = len(dest)
	return dest, err
}

// LoadFloat64s executes the query and returns the values appended to slice
// dest. It ignores and skips NULL values.
func (a *DBR) LoadFloat64s(ctx context.Context, dest []float64, args ...interface{}) (_ []float64, err error) {
	if a.base.Log != nil && a.base.Log.IsDebug() {
		// do not use fullSQL because we might log sensitive data
		defer log.WhenDone(a.base.Log).Debug("LoadFloat64s", log.String("id", a.base.id), log.Err(err))
	}

	var rows *sql.Rows
	if rows, err = a.query(ctx, args); err != nil {
		return nil, errors.WithStack(err)
	}
	defer func() {
		if errC := rows.Close(); errC != nil && err == nil {
			err = errors.WithStack(errC)
		}
	}()

	for rows.Next() {
		var nv sql.RawBytes
		if err = rows.Scan(&nv); err != nil {
			return nil, errors.WithStack(err)
		}
		if f64, ok, err := byteconv.ParseFloat(nv); ok && err == nil {
			dest = append(dest, f64)
		} else if err != nil {
			return nil, errors.WithStack(err)
		}
	}
	if err = rows.Err(); err != nil {
		return nil, errors.WithStack(err)
	}
	return dest, err
}

// LoadStrings executes the query and returns the values appended to slice
// dest. It ignores and skips NULL values.
func (a *DBR) LoadStrings(ctx context.Context, dest []string, args ...interface{}) (_ []string, err error) {
	var rowCount int
	if a.base.Log != nil && a.base.Log.IsDebug() {
		// do not use fullSQL because we might log sensitive data
		defer log.WhenDone(a.base.Log).Debug("LoadStrings", log.Int("row_count", rowCount), log.String("id", a.base.id), log.Err(err))
	}

	rows, err := a.query(ctx, args)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	defer func() {
		if errC := rows.Close(); errC != nil && err == nil {
			err = errors.WithStack(errC)
		}
	}()

	for rows.Next() {
		var value sql.RawBytes
		if err = rows.Scan(&value); err != nil {
			return nil, errors.WithStack(err)
		}
		if value != nil {
			dest = append(dest, string(value))
		}
	}
	if err = rows.Err(); err != nil {
		return nil, errors.WithStack(err)
	}
	rowCount = len(dest)
	return dest, err
}

func (a *DBR) query(ctx context.Context, args []interface{}) (rows *sql.Rows, err error) {
	sqlStr, args, err := a.prepareQueryAndArgs(args)
	if a.base.Log != nil && a.base.Log.IsDebug() {
		defer log.WhenDone(a.base.Log).Debug(
			"Query", log.String("sql", sqlStr), log.Int("length_args", len(args)), log.String("source", string(a.base.source)), log.Err(err))
	}
	if err != nil {
		return nil, errors.WithStack(err)
	}
	rows, err = a.base.db.QueryContext(ctx, sqlStr, args...)
	if err != nil {
		if sqlStr == "" {
			cachedSQL, _ := a.base.cachedSQL[a.base.cacheKey]
			sqlStr = "PREPARED:" + cachedSQL
		}
		return nil, errors.Wrapf(err, "[dml] Query.QueryContext with query %q", sqlStr)
	}
	return rows, err
}

func (a *DBR) exec(ctx context.Context, rawArgs []interface{}) (result sql.Result, err error) {
	sqlStr, args, err := a.prepareQueryAndArgs(rawArgs)
	if a.base.Log != nil && a.base.Log.IsDebug() {
		defer log.WhenDone(a.base.Log).Debug("Exec", log.String("sql", sqlStr),
			log.Int("length_args", len(args)), log.Int("length_raw_args", len(rawArgs)), log.String("source", string(a.base.source)),
			log.Err(err))
	}
	if err != nil {
		return nil, errors.WithStack(err)
	}

	result, err = a.base.db.ExecContext(ctx, sqlStr, args...)
	if err != nil {
		return nil, errors.Wrapf(err, "[dml] ExecContext with query %q", sqlStr) // err gets catched by the defer
	}
	lID, err := result.LastInsertId()
	if err != nil {
		return nil, errors.WithStack(err)
	}
	if lID == 0 {
		return // in case of non-insert statement
	}
	var j int64
	for _, arg := range rawArgs {
		switch a := arg.(type) {
		case LastInsertIDAssigner:
			a.AssignLastInsertID(lID + j)
			j++
		case QualifiedRecord:
			if rLIDA, ok := a.Record.(LastInsertIDAssigner); ok {
				rLIDA.AssignLastInsertID(lID + j)
				j++
			}
		}
	}
	return result, nil
}

// ExecValidateOneAffectedRow checks the sql.Result.RowsAffected if it returns
// one. If not returns an error of type NotValid. This function is
// useful for ExecContext function.
func ExecValidateOneAffectedRow(res sql.Result, err error) error {
	if err != nil {
		return errors.WithStack(err)
	}
	rowCount, err := res.RowsAffected()
	if err != nil {
		return errors.WithStack(err)
	}
	const expectedAffectedRows = 1
	if rowCount != expectedAffectedRows {
		return errors.NotValid.Newf("[dml] ExecValidateOneAffectedRow can't validate affected rows. Have %d Want %d", rowCount, expectedAffectedRows)
	}
	return nil
}
