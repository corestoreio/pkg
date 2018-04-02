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
	"runtime"
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

func TestWithDBStorage(t *testing.T) {
	t.Skip("todo")
}

func TestNotKeyNotFoundError(t *testing.T) {

	srv := config.MustNewService(config.NewInMemoryStore())
	assert.NotNil(t, srv)

	scopedSrv := srv.NewScoped(1, 1)

	flat, ok, err := scopedSrv.Value(scope.Default, "catalog/product/enable_flat")
	require.NoError(t, err)
	assert.False(t, ok, "Should not find the key")
	assert.Empty(t, flat)
	//assert.Exactly(t, scope.DefaultTypeID.String(), h.String())

	val, ok, err := scopedSrv.Value(scope.Store, "catalog")
	assert.Empty(t, val)
	assert.True(t, errors.NotValid.Match(err), "Error: %s", err)
	assert.False(t, errors.NotFound.Match(err), "Error: %s", err)
	//assert.Exactly(t, scope.TypeID(0).String(), h.String())
}

func TestService_Write(t *testing.T) {

	srv := config.MustNewService(config.NewInMemoryStore())
	assert.NotNil(t, srv)

	p1 := config.Path{}
	err := srv.Write(p1, []byte{})
	assert.True(t, errors.NotValid.Match(err), "Error: %s", err)
}

func TestService_Write_Get_Value_Success(t *testing.T) {

	runner := func(p config.Path, value []byte) func(*testing.T) {
		return func(t *testing.T) {
			srv := config.MustNewService(config.NewInMemoryStore())

			require.NoError(t, srv.Write(p, value), "Writing Value in Test %q should not fail", t.Name())

			srvVal, haveOK, haveErr := srv.Value(p)
			require.NoError(t, haveErr, "No error should occur when retrieving a value")
			require.True(t, haveOK, "The value should not be nil")

			haveStr, ok, haveErr := srvVal.Str()
			require.NoError(t, haveErr)
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

	for vi, wantVal := range vals {
		for xi, test := range tests {

			cg := config.NewMock(config.MockPathValue{
				test.fqpath: conv.ToString(wantVal),
			})

			sg := cg.NewScoped(test.websiteID, test.storeID)
			haveVal, haveOK, haveErr := sg.Value(test.perm, test.route)

			if test.wantErrKind > 0 {
				require.True(t, haveOK, "Index %d/%d scoped path value must be found ", vi, xi)
				// if d, ok := haveVal.(time.Duration); ok {
				// 	// oh that is so crap because time.Duration cannot be detected for zero value
				// 	assert.Empty(t, int64(d), "Index %d/%d => %v", vi, xi, wantVal)
				// } else {
				// 	assert.Empty(t, haveVal, "Index %d/%d => %v", vi, xi, wantVal)
				// }

				assert.True(t, test.wantErrKind.Match(haveErr), "Error: %s => %s", haveErr, test.desc)
				continue
			}
			require.NoError(t, haveErr, "Error: %+v\n\n%s", haveErr, test.desc)

			switch wantVal.(type) {
			case []byte:
				var buf strings.Builder
				_, err := haveVal.WriteTo(&buf)
				require.NoError(t, err, "Error: %+v\n\n%s", err, test.desc)
				assert.Exactly(t, wantVal, buf.String(), test.desc)

			case string:

			case bool:

			case float64:

			case int:

			case time.Time:

			case time.Duration:

			default:
				t.Fatalf("Unsupported type: %#v in vals index %d", wantVal, vi)
			}

			assert.Exactly(t, test.wantTypeIDs, cg.AllInvocations().ScopeIDs(), "Index %d/%d", vi, xi)
		}
	}
}

var benchmarkScopedServiceVal config.Value

// BenchmarkScopedServiceStringStore-4	 1000000	      2218 ns/op	     320 B/op	       9 allocs/op => Go 1.5.2
// BenchmarkScopedServiceStringStore-4	  500000	      2939 ns/op	     672 B/op	      17 allocs/op => Go 1.5.3 strings
// BenchmarkScopedServiceStringStore-4    500000	      2732 ns/op	     912 B/op	      17 allocs/op => cfgpath.Path with []ArgFunc
// BenchmarkScopedServiceStringStore-4	 1000000	      1821 ns/op	     336 B/op	       3 allocs/op => cfgpath.Path without []ArgFunc
// BenchmarkScopedServiceStringStore-4   1000000	      1747 ns/op	       0 B/op	       0 allocs/op => Go 1.6 sync.Pool cfgpath.Path without []ArgFunc
// BenchmarkScopedServiceStringStore-4    500000	      2604 ns/op	       0 B/op	       0 allocs/op => Go 1.7 with ScopeHash
func BenchmarkScopedServiceStringStore(b *testing.B) {
	benchmarkScopedServiceStringRun(b, 1, 1)
}

func BenchmarkScopedServiceStringWebsite(b *testing.B) {
	benchmarkScopedServiceStringRun(b, 1, 0)
}

func BenchmarkScopedServiceStringDefault(b *testing.B) {
	benchmarkScopedServiceStringRun(b, 0, 0)
}

func benchmarkScopedServiceStringRun(b *testing.B, websiteID, storeID int64) {
	route := "aa/bb/cc"
	want := strings.Repeat("Gopher", 100)
	sg := config.NewMock(config.MockPathValue{
		config.MustMakePath(route).String(): want,
	}).NewScoped(websiteID, storeID)

	runtime.GC()
	b.ResetTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		var err error
		var ok bool
		benchmarkScopedServiceVal, ok, err = sg.Value(scope.Store, route)
		if err != nil {
			b.Fatal(err)
		}
		if !ok {
			b.Fatal("must be ok")
		}
		s, _, _ := benchmarkScopedServiceVal.Str()
		if s != want {
			b.Errorf("Want %s Have %s", want, benchmarkScopedServiceVal)
		}
	}
}

func TestScopedServicePermission(t *testing.T) {

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
		have, ok, err := sm.NewScoped(1, 1).Value(test.s, "aa/bb/cc")
		if err != nil {
			t.Fatal("Index", i, "Error", err)
		}
		require.True(t, ok, "Scoped path value must be found")
		s, _, _ := have.Str()
		assert.Exactly(t, test.want, s, "Index %d", i)
		assert.Exactly(t, test.wantIDs, sm.Invokes().ScopeIDs(), "Index %d", i)
	}

	have, ok, err := sm.NewScoped(1, 1).Value(scope.Default, "aa/bb/cc")
	require.True(t, ok, "scoped path value must be found")
	require.NoError(t, err)
	s, _, _ := have.Str()
	assert.Exactly(t, "c", s) // because ScopedGetter bound to store scope
	assert.Exactly(t, []string{"default/0/aa/bb/cc", "stores/1/aa/bb/cc", "websites/1/aa/bb/cc"}, sm.Invokes().Paths())
}

func TestScopedService_Parent(t *testing.T) {
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
