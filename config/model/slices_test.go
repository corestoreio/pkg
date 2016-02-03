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

package model_test

import (
	"testing"

	"github.com/corestoreio/csfw/config"
	"github.com/corestoreio/csfw/config/model"
	"github.com/corestoreio/csfw/config/path"
	"github.com/corestoreio/csfw/config/source"
	"github.com/corestoreio/csfw/store/scope"
	"github.com/stretchr/testify/assert"
)

func TestStringCSV(t *testing.T) {
	t.Parallel()
	const pathWebCorsHeaders = "web/cors/exposed_headers"
	wantPath := path.MustNewByParts(pathWebCorsHeaders)
	b := model.NewStringCSV(
		"web/cors/exposed_headers",
		model.WithConfigStructure(configStructure),
		model.WithSourceByString(
			"Content-Type", "Content Type", "X-CoreStore-ID", "CoreStore Microservice ID",
		),
	)

	assert.NotEmpty(t, b.Options())

	assert.Exactly(t, []string{"Content-Type", "X-CoreStore-ID"}, b.Get(config.NewMockGetter().NewScoped(0, 0, 0)))

	assert.Exactly(t, []string{"Content-Application", "X-Gopher"}, b.Get(config.NewMockGetter(
		config.WithMockValues(config.MockPV{
			wantPath.String(): "Content-Application,X-Gopher",
		}),
	).NewScoped(0, 0, 0)))

	mw := &config.MockWrite{}
	b.Source.Merge(source.NewByString("a", "a", "b", "b", "c", "c"))

	assert.NoError(t, b.Write(mw, []string{"a", "b", "c"}, scope.DefaultID, 0))
	assert.Exactly(t, wantPath.String(), mw.ArgPath)
	assert.Exactly(t, "a,b,c", mw.ArgValue.(string))
}

func TestIntCSV(t *testing.T) {
	defer debugLogBuf.Reset()
	defer infoLogBuf.Reset()

	const pathWebCorsIntSlice = "web/cors/int_slice"

	b := model.NewIntCSV(
		pathWebCorsIntSlice,
		model.WithConfigStructure(configStructure),
		model.WithSourceByInt(source.Ints{
			{2014, "Year 2014"},
			{2015, "Year 2015"},
			{2016, "Year 2016"},
			{2017, "Year 2017"},
		}),
	)

	assert.Len(t, b.Options(), 4)

	assert.Exactly(t, []int{2014, 2015, 2016}, b.Get(config.NewMockGetter().NewScoped(0, 0, 4)))
	assert.Exactly(t, pathWebCorsIntSlice, b.String())

	wantPath := path.MustNewByParts(pathWebCorsIntSlice).Bind(scope.StoreID, 4)
	assert.Exactly(t, []int{}, b.Get(config.NewMockGetter(
		config.WithMockValues(config.MockPV{
			wantPath.String(): "3015,3016",
		}),
	).NewScoped(0, 0, 4)))

	assert.Contains(t, debugLogBuf.String(), "The value '3015' cannot be found within the allowed Options")
	assert.Contains(t, debugLogBuf.String(), "The value '3016' cannot be found within the allowed Options")

	assert.Exactly(t, []int{2015, 2017}, b.Get(config.NewMockGetter(
		config.WithMockValues(config.MockPV{
			wantPath.String(): "2015,2017",
		}),
	).NewScoped(0, 0, 4)))

	mw := &config.MockWrite{}
	b.Source.Merge(source.NewByInt(source.Ints{
		{2018, "Year 2018"},
	}))
	assert.NoError(t, b.Write(mw, []int{2016, 2017, 2018}, scope.StoreID, 4))
	assert.Exactly(t, wantPath.String(), mw.ArgPath)
	assert.Exactly(t, "2016,2017,2018", mw.ArgValue.(string))

	//t.Log("\n", debugLogBuf.String())

}
