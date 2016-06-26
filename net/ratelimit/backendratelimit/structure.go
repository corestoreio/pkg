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

package backendratelimit

import (
	"github.com/corestoreio/csfw/config/cfgpath"
	"github.com/corestoreio/csfw/config/element"
	"github.com/corestoreio/csfw/storage/text"
	"github.com/corestoreio/csfw/store/scope"
)

// NewConfigStructure global configuration structure for this package. Used in
// frontend (to display the user all the settings) and in backend (scope checks
// and default values). See the source code of this function for the overall
// available sections, groups and fields.
func NewConfigStructure() (element.SectionSlice, error) {
	sortIdx := 10
	var iter = func() int {
		sortIdx += 10
		return sortIdx
	}
	return element.NewConfiguration(
		element.Section{
			ID: cfgpath.NewRoute("net"),
			Groups: element.NewGroupSlice(
				element.Group{
					ID:        cfgpath.NewRoute("ratelimit"),
					Label:     text.Chars(`Rate throtteling`),
					SortOrder: 130,
					Scopes:    scope.PermStore,
					Fields: element.NewFieldSlice(
						element.Field{
							// Path: net/ratelimit/disabled
							ID:        cfgpath.NewRoute("disabled"),
							Label:     text.Chars(`Disabled`),
							Comment:   text.Chars(`Set to true to disable rate limiting.`),
							Type:      element.TypeSelect,
							SortOrder: iter(),
							Visible:   element.VisibleYes,
							Scopes:    scope.PermStore,
						},
						element.Field{
							// Path: net/ratelimit/burst
							ID:        cfgpath.NewRoute("burst"),
							Label:     text.Chars(`Burst`),
							Comment:   text.Chars(`Defines the number of requests that will be allowed to exceed the rate in a single burst and must be greater than or equal to zero`),
							Type:      element.TypeText,
							SortOrder: iter(),
							Visible:   element.VisibleYes,
							Scopes:    scope.PermStore,
							Default:   20,
						},
						element.Field{
							// Path: net/ratelimit/requests
							ID:        cfgpath.NewRoute("requests"),
							Label:     text.Chars(`Requests`),
							Comment:   text.Chars(`Number of requests allowed per time period`),
							Type:      element.TypeText,
							SortOrder: iter(),
							Visible:   element.VisibleYes,
							Scopes:    scope.PermStore,
							Default:   100,
						},
						element.Field{
							// Path: net/ratelimit/duration
							ID:        cfgpath.NewRoute("duration"),
							Label:     text.Chars(`Duration`),
							Comment:   text.Chars(`Per second (s), minute (i), hour (h) or day (d)`),
							Type:      element.TypeText,
							SortOrder: iter(),
							Visible:   element.VisibleYes,
							Scopes:    scope.PermStore,
							Default:   `h`,
						},
					),
				},
				element.Group{
					ID:        cfgpath.NewRoute("ratelimit_storage"),
					Label:     text.Chars(`Rate throtteling storage`),
					SortOrder: 140,
					Scopes:    scope.PermStore,
					Fields: element.NewFieldSlice(

						element.Field{
							// Path: net/ratelimit_storage/gcra_memory_enable
							ID:        cfgpath.NewRoute("gcra_memory_enable"),
							Label:     text.Chars(`Enable GCRA in-memory key storage`),
							Comment:   text.Chars(`Enables the in memory storage`),
							Type:      element.TypeSelect,
							SortOrder: iter(),
							Visible:   element.VisibleYes,
							Scopes:    scope.PermStore,
						},
						element.Field{
							// Path: net/ratelimit_storage/gcra_max_memory_keys
							ID:        cfgpath.NewRoute("gcra_max_memory_keys"),
							Label:     text.Chars(`GCRA max memory keys`),
							Comment:   text.Chars(`If maxKeys > 0, the number of different keys is restricted to the specified amount. In this case, it uses an LRU algorithm to evict older keys to make room for newer ones. If maxKeys <= 0, there is no limit on the number of keys, which may use an unbounded amount of memory.`),
							Type:      element.TypeText,
							SortOrder: iter(),
							Visible:   element.VisibleYes,
							Scopes:    scope.PermStore,
							Default:   65536,
						},
					),
				},
			),
		},
	)
}
