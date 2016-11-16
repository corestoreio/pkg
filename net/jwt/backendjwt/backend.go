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
	"github.com/corestoreio/csfw/config/cfgsource"
	"github.com/corestoreio/csfw/config/element"
	"github.com/corestoreio/csfw/net/jwt"
)

// Configuration just exported for the sake of documentation. See fields for
// more information.
type Configuration struct {
	*jwt.OptionFactories

	// Disabled if set to true disables the JWT validation.
	// Path: net/jwt/disabled
	Disabled cfgmodel.Bool

	// HmacPassword handles the password. Will panic if you
	// do not set the cfgmodel.Encryptor
	// Path: net/jwt/signing_method
	SigningMethod ConfigSigningMethod

	// Skew defines the time skew duration between verifier and signer.
	// Path: net/jwt/skew
	Skew cfgmodel.Duration

	// Expiration defines the duration in which a token expires.
	// Path: net/jwt/expiration
	Expiration cfgmodel.Duration

	// SingleTokenUsage if enabled a token can only be used once per request.
	// Path: net/jwt/single_usage
	SingleTokenUsage cfgmodel.Bool

	// HmacPassword handles the password. Will panic if you
	// do not set the cfgmodel.Encryptor
	// Path: net/jwt/hmac_password
	HmacPassword cfgmodel.Obscure

	// HmacPasswordPerUser if enable each logged in user will have their own
	// randomly generated password.
	// TODO(cs) think and implement. we also may need a map to map a user to his/her password and a 2nd field in config which defines the claim key for the username.
	// Path: net/jwt/hmac_password_per_user
	HmacPasswordPerUser cfgmodel.Bool

	// RSAKey handles the RSA private key. Will panic if you
	// do not set the cfgmodel.Encryptor
	// Path: net/jwt/rsa_key
	RSAKey cfgmodel.Obscure

	// RSAKeyPassword handles the password for the RSA private key.
	// Will panic if you do not set the cfgmodel.Encryptor
	// Path: net/jwt/rsa_key_password
	RSAKeyPassword cfgmodel.Obscure

	// ECDSAKey handles the ECDSA private key. Will panic if you
	// do not set the cfgmodel.Encryptor
	// Path: net/jwt/ecdsa_key
	ECDSAKey cfgmodel.Obscure

	// ECDSAKeyPassword handles the password for the ECDSA private key.
	// Will panic if you do not set the cfgmodel.Encryptor
	// Path: net/jwt/ecdsa_key_password
	ECDSAKeyPassword cfgmodel.Obscure
}

// New initializes the backend configuration models containing the cfgpath.Route
// variable to the appropriate entries in the storage. The argument SectionSlice
// and opts will be applied to all models.
func New(cfgStruct element.SectionSlice, opts ...cfgmodel.Option) *Configuration {
	be := &Configuration{
		OptionFactories: jwt.NewOptionFactories(),
	}

	opts = append(opts, cfgmodel.WithFieldFromSectionSlice(cfgStruct))

	be.Disabled = cfgmodel.NewBool(`net/jwt/disabled`, append(opts, cfgmodel.WithSource(cfgsource.EnableDisable))...)
	be.SigningMethod = NewConfigSigningMethod(`net/jwt/signing_method`, opts...)
	be.Expiration = cfgmodel.NewDuration(`net/jwt/expiration`, opts...)
	be.Skew = cfgmodel.NewDuration(`net/jwt/skew`, opts...)
	be.SingleTokenUsage = cfgmodel.NewBool(`net/jwt/single_usage`, append(opts, cfgmodel.WithSource(cfgsource.EnableDisable))...)
	be.HmacPassword = cfgmodel.NewObscure(`net/jwt/hmac_password`, opts...)
	be.HmacPasswordPerUser = cfgmodel.NewBool(`net/jwt/hmac_password_per_user`, append(opts, cfgmodel.WithSource(cfgsource.EnableDisable))...)
	be.RSAKey = cfgmodel.NewObscure(`net/jwt/rsa_key`, opts...)
	be.RSAKeyPassword = cfgmodel.NewObscure(`net/jwt/rsa_key_password`, opts...)
	be.ECDSAKey = cfgmodel.NewObscure(`net/jwt/ecdsa_key`, opts...)
	be.ECDSAKeyPassword = cfgmodel.NewObscure(`net/jwt/ecdsa_key_password`, opts...)

	return be
}
