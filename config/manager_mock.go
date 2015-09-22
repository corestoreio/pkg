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

// mockOptionFunc to initialize the NewMockReader
type mockOptionFunc func(*MockReader)

// MockReader used for testing. Contains functions which will be called in the
// appropriate methods of interface config.Reader.
type MockReader struct {
	s   func(path string) string
	b   func(path string) bool
	f64 func(path string) float64
	i   func(path string) int
	t   func(path string) time.Time
}

// MockPathScopeDefault creates for testing a fully qualified path for the
// default scope from a Scope ID and a path string (a/b/c)
func MockPathScopeDefault(id int64, path string) string {
	return scope.RangeDefault + PS + strconv.FormatInt(id, 10) + PS + path
}

// MockPathScopeWebsite creates for testing a fully qualified path for the
// website scope from a Scope ID and a path string (a/b/c)
func MockPathScopeWebsite(id int64, path string) string {
	return scope.RangeWebsites + PS + strconv.FormatInt(id, 10) + PS + path
}

// MockPathScopeStore creates for testing a fully qualified path for the
// store scope from a Scope ID and a path string (a/b/c)
func MockPathScopeStore(id int64, path string) string {
	return scope.RangeStores + PS + strconv.FormatInt(id, 10) + PS + path
}

// MockString returns a function which can be used in the NewMockReader().
// Your function returns a string value from a given path.
func MockString(f func(path string) string) mockOptionFunc {
	return func(mr *MockReader) {
		mr.s = f
	}
}

// MockBool returns a function which can be used in the NewMockReader().
// Your function returns a bool value from a given path.
func MockBool(f func(path string) bool) mockOptionFunc {
	return func(mr *MockReader) {
		mr.b = f
	}
}

// MockFloat64 returns a function which can be used in the NewMockReader().
// Your function returns a float64 value from a given path.
func MockFloat64(f func(path string) float64) mockOptionFunc {
	return func(mr *MockReader) {
		mr.f64 = f
	}
}

// MockInt returns a function which can be used in the NewMockReader().
// Your function returns an int value from a given path.
func MockInt(f func(path string) int) mockOptionFunc {
	return func(mr *MockReader) {
		mr.i = f
	}
}

// MockTime returns a function which can be used in the NewMockReader().
// Your function returns a Time value from a given path.
func MockTime(f func(path string) time.Time) mockOptionFunc {
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

func (sr *MockReader) GetString(opts ...ArgFunc) string {
	if sr.s == nil {
		return ""
	}
	return sr.s(mustNewArg(opts...).scopePath())
}

func (sr *MockReader) GetBool(opts ...ArgFunc) bool {
	if sr.b == nil {
		return false
	}
	return sr.b(mustNewArg(opts...).scopePath())
}
func (sr *MockReader) GetFloat64(opts ...ArgFunc) float64 {
	if sr.f64 == nil {
		return 0.0
	}
	return sr.f64(mustNewArg(opts...).scopePath())
}
func (sr *MockReader) GetInt(opts ...ArgFunc) int {
	if sr.i == nil {
		return 0
	}
	return sr.i(mustNewArg(opts...).scopePath())
}
func (sr *MockReader) GetDateTime(opts ...ArgFunc) time.Time {
	if sr.t == nil {
		return time.Time{}
	}
	return sr.t(mustNewArg(opts...).scopePath())
}
