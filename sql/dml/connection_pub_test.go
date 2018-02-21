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

func TestTx_Wrap(t *testing.T) {
	t.Parallel()

	t.Run("commit", func(t *testing.T) {
		dbc, dbMock := dmltest.MockDB(t)
		defer dmltest.MockClose(t, dbc, dbMock)

		dbMock.ExpectBegin()
		dbMock.ExpectExec("UPDATE `tableX` SET `value`").WithArgs().WillReturnResult(sqlmock.NewResult(0, 9))
		dbMock.ExpectCommit()

		require.NoError(t, dbc.Transaction(context.TODO(), nil, func(tx *dml.Tx) error {
			// this creates an interpolated statement
			res, err := tx.Update("tableX").Set(dml.Column("value").Int(5)).Where(dml.Column("scope").Str("default")).WithArgs().ExecContext(context.TODO())
			if err != nil {
				return err
			}
			af, err := res.RowsAffected()
			if err != nil {
				return err
			}
			assert.Exactly(t, int64(9), af)
			return nil
		}))
	})

	t.Run("rollback", func(t *testing.T) {
		dbc, dbMock := dmltest.MockDB(t)
		defer dmltest.MockClose(t, dbc, dbMock)

		dbMock.ExpectBegin()
		dbMock.ExpectExec("UPDATE `tableX` SET `value`").WithArgs().WillReturnError(errors.Aborted.Newf("Sorry dude"))
		dbMock.ExpectRollback()

		err := dbc.Transaction(context.TODO(), nil, func(tx *dml.Tx) error {
			// Interpolated statement
			res, err := tx.Update("tableX").Set(dml.Column("value").Int(5)).Where(dml.Column("scope").Str("default")).WithArgs().ExecContext(context.TODO())
			assert.Nil(t, res)
			return err
		})
		assert.True(t, errors.Aborted.Match(err))
	})
}

func TestWithRawSQL(t *testing.T) {
	t.Parallel()

	dbc, mock := dmltest.MockDB(t)
	defer dmltest.MockClose(t, dbc, mock)

	t.Run("ConnPool", func(t *testing.T) {

		compareToSQL(t,
			dbc.WithRawSQL("SELECT * FROM users WHERE x = ? AND y IN (?,?,?)").Int(9).Int(5).Int(6).Int(7),
			errors.NoKind,
			"SELECT * FROM users WHERE x = ? AND y IN (?,?,?)",
			"",
			int64(9), int64(5), int64(6), int64(7),
		)

		compareToSQL(t,
			dbc.WithRawSQL("SELECT * FROM users WHERE x = 1"),
			errors.NoKind,
			"SELECT * FROM users WHERE x = 1",
			"",
		)
		compareToSQL(t,
			dbc.WithRawSQL("SELECT * FROM users WHERE x = ? AND y IN ?").ExpandPlaceHolders().Int(9).Ints(5, 6, 7),
			errors.NoKind,
			"SELECT * FROM users WHERE x = ? AND y IN (?,?,?)",
			"",
			int64(9), int64(5), int64(6), int64(7),
		)
		compareToSQL(t,
			dbc.WithRawSQL("SELECT * FROM users WHERE x = ? AND y IN ?").Interpolate().Int(9).Ints(5, 6, 7),
			errors.NoKind,
			"SELECT * FROM users WHERE x = 9 AND y IN (5,6,7)",
			"",
		)
		compareToSQL(t,
			dbc.WithRawSQL("wat").Raw(9, 5, 6, 7),
			errors.NoKind,
			"wat",
			"",
			9, 5, 6, 7,
		)
	})

	t.Run("ConnSingle", func(t *testing.T) {
		c, err := dbc.Conn(context.TODO())
		defer dmltest.Close(t, c)
		if err != nil {
			t.Fatal(err)
		}
		compareToSQL(t,
			c.WithRawSQL("SELECT * FROM users WHERE x = ? AND y IN ?").Interpolate().Int(9).Ints(5, 6, 7),
			errors.NoKind,
			"SELECT * FROM users WHERE x = 9 AND y IN (5,6,7)",
			"",
		)
		compareToSQL(t,
			c.WithRawSQL("wat").Raw(9, 5, 6, 7),
			errors.NoKind,
			"wat",
			"",
			9, 5, 6, 7,
		)
	})

	t.Run("Tx", func(t *testing.T) {
		mock.ExpectBegin()
		mock.ExpectCommit()

		tx, err := dbc.BeginTx(context.TODO(), nil)
		if err != nil {
			t.Fatal(err)
		}
		defer func() { assert.NoError(t, tx.Commit()) }()
		compareToSQL(t,
			tx.WithRawSQL("SELECT * FROM users WHERE x = ? AND y IN ?").Interpolate().Int(9).Ints(5, 6, 7),
			errors.NoKind,
			"SELECT * FROM users WHERE x = 9 AND y IN (5,6,7)",
			"",
		)
		compareToSQL(t,
			tx.WithRawSQL("wat").Raw(9, 5, 6, 7),
			errors.NoKind,
			"wat",
			"",
			9, 5, 6, 7,
		)
	})
}
