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

package backendcors

import (
	"github.com/corestoreio/csfw/config/cfgmodel"
	"github.com/corestoreio/csfw/config/element"
	"github.com/corestoreio/csfw/config/source"
)

// Backend just exported for the sake of documentation. See fields
// for more information. The PkgBackend handles the reading and writing
// of configuration values within this package.
type Backend struct {
	cfgmodel.PkgBackend

	// NetCorsExposedHeaders indicates which headers are safe to
	// expose to the API of a CORS API specification.
	// Separate via line break (\n).
	//
	// Path: net/cors/exposed_headers
	NetCorsExposedHeaders cfgmodel.StringCSV

	// NetCorsAllowedOrigins is a list of origins a cross-domain request
	// can be executed from. If the special "*" value is present in the list, all origins
	// will be allowed. An origin may contain a wildcard (*) to replace 0 or more characters
	// (i.e.: http://*.domain.com). Usage of wildcards implies a small performance penality.
	// Only one wildcard can be used per origin. Default value is ["*"]
	// Separate via line break (\n).
	//
	// Path: net/cors/allowed_origins
	NetCorsAllowedOrigins cfgmodel.StringCSV

	// NetCorsAllowOriginRegex same as NetCorsAllowedOrigins but uses
	// a regex to check for the domains.
	//
	// Path: net/cors/allow_origin_regex
	NetCorsAllowOriginRegex cfgmodel.Str

	// NetCorsAllowedMethods a list of methods the client is allowed to
	// use with cross-domain requests. Default value is simple methods (GET and POST)
	// Separate via line break (\n).
	//
	// Path: net/cors/allowed_methods
	NetCorsAllowedMethods cfgmodel.StringCSV

	// NetCorsAllowedHeaders A list of non simple headers the client is
	// allowed to use with cross-domain requests. If the special "*" value is present
	// in the list, all headers will be allowed. Default value is [] but "Origin" is
	// always appended to the list.
	// Separate via line break (\n).
	//
	// Path: net/cors/allowed_headers
	NetCorsAllowedHeaders cfgmodel.StringCSV

	// NetCorsAllowCredentials Indicates whether the request can include
	// user credentials like cookies, HTTP authentication or client side SSL certificates.
	//
	// Path: net/cors/allow_credentials
	NetCorsAllowCredentials cfgmodel.Bool

	// NetCorsOptionsPassthrough instructs preflight to let other potential next handlers to
	// process the OPTIONS method. Turn this on if your application handles OPTIONS.
	//
	// Path: net/cors/options_passthrough
	NetCorsOptionsPassthrough cfgmodel.Bool

	// NetCorsMaxAge Indicates how long (in seconds) the results
	// of a preflight request can be cached.
	//
	// Path: net/cors/max_age
	NetCorsMaxAge cfgmodel.Duration
}

// New initializes the backend configuration models containing the
// cfgpath.Route variable to the appropriate entries.
// The function Load() will be executed to apply the SectionSlice
// to all models. See Load() for more details.
func New(cfgStruct element.SectionSlice, opts ...cfgmodel.Option) *Backend {
	return (&Backend{}).Load(cfgStruct, opts...)
}

// Load creates the configuration models for each PkgBackend field.
// Internal mutex will protect the fields during loading.
// The argument SectionSlice will be applied to all models.
func (pp *Backend) Load(cfgStruct element.SectionSlice, opts ...cfgmodel.Option) *Backend {
	pp.Lock()
	defer pp.Unlock()

	opts = append(opts, cfgmodel.WithFieldFromSectionSlice(cfgStruct))
	optsCSV := append([]cfgmodel.Option{}, opts...)
	optsCSV = append(optsCSV, cfgmodel.WithFieldFromSectionSlice(cfgStruct), cfgmodel.WithCSVSeparator('\n'))
	optsYN := append([]cfgmodel.Option{}, opts...)
	optsYN = append(optsYN, cfgmodel.WithFieldFromSectionSlice(cfgStruct), cfgmodel.WithSource(source.YesNo))

	pp.NetCorsExposedHeaders = cfgmodel.NewStringCSV(`net/cors/exposed_headers`, optsCSV...)
	pp.NetCorsAllowedOrigins = cfgmodel.NewStringCSV(`net/cors/allowed_origins`, optsCSV...)
	pp.NetCorsAllowOriginRegex = cfgmodel.NewStr(`net/cors/allow_origin_regex`, opts...)
	pp.NetCorsAllowedMethods = cfgmodel.NewStringCSV(`net/cors/allowed_methods`, optsCSV...)
	pp.NetCorsAllowedHeaders = cfgmodel.NewStringCSV(`net/cors/allowed_headers`, optsCSV...)
	pp.NetCorsAllowCredentials = cfgmodel.NewBool(`net/cors/allow_credentials`, optsYN...)
	pp.NetCorsOptionsPassthrough = cfgmodel.NewBool(`net/cors/options_passthrough`, optsYN...)
	pp.NetCorsMaxAge = cfgmodel.NewDuration(`net/cors/max_age`, opts...)

	return pp
}
