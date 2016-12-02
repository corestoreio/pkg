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
	"context"
	"database/sql"
	"testing"

	"github.com/corestoreio/csfw/storage/dbr"
	"github.com/corestoreio/csfw/util/cstesting"
	"github.com/corestoreio/csfw/util/errors"
	"github.com/stretchr/testify/assert"
)

var _ dbr.Preparer = (*sql.DB)(nil)
var _ dbr.Querier = (*sql.DB)(nil)
var _ dbr.Execer = (*sql.DB)(nil)
var _ dbr.QueryRower = (*sql.DB)(nil)

func TestWrapPrepareContext(t *testing.T) {

	dbc, dbMock := cstesting.MockDB(t)
	defer func() {
		dbMock.ExpectClose()
		assert.NoError(t, dbc.Close())
		if err := dbMock.ExpectationsWereMet(); err != nil {
			t.Error("there were unfulfilled expections", err)
		}
	}()
	dbMock.ExpectPrepare("INSERT INTO TABLE abc").WillReturnError(errors.New("Upssss"))

	pc := dbr.WrapPrepareContext(context.TODO(), dbc.DB)
	stmt, err := pc.Prepare("INSERT INTO TABLE abc")
	assert.Nil(t, stmt)
	assert.EqualError(t, err, "Upssss", "%+v", err)
}

func TestWrapQueryContext(t *testing.T) {

	dbc, dbMock := cstesting.MockDB(t)
	defer func() {
		dbMock.ExpectClose()
		assert.NoError(t, dbc.Close())
		if err := dbMock.ExpectationsWereMet(); err != nil {
			t.Error("there were unfulfilled expections", err)
		}
	}()
	dbMock.ExpectQuery("SELECT a FROM tableX").WithArgs(1).WillReturnError(errors.New("Upssss"))

	qc := dbr.WrapQueryContext(context.TODO(), dbc.DB)
	rows, err := qc.Query("SELECT a FROM tableX where b = ?", 1)
	assert.Nil(t, rows)
	assert.EqualError(t, err, "Upssss", "%+v", err)
}

func TestWrapQueryRowContext(t *testing.T) {

	dbc, dbMock := cstesting.MockDB(t)
	defer func() {
		dbMock.ExpectClose()
		assert.NoError(t, dbc.Close())
		if err := dbMock.ExpectationsWereMet(); err != nil {
			t.Error("there were unfulfilled expections", err)
		}
	}()
	dbMock.ExpectQuery("SELECT a FROM tableX").WithArgs(1).WillReturnError(errors.New("Upssss"))

	row := dbr.WrapQueryRowContext(context.TODO(), dbc.DB).QueryRow("SELECT a FROM tableX where b = ?", 1)
	var x string
	err := row.Scan(&x)
	assert.EqualError(t, err, "Upssss", "%+v", err)
}

func TestWrapExecContext(t *testing.T) {

	dbc, dbMock := cstesting.MockDB(t)
	defer func() {
		dbMock.ExpectClose()
		assert.NoError(t, dbc.Close())
		if err := dbMock.ExpectationsWereMet(); err != nil {
			t.Error("there were unfulfilled expections", err)
		}
	}()
	dbMock.ExpectExec("ALTER TABLE add").WithArgs(1).WillReturnError(errors.New("Upssss"))

	res, err := dbr.WrapExecContext(context.TODO(), dbc.DB).Exec("ALTER TABLE add a = ?", 1)
	assert.Nil(t, res)
	assert.EqualError(t, err, "Upssss", "%+v", err)
}
