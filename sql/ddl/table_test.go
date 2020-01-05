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
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/corestoreio/errors"
	"github.com/corestoreio/pkg/sql/ddl"
	"github.com/corestoreio/pkg/sql/dmltest"
	"github.com/corestoreio/pkg/storage/null"
	"github.com/corestoreio/pkg/util/assert"
)

var tableMap *ddl.Tables

func init() {
	tableMap = ddl.MustNewTables(
		ddl.WithTable(
			"catalog_category_anc_categs_index_idx",
			&ddl.Column{
				Field:      "category_id",
				ColumnType: "int(10) unsigned",
				Key:        "MUL",
				Default:    null.MakeString("0"),
				Extra:      "",
				Aliases:    []string{"entity_id"},
				Uniquified: true,
				StructTag:  `json:",omitempty"`,
			},
			&ddl.Column{
				Field:      "path",
				ColumnType: "varchar(255)",
				Null:       "YES",
				Key:        "MUL",
				Extra:      "",
			},
		),
		ddl.WithTable(
			"catalog_category_anc_categs_index_tmp",
			&ddl.Column{
				Field:      "category_id",
				ColumnType: "int(10) unsigned",
				Key:        "PRI",
				Default:    null.MakeString("0"),
				Extra:      "",
			},
			&ddl.Column{
				Field:      "path",
				ColumnType: "varchar(255)",
				Null:       "YES",
				Extra:      "",
			},
		),
	)

	if err := tableMap.Upsert(ddl.NewTable(
		"catalog_category_anc_products_index_idx",
		&ddl.Column{
			Field:      "category_id",
			ColumnType: "int(10) unsigned",
			Default:    null.MakeString("0"),
			Extra:      "",
		},
		&ddl.Column{
			Field:      "product_id",
			ColumnType: "int(10) unsigned",
			Key:        "",
			Default:    null.MakeString("0"),
			Extra:      "",
		},
		&ddl.Column{
			Field:      "position",
			ColumnType: "int(10) unsigned",
			Null:       "YES",
			Key:        "",
			Extra:      "",
		},
	),
	); err != nil {
		panic(err)
	}

	if err := tableMap.Upsert(ddl.NewTable(
		"admin_user",
		&ddl.Column{
			Field:      "user_id",
			ColumnType: "int(10) unsigned",
			Key:        "PRI",
			Extra:      "auto_increment",
		},
		&ddl.Column{
			Field:      "email",
			ColumnType: "varchar(128)",
			Null:       "YES",
			Key:        "",
			Extra:      "",
		},
		&ddl.Column{
			Field:      "first_name",
			ColumnType: "varchar(255)",
			Null:       "",
			Key:        "",
			Extra:      "",
		},
		&ddl.Column{
			Field:      "username",
			ColumnType: "varchar(40)",
			Null:       "YES",
			Key:        "UNI",
			Extra:      "",
		},
	),
	); err != nil {
		panic(err)
	}
}

func TestTable_HasColumn(t *testing.T) {
	t.Parallel()

	tbl := tableMap.MustTable("catalog_category_anc_categs_index_idx")
	assert.True(t, tbl.HasColumn("category_id"))
	assert.False(t, tbl.HasColumn("category_idx"))
	assert.False(t, tbl.HasColumn("category_ID"))
}

func TestTable_Select(t *testing.T) {
	t.Parallel()

	tbl := tableMap.MustTable("catalog_category_anc_categs_index_idx")
	sqlStr := tbl.Select("category_id", "path").String()
	assert.Exactly(t, "SELECT `category_id`, `path` FROM `catalog_category_anc_categs_index_idx` AS `main_table`", sqlStr)
}

func TestTableStructure(t *testing.T) {
	t.Parallel()

	sValid, err := tableMap.Table("catalog_category_anc_categs_index_idx")
	assert.NotNil(t, sValid)
	assert.NoError(t, err)

	assert.Equal(t, "catalog_category_anc_categs_index_tmp", tableMap.MustTable("catalog_category_anc_categs_index_tmp").Name)

	sInvalid, err := tableMap.Table("not_found")
	assert.Nil(t, sInvalid)
	assert.Error(t, err)
}

func TestTableStructureIn(t *testing.T) {
	t.Parallel()

	want := map[string]bool{
		"catalog_category_anc_categs_index_idx":   true,
		"catalog_category_anc_categs_index_tmp":   true,
		"catalog_category_anc_products_index_idx": false,
	}
	for tn, wantRes := range want {
		table := tableMap.MustTable(tn)
		assert.Exactly(t, wantRes, table.Columns.Contains("path"), "Table %s", table.Name)
	}

	want2 := map[string]bool{
		"catalog_category_anc_categs_index_idx":   true,
		"catalog_category_anc_categs_index_tmp":   true,
		"catalog_category_anc_products_index_idx": true,
	}
	for tn, wantRes := range want2 {
		table := tableMap.MustTable(tn)
		assert.Exactly(t, wantRes, table.Columns.Contains("category_id"), "Table %s", table.Name)
	}
}

func TestTable_Truncate(t *testing.T) {
	t.Parallel()

	t.Run("ok", func(t *testing.T) {
		dbc, dbMock := dmltest.MockDB(t)
		defer dmltest.MockClose(t, dbc, dbMock)

		dbMock.ExpectExec("TRUNCATE TABLE `catalog_category_anc_categs_index_tmp`").WillReturnResult(sqlmock.NewResult(0, 0))
		err := tableMap.MustTable("catalog_category_anc_categs_index_tmp").Truncate(context.TODO(), ddl.Options{Execer: dbc.DB})
		assert.NoError(t, err, "%+v", err)
	})
	t.Run("wait", func(t *testing.T) {
		dbc, dbMock := dmltest.MockDB(t)
		defer dmltest.MockClose(t, dbc, dbMock)

		dbMock.ExpectExec("TRUNCATE TABLE `catalog_category_anc_categs_index_tmp` WAIT 5 ").WillReturnResult(sqlmock.NewResult(0, 0))
		err := tableMap.MustTable("catalog_category_anc_categs_index_tmp").Truncate(context.TODO(), ddl.Options{Execer: dbc.DB, Wait: 5 * time.Second, Nowait: true})
		assert.NoError(t, err, "%+v", err)
	})
	t.Run("nowait", func(t *testing.T) {
		dbc, dbMock := dmltest.MockDB(t)
		defer dmltest.MockClose(t, dbc, dbMock)

		dbMock.ExpectExec("TRUNCATE TABLE `catalog_category_anc_categs_index_tmp` NOWAIT ").WillReturnResult(sqlmock.NewResult(0, 0))
		err := tableMap.MustTable("catalog_category_anc_categs_index_tmp").Truncate(context.TODO(), ddl.Options{Execer: dbc.DB, Nowait: true})
		assert.NoError(t, err, "%+v", err)
	})

	t.Run("Invalid table Name", func(t *testing.T) {
		tbl := ddl.NewTable("product")
		err := tbl.Rename(context.TODO(), "namecatalog_category_anc_categs_index_tmpcatalog_category_anc_categs_", ddl.Options{})
		assert.ErrorIsKind(t, errors.NotValid, err)
	})
}

func TestTable_Rename(t *testing.T) {
	t.Parallel()
	t.Run("ok", func(t *testing.T) {
		dbc, dbMock := dmltest.MockDB(t)
		defer dmltest.MockClose(t, dbc, dbMock)

		dbMock.ExpectExec("RENAME TABLE `catalog_category_anc_categs_index_tmp` TO `catalog_category_anc_categs`").
			WillReturnResult(sqlmock.NewResult(0, 0))
		err := tableMap.MustTable("catalog_category_anc_categs_index_tmp").Rename(
			context.TODO(), "catalog_category_anc_categs", ddl.Options{Execer: dbc.DB},
		)
		assert.NoError(t, err, "%+v", err)
	})
	t.Run("wait", func(t *testing.T) {
		dbc, dbMock := dmltest.MockDB(t)
		defer dmltest.MockClose(t, dbc, dbMock)

		dbMock.ExpectExec("RENAME TABLE `catalog_category_anc_categs_index_tmp`  WAIT 5  TO `catalog_category_anc_categs`").
			WillReturnResult(sqlmock.NewResult(0, 0))
		err := tableMap.MustTable("catalog_category_anc_categs_index_tmp").Rename(
			context.TODO(), "catalog_category_anc_categs", ddl.Options{Execer: dbc.DB, Wait: 5 * time.Second, Nowait: true},
		)
		assert.NoError(t, err, "%+v", err)
	})
	t.Run("nowait", func(t *testing.T) {
		dbc, dbMock := dmltest.MockDB(t)
		defer dmltest.MockClose(t, dbc, dbMock)

		dbMock.ExpectExec("RENAME TABLE `catalog_category_anc_categs_index_tmp`  NOWAIT  TO `catalog_category_anc_categs`").
			WillReturnResult(sqlmock.NewResult(0, 0))
		err := tableMap.MustTable("catalog_category_anc_categs_index_tmp").Rename(
			context.TODO(), "catalog_category_anc_categs", ddl.Options{Execer: dbc.DB, Nowait: true},
		)
		assert.NoError(t, err, "%+v", err)
	})

	t.Run("Invalid table Name", func(t *testing.T) {
		tbl := ddl.NewTable("product")
		err := tbl.Rename(context.TODO(), "namecatalog_category_anc_categs_index_tmpcatalog_category_anc_categs_", ddl.Options{})
		assert.ErrorIsKind(t, errors.NotValid, err)
	})
}

func TestTable_Swap(t *testing.T) {
	t.Parallel()

	t.Run("ok", func(t *testing.T) {
		dbc, dbMock := dmltest.MockDB(t)
		defer dmltest.MockClose(t, dbc, dbMock)

		dbMock.ExpectExec("RENAME TABLE `catalog_category_anc_categs_index_tmp` TO `catalog_category_anc_categs_index_tmp_[0-9]+`, `catalog_category_anc_categs_NEW` TO `catalog_category_anc_categs_index_tmp`,`catalog_category_anc_categs_index_tmp_[0-9]+` TO `catalog_category_anc_categs_NEW`").
			WillReturnResult(sqlmock.NewResult(0, 0))
		err := tableMap.MustTable("catalog_category_anc_categs_index_tmp").Swap(context.TODO(), "catalog_category_anc_categs_NEW", ddl.Options{Execer: dbc.DB})
		assert.NoError(t, err, "%+v", err)
	})

	t.Run("wait", func(t *testing.T) {
		dbc, dbMock := dmltest.MockDB(t)
		defer dmltest.MockClose(t, dbc, dbMock)

		dbMock.ExpectExec("RENAME TABLE `catalog_category_anc_categs_index_tmp`  WAIT 5  TO `catalog_category_anc_categs_index_tmp_[0-9]+`, `catalog_category_anc_categs_NEW` TO `catalog_category_anc_categs_index_tmp`,`catalog_category_anc_categs_index_tmp_[0-9]+` TO `catalog_category_anc_categs_NEW`").
			WillReturnResult(sqlmock.NewResult(0, 0))
		err := tableMap.MustTable("catalog_category_anc_categs_index_tmp").Swap(context.TODO(), "catalog_category_anc_categs_NEW", ddl.Options{Execer: dbc.DB, Wait: 5 * time.Second, Nowait: true})
		assert.NoError(t, err, "%+v", err)
	})

	t.Run("nowait", func(t *testing.T) {
		dbc, dbMock := dmltest.MockDB(t)
		defer dmltest.MockClose(t, dbc, dbMock)

		dbMock.ExpectExec("RENAME TABLE `catalog_category_anc_categs_index_tmp`  NOWAIT  TO `catalog_category_anc_categs_index_tmp_[0-9]+`, `catalog_category_anc_categs_NEW` TO `catalog_category_anc_categs_index_tmp`,`catalog_category_anc_categs_index_tmp_[0-9]+` TO `catalog_category_anc_categs_NEW`").
			WillReturnResult(sqlmock.NewResult(0, 0))
		err := tableMap.MustTable("catalog_category_anc_categs_index_tmp").Swap(context.TODO(), "catalog_category_anc_categs_NEW", ddl.Options{Execer: dbc.DB, Nowait: true})
		assert.NoError(t, err, "%+v", err)
	})

	t.Run("Invalid table Name", func(t *testing.T) {
		tbl := ddl.NewTable("product")
		err := tbl.Swap(context.TODO(), "namecatalog_category_anc_categs_index_tmpcatalog_category_anc_categs_", ddl.Options{})
		assert.ErrorIsKind(t, errors.NotValid, err)
	})
}

func TestTable_Drop(t *testing.T) {
	t.Parallel()
	t.Run("ok", func(t *testing.T) {
		dbc, dbMock := dmltest.MockDB(t)
		defer dmltest.MockClose(t, dbc, dbMock)

		dbMock.ExpectExec("DROP TABLE IF EXISTS `catalog_category_anc_categs_index_tmp`").
			WillReturnResult(sqlmock.NewResult(0, 0))
		err := tableMap.MustTable("catalog_category_anc_categs_index_tmp").Drop(context.TODO(),
			ddl.Options{Execer: dbc.DB},
		)
		assert.NoError(t, err, "%+v", err)
	})
	t.Run("wait", func(t *testing.T) {
		dbc, dbMock := dmltest.MockDB(t)
		defer dmltest.MockClose(t, dbc, dbMock)

		dbMock.ExpectExec("DROP TABLE IF EXISTS `catalog_category_anc_categs_index_tmp` WAIT 1 ").
			WillReturnResult(sqlmock.NewResult(0, 0))
		err := tableMap.MustTable("catalog_category_anc_categs_index_tmp").Drop(context.TODO(),
			ddl.Options{Execer: dbc.DB, Wait: time.Second, Nowait: true},
		)
		assert.NoError(t, err, "%+v", err)
	})
	t.Run("nowait", func(t *testing.T) {
		dbc, dbMock := dmltest.MockDB(t)
		defer dmltest.MockClose(t, dbc, dbMock)

		dbMock.ExpectExec("DROP TABLE IF EXISTS `catalog_category_anc_categs_index_tmp` NOWAIT ").
			WillReturnResult(sqlmock.NewResult(0, 0))
		err := tableMap.MustTable("catalog_category_anc_categs_index_tmp").Drop(context.TODO(),
			ddl.Options{Execer: dbc.DB, Nowait: true},
		)
		assert.NoError(t, err, "%+v", err)
	})
	t.Run("comment", func(t *testing.T) {
		dbc, dbMock := dmltest.MockDB(t)
		defer dmltest.MockClose(t, dbc, dbMock)

		dbMock.ExpectExec(dmltest.SQLMockQuoteMeta("DROP TABLE IF EXISTS /*X*/ `catalog_category_anc_categs_index_tmp` WAIT 2")).
			WillReturnResult(sqlmock.NewResult(0, 0))
		err := tableMap.MustTable("catalog_category_anc_categs_index_tmp").Drop(context.TODO(),
			ddl.Options{Execer: dbc.DB, Wait: 2 * time.Second, Comment: "X"},
		)
		assert.NoError(t, err, "%+v", err)
	})

	t.Run("Invalid table Name", func(t *testing.T) {
		tbl := ddl.NewTable("produ™€ct")
		err := tbl.Drop(context.TODO(), ddl.Options{})
		assert.ErrorIsKind(t, errors.NotValid, err)
	})
}

func TestTable_Optimize(t *testing.T) {
	t.Parallel()
	t.Run("ok", func(t *testing.T) {
		dbc, dbMock := dmltest.MockDB(t)
		defer dmltest.MockClose(t, dbc, dbMock)

		dbMock.ExpectExec("OPTIMIZE TABLE `catalog_category_anc_categs_index_tmp`").
			WillReturnResult(sqlmock.NewResult(0, 0))
		err := tableMap.MustTable("catalog_category_anc_categs_index_tmp").Optimize(context.TODO(),
			ddl.Options{Execer: dbc.DB},
		)
		assert.NoError(t, err, "%+v", err)
	})
	t.Run("wait", func(t *testing.T) {
		dbc, dbMock := dmltest.MockDB(t)
		defer dmltest.MockClose(t, dbc, dbMock)

		dbMock.ExpectExec("OPTIMIZE TABLE `catalog_category_anc_categs_index_tmp` WAIT 1 ").
			WillReturnResult(sqlmock.NewResult(0, 0))
		err := tableMap.MustTable("catalog_category_anc_categs_index_tmp").Optimize(context.TODO(),
			ddl.Options{Execer: dbc.DB, Wait: time.Second, Nowait: true},
		)
		assert.NoError(t, err, "%+v", err)
	})
	t.Run("nowait", func(t *testing.T) {
		dbc, dbMock := dmltest.MockDB(t)
		defer dmltest.MockClose(t, dbc, dbMock)

		dbMock.ExpectExec("OPTIMIZE TABLE `catalog_category_anc_categs_index_tmp` NOWAIT ").
			WillReturnResult(sqlmock.NewResult(0, 0))
		err := tableMap.MustTable("catalog_category_anc_categs_index_tmp").Optimize(context.TODO(),
			ddl.Options{Execer: dbc.DB, Nowait: true},
		)
		assert.NoError(t, err, "%+v", err)
	})

	t.Run("Invalid table Name", func(t *testing.T) {
		tbl := ddl.NewTable("produ™€ct")
		err := tbl.Optimize(context.TODO(), ddl.Options{})
		assert.ErrorIsKind(t, errors.NotValid, err)
	})
}

func TestTable_LoadDataInfile(t *testing.T) {
	t.Parallel()

	t.Run("default options", func(t *testing.T) {
		dbc, dbMock := dmltest.MockDB(t)
		defer dmltest.MockClose(t, dbc, dbMock)

		dbMock.ExpectExec(dmltest.SQLMockQuoteMeta("LOAD DATA LOCAL INFILE 'non-existent.csv' INTO TABLE `admin_user` (user_id,email,first_name,username) ;")).
			WillReturnResult(sqlmock.NewResult(0, 0))

		err := tableMap.MustTable("admin_user").LoadDataInfile(context.TODO(), "non-existent.csv", ddl.InfileOptions{Execer: dbc.DB})
		assert.NoError(t, err, "%+v", err)
	})

	t.Run("all options", func(t *testing.T) {
		dbc, dbMock := dmltest.MockDB(t)
		defer dmltest.MockClose(t, dbc, dbMock)

		dbMock.ExpectExec(dmltest.SQLMockQuoteMeta("LOAD DATA LOCAL INFILE 'non-existent.csv' REPLACE INTO TABLE `admin_user` FIELDS TERMINATED BY '|' OPTIONALLY ENCLOSED BY '+' ESCAPED BY '\"' LINES TERMINATED BY ' ' STARTING BY '###' IGNORE 1 LINES (user_id,@email,@username,) SET username=UPPER(@username), email=UPPER(@email);")).
			WillReturnResult(sqlmock.NewResult(0, 0))
		err := tableMap.MustTable("admin_user").LoadDataInfile(context.TODO(), "non-existent.csv", ddl.InfileOptions{
			Replace:                    true,
			FieldsTerminatedBy:         "|",
			FieldsOptionallyEnclosedBy: true,
			FieldsEnclosedBy:           '+',
			FieldsEscapedBy:            '"',
			LinesTerminatedBy:          "\r\n",
			LinesStartingBy:            "###",
			IgnoreLinesAtStart:         1,
			Columns:                    []string{"user_id", "@email", "@username"},
			Set:                        []string{"username", "UPPER(@username)", "email", "UPPER(@email)"},
			Execer:                     dbc.DB,
		})
		assert.NoError(t, err, "%+v", err)
	})
}

func TestTable_Artisan_Methods(t *testing.T) {
	dbc, dbMock := dmltest.MockDB(t)
	defer dmltest.MockClose(t, dbc, dbMock)
	assert.NoError(t, tableMap.Options(ddl.WithDB(dbc.DB)), "Cant set sql.DB on tableMap")

	tblAdmUser := tableMap.MustTable("admin_user")

	t.Run("Insert", func(t *testing.T) {
		dbMock.ExpectExec(dmltest.SQLMockQuoteMeta("INSERT INTO `admin_user` (`email`,`first_name`,`username`) VALUES (?,?,?),(?,?,?)")).
			WithArgs("a@b.c", "Franz", "franz", "d@e.f", "Sissi", "sissi").
			WillReturnResult(sqlmock.NewResult(11, 0))

		res, err := tblAdmUser.Insert().WithDBR().ExecContext(context.Background(),
			"a@b.c", "Franz", "franz",
			"d@e.f", "Sissi", "sissi",
		)
		assert.NoError(t, err)
		id, err := res.LastInsertId()
		assert.NoError(t, err)
		assert.Exactly(t, int64(11), id)
	})

	t.Run("DeleteByPK prepared", func(t *testing.T) {
		dbMock.ExpectExec(dmltest.SQLMockQuoteMeta("DELETE FROM `admin_user` WHERE (`user_id` IN (?,?))")).
			WithArgs("a@b.c", "d@e.f").
			WillReturnResult(sqlmock.NewResult(0, 2))

		res, err := tblAdmUser.DeleteByPK().WithDBR().ExpandPlaceHolders().
			ExecContext(context.Background(), []string{"a@b.c", "d@e.f"})
		assert.NoError(t, err)
		id, err := res.RowsAffected()
		assert.NoError(t, err)
		assert.Exactly(t, int64(2), id)
	})
	t.Run("DeleteByPK interpolated", func(t *testing.T) {
		dbMock.ExpectExec(dmltest.SQLMockQuoteMeta("DELETE FROM `admin_user` WHERE (`user_id` IN ('a@b.c','d@e.f'))")).
			WithArgs().
			WillReturnResult(sqlmock.NewResult(0, 2))

		res, err := tblAdmUser.DeleteByPK().WithDBR().Interpolate().
			ExecContext(context.Background(), []string{"a@b.c", "d@e.f"})
		assert.NoError(t, err, "%+v", err)
		id, err := res.RowsAffected()
		assert.NoError(t, err)
		assert.Exactly(t, int64(2), id)
	})

	t.Run("SelectByPK", func(t *testing.T) {
		dbMock.ExpectQuery(dmltest.SQLMockQuoteMeta("SELECT `user_id`, `email`, `first_name`, `username` FROM `admin_user` AS `main_table` WHERE (`user_id` IN (?,?,?))")).
			WithArgs(int64(234), int64(235), int64(236)).
			WillReturnRows(sqlmock.NewRows([]string{"user_id", "email", "first_name", "username"}))

		rows, err := tblAdmUser.SelectByPK("*").WithDBR().ExpandPlaceHolders().QueryContext(context.Background(), []int64{234, 235, 236})
		assert.NoError(t, err)
		assert.NoError(t, rows.Close())
	})

	t.Run("SelectAll", func(t *testing.T) {
		dbMock.ExpectQuery(dmltest.SQLMockQuoteMeta("SELECT `user_id`, `email`, `first_name`, `username` FROM `admin_user` AS `main_table`")).
			WithArgs().
			WillReturnRows(sqlmock.NewRows([]string{"user_id", "email", "first_name", "username"}))

		rows, err := tblAdmUser.Select("*").WithDBR().QueryContext(context.Background())
		assert.NoError(t, err)
		assert.NoError(t, rows.Close())
	})

	t.Run("UpdateByPK", func(t *testing.T) {
		dbMock.ExpectExec(dmltest.SQLMockQuoteMeta("UPDATE `admin_user` SET `email`=?, `first_name`=?, `username`=? WHERE (`user_id` = ?)")).
			WithArgs("a@b.c", "Franz", "franz", int64(3)).
			WillReturnResult(sqlmock.NewResult(0, 1))

		res, err := tblAdmUser.UpdateByPK().WithDBR().ExecContext(context.Background(),
			"a@b.c", "Franz", "franz", 3,
		)
		assert.NoError(t, err)
		id, err := res.RowsAffected()
		assert.NoError(t, err)
		assert.Exactly(t, int64(1), id)
	})
}

func TestTable_GeneratedColumns(t *testing.T) {
	dbc := dmltest.MustConnectDB(t)
	defer dmltest.Close(t, dbc)
	defer dmltest.SQLDumpLoad(t, "testdata/generated*.sql", nil).Deferred()

	tbls := ddl.MustNewTables(
		ddl.WithDB(dbc.DB),
		ddl.WithCreateTable(context.TODO(), "core_config_data_generated", ""),
	)

	// INSERT should contain only the non-generated columns.
	ins := tbls.MustTable("core_config_data_generated").Insert().BuildValues()
	assert.Exactly(t, "INSERT INTO `core_config_data_generated` (`type_id`,`expires`,`path`,`value`) VALUES (?,?,?,?)", ins.String())
}
