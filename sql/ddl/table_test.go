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
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var _ dml.QueryBuilder = (*ddl.Table)(nil)
var _ dml.ColumnMapper = (*ddl.Table)(nil)

var tableMap *ddl.Tables

func init() {
	tableMap = ddl.MustNewTables(
		ddl.WithTable(
			"catalog_category_anc_categs_index_idx",
			&ddl.Column{
				Field:      "category_id",
				ColumnType: "int(10) unsigned",
				Key:        "MUL",
				Default:    dml.MakeNullString("0"),
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
				Default:    dml.MakeNullString("0"),
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

	tableMap.Upsert(ddl.NewTable(
		"catalog_category_anc_products_index_idx",
		&ddl.Column{
			Field:      "category_id",
			ColumnType: "int(10) unsigned",
			Default:    dml.MakeNullString("0"),
			Extra:      "",
		},
		&ddl.Column{
			Field:      "product_id",
			ColumnType: "int(10) unsigned",
			Key:        "",
			Default:    dml.MakeNullString("0"),
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
	)
	tableMap.Upsert(ddl.NewTable(
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
	)
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
		err := tableMap.MustTable("catalog_category_anc_categs_index_tmp").Truncate(context.TODO(), dbc.DB)
		assert.NoError(t, err, "%+v", err)
	})

	t.Run("Invalid table Name", func(t *testing.T) {
		tbl := ddl.NewTable("product")
		tbl.IsView = true
		err := tbl.Rename(context.TODO(), nil, "namecatalog_category_anc_categs_index_tmpcatalog_category_anc_categs_")
		assert.True(t, errors.NotValid.Match(err), "%+v", err)
	})
}

func TestTable_Rename(t *testing.T) {
	t.Parallel()
	t.Run("ok", func(t *testing.T) {
		dbc, dbMock := dmltest.MockDB(t)
		defer dmltest.MockClose(t, dbc, dbMock)

		dbMock.ExpectExec("RENAME TABLE `catalog_category_anc_categs_index_tmp` TO `catalog_category_anc_categs`").
			WillReturnResult(sqlmock.NewResult(0, 0))
		err := tableMap.MustTable("catalog_category_anc_categs_index_tmp").Rename(context.TODO(), dbc.DB, "catalog_category_anc_categs")
		assert.NoError(t, err, "%+v", err)
	})

	t.Run("Invalid table Name", func(t *testing.T) {
		tbl := ddl.NewTable("product")
		tbl.IsView = true
		err := tbl.Rename(context.TODO(), nil, "namecatalog_category_anc_categs_index_tmpcatalog_category_anc_categs_")
		assert.True(t, errors.NotValid.Match(err), "%+v", err)
	})
}

func TestTable_Swap(t *testing.T) {
	t.Parallel()

	t.Run("ok", func(t *testing.T) {
		dbc, dbMock := dmltest.MockDB(t)
		defer dmltest.MockClose(t, dbc, dbMock)

		dbMock.ExpectExec("RENAME TABLE `catalog_category_anc_categs_index_tmp` TO `catalog_category_anc_categs_index_tmp_[0-9]+`, `catalog_category_anc_categs_NEW` TO `catalog_category_anc_categs_index_tmp`,`catalog_category_anc_categs_index_tmp_[0-9]+` TO `catalog_category_anc_categs_NEW`").
			WillReturnResult(sqlmock.NewResult(0, 0))
		err := tableMap.MustTable("catalog_category_anc_categs_index_tmp").Swap(context.TODO(), dbc.DB, "catalog_category_anc_categs_NEW")
		assert.NoError(t, err, "%+v", err)
	})

	t.Run("Invalid table Name", func(t *testing.T) {
		tbl := ddl.NewTable("product")
		tbl.IsView = true
		err := tbl.Swap(context.TODO(), nil, "namecatalog_category_anc_categs_index_tmpcatalog_category_anc_categs_")
		assert.True(t, errors.NotValid.Match(err), "%+v", err)
	})
}

func TestTable_Drop(t *testing.T) {
	t.Parallel()
	t.Run("ok", func(t *testing.T) {
		dbc, dbMock := dmltest.MockDB(t)
		defer dmltest.MockClose(t, dbc, dbMock)

		dbMock.ExpectExec("DROP TABLE IF EXISTS `catalog_category_anc_categs_index_tmp`").
			WillReturnResult(sqlmock.NewResult(0, 0))
		err := tableMap.MustTable("catalog_category_anc_categs_index_tmp").Drop(context.TODO(), dbc.DB)
		assert.NoError(t, err, "%+v", err)
	})
	t.Run("Invalid table Name", func(t *testing.T) {
		tbl := ddl.NewTable("produ™€ct")
		tbl.IsView = true
		err := tbl.Drop(context.TODO(), nil)
		assert.True(t, errors.NotValid.Match(err), "%+v", err)
	})
}

func TestTable_LoadDataInfile(t *testing.T) {
	t.Parallel()

	t.Run("default options", func(t *testing.T) {
		dbc, dbMock := dmltest.MockDB(t)
		defer dmltest.MockClose(t, dbc, dbMock)

		dbMock.ExpectExec(dmltest.SQLMockQuoteMeta("LOAD DATA LOCAL INFILE 'non-existent.csv' INTO TABLE `admin_user` (user_id,email,first_name,username) ;")).
			WillReturnResult(sqlmock.NewResult(0, 0))

		err := tableMap.MustTable("admin_user").LoadDataInfile(context.TODO(), dbc.DB, "non-existent.csv", ddl.InfileOptions{})
		assert.NoError(t, err, "%+v", err)
	})

	t.Run("all options", func(t *testing.T) {
		dbc, dbMock := dmltest.MockDB(t)
		defer dmltest.MockClose(t, dbc, dbMock)

		dbMock.ExpectExec(dmltest.SQLMockQuoteMeta("LOAD DATA LOCAL INFILE 'non-existent.csv' REPLACE INTO TABLE `admin_user` FIELDS TERMINATED BY '|' OPTIONALLY ENCLOSED BY '+' ESCAPED BY '\"' LINES TERMINATED BY ' ' STARTING BY '###' IGNORE 1 LINES (user_id,@email,@username,) SET username=UPPER(@username), email=UPPER(@email);")).
			WillReturnResult(sqlmock.NewResult(0, 0))
		err := tableMap.MustTable("admin_user").LoadDataInfile(context.TODO(), dbc.DB, "non-existent.csv", ddl.InfileOptions{
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
		})
		assert.NoError(t, err, "%+v", err)
	})
}

func TestTable_Artisan_Methods(t *testing.T) {
	t.Parallel()

	dbc, dbMock := dmltest.MockDB(t)
	defer dmltest.MockClose(t, dbc, dbMock)

	tblAdmUser := tableMap.MustTable("admin_user")
	tblAdmUser.DB = dbc.DB

	t.Run("Insert", func(t *testing.T) {
		dbMock.ExpectExec(dmltest.SQLMockQuoteMeta("INSERT INTO `admin_user` (`email`,`first_name`) VALUES (?,?),(?,?)")).
			WithArgs("a@b.c", "Franz", "d@e.f", "Sissi").
			WillReturnResult(sqlmock.NewResult(11, 0))

		res, err := tblAdmUser.Insert().WithArgs().
			String("a@b.c").String("Franz").
			String("d@e.f").String("Sissi").
			ExecContext(context.Background())
		require.NoError(t, err)
		id, err := res.LastInsertId()
		require.NoError(t, err)
		assert.Exactly(t, int64(11), id)
	})

	t.Run("DeleteByPK", func(t *testing.T) {
		dbMock.ExpectExec(dmltest.SQLMockQuoteMeta("DELETE FROM `admin_user` WHERE (`user_id` IN ?)")).
			WithArgs("a@b.c", "d@e.f").
			WillReturnResult(sqlmock.NewResult(0, 2))

		res, err := tblAdmUser.DeleteByPK().WithArgs().
			String("a@b.c").String("d@e.f").
			ExecContext(context.Background())
		require.NoError(t, err)
		id, err := res.RowsAffected()
		require.NoError(t, err)
		assert.Exactly(t, int64(2), id)
	})

	t.Run("SelectByPK", func(t *testing.T) {
		dbMock.ExpectQuery(dmltest.SQLMockQuoteMeta("SELECT `user_id`, `email`, `first_name`, `username` FROM `admin_user` AS `main_table` WHERE (`user_id` IN ?)")).
			WithArgs(int64(234)).
			WillReturnRows(sqlmock.NewRows([]string{"user_id", "email", "first_name", "username"}))

		rows, err := tblAdmUser.SelectByPK().WithArgs().Int64(234).QueryContext(context.Background())
		require.NoError(t, err)
		require.NoError(t, rows.Close())
	})

	t.Run("SelectAll", func(t *testing.T) {
		dbMock.ExpectQuery(dmltest.SQLMockQuoteMeta("SELECT `user_id`, `email`, `first_name`, `username` FROM `admin_user` AS `main_table`")).
			WithArgs().
			WillReturnRows(sqlmock.NewRows([]string{"user_id", "email", "first_name", "username"}))

		rows, err := tblAdmUser.SelectAll().WithArgs().QueryContext(context.Background())
		require.NoError(t, err)
		require.NoError(t, rows.Close())
	})

	t.Run("UpdateByPK", func(t *testing.T) {
		dbMock.ExpectExec(dmltest.SQLMockQuoteMeta("UPDATE `admin_user` SET `email`=?, `first_name`=? WHERE (`user_id` = ?)")).
			WithArgs("a@b.c", "Franz", int64(3)).
			WillReturnResult(sqlmock.NewResult(0, 1))

		res, err := tblAdmUser.UpdateByPK().WithArgs().
			String("a@b.c").String("Franz").
			Int64(3).
			ExecContext(context.Background())
		require.NoError(t, err)
		id, err := res.RowsAffected()
		require.NoError(t, err)
		assert.Exactly(t, int64(1), id)
	})

}
