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
	"github.com/corestoreio/pkg/storage/objcache"
	"github.com/corestoreio/pkg/util/assert"
)

func TestWithBigCache_Success(t *testing.T) {
	p, err := objcache.NewService(nil, objcache.NewBigCacheClient(bigcache.Config{}), newSrvOpt(JSONCodec{}))
	if err != nil {
		t.Fatal(err)
	}
	key := `key1`
	if err := p.Put(context.TODO(), key, math.Pi, 0); err != nil {
		t.Fatal(err)
	}

	var newVal float64
	if err := p.Get(context.TODO(), key, &newVal); err != nil {
		t.Fatal(err)
	}
	assert.Exactly(t, math.Pi, newVal)
}

func TestWithBigCache_Error(t *testing.T) {
	p, err := objcache.NewService(objcache.NewBigCacheClient(bigcache.Config{
		Shards: 3,
	}), nil, nil)
	assert.Nil(t, p)
	assert.EqualError(t, err, "Shards number must be power of two", "Error: %+v", err)
}

func TestBigCache_ServiceComplexParallel(t *testing.T) {
	newServiceComplexParallelTest(t, objcache.NewBigCacheClient(bigcache.Config{}), nil)
}

func TestWithBigCache_DecoderError(t *testing.T) {
	p, err := objcache.NewService(objcache.NewBlackHoleClient(nil), objcache.NewBigCacheClient(bigcache.Config{}), newSrvOpt(gobCodec{}))
	if err != nil {
		t.Fatal(err)
	}
	key := "key1"
	val1 := struct {
		Val string
	}{
		Val: "Gopher",
	}
	assert.NoError(t, p.Put(context.TODO(), key, val1, 0))

	var val2 struct {
		Val2 string
	}
	err = p.Get(context.TODO(), key, &val2)
	assert.EqualError(t, err, "[objcache] With key \"key1\" and dst type *struct { Val2 string }: gob: type mismatch: no fields matched compiling decoder for ", "Error: %s", err)
}

func TestWithBigCache_GetError(t *testing.T) {
	p, err := objcache.NewService(objcache.NewBlackHoleClient(nil), objcache.NewBigCacheClient(bigcache.Config{}), newSrvOpt(JSONCodec{}))
	assert.NoError(t, err)
	key := "key1"
	var ch struct {
		ErrChan string
	}
	err = p.Get(context.TODO(), key, &ch)
	assert.NoError(t, err, "Error: %s", err)
	assert.Empty(t, ch.ErrChan)
}

func TestWithBigCache_Delete(t *testing.T) {
	t.Parallel()
	newTestServiceDelete(t, objcache.NewBigCacheClient(bigcache.Config{}))
}
