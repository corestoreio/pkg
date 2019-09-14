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

// +build csall db

package storage_test

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/corestoreio/errors"
	"github.com/corestoreio/pkg/config"
	"github.com/corestoreio/pkg/config/storage"
	"github.com/corestoreio/pkg/sql/ddl"
	"github.com/corestoreio/pkg/sql/dmltest"
	"github.com/corestoreio/pkg/store/scope"
	"github.com/corestoreio/pkg/util/assert"
	"github.com/fortytw2/leaktest"
)

var _ config.Storager = (*storage.DB)(nil)

func mustNewTables(ctx context.Context, opts ...ddl.TableOption) (tm *ddl.Tables) {
	t, err := storage.NewTables(ctx, opts...)
	if err != nil {
		panic(err)
	}
	return t
}

func TestMustNewDB_Panic(t *testing.T) {
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
	_ = storage.MustNewDB(mustNewTables(context.TODO()), storage.DBOptions{
		TableName:            "non-existent",
		SkipSchemaValidation: true,
	})
}

func TestService_AllKeys_Mocked(t *testing.T) {
	defer leaktest.CheckTimeout(t, time.Second)()

	dbc, dbMock := dmltest.MockDB(t)
	defer dmltest.MockClose(t, dbc, dbMock)

	t.Run("table not found", func(t *testing.T) {
		dbs, err := storage.NewDB(mustNewTables(context.TODO()), storage.DBOptions{
			TableName:            "non-existent",
			SkipSchemaValidation: true,
		})
		assert.Nil(t, dbs)
		assert.True(t, errors.NotFound.Match(err), "%+v", err)
	})

	t.Run("no leaking goroutines", func(t *testing.T) {
		dbMock.ExpectQuery("SELECT.+FROM information_schema.COLUMNS").WithArgs().WillReturnRows(
			dmltest.MustMockRows(dmltest.WithFile("testdata", "core_configuration_columns.csv")),
		)

		dbs, err := storage.NewDB(mustNewTables(context.TODO(), ddl.WithConnPool(dbc)), storage.DBOptions{
			SkipSchemaValidation: true,
		})
		assert.NoError(t, err)
		assert.NoError(t, dbs.Close())
	})
}

var serviceMultiTests = []struct {
	path    string
	scopeID scope.TypeID
	value   []byte
}{
	{"testService/secure/base_url", scope.Website.WithID(10), []byte("http://corestore.io")},
	{"testService/log/active", scope.Store.WithID(9), []byte("https://crestre.i")},
	{"testService/checkout/multishipping", scope.DefaultTypeID, []byte("false")},
}

func TestService_Get(t *testing.T) {
	defer leaktest.CheckTimeout(t, time.Second)()

	testBody := func(t *testing.T, dbs *storage.DB, dbMock sqlmock.Sqlmock, sleep time.Duration) {
		prepSel := dbMock.ExpectPrepare(dmltest.SQLMockQuoteMeta("SELECT `value` FROM `core_configuration` AS `main_table` WHERE (`scope` = ?) AND (`scope_id` = ?) AND (`path` = ?)"))
		for _, test := range serviceMultiTests {
			scp, sID := test.scopeID.Unpack()
			prepSel.ExpectQuery().WithArgs(scp.StrType(), sID, test.path).WillReturnRows(sqlmock.NewRows([]string{"value"}))

			haveVal, haveOK, haveErr := dbs.Get(config.MustNewPathWithScope(test.scopeID, test.path))
			assert.NoError(t, haveErr)
			assert.False(t, haveOK, "%s Value with path %q should NOT be found", test.scopeID, test.path)
			assert.Exactly(t, []byte(nil), haveVal)
		}

		if sleep > 0 {
			time.Sleep(sleep)
			prepSel = dbMock.ExpectPrepare(dmltest.SQLMockQuoteMeta("SELECT `value` FROM `core_configuration` AS `main_table` WHERE (`scope` = ?) AND (`scope_id` = ?) AND (`path` = ?)"))
		}

		for _, test := range serviceMultiTests {
			scp, sID := test.scopeID.Unpack()
			prepSel.ExpectQuery().WithArgs(scp.StrType(), sID, test.path).WillReturnRows(sqlmock.NewRows([]string{"value"}).AddRow(test.value))

			haveVal, haveOK, haveErr := dbs.Get(config.MustNewPathWithScope(test.scopeID, test.path))
			assert.NoError(t, haveErr)
			assert.True(t, haveOK, "%s Value with path %q should be found", test.scopeID, test.path)
			assert.Exactly(t, test.value, haveVal)
		}
	}

	t.Run("no waiting", func(t *testing.T) {
		dbc, dbMock := dmltest.MockDB(t)
		defer dmltest.MockClose(t, dbc, dbMock)
		dbMock.MatchExpectationsInOrder(false)

		dbMock.ExpectQuery("SELECT.+FROM information_schema.COLUMNS").WithArgs().WillReturnRows(
			dmltest.MustMockRows(dmltest.WithFile("testdata", "core_configuration_columns.csv")),
		)

		dbs, err := storage.NewDB(mustNewTables(context.TODO(), ddl.WithConnPool(dbc)), storage.DBOptions{
			SkipSchemaValidation: true,
		})
		assert.NoError(t, err)
		defer dmltest.Close(t, dbs)
		testBody(t, dbs, dbMock, 0)
	})

	t.Run("wait and restart", func(t *testing.T) {
		dbc, dbMock := dmltest.MockDB(t)
		defer dmltest.MockClose(t, dbc, dbMock)
		dbMock.MatchExpectationsInOrder(false)

		dbMock.ExpectQuery("SELECT.+FROM information_schema.COLUMNS").WithArgs().WillReturnRows(
			dmltest.MustMockRows(dmltest.WithFile("testdata", "core_configuration_columns.csv")),
		)

		dbs, err := storage.NewDB(mustNewTables(context.TODO(), ddl.WithConnPool(dbc)), storage.DBOptions{
			IdleRead:             time.Millisecond * 50,
			IdleWrite:            time.Millisecond * 50,
			SkipSchemaValidation: true,
		})
		assert.NoError(t, err)
		defer dmltest.Close(t, dbs)

		testBody(t, dbs, dbMock, time.Millisecond*100)

		val, set := dbs.Statistics()
		assert.Exactly(t,
			"read storage.dbStats{Open:0x2, Close:0x1} write storage.dbStats{Open:0x0, Close:0x0}",
			fmt.Sprintf("read %#v write %#v", val, set),
		)
	})

	t.Run("query context timeout", func(t *testing.T) {
		dbc, dbMock := dmltest.MockDB(t)
		defer dmltest.MockClose(t, dbc, dbMock)
		dbMock.MatchExpectationsInOrder(false)

		dbMock.ExpectQuery("SELECT.+FROM information_schema.COLUMNS").WithArgs().WillReturnRows(
			dmltest.MustMockRows(dmltest.WithFile("testdata", "core_configuration_columns.csv")),
		)

		dbs, err := storage.NewDB(mustNewTables(context.TODO(), ddl.WithConnPool(dbc)), storage.DBOptions{
			ContextTimeoutRead:   time.Millisecond * 50,
			SkipSchemaValidation: true,
		})
		assert.NoError(t, err)
		defer dmltest.Close(t, dbs)

		prepSel := dbMock.ExpectPrepare(dmltest.SQLMockQuoteMeta("SELECT `value` FROM `core_configuration` AS `main_table` WHERE (`scope` = ?) AND (`scope_id` = ?) AND (`path` = ?)"))
		for _, test := range serviceMultiTests {
			scp, sID := test.scopeID.Unpack()
			prepSel.ExpectQuery().WithArgs(scp.StrType(), sID, test.path).WillDelayFor(time.Millisecond * 110).WillReturnRows(sqlmock.NewRows([]string{"value"}))

			haveVal, haveOK, haveErr := dbs.Get(config.MustNewPathWithScope(test.scopeID, test.path))
			assert.Nil(t, haveVal)
			assert.False(t, haveOK)
			causeErr := errors.Cause(haveErr)
			assert.EqualError(t, causeErr, "canceling query due to user request")
			return
		}
	})
}

func TestService_Set(t *testing.T) {
	defer leaktest.CheckTimeout(t, time.Second)()

	testBody := func(t *testing.T, dbs *storage.DB, dbMock sqlmock.Sqlmock, sleep time.Duration) {
		prepIns := dbMock.ExpectPrepare(dmltest.SQLMockQuoteMeta("INSERT INTO `core_configuration` (`scope`,`scope_id`,`path`,`value`) VALUES (?,?,?,?) ON DUPLICATE KEY UPDATE `value`=VALUES(`value`)"))

		for i, test := range serviceMultiTests {
			j := int64(i + 1)

			if sleep > 0 && i > 0 {
				prepIns = dbMock.ExpectPrepare(dmltest.SQLMockQuoteMeta("INSERT INTO `core_configuration` (`scope`,`scope_id`,`path`,`value`) VALUES (?,?,?,?) ON DUPLICATE KEY UPDATE `value`=VALUES(`value`)"))
			}

			prepIns.ExpectExec().
				WithArgs(test.scopeID, test.path, test.value).
				WillReturnResult(sqlmock.NewResult(j, 0))
			assert.NoError(t, dbs.Set(config.MustNewPathWithScope(test.scopeID, test.path), test.value))

			if sleep > 0 {
				time.Sleep(sleep)
			}
		}
	}

	t.Run("no waiting", func(t *testing.T) {
		dbc, dbMock := dmltest.MockDB(t)
		defer dmltest.MockClose(t, dbc, dbMock)
		dbMock.MatchExpectationsInOrder(false)

		dbMock.ExpectQuery("SELECT.+FROM information_schema.COLUMNS").WithArgs().WillReturnRows(
			dmltest.MustMockRows(dmltest.WithFile("testdata", "core_configuration_columns.csv")),
		)

		dbs, err := storage.NewDB(mustNewTables(context.TODO(), ddl.WithConnPool(dbc)), storage.DBOptions{
			SkipSchemaValidation: true,
		})
		assert.NoError(t, err)
		defer dmltest.Close(t, dbs)
		testBody(t, dbs, dbMock, 0)
	})

	t.Run("wait and restart", func(t *testing.T) {
		dbc, dbMock := dmltest.MockDB(t)
		defer dmltest.MockClose(t, dbc, dbMock)
		dbMock.MatchExpectationsInOrder(false)

		dbMock.ExpectQuery("SELECT.+FROM information_schema.COLUMNS").WithArgs().WillReturnRows(
			dmltest.MustMockRows(dmltest.WithFile("testdata", "core_configuration_columns.csv")),
		)

		dbs, err := storage.NewDB(mustNewTables(context.TODO(), ddl.WithConnPool(dbc)), storage.DBOptions{
			IdleRead:             time.Millisecond * 5,
			IdleWrite:            time.Millisecond * 5,
			SkipSchemaValidation: true,
		})
		assert.NoError(t, err)
		defer dmltest.Close(t, dbs)

		testBody(t, dbs, dbMock, time.Millisecond*8)

		val, set := dbs.Statistics()
		assert.Exactly(t,
			"read storage.dbStats{Open:0x0, Close:0x0} write storage.dbStats{Open:0x3, Close:0x3}",
			fmt.Sprintf("read %#v write %#v", val, set),
		)
	})

	t.Run("query context timeout", func(t *testing.T) {
		dbc, dbMock := dmltest.MockDB(t)
		defer dmltest.MockClose(t, dbc, dbMock)
		dbMock.MatchExpectationsInOrder(false)

		dbMock.ExpectQuery("SELECT.+FROM information_schema.COLUMNS").WithArgs().WillReturnRows(
			dmltest.MustMockRows(dmltest.WithFile("testdata", "core_configuration_columns.csv")),
		)

		dbs, err := storage.NewDB(mustNewTables(context.TODO(), ddl.WithConnPool(dbc)), storage.DBOptions{
			ContextTimeoutWrite:  time.Millisecond * 50,
			SkipSchemaValidation: true,
		})
		assert.NoError(t, err)
		defer dmltest.Close(t, dbs)

		prepIns := dbMock.ExpectPrepare(dmltest.SQLMockQuoteMeta("INSERT INTO `core_configuration` (`scope`,`scope_id`,`path`,`value`) VALUES (?,?,?,?) ON DUPLICATE KEY UPDATE `value`=VALUES(`value`)"))
		for i, test := range serviceMultiTests {

			prepIns.ExpectExec().
				WithArgs(test.scopeID, test.path, test.value).
				WillDelayFor(time.Millisecond * 110).
				WillReturnResult(sqlmock.NewResult(int64(i), 0))
			haveErr := dbs.Set(config.MustNewPathWithScope(test.scopeID, test.path), test.value)

			causeErr := errors.Cause(haveErr)
			assert.EqualError(t, causeErr, "canceling query due to user request")
		}
	})
}

// Test_WithApplyCoreConfigData reads from the MySQL core_configuration table and applies
// these value to the underlying storage. tries to get back the values from the
// underlying storage
func Test_WithCoreConfigData(t *testing.T) {
	t.Parallel()

	dbc, dbMock := dmltest.MockDB(t)
	defer dmltest.MockClose(t, dbc, dbMock)

	dbMock.ExpectQuery("SELECT.+FROM information_schema.COLUMNS").WithArgs().WillReturnRows(
		dmltest.MustMockRows(dmltest.WithFile("testdata", "core_configuration_columns.csv")),
	)

	dbMock.ExpectQuery("SELECT (.+) FROM `core_configuration` AS `main_table`").WillReturnRows(
		dmltest.MustMockRows(dmltest.WithFile("testdata", "core_configuration.csv")),
	)

	tbls := mustNewTables(context.TODO(), ddl.WithConnPool(dbc))

	im := storage.NewMap()
	s := config.MustNewService(
		im,
		config.Options{},
		storage.WithLoadFromDB(tbls, storage.DBOptions{}),
	)
	defer dmltest.Close(t, s)

	p1 := config.MustNewPath("web/secure/offloader_header").BindStore(987)
	assert.NoError(t, s.Set(p1, []byte("SSL_OFFLOADED")))

	v, ok, err := s.Get(p1).Str()
	assert.NoError(t, err)
	assert.True(t, ok)
	assert.Exactly(t, "SSL_OFFLOADED", v)

	p2 := config.MustNewPath("web/unsecure/base_skin_url").BindWebsite(44)
	v, ok, err = s.Get(p2).Str()
	assert.NoError(t, err)
	assert.True(t, ok)
	assert.Exactly(t, "{{unsecure_base_url}}skin/", v)
}
