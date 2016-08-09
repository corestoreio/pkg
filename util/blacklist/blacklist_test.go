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
	"hash/fnv"
	"testing"
	"time"

	"github.com/corestoreio/csfw/net/jwt"
	"github.com/corestoreio/csfw/util/blacklist"
	"github.com/stretchr/testify/assert"
)

var _ jwt.Blacklister = (*blacklist.InMemory)(nil)

func appendTo(b1 []byte, s string) []byte {
	bNew := make([]byte, len(b1)+len([]byte(s)))
	n := copy(bNew, b1)
	copy(bNew[n:], s)
	return bNew
}

func TestBlackLists(t *testing.T) {
	t.Parallel()
	tests := []struct {
		bl    jwt.Blacklister
		token []byte
	}{
		{
			blacklist.NewInMemory(fnv.New64a()),
			[]byte(`eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE0NTkxNTI3NTEsImlhdCI6MTQ1OTE0OTE1MSwibWFzY290IjoiZ29waGVyIn0.QzUJ5snl685Wmx4wXlCUykvBQMKn3OyL5MpnSaKrkdw`),
		},
		//{
		//	blacklist.NewFreeCache(0),
		//	[]byte(`eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE0NTkxNTI3NTEsImlhdCI6MTQ1OTE0OTE1MSwibWFzY290IjoiZ29waGVyIn0.QzUJ5snl685Wmx4wXlCUykvBQMKn3OyL5MpnSaKrkdw`),
		//},
	}
	for i, test := range tests {
		assert.NoError(t, test.bl.Set(test.token, time.Second*1), "Index %d", i)
		assert.NoError(t, test.bl.Set(appendTo(test.token, "2"), time.Second*2), "Index %d", i)
		assert.True(t, test.bl.Has(test.token), "Index %d", i)
		time.Sleep(time.Second * 3)
		assert.NoError(t, test.bl.Set(appendTo(test.token, "3"), time.Second*2), "Index %d", i)
		assert.False(t, test.bl.Has(test.token), "Index %d", i)
		assert.False(t, test.bl.Has(appendTo(test.token, "2")), "Index %d", i)
		assert.False(t, test.bl.Has(test.token), "Index %d", i)
		assert.True(t, test.bl.Has(appendTo(test.token, "3")), "Index %d", i)
	}
}
