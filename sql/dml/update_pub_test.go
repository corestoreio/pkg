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
	"context"
	"database/sql"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/corestoreio/errors"
	"github.com/corestoreio/log"
	"github.com/corestoreio/pkg/sql/dml"
	"github.com/corestoreio/pkg/sql/dmltest"
	"github.com/corestoreio/pkg/storage/null"
	"github.com/corestoreio/pkg/util/assert"
)

func TestUpdate_WithArgs(t *testing.T) {
	t.Parallel()

	t.Run("no columns provided", func(t *testing.T) {
		mu := dml.NewUpdate("catalog_product_entity").Where(dml.Column("entity_id").In().PlaceHolder())

		res, err := mu.WithArgs().ExecContext(context.TODO())
		assert.Nil(t, res)
		assert.True(t, errors.Empty.Match(err), "%+v", err)
	})

	t.Run("alias mismatch Exec", func(t *testing.T) {

		res, err := dml.NewUpdate("catalog_product_entity").
			AddColumns("sku", "updated_at").
			Where(dml.Column("entity_id").In().PlaceHolder()).WithDB(dbMock{
			prepareFn: func(query string) (*sql.Stmt, error) {
				return nil, nil
			},
		}).
			WithArgs().WithQualifiedColumnsAliases("update_sku").
			Record("", nil).ExecContext(context.TODO())
		assert.Nil(t, res)
		assert.True(t, errors.Mismatch.Match(err), "Should be of kind errors.Mismatch %+v", err)
	})

	t.Run("alias mismatch Prepare", func(t *testing.T) {
		dbc, dbMock := dmltest.MockDB(t)
		defer dmltest.MockClose(t, dbc, dbMock)

		_ = dbMock.ExpectPrepare(dmltest.SQLMockQuoteMeta("UPDATE `catalog_product_entity` SET `sku`=?, `updated_at`=? WHERE (`entity_id` IN ?)"))

		mu := dbc.Update("catalog_product_entity").
			AddColumns("sku", "updated_at").
			Where(dml.Column("entity_id").In().PlaceHolder())

		stmt, err := mu.Prepare(context.TODO())
		assert.NoError(t, err)

		p := &dmlPerson{
			ID:    1,
			Name:  "Alf",
			Email: null.MakeString("alf@m') -- el.mac"),
		}
		res, err := stmt.WithArgs().Record("", p).WithQualifiedColumnsAliases("update_sku").ExecContext(context.TODO())
		assert.True(t, errors.Mismatch.Match(err), "%+v", err)
		assert.Nil(t, res, "No result is expected")
	})

	t.Run("empty Records and RecordChan", func(t *testing.T) {
		mu := dml.NewUpdate("catalog_product_entity").AddColumns("sku", "updated_at").
			Where(dml.Column("entity_id").In().PlaceHolder()).WithDB(dbMock{
			error: errors.AlreadyClosed.Newf("Who closed myself?"),
		})

		res, err := mu.WithArgs().ExecContext(context.TODO())
		assert.Nil(t, res)
		assert.True(t, errors.AlreadyClosed.Match(err), "%+v", err)
	})

	t.Run("prepared success", func(t *testing.T) {
		dbc, dbMock := dmltest.MockDB(t)
		defer dmltest.MockClose(t, dbc, dbMock)

		records := []*dmlPerson{
			{
				ID:    1,
				Name:  "Alf",
				Email: null.MakeString("alf@m') -- el.mac"),
			},
			{
				ID:    2,
				Name:  "John",
				Email: null.MakeString("john@doe.com"),
			},
		}

		// interpolate must get ignored
		mu := dml.NewUpdate("customer_entity").Alias("ce").
			AddColumns("name", "email").
			Where(dml.Column("id").Equal().PlaceHolder()).
			WithDB(dbc.DB)

		prep := dbMock.ExpectPrepare(dmltest.SQLMockQuoteMeta("UPDATE `customer_entity` AS `ce` SET `name`=?, `email`=? WHERE (`id` = ?)"))
		prep.ExpectExec().WithArgs("Alf", "alf@m') -- el.mac", 1).WillReturnResult(sqlmock.NewResult(0, 1))
		prep.ExpectExec().WithArgs("John", "john@doe.com", 2).WillReturnResult(sqlmock.NewResult(0, 1))

		stmt, err := mu.Prepare(context.TODO())
		assert.NoError(t, err)
		for i, record := range records {
			results, err := stmt.WithArgs().Record("ce", record).ExecContext(context.TODO())
			assert.NoError(t, err)
			aff, err := results.RowsAffected()
			if err != nil {
				t.Fatalf("Result index %d with error: %s", i, err)
			}
			assert.Exactly(t, int64(1), aff)
		}
	})

	t.Run("ExecContext", func(t *testing.T) {
		dbc, dbMock := dmltest.MockDB(t)
		defer dmltest.MockClose(t, dbc, dbMock)

		prep := dbMock.ExpectPrepare(dmltest.SQLMockQuoteMeta("UPDATE `customer_entity` AS `ce` SET `name`=?, `email`=? WHERE (`id` = ?)"))
		prep.ExpectExec().WithArgs("Peter Gopher", "peter@gopher.go", 3456).WillReturnResult(sqlmock.NewResult(0, 9))

		stmt, err := dml.NewUpdate("customer_entity").Alias("ce").
			AddColumns("name", "email").
			Where(dml.Column("id").Equal().PlaceHolder()).
			WithDB(dbc.DB).Prepare(context.TODO())
		assert.NoError(t, err, "failed creating a prepared statement")
		defer func() {
			assert.NoError(t, stmt.Close(), "Close on a prepared statement")
		}()

		res, err := stmt.WithArgs().ExecContext(context.TODO(), "Peter Gopher", "peter@gopher.go", 3456)
		assert.NoError(t, err, "failed to execute ExecContext")

		ra, err := res.RowsAffected()
		if err != nil {
			t.Fatal(err)
		}
		assert.Exactly(t, int64(9), ra, "Different LastInsertIDs")
	})

	t.Run("ExecArgs One Row", func(t *testing.T) {
		dbc, dbMock := dmltest.MockDB(t)
		defer dmltest.MockClose(t, dbc, dbMock)

		prep := dbMock.ExpectPrepare(dmltest.SQLMockQuoteMeta("UPDATE `customer_entity` AS `ce` SET `name`=?, `email`=? WHERE (`id` = ?)"))
		prep.ExpectExec().WithArgs("Peter Gopher", "peter@gopher.go", 3456).WillReturnResult(sqlmock.NewResult(0, 11))
		prep.ExpectExec().WithArgs("Petra Gopher", "petra@gopher.go", 3457).WillReturnResult(sqlmock.NewResult(0, 21))

		stmt, err := dml.NewUpdate("customer_entity").Alias("ce").
			AddColumns("name", "email").
			Where(dml.Column("id").Equal().PlaceHolder()).
			WithDB(dbc.DB).Prepare(context.TODO())
		assert.NoError(t, err, "failed creating a prepared statement")
		defer func() {
			assert.NoError(t, stmt.Close(), "Close on a prepared statement")
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

		for i, test := range tests {
			res, err := stmt.WithArgs().String(test.name).String(test.email).Int(test.id).ExecContext(context.TODO())
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
var _ dml.ColumnMapper = (*salesInvoice)(nil)

// salesInvoice represents just a demo record.
type salesInvoice struct {
	EntityID   int64  // Auto Increment
	State      string // processing, pending, shipped,
	StoreID    int64
	CustomerID int64
	GrandTotal null.Float64
}

func (so *salesInvoice) MapColumns(cm *dml.ColumnMap) error {
	for cm.Next() {
		switch c := cm.Column(); c {
		case "entity_id":
			cm.Int64(&so.EntityID)
		case "state":
			cm.String(&so.State)
		case "store_id":
			cm.Int64(&so.StoreID)
		case "customer_id":
			cm.Int64(&so.CustomerID)
		case "alias_customer_id":
			// Here can be a special treatment implemented like encoding to JSON
			// or encryption
			cm.Int64(&so.CustomerID)
		case "grand_total":
			cm.NullFloat64(&so.GrandTotal)
		default:
			return errors.NotFound.Newf("[dml_test] Column %q not found", c)
		}
	}
	return cm.Err()
}

func TestArguments_WithQualifiedColumnsAliases(t *testing.T) {
	t.Parallel()

	dbc, dbMock := dmltest.MockDB(t)

	prep := dbMock.ExpectPrepare(dmltest.SQLMockQuoteMeta(
		"UPDATE `sales_invoice` SET `state`=?, `customer_id`=?, `grand_total`=? WHERE (`shipping_method` IN ('DHL','UPS')) AND (`entity_id` = ?)",
	))

	prep.ExpectExec().WithArgs(
		"pending", int64(5678), 31.41459, 21).
		WillReturnResult(sqlmock.NewResult(0, 1))

	prep.ExpectExec().WithArgs(
		"processing", int64(8912), nil, 32).
		WillReturnResult(sqlmock.NewResult(0, 1))
	// </ignore_this>

	// Our objects which should update the columns in the database table
	// `sales_invoice`.

	collection := []*salesInvoice{
		{21, "pending", 5, 5678, null.MakeFloat64(31.41459)},
		{32, "processing", 7, 8912, null.Float64{}},
	}

	// Create the multi update statement
	stmt, err := dbc.Update("sales_invoice").
		AddColumns("state", "customer_id", "grand_total").
		Where(
			dml.Column("shipping_method").In().Strs("DHL", "UPS"), // For all clauses the same restriction
			dml.Column("entity_id").PlaceHolder(),                 // Int64() acts as a place holder
		).
		Prepare(context.TODO())
	assert.NoError(t, err)

	stmtExec := stmt.WithArgs().WithQualifiedColumnsAliases("state", "alias_customer_id", "grand_total", "entity_id")

	for i, record := range collection {
		results, err := stmtExec.Record("sales_invoice", record).ExecContext(context.TODO())
		assert.NoError(t, err)
		ra, err := results.RowsAffected()
		assert.NoError(t, err, "Index %d", i)
		assert.Exactly(t, int64(1), ra, "Index %d", i)
		stmtExec.Reset()
	}

	dbMock.ExpectClose()
	dbc.Close()
	assert.NoError(t, dbMock.ExpectationsWereMet())

}

func TestUpdate_BindRecord(t *testing.T) {
	ce := &categoryEntity{
		EntityID:       678,
		AttributeSetID: 6,
		ParentID:       "p456",
		Path:           null.MakeString("3/4/5"),
	}

	t.Run("1 WHERE", func(t *testing.T) {
		u := dml.NewUpdate("catalog_category_entity").
			AddColumns("attribute_set_id", "parent_id", "path").
			Where(dml.Column("entity_id").Greater().PlaceHolder()).
			WithArgs().Record("", ce)

		compareToSQL(t, u, errors.NoKind,
			"UPDATE `catalog_category_entity` SET `attribute_set_id`=?, `parent_id`=?, `path`=? WHERE (`entity_id` > ?)",
			"UPDATE `catalog_category_entity` SET `attribute_set_id`=6, `parent_id`='p456', `path`='3/4/5' WHERE (`entity_id` > 678)",
			int64(6), "p456", "3/4/5", int64(678),
		)
	})

	t.Run("2 WHERE", func(t *testing.T) {
		u := dml.NewUpdate("catalog_category_entity").
			AddColumns("attribute_set_id", "parent_id", "path").
			Where(
				dml.Column("x").In().Int64s(66, 77),
				dml.Column("entity_id").Greater().PlaceHolder(),
			).WithArgs().Record("", ce)
		compareToSQL(t, u, errors.NoKind,
			"UPDATE `catalog_category_entity` SET `attribute_set_id`=?, `parent_id`=?, `path`=? WHERE (`x` IN (66,77)) AND (`entity_id` > ?)",
			"UPDATE `catalog_category_entity` SET `attribute_set_id`=6, `parent_id`='p456', `path`='3/4/5' WHERE (`x` IN (66,77)) AND (`entity_id` > 678)",
			int64(6), "p456", "3/4/5", int64(678),
		)
	})
	t.Run("3 WHERE", func(t *testing.T) {
		u := dml.NewUpdate("catalog_category_entity").
			AddColumns("attribute_set_id", "parent_id", "path").
			Where(
				dml.Column("entity_id").Greater().PlaceHolder(),
				dml.Column("x").In().Int64s(66, 77),
				dml.Column("y").Greater().Int64(99),
			).WithArgs().Record("", ce)
		compareToSQL(t, u, errors.NoKind,
			"UPDATE `catalog_category_entity` SET `attribute_set_id`=?, `parent_id`=?, `path`=? WHERE (`entity_id` > ?) AND (`x` IN (66,77)) AND (`y` > 99)",
			"UPDATE `catalog_category_entity` SET `attribute_set_id`=6, `parent_id`='p456', `path`='3/4/5' WHERE (`entity_id` > 678) AND (`x` IN (66,77)) AND (`y` > 99)",
			int64(6), "p456", "3/4/5", int64(678),
		)
	})

	t.Run("with alias table name", func(t *testing.T) {
		// A fictional table statement which already reflects future JOIN
		// implementation.
		u := dml.NewUpdate("catalog_category_entity").Alias("ce").
			AddColumns("attribute_set_id", "parent_id", "path").
			Where(
				dml.Column("ce.entity_id").Greater().PlaceHolder(), //678
				dml.Column("cpe.entity_id").In().Int64s(66, 77),
				dml.Column("cpei.attribute_set_id").Equal().PlaceHolder(), //6
			).WithArgs().Records(dml.Qualify("", ce), dml.Qualify("cpei", ce))
		compareToSQL(t, u, errors.NoKind,
			"UPDATE `catalog_category_entity` AS `ce` SET `attribute_set_id`=?, `parent_id`=?, `path`=? WHERE (`ce`.`entity_id` > ?) AND (`cpe`.`entity_id` IN (66,77)) AND (`cpei`.`attribute_set_id` = ?)",
			"UPDATE `catalog_category_entity` AS `ce` SET `attribute_set_id`=6, `parent_id`='p456', `path`='3/4/5' WHERE (`ce`.`entity_id` > 678) AND (`cpe`.`entity_id` IN (66,77)) AND (`cpei`.`attribute_set_id` = 6)",
			int64(6), "p456", "3/4/5", int64(678), int64(6),
		)
	})
}

func TestUpdate_Clone(t *testing.T) {
	t.Parallel()

	dbc, dbMock := dmltest.MockDB(t, dml.WithLogger(log.BlackHole{}, func() string { return "uniqueID" }))
	defer dmltest.MockClose(t, dbc, dbMock)

	t.Run("nil", func(t *testing.T) {
		var d *dml.Update
		d2 := d.Clone()
		assert.Nil(t, d)
		assert.Nil(t, d2)
	})

	t.Run("non-nil", func(t *testing.T) {
		d := dbc.Update("catalog_category_entity").Alias("ce").
			AddColumns("attribute_set_id", "parent_id", "path").
			Where(
				dml.Column("ce.entity_id").Greater().PlaceHolder(), //678
				dml.Column("cpe.entity_id").In().Int64s(66, 77),
				dml.Column("cpei.attribute_set_id").Equal().PlaceHolder(), //6
			)
		d2 := d.Clone()
		notEqualPointers(t, d, d2)
		notEqualPointers(t, d, d2)
		notEqualPointers(t, d.BuilderConditional.Wheres, d2.BuilderConditional.Wheres)
		notEqualPointers(t, d.SetClauses, d2.SetClauses)
		assert.Exactly(t, d.DB, d2.DB)
		assert.Exactly(t, d.Log, d2.Log)
	})
}
