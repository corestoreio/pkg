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

// +build bigcache csall

package objcache_test

import (
	"context"
	"math"
	"testing"

	"github.com/allegro/bigcache"
	"github.com/corestoreio/errors"
	"github.com/corestoreio/pkg/storage/objcache"
	"github.com/corestoreio/pkg/util/assert"
)

func TestWithBigCache_Success(t *testing.T) {
	p, err := objcache.NewService(objcache.WithBigCache(bigcache.Config{}), objcache.WithEncoder(JSONCodec{}))
	if err != nil {
		t.Fatal(err)
	}
	key := `key1`
	if err := p.Set(context.TODO(), objcache.NewItem(key, math.Pi)); err != nil {
		t.Fatal(err)
	}

	var newVal float64
	if err := p.Get(context.TODO(), objcache.NewItem(key, &newVal)); err != nil {
		t.Fatal(err)
	}
	assert.Exactly(t, math.Pi, newVal)
}

func TestWithBigCache_Error(t *testing.T) {
	p, err := objcache.NewService(objcache.WithBigCache(bigcache.Config{
		Shards: 3,
	}))
	assert.Nil(t, p)
	assert.EqualError(t, err, "[objcache] NewService applied options: Shards number must be power of two", "Error: %+v", err)
}

func TestProcessor_Parallel_GetSet_BigCache(t *testing.T) {
	newTestNewProcessor(t, objcache.WithBigCache(bigcache.Config{}))
}

func TestWithBigCache_DecoderError(t *testing.T) {
	p, err := objcache.NewService(objcache.WithPooledEncoder(gobCodec{}), objcache.WithBigCache(bigcache.Config{}))
	if err != nil {
		t.Fatal(err)
	}
	key := "key1"
	val1 := struct {
		Val string
	}{
		Val: "Gopher",
	}
	assert.NoError(t, p.Set(context.TODO(), objcache.NewItem(key, val1)))

	var val2 struct {
		Val2 string
	}
	err = p.Get(context.TODO(), objcache.NewItem(key, &val2))
	assert.EqualError(t, err, "[objcache] With key \"key1\" and dst type *struct { Val2 string }: gob: type mismatch: no fields matched compiling decoder for ", "Error: %s", err)
}

func TestWithBigCache_GetError(t *testing.T) {
	p, err := objcache.NewService(objcache.WithPooledEncoder(JSONCodec{}), objcache.WithBigCache(bigcache.Config{}))
	assert.NoError(t, err)
	key := "key1"
	var ch struct {
		ErrChan string
	}
	err = p.Get(context.TODO(), objcache.NewItem(key, ch))
	assert.True(t, errors.NotFound.Match(err), "Error: %s", err)
}

func TestWithBigCache_Delete(t *testing.T) {
	t.Parallel()
	newTestServiceDelete(t, objcache.WithBigCache(bigcache.Config{}))
}
