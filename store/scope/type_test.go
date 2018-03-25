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

package scope

import (
	"encoding/json"
	"testing"

	"github.com/corestoreio/errors"
	"github.com/stretchr/testify/assert"
)

var _ json.Marshaler = (*Type)(nil)
var _ json.Unmarshaler = (*Type)(nil)

func TestTypeBits(t *testing.T) {

	const (
		scope1 Type = iota + 1
		scope2
		scope3
		scope4
		scope5
	)

	tests := []struct {
		have    []Type
		want    Type
		notWant Type
		human   []string
		string
	}{
		{[]Type{scope1, scope2}, scope2, scope3, []string{"Default", "Website"}, "Default,Website"},
		{[]Type{scope3, scope4}, scope3, scope2, []string{"Group", "Store"}, "Group,Store"},
		{[]Type{scope4, scope5}, scope4, scope2, []string{"Store"}, "Store"},
	}

	for _, test := range tests {
		var b = Perm(0).Set(test.have...)
		if !b.Has(test.want) {
			t.Errorf("%d should contain %d", b, test.want)
		}
		if b.Has(test.notWant) {
			t.Errorf("%d should not contain %d", b, test.notWant)
		}
		assert.EqualValues(t, test.human, b.Human())
		assert.EqualValues(t, test.string, b.String())
	}
}

func TestFromString(t *testing.T) {

	tests := []struct {
		have string
		want Type
	}{
		{"asdasd", Default},
		{strDefault, Default},
		{strWebsites, Website},
		{strStores, Store},
	}
	for _, test := range tests {
		assert.Equal(t, test.want, FromString(test.have))
	}
}

func TestFromType(t *testing.T) {

	tests := []struct {
		have Type
		want TypeStr
	}{
		{Default, StrDefault},
		{Absent, StrDefault},
		{Group, StrDefault},
		{Website, StrWebsites},
		{Store, StrStores},
	}
	for _, test := range tests {
		assert.Exactly(t, test.want, FromType(test.have))
		assert.Exactly(t, test.want.String(), test.have.StrType())
	}
}

func TestStrType(t *testing.T) {

	assert.Equal(t, strDefault, StrDefault.String())
	assert.Equal(t, strWebsites, StrWebsites.String())
	assert.Equal(t, strStores, StrStores.String())

	assert.Exactly(t, Default, StrDefault.Type())
	assert.Exactly(t, Website, StrWebsites.Type())
	assert.Exactly(t, Store, StrStores.Type())
}

func TestValid(t *testing.T) {

	tests := []struct {
		have string
		want bool
	}{
		{"Rust", false},
		{"default", true},
		{"website", false},
		{"websites", true},
		{"stores", true},
		{"Stores", false},
	}
	for i, test := range tests {
		assert.Exactly(t, test.want, Valid(test.have), "Index %d", i)
	}
}

func TestFromBytes(t *testing.T) {
	tests := []struct {
		have []byte
		want Type
	}{
		{[]byte("asdasd"), Default},
		{[]byte(strDefault), Default},
		{[]byte(strWebsites), Website},
		{[]byte(strStores), Store},
	}
	for _, test := range tests {
		assert.Exactly(t, test.want, FromBytes(test.have))
	}
}

func TestValidBytes(t *testing.T) {
	tests := []struct {
		have []byte
		want bool
	}{
		{[]byte("Rust"), false},
		{[]byte("default"), true},
		{[]byte("website"), false},
		{[]byte("websites"), true},
		{[]byte("stores"), true},
		{[]byte("Stores"), false},
	}
	for i, test := range tests {
		assert.Exactly(t, test.want, ValidBytes(test.have), "Index %d", i)
	}
}

func TestStrTypeBytes(t *testing.T) {
	tests := []struct {
		id Type
	}{
		{Default},
		{Website},
		{Store},
		{44},
	}
	for i, test := range tests {
		assert.Exactly(t, test.id.StrType(), string(test.id.StrBytes()), "Index %d", i)
	}
}

func TestValidParent(t *testing.T) {
	tests := []struct {
		c    Type
		p    Type
		want bool
	}{
		{Default, Default, true},
		{Website, Default, true},
		{Store, Website, true},
		{Default, Website, false},
		{Absent, Absent, false},
		{Absent, Default, false},
		{Default, Absent, false},
	}
	for i, test := range tests {
		if have, want := ValidParent(test.c, test.p), test.want; have != want {
			t.Errorf("(%d) Have: %v Want: %v", i, have, want)
		}
	}
}

func TestType_MarshalJSON(t *testing.T) {
	tests := []struct {
		s    Type
		want []byte
	}{
		{Default, jsonDefault},
		{Website, jsonWebsite},
		{Group, jsonGroup},
		{Store, jsonStore},
		{Absent, jsonDefault},
		{128, jsonDefault},
	}
	for i, test := range tests {
		have, err := test.s.MarshalJSON()
		if err != nil {
			t.Fatal(i, err)
		}
		assert.Exactly(t, test.want, have, "Index %d", i)
	}
}

func TestType_UnmarshalJSON(t *testing.T) {
	tests := []struct {
		raw  []byte
		want Type
	}{
		{jsonDefault, Default},
		{jsonWebsite, Website},
		{jsonGroup, Group},
		{jsonStore, Store},
		{[]byte("Evi'l\x00"), Default},
	}
	for i, test := range tests {
		var have Type
		if err := have.UnmarshalJSON(test.raw); err != nil {
			t.Fatal(i, err)
		}
		assert.Exactly(t, test.want, have, "Index %d", i)
	}
}

func TestType_JSON(t *testing.T) {

	type x struct {
		Str string `json:"str"`
		Scp Type   `json:"myType"`
		ID  int64  `json:"id"`
	}

	var xt = x{
		Str: "Gopher",
		Scp: Website,
		ID:  3,
	}
	raw, err := json.Marshal(xt)
	if err != nil {
		t.Fatal(err)
	}

	var xt2 x
	if err := json.Unmarshal(raw, &xt2); err != nil {
		t.Fatal(err)
	}
	assert.Exactly(t, xt, xt2)
}

func TestType_Pack(t *testing.T) {
	tests := []struct {
		s    Type
		id   int64
		want TypeID
	}{
		{Website, 1, MakeTypeID(Website, 1)},
		{Store, 4, MakeTypeID(Store, 4)},
		{0, 0, 0},
	}
	for i, test := range tests {
		assert.Exactly(t, test.want, test.s.Pack(test.id), "Index %d", i)
	}
}

func TestType_IsValid(t *testing.T) {
	t1 := Store
	assert.NoError(t, t1.IsValid())
	t2 := Type(234)
	assert.True(t, errors.NotValid.Match(t2.IsValid()), "Should have error kind NotValid")
}
