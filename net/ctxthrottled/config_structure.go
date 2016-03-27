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

package ctxthrottled

import (
	"github.com/corestoreio/csfw/config/cfgpath"
	"github.com/corestoreio/csfw/config/element"
	"github.com/corestoreio/csfw/storage/text"
	"github.com/corestoreio/csfw/store/scope"
)

// NewConfigStructure global configuration structure for this package.
// Used in frontend (to display the user all the settings) and in
// backend (scope checks and default values). See the source code
// of this function for the overall available sections, groups and fields.
func NewConfigStructure() (element.SectionSlice, error) {
	return element.NewConfiguration(
		element.Section{
			ID: cfgpath.NewRoute("net"),
			Groups: element.NewGroupSlice(
				element.Group{
					ID:        cfgpath.NewRoute("ctxthrottled"),
					Label:     text.Chars(`Rate throtteling`),
					SortOrder: 40,
					Scopes:    scope.PermWebsite,
					Fields: element.NewFieldSlice(
						element.Field{
							// Path: net/ctxthrottled/burst
							ID:        cfgpath.NewRoute("burst"),
							Label:     text.Chars(`Burst`),
							Comment:   text.Chars(`Defines the number of requests that will be allowed to exceed the rate in a single burst and must be greater than or equal to zero`),
							Type:      element.TypeText,
							SortOrder: 10,
							Visible:   element.VisibleYes,
							Scopes:    scope.PermWebsite,
							Default:   5,
						},
						element.Field{
							// Path: net/ctxthrottled/requests
							ID:        cfgpath.NewRoute("requests"),
							Label:     text.Chars(`Requests`),
							Comment:   text.Chars(`Number of requests allowed per time period`),
							Type:      element.TypeText,
							SortOrder: 20,
							Visible:   element.VisibleYes,
							Scopes:    scope.PermWebsite,
							Default:   100,
						},
						element.Field{
							// Path: net/ctxthrottled/duration
							ID:        cfgpath.NewRoute("duration"),
							Label:     text.Chars(`Duration`),
							Comment:   text.Chars(`Per second (s), minute (i), hour (h) or day (d)`),
							Type:      element.TypeText,
							SortOrder: 30,
							Visible:   element.VisibleYes,
							Scopes:    scope.PermWebsite,
							Default:   `h`,
						},
					),
				},
			),
		},
	)
}
