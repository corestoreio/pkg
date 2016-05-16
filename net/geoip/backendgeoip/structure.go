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

package backendgeoip

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
			ID: cfgpath.NewRoute(`net`),
			Groups: element.NewGroupSlice(
				element.Group{
					ID:    cfgpath.NewRoute(`geoip`),
					Label: text.Chars(`Geo IP`),
					Comment: text.Chars(`Detects the country by an IP address and maybe restricts the access. Compatible to IPv4 and IPv6.
Uses the maxmind database or alternative country/city detectors.`),
					MoreURL:   text.Chars(`https://www.maxmind.com/en/geoip2-services-and-databases`),
					SortOrder: 170,
					Scopes:    scope.PermWebsite,
					Fields: element.NewFieldSlice(
						element.Field{
							// Path: `net/geoip/allowed_countries`,
							ID:    cfgpath.NewRoute(`allowed_countries`),
							Label: text.Chars(`Allowed countries`),
							Comment: text.Chars(`Indicates which headers are safe to
expose to the API of a CORS API specification. Separate via line break (\n)`),
							Type:      element.TypeSelect,
							SortOrder: 20,
							Visible:   element.VisibleYes,
							Scopes:    scope.PermStore,
						},
					),
				},
			),
		},
	)
}
