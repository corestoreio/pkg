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

	sqlmock "github.com/DATA-DOG/go-sqlmock"
	"github.com/corestoreio/csfw/storage/dbr"
	"github.com/corestoreio/csfw/util/cstesting"
	"github.com/corestoreio/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestTx_Wrap(t *testing.T) {
	t.Parallel()

	t.Run("commit", func(t *testing.T) {
		dbc, dbMock := cstesting.MockDB(t)
		defer cstesting.MockClose(t, dbc, dbMock)

		dbMock.ExpectBegin()
		dbMock.ExpectExec("UPDATE `tableX` SET `value`").WillReturnResult(sqlmock.NewResult(0, 9))
		dbMock.ExpectCommit()

		tx, err := dbc.Begin()
		require.NoError(t, err)

		require.NoError(t, tx.Wrap(func() error {
			res, err := tx.Update("tableX").Set("value", dbr.ArgInt(5)).Where(dbr.Column("scope", dbr.Equal.Str("default"))).Exec(context.TODO())
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
		dbc, dbMock := cstesting.MockDB(t)
		defer cstesting.MockClose(t, dbc, dbMock)

		dbMock.ExpectBegin()
		dbMock.ExpectExec("UPDATE `tableX` SET `value`").WillReturnError(errors.NewAbortedf("Sorry dude"))
		dbMock.ExpectRollback()

		tx, err := dbc.Begin()
		require.NoError(t, err)

		err = tx.Wrap(func() error {
			res, err := tx.Update("tableX").Set("value", dbr.ArgInt(5)).Where(dbr.Column("scope", dbr.Equal.Str("default"))).Exec(context.TODO())
			assert.Nil(t, res)
			return err
		})
		assert.True(t, errors.IsAborted(err))
	})

}
