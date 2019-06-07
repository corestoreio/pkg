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

package backendgeoip

import (
	"time"

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
	return element.MakeSectionsValidated(
		element.Section{
			ID: cfgpath.MakeRoute(`net`),
			Groups: element.MakeGroups(
				element.Group{
					ID:    cfgpath.MakeRoute(`geoip`),
					Label: text.Chars(`Geo IP`),
					Comment: text.Chars(`Detects the country by an IP address and maybe restricts the access. Compatible
to IPv4 and IPv6.`),
					SortOrder: 170,
					Scopes:    scope.PermStore,
					Fields: element.MakeFields(
						element.Field{
							// Path: `net/geoip/allowed_countries`,
							ID:    cfgpath.MakeRoute(`allowed_countries`),
							Label: text.Chars(`Allowed countries`),
							Comment: text.Chars(`Defines a list of ISO country codes which are allowed. Separated via comma,
e.g.: DE,CH,AT,AU,NZ`),
							Type:      element.TypeSelect,
							SortOrder: 20,
							Visible:   element.VisibleYes,
							Scopes:    scope.PermStore,
						},
						element.Field{
							// Path: `net/geoip/alternative_redirect`,
							ID:        cfgpath.MakeRoute(`alternative_redirect`),
							Label:     text.Chars(`Alternative Redirect URL`),
							Comment:   text.Chars(`Redirects the client to this URL if their country doesn't have access.`),
							Type:      element.TypeText,
							SortOrder: 30,
							Visible:   element.VisibleYes,
							Scopes:    scope.PermStore,
						},
						element.Field{
							// Path: `net/geoip/alternative_redirect_code`,
							ID:        cfgpath.MakeRoute(`alternative_redirect_code`),
							Label:     text.Chars(`Alternative Redirect HTTP Code`),
							Comment:   text.Chars(`Specifies the HTTP redirect code`),
							Type:      element.TypeSelect,
							SortOrder: 40,
							Visible:   element.VisibleYes,
							Scopes:    scope.PermStore,
							Default:   301,
						},
					),
				},

				element.Group{
					ID:    cfgpath.MakeRoute(`geoip_maxmind`),
					Label: text.Chars(`Geo IP (MaxMind)`),
					Comment: text.Chars(`Detects the country by an IP address and maybe restricts the access. Compatible
to IPv4 and IPv6. Uses the maxmind database from a file or the web service.`),
					MoreURL:   text.Chars(`https://www.maxmind.com/en/geoip2-services-and-databases`),
					HelpURL:   text.Chars(`http://dev.maxmind.com/geoip/geoip2/web-services/`),
					SortOrder: 170,
					Scopes:    scope.PermDefault,
					Fields: element.MakeFields(

						element.Field{
							// Path: `net/geoip_maxmind/data_source`,
							ID:    cfgpath.MakeRoute(`data_source`),
							Label: text.Chars(`Source geo location data`),
							Comment: text.Chars(`Choose from which source you would like to load the MaxMind geo location data.
Either from a "file" or from a "webservice".`),
							Type:      element.TypeSelect,
							SortOrder: 5,
							Visible:   element.VisibleYes,
							Scopes:    scope.PermDefault,
						},

						element.Field{
							// Path: `net/geoip_maxmind/local_file`,
							ID:    cfgpath.MakeRoute(`local_file`),
							Label: text.Chars(`Local MaxMind database file`),
							Comment: text.Chars(`Load a local MaxMind binary database file for extracting country information
from an IP address.`),
							Type:      element.TypeSelect,
							SortOrder: 10,
							Visible:   element.VisibleYes,
							Scopes:    scope.PermDefault,
						},

						element.Field{
							// Path: `net/geoip_maxmind/webservice_userid`,
							ID:    cfgpath.MakeRoute(`webservice_userid`),
							Label: text.Chars(`Webservice User ID`),
							//Comment:   text.Chars(``),
							Type:      element.TypeText,
							SortOrder: 30,
							Visible:   element.VisibleYes,
							Scopes:    scope.PermDefault,
						},
						element.Field{
							// Path: `net/geoip_maxmind/webservice_license`,
							ID:    cfgpath.MakeRoute(`webservice_license`),
							Label: text.Chars(`Webservice License`),
							//Comment:   text.Chars(``),
							Type:      element.TypeText,
							SortOrder: 40,
							Visible:   element.VisibleYes,
							Scopes:    scope.PermDefault,
						},
						element.Field{
							// Path: `net/geoip_maxmind/webservice_timeout`,
							ID:    cfgpath.MakeRoute(`webservice_timeout`),
							Label: text.Chars(`Webservice HTTP request timeout`),
							Comment: text.Chars(`A duration string is a possibly signed sequence of decimal numbers, each with
optional fraction and a unit suffix, such as "300s", "-1.5h" or "2h45m". Valid
time units are "s", "m", "h".`),
							Type:      element.TypeText,
							SortOrder: 50,
							Visible:   element.VisibleYes,
							Scopes:    scope.PermDefault,
							Default:   time.Second * 15,
						},
						element.Field{
							// Path: `net/geoip_maxmind/webservice_redisurl`,
							ID:    cfgpath.MakeRoute(`webservice_redisurl`),
							Label: text.Chars(`Webservice Redis URL`),
							Comment: text.Chars(`An URL to the Redis instance to be used as a cache. If empty the default cache
will be in-memory and limited to XX MB.

URL has not match the scheme: redis://localhost:6379/X. Where X is the database
number. Or redis://ignored:passw0rd@localhost:6379/3 to use the password
passw0rd with database 3.`),
							Type:      element.TypeText,
							SortOrder: 50,
							Visible:   element.VisibleYes,
							Scopes:    scope.PermDefault,
						},
					),
				},
			),
		},
	)
}
