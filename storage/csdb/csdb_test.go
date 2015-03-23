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

package csdb

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

const (
	table1 Index = iota
	table2
	table3
	table4
)

var (
	tableMap = TableStructureSlice{
		table1: NewTableStructure(
			"catalog_category_anc_categs_index_idx",
			[]string{},
			[]string{
				"category_id",
				"path",
			},
		),
		table2: NewTableStructure(
			"catalog_category_anc_categs_index_tmp",
			[]string{},
			[]string{
				"category_id",
				"path",
			},
		),
		table3: NewTableStructure(
			"catalog_category_anc_products_index_idx",
			[]string{},
			[]string{
				"category_id",
				"product_id",
				"position",
			},
		),
	}
)

func TestTableStructure(t *testing.T) {
	s, err := tableMap.Structure(table1)
	assert.NotNil(t, s)
	assert.NoError(t, err)

	assert.Equal(t, "catalog_category_anc_categs_index_tmp", tableMap.Name(table2))
	assert.Equal(t, "", tableMap.Name(table4))

	s, err = tableMap.Structure(table4)
	assert.Nil(t, s)
	assert.Error(t, err)
}
