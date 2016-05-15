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
	"regexp"

	"github.com/corestoreio/csfw/config"
	"github.com/corestoreio/csfw/config/cfgmodel"
	"github.com/corestoreio/csfw/net/mwcors"
	"github.com/corestoreio/csfw/util/errors"
)

// Default creates new mwcors.Option slice with the default configuration
// structure. It panics on error, so us it only during the app init phase.
func Default(opts ...cfgmodel.Option) mwcors.ScopedOptionFunc {
	cfgStruct, err := NewConfigStructure()
	if err != nil {
		panic(err)
	}
	return PrepareOptions(New(cfgStruct, opts...))
}

// PrepareOptions creates a closure around the type Backend. The closure will
// be used during a scoped request to figure out the configuration depending on
// the incoming scope. An option array will be returned by the closure.
func PrepareOptions(be *Backend) mwcors.ScopedOptionFunc {

	return func(sg config.ScopedGetter) []mwcors.Option {
		var opts [8]mwcors.Option
		var i int
		scp, id := sg.Scope()

		// EXPOSED HEADERS
		eh, err := be.NetCorsExposedHeaders.Get(sg)
		if err != nil {
			opts[i] = func(s *mwcors.Service) {
				s.AddError(errors.Wrap(err, "[backendcors] NetCorsExposedHeaders.Get"))
			}
			return opts[:]
		}
		opts[i] = mwcors.WithExposedHeaders(scp, id, eh...)
		i++

		// ALLOWED ORIGINS
		ao, err := be.NetCorsAllowedOrigins.Get(sg)
		if err != nil {
			opts[i] = func(s *mwcors.Service) {
				s.AddError(errors.Wrap(err, "[backendcors] NetCorsAllowedOrigins.Get"))
			}
			return opts[:]
		}
		opts[i] = mwcors.WithAllowedOrigins(scp, id, ao...)
		i++

		// ALLOW ORIGIN REGEX
		aor, err := be.NetCorsAllowOriginRegex.Get(sg)
		if err != nil {
			opts[i] = func(s *mwcors.Service) {
				s.AddError(errors.Wrap(err, "[backendcors] NetCorsAllowedOriginRegex.Get"))
			}
			return opts[:]
		}
		if len(aor) > 1 {
			r, err := regexp.Compile(aor)
			if err != nil {
				opts[i] = func(s *mwcors.Service) {
					s.AddError(errors.NewFatal(err, "[backendcors] NetCorsAllowedOriginRegex.regexp.Compile"))
				}
				return opts[:]
			}
			opts[i] = mwcors.WithAllowOriginFunc(scp, id, func(o string) bool {
				return r.MatchString(o)
			})
		}
		i++

		// ALLOWED METHODS
		am, err := be.NetCorsAllowedMethods.Get(sg)
		if err != nil {
			opts[i] = func(s *mwcors.Service) {
				s.AddError(errors.Wrap(err, "[backendcors] NetCorsAllowedMethods.Get"))
			}
			return opts[:]
		}
		opts[i] = mwcors.WithAllowedMethods(scp, id, am...)
		i++

		// ALLOWED HEADERS
		ah, err := be.NetCorsAllowedHeaders.Get(sg)
		if err != nil {
			opts[i] = func(s *mwcors.Service) {
				s.AddError(errors.Wrap(err, "[backendcors] NetCorsAllowedHeaders.Get"))
			}
			return opts[:]
		}
		opts[i] = mwcors.WithAllowedHeaders(scp, id, ah...)
		i++

		// ALLOW CREDENTIALS
		ac, err := be.NetCorsAllowCredentials.Get(sg)
		if err != nil {
			opts[i] = func(s *mwcors.Service) {
				s.AddError(errors.Wrap(err, "[backendcors] NetCorsAllowCredentials.Get"))
			}
			return opts[:]
		}
		opts[i] = mwcors.WithAllowCredentials(scp, id, ac)
		i++

		// OPTIONS PASSTHROUGH
		op, err := be.NetCorsOptionsPassthrough.Get(sg)
		if err != nil {
			opts[i] = func(s *mwcors.Service) {
				s.AddError(errors.Wrap(err, "[backendcors] NetCorsOptionsPassthrough.Get"))
			}
			return opts[:]
		}
		opts[i] = mwcors.WithOptionsPassthrough(scp, id, op)
		i++

		// MAX AGE
		ma, err := be.NetCorsMaxAge.Get(sg)
		if err != nil {
			opts[i] = func(s *mwcors.Service) {
				s.AddError(errors.Wrap(err, "[backendcors] NetCorsMaxAge.Get"))
			}
			return opts[:]
		}
		opts[i] = mwcors.WithMaxAge(scp, id, ma)
		i++

		return opts[:]
	}
}
