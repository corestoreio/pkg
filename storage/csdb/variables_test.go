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
	"database/sql/driver"
	"encoding/json"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/corestoreio/csfw/util/cstesting"
	"github.com/corestoreio/csfw/util/errors"
	"github.com/stretchr/testify/assert"
	"sort"
	"testing"
)

var _ sort.Interface = (*Variables)(nil)

func TestIsValidVarName(t *testing.T) {
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

func TestVariables_FindOne(t *testing.T) {
	vs := Variables{
		&Variable{Name: "a", Value: "1"},
		&Variable{Name: "b", Value: "2"},
	}
	assert.Exactly(t, "1", vs.FindOne("a").Value)
	assert.Exactly(t, "2", vs.FindOne("b").Value)
	assert.Exactly(t, "", vs.FindOne("c").Value)
}

func TestVariables_Sort(t *testing.T) {
	vs := Variables{
		&Variable{Name: "d", Value: "4"},
		&Variable{Name: "a", Value: "1"},
		&Variable{Name: "c", Value: "3"},
		&Variable{Name: "b", Value: "2"},
	}
	sort.Stable(vs)
	jb, err := json.Marshal(vs)
	assert.NoError(t, err)
	assert.Exactly(t, `[{"Name":"a","Value":"1"},{"Name":"b","Value":"2"},{"Name":"c","Value":"3"},{"Name":"d","Value":"4"}]`, string(jb))
}

func TestVariable_LoadOne(t *testing.T) {
	dbc, dbMock := cstesting.MockDB(t)
	defer func() {
		dbMock.ExpectClose()

		assert.NoError(t, dbc.Close())

		if err := dbMock.ExpectationsWereMet(); err != nil {
			t.Error("there were unfulfilled expections", err)
		}
	}()

	mockRows := sqlmock.NewRows([]string{"Variable_name", "Value"})
	mockRows.FromCSVString("test_x_var,helloWorld")
	dbMock.ExpectQuery("SHOW SESSn").
		WithArgs(driver.Value("test_x_var")).
		WillReturnRows(mockRows)

	v := &Variable{}
	v.LoadOne(context.TODO(), dbc.DB, "test_x_var")

	t.Logf("%#v", v)
}
