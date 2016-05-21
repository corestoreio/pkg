// Copyright 2015-2016, Cyrill @ Schumacher.fm and the CoreStore contributors
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

package tcredis

import (
	"math"
	"os"
	"testing"

	"github.com/corestoreio/csfw/storage/typecache"
	"github.com/corestoreio/csfw/util"
	"github.com/corestoreio/csfw/util/errors"
	"github.com/garyburd/redigo/redis"
	"github.com/rafaeljusto/redigomock"
	"github.com/stretchr/testify/assert"
)

var _ typecache.Cacher = (*wrapper)(nil)

func TestWithDial_SetGet_Success_Live(t *testing.T) {

	redConURL := os.Getenv("CS_REDIS_TEST") // redis://127.0.0.1:6379/3
	if redConURL == "" {
		t.Skip(`Skipping live test because environment CS_REDIS_TEST variable not found.
	export CS_REDIS_TEST="redis://127.0.0.1:6379/3"
		`)
	}

	p, err := typecache.NewProcessor(WithDialURL(redConURL))
	if err != nil {
		t.Fatal(err)
	}
	defer func() {
		if err := p.Cache.Close(); err != nil {
			t.Fatal(err)
		}
	}()

	var key = []byte(util.RandAlnum(30))
	if err := p.Set(key, math.Pi); err != nil {
		t.Fatalf("Key %q Error: %s", key, err)
	}

	var newVal float64
	if err := p.Get(key, &newVal); err != nil {
		t.Fatalf("Key %q Error: %s", key, err)
	}
	assert.Exactly(t, math.Pi, newVal)
}

func TestWithDial_Get_NotFound_Live(t *testing.T) {

	redConURL := os.Getenv("CS_REDIS_TEST") // redis://127.0.0.1:6379/3
	if redConURL == "" {
		t.Skip(`Skipping live test because environment CS_REDIS_TEST variable not found.
	export CS_REDIS_TEST="redis://127.0.0.1:6379/3"
		`)
	}

	p, err := typecache.NewProcessor(WithDialURL(redConURL))
	if err != nil {
		t.Fatal(err)
	}
	defer func() {
		if err := p.Cache.Close(); err != nil {
			t.Fatal(err)
		}
	}()

	var key = []byte(util.RandAlnum(30))

	var newVal float64
	err = p.Get(key, &newVal)
	assert.True(t, errors.IsNotFound(err), "Error: %s", err)
	assert.Empty(t, newVal)
}

func TestWithDial_SetGet_Success_Mock(t *testing.T) {
	c := redigomock.NewConn()

	p, err := typecache.NewProcessor(WithCon(c))
	if err != nil {
		t.Fatal(err)
	}
	defer func() {
		if err := p.Cache.Close(); err != nil {
			t.Fatal(err)
		}
	}()

	var key = []byte(util.RandAlnum(30))
	c.Command("SET", key, []uint8{0xb, 0x8, 0x0, 0xf8, 0x18, 0x2d, 0x44, 0x54, 0xfb, 0x21, 0x9, 0x40}).Expect([]uint8{0xb, 0x8, 0x0, 0xf8, 0x18, 0x2d, 0x44, 0x54, 0xfb, 0x21, 0x9, 0x40})
	if err := p.Set(key, math.Pi); err != nil {
		t.Fatalf("Key %q Error: %s", key, err)
	}

	var newVal float64
	c.Command("GET", key).Expect([]uint8{0xb, 0x8, 0x0, 0xf8, 0x18, 0x2d, 0x44, 0x54, 0xfb, 0x21, 0x9, 0x40})
	if err := p.Get(key, &newVal); err != nil {
		t.Fatalf("Key %q Error: %s", key, err)
	}
	assert.Exactly(t, math.Pi, newVal)
}

func TestWithDial_Get_NotFound_Mock(t *testing.T) {

	c := redigomock.NewConn()
	p, err := typecache.NewProcessor(WithCon(c))
	if err != nil {
		t.Fatal(err)
	}
	defer func() {
		if err := p.Cache.Close(); err != nil {
			t.Fatal(err)
		}
	}()

	var key = []byte(util.RandAlnum(30))
	c.Command("GET", key).Expect(nil)
	var newVal float64
	err = p.Get(key, &newVal)
	assert.True(t, errors.IsNotFound(err), "Error: %s", err)
	assert.Empty(t, newVal)
}

func TestWithDial_Get_Fatal_Mock(t *testing.T) {

	c := redigomock.NewConn()
	p, err := typecache.NewProcessor(WithCon(c))
	if err != nil {
		t.Fatal(err)
	}
	defer func() {
		if err := p.Cache.Close(); err != nil {
			t.Fatal(err)
		}
	}()

	var key = []byte(util.RandAlnum(30))
	c.Command("GET", key).ExpectError(errors.New("Some Error"))
	var newVal float64
	err = p.Get(key, &newVal)
	assert.True(t, errors.IsFatal(err), "Error: %s", err)
	assert.Empty(t, newVal)
}

func TestWithDial_ConFailure(t *testing.T) {
	p, err := typecache.NewProcessor(WithDial("tcp", "127.0.0.1:380", redis.DialDatabase(33)))
	assert.True(t, errors.IsFatal(err))
	assert.Nil(t, p)
}

func TestWithDialURL_ConFailure(t *testing.T) {

	var dialErrors = []struct {
		rawurl string
		errBhf errors.BehaviourFunc
	}{
		{
			"localhost",
			errors.IsFatal, // "invalid redis URL scheme",
		},
		// The error message for invalid hosts is diffferent in different
		// versions of Go, so just check that there is an error message.
		{
			"redis://weird url",
			errors.IsFatal,
		},
		{
			"redis://foo:bar:baz",
			errors.IsFatal,
		},
		{
			"http://www.google.com",
			errors.IsFatal, // "invalid redis URL scheme: http",
		},
		{
			"redis://localhost:6379/abc123",
			errors.IsFatal, // "invalid database: abc123",
		},
	}
	for i, test := range dialErrors {
		p, err := typecache.NewProcessor(WithDialURL(test.rawurl, redis.DialDatabase(33)))
		if test.errBhf != nil {
			assert.True(t, errors.IsFatal(err), "Index %d Error %s", i, err)
			assert.Nil(t, p, "Index %d", i)
		} else {
			assert.NoError(t, err, "Index %d", i)
			assert.NotNil(t, p, "Index %d", i)
		}
	}

}
