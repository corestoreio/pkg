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
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/corestoreio/errors"
	"github.com/corestoreio/pkg/sql/dml"
	"github.com/corestoreio/pkg/sql/dmltest"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDelete_Prepare(t *testing.T) {
	t.Parallel()

	t.Run("ToSQL Error", func(t *testing.T) {
		compareToSQL(t, dml.NewDelete("").Where(dml.Column("a").Int64(1)), errors.Empty, "", "")
	})

	t.Run("Prepare Error", func(t *testing.T) {
		d := &dml.Delete{
			BuilderBase: dml.BuilderBase{
				Table: dml.MakeIdentifier("table"),
			},
		}
		d.WithDB(dbMock{
			error: errors.AlreadyClosed.Newf("Who closed myself?"),
		})

		d.Where(dml.Column("a").Int(1))
		stmt, err := d.Prepare(context.TODO())
		assert.Nil(t, stmt)
		assert.True(t, errors.AlreadyClosed.Match(err), "%+v", err)
	})

	t.Run("ExecArgs One Row", func(t *testing.T) {
		dbc, dbMock := dmltest.MockDB(t)
		defer dmltest.MockClose(t, dbc, dbMock)

		prep := dbMock.ExpectPrepare(dmltest.SQLMockQuoteMeta("DELETE FROM `customer_entity` WHERE (`email` = ?) AND (`group_id` = ?)"))
		prep.ExpectExec().WithArgs("a@b.c", 33).WillReturnResult(sqlmock.NewResult(0, 1))
		prep.ExpectExec().WithArgs("x@y.z", 44).WillReturnResult(sqlmock.NewResult(0, 2))

		stmt, err := dml.NewDelete("customer_entity").
			Where(dml.Column("email").PlaceHolder(), dml.Column("group_id").PlaceHolder()).
			WithDB(dbc.DB).
			Prepare(context.TODO())
		require.NoError(t, err, "failed creating a prepared statement")
		defer dmltest.Close(t, stmt)

		tests := []struct {
			email   string
			groupID int
			affRows int64
		}{
			{"a@b.c", 33, 1},
			{"x@y.z", 44, 2},
		}

		args := dml.MakeArgs(3)
		for i, test := range tests {
			args = args.Reset()

			res, err := stmt.WithArgs().String(test.email).Int(test.groupID).ExecContext(context.TODO())
			if err != nil {
				t.Fatalf("Index %d => %+v", i, err)
			}
			ra, err := res.RowsAffected()
			if err != nil {
				t.Fatalf("Result index %d with error: %s", i, err)
			}
			assert.Exactly(t, test.affRows, ra, "Index %d has different RowsAffected", i)
		}
	})

	t.Run("ExecRecord One Row", func(t *testing.T) {
		dbc, dbMock := dmltest.MockDB(t)
		defer dmltest.MockClose(t, dbc, dbMock)

		prep := dbMock.ExpectPrepare(dmltest.SQLMockQuoteMeta("DELETE FROM `dml_person` WHERE (`name` = ?) AND (`email` = ?)"))
		prep.ExpectExec().WithArgs("Peter Gopher", "peter@gopher.go").WillReturnResult(sqlmock.NewResult(0, 4))
		prep.ExpectExec().WithArgs("John Doe", "john@doe.go").WillReturnResult(sqlmock.NewResult(0, 5))

		stmt, err := dml.NewDelete("dml_person").
			Where(dml.Column("name").PlaceHolder(), dml.Column("email").PlaceHolder()).
			WithDB(dbc.DB).
			Prepare(context.TODO())
		require.NoError(t, err, "failed creating a prepared statement")
		defer func() {
			require.NoError(t, stmt.Close(), "Close on a prepared statement")
		}()

		tests := []struct {
			name     string
			email    string
			insertID int64
		}{
			{"Peter Gopher", "peter@gopher.go", 4},
			{"John Doe", "john@doe.go", 5},
		}

		for i, test := range tests {

			p := &dmlPerson{
				Name:  test.name,
				Email: dml.MakeNullString(test.email),
			}

			res, err := stmt.WithArgs().Record("", p).ExecContext(context.TODO())
			if err != nil {
				t.Fatalf("Index %d => %+v", i, err)
			}
			lid, err := res.RowsAffected()
			if err != nil {
				t.Fatalf("Result index %d with error: %s", i, err)
			}
			assert.Exactly(t, test.insertID, lid, "Index %d has different RowsAffected", i)
		}
	})

	t.Run("ExecContext", func(t *testing.T) {
		dbc, dbMock := dmltest.MockDB(t)
		defer dmltest.MockClose(t, dbc, dbMock)

		prep := dbMock.ExpectPrepare(dmltest.SQLMockQuoteMeta("DELETE FROM `dml_person` WHERE (`name` = ?) AND (`email` = ?)"))
		prep.ExpectExec().WithArgs("Peter Gopher", "peter@gopher.go").WillReturnResult(sqlmock.NewResult(0, 4))

		stmt, err := dml.NewDelete("dml_person").
			Where(dml.Column("name").PlaceHolder(), dml.Column("email").PlaceHolder()).
			WithDB(dbc.DB).
			Prepare(context.TODO())
		require.NoError(t, err, "failed creating a prepared statement")
		defer func() {
			require.NoError(t, stmt.Close(), "Close on a prepared statement")
		}()

		res, err := stmt.WithArgs().ExecContext(context.TODO(), "Peter Gopher", "peter@gopher.go")
		require.NoError(t, err, "failed to execute ExecContext")

		lid, err := res.RowsAffected()
		if err != nil {
			t.Fatal(err)
		}
		assert.Exactly(t, int64(4), lid, "Different RowsAffected")
	})
}
