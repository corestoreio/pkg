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
	"sort"
	"testing"

	"github.com/corestoreio/errors"
	"github.com/corestoreio/pkg/store/scope"
	"github.com/corestoreio/pkg/util/naughtystrings"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var (
	_ encoding.TextMarshaler     = (*Path)(nil)
	_ encoding.TextUnmarshaler   = (*Path)(nil)
	_ encoding.BinaryMarshaler   = (*Path)(nil)
	_ encoding.BinaryUnmarshaler = (*Path)(nil)
	_ fmt.Stringer               = (*Path)(nil)
)

func assertStrErr(t *testing.T, want string, msgAndArgs ...interface{}) func(string, error) {
	return func(s string, err error) {
		require.NoError(t, err, "%+v", err)
		assert.Exactly(t, want, s, msgAndArgs...)
	}
}

func TestMakePathWithScope(t *testing.T) {
	t.Parallel()
	t.Run("ok", func(t *testing.T) {
		p, err := NewPathWithScope(scope.Store.WithID(23), "aa/bb/cc")
		require.NoError(t, err)
		assert.Exactly(t, "stores/23/aa/bb/cc", p.String())
	})
	t.Run("fails", func(t *testing.T) {
		p, err := NewPathWithScope(scope.Store.WithID(23), "")
		assert.True(t, errors.Empty.Match(err), "%+v", err)
		assert.Nil(t, p)
	})
}

func TestNewByParts(t *testing.T) {
	t.Parallel()
	tests := []struct {
		path        string
		want        string
		wantErrKind errors.Kind
	}{
		{"aa/bb/cc", "default/0/aa/bb/cc", errors.NoKind},
		{"aa/bb/c", "aa/bb/cc", errors.NotValid},
		{"", "", errors.Empty},
	}
	for i, test := range tests {
		haveP, haveErr := NewPath(test.path)
		if test.wantErrKind > 0 {
			assert.Nil(t, haveP, "Index %d", i)
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
	_ = MustNewPath("a/\x80/c")
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
	p := MustNewPath("aa/bb/cc")
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
		haveP, haveErr := NewPath(test.route)
		if test.wantErrKind > 0 {
			assert.True(t, test.wantErrKind.Match(haveErr), "Index %d", i)
			continue
		}
		haveP = haveP.Bind(test.s.WithID(test.id))
		assertStrErr(t, test.wantFQ, "Index %d", i)(haveP.FQ())
	}
}

func TestPath_FQ(t *testing.T) {
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
		{250, 0, "system/dev/debug", scope.StrDefault.String() + "/0/system/dev/debug", errors.NotValid},
		{250, 343, "system/dev/debug", scope.StrDefault.String() + "/0/system/dev/debug", errors.NotValid},
	}
	for i, test := range tests {
		p, pErr := NewPath(test.route)
		if pErr != nil {
			assert.True(t, test.wantErrKind.Match(pErr), "Index %d => %s", i, pErr)
			continue
		}
		p = p.Bind(test.scp.WithID(test.id))
		have, haveErr := p.FQ()
		if test.wantErrKind > 0 {
			assert.Empty(t, have, "Index %d", i)
			assert.True(t, test.wantErrKind.Match(haveErr), "Index %d => %s", i, haveErr)
			continue
		}
		assert.NoError(t, haveErr, "Index %d", i)
		assert.Exactly(t, test.want, have, "Index %d", i)
	}

	r := "catalog/frontend/list_allow_all"
	assert.Exactly(t, "stores/7475/catalog/frontend/list_allow_all", MustNewPath(r).BindStore(7475).String())
	p := MustNewPath(r).BindStore(5)
	assert.Exactly(t, "stores/5/catalog/frontend/list_allow_all", p.String())
}

func TestShouldNotPanicBecauseOfIncorrectStrScope(t *testing.T) {
	t.Parallel()
	assert.Exactly(t, "stores/345/xxxxx/yyyyy/zzzzz", MustNewPath("xxxxx/yyyyy/zzzzz").BindStore(345).String())
	defer func() {
		if r := recover(); r != nil {
			t.Fatal("Did not expect a panic")
		}
	}()
	_ = MustNewPath("xxxxx/yyyyy/zzzzz").Bind(345)
}

func TestShouldPanicIncorrectPath(t *testing.T) {
	t.Parallel()
	assert.Exactly(t, "default/0/xxxxx/yyyyy/zzzzz", MustNewPath("xxxxx/yyyyy/zzzzz").Bind(345).String())
	defer func() {
		if r := recover(); r != nil {
			err := r.(error)
			assert.True(t, errors.NotValid.Match(err))
		} else {
			t.Fatal("Expecting a panic")
		}
	}()
	assert.Exactly(t, "websites/345/xxxxx/yyyyy", MustNewPath("xxxxx/yyyyy").BindWebsite(345).String())
}

func TestPath_ParseStrings(t *testing.T) {
	t.Parallel()

	tests := []struct {
		scp, id, route string
		wantPath       string
		wantErr        errors.Kind
	}{
		{"default", "0", "aa/bb/cc", "default/0/aa/bb/cc", errors.NoKind},
		{"stores", "1", "aa/bb/cc", "stores/1/aa/bb/cc", errors.NoKind},
		{"websites", "1", "aa/bb/cc", "websites/1/aa/bb/cc", errors.NoKind},
		{"website", "1", "aa/bb/cc", "", errors.NotValid},
		{"websites", "-1", "aa/bb/cc", "", errors.CorruptData},
	}

	for i, test := range tests {
		p := new(Path)
		haveErr := p.ParseStrings(test.scp, test.id, test.route)
		if test.wantErr > 0 {
			assert.True(t, test.wantErr.Match(haveErr), "IDX:%d %+v", i, haveErr)
			continue
		}
		assert.Exactly(t, test.wantPath, p.String(), "index %d", i)
	}
}

func TestPath_Parse(t *testing.T) {
	// t.Parallel()
	tests := []struct {
		have        string
		wantScope   string
		wantScopeID int64
		wantPath    string
		wantErrKind errors.Kind
	}{
		{"catalog/frontend", "", 0, "", errors.NotValid},
		{"catalog/frontend/list_allow_all", "default", 0, "default/0/catalog/frontend/list_allow_all", errors.NoKind},
		{"catalog/frontend/list/allow_all", "default", 0, "default/0/catalog/frontend/list/allow_all", errors.NoKind},
		{"groups/1/catalog/frontend/list_allow_all", "default", 0, "", errors.NotSupported},
		{"stores/7475/catalog/frontend/list_allow_all", scope.StrStores.String(), 7475, "stores/7475/catalog/frontend/list_allow_all", errors.NoKind},
		{"stores/4/system/full_page_cache/varnish/backend_port", scope.StrStores.String(), 4, "stores/4/system/full_page_cache/varnish/backend_port", errors.NoKind},
		{"websites/1/catalog/frontend/list_allow_all", scope.StrWebsites.String(), 1, "websites/1/catalog/frontend/list_allow_all", errors.NoKind},
		{"default/0/catalog/frontend/list_allow_all", scope.StrDefault.String(), 0, "default/0/catalog/frontend/list_allow_all", errors.NoKind},
		{"default//catalog/frontend/list_allow_all", scope.StrDefault.String(), 0, "default/0/catalog/frontend/list_allow_all", errors.NotValid},
		{"stores/123/catalog/index", "default", 0, "", errors.NotValid},
	}
	havePath := new(Path)
	for i, test := range tests {
		haveErr := havePath.Parse(test.have)

		if test.wantErrKind > 0 {
			assert.True(t, test.wantErrKind.Match(haveErr), "Index %d => Error: %s", i, haveErr)
		} else {
			require.NoError(t, haveErr, "Test %v", test)
			assert.Exactly(t, test.wantScope, havePath.ScopeID.Type().StrType(), "Index %d", i)
			assert.Exactly(t, test.wantScopeID, havePath.ScopeID.ID(), "Index %d", i)
			ls, _ := havePath.Level(-1)
			assert.Exactly(t, test.wantPath, ls, "Index %d", i)
		}
	}
}

func TestPath_SplitFQ2(t *testing.T) {
	t.Parallel()
	p := new(Path)

	if err := p.Parse("websites/5/web/cors/allow_credentials"); err != nil {
		t.Fatalf("%+v", err)
	}
	assert.Exactly(t, scope.Website.WithID(5), p.ScopeID)

	if err := p.Parse("default/0/web/cors/allow_credentials"); err != nil {
		t.Fatalf("%+v", err)
	}
	assert.Exactly(t, scope.DefaultTypeID, p.ScopeID)
}

func TestPathIsValid(t *testing.T) {
	t.Parallel()
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
	t.Parallel()
	// nothing is valid from the list of naughty strings
	for _, str := range naughtystrings.Unencoded() {
		_, err := NewPath(str)
		if err == nil {
			t.Errorf("Should not be valid: %q", str)
		}
	}
}

func TestPathRouteIsValid(t *testing.T) {
	t.Parallel()
	t.Run("too short", func(t *testing.T) {
		p := Path{
			ScopeID: scope.MakeTypeID(scope.Store, 2),
			route:   `general/store_information`,
		}
		assert.True(t, errors.NotValid.Match(p.IsValid()))
	})
	t.Run("starts with default", func(t *testing.T) {
		p := Path{
			ScopeID: scope.MakeTypeID(scope.Store, 2),
			route:   `default/general/store_information`,
		}
		assert.True(t, errors.NotValid.Match(p.IsValid()))
	})
	t.Run("starts with websites", func(t *testing.T) {
		p := Path{
			ScopeID: scope.MakeTypeID(scope.Store, 2),
			route:   `websites/general/store_information`,
		}
		assert.True(t, errors.NotValid.Match(p.IsValid()))
	})
	t.Run("starts with stores", func(t *testing.T) {
		p := Path{
			ScopeID: scope.MakeTypeID(scope.Store, 2),
			route:   `stores/general/store_information`,
		}
		assert.True(t, errors.NotValid.Match(p.IsValid()))
	})

	t.Run("contains default is allowed", func(t *testing.T) {
		p := Path{
			ScopeID: scope.MakeTypeID(scope.Store, 2),
			route:   `general/default/store_information`,
		}
		assert.NoError(t, p.IsValid())
	})
	t.Run("contains websites is allowed", func(t *testing.T) {
		p := Path{
			ScopeID: scope.MakeTypeID(scope.Store, 2),
			route:   `general/websites/store_information`,
		}
		assert.NoError(t, p.IsValid())
	})
	t.Run("contains stores is allowed", func(t *testing.T) {
		p := Path{
			ScopeID: scope.MakeTypeID(scope.Store, 2),
			route:   `general/stores/store_information`,
		}
		assert.NoError(t, p.IsValid())
	})
}

func TestPathHashWebsite(t *testing.T) {
	t.Parallel()
	p := MustNewPath("general/single_store_mode/enabled").BindWebsite(33)
	hv := p.Hash64ByLevel(-1)

	t.Log(p.String())
	assert.Exactly(t, uint64(0xe79edc1df8f88eb0), hv, "hashes do not match")

}

func TestPath_Hash64ByLevel(t *testing.T) {
	t.Parallel()

	t.Run("level 4 different to full hash", func(t *testing.T) {
		p := MustNewPath("general/single_store_mode/enabled")
		hl := p.Hash64ByLevel(4)
		l, err := p.Level(4)
		require.NoError(t, err)
		assert.Exactly(t, "default/0/general/single_store_mode", l)
		assert.Exactly(t, uint64(0x5cbf04d81b7be97), hl)
		h := p.Hash64()
		assert.NotEqual(t, hl, h)
	})

	t.Run("different levels", func(t *testing.T) {
		tests := []struct {
			have      string
			level     int
			wantHash  uint64
			wantLevel string
		}{
			{"general/single_store_mode/enabled", 2, 0x9c8e9914e6663c89, "default/0"},
			{"general/single_store_mode/enabled", 1, 0x541e70f983cf2646, "default"},
			{"general/single_store_mode/enabled", 3, 0xe6fa13685c10e2da, "default/0/general"},
			{"general/single_store_mode/enabled", 4, 0x5cbf04d81b7be97, "default/0/general/single_store_mode"},
			{"general/single_store_mode/enabled", -1, 0xa317445783c664e3, "default/0/general/single_store_mode/enabled"},
			{"general/single_store_mode/enabled", 5, 0xa317445783c664e3, "default/0/general/single_store_mode/enabled"},
		}
		for i, test := range tests {
			p := &Path{
				route: test.have,
			}

			assert.Exactly(t, test.wantHash, p.Hash64ByLevel(test.level), "Index %d", i)

			if test.level < 0 {
				test.level = -3
			}
			xrl, err := p.Level(test.level)
			if err != nil {
				t.Fatal(err)
			}
			assert.Exactly(t, test.wantLevel, xrl, "Index %d", i)
		}
	})
}

func TestPath_BindStore(t *testing.T) {
	p := MustNewPath(`aa/bb/cc`)
	p = p.BindStore(33)
	assert.Exactly(t, scope.MakeTypeID(scope.Store, 33), p.ScopeID)
}

func TestPath_BindWebsite(t *testing.T) {
	p := MustNewPath(`aa/bb/cc`)
	p = p.BindWebsite(44)
	assert.Exactly(t, scope.MakeTypeID(scope.Website, 44), p.ScopeID)
}

var _ sort.Interface = (*PathSlice)(nil)

func TestPathSlice_Contains(t *testing.T) {
	t.Parallel()
	tests := []struct {
		paths  PathSlice
		search *Path
		want   bool
	}{
		{
			PathSlice{
				0: MustNewPath("aa/bb/cc").BindWebsite(3),
				1: MustNewPath("aa/bb/cc").BindWebsite(2),
			},
			MustNewPath("aa/bb/cc").BindWebsite(2),
			true,
		},
		{
			PathSlice{
				0: MustNewPath("aa/bb/cc").BindWebsite(3),
				1: MustNewPath("aa/bb/cc").BindWebsite(2),
			},
			MustNewPath("aa/bb/cc").BindStore(2),
			false,
		},
	}
	for i, test := range tests {
		assert.Exactly(t, test.want, test.paths.Contains(test.search), "Index %d", i)
	}
}

func TestPathSlice_Sort(t *testing.T) {
	t.Parallel()
	runner := func(have, want PathSlice) func(*testing.T) {
		return func(t *testing.T) {
			have.Sort()
			require.Exactly(t, want, have)
		}
	}

	t.Run("Default Scope", runner(
		PathSlice{
			MustNewPath("bb/cc/dd"),
			MustNewPath("xx/yy/zz"),
			MustNewPath("aa/bb/cc"),
		},
		PathSlice{
			&Path{route: `aa/bb/cc`, ScopeID: scope.DefaultTypeID},
			&Path{route: `bb/cc/dd`, ScopeID: scope.DefaultTypeID},
			&Path{route: `xx/yy/zz`, ScopeID: scope.DefaultTypeID},
		},
	))

	t.Run("Default+Website Scope", runner(
		PathSlice{
			MustNewPath("bb/cc/dd"),
			MustNewPathWithScope(scope.Website.WithID(3), "xx/yy/zz"),
			MustNewPathWithScope(scope.Website.WithID(1), "xx/yy/zz"),
			MustNewPathWithScope(scope.Website.WithID(2), "xx/yy/zz"),
			MustNewPathWithScope(scope.Website.WithID(1), "zz/aa/bb"),
			MustNewPathWithScope(scope.Website.WithID(2), "aa/bb/cc"),
			MustNewPath("aa/bb/cc"),
		},
		PathSlice{
			&Path{route: `aa/bb/cc`, ScopeID: scope.DefaultTypeID},
			&Path{route: `aa/bb/cc`, ScopeID: scope.Website.WithID(2)},
			&Path{route: `bb/cc/dd`, ScopeID: scope.DefaultTypeID},
			&Path{route: `xx/yy/zz`, ScopeID: scope.Website.WithID(1)},
			&Path{route: `xx/yy/zz`, ScopeID: scope.Website.WithID(2)},
			&Path{route: `xx/yy/zz`, ScopeID: scope.Website.WithID(3)},
			&Path{route: `zz/aa/bb`, ScopeID: scope.Website.WithID(1)},
		},
	))

	t.Run("Default+Website+Store Scope", runner(
		PathSlice{
			MustNewPathWithScope(scope.Store.WithID(3), "bb/cc/dd"),
			MustNewPathWithScope(scope.Store.WithID(2), "bb/cc/dd"),
			MustNewPath("bb/cc/dd"),
			MustNewPathWithScope(scope.Website.WithID(3), "xx/yy/zz"),
			MustNewPathWithScope(scope.Website.WithID(1), "xx/yy/zz"),
			MustNewPathWithScope(scope.Website.WithID(2), "xx/yy/zz"),
			MustNewPathWithScope(scope.Store.WithID(4), "zz/aa/bb"),
			MustNewPathWithScope(scope.Website.WithID(1), "zz/aa/bb"),
			MustNewPathWithScope(scope.Website.WithID(2), "aa/bb/cc"),
			MustNewPath("aa/bb/cc"),
		},
		PathSlice{
			&Path{route: `aa/bb/cc`, ScopeID: scope.DefaultTypeID},
			&Path{route: `aa/bb/cc`, ScopeID: scope.Website.WithID(2)},
			&Path{route: `bb/cc/dd`, ScopeID: scope.DefaultTypeID},
			&Path{route: `bb/cc/dd`, ScopeID: scope.Store.WithID(2)},
			&Path{route: `bb/cc/dd`, ScopeID: scope.Store.WithID(3)},
			&Path{route: `xx/yy/zz`, ScopeID: scope.Website.WithID(1)},
			&Path{route: `xx/yy/zz`, ScopeID: scope.Website.WithID(2)},
			&Path{route: `xx/yy/zz`, ScopeID: scope.Website.WithID(3)},
			&Path{route: `zz/aa/bb`, ScopeID: scope.Website.WithID(1)},
			&Path{route: `zz/aa/bb`, ScopeID: scope.Store.WithID(4)},
		},
	))

}

func TestPathEqual(t *testing.T) {
	t.Parallel()
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
	t.Parallel()
	tests := []struct {
		have  string
		depth int
		want  string
	}{
		{"general/single_store_mode/enabled", 0, ""},
		{"general/single_store_mode/enabled", 1, "default"},
		{"general/single_store_mode/enabled", 2, "default/0"},
		{"general/single_store_mode/enabled", 3, "default/0/general"},
		{"general/single_store_mode/enabled", 4, "default/0/general/single_store_mode"},
		{"general/single_store_mode/enabled", 5, "default/0/general/single_store_mode/enabled"},
		{"general/single_store_mode/enabled", -1, "default/0/general/single_store_mode/enabled"},
		{"system/full_page_cache/varnish/backend_port", 5, "default/0/system/full_page_cache/varnish"},
	}
	for i, test := range tests {
		p := MustNewPath(test.have)
		r, err := p.Level(test.depth)
		assert.NoError(t, err)
		assert.Exactly(t, test.want, r, "Index %d", i)
	}
}

func TestPathPartPosition(t *testing.T) {
	t.Parallel()

	t.Run("invalid route", func(t *testing.T) {
		p := &Path{
			route: "general/single_\x80store_mode/enabled",
		}
		part, haveErr := p.Part(0)
		assert.Empty(t, part)
		assert.True(t, errors.NotValid.Match(haveErr), "%+v", haveErr)
	})

	t.Run("valid routes", func(t *testing.T) {

		tests := []struct {
			have     string
			pos      int
			wantPart string
			wantErr  bool
		}{
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
			p := MustNewPath(test.have)
			part, haveErr := p.Part(test.pos)
			if test.wantErr {
				assert.Empty(t, part, "Index %d", i)
				assert.True(t, errors.NotValid.Match(haveErr), "Index %d => %s", i, haveErr)
				continue
			}
			assert.Exactly(t, test.wantPart, part, "Index %d", i)
		}
	})
}

func TestPathValidate(t *testing.T) {
	t.Parallel()
	tests := []struct {
		have        string
		wantErrKind errors.Kind
	}{
		{"//", errors.NotValid},
		{"general/store_information/city", errors.NoKind},
		{"general/store_information/city", errors.NoKind},
		{"system/full_page_cache/varnish/backend_port", errors.NoKind},
		{"", errors.Empty},
		{"general/store_information", errors.NotValid},
		////{MustNew("system/dev/debug".Bind(scope.WebsiteID, 22).String()), ErrIncorrectPath},
		{"groups/33/general/store_information/street", errors.NoKind},
		{"groups/33", errors.NotValid},
		{"system/dEv/inv˚lid", errors.NotValid},
		{"system/dEv/inv'lid", errors.NotValid},
		{"syst3m/dEv/invalid", errors.NoKind},
		{"", errors.Empty},
	}
	for i, test := range tests {
		_, haveErr := NewPath(test.have)
		if test.wantErrKind > 0 {
			assert.True(t, test.wantErrKind.Match(haveErr), "Index %d => %s", i, haveErr)
		} else {
			assert.NoError(t, haveErr, "Index %d", i)
		}
	}
}

func TestPathSplit(t *testing.T) {
	t.Parallel()
	t.Run("invalid", func(t *testing.T) {
		p := &Path{
			ScopeID: scope.Store.WithID(1),
			route:   "general",
		}
		sps, haveErr := p.Split()
		assert.Nil(t, sps)
		assert.True(t, errors.NotValid.Match(haveErr), "%+v", haveErr)
	})
	t.Run("valid paths", func(t *testing.T) {
		tests := []struct {
			have     string
			wantPart []string
		}{
			{"general/single_store_mode/xx", []string{"general", "single_store_mode", "xx"}},
			{"general/single_store_mode/enabled", []string{"general", "single_store_mode", "enabled"}},
			{"system/full_page_cache/varnish/backend_port", []string{"system", "full_page_cache", "varnish/backend_port"}},
		}
		for _, test := range tests {
			p := MustNewPath(test.have)
			sps, haveErr := p.Split()
			require.NoError(t, haveErr, "Path %q", test.have)
			for i, wantPart := range test.wantPart {
				assert.Exactly(t, wantPart, sps[i], "Index %d", i)
			}
		}
	})
}

func TestPath_MarshalText(t *testing.T) {
	t.Parallel()

	t.Run("two way, no errors", func(t *testing.T) {
		p, err := NewPathWithScope(scope.Store.WithID(4), "xx/yy/zz")
		require.NoError(t, err)
		txt, err := p.MarshalText()
		require.NoError(t, err)
		assert.Exactly(t, `stores/4/xx/yy/zz`, string(txt))

		var p2 Path
		err = p2.UnmarshalText(txt)
		require.NoError(t, err)
		assert.Exactly(t, `stores/4/xx/yy/zz`, p2.String())
	})

	t.Run("UnmarshalText invalid length", func(t *testing.T) {
		var p2 Path
		err := p2.UnmarshalText([]byte(`aa/bb`))
		assert.True(t, errors.NotValid.Match(err), "%+v", err)
		assert.EqualError(t, err, "[config] Incorrect fully qualified path: \"aa/bb\". Expecting: strScope/ID/aa/bb")
	})
	t.Run("UnmarshalText invalid scope", func(t *testing.T) {
		var p2 Path
		err := p2.UnmarshalText([]byte(`scopeX/123/aa/bb/cc/dd`))
		assert.True(t, errors.NotSupported.Match(err), "%+v", err)
		assert.EqualError(t, err, "[config] Unknown Scope: \"scopeX\"")
	})
	t.Run("UnmarshalText failed to parse scope id", func(t *testing.T) {
		var p2 Path
		err := p2.UnmarshalText([]byte(`websites/x/aa/bb/cc/dd`))
		assert.True(t, errors.NotValid.Match(err), "%+v", err)
		assert.EqualError(t, err, "[config] ParseInt: strconv.ParseInt: parsing \"x\": invalid syntax")
	})
}

func TestPath_MarshalBinary(t *testing.T) {
	t.Parallel()

	t.Run("two way, no errors", func(t *testing.T) {
		p, err := NewPathWithScope(scope.Store.WithID(4), "xx/yy/zz")
		require.NoError(t, err)
		txt, err := p.MarshalBinary()
		require.NoError(t, err)
		assert.Exactly(t, "\x04\x00\x00\x04\x00\x00\x00\x00xx/yy/zz", string(txt))

		var p2 Path
		err = p2.UnmarshalBinary(txt)
		require.NoError(t, err)
		assert.Exactly(t, `stores/4/xx/yy/zz`, p2.String())
	})

	t.Run("UnmarshalBinary invalid length", func(t *testing.T) {
		var p2 Path
		err := p2.UnmarshalBinary([]byte(`aa/bb`))
		assert.True(t, errors.TooShort.Match(err), "%+v", err)
		assert.EqualError(t, err, "[config] UnmarshalBinary: input data too short")
	})
	t.Run("UnmarshalBinary invalid scope", func(t *testing.T) {
		var p2 Path
		err := p2.UnmarshalBinary([]byte(`scopeX/123/aa/bb/cc/dd`))
		assert.True(t, errors.NotValid.Match(err), "%+v", err)
		assert.EqualError(t, err, "[config] Route \"23/aa/bb/cc/dd\" contains invalid ScopeID: \"Type(Type(112)) ID(7299955)\"")
	})
	t.Run("UnmarshalBinary failed to parse scope id", func(t *testing.T) {
		var p2 Path
		err := p2.UnmarshalBinary([]byte(`websites/x/aa/bb/cc/dd`))
		assert.True(t, errors.NotValid.Match(err), "%+v", err)
		assert.EqualError(t, err, "[config] Route \"/x/aa/bb/cc/dd\" contains invalid ScopeID: \"Type(Type(115)) ID(6448503)\"")
	})
}

func TestPath_IsEmpty(t *testing.T) {
	t.Parallel()
	var p Path
	p.ScopeID = scope.Store.WithID(122)
	assert.True(t, p.IsEmpty())
}

func TestPath_Equal(t *testing.T) {
	t.Parallel()

	p1, err := NewPathWithScope(scope.Store.WithID(4), "xx/yy/zz")
	require.NoError(t, err)

	p2, err := NewPathWithScope(scope.Store.WithID(4), "xx/yy/zz")
	require.NoError(t, err)

	assert.True(t, p1.Equal(p2))

	p2.ScopeID = scope.Website.WithID(5)
	assert.True(t, p1.EqualRoute(p2))
}

func TestPath_HasRoutePrefix(t *testing.T) {
	t.Parallel()
	p := &Path{route: `xx/yy/zz`, ScopeID: scope.Website.WithID(3)}

	assert.False(t, p.RouteHasPrefix(""))
	assert.True(t, p.RouteHasPrefix("xx"))
	assert.False(t, p.RouteHasPrefix("yy"))
	assert.True(t, p.RouteHasPrefix("xx/yy/zz"))
}

func TestPath_EnvName(t *testing.T) {
	t.Parallel()
	t.Run("String and Binary", func(t *testing.T) {
		p := MustNewPath("tt/ww/de").WithEnvSuffix()
		p.envSuffix = "STAGING"
		assertStrErr(t, `default/0/tt/ww/de/STAGING`)(p.FQ())

		d, err := p.BindStore(3).MarshalBinary()
		require.NoError(t, err)
		assert.Exactly(t, "\x03\x00\x00\x04\x00\x00\x00\x00tt/ww/de/STAGING", string(d))

		p.UseEnvSuffix = false
		assertStrErr(t, `default/0/tt/ww/de`)(p.FQ())
	})
	t.Run("Parse", func(t *testing.T) {
		p := &Path{
			envSuffix: "STAGING",
		}
		err := p.Parse(`stores/3/tt/ww/de/STAGING`)
		require.NoError(t, err, "%+v", err)
		assertStrErr(t, `stores/3/tt/ww/de`)(p.FQ())

		err = p.Parse(`stores/3/tt/ww/de`)
		require.NoError(t, err, "%+v", err)
		assertStrErr(t, `stores/3/tt/ww/de`)(p.FQ())
	})
	t.Run("ParseStrings", func(t *testing.T) {
		p := &Path{
			envSuffix: "STAGING",
		}
		err := p.ParseStrings(`stores`, "3", `tt/ww/de/STAGING`)
		require.NoError(t, err, "%+v", err)
		assertStrErr(t, `stores/3/tt/ww/de`)(p.FQ())

		err = p.ParseStrings(`stores`, "3", `tt/ww/de`)
		require.NoError(t, err, "%+v", err)
		assertStrErr(t, `stores/3/tt/ww/de`)(p.FQ())
	})
	t.Run("UnmarshalText", func(t *testing.T) {
		p := &Path{
			envSuffix: "STAGING",
		}
		err := p.UnmarshalText([]byte(`stores/3/tt/ww/de/STAGING`))
		require.NoError(t, err, "%+v", err)
		assertStrErr(t, `stores/3/tt/ww/de`)(p.FQ())

		err = p.UnmarshalText([]byte(`stores/3/tt/ww/de`))
		require.NoError(t, err, "%+v", err)
		assertStrErr(t, `stores/3/tt/ww/de`)(p.FQ())
	})
	t.Run("UnmarshalBinary", func(t *testing.T) {
		p := &Path{
			envSuffix: "STAGING",
		}
		err := p.UnmarshalBinary([]byte("\x03\x00\x00\x04\x00\x00\x00\x00tt/ww/de/STAGING"))
		require.NoError(t, err, "%+v", err)
		assertStrErr(t, `stores/3/tt/ww/de`)(p.FQ())

		err = p.UnmarshalBinary([]byte("\x03\x00\x00\x04\x00\x00\x00\x00tt/ww/de"))
		require.NoError(t, err, "%+v", err)
		assertStrErr(t, `stores/3/tt/ww/de`)(p.FQ())
	})
}
