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
	"github.com/corestoreio/csfw/config"
	"github.com/corestoreio/csfw/config/cfgmodel"
	"github.com/corestoreio/csfw/net/geoip"
)

// Default creates new geoip.Option slice with the default configuration
// structure. It panics on error, so us it only during the app init phase.
func Default(opts ...cfgmodel.Option) geoip.OptionFactoryFunc {
	cfgStruct, err := NewConfigStructure()
	if err != nil {
		panic(err)
	}
	return PrepareOptions(New(cfgStruct, opts...))
}

// PrepareOptions creates a closure around the type Backend. The closure will be
// used during a scoped request to figure out the configuration depending on the
// incoming scope. An option array will be returned by the closure.
func PrepareOptions(be *Backend) geoip.OptionFactoryFunc {

	return func(sg config.ScopedGetter) []geoip.Option {
		var opts [6]geoip.Option
		//var i int
		//scp, id := sg.Scope()
		//
		//acc, err := be.NetGeoipAllowedCountries.Get(sg)
		//if err != nil {
		//	return optError(errors.Wrap(err, "[backendgeoip] NetGeoipAllowedCountries.Get"))
		//}
		//opts[i] = geoip.WithAllowedCountryCodes(scp, id, acc...)
		//i++

		return opts[:]
	}
}

func optError(err error) []geoip.Option {
	return []geoip.Option{func(s *geoip.Service) error {
		return err
	}}
}
