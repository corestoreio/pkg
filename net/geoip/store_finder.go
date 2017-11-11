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

package geoip

import "github.com/corestoreio/cspkg/store/scope"

// StoreFinder see store.Finder for a description. Usage of this interface in
// WithRunMode() middleware.
type StoreFinder interface {
	// DefaultStoreID returns the default active store ID and its website ID
	// depending on the run mode. Error behaviour is mostly of type NotValid.
	DefaultStoreID(runMode scope.TypeID) (storeID, websiteID int64, err error)
	// StoreIDbyCode returns, depending on the runMode, for a storeCode its
	// active store ID and its website ID. An empty runMode hash falls back to
	// select the default website with its default group and the slice of
	// default stores. A not-found error behaviour gets returned if the code
	// cannot be found. If the runMode equals to scope.DefaultTypeID, the returned
	// ID is always 0 and error is nil.
	StoreIDbyCode(runMode scope.TypeID, storeCode string) (storeID, websiteID int64, err error)
}
