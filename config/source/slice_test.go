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

package source_test

import (
	"encoding/json"
	"math"
	"testing"

	"github.com/corestoreio/csfw/config/source"
	"github.com/corestoreio/csfw/util/errors"
	"github.com/stretchr/testify/assert"
)

var _ json.Marshaler = (*source.Pair)(nil)
var _ json.Unmarshaler = (*source.Pair)(nil)

func TestSliceStringPanic(t *testing.T) {

	defer func() {
		if r := recover(); r != nil {
			err := r.(error)
			assert.True(t, errors.IsNotValid(err), "Error: %s", err)
		}
	}()
	_ = source.NewByString("kb", "l2", "ka")
}

func TestNewByStringValue(t *testing.T) {

	sl := source.NewByStringValue("a", "b", "c")
	j, err := sl.ToJSON()
	assert.NoError(t, err)
	assert.Exactly(t, "[{\"Value\":\"a\",\"Label\":\"a\"},{\"Value\":\"b\",\"Label\":\"b\"},{\"Value\":\"c\",\"Label\":\"c\"}]\n", j)
}

func TestSliceString(t *testing.T) {

	// TODO(cs) go fuzz testing
	tests := []struct {
		have      source.Slice
		wantValue string
		wantLabel string
		order     int
	}{
		{
			source.NewByString("kb", "l2", "ka", "l1", "kc", "l3", "kY", "l5", "k0", "l4"),
			`[{"Value":"k0","Label":"l4"},{"Value":"kY","Label":"l5"},{"Value":"ka","Label":"l1"},{"Value":"kb","Label":"l2"},{"Value":"kc","Label":"l3"}]` + "\n",
			`[{"Value":"ka","Label":"l1"},{"Value":"kb","Label":"l2"},{"Value":"kc","Label":"l3"},{"Value":"k0","Label":"l4"},{"Value":"kY","Label":"l5"}]` + "\n",
			0,
		},
		{
			source.NewByString("x3", "l2", "xg", "l1", "xK", "l3", "x0", "l5", "x-", "l4"),
			`[{"Value":"xg","Label":"l1"},{"Value":"xK","Label":"l3"},{"Value":"x3","Label":"l2"},{"Value":"x0","Label":"l5"},{"Value":"x-","Label":"l4"}]` + "\n",
			`[{"Value":"x0","Label":"l5"},{"Value":"x-","Label":"l4"},{"Value":"xK","Label":"l3"},{"Value":"x3","Label":"l2"},{"Value":"xg","Label":"l1"}]` + "\n",
			1,
		},
		{
			source.NewByString("x'3", "l\"2", "xög", "l1", "x\"K", "l3", `x"0`, "l5", `™¢´ƒˆ∑`, "¢£•¥ü©∑üƒ"),
			`[{"Value":"™¢´ƒˆ∑","Label":"¢£•¥ü©∑üƒ"},{"Value":"xög","Label":"l1"},{"Value":"x'3","Label":"l\"2"},{"Value":"x\"K","Label":"l3"},{"Value":"x\"0","Label":"l5"}]` + "\n",
			`[{"Value":"™¢´ƒˆ∑","Label":"¢£•¥ü©∑üƒ"},{"Value":"x\"0","Label":"l5"},{"Value":"x\"K","Label":"l3"},{"Value":"xög","Label":"l1"},{"Value":"x'3","Label":"l\"2"}]` + "\n",
			1,
		},
	}

	for i, test := range tests {
		test.have.SortByValue(test.order)
		j, err := test.have.ToJSON()
		assert.NoError(t, err)
		assert.Exactly(t, test.wantValue, j, "SortByValue Index %d", i)

		test.have.SortByLabel(test.order)
		j, err = test.have.ToJSON()
		assert.NoError(t, err)
		assert.Exactly(t, test.wantLabel, j, "SortByLabel Index %d", i)
	}
}

func TestNewByIntValue(t *testing.T) {

	sl := source.NewByIntValue(-1, 0, 5, 3, 2, -1)
	j, err := sl.ToJSON()
	assert.NoError(t, err)
	assert.Exactly(t, "[{\"Value\":-1,\"Label\":\"-1\"},{\"Value\":0,\"Label\":\"0\"},{\"Value\":5,\"Label\":\"5\"},{\"Value\":3,\"Label\":\"3\"},{\"Value\":2,\"Label\":\"2\"},{\"Value\":-1,\"Label\":\"-1\"}]\n", j)
}

func TestSliceInt(t *testing.T) {

	tests := []struct {
		have      source.Slice
		wantValue string
		wantLabel string
		order     int
	}{
		{
			source.NewByInt(source.Ints{
				{0, "http"},
				{1, "https"},
				{2, "ftp"},
				{-1, "gopher"},
			}),
			`[{"Value":-1,"Label":"gopher"},{"Value":0,"Label":"http"},{"Value":1,"Label":"https"},{"Value":2,"Label":"ftp"}]` + "\n",
			`[{"Value":2,"Label":"ftp"},{"Value":-1,"Label":"gopher"},{"Value":0,"Label":"http"},{"Value":1,"Label":"https"}]` + "\n",
			0,
		},
		{
			source.NewByInt(source.Ints{
				{0, "http"},
				{1, "https"},
				{2, "ftp"},
				{-1, "gopher"},
			}),
			`[{"Value":2,"Label":"ftp"},{"Value":1,"Label":"https"},{"Value":0,"Label":"http"},{"Value":-1,"Label":"gopher"}]` + "\n",
			`[{"Value":1,"Label":"https"},{"Value":0,"Label":"http"},{"Value":-1,"Label":"gopher"},{"Value":2,"Label":"ftp"}]` + "\n",
			1,
		},
	}

	for i, test := range tests {
		test.have.SortByInt(test.order)
		j, err := test.have.ToJSON()
		assert.NoError(t, err)
		assert.Exactly(t, test.wantValue, j, "SortByInt Index %d", i)

		test.have.SortByLabel(test.order)
		j, err = test.have.ToJSON()
		assert.NoError(t, err)
		assert.Exactly(t, test.wantLabel, j, "SortByLabel Index %d", i)
	}

	vli := source.NewByInt(source.Ints{
		{-123, "gopher"},
	})
	assert.Equal(t, "-123", vli[0].Value())
}

func TestSliceFloat64(t *testing.T) {

	tests := []struct {
		have      source.Slice
		wantValue string
		wantLabel string
		order     int
	}{
		{
			source.NewByFloat64(source.F64s{
				{math.NaN(), "nan"},
				{33.44, "http"},
				{432.432342, "https"},
				{0.0001, "ftp"},
				{-2312.3234, "gopher"},
				{math.Inf(1), "positive inf"},
				{math.Inf(-1), "negative inf"},
			}),
			`[{"Value":0,"Label":"nan"},{"Value":-1.7976931348623157e+308,"Label":"negative inf"},{"Value":-2312.3234,"Label":"gopher"},{"Value":0.0001,"Label":"ftp"},{"Value":33.44,"Label":"http"},{"Value":432.432342,"Label":"https"},{"Value":1.7976931348623157e+308,"Label":"positive inf"}]` + "\n",
			`[{"Value":0.0001,"Label":"ftp"},{"Value":-2312.3234,"Label":"gopher"},{"Value":33.44,"Label":"http"},{"Value":432.432342,"Label":"https"},{"Value":0,"Label":"nan"},{"Value":-1.7976931348623157e+308,"Label":"negative inf"},{"Value":1.7976931348623157e+308,"Label":"positive inf"}]` + "\n",
			0,
		},
		{
			source.NewByFloat64(source.F64s{
				{math.NaN(), "nan"},
				{33.44, "http"},
				{432.432342, "https"},
				{0.0001, "ftp"},
				{-2312.3234, "gopher"},
				{math.Inf(1), "positive inf"},
				{math.Inf(-1), "negative inf"},
			}),
			`[{"Value":0,"Label":"nan"},{"Value":1.7976931348623157e+308,"Label":"positive inf"},{"Value":432.432342,"Label":"https"},{"Value":33.44,"Label":"http"},{"Value":0.0001,"Label":"ftp"},{"Value":-2312.3234,"Label":"gopher"},{"Value":-1.7976931348623157e+308,"Label":"negative inf"}]` + "\n",
			`[{"Value":1.7976931348623157e+308,"Label":"positive inf"},{"Value":-1.7976931348623157e+308,"Label":"negative inf"},{"Value":0,"Label":"nan"},{"Value":432.432342,"Label":"https"},{"Value":33.44,"Label":"http"},{"Value":-2312.3234,"Label":"gopher"},{"Value":0.0001,"Label":"ftp"}]` + "\n",
			1,
		},
	}

	for i, test := range tests {
		test.have.SortByFloat64(test.order)
		j, err := test.have.ToJSON()
		assert.NoError(t, err)
		assert.Exactly(t, test.wantValue, j, "SortByFloat64 Index %d", i)

		test.have.SortByLabel(test.order)
		j, err = test.have.ToJSON()
		assert.NoError(t, err)
		assert.Exactly(t, test.wantLabel, j, "SortByLabel Index %d", i)
	}

	vli := source.NewByFloat64(source.F64s{
		{-432.432342, "https"},
		{-432.432392, "https2"},
	})
	assert.Equal(t, "-432.4323", vli[0].Value())
	assert.Equal(t, "-432.4324", vli[1].Value())

}

func TestSliceBool(t *testing.T) {

	tests := []struct {
		have      source.Slice
		wantValue string
		wantLabel string
		order     int
	}{
		{
			source.NewByBool(source.Bools{
				{true, "yes"},
				{false, "no"},
				{false, "maybe"},
				{false, "possible"},
			}),
			`[{"Value":false,"Label":"maybe"},{"Value":false,"Label":"no"},{"Value":false,"Label":"possible"},{"Value":true,"Label":"yes"}]` + "\n",
			`[{"Value":false,"Label":"maybe"},{"Value":false,"Label":"no"},{"Value":false,"Label":"possible"},{"Value":true,"Label":"yes"}]` + "\n",
			0,
		},
		{
			source.NewByBool(source.Bools{
				{false, "maybe"},
				{true, "yes"},
				{false, "no"},
				{false, "possible"},
			}),
			`[{"Value":true,"Label":"yes"},{"Value":false,"Label":"possible"},{"Value":false,"Label":"no"},{"Value":false,"Label":"maybe"}]` + "\n",
			`[{"Value":true,"Label":"yes"},{"Value":false,"Label":"possible"},{"Value":false,"Label":"no"},{"Value":false,"Label":"maybe"}]` + "\n",
			1,
		},
	}

	for i, test := range tests {
		test.have.SortByBool(test.order)
		j, err := test.have.ToJSON()
		assert.NoError(t, err)
		assert.Exactly(t, test.wantValue, j, "SortByBool Index %d", i)

		test.have.SortByLabel(test.order)
		j, err = test.have.ToJSON()
		assert.NoError(t, err)
		assert.Exactly(t, test.wantLabel, j, "SortByLabel Index %d", i)
	}

	vli := source.NewByBool(source.Bools{
		{false, "yes"},
		{true, "no"},
	})
	assert.Equal(t, "false", vli[0].Value())
	assert.Equal(t, "true", vli[1].Value())

}

func TestSliceNull(t *testing.T) {

	nullSlice := source.Slice{
		source.Pair{},
		source.Pair{},
	}

	j, err := nullSlice.ToJSON()
	assert.NoError(t, err)
	assert.Equal(t, `[{"Value":"","Label":""},{"Value":"","Label":""}]`+"\n", j)

}

func TestSliceContainsValString(t *testing.T) {

	sl := source.NewByString("k1", "v1", "k2", "v2")
	assert.True(t, sl.ContainsValString("k1"), "Search for k1 failed")
	assert.False(t, sl.ContainsValString("k0"), "Found k0 despite it is not in the slice")
}

func TestSliceContainsValInt(t *testing.T) {

	sl := source.NewByInt(source.Ints{
		{1, "v1"},
		{2, "v2"},
		{3, "v3"},
	})
	assert.True(t, sl.ContainsValInt(1), "Search for 1 failed")
	assert.False(t, sl.ContainsValInt(0), "Found 0 despite it is not in the slice")
}

func TestSliceContainsValFloat64(t *testing.T) {

	sl := source.NewByFloat64(source.F64s{
		{1.0, "v1"},
		{2.2 * 0.3, "v2"},
		{0.4 * 3, "v3"},
	})
	assert.True(t, sl.ContainsValFloat64(0.66), "Search for 0.66 failed")
	assert.False(t, sl.ContainsValFloat64(0.1), "Found 0.1 despite it is not in the slice")
}

func TestSliceContainsValBool(t *testing.T) {

	sl := source.NewByBool(source.Bools{{true, "v1"}})
	assert.True(t, sl.ContainsValBool(true), "Search for k1 failed")
	assert.False(t, sl.ContainsValBool(false), "Found k0 despite it is not in the slice")
}

func TestSliceContainsLabel(t *testing.T) {

	sl := source.NewByInt(source.Ints{
		{1, "v1"},
		{2, "v2"},
		{3, "v3"},
	})
	assert.True(t, sl.ContainsLabel("v1"), "Search for v1 failed")
	assert.False(t, sl.ContainsLabel("v0"), "Found v0 despite it is not in the slice")
}

func TestSliceEquality(t *testing.T) {

	func(sl source.Slice, vlPairs ...source.Pair) {
		vlsl := source.Slice(vlPairs)
		assert.Exactly(t, sl, vlsl)
	}(source.YesNo, source.YesNo...)
}

func TestSliceMerge(t *testing.T) {

	tests := []struct {
		in    source.Slice
		merge source.Slice
		want  string
	}{
		{
			source.NewByString("k1", "v1", "k2", "v2"),
			source.NewByString("k0", "v0", "k3", "v3", "k2", "v2a"),
			`[{"Value":"k0","Label":"v0"},{"Value":"k1","Label":"v1"},{"Value":"k2","Label":"v2a"},{"Value":"k3","Label":"v3"}]` + "\n",
		},
		{
			source.NewByInt(source.Ints{{1, "v1"}, {2, "v2"}}),
			source.NewByInt(source.Ints{{0, "v0"}, {3, "v3"}, {2, "v2a"}}),
			`[{"Value":0,"Label":"v0"},{"Value":1,"Label":"v1"},{"Value":2,"Label":"v2a"},{"Value":3,"Label":"v3"}]` + "\n",
		},
		{
			source.NewByFloat64(source.F64s{{1.1, "v1"}, {2.2, "v2"}, {0.3 * 0.2, "v32"}}),
			source.NewByFloat64(source.F64s{{0.0, "v0"}, {3.3, "v3"}, {2.2, "v2a"}, {0.3 * 0.2, "v32a"}}),
			`[{"Value":0,"Label":"v0"},{"Value":0.06,"Label":"v32a"},{"Value":1.1,"Label":"v1"},{"Value":2.2,"Label":"v2a"},{"Value":3.3,"Label":"v3"}]` + "\n",
		},
		{
			source.NewByBool(source.Bools{{false, "v1"}, {false, "v2"}}),
			source.NewByBool(source.Bools{{true, "v0"}, {true, "v3"}, {false, "v2a"}}),
			`[{"Value":false,"Label":"v2a"},{"Value":false,"Label":"v2"},{"Value":true,"Label":"v3"}]` + "\n",
		},
	}
	for _, test := range tests {

		have, err := test.in.Merge(test.merge).SortByValue(0).ToJSON()
		assert.NoError(t, err)
		if test.want != have {
			t.Errorf("\nHave: %sWant: %s\n", have, test.want)
		}
	}
}

func TestSliceUnique(t *testing.T) {

	tests := []struct {
		in   source.Slice
		want string
	}{
		{
			source.NewByString("k2", "v20", "k1", "v1", "k2", "v21"),
			`[{"Value":"k2","Label":"v20"},{"Value":"k1","Label":"v1"}]` + "\n",
		},
		{
			source.NewByInt(source.Ints{{1, "v1"}, {2, "v20"}, {2, "v21"}}),
			`[{"Value":1,"Label":"v1"},{"Value":2,"Label":"v20"}]` + "\n",
		},
		{
			source.NewByFloat64(source.F64s{{0.3 * 0.2, "v31"}, {1.1, "v1"}, {2.2, "v2"}, {0.3 * 0.2, "v32"}}),
			`[{"Value":0.06,"Label":"v31"},{"Value":1.1,"Label":"v1"},{"Value":2.2,"Label":"v2"}]` + "\n",
		},
		{
			source.NewByBool(source.Bools{{false, "v1"}, {false, "v2"}, {true, "v3"}}),
			`[{"Value":false,"Label":"v1"},{"Value":true,"Label":"v3"}]` + "\n",
		},
	}
	for _, test := range tests {

		have, err := test.in.Unique().ToJSON()
		assert.NoError(t, err)
		if test.want != have {
			t.Errorf("\nHave: %sWant: %s\n", have, test.want)
		}
	}
}

func TestSliceUnmarshalJSON(t *testing.T) {

	tests := []struct {
		in      []byte
		want    source.Slice
		wantErr error
	}{
		{
			[]byte(`[{"Value":"k2","Label":"v20"},{"Value":"k1","Label":"v1"}]`),
			source.NewByString("k2", "v20", "k1", "v1"),
			nil,
		},
		{
			[]byte(`[{"Value":1,"Label":"v20"},{"Value":2,"Label":"v1"}]`),
			source.NewByInt(source.Ints{{1, "v20"}, {2, "v1"}}),
			nil,
		},
		{
			[]byte(`[{"Value":false,"Label":"v1"},{"Value":true,"Label":"v3"}]`),
			source.NewByBool(source.Bools{{false, "v1"}, {true, "v3"}}),
			nil,
		},
		{
			[]byte(`[{"Value":3.1415678,"Label":"pi"},{"Value":2.718281,"Label":"e"}]`),
			source.NewByFloat64(source.F64s{{3.1415678, "pi"}, {2.718281, "e"}}),
			nil,
		},
		{
			[]byte(`[{"Value":3.1415678,"Label":"pi"},{"Value":2374652873645287346523465,"Label":"overflow"}]`),
			source.NewByFloat64(source.F64s{{3.1415678, "pi"}, {2.3746528736452872e+24, "overflow"}}),
			nil,
		},
		{
			[]byte(`[{"Value":3.1415678,"Label":true} ]`),
			source.Slice{source.Pair{}},
			errors.New("[source] Unmarshal: \"{\\\"Value\\\":3.1415678,\\\"Label\\\":true}\": json: cannot unmarshal bool into Go value of type string"),
		},
	}
	for i, test := range tests {
		var have source.Slice
		err := json.Unmarshal(test.in, &have)
		if test.wantErr != nil {
			assert.EqualError(t, err, test.wantErr.Error(), "Index %d", i)
			continue
		}
		assert.NoError(t, err, "Index %d", i)
		assert.Exactly(t, test.want, have, "Index %d", i)
	}
}
