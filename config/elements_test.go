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

//
//import (
//	"encoding/json"
//	"testing"
//
//	"github.com/corestoreio/errors"
//	"github.com/corestoreio/pkg/config"
//	"github.com/corestoreio/pkg/store/scope"
//	"github.com/stretchr/testify/assert"
//	"github.com/stretchr/testify/require"
//)
//
//func TestSectionValidateDuplicate(t *testing.T) {
//	// for benchmark tests see package config_bm
//
//	ss := config.MakeSections(
//		&config.Section{
//			ID: `aa`,
//			Groups: config.MakeGroups(
//				&config.Group{
//					ID: `bb`,
//					Fields: config.MakeFields(
//						&config.Field{ID: `cc`},
//						&config.Field{ID: `cc`},
//					),
//				},
//			),
//		},
//	)
//	assert.True(t, errors.NotValid.Match(ss.Validate())) // "Duplicate entry for path aa/bb/cc :: [{\"ID\":\"aa\",\"Groups\":[{\"ID\":\"bb\",\"Fields\":[{\"ID\":\"cc\"},{\"ID\":\"cc\"}]}]}]\n"
//}
//
//func TestSectionValidateShortPath(t *testing.T) {
//
//	ss := config.MakeSections(
//		&config.Section{
//			//ID: `aa`,
//			Groups: config.MakeGroups(
//				&config.Group{
//					//ID: `b`,
//					Fields: config.MakeFields(
//						&config.Field{ID: `ca`},
//						&config.Field{ID: `cb`},
//						&config.Field{},
//					),
//				},
//			),
//		},
//	)
//
//	err := ss.Validate()
//	assert.True(t, errors.Empty.Match(err), "Error %s", err)
//}
//
//func TestSectionUpdateField(t *testing.T) {
//
//	ss := config.MakeSections(
//		&config.Section{
//			ID: `aa`,
//			Groups: config.MakeGroups(
//				&config.Group{
//					ID: `bb`,
//					Fields: config.MakeFields(
//						&config.Field{ID: `ca`},
//						&config.Field{ID: `cb`},
//					),
//				},
//			),
//		},
//	)
//
//	fr := `aa/bb/ca`
//	if err := ss.UpdateField(fr, &config.Field{
//		Label: "ca New Label",
//	}); err != nil {
//		t.Fatal(err)
//	}
//
//	f, _, err := ss.FindField(fr)
//	if err != nil {
//		t.Fatal(err)
//	}
//	assert.Exactly(t, `ca New Label`, f.Label)
//
//	err1 := ss.UpdateField(`a/b/c`, &config.Field{})
//	assert.True(t, errors.NotFound.Match(err1), "Error: %s", err1)
//
//	err2 := ss.UpdateField(`aa/b/c`, &config.Field{})
//	assert.True(t, errors.NotFound.Match(err2), "Error: %s", err2)
//
//	err3 := ss.UpdateField(`aa/bb/c`, &config.Field{})
//	assert.True(t, errors.NotFound.Match(err3), "Error: %s", err3)
//
//	err4 := ss.UpdateField(`aa_bb_c`, &config.Field{})
//	assert.True(t, errors.NotValid.Match(err4), "Error: %s", err4)
//
//}
//
//func TestNewConfiguration(t *testing.T) {
//
//	tests := []struct {
//		have       config.Sections
//		wantErrBhf errors.Kind
//		wantLen    int
//	}{
//		0: {
//			have:       nil,
//			wantErrBhf: errors.NotValid,
//		},
//		1: {
//			have: config.MakeSections(
//				&config.Section{
//					ID: `web`,
//					Groups: config.MakeGroups(
//						&config.Group{
//							ID:     `default`,
//							Fields: config.MakeFields(&config.Field{ID: `front`}, &config.Field{ID: `no_route`}),
//						},
//					),
//				},
//				&config.Section{
//					ID: `system`,
//					Groups: config.MakeGroups(
//						&config.Group{
//							ID:     `media_storage_configuration`,
//							Fields: config.MakeFields(&config.Field{ID: `allowed_resources`}),
//						},
//					),
//				},
//			),
//			wantErrBhf: 0,
//			wantLen:    3,
//		},
//		2: {
//			have:       config.MakeSections(&config.Section{ID: `aa`, Groups: config.MakeGroups()}),
//			wantErrBhf: 0,
//		},
//		3: {
//			have:       config.MakeSections(&config.Section{ID: `aa`, Groups: config.MakeGroups(&config.Group{ID: `bb`, Fields: nil})}),
//			wantErrBhf: 0,
//		},
//		4: {
//			have: config.MakeSections(
//				&config.Section{
//					ID: `aa`,
//					Groups: config.MakeGroups(
//						&config.Group{
//							ID:     `bb`,
//							Fields: config.MakeFields(&config.Field{ID: `cc`}, &config.Field{ID: `cc`}),
//						},
//					),
//				},
//			),
//			wantErrBhf: errors.NotValid, // `Duplicate entry for path aa/bb/cc :: [{"ID":"aa","Groups":[{"ID":"bb","Fields":[{"ID":"cc"},{"ID":"cc"}]}]}]`,
//		},
//	}
//
//	for i, test := range tests {
//		func(t *testing.T, have config.Sections, wantErr errors.Kind) {
//			defer func() {
//				if r := recover(); r != nil {
//					if err, ok := r.(error); ok {
//						assert.True(t, wantErr.Match(err), "Index %d => %s", i, err)
//					} else {
//						t.Errorf("Failed to convert to type error: %#v", r)
//					}
//				} else if wantErr > 0 {
//					t.Errorf("Cannot find panic: wantErr %v", wantErr)
//				}
//			}()
//
//			haveSlice := config.MustMakeSectionsValidate(have...)
//			if wantErr > 0 {
//				assert.Nil(t, haveSlice, "Index %d", i)
//			} else {
//				assert.NotNil(t, haveSlice, "Index %d", i)
//				assert.Len(t, haveSlice, len(have), "Index %d", i)
//			}
//			assert.Exactly(t, test.wantLen, haveSlice.TotalFields(), "Index %d", i)
//		}(t, test.have, test.wantErrBhf)
//	}
//}
//
//func TestSectionSliceMerge(t *testing.T) {
//
//	// Got stuck in comparing JSON?
//	// Use a Webservice to compare the JSON output!
//
//	tests := []struct {
//		have       []config.Sections
//		want       string
//		fieldCount int
//	}{
//		{
//			have: []config.Sections{
//				{
//					&config.Section{
//						ID: `a`,
//					},
//				},
//				{
//					&config.Section{ID: `a`, Label: `LabelA`, Groups: nil},
//				},
//			},
//			want:       `[{"ID":"a","Label":"LabelA","Groups":null}]` + "\n",
//			fieldCount: 0,
//		},
//		{
//			have: []config.Sections{
//				{
//					&config.Section{
//						ID: `a`,
//						Groups: config.MakeGroups(
//							&config.Group{
//								ID: `b`,
//								Fields: config.MakeFields(
//									&config.Field{ID: `c`, Default: `c`},
//								),
//							},
//							&config.Group{
//								ID: `b`,
//								Fields: config.MakeFields(
//									&config.Field{ID: `d`, Default: `d`},
//								),
//							},
//						),
//					},
//				},
//				{
//					&config.Section{ID: `a`, Label: `LabelA`, Groups: nil},
//				},
//			},
//			want:       `[{"ID":"a","Label":"LabelA","Groups":[{"ID":"b","Fields":[{"ID":"c","Default":"c"},{"ID":"d","Default":"d"}]}]}]` + "\n",
//			fieldCount: 2,
//		},
//		{
//			have: []config.Sections{
//				{
//					&config.Section{
//						ID:    `a`,
//						Label: `SectionLabelA`,
//						Groups: config.MakeGroups(
//							&config.Group{
//								ID:     `b`,
//								Scopes: scope.PermDefault,
//								Fields: config.MakeFields(
//									&config.Field{ID: `c`, Default: `c`},
//								),
//							},
//						),
//					},
//				},
//				{
//					&config.Section{
//						ID:     `a`,
//						Scopes: scope.PermWebsite,
//						Groups: config.MakeGroups(
//							&config.Group{ID: `b`, Label: `GroupLabelB1`},
//							&config.Group{ID: `b`, Label: `GroupLabelB2`},
//							&config.Group{
//								ID: `b2`,
//								Fields: config.MakeFields(
//									&config.Field{ID: `d`, Default: `d`},
//								),
//							},
//						),
//					},
//				},
//			},
//			want:       `[{"ID":"a","Label":"SectionLabelA","Scopes":["Default","Website"],"Groups":[{"ID":"b","Label":"GroupLabelB2","Scopes":["Default"],"Fields":[{"ID":"c","Default":"c"}]},{"ID":"b2","Fields":[{"ID":"d","Default":"d"}]}]}]` + "\n",
//			fieldCount: 2,
//		},
//		{
//			have: []config.Sections{
//				{
//					&config.Section{ID: `a`, Label: `SectionLabelA`, SortOrder: 20, Resource: 22},
//				},
//				{
//					&config.Section{ID: `a`, Scopes: scope.PermWebsite, SortOrder: 10, Resource: 3},
//				},
//			},
//			want: `[{"ID":"a","Label":"SectionLabelA","Scopes":["Default","Website"],"SortOrder":10,"Resource":3,"Groups":null}]` + "\n",
//		},
//		{
//			have: []config.Sections{
//				{
//					&config.Section{
//						ID:    `a`,
//						Label: `SectionLabelA`,
//						Groups: config.MakeGroups(
//							&config.Group{
//								ID:      `b`,
//								Label:   `SectionAGroupB`,
//								Comment: "SectionAGroupBComment",
//								Scopes:  scope.PermDefault,
//							},
//						),
//					},
//				},
//				{
//					&config.Section{
//						ID:        `a`,
//						SortOrder: 1000,
//						Scopes:    scope.PermWebsite,
//						Groups: config.MakeGroups(
//							&config.Group{ID: `b`, Label: `GroupLabelB1`, Scopes: scope.PermStore},
//							&config.Group{ID: `b`, Label: `GroupLabelB2`, Comment: "Section2AGroup3BComment", SortOrder: 100},
//							&config.Group{ID: `b2`},
//						),
//					},
//				},
//			},
//			want: `[{"ID":"a","Label":"SectionLabelA","Scopes":["Default","Website"],"SortOrder":1000,"Groups":[{"ID":"b","Label":"GroupLabelB2","Comment":"Section2AGroup3BComment","Scopes":["Default","Website","Store"],"SortOrder":100,"Fields":null},{"ID":"b2","Fields":null}]}]` + "\n",
//		},
//		{
//			have: []config.Sections{
//				{
//					&config.Section{
//						ID: `a`,
//						Groups: config.MakeGroups(
//							&config.Group{
//								ID:    `b`,
//								Label: `b1`,
//								Fields: config.MakeFields(
//									&config.Field{ID: `c`, Default: `c`, Type: config.TypeMultiselect, SortOrder: 1001},
//								),
//							},
//							&config.Group{
//								ID:    `b`,
//								Label: `b2`,
//								Fields: config.MakeFields(
//									&config.Field{ID: `d`, Default: `d`, Comment: "Ring of fire", Type: config.TypeObscure},
//									&config.Field{ID: `c`, Default: `haha`, Type: config.TypeSelect, Scopes: scope.PermWebsite},
//								),
//							},
//						),
//					},
//				},
//				{
//					&config.Section{
//						ID: `a`,
//						Groups: config.MakeGroups(
//							&config.Group{
//								ID:    `b`,
//								Label: `b3`,
//								Fields: config.MakeFields(
//									&config.Field{ID: `d`, Default: `overriddenD`, Label: `Sect2Group2Label4`, Comment: "LOTR"},
//									&config.Field{ID: `c`, Default: `overriddenHaha`, Type: config.TypeHidden},
//								),
//							},
//						),
//					},
//				},
//			},
//			want:       `[{"ID":"a","Groups":[{"ID":"b","Label":"b3","Fields":[{"ID":"c","Type":"hidden","Scopes":["Default","Website"],"SortOrder":1001,"Default":"overriddenHaha"},{"ID":"d","Type":"obscure","Label":"Sect2Group2Label4","Comment":"LOTR","Default":"overriddenD"}]}]}]` + "\n",
//			fieldCount: 2,
//		},
//		{
//			have: []config.Sections{
//				{
//					&config.Section{
//						ID: `a`,
//						Groups: config.MakeGroups(
//							&config.Group{
//								ID: `b`,
//								Fields: config.MakeFields(
//									&config.Field{
//										ID:      `c`,
//										Default: `c`,
//										Type:    config.TypeMultiselect,
//									},
//								),
//							},
//						),
//					},
//				},
//				{
//					&config.Section{
//						ID: `a`,
//						Groups: config.MakeGroups(
//							&config.Group{
//								ID: `b`,
//								Fields: config.MakeFields(
//									&config.Field{
//										ID:        `c`,
//										Default:   `overridenC`,
//										Type:      config.TypeSelect,
//										Label:     `Sect2Group2Label4`,
//										Comment:   "LOTR",
//										SortOrder: 100,
//										Visible:   true,
//									},
//								),
//							},
//						),
//					},
//				},
//			},
//			fieldCount: 1,
//			want:       `[{"ID":"a","Groups":[{"ID":"b","Fields":[{"ID":"c","Type":"select","Label":"Sect2Group2Label4","Comment":"LOTR","SortOrder":100,"Visible":true,"Default":"overridenC"}]}]}]` + "\n",
//		},
//	}
//
//	for i, test := range tests {
//
//		if len(test.have) == 0 {
//			test.want = "null\n"
//		}
//
//		var baseSl config.Sections
//		baseSl = baseSl.MergeMultiple(test.have...)
//		j, err := json.Marshal(baseSl)
//		require.NoError(t, err)
//
//		if string(j) != test.want {
//			t.Errorf("\nIndex: %d\nExpected: %s\nActual:   %s\n", i, test.want, j)
//		}
//
//		assert.Exactly(t, test.fieldCount, baseSl.TotalFields(), "Index %d", i)
//	}
//}
//
//func TestFieldSliceMerge(t *testing.T) {
//
//	tests := []struct {
//		have config.Fields
//		want string
//	}{
//		{
//			have: config.MakeFields(
//				&config.Field{ID: `d`, Default: `overrideMeD`, Comment: "Ring of fire", Type: config.TypeObscure},
//			),
//			want: `[{"ID":"d","Type":"obscure","Comment":"Ring of fire","Default":"overrideMeD"}]`,
//		},
//		{
//			have: config.MakeFields(
//				&config.Field{ID: `d`, Default: `overrideMeD`, Comment: "Ring of fire", Type: config.TypeObscure},
//				&config.Field{ID: `c`, Default: `overrideMeC`, Type: config.TypeSelect, Scopes: scope.PermWebsite},
//			),
//			want: `[{"ID":"d","Type":"obscure","Comment":"Ring of fire","Default":"overrideMeD"},{"ID":"c","Type":"select","Scopes":["Default","Website"],"Default":"overrideMeC"}]`,
//		},
//		{
//			have: config.MakeFields(
//				&config.Field{ID: `d`, Default: `overrideMeD`, Comment: "Ring of fire", Type: config.TypeObscure},
//				&config.Field{ID: `c`, Default: `overrideMeC`, Type: config.TypeSelect, Scopes: scope.PermWebsite},
//				&config.Field{ID: `d`, Default: `overrideMeE`, Type: config.TypeMultiselect},
//			),
//			want: `[{"ID":"d","Type":"multiselect","Comment":"Ring of fire","Default":"overrideMeE"},{"ID":"c","Type":"select","Scopes":["Default","Website"],"Default":"overrideMeC"}]`,
//		},
//		{
//			have: nil,
//			want: `null`,
//		},
//	}
//
//	for i, test := range tests {
//		var baseFsl config.Fields
//		baseFsl = baseFsl.Merge(test.have...)
//
//		fsj, err := json.Marshal(baseFsl)
//		if err != nil {
//			t.Fatal(err)
//		}
//		if string(fsj) != test.want {
//			t.Errorf("\nIndex: %d\nExpected: %s\nActual:   %s\n", i, test.want, string(fsj))
//		}
//
//	}
//}
//
//func TestGroupSliceMerge(t *testing.T) {
//
//	tests := []struct {
//		have config.Groups
//		want string
//	}{
//		{
//			have: config.MakeGroups(
//				&config.Group{
//					ID: `b`,
//					Fields: config.MakeFields(
//						&config.Field{ID: `c`, Default: `c`, Type: config.TypeMultiselect},
//					),
//				},
//				&config.Group{
//					ID: `b`,
//					Fields: config.MakeFields(
//						&config.Field{ID: `d`, Default: `overrideMeD`, Comment: "Ring of fire", Type: config.TypeObscure},
//						&config.Field{ID: `c`, Default: `overrideMeC`, Type: config.TypeSelect, Scopes: scope.PermWebsite},
//					),
//				},
//				&config.Group{
//					ID: `b`,
//					Fields: config.MakeFields(
//						&config.Field{ID: `d`, Default: `overriddenD`, Label: `Sect2Group2Label4`, Comment: "LOTR"},
//						&config.Field{ID: `c`, Default: `overriddenC`, Type: config.TypeHidden},
//					),
//				},
//			),
//			want: `[{"ID":"b","Fields":[{"ID":"c","Type":"hidden","Scopes":["Default","Website"],"Default":"overriddenC"},{"ID":"d","Type":"obscure","Label":"Sect2Group2Label4","Comment":"LOTR","Default":"overriddenD"}]}]` + "\n",
//		},
//		{
//			have: config.MakeGroups(
//				&config.Group{
//					ID:    `b`,
//					Label: `Single Field`,
//					Fields: config.MakeFields(
//						&config.Field{ID: `c`, Default: `c`, Type: config.TypeMultiselect},
//					),
//				},
//			),
//			want: `[{"ID":"b","Label":"Single Field","Fields":[{"ID":"c","Type":"multiselect","Default":"c"}]}]` + "\n",
//		},
//		{
//			have: config.MakeGroups(
//				&config.Group{
//					ID:    `b`,
//					Label: `Single Field`,
//					Fields: config.MakeFields(
//						&config.Field{ID: `c`, Default: `c`, Type: config.TypeMultiselect},
//					),
//				},
//				&config.Group{
//					ID:    `b`,
//					Label: `Single Field2`,
//					Fields: config.MakeFields(
//						&config.Field{ID: `c`, Default: `c`, Type: config.TypeMultiselect},
//					),
//				},
//			),
//			want: `[{"ID":"b","Label":"Single Field2","Fields":[{"ID":"c","Type":"multiselect","Default":"c"}]}]` + "\n",
//		},
//		{
//			have: config.MakeGroups(
//				&config.Group{
//					ID:    `b`,
//					Label: `Single Field`,
//					Fields: config.MakeFields(
//						&config.Field{ID: `c`, Default: `c`, Type: config.TypeMultiselect},
//					),
//				},
//				&config.Group{
//					ID:    `b`,
//					Label: `Single Field2`,
//					Fields: config.MakeFields(
//						&config.Field{ID: `c`, Default: `c2`, Type: config.TypeTextarea},
//					),
//				},
//			),
//			want: `[{"ID":"b","Label":"Single Field2","Fields":[{"ID":"c","Type":"textarea","Default":"c2"}]}]` + "\n",
//		},
//		{
//			have: config.MakeGroups(
//				&config.Group{
//					ID:    `b`,
//					Label: `Single Field`,
//					Fields: config.MakeFields(
//						&config.Field{ID: `c`, Default: `c`, Type: config.TypeMultiselect},
//					),
//				},
//				&config.Group{
//					ID: `b`,
//					Fields: config.MakeFields(
//						&config.Field{ID: `d`, Default: `d`, Type: config.TypeText},
//					),
//				},
//			),
//			want: `[{"ID":"b","Label":"Single Field","Fields":[{"ID":"c","Type":"multiselect","Default":"c"},{"ID":"d","Type":"text","Default":"d"}]}]` + "\n",
//		},
//		{
//			have: nil,
//			want: `null` + "\n",
//		},
//	}
//
//	for i, test := range tests {
//		var baseGsl config.Groups
//		baseGsl = baseGsl.Merge(test.have...)
//
//		j, err := json.Marshal(baseGsl)
//		require.NoError(t, err)
//		if string(j) != test.want {
//			t.Errorf("\nIndex: %d\nExpected: %s\nActual:   %s\n", i, test.want, j)
//		}
//
//	}
//}
//
//func TestSectionSliceFindGroupByID(t *testing.T) {
//
//	tests := []struct {
//		haveSlice  config.Sections
//		haveRoute  string
//		wantGID    string
//		wantErrBhf errors.Kind
//	}{
//		{
//			haveSlice:  config.MakeSections(&config.Section{ID: `a`, Groups: config.MakeGroups(&config.Group{ID: `b`}, &config.Group{ID: `bb`})}),
//			haveRoute:  "a/b",
//			wantGID:    "b",
//			wantErrBhf: errors.NoKind,
//		},
//		{
//			haveSlice:  config.MakeSections(&config.Section{ID: `a`, Groups: config.MakeGroups(&config.Group{ID: `b`}, &config.Group{ID: `bb`})}),
//			haveRoute:  "a/bc",
//			wantGID:    "b",
//			wantErrBhf: errors.NotFound,
//		},
//		{
//			haveSlice:  config.Sections{},
//			haveRoute:  "",
//			wantGID:    "b",
//			wantErrBhf: errors.NotFound,
//		},
//		{
//			haveSlice:  config.MakeSections(&config.Section{ID: `a`, Groups: config.MakeGroups(&config.Group{ID: `b`}, &config.Group{ID: `bb`})}),
//			haveRoute:  "a/bb/cc",
//			wantGID:    "bb",
//			wantErrBhf: errors.NoKind,
//		},
//		{
//			haveSlice:  config.MakeSections(&config.Section{ID: `a`, Groups: config.MakeGroups(&config.Group{ID: `b`}, &config.Group{ID: `bb`})}),
//			haveRoute:  "xa/bb/cc",
//			wantGID:    "",
//			wantErrBhf: errors.NotFound,
//		},
//	}
//
//	for i, test := range tests {
//		haveGroup, _, haveErr := test.haveSlice.FindGroup(test.haveRoute)
//		if test.wantErrBhf > 0 {
//			assert.Error(t, haveErr, "Index %d", i)
//			assert.Exactly(t, &config.Group{}, haveGroup)
//			assert.True(t, test.wantErrBhf.Match(haveErr), "Index %d => %s", i, haveErr)
//			continue
//		}
//
//		assert.NoError(t, haveErr, "Index %d", i)
//		assert.NotNil(t, haveGroup, "Index %d", i)
//		assert.Exactly(t, test.wantGID, haveGroup.ID)
//	}
//}
//
//func TestSectionSliceFindFieldByID(t *testing.T) {
//
//	tests := []struct {
//		haveSlice  config.Sections
//		haveRoute  string
//		wantFID    string
//		wantErrBhf errors.Kind
//	}{
//		{
//			haveSlice:  config.MakeSections(&config.Section{ID: `aa`, Groups: config.MakeGroups(&config.Group{ID: `bb`}, &config.Group{ID: `cc`})}),
//			haveRoute:  "",
//			wantFID:    "",
//			wantErrBhf: errors.NotValid,
//		},
//		{
//			haveSlice:  config.MakeSections(&config.Section{ID: `a`, Groups: config.MakeGroups(&config.Group{ID: `b`}, &config.Group{ID: `bb`})}),
//			haveRoute:  "a/b",
//			wantFID:    "b",
//			wantErrBhf: errors.NotFound,
//		},
//		{
//			haveSlice:  config.MakeSections(&config.Section{ID: `a`, Groups: config.MakeGroups(&config.Group{ID: `b`}, &config.Group{ID: `bb`})}),
//			haveRoute:  "a/bc",
//			wantFID:    "b",
//			wantErrBhf: errors.NotFound,
//		},
//		{
//			haveSlice:  config.MakeSections(),
//			haveRoute:  "",
//			wantFID:    "",
//			wantErrBhf: errors.NotValid,
//		},
//		{
//			haveSlice:  config.MakeSections(&config.Section{ID: `a`, Groups: config.MakeGroups(&config.Group{ID: `b`}, &config.Group{ID: `bb`})}),
//			haveRoute:  "a/bb/cc",
//			wantFID:    "bb",
//			wantErrBhf: errors.NotFound,
//		},
//		{
//			haveSlice:  config.MakeSections(&config.Section{ID: `a`, Groups: config.MakeGroups(&config.Group{ID: `b`}, &config.Group{ID: `bb`})}),
//			haveRoute:  "xa/bb/cc",
//			wantFID:    "",
//			wantErrBhf: errors.NotFound,
//		},
//		{
//			haveSlice:  config.MakeSections(&config.Section{ID: `a1`, Groups: config.MakeGroups(&config.Group{ID: `b1`, Fields: config.MakeFields(&config.Field{ID: `c1`})})}),
//			haveRoute:  "a1/b1/c1",
//			wantFID:    "c1",
//			wantErrBhf: errors.NoKind,
//		},
//	}
//
//	for i, test := range tests {
//		haveGroup, _, haveErr := test.haveSlice.FindField(test.haveRoute)
//		if test.wantErrBhf > 0 {
//			assert.Error(t, haveErr, "Index %d", i)
//			assert.Exactly(t, &config.Field{}, haveGroup, "Index %d", i)
//			assert.True(t, test.wantErrBhf.Match(haveErr), "Index %d => %s", i, haveErr)
//			continue
//		}
//		assert.NoError(t, haveErr, "Index %d", i)
//		assert.NotNil(t, haveGroup, "Index %d", i)
//		assert.Exactly(t, test.wantFID, haveGroup.ID, "Index %d", i)
//	}
//}
//
//func TestFieldSliceSort(t *testing.T) {
//
//	want := []int{-10, 1, 10, 11, 20}
//	fs := config.MakeFields(
//		&config.Field{ID: `aa`, SortOrder: 20},
//		&config.Field{ID: `bb`, SortOrder: -10},
//		&config.Field{ID: `cc`, SortOrder: 10},
//		&config.Field{ID: `dd`, SortOrder: 11},
//		&config.Field{ID: `ee`, SortOrder: 1},
//	)
//
//	for i, f := range fs.Sort() {
//		assert.EqualValues(t, want[i], f.SortOrder)
//	}
//}
//
//func TestGroupSliceSort(t *testing.T) {
//
//	want := []int{-10, 1, 10, 11, 20}
//	gs := config.MakeGroups(
//		&config.Group{ID: `aa`, SortOrder: 20},
//		&config.Group{ID: `bb`, SortOrder: -10},
//		&config.Group{ID: `cc`, SortOrder: 10},
//		&config.Group{ID: `dd`, SortOrder: 11},
//		&config.Group{ID: `ee`, SortOrder: 1},
//	)
//	for i, f := range gs.Sort() {
//		assert.EqualValues(t, want[i], f.SortOrder)
//	}
//}
//func TestSectionSliceSort(t *testing.T) {
//
//	want := []int{-10, 1, 10, 11, 20}
//	ss := config.MakeSections(
//		&config.Section{ID: `aa`, SortOrder: 20},
//		&config.Section{ID: `bb`, SortOrder: -10},
//		&config.Section{ID: `cc`, SortOrder: 10},
//		&config.Section{ID: `dd`, SortOrder: 11},
//		&config.Section{ID: `ee`, SortOrder: 1},
//	)
//	for i, f := range ss.Sort() {
//		assert.EqualValues(t, want[i], f.SortOrder)
//	}
//
//}
//
//func TestSectionSliceSortAll(t *testing.T) {
//
//	want := `[{"ID":"bb","SortOrder":-10,"Groups":null},{"ID":"ee","SortOrder":1,"Groups":null},{"ID":"cc","SortOrder":10,"Groups":null},{"ID":"aa","SortOrder":20,"Groups":[{"ID":"bb","SortOrder":-10,"Fields":[{"ID":"bb","SortOrder":-10},{"ID":"ee","SortOrder":1},{"ID":"cc","SortOrder":10},{"ID":"dd","SortOrder":11},{"ID":"aa","SortOrder":20}]},{"ID":"ee","SortOrder":1,"Fields":null},{"ID":"dd","SortOrder":11,"Fields":[{"ID":"bb","SortOrder":-10},{"ID":"ee","SortOrder":1},{"ID":"cc","SortOrder":10},{"ID":"dd","SortOrder":11},{"ID":"aa","SortOrder":20}]},{"ID":"aa","SortOrder":20,"Fields":[{"ID":"bb","SortOrder":-10},{"ID":"ee","SortOrder":1},{"ID":"cc","SortOrder":10},{"ID":"dd","SortOrder":11},{"ID":"aa","SortOrder":20}]}]}]` + "\n"
//	ss := config.MustMakeSectionsValidate(
//		&config.Section{ID: `aa`, SortOrder: 20, Groups: config.MakeGroups(
//			&config.Group{
//				ID:        `aa`,
//				SortOrder: 20,
//				Fields: config.MakeFields(
//					&config.Field{ID: `aa`, SortOrder: 20},
//					&config.Field{ID: `bb`, SortOrder: -10},
//					&config.Field{ID: `cc`, SortOrder: 10},
//					&config.Field{ID: `dd`, SortOrder: 11},
//					&config.Field{ID: `ee`, SortOrder: 1},
//				),
//			},
//			&config.Group{
//				ID:        `bb`,
//				SortOrder: -10,
//				Fields: config.MakeFields(
//					&config.Field{ID: `aa`, SortOrder: 20},
//					&config.Field{ID: `bb`, SortOrder: -10},
//					&config.Field{ID: `cc`, SortOrder: 10},
//					&config.Field{ID: `dd`, SortOrder: 11},
//					&config.Field{ID: `ee`, SortOrder: 1},
//				),
//			},
//			&config.Group{
//				ID:        `dd`,
//				SortOrder: 11,
//				Fields: config.MakeFields(
//					&config.Field{ID: `aa`, SortOrder: 20},
//					&config.Field{ID: `bb`, SortOrder: -10},
//					&config.Field{ID: `cc`, SortOrder: 10},
//					&config.Field{ID: `dd`, SortOrder: 11},
//					&config.Field{ID: `ee`, SortOrder: 1},
//				),
//			},
//			&config.Group{ID: `ee`, SortOrder: 1},
//		)},
//		&config.Section{ID: `bb`, SortOrder: -10},
//		&config.Section{ID: `cc`, SortOrder: 10},
//		&config.Section{ID: `ee`, SortOrder: 1},
//	)
//	assert.Exactly(t, 15, ss.TotalFields())
//	ss.SortAll()
//	have, err := json.Marshal(ss)
//	require.NoError(t, err)
//	if want != string(have) {
//		t.Errorf("\nWant: %s\nHave: %s\n", want, have)
//	}
//}
//
//func TestSectionSliceAppendFields(t *testing.T) {
//
//	want := `[{"ID":"aa","Groups":[{"ID":"aa","Fields":[{"ID":"aa"},{"ID":"bb"},{"ID":"cc"}]}]}]` + "\n"
//	ss := config.MustMakeSectionsValidate(
//		&config.Section{
//			ID: `aa`,
//			Groups: config.MakeGroups(
//				&config.Group{ID: `aa`,
//					Fields: config.MakeFields(
//						&config.Field{ID: `aa`},
//						&config.Field{ID: `bb`},
//					),
//				},
//			)},
//	)
//	ss, err := ss.AppendFields("aa/XX")
//	assert.True(t, errors.NotFound.Match(err))
//
//	ss, err = ss.AppendFields(("aa/aa"), &config.Field{ID: `cc`})
//	assert.NoError(t, err)
//	have, err := json.Marshal(ss)
//	require.NoError(t, err)
//	if want != string(have) {
//		t.Errorf("\nWant: %s\nHave: %s\n", want, have)
//	}
//}
