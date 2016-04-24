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

	"github.com/corestoreio/csfw/storage/csdb"
	"github.com/corestoreio/csfw/storage/dbr"
	"github.com/stretchr/testify/assert"
)

const (
	table1 csdb.Index = iota
	table2
	table3
	table4
	table5
)

var tableMap csdb.TableManager

// Returns a session that's not backed by a database
func createFakeSession() *dbr.Session {
	cxn, err := dbr.NewConnection()
	if err != nil {
		panic(err)
	}
	return cxn.NewSession()
}

func init() {
	tableMap = csdb.MustNewTableService(
		csdb.WithTable(
			table1,
			"catalog_category_anc_categs_index_idx",
			csdb.Column{
				Field:   dbr.NewNullString("category_id"),
				Type:    dbr.NewNullString("int(10) unsigned"),
				Null:    dbr.NewNullString("NO"),
				Key:     dbr.NewNullString("MUL"),
				Default: dbr.NewNullString("0"),
				Extra:   dbr.NewNullString(""),
			},
			csdb.Column{
				Field:   dbr.NewNullString("path"),
				Type:    dbr.NewNullString("varchar(255)"),
				Null:    dbr.NewNullString("YES"),
				Key:     dbr.NewNullString("MUL"),
				Default: dbr.NullString{},
				Extra:   dbr.NewNullString(""),
			},
		),
		csdb.WithTable(
			table2,
			"catalog_category_anc_categs_index_tmp",
			csdb.Column{
				Field:   dbr.NewNullString("category_id"),
				Type:    dbr.NewNullString("int(10) unsigned"),
				Null:    dbr.NewNullString("NO"),
				Key:     dbr.NewNullString("PRI"),
				Default: dbr.NewNullString("0"),
				Extra:   dbr.NewNullString(""),
			},
			csdb.Column{
				Field:   dbr.NewNullString("path"),
				Type:    dbr.NewNullString("varchar(255)"),
				Null:    dbr.NewNullString("YES"),
				Key:     dbr.NewNullString(nil),
				Default: dbr.NullString{},
				Extra:   dbr.NewNullString(""),
			},
		),
	)

	tableMap.Append(table3, csdb.NewTable(
		"catalog_category_anc_products_index_idx",
		csdb.Column{
			Field:   dbr.NewNullString("category_id"),
			Type:    dbr.NewNullString("int(10) unsigned"),
			Null:    dbr.NewNullString("NO"),
			Key:     dbr.NewNullString(nil),
			Default: dbr.NewNullString("0"),
			Extra:   dbr.NewNullString(""),
		},
		csdb.Column{
			Field:   dbr.NewNullString("product_id"),
			Type:    dbr.NewNullString("int(10) unsigned"),
			Null:    dbr.NewNullString("NO"),
			Key:     dbr.NewNullString(""),
			Default: dbr.NewNullString("0"),
			Extra:   dbr.NewNullString(""),
		},
		csdb.Column{
			Field:   dbr.NewNullString("position"),
			Type:    dbr.NewNullString("int(10) unsigned"),
			Null:    dbr.NewNullString("YES"),
			Key:     dbr.NewNullString(""),
			Default: dbr.NullString{},
			Extra:   dbr.NewNullString(""),
		},
	),
	)
	tableMap.Append(table4, csdb.NewTable(
		"admin_user",
		csdb.Column{
			Field:   dbr.NewNullString("user_id"),
			Type:    dbr.NewNullString("int(10) unsigned"),
			Null:    dbr.NewNullString("NO"),
			Key:     dbr.NewNullString("PRI"),
			Default: dbr.NullString{},
			Extra:   dbr.NewNullString("auto_increment"),
		},
		csdb.Column{
			Field:   dbr.NewNullString("email"),
			Type:    dbr.NewNullString("varchar(128)"),
			Null:    dbr.NewNullString("YES"),
			Key:     dbr.NewNullString(""),
			Default: dbr.NullString{},
			Extra:   dbr.NewNullString(""),
		},
		csdb.Column{
			Field:   dbr.NewNullString("username"),
			Type:    dbr.NewNullString("varchar(40)"),
			Null:    dbr.NewNullString("YES"),
			Key:     dbr.NewNullString("UNI"),
			Default: dbr.NullString{},
			Extra:   dbr.NewNullString(""),
		},
	),
	)
}

func mustStructure(i csdb.Index) *csdb.Table {
	st1, err := tableMap.Structure(i)
	if err != nil {
		panic(err)
	}
	return st1
}

func TestTableStructure(t *testing.T) {

	sValid, err := tableMap.Structure(table1)
	assert.NotNil(t, sValid)
	assert.NoError(t, err)

	assert.Equal(t, "catalog_category_anc_categs_index_tmp", tableMap.Name(table2))
	assert.Equal(t, "", tableMap.Name(table5))

	sInvalid, err := tableMap.Structure(table5)
	assert.Nil(t, sInvalid)
	assert.Error(t, err)

	dbrSess := createFakeSession()
	selectBuilder, err := sValid.Select(dbrSess)
	assert.NoError(t, err)
	selectString, _, err := selectBuilder.ToSql()
	assert.Equal(t, "SELECT `main_table`.`category_id`, `main_table`.`path` FROM `catalog_category_anc_categs_index_idx` AS `main_table`", selectString)
	assert.NoError(t, err)

	selectBuilder, err = sInvalid.Select(dbrSess)
	assert.Error(t, err)
	assert.Nil(t, selectBuilder)
}

func TestTableStructureTableAliasQuote(t *testing.T) {

	want := map[string]string{
		"catalog_category_anc_categs_index_idx":   "`catalog_category_anc_categs_index_idx` AS `alias`",
		"catalog_category_anc_categs_index_tmp":   "`catalog_category_anc_categs_index_tmp` AS `alias`",
		"catalog_category_anc_products_index_idx": "`catalog_category_anc_products_index_idx` AS `alias`",
		"admin_user":                              "`admin_user` AS `alias`",
	}
	for i := csdb.Index(0); tableMap.Next(i); i++ {
		table, err := tableMap.Structure(i)
		if err != nil {
			t.Error(err)
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
	for i := csdb.Index(0); tableMap.Next(i); i++ {
		table, err := tableMap.Structure(i)
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
	for i := csdb.Index(0); tableMap.Next(i); i++ {
		table, err := tableMap.Structure(i)
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
	for i := csdb.Index(0); tableMap.Next(i); i++ {
		table, err := tableMap.Structure(i)
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
	for i := csdb.Index(0); i < tableMap.Len(); i++ {
		table, err := tableMap.Structure(i)
		if err != nil {
			t.Error(err)
		}
		have := table.In("category_id")
		assert.EqualValues(t, want2[table.Name], have, "Table %s", table.Name)
	}
}
