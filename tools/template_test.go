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
	"testing"

	"github.com/stretchr/testify/assert"
)

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
			pkg: "catalog",
			tplCode: `package {{ .xPackage }}
		var Table{{ .Table | prepareVar }} = 1`,
			data: struct {
				Package, Table string
			}{"catalog", "catalog_product_entity"},
			expTpl: []byte(``),
			expErr: true,
		},
	}

	for _, test := range tests {
		actual, err := GenerateCode(test.pkg, test.tplCode, test.data)
		if test.expErr {
			assert.Error(t, err)
		} else {
			assert.Equal(t, test.expTpl, actual)
			//t.Logf("\nExp: %s\nAct: %s", test.expTpl, actual)
		}
	}
}
