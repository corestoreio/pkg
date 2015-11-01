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

package config

import (
	"fmt"
	"io"
	"io/ioutil"
	"reflect"
	"sync"
	"time"

	"github.com/corestoreio/csfw/config/scope"
	"github.com/dustin/gojson"
	"golang.org/x/net/context"
)

var _ Reader = (*MockReader)(nil)
var _ ReaderPubSuber = (*MockReader)(nil)

// mockOptionFunc to initialize the NewMockReader
type mockOptionFunc func(*MockReader)

// MockReader used for testing. Contains functions which will be called in the
// appropriate methods of interface config.Reader.
type MockReader struct {
	mu              sync.RWMutex
	mv              MockPV
	String          func(path string) (string, error)
	Bool            func(path string) (bool, error)
	F64             func(path string) (float64, error)
	Int             func(path string) (int, error)
	Time            func(path string) (time.Time, error)
	SubscriptionID  int
	SubscriptionErr error
}

// MockPV is a required type for an option function. PV = path => value.
// This map[string]interface{} is protected by a mutex.
type MockPV map[string]interface{}

// MockPathScopeDefault creates for testing a fully qualified path for the
// default scope and a path string (a/b/c)
func MockPathScopeDefault(path string) string {
	return scope.StrDefault.FQPathInt64(0, path)
}

// MockPathScopeWebsite creates for testing a fully qualified path for the
// website scope from a Scope ID and a path string (a/b/c)
func MockPathScopeWebsite(id int64, path string) string {
	return scope.StrWebsites.FQPathInt64(id, path)
}

// MockPathScopeStore creates for testing a fully qualified path for the
// store scope from a Scope ID and a path string (a/b/c)
func MockPathScopeStore(id int64, path string) string {
	return scope.StrStores.FQPathInt64(id, path)
}

// WithMockString returns a function which can be used in the NewMockReader().
// Your function returns a string value from a given path.
func WithMockString(f func(path string) (string, error)) mockOptionFunc {
	return func(mr *MockReader) { mr.String = f }
}

// WithMockBool returns a function which can be used in the NewMockReader().
// Your function returns a bool value from a given path.
func WithMockBool(f func(path string) (bool, error)) mockOptionFunc {
	return func(mr *MockReader) { mr.Bool = f }
}

// WithMockFloat64 returns a function which can be used in the NewMockReader().
// Your function returns a float64 value from a given path.
func WithMockFloat64(f func(path string) (float64, error)) mockOptionFunc {
	return func(mr *MockReader) { mr.F64 = f }
}

// WithMockInt returns a function which can be used in the NewMockReader().
// Your function returns an int value from a given path.
func WithMockInt(f func(path string) (int, error)) mockOptionFunc {
	return func(mr *MockReader) { mr.Int = f }
}

// WithMockTime returns a function which can be used in the NewMockReader().
// Your function returns a Time value from a given path.
func WithMockTime(f func(path string) (time.Time, error)) mockOptionFunc {
	return func(mr *MockReader) {
		mr.Time = f
	}
}

// WithMockValues lets you define a map of path and its values.
// Key is the fully qualified configuration path and value is the value.
// Value must be of the same type as returned by the functions.
func WithMockValues(pathValues MockPV) mockOptionFunc {
	return func(mr *MockReader) {
		mr.mu.Lock()
		mr.mv = pathValues
		mr.mu.Unlock()
	}
}

// WithMockValuesJSON same as WithMockValues but reads data from an io.Reader
// so you can read config from a JSON file.
//
func WithMockValuesJSON(r io.Reader) mockOptionFunc {
	rawJson, err := ioutil.ReadAll(r)
	if err != nil {
		panic(err)
	}

	var pathValues MockPV
	err = json.Unmarshal(rawJson, &pathValues)
	if err != nil {
		panic(err)
	}
	return func(mr *MockReader) {
		mr.mu.Lock()
		mr.mv = pathValues
		mr.mu.Unlock()
	}
}

// NewContextMockReader adds a MockReader to a context.
func NewContextMockReader(ctx context.Context, opts ...mockOptionFunc) context.Context {
	return context.WithValue(ctx, ctxKeyReader, NewMockReader(opts...))
}

// NewMockReader creates a new MockReader used in testing.
// Allows you to set different options duration creation or you can
// set the struct fields afterwards.
func NewMockReader(opts ...mockOptionFunc) *MockReader {
	mr := &MockReader{}
	for _, opt := range opts {
		opt(mr)
	}
	return mr
}

// UpdateValues adds or overwrites the internal path => value map.
func (mr *MockReader) UpdateValues(pathValues MockPV) {
	mr.mu.Lock()
	for k, v := range pathValues {
		mr.mv[k] = v
	}
	mr.mu.Unlock()
}

func (mr *MockReader) hasVal(path string) bool {
	mr.mu.Lock()
	defer mr.mu.Unlock()
	_, ok := mr.mv[path]
	return ok
}

func (mr *MockReader) getVal(path string) interface{} {
	mr.mu.Lock()
	defer mr.mu.Unlock()
	v := mr.mv[path]
	v = indirect(v)
	return v
}

func (mr *MockReader) valString(path string) (string, error) {
	switch s := mr.getVal(path).(type) {
	case string:
		return s, nil
	case nil:
		return "", nil
	case fmt.Stringer:
		return s.String(), nil
	default:
		return "", fmt.Errorf("Unable to Cast %#v to string", s)
	}
}

// GetString returns a string value
func (mr *MockReader) GetString(opts ...ArgFunc) (string, error) {
	path := mustNewArg(opts...).scopePath()
	switch {
	case mr.hasVal(path):
		return mr.valString(path)
	case mr.String != nil:
		return mr.String(path)
	default:
		return "", ErrKeyNotFound
	}
}

func (mr *MockReader) valBool(path string) (bool, error) {
	switch b := mr.getVal(path).(type) {
	case bool:
		return b, nil
	case int, int8, int16, int32, int64:
		if b != 0 {
			return true, nil
		}
		return false, nil
	case nil:
		return false, nil
	default:
		return false, fmt.Errorf("Unable to Cast %#v to bool", b)
	}
}

// GetBool returns a bool value
func (mr *MockReader) GetBool(opts ...ArgFunc) (bool, error) {
	path := mustNewArg(opts...).scopePath()
	switch {
	case mr.hasVal(path):
		return mr.valBool(path)
	case mr.Bool != nil:
		return mr.Bool(path)
	default:
		return false, ErrKeyNotFound
	}
}

func (sr *MockReader) valFloat64(path string) (float64, error) {
	switch s := sr.getVal(path).(type) {
	case float64:
		return s, nil
	case float32:
		return float64(s), nil
	default:
		return 0.0, fmt.Errorf("Unable to Cast %#v to float", s)
	}
}

// GetFloat64 returns a float64 value
func (sr *MockReader) GetFloat64(opts ...ArgFunc) (float64, error) {
	path := mustNewArg(opts...).scopePath()
	switch {
	case sr.hasVal(path):
		return sr.valFloat64(path)
	case sr.F64 != nil:
		return sr.F64(path)
	default:
		return 0.0, ErrKeyNotFound
	}
}

func (sr *MockReader) valInt(path string) (int, error) {
	switch s := sr.getVal(path).(type) {
	case int:
		return s, nil
	case nil:
		return 0, nil
	default:
		return 0, fmt.Errorf("Unable to Cast %#v to int", s)
	}
}

// GetInt returns an integer value
func (sr *MockReader) GetInt(opts ...ArgFunc) (int, error) {
	path := mustNewArg(opts...).scopePath()
	switch {
	case sr.hasVal(path):
		return sr.valInt(path)
	case sr.Int != nil:
		return sr.Int(path)
	default:
		return 0, ErrKeyNotFound
	}
}

func (sr *MockReader) valDateTime(path string) (time.Time, error) {
	switch s := sr.getVal(path).(type) {
	case time.Time:
		return s, nil
	default:
		return time.Time{}, fmt.Errorf("Unable to Cast %#v to Time\n", s)
	}
}

// GetDateTime returns a time value
func (sr *MockReader) GetDateTime(opts ...ArgFunc) (time.Time, error) {
	path := mustNewArg(opts...).scopePath()
	switch {
	case sr.hasVal(path):
		return sr.valDateTime(path)
	case sr.Time != nil:
		return sr.Time(path)
	default:
		return time.Time{}, ErrKeyNotFound
	}
}

// Subscribe returns the before applied SubscriptionID and SubscriptionErr
// Does not start any underlying Goroutines.
func (sr *MockReader) Subscribe(path string, s MessageReceiver) (subscriptionID int, err error) {
	return sr.SubscriptionID, sr.SubscriptionErr
}

// NewScoped creates a new config.ScopedReader which uses the underlying
// mocked paths and values.
func (sr *MockReader) NewScoped(websiteID, groupID, storeID int64) ScopedReader {
	return newScopedManager(sr, websiteID, groupID, storeID)
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
