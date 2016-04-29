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

package cfgmodel

import (
	"testing"

	"github.com/corestoreio/csfw/config"
	"github.com/corestoreio/csfw/config/cfgmock"
	"github.com/corestoreio/csfw/config/cfgpath"
	"github.com/corestoreio/csfw/config/element"
	"github.com/corestoreio/csfw/config/source"
	"github.com/corestoreio/csfw/storage/text"
	"github.com/corestoreio/csfw/store/scope"
	"github.com/corestoreio/csfw/util/errors"
	"github.com/stretchr/testify/assert"
)

var _ source.Optioner = (*baseValue)(nil)

// configStructure might be a duplicate of primitives_test but note that the
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
						// Path: `web/cors/float64`,
						ID:        cfgpath.NewRoute("float64"),
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
	wantPath := cfgpath.MustNewByParts(pathWebCorsHeaders).Bind(scope.Website, wantWebsiteID)

	mw := new(cfgmock.Write)
	assert.NoError(t, p1.Write(mw, "314159", scope.Website, wantWebsiteID))
	assert.Exactly(t, wantPath.String(), mw.ArgPath)
	assert.Exactly(t, "314159", mw.ArgValue.(string))

	sg := cfgmock.NewService().NewScoped(wantWebsiteID, 0)
	defaultStr, err := p1.Get(sg)
	assert.NoError(t, err)
	assert.Exactly(t, "Content-Type,X-CoreStore-ID", defaultStr)

	sg = cfgmock.NewService(
		cfgmock.WithPV(cfgmock.PathValue{
			wantPath.String(): "X-CoreStore-TOKEN",
		}),
	).NewScoped(wantWebsiteID, 0)

	customStr, err := p1.Get(sg)
	assert.NoError(t, err)
	assert.Exactly(t, "X-CoreStore-TOKEN", customStr)

	// now change a default value in the packageConfiguration and see it reflects to p1.
	// but this is not the way to go. You can directly change the field in p1
	// with p1.Field.Default
	if err := configStructure.UpdateField(wantPath.Route, element.Field{
		Default: "Content-Size,Y-CoreStore-ID",
	}); err != nil {
		t.Fatal(err)
	}

	// update p1 to apply the change field data
	p1.Option(WithFieldFromSectionSlice(configStructure))

	ws, err := p1.Get(cfgmock.NewService().NewScoped(wantWebsiteID, 0))
	assert.NoError(t, err)
	assert.Exactly(t, "Content-Size,Y-CoreStore-ID", ws)
}

func TestBaseValueInScope(t *testing.T) {
	t.Parallel()
	tests := []struct {
		sg         config.ScopedGetter
		p          scope.Perm
		wantErrBhf errors.BehaviourFunc
	}{
		{
			cfgmock.NewService().NewScoped(0, 0),
			scope.PermWebsite,
			nil,
		},
		{
			cfgmock.NewService().NewScoped(0, 4),
			scope.PermStore,
			nil,
		},
		{
			cfgmock.NewService().NewScoped(4, 0),
			scope.PermStore,
			nil,
		},
		{
			cfgmock.NewService().NewScoped(0, 4),
			scope.PermWebsite,
			errors.IsUnauthorized,
		},
		{
			cfgmock.NewService().NewScoped(4, 0),
			scope.PermDefault,
			errors.IsUnauthorized,
		},
	}
	for i, test := range tests {
		p1 := NewValue("a/b/c", WithField(&element.Field{
			ID:     cfgpath.NewRoute(`c`),
			Scopes: test.p,
		}))
		haveErr := p1.InScope(test.sg)

		if test.wantErrBhf != nil {
			assert.True(t, test.wantErrBhf(haveErr), "Index %d => %s", i, haveErr)
		} else {
			assert.NoError(t, haveErr, "Index %d", i)
		}
	}
}

func TestBaseValueFQ(t *testing.T) {
	t.Parallel()
	const pth = "aa/bb/cc"
	p := NewValue(pth)
	fq, err := p.FQ(scope.Store, 4)
	assert.NoError(t, err)
	assert.Exactly(t, cfgpath.MustNewByParts(pth).Bind(scope.Store, 4).String(), fq)
}

func TestBaseValueMustFQPanic(t *testing.T) {
	t.Parallel()
	defer func() {
		if r := recover(); r != nil {
			err := r.(error)
			assert.True(t, errors.IsNotValid(err), "Error: %s", err)
		} else {
			t.Fatal("Expecting a panic")
		}
	}()
	const pth = "a/b/c"
	p := NewValue(pth)
	fq := p.MustFQ(scope.Store, 4)
	assert.Empty(t, fq)
}

func TestBaseValueToPath(t *testing.T) {
	t.Parallel()
	tests := []struct {
		route      cfgpath.Route
		s          scope.Scope
		sid        int64
		wantErrBhf errors.BehaviourFunc
	}{
		{cfgpath.NewRoute("aa/bb/cc"), scope.Store, 23, nil},
		{cfgpath.NewRoute("a/bb/cc"), scope.Store, 23, errors.IsNotValid},
	}
	for i, test := range tests {
		bv := NewValue(test.route.String())
		havePath, haveErr := bv.ToPath(test.s, test.sid)
		if test.wantErrBhf != nil {
			assert.True(t, test.wantErrBhf(haveErr), "Index %d => %s", i, haveErr)
			continue
		}
		wantPath := cfgpath.MustNew(test.route).Bind(test.s, test.sid)
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
