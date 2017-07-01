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

package csdb_test

import (
	"context"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/corestoreio/csfw/storage/csdb"
	"github.com/corestoreio/csfw/storage/dbr"
	"github.com/corestoreio/csfw/util/cstesting"
	"github.com/corestoreio/errors"
	"github.com/stretchr/testify/assert"
)

var _ dbr.QueryBuilder = (*csdb.Table)(nil)
var _ dbr.Scanner = (*csdb.Table)(nil)

var tableMap *csdb.Tables

func init() {
	tableMap = csdb.MustNewTables(
		csdb.WithTable(
			"catalog_category_anc_categs_index_idx",
			&csdb.Column{
				Field:      "category_id",
				ColumnType: "int(10) unsigned",
				Key:        "MUL",
				Default:    dbr.MakeNullString("0"),
				Extra:      "",
			},
			&csdb.Column{
				Field:      "path",
				ColumnType: "varchar(255)",
				Null:       "YES",
				Key:        "MUL",
				Extra:      "",
			},
		),
		csdb.WithTable(
			"catalog_category_anc_categs_index_tmp",
			&csdb.Column{
				Field:      "category_id",
				ColumnType: "int(10) unsigned",
				Key:        "PRI",
				Default:    dbr.MakeNullString("0"),
				Extra:      "",
			},
			&csdb.Column{
				Field:      "path",
				ColumnType: "varchar(255)",
				Null:       "YES",
				Extra:      "",
			},
		),
	)

	tableMap.Upsert(csdb.NewTable(
		"catalog_category_anc_products_index_idx",
		&csdb.Column{
			Field:      "category_id",
			ColumnType: "int(10) unsigned",
			Default:    dbr.MakeNullString("0"),
			Extra:      "",
		},
		&csdb.Column{
			Field:      "product_id",
			ColumnType: "int(10) unsigned",
			Key:        "",
			Default:    dbr.MakeNullString("0"),
			Extra:      "",
		},
		&csdb.Column{
			Field:      "position",
			ColumnType: "int(10) unsigned",
			Null:       "YES",
			Key:        "",
			Extra:      "",
		},
	),
	)
	tableMap.Upsert(csdb.NewTable(
		"admin_user",
		&csdb.Column{
			Field:      "user_id",
			ColumnType: "int(10) unsigned",
			Key:        "PRI",
			Extra:      "auto_increment",
		},
		&csdb.Column{
			Field:      "email",
			ColumnType: "varchar(128)",
			Null:       "YES",
			Key:        "",
			Extra:      "",
		},
		&csdb.Column{
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
		dbc, dbMock := cstesting.MockDB(t)
		defer cstesting.MockClose(t, dbc, dbMock)

		dbMock.ExpectExec("TRUNCATE TABLE `catalog_category_anc_categs_index_tmp`").WillReturnResult(sqlmock.NewResult(0, 0))
		err := tableMap.MustTable("catalog_category_anc_categs_index_tmp").Truncate(context.TODO(), dbc.DB)
		assert.NoError(t, err, "%+v", err)
	})

	t.Run("Invalid table Name", func(t *testing.T) {
		tbl := csdb.NewTable("product")
		tbl.IsView = true
		err := tbl.Rename(context.TODO(), nil, "namecatalog_category_anc_categs_index_tmpcatalog_category_anc_categs_")
		assert.True(t, errors.IsNotValid(err), "%+v", err)
	})
}

func TestTable_Rename(t *testing.T) {
	t.Parallel()
	t.Run("ok", func(t *testing.T) {
		dbc, dbMock := cstesting.MockDB(t)
		defer cstesting.MockClose(t, dbc, dbMock)

		dbMock.ExpectExec("RENAME TABLE `catalog_category_anc_categs_index_tmp` TO `catalog_category_anc_categs`").
			WillReturnResult(sqlmock.NewResult(0, 0))
		err := tableMap.MustTable("catalog_category_anc_categs_index_tmp").Rename(context.TODO(), dbc.DB, "catalog_category_anc_categs")
		assert.NoError(t, err, "%+v", err)
	})

	t.Run("Invalid table Name", func(t *testing.T) {
		tbl := csdb.NewTable("product")
		tbl.IsView = true
		err := tbl.Rename(context.TODO(), nil, "namecatalog_category_anc_categs_index_tmpcatalog_category_anc_categs_")
		assert.True(t, errors.IsNotValid(err), "%+v", err)
	})
}

func TestTable_Swap(t *testing.T) {
	t.Parallel()

	t.Run("ok", func(t *testing.T) {
		dbc, dbMock := cstesting.MockDB(t)
		defer cstesting.MockClose(t, dbc, dbMock)

		dbMock.ExpectExec("RENAME TABLE `catalog_category_anc_categs_index_tmp` TO `catalog_category_anc_categs_index_tmp_[0-9]+`, `catalog_category_anc_categs_NEW` TO `catalog_category_anc_categs_index_tmp`,`catalog_category_anc_categs_index_tmp_[0-9]+` TO `catalog_category_anc_categs_NEW`").
			WillReturnResult(sqlmock.NewResult(0, 0))
		err := tableMap.MustTable("catalog_category_anc_categs_index_tmp").Swap(context.TODO(), dbc.DB, "catalog_category_anc_categs_NEW")
		assert.NoError(t, err, "%+v", err)
	})

	t.Run("Invalid table Name", func(t *testing.T) {
		tbl := csdb.NewTable("product")
		tbl.IsView = true
		err := tbl.Swap(context.TODO(), nil, "namecatalog_category_anc_categs_index_tmpcatalog_category_anc_categs_")
		assert.True(t, errors.IsNotValid(err), "%+v", err)
	})
}

func TestTable_Drop(t *testing.T) {
	t.Parallel()
	t.Run("ok", func(t *testing.T) {
		dbc, dbMock := cstesting.MockDB(t)
		defer cstesting.MockClose(t, dbc, dbMock)

		dbMock.ExpectExec("DROP TABLE IF EXISTS `catalog_category_anc_categs_index_tmp`").
			WillReturnResult(sqlmock.NewResult(0, 0))
		err := tableMap.MustTable("catalog_category_anc_categs_index_tmp").Drop(context.TODO(), dbc.DB)
		assert.NoError(t, err, "%+v", err)
	})
	t.Run("Invalid table Name", func(t *testing.T) {
		tbl := csdb.NewTable("produ™€ct")
		tbl.IsView = true
		err := tbl.Drop(context.TODO(), nil)
		assert.True(t, errors.IsNotValid(err), "%+v", err)
	})
}

func TestTable_LoadDataInfile(t *testing.T) {
	t.Parallel()

	t.Run("default options", func(t *testing.T) {
		dbc, dbMock := cstesting.MockDB(t)
		defer cstesting.MockClose(t, dbc, dbMock)

		dbMock.ExpectExec(cstesting.SQLMockQuoteMeta("LOAD DATA LOCAL INFILE 'non-existent.csv' INTO TABLE `admin_user` (user_id,email,username) ;")).
			WillReturnResult(sqlmock.NewResult(0, 0))

		err := tableMap.MustTable("admin_user").LoadDataInfile(context.TODO(), dbc.DB, "non-existent.csv", csdb.InfileOptions{})
		assert.NoError(t, err, "%+v", err)
	})

	t.Run("all options", func(t *testing.T) {
		dbc, dbMock := cstesting.MockDB(t)
		defer cstesting.MockClose(t, dbc, dbMock)

		dbMock.ExpectExec(cstesting.SQLMockQuoteMeta("LOAD DATA LOCAL INFILE 'non-existent.csv' REPLACE  INTO TABLE `admin_user` FIELDS TERMINATED BY '|' OPTIONALLY  ENCLOSED BY '+' ESCAPED BY '\"'\n LINES  TERMINATED BY '\r\n' STARTING BY '###'\nIGNORE 1 LINES\n (user_id,@email,@username)\nSET username=UPPER(@username),\nemail=UPPER(@email);")).
			WillReturnResult(sqlmock.NewResult(0, 0))
		err := tableMap.MustTable("admin_user").LoadDataInfile(context.TODO(), dbc.DB, "non-existent.csv", csdb.InfileOptions{
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
