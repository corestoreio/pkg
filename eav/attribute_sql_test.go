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

package eav_test

import (
	"testing"

	"github.com/corestoreio/csfw/codegen"
	"github.com/corestoreio/csfw/eav"
	"github.com/corestoreio/csfw/storage/csdb"
	"github.com/corestoreio/csfw/util/diff"
	"github.com/corestoreio/csfw/util/sqlbeautifier"
	"github.com/stretchr/testify/assert"
)

var testWantGetAttributeSelectSql string = "Please specify a build tag: mage1 or mage2\n$ go test -tags mageX -v ."

func TestGetAttributeSelectSql(t *testing.T) {
	dbc := csdb.MustConnectTest()
	defer dbc.Close()

	dbrSess := dbc.NewSession()
	dbrSelect, err := eav.GetAttributeSelectSql(dbrSess, codegen.NewAddAttrTables(dbc.DB, "customer"), 1, 4)
	if err != nil {
		t.Fatal(err)
	}
	sql, _, err := dbrSelect.ToSQL()
	assert.NoError(t, err)

	_, err = sqlbeautifier.FromString(sql) // check for syntax errors
	if err != nil {
		t.Fatalf("%s\n\n%s\n", err, sql)
	}

	if testWantGetAttributeSelectSql != sql {

		buf, err := sqlbeautifier.FromString(testWantGetAttributeSelectSql)
		if err != nil {
			t.Fatalf("%s\n%s\n", err, testWantGetAttributeSelectSql)
		}
		sql = sqlbeautifier.MustFromString(sql)
		println(diff.MustUnified(buf.String(), sql), "\n")
		t.Fatal(sql)
	}

	// @todo error is that we have column attribute_model in the select list but it should not occur
	// because in codegen it is defined that this column has no usage so we can skip it.
}
