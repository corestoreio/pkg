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
	"time"

	"github.com/corestoreio/csfw/storage/dbr"
	"github.com/corestoreio/csfw/util/wordwrap"
	"github.com/corestoreio/errors"
)

// iFaceToArgs unpacks the interface and creates an Value slice. Just a
// helper function for the examples.
func iFaceToArgs(values ...interface{}) dbr.Values {
	args := make(dbr.Values, 0, len(values))
	for _, val := range values {
		switch v := val.(type) {
		case float64:
			args = append(args, dbr.Float64(v))
		case int64:
			args = append(args, dbr.Int64(v))
		case int:
			args = append(args, dbr.Int64(v))
			args = append(args, dbr.Int64(v))
		case bool:
			args = append(args, dbr.Bool(v))
		case string:
			args = append(args, dbr.String(v))
		case []byte:
			args = append(args, dbr.Bytes(v))
		case time.Time:
			args = append(args, dbr.MakeTime(v))
		case *time.Time:
			if v != nil {
				args = append(args, dbr.MakeTime(*v))
			}
		case nil:
			args = append(args, dbr.NullValue())
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
	if len(args) > 0 {
		fmt.Printf("\nValues: %v\n\n", args)
	} else {
		fmt.Print("\n")
	}
	if len(args) == 0 {
		return
	}
	sqlStr, err = dbr.Interpolate(sqlStr, iFaceToArgs(args...)...)
	if err != nil {
		fmt.Printf("%+v\n", err)
		return
	}
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
	//Values: [1 2 Three <nil> 5 6 Seven 3.14156]
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
	//Values: [2046 33 3 2046 34 3 2046 35 3]
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
	//Values: [2046 33 3 2046 34 3 2046 35 3]
	//
	//Interpolated Statement:
	//INSERT INTO `catalog_product_link` VALUES (2046,33,3),(2046,34,3),(2046,35,3)
	//
	//Prepared Statement:
	//INSERT INTO `catalog_product_link`
	//(`product_id`,`linked_product_id`,`link_type_id`) VALUES (?,?,?),(?,?,?),(?,?,?)
	//Values: [2046 33 3 2046 34 3 2046 35 3]
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
		AddOnDuplicateKey("name", dbr.String("Pik3")).
		AddOnDuplicateKey("email", nil)
	writeToSQLAndInterpolate(i)

	// Output:
	//Prepared Statement:
	//INSERT INTO `dbr_people` (`id`,`name`,`email`) VALUES (?,?,?) ON DUPLICATE KEY
	//UPDATE `name`=?, `email`=VALUES(`email`)
	//Values: [1 Pik'e pikes@peak.com Pik3]
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
				dbr.Column("string").String("wat").Or(),
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
	//Values: [1 wat 1 2 3]
	//
	//Interpolated Statement:
	//INSERT INTO `tableA` SELECT `something_id`, `user_id`, `other` FROM `some_table`
	//WHERE ((`int64A` >= 1) OR (`string` = 'wat')) AND (`int64B` IN (1,2,3)) ORDER BY
	//`id` DESC LIMIT 20 OFFSET 0
}

func ExampleInsert_Pair() {
	ins := dbr.NewInsert("catalog_product_link").
		Pair("product_id", dbr.Int64(2046)).
		Pair("linked_product_id", dbr.Int64(33)).
		Pair("link_type_id", dbr.Int64(3)).
		// next row
		Pair("product_id", dbr.Int64(2046)).
		Pair("linked_product_id", dbr.Int64(34)).
		Pair("link_type_id", dbr.Int64(3))
	// next row ...
	writeToSQLAndInterpolate(ins)
	// Output:
	//Prepared Statement:
	//INSERT INTO `catalog_product_link`
	//(`product_id`,`linked_product_id`,`link_type_id`) VALUES (?,?,?),(?,?,?)
	//Values: [2046 33 3 2046 34 3]
	//
	//Interpolated Statement:
	//INSERT INTO `catalog_product_link`
	//(`product_id`,`linked_product_id`,`link_type_id`) VALUES (2046,33,3),(2046,34,3)
}

func ExampleNewDelete() {
	d := dbr.NewDelete("tableA").Where(
		dbr.Column("a").Like().String("b'%"),
		dbr.Column("b").In().Ints(3, 4, 5, 6),
	).
		Limit(1).OrderBy("id")
	writeToSQLAndInterpolate(d)
	// Output:
	//Prepared Statement:
	//DELETE FROM `tableA` WHERE (`a` LIKE ?) AND (`b` IN (?,?,?,?)) ORDER BY `id`
	//LIMIT 1
	//Values: [b'% 3 4 5 6]
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
		dbr.NewSelect().AddColumnsAlias("a1", "A", "a2", "B").From("tableA").Where(dbr.Column("a1").Int64(3)),
		dbr.NewSelect().AddColumnsAlias("b1", "A", "b2", "B").From("tableB").Where(dbr.Column("b1").Int64(4)),
	)
	// Maybe more of your code ...
	u.Append(
		dbr.NewSelect().AddColumnsExprAlias("concat(c1,?,c2)", "A").
			AddArguments(dbr.String("-")).
			AddColumnsAlias("c2", "B").
			From("tableC").Where(dbr.Column("c2").String("ArgForC2")),
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
	//ORDER BY `_preserve_result_set`, `A` ASC, `B` DESC
	//Values: [3 4 - ArgForC2]
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
	//ORDER BY `_preserve_result_set`, `A` ASC, `B` DESC
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
			FromAlias("catalog_product_entity_{type}", "t").
			Where(dbr.Column("entity_id").Int64(1561), dbr.Column("store_id").In().Int64s(1, 0)),
	).
		StringReplace("{type}", "varchar", "int", "decimal", "datetime", "text").
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
	//ORDER BY `_preserve_result_set`, `attribute_id` ASC, `store_id` ASC
	//Values: [1561 1 0 1561 1 0 1561 1 0 1561 1 0 1561 1 0]
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
	//ORDER BY `_preserve_result_set`, `attribute_id` ASC, `store_id` ASC
}

func ExampleInterpolate() {
	sqlStr, err := dbr.Interpolate("SELECT * FROM x WHERE a IN (?) AND b IN (?) AND c NOT IN (?) AND d BETWEEN ? AND ?",
		dbr.Ints{1},
		dbr.Ints{1, 2, 3},
		dbr.Int64s{5, 6, 7},
		dbr.String("wat"),
		dbr.String("ok"),
	)
	if err != nil {
		fmt.Printf("%+v\n", err)
		return
	}
	fmt.Printf("%s\n", sqlStr)
	// Output:
	// SELECT * FROM x WHERE a IN (1) AND b IN (1,2,3) AND c NOT IN (5,6,7) AND d BETWEEN 'wat' AND 'ok'
}

func ExampleRepeat() {
	sl := dbr.Strings{"a", "b", "c", "d", "e"}

	sqlStr, args, err := dbr.Repeat("SELECT * FROM `table` WHERE id IN (?) AND name IN (?)",
		dbr.Ints{5, 7, 9}, sl)

	if err != nil {
		fmt.Printf("%+v\n", err)
		return
	}

	fmt.Printf("%s\nValues: %v\n", sqlStr, args)
	// Output:
	// SELECT * FROM `table` WHERE id IN (?,?,?) AND name IN (?,?,?,?,?)
	// Values: [5 7 9 a b c d e]
}

func argPrinter(wf *dbr.WhereFragment) {
	sqlStr, args, err := dbr.NewSelect().AddColumns("a", "b").
		From("c").Where(wf).ToSQL()
	if err != nil {
		fmt.Printf("%+v\n", err)
	} else {
		fmt.Printf("%q", sqlStr)
		if len(args) > 0 {
			fmt.Printf(" Values: %v", args)
		}
		fmt.Print("\n")
	}
}

// ExampleArgument is duplicate of ExampleColumn
func ExampleArgument() {

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

	argPrinter(dbr.Column("d").Less().Int(32))
	argPrinter(dbr.Column("d").Greater().Int(33))
	argPrinter(dbr.Column("d").LessOrEqual().Int(34))
	argPrinter(dbr.Column("d").GreaterOrEqual().Int(35))

	argPrinter(dbr.Column("d").Like().String("Goph%"))
	argPrinter(dbr.Column("d").NotLike().String("Cat%"))

	//Output:
	//"SELECT `a`, `b` FROM `c` WHERE (`d` IS NULL)"
	//"SELECT `a`, `b` FROM `c` WHERE (`d` IS NOT NULL)"
	//"SELECT `a`, `b` FROM `c` WHERE (`d` = ?)" Values: [2]
	//"SELECT `a`, `b` FROM `c` WHERE (`d` IS NULL)"
	//"SELECT `a`, `b` FROM `c` WHERE (`d` IS NOT NULL)"
	//"SELECT `a`, `b` FROM `c` WHERE (`d` IN (?,?,?))" Values: [7 8 9]
	//"SELECT `a`, `b` FROM `c` WHERE (`d` NOT IN (?,?,?))" Values: [10 11 12]
	//"SELECT `a`, `b` FROM `c` WHERE (`d` BETWEEN ? AND ?)" Values: [13 14]
	//"SELECT `a`, `b` FROM `c` WHERE (`d` NOT BETWEEN ? AND ?)" Values: [15 16]
	//"SELECT `a`, `b` FROM `c` WHERE (`d` GREATEST (?,?,?))" Values: [17 18 19]
	//"SELECT `a`, `b` FROM `c` WHERE (`d` LEAST (?,?,?))" Values: [20 21 22]
	//"SELECT `a`, `b` FROM `c` WHERE (`d` = ?)" Values: [30]
	//"SELECT `a`, `b` FROM `c` WHERE (`d` != ?)" Values: [31]
	//"SELECT `a`, `b` FROM `c` WHERE (`d` < ?)" Values: [32]
	//"SELECT `a`, `b` FROM `c` WHERE (`d` > ?)" Values: [33]
	//"SELECT `a`, `b` FROM `c` WHERE (`d` <= ?)" Values: [34]
	//"SELECT `a`, `b` FROM `c` WHERE (`d` >= ?)" Values: [35]
	//"SELECT `a`, `b` FROM `c` WHERE (`d` LIKE ?)" Values: [Goph%]
	//"SELECT `a`, `b` FROM `c` WHERE (`d` NOT LIKE ?)" Values: [Cat%]
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

	argPrinter(dbr.Column("d").Less().Int(32))
	argPrinter(dbr.Column("d").Greater().Int(33))
	argPrinter(dbr.Column("d").LessOrEqual().Int(34))
	argPrinter(dbr.Column("d").GreaterOrEqual().Int(35))

	argPrinter(dbr.Column("d").Like().String("Goph%"))
	argPrinter(dbr.Column("d").NotLike().String("Cat%"))

	//Output:
	//"SELECT `a`, `b` FROM `c` WHERE (`d` IS NULL)"
	//"SELECT `a`, `b` FROM `c` WHERE (`d` IS NOT NULL)"
	//"SELECT `a`, `b` FROM `c` WHERE (`d` = ?)" Values: [2]
	//"SELECT `a`, `b` FROM `c` WHERE (`d` IS NULL)"
	//"SELECT `a`, `b` FROM `c` WHERE (`d` IS NOT NULL)"
	//"SELECT `a`, `b` FROM `c` WHERE (`d` IN (?,?,?))" Values: [7 8 9]
	//"SELECT `a`, `b` FROM `c` WHERE (`d` NOT IN (?,?,?))" Values: [10 11 12]
	//"SELECT `a`, `b` FROM `c` WHERE (`d` BETWEEN ? AND ?)" Values: [13 14]
	//"SELECT `a`, `b` FROM `c` WHERE (`d` NOT BETWEEN ? AND ?)" Values: [15 16]
	//"SELECT `a`, `b` FROM `c` WHERE (`d` GREATEST (?,?,?))" Values: [17 18 19]
	//"SELECT `a`, `b` FROM `c` WHERE (`d` LEAST (?,?,?))" Values: [20 21 22]
	//"SELECT `a`, `b` FROM `c` WHERE (`d` = ?)" Values: [30]
	//"SELECT `a`, `b` FROM `c` WHERE (`d` != ?)" Values: [31]
	//"SELECT `a`, `b` FROM `c` WHERE (`d` < ?)" Values: [32]
	//"SELECT `a`, `b` FROM `c` WHERE (`d` > ?)" Values: [33]
	//"SELECT `a`, `b` FROM `c` WHERE (`d` <= ?)" Values: [34]
	//"SELECT `a`, `b` FROM `c` WHERE (`d` >= ?)" Values: [35]
	//"SELECT `a`, `b` FROM `c` WHERE (`d` LIKE ?)" Values: [Goph%]
	//"SELECT `a`, `b` FROM `c` WHERE (`d` NOT LIKE ?)" Values: [Cat%]
}

func ExampleSubSelect() {
	s := dbr.NewSelect("sku", "type_id").
		From("catalog_product_entity").
		Where(dbr.SubSelect(
			"entity_id", dbr.In,
			dbr.NewSelect().From("catalog_category_product").
				AddColumns("entity_id").Where(dbr.Column("category_id").Int64(234)),
		))
	writeToSQLAndInterpolate(s)
	// Output:
	//Prepared Statement:
	//SELECT `sku`, `type_id` FROM `catalog_product_entity` WHERE (`entity_id` IN
	//(SELECT `entity_id` FROM `catalog_category_product` WHERE (`category_id` = ?)))
	//Values: [234]
	//
	//Interpolated Statement:
	//SELECT `sku`, `type_id` FROM `catalog_product_entity` WHERE (`entity_id` IN
	//(SELECT `entity_id` FROM `catalog_category_product` WHERE (`category_id` =
	//234)))
}

func ExampleNewSelectWithDerivedTable() {
	sel3 := dbr.NewSelect().FromAlias("sales_bestsellers_aggregated_daily", "t3").
		AddColumnsExprAlias("DATE_FORMAT(t3.period, '%Y-%m-01')", "period").
		AddColumns("t3.store_id", "t3.product_id", "t3.product_name").
		AddColumnsExprAlias("AVG(`t3`.`product_price`)", "avg_price", "SUM(t3.qty_ordered)", "total_qty").
		Where(dbr.Column("product_name").String("Canon%")).
		GroupBy("t3.store_id").
		GroupByExpr("DATE_FORMAT(t3.period, '%Y-%m-01')").
		GroupBy("t3.product_id", "t3.product_name").
		OrderBy("t3.store_id").
		OrderByExpr("DATE_FORMAT(t3.period, '%Y-%m-01')").
		OrderByDesc("total_qty")

	sel1 := dbr.NewSelectWithDerivedTable(sel3, "t1").
		AddColumns("t1.period", "t1.store_id", "t1.product_id", "t1.product_name", "t1.avg_price", "t1.qty_ordered").
		Where(dbr.Column("product_name").String("Sony%")).
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
	//Values: [Canon% Sony%]
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
	fmt.Println(dbr.SQLIfNull("column1"))
	fmt.Println(dbr.SQLIfNull("table1.column1"))
	fmt.Println(dbr.SQLIfNull("column1", "column2"))
	fmt.Println(dbr.SQLIfNull("table1.column1", "table2.column2"))
	fmt.Println(dbr.SQLIfNull("column2", "1/0", "alias"))
	fmt.Println(dbr.SQLIfNull("SELECT * FROM x", "8", "alias"))
	fmt.Println(dbr.SQLIfNull("SELECT * FROM x", "9 ", "alias"))
	fmt.Println(dbr.SQLIfNull("column1", "column2", "alias"))
	fmt.Println(dbr.SQLIfNull("table1.column1", "table2.column2", "alias"))
	fmt.Println(dbr.SQLIfNull("table1", "column1", "table2", "column2"))
	fmt.Println(dbr.SQLIfNull("table1", "column1", "table2", "column2", "alias"))
	fmt.Println(dbr.SQLIfNull("table1", "column1", "table2", "column2", "alias", "x"))
	fmt.Println(dbr.SQLIfNull("table1", "column1", "table2", "column2", "alias", "x", "y"))
	//Output:
	//IFNULL(`column1`,(NULL ))
	//IFNULL(`table1`.`column1`,(NULL ))
	//IFNULL(`column1`,`column2`)
	//IFNULL(`table1`.`column1`,`table2`.`column2`)
	//IFNULL(`column2`,(1/0)) AS `alias`
	//IFNULL((SELECT * FROM x),`8`) AS `alias`
	//IFNULL((SELECT * FROM x),(9 )) AS `alias`
	//IFNULL(`column1`,`column2`) AS `alias`
	//IFNULL(`table1`.`column1`,`table2`.`column2`) AS `alias`
	//IFNULL(`table1`.`column1`,`table2`.`column2`)
	//IFNULL(`table1`.`column1`,`table2`.`column2`) AS `alias`
	//IFNULL(`table1`.`column1`,`table2`.`column2`) AS `alias_x`
	//IFNULL(`table1`.`column1`,`table2`.`column2`) AS `alias_x_y`
}

func ExampleSQLIf() {
	s := dbr.NewSelect().AddColumns("a", "b", "c").
		From("table1").Where(
		dbr.Expression(dbr.SQLIf("a > 0", "b", "c")).Greater().Int(4711))
	writeToSQLAndInterpolate(s)

	// Output:
	//Prepared Statement:
	//SELECT `a`, `b`, `c` FROM `table1` WHERE (IF((a > 0), b, c) > ?)
	//Values: [4711]
	//
	//Interpolated Statement:
	//SELECT `a`, `b`, `c` FROM `table1` WHERE (IF((a > 0), b, c) > 4711)
}

func ExampleSQLCase_update() {
	u := dbr.NewUpdate("cataloginventory_stock_item").
		Set("qty", dbr.ExpressionValue(dbr.SQLCase("`product_id`", "qty",
			"3456", "qty+?",
			"3457", "qty+?",
			"3458", "qty+?",
		), dbr.Ints{3, 4, 5})).
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
	//Values: [3 4 5 345 567 897 6]
	//
	//Interpolated Statement:
	//UPDATE `cataloginventory_stock_item` SET `qty`=CASE `product_id` WHEN 3456 THEN
	//qty+3 WHEN 3457 THEN qty+4 WHEN 3458 THEN qty+5 ELSE qty END WHERE (`product_id`
	//IN (345,567,897)) AND (`website_id` = 6)
}

// ExampleSQLCase_select is a duplicate of ExampleSelect_AddArguments
func ExampleSQLCase_select() {
	// time stamp has no special meaning ;-)
	start := dbr.MakeTime(time.Unix(1257894000, 0))
	end := dbr.MakeTime(time.Unix(1257980400, 0))
	s := dbr.NewSelect().AddColumns("price", "sku", "name", "title", "description").
		AddColumnsExprAlias(
			dbr.SQLCase("", "`closed`",
				"date_start <= ? AND date_end >= ?", "`open`",
				"date_start > ? AND date_end > ?", "`upcoming`",
			),
			"is_on_sale",
		).
		AddArguments(start, end, start, end).
		From("catalog_promotions").Where(
		dbr.Column("promotion_id").NotIn().Ints(4711, 815, 42))
	writeToSQLAndInterpolate(s)

	// Output:
	//Prepared Statement:
	//SELECT `price`, `sku`, `name`, `title`, `description`, CASE  WHEN date_start <=
	//? AND date_end >= ? THEN `open` WHEN date_start > ? AND date_end > ? THEN
	//`upcoming` ELSE `closed` END AS `is_on_sale` FROM `catalog_promotions` WHERE
	//(`promotion_id` NOT IN (?,?,?))
	//Values: [2009-11-11 00:00:00 +0100 CET 2009-11-12 00:00:00 +0100 CET 2009-11-11 00:00:00 +0100 CET 2009-11-12 00:00:00 +0100 CET 4711 815 42]
	//
	//Interpolated Statement:
	//SELECT `price`, `sku`, `name`, `title`, `description`, CASE  WHEN date_start <=
	//'2009-11-11 00:00:00' AND date_end >= '2009-11-12 00:00:00' THEN `open` WHEN
	//date_start > '2009-11-11 00:00:00' AND date_end > '2009-11-12 00:00:00' THEN
	//`upcoming` ELSE `closed` END AS `is_on_sale` FROM `catalog_promotions` WHERE
	//(`promotion_id` NOT IN (4711,815,42))
}

// ExampleSelect_AddArguments is duplicate of ExampleSQLCase_select
func ExampleSelect_AddArguments() {
	// time stamp has no special meaning ;-)
	start := dbr.MakeTime(time.Unix(1257894000, 0))
	end := dbr.MakeTime(time.Unix(1257980400, 0))
	s := dbr.NewSelect().AddColumns("price", "sku", "name", "title", "description").
		AddColumnsExprAlias(
			dbr.SQLCase("", "`closed`",
				"date_start <= ? AND date_end >= ?", "`open`",
				"date_start > ? AND date_end > ?", "`upcoming`",
			),
			"is_on_sale",
		).
		AddArguments(start, end, start, end).
		From("catalog_promotions").Where(
		dbr.Column("promotion_id").NotIn().Ints(4711, 815, 42))
	writeToSQLAndInterpolate(s)

	// Output:
	//Prepared Statement:
	//SELECT `price`, `sku`, `name`, `title`, `description`, CASE  WHEN date_start <=
	//? AND date_end >= ? THEN `open` WHEN date_start > ? AND date_end > ? THEN
	//`upcoming` ELSE `closed` END AS `is_on_sale` FROM `catalog_promotions` WHERE
	//(`promotion_id` NOT IN (?,?,?))
	//Values: [2009-11-11 00:00:00 +0100 CET 2009-11-12 00:00:00 +0100 CET 2009-11-11 00:00:00 +0100 CET 2009-11-12 00:00:00 +0100 CET 4711 815 42]
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
			dbr.Column("e").String("wat").Or(),
			dbr.ParenthesisClose(),
			dbr.Column("f").Int(2),
		).
		GroupBy("ab").
		Having(
			dbr.Expression("j = k"),
			dbr.ParenthesisOpen(),
			dbr.Column("m").Int(33),
			dbr.Column("n").String("wh3r3").Or(),
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
	//Values: [1 wat 2 33 wh3r3]
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
			Select: dbr.NewSelect().AddColumnsExpr("Month(day_of_sale)", "Sum(amount)").From("sales_days").
				Where(dbr.Expression("Year(day_of_sale) = ?", dbr.Int(2015))).
				GroupByExpr("Month(day_of_sale))"),
		},
		dbr.WithCTE{Name: "best_month", Columns: []string{"month", "total", "award"},
			Select: dbr.NewSelect().AddColumns("month", "total").AddColumnsExpr(`"best"`).From("sales_by_month").
				Where(dbr.SubSelect("total", dbr.Equal, dbr.NewSelect().AddColumnsExpr("Max(total)").From("sales_by_month"))),
		},
		dbr.WithCTE{Name: "worst_month", Columns: []string{"month", "total", "award"},
			Select: dbr.NewSelect().AddColumns("month", "total").AddColumnsExpr(`"worst"`).From("sales_by_month").
				Where(dbr.SubSelect("total", dbr.Equal, dbr.NewSelect().AddColumnsExpr("Min(total)").From("sales_by_month"))),
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
	//Values: [2015]
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
