// Copyright 2015 CoreStore Authors
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

// package store implements the handling of websites, groups and stores
package store

import (
	"errors"

	"github.com/corestoreio/csfw/storage/csdb"
	"github.com/corestoreio/csfw/storage/dbr"
)

const (
	DefaultWebsiteId int64 = 0
)

// WebsiteIndex used for iota and for not mixing up indexes
type WebsiteIndex int

var (
	ErrWebsiteNotFound = errors.New("Website not found")
	websiteCollection  WebsiteSlice
)

// GetWebsite uses a GroupIndex to return a group or an error.
// One should not modify the group object.
func GetWebsite(i WebsiteIndex) (*Website, error) {
	if int(i) < len(websiteCollection) {
		return websiteCollection[i], nil
	}
	return nil, ErrWebsiteNotFound
}

// GetWebsites returns a copy of the main slice of store groups.
// One should not modify the slice and its content.
// @todo $withDefault bool
func GetWebsites() WebsiteSlice {
	return websiteCollection
}

// Load uses a dbr session to load all data from the core_website table into the current slice.
// The variadic 2nd argument can be a call back function to manipulate the select.
// Additional columns or joins cannot be added. This method receiver should only be used in development.
// @see app/code/Magento/Store/Model/Resource/Website/Collection.php::Load()
func (s *WebsiteSlice) Load(dbrSess dbr.SessionRunner, cbs ...csdb.DbrSelectCb) (int, error) {
	return loadSlice(dbrSess, TableWebsite, &(*s), append(cbs, func(sb *dbr.SelectBuilder) *dbr.SelectBuilder {
		return sb.OrderBy("main_table.sort_order ASC").OrderBy("main_table.name ASC")
	})...)
}

/*
	@todo implement Magento\Store\Model\Website
*/
