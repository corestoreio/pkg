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

package ddl_test

import (
	"bytes"
	"context"
	"io"
	"testing"

	sqlmock "github.com/DATA-DOG/go-sqlmock"
	"github.com/corestoreio/errors"
	"github.com/corestoreio/pkg/sql/ddl"
	"github.com/corestoreio/pkg/sql/dml"
	"github.com/corestoreio/pkg/sql/dmltest"
	"github.com/corestoreio/pkg/util/assert"
)

var _ dml.QueryBuilder = (*ddl.MasterStatus)(nil)
var _ dml.ColumnMapper = (*ddl.MasterStatus)(nil)
var _ io.WriterTo = (*ddl.MasterStatus)(nil)

func TestMasterStatus_Compare(t *testing.T) {
	t.Parallel()
	tests := []struct {
		left, right ddl.MasterStatus
		want        int
	}{
		{ddl.MasterStatus{File: "mysql-bin.000001", Position: 3}, ddl.MasterStatus{File: "mysql-bin.000001", Position: 4}, -1},
		{ddl.MasterStatus{File: "mysql-bin.000001", Position: 3}, ddl.MasterStatus{File: "mysql-bin.000001", Position: 3}, 0},
		{ddl.MasterStatus{File: "mysql-bin.000001", Position: 3}, ddl.MasterStatus{File: "mysql-bin.000001", Position: 2}, 1},
		{ddl.MasterStatus{File: "mysql-bin.000001", Position: 3}, ddl.MasterStatus{File: "mysql-bin.000002", Position: 2}, -1},
		{ddl.MasterStatus{File: "mysql-bin.000003", Position: 1}, ddl.MasterStatus{File: "mysql-bin.000002", Position: 2}, 1},
	}
	for i, test := range tests {
		have := test.left.Compare(test.right)
		assert.Exactly(t, test.want, have, "Index %d", i)
	}
}

func TestShowMasterStatus(t *testing.T) {
	t.Parallel()

	dbc, dbMock := dmltest.MockDB(t)
	defer dmltest.MockClose(t, dbc, dbMock)

	mockedRows := sqlmock.NewRows([]string{"File", "Position", "Binlog_Do_DB", "Binlog_Ignore_DB", "Executed_Gtid_Set"}).
		FromCSVString("mysql-bin.000001,3581378,test,mysql,123-456-789")

	dbMock.ExpectQuery("SHOW MASTER STATUS").WillReturnRows(mockedRows)

	v := new(ddl.MasterStatus)
	_, err := dbc.WithQueryBuilder(v).Load(context.TODO(), v)
	if err != nil {
		t.Fatalf("%+v", err)
	}
	assert.Exactly(t, "mysql-bin.000001", v.File)
	assert.Exactly(t, uint(3581378), v.Position)
	assert.Exactly(t, "123-456-789", v.ExecutedGTIDSet)
}

func TestMasterStatus_FromString(t *testing.T) {
	t.Parallel()
	tests := []struct {
		in           string
		wantFile     string
		wantPosition uint
		wantErrKind  errors.Kind
		wantString   string
	}{
		{"mysql-bin.000004;545460", "mysql-bin.000004", 545460, errors.NoKind, "mysql-bin.000004;545460"},
		{"mysql-bin.000004;", "", 0, errors.NotValid, ""},
		{"mysql-bin.000004", "", 0, errors.NotFound, ""},
	}
	for i, test := range tests {
		haveMS := &ddl.MasterStatus{}
		haveErr := haveMS.FromString(test.in)
		if !test.wantErrKind.Empty() {
			assert.True(t, test.wantErrKind.Match(haveErr), "Index %d: %+v", i, haveErr)
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

func TestMasterStatus_WriteTo(t *testing.T) {
	t.Parallel()

	ms := &ddl.MasterStatus{}
	err := ms.FromString("mysql-bin.000004;545460")
	assert.NoError(t, err)
	var buf bytes.Buffer

	_, err = ms.WriteTo(&buf)
	assert.NoError(t, err)

	assert.Exactly(t, "mysql-bin.000004;545460", buf.String())
}
