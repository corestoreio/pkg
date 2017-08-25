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

func TestUpdate_Prepare(t *testing.T) {
	t.Parallel()

	t.Run("no columns provided", func(t *testing.T) {
		mu := dbr.NewUpdate("catalog_product_entity").Where(dbr.Column("entity_id").In().PlaceHolder())

		res, err := mu.Exec(context.TODO())
		assert.Nil(t, res)
		assert.True(t, errors.IsEmpty(err), "%+v", err)
	})

	t.Run("alias mismatch Exec", func(t *testing.T) {
		mu := dbr.NewUpdate("catalog_product_entity").
			AddColumns("sku", "updated_at").
			Where(dbr.Column("entity_id").In().PlaceHolder()).WithDB(dbMock{})
		mu.SetClausAliases = []string{"update_sku"}
		res, err := mu.Exec(context.TODO())
		assert.Nil(t, res)
		assert.True(t, errors.IsMismatch(err), "%+v", err)
	})

	t.Run("alias mismatch Prepare", func(t *testing.T) {
		mu := dbr.NewUpdate("catalog_product_entity").
			AddColumns("sku", "updated_at").
			Where(dbr.Column("entity_id").In().PlaceHolder()).WithDB(dbMock{})
		mu.SetClausAliases = []string{"update_sku"}
		stmt, err := mu.Prepare(context.TODO())
		assert.Nil(t, stmt)
		assert.True(t, errors.IsMismatch(err), "%+v", err)
	})

	t.Run("empty Records and RecordChan", func(t *testing.T) {
		mu := dbr.NewUpdate("catalog_product_entity").AddColumns("sku", "updated_at").
			Where(dbr.Column("entity_id").In().PlaceHolder()).WithDB(dbMock{
			error: errors.NewAlreadyClosedf("Who closed myself?"),
		})

		res, err := mu.Exec(context.TODO())
		assert.Nil(t, res)
		assert.True(t, errors.IsAlreadyClosed(err), "%+v", err)
	})

	t.Run("prepared success", func(t *testing.T) {
		dbc, dbMock := cstesting.MockDB(t)
		defer cstesting.MockClose(t, dbc, dbMock)

		records := []*dbrPerson{
			{
				ID:    1,
				Name:  "Alf",
				Email: dbr.MakeNullString("alf@m') -- el.mac"),
			},
			{
				ID:    2,
				Name:  "John",
				Email: dbr.MakeNullString("john@doe.com"),
			},
		}

		// interpolate must get ignored
		mu := dbr.NewUpdate("customer_entity").Alias("ce").
			AddColumns("name", "email").
			Where(dbr.Column("id").Equal().PlaceHolder()).
			WithDB(dbc.DB).
			Interpolate() // gets ignored

		prep := dbMock.ExpectPrepare(cstesting.SQLMockQuoteMeta("UPDATE `customer_entity` AS `ce` SET `name`=?, `email`=? WHERE (`id` = ?)"))
		prep.ExpectExec().WithArgs("Alf", "alf@m') -- el.mac", 1).WillReturnResult(sqlmock.NewResult(0, 1))
		prep.ExpectExec().WithArgs("John", "john@doe.com", 2).WillReturnResult(sqlmock.NewResult(0, 1))

		stmt, err := mu.Prepare(context.TODO())
		require.NoError(t, err)
		for i, record := range records {
			results, err := stmt.WithRecords(dbr.Qualify("ce", record)).Exec(context.TODO())
			require.NoError(t, err)
			aff, err := results.RowsAffected()
			if err != nil {
				t.Fatalf("Result index %d with error: %s", i, err)
			}
			assert.Exactly(t, int64(1), aff)
		}
	})

	t.Run("ExecContext", func(t *testing.T) {
		dbc, dbMock := cstesting.MockDB(t)
		defer cstesting.MockClose(t, dbc, dbMock)

		prep := dbMock.ExpectPrepare(cstesting.SQLMockQuoteMeta("UPDATE `customer_entity` AS `ce` SET `name`=?, `email`=? WHERE (`id` = ?)"))
		prep.ExpectExec().WithArgs("Peter Gopher", "peter@gopher.go", 3456).WillReturnResult(sqlmock.NewResult(0, 9))

		stmt, err := dbr.NewUpdate("customer_entity").Alias("ce").
			AddColumns("name", "email").
			Where(dbr.Column("id").Equal().PlaceHolder()).
			WithDB(dbc.DB).Prepare(context.TODO())
		require.NoError(t, err, "failed creating a prepared statement")
		defer func() {
			require.NoError(t, stmt.Close(), "Close on a prepared statement")
		}()

		res, err := stmt.Exec(context.TODO(), "Peter Gopher", "peter@gopher.go", 3456)
		require.NoError(t, err, "failed to execute ExecContext")

		ra, err := res.RowsAffected()
		if err != nil {
			t.Fatal(err)
		}
		assert.Exactly(t, int64(9), ra, "Different LastInsertIDs")
	})

	t.Run("ExecArgs One Row", func(t *testing.T) {
		dbc, dbMock := cstesting.MockDB(t)
		defer cstesting.MockClose(t, dbc, dbMock)

		prep := dbMock.ExpectPrepare(cstesting.SQLMockQuoteMeta("UPDATE `customer_entity` AS `ce` SET `name`=?, `email`=? WHERE (`id` = ?)"))
		prep.ExpectExec().WithArgs("Peter Gopher", "peter@gopher.go", 3456).WillReturnResult(sqlmock.NewResult(0, 11))
		prep.ExpectExec().WithArgs("Petra Gopher", "petra@gopher.go", 3457).WillReturnResult(sqlmock.NewResult(0, 21))

		stmt, err := dbr.NewUpdate("customer_entity").Alias("ce").
			AddColumns("name", "email").
			Where(dbr.Column("id").Equal().PlaceHolder()).
			WithDB(dbc.DB).Prepare(context.TODO())
		require.NoError(t, err, "failed creating a prepared statement")
		defer func() {
			require.NoError(t, stmt.Close(), "Close on a prepared statement")
		}()

		tests := []struct {
			name         string
			email        string
			id           int
			rowsAffected int64
		}{
			{"Peter Gopher", "peter@gopher.go", 3456, 11},
			{"Petra Gopher", "petra@gopher.go", 3457, 21},
		}

		args := dbr.MakeArgs(3)
		for i, test := range tests {
			args = args[:0]

			res, err := stmt.WithArguments(args.Str(test.name).Str(test.email).Int(test.id)).Exec(context.TODO())
			if err != nil {
				t.Fatalf("Index %d => %+v", i, err)
			}
			ra, err := res.RowsAffected()
			if err != nil {
				t.Fatalf("Result index %d with error: %s", i, err)
			}
			assert.Exactly(t, test.rowsAffected, ra, "Index %d has different LastInsertIDs", i)
		}
	})

}

// Make sure that type salesInvoice implements interface.
var _ dbr.ArgumentsAppender = (*salesInvoice)(nil)

// salesInvoice represents just a demo record.
type salesInvoice struct {
	EntityID   int64  // Auto Increment
	State      string // processing, pending, shipped,
	StoreID    int64
	CustomerID int64
	GrandTotal dbr.NullFloat64
}

func (so salesInvoice) AppendArgs(args dbr.Arguments, columns []string) (dbr.Arguments, error) {
	column := columns[0]
	switch column {
	case "entity_id":
		args = args.Int64(so.EntityID)
	case "state":
		args = args.Str(so.State)
	case "store_id":
		args = args.Int64(so.StoreID)
	case "customer_id":
		args = args.Int64(so.CustomerID)
	case "alias_customer_id":
		// Here can be a special treatment implemented like encoding to JSON
		// or encryption
		args = args.Int64(so.CustomerID)
	case "grand_total":
		args = args.NullFloat64(so.GrandTotal)
	default:
		return nil, errors.NewNotFoundf("[dbr_test] Column %q not found", column)
	}
	return args, nil
}

func TestUpdate_SetClausAliases(t *testing.T) {
	t.Parallel()

	dbc, dbMock := cstesting.MockDB(t)

	prep := dbMock.ExpectPrepare(cstesting.SQLMockQuoteMeta(
		"UPDATE `sales_invoice` SET `state`=?, `customer_id`=?, `grand_total`=? WHERE (`shipping_method` = ?) AND (`entity_id` = ?)",
	))

	prep.ExpectExec().WithArgs(
		"pending", int64(5678), 31.41459, "DHL", 21).
		WillReturnResult(sqlmock.NewResult(0, 1))

	prep.ExpectExec().WithArgs(
		"processing", int64(8912), nil, "DHL", 32).
		WillReturnResult(sqlmock.NewResult(0, 1))
	// </ignore_this>

	// Our objects which should update the columns in the database table
	// `sales_invoice`.

	collection := []salesInvoice{
		{21, "pending", 5, 5678, dbr.MakeNullFloat64(31.41459)},
		{32, "processing", 7, 8912, dbr.NullFloat64{}},
	}

	// Create the multi update statement
	um := dbr.NewUpdate("sales_invoice").
		AddColumns("state", "customer_id", "grand_total").
		Where(
			// dbr.Column("shipping_method", dbr.In.Str("DHL", "UPS")), // For all clauses the same restriction TODO fix bug when using IN
			dbr.Column("shipping_method").Str("DHL"), // For all clauses the same restriction
			dbr.Column("entity_id").PlaceHolder(),    // Int64() acts as a place holder
		).WithDB(dbc.DB)

	um.SetClausAliases = []string{"state", "alias_customer_id", "grand_total"}

	stmt, err := um.Prepare(context.TODO())
	require.NoError(t, err)

	for i, record := range collection {
		results, err := stmt.WithRecords(dbr.Qualify("sales_invoice", record)).Exec(context.TODO())
		require.NoError(t, err)
		ra, err := results.RowsAffected()
		require.NoError(t, err, "Index %d", i)
		assert.Exactly(t, int64(1), ra, "Index %d", i)
	}

	dbMock.ExpectClose()
	dbc.Close()
	require.NoError(t, dbMock.ExpectationsWereMet())

}

func TestUpdate_BindRecord(t *testing.T) {
	ce := &categoryEntity{
		EntityID:       678,
		AttributeSetID: 6,
		ParentID:       "p456",
		Path:           dbr.MakeNullString("3/4/5"),
	}

	t.Run("1 WHERE", func(t *testing.T) {
		u := dbr.NewUpdate("catalog_category_entity").
			AddColumns("attribute_set_id", "parent_id", "path").
			BindRecord(dbr.Qualify("", ce)).
			Where(dbr.Column("entity_id").Greater().PlaceHolder())

		compareToSQL(t, u, nil,
			"UPDATE `catalog_category_entity` SET `attribute_set_id`=?, `parent_id`=?, `path`=? WHERE (`entity_id` > ?)",
			"UPDATE `catalog_category_entity` SET `attribute_set_id`=6, `parent_id`='p456', `path`='3/4/5' WHERE (`entity_id` > 678)",
			int64(6), "p456", "3/4/5", int64(678),
		)
	})

	t.Run("2 WHERE", func(t *testing.T) {
		u := dbr.NewUpdate("catalog_category_entity").
			AddColumns("attribute_set_id", "parent_id", "path").
			BindRecord(dbr.Qualify("", ce)).
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
			BindRecord(dbr.Qualify("", ce)).
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

	t.Run("with alias table name", func(t *testing.T) {
		// A fictional table statement which already reflects future JOIN
		// implementation.
		u := dbr.NewUpdate("catalog_category_entity").Alias("ce").
			AddColumns("attribute_set_id", "parent_id", "path").
			BindRecord(dbr.Qualify("", ce), dbr.Qualify("cpei", ce)).
			Where(
				dbr.Column("ce.entity_id").Greater().PlaceHolder(), //678
				dbr.Column("cpe.entity_id").In().Int64s(66, 77),
				dbr.Column("cpei.attribute_set_id").In().PlaceHolder(), //6
			)
		compareToSQL(t, u, nil,
			"UPDATE `catalog_category_entity` AS `ce` SET `attribute_set_id`=?, `parent_id`=?, `path`=? WHERE (`ce`.`entity_id` > ?) AND (`cpe`.`entity_id` IN (?,?)) AND (`cpei`.`attribute_set_id` IN (?))",
			"UPDATE `catalog_category_entity` AS `ce` SET `attribute_set_id`=6, `parent_id`='p456', `path`='3/4/5' WHERE (`ce`.`entity_id` > 678) AND (`cpe`.`entity_id` IN (66,77)) AND (`cpei`.`attribute_set_id` IN (6))",
			int64(6), "p456", "3/4/5", int64(678), int64(66), int64(77), int64(6),
		)
	})
}
