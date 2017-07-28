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

package dbr

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestUnionStmts(t *testing.T) {
	t.Parallel()

	t.Run("simple", func(t *testing.T) {
		u := NewUnion(
			NewSelect("a", "b").From("tableAB").Where(Column("a").Int64(3)),
		)
		u.Append(
			NewSelect("c", "d").From("tableCD").Where(Column("d").String("e")),
		)
		compareToSQL(t, u, nil,
			"(SELECT `a`, `b` FROM `tableAB` WHERE (`a` = ?))\nUNION\n(SELECT `c`, `d` FROM `tableCD` WHERE (`d` = ?))",
			"(SELECT `a`, `b` FROM `tableAB` WHERE (`a` = 3))\nUNION\n(SELECT `c`, `d` FROM `tableCD` WHERE (`d` = 'e'))",
			int64(3), "e",
		)
	})

	t.Run("simple all", func(t *testing.T) {
		u := NewUnion(
			NewSelect("c", "d").From("tableCD").Where(Column("d").String("e")),
			NewSelect("a", "b").From("tableAB").Where(Column("a").Int64(3)),
		).All()
		compareToSQL(t, u, nil,
			"(SELECT `c`, `d` FROM `tableCD` WHERE (`d` = ?))\nUNION ALL\n(SELECT `a`, `b` FROM `tableAB` WHERE (`a` = ?))",
			"(SELECT `c`, `d` FROM `tableCD` WHERE (`d` = 'e'))\nUNION ALL\n(SELECT `a`, `b` FROM `tableAB` WHERE (`a` = 3))",
			"e", int64(3),
		)
	})

	t.Run("order by", func(t *testing.T) {
		u := NewUnion(
			NewSelect("a").AddColumnsAlias("d", "b").From("tableAD").Where(Column("d").String("f")),
			NewSelect("a", "b").From("tableAB").Where(Column("a").Int64(3)),
		).All().OrderBy("a").OrderByDesc("b")

		compareToSQL(t, u, nil,
			"(SELECT `a`, `d` AS `b` FROM `tableAD` WHERE (`d` = ?))\nUNION ALL\n(SELECT `a`, `b` FROM `tableAB` WHERE (`a` = ?))\nORDER BY `a` ASC, `b` DESC",
			"(SELECT `a`, `d` AS `b` FROM `tableAD` WHERE (`d` = 'f'))\nUNION ALL\n(SELECT `a`, `b` FROM `tableAB` WHERE (`a` = 3))\nORDER BY `a` ASC, `b` DESC",
			"f", int64(3),
		)
	})

	t.Run("preserve result set", func(t *testing.T) {
		u := NewUnion(
			NewSelect("a").AddColumnsAlias("d", "b").From("tableAD"),
			NewSelect("a", "b").From("tableAB").Where(Column("c").Between().Int64s(3, 5)),
		).All().OrderBy("a").OrderByDesc("b").PreserveResultSet()

		// testing idempotent function ToSQL
		for i := 0; i < 3; i++ {
			compareToSQL(t, u, nil,
				"(SELECT `a`, `d` AS `b`, 0 AS `_preserve_result_set` FROM `tableAD`)\nUNION ALL\n(SELECT `a`, `b`, 1 AS `_preserve_result_set` FROM `tableAB` WHERE (`c` BETWEEN ? AND ?))\nORDER BY `_preserve_result_set`, `a` ASC, `b` DESC",
				"(SELECT `a`, `d` AS `b`, 0 AS `_preserve_result_set` FROM `tableAD`)\nUNION ALL\n(SELECT `a`, `b`, 1 AS `_preserve_result_set` FROM `tableAB` WHERE (`c` BETWEEN 3 AND 5))\nORDER BY `_preserve_result_set`, `a` ASC, `b` DESC",
				int64(3), int64(5),
			)
		}
	})
	t.Run("intersec", func(t *testing.T) {
		u := NewUnion(
			NewSelect("a").From("tableAD"),
			NewSelect("b").From("tableAB"),
		).All().Intersect().OrderBy("a").OrderByDesc("b")
		// All gets ignored

		compareToSQL(t, u, nil,
			"(SELECT `a` FROM `tableAD`)\nINTERSECT\n(SELECT `b` FROM `tableAB`)\nORDER BY `a` ASC, `b` DESC",
			"(SELECT `a` FROM `tableAD`)\nINTERSECT\n(SELECT `b` FROM `tableAB`)\nORDER BY `a` ASC, `b` DESC",
		)
	})
	t.Run("except", func(t *testing.T) {
		u := NewUnion(
			NewSelect("a").From("tableAD"),
			NewSelect("b").From("tableAB"),
		).All().Except()
		// All gets ignored

		compareToSQL(t, u, nil,
			"(SELECT `a` FROM `tableAD`)\nEXCEPT\n(SELECT `b` FROM `tableAB`)",
			"(SELECT `a` FROM `tableAD`)\nEXCEPT\n(SELECT `b` FROM `tableAB`)",
		)
	})
}

func TestUnion_UseBuildCache(t *testing.T) {
	t.Parallel()

	u := NewUnion(
		NewSelect("a").AddColumnsAlias("d", "b").From("tableAD"),
		NewSelect("a", "b").From("tableAB").Where(Column("b").Float64(3.14159)),
	).All().
		StringReplace("MyKey", "a", "b", "c"). // does nothing because more than one NewSelect functions
		OrderBy("a").OrderByDesc("b").OrderByExpr(`concat("c",b,"d")`).
		PreserveResultSet().BuildCache()

	const cachedSQLPlaceHolder = "(SELECT `a`, `d` AS `b`, 0 AS `_preserve_result_set` FROM `tableAD`)\nUNION ALL\n(SELECT `a`, `b`, 1 AS `_preserve_result_set` FROM `tableAB` WHERE (`b` = ?))\nORDER BY `_preserve_result_set`, `a` ASC, `b` DESC, concat(\"c\",b,\"d\")"
	t.Run("without interpolate", func(t *testing.T) {
		for i := 0; i < 3; i++ {
			compareToSQL(t, u, nil,
				cachedSQLPlaceHolder,
				"",
				float64(3.14159),
			)
			require.Equal(t, cachedSQLPlaceHolder, string(u.cacheSQL))
		}
	})

	t.Run("with interpolate", func(t *testing.T) {
		u.cacheSQL = nil

		const cachedSQLInterpolated = "(SELECT `a`, `d` AS `b`, 0 AS `_preserve_result_set` FROM `tableAD`)\nUNION ALL\n(SELECT `a`, `b`, 1 AS `_preserve_result_set` FROM `tableAB` WHERE (`b` = 3.14159))\nORDER BY `_preserve_result_set`, `a` ASC, `b` DESC, concat('c',b,'d')"
		for i := 0; i < 3; i++ {
			compareToSQL(t, u, nil,
				cachedSQLPlaceHolder,
				cachedSQLInterpolated,
				3.14159,
			)
			require.Equal(t, cachedSQLPlaceHolder, string(u.cacheSQL))
		}
	})
}

func TestNewUnionTemplate(t *testing.T) {
	t.Parallel()
	/*
		   (SELECT `t`.`value`,
				   `t`.`attribute_id`
			FROM   `catalog_product_entity_varchar` AS `t`
			WHERE  ( entity_id = '1561' )
				   AND ( `store_id` IN ( '1', 0 ) )
			ORDER  BY `t`.`store_id` DESC)
		   UNION ALL
		   (SELECT `t`.`value`,
				   `t`.`attribute_id`
			FROM   `catalog_product_entity_int` AS `t`
			WHERE  ( entity_id = '1561' )
				   AND ( `store_id` IN ( '1', 0 ) )
			ORDER  BY `t`.`store_id` DESC)
		   UNION ALL
		   (SELECT `t`.`value`,
				   `t`.`attribute_id`
			FROM   `catalog_product_entity_decimal` AS `t`
			WHERE  ( entity_id = '1561' )
				   AND ( `store_id` IN ( '1', 0 ) )
			ORDER  BY `t`.`store_id` DESC)
		   UNION ALL
		   (SELECT `t`.`value`,
				   `t`.`attribute_id`
			FROM   `catalog_product_entity_datetime` AS `t`
			WHERE  ( entity_id = '1561' )
				   AND ( `store_id` IN ( '1', 0 ) )
			ORDER  BY `t`.`store_id` DESC)
		   UNION ALL
		   (SELECT `t`.`value`,
				   `t`.`attribute_id`
			FROM   `catalog_product_entity_text` AS `t`
			WHERE  ( entity_id = '1561' )
				   AND ( `store_id` IN ( '1', 0 ) )
			ORDER  BY `t`.`store_id` DESC);
	*/

	t.Run("full statement EAV", func(t *testing.T) {
		u := NewUnion(
			NewSelect().AddColumns("t.value", "t.attribute_id").AddColumnsAlias("t.{column}", "col_type").
				FromAlias("catalog_product_entity_{type}", "t").
				Where(Column("entity_id").Int64(1561), Column("store_id").In().Int64s(1, 0)).
				OrderByDesc("t.{column}_store_id"),
		).
			StringReplace("{type}", "varchar", "int", "decimal", "datetime", "text").
			StringReplace("{column}", "varcharX", "intX", "decimalX", "datetimeX", "textX").
			PreserveResultSet().
			All().
			OrderByDesc("col_type")

		// testing idempotent function ToSQL
		for i := 0; i < 3; i++ {
			uStr, args, err := u.ToSQL()
			if err != nil {
				t.Fatalf("%+v", err)
			}

			wantArg := []interface{}{int64(1561), int64(1), int64(0)}
			haveArg := args
			assert.Exactly(t, wantArg, haveArg[:3])
			assert.Exactly(t, wantArg, haveArg[3:6])
			assert.Len(t, haveArg, 15)
			assert.Exactly(t,
				"(SELECT `t`.`value`, `t`.`attribute_id`, `t`.`varcharX` AS `col_type`, 0 AS `_preserve_result_set` FROM `catalog_product_entity_varchar` AS `t` WHERE (`entity_id` = ?) AND (`store_id` IN (?,?)) ORDER BY `t`.`varcharX_store_id` DESC)\nUNION ALL\n(SELECT `t`.`value`, `t`.`attribute_id`, `t`.`intX` AS `col_type`, 1 AS `_preserve_result_set` FROM `catalog_product_entity_int` AS `t` WHERE (`entity_id` = ?) AND (`store_id` IN (?,?)) ORDER BY `t`.`intX_store_id` DESC)\nUNION ALL\n(SELECT `t`.`value`, `t`.`attribute_id`, `t`.`decimalX` AS `col_type`, 2 AS `_preserve_result_set` FROM `catalog_product_entity_decimal` AS `t` WHERE (`entity_id` = ?) AND (`store_id` IN (?,?)) ORDER BY `t`.`decimalX_store_id` DESC)\nUNION ALL\n(SELECT `t`.`value`, `t`.`attribute_id`, `t`.`datetimeX` AS `col_type`, 3 AS `_preserve_result_set` FROM `catalog_product_entity_datetime` AS `t` WHERE (`entity_id` = ?) AND (`store_id` IN (?,?)) ORDER BY `t`.`datetimeX_store_id` DESC)\nUNION ALL\n(SELECT `t`.`value`, `t`.`attribute_id`, `t`.`textX` AS `col_type`, 4 AS `_preserve_result_set` FROM `catalog_product_entity_text` AS `t` WHERE (`entity_id` = ?) AND (`store_id` IN (?,?)) ORDER BY `t`.`textX_store_id` DESC)\nORDER BY `_preserve_result_set`, `col_type` DESC",
				uStr)
		}
	})
	t.Run("StringReplace 2nd call fewer values", func(t *testing.T) {
		defer func() {
			if r := recover(); r != nil {
				if err, ok := r.(error); ok {
					t.Log(err)
				} else {
					t.Errorf("Panic should contain an error but got:\n%+v", r)
				}
			} else {
				t.Error("Expecting a panic but got nothing")
			}
		}()

		NewUnion(
			NewSelect().AddColumns("t.value,t.attribute_id,t.{column} AS `col_type`").FromAlias("catalog_product_entity_{type}", "t"),
		).
			StringReplace("{type}", "varchar", "int", "decimal", "datetime", "text").
			StringReplace("{column}", "varcharX", "intX", "decimalX", "datetimeX")

	})
	t.Run("StringReplace 2nd call too many values and nothing should happen", func(t *testing.T) {
		u := NewUnion(
			NewSelect().AddColumns("t.value,t.attribute_id,t.{column} AS `col_type`").FromAlias("catalog_product_entity_{type}", "t"),
		).
			StringReplace("{type}", "varchar", "int", "decimal", "datetime", "text").
			StringReplace("{column}", "varcharX", "intX", "decimalX", "datetimeX", "textX", "bytesX")
		compareToSQL(t, u, nil,
			"(SELECT `t`.`value,t.attribute_id,t.varcharX AS col_type` FROM `catalog_product_entity_varchar` AS `t`)\nUNION\n(SELECT `t`.`value,t.attribute_id,t.intX AS col_type` FROM `catalog_product_entity_int` AS `t`)\nUNION\n(SELECT `t`.`value,t.attribute_id,t.decimalX AS col_type` FROM `catalog_product_entity_decimal` AS `t`)\nUNION\n(SELECT `t`.`value,t.attribute_id,t.datetimeX AS col_type` FROM `catalog_product_entity_datetime` AS `t`)\nUNION\n(SELECT `t`.`value,t.attribute_id,t.textX AS col_type` FROM `catalog_product_entity_text` AS `t`)",
			"",
		)
	})

	t.Run("Interpolated", func(t *testing.T) {
		// Hint(CyS): We use the following query all over this file. This EAV
		// query gets generated here:
		// app/code/Magento/Eav/Model/ResourceModel/ReadHandler.php::execute but
		// the ORDER BY clause in each SELECT gets ignored so the final UNION
		// query over all EAV attribute tables does not have any sorting. This
		// bug gets not discovered because the scoped values for each entity
		// gets inserted into the database table after the default scope has
		// been inserted. Usually MySQL sorts the data how they are inserted
		// into the table ... if a row gets deleted MySQL might insert new data
		// into the space of the old row and if you even delete the default
		// scoped data and then recreate it the bug would bite you in your a**.
		// The EAV UNION has been fixed with our statement that you first sort
		// by attribute_id ASC and then by store_id ASC. If in the original
		// buggy query of M1/M2 the sorting would be really DESC then you would
		// never see the scoped data because the default data overwrites
		// everything when loading the PHP array.

		u := NewUnion(
			NewSelect().AddColumns("t.value", "t.attribute_id", "t.store_id").FromAlias("catalog_product_entity_{type}", "t").
				Where(Column("entity_id").Int64(1561), Column("store_id").In().Int64s(1, 0)),
		).
			StringReplace("{type}", "varchar", "int", "decimal", "datetime", "text").
			PreserveResultSet().
			All().OrderBy("attribute_id", "store_id").
			Interpolate()
		compareToSQL(t, u, nil,
			"(SELECT `t`.`value`, `t`.`attribute_id`, `t`.`store_id`, 0 AS `_preserve_result_set` FROM `catalog_product_entity_varchar` AS `t` WHERE (`entity_id` = 1561) AND (`store_id` IN (1,0)))\nUNION ALL\n(SELECT `t`.`value`, `t`.`attribute_id`, `t`.`store_id`, 1 AS `_preserve_result_set` FROM `catalog_product_entity_int` AS `t` WHERE (`entity_id` = 1561) AND (`store_id` IN (1,0)))\nUNION ALL\n(SELECT `t`.`value`, `t`.`attribute_id`, `t`.`store_id`, 2 AS `_preserve_result_set` FROM `catalog_product_entity_decimal` AS `t` WHERE (`entity_id` = 1561) AND (`store_id` IN (1,0)))\nUNION ALL\n(SELECT `t`.`value`, `t`.`attribute_id`, `t`.`store_id`, 3 AS `_preserve_result_set` FROM `catalog_product_entity_datetime` AS `t` WHERE (`entity_id` = 1561) AND (`store_id` IN (1,0)))\nUNION ALL\n(SELECT `t`.`value`, `t`.`attribute_id`, `t`.`store_id`, 4 AS `_preserve_result_set` FROM `catalog_product_entity_text` AS `t` WHERE (`entity_id` = 1561) AND (`store_id` IN (1,0)))\nORDER BY `_preserve_result_set`, `attribute_id` ASC, `store_id` ASC",
			"(SELECT `t`.`value`, `t`.`attribute_id`, `t`.`store_id`, 0 AS `_preserve_result_set` FROM `catalog_product_entity_varchar` AS `t` WHERE (`entity_id` = 1561) AND (`store_id` IN (1,0)))\nUNION ALL\n(SELECT `t`.`value`, `t`.`attribute_id`, `t`.`store_id`, 1 AS `_preserve_result_set` FROM `catalog_product_entity_int` AS `t` WHERE (`entity_id` = 1561) AND (`store_id` IN (1,0)))\nUNION ALL\n(SELECT `t`.`value`, `t`.`attribute_id`, `t`.`store_id`, 2 AS `_preserve_result_set` FROM `catalog_product_entity_decimal` AS `t` WHERE (`entity_id` = 1561) AND (`store_id` IN (1,0)))\nUNION ALL\n(SELECT `t`.`value`, `t`.`attribute_id`, `t`.`store_id`, 3 AS `_preserve_result_set` FROM `catalog_product_entity_datetime` AS `t` WHERE (`entity_id` = 1561) AND (`store_id` IN (1,0)))\nUNION ALL\n(SELECT `t`.`value`, `t`.`attribute_id`, `t`.`store_id`, 4 AS `_preserve_result_set` FROM `catalog_product_entity_text` AS `t` WHERE (`entity_id` = 1561) AND (`store_id` IN (1,0)))\nORDER BY `_preserve_result_set`, `attribute_id` ASC, `store_id` ASC",
		)
	})
}

func TestUnionTemplate_UseBuildCache(t *testing.T) {
	t.Parallel()

	u := NewUnion(
		NewSelect().AddColumns("t.value", "t.attribute_id", "t.store_id").FromAlias("catalog_product_entity_{type}", "t").
			Where(Column("entity_id").Int64(1561), Column("store_id").In().Int64s(1, 0)),
	).
		StringReplace("{type}", "varchar", "int", "decimal", "datetime", "text").
		PreserveResultSet().
		All().OrderBy("attribute_id", "store_id").BuildCache()

	const cachedSQLPlaceHolder = "(SELECT `t`.`value`, `t`.`attribute_id`, `t`.`store_id`, 0 AS `_preserve_result_set` FROM `catalog_product_entity_varchar` AS `t` WHERE (`entity_id` = ?) AND (`store_id` IN (?,?)))\nUNION ALL\n(SELECT `t`.`value`, `t`.`attribute_id`, `t`.`store_id`, 1 AS `_preserve_result_set` FROM `catalog_product_entity_int` AS `t` WHERE (`entity_id` = ?) AND (`store_id` IN (?,?)))\nUNION ALL\n(SELECT `t`.`value`, `t`.`attribute_id`, `t`.`store_id`, 2 AS `_preserve_result_set` FROM `catalog_product_entity_decimal` AS `t` WHERE (`entity_id` = ?) AND (`store_id` IN (?,?)))\nUNION ALL\n(SELECT `t`.`value`, `t`.`attribute_id`, `t`.`store_id`, 3 AS `_preserve_result_set` FROM `catalog_product_entity_datetime` AS `t` WHERE (`entity_id` = ?) AND (`store_id` IN (?,?)))\nUNION ALL\n(SELECT `t`.`value`, `t`.`attribute_id`, `t`.`store_id`, 4 AS `_preserve_result_set` FROM `catalog_product_entity_text` AS `t` WHERE (`entity_id` = ?) AND (`store_id` IN (?,?)))\nORDER BY `_preserve_result_set`, `attribute_id` ASC, `store_id` ASC"
	t.Run("without interpolate", func(t *testing.T) {
		for i := 0; i < 3; i++ {
			compareToSQL(t, u, nil,
				cachedSQLPlaceHolder,
				"",
				int64(1561), int64(1), int64(0), int64(1561), int64(1), int64(0), int64(1561), int64(1), int64(0), int64(1561), int64(1), int64(0), int64(1561), int64(1), int64(0),
			)
			assert.Equal(t, cachedSQLPlaceHolder, string(u.cacheSQL))
		}
	})

	t.Run("with interpolate", func(t *testing.T) {
		u.cacheSQL = nil

		const cachedSQLInterpolated = "(SELECT `t`.`value`, `t`.`attribute_id`, `t`.`store_id`, 0 AS `_preserve_result_set` FROM `catalog_product_entity_varchar` AS `t` WHERE (`entity_id` = 1561) AND (`store_id` IN (1,0)))\nUNION ALL\n(SELECT `t`.`value`, `t`.`attribute_id`, `t`.`store_id`, 1 AS `_preserve_result_set` FROM `catalog_product_entity_int` AS `t` WHERE (`entity_id` = 1561) AND (`store_id` IN (1,0)))\nUNION ALL\n(SELECT `t`.`value`, `t`.`attribute_id`, `t`.`store_id`, 2 AS `_preserve_result_set` FROM `catalog_product_entity_decimal` AS `t` WHERE (`entity_id` = 1561) AND (`store_id` IN (1,0)))\nUNION ALL\n(SELECT `t`.`value`, `t`.`attribute_id`, `t`.`store_id`, 3 AS `_preserve_result_set` FROM `catalog_product_entity_datetime` AS `t` WHERE (`entity_id` = 1561) AND (`store_id` IN (1,0)))\nUNION ALL\n(SELECT `t`.`value`, `t`.`attribute_id`, `t`.`store_id`, 4 AS `_preserve_result_set` FROM `catalog_product_entity_text` AS `t` WHERE (`entity_id` = 1561) AND (`store_id` IN (1,0)))\nORDER BY `_preserve_result_set`, `attribute_id` ASC, `store_id` ASC"
		for i := 0; i < 3; i++ {
			compareToSQL(t, u, nil,
				cachedSQLPlaceHolder,
				cachedSQLInterpolated,
				int64(1561), int64(1), int64(0), int64(1561), int64(1), int64(0), int64(1561), int64(1), int64(0), int64(1561), int64(1), int64(0), int64(1561), int64(1), int64(0),
			)
			assert.Equal(t, cachedSQLPlaceHolder, string(u.cacheSQL))
		}
	})
}
