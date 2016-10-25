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

package csdb_test

import (
	"context"
	"testing"

	sqlmock "github.com/DATA-DOG/go-sqlmock"
	"github.com/corestoreio/csfw/storage/csdb"
	"github.com/corestoreio/csfw/util/cstesting"
	"github.com/stretchr/testify/assert"
)

func TestShowMasterStatus(t *testing.T) {
	dbc, dbMock := cstesting.MockDB(t)
	defer func() {
		dbMock.ExpectClose()
		assert.NoError(t, dbc.Close())
		if err := dbMock.ExpectationsWereMet(); err != nil {
			t.Error("there were unfulfilled expections", err)
		}
	}()

	var mockedRows = sqlmock.NewRows([]string{"File", "Position", "Binlog_Do_DB", "Binlog_Ignore_DB", "Executed_Gtid_Set"}).
		FromCSVString("mysql-bin.000001,3581378,test,mysql,123-456-789")

	dbMock.ExpectQuery("SHOW MASTER STATUS").WillReturnRows(mockedRows)

	v := new(csdb.MasterStatus)
	err := v.Load(context.TODO(), dbc.DB)
	if err != nil {
		t.Fatalf("%+v", err)
	}
	assert.Exactly(t, "mysql-bin.000001", v.File)
	assert.Exactly(t, uint(3581378), v.Position)
	assert.Exactly(t, "123-456-789", v.Executed_Gtid_Set)
}
