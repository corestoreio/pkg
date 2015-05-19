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

	"github.com/Sirupsen/logrus"
	"github.com/corestoreio/csfw/config"
	"github.com/stretchr/testify/assert"
)

func init() {
	l := logrus.New()
	l.Level = logrus.DebugLevel
	config.SetLogger(l)
}

func TestNewConfiguration(t *testing.T) {
	tests := []struct {
		have    []*config.Section
		wantErr string
	}{
		{
			have:    []*config.Section{},
			wantErr: "SectionSlice is empty",
		},
		{
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
		},
		{
			have:    []*config.Section{&config.Section{ID: "a", Groups: config.GroupSlice{}}},
			wantErr: "",
		},
		{
			have:    []*config.Section{&config.Section{ID: "a", Groups: config.GroupSlice{&config.Group{ID: "b", Fields: nil}}}},
			wantErr: "",
		},
		{
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
			wantErr: "Duplicate entry for path a/b/c",
		},
	}

	for _, test := range tests {
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
				assert.Nil(t, haveSlice)
			} else {
				assert.NotNil(t, haveSlice)
				assert.Len(t, haveSlice, len(have))
			}
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
		config.DefaultMap{"contact/email/recipient_email": "hello@example.com", "contact/email/sender_email_identity": 2.718281828459045, "contact/email/email_template": 4711, "contact/contact/enabled": true},
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
	}{
		0: {
			have: []config.SectionSlice{
				config.SectionSlice{
					&config.Section{
						ID: "a",
						Groups: config.GroupSlice{
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
			want:    `[{"ID":"a","Label":"LabelA","Permission":0,"Groups":[{"ID":"b","Fields":[{"ID":"c","Default":"c"}]},{"ID":"b","Fields":[{"ID":"d","Default":"d"}]}]}]` + "\n",
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
								Scope: config.NewScopePerm(config.ScopeDefault),
								Fields: config.FieldSlice{
									&config.Field{ID: "c", Default: `c`},
								},
							},
						},
					},
				},
				config.SectionSlice{
					&config.Section{
						ID:    "a",
						Scope: config.NewScopePerm(config.ScopeDefault, config.ScopeWebsite),
						Groups: config.GroupSlice{
							&config.Group{ID: "b", Label: "GroupLabelB1"},
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
			want:    `[{"ID":"a","Label":"SectionLabelA","Scope":["ScopeDefault","ScopeWebsite"],"Permission":0,"Groups":[{"ID":"b","Label":"GroupLabelB2","Scope":["ScopeDefault"],"Fields":[{"ID":"c","Default":"c"}]},{"ID":"b2","Fields":[{"ID":"d","Default":"d"}]}]}]` + "\n",
		},
		2: {
			have: []config.SectionSlice{
				config.SectionSlice{
					&config.Section{ID: "a", Label: "SectionLabelA", SortOrder: 20, Permission: 22},
				},
				config.SectionSlice{
					&config.Section{ID: "a", Scope: config.NewScopePerm(config.ScopeDefault, config.ScopeWebsite), SortOrder: 10, Permission: 3},
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
								Scope:   config.NewScopePerm(config.ScopeDefault),
							},
						},
					},
				},
				config.SectionSlice{
					&config.Section{
						ID:        "a",
						SortOrder: 1000,
						Scope:     config.NewScopePerm(config.ScopeDefault, config.ScopeWebsite),
						Groups: config.GroupSlice{
							&config.Group{ID: "b", Label: "GroupLabelB1", Scope: config.ScopePermAll},
							&config.Group{ID: "b", Label: "GroupLabelB2", Comment: "Section2AGroup3BComment", SortOrder: 100},
							&config.Group{ID: "b2"},
						},
					},
				},
			},
			wantErr: "",
			want:    `[{"ID":"a","Label":"SectionLabelA","Scope":["ScopeDefault","ScopeWebsite"],"SortOrder":1000,"Permission":0,"Groups":[{"ID":"b","Label":"GroupLabelB2","Comment":"Section2AGroup3BComment","Scope":["ScopeDefault","ScopeWebsite","ScopeStore"],"SortOrder":100,"Fields":null},{"ID":"b2","Fields":null}]}]` + "\n",
		},
		4: {
			have: []config.SectionSlice{
				config.SectionSlice{
					&config.Section{
						ID: "a",
						Groups: config.GroupSlice{
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
									&config.Field{ID: "c", Default: `haha`, Type: config.TypeSelect, Scope: config.NewScopePerm(config.ScopeDefault, config.ScopeWebsite)},
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
								ID: "b",
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
			want:    `[{"ID":"a","Permission":0,"Groups":[{"ID":"b","Fields":[{"ID":"c","Type":"select","Default":"overriddenHaha"},{"ID":"d","Label":"Sect2Group2Label4","Comment":"LOTR","Default":"overriddenD"}]},{"ID":"b","Fields":[{"ID":"d","Type":"obscure","Comment":"Ring of fire","Default":"d"},{"ID":"c","Type":"select","Scope":["ScopeDefault","ScopeWebsite"],"Default":"haha"}]}]}]` + "\n",
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
					&config.Section{
						ID: "a",
						Groups: config.GroupSlice{
							&config.Group{
								ID: "b",
								Fields: config.FieldSlice{
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
			want:    `[{"ID":"a","Permission":0,"Groups":[{"ID":"b","Fields":[{"ID":"c","Type":"select","Label":"Sect2Group2Label4","Comment":"LOTR","SortOrder":100,"Visible":true,"Default":"overridenC"}]}]}]` + "\n",
		},
	}

	for i, test := range tests {
		var baseSl config.SectionSlice
		haveErr := baseSl.MergeAll(test.have...)
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

	}
}
