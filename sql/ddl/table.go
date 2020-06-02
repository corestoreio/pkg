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
	"github.com/corestoreio/pkg/storage/null"
	"github.com/corestoreio/pkg/util/bufferpool"
)

// Table represents a table from a specific database with a bound default
// connection pool. Its fields are not secure to use in concurrent context and
// hence might cause a race condition if used not properly.
type Table struct {
	// dcp represents a Database Connection Pool inherited from the Tables type.
	// Shall not be nil.
	dcp *dml.ConnPool

	// Always def.
	Catalog string
	// Database name.
	Schema string
	// Table name.
	Name string
	// One of BASE TABLE for a regular table, VIEW for a view, SYSTEM VIEW for
	// Information Schema tables or SYSTEM VERSIONED for system-versioned tables.
	Type string
	// Storage Engine.
	Engine null.String
	// Version number from the table's .frm file
	Version null.Uint64
	// Row format (see InnoDB, Aria and MyISAM row formats).
	RowFormat null.String
	// Number of rows in the table. Some engines, such as XtraDB and InnoDB may
	// store an estimate.
	TableRows null.Uint64
	// Average row length in the table.
	AvgRowLength null.Uint64
	// For InnoDB/XtraDB, the index size, in pages, multiplied by the page size.
	// For Aria and MyISAM, length of the data file, in bytes. For MEMORY, the
	// approximate allocated memory.
	DataLength null.Uint64
	// Maximum length of the data file, ie the total number of bytes that could be
	// stored in the table. Not used in XtraDB and InnoDB.
	MaxDataLength null.Uint64
	// Length of the index file.
	IndexLength null.Uint64
	// Bytes allocated but unused. For InnoDB tables in a shared tablespace, the
	// free space of the shared tablespace with small safety margin. An estimate
	// in the case of partitioned tables - see the PARTITIONS table.
	DataFree null.Uint64
	// Next AUTO_INCREMENT value.
	AutoIncrement null.Uint64
	// Time the table was created.
	CreateTime null.Time
	// Time the table was last updated. On Windows, the timestamp is not updated on
	// update, so MyISAM values will be inaccurate. In InnoDB, if shared
	// tablespaces are used, will be NULL, while buffering can also delay the
	// update, so the value will differ from the actual time of the last UPDATE,
	// INSERT or DELETE.
	UpdateTime null.Time
	// Time the table was last checked. Not kept by all storage engines, in which
	// case will be NULL.
	CheckTime null.Time
	// Character set and collation.
	TableCollation null.String
	// Live checksum value, if any.
	Checksum null.Uint64
	// Extra CREATE TABLE options.
	CreateOptions null.String
	// Table comment provided when MariaDB created the table.
	TableComment   string
	MaxIndexLength null.Uint64
	// Columns all table columns. They do not get used to create or alter a
	// table.
	Columns Columns
	// optimized column selection for specific DML operations.
	columnsPK    []string // only primary key columns
	columnsNonPK []string // all columns, except PK and system-versioned
	columnsAll   []string // all columns, except system-versioned
	// columnsUpsert contains all non-current-timestamp, non-virtual, non-system
	// versioned and non auto_increment columns for update or insert operations.
	columnsUpsert []string
	// colset is a set to check case-sensitively if a table has a column.
	colset map[string]struct{}
}

// NewTable initializes a new table structure with minimal information and
// without a database connection.
func NewTable(tableName string, cs ...*Column) *Table {
	ts := &Table{
		Name:    tableName,
		Columns: Columns(cs),
	}
	return ts.update()
}

func newTable(rc *dml.ColumnMap) (*Table, error) {
	t := new(Table)
	for rc.Next() {
		switch col := rc.Column(); col {
		case "TABLE_CATALOG":
			rc.String(&t.Catalog)
		case "TABLE_SCHEMA":
			rc.String(&t.Schema)
		case "TABLE_NAME":
			rc.String(&t.Name)
		case "TABLE_TYPE":
			rc.String(&t.Type)
		case "ENGINE":
			rc.NullString(&t.Engine)
		case "VERSION":
			rc.NullUint64(&t.Version)
		case "ROW_FORMAT":
			rc.NullString(&t.RowFormat)
		case "TABLE_ROWS":
			rc.NullUint64(&t.TableRows)
		case "AVG_ROW_LENGTH":
			rc.NullUint64(&t.AvgRowLength)
		case "DATA_LENGTH":
			rc.NullUint64(&t.DataLength)
		case "MAX_DATA_LENGTH":
			rc.NullUint64(&t.MaxDataLength)
		case "INDEX_LENGTH":
			rc.NullUint64(&t.IndexLength)
		case "DATA_FREE":
			rc.NullUint64(&t.DataFree)
		case "AUTO_INCREMENT":
			rc.NullUint64(&t.AutoIncrement)
		case "CREATE_TIME":
			rc.NullTime(&t.CreateTime)
		case "UPDATE_TIME":
			rc.NullTime(&t.UpdateTime)
		case "CHECK_TIME":
			rc.NullTime(&t.CheckTime)
		case "TABLE_COLLATION":
			rc.NullString(&t.TableCollation)
		case "CHECKSUM":
			rc.NullUint64(&t.Checksum)
		case "CREATE_OPTIONS":
			rc.NullString(&t.CreateOptions)
		case "TABLE_COMMENT":
			rc.String(&t.TableComment)
		case "MAX_INDEX_LENGTH":
			rc.NullUint64(&t.MaxIndexLength)
		default:
			return nil, errors.NotSupported.Newf("[ddl] Column %q not supported", col)
		}
	}
	return t, errors.WithStack(rc.Err())
}

// IsView determines if a table is a view. Either via system attribute or via
// its table name.
func (t *Table) IsView() bool {
	return t.Type == "VIEW" || t.Type == "SYSTEM VIEW" ||
		strings.HasPrefix(t.Name, PrefixView) || strings.HasSuffix(t.Name, SuffixView)
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

// Insert creates a new INSERT statement with all non primary key columns. If
// OnDuplicateKey() gets called, the INSERT can be used as an update or create
// statement. Adding multiple VALUES section is allowed. Using this statement to
// prepare a query, a call to `BuildValues()` triggers building the VALUES
// clause, otherwise a SQL parse error will occur.
func (t *Table) Insert() *dml.Insert {
	return dml.NewInsert(t.Name).AddColumns(t.columnsUpsert...)
}

// Select creates a new SELECT statement. If "*" gets set as an argument, then
// all columns will be added to to list of columns.
func (t *Table) Select(columns ...string) *dml.Select {
	if len(columns) == 1 && columns[0] == "*" {
		columns = t.columnsAll
	}
	return dml.NewSelect(columns...).FromAlias(t.Name, MainTable)
}

// SelectByPK creates a new `SELECT columns FROM table WHERE id = ?`. If "*"
// gets set as an argument, then all columns will be added to to list of
// columns.
func (t *Table) SelectByPK(columns ...string) *dml.Select {
	if len(columns) == 1 && columns[0] == "*" {
		columns = t.columnsAll
	}
	s := dml.NewSelect(columns...).FromAlias(t.Name, MainTable)
	s.Wheres = t.WhereByPK(dml.Equal)
	return s
}

// DeleteByPK creates a new `DELETE FROM table WHERE id = ?`
func (t *Table) DeleteByPK() *dml.Delete {
	d := dml.NewDelete(t.Name)
	d.Wheres = t.WhereByPK(dml.Equal)
	return d
}

// Delete creates a new `DELETE FROM table` statement.
func (t *Table) Delete() *dml.Delete {
	return dml.NewDelete(t.Name)
}

// UpdateByPK creates a new `UPDATE table SET ... WHERE id = ?`. The SET clause
// contains all non primary columns.
func (t *Table) UpdateByPK() *dml.Update {
	u := dml.NewUpdate(t.Name).AddColumns(t.columnsUpsert...)
	u.Wheres = t.WhereByPK(dml.Equal)
	return u
}

// Update creates a new UPDATE statement without a WHERE clause.
func (t *Table) Update() *dml.Update {
	return dml.NewUpdate(t.Name).AddColumns(t.columnsUpsert...)
}

// WhereByPK puts the primary keys as WHERE clauses into a condition.
func (t *Table) WhereByPK(op dml.Op) dml.Conditions {
	cnds := make(dml.Conditions, 0, 1)
	for _, pk := range t.columnsPK {
		c := dml.Column(pk).PlaceHolder()
		c.Operator = op
		cnds = append(cnds, c)
	}
	return cnds
}

func (t *Table) runExec(ctx context.Context, o Options, qry string) error {
	var te dml.Execer
	if t.dcp != nil {
		te = t.dcp.DB
	}
	if _, err := o.exec(te).ExecContext(ctx, qry); err != nil {
		return errors.Wrapf(err, "[ddl] failed to exec %q", qry) // please do change this return signature, saves an alloc
	}
	return nil
}

// Truncate truncates the table. Removes all rows and sets the auto increment to
// zero. Just like a CREATE TABLE statement. To use a custom connection, call
// WithDB before.
func (t *Table) Truncate(ctx context.Context, o Options) error {
	if t.IsView() {
		return nil
	}
	if err := dml.IsValidIdentifier(t.Name); err != nil {
		return errors.WithStack(err)
	}
	var buf strings.Builder
	buf.WriteString("TRUNCATE TABLE ")
	buf.WriteString(dml.Quoter.QualifierName(t.Schema, t.Name))
	o.sqlAddShouldWait(&buf)
	return t.runExec(ctx, o, buf.String())
}

// Rename renames the current table to the new table name. Renaming is an atomic
// operation in the database. As long as two databases are on the same file
// system, you can use RENAME TABLE to move a table from one database to
// another. RENAME TABLE also works for views, as long as you do not try to
// rename a view into a different database. To use a custom connection.
// https://mariadb.com/kb/en/rename-table/
func (t *Table) Rename(ctx context.Context, newTableName string, o Options) error {
	if err := dml.IsValidIdentifier(t.Name); err != nil {
		return errors.WithStack(err)
	}
	if err := dml.IsValidIdentifier(newTableName); err != nil {
		return errors.WithStack(err)
	}
	var buf strings.Builder
	buf.WriteString("RENAME TABLE ")
	buf.WriteString(dml.Quoter.QualifierName(t.Schema, t.Name))
	o.sqlAddShouldWait(&buf)
	buf.WriteString(" TO ")
	buf.WriteString(dml.Quoter.NameAlias(newTableName, ""))
	return t.runExec(ctx, o, buf.String())
}

// Swap swaps the current table with the other table of the same structure.
// Renaming is an atomic operation in the database. Note: indexes won't get
// swapped! As long as two databases are on the same file system, you can use
// RENAME TABLE to move a table from one database to another. To use a custom
// connection, call WithDB before.
// RENAME TABLE has to wait for existing queries on
// the table to finish until it can be executed. That would be fine, but it also
// locks out other queries while waiting for RENAME to happen! This can cause a
// serious locking up of your database tables.
// https://mariadb.com/kb/en/rename-table/
func (t *Table) Swap(ctx context.Context, other string, o Options) error {
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
	o.sqlAddShouldWait(buf)
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
	return t.runExec(ctx, o, buf.String())
}

func (t *Table) getTyp() string {
	if t.IsView() {
		return "VIEW"
	}
	return "TABLE"
}

// Drop drops, if exists, the table or the view. To use a custom connection,
// call WithDB before.
func (t *Table) Drop(ctx context.Context, o Options) error {
	if err := dml.IsValidIdentifier(t.Name); err != nil {
		return errors.Wrap(err, "[ddl] Drop table name")
	}
	var buf strings.Builder
	buf.WriteString("DROP ")
	buf.WriteString(t.getTyp())
	buf.WriteString(" IF EXISTS ")
	if o.Comment != "" {
		buf.WriteString("/*")
		buf.WriteString(o.Comment)
		buf.WriteString("*/ ")
	}
	buf.WriteString(dml.Quoter.QualifierName(t.Schema, t.Name))
	o.sqlAddShouldWait(&buf)
	return t.runExec(ctx, o, buf.String())
}

// Optimize optimizes a table. https://mariadb.com/kb/en/optimize-table/
func (t *Table) Optimize(ctx context.Context, o Options) error {
	if err := dml.IsValidIdentifier(t.Name); err != nil {
		return errors.Wrap(err, "[ddl] Optimize table name")
	}
	var buf strings.Builder
	buf.WriteString("OPTIMIZE TABLE ")
	buf.WriteString(dml.Quoter.QualifierName(t.Schema, t.Name))
	o.sqlAddShouldWait(&buf)
	return t.runExec(ctx, o, buf.String())
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
	Log    log.Logger
	Execer dml.Execer
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
	if t.IsView() {
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
	return t.runExec(ctx, Options{Execer: o.Execer}, buf.String())
}
