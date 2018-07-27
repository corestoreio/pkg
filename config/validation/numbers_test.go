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
	"testing"

	"github.com/alecthomas/assert"
	"github.com/corestoreio/errors"
	"github.com/corestoreio/pkg/config"
	"github.com/corestoreio/pkg/config/validation"
)

var (
	_ config.Observer = (*validation.MinMaxInt64)(nil)
)

func TestMinMaxInt_Observe(t *testing.T) {
	t.Parallel()
	var p config.Path
	t.Run("parse failed", func(t *testing.T) {
		mm, err := validation.NewMinMaxInt64(1, 2)
		assert.NoError(t, err)
		_, err = mm.Observe(p, []byte("NAN"), false)
		assert.EqualError(t, err, "strconv.ParseInt: parsing \"NAN\": invalid syntax")
	})
	t.Run("null", func(t *testing.T) {
		mm, err := validation.NewMinMaxInt64(1, 2)
		assert.NoError(t, err)
		ret, err := mm.Observe(p, nil, false)
		assert.NoError(t, err)
		assert.Nil(t, ret)
	})
	t.Run("not in range1", func(t *testing.T) {
		mm, err := validation.NewMinMaxInt64(1, 2)
		assert.NoError(t, err)
		ret, err := mm.Observe(p, []byte(`3`), false)
		assert.True(t, errors.OutOfRange.Match(err), "%+v", err)
		assert.Nil(t, ret)
	})
	t.Run("not in range2", func(t *testing.T) {
		mm, err := validation.NewMinMaxInt64(2, 1)
		assert.NoError(t, err)
		ret, err := mm.Observe(p, []byte(`3`), false)
		assert.True(t, errors.OutOfRange.Match(err), "%+v", err)
		assert.Nil(t, ret)
	})
	t.Run("in range1", func(t *testing.T) {
		mm, err := validation.NewMinMaxInt64(1, 2)
		assert.NoError(t, err)
		data := []byte(`2`)
		ret, err := mm.Observe(p, data, false)
		assert.NoError(t, err)
		assert.Exactly(t, data, ret)
	})
	t.Run("in range2", func(t *testing.T) {
		mm, err := validation.NewMinMaxInt64(1, 2)
		assert.NoError(t, err)
		data := []byte(`2`)
		ret, err := mm.Observe(p, data, false)
		assert.NoError(t, err)
		assert.Exactly(t, data, ret)
	})

	t.Run("partial validation enabled success", func(t *testing.T) {
		mm, err := validation.NewMinMaxInt64(1, 2, 5, 6, 7, 8)
		assert.NoError(t, err)
		mm.PartialValidation = true
		data := []byte(`6`)
		ret, err := mm.Observe(p, data, false)
		assert.NoError(t, err)
		assert.Exactly(t, data, ret)
	})

	t.Run("partial validation disabled fails", func(t *testing.T) {
		mm, err := validation.NewMinMaxInt64(1, 2, 5, 6, 7, 8)
		assert.NoError(t, err)
		mm.PartialValidation = false
		data := []byte(`6`)
		ret, err := mm.Observe(p, data, false)
		assert.True(t, errors.OutOfRange.Match(err), "%+v", err)
		assert.Nil(t, ret)
	})
}
