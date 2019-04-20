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

package scope_test

import (
	"encoding"
	"fmt"
	"math"
	"sort"
	"sync"
	"testing"

	"github.com/corestoreio/errors"
	"github.com/corestoreio/pkg/store/scope"
	"github.com/corestoreio/pkg/util/assert"
)

var (
	_ fmt.Stringer               = (*scope.TypeID)(nil)
	_ fmt.GoStringer             = (*scope.TypeID)(nil)
	_ encoding.TextMarshaler     = (*scope.TypeID)(nil)
	_ encoding.TextUnmarshaler   = (*scope.TypeID)(nil)
	_ encoding.BinaryMarshaler   = (*scope.TypeID)(nil)
	_ encoding.BinaryUnmarshaler = (*scope.TypeID)(nil)
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
		id      uint32
		wantScp scope.Type
		wantID  uint32
	}{
		{scope.Default, 0, scope.Default, 0},
		{scope.Default, 1, scope.Default, 0},
		{scope.Store, 1, scope.Store, 1},
		{scope.Group, 4, scope.Group, 4},
		{scope.Website, scope.MaxID, scope.Website, scope.MaxID},
		{scope.Website, scope.MaxID + 1, scope.Absent, 0},
	}
	for i, test := range tests {
		t.Logf("ID: %d", scope.MakeTypeID(test.scp, test.id))
		haveScp, haveID := scope.MakeTypeID(test.scp, test.id).Unpack()
		if have, want := haveScp, test.wantScp; have != want {
			t.Fatalf("(IDX %d) Type Have: %v Want: %v", i, have, want)
		}
		if have, want := haveID, test.wantID; have != want {
			t.Fatalf("(IDX %d) ID Have: %v Want: %v", i, have, want)
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
		3:  {scope.MakeTypeID(scope.Default, 0), 0},
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
		16: {scope.MakeTypeID(scope.Type(7), 1), 1},
	}
	for i, test := range tests {
		if want, have := test.want, test.h.Segment(); want != have {
			t.Errorf("Index %03d: Want %03d Have %03d", i, want, have)
		}
	}
}

func TestMakeTypeIDError(t *testing.T) {

	h := scope.MakeTypeID(scope.Absent, 0)
	assert.Exactly(t, scope.TypeID(0), h)
}

func TestFromTypeIDError(t *testing.T) {

	scp, id := scope.TypeID(math.MaxUint32).Unpack()
	assert.Exactly(t, scope.Absent, scp)
	assert.Exactly(t, uint32(0), id)
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
			for id := uint32(0); id < scope.MaxID; id++ {
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
		8: {scope.MakeTypeID(scope.Store, 0), scope.MakeTypeID(scope.Website, 1), false},
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
		h       scope.TypeID
		id      uint32
		wantErr errors.Kind
	}{
		0: {0, 0, errors.NoKind},
		1: {scope.DefaultTypeID, 0, errors.NoKind},
		2: {scope.MakeTypeID(scope.Website, 33), 33, errors.NoKind},
		3: {scope.MakeTypeID(scope.Website, scope.MaxID), scope.MaxID, errors.NoKind},
		4: {scope.MakeTypeID(scope.Website, scope.MaxID+1), 0, errors.NoKind},
		5: {scope.TypeID(scope.Store)<<24 | scope.TypeID(scope.MaxID+1), 0, errors.Overflowed},
	}
	for i, test := range tests {
		hID, err := test.h.ID()
		if test.wantErr > errors.NoKind {
			assert.True(t, test.wantErr.Match(err), "IDX %d", i)
			continue
		}
		assert.NoError(t, err, "IDX %d", i)
		if have, want := hID, test.id; have != want {
			t.Errorf("IDX %d Have: %v Want: %v", i, have, want)
		}
	}
}

var benchmarkTypeID scope.TypeID
var benchmarkTypeIDUnpackTypeID = scope.TypeID(67112005)
var benchmarkTypeIDUnpackType scope.Type
var benchmarkTypeIDUnpackID uint32

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
		{scope.DefaultTypeID, scope.DefaultTypeID, true},
		{scope.MakeTypeID(scope.Website, 1), scope.DefaultTypeID, true},
		{scope.MakeTypeID(scope.Website, 0), scope.DefaultTypeID, true},
		{scope.MakeTypeID(scope.Store, 1), scope.MakeTypeID(scope.Website, 1), true},
		{scope.MakeTypeID(scope.Store, 1), scope.MakeTypeID(scope.Website, 0), true},
		{scope.MakeTypeID(scope.Store, 0), scope.MakeTypeID(scope.Website, 0), true},
		{scope.DefaultTypeID, scope.MakeTypeID(scope.Website, 1), false},
		{0, 0, true},
		{0, scope.DefaultTypeID, false},
		{scope.DefaultTypeID, 0, true},
		{scope.MakeTypeID(scope.Website, 0), scope.Store.WithID(0), false},
		{scope.MakeTypeID(scope.Website, 1), scope.Store.WithID(2), false},
		{scope.MakeTypeID(scope.Store, 1), scope.Type(5).WithID(2), false},
	}
	for i, test := range tests {
		if have, want := test.c.ValidParent(test.p), test.want; have != want {
			t.Fatalf("(%d) Have: %v Want: %v\nC: %d P: %d", i, have, want, test.c, test.p)
		}
	}
}

func TestTypeIDes_Lowest(t *testing.T) {
	tests := []struct {
		scope.TypeIDs
		wantTypeID  scope.TypeID
		wantErrKind errors.Kind
	}{
		0:  {scope.TypeIDs{scope.Store.WithID(1)}, scope.Store.WithID(1), errors.NoKind},
		1:  {scope.TypeIDs{scope.Store.WithID(1), scope.Store.WithID(2)}, 0, errors.NotValid},
		2:  {scope.TypeIDs{scope.Website.WithID(1), scope.Store.WithID(2)}, scope.Store.WithID(2), errors.NoKind},
		3:  {scope.TypeIDs{scope.Website.WithID(1), scope.Store.WithID(2), scope.Store.WithID(2)}, scope.Store.WithID(2), errors.NoKind},
		4:  {scope.TypeIDs{scope.Website.WithID(667), scope.Store.WithID(889), scope.Website.WithID(667), scope.Store.WithID(889)}, scope.Store.WithID(889), errors.NoKind},
		5:  {scope.TypeIDs{scope.Store.WithID(2), scope.Website.WithID(345), scope.Website.WithID(346), scope.Store.WithID(2)}, scope.Store.WithID(2), errors.NoKind},
		6:  {scope.TypeIDs{scope.Store.WithID(333), scope.Website.WithID(345), scope.Website.WithID(345), scope.Store.WithID(333)}, scope.Store.WithID(333), errors.NoKind},
		7:  {scope.TypeIDs{scope.Store.WithID(3), scope.DefaultTypeID, scope.Store.WithID(3)}, scope.Store.WithID(3), errors.NoKind},
		8:  {scope.TypeIDs{scope.Website.WithID(3), scope.DefaultTypeID, scope.Website.WithID(3)}, scope.Website.WithID(3), errors.NoKind},
		9:  {scope.TypeIDs{scope.DefaultTypeID}, scope.DefaultTypeID, errors.NoKind},
		10: {scope.TypeIDs{0, 1, 2}, scope.DefaultTypeID, errors.NoKind},
		11: {nil, scope.DefaultTypeID, errors.NoKind},
		12: {scope.TypeIDs{scope.MakeTypeID(55, 1), scope.MakeTypeID(55, 2), scope.MakeTypeID(56, 3)}, 0, errors.NotValid},
		13: {scope.TypeIDs{scope.MakeTypeID(55, 1), scope.MakeTypeID(55, 2), scope.MakeTypeID(55, 0)}, 0, errors.NotValid},
	}
	for i, test := range tests {
		haveTypeID, haveErr := test.TypeIDs.Lowest()
		assert.Exactly(t, test.wantTypeID, haveTypeID, "Index %d", i)
		if !test.wantErrKind.Empty() {
			assert.True(t, test.wantErrKind.Match(haveErr), "(IDX %d) %+v", i, haveErr)
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
		2: {scope.TypeIDs{scope.DefaultTypeID, scope.Website.WithID(3), scope.Store.WithID(2), scope.DefaultTypeID}, scope.DefaultTypeID, scope.TypeIDs{scope.DefaultTypeID}},
		3: {scope.TypeIDs{scope.Store.WithID(1)}, scope.Store.WithID(1), scope.TypeIDs{scope.DefaultTypeID}},
		4: {scope.TypeIDs{scope.Website.WithID(2), scope.Store.WithID(1)}, scope.Website.WithID(2), scope.TypeIDs{scope.DefaultTypeID}},
		5: {scope.TypeIDs{scope.Store.WithID(1), scope.Website.WithID(2)}, scope.Store.WithID(1), scope.TypeIDs{scope.Website.WithID(2), scope.DefaultTypeID}},
		6: {scope.TypeIDs{scope.Group.WithID(1), scope.Website.WithID(2)}, scope.Group.WithID(1), scope.TypeIDs{scope.Website.WithID(2), scope.DefaultTypeID}},
		7: {scope.TypeIDs{scope.DefaultTypeID, scope.Group.WithID(1), scope.Website.WithID(2)}, scope.DefaultTypeID, scope.TypeIDs{scope.DefaultTypeID}},
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
	ids := scope.TypeIDs{scope.Group.WithID(1), scope.Website.WithID(2)}
	target := scope.Group.WithID(1)
	parents := scope.TypeIDs{scope.Website.WithID(2), scope.DefaultTypeID}
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
		0: {scope.TypeIDs{scope.Store.WithID(1)}, "Type(Store) ID(1)"},
		1: {scope.TypeIDs{scope.Website.WithID(2), scope.Store.WithID(1)}, "Type(Website) ID(2); Type(Store) ID(1)"},
		2: {scope.TypeIDs{scope.DefaultTypeID, scope.Group.WithID(1), scope.Website.WithID(2)}, "Type(Default) ID(0); Type(Group) ID(1); Type(Website) ID(2)"},
		3: {nil, ""},
		4: {scope.TypeIDs{scope.DefaultTypeID}, "Type(Default) ID(0)"},
	}
	for i, test := range tests {
		assert.Exactly(t, test.want, test.TypeIDs.String(), "Index %d", i)
	}
}

func TestTypeID_Marshal(t *testing.T) {
	id := scope.MakeTypeID(scope.Store, 2)

	t.Run("binary", func(t *testing.T) {
		bin, err := id.MarshalBinary()
		assert.NoError(t, err)
		var idBin scope.TypeID
		assert.NoError(t, idBin.UnmarshalBinary(bin))
		assert.Exactly(t, id, idBin)

	})
	t.Run("text", func(t *testing.T) {
		txt, err := id.MarshalText()
		assert.NoError(t, err)

		var idTxt scope.TypeID
		assert.NoError(t, idTxt.UnmarshalText(txt))
		assert.Exactly(t, id, idTxt)
	})
	t.Run("UnmarshalText error", func(t *testing.T) {
		var idTxt scope.TypeID
		assert.EqualError(t, idTxt.UnmarshalText([]byte(`-1`)), "[scope] TypeID.UnmarshalText with text \"-1\": strconv.ParseUint: parsing \"-1\": invalid syntax")
	})
	t.Run("UnmarshalBinary error", func(t *testing.T) {
		var idTxt scope.TypeID
		assert.True(t, errors.BadEncoding.Match(idTxt.UnmarshalBinary([]byte(`x`))), "Error should have kind BadEncoding")
	})
	t.Run("ToIntString", func(t *testing.T) {
		assert.Exactly(t, "67108866", id.ToIntString())
	})
	t.Run("AppendBytes", func(t *testing.T) {
		data := []byte(`ID:`)
		assert.Exactly(t, []byte(`ID:67108866`), id.AppendBytes(data))
	})
}

func TestMakeTypeIDString(t *testing.T) {
	tid, err := scope.MakeTypeIDString("-1")
	assert.Exactly(t, scope.TypeID(0), tid)
	assert.EqualError(t, err, "[scope] MakeTypeIDString with text \"-1\": strconv.ParseUint: parsing \"-1\": invalid syntax")
}

func TestTypeID_IsValid(t *testing.T) {
	assert.NoError(t, scope.DefaultTypeID.IsValid())
	assert.NoError(t, scope.Website.WithID(3).IsValid())
	assert.NoError(t, scope.Group.WithID(4).IsValid())
	assert.NoError(t, scope.Store.WithID(5).IsValid())
	assert.True(t, errors.NotValid.Match(scope.TypeID(0).IsValid()))
	assert.True(t, errors.NotValid.Match(scope.TypeID(485968409).IsValid()))
}

func TestTypeID_AppendHuman(t *testing.T) {
	tests := []struct {
		sid  scope.TypeID
		want string
	}{
		{scope.DefaultTypeID, ""},
		{scope.Group.WithID(13), ""},
		{scope.Website.WithID(13), "websites/13"},
		{scope.Website.WithID(0), "websites/0"},
		{scope.Store.WithID(13), "stores/13"},
		{scope.Store.WithID(0), "stores/0"},
	}
	for _, test := range tests {
		assert.Exactly(t, test.want, string(test.sid.AppendHuman(nil, '/')))
	}
}
