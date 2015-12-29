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

package ctxcors

import (
	"github.com/corestoreio/csfw/config"
	"github.com/corestoreio/csfw/store/scope"
)

// Path constants defines the configuration paths in core_config_data
const (
	PathCorsExposedHeaders   = "web/cors/exposed_headers"    // csv
	PathCorsAllowedOrigins   = "web/cors/allowed_origins"    // csv
	PathCorsAllowedMethods   = "web/cors/allowed_methods"    // csv
	PathCorsAllowedHeaders   = "web/cors/allowed_headers"    // csv
	PathCorsAllowCredentials = "web/cors/allowe_credentials" // bool
)

func ConfigExposedHeaders(cg config.Getter, s scope.Scope, id int64) []string {
	//	fields, err := PackageConfiguration.FindFieldByPath(PathCorsExposedHeaders)
	//	fields.Default.(string)
	//	fields.BackendModel.AddData()
	//	headers, _ := cg.String(config.Path(PathCorsExposedHeaders), config.Scope(s, id))
	return nil
}

// PackageConfiguration contains the main configuration
var PackageConfiguration config.SectionSlice

func init() {
	PackageConfiguration = config.MustNewConfiguration(
		&config.Section{
			ID: "web", // defined in ?
			Groups: config.GroupSlice{
				&config.Group{
					ID:        "cors",
					Label:     `CORS Cross Origin Resource Sharing`,
					Comment:   ``,
					SortOrder: 150,
					Scope:     scope.NewPerm(scope.DefaultID),
					Fields: config.FieldSlice{
						&config.Field{
							// Path: `web/cors/exposed_headers`,
							ID:           "exposed_headers",
							Label:        `Exposed Headers`,
							Comment:      `Indicates which headers are safe to expose to the API of a CORS API specification. Separate via line break`,
							Type:         config.TypeTextarea,
							SortOrder:    10,
							Visible:      config.VisibleYes,
							Scope:        scope.NewPerm(scope.WebsiteID),
							Default:      nil,
							BackendModel: nil, // CSV
							SourceModel:  nil,
						},
						&config.Field{
							// Path: `web/cors/allowed_origins`,
							ID:    "allowed_origins",
							Label: `Allowed Origins`,
							Comment: `Is a list of origins a cross-domain request can be executed from.
If the special "*" value is present in the list, all origins will be allowed.
An origin may contain a wildcard (*) to replace 0 or more characters
(i.e.: http://*.domain.com). Usage of wildcards implies a small performance penality.
Only one wildcard can be used per origin.
Default value is ["*"]`,
							Type:         config.TypeTextarea,
							SortOrder:    20,
							Visible:      config.VisibleYes,
							Scope:        scope.NewPerm(scope.WebsiteID),
							Default:      nil,
							BackendModel: nil, // CSV
							SourceModel:  nil,
						},
						// TODO add other fields
					},
				},
			},
		},
	)
}
