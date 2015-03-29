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
	DefaultGroupId int64 = 0
)

// GroupIndex used for iota and for not mixing up indexes
type GroupIndex int

var (
	ErrStoreGroupNotFound = errors.New("Store Group not found")
	groupCollection       GroupSlice
)

// GetGroup uses a GroupIndex to return a group or an error.
// One should not modify the group object.
func GetGroup(i GroupIndex) (*Group, error) {
	if int(i) < len(groupCollection) {
		return groupCollection[i], nil
	}
	return nil, ErrStoreGroupNotFound
}

// GetGroups returns a copy of the main slice of store groups.
// One should not modify the slice and its content.
// @todo $withDefault bool
func GetGroups() GroupSlice {
	return groupCollection
}

// Load uses a dbr session to load all data from the core_store_group table into the current slice.
// The variadic 2nd argument can be a call back function to manipulate the select.
// Additional columns or joins cannot be added. This method receiver should only be used in development.
// @see app/code/Magento/Store/Model/Resource/Group/Collection.php::_beforeLoad()
func (s *GroupSlice) Load(dbrSess dbr.SessionRunner, cbs ...csdb.DbrSelectCb) (int, error) {
	return loadSlice(dbrSess, TableGroup, &(*s), append(cbs, func(sb *dbr.SelectBuilder) *dbr.SelectBuilder {
		return sb.OrderBy("main_table.name ASC")
	})...)
}

/*
	@todo implement Magento\Store\Model\Group
*/
