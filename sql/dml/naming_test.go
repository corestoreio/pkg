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

package dml

import (
	"testing"

	"github.com/corestoreio/pkg/util/assert"
)

func TestMakeAlias(t *testing.T) {
	assert.Exactly(t, "`table1`", MakeIdentifier("table1").String())
	assert.Exactly(t, "`table0` AS `table1`", MakeIdentifier("table0").Alias("table1").String())
	assert.Exactly(t, "`(table1)`", MakeIdentifier("(table1)").String())
	assert.Exactly(t, "`(table1)` AS `table2`", MakeIdentifier("(table1)").Alias("table2").String())
	assert.Exactly(t, "`(table1)`", MakeIdentifier("(table1)").String())
	assert.Exactly(t, "`table1`", MakeIdentifier("table1").String())
}

func TestIds_AppendColumns(t *testing.T) {
	tests := []struct {
		have, want ids
	}{
		{
			ids{}.AppendColumns(false, "aa", "bb"),
			ids{{Name: "aa"}, {Name: "bb"}},
		},
		{
			ids{}.AppendColumns(false, "aa ASC", "bb"),
			ids{{Name: "aa", Sort: sortAscending}, {Name: "bb"}},
		},
		{
			ids{}.AppendColumns(false, "aa ASC", "bb DESC"),
			ids{{Name: "aa", Sort: sortAscending}, {Name: "bb", Sort: sortDescending}},
		},
	}
	for i, test := range tests {
		assert.Exactly(t, test.want, test.have, "Index %d", i)
	}
}

// func TestMakeExpressionAlias(t *testing.T) {
//
//	assert.Exactly(t, "(table1)", MakeExpressionAlias("(table1)", "").String())
//	assert.Exactly(t, "(table1) AS `x`", MakeExpressionAlias("(table1)", "x").String())
//	assert.Exactly(t, "(table1)", MakeExpressionAlias("(table1)", "").String())
//}

func TestMysqlQuoter_QuoteAlias(t *testing.T) {
	tests := []struct {
		name, alias, want string
	}{
		0: {"a", "", "`a`"},
		1: {"a", "b", "`a` AS `b`"},
		2: {"a", "", "`a`"},
		3: {"`c`", "", "`c`"},
		4: {"d.e", "", "`d`.`e`"},
		5: {"`d`.`e`", "", "`d`.`e`"},
		6: {"f", "g_h", "`f` AS `g_h`"},
		7: {"f", "g_h`h", "`f` AS `g_hh`"},
	}
	for i, test := range tests {
		assert.Exactly(t, test.want, Quoter.NameAlias(test.name, test.alias), "Index %d", i)
	}
}

func TestMysqlQuoter_Name(t *testing.T) {
	assert.Exactly(t, "`tableName`", Quoter.Name("tableName"))
	assert.Exactly(t, "`tableName`", Quoter.Name("table`Name"))
	assert.Exactly(t, "``", Quoter.Name(""))
	assert.Exactly(t, "`databaseName`.`tableName`", Quoter.QualifierName("databaseName", "tableName"))
	assert.Exactly(t, "`tableName`", Quoter.QualifierName("", "tableName")) // qualifier is empty
	assert.Exactly(t, "`databaseName`.`tableName`", Quoter.QualifierName("database`Name", "table`Name"))
}

func TestIsValidIdentifier(t *testing.T) {
	tests := []struct {
		have string
		want int8
	}{
		{"*", 0},
		{"table.*", 0},
		{" table.*", 2},
		{"table.col", 0},
		{"table.col ", 2},
		{"*.*", 2},
		{"table.p*", 2},
		{"`table`.*", 2},     // not valid because of backticks
		{"`table`.`col`", 2}, // not valid because of backticks
		{"", 1},
		{"a", 0},
		{"a.", 1},
		{"a.b", 0},
		{".b", 1},
		{"", 2},
		{"花间一壶酒，独酌无相亲。", 2}, // no idea what this means but found it in x/text pkg
		{"独酌无相", 2},         // no idea what this means but found it in x/text pkg
		{"Goooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooopher", 1},
		{"Gooooooooooooooooooooooooooooooooooooooooooooooooooooooooooopher", 0},
		{"Gooooooooooooooooooooooooooooooooooooooooooooooooooooooooooopher.Gooooooooooooooooooooooooooooooooooooooooooooooooooooooooooopher", 0},
		{"Gooooooooooooooooooooooooooooooooooooooooooooooooooooooooooopher.Goooooooooooooooooooooooooooooooooooooooooooooooooooooooooo0opher", 1},
		{"Goooooooooooooooooooooooooooooooooooooooooooooooooooooooooopher.Gooooooooooooooooooooooooooooooooooooooooooooooooooooooo0oph€r", 2},
		{"DATE_FORMAT(t3.period, '%Y-%m-01')", 2},
		{"1ekdsdf", 3},
		{"9ekdsdf", 3},
	}
	for i, test := range tests {
		assert.Exactly(t, test.want, isValidIdentifier(test.have), "Index %d with %q", i, test.have)
	}
}

func TestIsValidIdentifier2(t *testing.T) {
	assert.Error(t, IsValidIdentifier("DATE_FORMAT(t3.period, '%Y-%m-01')"))
	assert.NoError(t, IsValidIdentifier("table.col"))
}

func TestIds_Clone(t *testing.T) {
	t.Run("non-nil", func(t *testing.T) {
		c := MakeIdentifier("c")
		c.DerivedTable = NewSelect("x", "y").From("z")
		names := ids{
			MakeIdentifier("a"),
			MakeIdentifier("b").Alias("b2"),
			c,
		}
		names2 := names.Clone()
		notEqualPointers(t, names, names2)
		notEqualPointers(t, names[2].DerivedTable, names2[2].DerivedTable)
	})

	t.Run("nil", func(t *testing.T) {
		var names ids
		names2 := names.Clone()
		assert.Nil(t, names2)
	})
}
