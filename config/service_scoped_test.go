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

	"github.com/corestoreio/csfw/config/cfgmock"
	"github.com/corestoreio/csfw/config/cfgpath"
	"github.com/corestoreio/csfw/store/scope"
	"github.com/corestoreio/csfw/util/errors"
	"github.com/stretchr/testify/assert"
)

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
	}{
		{
			"Default ScopedGetter should return default scope",
			basePath.String(), cfgpath.NewRoute("aa/bb/cc"), scope.Absent, 0, 0, nil,
		},
		{
			"Website ID 1 ScopedGetter should fall back to default scope",
			basePath.String(), cfgpath.NewRoute("aa/bb/cc"), scope.Website, 1, 0, nil,
		},
		{
			"Website ID 10 ScopedGetter should fall back to website 10 scope",
			basePath.Bind(scope.Website, 10).String(), cfgpath.NewRoute("aa/bb/cc"), scope.Website, 10, 0, nil,
		},
		{
			"Website ID 10 + Store 22 ScopedGetter should fall back to website 10 scope",
			basePath.Bind(scope.Website, 10).String(), cfgpath.NewRoute("aa/bb/cc"), scope.Store, 10, 22, nil,
		},
		{
			"Website ID 10 + Store 22 ScopedGetter should return Store 22 scope",
			basePath.Bind(scope.Store, 22).String(), cfgpath.NewRoute("aa/bb/cc"), scope.Store, 10, 22, nil,
		},
		{
			"Website ID 10 + Store 42 ScopedGetter should return nothing",
			basePath.Bind(scope.Store, 22).String(), cfgpath.NewRoute("aa/bb/cc"), scope.Store, 10, 42, errors.IsNotFound,
		},
		{
			"Path consists of only two elements which is incorrect",
			basePath.String(), cfgpath.NewRoute("aa", "bb"), scope.Store, 0, 0, errors.IsNotValid,
		},
	}

	// vals stores all possible types for which we have functions in config.ScopedGetter
	vals := []interface{}{"Gopher", true, float64(3.14159), int(2016), time.Now(), []byte(`Hellö Dear Goph€rs`)}

	for vi, wantVal := range vals {
		for _, test := range tests {

			cg := cfgmock.NewService(cfgmock.WithPV(cfgmock.PathValue{
				test.fqpath: wantVal,
			}))

			sg := cg.NewScoped(test.websiteID, test.storeID)

			var haveVal interface{}
			var haveErr error
			switch wantVal.(type) {
			case []byte:
				haveVal, haveErr = sg.Byte(test.route, test.perm)
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
			testScopedService(t, wantVal, haveVal, test.desc, test.wantErrBhf, haveErr)
		}
	}
}

func testScopedService(t *testing.T, want, have interface{}, desc string, wantErrBhf errors.BehaviourFunc, err error) {
	if wantErrBhf != nil {
		assert.Empty(t, have, desc)
		assert.True(t, wantErrBhf(err), "Error: %s => %s", err, desc)
		return
	}
	assert.NoError(t, err, desc)
	assert.Exactly(t, want, have, desc)
}

var benchmarkScopedServiceString string

// BenchmarkScopedServiceStringStore-4	 1000000	      2218 ns/op	     320 B/op	       9 allocs/op => Go 1.5.2
// BenchmarkScopedServiceStringStore-4	  500000	      2939 ns/op	     672 B/op	      17 allocs/op => Go 1.5.3 strings
// BenchmarkScopedServiceStringStore-4    500000	      2732 ns/op	     912 B/op	      17 allocs/op => cfgpath.Path with []ArgFunc
// BenchmarkScopedServiceStringStore-4	 1000000	      1821 ns/op	     336 B/op	       3 allocs/op => cfgpath.Path without []ArgFunc
// BenchmarkScopedServiceStringStore-4   1000000	      1747 ns/op	       0 B/op	       0 allocs/op => Go 1.6 sync.Pool cfgpath.Path without []ArgFunc
func BenchmarkScopedServiceStringStore(b *testing.B) {
	benchmarkScopedServiceStringRun(b, 1, 1, scope.Store)
}

func BenchmarkScopedServiceStringWebsite(b *testing.B) {
	benchmarkScopedServiceStringRun(b, 1, 0, scope.Website)
}

func BenchmarkScopedServiceStringDefault(b *testing.B) {
	benchmarkScopedServiceStringRun(b, 0, 0, scope.Default)
}

func benchmarkScopedServiceStringRun(b *testing.B, websiteID, storeID int64, s scope.Scope) {
	route := cfgpath.NewRoute("aa/bb/cc")
	want := strings.Repeat("Gopher", 100)
	sg := cfgmock.NewService(cfgmock.WithPV(cfgmock.PathValue{
		cfgpath.MustNew(route).String(): want,
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

	basePath := cfgpath.MustNewByParts("aa/bb/cc")

	sg := cfgmock.NewService(cfgmock.WithPV(cfgmock.PathValue{
		basePath.Bind(scope.Default, 0).String(): "a",
		basePath.Bind(scope.Website, 1).String(): "b",
		basePath.Bind(scope.Store, 1).String():   "c",
	})).NewScoped(1, 1)

	tests := []struct {
		s    scope.Scope
		want string
	}{
		{scope.Default, "a"},
		{scope.Website, "b"},
		{scope.Group, "a"},
		{scope.Store, "c"},
		{scope.Absent, "c"}, // because ScopedGetter bound to store scope
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
	assert.Exactly(t, "c", have) // because ScopedGetter bound to store scope
}
