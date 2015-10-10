// Copyright 2015, Cyrill @ Schumacher.fm and the CoreStore contributors
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

package valuelabel_test

import (
	"math"
	"testing"

	"github.com/corestoreio/csfw/config/valuelabel"
	"github.com/corestoreio/csfw/utils"
	"github.com/stretchr/testify/assert"
)

func TestValueLabelSliceStringPanic(t *testing.T) {
	defer func() {
		if r := recover(); r != nil {
			assert.EqualError(t, r.(error), valuelabel.ErrImbalancedPairs.Error())
		}
	}()
	_ = valuelabel.NewByString("kb", "l2", "ka")
}

func TestValueLabelSliceString(t *testing.T) {
	// TODO(cs) go fuzz testing
	tests := []struct {
		have      valuelabel.Slice
		wantValue string
		wantLabel string
		order     utils.SortDirection
	}{
		{
			valuelabel.NewByString("kb", "l2", "ka", "l1", "kc", "l3", "kY", "l5", "k0", "l4"),
			`[{"Value":"k0","Label":"l4"},{"Value":"kY","Label":"l5"},{"Value":"ka","Label":"l1"},{"Value":"kb","Label":"l2"},{"Value":"kc","Label":"l3"}]` + "\n",
			`[{"Value":"ka","Label":"l1"},{"Value":"kb","Label":"l2"},{"Value":"kc","Label":"l3"},{"Value":"k0","Label":"l4"},{"Value":"kY","Label":"l5"}]` + "\n",
			utils.SortAsc,
		},
		{
			valuelabel.NewByString("x3", "l2", "xg", "l1", "xK", "l3", "x0", "l5", "x-", "l4"),
			`[{"Value":"xg","Label":"l1"},{"Value":"xK","Label":"l3"},{"Value":"x3","Label":"l2"},{"Value":"x0","Label":"l5"},{"Value":"x-","Label":"l4"}]` + "\n",
			`[{"Value":"x0","Label":"l5"},{"Value":"x-","Label":"l4"},{"Value":"xK","Label":"l3"},{"Value":"x3","Label":"l2"},{"Value":"xg","Label":"l1"}]` + "\n",
			utils.SortDesc,
		},
		{
			valuelabel.NewByString("x'3", "l\"2", "xög", "l1", "x\"K", "l3", `x"0`, "l5", `™¢´ƒˆ∑`, "¢£•¥ü©∑üƒ"),
			`[{"Value":"™¢´ƒˆ∑","Label":"¢£•¥ü©∑üƒ"},{"Value":"xög","Label":"l1"},{"Value":"x'3","Label":"l\"2"},{"Value":"x\"K","Label":"l3"},{"Value":"x\"0","Label":"l5"}]` + "\n",
			`[{"Value":"™¢´ƒˆ∑","Label":"¢£•¥ü©∑üƒ"},{"Value":"x\"0","Label":"l5"},{"Value":"x\"K","Label":"l3"},{"Value":"xög","Label":"l1"},{"Value":"x'3","Label":"l\"2"}]` + "\n",
			utils.SortDesc,
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

func TestValueLabelSliceInt(t *testing.T) {
	tests := []struct {
		have      valuelabel.Slice
		wantValue string
		wantLabel string
		order     utils.SortDirection
	}{
		{
			valuelabel.NewByInt(valuelabel.Ints{
				{0, "http"},
				{1, "https"},
				{2, "ftp"},
				{-1, "gopher"},
			}),
			`[{"Value":-1,"Label":"gopher"},{"Value":0,"Label":"http"},{"Value":1,"Label":"https"},{"Value":2,"Label":"ftp"}]` + "\n",
			`[{"Value":2,"Label":"ftp"},{"Value":-1,"Label":"gopher"},{"Value":0,"Label":"http"},{"Value":1,"Label":"https"}]` + "\n",
			utils.SortAsc,
		},
		{
			valuelabel.NewByInt(valuelabel.Ints{
				{0, "http"},
				{1, "https"},
				{2, "ftp"},
				{-1, "gopher"},
			}),
			`[{"Value":2,"Label":"ftp"},{"Value":1,"Label":"https"},{"Value":0,"Label":"http"},{"Value":-1,"Label":"gopher"}]` + "\n",
			`[{"Value":1,"Label":"https"},{"Value":0,"Label":"http"},{"Value":-1,"Label":"gopher"},{"Value":2,"Label":"ftp"}]` + "\n",
			utils.SortDesc,
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

	vli := valuelabel.NewByInt(valuelabel.Ints{
		{-123, "gopher"},
	})
	assert.Equal(t, "-123", vli[0].Value())
}

func TestValueLabelSliceFloat64(t *testing.T) {
	tests := []struct {
		have      valuelabel.Slice
		wantValue string
		wantLabel string
		order     utils.SortDirection
	}{
		{
			valuelabel.NewByFloat64(valuelabel.F64s{
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
			utils.SortAsc,
		},
		{
			valuelabel.NewByFloat64(valuelabel.F64s{
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
			utils.SortDesc,
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

	vli := valuelabel.NewByFloat64(valuelabel.F64s{
		{-432.432342, "https"},
		{-432.432392, "https2"},
	})
	assert.Equal(t, "-432.4323", vli[0].Value())
	assert.Equal(t, "-432.4324", vli[1].Value())

}

func TestValueLabelSliceBool(t *testing.T) {
	tests := []struct {
		have      valuelabel.Slice
		wantValue string
		wantLabel string
		order     utils.SortDirection
	}{
		{
			valuelabel.NewByBool(valuelabel.Bools{
				{true, "yes"},
				{false, "no"},
				{false, "maybe"},
				{false, "possible"},
			}),
			`[{"Value":false,"Label":"maybe"},{"Value":false,"Label":"no"},{"Value":false,"Label":"possible"},{"Value":true,"Label":"yes"}]` + "\n",
			`[{"Value":false,"Label":"maybe"},{"Value":false,"Label":"no"},{"Value":false,"Label":"possible"},{"Value":true,"Label":"yes"}]` + "\n",
			utils.SortAsc,
		},
		{
			valuelabel.NewByBool(valuelabel.Bools{
				{false, "maybe"},
				{true, "yes"},
				{false, "no"},
				{false, "possible"},
			}),
			`[{"Value":true,"Label":"yes"},{"Value":false,"Label":"possible"},{"Value":false,"Label":"no"},{"Value":false,"Label":"maybe"}]` + "\n",
			`[{"Value":true,"Label":"yes"},{"Value":false,"Label":"possible"},{"Value":false,"Label":"no"},{"Value":false,"Label":"maybe"}]` + "\n",
			utils.SortDesc,
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

	vli := valuelabel.NewByBool(valuelabel.Bools{
		{false, "yes"},
		{true, "no"},
	})
	assert.Equal(t, "false", vli[0].Value())
	assert.Equal(t, "true", vli[1].Value())

}

func TestValueLabelSliceNull(t *testing.T) {
	nullSlice := valuelabel.Slice{
		valuelabel.Pair{},
		valuelabel.Pair{},
	}

	j, err := nullSlice.ToJSON()
	assert.NoError(t, err)
	assert.Equal(t, `[{"Value":"","Label":""},{"Value":"","Label":""}]`+"\n", j)

}
