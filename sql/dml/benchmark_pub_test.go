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

package dml_test

import (
	"context"
	"database/sql"
	"database/sql/driver"
	"math"
	"reflect"
	"testing"
	"time"

	"github.com/corestoreio/csfw/sql/dml"
	"github.com/corestoreio/errors"
)

var benchmarkGlobalVals []interface{}

type benchMockQuerier struct{}

func (benchMockQuerier) QueryContext(ctx context.Context, query string, args ...interface{}) (*sql.Rows, error) {
	return new(sql.Rows), nil
}
func (benchMockQuerier) PrepareContext(ctx context.Context, query string) (*sql.Stmt, error) {
	return new(sql.Stmt), nil
}

// BenchmarkSelect_Rows-4   	 1000000	      2188 ns/op	    1354 B/op	      19 allocs/op old
// BenchmarkSelect_Rows-4   	 1000000	      2223 ns/op	    1386 B/op	      20 allocs/op new
func BenchmarkSelect_Rows(b *testing.B) {

	tables := []string{"eav_attribute"}
	ctx := context.TODO()
	db := benchMockQuerier{}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {

		sel := dml.NewSelect("TABLE_NAME", "COLUMN_NAME", "ORDINAL_POSITION", "COLUMN_DEFAULT", "IS_NULLABLE",
			"DATA_TYPE", "CHARACTER_MAXIMUM_LENGTH", "NUMERIC_PRECISION", "NUMERIC_SCALE",
			"COLUMN_TYPE", "COLUMN_KEY", "EXTRA", "COLUMN_COMMENT").From("information_schema.COLUMNS").
			Where(dml.Expr(`TABLE_SCHEMA=DATABASE()`)).WithDB(db)

		if len(tables) > 0 {
			sel.Where(dml.Column("TABLE_NAME IN ?").In().Strs(tables...))
		}

		rows, err := sel.Query(ctx)
		if err != nil {
			b.Fatalf("%+v", err)
		}
		if rows == nil {
			b.Fatal("Query should not be nil")
		}
	}
}

var benchmarkSelectStr string

// BenchmarkSelectBasicSQL-4 	  500000	      2542 ns/op	    1512 B/op	      18 allocs/op
// BenchmarkSelectBasicSQL-4     1000000	      2395 ns/op	    1664 B/op	      17 allocs/op <== arg value ?
// BenchmarkSelectBasicSQL-4      500000	      3060 ns/op	    2089 B/op	      22 allocs/op <== Builder Structs
// BenchmarkSelectBasicSQL-4   	  200000	      9266 ns/op	    5875 B/op	      18 allocs/op <== arg union
// BenchmarkSelectBasicSQL-4   	  500000	      3385 ns/op	    3698 B/op	      19 allocs/op
// BenchmarkSelectBasicSQL-4   	  500000	      3216 ns/op	    3281 B/op	      20 allocs/op <= union type struct, large
// BenchmarkSelectBasicSQL-4   	  500000	      2937 ns/op	    2281 B/op	      20 allocs/op no pointers
// BenchmarkSelectBasicSQL-4   	  500000	      2849 ns/op	    2257 B/op	      18 allocs/op
func BenchmarkSelectBasicSQL(b *testing.B) {

	// Do some allocations outside the loop so they don't affect the results

	aVal := []int64{1, 2, 3}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, args, err := dml.NewSelect("something_id", "user_id", "other").
			From("some_table").
			Where(
				dml.Expr("d = ? OR e = ?").Int64(1).Str("wat"),
				dml.Column("a").In().Int64s(aVal...),
			).
			OrderByDesc("id").
			Paginate(1, 20).
			ToSQL()
		if err != nil {
			b.Fatalf("%+v", err)
		}
		benchmarkGlobalVals = args
	}
}

func BenchmarkSelectFullSQL(b *testing.B) {

	sqlObj := dml.NewSelect("a", "b", "z", "y", "x").From("c").
		Distinct().
		Where(
			dml.Expr("`d` = ? OR `e` = ?").Int64(1).Str("wat"),
			dml.Column("f").Int64(2),
			dml.Column("x").Str("hi"),
			dml.Column("g").Int64(3),
			dml.Column("h").In().Ints(1, 2, 3),
		).
		GroupBy("ab").GroupBy("ii").GroupBy("iii").
		Having(dml.Expr("j = k"), dml.Column("jj").Int64(1)).
		Having(dml.Column("jjj").Int64(2)).
		OrderBy("l1").OrderBy("l2").OrderBy("l3").
		Limit(7).Offset(8).Interpolate()

	b.ResetTimer()
	b.ReportAllocs()

	// BenchmarkSelectFullSQL/NewSelect-4             300000	      5849 ns/op	    3649 B/op	      38 allocs/op
	// BenchmarkSelectFullSQL/NewSelect-4          	  200000	      6307 ns/op	    4922 B/op	      45 allocs/op <== builder structs
	// BenchmarkSelectFullSQL/NewSelect-4         	  200000	      7084 ns/op	    8212 B/op	      44 allocs/op <== Arguments
	// BenchmarkSelectFullSQL/NewSelect-4         	  200000	      6449 ns/op	    5515 B/op	      44 allocs/op no pointers
	// BenchmarkSelectFullSQL/NewSelect-4         	  200000	      6268 ns/op	    5443 B/op	      37 allocs/op
	b.Run("NewSelect", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			_, args, err := dml.NewSelect("a", "b", "z", "y", "x").From("c").
				Distinct().
				Where(
					dml.Expr("`d` = ? OR `e` = ?").Int64(1).Str("wat"),
					dml.Column("f").Int64(2),
					dml.Column("x").Str("hi"),
					dml.Column("g").Int64(3),
					dml.Column("h").In().Ints(1, 2, 3),
				).
				GroupBy("ab").GroupBy("ii").GroupBy("iii").
				Having(dml.Expr("j = k"), dml.Column("jj").Int64(1)).
				Having(dml.Column("jjj").Int64(2)).
				OrderBy("l1").OrderBy("l2").OrderBy("l3").
				Limit(7).Offset(8).
				ToSQL()
			if err != nil {
				b.Fatalf("%+v", err)
			}
			benchmarkGlobalVals = args
		}
	})

	b.Run("ToSQL Interpolate NoCache", func(b *testing.B) {
		sqlObj.IsBuildCache = false
		for i := 0; i < b.N; i++ {
			_, args, err := sqlObj.ToSQL()
			if err != nil {
				b.Fatalf("%+v", err)
			}
			benchmarkGlobalVals = args
		}
	})

	b.Run("ToSQL Interpolate Cache", func(b *testing.B) {
		sqlObj.IsBuildCache = true
		for i := 0; i < b.N; i++ {
			_, args, err := sqlObj.ToSQL()
			if err != nil {
				b.Fatalf("%+v", err)
			}
			benchmarkGlobalVals = args
		}
	})
}

func BenchmarkSelect_Large_IN(b *testing.B) {

	// This tests simulates selecting many varchar attribute values for specific products.
	var entityIDs = make([]int64, 1024)
	for i := 0; i < 1024; i++ {
		entityIDs[i] = int64(i + 1600)
	}

	b.ResetTimer()
	b.ReportAllocs()

	b.Run("SQL", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			_, args, err := dml.NewSelect("entity_id", "attribute_id", "value").
				From("catalog_product_entity_varchar").
				Where(dml.Column("entity_type_id").Int64(4)).
				Where(dml.Column("entity_id").In().Int64s(entityIDs...)).
				Where(dml.Column("attribute_id").In().Int64s(174, 175)).
				Where(dml.Column("store_id").Int(0)).
				ToSQL()
			if err != nil {
				b.Fatalf("%+v", err)
			}
			benchmarkGlobalVals = args
		}
	})

	b.Run("interpolate", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			sqlStr, args, err := dml.NewSelect("entity_id", "attribute_id", "value").
				From("catalog_product_entity_varchar").
				Where(dml.Column("entity_type_id").Int64(4)).
				Where(dml.Column("entity_id").In().Int64s(entityIDs...)).
				Where(dml.Column("attribute_id").In().Int64s(174, 175)).
				Where(dml.Column("store_id").Int(0)).
				Interpolate().
				ToSQL()
			if err != nil {
				b.Fatalf("%+v", err)
			}
			if args != nil {
				b.Fatal("Args should be nil")
			}
			benchmarkSelectStr = sqlStr
		}
	})
}

func BenchmarkSelect_ComplexAddColumns(b *testing.B) {

	var haveSQL string
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		var args []interface{}
		var err error
		haveSQL, args, err = dml.NewSelect().
			AddColumns(" entity_id ,   value").
			AddColumns("cpev.entity_type_id", "cpev.attribute_id").
			AddColumnsAliases("(cpev.id*3)", "weirdID").
			AddColumnsAliases("cpev.value", "value2nd").
			FromAlias("catalog_product_entity_varchar", "cpev").
			Where(dml.Column("entity_type_id").Int64(4)).
			Where(dml.Column("attribute_id").In().Int64s(174, 175)).
			Where(dml.Column("store_id").Int64(0)).
			ToSQL()
		if err != nil {
			b.Fatalf("%+v", err)
		}
		benchmarkGlobalVals = args
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

// BenchmarkSelect_SQLCase-4      500000	      3451 ns/op	    2032 B/op	      21 allocs/op
// BenchmarkSelect_SQLCase-4   	  500000	      3690 ns/op	    2849 B/op	      24 allocs/op
// BenchmarkSelect_SQLCase-4   	  300000	      3784 ns/op	    2433 B/op	      26 allocs/op
func BenchmarkSelect_SQLCase(b *testing.B) {
	start := time.Unix(1257894000, 0)
	end := time.Unix(1257980400, 0)
	pid := []int{4711, 815, 42}

	var haveSQL string

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		var err error
		haveSQL, benchmarkGlobalVals, err = dml.NewSelect().
			AddColumns("price", "sku", "name").
			AddColumnsConditions(
				dml.SQLCase("", "`closed`",
					"date_start <= ? AND date_end >= ?", "`open`",
					"date_start > ? AND date_end > ?", "`upcoming`",
				).Alias("is_on_sale").Times(start, end, start, end),
			).
			From("catalog_promotions").
			Where(
				dml.Column("promotion_id").
					NotIn().
					Ints(pid...),
			).
			ToSQL()
		if err != nil {
			b.Fatalf("%+v", err)
		}
	}
	_ = haveSQL
}

func BenchmarkDeleteSQL(b *testing.B) {

	b.Run("NewDelete", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			var err error
			_, benchmarkGlobalVals, err = dml.NewDelete("alpha").Where(dml.Column("a").Str("b")).Limit(1).OrderBy("id").ToSQL()
			if err != nil {
				b.Fatalf("%+v", err)
			}
		}
	})

	sqlObj := dml.NewDelete("alpha").Where(dml.Column("a").Str("b")).Limit(1).OrderBy("id").Interpolate()
	b.Run("ToSQL no cache", func(b *testing.B) {
		sqlObj.IsBuildCache = false
		for i := 0; i < b.N; i++ {
			_, args, err := sqlObj.ToSQL()
			if err != nil {
				b.Fatalf("%+v", err)
			}
			benchmarkGlobalVals = args
		}
	})

	b.Run("ToSQL with cache", func(b *testing.B) {
		sqlObj.IsBuildCache = true
		for i := 0; i < b.N; i++ {
			_, args, err := sqlObj.ToSQL()
			if err != nil {
				b.Fatalf("%+v", err)
			}
			benchmarkGlobalVals = args
		}
	})
}

func BenchmarkInsertValuesSQL(b *testing.B) {

	b.Run("NewInsert", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			_, args, err := dml.NewInsert("alpha").AddColumns("something_id", "user_id", "other").AddArguments(
				dml.MakeArgs(3).Int64(1).Int64(2).Bool(true),
			).ToSQL()
			if err != nil {
				b.Fatal(err)
			}
			benchmarkGlobalVals = args
		}
	})

	sqlObj := dml.NewInsert("alpha").AddColumns("something_id", "user_id", "other").AddArguments(
		dml.MakeArgs(3).Int64(1).Int64(2).Bool(true),
	).Interpolate()
	b.Run("ToSQL no cache", func(b *testing.B) {
		sqlObj.IsBuildCache = false
		for i := 0; i < b.N; i++ {
			_, args, err := sqlObj.ToSQL()
			if err != nil {
				b.Fatalf("%+v", err)
			}
			benchmarkGlobalVals = args
		}
	})

	b.Run("ToSQL with cache", func(b *testing.B) {
		sqlObj.BuildCache()
		for i := 0; i < b.N; i++ {
			_, args, err := sqlObj.ToSQL()
			if err != nil {
				b.Fatalf("%+v", err)
			}
			benchmarkGlobalVals = args
		}
	})
}

func BenchmarkInsertRecordsSQL(b *testing.B) {

	obj := someRecord{SomethingID: 1, UserID: 99, Other: false}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, args, err := dml.NewInsert("alpha").
			AddColumns("something_id", "user_id", "other").
			BindRecord(obj).
			ToSQL()
		if err != nil {
			b.Fatal(err)
		}
		benchmarkGlobalVals = args
		// ifaces = args.Interfaces()
	}
}

func BenchmarkRepeat(b *testing.B) {

	b.Run("multi", func(b *testing.B) {
		args := dml.MakeArgs(3).Ints(5, 7, 9, 11).Strings("a", "b", "c", "d", "e").Int(22)
		const want = "SELECT * FROM `table` WHERE id IN (?,?,?,?) AND name IN (?,?,?,?,?) AND status IN (?)"
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			s, err := dml.Repeat("SELECT * FROM `table` WHERE id IN (?) AND name IN (?) AND status IN (?)", args)
			if err != nil {
				b.Fatalf("%+v", err)
			}
			if s != want {
				b.Fatalf("\nHave: %q\nWant: %q", s, want)
			}
		}
	})

	b.Run("single", func(b *testing.B) {
		args := dml.MakeArgs(1).Ints(9, 8, 7, 6)
		const want = "SELECT * FROM `table` WHERE id IN (?,?,?,?)"
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			s, err := dml.Repeat("SELECT * FROM `table` WHERE id IN (?)", args)
			if err != nil {
				b.Fatalf("%+v", err)
			}
			if s != want {
				b.Fatalf("\nHave: %q\nWant: %q", s, want)
			}
		}
	})
}

func BenchmarkQuoteAs(b *testing.B) {
	const want = "`e`.`entity_id` AS `ee`"
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		if have := dml.Quoter.NameAlias("e.entity_id", "ee"); have != want {
			b.Fatalf("Have %s\nWant %s\n", have, want)
		}
	}
}

func BenchmarkQuoteQuote(b *testing.B) {
	const want = "`databaseName`.`tableName`"

	b.ReportAllocs()
	b.ResetTimer()
	b.Run("Worse Case", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			if have := dml.Quoter.QualifierName("database`Name", "table`Name"); have != want {
				b.Fatalf("Have %s\nWant %s\n", have, want)
			}
		}
	})
	b.Run("Best Case", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			if have := dml.Quoter.QualifierName("databaseName", "tableName"); have != want {
				b.Fatalf("Have %s\nWant %s\n", have, want)
			}
		}
	})
}

var benchmarkIfNull *dml.Condition

func BenchmarkIfNull(b *testing.B) {
	runner := func(want string, have ...string) func(*testing.B) {
		return func(b *testing.B) {
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				var alias string
				if lh := len(have); lh%2 == 1 && lh > 1 {
					alias = have[lh-1]
					have = have[:lh-1]
				}
				ifn := dml.SQLIfNull(have...)
				if alias != "" {
					ifn = ifn.Alias(alias)
				}
				benchmarkIfNull = ifn
			}
		}
	}
	b.Run("3 args expression right", runner(
		"IFNULL(`c2`,(1/0)) AS `alias`",
		"c2", "1/0", "alias",
	))
	b.Run("3 args no qualifier", runner(
		"IFNULL(`c1`,`c2`) AS `alias`",
		"c1", "c2", "alias",
	))
	b.Run("3 args with qualifier", runner(
		"IFNULL(`t1`.`c1`,`t2`.`c2`) AS `alias`",
		"t1.c1", "t2.c2", "alias",
	))
	b.Run("4 args", runner(
		"IFNULL(`t1`.`c1`,`t2`.`c2`)",
		"t1", "c1", "t2", "c2",
	))
	b.Run("5 args", runner(
		"IFNULL(`t1`.`c1`,`t2`.`c2`) AS `alias`",
		"t1", "c1", "t2", "c2", "ali`as",
	))
}

func BenchmarkUnion(b *testing.B) {

	newUnion5 := func() *dml.Union {
		// not valid SQL
		return dml.NewUnion(
			dml.NewSelect().AddColumns("t.value", "t.attribute_id").AddColumnsAliases("'varchar'", "col_type").FromAlias("catalog_product_entity_varchar", "t").
				Where(dml.Column("entity_id").Int64(1561), dml.Column("store_id").In().Int64s(1, 0)).
				OrderByDesc("t.varchar_store_id"),
			dml.NewSelect().AddColumns("t.value", "t.attribute_id").AddColumnsAliases("'int'", "col_type").FromAlias("catalog_product_entity_int", "t").
				Where(dml.Column("entity_id").Int64(1561), dml.Column("store_id").In().Int64s(1, 0)).
				OrderByDesc("t.int_store_id"),
			dml.NewSelect().AddColumns("t.value", "t.attribute_id").AddColumnsAliases("'decimal'", "col_type").FromAlias("catalog_product_entity_decimal", "t").
				Where(dml.Column("entity_id").Int64(1561), dml.Column("store_id").In().Int64s(1, 0)).
				OrderByDesc("t.decimal_store_id"),
			dml.NewSelect().AddColumns("t.value", "t.attribute_id").AddColumnsAliases("'datetime'", "col_type").FromAlias("catalog_product_entity_datetime", "t").
				Where(dml.Column("entity_id").Int64(1561), dml.Column("store_id").In().Int64s(1, 0)).
				OrderByDesc("t.datetime_store_id"),
			dml.NewSelect().AddColumns("t.value", "t.attribute_id").AddColumnsAliases("'text'", "col_type").FromAlias("catalog_product_entity_text", "t").
				Where(dml.Column("entity_id").Int64(1561), dml.Column("store_id").In().Int64s(1, 0)).
				OrderByDesc("t.text_store_id"),
		).
			All().
			OrderBy("a").
			OrderByDesc("b").PreserveResultSet()
	}

	newUnionTpl := func() *dml.Union {
		return dml.NewUnion(
			dml.NewSelect().AddColumns("t.value", "t.attribute_id").AddColumnsAliases("'{column}'", "col_type").FromAlias("catalog_product_entity_{type}", "t").
				Where(dml.Column("entity_id").Int64(1561), dml.Column("store_id").In().Int64s(1, 0)).
				OrderByDesc("t.{column}_store_id"),
		).
			StringReplace("{type}", "varchar", "int", "decimal", "datetime", "text").
			StringReplace("{column}", "varcharX", "intX", "decimalX", "datetimeX", "textX").
			PreserveResultSet().
			All().
			OrderByDesc("col_type")
	}

	b.Run("5 SELECTs", func(b *testing.B) {
		u := newUnion5()
		for i := 0; i < b.N; i++ {
			_, args, err := u.ToSQL()
			if err != nil {
				b.Fatalf("%+v", err)
			}
			benchmarkGlobalVals = args
		}
	})
	b.Run("5 SELECTs interpolated", func(b *testing.B) {
		u := newUnion5()
		u.Interpolate()
		for i := 0; i < b.N; i++ {
			_, args, err := u.ToSQL()
			if err != nil {
				b.Fatalf("%+v", err)
			}
			benchmarkGlobalVals = args
		}
	})
	b.Run("Template", func(b *testing.B) {
		u := newUnionTpl()
		for i := 0; i < b.N; i++ {
			_, args, err := u.ToSQL()
			if err != nil {
				b.Fatalf("%+v", err)
			}
			benchmarkGlobalVals = args
		}
	})
	b.Run("Template interpolated", func(b *testing.B) {
		u := newUnionTpl().Interpolate()
		for i := 0; i < b.N; i++ {
			_, args, err := u.ToSQL()
			if err != nil {
				b.Fatalf("%+v", err)
			}
			benchmarkGlobalVals = args
		}
	})
}

func BenchmarkUpdateValuesSQL(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, args, err := dml.NewUpdate("alpha").
			Set(
				dml.Column("something_id").Int64(1),
			).Where(
			dml.Column("id").Int64(1),
		).ToSQL()
		if err != nil {
			b.Fatalf("%+v", err)
		}
		benchmarkGlobalVals = args
	}
}

func BenchmarkArgUnion(b *testing.B) {
	reflectIFaceContainer := make([]interface{}, 0, 25)
	var finalArgs = make([]interface{}, 0, 30)
	drvVal := []driver.Valuer{dml.MakeNullString("I'm a valid null string: See the License for the specific language governing permissions and See the License for the specific language governing permissions and See the License for the specific language governing permissions and")}
	argUnion := dml.MakeArgs(30)
	now1 := dml.Now.UTC()
	b.ResetTimer()

	b.Run("reflection all types", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			reflectIFaceContainer = append(reflectIFaceContainer,
				int64(5), []int64{6, 7, 8},
				uint64(9), []uint64{10, 11, 12},
				float64(3.14159), []float64{33.44, 55.66, 77.88},
				true, []bool{true, false, true},
				`Licensed under the Apache License, Version 2.0 (the "License");`,
				[]string{`Unless required by applicable law or agreed to in writing, software`, `Licensed under the Apache License, Version 2.0 (the "License");`},
				drvVal[0],
				nil,
				now1,
			)
			var err error
			finalArgs, err = encodePlaceholder(finalArgs, reflectIFaceContainer)
			if err != nil {
				b.Fatal(err)
			}
			//b.Fatal("%#v", finalArgs)
			reflectIFaceContainer = reflectIFaceContainer[:0]
			finalArgs = finalArgs[:0]
		}
	})
	b.Run("args all types", func(b *testing.B) {
		// two times faster than the reflection version

		finalArgs = finalArgs[:0]

		for i := 0; i < b.N; i++ {
			argUnion = argUnion.
				Int64(5).Int64s(6, 7, 8).
				Uint64(9).Uint64s(10, 11, 12).
				Float64(3.14159).Float64s(33.44, 55.66, 77.88).
				Bool(true).Bools(true, false, true).
				String(`Licensed under the Apache License, Version 2.0 (the "License");`).
				Strings(`Unless required by applicable law or agreed to in writing, software`, `Licensed under the Apache License, Version 2.0 (the "License");`).
				DriverValue(drvVal...).
				Null().
				Time(now1)

			finalArgs = argUnion.Interfaces(finalArgs...)
			//b.Fatal("%#v", finalArgs)
			argUnion = argUnion[:0]
			finalArgs = finalArgs[:0]
		}
	})

	b.Run("reflection numbers", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			reflectIFaceContainer = append(reflectIFaceContainer,
				int64(5), []int64{6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19},
				uint64(9), []uint64{10, 11, 12, 20, 21, 22, 23, 24, 25, 26, 27, 28, 29},
				float64(3.14159), []float64{33.44, 55.66, 77.88, 11.22, math.Pi, math.E, math.Sqrt2},
				nil,
			)
			var err error
			finalArgs, err = encodePlaceholder(finalArgs, reflectIFaceContainer)
			if err != nil {
				b.Fatal(err)
			}
			//b.Fatal("%#v", finalArgs)
			reflectIFaceContainer = reflectIFaceContainer[:0]
			finalArgs = finalArgs[:0]
		}
	})
	b.Run("args numbers", func(b *testing.B) {
		// three times faster than the reflection version

		finalArgs = finalArgs[:0]
		for i := 0; i < b.N; i++ {
			argUnion = argUnion.
				Int64(5).Int64s(6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19).
				Uint64(9).Uint64s(10, 11, 12, 20, 21, 22, 23, 24, 25, 26, 27, 28, 29).
				Float64(3.14159).Float64s(33.44, 55.66, 77.88, 11.22, math.Pi, math.E, math.Sqrt2).
				Null()

			finalArgs = argUnion.Interfaces(finalArgs...)
			//b.Fatal("%#v", finalArgs)
			argUnion = argUnion[:0]
			finalArgs = finalArgs[:0]
		}
	})

}

func encodePlaceholder(args []interface{}, value interface{}) ([]interface{}, error) {

	if valuer, ok := value.(driver.Valuer); ok {
		// get driver.Valuer's data
		var err error
		value, err = valuer.Value()
		if err != nil {
			return args, err
		}
	}

	if value == nil {
		return append(args, nil), nil
	}
	v := reflect.ValueOf(value)
	switch v.Kind() {
	case reflect.String:
		return append(args, v.String()), nil
	case reflect.Bool:
		return append(args, v.Bool()), nil
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return append(args, v.Int()), nil
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return append(args, v.Uint()), nil
	case reflect.Float32, reflect.Float64:
		return append(args, v.Float()), nil
	case reflect.Struct:
		if v.Type() == reflect.TypeOf(time.Time{}) {
			return append(args, v.Interface().(time.Time)), nil
		}
	case reflect.Slice:
		if v.Type().Elem().Kind() == reflect.Uint8 {
			// []byte
			return append(args, v.Bytes()), nil
		}
		if v.Len() == 0 {
			// FIXME: support zero-length slice
			return args, errors.NewNotValidf("invalid slice length")
		}

		for n := 0; n < v.Len(); n++ {
			var err error
			// recursion
			args, err = encodePlaceholder(args, v.Index(n).Interface())
			if err != nil {
				return args, err
			}
		}
		return args, nil
	case reflect.Ptr:
		if v.IsNil() {
			return append(args, nil), nil
		}
		return encodePlaceholder(args, v.Elem().Interface())

	}
	return args, errors.NewNotSupportedf("Type %#v not supported", value)
}

var benchmarkDecimal_String string

func BenchmarkDecimal_String(b *testing.B) {

	b.Run("123456789", func(b *testing.B) {
		d := dml.Decimal{
			Precision: 123456789,
			Valid:     true,
		}
		for i := 0; i < b.N; i++ {
			benchmarkDecimal_String = d.String()
		}
	})
	b.Run("-123456789", func(b *testing.B) {
		d := dml.Decimal{
			Precision: 123456789,
			Valid:     true,
			Negative:  true,
		}
		for i := 0; i < b.N; i++ {
			benchmarkDecimal_String = d.String()
		}
	})
	b.Run("12345.6789", func(b *testing.B) {
		d := dml.Decimal{
			Precision: 123456789,
			Scale:     4,
			Valid:     true,
		}
		for i := 0; i < b.N; i++ {
			benchmarkDecimal_String = d.String()
		}
	})
	b.Run("-12345.6789", func(b *testing.B) {
		d := dml.Decimal{
			Precision: 123456789,
			Valid:     true,
			Scale:     4,
			Negative:  true,
		}
		for i := 0; i < b.N; i++ {
			benchmarkDecimal_String = d.String()
		}
	})
	b.Run("-Scale140", func(b *testing.B) {
		d := dml.Decimal{
			Valid:     true,
			Precision: math.MaxUint64,
			Scale:     140,
			Negative:  true,
		}
		for i := 0; i < b.N; i++ {
			benchmarkDecimal_String = d.String()
		}
	})
}

var benchmarkDecimal_MarshalBinary []byte

func BenchmarkDecimal_Binary(b *testing.B) {

	b.Run("Marshal", func(b *testing.B) {
		d := dml.Decimal{
			Precision: 123456789,
			Valid:     true,
			Scale:     4,
			Negative:  true,
		}
		for i := 0; i < b.N; i++ {
			var err error
			benchmarkDecimal_MarshalBinary, err = d.MarshalBinary()
			if err != nil {
				b.Fatal(err)
			}
		}
	})

	var dUn dml.Decimal
	b.Run("Unmarshal", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			if err := dUn.UnmarshalBinary(benchmarkDecimal_MarshalBinary); err != nil {
				b.Fatal(err)
			}
		}
	})
}
