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

package storemock

import (
	"fmt"

	"github.com/corestoreio/csfw/config"
	"github.com/corestoreio/csfw/storage/dbr"
	"github.com/corestoreio/csfw/store"
	"github.com/corestoreio/csfw/store/scope"
	"github.com/corestoreio/csfw/util/errors"
)

// RequestedStoreAU satisfies the interface store.Requester and returns always
// the Australian store.
type RequestedStoreAU struct {
	// Getter will be applied to the function NewStoreAU()
	config.Getter
	// ScpOpt contains the scope which has been passed into the function
	// RequestedStore() for later inspection.
	ScpOpt scope.Option
}

func (rs *RequestedStoreAU) RequestedStore(so scope.Option) (*store.Store, error) {
	rs.ScpOpt = so
	return NewStoreAU(rs.Getter)
}

// NewRequestedStoreAU creates a new type. Only one argument may be passed.
func NewRequestedStoreAU(cgs ...config.Getter) *RequestedStoreAU {
	var cg config.Getter
	if len(cgs) > 0 {
		cg = cgs[0]
	}
	return &RequestedStoreAU{
		Getter: cg,
	}
}

// NewStoreAU creates a new Store with an attached config.
// Store ID 5, Code "au"; Website ID 2, Code "oz"; GroupID 3.
func NewStoreAU(cg config.Getter) (*store.Store, error) {

	st, err := store.NewStore(
		&store.TableStore{StoreID: 5, Code: dbr.NewNullString("au"), WebsiteID: 2, GroupID: 3, Name: "Australia", SortOrder: 10, IsActive: true},
		&store.TableWebsite{WebsiteID: 2, Code: dbr.NewNullString("oz"), Name: dbr.NewNullString("OZ"), SortOrder: 20, DefaultGroupID: 3, IsDefault: dbr.NewNullBool(false)},
		&store.TableGroup{GroupID: 3, WebsiteID: 2, Name: "Australia", RootCategoryID: 2, DefaultStoreID: 5},
		store.WithStoreConfig(cg),
	)

	return st, errors.Wrap(err, "[storemock] NewStoreAU")
}

// MustNewStoreAU creates a new Store with an attached config.
// Store ID 5, Code "au"; Website ID 2, Code "oz"; GroupID 3.
func MustNewStoreAU(cg config.Getter) *store.Store {
	st, err := NewStoreAU(cg)
	if err != nil {
		panic(fmt.Sprintf("%+v", err))
	}
	return st
}
