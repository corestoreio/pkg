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
	"strings"
	"sync"
	"time"

	"github.com/corestoreio/errors"
	"github.com/corestoreio/log"
	"github.com/corestoreio/pkg/storage/null"
	"github.com/corestoreio/pkg/util/bufferpool"
	"golang.org/x/sync/errgroup"
)

// Artisan prepares the SQL string from a DML type, collects and build a list of
// arguments for later sending and execution in the database server. Arguments
// are collections of primitive types or slices of primitive types. An Artisan
// type acts like a prepared statement. In fact it can contain under the hood
// different connection types. Artisan is optimized for reuse and allow saving
// memory allocations.
type Artisan struct {
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
	insertRowCount      uint
	insertIsBuildValues bool
	// isPrepared if true the cachedSQL field in base gets ignored
	isPrepared bool
	// Options like enable interpolation or expanding placeholders.
	Options uint
	// hasNamedArgs checks if the SQL string in the cachedSQL field contains
	// named arguments. 0 not yet checked, 1=does not contain, 2 = yes
	hasNamedArgs      uint8 // 0 not checked, 1=no, 2=yes
	nextUnnamedArgPos int
	raw               []interface{} // set by developer
	arguments
	recs []QualifiedRecord
}

const (
	argOptionExpandPlaceholder = 1 << iota
	argOptionInterpolate
)

// OrderBy appends columns to the ORDER BY statement for ascending sorting. A
// column gets always quoted if it is a valid identifier otherwise it will be
// treated as an expression. This ORDER BY clause gets appended to the current
// internal cached SQL string independently if the SQL statement supports it or
// not or if there exists already an ORDER BY clause.
func (a *Artisan) OrderBy(columns ...string) *Artisan {
	a.OrderBys = a.OrderBys.AppendColumns(false, columns...)
	return a
}

// OrderByDesc appends columns to the ORDER BY statement for descending sorting.
// A column gets always quoted if it is a valid identifier otherwise it will be
// treated as an expression. This ORDER BY clause gets appended to the current
// internal cached SQL string independently if the SQL statement supports it or
// not or if there exists already an ORDER BY clause.
func (a *Artisan) OrderByDesc(columns ...string) *Artisan {
	a.OrderBys = a.OrderBys.AppendColumns(false, columns...).applySort(len(columns), sortDescending)
	return a
}

// Limit sets a LIMIT clause for the statement; overrides any existing LIMIT.
// This LIMIT clause gets appended to the current internal cached SQL string
// independently if the SQL statement supports it or not or if there exists
// already a LIMIT clause.
func (a *Artisan) Limit(offset uint64, limit uint64) *Artisan {
	a.OffsetCount = offset
	a.LimitCount = limit
	a.OffsetValid = true
	a.LimitValid = true
	return a
}

// Paginate sets LIMIT/OFFSET for the statement based on the given page/perPage
// Assumes page/perPage are valid. Page and perPage must be >= 1
func (a *Artisan) Paginate(page, perPage uint64) *Artisan {
	a.Limit((page-1)*perPage, perPage)
	return a
}

// WithQualifiedColumnsAliases for documentation please see:
// Artisan.QualifiedColumnsAliases.
func (a *Artisan) WithQualifiedColumnsAliases(aliases ...string) *Artisan {
	a.QualifiedColumnsAliases = aliases
	return a
}

// ToSQL the returned interface slice is owned by the callee.
func (a *Artisan) ToSQL() (string, []interface{}, error) {
	return a.prepareArgs()
}

func (a *Artisan) CachedQueries(queries ...string) []string {
	return a.base.CachedQueries(queries...)
}

// WithCacheKey sets the currently used cache key when generating a SQL string.
// By setting a different cache key, a previous generated SQL query is
// accessible again. New cache keys allow to change the generated query of the
// current object. E.g. different where clauses or different row counts in
// INSERT ... VALUES statements. The empty string defines the default cache key.
// If the `args` argument contains values, then fmt.Sprintf gets used.
func (a *Artisan) WithCacheKey(key string, args ...interface{}) *Artisan {
	a.base.withCacheKey(key, args...)
	return a
}

// Interpolate if set stringyfies the arguments into the SQL string and returns
// pre-processed SQL command when calling the function ToSQL. Not suitable for
// prepared statements. ToSQLs second argument `args` will then be nil.
func (a *Artisan) Interpolate() *Artisan {
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
func (a *Artisan) ExpandPlaceHolders() *Artisan {
	a.Options = a.Options | argOptionExpandPlaceholder
	return a
}

func (a *Artisan) isEmpty() bool {
	if a == nil {
		return true
	}
	return len(a.raw) == 0 && len(a.arguments) == 0 && len(a.recs) == 0
}

// prepareArgs transforms mainly the Artisan into []interface{}. It appends
// its arguments to the `extArgs` arguments from the Exec+ or Query+ function.
// This allows for a developer to reuse the interface slice and save
// allocations. All method receivers are not thread safe. The returned interface
// slice is the same as `extArgs`.
func (a *Artisan) prepareArgs(extArgs ...interface{}) (_ string, _ []interface{}, err error) {
	if a.base.ärgErr != nil {
		return "", nil, errors.WithStack(a.base.ärgErr)
	}

	if a.base.source == dmlSourceInsert {
		return a.prepareArgsInsert(extArgs...)
	}
	cachedSQL, ok := a.base.cachedSQL[a.base.CacheKey]
	if !a.isPrepared && !ok {
		return "", nil, errors.Empty.Newf("[dml] Artisan: The SQL string is empty.")
	}

	if a.isEmpty() { // no arguments provided
		if a.isPrepared {
			return "", extArgs, nil
		}

		if len(a.OrderBys) == 0 && !a.LimitValid {
			return cachedSQL, extArgs, nil
		}
		buf := bufferpool.Get()
		defer bufferpool.Put(buf)
		buf.WriteString(cachedSQL)
		sqlWriteOrderBy(buf, a.OrderBys, false)
		sqlWriteLimitOffset(buf, a.LimitValid, a.OffsetValid, a.OffsetCount, a.LimitCount)
		return buf.String(), extArgs, nil
	}

	if !a.isPrepared && a.hasNamedArgs == 0 {
		var found bool
		a.hasNamedArgs = 1
		cachedSQL, a.base.qualifiedColumns, found = extractReplaceNamedArgs(cachedSQL, a.base.qualifiedColumns)
		if found {
			a.base.cachedSQLUpsert(a.base.CacheKey, cachedSQL)
		}

		switch la := len(a.arguments); true {
		case found:
			a.hasNamedArgs = 2
		case !found && len(a.recs) == 0 && la > 0:
			for _, arg := range a.arguments {
				if sn, ok := arg.value.(sql.NamedArg); ok && sn.Name != "" {
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
		return "", collectedArgs.toInterfaces(extArgs...), nil
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
	if a.Options > 0 && len(extArgs) > 0 && len(a.recs) == 0 && len(a.arguments) == 0 {
		return "", nil, errors.NotAllowed.Newf("[dml] Interpolation/ExpandPlaceholders supports only Records and Arguments and not yet an interface slice.")
	}

	// TODO more advanced caching of the final non-expanded SQL string

	if a.Options&argOptionExpandPlaceholder != 0 {
		phCount := bytes.Count(sqlBuf.First.Bytes(), placeHolderByte)
		if aLen, hasSlice := a.totalSliceLen(); phCount < aLen || hasSlice {
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

	extArgs = collectedArgs.toInterfaces(extArgs...)
	return sqlBuf.First.String(), extArgs, nil
}

func (a *Artisan) appendConvertedRecordsToArguments(collectedArgs arguments) (arguments, error) {
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
	cm := NewColumnMap(len(a.arguments)+len(a.recs), "") // can use an arg pool Artisan sync.Pool, nope.

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
func (a *Artisan) prepareArgsInsert(extArgs ...interface{}) (string, []interface{}, error) {
	// cm := pooledColumnMapGet()
	sqlBuf := bufferpool.GetTwin()
	defer bufferpool.PutTwin(sqlBuf)
	// defer pooledBufferColumnMapPut(cm, sqlBuf, nil)

	cm := NewColumnMap(16)
	cm.setColumns(a.base.qualifiedColumns)
	// defer bufferpool.PutTwin(sqlBuf)
	cm.arguments = append(cm.arguments, a.arguments...)
	lenInsertCachedSQL := len(a.insertCachedSQL)
	cachedSQL, _ := a.base.cachedSQL[a.base.CacheKey]
	{
		if lenInsertCachedSQL > 0 {
			cachedSQL = a.insertCachedSQL
		}
		if _, err := sqlBuf.First.WriteString(cachedSQL); err != nil {
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
		return "", cm.arguments.toInterfaces(extArgs...), nil
	}

	extArgs = append(extArgs, a.raw...)
	totalArgLen := uint(len(cm.arguments) + len(extArgs))

	if !a.insertIsBuildValues && lenInsertCachedSQL == 0 { // Write placeholder list e.g. "VALUES (?,?),(?,?)"
		odkPos := strings.Index(cachedSQL, onDuplicateKeyPartS)
		if odkPos > 0 {
			sqlBuf.First.Reset()
			sqlBuf.First.WriteString(cachedSQL[:odkPos])
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
			sqlBuf.First.WriteString(cachedSQL[odkPos:])
		}
	}
	if lenInsertCachedSQL == 0 {
		a.insertCachedSQL = sqlBuf.First.String()
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
	return a.insertCachedSQL, cm.arguments.toInterfaces(extArgs...), nil
}

// nextUnnamedArg returns an unnamed argument by its position.
func (a *Artisan) nextUnnamedArg() (argument, bool) {
	var unnamedCounter int
	lenArg := len(a.arguments)
	for i := 0; i < lenArg && a.nextUnnamedArgPos >= 0; i++ {
		if _, ok := a.arguments[i].value.(sql.NamedArg); !ok {
			if unnamedCounter == a.nextUnnamedArgPos {
				a.nextUnnamedArgPos++
				return a.arguments[i], true
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
func (a *Artisan) MapColumns(cm *ColumnMap) error {
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
			if c != "" {
				if sn, ok := arg.value.(sql.NamedArg); ok && sn.Name == c {
					cm.arguments = append(cm.arguments, arg)
					break
				}
			}
		}
	}
	return cm.Err()
}

func (a *Artisan) add(v interface{}) *Artisan {
	if a == nil {
		a = &Artisan{arguments: make(arguments, 0, defaultArgumentsCapacity)}
	}
	a.arguments = a.arguments.add(v)
	return a
}

// Record appends a record for argument extraction. Qualifier is the name of the
// table or view or procedure or their alias name. It must be a valid
// MySQL/MariaDB identifier. An empty qualifier gets assigned to the main table.
func (a *Artisan) Record(qualifier string, record ColumnMapper) *Artisan {
	a.recs = append(a.recs, Qualify(qualifier, record))
	return a
}

func (a *Artisan) Raw(raw ...interface{}) *Artisan { a.raw = raw; return a }

func (a *Artisan) Null() *Artisan                           { return a.add(nil) }
func (a *Artisan) Int(i int) *Artisan                       { return a.add(i) }
func (a *Artisan) Ints(i ...int) *Artisan                   { return a.add(i) }
func (a *Artisan) Int64(i int64) *Artisan                   { return a.add(i) }
func (a *Artisan) Int64s(i ...int64) *Artisan               { return a.add(i) }
func (a *Artisan) Uint(i uint) *Artisan                     { return a.add(uint64(i)) }
func (a *Artisan) Uints(i ...uint) *Artisan                 { return a.add(i) }
func (a *Artisan) Uint64(i uint64) *Artisan                 { return a.add(i) }
func (a *Artisan) Uint64s(i ...uint64) *Artisan             { return a.add(i) }
func (a *Artisan) Float64(f float64) *Artisan               { return a.add(f) }
func (a *Artisan) Float64s(f ...float64) *Artisan           { return a.add(f) }
func (a *Artisan) Bool(b bool) *Artisan                     { return a.add(b) }
func (a *Artisan) Bools(b ...bool) *Artisan                 { return a.add(b) }
func (a *Artisan) String(s string) *Artisan                 { return a.add(s) }
func (a *Artisan) Strings(s ...string) *Artisan             { return a.add(s) }
func (a *Artisan) Time(t time.Time) *Artisan                { return a.add(t) }
func (a *Artisan) Times(t ...time.Time) *Artisan            { return a.add(t) }
func (a *Artisan) Bytes(b []byte) *Artisan                  { return a.add(b) }
func (a *Artisan) BytesSlice(b ...[]byte) *Artisan          { return a.add(b) }
func (a *Artisan) NullString(nv null.String) *Artisan       { return a.add(nv) }
func (a *Artisan) NullStrings(nv ...null.String) *Artisan   { return a.add(nv) }
func (a *Artisan) NullFloat64(nv null.Float64) *Artisan     { return a.add(nv) }
func (a *Artisan) NullFloat64s(nv ...null.Float64) *Artisan { return a.add(nv) }
func (a *Artisan) NullInt64(nv null.Int64) *Artisan         { return a.add(nv) }
func (a *Artisan) NullInt64s(nv ...null.Int64) *Artisan     { return a.add(nv) }
func (a *Artisan) NullBool(nv null.Bool) *Artisan           { return a.add(nv) }
func (a *Artisan) NullBools(nv ...null.Bool) *Artisan       { return a.add(nv) }
func (a *Artisan) NullTime(nv null.Time) *Artisan           { return a.add(nv) }
func (a *Artisan) NullTimes(nv ...null.Time) *Artisan       { return a.add(nv) }

// NamedArg converts to sql.NamedArg and as go-sql-driver/mysql does not (yet)
// support named args, they get resolved, converted to question mark place
// holders.
func (a *Artisan) NamedArg(name string, value interface{}) *Artisan {
	a.arguments = append(a.arguments, argument{isSet: true, value: sql.Named(name, value)})
	return a
}

// NamedArgs appends multiple
func (a *Artisan) NamedArgs(sns ...sql.NamedArg) *Artisan {
	for _, sn := range sns {
		a.arguments = append(a.arguments, argument{isSet: true, value: sn})
	}
	return a
}

// Reset resets the internal slices for new usage retaining the already
// allocated memory. Reset gets called automatically in many Load* functions. In
// case of an INSERT statement, Reset triggers a new build of the VALUES part.
// This function must be called when the number of argument changes.
func (a *Artisan) Reset() *Artisan {
	for i := range a.recs {
		a.recs[i].Qualifier = ""
		a.recs[i].Record = nil // remove pointers for GC
	}
	a.recs = a.recs[:0]
	a.arguments = a.arguments[:0]
	a.raw = a.raw[:0]
	a.nextUnnamedArgPos = 0
	a.insertIsBuildValues = false
	a.insertCachedSQL = a.insertCachedSQL[:0]
	return a
}

// DriverValue adds multiple of the same underlying values to the argument
// slice. When using different values, the last applied value wins and gets
// added to the argument slice. For example driver.Values of type `int` will
// result in []int.
func (a *Artisan) DriverValue(dvs ...driver.Valuer) *Artisan {
	if a == nil {
		a = &Artisan{arguments: make(arguments, 0, len(dvs))}
	}
	a.arguments, a.base.ärgErr = driverValue(a.arguments, dvs...)
	return a
}

// DriverValues adds each driver.Value as its own argument to the argument
// slice. It panics if the underlying type is not one of the allowed of
// interface driver.Valuer.
func (a *Artisan) DriverValues(dvs ...driver.Valuer) *Artisan {
	if a == nil {
		a = &Artisan{arguments: make(arguments, 0, len(dvs))}
	}
	a.arguments, a.base.ärgErr = driverValues(a.arguments, dvs...)
	return a
}

// WithDB sets the database query object.
func (a *Artisan) WithDB(db QueryExecPreparer) *Artisan {
	a.base.DB = db
	return a
}

// WithPreparedStmt uses a SQL statement as DB connection.
func (a *Artisan) WithPreparedStmt(stmt *sql.Stmt) *Artisan {
	a.base.DB = stmtWrapper{stmt: stmt}
	return a
}

// WithTx sets the transaction query executor and the logger to run this query
// within a transaction.
func (a *Artisan) WithTx(tx *Tx) *Artisan {
	if a.base.id == "" {
		a.base.id = tx.makeUniqueID()
	}
	a.base.Log = tx.Log
	a.base.DB = tx.DB
	return a
}

// Clone creates a shallow clone of the current pointer and sets the field `DB`
// to nil. The logger gets copied. Some underlying slices for the cached SQL
// statements are still referring to the source Artisan object.
func (a *Artisan) Clone() *Artisan {
	c := new(Artisan)
	*c = *a

	c.raw = nil
	c.arguments = make(arguments, 0, len(a.arguments))
	c.recs = nil
	c.base.DB = nil

	return c
}

// Close tries to close the underlying DB connection. Useful in cases of
// prepared statements. If the underlying DB connection does not implement
// io.Closer, nothing will happen.
func (a *Artisan) Close() error {
	if a.base.ärgErr != nil {
		return errors.WithStack(a.base.ärgErr)
	}
	if c, ok := a.base.DB.(ioCloser); ok {
		return errors.WithStack(c.Close())
	}
	return nil
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

const argumentPoolMaxSize = 32

// regarding the returned slices in both pools: https://github.com/golang/go/blob/7e394a2/src/net/http/h2_bundle.go#L998-L1043
// they also uses a []byte slice in the pool and not a pointer

var pooledArguments = sync.Pool{
	New: func() interface{} {
		var a [argumentPoolMaxSize]argument
		return arguments(a[:0])
	},
}

func pooledArgumentsGet() arguments {
	return pooledArguments.Get().(arguments)
}

func pooledArgumentsPut(a arguments, buf *bufferpool.TwinBuffer) {
	// @see https://go-review.googlesource.com/c/go/+/136116/4/src/fmt/print.go
	// Proper usage of a sync.Pool requires each entry to have approximately
	// the same memory cost. To obtain this property when the stored type
	// contains a variably-sized buffer, we add a hard limit on the maximum buffer
	// to place back in the pool.
	//
	// See https://golang.org/issue/23199
	if cap(a) <= argumentPoolMaxSize {
		a = a[:0]
		pooledArguments.Put(a)
	}
	if buf != nil {
		bufferpool.PutTwin(buf)
	}
}

var pooledInterfaces = sync.Pool{
	New: func() interface{} {
		var a [argumentPoolMaxSize]interface{}
		return a[:0]
	},
}

func pooledInterfacesGet() []interface{} {
	return pooledInterfaces.Get().([]interface{})
}

func pooledInterfacesPut(a []interface{}) {
	if cap(a) <= argumentPoolMaxSize {
		a = a[:0]
		pooledInterfaces.Put(a)
	}
}

// ExecContext executes the statement represented by the Update/Insert object.
// It returns the raw database/sql Result or an error if there was one.
// Regarding LastInsertID(): If you insert multiple rows using a single INSERT
// statement, LAST_INSERT_ID() returns the value generated for the first
// inserted row only. The reason for this at to make it possible to reproduce
// easily the same INSERT statement against some other server. If a record resp.
// and object implements the interface LastInsertIDAssigner then the
// LastInsertID gets assigned incrementally to the objects.
func (a *Artisan) ExecContext(ctx context.Context, args ...interface{}) (sql.Result, error) {
	return a.exec(ctx, args...)
}

// QueryContext traditional way of the databasel/sql package.
func (a *Artisan) QueryContext(ctx context.Context, args ...interface{}) (*sql.Rows, error) {
	return a.query(ctx, args...)
}

// QueryRowContext traditional way of the databasel/sql package.
func (a *Artisan) QueryRowContext(ctx context.Context, args ...interface{}) *sql.Row {
	sqlStr, args, err := a.prepareArgs(args...)
	if a.base.Log != nil && a.base.Log.IsDebug() {
		defer log.WhenDone(a.base.Log).Debug("QueryRowContext", log.String("sql", sqlStr), log.String("source", string(a.base.source)), log.Err(err))
	}
	return a.base.DB.QueryRowContext(ctx, sqlStr, args...)
}

// IterateSerial iterates in serial order over the result set by loading one row each
// iteration and then discarding it. Handles records one by one. The context
// gets only used in the Query function.
func (a *Artisan) IterateSerial(ctx context.Context, callBack func(*ColumnMap) error, args ...interface{}) (err error) {
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
func (a *Artisan) IterateParallel(ctx context.Context, concurrencyLevel int, callBack func(*ColumnMap) error, args ...interface{}) (err error) {
	if a.base.Log != nil && a.base.Log.IsDebug() {
		defer log.WhenDone(a.base.Log).Debug("IterateParallel", log.String("id", a.base.id), log.Err(err))
	}
	if concurrencyLevel < 1 {
		return errors.OutOfRange.Newf("[dml] Artisan.IterateParallel concurrencyLevel %d for query ID %q cannot be smaller zero.", concurrencyLevel, a.base.id)
	}

	r, err := a.query(ctx, args...)
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
// muliple-rows. It checks on top if ColumnMapper `s` implements io.Closer, to
// call the custom close function. This is useful for e.g. unlocking a mutex.
func (a *Artisan) Load(ctx context.Context, s ColumnMapper, args ...interface{}) (rowCount uint64, err error) {
	if a.base.Log != nil && a.base.Log.IsDebug() {
		defer log.WhenDone(a.base.Log).Debug("Load", log.String("id", a.base.id), log.Err(err), log.ObjectTypeOf("ColumnMapper", s), log.Uint64("row_count", rowCount))
	}

	r, err := a.query(ctx, args...)
	if err != nil {
		err = errors.Wrapf(err, "[dml] Artisan.Load.QueryContext failed with queryID %q and ColumnMapper %T", a.base.id, s)
		return
	}
	cm := pooledColumnMapGet()
	defer pooledBufferColumnMapPut(cm, nil, func() {
		a.Reset()
		// Not testable with the sqlmock package :-(
		if err2 := r.Close(); err2 != nil && err == nil {
			err = errors.Wrap(err2, "[dml] Artisan.Load.Rows.Close")
		}
		if rc, ok := s.(ioCloser); ok {
			if err2 := rc.Close(); err2 != nil && err == nil {
				err = errors.Wrap(err2, "[dml] Artisan.Load.ColumnMapper.Close")
			}
		}
	})

	for r.Next() {
		if err = cm.Scan(r); err != nil {
			return 0, errors.WithStack(err)
		}
		if err = s.MapColumns(cm); err != nil {
			return 0, errors.Wrapf(err, "[dml] Artisan.Load failed with queryID %q and ColumnMapper %T", a.base.id, s)
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
func (a *Artisan) LoadNullInt64(ctx context.Context, args ...interface{}) (nv null.Int64, found bool, err error) {
	found, err = a.loadPrimitive(ctx, &nv, args...)
	return
}

// LoadNullUint64 executes the query and returns the first row parsed into the
// current type. `Found` might be false if there are no matching rows.
// This function with ptr type uint64 comes in handy when performing
// a COUNT(*) query. See function `Select.Count`.
func (a *Artisan) LoadNullUint64(ctx context.Context, args ...interface{}) (nv null.Uint64, found bool, err error) {
	found, err = a.loadPrimitive(ctx, &nv, args...)
	return
}

// LoadNullFloat64 executes the query and returns the first row parsed into the
// current type. `Found` might be false if there are no matching rows.
func (a *Artisan) LoadNullFloat64(ctx context.Context, args ...interface{}) (nv null.Float64, found bool, err error) {
	found, err = a.loadPrimitive(ctx, &nv, args...)
	return
}

// LoadNullString executes the query and returns the first row parsed into the
// current type. `Found` might be false if there are no matching rows.
func (a *Artisan) LoadNullString(ctx context.Context, args ...interface{}) (nv null.String, found bool, err error) {
	found, err = a.loadPrimitive(ctx, &nv, args...)
	return
}

// LoadNullTime executes the query and returns the first row parsed into the
// current type. `Found` might be false if there are no matching rows.
func (a *Artisan) LoadNullTime(ctx context.Context, args ...interface{}) (nv null.Time, found bool, err error) {
	found, err = a.loadPrimitive(ctx, &nv, args...)
	return
}

// LoadDecimal executes the query and returns the first row parsed into the
// current type. `Found` might be false if there are no matching rows.
func (a *Artisan) LoadDecimal(ctx context.Context, args ...interface{}) (nv null.Decimal, found bool, err error) {
	found, err = a.loadPrimitive(ctx, &nv, args...)
	return
}

func (a *Artisan) loadPrimitive(ctx context.Context, ptr interface{}, args ...interface{}) (found bool, err error) {
	if a.base.Log != nil && a.base.Log.IsDebug() {
		// do not use fullSQL because we might log sensitive data
		defer log.WhenDone(a.base.Log).Debug("LoadPrimitive", log.String("id", a.base.id), log.Err(err), log.ObjectTypeOf("ptr_type", ptr))
	}
	var rows *sql.Rows
	rows, err = a.query(ctx, args...)
	if err != nil {
		err = errors.WithStack(err)
		return
	}
	defer func() {
		a.Reset() // reset the internal slices to avoid adding more and more arguments when a query gets executed
		if errC := rows.Close(); err == nil && errC != nil {
			err = errors.WithStack(errC)
		}
	}()

	for rows.Next() && !found {
		if err = rows.Scan(ptr); err != nil {
			err = errors.WithStack(err)
			return
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
func (a *Artisan) LoadInt64s(ctx context.Context, dest []int64, args ...interface{}) (_ []int64, err error) {
	var rowCount int
	if a.base.Log != nil && a.base.Log.IsDebug() {
		// do not use fullSQL because we might log sensitive data
		defer log.WhenDone(a.base.Log).Debug("LoadInt64s", log.Int("row_count", rowCount), log.Err(err))
	}
	var r *sql.Rows
	r, err = a.query(ctx, args...)
	if err != nil {
		err = errors.WithStack(err)
		return
	}
	defer func() {
		a.Reset() // reset the internal slices to avoid adding more and more arguments when a query gets executed
		if cErr := r.Close(); err == nil && cErr != nil {
			err = errors.WithStack(cErr)
		}
	}()
	for r.Next() {
		var nv null.Int64
		if err = r.Scan(&nv); err != nil {
			err = errors.WithStack(err)
			return
		}
		if nv.Valid {
			dest = append(dest, nv.Int64)
		}
	}
	if err = r.Err(); err != nil {
		err = errors.WithStack(err)
		return
	}

	rowCount = len(dest)
	return dest, err
}

// LoadUint64s executes the query and returns the values appended to slice
// dest. It ignores and skips NULL values.
func (a *Artisan) LoadUint64s(ctx context.Context, dest []uint64, args ...interface{}) (_ []uint64, err error) {
	var rowCount int
	if a.base.Log != nil && a.base.Log.IsDebug() {
		// do not use fullSQL because we might log sensitive data
		defer log.WhenDone(a.base.Log).Debug("LoadUint64s", log.Int("row_count", rowCount), log.String("id", a.base.id), log.Err(err))
	}

	rows, err := a.query(ctx, args...)
	if err != nil {
		err = errors.WithStack(err)
		return
	}
	defer func() {
		a.Reset() // reset the internal slices to avoid adding more and more arguments when a query gets executed
		if errC := rows.Close(); errC != nil && err == nil {
			err = errors.WithStack(errC)
		}
	}()

	for rows.Next() {
		var nv null.Uint64
		if err = rows.Scan(&nv); err != nil {
			err = errors.WithStack(err)
			return
		}
		if nv.Valid {
			dest = append(dest, nv.Uint64)
		}
	}
	if err = rows.Err(); err != nil {
		err = errors.WithStack(err)
		return
	}
	rowCount = len(dest)
	return dest, err
}

// LoadFloat64s executes the query and returns the values appended to slice
// dest. It ignores and skips NULL values.
func (a *Artisan) LoadFloat64s(ctx context.Context, dest []float64, args ...interface{}) (_ []float64, err error) {
	if a.base.Log != nil && a.base.Log.IsDebug() {
		// do not use fullSQL because we might log sensitive data
		defer log.WhenDone(a.base.Log).Debug("LoadFloat64s", log.String("id", a.base.id), log.Err(err))
	}

	var rows *sql.Rows
	if rows, err = a.query(ctx, args...); err != nil {
		err = errors.WithStack(err)
		return
	}
	defer func() {
		a.Reset() // reset the internal slices to avoid adding more and more arguments when a query gets executed
		if errC := rows.Close(); errC != nil && err == nil {
			err = errors.WithStack(errC)
		}
	}()

	for rows.Next() {
		var nv null.Float64
		if err = rows.Scan(&nv); err != nil {
			err = errors.WithStack(err)
			return
		}
		if nv.Valid {
			dest = append(dest, nv.Float64)
		}
	}
	if err = rows.Err(); err != nil {
		err = errors.WithStack(err)
		return
	}
	return dest, err
}

// LoadStrings executes the query and returns the values appended to slice
// dest. It ignores and skips NULL values.
func (a *Artisan) LoadStrings(ctx context.Context, dest []string, args ...interface{}) (_ []string, err error) {
	var rowCount int
	if a.base.Log != nil && a.base.Log.IsDebug() {
		// do not use fullSQL because we might log sensitive data
		defer log.WhenDone(a.base.Log).Debug("LoadStrings", log.Int("row_count", rowCount), log.String("id", a.base.id), log.Err(err))
	}

	rows, err := a.query(ctx, args...)
	if err != nil {
		err = errors.WithStack(err)
		return
	}
	defer func() {
		a.Reset() // reset the internal slices to avoid adding more and more arguments when a query gets executed
		if errC := rows.Close(); errC != nil && err == nil {
			err = errors.WithStack(errC)
		}
	}()

	for rows.Next() {
		var value null.String
		if err = rows.Scan(&value); err != nil {
			err = errors.WithStack(err)
			return
		}
		if value.Valid {
			dest = append(dest, value.String)
		}
	}
	if err = rows.Err(); err != nil {
		err = errors.WithStack(err)
		return
	}
	rowCount = len(dest)
	return dest, err
}

func (a *Artisan) query(ctx context.Context, args ...interface{}) (rows *sql.Rows, err error) {
	pArgs := pooledInterfacesGet()
	defer pooledInterfacesPut(pArgs)
	pArgs = append(pArgs, args...)

	var sqlStr string
	sqlStr, pArgs, err = a.prepareArgs(pArgs...)
	if a.base.Log != nil && a.base.Log.IsDebug() {
		defer log.WhenDone(a.base.Log).Debug("Query", log.String("sql", sqlStr), log.String("source", string(a.base.source)), log.Err(err))
	}
	if err != nil {
		return nil, errors.WithStack(err)
	}
	rows, err = a.base.DB.QueryContext(ctx, sqlStr, pArgs...)
	if err != nil {
		if sqlStr == "" {
			cachedSQL, _ := a.base.cachedSQL[a.base.CacheKey]
			sqlStr = "PREPARED:" + cachedSQL
		}
		err = errors.Wrapf(err, "[dml] Query.QueryContext with query %q", sqlStr)
	}
	return
}

func (a *Artisan) exec(ctx context.Context, args ...interface{}) (result sql.Result, err error) {
	pArgs := pooledInterfacesGet()
	defer pooledInterfacesPut(pArgs)
	pArgs = append(pArgs, args...)

	var sqlStr string
	sqlStr, pArgs, err = a.prepareArgs(pArgs...)
	if a.base.Log != nil && a.base.Log.IsDebug() {
		defer log.WhenDone(a.base.Log).Debug("Exec", log.String("sql", sqlStr), log.String("source", string(a.base.source)), log.Err(err))
	}
	if err != nil {
		return nil, errors.WithStack(err)
	}

	result, err = a.base.DB.ExecContext(ctx, sqlStr, pArgs...)
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
	if lID == 0 {
		return // in case of non-insert statement
	}
	for i, rec := range a.recs {
		if a, ok := rec.Record.(LastInsertIDAssigner); ok {
			a.AssignLastInsertID(lID + int64(i))
		}
	}
	return
}
