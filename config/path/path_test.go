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
	"reflect"
	"strconv"
	"testing"

	"github.com/corestoreio/csfw/config/path"
	"github.com/corestoreio/csfw/store/scope"
	"github.com/stretchr/testify/assert"
)

func TestPathSplit(t *testing.T) {
	t.Parallel()
	tests := []struct {
		parts   []string
		want    string
		wantErr error
	}{
		{[]string{"general", "single_store_mode", "enabled"}, "general/single_store_mode/enabled", nil},
		{[]string{"general", "single_store_mode"}, "general/single_store_mode/enabled", errors.New("Incorrect number of paths elements: want 3, have 2, Path: [general single_store_mode]")},
		{[]string{"general/single_store_mode/enabled"}, "general/single_store_mode/enabled", nil},
		{[]string{"general/singlestore_mode/enabled"}, "general/single_store_mode/enabled", errors.New("This character \"\\uf8ff\" is not allowed in Parts []string{\"general\", \"single\\uf8ffstore_mode\", \"enabled\"}")},
	}
	for i, test := range tests {
		haveP, haveErr := path.NewSplit(test.parts...)
		if test.wantErr != nil {
			assert.EqualError(t, haveErr, test.wantErr.Error(), "Index %d", i)
			continue
		}
		assert.Exactly(t, test.want, haveP.Level(0), "Index %d", i)
	}
}

func TestPath(t *testing.T) {
	t.Parallel()
	tests := []struct {
		parts      []string
		s          scope.Scope
		id         int64
		wantFQ     string
		wantNewErr error
	}{
		{[]string{"ab/ba/cd"}, scope.WebsiteID, 3, "websites/3/ab/ba/cd", nil},
		{[]string{"ad/ba/ca/sd"}, scope.WebsiteID, 3, "websites/3/a/b/c/d", path.ErrIncorrectPath},
		{[]string{"as/sb"}, scope.WebsiteID, 3, "websites/3/a/b/c/d", path.ErrIncorrectPath},
	}
	for i, test := range tests {
		haveP, haveErr := path.New(test.parts...)
		haveP = haveP.Bind(test.s, test.id)
		if test.wantNewErr != nil {
			assert.EqualError(t, haveErr, test.wantNewErr.Error(), "Index %d", i)
			continue
		}
		fq, fqErr := haveP.FQ()
		assert.NoError(t, fqErr, "Index %d", i)
		assert.Exactly(t, test.wantFQ, fq, "Index %d", i)
	}
}

func TestFQ(t *testing.T) {
	t.Parallel()
	tests := []struct {
		str     scope.StrScope
		id      int64
		path    []string
		want    string
		wantErr error
	}{
		{scope.StrDefault, 0, nil, "", errors.New("Parts are empty")},
		{scope.StrDefault, 0, []string{}, "", errors.New("Parts are empty")},
		{scope.StrDefault, 0, []string{""}, "", path.ErrIncorrectPath},
		{scope.StrDefault, 0, []string{"system/dev/debug"}, scope.StrDefault.String() + "/0/system/dev/debug", nil},
		{scope.StrDefault, 33, []string{"system", "dev", "debug"}, scope.StrDefault.String() + "/0/system/dev/debug", nil},
		{scope.StrWebsites, 0, []string{"system/dev/debug"}, scope.StrWebsites.String() + "/0/system/dev/debug", nil},
		{scope.StrWebsites, 343, []string{"system", "dev", "debug"}, scope.StrWebsites.String() + "/343/system/dev/debug", nil},
		{scope.StrScope("hello"), 343, []string{"system", "dev", "debug"}, scope.StrDefault.String() + "/0/system/dev/debug", nil},
	}
	for i, test := range tests {
		p, pErr := path.New(test.path...)
		p = p.BindStr(test.str, test.id)
		have, haveErr := p.FQ()
		if test.wantErr != nil {
			assert.Empty(t, have, "Index %d", i)
			if pErr != nil {
				assert.EqualError(t, pErr, test.wantErr.Error(), "Index %d", i)
				continue
			}
			assert.EqualError(t, haveErr, test.wantErr.Error(), "Index %d", i)
			continue
		}
		assert.NoError(t, haveErr, "Index %d", i)
		assert.Equal(t, test.want, have, "Index %d", i)
	}
	assert.Equal(t, "stores/7475/catalog/frontend/list_allow_all", path.MustNew("catalog", "frontend", "list_allow_all").BindStr(scope.StrStores, 7475).String())
	assert.Equal(t, "stores/5/catalog/frontend/list_allow_all", path.MustNew("catalog", "frontend", "list_allow_all").BindStr(scope.StrStores, 5).String())
}

func TestShouldNotPanicBecauseOfIncorrectStrScope(t *testing.T) {
	t.Parallel()
	assert.Exactly(t, "stores/345/xxxxx/yyyyy/zzzzz", path.MustNew("xxxxx", "yyyyy", "zzzzz").BindStr(scope.StrStores, 345).String())
	defer func() {
		if r := recover(); r != nil {
			t.Fatal("Did not expect a panic")
		}
	}()
	_ = path.MustNew("xxxxx", "yyyyy", "zzzzz").BindStr(scope.StrScope("invalid"), 345)
}

func TestShouldPanicIncorrectPath(t *testing.T) {
	t.Parallel()
	assert.Exactly(t, "default/0/xxxxx/yyyyy/zzzzz", path.MustNew("xxxxx", "yyyyy", "zzzzz").BindStr(scope.StrDefault, 345).String())
	defer func() {
		if r := recover(); r != nil {
			assert.EqualError(t, r.(error), "All arguments must be valid! Min want: 3. Have: 2. Parts []string{\"xxxxx\", \"yyyyy\"}")
		} else {
			t.Fatal("Expecting a panic")
		}
	}()
	assert.Exactly(t, "websites/345/xxxxx/yyyyy", path.MustNew("xxxxx", "yyyyy").BindStr(scope.StrWebsites, 345).String())
}

var benchmarkStrScopeFQPath string

func benchmarkFQ(scopeID int64, b *testing.B) {
	want := scope.StrWebsites.String() + "/" + strconv.FormatInt(scopeID, 10) + "/system/dev/debug"
	p := []string{"system/dev/debug"}
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		var err error
		benchmarkStrScopeFQPath, err = path.MustNew(p...).BindStr(scope.StrWebsites, scopeID).FQ()
		if err != nil {
			b.Error(err)
		}
	}
	if benchmarkStrScopeFQPath != want {
		b.Errorf("Want: %s; Have, %s", want, benchmarkStrScopeFQPath)
	}
}

// BenchmarkFQ__Cached-4	 3000000	       427 ns/op	      32 B/op	       1 allocs/op
func BenchmarkFQ__Cached(b *testing.B) {
	benchmarkFQ(4, b)
}

// BenchmarkFQUnCached-4	 3000000	       513 ns/op	      48 B/op	       2 allocs/op
func BenchmarkFQUnCached(b *testing.B) {
	benchmarkFQ(40, b)
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
		{"groups/1/catalog/frontend/list_allow_all", "default", 0, "", scope.ErrUnsupportedScope},
		{"stores/7475/catalog/frontend/list_allow_all", scope.StrStores.String(), 7475, "catalog/frontend/list_allow_all", nil},
		{"websites/1/catalog/frontend/list_allow_all", scope.StrWebsites.String(), 1, "catalog/frontend/list_allow_all", nil},
		{"default/0/catalog/frontend/list_allow_all", scope.StrDefault.String(), 0, "catalog/frontend/list_allow_all", nil},
		{"default//catalog/frontend/list_allow_all", scope.StrDefault.String(), 0, "catalog/frontend/list_allow_all", errors.New("strconv.ParseInt: parsing \"\\uf8ff\": invalid syntax")},
		{"stores/123/catalog/index", "default", 0, "", errors.New("Incorrect fully qualified path: \"stores/123/catalog/index\"")},
	}
	for _, test := range tests {
		havePath, haveErr := path.SplitFQ(test.have)

		if test.wantErr != nil {
			assert.EqualError(t, haveErr, test.wantErr.Error(), "Test %v", test)
		} else {
			assert.NoError(t, haveErr, "Test %v", test)
		}
		assert.Exactly(t, test.wantScope, havePath.StrScope(), "Test %v", test)
		assert.Exactly(t, test.wantScopeID, havePath.ID, "Test %v", test)
		assert.Exactly(t, test.wantPath, havePath.Level(-1), "Test %v", test)
	}
}

var benchmarkReverseFQPath path.Path

// BenchmarkReverseFQPath-4	10000000	       175 ns/op	      16 B/op	       1 allocs/op
func BenchmarkReverseFQPath(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		var err error
		benchmarkReverseFQPath, err = path.SplitFQ("stores/7475/catalog/frontend/list_allow_all")
		if err != nil {
			b.Error(err)
		}
	}
}

func TestLevel(t *testing.T) {
	t.Parallel()
	tests := []struct {
		have  []string
		level int
		want  string
	}{
		{[]string{"general", "single_store_mode", "enabled"}, -1, "general/single_store_mode/enabled"},
		{[]string{"general", "single_store_mode", "enabled"}, 0, "general/single_store_mode/enabled"},
		{[]string{"general", "single_store_mode", "enabled"}, 5, "general/single_store_mode/enabled"},
		{[]string{"general", "single_store_mode", "enabled"}, 4, "general/single_store_mode/enabled"},
		{[]string{"general", "single_store_mode", "enabled"}, 3, "general/single_store_mode/enabled"},
		{[]string{"general", "single_store_mode", "enabled"}, 2, "general/single_store_mode"},
		{[]string{"general", "single_store_mode", "enabled"}, 1, "general"},
		{[]string{"general/single_store_mode/enabled"}, -1, "general/single_store_mode/enabled"},
	}
	for i, test := range tests {
		assert.Exactly(t, test.want, path.MustNew(test.have...).Level(test.level), "Index %d", i)
	}
}

var benchmarkLevel string

// BenchmarkLevel-4           	10000000	       238 ns/op	      16 B/op	       1 allocs/op => Go 1.5.2
func BenchmarkLevel(b *testing.B) {
	have := []string{"system", "dev", "debug"}
	want := "system/dev/debug"

	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		benchmarkLevel = path.MustNew(have...).Level(-1)
	}
	if benchmarkLevel != want {
		b.Errorf("Want: %s; Have, %s", want, benchmarkLevel)
	}
}

func TestSplit(t *testing.T) {
	t.Parallel()
	tests := []struct {
		have string
		want []string
	}{
		{"system/dev/debug", []string{"system", "dev", "debug"}},
		{"a/b", []string{"a", "b"}},
		{"a/b/c", []string{"a", "b", "c"}},
		{"/a/b/c", []string{"a", "b", "c"}},
		{"a/b/c/d/e", []string{"a", "b", "c", "d", "e"}},
	}
	for i, test := range tests {
		assert.Exactly(t, test.want, path.Split(test.have), "Index %d", i)
	}
}

var benchmarkSplit []string

// BenchmarkSplit-4          	 5000000	       290 ns/op	      48 B/op	       1 allocs/op => Go 1.5.2
func BenchmarkSplit(b *testing.B) {
	want := []string{"system", "dev", "debug"}
	have := "system/dev/debug"

	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		benchmarkSplit = path.Split(have)
	}
	if false == reflect.DeepEqual(benchmarkSplit, want) {
		b.Errorf("Want: %s; Have, %s", want, benchmarkLevel)
	}
}

func TestIsValid(t *testing.T) {
	tests := []struct {
		have []string
		want error
	}{
		{[]string{"//"}, path.ErrIncorrectPath},
		{[]string{"general/store_information/city"}, nil},
		{[]string{"", "", ""}, errors.New("This path part \"\" is too short. Parts: []string{\"\", \"\", \"\"}")},
		{[]string{"general", "store_information", "name"}, nil},
		{[]string{"general", "store_information"}, errors.New("All arguments must be valid! Min want: 3. Have: 2. Parts []string{\"general\", \"store_information\"}")},
		{[]string{path.MustNew("system", "dev", "debug").Bind(scope.WebsiteID, 22).String()}, path.ErrIncorrectPath},
		{[]string{"groups/33/general/store_information/street"}, path.ErrIncorrectPath},
		{[]string{"groups/33"}, path.ErrIncorrectPath},
		{[]string{"system/dEv/inv˚lid"}, errors.New("This character \"˚\" is not allowed in Parts []string{\"system/dEv/inv˚lid\"}")},
		{[]string{"syst3m/dEv/invalid"}, nil},
		{nil, errors.New("Parts are empty")},
	}
	for i, test := range tests {
		p := path.Path{
			Parts: test.have,
		}
		haveErr := p.IsValid()
		if test.want != nil {
			assert.EqualError(t, haveErr, test.want.Error(), "Index %d", i)
		} else {
			assert.NoError(t, haveErr, "Index %d", i)
		}
	}
}

var benchmarkIsValid error

// BenchmarkIsValid-4        	20000000	       116 ns/op	       0 B/op	       0 allocs/op
func BenchmarkIsValid(b *testing.B) {
	have := []string{"system", "dEv", "d3bug"}
	want := "system/dev/debug"

	p := path.Path{
		Parts: have,
	}
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		benchmarkIsValid = p.IsValid()
		if nil != benchmarkIsValid {
			b.Errorf("Want: %s; Have: %v", want, p.Parts)
		}
	}
}
