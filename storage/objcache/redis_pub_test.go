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

	"github.com/alicebob/miniredis"
	"github.com/corestoreio/errors"
	"github.com/corestoreio/pkg/storage/objcache"
	"github.com/corestoreio/pkg/util/assert"
	"github.com/corestoreio/pkg/util/strs"
	"github.com/garyburd/redigo/redis"
)

func TestWithDial_SetGet_Success_Live(t *testing.T) {
	t.Parallel()

	mr := miniredis.NewMiniRedis()
	if err := mr.Start(); err != nil {
		t.Fatal(err)
	}
	defer mr.Close()
	redConURL := "redis://" + mr.Addr()

	p, err := objcache.NewManager(objcache.WithRedisURL(redConURL), objcache.WithRedisPing(), objcache.WithEncoder(JSONCodec{}))
	if err != nil {
		t.Fatal(err)
	}
	defer func() {
		if err := p.Close(); err != nil {
			t.Fatal(err)
		}
	}()

	key := strs.RandAlnum(30)
	if err := p.Set(context.TODO(), key, math.Pi, nil); err != nil {
		t.Fatalf("Key %q Error: %s", key, err)
	}

	var newVal float64
	if err := p.Get(context.TODO(), key, &newVal, nil); err != nil {
		t.Fatalf("Key %q Error: %s", key, err)
	}
	assert.Exactly(t, math.Pi, newVal)
}

func TestWithDial_Get_NotFound_Live(t *testing.T) {
	t.Parallel()

	mr := miniredis.NewMiniRedis()
	if err := mr.Start(); err != nil {
		t.Fatal(err)
	}
	defer mr.Close()
	redConURL := "redis://" + mr.Addr()

	p, err := objcache.NewManager(objcache.WithRedisPing(), objcache.WithRedisURL(redConURL), objcache.WithEncoder(JSONCodec{}))
	if err != nil {
		t.Fatal(err)
	}
	defer func() {
		if err := p.Close(); err != nil {
			t.Fatal(err)
		}
	}()

	key := strs.RandAlnum(30)
	var newVal float64
	err = p.Get(context.TODO(), key, &newVal, nil)
	assert.True(t, errors.NotFound.Match(err), "%+v", err)
	assert.Empty(t, newVal)
}

func TestWithURL_SetGet_Success_Mock(t *testing.T) {
	t.Parallel()

	mr := miniredis.NewMiniRedis()
	if err := mr.Start(); err != nil {
		t.Fatal(err)
	}
	defer mr.Close()

	p, err := objcache.NewManager(objcache.WithRedisURL("redis://"+mr.Addr()), objcache.WithRedisPing(), objcache.WithEncoder(JSONCodec{}))
	if err != nil {
		t.Fatal(err)
	}
	defer func() {
		if err := p.Close(); err != nil {
			t.Fatal(err)
		}
	}()

	key := strs.RandAlnum(30)

	if err := p.Set(context.TODO(), key, math.Pi, nil); err != nil {
		t.Fatalf("Key %q Error: %s", key, err)
	}

	var newVal float64
	if err := p.Get(context.TODO(), key, &newVal, nil); err != nil {
		t.Fatalf("Key %q Error: %s", key, err)
	}
	assert.Exactly(t, math.Pi, newVal)
}

func TestWithDial_Get_NotFound_Mock(t *testing.T) {
	t.Parallel()

	mr := miniredis.NewMiniRedis()
	if err := mr.Start(); err != nil {
		t.Fatal(err)
	}
	defer mr.Close()

	p, err := objcache.NewManager(objcache.WithRedisURL("redis://"+mr.Addr()), objcache.WithRedisPing(), objcache.WithEncoder(JSONCodec{}))
	if err != nil {
		t.Fatal(err)
	}
	defer func() {
		if err := p.Close(); err != nil {
			t.Fatal(err)
		}
	}()

	key := strs.RandAlnum(30)

	var newVal float64
	err = p.Get(context.TODO(), key, &newVal, nil)
	assert.True(t, errors.NotFound.Match(err), "Error: %s", err)
	assert.Empty(t, newVal)
}

func TestWithDial_Get_Fatal_Mock(t *testing.T) {
	t.Parallel()

	mr := miniredis.NewMiniRedis()
	if err := mr.Start(); err != nil {
		t.Fatal(err)
	}
	defer mr.Close()

	p, err := objcache.NewManager(objcache.WithRedisURL("redis://"+mr.Addr()), objcache.WithRedisPing(), objcache.WithEncoder(JSONCodec{}))
	if err != nil {
		t.Fatal(err)
	}
	defer func() {
		if err := p.Close(); err != nil {
			t.Fatal(err)
		}
	}()

	key := strs.RandAlnum(30)

	var newVal float64
	err = p.Get(context.TODO(), key, &newVal, nil)
	assert.True(t, errors.NotFound.Match(err), "Error: %+v", err)
	assert.Empty(t, newVal)
}

func TestWithDial_ConFailure(t *testing.T) {
	t.Parallel()

	p, err := objcache.NewManager(objcache.WithRedisPing(), objcache.WithRedisClient(&redis.Pool{
		Dial: func() (redis.Conn, error) { return redis.Dial("tcp", "127.0.0.1:3344") }, // random port
	}), objcache.WithEncoder(JSONCodec{}))
	assert.True(t, errors.Fatal.Match(err), "Error: %s", err)
	assert.True(t, p == nil, "p is not nil")
}

func TestWithDialURL_ConFailure(t *testing.T) {
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
		p, err := objcache.NewManager(objcache.WithRedisURL(test.rawurl), objcache.WithRedisPing(), objcache.WithEncoder(JSONCodec{}))
		if test.errBhf > 0 {
			assert.True(t, test.errBhf.Match(err), "Index %d Error %+v", i, err)
			assert.Nil(t, p, "Index %d", i)
		} else {
			assert.NoError(t, err, "Index %d", i)
			assert.NotNil(t, p, "Index %d", i)
		}
	}

}

func TestProcessor_Parallel_GetSet_Redis(t *testing.T) {
	mr := miniredis.NewMiniRedis()
	if err := mr.Start(); err != nil {
		t.Fatalf("%+v", err)
	}
	defer mr.Close()
	redConURL := fmt.Sprintf("redis://%s/2", mr.Addr())
	newTestNewProcessor(t, objcache.WithRedisURL(redConURL))
}
