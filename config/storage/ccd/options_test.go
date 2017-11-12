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

	"github.com/corestoreio/pkg/config"
	"github.com/corestoreio/pkg/config/cfgpath"
	"github.com/corestoreio/pkg/config/storage/ccd"
	"github.com/corestoreio/pkg/util/cstesting"
	"github.com/stretchr/testify/assert"
)

// Not needed because columns have been supplied via csdb.WithTable() function
//func init() {
//	if _, err := csdb.GetDSN(); err == csdb.ErrDSNNotFound {
//		println("init()", err.Error(), "will skip loading of TableCollection")
//		return
//	}
//
//	dbc := csdb.MustConnectTest()
//	if err := ccd.TableCollection.Init(dbc.NewSession()); err != nil {
//		panic(err)
//	}
//	if err := dbc.Close(); err != nil {
//		panic(err)
//	}
//}

// Test_WithApplyCoreConfigData reads from the MySQL core_config_data table and applies
// these value to the underlying storage. tries to get back the values from the
// underlying storage
func Test_WithCoreConfigData(t *testing.T) {
	t.Parallel()

	dbc, dbMock := cstesting.MockDB(t)
	defer func() {
		dbMock.ExpectClose()

		assert.NoError(t, dbc.Close())

		if err := dbMock.ExpectationsWereMet(); err != nil {
			t.Error("there were unfulfilled expections", err)
		}
	}()

	sess := dbc.NewSession()

	dbMock.ExpectQuery("SELECT (.+) FROM `core_config_data` AS `main_table`").WillReturnRows(
		cstesting.MustMockRows(cstesting.WithFile("testdata", "core_config_data.csv")),
	)

	im := config.NewInMemoryStore()
	s := config.MustNewService(
		im,
		ccd.WithCoreConfigData(sess),
	)
	defer func() { assert.NoError(t, s.Close()) }()

	assert.NoError(t, s.Write(cfgpath.MustNewByParts("web/secure/offloader_header"), "SSL_OFFLOADED"))

	h, err := s.String(cfgpath.MustNewByParts("web/secure/offloader_header"))
	assert.NoError(t, err)
	assert.Exactly(t, "SSL_OFFLOADED", h)

	allKeys, err := im.AllKeys()
	assert.NoError(t, err)
	//for i, ak := range allKeys {
	//	t.Log(i, ak.String())
	//}
	assert.Len(t, allKeys, 21)

}
