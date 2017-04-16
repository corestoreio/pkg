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
	"testing"

	"github.com/corestoreio/csfw/storage/dbr"
	"github.com/corestoreio/errors"
	"github.com/stretchr/testify/assert"
)

func TestUnion(t *testing.T) {
	t.Parallel()
	t.Run("simple", func(t *testing.T) {
		u := dbr.NewUnion(
			dbr.NewSelect("a", "b").From("tableAB").Where(dbr.Condition("a", dbr.ArgInt64(3))),
		)
		u.Append(
			dbr.NewSelect("c", "d").From("tableCD").Where(dbr.Condition("d", dbr.ArgString("e"))),
		)

		uStr, args, err := u.ToSQL()
		assert.NoError(t, err, "%+v", err)
		assert.Exactly(t, []interface{}{int64(3), "e"}, args.Interfaces())
		assert.Exactly(t,
			"(SELECT a, b FROM `tableAB` WHERE (`a` = ?)) UNION (SELECT c, d FROM `tableCD` WHERE (`d` = ?))",
			uStr)
	})

	t.Run("simple all", func(t *testing.T) {
		u := dbr.NewUnion(
			dbr.NewSelect("c", "d").From("tableCD").Where(dbr.Condition("d", dbr.ArgString("e"))),
			dbr.NewSelect("a", "b").From("tableAB").Where(dbr.Condition("a", dbr.ArgInt64(3))),
		).All()

		uStr, args, err := u.ToSQL()
		assert.NoError(t, err, "%+v", err)
		assert.Exactly(t, []interface{}{"e", int64(3)}, args.Interfaces())
		assert.Exactly(t,
			"(SELECT c, d FROM `tableCD` WHERE (`d` = ?)) UNION ALL (SELECT a, b FROM `tableAB` WHERE (`a` = ?))",
			uStr)
	})

	t.Run("order by", func(t *testing.T) {
		u := dbr.NewUnion(
			dbr.NewSelect("a").AddColumnsQuotedAlias("d", "b").From("tableAD").Where(dbr.Condition("d", dbr.ArgString("f"))),
			dbr.NewSelect("a", "b").From("tableAB").Where(dbr.Condition("a", dbr.ArgInt64(3))),
		).All().OrderBy("a").OrderDir("b", false)

		uStr, args, err := u.ToSQL()
		assert.NoError(t, err, "%+v", err)
		assert.Exactly(t, []interface{}{"f", int64(3)}, args.Interfaces())
		assert.Exactly(t,
			"(SELECT a, `d` AS `b` FROM `tableAD` WHERE (`d` = ?)) UNION ALL (SELECT a, b FROM `tableAB` WHERE (`a` = ?)) ORDER BY a, b DESC",
			uStr)
	})

	t.Run("preserve result set", func(t *testing.T) {
		u := dbr.NewUnion(
			dbr.NewSelect("a").AddColumnsQuotedAlias("d", "b").From("tableAD"),
			dbr.NewSelect("a", "b").From("tableAB"),
		).All().OrderBy("a").OrderDir("b", false).PreserveResultSet()

		// testing idempotent function ToSQL
		for i := 0; i < 3; i++ {
			uStr, args, err := u.ToSQL()
			assert.NoError(t, err, "%+v", err)
			assert.True(t, args.Interfaces() == nil)
			assert.Exactly(t,
				"(SELECT a, `d` AS `b`, 0 AS `_preserve_result_set` FROM `tableAD`) UNION ALL (SELECT a, b, 1 AS `_preserve_result_set` FROM `tableAB`) ORDER BY `_preserve_result_set`, a, b DESC",
				uStr)
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
		u := dbr.NewUnionTemplate(
			dbr.NewSelect().AddColumnsQuoted("t.value,t.attribute_id,t.{column} AS `col_type`").From("catalog_product_entity_{type}", "t").
				Where(dbr.Condition("entity_id", dbr.ArgInt64(1561)), dbr.Condition("store_id", dbr.ArgInt64(1, 0).Operator(dbr.OperatorIn))).
				OrderDir("t.{column}_store_id", false),
		).
			StringReplace("{type}", "varchar", "int", "decimal", "datetime", "text").
			StringReplace("{column}", "varcharX", "intX", "decimalX", "datetimeX", "textX").
			PreserveResultSet().
			All().
			OrderDir("col_type", false)

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
				"(SELECT `t`.`value`, `t`.`attribute_id`, t.varcharX AS `col_type`, 0 AS `_preserve_result_set` FROM `catalog_product_entity_varchar` AS `t` WHERE (`entity_id` = ?) AND (`store_id` IN ?) ORDER BY t.varcharX_store_id DESC) UNION ALL (SELECT `t`.`value`, `t`.`attribute_id`, t.intX AS `col_type`, 1 AS `_preserve_result_set` FROM `catalog_product_entity_int` AS `t` WHERE (`entity_id` = ?) AND (`store_id` IN ?) ORDER BY t.intX_store_id DESC) UNION ALL (SELECT `t`.`value`, `t`.`attribute_id`, t.decimalX AS `col_type`, 2 AS `_preserve_result_set` FROM `catalog_product_entity_decimal` AS `t` WHERE (`entity_id` = ?) AND (`store_id` IN ?) ORDER BY t.decimalX_store_id DESC) UNION ALL (SELECT `t`.`value`, `t`.`attribute_id`, t.datetimeX AS `col_type`, 3 AS `_preserve_result_set` FROM `catalog_product_entity_datetime` AS `t` WHERE (`entity_id` = ?) AND (`store_id` IN ?) ORDER BY t.datetimeX_store_id DESC) UNION ALL (SELECT `t`.`value`, `t`.`attribute_id`, t.textX AS `col_type`, 4 AS `_preserve_result_set` FROM `catalog_product_entity_text` AS `t` WHERE (`entity_id` = ?) AND (`store_id` IN ?) ORDER BY t.textX_store_id DESC) ORDER BY `_preserve_result_set`, col_type DESC",
				uStr)
		}
	})
	t.Run("StringReplace 2nd call fewer values", func(t *testing.T) {
		u := dbr.NewUnionTemplate(
			dbr.NewSelect().AddColumnsQuoted("t.value,t.attribute_id,t.{column} AS `col_type`").From("catalog_product_entity_{type}", "t"),
		).
			StringReplace("{type}", "varchar", "int", "decimal", "datetime", "text").
			StringReplace("{column}", "varcharX", "intX", "decimalX", "datetimeX")

		uStr, args, err := u.ToSQL()
		assert.Empty(t, args.Interfaces())
		assert.Empty(t, uStr)
		assert.True(t, errors.IsNotValid(err), "%+v")
	})
	t.Run("StringReplace 2nd call too many values", func(t *testing.T) {
		u := dbr.NewUnionTemplate(
			dbr.NewSelect().AddColumnsQuoted("t.value,t.attribute_id,t.{column} AS `col_type`").From("catalog_product_entity_{type}", "t"),
		).
			StringReplace("{type}", "varchar", "int", "decimal", "datetime", "text").
			StringReplace("{column}", "varcharX", "intX", "decimalX", "datetimeX", "textX", "bytesX")

		uStr, args, err := u.ToSQL()
		assert.Empty(t, args.Interfaces())
		assert.Empty(t, uStr)
		assert.True(t, errors.IsNotValid(err), "%+v")
	})
}

var benchmarkGlobalArgs dbr.Arguments

// BenchmarkUnion_AllOptions-4           	   10000	   1207331 ns/op	  950890 B/op	      20 allocs/op
func BenchmarkUnion_AllOptions(b *testing.B) {

	u := dbr.NewUnion(
		dbr.NewSelect().AddColumnsQuoted("t.value,t.attribute_id,t.varchar AS `col_type`").From("catalog_product_entity_varchar", "t").
			Where(dbr.Condition("entity_id", dbr.ArgInt64(1561)), dbr.Condition("store_id", dbr.ArgInt64(1, 0).Operator(dbr.OperatorIn))).
			OrderDir("t.varchar_store_id", false),
		dbr.NewSelect().AddColumnsQuoted("t.value,t.attribute_id,t.int AS `col_type`").From("catalog_product_entity_int", "t").
			Where(dbr.Condition("entity_id", dbr.ArgInt64(1561)), dbr.Condition("store_id", dbr.ArgInt64(1, 0).Operator(dbr.OperatorIn))).
			OrderDir("t.int_store_id", false),
		dbr.NewSelect().AddColumnsQuoted("t.value,t.attribute_id,t.decimal AS `col_type`").From("catalog_product_entity_decimal", "t").
			Where(dbr.Condition("entity_id", dbr.ArgInt64(1561)), dbr.Condition("store_id", dbr.ArgInt64(1, 0).Operator(dbr.OperatorIn))).
			OrderDir("t.decimal_store_id", false),
		dbr.NewSelect().AddColumnsQuoted("t.value,t.attribute_id,t.datetime AS `col_type`").From("catalog_product_entity_datetime", "t").
			Where(dbr.Condition("entity_id", dbr.ArgInt64(1561)), dbr.Condition("store_id", dbr.ArgInt64(1, 0).Operator(dbr.OperatorIn))).
			OrderDir("t.datetime_store_id", false),
		dbr.NewSelect().AddColumnsQuoted("t.value,t.attribute_id,t.text AS `col_type`").From("catalog_product_entity_text", "t").
			Where(dbr.Condition("entity_id", dbr.ArgInt64(1561)), dbr.Condition("store_id", dbr.ArgInt64(1, 0).Operator(dbr.OperatorIn))).
			OrderDir("t.text_store_id", false),
	).All().OrderBy("a").OrderDir("b", false).PreserveResultSet()

	for i := 0; i < b.N; i++ {
		_, args, err := u.ToSQL()
		if err != nil {
			b.Fatalf("%+v", err)
		}
		benchmarkGlobalArgs = args
	}
}

// BenchmarkUnionTemplate_AllOptions-4   	  100000	     15432 ns/op	    7793 B/op	      59 allocs/op
func BenchmarkUnionTemplate_AllOptions(b *testing.B) {
	u := dbr.NewUnionTemplate(
		dbr.NewSelect().AddColumnsQuoted("t.value,t.attribute_id,t.{column} AS `col_type`").From("catalog_product_entity_{type}", "t").
			Where(dbr.Condition("entity_id", dbr.ArgInt64(1561)), dbr.Condition("store_id", dbr.ArgInt64(1, 0).Operator(dbr.OperatorIn))).
			OrderDir("t.{column}_store_id", false),
	).
		StringReplace("{type}", "varchar", "int", "decimal", "datetime", "text").
		StringReplace("{column}", "varcharX", "intX", "decimalX", "datetimeX", "textX").
		PreserveResultSet().
		All().
		OrderDir("col_type", false)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, args, err := u.ToSQL()
		if err != nil {
			b.Fatalf("%+v", err)
		}
		benchmarkGlobalArgs = args
	}
}
