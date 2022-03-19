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
	"context"
	"database/sql"
	"testing"
	"time"
)

var (
	benchmarkGlobalVals []any
	benchmarkSelectStr  string
)

// BenchmarkSelect_Rows-4		 1000000	      2276 ns/op	    4344 B/op	      12 allocs/op git commit 609db6db
// BenchmarkSelect_Rows-4   	  500000	      2919 ns/op	    5411 B/op	      18 allocs/op
// BenchmarkSelect_Rows-4   	  500000	      2504 ns/op	    4239 B/op	      17 allocs/op
func BenchmarkSelect_Rows(b *testing.B) {
	tables := []string{"eav_attribute"}
	ctx := context.TODO()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {

		sel := NewSelect("TABLE_NAME", "COLUMN_NAME", "ORDINAL_POSITION", "COLUMN_DEFAULT", "IS_NULLABLE",
			"DATA_TYPE", "CHARACTER_MAXIMUM_LENGTH", "NUMERIC_PRECISION", "NUMERIC_SCALE",
			"COLUMN_TYPE", "COLUMN_KEY", "EXTRA", "COLUMN_COMMENT").From("information_schema.COLUMNS").
			Where(Expr(`TABLE_SCHEMA=DATABASE()`))

		sel.Where(Column("TABLE_NAME").In().PlaceHolder())

		rows, err := sel.WithDBR(dbMock{}).QueryContext(ctx, tables)
		if err != nil {
			b.Fatalf("%+v", err)
		}
		if rows == nil {
			b.Fatal("Query should not be nil")
		}
	}
}

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
		var err error
		benchmarkSelectStr, benchmarkGlobalVals, err = NewSelect("something_id", "user_id", "other").
			From("some_table").
			Where(
				Expr("d = ? OR e = ?").Int64(1).Str("wat"),
				Column("a").In().Int64s(aVal...),
			).
			OrderByDesc("id").
			Paginate(1, 20).
			ToSQL()
		if err != nil {
			b.Fatalf("%+v", err)
		}
	}
}

// Type Condition is not a pointer. The GC pressure is much less but the program
// is a bit slower due to copying.
// BenchmarkSelectExcessConditions-4   	  200000	      7055 ns/op	    4984 B/op	      26 allocs/op
// BenchmarkSelectExcessConditions-4   	  200000	      7118 ns/op	    4984 B/op	      26 allocs/op
// The next line gives the results for the implementation when using type
// Condition as a pointer. We have more allocs, more pressure on the GC but it's
// a bit faster.
// BenchmarkSelectExcessConditions-4   	  200000	      6297 ns/op	    4920 B/op	      35 allocs/op
// For now we stick with the pointers. Reason: the SQL statements are getting
// usually only created once and not in a loop.
func BenchmarkSelectExcessConditions(b *testing.B) {
	i64Vals := []int64{1, 2, 3}
	b.ResetTimer()
	var err error
	for i := 0; i < b.N; i++ {
		benchmarkSelectStr, benchmarkGlobalVals, err = NewSelect("entity_id", "name", "q.total_sum").
			FromAlias("customer", "c").
			Join(MakeIdentifier("quote").Alias("q"),
				Column("c.entity_id").Column("q.customer_id"),
				Column("c.group_id").Column("q.group_id"),
				Column("c.shop_id").NotEqual().PlaceHolders(10),
			).
			Where(
				Expr("d = ? OR e = ?").Int64(1).Str("wat"),
				Column("a").In().PlaceHolder(),
				Column("q.product_id").In().Int64s(i64Vals...),
				Column("q.sub_total").NotBetween().Float64s(3.141, 6.2831),
				Column("q.qty").GreaterOrEqual().PlaceHolder(),
				Column("q.coupon_code").SpaceShip().Str("400EA8BBE4"),
			).
			OrderByDesc("id").
			Paginate(1, 20).
			ToSQL()
		if err != nil {
			b.Fatalf("%+v", err)
		}
	}
}

func BenchmarkSelectFullSQL(b *testing.B) {
	sqlObj := NewSelect("a", "b", "z", "y", "x").From("c").
		Distinct().
		Where(
			Expr("`d` = ? OR `e` = ?").Int64(1).Str("wat"),
			Column("f").Int64(2),
			Column("x").Str("hi"),
			Column("g").Int64(3),
			Column("h").In().Ints(1, 2, 3),
		).
		GroupBy("ab").GroupBy("ii").GroupBy("iii").
		Having(Expr("j = k"), Column("jj").Int64(1)).
		Having(Column("jjj").Int64(2)).
		OrderBy("l1").OrderBy("l2").OrderBy("l3").
		Limit(8, 7)
	sqlObjDBR := sqlObj.WithDBR(nil)
	b.ResetTimer()
	b.ReportAllocs()

	// BenchmarkSelectFullSQL/NewSelect-4             300000	      5849 ns/op	    3649 B/op	      38 allocs/op
	// BenchmarkSelectFullSQL/NewSelect-4          	  200000	      6307 ns/op	    4922 B/op	      45 allocs/op <== builder structs
	// BenchmarkSelectFullSQL/NewSelect-4         	  200000	      7084 ns/op	    8212 B/op	      44 allocs/op <== DBR
	// BenchmarkSelectFullSQL/NewSelect-4         	  200000	      6449 ns/op	    5515 B/op	      44 allocs/op no pointers
	// BenchmarkSelectFullSQL/NewSelect-4         	  200000	      6268 ns/op	    5443 B/op	      37 allocs/op
	// BenchmarkSelectFullSQL/NewSelect-4          	  342667	      3497 ns/op	    3672 B/op	      29 allocs/op
	b.Run("NewSelect", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			var err error
			benchmarkSelectStr, benchmarkGlobalVals, err = NewSelect("a", "b", "z", "y", "x").From("c").
				Distinct().
				Where(
					Expr("`d` = ? OR `e` = ?").Int64(1).Str("wat"),
					Column("f").Int64(2),
					Column("x").Str("hi"),
					Column("g").Int64(3),
					Column("h").In().Ints(1, 2, 3),
				).
				GroupBy("ab").GroupBy("ii").GroupBy("iii").
				Having(Expr("j = k"), Column("jj").Int64(1)).
				Having(Column("jjj").Int64(2)).
				OrderBy("l1").OrderBy("l2").OrderBy("l3").
				Limit(8, 7).
				ToSQL()
			if err != nil {
				b.Fatalf("%+v", err)
			}
		}
	})

	b.Run("ToSQL Interpolate Cache", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			var err error
			benchmarkSelectStr, benchmarkGlobalVals, err = sqlObjDBR.ToSQL()
			if err != nil {
				b.Fatalf("%+v", err)
			}
		}
	})
}

func BenchmarkSelect_Large_IN(b *testing.B) {
	// This tests simulates selecting many varchar attribute values for specific products.
	entityIDs := make([]int64, 1024)
	for i := 0; i < 1024; i++ {
		entityIDs[i] = int64(i + 1600)
	}

	b.ResetTimer()
	b.ReportAllocs()

	b.Run("SQL", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			sel := NewSelect("entity_id", "attribute_id", "value").
				From("catalog_product_entity_varchar").
				Where(Column("entity_type_id").Int64(4)).
				Where(Column("entity_id").In().Int64s(entityIDs...)).
				Where(Column("attribute_id").In().Int64s(174, 175)).
				Where(Column("store_id").Int(0))
			var err error
			benchmarkSelectStr, benchmarkGlobalVals, err = sel.ToSQL()
			if err != nil {
				b.Fatalf("%+v", err)
			}
			if benchmarkGlobalVals != nil {
				b.Fatal("Args should be nil")
			}
		}
	})

	b.Run("interpolate", func(b *testing.B) {
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			sel := NewSelect("entity_id", "attribute_id", "value").
				From("catalog_product_entity_varchar").
				Where(Column("entity_type_id").PlaceHolder()).
				Where(Column("entity_id").In().PlaceHolder()).
				Where(Column("attribute_id").In().PlaceHolder()).
				Where(Column("store_id").PlaceHolder())

			var err error
			benchmarkSelectStr, benchmarkGlobalVals, err = sel.WithDBR(dbMock{}).
				Interpolate().
				testWithArgs(4, entityIDs, []int64{174, 175}, 0).
				ToSQL()
			if err != nil {
				b.Fatalf("%+v", err)
			}
			if benchmarkGlobalVals != nil {
				b.Fatal("Args should be nil")
			}
		}
	})

	b.Run("interpolate named", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			sel := NewSelect("entity_id", "attribute_id", "value").
				From("catalog_product_entity_varchar").
				Where(Column("entity_type_id").NamedArg("EntityTypeId")).
				Where(Column("entity_id").In().NamedArg("EntityId")).
				Where(Column("attribute_id").In().NamedArg("AttributeId")).
				Where(Column("store_id").NamedArg("StoreId"))

			var err error
			benchmarkSelectStr, benchmarkGlobalVals, err = sel.WithDBR(dbMock{}).Interpolate().testWithArgs(
				sql.Named("EntityTypeId", int64(4)),
				sql.Named("EntityId", entityIDs),
				sql.Named("AttributeId", []int64{174, 175}),
				sql.Named("StoreId", 0),
			).ToSQL()
			if err != nil {
				b.Fatalf("%+v", err)
			}
			if benchmarkGlobalVals != nil {
				b.Fatal("Args should be nil")
			}
		}
	})

	b.Run("interpolate optimized", func(b *testing.B) {
		sel := NewSelect("entity_id", "attribute_id", "value").
			From("catalog_product_entity_varchar").
			Where(Column("entity_type_id").PlaceHolder()).
			Where(Column("entity_id").In().PlaceHolder()).
			Where(Column("attribute_id").In().PlaceHolder()).
			Where(Column("store_id").PlaceHolder())

		selA := sel.WithDBR(dbMock{}).Interpolate()

		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			var err error
			// the generated string benchmarkSelectStr is 5300 characters long
			benchmarkSelectStr, benchmarkGlobalVals, err = selA.testWithArgs(4, entityIDs, []int64{174, 175}, 0).ToSQL()
			if err != nil {
				b.Fatalf("%+v", err)
			}
			if benchmarkGlobalVals != nil {
				b.Fatal("Args should be nil")
			}
			selA.Reset()
		}
	})
}

func BenchmarkSelect_ComplexAddColumns(b *testing.B) {
	var haveSQL string
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		var args []any
		var err error
		haveSQL, args, err = NewSelect().
			AddColumns(" entity_id ,   value").
			AddColumns("cpev.entity_type_id", "cpev.attribute_id").
			AddColumnsAliases("(cpev.id*3)", "weirdID").
			AddColumnsAliases("cpev.value", "value2nd").
			FromAlias("catalog_product_entity_varchar", "cpev").
			Where(Column("entity_type_id").Int64(4)).
			Where(Column("attribute_id").In().Int64s(174, 175)).
			Where(Column("store_id").Int64(0)).
			ToSQL()
		if err != nil {
			b.Fatalf("%+v", err)
		}
		benchmarkGlobalVals = args
	}
	_ = haveSQL
	// b.Logf("%s", haveSQL)
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
		haveSQL, benchmarkGlobalVals, err = NewSelect().
			AddColumns("price", "sku", "name").
			AddColumnsConditions(
				SQLCase("", "`closed`",
					"date_start <= ? AND date_end >= ?", "`open`",
					"date_start > ? AND date_end > ?", "`upcoming`",
				).Alias("is_on_sale").Time(start).Time(end).Time(start).Time(end),
			).
			From("catalog_promotions").
			Where(
				Column("promotion_id").
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
			_, benchmarkGlobalVals, err = NewDelete("alpha").Where(Column("a").Str("b")).Limit(1).OrderBy("id").ToSQL()
			if err != nil {
				b.Fatalf("%+v", err)
			}
		}
	})

	sqlObj := NewDelete("alpha").Where(Column("a").Str("b")).Limit(1).OrderBy("id")
	b.Run("ToSQL no cache", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			var err error
			benchmarkSelectStr, benchmarkGlobalVals, err = sqlObj.ToSQL()
			if err != nil {
				b.Fatalf("%+v", err)
			}
		}
	})

	sqlObjDBR := sqlObj.WithDBR(nil)
	b.Run("ToSQL with cache", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			var err error
			benchmarkSelectStr, benchmarkGlobalVals, err = sqlObjDBR.ToSQL()
			if err != nil {
				b.Fatalf("%+v", err)
			}
		}
	})
}

func BenchmarkInsertValuesSQL(b *testing.B) {
	b.Run("NewInsert", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			var err error
			benchmarkSelectStr, benchmarkGlobalVals, err = NewInsert("alpha").
				AddColumns("something_id", "user_id", "other").
				WithDBR(dbMock{}).testWithArgs(1, 2, true).ToSQL()
			if err != nil {
				b.Fatal(err)
			}
		}
	})

	b.Run("ToSQL no cache", func(b *testing.B) {
		sqlObj := NewInsert("alpha").AddColumns("something_id", "user_id", "other")
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			var err error
			sqlObjA := sqlObj.WithDBR(dbMock{}).testWithArgs(1, 2, true)
			benchmarkSelectStr, benchmarkGlobalVals, err = sqlObjA.ToSQL()
			if err != nil {
				b.Fatalf("%+v", err)
			}
		}
	})

	b.Run("ToSQL with cache", func(b *testing.B) {
		sqlObj := NewInsert("alpha").AddColumns("something_id", "user_id", "other")
		delA := sqlObj.WithDBR(dbMock{})
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			var err error
			benchmarkSelectStr, benchmarkGlobalVals, err = delA.testWithArgs(1, 2, true).ToSQL()
			if err != nil {
				b.Fatalf("%+v", err)
			}
			delA.Reset()
		}
	})
}

func BenchmarkInsertRecordsSQL(b *testing.B) {
	obj := someRecord{SomethingID: 1, UserID: 99, Other: false}
	insA := NewInsert("alpha").
		AddColumns("something_id", "user_id", "other").
		WithDBR(dbMock{})

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		var err error
		benchmarkSelectStr, benchmarkGlobalVals, err = insA.testWithArgs(Qualify("", obj)).ToSQL()
		if err != nil {
			b.Fatal(err)
		}
		insA.Reset()
	}
}

func BenchmarkRepeat(b *testing.B) {
	cp, err := NewConnPool()
	if err != nil {
		b.Fatal(err)
	}
	b.Run("multi", func(b *testing.B) {
		const want = "SELECT * FROM `table` WHERE id IN (?,?,?,?) AND name IN (?,?,?,?,?) AND status = ?"
		dbr := cp.WithQueryBuilder(QuerySQL("SELECT * FROM `table` WHERE id IN ? AND name IN ? AND status = ?")).
			ExpandPlaceHolders().
			testWithArgs([]int{5, 7, 9, 11}, []string{"a", "b", "c", "d", "e"}, 22)
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			s, _, err := dbr.ToSQL()
			if err != nil {
				b.Fatalf("%+v", err)
			}
			if s != want {
				b.Fatalf("\nHave: %q\nWant: %q", s, want)
			}
		}
	})

	b.Run("single", func(b *testing.B) {
		const want = "SELECT * FROM `table` WHERE id IN (?,?,?,?)"
		dbr := cp.WithQueryBuilder(QuerySQL("SELECT * FROM `table` WHERE id IN ?")).
			ExpandPlaceHolders().
			testWithArgs([]int{9, 8, 7, 6})
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			s, _, err := dbr.ToSQL()
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
		if have := Quoter.NameAlias("e.entity_id", "ee"); have != want {
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
			if have := Quoter.QualifierName("database`Name", "table`Name"); have != want {
				b.Fatalf("Have %s\nWant %s\n", have, want)
			}
		}
	})
	b.Run("Best Case", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			if have := Quoter.QualifierName("databaseName", "tableName"); have != want {
				b.Fatalf("Have %s\nWant %s\n", have, want)
			}
		}
	})
}

var benchmarkIfNull *Condition

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
				ifn := SQLIfNull(have...)
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
	newUnion5 := func() *Union {
		// not valid SQL
		return NewUnion(
			NewSelect().AddColumns("t.value", "t.attribute_id").AddColumnsAliases("'varchar'", "col_type").FromAlias("catalog_product_entity_varchar", "t").
				Where(Column("entity_id").Int64(1561), Column("store_id").In().Int64s(1, 0)).
				OrderByDesc("t.varchar_store_id"),
			NewSelect().AddColumns("t.value", "t.attribute_id").AddColumnsAliases("'int'", "col_type").FromAlias("catalog_product_entity_int", "t").
				Where(Column("entity_id").Int64(1561), Column("store_id").In().Int64s(1, 0)).
				OrderByDesc("t.int_store_id"),
			NewSelect().AddColumns("t.value", "t.attribute_id").AddColumnsAliases("'decimal'", "col_type").FromAlias("catalog_product_entity_decimal", "t").
				Where(Column("entity_id").Int64(1561), Column("store_id").In().Int64s(1, 0)).
				OrderByDesc("t.decimal_store_id"),
			NewSelect().AddColumns("t.value", "t.attribute_id").AddColumnsAliases("'datetime'", "col_type").FromAlias("catalog_product_entity_datetime", "t").
				Where(Column("entity_id").Int64(1561), Column("store_id").In().Int64s(1, 0)).
				OrderByDesc("t.datetime_store_id"),
			NewSelect().AddColumns("t.value", "t.attribute_id").AddColumnsAliases("'text'", "col_type").FromAlias("catalog_product_entity_text", "t").
				Where(Column("entity_id").Int64(1561), Column("store_id").In().Int64s(1, 0)).
				OrderByDesc("t.text_store_id"),
		).
			All().
			OrderBy("a").
			OrderByDesc("b").PreserveResultSet()
	}

	newUnionTpl := func() *Union {
		return NewUnion(
			NewSelect().AddColumns("t.value", "t.attribute_id").AddColumnsAliases("'{column}'", "col_type").FromAlias("catalog_product_entity_{type}", "t").
				Where(Column("entity_id").Int64(1561), Column("store_id").In().Int64s(1, 0)).
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
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			var err error
			benchmarkSelectStr, benchmarkGlobalVals, err = u.ToSQL()
			if err != nil {
				b.Fatalf("%+v", err)
			}
		}
	})
	b.Run("5 SELECTs not cached", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			var err error
			benchmarkSelectStr, benchmarkGlobalVals, err = newUnion5().ToSQL()
			if err != nil {
				b.Fatalf("%+v", err)
			}
		}
	})
	b.Run("5 SELECTs WithDBR", func(b *testing.B) {
		u := newUnion5().WithDBR(dbMock{})
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			var err error
			benchmarkSelectStr, benchmarkGlobalVals, err = u.ToSQL()
			if err != nil {
				b.Fatalf("%+v", err)
			}
		}
	})
	b.Run("Template", func(b *testing.B) {
		u := newUnionTpl()
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			var err error
			benchmarkSelectStr, benchmarkGlobalVals, err = u.ToSQL()
			if err != nil {
				b.Fatalf("%+v", err)
			}
		}
	})
	b.Run("Template not cached", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			var err error
			benchmarkSelectStr, benchmarkGlobalVals, err = newUnionTpl().ToSQL()
			if err != nil {
				b.Fatalf("%+v", err)
			}
		}
	})
	b.Run("Template interpolated", func(b *testing.B) {
		u := newUnionTpl().WithDBR(nil)
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			var err error
			benchmarkSelectStr, benchmarkGlobalVals, err = u.ToSQL()
			if err != nil {
				b.Fatalf("%+v", err)
			}
		}
	})
}

func BenchmarkUpdateValuesSQL(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		var err error
		benchmarkSelectStr, benchmarkGlobalVals, err = NewUpdate("alpha").
			AddClauses(
				Column("something_id").Int64(1),
			).Where(
			Column("id").Int64(1),
		).ToSQL()
		if err != nil {
			b.Fatalf("%+v", err)
		}
	}
}
