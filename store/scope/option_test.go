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
	"testing"

	"github.com/corestoreio/csfw/store/scope"
	"github.com/stretchr/testify/assert"
	"hash/fnv"
)

func TestMustSetByCode(t *testing.T) {
	t.Parallel()
	defer func() {
		if r := recover(); r != nil {
			assert.EqualError(t, r.(error), scope.ErrUnsupportedScopeID.Error())
		} else {
			t.Fatal("Expecting a panic")
		}
	}()
	scope.MustSetByCode(99, "Gopher")
}

func TestMustSetByID(t *testing.T) {
	t.Parallel()
	defer func() {
		if r := recover(); r != nil {
			assert.EqualError(t, r.(error), scope.ErrUnsupportedScopeID.Error())
		} else {
			t.Fatal("Expecting a panic")
		}
	}()
	scope.MustSetByID(99, 444)
}

func TestApplyCode(t *testing.T) {
	t.Parallel()
	tests := []struct {
		wantStoreCode   string
		wantWebsiteCode string
		haveCode        string
		s               scope.Scope
		err             error
	}{
		{"", "de1", "de1", scope.WebsiteID, nil},
		{"de2", "", "de2", scope.StoreID, nil},
		{"", "", "de3", scope.GroupID, scope.ErrUnsupportedScopeID},
		{"", "", "de4", scope.AbsentID, scope.ErrUnsupportedScopeID},
	}

	for _, test := range tests {
		so, err := scope.SetByCode(test.s, test.haveCode)
		assert.NotNil(t, so)
		if test.err != nil {
			assert.EqualError(t, err, test.err.Error())
		} else {
			assert.NoError(t, err)
			assert.Equal(t, test.s, so.Scope())
			assert.Equal(t, test.wantStoreCode, so.StoreCode())
			assert.Equal(t, test.wantWebsiteCode, so.WebsiteCode())
		}
	}
}

func TestApplyID(t *testing.T) {
	t.Parallel()
	tests := []struct {
		wantWebsiteID scope.WebsiteIDer
		wantGroupID   scope.GroupIDer
		wantStoreID   scope.StoreIDer

		haveID int64
		s      scope.Scope
		err    error
	}{
		{scope.MockID(1), nil, nil, 1, scope.WebsiteID, nil},
		{nil, scope.MockID(3), nil, 3, scope.GroupID, nil},
		{nil, nil, scope.MockID(2), 2, scope.StoreID, nil},
		{nil, nil, nil, 4, scope.AbsentID, scope.ErrUnsupportedScopeID},
	}

	for _, test := range tests {
		so, err := scope.SetByID(test.s, test.haveID)
		assert.NotNil(t, so)
		if test.err != nil {
			assert.EqualError(t, err, test.err.Error())
			assert.Nil(t, so.Website)
			assert.Nil(t, so.Group)
			assert.Nil(t, so.Store)
		} else {
			assert.NoError(t, err)
			assert.Equal(t, test.s, so.Scope())
			assert.Equal(t, "", so.StoreCode())
			assert.Equal(t, "", so.WebsiteCode())

			if test.wantWebsiteID != nil {
				assert.Equal(t, test.wantWebsiteID.WebsiteID(), so.Website.WebsiteID())
			} else {
				assert.Nil(t, test.wantWebsiteID)
			}

			if test.wantGroupID != nil {
				assert.Equal(t, test.wantGroupID.GroupID(), so.Group.GroupID())
			} else {
				assert.Nil(t, test.wantGroupID)
			}

			if test.wantStoreID != nil {
				assert.Equal(t, test.wantStoreID.StoreID(), so.Store.StoreID())
			} else {
				assert.Nil(t, test.wantStoreID)
			}
		}
	}
}

func TestApplyWebsite(t *testing.T) {
	t.Parallel()
	so := scope.Option{Website: scope.MockID(3)}
	assert.NotNil(t, so)
	assert.Equal(t, int64(3), so.Website.WebsiteID())
	assert.Nil(t, so.Group)
	assert.Nil(t, so.Store)
	assert.Exactly(t, scope.WebsiteID.String(), so.String())
}

func TestApplyGroup(t *testing.T) {
	t.Parallel()
	so := scope.Option{Group: scope.MockID(3)}
	assert.NotNil(t, so)
	assert.Equal(t, int64(3), so.Group.GroupID())
	assert.Nil(t, so.Website)
	assert.Nil(t, so.Store)
	assert.Exactly(t, scope.GroupID.String(), so.String())
}

func TestApplyStore(t *testing.T) {
	t.Parallel()
	so := scope.Option{Store: scope.MockID(3)}
	assert.NotNil(t, so)
	assert.Equal(t, int64(3), so.Store.StoreID())
	assert.Nil(t, so.Website)
	assert.Nil(t, so.Group)
	assert.Exactly(t, scope.StoreID.String(), so.String())
}

func TestApplyDefault(t *testing.T) {
	t.Parallel()
	so := scope.Option{}
	assert.NotNil(t, so)
	assert.Exactly(t, scope.DefaultID, so.Scope())
}

func TestToUint32ByID(t *testing.T) {
	tests := []struct {
		so   scope.Option
		want uint32
	}{
		{scope.Option{}, 0},
		{scope.MustSetByID(scope.WebsiteID, 11), 11},
		{scope.MustSetByID(scope.GroupID, 12), 12},
		{scope.MustSetByID(scope.StoreID, 13), 13},
	}
	for _, test := range tests {
		if have := test.so.ToUint32(); have != test.want {
			t.Errorf("want %d have %d", test.want, have)
		}
	}
}

func TestToUint32ByCode(t *testing.T) {
	tests := []struct {
		so   scope.Option
		want string
	}{
		{scope.Option{}, ""},
		{scope.MustSetByCode(scope.WebsiteID, "ch"), "ch"},
		{scope.MustSetByCode(scope.StoreID, "de_DE"), "de_DE"},
		{scope.MustSetByCode(scope.StoreID, "deDE"), "deDE"},
	}
	for _, test := range tests {

		var want uint32
		if test.want != "" {
			h := fnv.New32a()
			if _, err := h.Write([]byte(test.want)); err != nil {
				t.Fatal(err)
			}
			want = h.Sum32()
		}

		if have := test.so.ToUint32(); have != want {
			t.Errorf("want %d have %d", want, have)
		}
	}
}

var benchmarkToUint32 uint32
var benchmarkToUintID = scope.MockID(3141)
var benchmarkToUintCode = scope.MockCode("G0ph€r")

func BenchmarkWebsiteToUint32ByID(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		benchmarkToUint32 = scope.Option{Website: benchmarkToUintID}.ToUint32()
	}
	if benchmarkToUint32 != 3141 {
		b.Fatal("Expecting result of uint32(3141)")
	}
}

func BenchmarkGroupToUint32ByID(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		benchmarkToUint32 = scope.Option{Group: benchmarkToUintID}.ToUint32()
	}
	if benchmarkToUint32 != 3141 {
		b.Fatal("Expecting result of uint32(3141)")
	}
}
func BenchmarkStoreToUint32ByID(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		benchmarkToUint32 = scope.Option{Store: benchmarkToUintID}.ToUint32()
	}
	if benchmarkToUint32 != 3141 {
		b.Fatal("Expecting result of uint32(3141)")
	}
}

func BenchmarkWebsiteToUint32ByCode(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		benchmarkToUint32 = scope.Option{Website: benchmarkToUintCode}.ToUint32()
	}
	if benchmarkToUint32 != 1816471052 {
		b.Fatalf("Expecting result of uint32(%d)", benchmarkToUint32)
	}
}

func BenchmarkStoreToUint32ByCode(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		benchmarkToUint32 = scope.Option{Store: benchmarkToUintCode}.ToUint32()
	}
	if benchmarkToUint32 != 1816471052 {
		b.Fatal("Expecting result of uint32(1816471052)")
	}
}
