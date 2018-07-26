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
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestStrings_UnmarshalJSON(t *testing.T) {
	t.Parallel()
	rawJSON := []byte(`{"validators":["Locale"],"csv_comma":"|","additional_allowed_values":["Klingon","Vulcan"]}`)

	var data json.Strings
	require.NoError(t, data.UnmarshalJSON(rawJSON))
	assert.Exactly(t, json.Strings{
		Strings: validation.Strings{
			Validators:              []string{"Locale"},
			CSVComma:                "|",
			AdditionalAllowedValues: []string{"Klingon", "Vulcan"},
		},
	}, data)

	sv, err := validation.NewStrings(data.Strings)
	require.NoError(t, err)
	value := []byte(`de_DE|Vulcan`)
	nValue, err := sv.Observe(config.Path{}, value, true)
	require.NoError(t, err)
	assert.Exactly(t, value, nValue)

}
