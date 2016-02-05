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

	"github.com/corestoreio/csfw/config"
	"github.com/corestoreio/csfw/config/ccd"
	"github.com/corestoreio/csfw/config/path"
	"github.com/corestoreio/csfw/storage/csdb"
	"github.com/stretchr/testify/assert"
)

func init() {
	if _, err := csdb.GetDSNTest(); err == csdb.ErrDSNTestNotFound {
		println("init()", err.Error(), "will skip loading of TableCollection")
		return
	}

	dbc := csdb.MustConnectTest()
	if err := ccd.TableCollection.Init(dbc.NewSession()); err != nil {
		panic(err)
	}
	if err := dbc.Close(); err != nil {
		panic(err)
	}
}

// Test_WithApplyCoreConfigData reads from the MySQL core_config_data table and applies
// these value to the underlying storage. tries to get back the values from the
// underlying storage
func Test_WithCoreConfigData(t *testing.T) {
	defer debugLogBuf.Reset()
	defer infoLogBuf.Reset()
	if _, err := csdb.GetDSNTest(); err == csdb.ErrDSNTestNotFound {
		t.Skip(err)
	}

	dbc := csdb.MustConnectTest()
	defer func() { assert.NoError(t, dbc.Close()) }()
	sess := dbc.NewSession(nil) // nil tricks the NewSession ;-)

	s := config.MustNewService(
		ccd.WithCoreConfigData(sess),
	)
	defer func() { assert.NoError(t, s.Close()) }()

	//	println("\n", debugLogBuf.String(), "\n")
	//	println("\n", infoLogBuf.String(), "\n")

	assert.NoError(t, s.Write(path.MustNewByParts("web/secure/offloader_header"), "SSL_OFFLOADED"))

	h, err := s.String(path.MustNewByParts("web/secure/offloader_header"))
	assert.NoError(t, err)
	assert.Exactly(t, "SSL_OFFLOADED", h)

	allKeys, err := s.Storage.AllKeys()
	assert.NoError(t, err)

	assert.True(t, len(allKeys) > 170) // TODO: refactor this if else and use a clean database ...

}
