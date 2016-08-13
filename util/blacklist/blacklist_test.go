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

package blacklist_test

import (
	"github.com/corestoreio/csfw/net/jwt"
	"github.com/corestoreio/csfw/util/blacklist"
	"github.com/pierrec/xxHash/xxHash64"
	"github.com/stretchr/testify/assert"
	"hash"
	"hash/fnv"
	"testing"
	"time"
)

var _ jwt.Blacklister = (*blacklist.InMemory)(nil)

func appendTo(b1 []byte, s string) []byte {
	bNew := make([]byte, len(b1)+len([]byte(s)))
	n := copy(bNew, b1)
	copy(bNew[n:], s)
	return bNew
}

func TestBlackLists(t *testing.T) {

	tests := []struct {
		bl jwt.Blacklister
	}{
		{blacklist.NewInMemory(fnv.New64a)},
		{blacklist.NewInMemory(func() hash.Hash64 { return xxHash64.New(33) })},
	}
	for i, test := range tests {

		id := []byte(`eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9`)

		assert.NoError(t, test.bl.Set(id, time.Second*1), "Index %d", i)
		assert.NoError(t, test.bl.Set(appendTo(id, "2"), time.Second*2), "Index %d", i)
		assert.True(t, test.bl.Has(id), "Index %d", i)
		time.Sleep(time.Second * 3)
		assert.NoError(t, test.bl.Set(appendTo(id, "3"), time.Second*2), "Index %d", i)
		assert.False(t, test.bl.Has(id), "Index %d", i)
		assert.False(t, test.bl.Has(appendTo(id, "2")), "Index %d", i)
		assert.False(t, test.bl.Has(id), "Index %d", i)
		assert.True(t, test.bl.Has(appendTo(id, "3")), "Index %d", i)
	}
}
