// Copyright 2015-2017, Cyrill @ Schumacher.fm and the CoreStore contributors
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

package dbr_test

import (
	"context"
	"database/sql"
	"testing"

	"github.com/corestoreio/csfw/storage/dbr"
)

var _ dbr.Querier = (*benchMockQuerier)(nil)

type benchMockQuerier struct{}

func (benchMockQuerier) QueryContext(ctx context.Context, query string, args ...interface{}) (*sql.Rows, error) {
	return new(sql.Rows), nil
}

// BenchmarkSelect_Rows-4   	 1000000	      2188 ns/op	    1354 B/op	      19 allocs/op old
// BenchmarkSelect_Rows-4   	 1000000	      2223 ns/op	    1386 B/op	      20 allocs/op new
func BenchmarkSelect_Rows(b *testing.B) {

	tables := []string{"eav_attribute"}
	ctx := context.TODO()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {

		sel := dbr.NewSelect("TABLE_NAME", "COLUMN_NAME", "ORDINAL_POSITION", "COLUMN_DEFAULT", "IS_NULLABLE",
			"DATA_TYPE", "CHARACTER_MAXIMUM_LENGTH", "NUMERIC_PRECISION", "NUMERIC_SCALE",
			"COLUMN_TYPE", "COLUMN_KEY", "EXTRA", "COLUMN_COMMENT").From("information_schema.COLUMNS").
			Where(dbr.Condition(`TABLE_SCHEMA=DATABASE()`))
		sel.DB.Querier = benchMockQuerier{}
		if len(tables) > 0 {
			sel.Where(dbr.Condition("TABLE_NAME IN ?", dbr.ArgString(tables...)))
		}

		rows, err := sel.Rows(ctx)
		if err != nil {
			b.Fatalf("%+v", err)
		}
		if rows == nil {
			b.Fatal("Rows should not be nil")
		}
	}
}

var benchmarkSelectBasicSQL dbr.Arguments

func BenchmarkSelectBasicSQL(b *testing.B) {

	// Do some allocations outside the loop so they don't affect the results
	argEq := dbr.Eq{"a": dbr.ArgInt64(1, 2, 3).Operator(dbr.OperatorIn)}
	args := dbr.Arguments{dbr.ArgInt64(1), dbr.ArgString("wat")}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, args, err := dbr.NewSelect("something_id", "user_id", "other").
			From("some_table").
			Where(dbr.Condition("d = ? OR e = ?", args...)).
			Where(argEq).
			OrderByDesc("id").
			Paginate(1, 20).
			ToSQL()
		if err != nil {
			b.Fatalf("%+v", err)
		}
		benchmarkSelectBasicSQL = args
	}
}

func BenchmarkSelectFullSQL(b *testing.B) {

	// Do some allocations outside the loop so they don't affect the results
	argEq1 := dbr.Eq{"f": dbr.ArgInt64(2), "x": dbr.ArgString("hi")}
	argEq2 := dbr.Eq{"g": dbr.ArgInt64(3)}
	argEq3 := dbr.Eq{"h": dbr.ArgInt(1, 2, 3)}
	args := dbr.Arguments{dbr.ArgInt64(1), dbr.ArgString("wat")}

	b.ResetTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		_, args, err := dbr.NewSelect("a", "b", "z", "y", "x").From("c").
			Distinct().
			Where(dbr.Condition("d = ? OR e = ?", args...)).
			Where(argEq1).
			Where(argEq2).
			Where(argEq3).
			GroupBy("ab").
			GroupBy("ii").
			GroupBy("iii").
			Having(dbr.Condition("j = k"), dbr.Condition("jj = ?", dbr.ArgInt64(1))).
			Having(dbr.Condition("jjj = ?", dbr.ArgInt64(2))).
			OrderBy("l").
			OrderBy("l").
			OrderBy("l").
			Limit(7).
			Offset(8).
			ToSQL()
		if err != nil {
			b.Fatalf("%+v", err)
		}
		benchmarkSelectBasicSQL = args
	}
}

// BenchmarkSelect_Large_IN-4   	  500000	      2807 ns/op	    1216 B/op	      27 allocs/op
func BenchmarkSelect_Large_IN(b *testing.B) {

	// This tests simulates selecting many varchar attribute values for specific products.

	var entityIDs = make([]int64, 1024)
	for i := 0; i < 1024; i++ {
		entityIDs[i] = int64(i + 600)
	}

	b.ResetTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		_, args, err := dbr.NewSelect("entity_id", "attribute_id", "value").
			From("catalog_product_entity_varchar").
			Where(dbr.Condition("entity_type_id", dbr.ArgInt64(4))).
			Where(dbr.Condition("entity_id", dbr.ArgInt64(entityIDs...).Operator(dbr.OperatorIn))).
			Where(dbr.Condition("attribute_id", dbr.ArgInt64(174, 175).Operator(dbr.OperatorIn))).
			Where(dbr.Condition("store_id", dbr.ArgInt64(0))).
			ToSQL()
		if err != nil {
			b.Fatalf("%+v", err)
		}
		benchmarkSelectBasicSQL = args
	}
}

func BenchmarkSelect_ComplexAddColumns(b *testing.B) {

	var haveSQL string
	b.ResetTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		var args dbr.Arguments
		var err error
		haveSQL, args, err = dbr.NewSelect().
			AddColumns(" entity_id ,   value").
			AddColumnsQuoted("cpev.entity_type_id", "cpev.attribute_id").
			AddColumnsExprAlias("(cpev.id*3)", "weirdID").
			AddColumnsQuotedAlias("cpev.value", "value2nd").
			From("catalog_product_entity_varchar", "cpev").
			Where(dbr.Condition("entity_type_id", dbr.ArgInt64(4))).
			Where(dbr.Condition("attribute_id", dbr.ArgInt64(174, 175).Operator(dbr.OperatorIn))).
			Where(dbr.Condition("store_id", dbr.ArgInt64(0))).
			ToSQL()
		if err != nil {
			b.Fatalf("%+v", err)
		}
		benchmarkSelectBasicSQL = args
	}
	_ = haveSQL
	//b.Logf("%s", haveSQL)
	/*
		SELECT entity_id,
		       value,
		       `cpev`.`entity_type_id`,
		       `cpev`.`attribute_id`,
		       ( cpev.id * 3 ) AS `weirdID`,
		       `cpev`.`value`  AS `value2nd`
		FROM   `catalog_product_entity_varchar` AS `cpev`
		WHERE  ( `entity_type_id` = ? )
		       AND ( `attribute_id` IN ? )
		       AND ( `store_id` = ? )
	*/
}
