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
	"math"
	"testing"

	"github.com/corestoreio/pkg/config/validation"
	"github.com/corestoreio/pkg/config/validation/json"
	"github.com/stretchr/testify/assert"
)

func TestMinMaxInt64_MarshalJSON(t *testing.T) {
	t.Parallel()

	mm := json.MinMaxInt64{
		MinMaxInt64: validation.MinMaxInt64{Conditions: []int64{-math.MaxInt64, math.MaxInt64}},
	}

	data, err := mm.MarshalJSON()
	assert.NoError(t, err)
	assert.Exactly(t, "{\"conditions\":[-9223372036854775807,9223372036854775807]}", string(data))
}

func TestMinMaxInt64_UnmarshalJSON(t *testing.T) {
	t.Parallel()

	mm := new(json.MinMaxInt64)

	assert.NoError(t, mm.UnmarshalJSON([]byte("{\"conditions\":[-9223372036854775806,9223372036854775806]}")))
	assert.Exactly(t, &json.MinMaxInt64{
		MinMaxInt64: validation.MinMaxInt64{
			Conditions: []int64{-math.MaxInt64 + 1, math.MaxInt64 - 1},
		},
	}, mm)
}
