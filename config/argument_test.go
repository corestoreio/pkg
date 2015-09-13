// Copyright 2015, Cyrill @ Schumacher.fm and the CoreStore contributors
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
	"strings"
	"testing"

	"github.com/corestoreio/csfw/utils/log"
	"github.com/stretchr/testify/assert"
)

func init() {
	log.Set(log.NewStdLogger())
	log.SetLevel(log.StdLevelError)
}

func TestScopeKey(t *testing.T) {
	tests := []struct {
		haveArg []ArgFunc
		want    string
	}{
		{[]ArgFunc{Path("a/b/c")}, ScopeRangeDefault + "/0/a/b/c"},
		{[]ArgFunc{Path("")}, ""},
		{[]ArgFunc{Path()}, ""},
		{[]ArgFunc{Scope(ScopeDefaultID, nil)}, ""},
		{[]ArgFunc{Scope(ScopeWebsiteID, nil)}, ""},
		{[]ArgFunc{Scope(ScopeStoreID, nil)}, ""},
		{[]ArgFunc{Path("a/b/c"), Scope(ScopeWebsiteID, nil)}, ScopeRangeDefault + "/0/a/b/c"},
		{[]ArgFunc{Path("a/b/c"), Scope(ScopeWebsiteID, ScopeID(2))}, ScopeRangeWebsites + "/2/a/b/c"},
		{[]ArgFunc{Path("a", "b", "c"), Scope(ScopeWebsiteID, ScopeID(200))}, ScopeRangeWebsites + "/200/a/b/c"},
		{[]ArgFunc{Path("a", "b", "c"), Scope(ScopeStoreID, ScopeID(4))}, ScopeRangeStores + "/4/a/b/c"},
		{[]ArgFunc{Path("a", "b"), Scope(ScopeStoreID, ScopeID(4))}, ScopeRangeStores + "/4/a"},
		{[]ArgFunc{nil, Scope(ScopeStoreID, ScopeID(4))}, ""},
		{[]ArgFunc{Path("a", "b", "c"), ScopeStore(ScopeID(5))}, ScopeRangeStores + "/5/a/b/c"},
		{[]ArgFunc{Path("a", "b", "c"), ScopeStore(nil)}, ScopeRangeDefault + "/0/a/b/c"},
		{[]ArgFunc{Path("a", "b", "c"), ScopeWebsite(ScopeID(50))}, ScopeRangeWebsites + "/50/a/b/c"},
		{[]ArgFunc{Path("a", "b", "c"), ScopeWebsite(nil)}, ScopeRangeDefault + "/0/a/b/c"},
		{nil, ""},
	}

	for _, test := range tests {
		arg := newArg(test.haveArg...)
		actualPath := arg.scopePath()
		assert.EqualValues(t, test.want, actualPath, "Test: %#v", test)
	}
}

func TestScopeKeyValue(t *testing.T) {
	tests := []struct {
		haveArg []ArgFunc
		want    string
	}{
		{[]ArgFunc{Value(1), Path("a/b/c")}, ScopeRangeDefault + "/0/a/b/c"},
		{[]ArgFunc{Value("1"), Path("")}, ""},
		{[]ArgFunc{Value(1.1), Path()}, ""},
		{[]ArgFunc{Value(1), Scope(ScopeDefaultID, nil)}, ""},
		{[]ArgFunc{Value(1), Scope(ScopeWebsiteID, nil)}, ""},
		{[]ArgFunc{Value(1), Scope(ScopeStoreID, nil)}, ""},
		{[]ArgFunc{Value(1), Path("a/b/c"), Scope(ScopeWebsiteID, nil)}, ScopeRangeDefault + "/0/a/b/c"},
		{[]ArgFunc{Value(1), Path("a/b/c"), Scope(ScopeWebsiteID, ScopeID(2))}, ScopeRangeWebsites + "/2/a/b/c"},
		{[]ArgFunc{Value(1), Path("a", "b", "c"), Scope(ScopeWebsiteID, ScopeID(200))}, ScopeRangeWebsites + "/200/a/b/c"},
		{[]ArgFunc{Value(1), Path("a", "b", "c"), Scope(ScopeStoreID, ScopeID(4))}, ScopeRangeStores + "/4/a/b/c"},
		{[]ArgFunc{Value(1), Path("a", "b"), Scope(ScopeStoreID, ScopeID(4))}, ScopeRangeStores + "/4/a"},
		{[]ArgFunc{Value(1), nil, Scope(ScopeStoreID, ScopeID(4))}, ""},
		{[]ArgFunc{Value(1), Path("a", "b", "c"), ScopeStore(ScopeID(5))}, ScopeRangeStores + "/5/a/b/c"},
		{[]ArgFunc{Value(1.2), Path("a", "b", "c"), ScopeStore(nil)}, ScopeRangeDefault + "/0/a/b/c"},
		{[]ArgFunc{Value(1.3), Path("a", "b", "c"), ScopeWebsite(ScopeID(50))}, ScopeRangeWebsites + "/50/a/b/c"},
		{[]ArgFunc{ValueReader(strings.NewReader("a config value")), Path("a", "b", "c"), ScopeWebsite(nil)}, ScopeRangeDefault + "/0/a/b/c"},
		{nil, ""},
	}

	for _, test := range tests {
		a := newArg(test.haveArg...)
		actualPath, actualVal := a.scopePath(), a.v
		assert.EqualValues(t, test.want, actualPath, "Test: %#v", test)
		if test.haveArg != nil {
			assert.NotEmpty(t, actualVal, "Test: %#v", test)
		} else {
			assert.Empty(t, actualVal, "Test: %#v", test)
		}
	}
}

// All benchmarks MacBook Air (13-inch, Mid 2012); 1.8 GHz Intel Core i5; 8 GB 1600 MHz DDR3

var benchmarkScopeKey string

// BenchmarkScopeKey____InMap	 2000000	       936 ns/op	     176 B/op	       9 allocs/op
func BenchmarkScopeKey____InMap(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		arg := newArg(Path("a", "b", "c"), Scope(ScopeWebsiteID, ScopeID(4)))
		benchmarkScopeKey = arg.scopePath()
	}
}

// BenchmarkScopeKey_NotInMap	 2000000	       992 ns/op	     200 B/op	      10 allocs/op
func BenchmarkScopeKey_NotInMap(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		arg := newArg(Path("a", "b", "c"), Scope(ScopeWebsiteID, ScopeID(40)))
		benchmarkScopeKey = arg.scopePath()
	}
}

// BenchmarkScopeKey____InMapNoJoin	 2000000	       824 ns/op	     176 B/op	       8 allocs/op
func BenchmarkScopeKey____InMapNoJoin(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		arg := newArg(Path("a/b/c"), Scope(ScopeWebsiteID, ScopeID(3)))
		benchmarkScopeKey = arg.scopePath()
	}
}
