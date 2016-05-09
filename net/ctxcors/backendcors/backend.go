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

	// NetCtxcorsExposedHeaders indicates which headers are safe to
	// expose to the API of a CORS API specification.
	// Separate via line break (\n).
	//
	// Path: net/ctxcors/exposed_headers
	NetCtxcorsExposedHeaders cfgmodel.StringCSV

	// NetCtxcorsAllowedOrigins is a list of origins a cross-domain request
	// can be executed from. If the special "*" value is present in the list, all origins
	// will be allowed. An origin may contain a wildcard (*) to replace 0 or more characters
	// (i.e.: http://*.domain.com). Usage of wildcards implies a small performance penality.
	// Only one wildcard can be used per origin. Default value is ["*"]
	// Separate via line break (\n).
	//
	// Path: net/ctxcors/allowed_origins
	NetCtxcorsAllowedOrigins cfgmodel.StringCSV

	// NetCtxcorsAllowedMethods a list of methods the client is allowed to
	// use with cross-domain requests. Default value is simple methods (GET and POST)
	// Separate via line break (\n).
	//
	// Path: net/ctxcors/allowed_methods
	NetCtxcorsAllowedMethods cfgmodel.StringCSV

	// NetCtxcorsAllowedHeaders A list of non simple headers the client is
	// allowed to use with cross-domain requests. If the special "*" value is present
	// in the list, all headers will be allowed. Default value is [] but "Origin" is
	// always appended to the list.
	// Separate via line break (\n).
	//
	// Path: net/ctxcors/allowed_headers
	NetCtxcorsAllowedHeaders cfgmodel.StringCSV

	// NetCtxcorsAllowCredentials Indicates whether the request can include
	// user credentials like cookies, HTTP authentication or client side SSL certificates.
	//
	// Path: net/ctxcors/allow_credentials
	NetCtxcorsAllowCredentials cfgmodel.Bool

	// NetCtxcorsOptionsPassthrough instructs preflight to let other potential next handlers to
	// process the OPTIONS method. Turn this on if your application handles OPTIONS.
	//
	// Path: net/ctxcors/options_passthrough
	NetCtxcorsOptionsPassthrough cfgmodel.Bool

	// NetCtxcorsMaxAge Indicates how long (in seconds) the results
	// of a preflight request can be cached.
	//
	// Path: net/ctxcors/max_age
	NetCtxcorsMaxAge cfgmodel.Duration
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
func (pp *Backend) Load(cfgStruct element.SectionSlice) *Backend {
	pp.Lock()
	defer pp.Unlock()

	opt := cfgmodel.WithFieldFromSectionSlice(cfgStruct)

	pp.NetCtxcorsExposedHeaders = cfgmodel.NewStringCSV(`net/ctxcors/exposed_headers`, opt, cfgmodel.WithCSVSeparator('\n'))
	pp.NetCtxcorsAllowedOrigins = cfgmodel.NewStringCSV(`net/ctxcors/allowed_origins`, opt, cfgmodel.WithCSVSeparator('\n'))
	pp.NetCtxcorsAllowedMethods = cfgmodel.NewStringCSV(`net/ctxcors/allowed_methods`, opt, cfgmodel.WithCSVSeparator('\n'))
	pp.NetCtxcorsAllowedHeaders = cfgmodel.NewStringCSV(`net/ctxcors/allowed_headers`, opt, cfgmodel.WithCSVSeparator('\n'))
	pp.NetCtxcorsAllowCredentials = cfgmodel.NewBool(`net/ctxcors/allow_credentials`, opt, cfgmodel.WithSource(source.YesNo))
	pp.NetCtxcorsOptionsPassthrough = cfgmodel.NewBool(`net/ctxcors/allow_credentials`, opt, cfgmodel.WithSource(source.YesNo))
	pp.NetCtxcorsMaxAge = cfgmodel.NewDuration(`net/ctxcors/max_age`, opt)

	return pp
}
