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

	sqlmock "github.com/DATA-DOG/go-sqlmock"
	"github.com/corestoreio/csfw/sql/dml"
	"github.com/corestoreio/csfw/util/cstesting"
	"github.com/corestoreio/errors"
)

// Make sure that type salesCreditMemo implements interface.
var _ dml.ColumnMapper = (*salesCreditMemoCollection)(nil)
var _ dml.ColumnMapper = (*salesCreditMemo)(nil)

// salesCreditMemo represents just a demo record.
type salesCreditMemoCollection struct {
	Data           []*salesCreditMemo
	EventAfterScan []func(*salesCreditMemo)
}

// salesCreditMemo represents just a demo record.
type salesCreditMemo struct {
	EntityID   uint64 // Auto Increment, supports until MaxUint64
	State      string // processing, pending, shipped,
	StoreID    uint16
	CustomerID int64
	GrandTotal dml.NullFloat64
	// VoucherCodes contains list of refunded codes, stored as CSV. Or even
	// stored in another table or even encrypted and the function decrypts it on
	// load. Same as the M2 EAV models
	VoucherCodes       []string
	VoucherCodeEncoder interface {
		Encode([]string) (string, error)
		Decode(string) ([]string, error)
	}
}

type globalExampleStringSliceEncoder struct{}

func (globalExampleStringSliceEncoder) Encode(sl []string) (string, error) {
	return strings.Join(sl, "|"), nil
}
func (globalExampleStringSliceEncoder) Decode(s string) ([]string, error) {
	return strings.Split(s, "|"), nil
}

func newSalesCreditMemo() *salesCreditMemo {
	return &salesCreditMemo{
		VoucherCodeEncoder: globalExampleStringSliceEncoder{},
	}
}

// MapColumns implements interface ColumnMapper only partially.
func (p *salesCreditMemo) MapColumns(cm *dml.ColumnMap) error {
	if cm.Mode() == dml.ColumnMapEntityReadAll {
		voucherCodes, _ := p.VoucherCodeEncoder.Encode(p.VoucherCodes)
		return cm.Uint64(&p.EntityID).String(&p.State).Uint16(&p.StoreID).Int64(&p.CustomerID).NullFloat64(&p.GrandTotal).String(&voucherCodes).Err()
	}
	for cm.Next() {
		switch c := cm.Column(); c {
		case "entity_id":
			cm.Uint64(&p.EntityID)
		case "state":
			cm.String(&p.State)
		case "store_id":
			cm.Uint16(&p.StoreID)
		case "customer_id":
			cm.Int64(&p.CustomerID)
		case "grand_total":
			cm.NullFloat64(&p.GrandTotal)
		case "voucher_codes":
			// TODO(CyS) fix bug in case we're reading or writing
			var voucherCodes string
			cm.String(&voucherCodes)
			p.VoucherCodes, _ = p.VoucherCodeEncoder.Decode(voucherCodes)
		default:
			return errors.NewNotFoundf("[dml_test] dmlPerson Column %q not found", c)
		}
	}
	return cm.Err()
}

func (cc *salesCreditMemoCollection) MapColumns(cm *dml.ColumnMap) error {
	switch m := cm.Mode(); m {
	case 'a', 'R':
		// INSERT STATEMENT requesting all columns aka arguments or SELECT
		// requesting specific columns.
		for _, p := range cc.Data {
			if err := p.MapColumns(cm); err != nil {
				return errors.WithStack(err)
			}
		}
	case dml.ColumnMapCollectionCreate:
		// case for scanning when loading certain rows, hence we write data from
		// the DB into the struct in each for-loop.
		if cm.Count == 0 {
			cc.Data = cc.Data[:0]
		}
		p := newSalesCreditMemo()
		if err := p.MapColumns(cm); err != nil {
			return errors.WithStack(err)
		}
		for _, fn := range cc.EventAfterScan {
			fn(p)
		}
		cc.Data = append(cc.Data, p)
	case 'r':
		// SELECT, DELETE or UPDATE or INSERT with n columns
		// omitted because not needed in this example.
	default:
		return errors.NewNotSupportedf("[dml] Unknown Mode: %q", string(m))
	}
	return cm.Err()
}

// ExampleColumnMapper implementation POC for interface ColumnMapper.
func ExampleColumnMapper() {
	// <ignore_this>
	dbc, dbMock := cstesting.MockDB(nil)
	defer cstesting.MockClose(nil, dbc, dbMock)

	r := sqlmock.NewRows([]string{"entity_id", "state", "store_id", "customer_id", "grand_total", "voucher_codes"}).
		FromCSVString("18446744073700551613,shipped,7,98765,47.11,1FE9983E|28E76FBC\n18446744073700551614,shipped,7,12345,28.94,4FE7787E|15E59FBB|794EFDE8\n")

	dbMock.ExpectQuery(cstesting.SQLMockQuoteMeta("SELECT * FROM `sales_creditmemo` WHERE (`state` = 'shipped')")).WillReturnRows(r)
	// </ignore_this>

	s := dml.NewSelect("*").From("sales_creditmemo").
		Where(dml.Column("state").Str("shipped")).
		WithDB(dbc.DB).Interpolate()

	cmc := &salesCreditMemoCollection{}
	_, err := s.Load(context.TODO(), cmc)
	if err != nil {
		panic(err)
	}

	fmt.Print("[entity_id state store_id customer_id grand_total voucher_codes]\n")
	for _, c := range cmc.Data {
		fmt.Printf("%v\n", *c)
	}

	// Output:
	//[entity_id state store_id customer_id grand_total voucher_codes]
	//{18446744073700551613 shipped 7 98765 {{47.11 true}} [1FE9983E 28E76FBC] {}}
	//{18446744073700551614 shipped 7 12345 {{28.94 true}} [4FE7787E 15E59FBB 794EFDE8] {}}
}
