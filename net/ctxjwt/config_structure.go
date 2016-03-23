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

package ctxjwt

import (
	"github.com/corestoreio/csfw/config/cfgpath"
	"github.com/corestoreio/csfw/config/element"
	"github.com/corestoreio/csfw/storage/text"
	"github.com/corestoreio/csfw/store/scope"
)

// DefaultSigningMethod HMAC-SHA signing with 512 bits. Gets applied if the
// ConfigSigningMethod model returns an empty string.
const DefaultSigningMethod = "HS512"

// NewConfigStructure global configuration structure for this package.
// Used in frontend (to display the user all the settings) and in
// backend (scope checks and default values). See the source code
// of this function for the overall available sections, groups and fields.
func NewConfigStructure() (element.SectionSlice, error) {
	return element.NewConfiguration(
		&element.Section{
			ID: cfgpath.NewRoute("net"),
			Groups: element.NewGroupSlice(
				&element.Group{
					ID:        cfgpath.NewRoute("ctxjwt"),
					Label:     text.Chars(`JSON Web Token (JWT)`),
					SortOrder: 40,
					Scopes:    scope.PermWebsite,
					Fields: element.NewFieldSlice(
						&element.Field{
							// Path: net/ctxjwt/signing_method
							ID:        cfgpath.NewRoute("signing_method"),
							Label:     text.Chars(`Token Signing Algorithm`),
							Type:      element.TypeSelect,
							SortOrder: 10,
							Visible:   element.VisibleYes,
							Scopes:    scope.PermWebsite,
							Default:   DefaultSigningMethod,
						},
						&element.Field{
							// Path: net/ctxjwt/hmac_password
							ID:        cfgpath.NewRoute("hmac_password"),
							Label:     text.Chars(`HMAC Token Password`),
							Type:      element.TypeObscure,
							SortOrder: 10,
							Visible:   element.VisibleYes,
							Scopes:    scope.PermWebsite,
						},
						&element.Field{
							// Path: net/ctxjwt/rsa_key
							ID:        cfgpath.NewRoute("rsa_key"),
							Label:     text.Chars(`Private RSA Key`),
							Type:      element.TypeObscure,
							SortOrder: 20,
							Visible:   element.VisibleYes,
							Scopes:    scope.PermWebsite,
						},
						&element.Field{
							// Path: net/ctxjwt/rsa_key_password
							ID:        cfgpath.NewRoute("rsa_key_password"),
							Label:     text.Chars(`Private RSA Key Password`),
							Comment:   text.Chars(`If the key has been secured via a password, provide it here.`),
							Type:      element.TypeObscure,
							SortOrder: 30,
							Visible:   element.VisibleYes,
							Scopes:    scope.PermWebsite,
						},
						&element.Field{
							// Path: net/ctxjwt/ecdsa_key
							ID:        cfgpath.NewRoute("ecdsa_key"),
							Label:     text.Chars(`Private ECDSA Key`),
							Comment:   text.Chars(`Elliptic Curve Digital Signature Algorithm, as defined in FIPS 186-3.`),
							Type:      element.TypeObscure,
							SortOrder: 40,
							Visible:   element.VisibleYes,
							Scopes:    scope.PermWebsite,
						},
						&element.Field{
							// Path: net/ctxjwt/ecdsa_key_password
							ID:        cfgpath.NewRoute("ecdsa_key_password"),
							Label:     text.Chars(`Private ECDSA Key Password`),
							Comment:   text.Chars(`If the key has been secured via a password, provide it here.`),
							Type:      element.TypeObscure,
							SortOrder: 50,
							Visible:   element.VisibleYes,
							Scopes:    scope.PermWebsite,
						},
					),
				},
			),
		},
	)
}
