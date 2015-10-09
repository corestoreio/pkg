// Copyright 2015, Cyrill @ Schumacher.fm and the CoreStore contributors
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

package store

import (
	"github.com/corestoreio/csfw/storage/csdb"
	"github.com/corestoreio/csfw/storage/dbr"
)

/*
	TableWebsite and TableWebsiteSlice method receivers
*/

// Load uses a dbr session to load all data from the core_website table into the current slice.
// The variadic 2nd argument can be a call back function to manipulate the select.
// Additional columns or joins cannot be added. This method receiver should only be used in development.
// @see app/code/Magento/Store/Model/Resource/Website/Collection.php::Load()
func (s *TableWebsiteSlice) SQLSelect(dbrSess dbr.SessionRunner, cbs ...csdb.DbrSelectCb) (int, error) {
	return s.parentSQLSelect(dbrSess, append(append([]csdb.DbrSelectCb{nil}, func(sb *dbr.SelectBuilder) *dbr.SelectBuilder {
		return sb.OrderBy("main_table.sort_order ASC").OrderBy("main_table.name ASC")
	}), cbs...)...)
}

// @todo review Magento code because of column is_default
//func (s TableWebsite) IsDefault() bool {
//	return s.WebsiteID == DefaultWebsiteId
//}
