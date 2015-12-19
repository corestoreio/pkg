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
	"github.com/stretchr/testify/assert"
)

var packageConfiguration = config.MustNewConfiguration(
	&config.Section{
		ID: "web",
		Groups: config.GroupSlice{
			&config.Group{
				ID:        "cors",
				Label:     `CORS Cross Origin Resource Sharing`,
				Comment:   ``,
				SortOrder: 150,
				Scope:     scope.NewPerm(scope.DefaultID),
				Fields: config.FieldSlice{
					&config.Field{
						// Path: `web/cors/exposed_headers`,
						ID:        "exposed_headers",
						Label:     `Exposed Headers`,
						Comment:   `Indicates which headers are safe to expose to the API of a CORS API specification. Separate via line break`,
						Type:      config.TypeTextarea,
						SortOrder: 10,
						Visible:   config.VisibleYes,
						Scope:     scope.NewPerm(scope.DefaultID, scope.WebsiteID),
						Default:   "Content-Type,X-CoreStore-ID",
					},
					&config.Field{
						// Path: `web/cors/allowed_origins`,
						ID:        "allowed_origins",
						Label:     `Allowed Origins`,
						Comment:   `Is a list of origins a cross-domain request can be executed from.`,
						Type:      config.TypeTextarea,
						SortOrder: 20,
						Visible:   config.VisibleYes,
						Scope:     scope.NewPerm(scope.DefaultID, scope.WebsiteID),
						Default:   "corestore.io,cs.io",
					},
					&config.Field{
						// Path: `web/cors/allow_credentials`,
						ID:        "allow_credentials",
						Label:     `Allowed Credentials`,
						Comment:   ``,
						Type:      config.TypeSelect,
						SortOrder: 30,
						Visible:   config.VisibleYes,
						Scope:     scope.NewPerm(scope.DefaultID, scope.WebsiteID),
						Default:   "true",
					},
					&config.Field{
						// Path: `web/cors/int`,
						ID:        "int",
						Type:      config.TypeText,
						SortOrder: 30,
						Visible:   config.VisibleYes,
						Scope:     scope.NewPerm(scope.DefaultID, scope.WebsiteID),
						Default:   2015,
					},
					&config.Field{
						// Path: `web/cors/int_slice`,
						ID:        "int_slice",
						Type:      config.TypeSelect,
						SortOrder: 30,
						Visible:   config.VisibleYes,
						Scope:     scope.NewPerm(scope.DefaultID, scope.WebsiteID),
						Default:   "2014,2015,2016",
					},
					&config.Field{
						// Path: `web/cors/float64`,
						ID:        "float64",
						Type:      config.TypeSelect,
						SortOrder: 30,
						Visible:   config.VisibleYes,
						Scope:     scope.NewPerm(scope.DefaultID, scope.WebsiteID),
						Default:   2015.1000001,
					},
				},
			},

			&config.Group{
				ID:        "unsecure",
				Label:     `Base URLs`,
				Comment:   `Any of the fields allow fully qualified URLs that end with '/' (slash) e.g. http://example.com/magento/`,
				SortOrder: 10,
				Scope:     scope.PermAll,
				Fields: config.FieldSlice{
					&config.Field{
						// Path: `web/unsecure/base_url`,
						ID:        "base_url",
						Label:     `Base URL`,
						Comment:   `Specify URL or {{base_url}} placeholder.`,
						Type:      config.TypeText,
						SortOrder: 10,
						Visible:   config.VisibleYes,
						Scope:     scope.PermAll,
						Default:   "{{base_url}}",
						//BackendModel: nil, // Magento\Config\Model\Config\Backend\Baseurl
					},

					&config.Field{
						// Path: `web/unsecure/base_link_url`,
						ID:        "base_link_url",
						Label:     `Base Link URL`,
						Comment:   `May start with {{unsecure_base_url}} placeholder.`,
						Type:      config.TypeText,
						SortOrder: 20,
						Visible:   config.VisibleYes,
						Scope:     scope.PermAll,
						Default:   "{{unsecure_base_url}}",
						//BackendModel: nil, // Magento\Config\Model\Config\Backend\Baseurl
					},

					&config.Field{
						// Path: `web/unsecure/base_static_url`,
						ID:        "base_static_url",
						Label:     `Base URL for Static View Files`,
						Comment:   `May be empty or start with {{unsecure_base_url}} placeholder.`,
						Type:      config.TypeText,
						SortOrder: 25,
						Visible:   config.VisibleYes,
						Scope:     scope.PermAll,
						Default:   nil,
						//BackendModel: nil, // Magento\Config\Model\Config\Backend\Baseurl
					},
				},
			},
		},
	},
)

func TestPath(t *testing.T) {
	p1 := model.NewPath("web/cors/exposed_headers")
	assert.Exactly(t, "web/cors/exposed_headers", p1.String())

	wantPath := scope.StrWebsites.FQPathInt64(2, "web/cors/exposed_headers")
	wantWebsiteID := int64(2) // This number 2 is usually stored in core_website/store_website table in column website_id

	mw := new(config.MockWrite)
	assert.NoError(t, p1.Write(mw, 314159, scope.WebsiteID, wantWebsiteID))
	assert.Exactly(t, wantPath, mw.ArgPath)
	assert.Exactly(t, 314159, mw.ArgValue.(int))

	sg := config.NewMockGetter().NewScoped(wantWebsiteID, 0, 0)
	defaultStr := p1.LookupString(packageConfiguration, sg)
	assert.Exactly(t, "Content-Type,X-CoreStore-ID", defaultStr)

	sg = config.NewMockGetter(
		config.WithMockValues(config.MockPV{
			wantPath: "X-CoreStore-TOKEN",
		}),
	).NewScoped(wantWebsiteID, 0, 0)

	customStr := p1.LookupString(packageConfiguration, sg)
	assert.Exactly(t, "X-CoreStore-TOKEN", customStr)

	assert.True(t, p1.InScope(&config.Field{
		Scope: scope.NewPerm(scope.DefaultID, scope.WebsiteID),
	}, sg))

	assert.False(t, p1.InScope(&config.Field{
		Scope: scope.NewPerm(scope.StoreID),
	}, sg))
}
