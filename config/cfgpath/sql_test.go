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

package cfgpath_test

import (
	"database/sql"
	"database/sql/driver"
	"strconv"
	"testing"
	"time"

	"github.com/corestoreio/csfw/config/cfgpath"
	"github.com/corestoreio/csfw/storage/dbr"
	"github.com/corestoreio/csfw/util"
	"github.com/corestoreio/csfw/util/cstesting"
	"github.com/stretchr/testify/assert"
)

var _ sql.Scanner = (*cfgpath.Route)(nil)
var _ driver.Valuer = (*cfgpath.Route)(nil)

func TestIntegrationSQLType(t *testing.T) {

	dbCon, version := cstesting.MustConnectDB()
	if version < 0 {
		t.Skip("Environment DB DSN not found")
	}
	defer func() { assert.NoError(t, dbCon.Close()) }()

	var testPath = `system/full_page_cache/varnish/` + util.RandAlnum(5)
	var insPath = cfgpath.NewRoute(testPath)
	var insVal = time.Now().Unix()

	// just for testing! TODO refactor test
	stmt, err := dbCon.DB.Prepare("INSERT INTO `" + tableCollection.Name(tableIndexCoreConfigData) + "` (path,value) values (?,?)")
	if false == assert.NoError(t, err) {
		t.Fatal("Stopping ...")
	}

	// yay! writing bytes instead of boring slow strings!
	res, err := stmt.Exec(insPath.Bytes(), insVal)
	if false == assert.NoError(t, err) {
		t.Fatal("Stopping ...")
	}

	id, err := res.LastInsertId()
	assert.NoError(t, err)
	assert.NotEmpty(t, id)

	var ccds TableCoreConfigDataSlice
	tbl := tableCollection.MustTable(tableIndexCoreConfigData)
	rows, err := tbl.LoadSlice(dbCon.NewSession(), &ccds, func(sb *dbr.Select) *dbr.Select {
		sb.Where(dbr.Condition("config_id=?", id))
		return sb
	})

	assert.NoError(t, err)
	assert.NotEmpty(t, rows)

	if false == assert.Exactly(t, int(1), rows) {
		t.Fatal("No rows loaded from the database!")
	}

	assert.Exactly(t, testPath, ccds[0].Path.String())
	haveI64, err := strconv.ParseInt(ccds[0].Value.String, 10, 64)
	assert.NoError(t, err)
	assert.Exactly(t, insVal, haveI64)
}
