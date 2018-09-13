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
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/corestoreio/errors"
	"github.com/corestoreio/pkg/sql/ddl"
	"github.com/corestoreio/pkg/sql/dml"
	"github.com/corestoreio/pkg/sql/dmltest"
	"github.com/corestoreio/pkg/storage/null"
	"github.com/corestoreio/pkg/util/assert"
)

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
		ddl.WithCreateTable(context.TODO(), nil, ""),
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

	ts := ddl.MustNewTables(ddl.WithCreateTable(context.TODO(), nil, "a3", "", "b5", "", "c7", ""))
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

	ts := ddl.MustNewTables(ddl.WithCreateTable(context.TODO(), nil, "a3", "", "b5", "", "c7", ""))
	ts.DeleteAllFromCache()
	assert.Exactly(t, 0, ts.Len())
}

func TestTables_Upsert_Update(t *testing.T) {
	t.Parallel()

	ts := ddl.MustNewTables(ddl.WithCreateTable(context.TODO(), nil, "a3", "", "b5", "", "c7", ""))
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
			assert.True(t, errors.NotValid.Match(err), "%+v", err)
		} else {
			t.Error("Expecting a panic")
		}
	}()

	ts := ddl.MustNewTables(ddl.WithCreateTable(context.TODO(), nil, "a3"))
	tbl := ts.MustTable("a3")
	assert.NotNil(t, tbl)
	tbl = ts.MustTable("a44")
	assert.Nil(t, tbl)
}

func TestWithTableNames(t *testing.T) {
	t.Parallel()

	ts := ddl.MustNewTables(ddl.WithCreateTable(context.TODO(), nil, "a3", "", "b5", "", "c7", ""))
	t.Run("Ok", func(t *testing.T) {
		assert.Exactly(t, "a3", ts.MustTable("a3").Name)
		assert.Exactly(t, "b5", ts.MustTable("b5").Name)
		assert.Exactly(t, "c7", ts.MustTable("c7").Name)
	})

	t.Run("Invalid Identifier", func(t *testing.T) {
		err := ts.Options(ddl.WithCreateTable(context.TODO(), nil, "x1", ""))
		assert.True(t, errors.NotValid.Match(err), "%+v", err)
		assert.Contains(t, err.Error(), `identifier "x\uf8ff1" (Case 2)`)
	})
}

func TestWithCreateTable_Mock_DoesNotCreateTable(t *testing.T) {
	t.Parallel()

	dbc, dbMock := dmltest.MockDB(t)
	defer dmltest.MockClose(t, dbc, dbMock)

	rows := sqlmock.NewRows([]string{"Table", "Create Table"}).
		FromCSVString(
			`"admin_user","CREATE TABLE admin_user ( user_id int(10),  PRIMARY KEY (user_id))"
`)
	dbMock.ExpectQuery("SHOW CREATE TABLE `admin_user`").
		WillReturnRows(rows)

	rows = sqlmock.NewRows([]string{"TABLE_NAME", "COLUMN_NAME", "ORDINAL_POSITION", "COLUMN_DEFAULT", "IS_NULLABLE", "DATA_TYPE", "CHARACTER_MAXIMUM_LENGTH", "NUMERIC_PRECISION", "NUMERIC_SCALE", "COLUMN_TYPE", "COLUMN_KEY", "EXTRA", "COLUMN_COMMENT"}).
		FromCSVString(
			`"admin_user","user_id",1,0,"NO","int",0,10,0,"int(10) unsigned","PRI","auto_increment","User ID"
"admin_user","firsname",2,NULL,"YES","varchar",32,0,0,"varchar(32)","","","User First Name"
"admin_user","modified",8,"CURRENT_TIMESTAMP","NO","timestamp",0,0,0,"timestamp","","on update CURRENT_TIMESTAMP","User Modified Time"
`)

	dbMock.ExpectQuery("SELECT.+FROM information_schema.COLUMNS WHERE TABLE_SCHEMA=DATABASE\\(\\) AND TABLE_NAME.+").
		WillReturnRows(rows)

	tm0, err := ddl.NewTables(
		ddl.WithCreateTable(context.TODO(), dbc.DB, "admin_user", "CREATE BUGGY TABLE"),
	)
	assert.NoError(t, err, "%+v", err)

	table := tm0.MustTable("admin_user")
	assert.Exactly(t, []string{"user_id", "firsname", "modified"}, table.Columns.FieldNames())
	assert.Exactly(t, "CREATE TABLE admin_user ( user_id int(10),  PRIMARY KEY (user_id))", table.CreateSyntax)
	//t.Log(table.Columns.GoString())
}

func TestWithCreateTable_Mock_DoesCreateTable(t *testing.T) {
	t.Parallel()

	dbc, dbMock := dmltest.MockDB(t)
	defer dmltest.MockClose(t, dbc, dbMock)

	dbMock.ExpectExec(dmltest.SQLMockQuoteMeta("CREATE TABLE `admin_user` ( user_id int(10),  PRIMARY KEY (user_id))")).
		WillReturnResult(sqlmock.NewResult(0, 0))

	dbMock.ExpectQuery("SHOW CREATE TABLE `admin_user`").
		WillReturnRows(sqlmock.NewRows([]string{"Table", "Create Table"}).
			FromCSVString(`"admin_user","CREATE TABLE admin_user ( user_id int(10),  PRIMARY KEY (user_id))"` + "\n"))

	dbMock.ExpectQuery("SELECT.+FROM information_schema.COLUMNS WHERE TABLE_SCHEMA=DATABASE\\(\\) AND TABLE_NAME.+").
		WillReturnRows(sqlmock.NewRows([]string{"TABLE_NAME", "COLUMN_NAME", "ORDINAL_POSITION", "COLUMN_DEFAULT", "IS_NULLABLE", "DATA_TYPE", "CHARACTER_MAXIMUM_LENGTH", "NUMERIC_PRECISION", "NUMERIC_SCALE", "COLUMN_TYPE", "COLUMN_KEY", "EXTRA", "COLUMN_COMMENT"}).
			FromCSVString(`"admin_user","user_id",1,0,"NO","int",0,10,0,"int(10) unsigned","PRI","auto_increment","User ID"
"admin_user","firsname",2,NULL,"YES","varchar",32,0,0,"varchar(32)","","","User First Name"
"admin_user","modified",8,"CURRENT_TIMESTAMP","NO","timestamp",0,0,0,"timestamp","","on update CURRENT_TIMESTAMP","User Modified Time"
`))

	tm0, err := ddl.NewTables(
		ddl.WithCreateTable(context.TODO(), dbc.DB, "admin_user", "CREATE TABLE `admin_user` ( user_id int(10),  PRIMARY KEY (user_id))"),
	)
	assert.NoError(t, err, "%+v", err)

	table := tm0.MustTable("admin_user")
	assert.Exactly(t, []string{"user_id", "firsname", "modified"}, table.Columns.FieldNames())
	assert.Exactly(t, "CREATE TABLE admin_user ( user_id int(10),  PRIMARY KEY (user_id))", table.CreateSyntax)
	//t.Log(table.Columns.GoString())
}

func TestWithDropTable(t *testing.T) {
	t.Parallel()

	t.Run("not previously registered", func(t *testing.T) {

		t.Run("with DISABLE_FOREIGN_KEY_CHECKS", func(t *testing.T) {
			dbc, dbMock := dmltest.MockDB(t)
			defer dmltest.MockClose(t, dbc, dbMock)

			dbMock.ExpectExec(`SET FOREIGN_KEY_CHECKS=0`).WillReturnResult(sqlmock.NewResult(0, 0))
			dbMock.ExpectExec("DROP TABLE IF EXISTS `admin_user`").WillReturnResult(sqlmock.NewResult(0, 0))
			dbMock.ExpectExec("DROP TABLE IF EXISTS `admin_role`").WillReturnResult(sqlmock.NewResult(0, 0))
			dbMock.ExpectExec("DROP VIEW IF EXISTS `view_admin_perm`").WillReturnResult(sqlmock.NewResult(0, 0))
			dbMock.ExpectExec(`SET FOREIGN_KEY_CHECKS=1`).WillReturnResult(sqlmock.NewResult(0, 0))

			tm0, err := ddl.NewTables(
				ddl.WithDropTable(context.TODO(), dbc.DB, "DISABLE_FOREIGN_KEY_CHECKS", "admin_user", "admin_role", "view_admin_perm"),
			)
			assert.NoError(t, err, "%+v", err)
			_ = tm0
		})

		t.Run("with DISABLE_FOREIGN_KEY_CHECKS, invalid identifier", func(t *testing.T) {
			dbc, dbMock := dmltest.MockDB(t)
			defer dmltest.MockClose(t, dbc, dbMock)

			dbMock.ExpectExec(`SET FOREIGN_KEY_CHECKS=0`).WillReturnResult(sqlmock.NewResult(0, 0))

			dbMock.ExpectExec(`SET FOREIGN_KEY_CHECKS=1`).WillReturnResult(sqlmock.NewResult(0, 0))

			tm0, err := ddl.NewTables(
				ddl.WithDropTable(context.TODO(), dbc.DB, "DISABLE_FOREIGN_KEY_CHECKS", "admin_user"),
			)
			assert.Nil(t, tm0)
			assert.True(t, errors.NotValid.Match(err), "%+v", err)
		})

		t.Run("without DISABLE_FOREIGN_KEY_CHECKS", func(t *testing.T) {
			dbc, dbMock := dmltest.MockDB(t)
			defer dmltest.MockClose(t, dbc, dbMock)

			dbMock.ExpectExec("DROP TABLE IF EXISTS `admin_user`").WillReturnResult(sqlmock.NewResult(0, 0))
			dbMock.ExpectExec("DROP TABLE IF EXISTS `admin_role`").WillReturnResult(sqlmock.NewResult(0, 0))

			tm0, err := ddl.NewTables(
				ddl.WithDropTable(context.TODO(), dbc.DB, "", "admin_user", "admin_role"),
			)
			assert.NoError(t, err, "%+v", err)
			_ = tm0
		})
	})

	t.Run("previously registered at Tables", func(t *testing.T) {

		t.Run("with DISABLE_FOREIGN_KEY_CHECKS", func(t *testing.T) {
			dbc, dbMock := dmltest.MockDB(t)
			defer dmltest.MockClose(t, dbc, dbMock)

			dbMock.ExpectExec(`SET FOREIGN_KEY_CHECKS=0`).WillReturnResult(sqlmock.NewResult(0, 0))
			dbMock.ExpectExec("DROP TABLE IF EXISTS `admin_user`").WillReturnResult(sqlmock.NewResult(0, 0))
			dbMock.ExpectExec("DROP TABLE IF EXISTS `admin_role`").WillReturnResult(sqlmock.NewResult(0, 0))
			dbMock.ExpectExec(`SET FOREIGN_KEY_CHECKS=1`).WillReturnResult(sqlmock.NewResult(0, 0))

			tm0, err := ddl.NewTables(
				ddl.WithTable("admin_user"),
				ddl.WithTable("admin_role"),
				ddl.WithDropTable(context.TODO(), dbc.DB, "DISABLE_FOREIGN_KEY_CHECKS", "admin_user", "admin_role"),
			)

			assert.NoError(t, err, "%+v", err)
			_ = tm0
		})
	})
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
		assert.Exactly(t, "TeschtT", tbl.Name)
	})

	t.Run("Nil Table and after WithTable call", func(*testing.T) {
		ts := ddl.MustNewTables(
			ddl.WithTableDMLListeners("TeschtU", ev, ev),
			ddl.WithTable("TeschtU", &ddl.Column{Field: "col1"}),
		) // +=2
		tbl := ts.MustTable("TeschtU")
		assert.Exactly(t, "TeschtU", tbl.Name)
	})
}

func TestWithCreateTable_FromQuery(t *testing.T) {
	t.Parallel()

	t.Run("create table fails", func(t *testing.T) {
		dbc, dbMock := dmltest.MockDB(t)
		defer dmltest.MockClose(t, dbc, dbMock)

		xErr := errors.AlreadyClosed.Newf("Connection already closed")
		dbMock.ExpectExec(dmltest.SQLMockQuoteMeta("CREATE TABLE `testTable` AS SELECT * FROM catalog_product_entity")).WillReturnError(xErr)

		tbls, err := ddl.NewTables(ddl.WithCreateTable(context.TODO(), dbc.DB, "testTable", "CREATE TABLE `testTable` AS SELECT * FROM catalog_product_entity"))
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

		dbMock.ExpectQuery("SHOW CREATE TABLE `testTable`").
			WillReturnRows(sqlmock.NewRows([]string{"Table", "Create Table"}).
				FromCSVString(`"testTable","CREATE TABLE testTable ( user_id int(10),  PRIMARY KEY (user_id))"` + "\n"))

		dbMock.ExpectQuery("SELEC.+\\s+FROM\\s+information_schema\\.COLUMNS").WillReturnError(xErr)

		tbls, err := ddl.NewTables(ddl.WithCreateTable(context.TODO(), dbc.DB, "testTable", "CREATE TABLE `testTable` AS SELECT * FROM catalog_product_entity"))
		assert.Nil(t, tbls)
		assert.True(t, errors.AlreadyClosed.Match(err), "%+v", err)
	})
}

func TestWithTableDB(t *testing.T) {
	t.Parallel()
	dbc, dbMock := dmltest.MockDB(t)
	defer dmltest.MockClose(t, dbc, dbMock)

	ts := ddl.MustNewTables(
		ddl.WithDB(dbc.DB),
		ddl.WithTable("tableA"),
		ddl.WithCreateTable(context.TODO(), nil, "tableB", ""),
	) // +=2

	assert.Exactly(t, dbc.DB, ts.MustTable("tableA").DB)
	assert.Exactly(t, dbc.DB, ts.MustTable("tableB").DB)
	assert.Exactly(t, []string{"tableA", "tableB"}, ts.Tables())
}

func newCCD() *ddl.Tables {
	return ddl.MustNewTables(
		ddl.WithTable(
			"core_config_data",
			&ddl.Column{Field: `config_id`, ColumnType: `int(10) unsigned`, Null: `NO`, Key: `PRI`, Extra: `auto_increment`},
			&ddl.Column{Field: `scope`, ColumnType: `varchar(8)`, Null: `NO`, Key: `MUL`, Default: null.MakeString(`default`), Extra: ""},
			&ddl.Column{Field: `scope_id`, ColumnType: `int(11)`, Null: `NO`, Key: "", Default: null.MakeString(`0`), Extra: ""},
			&ddl.Column{Field: `path`, ColumnType: `varchar(255)`, Null: `NO`, Key: "", Default: null.MakeString(`general`), Extra: ""},
			&ddl.Column{Field: `value`, ColumnType: `text`, Null: `YES`, Key: ``, Extra: ""},
		),
	)
}

func TestTables_Validate(t *testing.T) {
	dbc, dbMock := dmltest.MockDB(t)
	defer dmltest.MockClose(t, dbc, dbMock)
	tbls := newCCD()
	tbls.DB = dbc.DB

	t.Run("context timeout", func(t *testing.T) {
		ctx, cancel := context.WithCancel(context.Background())
		cancel()
		err := tbls.Validate(ctx)
		err = errors.Cause(err)
		assert.EqualError(t, err, "context canceled")
	})
	t.Run("Validation OK", func(t *testing.T) {
		dbMock.ExpectQuery("SELECT.+FROM information_schema.COLUMNS WHERE").
			WillReturnRows(
				dmltest.MustMockRows(dmltest.WithFile("testdata/core_config_data_columns.csv")))
		err := tbls.Validate(context.Background())
		assert.NoError(t, err, "Validation should succeed")
	})
	t.Run("mismatch field name", func(t *testing.T) {
		tbls.MustTable("core_config_data").Columns[0].Field = "configID"
		dbMock.ExpectQuery("SELECT.+FROM information_schema.COLUMNS WHERE").
			WillReturnRows(
				dmltest.MustMockRows(dmltest.WithFile("testdata/core_config_data_columns.csv")))
		err := tbls.Validate(context.Background())

		assert.True(t, errors.Mismatch.Match(err), "should have kind mismatch")
		assert.EqualError(t, err, "[ddl] Table \"core_config_data\" with column name \"configID\" at index 0 does not match database column name \"config_id\"")
		tbls.MustTable("core_config_data").Columns[0].Field = "config_id"
	})
	t.Run("mismatch column type", func(t *testing.T) {
		tbls.MustTable("core_config_data").Columns[0].ColumnType = "varchar(XX)"
		dbMock.ExpectQuery("SELECT.+FROM information_schema.COLUMNS WHERE").
			WillReturnRows(
				dmltest.MustMockRows(dmltest.WithFile("testdata/core_config_data_columns.csv")))
		err := tbls.Validate(context.Background())

		assert.True(t, errors.Mismatch.Match(err), "should have kind mismatch")
		assert.EqualError(t, err, "[ddl] Table \"core_config_data\" with Go column name \"config_id\" does not match MySQL column type. MySQL: \"int(10) unsigned\" Go: \"varchar(XX)\".")
		tbls.MustTable("core_config_data").Columns[0].ColumnType = "int(10) unsigned"
	})
	t.Run("mismatch null property", func(t *testing.T) {
		tbls.MustTable("core_config_data").Columns[0].Null = "YES"
		dbMock.ExpectQuery("SELECT.+FROM information_schema.COLUMNS WHERE").
			WillReturnRows(
				dmltest.MustMockRows(dmltest.WithFile("testdata/core_config_data_columns.csv")))
		err := tbls.Validate(context.Background())

		assert.True(t, errors.Mismatch.Match(err), "should have kind mismatch")
		assert.EqualError(t, err, "[ddl] Table \"core_config_data\" with column name \"config_id\" does not match MySQL null types. MySQL: \"NO\" Go: \"YES\"")
		tbls.MustTable("core_config_data").Columns[0].Null = "NO"
	})

	t.Run("too many tables", func(t *testing.T) {
		tbls := ddl.MustNewTables(
			ddl.WithTable("core_config_data",
				&ddl.Column{Field: `config_id`, ColumnType: `int(10) unsigned`, Null: `NO`, Key: `PRI`, Extra: `auto_increment`},
			),
			ddl.WithTable("customer_entity",
				&ddl.Column{Field: `config_id`, ColumnType: `int(10) unsigned`, Null: `NO`, Key: `PRI`, Extra: `auto_increment`},
			),
			ddl.WithDB(dbc.DB),
		)

		dbMock.ExpectQuery("SELECT.+FROM information_schema.COLUMNS WHERE").
			WillReturnRows(
				dmltest.MustMockRows(dmltest.WithFile("testdata/core_config_data_columns.csv")))
		err := tbls.Validate(context.Background())

		assert.True(t, errors.Mismatch.Match(err), "should have kind mismatch")
		assert.EqualError(t, err, "[ddl] Tables count 2 does not match table count 1 in database.")
	})

	t.Run("less columns", func(t *testing.T) {

		dbMock.ExpectQuery("SELECT.+FROM information_schema.COLUMNS WHERE").
			WillReturnRows(
				dmltest.MustMockRows(dmltest.WithFile("testdata/core_config_data_columns_less.csv")))
		err := tbls.Validate(context.Background())

		assert.True(t, errors.Mismatch.Match(err), "should have kind mismatch")
		assert.EqualError(t, err, "[ddl] Table \"core_config_data\" has more columns (count 5) than its object (column count 4) in the database.")
	})

	t.Run("more columns", func(t *testing.T) {

		dbMock.ExpectQuery("SELECT.+FROM information_schema.COLUMNS WHERE").
			WillReturnRows(
				dmltest.MustMockRows(dmltest.WithFile("testdata/core_config_data_columns_more.csv")))
		err := tbls.Validate(context.Background())
		assert.NoError(t, err)
	})

}
