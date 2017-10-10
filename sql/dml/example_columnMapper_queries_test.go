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

package dml_test

import "github.com/corestoreio/csfw/sql/dml"

var cmCustomers = &customerCollection{
	Data: []*customerEntity{
		{EntityID: 11, Firstname: "Karl Gopher", StoreID: 0x7, LifetimeSales: dml.MakeNullFloat64(47.11), VoucherCodes: exampleStringSlice{"1FE9983E", "28E76FBC"}},
		{EntityID: 12, Firstname: "Fung Go Roo", StoreID: 0x7, LifetimeSales: dml.MakeNullFloat64(28.94), VoucherCodes: exampleStringSlice{"4FE7787E", "15E59FBB", "794EFDE8"}},
		{EntityID: 13, Firstname: "John Doe", StoreID: 0x6, LifetimeSales: dml.MakeNullFloat64(138.54), VoucherCodes: exampleStringSlice{""}},
	},
}

// ExampleColumnMapper_selectWhereInCollection uses a customer collection to
// retrieve all entity_ids to be used in an IN condition. The customer
// collection does not get qualified because SELECT happens from one table
// without an alias.
func ExampleColumnMapper_selectWhereInCollection() {

	q := dml.NewSelect("entity_id", "firstname", "lifetime_sales").From("customer_entity").
		Where(
			dml.Column("entity_id").In().PlaceHolder(),
		).
		// for variable customers see ExampleColumnMapper
		BindRecord(dml.Qualify("", cmCustomers))
	writeToSQLAndInterpolate(q)

	// Output:
	//Prepared Statement:
	//SELECT `entity_id`, `firstname`, `lifetime_sales` FROM `customer_entity` WHERE
	//(`entity_id` IN (?))
	//Arguments: [11 12 13]
	//
	//Interpolated Statement:
	//SELECT `entity_id`, `firstname`, `lifetime_sales` FROM `customer_entity` WHERE
	//(`entity_id` IN (11,12,13))
}

// ExampleColumnMapper_selectJoinCollection uses a qualified customer
// collection. The qualifier maps to the alias name of the customer_entity
// table.
func ExampleColumnMapper_selectJoinCollection() {

	q := dml.NewSelect("ce.entity_id", "ce.firstname", "cg.customer_group_code", "cg.tax_class_id").FromAlias("customer_entity", "ce").
		Join(dml.MakeIdentifier("customer_group").Alias("cg"),
			dml.Column("ce.group_id").Equal().Column("cg.customer_group_id"),
		).
		Where(
			dml.Column("ce.entity_id").In().PlaceHolder(),
		).
		BindRecord(dml.Qualify("ce", cmCustomers))

	writeToSQLAndInterpolate(q)

	// Output:
	//Prepared Statement:
	//SELECT `ce`.`entity_id`, `ce`.`firstname`, `cg`.`customer_group_code`,
	//`cg`.`tax_class_id` FROM `customer_entity` AS `ce` INNER JOIN `customer_group`
	//AS `cg` ON (`ce`.`group_id` = `cg`.`customer_group_id`) WHERE (`ce`.`entity_id`
	//IN (?))
	//Arguments: [11 12 13]
	//
	//Interpolated Statement:
	//SELECT `ce`.`entity_id`, `ce`.`firstname`, `cg`.`customer_group_code`,
	//`cg`.`tax_class_id` FROM `customer_entity` AS `ce` INNER JOIN `customer_group`
	//AS `cg` ON (`ce`.`group_id` = `cg`.`customer_group_id`) WHERE (`ce`.`entity_id`
	//IN (11,12,13))
}

// ExampleColumnMapper_updateEntity updates an entity with the defined fields.
func ExampleColumnMapper_updateEntity() {

	q := dml.NewUpdate("customer_entity").AddColumns("firstname", "lifetime_sales", "voucher_codes").
		BindRecord(dml.Qualify("", cmCustomers.Data[0])).
		Where(dml.Column("entity_id").Equal().PlaceHolder())

	writeToSQLAndInterpolate(q)
	// Output:
	//Prepared Statement:
	//UPDATE `customer_entity` SET `firstname`=?, `lifetime_sales`=?,
	//`voucher_codes`=? WHERE (`entity_id` = ?)
	//Arguments: [Karl Gopher 47.11 1FE9983E|28E76FBC 11]
	//
	//Interpolated Statement:
	//UPDATE `customer_entity` SET `firstname`='Karl Gopher', `lifetime_sales`=47.11,
	//`voucher_codes`='1FE9983E|28E76FBC' WHERE (`entity_id` = 11)
}

// ExampleColumnMapper_insertEntities inserts multiple entities into a table.
// Collection not yet supported.
func ExampleColumnMapper_insertEntitiesWithColumns() {

	q := dml.NewInsert("customer_entity").AddColumns("firstname", "lifetime_sales", "store_id", "voucher_codes").
		BindRecord(cmCustomers.Data[0], cmCustomers.Data[1], cmCustomers.Data[2])

	writeToSQLAndInterpolate(q)
	// Output:
	//Prepared Statement:
	//INSERT INTO `customer_entity`
	//(`firstname`,`lifetime_sales`,`store_id`,`voucher_codes`) VALUES
	//(?,?,?,?),(?,?,?,?),(?,?,?,?)
	//Arguments: [Karl Gopher 47.11 7 1FE9983E|28E76FBC Fung Go Roo 28.94 7 4FE7787E|15E59FBB|794EFDE8 John Doe 138.54 6 ]
	//
	//Interpolated Statement:
	//INSERT INTO `customer_entity`
	//(`firstname`,`lifetime_sales`,`store_id`,`voucher_codes`) VALUES ('Karl
	//Gopher',47.11,7,'1FE9983E|28E76FBC'),('Fung Go
	//Roo',28.94,7,'4FE7787E|15E59FBB|794EFDE8'),('John Doe',138.54,6,'')
}

// ExampleColumnMapper_insertEntitiesWithoutColumns inserts multiple entities
// into a table. It includes all fields in the sruct. In this case 5 fields
// including the autoincrement field.
func ExampleColumnMapper_insertEntitiesWithoutColumns() {

	q := dml.NewInsert("customer_entity").
		// SetRecordPlaceHolderCount mandatory because no columns provided!
		// customerEntity has five fields and all fields are requested. For
		// now a hard coded 5.
		SetRecordPlaceHolderCount(5).
		BindRecord(cmCustomers.Data[0], cmCustomers.Data[1], cmCustomers.Data[2])

	writeToSQLAndInterpolate(q)
	// Output:
	//Prepared Statement:
	//INSERT INTO `customer_entity` VALUES (?,?,?,?,?),(?,?,?,?,?),(?,?,?,?,?)
	//Arguments: [11 Karl Gopher 7 47.11 1FE9983E|28E76FBC 12 Fung Go Roo 7 28.94 4FE7787E|15E59FBB|794EFDE8 13 John Doe 6 138.54 ]
	//
	//Interpolated Statement:
	//INSERT INTO `customer_entity` VALUES (11,'Karl
	//Gopher',7,47.11,'1FE9983E|28E76FBC'),(12,'Fung Go
	//Roo',7,28.94,'4FE7787E|15E59FBB|794EFDE8'),(13,'John Doe',6,138.54,'')
}

func ExampleColumnMapper_insertCollectionWithoutColumns() {

	q := dml.NewInsert("customer_entity"). //AddColumns("firstname", "lifetime_sales", "store_id", "voucher_codes").
						SetRecordPlaceHolderCount(5).
						SetRowCount(len(cmCustomers.Data)).BindRecord(cmCustomers)

	writeToSQLAndInterpolate(q)
	// Output:
	//Prepared Statement:
	//INSERT INTO `customer_entity` VALUES (?,?,?,?,?),(?,?,?,?,?),(?,?,?,?,?)
	//Arguments: [11 Karl Gopher 7 47.11 1FE9983E|28E76FBC 12 Fung Go Roo 7 28.94 4FE7787E|15E59FBB|794EFDE8 13 John Doe 6 138.54 ]
	//
	//Interpolated Statement:
	//INSERT INTO `customer_entity` VALUES (11,'Karl
	//Gopher',7,47.11,'1FE9983E|28E76FBC'),(12,'Fung Go
	//Roo',7,28.94,'4FE7787E|15E59FBB|794EFDE8'),(13,'John Doe',6,138.54,'')
}

// ExampleColumnMapper_selectSalesOrdersFromSpecificCustomers this query should
// return all sales orders from different customers which are loaded within a
// collection. The challenge depict to map the customer_entity.entity_id column
// to the sales_order_entity.customer_id column.
func ExampleColumnMapper_selectSalesOrdersFromSpecificCustomers() {

	// Column `customer_id` has been hard coded into the switch statement of the
	// ColumnMapper in customerCollection and customerEntity. `customer_id` acts
	// as an alias to `entity_id`.
	q := dml.NewSelect("entity_id", "status", "increment_id", "grand_total", "tax_total").From("sales_order_entity").
		Where(dml.Column("customer_id").In().PlaceHolder()).BindRecord(
		dml.Qualify("", cmCustomers),
	)

	writeToSQLAndInterpolate(q)
	// Output:
	//Prepared Statement:
	//SELECT `entity_id`, `status`, `increment_id`, `grand_total`, `tax_total` FROM
	//`sales_order_entity` WHERE (`customer_id` IN (?))
	//Arguments: [11 12 13]
	//
	//Interpolated Statement:
	//SELECT `entity_id`, `status`, `increment_id`, `grand_total`, `tax_total` FROM
	//`sales_order_entity` WHERE (`customer_id` IN (11,12,13))
}
