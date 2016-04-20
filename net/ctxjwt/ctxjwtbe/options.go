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

package ctxjwtbe

import (
	"github.com/corestoreio/csfw/config"
	"github.com/corestoreio/csfw/config/cfgmodel"
	"github.com/corestoreio/csfw/net/ctxjwt"
	"github.com/corestoreio/csfw/util/csjwt"
	"github.com/juju/errors"
)

// DefaultBackend creates new ctxjwt.Option slice with the default configuration
// structure and a noop encryptor/decryptor IF no option arguments have been
// provided. It panics on error, so us it only during the app init phase.
func DefaultBackend(opts ...cfgmodel.Option) ctxjwt.ScopedOption {
	cfgStruct, err := NewConfigStructure()
	if err != nil {
		panic(err)
	}
	if len(opts) == 0 {
		opts = append(opts, cfgmodel.WithEncryptor(cfgmodel.NoopEncryptor{}))
	}

	return BackendOptions(NewBackend(cfgStruct, opts...))
}

// BackendOptions creates a closure around the PkgBackend. The closure will
// be used during a scoped request to figure out the configuration depending on
// the scope. An option array will be returned by the closure.
func BackendOptions(be *PkgBackend) ctxjwt.ScopedOption {

	return func(sg config.ScopedGetter) (opts []ctxjwt.Option) {

		scp, id := sg.Scope()

		exp, err := be.NetCtxjwtExpiration.Get(sg)
		if err != nil {
			return append(opts, func(s *ctxjwt.Service) {
				s.MultiErr = s.AppendErrors(errors.Mask(err))
			})
		}
		opts = append(opts, ctxjwt.WithExpiration(scp, id, exp))

		isJTI, err := be.NetCtxjwtEnableJTI.Get(sg)
		if err != nil {
			return append(opts, func(s *ctxjwt.Service) {
				s.MultiErr = s.AppendErrors(errors.Mask(err))
			})
		}
		opts = append(opts, ctxjwt.WithTokenID(scp, id, isJTI))

		signingMethod, err := be.NetCtxjwtSigningMethod.Get(sg)
		if err != nil {
			return append(opts, func(s *ctxjwt.Service) {
				s.MultiErr = s.AppendErrors(errors.Mask(err))
			})
		}

		var key csjwt.Key

		switch signingMethod.Alg() {
		case csjwt.RS256, csjwt.RS384, csjwt.RS512:

			rsaKey, err := be.NetCtxjwtRSAKey.Get(sg)
			if err != nil {
				return append(opts, func(s *ctxjwt.Service) {
					s.MultiErr = s.AppendErrors(errors.Mask(err))
				})
			}
			rsaPW, err := be.NetCtxjwtRSAKeyPassword.Get(sg)
			if err != nil {
				return append(opts, func(s *ctxjwt.Service) {
					s.MultiErr = s.AppendErrors(errors.Mask(err))
				})
			}
			key = csjwt.WithRSAPrivateKeyFromPEM(rsaKey, rsaPW)

		case csjwt.ES256, csjwt.ES384, csjwt.ES512:

			ecdsaKey, err := be.NetCtxjwtECDSAKey.Get(sg)
			if err != nil {
				return append(opts, func(s *ctxjwt.Service) {
					s.MultiErr = s.AppendErrors(errors.Mask(err))
				})
			}
			ecdsaPW, err := be.NetCtxjwtECDSAKeyPassword.Get(sg)
			if err != nil {
				return append(opts, func(s *ctxjwt.Service) {
					s.MultiErr = s.AppendErrors(errors.Mask(err))
				})
			}
			key = csjwt.WithECPrivateKeyFromPEM(ecdsaKey, ecdsaPW)

		case csjwt.HS256, csjwt.HS384, csjwt.HS512:

			password, err := be.NetCtxjwtHmacPassword.Get(sg)
			if err != nil {
				return append(opts, func(s *ctxjwt.Service) {
					s.MultiErr = s.AppendErrors(errors.Mask(err))
				})
			}
			key = csjwt.WithPassword(password)

		default:
			opts = append(opts, func(s *ctxjwt.Service) {
				s.MultiErr = s.AppendErrors(errors.Errorf("[ctxjwt] Unknown signing method: %q", signingMethod.Alg()))
			})
		}

		// WithSigningMethod must be added at the end of the slice to overwrite default signing methods
		return append(opts, ctxjwt.WithKey(scp, id, key), ctxjwt.WithSigningMethod(scp, id, signingMethod))
	}
}
