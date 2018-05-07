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

package cfgdb_test

import (
	"fmt"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/corestoreio/errors"
	"github.com/corestoreio/pkg/config"
	"github.com/corestoreio/pkg/config/storage/cfgdb"
	"github.com/corestoreio/pkg/sql/dmltest"
	"github.com/corestoreio/pkg/store/scope"
	"github.com/fortytw2/leaktest"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var _ config.Storager = (*cfgdb.DBStorage)(nil)

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
	_ = cfgdb.MustNewDBStorage(cfgdb.NewTableCollection(nil), cfgdb.Options{
		TableName:            "non-existent",
		SkipSchemaValidation: true,
	})
}

func TestDBStorage_AllKeys_Mocked(t *testing.T) {
	defer leaktest.CheckTimeout(t, time.Second)()

	dbc, dbMock := dmltest.MockDB(t)
	defer dmltest.MockClose(t, dbc, dbMock)

	t.Run("table not found", func(t *testing.T) {
		dbs, err := cfgdb.NewDBStorage(cfgdb.NewTableCollection(dbc.DB), cfgdb.Options{
			TableName:            "non-existent",
			SkipSchemaValidation: true,
		})
		assert.Nil(t, dbs)
		assert.True(t, errors.NotFound.Match(err), "%+v", err)
	})

	t.Run("no leaking goroutines", func(t *testing.T) {
		dbs, err := cfgdb.NewDBStorage(cfgdb.NewTableCollection(dbc.DB), cfgdb.Options{
			SkipSchemaValidation: true,
		})
		require.NoError(t, err)
		assert.NoError(t, dbs.Close())
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

func TestDBStorage_Get(t *testing.T) {
	defer leaktest.CheckTimeout(t, time.Second)()

	testBody := func(t *testing.T, dbs *cfgdb.DBStorage, dbMock sqlmock.Sqlmock, sleep time.Duration) {

		prepSel := dbMock.ExpectPrepare(dmltest.SQLMockQuoteMeta("SELECT `value` FROM `core_config_data` AS `main_table` WHERE (`scope` = ?) AND (`scope_id` = ?) AND (`path` = ?)"))
		for _, test := range dbStorageMultiTests {
			scp, sID := test.scopeID.Unpack()
			prepSel.ExpectQuery().WithArgs(scp.StrType(), sID, test.path).WillReturnRows(sqlmock.NewRows([]string{"value"}))

			haveVal, haveOK, haveErr := dbs.Get(test.scopeID, test.path)
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

			haveVal, haveOK, haveErr := dbs.Get(test.scopeID, test.path)
			require.NoError(t, haveErr)
			require.True(t, haveOK, "%s Value with path %q should be found", test.scopeID, test.path)
			assert.Exactly(t, test.value, haveVal)
		}
	}

	t.Run("no waiting", func(t *testing.T) {
		dbc, dbMock := dmltest.MockDB(t)
		defer dmltest.MockClose(t, dbc, dbMock)
		dbMock.MatchExpectationsInOrder(false)

		dbs, err := cfgdb.NewDBStorage(cfgdb.NewTableCollection(dbc.DB), cfgdb.Options{
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

		dbs, err := cfgdb.NewDBStorage(cfgdb.NewTableCollection(dbc.DB), cfgdb.Options{
			IdleRead:             time.Millisecond * 50,
			IdleWrite:            time.Millisecond * 50,
			SkipSchemaValidation: true,
		})
		require.NoError(t, err)
		defer dmltest.Close(t, dbs)

		testBody(t, dbs, dbMock, time.Millisecond*100)

		val, set := dbs.Statistics()
		assert.Exactly(t,
			"read cfgdb.stats{Open:0x2, Close:0x1} write cfgdb.stats{Open:0x0, Close:0x0}",
			fmt.Sprintf("read %#v write %#v", val, set),
		)
	})

	t.Run("query context timeout", func(t *testing.T) {
		dbc, dbMock := dmltest.MockDB(t)
		defer dmltest.MockClose(t, dbc, dbMock)
		dbMock.MatchExpectationsInOrder(false)

		dbs, err := cfgdb.NewDBStorage(cfgdb.NewTableCollection(dbc.DB), cfgdb.Options{
			ContextTimeoutRead:   time.Millisecond * 50,
			SkipSchemaValidation: true,
		})
		require.NoError(t, err)
		defer dmltest.Close(t, dbs)

		prepSel := dbMock.ExpectPrepare(dmltest.SQLMockQuoteMeta("SELECT `value` FROM `core_config_data` AS `main_table` WHERE (`scope` = ?) AND (`scope_id` = ?) AND (`path` = ?)"))
		for _, test := range dbStorageMultiTests {
			scp, sID := test.scopeID.Unpack()
			prepSel.ExpectQuery().WithArgs(scp.StrType(), sID, test.path).WillDelayFor(time.Millisecond * 110).WillReturnRows(sqlmock.NewRows([]string{"value"}))

			haveVal, haveOK, haveErr := dbs.Get(test.scopeID, test.path)
			assert.Nil(t, haveVal)
			assert.False(t, haveOK)
			causeErr := errors.Cause(haveErr)
			require.EqualError(t, causeErr, "canceling query due to user request")
			return
		}

	})
}

func TestDBStorage_Set(t *testing.T) {
	defer leaktest.CheckTimeout(t, time.Second)()

	testBody := func(t *testing.T, dbs *cfgdb.DBStorage, dbMock sqlmock.Sqlmock, sleep time.Duration) {

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

		dbs, err := cfgdb.NewDBStorage(cfgdb.NewTableCollection(dbc.DB), cfgdb.Options{
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

		dbs, err := cfgdb.NewDBStorage(cfgdb.NewTableCollection(dbc.DB), cfgdb.Options{
			IdleRead:             time.Millisecond * 5,
			IdleWrite:            time.Millisecond * 5,
			SkipSchemaValidation: true,
		})
		require.NoError(t, err)
		defer dmltest.Close(t, dbs)

		testBody(t, dbs, dbMock, time.Millisecond*8)

		val, set := dbs.Statistics()
		assert.Exactly(t,
			"read cfgdb.stats{Open:0x0, Close:0x0} write cfgdb.stats{Open:0x3, Close:0x3}",
			fmt.Sprintf("read %#v write %#v", val, set),
		)
	})

	t.Run("query context timeout", func(t *testing.T) {
		dbc, dbMock := dmltest.MockDB(t)
		defer dmltest.MockClose(t, dbc, dbMock)
		dbMock.MatchExpectationsInOrder(false)

		dbs, err := cfgdb.NewDBStorage(cfgdb.NewTableCollection(dbc.DB), cfgdb.Options{
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
