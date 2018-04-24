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

package config_test

import (
	"strings"
	"testing"
	"time"

	"github.com/corestoreio/errors"
	"github.com/corestoreio/pkg/config"
	"github.com/corestoreio/pkg/store/scope"
	"github.com/corestoreio/pkg/util/conv"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var (
	_ config.Getter     = (*config.Service)(nil)
	_ config.Writer     = (*config.Service)(nil)
	_ config.Subscriber = (*config.Service)(nil)
)

func TestNewServiceStandard(t *testing.T) {

	srv := config.MustNewService(config.NewInMemoryStore())
	assert.NotNil(t, srv)
}

func TestNotKeyNotFoundError(t *testing.T) {

	srv := config.MustNewService(config.NewInMemoryStore())
	assert.NotNil(t, srv)

	scopedSrv := srv.NewScoped(1, 1)

	flat := scopedSrv.Value(scope.Default, "catalog/product/enable_flat")
	assert.False(t, flat.IsValid(), "Should not find the key")
	assert.True(t, flat.IsEmpty(), "should be empty")

	val := scopedSrv.Value(scope.Store, "catalog")
	assert.False(t, val.IsValid())
	assert.True(t, val.IsEmpty(), "should be empty")
}

func TestService_Write(t *testing.T) {

	srv := config.MustNewService(config.NewInMemoryStore())
	assert.NotNil(t, srv)

	p1 := config.Path{}
	err := srv.Write(p1, []byte{})
	assert.True(t, errors.Empty.Match(err), "Error: %s", err)
}

func TestService_Write_Get_Value_Success(t *testing.T) {

	runner := func(p config.Path, value []byte) func(*testing.T) {
		return func(t *testing.T) {
			srv := config.MustNewService(config.NewInMemoryStore())

			require.NoError(t, srv.Write(p, value), "Writing Value in Test %q should not fail", t.Name())

			haveStr, ok, haveErr := srv.Value(p).Str()
			require.NoError(t, haveErr, "No error should occur when retrieving a value")
			assert.True(t, ok)
			assert.Exactly(t, string(value), haveStr)

		}
	}

	basePath := config.MustMakePath("aa/bb/cc")

	t.Run("stringDefault", runner(basePath, []byte("Gopher")))
	t.Run("stringWebsite", runner(basePath.BindWebsite(10), []byte("Gopher")))
	t.Run("stringStore", runner(basePath.BindStore(22), []byte("Gopher")))

}

func TestScoped_ScopeIDs(t *testing.T) {
	scp := config.Scoped{WebsiteID: 3, StoreID: 4}
	assert.Exactly(t, scope.TypeIDs{scope.Store.Pack(4), scope.Website.Pack(3)}, scp.ScopeIDs())
}

func TestScoped_IsValid(t *testing.T) {
	t.Parallel()
	cfg := config.NewMock()
	tests := []struct {
		s    config.Scoped
		want bool
	}{
		{config.Scoped{}, false},
		{config.Scoped{WebsiteID: 1, StoreID: 0}, false},
		{config.Scoped{WebsiteID: 1, StoreID: 1}, false},
		{config.Scoped{WebsiteID: 0, StoreID: 1}, false},
		{config.Scoped{Root: cfg}, true},
		{config.Scoped{Root: cfg, WebsiteID: 1, StoreID: 0}, true},
		{config.Scoped{Root: cfg, WebsiteID: 1, StoreID: 1}, true},
		{config.Scoped{Root: cfg, WebsiteID: 0, StoreID: 1}, false},
		{config.Scoped{Root: cfg, WebsiteID: 1, StoreID: -1}, false},
		{config.Scoped{Root: cfg, WebsiteID: -1, StoreID: -1}, false},
	}
	for i, test := range tests {
		if have, want := test.s.IsValid(), test.want; have != want {
			t.Errorf("Idx %d => Have: %v Want: %v", i, have, want)
		}
	}
}

func TestScopedServiceScope(t *testing.T) {
	t.Parallel()
	tests := []struct {
		websiteID, storeID int64
		wantScope          scope.Type
		wantID             int64
	}{
		{0, 0, scope.Default, 0},
		{1, 0, scope.Website, 1},
		{1, 3, scope.Store, 3},
		{0, 3, scope.Store, 3},
	}
	for i, test := range tests {
		sg := config.NewMock().NewScoped(test.websiteID, test.storeID)
		haveScope, haveID := sg.ScopeID().Unpack()
		assert.Exactly(t, test.wantScope, haveScope, "Index %d", i)
		assert.Exactly(t, test.wantID, haveID, "Index %d", i)
	}
}

func TestScopedServicePath(t *testing.T) {
	t.Parallel()
	basePath := config.MustMakePath("aa/bb/cc")
	tests := []struct {
		desc               string
		fqpath             string
		route              string
		perm               scope.Type
		websiteID, storeID int64
		wantErrKind        errors.Kind
		wantTypeIDs        scope.TypeIDs
	}{
		{
			"Default ScopedGetter should return default scope",
			basePath.String(), "aa/bb/cc", scope.Absent, 0, 0, errors.NoKind, scope.TypeIDs{scope.DefaultTypeID},
		},
		{
			"Website ID 1 ScopedGetter should fall back to default scope",
			basePath.String(), "aa/bb/cc", scope.Website, 1, 0, errors.NoKind, scope.TypeIDs{scope.DefaultTypeID, scope.Website.Pack(1)},
		},
		{
			"Website ID 10 ScopedGetter should fall back to website 10 scope",
			basePath.BindWebsite(10).String(), "aa/bb/cc", scope.Website, 10, 0, errors.NoKind, scope.TypeIDs{scope.Website.Pack(10)},
		},
		{
			"Website ID 10 + Store 22 ScopedGetter should fall back to website 10 scope",
			basePath.BindWebsite(10).String(), "aa/bb/cc", scope.Store, 10, 22, errors.NoKind, scope.TypeIDs{scope.Website.Pack(10), scope.Store.Pack(22)},
		},
		{
			"Website ID 10 + Store 22 ScopedGetter should return Store 22 scope",
			basePath.BindStore(22).String(), "aa/bb/cc", scope.Store, 10, 22, errors.NoKind, scope.TypeIDs{scope.Store.Pack(22)},
		},
		{
			"Website ID 10 + Store 42 ScopedGetter should return nothing",
			basePath.BindStore(22).String(), "aa/bb/cc", scope.Store, 10, 42, errors.NotFound, scope.TypeIDs{scope.DefaultTypeID},
		},
		{
			"Path consists of only two elements which is incorrect",
			basePath.String(), "aa/bb", scope.Store, 0, 0, errors.NotValid, nil,
		},
	}

	// vals stores all possible types for which we have functions in config.ScopedGetter
	vals := []interface{}{"Gopher", true, float64(3.14159), int(2016), time.Now(), []byte(`Hellö Dear Goph€rs`), time.Hour}
	const customTestTimeFormat = "2006-01-02 15:04:05"

	for vi, wantVal := range vals {
		for xi, test := range tests {

			fqVal := conv.ToString(wantVal)
			if wv, ok := wantVal.(time.Time); ok {
				fqVal = wv.Format(customTestTimeFormat)
			}
			cg := config.NewMock(config.MockPathValue{
				test.fqpath: fqVal,
			})

			sg := cg.NewScoped(test.websiteID, test.storeID)
			haveVal := sg.Value(test.perm, test.route)

			if test.wantErrKind > 0 {
				require.False(t, haveVal.IsValid(), "Index %d/%d scoped path value must be found ", vi, xi)
				continue
			}

			switch wv := wantVal.(type) {
			case []byte:
				var buf strings.Builder
				_, err := haveVal.WriteTo(&buf)
				require.NoError(t, err, "Error: %+v\n\n%s", err, test.desc)
				assert.Exactly(t, string(wv), buf.String(), test.desc)

			case string:
				hs, _, err := haveVal.Str()
				require.NoError(t, err)
				assert.Exactly(t, wv, hs)
			case bool:
				hs, _, err := haveVal.Bool()
				require.NoError(t, err)
				assert.Exactly(t, wv, hs)
			case float64:
				hs, _, err := haveVal.Float64()
				require.NoError(t, err)
				assert.Exactly(t, wv, hs)
			case int:
				hs, _, err := haveVal.Int()
				require.NoError(t, err)
				assert.Exactly(t, wv, hs)
			case time.Time:
				hs, _, err := haveVal.Time()
				require.NoError(t, err)
				assert.Exactly(t, wv.Format(customTestTimeFormat), hs.Format(customTestTimeFormat))
			case time.Duration:
				hs, _, err := haveVal.Duration()
				require.NoError(t, err)
				assert.Exactly(t, wv, hs)

			default:
				t.Fatalf("Unsupported type: %#v in vals index %d", wantVal, vi)
			}

			assert.Exactly(t, test.wantTypeIDs, cg.AllInvocations().ScopeIDs(), "Index %d/%d", vi, xi)
		}
	}
}

func TestScopedServicePermission_All(t *testing.T) {
	t.Parallel()
	basePath := config.MustMakePath("aa/bb/cc")

	sm := config.NewMock(config.MockPathValue{
		basePath.Bind(scope.DefaultTypeID).String(): "a",
		basePath.BindWebsite(1).String():            "b",
		basePath.BindStore(1).String():              "c",
	})

	tests := []struct {
		s       scope.Type
		want    string
		wantIDs scope.TypeIDs
	}{
		{scope.Default, "a", scope.TypeIDs{scope.DefaultTypeID}},
		{scope.Website, "b", scope.TypeIDs{scope.DefaultTypeID, scope.Website.Pack(1)}},
		{scope.Group, "a", scope.TypeIDs{scope.DefaultTypeID, scope.Website.Pack(1)}},
		{scope.Store, "c", scope.TypeIDs{scope.DefaultTypeID, scope.Website.Pack(1), scope.Store.Pack(1)}},
		{scope.Absent, "c", scope.TypeIDs{scope.DefaultTypeID, scope.Website.Pack(1), scope.Store.Pack(1)}}, // because ScopedGetter bound to store scope
	}
	for i, test := range tests {
		haveVal := sm.NewScoped(1, 1).Value(test.s, "aa/bb/cc")
		s, ok, err := haveVal.Str()
		if err != nil {
			t.Fatal("Index", i, "Error", err)
		}
		require.True(t, ok, "Scoped path value must be found")
		assert.Exactly(t, test.want, s, "Index %d", i)
		assert.Exactly(t, test.wantIDs, sm.Invokes().ScopeIDs(), "Index %d", i)
	}
	assert.Exactly(t, []string{"default/0/aa/bb/cc", "stores/1/aa/bb/cc", "websites/1/aa/bb/cc"}, sm.Invokes().Paths())
}

func TestScopedServicePermission_One(t *testing.T) {
	t.Parallel()
	basePath1 := config.MustMakePath("aa/bb/cc")
	basePath2 := config.MustMakePath("dd/ee/ff")
	basePath3 := config.MustMakePath("dd/ee/gg")

	const WebsiteID = 3
	const StoreID = 5

	sm := config.NewMock(config.MockPathValue{
		basePath1.Bind(scope.DefaultTypeID).String(): "a",
		basePath1.BindWebsite(WebsiteID).String():    "b",
		basePath1.BindStore(StoreID).String():        "c",

		basePath2.BindWebsite(WebsiteID).String(): "bb2",

		basePath3.String(): "cc3",
	})

	t.Run("query1 by scope.Default, matches default", func(t *testing.T) {
		s, ok, err := sm.NewScoped(WebsiteID, StoreID).Value(scope.Default, "aa/bb/cc").Str()
		require.True(t, ok, "scoped path value must be found")
		require.NoError(t, err)
		assert.Exactly(t, "a", s) // because ScopedGetter bound to store scope
	})

	t.Run("query1 by scope.Website, matches website", func(t *testing.T) {
		s, ok, err := sm.NewScoped(WebsiteID, StoreID).Value(scope.Website, "aa/bb/cc").Str()
		require.True(t, ok, "scoped path value must be found")
		require.NoError(t, err)
		assert.Exactly(t, "b", s) // because ScopedGetter bound to store scope
	})

	t.Run("query1 by scope.Store, matches store", func(t *testing.T) {
		s, ok, err := sm.NewScoped(WebsiteID, StoreID).Value(scope.Store, "aa/bb/cc").Str()
		require.True(t, ok, "scoped path value must be found")
		require.NoError(t, err)
		assert.Exactly(t, "c", s) // because ScopedGetter bound to store scope
	})
	t.Run("query1 by scope.Absent, matches store", func(t *testing.T) {
		s, ok, err := sm.NewScoped(WebsiteID, StoreID).Value(scope.Absent, "aa/bb/cc").Str()
		require.True(t, ok, "scoped path value must be found")
		require.NoError(t, err)
		assert.Exactly(t, "c", s) // because ScopedGetter bound to store scope
	})

	t.Run("query2 by scope.Store, fallback to website", func(t *testing.T) {
		s, ok, err := sm.NewScoped(WebsiteID, StoreID).Value(scope.Store, "dd/ee/ff").Str()
		require.True(t, ok, "scoped path value must be found")
		require.NoError(t, err)
		assert.Exactly(t, "bb2", s) // because ScopedGetter bound to store scope
	})

	t.Run("query3 by scope.Store, fallback to default", func(t *testing.T) {
		s, ok, err := sm.NewScoped(WebsiteID, StoreID).Value(scope.Store, "dd/ee/gg").Str()
		require.True(t, ok, "scoped path value must be found")
		require.NoError(t, err)
		assert.Exactly(t, "cc3", s) // because ScopedGetter bound to store scope
	})
	t.Run("query3 by scope.Website, fallback to default", func(t *testing.T) {
		s, ok, err := sm.NewScoped(WebsiteID, StoreID).Value(scope.Website, "dd/ee/gg").Str()
		require.True(t, ok, "scoped path value must be found")
		require.NoError(t, err)
		assert.Exactly(t, "cc3", s) // because ScopedGetter bound to store scope
	})

}

func TestScopedService_Parent(t *testing.T) {
	t.Parallel()
	tests := []struct {
		sg               config.Scoped
		wantCurrentScope scope.Type
		wantCurrentId    int64
		wantParentScope  scope.Type
		wantParentID     int64
	}{
		{config.NewScoped(nil, 33, 1), scope.Store, 1, scope.Website, 33},
		{config.NewScoped(nil, 3, 0), scope.Website, 3, scope.Default, 0},
		{config.NewScoped(nil, 0, 0), scope.Default, 0, scope.Default, 0},
	}
	for _, test := range tests {
		haveScp, haveID := test.sg.ParentID().Unpack()
		if have, want := haveScp, test.wantParentScope; have != want {
			t.Errorf("ParentScope: Have: %v Want: %v", have, want)
		}
		if have, want := haveID, test.wantParentID; have != want {
			t.Errorf("ParentScopeID: Have: %v Want: %v", have, want)
		}

		haveScp, haveID = test.sg.ScopeID().Unpack()
		if have, want := haveScp, test.wantCurrentScope; have != want {
			t.Errorf("Scope: Have: %v Want: %v", have, want)
		}
		if have, want := haveID, test.wantCurrentId; have != want {
			t.Errorf("ScopeID: Have: %v Want: %v", have, want)
		}

	}
}
