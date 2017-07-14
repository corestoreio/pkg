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
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/corestoreio/csfw/storage/dbr"
	"github.com/corestoreio/csfw/util/cstesting"
	"github.com/corestoreio/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// due to import cycle with the cstesting package, we must test externally

func TestUpdateMulti_Exec(t *testing.T) {
	t.Parallel()

	t.Run("no columns provided", func(t *testing.T) {
		mu := dbr.NewUpdateMulti(dbr.NewUpdate("catalog_product_entity").
			Where(dbr.Column("entity_id").In().PlaceHolder()), // ArgInt64 must be without arguments
		)
		res, err := mu.Exec(context.TODO())
		assert.Nil(t, res)
		assert.True(t, errors.IsEmpty(err), "%+v", err)
	})

	t.Run("alias mismatch", func(t *testing.T) {
		mu := dbr.NewUpdateMulti(dbr.NewUpdate("catalog_product_entity").AddColumns("sku", "updated_at").Where(dbr.Column("entity_id").In().PlaceHolder()))
		mu.ColumnAliases = []string{"update_sku"}
		res, err := mu.Exec(context.TODO())
		assert.Nil(t, res)
		assert.True(t, errors.IsMismatch(err), "%+v", err)
	})

	t.Run("empty Records and RecordChan", func(t *testing.T) {
		mu := dbr.NewUpdateMulti(
			dbr.NewUpdate("catalog_product_entity").AddColumns("sku", "updated_at").
				Where(dbr.Column("entity_id").In().PlaceHolder()),
		).WithDB(dbMock{
			error: errors.NewAlreadyClosedf("Who closed myself?"),
		})

		res, err := mu.Exec(context.TODO())
		assert.Nil(t, res)
		assert.True(t, errors.IsAlreadyClosed(err), "%+v", err)
	})

	records := []dbr.ArgumentAssembler{
		&dbrPerson{
			ID:    1,
			Name:  "Alf",
			Email: dbr.MakeNullString("alf@m') -- el.mac"),
		},
		&dbrPerson{
			ID:    2,
			Name:  "John",
			Email: dbr.MakeNullString("john@doe.com"),
		},
	}

	mu := dbr.NewUpdateMulti(
		dbr.NewUpdate("customer_entity").Alias("ce").AddColumns("name", "email").Where(dbr.Column("id").Equal().PlaceHolder()).Interpolate(),
	)

	// SM = SQL Mock
	setSQLMockInterpolate := func(m sqlmock.Sqlmock) {
		m.ExpectExec(cstesting.SQLMockQuoteMeta("UPDATE `customer_entity` AS `ce` SET `name`='Alf', `email`='alf@m\\') -- el.mac' WHERE (`id` = 1)")).
			WillReturnResult(sqlmock.NewResult(0, 1))
		m.ExpectExec(cstesting.SQLMockQuoteMeta("UPDATE `customer_entity` AS `ce` SET `name`='John', `email`='john@doe.com' WHERE (`id` = 2)")).
			WillReturnResult(sqlmock.NewResult(0, 1))
	}
	setSMPrepared := func(m sqlmock.Sqlmock) {
		prep := m.ExpectPrepare(cstesting.SQLMockQuoteMeta("UPDATE `customer_entity` AS `ce` SET `name`=?, `email`=? WHERE (`id` = ?)"))
		prep.ExpectExec().WithArgs("Alf", "alf@m') -- el.mac", 1).WillReturnResult(sqlmock.NewResult(0, 1))
		prep.ExpectExec().WithArgs("John", "john@doe.com", 2).WillReturnResult(sqlmock.NewResult(0, 1))
	}

	t.Run("preprocess no transaction", func(t *testing.T) {
		dbc, dbMock := cstesting.MockDB(t)
		defer cstesting.MockClose(t, dbc, dbMock)

		setSQLMockInterpolate(dbMock)

		mu.WithDB(dbc.DB)

		results, err := mu.Exec(context.TODO(), records...)
		if err != nil {
			t.Fatalf("%+v", err)
		}
		assert.Len(t, results, 2)
		for i, res := range results {
			aff, err := res.RowsAffected()
			if err != nil {
				t.Fatalf("Result index %d with error: %s", i, err)
			}
			assert.Exactly(t, int64(1), aff)
		}
	})

	t.Run("prepared no transaction", func(t *testing.T) {
		dbc, dbMock := cstesting.MockDB(t)
		defer cstesting.MockClose(t, dbc, dbMock)

		setSMPrepared(dbMock)

		mu.Update.IsInterpolate = false
		mu.WithDB(dbc.DB)

		results, err := mu.Exec(context.TODO(), records...)
		if err != nil {
			t.Fatalf("%+v", err)
		}
		assert.Len(t, results, 2)
		for i, res := range results {
			aff, err := res.RowsAffected()
			if err != nil {
				t.Fatalf("Result index %d with error: %s", i, err)
			}
			assert.Exactly(t, int64(1), aff)
		}
	})

	t.Run("prepared with transaction", func(t *testing.T) {
		dbc, dbMock := cstesting.MockDB(t)
		defer cstesting.MockClose(t, dbc, dbMock)

		dbMock.ExpectBegin()
		setSMPrepared(dbMock)
		dbMock.ExpectCommit()

		mu.Tx = dbc.DB
		mu.Transaction()
		mu.Update.IsInterpolate = false
		mu.WithDB(dbc.DB)

		results, err := mu.Exec(context.TODO(), records...)
		if err != nil {
			t.Fatalf("%+v", err)
		}
		assert.Len(t, results, 2)
		for i, res := range results {
			aff, err := res.RowsAffected()
			if err != nil {
				t.Fatalf("Result index %d with error: %s", i, err)
			}
			assert.Exactly(t, int64(1), aff)
		}
	})

	t.Run("preprocess with transaction", func(t *testing.T) {
		dbc, dbMock := cstesting.MockDB(t)
		defer cstesting.MockClose(t, dbc, dbMock)

		dbMock.ExpectBegin()
		setSQLMockInterpolate(dbMock)
		dbMock.ExpectCommit()

		mu.Tx = dbc.DB
		mu.Transaction()
		mu.Update.IsInterpolate = true
		mu.WithDB(dbc.DB)

		results, err := mu.Exec(context.TODO(), records...)
		if err != nil {
			t.Fatalf("%+v", err)
		}
		assert.Len(t, results, 2)
		for i, res := range results {
			aff, err := res.RowsAffected()
			if err != nil {
				t.Fatalf("Result index %d with error: %s", i, err)
			}
			assert.Exactly(t, int64(1), aff)
		}
	})
}

// Make sure that type salesInvoice implements interface.
var _ dbr.ArgumentAssembler = (*salesInvoice)(nil)

// salesInvoice represents just a demo record.
type salesInvoice struct {
	EntityID   int64  // Auto Increment
	State      string // processing, pending, shipped,
	StoreID    dbr.ArgInt64
	CustomerID int64
	GrandTotal dbr.NullFloat64
}

func (so salesInvoice) AssembleArguments(stmtType int, args dbr.Arguments, columns []string) (dbr.Arguments, error) {
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
		case "alias_customer_id":
			// Here can be a special treatment implement like encoding to JSON or encryption
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

func TestUpdateMulti_ColumnAliases(t *testing.T) {
	t.Parallel()

	dbc, dbMock := cstesting.MockDB(t)

	prep := dbMock.ExpectPrepare(cstesting.SQLMockQuoteMeta(
		"UPDATE `sales_invoice` SET `state`=?, `customer_id`=?, `grand_total`=? WHERE (`shipping_method` = ?) AND (`entity_id` = ?)",
	))

	prep.ExpectExec().WithArgs(
		"pending", int64(5678), 31.41459, "DHL", 1).
		WillReturnResult(sqlmock.NewResult(0, 1))

	prep.ExpectExec().WithArgs(
		"processing", int64(8912), nil, "DHL", 2).
		WillReturnResult(sqlmock.NewResult(0, 1))
	// </ignore_this>

	// Our objects which should update the columns in the database table
	// `sales_invoice`.
	so1 := salesInvoice{1, "pending", 5, 5678, dbr.MakeNullFloat64(31.41459)}
	so2 := salesInvoice{2, "processing", 7, 8912, dbr.NullFloat64{}}

	// Create the multi update statement
	um := dbr.NewUpdateMulti(dbr.NewUpdate("sales_invoice").
		AddColumns("state", "customer_id", "grand_total").
		Where(
			// dbr.Column("shipping_method", dbr.In.Str("DHL", "UPS")), // For all clauses the same restriction TODO fix bug when using IN
			dbr.Column("shipping_method").String("DHL"), // For all clauses the same restriction
			dbr.Column("entity_id").PlaceHolder(),       // Int64() acts as a place holder
		), // Our template statement
	).WithDB(dbc.DB)

	um.ColumnAliases = []string{"state", "alias_customer_id", "grand_total"}

	results, err := um.Exec(context.Background(), so1, so2)
	require.NoError(t, err)

	dbMock.ExpectClose()
	dbc.Close()
	require.NoError(t, dbMock.ExpectationsWereMet())

	for i, r := range results {
		ra, err := r.RowsAffected()
		require.NoError(t, err, "Index %d", i)
		assert.Exactly(t, int64(1), ra, "Index %d", i)
	}
}

func TestUpdate_SetRecord_Arguments(t *testing.T) {
	ce := &categoryEntity{
		EntityID:       678,
		AttributeSetID: 6,
		ParentID:       "p456",
		Path:           dbr.MakeNullString("3/4/5"),
	}

	t.Run("1 WHERE", func(t *testing.T) {
		u := dbr.NewUpdate("catalog_category_entity").
			AddColumns("attribute_set_id", "parent_id", "path").
			SetRecord(ce).
			Where(dbr.Column("entity_id").Greater().PlaceHolder()) // No Arguments in Int64 because we need a place holder.

		compareToSQL(t, u, nil,
			"UPDATE `catalog_category_entity` SET `attribute_set_id`=?, `parent_id`=?, `path`=? WHERE (`entity_id` > ?)",
			"UPDATE `catalog_category_entity` SET `attribute_set_id`=6, `parent_id`='p456', `path`='3/4/5' WHERE (`entity_id` > 678)",
			int64(6), "p456", "3/4/5", int64(678),
		)
	})

	t.Run("2 WHERE", func(t *testing.T) {
		u := dbr.NewUpdate("catalog_category_entity").
			AddColumns("attribute_set_id", "parent_id", "path").
			SetRecord(ce).
			Where(
				dbr.Column("x").In().Int64s(66, 77),
				dbr.Column("entity_id").Greater().PlaceHolder(),
			)
		compareToSQL(t, u, nil,
			"UPDATE `catalog_category_entity` SET `attribute_set_id`=?, `parent_id`=?, `path`=? WHERE (`x` IN (?,?)) AND (`entity_id` > ?)",
			"UPDATE `catalog_category_entity` SET `attribute_set_id`=6, `parent_id`='p456', `path`='3/4/5' WHERE (`x` IN (66,77)) AND (`entity_id` > 678)",
			int64(6), "p456", "3/4/5", int64(66), int64(77), int64(678),
		)
	})
	t.Run("3 WHERE", func(t *testing.T) {
		u := dbr.NewUpdate("catalog_category_entity").
			AddColumns("attribute_set_id", "parent_id", "path").
			SetRecord(ce).
			Where(
				dbr.Column("entity_id").Greater().PlaceHolder(),
				dbr.Column("x").In().Int64s(66, 77),
				dbr.Column("y").Greater().Int64(99),
			)
		compareToSQL(t, u, nil,
			"UPDATE `catalog_category_entity` SET `attribute_set_id`=?, `parent_id`=?, `path`=? WHERE (`entity_id` > ?) AND (`x` IN (?,?)) AND (`y` > ?)",
			"UPDATE `catalog_category_entity` SET `attribute_set_id`=6, `parent_id`='p456', `path`='3/4/5' WHERE (`entity_id` > 678) AND (`x` IN (66,77)) AND (`y` > 99)",
			int64(6), "p456", "3/4/5", int64(678), int64(66), int64(77), int64(99),
		)
	})

}
