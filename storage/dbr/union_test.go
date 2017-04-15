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

		uStr, args, err := u.ToSQL()
		assert.NoError(t, err, "%+v", err)
		assert.True(t, args.Interfaces() == nil)
		assert.Exactly(t,
			"(SELECT a, `d` AS `b`, 0 AS `_preserve_result_set` FROM `tableAD`) UNION ALL (SELECT a, b, 1 AS `_preserve_result_set` FROM `tableAB`) ORDER BY `_preserve_result_set`, a, b DESC",
			uStr)
	})

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
	//t.Run("complex", func(t *testing.T) {
	//
	//	u := dbr.NewUnion(
	//		dbr.NewSelect().AddColumnsQuoted("t.value,t.attribute_id").From("catalog_product_entity_varchar", "t").
	//			Where(dbr.Condition("entity_id", dbr.ArgInt64(1561)), dbr.Condition("store_id", dbr.ArgInt64(1, 0).Operator(dbr.OperatorIn))).
	//			OrderDir("t.store_id", false),
	//	).All()
	//
	//	uStr, args, err := u.ToSQL()
	//	assert.NoError(t, err, "%+v", err)
	//	assert.True(t, args.Interfaces() == nil)
	//	assert.Exactly(t,
	//		"(SELECT a, `d` AS `b`, 0 AS `_preserve_result_set` FROM `tableAD`) UNION ALL (SELECT a, b, 1 AS `_preserve_result_set` FROM `tableAB`) ORDER BY `_preserve_result_set`, a, b DESC",
	//		uStr)
	//
	//})
}
