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

package ctxjwt

import (
	"github.com/corestoreio/csfw/config"
	"github.com/corestoreio/csfw/config/cfgmodel"
	"github.com/corestoreio/csfw/util/csjwt"
	"github.com/juju/errors"
)

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
				csjwt.SigningMethodRS256.Alg(), "RSA 256",
				csjwt.SigningMethodRS384.Alg(), "RSA 384",
				csjwt.SigningMethodRS512.Alg(), "RSA 512",

				csjwt.SigningMethodES256.Alg(), "ECDSA 256",
				csjwt.SigningMethodES384.Alg(), "ECDSA 384",
				csjwt.SigningMethodES512.Alg(), "ECDSA 512",

				csjwt.SigningMethodHS256.Alg(), "HMAC-SHA 256",
				csjwt.SigningMethodHS384.Alg(), "HMAC-SHA 384",
				csjwt.SigningMethodHS512.Alg(), "HMAC-SHA 512",
			),
		)...),
	}
}

// Get returns a signing method definied for a scope.
func (cc ConfigSigningMethod) Get(sg config.ScopedGetter) (sm csjwt.Signer, err error) {
	raw, err := cc.Str.Get(sg)
	if err != nil {
		err = errors.Mask(err)
		return
	}

	if raw == "" {
		raw = DefaultSigningMethod
	}

	switch raw {
	case csjwt.SigningMethodRS256.Alg():
		sm = csjwt.SigningMethodRS256
	case csjwt.SigningMethodRS384.Alg():
		sm = csjwt.SigningMethodRS384
	case csjwt.SigningMethodRS512.Alg():
		sm = csjwt.SigningMethodRS512

	case csjwt.SigningMethodES256.Alg():
		sm = csjwt.SigningMethodES256
	case csjwt.SigningMethodES384.Alg():
		sm = csjwt.SigningMethodES384
	case csjwt.SigningMethodES512.Alg():
		sm = csjwt.SigningMethodES512

	case csjwt.SigningMethodHS256.Alg():
		sm = csjwt.SigningMethodHS256
	case csjwt.SigningMethodHS384.Alg():
		sm = csjwt.SigningMethodHS384
	case csjwt.SigningMethodHS512.Alg():
		sm = csjwt.SigningMethodHS512
	default:
		err = errors.Errorf("ctxjwt.ConfigSigningMethod: Unknown algorithm %s", raw)
	}
	return
}
