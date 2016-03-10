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
	"github.com/corestoreio/csfw/config/path"
	"github.com/corestoreio/csfw/config/storage"
	"github.com/corestoreio/csfw/store/scope"
	"github.com/corestoreio/csfw/util/log"
	"github.com/stretchr/testify/assert"
)

func TestScopedServiceScope(t *testing.T) {
	t.Parallel()
	tests := []struct {
		websiteID, storeID int64
		wantScope          scope.Scope
		wantID             int64
	}{
		{0, 0, scope.DefaultID, 0},
		{1, 0, scope.WebsiteID, 1},
		{1, 3, scope.StoreID, 3},
		{0, 3, scope.StoreID, 3},
	}
	for i, test := range tests {
		sg := cfgmock.NewService().NewScoped(test.websiteID, test.storeID)
		haveScope, haveID := sg.Scope()
		assert.Exactly(t, test.wantScope, haveScope, "Index %d", i)
		assert.Exactly(t, test.wantID, haveID, "Index %d", i)
	}
}

func TestScopedServicePath(t *testing.T) {
	t.Parallel()
	basePath := path.MustNewByParts("aa/bb/cc")
	tests := []struct {
		desc               string
		fqpath             string
		route              path.Route
		perm               scope.Scope
		websiteID, storeID int64
		err                error
	}{
		{
			"Default ScopedGetter should return default scope",
			basePath.String(), path.NewRoute("aa/bb/cc"), scope.AbsentID, 0, 0, nil,
		},
		{
			"Website ID 1 ScopedGetter should fall back to default scope",
			basePath.String(), path.NewRoute("aa/bb/cc"), scope.WebsiteID, 1, 0, nil,
		},
		{
			"Website ID 10 ScopedGetter should fall back to website 10 scope",
			basePath.Bind(scope.WebsiteID, 10).String(), path.NewRoute("aa/bb/cc"), scope.WebsiteID, 10, 0, nil,
		},
		{
			"Website ID 10 + Store 22 ScopedGetter should fall back to website 10 scope",
			basePath.Bind(scope.WebsiteID, 10).String(), path.NewRoute("aa/bb/cc"), scope.StoreID, 10, 22, nil,
		},
		{
			"Website ID 10 + Store 22 ScopedGetter should return Store 22 scope",
			basePath.Bind(scope.StoreID, 22).String(), path.NewRoute("aa/bb/cc"), scope.StoreID, 10, 22, nil,
		},
		{
			"Website ID 10 + Store 42 ScopedGetter should return nothing",
			basePath.Bind(scope.StoreID, 22).String(), path.NewRoute("aa/bb/cc"), scope.StoreID, 10, 42, storage.ErrKeyNotFound,
		},
		{
			"Path consists of only two elements which is incorrect",
			basePath.String(), path.NewRoute("aa", "bb"), scope.StoreID, 0, 0, path.ErrIncorrectPath,
		},
	}

	// vals stores all possible types for which we have functions in config.ScopedGetter
	//vals := []interface{}{"Gopher", true, float64(3.14159), int(2016), time.Now()}
	vals := []interface{}{"Gopher"}

	for vi, wantVal := range vals {
		for _, test := range tests {

			cg := cfgmock.NewService(cfgmock.WithPV(cfgmock.PathValue{
				test.fqpath: wantVal,
			}))

			sg := cg.NewScoped(test.websiteID, test.storeID)

			var haveVal interface{}
			var haveErr error
			switch wantVal.(type) {
			case string:
				haveVal, haveErr = sg.String(test.route, test.perm)
			case bool:
				haveVal, haveErr = sg.Bool(test.route, test.perm)
			case float64:
				haveVal, haveErr = sg.Float64(test.route, test.perm)
			case int:
				haveVal, haveErr = sg.Int(test.route, test.perm)
			case time.Time:
				haveVal, haveErr = sg.Time(test.route, test.perm)
			default:
				t.Fatalf("Unsupported type: %#v in vals index %d", wantVal, vi)
			}
			testScopedService(t, wantVal, haveVal, test.desc, test.err, haveErr)
		}
	}
}

func testScopedService(t *testing.T, want, have interface{}, desc string, wantErr, err error) {
	if wantErr != nil {
		// if this fails for time.Time{} then my PR to assert pkg has not yet been merged :-(
		// https://github.com/stretchr/testify/pull/259
		assert.Empty(t, have, desc)
		assert.EqualError(t, err, wantErr.Error(), desc)
		return
	}
	assert.NoError(t, err, desc)
	assert.Exactly(t, want, have, desc)
}

var benchmarkScopedServiceString string

// BenchmarkScopedServiceStringStore-4	 1000000	      2218 ns/op	     320 B/op	       9 allocs/op => Go 1.5.2
// BenchmarkScopedServiceStringStore-4	  500000	      2939 ns/op	     672 B/op	      17 allocs/op => Go 1.5.3 strings
// BenchmarkScopedServiceStringStore-4    500000	      2732 ns/op	     912 B/op	      17 allocs/op => path.Path with []ArgFunc
// BenchmarkScopedServiceStringStore-4	 1000000	      1821 ns/op	     336 B/op	       3 allocs/op => path.Path without []ArgFunc
// BenchmarkScopedServiceStringStore-4   1000000	      1747 ns/op	       0 B/op	       0 allocs/op => Go 1.6 sync.Pool path.Path without []ArgFunc
func BenchmarkScopedServiceStringStore(b *testing.B) {
	benchmarkScopedServiceStringRun(b, 1, 1, scope.StoreID)
}

func BenchmarkScopedServiceStringWebsite(b *testing.B) {
	benchmarkScopedServiceStringRun(b, 1, 0, scope.WebsiteID)
}

func BenchmarkScopedServiceStringDefault(b *testing.B) {
	benchmarkScopedServiceStringRun(b, 0, 0, scope.DefaultID)
}

func benchmarkScopedServiceStringRun(b *testing.B, websiteID, storeID int64, s scope.Scope) {
	config.PkgLog.SetLevel(log.StdLevelFatal)
	route := path.NewRoute("aa/bb/cc")
	want := strings.Repeat("Gopher", 100)
	sg := cfgmock.NewService(cfgmock.WithPV(cfgmock.PathValue{
		path.MustNew(route).String(): want,
	})).NewScoped(websiteID, storeID)

	runtime.GC()
	b.ResetTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		var err error
		benchmarkScopedServiceString, err = sg.String(route)
		if err != nil {
			b.Error(err)
		}
		if benchmarkScopedServiceString != want {
			b.Errorf("Want %s Have %s", want, benchmarkScopedServiceString)
		}
	}
}

func TestScopedServicePermission(t *testing.T) {
	t.Parallel()

	basePath := path.MustNewByParts("aa/bb/cc")

	sg := cfgmock.NewService(cfgmock.WithPV(cfgmock.PathValue{
		basePath.Bind(scope.DefaultID, 0).String(): "a",
		basePath.Bind(scope.WebsiteID, 1).String(): "b",
		basePath.Bind(scope.StoreID, 1).String():   "c",
	})).NewScoped(1, 1)

	tests := []struct {
		s    scope.Scope
		want string
	}{
		{scope.DefaultID, "a"},
		{scope.WebsiteID, "b"},
		{scope.GroupID, "a"},
		{scope.StoreID, "c"},
		{scope.AbsentID, "a"},
	}
	for i, test := range tests {
		have, err := sg.String(basePath.Route, test.s)
		if err != nil {
			t.Fatal("Index", i, "Error", err)
		}
		assert.Exactly(t, test.want, have, "Index %d", i)
	}

	var ss = []scope.Scope{}
	ss = nil
	have, err := sg.String(basePath.Route, ss...)
	assert.NoError(t, err)
	assert.Exactly(t, "a", have)
}
