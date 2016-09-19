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
			opts      [1]cors.Option
			settings  cors.Settings
			scpHashes = make(scope.Hashes, 0, 8)
		)

		// For now the scope for all options depends on the scope of the
		// setting: NetCorsExposedHeaders

		// EXPOSED HEADERS
		eh, scp, err := be.ExposedHeaders.Get(sg)
		if err != nil {
			return cors.OptionsError(errors.Wrap(err, "[backendcors] NetCorsExposedHeaders.Get"))
		}
		settings.ExposedHeaders = eh
		scpHashes = append(scpHashes, scp)

		// ALLOWED ORIGINS
		ao, scp, err := be.AllowedOrigins.Get(sg)
		if err != nil {
			return cors.OptionsError(errors.Wrap(err, "[backendcors] NetCorsAllowedOrigins.Get"))
		}
		settings.AllowedOrigins = ao
		scpHashes = append(scpHashes, scp)

		// ALLOW ORIGIN REGEX
		aor, scp, err := be.AllowOriginRegex.Get(sg)
		if err != nil {
			return cors.OptionsError(errors.Wrap(err, "[backendcors] NetCorsAllowedOriginRegex.Get"))
		}
		if len(aor) > 1 {
			r, err := regexp.Compile(aor)
			if err != nil {
				return cors.OptionsError(errors.NewFatalf("[backendcors] NetCorsAllowedOriginRegex.regexp.Compile: %s", err))
			}
			settings.AllowOriginFunc = func(o string) bool {
				return r.MatchString(o)
			}
			scpHashes = append(scpHashes, scp)
		}

		// ALLOWED METHODS
		am, scp, err := be.AllowedMethods.Get(sg)
		if err != nil {
			return cors.OptionsError(errors.Wrap(err, "[backendcors] NetCorsAllowedMethods.Get"))
		}
		settings.AllowedMethods = am
		scpHashes = append(scpHashes, scp)

		// ALLOWED HEADERS
		ah, scp, err := be.AllowedHeaders.Get(sg)
		if err != nil {
			return cors.OptionsError(errors.Wrap(err, "[backendcors] NetCorsAllowedHeaders.Get"))
		}
		settings.AllowedHeaders = ah
		scpHashes = append(scpHashes, scp)

		// ALLOW CREDENTIALS
		ac, scp, err := be.AllowCredentials.Get(sg)
		if err != nil {
			return cors.OptionsError(errors.Wrap(err, "[backendcors] NetCorsAllowCredentials.Get"))
		}
		settings.AllowCredentials = ac
		scpHashes = append(scpHashes, scp)

		// OPTIONS PASSTHROUGH
		op, scp, err := be.OptionsPassthrough.Get(sg)
		if err != nil {
			return cors.OptionsError(errors.Wrap(err, "[backendcors] NetCorsOptionsPassthrough.Get"))
		}
		settings.OptionsPassthrough = op
		scpHashes = append(scpHashes, scp)

		// MAX AGE
		ma, scp, err := be.MaxAge.Get(sg)
		if err != nil {
			return cors.OptionsError(errors.Wrap(err, "[backendcors] NetCorsMaxAge.Get"))
		}
		settings.MaxAge = ma
		scpHashes = append(scpHashes, scp)

		// check the scope
		appliedScp, err := scpHashes.Lowest()
		if err != nil {
			return cors.OptionsError(errors.Wrap(err, "[backendcors] ScopeHashes.Lowest"))
		}

		opts[0] = cors.WithSettings(appliedScp, settings)
		return opts[:]
	}
}
