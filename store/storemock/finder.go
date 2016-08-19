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

package storemock

import (
	"github.com/corestoreio/csfw/store"
	"github.com/corestoreio/csfw/store/scope"
)

var _ store.Finder = (*Find)(nil)

// Find implements interface store.Finder for mocking in tests.
type Find struct {
	// Next three get returned by function DefaultStoreID()
	StoreIDDefault   int64
	WebsiteIDDefault int64
	StoreIDError     error

	// Next three gets returned by function StoreIDbyCode()
	IDByCodeStoreID   int64
	IDByCodeWebsiteID int64
	IDByCodeError     error
}

func (s Find) DefaultStoreID(runMode scope.Hash) (storeID, websiteID int64, err error) {
	return s.StoreIDDefault, s.WebsiteIDDefault, s.StoreIDError
}

func (s Find) StoreIDbyCode(runMode scope.Hash, storeCode string) (storeID, websiteID int64, err error) {
	return s.IDByCodeStoreID, s.IDByCodeWebsiteID, s.IDByCodeError
}
