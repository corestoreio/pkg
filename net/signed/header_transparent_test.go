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

package signed_test

import (
	"testing"
	"time"

	"github.com/corestoreio/pkg/net/signed"
	"github.com/corestoreio/pkg/storage/containable"
	"github.com/stretchr/testify/assert"
)

var _ signed.Cacher = (*set.InMemory)(nil)
var _ signed.Cacher = (*set.Mock)(nil)

func TestMakeTransparent(t *testing.T) {

	haveHash := []byte(`I'm your testing hash value`)
	haveTTL := time.Millisecond * 333

	cm := set.Mock{
		SetFn: func(hash []byte, ttl time.Duration) error {
			assert.Exactly(t, haveHash, hash)
			assert.Exactly(t, haveTTL, ttl)
			return nil
		},
		HasFn: func(hash []byte) bool {
			assert.Exactly(t, haveHash, hash)
			return false
		},
	}

	tp := signed.MakeTransparent(cm, haveTTL)
	assert.Empty(t, tp.HeaderKey())
	tp.Write(nil, haveHash)
	b, err := tp.Parse(nil)
	assert.Nil(t, b)
	assert.NoError(t, err)
}
