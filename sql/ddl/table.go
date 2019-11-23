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

package ddl

import (
	"bytes"
	"context"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/corestoreio/errors"
	"github.com/corestoreio/log"
	"github.com/corestoreio/pkg/sql/dml"
	"github.com/corestoreio/pkg/util/bufferpool"
)

// Table represents a table from a specific database with a bound default
// connection pool. Its fields are not secure to use in concurrent context and
// hence might cause a race condition if used not properly.
type Table struct {
	// dcp represents a Database Connection Pool. Shall not be nil
	dcp      *dml.ConnPool
	customDB dml.QueryExecPreparer // if set we have a shallow copy
	// Schema represents the name of the database. Might be empty.
	Schema string
	// Name of the table
	Name string
	// Columns all table columns. They do not get used to create or alter a
	// table.
	Columns Columns
	// IsView set to true to mark if the table is a view.
	IsView bool
	// optimized column selection for specific DML operations.
	columnsPK    []string // only primary key columns
	columnsNonPK []string // all columns, except PK and system-versioned
	columnsAll   []string // all columns, except system-versioned
	// columnsIsEligibleForUpsert contains all non-current-timestamp, non-virtual, non-system
	// versioned and non auto_increment columns for update or insert operations.
	columnsUpsert []string
	colset        map[string]struct{}
}

// NewTable initializes a new table structure
func NewTable(tableName string, cs ...*Column) *Table {
	ts := &Table{
		Name:    tableName,
		Columns: Columns(cs),
	}
	return ts.update()
}

// update recalculates the internal cached columns
func (t *Table) update() *Table {
	if t.Columns.Len() == 0 {
		return t
	}

	t.columnsNonPK = t.columnsNonPK[:0]
	t.columnsNonPK = t.Columns.NonPrimaryColumns().FieldNames(t.columnsNonPK...)

	t.columnsPK = t.columnsPK[:0]
	t.columnsPK = t.Columns.PrimaryKeys().FieldNames(t.columnsPK...)

	t.columnsAll = t.columnsAll[:0]
	t.columnsAll = t.Columns.Filter(colIsNotSysVers).FieldNames(t.columnsAll...)

	t.columnsUpsert = t.columnsUpsert[:0]
	t.columnsUpsert = t.Columns.Filter(columnsIsEligibleForUpsert).FieldNames(t.columnsUpsert...)

	if t.colset == nil {
		t.colset = make(map[string]struct{}, t.Columns.Len())
	}
	t.Columns.Each(func(c *Column) {
		t.colset[c.Field] = struct{}{}
	})

	return t
}

// WithDB creates a shallow clone of the current Table object and uses argument
// `db` as the current connection. `db` can be a connection pool, a single
// connection or a transaction. This method might cause a race condition if use
// not properly. One shall not modify the slices in the returned *Table.
func (t *Table) WithDB(db dml.QueryExecPreparer) *Table {
	t2 := new(Table)
	*t2 = *t // dereference pointer and copy object. retains the slice storage.
	t2.customDB = db
	return t2
}

// Insert creates a new INSERT statement with all non primary key columns. If
// OnDuplicateKey() gets called, the INSERT can be used as an update or create
// statement. Adding multiple VALUES section is allowed. Using this statement to
// prepare a query, a call to `BuildValues()` triggers building the VALUES
// clause, otherwise a SQL parse error will occur.
func (t *Table) Insert() *dml.Insert {
	i := t.dcp.InsertInto(t.Name).AddColumns(t.columnsUpsert...)
	if t.customDB != nil {
		i.DB = t.customDB
	}
	return i
}

// Select creates a new SELECT statement. If "*" gets set as an argument, then
// all columns will be added to to list of columns.
func (t *Table) Select(columns ...string) *dml.Select {
	if len(columns) == 1 && columns[0] == "*" {
		columns = t.columnsAll
	}
	s := t.dcp.SelectFrom(t.Name, MainTable).AddColumns(columns...)
	if t.customDB != nil {
		s.DB = t.customDB
	}
	return s
}

// SelectByPK creates a new `SELECT columns FROM table WHERE id IN (?)`. If "*"
// gets set as an argument, then all columns will be added to to list of
// columns.
func (t *Table) SelectByPK(columns ...string) *dml.Select {
	if len(columns) == 1 && columns[0] == "*" {
		columns = t.columnsAll
	}
	s := t.dcp.SelectFrom(t.Name, MainTable).AddColumns(columns...)
	s.Wheres = t.whereByPK(dml.In)
	if t.customDB != nil {
		s.DB = t.customDB
	}
	return s
}

// DeleteByPK creates a new `DELETE FROM table WHERE id IN (?)`
func (t *Table) DeleteByPK() *dml.Delete {
	d := t.dcp.DeleteFrom(t.Name)
	d.Wheres = t.whereByPK(dml.In)
	if t.customDB != nil {
		d.DB = t.customDB
	}
	return d
}

// Delete creates a new `DELETE FROM table` statement.
func (t *Table) Delete() *dml.Delete {
	d := t.dcp.DeleteFrom(t.Name)
	if t.customDB != nil {
		d.DB = t.customDB
	}
	return d
}

// UpdateByPK creates a new `UPDATE table SET ... WHERE id = ?`. The SET clause
// contains all non primary columns.
func (t *Table) UpdateByPK() *dml.Update {
	u := t.dcp.Update(t.Name).AddColumns(t.columnsUpsert...)
	u.Wheres = t.whereByPK(dml.Equal)
	if t.customDB != nil {
		u.DB = t.customDB
	}
	return u
}

func (t *Table) whereByPK(op dml.Op) dml.Conditions {
	cnds := make(dml.Conditions, 0, 1)
	for _, pk := range t.columnsPK {
		c := dml.Column(pk).PlaceHolder()
		c.Operator = op
		cnds = append(cnds, c)
	}
	return cnds
}

func (t *Table) runExec(ctx context.Context, qry string) error {
	var db dml.QueryExecPreparer
	if t.dcp != nil {
		db = t.dcp.DB
	}
	if t.customDB != nil {
		db = t.customDB
	}
	if _, err := db.ExecContext(ctx, qry); err != nil {
		return errors.Wrapf(err, "[ddl] failed to exec %q", qry) // please do change this return signature, saves an alloc
	}
	return nil
}

// Truncate truncates the table. Removes all rows and sets the auto increment to
// zero. Just like a CREATE TABLE statement. To use a custom connection, call
// WithDB before.
func (t *Table) Truncate(ctx context.Context) error {
	if t.IsView {
		return nil
	}
	if err := dml.IsValidIdentifier(t.Name); err != nil {
		return errors.WithStack(err)
	}
	return t.runExec(ctx, "TRUNCATE TABLE "+dml.Quoter.QualifierName(t.Schema, t.Name))
}

// Rename renames the current table to the new table name. Renaming is an atomic
// operation in the database. As long as two databases are on the same file
// system, you can use RENAME TABLE to move a table from one database to
// another. RENAME TABLE also works for views, as long as you do not try to
// rename a view into a different database. To use a custom connection, call
// WithDB before.
func (t *Table) Rename(ctx context.Context, newTableName string) error {
	if err := dml.IsValidIdentifier(t.Name); err != nil {
		return errors.WithStack(err)
	}
	if err := dml.IsValidIdentifier(newTableName); err != nil {
		return errors.WithStack(err)
	}
	return t.runExec(ctx, "RENAME TABLE "+dml.Quoter.QualifierName(t.Schema, t.Name)+" TO "+dml.Quoter.NameAlias(newTableName, ""))
}

// Swap swaps the current table with the other table of the same structure.
// Renaming is an atomic operation in the database. Note: indexes won't get
// swapped! As long as two databases are on the same file system, you can use
// RENAME TABLE to move a table from one database to another. To use a custom
// connection, call WithDB before.
func (t *Table) Swap(ctx context.Context, other string) error {
	if err := dml.IsValidIdentifier(t.Name); err != nil {
		return errors.WithStack(err)
	}
	if err := dml.IsValidIdentifier(other); err != nil {
		return errors.WithStack(err)
	}

	tmp := TableName("", t.Name, strconv.FormatInt(time.Now().UnixNano(), 10))

	buf := bufferpool.Get()
	defer bufferpool.Put(buf)
	buf.WriteString("RENAME TABLE ")
	dml.Quoter.WriteQualifierName(buf, t.Schema, t.Name)
	buf.WriteString(" TO ")
	dml.Quoter.WriteIdentifier(buf, tmp)
	buf.WriteString(", ")
	dml.Quoter.WriteIdentifier(buf, other)
	buf.WriteString(" TO ")
	dml.Quoter.WriteQualifierName(buf, t.Schema, t.Name)
	buf.WriteByte(',')
	dml.Quoter.WriteIdentifier(buf, tmp)
	buf.WriteString(" TO ")
	dml.Quoter.WriteIdentifier(buf, other)
	return t.runExec(ctx, buf.String())
}

func (t *Table) getTyp() string {
	if t.IsView || strings.HasPrefix(t.Name, PrefixView) {
		return "VIEW"
	}
	return "TABLE"
}

// Drop drops, if exists, the table or the view. To use a custom connection,
// call WithDB before.
func (t *Table) Drop(ctx context.Context) error {
	if err := dml.IsValidIdentifier(t.Name); err != nil {
		return errors.Wrap(err, "[ddl] Drop table name")
	}
	return t.runExec(ctx, "DROP "+t.getTyp()+" IF EXISTS "+dml.Quoter.QualifierName(t.Schema, t.Name))
}

// HasColumn uses the internal cache to check if a column exists in a table and
// if so returns true. Case sensitive.
func (t *Table) HasColumn(columnName string) bool {
	_, ok := t.colset[columnName]
	return ok
}

// InfileOptions provides options for the function LoadDataInfile. Some columns
// are self-describing.
type InfileOptions struct {
	// IsNotLocal disables LOCAL load file. If LOCAL is specified, the file is read
	// by the client program on the client host and sent to the server. If LOCAL
	// is not specified, the file must be located on the server host and is read
	// directly by the server.
	// See security issues in https://dev.mysql.com/doc/refman/5.7/en/load-data-local.html
	IsNotLocal bool
	// Replace, input rows replace existing rows. In other words, rows that have
	// the same value for a primary key or unique index as an existing row.
	Replace bool
	// Ignore, rows that duplicate an existing row on a unique key value are
	// discarded.
	Ignore bool
	// FieldsOptionallyEnclosedBy set true if not all columns are enclosed.
	FieldsOptionallyEnclosedBy bool
	FieldsEnclosedBy           rune
	FieldsEscapedBy            rune
	LinesTerminatedBy          string
	FieldsTerminatedBy         string
	// LinesStartingBy: If all the lines you want to read in have a common
	// prefix that you want to ignore, you can use LINES STARTING BY
	// 'prefix_string' to skip over the prefix, and anything before it. If a
	// line does not include the prefix, the entire line is skipped.
	LinesStartingBy string
	// IgnoreLinesAtStart can be used to ignore lines at the start of the file.
	// For example, you can use IGNORE 1 LINES to skip over an initial header
	// line containing column names.
	IgnoreLinesAtStart int
	// Set must be a balanced key,value slice. The column list (field Columns)
	// can contain either column names or user variables. With user variables,
	// the SET clause enables you to perform transformations on their values
	// before assigning the result to columns. The SET clause can be used to
	// supply values not derived from the input file. e.g. SET column3 =
	// CURRENT_TIMESTAMP For more details please read
	// https://dev.mysql.com/doc/refman/5.7/en/load-data.html
	Set []string
	// Columns optional custom columns if the default columns of the table
	// differs from the CSV file. Column names do NOT get automatically quoted.
	Columns []string
	// Log optional logger for debugging purposes
	Log log.Logger
}

// LoadDataInfile loads a local CSV file into a MySQL table. For more details
// please read https://dev.mysql.com/doc/refman/5.7/en/load-data.html Files must
// be whitelisted by registering them with mysql.RegisterLocalFile(filepath)
// (recommended) or the Whitelist check must be deactivated by using the DSN
// parameter allowAllFiles=true (Might be insecure!). For more details
// https://godoc.org/github.com/go-sql-driver/mysql#RegisterLocalFile. To ignore
// foreign key constraints during the load operation, issue a SET
// foreign_key_checks = 0 statement before executing LOAD DATA.
func (t *Table) LoadDataInfile(ctx context.Context, filePath string, o InfileOptions) error {
	if t.IsView {
		return nil
	}
	if o.Log == nil {
		o.Log = log.BlackHole{}
	}

	var buf bytes.Buffer
	buf.WriteString("LOAD DATA ")
	if !o.IsNotLocal {
		buf.WriteString("LOCAL")
	}
	buf.WriteString(" INFILE '")
	buf.WriteString(filePath)
	buf.WriteRune('\'')
	switch {
	case o.Replace:
		buf.WriteString(" REPLACE ")
	case o.Ignore:
		buf.WriteString(" IGNORE ")
	}
	buf.WriteString(" INTO TABLE ")
	dml.Quoter.WriteQualifierName(&buf, t.Schema, t.Name)

	var hasFields bool
	if o.FieldsEscapedBy > 0 || o.FieldsTerminatedBy != "" || o.FieldsEnclosedBy > 0 {
		buf.WriteString(" FIELDS ")
		hasFields = true
	}
	if o.FieldsTerminatedBy != "" {
		buf.WriteString("TERMINATED BY '")
		buf.WriteString(o.FieldsTerminatedBy) // todo fix if it contains a single quote
		buf.WriteRune('\'')
	}
	if o.FieldsEnclosedBy > 0 {
		if o.FieldsOptionallyEnclosedBy {
			buf.WriteString(" OPTIONALLY ")
		}
		buf.WriteString(" ENCLOSED BY '")
		buf.WriteRune(o.FieldsEnclosedBy) // todo fix if it contains a single quote
		buf.WriteRune('\'')
	}
	if o.FieldsEscapedBy > 0 {
		buf.WriteString(" ESCAPED BY '")
		buf.WriteRune(o.FieldsEscapedBy) // todo fix if it contains a single quote
		buf.WriteRune('\'')
	}
	if hasFields {
		buf.WriteRune('\n')
	}

	var hasLines bool
	if o.LinesTerminatedBy != "" || o.LinesStartingBy != "" {
		buf.WriteString(" LINES ")
		hasLines = true
	}

	if o.LinesTerminatedBy != "" {
		buf.WriteString(" TERMINATED BY '")
		buf.WriteString(o.LinesTerminatedBy) // todo fix if it contains a single quote
		buf.WriteRune('\'')
	}
	if o.LinesStartingBy != "" {
		buf.WriteString(" STARTING BY '")
		buf.WriteString(o.LinesStartingBy) // todo fix if it contains a single quote
		buf.WriteRune('\'')
	}
	if hasLines {
		buf.WriteRune('\n')
	}

	if o.IgnoreLinesAtStart > 0 {
		fmt.Fprintf(&buf, "IGNORE %d LINES\n", o.IgnoreLinesAtStart)
	}

	// write COLUMNS
	buf.WriteString(" (")
	if len(o.Columns) == 0 {
		o.Columns = t.Columns.FieldNames()
	}
	for i, c := range o.Columns {
		if c != "" {
			buf.WriteString(c) // do not quote because custom columns or variables
		}
		if i < len(t.Columns)-1 {
			buf.WriteRune(',')
		}
	}
	buf.WriteString(")\n")

	if ls := len(o.Set); ls > 0 && ls%2 == 0 {
		buf.WriteString("SET ")
		for i := 0; i < ls; i += 2 {
			buf.WriteString(o.Set[i])
			buf.WriteRune('=')
			buf.WriteString(o.Set[i+1])
			if i+1 < ls-1 {
				buf.WriteRune(',')
				buf.WriteRune('\n')
			}
		}
	}
	buf.WriteRune(';')

	if o.Log.IsDebug() {
		o.Log.Debug("ddl.Table.Infile.SQL", log.String("sql", buf.String()))
	}
	return t.runExec(ctx, buf.String())
}
