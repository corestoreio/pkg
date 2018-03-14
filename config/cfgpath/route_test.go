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
	"encoding"
	"encoding/json"
	"fmt"
	"hash/fnv"
	"testing"

	"github.com/corestoreio/errors"
	"github.com/corestoreio/pkg/config/cfgpath"
	"github.com/corestoreio/pkg/util/naughtystrings"
	"github.com/corestoreio/pkg/util/slices"
	"github.com/stretchr/testify/assert"
)

// These checks, if a type implements an interface, belong into the test package
// and not into its "main" package. Otherwise you would also compile each time
// all the packages with their interfaces, etc.
var _ encoding.TextMarshaler = (*cfgpath.Route)(nil)
var _ encoding.TextUnmarshaler = (*cfgpath.Route)(nil)

var _ fmt.Stringer = (*cfgpath.Route)(nil)
var _ cfgpath.SelfRouter = (*cfgpath.Route)(nil)

func TestRouteRouteSelfer(t *testing.T) {

	r := cfgpath.MakeRoute("a/b/c")
	assert.Exactly(t, r, r.SelfRoute())
}

func TestRouteEqual(t *testing.T) {

	tests := []struct {
		a    cfgpath.Route
		b    cfgpath.Route
		want bool
	}{
		{cfgpath.Route{}, cfgpath.Route{}, true},
		{cfgpath.MakeRoute("a"), cfgpath.MakeRoute("a"), true},
		{cfgpath.MakeRoute("a"), cfgpath.MakeRoute("b"), false},
		{cfgpath.MakeRoute("a\x80"), cfgpath.MakeRoute("a"), false},
		{cfgpath.MakeRoute("general/single_\x80store_mode/enabled"), cfgpath.MakeRoute("general/single_store_mode/enabled"), false},
		{cfgpath.MakeRoute("general/single_store_mode/enabled"), cfgpath.MakeRoute("general/single_store_mode/enabled"), true},
		{cfgpath.MakeRoute(""), cfgpath.MakeRoute(""), true},
	}
	for i, test := range tests {
		assert.Exactly(t, test.want, test.a.Equal(test.b), "Index %d", i)
	}
}

func TestRouteAppend(t *testing.T) {

	tests := []struct {
		a           cfgpath.Route
		b           cfgpath.Route
		want        string
		wantErrKind errors.Kind
	}{
		0:  {cfgpath.MakeRoute("aa/bb/cc"), cfgpath.MakeRoute("dd"), "aa/bb/cc/dd", errors.NoKind},
		1:  {cfgpath.MakeRoute("aa"), cfgpath.MakeRoute("bb/cc"), "aa/bb/cc", errors.NoKind},
		2:  {cfgpath.MakeRoute("aa"), cfgpath.MakeRoute("bbcc"), "aa/bbcc", errors.NoKind},
		3:  {cfgpath.MakeRoute("aa"), cfgpath.MakeRoute("bb\x80cc"), "", errors.NotValid},
		4:  {cfgpath.MakeRoute("aa/"), cfgpath.MakeRoute("bbcc"), "aa/bbcc", errors.NoKind},
		5:  {cfgpath.MakeRoute("aa"), cfgpath.MakeRoute("/bbcc"), "aa/bbcc", errors.NoKind},
		6:  {cfgpath.MakeRoute("aa"), cfgpath.MakeRoute("/bb\x00cc"), "", errors.NotValid},
		7:  {cfgpath.MakeRoute("ag"), cfgpath.MakeRoute("b"), "ag/b", errors.NoKind},
		8:  {cfgpath.MakeRoute("c/"), cfgpath.MakeRoute("/b"), "c/b", errors.NoKind},
		9:  {cfgpath.MakeRoute("d/"), cfgpath.MakeRoute(""), "d", errors.NoKind},
		10: {cfgpath.MakeRoute("e"), cfgpath.MakeRoute(""), "e", errors.NoKind},
		11: {cfgpath.MakeRoute(""), cfgpath.MakeRoute("f"), "f", errors.NoKind},
	}
	for i, test := range tests {
		test.a = test.a.Append(test.b)
		haveErr := test.a.Validate()
		if test.wantErrKind > 0 {
			assert.True(t, test.wantErrKind.Match(haveErr), "Index %d => %s", i, haveErr)
			continue
		}
		assert.NoError(t, haveErr, "Index %d", i)
		assert.Exactly(t, test.want, test.a.String(), "Index %d", i)
	}
}

func TestRouteVariadicAppend(t *testing.T) {

	tests := []struct {
		a           cfgpath.Route
		routes      []cfgpath.Route
		want        string
		wantErrKind errors.Kind
	}{
		//{cfgpath.MakeRoute("aa"), []cfgpath.Route{cfgpath.MakeRoute("bb"), cfgpath.MakeRoute("cc"), cfgpath.MakeRoute("dd")}, "aa/bb/cc/dd", errors.NoKind},
		{cfgpath.MakeRoute("aa"), []cfgpath.Route{{}, {}, {Data: `cb`, Valid: true}}, "aa/cb", errors.NoKind},
		//{cfgpath.Route{}, []cfgpath.Route{{}, {}, {Data: `cb`, Valid: true}}, "cb", errors.NoKind},
		//{cfgpath.Route{}, []cfgpath.Route{{}, {}, {}}, "", errors.Empty},
	}
	for i, test := range tests {
		test.a = test.a.Append(test.routes...)
		haveErr := test.a.Validate()
		if test.wantErrKind > 0 {
			assert.True(t, test.wantErrKind.Match(haveErr), "Index %d => %s", i, haveErr)
			continue
		}
		assert.NoError(t, haveErr, "Index %d", i)
		assert.Exactly(t, test.want, test.a.String(), "Index %d", i)
	}

}

var benchmarkRouteAppendWant = cfgpath.MakeRoute("general/single_store_mode/enabled")

// BenchmarkRouteAppend-4   	 5000000	       376 ns/op	      64 B/op	       2 allocs/op
func BenchmarkRouteAppend(b *testing.B) {
	havePartials := []cfgpath.Route{cfgpath.MakeRoute("single_store_mode"), cfgpath.MakeRoute("enabled")}
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		have := cfgpath.MakeRoute("general/").Append(havePartials...)
		if !benchmarkRouteAppendWant.Equal(have) {
			b.Errorf("Want: %s; Have: %s", benchmarkRouteAppendWant, have)
		}
	}
}

func TestRoute_Append(t *testing.T) {
	havePartials := []cfgpath.Route{cfgpath.MakeRoute("single_store_mode"), cfgpath.MakeRoute("enabled")}
	have := cfgpath.MakeRoute("general/").Append(havePartials...)
	assert.Exactly(t, `general/single_store_mode/enabled`, have.String())
}

func TestRouteTextMarshal(t *testing.T) {

	r := cfgpath.MakeRoute("admin/security/password_lifetime")
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
		{cfgpath.MakeRoute("general/single_store_mode/enabled"), 0, "<nil>"},
		{cfgpath.MakeRoute("general/single_store_mode/enabled"), 1, "general"},
		{cfgpath.MakeRoute("general/single_store_mode/enabled"), 2, "general/single_store_mode"},
		{cfgpath.MakeRoute("general/single_store_mode/enabled"), 3, "general/single_store_mode/enabled"},
		{cfgpath.MakeRoute("general/single_store_mode/enabled"), -1, "general/single_store_mode/enabled"},
		{cfgpath.MakeRoute("general/single_store_mode/enabled"), 5, "general/single_store_mode/enabled"},
		{cfgpath.MakeRoute("general/single_store_mode/enabled"), 4, "general/single_store_mode/enabled"},
		{cfgpath.MakeRoute("system/full_page_cache/varnish/backend_port"), 3, "system/full_page_cache/varnish"},
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
	benchmarkRouteLevelRun(b, 1, cfgpath.MakeRoute("system/dev/debug"), cfgpath.MakeRoute("system"))
}

// BenchmarkRouteLevel_Two-4	 5000000	       332 ns/op	      16 B/op	       1 allocs/op
func BenchmarkRouteLevel_Two(b *testing.B) {
	benchmarkRouteLevelRun(b, 2, cfgpath.MakeRoute("system/dev/debug"), cfgpath.MakeRoute("system/dev"))
}

// BenchmarkRouteLevel_All-4	 5000000	       379 ns/op	      16 B/op	       1 allocs/op
func BenchmarkRouteLevel_All(b *testing.B) {
	benchmarkRouteLevelRun(b, -1, cfgpath.MakeRoute("system/dev/debug"), cfgpath.MakeRoute("system/dev/debug"))
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
	if !benchmarkRouteLevel.Equal(want) {
		b.Errorf("Want: %s; Have, %s", want, benchmarkRouteLevel)
	}
}

func TestRouteHash(t *testing.T) {
	t.Parallel()
	tests := []struct {
		have        cfgpath.Route
		depth       int
		wantHash    uint32
		wantErrKind errors.Kind
		wantLevel   string
	}{
		{cfgpath.MakeRoute("general/single_\x80store_mode/enabled"), 0, 0, errors.NotValid, ""},
		{cfgpath.MakeRoute("general/single_store_mode/enabled"), 0, 2166136261, errors.NoKind, "<nil>"},
		{cfgpath.MakeRoute("general/single_store_mode/enabled"), 1, 616112491, errors.NoKind, "general"},
		{cfgpath.MakeRoute("general/single_store_mode/enabled"), 2, 2274889228, errors.NoKind, "general/single_store_mode"},
		{cfgpath.MakeRoute("general/single_store_mode/enabled"), 3, 1644245266, errors.NoKind, "general/single_store_mode/enabled"},
		{cfgpath.MakeRoute("general/single_store_mode/enabled"), -1, 1644245266, errors.NoKind, "general/single_store_mode/enabled"},
		{cfgpath.MakeRoute("general/single_store_mode/enabled"), 5, 1644245266, errors.NoKind, "general/single_store_mode/enabled"},
		{cfgpath.MakeRoute("general/single_store_mode/enabled"), 4, 1644245266, errors.NoKind, "general/single_store_mode/enabled"},
	}
	for i, test := range tests {

		hv, err := test.have.Hash(test.depth)
		if test.wantErrKind > 0 {
			assert.True(t, test.wantErrKind.Match(err), "Index %d => %s", i, err)
			assert.Empty(t, hv, "Index %d", i)
			continue
		}
		assert.NoError(t, err, "Index %d", i)

		check := fnv.New32a()
		if test.wantLevel == "<nil>" {
			_, cErr := check.Write(nil)
			assert.NoError(t, cErr)
		} else {
			_, cErr := check.Write([]byte(test.wantLevel))
			assert.NoError(t, cErr)
		}

		assert.Exactly(t, check.Sum32(), hv, "Have %d Want %d Index %d", check.Sum32(), hv, i)

		l, err := test.have.Level(test.depth)
		assert.NoError(t, err)
		assert.Exactly(t, test.wantLevel, l.String(), "Index %d", i)
		assert.Exactly(t, test.wantHash, hv, "Expected %d Actual %d Index %d", test.wantHash, hv, i)
	}
}

var benchmarkRouteHash uint32

// BenchmarkRouteHash-4     	 5000000	       287 ns/op	       0 B/op	       0 allocs/op
func BenchmarkRouteHash(b *testing.B) {
	have := cfgpath.MakeRoute("general/single_store_mode/enabled")
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
	have := cfgpath.MakeRoute("general/single_store_mode/enabled")
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
		{cfgpath.MakeRoute("general/single_\x80store_mode/enabled"), 0, "", true},
		{cfgpath.MakeRoute("general/single_store_mode/enabled"), 0, "", true},
		{cfgpath.MakeRoute("general/single_store_mode/enabled"), 1, "general", false},
		{cfgpath.MakeRoute("general/single_store_mode/enabled"), 2, "single_store_mode", false},
		{cfgpath.MakeRoute("general/single_store_mode/enabled"), 3, "enabled", false},
		{cfgpath.MakeRoute("system/full_page_cache/varnish/backend_port"), 4, "backend_port", false},
		{cfgpath.MakeRoute("general/single_store_mode/enabled"), -1, "", true},
		{cfgpath.MakeRoute("general/single_store_mode/enabled"), 5, "", true},
		{cfgpath.MakeRoute("general/single/store/website/group/mode/enabled/disabled/default"), 5, "group", false},
	}
	for i, test := range tests {
		part, haveErr := test.have.Part(test.pos)
		if test.wantErr {
			assert.Empty(t, part.Data, "Index %d", i)
			assert.True(t, errors.NotValid.Match(haveErr), "Index %d => %s", i, haveErr)
			continue
		}
		assert.Exactly(t, test.wantPart, part.String(), "Index %d", i)
	}
}

var benchmarkRoutePart cfgpath.Route

// BenchmarkRoutePart-4	 5000000	       240 ns/op	       0 B/op	       0 allocs/op
func BenchmarkRoutePart(b *testing.B) {
	have := cfgpath.MakeRoute("general/single_store_mode/enabled")
	want := "enabled"

	var err error
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		benchmarkRoutePart, err = have.Part(3)
		if err != nil {
			b.Error(err)
		}
		if benchmarkRoutePart.Data == "" {
			b.Error("benchmarkRoutePart is nil! Unexpected")
		}
	}
	if want != benchmarkRoutePart.String() {
		b.Errorf("Want: %q; Have: %q", want, benchmarkRoutePart.String())
	}
}

func TestRouteValidate(t *testing.T) {

	tests := []struct {
		have        cfgpath.Route
		wantErrKind errors.Kind
	}{
		{cfgpath.MakeRoute("//"), errors.NotValid},
		{cfgpath.MakeRoute("general/store_information/city"), errors.NoKind},
		{cfgpath.MakeRoute("general/store_information/city"), errors.NoKind},
		{cfgpath.MakeRoute("system/full_page_cache/varnish/backend_port"), errors.NoKind},
		{cfgpath.MakeRoute(""), errors.Empty},
		{cfgpath.MakeRoute("general/store_information"), errors.NoKind},
		////{cfgpath.MakeRoute(cfgpath.MustNew("system/dev/debug").Bind(scope.WebsiteID, 22).String()), cfgpath.ErrIncorrectPath},
		{cfgpath.MakeRoute("groups/33/general/store_information/street"), errors.NoKind},
		{cfgpath.MakeRoute("groups/33"), errors.NoKind},
		{cfgpath.MakeRoute("system/dEv/inv˚lid"), errors.NotValid},
		{cfgpath.MakeRoute("system/dEv/inv'lid"), errors.NotValid},
		{cfgpath.MakeRoute("syst3m/dEv/invalid"), errors.NoKind},
		{cfgpath.Route{}, errors.Empty},
	}
	for i, test := range tests {
		haveErr := test.have.Validate()
		if test.wantErrKind > 0 {
			assert.True(t, test.wantErrKind.Match(haveErr), "Index %d => %s", i, haveErr)
		} else {
			assert.NoError(t, haveErr, "Index %d", i)
		}
	}
}

func TestRouteValidateNaughtyStrings(t *testing.T) {
	var valids slices.String = []string{"undefined", "undef", "null", "NULL", "nil", "NIL", "true", "false", "True", "False", "None", "hasOwnProperty", "0", "1", "1/2", "1E2", "1E02", "1/0", "0/0", "999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999", "NaN", "Infinity", "INF", "0x0", "0xffffffff", "0xffffffffffffffff", "0xabad1dea", "123456789012345678901234567890123456789", "01000", "08", "09", "_", "CON", "PRN", "AUX", "NUL", "COM1", "LPT1", "LPT2", "LPT3", "COM2", "COM3", "COM4", "evaluate", "mocha", "expression", "classic", "basement"}
	for _, str := range naughtystrings.Unencoded() {
		r := cfgpath.MakeRoute(str)
		if err := r.Validate(); valids.Contains(str) && err != nil {
			t.Errorf("Should be valid %q but error: %+v", str, err)
		}
	}
}

var benchmarkRouteValidate error

// BenchmarkRouteValidate-4	20000000	        83.5 ns/op	       0 B/op	       0 allocs/op
func BenchmarkRouteValidate(b *testing.B) {
	have := cfgpath.MakeRoute("system/dEv/d3bug")
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
		have        cfgpath.Route
		wantPart    []string
		wantErrKind errors.Kind
	}{
		{cfgpath.MakeRoute("general/single_\x80store_mode"), []string{"general", "single_\x80store_mode"}, errors.NoKind},
		{cfgpath.MakeRoute("general/single_store_mode"), []string{"general", "single_store_mode"}, errors.NoKind},
		{cfgpath.MakeRoute("general"), nil, errors.NotValid},
		{cfgpath.MakeRoute("general/single_store_mode/enabled"), []string{"general", "single_store_mode", "enabled"}, errors.NoKind},
		{cfgpath.MakeRoute("system/full_page_cache/varnish/backend_port"), []string{"system", "full_page_cache", "varnish/backend_port"}, errors.NoKind},
	}
	for i, test := range tests {
		sps, haveErr := test.have.Split()
		if test.wantErrKind > 0 {
			assert.True(t, test.wantErrKind.Match(haveErr), "Index %d => %s", i, haveErr)
			assert.Empty(t, sps[0].Data, "Index %d", i)
			assert.Empty(t, sps[1].Data, "Index %d", i)
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
	have := cfgpath.MakeRoute("general/single_store_mode/enabled")
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		var err error
		benchmarkRouteSplit, err = have.Split()
		if err != nil {
			b.Error(err)
		}
		if benchmarkRouteSplit[1].Data == "" {
			b.Error("benchmarkRouteSplit[1] is nil! Unexpected")
		}
	}
}
