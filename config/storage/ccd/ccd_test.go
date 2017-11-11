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

package ccd_test

import (
	"database/sql/driver"
	"fmt"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/corestoreio/cspkg/config"
	"github.com/corestoreio/cspkg/config/cfgpath"
	"github.com/corestoreio/cspkg/config/storage/ccd"
	"github.com/corestoreio/cspkg/store/scope"
	"github.com/corestoreio/cspkg/util/cstesting"
	"github.com/stretchr/testify/assert"
)

var _ config.Storager = (*ccd.DBStorage)(nil)

func TestDBStorageOneStmt(t *testing.T) {
	t.Parallel()

	dbc, dbMock := cstesting.MockDB(t)
	defer func() {
		dbMock.ExpectClose()

		assert.NoError(t, dbc.Close())

		if err := dbMock.ExpectationsWereMet(); err != nil {
			t.Error("there were unfulfilled expections", err)
		}
	}()

	sdb := ccd.MustNewDBStorage(dbc.DB).Start()

	// Stop() would only be called under rare circumstances on a production system
	defer func() { assert.NoError(t, sdb.Stop()) }()

	tests := []struct {
		key       cfgpath.Path
		value     interface{}
		wantNil   bool
		wantValue string
	}{
		{cfgpath.MustNewByParts("testDBStorage/secure/base_url").BindStore(1), "http://corestore.io", false, "http://corestore.io"},
		{cfgpath.MustNewByParts("testDBStorage/log/active").BindStore(2), 1, false, "1"},
		{cfgpath.MustNewByParts("testDBStorage/log/clean").BindStore(99999), 19.999, false, "19.999"},
		{cfgpath.MustNewByParts("testDBStorage/log/clean").BindStore(99999), 29.999, false, "29.999"},
		{cfgpath.MustNewByParts("testDBStorage/catalog/purge").Bind(scope.DefaultTypeID), true, false, "true"},
		{cfgpath.MustNewByParts("testDBStorage/catalog/clean").Bind(scope.DefaultTypeID), 0, false, "0"},
	}

	prepIns := dbMock.ExpectPrepare("INSERT INTO `[^`]+` \\(.+\\) VALUES \\(\\?,\\?,\\?,\\?\\) ON DUPLICATE KEY UPDATE `value`=\\?")
	for i, test := range tests {

		prepIns.ExpectExec().WithArgs(
			driver.Value(test.key.ScopeID.Type().StrType()),
			driver.Value(test.key.ScopeID.ID()),
			driver.Value(test.key.Bytes()),
			driver.Value(test.wantValue), // value
			driver.Value(test.wantValue), // value fall back on duplicate key
		).WillReturnResult(sqlmock.NewResult(0, 1))

		if err := sdb.Set(test.key, test.value); err != nil {
			t.Fatal("Index", i, " => ", err)
		}
	}

	// both prepared statements cannot run within one for loop :-(

	prepSel := dbMock.ExpectPrepare("SELECT `value` FROM `[^`]+` WHERE `scope`=\\? AND `scope_id`=\\? AND `path`=\\?")
	for i, test := range tests {

		prepSel.ExpectQuery().WithArgs(
			driver.Value(test.key.ScopeID.Type().StrType()),
			driver.Value(test.key.ScopeID.ID()),
			driver.Value(test.key.Bytes()),
		).WillReturnRows(sqlmock.NewRows([]string{"value"}).FromCSVString(test.wantValue))

		if test.wantNil {
			g, err := sdb.Get(test.key)
			assert.NoError(t, err, "Index %d", i)
			assert.Nil(t, g, "Index %d", i)
		} else {
			g, err := sdb.Get(test.key)
			assert.NoError(t, err, "Index %d", i)
			assert.Exactly(t, test.wantValue, g, "Index %d", i)
		}
	}

	//assert.Exactly(t, 1, strings.Count(debugLogBuf.String(), `csdb.ResurrectStmt.stmt.Prepare SQL: "INSERT INTO`))
	//assert.Exactly(t, 1, strings.Count(debugLogBuf.String(), "csdb.ResurrectStmt.stmt.Prepare SQL: \"SELECT `value` FROM"))

	mockRows := sqlmock.NewRows([]string{"scope", "scope_id", "path"})
	for _, test := range tests {
		mockRows.FromCSVString(fmt.Sprintf("%s,%d,%s", test.key.ScopeID.Type().StrType(), test.key.ScopeID.ID(), test.key.Chars))
	}
	prepAll := dbMock.ExpectPrepare("SELECT scope,scope_id,path FROM `[^`]+` ORDER BY scope,scope_id,path")
	prepAll.ExpectQuery().WillReturnRows(mockRows)

	allKeys, err := sdb.AllKeys()
	assert.NoError(t, err)

	for i, test := range tests {
		assert.True(t, allKeys.Contains(test.key), "Missing Key: %s\nIndex %d", test.key, i)
	}
	//assert.Exactly(t, 1, strings.Count(debugLogBuf.String(), `SELECT scope,scope_id,path FROM `))
}

var dbStorageMultiTests = []struct {
	key       cfgpath.Path
	value     interface{}
	wantValue string
}{
	{cfgpath.MustNewByParts("testDBStorage/secure/base_url").BindWebsite(10), "http://corestore.io", "http://corestore.io"},
	{cfgpath.MustNewByParts("testDBStorage/log/active").BindWebsite(10), 1, "1"},
	{cfgpath.MustNewByParts("testDBStorage/log/clean").BindWebsite(20), 19.999, "19.999"},
	{cfgpath.MustNewByParts("testDBStorage/product/shipping").BindWebsite(20), 29.999, "29.999"},
	{cfgpath.MustNewByParts("testDBStorage/checkout/multishipping"), false, "false"},
	{cfgpath.MustNewByParts("testDBStorage/shipping/rate").BindStore(321), 3.14159, "3.14159"},
}

func TestDBStorageMultipleStmt_Set(t *testing.T) {
	t.Parallel()
	dbc, dbMock := cstesting.MockDB(t)
	defer func() {
		dbMock.ExpectClose()

		assert.NoError(t, dbc.Close())

		if err := dbMock.ExpectationsWereMet(); err != nil {
			t.Error("there were unfulfilled expections", err)
		}
	}()

	sdb := ccd.MustNewDBStorage(dbc.DB)
	sdb.Write.Idle = time.Second * 1

	sdb.Start()

	var prepIns *sqlmock.ExpectedPrepare
	for i, test := range dbStorageMultiTests {
		if i < 3 {
			prepIns = dbMock.ExpectPrepare("INSERT INTO `[^`]+` \\(.+\\) VALUES \\(\\?,\\?,\\?,\\?\\) ON DUPLICATE KEY UPDATE `value`=\\?")
		}

		prepIns.ExpectExec().WithArgs(
			driver.Value(test.key.ScopeID.Type().StrType()),
			driver.Value(test.key.ScopeID.ID()),
			driver.Value(test.key.Bytes()),
			driver.Value(test.wantValue), // value
			driver.Value(test.wantValue), // value fall back on duplicate key
		).WillReturnResult(sqlmock.NewResult(0, 1))

		if err := sdb.Set(test.key, test.value); err != nil {
			t.Fatal("Index", i, "with error:", err)
		}

		if i < 2 {
			// last two iterations reopen a new statement, not closing it and reusing it
			time.Sleep(time.Millisecond * 1500) // trigger ticker to close statements
		}
	}

	assert.NoError(t, sdb.Stop())
}

func TestDBStorageMultipleStmt_Get(t *testing.T) {
	t.Parallel()
	dbc, dbMock := cstesting.MockDB(t)
	defer func() {
		dbMock.ExpectClose()

		assert.NoError(t, dbc.Close())

		if err := dbMock.ExpectationsWereMet(); err != nil {
			t.Error("there were unfulfilled expections", err)
		}
	}()

	sdb := ccd.MustNewDBStorage(dbc.DB)
	sdb.Read.Idle = time.Second * 1

	sdb.Start()

	// both prepared statements cannot run within one for loop :-(

	var prepSel *sqlmock.ExpectedPrepare
	for i, test := range dbStorageMultiTests {
		if i < 3 {
			prepSel = dbMock.ExpectPrepare("SELECT `value` FROM `[^`]+` WHERE `scope`=\\? AND `scope_id`=\\? AND `path`=\\?")
		}

		prepSel.ExpectQuery().WithArgs(
			driver.Value(test.key.ScopeID.Type().StrType()),
			driver.Value(test.key.ScopeID.ID()),
			driver.Value(test.key.Bytes()),
		).WillReturnRows(sqlmock.NewRows([]string{"value"}).FromCSVString(test.wantValue))

		g, err := sdb.Get(test.key)
		assert.NoError(t, err, "Index %d", i)
		assert.Exactly(t, test.wantValue, g, "Index %d", i)

		if i < 2 {
			// last two iterations reopen a new statement, not closing it and reusing it
			time.Sleep(time.Millisecond * 1500) // trigger ticker to close statements
		}
	}

	assert.NoError(t, sdb.Stop())
}

func TestDBStorageMultipleStmt_All(t *testing.T) {
	t.Parallel()

	dbc, dbMock := cstesting.MockDB(t)
	defer func() {
		dbMock.ExpectClose()

		assert.NoError(t, dbc.Close())

		if err := dbMock.ExpectationsWereMet(); err != nil {
			t.Error("there were unfulfilled expections", err)
		}
	}()

	sdb := ccd.MustNewDBStorage(dbc.DB)
	sdb.All.Idle = time.Second * 1

	sdb.Start()

	var prepAll *sqlmock.ExpectedPrepare
	for i, test := range dbStorageMultiTests {

		if i < 3 {
			prepAll = dbMock.ExpectPrepare("SELECT scope,scope_id,path FROM `[^`]+` ORDER BY scope,scope_id,path")
		}

		mockRows := sqlmock.NewRows([]string{"scope", "scope_id", "path"})
		for _, test := range dbStorageMultiTests {
			mockRows.FromCSVString(fmt.Sprintf("%s,%d,%s", test.key.ScopeID.Type().StrType(), test.key.ScopeID.ID(), test.key.Chars))
		}
		prepAll.ExpectQuery().WillReturnRows(mockRows)

		allKeys, err := sdb.AllKeys()
		if err != nil {
			t.Fatal("Index", i, "with error AllKeys:", err)
		}

		assert.True(t, allKeys.Contains(test.key), "Missing Key: %s", test.key)

		if i < 2 {
			time.Sleep(time.Millisecond * 1500) // trigger ticker to close statements
		}
	}
	assert.NoError(t, sdb.Stop())

}
