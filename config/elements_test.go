// Copyright 2015 CoreStore Authors
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
	"testing"

	"github.com/corestoreio/csfw/config"
	"github.com/stretchr/testify/assert"
)

func TestNewConfiguration(t *testing.T) {
	tests := []struct {
		have    []*config.Section
		wantErr string
		wantLen int
	}{
		0: {
			have:    []*config.Section{},
			wantErr: "SectionSlice is empty",
		},
		1: {
			have: []*config.Section{
				&config.Section{
					ID: "web",
					Groups: config.GroupSlice{
						&config.Group{
							ID:     "default",
							Fields: config.FieldSlice{&config.Field{ID: "front"}, &config.Field{ID: "no_route"}},
						},
					},
				},
				&config.Section{
					ID: "system",
					Groups: config.GroupSlice{
						&config.Group{
							ID:     "media_storage_configuration",
							Fields: config.FieldSlice{&config.Field{ID: "allowed_resources"}},
						},
					},
				},
			},
			wantErr: "",
			wantLen: 3,
		},
		2: {
			have:    []*config.Section{&config.Section{ID: "a", Groups: config.GroupSlice{}}},
			wantErr: "",
		},
		3: {
			have:    []*config.Section{&config.Section{ID: "a", Groups: config.GroupSlice{&config.Group{ID: "b", Fields: nil}}}},
			wantErr: "",
		},
		4: {
			have: []*config.Section{
				&config.Section{
					ID: "a",
					Groups: config.GroupSlice{
						&config.Group{
							ID:     "b",
							Fields: config.FieldSlice{&config.Field{ID: "c"}, &config.Field{ID: "c"}},
						},
					},
				},
			},
			wantErr: "Duplicate entry for path default/0/a/b/c",
		},
	}

	for i, test := range tests {
		func(t *testing.T, have []*config.Section, wantErr string) {
			defer func() {
				if r := recover(); r != nil {
					if err, ok := r.(error); ok {
						assert.Contains(t, err.Error(), wantErr)
					} else {
						t.Errorf("Failed to convert to type error: %#v", err)
					}
				} else if wantErr != "" {
					t.Errorf("Cannot find panic: wantErr %s", wantErr)
				}
			}()

			haveSlice := config.NewConfiguration(have...)
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
	pkgCfg := config.NewConfiguration(
		&config.Section{
			ID: "contact",
			Groups: config.GroupSlice{
				&config.Group{
					ID: "contact",
					Fields: config.FieldSlice{
						&config.Field{
							// Path: `contact/contact/enabled`,
							ID:      "enabled",
							Default: true,
						},
					},
				},
				&config.Group{
					ID: "email",
					Fields: config.FieldSlice{
						&config.Field{
							// Path: `contact/email/recipient_email`,
							ID:      "recipient_email",
							Default: `hello@example.com`,
						},
						&config.Field{
							// Path: `contact/email/sender_email_identity`,
							ID:      "sender_email_identity",
							Default: 2.7182818284590452353602874713527,
						},
						&config.Field{
							// Path: `contact/email/email_template`,
							ID:      "email_template",
							Default: 4711,
						},
					},
				},
			},
		},
	)

	assert.Exactly(
		t,
		config.DefaultMap{"default/0/contact/email/sender_email_identity": 2.718281828459045, "default/0/contact/email/email_template": 4711, "default/0/contact/contact/enabled": true, "default/0/contact/email/recipient_email": "hello@example.com"},
		pkgCfg.Defaults(),
	)
}

func TestSectionSliceMerge(t *testing.T) {

	// Got stuck in comparing JSON?
	// Use a Webservice to compare the JSON output!

	tests := []struct {
		have    []config.SectionSlice
		wantErr string
		want    string
		wantLen int
	}{
		0: {
			have: []config.SectionSlice{
				nil,
				config.SectionSlice{
					nil,
					&config.Section{
						ID: "a",
						Groups: config.GroupSlice{
							nil,
							&config.Group{
								ID: "b",
								Fields: config.FieldSlice{
									&config.Field{ID: "c", Default: `c`},
								},
							},
							&config.Group{
								ID: "b",
								Fields: config.FieldSlice{
									&config.Field{ID: "d", Default: `d`},
								},
							},
						},
					},
				},
				config.SectionSlice{
					&config.Section{ID: "a", Label: "LabelA", Groups: nil},
				},
			},
			wantErr: "",
			want:    `[{"ID":"a","Label":"LabelA","Groups":[{"ID":"b","Fields":[{"ID":"c","Default":"c"},{"ID":"d","Default":"d"}]}]}]` + "\n",
			wantLen: 2,
		},
		1: {
			have: []config.SectionSlice{
				config.SectionSlice{
					&config.Section{
						ID:    "a",
						Label: "SectionLabelA",
						Groups: config.GroupSlice{
							&config.Group{
								ID:    "b",
								Scope: config.NewScopePerm(config.ScopeDefaultID),
								Fields: config.FieldSlice{
									&config.Field{ID: "c", Default: `c`},
								},
							},
							nil,
						},
					},
				},
				config.SectionSlice{
					&config.Section{
						ID:    "a",
						Scope: config.NewScopePerm(config.ScopeDefaultID, config.ScopeWebsiteID),
						Groups: config.GroupSlice{
							&config.Group{ID: "b", Label: "GroupLabelB1"},
							nil,
							&config.Group{ID: "b", Label: "GroupLabelB2"},
							&config.Group{
								ID: "b2",
								Fields: config.FieldSlice{
									&config.Field{ID: "d", Default: `d`},
								},
							},
						},
					},
				},
			},
			wantErr: "",
			want:    `[{"ID":"a","Label":"SectionLabelA","Scope":["ScopeDefault","ScopeWebsite"],"Groups":[{"ID":"b","Label":"GroupLabelB2","Scope":["ScopeDefault"],"Fields":[{"ID":"c","Default":"c"}]},{"ID":"b2","Fields":[{"ID":"d","Default":"d"}]}]}]` + "\n",
			wantLen: 2,
		},
		2: {
			have: []config.SectionSlice{
				config.SectionSlice{
					&config.Section{ID: "a", Label: "SectionLabelA", SortOrder: 20, Permission: 22},
				},
				config.SectionSlice{
					&config.Section{ID: "a", Scope: config.NewScopePerm(config.ScopeDefaultID, config.ScopeWebsiteID), SortOrder: 10, Permission: 3},
				},
			},
			wantErr: "",
			want:    `[{"ID":"a","Label":"SectionLabelA","Scope":["ScopeDefault","ScopeWebsite"],"SortOrder":10,"Permission":3,"Groups":null}]` + "\n",
		},
		3: {
			have: []config.SectionSlice{
				config.SectionSlice{
					&config.Section{
						ID:    "a",
						Label: "SectionLabelA",
						Groups: config.GroupSlice{
							&config.Group{
								ID:      "b",
								Label:   "SectionAGroupB",
								Comment: "SectionAGroupBComment",
								Scope:   config.NewScopePerm(config.ScopeDefaultID),
							},
						},
					},
				},
				config.SectionSlice{
					&config.Section{
						ID:        "a",
						SortOrder: 1000,
						Scope:     config.NewScopePerm(config.ScopeDefaultID, config.ScopeWebsiteID),
						Groups: config.GroupSlice{
							&config.Group{ID: "b", Label: "GroupLabelB1", Scope: config.ScopePermAll},
							&config.Group{ID: "b", Label: "GroupLabelB2", Comment: "Section2AGroup3BComment", SortOrder: 100},
							&config.Group{ID: "b2"},
						},
					},
				},
			},
			wantErr: "",
			want:    `[{"ID":"a","Label":"SectionLabelA","Scope":["ScopeDefault","ScopeWebsite"],"SortOrder":1000,"Groups":[{"ID":"b","Label":"GroupLabelB2","Comment":"Section2AGroup3BComment","Scope":["ScopeDefault","ScopeWebsite","ScopeStore"],"SortOrder":100,"Fields":null},{"ID":"b2","Fields":null}]}]` + "\n",
		},
		4: {
			have: []config.SectionSlice{
				config.SectionSlice{
					&config.Section{
						ID: "a",
						Groups: config.GroupSlice{
							&config.Group{
								ID:    "b",
								Label: "b1",
								Fields: config.FieldSlice{
									&config.Field{ID: "c", Default: `c`, Type: config.TypeMultiselect, SortOrder: 1001},
								},
							},
							&config.Group{
								ID:    "b",
								Label: "b2",
								Fields: config.FieldSlice{
									nil,
									&config.Field{ID: "d", Default: `d`, Comment: "Ring of fire", Type: config.TypeObscure},
									&config.Field{ID: "c", Default: `haha`, Type: config.TypeSelect, Scope: config.NewScopePerm(config.ScopeDefaultID, config.ScopeWebsiteID)},
								},
							},
						},
					},
				},
				config.SectionSlice{
					&config.Section{
						ID: "a",
						Groups: config.GroupSlice{
							&config.Group{
								ID:    "b",
								Label: "b3",
								Fields: config.FieldSlice{
									&config.Field{ID: "d", Default: `overriddenD`, Label: "Sect2Group2Label4", Comment: "LOTR"},
									&config.Field{ID: "c", Default: `overriddenHaha`, Type: config.TypeHidden},
								},
							},
						},
					},
				},
			},
			wantErr: "",
			want:    `[{"ID":"a","Groups":[{"ID":"b","Label":"b3","Fields":[{"ID":"c","Type":"hidden","Scope":["ScopeDefault","ScopeWebsite"],"SortOrder":1001,"Default":"overriddenHaha"},{"ID":"d","Type":"obscure","Label":"Sect2Group2Label4","Comment":"LOTR","Default":"overriddenD"}]}]}]` + "\n",
			wantLen: 2,
		},
		5: {
			have: []config.SectionSlice{
				config.SectionSlice{
					&config.Section{
						ID: "a",
						Groups: config.GroupSlice{
							&config.Group{
								ID: "b",
								Fields: config.FieldSlice{
									&config.Field{
										ID:      "c",
										Default: `c`,
										Type:    config.TypeMultiselect,
									},
								},
							},
						},
					},
				},
				config.SectionSlice{
					nil,
					&config.Section{
						ID: "a",
						Groups: config.GroupSlice{
							&config.Group{
								ID: "b",
								Fields: config.FieldSlice{
									nil,
									&config.Field{
										ID:        "c",
										Default:   `overridenC`,
										Type:      config.TypeSelect,
										Label:     "Sect2Group2Label4",
										Comment:   "LOTR",
										SortOrder: 100,
										Visible:   config.VisibleYes,
									},
								},
							},
						},
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

		var baseSl config.SectionSlice
		haveErr := baseSl.MergeMultiple(test.have...)
		if test.wantErr != "" {
			assert.Len(t, baseSl, 0)
			assert.Error(t, haveErr)
			assert.Contains(t, haveErr.Error(), test.wantErr)
		} else {
			assert.NoError(t, haveErr)
			j := baseSl.ToJson()
			if j != test.want {
				t.Errorf("\nIndex: %d\nExpected: %s\nActual:   %s\n", i, test.want, j)
			}
		}
		assert.Exactly(t, test.wantLen, baseSl.TotalFields(), "Index %d", i)
	}
}

func TestGroupSliceMerge(t *testing.T) {

	tests := []struct {
		have    []*config.Group
		wantErr error
		want    string
	}{
		{
			have: []*config.Group{
				&config.Group{
					ID: "b",
					Fields: config.FieldSlice{
						&config.Field{ID: "c", Default: `c`, Type: config.TypeMultiselect},
					},
				},
				&config.Group{
					ID: "b",
					Fields: config.FieldSlice{
						&config.Field{ID: "d", Default: `d`, Comment: "Ring of fire", Type: config.TypeObscure},
						&config.Field{ID: "c", Default: `haha`, Type: config.TypeSelect, Scope: config.NewScopePerm(config.ScopeDefaultID, config.ScopeWebsiteID)},
					},
				},
				&config.Group{
					ID: "b",
					Fields: config.FieldSlice{
						&config.Field{ID: "d", Default: `overriddenD`, Label: "Sect2Group2Label4", Comment: "LOTR"},
						&config.Field{ID: "c", Default: `overriddenHaha`, Type: config.TypeHidden},
					},
				},
			},
			wantErr: nil,
			want:    `[{"ID":"b","Fields":[{"ID":"c","Type":"hidden","Scope":["ScopeDefault","ScopeWebsite"],"Default":"overriddenHaha"},{"ID":"d","Type":"obscure","Label":"Sect2Group2Label4","Comment":"LOTR","Default":"overriddenD"}]}]` + "\n",
		},
		{
			have:    nil,
			wantErr: nil,
			want:    `null` + "\n",
		},
	}

	for i, test := range tests {
		var baseGsl config.GroupSlice
		haveErr := baseGsl.Merge(test.have...)
		if test.wantErr != nil {
			assert.Len(t, baseGsl, 0)
			assert.Error(t, haveErr)
			assert.Contains(t, haveErr.Error(), test.wantErr)
		} else {
			assert.NoError(t, haveErr)
			j := baseGsl.ToJson()
			if j != test.want {
				t.Errorf("\nIndex: %d\nExpected: %s\nActual:   %s\n", i, test.want, j)
			}
		}
	}
}

func TestSectionSliceFindGroupByPath(t *testing.T) {
	tests := []struct {
		haveSlice config.SectionSlice
		havePath  []string
		wantGID   string
		wantErr   error
	}{
		0: {
			haveSlice: config.SectionSlice{&config.Section{ID: "a", Groups: config.GroupSlice{&config.Group{ID: "b"}, &config.Group{ID: "bb"}}}},
			havePath:  []string{"a/b"},
			wantGID:   "b",
			wantErr:   nil,
		},
		1: {
			haveSlice: config.SectionSlice{&config.Section{ID: "a", Groups: config.GroupSlice{&config.Group{ID: "b"}, &config.Group{ID: "bb"}}}},
			havePath:  []string{"a/bc"},
			wantGID:   "b",
			wantErr:   config.ErrGroupNotFound,
		},
		2: {
			haveSlice: config.SectionSlice{},
			havePath:  nil,
			wantGID:   "b",
			wantErr:   config.ErrGroupNotFound,
		},
		3: {
			haveSlice: config.SectionSlice{&config.Section{ID: "a", Groups: config.GroupSlice{&config.Group{ID: "b"}, &config.Group{ID: "bb"}}}},
			havePath:  []string{"a", "bb", "cc"},
			wantGID:   "bb",
			wantErr:   nil,
		},
		4: {
			haveSlice: config.SectionSlice{&config.Section{ID: "a", Groups: config.GroupSlice{&config.Group{ID: "b"}, &config.Group{ID: "bb"}}}},
			havePath:  []string{"xa", "bb", "cc"},
			wantGID:   "",
			wantErr:   config.ErrSectionNotFound,
		},
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
		haveSlice config.SectionSlice
		havePath  []string
		wantFID   string
		wantErr   error
	}{
		0: {
			haveSlice: config.SectionSlice{&config.Section{ID: "a", Groups: config.GroupSlice{&config.Group{ID: "b"}, &config.Group{ID: "bb"}}}},
			havePath:  []string{"a/b"},
			wantFID:   "b",
			wantErr:   config.ErrFieldNotFound,
		},
		1: {
			haveSlice: config.SectionSlice{&config.Section{ID: "a", Groups: config.GroupSlice{&config.Group{ID: "b"}, &config.Group{ID: "bb"}}}},
			havePath:  []string{"a/bc"},
			wantFID:   "b",
			wantErr:   config.ErrFieldNotFound,
		},
		2: {
			haveSlice: config.SectionSlice{nil},
			havePath:  nil,
			wantFID:   "b",
			wantErr:   config.ErrFieldNotFound,
		},
		3: {
			haveSlice: config.SectionSlice{&config.Section{ID: "a", Groups: config.GroupSlice{nil, &config.Group{ID: "b"}, &config.Group{ID: "bb"}}}},
			havePath:  []string{"a", "bb", "cc"},
			wantFID:   "bb",
			wantErr:   config.ErrFieldNotFound,
		},
		4: {
			haveSlice: config.SectionSlice{&config.Section{ID: "a", Groups: config.GroupSlice{&config.Group{ID: "b"}, &config.Group{ID: "bb"}}}},
			havePath:  []string{"xa", "bb", "cc"},
			wantFID:   "",
			wantErr:   config.ErrSectionNotFound,
		},
		5: {
			haveSlice: config.SectionSlice{&config.Section{ID: "a", Groups: config.GroupSlice{&config.Group{ID: "b", Fields: config.FieldSlice{
				&config.Field{ID: "c"}, nil,
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
	fs := config.FieldSlice{
		&config.Field{ID: "a", SortOrder: 20},
		&config.Field{ID: "b", SortOrder: -10},
		&config.Field{ID: "c", SortOrder: 10},
		&config.Field{ID: "d", SortOrder: 11},
		&config.Field{ID: "e", SortOrder: 1},
	}

	for i, f := range *(fs.Sort()) {
		assert.EqualValues(t, want[i], f.SortOrder)
	}
}

func TestGroupSliceSort(t *testing.T) {
	want := []int{-10, 1, 10, 11, 20}
	gs := config.GroupSlice{
		&config.Group{ID: "a", SortOrder: 20},
		&config.Group{ID: "b", SortOrder: -10},
		&config.Group{ID: "c", SortOrder: 10},
		&config.Group{ID: "d", SortOrder: 11},
		&config.Group{ID: "e", SortOrder: 1},
	}
	for i, f := range *(gs.Sort()) {
		assert.EqualValues(t, want[i], f.SortOrder)
	}
}
func TestSectionSliceSort(t *testing.T) {
	want := []int{-10, 1, 10, 11, 20}
	ss := config.SectionSlice{
		&config.Section{ID: "a", SortOrder: 20},
		&config.Section{ID: "b", SortOrder: -10},
		&config.Section{ID: "c", SortOrder: 10},
		&config.Section{ID: "d", SortOrder: 11},
		&config.Section{ID: "e", SortOrder: 1},
	}
	for i, f := range *(ss.Sort()) {
		assert.EqualValues(t, want[i], f.SortOrder)
	}

}

func TestSectionSliceSortAll(t *testing.T) {
	want := `[{"ID":"b","SortOrder":-10,"Groups":null},{"ID":"e","SortOrder":1,"Groups":null},{"ID":"c","SortOrder":10,"Groups":null},{"ID":"a","SortOrder":20,"Groups":[{"ID":"b","SortOrder":-10,"Fields":[{"ID":"b","SortOrder":-10},{"ID":"e","SortOrder":1},{"ID":"c","SortOrder":10},{"ID":"d","SortOrder":11},{"ID":"a","SortOrder":20}]},{"ID":"e","SortOrder":1,"Fields":null},{"ID":"d","SortOrder":11,"Fields":[{"ID":"b","SortOrder":-10},{"ID":"e","SortOrder":1},{"ID":"c","SortOrder":10},{"ID":"d","SortOrder":11},{"ID":"a","SortOrder":20}]},{"ID":"a","SortOrder":20,"Fields":[{"ID":"b","SortOrder":-10},{"ID":"e","SortOrder":1},{"ID":"c","SortOrder":10},{"ID":"d","SortOrder":11},{"ID":"a","SortOrder":20}]}]}]` + "\n"
	ss := config.SectionSlice{
		&config.Section{ID: "a", SortOrder: 20, Groups: config.GroupSlice{
			&config.Group{ID: "a", SortOrder: 20, Fields: config.FieldSlice{&config.Field{ID: "a", SortOrder: 20}, &config.Field{ID: "b", SortOrder: -10}, &config.Field{ID: "c", SortOrder: 10}, &config.Field{ID: "d", SortOrder: 11}, &config.Field{ID: "e", SortOrder: 1}}},
			&config.Group{ID: "b", SortOrder: -10, Fields: config.FieldSlice{&config.Field{ID: "a", SortOrder: 20}, &config.Field{ID: "b", SortOrder: -10}, &config.Field{ID: "c", SortOrder: 10}, &config.Field{ID: "d", SortOrder: 11}, &config.Field{ID: "e", SortOrder: 1}}},
			&config.Group{ID: "d", SortOrder: 11, Fields: config.FieldSlice{&config.Field{ID: "a", SortOrder: 20}, &config.Field{ID: "b", SortOrder: -10}, &config.Field{ID: "c", SortOrder: 10}, &config.Field{ID: "d", SortOrder: 11}, &config.Field{ID: "e", SortOrder: 1}}},
			&config.Group{ID: "e", SortOrder: 1},
		}},
		&config.Section{ID: "b", SortOrder: -10},
		&config.Section{ID: "c", SortOrder: 10},
		&config.Section{ID: "e", SortOrder: 1},
	}
	ss.SortAll()
	have := ss.ToJson()
	if want != have {
		t.Errorf("\nWant: %s\nHave: %s\n", want, have)
	}
}
