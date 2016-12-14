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

package dbr_test

import (
	"testing"

	"github.com/corestoreio/csfw/storage/dbr"
	"github.com/stretchr/testify/assert"
)

func TestStackIfNull(t *testing.T) {
	tests := []struct {
		alias      string
		columnName string
		defaultVal string
		want       string
	}{
		{
			"manufacturer", "value", "",
			"IFNULL(`manufacturerStore`.`value`,IFNULL(`manufacturerGroup`.`value`,IFNULL(`manufacturerWebsite`.`value`,IFNULL(`manufacturerDefault`.`value`,'')))) AS `manufacturer`",
		},
		{
			"manufacturer", "value", "0",
			"IFNULL(`manufacturerStore`.`value`,IFNULL(`manufacturerGroup`.`value`,IFNULL(`manufacturerWebsite`.`value`,IFNULL(`manufacturerDefault`.`value`,0)))) AS `manufacturer`",
		},
	}
	for i, test := range tests {
		assert.Exactly(t, test.want, dbr.EAVIfNull(test.alias, test.columnName, test.defaultVal), "Index %d", i)
	}
}
