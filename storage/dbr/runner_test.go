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

	"github.com/corestoreio/csfw/storage/dbr"
	"github.com/corestoreio/csfw/util/cstesting"
	"github.com/corestoreio/errors"
	"github.com/stretchr/testify/assert"
)

type myToSQL struct {
	sql  string
	args []interface{}
	error
}

func (m myToSQL) ToSQL() (string, []interface{}, error) {
	return m.sql, m.args, m.error
}

func TestExec(t *testing.T) {
	t.Parallel()
	haveErr := errors.NewAlreadyClosedf("Who closed myself?")

	t.Run("ToSQL error", func(t *testing.T) {
		stmt, err := dbr.Exec(context.TODO(), nil, myToSQL{error: haveErr})
		assert.Nil(t, stmt)
		assert.True(t, errors.IsAlreadyClosed(err), "%+v", err)
	})
}

func TestPrepare(t *testing.T) {
	t.Parallel()
	haveErr := errors.NewAlreadyClosedf("Who closed myself?")

	t.Run("ToSQL error", func(t *testing.T) {
		stmt, err := dbr.Prepare(context.TODO(), nil, myToSQL{error: haveErr})
		assert.Nil(t, stmt)
		assert.True(t, errors.IsAlreadyClosed(err), "%+v", err)
	})
	t.Run("ToSQL prepare error", func(t *testing.T) {
		dbc, dbMock := cstesting.MockDB(t)
		defer func() {
			dbMock.ExpectClose()
			assert.NoError(t, dbc.Close())
			if err := dbMock.ExpectationsWereMet(); err != nil {
				t.Error("there were unfulfilled expections", err)
			}
		}()
		dbMock.ExpectPrepare("SELECT `a` FROM `b`").WillReturnError(haveErr)

		stmt, err := dbr.Prepare(context.TODO(), dbc.DB, myToSQL{sql: "SELECT `a` FROM `b`"})
		assert.Nil(t, stmt)
		assert.True(t, errors.IsAlreadyClosed(err), "%+v", err)
	})
}
