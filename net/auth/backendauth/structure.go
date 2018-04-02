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

package backendauth

import (
	"github.com/corestoreio/pkg/config/cfgpath"
	"github.com/corestoreio/pkg/config/element"
	"github.com/corestoreio/pkg/storage/text"
	"github.com/corestoreio/pkg/store/scope"
)

// NewConfigStructure global configuration structure for this package.
// Used in frontend (to display the user all the settings) and in
// backend (scope checks and default values). See the source code
// of this function for the overall available sections, groups and fields.
func NewConfigStructure() (element.Sections, error) {
	return element.MakeSectionsValidated(
		element.Section{
			ID: cfgpath.MakeRoute(`net`),
			Groups: element.MakeGroups(
				element.Group{
					ID:        cfgpath.MakeRoute(`auth`),
					Label:     text.Chars(`Authentication`),
					Comment:   text.Chars(`TODO`),
					SortOrder: 160,
					Scopes:    scope.PermWebsite,
					Fields: element.MakeFields(
						element.Field{
							// Path: net/auth/disabled
							ID:        cfgpath.MakeRoute("disabled"),
							Label:     text.Chars(`Authentication is disabled`),
							Comment:   text.Chars(`Disables completely the authentication. Set to true/enable to activate the disabling.`),
							Type:      element.TypeSelect,
							SortOrder: 10,
							Visible:   element.VisibleYes,
							Scopes:    scope.PermWebsite,
							Default:   false,
						},
					),
				},
			),
		},
	)
}
