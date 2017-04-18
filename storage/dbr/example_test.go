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
	"fmt"

	"github.com/corestoreio/csfw/storage/dbr"
)

func ExampleNewInsert() {
	sqlStr, args, err := dbr.NewInsert("tableA").
		AddColumns("b", "c", "d", "e").
		AddValues(dbr.ArgInt(1), dbr.ArgInt64(2), dbr.ArgString("Three"), dbr.ArgNull()).
		AddValues(dbr.ArgInt(5), dbr.ArgInt64(6), dbr.ArgString("Seven"), dbr.ArgFloat64(3.14156)).
		ToSQL()
	if err != nil {
		fmt.Printf("%+v\n", err)
		return
	}
	fmt.Printf("%s\nArguments: %v\n", sqlStr, args.Interfaces())
	// Output:
	// INSERT INTO `tableA` (`b`,`c`,`d`,`e`) VALUES (?,?,?,?),(?,?,?,?)
	// Arguments: [1 2 Three <nil> 5 6 Seven 3.14156]
}

func ExampleInsert_AddOnDuplicateKey() {
	sqlStr, args, err := dbr.NewInsert("dbr_people").
		AddColumns("id", "name", "email").
		AddValues(dbr.ArgInt64(1), dbr.ArgStrings("Pik'e"), dbr.ArgStrings("pikes@peak.com")).
		AddOnDuplicateKey("name", dbr.ArgStrings("Pik3")).
		AddOnDuplicateKey("email", nil).
		ToSQL()
	if err != nil {
		fmt.Printf("%+v\n", err)
		return
	}

	sqlPre, err := dbr.Preprocess(sqlStr, args...)
	if err != nil {
		fmt.Printf("%+v\n", err)
		return
	}
	fmt.Printf("%s\nArguments: %v\nProcessed: %s\n", sqlStr, args.Interfaces(), sqlPre)

	// Output:
	// INSERT INTO `dbr_people` (`id`,`name`,`email`) VALUES (?,?,?) ON DUPLICATE KEY UPDATE `name`=?, `email`=VALUES(`email`)
	// Arguments: [1 Pik'e pikes@peak.com Pik3]
	// Processed: INSERT INTO `dbr_people` (`id`,`name`,`email`) VALUES (1,'Pik\'e','pikes@peak.com') ON DUPLICATE KEY UPDATE `name`='Pik3', `email`=VALUES(`email`)
}

func ExampleInsert_FromSelect() {
	ins := dbr.NewInsert("tableA")

	argEq := dbr.Eq{"int64B": dbr.ArgInt64(1, 2, 3).Operator(dbr.OperatorIn)}

	sqlStr, args, err := ins.FromSelect(dbr.NewSelect().AddColumnsQuoted("something_id,user_id,other").
		From("some_table").
		Where(dbr.Condition("int64A = ? OR string = ?", dbr.ArgInt64(1), dbr.ArgStrings("wat"))).
		Where(argEq).
		OrderByDesc("id").
		Paginate(1, 20))
	if err != nil {
		fmt.Printf("%+v\n", err)
		return
	}

	sqlPre, err := dbr.Preprocess(sqlStr, args...)
	if err != nil {
		fmt.Printf("%+v\n", err)
		return
	}

	fmt.Printf("%s\nArguments: %v\nProcessed: %s\n", sqlStr, args.Interfaces(), sqlPre)
	// Output:
	// INSERT INTO `tableA` SELECT `something_id`, `user_id`, `other` FROM `some_table` WHERE (int64A = ? OR string = ?) AND (`int64B` IN ?) ORDER BY id DESC LIMIT 20 OFFSET 0
	// Arguments: [1 wat 1 2 3]
	// Processed: INSERT INTO `tableA` SELECT `something_id`, `user_id`, `other` FROM `some_table` WHERE (int64A = 1 OR string = 'wat') AND (`int64B` IN (1,2,3)) ORDER BY id DESC LIMIT 20 OFFSET 0
}

func ExampleNewDelete() {
	sqlStr, args, err := dbr.NewDelete("tableA").Where(
		dbr.Condition("a", dbr.ArgStrings("b'%").Operator(dbr.OperatorLike)),
		dbr.Condition("b", dbr.ArgInt(3, 4, 5, 6).Operator(dbr.OperatorIn)),
	).
		Limit(1).OrderBy("id").
		ToSQL()
	if err != nil {
		fmt.Printf("%+v\n", err)
		return
	}

	sqlPre, err := dbr.Preprocess(sqlStr, args...)
	if err != nil {
		fmt.Printf("%+v\n", err)
		return
	}

	fmt.Printf("%s\nArguments: %v\nProcessed: %s\n", sqlStr, args.Interfaces(), sqlPre)
	// Output:
	// DELETE FROM `tableA` WHERE (`a` LIKE ?) AND (`b` IN ?) ORDER BY id LIMIT 1
	// Arguments: [b'% 3 4 5 6]
	// Processed: DELETE FROM `tableA` WHERE (`a` LIKE 'b\'%') AND (`b` IN (3,4,5,6)) ORDER BY id LIMIT 1
}

// ExampleNewUnion constructs a UNION with three SELECTs. It preserves the
// results sets of each SELECT by simply adding an internal index to the columns
// list and sort ascending with the internal index.
func ExampleNewUnion() {

	u := dbr.NewUnion(
		dbr.NewSelect().AddColumnsQuotedAlias("a1", "A", "a2", "B").From("tableA").Where(dbr.Condition("a1", dbr.ArgInt64(3))),
		dbr.NewSelect().AddColumnsQuotedAlias("b1", "A", "b2", "B").From("tableB").Where(dbr.Condition("b1", dbr.ArgInt64(4))),
	)
	// Maybe more of your code ...
	u.Append(
		dbr.NewSelect().AddColumnsExprAlias("concat(c1,'-',c2)", "A").
			AddColumnsQuotedAlias("c2", "B").
			From("tableC").Where(dbr.Condition("c2", dbr.ArgString("ArgForC2"))),
	).
		OrderBy("A").       // Ascending by A
		OrderByDesc("B").   // Descending by B
		All().              // Enables UNION ALL syntax
		PreserveResultSet() // Maintains the correct order of the result set for all SELECTs.
	// Note that the final ORDER BY statement of a UNION creates a temporary
	// table in MySQL.
	sqlStr, args, err := u.ToSQL()
	if err != nil {
		fmt.Printf("%+v\n", err)
		return
	}

	fmt.Printf("%s\nArguments: %v\n", sqlStr, args.Interfaces())
	// Output:
	// (SELECT `a1` AS `A`, `a2` AS `B`, 0 AS `_preserve_result_set` FROM `tableA` WHERE (`a1` = ?))
	// UNION ALL
	// (SELECT `b1` AS `A`, `b2` AS `B`, 1 AS `_preserve_result_set` FROM `tableB` WHERE (`b1` = ?))
	// UNION ALL
	// (SELECT concat(c1,'-',c2) AS `A`, `c2` AS `B`, 2 AS `_preserve_result_set` FROM `tableC` WHERE (`c2` = ?))
	// ORDER BY `_preserve_result_set`, A, B DESC
	// Arguments: [3 4 ArgForC2]
}

func ExampleNewUnionTemplate() {

	u := dbr.NewUnionTemplate(
		dbr.NewSelect().AddColumnsQuoted("t.value,t.attribute_id,t.store_id").From("catalog_product_entity_{type}", "t").
			Where(dbr.Condition("entity_id", dbr.ArgInt64(1561)), dbr.Condition("store_id", dbr.ArgInt64(1, 0).Operator(dbr.OperatorIn))),
	).
		StringReplace("{type}", "varchar", "int", "decimal", "datetime", "text").
		PreserveResultSet().
		All().OrderBy("attribute_id", "store_id")

	sqlStr, args, err := u.ToSQL()
	if err != nil {
		fmt.Printf("%+v\n", err)
		return
	}

	fmt.Printf("%s\nArguments: %v\n", sqlStr, args.Interfaces())
	// Output:
	// (SELECT `t`.`value`, `t`.`attribute_id`, `t`.`store_id`, 0 AS `_preserve_result_set` FROM `catalog_product_entity_varchar` AS `t` WHERE (`entity_id` = ?) AND (`store_id` IN ?))
	// UNION ALL
	// (SELECT `t`.`value`, `t`.`attribute_id`, `t`.`store_id`, 1 AS `_preserve_result_set` FROM `catalog_product_entity_int` AS `t` WHERE (`entity_id` = ?) AND (`store_id` IN ?))
	// UNION ALL
	// (SELECT `t`.`value`, `t`.`attribute_id`, `t`.`store_id`, 2 AS `_preserve_result_set` FROM `catalog_product_entity_decimal` AS `t` WHERE (`entity_id` = ?) AND (`store_id` IN ?))
	// UNION ALL
	// (SELECT `t`.`value`, `t`.`attribute_id`, `t`.`store_id`, 3 AS `_preserve_result_set` FROM `catalog_product_entity_datetime` AS `t` WHERE (`entity_id` = ?) AND (`store_id` IN ?))
	// UNION ALL
	// (SELECT `t`.`value`, `t`.`attribute_id`, `t`.`store_id`, 4 AS `_preserve_result_set` FROM `catalog_product_entity_text` AS `t` WHERE (`entity_id` = ?) AND (`store_id` IN ?))
	// ORDER BY `_preserve_result_set`, attribute_id, store_id
	// Arguments: [1561 1 0 1561 1 0 1561 1 0 1561 1 0 1561 1 0]
}

// ExampleUnionTemplate_Preprocess interpolates the SQL string with its
// placeholders and puts for each placeholder the correct encoded and escaped
// value into it. Eliminates the need for prepared statements by sending in one
// round trip the query and its arguments directly to the database server. If
// you execute a query multiple times within a short time you should use
// prepared statements.
func ExampleUnionTemplate_Preprocess() {

	u := dbr.NewUnionTemplate(
		dbr.NewSelect().AddColumnsQuoted("t.value,t.attribute_id,t.store_id").From("catalog_product_entity_{type}", "t").
			Where(dbr.Condition("entity_id", dbr.ArgInt64(1561)), dbr.Condition("store_id", dbr.ArgInt64(1, 0).Operator(dbr.OperatorIn))),
	).
		StringReplace("{type}", "varchar", "int", "decimal", "datetime", "text").
		PreserveResultSet().
		All().OrderBy("attribute_id", "store_id")

	sqlStr, err := u.Preprocess()
	if err != nil {
		fmt.Printf("%+v\n", err)
		return
	}

	fmt.Printf("%s\n", sqlStr)
	// Output:
	// (SELECT `t`.`value`, `t`.`attribute_id`, `t`.`store_id`, 0 AS `_preserve_result_set` FROM `catalog_product_entity_varchar` AS `t` WHERE (`entity_id` = 1561) AND (`store_id` IN (1,0)))
	// UNION ALL
	// (SELECT `t`.`value`, `t`.`attribute_id`, `t`.`store_id`, 1 AS `_preserve_result_set` FROM `catalog_product_entity_int` AS `t` WHERE (`entity_id` = 1561) AND (`store_id` IN (1,0)))
	// UNION ALL
	// (SELECT `t`.`value`, `t`.`attribute_id`, `t`.`store_id`, 2 AS `_preserve_result_set` FROM `catalog_product_entity_decimal` AS `t` WHERE (`entity_id` = 1561) AND (`store_id` IN (1,0)))
	// UNION ALL
	// (SELECT `t`.`value`, `t`.`attribute_id`, `t`.`store_id`, 3 AS `_preserve_result_set` FROM `catalog_product_entity_datetime` AS `t` WHERE (`entity_id` = 1561) AND (`store_id` IN (1,0)))
	// UNION ALL
	// (SELECT `t`.`value`, `t`.`attribute_id`, `t`.`store_id`, 4 AS `_preserve_result_set` FROM `catalog_product_entity_text` AS `t` WHERE (`entity_id` = 1561) AND (`store_id` IN (1,0)))
	// ORDER BY `_preserve_result_set`, attribute_id, store_id
}

func ExamplePreprocess() {
	sqlStr, err := dbr.Preprocess("SELECT * FROM x WHERE a IN ? AND b IN ? AND c NOT IN ? AND d BETWEEN ? AND ?",
		dbr.ArgInt(1).Operator(dbr.OperatorIn),
		dbr.ArgInt(1, 2, 3).Operator(dbr.OperatorIn),
		dbr.ArgInt64(5, 6, 7).Operator(dbr.OperatorIn),
		dbr.ArgStrings("wat", "ok").Operator(dbr.OperatorBetween),
	)
	if err != nil {
		fmt.Printf("%+v\n", err)
		return
	}
	fmt.Printf("%s\n", sqlStr)
	// Output:
	// SELECT * FROM x WHERE a IN 1 AND b IN (1,2,3) AND c NOT IN (5,6,7) AND d BETWEEN 'wat' AND 'ok'
}

func ExampleRepeat() {
	sl := []string{"a", "b", "c", "d", "e"}

	sqlStr, args, err := dbr.Repeat("SELECT * FROM `table` WHERE id IN (?) AND name IN (?)",
		dbr.ArgInt(5, 7, 9), dbr.ArgStrings(sl...))

	if err != nil {
		fmt.Printf("%+v\n", err)
		return
	}

	fmt.Printf("%s\nArguments: %v\n", sqlStr, args)
	// Output:
	// SELECT * FROM `table` WHERE id IN (?,?,?) AND name IN (?,?,?,?,?)
	// Arguments: [5 7 9 a b c d e]
}

func ExampleCondition() {
	argPrinter := func(arg ...dbr.Argument) {
		sqlStr, args, err := dbr.NewSelect().AddColumnsQuoted("a", "b").
			From("c").Where(dbr.Condition("d", arg...)).ToSQL()
		if err != nil {
			fmt.Printf("%+v\n", err)
		} else {
			fmt.Printf("%s\n", sqlStr)
			if len(args) > 0 {
				fmt.Printf("Arguments: %v\n", args)
			}
		}
	}

	//OperatorNull       byte = 'n' // IS NULL
	//OperatorNotNull    byte = 'N' // IS NOT NULL
	//OperatorIn         byte = 'i' // IN ?
	//OperatorNotIn      byte = 'I' // NOT IN ?
	//OperatorBetween    byte = 'b' // BETWEEN ? AND ?
	//OperatorNotBetween byte = 'B' // NOT BETWEEN ? AND ?
	//OperatorLike       byte = 'l' // LIKE ?
	//OperatorNotLike    byte = 'L' // NOT LIKE ?
	//OperatorGreatest   byte = 'g' // GREATEST(?,?,?)
	//OperatorLeast      byte = 'a' // LEAST(?,?,?)
	//OperatorEqual      byte = '=' // = ?
	//OperatorNotEqual   byte = '!' // != ?
	//OperatorExists     byte = 'e' // EXISTS(subquery)
	//OperatorNotExists  byte = 'E' // NOT EXISTS(subquery)

	argPrinter(dbr.ArgNull())
	argPrinter(dbr.ArgNotNull())

	argPrinter(dbr.ArgInt(3).Operator(dbr.OperatorNull))
	argPrinter(dbr.ArgInt(4).Operator(dbr.OperatorNotNull))
	argPrinter(dbr.ArgInt(5).Operator(dbr.OperatorEqual))
	argPrinter(dbr.ArgInt(6).Operator(dbr.OperatorNotEqual))

	// Output:
	// SELECT `a`, `b` FROM `c` WHERE (`d` IS NULL)
	// SELECT `a`, `b` FROM `c` WHERE (`d` IS NOT NULL)
}
