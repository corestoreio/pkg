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
)

// due to import cycle with the cstesting package, we must test externally

func TestUpdateMulti_Exec(t *testing.T) {

	t.Run("no columns provided", func(t *testing.T) {
		mu := dbr.NewUpdateMulti("catalog_product_entity", "cpe")
		mu.Update.Where(dbr.Column("entity_id", dbr.In.Int64())) // ArgInt64 must be without arguments
		res, err := mu.Exec(context.TODO())
		assert.Nil(t, res)
		assert.True(t, errors.IsEmpty(err), "%+v", err)
	})
	t.Run("alias mismatch", func(t *testing.T) {
		mu := dbr.NewUpdateMulti("catalog_product_entity", "cpe")
		mu.Update.SetClauses.Columns = []string{"sku", "updated_at"}
		mu.Update.Where(dbr.Column("entity_id", dbr.In.Int64())) // ArgInt64 must be without arguments
		mu.Alias = []string{"update_sku"}
		res, err := mu.Exec(context.TODO())
		assert.Nil(t, res)
		assert.True(t, errors.IsMismatch(err), "%+v", err)
	})
	t.Run("empty Records and RecordChan", func(t *testing.T) {
		mu := dbr.NewUpdateMulti("catalog_product_entity", "cpe")
		mu.Update.SetClauses.Columns = []string{"sku", "updated_at"}
		mu.Update.Where(dbr.Column("entity_id", dbr.In.Int64())) // ArgInt64 must be without arguments
		res, err := mu.Exec(context.TODO())
		assert.Nil(t, res)
		assert.True(t, errors.IsEmpty(err), "%+v", err)
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

	mu := dbr.NewUpdateMulti("customer_entity", "ce")
	mu.Update.SetClauses.Columns = []string{"name", "email"}
	mu.Update.Where(dbr.Column("id", dbr.Equal.Int64())) // Int64 must be without arguments, because Placeholder
	mu.Update.Interpolate()

	mu.Records = append(mu.Records, records...)

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
		defer func() {
			dbMock.ExpectClose()
			assert.NoError(t, dbc.Close())
			if err := dbMock.ExpectationsWereMet(); err != nil {
				t.Error("there were unfulfilled expections", err)
			}
		}()

		setSQLMockInterpolate(dbMock)

		mu.Update.DB.Execer = dbc.DB
		mu.Update.DB.Preparer = nil

		results, err := mu.Exec(context.TODO())
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
		defer func() {
			dbMock.ExpectClose()
			assert.NoError(t, dbc.Close())
			if err := dbMock.ExpectationsWereMet(); err != nil {
				t.Error("there were unfulfilled expections", err)
			}
		}()

		setSMPrepared(dbMock)

		mu.Update.IsInterpolate = false
		mu.Update.DB.Execer = nil
		mu.Update.DB.Preparer = dbc.DB

		results, err := mu.Exec(context.TODO())
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
		defer func() {
			dbMock.ExpectClose()
			assert.NoError(t, dbc.Close())
			if err := dbMock.ExpectationsWereMet(); err != nil {
				t.Error("there were unfulfilled expections", err)
			}
		}()

		dbMock.ExpectBegin()
		setSMPrepared(dbMock)
		dbMock.ExpectCommit()

		mu.Tx = dbc.DB
		mu.Transaction()
		mu.Update.IsInterpolate = false
		mu.Update.DB.Execer = nil
		mu.Update.DB.Preparer = dbc.DB

		results, err := mu.Exec(context.TODO())
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
		defer func() {
			dbMock.ExpectClose()
			assert.NoError(t, dbc.Close())
			if err := dbMock.ExpectationsWereMet(); err != nil {
				t.Error("there were unfulfilled expections", err)
			}
		}()

		dbMock.ExpectBegin()
		setSQLMockInterpolate(dbMock)
		dbMock.ExpectCommit()

		mu.Tx = dbc.DB
		mu.Transaction()
		mu.Update.IsInterpolate = true
		mu.Update.DB.Execer = nil
		mu.Update.DB.Preparer = dbc.DB

		results, err := mu.Exec(context.TODO())
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

func TestUpdate_SetRecord_Arguments(t *testing.T) {
	ce := &categoryEntity{
		EntityID:       678,
		AttributeSetID: 6,
		ParentID:       "p456",
		Path:           dbr.MakeNullString("3/4/5"),
	}

	t.Run("1 WHERE", func(t *testing.T) {
		u := dbr.NewUpdate("catalog_category_entity").
			SetRecord([]string{"attribute_set_id", "parent_id", "path"}, ce).
			Where(dbr.Column("entity_id", dbr.Greater.Int64())) // No Arguments in Int64 because we need a place holder.

		compareToSQL(t, u, nil,
			"UPDATE `catalog_category_entity` SET `attribute_set_id`=?, `parent_id`=?, `path`=? WHERE (`entity_id` > ?)",
			"UPDATE `catalog_category_entity` SET `attribute_set_id`=6, `parent_id`='p456', `path`='3/4/5' WHERE (`entity_id` > 678)",
			int64(6), "p456", "3/4/5", int64(678),
		)
	})

	t.Run("2 WHERE", func(t *testing.T) {
		u := dbr.NewUpdate("catalog_category_entity").
			SetRecord([]string{"attribute_set_id", "parent_id", "path"}, ce).
			Where(
				dbr.Column("x", dbr.In.Int64(66, 77)),
				dbr.Column("entity_id", dbr.Greater.Int64()),
			)
		compareToSQL(t, u, nil,
			"UPDATE `catalog_category_entity` SET `attribute_set_id`=?, `parent_id`=?, `path`=? WHERE (`x` IN ?) AND (`entity_id` > ?)",
			"UPDATE `catalog_category_entity` SET `attribute_set_id`=6, `parent_id`='p456', `path`='3/4/5' WHERE (`x` IN (66,77)) AND (`entity_id` > 678)",
			int64(6), "p456", "3/4/5", int64(66), int64(77), int64(678),
		)
	})
	t.Run("3 WHERE", func(t *testing.T) {
		u := dbr.NewUpdate("catalog_category_entity").
			SetRecord([]string{"attribute_set_id", "parent_id", "path"}, ce).
			Where(
				dbr.Column("entity_id", dbr.Greater.Int64()),
				dbr.Column("x", dbr.In.Int64(66, 77)),
				dbr.Column("y", dbr.Greater.Int64(99)),
			)
		compareToSQL(t, u, nil,
			"UPDATE `catalog_category_entity` SET `attribute_set_id`=?, `parent_id`=?, `path`=? WHERE (`entity_id` > ?) AND (`x` IN ?) AND (`y` > ?)",
			"UPDATE `catalog_category_entity` SET `attribute_set_id`=6, `parent_id`='p456', `path`='3/4/5' WHERE (`entity_id` > 678) AND (`x` IN (66,77)) AND (`y` > 99)",
			int64(6), "p456", "3/4/5", int64(678), int64(66), int64(77), int64(99),
		)
	})

}
