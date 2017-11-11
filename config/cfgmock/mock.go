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

package cfgmock

import (
	"fmt"
	"reflect"
	"sort"
	"sync"
	"time"

	"github.com/corestoreio/cspkg/config"
	"github.com/corestoreio/cspkg/config/cfgpath"
	"github.com/corestoreio/cspkg/store/scope"
	"github.com/corestoreio/cspkg/util/bufferpool"
	"github.com/corestoreio/cspkg/util/conv"
	"github.com/corestoreio/errors"
)

// keyNotFound for performance and allocs reasons in benchmarks to test properly
// the cfg* code and not the configuration Service. The NotFound error has been
// hard coded which does not record the position where the error happens. We can
// maybe add the path which was not found but that will trigger 2 allocs because
// of the sprintf ... which could be bypassed with a bufferpool ;-)
type keyNotFound struct{}

func (a keyNotFound) Error() string  { return "[cfgmock] Get() Path not found" }
func (a keyNotFound) NotFound() bool { return true }

// Write used for testing when writing configuration values.
type Write struct {
	// WriteError gets always returned by Write
	WriteError error
	// ArgPath will be set after calling write to export the config path.
	// Values you enter here will be overwritten when calling Write
	ArgPath string
	// ArgValue contains the written data
	ArgValue interface{}
}

// Write writes to a black hole, may return an error
func (w *Write) Write(p cfgpath.Path, v interface{}) error {
	w.ArgPath = p.String()
	w.ArgValue = v
	return w.WriteError
}

// Invocations represents a list containing the fully qualified configuration
// path and number of invocations. This type has attached some helper functions.
type Invocations map[string]int

// Sum returns the total calls for the current method receiver type. E.g. All
// calls to String() with different or same paths.
func (iv Invocations) Sum() int {
	s := 0
	for _, v := range iv {
		s += v
	}
	return s
}

// PathCount returns the number of different paths.
func (iv Invocations) PathCount() int {
	return len(iv)
}

// Paths returns all paths in sorted ascending order.
func (iv Invocations) Paths() []string {
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
func (iv Invocations) ScopeIDs() scope.TypeIDs {
	ids := make(scope.TypeIDs, len(iv))
	i := 0
	for k := range iv {
		p, err := cfgpath.SplitFQ(k)
		if err != nil {
			panic(fmt.Sprintf("[cfgmock] Path: %q with error: %+v", k, err))
		}
		ids[i] = p.ScopeID
		i++
	}
	sort.Sort(ids)
	return ids
}

// Service used for testing. Contains functions which will be called in the
// appropriate methods of interface config.Getter. Field DB has precedence over
// the applied functions.
type Service struct {
	Storage          config.Storager
	mu               sync.Mutex
	ByteFn           func(path string) ([]byte, error)
	byteInvokes      Invocations // contains path and count of how many times the typed function has been called
	StringFn         func(path string) (string, error)
	stringInvokes    Invocations
	BoolFn           func(path string) (bool, error)
	boolInvokes      Invocations
	Float64Fn        func(path string) (float64, error)
	float64Invokes   Invocations
	IntFn            func(path string) (int, error)
	intInvokes       Invocations
	TimeFn           func(path string) (time.Time, error)
	timeInvokes      Invocations
	DurationFn       func(path string) (time.Duration, error)
	durationInvokes  Invocations
	SubscribeFn      func(cfgpath.Route, config.MessageReceiver) (subscriptionID int, err error)
	SubscribeInvokes int32
}

// PathValue is a required type for an option function. PV = path => value. This
// map[string]interface{} is protected by a mutex.
type PathValue map[string]interface{}

func (pv PathValue) set(db config.Storager) {
	for fq, v := range pv {
		p, err := cfgpath.SplitFQ(fq)
		if err != nil {
			panic(err)
		}
		if err := db.Set(p, v); err != nil {
			panic(err)
		}
	}
}

// GoString creates a sorted Go syntax valid map representation. This function
// panics if it fails to write to the internal buffer. Panicing permitted here
// because this function is only used in testing.
func (pv PathValue) GoString() string {
	keys := make(sort.StringSlice, len(pv))
	i := 0
	for k := range pv {
		keys[i] = k
		i++
	}
	keys.Sort()

	buf := bufferpool.Get()
	defer bufferpool.Put(buf)
	if _, err := buf.WriteString("cfgmock.PathValue{\n"); err != nil {
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

// NewService creates a new mocked Service for testing usage. Initializes a
// simple in memory key/value storage.
func NewService(pvs ...PathValue) *Service {
	mr := &Service{
		Storage: config.NewInMemoryStore(),
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
func (s *Service) AllInvocations() Invocations {
	s.mu.Lock()
	defer s.mu.Unlock()
	ret := make(Invocations)
	add := func(iv Invocations) {
		for k, v := range iv {
			ret[k] += v
		}
	}
	add(s.byteInvokes)
	add(s.stringInvokes)
	add(s.boolInvokes)
	add(s.float64Invokes)
	add(s.intInvokes)
	add(s.timeInvokes)
	add(s.durationInvokes)
	return ret
}

// UpdateValues adds or overwrites the internal path => value map.
func (s *Service) UpdateValues(pv PathValue) {
	pv.set(s.Storage)
}

func (s *Service) hasVal(p cfgpath.Path) bool {
	if s.Storage == nil {
		return false
	}
	v, err := s.Storage.Get(p)
	if err != nil && !errors.IsNotFound(err) {
		println("Mock.Service.hasVal error:", err.Error(), "path", p.String())
	}
	return v != nil && err == nil
}

func (s *Service) getVal(p cfgpath.Path) interface{} {
	v, err := s.Storage.Get(p)
	if err != nil && !errors.IsNotFound(err) {
		println("Mock.Service.getVal error:", err.Error(), "path", p.String())
		return nil
	}
	v = indirect(v)
	return v
}

// Byte returns a byte slice value
func (s *Service) Byte(p cfgpath.Path) ([]byte, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.byteInvokes == nil {
		s.byteInvokes = make(Invocations)
	}
	ps := p.String()
	s.byteInvokes[ps]++

	switch {
	case s.hasVal(p):
		return conv.ToByteE(s.getVal(p))
	case s.ByteFn != nil:
		return s.ByteFn(ps)
	default:
		return nil, keyNotFound{}
	}
}

// ByteInvokes returns the number of Byte() invocations.
func (s *Service) ByteInvokes() Invocations {
	return s.byteInvokes
}

// String returns a string value
func (s *Service) String(p cfgpath.Path) (string, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.stringInvokes == nil {
		s.stringInvokes = make(Invocations)
	}
	ps := p.String()
	s.stringInvokes[ps]++

	switch {
	case s.hasVal(p):
		return conv.ToStringE(s.getVal(p))
	case s.StringFn != nil:
		return s.StringFn(ps)
	default:
		return "", keyNotFound{}
	}
}

// StringInvokes returns the number of String() invocations.
func (s *Service) StringInvokes() Invocations {
	return s.stringInvokes
}

// Bool returns a bool value
func (s *Service) Bool(p cfgpath.Path) (bool, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.boolInvokes == nil {
		s.boolInvokes = make(Invocations)
	}
	ps := p.String()
	s.boolInvokes[ps]++

	switch {
	case s.hasVal(p):
		return conv.ToBoolE(s.getVal(p))
	case s.BoolFn != nil:
		return s.BoolFn(ps)
	default:
		return false, keyNotFound{}
	}
}

// BoolInvokes returns the number of Bool() invocations.
func (s *Service) BoolInvokes() Invocations {
	return s.boolInvokes
}

// Float64 returns a float64 value
func (s *Service) Float64(p cfgpath.Path) (float64, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.float64Invokes == nil {
		s.float64Invokes = make(Invocations)
	}
	ps := p.String()
	s.float64Invokes[ps]++

	switch {
	case s.hasVal(p):
		return conv.ToFloat64E(s.getVal(p))
	case s.Float64Fn != nil:
		return s.Float64Fn(ps)
	default:
		return 0.0, keyNotFound{}
	}
}

// Float64Invokes returns the number of Float64() invocations.
func (s *Service) Float64Invokes() Invocations {
	return s.float64Invokes
}

// Int returns an integer value
func (s *Service) Int(p cfgpath.Path) (int, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.intInvokes == nil {
		s.intInvokes = make(Invocations)
	}
	ps := p.String()
	s.intInvokes[ps]++

	switch {
	case s.hasVal(p):
		return conv.ToIntE(s.getVal(p))
	case s.IntFn != nil:
		return s.IntFn(ps)
	default:
		return 0, keyNotFound{}
	}
}

// IntInvokes returns the number of Int() invocations.
func (s *Service) IntInvokes() Invocations {
	return s.intInvokes
}

// Time returns a time value
func (s *Service) Time(p cfgpath.Path) (time.Time, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.timeInvokes == nil {
		s.timeInvokes = make(Invocations)
	}
	ps := p.String()
	s.timeInvokes[ps]++

	switch {
	case s.hasVal(p):
		return conv.ToTimeE(s.getVal(p))
	case s.TimeFn != nil:
		return s.TimeFn(ps)
	default:
		return time.Time{}, keyNotFound{}
	}
}

// TimeInvokes returns the number of Time() invocations.
func (s *Service) TimeInvokes() Invocations {
	return s.timeInvokes
}

// Duration returns a duration value or a NotFound error.
func (s *Service) Duration(p cfgpath.Path) (time.Duration, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.durationInvokes == nil {
		s.durationInvokes = make(Invocations)
	}
	ps := p.String()
	s.durationInvokes[ps]++

	switch {
	case s.hasVal(p):
		return conv.ToDurationE(s.getVal(p))
	case s.DurationFn != nil:
		return s.DurationFn(ps)
	default:
		return 0, keyNotFound{}
	}
}

// DurationInvokes returns the number of Duration() invocations.
func (s *Service) DurationInvokes() Invocations {
	return s.durationInvokes
}

// Subscribe returns the before applied SubscriptionID and SubscriptionErr
// Does not start any underlying Goroutines.
func (s *Service) Subscribe(p cfgpath.Route, mr config.MessageReceiver) (subscriptionID int, err error) {
	s.SubscribeInvokes++
	return s.SubscribeFn(p, mr)
}

// NewScoped creates a new config.ScopedReader which uses the underlying
// mocked paths and values.
func (s *Service) NewScoped(websiteID, storeID int64) config.Scoped {
	return config.NewScoped(s, websiteID, storeID)
}

// From html/template/content.go
// Copyright 2011 The Go Authors. All rights reserved.
// indirect returns the value, after dereferencing as many times
// as necessary to reach the base type (or nil).
func indirect(a interface{}) interface{} {
	if a == nil {
		return nil
	}
	if t := reflect.TypeOf(a); t.Kind() != reflect.Ptr {
		// Avoid creating a reflect.Value if it's not a pointer.
		return a
	}
	v := reflect.ValueOf(a)
	for v.Kind() == reflect.Ptr && !v.IsNil() {
		v = v.Elem()
	}
	return v.Interface()
}
