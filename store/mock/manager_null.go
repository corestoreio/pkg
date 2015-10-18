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

package mock

import (
	"github.com/corestoreio/csfw/config/scope"
	"github.com/corestoreio/csfw/store"
)

// NullManager does nothing and returns only errors.
type NullManager struct{}

func (m *NullManager) IsSingleStoreMode() bool { return false }
func (m *NullManager) HasSingleStore() bool    { return false }
func (m *NullManager) Website(r ...scope.WebsiteIDer) (*store.Website, error) {
	return nil, store.ErrWebsiteNotFound
}
func (m *NullManager) Websites() (store.WebsiteSlice, error) { return nil, store.ErrWebsiteNotFound }
func (m *NullManager) Group(r ...scope.GroupIDer) (*store.Group, error) {
	return nil, store.ErrGroupNotFound
}
func (m *NullManager) Groups() (store.GroupSlice, error) { return nil, store.ErrGroupNotFound }
func (m *NullManager) Store(r ...scope.StoreIDer) (*store.Store, error) {
	return nil, store.ErrStoreNotFound
}
func (m *NullManager) Stores() (store.StoreSlice, error)       { return nil, store.ErrStoreNotFound }
func (m *NullManager) DefaultStoreView() (*store.Store, error) { return nil, store.ErrStoreNotFound }

// NewNullManager creates a new NullManager
func NewNullManager() *NullManager {
	return &NullManager{}
}

var _ store.ManagerReader = (*NullManager)(nil)
