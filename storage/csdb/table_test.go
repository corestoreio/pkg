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
	"testing"

	"regexp"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/corestoreio/csfw/storage/csdb"
	"github.com/corestoreio/csfw/storage/dbr"
	"github.com/corestoreio/csfw/util/cstesting"
	"github.com/corestoreio/csfw/util/null"
	"github.com/stretchr/testify/assert"
)

const (
	table1 = iota // must start with 0 because of the for loops
	table2
	table3
	table4
	table5
)

var tableMap *csdb.Tables

// Returns a session that's not backed by a database
func createFakeSession() *dbr.Session {
	cxn, err := dbr.NewConnection()
	if err != nil {
		panic(err)
	}
	return cxn.NewSession()
}

func init() {
	tableMap = csdb.MustNewTables(
		csdb.WithTable(
			table1,
			"catalog_category_anc_categs_index_idx",
			&csdb.Column{
				Field:      ("category_id"),
				ColumnType: ("int(10) unsigned"),
				Key:        ("MUL"),
				Default:    null.StringFrom("0"),
				Extra:      (""),
			},
			&csdb.Column{
				Field:      ("path"),
				ColumnType: ("varchar(255)"),
				Null:       "YES",
				Key:        ("MUL"),
				Extra:      (""),
			},
		),
		csdb.WithTable(
			table2,
			"catalog_category_anc_categs_index_tmp",
			&csdb.Column{
				Field:      ("category_id"),
				ColumnType: ("int(10) unsigned"),
				Key:        ("PRI"),
				Default:    null.StringFrom("0"),
				Extra:      (""),
			},
			&csdb.Column{
				Field:      ("path"),
				ColumnType: ("varchar(255)"),
				Null:       "YES",
				Extra:      (""),
			},
		),
	)

	tableMap.Upsert(table3, csdb.NewTable(
		"catalog_category_anc_products_index_idx",
		&csdb.Column{
			Field:      ("category_id"),
			ColumnType: ("int(10) unsigned"),
			Default:    null.StringFrom("0"),
			Extra:      (""),
		},
		&csdb.Column{
			Field:      ("product_id"),
			ColumnType: ("int(10) unsigned"),
			Key:        (""),
			Default:    null.StringFrom("0"),
			Extra:      (""),
		},
		&csdb.Column{
			Field:      ("position"),
			ColumnType: ("int(10) unsigned"),
			Null:       "YES",
			Key:        (""),
			Extra:      (""),
		},
	),
	)
	tableMap.Upsert(table4, csdb.NewTable(
		"admin_user",
		&csdb.Column{
			Field:      ("user_id"),
			ColumnType: ("int(10) unsigned"),
			Key:        ("PRI"),
			Extra:      ("auto_increment"),
		},
		&csdb.Column{
			Field:      ("email"),
			ColumnType: ("varchar(128)"),
			Null:       "YES",
			Key:        (""),
			Extra:      (""),
		},
		&csdb.Column{
			Field:      ("username"),
			ColumnType: ("varchar(40)"),
			Null:       "YES",
			Key:        ("UNI"),
			Extra:      (""),
		},
	),
	)
}

func mustStructure(i int) *csdb.Table {
	st1, err := tableMap.Table(i)
	if err != nil {
		panic(err)
	}
	return st1
}

func TestTableStructure(t *testing.T) {

	sValid, err := tableMap.Table(table1)
	assert.NotNil(t, sValid)
	assert.NoError(t, err)

	assert.Equal(t, "catalog_category_anc_categs_index_tmp", tableMap.Name(table2))
	assert.Equal(t, "", tableMap.Name(table5))

	sInvalid, err := tableMap.Table(table5)
	assert.Nil(t, sInvalid)
	assert.Error(t, err)

	selectBuilder := sValid.Select()
	selectString, _, err := selectBuilder.ToSQL()
	assert.Equal(t, "SELECT `main_table`.`category_id`, `main_table`.`path` FROM `catalog_category_anc_categs_index_idx` AS `main_table`", selectString)
	assert.NoError(t, err)
}

func TestTableStructureTableAliasQuote(t *testing.T) {

	want := map[string]string{
		"catalog_category_anc_categs_index_idx":   "`catalog_category_anc_categs_index_idx` AS `alias`",
		"catalog_category_anc_categs_index_tmp":   "`catalog_category_anc_categs_index_tmp` AS `alias`",
		"catalog_category_anc_products_index_idx": "`catalog_category_anc_products_index_idx` AS `alias`",
		"admin_user":                              "`admin_user` AS `alias`",
	}
	for i := 0; i < tableMap.Len(); i++ {
		table, err := tableMap.Table(i)
		if err != nil {
			t.Fatalf("%+v", err)
		}
		have := table.TableAliasQuote("alias")
		assert.EqualValues(t, want[table.Name], have, "Table %s", table.Name)
	}
}

func TestTableStructureColumnAliasQuote(t *testing.T) {

	want := map[string][]string{
		"catalog_category_anc_categs_index_idx":   {"`alias`.`category_id`", "`alias`.`path`"},
		"catalog_category_anc_categs_index_tmp":   {"`alias`.`path`"},
		"catalog_category_anc_products_index_idx": {"`alias`.`category_id`", "`alias`.`product_id`", "`alias`.`position`"},
		"admin_user":                              {"`alias`.`email`", "`alias`.`username`"},
	}
	for i := 0; i < tableMap.Len(); i++ {
		table, err := tableMap.Table(i)
		if err != nil {
			t.Error(err)
		}
		have := table.ColumnAliasQuote("alias")
		assert.EqualValues(t, want[table.Name], have, "Table %s", table.Name)
	}
}

func TestTableStructureAllColumnAliasQuote(t *testing.T) {

	want := map[string][]string{
		"catalog_category_anc_categs_index_idx":   {"`alias`.`category_id`", "`alias`.`path`"},
		"catalog_category_anc_categs_index_tmp":   {"`alias`.`category_id`", "`alias`.`path`"},
		"catalog_category_anc_products_index_idx": {"`alias`.`category_id`", "`alias`.`product_id`", "`alias`.`position`"},
		"admin_user":                              {"`alias`.`user_id`", "`alias`.`email`", "`alias`.`username`"},
	}
	for i := 0; i < tableMap.Len(); i++ {
		table, err := tableMap.Table(i)
		if err != nil {
			t.Error(err)
		}
		have := table.AllColumnAliasQuote("alias")
		assert.EqualValues(t, want[table.Name], have, "Table %s", table.Name)
	}
}

func TestTableStructureIn(t *testing.T) {

	want := map[string]bool{
		"catalog_category_anc_categs_index_idx":   true,
		"catalog_category_anc_categs_index_tmp":   true,
		"catalog_category_anc_products_index_idx": false,
	}
	for i := 0; i < tableMap.Len(); i++ {
		table, err := tableMap.Table(i)
		if err != nil {
			t.Error(err)
		}
		have := table.In("path")
		assert.EqualValues(t, want[table.Name], have, "Table %s", table.Name)
	}

	want2 := map[string]bool{
		"catalog_category_anc_categs_index_idx":   true,
		"catalog_category_anc_categs_index_tmp":   true,
		"catalog_category_anc_products_index_idx": true,
	}
	for i := 0; i < tableMap.Len(); i++ {
		table, err := tableMap.Table(i)
		if err != nil {
			t.Error(err)
		}
		have := table.In("category_id")
		assert.EqualValues(t, want2[table.Name], have, "Table %s", table.Name)
	}
}

func TestTable_Truncate(t *testing.T) {
	t.Parallel()

	dbc, dbMock := cstesting.MockDB(t)
	defer func() {
		dbMock.ExpectClose()
		assert.NoError(t, dbc.Close())
		if err := dbMock.ExpectationsWereMet(); err != nil {
			t.Error("there were unfulfilled expections", err)
		}
	}()

	dbMock.ExpectExec("TRUNCATE TABLE `catalog_category_anc_categs_index_tmp`").WillReturnResult(sqlmock.NewResult(0, 0))
	err := tableMap.MustTable(table2).Truncate(dbc.DB)
	assert.NoError(t, err, "%+v", err)
}

func TestTable_Rename(t *testing.T) {
	t.Parallel()

	dbc, dbMock := cstesting.MockDB(t)
	defer func() {
		dbMock.ExpectClose()
		assert.NoError(t, dbc.Close())
		if err := dbMock.ExpectationsWereMet(); err != nil {
			t.Error("there were unfulfilled expections", err)
		}
	}()

	dbMock.ExpectExec("RENAME TABLE `catalog_category_anc_categs_index_tmp` TO `catalog_category_anc_categs`").
		WillReturnResult(sqlmock.NewResult(0, 0))
	err := tableMap.MustTable(table2).Rename(dbc.DB, "catalog_category_anc_categs")
	assert.NoError(t, err, "%+v", err)
}

func TestTable_Swap(t *testing.T) {
	t.Parallel()

	dbc, dbMock := cstesting.MockDB(t)
	defer func() {
		dbMock.ExpectClose()
		assert.NoError(t, dbc.Close())
		if err := dbMock.ExpectationsWereMet(); err != nil {
			t.Error("there were unfulfilled expections", err)
		}
	}()

	dbMock.ExpectExec("RENAME TABLE `catalog_category_anc_categs_index_tmp` TO `catalog_category_anc_categs_index_tmp_swap_[0-9]+`, `catalog_category_anc_categs_NEW` TO `catalog_category_anc_categs_index_tmp`,`catalog_category_anc_categs_index_tmp_swap_[0-9]+` TO `catalog_category_anc_categs_NEW`").
		WillReturnResult(sqlmock.NewResult(0, 0))
	err := tableMap.MustTable(table2).Swap(dbc.DB, "catalog_category_anc_categs_NEW")
	assert.NoError(t, err, "%+v", err)
}

func TestTable_Drop(t *testing.T) {
	t.Parallel()

	dbc, dbMock := cstesting.MockDB(t)
	defer func() {
		dbMock.ExpectClose()
		assert.NoError(t, dbc.Close())
		if err := dbMock.ExpectationsWereMet(); err != nil {
			t.Error("there were unfulfilled expections", err)
		}
	}()

	dbMock.ExpectExec("DROP TABLE IF EXISTS `catalog_category_anc_categs_index_tmp`").
		WillReturnResult(sqlmock.NewResult(0, 0))
	err := tableMap.MustTable(table2).Drop(dbc.DB)
	assert.NoError(t, err, "%+v", err)
}

func TestTable_LoadDataInfile(t *testing.T) {
	t.Parallel()

	t.Run("default options", func(t *testing.T) {
		dbc, dbMock := cstesting.MockDB(t)
		defer func() {
			dbMock.ExpectClose()
			assert.NoError(t, dbc.Close())
			if err := dbMock.ExpectationsWereMet(); err != nil {
				t.Error("there were unfulfilled expections", err)
			}
		}()

		dbMock.ExpectExec(regexp.QuoteMeta("LOAD DATA LOCAL INFILE 'non-existent.csv' INTO TABLE `admin_user` (user_id,email,username) ;")).
			WillReturnResult(sqlmock.NewResult(0, 0))
		err := tableMap.MustTable(table4).LoadDataInfile(dbc.DB, "non-existent.csv", csdb.InfileOptions{})
		assert.NoError(t, err, "%+v", err)
	})

	t.Run("all options", func(t *testing.T) {
		dbc, dbMock := cstesting.MockDB(t)
		defer func() {
			dbMock.ExpectClose()
			assert.NoError(t, dbc.Close())
			if err := dbMock.ExpectationsWereMet(); err != nil {
				t.Error("there were unfulfilled expections", err)
			}
		}()

		dbMock.ExpectExec(regexp.QuoteMeta("LOAD DATA LOCAL INFILE 'non-existent.csv' REPLACE  INTO TABLE `admin_user` FIELDS TERMINATED BY '|' OPTIONALLY  ENCLOSED BY '+' ESCAPED BY '\"'\n LINES  TERMINATED BY '\r\n' STARTING BY '###'\nIGNORE 1 LINES\n (user_id,@email,@username)\nSET username=UPPER(@username),\nemail=UPPER(@email);")).
			WillReturnResult(sqlmock.NewResult(0, 0))
		err := tableMap.MustTable(table4).LoadDataInfile(dbc.DB, "non-existent.csv", csdb.InfileOptions{
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
