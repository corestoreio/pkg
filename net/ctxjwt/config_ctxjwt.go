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
	"github.com/corestoreio/csfw/config/element"
	"github.com/corestoreio/csfw/config/path"
	"github.com/corestoreio/csfw/storage/text"
	"github.com/corestoreio/csfw/store/scope"
)

const PathJWTHMACPassword = `corestore/jwt/hmac_password`

// ConfigStructure global configuration structure for this package.
// Used in frontend and backend. See init() for details.
var ConfigStructure element.SectionSlice

func init() {
	ConfigStructure = element.MustNewConfiguration(
		&element.Section{
			ID: path.NewRoute("corestore"),
			Groups: element.NewGroupSlice(
				&element.Group{
					ID:        path.NewRoute("jwt"),
					Label:     text.Chars(`JSON Web Token (JWT)`),
					SortOrder: 40,
					Scopes:    scope.PermStore,
					Fields: element.NewFieldSlice(
						&element.Field{
							// Path: corestore/jwt/hmac_password
							ID:        path.NewRoute("hmac_password"),
							Label:     text.Chars(`Token Password`),
							Type:      element.TypeObscure,
							SortOrder: 10,
							Visible:   element.VisibleYes,
							Scopes:    scope.PermStore,
						},
					),
				},
			),
		},
	)
}
