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

package ratelimit

import "github.com/corestoreio/csfw/store/scope"

// auto generated: do not edit. See net/gen eric package

// scopedConfigGeneric private internal scoped based configuration used for
// embedding into scopedConfig type.
type scopedConfigGeneric struct {
	// lastErr used during selecting the config from the scopeCache map and gets
	// filled if an entry cannot be found.
	lastErr error
	// scopeHash defines the scope to which this configuration is bound to.
	scopeHash scope.Hash
}

func newScopedConfigError(err error) *scopedConfig {
	return &scopedConfig{
		scopedConfigGeneric: scopedConfigGeneric{
			lastErr: err,
		},
	}
}

func (sc *scopedConfig) printScope() string {
	if sc == nil {
		return "<nil>"
	}
	return sc.scopeHash.String()
}
