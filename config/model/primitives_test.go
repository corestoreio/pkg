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
	"github.com/corestoreio/csfw/config/element"
	"github.com/corestoreio/csfw/config/model"
	"github.com/corestoreio/csfw/config/path"
	"github.com/corestoreio/csfw/config/source"
	"github.com/corestoreio/csfw/storage/text"
	"github.com/corestoreio/csfw/store/scope"
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
				Scope:     scope.NewPerm(scope.DefaultID),
				Fields: element.NewFieldSlice(
					&element.Field{
						// Path: `web/cors/exposed_headers`,
						ID:        path.NewRoute("exposed_headers"),
						Label:     text.Chars(`Exposed Headers`),
						Comment:   text.Chars(`Indicates which headers are safe to expose to the API of a CORS API specification. Separate via line break`),
						Type:      element.TypeTextarea,
						SortOrder: 10,
						Visible:   element.VisibleYes,
						Scope:     scope.NewPerm(scope.DefaultID, scope.WebsiteID),
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
						Scope:     scope.NewPerm(scope.DefaultID, scope.WebsiteID),
						Default:   "corestore.io,cs.io",
					},
					&element.Field{
						// Path: `web/cors/allow_credentials`,
						ID:        path.NewRoute("allow_credentials"),
						Label:     text.Chars(`Allowed Credentials`),
						Type:      element.TypeSelect,
						SortOrder: 30,
						Visible:   element.VisibleYes,
						Scope:     scope.NewPerm(scope.DefaultID, scope.WebsiteID),
						Default:   "true",
					},
					&element.Field{
						// Path: `web/cors/int`,
						ID:        path.NewRoute("int"),
						Type:      element.TypeText,
						SortOrder: 30,
						Visible:   element.VisibleYes,
						Scope:     scope.NewPerm(scope.DefaultID, scope.WebsiteID),
						Default:   2015,
					},
					&element.Field{
						// Path: `web/cors/int_slice`,
						ID:        path.NewRoute("int_slice"),
						Type:      element.TypeSelect,
						SortOrder: 30,
						Visible:   element.VisibleYes,
						Scope:     scope.PermAll,
						Default:   "2014,2015,2016",
					},
					&element.Field{
						// Path: `web/cors/float64`,
						ID:        path.NewRoute("float64"),
						Type:      element.TypeSelect,
						SortOrder: 30,
						Visible:   element.VisibleYes,
						Scope:     scope.NewPerm(scope.DefaultID, scope.WebsiteID),
						Default:   2015.1000001,
					},
				),
			},

			&element.Group{
				ID:        path.NewRoute("unsecure"),
				Label:     text.Chars(`Base URLs`),
				Comment:   text.Chars(`Any of the fields allow fully qualified URLs that end with '/' (slash) e.g. http://example.com/magento/`),
				SortOrder: 10,
				Scope:     scope.PermAll,
				Fields: element.NewFieldSlice(
					&element.Field{
						// Path: `web/unsecure/base_url`,
						ID:        path.NewRoute("base_url"),
						Label:     text.Chars(`Base URL`),
						Comment:   text.Chars(`Specify URL or {{base_url}} placeholder.`),
						Type:      element.TypeText,
						SortOrder: 10,
						Visible:   element.VisibleYes,
						Scope:     scope.PermAll,
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
						Scope:     scope.PermAll,
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
						Scope:     scope.PermAll,
						Default:   nil,
						//BackendModel: nil, // Magento\Config\Model\Config\Backend\Baseurl
					},
				),
			},
		),
	},
)

func TestBool(t *testing.T) {
	t.Parallel()
	const pathWebCorsCred = "web/cors/allow_credentials"
	wantPath := path.MustNewByParts(pathWebCorsCred).Bind(scope.WebsiteID, 3)
	b := model.NewBool(pathWebCorsCred, model.WithConfigStructure(configStructure), model.WithSource(source.YesNo))

	assert.Exactly(t, source.YesNo, b.Options())
	// because default value in packageConfiguration is "true"
	assert.True(t, b.Get(config.NewMockGetter().NewScoped(0, 0, 0)))

	assert.False(t, b.Get(config.NewMockGetter(
		config.WithMockValues(config.MockPV{
			wantPath.String(): 0,
		}),
	).NewScoped(0, 0, 3)))

	mw := &config.MockWrite{}
	assert.EqualError(t, b.Write(mw, true, scope.StoreID, 3), "Scope permission insufficient: Have 'Store'; Want 'Default,Website'")
	assert.NoError(t, b.Write(mw, true, scope.WebsiteID, 3))
	assert.Exactly(t, wantPath.String(), mw.ArgPath)
	assert.Exactly(t, true, mw.ArgValue.(bool))
}

func TestStr(t *testing.T) {
	t.Parallel()
	const pathWebCorsHeaders = "web/cors/exposed_headers"
	wantPath := path.MustNewByParts(pathWebCorsHeaders)
	b := model.NewStr(pathWebCorsHeaders, model.WithConfigStructure(configStructure))

	assert.Empty(t, b.Options())

	assert.Exactly(t, "Content-Type,X-CoreStore-ID", b.Get(config.NewMockGetter().NewScoped(0, 0, 0)))

	assert.Exactly(t, "X-Gopher", b.Get(config.NewMockGetter(
		config.WithMockValues(config.MockPV{
			wantPath.String(): "X-Gopher",
		}),
	).NewScoped(0, 0, 0)))

	mw := &config.MockWrite{}
	assert.NoError(t, b.Write(mw, "dude", scope.DefaultID, 0))
	assert.Exactly(t, wantPath.String(), mw.ArgPath)
	assert.Exactly(t, "dude", mw.ArgValue.(string))
}

func TestInt(t *testing.T) {
	t.Parallel()
	const pathWebCorsInt = "web/cors/int"
	wantPath := path.MustNewByParts(pathWebCorsInt).Bind(scope.WebsiteID, 10)
	b := model.NewInt(pathWebCorsInt, model.WithConfigStructure(configStructure))

	assert.Empty(t, b.Options())

	assert.Exactly(t, 2015, b.Get(config.NewMockGetter().NewScoped(0, 0, 0)))

	assert.Exactly(t, 2016, b.Get(config.NewMockGetter(
		config.WithMockValues(config.MockPV{
			wantPath.String(): 2016,
		}),
	).NewScoped(10, 0, 0)))

	mw := &config.MockWrite{}
	assert.NoError(t, b.Write(mw, 27182, scope.WebsiteID, 10))
	assert.Exactly(t, wantPath.String(), mw.ArgPath)
	assert.Exactly(t, 27182, mw.ArgValue.(int))
}

func TestFloat64(t *testing.T) {
	t.Parallel()
	const pathWebCorsF64 = "web/cors/float64"
	wantPath := path.MustNewByParts(pathWebCorsF64).Bind(scope.WebsiteID, 10)
	b := model.NewFloat64("web/cors/float64", model.WithConfigStructure(configStructure))

	assert.Empty(t, b.Options())

	assert.Exactly(t, 2015.1000001, b.Get(config.NewMockGetter().NewScoped(0, 0, 0)))

	assert.Exactly(t, 2016.1000001, b.Get(config.NewMockGetter(
		config.WithMockValues(config.MockPV{
			wantPath.String(): 2016.1000001,
		}),
	).NewScoped(10, 0, 0)))

	mw := &config.MockWrite{}
	assert.NoError(t, b.Write(mw, 1.123456789, scope.WebsiteID, 10))
	assert.Exactly(t, wantPath.String(), mw.ArgPath)
	assert.Exactly(t, 1.12345678900000, mw.ArgValue.(float64))
}

func TestRecursiveOption(t *testing.T) {
	t.Parallel()
	b := model.NewInt(
		"web/cors/int",
		model.WithConfigStructure(configStructure),
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
