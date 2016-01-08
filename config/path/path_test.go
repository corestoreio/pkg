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

package path_test

import (
	"errors"
	"strconv"
	"testing"

	"github.com/corestoreio/csfw/config/path"
	"github.com/corestoreio/csfw/store/scope"
	"github.com/stretchr/testify/assert"
)

func TestFQ(t *testing.T) {
	t.Parallel()
	tests := []struct {
		str     scope.StrScope
		id      string
		path    []string
		want    string
		wantErr error
	}{
		{scope.StrDefault, "0", []string{"system/dev/debug"}, scope.StrDefault.String() + "/0/system/dev/debug", nil},
		{scope.StrDefault, "33", []string{"system", "dev", "debug"}, scope.StrDefault.String() + "/0/system/dev/debug", nil},
		{scope.StrWebsites, "0", []string{"system/dev/debug"}, scope.StrWebsites.String() + "/0/system/dev/debug", nil},
		{scope.StrWebsites, "343", []string{"system", "dev", "debug"}, scope.StrWebsites.String() + "/343/system/dev/debug", nil},
		{scope.StrScope("hello"), "343", []string{"system", "dev", "debug"}, scope.StrWebsites.String() + "/343/system/dev/debug", scope.ErrUnsupportedScope},
	}
	for i, test := range tests {
		have, haveErr := path.FQ(test.str, test.id, test.path...)
		if test.wantErr != nil {
			assert.Empty(t, have, "Index %d", i)
			assert.EqualError(t, haveErr, test.wantErr.Error(), "Index %d", i)
			continue
		}
		assert.NoError(t, haveErr, "Index %d", i)
		assert.Equal(t, test.want, have, "Index %d", i)
	}
	assert.Equal(t, "stores/7475/catalog/frontend/list_allow_all", path.MustFQInt64(scope.StrStores, 7475, "catalog", "frontend", "list_allow_all"))
	assert.Equal(t, "stores/5/catalog/frontend/list_allow_all", path.MustFQInt64(scope.StrStores, 5, "catalog", "frontend", "list_allow_all"))
}

func TestMustFQInt64_01(t *testing.T) {
	t.Parallel()
	assert.Exactly(t, "stores/345/x/y/z", path.MustFQInt64(scope.StrStores, 345, "x", "y", "z"))
	defer func() {
		if r := recover(); r != nil {
			assert.EqualError(t, r.(error), scope.ErrUnsupportedScope.Error())
		}
	}()
	_ = path.MustFQInt64(scope.StrScope("invalid"), 345, "x", "y", "z")
}

func TestMustFQ_01(t *testing.T) {
	t.Parallel()
	assert.Exactly(t, "stores/345/x/y/z", path.MustFQ(scope.StrStores, "345", "x", "y", "z"))
	defer func() {
		if r := recover(); r != nil {
			assert.EqualError(t, r.(error), scope.ErrUnsupportedScope.Error())
		}
	}()
	_ = path.MustFQ(scope.StrScope("invalid"), "345", "x", "y", "z")
}

func TestMustFQ_02(t *testing.T) {
	t.Parallel()
	assert.Exactly(t, "default/0/x/y/z", path.MustFQ(scope.StrDefault, "345", "x", "y", "z"))
	// the next assertion may can get adjusted because invalid path parts
	assert.Exactly(t, "websites/345/x/y", path.MustFQ(scope.StrWebsites, "345", "x", "y"))
}

var benchmarkStrScopeFQPath string

// BenchmarkStrScopeFQPath-4	 5000000	       384 ns/op	      32 B/op	       1 allocs/op
func BenchmarkStrScopeFQPath(b *testing.B) {
	want := scope.StrWebsites.String() + "/4/system/dev/debug"
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		benchmarkStrScopeFQPath, _ = path.FQ(scope.StrWebsites, "4", "system", "dev", "debug")
	}
	if benchmarkStrScopeFQPath != want {
		b.Errorf("Want: %s; Have, %s", want, benchmarkStrScopeFQPath)
	}
}

func benchmarkStrScopeFQPathInt64(scopeID int64, b *testing.B) {
	want := scope.StrWebsites.String() + "/" + strconv.FormatInt(scopeID, 10) + "/system/dev/debug"
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		benchmarkStrScopeFQPath, _ = path.FQInt64(scope.StrWebsites, scopeID, "system", "dev", "debug")
	}
	if benchmarkStrScopeFQPath != want {
		b.Errorf("Want: %s; Have, %s", want, benchmarkStrScopeFQPath)
	}
}

// BenchmarkStrScopeFQPathInt64__Cached-4	 5000000	       383 ns/op	      32 B/op	       1 allocs/op
func BenchmarkStrScopeFQPathInt64__Cached(b *testing.B) {
	benchmarkStrScopeFQPathInt64(4, b)
}

// BenchmarkStrScopeFQPathInt64UnCached-4	 3000000	       438 ns/op	      48 B/op	       2 allocs/op
func BenchmarkStrScopeFQPathInt64UnCached(b *testing.B) {
	benchmarkStrScopeFQPathInt64(40, b)
}

func TestSplitFQPath(t *testing.T) {
	t.Parallel()
	tests := []struct {
		have        string
		wantScope   string
		wantScopeID int64
		wantPath    string
		wantErr     error
	}{
		{"groups/1/catalog/frontend/list_allow_all", "groups", 0, "", scope.ErrUnsupportedScope},
		{"stores/7475/catalog/frontend/list_allow_all", scope.StrStores.String(), 7475, "catalog/frontend/list_allow_all", nil},
		{"websites/1/catalog/frontend/list_allow_all", scope.StrWebsites.String(), 1, "catalog/frontend/list_allow_all", nil},
		{"default/0/catalog/frontend/list_allow_all", scope.StrDefault.String(), 0, "catalog/frontend/list_allow_all", nil},
		{"default//catalog/frontend/list_allow_all", scope.StrDefault.String(), 0, "catalog/frontend/list_allow_all", errors.New("strconv.ParseInt: parsing \"\\uf8ff\": invalid syntax")},
		{"stores/123/catalog/index", "", 0, "", errors.New("Incorrect fully qualified path: \"stores/123/catalog/index\"")},
	}
	for _, test := range tests {
		haveScope, haveScopeID, havePath, haveErr := path.SplitFQ(test.have)

		if test.wantErr != nil {
			assert.EqualError(t, haveErr, test.wantErr.Error(), "Test %v", test)
		} else {
			assert.NoError(t, haveErr, "Test %v", test)
		}
		assert.Exactly(t, test.wantScope, haveScope, "Test %v", test)
		assert.Exactly(t, test.wantScopeID, haveScopeID, "Test %v", test)
		assert.Exactly(t, test.wantPath, havePath, "Test %v", test)
	}
}

var benchmarkReverseFQPath = struct {
	scope   string
	scopeID int64
	path    string
	err     error
}{}

// BenchmarkReverseFQPath-4 	10000000	       121 ns/op	       0 B/op	       0 allocs/op
func BenchmarkReverseFQPath(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		benchmarkReverseFQPath.scope, benchmarkReverseFQPath.scopeID, benchmarkReverseFQPath.path, benchmarkReverseFQPath.err = path.SplitFQ("stores/7475/catalog/frontend/list_allow_all")
		if benchmarkReverseFQPath.err != nil {
			b.Error(benchmarkReverseFQPath.err)
		}
	}
}

func TestJoin(t *testing.T) {
	t.Parallel()
	tests := []struct {
		have []string
		want string
	}{
		{[]string{"a", "b", "c"}, "a/b/c"},
		{[]string{"a/b", "c"}, "a/b/c"},
		{[]string{"a", "b/c"}, "a/b/c"},
		{[]string{"a/b/c"}, "a/b/c"},
	}
	for i, test := range tests {
		assert.Exactly(t, test.want, path.Join(test.have...), "Index %d", i)
	}
}

func TestSplit(t *testing.T) {
	t.Parallel()
	tests := []struct {
		have string
		want []string
	}{
		{"a/b/c", []string{"a", "b", "c"}},
	}
	for i, test := range tests {
		assert.Exactly(t, test.want, path.Split(test.have), "Index %d", i)
	}
}
