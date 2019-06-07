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

package backendsigned

import (
	"github.com/corestoreio/pkg/config/cfgpath"
	"github.com/corestoreio/pkg/config/element"
	"github.com/corestoreio/pkg/storage/text"
	"github.com/corestoreio/pkg/store/scope"
)

// NewConfigStructure global configuration structure for this package. Used in
// frontend (to display the user all the settings) and in backend (scope checks
// and default values). See the source code of this function for the overall
// available sections, groups and fields.
func NewConfigStructure() (element.Sections, error) {
	sortIdx := 10
	var iter = func() int {
		sortIdx += 10
		return sortIdx
	}
	return element.MakeSectionsValidated(
		element.Section{
			ID: cfgpath.MakeRoute("net"),
			Groups: element.MakeGroups(
				element.Group{
					ID:        cfgpath.MakeRoute("signed"),
					Label:     text.Chars(`Rate throtteling`),
					SortOrder: 130,
					Scopes:    scope.PermStore,
					Fields: element.MakeFields(
						element.Field{
							// Path: net/signed/disabled
							ID:        cfgpath.MakeRoute("disabled"),
							Label:     text.Chars(`Disabled`),
							Comment:   text.Chars(`Set to true to disable all signed middlewares.`),
							Type:      element.TypeSelect,
							SortOrder: iter(),
							Visible:   element.VisibleYes,
							Scopes:    scope.PermStore,
						},
						element.Field{
							// Path: net/signed/in_trailer
							ID:        cfgpath.MakeRoute("in_trailer"),
							Label:     text.Chars(`In Trailer`),
							Comment:   text.Chars(`If true uses a stream based approach to calculate the hash and appends the hash to the HTTP trailer. Note: not all clients can read a HTTP trailer.`),
							Type:      element.TypeSelect,
							SortOrder: iter(),
							Visible:   element.VisibleYes,
							Scopes:    scope.PermStore,
							Default:   true,
						},
						element.Field{
							// Path: net/signed/allowed_methods
							ID:        cfgpath.MakeRoute("allowed_methods"),
							Label:     text.Chars(`Allowed HTTP methods`),
							Comment:   text.Chars(`Limit the validation middleware to the listed HTTP methods to verify a hash.`),
							Type:      element.TypeSelect,
							SortOrder: iter(),
							Visible:   element.VisibleYes,
							Scopes:    scope.PermStore,
							Default:   `POST,PUT,PATCH`,
						},
						element.Field{
							// Path: net/signed/key
							ID:        cfgpath.MakeRoute("key"),
							Label:     text.Chars(`Key / Password`),
							Comment:   text.Chars(`Your key or password to calculate the hash value with an HMAC function. Longer and weirder keys are winners.`),
							Type:      element.TypeText,
							SortOrder: iter(),
							Visible:   element.VisibleYes,
							Scopes:    scope.PermStore,
						},
						element.Field{
							// Path: net/signed/algorithm
							ID:        cfgpath.MakeRoute("algorithm"),
							Label:     text.Chars(`Algorithm`),
							Comment:   text.Chars(`Currently supported algorithms: sha256, sha512`),
							Type:      element.TypeSelect,
							SortOrder: iter(),
							Visible:   element.VisibleYes,
							Scopes:    scope.PermStore,
							Default:   `sha256`,
						},
						element.Field{
							// Path: net/signed/http_header_type
							ID:        cfgpath.MakeRoute("http_header_type"),
							Label:     text.Chars(`HTTP Header Type`),
							Comment:   text.Chars(`Sets the type of the HTTP header to either Content-HMAC or Content-Signature.`),
							Type:      element.TypeSelect,
							SortOrder: iter(),
							Visible:   element.VisibleYes,
							Scopes:    scope.PermStore,
							Default:   `hmac`,
						},
						element.Field{
							// Path: net/signed/key_id
							ID:        cfgpath.MakeRoute("key_id"),
							Label:     text.Chars(`Key ID`),
							Comment:   text.Chars(`KeyID name or ID of the key which is used in the HMAC algorithm. Only usable when HTTPHeaderType has been set to "signature"`),
							Type:      element.TypeText,
							SortOrder: iter(),
							Visible:   element.VisibleYes,
							Scopes:    scope.PermStore,
						},
					),
				},
				element.Group{
					ID:        cfgpath.MakeRoute("signed_algorithm"),
					Label:     text.Chars(`Signature algortihm`),
					SortOrder: 140,
					Scopes:    scope.PermStore,
					Fields: element.MakeFields(
						element.Field{
							// Path: net/signed_storage/gcra_name
							ID:        cfgpath.MakeRoute("gcra_name"),
							Label:     text.Chars(`Name of the registered GCRA`),
							Comment:   text.Chars(`Insert the name of the registered GCRA during program initialization with the function Backend.Register().`),
							Type:      element.TypeText,
							SortOrder: iter(),
							Visible:   element.VisibleYes,
							Scopes:    scope.PermStore,
						},
						element.Field{
							// Path: net/signed_storage/enable_gcra_memory
							ID:        cfgpath.MakeRoute("enable_gcra_memory"),
							Label:     text.Chars(`Use GCRA in-memory (max keys)`),
							Comment:   text.Chars(`If maxKeys > 0 in-memory key storage will be enabled. The max keys  number of different keys is restricted to the specified amount (65536). In this case, it uses an LRU algorithm to evict older keys to make room for newer ones.`),
							Type:      element.TypeText,
							SortOrder: iter(),
							Visible:   element.VisibleYes,
							Scopes:    scope.PermStore,
							Default:   0,
						},
						element.Field{
							// Path: net/signed_storage/enable_gcra_redis
							ID:        cfgpath.MakeRoute("enable_gcra_redis"),
							Label:     text.Chars(`Use GCRA Redis`),
							Comment:   text.Chars(`If a Redis URL is provided a Redis server will be used for key storage. Setting both entries (in-memory and Redis) then only Redis will be applied. URLs should follow the draft IANA specification for the scheme (https://www.iana.org/assignments/uri-schemes/prov/redis). For example: redis://localhost:6379/3 |  redis://:6380/0 => connects to localhost:6380 | redis:// => connects to localhost:6379 with DB 0 | redis://empty:myPassword@clusterName.xxxxxx.0001.usw2.cache.amazonaws.com:6379/0`),
							Type:      element.TypeText,
							SortOrder: iter(),
							Visible:   element.VisibleYes,
							Scopes:    scope.PermStore,
						},
					),
				},
			),
		},
	)
}
