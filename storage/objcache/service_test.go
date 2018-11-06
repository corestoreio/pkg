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

package objcache

import (
	"context"
	"encoding"
	"testing"

	"github.com/corestoreio/errors"
	"github.com/corestoreio/pkg/util/assert"
)

func TestNewProcessor_NewError(t *testing.T) {

	t.Run("level1 error", func(t *testing.T) {
		p, err := NewService(
			func() (Storager, error) { return nil, errors.NotImplemented.Newf("ups") },
			NewBlackHoleClient(nil),
			nil,
		)
		assert.Nil(t, p)
		assert.True(t, errors.NotImplemented.Match(err), "Error: %s", err)
	})

	t.Run("level2 error", func(t *testing.T) {
		p, err := NewService(
			NewBlackHoleClient(nil),
			func() (Storager, error) { return nil, errors.NotImplemented.Newf("ups") },
			nil,
		)
		assert.Nil(t, p)
		assert.True(t, errors.NotImplemented.Match(err), "Error: %s", err)
	})
}

var (
	_ encoding.TextMarshaler     = (*encodingText)(nil)
	_ encoding.TextUnmarshaler   = (*encodingText)(nil)
	_ encoding.BinaryUnmarshaler = (*encodingBinary)(nil)
	_ encoding.BinaryMarshaler   = (*encodingBinary)(nil)
)

type encodingText string

func (e *encodingText) UnmarshalText(text []byte) error {
	*e = encodingText(text)
	return nil
}
func (e encodingText) MarshalText() (text []byte, err error) { return []byte(e), nil }

type encodingBinary string

func (e *encodingBinary) UnmarshalBinary(text []byte) error {
	*e = encodingBinary(text)
	return nil
}
func (e encodingBinary) MarshalBinary() (text []byte, err error) { return []byte(e), nil }

func TestEncoding_Text_Binary(t *testing.T) {
	t.Parallel()

	// Not using any codec
	p, err := NewService(NewCacheSimpleInmemory, NewCacheSimpleInmemory, &ServiceOptions{Codec: nil})
	assert.NoError(t, err)
	defer assert.NoError(t, p.Close())

	ctx := context.TODO()
	t.Run("Text", func(t *testing.T) {
		obj := encodingText("Hello World 🎉")
		err := p.Set(ctx, "kt", obj, 0)
		assert.NoError(t, err)

		var obj2 encodingText
		err = p.Get(ctx, "kt", &obj2)
		assert.NoError(t, err)
		assert.Exactly(t, obj, obj2)
		assert.NotEmpty(t, obj2)
	})
	t.Run("Binary", func(t *testing.T) {
		obj := encodingBinary("Hello World 🎉")
		err := p.Set(ctx, "kt", obj, 0)
		assert.NoError(t, err)

		var obj2 encodingBinary
		err = p.Get(ctx, "kt", &obj2)
		assert.NoError(t, err)
		assert.Exactly(t, obj, obj2)
		assert.NotEmpty(t, obj2)
	})
}
