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
	"github.com/corestoreio/csfw/config/cfgmodel"
	"github.com/corestoreio/csfw/config/cfgsource"
	"github.com/corestoreio/csfw/config/element"
	"github.com/corestoreio/csfw/net/auth"
)

// Configuration just exported for the sake of documentation. See fields
// for more information. The PkgBackend handles the reading and writing
// of configuration values within this package.
type Configuration struct {
	*auth.OptionFactories

	// Disabled indicates whether authentication has been enabled or not.
	// Path: net/auth/disabled
	Disabled cfgmodel.Bool

	// AllowedIP indicates which IPs are allowed.
	// Separate via line break (\n).
	//
	// Path: net/auth/allowed_ips
	AllowedIPs cfgmodel.StringCSV

	// DeniedIPs indicates which IPs are denied.
	// Separate via line break (\n).
	//
	// Path: net/auth/denied_ips
	DeniedIPs cfgmodel.StringCSV

	// AllowedIPRange indicates which IP ranges are denied.
	// Separate via line break (\n).
	//
	// Path: net/auth/denied_ips
	AllowedIPRange ConfigIPRange

	// DeniedIPRange indicates which IP ranges are denied.
	// Separate via line break (\n).
	//
	// Path: net/auth/denied_ips
	DeniedIPRange ConfigIPRange

	// and so on
	// range based allowances and denies
}

// New initializes the backend configuration models containing the cfgpath.Route
// variable to the appropriate entries in the storage. The argument SectionSlice
// and opts will be applied to all models.
func New(cfgStruct element.SectionSlice, opts ...cfgmodel.Option) *Configuration {
	be := &Configuration{
		OptionFactories: auth.NewOptionFactories(),
	}

	opts = append(opts, cfgmodel.WithFieldFromSectionSlice(cfgStruct))

	be.Disabled = cfgmodel.NewBool(`net/auth/disabled`, append(opts, cfgmodel.WithSource(cfgsource.EnableDisable))...)

	//opts = append(opts, cfgmodel.WithFieldFromSectionSlice(cfgStruct))
	//optsCSV := append([]cfgmodel.Option{}, opts...)
	//optsCSV = append(optsCSV, cfgmodel.WithFieldFromSectionSlice(cfgStruct), cfgmodel.WithCSVComma('\n'))
	//optsYN := append([]cfgmodel.Option{}, opts...)
	//optsYN = append(optsYN, cfgmodel.WithFieldFromSectionSlice(cfgStruct), cfgmodel.WithSource(cfgsource.YesNo))
	//
	//pp.AllowedIPs = cfgmodel.NewStringCSV(`net/auth/exposed_headers`, optsCSV...)

	return be
}
