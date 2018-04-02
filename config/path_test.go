// Copyright 2015-present, Cyrill @ Schumacher.fm and the CoreStore contributors
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
	"encoding"
	"fmt"
	"hash/fnv"
	"sort"
	"testing"

	"github.com/corestoreio/errors"
	"github.com/corestoreio/pkg/store/scope"
	"github.com/corestoreio/pkg/util/naughtystrings"
	"github.com/stretchr/testify/assert"
)

var (
	_ encoding.TextMarshaler   = (*Path)(nil)
	_ encoding.TextUnmarshaler = (*Path)(nil)
	_ fmt.Stringer             = (*Path)(nil)
)

func TestNewByParts(t *testing.T) {
	t.Parallel()
	tests := []struct {
		path        string
		want        string
		wantErrKind errors.Kind
	}{
		{"aa/bb/cc", "aa/bb/cc", errors.NoKind},
		{"aa/bb/c", "aa/bb/cc", errors.NotValid},
		{"", "", errors.Empty},
	}
	for i, test := range tests {
		haveP, haveErr := MakePath(test.path)
		if test.wantErrKind > 0 {
			assert.Empty(t, haveP.route, "Index %d", i)
			assert.True(t, test.wantErrKind.Match(haveErr), "Index %d => %s", i, haveErr)
			continue
		}
		l, err := haveP.Level(-1)
		assert.NoError(t, err, "Index %d", i)
		assert.Exactly(t, test.want, l, "Index %d", i)
	}
}

func TestMustNewByPartsPanic(t *testing.T) {
	t.Parallel()
	defer func() {
		if r := recover(); r != nil {
			err := r.(error)
			assert.True(t, errors.NotValid.Match(err), "Error => %s", err)
		} else {
			t.Fatal("Expecting a panic")
		}
	}()
	_ = MustMakePath("a/\x80/c")
}

func TestMustNewByPartsNoPanic(t *testing.T) {
	t.Parallel()
	defer func() {
		if r := recover(); r != nil {
			assert.NotNil(t, r, "Did not expect a panic")
		} else {
			assert.Nil(t, r, "Why is here a panic")
		}
	}()
	p := MustMakePath("aa/bb/cc")
	assert.Exactly(t, "default/0/aa/bb/cc", p.String())
}

func TestMakePath(t *testing.T) {
	t.Parallel()
	tests := []struct {
		route       string
		s           scope.Type
		id          int64
		wantFQ      string
		wantErrKind errors.Kind
	}{
		{"ab/b\x80/cd", scope.Website, 3, "websites/3/ab/ba/cd", errors.NotValid},
		{"ab/ba/cd", scope.Website, 3, "websites/3/ab/ba/cd", errors.NoKind},
		{"ad/ba/ca/sd", scope.Website, 3, "websites/3/ad/ba/ca/sd", errors.NoKind},
		{"as/sb", scope.Website, 3, "websites/3/a/b/c/d", errors.NotValid},
		{"aa/bb/cc", scope.Group, 3, "default/0/aa/bb/cc", errors.NoKind},
		{"aa/bb/cc", scope.Store, 3, "stores/3/aa/bb/cc", errors.NoKind},
	}
	for i, test := range tests {
		haveP, haveErr := MakePath(test.route)
		haveP = haveP.Bind(test.s.Pack(test.id))
		if test.wantErrKind > 0 {
			assert.True(t, test.wantErrKind.Match(haveErr), "Index %d", i)
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
		scp         scope.Type
		id          int64
		route       string
		want        string
		wantErrKind errors.Kind
	}{
		{scope.Default, 0, "", "", errors.Empty},
		{scope.Default, 0, "", "", errors.Empty},
		{scope.Default, 0, "system/dev/debug", scope.StrDefault.String() + "/0/system/dev/debug", errors.NoKind},
		{scope.Default, 44, "system/dev/debug", scope.StrDefault.String() + "/0/system/dev/debug", errors.NoKind},
		{scope.Website, 0, "system/dev/debug", scope.StrWebsites.String() + "/0/system/dev/debug", errors.NoKind},
		{scope.Website, 343, "system/dev/debug", scope.StrWebsites.String() + "/343/system/dev/debug", errors.NoKind},
		{250, 0, "system/dev/debug", scope.StrDefault.String() + "/0/system/dev/debug", errors.NoKind},
		{250, 343, "system/dev/debug", scope.StrDefault.String() + "/0/system/dev/debug", errors.NoKind},
	}
	for i, test := range tests {
		p, pErr := MakePath(test.route)
		p = p.Bind(test.scp.Pack(test.id))
		have, haveErr := p.FQ()
		if test.wantErrKind > 0 {
			assert.Empty(t, have, "Index %d", i)
			if pErr != nil {
				assert.True(t, test.wantErrKind.Match(pErr), "Index %d => %s", i, pErr)
				continue
			}
			assert.True(t, test.wantErrKind.Match(haveErr), "Index %d => %s", i, haveErr)
			continue
		}
		assert.NoError(t, haveErr, "Index %d", i)
		assert.Exactly(t, test.want, have, "Index %d", i)
	}

	r := "catalog/frontend/list_allow_all"
	assert.Exactly(t, "stores/7475/catalog/frontend/list_allow_all", MustMakePath(r).BindStore(7475).String())
	p := MustMakePath(r).BindStore(5)
	assert.Exactly(t, "stores/5/catalog/frontend/list_allow_all", p.String())
	assert.Exactly(t, "Path{ Route:MakeRoute(\"catalog/frontend/list_allow_all\"), ScopeHash: 67108869 }", p.GoString())
}

func TestShouldNotPanicBecauseOfIncorrectStrScope(t *testing.T) {
	t.Parallel()
	assert.Exactly(t, "stores/345/xxxxx/yyyyy/zzzzz", MustMakePath("xxxxx/yyyyy/zzzzz").BindStore(345).String())
	defer func() {
		if r := recover(); r != nil {
			t.Fatal("Did not expect a panic")
		}
	}()
	_ = MustMakePath("xxxxx/yyyyy/zzzzz").Bind(345)
}

func TestShouldPanicIncorrectPath(t *testing.T) {
	t.Parallel()
	assert.Exactly(t, "default/0/xxxxx/yyyyy/zzzzz", MustMakePath("xxxxx/yyyyy/zzzzz").Bind(345).String())
	defer func() {
		if r := recover(); r != nil {
			err := r.(error)
			assert.True(t, errors.NotValid.Match(err))
		} else {
			t.Fatal("Expecting a panic")
		}
	}()
	assert.Exactly(t, "websites/345/xxxxx/yyyyy", MustMakePath("xxxxx/yyyyy").BindWebsite(345).String())
}

func TestSplitFQ(t *testing.T) {

	tests := []struct {
		have        string
		wantScope   string
		wantScopeID int64
		wantPath    string
		wantErrKind errors.Kind
	}{
		{"groups/1/catalog/frontend/list_allow_all", "default", 0, "<nil>", errors.NotSupported},
		{"stores/7475/catalog/frontend/list_allow_all", scope.StrStores.String(), 7475, "catalog/frontend/list_allow_all", errors.NoKind},
		{"stores/4/system/full_page_cache/varnish/backend_port", scope.StrStores.String(), 4, "system/full_page_cache/varnish/backend_port", errors.NoKind},
		{"websites/1/catalog/frontend/list_allow_all", scope.StrWebsites.String(), 1, "catalog/frontend/list_allow_all", errors.NoKind},
		{"default/0/catalog/frontend/list_allow_all", scope.StrDefault.String(), 0, "catalog/frontend/list_allow_all", errors.NoKind},
		{"default//catalog/frontend/list_allow_all", scope.StrDefault.String(), 0, "catalog/frontend/list_allow_all", errors.NotValid},
		{"stores/123/catalog/index", "default", 0, "<nil>", errors.NotValid},
	}
	for i, test := range tests {
		havePath, haveErr := SplitFQ(test.have)

		if test.wantErrKind > 0 {
			assert.True(t, test.wantErrKind.Match(haveErr), "Index %d => Error: %s", i, haveErr)
		} else {
			assert.NoError(t, haveErr, "Test %v", test)
		}
		assert.Exactly(t, test.wantScope, havePath.ScopeID.Type().StrType(), "Index %d", i)
		assert.Exactly(t, test.wantScopeID, havePath.ScopeID.ID(), "Index %d", i)
		ls, _ := havePath.Level(-1)
		assert.Exactly(t, test.wantPath, ls, "Index %d", i)
	}
}

func TestSplitFQ2(t *testing.T) {
	p, err := SplitFQ("websites/5/web/cors/allow_credentials")
	if err != nil {
		t.Fatalf("%+v", err)
	}
	assert.Exactly(t, scope.Website.Pack(5), p.ScopeID)

	p, err = SplitFQ("default/0/web/cors/allow_credentials")
	if err != nil {
		t.Fatalf("%+v", err)
	}
	assert.Exactly(t, scope.DefaultTypeID, p.ScopeID)
}

func TestPathIsValid(t *testing.T) {

	tests := []struct {
		s           scope.Type
		id          int64
		have        string
		wantErrKind errors.Kind
	}{
		{scope.Default, 0, "//", errors.NotValid},
		{scope.Default, 0, "general/store_information/city", errors.NoKind},
		{scope.Default, 33, "general/store_information/city", errors.NoKind},
		{scope.Website, 33, "system/full_page_cache/varnish/backend_port", errors.NoKind},
		{scope.Default, 0, "", errors.Empty},
		{scope.Default, 0, "general/store_information", errors.NotValid},
		////{MustNew("system/dev/debug".Bind(scope.WebsiteID, 22).String()), ErrIncorrectPath},
		{scope.Default, 0, "groups/33/general/store_information/street", errors.NoKind},
		{scope.Default, 0, "groups/33", errors.NotValid},
		{scope.Default, 0, "system/dEv/inv˚lid", errors.NotValid},
		{scope.Default, 0, "system/dEv/inv'lid", errors.NotValid},
		{scope.Default, 0, "syst3m/dEv/invalid", errors.NoKind},
		{scope.Default, 0, "", errors.Empty},
	}
	for i, test := range tests {
		p := Path{
			ScopeID: scope.MakeTypeID(test.s, test.id),
			route:   test.have,
		}
		haveErr := p.IsValid()
		if test.wantErrKind > 0 {
			assert.True(t, test.wantErrKind.Match(haveErr), "Index %d => %s", i, haveErr)
		} else {
			assert.NoError(t, haveErr, "Index %d", i)
		}
	}
}

func TestPathValidateNaughtyStrings(t *testing.T) {
	// nothing is valid from the list of naughty strings
	for _, str := range naughtystrings.Unencoded() {
		_, err := MakePath(str)
		if err == nil {
			t.Errorf("Should not be valid: %q", str)
		}
	}
}

// func TestPathValidateNaughtyStrings(t *testing.T) {
// 	var valids slices.String = []string{"undefined", "undef", "null", "NULL", "nil", "NIL", "true", "false", "True", "False", "None", "hasOwnProperty", "0", "1", "1/2", "1E2", "1E02", "1/0", "0/0", "999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999", "NaN", "Infinity", "INF", "0x0", "0xffffffff", "0xffffffffffffffff", "0xabad1dea", "123456789012345678901234567890123456789", "01000", "08", "09", "_", "CON", "PRN", "AUX", "NUL", "COM1", "LPT1", "LPT2", "LPT3", "COM2", "COM3", "COM4", "evaluate", "mocha", "expression", "classic", "basement"}
// 	for _, str := range naughtystrings.Unencoded() {
// 		_, err := MakePath(str)
// 		if valids.Contains(str) && err != nil {
// 			t.Errorf("Should be valid %q but error: %+v", str, err)
// 		}
// 	}
// }

func TestPathRouteIsValid(t *testing.T) {

	p := Path{
		ScopeID: scope.MakeTypeID(scope.Store, 2),
		route:   `general/store_information`,
	}
	assert.True(t, errors.NotValid.Match(p.IsValid()))

	p = Path{
		ScopeID: scope.MakeTypeID(scope.Store, 2),
		route:   `general/store_information`,
	}
	assert.NoError(t, p.IsValid())
}

func TestPathHashWebsite(t *testing.T) {

	p := MustMakePath("general/single_store_mode/enabled").BindWebsite(33)
	hv, err := p.Hash(-1)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(p.String())
	check := fnv.New32a()
	_, cErr := check.Write([]byte(p.String()))
	assert.NoError(t, cErr)
	assert.Exactly(t, check.Sum32(), hv, "Have %d want %d", hv, check.Sum32())

}

func TestPathHashDefault(t *testing.T) {

	tests := []struct {
		have        string
		level       int
		wantHash    uint32
		wantErrKind errors.Kind
		wantLevel   string
	}{
		{"general/single_\x80store_mode/enabled", 0, 0, errors.NotValid, ""},
		{"general/single_store_mode/enabled", 0, 453736105, errors.NoKind, "default/0"},
		{"general/single_store_mode/enabled", 1, 2243014074, errors.NoKind, "default/0/general"},
		{"general/single_store_mode/enabled", 2, 4182795913, errors.NoKind, "default/0/general/single_store_mode"},
		{"general/single_store_mode/enabled", 3, 1584651487, errors.NoKind, "default/0/general/single_store_mode/enabled"},
		{"general/single_store_mode/enabled", -1, 1584651487, errors.NoKind, "default/0/general/single_store_mode/enabled"}, // 5
		{"general/single_store_mode/enabled", 5, 1584651487, errors.NoKind, "default/0/general/single_store_mode/enabled"},  // 6
		{"general/single_store_mode/enabled", 4, 1584651487, errors.NoKind, "default/0/general/single_store_mode/enabled"},  // 7
	}
	for i, test := range tests {
		p := Path{
			route: test.have,
		}

		hv, err := p.Hash(test.level)
		if test.wantErrKind > 0 {
			assert.True(t, test.wantErrKind.Match(err), "Index %d => %s", i, err)
			assert.Empty(t, hv, "Index %d", i)
			continue
		}
		assert.NoError(t, err, "Index %d", i)

		check := fnv.New32a()
		_, cErr := check.Write([]byte(test.wantLevel))
		assert.NoError(t, cErr)
		assert.Exactly(t, check.Sum32(), hv, "Want %d Have %d Index %d", check.Sum32(), hv, i)

		if test.level < 0 {
			test.level = -3
		}
		xrl, err := p.Level(test.level + 2)
		if err != nil {
			t.Fatal(err)
		}
		assert.Exactly(t, test.wantLevel, xrl, "Index %d", i)
		assert.Exactly(t, test.wantHash, hv, "Want %d Have %d Index %d", test.wantHash, hv, i)
	}
}

func TestPathCloneAppend(t *testing.T) {

	rs := "aa/bb/cc"
	pOrg := MustMakePath(rs)
	pOrg = pOrg.BindStore(3141)

	pAssigned := pOrg
	assert.Exactly(t, pOrg, pAssigned)
	pOrg = pOrg.Append("dd")
	assert.NotEqual(t, pOrg, pAssigned)
}

func TestPath_BindStore(t *testing.T) {
	p := MustMakePath(`aa/bb/cc`)
	p = p.BindStore(33)
	assert.Exactly(t, scope.MakeTypeID(scope.Store, 33), p.ScopeID)
}

func TestPath_BindWebsite(t *testing.T) {
	p := MustMakePath(`aa/bb/cc`)
	p = p.BindWebsite(44)
	assert.Exactly(t, scope.MakeTypeID(scope.Website, 44), p.ScopeID)
}

var _ sort.Interface = (*PathSlice)(nil)

func TestPathSlice_Contains(t *testing.T) {

	tests := []struct {
		paths  PathSlice
		search Path
		want   bool
	}{
		{
			PathSlice{
				0: MustMakePath("aa/bb/cc").BindWebsite(3),
				1: MustMakePath("aa/bb/cc").BindWebsite(2),
			},
			MustMakePath("aa/bb/cc").BindWebsite(2),
			true,
		},
		{
			PathSlice{
				0: MustMakePath("aa/bb/cc").BindWebsite(3),
				1: MustMakePath("aa/bb/cc").BindWebsite(2),
			},
			MustMakePath("aa/bb/cc").BindStore(2),
			false,
		},
	}
	for i, test := range tests {
		assert.Exactly(t, test.want, test.paths.Contains(test.search), "Index %d", i)
	}
}

func TestPathSlice_Sort(t *testing.T) {

	ps := PathSlice{
		MustMakePath("bb/cc/dd"),
		MustMakePath("xx/yy/zz"),
		MustMakePath("aa/bb/cc"),
	}
	ps.Sort()
	want := PathSlice{
		Path{route: `aa/bb/cc`, ScopeID: scope.DefaultTypeID},
		Path{route: `bb/cc/dd`, ScopeID: scope.DefaultTypeID},
		Path{route: `xx/yy/zz`, ScopeID: scope.DefaultTypeID},
	}
	assert.Exactly(t, want, ps)
}

// BenchmarkPathSlice_Sort-4	 1000000	      1987 ns/op	     480 B/op	       8 allocs/op
func BenchmarkPathSlice_Sort(b *testing.B) {
	// allocs are here uninteresting
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		ps := PathSlice{
			MustMakePath("rr/ss/tt"),
			MustMakePath("bb/cc/dd"),
			MustMakePath("xx/yy/zz"),
			MustMakePath("aa/bb/cc"),
			MustMakePath("ff/gg/hh"),
			MustMakePath("cc/dd/ee"),
		}
		ps.Sort()
		if len(ps) != 6 {
			b.Fatal("Incorrect length of ps variable after sorting")
		}
	}

}

func TestPathEqual(t *testing.T) {

	tests := []struct {
		a    string
		b    string
		want bool
	}{
		{"", "", true},
		{"a", "a", true},
		{"a", "b", false},
		{"a\x80", "a", false},
		{"general/single_\x80store_mode/enabled", "general/single_store_mode/enabled", false},
		{"general/single_store_mode/enabled", "general/single_store_mode/enabled", true},
		{"", "", true},
	}
	for i, test := range tests {
		assert.Exactly(t, test.want, test.a == test.b, "Index %d", i)
	}
}

func TestPathLevel(t *testing.T) {

	tests := []struct {
		have  string
		depth int
		want  string
	}{
		{"general/single_store_mode/enabled", 0, "<nil>"},
		{"general/single_store_mode/enabled", 1, "general"},
		{"general/single_store_mode/enabled", 2, "general/single_store_mode"},
		{"general/single_store_mode/enabled", 3, "general/single_store_mode/enabled"},
		{"general/single_store_mode/enabled", -1, "general/single_store_mode/enabled"},
		{"general/single_store_mode/enabled", 5, "general/single_store_mode/enabled"},
		{"general/single_store_mode/enabled", 4, "general/single_store_mode/enabled"},
		{"system/full_page_cache/varnish/backend_port", 3, "system/full_page_cache/varnish"},
	}
	for i, test := range tests {
		p := MustMakePath(test.have)
		r, err := p.Level(test.depth)
		assert.NoError(t, err)
		assert.Exactly(t, test.want, r, "Index %d", i)
	}
}

func TestPathHash(t *testing.T) {
	t.Parallel()
	tests := []struct {
		have        string
		depth       int
		wantHash    uint32
		wantErrKind errors.Kind
		wantLevel   string
	}{
		{"general/single_\x80store_mode/enabled", 0, 0, errors.NotValid, ""},
		{"general/single_store_mode/enabled", 0, 2166136261, errors.NoKind, "<nil>"},
		{"general/single_store_mode/enabled", 1, 616112491, errors.NoKind, "general"},
		{"general/single_store_mode/enabled", 2, 2274889228, errors.NoKind, "general/single_store_mode"},
		{"general/single_store_mode/enabled", 3, 1644245266, errors.NoKind, "general/single_store_mode/enabled"},
		{"general/single_store_mode/enabled", -1, 1644245266, errors.NoKind, "general/single_store_mode/enabled"},
		{"general/single_store_mode/enabled", 5, 1644245266, errors.NoKind, "general/single_store_mode/enabled"},
		{"general/single_store_mode/enabled", 4, 1644245266, errors.NoKind, "general/single_store_mode/enabled"},
	}
	for i, test := range tests {
		p := MustMakePath(test.have)
		hv, err := p.Hash(test.depth)
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

		l, err := p.Level(test.depth)
		assert.NoError(t, err)
		assert.Exactly(t, test.wantLevel, l, "Index %d", i)
		assert.Exactly(t, test.wantHash, hv, "Expected %d Actual %d Index %d", test.wantHash, hv, i)
	}
}

func TestPathPartPosition(t *testing.T) {

	tests := []struct {
		have     string
		pos      int
		wantPart string
		wantErr  bool
	}{
		{"general/single_\x80store_mode/enabled", 0, "", true},
		{"general/single_store_mode/enabled", 0, "", true},
		{"general/single_store_mode/enabled", 1, "general", false},
		{"general/single_store_mode/enabled", 2, "single_store_mode", false},
		{"general/single_store_mode/enabled", 3, "enabled", false},
		{"system/full_page_cache/varnish/backend_port", 4, "backend_port", false},
		{"general/single_store_mode/enabled", -1, "", true},
		{"general/single_store_mode/enabled", 5, "", true},
		{"general/single/store/website/group/mode/enabled/disabled/default", 5, "group", false},
	}
	for i, test := range tests {
		p := MustMakePath(test.have)
		part, haveErr := p.Part(test.pos)
		if test.wantErr {
			assert.Empty(t, part, "Index %d", i)
			assert.True(t, errors.NotValid.Match(haveErr), "Index %d => %s", i, haveErr)
			continue
		}
		assert.Exactly(t, test.wantPart, part, "Index %d", i)
	}
}

func TestPathValidate(t *testing.T) {

	tests := []struct {
		have        string
		wantErrKind errors.Kind
	}{
		{"//", errors.NotValid},
		{"general/store_information/city", errors.NoKind},
		{"general/store_information/city", errors.NoKind},
		{"system/full_page_cache/varnish/backend_port", errors.NoKind},
		{"", errors.Empty},
		{"general/store_information", errors.NoKind},
		////{MustNew("system/dev/debug".Bind(scope.WebsiteID, 22).String()), ErrIncorrectPath},
		{"groups/33/general/store_information/street", errors.NoKind},
		{"groups/33", errors.NoKind},
		{"system/dEv/inv˚lid", errors.NotValid},
		{"system/dEv/inv'lid", errors.NotValid},
		{"syst3m/dEv/invalid", errors.NoKind},
		{"", errors.Empty},
	}
	for i, test := range tests {
		_, haveErr := MakePath(test.have)
		if test.wantErrKind > 0 {
			assert.True(t, test.wantErrKind.Match(haveErr), "Index %d => %s", i, haveErr)
		} else {
			assert.NoError(t, haveErr, "Index %d", i)
		}
	}
}

func TestPathSplit(t *testing.T) {

	tests := []struct {
		have        string
		wantPart    []string
		wantErrKind errors.Kind
	}{
		{"general/single_\x80store_mode", []string{"general", "single_\x80store_mode"}, errors.NoKind},
		{"general/single_store_mode", []string{"general", "single_store_mode"}, errors.NoKind},
		{"general", nil, errors.NotValid},
		{"general/single_store_mode/enabled", []string{"general", "single_store_mode", "enabled"}, errors.NoKind},
		{"system/full_page_cache/varnish/backend_port", []string{"system", "full_page_cache", "varnish/backend_port"}, errors.NoKind},
	}
	for i, test := range tests {
		p := MustMakePath(test.have)
		sps, haveErr := p.Split()
		if test.wantErrKind > 0 {
			assert.True(t, test.wantErrKind.Match(haveErr), "Index %d => %s", i, haveErr)
			assert.Empty(t, sps[0], "Index %d", i)
			assert.Empty(t, sps[1], "Index %d", i)
			continue
		}
		for i, wantPart := range test.wantPart {
			assert.Exactly(t, wantPart, sps[i], "Index %d", i)
		}
	}
}
