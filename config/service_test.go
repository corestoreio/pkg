// Copyright 2015, Cyrill @ Schumacher.fm and the CoreStore contributors
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
	"github.com/corestoreio/csfw/config/element"
	"github.com/corestoreio/csfw/storage/csdb"
	"github.com/stretchr/testify/assert"
)

func init() {
	dbc := csdb.MustConnectTest()
	if err := config.TableCollection.Init(dbc.NewSession()); err != nil {
		panic(err)
	}
	if err := dbc.Close(); err != nil {
		panic(err)
	}
}

func TestScopeApplyDefaults(t *testing.T) {
	defer debugLogBuf.Reset()
	defer infoLogBuf.Reset()

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
	s := config.NewService()
	s.ApplyDefaults(pkgCfg)
	cer, err := pkgCfg.FindFieldByPath("contact", "email", "recipient_email")
	if err != nil {
		t.Error(err)
		return
	}
	sval, err := s.String(config.Path("contact/email/recipient_email"))
	assert.NoError(t, err)
	assert.Exactly(t, cer.Default.(string), sval)
	assert.NoError(t, s.Close())
}

// TestApplyCoreConfigData reads from the MySQL core_config_data table and applies
// these value to the underlying storage. tries to get back the values from the
// underlying storage
func TestApplyCoreConfigData(t *testing.T) {
	defer debugLogBuf.Reset()
	defer infoLogBuf.Reset()

	dbc := csdb.MustConnectTest()
	defer func() { assert.NoError(t, dbc.Close()) }()
	sess := dbc.NewSession(nil) // nil tricks the NewSession ;-)

	s := config.NewService()
	defer func() { assert.NoError(t, s.Close()) }()

	loadedRows, writtenRows, err := s.ApplyCoreConfigData(sess)
	if err != nil {
		t.Fatal(err)
	}
	assert.True(t, loadedRows > 9, "loadedRows %d", loadedRows)
	assert.True(t, writtenRows > 9, "writtenRows %d", writtenRows)

	//	println("\n", debugLogBuf.String(), "\n")
	//	println("\n", infoLogBuf.String(), "\n")

	assert.NoError(t, s.Write(config.Path("web/secure/offloader_header"), config.ScopeDefault(), config.Value("SSL_OFFLOADED")))

	h, err := s.String(config.Path("web/secure/offloader_header"), config.ScopeDefault())
	assert.NoError(t, err)
	assert.Exactly(t, "SSL_OFFLOADED", h)

	assert.Len(t, s.Storage.AllKeys(), writtenRows)
}
