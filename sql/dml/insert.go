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
	"fmt"

	"github.com/corestoreio/errors"
)

/*
TODO: Using big transactions When doing many inserts in a row, you should wrap
 them with BEGIN / END to avoid doing a full transaction (which includes a disk
 sync) for every row. For example, doing a begin/end every 1000 inserts will
speed up your inserts by almost 1000 times.
BEGIN;
INSERT ...
INSERT ...
END;
BEGIN;
INSERT ...
INSERT ...
END;
...
The reason why you may want to have many BEGIN/END statements instead of just
one is that the former will use up less transaction log space.

Multi-value inserts
You can insert many rows at once with multi-value row inserts:

INSERT INTO table_name values(1,"row 1"),(2, "row 2"),...;
The limit for how much data you can have in one statement is controlled by the
max_allowed_packet server variable.
*/

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
	// OnDuplicateKeys updates the referenced columns. See documentation for
	// type `Conditions`. For more details
	// https://dev.mysql.com/doc/refman/5.7/en/insert-on-duplicate.html
	// Conditions contains the column/argument association for either the SET
	// clause in an UPDATE statement or to be used in an INSERT ... ON DUPLICATE
	// KEY statement. For each column there must be one argument which can
	// either be nil or has an actual value.
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
	// VALUES do not need to get build by default because mostly WithDBR gets
	// called to build the VALUES part dynamically.
	IsBuildValues bool
}

// NewInsert creates a new Insert object.
func NewInsert(into string) *Insert {
	return &Insert{
		Into: into,
	}
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
// Only needed in case the SQL string gets build without any arguments.
//		INSERT INTO tableX (?,?,?)
// SetRecordPlaceHolderCount would now be 3 because of the three place holders.
func (b *Insert) SetRecordPlaceHolderCount(valueCount int) *Insert {
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
// multiple times with the same column name will trigger an error.
// Slice values and right/left side expressions are not supported and ignored.
// You must call WithDBR afterwards to activate the `Pairs`.
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

// ToSQL serialized the Insert to a SQL string
// It returns the string with placeholders and a slice of query arguments.
func (b *Insert) ToSQL() (string, []interface{}, error) {
	rawSQL, err := b.buildToSQL(b)
	if err != nil {
		return "", nil, errors.WithStack(err)
	}
	return rawSQL, nil, nil
}

func (b *Insert) toSQL(buf *bytes.Buffer, placeHolders []string) ([]string, error) {
	for _, cv := range b.Pairs {
		if !strInSlice(cv.Left, b.Columns) {
			b.Columns = append(b.Columns, cv.Left)
		}
	}

	if b.Into == "" {
		return nil, errors.Empty.Newf("[dml] Inserted table is missing")
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
				case cv.Right.arg != nil:
					if err := writeInterfaceValue(cv.Right.arg, buf, 0); err != nil {
						return nil, errors.WithStack(err)
					}
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
			writeTuplePlaceholders(buf, uint(rowCount), uint(argCount0))
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
