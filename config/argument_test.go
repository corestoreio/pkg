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

func TestScopeKeyPath(t *testing.T) {
	tests := []struct {
		haveArg []ArgFunc
		want    string
		wantErr error
	}{
		{[]ArgFunc{Path("a/b/c")}, scope.RangeDefault + "/0/a/b/c", nil},
		{[]ArgFunc{Path("")}, "", ErrPathEmpty},
		{[]ArgFunc{Path()}, "", ErrPathEmpty},
		{[]ArgFunc{Scope(scope.DefaultID, nil)}, "", nil},
		{[]ArgFunc{Scope(scope.WebsiteID, nil)}, "", nil},
		{[]ArgFunc{Scope(scope.StoreID, nil)}, "", nil},
		{[]ArgFunc{Path("a/b/c"), Scope(scope.WebsiteID, nil)}, scope.RangeDefault + "/0/a/b/c", nil},
		{[]ArgFunc{Path("a/b/c"), Scope(scope.WebsiteID, scope.ID(2))}, scope.RangeWebsites + "/2/a/b/c", nil},
		{[]ArgFunc{Path("a", "b", "c"), Scope(scope.WebsiteID, scope.ID(200))}, scope.RangeWebsites + "/200/a/b/c", nil},
		{[]ArgFunc{Path("a", "b", "c"), Scope(scope.StoreID, scope.ID(4))}, scope.RangeStores + "/4/a/b/c", nil},
		{[]ArgFunc{Path("a", "b"), Scope(scope.StoreID, scope.ID(4))}, "", errors.New("Incorrect number of paths elements: want 3, have 1, Path: [a b]")},
		{[]ArgFunc{nil, Scope(scope.StoreID, scope.ID(4))}, "", nil},
		{[]ArgFunc{Path("a", "b", "c"), ScopeStore(scope.ID(5))}, scope.RangeStores + "/5/a/b/c", nil},
		{[]ArgFunc{Path("a", "b", "c"), ScopeStore(nil)}, scope.RangeDefault + "/0/a/b/c", nil},
		{[]ArgFunc{Path("a", "b", "c"), ScopeWebsite(scope.ID(50))}, scope.RangeWebsites + "/50/a/b/c", nil},
		{[]ArgFunc{Path("a", "b", "c"), ScopeWebsite(nil)}, scope.RangeDefault + "/0/a/b/c", nil},
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
		{[]ArgFunc{Value(1), Path("a/b/c")}, scope.RangeDefault + "/0/a/b/c", nil},
		{[]ArgFunc{Value("1"), Path("")}, "", ErrPathEmpty},
		{[]ArgFunc{Value(1.1), Path()}, "", ErrPathEmpty},
		{[]ArgFunc{Value(1), Scope(scope.DefaultID, nil)}, "", nil},
		{[]ArgFunc{Value(1), Scope(scope.WebsiteID, nil)}, "", nil},
		{[]ArgFunc{Value(1), Scope(scope.StoreID, nil)}, "", nil},
		{[]ArgFunc{Value(1), Path("a/b/c"), Scope(scope.WebsiteID, nil)}, scope.RangeDefault + "/0/a/b/c", nil},
		{[]ArgFunc{Value(1), Path("a/b/c"), Scope(scope.WebsiteID, scope.ID(2))}, scope.RangeWebsites + "/2/a/b/c", nil},
		{[]ArgFunc{Value(1), Path("a", "b", "c"), Scope(scope.WebsiteID, scope.ID(200))}, scope.RangeWebsites + "/200/a/b/c", nil},
		{[]ArgFunc{Value(1), Path("a", "b", "c"), Scope(scope.StoreID, scope.ID(4))}, scope.RangeStores + "/4/a/b/c", nil},
		{[]ArgFunc{Value(1), Path("a", "b"), Scope(scope.StoreID, scope.ID(4))}, "", errors.New("Incorrect number of paths elements: want 3, have 1, Path: [a b]")},
		{[]ArgFunc{Value(1), nil, Scope(scope.StoreID, scope.ID(4))}, "", nil},
		{[]ArgFunc{Value(1), Path("a", "b", "c"), ScopeStore(scope.ID(5))}, scope.RangeStores + "/5/a/b/c", nil},
		{[]ArgFunc{Value(1.2), Path("a", "b", "c"), ScopeStore(nil)}, scope.RangeDefault + "/0/a/b/c", nil},
		{[]ArgFunc{Value(1.3), Path("a", "b", "c"), ScopeWebsite(scope.ID(50))}, scope.RangeWebsites + "/50/a/b/c", nil},
		{[]ArgFunc{ValueReader(strings.NewReader("a config value")), Path("a", "b", "c"), ScopeWebsite(nil)}, scope.RangeDefault + "/0/a/b/c", nil},
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

// BenchmarkScopeKey____InMap-4      	 2000000	       663 ns/op	     288 B/op	       7 allocs/op
func BenchmarkScopeKey____InMap(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		arg, _ := newArg(Path("a", "b", "c"), Scope(scope.WebsiteID, scope.ID(4)))
		benchmarkScopeKey = arg.scopePath()
	}
}

// BenchmarkScopeKey_NotInMap-4      	 2000000	       769 ns/op	     320 B/op	       8 allocs/op
func BenchmarkScopeKey_NotInMap(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		arg, _ := newArg(Path("a", "b", "c"), Scope(scope.WebsiteID, scope.ID(40)))
		benchmarkScopeKey = arg.scopePath()
	}
}

// BenchmarkScopeKey____InMapNoJoin-4	 2000000	       687 ns/op	     288 B/op	       7 allocs/op
func BenchmarkScopeKey____InMapNoJoin(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		arg, _ := newArg(Path("a/b/c"), Scope(scope.WebsiteID, scope.ID(3)))
		benchmarkScopeKey = arg.scopePath()
	}
}
