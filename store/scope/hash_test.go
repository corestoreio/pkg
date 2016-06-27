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

	"sync"

	"github.com/corestoreio/csfw/store/scope"
	"github.com/stretchr/testify/assert"
)

var benchmarkHashString string

func BenchmarkHashString(b *testing.B) {
	s := scope.NewHash(scope.Store, 33)
	b.ResetTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		benchmarkHashString = s.String()
	}
}

func TestNewHash(t *testing.T) {
	tests := []struct {
		scp     scope.Scope
		id      int64
		wantScp scope.Scope
		wantID  int64
	}{
		{scope.Default, 0, scope.Default, 0},
		{scope.Default, -1, scope.Default, 0},
		{scope.Default, 1, scope.Default, 0},
		{scope.Store, 1, scope.Store, 1},
		{scope.Group, 4, scope.Group, 4},
		{scope.Group, -4, scope.Absent, 0},
		{scope.Website, scope.MaxStoreID, scope.Website, scope.MaxStoreID},
		{scope.Website, -scope.MaxStoreID, scope.Absent, 0},
		{scope.Website, scope.MaxStoreID + 1, scope.Absent, 0},
	}
	for i, test := range tests {
		haveScp, haveID := scope.NewHash(test.scp, test.id).Unpack()
		if have, want := haveScp, test.wantScp; have != want {
			t.Errorf("(IDX %d) Scope Have: %v Want: %v", i, have, want)
		}
		if have, want := haveID, test.wantID; have != want {
			t.Errorf("(IDX %d) ID Have: %v Want: %v", i, have, want)
		}
	}
}

func TestHashString(t *testing.T) {
	s := scope.NewHash(scope.Store, 33).String()
	assert.Exactly(t, "Scope(Store) ID(33)", s)
}

func TestHashSegment(t *testing.T) {

	tests := []struct {
		h    scope.Hash
		want uint8
	}{
		{scope.DefaultHash, 0},
		{scope.NewHash(scope.Scope(0), 0), 0},
		{scope.NewHash(scope.Scope(1), 0), 0},
		{scope.NewHash(scope.Default, -1), 0},
		{scope.NewHash(scope.Default, 1), 0},
		{scope.NewHash(scope.Default, 0), 0},
		{scope.NewHash(scope.Store, 0), 0},
		{scope.NewHash(scope.Store, 1), 1},
		{scope.NewHash(scope.Store, 2), 2},
		{scope.NewHash(scope.Store, 255), 255},
		{scope.NewHash(scope.Store, 256), 0},
		{scope.NewHash(scope.Store, 257), 1},
		{scope.NewHash(scope.Store, scope.MaxStoreID-1), 254},
		{scope.NewHash(scope.Store, scope.MaxStoreID), 255},
		{scope.NewHash(scope.Store, scope.MaxStoreID+1), 0},
		{scope.NewHash(scope.Store, scope.MaxStoreID+2), 0},
		{scope.NewHash(scope.Store, -scope.MaxStoreID), 0},
		{scope.NewHash(scope.Scope(7), 1), 1},
	}
	for i, test := range tests {
		if want, have := test.want, test.h.Segment(); want != have {
			t.Errorf("Index %03d: Want %03d Have %03d", i, want, have)
		}
	}
}

func TestNewHashError(t *testing.T) {

	h := scope.NewHash(scope.Absent, -1)
	assert.Exactly(t, scope.Hash(0), h)
}

func TestFromHashError(t *testing.T) {

	scp, id := scope.Hash(math.MaxUint32).Unpack()
	assert.Exactly(t, scope.Absent, scp)
	assert.Exactly(t, int64(-1), id)
}

func TestHashValid(t *testing.T) {

	t.Logf("[Info] Max Store ID: %d", scope.MaxStoreID)

	if testing.Short() {
		t.Skip("Skipping in short mode")
	}

	//var collisionCheck = make(map[scope.Hash]bool) // just in case ...
	var wg sync.WaitGroup
	var scp = scope.Absent
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
				if theScp < scope.Website {
					wantID = 0
				}
				if haveID != wantID {
					t.Fatalf("Have Scope(%d) ScopeID: %d, Want: Scope(%d) ScopeID: %d", haveScp, haveID, scp, wantID)
				}
				if haveHash > 0 && haveHash.ToUint64() < 1 { // stupid test
					t.Fatal("haveHash.ToUint64() cannot return zero")
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
		benchmarkHash = scope.NewHash(scope.Store, 3141)
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
	if benchmarkHashUnpackScope != scope.Store {
		b.Fatal("Expecting scope store")
	}
	if benchmarkHashUnpackID != 3141 {
		b.Fatal("Expecting ID 3141")
	}
}
