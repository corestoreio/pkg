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
	"fmt"
	"sort"
	"sync"

	"github.com/corestoreio/errors"
	"github.com/corestoreio/pkg/store/scope"
	"github.com/corestoreio/pkg/util/bufferpool"
)

// MockWrite used for testing when writing configuration values.
// deprecated no replacement
type MockWrite struct {
	// WriteError gets always returned by MockWrite
	WriteError error
	// ArgPath will be set after calling write to export the config path.
	// Values you enter here will be overwritten when calling MockWrite
	ArgPath string
	// ArgValue contains the written data
	ArgValue []byte
}

// MockWrite writes to a black hole, may return an error
func (w *MockWrite) Write(p Path, value []byte) error {
	w.ArgPath = p.String()
	w.ArgValue = value
	return w.WriteError
}

// invocations represents a list containing the fully qualified configuration
// path and number of invocations. This type has attached some helper functions.
type invocations map[string]int

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
		p[i] = k
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
		p, err := SplitFQ(k)
		if err != nil {
			panic(fmt.Sprintf("[cfgmock] Path: %q with error: %+v", k, err))
		}
		ids[i] = p.ScopeID
		i++
	}
	sort.Sort(ids)
	return ids
}

// Mock used for testing. Contains functions which will be called in the
// appropriate methods of interface Getter. Field DB has precedence over
// the applied functions.
type Mock struct {
	Storage Storager
	mu      sync.Mutex
	// GetFn can be set optionally. If set then Storage will be ignored.
	GetFn       func(p Path) (v Value)
	invocations invocations // contains path and count of how many times the typed function has been called

	SubscribeFn      func(string, MessageReceiver) (subscriptionID int, err error)
	SubscribeInvokes int32
}

// MockPathValue is a required type for an option function. PV = path => value. This
// map[string]interface{} is protected by a mutex.
type MockPathValue map[string]string

func (pv MockPathValue) set(db Storager) {
	for fq, v := range pv {
		if err := db.Set(0, fq, []byte(v)); err != nil {
			panic(err)
		}
	}
}

// GoString creates a sorted Go syntax valid map representation. This function
// panics if it fails to write to the internal buffer. Panicing permitted here
// because this function is only used in testing.
func (pv MockPathValue) GoString() string {
	keys := make(sort.StringSlice, len(pv))
	i := 0
	for k := range pv {
		keys[i] = k
		i++
	}
	keys.Sort()

	buf := bufferpool.Get()
	defer bufferpool.Put(buf)
	if _, err := buf.WriteString("config.MockPathValue{\n"); err != nil {
		panic(err)
	}

	for _, p := range keys {
		if _, err := fmt.Fprintf(buf, "%q: %#v,\n", p, pv[p]); err != nil {
			panic(err)
		}
	}
	if _, err := buf.WriteRune('}'); err != nil {
		panic(err)
	}
	return buf.String()
}

// NewMock creates a new mocked Mock for testing usage. Initializes a
// simple in memory key/value storage.
func NewMock(pvs ...MockPathValue) *Mock {
	mr := &Mock{
		Storage: NewInMemoryStore(),
	}
	if len(pvs) > 0 {
		for _, pv := range pvs {
			pv.set(mr.Storage)
		}
	}
	return mr
}

// AllInvocations returns all called paths and increases the counter for
// duplicated paths.
func (s *Mock) AllInvocations() invocations {
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

// UpdateValues adds or overwrites the internal path => value map.
func (s *Mock) UpdateValues(pv MockPathValue) {
	pv.set(s.Storage)
}

// Value looks up a configuration value for a given path.
func (s *Mock) Value(p Path) Value {
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.invocations == nil {
		s.invocations = make(invocations)
	}
	ps := p.String()
	s.invocations[ps]++

	if s.GetFn != nil {
		return s.GetFn(p)
	}

	vb, ok, err := s.Storage.Value(0, p.String())
	var found uint8
	if ok {
		found = valFoundYes
	}
	return Value{
		Path:    p,
		data:    vb,
		found:   found,
		lastErr: errors.WithStack(err),
	}
}

// Invokes returns statistics about invocations
func (s *Mock) Invokes() invocations {
	return s.invocations
}

// Subscribe returns the before applied SubscriptionID and SubscriptionErr
// Does not start any underlying Goroutines.
func (s *Mock) Subscribe(path string, mr MessageReceiver) (subscriptionID int, err error) {
	s.SubscribeInvokes++
	return s.SubscribeFn(path, mr)
}

// NewScoped creates a new ScopedReader which uses the underlying
// mocked paths and values.
func (s *Mock) NewScoped(websiteID, storeID int64) Scoped {
	return NewScoped(s, websiteID, storeID)
}
