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

package csdb

import (
	"context"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/corestoreio/csfw/storage/dbr"
	"github.com/corestoreio/csfw/util/cstesting"
	"github.com/corestoreio/errors"
	"github.com/stretchr/testify/assert"
)

var _ dbr.Scanner = (*Variables)(nil)
var _ dbr.QueryBuilder = (*Variables)(nil)

func TestIsValidVarName(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name         string
		allowPercent bool
		errBhf       errors.BehaviourFunc
	}{
		{"auto_increment_offset", false, nil},
		{"auto_increment_offset", true, nil},
		{"auto_increment%", true, nil},
		{"auto_increment%", false, errors.IsNotValid},
		{"auto_in`crement%", true, errors.IsNotValid},
		{"auto_increment%", true, errors.IsNotValid},
		{"auto_in'crement%", true, errors.IsNotValid},
		{"", true, nil},
		{"", false, nil},
	}
	for _, test := range tests {

		haveErr := isValidVarName(test.name, test.allowPercent)
		if test.errBhf == nil {
			assert.NoError(t, haveErr, "%+v", haveErr)
			continue
		}
		assert.True(t, test.errBhf(haveErr), "%+v", haveErr)
	}
}

func TestNewVariables_Integration(t *testing.T) {
	t.Parallel()

	db := cstesting.MustConnectDB(t)
	defer cstesting.Close(t, db)

	vs := NewVariables()
	_, err := dbr.Load(context.TODO(), db.DB, vs, vs)
	if err != nil {
		t.Fatalf("%+v", err)
	}
	assert.Exactly(t, "InnoDB", vs.Data["storage_engine"])
	assert.True(t, len(vs.Data) > 400, "Should have more than 400 map entries")
}

func TestNewVariables_Mock(t *testing.T) {
	t.Parallel()

	dbc, dbMock := cstesting.MockDB(t)
	defer cstesting.MockClose(t, dbc, dbMock)

	t.Run("one with LIKE", func(t *testing.T) {
		var mockedRows = sqlmock.NewRows([]string{"Variable_name", "Value"}).
			FromCSVString("keyVal11,helloAustralia")

		dbMock.ExpectQuery(cstesting.SQLMockQuoteMeta("SHOW VARIABLES WHERE (`Variable_name` LIKE 'keyVal11')")).
			WillReturnRows(mockedRows)

		vs := NewVariables("keyVal11")
		rc, err := dbr.Load(context.TODO(), dbc.DB, vs, vs)
		if err != nil {
			t.Fatalf("%+v", err)
		}
		assert.Exactly(t, int64(1), rc, "Should load one row")

		assert.Exactly(t, `helloAustralia`, vs.Data["keyVal11"])
		assert.Len(t, vs.Data, 1)
	})

	t.Run("many with WHERE", func(t *testing.T) {
		var mockedRows = sqlmock.NewRows([]string{"Variable_name", "Value"}).
			FromCSVString("keyVal11,helloAustralia\nkeyVal22,helloNewZealand")

		dbMock.ExpectQuery(cstesting.SQLMockQuoteMeta("SHOW VARIABLES WHERE (`Variable_name` IN ('keyVal11','keyVal22'))")).
			WillReturnRows(mockedRows)

		vs := NewVariables("keyVal11", "keyVal22")
		rc, err := dbr.Load(context.TODO(), dbc.DB, vs, vs)
		if err != nil {
			t.Fatalf("%+v", err)
		}
		assert.Exactly(t, int64(2), rc, "Shoud load two rows")

		assert.Exactly(t, `helloAustralia`, vs.Data["keyVal11"])
		assert.Exactly(t, `helloNewZealand`, vs.Data["keyVal22"])
		assert.Len(t, vs.Data, 2)
	})
}
