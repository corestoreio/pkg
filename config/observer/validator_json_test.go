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

// +build csall json

package observer_test

import (
	"math"
	"testing"

	"github.com/corestoreio/pkg/config"
	"github.com/corestoreio/pkg/config/observer"
	"github.com/corestoreio/pkg/util/assert"
	"github.com/mailru/easyjson"
)

var (
	_ easyjson.Marshaler = (*observer.ValidateMinMaxInt)(nil)
)

func TestMinMaxInt64_MarshalJSON(t *testing.T) {
	t.Parallel()

	mm := observer.ValidateMinMaxInt{
		Conditions: []int64{-math.MaxInt64, math.MaxInt64},
	}

	data, err := mm.MarshalJSON()
	assert.NoError(t, err)
	assert.Exactly(t, "{\"conditions\":[-9223372036854775807,9223372036854775807]}", string(data))
}

func TestMinMaxInt64_UnmarshalJSON(t *testing.T) {
	t.Parallel()

	mm := new(observer.ValidateMinMaxInt)

	assert.NoError(t, mm.UnmarshalJSON([]byte("{\"conditions\":[-9223372036854775806,9223372036854775806]}")))
	assert.Exactly(t, &observer.ValidateMinMaxInt{
		Conditions: []int64{-math.MaxInt64 + 1, math.MaxInt64 - 1},
	}, mm)
}

func TestStrings_UnmarshalJSON(t *testing.T) {
	t.Parallel()
	rawJSON := []byte(`{"funcs":["Locale"],"csv_comma":"|","additional_allowed_values":["Klingon","Vulcan"]}`)

	var data observer.ValidatorArg
	assert.NoError(t, data.UnmarshalJSON(rawJSON))
	assert.Exactly(t, observer.ValidatorArg{
		Funcs:                   []string{"Locale"},
		CSVComma:                "|",
		AdditionalAllowedValues: []string{"Klingon", "Vulcan"},
	}, data)

	sv, err := observer.NewValidator(data)
	assert.NoError(t, err)
	value := []byte(`de_DE|Vulcan`)
	nValue, err := sv.Observe(config.Path{}, value, true)
	assert.NoError(t, err)
	assert.Exactly(t, value, nValue)

}

func TestStrings_MarshalJSON(t *testing.T) {
	t.Parallel()

	strs := observer.ValidatorArg{
		Funcs:                   []string{"Locale"},
		CSVComma:                "|",
		AdditionalAllowedValues: []string{"Klingon", "Vulcan"},
	}
	data, err := strs.MarshalJSON()
	assert.NoError(t, err)
	assert.Exactly(t, "{\"funcs\":[\"Locale\"],\"csv_comma\":\"|\",\"additional_allowed_values\":[\"Klingon\",\"Vulcan\"]}", string(data))
}
