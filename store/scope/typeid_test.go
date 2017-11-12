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

	"github.com/corestoreio/pkg/store/scope"
	"github.com/corestoreio/errors"
	"github.com/stretchr/testify/assert"
)

var (
	_ fmt.Stringer            = (*scope.TypeID)(nil)
	_ fmt.GoStringer          = (*scope.TypeID)(nil)
	_ scope.RunModeCalculater = (*scope.TypeID)(nil)
)

var benchmarkTypeIDString string

func BenchmarkTypeIDString(b *testing.B) {
	s := scope.MakeTypeID(scope.Store, 33)
	b.ResetTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		benchmarkTypeIDString = s.String()
	}
}

func TestMakeTypeID(t *testing.T) {
	tests := []struct {
		scp     scope.Type
		id      int64
		wantScp scope.Type
		wantID  int64
	}{
		0: {scope.Default, 0, scope.Default, 0},
		1: {scope.Default, -1, scope.Default, 0},
		2: {scope.Default, 1, scope.Default, 0},
		3: {scope.Store, 1, scope.Store, 1},
		4: {scope.Group, 4, scope.Group, 4},
		5: {scope.Group, -4, scope.Absent, 0},
		6: {scope.Website, scope.MaxID, scope.Website, scope.MaxID},
		7: {scope.Website, -scope.MaxID, scope.Absent, 0},
		8: {scope.Website, scope.MaxID + 1, scope.Absent, 0},
	}
	for i, test := range tests {
		haveScp, haveID := scope.MakeTypeID(test.scp, test.id).Unpack()
		if have, want := haveScp, test.wantScp; have != want {
			t.Errorf("(IDX %d) Type Have: %v Want: %v", i, have, want)
		}
		if have, want := haveID, test.wantID; have != want {
			t.Errorf("(IDX %d) ID Have: %v Want: %v", i, have, want)
		}
	}
}

func TestTypeID_String(t *testing.T) {
	tests := []struct {
		h    scope.TypeID
		want string
	}{
		0: {scope.DefaultTypeID, "Type(Default) ID(0)"},
		1: {scope.MakeTypeID(scope.Website, 1), "Type(Website) ID(1)"},
		2: {scope.MakeTypeID(scope.Store, 2), "Type(Store) ID(2)"},
		3: {scope.MakeTypeID(scope.Group, 4), "Type(Group) ID(4)"},
		4: {scope.MakeTypeID(6, 5), "Type(Type(6)) ID(5)"},
		5: {0, "Type(Absent) ID(0)"},
	}
	for i, test := range tests {
		assert.Exactly(t, test.want, test.h.String(), "Index %d", i)
	}
}

func TestTypeID_GoString(t *testing.T) {
	tests := []struct {
		h    scope.TypeID
		want string
	}{
		0: {scope.DefaultTypeID, "scope.MakeTypeID(scope.Default, 0)"},
		1: {scope.MakeTypeID(scope.Website, 1), "scope.MakeTypeID(scope.Website, 1)"},
		2: {scope.MakeTypeID(scope.Store, 2), "scope.MakeTypeID(scope.Store, 2)"},
		3: {scope.MakeTypeID(scope.Group, 4), "scope.MakeTypeID(scope.Group, 4)"},
		4: {scope.MakeTypeID(6, 5), "scope.MakeTypeID(scope.Type(6), 5)"},
		5: {0, "scope.MakeTypeID(scope.Absent, 0)"},
	}
	for i, test := range tests {
		assert.Exactly(t, test.want, test.h.GoString(), "Index %d", i)
	}
}

func TestTypeIDSegment(t *testing.T) {

	tests := []struct {
		h    scope.TypeID
		want uint8
	}{
		0:  {scope.DefaultTypeID, 0},
		1:  {scope.MakeTypeID(scope.Type(0), 0), 0},
		2:  {scope.MakeTypeID(scope.Type(1), 0), 0},
		3:  {scope.MakeTypeID(scope.Default, -1), 0},
		4:  {scope.MakeTypeID(scope.Default, 1), 0},
		5:  {scope.MakeTypeID(scope.Default, 0), 0},
		6:  {scope.MakeTypeID(scope.Store, 0), 0},
		7:  {scope.MakeTypeID(scope.Store, 1), 1},
		8:  {scope.MakeTypeID(scope.Store, 2), 2},
		9:  {scope.MakeTypeID(scope.Store, 255), 255},
		10: {scope.MakeTypeID(scope.Store, 256), 0},
		11: {scope.MakeTypeID(scope.Store, 257), 1},
		12: {scope.MakeTypeID(scope.Store, scope.MaxID-1), 254},
		13: {scope.MakeTypeID(scope.Store, scope.MaxID), 255},
		14: {scope.MakeTypeID(scope.Store, scope.MaxID+1), 0},
		15: {scope.MakeTypeID(scope.Store, scope.MaxID+2), 0},
		16: {scope.MakeTypeID(scope.Store, -scope.MaxID), 0},
		17: {scope.MakeTypeID(scope.Type(7), 1), 1},
	}
	for i, test := range tests {
		if want, have := test.want, test.h.Segment(); want != have {
			t.Errorf("Index %03d: Want %03d Have %03d", i, want, have)
		}
	}
}

func TestMakeTypeIDError(t *testing.T) {

	h := scope.MakeTypeID(scope.Absent, -1)
	assert.Exactly(t, scope.TypeID(0), h)
}

func TestFromTypeIDError(t *testing.T) {

	scp, id := scope.TypeID(math.MaxUint32).Unpack()
	assert.Exactly(t, scope.Absent, scp)
	assert.Exactly(t, int64(-1), id)
}

func TestTypeIDValid(t *testing.T) {

	t.Logf("[Info] Max Store ID: %d", scope.MaxID)

	if testing.Short() {
		t.Skip("Skipping in short mode")
	}

	//var collisionCheck = make(map[scope.TypeID]bool) // just in case ...
	var wg sync.WaitGroup
	var scp = scope.Absent
	for ; scp < math.MaxUint8; scp++ {
		wg.Add(1)
		go func(theScp scope.Type) {
			defer wg.Done()
			for id := int64(0); id < scope.MaxID; id++ {
				haveTypeID := scope.MakeTypeID(theScp, id)

				haveScp, haveID := haveTypeID.Unpack()
				if haveScp != theScp {
					t.Fatalf("Have Type: %d, Want Type: %d", haveScp, theScp)
				}

				wantID := id
				if theScp < scope.Website {
					wantID = 0
				}
				if haveID != wantID {
					t.Fatalf("Have Type(%d) TypeID: %d, Want: Type(%d) TypeID: %d", haveScp, haveID, scp, wantID)
				}
				if haveTypeID > 0 && haveTypeID.ToUint64() < 1 { // stupid test
					t.Fatal("haveTypeID.ToUint64() cannot return zero")
				}

				//if ok := collisionCheck[haveTypeID]; ok {
				//	t.Fatalf("Collision Detected: %d", haveTypeID)
				//}
				//collisionCheck[haveTypeID] = true
			}
		}(scp)
	}
	wg.Wait()
	//t.Logf("[Info] Collision Map length: %d", len(collisionCheck))
}

func TestTypeID_EqualTypes(t *testing.T) {
	tests := []struct {
		h1        scope.TypeID
		h2        scope.TypeID
		wantEqual bool
	}{
		0: {0, 0, false},
		1: {0, scope.DefaultTypeID, false},
		2: {scope.DefaultTypeID, 0, false},
		3: {scope.DefaultTypeID, scope.DefaultTypeID, true},
		4: {scope.MakeTypeID(scope.Absent, 1), scope.MakeTypeID(scope.Absent, 1), false},
		5: {scope.MakeTypeID(scope.Store, scope.MaxID), scope.MakeTypeID(scope.Store, scope.MaxID), true},
		6: {scope.MakeTypeID(scope.Store, scope.MaxID), scope.MakeTypeID(scope.Store, scope.MaxID+1), false},
		7: {scope.MakeTypeID(scope.Store, scope.MaxID+1), scope.MakeTypeID(scope.Store, scope.MaxID), false},
		8: {scope.MakeTypeID(scope.Website, -1), scope.MakeTypeID(scope.Website, 1), false},
	}
	for i, test := range tests {
		if have, want := test.h1.EqualTypes(test.h2), test.wantEqual; have != want {
			t.Errorf("IDX %d Have: %v Want: %v", i, have, want)
		}
	}
}

func TestTypeID_Type(t *testing.T) {
	tests := []struct {
		h scope.TypeID
		s scope.Type
	}{
		0: {0, 0},
		1: {scope.DefaultTypeID, scope.Default},
		2: {scope.MakeTypeID(scope.Store, 1), scope.Store},
		3: {scope.MakeTypeID(254, 1), 254},
	}
	for i, test := range tests {
		if have, want := test.h.Type(), test.s; have != want {
			t.Errorf("IDX %d Have: %v Want: %v", i, have, want)
		}
	}
}

func TestTypeID_ID(t *testing.T) {
	tests := []struct {
		h  scope.TypeID
		id int64
	}{
		0: {0, 0},
		1: {scope.DefaultTypeID, 0},
		2: {scope.MakeTypeID(scope.Website, 33), 33},
		3: {scope.MakeTypeID(scope.Website, scope.MaxID), scope.MaxID},
		4: {scope.MakeTypeID(scope.Website, scope.MaxID+1), 0},
		5: {scope.TypeID(scope.Store)<<24 | scope.TypeID(scope.MaxID+1), -1},
	}
	for i, test := range tests {
		if have, want := test.h.ID(), test.id; have != want {
			t.Errorf("IDX %d Have: %v Want: %v", i, have, want)
		}
	}
}

var benchmarkTypeID scope.TypeID
var benchmarkTypeIDUnpackTypeID = scope.TypeID(67112005)
var benchmarkTypeIDUnpackType scope.Type
var benchmarkTypeIDUnpackID int64

func BenchmarkTypeIDPack(b *testing.B) {
	const want scope.TypeID = 67112005
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		benchmarkTypeID = scope.MakeTypeID(scope.Store, 3141)
	}
	if benchmarkTypeID != want {
		b.Fatalf("want %d have %d", want, benchmarkTypeID)
	}
}

func BenchmarkTypeIDUnPack(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		benchmarkTypeIDUnpackType, benchmarkTypeIDUnpackID = benchmarkTypeIDUnpackTypeID.Unpack()
	}
	if benchmarkTypeIDUnpackType != scope.Store {
		b.Fatal("Expecting scope store")
	}
	if benchmarkTypeIDUnpackID != 3141 {
		b.Fatal("Expecting ID 3141")
	}
}

var benchmarkTypeID_ValidParent bool

func BenchmarkTypeID_ValidParent(b *testing.B) {
	c := scope.MakeTypeID(scope.Store, 33)
	p := scope.MakeTypeID(scope.Website, 44)
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		benchmarkTypeID_ValidParent = c.ValidParent(p)
		if have, want := benchmarkTypeID_ValidParent, true; have != want {
			b.Errorf("Have: %v Want: %v", have, want)
		}
	}
}

func TestTypeIDes_Sort(t *testing.T) {
	hs := scope.TypeIDs{
		scope.MakeTypeID(scope.Store, 3),
		scope.MakeTypeID(scope.Website, 1),
		scope.DefaultTypeID,
		scope.MakeTypeID(scope.Store, 4),
		scope.MakeTypeID(scope.Website, 2),
	}
	sort.Sort(hs)
	assert.Exactly(t, scope.TypeIDs{0x1000000, 0x2000001, 0x2000002, 0x4000003, 0x4000004}, hs)
}

func TestTypeID_ValidParent(t *testing.T) {
	tests := []struct {
		c    scope.TypeID
		p    scope.TypeID
		want bool
	}{
		0:  {scope.DefaultTypeID, scope.DefaultTypeID, true},
		1:  {scope.MakeTypeID(scope.Website, 1), scope.DefaultTypeID, true},
		2:  {scope.MakeTypeID(scope.Website, 0), scope.DefaultTypeID, true},
		3:  {scope.MakeTypeID(scope.Store, 1), scope.MakeTypeID(scope.Website, 1), true},
		4:  {scope.MakeTypeID(scope.Store, -1), scope.MakeTypeID(scope.Website, 1), false},
		5:  {scope.MakeTypeID(scope.Store, 1), scope.MakeTypeID(scope.Website, -1), true},
		6:  {scope.MakeTypeID(scope.Store, 0), scope.MakeTypeID(scope.Website, 0), true},
		7:  {scope.DefaultTypeID, scope.MakeTypeID(scope.Website, 1), false},
		8:  {0, 0, true},
		9:  {0, scope.DefaultTypeID, false},
		10: {scope.DefaultTypeID, 0, true},
		11: {scope.MakeTypeID(scope.Website, 0), scope.Store.Pack(0), false},
		12: {scope.MakeTypeID(scope.Website, 1), scope.Store.Pack(2), false},
		13: {scope.MakeTypeID(scope.Store, 1), scope.Type(5).Pack(2), false},
	}
	for i, test := range tests {
		if have, want := test.c.ValidParent(test.p), test.want; have != want {
			t.Errorf("(%d) Have: %v Want: %v\nC: %d P: %d", i, have, want, test.c, test.p)
		}
	}
}

func TestTypeIDes_Lowest(t *testing.T) {
	tests := []struct {
		scope.TypeIDs
		wantTypeID scope.TypeID
		wantErrBhf errors.BehaviourFunc
	}{
		0:  {scope.TypeIDs{scope.Store.Pack(1)}, scope.Store.Pack(1), nil},
		1:  {scope.TypeIDs{scope.Store.Pack(1), scope.Store.Pack(2)}, 0, errors.IsNotValid},
		2:  {scope.TypeIDs{scope.Website.Pack(1), scope.Store.Pack(2)}, scope.Store.Pack(2), nil},
		3:  {scope.TypeIDs{scope.Website.Pack(1), scope.Store.Pack(2), scope.Store.Pack(2)}, scope.Store.Pack(2), nil},
		4:  {scope.TypeIDs{scope.Website.Pack(667), scope.Store.Pack(889), scope.Website.Pack(667), scope.Store.Pack(889)}, scope.Store.Pack(889), nil},
		5:  {scope.TypeIDs{scope.Store.Pack(2), scope.Website.Pack(345), scope.Website.Pack(346), scope.Store.Pack(2)}, scope.Store.Pack(2), nil},
		6:  {scope.TypeIDs{scope.Store.Pack(333), scope.Website.Pack(345), scope.Website.Pack(345), scope.Store.Pack(333)}, scope.Store.Pack(333), nil},
		7:  {scope.TypeIDs{scope.Store.Pack(3), scope.DefaultTypeID, scope.Store.Pack(3)}, scope.Store.Pack(3), nil},
		8:  {scope.TypeIDs{scope.Website.Pack(3), scope.DefaultTypeID, scope.Website.Pack(3)}, scope.Website.Pack(3), nil},
		9:  {scope.TypeIDs{scope.DefaultTypeID}, scope.DefaultTypeID, nil},
		10: {scope.TypeIDs{0, 1, 2}, scope.DefaultTypeID, nil},
		11: {nil, scope.DefaultTypeID, nil},
		12: {scope.TypeIDs{scope.MakeTypeID(55, 1), scope.MakeTypeID(55, 2), scope.MakeTypeID(56, 3)}, 0, errors.IsNotValid},
		13: {scope.TypeIDs{scope.MakeTypeID(55, 1), scope.MakeTypeID(55, 2), scope.MakeTypeID(55, 0)}, 0, errors.IsNotValid},
	}
	for i, test := range tests {
		haveTypeID, haveErr := test.TypeIDs.Lowest()
		assert.Exactly(t, test.wantTypeID, haveTypeID, "Index %d", i)
		if test.wantErrBhf != nil {
			assert.True(t, test.wantErrBhf(haveErr), "(IDX %d) %+v", i, haveErr)
			continue
		}
		assert.NoError(t, haveErr, "(IDX %d) %+v", i, haveErr)
	}
}

func TestTypeIDs_TargetAndParents(t *testing.T) {
	tests := []struct {
		scope.TypeIDs
		wantTarget  scope.TypeID
		wantParents scope.TypeIDs
	}{
		0: {scope.TypeIDs{scope.DefaultTypeID, scope.DefaultTypeID}, scope.DefaultTypeID, scope.TypeIDs{scope.DefaultTypeID}},
		1: {scope.TypeIDs{scope.DefaultTypeID, scope.DefaultTypeID, scope.DefaultTypeID}, scope.DefaultTypeID, scope.TypeIDs{scope.DefaultTypeID}},
		2: {scope.TypeIDs{scope.DefaultTypeID, scope.Website.Pack(3), scope.Store.Pack(2), scope.DefaultTypeID}, scope.DefaultTypeID, scope.TypeIDs{scope.DefaultTypeID}},
		3: {scope.TypeIDs{scope.Store.Pack(1)}, scope.Store.Pack(1), scope.TypeIDs{scope.DefaultTypeID}},
		4: {scope.TypeIDs{scope.Website.Pack(2), scope.Store.Pack(1)}, scope.Website.Pack(2), scope.TypeIDs{scope.DefaultTypeID}},
		5: {scope.TypeIDs{scope.Store.Pack(1), scope.Website.Pack(2)}, scope.Store.Pack(1), scope.TypeIDs{scope.Website.Pack(2), scope.DefaultTypeID}},
		6: {scope.TypeIDs{scope.Group.Pack(1), scope.Website.Pack(2)}, scope.Group.Pack(1), scope.TypeIDs{scope.Website.Pack(2), scope.DefaultTypeID}},
		7: {scope.TypeIDs{scope.DefaultTypeID, scope.Group.Pack(1), scope.Website.Pack(2)}, scope.DefaultTypeID, scope.TypeIDs{scope.DefaultTypeID}},
		8: {nil, scope.DefaultTypeID, scope.TypeIDs{scope.DefaultTypeID}},
		9: {scope.TypeIDs{scope.DefaultTypeID}, scope.DefaultTypeID, scope.TypeIDs{scope.DefaultTypeID}},
	}
	for i, test := range tests {
		haveTarget, haveParents := test.TypeIDs.TargetAndParents()
		assert.Exactly(t, test.wantTarget, haveTarget, "Index %d", i)
		assert.Exactly(t, test.wantParents, haveParents, "Index %d", i)
	}
}

func BenchmarkTypeIDs_TargetAndParents(b *testing.B) {
	ids := scope.TypeIDs{scope.Group.Pack(1), scope.Website.Pack(2)}
	target := scope.Group.Pack(1)
	parents := scope.TypeIDs{scope.Website.Pack(2), scope.DefaultTypeID}
	var haveT scope.TypeID
	var haveP scope.TypeIDs
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		haveT, haveP = ids.TargetAndParents()
	}
	if haveT != target {
		b.Fatalf("Have %s Want %s", haveT, target)
	}
	assert.Exactly(b, parents, haveP)
}

func TestTypeIDs_String(t *testing.T) {
	tests := []struct {
		scope.TypeIDs
		want string
	}{
		0: {scope.TypeIDs{scope.Store.Pack(1)}, "Type(Store) ID(1)"},
		1: {scope.TypeIDs{scope.Website.Pack(2), scope.Store.Pack(1)}, "Type(Website) ID(2); Type(Store) ID(1)"},
		2: {scope.TypeIDs{scope.DefaultTypeID, scope.Group.Pack(1), scope.Website.Pack(2)}, "Type(Default) ID(0); Type(Group) ID(1); Type(Website) ID(2)"},
		3: {nil, ""},
		4: {scope.TypeIDs{scope.DefaultTypeID}, "Type(Default) ID(0)"},
	}
	for i, test := range tests {
		assert.Exactly(t, test.want, test.TypeIDs.String(), "Index %d", i)
	}
}
