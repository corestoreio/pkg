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

	"encoding/json"

	"github.com/corestoreio/csfw/config/scope"
	"golang.org/x/net/context"
)

var _ Getter = (*MockGet)(nil)
var _ Writer = (*MockWrite)(nil)
var _ GetterPubSuber = (*MockGet)(nil)

// MockWrite used for testing to testing configuration writing
type MockWrite struct {
	// WriteError gets always returned by Write
	WriteError error
	// ArgPath will be set after calling write to export the config path.
	// Values you enter here will be overwritten when calling Write
	ArgPath string
}

// Write writes to a black hole, may return an error
func (w *MockWrite) Write(o ...ArgFunc) error {
	a, err := newArg(o...)
	if err != nil {
		return err
	}
	w.ArgPath = a.scopePath()
	return w.WriteError
}

// mockOptionFunc to initialize the NewMockGetter
type mockOptionFunc func(*MockGet)

// MockGet used for testing. Contains functions which will be called in the
// appropriate methods of interface config.Getter.
type MockGet struct {
	mu              sync.RWMutex
	mv              MockPV
	FString         func(path string) (string, error)
	FBool           func(path string) (bool, error)
	FFloat64        func(path string) (float64, error)
	FInt            func(path string) (int, error)
	FTime           func(path string) (time.Time, error)
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

// WithMockString returns a function which can be used in the NewMockGetter().
// Your function returns a string value from a given path.
func WithMockString(f func(path string) (string, error)) mockOptionFunc {
	return func(mr *MockGet) { mr.FString = f }
}

// WithMockBool returns a function which can be used in the NewMockGetter().
// Your function returns a bool value from a given path.
func WithMockBool(f func(path string) (bool, error)) mockOptionFunc {
	return func(mr *MockGet) { mr.FBool = f }
}

// WithMockFloat64 returns a function which can be used in the NewMockGetter().
// Your function returns a float64 value from a given path.
func WithMockFloat64(f func(path string) (float64, error)) mockOptionFunc {
	return func(mr *MockGet) { mr.FFloat64 = f }
}

// WithMockInt returns a function which can be used in the NewMockGetter().
// Your function returns an int value from a given path.
func WithMockInt(f func(path string) (int, error)) mockOptionFunc {
	return func(mr *MockGet) { mr.FInt = f }
}

// WithMockTime returns a function which can be used in the NewMockGetter().
// Your function returns a Time value from a given path.
func WithMockTime(f func(path string) (time.Time, error)) mockOptionFunc {
	return func(mr *MockGet) {
		mr.FTime = f
	}
}

// WithMockValues lets you define a map of path and its values.
// Key is the fully qualified configuration path and value is the value.
// Value must be of the same type as returned by the functions.
func WithMockValues(pathValues MockPV) mockOptionFunc {
	return func(mr *MockGet) {
		mr.mu.Lock()
		mr.mv = pathValues
		mr.mu.Unlock()
	}
}

// WithMockValuesJSON same as WithMockValues but reads data from an io.Reader
// so you can read config from a JSON file.
//
func WithMockValuesJSON(r io.Reader) mockOptionFunc {
	rawJSON, err := ioutil.ReadAll(r)
	if err != nil {
		panic(err)
	}

	var pathValues MockPV
	err = json.Unmarshal(rawJSON, &pathValues)
	if err != nil {
		panic(err)
	}
	return func(mr *MockGet) {
		mr.mu.Lock()
		mr.mv = pathValues
		mr.mu.Unlock()
	}
}

// WithContextMockGetter adds a MockGetter to a context.
func WithContextMockGetter(ctx context.Context, opts ...mockOptionFunc) context.Context {
	return context.WithValue(ctx, ctxKeyGetter{}, NewMockGetter(opts...))
}

// NewMockGetter creates a new MockGetter used in testing.
// Allows you to set different options duration creation or you can
// set the struct fields afterwards.
func NewMockGetter(opts ...mockOptionFunc) *MockGet {
	mr := &MockGet{}
	for _, opt := range opts {
		opt(mr)
	}
	return mr
}

// UpdateValues adds or overwrites the internal path => value map.
func (mr *MockGet) UpdateValues(pathValues MockPV) {
	mr.mu.Lock()
	for k, v := range pathValues {
		mr.mv[k] = v
	}
	mr.mu.Unlock()
}

func (mr *MockGet) hasVal(path string) bool {
	mr.mu.Lock()
	defer mr.mu.Unlock()
	_, ok := mr.mv[path]
	return ok
}

func (mr *MockGet) getVal(path string) interface{} {
	mr.mu.Lock()
	defer mr.mu.Unlock()
	v := mr.mv[path]
	v = indirect(v)
	return v
}

func (mr *MockGet) valString(path string) (string, error) {
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

// String returns a string value
func (mr *MockGet) String(opts ...ArgFunc) (string, error) {
	path := mustNewArg(opts...).scopePath()
	switch {
	case mr.hasVal(path):
		return mr.valString(path)
	case mr.FString != nil:
		return mr.FString(path)
	default:
		return "", ErrKeyNotFound
	}
}

func (mr *MockGet) valBool(path string) (bool, error) {
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

// Bool returns a bool value
func (mr *MockGet) Bool(opts ...ArgFunc) (bool, error) {
	path := mustNewArg(opts...).scopePath()
	switch {
	case mr.hasVal(path):
		return mr.valBool(path)
	case mr.FBool != nil:
		return mr.FBool(path)
	default:
		return false, ErrKeyNotFound
	}
}

func (mr *MockGet) valFloat64(path string) (float64, error) {
	switch s := mr.getVal(path).(type) {
	case float64:
		return s, nil
	case float32:
		return float64(s), nil
	default:
		return 0.0, fmt.Errorf("Unable to Cast %#v to float", s)
	}
}

// Float64 returns a float64 value
func (mr *MockGet) Float64(opts ...ArgFunc) (float64, error) {
	path := mustNewArg(opts...).scopePath()
	switch {
	case mr.hasVal(path):
		return mr.valFloat64(path)
	case mr.FFloat64 != nil:
		return mr.FFloat64(path)
	default:
		return 0.0, ErrKeyNotFound
	}
}

func (mr *MockGet) valInt(path string) (int, error) {
	switch s := mr.getVal(path).(type) {
	case int:
		return s, nil
	case nil:
		return 0, nil
	default:
		return 0, fmt.Errorf("Unable to Cast %#v to int", s)
	}
}

// Int returns an integer value
func (mr *MockGet) Int(opts ...ArgFunc) (int, error) {
	path := mustNewArg(opts...).scopePath()
	switch {
	case mr.hasVal(path):
		return mr.valInt(path)
	case mr.FInt != nil:
		return mr.FInt(path)
	default:
		return 0, ErrKeyNotFound
	}
}

func (mr *MockGet) valDateTime(path string) (time.Time, error) {
	switch s := mr.getVal(path).(type) {
	case time.Time:
		return s, nil
	default:
		return time.Time{}, fmt.Errorf("Unable to Cast %#v to Time\n", s)
	}
}

// DateTime returns a time value
func (mr *MockGet) DateTime(opts ...ArgFunc) (time.Time, error) {
	path := mustNewArg(opts...).scopePath()
	switch {
	case mr.hasVal(path):
		return mr.valDateTime(path)
	case mr.FTime != nil:
		return mr.FTime(path)
	default:
		return time.Time{}, ErrKeyNotFound
	}
}

// Subscribe returns the before applied SubscriptionID and SubscriptionErr
// Does not start any underlying Goroutines.
func (mr *MockGet) Subscribe(path string, s MessageReceiver) (subscriptionID int, err error) {
	return mr.SubscriptionID, mr.SubscriptionErr
}

// NewScoped creates a new config.ScopedReader which uses the underlying
// mocked paths and values.
func (mr *MockGet) NewScoped(websiteID, groupID, storeID int64) ScopedGetter {
	return newScopedService(mr, websiteID, groupID, storeID)
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
