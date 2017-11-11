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
	"encoding/json"
	"testing"

	"github.com/corestoreio/cspkg/config/cfgpath"
	"github.com/corestoreio/cspkg/config/element"
	"github.com/corestoreio/cspkg/storage/text"
	"github.com/corestoreio/cspkg/store/scope"
	"github.com/corestoreio/errors"
	"github.com/stretchr/testify/assert"
)

func TestNewConfiguration(t *testing.T) {

	tests := []struct {
		have       element.SectionSlice
		wantErrBhf errors.BehaviourFunc
		wantLen    int
	}{
		0: {
			have:       nil,
			wantErrBhf: errors.IsNotValid,
		},
		1: {
			have: element.NewSectionSlice(
				element.Section{
					ID: cfgpath.NewRoute(`web`),
					Groups: element.NewGroupSlice(
						element.Group{
							ID:     cfgpath.NewRoute(`default`),
							Fields: element.NewFieldSlice(element.Field{ID: cfgpath.NewRoute(`front`)}, element.Field{ID: cfgpath.NewRoute(`no_route`)}),
						},
					),
				},
				element.Section{
					ID: cfgpath.NewRoute(`system`),
					Groups: element.NewGroupSlice(
						element.Group{
							ID:     cfgpath.NewRoute(`media_storage_configuration`),
							Fields: element.NewFieldSlice(element.Field{ID: cfgpath.NewRoute(`allowed_resources`)}),
						},
					),
				},
			),
			wantErrBhf: nil,
			wantLen:    3,
		},
		2: {
			have:       element.NewSectionSlice(element.Section{ID: cfgpath.NewRoute(`aa`), Groups: element.NewGroupSlice()}),
			wantErrBhf: nil,
		},
		3: {
			have:       element.NewSectionSlice(element.Section{ID: cfgpath.NewRoute(`aa`), Groups: element.NewGroupSlice(element.Group{ID: cfgpath.NewRoute(`bb`), Fields: nil})}),
			wantErrBhf: nil,
		},
		4: {
			have: element.NewSectionSlice(
				element.Section{
					ID: cfgpath.NewRoute(`aa`),
					Groups: element.NewGroupSlice(
						element.Group{
							ID:     cfgpath.NewRoute(`bb`),
							Fields: element.NewFieldSlice(element.Field{ID: cfgpath.NewRoute(`cc`)}, element.Field{ID: cfgpath.NewRoute(`cc`)}),
						},
					),
				},
			),
			wantErrBhf: errors.IsNotValid, // `Duplicate entry for path aa/bb/cc :: [{"ID":"aa","Groups":[{"ID":"bb","Fields":[{"ID":"cc"},{"ID":"cc"}]}]}]`,
		},
	}

	for i, test := range tests {
		func(t *testing.T, have element.SectionSlice, wantErr errors.BehaviourFunc) {
			defer func() {
				if r := recover(); r != nil {
					if err, ok := r.(error); ok {
						assert.True(t, wantErr(err), "Index %d => %s", i, err)
					} else {
						t.Errorf("Failed to convert to type error: %#v", r)
					}
				} else if wantErr != nil {
					t.Errorf("Cannot find panic: wantErr %v", wantErr)
				}
			}()

			haveSlice := element.MustNewConfiguration(have...)
			if wantErr != nil {
				assert.Nil(t, haveSlice, "Index %d", i)
			} else {
				assert.NotNil(t, haveSlice, "Index %d", i)
				assert.Len(t, haveSlice, len(have), "Index %d", i)
			}
			assert.Exactly(t, test.wantLen, haveSlice.TotalFields(), "Index %d", i)
		}(t, test.have, test.wantErrBhf)
	}
}

func TestSectionSliceDefaults(t *testing.T) {

	pkgCfg := element.MustNewConfiguration(
		element.Section{
			ID: cfgpath.NewRoute(`contact`),
			Groups: element.NewGroupSlice(
				element.Group{
					ID: cfgpath.NewRoute(`contact`),
					Fields: element.NewFieldSlice(
						element.Field{
							// Path: `contact/contact/enabled`,
							ID:      cfgpath.NewRoute(`enabled`),
							Default: true,
						},
					),
				},
				element.Group{
					ID: cfgpath.NewRoute(`email`),
					Fields: element.NewFieldSlice(
						element.Field{
							// Path: `contact/email/recipient_email`,
							ID:      cfgpath.NewRoute(`recipient_email`),
							Default: `hello@example.com`,
						},
						element.Field{
							// Path: `contact/email/sender_email_identity`,
							ID:      cfgpath.NewRoute(`sender_email_identity`),
							Default: 2.7182818284590452353602874713527,
						},
						element.Field{
							// Path: `contact/email/email_template`,
							ID:      cfgpath.NewRoute(`email_template`),
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

	// Got stuck in comparing JSON?
	// Use a Webservice to compare the JSON output!

	tests := []struct {
		have       []element.SectionSlice
		wantErr    string
		want       string
		fieldCount int
	}{
		{
			have: []element.SectionSlice{
				{
					element.Section{
						ID: cfgpath.NewRoute(`a`),
					},
				},
				{
					element.Section{ID: cfgpath.NewRoute(`a`), Label: text.Chars(`LabelA`), Groups: nil},
				},
			},
			wantErr:    "",
			want:       `[{"ID":"a","Label":"LabelA","Groups":null}]` + "\n",
			fieldCount: 0,
		},
		{
			have: []element.SectionSlice{
				{
					element.Section{
						ID: cfgpath.NewRoute(`a`),
						Groups: element.NewGroupSlice(
							element.Group{
								ID: cfgpath.NewRoute(`b`),
								Fields: element.NewFieldSlice(
									element.Field{ID: cfgpath.NewRoute(`c`), Default: `c`},
								),
							},
							element.Group{
								ID: cfgpath.NewRoute(`b`),
								Fields: element.NewFieldSlice(
									element.Field{ID: cfgpath.NewRoute(`d`), Default: `d`},
								),
							},
						),
					},
				},
				{
					element.Section{ID: cfgpath.NewRoute(`a`), Label: text.Chars(`LabelA`), Groups: nil},
				},
			},
			wantErr:    "",
			want:       `[{"ID":"a","Label":"LabelA","Groups":[{"ID":"b","Fields":[{"ID":"c","Default":"c"},{"ID":"d","Default":"d"}]}]}]` + "\n",
			fieldCount: 2,
		},
		{
			have: []element.SectionSlice{
				{
					element.Section{
						ID:    cfgpath.NewRoute(`a`),
						Label: text.Chars(`SectionLabelA`),
						Groups: element.NewGroupSlice(
							element.Group{
								ID:     cfgpath.NewRoute(`b`),
								Scopes: scope.PermDefault,
								Fields: element.NewFieldSlice(
									element.Field{ID: cfgpath.NewRoute(`c`), Default: `c`},
								),
							},
						),
					},
				},
				{
					element.Section{
						ID:     cfgpath.NewRoute(`a`),
						Scopes: scope.PermWebsite,
						Groups: element.NewGroupSlice(
							element.Group{ID: cfgpath.NewRoute(`b`), Label: text.Chars(`GroupLabelB1`)},
							element.Group{ID: cfgpath.NewRoute(`b`), Label: text.Chars(`GroupLabelB2`)},
							element.Group{
								ID: cfgpath.NewRoute(`b2`),
								Fields: element.NewFieldSlice(
									element.Field{ID: cfgpath.NewRoute(`d`), Default: `d`},
								),
							},
						),
					},
				},
			},
			wantErr:    "",
			want:       `[{"ID":"a","Label":"SectionLabelA","Scopes":["Default","Website"],"Groups":[{"ID":"b","Label":"GroupLabelB2","Scopes":["Default"],"Fields":[{"ID":"c","Default":"c"}]},{"ID":"b2","Fields":[{"ID":"d","Default":"d"}]}]}]` + "\n",
			fieldCount: 2,
		},
		{
			have: []element.SectionSlice{
				{
					element.Section{ID: cfgpath.NewRoute(`a`), Label: text.Chars(`SectionLabelA`), SortOrder: 20, Resource: 22},
				},
				{
					element.Section{ID: cfgpath.NewRoute(`a`), Scopes: scope.PermWebsite, SortOrder: 10, Resource: 3},
				},
			},
			wantErr: "",
			want:    `[{"ID":"a","Label":"SectionLabelA","Scopes":["Default","Website"],"SortOrder":10,"Resource":3,"Groups":null}]` + "\n",
		},
		{
			have: []element.SectionSlice{
				{
					element.Section{
						ID:    cfgpath.NewRoute(`a`),
						Label: text.Chars(`SectionLabelA`),
						Groups: element.NewGroupSlice(
							element.Group{
								ID:      cfgpath.NewRoute(`b`),
								Label:   text.Chars(`SectionAGroupB`),
								Comment: text.Chars("SectionAGroupBComment"),
								Scopes:  scope.PermDefault,
							},
						),
					},
				},
				{
					element.Section{
						ID:        cfgpath.NewRoute(`a`),
						SortOrder: 1000,
						Scopes:    scope.PermWebsite,
						Groups: element.NewGroupSlice(
							element.Group{ID: cfgpath.NewRoute(`b`), Label: text.Chars(`GroupLabelB1`), Scopes: scope.PermStore},
							element.Group{ID: cfgpath.NewRoute(`b`), Label: text.Chars(`GroupLabelB2`), Comment: text.Chars("Section2AGroup3BComment"), SortOrder: 100},
							element.Group{ID: cfgpath.NewRoute(`b2`)},
						),
					},
				},
			},
			wantErr: "",
			want:    `[{"ID":"a","Label":"SectionLabelA","Scopes":["Default","Website"],"SortOrder":1000,"Groups":[{"ID":"b","Label":"GroupLabelB2","Comment":"Section2AGroup3BComment","Scopes":["Default","Website","Store"],"SortOrder":100,"Fields":null},{"ID":"b2","Fields":null}]}]` + "\n",
		},
		{
			have: []element.SectionSlice{
				{
					element.Section{
						ID: cfgpath.NewRoute(`a`),
						Groups: element.NewGroupSlice(
							element.Group{
								ID:    cfgpath.NewRoute(`b`),
								Label: text.Chars(`b1`),
								Fields: element.NewFieldSlice(
									element.Field{ID: cfgpath.NewRoute(`c`), Default: `c`, Type: element.TypeMultiselect, SortOrder: 1001},
								),
							},
							element.Group{
								ID:    cfgpath.NewRoute(`b`),
								Label: text.Chars(`b2`),
								Fields: element.NewFieldSlice(
									element.Field{ID: cfgpath.NewRoute(`d`), Default: `d`, Comment: text.Chars("Ring of fire"), Type: element.TypeObscure},
									element.Field{ID: cfgpath.NewRoute(`c`), Default: `haha`, Type: element.TypeSelect, Scopes: scope.PermWebsite},
								),
							},
						),
					},
				},
				{
					element.Section{
						ID: cfgpath.NewRoute(`a`),
						Groups: element.NewGroupSlice(
							element.Group{
								ID:    cfgpath.NewRoute(`b`),
								Label: text.Chars(`b3`),
								Fields: element.NewFieldSlice(
									element.Field{ID: cfgpath.NewRoute(`d`), Default: `overriddenD`, Label: text.Chars(`Sect2Group2Label4`), Comment: text.Chars("LOTR")},
									element.Field{ID: cfgpath.NewRoute(`c`), Default: `overriddenHaha`, Type: element.TypeHidden},
								),
							},
						),
					},
				},
			},
			wantErr:    "",
			want:       `[{"ID":"a","Groups":[{"ID":"b","Label":"b3","Fields":[{"ID":"c","Type":"hidden","Scopes":["Default","Website"],"SortOrder":1001,"Default":"overriddenHaha"},{"ID":"d","Type":"obscure","Label":"Sect2Group2Label4","Comment":"LOTR","Default":"overriddenD"}]}]}]` + "\n",
			fieldCount: 2,
		},
		{
			have: []element.SectionSlice{
				{
					element.Section{
						ID: cfgpath.NewRoute(`a`),
						Groups: element.NewGroupSlice(
							element.Group{
								ID: cfgpath.NewRoute(`b`),
								Fields: element.NewFieldSlice(
									element.Field{
										ID:      cfgpath.NewRoute(`c`),
										Default: `c`,
										Type:    element.TypeMultiselect,
									},
								),
							},
						),
					},
				},
				{
					element.Section{
						ID: cfgpath.NewRoute(`a`),
						Groups: element.NewGroupSlice(
							element.Group{
								ID: cfgpath.NewRoute(`b`),
								Fields: element.NewFieldSlice(
									element.Field{
										ID:        cfgpath.NewRoute(`c`),
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
			wantErr:    "",
			fieldCount: 1,
			want:       `[{"ID":"a","Groups":[{"ID":"b","Fields":[{"ID":"c","Type":"select","Label":"Sect2Group2Label4","Comment":"LOTR","SortOrder":100,"Visible":true,"Default":"overridenC"}]}]}]` + "\n",
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
		assert.Exactly(t, test.fieldCount, baseSl.TotalFields(), "Index %d", i)
	}
}

func TestFieldSliceMerge(t *testing.T) {

	tests := []struct {
		have    element.FieldSlice
		wantErr error
		want    string
	}{
		{
			have: element.NewFieldSlice(
				element.Field{ID: cfgpath.NewRoute(`d`), Default: `overrideMeD`, Comment: text.Chars("Ring of fire"), Type: element.TypeObscure},
			),
			wantErr: nil,
			want:    `[{"ID":"d","Type":"obscure","Comment":"Ring of fire","Default":"overrideMeD"}]`,
		},
		{
			have: element.NewFieldSlice(
				element.Field{ID: cfgpath.NewRoute(`d`), Default: `overrideMeD`, Comment: text.Chars("Ring of fire"), Type: element.TypeObscure},
				element.Field{ID: cfgpath.NewRoute(`c`), Default: `overrideMeC`, Type: element.TypeSelect, Scopes: scope.PermWebsite},
			),
			wantErr: nil,
			want:    `[{"ID":"d","Type":"obscure","Comment":"Ring of fire","Default":"overrideMeD"},{"ID":"c","Type":"select","Scopes":["Default","Website"],"Default":"overrideMeC"}]`,
		},
		{
			have: element.NewFieldSlice(
				element.Field{ID: cfgpath.NewRoute(`d`), Default: `overrideMeD`, Comment: text.Chars("Ring of fire"), Type: element.TypeObscure},
				element.Field{ID: cfgpath.NewRoute(`c`), Default: `overrideMeC`, Type: element.TypeSelect, Scopes: scope.PermWebsite},
				element.Field{ID: cfgpath.NewRoute(`d`), Default: `overrideMeE`, Type: element.TypeMultiselect},
			),
			wantErr: nil,
			want:    `[{"ID":"d","Type":"multiselect","Comment":"Ring of fire","Default":"overrideMeE"},{"ID":"c","Type":"select","Scopes":["Default","Website"],"Default":"overrideMeC"}]`,
		},
		{
			have:    nil,
			wantErr: nil,
			want:    `null`,
		},
	}

	for i, test := range tests {
		var baseFsl element.FieldSlice
		haveErr := baseFsl.Merge(test.have...)
		if test.wantErr != nil {
			assert.Len(t, baseFsl, 0)
			assert.Error(t, haveErr)
			assert.Contains(t, haveErr.Error(), test.wantErr)
		} else {
			assert.NoError(t, haveErr)
			fsj, err := json.Marshal(baseFsl)
			if err != nil {
				t.Fatal(err)
			}
			if string(fsj) != test.want {
				t.Errorf("\nIndex: %d\nExpected: %s\nActual:   %s\n", i, test.want, string(fsj))
			}
		}
	}
}

func TestGroupSliceMerge(t *testing.T) {

	tests := []struct {
		have    element.GroupSlice
		wantErr error
		want    string
	}{
		{
			have: element.NewGroupSlice(
				element.Group{
					ID: cfgpath.NewRoute(`b`),
					Fields: element.NewFieldSlice(
						element.Field{ID: cfgpath.NewRoute(`c`), Default: `c`, Type: element.TypeMultiselect},
					),
				},
				element.Group{
					ID: cfgpath.NewRoute(`b`),
					Fields: element.NewFieldSlice(
						element.Field{ID: cfgpath.NewRoute(`d`), Default: `overrideMeD`, Comment: text.Chars("Ring of fire"), Type: element.TypeObscure},
						element.Field{ID: cfgpath.NewRoute(`c`), Default: `overrideMeC`, Type: element.TypeSelect, Scopes: scope.PermWebsite},
					),
				},
				element.Group{
					ID: cfgpath.NewRoute(`b`),
					Fields: element.NewFieldSlice(
						element.Field{ID: cfgpath.NewRoute(`d`), Default: `overriddenD`, Label: text.Chars(`Sect2Group2Label4`), Comment: text.Chars("LOTR")},
						element.Field{ID: cfgpath.NewRoute(`c`), Default: `overriddenC`, Type: element.TypeHidden},
					),
				},
			),
			wantErr: nil,
			want:    `[{"ID":"b","Fields":[{"ID":"c","Type":"hidden","Scopes":["Default","Website"],"Default":"overriddenC"},{"ID":"d","Type":"obscure","Label":"Sect2Group2Label4","Comment":"LOTR","Default":"overriddenD"}]}]` + "\n",
		},
		{
			have: element.NewGroupSlice(
				element.Group{
					ID:    cfgpath.NewRoute(`b`),
					Label: text.Chars(`Single Field`),
					Fields: element.NewFieldSlice(
						element.Field{ID: cfgpath.NewRoute(`c`), Default: `c`, Type: element.TypeMultiselect},
					),
				},
			),
			wantErr: nil,
			want:    `[{"ID":"b","Label":"Single Field","Fields":[{"ID":"c","Type":"multiselect","Default":"c"}]}]` + "\n",
		},
		{
			have: element.NewGroupSlice(
				element.Group{
					ID:    cfgpath.NewRoute(`b`),
					Label: text.Chars(`Single Field`),
					Fields: element.NewFieldSlice(
						element.Field{ID: cfgpath.NewRoute(`c`), Default: `c`, Type: element.TypeMultiselect},
					),
				},
				element.Group{
					ID:    cfgpath.NewRoute(`b`),
					Label: text.Chars(`Single Field2`),
					Fields: element.NewFieldSlice(
						element.Field{ID: cfgpath.NewRoute(`c`), Default: `c`, Type: element.TypeMultiselect},
					),
				},
			),
			wantErr: nil,
			want:    `[{"ID":"b","Label":"Single Field2","Fields":[{"ID":"c","Type":"multiselect","Default":"c"}]}]` + "\n",
		},
		{
			have: element.NewGroupSlice(
				element.Group{
					ID:    cfgpath.NewRoute(`b`),
					Label: text.Chars(`Single Field`),
					Fields: element.NewFieldSlice(
						element.Field{ID: cfgpath.NewRoute(`c`), Default: `c`, Type: element.TypeMultiselect},
					),
				},
				element.Group{
					ID:    cfgpath.NewRoute(`b`),
					Label: text.Chars(`Single Field2`),
					Fields: element.NewFieldSlice(
						element.Field{ID: cfgpath.NewRoute(`c`), Default: `c2`, Type: element.TypeTextarea},
					),
				},
			),
			wantErr: nil,
			want:    `[{"ID":"b","Label":"Single Field2","Fields":[{"ID":"c","Type":"textarea","Default":"c2"}]}]` + "\n",
		},
		{
			have: element.NewGroupSlice(
				element.Group{
					ID:    cfgpath.NewRoute(`b`),
					Label: text.Chars(`Single Field`),
					Fields: element.NewFieldSlice(
						element.Field{ID: cfgpath.NewRoute(`c`), Default: `c`, Type: element.TypeMultiselect},
					),
				},
				element.Group{
					ID: cfgpath.NewRoute(`b`),
					Fields: element.NewFieldSlice(
						element.Field{ID: cfgpath.NewRoute(`d`), Default: `d`, Type: element.TypeText},
					),
				},
			),
			wantErr: nil,
			want:    `[{"ID":"b","Label":"Single Field","Fields":[{"ID":"c","Type":"multiselect","Default":"c"},{"ID":"d","Type":"text","Default":"d"}]}]` + "\n",
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

	tests := []struct {
		haveSlice  element.SectionSlice
		haveRoute  cfgpath.Route
		wantGID    string
		wantErrBhf errors.BehaviourFunc
	}{
		{
			haveSlice:  element.NewSectionSlice(element.Section{ID: cfgpath.NewRoute(`a`), Groups: element.NewGroupSlice(element.Group{ID: cfgpath.NewRoute(`b`)}, element.Group{ID: cfgpath.NewRoute(`bb`)})}),
			haveRoute:  cfgpath.NewRoute("a/b"),
			wantGID:    "b",
			wantErrBhf: nil,
		},
		{
			haveSlice:  element.NewSectionSlice(element.Section{ID: cfgpath.NewRoute(`a`), Groups: element.NewGroupSlice(element.Group{ID: cfgpath.NewRoute(`b`)}, element.Group{ID: cfgpath.NewRoute(`bb`)})}),
			haveRoute:  cfgpath.NewRoute("a/bc"),
			wantGID:    "b",
			wantErrBhf: errors.IsNotFound,
		},
		{
			haveSlice:  element.SectionSlice{},
			haveRoute:  cfgpath.Route{},
			wantGID:    "b",
			wantErrBhf: errors.IsNotFound,
		},
		{
			haveSlice:  element.NewSectionSlice(element.Section{ID: cfgpath.NewRoute(`a`), Groups: element.NewGroupSlice(element.Group{ID: cfgpath.NewRoute(`b`)}, element.Group{ID: cfgpath.NewRoute(`bb`)})}),
			haveRoute:  cfgpath.NewRoute("a", "bb", "cc"),
			wantGID:    "bb",
			wantErrBhf: nil,
		},
		{
			haveSlice:  element.NewSectionSlice(element.Section{ID: cfgpath.NewRoute(`a`), Groups: element.NewGroupSlice(element.Group{ID: cfgpath.NewRoute(`b`)}, element.Group{ID: cfgpath.NewRoute(`bb`)})}),
			haveRoute:  cfgpath.NewRoute("xa", "bb", "cc"),
			wantGID:    "",
			wantErrBhf: errors.IsNotFound,
		},
	}

	for i, test := range tests {
		haveGroup, _, haveErr := test.haveSlice.FindGroup(test.haveRoute)
		if test.wantErrBhf != nil {
			assert.Error(t, haveErr, "Index %d", i)
			assert.Exactly(t, element.Group{}, haveGroup)
			assert.True(t, test.wantErrBhf(haveErr), "Index %d => %s", i, haveErr)
			continue
		}

		assert.NoError(t, haveErr, "Index %d", i)
		assert.NotNil(t, haveGroup, "Index %d", i)
		assert.Exactly(t, test.wantGID, haveGroup.ID.String())
	}
}

func TestSectionSliceFindFieldByID(t *testing.T) {

	tests := []struct {
		haveSlice  element.SectionSlice
		haveRoute  cfgpath.Route
		wantFID    string
		wantErrBhf errors.BehaviourFunc
	}{
		{
			haveSlice:  element.NewSectionSlice(element.Section{ID: cfgpath.NewRoute(`aa`), Groups: element.NewGroupSlice(element.Group{ID: cfgpath.NewRoute(`bb`)}, element.Group{ID: cfgpath.NewRoute(`cc`)})}),
			haveRoute:  cfgpath.Route{},
			wantFID:    "",
			wantErrBhf: errors.IsNotValid,
		},
		{
			haveSlice:  element.NewSectionSlice(element.Section{ID: cfgpath.NewRoute(`a`), Groups: element.NewGroupSlice(element.Group{ID: cfgpath.NewRoute(`b`)}, element.Group{ID: cfgpath.NewRoute(`bb`)})}),
			haveRoute:  cfgpath.NewRoute("a/b"),
			wantFID:    "b",
			wantErrBhf: errors.IsNotFound,
		},
		{
			haveSlice:  element.NewSectionSlice(element.Section{ID: cfgpath.NewRoute(`a`), Groups: element.NewGroupSlice(element.Group{ID: cfgpath.NewRoute(`b`)}, element.Group{ID: cfgpath.NewRoute(`bb`)})}),
			haveRoute:  cfgpath.NewRoute("a/bc"),
			wantFID:    "b",
			wantErrBhf: errors.IsNotFound,
		},
		{
			haveSlice:  element.NewSectionSlice(),
			haveRoute:  cfgpath.Route{},
			wantFID:    "",
			wantErrBhf: errors.IsNotValid,
		},
		{
			haveSlice:  element.NewSectionSlice(element.Section{ID: cfgpath.NewRoute(`a`), Groups: element.NewGroupSlice(element.Group{ID: cfgpath.NewRoute(`b`)}, element.Group{ID: cfgpath.NewRoute(`bb`)})}),
			haveRoute:  cfgpath.NewRoute("a", "bb", "cc"),
			wantFID:    "bb",
			wantErrBhf: errors.IsNotFound,
		},
		{
			haveSlice:  element.NewSectionSlice(element.Section{ID: cfgpath.NewRoute(`a`), Groups: element.NewGroupSlice(element.Group{ID: cfgpath.NewRoute(`b`)}, element.Group{ID: cfgpath.NewRoute(`bb`)})}),
			haveRoute:  cfgpath.NewRoute("xa", "bb", "cc"),
			wantFID:    "",
			wantErrBhf: errors.IsNotFound,
		},
		{
			haveSlice:  element.NewSectionSlice(element.Section{ID: cfgpath.NewRoute(`a1`), Groups: element.NewGroupSlice(element.Group{ID: cfgpath.NewRoute(`b1`), Fields: element.NewFieldSlice(element.Field{ID: cfgpath.NewRoute(`c1`)})})}),
			haveRoute:  cfgpath.NewRoute("a1", "b1", "c1"),
			wantFID:    "c1",
			wantErrBhf: nil,
		},
	}

	for i, test := range tests {
		haveGroup, _, haveErr := test.haveSlice.FindField(test.haveRoute)
		if test.wantErrBhf != nil {
			assert.Error(t, haveErr, "Index %d", i)
			assert.Exactly(t, element.Field{}, haveGroup, "Index %d", i)
			assert.True(t, test.wantErrBhf(haveErr), "Index %d => %s", i, haveErr)
			continue
		}
		assert.NoError(t, haveErr, "Index %d", i)
		assert.NotNil(t, haveGroup, "Index %d", i)
		assert.Exactly(t, test.wantFID, haveGroup.ID.String(), "Index %d", i)
	}
}

func TestFieldSliceSort(t *testing.T) {

	want := []int{-10, 1, 10, 11, 20}
	fs := element.NewFieldSlice(
		element.Field{ID: cfgpath.NewRoute(`aa`), SortOrder: 20},
		element.Field{ID: cfgpath.NewRoute(`bb`), SortOrder: -10},
		element.Field{ID: cfgpath.NewRoute(`cc`), SortOrder: 10},
		element.Field{ID: cfgpath.NewRoute(`dd`), SortOrder: 11},
		element.Field{ID: cfgpath.NewRoute(`ee`), SortOrder: 1},
	)

	for i, f := range fs.Sort() {
		assert.EqualValues(t, want[i], f.SortOrder)
	}
}

func TestGroupSliceSort(t *testing.T) {

	want := []int{-10, 1, 10, 11, 20}
	gs := element.NewGroupSlice(
		element.Group{ID: cfgpath.NewRoute(`aa`), SortOrder: 20},
		element.Group{ID: cfgpath.NewRoute(`bb`), SortOrder: -10},
		element.Group{ID: cfgpath.NewRoute(`cc`), SortOrder: 10},
		element.Group{ID: cfgpath.NewRoute(`dd`), SortOrder: 11},
		element.Group{ID: cfgpath.NewRoute(`ee`), SortOrder: 1},
	)
	for i, f := range gs.Sort() {
		assert.EqualValues(t, want[i], f.SortOrder)
	}
}
func TestSectionSliceSort(t *testing.T) {

	want := []int{-10, 1, 10, 11, 20}
	ss := element.NewSectionSlice(
		element.Section{ID: cfgpath.NewRoute(`aa`), SortOrder: 20},
		element.Section{ID: cfgpath.NewRoute(`bb`), SortOrder: -10},
		element.Section{ID: cfgpath.NewRoute(`cc`), SortOrder: 10},
		element.Section{ID: cfgpath.NewRoute(`dd`), SortOrder: 11},
		element.Section{ID: cfgpath.NewRoute(`ee`), SortOrder: 1},
	)
	for i, f := range ss.Sort() {
		assert.EqualValues(t, want[i], f.SortOrder)
	}

}

func TestSectionSliceSortAll(t *testing.T) {

	want := `[{"ID":"bb","SortOrder":-10,"Groups":null},{"ID":"ee","SortOrder":1,"Groups":null},{"ID":"cc","SortOrder":10,"Groups":null},{"ID":"aa","SortOrder":20,"Groups":[{"ID":"bb","SortOrder":-10,"Fields":[{"ID":"bb","SortOrder":-10},{"ID":"ee","SortOrder":1},{"ID":"cc","SortOrder":10},{"ID":"dd","SortOrder":11},{"ID":"aa","SortOrder":20}]},{"ID":"ee","SortOrder":1,"Fields":null},{"ID":"dd","SortOrder":11,"Fields":[{"ID":"bb","SortOrder":-10},{"ID":"ee","SortOrder":1},{"ID":"cc","SortOrder":10},{"ID":"dd","SortOrder":11},{"ID":"aa","SortOrder":20}]},{"ID":"aa","SortOrder":20,"Fields":[{"ID":"bb","SortOrder":-10},{"ID":"ee","SortOrder":1},{"ID":"cc","SortOrder":10},{"ID":"dd","SortOrder":11},{"ID":"aa","SortOrder":20}]}]}]` + "\n"
	ss := element.MustNewConfiguration(
		element.Section{ID: cfgpath.NewRoute(`aa`), SortOrder: 20, Groups: element.NewGroupSlice(
			element.Group{
				ID:        cfgpath.NewRoute(`aa`),
				SortOrder: 20,
				Fields: element.NewFieldSlice(
					element.Field{ID: cfgpath.NewRoute(`aa`), SortOrder: 20},
					element.Field{ID: cfgpath.NewRoute(`bb`), SortOrder: -10},
					element.Field{ID: cfgpath.NewRoute(`cc`), SortOrder: 10},
					element.Field{ID: cfgpath.NewRoute(`dd`), SortOrder: 11},
					element.Field{ID: cfgpath.NewRoute(`ee`), SortOrder: 1},
				),
			},
			element.Group{
				ID:        cfgpath.NewRoute(`bb`),
				SortOrder: -10,
				Fields: element.NewFieldSlice(
					element.Field{ID: cfgpath.NewRoute(`aa`), SortOrder: 20},
					element.Field{ID: cfgpath.NewRoute(`bb`), SortOrder: -10},
					element.Field{ID: cfgpath.NewRoute(`cc`), SortOrder: 10},
					element.Field{ID: cfgpath.NewRoute(`dd`), SortOrder: 11},
					element.Field{ID: cfgpath.NewRoute(`ee`), SortOrder: 1},
				),
			},
			element.Group{
				ID:        cfgpath.NewRoute(`dd`),
				SortOrder: 11,
				Fields: element.NewFieldSlice(
					element.Field{ID: cfgpath.NewRoute(`aa`), SortOrder: 20},
					element.Field{ID: cfgpath.NewRoute(`bb`), SortOrder: -10},
					element.Field{ID: cfgpath.NewRoute(`cc`), SortOrder: 10},
					element.Field{ID: cfgpath.NewRoute(`dd`), SortOrder: 11},
					element.Field{ID: cfgpath.NewRoute(`ee`), SortOrder: 1},
				),
			},
			element.Group{ID: cfgpath.NewRoute(`ee`), SortOrder: 1},
		)},
		element.Section{ID: cfgpath.NewRoute(`bb`), SortOrder: -10},
		element.Section{ID: cfgpath.NewRoute(`cc`), SortOrder: 10},
		element.Section{ID: cfgpath.NewRoute(`ee`), SortOrder: 1},
	)
	assert.Exactly(t, 15, ss.TotalFields())
	ss.SortAll()
	have := ss.ToJSON()
	if want != have {
		t.Errorf("\nWant: %s\nHave: %s\n", want, have)
	}
}

func TestSectionSliceAppendFields(t *testing.T) {

	want := `[{"ID":"aa","Groups":[{"ID":"aa","Fields":[{"ID":"aa"},{"ID":"bb"},{"ID":"cc"}]}]}]` + "\n"
	ss := element.MustNewConfiguration(
		element.Section{
			ID: cfgpath.NewRoute(`aa`),
			Groups: element.NewGroupSlice(
				element.Group{ID: cfgpath.NewRoute(`aa`),
					Fields: element.NewFieldSlice(
						element.Field{ID: cfgpath.NewRoute(`aa`)},
						element.Field{ID: cfgpath.NewRoute(`bb`)},
					),
				},
			)},
	)
	assert.True(t, errors.IsNotFound(ss.AppendFields(cfgpath.NewRoute("aa/XX"))))

	assert.NoError(t, ss.AppendFields(cfgpath.NewRoute("aa/aa"), element.Field{ID: cfgpath.NewRoute(`cc`)}))
	have := ss.ToJSON()
	if want != have {
		t.Errorf("\nWant: %s\nHave: %s\n", want, have)
	}
}
