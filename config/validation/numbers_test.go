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

package validation_test

import (
	"math"
	"testing"

	"github.com/corestoreio/errors"
	"github.com/corestoreio/pkg/config"
	"github.com/corestoreio/pkg/config/validation"
	"github.com/stretchr/testify/assert"
)

var (
	_ config.Observer = (*validation.MinMaxInt64)(nil)
)

func TestMinMaxInt_Observe(t *testing.T) {
	t.Parallel()
	var p config.Path
	t.Run("parse failed", func(t *testing.T) {
		mm := &validation.MinMaxInt64{
			Min: 1,
			Max: 2,
		}
		_, err := mm.Observe(p, []byte("NAN"), false)
		assert.EqualError(t, err, "strconv.ParseInt: parsing \"NAN\": invalid syntax")
	})
	t.Run("null", func(t *testing.T) {
		mm := &validation.MinMaxInt64{
			Min: 1,
			Max: 2,
		}
		ret, err := mm.Observe(p, nil, false)
		assert.NoError(t, err)
		assert.Nil(t, ret)
	})
	t.Run("not in range1", func(t *testing.T) {
		mm := &validation.MinMaxInt64{
			Min: 1,
			Max: 2,
		}
		ret, err := mm.Observe(p, []byte(`3`), false)
		assert.True(t, errors.OutOfRange.Match(err), "%+v", err)
		assert.Nil(t, ret)
	})
	t.Run("not in range2", func(t *testing.T) {
		mm := &validation.MinMaxInt64{
			Min: 2,
			Max: 1,
		}
		ret, err := mm.Observe(p, []byte(`3`), false)
		assert.True(t, errors.OutOfRange.Match(err), "%+v", err)
		assert.Nil(t, ret)
	})
	t.Run("in range1", func(t *testing.T) {
		mm := &validation.MinMaxInt64{
			Min: 1,
			Max: 2,
		}
		data := []byte(`2`)
		ret, err := mm.Observe(p, data, false)
		assert.NoError(t, err)
		assert.Exactly(t, ret, data)
	})
	t.Run("in range2", func(t *testing.T) {
		mm := &validation.MinMaxInt64{
			Min: 1,
			Max: 2,
		}
		data := []byte(`2`)
		ret, err := mm.Observe(p, data, false)
		assert.NoError(t, err)
		assert.Exactly(t, ret, data)
	})
}

func TestMinMaxInt64_MarshalJSON(t *testing.T) {
	t.Parallel()

	mm := &validation.MinMaxInt64{
		Min: -math.MaxInt64,
		Max: math.MaxInt64,
	}
	data, err := mm.MarshalJSON()
	assert.NoError(t, err)
	assert.Exactly(t, "{\"min\":-9223372036854775807,\"max\":9223372036854775807}", string(data))
}

func TestMinMaxInt64_UnmarshalJSON(t *testing.T) {
	t.Parallel()

	mm := new(validation.MinMaxInt64)

	assert.NoError(t, mm.UnmarshalJSON([]byte("{\"min\":-9223372036854775806,\"max\":9223372036854775806}")))
	assert.Exactly(t, &validation.MinMaxInt64{
		Min: -math.MaxInt64 + 1,
		Max: math.MaxInt64 - 1,
	}, mm)
}
