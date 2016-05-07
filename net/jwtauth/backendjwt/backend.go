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

	// NetCtxjwtHmacPassword handles the password. Will panic if you
	// do not set the cfgmodel.Encryptor
	// Path: net/jwtauth/signing_method
	NetCtxjwtSigningMethod ConfigSigningMethod

	// NetCtxjwtExpiration defines the duration in which a token expires.
	// Path: net/jwtauth/expiration
	NetCtxjwtExpiration cfgmodel.Duration

	// NetCtxjwtEnableJTI if enabled a new token ID will be generated.
	// Path: net/jwtauth/enable_jti
	NetCtxjwtEnableJTI cfgmodel.Bool

	// NetCtxjwtHmacPassword handles the password. Will panic if you
	// do not set the cfgmodel.Encryptor
	// Path: net/jwtauth/hmac_password
	NetCtxjwtHmacPassword cfgmodel.Obscure

	// NetCtxjwtRSAKey handles the RSA private key. Will panic if you
	// do not set the cfgmodel.Encryptor
	// Path: net/jwtauth/rsa_key
	NetCtxjwtRSAKey cfgmodel.Obscure

	// NetCtxjwtRSAKeyPassword handles the password for the RSA private key.
	// Will panic if you do not set the cfgmodel.Encryptor
	// Path: net/jwtauth/rsa_key_password
	NetCtxjwtRSAKeyPassword cfgmodel.Obscure

	// NetCtxjwtECDSAKey handles the ECDSA private key. Will panic if you
	// do not set the cfgmodel.Encryptor
	// Path: net/jwtauth/ecdsa_key
	NetCtxjwtECDSAKey cfgmodel.Obscure

	// NetCtxjwtECDSAKeyPassword handles the password for the ECDSA private key.
	// Will panic if you do not set the cfgmodel.Encryptor
	// Path: net/jwtauth/ecdsa_key_password
	NetCtxjwtECDSAKeyPassword cfgmodel.Obscure
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

	pp.NetCtxjwtSigningMethod = NewConfigSigningMethod(`net/jwtauth/signing_method`, opts...)
	pp.NetCtxjwtExpiration = cfgmodel.NewDuration(`net/jwtauth/expiration`, opts...)
	pp.NetCtxjwtEnableJTI = cfgmodel.NewBool(`net/jwtauth/enable_jti`, append(opts, cfgmodel.WithSource(source.EnableDisable))...)
	pp.NetCtxjwtHmacPassword = cfgmodel.NewObscure(`net/jwtauth/hmac_password`, opts...)
	pp.NetCtxjwtRSAKey = cfgmodel.NewObscure(`net/jwtauth/rsa_key`, opts...)
	pp.NetCtxjwtRSAKeyPassword = cfgmodel.NewObscure(`net/jwtauth/rsa_key_password`, opts...)
	pp.NetCtxjwtECDSAKey = cfgmodel.NewObscure(`net/jwtauth/ecdsa_key`, opts...)
	pp.NetCtxjwtECDSAKeyPassword = cfgmodel.NewObscure(`net/jwtauth/ecdsa_key_password`, opts...)

	return pp
}
