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

package jwt_test

import (
	"github.com/corestoreio/csfw/net/jwt"
	"github.com/corestoreio/csfw/store"
	"github.com/corestoreio/csfw/store/scope"
)

type storeFinderMock struct{}

func (storeFinderMock) DefaultStoreID(runMode scope.TypeID) (storeID, websiteID int64, err error) {
	return
}

func (storeFinderMock) StoreIDbyCode(runMode scope.TypeID, storeCode string) (storeID, websiteID int64, err error) {
	return
}

var _ jwt.StoreFinder = (*storeFinderMock)(nil)
var _ store.Finder = (*storeFinderMock)(nil)
