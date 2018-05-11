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
	"testing"

	"github.com/corestoreio/pkg/config"
	"github.com/corestoreio/pkg/config/storage"
	"github.com/corestoreio/pkg/config/storage/cfgdb"
	"github.com/corestoreio/pkg/sql/dmltest"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// Test_WithApplyCoreConfigData reads from the MySQL core_config_data table and applies
// these value to the underlying storage. tries to get back the values from the
// underlying storage
func Test_WithCoreConfigData(t *testing.T) {
	t.Parallel()

	dbc, dbMock := dmltest.MockDB(t)
	defer dmltest.MockClose(t, dbc, dbMock)

	dbMock.ExpectQuery("SELECT (.+) FROM `core_config_data` AS `main_table`").WillReturnRows(
		dmltest.MustMockRows(dmltest.WithFile("testdata", "core_config_data.csv")),
	)

	tbls := cfgdb.NewTableCollection(dbc.DB)

	im := storage.NewMap()
	s := config.MustNewService(
		im,
		config.Options{},
		cfgdb.WithLoadFromDB(tbls, cfgdb.Options{}),
	)
	defer dmltest.Close(t, s)

	p1 := config.MustNewPath("web/secure/offloader_header").BindStore(987)
	assert.NoError(t, s.Set(p1, []byte("SSL_OFFLOADED")))

	v, ok, err := s.Get(p1).Str()
	require.NoError(t, err)
	require.True(t, ok)
	assert.Exactly(t, "SSL_OFFLOADED", v)

	p2 := config.MustNewPath("web/unsecure/base_skin_url").BindWebsite(44)
	v, ok, err = s.Get(p2).Str()
	require.NoError(t, err)
	require.True(t, ok)
	assert.Exactly(t, "{{unsecure_base_url}}skin/", v)
}
