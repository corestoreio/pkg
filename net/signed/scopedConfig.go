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

package signed

import (
	"github.com/corestoreio/csfw/util/errors"
	"github.com/corestoreio/csfw/util/hashpool"
)

// ScopedConfig scoped based configuration and should not be embedded into your
// own types. Call ScopedConfig.ScopeHash to know to which scope this
// configuration has been bound to.
type ScopedConfig struct {
	scopedConfigGeneric

	// start of package specific config values

	// Disabled set to true to disable rate limiting
	Disabled bool

	hashPool hashpool.Tank
}

// newScopedConfig creates a new object with the minimum needed configuration.
func newScopedConfig() *ScopedConfig {
	return &ScopedConfig{
		scopedConfigGeneric: newScopedConfigGeneric(),
	}
}

// IsValid a configuration for a scope is only then valid when several fields
// are not empty: RateLimiter, DeniedHandler and VaryByer.
func (sc ScopedConfig) IsValid() error {
	if sc.lastErr != nil {
		return errors.Wrap(sc.lastErr, "[signed] scopedConfig.isValid has an lastErr")
	}
	if sc.ScopeHash == 0 {
		return errors.NewNotValidf(errScopedConfigNotValid, sc.ScopeHash)
	}
	return nil
}
