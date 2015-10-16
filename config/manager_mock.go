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
	"strconv"
	"time"

	"github.com/corestoreio/csfw/config/scope"
)

var _ Reader = (*MockReader)(nil)
var _ ReaderPubSuber = (*MockReader)(nil)

// mockOptionFunc to initialize the NewMockReader
type mockOptionFunc func(*MockReader)

// MockReader used for testing. Contains functions which will be called in the
// appropriate methods of interface config.Reader.
type MockReader struct {
	s               func(path string) (string, error)
	b               func(path string) (bool, error)
	f64             func(path string) (float64, error)
	i               func(path string) (int, error)
	t               func(path string) (time.Time, error)
	SubscriptionID  int
	SubscriptionErr error
}

// MockPathScopeDefault creates for testing a fully qualified path for the
// default scope from a Scope ID and a path string (a/b/c)
func MockPathScopeDefault(id int64, path string) string {
	return scope.StrDefault.FQPath(strconv.FormatInt(id, 10), path)
}

// MockPathScopeWebsite creates for testing a fully qualified path for the
// website scope from a Scope ID and a path string (a/b/c)
func MockPathScopeWebsite(id int64, path string) string {
	return scope.StrWebsites.FQPath(strconv.FormatInt(id, 10), path)
}

// MockPathScopeStore creates for testing a fully qualified path for the
// store scope from a Scope ID and a path string (a/b/c)
func MockPathScopeStore(id int64, path string) string {
	return scope.StrStores.FQPath(strconv.FormatInt(id, 10), path)
}

// MockString returns a function which can be used in the NewMockReader().
// Your function returns a string value from a given path.
func MockString(f func(path string) (string, error)) mockOptionFunc {
	return func(mr *MockReader) { mr.s = f }
}

// MockBool returns a function which can be used in the NewMockReader().
// Your function returns a bool value from a given path.
func MockBool(f func(path string) (bool, error)) mockOptionFunc {
	return func(mr *MockReader) { mr.b = f }
}

// MockFloat64 returns a function which can be used in the NewMockReader().
// Your function returns a float64 value from a given path.
func MockFloat64(f func(path string) (float64, error)) mockOptionFunc {
	return func(mr *MockReader) { mr.f64 = f }
}

// MockInt returns a function which can be used in the NewMockReader().
// Your function returns an int value from a given path.
func MockInt(f func(path string) (int, error)) mockOptionFunc {
	return func(mr *MockReader) { mr.i = f }
}

// MockTime returns a function which can be used in the NewMockReader().
// Your function returns a Time value from a given path.
func MockTime(f func(path string) (time.Time, error)) mockOptionFunc {
	return func(mr *MockReader) {
		mr.t = f
	}
}

// NewMockReader used for testing
func NewMockReader(opts ...mockOptionFunc) *MockReader {
	mr := &MockReader{}
	for _, opt := range opts {
		opt(mr)
	}
	return mr
}

func (sr *MockReader) GetString(opts ...ArgFunc) (string, error) {
	if sr.s == nil {
		return "", ErrKeyNotFound
	}
	return sr.s(mustNewArg(opts...).scopePath())
}

func (sr *MockReader) GetBool(opts ...ArgFunc) (bool, error) {
	if sr.b == nil {
		return false, ErrKeyNotFound
	}
	return sr.b(mustNewArg(opts...).scopePath())
}
func (sr *MockReader) GetFloat64(opts ...ArgFunc) (float64, error) {
	if sr.f64 == nil {
		return 0.0, ErrKeyNotFound
	}
	return sr.f64(mustNewArg(opts...).scopePath())
}
func (sr *MockReader) GetInt(opts ...ArgFunc) (int, error) {
	if sr.i == nil {
		return 0, ErrKeyNotFound
	}
	return sr.i(mustNewArg(opts...).scopePath())
}
func (sr *MockReader) GetDateTime(opts ...ArgFunc) (time.Time, error) {
	if sr.t == nil {
		return time.Time{}, ErrKeyNotFound
	}
	return sr.t(mustNewArg(opts...).scopePath())
}

func (sr *MockReader) Subscribe(path string, s MessageReceiver) (subscriptionID int, err error) {
	return sr.SubscriptionID, sr.SubscriptionErr
}

func (sr *MockReader) NewScoped(websiteID, groupID, storeID int64) ScopedReader {
	return newScopedManager(sr, websiteID, groupID, storeID)
}
