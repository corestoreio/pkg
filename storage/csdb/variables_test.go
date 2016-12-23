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
	"encoding/json"
	"sort"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/corestoreio/csfw/util/cstesting"
	"github.com/corestoreio/errors"
	"github.com/stretchr/testify/assert"
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

	var mockedRows = sqlmock.NewRows([]string{"Variable_name", "Value"}).
		FromCSVString("keyVal01,helloWorld")

	dbMock.ExpectQuery("SHOW SESSION VARIABLES.+").
		WithArgs(("keyVal01")).
		WillReturnRows(mockedRows)

	v := &Variable{}
	err := v.LoadOne(dbc.DB, "keyVal01")
	if err != nil {
		t.Fatalf("%+v", err)
	}
	assert.Exactly(t, "keyVal01", v.Name)
	assert.Exactly(t, "helloWorld", v.Value)

	err = v.LoadOne(nil, "test__var")
	assert.True(t, errors.IsNotValid(err), "%+v", err)
}

func TestVariables_AppendFiltered(t *testing.T) {
	dbc, dbMock := cstesting.MockDB(t)

	t.Run("One", func(t *testing.T) {
		var mockedRows = sqlmock.NewRows([]string{"Variable_name", "Value"}).
			FromCSVString("keyVal11,helloAustralia")

		dbMock.ExpectQuery("SHOW SESSION VARIABLES LIKE.+").
			WithArgs(("keyVal11")).WillReturnRows(mockedRows)

		var vs Variables
		if err := vs.AppendFiltered(dbc.DB, "keyVal11"); err != nil {
			t.Fatalf("%+v", err)
		}
		assert.Exactly(t, `keyVal11`, vs.FindOne("keyVal11").Name)
		assert.Exactly(t, `helloAustralia`, vs.FindOne("keyVal11").Value)
		assert.Len(t, vs, 1)
	})

	t.Run("All", func(t *testing.T) {
		var mockedRows = sqlmock.NewRows([]string{"Variable_name", "Value"}).
			FromCSVString("keyVal01,helloWorld\nkeyVal10,helloGermany\nkeyVal11,helloAustralia")

		dbMock.ExpectQuery("SHOW SESSION VARIABLES").
			WillReturnRows(mockedRows)

		var vs Variables
		if err := vs.AppendFiltered(dbc.DB, ""); err != nil {
			t.Fatalf("%+v", err)
		}
		js, err := json.Marshal(vs)
		if err != nil {
			t.Fatal(err)
		}
		assert.Exactly(t, `[{"Name":"keyVal01","Value":"helloWorld"},{"Name":"keyVal10","Value":"helloGermany"},{"Name":"keyVal11","Value":"helloAustralia"}]`, string(js))
		assert.Len(t, vs, 3)
	})

	t.Run("Invalid Name", func(t *testing.T) {
		var vs Variables
		err := vs.AppendFiltered(dbc.DB, "")
		assert.True(t, errors.IsNotValid(err), "%+v", err)
	})

	dbMock.ExpectClose()
	assert.NoError(t, dbc.Close())
	if err := dbMock.ExpectationsWereMet(); err != nil {
		t.Error("there were unfulfilled expections", err)
	}
}
