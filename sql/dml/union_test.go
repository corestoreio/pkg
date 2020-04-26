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
	"database/sql"
	"testing"

	"github.com/corestoreio/errors"
	"github.com/corestoreio/pkg/util/assert"
)

func TestUnion_Basics(t *testing.T) {
	t.Run("simple", func(t *testing.T) {
		u := NewUnion(
			NewSelect("a", "b").From("tableAB").Where(Column("a").Int64(3)),
		)
		u.Append(
			NewSelect("c", "d").From("tableCD").Where(Column("d").Str("e")),
		)
		compareToSQL(t, u, errors.NoKind,
			"(SELECT `a`, `b` FROM `tableAB` WHERE (`a` = 3))\nUNION\n(SELECT `c`, `d` FROM `tableCD` WHERE (`d` = 'e'))",
			"(SELECT `a`, `b` FROM `tableAB` WHERE (`a` = 3))\nUNION\n(SELECT `c`, `d` FROM `tableCD` WHERE (`d` = 'e'))",
		)
	})

	t.Run("simple all", func(t *testing.T) {
		u := NewUnion(
			NewSelect("c", "d").From("tableCD").Where(Column("d").Str("e")),
			NewSelect("a", "b").From("tableAB").Where(Column("a").Int64(3)),
		).All()
		compareToSQL(t, u, errors.NoKind,
			"(SELECT `c`, `d` FROM `tableCD` WHERE (`d` = 'e'))\nUNION ALL\n(SELECT `a`, `b` FROM `tableAB` WHERE (`a` = 3))",
			"(SELECT `c`, `d` FROM `tableCD` WHERE (`d` = 'e'))\nUNION ALL\n(SELECT `a`, `b` FROM `tableAB` WHERE (`a` = 3))",
		)
	})

	t.Run("order by", func(t *testing.T) {
		u := NewUnion(
			NewSelect("a").AddColumnsAliases("d", "b").From("tableAD").Where(Column("d").Str("f")),
			NewSelect("a", "b").From("tableAB").Where(Column("a").Int64(3)),
		).All().OrderBy("a").OrderByDesc("b")

		compareToSQL(t, u, errors.NoKind,
			"(SELECT `a`, `d` AS `b` FROM `tableAD` WHERE (`d` = 'f'))\nUNION ALL\n(SELECT `a`, `b` FROM `tableAB` WHERE (`a` = 3))\nORDER BY `a`, `b` DESC",
			"(SELECT `a`, `d` AS `b` FROM `tableAD` WHERE (`d` = 'f'))\nUNION ALL\n(SELECT `a`, `b` FROM `tableAB` WHERE (`a` = 3))\nORDER BY `a`, `b` DESC",
		)
	})

	t.Run("preserve result set", func(t *testing.T) {
		u := NewUnion(
			NewSelect("a").AddColumnsAliases("d", "b").From("tableAD"),
			NewSelect("a", "b").From("tableAB").Where(Column("c").Between().Int64s(3, 5)),
		).All().OrderBy("a").OrderByDesc("b").PreserveResultSet()

		// testing idempotent function ToSQL
		for i := 0; i < 3; i++ {
			compareToSQL(t, u, errors.NoKind,
				"(SELECT `a`, `d` AS `b`, 0 AS `_preserve_result_set` FROM `tableAD`)\nUNION ALL\n(SELECT `a`, `b`, 1 AS `_preserve_result_set` FROM `tableAB` WHERE (`c` BETWEEN 3 AND 5))\nORDER BY `_preserve_result_set`, `a`, `b` DESC",
				"(SELECT `a`, `d` AS `b`, 0 AS `_preserve_result_set` FROM `tableAD`)\nUNION ALL\n(SELECT `a`, `b`, 1 AS `_preserve_result_set` FROM `tableAB` WHERE (`c` BETWEEN 3 AND 5))\nORDER BY `_preserve_result_set`, `a`, `b` DESC",
			)
		}
	})
	t.Run("intersec", func(t *testing.T) {
		u := NewUnion(
			NewSelect("a").From("tableAD"),
			NewSelect("b").From("tableAB"),
		).All().Intersect().OrderBy("a").OrderByDesc("b")
		// All gets ignored

		compareToSQL(t, u, errors.NoKind,
			"(SELECT `a` FROM `tableAD`)\nINTERSECT\n(SELECT `b` FROM `tableAB`)\nORDER BY `a`, `b` DESC",
			"(SELECT `a` FROM `tableAD`)\nINTERSECT\n(SELECT `b` FROM `tableAB`)\nORDER BY `a`, `b` DESC",
		)
	})
	t.Run("except", func(t *testing.T) {
		u := NewUnion(
			NewSelect("a").From("tableAD"),
			NewSelect("b").From("tableAB"),
		).All().Except()
		// All gets ignored

		compareToSQL(t, u, errors.NoKind,
			"(SELECT `a` FROM `tableAD`)\nEXCEPT\n(SELECT `b` FROM `tableAB`)",
			"(SELECT `a` FROM `tableAD`)\nEXCEPT\n(SELECT `b` FROM `tableAB`)",
		)
	})

	t.Run("placeholder question mark", func(t *testing.T) {
		u := NewUnion(
			NewSelect("a", "b").From("tableAD").Where(Column("a").Like().PlaceHolder()),
			NewSelect("a", "b").From("tableAB").Where(Column("c").Between().PlaceHolder()),
		).All().OrderBy("a").OrderByDesc("b").PreserveResultSet().
			WithDBR()
		// TODO(CYS) there is a small bug when using BETWEEN operator together
		// with named arguments. fix it. write a 2nd test function like this
		// used NamedArg.
		// testing idempotent function ToSQL
		for i := 0; i < 3; i++ {
			compareToSQL(t, u.TestWithArgs("XMEEN", 3.141, 6.283), errors.NoKind,
				"(SELECT `a`, `b`, 0 AS `_preserve_result_set` FROM `tableAD` WHERE (`a` LIKE ?))\nUNION ALL\n(SELECT `a`, `b`, 1 AS `_preserve_result_set` FROM `tableAB` WHERE (`c` BETWEEN ? AND ?))\nORDER BY `_preserve_result_set`, `a`, `b` DESC",
				"(SELECT `a`, `b`, 0 AS `_preserve_result_set` FROM `tableAD` WHERE (`a` LIKE 'XMEEN'))\nUNION ALL\n(SELECT `a`, `b`, 1 AS `_preserve_result_set` FROM `tableAB` WHERE (`c` BETWEEN 3.141 AND 6.283))\nORDER BY `_preserve_result_set`, `a`, `b` DESC",
				"XMEEN", 3.141, 6.283,
			)
		}
		assert.Exactly(t, []string{"a", "c"}, u.base.qualifiedColumns)
	})
}

func TestUnion_DisableBuildCache(t *testing.T) {
	u := NewUnion(
		NewSelect("a").AddColumnsAliases("d", "b").From("tableAD"),
		NewSelect("a", "b").From("tableAB").Where(Column("b").Float64(3.14159)),
	).
		All().
		Unsafe().
		StringReplace("MyKey", "a", "b", "c"). // does nothing because more than one NewSelect functions
		OrderBy("a").OrderByDesc("b").OrderBy(`concat("c",b,"d")`).
		PreserveResultSet()

	const cachedSQLPlaceHolder = "(SELECT `a`, `d` AS `b`, 0 AS `_preserve_result_set` FROM `tableAD`)\n" +
		"UNION ALL\n" +
		"(SELECT `a`, `b`, 1 AS `_preserve_result_set` FROM `tableAB` WHERE (`b` = 3.14159))\n" +
		"ORDER BY `_preserve_result_set`, `a`, `b` DESC, concat(\"c\",b,\"d\")"
	t.Run("without interpolate", func(t *testing.T) {
		for i := 0; i < 3; i++ {
			compareToSQL(t, u, errors.NoKind,
				cachedSQLPlaceHolder,
				"",
			)
		}
	})

	t.Run("with interpolate", func(t *testing.T) {
		const cachedSQLInterpolated = "(SELECT `a`, `d` AS `b`, 0 AS `_preserve_result_set` FROM `tableAD`)\nUNION ALL\n(SELECT `a`, `b`, 1 AS `_preserve_result_set` FROM `tableAB` WHERE (`b` = 3.14159))\nORDER BY `_preserve_result_set`, `a`, `b` DESC, concat(\"c\",b,\"d\")"
		for i := 0; i < 3; i++ {
			compareToSQL(t, u, errors.NoKind,
				cachedSQLPlaceHolder,
				cachedSQLInterpolated,
			)
		}
	})
	assert.Exactly(t, []string{"", "(SELECT `a`, `d` AS `b`, 0 AS `_preserve_result_set` FROM `tableAD`)\nUNION ALL\n(SELECT `a`, `b`, 1 AS `_preserve_result_set` FROM `tableAB` WHERE (`b` = 3.14159))\nORDER BY `_preserve_result_set`, `a`, `b` DESC, concat(\"c\",b,\"d\")"},
		u.CachedQueries())
}

func TestNewUnionTemplate(t *testing.T) {
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

	t.Run("full statement EAV no placeholder", func(t *testing.T) {
		u := NewUnion(
			NewSelect().AddColumns("t.value", "t.attribute_id").AddColumnsAliases("t.$column$", "col_type").
				FromAlias("catalog_product_entity_$type$", "t").
				Where(Column("entity_id").Int64(1561), Column("store_id").In().Int64s(1, 0)).
				OrderByDesc("t.$column$_store_id"),
		).
			StringReplace("$type$", "varchar", "int", "decimal", "datetime", "text").
			StringReplace("$column$", "varcharX", "intX", "decimalX", "datetimeX", "textX").
			PreserveResultSet().
			All().
			OrderByDesc("col_type")

		// testing idempotent function ToSQL
		for i := 0; i < 3; i++ {
			compareToSQL2(t, u, errors.NoKind,
				"(SELECT `t`.`value`, `t`.`attribute_id`, `t`.`varcharX` AS `col_type`, 0 AS `_preserve_result_set` FROM `catalog_product_entity_varchar` AS `t` WHERE (`entity_id` = 1561) AND (`store_id` IN (1,0)) ORDER BY `t`.`varcharX_store_id` DESC)\n"+
					"UNION ALL\n"+
					"(SELECT `t`.`value`, `t`.`attribute_id`, `t`.`intX` AS `col_type`, 1 AS `_preserve_result_set` FROM `catalog_product_entity_int` AS `t` WHERE (`entity_id` = 1561) AND (`store_id` IN (1,0)) ORDER BY `t`.`intX_store_id` DESC)\n"+
					"UNION ALL\n"+
					"(SELECT `t`.`value`, `t`.`attribute_id`, `t`.`decimalX` AS `col_type`, 2 AS `_preserve_result_set` FROM `catalog_product_entity_decimal` AS `t` WHERE (`entity_id` = 1561) AND (`store_id` IN (1,0)) ORDER BY `t`.`decimalX_store_id` DESC)\n"+
					"UNION ALL\n"+
					"(SELECT `t`.`value`, `t`.`attribute_id`, `t`.`datetimeX` AS `col_type`, 3 AS `_preserve_result_set` FROM `catalog_product_entity_datetime` AS `t` WHERE (`entity_id` = 1561) AND (`store_id` IN (1,0)) ORDER BY `t`.`datetimeX_store_id` DESC)\n"+
					"UNION ALL\n"+
					"(SELECT `t`.`value`, `t`.`attribute_id`, `t`.`textX` AS `col_type`, 4 AS `_preserve_result_set` FROM `catalog_product_entity_text` AS `t` WHERE (`entity_id` = 1561) AND (`store_id` IN (1,0)) ORDER BY `t`.`textX_store_id` DESC)\n"+
					"ORDER BY `_preserve_result_set`, `col_type` DESC",
			)
		}
	})

	t.Run("full statement EAV with placeholder", func(t *testing.T) {
		u := NewUnion(
			NewSelect().AddColumns("t.value", "t.attribute_id").AddColumnsAliases("t.$column$", "col_type").
				FromAlias("catalog_product_entity_$type$", "t").
				Where(Column("entity_id").PlaceHolder(), Column("store_id").In().PlaceHolder()).
				OrderByDesc("t.$column$_store_id"),
		).
			StringReplace("$type$", "varchar", "int", "decimal", "datetime", "text").
			StringReplace("$column$", "varcharX", "intX", "decimalX", "datetimeX", "textX").
			PreserveResultSet().
			All().
			OrderByDesc("col_type")

		ua := u.WithDBR().TestWithArgs(1563, []int{3, 4})

		// testing idempotent function ToSQL
		for i := 0; i < 3; i++ {
			compareToSQL(t, ua, errors.NoKind,
				"(SELECT `t`.`value`, `t`.`attribute_id`, `t`.`varcharX` AS `col_type`, 0 AS `_preserve_result_set` FROM `catalog_product_entity_varchar` AS `t` WHERE (`entity_id` = ?) AND (`store_id` IN ?) ORDER BY `t`.`varcharX_store_id` DESC)\n"+
					"UNION ALL\n"+
					"(SELECT `t`.`value`, `t`.`attribute_id`, `t`.`intX` AS `col_type`, 1 AS `_preserve_result_set` FROM `catalog_product_entity_int` AS `t` WHERE (`entity_id` = ?) AND (`store_id` IN ?) ORDER BY `t`.`intX_store_id` DESC)\n"+
					"UNION ALL\n"+
					"(SELECT `t`.`value`, `t`.`attribute_id`, `t`.`decimalX` AS `col_type`, 2 AS `_preserve_result_set` FROM `catalog_product_entity_decimal` AS `t` WHERE (`entity_id` = ?) AND (`store_id` IN ?) ORDER BY `t`.`decimalX_store_id` DESC)\n"+
					"UNION ALL\n"+
					"(SELECT `t`.`value`, `t`.`attribute_id`, `t`.`datetimeX` AS `col_type`, 3 AS `_preserve_result_set` FROM `catalog_product_entity_datetime` AS `t` WHERE (`entity_id` = ?) AND (`store_id` IN ?) ORDER BY `t`.`datetimeX_store_id` DESC)\n"+
					"UNION ALL\n"+
					"(SELECT `t`.`value`, `t`.`attribute_id`, `t`.`textX` AS `col_type`, 4 AS `_preserve_result_set` FROM `catalog_product_entity_text` AS `t` WHERE (`entity_id` = ?) AND (`store_id` IN ?) ORDER BY `t`.`textX_store_id` DESC)\n"+
					"ORDER BY `_preserve_result_set`, `col_type` DESC",
				"(SELECT `t`.`value`, `t`.`attribute_id`, `t`.`varcharX` AS `col_type`, 0 AS `_preserve_result_set` FROM `catalog_product_entity_varchar` AS `t` WHERE (`entity_id` = 1563) AND (`store_id` IN (3,4)) ORDER BY `t`.`varcharX_store_id` DESC)\n"+
					"UNION ALL\n"+
					"(SELECT `t`.`value`, `t`.`attribute_id`, `t`.`intX` AS `col_type`, 1 AS `_preserve_result_set` FROM `catalog_product_entity_int` AS `t` WHERE (`entity_id` = 1563) AND (`store_id` IN (3,4)) ORDER BY `t`.`intX_store_id` DESC)\n"+
					"UNION ALL\n"+
					"(SELECT `t`.`value`, `t`.`attribute_id`, `t`.`decimalX` AS `col_type`, 2 AS `_preserve_result_set` FROM `catalog_product_entity_decimal` AS `t` WHERE (`entity_id` = 1563) AND (`store_id` IN (3,4)) ORDER BY `t`.`decimalX_store_id` DESC)\n"+
					"UNION ALL\n"+
					"(SELECT `t`.`value`, `t`.`attribute_id`, `t`.`datetimeX` AS `col_type`, 3 AS `_preserve_result_set` FROM `catalog_product_entity_datetime` AS `t` WHERE (`entity_id` = 1563) AND (`store_id` IN (3,4)) ORDER BY `t`.`datetimeX_store_id` DESC)\n"+
					"UNION ALL\n"+
					"(SELECT `t`.`value`, `t`.`attribute_id`, `t`.`textX` AS `col_type`, 4 AS `_preserve_result_set` FROM `catalog_product_entity_text` AS `t` WHERE (`entity_id` = 1563) AND (`store_id` IN (3,4)) ORDER BY `t`.`textX_store_id` DESC)\n"+
					"ORDER BY `_preserve_result_set`, `col_type` DESC",
				int64(1563), int64(3), int64(4), int64(1563), int64(3), int64(4), int64(1563), int64(3), int64(4), int64(1563), int64(3), int64(4), int64(1563), int64(3), int64(4),
			)
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
			NewSelect().AddColumns("t.value", "t.attribute_id").AddColumnsAliases("t.$column$", "col_type").
				FromAlias("catalog_product_entity_$type$", "t"),
		).
			StringReplace("$type$", "varchar", "int", "decimal", "datetime", "text").
			StringReplace("$column$", "varcharX", "intX", "decimalX", "datetimeX")
	})
	t.Run("StringReplace 2nd call too many values and nothing should happen", func(t *testing.T) {
		u := NewUnion(
			NewSelect().AddColumns("t.value", "t.attribute_id").AddColumnsAliases("t.$column$", "col_type").
				FromAlias("catalog_product_entity_$type$", "t"),
		).
			StringReplace("$type$", "varchar", "int", "decimal", "datetime", "text").
			StringReplace("$column$", "varcharX", "intX", "decimalX", "datetimeX", "textX", "bytesX")
		compareToSQL(t, u, errors.NoKind,
			"(SELECT `t`.`value`, `t`.`attribute_id`, `t`.`varcharX` AS `col_type` FROM `catalog_product_entity_varchar` AS `t`)\nUNION\n(SELECT `t`.`value`, `t`.`attribute_id`, `t`.`intX` AS `col_type` FROM `catalog_product_entity_int` AS `t`)\nUNION\n(SELECT `t`.`value`, `t`.`attribute_id`, `t`.`decimalX` AS `col_type` FROM `catalog_product_entity_decimal` AS `t`)\nUNION\n(SELECT `t`.`value`, `t`.`attribute_id`, `t`.`datetimeX` AS `col_type` FROM `catalog_product_entity_datetime` AS `t`)\nUNION\n(SELECT `t`.`value`, `t`.`attribute_id`, `t`.`textX` AS `col_type` FROM `catalog_product_entity_text` AS `t`)",
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
		// never see the scoped value because the default value overwrites
		// everything when loading the PHP array.

		u := NewUnion(
			NewSelect().AddColumns("t.value", "t.attribute_id", "t.store_id").FromAlias("catalog_product_entity_$type$", "t").
				Where(Column("entity_id").Int64(1561), Column("store_id").In().Int64s(1, 0)),
		).
			StringReplace("$type$", "varchar", "int", "decimal", "datetime", "text").
			PreserveResultSet().
			All().OrderBy("attribute_id", "store_id")
		compareToSQL(t, u, errors.NoKind,
			"(SELECT `t`.`value`, `t`.`attribute_id`, `t`.`store_id`, 0 AS `_preserve_result_set` FROM `catalog_product_entity_varchar` AS `t` WHERE (`entity_id` = 1561) AND (`store_id` IN (1,0)))\nUNION ALL\n(SELECT `t`.`value`, `t`.`attribute_id`, `t`.`store_id`, 1 AS `_preserve_result_set` FROM `catalog_product_entity_int` AS `t` WHERE (`entity_id` = 1561) AND (`store_id` IN (1,0)))\nUNION ALL\n(SELECT `t`.`value`, `t`.`attribute_id`, `t`.`store_id`, 2 AS `_preserve_result_set` FROM `catalog_product_entity_decimal` AS `t` WHERE (`entity_id` = 1561) AND (`store_id` IN (1,0)))\nUNION ALL\n(SELECT `t`.`value`, `t`.`attribute_id`, `t`.`store_id`, 3 AS `_preserve_result_set` FROM `catalog_product_entity_datetime` AS `t` WHERE (`entity_id` = 1561) AND (`store_id` IN (1,0)))\nUNION ALL\n(SELECT `t`.`value`, `t`.`attribute_id`, `t`.`store_id`, 4 AS `_preserve_result_set` FROM `catalog_product_entity_text` AS `t` WHERE (`entity_id` = 1561) AND (`store_id` IN (1,0)))\nORDER BY `_preserve_result_set`, `attribute_id`, `store_id`",
			"(SELECT `t`.`value`, `t`.`attribute_id`, `t`.`store_id`, 0 AS `_preserve_result_set` FROM `catalog_product_entity_varchar` AS `t` WHERE (`entity_id` = 1561) AND (`store_id` IN (1,0)))\nUNION ALL\n(SELECT `t`.`value`, `t`.`attribute_id`, `t`.`store_id`, 1 AS `_preserve_result_set` FROM `catalog_product_entity_int` AS `t` WHERE (`entity_id` = 1561) AND (`store_id` IN (1,0)))\nUNION ALL\n(SELECT `t`.`value`, `t`.`attribute_id`, `t`.`store_id`, 2 AS `_preserve_result_set` FROM `catalog_product_entity_decimal` AS `t` WHERE (`entity_id` = 1561) AND (`store_id` IN (1,0)))\nUNION ALL\n(SELECT `t`.`value`, `t`.`attribute_id`, `t`.`store_id`, 3 AS `_preserve_result_set` FROM `catalog_product_entity_datetime` AS `t` WHERE (`entity_id` = 1561) AND (`store_id` IN (1,0)))\nUNION ALL\n(SELECT `t`.`value`, `t`.`attribute_id`, `t`.`store_id`, 4 AS `_preserve_result_set` FROM `catalog_product_entity_text` AS `t` WHERE (`entity_id` = 1561) AND (`store_id` IN (1,0)))\nORDER BY `_preserve_result_set`, `attribute_id`, `store_id`",
		)
	})
}

func TestUnionTemplate_DisableBuildCache(t *testing.T) {
	u := NewUnion(
		NewSelect().AddColumns("t.value", "t.attribute_id", "t.store_id").FromAlias("catalog_product_entity_$type$", "t").
			Where(Column("entity_id").Int64(1561), Column("store_id").In().Int64s(1, 0)),
	).
		StringReplace("$type$", "varchar", "int", "decimal", "datetime", "text").
		PreserveResultSet().
		All().OrderBy("attribute_id", "store_id")

	const cachedSQLPlaceHolder = "(SELECT `t`.`value`, `t`.`attribute_id`, `t`.`store_id`, 0 AS `_preserve_result_set` FROM `catalog_product_entity_varchar` AS `t` WHERE (`entity_id` = 1561) AND (`store_id` IN (1,0)))\nUNION ALL\n(SELECT `t`.`value`, `t`.`attribute_id`, `t`.`store_id`, 1 AS `_preserve_result_set` FROM `catalog_product_entity_int` AS `t` WHERE (`entity_id` = 1561) AND (`store_id` IN (1,0)))\nUNION ALL\n(SELECT `t`.`value`, `t`.`attribute_id`, `t`.`store_id`, 2 AS `_preserve_result_set` FROM `catalog_product_entity_decimal` AS `t` WHERE (`entity_id` = 1561) AND (`store_id` IN (1,0)))\nUNION ALL\n(SELECT `t`.`value`, `t`.`attribute_id`, `t`.`store_id`, 3 AS `_preserve_result_set` FROM `catalog_product_entity_datetime` AS `t` WHERE (`entity_id` = 1561) AND (`store_id` IN (1,0)))\nUNION ALL\n(SELECT `t`.`value`, `t`.`attribute_id`, `t`.`store_id`, 4 AS `_preserve_result_set` FROM `catalog_product_entity_text` AS `t` WHERE (`entity_id` = 1561) AND (`store_id` IN (1,0)))\nORDER BY `_preserve_result_set`, `attribute_id`, `store_id`"
	t.Run("without interpolate", func(t *testing.T) {
		for i := 0; i < 3; i++ {
			u.WithCacheKey("index_%d", i)
			compareToSQL(t, u, errors.NoKind,
				cachedSQLPlaceHolder,
				"",
			)
		}
		assert.Exactly(t, []string{"index_0", "(SELECT `t`.`value`, `t`.`attribute_id`, `t`.`store_id`, 0 AS `_preserve_result_set` FROM `catalog_product_entity_varchar` AS `t` WHERE (`entity_id` = 1561) AND (`store_id` IN (1,0)))\nUNION ALL\n(SELECT `t`.`value`, `t`.`attribute_id`, `t`.`store_id`, 1 AS `_preserve_result_set` FROM `catalog_product_entity_int` AS `t` WHERE (`entity_id` = 1561) AND (`store_id` IN (1,0)))\nUNION ALL\n(SELECT `t`.`value`, `t`.`attribute_id`, `t`.`store_id`, 2 AS `_preserve_result_set` FROM `catalog_product_entity_decimal` AS `t` WHERE (`entity_id` = 1561) AND (`store_id` IN (1,0)))\nUNION ALL\n(SELECT `t`.`value`, `t`.`attribute_id`, `t`.`store_id`, 3 AS `_preserve_result_set` FROM `catalog_product_entity_datetime` AS `t` WHERE (`entity_id` = 1561) AND (`store_id` IN (1,0)))\nUNION ALL\n(SELECT `t`.`value`, `t`.`attribute_id`, `t`.`store_id`, 4 AS `_preserve_result_set` FROM `catalog_product_entity_text` AS `t` WHERE (`entity_id` = 1561) AND (`store_id` IN (1,0)))\nORDER BY `_preserve_result_set`, `attribute_id`, `store_id`", "index_1", "(SELECT `t`.`value`, `t`.`attribute_id`, `t`.`store_id`, 0 AS `_preserve_result_set` FROM `catalog_product_entity_varchar` AS `t` WHERE (`entity_id` = 1561) AND (`store_id` IN (1,0)))\nUNION ALL\n(SELECT `t`.`value`, `t`.`attribute_id`, `t`.`store_id`, 1 AS `_preserve_result_set` FROM `catalog_product_entity_int` AS `t` WHERE (`entity_id` = 1561) AND (`store_id` IN (1,0)))\nUNION ALL\n(SELECT `t`.`value`, `t`.`attribute_id`, `t`.`store_id`, 2 AS `_preserve_result_set` FROM `catalog_product_entity_decimal` AS `t` WHERE (`entity_id` = 1561) AND (`store_id` IN (1,0)))\nUNION ALL\n(SELECT `t`.`value`, `t`.`attribute_id`, `t`.`store_id`, 3 AS `_preserve_result_set` FROM `catalog_product_entity_datetime` AS `t` WHERE (`entity_id` = 1561) AND (`store_id` IN (1,0)))\nUNION ALL\n(SELECT `t`.`value`, `t`.`attribute_id`, `t`.`store_id`, 4 AS `_preserve_result_set` FROM `catalog_product_entity_text` AS `t` WHERE (`entity_id` = 1561) AND (`store_id` IN (1,0)))\nORDER BY `_preserve_result_set`, `attribute_id`, `store_id`", "index_2", "(SELECT `t`.`value`, `t`.`attribute_id`, `t`.`store_id`, 0 AS `_preserve_result_set` FROM `catalog_product_entity_varchar` AS `t` WHERE (`entity_id` = 1561) AND (`store_id` IN (1,0)))\nUNION ALL\n(SELECT `t`.`value`, `t`.`attribute_id`, `t`.`store_id`, 1 AS `_preserve_result_set` FROM `catalog_product_entity_int` AS `t` WHERE (`entity_id` = 1561) AND (`store_id` IN (1,0)))\nUNION ALL\n(SELECT `t`.`value`, `t`.`attribute_id`, `t`.`store_id`, 2 AS `_preserve_result_set` FROM `catalog_product_entity_decimal` AS `t` WHERE (`entity_id` = 1561) AND (`store_id` IN (1,0)))\nUNION ALL\n(SELECT `t`.`value`, `t`.`attribute_id`, `t`.`store_id`, 3 AS `_preserve_result_set` FROM `catalog_product_entity_datetime` AS `t` WHERE (`entity_id` = 1561) AND (`store_id` IN (1,0)))\nUNION ALL\n(SELECT `t`.`value`, `t`.`attribute_id`, `t`.`store_id`, 4 AS `_preserve_result_set` FROM `catalog_product_entity_text` AS `t` WHERE (`entity_id` = 1561) AND (`store_id` IN (1,0)))\nORDER BY `_preserve_result_set`, `attribute_id`, `store_id`"}, u.CachedQueries())
	})

	t.Run("with interpolate", func(t *testing.T) {
		const cachedSQLInterpolated = "(SELECT `t`.`value`, `t`.`attribute_id`, `t`.`store_id`, 0 AS `_preserve_result_set` FROM `catalog_product_entity_varchar` AS `t` WHERE (`entity_id` = 1561) AND (`store_id` IN (1,0)))\nUNION ALL\n(SELECT `t`.`value`, `t`.`attribute_id`, `t`.`store_id`, 1 AS `_preserve_result_set` FROM `catalog_product_entity_int` AS `t` WHERE (`entity_id` = 1561) AND (`store_id` IN (1,0)))\nUNION ALL\n(SELECT `t`.`value`, `t`.`attribute_id`, `t`.`store_id`, 2 AS `_preserve_result_set` FROM `catalog_product_entity_decimal` AS `t` WHERE (`entity_id` = 1561) AND (`store_id` IN (1,0)))\nUNION ALL\n(SELECT `t`.`value`, `t`.`attribute_id`, `t`.`store_id`, 3 AS `_preserve_result_set` FROM `catalog_product_entity_datetime` AS `t` WHERE (`entity_id` = 1561) AND (`store_id` IN (1,0)))\nUNION ALL\n(SELECT `t`.`value`, `t`.`attribute_id`, `t`.`store_id`, 4 AS `_preserve_result_set` FROM `catalog_product_entity_text` AS `t` WHERE (`entity_id` = 1561) AND (`store_id` IN (1,0)))\nORDER BY `_preserve_result_set`, `attribute_id`, `store_id`"
		for i := 0; i < 3; i++ {
			u.WithCacheKey("index_%d", i)
			compareToSQL(t, u, errors.NoKind,
				cachedSQLPlaceHolder,
				cachedSQLInterpolated,
			)
		}
		assert.Exactly(t, []string{"index_0", "(SELECT `t`.`value`, `t`.`attribute_id`, `t`.`store_id`, 0 AS `_preserve_result_set` FROM `catalog_product_entity_varchar` AS `t` WHERE (`entity_id` = 1561) AND (`store_id` IN (1,0)))\nUNION ALL\n(SELECT `t`.`value`, `t`.`attribute_id`, `t`.`store_id`, 1 AS `_preserve_result_set` FROM `catalog_product_entity_int` AS `t` WHERE (`entity_id` = 1561) AND (`store_id` IN (1,0)))\nUNION ALL\n(SELECT `t`.`value`, `t`.`attribute_id`, `t`.`store_id`, 2 AS `_preserve_result_set` FROM `catalog_product_entity_decimal` AS `t` WHERE (`entity_id` = 1561) AND (`store_id` IN (1,0)))\nUNION ALL\n(SELECT `t`.`value`, `t`.`attribute_id`, `t`.`store_id`, 3 AS `_preserve_result_set` FROM `catalog_product_entity_datetime` AS `t` WHERE (`entity_id` = 1561) AND (`store_id` IN (1,0)))\nUNION ALL\n(SELECT `t`.`value`, `t`.`attribute_id`, `t`.`store_id`, 4 AS `_preserve_result_set` FROM `catalog_product_entity_text` AS `t` WHERE (`entity_id` = 1561) AND (`store_id` IN (1,0)))\nORDER BY `_preserve_result_set`, `attribute_id`, `store_id`", "index_1", "(SELECT `t`.`value`, `t`.`attribute_id`, `t`.`store_id`, 0 AS `_preserve_result_set` FROM `catalog_product_entity_varchar` AS `t` WHERE (`entity_id` = 1561) AND (`store_id` IN (1,0)))\nUNION ALL\n(SELECT `t`.`value`, `t`.`attribute_id`, `t`.`store_id`, 1 AS `_preserve_result_set` FROM `catalog_product_entity_int` AS `t` WHERE (`entity_id` = 1561) AND (`store_id` IN (1,0)))\nUNION ALL\n(SELECT `t`.`value`, `t`.`attribute_id`, `t`.`store_id`, 2 AS `_preserve_result_set` FROM `catalog_product_entity_decimal` AS `t` WHERE (`entity_id` = 1561) AND (`store_id` IN (1,0)))\nUNION ALL\n(SELECT `t`.`value`, `t`.`attribute_id`, `t`.`store_id`, 3 AS `_preserve_result_set` FROM `catalog_product_entity_datetime` AS `t` WHERE (`entity_id` = 1561) AND (`store_id` IN (1,0)))\nUNION ALL\n(SELECT `t`.`value`, `t`.`attribute_id`, `t`.`store_id`, 4 AS `_preserve_result_set` FROM `catalog_product_entity_text` AS `t` WHERE (`entity_id` = 1561) AND (`store_id` IN (1,0)))\nORDER BY `_preserve_result_set`, `attribute_id`, `store_id`", "index_2", "(SELECT `t`.`value`, `t`.`attribute_id`, `t`.`store_id`, 0 AS `_preserve_result_set` FROM `catalog_product_entity_varchar` AS `t` WHERE (`entity_id` = 1561) AND (`store_id` IN (1,0)))\nUNION ALL\n(SELECT `t`.`value`, `t`.`attribute_id`, `t`.`store_id`, 1 AS `_preserve_result_set` FROM `catalog_product_entity_int` AS `t` WHERE (`entity_id` = 1561) AND (`store_id` IN (1,0)))\nUNION ALL\n(SELECT `t`.`value`, `t`.`attribute_id`, `t`.`store_id`, 2 AS `_preserve_result_set` FROM `catalog_product_entity_decimal` AS `t` WHERE (`entity_id` = 1561) AND (`store_id` IN (1,0)))\nUNION ALL\n(SELECT `t`.`value`, `t`.`attribute_id`, `t`.`store_id`, 3 AS `_preserve_result_set` FROM `catalog_product_entity_datetime` AS `t` WHERE (`entity_id` = 1561) AND (`store_id` IN (1,0)))\nUNION ALL\n(SELECT `t`.`value`, `t`.`attribute_id`, `t`.`store_id`, 4 AS `_preserve_result_set` FROM `catalog_product_entity_text` AS `t` WHERE (`entity_id` = 1561) AND (`store_id` IN (1,0)))\nORDER BY `_preserve_result_set`, `attribute_id`, `store_id`"}, u.CachedQueries())
	})
}

func TestUnionTemplate_ReuseArgs(t *testing.T) {
	u := NewUnion(
		NewSelect().AddColumns("t.value", "t.attribute_id", "t.store_id").FromAlias("catalog_product_entity_$type$", "t").
			Where(Column("entity_id").NamedArg("entityID"), Column("store_id").In().NamedArg("storeID")),
	).
		StringReplace("$type$", "varchar", "int", "decimal", "datetime", "text").
		PreserveResultSet().
		All().OrderBy("attribute_id", "store_id")

	const wantSQLPH = "(SELECT `t`.`value`, `t`.`attribute_id`, `t`.`store_id`, 0 AS `_preserve_result_set` FROM `catalog_product_entity_varchar` AS `t` WHERE (`entity_id` = ?) AND (`store_id` IN ?))\n" +
		"UNION ALL\n" +
		"(SELECT `t`.`value`, `t`.`attribute_id`, `t`.`store_id`, 1 AS `_preserve_result_set` FROM `catalog_product_entity_int` AS `t` WHERE (`entity_id` = ?) AND (`store_id` IN ?))\n" +
		"UNION ALL\n" +
		"(SELECT `t`.`value`, `t`.`attribute_id`, `t`.`store_id`, 2 AS `_preserve_result_set` FROM `catalog_product_entity_decimal` AS `t` WHERE (`entity_id` = ?) AND (`store_id` IN ?))\n" +
		"UNION ALL\n" +
		"(SELECT `t`.`value`, `t`.`attribute_id`, `t`.`store_id`, 3 AS `_preserve_result_set` FROM `catalog_product_entity_datetime` AS `t` WHERE (`entity_id` = ?) AND (`store_id` IN ?))\n" +
		"UNION ALL\n" +
		"(SELECT `t`.`value`, `t`.`attribute_id`, `t`.`store_id`, 4 AS `_preserve_result_set` FROM `catalog_product_entity_text` AS `t` WHERE (`entity_id` = ?) AND (`store_id` IN ?))\n" +
		"ORDER BY `_preserve_result_set`, `attribute_id`, `store_id`"

	t.Run("run1", func(t *testing.T) {
		compareToSQL(t, u.WithDBR().TestWithArgs(sql.Named("storeID", []int64{4, 6}), sql.Named("entityID", 5)), errors.NoKind,
			wantSQLPH,
			"(SELECT `t`.`value`, `t`.`attribute_id`, `t`.`store_id`, 0 AS `_preserve_result_set` FROM `catalog_product_entity_varchar` AS `t` WHERE (`entity_id` = 5) AND (`store_id` IN (4,6)))\n"+
				"UNION ALL\n"+
				"(SELECT `t`.`value`, `t`.`attribute_id`, `t`.`store_id`, 1 AS `_preserve_result_set` FROM `catalog_product_entity_int` AS `t` WHERE (`entity_id` = 5) AND (`store_id` IN (4,6)))\n"+
				"UNION ALL\n"+
				"(SELECT `t`.`value`, `t`.`attribute_id`, `t`.`store_id`, 2 AS `_preserve_result_set` FROM `catalog_product_entity_decimal` AS `t` WHERE (`entity_id` = 5) AND (`store_id` IN (4,6)))\n"+
				"UNION ALL\n"+
				"(SELECT `t`.`value`, `t`.`attribute_id`, `t`.`store_id`, 3 AS `_preserve_result_set` FROM `catalog_product_entity_datetime` AS `t` WHERE (`entity_id` = 5) AND (`store_id` IN (4,6)))\n"+
				"UNION ALL\n"+
				"(SELECT `t`.`value`, `t`.`attribute_id`, `t`.`store_id`, 4 AS `_preserve_result_set` FROM `catalog_product_entity_text` AS `t` WHERE (`entity_id` = 5) AND (`store_id` IN (4,6)))\n"+
				"ORDER BY `_preserve_result_set`, `attribute_id`, `store_id`",
			int64(5), int64(4), int64(6), // varchar
			int64(5), int64(4), int64(6), // int
			int64(5), int64(4), int64(6), // decimal
			int64(5), int64(4), int64(6), // datetime
			int64(5), int64(4), int64(6), // text
		)
		assert.Exactly(t, []string{":entityID", ":storeID"}, u.qualifiedColumns)
	})

	t.Run("with interpolate", func(t *testing.T) {
		compareToSQL(t, u.WithDBR().TestWithArgs(sql.Named("entityID", 4), sql.Named("storeID", []int64{8, 11})), errors.NoKind,
			wantSQLPH,
			"(SELECT `t`.`value`, `t`.`attribute_id`, `t`.`store_id`, 0 AS `_preserve_result_set` FROM `catalog_product_entity_varchar` AS `t` WHERE (`entity_id` = 4) AND (`store_id` IN (8,11)))\n"+
				"UNION ALL\n"+
				"(SELECT `t`.`value`, `t`.`attribute_id`, `t`.`store_id`, 1 AS `_preserve_result_set` FROM `catalog_product_entity_int` AS `t` WHERE (`entity_id` = 4) AND (`store_id` IN (8,11)))\n"+
				"UNION ALL\n"+
				"(SELECT `t`.`value`, `t`.`attribute_id`, `t`.`store_id`, 2 AS `_preserve_result_set` FROM `catalog_product_entity_decimal` AS `t` WHERE (`entity_id` = 4) AND (`store_id` IN (8,11)))\n"+
				"UNION ALL\n"+
				"(SELECT `t`.`value`, `t`.`attribute_id`, `t`.`store_id`, 3 AS `_preserve_result_set` FROM `catalog_product_entity_datetime` AS `t` WHERE (`entity_id` = 4) AND (`store_id` IN (8,11)))\n"+
				"UNION ALL\n"+
				"(SELECT `t`.`value`, `t`.`attribute_id`, `t`.`store_id`, 4 AS `_preserve_result_set` FROM `catalog_product_entity_text` AS `t` WHERE (`entity_id` = 4) AND (`store_id` IN (8,11)))\n"+
				"ORDER BY `_preserve_result_set`, `attribute_id`, `store_id`",
			int64(4), int64(8), int64(11), // varchar
			int64(4), int64(8), int64(11), // int
			int64(4), int64(8), int64(11), // decimal
			int64(4), int64(8), int64(11), // datetime
			int64(4), int64(8), int64(11), // text
		)
		assert.Exactly(t, []string{":entityID", ":storeID"}, u.qualifiedColumns)
	})
}
