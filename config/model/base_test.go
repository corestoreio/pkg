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

package model

import (
	"testing"

	"errors"

	"github.com/corestoreio/csfw/config"
	"github.com/corestoreio/csfw/config/element"
	"github.com/corestoreio/csfw/store/scope"
	"github.com/stretchr/testify/assert"
)

// configStructure might be a duplicate of primitives_test but note that the
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
		),
	},
)

func TestBasePathString(t *testing.T) {
	t.Parallel()
	const path = "web/cors/exposed_headers"
	p1 := NewPath(path, WithConfigStructure(configStructure))
	assert.Exactly(t, path, p1.String())

	wantPath := scope.StrWebsites.FQPathInt64(2, "web/cors/exposed_headers")
	wantWebsiteID := int64(2) // This number 2 is usually stored in core_website/store_website table in column website_id

	mw := new(config.MockWrite)
	assert.NoError(t, p1.Write(mw, 314159, scope.WebsiteID, wantWebsiteID))
	assert.Exactly(t, wantPath, mw.ArgPath)
	assert.Exactly(t, 314159, mw.ArgValue.(int))

	sg := config.NewMockGetter().NewScoped(wantWebsiteID, 0, 0)
	defaultStr, err := p1.lookupString(sg)
	assert.NoError(t, err)
	assert.Exactly(t, "Content-Type,X-CoreStore-ID", defaultStr)

	sg = config.NewMockGetter(
		config.WithMockValues(config.MockPV{
			wantPath: "X-CoreStore-TOKEN",
		}),
	).NewScoped(wantWebsiteID, 0, 0)

	customStr, err := p1.lookupString(sg)
	assert.NoError(t, err)
	assert.Exactly(t, "X-CoreStore-TOKEN", customStr)

	// now change a default value in the packageConfiguration and see it reflects to p1
	f, err := configStructure.FindFieldByPath(path)
	assert.NoError(t, err)
	f.Default = "Content-Size,Y-CoreStore-ID"

	ws, err := p1.lookupString(config.NewMockGetter().NewScoped(wantWebsiteID, 0, 0))
	assert.NoError(t, err)
	assert.Exactly(t, "Content-Size,Y-CoreStore-ID", ws)
}

func TestBasePathInScope(t *testing.T) {
	t.Parallel()
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
		{
			config.NewMockGetter().NewScoped(0, 0, 4),
			scope.NewPerm(scope.StoreID),
			nil,
		},
		{
			config.NewMockGetter().NewScoped(0, 4, 0),
			scope.NewPerm(scope.StoreID),
			errors.New("Scope permission insufficient: Have 'Group'; Want 'Store'"),
		},
	}
	for _, test := range tests {
		p1 := NewPath("a/b/c", WithField(&element.Field{
			Scope: test.p,
		}))
		haveErr := p1.InScope(test.sg)

		if test.wantErr != nil {
			assert.EqualError(t, haveErr, test.wantErr.Error())
		} else {
			assert.NoError(t, haveErr)
		}
	}
}

func TestFQPathInt64(t *testing.T) {
	t.Parallel()
	p := NewPath("a/b/c")
	assert.Exactly(t, scope.StrStores.FQPathInt64(4, "a/b/c"), p.FQ(scope.StrStores, 4))
}
