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
var _ dbr.ArgumentAssembler = (*salesOrder)(nil)

// salesOrder represents just a demo record.
type salesOrder struct {
	EntityID   int64  // Auto Increment
	State      string // processing, pending, shipped,
	StoreID    dbr.ArgInt64
	CustomerID int64
	GrandTotal dbr.NullFloat64
}

func (so salesOrder) AssembleArguments(stmtType int, args dbr.Arguments, columns []string) (dbr.Arguments, error) {
	for _, c := range columns {
		switch c {
		case "entity_id":
			args = append(args, dbr.ArgInt64(so.EntityID))
		case "state":
			args = append(args, dbr.ArgString(so.State))
		case "store_id":
			args = append(args, so.StoreID)
		case "customer_id":
			args = append(args, dbr.ArgInt64(so.CustomerID))
		case "grand_total":
			args = append(args, so.GrandTotal)
		default:
			return nil, errors.NewNotFoundf("[dbr_test] Column %q not found", c)
		}
	}
	if len(columns) == 0 && stmtType&(dbr.SQLPartValues) != 0 {
		args = append(args,
			dbr.ArgInt64(so.EntityID),
			dbr.ArgString(so.State),
			so.StoreID,
			dbr.ArgInt64(so.CustomerID),
			so.GrandTotal,
		)
	}
	return args, nil
}

func ExampleUpdateMulti() {
	// <ignore_this>
	dbc, dbMock := cstesting.MockDB(fatalLog{})
	defer func() {
		dbMock.ExpectClose()
		dbc.Close()
		if err := dbMock.ExpectationsWereMet(); err != nil {
			fmt.Printf("dbMock Error: %+v\n", err)
		}
	}()

	prep := dbMock.ExpectPrepare(cstesting.SQLMockQuoteMeta(
		"UPDATE `sales_order` SET `state`=?, `customer_id`=?, `grand_total`=? WHERE (`shipping_method` = ?) AND (`entity_id` = ?)",
	))

	prep.ExpectExec().WithArgs(
		"pending", int64(5678), 31.41459, "DHL", 1).
		WillReturnResult(sqlmock.NewResult(0, 1))

	prep.ExpectExec().WithArgs(
		"processing", int64(8912), nil, "DHL", 2).
		WillReturnResult(sqlmock.NewResult(0, 1))
	// </ignore_this>

	// Our objects which should update the columns in the database table
	// `sales_order`.
	so1 := salesOrder{1, "pending", 5, 5678, dbr.MakeNullFloat64(31.41459)}
	so2 := salesOrder{2, "processing", 7, 8912, dbr.NullFloat64{}}

	// Create the multi update statement
	um := dbr.NewUpdateMulti(
		dbr.NewUpdate("sales_order").
			AddColumns("state", "customer_id", "grand_total").
			Where(
				// dbr.Column("shipping_method", dbr.In.Str("DHL", "UPS")), // For all clauses the same restriction TODO fix bug when using IN
				dbr.Column("shipping_method", dbr.Equal.Str("DHL")), // For all clauses the same restriction
				dbr.Column("entity_id", dbr.Equal.Int64()),          // Int64() acts as a place holder
			).
			WithDB(dbc.DB), // Our template statement
	).AddRecords(so1, so2)

	results, err := um.Exec(context.Background())
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

type fatalLog struct{}

func (fatalLog) Fatalf(format string, args ...interface{}) {
	panic(fmt.Sprintf(format, args...))
}
