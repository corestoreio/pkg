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

package scope

import (
	"errors"
	"testing"

	"strconv"

	"github.com/stretchr/testify/assert"
)

func TestScopeBits(t *testing.T) {
	t.Parallel()
	const (
		scope1 Scope = iota + 1
		scope2
		scope3
		scope4
		scope5
	)

	tests := []struct {
		have    []Scope
		want    Scope
		notWant Scope
		human   []string
		string
	}{
		{[]Scope{scope1, scope2}, scope2, scope3, []string{"Default", "Website"}, "Default,Website"},
		{[]Scope{scope3, scope4}, scope3, scope2, []string{"Group", "Store"}, "Group,Store"},
		{[]Scope{scope4, scope5}, scope4, scope2, []string{"Store", "Scope(5)"}, "Store,Scope(5)"},
	}

	for _, test := range tests {
		var b Perm
		b.Set(test.have...)
		if b.Has(test.want) == false {
			t.Errorf("%d should contain %d", b, test.want)
		}
		if b.Has(test.notWant) {
			t.Errorf("%d should not contain %d", b, test.notWant)
		}
		assert.EqualValues(t, test.human, b.Human())
		assert.EqualValues(t, test.string, b.String())
	}
}

func TestFromString(t *testing.T) {
	t.Parallel()
	tests := []struct {
		have string
		want Scope
	}{
		{"asdasd", DefaultID},
		{strDefault, DefaultID},
		{strWebsites, WebsiteID},
		{strStores, StoreID},
	}
	for _, test := range tests {
		assert.Equal(t, test.want, FromString(test.have))
	}
}

func TestFromScope(t *testing.T) {
	t.Parallel()
	tests := []struct {
		have Scope
		want StrScope
	}{
		{DefaultID, StrDefault},
		{WebsiteID, StrWebsites},
		{StoreID, StrStores},
	}
	for _, test := range tests {
		assert.Equal(t, test.want, FromScope(test.have))
	}
}

func TestStrScope(t *testing.T) {
	assert.Equal(t, strDefault, StrDefault.String())
	assert.Equal(t, strWebsites, StrWebsites.String())
	assert.Equal(t, strStores, StrStores.String())
}

func TestStrScopeFQPath(t *testing.T) {
	t.Parallel()
	tests := []struct {
		str  StrScope
		id   string
		path []string
		want string
	}{
		{StrDefault, "0", []string{"system/dev/debug"}, strDefault + "/0/system/dev/debug"},
		{StrDefault, "33", []string{"system", "dev", "debug"}, strDefault + "/0/system/dev/debug"},
		{StrWebsites, "0", []string{"system/dev/debug"}, strWebsites + "/0/system/dev/debug"},
		{StrWebsites, "343", []string{"system", "dev", "debug"}, strWebsites + "/343/system/dev/debug"},
	}
	for _, test := range tests {
		assert.Equal(t, test.want, test.str.FQPath(test.id, test.path...))
	}
	assert.Equal(t, "stores/7475/catalog/frontend/list_allow_all", StrStores.FQPathInt64(7475, "catalog", "frontend", "list_allow_all"))
	assert.Equal(t, "stores/5/catalog/frontend/list_allow_all", StrStores.FQPathInt64(5, "catalog", "frontend", "list_allow_all"))
}

var benchmarkStrScopeFQPath string

// BenchmarkStrScopeFQPath-4	 5000000	       384 ns/op	      32 B/op	       1 allocs/op
func BenchmarkStrScopeFQPath(b *testing.B) {
	want := strWebsites + "/4/system/dev/debug"
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		benchmarkStrScopeFQPath = StrWebsites.FQPath("4", "system", "dev", "debug")
	}
	if benchmarkStrScopeFQPath != want {
		b.Errorf("Want: %s; Have, %s", want, benchmarkStrScopeFQPath)
	}
}

func benchmarkStrScopeFQPathInt64(scopeID int64, b *testing.B) {
	want := strWebsites + "/" + strconv.FormatInt(scopeID, 10) + "/system/dev/debug"
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		benchmarkStrScopeFQPath = StrWebsites.FQPathInt64(scopeID, "system", "dev", "debug")
	}
	if benchmarkStrScopeFQPath != want {
		b.Errorf("Want: %s; Have, %s", want, benchmarkStrScopeFQPath)
	}
}

// BenchmarkStrScopeFQPathInt64__Cached-4	 5000000	       383 ns/op	      32 B/op	       1 allocs/op
func BenchmarkStrScopeFQPathInt64__Cached(b *testing.B) {
	benchmarkStrScopeFQPathInt64(4, b)
}

// BenchmarkStrScopeFQPathInt64UnCached-4	 3000000	       438 ns/op	      48 B/op	       2 allocs/op
func BenchmarkStrScopeFQPathInt64UnCached(b *testing.B) {
	benchmarkStrScopeFQPathInt64(40, b)
}

func TestSplitFQPath(t *testing.T) {
	t.Parallel()
	tests := []struct {
		have        string
		wantScope   string
		wantScopeID int64
		wantPath    string
		wantErr     error
	}{
		{"groups/1/catalog/frontend/list_allow_all", "groups", 0, "", ErrUnsupportedScope},
		{"stores/7475/catalog/frontend/list_allow_all", strStores, 7475, "catalog/frontend/list_allow_all", nil},
		{"websites/1/catalog/frontend/list_allow_all", strWebsites, 1, "catalog/frontend/list_allow_all", nil},
		{"default/0/catalog/frontend/list_allow_all", strDefault, 0, "catalog/frontend/list_allow_all", nil},
		{"default//catalog/frontend/list_allow_all", strDefault, 0, "catalog/frontend/list_allow_all", errors.New("strconv.ParseInt: parsing \"\\uf8ff\": invalid syntax")},
		{"stores/123/catalog/index", "", 0, "", errors.New("Incorrect fully qualified path: \"stores/123/catalog/index\"")},
	}
	for _, test := range tests {
		haveScope, haveScopeID, havePath, haveErr := SplitFQPath(test.have)

		if test.wantErr != nil {
			assert.EqualError(t, haveErr, test.wantErr.Error(), "Test %v", test)
		} else {
			assert.NoError(t, haveErr, "Test %v", test)
		}
		assert.Exactly(t, test.wantScope, haveScope, "Test %v", test)
		assert.Exactly(t, test.wantScopeID, haveScopeID, "Test %v", test)
		assert.Exactly(t, test.wantPath, havePath, "Test %v", test)
	}
}

var benchmarkReverseFQPath = struct {
	scope   string
	scopeID int64
	path    string
	err     error
}{}

// BenchmarkReverseFQPath-4 	10000000	       121 ns/op	       0 B/op	       0 allocs/op
func BenchmarkReverseFQPath(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		benchmarkReverseFQPath.scope, benchmarkReverseFQPath.scopeID, benchmarkReverseFQPath.path, benchmarkReverseFQPath.err = SplitFQPath("stores/7475/catalog/frontend/list_allow_all")
		if benchmarkReverseFQPath.err != nil {
			b.Error(benchmarkReverseFQPath.err)
		}
	}
}
