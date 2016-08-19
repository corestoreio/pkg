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
	"github.com/corestoreio/csfw/config/cfgmodel"
	"github.com/corestoreio/csfw/config/element"
	"github.com/corestoreio/csfw/config/source"
	"github.com/corestoreio/csfw/net/jwt"
)

// Configuration just exported for the sake of documentation. See fields for
// more information.
type Configuration struct {
	*jwt.OptionFactories

	// NetJwtDisabled if set to true disables the JWT validation.
	// Path: net/jwt/disabled
	NetJwtDisabled cfgmodel.Bool

	// NetJwtHmacPassword handles the password. Will panic if you
	// do not set the cfgmodel.Encryptor
	// Path: net/jwt/signing_method
	NetJwtSigningMethod ConfigSigningMethod

	// NetJwtSkew defines the time skew duration between verifier and signer.
	// Path: net/jwt/skew
	NetJwtSkew cfgmodel.Duration

	// NetJwtExpiration defines the duration in which a token expires.
	// Path: net/jwt/expiration
	NetJwtExpiration cfgmodel.Duration

	// NetJwtSingleTokenUsage if enabled a token can only be used once per request.
	// Path: net/jwt/single_usage
	NetJwtSingleTokenUsage cfgmodel.Bool

	// NetJwtHmacPassword handles the password. Will panic if you
	// do not set the cfgmodel.Encryptor
	// Path: net/jwt/hmac_password
	NetJwtHmacPassword cfgmodel.Obscure

	// NetJwtHmacPasswordPerUser if enable each logged in user will have their own
	// randomly generated password.
	// TODO(cs) think and implement. we also may need a map to map a user to his/her password and a 2nd field in config which defines the claim key for the username.
	// Path: net/jwt/hmac_password_per_user
	NetJwtHmacPasswordPerUser cfgmodel.Bool

	// NetJwtRSAKey handles the RSA private key. Will panic if you
	// do not set the cfgmodel.Encryptor
	// Path: net/jwt/rsa_key
	NetJwtRSAKey cfgmodel.Obscure

	// NetJwtRSAKeyPassword handles the password for the RSA private key.
	// Will panic if you do not set the cfgmodel.Encryptor
	// Path: net/jwt/rsa_key_password
	NetJwtRSAKeyPassword cfgmodel.Obscure

	// NetJwtECDSAKey handles the ECDSA private key. Will panic if you
	// do not set the cfgmodel.Encryptor
	// Path: net/jwt/ecdsa_key
	NetJwtECDSAKey cfgmodel.Obscure

	// NetJwtECDSAKeyPassword handles the password for the ECDSA private key.
	// Will panic if you do not set the cfgmodel.Encryptor
	// Path: net/jwt/ecdsa_key_password
	NetJwtECDSAKeyPassword cfgmodel.Obscure
}

// New initializes the backend configuration models containing the cfgpath.Route
// variable to the appropriate entries in the storage. The argument SectionSlice
// and opts will be applied to all models.
func New(cfgStruct element.SectionSlice, opts ...cfgmodel.Option) *Configuration {
	be := &Configuration{
		OptionFactories: jwt.NewOptionFactories(),
	}

	opts = append(opts, cfgmodel.WithFieldFromSectionSlice(cfgStruct))

	be.NetJwtDisabled = cfgmodel.NewBool(`net/jwt/disabled`, append(opts, cfgmodel.WithSource(source.EnableDisable))...)
	be.NetJwtSigningMethod = NewConfigSigningMethod(`net/jwt/signing_method`, opts...)
	be.NetJwtExpiration = cfgmodel.NewDuration(`net/jwt/expiration`, opts...)
	be.NetJwtSkew = cfgmodel.NewDuration(`net/jwt/skew`, opts...)
	be.NetJwtSingleTokenUsage = cfgmodel.NewBool(`net/jwt/single_usage`, append(opts, cfgmodel.WithSource(source.EnableDisable))...)
	be.NetJwtHmacPassword = cfgmodel.NewObscure(`net/jwt/hmac_password`, opts...)
	be.NetJwtHmacPasswordPerUser = cfgmodel.NewBool(`net/jwt/hmac_password_per_user`, append(opts, cfgmodel.WithSource(source.EnableDisable))...)
	be.NetJwtRSAKey = cfgmodel.NewObscure(`net/jwt/rsa_key`, opts...)
	be.NetJwtRSAKeyPassword = cfgmodel.NewObscure(`net/jwt/rsa_key_password`, opts...)
	be.NetJwtECDSAKey = cfgmodel.NewObscure(`net/jwt/ecdsa_key`, opts...)
	be.NetJwtECDSAKeyPassword = cfgmodel.NewObscure(`net/jwt/ecdsa_key_password`, opts...)

	return be
}
