// Copyright 2015-present, Cyrill @ Schumacher.fm and the CoreStore contributors
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

package config_test

import (
	"encoding/json"
	"testing"

	"github.com/corestoreio/errors"
	"github.com/corestoreio/pkg/config"
	"github.com/corestoreio/pkg/store/scope"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSections_Validate(t *testing.T) {
	t.Parallel()
	t.Run("Duplicate", func(t *testing.T) {
		ss := config.MakeSections(
			&config.Section{
				ID: `aa`,
				Groups: config.MakeGroups(
					&config.Group{
						ID: `bb`,
						Fields: config.MakeFields(
							&config.Field{ID: `cc`},
							&config.Field{ID: `cc`},
						),
					},
				),
			},
		)
		err := ss.Validate()
		assert.True(t, errors.Duplicated.Match(err), "%+v", err)
	})
	t.Run("malformed", func(t *testing.T) {
		ss := config.MakeSections(
			&config.Section{
				ID: `aa`,
				Groups: config.MakeGroups(
					&config.Group{
						ID: `bb`,
						Fields: config.MakeFields(
							&config.Field{ID: `c`},
						),
					},
				),
			},
		)
		err := ss.Validate()
		assert.True(t, errors.NotValid.Match(err), "%+v", err)
	})
	t.Run("length zero", func(t *testing.T) {
		ss := config.MakeSections()
		err := ss.Validate()
		assert.NoError(t, err)
	})
}

func TestSectionValidateShortPath(t *testing.T) {
	t.Parallel()

	ss := config.MakeSections(
		&config.Section{
			//ID: `aa`,
			Groups: config.MakeGroups(
				&config.Group{
					//ID: `b`,
					Fields: config.MakeFields(
						&config.Field{ID: `ca`},
						&config.Field{ID: `cb`},
						&config.Field{},
					),
				},
			),
		},
	)

	err := ss.Validate()
	assert.True(t, errors.NotValid.Match(err), "%+v", err)
}

func TestSectionUpdateField(t *testing.T) {

	ss := config.MakeSections(
		&config.Section{
			ID: `aa`,
			Groups: config.MakeGroups(
				&config.Group{
					ID: `bb`,
					Fields: config.MakeFields(
						&config.Field{ID: `ca`},
						&config.Field{ID: `cb`},
					),
				},
			),
		},
	)

	fr := `aa/bb/ca`
	if idx := ss.UpdateField(fr, &config.Field{
		Label: "ca New Label",
	}); idx < 0 {
		t.Fatal("Field not found")
	}

	f, idx := ss.FindField(fr)
	if idx < 0 {
		t.Fatalf("Not found %q", fr)
	}
	assert.Exactly(t, `ca New Label`, f.Label)

	idx = ss.UpdateField(`a/b/c`, &config.Field{})
	assert.Exactly(t, -1, idx, "Field not found")

	idx = ss.UpdateField(`aa/b/c`, &config.Field{})
	assert.Exactly(t, -2, idx, "Group not found")

	idx = ss.UpdateField(`aa/bb/c`, &config.Field{})
	assert.Exactly(t, -3, idx, "Field not found")

	idx = ss.UpdateField(`aa_bb_c`, &config.Field{})
	assert.Exactly(t, -10, idx, "Invalid path")
}

func TestNewConfiguration(t *testing.T) {

	tests := []struct {
		have       config.Sections
		wantErrBhf errors.Kind
		wantLen    int
	}{
		{
			have: config.MakeSections(
				&config.Section{
					ID: `web`,
					Groups: config.MakeGroups(
						&config.Group{
							ID:     `default`,
							Fields: config.MakeFields(&config.Field{ID: `front`}, &config.Field{ID: `no_route`}),
						},
					),
				},
				&config.Section{
					ID: `system`,
					Groups: config.MakeGroups(
						&config.Group{
							ID:     `media_storage_configuration`,
							Fields: config.MakeFields(&config.Field{ID: `allowed_resources`}),
						},
					),
				},
			),
			wantErrBhf: 0,
			wantLen:    3,
		},
		{
			have:       config.MakeSections(&config.Section{ID: `aa`, Groups: config.MakeGroups()}),
			wantErrBhf: 0,
		},
		{
			have:       config.MakeSections(&config.Section{ID: `aa`, Groups: config.MakeGroups(&config.Group{ID: `bb`, Fields: nil})}),
			wantErrBhf: 0,
		},
		{
			have: config.MakeSections(
				&config.Section{
					ID: `aa`,
					Groups: config.MakeGroups(
						&config.Group{
							ID:     `bb`,
							Fields: config.MakeFields(&config.Field{ID: `cc`}, &config.Field{ID: `cc`}),
						},
					),
				},
			),
			wantErrBhf: errors.Duplicated,
		},
	}

	for i, test := range tests {

		haveSlice, err := config.MakeSectionsValidated(test.have...)
		if test.wantErrBhf > 0 {
			assert.Nil(t, haveSlice, "Index %d", i)
			assert.True(t, test.wantErrBhf.Match(err), "IDX %d: %+v", i, err)
		} else {
			assert.NotNil(t, haveSlice, "Index %d", i)
			assert.Len(t, haveSlice, len(test.have), "Index %d", i)
		}
		assert.Exactly(t, test.wantLen, haveSlice.TotalFields(), "Index %d", i)
	}
}

func TestSectionSliceMerge(t *testing.T) {

	// Got stuck in comparing JSON?
	// Use a Webservice to compare the JSON output!

	tests := []struct {
		have       []config.Sections
		want       string
		fieldCount int
	}{
		{
			have: []config.Sections{
				{
					&config.Section{
						ID: `a`,
					},
				},
				{
					&config.Section{ID: `a`, Label: `LabelA`, Groups: nil},
				},
			},
			want:       `[{"ID":"a","Label":"LabelA"}]`,
			fieldCount: 0,
		},
		{
			have: []config.Sections{
				{
					&config.Section{
						ID: `a`,
						Groups: config.MakeGroups(
							&config.Group{
								ID: `b`,
								Fields: config.MakeFields(
									&config.Field{ID: `c`, Default: `c`},
								),
							},
							&config.Group{
								ID: `b`,
								Fields: config.MakeFields(
									&config.Field{ID: `d`, Default: `d`},
								),
							},
						),
					},
				},
				{
					&config.Section{ID: `a`, Label: `LabelA`, Groups: nil},
				},
			},
			want:       `[{"ID":"a","Label":"LabelA","Groups":[{"ID":"b","Fields":[{"ID":"c","Default":"c"},{"ID":"d","Default":"d"}]}]}]`,
			fieldCount: 2,
		},
		{
			have: []config.Sections{
				{
					&config.Section{
						ID:    `a`,
						Label: `SectionLabelA`,
						Groups: config.MakeGroups(
							&config.Group{
								ID:     `b`,
								Scopes: scope.PermDefault,
								Fields: config.MakeFields(
									&config.Field{ID: `c`, Default: `c`},
								),
							},
						),
					},
				},
				{
					&config.Section{
						ID:     `a`,
						Scopes: scope.PermWebsite,
						Groups: config.MakeGroups(
							&config.Group{ID: `b`, Label: `GroupLabelB1`},
							&config.Group{ID: `b`, Label: `GroupLabelB2`},
							&config.Group{
								ID: `b2`,
								Fields: config.MakeFields(
									&config.Field{ID: `d`, Default: `d`},
								),
							},
						),
					},
				},
			},
			want:       `[{"ID":"a","Label":"SectionLabelA","Scopes":"websites","Groups":[{"ID":"b","Label":"GroupLabelB2","Scopes":"default","Fields":[{"ID":"c","Default":"c"}]},{"ID":"b2","Fields":[{"ID":"d","Default":"d"}]}]}]`,
			fieldCount: 2,
		},
		{
			have: []config.Sections{
				{
					&config.Section{ID: `a`, Label: `SectionLabelA`, SortOrder: 20, Resource: 22},
				},
				{
					&config.Section{ID: `a`, Scopes: scope.PermWebsite, SortOrder: 10, Resource: 3},
				},
			},
			want: `[{"ID":"a","Label":"SectionLabelA","Scopes":"websites","SortOrder":10,"Resource":3}]`,
		},
		{
			have: []config.Sections{
				{
					&config.Section{
						ID:    `a`,
						Label: `SectionLabelA`,
						Groups: config.MakeGroups(
							&config.Group{
								ID:      `b`,
								Label:   `SectionAGroupB`,
								Comment: "SectionAGroupBComment",
								Scopes:  scope.PermDefault,
							},
						),
					},
				},
				{
					&config.Section{
						ID:        `a`,
						SortOrder: 1000,
						Scopes:    scope.PermWebsite,
						Groups: config.MakeGroups(
							&config.Group{ID: `b`, Label: `GroupLabelB1`, Scopes: scope.PermStore},
							&config.Group{ID: `b`, Label: `GroupLabelB2`, Comment: "Section2AGroup3BComment", SortOrder: 100},
							&config.Group{ID: `b2`},
						),
					},
				},
			},
			want: `[{"ID":"a","Label":"SectionLabelA","Scopes":"websites","SortOrder":1000,"Groups":[{"ID":"b","Label":"GroupLabelB2","Comment":"Section2AGroup3BComment","Scopes":"stores","SortOrder":100},{"ID":"b2"}]}]`,
		},
		{
			have: []config.Sections{
				{
					&config.Section{
						ID: `a`,
						Groups: config.MakeGroups(
							&config.Group{
								ID:    `b`,
								Label: `b1`,
								Fields: config.MakeFields(
									&config.Field{ID: `c`, Default: `c`, Type: config.TypeMultiselect, SortOrder: 1001},
								),
							},
							&config.Group{
								ID:    `b`,
								Label: `b2`,
								Fields: config.MakeFields(
									&config.Field{ID: `d`, Default: `d`, Comment: "Ring of fire", Type: config.TypeObscure},
									&config.Field{ID: `c`, Default: `haha`, Type: config.TypeSelect, Scopes: scope.PermWebsite},
								),
							},
						),
					},
				},
				{
					&config.Section{
						ID: `a`,
						Groups: config.MakeGroups(
							&config.Group{
								ID:    `b`,
								Label: `b3`,
								Fields: config.MakeFields(
									&config.Field{ID: `d`, Default: `overriddenD`, Label: `Sect2Group2Label4`, Comment: "LOTR"},
									&config.Field{ID: `c`, Default: `overriddenHaha`, Type: config.TypeHidden},
								),
							},
						),
					},
				},
			},
			want:       `[{"ID":"a","Groups":[{"ID":"b","Label":"b3","Fields":[{"ID":"c","Type":"hidden","SortOrder":1001,"Scopes":"websites","Default":"overriddenHaha"},{"ID":"d","Type":"obscure","Label":"Sect2Group2Label4","Comment":"LOTR","Default":"overriddenD"}]}]}]`,
			fieldCount: 2,
		},
		{
			have: []config.Sections{
				{
					&config.Section{
						ID: `a`,
						Groups: config.MakeGroups(
							&config.Group{
								ID: `b`,
								Fields: config.MakeFields(
									&config.Field{
										ID:      `c`,
										Default: `c`,
										Type:    config.TypeMultiselect,
									},
								),
							},
						),
					},
				},
				{
					&config.Section{
						ID: `a`,
						Groups: config.MakeGroups(
							&config.Group{
								ID: `b`,
								Fields: config.MakeFields(
									&config.Field{
										ID:        `c`,
										Default:   `overridenC`,
										Type:      config.TypeSelect,
										Label:     `Sect2Group2Label4`,
										Comment:   "LOTR",
										SortOrder: 100,
										Visible:   true,
									},
								),
							},
						),
					},
				},
			},
			fieldCount: 1,
			want:       `[{"ID":"a","Groups":[{"ID":"b","Fields":[{"ID":"c","Type":"select","Label":"Sect2Group2Label4","Comment":"LOTR","SortOrder":100,"Visible":true,"Default":"overridenC"}]}]}]`,
		},
	}

	for i, test := range tests {

		if len(test.have) == 0 {
			test.want = "null\n"
		}

		var baseSl config.Sections
		baseSl = baseSl.MergeMultiple(test.have...)
		j, err := json.Marshal(baseSl)
		require.NoError(t, err)

		if string(j) != test.want {
			t.Errorf("\nIndex: %d\nExpected: %s\nActual:   %s\n", i, test.want, j)
		}

		assert.Exactly(t, test.fieldCount, baseSl.TotalFields(), "Index %d", i)
	}
}

func TestFieldSliceMerge(t *testing.T) {

	tests := []struct {
		have config.Fields
		want string
	}{
		{
			have: config.MakeFields(
				&config.Field{ID: `d`, Default: `overrideMeD`, Comment: "Ring of fire", Type: config.TypeObscure},
			),
			want: `[{"ID":"d","Type":"obscure","Comment":"Ring of fire","Default":"overrideMeD"}]`,
		},
		{
			have: config.MakeFields(
				&config.Field{ID: `d`, Default: `overrideMeD`, Comment: "Ring of fire", Type: config.TypeObscure},
				&config.Field{ID: `c`, Default: `overrideMeC`, Type: config.TypeSelect, Scopes: scope.PermWebsite},
			),
			want: `[{"ID":"d","Type":"obscure","Comment":"Ring of fire","Default":"overrideMeD"},{"ID":"c","Type":"select","Scopes":"websites","Default":"overrideMeC"}]`,
		},
		{
			have: config.MakeFields(
				&config.Field{ID: `d`, Default: `overrideMeD`, Comment: "Ring of fire", Type: config.TypeObscure},
				&config.Field{ID: `c`, Default: `overrideMeC`, Type: config.TypeSelect, Scopes: scope.PermWebsite},
				&config.Field{ID: `d`, Default: `overrideMeE`, Type: config.TypeMultiselect},
			),
			want: `[{"ID":"d","Type":"multiselect","Comment":"Ring of fire","Default":"overrideMeE"},{"ID":"c","Type":"select","Scopes":"websites","Default":"overrideMeC"}]`,
		},
		{
			have: nil,
			want: `null`,
		},
	}

	for i, test := range tests {
		var baseFsl config.Fields
		baseFsl = baseFsl.Merge(test.have...)

		fsj, err := json.Marshal(baseFsl)
		if err != nil {
			t.Fatal(err)
		}
		if string(fsj) != test.want {
			t.Errorf("\nIndex: %d\nExpected: %s\nActual:   %s\n", i, test.want, string(fsj))
		}

	}
}

func TestGroupSliceMerge(t *testing.T) {

	tests := []struct {
		have config.Groups
		want string
	}{
		{
			have: config.MakeGroups(
				&config.Group{
					ID: `b`,
					Fields: config.MakeFields(
						&config.Field{ID: `c`, Default: `c`, Type: config.TypeMultiselect},
					),
				},
				&config.Group{
					ID: `b`,
					Fields: config.MakeFields(
						&config.Field{ID: `d`, Default: `overrideMeD`, Comment: "Ring of fire", Type: config.TypeObscure},
						&config.Field{ID: `c`, Default: `overrideMeC`, Type: config.TypeSelect, Scopes: scope.PermWebsite},
					),
				},
				&config.Group{
					ID: `b`,
					Fields: config.MakeFields(
						&config.Field{ID: `d`, Default: `overriddenD`, Label: `Sect2Group2Label4`, Comment: "LOTR"},
						&config.Field{ID: `c`, Default: `overriddenC`, Type: config.TypeHidden},
					),
				},
			),
			want: `[{"ID":"b","Fields":[{"ID":"c","Type":"hidden","Scopes":"websites","Default":"overriddenC"},{"ID":"d","Type":"obscure","Label":"Sect2Group2Label4","Comment":"LOTR","Default":"overriddenD"}]}]`,
		},
		{
			have: config.MakeGroups(
				&config.Group{
					ID:    `b`,
					Label: `Single Field`,
					Fields: config.MakeFields(
						&config.Field{ID: `c`, Default: `c`, Type: config.TypeMultiselect},
					),
				},
			),
			want: `[{"ID":"b","Label":"Single Field","Fields":[{"ID":"c","Type":"multiselect","Default":"c"}]}]`,
		},
		{
			have: config.MakeGroups(
				&config.Group{
					ID:    `b`,
					Label: `Single Field`,
					Fields: config.MakeFields(
						&config.Field{ID: `c`, Default: `c`, Type: config.TypeMultiselect},
					),
				},
				&config.Group{
					ID:    `b`,
					Label: `Single Field2`,
					Fields: config.MakeFields(
						&config.Field{ID: `c`, Default: `c`, Type: config.TypeMultiselect},
					),
				},
			),
			want: `[{"ID":"b","Label":"Single Field2","Fields":[{"ID":"c","Type":"multiselect","Default":"c"}]}]`,
		},
		{
			have: config.MakeGroups(
				&config.Group{
					ID:    `b`,
					Label: `Single Field`,
					Fields: config.MakeFields(
						&config.Field{ID: `c`, Default: `c`, Type: config.TypeMultiselect},
					),
				},
				&config.Group{
					ID:    `b`,
					Label: `Single Field2`,
					Fields: config.MakeFields(
						&config.Field{ID: `c`, Default: `c2`, Type: config.TypeTextarea},
					),
				},
			),
			want: `[{"ID":"b","Label":"Single Field2","Fields":[{"ID":"c","Type":"textarea","Default":"c2"}]}]`,
		},
		{
			have: config.MakeGroups(
				&config.Group{
					ID:    `b`,
					Label: `Single Field`,
					Fields: config.MakeFields(
						&config.Field{ID: `c`, Default: `c`, Type: config.TypeMultiselect},
					),
				},
				&config.Group{
					ID: `b`,
					Fields: config.MakeFields(
						&config.Field{ID: `d`, Default: `d`, Type: config.TypeText},
					),
				},
			),
			want: `[{"ID":"b","Label":"Single Field","Fields":[{"ID":"c","Type":"multiselect","Default":"c"},{"ID":"d","Type":"text","Default":"d"}]}]`,
		},
		{
			have: nil,
			want: `null`,
		},
	}

	for i, test := range tests {
		var baseGsl config.Groups
		baseGsl = baseGsl.Merge(test.have...)

		j, err := json.Marshal(baseGsl)
		require.NoError(t, err)
		if string(j) != test.want {
			t.Errorf("\nIndex: %d\nExpected: %s\nActual:   %s\n", i, test.want, j)
		}

	}
}

func TestSectionSliceFindGroupByID(t *testing.T) {
	t.Parallel()
	tests := []struct {
		haveSlice config.Sections
		haveRoute string
		wantGID   string
		wantIdx   int
	}{
		{
			haveSlice: config.MakeSections(&config.Section{ID: `a`, Groups: config.MakeGroups(&config.Group{ID: `b`}, &config.Group{ID: `bb`})}),
			haveRoute: "a/b",
			wantGID:   "b",
			wantIdx:   -10,
		},
		{
			haveSlice: config.MakeSections(&config.Section{ID: `a`, Groups: config.MakeGroups(&config.Group{ID: `b`}, &config.Group{ID: `bb`})}),
			haveRoute: "a/b/c",
			wantGID:   "b",
		},
		{
			haveSlice: config.MakeSections(&config.Section{ID: `a`, Groups: config.MakeGroups(&config.Group{ID: `b`}, &config.Group{ID: `bb`})}),
			haveRoute: "a/bc/d",
			wantGID:   "b",
			wantIdx:   -2,
		},
		{
			haveSlice: config.Sections{},
			haveRoute: "",
			wantGID:   "b",
			wantIdx:   -10,
		},
		{
			haveSlice: config.MakeSections(&config.Section{ID: `a`, Groups: config.MakeGroups(&config.Group{ID: `b`}, &config.Group{ID: `bb`})}),
			haveRoute: "a/bb/cc",
			wantGID:   "bb",
		},
		{
			haveSlice: config.MakeSections(&config.Section{ID: `a`, Groups: config.MakeGroups(&config.Group{ID: `b`}, &config.Group{ID: `bb`})}),
			haveRoute: "xa/bb/cc",
			wantGID:   "",
			wantIdx:   -1,
		},
	}

	for i, test := range tests {
		haveGroup, haveIdx := test.haveSlice.FindGroup(test.haveRoute)
		if test.wantIdx < 0 {
			assert.Exactly(t, test.wantIdx, haveIdx, "Index %d", i)
			assert.Nil(t, haveGroup, "Index %d", i)
			continue
		}

		assert.True(t, haveIdx >= 0, "Index %d", i)
		require.NotNil(t, haveGroup, "Index %d", i)
		assert.Exactly(t, test.wantGID, haveGroup.ID)
	}
}

func TestSectionSliceFindFieldByID(t *testing.T) {

	tests := []struct {
		haveSlice config.Sections
		haveRoute string
		wantFID   string
		wantIdx   int
	}{
		{
			haveSlice: config.MakeSections(&config.Section{ID: `aa`, Groups: config.MakeGroups(&config.Group{ID: `bb`}, &config.Group{ID: `cc`})}),
			haveRoute: "",
			wantFID:   "",
			wantIdx:   -10,
		},
		{
			haveSlice: config.MakeSections(&config.Section{ID: `a`, Groups: config.MakeGroups(&config.Group{ID: `b`}, &config.Group{ID: `bb`})}),
			haveRoute: "a/b",
			wantFID:   "b",
			wantIdx:   -10,
		},
		{
			haveSlice: config.MakeSections(&config.Section{ID: `a`, Groups: config.MakeGroups(&config.Group{ID: `b`}, &config.Group{ID: `bb`})}),
			haveRoute: "a/bc",
			wantFID:   "b",
			wantIdx:   -10,
		},
		{
			haveSlice: config.MakeSections(),
			haveRoute: "",
			wantFID:   "",
			wantIdx:   -10,
		},
		{
			haveSlice: config.MakeSections(&config.Section{ID: `a`, Groups: config.MakeGroups(&config.Group{ID: `b`}, &config.Group{ID: `bb`})}),
			haveRoute: "a/bb/cc",
			wantFID:   "bb",
			wantIdx:   -3,
		},
		{
			haveSlice: config.MakeSections(&config.Section{ID: `a`, Groups: config.MakeGroups(&config.Group{ID: `b`}, &config.Group{ID: `bb`})}),
			haveRoute: "xa/bb/cc",
			wantFID:   "",
			wantIdx:   -1,
		},
		{
			haveSlice: config.MakeSections(&config.Section{ID: `a1`, Groups: config.MakeGroups(&config.Group{ID: `b1`, Fields: config.MakeFields(&config.Field{ID: `c1`})})}),
			haveRoute: "a1/b1/c1",
			wantFID:   "c1",
		},
	}

	for i, test := range tests {
		haveGroup, haveIdx := test.haveSlice.FindField(test.haveRoute)
		if test.wantIdx < 0 {
			assert.Exactly(t, test.wantIdx, haveIdx, "Index %d", i)
			assert.Nil(t, haveGroup, "Index %d", i)
			continue
		}
		assert.True(t, haveIdx >= 0, "Index %d", i)
		require.NotNil(t, haveGroup)
		assert.Exactly(t, test.wantFID, haveGroup.ID, "Index %d", i)
	}
}

func TestFieldSliceSort(t *testing.T) {

	want := []int{-10, 1, 10, 11, 20}
	fs := config.MakeFields(
		&config.Field{ID: `aa`, SortOrder: 20},
		&config.Field{ID: `bb`, SortOrder: -10},
		&config.Field{ID: `cc`, SortOrder: 10},
		&config.Field{ID: `dd`, SortOrder: 11},
		&config.Field{ID: `ee`, SortOrder: 1},
	)

	for i, f := range fs.Sort() {
		assert.EqualValues(t, want[i], f.SortOrder)
	}
}

func TestGroupSliceSort(t *testing.T) {

	want := []int{-10, 1, 10, 11, 20}
	gs := config.MakeGroups(
		&config.Group{ID: `aa`, SortOrder: 20},
		&config.Group{ID: `bb`, SortOrder: -10},
		&config.Group{ID: `cc`, SortOrder: 10},
		&config.Group{ID: `dd`, SortOrder: 11},
		&config.Group{ID: `ee`, SortOrder: 1},
	)
	for i, f := range gs.Sort() {
		assert.EqualValues(t, want[i], f.SortOrder)
	}
}
func TestSectionSliceSort(t *testing.T) {
	t.Parallel()
	want := []int{-10, 1, 10, 11, 20}
	ss := config.MakeSections(
		&config.Section{ID: `aa`, SortOrder: 20},
		&config.Section{ID: `bb`, SortOrder: -10},
		&config.Section{ID: `cc`, SortOrder: 10},
		&config.Section{ID: `dd`, SortOrder: 11},
		&config.Section{ID: `ee`, SortOrder: 1},
	)
	for i, f := range ss.Sort() {
		assert.EqualValues(t, want[i], f.SortOrder)
	}

}

func TestSectionSliceSortAll(t *testing.T) {
	t.Parallel()
	want := `[{"ID":"bb","SortOrder":-10},{"ID":"ee","SortOrder":1},{"ID":"cc","SortOrder":10},{"ID":"aa","SortOrder":20,"Groups":[{"ID":"bb","SortOrder":-10,"Fields":[{"ID":"bb","SortOrder":-10},{"ID":"ee","SortOrder":1},{"ID":"cc","SortOrder":10},{"ID":"dd","SortOrder":11},{"ID":"aa","SortOrder":20}]},{"ID":"ee","SortOrder":1},{"ID":"dd","SortOrder":11,"Fields":[{"ID":"bb","SortOrder":-10},{"ID":"ee","SortOrder":1},{"ID":"cc","SortOrder":10},{"ID":"dd","SortOrder":11},{"ID":"aa","SortOrder":20}]},{"ID":"aa","SortOrder":20,"Fields":[{"ID":"bb","SortOrder":-10},{"ID":"ee","SortOrder":1},{"ID":"cc","SortOrder":10},{"ID":"dd","SortOrder":11},{"ID":"aa","SortOrder":20}]}]}]`
	ss := config.MustMakeSectionsValidate(
		&config.Section{ID: `aa`, SortOrder: 20, Groups: config.MakeGroups(
			&config.Group{
				ID:        `aa`,
				SortOrder: 20,
				Fields: config.MakeFields(
					&config.Field{ID: `aa`, SortOrder: 20},
					&config.Field{ID: `bb`, SortOrder: -10},
					&config.Field{ID: `cc`, SortOrder: 10},
					&config.Field{ID: `dd`, SortOrder: 11},
					&config.Field{ID: `ee`, SortOrder: 1},
				),
			},
			&config.Group{
				ID:        `bb`,
				SortOrder: -10,
				Fields: config.MakeFields(
					&config.Field{ID: `aa`, SortOrder: 20},
					&config.Field{ID: `bb`, SortOrder: -10},
					&config.Field{ID: `cc`, SortOrder: 10},
					&config.Field{ID: `dd`, SortOrder: 11},
					&config.Field{ID: `ee`, SortOrder: 1},
				),
			},
			&config.Group{
				ID:        `dd`,
				SortOrder: 11,
				Fields: config.MakeFields(
					&config.Field{ID: `aa`, SortOrder: 20},
					&config.Field{ID: `bb`, SortOrder: -10},
					&config.Field{ID: `cc`, SortOrder: 10},
					&config.Field{ID: `dd`, SortOrder: 11},
					&config.Field{ID: `ee`, SortOrder: 1},
				),
			},
			&config.Group{ID: `ee`, SortOrder: 1},
		)},
		&config.Section{ID: `bb`, SortOrder: -10},
		&config.Section{ID: `cc`, SortOrder: 10},
		&config.Section{ID: `ee`, SortOrder: 1},
	)
	assert.Exactly(t, 15, ss.TotalFields())
	ss.SortAll()
	have, err := json.Marshal(ss)
	require.NoError(t, err)
	if want != string(have) {
		t.Errorf("\nWant: %s\nHave: %s\n", want, have)
	}
}

func TestSectionSliceAppendFields(t *testing.T) {

	want := `[{"ID":"aa","Groups":[{"ID":"aa","Fields":[{"ID":"aa"},{"ID":"bb"},{"ID":"cc"}]}]}]`
	ss := config.MustMakeSectionsValidate(
		&config.Section{
			ID: `aa`,
			Groups: config.MakeGroups(
				&config.Group{ID: `aa`,
					Fields: config.MakeFields(
						&config.Field{ID: `aa`},
						&config.Field{ID: `bb`},
					),
				},
			)},
	)
	ss, idx := ss.AppendFields("aa/XX/YY")
	require.Exactly(t, -2, idx)

	ss, idx = ss.AppendFields("aa/aa/cc", &config.Field{ID: `cc`})
	require.Exactly(t, 3, idx)

	have, err := json.Marshal(ss)
	require.NoError(t, err)
	if want != string(have) {
		t.Errorf("\nWant: %s\nHave: %s\n", want, have)
	}
}
