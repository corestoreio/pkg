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
	"github.com/corestoreio/csfw/config"
	"github.com/corestoreio/csfw/net/geoip"
	"github.com/corestoreio/csfw/util/errors"
)

// PrepareOptionFactory creates a closure around the type Backend. The closure
// will be used during a scoped request to figure out the configuration
// depending on the incoming scope. An option array will be returned by the
// closure.
func (be *Configuration) PrepareOptionFactory() geoip.OptionFactoryFunc {
	return func(sg config.Scoped) []geoip.Option {
		var (
			opts [6]geoip.Option
			i    int // used as index in opts
		)

		acc, err := be.AllowedCountries.Get(sg)
		if err != nil {
			return geoip.OptionsError(errors.Wrap(err, "[backendgeoip] NetGeoipAllowedCountries.Get"))
		}
		opts[i] = geoip.WithAllowedCountryCodes(acc, sg.ScopeIDs()...)
		i++

		// REDIRECT TO ALTERNATIVE URL
		arURL, err := be.AlternativeRedirect.Get(sg)
		if err != nil {
			return geoip.OptionsError(errors.Wrap(err, "[backendgeoip] NetGeoipAlternativeRedirect.Get"))
		}
		arCode, err := be.AlternativeRedirectCode.Get(sg)
		if err != nil {
			return geoip.OptionsError(errors.Wrap(err, "[backendgeoip] NetGeoipAlternativeRedirectCode.Get"))
		}
		if arCode > 0 && arURL != nil {
			opts[i] = geoip.WithAlternativeRedirect(arURL.String(), arCode, sg.ScopeIDs()...)
		}
		i++

		source, err := be.DataSource.Get(sg)
		if err != nil {
			return geoip.OptionsError(errors.Wrap(err, "[backendgeoip] DataSource.Get"))
		}

		// source contains the configured geo location data source either file
		// or webservice.
		ofFnc, err := be.Lookup(source) // off = OptionFactoryFunc
		if err != nil {
			return geoip.OptionsError(errors.Wrap(err, "[backendgeoip] Backend.Lookup"))
		}
		return append(opts[:], ofFnc(sg)...)
	}
}
