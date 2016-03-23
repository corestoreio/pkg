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
	"github.com/dgrijalva/jwt-go"
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
				jwt.SigningMethodRS256.Alg(), "RSA 256",
				jwt.SigningMethodRS384.Alg(), "RSA 384",
				jwt.SigningMethodRS512.Alg(), "RSA 512",

				jwt.SigningMethodES256.Alg(), "ECDSA 256",
				jwt.SigningMethodES384.Alg(), "ECDSA 384",
				jwt.SigningMethodES512.Alg(), "ECDSA 512",

				jwt.SigningMethodHS256.Alg(), "HMAC-SHA 256",
				jwt.SigningMethodHS384.Alg(), "HMAC-SHA 384",
				jwt.SigningMethodHS512.Alg(), "HMAC-SHA 512",
			),
		)...),
	}
}

// Get returns a signing method definied for a scope.
func (cc ConfigSigningMethod) Get(sg config.ScopedGetter) (sm jwt.SigningMethod, err error) {
	raw, err := cc.Str.Get(sg)
	if err != nil {
		err = errors.Mask(err)
		return
	}

	if raw == "" {
		raw = DefaultSigningMethod
	}

	switch raw {
	case jwt.SigningMethodRS256.Alg():
		sm = jwt.SigningMethodRS256
	case jwt.SigningMethodRS384.Alg():
		sm = jwt.SigningMethodRS384
	case jwt.SigningMethodRS512.Alg():
		sm = jwt.SigningMethodRS512

	case jwt.SigningMethodES256.Alg():
		sm = jwt.SigningMethodES256
	case jwt.SigningMethodES384.Alg():
		sm = jwt.SigningMethodES384
	case jwt.SigningMethodES512.Alg():
		sm = jwt.SigningMethodES512

	case jwt.SigningMethodHS256.Alg():
		sm = jwt.SigningMethodHS256
	case jwt.SigningMethodHS384.Alg():
		sm = jwt.SigningMethodHS384
	case jwt.SigningMethodHS512.Alg():
		sm = jwt.SigningMethodHS512
	default:
		err = errors.Errorf("ctxjwt.ConfigSigningMethod: Unknown algorithm %s", raw)
	}
	return
}
