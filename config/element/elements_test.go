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
	goerr "errors"
	"testing"

	"github.com/corestoreio/csfw/config/element"
	"github.com/corestoreio/csfw/config/path"
	"github.com/corestoreio/csfw/storage/text"
	"github.com/corestoreio/csfw/store/scope"
	"github.com/juju/errors"
	"github.com/stretchr/testify/assert"
)

func TestNewConfiguration(t *testing.T) {
	t.Parallel()
	tests := []struct {
		have    []*element.Section
		wantErr string
		wantLen int
	}{
		0: {
			have:    []*element.Section{},
			wantErr: "SectionSlice is empty",
		},
		1: {
			have: []*element.Section{
				{
					ID: path.NewRoute(`web`),
					Groups: element.NewGroupSlice(
						&element.Group{
							ID:     path.NewRoute(`default`),
							Fields: element.FieldSlice{&element.Field{ID: path.NewRoute(`front`)}, &element.Field{ID: path.NewRoute(`no_route`)}},
						},
					),
				},
				{
					ID: path.NewRoute(`system`),
					Groups: element.NewGroupSlice(
						&element.Group{
							ID:     path.NewRoute(`media_storage_configuration`),
							Fields: element.FieldSlice{&element.Field{ID: path.NewRoute(`allowed_resources`)}},
						},
					),
				},
			},
			wantErr: "",
			wantLen: 3,
		},
		2: {
			have:    []*element.Section{{ID: path.NewRoute(`aa`), Groups: element.GroupSlice{}}},
			wantErr: "",
		},
		3: {
			have:    []*element.Section{{ID: path.NewRoute(`aa`), Groups: element.GroupSlice{&element.Group{ID: path.NewRoute(`bb`), Fields: nil}}}},
			wantErr: "",
		},
		4: {
			have: []*element.Section{
				{
					ID: path.NewRoute(`aa`),
					Groups: element.NewGroupSlice(
						&element.Group{
							ID:     path.NewRoute(`bb`),
							Fields: element.FieldSlice{&element.Field{ID: path.NewRoute(`cc`)}, &element.Field{ID: path.NewRoute(`cc`)}},
						},
					),
				},
			},
			wantErr: `Duplicate entry for path aa/bb/cc :: [{"ID":"aa","Groups":[{"ID":"bb","Fields":[{"ID":"cc"},{"ID":"cc"}]}]}]`,
		},
	}

	for i, test := range tests {
		func(t *testing.T, have []*element.Section, wantErr string) {
			defer func() {
				if r := recover(); r != nil {
					if err, ok := r.(error); ok {
						assert.Contains(t, err.Error(), wantErr, "Index %d", i)
					} else {
						t.Errorf("Failed to convert to type error: %#v", err)
					}
				} else if wantErr != "" {
					t.Errorf("Cannot find panic: wantErr %s", wantErr)
				}
			}()

			haveSlice := element.MustNewConfiguration(have...)
			if wantErr != "" {
				assert.Nil(t, haveSlice, "Index %d", i)
			} else {
				assert.NotNil(t, haveSlice, "Index %d", i)
				assert.Len(t, haveSlice, len(have), "Index %d", i)
			}
			assert.Exactly(t, test.wantLen, haveSlice.TotalFields(), "Index %d", i)
		}(t, test.have, test.wantErr)
	}
}

func TestSectionSliceDefaults(t *testing.T) {
	t.Parallel()
	pkgCfg := element.MustNewConfiguration(
		&element.Section{
			ID: path.NewRoute(`contact`),
			Groups: element.NewGroupSlice(
				&element.Group{
					ID: path.NewRoute(`contact`),
					Fields: element.NewFieldSlice(
						&element.Field{
							// Path: `contact/contact/enabled`,
							ID:      path.NewRoute(`enabled`),
							Default: true,
						},
					),
				},
				&element.Group{
					ID: path.NewRoute(`email`),
					Fields: element.NewFieldSlice(
						&element.Field{
							// Path: `contact/email/recipient_email`,
							ID:      path.NewRoute(`recipient_email`),
							Default: `hello@example.com`,
						},
						&element.Field{
							// Path: `contact/email/sender_email_identity`,
							ID:      path.NewRoute(`sender_email_identity`),
							Default: 2.7182818284590452353602874713527,
						},
						&element.Field{
							// Path: `contact/email/email_template`,
							ID:      path.NewRoute(`email_template`),
							Default: 4711,
						},
					),
				},
			),
		},
	)

	dm, err := pkgCfg.Defaults()
	assert.NoError(t, err)
	assert.Exactly(
		t,
		element.DefaultMap{"contact/email/sender_email_identity": 2.718281828459045, "contact/email/email_template": 4711, "contact/contact/enabled": true, "contact/email/recipient_email": "hello@example.com"},
		dm,
	)
}

func TestSectionSliceMerge(t *testing.T) {
	t.Parallel()
	// Got stuck in comparing JSON?
	// Use a Webservice to compare the JSON output!

	tests := []struct {
		have    []element.SectionSlice
		wantErr string
		want    string
		wantLen int
	}{
		0: {
			have: []element.SectionSlice{
				nil,
				{
					nil,
					&element.Section{
						ID: path.NewRoute(`a`),
						Groups: element.NewGroupSlice(
							nil,
							&element.Group{
								ID: path.NewRoute(`b`),
								Fields: element.NewFieldSlice(
									&element.Field{ID: path.NewRoute(`c`), Default: `c`},
								),
							},
							&element.Group{
								ID: path.NewRoute(`b`),
								Fields: element.NewFieldSlice(
									&element.Field{ID: path.NewRoute(`d`), Default: `d`},
								),
							},
						),
					},
				},
				{
					&element.Section{ID: path.NewRoute(`a`), Label: text.Chars(`LabelA`), Groups: nil},
				},
			},
			wantErr: "",
			want:    `[{"ID":"a","Label":"LabelA","Groups":[{"ID":"b","Fields":[{"ID":"c","Default":"c"},{"ID":"d","Default":"d"}]}]}]` + "\n",
			wantLen: 2,
		},
		1: {
			have: []element.SectionSlice{
				{
					&element.Section{
						ID:    path.NewRoute(`a`),
						Label: text.Chars(`SectionLabelA`),
						Groups: element.NewGroupSlice(
							&element.Group{
								ID:    path.NewRoute(`b`),
								Scope: scope.PermDefault,
								Fields: element.NewFieldSlice(
									&element.Field{ID: path.NewRoute(`c`), Default: `c`},
								),
							},
							nil,
						),
					},
				},
				{
					&element.Section{
						ID:    path.NewRoute(`a`),
						Scope: scope.PermWebsite,
						Groups: element.NewGroupSlice(
							&element.Group{ID: path.NewRoute(`b`), Label: text.Chars(`GroupLabelB1`)},
							nil,
							&element.Group{ID: path.NewRoute(`b`), Label: text.Chars(`GroupLabelB2`)},
							&element.Group{
								ID: path.NewRoute(`b2`),
								Fields: element.NewFieldSlice(
									&element.Field{ID: path.NewRoute(`d`), Default: `d`},
								),
							},
						),
					},
				},
			},
			wantErr: "",
			want:    `[{"ID":"a","Label":"SectionLabelA","Scope":["Default","Website"],"Groups":[{"ID":"b","Label":"GroupLabelB2","Scope":["Default"],"Fields":[{"ID":"c","Default":"c"}]},{"ID":"b2","Fields":[{"ID":"d","Default":"d"}]}]}]` + "\n",
			wantLen: 2,
		},
		2: {
			have: []element.SectionSlice{
				{
					&element.Section{ID: path.NewRoute(`a`), Label: text.Chars(`SectionLabelA`), SortOrder: 20, Resource: 22},
				},
				{
					&element.Section{ID: path.NewRoute(`a`), Scope: scope.PermWebsite, SortOrder: 10, Resource: 3},
				},
			},
			wantErr: "",
			want:    `[{"ID":"a","Label":"SectionLabelA","Scope":["Default","Website"],"SortOrder":10,"Resource":3,"Groups":null}]` + "\n",
		},
		3: {
			have: []element.SectionSlice{
				{
					&element.Section{
						ID:    path.NewRoute(`a`),
						Label: text.Chars(`SectionLabelA`),
						Groups: element.NewGroupSlice(
							&element.Group{
								ID:      path.NewRoute(`b`),
								Label:   text.Chars(`SectionAGroupB`),
								Comment: text.Chars("SectionAGroupBComment"),
								Scope:   scope.PermDefault,
							},
						),
					},
				},
				{
					&element.Section{
						ID:        path.NewRoute(`a`),
						SortOrder: 1000,
						Scope:     scope.PermWebsite,
						Groups: element.NewGroupSlice(
							&element.Group{ID: path.NewRoute(`b`), Label: text.Chars(`GroupLabelB1`), Scope: scope.PermStore},
							&element.Group{ID: path.NewRoute(`b`), Label: text.Chars(`GroupLabelB2`), Comment: text.Chars("Section2AGroup3BComment"), SortOrder: 100},
							&element.Group{ID: path.NewRoute(`b2`)},
						),
					},
				},
			},
			wantErr: "",
			want:    `[{"ID":"a","Label":"SectionLabelA","Scope":["Default","Website"],"SortOrder":1000,"Groups":[{"ID":"b","Label":"GroupLabelB2","Comment":"Section2AGroup3BComment","Scope":["Default","Website","Store"],"SortOrder":100,"Fields":null},{"ID":"b2","Fields":null}]}]` + "\n",
		},
		4: {
			have: []element.SectionSlice{
				{
					&element.Section{
						ID: path.NewRoute(`a`),
						Groups: element.NewGroupSlice(
							&element.Group{
								ID:    path.NewRoute(`b`),
								Label: text.Chars(`b1`),
								Fields: element.NewFieldSlice(
									&element.Field{ID: path.NewRoute(`c`), Default: `c`, Type: element.TypeMultiselect, SortOrder: 1001},
								),
							},
							&element.Group{
								ID:    path.NewRoute(`b`),
								Label: text.Chars(`b2`),
								Fields: element.NewFieldSlice(
									nil,
									&element.Field{ID: path.NewRoute(`d`), Default: `d`, Comment: text.Chars("Ring of fire"), Type: element.TypeObscure},
									&element.Field{ID: path.NewRoute(`c`), Default: `haha`, Type: element.TypeSelect, Scope: scope.PermWebsite},
								),
							},
						),
					},
				},
				{
					&element.Section{
						ID: path.NewRoute(`a`),
						Groups: element.NewGroupSlice(
							&element.Group{
								ID:    path.NewRoute(`b`),
								Label: text.Chars(`b3`),
								Fields: element.NewFieldSlice(
									&element.Field{ID: path.NewRoute(`d`), Default: `overriddenD`, Label: text.Chars(`Sect2Group2Label4`), Comment: text.Chars("LOTR")},
									&element.Field{ID: path.NewRoute(`c`), Default: `overriddenHaha`, Type: element.TypeHidden},
								),
							},
						),
					},
				},
			},
			wantErr: "",
			want:    `[{"ID":"a","Groups":[{"ID":"b","Label":"b3","Fields":[{"ID":"c","Type":"hidden","Scope":["Default","Website"],"SortOrder":1001,"Default":"overriddenHaha"},{"ID":"d","Type":"obscure","Label":"Sect2Group2Label4","Comment":"LOTR","Default":"overriddenD"}]}]}]` + "\n",
			wantLen: 2,
		},
		5: {
			have: []element.SectionSlice{
				{
					&element.Section{
						ID: path.NewRoute(`a`),
						Groups: element.NewGroupSlice(
							&element.Group{
								ID: path.NewRoute(`b`),
								Fields: element.NewFieldSlice(
									&element.Field{
										ID:      path.NewRoute(`c`),
										Default: `c`,
										Type:    element.TypeMultiselect,
									},
								),
							},
						),
					},
				},
				{
					nil,
					&element.Section{
						ID: path.NewRoute(`a`),
						Groups: element.NewGroupSlice(
							&element.Group{
								ID: path.NewRoute(`b`),
								Fields: element.NewFieldSlice(
									nil,
									&element.Field{
										ID:        path.NewRoute(`c`),
										Default:   `overridenC`,
										Type:      element.TypeSelect,
										Label:     text.Chars(`Sect2Group2Label4`),
										Comment:   text.Chars("LOTR"),
										SortOrder: 100,
										Visible:   element.VisibleYes,
									},
								),
							},
						),
					},
				},
			},
			wantErr: "",
			wantLen: 1,
			want:    `[{"ID":"a","Groups":[{"ID":"b","Fields":[{"ID":"c","Type":"select","Label":"Sect2Group2Label4","Comment":"LOTR","SortOrder":100,"Visible":true,"Default":"overridenC"}]}]}]` + "\n",
		},
	}

	for i, test := range tests {

		if len(test.have) == 0 {
			test.want = "null\n"
		}

		var baseSl element.SectionSlice
		haveErr := baseSl.MergeMultiple(test.have...)
		if test.wantErr != "" {
			assert.Len(t, baseSl, 0)
			assert.Error(t, haveErr)
			assert.Contains(t, haveErr.Error(), test.wantErr)
		} else {
			assert.NoError(t, haveErr)
			j := baseSl.ToJSON()
			if j != test.want {
				t.Errorf("\nIndex: %d\nExpected: %s\nActual:   %s\n", i, test.want, j)
			}
		}
		assert.Exactly(t, test.wantLen, baseSl.TotalFields(), "Index %d", i)
	}
}

func TestGroupSliceMerge(t *testing.T) {
	t.Parallel()
	tests := []struct {
		have    []*element.Group
		wantErr error
		want    string
	}{
		{
			have: []*element.Group{
				{
					ID: path.NewRoute(`b`),
					Fields: element.NewFieldSlice(
						&element.Field{ID: path.NewRoute(`c`), Default: `c`, Type: element.TypeMultiselect},
					),
				},
				{
					ID: path.NewRoute(`b`),
					Fields: element.NewFieldSlice(
						&element.Field{ID: path.NewRoute(`d`), Default: `d`, Comment: text.Chars("Ring of fire"), Type: element.TypeObscure},
						&element.Field{ID: path.NewRoute(`c`), Default: `haha`, Type: element.TypeSelect, Scope: scope.PermWebsite},
					),
				},
				{
					ID: path.NewRoute(`b`),
					Fields: element.NewFieldSlice(
						&element.Field{ID: path.NewRoute(`d`), Default: `overriddenD`, Label: text.Chars(`Sect2Group2Label4`), Comment: text.Chars("LOTR")},
						&element.Field{ID: path.NewRoute(`c`), Default: `overriddenHaha`, Type: element.TypeHidden},
					),
				},
			},
			wantErr: nil,
			want:    `[{"ID":"b","Fields":[{"ID":"c","Type":"hidden","Scope":["Default","Website"],"Default":"overriddenHaha"},{"ID":"d","Type":"obscure","Label":"Sect2Group2Label4","Comment":"LOTR","Default":"overriddenD"}]}]` + "\n",
		},
		{
			have:    nil,
			wantErr: nil,
			want:    `null` + "\n",
		},
	}

	for i, test := range tests {
		var baseGsl element.GroupSlice
		haveErr := baseGsl.Merge(test.have...)
		if test.wantErr != nil {
			assert.Len(t, baseGsl, 0)
			assert.Error(t, haveErr)
			assert.Contains(t, haveErr.Error(), test.wantErr)
		} else {
			assert.NoError(t, haveErr)
			j := baseGsl.ToJSON()
			if j != test.want {
				t.Errorf("\nIndex: %d\nExpected: %s\nActual:   %s\n", i, test.want, j)
			}
		}
	}
}

func TestSectionSliceFindGroupByID(t *testing.T) {
	t.Parallel()
	tests := []struct {
		haveSlice element.SectionSlice
		haveRoute path.Route
		wantGID   string
		wantErr   error
	}{
		{
			haveSlice: element.SectionSlice{&element.Section{ID: path.NewRoute(`a`), Groups: element.NewGroupSlice(&element.Group{ID: path.NewRoute(`b`)}, &element.Group{ID: path.NewRoute(`bb`)})}},
			haveRoute: path.NewRoute("a/b"),
			wantGID:   "b",
			wantErr:   nil,
		},
		{
			haveSlice: element.SectionSlice{&element.Section{ID: path.NewRoute(`a`), Groups: element.NewGroupSlice(&element.Group{ID: path.NewRoute(`b`)}, &element.Group{ID: path.NewRoute(`bb`)})}},
			haveRoute: path.NewRoute("a/bc"),
			wantGID:   "b",
			wantErr:   element.ErrGroupNotFound,
		},
		{
			haveSlice: element.SectionSlice{},
			haveRoute: path.Route{},
			wantGID:   "b",
			wantErr:   element.ErrGroupNotFound,
		},
		{
			haveSlice: element.SectionSlice{&element.Section{ID: path.NewRoute(`a`), Groups: element.GroupSlice{&element.Group{ID: path.NewRoute(`b`)}, &element.Group{ID: path.NewRoute(`bb`)}}}},
			haveRoute: path.NewRoute("a", "bb", "cc"),
			wantGID:   "bb",
			wantErr:   nil,
		},
		{
			haveSlice: element.SectionSlice{&element.Section{ID: path.NewRoute(`a`), Groups: element.GroupSlice{&element.Group{ID: path.NewRoute(`b`)}, &element.Group{ID: path.NewRoute(`bb`)}}}},
			haveRoute: path.NewRoute("xa", "bb", "cc"),
			wantGID:   "",
			wantErr:   element.ErrSectionNotFound,
		},
	}

	for i, test := range tests {
		haveGroup, haveErr := test.haveSlice.FindGroupByID(test.haveRoute)
		if test.wantErr != nil {
			assert.Error(t, haveErr, "Index %d", i)
			assert.Nil(t, haveGroup)
			assert.EqualError(t, haveErr, test.wantErr.Error())
			continue
		}

		assert.NoError(t, haveErr, "Index %d", i)
		assert.NotNil(t, haveGroup, "Index %d", i)
		assert.Exactly(t, test.wantGID, haveGroup.ID.String())
	}
}

func TestSectionSliceFindFieldByID(t *testing.T) {
	t.Parallel()
	tests := []struct {
		haveSlice element.SectionSlice
		haveRoute path.Route
		wantFID   string
		wantErr   error
	}{
		{
			haveSlice: element.SectionSlice{&element.Section{ID: path.NewRoute(`aa`), Groups: element.GroupSlice{&element.Group{ID: path.NewRoute(`bb`)}, &element.Group{ID: path.NewRoute(`cc`)}}}},
			haveRoute: path.Route{},
			wantFID:   "",
			wantErr:   path.ErrIncorrectPath,
		},
		{
			haveSlice: element.SectionSlice{&element.Section{ID: path.NewRoute(`a`), Groups: element.GroupSlice{&element.Group{ID: path.NewRoute(`b`)}, &element.Group{ID: path.NewRoute(`bb`)}}}},
			haveRoute: path.NewRoute("a/b"),
			wantFID:   "b",
			wantErr:   element.ErrFieldNotFound,
		},
		{
			haveSlice: element.SectionSlice{&element.Section{ID: path.NewRoute(`a`), Groups: element.GroupSlice{&element.Group{ID: path.NewRoute(`b`)}, &element.Group{ID: path.NewRoute(`bb`)}}}},
			haveRoute: path.NewRoute("a/bc"),
			wantFID:   "b",
			wantErr:   element.ErrGroupNotFound,
		},
		{
			haveSlice: element.SectionSlice{nil},
			haveRoute: path.Route{},
			wantFID:   "",
			wantErr:   path.ErrIncorrectPath,
		},
		{
			haveSlice: element.SectionSlice{&element.Section{ID: path.NewRoute(`a`), Groups: element.GroupSlice{nil, &element.Group{ID: path.NewRoute(`b`)}, &element.Group{ID: path.NewRoute(`bb`)}}}},
			haveRoute: path.NewRoute("a", "bb", "cc"),
			wantFID:   "bb",
			wantErr:   element.ErrFieldNotFound,
		},
		{
			haveSlice: element.SectionSlice{&element.Section{ID: path.NewRoute(`a`), Groups: element.GroupSlice{&element.Group{ID: path.NewRoute(`b`)}, &element.Group{ID: path.NewRoute(`bb`)}}}},
			haveRoute: path.NewRoute("xa", "bb", "cc"),
			wantFID:   "",
			wantErr:   element.ErrSectionNotFound,
		},
		{
			haveSlice: element.SectionSlice{&element.Section{ID: path.NewRoute(`a1`), Groups: element.GroupSlice{&element.Group{ID: path.NewRoute(`b1`), Fields: element.FieldSlice{
				&element.Field{ID: path.NewRoute(`c1`)}, nil,
			}}}}},
			haveRoute: path.NewRoute("a1", "b1", "c1"),
			wantFID:   "c1",
			wantErr:   nil,
		},
	}

	for i, test := range tests {
		haveGroup, haveErr := test.haveSlice.FindFieldByID(test.haveRoute)
		if test.wantErr != nil {
			assert.Error(t, haveErr, "Index %d", i)
			assert.Nil(t, haveGroup, "Index %d", i)
			assert.EqualError(t, haveErr, test.wantErr.Error(), "Index %d", i)
			continue
		}
		assert.NoError(t, haveErr, "Index %d", i)
		assert.NotNil(t, haveGroup, "Index %d", i)
		assert.Exactly(t, test.wantFID, haveGroup.ID.String(), "Index %d", i)
	}
}

func TestFieldSliceSort(t *testing.T) {
	t.Parallel()
	want := []int{-10, 1, 10, 11, 20}
	fs := element.FieldSlice{
		&element.Field{ID: path.NewRoute(`aa`), SortOrder: 20},
		&element.Field{ID: path.NewRoute(`bb`), SortOrder: -10},
		&element.Field{ID: path.NewRoute(`cc`), SortOrder: 10},
		&element.Field{ID: path.NewRoute(`dd`), SortOrder: 11},
		&element.Field{ID: path.NewRoute(`ee`), SortOrder: 1},
	}

	for i, f := range *(fs.Sort()) {
		assert.EqualValues(t, want[i], f.SortOrder)
	}
}

func TestGroupSliceSort(t *testing.T) {
	t.Parallel()
	want := []int{-10, 1, 10, 11, 20}
	gs := element.GroupSlice{
		&element.Group{ID: path.NewRoute(`aa`), SortOrder: 20},
		&element.Group{ID: path.NewRoute(`bb`), SortOrder: -10},
		&element.Group{ID: path.NewRoute(`cc`), SortOrder: 10},
		&element.Group{ID: path.NewRoute(`dd`), SortOrder: 11},
		&element.Group{ID: path.NewRoute(`ee`), SortOrder: 1},
	}
	for i, f := range *(gs.Sort()) {
		assert.EqualValues(t, want[i], f.SortOrder)
	}
}
func TestSectionSliceSort(t *testing.T) {
	t.Parallel()
	want := []int{-10, 1, 10, 11, 20}
	ss := element.SectionSlice{
		&element.Section{ID: path.NewRoute(`aa`), SortOrder: 20},
		&element.Section{ID: path.NewRoute(`bb`), SortOrder: -10},
		&element.Section{ID: path.NewRoute(`cc`), SortOrder: 10},
		&element.Section{ID: path.NewRoute(`dd`), SortOrder: 11},
		&element.Section{ID: path.NewRoute(`ee`), SortOrder: 1},
	}
	for i, f := range *(ss.Sort()) {
		assert.EqualValues(t, want[i], f.SortOrder)
	}

}

func TestSectionSliceSortAll(t *testing.T) {
	t.Parallel()
	want := `[{"ID":"bb","SortOrder":-10,"Groups":null},{"ID":"ee","SortOrder":1,"Groups":null},{"ID":"cc","SortOrder":10,"Groups":null},{"ID":"aa","SortOrder":20,"Groups":[{"ID":"bb","SortOrder":-10,"Fields":[{"ID":"bb","SortOrder":-10},{"ID":"ee","SortOrder":1},{"ID":"cc","SortOrder":10},{"ID":"dd","SortOrder":11},{"ID":"aa","SortOrder":20}]},{"ID":"ee","SortOrder":1,"Fields":null},{"ID":"dd","SortOrder":11,"Fields":[{"ID":"bb","SortOrder":-10},{"ID":"ee","SortOrder":1},{"ID":"cc","SortOrder":10},{"ID":"dd","SortOrder":11},{"ID":"aa","SortOrder":20}]},{"ID":"aa","SortOrder":20,"Fields":[{"ID":"bb","SortOrder":-10},{"ID":"ee","SortOrder":1},{"ID":"cc","SortOrder":10},{"ID":"dd","SortOrder":11},{"ID":"aa","SortOrder":20}]}]}]` + "\n"
	ss := element.MustNewConfiguration(
		&element.Section{ID: path.NewRoute(`aa`), SortOrder: 20, Groups: element.NewGroupSlice(
			&element.Group{
				ID:        path.NewRoute(`aa`),
				SortOrder: 20,
				Fields: element.NewFieldSlice(
					&element.Field{ID: path.NewRoute(`aa`), SortOrder: 20},
					&element.Field{ID: path.NewRoute(`bb`), SortOrder: -10},
					&element.Field{ID: path.NewRoute(`cc`), SortOrder: 10},
					&element.Field{ID: path.NewRoute(`dd`), SortOrder: 11},
					&element.Field{ID: path.NewRoute(`ee`), SortOrder: 1},
				),
			},
			&element.Group{
				ID:        path.NewRoute(`bb`),
				SortOrder: -10,
				Fields: element.NewFieldSlice(
					&element.Field{ID: path.NewRoute(`aa`), SortOrder: 20},
					&element.Field{ID: path.NewRoute(`bb`), SortOrder: -10},
					&element.Field{ID: path.NewRoute(`cc`), SortOrder: 10},
					&element.Field{ID: path.NewRoute(`dd`), SortOrder: 11},
					&element.Field{ID: path.NewRoute(`ee`), SortOrder: 1},
				),
			},
			&element.Group{
				ID:        path.NewRoute(`dd`),
				SortOrder: 11,
				Fields: element.NewFieldSlice(
					&element.Field{ID: path.NewRoute(`aa`), SortOrder: 20},
					&element.Field{ID: path.NewRoute(`bb`), SortOrder: -10},
					&element.Field{ID: path.NewRoute(`cc`), SortOrder: 10},
					&element.Field{ID: path.NewRoute(`dd`), SortOrder: 11},
					&element.Field{ID: path.NewRoute(`ee`), SortOrder: 1},
				),
			},
			&element.Group{ID: path.NewRoute(`ee`), SortOrder: 1},
		)},
		&element.Section{ID: path.NewRoute(`bb`), SortOrder: -10},
		&element.Section{ID: path.NewRoute(`cc`), SortOrder: 10},
		&element.Section{ID: path.NewRoute(`ee`), SortOrder: 1},
	)
	assert.Exactly(t, 15, ss.TotalFields())
	ss.SortAll()
	have := ss.ToJSON()
	if want != have {
		t.Errorf("\nWant: %s\nHave: %s\n", want, have)
	}
}

func TestSectionSliceAppendFields(t *testing.T) {
	t.Parallel()
	want := `[{"ID":"aa","Groups":[{"ID":"aa","Fields":[{"ID":"aa"},{"ID":"bb"},{"ID":"cc"}]}]}]` + "\n"
	ss := element.MustNewConfiguration(
		&element.Section{
			ID: path.NewRoute(`aa`),
			Groups: element.NewGroupSlice(
				&element.Group{ID: path.NewRoute(`aa`),
					Fields: element.NewFieldSlice(
						&element.Field{ID: path.NewRoute(`aa`)},
						&element.Field{ID: path.NewRoute(`bb`)},
					),
				},
			)},
	)
	assert.EqualError(t, ss.AppendFields(path.NewRoute("aa/XX")), element.ErrGroupNotFound.Error())

	assert.NoError(t, ss.AppendFields(path.NewRoute("aa/aa"), &element.Field{ID: path.NewRoute(`cc`)}))
	have := ss.ToJSON()
	if want != have {
		t.Errorf("\nWant: %s\nHave: %s\n", want, have)
	}
}

func TestNotNotFoundError(t *testing.T) {
	t.Parallel()
	tests := []struct {
		have error
		want bool
	}{
		{goerr.New("PHP"), true},
		{errors.Mask(goerr.New("Scala")), true},
		{errors.New("Java"), true},
		{errors.Mask(errors.New("Java")), true},
		{element.ErrSectionNotFound, false},
		{element.ErrGroupNotFound, false},
		{element.ErrFieldNotFound, false},
		{errors.Mask(element.ErrFieldNotFound), false},
		{errors.Maskf(errors.Mask(element.ErrFieldNotFound), "A field not found error"), false},
		{nil, false},
		{errors.Mask(nil), false},
	}
	for i, test := range tests {
		assert.Exactly(t, test.want, element.NotNotFoundError(test.have), "Index %d", i)
	}
}
