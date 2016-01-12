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

package config

import (
	"errors"
	"strings"
	"testing"

	"github.com/corestoreio/csfw/config/path"
	"github.com/corestoreio/csfw/store/scope"
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
		{[]ArgFunc{Path("aa/bb/cc")}, path.MustNew("aa/bb/cc").String(), nil},
		{[]ArgFunc{Path("")}, "", errors.New("Incorrect number of paths elements: want 3, have 1, Path: []")},
		{[]ArgFunc{Path()}, "", errors.New("Incorrect number of paths elements: want 3, have 0, Path: []")},
		{[]ArgFunc{Scope(scope.DefaultID, -1)}, "", nil},
		{[]ArgFunc{Scope(scope.WebsiteID, -1)}, "", nil},
		{[]ArgFunc{Scope(scope.StoreID, -1)}, "", nil},
		{[]ArgFunc{Path("aa/bb/cc"), Scope(scope.WebsiteID, -2)}, path.MustNew("aa/bb/cc").String(), nil},
		{[]ArgFunc{Path("aa/bb/cc"), Scope(scope.WebsiteID, 2)}, path.MustNew("aa/bb/cc").Bind(scope.WebsiteID, 2).String(), nil},
		{[]ArgFunc{Path("aa", "bb", "cc"), Scope(scope.WebsiteID, 200)}, path.MustNew("aa/bb/cc").Bind(scope.WebsiteID, 200).String(), nil},
		{[]ArgFunc{Path("aa", "bb", "cc"), Scope(scope.StoreID, 4)}, path.MustNew("aa/bb/cc").Bind(scope.StoreID, 4).String(), nil},
		{[]ArgFunc{Path("a", "b"), Scope(scope.StoreID, 4)}, "", errors.New("Incorrect number of paths elements: want 3, have 2, Path: [a b]")},
		{[]ArgFunc{Path("a/b"), Scope(scope.StoreID, 4)}, "", errors.New("This path part \"a\" is too short. Parts: []string{\"a\", \"b\"}")},
		{[]ArgFunc{nil, Scope(scope.StoreID, 4)}, "", nil},
		{[]ArgFunc{Path("aa", "bb", "cc"), ScopeStore(5)}, path.MustNew("aa/bb/cc").Bind(scope.StoreID, 5).String(), nil},
		{[]ArgFunc{Path("aa", "bb", "cc"), ScopeStore(0)}, path.MustNew("aa/bb/cc").String(), nil},
		{[]ArgFunc{Path("aa", "bb", "cc"), ScopeWebsite(50)}, path.MustNew("aa/bb/cc").Bind(scope.WebsiteID, 50).String(), nil},
		{[]ArgFunc{Path("aa", "bb", "cc"), ScopeWebsite(0)}, path.MustNew("aa/bb/cc").String(), nil},
		{nil, "", nil},
	}

	for i, test := range tests {
		a, err := newArg(test.haveArg...)
		if test.wantErr == nil {
			assert.NoError(t, err, "test IDX: %d", i)
		} else {
			assert.EqualError(t, err, test.wantErr.Error(), "test IDX: %d", i)
		}
		actualPath := a.String()
		assert.EqualValues(t, test.want, actualPath, "Test: %#v", test)
	}
}

func TestScopeKeyValue(t *testing.T) {
	tests := []struct {
		haveArg []ArgFunc
		want    string
		wantErr error
	}{
		{[]ArgFunc{Value(1), Path("aa/bb/cc")}, path.MustNew("aa/bb/cc").String(), nil},
		{[]ArgFunc{Value("123"), Path("")}, "", errors.New("Incorrect number of paths elements: want 3, have 1, Path: []")},
		{[]ArgFunc{Value(1.321), Path()}, "", errors.New("Incorrect number of paths elements: want 3, have 0, Path: []")},
		{[]ArgFunc{Value(1), Scope(scope.DefaultID, -9)}, "", nil},
		{[]ArgFunc{Value(1), Scope(scope.WebsiteID, 0)}, "", nil},
		{[]ArgFunc{Value(1), Scope(scope.StoreID, 0)}, "", nil},
		{[]ArgFunc{Value(1), Path("aa/bb/cc"), Scope(scope.WebsiteID, 0)}, path.MustNew("aa/bb/cc").String(), nil},
		{[]ArgFunc{Value(1), Path("aa/bb/cc"), Scope(scope.WebsiteID, 2)}, path.MustNew("aa/bb/cc").Bind(scope.WebsiteID, 2).String(), nil},
		{[]ArgFunc{Value(8), Path("aa", "bb", "cc"), Scope(scope.WebsiteID, 200)}, path.MustNew("aa/bb/cc").Bind(scope.WebsiteID, 200).String(), nil},
		{[]ArgFunc{Value(9), Path("aa", "bb", "cc"), Scope(scope.StoreID, 4)}, path.MustNew("aa/bb/cc").Bind(scope.StoreID, 4).String(), nil},
		{[]ArgFunc{Value(10), Path("a", "b"), Scope(scope.StoreID, 4)}, "", errors.New("Incorrect number of paths elements: want 3, have 2, Path: [a b]")},
		{[]ArgFunc{Value(11), Path("a/b"), Scope(scope.StoreID, 4)}, "", errors.New("This path part \"a\" is too short. Parts: []string{\"a\", \"b\"}")},
		{[]ArgFunc{Value(12), nil, Scope(scope.StoreID, 4)}, "", nil},
		{[]ArgFunc{Value(1), Path("aa", "bb", "cc"), ScopeStore(5)}, path.MustNew("aa/bb/cc").Bind(scope.StoreID, 5).String(), nil},
		{[]ArgFunc{Value(1.2), Path("aa", "bb", "cc"), ScopeStore(-1)}, path.MustNew("aa/bb/cc").String(), nil},
		{[]ArgFunc{Value(1.3), Path("aa", "bb", "cc"), ScopeWebsite(50)}, path.MustNew("aa/bb/cc").Bind(scope.WebsiteID, 50).String(), nil},
		{[]ArgFunc{ValueReader(strings.NewReader("a config value")), Path("aa", "bb", "cc"), ScopeWebsite(0)}, path.MustNew("aa/bb/cc").String(), nil},
		{nil, "", nil},
	}

	for i, test := range tests {
		a, err := newArg(test.haveArg...)
		if test.wantErr == nil {
			assert.NoError(t, err, "test IDX: %d", i)
		} else {
			assert.EqualError(t, err, test.wantErr.Error(), "test IDX: %d", i)
		}
		actualPath, actualVal := a.String(), a.v
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
				assert.EqualError(t, err, "All arguments must be valid! Min want: 3. Have: 2. Parts []string{\"aa\", \"bb\"}")
			}
		}
	}()
	a := mustNewArg(Path("aa/bb"))
	assert.NotNil(t, a)
}

type testMultiErrorReader struct{}

func (testMultiErrorReader) Read(p []byte) (n int, err error) {
	return 0, errors.New("testMultiErrorReader error")
}

func TestMultiError(t *testing.T) {
	a, err := newArg(Path("a/b"), ValueReader(testMultiErrorReader{}))
	assert.NotNil(t, a)
	assert.EqualError(t, err, "This path part \"a\" is too short. Parts: []string{\"a\", \"b\"}\nValueReader error testMultiErrorReader error")
}

var benchmarkScopeKey string

// BenchmarkScopeKey____InMap-4      	 2000000	       622 ns/op	     336 B/op	       6 allocs/op
// BenchmarkScopeKey____InMap-4       	 1000000	      1139 ns/op	     400 B/op	       9 allocs/op
func BenchmarkScopeKey____InMap(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		a, _ := newArg(Path("aa", "bb", "cc"), Scope(scope.WebsiteID, 4))
		benchmarkScopeKey = a.String()
	}
}

// BenchmarkScopeKey_NotInMap-4      	 2000000	       687 ns/op	     368 B/op	       7 allocs/op
// BenchmarkScopeKey_NotInMap-4       	 1000000	      1241 ns/op	     416 B/op	      10 allocs/op
func BenchmarkScopeKey_NotInMap(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		a, _ := newArg(Path("aa", "bb", "cc"), Scope(scope.WebsiteID, 40))
		benchmarkScopeKey = a.String()
	}
}

// BenchmarkScopeKey____InMapNoJoin-4	 2000000	       768 ns/op	     352 B/op	       7 allocs/op
// BenchmarkScopeKey____InMapNoJoin-4 	 1000000	      1337 ns/op	     416 B/op	      10 allocs/op
func BenchmarkScopeKey____InMapNoJoin(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		a, _ := newArg(Path("aa/bb/cc"), Scope(scope.WebsiteID, 3))
		benchmarkScopeKey = a.String()
	}
}
