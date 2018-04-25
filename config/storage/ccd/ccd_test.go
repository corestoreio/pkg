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
