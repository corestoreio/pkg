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

package model_test

import (
	"testing"

	"github.com/corestoreio/csfw/config"
	"github.com/corestoreio/csfw/config/model"
	"github.com/corestoreio/csfw/config/scope"
	"github.com/corestoreio/csfw/config/valuelabel"
	"github.com/stretchr/testify/assert"
)

func TestStringCSV(t *testing.T) {

	wantPath := scope.StrDefault.FQPathInt64(0, "web/cors/exposed_headers")
	b := model.NewStringCSV(
		"web/cors/exposed_headers",
		valuelabel.NewByString("Content-Type", "Content Type", "X-CoreStore-ID", "CoreStore Microservice ID")...,
	)

	assert.NotEmpty(t, b.Options())

	assert.Exactly(t, []string{"Content-Type", "X-CoreStore-ID"}, b.Get(packageConfiguration, config.NewMockGetter().NewScoped(0, 0, 0)))

	assert.Exactly(t, []string{"Content-Application", "X-Gopher"}, b.Get(packageConfiguration, config.NewMockGetter(
		config.WithMockValues(config.MockPV{
			wantPath: "Content-Application,X-Gopher",
		}),
	).NewScoped(0, 0, 0)))

	mw := &config.MockWrite{}
	assert.NoError(t, b.Write(mw, []string{"a", "b", "c"}, scope.DefaultID, 0))
	assert.Exactly(t, wantPath, mw.ArgPath)
	assert.Exactly(t, "a,b,c", mw.ArgValue.(string))
}

func TestIntCSV(t *testing.T) {
	defer debugLogBuf.Reset()
	defer infoLogBuf.Reset()

	b := model.NewIntCSV(
		"web/cors/int_slice",
		valuelabel.NewByInt(valuelabel.Ints{
			{2014, "Year 2014"},
			{2015, "Year 2015"},
			{2016, "Year 2016"},
			{2017, "Year 2017"},
		})...,
	)

	assert.Len(t, b.Options(), 4)

	assert.Exactly(t, []int{2014, 2015, 2016}, b.Get(packageConfiguration, config.NewMockGetter().NewScoped(0, 0, 4)))
	assert.Exactly(t, "web/cors/int_slice", b.String())

	wantPath := scope.StrStores.FQPathInt64(4, "web/cors/int_slice")
	assert.Exactly(t, []int{}, b.Get(packageConfiguration, config.NewMockGetter(
		config.WithMockValues(config.MockPV{
			wantPath: "3015,3016",
		}),
	).NewScoped(0, 0, 4)))

	assert.Contains(t, debugLogBuf.String(), "The value '3015' cannot be found within the allowed Options")
	assert.Contains(t, debugLogBuf.String(), "The value '3016' cannot be found within the allowed Options")

	assert.Exactly(t, []int{2015, 2017}, b.Get(packageConfiguration, config.NewMockGetter(
		config.WithMockValues(config.MockPV{
			wantPath: "2015,2017",
		}),
	).NewScoped(0, 0, 4)))

	mw := &config.MockWrite{}
	assert.NoError(t, b.Write(mw, []int{2016, 2017, 2018}, scope.StoreID, 4))
	assert.Exactly(t, wantPath, mw.ArgPath)
	assert.Exactly(t, "2016,2017,2018", mw.ArgValue.(string))

	//t.Log("\n", debugLogBuf.String())

}
