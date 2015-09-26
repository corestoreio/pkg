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

package scope

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestScopeBits(t *testing.T) {
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
	}{
		{[]Scope{scope1, scope2}, scope2, scope3, []string{"ScopeDefault", "ScopeWebsite"}},
		{[]Scope{scope3, scope4}, scope3, scope2, []string{"ScopeGroup", "ScopeStore"}},
		{[]Scope{scope4, scope5}, scope4, scope2, []string{"ScopeStore", "ScopeGroup(5)"}},
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
	}
}

func TestFromString(t *testing.T) {
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
	tests := []struct {
		str  StrScope
		id   string
		path []string
		want string
	}{
		{StrDefault, "0", []string{"system/dev/debug"}, strDefault + "/0/system/dev/debug"},
		{StrDefault, "33", []string{"system", "dev", "debug"}, strDefault + "/33/system/dev/debug"},
		{StrWebsites, "0", []string{"system/dev/debug"}, strWebsites + "/0/system/dev/debug"},
		{StrWebsites, "343", []string{"system", "dev", "debug"}, strWebsites + "/343/system/dev/debug"},
	}
	for _, test := range tests {
		assert.Equal(t, test.want, test.str.FQPath(test.id, test.path...))
	}
}

var benchmarkStrScopeFQPath string

// BenchmarkStrScopeFQPath-4	 5000000	       286 ns/op	     144 B/op	       2 allocs/op
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
