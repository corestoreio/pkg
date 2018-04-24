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

package dmltest_test

import (
	"bytes"
	"database/sql/driver"
	"fmt"
	"testing"

	"github.com/corestoreio/pkg/sql/dmltest"
	"github.com/stretchr/testify/assert"
)

func rowsToString(rows [][]driver.Value) string {
	var buf bytes.Buffer
	for _, row := range rows {
		for _, col := range row {
			switch v := col.(type) {
			case []byte:
				buf.Write(v)
			default:
				buf.WriteString("NULL")
			}
			buf.WriteRune('|')
		}
		buf.WriteRune('\n')
	}
	return buf.String()
}

func TestLoadCSVWithFile(t *testing.T) {
	t.Parallel()
	cols, rows, err := dmltest.LoadCSV(
		dmltest.WithFile("testdata", "core_config_data1.csv"),
		dmltest.WithTestMode(),
	)
	assert.NoError(t, err)
	assert.Exactly(t, []string{"config_id", "scope", "scope_id", "path", "value"}, cols)
	assert.Len(t, rows, 20)

	want := "1|default|0|cms/wysiwyg/enabled|disabled|\n2|default|0|general/region/display_all|1|\n3|default|0|general/region/state_required|AT,CA,CH,DE,EE,ES,FI,FR,LT,LV,RO,US|\n3|stores|2|general/region/state_required|AT|\n5|default|0|web/url/redirect_to_base|1|\n7|default|0|web/unsecure/base_url|http://magento-1-8.local/|\n7|websites|1|web/unsecure/base_url|http://magento-1-8a.dev/|\n8|default|0|web/unsecure/base_link_url|{{unsecure_base_url}}|\n9|default|0|web/unsecure/base_skin_url|{{unsecure_base_url}}skin/|\n10|default|0|web/unsecure/base_media_url|http://localhost:4711/media/|\n11|default|0|web/unsecure/base_js_url|{{unsecure_base_url}}js/|\n12|default|0|web/secure/base_url|http://magento-1-8.local/|\n13|default|0|web/secure/base_link_url|{{secure_base_url}}|\n14|default|0|web/secure/base_skin_url|{{secure_base_url}}skin/|\n15|default|0|web/secure/base_media_url|http://localhost:4711/media/|\n16|default|0|web/secure/base_js_url|{{secure_base_url}}js/|\n17|default|0|web/secure/use_in_frontend|0|\n18|default|0|web/secure/use_in_adminhtml|0|\n19|default|0|web/secure/offloader_header|SSL_OFFLOADED|\n20|default|0|web/default/front|NULL|\n"
	assert.Exactly(t, want, rowsToString(rows))
}

func TestLoadCSVWithReaderConfig(t *testing.T) {
	t.Parallel()
	cols, rows, err := dmltest.LoadCSV(
		dmltest.WithTestMode(),
		dmltest.WithFile("testdata", "core_config_data3.csv"),
		dmltest.WithReaderConfig(dmltest.CSVConfig{Comma: '|'}),
	)
	assert.NoError(t, err)
	assert.Exactly(t, []string{"config_id", "scope", "scope_id", "path", "value"}, cols)
	assert.Len(t, rows, 5)

	want := "1|default|0|cms/wysiwyg/enabled|disabled|\n2|default|0|general/region/display_all|1|\n3|default|0|general/region/state_required|AT,CA,CH,DE,EE,ES,FI,FR,LT,LV,RO,US|\n3|stores|2|general/region/state_required|AT|\n5|default|0|NULL|1|\n"
	assert.Exactly(t, want, rowsToString(rows))
}

func TestLoadCSVFileError(t *testing.T) {
	t.Parallel()
	cols, rows, err := dmltest.LoadCSV(
		dmltest.WithTestMode(),
		dmltest.WithFile("testdata", "core_config_dataXX.csv"),
	)
	assert.Nil(t, cols)
	assert.Nil(t, rows)
	assert.Contains(t, err.Error(), "core_config_dataXX.csv: no such file or directory")
}

func TestLoadCSVReadError(t *testing.T) {
	t.Parallel()
	cols, rows, err := dmltest.LoadCSV(
		dmltest.WithFile("testdata", "core_config_data2.csv"),
		dmltest.WithTestMode(),
	)
	assert.Exactly(t, []string{"config_id", "scope", "scope_id", "path", "value"}, cols)
	assert.Len(t, rows, 5)
	assert.EqualError(t, err, "[cstesting] csvReader.Read: record on line 7; parse error on line 8, column 0: extraneous or missing \" in quoted-field")
}

func TestMockRowsError(t *testing.T) {
	t.Parallel()
	r, err := dmltest.MockRows(dmltest.WithFile("non", "existent.csv"))
	assert.Nil(t, r)
	assert.Contains(t, err.Error(), "non/existent.csv: no such file or directory")
}

func TestMockRowsLoaded(t *testing.T) {
	t.Parallel()
	rows, err := dmltest.MockRows(
		dmltest.WithReaderConfig(dmltest.CSVConfig{Comma: '|'}),
		dmltest.WithFile("testdata", "core_config_data3.csv"),
		dmltest.WithTestMode(),
	)
	assert.NoError(t, err)
	assert.NotNil(t, rows)

	// Sorry for this test, but they removed the .Columns() function
	assert.Contains(t, fmt.Sprintf("%#v", rows), `sqlmock.Rows{cols:[]string{"config_id", "scope", "scope_id", "path", "value"}`)
}

func TestMustMockRows(t *testing.T) {
	t.Parallel()
	defer func() {
		if r := recover(); r != nil {
			assert.Contains(t, r.(error).Error(), "non/existent.csv: no such file or directory")
		} else {
			t.Fatal("Expecting a panic")
		}
	}()

	r := dmltest.MustMockRows(dmltest.WithFile("non", "existent.csv"))
	assert.Nil(t, r)
}
