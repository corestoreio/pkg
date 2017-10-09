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

import (
	"context"
	"fmt"
	"strings"

	"github.com/corestoreio/csfw/sql/dml"
	"github.com/corestoreio/csfw/util/cstesting"
	"github.com/corestoreio/errors"
)

// Make sure that type customerEntity implements interface.
var _ dml.ColumnMapper = (*customerCollection)(nil)
var _ dml.ColumnMapper = (*customerEntity)(nil)

// customerCollection a slice of customer entities.
type customerCollection struct {
	Data []*customerEntity
	// AfterScan gets called in mode ColumnMapScan after the new
	// customerEntity has been created and assigned with values from the query.
	AfterScan []func(*customerEntity)
}

// customerEntity has been generated from the SQL table customer_entities.
type customerEntity struct {
	EntityID      uint64 // Auto Increment, supports until MaxUint64
	Firstname     string
	StoreID       uint16
	LifetimeSales dml.NullFloat64
	// VoucherCodes contains list of refunded codes, stored as CSV. Or even
	// stored in another table or even encrypted and the function decrypts it on
	// load. Same as the M2 EAV models.
	VoucherCodes exampleStringSlice
}

type exampleStringSlice []string

func (sl exampleStringSlice) ToString() string {
	return strings.Join(sl, "|")
}
func (sl exampleStringSlice) FromString(s string) []string {
	return strings.Split(s, "|")
}

func newCustomerEntity() *customerEntity {
	return &customerEntity{}
}

// MapColumns implements interface ColumnMapper only partially.
func (p *customerEntity) MapColumns(cm *dml.ColumnMap) error {
	if cm.Mode() == dml.ColumnMapEntityReadAll {
		voucherCodes := p.VoucherCodes.ToString()
		return cm.Uint64(&p.EntityID).String(&p.Firstname).Uint16(&p.StoreID).NullFloat64(&p.LifetimeSales).String(&voucherCodes).Err()
	}
	for cm.Next() {
		switch c := cm.Column(); c {
		case "entity_id":
			cm.Uint64(&p.EntityID)
		case "firstname":
			cm.String(&p.Firstname)
		case "store_id":
			cm.Uint16(&p.StoreID)
		case "lifetime_sales":
			cm.NullFloat64(&p.LifetimeSales)
		case "voucher_codes":
			if cm.Mode() == dml.ColumnMapScan {
				var voucherCodes string
				cm.String(&voucherCodes)
				p.VoucherCodes = p.VoucherCodes.FromString(voucherCodes)
			} else {
				voucherCodes := p.VoucherCodes.ToString()
				cm.String(&voucherCodes)
			}
		default:
			return errors.NewNotFoundf("[dml_test] customerEntity Column %q not found", c)
		}
	}
	return cm.Err()
}

func (cc *customerCollection) MapColumns(cm *dml.ColumnMap) error {
	switch m := cm.Mode(); m {
	case dml.ColumnMapEntityReadAll, dml.ColumnMapEntityReadSet:
		for _, p := range cc.Data {
			if err := p.MapColumns(cm); err != nil {
				return errors.WithStack(err)
			}
		}
	case dml.ColumnMapScan:
		if cm.Count == 0 {
			cc.Data = cc.Data[:0]
		}
		p := newCustomerEntity()
		if err := p.MapColumns(cm); err != nil {
			return errors.WithStack(err)
		}
		for _, fn := range cc.AfterScan {
			fn(p)
		}
		cc.Data = append(cc.Data, p)
	case dml.ColumnMapCollectionReadSet:
		for cm.Next() {
			switch c := cm.Column(); c {
			case "entity_id":
				cm.Args = cm.Args.Uint64s(cc.EntityIDs()...)
			case "firstname":
				cm.Args = cm.Args.Strings(cc.Firstnames()...)
			default:
				return errors.NewNotFoundf("[dml_test] customerCollection Column %q not found", c)
			}
		}
	default:
		return errors.NewNotSupportedf("[dml] Unknown Mode: %q", string(m))
	}
	return cm.Err()
}

func (ps *customerCollection) EntityIDs(ret ...uint64) []uint64 {
	if ret == nil {
		ret = make([]uint64, 0, len(ps.Data))
	}
	for _, p := range ps.Data {
		ret = append(ret, p.EntityID)
	}
	return ret
}

func (ps *customerCollection) Firstnames(ret ...string) []string {
	if ret == nil {
		ret = make([]string, 0, len(ps.Data))
	}
	for _, p := range ps.Data {
		ret = append(ret, p.Firstname)
	}
	return ret // can be made unique
}

// ExampleColumnMapper implementation POC for interface ColumnMapper.
func ExampleColumnMapper() {
	// <ignore_this>
	dbc, dbMock := cstesting.MockDB(nil)
	defer cstesting.MockClose(nil, dbc, dbMock)

	r := cstesting.MustMockRows(cstesting.WithFile("testdata", "customer_entity_example.csv"))
	dbMock.ExpectQuery(cstesting.SQLMockQuoteMeta("SELECT * FROM `customer_entity`")).WillReturnRows(r)
	// </ignore_this>

	customers := new(customerCollection)
	// Example Query 10 retrieving rows from a database mock
	{
		s := dml.NewSelect("*").From("customer_entity").WithDB(dbc.DB)
		_, err := s.Load(context.TODO(), customers)
		if err != nil {
			panic(err)
		}
		fmt.Printf("Result of %q query:\n", s)
		fmt.Println("[entity_id firstname store_id lifetime_sales voucher_codes]")
		for _, c := range customers.Data {
			fmt.Printf("%v\n", *c)
		}
	}

	// Example SELECT Query 20
	{
		s := dml.NewSelect("entity_id", "firstname", "lifetime_sales").From("customer_entity").
			Where(
				dml.Column("entity_id").In().PlaceHolder(),
			).
			BindRecord(dml.Qualify("", customers)).Interpolate()
		writeToSQL("\nExample query: SELECT using data from an unqualified customer collection.", s)
	}

	// Example SELECT Query 30
	{
		s := dml.NewSelect("ce.entity_id", "ce.firstname", "cg.customer_group_code", "cg.tax_class_id").FromAlias("customer_entity", "ce").
			Join(dml.MakeIdentifier("customer_group").Alias("cg"),
				dml.Column("ce.group_id").Equal().Column("cg.customer_group_id"),
			).
			Where(
				dml.Column("ce.entity_id").In().PlaceHolder(),
			).
			BindRecord(dml.Qualify("ce", customers)).Interpolate()
		writeToSQL("\nExample query: SELECT using data from a qualified customer collection.", s)
	}

	// TODO somehow variable customers must know that it should return all entity_ids for column customer_id
	//s = dml.NewSelect("entity_id", "status", "increment_id", "grand_total", "tax_total").From("sales_order_entity").
	//	Where(dml.Column("customer_id").In().PlaceHolder()).BindRecord(
	//	dml.Qualify("customer_entity", customers),
	//)

	// Example UPDATE Query 40
	{
		u := dml.NewUpdate("customer_entity").AddColumns("firstname", "lifetime_sales", "voucher_codes").
			BindRecord(dml.Qualify("", customers.Data[0])).
			Where(dml.Column("entity_id").Equal().PlaceHolder()).Interpolate()
		writeToSQL("\nExample query: UPDATE using data from an unqualified customer entity.", u)
	}

	// Example INSERT Query 50
	{
		in := dml.NewInsert("customer_entity").AddColumns("firstname", "lifetime_sales", "store_id", "voucher_codes").
			BindRecord(customers.Data[0], customers.Data[1]).Interpolate()
		writeToSQL("\nExample query: INSERT using data from customer entities.", in)
	}
	//{
	//	in := dml.NewInsert("customer_entity").
	//		BindRecord(customers.Data[0], customers.Data[1]).Interpolate()
	//	writeToSQL("\nExample query: INSERT using data from customer entities.", in)
	//}

	// Output:
	//Result of "SELECT * FROM `customer_entity`" query:
	//[entity_id firstname store_id lifetime_sales voucher_codes]
	//{18446744073700551613 Karl Gopher 7 {{47.11 true}} [1FE9983E 28E76FBC]}
	//{18446744073700551614 Fung Go Roo 7 {{28.94 true}} [4FE7787E 15E59FBB 794EFDE8]}
	//{18446744073700551615 John Doe 6 {{138.54 true}} []}
	//
	//Example query: SELECT using data from an unqualified customer collection. Statement:
	//SELECT `entity_id`, `firstname`, `lifetime_sales` FROM `customer_entity` WHERE
	//(`entity_id` IN
	//(18446744073700551613,18446744073700551614,18446744073700551615))
	//
	//Example query: SELECT using data from a qualified customer collection. Statement:
	//SELECT `ce`.`entity_id`, `ce`.`firstname`, `cg`.`customer_group_code`,
	//	`cg`.`tax_class_id` FROM `customer_entity` AS `ce` INNER JOIN `customer_group`
	//AS `cg` ON (`ce`.`group_id` = `cg`.`customer_group_id`) WHERE (`ce`.`entity_id`
	//IN (18446744073700551613,18446744073700551614,18446744073700551615))
	//
	//Example query: UPDATE using data from an unqualified customer entity. Statement:
	//UPDATE `customer_entity` SET `firstname`='Karl Gopher', `lifetime_sales`=47.11,
	//	`voucher_codes`='1FE9983E|28E76FBC' WHERE (`entity_id` = 18446744073700551613)
	//
	//Example query: INSERT using data from customer entities. Statement:
	//INSERT INTO `customer_entity`
	//(`firstname`,`lifetime_sales`,`store_id`,`voucher_codes`) VALUES ('Karl
	//Gopher',47.11,7,'1FE9983E|28E76FBC'),('Fung Go
	//Roo',28.94,7,'4FE7787E|15E59FBB|794EFDE8')
}
