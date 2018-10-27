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
	"encoding/json"
	"io"
	"sort"
	"testing"

	"github.com/corestoreio/errors"
	"github.com/corestoreio/pkg/util/assert"
)

var (
	_ sort.Interface = (*options)(nil)
	_ errors.Kinder  = (*ErrKeyNotFound)(nil)
	_ error          = (*ErrKeyNotFound)(nil)
)

var _ Codecer = (*JSONCodec)(nil)

type JSONCodec struct{}

func (c JSONCodec) NewEncoder(w io.Writer) Encoder {
	return json.NewEncoder(w)
}

func (c JSONCodec) NewDecoder(r io.Reader) Decoder {
	return json.NewDecoder(r)
}

func TestWithSimpleSlowCacheMap(t *testing.T) {
	p, err := NewService(WithPooledEncoder(JSONCodec{}), WithSimpleSlowCacheMap())
	assert.NoError(t, err)
	defer assert.NoError(t, p.Close())

	t.Run("key not found", func(t *testing.T) {
		err := p.Get(context.TODO(), NewItem("upppsss", nil))
		assert.True(t, errors.NotFound.Match(err), "should have kind not found, but got: %+v", err)
	})
}
