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
	"github.com/corestoreio/csfw/util/errors"
)

// PrepareOptionFactory creates a closure around the type Backend. The closure
// will be used during a scoped request to figure out the configuration
// depending on the incoming scope. An option array will be returned by the
// closure.
func (be *Configuration) PrepareOptionFactory() cors.OptionFactoryFunc {
	return func(sg config.Scoped) []cors.Option {
		var (
			opts     [2]cors.Option
			settings cors.Settings
			err      error
		)

		// EXPOSED HEADERS
		settings.ExposedHeaders, err = be.ExposedHeaders.Get(sg)
		if err != nil {
			return cors.OptionsError(errors.Wrap(err, "[backendcors] NetCorsExposedHeaders.Get"))
		}

		// ALLOWED ORIGINS
		settings.AllowedOrigins, err = be.AllowedOrigins.Get(sg)
		if err != nil {
			return cors.OptionsError(errors.Wrap(err, "[backendcors] NetCorsAllowedOrigins.Get"))
		}

		// ALLOW ORIGIN REGEX
		aor, err := be.AllowOriginRegex.Get(sg)
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
		}

		// ALLOWED METHODS
		settings.AllowedMethods, err = be.AllowedMethods.Get(sg)
		if err != nil {
			return cors.OptionsError(errors.Wrap(err, "[backendcors] NetCorsAllowedMethods.Get"))
		}

		// ALLOWED HEADERS
		settings.AllowedHeaders, err = be.AllowedHeaders.Get(sg)
		if err != nil {
			return cors.OptionsError(errors.Wrap(err, "[backendcors] NetCorsAllowedHeaders.Get"))
		}

		// ALLOW CREDENTIALS
		settings.AllowCredentials, err = be.AllowCredentials.Get(sg)
		if err != nil {
			return cors.OptionsError(errors.Wrap(err, "[backendcors] NetCorsAllowCredentials.Get"))
		}

		// OPTIONS PASSTHROUGH
		settings.OptionsPassthrough, err = be.OptionsPassthrough.Get(sg)
		if err != nil {
			return cors.OptionsError(errors.Wrap(err, "[backendcors] NetCorsOptionsPassthrough.Get"))
		}

		// MAX AGE
		settings.MaxAge, err = be.MaxAge.Get(sg)
		if err != nil {
			return cors.OptionsError(errors.Wrap(err, "[backendcors] NetCorsMaxAge.Get"))
		}

		// in case someone marks the config as partially applied now it's time to revert
		// it.
		opts[0] = cors.WithMarkPartiallyApplied(false, sg.ScopeIDs()...)
		opts[1] = cors.WithSettings(settings, sg.ScopeIDs()...)
		return opts[:]
	}
}
