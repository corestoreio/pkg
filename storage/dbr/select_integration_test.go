// Copyright 2015-2016, Cyrill @ Schumacher.fm and the CoreStore contributors
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
	"testing"

	"github.com/corestoreio/csfw/storage/dbr"
	"github.com/corestoreio/csfw/util/cstesting"
	"github.com/corestoreio/csfw/util/errors"
	"github.com/stretchr/testify/assert"
)

func TestSelect_Rows(t *testing.T) {

	t.Run("ToSQL Error", func(t *testing.T) {
		sel := &dbr.Select{}
		sel.Columns = []string{"a", "b"}
		rows, err := sel.Rows()
		assert.Nil(t, rows)
		assert.True(t, errors.IsEmpty(err))
	})

	t.Run("Query Error", func(t *testing.T) {
		sel := &dbr.Select{
			FromTable: dbr.MakeAlias("tableX"),
			Columns:   []string{"a", "b"},
			Querier: dbMock{
				error: errors.NewAlreadyClosedf("Who closed myself?"),
			},
		}

		rows, err := sel.Rows()
		assert.Nil(t, rows)
		assert.True(t, errors.IsAlreadyClosed(err), "%+v", err)
	})
}

func TestSelect_Row(t *testing.T) {

	dbc, dbMock := cstesting.MockDB(t)
	defer func() {
		dbMock.ExpectClose()
		assert.NoError(t, dbc.Close())
		if err := dbMock.ExpectationsWereMet(); err != nil {
			t.Error("there were unfulfilled expections", err)
		}
	}()
	dbMock.ExpectQuery("SELECT a, b FROM `tableX`").WillReturnError(errors.NewAlreadyClosedf("Who closed myself?"))

	sel := &dbr.Select{
		FromTable:  dbr.MakeAlias("tableX"),
		Columns:    []string{"a", "b"},
		QueryRower: dbc.DB,
	}
	row := sel.Row()
	var x string
	err := row.Scan(&x)
	assert.True(t, errors.IsAlreadyClosed(err), "%+v", err)
}

func TestSelect_Prepare(t *testing.T) {

	t.Run("ToSQL Error", func(t *testing.T) {
		sel := &dbr.Select{}
		sel.Columns = []string{"a", "b"}
		stmt, err := sel.Prepare()
		assert.Nil(t, stmt)
		assert.True(t, errors.IsEmpty(err))
	})

	t.Run("Prepare Error", func(t *testing.T) {
		dbc, dbMock := cstesting.MockDB(t)
		defer func() {
			dbMock.ExpectClose()
			assert.NoError(t, dbc.Close())
			if err := dbMock.ExpectationsWereMet(); err != nil {
				t.Error("there were unfulfilled expections", err)
			}
		}()
		dbMock.ExpectPrepare("SELECT a, b FROM `tableX`").WillReturnError(errors.NewAlreadyClosedf("Who closed myself?"))

		sel := &dbr.Select{
			FromTable: dbr.MakeAlias("tableX"),
			Columns:   []string{"a", "b"},
			Preparer:  dbc.DB,
		}
		stmt, err := sel.Prepare()
		assert.Nil(t, stmt)
		assert.True(t, errors.IsAlreadyClosed(err), "%+v", err)
	})

}
