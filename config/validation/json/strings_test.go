// Copyright 2015-present, Cyrill @ Schumacher.fm and the CoreStore contributors
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

package json_test

import (
	"testing"

	"github.com/corestoreio/pkg/config"
	"github.com/corestoreio/pkg/config/validation"
	"github.com/corestoreio/pkg/config/validation/json"
	"github.com/corestoreio/pkg/util/assert"
	"github.com/mailru/easyjson"
)

var (
	_ easyjson.Marshaler = (*json.Strings)(nil)
)

func TestStrings_UnmarshalJSON(t *testing.T) {
	t.Parallel()
	rawJSON := []byte(`{"validators":["Locale"],"csv_comma":"|","additional_allowed_values":["Klingon","Vulcan"]}`)

	var data json.Strings
	assert.NoError(t, data.UnmarshalJSON(rawJSON))
	assert.Exactly(t, json.Strings{
		Strings: validation.Strings{
			Validators:              []string{"Locale"},
			CSVComma:                "|",
			AdditionalAllowedValues: []string{"Klingon", "Vulcan"},
		},
	}, data)

	sv, err := validation.NewStrings(data.Strings)
	assert.NoError(t, err)
	value := []byte(`de_DE|Vulcan`)
	nValue, err := sv.Observe(config.Path{}, value, true)
	assert.NoError(t, err)
	assert.Exactly(t, value, nValue)

}

func TestStrings_MarshalJSON(t *testing.T) {
	t.Parallel()

	strs := json.Strings{
		Strings: validation.Strings{
			Validators:              []string{"Locale"},
			CSVComma:                "|",
			AdditionalAllowedValues: []string{"Klingon", "Vulcan"},
		},
	}
	data, err := strs.MarshalJSON()
	assert.NoError(t, err)
	assert.Exactly(t, "{\"validators\":[\"Locale\"],\"csv_comma\":\"|\",\"additional_allowed_values\":[\"Klingon\",\"Vulcan\"]}", string(data))
}
