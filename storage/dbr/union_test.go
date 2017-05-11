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

	"github.com/corestoreio/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestUnion(t *testing.T) {
	t.Parallel()
	t.Run("simple", func(t *testing.T) {
		u := NewUnion(
			NewSelect("a", "b").From("tableAB").Where(Column("a", ArgInt64(3))),
		)
		u.Append(
			NewSelect("c", "d").From("tableCD").Where(Column("d", ArgString("e"))),
		)

		uStr, args, err := u.ToSQL()
		assert.NoError(t, err, "%+v", err)
		assert.Exactly(t, []interface{}{int64(3), "e"}, args.Interfaces())
		assert.Exactly(t,
			"(SELECT `a`, `b` FROM `tableAB` WHERE (`a` = ?))\nUNION\n(SELECT `c`, `d` FROM `tableCD` WHERE (`d` = ?))",
			uStr)
	})

	t.Run("simple all", func(t *testing.T) {
		u := NewUnion(
			NewSelect("c", "d").From("tableCD").Where(Column("d", ArgString("e"))),
			NewSelect("a", "b").From("tableAB").Where(Column("a", ArgInt64(3))),
		).All()

		uStr, args, err := u.ToSQL()
		assert.NoError(t, err, "%+v", err)
		assert.Exactly(t, []interface{}{"e", int64(3)}, args.Interfaces())
		assert.Exactly(t,
			"(SELECT `c`, `d` FROM `tableCD` WHERE (`d` = ?))\nUNION ALL\n(SELECT `a`, `b` FROM `tableAB` WHERE (`a` = ?))",
			uStr)
	})

	t.Run("order by", func(t *testing.T) {
		u := NewUnion(
			NewSelect("a").AddColumnsAlias("d", "b").From("tableAD").Where(Column("d", ArgString("f"))),
			NewSelect("a", "b").From("tableAB").Where(Column("a", ArgInt64(3))),
		).All().OrderBy("a").OrderByDesc("b")

		uStr, args, err := u.ToSQL()
		assert.NoError(t, err, "%+v", err)
		assert.Exactly(t, []interface{}{"f", int64(3)}, args.Interfaces())
		assert.Exactly(t,
			"(SELECT `a`, `d` AS `b` FROM `tableAD` WHERE (`d` = ?))\nUNION ALL\n(SELECT `a`, `b` FROM `tableAB` WHERE (`a` = ?))\nORDER BY `a` ASC, `b` DESC",
			uStr)
	})

	t.Run("preserve result set", func(t *testing.T) {
		u := NewUnion(
			NewSelect("a").AddColumnsAlias("d", "b").From("tableAD"),
			NewSelect("a", "b").From("tableAB"),
		).All().OrderBy("a").OrderByDesc("b").PreserveResultSet()

		// testing idempotent function ToSQL
		for i := 0; i < 3; i++ {
			uStr, args, err := u.ToSQL()
			assert.NoError(t, err, "%+v", err)
			assert.True(t, args.Interfaces() == nil)
			assert.Exactly(t,
				"(SELECT `a`, `d` AS `b`, 0 AS `_preserve_result_set` FROM `tableAD`)\nUNION ALL\n(SELECT `a`, `b`, 1 AS `_preserve_result_set` FROM `tableAB`)\nORDER BY `_preserve_result_set`, `a` ASC, `b` DESC",
				uStr)
		}
	})
}

func TestUnion_UseBuildCache(t *testing.T) {
	t.Parallel()

	u := NewUnion(
		NewSelect("a").AddColumnsAlias("d", "b").From("tableAD"),
		NewSelect("a", "b").From("tableAB").Where(Column("b", ArgFloat64(3.14159))),
	).All().OrderBy("a").OrderByDesc("b").PreserveResultSet()
	u.UseBuildCache = true

	const cachedSQLPlaceHolder = "(SELECT `a`, `d` AS `b`, 0 AS `_preserve_result_set` FROM `tableAD`)\nUNION ALL\n(SELECT `a`, `b`, 1 AS `_preserve_result_set` FROM `tableAB` WHERE (`b` = ?))\nORDER BY `_preserve_result_set`, `a` ASC, `b` DESC"
	t.Run("without interpolate", func(t *testing.T) {
		for i := 0; i < 3; i++ {
			sql, args, err := u.ToSQL()
			require.NoError(t, err, "%+v", err)
			require.Equal(t, cachedSQLPlaceHolder, sql)
			assert.Equal(t, []interface{}{float64(3.14159)}, args.Interfaces())
			assert.Equal(t, cachedSQLPlaceHolder, string(u.buildCache))
		}
	})

	t.Run("with interpolate", func(t *testing.T) {
		u.Interpolate()
		u.buildCache = nil
		u.RawArguments = nil

		const cachedSQLInterpolated = "(SELECT `a`, `d` AS `b`, 0 AS `_preserve_result_set` FROM `tableAD`)\nUNION ALL\n(SELECT `a`, `b`, 1 AS `_preserve_result_set` FROM `tableAB` WHERE (`b` = 3.14159))\nORDER BY `_preserve_result_set`, `a` ASC, `b` DESC"
		for i := 0; i < 3; i++ {
			sql, args, err := u.ToSQL()
			assert.Equal(t, cachedSQLPlaceHolder, string(u.buildCache))
			require.NoError(t, err, "%+v", err)
			require.Equal(t, cachedSQLInterpolated, sql)
			assert.Nil(t, args)
		}
	})
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

	t.Run("full statement EAV", func(t *testing.T) {
		u := NewUnionTemplate(
			NewSelect().AddColumns("t.value", "t.attribute_id").AddColumnsAlias("t.{column}", "col_type").
				From("catalog_product_entity_{type}", "t").
				Where(Column("entity_id", ArgInt64(1561)), Column("store_id", ArgInt64(1, 0).Operator(In))).
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
			haveArg := args.Interfaces()
			assert.Exactly(t, wantArg, haveArg[:3])
			assert.Exactly(t, wantArg, haveArg[3:6])
			assert.Len(t, haveArg, 15)
			assert.Exactly(t,
				"(SELECT `t`.`value`, `t`.`attribute_id`, `t`.`varcharX` AS `col_type`, 0 AS `_preserve_result_set` FROM `catalog_product_entity_varchar` AS `t` WHERE (`entity_id` = ?) AND (`store_id` IN ?) ORDER BY `t`.`varcharX_store_id` DESC)\nUNION ALL\n(SELECT `t`.`value`, `t`.`attribute_id`, `t`.`intX` AS `col_type`, 1 AS `_preserve_result_set` FROM `catalog_product_entity_int` AS `t` WHERE (`entity_id` = ?) AND (`store_id` IN ?) ORDER BY `t`.`intX_store_id` DESC)\nUNION ALL\n(SELECT `t`.`value`, `t`.`attribute_id`, `t`.`decimalX` AS `col_type`, 2 AS `_preserve_result_set` FROM `catalog_product_entity_decimal` AS `t` WHERE (`entity_id` = ?) AND (`store_id` IN ?) ORDER BY `t`.`decimalX_store_id` DESC)\nUNION ALL\n(SELECT `t`.`value`, `t`.`attribute_id`, `t`.`datetimeX` AS `col_type`, 3 AS `_preserve_result_set` FROM `catalog_product_entity_datetime` AS `t` WHERE (`entity_id` = ?) AND (`store_id` IN ?) ORDER BY `t`.`datetimeX_store_id` DESC)\nUNION ALL\n(SELECT `t`.`value`, `t`.`attribute_id`, `t`.`textX` AS `col_type`, 4 AS `_preserve_result_set` FROM `catalog_product_entity_text` AS `t` WHERE (`entity_id` = ?) AND (`store_id` IN ?) ORDER BY `t`.`textX_store_id` DESC)\nORDER BY `_preserve_result_set`, `col_type` DESC",
				uStr)
		}
	})
	t.Run("StringReplace 2nd call fewer values", func(t *testing.T) {
		u := NewUnionTemplate(
			NewSelect().AddColumns("t.value,t.attribute_id,t.{column} AS `col_type`").From("catalog_product_entity_{type}", "t"),
		).
			StringReplace("{type}", "varchar", "int", "decimal", "datetime", "text").
			StringReplace("{column}", "varcharX", "intX", "decimalX", "datetimeX")

		uStr, args, err := u.ToSQL()
		assert.Empty(t, args.Interfaces())
		assert.Empty(t, uStr)
		assert.True(t, errors.IsNotValid(err), "%+v")
	})
	t.Run("StringReplace 2nd call too many values", func(t *testing.T) {
		u := NewUnionTemplate(
			NewSelect().AddColumns("t.value,t.attribute_id,t.{column} AS `col_type`").From("catalog_product_entity_{type}", "t"),
		).
			StringReplace("{type}", "varchar", "int", "decimal", "datetime", "text").
			StringReplace("{column}", "varcharX", "intX", "decimalX", "datetimeX", "textX", "bytesX")

		uStr, args, err := u.ToSQL()
		assert.Empty(t, args.Interfaces())
		assert.Empty(t, uStr)
		assert.True(t, errors.IsNotValid(err), "%+v")
	})

	t.Run("Preprocessed", func(t *testing.T) {
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

		u := NewUnionTemplate(
			NewSelect().AddColumns("t.value", "t.attribute_id", "t.store_id").From("catalog_product_entity_{type}", "t").
				Where(Column("entity_id", ArgInt64(1561)), Column("store_id", ArgInt64(1, 0).Operator(In))),
		).
			StringReplace("{type}", "varchar", "int", "decimal", "datetime", "text").
			PreserveResultSet().
			All().OrderBy("attribute_id", "store_id").
			Interpolate()
		sqlStr, args, err := u.ToSQL()
		require.NoError(t, err)
		assert.Nil(t, args)
		assert.Exactly(t,
			"(SELECT `t`.`value`, `t`.`attribute_id`, `t`.`store_id`, 0 AS `_preserve_result_set` FROM `catalog_product_entity_varchar` AS `t` WHERE (`entity_id` = 1561) AND (`store_id` IN (1,0)))\nUNION ALL\n(SELECT `t`.`value`, `t`.`attribute_id`, `t`.`store_id`, 1 AS `_preserve_result_set` FROM `catalog_product_entity_int` AS `t` WHERE (`entity_id` = 1561) AND (`store_id` IN (1,0)))\nUNION ALL\n(SELECT `t`.`value`, `t`.`attribute_id`, `t`.`store_id`, 2 AS `_preserve_result_set` FROM `catalog_product_entity_decimal` AS `t` WHERE (`entity_id` = 1561) AND (`store_id` IN (1,0)))\nUNION ALL\n(SELECT `t`.`value`, `t`.`attribute_id`, `t`.`store_id`, 3 AS `_preserve_result_set` FROM `catalog_product_entity_datetime` AS `t` WHERE (`entity_id` = 1561) AND (`store_id` IN (1,0)))\nUNION ALL\n(SELECT `t`.`value`, `t`.`attribute_id`, `t`.`store_id`, 4 AS `_preserve_result_set` FROM `catalog_product_entity_text` AS `t` WHERE (`entity_id` = 1561) AND (`store_id` IN (1,0)))\nORDER BY `_preserve_result_set`, `attribute_id` ASC, `store_id` ASC",
			sqlStr)
	})
}

func TestUnionTemplate_UseBuildCache(t *testing.T) {
	t.Parallel()

	u := NewUnionTemplate(
		NewSelect().AddColumns("t.value", "t.attribute_id", "t.store_id").From("catalog_product_entity_{type}", "t").
			Where(Column("entity_id", ArgInt64(1561)), Column("store_id", ArgInt64(1, 0).Operator(In))),
	).
		StringReplace("{type}", "varchar", "int", "decimal", "datetime", "text").
		PreserveResultSet().
		All().OrderBy("attribute_id", "store_id")
	u.UseBuildCache = true

	const cachedSQLPlaceHolder = "(SELECT `t`.`value`, `t`.`attribute_id`, `t`.`store_id`, 0 AS `_preserve_result_set` FROM `catalog_product_entity_varchar` AS `t` WHERE (`entity_id` = ?) AND (`store_id` IN ?))\nUNION ALL\n(SELECT `t`.`value`, `t`.`attribute_id`, `t`.`store_id`, 1 AS `_preserve_result_set` FROM `catalog_product_entity_int` AS `t` WHERE (`entity_id` = ?) AND (`store_id` IN ?))\nUNION ALL\n(SELECT `t`.`value`, `t`.`attribute_id`, `t`.`store_id`, 2 AS `_preserve_result_set` FROM `catalog_product_entity_decimal` AS `t` WHERE (`entity_id` = ?) AND (`store_id` IN ?))\nUNION ALL\n(SELECT `t`.`value`, `t`.`attribute_id`, `t`.`store_id`, 3 AS `_preserve_result_set` FROM `catalog_product_entity_datetime` AS `t` WHERE (`entity_id` = ?) AND (`store_id` IN ?))\nUNION ALL\n(SELECT `t`.`value`, `t`.`attribute_id`, `t`.`store_id`, 4 AS `_preserve_result_set` FROM `catalog_product_entity_text` AS `t` WHERE (`entity_id` = ?) AND (`store_id` IN ?))\nORDER BY `_preserve_result_set`, `attribute_id` ASC, `store_id` ASC"
	t.Run("without interpolate", func(t *testing.T) {
		for i := 0; i < 3; i++ {
			sql, args, err := u.ToSQL()
			require.NoError(t, err, "%+v", err)
			require.Equal(t, cachedSQLPlaceHolder, sql)
			assert.Equal(t, []interface{}{int64(1561), int64(1), int64(0), int64(1561), int64(1), int64(0), int64(1561), int64(1), int64(0), int64(1561), int64(1), int64(0), int64(1561), int64(1), int64(0)}, args.Interfaces())
			assert.Equal(t, cachedSQLPlaceHolder, string(u.buildCache))
		}
	})

	t.Run("with interpolate", func(t *testing.T) {
		u.Interpolate()
		u.buildCache = nil
		u.RawArguments = nil

		const cachedSQLInterpolated = "(SELECT `t`.`value`, `t`.`attribute_id`, `t`.`store_id`, 0 AS `_preserve_result_set` FROM `catalog_product_entity_varchar` AS `t` WHERE (`entity_id` = 1561) AND (`store_id` IN (1,0)))\nUNION ALL\n(SELECT `t`.`value`, `t`.`attribute_id`, `t`.`store_id`, 1 AS `_preserve_result_set` FROM `catalog_product_entity_int` AS `t` WHERE (`entity_id` = 1561) AND (`store_id` IN (1,0)))\nUNION ALL\n(SELECT `t`.`value`, `t`.`attribute_id`, `t`.`store_id`, 2 AS `_preserve_result_set` FROM `catalog_product_entity_decimal` AS `t` WHERE (`entity_id` = 1561) AND (`store_id` IN (1,0)))\nUNION ALL\n(SELECT `t`.`value`, `t`.`attribute_id`, `t`.`store_id`, 3 AS `_preserve_result_set` FROM `catalog_product_entity_datetime` AS `t` WHERE (`entity_id` = 1561) AND (`store_id` IN (1,0)))\nUNION ALL\n(SELECT `t`.`value`, `t`.`attribute_id`, `t`.`store_id`, 4 AS `_preserve_result_set` FROM `catalog_product_entity_text` AS `t` WHERE (`entity_id` = 1561) AND (`store_id` IN (1,0)))\nORDER BY `_preserve_result_set`, `attribute_id` ASC, `store_id` ASC"
		for i := 0; i < 3; i++ {
			sql, args, err := u.ToSQL()
			assert.Equal(t, cachedSQLPlaceHolder, string(u.buildCache))
			require.NoError(t, err, "%+v", err)
			require.Equal(t, cachedSQLInterpolated, sql)
			assert.Nil(t, args)
		}
	})
}
