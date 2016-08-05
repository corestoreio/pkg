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

	aust := storemock.MustNewStoreAU(cfgmock.NewService(cfgmock.WithPV(cfgmock.PathValue{
		configPath.Bind(scope.Default, 0).String(): "DefaultScopeString",
		configPath.Bind(scope.Website, 2).String(): "WebsiteScopeString",
		configPath.Bind(scope.Store, 5).String():   "StoreScopeString",
	})))

	haveS, scp, err := aust.Website.Config.String(configPath.Route)
	if err != nil {
		t.Fatalf("%+v", err)
	}
	assert.Exactly(t, "WebsiteScopeString", haveS)
	assert.Exactly(t, scope.NewHash(scope.Website, 2), scp)

	haveS, scp, err = aust.Website.Config.String(configPath.Route, scope.Default)
	if err != nil {
		t.Fatalf("%+v", err)
	}
	assert.Exactly(t, "DefaultScopeString", haveS)
	assert.Exactly(t, scope.DefaultHash, scp)

	haveS, scp, err = aust.Config.String(configPath.Route)
	if err != nil {
		t.Fatalf("%+v", err)
	}
	assert.Exactly(t, "StoreScopeString", haveS)
	assert.Exactly(t, scope.NewHash(scope.Store, 5), scp)

	haveS, scp, err = aust.Config.String(configPath.Route, scope.Default)
	if err != nil {
		t.Fatalf("%+v", err)
	}
	assert.Exactly(t, "DefaultScopeString", haveS)
	assert.Exactly(t, scope.DefaultHash, scp)

}
