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

package cfgpath_test

import (
	"database/sql"
	"database/sql/driver"
	"encoding"
	"encoding/json"
	"fmt"
	"hash/fnv"
	"testing"

	"github.com/corestoreio/csfw/config/cfgpath"
	"github.com/corestoreio/errors"
	"github.com/stretchr/testify/assert"

	"github.com/corestoreio/csfw/util/naughtystrings"
	"github.com/corestoreio/csfw/util/slices"
)

// These checks, if a type implements an interface, belong into the test package
// and not into its "main" package. Otherwise you would also compile each time
// all the packages with their interfaces, etc.
var _ encoding.TextMarshaler = (*cfgpath.Route)(nil)
var _ encoding.TextUnmarshaler = (*cfgpath.Route)(nil)
var _ sql.Scanner = (*cfgpath.Route)(nil)
var _ driver.Valuer = (*cfgpath.Route)(nil)
var _ fmt.GoStringer = (*cfgpath.Route)(nil)
var _ fmt.Stringer = (*cfgpath.Route)(nil)
var _ cfgpath.SelfRouter = (*cfgpath.Route)(nil)

func TestRouteRouteSelfer(t *testing.T) {

	r := cfgpath.NewRoute("a/b/c")
	assert.Exactly(t, r, r.SelfRoute())
}

func TestRouteGoString(t *testing.T) {

	tests := []struct {
		have cfgpath.Route
		want string
	}{
		{cfgpath.NewRoute("a"), "cfgpath.Route{Chars:[]byte(`a`)}"},
		{cfgpath.NewRoute(""), "cfgpath.Route{}"},
		{cfgpath.Route{}, "cfgpath.Route{}"},
	}
	for i, test := range tests {
		assert.Exactly(t, test.want, test.have.GoString(), "Index %d", i)
	}
}

func TestRouteEqual(t *testing.T) {

	tests := []struct {
		a    cfgpath.Route
		b    cfgpath.Route
		want bool
	}{
		{cfgpath.Route{}, cfgpath.Route{}, true},
		{cfgpath.NewRoute("a"), cfgpath.NewRoute("a"), true},
		{cfgpath.NewRoute("a"), cfgpath.NewRoute("b"), false},
		{cfgpath.NewRoute("a\x80"), cfgpath.NewRoute("a"), false},
		{cfgpath.NewRoute("general/single_\x80store_mode/enabled"), cfgpath.NewRoute("general/single_store_mode/enabled"), false},
		{cfgpath.NewRoute("general/single_store_mode/enabled"), cfgpath.NewRoute("general/single_store_mode/enabled"), true},
		{cfgpath.NewRoute(""), cfgpath.NewRoute(""), true},
		{cfgpath.NewRoute(""), cfgpath.NewRoute(), true},
	}
	for i, test := range tests {
		assert.Exactly(t, test.want, test.a.Equal(test.b), "Index %d", i)
	}
}

func TestRouteAppend(t *testing.T) {

	tests := []struct {
		a          cfgpath.Route
		b          cfgpath.Route
		want       string
		wantErrBhf errors.BehaviourFunc
	}{
		{cfgpath.NewRoute("aa/bb/cc"), cfgpath.NewRoute("dd"), "aa/bb/cc/dd", nil},
		{cfgpath.NewRoute("aa"), cfgpath.NewRoute("bb/cc"), "aa/bb/cc", nil},
		{cfgpath.NewRoute("aa"), cfgpath.NewRoute("bbcc"), "aa/bbcc", nil},
		{cfgpath.NewRoute("aa"), cfgpath.NewRoute("bb\x80cc"), "", errors.IsNotValid},
		{cfgpath.NewRoute("aa/"), cfgpath.NewRoute("bbcc"), "aa/bbcc", nil},
		{cfgpath.NewRoute("aa"), cfgpath.NewRoute("/bbcc"), "aa/bbcc", nil},
		{cfgpath.NewRoute("aa"), cfgpath.NewRoute("/bb\x00cc"), "aa/bb", nil},
		{cfgpath.NewRoute("ag"), cfgpath.NewRoute("b"), "ag/b", nil},
		{cfgpath.NewRoute("c/"), cfgpath.NewRoute("/b"), "c/b", nil},
		{cfgpath.NewRoute("d/"), cfgpath.NewRoute(""), "d", nil},
		{cfgpath.NewRoute("e"), cfgpath.NewRoute(""), "e", nil},
		{cfgpath.NewRoute(""), cfgpath.NewRoute("f"), "f", nil},
	}
	for i, test := range tests {
		haveErr := test.a.Append(test.b)
		if test.wantErrBhf != nil {
			assert.True(t, test.wantErrBhf(haveErr), "Index %d => %s", i, haveErr)
			continue
		}
		assert.NoError(t, haveErr, "Index %d", i)
		assert.Exactly(t, test.want, test.a.String(), "Index %d", i)
	}
}

func TestRouteVariadicAppend(t *testing.T) {

	tests := []struct {
		a          cfgpath.Route
		routes     []cfgpath.Route
		want       string
		wantErrBhf errors.BehaviourFunc
	}{
		{cfgpath.NewRoute("aa"), []cfgpath.Route{cfgpath.NewRoute("bb"), cfgpath.NewRoute("cc"), cfgpath.NewRoute("dd")}, "aa/bb/cc/dd", nil},
		{cfgpath.NewRoute("aa"), []cfgpath.Route{{}, {}, {Chars: []byte(`cb`)}}, "aa/cb", nil},
		{cfgpath.Route{}, []cfgpath.Route{{}, {}, {Chars: []byte(`cb`)}}, "cb", nil},
		{cfgpath.Route{}, []cfgpath.Route{{}, {}, {}}, "", errors.IsEmpty},
	}
	for i, test := range tests {
		haveErr := test.a.Append(test.routes...)
		if test.wantErrBhf != nil {
			assert.True(t, test.wantErrBhf(haveErr), "Index %d => %s", i, haveErr)
			continue
		}
		assert.NoError(t, haveErr, "Index %d", i)
		assert.Exactly(t, test.want, test.a.String(), "Index %d", i)
	}

}

var benchmarkRouteAppendWant = cfgpath.NewRoute("general/single_store_mode/enabled")

// BenchmarkRouteAppend-4   	 5000000	       376 ns/op	      64 B/op	       2 allocs/op
func BenchmarkRouteAppend(b *testing.B) {
	havePartials := []cfgpath.Route{cfgpath.NewRoute("single_store_mode"), cfgpath.NewRoute("enabled")}
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		var have = cfgpath.NewRoute("general/")
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

	r := cfgpath.NewRoute("admin/security/password_lifetime")
	j, err := json.Marshal(r)
	assert.NoError(t, err)
	assert.Exactly(t, "\"admin/security/password_lifetime\"", string(j))
}

func TestRouteUnmarshalTextOk(t *testing.T) {

	var r cfgpath.Route
	err := json.Unmarshal([]byte(`"admin/security/password_lifetime"`), &r)
	assert.NoError(t, err)
	assert.Exactly(t, "admin/security/password_lifetime", r.String())
}

func TestRouteLevel(t *testing.T) {

	tests := []struct {
		have  cfgpath.Route
		depth int
		want  string
	}{
		{cfgpath.NewRoute("general/single_store_mode/enabled"), 0, ""},
		{cfgpath.NewRoute("general/single_store_mode/enabled"), 1, "general"},
		{cfgpath.NewRoute("general/single_store_mode/enabled"), 2, "general/single_store_mode"},
		{cfgpath.NewRoute("general/single_store_mode/enabled"), 3, "general/single_store_mode/enabled"},
		{cfgpath.NewRoute("general/single_store_mode/enabled"), -1, "general/single_store_mode/enabled"},
		{cfgpath.NewRoute("general/single_store_mode/enabled"), 5, "general/single_store_mode/enabled"},
		{cfgpath.NewRoute("general/single_store_mode/enabled"), 4, "general/single_store_mode/enabled"},
		{cfgpath.NewRoute("system/full_page_cache/varnish/backend_port"), 3, "system/full_page_cache/varnish"},
	}
	for i, test := range tests {
		r, err := test.have.Level(test.depth)
		assert.NoError(t, err)
		assert.Exactly(t, test.want, r.String(), "Index %d", i)
	}
}

var benchmarkRouteLevel cfgpath.Route

// BenchmarkRouteLevel_One-4	 5000000	       297 ns/op	      16 B/op	       1 allocs/op
func BenchmarkRouteLevel_One(b *testing.B) {
	benchmarkRouteLevelRun(b, 1, cfgpath.NewRoute("system/dev/debug"), cfgpath.NewRoute("system"))
}

// BenchmarkRouteLevel_Two-4	 5000000	       332 ns/op	      16 B/op	       1 allocs/op
func BenchmarkRouteLevel_Two(b *testing.B) {
	benchmarkRouteLevelRun(b, 2, cfgpath.NewRoute("system/dev/debug"), cfgpath.NewRoute("system/dev"))
}

// BenchmarkRouteLevel_All-4	 5000000	       379 ns/op	      16 B/op	       1 allocs/op
func BenchmarkRouteLevel_All(b *testing.B) {
	benchmarkRouteLevelRun(b, -1, cfgpath.NewRoute("system/dev/debug"), cfgpath.NewRoute("system/dev/debug"))
}

func benchmarkRouteLevelRun(b *testing.B, level int, have, want cfgpath.Route) {
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

	tests := []struct {
		have       cfgpath.Route
		depth      int
		wantHash   uint32
		wantErrBHf errors.BehaviourFunc
		wantLevel  string
	}{
		{cfgpath.NewRoute("general/single_\x80store_mode/enabled"), 0, 0, errors.IsNotValid, ""},
		{cfgpath.NewRoute("general/single_store_mode/enabled"), 0, 2166136261, nil, ""},
		{cfgpath.NewRoute("general/single_store_mode/enabled"), 1, 616112491, nil, "general"},
		{cfgpath.NewRoute("general/single_store_mode/enabled"), 2, 2274889228, nil, "general/single_store_mode"},
		{cfgpath.NewRoute("general/single_store_mode/enabled"), 3, 1644245266, nil, "general/single_store_mode/enabled"},
		{cfgpath.NewRoute("general/single_store_mode/enabled"), -1, 1644245266, nil, "general/single_store_mode/enabled"},
		{cfgpath.NewRoute("general/single_store_mode/enabled"), 5, 1644245266, nil, "general/single_store_mode/enabled"},
		{cfgpath.NewRoute("general/single_store_mode/enabled"), 4, 1644245266, nil, "general/single_store_mode/enabled"},
	}
	for i, test := range tests {

		hv, err := test.have.Hash(test.depth)
		if test.wantErrBHf != nil {
			assert.True(t, test.wantErrBHf(err), "Index %d => %s", i, err)
			assert.Empty(t, hv, "Index %d", i)
			continue
		}
		assert.NoError(t, err, "Index %d", i)

		check := fnv.New32a()
		_, cErr := check.Write([]byte(test.wantLevel))
		assert.NoError(t, cErr)
		assert.Exactly(t, check.Sum32(), hv, "Have %d Want %d Index %d", check.Sum32(), hv, i)

		l, err := test.have.Level(test.depth)
		assert.Exactly(t, test.wantLevel, l.String(), "Index %d", i)
		assert.Exactly(t, test.wantHash, hv, "Have %d Want %d Index %d", test.wantHash, hv, i)
	}
}

func TestRouteHash32(t *testing.T) {

	tests := []struct {
		have     cfgpath.Route
		wantHash uint32
	}{
		{cfgpath.NewRoute("general/single_\x80store_mode/enabled"), 1310908924},
		{cfgpath.NewRoute("general/single_store_mode/enabled"), 1644245266},
		{cfgpath.NewRoute(""), 2166136261},
		{cfgpath.Route{}, 2166136261},
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
	have := cfgpath.NewRoute("general/single_store_mode/enabled")
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
	have := cfgpath.NewRoute("general/single_store_mode/enabled")
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

	tests := []struct {
		have     cfgpath.Route
		pos      int
		wantPart string
		wantErr  bool
	}{
		{cfgpath.NewRoute("general/single_\x80store_mode/enabled"), 0, "", true},
		{cfgpath.NewRoute("general/single_store_mode/enabled"), 0, "", true},
		{cfgpath.NewRoute("general/single_store_mode/enabled"), 1, "general", false},
		{cfgpath.NewRoute("general/single_store_mode/enabled"), 2, "single_store_mode", false},
		{cfgpath.NewRoute("general/single_store_mode/enabled"), 3, "enabled", false},
		{cfgpath.NewRoute("system/full_page_cache/varnish/backend_port"), 4, "backend_port", false},
		{cfgpath.NewRoute("general/single_store_mode/enabled"), -1, "", true},
		{cfgpath.NewRoute("general/single_store_mode/enabled"), 5, "", true},
		{cfgpath.NewRoute("general/single/store/website/group/mode/enabled/disabled/default"), 5, "group", false},
	}
	for i, test := range tests {
		part, haveErr := test.have.Part(test.pos)
		if test.wantErr {
			assert.Nil(t, part.Chars, "Index %d", i)
			assert.True(t, errors.IsNotValid(haveErr), "Index %d => %s", i, haveErr)
			continue
		}
		assert.Exactly(t, test.wantPart, part.String(), "Index %d", i)
	}
}

var benchmarkRoutePart cfgpath.Route

// BenchmarkRoutePart-4	 5000000	       240 ns/op	       0 B/op	       0 allocs/op
func BenchmarkRoutePart(b *testing.B) {
	have := cfgpath.NewRoute("general/single_store_mode/enabled")
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

	tests := []struct {
		have    cfgpath.Route
		wantBhf errors.BehaviourFunc
	}{
		{cfgpath.NewRoute("//"), errors.IsNotValid},
		{cfgpath.NewRoute("general/store_information/city"), nil},
		{cfgpath.NewRoute("general/store_information/city"), nil},
		{cfgpath.NewRoute("system/full_page_cache/varnish/backend_port"), nil},
		{cfgpath.NewRoute(""), errors.IsEmpty},
		{cfgpath.NewRoute("general/store_information"), nil},
		////{cfgpath.NewRoute(cfgpath.MustNew("system/dev/debug").Bind(scope.WebsiteID, 22).String()), cfgpath.ErrIncorrectPath},
		{cfgpath.NewRoute("groups/33/general/store_information/street"), nil},
		{cfgpath.NewRoute("groups/33"), nil},
		{cfgpath.NewRoute("system/dEv/inv˚lid"), errors.IsNotValid},
		{cfgpath.NewRoute("system/dEv/inv'lid"), errors.IsNotValid},
		{cfgpath.NewRoute("syst3m/dEv/invalid"), nil},
		{cfgpath.Route{}, errors.IsEmpty},
	}
	for i, test := range tests {
		haveErr := test.have.Validate()
		if test.wantBhf != nil {
			assert.True(t, test.wantBhf(haveErr), "Index %d => %s", i, haveErr)
		} else {
			assert.NoError(t, haveErr, "Index %d", i)
		}
	}
}

func TestRouteValidateNaughtyStrings(t *testing.T) {
	var valids slices.String = []string{"undefined", "undef", "null", "NULL", "nil", "NIL", "true", "false", "True", "False", "None", "hasOwnProperty", "0", "1", "1/2", "1E2", "1E02", "1/0", "0/0", "999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999", "NaN", "Infinity", "INF", "0x0", "0xffffffff", "0xffffffffffffffff", "0xabad1dea", "123456789012345678901234567890123456789", "01000", "08", "09", "_", "CON", "PRN", "AUX", "NUL", "COM1", "LPT1", "LPT2", "LPT3", "COM2", "COM3", "COM4", "evaluate", "mocha", "expression", "classic", "basement"}
	for _, str := range naughtystrings.Unencoded() {
		r := cfgpath.NewRoute(str)
		if err := r.Validate(); valids.Contains(str) && err != nil {
			t.Errorf("Should be valid %q but error: %+v", str, err)
		}
	}
}

var benchmarkRouteValidate error

// BenchmarkRouteValidate-4	20000000	        83.5 ns/op	       0 B/op	       0 allocs/op
func BenchmarkRouteValidate(b *testing.B) {
	have := cfgpath.NewRoute("system/dEv/d3bug")
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

	tests := []struct {
		have       cfgpath.Route
		wantPart   []string
		wantErrBhf errors.BehaviourFunc
	}{
		{cfgpath.NewRoute("general/single_\x80store_mode"), []string{"general", "single_\x80store_mode"}, nil},
		{cfgpath.NewRoute("general/single_store_mode"), []string{"general", "single_store_mode"}, nil},
		{cfgpath.NewRoute("general"), nil, errors.IsNotValid},
		{cfgpath.NewRoute("general/single_store_mode/enabled"), []string{"general", "single_store_mode", "enabled"}, nil},
		{cfgpath.NewRoute("system/full_page_cache/varnish/backend_port"), []string{"system", "full_page_cache", "varnish/backend_port"}, nil},
	}
	for i, test := range tests {
		sps, haveErr := test.have.Split()
		if test.wantErrBhf != nil {
			assert.True(t, test.wantErrBhf(haveErr), "Index %d => %s", i, haveErr)
			assert.Nil(t, sps[0].Chars, "Index %d", i)
			assert.Nil(t, sps[1].Chars, "Index %d", i)
			continue
		}
		for i, wantPart := range test.wantPart {
			assert.Exactly(t, wantPart, sps[i].String(), "Index %d", i)
		}
	}
}

var benchmarkRouteSplit [cfgpath.Levels]cfgpath.Route

// BenchmarkRouteSplit-4    	 5000000	       286 ns/op	       0 B/op	       0 allocs/op
func BenchmarkRouteSplit(b *testing.B) {
	have := cfgpath.NewRoute("general/single_store_mode/enabled")
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
