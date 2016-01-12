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
	"errors"
	"strings"
	"testing"
	"time"

	"github.com/corestoreio/csfw/config"
	"github.com/corestoreio/csfw/config/path"
	"github.com/corestoreio/csfw/store/scope"
	"github.com/stretchr/testify/assert"
)

func TestScopedServiceScope(t *testing.T) {
	tests := []struct {
		websiteID, groupID, storeID int64
		wantScope                   scope.Scope
		wantID                      int64
	}{
		{0, 0, 0, scope.DefaultID, 0},
		{1, 0, 0, scope.WebsiteID, 1},
		{1, 2, 0, scope.GroupID, 2},
		{1, 2, 3, scope.StoreID, 3},
		{0, 0, 3, scope.StoreID, 3},
		{0, 2, 0, scope.GroupID, 2},
	}
	for i, test := range tests {
		sg := config.NewMockGetter().NewScoped(test.websiteID, test.groupID, test.storeID)
		haveScope, haveID := sg.Scope()
		assert.Exactly(t, test.wantScope, haveScope, "Index %d", i)
		assert.Exactly(t, test.wantID, haveID, "Index %d", i)
	}
}

func TestScopedServicePath(t *testing.T) {
	tests := []struct {
		desc                        string
		fqpath                      string
		path                        []string
		websiteID, groupID, storeID int64
		err                         error
	}{
		{
			"Default ScopedGetter should return default scope",
			path.MustNew("aa/bb/cc").String(), []string{"aa/bb/cc"}, 0, 0, 0, nil,
		},
		{
			"Website ID 1 ScopedGetter should fall back to default scope",
			path.MustNew("aa/bb/cc").String(), []string{"aa/bb/cc"}, 1, 0, 0, nil,
		},
		{
			"Website ID 10 + Group ID 12 ScopedGetter should fall back to website 10 scope",
			path.MustNew("aa/bb/cc").Bind(scope.WebsiteID, 10).String(), []string{"aa/bb/cc"}, 10, 12, 0, nil,
		},
		{
			"Website ID 10 + Group ID 12 + Store 22 ScopedGetter should fall back to website 10 scope",
			path.MustNew("aa/bb/cc").Bind(scope.WebsiteID, 10).String(), []string{"aa/bb/cc"}, 10, 12, 22, nil,
		},
		{
			"Website ID 10 + Group ID 12 + Store 22 ScopedGetter should return Store 22 scope",
			path.MustNew("aa/bb/cc").Bind(scope.StoreID, 22).String(), []string{"aa/bb/cc"}, 10, 12, 22, nil,
		},
		{
			"Website ID 10 + Group ID 12 + Store 42 ScopedGetter should return nothing",
			path.MustNew("aa/bb/cc").Bind(scope.StoreID, 22).String(), []string{"aa/bb/cc"}, 10, 12, 42, config.ErrKeyNotFound,
		},
		{
			"Path consists of only two elements which is incorrect",
			path.MustNew("aa/bb/cc").String(), []string{"aa", "bb"}, 0, 0, 0, errors.New("Incorrect number of paths elements: want 3, have 2, Path: [aa bb]"),
		},
	}

	// vals stores all possible types for which we have functions in config.ScopedGetter
	vals := []interface{}{"Gopher", true, float64(3.14159), int(2016), time.Now()}

	for vi, wantVal := range vals {
		for _, test := range tests {

			cg := config.NewMockGetter(config.WithMockValues(config.MockPV{
				test.fqpath: wantVal,
			}))

			sg := cg.NewScoped(test.websiteID, test.groupID, test.storeID)

			var haveVal interface{}
			var haveErr error
			switch wantVal.(type) {
			case string:
				haveVal, haveErr = sg.String(test.path...)
			case bool:
				haveVal, haveErr = sg.Bool(test.path...)
			case float64:
				haveVal, haveErr = sg.Float64(test.path...)
			case int:
				haveVal, haveErr = sg.Int(test.path...)
			case time.Time:
				haveVal, haveErr = sg.DateTime(test.path...)
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

// BenchmarkScopedServiceStringStore-4	 1000000	      2218 ns/op	     320 B/op	       9 allocs/op => Go 1.5.2
// BenchmarkScopedServiceStringStore-4	  500000	      2939 ns/op	     672 B/op	      17 allocs/op
func BenchmarkScopedServiceStringStore(b *testing.B) {

	var want = strings.Repeat("Gopher", 100)
	sg := config.NewMockGetter(config.WithMockValues(config.MockPV{
		path.MustNew("aa/bb/cc").String(): want,
	})).NewScoped(1, 1, 1)

	b.ResetTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		v, err := sg.String("aa/bb/cc")
		if err != nil {
			b.Error(err)
		}
		if v != want {
			b.Errorf("Want %s Have %s", want, v)
		}
	}
}
