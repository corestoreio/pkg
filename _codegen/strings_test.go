// Copyright 2015-2017, Cyrill @ Schumacher.fm and the CoreStore contributors
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

package codegen

// global variables in benchmark tests are used to disable compiler optimizations

import (
	"os"
	"testing"
	"text/template"

	"github.com/stretchr/testify/assert"
)

func TestOFile(t *testing.T) {
	s := string(os.PathSeparator)
	tests := []struct {
		base, dir, file, want string
	}{
		{"", "", "hello", "hello.go"},
		{s + "usr", "local", "hello", s + "usr" + s + "localhello.go"},
		{s + "usr", "local" + s, "hello", s + "usr" + s + "localhello.go"},
		{s + "usr", "local" + s + "world_", "hello", s + "usr" + s + "local" + s + "world_hello.go"},
	}
	for _, test := range tests {
		have := NewOFile(test.base).AppendDir(test.dir).AppendName(test.file).String()
		assert.Equal(t, test.want, have)
	}
}

func TestGenerateCode(t *testing.T) {
	tests := []struct {
		pkg, tplCode string
		data         interface{}
		expTpl       []byte
		expErr       bool
		fncMap       template.FuncMap
	}{
		{
			pkg: "catalog",
			tplCode: `package {{ .Package }}
		var Table{{ .Table | prepareVar }} = {{ "Gopher" | quote }}`,
			data: struct {
				Package, Table string
			}{"catalog", "catalog_product_entity"},
			expTpl: []byte(`package catalog

var TableProductEntity = ` + "`Gopher`" + `
`),
			expErr: false,
			fncMap: nil,
		},
		{
			pkg: "aa",
			tplCode: `package {{ .Package }}
		var Table{{ .Table | prepareVar }} = {{ "Gopher" | quote }}`,
			data: struct {
				Package, Table string
			}{"aa", "aa"},
			expTpl: []byte(`package aa

var TableAa = ` + "`Gopher`\n"),
			expErr: false,
			fncMap: nil,
		},
		{
			pkg: "store",
			tplCode: `package {{ .Package }}
		var Table{{ prepareVarIndex 5 .Table }} = {{ "Gopher" | quote }}
		var Gogento = "{{ "GoGentö" | toLowerFirstTest}}"{{"" | toLowerFirstTest}}`,
			data: struct {
				Package, Table string
			}{"store", "core_store_group-01"},
			expTpl: []byte(`package store

var Table005CoreStoreGroup01 = ` + "`Gopher`" + `
var Gogento = "goGentö"` + "\n"),
			expErr: false,
			fncMap: template.FuncMap{
				"toLowerFirstTest": toLowerFirst,
			},
		},
		{
			pkg: "catalog",
			tplCode: `package {{ .xPackage }}
		var Table{{ .Table | prepareVar }} = 1`,
			data: struct {
				Package, Table string
			}{"catalog", "catalog_product_entity"},
			expTpl: []byte(``),
			expErr: true,
			fncMap: nil,
		},
		{
			pkg: "catalog",
			tplCode: `package {{ .Package }}
		var Table{{ .Table | prepareVar }} = ""
		""
		struct
		`,
			data: struct {
				Package, Table string
			}{"catalog", "catalog_product_entity"},
			expTpl: []byte(``),
			expErr: true,
			fncMap: nil,
		},
	}

	for _, test := range tests {
		actual, err := GenerateCode(test.pkg, test.tplCode, test.data, test.fncMap)
		if test.expErr {
			assert.Error(t, err)
		} else {
			assert.Equal(t, test.expTpl, actual)
			//t.Logf("\nExp: %s\nAct: %s", test.expTpl, actual)
		}
	}
}

func TestRandSeq(t *testing.T) {
	assert.Len(t, randSeq(10), 10)
}

var benchRandSeq = ""

// BenchmarkRandSeq	  100000	     12901 ns/op	     112 B/op	       2 allocs/op
func BenchmarkRandSeq(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		benchRandSeq = randSeq(20)
	}
}

const testStringTablePrefix = `SELECT * FROM {{tableprefix}}store
		JOIN {{tableprefix}}store_group ON col1=col2
		WHERE {{tableprefix}}store = 1`

func TestReplaceTablePrefix(t *testing.T) {
	assert.Equal(t, ReplaceTablePrefix(testStringTablePrefix), `SELECT * FROM store
		JOIN store_group ON col1=col2
		WHERE store = 1`)

	TablePrefix = "mystore1_"
	assert.Equal(t, ReplaceTablePrefix(testStringTablePrefix), `SELECT * FROM mystore1_store
		JOIN mystore1_store_group ON col1=col2
		WHERE mystore1_store = 1`)
	TablePrefix = ""
}

var benchReplaceTablePrefix = ""

// BenchmarkReplaceTablePrefix	 1000000	      1031 ns/op	     160 B/op	       2 allocs/op
func BenchmarkReplaceTablePrefix(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		benchReplaceTablePrefix = ReplaceTablePrefix(testStringTablePrefix)
	}
}

func TestPrepareVar(t *testing.T) {
	tests := []struct {
		pkg  string
		have string
		want string
	}{
		{"catalog", "catalog_product", "Product"},
		{"store", "AU+NZ", "AuNz"},
		{"store", "Madison Island", "MadisonIsland"},
	}
	for _, test := range tests {
		assert.Equal(t, test.want, prepareVar(test.pkg)(test.have))
	}
}

var benchPrepareVar = ""

// BenchmarkPrepareVar	  500000	      3078 ns/op	     288 B/op	      14 allocs/op
func BenchmarkPrepareVar(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		benchPrepareVar = prepareVar("catalog")(`catalog_product_attribute_collection `)
	}
}

func TestExtractPathType(t *testing.T) {
	tests := []struct {
		have   string
		wantIP string // IP import path
		wantFT string // FT FuncType
		isErr  bool
	}{
		{
			have:   "",
			wantIP: "",
			wantFT: "",
			isErr:  false,
		},
		{
			have:   "gopkg.in_corestoreio)catalog.v3.Product()",
			wantIP: "",
			wantFT: "",
			isErr:  true,
		},
		{
			have:   "gopkg.in/corestoreio/catalog.v3.Product()",
			wantIP: "gopkg.in/corestoreio/catalog.v3",
			wantFT: "catalog.Product()",
			isErr:  false,
		},
		{
			have:   "gopkg.in/corestoreio/catalogv3Product()",
			wantIP: "",
			wantFT: "",
			isErr:  true,
		},
		{
			have:   "gopkg.in/yaml.v2.Unmarshal()",
			wantIP: "gopkg.in/yaml.v2",
			wantFT: "yaml.Unmarshal()",
			isErr:  false,
		},
		{
			have:   "gopkg.in/yaml.v2.1.Unmarshal()",
			wantIP: "gopkg.in/yaml.v2.1",
			wantFT: "yaml.Unmarshal()",
			isErr:  false,
		},
		{
			have:   "gopkg.in/yaml.v2.1.30.Unmarshal()",
			wantIP: "gopkg.in/yaml.v2.1.30",
			wantFT: "yaml.Unmarshal()",
			isErr:  false,
		},
		{
			have:   "github.com/corestoreio/cspkg/catalog.Product()",
			wantIP: "github.com/corestoreio/cspkg/catalog",
			wantFT: "catalog.Product()",
			isErr:  false,
		},
		{
			have:   "github.com/corestoreio/cspkg/catalog/catattr.NewHandler({{.EntityTypeID}})",
			wantIP: "github.com/corestoreio/cspkg/catalog/catattr",
			wantFT: "catattr.NewHandler({{.EntityTypeID}})",
			isErr:  false,
		},
		{
			have:   "github.com/corestoreio.ARandomType",
			wantIP: "github.com/corestoreio",
			wantFT: "corestoreio.ARandomType",
			isErr:  false,
		},
	}
	for _, test := range tests {
		ip, errIP := ExtractImportPath(test.have)
		ft, errFT := ExtractFuncType(test.have)
		if test.isErr {
			assert.Error(t, errIP)
			assert.Error(t, errFT)
		} else {
			assert.NoError(t, errIP)
			assert.NoError(t, errFT)
		}
		assert.Equal(t, test.wantIP, ip)
		assert.Equal(t, test.wantFT, ft)
	}
}

var benchExtractImportPath string

// BenchmarkExtractImportPath	 1000000	      1133 ns/op	     176 B/op	       6 allocs/op
func BenchmarkExtractImportPath(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		benchExtractImportPath, _ = ExtractImportPath("gopkg.in/yaml.v2.1.30.Unmarshal()")
	}
}

var benchExtractFuncType string

// BenchmarkExtractFuncType	 2000000	       881 ns/op	     144 B/op	       4 allocs/op
func BenchmarkExtractFuncType(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		benchExtractFuncType, _ = ExtractFuncType("gopkg.in/yaml.v2.1.30.Unmarshal()")
	}
}
