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

		testExpiration(t, func() {
			mr.FastForward(time.Second * 2)
		}, objcache.NewRedisByURLClient(redConURL), newSrvOpt(JSONCodec{}))
	})

	t.Run("real redis integration", func(t *testing.T) {
		redConURL := lookupRedisEnv(t)
		testExpiration(t, func() {
			time.Sleep(time.Second * 2)
		}, objcache.NewRedisByURLClient(redConURL), newSrvOpt(JSONCodec{}))
	})
}

func TestWithRedisURL_Get_NotFound_Mock(t *testing.T) {
	t.Parallel()

	mr := miniredis.NewMiniRedis()
	assert.NoError(t, mr.Start())
	defer mr.Close()

	p, err := objcache.NewService(nil, objcache.NewRedisByURLClient("redis://"+mr.Addr()), newSrvOpt(JSONCodec{}))
	assert.NoError(t, err)
	defer func() {
		assert.NoError(t, p.Close())
	}()

	key := strs.RandAlnum(30)

	var newVal float64
	err = p.Get(context.TODO(), key, &newVal)
	assert.NoError(t, err, "Error: %+v", err)
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
		p, err := objcache.NewService(nil, objcache.NewRedisByURLClient(test.rawurl), newSrvOpt(JSONCodec{}))
		if test.errBhf > 0 {
			assert.True(t, test.errBhf.Match(err), "Index %d Error %+v", i, err)
			assert.Nil(t, p, "Index %d", i)
		} else {
			assert.NoError(t, err, "Index %d", i)
			assert.NotNil(t, p, "Index %d", i)
		}
	}

}

func TestWithRedisURL_ComplexParallel(t *testing.T) {
	mr := miniredis.NewMiniRedis()
	assert.NoError(t, mr.Start())
	defer mr.Close()
	redConURL := fmt.Sprintf("redis://%s/?db=2", mr.Addr())
	newServiceComplexParallelTest(t, objcache.NewRedisByURLClient(redConURL), nil)
}

func TestWithRedisURLMock_Delete(t *testing.T) {
	mr := miniredis.NewMiniRedis()
	assert.NoError(t, mr.Start())
	defer mr.Close()
	redConURL := fmt.Sprintf("redis://%s/?db=2", mr.Addr())
	newTestServiceDelete(t, objcache.NewRedisByURLClient(redConURL))
}

func TestWithRedisURLReal_Delete(t *testing.T) {
	redConURL := lookupRedisEnv(t)
	newTestServiceDelete(t, objcache.NewRedisByURLClient(redConURL))
}
