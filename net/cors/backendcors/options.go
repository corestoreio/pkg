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
	"github.com/corestoreio/csfw/net/cors"
	"github.com/corestoreio/csfw/util/errors"
)

// Default creates new cors.Option slice with the default configuration
// structure. It panics on error, so us it only during the app init phase.
func Default(opts ...cfgmodel.Option) cors.OptionFactoryFunc {
	cfgStruct, err := NewConfigStructure()
	if err != nil {
		panic(err)
	}
	return PrepareOptions(New(cfgStruct, opts...))
}

// PrepareOptions creates a closure around the type Backend. The closure will
// be used during a scoped request to figure out the configuration depending on
// the incoming scope. An option array will be returned by the closure.
func PrepareOptions(be *Backend) cors.OptionFactoryFunc {

	return func(sg config.ScopedGetter) []cors.Option {
		var opts [8]cors.Option
		var i int
		scp, id := sg.Scope()

		// EXPOSED HEADERS
		eh, err := be.NetCorsExposedHeaders.Get(sg)
		if err != nil {
			return optError(errors.Wrap(err, "[backendcors] NetCorsExposedHeaders.Get"))
		}
		opts[i] = cors.WithExposedHeaders(scp, id, eh...)
		i++

		// ALLOWED ORIGINS
		ao, err := be.NetCorsAllowedOrigins.Get(sg)
		if err != nil {
			return optError(errors.Wrap(err, "[backendcors] NetCorsAllowedOrigins.Get"))
		}
		opts[i] = cors.WithAllowedOrigins(scp, id, ao...)
		i++

		// ALLOW ORIGIN REGEX
		aor, err := be.NetCorsAllowOriginRegex.Get(sg)
		if err != nil {
			return optError(errors.Wrap(err, "[backendcors] NetCorsAllowedOriginRegex.Get"))
		}
		if len(aor) > 1 {
			r, err := regexp.Compile(aor)
			if err != nil {
				return optError(errors.NewFatalf("[backendcors] NetCorsAllowedOriginRegex.regexp.Compile: %s", err))
			}
			opts[i] = cors.WithAllowOriginFunc(scp, id, func(o string) bool {
				return r.MatchString(o)
			})
		}
		i++

		// ALLOWED METHODS
		am, err := be.NetCorsAllowedMethods.Get(sg)
		if err != nil {
			return optError(errors.Wrap(err, "[backendcors] NetCorsAllowedMethods.Get"))
		}
		opts[i] = cors.WithAllowedMethods(scp, id, am...)
		i++

		// ALLOWED HEADERS
		ah, err := be.NetCorsAllowedHeaders.Get(sg)
		if err != nil {
			return optError(errors.Wrap(err, "[backendcors] NetCorsAllowedHeaders.Get"))
		}
		opts[i] = cors.WithAllowedHeaders(scp, id, ah...)
		i++

		// ALLOW CREDENTIALS
		ac, err := be.NetCorsAllowCredentials.Get(sg)
		if err != nil {
			return optError(errors.Wrap(err, "[backendcors] NetCorsAllowCredentials.Get"))
		}
		opts[i] = cors.WithAllowCredentials(scp, id, ac)
		i++

		// OPTIONS PASSTHROUGH
		op, err := be.NetCorsOptionsPassthrough.Get(sg)
		if err != nil {
			return optError(errors.Wrap(err, "[backendcors] NetCorsOptionsPassthrough.Get"))
		}
		opts[i] = cors.WithOptionsPassthrough(scp, id, op)
		i++

		// MAX AGE
		ma, err := be.NetCorsMaxAge.Get(sg)
		if err != nil {
			return optError(errors.Wrap(err, "[backendcors] NetCorsMaxAge.Get"))
		}
		opts[i] = cors.WithMaxAge(scp, id, ma)
		i++

		return opts[:]
	}
}

func optError(err error) []cors.Option {
	return []cors.Option{func(s *cors.Service) error {
		return err
	}}
}
