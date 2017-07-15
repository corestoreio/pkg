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
	"context"
	"database/sql"
	"fmt"

	sqlmock "github.com/DATA-DOG/go-sqlmock"
	"github.com/corestoreio/csfw/storage/dbr"
	"github.com/corestoreio/csfw/util/cstesting"
	"github.com/corestoreio/errors"
)

// Make sure that type salesCreditMemo implements interface.
var _ dbr.Scanner = (*salesCreditMemoCollection)(nil)

// salesCreditMemo represents just a demo record.
type salesCreditMemoCollection struct {
	Convert        dbr.RowConvert
	Data           []*salesCreditMemo
	EventAfterScan []func(*salesCreditMemo)
}

// salesCreditMemo represents just a demo record.
type salesCreditMemo struct {
	EntityID   uint64 // Auto Increment, supports until MaxUint64
	State      string // processing, pending, shipped,
	StoreID    uint16
	CustomerID int64
	GrandTotal sql.NullFloat64
}

func (cc *salesCreditMemoCollection) RowScan(r *sql.Rows) error {
	if err := cc.Convert.Scan(r); err != nil {
		return err
	}

	o := new(salesCreditMemo)
	for i, col := range cc.Convert.Columns {
		if cc.Convert.Alias != nil {
			if orgCol, ok := cc.Convert.Alias[col]; ok {
				col = orgCol
			}
		}
		b := cc.Convert.Index(i)
		var err error
		switch col {
		case "entity_id":
			o.EntityID, err = b.Uint64()
		case "state":
			o.State, err = b.Str()
		case "store_id":
			o.StoreID, err = b.Uint16()
		case "customer_id":
			o.CustomerID, err = b.Int64()
		case "grand_total":
			o.GrandTotal, err = b.NullFloat64()
		}
		if err != nil {
			return errors.Wrapf(err, "[dbr] Failed to convert value at row % with column index %d", cc.Convert.Count, i)
		}
	}
	// For example to implement an event after scanning has been performed.
	for _, fn := range cc.EventAfterScan {
		fn(o)
	}

	cc.Data = append(cc.Data, o)
	return nil
}

func ExampleRowConvert() {
	// <ignore_this>
	dbc, dbMock := cstesting.MockDB(nil)
	defer cstesting.MockClose(nil, dbc, dbMock)

	r := sqlmock.NewRows([]string{"entity_id", "state", "store_id", "customer_id", "grand_total"}).
		FromCSVString("18446744073700551613,shipped,7,98765,47.11\n18446744073700551614,shipped,7,12345,28.94\n")

	dbMock.ExpectQuery(cstesting.SQLMockQuoteMeta("SELECT * FROM `sales_creditmemo` WHERE (`state` = 'shipped')")).WillReturnRows(r)
	// </ignore_this>

	s := dbr.NewSelect("*").From("sales_creditmemo").
		Where(dbr.Column("state").String("shipped")).
		WithDB(dbc.DB).Interpolate()

	cmc := &salesCreditMemoCollection{}
	_, err := s.Load(context.TODO(), cmc)
	if err != nil {
		panic(err)
	}

	fmt.Printf("%v\n", cmc.Convert.Columns)
	for _, c := range cmc.Data {
		fmt.Printf("%v\n", *c)
	}

	// Output:
	//[entity_id state store_id customer_id grand_total]
	//{18446744073700551613 shipped 7 98765 {47.11 true}}
	//{18446744073700551614 shipped 7 12345 {28.94 true}}
}
