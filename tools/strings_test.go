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

package tools

import (
	"errors"
	"log"
	"testing"

	"github.com/juju/errgo"
	"github.com/stretchr/testify/assert"
)

func TestCamelize(t *testing.T) {
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
		{"idx_Tmp_cs", "IDXTMPCS"},
	}
	for _, test := range tests {
		assert.Equal(t, test.expected, Camelize(test.actual))
	}
}

// LogFatal logs an error as fatal with printed location and exists the program.
func TestLogFatal(t *testing.T) {
	defer func() { logFatalln = log.Fatalln }()
	var err error
	err = errors.New("Test")
	logFatalln = func(v ...interface{}) {
		assert.Contains(t, v[0].(string), "Error: Test")
	}
	LogFatal(err)

	err = errgo.New("Test")
	LogFatal(err)

	err = nil
	LogFatal(err)
}

func TestGenerateCode(t *testing.T) {
	tests := []struct {
		pkg, tplCode string
		data         interface{}
		expTpl       []byte
		expErr       bool
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
		},
		{
			pkg: "store",
			tplCode: `package {{ .Package }}
		var Table{{ prepareVarIndex 5 .Table }} = {{ "Gopher" | quote }}`,
			data: struct {
				Package, Table string
			}{"store", "core_store_group-01"},
			expTpl: []byte(`package store

var Table005CoreStoreGroup01 = ` + "`Gopher`\n"),
			expErr: false,
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
		},
	}

	for _, test := range tests {
		actual, err := GenerateCode(test.pkg, test.tplCode, test.data, nil)
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
