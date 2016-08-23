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

package config_test

import (
	"runtime"
	"strings"
	"testing"
	"time"

	"github.com/corestoreio/csfw/config"
	"github.com/corestoreio/csfw/config/cfgmock"
	"github.com/corestoreio/csfw/config/cfgpath"
	"github.com/corestoreio/csfw/store/scope"
	"github.com/corestoreio/csfw/util/errors"
	"github.com/stretchr/testify/assert"
)

func TestScoped_IsValid(t *testing.T) {
	cfg := cfgmock.NewService()
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
		wantScope          scope.Scope
		wantID             int64
	}{
		{0, 0, scope.Default, 0},
		{1, 0, scope.Website, 1},
		{1, 3, scope.Store, 3},
		{0, 3, scope.Store, 3},
	}
	for i, test := range tests {
		sg := cfgmock.NewService().NewScoped(test.websiteID, test.storeID)
		haveScope, haveID := sg.Scope()
		assert.Exactly(t, test.wantScope, haveScope, "Index %d", i)
		assert.Exactly(t, test.wantID, haveID, "Index %d", i)
	}
}

func TestScopedServicePath(t *testing.T) {

	basePath := cfgpath.MustNewByParts("aa/bb/cc")
	tests := []struct {
		desc               string
		fqpath             string
		route              cfgpath.Route
		perm               scope.Scope
		websiteID, storeID int64
		wantErrBhf         errors.BehaviourFunc
		wantHash           scope.Hash
	}{
		{
			"Default ScopedGetter should return default scope",
			basePath.String(), cfgpath.NewRoute("aa/bb/cc"), scope.Absent, 0, 0, nil, scope.DefaultHash,
		},
		{
			"Website ID 1 ScopedGetter should fall back to default scope",
			basePath.String(), cfgpath.NewRoute("aa/bb/cc"), scope.Website, 1, 0, nil, scope.DefaultHash,
		},
		{
			"Website ID 10 ScopedGetter should fall back to website 10 scope",
			basePath.BindWebsite(10).String(), cfgpath.NewRoute("aa/bb/cc"), scope.Website, 10, 0, nil, scope.NewHash(scope.Website, 10),
		},
		{
			"Website ID 10 + Store 22 ScopedGetter should fall back to website 10 scope",
			basePath.BindWebsite(10).String(), cfgpath.NewRoute("aa/bb/cc"), scope.Store, 10, 22, nil, scope.NewHash(scope.Website, 10),
		},
		{
			"Website ID 10 + Store 22 ScopedGetter should return Store 22 scope",
			basePath.BindStore(22).String(), cfgpath.NewRoute("aa/bb/cc"), scope.Store, 10, 22, nil, scope.NewHash(scope.Store, 22),
		},
		{
			"Website ID 10 + Store 42 ScopedGetter should return nothing",
			basePath.BindStore(22).String(), cfgpath.NewRoute("aa/bb/cc"), scope.Store, 10, 42, errors.IsNotFound, scope.DefaultHash,
		},
		{
			"Path consists of only two elements which is incorrect",
			basePath.String(), cfgpath.NewRoute("aa", "bb"), scope.Store, 0, 0, errors.IsNotValid, 0,
		},
	}

	// vals stores all possible types for which we have functions in config.ScopedGetter
	vals := []interface{}{"Gopher", true, float64(3.14159), int(2016), time.Now(), []byte(`Hellö Dear Goph€rs`)}

	for vi, wantVal := range vals {
		for _, test := range tests {

			cg := cfgmock.NewService(cfgmock.PathValue{
				test.fqpath: wantVal,
			})

			sg := cg.NewScoped(test.websiteID, test.storeID)

			var haveVal interface{}
			var haveErr error
			var haveHash scope.Hash
			switch wantVal.(type) {
			case []byte:
				haveVal, haveHash, haveErr = sg.Byte(test.route, test.perm)
			case string:
				haveVal, haveHash, haveErr = sg.String(test.route, test.perm)
			case bool:
				haveVal, haveHash, haveErr = sg.Bool(test.route, test.perm)
			case float64:
				haveVal, haveHash, haveErr = sg.Float64(test.route, test.perm)
			case int:
				haveVal, haveHash, haveErr = sg.Int(test.route, test.perm)
			case time.Time:
				haveVal, haveHash, haveErr = sg.Time(test.route, test.perm)
			default:
				t.Fatalf("Unsupported type: %#v in vals index %d", wantVal, vi)
			}
			testScopedService(t, wantVal, haveVal, test.wantHash, haveHash, test.desc, test.wantErrBhf, haveErr)
		}
	}
}

func testScopedService(t *testing.T, want, have interface{}, wantHash, haveHash scope.Hash, desc string, wantErrBhf errors.BehaviourFunc, err error) {
	assert.Exactly(t, wantHash.String(), haveHash.String(), desc)
	if wantErrBhf != nil {
		assert.Empty(t, have, desc)
		assert.True(t, wantErrBhf(err), "Error: %s => %s", err, desc)
		return
	}
	assert.NoError(t, err, "Error: %+v\n\n%s", err, desc)
	assert.Exactly(t, want, have, desc)
}

var benchmarkScopedServiceString string

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
	route := cfgpath.NewRoute("aa/bb/cc")
	want := strings.Repeat("Gopher", 100)
	sg := cfgmock.NewService(cfgmock.PathValue{
		cfgpath.MustNew(route).String(): want,
	}).NewScoped(websiteID, storeID)

	runtime.GC()
	b.ResetTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		var err error
		var h scope.Hash
		benchmarkScopedServiceString, h, err = sg.String(route)
		if err != nil {
			b.Error(err)
		}
		if benchmarkScopedServiceString != want {
			b.Errorf("Want %s Have %s", want, benchmarkScopedServiceString)
		}
		if h != scope.DefaultHash {
			b.Errorf("Want %s Have %s", scope.DefaultHash, h)
		}
	}
}

func TestScopedServicePermission(t *testing.T) {

	basePath := cfgpath.MustNewByParts("aa/bb/cc")

	sg := cfgmock.NewService(cfgmock.PathValue{
		basePath.Bind(scope.Default, 0).String(): "a",
		basePath.BindWebsite(1).String():         "b",
		basePath.BindStore(1).String():           "c",
	}).NewScoped(1, 1)

	tests := []struct {
		s        scope.Scope
		want     string
		wantHash scope.Hash
	}{
		{scope.Default, "a", scope.DefaultHash},
		{scope.Website, "b", scope.NewHash(scope.Website, 1)},
		{scope.Group, "a", scope.DefaultHash},
		{scope.Store, "c", scope.NewHash(scope.Store, 1)},
		{scope.Absent, "c", scope.NewHash(scope.Store, 1)}, // because ScopedGetter bound to store scope
	}
	for i, test := range tests {
		have, haveH, err := sg.String(basePath.Route, test.s)
		if err != nil {
			t.Fatal("Index", i, "Error", err)
		}
		assert.Exactly(t, test.want, have, "Index %d", i)
		assert.Exactly(t, test.wantHash.String(), haveH.String(), "Index %d", i)
	}

	var ss = []scope.Scope{}
	ss = nil
	have, haveH, err := sg.String(basePath.Route, ss...)
	assert.NoError(t, err)
	assert.Exactly(t, "c", have) // because ScopedGetter bound to store scope
	assert.Exactly(t, scope.NewHash(scope.Store, 1).String(), haveH.String())
}

func TestScopedService_Parent(t *testing.T) {
	tests := []struct {
		sg               config.Scoped
		wantCurrentScope scope.Scope
		wantCurrentId    int64
		wantParentScope  scope.Scope
		wantParentID     int64
	}{
		{config.NewScoped(nil, 33, 1), scope.Store, 1, scope.Website, 33},
		{config.NewScoped(nil, 3, 0), scope.Website, 3, scope.Default, 0},
		{config.NewScoped(nil, 0, 0), scope.Default, 0, scope.Default, 0},
	}
	for _, test := range tests {
		haveScp, haveID := test.sg.Parent()
		if have, want := haveScp, test.wantParentScope; have != want {
			t.Errorf("ParentScope: Have: %v Want: %v", have, want)
		}
		if have, want := haveID, test.wantParentID; have != want {
			t.Errorf("ParentScopeID: Have: %v Want: %v", have, want)
		}

		haveScp, haveID = test.sg.Scope()
		if have, want := haveScp, test.wantCurrentScope; have != want {
			t.Errorf("Scope: Have: %v Want: %v", have, want)
		}
		if have, want := haveID, test.wantCurrentId; have != want {
			t.Errorf("ScopeID: Have: %v Want: %v", have, want)
		}

	}
}
