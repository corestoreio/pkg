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

// NullService does nothing and returns only errors.
type NullService struct{}

func (m *NullService) IsSingleStoreMode() bool { return false }
func (m *NullService) HasSingleStore() bool    { return false }
func (m *NullService) Website(r ...scope.WebsiteIDer) (*store.Website, error) {
	return nil, store.ErrWebsiteNotFound
}
func (m *NullService) Websites() (store.WebsiteSlice, error) { return nil, store.ErrWebsiteNotFound }
func (m *NullService) Group(r ...scope.GroupIDer) (*store.Group, error) {
	return nil, store.ErrGroupNotFound
}
func (m *NullService) Groups() (store.GroupSlice, error) { return nil, store.ErrGroupNotFound }
func (m *NullService) Store(r ...scope.StoreIDer) (*store.Store, error) {
	return nil, store.ErrStoreNotFound
}
func (m *NullService) Stores() (store.StoreSlice, error)       { return nil, store.ErrStoreNotFound }
func (m *NullService) DefaultStoreView() (*store.Store, error) { return nil, store.ErrStoreNotFound }
func (m *NullService) RequestedStore(scope.Option) (*store.Store, error) {
	return nil, store.ErrStoreNotFound
}

// NewNullService creates a new NullService
func NewNullService() *NullService {
	return &NullService{}
}

var _ store.Reader = (*NullService)(nil)
