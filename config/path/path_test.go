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
	"bytes"
	"errors"
	"strconv"
	"testing"

	"github.com/corestoreio/csfw/config/path"
	"github.com/corestoreio/csfw/store/scope"
	"github.com/stretchr/testify/assert"
)

func TestPathNew(t *testing.T) {
	t.Parallel()
	tests := []struct {
		route      path.Route
		s          scope.Scope
		id         int64
		wantFQ     path.Route
		wantNewErr error
	}{
		{path.Route("ab/ba/cd"), scope.WebsiteID, 3, path.Route("websites/3/ab/ba/cd"), nil},
		{path.Route("ad/ba/ca/sd"), scope.WebsiteID, 3, path.Route("websites/3/a/b/c/d"), path.ErrIncorrectPath},
		{path.Route("as/sb"), scope.WebsiteID, 3, path.Route("websites/3/a/b/c/d"), path.ErrIncorrectPath},
	}
	for i, test := range tests {
		haveP, haveErr := path.New(test.route)
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
		route   path.Route
		want    string
		wantErr error
	}{
		{scope.StrDefault, 0, nil, "", path.ErrRouteEmpty},
		{scope.StrDefault, 0, path.Route(""), "", path.ErrRouteEmpty},
		{scope.StrDefault, 0, path.Route("system/dev/debug"), scope.StrDefault.String() + "/0/system/dev/debug", nil},
		{scope.StrDefault, 44, path.Route("system/dev/debug"), scope.StrDefault.String() + "/0/system/dev/debug", nil},
		{scope.StrWebsites, 0, path.Route("system/dev/debug"), scope.StrWebsites.String() + "/0/system/dev/debug", nil},
		{scope.StrWebsites, 343, path.Route("system/dev/debug"), scope.StrWebsites.String() + "/343/system/dev/debug", nil},
		{scope.StrScope("hello"), 0, path.Route("system/dev/debug"), scope.StrDefault.String() + "/0/system/dev/debug", nil},
		{scope.StrScope("hello"), 343, path.Route("system/dev/debug"), scope.StrDefault.String() + "/0/system/dev/debug", nil},
	}
	for i, test := range tests {
		p, pErr := path.New(test.route)
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
		assert.Exactly(t, test.want, have.String(), "Index %d", i)
	}

	r := path.Route("catalog/frontend/list_allow_all")
	assert.Exactly(t, "stores/7475/catalog/frontend/list_allow_all", path.MustNew(r).BindStr(scope.StrStores, 7475).String())
	assert.Exactly(t, "stores/5/catalog/frontend/list_allow_all", path.MustNew(r).BindStr(scope.StrStores, 5).String())
}

func TestShouldNotPanicBecauseOfIncorrectStrScope(t *testing.T) {
	t.Parallel()
	assert.Exactly(t, "stores/345/xxxxx/yyyyy/zzzzz", path.MustNew(path.Route("xxxxx/yyyyy/zzzzz")).BindStr(scope.StrStores, 345).String())
	defer func() {
		if r := recover(); r != nil {
			t.Fatal("Did not expect a panic")
		}
	}()
	_ = path.MustNew(path.Route("xxxxx/yyyyy/zzzzz")).BindStr(scope.StrScope("invalid"), 345)
}

func TestShouldPanicIncorrectPath(t *testing.T) {
	t.Parallel()
	assert.Exactly(t, "default/0/xxxxx/yyyyy/zzzzz", path.MustNew(path.Route("xxxxx/yyyyy/zzzzz")).BindStr(scope.StrDefault, 345).String())
	defer func() {
		if r := recover(); r != nil {
			assert.EqualError(t, r.(error), path.ErrIncorrectPath.Error())
		} else {
			t.Fatal("Expecting a panic")
		}
	}()
	assert.Exactly(t, "websites/345/xxxxx/yyyyy", path.MustNew(path.Route("xxxxx/yyyyy")).BindStr(scope.StrWebsites, 345).String())
}

var benchmarkStrScopeFQPath path.Route

func benchmarkFQ(scopeID int64, b *testing.B) {
	want := path.Route(scope.StrWebsites.String() + "/" + strconv.FormatInt(scopeID, 10) + "/system/dev/debug")
	p := path.Route("system/dev/debug")
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		var err error
		benchmarkStrScopeFQPath, err = path.MustNew(p).BindStr(scope.StrWebsites, scopeID).FQ()
		if err != nil {
			b.Error(err)
		}
	}
	if bytes.Compare(benchmarkStrScopeFQPath, want) != 0 {
		b.Errorf("Want: %s; Have, %s", want, benchmarkStrScopeFQPath)
	}
}

// BenchmarkFQ-4     	 3000000	       401 ns/op	     112 B/op	       1 allocs/op
func BenchmarkFQ(b *testing.B) {
	benchmarkFQ(11, b)
}

func TestSplitFQ(t *testing.T) {
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
		havePath, haveErr := path.SplitFQ(path.Route(test.have))

		if test.wantErr != nil {
			assert.EqualError(t, haveErr, test.wantErr.Error(), "Test %v", test)
		} else {
			assert.NoError(t, haveErr, "Test %v", test)
		}
		assert.Exactly(t, test.wantScope, havePath.StrScope(), "Test %v", test)
		assert.Exactly(t, test.wantScopeID, havePath.ID, "Test %v", test)
		l, _ := havePath.Level(-1)
		assert.Exactly(t, test.wantPath, l.String(), "Test %v", test)
	}
}

var benchmarkSplitFQ path.Path

// BenchmarkSplitFQ-4	10000000	       175 ns/op	      16 B/op	       1 allocs/op strings
// BenchmarkSplitFQ-4  	10000000	       186 ns/op	      32 B/op	       1 allocs/op strings
// BenchmarkSplitFQ-4  	 5000000	       364 ns/op	      16 B/op	       1 allocs/op bytes
func BenchmarkSplitFQ(b *testing.B) {
	r := path.Route("stores/7475/catalog/frontend/list_allow_all")
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		var err error
		benchmarkSplitFQ, err = path.SplitFQ(r)
		if err != nil {
			b.Error(err)
		}
	}
}

func TestLevel(t *testing.T) {
	t.Parallel()
	tests := []struct {
		have  path.Route
		level int
		want  string
	}{
		{path.Route("general/single_store_mode/enabled"), 0, ""},
		{path.Route("general/single_store_mode/enabled"), 1, "general"},
		{path.Route("general/single_store_mode/enabled"), 2, "general/single_store_mode"},
		{path.Route("general/single_store_mode/enabled"), 3, "general/single_store_mode/enabled"},
		{path.Route("general/single_store_mode/enabled"), -1, "general/single_store_mode/enabled"},
		{path.Route("general/single_store_mode/enabled"), 5, "general/single_store_mode/enabled"},
		{path.Route("general/single_store_mode/enabled"), 4, "general/single_store_mode/enabled"},
	}
	for i, test := range tests {
		r, err := path.MustNew(test.have).Level(test.level)
		assert.NoError(t, err)
		assert.Exactly(t, test.want, r.String(), "Index %d", i)
	}
}

var benchmarkLevel path.Route

// BenchmarkLevel_One-4	 5000000	       297 ns/op	      16 B/op	       1 allocs/op
func BenchmarkLevel_One(b *testing.B) {
	benchmarkLevelRun(b, 1, path.Route("system/dev/debug"), path.Route("system"))
}

// BenchmarkLevel_Two-4	 5000000	       332 ns/op	      16 B/op	       1 allocs/op
func BenchmarkLevel_Two(b *testing.B) {
	benchmarkLevelRun(b, 2, path.Route("system/dev/debug"), path.Route("system/dev"))
}

// BenchmarkLevel_All-4	 5000000	       379 ns/op	      16 B/op	       1 allocs/op
func BenchmarkLevel_All(b *testing.B) {
	benchmarkLevelRun(b, -1, path.Route("system/dev/debug"), path.Route("system/dev/debug"))
}

func benchmarkLevelRun(b *testing.B, level int, have, want path.Route) {
	b.ReportAllocs()
	b.ResetTimer()
	var err error
	for i := 0; i < b.N; i++ {
		benchmarkLevel, err = path.MustNew(have).Level(level)
	}
	if err != nil {
		b.Error(err)
	}
	if bytes.Compare(benchmarkLevel, want) != 0 {
		b.Errorf("Want: %s; Have, %s", want, benchmarkLevel)
	}
}

func TestIsValid(t *testing.T) {
	tests := []struct {
		s    scope.Scope
		id   int64
		have path.Route
		want error
	}{
		{scope.DefaultID, 0, path.Route("//"), path.ErrIncorrectPath},
		{scope.DefaultID, 0, path.Route("general/store_information/city"), nil},
		{scope.DefaultID, 33, path.Route("general/store_information/city"), nil},
		{scope.DefaultID, 0, path.Route(""), path.ErrRouteEmpty},
		{scope.DefaultID, 0, path.Route("general/store_information"), path.ErrIncorrectPath},
		////{path.Route(path.MustNew("system/dev/debug").Bind(scope.WebsiteID, 22).String()), path.ErrIncorrectPath},
		{scope.DefaultID, 0, path.Route("groups/33/general/store_information/street"), path.ErrIncorrectPath},
		{scope.DefaultID, 0, path.Route("groups/33"), path.ErrIncorrectPath},
		{scope.DefaultID, 0, path.Route("system/dEv/inv˚lid"), errors.New("This character \"˚\" is not allowed in Route system/dEv/inv˚lid")},
		{scope.DefaultID, 0, path.Route("system/dEv/inv'lid"), errors.New("This character \"'\" is not allowed in Route system/dEv/inv'lid")},
		{scope.DefaultID, 0, path.Route("syst3m/dEv/invalid"), nil},
		{scope.DefaultID, 0, nil, path.ErrRouteEmpty},
	}
	for i, test := range tests {
		p := path.Path{
			Scope: test.s,
			ID:    test.id,
			Route: test.have,
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

// BenchmarkIsValid-4	20000000	        83.5 ns/op	       0 B/op	       0 allocs/op
func BenchmarkIsValid(b *testing.B) {
	have := path.Route("system/dEv/d3bug")
	want := "system/dev/debug"

	p := path.Path{
		Route: have,
	}
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		benchmarkIsValid = p.IsValid()
		if nil != benchmarkIsValid {
			b.Errorf("Want: %s; Have: %v", want, p.Route)
		}
	}
}
