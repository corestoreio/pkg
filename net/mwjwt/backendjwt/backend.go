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
)

// Backend just exported for the sake of documentation. See fields
// for more information. The PkgBackend handles the reading and writing
// of configuration values within this package.
type Backend struct {
	cfgmodel.PkgBackend

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

	// NetJwtEnableJTI if enabled a new token ID will be generated.
	// Path: net/jwt/enable_jti
	NetJwtEnableJTI cfgmodel.Bool

	// NetJwtHmacPassword handles the password. Will panic if you
	// do not set the cfgmodel.Encryptor
	// Path: net/jwt/hmac_password
	NetJwtHmacPassword cfgmodel.Obscure

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

// New initializes the backend configuration models containing the
// cfgpath.Route variable to the appropriate entries.
// The function Load() will be executed to apply the SectionSlice
// to all models. See Load() for more details.
func New(cfgStruct element.SectionSlice, opts ...cfgmodel.Option) *Backend {
	return (&Backend{}).Load(cfgStruct, opts...)
}

// Load creates the configuration models for each PkgBackend field.
// Internal mutex will protect the fields during loading.
// The argument SectionSlice will be applied to all models.
// Obscure types needs the cfgmodel.Encryptor to be set.
func (pp *Backend) Load(cfgStruct element.SectionSlice, opts ...cfgmodel.Option) *Backend {
	pp.Lock()
	defer pp.Unlock()

	opts = append(opts, cfgmodel.WithFieldFromSectionSlice(cfgStruct))

	pp.NetJwtSigningMethod = NewConfigSigningMethod(`net/jwt/signing_method`, opts...)
	pp.NetJwtExpiration = cfgmodel.NewDuration(`net/jwt/expiration`, opts...)
	pp.NetJwtSkew = cfgmodel.NewDuration(`net/jwt/skew`, opts...)
	pp.NetJwtEnableJTI = cfgmodel.NewBool(`net/jwt/enable_jti`, append(opts, cfgmodel.WithSource(source.EnableDisable))...)
	pp.NetJwtHmacPassword = cfgmodel.NewObscure(`net/jwt/hmac_password`, opts...)
	pp.NetJwtRSAKey = cfgmodel.NewObscure(`net/jwt/rsa_key`, opts...)
	pp.NetJwtRSAKeyPassword = cfgmodel.NewObscure(`net/jwt/rsa_key_password`, opts...)
	pp.NetJwtECDSAKey = cfgmodel.NewObscure(`net/jwt/ecdsa_key`, opts...)
	pp.NetJwtECDSAKeyPassword = cfgmodel.NewObscure(`net/jwt/ecdsa_key_password`, opts...)

	return pp
}
