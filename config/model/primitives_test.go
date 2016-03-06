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

	"time"

	"github.com/corestoreio/csfw/config"
	"github.com/corestoreio/csfw/config/element"
	"github.com/corestoreio/csfw/config/mock"
	"github.com/corestoreio/csfw/config/model"
	"github.com/corestoreio/csfw/config/path"
	"github.com/corestoreio/csfw/config/source"
	"github.com/corestoreio/csfw/storage/text"
	"github.com/corestoreio/csfw/store/scope"
	"github.com/corestoreio/csfw/util/conv"
	"github.com/corestoreio/csfw/util/cserr"
	"github.com/juju/errors"
	"github.com/stretchr/testify/assert"
)

// configStructure might be a duplicate of base_test but note that the
// test package names are different.
var configStructure = element.MustNewConfiguration(
	&element.Section{
		ID: path.NewRoute("web"),
		Groups: element.NewGroupSlice(
			&element.Group{
				ID:        path.NewRoute("cors"),
				Label:     text.Chars(`CORS Cross Origin Resource Sharing`),
				SortOrder: 150,
				Scopes:    scope.PermDefault,
				Fields: element.NewFieldSlice(
					&element.Field{
						// Path: `web/cors/exposed_headers`,
						ID:        path.NewRoute("exposed_headers"),
						Label:     text.Chars(`Exposed Headers`),
						Comment:   text.Chars(`Indicates which headers are safe to expose to the API of a CORS API specification. Separate via line break`),
						Type:      element.TypeTextarea,
						SortOrder: 10,
						Visible:   element.VisibleYes,
						Scopes:    scope.PermWebsite,
						Default:   "Content-Type,X-CoreStore-ID",
					},
					&element.Field{
						// Path: `web/cors/allowed_origins`,
						ID:        path.NewRoute("allowed_origins"),
						Label:     text.Chars(`Allowed Origins`),
						Comment:   text.Chars(`Is a list of origins a cross-domain request can be executed from.`),
						Type:      element.TypeTextarea,
						SortOrder: 20,
						Visible:   element.VisibleYes,
						Scopes:    scope.PermWebsite,
						Default:   "corestore.io,cs.io",
					},
					&element.Field{
						// Path: `web/cors/allow_credentials`,
						ID:        path.NewRoute("allow_credentials"),
						Label:     text.Chars(`Allowed Credentials`),
						Type:      element.TypeSelect,
						SortOrder: 30,
						Visible:   element.VisibleYes,
						Scopes:    scope.PermWebsite,
						Default:   "true",
					},
					&element.Field{
						// Path: `web/cors/int`,
						ID:        path.NewRoute("int"),
						Type:      element.TypeText,
						SortOrder: 30,
						Visible:   element.VisibleYes,
						Scopes:    scope.PermWebsite,
						Default:   2015,
					},
					&element.Field{
						// Path: `web/cors/int_slice`,
						ID:        path.NewRoute("int_slice"),
						Type:      element.TypeSelect,
						SortOrder: 30,
						Visible:   element.VisibleYes,
						Scopes:    scope.PermStore,
						Default:   "2014,2015,2016",
					},
					&element.Field{
						// Path: `web/cors/float64`,
						ID:        path.NewRoute("float64"),
						Type:      element.TypeText,
						SortOrder: 50,
						Visible:   element.VisibleYes,
						Scopes:    scope.PermWebsite,
						Default:   2015.1000001,
					},
					&element.Field{
						// Path: `web/cors/time`,
						ID:        path.NewRoute("time"),
						Type:      element.TypeText,
						SortOrder: 90,
						Visible:   element.VisibleYes,
						Scopes:    scope.PermStore,
						Default:   "2012-08-23 09:20:13",
					},
				),
			},

			&element.Group{
				ID:        path.NewRoute("unsecure"),
				Label:     text.Chars(`Base URLs`),
				Comment:   text.Chars(`Any of the fields allow fully qualified URLs that end with '/' (slash) e.g. http://example.com/magento/`),
				SortOrder: 10,
				Scopes:    scope.PermStore,
				Fields: element.NewFieldSlice(
					&element.Field{
						// Path: `web/unsecure/base_url`,
						ID:        path.NewRoute("base_url"),
						Label:     text.Chars(`Base URL`),
						Comment:   text.Chars(`Specify URL or {{base_url}} placeholder.`),
						Type:      element.TypeText,
						SortOrder: 10,
						Visible:   element.VisibleYes,
						Scopes:    scope.PermStore,
						Default:   "{{base_url}}",
						//BackendModel: nil, // Magento\Config\Model\Config\Backend\Baseurl
					},

					&element.Field{
						// Path: `web/unsecure/base_link_url`,
						ID:        path.NewRoute("base_link_url"),
						Label:     text.Chars(`Base Link URL`),
						Comment:   text.Chars(`May start with {{unsecure_base_url}} placeholder.`),
						Type:      element.TypeText,
						SortOrder: 20,
						Visible:   element.VisibleYes,
						Scopes:    scope.PermStore,
						Default:   "{{unsecure_base_url}}",
						//BackendModel: nil, // Magento\Config\Model\Config\Backend\Baseurl
					},

					&element.Field{
						// Path: `web/unsecure/base_static_url`,
						ID:        path.NewRoute("base_static_url"),
						Label:     text.Chars(`Base URL for Static View Files`),
						Comment:   text.Chars(`May be empty or start with {{unsecure_base_url}} placeholder.`),
						Type:      element.TypeText,
						SortOrder: 25,
						Visible:   element.VisibleYes,
						Scopes:    scope.PermStore,
						Default:   nil,
						//BackendModel: nil, // Magento\Config\Model\Config\Backend\Baseurl
					},
				),
			},
		),
	},
)

func TestBoolGetWithCfgStruct(t *testing.T) {
	t.Parallel()
	const pathWebCorsCred = "web/cors/allow_credentials"
	wantPath := path.MustNewByParts(pathWebCorsCred).Bind(scope.WebsiteID, 3)
	b := model.NewBool(pathWebCorsCred, model.WithFieldFromSectionSlice(configStructure), model.WithSource(source.YesNo))

	assert.Exactly(t, source.YesNo, b.Options())

	tests := []struct {
		sg   config.ScopedGetter
		want bool
	}{
		{mock.NewService().NewScoped(0, 0), true}, // because default value in packageConfiguration is "true"
		{mock.NewService().NewScoped(5, 4), true}, // because default value in packageConfiguration is "true"
		{mock.NewService(mock.WithPV(mock.PathValue{wantPath.String(): 0})).NewScoped(3, 0), false},
	}
	for i, test := range tests {
		gb, err := b.Get(test.sg)
		if err != nil {
			t.Fatal("Index", i, err)
		}
		assert.Exactly(t, test.want, gb, "Index %d", i)
	}
}

func TestBoolGetWithoutCfgStruct(t *testing.T) {
	t.Parallel()
	const pathWebCorsCred = "web/cors/allow_credentials"
	wantPath := path.MustNewByParts(pathWebCorsCred).Bind(scope.WebsiteID, 4)
	b := model.NewBool(pathWebCorsCred)

	tests := []struct {
		sg   config.ScopedGetter
		want bool
	}{
		{mock.NewService().NewScoped(0, 0), false},
		{mock.NewService().NewScoped(5, 4), false},
		{mock.NewService(mock.WithPV(mock.PathValue{wantPath.String(): 1})).NewScoped(4, 0), false}, // not allowed because DefaultID scope because there has not been set a *element.Field!
		{mock.NewService(mock.WithPV(mock.PathValue{wantPath.Bind(scope.DefaultID, 0).String(): 1})).NewScoped(4, 0), true},
	}
	for i, test := range tests {
		gb, err := b.Get(test.sg)
		if err != nil {
			t.Fatal("Index", i, err)
		}
		assert.Exactly(t, test.want, gb, "Index %d", i)
	}
}

func TestBoolGetWithoutCfgStructShouldReturnUnexpectedError(t *testing.T) {
	t.Parallel()

	b := model.NewBool("web/cors/allow_credentials")
	haveErr := errors.New("Unexpected error")

	gb, err := b.Get(mock.NewService(
		mock.WithBool(func(path string) (bool, error) {
			return false, haveErr
		}),
	).NewScoped(1, 1))
	assert.Empty(t, gb)
	assert.Exactly(t, haveErr, cserr.UnwrapMasked(err))
}

func TestBoolWrite(t *testing.T) {
	t.Parallel()
	const pathWebCorsCred = "web/cors/allow_credentials"
	wantPath := path.MustNewByParts(pathWebCorsCred).Bind(scope.WebsiteID, 3)
	b := model.NewBool(pathWebCorsCred, model.WithFieldFromSectionSlice(configStructure), model.WithSource(source.YesNo))

	mw := &mock.Write{}
	assert.EqualError(t, b.Write(mw, true, scope.StoreID, 3), "Scope permission insufficient: Have 'Store'; Want 'Default,Website'")
	assert.NoError(t, b.Write(mw, true, scope.WebsiteID, 3))
	assert.Exactly(t, wantPath.String(), mw.ArgPath)
	assert.Exactly(t, true, mw.ArgValue.(bool))
}

func TestStrGetWithCfgStruct(t *testing.T) {
	t.Parallel()
	const pathWebCorsHeaders = "web/cors/exposed_headers"
	b := model.NewStr(pathWebCorsHeaders, model.WithFieldFromSectionSlice(configStructure))
	assert.Empty(t, b.Options())

	wantPath := path.MustNewByParts(pathWebCorsHeaders)
	tests := []struct {
		sg   config.ScopedGetter
		want string
	}{
		{mock.NewService().NewScoped(0, 0), "Content-Type,X-CoreStore-ID"}, // because default value in packageConfiguration
		{mock.NewService().NewScoped(5, 4), "Content-Type,X-CoreStore-ID"}, // because default value in packageConfiguration
		{mock.NewService(mock.WithPV(mock.PathValue{wantPath.String(): "X-Gopher"})).NewScoped(0, 0), "X-Gopher"},
		{mock.NewService(mock.WithPV(mock.PathValue{wantPath.String(): "X-Gopher"})).NewScoped(3, 5), "X-Gopher"},
		{mock.NewService(mock.WithPV(mock.PathValue{
			wantPath.String():                         "X-Gopher262",
			wantPath.Bind(scope.StoreID, 44).String(): "X-Gopher44", // because Field.Scopes has PermWebsite
		})).NewScoped(3, 44), "X-Gopher262"},
		{mock.NewService(mock.WithPV(mock.PathValue{
			wantPath.String():                           "X-Gopher",
			wantPath.Bind(scope.WebsiteID, 33).String(): "X-Gopher33",
			wantPath.Bind(scope.WebsiteID, 43).String(): "X-GopherW43",
			wantPath.Bind(scope.StoreID, 44).String():   "X-Gopher44",
		})).NewScoped(33, 43), "X-Gopher33"},
	}
	for i, test := range tests {
		gb, err := b.Get(test.sg)
		if err != nil {
			t.Fatal("Index", i, err)
		}
		assert.Exactly(t, test.want, gb, "Index %d", i)
	}
}

func TestStrGetWithoutCfgStruct(t *testing.T) {
	t.Parallel()
	const pathWebCorsHeaders = "web/cors/exposed_headers"
	b := model.NewStr(pathWebCorsHeaders)
	assert.Empty(t, b.Options())

	wantPath := path.MustNewByParts(pathWebCorsHeaders)
	tests := []struct {
		sg   config.ScopedGetter
		want string
	}{
		{mock.NewService().NewScoped(0, 0), ""},
		{mock.NewService().NewScoped(5, 4), ""},
		{mock.NewService(mock.WithPV(mock.PathValue{wantPath.String(): "X-Gopher"})).NewScoped(0, 0), "X-Gopher"},
	}
	for i, test := range tests {
		gb, err := b.Get(test.sg)
		if err != nil {
			t.Fatal("Index", i, err)
		}
		assert.Exactly(t, test.want, gb, "Index %d", i)
	}
}

func TestStrGetWithoutCfgStructShouldReturnUnexpectedError(t *testing.T) {
	t.Parallel()

	b := model.NewStr("web/cors/exposed_headers")
	assert.Empty(t, b.Options())

	haveErr := errors.New("Unexpected error")
	gb, err := b.Get(mock.NewService(
		mock.WithString(func(path string) (string, error) {
			return "", haveErr
		}),
	).NewScoped(1, 1))
	assert.Empty(t, gb)
	assert.Exactly(t, haveErr, cserr.UnwrapMasked(err))
}

func TestStrWrite(t *testing.T) {
	t.Parallel()
	const pathWebCorsHeaders = "web/cors/exposed_headers"
	wantPath := path.MustNewByParts(pathWebCorsHeaders)
	b := model.NewStr(pathWebCorsHeaders, model.WithFieldFromSectionSlice(configStructure))

	mw := &mock.Write{}
	assert.NoError(t, b.Write(mw, "dude", scope.DefaultID, 0))
	assert.Exactly(t, wantPath.String(), mw.ArgPath)
	assert.Exactly(t, "dude", mw.ArgValue.(string))
}

func TestIntGetWithCfgStruct(t *testing.T) {
	t.Parallel()
	const pathWebCorsInt = "web/cors/int"
	b := model.NewInt(pathWebCorsInt, model.WithFieldFromSectionSlice(configStructure))
	assert.Empty(t, b.Options())

	wantPath := path.MustNewByParts(pathWebCorsInt)
	tests := []struct {
		sg   config.ScopedGetter
		want int
	}{
		{mock.NewService().NewScoped(0, 0), 2015}, // because default value in packageConfiguration
		{mock.NewService().NewScoped(0, 1), 2015}, // because default value in packageConfiguration
		{mock.NewService().NewScoped(1, 1), 2015}, // because default value in packageConfiguration
		{mock.NewService(mock.WithPV(mock.PathValue{wantPath.Bind(scope.WebsiteID, 10).String(): 2016})).NewScoped(10, 0), 2016},
		{mock.NewService(mock.WithPV(mock.PathValue{wantPath.Bind(scope.WebsiteID, 10).String(): 2016})).NewScoped(10, 1), 2016},
		{mock.NewService(mock.WithPV(mock.PathValue{
			wantPath.String():                         3017,
			wantPath.Bind(scope.StoreID, 11).String(): 2016, // because Field.Scopes set to PermWebsite
		})).NewScoped(10, 11), 3017},
		{mock.NewService(mock.WithPV(mock.PathValue{
			wantPath.String():                           3017,
			wantPath.Bind(scope.WebsiteID, 10).String(): 4018,
			wantPath.Bind(scope.StoreID, 11).String():   2016, // because Field.Scopes set to PermWebsite
		})).NewScoped(10, 11), 4018},
	}
	for i, test := range tests {
		gb, err := b.Get(test.sg)
		if err != nil {
			t.Fatal("Index", i, err)
		}
		assert.Exactly(t, test.want, gb, "Index %d", i)
	}
}

func TestIntGetWithoutCfgStruct(t *testing.T) {
	t.Parallel()

	const pathWebCorsInt = "web/cors/int"
	b := model.NewInt(pathWebCorsInt) // no *element.Field has been set. So Default Scope will be enforced
	assert.Empty(t, b.Options())

	wantPath := path.MustNewByParts(pathWebCorsInt)
	tests := []struct {
		sg   config.ScopedGetter
		want int
	}{
		{mock.NewService().NewScoped(1, 1), 0},
		{mock.NewService(mock.WithPV(mock.PathValue{wantPath.Bind(scope.WebsiteID, 10).String(): 2016})).NewScoped(10, 0), 0},
		{mock.NewService(mock.WithPV(mock.PathValue{wantPath.Bind(scope.DefaultID, 0).String(): 2019})).NewScoped(10, 0), 2019},
	}
	for i, test := range tests {
		gb, err := b.Get(test.sg)
		if err != nil {
			t.Fatal("Index", i, err)
		}
		assert.Exactly(t, test.want, gb, "Index %d", i)
	}
}

func TestIntGetWithoutCfgStructShouldReturnUnexpectedError(t *testing.T) {
	t.Parallel()

	b := model.NewInt("web/cors/int")

	haveErr := errors.New("Unexpected error")
	gb, err := b.Get(mock.NewService(
		mock.WithInt(func(path string) (int, error) {
			return 0, haveErr
		}),
	).NewScoped(1, 1))
	assert.Empty(t, gb)
	assert.Exactly(t, haveErr, cserr.UnwrapMasked(err))
}

func TestIntWrite(t *testing.T) {
	t.Parallel()
	const pathWebCorsInt = "web/cors/int"
	wantPath := path.MustNewByParts(pathWebCorsInt).Bind(scope.WebsiteID, 10)
	b := model.NewInt(pathWebCorsInt, model.WithFieldFromSectionSlice(configStructure))

	mw := &mock.Write{}
	assert.NoError(t, b.Write(mw, 27182, scope.WebsiteID, 10))
	assert.Exactly(t, wantPath.String(), mw.ArgPath)
	assert.Exactly(t, 27182, mw.ArgValue.(int))
}

func TestFloat64GetWithCfgStruct(t *testing.T) {
	t.Parallel()
	const pathWebCorsF64 = "web/cors/float64"
	b := model.NewFloat64("web/cors/float64", model.WithFieldFromSectionSlice(configStructure))
	assert.Empty(t, b.Options())

	wantPath := path.MustNewByParts(pathWebCorsF64).Bind(scope.WebsiteID, 10)
	tests := []struct {
		sg   config.ScopedGetter
		want float64
	}{
		{mock.NewService().NewScoped(0, 0), 2015.1000001}, // because default value in packageConfiguration
		{mock.NewService().NewScoped(0, 1), 2015.1000001}, // because default value in packageConfiguration
		{mock.NewService().NewScoped(1, 1), 2015.1000001}, // because default value in packageConfiguration
		{mock.NewService(mock.WithPV(mock.PathValue{wantPath.Bind(scope.WebsiteID, 10).String(): 2016.1000001})).NewScoped(10, 0), 2016.1000001},
		{mock.NewService(mock.WithPV(mock.PathValue{wantPath.Bind(scope.WebsiteID, 10).String(): 2016.1000001})).NewScoped(10, 1), 2016.1000001},
		{mock.NewService(mock.WithPV(mock.PathValue{
			wantPath.String():                         2017.1000001,
			wantPath.Bind(scope.StoreID, 11).String(): 2016.1000021,
		})).NewScoped(10, 11), 2017.1000001},
		{mock.NewService(mock.WithPV(mock.PathValue{
			wantPath.String():                           2017.1000001,
			wantPath.Bind(scope.WebsiteID, 13).String(): 2018.2000001,
			wantPath.Bind(scope.StoreID, 11).String():   2016.1000021,
		})).NewScoped(13, 11), 2018.2000001},
	}
	for i, test := range tests {
		gb, err := b.Get(test.sg)
		if err != nil {
			t.Fatal("Index", i, err)
		}
		assert.Exactly(t, test.want, gb, "Index %d", i)
	}

}

func TestFloat64GetWithoutCfgStruct(t *testing.T) {
	t.Parallel()
	const pathWebCorsF64 = "web/cors/float64"
	b := model.NewFloat64(pathWebCorsF64) // no *element.Field has been set. So Default Scope will be enforced
	assert.Empty(t, b.Options())

	wantPath := path.MustNewByParts(pathWebCorsF64).Bind(scope.WebsiteID, 10)
	tests := []struct {
		sg   config.ScopedGetter
		want float64
	}{
		{mock.NewService().NewScoped(0, 0), 0},
		{mock.NewService(mock.WithPV(mock.PathValue{wantPath.Bind(scope.WebsiteID, 10).String(): 2016.1000001})).NewScoped(10, 0), 0},
		{mock.NewService(mock.WithPV(mock.PathValue{wantPath.Bind(scope.DefaultID, 0).String(): 2016.1000001})).NewScoped(10, 0), 2016.1000001},
	}
	for i, test := range tests {
		gb, err := b.Get(test.sg)
		if err != nil {
			t.Fatal("Index", i, err)
		}
		assert.Exactly(t, test.want, gb, "Index %d", i)
	}
}

func TestFloat64GetWithoutCfgStructShouldReturnUnexpectedError(t *testing.T) {
	t.Parallel()

	b := model.NewFloat64("web/cors/float64")

	haveErr := errors.New("Unexpected error")
	gb, err := b.Get(mock.NewService(

		mock.WithFloat64(func(path string) (float64, error) {
			return 0, haveErr
		}),
	).NewScoped(1, 1))
	assert.Empty(t, gb)
	assert.Exactly(t, haveErr, cserr.UnwrapMasked(err))

}

func TestFloat64Write(t *testing.T) {
	t.Parallel()
	const pathWebCorsF64 = "web/cors/float64"
	wantPath := path.MustNewByParts(pathWebCorsF64).Bind(scope.WebsiteID, 10)
	b := model.NewFloat64("web/cors/float64", model.WithFieldFromSectionSlice(configStructure))

	mw := &mock.Write{}
	assert.NoError(t, b.Write(mw, 1.123456789, scope.WebsiteID, 10))
	assert.Exactly(t, wantPath.String(), mw.ArgPath)
	assert.Exactly(t, 1.12345678900000, mw.ArgValue.(float64))
}

func TestRecursiveOption(t *testing.T) {
	t.Parallel()
	b := model.NewInt(
		"web/cors/int",
		model.WithFieldFromSectionSlice(configStructure),
		model.WithSourceByString("a", "A", "b", "b"),
	)

	assert.Exactly(t, source.NewByString("a", "A", "b", "b"), b.Source)

	previous := b.Option(model.WithSourceByString(
		"1", "One", "2", "Two",
	))
	assert.Exactly(t, source.NewByString("1", "One", "2", "Two"), b.Source)

	b.Option(previous)
	assert.Exactly(t, source.NewByString("a", "A", "b", "b"), b.Source)
}

func mustParseTime(s string) time.Time {
	t, err := conv.StringToDate(s, nil)
	if err != nil {
		panic(err)
	}
	return t
}

func TestTimeGetWithCfgStruct(t *testing.T) {
	t.Parallel()
	const pathWebCorsTime = "web/cors/time"
	tm := model.NewTime("web/cors/time", model.WithFieldFromSectionSlice(configStructure))
	assert.Empty(t, tm.Options())

	wantPath := path.MustNewByParts(pathWebCorsTime).Bind(scope.WebsiteID, 10)
	defaultTime := mustParseTime("2012-08-23 09:20:13")
	tests := []struct {
		sg   config.ScopedGetter
		want time.Time
	}{
		{mock.NewService().NewScoped(0, 0), defaultTime}, // because default value in packageConfiguration
		{mock.NewService().NewScoped(0, 1), defaultTime}, // because default value in packageConfiguration
		{mock.NewService().NewScoped(1, 1), defaultTime}, // because default value in packageConfiguration
		{mock.NewService(mock.WithPV(mock.PathValue{wantPath.Bind(scope.WebsiteID, 10).String(): defaultTime.Add(time.Second * 2)})).NewScoped(10, 0), defaultTime.Add(time.Second * 2)},
		{mock.NewService(mock.WithPV(mock.PathValue{wantPath.Bind(scope.WebsiteID, 10).String(): defaultTime.Add(time.Second * 3)})).NewScoped(10, 1), defaultTime.Add(time.Second * 3)},
		{mock.NewService(mock.WithPV(mock.PathValue{
			wantPath.String():                         defaultTime.Add(time.Second * 5),
			wantPath.Bind(scope.StoreID, 11).String(): defaultTime.Add(time.Second * 6),
		})).NewScoped(10, 11), defaultTime.Add(time.Second * 6)},
	}
	for i, test := range tests {
		gb, err := tm.Get(test.sg)
		if err != nil {
			t.Fatal("Index", i, err)
		}
		assert.Exactly(t, test.want, gb, "Index %d", i)
	}
}

func TestTimeGetWithoutCfgStruct(t *testing.T) {
	t.Parallel()
	const pathWebCorsTime = "web/cors/time"
	b := model.NewTime(pathWebCorsTime)
	assert.Empty(t, b.Options())

	wantPath := path.MustNewByParts(pathWebCorsTime).Bind(scope.WebsiteID, 10)
	defaultTime := mustParseTime("2012-08-23 09:20:13")
	tests := []struct {
		sg   config.ScopedGetter
		want time.Time
	}{
		{mock.NewService().NewScoped(1, 1), time.Time{}}, // because default value in packageConfiguration
		{mock.NewService(mock.WithPV(mock.PathValue{wantPath.String(): defaultTime.Add(time.Second * 2)})).NewScoped(10, 0), time.Time{}},
		{mock.NewService(mock.WithPV(mock.PathValue{wantPath.String(): defaultTime.Add(time.Second * 3)})).NewScoped(10, 1), time.Time{}},
		{mock.NewService(mock.WithPV(mock.PathValue{wantPath.Bind(scope.DefaultID, 0).String(): defaultTime.Add(time.Second * 3)})).NewScoped(0, 0), defaultTime.Add(time.Second * 3)},
		{mock.NewService(mock.WithPV(mock.PathValue{
			wantPath.Bind(scope.DefaultID, 0).String(): defaultTime.Add(time.Second * 5),
			wantPath.Bind(scope.StoreID, 11).String():  defaultTime.Add(time.Second * 6),
		})).NewScoped(10, 11), defaultTime.Add(time.Second * 5)},
	}
	for i, test := range tests {
		gb, err := b.Get(test.sg)
		if err != nil {
			t.Fatal("Index", i, err)
		}
		assert.Exactly(t, test.want, gb, "Index %d", i)
	}
}

func TestTimeGetWithoutCfgStructShouldReturnUnexpectedError(t *testing.T) {
	t.Parallel()

	b := model.NewTime("web/cors/time")
	assert.Empty(t, b.Options())

	haveErr := errors.New("Unexpected error")
	gb, err := b.Get(mock.NewService(
		mock.WithTime(func(path string) (time.Time, error) {
			return time.Time{}, haveErr
		}),
	).NewScoped(1, 1))
	assert.Empty(t, gb)
	assert.Exactly(t, haveErr, cserr.UnwrapMasked(err))
}

func TestTimeWrite(t *testing.T) {
	t.Parallel()
	const pathWebCorsF64 = "web/cors/time"
	wantPath := path.MustNewByParts(pathWebCorsF64).Bind(scope.WebsiteID, 10)
	haveTime := mustParseTime("2000-08-23 09:20:13")

	b := model.NewTime("web/cors/time", model.WithFieldFromSectionSlice(configStructure))

	mw := &mock.Write{}
	assert.NoError(t, b.Write(mw, haveTime, scope.WebsiteID, 10))
	assert.Exactly(t, wantPath.String(), mw.ArgPath)
	assert.Exactly(t, haveTime, mw.ArgValue.(time.Time))
}
