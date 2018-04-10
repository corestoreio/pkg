// Copyright 2015-present, Cyrill @ Schumacher.fm and the CoreStore contributors
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
	"github.com/corestoreio/errors"
	"github.com/corestoreio/pkg/config"
	"github.com/corestoreio/pkg/config/storage/ccd"
	"github.com/corestoreio/pkg/sql/dmltest"
	"github.com/corestoreio/pkg/sync/bgwork"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

var _ config.Storager = (*ccd.DBStorage)(nil)

func TestMustNewDBStorage_Panic(t *testing.T) {
	t.Parallel()
	defer func() {
		if r := recover(); r != nil {
			if err, ok := r.(error); ok {
				assert.True(t, errors.NotFound.Match(err), "%+v", err)
			} else {
				t.Errorf("Panic should contain an error but got:\n%+v", r)
			}
		} else {
			t.Error("Expecting a panic but got nothing")
		}
	}()
	_ = ccd.MustNewDBStorage(ccd.NewTableCollection(nil), ccd.Options{
		TableName: "non-existent",
	})
}

func TestDBStorage_AllKeys_Mocked(t *testing.T) {
	t.Parallel()

	dbc, dbMock := dmltest.MockDB(t)
	defer dmltest.MockClose(t, dbc, dbMock)

	t.Run("table not found", func(t *testing.T) {
		dbs, err := ccd.NewDBStorage(ccd.NewTableCollection(dbc.DB), ccd.Options{
			TableName: "non-existent",
		})
		assert.Nil(t, dbs)
		assert.True(t, errors.NotFound.Match(err), "%+v", err)
	})

	t.Run("no leaking goroutines", func(t *testing.T) {
		// TODO use package leak test
		dbs, err := ccd.NewDBStorage(ccd.NewTableCollection(dbc.DB), ccd.Options{})
		require.NoError(t, err)
		assert.NoError(t, dbs.Close())
	})

	t.Run("return all keys, no waiting", func(t *testing.T) {
		prepQry := dbMock.ExpectPrepare(dmltest.SQLMockQuoteMeta("SELECT `scope`, `scope_id`, `path` FROM `core_config_data` AS `main_table` ORDER BY `scope`, `scope_id`, `path`")).ExpectQuery()
		rows, err := dmltest.MockRows(dmltest.WithFile("testdata", "core_config_data.csv"))
		require.NoError(t, err)
		prepQry.WithArgs().WillReturnRows(rows)

		dbs, err := ccd.NewDBStorage(ccd.NewTableCollection(dbc.DB), ccd.Options{})
		require.NoError(t, err)
		defer dmltest.Close(t, dbs)

		scps, paths, err := dbs.AllKeys()
		require.NoError(t, err)
		assert.Exactly(t, []string{"cms/wysiwyg/enabled", "general/region/display_all", "general/region/state_required", "general/region/state_required", "web/url/redirect_to_base", "web/unsecure/base_url", "web/unsecure/base_url", "web/unsecure/base_link_url", "web/unsecure/base_skin_url", "web/unsecure/base_media_url"},
			paths)
		assert.Exactly(t, "Type(Default) ID(0); Type(Store) ID(4); Type(Default) ID(0); Type(Store) ID(2); Type(Default) ID(0); Type(Default) ID(0); Type(Website) ID(1); Type(Default) ID(0); Type(Website) ID(44); Type(Default) ID(0)",
			scps.String())
	})

	t.Run("return all keys, waiting and reprepare", func(t *testing.T) {

		dbs, err := ccd.NewDBStorage(ccd.NewTableCollection(dbc.DB), ccd.Options{
			IdleAllKeys: time.Millisecond * 5,
		})
		require.NoError(t, err)
		defer dmltest.Close(t, dbs)

		for i := 0; i < 4; i++ {
			prepQry := dbMock.ExpectPrepare(dmltest.SQLMockQuoteMeta("SELECT `scope`, `scope_id`, `path` FROM `core_config_data` AS `main_table` ORDER BY `scope`, `scope_id`, `path`")).ExpectQuery()
			rows, err := dmltest.MockRows(dmltest.WithFile("testdata", "core_config_data.csv"))
			require.NoError(t, err)
			prepQry.WithArgs().WillReturnRows(rows)

			scps, paths, err := dbs.AllKeys()
			require.NoError(t, err)
			assert.Exactly(t, []string{"cms/wysiwyg/enabled", "general/region/display_all", "general/region/state_required", "general/region/state_required", "web/url/redirect_to_base", "web/unsecure/base_url", "web/unsecure/base_url", "web/unsecure/base_link_url", "web/unsecure/base_skin_url", "web/unsecure/base_media_url"},
				paths)
			assert.Exactly(t, "Type(Default) ID(0); Type(Store) ID(4); Type(Default) ID(0); Type(Store) ID(2); Type(Default) ID(0); Type(Default) ID(0); Type(Website) ID(1); Type(Default) ID(0); Type(Website) ID(44); Type(Default) ID(0)",
				scps.String())

			time.Sleep(time.Millisecond * 8)
		}
	})

	//tests := []struct {
	//	key       config.Path
	//	value     interface{}
	//	wantNil   bool
	//	wantValue string
	//}{
	//	{config.MustMakePath("testDBStorage/secure/base_url").BindStore(1), "http://corestore.io", false, "http://corestore.io"},
	//	{config.MustMakePath("testDBStorage/log/active").BindStore(2), 1, false, "1"},
	//	{config.MustMakePath("testDBStorage/log/clean").BindStore(99999), 19.999, false, "19.999"},
	//	{config.MustMakePath("testDBStorage/log/clean").BindStore(99999), 29.999, false, "29.999"},
	//	{config.MustMakePath("testDBStorage/catalog/purge").Bind(scope.DefaultTypeID), true, false, "true"},
	//	{config.MustMakePath("testDBStorage/catalog/clean").Bind(scope.DefaultTypeID), 0, false, "0"},
	//}
	//
	//prepIns := dbMock.ExpectPrepare("INSERT INTO `[^`]+` \\(.+\\) VALUES \\(\\?,\\?,\\?,\\?\\) ON DUPLICATE KEY UPDATE `value`=\\?")
	//for i, test := range tests {
	//
	//	prepIns.ExpectExec().WithArgs(
	//		driver.Value(test.key.ScopeID.Type().StrType()),
	//		driver.Value(test.key.ScopeID.ID()),
	//		driver.Value(test.key.Bytes()),
	//		driver.Value(test.wantValue), // value
	//		driver.Value(test.wantValue), // value fall back on duplicate key
	//	).WillReturnResult(sqlmock.NewResult(0, 1))
	//
	//	if err := sdb.Set(test.key, test.value); err != nil {
	//		t.Fatal("Index", i, " => ", err)
	//	}
	//}
	//
	//// both prepared statements cannot run within one for loop :-(
	//
	//prepSel := dbMock.ExpectPrepare("SELECT `value` FROM `[^`]+` WHERE `scope`=\\? AND `scope_id`=\\? AND `path`=\\?")
	//for i, test := range tests {
	//
	//	prepSel.ExpectQuery().WithArgs(
	//		driver.Value(test.key.ScopeID.Type().StrType()),
	//		driver.Value(test.key.ScopeID.ID()),
	//		driver.Value(test.key.Bytes()),
	//	).WillReturnRows(sqlmock.NewRows([]string{"value"}).FromCSVString(test.wantValue))
	//
	//	if test.wantNil {
	//		g, err := sdb.Value(test.key)
	//		assert.NoError(t, err, "Index %d", i)
	//		assert.Nil(t, g, "Index %d", i)
	//	} else {
	//		g, err := sdb.Value(test.key)
	//		assert.NoError(t, err, "Index %d", i)
	//		assert.Exactly(t, test.wantValue, g, "Index %d", i)
	//	}
	//}
}

func TestDBStorage_AllKeys_Integration(t *testing.T) {
	t.Parallel()

	dbc := dmltest.MustConnectDB(t)
	defer dmltest.Close(t, dbc)

	dbs, err := ccd.NewDBStorage(ccd.NewTableCollection(dbc.DB), ccd.Options{
		IdleAllKeys: time.Millisecond * 2,
	})
	require.NoError(t, err)
	defer dmltest.Close(t, dbs)

	bgwork.Wait(10, func(idx int) {
		scps, paths, err := dbs.AllKeys()
		require.NoError(t, err)
		assert.Exactly(t, len(scps), len(paths))
		assert.True(t, len(paths) > 5, "path string slice should contain at least 5 items")
		time.Sleep(time.Millisecond * time.Duration(idx))
	})
}

var dbStorageMultiTests = []struct {
	key       config.Path
	value     interface{}
	wantValue string
}{
	{config.MustMakePath("testDBStorage/secure/base_url").BindWebsite(10), "http://corestore.io", "http://corestore.io"},
	{config.MustMakePath("testDBStorage/log/active").BindWebsite(10), 1, "1"},
	{config.MustMakePath("testDBStorage/log/clean").BindWebsite(20), 19.999, "19.999"},
	{config.MustMakePath("testDBStorage/product/shipping").BindWebsite(20), 29.999, "29.999"},
	{config.MustMakePath("testDBStorage/checkout/multishipping"), false, "false"},
	{config.MustMakePath("testDBStorage/shipping/rate").BindStore(321), 3.14159, "3.14159"},
}

//func TestDBStorageMultipleStmt_Set(t *testing.T) {
//	t.Parallel()
//	dbc, dbMock := cstesting.MockDB(t)
//	defer func() {
//		dbMock.ExpectClose()
//
//		assert.NoError(t, dbc.Close())
//
//		if err := dbMock.ExpectationsWereMet(); err != nil {
//			t.Error("there were unfulfilled expections", err)
//		}
//	}()
//
//	sdb := ccd.MustNewDBStorage(dbc.DB)
//	sdb.Write.Idle = time.Second * 1
//
//	sdb.Start()
//
//	var prepIns *sqlmock.ExpectedPrepare
//	for i, test := range dbStorageMultiTests {
//		if i < 3 {
//			prepIns = dbMock.ExpectPrepare("INSERT INTO `[^`]+` \\(.+\\) VALUES \\(\\?,\\?,\\?,\\?\\) ON DUPLICATE KEY UPDATE `value`=\\?")
//		}
//
//		prepIns.ExpectExec().WithArgs(
//			driver.Value(test.key.ScopeID.Type().StrType()),
//			driver.Value(test.key.ScopeID.ID()),
//			driver.Value(test.key.Bytes()),
//			driver.Value(test.wantValue), // value
//			driver.Value(test.wantValue), // value fall back on duplicate key
//		).WillReturnResult(sqlmock.NewResult(0, 1))
//
//		if err := sdb.Set(test.key, test.value); err != nil {
//			t.Fatal("Index", i, "with error:", err)
//		}
//
//		if i < 2 {
//			// last two iterations reopen a new statement, not closing it and reusing it
//			time.Sleep(time.Millisecond * 1500) // trigger ticker to close statements
//		}
//	}
//
//	assert.NoError(t, sdb.Stop())
//}

//func TestDBStorageMultipleStmt_Get(t *testing.T) {
//	t.Parallel()
//	dbc, dbMock := cstesting.MockDB(t)
//	defer func() {
//		dbMock.ExpectClose()
//
//		assert.NoError(t, dbc.Close())
//
//		if err := dbMock.ExpectationsWereMet(); err != nil {
//			t.Error("there were unfulfilled expections", err)
//		}
//	}()
//
//	sdb := ccd.MustNewDBStorage(dbc.DB)
//	sdb.Read.Idle = time.Second * 1
//
//	sdb.Start()
//
//	// both prepared statements cannot run within one for loop :-(
//
//	var prepSel *sqlmock.ExpectedPrepare
//	for i, test := range dbStorageMultiTests {
//		if i < 3 {
//			prepSel = dbMock.ExpectPrepare("SELECT `value` FROM `[^`]+` WHERE `scope`=\\? AND `scope_id`=\\? AND `path`=\\?")
//		}
//
//		prepSel.ExpectQuery().WithArgs(
//			driver.Value(test.key.ScopeID.Type().StrType()),
//			driver.Value(test.key.ScopeID.ID()),
//			driver.Value(test.key.Bytes()),
//		).WillReturnRows(sqlmock.NewRows([]string{"value"}).FromCSVString(test.wantValue))
//
//		g, err := sdb.Value(test.key)
//		assert.NoError(t, err, "Index %d", i)
//		assert.Exactly(t, test.wantValue, g, "Index %d", i)
//
//		if i < 2 {
//			// last two iterations reopen a new statement, not closing it and reusing it
//			time.Sleep(time.Millisecond * 1500) // trigger ticker to close statements
//		}
//	}
//
//	assert.NoError(t, sdb.Stop())
//}

func TestDBStorageMultipleStmt_All(t *testing.T) {
	t.Parallel()
	//
	//dbc, dbMock := cstesting.MockDB(t)
	//defer func() {
	//	dbMock.ExpectClose()
	//
	//	assert.NoError(t, dbc.Close())
	//
	//	if err := dbMock.ExpectationsWereMet(); err != nil {
	//		t.Error("there were unfulfilled expections", err)
	//	}
	//}()
	//
	//sdb := ccd.MustNewDBStorage(dbc.DB)
	//sdb.All.Idle = time.Second * 1
	//
	//sdb.Start()
	//
	//var prepAll *sqlmock.ExpectedPrepare
	//for i, test := range dbStorageMultiTests {
	//
	//	if i < 3 {
	//		prepAll = dbMock.ExpectPrepare("SELECT scope,scope_id,path FROM `[^`]+` ORDER BY scope,scope_id,path")
	//	}
	//
	//	mockRows := sqlmock.NewRows([]string{"scope", "scope_id", "path"})
	//	for _, test := range dbStorageMultiTests {
	//		mockRows.FromCSVString(fmt.Sprintf("%s,%d,%s", test.key.ScopeID.Type().StrType(), test.key.ScopeID.ID(), test.key.Chars))
	//	}
	//	prepAll.ExpectQuery().WillReturnRows(mockRows)
	//
	//	allKeys, err := sdb.AllKeys()
	//	if err != nil {
	//		t.Fatal("Index", i, "with error AllKeys:", err)
	//	}
	//
	//	assert.True(t, allKeys.Contains(test.key), "Missing Key: %s", test.key)
	//
	//	if i < 2 {
	//		time.Sleep(time.Millisecond * 1500) // trigger ticker to close statements
	//	}
	//}
	//assert.NoError(t, sdb.Stop())

}

// TestIntegrationSQLType is not a real test for the type Route
//func TestIntegrationSQLType(t *testing.T) {
//
//	dbCon, dbMock := dmltest.MockDB(t)
//	defer dmltest.MockClose(t, dbCon, dbMock)
//
//	var testPath = `system/full_page_cache/varnish/` + strs.RandAlnum(5)
//	//var insPath = cfgpath.MakeRoute(testPath)
//	var insVal = time.Now().Unix()
//
//	dbMock.ExpectQuery(dmltest.SQLMockQuoteMeta("SELECT `config_id`, `scope`, `scope_id`, `path`, `value` FROM `core_config_data` AS `main_table`")).
//		WithArgs(testPath).
//		WillReturnRows(
//			sqlmock.NewRows([]string{"config_id", "scope", "scope_id", "path", "value"}).
//				AddRow(1, "default", 0, testPath, fmt.Sprintf("%d", insVal)),
//		)
//
//	var ccd ccd.TableCoreConfigData
//	tbl := tableCollection.MustTable(TableNameCoreConfigData)
//	rc, err := tbl.SelectAll().WithDB(dbCon.DB).WithArgs().String(testPath).Load(context.TODO(), &ccd)
//	require.NoError(t, err)
//	assert.Exactly(t, uint64(1), rc)
//
//	assert.Exactly(t, testPath, ccd.Path.String())
//	haveI64, err := strconv.ParseInt(ccd.Value.String, 10, 64)
//	assert.NoError(t, err)
//	assert.Exactly(t, insVal, haveI64)
//}
