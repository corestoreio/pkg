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
	"github.com/juju/errgo"
)

const (
	DefaultGroupId int64 = 0
)

type (
	// GroupIndex used for iota and for not mixing up indexes
	GroupIndex  int
	GroupGetter interface {
		// ByID returns a GroupIndex using the GroupID. This GroupIndex identifies a group within a GroupSlice.
		ByID(id int64) (GroupIndex, error)
	}
)

var (
	ErrGroupNotFound = errors.New("Store Group not found")
	groupCollection  TableGroupSlice
	groupGetter      GroupGetter
)

func SetGroupCollection(gc TableGroupSlice) {
	if len(gc) == 0 {
		panic("StoreSlice is empty")
	}
	groupCollection = gc
}

func SetGroupGetter(g GroupGetter) {
	if g == nil {
		panic("GroupGetter cannot be nil")
	}
	groupGetter = g
}

// GetGroup uses a GroupIndex to return a group or an error.
// One should not modify the group object.
func GetGroup(i GroupIndex) (*TableGroup, error) {
	if int(i) < len(groupCollection) {
		return groupCollection[i], nil
	}
	return nil, ErrGroupNotFound
}

// GetGroups returns a copy of the main slice of store groups.
// One should not modify the slice and its content.
func GetGroups() TableGroupSlice {
	return groupCollection
}

// Load uses a dbr session to load all data from the core_store_group table into the current slice.
// The variadic 2nd argument can be a call back function to manipulate the select.
// Additional columns or joins cannot be added. This method receiver should only be used in development.
// @see app/code/Magento/Store/Model/Resource/Group/Collection.php::_beforeLoad()
func (s *TableGroupSlice) Load(dbrSess dbr.SessionRunner, cbs ...csdb.DbrSelectCb) (int, error) {
	return loadSlice(dbrSess, TableIndexGroup, &(*s), append(cbs, func(sb *dbr.SelectBuilder) *dbr.SelectBuilder {
		return sb.OrderBy("main_table.name ASC")
	})...)
}

func (s TableGroupSlice) ByID(id int64) (*TableGroup, error) {
	i, err := groupGetter.ByID(id)
	if err != nil {
		return nil, errgo.Mask(err)
	}
	return s[i], nil
}

/*
	@todo implement Magento\Store\Model\Group
*/
