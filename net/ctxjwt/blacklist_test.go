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

package ctxjwt_test

import (
	"testing"
	"time"

	"github.com/corestoreio/csfw/net/ctxjwt"
	"github.com/corestoreio/csfw/store/scope"
	"github.com/stretchr/testify/assert"
)

type testBL struct {
	*testing.T
	theToken string
	exp      time.Duration
}

func (b *testBL) Set(theToken string, exp time.Duration) error {
	b.theToken = theToken
	b.exp = exp
	return nil
}
func (b *testBL) Has(_ string) bool { return false }

type testRealBL struct {
	theToken string
	exp      time.Duration
}

func (b *testRealBL) Set(t string, exp time.Duration) error {
	b.theToken = t
	b.exp = exp
	return nil
}
func (b *testRealBL) Has(t string) bool { return b.theToken == t }

func TestBlackLists(t *testing.T) {
	t.Parallel()
	tests := []struct {
		bl    ctxjwt.Blacklister
		token string
	}{
		{
			ctxjwt.NewBlackListSimpleMap(),
			`eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE0NTkxNTI3NTEsImlhdCI6MTQ1OTE0OTE1MSwibWFzY290IjoiZ29waGVyIn0.QzUJ5snl685Wmx4wXlCUykvBQMKn3OyL5MpnSaKrkdw`,
		},
		{
			ctxjwt.NewBlackListFreeCache(0),
			`eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE0NTkxNTI3NTEsImlhdCI6MTQ1OTE0OTE1MSwibWFzY290IjoiZ29waGVyIn0.QzUJ5snl685Wmx4wXlCUykvBQMKn3OyL5MpnSaKrkdw`,
		},
	}
	for i, test := range tests {
		assert.NoError(t, test.bl.Set(test.token, time.Second*1), "Index %d", i)
		assert.NoError(t, test.bl.Set(test.token+"2", time.Second*2), "Index %d", i)
		assert.True(t, test.bl.Has(test.token), "Index %d", i)
		time.Sleep(time.Second * 3)
		assert.NoError(t, test.bl.Set(test.token+"3", time.Second*2), "Index %d", i)
		assert.False(t, test.bl.Has(test.token), "Index %d", i)
		assert.False(t, test.bl.Has(test.token+"2"), "Index %d", i)
		assert.False(t, test.bl.Has(test.token), "Index %d", i)
		assert.True(t, test.bl.Has(test.token+"3"), "Index %d", i)
	}
}

const benchTokenCount = 100

func benchBlackList(b *testing.B, bl ctxjwt.Blacklister) {
	jwts := ctxjwt.MustNewService()
	var tokens [benchTokenCount]string

	for i := 0; i < benchTokenCount; i++ {
		claim := map[string]interface{}{
			"someKey": i,
		}
		var err error
		tokens[i], _, err = jwts.GenerateToken(scope.DefaultID, 0, claim)
		if err != nil {
			b.Fatal(err)
		}
	}

	b.ReportAllocs()
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			for i := 0; i < benchTokenCount; i++ {
				if err := bl.Set(tokens[i], time.Minute); err != nil {
					b.Fatal(err)
				}
				if bl.Has(tokens[i]) == false {
					b.Fatalf("Cannot find token %s with index %d", tokens[i], i)
				}
			}
		}
	})
}

// BenchmarkBlackListMap_Parallel-4      	    2000	    586726 ns/op	   31686 B/op	     200 allocs/op
func BenchmarkBlackListMap_Parallel(b *testing.B) {
	bl := ctxjwt.NewBlackListSimpleMap()
	benchBlackList(b, bl)
}

// BenchmarkBlackListFreeCache_Parallel-4	   30000	     59542 ns/op	   31781 B/op	     300 allocs/op
func BenchmarkBlackListFreeCache_Parallel(b *testing.B) {
	bl := ctxjwt.NewBlackListFreeCache(0)
	benchBlackList(b, bl)
}
