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
	"github.com/corestoreio/errors"
	"github.com/corestoreio/pkg/config"
	"github.com/corestoreio/pkg/net/auth"
)

// PrepareOptionFactory creates a closure around the type Backend. The closure
// will be used during a scoped request to figure out the configuration
// depending on the incoming scope. An option array will be returned by the
// closure.
func (be *Configuration) PrepareOptionFactory() auth.OptionFactoryFunc {
	return func(sg config.Scoped) []auth.Option {
		var opts [8]auth.Option
		var i int

		off, err := be.Disabled.Get(sg)
		if err != nil {
			return auth.OptionsError(errors.Wrap(err, "[backendauth] Disabled.Get"))
		}
		opts[i] = auth.WithDisable(off, sg.ScopeIDs()...)
		i++
		if off {
			return opts[:i]
		}

		return opts[:]
	}
}
