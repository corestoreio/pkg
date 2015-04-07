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

// Package store implements the handling of websites, groups and stores
package store

import (
	"errors"

	"github.com/corestoreio/csfw/storage/csdb"
	"github.com/corestoreio/csfw/storage/dbr"
	"github.com/juju/errgo"
)

const (
	DefaultWebsiteId int64 = 0
)

type (
	// WebsiteIndex used for iota and for not mixing up indexes
	WebsiteIndex  int
	WebsiteGetter interface {
		// ByID returns a WebsiteIndex using the StoreID. This WebsiteIndex identifies a website within a WebsiteSlice.
		ByID(id int64) (WebsiteIndex, error)
		// ByCode returns a WebsiteIndex using the code. This WebsiteIndex identifies a website within a WebsiteSlice.
		ByCode(code string) (WebsiteIndex, error)
	}
)

var (
	ErrWebsiteNotFound = errors.New("Website not found")
	websiteCollection  TableWebsiteSlice
	websiteGetter      WebsiteGetter
)

func SetWebsiteCollection(sc TableWebsiteSlice) {
	if len(sc) == 0 {
		panic("WebsiteSlice is empty")
	}
	websiteCollection = sc
}

func SetWebsiteGetter(g WebsiteGetter) {
	if g == nil {
		panic("WebsiteGetter cannot be nil")
	}
	websiteGetter = g
}

// GetWebsite uses a GroupIndex to return a group or an error.
// One should not modify the group object.
func GetWebsite(i WebsiteIndex) (*TableWebsite, error) {
	if int(i) < len(websiteCollection) {
		return websiteCollection[i], nil
	}
	return nil, ErrWebsiteNotFound
}

func GetWebsiteByID(id int64) (*TableWebsite, error) {
	return websiteCollection.ByID(id)
}

func GetWebsiteByCode(code string) (*TableWebsite, error) {
	return websiteCollection.ByCode(code)
}

// GetWebsites returns a copy of the main slice of websites.
// One should not modify the slice and its content.
func GetWebsites() TableWebsiteSlice {
	return websiteCollection
}

// Load uses a dbr session to load all data from the core_website table into the current slice.
// The variadic 2nd argument can be a call back function to manipulate the select.
// Additional columns or joins cannot be added. This method receiver should only be used in development.
// @see app/code/Magento/Store/Model/Resource/Website/Collection.php::Load()
func (s *TableWebsiteSlice) Load(dbrSess dbr.SessionRunner, cbs ...csdb.DbrSelectCb) (int, error) {
	return loadSlice(dbrSess, TableIndexWebsite, &(*s), append(cbs, func(sb *dbr.SelectBuilder) *dbr.SelectBuilder {
		return sb.OrderBy("main_table.sort_order ASC").OrderBy("main_table.name ASC")
	})...)
}

func (s TableWebsiteSlice) ByID(id int64) (*TableWebsite, error) {
	i, err := websiteGetter.ByID(id)
	if err != nil {
		return nil, errgo.Mask(err)
	}
	return s[i], nil
}

func (s TableWebsiteSlice) ByCode(code string) (*TableWebsite, error) {
	i, err := websiteGetter.ByCode(code)
	if err != nil {
		return nil, errgo.Mask(err)
	}
	return s[i], nil
}

// @todo review Magento code because of column is_default
//func (s Website) IsDefault() bool {
//	return s.WebsiteID == DefaultWebsiteId
//}

/*
	@todo implement Magento\Store\Model\Website
*/
