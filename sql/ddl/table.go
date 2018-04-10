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
	"time"

	"github.com/corestoreio/errors"
	"github.com/corestoreio/log"
	"github.com/corestoreio/pkg/sql/dml"
	"github.com/corestoreio/pkg/util/bufferpool"
)

// Table represents a table from a specific database.
type Table struct {
	DB dml.QueryExecPreparer
	// Schema represents the name of the database. Might be empty.
	Schema string
	// Name of the table
	Name string
	// Columns all table columns
	Columns Columns
	// Listeners specific pre defined listeners which gets dispatches to each
	// DML statement (SELECT, INSERT, UPDATE or DELETE).
	Listeners dml.ListenerBucket
	// IsView set to true to mark if the table is a view
	IsView       bool
	columnsPK    []string
	columnsNonPK []string
	columnsAll   []string
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
	if len(t.Columns) == 0 {
		return t
	}
	t.columnsNonPK = t.columnsNonPK[:0]
	t.columnsNonPK = t.Columns.NonPrimaryColumns().FieldNames(t.columnsNonPK...)

	t.columnsPK = t.columnsPK[:0]
	t.columnsPK = t.Columns.PrimaryKeys().FieldNames(t.columnsPK...)

	t.columnsAll = t.columnsAll[:0]
	t.columnsAll = t.Columns.FieldNames(t.columnsAll...)

	return t
}

// Insert creates a new INSERT statement with all non primary key columns. If
// OnDuplicateKey() gets called, the INSERT can be used as an update or create
// statement. Adding multiple VALUES section is allowed.
func (t *Table) Insert() *dml.Insert {
	i := dml.NewInsert(t.Name).AddColumns(t.columnsNonPK...)
	i.RecordPlaceHolderCount = len(i.Columns)
	i.Listeners = i.Listeners.Merge(t.Listeners.Insert)
	return i.WithDB(t.DB)
}

// SelectAll creates a new `SELECT column1,column2,cloumnX FROM table` without a
// WHERE clause.
func (t *Table) SelectAll() *dml.Select {
	s := dml.NewSelect(t.columnsAll...).
		FromAlias(t.Name, MainTable)
	s.Listeners = s.Listeners.Merge(t.Listeners.Select)
	return s.WithDB(t.DB)
}

// Select creates a new SELECT statement with a set FROM clause.
func (t *Table) Select(columns ...string) *dml.Select {
	s := dml.NewSelect(columns...).
		FromAlias(t.Name, MainTable)
	s.Listeners = s.Listeners.Merge(t.Listeners.Select)
	return s.WithDB(t.DB)
}

// SelectByPK creates a new `SELECT * FROM table WHERE id IN (?)`
func (t *Table) SelectByPK() *dml.Select {
	s := dml.NewSelect(t.columnsAll...).FromAlias(t.Name, MainTable)
	s.Wheres = t.whereByPK(dml.In)
	s.Listeners = s.Listeners.Merge(t.Listeners.Select)
	return s.WithDB(t.DB)
}

// DeleteByPK creates a new `DELETE FROM table WHERE id IN (?)`
func (t *Table) DeleteByPK() *dml.Delete {
	d := dml.NewDelete(t.Name)
	d.Wheres = t.whereByPK(dml.In)
	d.Listeners = d.Listeners.Merge(t.Listeners.Delete)
	return d.WithDB(t.DB)
}

// UpdateByPK creates a new `UPDATE table SET ... WHERE id = ?`. The SET clause
// contains all non primary columns.
func (t *Table) UpdateByPK() *dml.Update {
	u := dml.NewUpdate(t.Name).AddColumns(t.columnsNonPK...)
	u.Wheres = t.whereByPK(dml.Equal)
	u.Listeners = u.Listeners.Merge(t.Listeners.Update)
	return u.WithDB(t.DB)
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

func (t *Table) resetColumns() {
	if cap(t.Columns) == 0 {
		t.Columns = make(Columns, 0, 10)
	}
	for i := range t.Columns {
		// Pointer must be nilled to remove a reference and avoid a memory
		// leak, AFAIK.
		t.Columns[i] = nil
	}
	t.Columns = t.Columns[:0]
}

// MapColumns implements dml.ColumnMapper interface to read column values from a
// query with table information_schema.COLUMNS.
func (t *Table) MapColumns(rc *dml.ColumnMap) error {
	if rc.Count == 0 {
		t.resetColumns()
	}

	c, tableName, err := NewColumn(rc)
	if err != nil {
		return errors.Wrapf(err, "[ddl] Table.RowScan. Table %q\n", t.Name)
	}

	if t.Name == "" {
		t.Name = tableName
	}

	t.Columns = append(t.Columns, c)
	t.update()
	return nil
}

// ToSQL creates a SQL query for loading all columns for the current table.
func (t *Table) ToSQL() (string, []interface{}, error) {
	sqlStr, _, err := dml.Interpolate(selTablesColumns).Str(t.Name).ToSQL()
	if err != nil {
		return "", nil, errors.Wrapf(err, "[ddl] Table.ToSQL.Interpolate for table %q", t.Name)
	}
	return sqlStr, nil, nil
}

// Truncate truncates the tables. Removes all rows and sets the auto increment
// to zero. Just like a CREATE TABLE statement.
func (t *Table) Truncate(ctx context.Context, execer dml.Execer) error {
	if t.IsView {
		return nil
	}
	if err := dml.IsValidIdentifier(t.Name); err != nil {
		return errors.WithStack(err)
	}
	ddl := "TRUNCATE TABLE " + dml.Quoter.QualifierName(t.Schema, t.Name)
	_, err := execer.ExecContext(ctx, ddl)
	return errors.Wrapf(err, "[ddl] failed to truncate table %q", ddl)
}

// Rename renames the current table to the new table name. Renaming is an atomic
// operation in the database. As long as two databases are on the same file
// system, you can use RENAME TABLE to move a table from one database to
// another. RENAME TABLE also works for views, as long as you do not try to
// rename a view into a different database.
func (t *Table) Rename(ctx context.Context, execer dml.Execer, new string) error {
	if err := dml.IsValidIdentifier(t.Name); err != nil {
		return errors.WithStack(err)
	}
	if err := dml.IsValidIdentifier(new); err != nil {
		return errors.WithStack(err)
	}
	ddl := "RENAME TABLE " + dml.Quoter.QualifierName(t.Schema, t.Name) + " TO " + dml.Quoter.NameAlias(new, "")
	_, err := execer.ExecContext(ctx, ddl)
	return errors.Wrapf(err, "[ddl] failed to rename table %q", ddl)
}

// Swap swaps the current table with the other table of the same structure.
// Renaming is an atomic operation in the database. Note: indexes won't get
// swapped! As long as two databases are on the same file system, you can use
// RENAME TABLE to move a table from one database to another.
func (t *Table) Swap(ctx context.Context, execer dml.Execer, other string) error {
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

	if _, err := execer.ExecContext(ctx, buf.String()); err != nil {
		// only allocs in case of an error ;-)
		return errors.Wrapf(err, "[ddl] Failed to swap table %q", buf.String())
	}
	return nil
}

// Drop drops, if exists, the table or the view.
func (t *Table) Drop(ctx context.Context, execer dml.Execer) error {
	typ := "TABLE"
	if t.IsView {
		typ = "VIEW"
	}
	if err := dml.IsValidIdentifier(t.Name); err != nil {
		return errors.Wrap(err, "[ddl] Drop table name")
	}
	_, err := execer.ExecContext(ctx, "DROP "+typ+" IF EXISTS "+dml.Quoter.QualifierName(t.Schema, t.Name))
	return errors.Wrapf(err, "[ddl] failed to drop table %q", t.Name)
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
func (t *Table) LoadDataInfile(ctx context.Context, db dml.Execer, filePath string, o InfileOptions) error {
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
		for i := 0; i < ls; i = i + 2 {
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

	_, err := db.ExecContext(ctx, buf.String())
	return errors.Fatal.New(err, "[csb] Infile for table %q failed with query: %q", t.Name, buf.String())
}
