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
	"testing"

	"strings"

	"github.com/corestoreio/csfw/config/path"
	"github.com/corestoreio/csfw/store/scope"
	"github.com/corestoreio/csfw/util/log"
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
	t.Parallel()
	mainRoute := path.NewRoute("aa/bb/cc")
	mainPath := path.MustNewByParts("aa", "bb", "cc")
	tests := []struct {
		haveArg []ArgFunc
		want    string
		wantErr error
	}{
		{[]ArgFunc{Route(mainRoute)}, path.MustNewByParts("aa/bb/cc").String(), nil},
		{[]ArgFunc{Route(path.NewRoute(""))}, "", path.ErrRouteEmpty},
		{[]ArgFunc{Route(path.Route{})}, "", path.ErrRouteEmpty},
		{[]ArgFunc{Scope(scope.DefaultID, -1)}, "", nil},
		{[]ArgFunc{Scope(scope.WebsiteID, -1)}, "", nil},
		{[]ArgFunc{Scope(scope.StoreID, -1)}, "", nil},
		{[]ArgFunc{Route(mainRoute), Scope(scope.WebsiteID, -2)}, path.MustNewByParts("aa/bb/cc").String(), nil},
		{[]ArgFunc{Route(mainRoute), Scope(scope.WebsiteID, 2)}, path.MustNewByParts("aa/bb/cc").Bind(scope.WebsiteID, 2).String(), nil},
		{[]ArgFunc{Path(mainPath), Scope(scope.WebsiteID, 200)}, path.MustNewByParts("aa/bb/cc").Bind(scope.WebsiteID, 200).String(), nil},
		{[]ArgFunc{Path(mainPath), Scope(scope.StoreID, 4)}, path.MustNewByParts("aa/bb/cc").Bind(scope.StoreID, 4).String(), nil},
		{[]ArgFunc{Path(path.Path{Route: path.NewRoute("a", "b")}), Scope(scope.StoreID, 4)}, "", path.ErrIncorrectPath},
		{[]ArgFunc{Path(path.Path{Route: path.NewRoute("a/b")}), Scope(scope.StoreID, 4)}, "", path.ErrIncorrectPath},
		{[]ArgFunc{nil, Scope(scope.StoreID, 4)}, "", nil},
		{[]ArgFunc{Path(mainPath), ScopeStore(5)}, path.MustNewByParts("aa/bb/cc").Bind(scope.StoreID, 5).String(), nil},
		{[]ArgFunc{Path(mainPath), ScopeStore(0)}, path.MustNewByParts("aa/bb/cc").String(), nil},
		{[]ArgFunc{Path(mainPath), ScopeWebsite(50)}, path.MustNewByParts("aa/bb/cc").Bind(scope.WebsiteID, 50).String(), nil},
		{[]ArgFunc{Path(mainPath), ScopeWebsite(0)}, path.MustNewByParts("aa/bb/cc").String(), nil},
		{[]ArgFunc{PathScoped("rr/ss/tt", scope.WebsiteID, 2)}, path.MustNewByParts("rr/ss/tt").Bind(scope.WebsiteID, 2).String(), nil},
		{[]ArgFunc{PathScoped("rr/ss", scope.WebsiteID, 2)}, "", path.ErrIncorrectPath},
		{[]ArgFunc{PathScoped("", scope.WebsiteID, 2)}, "", path.ErrRouteEmpty},
		{nil, "", nil},
	}

	for i, test := range tests {
		a, err := newArg(test.haveArg...)
		if test.wantErr == nil {
			assert.NoError(t, err, "Index: %d", i)
		} else {
			assert.EqualError(t, err, test.wantErr.Error(), "Index: %d", i)
		}
		actualPath := a.String()
		assert.EqualValues(t, test.want, actualPath, "Index: %d", i)
	}
}

func TestScopeKeyValue(t *testing.T) {
	t.Parallel()
	defaultPath := path.MustNew(path.NewRoute("aa/bb/cc"))
	tests := []struct {
		haveArg []ArgFunc
		want    string
		wantErr error
	}{
		{[]ArgFunc{Value(1), PathScoped("aa/bb/cc", 0, 0)}, defaultPath.String(), nil},
		{[]ArgFunc{Value("123"), PathScoped("", 0, 0)}, "", path.ErrRouteEmpty},
		{[]ArgFunc{Value(1.321), Path(path.Path{})}, "", path.ErrRouteEmpty},
		{[]ArgFunc{Value(1), Scope(scope.DefaultID, -9)}, "", nil},
		{[]ArgFunc{Value(1), Scope(scope.WebsiteID, 0)}, "", nil},
		{[]ArgFunc{Value(1), Scope(scope.StoreID, 0)}, "", nil},
		{[]ArgFunc{Value(1), PathScoped("aa/bb/cc", scope.WebsiteID, 0)}, `websites/0/aa/bb/cc`, nil},
		{[]ArgFunc{Value(1), Path(path.MustNewByParts("aa/bb/cc")), Scope(scope.WebsiteID, 2)}, defaultPath.Bind(scope.WebsiteID, 2).String(), nil},
		{[]ArgFunc{Value(8), Path(path.MustNewByParts("aa", "bb", "cc")), Scope(scope.WebsiteID, 200)}, defaultPath.Bind(scope.WebsiteID, 200).String(), nil},
		{[]ArgFunc{Value(9), Path(path.MustNewByParts("aa", "bb", "cc")), Scope(scope.StoreID, 4)}, defaultPath.Bind(scope.StoreID, 4).String(), nil},
		{[]ArgFunc{Value(10), PathScoped("a/b", scope.StoreID, 4)}, "", path.ErrIncorrectPath},
		{[]ArgFunc{Value(12), nil, Scope(scope.StoreID, 4)}, "", nil},
		{[]ArgFunc{Value(1), Path(path.MustNewByParts("aa", "bb", "cc")), ScopeStore(5)}, defaultPath.Bind(scope.StoreID, 5).String(), nil},
		{[]ArgFunc{Value(1.2), Path(path.MustNewByParts("aa", "bb", "cc")), ScopeStore(-1)}, defaultPath.String(), nil},
		{[]ArgFunc{Value(1.3), Path(path.MustNewByParts("aa", "bb", "cc")), ScopeWebsite(50)}, defaultPath.Bind(scope.WebsiteID, 50).String(), nil},
		{[]ArgFunc{ValueReader(strings.NewReader("a config value")), PathScoped("aa/bb/cc", scope.WebsiteID, 0)}, `websites/0/aa/bb/cc`, nil},
		{nil, "", nil},
	}

	for i, test := range tests {
		a, err := newArg(test.haveArg...)
		if test.wantErr == nil {
			assert.NoError(t, err, "Index: %d", i)
		} else {
			assert.EqualError(t, err, test.wantErr.Error(), "Index: %d", i)
		}
		actualPath, actualVal := a.String(), a.v
		assert.EqualValues(t, test.want, actualPath, "Index: %d", i)
		if test.haveArg != nil && test.wantErr == nil {
			assert.NotEmpty(t, actualVal, "Index: %d", i)
		} else {
			assert.Empty(t, actualVal, "Index: %d", i)
		}
	}
}

func TestMustNewArg(t *testing.T) {
	t.Parallel()
	defer func() { // protect ... you'll never know
		if r := recover(); r != nil {
			if err, ok := r.(error); ok {
				assert.EqualError(t, err, path.ErrIncorrectPath.Error())
			}
		}
	}()
	a := mustNewArg(Path(path.MustNewByParts("aa/bb")))
	assert.NotNil(t, a)
}

type testMultiErrorReader struct{}

func (testMultiErrorReader) Read(p []byte) (n int, err error) {
	return 0, errors.New("testMultiErrorReader error")
}

func TestMultiError(t *testing.T) {
	t.Parallel()
	a, err := newArg(Path(path.Path{Route: path.NewRoute("a/b")}), ValueReader(testMultiErrorReader{}))
	assert.NotNil(t, a)
	assert.EqualError(t, err, "Incorrect Path. Either to short or missing path separator.\nValueReader error testMultiErrorReader error")
}

var benchmarkScopeKey string

const benchmarkScopeKeyWant = "websites/40/aa/bb/cc"

// BenchmarkScopeKey____InMap-4      	 2000000	       622 ns/op	     336 B/op	       6 allocs/op
// BenchmarkScopeKey____InMap-4       	 1000000	      1139 ns/op	     400 B/op	       9 allocs/op
// BenchmarkScopeKey____InMap___NewPath-4	 1000000	      1341 ns/op	     528 B/op	      10 allocs/op
func BenchmarkScopeKey____InMap___NewPath(b *testing.B) {
	PkgLog.SetLevel(log.StdLevelFatal)
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		a, _ := newArg(Path(path.MustNewByParts("aa", "bb", "cc")), Scope(scope.WebsiteID, 40))
		benchmarkScopeKey = a.String()
	}
	if benchmarkScopeKey != benchmarkScopeKeyWant {
		b.Errorf("Want %s Have %s", benchmarkScopeKeyWant, benchmarkScopeKey)
	}
}

// BenchmarkScopeKey____InMap_ExistPath-4	 1000000	      1023 ns/op	     480 B/op	       8 allocs/op
func BenchmarkScopeKey____InMap_ExistPath(b *testing.B) {
	p := path.MustNewByParts("aa", "bb", "cc").Bind(scope.WebsiteID, 40)
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		a, _ := newArg(Path(p))
		benchmarkScopeKey = a.String()
	}
	if benchmarkScopeKey != benchmarkScopeKeyWant {
		b.Errorf("Want %s Have %s", benchmarkScopeKeyWant, benchmarkScopeKey)
	}
}

// BenchmarkScopeKey_NotInMap-4      	 2000000	       687 ns/op	     368 B/op	       7 allocs/op
// BenchmarkScopeKey_NotInMap-4       	 1000000	      1241 ns/op	     416 B/op	      10 allocs/op
// BenchmarkScopeKey_NotInMap___NewPath-4	 1000000	      1396 ns/op	     528 B/op	      10 allocs/op
func BenchmarkScopeKey_NotInMap___NewPath(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		a, _ := newArg(Path(path.MustNewByParts("aa", "bb", "cc")), Scope(scope.WebsiteID, 40))
		benchmarkScopeKey = a.String()
	}
	if benchmarkScopeKey != benchmarkScopeKeyWant {
		b.Errorf("Want %s Have %s", benchmarkScopeKeyWant, benchmarkScopeKey)
	}
}

// BenchmarkScopeKey_NotInMap_ExistPath-4	 1000000	      1276 ns/op	     480 B/op	       8 allocs/op
func BenchmarkScopeKey_NotInMap_ExistPath(b *testing.B) {
	p := path.MustNewByParts("aa", "bb", "cc").Bind(scope.WebsiteID, 40)
	b.ResetTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		a, _ := newArg(Path(p))
		benchmarkScopeKey = a.String()
	}
	if benchmarkScopeKey != benchmarkScopeKeyWant {
		b.Errorf("Want %s Have %s", benchmarkScopeKeyWant, benchmarkScopeKey)
	}
}
