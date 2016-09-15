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

package containable_test

import (
	"github.com/corestoreio/csfw/net/jwt"
	"github.com/corestoreio/csfw/storage/containable"
	"github.com/stretchr/testify/assert"
	"strconv"
	"testing"
	"time"
)

var _ containable.Container = (*containable.InMemory)(nil)
var _ containable.Container = (*containable.Mock)(nil)

func appendTo(b1 []byte, s string) []byte {
	bNew := make([]byte, len(b1)+len([]byte(s)))
	n := copy(bNew, b1)
	copy(bNew[n:], s)
	return bNew
}

func TestNewInMemory(t *testing.T) {
	t.Parallel()
	tests := []struct {
		bl jwt.Blacklister
	}{
		{containable.NewInMemory()},
		{containable.NewInMemory()},
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

func TestNewInMemory_Purge(t *testing.T) {
	t.Parallel()
	m := containable.NewInMemory()
	for i := 0; i < 6; i++ {
		id := []byte(`eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9`)
		id = strconv.AppendInt(id, int64(i), 10)
		assert.NoError(t, m.Set(id, time.Second))
		time.Sleep(time.Second) // bit lame this test but so far ok, can be refactored one day.
	}
	assert.Exactly(t, 3, m.Len())

}
