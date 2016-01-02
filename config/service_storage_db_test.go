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

package config_test

import (
	"testing"
	"time"

	"fmt"
	"strings"

	"github.com/corestoreio/csfw/config"
	"github.com/corestoreio/csfw/storage/csdb"
	"github.com/corestoreio/csfw/store/scope"
	"github.com/corestoreio/csfw/util"
	"github.com/stretchr/testify/assert"
)

func TestDBStorageOneStmt(t *testing.T) {
	debugLogBuf.Reset()
	defer debugLogBuf.Reset()
	defer infoLogBuf.Reset()

	dbc := csdb.MustConnectTest()
	defer func() { assert.NoError(t, dbc.Close()) }()

	sdb := config.MustNewDBStorage(dbc.DB).Start()

	// Stop() would only be called under rare circumstances on a production system
	defer func() { assert.NoError(t, sdb.Stop()) }()

	tests := []struct {
		key       string
		value     interface{}
		wantNil   bool
		wantValue string
	}{
		{"stores/1/testDBStorage/secure/base_url", "http://corestore.io", false, "http://corestore.io"},
		{"stores/1/testDBStorage/log/active", 1, false, "1"},
		{"stores/99999/testDBStorage/log/clean", 19.999, false, "19.999"},
		{"stores/99999/testDBStorage/log/clean", 29.999, false, "29.999"},
		{"default/1/testDBStorage/catalog/purge", true, false, "true"},
		{"default/1/testDBStorage/catalog/clean", 0, false, "0"},
	}
	for _, test := range tests {
		sdb.Set(test.key, test.value)
		if test.wantNil {
			assert.Nil(t, sdb.Get(test.key), "Test: %v", test)
		} else {
			assert.Exactly(t, test.wantValue, sdb.Get(test.key), "Test: %v", test)
		}
	}

	assert.Exactly(t, 1, strings.Count(debugLogBuf.String(), `csdb.ResurrectStmt.stmt.Prepare SQL: "INSERT INTO`))
	assert.Exactly(t, 1, strings.Count(debugLogBuf.String(), "csdb.ResurrectStmt.stmt.Prepare SQL: \"SELECT `value` FROM"))

	for _, test := range tests {
		ak := util.StringSlice(sdb.AllKeys()) // trigger many queries with one statement
		assert.True(t, ak.Include(test.key), "Missing Key: %s", test.key)
	}
	assert.Exactly(t, 1, strings.Count(debugLogBuf.String(), fmt.Sprintf("CONCAT(scope,'%s',scope_id,'%s',path) AS `fqpath`", scope.PS, scope.PS)))
}

func TestDBStorageMultipleStmt(t *testing.T) {
	debugLogBuf.Reset()
	defer debugLogBuf.Reset() // contains only data from the debug level, info level will be dumped to os.Stdout
	defer infoLogBuf.Reset()

	if testing.Short() {
		t.Skip("Test skipped in short mode")
	}
	dbc := csdb.MustConnectTest()
	defer func() { assert.NoError(t, dbc.Close()) }()

	sdb := config.MustNewDBStorage(dbc.DB)
	sdb.All.Idle = time.Second * 1
	sdb.Read.Idle = time.Second * 1
	sdb.Write.Idle = time.Second * 1

	sdb.Start()

	tests := []struct {
		key       string
		value     interface{}
		wantValue string
	}{
		{"websites/10/testDBStorage/secure/base_url", "http://corestore.io", "http://corestore.io"},
		{"websites/10/testDBStorage/log/active", 1, "1"},
		{"websites/20/testDBStorage/log/clean", 19.999, "19.999"},
		{"websites/20/testDBStorage/product/shipping", 29.999, "29.999"},
		{"default/0/testDBStorage/checkout/multishipping", false, "false"},
	}
	for i, test := range tests {
		sdb.Set(test.key, test.value)
		assert.Exactly(t, test.wantValue, sdb.Get(test.key), "Test: %v", test)
		if i < 2 {
			// last two iterations reopen a new statement, not closing it and reusing it
			time.Sleep(time.Millisecond * 1500) // trigger ticker to close statements
		}
	}

	for i, test := range tests {
		ak := util.StringSlice(sdb.AllKeys())
		assert.True(t, ak.Include(test.key), "Missing Key: %s", test.key)
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
	assert.Exactly(t, 6, strings.Count(logStr, fmt.Sprintf("CONCAT(scope,'%s',scope_id,'%s',path) AS `fqpath`", scope.PS, scope.PS)))
}
