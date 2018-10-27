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

// +build bigcache redis gob csall

package objcache_test

import (
	"encoding/gob"
	"encoding/json"
	"io"
	"os"
	"testing"

	"github.com/corestoreio/pkg/storage/objcache"
)

var _ objcache.Codecer = (*JSONCodec)(nil)

type JSONCodec struct{}

func (c JSONCodec) NewEncoder(w io.Writer) objcache.Encoder {
	return json.NewEncoder(w)
}

func (c JSONCodec) NewDecoder(r io.Reader) objcache.Decoder {
	return json.NewDecoder(r)
}

var _ objcache.Codecer = gobCodec{}

type gobCodec struct{}

func (c gobCodec) NewEncoder(w io.Writer) objcache.Encoder {
	return gob.NewEncoder(w)
}

func (c gobCodec) NewDecoder(r io.Reader) objcache.Decoder {
	return gob.NewDecoder(r)
}

func TestWithSimpleSlowCacheMap_Delete(t *testing.T) {
	t.Parallel()
	newTestServiceDelete(t, objcache.WithSimpleSlowCacheMap())
}

func lookupRedisEnv(t testing.TB) string {
	redConURL := os.Getenv("CS_REDIS_TEST")
	if redConURL == "" {
		t.Skip(`Skipping live test because environment CS_REDIS_TEST variable not found.
	export CS_REDIS_TEST="redis://127.0.0.1:6379/?db=3"
		`)
	}
	return redConURL
}
