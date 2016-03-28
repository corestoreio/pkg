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

package scope_test

import (
	"math"
	"testing"

	"github.com/corestoreio/csfw/store/scope"
	"github.com/stretchr/testify/assert"
	"sync"
)

func TestHashString(t *testing.T) {
	t.Parallel()
	s := scope.NewHash(scope.StoreID, 33).String()
	assert.Exactly(t, "Scope(Store) ID(33)", s)
}

func TestHashSegment(t *testing.T) {
	t.Parallel()
	tests := []struct {
		h    scope.Hash
		want uint8
	}{
		{scope.DefaultHash, 0},
		{scope.NewHash(scope.Scope(0), 0), 0},
		{scope.NewHash(scope.Scope(1), 0), 0},
		{scope.NewHash(scope.DefaultID, -1), 0},
		{scope.NewHash(scope.DefaultID, 1), 0},
		{scope.NewHash(scope.DefaultID, 0), 0},
		{scope.NewHash(scope.StoreID, 0), 0},
		{scope.NewHash(scope.StoreID, 1), 1},
		{scope.NewHash(scope.StoreID, 2), 2},
		{scope.NewHash(scope.StoreID, 255), 255},
		{scope.NewHash(scope.StoreID, 256), 0},
		{scope.NewHash(scope.StoreID, 257), 1},
		{scope.NewHash(scope.StoreID, scope.MaxStoreID-1), 254},
		{scope.NewHash(scope.StoreID, scope.MaxStoreID), 255},
		{scope.NewHash(scope.StoreID, scope.MaxStoreID+1), 0},
		{scope.NewHash(scope.StoreID, scope.MaxStoreID+2), 0},
		{scope.NewHash(scope.StoreID, -scope.MaxStoreID), 0},
		{scope.NewHash(scope.Scope(7), 1), 1},
	}
	for i, test := range tests {
		if want, have := test.want, test.h.Segment(); want != have {
			t.Errorf("Index %03d: Want %03d Have %03d", i, want, have)
		}
	}
}

func TestNewHashError(t *testing.T) {
	t.Parallel()
	h := scope.NewHash(scope.AbsentID, -1)
	assert.Exactly(t, scope.Hash(0), h)
}

func TestFromHashError(t *testing.T) {
	t.Parallel()
	scp, id := scope.Hash(math.MaxUint32).Unpack()
	assert.Exactly(t, scope.AbsentID, scp)
	assert.Exactly(t, int64(-1), id)
}

func TestHashValid(t *testing.T) {
	t.Parallel()
	t.Logf("[Info] Max Store ID: %d", scope.MaxStoreID)

	if testing.Short() {
		t.Skip("Skipping in short mode")
	}

	//var collisionCheck = make(map[scope.Hash]bool) // just in case ...
	var wg sync.WaitGroup
	var scp = scope.AbsentID
	for ; scp < math.MaxUint8; scp++ {
		wg.Add(1)
		go func(theScp scope.Scope) {
			defer wg.Done()
			for id := int64(0); id < scope.MaxStoreID; id++ {
				haveHash := scope.NewHash(theScp, id)

				haveScp, haveID := haveHash.Unpack()
				if haveScp != theScp {
					t.Fatalf("Have Scope: %d, Want Scope: %d", haveScp, theScp)
				}

				wantID := id
				if theScp < scope.WebsiteID {
					wantID = 0
				}
				if haveID != wantID {
					t.Fatalf("Have Scope(%d) ScopeID: %d, Want: Scope(%d) ScopeID: %d", haveScp, haveID, scp, wantID)
				}

				//if ok := collisionCheck[haveHash]; ok {
				//	t.Fatalf("Collision Detected: %d", haveHash)
				//}
				//collisionCheck[haveHash] = true
			}
		}(scp)
	}
	wg.Wait()
	//t.Logf("[Info] Collision Map length: %d", len(collisionCheck))
}

var benchmarkHash scope.Hash
var benchmarkHashUnpackHash = scope.Hash(67112005)
var benchmarkHashUnpackScope scope.Scope
var benchmarkHashUnpackID int64

func BenchmarkHashPack(b *testing.B) {
	const want scope.Hash = 67112005
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		benchmarkHash = scope.NewHash(scope.StoreID, 3141)
	}
	if benchmarkHash != want {
		b.Fatalf("want %d have %d", want, benchmarkHash)
	}
}

func BenchmarkHashUnPack(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		benchmarkHashUnpackScope, benchmarkHashUnpackID = benchmarkHashUnpackHash.Unpack()
	}
	if benchmarkHashUnpackScope != scope.StoreID {
		b.Fatal("Expecting scope store")
	}
	if benchmarkHashUnpackID != 3141 {
		b.Fatal("Expecting ID 3141")
	}
}
