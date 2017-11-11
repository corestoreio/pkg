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

package strs_test

import (
	"testing"

	"github.com/corestoreio/cspkg/util/strs"
	"github.com/stretchr/testify/assert"
)

func TestStrIsAlNum(t *testing.T) {
	t.Parallel()

	tests := []struct {
		have string
		want bool
	}{
		{"Hello World", false},
		{"HelloWorld", true},
		{"Hello1World", true},
		{"Hello0123456789", true},
		{"Hello0123456789€", false},
		{" Hello0123456789", false},
	}

	for _, test := range tests {
		assert.True(t, strs.IsAlNum(test.have) == test.want, "%#v", test)
	}
}

var benchStrIsAlNum bool

// BenchmarkStrIsAlNum	10000000	       132 ns/op	       0 B/op	       0 allocs/op
func BenchmarkStrIsAlNum(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		benchStrIsAlNum = strs.IsAlNum("Hello1WorldOfGophers")
	}
}

func TestUnderscoreToCamelCase(t *testing.T) {
	t.Parallel()
	tests := []struct {
		in  string
		out string
	}{
		{"catalog_product_entity", "CatalogProductEntity"},
		{"_catalog__product_entity", "CatalogProductEntity"},
		{"catalog_____product_entity_", "CatalogProductEntity"},
	}
	for _, test := range tests {
		assert.Equal(t, test.out, strs.ToCamelCase(test.in), "%#v", test)
	}
}

var benchCases string

// BenchmarkUnderscoreToCamelCase-4    	 1000000	      1457 ns/op	     192 B/op	       6 allocs/op
func BenchmarkUnderscoreToCamelCase(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		benchCases = strs.ToCamelCase("_catalog__product_entity")
	}
}

func TestCamelCaseToUnderscore(t *testing.T) {
	t.Parallel()
	tests := []struct {
		in  string
		out string
	}{
		{"CatalogProductEntity", "catalog_product_entity"},
		{"CatalogPPoductEntity", "catalog_p_poduct_entity"},
		{"catalogProductEntityE", "catalog_product_entity_e"},
		{"  catalogProductEntityE", "  catalog_product_entity_e"},
		{"  CatalogProductEntityE", "  _catalog_product_entity_e"}, // the leading underscore is a bug ... :-\
	}
	for _, test := range tests {
		assert.Equal(t, test.out, strs.FromCamelCase(test.in), "%#v", test)
	}
}

// BenchmarkCamelCaseToUnderscore-4    	 2000000	       928 ns/op	     288 B/op	       6 allocs/op
func BenchmarkCamelCaseToUnderscore(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		benchCases = strs.FromCamelCase("CatalogPPoductEntity")
	}
}

// BenchmarkCamelize-4                 	  500000	      2906 ns/op	     368 B/op	      15 allocs/op
func BenchmarkCamelize(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		benchCases = strs.ToGoCamelCase("catalog_____product_entity_")
	}
}

func TestCamelize(t *testing.T) {
	t.Parallel()
	tests := []struct {
		actual, expected string
	}{
		{"hello", "Hello"},
		{"hello_gopher", "HelloGopher"},
		{"hello_gopher_", "HelloGopher"},
		{"hello_gopher_id", "HelloGopherID"},
		{"hello_gopher_idx", "HelloGopherIDX"},
		{"idx_id", "IDXID"},
		{"idx_eav_id", "IDXEAVID"},
		{"idxeav_id", "IdxeavID"},
		{"idxeav_cs", "IdxeavCS"},
		{"idx_eav_cs", "IDXEAVCS"},
		{"idx_eav_cs_url", "IDXEAVCSURL"},
		{"hello_eav_idx_cs", "HelloEAVIDXCS"},
		{"hello_idx_Tmp_cs", "HelloIDXTMPCS"},
	}
	for _, test := range tests {
		assert.Equal(t, test.expected, strs.ToGoCamelCase(test.actual))
	}
}

func TestLintName(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name, want string
	}{
		{"foo_bar", "fooBar"},
		{"foo_bar_baz", "fooBarBaz"},
		{"Foo_bar", "FooBar"},
		{"foo_WiFi", "fooWiFi"},
		{"id", "id"},
		{"Id", "ID"},
		{"foo_id", "fooID"},
		{"fooId", "fooID"},
		{"fooUid", "fooUID"},
		{"idFoo", "idFoo"},
		{"uidFoo", "uidFoo"},
		{"midIdDle", "midIDDle"},
		{"APIProxy", "APIProxy"},
		{"ApiProxy", "APIProxy"},
		{"apiProxy", "apiProxy"},
		{"_Leading", "_Leading"},
		{"___Leading", "_Leading"},
		{"trailing_", "trailing"},
		{"trailing___", "trailing"},
		{"a_b", "aB"},
		{"a__b", "aB"},
		{"a___b", "aB"},
		{"Rpc1150", "RPC1150"},
		{"case3_1", "case3_1"},
		{"case3__1", "case3_1"},
		{"IEEE802_16bit", "IEEE802_16bit"},
		{"IEEE802_16Bit", "IEEE802_16Bit"},
		{"TableIndexUrlRewriteProductCategory", "TableIndexURLRewriteProductCategory"},
		{"IsHtmlAllowedOnFront", "IsHTMLAllowedOnFront"},
		{"UrlRewriteID", "URLRewriteID"},
	}
	for _, test := range tests {
		got := strs.LintName(test.name)
		if got != test.want {
			t.Errorf("lintName(%q) = %q, want %q", test.name, got, test.want)
		}
	}
}

// BenchmarkLintName-4                 	 1000000	      1485 ns/op	     144 B/op	       9 allocs/op
func BenchmarkLintName(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		benchCases = strs.LintName("____ApiProxyId")
	}
}

func TestRandAlnum(t *testing.T) {
	t.Parallel()
	s := strs.RandAlnum(18)
	assert.Len(t, s, 18)
}

func BenchmarkRandAlnum_18(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		benchCases = strs.RandAlnum(18)
	}
}
