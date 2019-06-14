// Copyright 2015-present, Cyrill @ Schumacher.fm and the CoreStore contributors
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
	"sync/atomic"

	"github.com/corestoreio/pkg/store"
	"github.com/corestoreio/pkg/store/scope"
)

var _ store.Finder = (*Find)(nil)

// Find implements interface store.Finder for mocking in tests. Thread safe.
type Find struct {
	DefaultStoreIDFn      func(runMode scope.TypeID) (websiteID, storeID uint32, err error)
	defaultStoreIDInvoked int32

	StoreIDbyCodeFn      func(runMode scope.TypeID, storeCode string) (websiteID, storeID uint32, err error)
	storeIDbyCodeInvoked int32
}

// NewFindDefaultStoreID creates a new closure for the function DefaultStoreID.
// The last variadic argument allows to append the other NewFind*() function.
func NewDefaultStoreID(websiteID, storeID uint32, err error, fs ...*Find) *Find {
	f := func(runMode scope.TypeID, storeCode string) (uint32, uint32, error) {
		return 0, 0, nil
	}
	if len(fs) == 1 && fs[0] != nil {
		f = fs[0].StoreIDbyCodeFn
	}
	return &Find{
		DefaultStoreIDFn: func(runMode scope.TypeID) (uint32, uint32, error) {
			return websiteID, storeID, err
		},
		StoreIDbyCodeFn: f,
	}
}

// NewStoreIDbyCode creates a new closure for the function StoreIDbyCode.
// The last variadic argument allows to append the other NewFind*() function.
func NewStoreIDbyCode(websiteID, storeID uint32, err error, fs ...*Find) *Find {
	f := func(runMode scope.TypeID) (uint32, uint32, error) {
		return 0, 0, nil
	}
	if len(fs) == 1 && fs[0] != nil {
		f = fs[0].DefaultStoreIDFn
	}
	return &Find{
		DefaultStoreIDFn: f,
		StoreIDbyCodeFn: func(runMode scope.TypeID, storeCode string) (uint32, uint32, error) {
			return websiteID, storeID, err
		},
	}
}

func (s *Find) DefaultStoreID(runMode scope.TypeID) (websiteID, storeID uint32, err error) {
	atomic.AddInt32(&s.defaultStoreIDInvoked, 1)
	return s.DefaultStoreIDFn(runMode)
}

// DefaultStoreIDInvoked returns the number of DefaultStoreID() call invocations.
func (s *Find) DefaultStoreIDInvoked() int {
	return int(atomic.LoadInt32(&s.defaultStoreIDInvoked))
}

func (s *Find) StoreIDbyCode(runMode scope.TypeID, storeCode string) (websiteID, storeID uint32, err error) {
	atomic.AddInt32(&s.storeIDbyCodeInvoked, 1)
	return s.StoreIDbyCodeFn(runMode, storeCode)
}

// StoreIDbyCodeInvoked returns the number of StoreIDbyCode() call invocations.
func (s *Find) StoreIDbyCodeInvoked() int {
	return int(atomic.LoadInt32(&s.storeIDbyCodeInvoked))
}
