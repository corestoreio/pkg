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
	"github.com/corestoreio/csfw/store/scope"
	"github.com/corestoreio/csfw/store/storemock"
	"github.com/stretchr/testify/assert"
)

func TestMustNewStoreAU_ConfigNil(t *testing.T) {
	sAU := storemock.MustNewStoreAU(cfgmock.NewService())
	assert.NotNil(t, sAU)
	assert.NotNil(t, sAU.Config)
	assert.NotNil(t, sAU.Website.Config)

	assert.Exactly(t, int64(5), sAU.Config.StoreID)
	assert.Exactly(t, int64(2), sAU.Config.WebsiteID)

	assert.Exactly(t, int64(0), sAU.Website.Config.StoreID)
	assert.Exactly(t, int64(2), sAU.Website.Config.WebsiteID)

}

func TestMustNewStoreAU_ConfigNonNil(t *testing.T) {
	sAU := storemock.MustNewStoreAU(cfgmock.NewService())
	assert.NotNil(t, sAU)
	assert.NotNil(t, sAU.Config)
	assert.NotNil(t, sAU.Website.Config)
}

func TestMustNewStoreAU_Config(t *testing.T) {
	var configPath = cfgpath.MustNewByParts("aa/bb/cc")

	sm := cfgmock.NewService(cfgmock.PathValue{
		configPath.String():                "DefaultScopeString",
		configPath.BindWebsite(2).String(): "WebsiteScopeString",
		configPath.BindStore(5).String():   "StoreScopeString",
	})
	aust := storemock.MustNewStoreAU(sm)

	haveS, err := aust.Website.Config.String(configPath.Route)
	if err != nil {
		t.Fatalf("%+v", err)
	}
	assert.Exactly(t, "WebsiteScopeString", haveS)

	haveS, err = aust.Website.Config.String(configPath.Route, scope.Default)
	if err != nil {
		t.Fatalf("%+v", err)
	}
	assert.Exactly(t, "DefaultScopeString", haveS)

	haveS, err = aust.Config.String(configPath.Route)
	if err != nil {
		t.Fatalf("%+v", err)
	}
	assert.Exactly(t, "StoreScopeString", haveS)

	haveS, err = aust.Config.String(configPath.Route, scope.Default)
	if err != nil {
		t.Fatalf("%+v", err)
	}
	assert.Exactly(t, "DefaultScopeString", haveS)

	assert.Exactly(t, scope.TypeIDs{scope.DefaultTypeID, scope.Website.Pack(2), scope.Store.Pack(5)}, sm.AllInvocations().TypeIDs())
	assert.Exactly(t, 3, sm.AllInvocations().PathCount())

}
