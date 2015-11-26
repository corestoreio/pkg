// Copyright 2015, Cyrill @ Schumacher.fm and the CoreStore contributors
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

	"github.com/corestoreio/csfw/config"
	"github.com/corestoreio/csfw/storage/csdb"
	"github.com/corestoreio/csfw/utils"
	"github.com/stretchr/testify/assert"
	"strings"
)

func TestDBStorageOneStmt(t *testing.T) {
	defer debugLogBuf.Reset() // contains only data from the debug level, info level will be dumped to os.Stdout

	dbc := csdb.MustConnectTest()
	defer func() { assert.NoError(t, dbc.Close()) }()

	sdb := config.NewDBStorage(dbc.DB)
	sdb.Read.Idle = time.Second * 2
	sdb.Write.Idle = time.Second * 3

	sdb.Start()
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

	assert.Exactly(t, 1, strings.Count(debugLogBuf.String(), "(`scope`,`scope_id`,`path`,`value`) VALUES"))
	assert.Exactly(t, 1, strings.Count(debugLogBuf.String(), "WHERE `scope`=? AND `scope_id`=? AND `path`=?"))

	ak := utils.StringSlice(sdb.AllKeys())
	for _, test := range tests {
		assert.True(t, ak.Include(test.key), "Missing Key: %s", test.key)
	}
	assert.Exactly(t, 1, strings.Count(debugLogBuf.String(), "CONCAT(scope,'%s',scope_id,'%s',path) AS `fqpath`"))
}
