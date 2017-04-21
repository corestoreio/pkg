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
		AddValues(dbr.ArgInt64(1), dbr.ArgString("Pik'e"), dbr.ArgString("pikes@peak.com")).
		AddOnDuplicateKey("name", dbr.ArgString("Pik3")).
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

	argEq := dbr.Eq{"int64B": dbr.ArgInt64(1, 2, 3).Operator(dbr.In)}

	sqlStr, args, err := ins.FromSelect(dbr.NewSelect().AddColumnsQuoted("something_id,user_id,other").
		From("some_table").
		Where(dbr.Condition("int64A = ? OR string = ?", dbr.ArgInt64(1), dbr.ArgString("wat"))).
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
		dbr.Condition("a", dbr.ArgString("b'%").Operator(dbr.Like)),
		dbr.Condition("b", dbr.ArgInt(3, 4, 5, 6).Operator(dbr.In)),
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
			Where(dbr.Condition("entity_id", dbr.ArgInt64(1561)), dbr.Condition("store_id", dbr.ArgInt64(1, 0).Operator(dbr.In))),
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
			Where(dbr.Condition("entity_id", dbr.ArgInt64(1561)), dbr.Condition("store_id", dbr.ArgInt64(1, 0).Operator(dbr.In))),
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
		dbr.ArgInt(1).Operator(dbr.In),
		dbr.ArgInt(1, 2, 3).Operator(dbr.In),
		dbr.ArgInt64(5, 6, 7).Operator(dbr.In),
		dbr.ArgString("wat", "ok").Operator(dbr.Between),
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
	sl := []string{"a", "b", "c", "d", "e"}

	sqlStr, args, err := dbr.Repeat("SELECT * FROM `table` WHERE id IN (?) AND name IN (?)",
		dbr.ArgInt(5, 7, 9), dbr.ArgString(sl...))

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
			fmt.Printf("%q", sqlStr)
			if len(args) > 0 {
				fmt.Printf(" Arguments: %v", args.Interfaces())
			}
			fmt.Print("\n")
		}
	}

	argPrinter(dbr.ArgNull())
	argPrinter(dbr.ArgNotNull())
	argPrinter(dbr.ArgInt(2))
	argPrinter(dbr.ArgInt(3).Operator(dbr.Null))
	argPrinter(dbr.ArgInt(4).Operator(dbr.NotNull))
	argPrinter(dbr.ArgInt(7, 8, 9).Operator(dbr.In))
	argPrinter(dbr.ArgInt(10, 11, 12).Operator(dbr.NotIn))
	argPrinter(dbr.ArgInt(13, 14).Operator(dbr.Between))
	argPrinter(dbr.ArgInt(15, 16).Operator(dbr.NotBetween))
	argPrinter(dbr.ArgInt(17, 18, 19).Operator(dbr.Greatest))
	argPrinter(dbr.ArgInt(20, 21, 22).Operator(dbr.Least))
	argPrinter(dbr.ArgInt(30).Operator(dbr.Equal))
	argPrinter(dbr.ArgInt(31).Operator(dbr.NotEqual))

	argPrinter(dbr.ArgInt(32).Operator(dbr.Less))
	argPrinter(dbr.ArgInt(33).Operator(dbr.Greater))
	argPrinter(dbr.ArgInt(34).Operator(dbr.LessOrEqual))
	argPrinter(dbr.ArgInt(35).Operator(dbr.GreaterOrEqual))

	argPrinter(dbr.ArgString("Goph%").Operator(dbr.Like))
	argPrinter(dbr.ArgString("Cat%").Operator(dbr.NotLike))

	//Output:
	//"SELECT `a`, `b` FROM `c` WHERE (`d` IS NULL)"
	//"SELECT `a`, `b` FROM `c` WHERE (`d` IS NOT NULL)"
	//"SELECT `a`, `b` FROM `c` WHERE (`d` = ?)" Arguments: [2]
	//"SELECT `a`, `b` FROM `c` WHERE (`d` IS NULL)"
	//"SELECT `a`, `b` FROM `c` WHERE (`d` IS NOT NULL)"
	//"SELECT `a`, `b` FROM `c` WHERE (`d` IN ?)" Arguments: [7 8 9]
	//"SELECT `a`, `b` FROM `c` WHERE (`d` NOT IN ?)" Arguments: [10 11 12]
	//"SELECT `a`, `b` FROM `c` WHERE (`d` BETWEEN ? AND ?)" Arguments: [13 14]
	//"SELECT `a`, `b` FROM `c` WHERE (`d` NOT BETWEEN ? AND ?)" Arguments: [15 16]
	//"SELECT `a`, `b` FROM `c` WHERE (`d` GREATEST (?))" Arguments: [17 18 19]
	//"SELECT `a`, `b` FROM `c` WHERE (`d` LEAST (?))" Arguments: [20 21 22]
	//"SELECT `a`, `b` FROM `c` WHERE (`d` = ?)" Arguments: [30]
	//"SELECT `a`, `b` FROM `c` WHERE (`d` != ?)" Arguments: [31]
	//"SELECT `a`, `b` FROM `c` WHERE (`d` < ?)" Arguments: [32]
	//"SELECT `a`, `b` FROM `c` WHERE (`d` > ?)" Arguments: [33]
	//"SELECT `a`, `b` FROM `c` WHERE (`d` <= ?)" Arguments: [34]
	//"SELECT `a`, `b` FROM `c` WHERE (`d` >= ?)" Arguments: [35]
	//"SELECT `a`, `b` FROM `c` WHERE (`d` LIKE ?)" Arguments: [Goph%]
	//"SELECT `a`, `b` FROM `c` WHERE (`d` NOT LIKE ?)" Arguments: [Cat%]

}

func ExampleSubSelect() {
	s := dbr.NewSelect("*").
		From("catalog_product_entity").
		Where(dbr.SubSelect(
			"entity_id", dbr.In,
			dbr.NewSelect().From("catalog_category_product").
				AddColumnsQuoted("entity_id").Where(dbr.Condition("category_id", dbr.ArgInt64(234))),
		))

	sqlStr, args, err := s.ToSQL()
	if err != nil {
		fmt.Printf("%+v\n", err)
		return
	}
	fmt.Printf("%s\nArguments: %v\n", sqlStr, args)
	// Output:
	//SELECT * FROM `catalog_product_entity` WHERE (`entity_id` IN (SELECT `entity_id` FROM `catalog_category_product` WHERE (`category_id` = ?)))
	//Arguments: [234]
}

func ExampleNewSelectFromSub() {
	sel3 := dbr.NewSelect().From("sales_bestsellers_aggregated_daily", "t3").
		AddColumnsExprAlias("DATE_FORMAT(t3.period, '%Y-%m-01')", "period").
		AddColumns("`t3`.`store_id`,`t3`.`product_id`,`t3`.`product_name`").
		AddColumnsExprAlias("AVG(`t3`.`product_price`)", "avg_price", "SUM(t3.qty_ordered)", "total_qty").
		Where(dbr.Condition("product_name", dbr.ArgString("Canon%"))).
		GroupBy("`t3`.`store_id`", "DATE_FORMAT(t3.period, '%Y-%m-01')", "`t3`.`product_id`", "`t3`.`product_name`").
		OrderBy("`t3`.`store_id`", "DATE_FORMAT(t3.period, '%Y-%m-01')", "`total_qty` DESC")

	sel1 := dbr.NewSelectFromSub(sel3, "t1").
		AddColumns("`t1`.`period`,`t1`.`store_id`,`t1`.`product_id`,`t1`.`product_name`,`t1`.`avg_price`,`t1`.`qty_ordered`").
		Where(dbr.Condition("product_name", dbr.ArgString("Sony%"))).
		OrderBy("`t1`.period", "`t1`.product_id")

	sqlStr, args, err := sel1.ToSQL()
	if err != nil {
		fmt.Printf("%+v\n", err)
		return
	}
	// To see the correct SQL in a nicely formatted way use an online SQL formatter ;-)
	fmt.Printf("%s\nArguments: %v\n", sqlStr, args)
	// Output:
	//SELECT `t1`.`period`, `t1`.`store_id`, `t1`.`product_id`, `t1`.`product_name`, `t1`.`avg_price`, `t1`.`qty_ordered` FROM (SELECT DATE_FORMAT(t3.period, '%Y-%m-01') AS `period`, `t3`.`store_id`, `t3`.`product_id`, `t3`.`product_name`, AVG(`t3`.`product_price`) AS `avg_price`, SUM(t3.qty_ordered) AS `total_qty` FROM `sales_bestsellers_aggregated_daily` AS `t3` WHERE (`product_name` = ?) GROUP BY `t3`.`store_id`, DATE_FORMAT(t3.period, '%Y-%m-01'), `t3`.`product_id`, `t3`.`product_name` ORDER BY `t3`.`store_id`, DATE_FORMAT(t3.period, '%Y-%m-01'), `total_qty` DESC) AS `t1` WHERE (`product_name` = ?) ORDER BY `t1`.period, `t1`.product_id
	//Arguments: [Canon% Sony%]
}

func ExampleIfNull() {
	fmt.Println(dbr.IfNull("column1"))
	fmt.Println(dbr.IfNull("table1.column1"))
	fmt.Println(dbr.IfNull("column1", "column2"))
	fmt.Println(dbr.IfNull("table1.column1", "table2.column2"))
	fmt.Println(dbr.IfNull("column2", "1/0", "alias"))
	fmt.Println(dbr.IfNull("SELECT * FROM x", "8", "alias"))
	fmt.Println(dbr.IfNull("SELECT * FROM x", "9 ", "alias"))
	fmt.Println(dbr.IfNull("column1", "column2", "alias"))
	fmt.Println(dbr.IfNull("table1.column1", "table2.column2", "alias"))
	fmt.Println(dbr.IfNull("table1", "column1", "table2", "column2"))
	fmt.Println(dbr.IfNull("table1", "column1", "table2", "column2", "alias"))
	fmt.Println(dbr.IfNull("table1", "column1", "table2", "column2", "alias", "x"))
	fmt.Println(dbr.IfNull("table1", "column1", "table2", "column2", "alias", "x", "y"))
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
