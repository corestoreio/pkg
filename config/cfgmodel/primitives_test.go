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

	"github.com/corestoreio/csfw/config"
	"github.com/corestoreio/csfw/config/cfgmock"
	"github.com/corestoreio/csfw/config/cfgmodel"
	"github.com/corestoreio/csfw/config/cfgpath"
	"github.com/corestoreio/csfw/config/element"
	"github.com/corestoreio/csfw/config/source"
	"github.com/corestoreio/csfw/storage/text"
	"github.com/corestoreio/csfw/store/scope"
	"github.com/corestoreio/csfw/util/errors"
	"github.com/stretchr/testify/assert"
)

// configStructure might be a duplicate of base_test but note that the
// test package names are different.
var configStructure = element.MustNewConfiguration(
	element.Section{
		ID: cfgpath.NewRoute("web"),
		Groups: element.NewGroupSlice(
			element.Group{
				ID:        cfgpath.NewRoute("cors"),
				Label:     text.Chars(`CORS Cross Origin Resource Sharing`),
				SortOrder: 150,
				Scopes:    scope.PermDefault,
				Fields: element.NewFieldSlice(
					element.Field{
						// Path: `web/cors/exposed_headers`,
						ID:        cfgpath.NewRoute("exposed_headers"),
						Label:     text.Chars(`Exposed Headers`),
						Comment:   text.Chars(`Indicates which headers are safe to expose to the API of a CORS API specification. Separate via line break`),
						Type:      element.TypeTextarea,
						SortOrder: 10,
						Visible:   element.VisibleYes,
						Scopes:    scope.PermWebsite,
						Default:   "Content-Type,X-CoreStore-ID",
					},
					element.Field{
						// Path: `web/cors/allowed_origins`,
						ID:        cfgpath.NewRoute("allowed_origins"),
						Label:     text.Chars(`Allowed Origins`),
						Comment:   text.Chars(`Is a list of origins a cross-domain request can be executed from.`),
						Type:      element.TypeTextarea,
						SortOrder: 20,
						Visible:   element.VisibleYes,
						Scopes:    scope.PermWebsite,
						Default:   "corestore.io,cs.io",
					},
					element.Field{
						// Path: `web/cors/allow_credentials`,
						ID:        cfgpath.NewRoute("allow_credentials"),
						Label:     text.Chars(`Allowed Credentials`),
						Type:      element.TypeSelect,
						SortOrder: 30,
						Visible:   element.VisibleYes,
						Scopes:    scope.PermWebsite,
						Default:   "true",
					},
					element.Field{
						// Path: `web/cors/int`,
						ID:        cfgpath.NewRoute("int"),
						Type:      element.TypeText,
						SortOrder: 30,
						Visible:   element.VisibleYes,
						Scopes:    scope.PermWebsite,
						Default:   2015,
					},
					element.Field{
						// Path: `web/cors/int_slice`,
						ID:        cfgpath.NewRoute("int_slice"),
						Type:      element.TypeSelect,
						SortOrder: 30,
						Visible:   element.VisibleYes,
						Scopes:    scope.PermStore,
						Default:   "2014,2015,2016",
					},
					element.Field{
						// Path: `web/cors/float64`,
						ID:        cfgpath.NewRoute("float64"),
						Type:      element.TypeText,
						SortOrder: 50,
						Visible:   element.VisibleYes,
						Scopes:    scope.PermWebsite,
						Default:   2015.1000001,
					},
					element.Field{
						// Path: `web/cors/time`,
						ID:        cfgpath.NewRoute("time"),
						Type:      element.TypeText,
						SortOrder: 90,
						Visible:   element.VisibleYes,
						Scopes:    scope.PermStore,
						Default:   "2012-08-23 09:20:13",
					},
					element.Field{
						// Path: `web/cors/duration`,
						ID:        cfgpath.NewRoute("duration"),
						Type:      element.TypeText,
						SortOrder: 100,
						Visible:   element.VisibleYes,
						Scopes:    scope.PermStore,
						Default:   "1h45m",
					},
					element.Field{
						// Path: `web/cors/byte`,
						ID:        cfgpath.NewRoute("byte"),
						Type:      element.TypeText,
						SortOrder: 110,
						Visible:   element.VisibleYes,
						Scopes:    scope.PermWebsite,
						Default:   []byte(`Hello Dud€`),
					},
					element.Field{
						// Path: `web/cors/csv`,
						ID:        cfgpath.NewRoute("csv"),
						Type:      element.TypeTextarea,
						SortOrder: 30,
						Visible:   element.VisibleYes,
						Scopes:    scope.PermWebsite,
						Default: `0|"""Did you mean..."" Suggestions"|"""meinten Sie...?"""
1|"Accuracy for Suggestions"|"Genauigkeit der Vorschläge"
2|"After switching please reindex the<br /><em>Catalog Search Index</em>."|"Nach dem Umschalten reindexieren Sie bitte den <br /><em>Katalog Suchindex</em>."
3|"CATALOG"|"KATALOG"`,
					},
				),
			},

			element.Group{
				ID:        cfgpath.NewRoute("unsecure"),
				Label:     text.Chars(`Base URLs`),
				Comment:   text.Chars(`Any of the fields allow fully qualified URLs that end with '/' (slash) e.g. http://example.com/magento/`),
				SortOrder: 10,
				Scopes:    scope.PermStore,
				Fields: element.NewFieldSlice(
					element.Field{
						// Path: `web/unsecure/url`,
						ID:        cfgpath.NewRoute("url"),
						Label:     text.Chars(`Just an URL`),
						Type:      element.TypeText,
						SortOrder: 9,
						Visible:   element.VisibleYes,
						Scopes:    scope.PermStore,
						Default:   `http://john%20doe@corestore.io/?q=go+language#foo%26bar`,
					},
					element.Field{
						// Path: `web/unsecure/base_url`,
						ID:        cfgpath.NewRoute("base_url"),
						Label:     text.Chars(`Base URL`),
						Comment:   text.Chars(`Specify URL or {{base_url}} placeholder.`),
						Type:      element.TypeText,
						SortOrder: 10,
						Visible:   element.VisibleYes,
						Scopes:    scope.PermStore,
						Default:   "{{base_url}}",
						//BackendModel: nil, // Magento\Config\Model\Config\Backend\Baseurl
					},

					element.Field{
						// Path: `web/unsecure/base_link_url`,
						ID:        cfgpath.NewRoute("base_link_url"),
						Label:     text.Chars(`Base Link URL`),
						Comment:   text.Chars(`May start with {{unsecure_base_url}} placeholder.`),
						Type:      element.TypeText,
						SortOrder: 20,
						Visible:   element.VisibleYes,
						Scopes:    scope.PermStore,
						Default:   "{{unsecure_base_url}}",
						//BackendModel: nil, // Magento\Config\Model\Config\Backend\Baseurl
					},

					element.Field{
						// Path: `web/unsecure/base_static_url`,
						ID:        cfgpath.NewRoute("base_static_url"),
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

	const pathWebCorsCred = "web/cors/allow_credentials"
	wantPath := cfgpath.MustNewByParts(pathWebCorsCred).Bind(scope.Website, 3)
	b := cfgmodel.NewBool(pathWebCorsCred, cfgmodel.WithFieldFromSectionSlice(configStructure), cfgmodel.WithSource(source.YesNo))

	assert.Exactly(t, source.YesNo, b.Options())

	tests := []struct {
		sg       config.Scoped
		wantHash scope.Hash
		want     bool
	}{
		{cfgmock.NewService().NewScoped(0, 0), scope.DefaultHash, true}, // because default value in packageConfiguration is "true"
		{cfgmock.NewService().NewScoped(5, 4), scope.DefaultHash, true}, // because default value in packageConfiguration is "true"
		{cfgmock.NewService(cfgmock.WithPV(cfgmock.PathValue{wantPath.String(): 0})).NewScoped(3, 0), scope.NewHash(scope.Website, 3), false},
		{cfgmock.NewService(cfgmock.WithPV(cfgmock.PathValue{wantPath.String(): 0})).NewScoped(3, 5), scope.NewHash(scope.Website, 3), false},
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

func TestBoolGetWithoutCfgStruct(t *testing.T) {

	const pathWebCorsCred = "web/cors/allow_credentials"
	wantPath := cfgpath.MustNewByParts(pathWebCorsCred).Bind(scope.Website, 4)
	b := cfgmodel.NewBool(pathWebCorsCred)

	tests := []struct {
		sg       config.Scoped
		wantHash scope.Hash
		want     bool
	}{
		{cfgmock.NewService().NewScoped(0, 0), scope.DefaultHash, false},
		{cfgmock.NewService().NewScoped(5, 4), scope.DefaultHash, false},
		{cfgmock.NewService(cfgmock.WithPV(cfgmock.PathValue{wantPath.String(): 1})).NewScoped(4, 0), scope.DefaultHash, false}, // not allowed because DefaultID scope because there has not been set a *element.Field!
		{cfgmock.NewService(cfgmock.WithPV(cfgmock.PathValue{wantPath.Bind(scope.Default, 0).String(): 1})).NewScoped(4, 0), scope.DefaultHash, true},
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

func TestBoolGetWithoutCfgStructShouldReturnUnexpectedError(t *testing.T) {

	b := cfgmodel.NewBool("web/cors/allow_credentials")

	gb, h, haveErr := b.Get(cfgmock.NewService(
		cfgmock.WithBool(func(path string) (bool, error) {
			return false, errors.NewFatalf("Unexpected error")
		}),
	).NewScoped(1, 1))
	assert.Empty(t, gb)
	assert.True(t, errors.IsFatal(haveErr), "Error: %s", haveErr)
	assert.Exactly(t, scope.DefaultHash.String(), h.String())
}

func TestBoolIgnoreNilDefaultValues(t *testing.T) {

	b := cfgmodel.NewBool("web/cors/bool", cfgmodel.WithField(nil))
	gb, h, err := b.Get(cfgmock.NewService().NewScoped(0, 0))
	if err != nil {
		t.Fatal(err)
	}
	assert.Exactly(t, false, gb)
	assert.Exactly(t, scope.DefaultHash.String(), h.String())
}

func TestBoolWrite(t *testing.T) {

	const pathWebCorsCred = "web/cors/allow_credentials"
	wantPath := cfgpath.MustNewByParts(pathWebCorsCred).Bind(scope.Website, 3)
	b := cfgmodel.NewBool(pathWebCorsCred, cfgmodel.WithFieldFromSectionSlice(configStructure), cfgmodel.WithSource(source.YesNo))

	mw := &cfgmock.Write{}
	err := b.Write(mw, true, scope.Store, 3)
	assert.True(t, errors.IsUnauthorized(err), "Error: %s", err)
	assert.NoError(t, b.Write(mw, true, scope.Website, 3))
	assert.Exactly(t, wantPath.String(), mw.ArgPath)
	assert.Exactly(t, true, mw.ArgValue.(bool))
}

func TestByteGetWithCfgStruct(t *testing.T) {

	const pathWebCorsByte = "web/cors/byte"
	var defaultWebCorsByte = []byte(`Hello Dud€`)
	b := cfgmodel.NewByte(pathWebCorsByte, cfgmodel.WithFieldFromSectionSlice(configStructure))
	assert.Empty(t, b.Options())

	wantPath := cfgpath.MustNewByParts(pathWebCorsByte)
	tests := []struct {
		sg       config.Scoped
		wantHash scope.Hash
		want     []byte
	}{
		{cfgmock.NewService().NewScoped(0, 0), scope.DefaultHash, defaultWebCorsByte}, // because default value in packageConfiguration
		{cfgmock.NewService().NewScoped(5, 4), scope.DefaultHash, defaultWebCorsByte}, // because default value in packageConfiguration
		{cfgmock.NewService(cfgmock.WithPV(cfgmock.PathValue{wantPath.String(): []byte("X-Gopher")})).NewScoped(0, 0), scope.DefaultHash, []byte("X-Gopher")},
		{cfgmock.NewService(cfgmock.WithPV(cfgmock.PathValue{wantPath.String(): []byte("X-Gopher")})).NewScoped(3, 5), scope.DefaultHash, []byte("X-Gopher")},
		{cfgmock.NewService(cfgmock.WithPV(cfgmock.PathValue{
			wantPath.String():                       []byte("X-Gopher262"),
			wantPath.Bind(scope.Store, 44).String(): []byte("X-Gopher44"), // because Field.Scopes has PermWebsite
		})).NewScoped(3, 44), scope.DefaultHash, []byte("X-Gopher262")},
		{cfgmock.NewService(cfgmock.WithPV(cfgmock.PathValue{
			wantPath.String():                         []byte("X-Gopher"),
			wantPath.Bind(scope.Website, 33).String(): []byte("X-Gopher33"),
			wantPath.Bind(scope.Website, 43).String(): []byte("X-GopherW43"),
			wantPath.Bind(scope.Store, 44).String():   []byte("X-Gopher44"),
		})).NewScoped(33, 43), scope.NewHash(scope.Website, 33), []byte("X-Gopher33")},
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

func TestByteGetWithoutCfgStruct(t *testing.T) {

	const pathWebCorsHeaders = "web/cors/byte"
	b := cfgmodel.NewByte(pathWebCorsHeaders)
	assert.Empty(t, b.Options())

	wantPath := cfgpath.MustNewByParts(pathWebCorsHeaders)
	tests := []struct {
		sg       config.Scoped
		wantHash scope.Hash
		want     []byte
	}{
		{cfgmock.NewService().NewScoped(0, 0), scope.DefaultHash, nil},
		{cfgmock.NewService().NewScoped(5, 4), scope.DefaultHash, nil},
		{cfgmock.NewService(cfgmock.WithPV(cfgmock.PathValue{wantPath.String(): []byte(`Hello Dud€`)})).NewScoped(0, 0), scope.DefaultHash, []byte(`Hello Dud€`)},
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

func TestByteGetWithoutCfgStructShouldReturnUnexpectedError(t *testing.T) {

	b := cfgmodel.NewByte("web/cors/byte")
	assert.Empty(t, b.Options())

	gb, h, haveErr := b.Get(cfgmock.NewService(
		cfgmock.WithByte(func(path string) ([]byte, error) {
			return nil, errors.NewFatalf("Unexpected error")
		}),
	).NewScoped(1, 1))
	assert.Empty(t, gb)
	assert.True(t, errors.IsFatal(haveErr), "Error: %s", haveErr)
	assert.Exactly(t, scope.DefaultHash.String(), h.String())
}

func TestByteIgnoreNilDefaultValues(t *testing.T) {

	b := cfgmodel.NewByte("web/cors/byte", cfgmodel.WithField(&element.Field{}))
	gb, h, err := b.Get(cfgmock.NewService().NewScoped(0, 0))
	if err != nil {
		t.Fatal(err)
	}
	assert.Exactly(t, []byte(nil), gb)
	assert.Exactly(t, scope.DefaultHash.String(), h.String())
}
func TestByteWrite(t *testing.T) {

	const pathWebCorsHeaders = "web/cors/byte"
	wantPath := cfgpath.MustNewByParts(pathWebCorsHeaders)
	b := cfgmodel.NewByte(pathWebCorsHeaders, cfgmodel.WithFieldFromSectionSlice(configStructure))

	mw := &cfgmock.Write{}
	assert.NoError(t, b.Write(mw, []byte("dude"), scope.Default, 0))
	assert.Exactly(t, wantPath.String(), mw.ArgPath)
	assert.Exactly(t, []byte("dude"), mw.ArgValue.([]byte))
}

func TestStrGetWithCfgStruct(t *testing.T) {

	const pathWebCorsHeaders = "web/cors/exposed_headers"
	b := cfgmodel.NewStr(pathWebCorsHeaders, cfgmodel.WithFieldFromSectionSlice(configStructure))
	assert.Empty(t, b.Options())

	wantPath := cfgpath.MustNewByParts(pathWebCorsHeaders)
	tests := []struct {
		sg       config.Scoped
		wantHash scope.Hash
		want     string
	}{
		{cfgmock.NewService().NewScoped(0, 0), scope.DefaultHash, "Content-Type,X-CoreStore-ID"}, // because default value in packageConfiguration
		{cfgmock.NewService().NewScoped(5, 4), scope.DefaultHash, "Content-Type,X-CoreStore-ID"}, // because default value in packageConfiguration
		{cfgmock.NewService(cfgmock.WithPV(cfgmock.PathValue{wantPath.String(): "X-Gopher"})).NewScoped(0, 0), scope.DefaultHash, "X-Gopher"},
		{cfgmock.NewService(cfgmock.WithPV(cfgmock.PathValue{wantPath.String(): "X-Gopher"})).NewScoped(3, 5), scope.DefaultHash, "X-Gopher"},
		{cfgmock.NewService(cfgmock.WithPV(cfgmock.PathValue{
			wantPath.String():                       "X-Gopher262",
			wantPath.Bind(scope.Store, 44).String(): "X-Gopher44", // because Field.Scopes has PermWebsite
		})).NewScoped(3, 44), scope.DefaultHash, "X-Gopher262"},
		{cfgmock.NewService(cfgmock.WithPV(cfgmock.PathValue{
			wantPath.String():                         "X-Gopher",
			wantPath.Bind(scope.Website, 33).String(): "X-Gopher33",
			wantPath.Bind(scope.Website, 43).String(): "X-GopherW43",
			wantPath.Bind(scope.Store, 44).String():   "X-Gopher44",
		})).NewScoped(33, 43), scope.NewHash(scope.Website, 33), "X-Gopher33"},
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

func TestStrGetWithoutCfgStruct(t *testing.T) {

	const pathWebCorsHeaders = "web/cors/exposed_headers"
	b := cfgmodel.NewStr(pathWebCorsHeaders)
	assert.Empty(t, b.Options())

	wantPath := cfgpath.MustNewByParts(pathWebCorsHeaders)
	tests := []struct {
		sg       config.Scoped
		wantHash scope.Hash
		want     string
	}{
		{cfgmock.NewService().NewScoped(0, 0), scope.DefaultHash, ""},
		{cfgmock.NewService().NewScoped(5, 4), scope.DefaultHash, ""},
		{cfgmock.NewService(cfgmock.WithPV(cfgmock.PathValue{wantPath.String(): "X-Gopher"})).NewScoped(0, 0), scope.DefaultHash, "X-Gopher"},
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

func TestStrGetWithoutCfgStructShouldReturnUnexpectedError(t *testing.T) {

	b := cfgmodel.NewStr("web/cors/exposed_headers")
	assert.Empty(t, b.Options())

	gb, h, haveErr := b.Get(cfgmock.NewService(
		cfgmock.WithString(func(path string) (string, error) {
			return "", errors.NewFatalf("Unexpected error")
		}),
	).NewScoped(1, 1))
	assert.Empty(t, gb)
	assert.True(t, errors.IsFatal(haveErr), "Error: %s", haveErr)
	assert.Exactly(t, scope.DefaultHash.String(), h.String())
}

func TestStrIgnoreNilDefaultValues(t *testing.T) {

	b := cfgmodel.NewStr("web/cors/str", cfgmodel.WithField(nil))
	gb, h, err := b.Get(cfgmock.NewService().NewScoped(0, 0))
	if err != nil {
		t.Fatal(err)
	}
	assert.Exactly(t, "", gb)
	assert.Exactly(t, scope.DefaultHash.String(), h.String())
}
func TestStrWrite(t *testing.T) {

	const pathWebCorsHeaders = "web/cors/exposed_headers"
	wantPath := cfgpath.MustNewByParts(pathWebCorsHeaders)
	b := cfgmodel.NewStr(pathWebCorsHeaders, cfgmodel.WithFieldFromSectionSlice(configStructure))

	mw := &cfgmock.Write{}
	assert.NoError(t, b.Write(mw, "dude", scope.Default, 0))
	assert.Exactly(t, wantPath.String(), mw.ArgPath)
	assert.Exactly(t, "dude", mw.ArgValue.(string))
}

func TestIntGetWithCfgStruct(t *testing.T) {

	const pathWebCorsInt = "web/cors/int"
	b := cfgmodel.NewInt(pathWebCorsInt, cfgmodel.WithFieldFromSectionSlice(configStructure))
	assert.Empty(t, b.Options())

	wantPath := cfgpath.MustNewByParts(pathWebCorsInt)
	tests := []struct {
		sg       config.Scoped
		wantHash scope.Hash
		want     int
	}{
		{cfgmock.NewService().NewScoped(0, 0), scope.DefaultHash, 2015}, // because default value in packageConfiguration
		{cfgmock.NewService().NewScoped(0, 1), scope.DefaultHash, 2015}, // because default value in packageConfiguration
		{cfgmock.NewService().NewScoped(1, 1), scope.DefaultHash, 2015}, // because default value in packageConfiguration
		{cfgmock.NewService(cfgmock.WithPV(cfgmock.PathValue{wantPath.Bind(scope.Website, 10).String(): 2016})).NewScoped(10, 0), scope.NewHash(scope.Website, 10), 2016},
		{cfgmock.NewService(cfgmock.WithPV(cfgmock.PathValue{wantPath.Bind(scope.Website, 10).String(): 2016})).NewScoped(10, 1), scope.NewHash(scope.Website, 10), 2016},
		{cfgmock.NewService(cfgmock.WithPV(cfgmock.PathValue{
			wantPath.String():                       3017,
			wantPath.Bind(scope.Store, 11).String(): 2016, // because Field.Scopes set to PermWebsite
		})).NewScoped(10, 11), scope.DefaultHash, 3017},
		{cfgmock.NewService(cfgmock.WithPV(cfgmock.PathValue{
			wantPath.String():                         3017,
			wantPath.Bind(scope.Website, 10).String(): 4018,
			wantPath.Bind(scope.Store, 11).String():   2016, // because Field.Scopes set to PermWebsite
		})).NewScoped(10, 11), scope.NewHash(scope.Website, 10), 4018},
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

func TestIntGetWithoutCfgStruct(t *testing.T) {

	const pathWebCorsInt = "web/cors/int"
	b := cfgmodel.NewInt(pathWebCorsInt) // no *element.Field has been set. So Default Scope will be enforced
	assert.Empty(t, b.Options())

	wantPath := cfgpath.MustNewByParts(pathWebCorsInt)
	tests := []struct {
		sg       config.Scoped
		wantHash scope.Hash
		want     int
	}{
		{cfgmock.NewService().NewScoped(1, 1), scope.DefaultHash, 0},
		{cfgmock.NewService(cfgmock.WithPV(cfgmock.PathValue{wantPath.Bind(scope.Website, 10).String(): 2016})).NewScoped(10, 0), scope.DefaultHash, 0},
		{cfgmock.NewService(cfgmock.WithPV(cfgmock.PathValue{wantPath.Bind(scope.Default, 0).String(): 2019})).NewScoped(10, 0), scope.DefaultHash, 2019},
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

func TestIntGetWithoutCfgStructShouldReturnUnexpectedError(t *testing.T) {

	b := cfgmodel.NewInt("web/cors/int")

	gb, h, haveErr := b.Get(cfgmock.NewService(
		cfgmock.WithInt(func(path string) (int, error) {
			return 0, errors.NewFatalf("Unexpected error")
		}),
	).NewScoped(1, 1))
	assert.Empty(t, gb)
	assert.True(t, errors.IsFatal(haveErr), "Error: %s", haveErr)
	assert.Exactly(t, scope.DefaultHash.String(), h.String())
}

func TestIntIgnoreNilDefaultValues(t *testing.T) {

	b := cfgmodel.NewInt("web/cors/int", cfgmodel.WithField(&element.Field{}))
	gb, h, err := b.Get(cfgmock.NewService().NewScoped(0, 0))
	if err != nil {
		t.Fatal(err)
	}
	assert.Exactly(t, int(0), gb)
	assert.Exactly(t, scope.DefaultHash.String(), h.String())
}

func TestIntWrite(t *testing.T) {

	const pathWebCorsInt = "web/cors/int"
	wantPath := cfgpath.MustNewByParts(pathWebCorsInt).Bind(scope.Website, 10)
	b := cfgmodel.NewInt(pathWebCorsInt, cfgmodel.WithFieldFromSectionSlice(configStructure))

	mw := &cfgmock.Write{}
	assert.NoError(t, b.Write(mw, 27182, scope.Website, 10))
	assert.Exactly(t, wantPath.String(), mw.ArgPath)
	assert.Exactly(t, 27182, mw.ArgValue.(int))
}

func TestFloat64GetWithCfgStruct(t *testing.T) {

	const pathWebCorsF64 = "web/cors/float64"
	b := cfgmodel.NewFloat64("web/cors/float64", cfgmodel.WithFieldFromSectionSlice(configStructure))
	assert.Empty(t, b.Options())

	wantPath := cfgpath.MustNewByParts(pathWebCorsF64).Bind(scope.Website, 10)
	tests := []struct {
		sg       config.Scoped
		wantHash scope.Hash
		want     float64
	}{
		{cfgmock.NewService().NewScoped(0, 0), scope.DefaultHash, 2015.1000001}, // because default value in packageConfiguration
		{cfgmock.NewService().NewScoped(0, 1), scope.DefaultHash, 2015.1000001}, // because default value in packageConfiguration
		{cfgmock.NewService().NewScoped(1, 1), scope.DefaultHash, 2015.1000001}, // because default value in packageConfiguration
		{cfgmock.NewService(cfgmock.WithPV(cfgmock.PathValue{wantPath.Bind(scope.Website, 10).String(): 2016.1000001})).NewScoped(10, 0), scope.NewHash(scope.Website, 10), 2016.1000001},
		{cfgmock.NewService(cfgmock.WithPV(cfgmock.PathValue{wantPath.Bind(scope.Website, 10).String(): 2016.1000001})).NewScoped(10, 1), scope.NewHash(scope.Website, 10), 2016.1000001},
		{cfgmock.NewService(cfgmock.WithPV(cfgmock.PathValue{
			wantPath.String():                       2017.1000001,
			wantPath.Bind(scope.Store, 11).String(): 2016.1000021,
		})).NewScoped(10, 11), scope.NewHash(scope.Website, 10), 2017.1000001},
		{cfgmock.NewService(cfgmock.WithPV(cfgmock.PathValue{
			wantPath.String():                         2017.1000001,
			wantPath.Bind(scope.Website, 13).String(): 2018.2000001,
			wantPath.Bind(scope.Store, 11).String():   2016.1000021,
		})).NewScoped(13, 11), scope.NewHash(scope.Website, 13), 2018.2000001},
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

func TestFloat64GetWithoutCfgStruct(t *testing.T) {

	const pathWebCorsF64 = "web/cors/float64"
	b := cfgmodel.NewFloat64(pathWebCorsF64) // no *element.Field has been set. So Default Scope will be enforced
	assert.Empty(t, b.Options())

	wantPath := cfgpath.MustNewByParts(pathWebCorsF64).Bind(scope.Website, 10)
	tests := []struct {
		sg       config.Scoped
		wantHash scope.Hash
		want     float64
	}{
		{cfgmock.NewService().NewScoped(0, 0), scope.DefaultHash, 0},
		{cfgmock.NewService(cfgmock.WithPV(cfgmock.PathValue{wantPath.Bind(scope.Website, 10).String(): 2016.1000001})).NewScoped(10, 0), scope.DefaultHash, 0},
		{cfgmock.NewService(cfgmock.WithPV(cfgmock.PathValue{wantPath.Bind(scope.Default, 0).String(): 2016.1000001})).NewScoped(10, 0), scope.DefaultHash, 2016.1000001},
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

func TestFloat64GetWithoutCfgStructShouldReturnUnexpectedError(t *testing.T) {

	b := cfgmodel.NewFloat64("web/cors/float64")

	gb, h, haveErr := b.Get(cfgmock.NewService(

		cfgmock.WithFloat64(func(path string) (float64, error) {
			return 0, errors.NewFatalf("Unexpected error")
		}),
	).NewScoped(1, 1))
	assert.Empty(t, gb)
	assert.True(t, errors.IsFatal(haveErr), "Error: %s", haveErr)
	assert.Exactly(t, scope.DefaultHash.String(), h.String())
}

func TestFloat64IgnoreNilDefaultValues(t *testing.T) {

	b := cfgmodel.NewFloat64("web/cors/float64", cfgmodel.WithField(&element.Field{}))
	gb, h, err := b.Get(cfgmock.NewService().NewScoped(0, 0))
	if err != nil {
		t.Fatal(err)
	}
	assert.Exactly(t, float64(0), gb)
	assert.Exactly(t, scope.DefaultHash.String(), h.String())
}

func TestFloat64Write(t *testing.T) {

	const pathWebCorsF64 = "web/cors/float64"
	wantPath := cfgpath.MustNewByParts(pathWebCorsF64).Bind(scope.Website, 10)
	b := cfgmodel.NewFloat64("web/cors/float64", cfgmodel.WithFieldFromSectionSlice(configStructure))

	mw := &cfgmock.Write{}
	assert.NoError(t, b.Write(mw, 1.123456789, scope.Website, 10))
	assert.Exactly(t, wantPath.String(), mw.ArgPath)
	assert.Exactly(t, 1.12345678900000, mw.ArgValue.(float64))
}

func TestNewInt_Option_Error(t *testing.T) {
	b := cfgmodel.NewInt(
		"web/cors/int",
		cfgmodel.WithFieldFromSectionSlice(configStructure),
		cfgmodel.WithSourceByString("a", "A", "b", "b"),
	)

	assert.Exactly(t, source.MustNewByString("a", "A", "b", "b"), b.Source)

	err := b.Option(cfgmodel.WithSourceByString(
		"One", "2", "Two",
	))
	assert.True(t, errors.IsNotValid(err), "Error: %s", err)
}

func TestBaseValue_LastError(t *testing.T) {

	wantErr := errors.NewNotImplementedf("Something has not been implemented, e.g. Generics")

	b := cfgmodel.NewBool("a/b/c")
	b.LastError = wantErr
	_, _, haveErr := b.Get(config.Scoped{})
	assert.True(t, errors.IsNotImplemented(haveErr))

	by := cfgmodel.NewByte(`a/b/c`)
	by.LastError = wantErr
	_, _, haveErr = by.Get(config.Scoped{})
	assert.True(t, errors.IsNotImplemented(haveErr))

	str := cfgmodel.NewStr(`a/b/c`)
	str.LastError = wantErr
	_, _, haveErr = str.Get(config.Scoped{})
	assert.True(t, errors.IsNotImplemented(haveErr))

	i := cfgmodel.NewInt(`a/b/c`)
	i.LastError = wantErr
	_, _, haveErr = i.Get(config.Scoped{})
	assert.True(t, errors.IsNotImplemented(haveErr))

	f := cfgmodel.NewFloat64(`a/b/c`)
	f.LastError = wantErr
	_, _, haveErr = f.Get(config.Scoped{})
	assert.True(t, errors.IsNotImplemented(haveErr))

}
