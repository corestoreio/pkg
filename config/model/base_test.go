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

package model

import (
	"testing"

	"github.com/corestoreio/csfw/config"
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
		},
	},
)

func TestBasePathString(t *testing.T) {
	p1 := NewPath("web/cors/exposed_headers")
	assert.Exactly(t, "web/cors/exposed_headers", p1.String())

	wantPath := scope.StrWebsites.FQPathInt64(2, "web/cors/exposed_headers")
	wantWebsiteID := int64(2) // This number 2 is usually stored in core_website/store_website table in column website_id

	mw := new(config.MockWrite)
	assert.NoError(t, p1.Write(mw, 314159, scope.WebsiteID, wantWebsiteID))
	assert.Exactly(t, wantPath, mw.ArgPath)
	assert.Exactly(t, 314159, mw.ArgValue.(int))

	sg := config.NewMockGetter().NewScoped(wantWebsiteID, 0, 0)
	defaultStr, err := p1.lookupString(packageConfiguration, sg)
	assert.NoError(t, err)
	assert.Exactly(t, "Content-Type,X-CoreStore-ID", defaultStr)

	sg = config.NewMockGetter(
		config.WithMockValues(config.MockPV{
			wantPath: "X-CoreStore-TOKEN",
		}),
	).NewScoped(wantWebsiteID, 0, 0)

	customStr, err := p1.lookupString(packageConfiguration, sg)
	assert.NoError(t, err)
	assert.Exactly(t, "X-CoreStore-TOKEN", customStr)

}

func TestBasePathInScope(t *testing.T) {

	tests := []struct {
		sg      config.ScopedGetter
		p       scope.Perm
		wantErr error
	}{
		{
			config.NewMockGetter().NewScoped(0, 0, 0),
			scope.NewPerm(scope.DefaultID, scope.WebsiteID),
			nil,
		},
	}
	for _, test := range tests {
		p1 := NewPath("a/b/c")
		haveErr := p1.InScope(&config.Field{
			Scope: test.p,
		}, test.sg)

		if test.wantErr != nil {
			assert.EqualError(t, haveErr, test.wantErr.Error())
		} else {
			assert.NoError(t, haveErr)
		}
	}
}
