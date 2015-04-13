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
)

const (
	DefaultGroupId int64 = 0
)

type (
	// GroupIndex used for iota and for not mixing up indexes
	GroupIndex      uint
	GroupIndexIDMap map[int64]GroupIndex
	// GroupBucket contains two maps for faster retrieving of the store index and the store collection
	// Only used in generated code. Implements interface GroupGetter.
	GroupBucket struct {
		// store collection
		s TableGroupSlice
		// i map by store_id
		i GroupIndexIDMap
	}
)

var (
	ErrGroupNotFound = errors.New("Store Group not found")
)

// NewGroupBucket returns a new pointer to a GroupBucket.
func NewGroupBucket(s TableGroupSlice, i GroupIndexIDMap) *GroupBucket {
	// @todo idea if i and c is nil generate them from s.
	return &GroupBucket{
		i: i,
		s: s,
	}
}

// ByID uses the database store id to return a TableGroup struct.
func (s *GroupBucket) ByID(id int64) (*TableGroup, error) {
	if i, ok := s.i[id]; ok && id < int64(s.s.Len()) {
		return s.s[i], nil
	}
	return nil, ErrGroupNotFound
}

// ByIndex returns a TableGroup struct using the slice index
func (s *GroupBucket) ByIndex(i GroupIndex) (*TableGroup, error) {
	if int(i) < s.s.Len() {
		return s.s[i], nil
	}
	return nil, ErrGroupNotFound
}

// ByIndex returns the TableGroupSlice
func (s *GroupBucket) Group() TableGroupSlice { return s.s }

// Load uses a dbr session to load all data from the core_store_group table into the current slice.
// The variadic 2nd argument can be a call back function to manipulate the select.
// Additional columns or joins cannot be added. This method receiver should only be used in development.
// @see app/code/Magento/Store/Model/Resource/Group/Collection.php::_beforeLoad()
func (s *TableGroupSlice) Load(dbrSess dbr.SessionRunner, cbs ...csdb.DbrSelectCb) (int, error) {
	return loadSlice(dbrSess, TableIndexGroup, &(*s), append(cbs, func(sb *dbr.SelectBuilder) *dbr.SelectBuilder {
		return sb.OrderBy("main_table.name ASC")
	})...)
}

// Len returns the length
func (s TableGroupSlice) Len() int { return len(s) }

/*
	@todo implement Magento\Store\Model\Group
*/
