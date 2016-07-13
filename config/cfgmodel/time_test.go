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

package cfgmodel_test

import (
	"testing"
	"time"

	"github.com/corestoreio/csfw/config"
	"github.com/corestoreio/csfw/config/cfgmock"
	"github.com/corestoreio/csfw/config/cfgmodel"
	"github.com/corestoreio/csfw/config/cfgpath"
	"github.com/corestoreio/csfw/config/element"
	"github.com/corestoreio/csfw/store/scope"
	"github.com/corestoreio/csfw/util/conv"
	"github.com/corestoreio/csfw/util/errors"
	"github.com/stretchr/testify/assert"
)

func mustParseTime(s string) time.Time {
	t, err := conv.StringToDate(s, nil)
	if err != nil {
		panic(err)
	}
	return t
}

func TestTimeGetWithCfgStruct(t *testing.T) {

	const pathWebCorsTime = "web/cors/time"
	tm := cfgmodel.NewTime("web/cors/time", cfgmodel.WithFieldFromSectionSlice(configStructure))
	assert.Empty(t, tm.Options())

	wantPath := cfgpath.MustNewByParts(pathWebCorsTime).Bind(scope.Website, 10)
	defaultTime := mustParseTime("2012-08-23 09:20:13")
	tests := []struct {
		sg       config.Scoped
		wantHash scope.Hash
		want     time.Time
	}{
		{cfgmock.NewService().NewScoped(0, 0), scope.DefaultHash, defaultTime}, // because default value in packageConfiguration
		{cfgmock.NewService().NewScoped(0, 1), scope.DefaultHash, defaultTime}, // because default value in packageConfiguration
		{cfgmock.NewService().NewScoped(1, 1), scope.DefaultHash, defaultTime}, // because default value in packageConfiguration
		{cfgmock.NewService(cfgmock.WithPV(cfgmock.PathValue{wantPath.Bind(scope.Website, 10).String(): defaultTime.Add(time.Second * 2)})).NewScoped(10, 0), scope.NewHash(scope.Website, 10), defaultTime.Add(time.Second * 2)},
		{cfgmock.NewService(cfgmock.WithPV(cfgmock.PathValue{wantPath.Bind(scope.Website, 10).String(): defaultTime.Add(time.Second * 3)})).NewScoped(10, 1), scope.NewHash(scope.Website, 10), defaultTime.Add(time.Second * 3)},
		{cfgmock.NewService(cfgmock.WithPV(cfgmock.PathValue{
			wantPath.String():                       defaultTime.Add(time.Second * 5),
			wantPath.Bind(scope.Store, 11).String(): defaultTime.Add(time.Second * 6),
		})).NewScoped(10, 11), scope.NewHash(scope.Store, 11), defaultTime.Add(time.Second * 6)},
	}
	for i, test := range tests {
		gb, h, err := tm.Get(test.sg)
		if err != nil {
			t.Fatal("Index", i, err)
		}
		assert.Exactly(t, test.want, gb, "Index %d", i)
		assert.Exactly(t, test.wantHash.String(), h.String(), "Index %d", i)
	}
}

func TestTimeGetWithoutCfgStruct(t *testing.T) {

	const pathWebCorsTime = "web/cors/time"
	b := cfgmodel.NewTime(pathWebCorsTime)
	assert.Empty(t, b.Options())

	wantPath := cfgpath.MustNewByParts(pathWebCorsTime).Bind(scope.Website, 10)
	defaultTime := mustParseTime("2012-08-23 09:20:13")
	tests := []struct {
		sg       config.Scoped
		wantHash scope.Hash
		want     time.Time
	}{
		{cfgmock.NewService().NewScoped(1, 1), scope.DefaultHash, time.Time{}}, // because default value in packageConfiguration
		{cfgmock.NewService(cfgmock.WithPV(cfgmock.PathValue{wantPath.String(): defaultTime.Add(time.Second * 2)})).NewScoped(10, 0), scope.DefaultHash, time.Time{}},
		{cfgmock.NewService(cfgmock.WithPV(cfgmock.PathValue{wantPath.String(): defaultTime.Add(time.Second * 3)})).NewScoped(10, 1), scope.DefaultHash, time.Time{}},
		{cfgmock.NewService(cfgmock.WithPV(cfgmock.PathValue{wantPath.Bind(scope.Default, 0).String(): defaultTime.Add(time.Second * 3)})).NewScoped(0, 0), scope.DefaultHash, defaultTime.Add(time.Second * 3)},
		{cfgmock.NewService(cfgmock.WithPV(cfgmock.PathValue{
			wantPath.Bind(scope.Default, 0).String(): defaultTime.Add(time.Second * 5),
			wantPath.Bind(scope.Store, 11).String():  defaultTime.Add(time.Second * 6),
		})).NewScoped(10, 11), scope.DefaultHash, defaultTime.Add(time.Second * 5)},
	}
	for i, test := range tests {
		gb, h, err := b.Get(test.sg)
		if err != nil {
			t.Fatal("Index", i, err)
		}
		assert.Exactly(t, test.want, gb, "Index %d", i)
		assert.Exactly(t, test.wantHash.String(), h.String(), "Index %d", i)
	}
}

func TestTimeGetWithoutCfgStructShouldReturnUnexpectedError(t *testing.T) {

	b := cfgmodel.NewTime("web/cors/time")
	assert.Empty(t, b.Options())

	gb, h, haveErr := b.Get(cfgmock.NewService(
		cfgmock.WithTime(func(path string) (time.Time, error) {
			return time.Time{}, errors.NewFatalf("Unexpected error")
		}),
	).NewScoped(1, 1))
	assert.Empty(t, gb)
	assert.True(t, errors.IsFatal(haveErr), "Error: %s", haveErr)
	assert.Exactly(t, scope.DefaultHash.String(), h.String())
}

func TestTimeIgnoreNilDefaultValues(t *testing.T) {

	b := cfgmodel.NewTime("web/cors/time", cfgmodel.WithField(&element.Field{}))
	gb, h, err := b.Get(cfgmock.NewService().NewScoped(1, 1))
	if err != nil {
		t.Fatal(err)
	}
	assert.Exactly(t, time.Time{}, gb)
	assert.Exactly(t, scope.DefaultHash.String(), h.String())
}

func TestTimeWrite(t *testing.T) {

	const pathWebCorsF64 = "web/cors/time"
	wantPath := cfgpath.MustNewByParts(pathWebCorsF64).Bind(scope.Website, 10)
	haveTime := mustParseTime("2000-08-23 09:20:13")

	b := cfgmodel.NewTime("web/cors/time", cfgmodel.WithFieldFromSectionSlice(configStructure))

	mw := &cfgmock.Write{}
	assert.NoError(t, b.Write(mw, haveTime, scope.Website, 10))
	assert.Exactly(t, wantPath.String(), mw.ArgPath)
	assert.Exactly(t, haveTime, mw.ArgValue.(time.Time))
}

//Scopes:    scope.PermStore,
//Default:   "1h45m",

func mustParseDuration(s string) time.Duration {
	t, err := conv.ToDurationE(s)
	if err != nil {
		panic(err)
	}
	return t
}

func TestDurationGetWithCfgStruct(t *testing.T) {

	const pathWebCorsDuration = "web/cors/duration"
	tm := cfgmodel.NewDuration("web/cors/duration", cfgmodel.WithFieldFromSectionSlice(configStructure))
	assert.Empty(t, tm.Options())

	wantPath := cfgpath.MustNewByParts(pathWebCorsDuration).Bind(scope.Website, 10)
	defaultDuration := mustParseDuration("1h45m") // default as in the configStructure slice

	tests := []struct {
		sg       config.Scoped
		wantHash scope.Hash
		want     time.Duration
	}{
		{cfgmock.NewService().NewScoped(0, 0), scope.DefaultHash, defaultDuration}, // because default value in packageConfiguration
		{cfgmock.NewService().NewScoped(0, 1), scope.DefaultHash, defaultDuration}, // because default value in packageConfiguration
		{cfgmock.NewService().NewScoped(1, 1), scope.DefaultHash, defaultDuration}, // because default value in packageConfiguration
		{cfgmock.NewService(cfgmock.WithPV(cfgmock.PathValue{wantPath.Bind(scope.Website, 10).String(): defaultDuration * (time.Second * 2)})).NewScoped(10, 0), scope.NewHash(scope.Website, 10), defaultDuration * (time.Second * 2)},
		{cfgmock.NewService(cfgmock.WithPV(cfgmock.PathValue{wantPath.Bind(scope.Website, 10).String(): defaultDuration * (time.Second * 3)})).NewScoped(10, 1), scope.NewHash(scope.Website, 10), defaultDuration * (time.Second * 3)},
		{cfgmock.NewService(cfgmock.WithPV(cfgmock.PathValue{
			wantPath.String():                       defaultDuration * (time.Second * 5),
			wantPath.Bind(scope.Store, 11).String(): defaultDuration * (time.Second * 6),
		})).NewScoped(10, 11), scope.NewHash(scope.Store, 11), defaultDuration * (time.Second * 6)},
	}
	for i, test := range tests {
		gb, h, err := tm.Get(test.sg)
		if err != nil {
			t.Fatal("Index", i, err)
		}
		assert.Exactly(t, test.want, gb, "Index %d", i)
		assert.Exactly(t, test.wantHash.String(), h.String(), "Index %d", i)
	}
}

func TestDurationGetWithoutCfgStruct(t *testing.T) {

	const pathWebCorsDuration = "web/cors/duration"
	b := cfgmodel.NewDuration(pathWebCorsDuration)
	assert.Empty(t, b.Options())

	wantPath := cfgpath.MustNewByParts(pathWebCorsDuration).Bind(scope.Website, 10)
	defaultDuration := mustParseDuration("2h44m")
	tests := []struct {
		sg       config.Scoped
		wantHash scope.Hash
		want     time.Duration
	}{
		{cfgmock.NewService().NewScoped(1, 1), scope.DefaultHash, 0}, // because default value in packageConfiguration
		{cfgmock.NewService(cfgmock.WithPV(cfgmock.PathValue{wantPath.String(): defaultDuration * (time.Second * 2)})).NewScoped(10, 0), scope.DefaultHash, 0},
		{cfgmock.NewService(cfgmock.WithPV(cfgmock.PathValue{wantPath.String(): defaultDuration * (time.Second * 3)})).NewScoped(10, 1), scope.DefaultHash, 0},
		{cfgmock.NewService(cfgmock.WithPV(cfgmock.PathValue{wantPath.Bind(scope.Default, 0).String(): defaultDuration * (time.Second * 3)})).NewScoped(0, 0), scope.DefaultHash, defaultDuration * (time.Second * 3)},
		{cfgmock.NewService(cfgmock.WithPV(cfgmock.PathValue{
			wantPath.Bind(scope.Default, 0).String(): defaultDuration * (time.Second * 5),
			wantPath.Bind(scope.Store, 11).String():  defaultDuration * (time.Second * 6),
		})).NewScoped(10, 11), scope.DefaultHash, defaultDuration * (time.Second * 5)},
	}
	for i, test := range tests {
		gb, h, err := b.Get(test.sg)
		if err != nil {
			t.Fatal("Index", i, err)
		}
		assert.Exactly(t, test.want, gb, "Index %d", i)
		assert.Exactly(t, test.wantHash.String(), h.String(), "Index %d", i)
	}
}

func TestDurationGetWithoutCfgStructShouldReturnUnexpectedError(t *testing.T) {

	b := cfgmodel.NewDuration("web/cors/duration")
	assert.Empty(t, b.Options())

	gb, h, haveErr := b.Get(cfgmock.NewService(
		cfgmock.WithString(func(path string) (string, error) {
			return "", errors.NewFatalf("Unexpected error")
		}),
	).NewScoped(1, 1))
	assert.Exactly(t, time.Duration(0), gb)
	assert.True(t, errors.IsFatal(haveErr), "Error: %s", haveErr)
	assert.Exactly(t, scope.DefaultHash.String(), h.String())
}

func TestDurationIgnoreNilDefaultValues(t *testing.T) {

	b := cfgmodel.NewDuration("web/cors/duration", cfgmodel.WithField(nil))
	gb, h, err := b.Get(cfgmock.NewService().NewScoped(1, 1))
	if err != nil {
		t.Fatal(err)
	}
	assert.Exactly(t, time.Duration(0), gb)
	assert.Exactly(t, scope.DefaultHash.String(), h.String())
}

func TestDurationWrite(t *testing.T) {

	const pathWebCorsF64 = "web/cors/duration"
	wantPath := cfgpath.MustNewByParts(pathWebCorsF64).Bind(scope.Website, 10)
	haveDuration := mustParseDuration("4h33m")

	b := cfgmodel.NewDuration("web/cors/duration", cfgmodel.WithFieldFromSectionSlice(configStructure))

	mw := &cfgmock.Write{}
	assert.NoError(t, b.Write(mw, haveDuration, scope.Website, 10))
	assert.Exactly(t, wantPath.String(), mw.ArgPath)
	assert.Exactly(t, haveDuration.String(), mw.ArgValue.(string))
}
