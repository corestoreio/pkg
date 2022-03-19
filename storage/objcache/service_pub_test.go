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

package objcache_test

import (
	"context"
	"io"
	"io/ioutil"
	"net"
	"reflect"
	"testing"

	"github.com/corestoreio/pkg/util/assert"
)

var _ io.Closer = (*objcache.Service[string])(nil)

func TestNewProcessor_EncoderError(t *testing.T) {
	t.Parallel()
	p, err := objcache.NewService[string](nil, objcache.NewCacheSimpleInmemory[string], internal.NewSrvOpt(internal.GobCodec{}))
	assert.NoError(t, err)

	ch := struct {
		ErrChan chan error
	}{
		ErrChan: make(chan error),
	}
	err = p.Set(context.TODO(), "key1", ch, 0)
	assert.EqualError(t, err, "[objcache] 1643662314576 With key key1 and dst type struct { ErrChan chan error }", "Error: %s", err)
}

type myString struct {
	data string
	err  error
}

func (ms *myString) Unmarshal(data []byte) error {
	ms.data = string(data)
	return ms.err
}

func (ms *myString) Marshal() ([]byte, error) {
	return []byte(ms.data), ms.err
}

func TestService_Encoding(t *testing.T) {
	p, err := objcache.NewService[string](nil, objcache.NewCacheSimpleInmemory[string], internal.NewSrvOpt(internal.GobCodec{}))
	assert.NoError(t, err)
	defer assert.NoError(t, p.Close())

	t.Run("marshal error", func(t *testing.T) {
		dErr := &myString{err: errors.New("Bad encoding")}
		err := p.Set(context.TODO(), "dErr", dErr, 0)
		assert.EqualError(t, err, "[objcache] 1643662002029 With key dErr and dst type *objcache_test.myString: Bad encoding")
	})
	t.Run("unmarshal error", func(t *testing.T) {
		err := p.Set(context.TODO(), "dErr2", 1, 0)
		assert.NoError(t, err)
		dErr := &myString{err: errors.New("Bad encoding")}
		err = p.Get(context.TODO(), "dErr2", dErr)
		assert.EqualError(t, err, "[objcache] 1643662475854 With key dErr2 and dst type *objcache_test.myString: Bad encoding")
	})

	t.Run("marshal success", func(t *testing.T) {
		d1 := &myString{data: "HelloWorld"}
		d2 := &myString{data: "HalloWelt"}

		err = p.SetMulti(context.TODO(), []string{"d1x", "d2x"}, []any{d1, d2}, nil)
		assert.NoError(t, err)

		d1.data = ""
		d2.data = ""
		err = p.GetMulti(context.TODO(), []string{"d1x", "d2x"}, []any{d1, d2})
		assert.NoError(t, err)

		assert.Exactly(t, "HelloWorld", d1.data)
		assert.Exactly(t, "HalloWelt", d2.data)
	})
}
