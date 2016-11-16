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
	"github.com/corestoreio/csfw/config/cfgsource"
	"github.com/corestoreio/csfw/config/element"
	"github.com/corestoreio/csfw/storage/text"
	"github.com/corestoreio/csfw/store/scope"
	"github.com/corestoreio/csfw/util/errors"
	"github.com/stretchr/testify/assert"
)

var _ cfgsource.Optioner = (*baseValue)(nil)

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

	const pathWebCorsHeaders = "web/cors/exposed_headers" // also perm website
	p1 := NewStr(pathWebCorsHeaders, WithFieldFromSectionSlice(configStructure))
	assert.Exactly(t, pathWebCorsHeaders, p1.String())

	wantWebsiteID := int64(2) // This number 2 is usually stored in core_website/store_website table in column website_id
	wantPath := cfgpath.MustNewByParts(pathWebCorsHeaders).BindWebsite(wantWebsiteID)

	mw := new(cfgmock.Write)
	err := p1.Write(mw, "314159", scope.Website.Pack(wantWebsiteID))
	assert.NoError(t, err, "%+v", err)
	assert.Exactly(t, wantPath.String(), mw.ArgPath)

	assert.Exactly(t, "314159", mw.ArgValue.(string))

	sg := cfgmock.NewService().NewScoped(wantWebsiteID, 0)
	defaultStr, err := p1.Get(sg)
	assert.NoError(t, err)
	assert.Exactly(t, "Content-Type,X-CoreStore-ID", defaultStr)

	sg = cfgmock.NewService(cfgmock.PathValue{
		wantPath.String(): "X-CoreStore-TOKEN",
	}).NewScoped(wantWebsiteID, 0)

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

func TestBaseValue_InScope(t *testing.T) {

	tests := []struct {
		sg         config.Scoped
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
		p1 := newBaseValue("a/b/c", WithField(&element.Field{
			ID:     cfgpath.NewRoute(`c`),
			Scopes: test.p,
		}))
		haveErr := p1.InScope(test.sg.ScopeID())

		if test.wantErrBhf != nil {
			assert.True(t, test.wantErrBhf(haveErr), "Index %d => %s", i, haveErr)
		} else {
			assert.NoError(t, haveErr, "Index %d", i)
		}
	}
}

func TestBaseValue_InScope_Perm(t *testing.T) {
	bv := newBaseValue("x/y/z", WithScopeStore())
	assert.NoError(t, bv.inScope(scope.Store.Pack(0)))
	assert.NoError(t, bv.inScope(scope.Website.Pack(0)))

	bv = newBaseValue("x/y/z", WithScopeWebsite())
	assert.Error(t, bv.inScope(scope.Store.Pack(0)))
	assert.NoError(t, bv.inScope(scope.Website.Pack(0)))

	bv = newBaseValue("x/y/z")
	assert.Error(t, bv.inScope(scope.Store.Pack(0)))
	assert.Error(t, bv.inScope(scope.Website.Pack(0)))
}

func TestBaseValue_FQ(t *testing.T) {

	const pth = "aa/bb/cc"
	p := newBaseValue(pth, WithScopeStore())
	fq, err := p.FQ(scope.Store.Pack(4))
	assert.NoError(t, err, "%+v", err)
	assert.Exactly(t, cfgpath.MustNewByParts(pth).BindStore(4).String(), fq)
}

func TestBaseValueMustFQPanic(t *testing.T) {
	defer func() {
		if r := recover(); r != nil {
			err := r.(error)
			assert.True(t, errors.IsNotValid(err), "Error: %s", err)
		} else {
			t.Fatal("Expecting a panic")
		}
	}()
	const pth = "a/b/c"
	p := newBaseValue(pth, WithScopeStore())
	fq := p.MustFQ(scope.Website.Pack(4))
	assert.Empty(t, fq)
}

func TestBaseValueToPath(t *testing.T) {
	t.Run("Valid Route", testBaseValueToPath(cfgpath.NewRoute("aa/bb/cc"), scope.Website.Pack(23), nil))
	t.Run("Invalid Route", testBaseValueToPath(cfgpath.NewRoute("a/bb/cc"), scope.Website.Pack(23), errors.IsNotValid))
	t.Run("Unauthorized Route", testBaseValueToPath(cfgpath.NewRoute("aa/bb/cc"), scope.Store.Pack(22), errors.IsUnauthorized))
}

func testBaseValueToPath(route cfgpath.Route, h scope.TypeID, wantErrBhf errors.BehaviourFunc) func(*testing.T) {
	return func(t *testing.T) {
		bv := newBaseValue(route.String())

		bv.Field = &element.Field{
			ID:     cfgpath.NewRoute("cc"),
			Scopes: scope.PermWebsite, // only scope default and website are allowed
		}

		havePath, haveErr := bv.ToPath(h)
		if wantErrBhf != nil {
			// t.Log(haveErr)
			assert.True(t, wantErrBhf(haveErr), "Error: %s", haveErr)
			return
		}
		wantPath := cfgpath.MustNew(route).Bind(h)
		assert.Exactly(t, wantPath, havePath)
	}
}

func TestBaseValueRoute(t *testing.T) {

	org := newBaseValue("aa/bb/cc")
	clone := org.Route()

	if &(org.route) == &clone { // comparing pointer addresses
		// is there a better way to test of the slice headers points to a different location?
		// because clone should be a clone ;-)
		t.Error("Should not be equal")
	}
}

func TestBaseValue_IsSet(t *testing.T) {
	r := newBaseValue("aa/bb/cc")
	assert.True(t, r.IsSet())

	r = baseValue{}
	assert.False(t, r.IsSet())
}

func TestBaseValue_MustFQWebsite_Panic(t *testing.T) {
	defer func() {
		if r := recover(); r != nil {
			err := r.(error)
			assert.True(t, errors.IsNotValid(err), "Error: %s", err)
		} else {
			t.Fatal("Expecting a panic")
		}
	}()
	const pth = "a/b/c"
	p := newBaseValue(pth, WithScopeStore())
	fq := p.MustFQWebsite(4)
	assert.Empty(t, fq)
}

func TestBaseValue_MustFQStore_Panic(t *testing.T) {
	defer func() {
		if r := recover(); r != nil {
			err := r.(error)
			assert.True(t, errors.IsNotValid(err), "Error: %s", err)
		} else {
			t.Fatal("Expecting a panic")
		}
	}()
	const pth = "a/b/c"
	p := newBaseValue(pth, WithScopeStore())
	fq := p.MustFQStore(5)
	assert.Empty(t, fq)
}

func TestBaseValue_MustFQWebsite(t *testing.T) {
	const pth = "aa/bb/cc"
	p := newBaseValue(pth, WithScopeStore())
	fq := p.MustFQWebsite(4)
	assert.Exactly(t, `websites/4/aa/bb/cc`, fq)
}

func TestBaseValue_MustFQStore(t *testing.T) {
	const pth = "aa/bb/cc"
	p := newBaseValue(pth, WithScopeStore())
	fq := p.MustFQStore(5)
	assert.Exactly(t, `stores/5/aa/bb/cc`, fq)
}
