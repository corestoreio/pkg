// Copyright 2015, Cyrill @ Schumacher.fm and the CoreStore contributors
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

	"database/sql"

	"github.com/corestoreio/csfw/storage/csdb"
	"github.com/stretchr/testify/assert"
)

const (
	table1 csdb.Index = iota
	table2
	table3
	table4
)

var tableMap csdb.TableStructurer

func init() {
	tableMap = csdb.TableStructureSlice{
		table1: csdb.NewTableStructure(
			"catalog_category_anc_categs_index_idx",
			csdb.Column{
				Field:   sql.NullString{String: "category_id", Valid: true},
				Type:    sql.NullString{String: "int(10) unsigned", Valid: true},
				Null:    sql.NullString{String: "NO", Valid: true},
				Key:     sql.NullString{String: "MUL", Valid: true},
				Default: sql.NullString{String: "0", Valid: true},
				Extra:   sql.NullString{String: "", Valid: true},
			},
			csdb.Column{
				Field:   sql.NullString{String: "path", Valid: true},
				Type:    sql.NullString{String: "varchar(255)", Valid: true},
				Null:    sql.NullString{String: "YES", Valid: true},
				Key:     sql.NullString{String: "MUL", Valid: true},
				Default: sql.NullString{},
				Extra:   sql.NullString{String: "", Valid: true},
			},
		),
		table2: csdb.NewTableStructure(
			"catalog_category_anc_categs_index_tmp",
			csdb.Column{
				Field:   sql.NullString{String: "category_id", Valid: true},
				Type:    sql.NullString{String: "int(10) unsigned", Valid: true},
				Null:    sql.NullString{String: "NO", Valid: true},
				Key:     sql.NullString{String: "PRI", Valid: true},
				Default: sql.NullString{String: "0", Valid: true},
				Extra:   sql.NullString{String: "", Valid: true},
			},
			csdb.Column{
				Field:   sql.NullString{String: "path", Valid: true},
				Type:    sql.NullString{String: "varchar(255)", Valid: true},
				Null:    sql.NullString{String: "YES", Valid: true},
				Key:     sql.NullString{String: "", Valid: false},
				Default: sql.NullString{},
				Extra:   sql.NullString{String: "", Valid: true},
			},
		),
		table3: csdb.NewTableStructure(
			"catalog_category_anc_products_index_idx",
			csdb.Column{
				Field:   sql.NullString{String: "category_id", Valid: true},
				Type:    sql.NullString{String: "int(10) unsigned", Valid: true},
				Null:    sql.NullString{String: "NO", Valid: true},
				Key:     sql.NullString{String: "", Valid: false},
				Default: sql.NullString{String: "0", Valid: true},
				Extra:   sql.NullString{String: "", Valid: true},
			},
			csdb.Column{
				Field:   sql.NullString{String: "product_id", Valid: true},
				Type:    sql.NullString{String: "int(10) unsigned", Valid: true},
				Null:    sql.NullString{String: "NO", Valid: true},
				Key:     sql.NullString{String: "", Valid: true},
				Default: sql.NullString{String: "0", Valid: true},
				Extra:   sql.NullString{String: "", Valid: true},
			},
			csdb.Column{
				Field:   sql.NullString{String: "position", Valid: true},
				Type:    sql.NullString{String: "int(10) unsigned", Valid: true},
				Null:    sql.NullString{String: "YES", Valid: true},
				Key:     sql.NullString{String: "", Valid: true},
				Default: sql.NullString{},
				Extra:   sql.NullString{String: "", Valid: true},
			},
		),
	}
}

func TestColumn(t *testing.T) {
	t.Error("@todo")
}

func TestColumns(t *testing.T) {
	t.Error("@todo")
}

func TestTableStructure(t *testing.T) {
	dbc := csdb.MustConnectTest()
	defer dbc.Close()

	sValid, err := tableMap.Structure(table1)
	assert.NotNil(t, sValid)
	assert.NoError(t, err)

	assert.Equal(t, "catalog_category_anc_categs_index_tmp", tableMap.Name(table2))
	assert.Equal(t, "", tableMap.Name(table4))

	sInvalid, err := tableMap.Structure(table4)
	assert.Nil(t, sInvalid)
	assert.Error(t, err)

	dbrSess := dbc.NewSession()
	selectBuilder, err := sValid.Select(dbrSess)
	assert.NoError(t, err)
	selectString, _ := selectBuilder.ToSql()
	assert.Equal(t, "SELECT `main_table`.`category_id`, `main_table`.`path` FROM `catalog_category_anc_categs_index_idx` AS `main_table`", selectString)

	selectBuilder, err = sInvalid.Select(dbrSess)
	assert.Error(t, err)
	assert.Nil(t, selectBuilder)
}

func TestTableStructureTableAliasQuote(t *testing.T) {
	want := map[string]string{
		"catalog_category_anc_categs_index_idx":   "`catalog_category_anc_categs_index_idx` AS `alias`",
		"catalog_category_anc_categs_index_tmp":   "`catalog_category_anc_categs_index_tmp` AS `alias`",
		"catalog_category_anc_products_index_idx": "`catalog_category_anc_products_index_idx` AS `alias`",
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
		"catalog_category_anc_categs_index_idx":   []string{"`alias`.`category_id`", "`alias`.`path`"},
		"catalog_category_anc_categs_index_tmp":   []string{"`alias`.`path`"},
		"catalog_category_anc_products_index_idx": []string{"`alias`.`category_id`", "`alias`.`product_id`", "`alias`.`position`"},
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
		"catalog_category_anc_categs_index_idx":   []string{"`alias`.`category_id`", "`alias`.`path`"},
		"catalog_category_anc_categs_index_tmp":   []string{"`alias`.`category_id`", "`alias`.`path`"},
		"catalog_category_anc_products_index_idx": []string{"`alias`.`category_id`", "`alias`.`product_id`", "`alias`.`position`"},
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
