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
	"github.com/corestoreio/csfw/store/scope"
	"github.com/corestoreio/csfw/util/errors"
)

// PrepareOptions creates a closure around the type Backend. The closure will
// be used during a scoped request to figure out the configuration depending on
// the incoming scope. An option array will be returned by the closure.
func PrepareOptions(be *Backend) cors.OptionFactoryFunc {
	return func(sg config.Scoped) []cors.Option {
		var (
			opts  [8]cors.Option
			i     int // used as index in opts
			scp   scope.Scope
			scpID int64
		)

		// EXPOSED HEADERS
		eh, h, err := be.NetCorsExposedHeaders.Get(sg)
		if err != nil {
			return optError(errors.Wrap(err, "[backendcors] NetCorsExposedHeaders.Get"))
		}
		scp, scpID = h.Unpack()
		opts[i] = cors.WithExposedHeaders(scp, scpID, eh...)
		i++

		// ALLOWED ORIGINS
		ao, h, err := be.NetCorsAllowedOrigins.Get(sg)
		if err != nil {
			return optError(errors.Wrap(err, "[backendcors] NetCorsAllowedOrigins.Get"))
		}

		scp, scpID = h.Unpack()
		opts[i] = cors.WithAllowedOrigins(scp, scpID, ao...)
		i++

		// ALLOW ORIGIN REGEX
		aor, h, err := be.NetCorsAllowOriginRegex.Get(sg)
		if err != nil {
			return optError(errors.Wrap(err, "[backendcors] NetCorsAllowedOriginRegex.Get"))
		}
		if len(aor) > 1 {
			r, err := regexp.Compile(aor)
			if err != nil {
				return optError(errors.NewFatalf("[backendcors] NetCorsAllowedOriginRegex.regexp.Compile: %s", err))
			}
			scp, scpID = h.Unpack()
			opts[i] = cors.WithAllowOriginFunc(scp, scpID, func(o string) bool {
				return r.MatchString(o)
			})
		}
		i++

		// ALLOWED METHODS
		am, h, err := be.NetCorsAllowedMethods.Get(sg)
		if err != nil {
			return optError(errors.Wrap(err, "[backendcors] NetCorsAllowedMethods.Get"))
		}
		scp, scpID = h.Unpack()
		opts[i] = cors.WithAllowedMethods(scp, scpID, am...)
		i++

		// ALLOWED HEADERS
		ah, h, err := be.NetCorsAllowedHeaders.Get(sg)
		if err != nil {
			return optError(errors.Wrap(err, "[backendcors] NetCorsAllowedHeaders.Get"))
		}
		scp, scpID = h.Unpack()
		opts[i] = cors.WithAllowedHeaders(scp, scpID, ah...)
		i++

		// ALLOW CREDENTIALS
		ac, h, err := be.NetCorsAllowCredentials.Get(sg)
		if err != nil {
			return optError(errors.Wrap(err, "[backendcors] NetCorsAllowCredentials.Get"))
		}
		scp, scpID = h.Unpack()
		opts[i] = cors.WithAllowCredentials(scp, scpID, ac)
		i++

		// OPTIONS PASSTHROUGH
		op, h, err := be.NetCorsOptionsPassthrough.Get(sg)
		if err != nil {
			return optError(errors.Wrap(err, "[backendcors] NetCorsOptionsPassthrough.Get"))
		}
		scp, scpID = h.Unpack()
		opts[i] = cors.WithOptionsPassthrough(scp, scpID, op)
		i++

		// MAX AGE
		ma, h, err := be.NetCorsMaxAge.Get(sg)
		if err != nil {
			return optError(errors.Wrap(err, "[backendcors] NetCorsMaxAge.Get"))
		}
		scp, scpID = h.Unpack()
		opts[i] = cors.WithMaxAge(scp, scpID, ma)
		i++

		return opts[:]
	}
}
