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
	"github.com/corestoreio/csfw/config/source"
	"github.com/corestoreio/csfw/store/scope"
	"github.com/stretchr/testify/assert"
)

// configStructure might be a duplicate of base_test but note that the
// test package names are different.
var configStructure = element.MustNewConfiguration(
	&element.Section{
		ID: "web",
		Groups: element.NewGroupSlice(
			&element.Group{
				ID:        "cors",
				Label:     `CORS Cross Origin Resource Sharing`,
				SortOrder: 150,
				Scope:     scope.NewPerm(scope.DefaultID),
				Fields: element.NewFieldSlice(
					&element.Field{
						// Path: `web/cors/exposed_headers`,
						ID:        "exposed_headers",
						Label:     `Exposed Headers`,
						Comment:   text.Long(`Indicates which headers are safe to expose to the API of a CORS API specification. Separate via line break`),
						Type:      element.TypeTextarea,
						SortOrder: 10,
						Visible:   element.VisibleYes,
						Scope:     scope.NewPerm(scope.DefaultID, scope.WebsiteID),
						Default:   "Content-Type,X-CoreStore-ID",
					},
					&element.Field{
						// Path: `web/cors/allowed_origins`,
						ID:        "allowed_origins",
						Label:     `Allowed Origins`,
						Comment:   text.Long(`Is a list of origins a cross-domain request can be executed from.`),
						Type:      element.TypeTextarea,
						SortOrder: 20,
						Visible:   element.VisibleYes,
						Scope:     scope.NewPerm(scope.DefaultID, scope.WebsiteID),
						Default:   "corestore.io,cs.io",
					},
					&element.Field{
						// Path: `web/cors/allow_credentials`,
						ID:        "allow_credentials",
						Label:     `Allowed Credentials`,
						Type:      element.TypeSelect,
						SortOrder: 30,
						Visible:   element.VisibleYes,
						Scope:     scope.NewPerm(scope.DefaultID, scope.WebsiteID),
						Default:   "true",
					},
					&element.Field{
						// Path: `web/cors/int`,
						ID:        "int",
						Type:      element.TypeText,
						SortOrder: 30,
						Visible:   element.VisibleYes,
						Scope:     scope.NewPerm(scope.DefaultID, scope.WebsiteID),
						Default:   2015,
					},
					&element.Field{
						// Path: `web/cors/int_slice`,
						ID:        "int_slice",
						Type:      element.TypeSelect,
						SortOrder: 30,
						Visible:   element.VisibleYes,
						Scope:     scope.PermAll,
						Default:   "2014,2015,2016",
					},
					&element.Field{
						// Path: `web/cors/float64`,
						ID:        "float64",
						Type:      element.TypeSelect,
						SortOrder: 30,
						Visible:   element.VisibleYes,
						Scope:     scope.NewPerm(scope.DefaultID, scope.WebsiteID),
						Default:   2015.1000001,
					},
				),
			},

			&element.Group{
				ID:        "unsecure",
				Label:     `Base URLs`,
				Comment:   text.Long(`Any of the fields allow fully qualified URLs that end with '/' (slash) e.g. http://example.com/magento/`),
				SortOrder: 10,
				Scope:     scope.PermAll,
				Fields: element.NewFieldSlice(
					&element.Field{
						// Path: `web/unsecure/base_url`,
						ID:        "base_url",
						Label:     `Base URL`,
						Comment:   text.Long(`Specify URL or {{base_url}} placeholder.`),
						Type:      element.TypeText,
						SortOrder: 10,
						Visible:   element.VisibleYes,
						Scope:     scope.PermAll,
						Default:   "{{base_url}}",
						//BackendModel: nil, // Magento\Config\Model\Config\Backend\Baseurl
					},

					&element.Field{
						// Path: `web/unsecure/base_link_url`,
						ID:        "base_link_url",
						Label:     `Base Link URL`,
						Comment:   text.Long(`May start with {{unsecure_base_url}} placeholder.`),
						Type:      element.TypeText,
						SortOrder: 20,
						Visible:   element.VisibleYes,
						Scope:     scope.PermAll,
						Default:   "{{unsecure_base_url}}",
						//BackendModel: nil, // Magento\Config\Model\Config\Backend\Baseurl
					},

					&element.Field{
						// Path: `web/unsecure/base_static_url`,
						ID:        "base_static_url",
						Label:     `Base URL for Static View Files`,
						Comment:   text.Long(`May be empty or start with {{unsecure_base_url}} placeholder.`),
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
	wantPath := scope.StrWebsites.FQPathInt64(3, "web/cors/allow_credentials")
	b := model.NewBool("web/cors/allow_credentials", model.WithConfigStructure(configStructure), model.WithSource(source.YesNo))

	assert.Exactly(t, source.YesNo, b.Options())
	// because default value in packageConfiguration is "true"
	assert.True(t, b.Get(config.NewMockGetter().NewScoped(0, 0, 0)))

	assert.False(t, b.Get(config.NewMockGetter(
		config.WithMockValues(config.MockPV{
			wantPath: 0,
		}),
	).NewScoped(0, 0, 3)))

	mw := &config.MockWrite{}
	assert.EqualError(t, b.Write(mw, true, scope.StoreID, 3), "Scope permission insufficient: Have 'Store'; Want 'Default,Website'")
	assert.NoError(t, b.Write(mw, true, scope.WebsiteID, 3))
	assert.Exactly(t, wantPath, mw.ArgPath)
	assert.Exactly(t, "true", mw.ArgValue.(string))
}

func TestStr(t *testing.T) {
	t.Parallel()
	wantPath := scope.StrDefault.FQPathInt64(0, "web/cors/exposed_headers")
	b := model.NewStr("web/cors/exposed_headers", model.WithConfigStructure(configStructure))

	assert.Empty(t, b.Options())

	assert.Exactly(t, "Content-Type,X-CoreStore-ID", b.Get(config.NewMockGetter().NewScoped(0, 0, 0)))

	assert.Exactly(t, "X-Gopher", b.Get(config.NewMockGetter(
		config.WithMockValues(config.MockPV{
			wantPath: "X-Gopher",
		}),
	).NewScoped(0, 0, 0)))

	mw := &config.MockWrite{}
	assert.NoError(t, b.Write(mw, "dude", scope.DefaultID, 0))
	assert.Exactly(t, wantPath, mw.ArgPath)
	assert.Exactly(t, "dude", mw.ArgValue.(string))
}

func TestInt(t *testing.T) {
	t.Parallel()
	wantPath := scope.StrWebsites.FQPathInt64(10, "web/cors/int")
	b := model.NewInt("web/cors/int", model.WithConfigStructure(configStructure))

	assert.Empty(t, b.Options())

	assert.Exactly(t, 2015, b.Get(config.NewMockGetter().NewScoped(0, 0, 0)))

	assert.Exactly(t, 2016, b.Get(config.NewMockGetter(
		config.WithMockValues(config.MockPV{
			wantPath: 2016,
		}),
	).NewScoped(10, 0, 0)))

	mw := &config.MockWrite{}
	assert.NoError(t, b.Write(mw, 1, scope.WebsiteID, 10))
	assert.Exactly(t, wantPath, mw.ArgPath)
	assert.Exactly(t, "1", mw.ArgValue.(string))
}

func TestFloat64(t *testing.T) {
	t.Parallel()
	wantPath := scope.StrWebsites.FQPathInt64(10, "web/cors/float64")
	b := model.NewFloat64("web/cors/float64", model.WithConfigStructure(configStructure))

	assert.Empty(t, b.Options())

	assert.Exactly(t, 2015.1000001, b.Get(config.NewMockGetter().NewScoped(0, 0, 0)))

	assert.Exactly(t, 2016.1000001, b.Get(config.NewMockGetter(
		config.WithMockValues(config.MockPV{
			wantPath: 2016.1000001,
		}),
	).NewScoped(10, 0, 0)))

	mw := &config.MockWrite{}
	assert.NoError(t, b.Write(mw, 1.123456789, scope.WebsiteID, 10))
	assert.Exactly(t, wantPath, mw.ArgPath)
	assert.Exactly(t, "1.12345678900000", mw.ArgValue.(string))
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
