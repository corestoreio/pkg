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

// +build redis csall

package objcache_test

import (
	"context"
	"fmt"
	"math"
	"testing"
	"time"

	"github.com/alicebob/miniredis"
	"github.com/corestoreio/errors"
	"github.com/corestoreio/pkg/storage/objcache"
	"github.com/corestoreio/pkg/util/assert"
	"github.com/corestoreio/pkg/util/strs"
	"github.com/gomodule/redigo/redis"
)

func TestWithRedisURL_PutGet_Success(t *testing.T) {
	t.Parallel()

	t.Run("miniredis", func(t *testing.T) {
		mr := miniredis.NewMiniRedis()
		if err := mr.Start(); err != nil {
			t.Fatal(err)
		}
		defer mr.Close()
		redConURL := "redis://" + mr.Addr()

		testWithRedisURL_PutGet_Success(t, func() {
			mr.FastForward(time.Second * 2)
		}, objcache.WithRedisURL(redConURL), newSrvOpt(JSONCodec{}))
	})

	t.Run("real redis integration", func(t *testing.T) {
		redConURL := lookupRedisEnv(t)
		testWithRedisURL_PutGet_Success(t, func() {
			time.Sleep(time.Second * 2)
		}, objcache.WithRedisURL(redConURL), newSrvOpt(JSONCodec{}))
	})
}

func testWithRedisURL_PutGet_Success(t *testing.T, cb func(), level2 objcache.NewStorageFn, so *objcache.ServiceOptions) {
	p, err := objcache.NewService(nil, level2, so)
	if err != nil {
		t.Fatal(err)
	}
	defer func() {
		assert.NoError(t, p.Close())
	}()

	key := strs.RandAlnum(30)
	if err := p.Put(context.TODO(), key, math.Pi, time.Second); err != nil {
		t.Fatalf("Key %q Error: %s", key, err)
	}

	var newVal float64
	err = p.Get(context.TODO(), key, &newVal)
	assert.NoError(t, err, "%+v", err)
	assert.Exactly(t, math.Pi, newVal)

	cb()

	newVal = 0
	err = p.Get(context.TODO(), key, &newVal)
	assert.True(t, errors.NotFound.Match(err), "%+v", err)
}

func TestWithRedisURL_Get_NotFound_Mock(t *testing.T) {
	t.Parallel()

	mr := miniredis.NewMiniRedis()
	assert.NoError(t, mr.Start())
	defer mr.Close()

	p, err := objcache.NewService(nil, objcache.WithRedisURL("redis://"+mr.Addr()), newSrvOpt(JSONCodec{}))
	assert.NoError(t, err)
	defer func() {
		assert.NoError(t, p.Close())
	}()

	key := strs.RandAlnum(30)

	var newVal float64
	err = p.Get(context.TODO(), key, &newVal)
	assert.True(t, errors.NotFound.Match(err), "Error: %+v", err)
	assert.Empty(t, newVal)
}

func TestWithRedisURLURL_ConFailure_Dial(t *testing.T) {
	t.Parallel()

	p, err := objcache.NewService(nil, objcache.NewRedisClient(&redis.Pool{
		Dial: func() (redis.Conn, error) { return redis.Dial("tcp", "127.0.0.1:53344") }, // random port
	}, &objcache.RedisOption{}), newSrvOpt(JSONCodec{}))
	assert.True(t, errors.Fatal.Match(err), "Error: %s", err)
	assert.True(t, p == nil, "p is not nil")
}

func TestWithRedisURL_ConFailure(t *testing.T) {
	t.Parallel()

	var dialErrors = []struct {
		rawurl string
		errBhf errors.Kind
	}{
		{
			"localhost",
			errors.NotSupported, // "invalid redis URL scheme",
		},
		// The error message for invalid hosts is different in different
		// versions of Go, so just check that there is an error message.
		{
			"redis://weird url",
			errors.Fatal,
		},
		{
			"redis://foo:bar:baz",
			errors.Fatal,
		},
		{
			"http://www.google.com",
			errors.NotSupported, // "invalid redis URL scheme: http",
		},
		{
			"redis://localhost:6379?db=",
			errors.Fatal, // "invalid database: abc123",
		},
	}
	for i, test := range dialErrors {
		p, err := objcache.NewService(nil, objcache.WithRedisURL(test.rawurl), newSrvOpt(JSONCodec{}))
		if test.errBhf > 0 {
			assert.True(t, test.errBhf.Match(err), "Index %d Error %+v", i, err)
			assert.Nil(t, p, "Index %d", i)
		} else {
			assert.NoError(t, err, "Index %d", i)
			assert.NotNil(t, p, "Index %d", i)
		}
	}

}

func TestWithRedisURL_Parallel_GetPut(t *testing.T) {
	mr := miniredis.NewMiniRedis()
	assert.NoError(t, mr.Start())
	defer mr.Close()
	redConURL := fmt.Sprintf("redis://%s/?db=2", mr.Addr())
	newTestNewProcessor(t, objcache.WithRedisURL(redConURL))
}

func TestWithRedisURLMock_Delete(t *testing.T) {
	mr := miniredis.NewMiniRedis()
	assert.NoError(t, mr.Start())
	defer mr.Close()
	redConURL := fmt.Sprintf("redis://%s/?db=2", mr.Addr())
	newTestServiceDelete(t, objcache.WithRedisURL(redConURL))
}

func TestWithRedisURLReal_Delete(t *testing.T) {
	redConURL := lookupRedisEnv(t)
	newTestServiceDelete(t, objcache.WithRedisURL(redConURL))
}
