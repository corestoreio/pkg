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
	"bytes"
	"context"
	"fmt"
	"sync/atomic"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/corestoreio/errors"
	"github.com/corestoreio/log/logw"
	"github.com/corestoreio/pkg/sql/dml"
	"github.com/corestoreio/pkg/sql/dmltest"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
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
		mu := dml.NewUpdate("catalog_product_entity").
			AddColumns("sku", "updated_at").
			Where(dml.Column("entity_id").In().PlaceHolder()).WithDB(dbMock{})
		mu.SetClausAliases = []string{"update_sku"}
		res, err := mu.WithArgs().ExecContext(context.TODO())
		assert.Nil(t, res)
		assert.True(t, errors.Mismatch.Match(err), "%+v", err)
	})

	t.Run("alias mismatch Prepare", func(t *testing.T) {
		mu := dml.NewUpdate("catalog_product_entity").
			AddColumns("sku", "updated_at").
			Where(dml.Column("entity_id").In().PlaceHolder()).WithDB(dbMock{})
		mu.SetClausAliases = []string{"update_sku"}
		stmt, err := mu.Prepare(context.TODO())
		assert.Nil(t, stmt)
		assert.True(t, errors.Mismatch.Match(err), "%+v", err)
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
				Email: dml.MakeNullString("alf@m') -- el.mac"),
			},
			{
				ID:    2,
				Name:  "John",
				Email: dml.MakeNullString("john@doe.com"),
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
		require.NoError(t, err)
		for i, record := range records {
			results, err := stmt.WithArgs().Record("ce", record).ExecContext(context.TODO())
			require.NoError(t, err)
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
		require.NoError(t, err, "failed creating a prepared statement")
		defer func() {
			require.NoError(t, stmt.Close(), "Close on a prepared statement")
		}()

		res, err := stmt.WithArgs().ExecContext(context.TODO(), "Peter Gopher", "peter@gopher.go", 3456)
		require.NoError(t, err, "failed to execute ExecContext")

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
	GrandTotal dml.NullFloat64
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

func TestUpdate_SetClausAliases(t *testing.T) {
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
		{21, "pending", 5, 5678, dml.MakeNullFloat64(31.41459)},
		{32, "processing", 7, 8912, dml.NullFloat64{}},
	}

	// Create the multi update statement
	um := dml.NewUpdate("sales_invoice").
		AddColumns("state", "customer_id", "grand_total").
		Where(
			dml.Column("shipping_method").In().Strs("DHL", "UPS"), // For all clauses the same restriction
			dml.Column("entity_id").PlaceHolder(),                 // Int64() acts as a place holder
		).WithDB(dbc.DB)

	um.SetClausAliases = []string{"state", "alias_customer_id", "grand_total"}

	stmt, err := um.Prepare(context.TODO())
	require.NoError(t, err)

	for i, record := range collection {
		results, err := stmt.WithArgs().Record("sales_invoice", record).ExecContext(context.TODO())
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
		Path:           dml.MakeNullString("3/4/5"),
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

func TestUpdate_WithLogger(t *testing.T) {

	t.Skip("TODO Check if duplicated by other WithLogger tests and then remove")

	uniID := new(int32)
	rConn := createRealSession(t)
	defer dmltest.Close(t, rConn)

	var uniqueIDFunc = func() string {
		return fmt.Sprintf("UNIQ%02d", atomic.AddInt32(uniID, 3))
	}

	buf := new(bytes.Buffer)
	lg := logw.NewLog(
		logw.WithLevel(logw.LevelDebug),
		logw.WithWriter(buf),
		logw.WithFlag(0), // no flags at all
	)
	require.NoError(t, rConn.Options(dml.WithLogger(lg, uniqueIDFunc)))

	t.Run("ConnPool", func(t *testing.T) {
		d := rConn.Update("dml_people").Set(
			dml.Column("email").Str("new@email.com"),
		).Where(dml.Column("id").GreaterOrEqual().Float64(78.31))

		t.Run("Exec", func(t *testing.T) {
			defer buf.Reset()
			_, err := d.WithArgs().ExecContext(context.TODO())
			require.NoError(t, err)

			assert.Exactly(t, "DEBUG Exec conn_pool_id: \"UNIQ03\" update_id: \"UNIQ06\" table: \"dml_people\" duration: 0 sql: \"UPDATE /*ID:UNIQ06*/ `dml_people` SET `email`='new@email.com' WHERE (`id` >= 78.31)\"\n",
				buf.String())
		})

		t.Run("Prepare", func(t *testing.T) {
			defer buf.Reset()
			stmt, err := d.Prepare(context.TODO())
			require.NoError(t, err)
			defer stmt.Close()

			assert.Exactly(t, "DEBUG Prepare conn_pool_id: \"UNIQ03\" update_id: \"UNIQ06\" table: \"dml_people\" duration: 0 sql: \"UPDATE /*ID:UNIQ06*/ `dml_people` SET `email`='new@email.com' WHERE (`id` >= 78.31)\"\n",
				buf.String())
		})

		t.Run("Tx Commit", func(t *testing.T) {
			defer buf.Reset()
			tx, err := rConn.BeginTx(context.TODO(), nil)
			require.NoError(t, err)
			require.NoError(t, tx.Wrap(func() error {
				_, err := tx.Update("dml_people").Set(
					dml.Column("email").Str("new@email.com"),
				).Where(dml.Column("id").GreaterOrEqual().Float64(36.56)).WithArgs().ExecContext(context.TODO())
				return err
			}))
			assert.Exactly(t, "DEBUG BeginTx conn_pool_id: \"UNIQ03\" tx_id: \"UNIQ09\"\nDEBUG Exec conn_pool_id: \"UNIQ03\" tx_id: \"UNIQ09\" update_id: \"UNIQ12\" table: \"dml_people\" duration: 0 sql: \"UPDATE /*ID:UNIQ12*/ `dml_people` SET `email`='new@email.com' WHERE (`id` >= 36.56)\"\nDEBUG Commit conn_pool_id: \"UNIQ03\" tx_id: \"UNIQ09\" duration: 0\n",
				buf.String())
		})
	})

	t.Run("Conn", func(t *testing.T) {
		conn, err := rConn.Conn(context.TODO())
		require.NoError(t, err)

		d := conn.Update("dml_people").Set(
			dml.Column("email").Str("new@email.com"),
		).Where(dml.Column("id").GreaterOrEqual().Float64(21.56))

		t.Run("Exec", func(t *testing.T) {
			defer buf.Reset()

			_, err := d.WithArgs().ExecContext(context.TODO())
			require.NoError(t, err)

			assert.Exactly(t, "DEBUG Exec conn_pool_id: \"UNIQ03\" conn_id: \"UNIQ15\" update_id: \"UNIQ18\" table: \"dml_people\" duration: 0 sql: \"UPDATE /*ID:UNIQ18*/ `dml_people` SET `email`='new@email.com' WHERE (`id` >= 21.56)\"\n",
				buf.String())
		})

		t.Run("Prepare", func(t *testing.T) {
			defer buf.Reset()

			stmt, err := d.Prepare(context.TODO())
			require.NoError(t, err)
			defer stmt.Close()

			assert.Exactly(t, "DEBUG Prepare conn_pool_id: \"UNIQ03\" conn_id: \"UNIQ15\" update_id: \"UNIQ18\" table: \"dml_people\" duration: 0 sql: \"UPDATE /*ID:UNIQ18*/ `dml_people` SET `email`='new@email.com' WHERE (`id` >= 21.56)\"\n",
				buf.String())
		})

		t.Run("Prepare Exec", func(t *testing.T) {
			defer buf.Reset()

			stmt, err := d.Prepare(context.TODO())
			require.NoError(t, err)
			defer stmt.Close()

			_, err = stmt.WithArgs().ExecContext(context.TODO())
			require.NoError(t, err)

			assert.Exactly(t, "DEBUG Prepare conn_pool_id: \"UNIQ03\" conn_id: \"UNIQ15\" update_id: \"UNIQ18\" table: \"dml_people\" duration: 0 sql: \"UPDATE /*ID:UNIQ18*/ `dml_people` SET `email`='new@email.com' WHERE (`id` >= 21.56)\"\nDEBUG Exec conn_pool_id: \"UNIQ03\" conn_id: \"UNIQ15\" update_id: \"UNIQ18\" table: \"dml_people\" duration: 0 arg_len: 0\n",
				buf.String())
		})

		t.Run("Tx Commit", func(t *testing.T) {
			defer buf.Reset()
			tx, err := conn.BeginTx(context.TODO(), nil)
			require.NoError(t, err)
			require.NoError(t, tx.Wrap(func() error {
				_, err := tx.Update("dml_people").Set(
					dml.Column("email").Str("new@email.com"),
				).Where(dml.Column("id").GreaterOrEqual().Float64(39.56)).WithArgs().ExecContext(context.TODO())
				return err
			}))

			assert.Exactly(t, "DEBUG BeginTx conn_pool_id: \"UNIQ03\" conn_id: \"UNIQ15\" tx_id: \"UNIQ21\"\nDEBUG Exec conn_pool_id: \"UNIQ03\" conn_id: \"UNIQ15\" tx_id: \"UNIQ21\" update_id: \"UNIQ24\" table: \"dml_people\" duration: 0 sql: \"UPDATE /*ID:UNIQ24*/ `dml_people` SET `email`='new@email.com' WHERE (`id` >= 39.56)\"\nDEBUG Commit conn_pool_id: \"UNIQ03\" conn_id: \"UNIQ15\" tx_id: \"UNIQ21\" duration: 0\n",
				buf.String())
		})

		t.Run("Tx Rollback", func(t *testing.T) {
			defer buf.Reset()
			tx, err := conn.BeginTx(context.TODO(), nil)
			require.NoError(t, err)
			require.Error(t, tx.Wrap(func() error {
				_, err := tx.Update("dml_people").Set(
					dml.Column("email").Str("new@email.com"),
				).Where(dml.Column("id").GreaterOrEqual().PlaceHolder()).WithArgs().ExecContext(context.TODO())
				return err
			}))

			assert.Exactly(t, "DEBUG BeginTx conn_pool_id: \"UNIQ03\" conn_id: \"UNIQ15\" tx_id: \"UNIQ27\"\nDEBUG Exec conn_pool_id: \"UNIQ03\" conn_id: \"UNIQ15\" tx_id: \"UNIQ27\" update_id: \"UNIQ30\" table: \"dml_people\" duration: 0 sql: \"UPDATE /*ID:UNIQ30*/ `dml_people` SET `email`='new@email.com' WHERE (`id` >= ?)\"\nDEBUG Rollback conn_pool_id: \"UNIQ03\" conn_id: \"UNIQ15\" tx_id: \"UNIQ27\" duration: 0\n",
				buf.String())
		})
	})
}
