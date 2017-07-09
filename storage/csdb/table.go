// Copyright 2015-2016, Cyrill @ Schumacher.fm and the CoreStore contributors
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

package csdb

import (
	"bytes"
	"context"
	"database/sql"
	"fmt"
	"strconv"
	"time"

	"github.com/corestoreio/csfw/storage/dbr"
	"github.com/corestoreio/csfw/util/bufferpool"
	"github.com/corestoreio/errors"
	"github.com/corestoreio/log"
)

// Table represents a table from a specific database.
type Table struct {
	Convert dbr.RowConvert
	// Schema represents the name of the database. Might be empty.
	Schema string
	// Name of the table
	Name string
	// Columns all table columns
	Columns Columns
	// Listeners specific pre defined listeners which gets dispatches to each
	// DML statement (SELECT, INSERT, UPDATE or DELETE).
	Listeners dbr.ListenerBucket
	// IsView set to true to mark if the table is a view
	IsView bool

	// selectAllCache no quite sure about this one .... maybe remove it
	selectAllCache *dbr.Select
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

	t.selectAllCache = &dbr.Select{
	// Columns: t.AllColumnAliasQuote(MainTable), // TODO refactor
	//Table: dbr.MakeNameAlias(t.Name, MainTable),
	}

	return t
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

// RowScan implements dbr.Scanner interface
func (t *Table) RowScan(r *sql.Rows) error {
	if t.Convert.Count == 0 {
		t.resetColumns()
	}
	if err := t.Convert.Scan(r); err != nil {
		return err
	}

	c, tableName, err := NewColumn(&t.Convert)
	if err != nil {
		return errors.Wrapf(err, "[csdb] Table.RowScan. Table %q Columns %v\n", t.Name, t.Convert.Columns)
	}

	if t.Name == "" {
		t.Name = tableName
	}

	t.Columns = append(t.Columns, c)
	return nil
}

func (t *Table) ScanClose() error {
	t.update()
	return nil
}

func (t *Table) ToSQL() (string, []interface{}, error) {
	sql, err := dbr.Interpolate(selTablesColumns, dbr.ArgString(t.Name))
	if err != nil {
		return "", nil, errors.Wrapf(err, "[csdb] Table.ToSQL.Interpolate for table %q", t.Name)
	}
	return sql, nil, nil
}

// Truncate truncates the tables. Removes all rows and sets the auto increment
// to zero. Just like a CREATE TABLE statement.
func (t *Table) Truncate(ctx context.Context, execer dbr.Execer) error {
	if t.IsView {
		return nil
	}
	if err := IsValidIdentifier(t.Name); err != nil {
		return errors.Wrap(err, "[csdb] Truncate table name")
	}
	ddl := "TRUNCATE TABLE " + dbr.Quoter.QualifierName(t.Schema, t.Name)
	_, err := execer.ExecContext(ctx, ddl)
	return errors.Wrapf(err, "[csdb] failed to truncate table %q", ddl)
}

// Rename renames the current table to the new table name. Renaming is an atomic
// operation in the database. As long as two databases are on the same file
// system, you can use RENAME TABLE to move a table from one database to
// another. RENAME TABLE also works for views, as long as you do not try to
// rename a view into a different database.
func (t *Table) Rename(ctx context.Context, execer dbr.Execer, new string) error {
	if err := IsValidIdentifier(t.Name, new); err != nil {
		return errors.Wrap(err, "[csdb] Rename table name")
	}
	ddl := "RENAME TABLE " + dbr.Quoter.QualifierName(t.Schema, t.Name) + " TO " + dbr.Quoter.NameAlias(new, "")
	_, err := execer.ExecContext(ctx, ddl)
	return errors.Wrapf(err, "[csdb] failed to rename table %q", ddl)
}

// Swap swaps the current table with the other table of the same structure.
// Renaming is an atomic operation in the database. Note: indexes won't get
// swapped! As long as two databases are on the same file system, you can use
// RENAME TABLE to move a table from one database to another.
func (t *Table) Swap(ctx context.Context, execer dbr.Execer, other string) error {
	if err := IsValidIdentifier(t.Name, other); err != nil {
		return errors.Wrap(err, "[csdb] Swap table name")
	}

	tmp := TableName("", t.Name, strconv.FormatInt(time.Now().UnixNano(), 10))

	buf := bufferpool.Get()
	defer bufferpool.Put(buf)
	buf.WriteString("RENAME TABLE ")
	dbr.Quoter.WriteQualifierName(buf, t.Schema, t.Name)
	buf.WriteString(" TO ")
	dbr.Quoter.WriteNameAlias(buf, tmp, "")
	buf.WriteString(", ")
	dbr.Quoter.WriteNameAlias(buf, other, "")
	buf.WriteString(" TO ")
	dbr.Quoter.WriteQualifierName(buf, t.Schema, t.Name)
	buf.WriteByte(',')
	dbr.Quoter.WriteNameAlias(buf, tmp, "")
	buf.WriteString(" TO ")
	dbr.Quoter.WriteNameAlias(buf, other, "")

	if _, err := execer.ExecContext(ctx, buf.String()); err != nil {
		// only allocs in case of an error ;-)
		return errors.Wrapf(err, "[csdb] Failed to swap table %q", buf.String())
	}
	return nil
}

// Drop drops, if exists, the table or the view.
func (t *Table) Drop(ctx context.Context, execer dbr.Execer) error {
	typ := "TABLE"
	if t.IsView {
		typ = "VIEW"
	}
	if err := IsValidIdentifier(t.Name); err != nil {
		return errors.Wrap(err, "[csdb] Drop table name")
	}
	_, err := execer.ExecContext(ctx, "DROP "+typ+" IF EXISTS "+dbr.Quoter.QualifierName(t.Schema, t.Name))
	return errors.Wrapf(err, "[csdb] failed to drop table %q", t.Name)
}

// Load performs a SELECT * FROM `tableName` query and puts the results
// into the pointer slice `dest`. Returns the number of loaded rows and nil or 0
// and an error. The variadic third arguments can modify the SQL query.
func (t *Table) Load(ctx context.Context, db dbr.Querier, dest interface{}, listeners ...dbr.Listen) (int, error) {
	//sb := t.Select()
	//sb.DB.Querier = db
	//sb.Listeners.Merge(t.Listeners.Select)
	//sb.Listeners.Add(listeners...)
	//return sb.LoadStructs(ctx, dest)
	return 0, nil
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
	Ignore             bool
	FieldsTerminatedBy string
	// FieldsOptionallyEnclosedBy set true if not all columns are enclosed.
	FieldsOptionallyEnclosedBy bool
	FieldsEnclosedBy           rune
	FieldsEscapedBy            rune
	LinesTerminatedBy          string
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
func (t *Table) LoadDataInfile(ctx context.Context, db dbr.Execer, filePath string, o InfileOptions) error {
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
	dbr.Quoter.WriteQualifierName(&buf, t.Schema, t.Name)

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
		o.Log.Debug("csdb.Table.Infile.SQL", log.String("sql", buf.String()))
	}

	_, err := db.ExecContext(ctx, buf.String())
	return errors.NewFatal(err, "[csb] Infile for table %q failed with query: %q", t.Name, buf.String())
}
