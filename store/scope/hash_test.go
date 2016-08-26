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
	"fmt"
	"math"
	"sort"
	"sync"
	"testing"

	"github.com/corestoreio/csfw/store/scope"
	"github.com/corestoreio/csfw/util/errors"
	"github.com/stretchr/testify/assert"
)

var (
	_ fmt.Stringer            = (*scope.Hash)(nil)
	_ fmt.GoStringer          = (*scope.Hash)(nil)
	_ scope.RunModeCalculater = (*scope.Hash)(nil)
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

func TestHash_String(t *testing.T) {
	tests := []struct {
		h    scope.Hash
		want string
	}{
		{scope.DefaultHash, "Scope(Default) ID(0)"},
		{scope.NewHash(scope.Website, 1), "Scope(Website) ID(1)"},
		{scope.NewHash(scope.Store, 2), "Scope(Store) ID(2)"},
		{scope.NewHash(scope.Group, 4), "Scope(Group) ID(4)"},
		{scope.NewHash(6, 5), "Scope(Scope(6)) ID(5)"},
		{0, "Scope(Absent) ID(0)"},
	}
	for i, test := range tests {
		assert.Exactly(t, test.want, test.h.String(), "Index %d", i)
	}
}

func TestHash_GoString(t *testing.T) {
	tests := []struct {
		h    scope.Hash
		want string
	}{
		{scope.DefaultHash, "scope.NewHash(scope.Default, 0)"},
		{scope.NewHash(scope.Website, 1), "scope.NewHash(scope.Website, 1)"},
		{scope.NewHash(scope.Store, 2), "scope.NewHash(scope.Store, 2)"},
		{scope.NewHash(scope.Group, 4), "scope.NewHash(scope.Group, 4)"},
		{scope.NewHash(6, 5), "scope.NewHash(scope.Scope(6), 5)"},
		{0, "scope.NewHash(scope.Absent, 0)"},
	}
	for i, test := range tests {
		assert.Exactly(t, test.want, test.h.GoString(), "Index %d", i)
	}
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

func TestHash_EqualScope(t *testing.T) {
	tests := []struct {
		h1        scope.Hash
		h2        scope.Hash
		wantEqual bool
	}{
		{0, 0, false},
		{0, scope.DefaultHash, false},
		{scope.DefaultHash, 0, false},
		{scope.DefaultHash, scope.DefaultHash, true},
		{scope.NewHash(scope.Absent, 1), scope.NewHash(scope.Absent, 1), false},
		{scope.NewHash(scope.Store, scope.MaxStoreID), scope.NewHash(scope.Store, scope.MaxStoreID), true},
		{scope.NewHash(scope.Store, scope.MaxStoreID), scope.NewHash(scope.Store, scope.MaxStoreID+1), false},
		{scope.NewHash(scope.Store, scope.MaxStoreID+1), scope.NewHash(scope.Store, scope.MaxStoreID), false},
		{scope.NewHash(scope.Website, -1), scope.NewHash(scope.Website, 1), false},
	}
	for i, test := range tests {
		if have, want := test.h1.EqualScope(test.h2), test.wantEqual; have != want {
			t.Errorf("IDX %d Have: %v Want: %v", i, have, want)
		}
	}
}

func TestHash_Scope(t *testing.T) {
	tests := []struct {
		h scope.Hash
		s scope.Scope
	}{
		{0, 0},
		{scope.DefaultHash, scope.Default},
		{scope.NewHash(scope.Store, 1), scope.Store},
		{scope.NewHash(254, 1), 254},
	}
	for i, test := range tests {
		if have, want := test.h.Scope(), test.s; have != want {
			t.Errorf("IDX %d Have: %v Want: %v", i, have, want)
		}
	}
}

func TestHash_ID(t *testing.T) {
	tests := []struct {
		h  scope.Hash
		id int64
	}{
		{0, 0},
		{scope.DefaultHash, 0},
		{scope.NewHash(scope.Website, 33), 33},
		{scope.NewHash(scope.Website, scope.MaxStoreID), scope.MaxStoreID},
		{scope.NewHash(scope.Website, scope.MaxStoreID+1), 0},
		{scope.Hash(scope.Store)<<24 | scope.Hash(scope.MaxStoreID+1), -1},
	}
	for i, test := range tests {
		if have, want := test.h.ID(), test.id; have != want {
			t.Errorf("IDX %d Have: %v Want: %v", i, have, want)
		}
	}
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

var benchmarkHash_ValidParent bool

func BenchmarkHash_ValidParent(b *testing.B) {
	c := scope.NewHash(scope.Store, 33)
	p := scope.NewHash(scope.Website, 44)
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		benchmarkHash_ValidParent = c.ValidParent(p)
		if have, want := benchmarkHash_ValidParent, true; have != want {
			b.Errorf("Have: %v Want: %v", have, want)
		}
	}
}

func TestHashes_Sort(t *testing.T) {
	hs := scope.Hashes{
		scope.NewHash(scope.Store, 3),
		scope.NewHash(scope.Website, 1),
		scope.DefaultHash,
		scope.NewHash(scope.Store, 4),
		scope.NewHash(scope.Website, 2),
	}
	sort.Sort(hs)
	assert.Exactly(t, scope.Hashes{0x1000000, 0x2000001, 0x2000002, 0x4000003, 0x4000004}, hs)
}

func TestHash_ValidParent(t *testing.T) {
	tests := []struct {
		c    scope.Hash
		p    scope.Hash
		want bool
	}{
		{scope.DefaultHash, scope.DefaultHash, true},
		{scope.NewHash(scope.Website, 1), scope.DefaultHash, true},
		{scope.NewHash(scope.Website, 0), scope.DefaultHash, true},
		{scope.NewHash(scope.Store, 1), scope.NewHash(scope.Website, 1), true},
		{scope.NewHash(scope.Store, -1), scope.NewHash(scope.Website, 1), false},
		{scope.NewHash(scope.Store, 1), scope.NewHash(scope.Website, -1), false},
		{scope.NewHash(scope.Store, 0), scope.NewHash(scope.Website, 0), true},
		{scope.DefaultHash, scope.NewHash(scope.Website, 1), false},
		{0, 0, false},
		{0, scope.DefaultHash, false},
		{scope.DefaultHash, 0, false},
	}
	for i, test := range tests {
		if have, want := test.c.ValidParent(test.p), test.want; have != want {
			t.Errorf("(%d) Have: %v Want: %v", i, have, want)
		}
	}
}

func TestHashes_Lowest(t *testing.T) {
	tests := []struct {
		scope.Hashes
		wantHash   scope.Hash
		wantErrBhf errors.BehaviourFunc
	}{
		{scope.Hashes{scope.Store.ToHash(1)}, scope.Store.ToHash(1), nil},
		{scope.Hashes{scope.Store.ToHash(1), scope.Store.ToHash(2)}, 0, errors.IsNotValid},
		{scope.Hashes{scope.Website.ToHash(1), scope.Store.ToHash(2)}, scope.Store.ToHash(2), nil},
		{scope.Hashes{scope.Website.ToHash(1), scope.Store.ToHash(2), scope.Store.ToHash(2)}, scope.Store.ToHash(2), nil},
		{scope.Hashes{scope.Website.ToHash(667), scope.Store.ToHash(889), scope.Website.ToHash(667), scope.Store.ToHash(889)}, scope.Store.ToHash(889), nil},
		{scope.Hashes{scope.Store.ToHash(2), scope.Website.ToHash(345), scope.Website.ToHash(346), scope.Store.ToHash(2)}, scope.Store.ToHash(2), nil},
		{scope.Hashes{scope.Store.ToHash(333), scope.Website.ToHash(345), scope.Website.ToHash(345), scope.Store.ToHash(333)}, scope.Store.ToHash(333), nil},
		{scope.Hashes{scope.Store.ToHash(3), scope.DefaultHash, scope.Store.ToHash(3)}, scope.Store.ToHash(3), nil},
		{scope.Hashes{scope.Website.ToHash(3), scope.DefaultHash, scope.Website.ToHash(3)}, scope.Website.ToHash(3), nil},
		{scope.Hashes{scope.DefaultHash}, scope.DefaultHash, nil},
		{scope.Hashes{0, 1, 2}, scope.DefaultHash, nil},
		{nil, scope.DefaultHash, nil},
		{scope.Hashes{scope.NewHash(55, 1), scope.NewHash(55, 2), scope.NewHash(56, 3)}, 0, errors.IsNotValid},
		{scope.Hashes{scope.NewHash(55, 1), scope.NewHash(55, 2), scope.NewHash(55, 0)}, 0, errors.IsNotValid},
	}
	for i, test := range tests {
		haveHash, haveErr := test.Hashes.Lowest()
		assert.Exactly(t, test.wantHash, haveHash, "Index %d", i)
		if test.wantErrBhf != nil {
			assert.True(t, test.wantErrBhf(haveErr), "(IDX %d) %+v", i, haveErr)
			continue
		}
		assert.NoError(t, haveErr, "(IDX %d) %+v", i, haveErr)
	}
}
