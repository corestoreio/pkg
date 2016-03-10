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

	"errors"

	"github.com/corestoreio/csfw/config/cfgmock"
	"github.com/corestoreio/csfw/config/cfgpath"
	"github.com/corestoreio/csfw/config/model"
	"github.com/corestoreio/csfw/config/source"
	"github.com/corestoreio/csfw/store/scope"
	"github.com/stretchr/testify/assert"
)

func TestStringCSVGet(t *testing.T) {
	t.Parallel()
	const pathWebCorsHeaders = "web/cors/exposed_headers"
	wantPath := cfgpath.MustNewByParts(pathWebCorsHeaders).String()
	b := model.NewStringCSV(
		"web/cors/exposed_headers",
		model.WithFieldFromSectionSlice(configStructure),
		model.WithSourceByString(
			"Content-Type", "Content Type", "X-CoreStore-ID", "CoreStore Microservice ID",
		),
	)
	assert.NotEmpty(t, b.Options())

	sl, err := b.Get(cfgmock.NewService().NewScoped(0, 0))
	assert.NoError(t, err)
	assert.Exactly(t, []string{"Content-Type", "X-CoreStore-ID"}, sl) // default values from variable configStructure

	tests := []struct {
		have    string
		want    []string
		wantErr error
	}{
		{"Content-Type,X-CoreStore-ID", []string{"Content-Type", "X-CoreStore-ID"}, nil},
		{"", nil, nil},
		{"X-CoreStore-ID", []string{"X-CoreStore-ID"}, nil},
		{"Content-Type,X-CS", []string{"Content-Type", "X-CS"}, nil},
		// todo add errors
	}
	for i, test := range tests {
		haveSL, haveErr := b.Get(cfgmock.NewService(
			cfgmock.WithPV(cfgmock.PathValue{
				wantPath: test.have,
			}),
		).NewScoped(1, 0)) // 1,0 because scope of pathWebCorsHeaders is default,website

		assert.Exactly(t, test.want, haveSL, "Index %d", i)
		if test.wantErr != nil {
			assert.EqualError(t, haveErr, test.wantErr.Error(), "Index %d", i)
			continue
		}
		assert.NoError(t, haveErr, "Index %d", i)
	}
}

func TestStringCSVWrite(t *testing.T) {
	t.Parallel()
	const pathWebCorsHeaders = "web/cors/exposed_headers"
	wantPath := cfgpath.MustNewByParts(pathWebCorsHeaders).String()
	b := model.NewStringCSV(
		"web/cors/exposed_headers",
		model.WithFieldFromSectionSlice(configStructure),
		model.WithSourceByString(
			"Content-Type", "Content Type", "X-CoreStore-ID", "CoreStore Microservice ID",
		),
	)

	mw := &cfgmock.Write{}
	b.Source.Merge(source.NewByString("a", "a", "b", "b", "c", "c"))

	assert.NoError(t, b.Write(mw, []string{"a", "b", "c"}, scope.DefaultID, 0))
	assert.Exactly(t, wantPath, mw.ArgPath)
	assert.Exactly(t, "a,b,c", mw.ArgValue.(string))
	assert.EqualError(t,
		b.Write(mw, []string{"abc"}, scope.DefaultID, 0),
		"The value 'abc' cannot be found within the allowed Options():\n[{\"Value\":\"Content-Type\",\"Label\":\"Content Type\"},{\"Value\":\"X-CoreStore-ID\",\"Label\":\"CoreStore Microservice ID\"},{\"Value\":\"a\",\"Label\":\"a\"},{\"Value\":\"b\",\"Label\":\"b\"},{\"Value\":\"c\",\"Label\":\"c\"}]\n",
	)
}

func TestStringCSVCustomSeparator(t *testing.T) {
	t.Parallel()

	const cfgPath = "aa/bb/cc"

	b := model.NewStringCSV(
		cfgPath,
		model.WithSourceByString(
			"2014", "Year 2014",
			"2015", "Year 2015",
			"2016", "Year 2016",
			"2017", "Year 2017",
		),
		model.WithCSVSeparator(''),
	)
	wantPath := cfgpath.MustNewByParts(cfgPath).String() // Default Scope

	haveSL, haveErr := b.Get(cfgmock.NewService(
		cfgmock.WithPV(cfgmock.PathValue{
			wantPath: `20152016`,
		}),
	).NewScoped(34, 4))
	if haveErr != nil {
		t.Fatal(haveErr)
	}

	assert.Exactly(t, []string{"2015", "2016"}, haveSL)
}

func TestIntCSV(t *testing.T) {
	t.Parallel()

	const pathWebCorsIntSlice = "web/cors/int_slice"

	b := model.NewIntCSV(
		pathWebCorsIntSlice,
		model.WithFieldFromSectionSlice(configStructure),
		model.WithSourceByInt(source.Ints{
			{2014, "Year 2014"},
			{2015, "Year 2015"},
			{2016, "Year 2016"},
			{2017, "Year 2017"},
		}),
	)
	assert.Len(t, b.Options(), 4)
	assert.Exactly(t, pathWebCorsIntSlice, b.String())
	// default values:
	sl, err := b.Get(cfgmock.NewService().NewScoped(0, 4))
	assert.NoError(t, err)
	assert.Exactly(t, []int{2014, 2015, 2016}, sl) // three years are defined in variable configStructure

	wantPath := cfgpath.MustNewByParts(pathWebCorsIntSlice).Bind(scope.StoreID, 4).String()

	tests := []struct {
		lenient bool
		have    string
		want    []int
		wantErr error
	}{
		{false, "3015,3016", []int{3015, 3016}, nil},
		{false, "2015,2017", []int{2015, 2017}, nil},
		{false, "", nil, nil},
		{false, "2015,,20x17", []int{2015}, errors.New("strconv.ParseInt: parsing \"20x17\": invalid syntax")},
		{true, "2015,,2017", []int{2015, 2017}, nil},
	}
	for i, test := range tests {
		b.Lenient = test.lenient
		haveSL, haveErr := b.Get(cfgmock.NewService(
			cfgmock.WithPV(cfgmock.PathValue{
				wantPath: test.have,
			}),
		).NewScoped(0, 4))

		assert.Exactly(t, test.want, haveSL, "Index %d", i)
		if test.wantErr != nil {
			assert.EqualError(t, haveErr, test.wantErr.Error(), "Index %d", i)
			continue
		}
		assert.NoError(t, haveErr, "Index %d", i)
	}
}

func TestIntCSVWrite(t *testing.T) {
	t.Parallel()

	const pathWebCorsIntSlice = "web/cors/int_slice"

	b := model.NewIntCSV(
		pathWebCorsIntSlice,
		model.WithFieldFromSectionSlice(configStructure),
		model.WithSourceByInt(source.Ints{
			{2014, "Year 2014"},
			{2015, "Year 2015"},
			{2016, "Year 2016"},
			{2017, "Year 2017"},
		}),
	)
	wantPath := cfgpath.MustNewByParts(pathWebCorsIntSlice).Bind(scope.StoreID, 4).String()

	mw := &cfgmock.Write{}
	b.Source.Merge(source.NewByInt(source.Ints{
		{2018, "Year 2018"},
	}))
	assert.NoError(t, b.Write(mw, []int{2016, 2017, 2018}, scope.StoreID, 4))
	assert.Exactly(t, wantPath, mw.ArgPath)
	assert.Exactly(t, "2016,2017,2018", mw.ArgValue.(string))
	assert.EqualError(t,
		b.Write(mw, []int{2019}, scope.StoreID, 4),
		"The value '2019' cannot be found within the allowed Options():\n[{\"Value\":2014,\"Label\":\"Year 2014\"},{\"Value\":2015,\"Label\":\"Year 2015\"},{\"Value\":2016,\"Label\":\"Year 2016\"},{\"Value\":2017,\"Label\":\"Year 2017\"},{\"Value\":2018,\"Label\":\"Year 2018\"}]\n",
	)
}

func TestIntCSVCustomSeparator(t *testing.T) {
	t.Parallel()

	const pathWebCorsIntSlice = "web/cors/int_slice"

	b := model.NewIntCSV(
		pathWebCorsIntSlice,
		model.WithFieldFromSectionSlice(configStructure),
		model.WithSourceByInt(source.Ints{
			{2014, "Year 2014"},
			{2015, "Year 2015"},
			{2016, "Year 2016"},
			{2017, "Year 2017"},
		}),
		model.WithCSVSeparator('|'),
	)
	wantPath := cfgpath.MustNewByParts(pathWebCorsIntSlice).Bind(scope.WebsiteID, 34).String()

	haveSL, haveErr := b.Get(cfgmock.NewService(
		cfgmock.WithPV(cfgmock.PathValue{
			wantPath: `2015|2016|`,
		}),
	).NewScoped(34, 4))
	if haveErr != nil {
		t.Fatal(haveErr)
	}

	assert.Exactly(t, []int{2015, 2016}, haveSL)
}
