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

	"github.com/corestoreio/csfw/config/element"
	"github.com/corestoreio/csfw/store/scope"
	"github.com/stretchr/testify/assert"
)

func TestNewConfiguration(t *testing.T) {
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
					ID: "web",
					Groups: element.NewGroupSlice(
						&element.Group{
							ID:     "default",
							Fields: element.FieldSlice{&element.Field{ID: "front"}, &element.Field{ID: "no_route"}},
						},
					),
				},
				{
					ID: "system",
					Groups: element.NewGroupSlice(
						&element.Group{
							ID:     "media_storage_configuration",
							Fields: element.FieldSlice{&element.Field{ID: "allowed_resources"}},
						},
					),
				},
			},
			wantErr: "",
			wantLen: 3,
		},
		2: {
			have:    []*element.Section{{ID: "aa", Groups: element.GroupSlice{}}},
			wantErr: "",
		},
		3: {
			have:    []*element.Section{{ID: "aa", Groups: element.GroupSlice{&element.Group{ID: "bb", Fields: nil}}}},
			wantErr: "",
		},
		4: {
			have: []*element.Section{
				{
					ID: "aa",
					Groups: element.NewGroupSlice(
						&element.Group{
							ID:     "bb",
							Fields: element.FieldSlice{&element.Field{ID: "cc"}, &element.Field{ID: "cc"}},
						},
					),
				},
			},
			wantErr: `Duplicate entry for path default/0/aa/bb/cc :: [{"ID":"aa","Groups":[{"ID":"bb","Fields":[{"ID":"cc"},{"ID":"cc"}]}]}]`,
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
	pkgCfg := element.MustNewConfiguration(
		&element.Section{
			ID: "contact",
			Groups: element.NewGroupSlice(
				&element.Group{
					ID: "contact",
					Fields: element.NewFieldSlice(
						&element.Field{
							// Path: `contact/contact/enabled`,
							ID:      "enabled",
							Default: true,
						},
					),
				},
				&element.Group{
					ID: "email",
					Fields: element.NewFieldSlice(
						&element.Field{
							// Path: `contact/email/recipient_email`,
							ID:      "recipient_email",
							Default: `hello@example.com`,
						},
						&element.Field{
							// Path: `contact/email/sender_email_identity`,
							ID:      "sender_email_identity",
							Default: 2.7182818284590452353602874713527,
						},
						&element.Field{
							// Path: `contact/email/email_template`,
							ID:      "email_template",
							Default: 4711,
						},
					),
				},
			),
		},
	)

	assert.Exactly(
		t,
		element.DefaultMap{"default/0/contact/email/sender_email_identity": 2.718281828459045, "default/0/contact/email/email_template": 4711, "default/0/contact/contact/enabled": true, "default/0/contact/email/recipient_email": "hello@example.com"},
		pkgCfg.Defaults(),
	)
}

func TestSectionSliceMerge(t *testing.T) {

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
						ID: "a",
						Groups: element.NewGroupSlice(
							nil,
							&element.Group{
								ID: "b",
								Fields: element.NewFieldSlice(
									&element.Field{ID: "c", Default: `c`},
								),
							},
							&element.Group{
								ID: "b",
								Fields: element.NewFieldSlice(
									&element.Field{ID: "d", Default: `d`},
								),
							},
						),
					},
				},
				{
					&element.Section{ID: "a", Label: "LabelA", Groups: nil},
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
						ID:    "a",
						Label: "SectionLabelA",
						Groups: element.NewGroupSlice(
							&element.Group{
								ID:    "b",
								Scope: scope.NewPerm(scope.DefaultID),
								Fields: element.NewFieldSlice(
									&element.Field{ID: "c", Default: `c`},
								),
							},
							nil,
						),
					},
				},
				{
					&element.Section{
						ID:    "a",
						Scope: scope.NewPerm(scope.DefaultID, scope.WebsiteID),
						Groups: element.NewGroupSlice(
							&element.Group{ID: "b", Label: "GroupLabelB1"},
							nil,
							&element.Group{ID: "b", Label: "GroupLabelB2"},
							&element.Group{
								ID: "b2",
								Fields: element.NewFieldSlice(
									&element.Field{ID: "d", Default: `d`},
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
					&element.Section{ID: "a", Label: "SectionLabelA", SortOrder: 20, Resource: 22},
				},
				{
					&element.Section{ID: "a", Scope: scope.NewPerm(scope.DefaultID, scope.WebsiteID), SortOrder: 10, Resource: 3},
				},
			},
			wantErr: "",
			want:    `[{"ID":"a","Label":"SectionLabelA","Scope":["Default","Website"],"SortOrder":10,"Resource":3,"Groups":null}]` + "\n",
		},
		3: {
			have: []element.SectionSlice{
				{
					&element.Section{
						ID:    "a",
						Label: "SectionLabelA",
						Groups: element.NewGroupSlice(
							&element.Group{
								ID:      "b",
								Label:   "SectionAGroupB",
								Comment: element.LongText("SectionAGroupBComment"),
								Scope:   scope.NewPerm(scope.DefaultID),
							},
						),
					},
				},
				{
					&element.Section{
						ID:        "a",
						SortOrder: 1000,
						Scope:     scope.NewPerm(scope.DefaultID, scope.WebsiteID),
						Groups: element.NewGroupSlice(
							&element.Group{ID: "b", Label: "GroupLabelB1", Scope: scope.PermAll},
							&element.Group{ID: "b", Label: "GroupLabelB2", Comment: element.LongText("Section2AGroup3BComment"), SortOrder: 100},
							&element.Group{ID: "b2"},
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
						ID: "a",
						Groups: element.NewGroupSlice(
							&element.Group{
								ID:    "b",
								Label: "b1",
								Fields: element.NewFieldSlice(
									&element.Field{ID: "c", Default: `c`, Type: element.TypeMultiselect, SortOrder: 1001},
								),
							},
							&element.Group{
								ID:    "b",
								Label: "b2",
								Fields: element.NewFieldSlice(
									nil,
									&element.Field{ID: "d", Default: `d`, Comment: element.LongText("Ring of fire"), Type: element.TypeObscure},
									&element.Field{ID: "c", Default: `haha`, Type: element.TypeSelect, Scope: scope.NewPerm(scope.DefaultID, scope.WebsiteID)},
								),
							},
						),
					},
				},
				{
					&element.Section{
						ID: "a",
						Groups: element.NewGroupSlice(
							&element.Group{
								ID:    "b",
								Label: "b3",
								Fields: element.NewFieldSlice(
									&element.Field{ID: "d", Default: `overriddenD`, Label: "Sect2Group2Label4", Comment: element.LongText("LOTR")},
									&element.Field{ID: "c", Default: `overriddenHaha`, Type: element.TypeHidden},
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
						ID: "a",
						Groups: element.NewGroupSlice(
							&element.Group{
								ID: "b",
								Fields: element.NewFieldSlice(
									&element.Field{
										ID:      "c",
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
						ID: "a",
						Groups: element.NewGroupSlice(
							&element.Group{
								ID: "b",
								Fields: element.NewFieldSlice(
									nil,
									&element.Field{
										ID:        "c",
										Default:   `overridenC`,
										Type:      element.TypeSelect,
										Label:     "Sect2Group2Label4",
										Comment:   element.LongText("LOTR"),
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

	tests := []struct {
		have    []*element.Group
		wantErr error
		want    string
	}{
		{
			have: []*element.Group{
				{
					ID: "b",
					Fields: element.NewFieldSlice(
						&element.Field{ID: "c", Default: `c`, Type: element.TypeMultiselect},
					),
				},
				{
					ID: "b",
					Fields: element.NewFieldSlice(
						&element.Field{ID: "d", Default: `d`, Comment: element.LongText("Ring of fire"), Type: element.TypeObscure},
						&element.Field{ID: "c", Default: `haha`, Type: element.TypeSelect, Scope: scope.NewPerm(scope.DefaultID, scope.WebsiteID)},
					),
				},
				{
					ID: "b",
					Fields: element.NewFieldSlice(
						&element.Field{ID: "d", Default: `overriddenD`, Label: "Sect2Group2Label4", Comment: element.LongText("LOTR")},
						&element.Field{ID: "c", Default: `overriddenHaha`, Type: element.TypeHidden},
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

func TestSectionSliceFindGroupByPath(t *testing.T) {
	tests := []struct {
		haveSlice element.SectionSlice
		havePath  []string
		wantGID   string
		wantErr   error
	}{
		0: {
			haveSlice: element.SectionSlice{&element.Section{ID: "a", Groups: element.NewGroupSlice(&element.Group{ID: "b"}, &element.Group{ID: "bb"})}},
			havePath:  []string{"a/b"},
			wantGID:   "b",
			wantErr:   nil,
		},
		//1: {
		//	haveSlice: element.SectionSlice{&element.Section{ID: "a", Groups: element.NewGroupSlice(&element.Group{ID: "b"}, &element.Group{ID: "bb"})}},
		//	havePath:  []string{"a/bc"},
		//	wantGID:   "b",
		//	wantErr:   element.ErrGroupNotFound,
		//},
		//2: {
		//	haveSlice: element.SectionSlice{},
		//	havePath:  nil,
		//	wantGID:   "b",
		//	wantErr:   element.ErrGroupNotFound,
		//},
		//3: {
		//	haveSlice: element.SectionSlice{&element.Section{ID: "a", Groups: element.GroupSlice{&element.Group{ID: "b"}, &element.Group{ID: "bb"}}}},
		//	havePath:  []string{"a", "bb", "cc"},
		//	wantGID:   "bb",
		//	wantErr:   nil,
		//},
		//4: {
		//	haveSlice: element.SectionSlice{&element.Section{ID: "a", Groups: element.GroupSlice{&element.Group{ID: "b"}, &element.Group{ID: "bb"}}}},
		//	havePath:  []string{"xa", "bb", "cc"},
		//	wantGID:   "",
		//	wantErr:   element.ErrSectionNotFound,
		//},
	}

	for i, test := range tests {
		haveGroup, haveErr := test.haveSlice.FindGroupByPath(test.havePath...)
		if test.wantErr != nil {
			assert.Error(t, haveErr, "Index %d", i)
			assert.Nil(t, haveGroup)
			assert.EqualError(t, haveErr, test.wantErr.Error())
		} else {
			assert.NoError(t, haveErr, "Index %d", i)
			assert.NotNil(t, haveGroup, "Index %d", i)
			assert.Exactly(t, test.wantGID, haveGroup.ID)
		}
	}
}

func TestSectionSliceFindFieldByPath(t *testing.T) {

	tests := []struct {
		haveSlice element.SectionSlice
		havePath  []string
		wantFID   string
		wantErr   error
	}{
		0: {
			haveSlice: element.SectionSlice{&element.Section{ID: "a", Groups: element.GroupSlice{&element.Group{ID: "b"}, &element.Group{ID: "bb"}}}},
			havePath:  []string{"a/b"},
			wantFID:   "b",
			wantErr:   element.ErrFieldNotFound,
		},
		1: {
			haveSlice: element.SectionSlice{&element.Section{ID: "a", Groups: element.GroupSlice{&element.Group{ID: "b"}, &element.Group{ID: "bb"}}}},
			havePath:  []string{"a/bc"},
			wantFID:   "b",
			wantErr:   element.ErrFieldNotFound,
		},
		2: {
			haveSlice: element.SectionSlice{nil},
			havePath:  nil,
			wantFID:   "b",
			wantErr:   element.ErrFieldNotFound,
		},
		3: {
			haveSlice: element.SectionSlice{&element.Section{ID: "a", Groups: element.GroupSlice{nil, &element.Group{ID: "b"}, &element.Group{ID: "bb"}}}},
			havePath:  []string{"a", "bb", "cc"},
			wantFID:   "bb",
			wantErr:   element.ErrFieldNotFound,
		},
		4: {
			haveSlice: element.SectionSlice{&element.Section{ID: "a", Groups: element.GroupSlice{&element.Group{ID: "b"}, &element.Group{ID: "bb"}}}},
			havePath:  []string{"xa", "bb", "cc"},
			wantFID:   "",
			wantErr:   element.ErrSectionNotFound,
		},
		5: {
			haveSlice: element.SectionSlice{&element.Section{ID: "a", Groups: element.GroupSlice{&element.Group{ID: "b", Fields: element.FieldSlice{
				&element.Field{ID: "c"}, nil,
			}}}}},
			havePath: []string{"a", "b", "c"},
			wantFID:  "c",
			wantErr:  nil,
		},
	}

	for i, test := range tests {
		haveGroup, haveErr := test.haveSlice.FindFieldByPath(test.havePath...)
		if test.wantErr != nil {
			assert.Error(t, haveErr, "Index %d", i)
			assert.Nil(t, haveGroup, "Index %d", i)
			assert.EqualError(t, haveErr, test.wantErr.Error(), "Index %d", i)
		} else {
			assert.NoError(t, haveErr, "Index %d", i)
			assert.NotNil(t, haveGroup, "Index %d", i)
			assert.Exactly(t, test.wantFID, haveGroup.ID, "Index %d", i)
		}
	}
}

func TestFieldSliceSort(t *testing.T) {
	want := []int{-10, 1, 10, 11, 20}
	fs := element.FieldSlice{
		&element.Field{ID: "a", SortOrder: 20},
		&element.Field{ID: "b", SortOrder: -10},
		&element.Field{ID: "c", SortOrder: 10},
		&element.Field{ID: "d", SortOrder: 11},
		&element.Field{ID: "e", SortOrder: 1},
	}

	for i, f := range *(fs.Sort()) {
		assert.EqualValues(t, want[i], f.SortOrder)
	}
}

func TestGroupSliceSort(t *testing.T) {
	want := []int{-10, 1, 10, 11, 20}
	gs := element.GroupSlice{
		&element.Group{ID: "a", SortOrder: 20},
		&element.Group{ID: "b", SortOrder: -10},
		&element.Group{ID: "c", SortOrder: 10},
		&element.Group{ID: "d", SortOrder: 11},
		&element.Group{ID: "e", SortOrder: 1},
	}
	for i, f := range *(gs.Sort()) {
		assert.EqualValues(t, want[i], f.SortOrder)
	}
}
func TestSectionSliceSort(t *testing.T) {
	want := []int{-10, 1, 10, 11, 20}
	ss := element.SectionSlice{
		&element.Section{ID: "a", SortOrder: 20},
		&element.Section{ID: "b", SortOrder: -10},
		&element.Section{ID: "c", SortOrder: 10},
		&element.Section{ID: "d", SortOrder: 11},
		&element.Section{ID: "e", SortOrder: 1},
	}
	for i, f := range *(ss.Sort()) {
		assert.EqualValues(t, want[i], f.SortOrder)
	}

}

func TestSectionSliceSortAll(t *testing.T) {
	want := `[{"ID":"b","SortOrder":-10,"Groups":null},{"ID":"e","SortOrder":1,"Groups":null},{"ID":"c","SortOrder":10,"Groups":null},{"ID":"a","SortOrder":20,"Groups":[{"ID":"b","SortOrder":-10,"Fields":[{"ID":"b","SortOrder":-10},{"ID":"e","SortOrder":1},{"ID":"c","SortOrder":10},{"ID":"d","SortOrder":11},{"ID":"a","SortOrder":20}]},{"ID":"e","SortOrder":1,"Fields":null},{"ID":"d","SortOrder":11,"Fields":[{"ID":"b","SortOrder":-10},{"ID":"e","SortOrder":1},{"ID":"c","SortOrder":10},{"ID":"d","SortOrder":11},{"ID":"a","SortOrder":20}]},{"ID":"a","SortOrder":20,"Fields":[{"ID":"b","SortOrder":-10},{"ID":"e","SortOrder":1},{"ID":"c","SortOrder":10},{"ID":"d","SortOrder":11},{"ID":"a","SortOrder":20}]}]}]` + "\n"
	ss := element.MustNewConfiguration(
		&element.Section{ID: "a", SortOrder: 20, Groups: element.NewGroupSlice(
			&element.Group{ID: "a", SortOrder: 20, Fields: element.FieldSlice{&element.Field{ID: "a", SortOrder: 20}, &element.Field{ID: "b", SortOrder: -10}, &element.Field{ID: "c", SortOrder: 10}, &element.Field{ID: "d", SortOrder: 11}, &element.Field{ID: "e", SortOrder: 1}}},
			&element.Group{ID: "b", SortOrder: -10, Fields: element.FieldSlice{&element.Field{ID: "a", SortOrder: 20}, &element.Field{ID: "b", SortOrder: -10}, &element.Field{ID: "c", SortOrder: 10}, &element.Field{ID: "d", SortOrder: 11}, &element.Field{ID: "e", SortOrder: 1}}},
			&element.Group{ID: "d", SortOrder: 11, Fields: element.FieldSlice{&element.Field{ID: "a", SortOrder: 20}, &element.Field{ID: "b", SortOrder: -10}, &element.Field{ID: "c", SortOrder: 10}, &element.Field{ID: "d", SortOrder: 11}, &element.Field{ID: "e", SortOrder: 1}}},
			&element.Group{ID: "e", SortOrder: 1},
		)},
		&element.Section{ID: "b", SortOrder: -10},
		&element.Section{ID: "c", SortOrder: 10},
		&element.Section{ID: "e", SortOrder: 1},
	)
	ss.SortAll()
	have := ss.ToJSON()
	if want != have {
		t.Errorf("\nWant: %s\nHave: %s\n", want, have)
	}
}

func TestSectionSliceAppendFields(t *testing.T) {
	want := `[{"ID":"a","Groups":[{"ID":"a","Fields":[{"ID":"a"},{"ID":"b"},{"ID":"c"}]}]}]` + "\n"
	ss := element.MustNewConfiguration(
		&element.Section{
			ID: "a",
			Groups: element.NewGroupSlice(
				&element.Group{ID: "a",
					Fields: element.NewFieldSlice(
						&element.Field{ID: "a"},
						&element.Field{ID: "b"},
					),
				},
			)},
	)
	assert.EqualError(t, ss.AppendFields("a/XX"), element.ErrGroupNotFound.Error())

	assert.NoError(t, ss.AppendFields("a/a", &element.Field{ID: "c"}))
	have := ss.ToJSON()
	if want != have {
		t.Errorf("\nWant: %s\nHave: %s\n", want, have)
	}
}
