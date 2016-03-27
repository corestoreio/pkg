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

package element_test

import (
	"testing"

	"github.com/corestoreio/csfw/config/cfgpath"
	"github.com/corestoreio/csfw/config/element"
	"github.com/corestoreio/csfw/storage/text"
	"github.com/stretchr/testify/assert"
)

var _ element.Sectioner = (*element.SectionSlice)(nil)

func TestSectionValidateDuplicate(t *testing.T) {
	// for benchmark tests see package config_bm
	t.Parallel()
	ss := element.NewSectionSlice(
		element.Section{
			ID: cfgpath.NewRoute(`aa`),
			Groups: element.NewGroupSlice(
				element.Group{
					ID: cfgpath.NewRoute(`bb`),
					Fields: element.NewFieldSlice(
						element.Field{ID: cfgpath.NewRoute(`cc`)},
						element.Field{ID: cfgpath.NewRoute(`cc`)},
					),
				},
			),
		},
	)

	err := ss.Validate()
	assert.EqualError(t, err, "Duplicate entry for path aa/bb/cc :: [{\"ID\":\"aa\",\"Groups\":[{\"ID\":\"bb\",\"Fields\":[{\"ID\":\"cc\"},{\"ID\":\"cc\"}]}]}]\n")
}

func TestSectionValidateShortPath(t *testing.T) {
	t.Parallel()
	ss := element.NewSectionSlice(
		element.Section{
			//ID: cfgpath.NewRoute(`aa`),
			Groups: element.NewGroupSlice(
				element.Group{
					//ID: cfgpath.NewRoute(`b`),
					Fields: element.NewFieldSlice(
						element.Field{ID: cfgpath.NewRoute(`ca`)},
						element.Field{ID: cfgpath.NewRoute(`cb`)},
						element.Field{},
					),
				},
			),
		},
	)

	err := ss.Validate()
	assert.EqualError(t, err, cfgpath.ErrRouteEmpty.Error())

	if e2, ok := err.(*element.FieldError); ok {
		assert.Exactly(t, "", e2.Field.ID.String())
		assert.Exactly(t, "", e2.RenderRoutes())
	} else {
		t.Fatal("Cannot type assert to *element.FieldError in err variable")
	}
}

func TestSectionUpdateField(t *testing.T) {
	t.Parallel()
	ss := element.NewSectionSlice(
		element.Section{
			ID: cfgpath.NewRoute(`aa`),
			Groups: element.NewGroupSlice(
				element.Group{
					ID: cfgpath.NewRoute(`bb`),
					Fields: element.NewFieldSlice(
						element.Field{ID: cfgpath.NewRoute(`ca`)},
						element.Field{ID: cfgpath.NewRoute(`cb`)},
					),
				},
			),
		},
	)

	fr := cfgpath.NewRoute(`aa/bb/ca`)
	if err := ss.UpdateField(fr, element.Field{
		Label: text.Chars("ca New Label"),
	}); err != nil {
		t.Fatal(err)
	}

	f, _, err := ss.FindField(fr)
	if err != nil {
		t.Fatal(err)
	}
	assert.Exactly(t, `ca New Label`, f.Label.String())

	err1 := ss.UpdateField(cfgpath.NewRoute(`a/b/c`), element.Field{})
	assert.EqualError(t, err1, element.ErrSectionNotFound.Error())

	err2 := ss.UpdateField(cfgpath.NewRoute(`aa/b/c`), element.Field{})
	assert.EqualError(t, err2, element.ErrGroupNotFound.Error())

	err3 := ss.UpdateField(cfgpath.NewRoute(`aa/bb/c`), element.Field{})
	assert.EqualError(t, err3, element.ErrFieldNotFound.Error())

	err4 := ss.UpdateField(cfgpath.NewRoute(`aa_bb_c`), element.Field{})
	assert.EqualError(t, err4, cfgpath.ErrIncorrectPath.Error())
}
