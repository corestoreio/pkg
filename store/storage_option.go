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

package store

import (
	"github.com/corestoreio/csfw/config"
	"github.com/corestoreio/csfw/storage/dbr"
)

// StorageOption option func for NewStorage()
type StorageOption func(*storage)

// SetStorageWebsites adds the TableWebsiteSlice to the Storage. By default, the slice is nil.
func SetStorageWebsites(tws ...*TableWebsite) StorageOption {
	return func(s *storage) { s.websites = TableWebsiteSlice(tws) }
}

// SetStorageGroups adds the TableGroupSlice to the Storage. By default, the slice is nil.
func SetStorageGroups(tgs ...*TableGroup) StorageOption {
	return func(s *storage) { s.groups = TableGroupSlice(tgs) }
}

// SetStorageStores adds the TableStoreSlice to the Storage. By default, the slice is nil.
func SetStorageStores(tss ...*TableStore) StorageOption {
	return func(s *storage) { s.stores = TableStoreSlice(tss) }
}

// WithStorageConfig sets the configuration Getter. Optional.
// If not set all websites, groups and stores have a nil Config.
func WithStorageConfig(cr config.Getter) StorageOption {
	return func(s *storage) { s.cr = cr }
}

// WithDatabaseInit triggers the ReInit function to load the data from the
// database.
func WithDatabaseInit(dbrSess dbr.SessionRunner, cbs ...dbr.SelectCb) StorageOption {
	return func(s *storage) {
		if err := s.ReInit(dbrSess, cbs...); err != nil {
			s.MultiErr = s.AppendErrors(err)
		}
	}
}
