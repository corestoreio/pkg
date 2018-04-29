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
	"fmt"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/corestoreio/errors"
	"github.com/corestoreio/pkg/config"
	"github.com/corestoreio/pkg/config/storage/ccd"
	"github.com/corestoreio/pkg/sql/dmltest"
	"github.com/corestoreio/pkg/store/scope"
	"github.com/corestoreio/pkg/sync/bgwork"
	"github.com/fortytw2/leaktest"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
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
		TableName:            "non-existent",
		SkipSchemaValidation: true,
	})
}

func TestDBStorage_AllKeys_Mocked(t *testing.T) {
	defer leaktest.CheckTimeout(t, time.Second)()

	dbc, dbMock := dmltest.MockDB(t)
	defer dmltest.MockClose(t, dbc, dbMock)

	t.Run("table not found", func(t *testing.T) {
		dbs, err := ccd.NewDBStorage(ccd.NewTableCollection(dbc.DB), ccd.Options{
			TableName:            "non-existent",
			SkipSchemaValidation: true,
		})
		assert.Nil(t, dbs)
		assert.True(t, errors.NotFound.Match(err), "%+v", err)
	})

	t.Run("no leaking goroutines", func(t *testing.T) {
		dbs, err := ccd.NewDBStorage(ccd.NewTableCollection(dbc.DB), ccd.Options{
			SkipSchemaValidation: true,
		})
		require.NoError(t, err)
		assert.NoError(t, dbs.Close())
	})

	t.Run("return all keys, no waiting", func(t *testing.T) {
		prepQry := dbMock.ExpectPrepare(dmltest.SQLMockQuoteMeta("SELECT `scope`, `scope_id`, `path` FROM `core_config_data` AS `main_table` ORDER BY `scope`, `scope_id`, `path`")).ExpectQuery()
		rows, err := dmltest.MockRows(dmltest.WithFile("testdata", "core_config_data.csv"))
		require.NoError(t, err)
		prepQry.WithArgs().WillReturnRows(rows)

		dbs, err := ccd.NewDBStorage(ccd.NewTableCollection(dbc.DB), ccd.Options{
			SkipSchemaValidation: true,
		})
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
			IdleAllKeys:          time.Millisecond * 5,
			SkipSchemaValidation: true,
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
}

func TestDBStorage_AllKeys_Integration(t *testing.T) {
	defer leaktest.Check(t)()

	dbc := dmltest.MustConnectDB(t)
	defer dmltest.Close(t, dbc)

	dbs, err := ccd.NewDBStorage(ccd.NewTableCollection(dbc.DB), ccd.Options{
		IdleAllKeys:          time.Millisecond * 2,
		SkipSchemaValidation: false,
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
	path    string
	scopeID scope.TypeID
	value   []byte
}{
	{"testDBStorage/secure/base_url", scope.Website.Pack(10), []byte("http://corestore.io")},
	{"testDBStorage/log/active", scope.Store.Pack(9), []byte("https://crestre.i")},
	{"testDBStorage/checkout/multishipping", scope.DefaultTypeID, []byte("false")},
}

func TestDBStorage_Value(t *testing.T) {
	defer leaktest.CheckTimeout(t, time.Second)()

	testBody := func(t *testing.T, dbs *ccd.DBStorage, dbMock sqlmock.Sqlmock, sleep time.Duration) {

		prepSel := dbMock.ExpectPrepare(dmltest.SQLMockQuoteMeta("SELECT `value` FROM `core_config_data` AS `main_table` WHERE (`scope` = ?) AND (`scope_id` = ?) AND (`path` = ?)"))
		for _, test := range dbStorageMultiTests {
			scp, sID := test.scopeID.Unpack()
			prepSel.ExpectQuery().WithArgs(scp.StrType(), sID, test.path).WillReturnRows(sqlmock.NewRows([]string{"value"}))

			haveVal, haveOK, haveErr := dbs.Value(test.scopeID, test.path)
			require.NoError(t, haveErr)
			require.False(t, haveOK, "%s Value with path %q should NOT be found", test.scopeID, test.path)
			assert.Exactly(t, []byte(nil), haveVal)
		}

		if sleep > 0 {
			time.Sleep(sleep)
			prepSel = dbMock.ExpectPrepare(dmltest.SQLMockQuoteMeta("SELECT `value` FROM `core_config_data` AS `main_table` WHERE (`scope` = ?) AND (`scope_id` = ?) AND (`path` = ?)"))
		}

		for _, test := range dbStorageMultiTests {
			scp, sID := test.scopeID.Unpack()
			prepSel.ExpectQuery().WithArgs(scp.StrType(), sID, test.path).WillReturnRows(sqlmock.NewRows([]string{"value"}).AddRow(test.value))

			haveVal, haveOK, haveErr := dbs.Value(test.scopeID, test.path)
			require.NoError(t, haveErr)
			require.True(t, haveOK, "%s Value with path %q should be found", test.scopeID, test.path)
			assert.Exactly(t, test.value, haveVal)
		}
	}

	t.Run("no waiting", func(t *testing.T) {
		dbc, dbMock := dmltest.MockDB(t)
		defer dmltest.MockClose(t, dbc, dbMock)
		dbMock.MatchExpectationsInOrder(false)

		dbs, err := ccd.NewDBStorage(ccd.NewTableCollection(dbc.DB), ccd.Options{
			SkipSchemaValidation: true,
		})
		require.NoError(t, err)
		defer dmltest.Close(t, dbs)
		testBody(t, dbs, dbMock, 0)
	})

	t.Run("wait and restart", func(t *testing.T) {
		dbc, dbMock := dmltest.MockDB(t)
		defer dmltest.MockClose(t, dbc, dbMock)
		dbMock.MatchExpectationsInOrder(false)

		dbs, err := ccd.NewDBStorage(ccd.NewTableCollection(dbc.DB), ccd.Options{
			IdleRead:             time.Millisecond * 50,
			IdleWrite:            time.Millisecond * 50,
			SkipSchemaValidation: true,
		})
		require.NoError(t, err)
		defer dmltest.Close(t, dbs)

		testBody(t, dbs, dbMock, time.Millisecond*100)

		val, set, all := dbs.Statistics()
		assert.Exactly(t,
			"read ccd.stats{Open:0x2, Close:0x1} write ccd.stats{Open:0x0, Close:0x0} all ccd.stats{Open:0x0, Close:0x0}",
			fmt.Sprintf("read %#v write %#v all %#v", val, set, all),
		)
	})

	t.Run("query context timeout", func(t *testing.T) {
		dbc, dbMock := dmltest.MockDB(t)
		defer dmltest.MockClose(t, dbc, dbMock)
		dbMock.MatchExpectationsInOrder(false)

		dbs, err := ccd.NewDBStorage(ccd.NewTableCollection(dbc.DB), ccd.Options{
			ContextTimeoutRead:   time.Millisecond * 50,
			SkipSchemaValidation: true,
		})
		require.NoError(t, err)
		defer dmltest.Close(t, dbs)

		prepSel := dbMock.ExpectPrepare(dmltest.SQLMockQuoteMeta("SELECT `value` FROM `core_config_data` AS `main_table` WHERE (`scope` = ?) AND (`scope_id` = ?) AND (`path` = ?)"))
		for _, test := range dbStorageMultiTests {
			scp, sID := test.scopeID.Unpack()
			prepSel.ExpectQuery().WithArgs(scp.StrType(), sID, test.path).WillDelayFor(time.Millisecond * 110).WillReturnRows(sqlmock.NewRows([]string{"value"}))

			haveVal, haveOK, haveErr := dbs.Value(test.scopeID, test.path)
			assert.Nil(t, haveVal)
			assert.False(t, haveOK)
			causeErr := errors.Cause(haveErr)
			require.EqualError(t, causeErr, "canceling query due to user request")
		}

	})
}

func TestDBStorage_Set(t *testing.T) {
	defer leaktest.CheckTimeout(t, time.Second)()

	testBody := func(t *testing.T, dbs *ccd.DBStorage, dbMock sqlmock.Sqlmock, sleep time.Duration) {

		prepIns := dbMock.ExpectPrepare(dmltest.SQLMockQuoteMeta("INSERT INTO `core_config_data` (`scope`,`scope_id`,`path`,`value`) VALUES (?,?,?,?) ON DUPLICATE KEY UPDATE `value`=VALUES(`value`)"))

		for i, test := range dbStorageMultiTests {
			j := int64(i + 1)

			if sleep > 0 && i > 0 {
				prepIns = dbMock.ExpectPrepare(dmltest.SQLMockQuoteMeta("INSERT INTO `core_config_data` (`scope`,`scope_id`,`path`,`value`) VALUES (?,?,?,?) ON DUPLICATE KEY UPDATE `value`=VALUES(`value`)"))
			}

			prepIns.ExpectExec().
				WithArgs(test.scopeID, test.path, test.value).
				WillReturnResult(sqlmock.NewResult(j, 0))
			require.NoError(t, dbs.Set(test.scopeID, test.path, test.value))

			if sleep > 0 {
				time.Sleep(sleep)
			}
		}
	}

	t.Run("no waiting", func(t *testing.T) {
		dbc, dbMock := dmltest.MockDB(t)
		defer dmltest.MockClose(t, dbc, dbMock)
		dbMock.MatchExpectationsInOrder(false)

		dbs, err := ccd.NewDBStorage(ccd.NewTableCollection(dbc.DB), ccd.Options{
			SkipSchemaValidation: true,
		})
		require.NoError(t, err)
		defer dmltest.Close(t, dbs)
		testBody(t, dbs, dbMock, 0)
	})

	t.Run("wait and restart", func(t *testing.T) {
		dbc, dbMock := dmltest.MockDB(t)
		defer dmltest.MockClose(t, dbc, dbMock)
		dbMock.MatchExpectationsInOrder(false)

		dbs, err := ccd.NewDBStorage(ccd.NewTableCollection(dbc.DB), ccd.Options{
			IdleRead:             time.Millisecond * 5,
			IdleWrite:            time.Millisecond * 5,
			SkipSchemaValidation: true,
		})
		require.NoError(t, err)
		defer dmltest.Close(t, dbs)

		testBody(t, dbs, dbMock, time.Millisecond*8)

		val, set, all := dbs.Statistics()
		assert.Exactly(t,
			"read ccd.stats{Open:0x0, Close:0x0} write ccd.stats{Open:0x3, Close:0x3} all ccd.stats{Open:0x0, Close:0x0}",
			fmt.Sprintf("read %#v write %#v all %#v", val, set, all),
		)
	})

	t.Run("query context timeout", func(t *testing.T) {
		dbc, dbMock := dmltest.MockDB(t)
		defer dmltest.MockClose(t, dbc, dbMock)
		dbMock.MatchExpectationsInOrder(false)

		dbs, err := ccd.NewDBStorage(ccd.NewTableCollection(dbc.DB), ccd.Options{
			ContextTimeoutWrite:  time.Millisecond * 50,
			SkipSchemaValidation: true,
		})
		require.NoError(t, err)
		defer dmltest.Close(t, dbs)

		prepIns := dbMock.ExpectPrepare(dmltest.SQLMockQuoteMeta("INSERT INTO `core_config_data` (`scope`,`scope_id`,`path`,`value`) VALUES (?,?,?,?) ON DUPLICATE KEY UPDATE `value`=VALUES(`value`)"))
		for i, test := range dbStorageMultiTests {

			prepIns.ExpectExec().
				WithArgs(test.scopeID, test.path, test.value).
				WillDelayFor(time.Millisecond * 110).
				WillReturnResult(sqlmock.NewResult(int64(i), 0))
			haveErr := dbs.Set(test.scopeID, test.path, test.value)

			causeErr := errors.Cause(haveErr)
			require.EqualError(t, causeErr, "canceling query due to user request")
		}

	})

}
