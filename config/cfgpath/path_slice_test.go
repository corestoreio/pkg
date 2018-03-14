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

	"github.com/corestoreio/pkg/config/cfgpath"
	"github.com/corestoreio/pkg/store/scope"
	"github.com/stretchr/testify/assert"
)

var _ sort.Interface = (*cfgpath.PathSlice)(nil)

func TestPathSlice_Contains(t *testing.T) {

	tests := []struct {
		paths  cfgpath.PathSlice
		search cfgpath.Path
		want   bool
	}{
		{
			cfgpath.PathSlice{
				0: cfgpath.MustMakeByString("aa/bb/cc").BindWebsite(3),
				1: cfgpath.MustMakeByString("aa/bb/cc").BindWebsite(2),
			},
			cfgpath.MustMakeByString("aa/bb/cc").BindWebsite(2),
			true,
		},
		{
			cfgpath.PathSlice{
				0: cfgpath.MustMakeByString("aa/bb/cc").BindWebsite(3),
				1: cfgpath.MustMakeByString("aa/bb/cc").BindWebsite(2),
			},
			cfgpath.MustMakeByString("aa/bb/cc").BindStore(2),
			false,
		},
	}
	for i, test := range tests {
		assert.Exactly(t, test.want, test.paths.Contains(test.search), "Index %d", i)
	}
}

func TestPathSlice_Sort(t *testing.T) {

	ps := cfgpath.PathSlice{
		cfgpath.MustMakeByString("bb/cc/dd"),
		cfgpath.MustMakeByString("xx/yy/zz"),
		cfgpath.MustMakeByString("aa/bb/cc"),
	}
	ps.Sort()
	want := cfgpath.PathSlice{
		cfgpath.Path{Route: cfgpath.MakeRoute(`aa/bb/cc`), ScopeID: scope.DefaultTypeID},
		cfgpath.Path{Route: cfgpath.MakeRoute(`bb/cc/dd`), ScopeID: scope.DefaultTypeID},
		cfgpath.Path{Route: cfgpath.MakeRoute(`xx/yy/zz`), ScopeID: scope.DefaultTypeID},
	}
	assert.Exactly(t, want, ps)
}

// BenchmarkPathSlice_Sort-4	 1000000	      1987 ns/op	     480 B/op	       8 allocs/op
func BenchmarkPathSlice_Sort(b *testing.B) {
	// allocs are here uninteresting
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		ps := cfgpath.PathSlice{
			cfgpath.MustMakeByString("rr/ss/tt"),
			cfgpath.MustMakeByString("bb/cc/dd"),
			cfgpath.MustMakeByString("xx/yy/zz"),
			cfgpath.MustMakeByString("aa/bb/cc"),
			cfgpath.MustMakeByString("ff/gg/hh"),
			cfgpath.MustMakeByString("cc/dd/ee"),
		}
		ps.Sort()
		if len(ps) != 6 {
			b.Fatal("Incorrect length of ps variable after sorting")
		}
	}

}
