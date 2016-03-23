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
	"github.com/corestoreio/csfw/config/cfgmodel"
	"github.com/corestoreio/csfw/config/element"
)

// PkgBackend just exported for the sake of documentation. See fields
// for more information. The PkgBackend handles the reading and writing
// of configuration values within this package.
type PkgBackend struct {
	cfgmodel.PkgBackend

	// NetCtxjwtHmacPassword handles the password. Will panic if you
	// do not set the cfgmodel.Encryptor
	// Path: net/ctxjwt/signing_method @todo
	NetCtxjwtSigningMethod ConfigSigningMethod

	// NetCtxjwtHmacPassword handles the password. Will panic if you
	// do not set the cfgmodel.Encryptor
	// Path: net/ctxjwt/hmac_password
	NetCtxjwtHmacPassword cfgmodel.Obscure

	// NetCtxjwtRSAKey handles the RSA private key. Will panic if you
	// do not set the cfgmodel.Encryptor
	// Path: net/ctxjwt/rsa_key @todo implement
	NetCtxjwtRSAKey cfgmodel.Obscure

	// NetCtxjwtRSAKeyPassword handles the password for the RSA private key.
	// Will panic if you do not set the cfgmodel.Encryptor
	// Path: net/ctxjwt/rsa_key_password @todo implement
	NetCtxjwtRSAKeyPassword cfgmodel.Obscure

	// NetCtxjwtECDSAKey handles the ECDSA private key. Will panic if you
	// do not set the cfgmodel.Encryptor
	// Path: net/ctxjwt/ecdsa_key @todo implement
	NetCtxjwtECDSAKey cfgmodel.Obscure

	// NetCtxjwtECDSAKeyPassword handles the password for the ECDSA private key.
	// Will panic if you do not set the cfgmodel.Encryptor
	// Path: net/ctxjwt/ecdsa_key_password @todo implement
	NetCtxjwtECDSAKeyPassword cfgmodel.Obscure
}

// NewBackend initializes the global configuration models containing the
// cfgpath.Route variable to the appropriate entry.
// The function Load() will be executed to apply the SectionSlice
// to all models. See Load() for more details.
func NewBackend(cfgStruct element.SectionSlice) *PkgBackend {
	return (&PkgBackend{}).Load(cfgStruct)
}

// Load creates the configuration models for each PkgBackend field.
// Internal mutex will protect the fields during loading.
// The argument SectionSlice will be applied to all models.
func (pp *PkgBackend) Load(cfgStruct element.SectionSlice) *PkgBackend {
	pp.Lock()
	defer pp.Unlock()

	opt := cfgmodel.WithFieldFromSectionSlice(cfgStruct)

	pp.NetCtxjwtSigningMethod = NewConfigSigningMethod(`net/ctxjwt/signing_method`, opt)
	pp.NetCtxjwtHmacPassword = cfgmodel.NewObscure(`net/ctxjwt/hmac_password`, opt)
	pp.NetCtxjwtRSAKey = cfgmodel.NewObscure(`net/ctxjwt/rsa_key`, opt)
	pp.NetCtxjwtRSAKeyPassword = cfgmodel.NewObscure(`net/ctxjwt/rsa_key_password`, opt)
	pp.NetCtxjwtECDSAKey = cfgmodel.NewObscure(`net/ctxjwt/ecdsa_key`, opt)
	pp.NetCtxjwtECDSAKeyPassword = cfgmodel.NewObscure(`net/ctxjwt/ecdsa_key_password`, opt)

	return pp
}
