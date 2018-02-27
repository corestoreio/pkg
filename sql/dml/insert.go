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
	"fmt"
	"sync"

	"github.com/corestoreio/errors"
	"github.com/corestoreio/log"
)

// LastInsertIDAssigner assigns the last insert ID of an auto increment
// column back to the objects.
type LastInsertIDAssigner interface {
	AssignLastInsertID(int64)
}

// Insert contains the clauses for an INSERT statement
type Insert struct {
	BuilderBase
	Into    string
	Columns []string
	// RowCount defines the number of expected rows.
	RowCount int // See SetRowCount()
	// RecordPlaceHolderCount defines the number of place holders for each set
	// within the brackets. Must only be set when Records have been applied
	// and `Columns` field has been omitted.
	RecordPlaceHolderCount int
	// Select used to create an "INSERT INTO `table` SELECT ..." statement.
	Select *Select
	Pairs  Conditions
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
	// IsBuildValues if true the VALUES part gets build when calling ToSQL.
	// VALUES do not need to get build by default because mostly WithArgs gets
	// called to build the VALUES part dynamically.
	IsBuildValues bool
	// Listeners allows to dispatch certain functions in different
	// situations.
	Listeners ListenersInsert
}

// NewInsert creates a new Insert object.
func NewInsert(into string) *Insert {
	var rwmu sync.RWMutex
	return &Insert{
		BuilderBase: BuilderBase{
			rwmu: &rwmu,
		},
		Into: into,
	}
}

func newInsertInto(db QueryExecPreparer, cCom *connCommon, into string) *Insert {
	id := cCom.makeUniqueID()
	into = cCom.mapTableName(into)
	l := cCom.Log
	if l != nil {
		l = l.With(log.String("insert_id", id), log.String("table", into))
	}
	var rwmu sync.RWMutex
	return &Insert{
		BuilderBase: BuilderBase{
			rwmu: &rwmu,
			builderCommon: builderCommon{
				id:  id,
				Log: l,
				DB:  db,
			},
		},
		Into: into,
	}
}

// InsertInto instantiates a Insert for the given table. Mapping the table name
// is supported.
func (c *ConnPool) InsertInto(into string) *Insert {
	return newInsertInto(c.DB, &c.connCommon, into)
}

// InsertInto instantiates a Insert for the given table. Mapping the table name
// is supported.
func (c *Conn) InsertInto(into string) *Insert {
	return newInsertInto(c.DB, &c.connCommon, into)
}

// InsertInto instantiates a Insert for the given table bound to a transaction.
// Mapping the table name is supported.
func (tx *Tx) InsertInto(into string) *Insert {
	return newInsertInto(tx.DB, &tx.connCommon, into)
}

// WithDB sets the database query object.
func (b *Insert) WithDB(db QueryExecPreparer) *Insert {
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

// BuildValues see IsBuildValues.
func (b *Insert) BuildValues() *Insert {
	b.IsBuildValues = true
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

// WithPairs appends a column/value pair to the statement. Calling this function
// multiple times with the same column name produces next rows for insertion.
// Slice values and right/left side expressions are not supported and ignored.
// You must call WithArgs afterwards to activate the `Pairs`.
func (b *Insert) WithPairs(cvs ...*Condition) *Insert {
	b.Pairs = append(b.Pairs, cvs...)
	return b
}

// FromSelect creates an "INSERT INTO `table` SELECT ..." statement from a
// previously created SELECT statement.
func (b *Insert) FromSelect(s *Select) *Insert {
	b.Select = s
	return b
}

// WithArgs returns a new Artisan type to support multiple executions of the
// underlying SQL statement and reuse of memory allocations for the arguments.
// WithArgs builds the SQL string in a thread safe way. It copies the underlying
// connection and settings from the current DML type (Delete, Insert, Select,
// Update, Union, With, etc.). The field DB can still be overwritten.
// Interpolation does not support the raw interfaces. It's an architecture bug
// to use WithArgs inside a loop.
// In case of INSERT statement, WithArgs figures automatically out how the
// VALUES section must look like depending on the number of arguments. In some
// cases type Insert needs to know the RowCount to build the appropriate amount
// of placeholders.
func (b *Insert) WithArgs() *Artisan {
	var pairArgs arguments
	b.rwmu.RLock()
	isSelect := b.Select != nil // b.withArtisan unsets the Select field if caching is enabled
	for _, cv := range b.Pairs {
		pairArgs = append(pairArgs, cv.Right.arg)
	}
	b.rwmu.RUnlock()

	a := b.withArtisan(b)
	a.base.source = dmlSourceInsert

	if isSelect {
		// Must change to this source because to trigger a different argument
		// collector in Artisan.prepareArgs. It is not a real INSERT statement
		// anymore.
		a.base.source = dmlSourceInsertSelect
		return a
	}

	a.arguments = append(a.arguments, pairArgs...)
	a.insertColumnCount = uint(len(b.Columns))
	if b.RecordPlaceHolderCount > 0 {
		a.insertColumnCount = uint(b.RecordPlaceHolderCount)
	}
	a.insertRowCount = uint(b.RowCount)
	a.insertIsBuildValues = b.IsBuildValues
	return a
}

// ToSQL serialized the Insert to a SQL string
// It returns the string with placeholders and a slice of query arguments
func (b *Insert) ToSQL() (string, []interface{}, error) {
	b.source = dmlSourceInsert
	rawSQL, err := b.buildToSQL(b)
	if err != nil {
		return "", nil, errors.WithStack(err)
	}
	return string(rawSQL), nil, nil
}

func (b *Insert) writeBuildCache(sql []byte, qualifiedColumns []string) {
	if !b.IsBuildCacheDisabled {
		b.cachedSQL = sql
		b.Select = nil
		b.Pairs = nil
		b.OnDuplicateKeys = nil
		b.OnDuplicateKeyExclude = nil
	}
	b.qualifiedColumns = qualifiedColumns
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

	for _, cv := range b.Pairs {
		if !strInSlice(cv.Left, b.Columns) {
			b.Columns = append(b.Columns, cv.Left)
		}
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
		if err != nil {
			return nil, errors.WithStack(err)
		}
		return b.writeOnDuplicateKey(buf, ph)
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

	if argCount0 := len(b.Columns); argCount0 > 0 && b.IsBuildValues {
		rowCount := 1
		if b.RowCount > 0 {
			rowCount = b.RowCount
		}
		if b.RecordPlaceHolderCount > 0 {
			argCount0 = b.RecordPlaceHolderCount
		}

		if lPairs := len(b.Pairs); lPairs > 0 { // monster IF, must be refactored
			if argCount0 > 0 {
				rowCount = argCount0
			}
			buf.WriteByte('(')
			for i := 0; i < argCount0 && lPairs <= rowCount; i++ {
				buf.WriteString("?,")
			}

			for i, cv := range b.Pairs {
				if i > 0 && i%rowCount == 0 {
					buf.WriteString("),(")
				} else if i > 0 {
					buf.WriteByte(',')
				}
				switch {
				case cv.Right.arg.isSet:
					cv.Right.arg.writeTo(buf, 0)
				case cv.Right.IsExpression:
					buf.WriteString(cv.Right.Column)
				case cv.Right.Sub != nil:
					var err error
					buf.WriteByte('(')
					placeHolders, err = cv.Right.Sub.toSQL(buf, placeHolders)
					if err != nil {
						return nil, errors.WithStack(err)
					}
					buf.WriteByte(')')
				default:
					fmt.Printf("%#v\n\n", cv.Right)
				}
			}
			buf.WriteByte(')')
		} else {
			writeInsertPlaceholders(buf, uint(rowCount), uint(argCount0))
		}
	}

	return b.writeOnDuplicateKey(buf, placeHolders)
}

func (b *Insert) writeOnDuplicateKey(buf *bytes.Buffer, placeHolders []string) ([]string, error) {
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

	return b.OnDuplicateKeys.writeOnDuplicateKey(buf, placeHolders)
}

func strInSlice(search string, sl []string) bool {
	for _, s := range sl {
		if s == search {
			return true
		}
	}
	return false
}

// Prepare executes the statement represented by the Insert to create a prepared
// statement. It returns a custom statement type or an error if there was one.
// Provided arguments or recs in the Insert are getting ignored. The provided
// context is used for the preparation of the statement, not for the execution
// of the statement. The returned Stmter is not safe for concurrent use, despite
// the underlying *sql.Stmt is.
func (b *Insert) Prepare(ctx context.Context) (*Stmt, error) {
	return b.prepare(ctx, b.DB, b, dmlSourceInsert)
}

// Clone creates a clone of the current object, leaving fields DB and Log
// untouched.
func (b *Insert) Clone() *Insert {
	if b == nil {
		return nil
	}
	c := *b
	c.BuilderBase = b.BuilderBase.Clone()
	c.Columns = cloneStringSlice(b.Columns)
	c.OnDuplicateKeyExclude = cloneStringSlice(b.OnDuplicateKeyExclude)
	c.OnDuplicateKeys = b.OnDuplicateKeys.Clone()
	c.Select = b.Select.Clone()
	c.Pairs = b.Pairs.Clone()
	return &c
}
