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
	"testing"

	"errors"

	"strings"

	"github.com/corestoreio/csfw/config/scope"
	"github.com/stretchr/testify/assert"
)

var _ scope.WebsiteIDer = (*mockScopeID)(nil)
var _ scope.GroupIDer = (*mockScopeID)(nil)
var _ scope.StoreIDer = (*mockScopeID)(nil)

type mockScopeID int64

func (id mockScopeID) StoreID() int64 {
	return int64(id)
}
func (id mockScopeID) GroupID() int64 {
	return int64(id)
}
func (id mockScopeID) WebsiteID() int64 {
	return int64(id)
}

func TestScopeKeyPath(t *testing.T) {
	tests := []struct {
		haveArg []ArgFunc
		want    string
		wantErr error
	}{
		{[]ArgFunc{Path("a/b/c")}, scope.StrDefault.FQPath("0", "a/b/c"), nil},
		{[]ArgFunc{Path("")}, "", ErrPathEmpty},
		{[]ArgFunc{Path()}, "", ErrPathEmpty},
		{[]ArgFunc{Scope(scope.DefaultID, -1)}, "", nil},
		{[]ArgFunc{Scope(scope.WebsiteID, -1)}, "", nil},
		{[]ArgFunc{Scope(scope.StoreID, -1)}, "", nil},
		{[]ArgFunc{Path("a/b/c"), Scope(scope.WebsiteID, -2)}, scope.StrDefault.FQPath("0", "a/b/c"), nil},
		{[]ArgFunc{Path("a/b/c"), Scope(scope.WebsiteID, 2)}, scope.StrWebsites.FQPath("2", "a/b/c"), nil},
		{[]ArgFunc{Path("a", "b", "c"), Scope(scope.WebsiteID, 200)}, scope.StrWebsites.FQPath("200", "a/b/c"), nil},
		{[]ArgFunc{Path("a", "b", "c"), Scope(scope.StoreID, 4)}, scope.StrStores.FQPath("4", "a/b/c"), nil},
		{[]ArgFunc{Path("a", "b"), Scope(scope.StoreID, 4)}, "", ErrPathEmpty},
		{[]ArgFunc{Path("a/b"), Scope(scope.StoreID, 4)}, "", errors.New("Incorrect number of paths elements: want 3, have 2, Path: [a/b]")},
		{[]ArgFunc{nil, Scope(scope.StoreID, 4)}, "", nil},
		{[]ArgFunc{Path("a", "b", "c"), ScopeStore(5)}, scope.StrStores.FQPath("5", "a/b/c"), nil},
		{[]ArgFunc{Path("a", "b", "c"), ScopeStore(0)}, scope.StrDefault.FQPath("0", "a/b/c"), nil},
		{[]ArgFunc{Path("a", "b", "c"), ScopeWebsite(50)}, scope.StrWebsites.FQPath("50", "a/b/c"), nil},
		{[]ArgFunc{Path("a", "b", "c"), ScopeWebsite(0)}, scope.StrDefault.FQPath("0", "a/b/c"), nil},
		{nil, "", nil},
	}

	for i, test := range tests {
		arg, err := newArg(test.haveArg...)
		if test.wantErr == nil {
			assert.NoError(t, err, "test IDX: %d", i)
		} else {
			assert.EqualError(t, err, test.wantErr.Error(), "test IDX: %d", i)
		}
		actualPath := arg.scopePath()
		assert.EqualValues(t, test.want, actualPath, "Test: %#v", test)
	}
}

func TestScopeKeyValue(t *testing.T) {
	tests := []struct {
		haveArg []ArgFunc
		want    string
		wantErr error
	}{
		{[]ArgFunc{Value(1), Path("a/b/c")}, scope.StrDefault.FQPath("0", "a/b/c"), nil},
		{[]ArgFunc{Value("1"), Path("")}, "", ErrPathEmpty},
		{[]ArgFunc{Value(1.1), Path()}, "", ErrPathEmpty},
		{[]ArgFunc{Value(1), Scope(scope.DefaultID, -9)}, "", nil},
		{[]ArgFunc{Value(1), Scope(scope.WebsiteID, 0)}, "", nil},
		{[]ArgFunc{Value(1), Scope(scope.StoreID, 0)}, "", nil},
		{[]ArgFunc{Value(1), Path("a/b/c"), Scope(scope.WebsiteID, 0)}, scope.StrDefault.FQPath("0", "a/b/c"), nil},
		{[]ArgFunc{Value(1), Path("a/b/c"), Scope(scope.WebsiteID, 2)}, scope.StrWebsites.FQPath("2", "a/b/c"), nil},
		{[]ArgFunc{Value(1), Path("a", "b", "c"), Scope(scope.WebsiteID, 200)}, scope.StrWebsites.FQPath("200", "a/b/c"), nil},
		{[]ArgFunc{Value(1), Path("a", "b", "c"), Scope(scope.StoreID, 4)}, scope.StrStores.FQPath("4", "a/b/c"), nil},
		{[]ArgFunc{Value(1), Path("a", "b"), Scope(scope.StoreID, 4)}, "", ErrPathEmpty},
		{[]ArgFunc{Value(1), Path("a/b"), Scope(scope.StoreID, 4)}, "", errors.New("Incorrect number of paths elements: want 3, have 2, Path: [a/b]")},
		{[]ArgFunc{Value(1), nil, Scope(scope.StoreID, 4)}, "", nil},
		{[]ArgFunc{Value(1), Path("a", "b", "c"), ScopeStore(5)}, scope.StrStores.FQPath("5", "a/b/c"), nil},
		{[]ArgFunc{Value(1.2), Path("a", "b", "c"), ScopeStore(-1)}, scope.StrDefault.FQPath("0", "a/b/c"), nil},
		{[]ArgFunc{Value(1.3), Path("a", "b", "c"), ScopeWebsite(50)}, scope.StrWebsites.FQPath("50", "a/b/c"), nil},
		{[]ArgFunc{ValueReader(strings.NewReader("a config value")), Path("a", "b", "c"), ScopeWebsite(0)}, scope.StrDefault.FQPath("0", "a/b/c"), nil},
		{nil, "", nil},
	}

	for i, test := range tests {
		a, err := newArg(test.haveArg...)
		if test.wantErr == nil {
			assert.NoError(t, err, "test IDX: %d", i)
		} else {
			assert.EqualError(t, err, test.wantErr.Error(), "test IDX: %d", i)
		}
		actualPath, actualVal := a.scopePath(), a.v
		assert.EqualValues(t, test.want, actualPath, "Test: %#v", test)
		if test.haveArg != nil && test.wantErr == nil {
			assert.NotEmpty(t, actualVal, "test IDX: %d", i)
		} else {
			assert.Empty(t, actualVal, "test IDX: %d", i)
		}
	}
}

func TestMustNewArg(t *testing.T) {
	defer func() { // protect ... you'll never know
		if r := recover(); r != nil {
			if err, ok := r.(error); ok {
				assert.EqualError(t, err, "Incorrect number of paths elements: want 3, have 2, Path: [a/b]")
			}
		}
	}()
	a := mustNewArg(Path("a/b"))
	assert.NotNil(t, a)
}

type testMultiErrorReader struct{}

func (testMultiErrorReader) Read(p []byte) (n int, err error) {
	return 0, errors.New("testMultiErrorReader error")
}

func TestMultiError(t *testing.T) {

	a, err := newArg(Path("a/b"), ValueReader(testMultiErrorReader{}))
	assert.NotNil(t, a)
	assert.EqualError(t, err, "Incorrect number of paths elements: want 3, have 2, Path: [a/b]\nValueReader error testMultiErrorReader error")
}

var benchmarkScopeKey string

// BenchmarkScopeKey____InMap-4      	 2000000	       622 ns/op	     336 B/op	       6 allocs/op
func BenchmarkScopeKey____InMap(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		arg, _ := newArg(Path("a", "b", "c"), Scope(scope.WebsiteID, 4))
		benchmarkScopeKey = arg.scopePath()
	}
}

// BenchmarkScopeKey_NotInMap-4      	 2000000	       687 ns/op	     368 B/op	       7 allocs/op
func BenchmarkScopeKey_NotInMap(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		arg, _ := newArg(Path("a", "b", "c"), Scope(scope.WebsiteID, 40))
		benchmarkScopeKey = arg.scopePath()
	}
}

// BenchmarkScopeKey____InMapNoJoin-4	 2000000	       768 ns/op	     352 B/op	       7 allocs/op
func BenchmarkScopeKey____InMapNoJoin(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		arg, _ := newArg(Path("a/b/c"), Scope(scope.WebsiteID, 3))
		benchmarkScopeKey = arg.scopePath()
	}
}
