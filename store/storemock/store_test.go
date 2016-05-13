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

package storemock_test

import (
	"testing"

	"github.com/corestoreio/csfw/config/cfgmock"
	"github.com/corestoreio/csfw/config/cfgpath"
	"github.com/corestoreio/csfw/store"
	"github.com/corestoreio/csfw/store/scope"
	"github.com/corestoreio/csfw/store/storemock"
	"github.com/corestoreio/csfw/util/errors"
	"github.com/stretchr/testify/assert"
)

var _ store.Requester = (*storemock.RequestedStoreAU)(nil)

func TestMustNewStoreAU_ConfigNil(t *testing.T) {
	sAU := storemock.MustNewStoreAU(nil)
	assert.NotNil(t, sAU)
	assert.Nil(t, sAU.Config)
	assert.Nil(t, sAU.Website.Config)
}

func TestMustNewStoreAU_ConfigNonNil(t *testing.T) {
	sAU := storemock.MustNewStoreAU(cfgmock.NewService())
	assert.NotNil(t, sAU)
	assert.NotNil(t, sAU.Config)
	assert.NotNil(t, sAU.Website.Config)
}

func TestRequestedStoreAU(t *testing.T) {

	rsau := &storemock.RequestedStoreAU{
		Getter: cfgmock.NewService(),
	}
	rStore, err := rsau.RequestedStore(scope.Option{Store: scope.MockCode("unimportant")})
	if err != nil {
		t.Fatal(err)
	}
	assert.NotNil(t, rStore)
	assert.NotNil(t, rStore.Config)
	assert.NotNil(t, rStore.Website.Config)
}

func TestMustNewStoreAU_Config(t *testing.T) {
	var configPath = cfgpath.MustNewByParts("aa/bb/cc")

	aust := storemock.MustNewStoreAU(cfgmock.NewService(cfgmock.WithPV(cfgmock.PathValue{
		configPath.Bind(scope.Default, 0).String(): "DefaultScopeString",
		configPath.Bind(scope.Website, 2).String(): "WebsiteScopeString",
		configPath.Bind(scope.Store, 5).String():   "StoreScopeString",
	})))

	haveS, err := aust.Website.Config.String(configPath.Route)
	if err != nil {
		t.Fatal("fatal", errors.PrintLoc(err))
	}
	assert.Exactly(t, "WebsiteScopeString", haveS)

	haveS, err = aust.Website.Config.String(configPath.Route, scope.Default)
	if err != nil {
		t.Fatal("fatal", errors.PrintLoc(err))
	}
	assert.Exactly(t, "DefaultScopeString", haveS)

	haveS, err = aust.Config.String(configPath.Route)
	if err != nil {
		t.Fatal("fatal", errors.PrintLoc(err))
	}
	assert.Exactly(t, "StoreScopeString", haveS)

	haveS, err = aust.Config.String(configPath.Route, scope.Default)
	if err != nil {
		t.Fatal("fatal", errors.PrintLoc(err))
	}
	assert.Exactly(t, "DefaultScopeString", haveS)

}
