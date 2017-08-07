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
	"fmt"

	sqlmock "github.com/DATA-DOG/go-sqlmock"
	"github.com/corestoreio/csfw/storage/dbr"
	"github.com/corestoreio/csfw/util/cstesting"
	"github.com/corestoreio/errors"
)

// Make sure that type salesOrder implements interface.
var _ dbr.ArgumentsAppender = (*salesOrder)(nil)

// salesOrder represents just a demo record.
type salesOrder struct {
	EntityID   int64  // Auto Increment
	State      string // processing, pending, shipped,
	StoreID    int64
	CustomerID int64
	GrandTotal dbr.NullFloat64
}

func (so salesOrder) AppendArguments(stmtType int, args dbr.ArgUnions, columns []string) (dbr.ArgUnions, error) {
	for _, c := range columns {
		switch c {
		case "entity_id":
			args = args.Int64(so.EntityID)
		case "state":
			args = args.Str(so.State)
		case "store_id":
			args = args.Int64(so.StoreID)
		case "customer_id":
			args = args.Int64(so.CustomerID)
		case "grand_total":
			args = args.NullFloat64(so.GrandTotal)
		default:
			return nil, errors.NewNotFoundf("[dbr_test] Column %q not found", c)
		}
	}
	if len(columns) == 0 && stmtType&(dbr.SQLPartValues) != 0 {
		args = args.Int64(so.EntityID).Str(so.State).Int64(so.StoreID).Int64(so.CustomerID).NullFloat64(so.GrandTotal)
	}
	return args, nil
}

// ExampleUpdateMulti can create a prepared statement or interpolated statements
// to run updates on table  `sales_order` with different objects. The SQL UPDATE
// statement acts as a template.
func ExampleUpdateMulti() {
	// <ignore_this>
	dbc, dbMock := cstesting.MockDB(nil)
	defer cstesting.MockClose(nil, dbc, dbMock)

	prep := dbMock.ExpectPrepare(cstesting.SQLMockQuoteMeta(
		"UPDATE `sales_order` SET `state`=?, `customer_id`=?, `grand_total`=? WHERE (`shipping_method` IN (?,?)) AND (`entity_id` = ?)",
	))

	prep.ExpectExec().WithArgs(
		"pending", int64(5678), 31.41459, "DHL", "UPS", 1).
		WillReturnResult(sqlmock.NewResult(0, 1))

	prep.ExpectExec().WithArgs(
		"processing", int64(8912), nil, "DHL", "UPS", 2).
		WillReturnResult(sqlmock.NewResult(0, 1))
	// </ignore_this>

	// Create the multi update statement
	um := dbr.NewUpdateMulti(
		dbr.NewUpdate("sales_order").
			AddColumns("state", "customer_id", "grand_total").
			Where(
				dbr.Column("shipping_method").In().Strs("DHL", "UPS"),
				dbr.Column("entity_id").PlaceHolder(),
			), // Our template statement
	).WithDB(dbc.DB)

	// Our objects which should update the columns in the database table
	// `sales_order`.
	so1 := salesOrder{1, "pending", 5, 5678, dbr.MakeNullFloat64(31.41459)}
	so2 := salesOrder{2, "processing", 7, 8912, dbr.NullFloat64{}}

	results, err := um.Exec(context.Background(), so1, so2)
	if err != nil {
		fmt.Printf("Exec Error: %+v\n", err)
		return
	}
	for i, r := range results {
		ra, err := r.RowsAffected()
		if err != nil {
			fmt.Printf("Index %d RowsAffected Error: %+v\n", i, err)
			return
		}
		fmt.Printf("Index %d RowsAffected %d\n", i, ra)
	}

	// Output:
	//Index 0 RowsAffected 1
	//Index 1 RowsAffected 1
}
