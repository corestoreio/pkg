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
	"database/sql"
	"database/sql/driver"
	"encoding"
	"encoding/json"
	"errors"
	"fmt"
	"hash/fnv"
	"testing"

	"github.com/corestoreio/csfw/config/path"
	"github.com/stretchr/testify/assert"
)

// These checks, if a type implements an interface, belong into the test package
// and not into its "main" package. Otherwise you would also compile each time
// all the packages with their interfaces, etc.
var _ encoding.TextMarshaler = (*path.Route)(nil)
var _ encoding.TextUnmarshaler = (*path.Route)(nil)
var _ sql.Scanner = (*path.Route)(nil)
var _ driver.Valuer = (*path.Route)(nil)
var _ fmt.GoStringer = (*path.Route)(nil)
var _ fmt.Stringer = (*path.Route)(nil)
var _ path.Router = (*path.Route)(nil)

func TestRouteRouter(t *testing.T) {
	t.Parallel()
	r := path.NewRoute("a/b/c")
	assert.Exactly(t, r, r.Route())
}

func TestRouteGoString(t *testing.T) {
	t.Parallel()
	tests := []struct {
		have path.Route
		want string
	}{
		{path.NewRoute("a"), "path.Route{Chars:[]byte(`a`)}"},
		{path.NewRoute(""), "path.Route{}"},
		{path.Route{}, "path.Route{}"},
	}
	for i, test := range tests {
		assert.Exactly(t, test.want, test.have.GoString(), "Index %d", i)
	}
}

func TestRouteEqual(t *testing.T) {
	t.Parallel()
	tests := []struct {
		a    path.Route
		b    path.Route
		want bool
	}{
		{path.Route{}, path.Route{}, true},
		{path.NewRoute("a"), path.NewRoute("a"), true},
		{path.NewRoute("a"), path.NewRoute("b"), false},
		{path.NewRoute("a\x80"), path.NewRoute("a"), false},
		{path.NewRoute("general/single_\x80store_mode/enabled"), path.NewRoute("general/single_store_mode/enabled"), false},
		{path.NewRoute("general/single_store_mode/enabled"), path.NewRoute("general/single_store_mode/enabled"), true},
		{path.NewRoute(""), path.NewRoute(""), true},
		{path.NewRoute(""), path.NewRoute(), true},
	}
	for i, test := range tests {
		assert.Exactly(t, test.want, test.a.Equal(test.b), "Index %d", i)
	}
}

func TestRouteAppend(t *testing.T) {
	t.Parallel()
	tests := []struct {
		a       path.Route
		b       path.Route
		want    string
		wantErr error
	}{
		{path.NewRoute("aa"), path.NewRoute("bb/cc"), "aa/bb/cc", nil},
		{path.NewRoute("aa"), path.NewRoute("bbcc"), "aa/bbcc", nil},
		{path.NewRoute("aa"), path.NewRoute("bb\x80cc"), "", path.ErrRouteInvalidBytes},
		{path.NewRoute("aa/"), path.NewRoute("bbcc"), "aa/bbcc", nil},
		{path.NewRoute("aa"), path.NewRoute("/bbcc"), "aa/bbcc", nil},
		{path.NewRoute("aa"), path.NewRoute("/bb\x00cc"), "aa/bb", nil},
		{path.NewRoute("ag"), path.NewRoute("b"), "ag/b", nil},
		{path.NewRoute("c/"), path.NewRoute("/b"), "c/b", nil},
		{path.NewRoute("d/"), path.NewRoute(""), "d", nil},
		{path.NewRoute("e"), path.NewRoute(""), "e", nil},
		{path.NewRoute(""), path.NewRoute("f"), "f", nil},
	}
	for i, test := range tests {
		haveErr := test.a.Append(test.b)
		if test.wantErr != nil {
			assert.EqualError(t, haveErr, test.wantErr.Error(), "Index %d", i)
			continue
		}
		assert.NoError(t, haveErr, "Index %d", i)
		assert.Exactly(t, test.want, test.a.String(), "Index %d", i)
	}
}

func TestRouteVariadicAppend(t *testing.T) {
	t.Parallel()

	tests := []struct {
		a       path.Route
		routes  []path.Route
		want    string
		wantErr error
	}{
		{path.NewRoute("aa"), []path.Route{path.NewRoute("bb"), path.NewRoute("cc"), path.NewRoute("dd")}, "aa/bb/cc/dd", nil},
		{path.NewRoute("aa"), []path.Route{{}, {}, {Chars: []byte(`cb`)}}, "aa/cb", nil},
		{path.Route{}, []path.Route{{}, {}, {Chars: []byte(`cb`)}}, "cb", nil},
		{path.Route{}, []path.Route{{}, {}, {}}, "", path.ErrRouteEmpty},
	}
	for i, test := range tests {
		haveErr := test.a.Append(test.routes...)
		if test.wantErr != nil {
			assert.EqualError(t, haveErr, test.wantErr.Error(), "Index %d", i)
			continue
		}
		assert.NoError(t, haveErr, "Index %d", i)
		assert.Exactly(t, test.want, test.a.String(), "Index %d", i)
	}

}

var benchmarkRouteAppendWant = path.NewRoute("general/single_store_mode/enabled")

// BenchmarkRouteAppend-4   	 5000000	       376 ns/op	      64 B/op	       2 allocs/op
func BenchmarkRouteAppend(b *testing.B) {
	havePartials := []path.Route{path.NewRoute("single_store_mode"), path.NewRoute("enabled")}
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		var have = path.NewRoute("general/")
		err := have.Append(havePartials...)
		if err != nil {
			b.Error(err)
		}
		if benchmarkRouteAppendWant.Equal(have) == false {
			b.Errorf("Want: %s; Have: %s", benchmarkRouteAppendWant, have)
		}
	}
}

func TestRouteTextMarshal(t *testing.T) {
	t.Parallel()
	r := path.NewRoute("admin/security/password_lifetime")
	j, err := json.Marshal(r)
	assert.NoError(t, err)
	assert.Exactly(t, "\"admin/security/password_lifetime\"", string(j))
}

func TestRouteUnmarshalTextOk(t *testing.T) {
	t.Parallel()
	var r path.Route
	err := json.Unmarshal([]byte(`"admin/security/password_lifetime"`), &r)
	assert.NoError(t, err)
	assert.Exactly(t, "admin/security/password_lifetime", r.String())
}

func TestRouteLevel(t *testing.T) {
	t.Parallel()
	tests := []struct {
		have  path.Route
		level int
		want  string
	}{
		{path.NewRoute("general/single_store_mode/enabled"), 0, ""},
		{path.NewRoute("general/single_store_mode/enabled"), 1, "general"},
		{path.NewRoute("general/single_store_mode/enabled"), 2, "general/single_store_mode"},
		{path.NewRoute("general/single_store_mode/enabled"), 3, "general/single_store_mode/enabled"},
		{path.NewRoute("general/single_store_mode/enabled"), -1, "general/single_store_mode/enabled"},
		{path.NewRoute("general/single_store_mode/enabled"), 5, "general/single_store_mode/enabled"},
		{path.NewRoute("general/single_store_mode/enabled"), 4, "general/single_store_mode/enabled"},
		{path.NewRoute("system/full_page_cache/varnish/backend_port"), 3, "system/full_page_cache/varnish"},
	}
	for i, test := range tests {
		r, err := test.have.Level(test.level)
		assert.NoError(t, err)
		assert.Exactly(t, test.want, r.String(), "Index %d", i)
	}
}

var benchmarkRouteLevel path.Route

// BenchmarkRouteLevel_One-4	 5000000	       297 ns/op	      16 B/op	       1 allocs/op
func BenchmarkRouteLevel_One(b *testing.B) {
	benchmarkRouteLevelRun(b, 1, path.NewRoute("system/dev/debug"), path.NewRoute("system"))
}

// BenchmarkRouteLevel_Two-4	 5000000	       332 ns/op	      16 B/op	       1 allocs/op
func BenchmarkRouteLevel_Two(b *testing.B) {
	benchmarkRouteLevelRun(b, 2, path.NewRoute("system/dev/debug"), path.NewRoute("system/dev"))
}

// BenchmarkRouteLevel_All-4	 5000000	       379 ns/op	      16 B/op	       1 allocs/op
func BenchmarkRouteLevel_All(b *testing.B) {
	benchmarkRouteLevelRun(b, -1, path.NewRoute("system/dev/debug"), path.NewRoute("system/dev/debug"))
}

func benchmarkRouteLevelRun(b *testing.B, level int, have, want path.Route) {
	b.ReportAllocs()
	b.ResetTimer()
	var err error
	for i := 0; i < b.N; i++ {
		benchmarkRouteLevel, err = have.Level(level)
	}
	if err != nil {
		b.Error(err)
	}
	if benchmarkRouteLevel.Equal(want) == false {
		b.Errorf("Want: %s; Have, %s", want, benchmarkRouteLevel)
	}
}

func TestRouteHash(t *testing.T) {
	t.Parallel()
	tests := []struct {
		have      path.Route
		level     int
		wantHash  uint32
		wantErr   error
		wantLevel string
	}{
		{path.NewRoute("general/single_\x80store_mode/enabled"), 0, 0, path.ErrRouteInvalidBytes, ""},
		{path.NewRoute("general/single_store_mode/enabled"), 0, 2166136261, nil, ""},
		{path.NewRoute("general/single_store_mode/enabled"), 1, 616112491, nil, "general"},
		{path.NewRoute("general/single_store_mode/enabled"), 2, 2274889228, nil, "general/single_store_mode"},
		{path.NewRoute("general/single_store_mode/enabled"), 3, 1644245266, nil, "general/single_store_mode/enabled"},
		{path.NewRoute("general/single_store_mode/enabled"), -1, 1644245266, nil, "general/single_store_mode/enabled"},
		{path.NewRoute("general/single_store_mode/enabled"), 5, 1644245266, nil, "general/single_store_mode/enabled"},
		{path.NewRoute("general/single_store_mode/enabled"), 4, 1644245266, nil, "general/single_store_mode/enabled"},
	}
	for i, test := range tests {

		hv, err := test.have.Hash(test.level)
		if test.wantErr != nil {
			assert.EqualError(t, err, test.wantErr.Error(), "Index %d", i)
			assert.Empty(t, hv, "Index %d", i)
			continue
		}
		assert.NoError(t, err, "Index %d", i)

		check := fnv.New32a()
		_, cErr := check.Write([]byte(test.wantLevel))
		assert.NoError(t, cErr)
		assert.Exactly(t, check.Sum32(), hv, "Have %d Want %d Index %d", check.Sum32(), hv, i)

		l, err := test.have.Level(test.level)
		assert.Exactly(t, test.wantLevel, l.String(), "Index %d", i)
		assert.Exactly(t, test.wantHash, hv, "Have %d Want %d Index %d", test.wantHash, hv, i)
	}
}

func TestRouteHash32(t *testing.T) {
	t.Parallel()
	tests := []struct {
		have     path.Route
		wantHash uint32
	}{
		{path.NewRoute("general/single_\x80store_mode/enabled"), 1310908924},
		{path.NewRoute("general/single_store_mode/enabled"), 1644245266},
		{path.NewRoute(""), 2166136261},
		{path.Route{}, 2166136261},
	}
	for i, test := range tests {
		if test.have.Sum32 > 0 {
			assert.Exactly(t, test.wantHash, test.have.Sum32, "Want %d Have %d Index %d => %s", test.wantHash, test.have.Sum32, i, test.have)
		}
		hv := test.have.Hash32()
		check := fnv.New32a()
		_, cErr := check.Write(test.have.Bytes())
		assert.NoError(t, cErr)
		assert.Exactly(t, check.Sum32(), hv, "Have %d Want %d Index %d", check.Sum32(), hv, i)
		assert.Exactly(t, test.wantHash, hv, "Have %d Want %d Index %d", test.wantHash, hv, i)
	}
}

var benchmarkRouteHash uint32

// BenchmarkRouteHash-4     	 5000000	       287 ns/op	       0 B/op	       0 allocs/op
func BenchmarkRouteHash(b *testing.B) {
	have := path.NewRoute("general/single_store_mode/enabled")
	want := uint32(1644245266)

	var err error
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		benchmarkRouteHash, err = have.Hash(3)
		if err != nil {
			b.Error(err)
		}
		if want != benchmarkRouteHash {
			b.Errorf("Want: %d; Have: %d", want, benchmarkRouteHash)
		}
	}
}

// BenchmarkRouteHash32-4   	50000000	        37.7 ns/op	       0 B/op	       0 allocs/op
func BenchmarkRouteHash32(b *testing.B) {
	have := path.NewRoute("general/single_store_mode/enabled")
	want := uint32(1644245266)

	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		benchmarkRouteHash = have.Hash32()
		if want != benchmarkRouteHash {
			b.Errorf("Want: %d; Have: %d", want, benchmarkRouteHash)
		}
	}
}

func TestRoutePartPosition(t *testing.T) {
	t.Parallel()
	tests := []struct {
		have     path.Route
		level    int
		wantPart string
		wantErr  error
	}{
		{path.NewRoute("general/single_\x80store_mode/enabled"), 0, "", path.ErrIncorrectPosition},
		{path.NewRoute("general/single_store_mode/enabled"), 0, "", path.ErrIncorrectPosition},
		{path.NewRoute("general/single_store_mode/enabled"), 1, "general", nil},
		{path.NewRoute("general/single_store_mode/enabled"), 2, "single_store_mode", nil},
		{path.NewRoute("general/single_store_mode/enabled"), 3, "enabled", nil},
		{path.NewRoute("system/full_page_cache/varnish/backend_port"), 4, "backend_port", nil},
		{path.NewRoute("general/single_store_mode/enabled"), -1, "", path.ErrIncorrectPosition},
		{path.NewRoute("general/single_store_mode/enabled"), 5, "", path.ErrIncorrectPosition},
		{path.NewRoute("general/single/store/website/group/mode/enabled/disabled/default"), 5, "group", nil},
	}
	for i, test := range tests {
		part, haveErr := test.have.Part(test.level)
		if test.wantErr != nil {
			assert.EqualError(t, haveErr, test.wantErr.Error(), "Index %d", i)
			assert.Nil(t, part.Chars, "Index %d", i)
			continue
		}
		assert.Exactly(t, test.wantPart, part.String(), "Index %d", i)
	}
}

var benchmarkRoutePart path.Route

// BenchmarkRoutePart-4	 5000000	       240 ns/op	       0 B/op	       0 allocs/op
func BenchmarkRoutePart(b *testing.B) {
	have := path.NewRoute("general/single_store_mode/enabled")
	want := "enabled"

	var err error
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		benchmarkRoutePart, err = have.Part(3)
		if err != nil {
			b.Error(err)
		}
		if benchmarkRoutePart.Chars == nil {
			b.Error("benchmarkRoutePart is nil! Unexpected")
		}
	}
	if want != benchmarkRoutePart.String() {
		b.Errorf("Want: %d; Have: %d", want, benchmarkRoutePart.String())
	}
}

func TestRouteValidate(t *testing.T) {
	t.Parallel()
	tests := []struct {
		have path.Route
		want error
	}{
		{path.NewRoute("//"), path.ErrIncorrectPath},
		{path.NewRoute("general/store_information/city"), nil},
		{path.NewRoute("general/store_information/city"), nil},
		{path.NewRoute("system/full_page_cache/varnish/backend_port"), nil},
		{path.NewRoute(""), path.ErrRouteEmpty},
		{path.NewRoute("general/store_information"), nil},
		////{path.NewRoute(path.MustNew("system/dev/debug").Bind(scope.WebsiteID, 22).String()), path.ErrIncorrectPath},
		{path.NewRoute("groups/33/general/store_information/street"), nil},
		{path.NewRoute("groups/33"), nil},
		{path.NewRoute("system/dEv/inv˚lid"), errors.New("This character \"˚\" is not allowed in Route system/dEv/inv˚lid")},
		{path.NewRoute("system/dEv/inv'lid"), errors.New("This character \"'\" is not allowed in Route system/dEv/inv'lid")},
		{path.NewRoute("syst3m/dEv/invalid"), nil},
		{path.Route{}, path.ErrRouteEmpty},
	}
	for i, test := range tests {
		haveErr := test.have.Validate()
		if test.want != nil {
			assert.EqualError(t, haveErr, test.want.Error(), "Index %d", i)
		} else {
			assert.NoError(t, haveErr, "Index %d", i)
		}
	}
}

var benchmarkRouteValidate error

// BenchmarkRouteValidate-4	20000000	        83.5 ns/op	       0 B/op	       0 allocs/op
func BenchmarkRouteValidate(b *testing.B) {
	have := path.NewRoute("system/dEv/d3bug")
	want := "system/dev/debug"

	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		benchmarkRouteValidate = have.Validate()
		if nil != benchmarkRouteValidate {
			b.Errorf("Want: %s; Have: %v", want, have)
		}
	}
}

func TestRouteSplit(t *testing.T) {
	t.Parallel()
	tests := []struct {
		have     path.Route
		wantPart []string
		wantErr  error
	}{
		{path.NewRoute("general/single_\x80store_mode"), []string{"general", "single_\x80store_mode"}, nil},
		{path.NewRoute("general/single_store_mode"), []string{"general", "single_store_mode"}, nil},
		{path.NewRoute("general"), nil, path.ErrIncorrectPath},
		{path.NewRoute("general/single_store_mode/enabled"), []string{"general", "single_store_mode", "enabled"}, nil},
		{path.NewRoute("system/full_page_cache/varnish/backend_port"), []string{"system", "full_page_cache", "varnish/backend_port"}, nil},
	}
	for i, test := range tests {
		sps, haveErr := test.have.Split()
		if test.wantErr != nil {
			assert.EqualError(t, haveErr, test.wantErr.Error(), "Index %d", i)
			assert.Nil(t, sps[0].Chars, "Index %d", i)
			assert.Nil(t, sps[1].Chars, "Index %d", i)
			continue
		}
		for i, wantPart := range test.wantPart {
			assert.Exactly(t, wantPart, sps[i].String(), "Index %d", i)
		}
	}
}

var benchmarkRouteSplit [path.Levels]path.Route

// BenchmarkRouteSplit-4    	 5000000	       286 ns/op	       0 B/op	       0 allocs/op
func BenchmarkRouteSplit(b *testing.B) {
	have := path.NewRoute("general/single_store_mode/enabled")
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		var err error
		benchmarkRouteSplit, err = have.Split()
		if err != nil {
			b.Error(err)
		}
		if benchmarkRouteSplit[1].Chars == nil {
			b.Error("benchmarkRouteSplit[1] is nil! Unexpected")
		}
	}
}
