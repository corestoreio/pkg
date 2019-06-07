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

package backendcors

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
					ID:    cfgpath.MakeRoute(`cors`),
					Label: text.Chars(`CORS Cross Origin Resource Sharing`),
					Comment: text.Chars(`CORS describes the CrossOriginResourceSharing
which is used to create a Container Filter that implements CORS. Cross-origin
resource sharing (CORS) is a mechanism that allows JavaScript on a web page to
make XMLHttpRequests to another domain, not the domain the JavaScript originated
from.`),
					MoreURL:   text.Chars(`http://en.wikipedia.org/wiki/Cross-origin_resource_sharing|http://enable-cors.org/server.html|http://www.html5rocks.com/en/tutorials/cors/#toc-handling-a-not-so-simple-request`),
					SortOrder: 160,
					Scopes:    scope.PermWebsite,
					Fields: element.MakeFields(
						element.Field{
							// Path: `net/cors/exposed_headers`,
							ID:    cfgpath.MakeRoute(`exposed_headers`),
							Label: text.Chars(`Exposed Headers`),
							Comment: text.Chars(`Indicates which headers are safe to
expose to the API of a CORS API specification. Separate via line break (\n)`),
							Type:      element.TypeTextarea,
							SortOrder: 10,
							Visible:   element.VisibleYes,
							Scopes:    scope.PermWebsite,
						},
						element.Field{
							// Path: `net/cors/allowed_origins`,
							ID:    cfgpath.MakeRoute(`allowed_origins`),
							Label: text.Chars(`Allowed Origins`),
							Comment: text.Chars(`Is a list of origins a cross-domain request
can be executed from. If the special "*" value is present in the list, all origins
will be allowed. An origin may contain a wildcard (*) to replace 0 or more characters
(i.e.: http://*.domain.com). Usage of wildcards implies a small performance penality.
Only one wildcard can be used per origin. Default value is ["*"]. Separate via line break (\n)`),
							Type:      element.TypeTextarea,
							SortOrder: 20,
							Visible:   element.VisibleYes,
							Scopes:    scope.PermWebsite,
							Default:   `*`,
						},
						element.Field{
							// Path: `net/cors/allow_origin_regex`,
							ID:        cfgpath.MakeRoute(`allow_origin_regex`),
							Label:     text.Chars(`Allowed Origin Regex`),
							Comment:   text.Chars(`Same as Allowed Origins but uses a PERL regular expression.`),
							Type:      element.TypeText,
							SortOrder: 25,
							Visible:   element.VisibleYes,
							Scopes:    scope.PermWebsite,
						},
						element.Field{
							// Path: `net/cors/allowed_methods`,
							ID:    cfgpath.MakeRoute(`allowed_methods`),
							Label: text.Chars(`Allowed Methods`),
							Comment: text.Chars(`A list of methods the client is allowed to
use with cross-domain requests. Default value is simple methods (GET and POST). Separate via line break (\n)`),
							Type:      element.TypeText,
							SortOrder: 30,
							Visible:   element.VisibleYes,
							Scopes:    scope.PermWebsite,
							Default:   "GET\nPOST",
						},
						element.Field{
							// Path: `net/cors/allowed_headers`,
							ID:    cfgpath.MakeRoute(`allowed_headers`),
							Label: text.Chars(`Allowed Headers`),
							Comment: text.Chars(`A list of non simple headers the client is
allowed to use with cross-domain requests. If the special "*" value is present
in the list, all headers will be allowed. Default value is [] but "Origin" is
always appended to the list. Separate via line break (\n)`),
							Type:      element.TypeText,
							SortOrder: 40,
							Visible:   element.VisibleYes,
							Scopes:    scope.PermWebsite,
							Default:   "Origin\nAccept\nContent-Type",
						},
						element.Field{
							// Path: `net/cors/allow_credentials`,
							ID:    cfgpath.MakeRoute(`allow_credentials`),
							Label: text.Chars(`Allow Credentials`),
							Comment: text.Chars(`Indicates whether the request can include
user credentials like cookies, HTTP authentication or client side SSL certificates.`),
							Type:      element.TypeSelect,
							SortOrder: 50,
							Visible:   element.VisibleYes,
							Scopes:    scope.PermWebsite,
							Default:   `false`,
						},
						element.Field{
							// Path: `net/cors/options_passthrough`,
							ID:    cfgpath.MakeRoute(`options_passthrough`),
							Label: text.Chars(`Options Passthrough`),
							Comment: text.Chars(`OptionsPassthrough instructs preflight to let other potential next handlers to
process the OPTIONS method. Turn this on if your application handles OPTIONS.`),
							Type:      element.TypeSelect,
							SortOrder: 60,
							Visible:   element.VisibleYes,
							Scopes:    scope.PermWebsite,
							Default:   `false`,
						},
						element.Field{
							// Path: `net/cors/max_age`,
							ID:    cfgpath.MakeRoute(`max_age`),
							Label: text.Chars(`Max Age`),
							Comment: text.Chars(`Indicates how long (in seconds) the results
of a preflight request can be cached.`),
							Tooltip: text.Chars(`A duration string is a possibly signed sequence of
decimal numbers, each with optional fraction and a unit suffix,
such as "300ms", "-1.5h" or "2h45m".
Valid time units are "ns", "us" (or "µs"), "ms", "s", "m", "h".`),
							Type:      element.TypeText,
							SortOrder: 70,
							Visible:   element.VisibleYes,
							Scopes:    scope.PermWebsite,
						},
					),
				},
			),
		},
	)
}
