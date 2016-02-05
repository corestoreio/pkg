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
	"testing"
	"time"

	"strings"

	"github.com/corestoreio/csfw/config"
	"github.com/corestoreio/csfw/config/ccd"
	"github.com/corestoreio/csfw/config/path"
	"github.com/corestoreio/csfw/storage/csdb"
	"github.com/corestoreio/csfw/store/scope"
	"github.com/stretchr/testify/assert"
)

var _ config.Storager = (*ccd.DBStorage)(nil)

func TestDBStorageOneStmt(t *testing.T) {
	debugLogBuf.Reset()
	defer debugLogBuf.Reset()
	defer infoLogBuf.Reset()
	if _, err := csdb.GetDSNTest(); err == csdb.ErrDSNTestNotFound {
		t.Skip(err)
	}

	dbc := csdb.MustConnectTest()
	defer func() { assert.NoError(t, dbc.Close()) }()

	sdb := ccd.MustNewDBStorage(dbc.DB).Start()

	// Stop() would only be called under rare circumstances on a production system
	defer func() { assert.NoError(t, sdb.Stop()) }()

	tests := []struct {
		key       path.Path
		value     interface{}
		wantNil   bool
		wantValue string
	}{
		{path.MustNewByParts("testDBStorage/secure/base_url").Bind(scope.StoreID, 1), "http://corestore.io", false, "http://corestore.io"},
		{path.MustNewByParts("testDBStorage/log/active").Bind(scope.StoreID, 1), 1, false, "1"},
		{path.MustNewByParts("testDBStorage/log/clean").Bind(scope.StoreID, 99999), 19.999, false, "19.999"},
		{path.MustNewByParts("testDBStorage/log/clean").Bind(scope.StoreID, 99999), 29.999, false, "29.999"},
		{path.MustNewByParts("testDBStorage/catalog/purge").Bind(scope.DefaultID, 1), true, false, "true"},
		{path.MustNewByParts("testDBStorage/catalog/clean").Bind(scope.DefaultID, 1), 0, false, "0"},
	}
	for i, test := range tests {
		sdb.Set(test.key, test.value)
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

	assert.Exactly(t, 1, strings.Count(debugLogBuf.String(), `csdb.ResurrectStmt.stmt.Prepare SQL: "INSERT INTO`))
	assert.Exactly(t, 1, strings.Count(debugLogBuf.String(), "csdb.ResurrectStmt.stmt.Prepare SQL: \"SELECT `value` FROM"))

	for i, test := range tests {
		allKeys, err := sdb.AllKeys()
		assert.NoError(t, err, "Index %d", i)
		assert.True(t, allKeys.Contains(test.key), "Missing Key: %s\nIndex %d", test.key, i)
	}

	assert.Exactly(t, 1, strings.Count(debugLogBuf.String(), `SELECT scope,scope_id,path FROM `))
}

func TestDBStorageMultipleStmt(t *testing.T) {
	debugLogBuf.Reset()
	defer debugLogBuf.Reset() // contains only data from the debug level, info level will be dumped to os.Stdout
	defer infoLogBuf.Reset()
	if _, err := csdb.GetDSNTest(); err == csdb.ErrDSNTestNotFound {
		t.Skip(err)
	}

	if testing.Short() {
		t.Skip("Test skipped in short mode")
	}
	dbc := csdb.MustConnectTest()
	defer func() { assert.NoError(t, dbc.Close()) }()

	sdb := ccd.MustNewDBStorage(dbc.DB)
	sdb.All.Idle = time.Second * 1
	sdb.Read.Idle = time.Second * 1
	sdb.Write.Idle = time.Second * 1

	sdb.Start()

	tests := []struct {
		key       path.Path
		value     interface{}
		wantValue string
	}{
		{path.MustNewByParts("testDBStorage/secure/base_url").Bind(scope.WebsiteID, 10), "http://corestore.io", "http://corestore.io"},
		{path.MustNewByParts("testDBStorage/log/active").Bind(scope.WebsiteID, 10), 1, "1"},
		{path.MustNewByParts("testDBStorage/log/clean").Bind(scope.WebsiteID, 20), 19.999, "19.999"},
		{path.MustNewByParts("testDBStorage/product/shipping").Bind(scope.WebsiteID, 20), 29.999, "29.999"},
		{path.MustNewByParts("testDBStorage/checkout/multishipping"), false, "false"},
	}
	for i, test := range tests {
		assert.NoError(t, sdb.Set(test.key, test.value), "Index %d", i)
		g, err := sdb.Get(test.key)
		assert.NoError(t, err, "Index %d", i)
		assert.Exactly(t, test.wantValue, g, "Index %d", i)
		if i < 2 {
			// last two iterations reopen a new statement, not closing it and reusing it
			time.Sleep(time.Millisecond * 1500) // trigger ticker to close statements
		}
	}

	for i, test := range tests {
		allKeys, err := sdb.AllKeys()
		assert.NoError(t, err, "Index %d", i)
		assert.True(t, allKeys.Contains(test.key), "Missing Key: %s", test.key)
		if i < 2 {
			time.Sleep(time.Millisecond * 1500) // trigger ticker to close statements
		}
	}
	assert.NoError(t, sdb.Stop())

	logStr := debugLogBuf.String()
	assert.Exactly(t, 3, strings.Count(logStr, `csdb.ResurrectStmt.stmt.Prepare SQL: "INSERT INTO`))
	assert.Exactly(t, 3, strings.Count(logStr, "csdb.ResurrectStmt.stmt.Prepare SQL: \"SELECT `value` FROM"))

	assert.Exactly(t, 4, strings.Count(logStr, `csdb.ResurrectStmt.stmt.Close SQL: "INSERT INTO`), "\n%s\n", logStr)
	assert.Exactly(t, 4, strings.Count(logStr, "csdb.ResurrectStmt.stmt.Close SQL: \"SELECT `value` FROM"))

	//println("\n", logStr, "\n")

	// 6 is: open close for iteration 0+1, open in iteration 2 and close in iteration 4
	assert.Exactly(t, 6, strings.Count(logStr, `SELECT scope,scope_id,path FROM `))
}
