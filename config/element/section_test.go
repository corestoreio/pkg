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

package element_test

import (
	"testing"

	"github.com/corestoreio/errors"
	"github.com/corestoreio/pkg/config"
	"github.com/corestoreio/pkg/config/cfgpath"
	"github.com/corestoreio/pkg/config/element"
	"github.com/corestoreio/pkg/storage/text"
	"github.com/stretchr/testify/assert"
)

var _ element.Sectioner = (*element.Sections)(nil)

func TestSectionValidateDuplicate(t *testing.T) {
	// for benchmark tests see package config_bm

	ss := element.MakeSections(
		element.Section{
			ID: cfgpath.MakeRoute(`aa`),
			Groups: element.MakeGroups(
				element.Group{
					ID: cfgpath.MakeRoute(`bb`),
					Fields: element.MakeFields(
						element.Field{ID: cfgpath.MakeRoute(`cc`)},
						element.Field{ID: cfgpath.MakeRoute(`cc`)},
					),
				},
			),
		},
	)
	assert.True(t, errors.IsNotValid(ss.Validate())) // "Duplicate entry for path aa/bb/cc :: [{\"ID\":\"aa\",\"Groups\":[{\"ID\":\"bb\",\"Fields\":[{\"ID\":\"cc\"},{\"ID\":\"cc\"}]}]}]\n"
}

func TestSectionValidateShortPath(t *testing.T) {

	ss := element.MakeSections(
		element.Section{
			//ID: cfgpath.MakeRoute(`aa`),
			Groups: element.MakeGroups(
				element.Group{
					//ID: cfgpath.MakeRoute(`b`),
					Fields: element.MakeFields(
						element.Field{ID: cfgpath.MakeRoute(`ca`)},
						element.Field{ID: cfgpath.MakeRoute(`cb`)},
						element.Field{},
					),
				},
			),
		},
	)

	err := ss.Validate()
	assert.True(t, errors.IsEmpty(err), "Error %s", err)
}

func TestSectionUpdateField(t *testing.T) {

	ss := element.MakeSections(
		element.Section{
			ID: cfgpath.MakeRoute(`aa`),
			Groups: element.MakeGroups(
				element.Group{
					ID: cfgpath.MakeRoute(`bb`),
					Fields: element.MakeFields(
						element.Field{ID: cfgpath.MakeRoute(`ca`)},
						element.Field{ID: cfgpath.MakeRoute(`cb`)},
					),
				},
			),
		},
	)

	fr := cfgpath.MakeRoute(`aa/bb/ca`)
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

	err1 := ss.UpdateField(cfgpath.MakeRoute(`a/b/c`), element.Field{})
	assert.True(t, errors.NotFound.Match(err1), "Error: %s", err1)

	err2 := ss.UpdateField(cfgpath.MakeRoute(`aa/b/c`), element.Field{})
	assert.True(t, errors.NotFound.Match(err2), "Error: %s", err2)

	err3 := ss.UpdateField(cfgpath.MakeRoute(`aa/bb/c`), element.Field{})
	assert.True(t, errors.NotFound.Match(err3), "Error: %s", err3)

	err4 := ss.UpdateField(cfgpath.MakeRoute(`aa_bb_c`), element.Field{})
	assert.True(t, errors.IsNotValid(err4), "Error: %s", err4)

}

var _ element.ConfigurationWriter = (*config.Service)(nil)
var _ config.Writer = (*config.Service)(nil)

func TestService_ApplyDefaults(t *testing.T) {

	pkgCfg := element.MustMakeSectionsValidate(
		element.Section{
			ID: cfgpath.MakeRoute("contact"),
			Groups: element.MakeGroups(
				element.Group{
					ID: cfgpath.MakeRoute("contact"),
					Fields: element.MakeFields(
						element.Field{
							// Path: `contact/contact/enabled`,
							ID:      cfgpath.MakeRoute("enabled"),
							Default: true,
						},
					),
				},
				element.Group{
					ID: cfgpath.MakeRoute("email"),
					Fields: element.MakeFields(
						element.Field{
							// Path: `contact/email/recipient_email`,
							ID:      cfgpath.MakeRoute("recipient_email"),
							Default: `hello@example.com`,
						},
						element.Field{
							// Path: `contact/email/sender_email_identity`,
							ID:      cfgpath.MakeRoute("sender_email_identity"),
							Default: 2.7182818284590452353602874713527,
						},
						element.Field{
							// Path: `contact/email/email_template`,
							ID:      cfgpath.MakeRoute("email_template"),
							Default: 4711,
						},
					),
				},
			),
		},
	)
	s := config.MustNewService(config.NewInMemoryStore())
	if _, err := pkgCfg.ApplyDefaults(s); err != nil {
		t.Fatal(err)
	}
	cer, _, err := pkgCfg.FindField(cfgpath.MakeRoute("contact", "email", "recipient_email"))
	if err != nil {
		t.Fatal(err)
	}
	email, err := s.String(cfgpath.MustMakeByString("contact/email/recipient_email")) // default scope
	assert.NoError(t, err)
	assert.Exactly(t, cer.Default.(string), email)
	assert.NoError(t, s.Close())
}
