// Copyright 2015 CoreStore Authors
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
	"testing"

	"github.com/stretchr/testify/assert"
)

type ID int64

// ID is convenience helper to satisfy the interface Retriever
func (i ID) ID() int64 { return int64(i) }

func TestGetScopePath(t *testing.T) {
	tests := []struct {
		haveArg []OptionFunc
		want    string
	}{
		{[]OptionFunc{Path("a/b/c")}, DataScopeDefault + "/0/a/b/c"},
		{[]OptionFunc{Path("")}, DataScopeDefault + "/0/"},
		{[]OptionFunc{Path()}, DataScopeDefault + "/0/"},
		{[]OptionFunc{Scope(IDScopeDefault)}, DataScopeDefault + "/0/"},
		{[]OptionFunc{Scope(IDScopeWebsite)}, DataScopeDefault + "/0/"},
		{[]OptionFunc{Scope(IDScopeStore)}, DataScopeDefault + "/0/"},
		{[]OptionFunc{Path("a/b/c"), Scope(IDScopeWebsite)}, DataScopeDefault + "/0/a/b/c"},
		{[]OptionFunc{Path("a/b/c"), Scope(IDScopeWebsite, ID(2))}, DataScopeWebsites + "/2/a/b/c"},
		{[]OptionFunc{Path("a", "b", "c"), Scope(IDScopeWebsite, ID(200))}, DataScopeWebsites + "/200/a/b/c"},
		{[]OptionFunc{Path("a", "b", "c"), Scope(IDScopeStore, ID(4))}, DataScopeStores + "/4/a/b/c"},
		{[]OptionFunc{Path("a", "b"), Scope(IDScopeStore, ID(4))}, DataScopeStores + "/4/a"},
		{[]OptionFunc{nil, Scope(IDScopeStore, ID(4))}, DataScopeStores + "/4/"},
		{[]OptionFunc{Path("a", "b", "c"), ScopeStore(ID(5))}, DataScopeStores + "/5/a/b/c"},
		{[]OptionFunc{Path("a", "b", "c"), ScopeStore()}, DataScopeDefault + "/0/a/b/c"},
		{[]OptionFunc{Path("a", "b", "c"), ScopeWebsite(ID(50))}, DataScopeWebsites + "/50/a/b/c"},
		{[]OptionFunc{Path("a", "b", "c"), ScopeWebsite()}, DataScopeDefault + "/0/a/b/c"},
		{nil, ""},
	}

	for _, test := range tests {
		actualPath := getScopePath(test.haveArg...)
		assert.EqualValues(t, test.want, actualPath, "Test: %#v", test)
	}
}

// All benchmarks MacBook Air (13-inch, Mid 2012); 1.8 GHz Intel Core i5; 8 GB 1600 MHz DDR3

var benchmarkGetScopePath string

// BenchmarkGetScopePath____InMap	 2000000	       903 ns/op	     160 B/op	       9 allocs/op
func BenchmarkGetScopePath____InMap(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		benchmarkGetScopePath = getScopePath(Path("a", "b", "c"), Scope(IDScopeWebsite, ID(4)))
	}
}

// BenchmarkGetScopePath_NotInMap	 2000000	       959 ns/op	     184 B/op	      10 allocs/op
func BenchmarkGetScopePath_NotInMap(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		benchmarkGetScopePath = getScopePath(Path("a", "b", "c"), Scope(IDScopeWebsite, ID(40)))
	}
}

// BenchmarkGetScopePath____InMapNoJoin	 2000000	       723 ns/op	     160 B/op	       8 allocs/op
func BenchmarkGetScopePath____InMapNoJoin(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		benchmarkGetScopePath = getScopePath(Path("a/b/c"), Scope(IDScopeWebsite, ID(3)))
	}
}
