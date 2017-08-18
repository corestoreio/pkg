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
	"os"
	"strings"
	"time"

	"github.com/corestoreio/csfw/storage/dbr"
	"github.com/corestoreio/csfw/util/wordwrap"
	"github.com/corestoreio/errors"
)

// iFaceToArgs unpacks the interface and creates an Argument slice. Just a
// helper function for the examples.
func iFaceToArgs(values ...interface{}) dbr.Arguments {
	args := make(dbr.Arguments, 0, len(values))
	for _, val := range values {
		switch v := val.(type) {
		case float64:
			args = args.Float64(v)
		case int64:
			args = args.Int64(v)
		case int:
			args = args.Int(v)
			args = args.Int(v)
		case bool:
			args = args.Bool(v)
		case string:
			args = args.Str(v)
		case []byte:
			args = args.Bytes(v)
		case time.Time:
			args = args.Time(v)
		case *time.Time:
			if v != nil {
				args = args.Time(*v)
			}
		case nil:
			args = args.Null()
		default:
			panic(errors.NewNotSupportedf("[dbr] iFaceToArgs type %#v not yet supported", v))
		}
	}
	return args
}

func writeToSQLAndInterpolate(qb dbr.QueryBuilder) {
	sqlStr, args, err := qb.ToSQL()
	if err != nil {
		fmt.Printf("%+v\n", err)
		return
	}
	fmt.Println("Prepared Statement:")
	wordwrap.Fstring(os.Stdout, sqlStr, 80)
	fmt.Print("\n")
	if len(args) > 0 {
		fmt.Printf("Arguments: %v\n\n", args)
	}
	if len(args) == 0 {
		return
	}
	// iFaceToArgs is a hacky way, but works for examples, not in production!
	sqlStr = dbr.Interpolate(sqlStr).ArgUnions(iFaceToArgs(args...)).String()
	fmt.Println("Interpolated Statement:")
	wordwrap.Fstring(os.Stdout, sqlStr, 80)
}

func ExampleNewInsert() {
	i := dbr.NewInsert("tableA").
		AddColumns("b", "c", "d", "e").
		AddValues(1, 2, "Three", nil).
		AddValues(5, 6, "Seven", 3.14156)
	writeToSQLAndInterpolate(i)

	// Output:
	//Prepared Statement:
	//INSERT INTO `tableA` (`b`,`c`,`d`,`e`) VALUES (?,?,?,?),(?,?,?,?)
	//Arguments: [1 2 Three <nil> 5 6 Seven 3.14156]
	//
	//Interpolated Statement:
	//INSERT INTO `tableA` (`b`,`c`,`d`,`e`) VALUES
	//(1,2,'Three',NULL),(5,6,'Seven',3.14156)
}

func ExampleNewInsert_withoutColumns() {
	i := dbr.NewInsert("catalog_product_link").
		AddValues(2046, 33, 3).
		AddValues(2046, 34, 3).
		AddValues(2046, 35, 3)
	writeToSQLAndInterpolate(i)

	// Output:
	//Prepared Statement:
	//INSERT INTO `catalog_product_link` VALUES (?,?,?),(?,?,?),(?,?,?)
	//Arguments: [2046 33 3 2046 34 3 2046 35 3]
	//
	//Interpolated Statement:
	//INSERT INTO `catalog_product_link` VALUES (2046,33,3),(2046,34,3),(2046,35,3)
}

func ExampleInsert_AddValues() {
	// Without any columns you must for each row call AddValues. Here we insert
	// three rows at once.
	i := dbr.NewInsert("catalog_product_link").
		AddValues(2046, 33, 3).
		AddValues(2046, 34, 3).
		AddValues(2046, 35, 3)
	writeToSQLAndInterpolate(i)
	fmt.Print("\n\n")

	// Specifying columns allows to call only one time AddValues but inserting
	// three rows at once. Of course you can also insert only one row ;-)
	i = dbr.NewInsert("catalog_product_link").
		AddColumns("product_id", "linked_product_id", "link_type_id").
		AddValues(
			2046, 33, 3,
			2046, 34, 3,
			2046, 35, 3,
		)
	writeToSQLAndInterpolate(i)

	// Output:
	//Prepared Statement:
	//INSERT INTO `catalog_product_link` VALUES (?,?,?),(?,?,?),(?,?,?)
	//Arguments: [2046 33 3 2046 34 3 2046 35 3]
	//
	//Interpolated Statement:
	//INSERT INTO `catalog_product_link` VALUES (2046,33,3),(2046,34,3),(2046,35,3)
	//
	//Prepared Statement:
	//INSERT INTO `catalog_product_link`
	//(`product_id`,`linked_product_id`,`link_type_id`) VALUES (?,?,?),(?,?,?),(?,?,?)
	//Arguments: [2046 33 3 2046 34 3 2046 35 3]
	//
	//Interpolated Statement:
	//INSERT INTO `catalog_product_link`
	//(`product_id`,`linked_product_id`,`link_type_id`) VALUES
	//(2046,33,3),(2046,34,3),(2046,35,3)
}

func ExampleInsert_AddOnDuplicateKey() {
	i := dbr.NewInsert("dbr_people").
		AddColumns("id", "name", "email").
		AddValues(1, "Pik'e", "pikes@peak.com").
		AddOnDuplicateKey(
			dbr.Column("name").Str("Pik3"),
			dbr.Column("email").Values(),
		)
	writeToSQLAndInterpolate(i)

	// Output:
	//Prepared Statement:
	//INSERT INTO `dbr_people` (`id`,`name`,`email`) VALUES (?,?,?) ON DUPLICATE KEY
	//UPDATE `name`=?, `email`=VALUES(`email`)
	//Arguments: [1 Pik'e pikes@peak.com Pik3]
	//
	//Interpolated Statement:
	//INSERT INTO `dbr_people` (`id`,`name`,`email`) VALUES
	//(1,'Pik\'e','pikes@peak.com') ON DUPLICATE KEY UPDATE `name`='Pik3',
	//`email`=VALUES(`email`)
}

func ExampleInsert_SetRowCount() {
	// RowCount of 4 allows to insert four rows with a single INSERT query.
	// Useful when creating prepared statements.
	i := dbr.NewInsert("dbr_people").AddColumns("id", "name", "email").SetRowCount(4)
	writeToSQLAndInterpolate(i)

	// Output:
	//Prepared Statement:
	//INSERT INTO `dbr_people` (`id`,`name`,`email`) VALUES
	//(?,?,?),(?,?,?),(?,?,?),(?,?,?)
}

func ExampleInsert_FromSelect() {
	ins := dbr.NewInsert("tableA")

	ins.FromSelect(
		dbr.NewSelect().AddColumns("something_id", "user_id").
			AddColumns("other").
			From("some_table").
			Where(
				dbr.ParenthesisOpen(),
				dbr.Column("int64A").GreaterOrEqual().Int64(1),
				dbr.Column("string").Str("wat").Or(),
				dbr.ParenthesisClose(),
				dbr.Column("int64B").In().Int64s(1, 2, 3),
			).
			OrderByDesc("id").
			Paginate(1, 20),
	)
	writeToSQLAndInterpolate(ins)
	// Output:
	//Prepared Statement:
	//INSERT INTO `tableA` SELECT `something_id`, `user_id`, `other` FROM `some_table`
	//WHERE ((`int64A` >= ?) OR (`string` = ?)) AND (`int64B` IN (?,?,?)) ORDER BY
	//`id` DESC LIMIT 20 OFFSET 0
	//Arguments: [1 wat 1 2 3]
	//
	//Interpolated Statement:
	//INSERT INTO `tableA` SELECT `something_id`, `user_id`, `other` FROM `some_table`
	//WHERE ((`int64A` >= 1) OR (`string` = 'wat')) AND (`int64B` IN (1,2,3)) ORDER BY
	//`id` DESC LIMIT 20 OFFSET 0
}

func ExampleInsert_Pair() {
	ins := dbr.NewInsert("catalog_product_link").
		Pair(
			// First row
			dbr.Column("product_id").Int64(2046),
			dbr.Column("linked_product_id").Int64(33),
			dbr.Column("link_type_id").Int64(3),

			// second row
			dbr.Column("product_id").Int64(2046),
			dbr.Column("linked_product_id").Int64(34),
			dbr.Column("link_type_id").Int64(3),
		)
	writeToSQLAndInterpolate(ins)
	// Output:
	//Prepared Statement:
	//INSERT INTO `catalog_product_link`
	//(`product_id`,`linked_product_id`,`link_type_id`) VALUES (?,?,?),(?,?,?)
	//Arguments: [2046 33 3 2046 34 3]
	//
	//Interpolated Statement:
	//INSERT INTO `catalog_product_link`
	//(`product_id`,`linked_product_id`,`link_type_id`) VALUES (2046,33,3),(2046,34,3)
}

func ExampleNewDelete() {
	d := dbr.NewDelete("tableA").Where(
		dbr.Column("a").Like().Str("b'%"),
		dbr.Column("b").In().Ints(3, 4, 5, 6),
	).
		Limit(1).OrderBy("id")
	writeToSQLAndInterpolate(d)
	// Output:
	//Prepared Statement:
	//DELETE FROM `tableA` WHERE (`a` LIKE ?) AND (`b` IN (?,?,?,?)) ORDER BY `id`
	//LIMIT 1
	//Arguments: [b'% 3 4 5 6]
	//
	//Interpolated Statement:
	//DELETE FROM `tableA` WHERE (`a` LIKE 'b\'%') AND (`b` IN (3,4,5,6)) ORDER BY
	//`id` LIMIT 1
}

// ExampleNewUnion constructs a UNION with three SELECTs. It preserves the
// results sets of each SELECT by simply adding an internal index to the columns
// list and sort ascending with the internal index.
func ExampleNewUnion() {

	u := dbr.NewUnion(
		dbr.NewSelect().AddColumnsAliases("a1", "A", "a2", "B").From("tableA").Where(dbr.Column("a1").Int64(3)),
		dbr.NewSelect().AddColumnsAliases("b1", "A", "b2", "B").From("tableB").Where(dbr.Column("b1").Int64(4)),
	)
	// Maybe more of your code ...
	u.Append(
		dbr.NewSelect().AddColumnsConditions(
			dbr.Expr("concat(c1,?,c2)").Alias("A").Str("-"),
		).
			AddColumnsAliases("c2", "B").
			From("tableC").Where(dbr.Column("c2").Str("ArgForC2")),
	).
		OrderBy("A").       // Ascending by A
		OrderByDesc("B").   // Descending by B
		All().              // Enables UNION ALL syntax
		PreserveResultSet() // Maintains the correct order of the result set for all SELECTs.
	// Note that the final ORDER BY statement of a UNION creates a temporary
	// table in MySQL.
	writeToSQLAndInterpolate(u)
	// Output:
	//Prepared Statement:
	//(SELECT `a1` AS `A`, `a2` AS `B`, 0 AS `_preserve_result_set` FROM `tableA`
	//WHERE (`a1` = ?))
	//UNION ALL
	//(SELECT `b1` AS `A`, `b2` AS `B`, 1 AS `_preserve_result_set` FROM `tableB`
	//WHERE (`b1` = ?))
	//UNION ALL
	//(SELECT concat(c1,?,c2) AS `A`, `c2` AS `B`, 2 AS `_preserve_result_set` FROM
	//`tableC` WHERE (`c2` = ?))
	//ORDER BY `_preserve_result_set`, `A`, `B` DESC
	//Arguments: [3 4 - ArgForC2]
	//
	//Interpolated Statement:
	//(SELECT `a1` AS `A`, `a2` AS `B`, 0 AS `_preserve_result_set` FROM `tableA`
	//WHERE (`a1` = 3))
	//UNION ALL
	//(SELECT `b1` AS `A`, `b2` AS `B`, 1 AS `_preserve_result_set` FROM `tableB`
	//WHERE (`b1` = 4))
	//UNION ALL
	//(SELECT concat(c1,'-',c2) AS `A`, `c2` AS `B`, 2 AS `_preserve_result_set` FROM
	//`tableC` WHERE (`c2` = 'ArgForC2'))
	//ORDER BY `_preserve_result_set`, `A`, `B` DESC
}

// ExampleNewUnion_template interpolates the SQL string with its placeholders
// and puts for each placeholder the correct encoded and escaped value into it.
// Eliminates the need for prepared statements. Avoids an additional round trip
// to the database server by sending the query and its arguments directly. If
// you execute a query multiple times within a short time you should use
// prepared statements.
func ExampleNewUnion_template() {

	u := dbr.NewUnion(
		dbr.NewSelect().AddColumns("t.value", "t.attribute_id", "t.store_id").
			FromAlias("catalog_product_entity_$type$", "t").
			Where(dbr.Column("entity_id").Int64(1561), dbr.Column("store_id").In().Int64s(1, 0)),
	).
		StringReplace("$type$", "varchar", "int", "decimal", "datetime", "text").
		PreserveResultSet().
		All().OrderBy("attribute_id", "store_id")
	writeToSQLAndInterpolate(u)
	// Output:
	//Prepared Statement:
	//(SELECT `t`.`value`, `t`.`attribute_id`, `t`.`store_id`, 0 AS
	//`_preserve_result_set` FROM `catalog_product_entity_varchar` AS `t` WHERE
	//(`entity_id` = ?) AND (`store_id` IN (?,?)))
	//UNION ALL
	//(SELECT `t`.`value`, `t`.`attribute_id`, `t`.`store_id`, 1 AS
	//`_preserve_result_set` FROM `catalog_product_entity_int` AS `t` WHERE
	//(`entity_id` = ?) AND (`store_id` IN (?,?)))
	//UNION ALL
	//(SELECT `t`.`value`, `t`.`attribute_id`, `t`.`store_id`, 2 AS
	//`_preserve_result_set` FROM `catalog_product_entity_decimal` AS `t` WHERE
	//(`entity_id` = ?) AND (`store_id` IN (?,?)))
	//UNION ALL
	//(SELECT `t`.`value`, `t`.`attribute_id`, `t`.`store_id`, 3 AS
	//`_preserve_result_set` FROM `catalog_product_entity_datetime` AS `t` WHERE
	//(`entity_id` = ?) AND (`store_id` IN (?,?)))
	//UNION ALL
	//(SELECT `t`.`value`, `t`.`attribute_id`, `t`.`store_id`, 4 AS
	//`_preserve_result_set` FROM `catalog_product_entity_text` AS `t` WHERE
	//(`entity_id` = ?) AND (`store_id` IN (?,?)))
	//ORDER BY `_preserve_result_set`, `attribute_id`, `store_id`
	//Arguments: [1561 1 0 1561 1 0 1561 1 0 1561 1 0 1561 1 0]
	//
	//Interpolated Statement:
	//(SELECT `t`.`value`, `t`.`attribute_id`, `t`.`store_id`, 0 AS
	//`_preserve_result_set` FROM `catalog_product_entity_varchar` AS `t` WHERE
	//(`entity_id` = 1561) AND (`store_id` IN (1,0)))
	//UNION ALL
	//(SELECT `t`.`value`, `t`.`attribute_id`, `t`.`store_id`, 1 AS
	//`_preserve_result_set` FROM `catalog_product_entity_int` AS `t` WHERE
	//(`entity_id` = 1561) AND (`store_id` IN (1,0)))
	//UNION ALL
	//(SELECT `t`.`value`, `t`.`attribute_id`, `t`.`store_id`, 2 AS
	//`_preserve_result_set` FROM `catalog_product_entity_decimal` AS `t` WHERE
	//(`entity_id` = 1561) AND (`store_id` IN (1,0)))
	//UNION ALL
	//(SELECT `t`.`value`, `t`.`attribute_id`, `t`.`store_id`, 3 AS
	//`_preserve_result_set` FROM `catalog_product_entity_datetime` AS `t` WHERE
	//(`entity_id` = 1561) AND (`store_id` IN (1,0)))
	//UNION ALL
	//(SELECT `t`.`value`, `t`.`attribute_id`, `t`.`store_id`, 4 AS
	//`_preserve_result_set` FROM `catalog_product_entity_text` AS `t` WHERE
	//(`entity_id` = 1561) AND (`store_id` IN (1,0)))
	//ORDER BY `_preserve_result_set`, `attribute_id`, `store_id`
}

func ExampleInterpolate() {
	sqlStr := dbr.Interpolate("SELECT * FROM x WHERE a IN (?) AND b IN (?) AND c NOT IN (?) AND d BETWEEN ? AND ?").
		Ints(1).
		Ints(1, 2, 3).
		Int64s(5, 6, 7).
		Str("wat").
		Str("ok").
		// `MustString` panics on error, or use `String` which prints the error into
		// the returned string and hence generates invalid SQL. Alternatively use
		// `ToSQL`.
		MustString()

	fmt.Printf("%s\n", sqlStr)
	// Output:
	// SELECT * FROM x WHERE a IN (1) AND b IN (1,2,3) AND c NOT IN (5,6,7) AND d BETWEEN 'wat' AND 'ok'
}

func ExampleRepeat() {
	args := dbr.MakeArgs(2).Ints(5, 7, 9).Strs("a", "b", "c", "d", "e")
	sqlStr, err := dbr.Repeat("SELECT * FROM `table` WHERE id IN (?) AND name IN (?)", args)
	if err != nil {
		fmt.Printf("%+v\n", err)
		return
	}

	fmt.Printf("%s\nArguments: %v\n", sqlStr, args.Interfaces())
	// Output:
	// SELECT * FROM `table` WHERE id IN (?,?,?) AND name IN (?,?,?,?,?)
	// Arguments: [5 7 9 a b c d e]
}

func argPrinter(wf *dbr.Condition) {
	sqlStr, args, err := dbr.NewSelect().AddColumns("a", "b").
		From("c").Where(wf).ToSQL()
	if err != nil {
		fmt.Printf("%+v\n", err)
	} else {
		fmt.Printf("%q", sqlStr)
		if len(args) > 0 {
			fmt.Printf(" Arguments: %v", args)
		}
		fmt.Print("\n")
	}
}

// ExampleArguments is duplicate of ExampleColumn
func ExampleArguments() {

	argPrinter(dbr.Column("d").Null())
	argPrinter(dbr.Column("d").NotNull())
	argPrinter(dbr.Column("d").Int(2))
	argPrinter(dbr.Column("d").Int(3).Null())
	argPrinter(dbr.Column("d").Int(4).NotNull())
	argPrinter(dbr.Column("d").In().Ints(7, 8, 9))
	argPrinter(dbr.Column("d").NotIn().Ints(10, 11, 12))
	argPrinter(dbr.Column("d").Between().Ints(13, 14))
	argPrinter(dbr.Column("d").NotBetween().Ints(15, 16))
	argPrinter(dbr.Column("d").Greatest().Ints(17, 18, 19))
	argPrinter(dbr.Column("d").Least().Ints(20, 21, 22))
	argPrinter(dbr.Column("d").Equal().Int(30))
	argPrinter(dbr.Column("d").NotEqual().Int(31))
	argPrinter(dbr.Column("alias.column").SpaceShip().Float64(3.14159))

	argPrinter(dbr.Column("d").Less().Int(32))
	argPrinter(dbr.Column("d").Greater().Int(33))
	argPrinter(dbr.Column("d").LessOrEqual().Int(34))
	argPrinter(dbr.Column("d").GreaterOrEqual().Int(35))

	argPrinter(dbr.Column("d").Like().Str("Goph%"))
	argPrinter(dbr.Column("d").NotLike().Str("Cat%"))

	//Output:
	//"SELECT `a`, `b` FROM `c` WHERE (`d` IS NULL)"
	//"SELECT `a`, `b` FROM `c` WHERE (`d` IS NOT NULL)"
	//"SELECT `a`, `b` FROM `c` WHERE (`d` = ?)" Arguments: [2]
	//"SELECT `a`, `b` FROM `c` WHERE (`d` IS NULL)"
	//"SELECT `a`, `b` FROM `c` WHERE (`d` IS NOT NULL)"
	//"SELECT `a`, `b` FROM `c` WHERE (`d` IN (?,?,?))" Arguments: [7 8 9]
	//"SELECT `a`, `b` FROM `c` WHERE (`d` NOT IN (?,?,?))" Arguments: [10 11 12]
	//"SELECT `a`, `b` FROM `c` WHERE (`d` BETWEEN ? AND ?)" Arguments: [13 14]
	//"SELECT `a`, `b` FROM `c` WHERE (`d` NOT BETWEEN ? AND ?)" Arguments: [15 16]
	//"SELECT `a`, `b` FROM `c` WHERE (`d` GREATEST (?,?,?))" Arguments: [17 18 19]
	//"SELECT `a`, `b` FROM `c` WHERE (`d` LEAST (?,?,?))" Arguments: [20 21 22]
	//"SELECT `a`, `b` FROM `c` WHERE (`d` = ?)" Arguments: [30]
	//"SELECT `a`, `b` FROM `c` WHERE (`d` != ?)" Arguments: [31]
	//"SELECT `a`, `b` FROM `c` WHERE (`alias`.`column` <=> ?)" Arguments: [3.14159]
	//"SELECT `a`, `b` FROM `c` WHERE (`d` < ?)" Arguments: [32]
	//"SELECT `a`, `b` FROM `c` WHERE (`d` > ?)" Arguments: [33]
	//"SELECT `a`, `b` FROM `c` WHERE (`d` <= ?)" Arguments: [34]
	//"SELECT `a`, `b` FROM `c` WHERE (`d` >= ?)" Arguments: [35]
	//"SELECT `a`, `b` FROM `c` WHERE (`d` LIKE ?)" Arguments: [Goph%]
	//"SELECT `a`, `b` FROM `c` WHERE (`d` NOT LIKE ?)" Arguments: [Cat%]
}

// ExampleColumn is a duplicate of ExampleArgument
func ExampleColumn() {

	argPrinter(dbr.Column("d").Null())
	argPrinter(dbr.Column("d").NotNull())
	argPrinter(dbr.Column("d").Int(2))
	argPrinter(dbr.Column("d").Int(3).Null())
	argPrinter(dbr.Column("d").Int(4).NotNull())
	argPrinter(dbr.Column("d").In().Ints(7, 8, 9))
	argPrinter(dbr.Column("d").NotIn().Ints(10, 11, 12))
	argPrinter(dbr.Column("d").Between().Ints(13, 14))
	argPrinter(dbr.Column("d").NotBetween().Ints(15, 16))
	argPrinter(dbr.Column("d").Greatest().Ints(17, 18, 19))
	argPrinter(dbr.Column("d").Least().Ints(20, 21, 22))
	argPrinter(dbr.Column("d").Equal().Int(30))
	argPrinter(dbr.Column("d").NotEqual().Int(31))
	argPrinter(dbr.Column("alias.column").SpaceShip().Float64(3.14159))

	argPrinter(dbr.Column("d").Less().Int(32))
	argPrinter(dbr.Column("d").Greater().Int(33))
	argPrinter(dbr.Column("d").LessOrEqual().Int(34))
	argPrinter(dbr.Column("d").GreaterOrEqual().Int(35))

	argPrinter(dbr.Column("d").Like().Str("Goph%"))
	argPrinter(dbr.Column("d").NotLike().Str("Cat%"))

	//Output:
	//"SELECT `a`, `b` FROM `c` WHERE (`d` IS NULL)"
	//"SELECT `a`, `b` FROM `c` WHERE (`d` IS NOT NULL)"
	//"SELECT `a`, `b` FROM `c` WHERE (`d` = ?)" Arguments: [2]
	//"SELECT `a`, `b` FROM `c` WHERE (`d` IS NULL)"
	//"SELECT `a`, `b` FROM `c` WHERE (`d` IS NOT NULL)"
	//"SELECT `a`, `b` FROM `c` WHERE (`d` IN (?,?,?))" Arguments: [7 8 9]
	//"SELECT `a`, `b` FROM `c` WHERE (`d` NOT IN (?,?,?))" Arguments: [10 11 12]
	//"SELECT `a`, `b` FROM `c` WHERE (`d` BETWEEN ? AND ?)" Arguments: [13 14]
	//"SELECT `a`, `b` FROM `c` WHERE (`d` NOT BETWEEN ? AND ?)" Arguments: [15 16]
	//"SELECT `a`, `b` FROM `c` WHERE (`d` GREATEST (?,?,?))" Arguments: [17 18 19]
	//"SELECT `a`, `b` FROM `c` WHERE (`d` LEAST (?,?,?))" Arguments: [20 21 22]
	//"SELECT `a`, `b` FROM `c` WHERE (`d` = ?)" Arguments: [30]
	//"SELECT `a`, `b` FROM `c` WHERE (`d` != ?)" Arguments: [31]
	//"SELECT `a`, `b` FROM `c` WHERE (`alias`.`column` <=> ?)" Arguments: [3.14159]
	//"SELECT `a`, `b` FROM `c` WHERE (`d` < ?)" Arguments: [32]
	//"SELECT `a`, `b` FROM `c` WHERE (`d` > ?)" Arguments: [33]
	//"SELECT `a`, `b` FROM `c` WHERE (`d` <= ?)" Arguments: [34]
	//"SELECT `a`, `b` FROM `c` WHERE (`d` >= ?)" Arguments: [35]
	//"SELECT `a`, `b` FROM `c` WHERE (`d` LIKE ?)" Arguments: [Goph%]
	//"SELECT `a`, `b` FROM `c` WHERE (`d` NOT LIKE ?)" Arguments: [Cat%]
}

func ExampleCondition_Sub() {
	s := dbr.NewSelect("sku", "type_id").
		From("catalog_product_entity").
		Where(dbr.Column("entity_id").In().Sub(
			dbr.NewSelect().From("catalog_category_product").
				AddColumns("entity_id").Where(dbr.Column("category_id").Int64(234)),
		))
	writeToSQLAndInterpolate(s)
	// Output:
	//Prepared Statement:
	//SELECT `sku`, `type_id` FROM `catalog_product_entity` WHERE (`entity_id` IN
	//(SELECT `entity_id` FROM `catalog_category_product` WHERE (`category_id` = ?)))
	//Arguments: [234]
	//
	//Interpolated Statement:
	//SELECT `sku`, `type_id` FROM `catalog_product_entity` WHERE (`entity_id` IN
	//(SELECT `entity_id` FROM `catalog_category_product` WHERE (`category_id` =
	//234)))
}

func ExampleNewSelectWithDerivedTable() {
	sel3 := dbr.NewSelect().Unsafe().FromAlias("sales_bestsellers_aggregated_daily", "t3").
		AddColumnsAliases("DATE_FORMAT(t3.period, '%Y-%m-01')", "period").
		AddColumns("t3.store_id", "t3.product_id", "t3.product_name").
		AddColumnsAliases("AVG(`t3`.`product_price`)", "avg_price", "SUM(t3.qty_ordered)", "total_qty").
		Where(dbr.Column("product_name").Str("Canon%")).
		GroupBy("t3.store_id").
		GroupBy("DATE_FORMAT(t3.period, '%Y-%m-01')").
		GroupBy("t3.product_id", "t3.product_name").
		OrderBy("t3.store_id").
		OrderBy("DATE_FORMAT(t3.period, '%Y-%m-01')").
		OrderByDesc("total_qty")

	sel1 := dbr.NewSelectWithDerivedTable(sel3, "t1").
		AddColumns("t1.period", "t1.store_id", "t1.product_id", "t1.product_name", "t1.avg_price", "t1.qty_ordered").
		Where(dbr.Column("product_name").Str("Sony%")).
		OrderBy("t1.period", "t1.product_id")
	writeToSQLAndInterpolate(sel1)
	// Output:
	//Prepared Statement:
	//SELECT `t1`.`period`, `t1`.`store_id`, `t1`.`product_id`, `t1`.`product_name`,
	//`t1`.`avg_price`, `t1`.`qty_ordered` FROM (SELECT DATE_FORMAT(t3.period,
	//'%Y-%m-01') AS `period`, `t3`.`store_id`, `t3`.`product_id`,
	//`t3`.`product_name`, AVG(`t3`.`product_price`) AS `avg_price`,
	//SUM(t3.qty_ordered) AS `total_qty` FROM `sales_bestsellers_aggregated_daily` AS
	//`t3` WHERE (`product_name` = ?) GROUP BY `t3`.`store_id`, DATE_FORMAT(t3.period,
	//'%Y-%m-01'), `t3`.`product_id`, `t3`.`product_name` ORDER BY `t3`.`store_id`,
	//DATE_FORMAT(t3.period, '%Y-%m-01'), `total_qty` DESC) AS `t1` WHERE
	//(`product_name` = ?) ORDER BY `t1`.`period`, `t1`.`product_id`
	//Arguments: [Canon% Sony%]
	//
	//Interpolated Statement:
	//SELECT `t1`.`period`, `t1`.`store_id`, `t1`.`product_id`, `t1`.`product_name`,
	//`t1`.`avg_price`, `t1`.`qty_ordered` FROM (SELECT DATE_FORMAT(t3.period,
	//'%Y-%m-01') AS `period`, `t3`.`store_id`, `t3`.`product_id`,
	//`t3`.`product_name`, AVG(`t3`.`product_price`) AS `avg_price`,
	//SUM(t3.qty_ordered) AS `total_qty` FROM `sales_bestsellers_aggregated_daily` AS
	//`t3` WHERE (`product_name` = 'Canon%') GROUP BY `t3`.`store_id`,
	//DATE_FORMAT(t3.period, '%Y-%m-01'), `t3`.`product_id`, `t3`.`product_name` ORDER
	//BY `t3`.`store_id`, DATE_FORMAT(t3.period, '%Y-%m-01'), `total_qty` DESC) AS
	//`t1` WHERE (`product_name` = 'Sony%') ORDER BY `t1`.`period`, `t1`.`product_id`
}

func ExampleSQLIfNull() {
	s := dbr.NewSelect().AddColumnsConditions(
		dbr.SQLIfNull("column1"),
		dbr.SQLIfNull("table1.column1"),
		dbr.SQLIfNull("column1", "column2"),
		dbr.SQLIfNull("table1.column1", "table2.column2"),
		dbr.SQLIfNull("column2", "1/0").Alias("alias"),
		dbr.SQLIfNull("SELECT * FROM x", "8").Alias("alias"),
		dbr.SQLIfNull("SELECT * FROM x", "9 ").Alias("alias"),
		dbr.SQLIfNull("column1", "column2").Alias("alias"),
		dbr.SQLIfNull("table1.column1", "table2.column2").Alias("alias"),
		dbr.SQLIfNull("table1", "column1", "table2", "column2"),
		dbr.SQLIfNull("table1", "column1", "table2", "column2").Alias("alias"),
	).From("table1")
	sStr, _, _ := s.ToSQL()
	fmt.Print(strings.Replace(sStr, ", ", ",\n", -1))

	//Output:
	//SELECT IFNULL(`column1`,NULL),
	//IFNULL(`table1`.`column1`,NULL),
	//IFNULL(`column1`,`column2`),
	//IFNULL(`table1`.`column1`,`table2`.`column2`),
	//IFNULL(`column2`,1/0) AS `alias`,
	//IFNULL(SELECT * FROM x,8) AS `alias`,
	//IFNULL(SELECT * FROM x,9 ) AS `alias`,
	//IFNULL(`column1`,`column2`) AS `alias`,
	//IFNULL(`table1`.`column1`,`table2`.`column2`) AS `alias`,
	//IFNULL(`table1`.`column1`,`table2`.`column2`),
	//IFNULL(`table1`.`column1`,`table2`.`column2`) AS `alias` FROM `table1`
}

func ExampleSQLIf() {
	s := dbr.NewSelect().
		AddColumns("a", "b", "c").
		From("table1").
		Where(
			dbr.SQLIf("a > 0", "b", "c").Greater().Int(4711),
		)
	writeToSQLAndInterpolate(s)

	// Output:
	//Prepared Statement:
	//SELECT `a`, `b`, `c` FROM `table1` WHERE (IF((a > 0), b, c) > ?)
	//Arguments: [4711]
	//
	//Interpolated Statement:
	//SELECT `a`, `b`, `c` FROM `table1` WHERE (IF((a > 0), b, c) > 4711)
}

func ExampleSQLCase_update() {
	u := dbr.NewUpdate("cataloginventory_stock_item").
		Set(dbr.Column("qty").SQLCase("`product_id`", "qty",
			"3456", "qty+?",
			"3457", "qty+?",
			"3458", "qty+?",
		).Ints(3, 4, 5)).
		Where(
			dbr.Column("product_id").In().Int64s(345, 567, 897),
			dbr.Column("website_id").Int64(6),
		)
	writeToSQLAndInterpolate(u)

	// Output:
	//Prepared Statement:
	//UPDATE `cataloginventory_stock_item` SET `qty`=CASE `product_id` WHEN 3456 THEN
	//qty+? WHEN 3457 THEN qty+? WHEN 3458 THEN qty+? ELSE qty END WHERE (`product_id`
	//IN (?,?,?)) AND (`website_id` = ?)
	//Arguments: [3 4 5 345 567 897 6]
	//
	//Interpolated Statement:
	//UPDATE `cataloginventory_stock_item` SET `qty`=CASE `product_id` WHEN 3456 THEN
	//qty+3 WHEN 3457 THEN qty+4 WHEN 3458 THEN qty+5 ELSE qty END WHERE (`product_id`
	//IN (345,567,897)) AND (`website_id` = 6)
}

// ExampleSQLCase_select is a duplicate of ExampleSelect_AddArguments
func ExampleSQLCase_select() {
	// time stamp has no special meaning ;-)
	start := time.Unix(1257894000, 0)
	end := time.Unix(1257980400, 0)

	s := dbr.NewSelect().AddColumns("price", "sku", "name", "title", "description").
		AddColumnsConditions(
			dbr.SQLCase("", "`closed`",
				"date_start <= ? AND date_end >= ?", "`open`",
				"date_start > ? AND date_end > ?", "`upcoming`",
			).Alias("is_on_sale").Times(start, end, start, end),
		).
		From("catalog_promotions").Where(
		dbr.Column("promotion_id").NotIn().Ints(4711, 815, 42))
	writeToSQLAndInterpolate(s)

	// Output:
	//Prepared Statement:
	//SELECT `price`, `sku`, `name`, `title`, `description`, CASE  WHEN date_start <=
	//? AND date_end >= ? THEN `open` WHEN date_start > ? AND date_end > ? THEN
	//`upcoming` ELSE `closed` END AS `is_on_sale` FROM `catalog_promotions` WHERE
	//(`promotion_id` NOT IN (?,?,?))
	//Arguments: [2009-11-11 00:00:00 +0100 CET 2009-11-12 00:00:00 +0100 CET 2009-11-11 00:00:00 +0100 CET 2009-11-12 00:00:00 +0100 CET 4711 815 42]
	//
	//Interpolated Statement:
	//SELECT `price`, `sku`, `name`, `title`, `description`, CASE  WHEN date_start <=
	//'2009-11-11 00:00:00' AND date_end >= '2009-11-12 00:00:00' THEN `open` WHEN
	//date_start > '2009-11-11 00:00:00' AND date_end > '2009-11-12 00:00:00' THEN
	//`upcoming` ELSE `closed` END AS `is_on_sale` FROM `catalog_promotions` WHERE
	//(`promotion_id` NOT IN (4711,815,42))
}

// ExampleSelect_AddColumnsConditions is duplicate of ExampleSQLCase_select
func ExampleSelect_AddColumnsConditions() {

	start := time.Unix(1257894000, 0)
	end := time.Unix(1257980400, 0)

	s := dbr.NewSelect().AddColumns("price", "sku", "name", "title", "description").
		AddColumnsConditions(
			dbr.SQLCase("", "`closed`",
				"date_start <= ? AND date_end >= ?", "`open`",
				"date_start > ? AND date_end > ?", "`upcoming`",
			).Alias("is_on_sale").Times(start, end, start, end),
		).
		From("catalog_promotions").Where(
		dbr.Column("promotion_id").NotIn().Ints(4711, 815, 42))
	writeToSQLAndInterpolate(s)

	// Output:
	//Prepared Statement:
	//SELECT `price`, `sku`, `name`, `title`, `description`, CASE  WHEN date_start <=
	//? AND date_end >= ? THEN `open` WHEN date_start > ? AND date_end > ? THEN
	//`upcoming` ELSE `closed` END AS `is_on_sale` FROM `catalog_promotions` WHERE
	//(`promotion_id` NOT IN (?,?,?))
	//Arguments: [2009-11-11 00:00:00 +0100 CET 2009-11-12 00:00:00 +0100 CET 2009-11-11 00:00:00 +0100 CET 2009-11-12 00:00:00 +0100 CET 4711 815 42]
	//
	//Interpolated Statement:
	//SELECT `price`, `sku`, `name`, `title`, `description`, CASE  WHEN date_start <=
	//'2009-11-11 00:00:00' AND date_end >= '2009-11-12 00:00:00' THEN `open` WHEN
	//date_start > '2009-11-11 00:00:00' AND date_end > '2009-11-12 00:00:00' THEN
	//`upcoming` ELSE `closed` END AS `is_on_sale` FROM `catalog_promotions` WHERE
	//(`promotion_id` NOT IN (4711,815,42))
}

func ExampleParenthesisOpen() {
	s := dbr.NewSelect("columnA", "columnB").
		Distinct().
		FromAlias("tableC", "ccc").
		Where(
			dbr.ParenthesisOpen(),
			dbr.Column("d").Int(1),
			dbr.Column("e").Str("wat").Or(),
			dbr.ParenthesisClose(),
			dbr.Column("f").Int(2),
		).
		GroupBy("ab").
		Having(
			dbr.Expr("j = k"),
			dbr.ParenthesisOpen(),
			dbr.Column("m").Int(33),
			dbr.Column("n").Str("wh3r3").Or(),
			dbr.ParenthesisClose(),
		).
		OrderBy("l").
		Limit(7).
		Offset(8)
	writeToSQLAndInterpolate(s)

	// Output:
	//Prepared Statement:
	//SELECT DISTINCT `columnA`, `columnB` FROM `tableC` AS `ccc` WHERE ((`d` = ?) OR
	//(`e` = ?)) AND (`f` = ?) GROUP BY `ab` HAVING (j = k) AND ((`m` = ?) OR (`n` =
	//?)) ORDER BY `l` LIMIT 7 OFFSET 8
	//Arguments: [1 wat 2 33 wh3r3]
	//
	//Interpolated Statement:
	//SELECT DISTINCT `columnA`, `columnB` FROM `tableC` AS `ccc` WHERE ((`d` = 1) OR
	//(`e` = 'wat')) AND (`f` = 2) GROUP BY `ab` HAVING (j = k) AND ((`m` = 33) OR
	//(`n` = 'wh3r3')) ORDER BY `l` LIMIT 7 OFFSET 8
}

func ExampleWith_Union() {
	// Non-recursive CTE
	// Sales: Find best and worst month:
	cte := dbr.NewWith(
		dbr.WithCTE{Name: "sales_by_month", Columns: []string{"month", "total"},
			Select: dbr.NewSelect().Unsafe().AddColumns("Month(day_of_sale)", "Sum(amount)").From("sales_days").
				Where(dbr.Expr("Year(day_of_sale) = ?").Int(2015)).
				GroupBy("Month(day_of_sale))"),
		},
		dbr.WithCTE{Name: "best_month", Columns: []string{"month", "total", "award"},
			Select: dbr.NewSelect().Unsafe().AddColumns("month", "total").AddColumns(`"best"`).From("sales_by_month").
				Where(dbr.Column("total").Equal().Sub(dbr.NewSelect().Unsafe().AddColumns("Max(total)").From("sales_by_month"))),
		},
		dbr.WithCTE{Name: "worst_month", Columns: []string{"month", "total", "award"},
			Select: dbr.NewSelect().Unsafe().AddColumns("month", "total").AddColumns(`"worst"`).From("sales_by_month").
				Where(dbr.Column("total").Equal().Sub(dbr.NewSelect().Unsafe().AddColumns("Min(total)").From("sales_by_month"))),
		},
	).Union(dbr.NewUnion(
		dbr.NewSelect().Star().From("best_month"),
		dbr.NewSelect().Star().From("worst_month"),
	).All())
	writeToSQLAndInterpolate(cte)

	//Result:
	//+-------+-------+-------+
	//| month | total | award |
	//+-------+-------+-------+
	//|     1 |   300 | best  |
	//|     3 |    11 | worst |
	//+-------+-------+-------+

	// Output:
	//Prepared Statement:
	//WITH
	//`sales_by_month` (`month`,`total`) AS (SELECT Month(day_of_sale), Sum(amount)
	//FROM `sales_days` WHERE (Year(day_of_sale) = ?) GROUP BY Month(day_of_sale))),
	//`best_month` (`month`,`total`,`award`) AS (SELECT `month`, `total`, "best" FROM
	//`sales_by_month` WHERE (`total` = (SELECT Max(total) FROM `sales_by_month`))),
	//`worst_month` (`month`,`total`,`award`) AS (SELECT `month`, `total`, "worst"
	//FROM `sales_by_month` WHERE (`total` = (SELECT Min(total) FROM
	//`sales_by_month`)))
	//(SELECT * FROM `best_month`)
	//UNION ALL
	//(SELECT * FROM `worst_month`)
	//Arguments: [2015]
	//
	//Interpolated Statement:
	//WITH
	//`sales_by_month` (`month`,`total`) AS (SELECT Month(day_of_sale), Sum(amount)
	//FROM `sales_days` WHERE (Year(day_of_sale) = 2015) GROUP BY
	//Month(day_of_sale))),
	//`best_month` (`month`,`total`,`award`) AS (SELECT `month`, `total`, 'best' FROM
	//`sales_by_month` WHERE (`total` = (SELECT Max(total) FROM `sales_by_month`))),
	//`worst_month` (`month`,`total`,`award`) AS (SELECT `month`, `total`, 'worst'
	//FROM `sales_by_month` WHERE (`total` = (SELECT Min(total) FROM
	//`sales_by_month`)))
	//(SELECT * FROM `best_month`)
	//UNION ALL
	//(SELECT * FROM `worst_month`)
}
