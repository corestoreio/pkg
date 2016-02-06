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

package cstesting_test

import (
	"github.com/corestoreio/csfw/util/cstesting"

	"encoding/json"
	"github.com/stretchr/testify/assert"
	"path/filepath"
	"testing"
)

func TestLoadCSVOk(t *testing.T) {

	dataFile := filepath.Join(cstesting.RootPath, "util", "cstesting", "testdata", "core_config_data1.csv")
	cols, rows, err := cstesting.LoadCSV(dataFile)
	assert.NoError(t, err)
	assert.Exactly(t, []string{"config_id", "scope", "scope_id", "path", "value"}, cols)
	assert.Len(t, rows, 20)

	jData, err := json.Marshal(rows)
	assert.NoError(t, err)
	assert.Exactly(t,
		`[["1","default","0","cms/wysiwyg/enabled","disabled"],["2","default","0","general/region/display_all","1"],["3","default","0","general/region/state_required","AT,CA,CH,DE,EE,ES,FI,FR,LT,LV,RO,US"],["3","stores","2","general/region/state_required","AT"],["5","default","0","web/url/redirect_to_base","1"],["7","default","0","web/unsecure/base_url","http://magento-1-8.local/"],["7","websites","1","web/unsecure/base_url","http://magento-1-8a.dev/"],["8","default","0","web/unsecure/base_link_url","{{unsecure_base_url}}"],["9","default","0","web/unsecure/base_skin_url","{{unsecure_base_url}}skin/"],["10","default","0","web/unsecure/base_media_url","http://localhost:4711/media/"],["11","default","0","web/unsecure/base_js_url","{{unsecure_base_url}}js/"],["12","default","0","web/secure/base_url","http://magento-1-8.local/"],["13","default","0","web/secure/base_link_url","{{secure_base_url}}"],["14","default","0","web/secure/base_skin_url","{{secure_base_url}}skin/"],["15","default","0","web/secure/base_media_url","http://localhost:4711/media/"],["16","default","0","web/secure/base_js_url","{{secure_base_url}}js/"],["17","default","0","web/secure/use_in_frontend","0"],["18","default","0","web/secure/use_in_adminhtml","0"],["19","default","0","web/secure/offloader_header","SSL_OFFLOADED"],["20","default","0","web/default/front",""]]`,
		string(jData),
	)
}

func TestLoadCSVFileError(t *testing.T) {
	dataFile := filepath.Join(cstesting.RootPath, "util", "cstesting", "testdata", "core_config_dataXX.csv")
	cols, rows, err := cstesting.LoadCSV(dataFile)
	assert.Nil(t, cols)
	assert.Nil(t, rows)
	assert.Contains(t, err.Error(), "core_config_dataXX.csv: no such file or directory")
}

func TestLoadCSVReadError(t *testing.T) {
	dataFile := filepath.Join(cstesting.RootPath, "util", "cstesting", "testdata", "core_config_data2.csv")
	cols, rows, err := cstesting.LoadCSV(dataFile)
	assert.Exactly(t, []string{"config_id", "scope", "scope_id", "path", "value"}, cols)
	assert.Len(t, rows, 5)
	assert.EqualError(t, err, `line 8, column 0: extraneous " in field`)
}
