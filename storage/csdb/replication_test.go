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
	"testing"

	sqlmock "github.com/DATA-DOG/go-sqlmock"
	"github.com/corestoreio/csfw/storage/csdb"
	"github.com/corestoreio/csfw/util/cstesting"
	"github.com/corestoreio/csfw/util/errors"
	"github.com/stretchr/testify/assert"
)

func TestMasterStatus_Compare(t *testing.T) {
	tests := []struct {
		left, right csdb.MasterStatus
		want        int
	}{
		{csdb.MasterStatus{File: "mysql-bin.000001", Position: 3}, csdb.MasterStatus{File: "mysql-bin.000001", Position: 4}, -1},
		{csdb.MasterStatus{File: "mysql-bin.000001", Position: 3}, csdb.MasterStatus{File: "mysql-bin.000001", Position: 3}, 0},
		{csdb.MasterStatus{File: "mysql-bin.000001", Position: 3}, csdb.MasterStatus{File: "mysql-bin.000001", Position: 2}, 1},
		{csdb.MasterStatus{File: "mysql-bin.000001", Position: 3}, csdb.MasterStatus{File: "mysql-bin.000002", Position: 2}, -1},
		{csdb.MasterStatus{File: "mysql-bin.000003", Position: 1}, csdb.MasterStatus{File: "mysql-bin.000002", Position: 2}, 1},
	}
	for i, test := range tests {
		have := test.left.Compare(test.right)
		assert.Exactly(t, test.want, have, "Index %d", i)
	}
}

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
	err := v.Load(dbc.DB)
	if err != nil {
		t.Fatalf("%+v", err)
	}
	assert.Exactly(t, "mysql-bin.000001", v.File)
	assert.Exactly(t, uint(3581378), v.Position)
	assert.Exactly(t, "123-456-789", v.ExecutedGTIDSet)
}

func TestMasterStatus_FromString(t *testing.T) {
	tests := []struct {
		in           string
		wantFile     string
		wantPosition uint
		wantErr      errors.BehaviourFunc
		wantString   string
	}{
		{"mysql-bin.000004;545460", "mysql-bin.000004", 545460, nil, "mysql-bin.000004;545460"},
		{"mysql-bin.000004;", "", 0, errors.IsNotValid, ""},
		{"mysql-bin.000004", "", 0, errors.IsNotFound, ""},
	}
	for i, test := range tests {
		var haveMS = &csdb.MasterStatus{}
		haveErr := haveMS.FromString(test.in)
		if test.wantErr != nil {
			assert.True(t, test.wantErr(haveErr), "Index %d: %+v", i, haveErr)
			assert.Empty(t, haveMS.File, "Index %d", i)
			assert.Empty(t, haveMS.Position, "Index %d", i)
			assert.Empty(t, haveMS.String(), "Index %d", i)
			continue
		}
		assert.NoError(t, haveErr, "Index %d", i)
		assert.Exactly(t, test.wantFile, haveMS.File, "Index %d", i)
		assert.Exactly(t, test.wantPosition, haveMS.Position, "Index %d", i)
		assert.Exactly(t, test.wantString, haveMS.String())
	}
}
