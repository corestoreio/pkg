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

package dml_test

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/corestoreio/pkg/sql/dml"
	"github.com/corestoreio/pkg/util/strs"
)

func writeToSQL(qb dml.QueryBuilder) {
	sqlStr, args, err := qb.ToSQL()
	if err != nil {
		fmt.Printf("%+v\n", err)
		return
	}
	if len(args) > 0 {
		fmt.Print("Prepared ")
	}
	fmt.Println("Statement:")
	strs.FwordWrap(os.Stdout, sqlStr, 80)
	fmt.Print("\n")
	if len(args) > 0 {
		fmt.Printf("Arguments: %v\n\n", args)
	}
}

func writeToSQLAndInterpolate(qb dml.QueryBuilder) {
	sqlStr, args, err := qb.ToSQL()
	if err != nil {
		fmt.Printf("%+v\n", err)
		return
	}
	if len(args) > 0 {
		fmt.Print("Prepared ")
	}
	fmt.Println("Statement:")
	strs.FwordWrap(os.Stdout, sqlStr, 80)
	fmt.Print("\n")
	if len(args) > 0 {
		fmt.Printf("Arguments: %v\n\n", args)
	} else {
		return
	}

	switch dmlArg := qb.(type) {
	case *dml.Artisan:
		prev := dmlArg.Options
		qb = dmlArg.Interpolate()
		defer func() { dmlArg.Options = prev; qb = dmlArg }()
	default:
		panic(fmt.Sprintf("func compareToSQL: the type %#v is not (yet) supported.", qb))
	}

	sqlStr, args, err = qb.ToSQL()
	if err != nil {
		fmt.Printf("%+v\n", err)
		return
	}
	if len(args) > 0 {
		panic(fmt.Sprintf("func compareToSQL should not return arguments when interpolation is enabled, got: %#v\n\n", args))
	}
	fmt.Println("Interpolated Statement:")
	strs.FwordWrap(os.Stdout, sqlStr, 80)
}

func ExampleNewInsert() {
	i := dml.NewInsert("tableA").
		AddColumns("b", "c", "d", "e").SetRowCount(2).WithArgs().
		Int(1).Int(2).String("Three").Null().
		Int(5).Int(6).String("Seven").Float64(3.14156)
	writeToSQLAndInterpolate(i)

	// Output:
	// Prepared Statement:
	// INSERT INTO `tableA` (`b`,`c`,`d`,`e`) VALUES (?,?,?,?),(?,?,?,?)
	// Arguments: [1 2 Three <nil> 5 6 Seven 3.14156]
	//
	// Interpolated Statement:
	// INSERT INTO `tableA` (`b`,`c`,`d`,`e`) VALUES
	//(1,2,'Three',NULL),(5,6,'Seven',3.14156)
}

func ExampleInsert_SetRowCount() {
	// RowCount of 4 allows to insert four rows with a single INSERT query.
	// Useful when creating prepared statements.
	i := dml.NewInsert("dml_people").AddColumns("id", "name", "email").SetRowCount(4).BuildValues()
	writeToSQLAndInterpolate(i)

	// Output:
	// Statement:
	// INSERT INTO `dml_people` (`id`,`name`,`email`) VALUES
	//(?,?,?),(?,?,?),(?,?,?),(?,?,?)
}

func ExampleInsert_SetRowCount_withdata() {
	i := dml.NewInsert("catalog_product_link").SetRowCount(3).WithArgs().
		Int(2046).Int(33).Int(3).
		Int(2046).Int(34).Int(3).
		Int(2046).Int(35).Int(3)
	writeToSQLAndInterpolate(i)

	// Output:
	// Prepared Statement:
	// INSERT INTO `catalog_product_link` VALUES (?,?,?),(?,?,?),(?,?,?)
	// Arguments: [2046 33 3 2046 34 3 2046 35 3]
	//
	// Interpolated Statement:
	// INSERT INTO `catalog_product_link` VALUES (2046,33,3),(2046,34,3),(2046,35,3)
}

// ExampleInsert_WithArgs_rawData cannot interpolate because raw interfaces are
// not supported.
func ExampleInsert_WithArgs_rawData() {
	// Without any columns you must for each row call AddArgs. Here we insert
	// three rows at once.
	i := dml.NewInsert("catalog_product_link").SetRowCount(3).WithArgs().Raw(
		2046, 33, 3,
		2046, 34, 3,
		2046, 35, 3,
	)
	writeToSQL(i)

	// Specifying columns allows to call only one time AddArgs but inserting
	// three rows at once. Of course you can also insert only one row ;-)
	i = dml.NewInsert("catalog_product_link").
		AddColumns("product_id", "linked_product_id", "link_type_id").
		WithArgs().Raw(
		2046, 33, 3,
		2046, 34, 3,
		2046, 35, 3)
	writeToSQL(i)

	// Output:
	// Prepared Statement:
	// INSERT INTO `catalog_product_link` VALUES (?,?,?),(?,?,?),(?,?,?)
	// Arguments: [2046 33 3 2046 34 3 2046 35 3]
	//
	// Prepared Statement:
	// INSERT INTO `catalog_product_link`
	//(`product_id`,`linked_product_id`,`link_type_id`) VALUES (?,?,?),(?,?,?),(?,?,?)
	// Arguments: [2046 33 3 2046 34 3 2046 35 3]
}

// ExampleInsert_AddOnDuplicateKey this example assumes you are not using a any
// place holders. Be aware of SQL injections.
func ExampleInsert_AddOnDuplicateKey() {
	i := dml.NewInsert("dml_people").
		AddColumns("id", "name", "email").
		AddOnDuplicateKey(
			dml.Column("name").Str("Pik3"),
			dml.Column("email").Values(),
		).WithArgs().Int(1).String("Pik'e").String("pikes@peak.com")
	writeToSQLAndInterpolate(i)

	// Output:
	// Prepared Statement:
	// INSERT INTO `dml_people` (`id`,`name`,`email`) VALUES (?,?,?) ON DUPLICATE KEY
	// UPDATE `name`='Pik3', `email`=VALUES(`email`)
	// Arguments: [1 Pik'e pikes@peak.com]
	//
	// Interpolated Statement:
	// INSERT INTO `dml_people` (`id`,`name`,`email`) VALUES
	//(1,'Pik\'e','pikes@peak.com') ON DUPLICATE KEY UPDATE `name`='Pik3',
	//`email`=VALUES(`email`)
}

func ExampleInsert_FromSelect_withPlaceHolders() {
	ins := dml.NewInsert("tableA").FromSelect(
		dml.NewSelect().AddColumns("something_id", "user_id").
			AddColumns("other").
			From("some_table").
			Where(
				dml.ParenthesisOpen(),
				dml.Column("int64A").GreaterOrEqual().PlaceHolder(),
				dml.Column("string").Str("wat").Or(),
				dml.ParenthesisClose(),
				dml.Column("int64B").In().NamedArg("i64BIn"),
			).
			OrderByDesc("id").
			Paginate(1, 20),
	).WithArgs().Int64(4).NamedArg("i64BIn", []int64{9, 8, 7})
	writeToSQLAndInterpolate(ins)
	// Output:
	// Prepared Statement:
	// INSERT INTO `tableA` SELECT `something_id`, `user_id`, `other` FROM `some_table`
	// WHERE ((`int64A` >= ?) OR (`string` = 'wat')) AND (`int64B` IN ?) ORDER BY `id`
	// DESC LIMIT 0,20
	// Arguments: [4 9 8 7]
	//
	// Interpolated Statement:
	// INSERT INTO `tableA` SELECT `something_id`, `user_id`, `other` FROM `some_table`
	// WHERE ((`int64A` >= 4) OR (`string` = 'wat')) AND (`int64B` IN (9,8,7)) ORDER BY
	//`id` DESC LIMIT 0,20
}

func ExampleInsert_FromSelect_withoutPlaceHolders() {
	ins := dml.NewInsert("tableA")

	ins.FromSelect(
		dml.NewSelect().AddColumns("something_id", "user_id").
			AddColumns("other").
			From("some_table").
			Where(
				dml.ParenthesisOpen(),
				dml.Column("int64A").GreaterOrEqual().Int64(1),
				dml.Column("string").Str("wat").Or(),
				dml.ParenthesisClose(),
				dml.Column("int64B").In().Int64s(1, 2, 3),
			).
			OrderByDesc("id").
			Paginate(1, 20),
	)
	writeToSQLAndInterpolate(ins)
	// Output:
	// Statement:
	// INSERT INTO `tableA` SELECT `something_id`, `user_id`, `other` FROM `some_table`
	// WHERE ((`int64A` >= 1) OR (`string` = 'wat')) AND (`int64B` IN (1,2,3)) ORDER BY
	//`id` DESC LIMIT 0,20
}

// ExampleInsert_WithPairs this example uses WithArgs to build the final SQL
// string.
func ExampleInsert_WithPairs() {
	ins := dml.NewInsert("catalog_product_link").
		WithPairs(
			// First row
			dml.Column("product_id").Int64(2046),
			dml.Column("linked_product_id").Int64(33),
			dml.Column("link_type_id").Int64(3),

			// second row
			dml.Column("product_id").Int64(2046),
			dml.Column("linked_product_id").Int64(34),
			dml.Column("link_type_id").Int64(3),
		).WithArgs()
	writeToSQLAndInterpolate(ins)
	// Output:
	// Prepared Statement:
	// INSERT INTO `catalog_product_link`
	//(`product_id`,`linked_product_id`,`link_type_id`) VALUES (?,?,?),(?,?,?)
	// Arguments: [2046 33 3 2046 34 3]
	//
	// Interpolated Statement:
	// INSERT INTO `catalog_product_link`
	//(`product_id`,`linked_product_id`,`link_type_id`) VALUES (2046,33,3),(2046,34,3)
}

// ExampleInsert_BuildValues does not call WithArgs but call to BuildValues must
// be made to enable building the VALUES part.
func ExampleInsert_BuildValues() {
	ins := dml.NewInsert("catalog_product_link").
		WithPairs(
			// First row
			dml.Column("product_id").Int64(2046),
			dml.Column("linked_product_id").Int64(33),
			dml.Column("link_type_id").Int64(3),

			// second row
			dml.Column("product_id").Int64(2046),
			dml.Column("linked_product_id").Int64(34),
			dml.Column("link_type_id").Int64(3),
		).BuildValues()
	writeToSQLAndInterpolate(ins)
	// Output:
	// Statement:
	// INSERT INTO `catalog_product_link`
	//(`product_id`,`linked_product_id`,`link_type_id`) VALUES (2046,33,3),(2046,34,3)
}

// ExampleInsert_expressionInVALUES contains an expression in the VALUES part.
// You must provide the column names.
func ExampleInsert_expressionInVALUES() {
	ins := dml.NewInsert("catalog_product_customer_relation").
		AddColumns("product_id", "sort_order").
		WithPairs(
			dml.Column("customer_id").Expr("IFNULL(SELECT entity_id FROM customer_entity WHERE email like ?,0)"),
			dml.Column("group_id").Sub(
				dml.NewSelect("group_id").From("customer_group").Where(
					dml.Column("name").Equal().PlaceHolder(),
				),
			),
		).BuildValues()
	writeToSQLAndInterpolate(ins)
	// Output:
	// Statement:
	// INSERT INTO `catalog_product_customer_relation`
	//(`product_id`,`sort_order`,`customer_id`,`group_id`) VALUES (?,?,IFNULL(SELECT
	// entity_id FROM customer_entity WHERE email like ?,0),(SELECT `group_id` FROM
	//`customer_group` WHERE (`name` = ?)))
}

func ExampleDelete() {
	d := dml.NewDelete("tableA").Where(
		dml.Column("a").Like().Str("b'%"),
		dml.Column("b").In().Ints(3, 4, 5, 6),
	).
		Limit(1).OrderBy("id")
	writeToSQLAndInterpolate(d)
	// Output:
	// Statement:
	// DELETE FROM `tableA` WHERE (`a` LIKE 'b\'%') AND (`b` IN (3,4,5,6)) ORDER BY
	//`id` LIMIT 1
}

func ExampleDelete_FromTables() {
	d := dml.NewDelete("customer_entity").Alias("ce").
		FromTables("customer_address", "customer_company").
		Join(
			dml.MakeIdentifier("customer_company").Alias("cc"),
			dml.Columns("ce.entity_id", "cc.customer_id"),
		).
		RightJoin(
			dml.MakeIdentifier("customer_address").Alias("ca"),
			dml.Column("ce.entity_id").Equal().Column("ca.parent_id"),
		).
		Where(
			dml.Column("ce.created_at").Less().PlaceHolder(),
		).
		Limit(1).OrderBy("id")
	writeToSQLAndInterpolate(d)
	// Output:
	// Statement:
	// DELETE `ce`,`customer_address`,`customer_company` FROM `customer_entity` AS `ce`
	// INNER JOIN `customer_company` AS `cc` USING (`ce.entity_id`,`cc.customer_id`)
	// RIGHT JOIN `customer_address` AS `ca` ON (`ce`.`entity_id` = `ca`.`parent_id`)
	// WHERE (`ce`.`created_at` < ?) ORDER BY `id` LIMIT 1
}

// ExampleNewUnion constructs a UNION with three SELECTs. It preserves the
// results sets of each SELECT by simply adding an internal index to the columns
// list and sort ascending with the internal index.
func ExampleNewUnion() {
	u := dml.NewUnion(
		dml.NewSelect().AddColumnsAliases("a1", "A", "a2", "B").From("tableA").Where(dml.Column("a1").Int64(3)),
		dml.NewSelect().AddColumnsAliases("b1", "A", "b2", "B").From("tableB").Where(dml.Column("b1").Int64(4)),
	)
	// Maybe more of your code ...
	u.Append(
		dml.NewSelect().AddColumnsConditions(
			dml.Expr("concat(c1,?,c2)").Alias("A").Str("-"),
		).
			AddColumnsAliases("c2", "B").
			From("tableC").Where(dml.Column("c2").Str("ArgForC2")),
	).
		OrderBy("A").       // Ascending by A
		OrderByDesc("B").   // Descending by B
		All().              // Enables UNION ALL syntax
		PreserveResultSet() // Maintains the correct order of the result set for all SELECTs.
	// Note that the final ORDER BY statement of a UNION creates a temporary
	// table in MySQL.
	writeToSQLAndInterpolate(u)
	// Output:
	// Statement:
	//(SELECT `a1` AS `A`, `a2` AS `B`, 0 AS `_preserve_result_set` FROM `tableA`
	// WHERE (`a1` = 3))
	// UNION ALL
	//(SELECT `b1` AS `A`, `b2` AS `B`, 1 AS `_preserve_result_set` FROM `tableB`
	// WHERE (`b1` = 4))
	// UNION ALL
	//(SELECT concat(c1,'-',c2) AS `A`, `c2` AS `B`, 2 AS `_preserve_result_set` FROM
	//`tableC` WHERE (`c2` = 'ArgForC2'))
	// ORDER BY `_preserve_result_set`, `A`, `B` DESC
}

// ExampleNewUnion_template interpolates the SQL string with its placeholders
// and puts for each placeholder the correct encoded and escaped value into it.
// Eliminates the need for prepared statements. Avoids an additional round trip
// to the database server by sending the query and its arguments directly. If
// you execute a query multiple times within a short time you should use
// prepared statements.
func ExampleNewUnion_template() {
	u := dml.NewUnion(
		dml.NewSelect().AddColumns("t.value", "t.attribute_id", "t.store_id").
			FromAlias("catalog_product_entity_$type$", "t").
			Where(dml.Column("entity_id").Int64(1561), dml.Column("store_id").In().Int64s(1, 0)),
	).
		StringReplace("$type$", "varchar", "int", "decimal", "datetime", "text").
		PreserveResultSet().
		All().OrderBy("attribute_id", "store_id")
	writeToSQLAndInterpolate(u)
	// Output:
	// Statement:
	//(SELECT `t`.`value`, `t`.`attribute_id`, `t`.`store_id`, 0 AS
	//`_preserve_result_set` FROM `catalog_product_entity_varchar` AS `t` WHERE
	//(`entity_id` = 1561) AND (`store_id` IN (1,0)))
	// UNION ALL
	//(SELECT `t`.`value`, `t`.`attribute_id`, `t`.`store_id`, 1 AS
	//`_preserve_result_set` FROM `catalog_product_entity_int` AS `t` WHERE
	//(`entity_id` = 1561) AND (`store_id` IN (1,0)))
	// UNION ALL
	//(SELECT `t`.`value`, `t`.`attribute_id`, `t`.`store_id`, 2 AS
	//`_preserve_result_set` FROM `catalog_product_entity_decimal` AS `t` WHERE
	//(`entity_id` = 1561) AND (`store_id` IN (1,0)))
	// UNION ALL
	//(SELECT `t`.`value`, `t`.`attribute_id`, `t`.`store_id`, 3 AS
	//`_preserve_result_set` FROM `catalog_product_entity_datetime` AS `t` WHERE
	//(`entity_id` = 1561) AND (`store_id` IN (1,0)))
	// UNION ALL
	//(SELECT `t`.`value`, `t`.`attribute_id`, `t`.`store_id`, 4 AS
	//`_preserve_result_set` FROM `catalog_product_entity_text` AS `t` WHERE
	//(`entity_id` = 1561) AND (`store_id` IN (1,0)))
	// ORDER BY `_preserve_result_set`, `attribute_id`, `store_id`
}

func ExampleInterpolate() {
	sqlStr := dml.Interpolate("SELECT * FROM x WHERE a IN ? AND b IN ? AND c NOT IN ? AND d BETWEEN ? AND ?").
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

func ExampleExpandPlaceHolders() {
	cp, err := dml.NewConnPool()
	if err != nil {
		panic(err)
	}
	sqlStr, args, err := cp.WithRawSQL("SELECT * FROM `table` WHERE id IN ? AND name IN ?").
		ExpandPlaceHolders().
		Ints(5, 7, 9).Strings("a", "b", "c", "d", "e").ToSQL()
	if err != nil {
		fmt.Printf("%+v\n", err)
		return
	}

	fmt.Printf("%s\nArguments: %v\n", sqlStr, args)
	// Output:
	// SELECT * FROM `table` WHERE id IN (?,?,?) AND name IN (?,?,?,?,?)
	// Arguments: [5 7 9 a b c d e]
}

func argPrinter(wf *dml.Condition) {
	sqlStr, args, err := dml.NewSelect().AddColumns("a", "b").
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

// ExampleCondition is duplicate of ExampleColumn
func ExampleCondition() {
	argPrinter(dml.Column("d").Null())
	argPrinter(dml.Column("d").NotNull())
	argPrinter(dml.Column("d").Int(2))
	argPrinter(dml.Column("d").Int(3).Null())
	argPrinter(dml.Column("d").Int(4).NotNull())
	argPrinter(dml.Column("d").In().Ints(7, 8, 9))
	argPrinter(dml.Column("d").NotIn().Ints(10, 11, 12))
	argPrinter(dml.Column("d").Between().Ints(13, 14))
	argPrinter(dml.Column("d").NotBetween().Ints(15, 16))
	argPrinter(dml.Column("d").Greatest().Ints(17, 18, 19))
	argPrinter(dml.Column("d").Least().Ints(20, 21, 22))
	argPrinter(dml.Column("d").Equal().Int(30))
	argPrinter(dml.Column("d").NotEqual().Int(31))
	argPrinter(dml.Column("alias.column").SpaceShip().Float64(3.14159))

	argPrinter(dml.Column("d").Less().Int(32))
	argPrinter(dml.Column("d").Greater().Int(33))
	argPrinter(dml.Column("d").LessOrEqual().Int(34))
	argPrinter(dml.Column("d").GreaterOrEqual().Int(35))

	argPrinter(dml.Column("d").Like().Str("Goph%"))
	argPrinter(dml.Column("d").NotLike().Str("Cat%"))

	// Output:
	//"SELECT `a`, `b` FROM `c` WHERE (`d` IS NULL)"
	//"SELECT `a`, `b` FROM `c` WHERE (`d` IS NOT NULL)"
	//"SELECT `a`, `b` FROM `c` WHERE (`d` = 2)"
	//"SELECT `a`, `b` FROM `c` WHERE (`d` IS NULL)"
	//"SELECT `a`, `b` FROM `c` WHERE (`d` IS NOT NULL)"
	//"SELECT `a`, `b` FROM `c` WHERE (`d` IN (7,8,9))"
	//"SELECT `a`, `b` FROM `c` WHERE (`d` NOT IN (10,11,12))"
	//"SELECT `a`, `b` FROM `c` WHERE (`d` BETWEEN 13 AND 14)"
	//"SELECT `a`, `b` FROM `c` WHERE (`d` NOT BETWEEN 15 AND 16)"
	//"SELECT `a`, `b` FROM `c` WHERE (`d` GREATEST (17,18,19))"
	//"SELECT `a`, `b` FROM `c` WHERE (`d` LEAST (20,21,22))"
	//"SELECT `a`, `b` FROM `c` WHERE (`d` = 30)"
	//"SELECT `a`, `b` FROM `c` WHERE (`d` != 31)"
	//"SELECT `a`, `b` FROM `c` WHERE (`alias`.`column` <=> 3.14159)"
	//"SELECT `a`, `b` FROM `c` WHERE (`d` < 32)"
	//"SELECT `a`, `b` FROM `c` WHERE (`d` > 33)"
	//"SELECT `a`, `b` FROM `c` WHERE (`d` <= 34)"
	//"SELECT `a`, `b` FROM `c` WHERE (`d` >= 35)"
	//"SELECT `a`, `b` FROM `c` WHERE (`d` LIKE 'Goph%')"
	//"SELECT `a`, `b` FROM `c` WHERE (`d` NOT LIKE 'Cat%')"
}

// ExampleColumn is a duplicate of ExampleCondition
func ExampleColumn() {
	argPrinter(dml.Column("d").Null())
	argPrinter(dml.Column("d").NotNull())
	argPrinter(dml.Column("d").Int(2))
	argPrinter(dml.Column("d").Int(3).Null())
	argPrinter(dml.Column("d").Int(4).NotNull())
	argPrinter(dml.Column("d").In().Ints(7, 8, 9))
	argPrinter(dml.Column("d").NotIn().Ints(10, 11, 12))
	argPrinter(dml.Column("d").Between().Ints(13, 14))
	argPrinter(dml.Column("d").NotBetween().Ints(15, 16))
	argPrinter(dml.Column("d").Greatest().Ints(17, 18, 19))
	argPrinter(dml.Column("d").Least().Ints(20, 21, 22))
	argPrinter(dml.Column("d").Equal().Int(30))
	argPrinter(dml.Column("d").NotEqual().Int(31))
	argPrinter(dml.Column("alias.column").SpaceShip().Float64(3.14159))

	argPrinter(dml.Column("d").Less().Int(32))
	argPrinter(dml.Column("d").Greater().Int(33))
	argPrinter(dml.Column("d").LessOrEqual().Int(34))
	argPrinter(dml.Column("d").GreaterOrEqual().Int(35))

	argPrinter(dml.Column("d").Like().Str("Goph%"))
	argPrinter(dml.Column("d").NotLike().Str("Cat%"))

	// Output:
	//"SELECT `a`, `b` FROM `c` WHERE (`d` IS NULL)"
	//"SELECT `a`, `b` FROM `c` WHERE (`d` IS NOT NULL)"
	//"SELECT `a`, `b` FROM `c` WHERE (`d` = 2)"
	//"SELECT `a`, `b` FROM `c` WHERE (`d` IS NULL)"
	//"SELECT `a`, `b` FROM `c` WHERE (`d` IS NOT NULL)"
	//"SELECT `a`, `b` FROM `c` WHERE (`d` IN (7,8,9))"
	//"SELECT `a`, `b` FROM `c` WHERE (`d` NOT IN (10,11,12))"
	//"SELECT `a`, `b` FROM `c` WHERE (`d` BETWEEN 13 AND 14)"
	//"SELECT `a`, `b` FROM `c` WHERE (`d` NOT BETWEEN 15 AND 16)"
	//"SELECT `a`, `b` FROM `c` WHERE (`d` GREATEST (17,18,19))"
	//"SELECT `a`, `b` FROM `c` WHERE (`d` LEAST (20,21,22))"
	//"SELECT `a`, `b` FROM `c` WHERE (`d` = 30)"
	//"SELECT `a`, `b` FROM `c` WHERE (`d` != 31)"
	//"SELECT `a`, `b` FROM `c` WHERE (`alias`.`column` <=> 3.14159)"
	//"SELECT `a`, `b` FROM `c` WHERE (`d` < 32)"
	//"SELECT `a`, `b` FROM `c` WHERE (`d` > 33)"
	//"SELECT `a`, `b` FROM `c` WHERE (`d` <= 34)"
	//"SELECT `a`, `b` FROM `c` WHERE (`d` >= 35)"
	//"SELECT `a`, `b` FROM `c` WHERE (`d` LIKE 'Goph%')"
	//"SELECT `a`, `b` FROM `c` WHERE (`d` NOT LIKE 'Cat%')"
}

func ExampleCondition_Sub() {
	s := dml.NewSelect("sku", "type_id").
		From("catalog_product_entity").
		Where(dml.Column("entity_id").In().Sub(
			dml.NewSelect().From("catalog_category_product").
				AddColumns("entity_id").Where(dml.Column("category_id").Int64(234)),
		))
	writeToSQLAndInterpolate(s)
	// Output:
	// Statement:
	// SELECT `sku`, `type_id` FROM `catalog_product_entity` WHERE (`entity_id` IN
	//(SELECT `entity_id` FROM `catalog_category_product` WHERE (`category_id` =
	// 234)))
}

func ExampleNewSelectWithDerivedTable() {
	sel3 := dml.NewSelect().Unsafe().FromAlias("sales_bestsellers_aggregated_daily", "t3").
		AddColumnsAliases("DATE_FORMAT(t3.period, '%Y-%m-01')", "period").
		AddColumns("t3.store_id", "t3.product_id", "t3.product_name").
		AddColumnsAliases("AVG(`t3`.`product_price`)", "avg_price", "SUM(t3.qty_ordered)", "total_qty").
		Where(dml.Column("product_name").Str("Canon%")).
		GroupBy("t3.store_id").
		GroupBy("DATE_FORMAT(t3.period, '%Y-%m-01')").
		GroupBy("t3.product_id", "t3.product_name").
		OrderBy("t3.store_id").
		OrderBy("DATE_FORMAT(t3.period, '%Y-%m-01')").
		OrderByDesc("total_qty")

	sel1 := dml.NewSelectWithDerivedTable(sel3, "t1").
		AddColumns("t1.period", "t1.store_id", "t1.product_id", "t1.product_name", "t1.avg_price", "t1.qty_ordered").
		Where(dml.Column("product_name").Str("Sony%")).
		OrderBy("t1.period", "t1.product_id")
	writeToSQLAndInterpolate(sel1)
	// Output:
	// Statement:
	// SELECT `t1`.`period`, `t1`.`store_id`, `t1`.`product_id`, `t1`.`product_name`,
	//`t1`.`avg_price`, `t1`.`qty_ordered` FROM (SELECT DATE_FORMAT(t3.period,
	//'%Y-%m-01') AS `period`, `t3`.`store_id`, `t3`.`product_id`,
	//`t3`.`product_name`, AVG(`t3`.`product_price`) AS `avg_price`,
	// SUM(t3.qty_ordered) AS `total_qty` FROM `sales_bestsellers_aggregated_daily` AS
	//`t3` WHERE (`product_name` = 'Canon%') GROUP BY `t3`.`store_id`,
	// DATE_FORMAT(t3.period, '%Y-%m-01'), `t3`.`product_id`, `t3`.`product_name` ORDER
	// BY `t3`.`store_id`, DATE_FORMAT(t3.period, '%Y-%m-01'), `total_qty` DESC) AS
	//`t1` WHERE (`product_name` = 'Sony%') ORDER BY `t1`.`period`, `t1`.`product_id`
}

func ExampleSQLIfNull() {
	s := dml.NewSelect().AddColumnsConditions(
		dml.SQLIfNull("column1"),
		dml.SQLIfNull("table1.column1"),
		dml.SQLIfNull("column1", "column2"),
		dml.SQLIfNull("table1.column1", "table2.column2"),
		dml.SQLIfNull("column2", "1/0").Alias("alias"),
		dml.SQLIfNull("SELECT * FROM x", "8").Alias("alias"),
		dml.SQLIfNull("SELECT * FROM x", "9 ").Alias("alias"),
		dml.SQLIfNull("column1", "column2").Alias("alias"),
		dml.SQLIfNull("table1.column1", "table2.column2").Alias("alias"),
		dml.SQLIfNull("table1", "column1", "table2", "column2"),
		dml.SQLIfNull("table1", "column1", "table2", "column2").Alias("alias"),
	).From("table1")
	sStr, _, _ := s.ToSQL()
	fmt.Print(strings.Replace(sStr, ", ", ",\n", -1))

	// Output:
	// SELECT IFNULL(`column1`,NULL),
	// IFNULL(`table1`.`column1`,NULL),
	// IFNULL(`column1`,`column2`),
	// IFNULL(`table1`.`column1`,`table2`.`column2`),
	// IFNULL(`column2`,1/0) AS `alias`,
	// IFNULL(SELECT * FROM x,8) AS `alias`,
	// IFNULL(SELECT * FROM x,9 ) AS `alias`,
	// IFNULL(`column1`,`column2`) AS `alias`,
	// IFNULL(`table1`.`column1`,`table2`.`column2`) AS `alias`,
	// IFNULL(`table1`.`column1`,`table2`.`column2`),
	// IFNULL(`table1`.`column1`,`table2`.`column2`) AS `alias` FROM `table1`
}

func ExampleSQLIf() {
	s := dml.NewSelect().
		AddColumns("a", "b", "c").
		From("table1").
		Where(
			dml.SQLIf("a > 0", "b", "c").Greater().Int(4711),
		)
	writeToSQLAndInterpolate(s)

	// Output:
	// Statement:
	// SELECT `a`, `b`, `c` FROM `table1` WHERE (IF((a > 0), b, c) > 4711)
}

func ExampleSQLCase_update() {
	u := dml.NewUpdate("cataloginventory_stock_item").
		AddClauses(dml.Column("qty").SQLCase("`product_id`", "qty",
			"3456", "qty+?",
			"3457", "qty+?",
			"3458", "qty+?",
		).Int(3).Int(4).Int(5)).
		Where(
			dml.Column("product_id").In().Int64s(345, 567, 897),
			dml.Column("website_id").Int64(6),
		)
	writeToSQLAndInterpolate(u)

	// Output:
	// Statement:
	// UPDATE `cataloginventory_stock_item` SET `qty`=CASE `product_id` WHEN 3456 THEN
	// qty+3 WHEN 3457 THEN qty+4 WHEN 3458 THEN qty+5 ELSE qty END WHERE (`product_id`
	// IN (345,567,897)) AND (`website_id` = 6)
}

// ExampleSQLCase_select is a duplicate of ExampleSelect_AddArguments
func ExampleSQLCase_select() {
	// time stamp has no special meaning ;-)
	start := time.Unix(1257894000, 0).In(time.UTC)
	end := time.Unix(1257980400, 0).In(time.UTC)

	s := dml.NewSelect().AddColumns("price", "sku", "name", "title", "description").
		AddColumnsConditions(
			dml.SQLCase("", "`closed`",
				"date_start <= ? AND date_end >= ?", "`open`",
				"date_start > ? AND date_end > ?", "`upcoming`",
			).Alias("is_on_sale"),
		).
		From("catalog_promotions").Where(
		dml.Column("promotion_id").NotIn().PlaceHolders(3)).
		WithArgs().Time(start).Time(end).Time(start).Time(end).Int(4711).Int(815).Int(42)
	writeToSQLAndInterpolate(s)

	// Output:
	// Prepared Statement:
	// SELECT `price`, `sku`, `name`, `title`, `description`, CASE  WHEN date_start <=
	//? AND date_end >= ? THEN `open` WHEN date_start > ? AND date_end > ? THEN
	//`upcoming` ELSE `closed` END AS `is_on_sale` FROM `catalog_promotions` WHERE
	//(`promotion_id` NOT IN (?,?,?))
	// Arguments: [2009-11-10 23:00:00 +0000 UTC 2009-11-11 23:00:00 +0000 UTC 2009-11-10 23:00:00 +0000 UTC 2009-11-11 23:00:00 +0000 UTC 4711 815 42]
	//
	// Interpolated Statement:
	// SELECT `price`, `sku`, `name`, `title`, `description`, CASE  WHEN date_start <=
	//'2009-11-10 23:00:00' AND date_end >= '2009-11-11 23:00:00' THEN `open` WHEN
	// date_start > '2009-11-10 23:00:00' AND date_end > '2009-11-11 23:00:00' THEN
	//`upcoming` ELSE `closed` END AS `is_on_sale` FROM `catalog_promotions` WHERE
	//(`promotion_id` NOT IN (4711,815,42))
}

// ExampleSelect_AddColumnsConditions is duplicate of ExampleSQLCase_select
func ExampleSelect_AddColumnsConditions() {
	start := time.Unix(1257894000, 0).In(time.UTC)
	end := time.Unix(1257980400, 0).In(time.UTC)

	s := dml.NewSelect().AddColumns("price", "sku", "name", "title", "description").
		AddColumnsConditions(
			dml.SQLCase("", "`closed`",
				"date_start <= ? AND date_end >= ?", "`open`",
				"date_start > ? AND date_end > ?", "`upcoming`",
			).Alias("is_on_sale").Time(start).Time(end).Time(start).Time(end),
		).
		From("catalog_promotions").Where(
		dml.Column("promotion_id").NotIn().Ints(4711, 815, 42))
	writeToSQLAndInterpolate(s.WithArgs())

	// Output:
	// Statement:
	// SELECT `price`, `sku`, `name`, `title`, `description`, CASE  WHEN date_start <=
	//'2009-11-10 23:00:00' AND date_end >= '2009-11-11 23:00:00' THEN `open` WHEN
	// date_start > '2009-11-10 23:00:00' AND date_end > '2009-11-11 23:00:00' THEN
	//`upcoming` ELSE `closed` END AS `is_on_sale` FROM `catalog_promotions` WHERE
	//(`promotion_id` NOT IN (4711,815,42))
}

func ExampleParenthesisOpen() {
	s := dml.NewSelect("columnA", "columnB").
		Distinct().
		FromAlias("tableC", "ccc").
		Where(
			dml.ParenthesisOpen(),
			dml.Column("d").Int(1),
			dml.Column("e").Str("wat").Or(),
			dml.ParenthesisClose(),
			dml.Column("f").Int(2),
		).
		GroupBy("ab").
		Having(
			dml.Expr("j = k"),
			dml.ParenthesisOpen(),
			dml.Column("m").Int(33),
			dml.Column("n").Str("wh3r3").Or(),
			dml.ParenthesisClose(),
		).
		OrderBy("l").
		Limit(8, 7)
	writeToSQLAndInterpolate(s)

	// Output:
	// Statement:
	// SELECT DISTINCT `columnA`, `columnB` FROM `tableC` AS `ccc` WHERE ((`d` = 1) OR
	//(`e` = 'wat')) AND (`f` = 2) GROUP BY `ab` HAVING (j = k) AND ((`m` = 33) OR
	//(`n` = 'wh3r3')) ORDER BY `l` LIMIT 8,7
}

func ExampleWith_Union() {
	// Non-recursive CTE
	// Sales: Find best and worst month:
	cte := dml.NewWith(
		dml.WithCTE{
			Name: "sales_by_month", Columns: []string{"month", "total"},
			Select: dml.NewSelect().Unsafe().AddColumns("Month(day_of_sale)", "Sum(amount)").From("sales_days").
				Where(dml.Expr("Year(day_of_sale) = ?").Int(2015)).
				GroupBy("Month(day_of_sale))"),
		},
		dml.WithCTE{
			Name: "best_month", Columns: []string{"month", "total", "award"},
			Select: dml.NewSelect().Unsafe().AddColumns("month", "total").AddColumns(`"best"`).From("sales_by_month").
				Where(dml.Column("total").Equal().Sub(dml.NewSelect().Unsafe().AddColumns("Max(total)").From("sales_by_month"))),
		},
		dml.WithCTE{
			Name: "worst_month", Columns: []string{"month", "total", "award"},
			Select: dml.NewSelect().Unsafe().AddColumns("month", "total").AddColumns(`"worst"`).From("sales_by_month").
				Where(dml.Column("total").Equal().Sub(dml.NewSelect().Unsafe().AddColumns("Min(total)").From("sales_by_month"))),
		},
	).Union(dml.NewUnion(
		dml.NewSelect().Star().From("best_month"),
		dml.NewSelect().Star().From("worst_month"),
	).All())
	writeToSQLAndInterpolate(cte)

	// Result:
	//+-------+-------+-------+
	//| month | total | award |
	//+-------+-------+-------+
	//|     1 |   300 | best  |
	//|     3 |    11 | worst |
	//+-------+-------+-------+

	// Output:
	// Statement:
	// WITH `sales_by_month` (`month`,`total`) AS (SELECT Month(day_of_sale),
	// Sum(amount) FROM `sales_days` WHERE (Year(day_of_sale) = 2015) GROUP BY
	// Month(day_of_sale))),
	//`best_month` (`month`,`total`,`award`) AS (SELECT `month`, `total`, "best" FROM
	//`sales_by_month` WHERE (`total` = (SELECT Max(total) FROM `sales_by_month`))),
	//`worst_month` (`month`,`total`,`award`) AS (SELECT `month`, `total`, "worst"
	// FROM `sales_by_month` WHERE (`total` = (SELECT Min(total) FROM
	//`sales_by_month`)))
	//(SELECT * FROM `best_month`)
	// UNION ALL
	//(SELECT * FROM `worst_month`)
}
