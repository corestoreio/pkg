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
	"testing"

	"sort"

	"github.com/corestoreio/csfw/config/cfgpath"
	"github.com/corestoreio/csfw/store/scope"
	"github.com/stretchr/testify/assert"
)

var _ sort.Interface = (*cfgpath.PathSlice)(nil)

func TestPathSlice_Contains(t *testing.T) {
	t.Parallel()
	tests := []struct {
		paths  cfgpath.PathSlice
		search cfgpath.Path
		want   bool
	}{
		{
			cfgpath.PathSlice{
				0: cfgpath.MustNewByParts("aa/bb/cc").Bind(scope.WebsiteID, 3),
				1: cfgpath.MustNewByParts("aa/bb/cc").Bind(scope.WebsiteID, 2),
			},
			cfgpath.MustNewByParts("aa/bb/cc").Bind(scope.WebsiteID, 2),
			true,
		},
		{
			cfgpath.PathSlice{
				0: cfgpath.MustNewByParts("aa/bb/cc").Bind(scope.WebsiteID, 3),
				1: cfgpath.MustNewByParts("aa/bb/cc").Bind(scope.WebsiteID, 2),
			},
			cfgpath.MustNewByParts("aa/bb/cc").Bind(scope.StoreID, 2),
			false,
		},
	}
	for i, test := range tests {
		assert.Exactly(t, test.want, test.paths.Contains(test.search), "Index %d", i)
	}
}

func TestPathSlice_Sort(t *testing.T) {
	t.Parallel()
	ps := cfgpath.PathSlice{
		cfgpath.MustNewByParts("bb/cc/dd"),
		cfgpath.MustNewByParts("xx/yy/zz"),
		cfgpath.MustNewByParts("aa/bb/cc"),
	}
	ps.Sort()
	want := cfgpath.PathSlice{cfgpath.Path{Route: cfgpath.NewRoute(`aa/bb/cc`), Scope: 1, ID: 0}, cfgpath.Path{Route: cfgpath.NewRoute(`bb/cc/dd`), Scope: 1, ID: 0}, cfgpath.Path{Route: cfgpath.NewRoute(`xx/yy/zz`), Scope: 1, ID: 0}}
	assert.Exactly(t, want, ps)
}

// BenchmarkPathSlice_Sort-4	 1000000	      1987 ns/op	     480 B/op	       8 allocs/op
func BenchmarkPathSlice_Sort(b *testing.B) {
	// allocs are here uninteresting
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		ps := cfgpath.PathSlice{
			cfgpath.MustNewByParts("rr/ss/tt"),
			cfgpath.MustNewByParts("bb/cc/dd"),
			cfgpath.MustNewByParts("xx/yy/zz"),
			cfgpath.MustNewByParts("aa/bb/cc"),
			cfgpath.MustNewByParts("ff/gg/hh"),
			cfgpath.MustNewByParts("cc/dd/ee"),
		}
		ps.Sort()
		if len(ps) != 6 {
			b.Fatal("Incorrect length of ps variable after sorting")
		}
	}

}
