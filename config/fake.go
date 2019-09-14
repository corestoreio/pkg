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

package config

import (
	"sort"
	"sync"

	"github.com/corestoreio/errors"
	"github.com/corestoreio/pkg/store/scope"
)

// FakeWrite used for testing when writing configuration values.
// deprecated no replacement
type FakeWrite struct {
	// WriteError gets always returned by FakeWrite
	WriteError error
	// ArgPath will be set after calling write to export the config path.
	// Values you enter here will be overwritten when calling FakeWrite
	ArgPath string
	// ArgValue contains the written data
	ArgValue []byte
}

// Set writes to a black hole, may return an error
func (w *FakeWrite) Set(p *Path, value []byte) error {
	w.ArgPath = p.String()
	w.ArgValue = value
	return w.WriteError
}

// invocations represents a list containing the fully qualified configuration
// path and number of invocations. This type has attached some helper functions.
type invocations map[Path]int

// Sum returns the total calls for the current method receiver type. E.g. All
// calls to String() with different or same paths.
func (iv invocations) Sum() int {
	s := 0
	for _, v := range iv {
		s += v
	}
	return s
}

// PathCount returns the number of different paths.
func (iv invocations) PathCount() int {
	return len(iv)
}

// Paths returns all paths in sorted ascending order.
func (iv invocations) Paths() []string {
	p := make([]string, len(iv))
	i := 0
	for k := range iv {
		p[i] = k.String()
		i++
	}
	sort.Strings(p)
	return p
}

// ScopeIDs extracts all scope.TypeID from all paths in sorted ascending order.
// The returned length of scope.TypeIDs is equal to PathCount().
func (iv invocations) ScopeIDs() scope.TypeIDs {
	ids := make(scope.TypeIDs, len(iv))
	i := 0
	for k := range iv {
		ids[i] = k.ScopeID
		i++
	}
	sort.Sort(ids)
	return ids
}

// FakeService used for testing. Contains functions which will be called in the
// appropriate methods of interface Getter. Field DB has precedence over
// the applied functions.
type FakeService struct {
	storage Storager
	mu      sync.Mutex
	// GetFn can be set optionally. If set then Storage will be ignored.
	// Must return a guaranteed non-nil Value.
	GetFn       func(p *Path) (v *Value)
	invocations invocations // contains path and count of how many times the typed function has been called

	SubscribeFn      func(string, MessageReceiver) (subscriptionID int, err error)
	SubscribeInvokes int32
}

// NewFakeService creates a new mocked Service for testing usage. Initializes a
// simple in memory key/value storage.
func NewFakeService(s Storager) *FakeService {
	return &FakeService{
		storage: s,
	}
}

// AllInvocations returns all called paths and increases the counter for
// duplicated paths.
func (s *FakeService) AllInvocations() invocations {
	s.mu.Lock()
	defer s.mu.Unlock()
	ret := make(invocations)
	add := func(iv invocations) {
		for k, v := range iv {
			ret[k] += v
		}
	}
	add(s.invocations)
	return ret
}

// Get looks up a configuration value for a given path.
func (s *FakeService) Get(p *Path) *Value {
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.invocations == nil {
		s.invocations = make(invocations)
	}

	s.invocations[*p]++

	if s.GetFn != nil {
		return s.GetFn(p)
	}

	vb, ok, err := s.storage.Get(p)
	if err != nil {
		err = errors.WithStack(err)
	}
	var found uint8
	if ok {
		found = valFoundL2
	}
	return &Value{
		Path:    *p,
		data:    vb,
		found:   found,
		lastErr: err,
	}
}

// Invokes returns statistics about invocations
func (s *FakeService) Invokes() invocations {
	return s.invocations
}

// Subscribe returns the before applied SubscriptionID and SubscriptionErr
// Does not start any underlying Goroutines.
func (s *FakeService) Subscribe(path string, mr MessageReceiver) (subscriptionID int, err error) {
	s.SubscribeInvokes++
	return s.SubscribeFn(path, mr)
}

// Scoped creates a new ScopedReader which uses the underlying
// mocked paths and values.
func (s *FakeService) Scoped(websiteID, storeID uint32) Scoped {
	return makeScoped(s, websiteID, storeID)
}
