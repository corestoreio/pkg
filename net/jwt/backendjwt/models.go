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

package backendjwt

import (
	"github.com/corestoreio/csfw/config"
	"github.com/corestoreio/csfw/config/cfgmodel"
	"github.com/corestoreio/csfw/net/jwt"
	"github.com/corestoreio/csfw/util/csjwt"
	"github.com/corestoreio/errors"
)

// todo add blake2 hash

// ConfigSigningMethod signing method type for the JWT.
type ConfigSigningMethod struct {
	cfgmodel.Str
}

// NewConfigSigningMethod creates a new signing method configuration type.
func NewConfigSigningMethod(path string, opts ...cfgmodel.Option) ConfigSigningMethod {
	return ConfigSigningMethod{
		Str: cfgmodel.NewStr(path, append(
			opts,
			cfgmodel.WithSourceByString(
				csjwt.RS256, "RSA 256",
				csjwt.RS384, "RSA 384",
				csjwt.RS512, "RSA 512",

				csjwt.ES256, "ECDSA 256",
				csjwt.ES384, "ECDSA 384",
				csjwt.ES512, "ECDSA 512",

				csjwt.HS256, "HMAC-SHA 256",
				csjwt.HS384, "HMAC-SHA 384",
				csjwt.HS512, "HMAC-SHA 512",
			),
		)...),
	}
}

// Get returns a signing method definied for a scope.
// Error behaviour: NotImplemented
func (cc ConfigSigningMethod) Get(sg config.Scoped) (csjwt.Signer, error) {
	raw, err := cc.Str.Get(sg)
	if err != nil {
		return nil, errors.Wrap(err, "[backendjwt] Str.Get")
	}

	if raw == "" {
		raw = jwt.DefaultSigningMethod
	}

	var sm csjwt.Signer
	switch raw {
	case csjwt.RS256:
		sm = csjwt.NewSigningMethodRS256()
	case csjwt.RS384:
		sm = csjwt.NewSigningMethodRS384()
	case csjwt.RS512:
		sm = csjwt.NewSigningMethodRS512()

	case csjwt.ES256:
		sm = csjwt.NewSigningMethodES256()
	case csjwt.ES384:
		sm = csjwt.NewSigningMethodES384()
	case csjwt.ES512:
		sm = csjwt.NewSigningMethodES512()

	case csjwt.HS256:
		sm = csjwt.NewSigningMethodHS256()
	case csjwt.HS384:
		sm = csjwt.NewSigningMethodHS384()
	case csjwt.HS512:
		sm = csjwt.NewSigningMethodHS512()
	default:
		err = errors.NewNotImplementedf("[jwt] ConfigSigningMethod: Unknown algorithm %s", raw)
	}
	return sm, err
}
