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

package signedheader

import (
	"github.com/corestoreio/csfw/config"
	"github.com/corestoreio/csfw/net/signed"
	"github.com/corestoreio/csfw/net/signed/backendsigned"
	"github.com/corestoreio/csfw/util/errors"
)

// OptionName identifies this package within the register of the
// backendsigned.Backend type.
const OptionName = `sha256`

// NewOptionFactory creates a new option factory function for the memstore in the
// backend package to be used for automatic scope based configuration
// initialization. Configuration values are read from argument `be`.
func NewOptionFactory(be *backendsigned.Configuration) (string, signed.OptionFactoryFunc) {
	return OptionName, func(sg config.Scoped) []signed.Option {

		burst, _, err := be.Burst.Get(sg)
		if err != nil {
			return signed.OptionsError(errors.Wrap(err, "[memstore] RateLimitBurst.Get"))
		}
		req, _, err := be.Requests.Get(sg)
		if err != nil {
			return signed.OptionsError(errors.Wrap(err, "[memstore] RateLimitRequests.Get"))
		}
		durRaw, _, err := be.Duration.Get(sg)
		if err != nil {
			return signed.OptionsError(errors.Wrap(err, "[memstore] RateLimitDuration.Get"))
		}

		if len(durRaw) != 1 {
			return signed.OptionsError(errors.NewFatalf("[memstore] RateLimitDuration invalid character count: %q. Should be one character long.", durRaw))
		}

		dur := rune(durRaw[0])

		useInMemMaxKeys, scpHash, err := be.StorageGcraMaxMemoryKeys.Get(sg)
		if err != nil {
			return signed.OptionsError(errors.Wrap(err, "[memstore] RateLimitStorageGcraMaxMemoryKeys.Get"))
		} else if useInMemMaxKeys > 0 {
			scp, scpID := scpHash.Unpack()
			return []signed.Option{
				WithGCRA(scp, scpID, useInMemMaxKeys, dur, req, burst),
			}
		}
		return signed.OptionsError(errors.NewEmptyf("[memstore] Memstore not active because RateLimitStorageGcraMaxMemoryKeys is %d.", useInMemMaxKeys))
	}
}
