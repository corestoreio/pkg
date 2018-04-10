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

package ddl_test

import (
	"context"
	"io"
	"testing"

	sqlmock "github.com/DATA-DOG/go-sqlmock"
	"github.com/corestoreio/errors"
	"github.com/corestoreio/pkg/sql/ddl"
	"github.com/corestoreio/pkg/sql/dml"
	"github.com/corestoreio/pkg/sql/dmltest"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var _ dml.QueryBuilder = (*ddl.Tables)(nil)
var _ dml.ColumnMapper = (*ddl.Tables)(nil)
var _ io.Closer = (*ddl.Tables)(nil)

func TestNewTableService(t *testing.T) {
	t.Parallel()
	assert.Equal(t, ddl.MustNewTables().Len(), 0)

	tm1 := ddl.MustNewTables(
		ddl.WithTable("store"),
		ddl.WithTable("store_group"),
		ddl.WithTable("store_website"),
	)
	assert.Equal(t, 3, tm1.Len())
}

func TestNewTableServicePanic(t *testing.T) {
	t.Parallel()

	defer func() {
		if r := recover(); r != nil {
			err := r.(error)
			assert.True(t, errors.NotValid.Match(err), "%+v", err)
		} else {
			t.Error("Expecting a panic")
		}
	}()

	_ = ddl.MustNewTables(
		ddl.WithTable(""),
	)
}

func TestTables_Upsert_Insert(t *testing.T) {
	t.Parallel()

	ts := ddl.MustNewTables()

	t.Run("Insert OK", func(t *testing.T) {
		assert.NoError(t, ts.Upsert(ddl.NewTable("test1")))
		assert.Equal(t, 1, ts.Len())
	})
}

func TestTables_DeleteFromCache(t *testing.T) {
	t.Parallel()

	ts := ddl.MustNewTables(ddl.WithTableNames("a3", "b5", "c7"))
	t.Run("Delete One", func(t *testing.T) {
		ts.DeleteFromCache("b5")
		assert.Exactly(t, 2, ts.Len())
	})
	t.Run("Delete All does nothing", func(t *testing.T) {
		ts.DeleteFromCache()
		assert.Exactly(t, 2, ts.Len())
	})
}

func TestTables_DeleteAllFromCache(t *testing.T) {
	t.Parallel()

	ts := ddl.MustNewTables(ddl.WithTableNames("a3", "b5", "c7"))
	ts.DeleteAllFromCache()
	assert.Exactly(t, 0, ts.Len())
}

func TestTables_Upsert_Update(t *testing.T) {
	t.Parallel()

	ts := ddl.MustNewTables(ddl.WithTableNames("a3", "b5", "c7"))
	t.Run("One", func(t *testing.T) {
		ts.Upsert(ddl.NewTable("x5"))
		assert.Exactly(t, 4, ts.Len())
		tb, err := ts.Table("x5")
		assert.NoError(t, err, "%+v", err)
		assert.Exactly(t, `x5`, tb.Name)
	})
}

func TestTables_MustTable(t *testing.T) {
	t.Parallel()

	defer func() {
		if r := recover(); r != nil {
			err := r.(error)
			assert.True(t, errors.NotFound.Match(err), "%+v", err)
		} else {
			t.Error("Expecting a panic")
		}
	}()

	ts := ddl.MustNewTables(ddl.WithTableNames("a3"))
	tbl := ts.MustTable("a3")
	assert.NotNil(t, tbl)
	tbl = ts.MustTable("a44")
	assert.Nil(t, tbl)
}

func TestWithTableNames(t *testing.T) {
	t.Parallel()

	ts := ddl.MustNewTables(ddl.WithTableNames("a3", "b5", "c7"))
	t.Run("Ok", func(t *testing.T) {
		assert.Exactly(t, "a3", ts.MustTable("a3").Name)
		assert.Exactly(t, "b5", ts.MustTable("b5").Name)
		assert.Exactly(t, "c7", ts.MustTable("c7").Name)
	})

	t.Run("Invalid Identifier", func(t *testing.T) {
		err := ts.Options(ddl.WithTableNames("x1"))
		assert.True(t, errors.NotValid.Match(err), "%+v", err)
		assert.Contains(t, err.Error(), `identifier "x\uf8ff1" (Case 2)`)
	})
}

func TestTables_RowScan_Integration(t *testing.T) {
	t.Parallel()

	dbc := dmltest.MustConnectDB(t)
	defer dmltest.Close(t, dbc)

	tm0 := ddl.MustNewTables(
		ddl.WithTable("admin_user"),
	)
	_, err := dbc.WithQueryBuilder(tm0).Load(context.TODO(), tm0)
	require.NoError(t, err)

	table := tm0.MustTable("admin_user")

	assert.True(t, len(table.Columns.FieldNames()) >= 15)
	// t.Logf("%+v", table.Columns)
}

func TestTables_RowScan_Mock(t *testing.T) {
	t.Parallel()

	dbc, dbMock := dmltest.MockDB(t)
	defer dmltest.MockClose(t, dbc, dbMock)

	rows := sqlmock.NewRows([]string{"TABLE_NAME", "COLUMN_NAME", "ORDINAL_POSITION", "COLUMN_DEFAULT", "IS_NULLABLE", "DATA_TYPE", "CHARACTER_MAXIMUM_LENGTH", "NUMERIC_PRECISION", "NUMERIC_SCALE", "COLUMN_TYPE", "COLUMN_KEY", "EXTRA", "COLUMN_COMMENT"}).
		FromCSVString(
			`"admin_user","user_id",1,0,"NO","int",0,10,0,"int(10) unsigned","PRI","auto_increment","User ID"
"admin_user","firsname",2,NULL,"YES","varchar",32,0,0,"varchar(32)","","","User First Name"
"admin_user","modified",8,"CURRENT_TIMESTAMP","NO","timestamp",0,0,0,"timestamp","","on update CURRENT_TIMESTAMP","User Modified Time"
`)

	dbMock.ExpectQuery("SELECT.+FROM information_schema.COLUMNS WHERE TABLE_SCHEMA=DATABASE\\(\\) AND TABLE_NAME.+").
		WillReturnRows(rows)

	tm0 := ddl.MustNewTables(
		ddl.WithTable("admin_user"),
	)
	_, err := dbc.WithQueryBuilder(tm0).Load(context.TODO(), tm0)
	require.NoError(t, err)

	table := tm0.MustTable("admin_user")
	assert.Exactly(t, []string{"user_id", "firsname", "modified"}, table.Columns.FieldNames())
	//t.Log(table.Columns.GoString())
}

func TestWithTableDMLListeners(t *testing.T) {
	t.Parallel()

	counter := 0
	ev := dml.MustNewListenerBucket(
		dml.Listen{
			Name:           "l1",
			EventType:      dml.OnBeforeToSQL,
			ListenSelectFn: func(_ *dml.Select) { counter++ },
		},
		dml.Listen{
			Name:           "l2",
			EventType:      dml.OnBeforeToSQL,
			ListenSelectFn: func(_ *dml.Select) { counter++ },
		},
	)

	t.Run("Nil Table / No-WithTable", func(*testing.T) {
		ts := ddl.MustNewTables(
			ddl.WithTableDMLListeners("tableA", ev, ev),
			ddl.WithTable("tableA"),
		) // +=2
		tbl := ts.MustTable("tableA")
		sel := dml.NewSelect().From("tableA")
		sel.Listeners.Merge(tbl.Listeners.Select) // +=2
		sel.AddColumns("a", "b")
		assert.Exactly(t, "SELECT `a`, `b` FROM `tableA`", sel.String())
		assert.Exactly(t, 4, counter) // yes 4 is correct
	})

	t.Run("Non Nil Table", func(*testing.T) {
		ts := ddl.MustNewTables(
			ddl.WithTable("TeschtT", &ddl.Column{Field: "col1"}),
			ddl.WithTableDMLListeners("TeschtT", ev, ev),
		) // +=2
		tbl := ts.MustTable("TeschtT")
		require.Exactly(t, "TeschtT", tbl.Name)
	})

	t.Run("Nil Table and after WithTable call", func(*testing.T) {
		ts := ddl.MustNewTables(
			ddl.WithTableDMLListeners("TeschtU", ev, ev),
			ddl.WithTable("TeschtU", &ddl.Column{Field: "col1"}),
		) // +=2
		tbl := ts.MustTable("TeschtU")
		require.Exactly(t, "TeschtU", tbl.Name)
	})
}

func TestWithTableLoadColumns(t *testing.T) {
	t.Parallel()

	t.Run("Invalid Identifier", func(t *testing.T) {
		tbls, err := ddl.NewTables(ddl.WithTableLoadColumns(context.TODO(), nil, "H€llo"))
		assert.Nil(t, tbls)
		assert.True(t, errors.NotValid.Match(err), "%+v", err)
	})

	t.Run("Ok", func(t *testing.T) {

		dbc, dbMock := dmltest.MockDB(t)
		defer dmltest.MockClose(t, dbc, dbMock)

		rows := sqlmock.NewRows([]string{"TABLE_NAME", "COLUMN_NAME", "ORDINAL_POSITION", "COLUMN_DEFAULT", "IS_NULLABLE", "DATA_TYPE", "CHARACTER_MAXIMUM_LENGTH", "NUMERIC_PRECISION", "NUMERIC_SCALE", "COLUMN_TYPE", "COLUMN_KEY", "EXTRA", "COLUMN_COMMENT"}).
			FromCSVString(
				`"admin_user","user_id",1,0,"NO","int",0,10,0,"int(10) unsigned","PRI","auto_increment","User ID"
"admin_user","firstname",2,NULL,"YES","varchar",32,0,0,"varchar(32)","","","User First Name"
"admin_user","modified",8,"CURRENT_TIMESTAMP","NO","timestamp",0,0,0,"timestamp","","on update CURRENT_TIMESTAMP","User Modified Time"
`)

		dbMock.ExpectQuery("SELECT.+FROM information_schema.COLUMNS WHERE TABLE_SCHEMA=DATABASE\\(\\) AND TABLE_NAME.+").
			WillReturnRows(rows)

		tm0 := ddl.MustNewTables(
			ddl.WithTableLoadColumns(context.TODO(), dbc.DB, "admin_user"),
		)

		table := tm0.MustTable("admin_user")
		assert.Exactly(t, []string{"user_id", "firstname", "modified"},
			table.Columns.FieldNames())
	})
}

func TestWithTableOrViewFromQuery(t *testing.T) {
	t.Parallel()

	t.Run("Invalid type", func(t *testing.T) {
		tbls, err := ddl.NewTables(ddl.WithTableOrViewFromQuery(context.TODO(), nil, "proc", "asdasd", "SELECT * from"))
		assert.Nil(t, tbls)
		assert.True(t, errors.Unavailable.Match(err), "%+v", err)
	})

	t.Run("Invalid object name", func(t *testing.T) {
		tbls, err := ddl.NewTables(ddl.WithTableOrViewFromQuery(context.TODO(), nil, "proc", "asdasd", "SELECT * from"))
		assert.Nil(t, tbls)
		assert.True(t, errors.NotValid.Match(err), "%+v", err)
	})

	t.Run("drop table fails", func(t *testing.T) {
		dbc, dbMock := dmltest.MockDB(t)
		defer dmltest.MockClose(t, dbc, dbMock)

		xErr := errors.AlreadyClosed.Newf("Connection already closed")
		dbMock.ExpectExec("DROP TABLE IF EXISTS `testTable`").WillReturnError(xErr)

		tbls, err := ddl.NewTables(ddl.WithTableOrViewFromQuery(context.TODO(), dbc.DB, "table", "testTable", "SELECT * FROM catalog_product_entity", true))
		assert.Nil(t, tbls)
		assert.True(t, errors.AlreadyClosed.Match(err), "%+v", err)
	})

	t.Run("create table fails", func(t *testing.T) {
		dbc, dbMock := dmltest.MockDB(t)
		defer dmltest.MockClose(t, dbc, dbMock)

		xErr := errors.AlreadyClosed.Newf("Connection already closed")
		dbMock.ExpectExec(dmltest.SQLMockQuoteMeta("CREATE TABLE `testTable` AS SELECT * FROM catalog_product_entity")).WillReturnError(xErr)

		tbls, err := ddl.NewTables(ddl.WithTableOrViewFromQuery(context.TODO(), dbc.DB, "table", "testTable", "SELECT * FROM catalog_product_entity", false))
		assert.Nil(t, tbls)
		assert.True(t, errors.AlreadyClosed.Match(err), "%+v", err)
	})

	t.Run("load columns fails", func(t *testing.T) {
		dbc, dbMock := dmltest.MockDB(t)
		defer dmltest.MockClose(t, dbc, dbMock)

		xErr := errors.AlreadyClosed.Newf("Connection already closed")
		dbMock.
			ExpectExec(dmltest.SQLMockQuoteMeta("CREATE TABLE `testTable` AS SELECT * FROM catalog_product_entity")).
			WillReturnResult(sqlmock.NewResult(0, 0))

		dbMock.ExpectQuery("SELEC.+\\s+FROM\\s+information_schema\\.COLUMNS").WillReturnError(xErr)

		tbls, err := ddl.NewTables(ddl.WithTableOrViewFromQuery(context.TODO(), dbc.DB, "table", "testTable", "SELECT * FROM catalog_product_entity", false))
		assert.Nil(t, tbls)
		assert.True(t, errors.AlreadyClosed.Match(err), "%+v", err)
	})

	t.Run("create view", func(t *testing.T) {
		dbc, dbMock := dmltest.MockDB(t)
		defer dmltest.MockClose(t, dbc, dbMock)

		dbMock.
			ExpectExec(dmltest.SQLMockQuoteMeta("CREATE VIEW `testTable` AS SELECT * FROM core_config_data")).
			WillReturnResult(sqlmock.NewResult(0, 0))

		dbMock.ExpectQuery("SELECT.+FROM information_schema.COLUMNS WHERE").
			WillReturnRows(
				dmltest.MustMockRows(dmltest.WithFile("testdata/core_config_data_columns.csv")))

		tbls, err := ddl.NewTables(ddl.WithTableOrViewFromQuery(context.TODO(), dbc.DB, "view", "testTable", "SELECT * FROM core_config_data", false))
		if err != nil {
			t.Fatalf("%+v", err)
		}
		assert.Exactly(t, "testTable", tbls.MustTable("testTable").Name)
		assert.True(t, tbls.MustTable("testTable").IsView, "Table should be a view")
	})
}

func TestWithTableDB(t *testing.T) {
	t.Parallel()
	dbc, dbMock := dmltest.MockDB(t)
	defer dmltest.MockClose(t, dbc, dbMock)

	ts := ddl.MustNewTables(
		ddl.WithDB(dbc.DB),
		ddl.WithTable("tableA"),
		ddl.WithTable("tableB"),
	) // +=2

	assert.Exactly(t, dbc.DB, ts.MustTable("tableA").DB)
	assert.Exactly(t, dbc.DB, ts.MustTable("tableB").DB)
}
