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

package modification_test

import (
	"testing"

	"github.com/corestoreio/errors"
	"github.com/corestoreio/pkg/config"
	"github.com/corestoreio/pkg/config/modification"
	"github.com/corestoreio/pkg/util/assert"
)

func TestMustNewStrings(t *testing.T) {
	t.Parallel()
	defer func() {
		if r := recover(); r != nil {
			if err, ok := r.(error); ok {
				assert.True(t, errors.NotSupported.Match(err))
			} else {
				t.Errorf("Panic should contain an error but got:\n%+v", r)
			}
		} else {
			t.Error("Expecting a panic but got nothing")
		}
	}()
	_ = modification.MustNewStrings(modification.Strings{
		Modificators: []string{"neverGonnaGiveYouUp"},
	})
}

func TestNewStrings(t *testing.T) {
	t.Parallel()

	t.Run("trim upper", func(t *testing.T) {
		ms := modification.MustNewStrings(modification.Strings{
			Modificators: []string{"trim", "upper"},
		})

		var p config.Path
		data := []byte(" \thello\n \t")
		have, err := ms.Observe(p, data, true)
		assert.NoError(t, err)
		assert.Exactly(t, "HELLO", string(have))
	})

	t.Run("custom operator returns error ", func(t *testing.T) {
		modification.RegisterOperator("csx", func(*config.Path, []byte) ([]byte, error) {
			return nil, errors.New("An error")
		})

		ms := modification.MustNewStrings(modification.Strings{
			Modificators: []string{"trim", "csx", "upper"},
		})

		var p config.Path
		data := []byte(" \thello\n \t")
		have, err := ms.Observe(p, data, true)
		assert.Nil(t, have)
		assert.True(t, errors.Interrupted.Match(err), "%+v", err)

	})
}
