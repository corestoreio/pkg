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

package path_test

import (
	"database/sql"
	"database/sql/driver"
	"testing"

	"github.com/corestoreio/csfw/config/path"
	"github.com/corestoreio/csfw/storage/csdb"
	"github.com/stretchr/testify/assert"
	"time"
)

var _ sql.Scanner = (*path.Route)(nil)
var _ driver.Valuer = (*path.Route)(nil)

func TestSQLType(t *testing.T) {
	//if false == testing.Short() {
	//	t.Skip("Only run in short test mode ...")
	//}

	dbCon := csdb.MustConnectTest()
	defer func() { assert.NoError(t, dbCon.Close()) }()

	assert.NoError(t, tableCollection.Init(dbCon.NewSession()))

	const testPath = `system/full_page_cache/varnish/backend_port`
	var rTestPath = path.Route(testPath)

	var insertVal = time.Now().Unix()
	ib := dbCon.NewSession().InsertInto(tableCollection.Name(tableIndexCoreConfigData))
	ib.Pair("path", rTestPath)
	ib.Pair("value", insertVal)

	res, err := ib.Exec()

	var ccds TableCoreConfigDataSlice

	rows, err := csdb.LoadSlice(dbCon.NewSession(), tableCollection, tableIndexCoreConfigData, &ccds)
	assert.NoError(t, err)
	assert.NotEmpty(t, rows)

	// test 1. write data 2. select this single row

	for _, ccd := range ccds {

		t.Logf("%#v", ccd)
	}
}
