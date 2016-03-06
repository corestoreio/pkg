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
	"errors"
	"testing"

	"github.com/corestoreio/csfw/config"
	"github.com/corestoreio/csfw/config/element"
	"github.com/corestoreio/csfw/config/mock"
	"github.com/corestoreio/csfw/config/path"
	"github.com/corestoreio/csfw/config/source"
	"github.com/corestoreio/csfw/storage/text"
	"github.com/corestoreio/csfw/store/scope"
	"github.com/stretchr/testify/assert"
)

var _ source.Optioner = (*baseValue)(nil)

// configStructure might be a duplicate of primitives_test but note that the
// test package names are different.
var configStructure = element.MustNewConfiguration(
	&element.Section{
		ID: path.NewRoute("web"),
		Groups: element.NewGroupSlice(
			&element.Group{
				ID:        path.NewRoute("cors"),
				Label:     text.Chars(`CORS Cross Origin Resource Sharing`),
				SortOrder: 150,
				Scope:     scope.PermDefault,
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
						// Path: `web/cors/float64`,
						ID:        path.NewRoute("float64"),
						Type:      element.TypeSelect,
						SortOrder: 30,
						Visible:   element.VisibleYes,
						Scopes:    scope.PermWebsite,
						Default:   2015.1000001,
					},
				),
			},
		),
	},
)

func TestBaseValueString(t *testing.T) {
	t.Parallel()
	const pathWebCorsHeaders = "web/cors/exposed_headers"
	p1 := NewStr(pathWebCorsHeaders, WithFieldFromSectionSlice(configStructure))
	assert.Exactly(t, pathWebCorsHeaders, p1.String())

	wantWebsiteID := int64(2) // This number 2 is usually stored in core_website/store_website table in column website_id
	wantPath := path.MustNewByParts(pathWebCorsHeaders).Bind(scope.WebsiteID, wantWebsiteID)

	mw := new(mock.Write)
	assert.NoError(t, p1.Write(mw, "314159", scope.WebsiteID, wantWebsiteID))
	assert.Exactly(t, wantPath.String(), mw.ArgPath)
	assert.Exactly(t, "314159", mw.ArgValue.(string))

	sg := mock.NewService().NewScoped(wantWebsiteID, 0)
	defaultStr, err := p1.Get(sg)
	assert.NoError(t, err)
	assert.Exactly(t, "Content-Type,X-CoreStore-ID", defaultStr)

	sg = mock.NewService(
		mock.WithPV(mock.PathValue{
			wantPath.String(): "X-CoreStore-TOKEN",
		}),
	).NewScoped(wantWebsiteID, 0)

	customStr, err := p1.Get(sg)
	assert.NoError(t, err)
	assert.Exactly(t, "X-CoreStore-TOKEN", customStr)

	// now change a default value in the packageConfiguration and see it reflects to p1
	f, err := configStructure.FindFieldByID(wantPath.Route)
	assert.NoError(t, err)
	f.Default = "Content-Size,Y-CoreStore-ID"

	ws, err := p1.Get(mock.NewService().NewScoped(wantWebsiteID, 0))
	assert.NoError(t, err)
	assert.Exactly(t, "Content-Size,Y-CoreStore-ID", ws)
}

func TestBaseValueInScope(t *testing.T) {
	t.Parallel()
	tests := []struct {
		sg      config.ScopedGetter
		p       scope.Perm
		wantErr error
	}{
		{
			mock.NewService().NewScoped(0, 0),
			scope.PermWebsite,
			nil,
		},
		{
			mock.NewService().NewScoped(0, 4),
			scope.PermStore,
			nil,
		},
		{
			mock.NewService().NewScoped(4, 0),
			scope.PermStore,
			nil,
		},
		{
			mock.NewService().NewScoped(0, 4),
			scope.PermWebsite,
			errors.New("Scope permission insufficient: Have 'Store'; Want 'Default,Website'"),
		},
		{
			mock.NewService().NewScoped(4, 0),
			scope.PermDefault,
			errors.New("Scope permission insufficient: Have 'Website'; Want 'Default'"),
		},
	}
	for i, test := range tests {
		p1 := NewValue("a/b/c", WithField(&element.Field{
			Scopes: test.p,
		}))
		haveErr := p1.InScope(test.sg)

		if test.wantErr != nil {
			assert.EqualError(t, haveErr, test.wantErr.Error(), "Index %d", i)
		} else {
			assert.NoError(t, haveErr, "Index %d", i)
		}
	}
}

func TestBaseValueFQ(t *testing.T) {
	t.Parallel()
	const pth = "aa/bb/cc"
	p := NewValue(pth)
	fq, err := p.FQ(scope.StoreID, 4)
	assert.NoError(t, err)
	assert.Exactly(t, path.MustNewByParts(pth).Bind(scope.StoreID, 4).String(), fq)
}

func TestBaseValueToPath(t *testing.T) {
	t.Parallel()
	tests := []struct {
		route   path.Route
		s       scope.Scope
		sid     int64
		wantErr error
	}{
		{path.NewRoute("aa/bb/cc"), scope.StoreID, 23, nil},
		{path.NewRoute("a/bb/cc"), scope.StoreID, 23, path.ErrIncorrectPath},
	}
	for i, test := range tests {
		bv := NewValue(test.route.String())
		havePath, haveErr := bv.ToPath(test.s, test.sid)
		if test.wantErr != nil {
			assert.EqualError(t, haveErr, test.wantErr.Error(), "Index %d", i)
			continue
		}
		assert.NoError(t, test.wantErr, "Index %d", i)
		wantPath := path.MustNew(test.route).Bind(test.s, test.sid)
		assert.Exactly(t, wantPath, havePath, "Index %d", i)
	}
}

func TestBaseValueRoute(t *testing.T) {
	t.Parallel()
	org := NewValue("aa/bb/cc")
	clone := org.Route()

	if &(org.route) == &clone { // comparing pointer addresses
		// is there a better way to test of the slice headers points to a different location?
		// because clone should be a clone ;-)
		t.Error("Should not be equal")
	}
}
