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

package backendcors

import (
	"github.com/corestoreio/csfw/config"
	"github.com/corestoreio/csfw/config/cfgmodel"
	"github.com/corestoreio/csfw/net/ctxcors"
)

// DefaultBackend creates new ctxcors.Option slice with the default configuration
// structure and a noop encryptor/decryptor IF no option arguments have been
// provided. It panics on error, so us it only during the app init phase.
func DefaultBackend(opts ...cfgmodel.Option) ctxcors.ScopedOptionFunc {
	cfgStruct, err := NewConfigStructure()
	if err != nil {
		panic(err)
	}
	if len(opts) == 0 {
		opts = append(opts, cfgmodel.WithEncryptor(cfgmodel.NoopEncryptor{}))
	}

	return BackendOptions(New(cfgStruct, opts...))
}

// BackendOptions creates a closure around the PkgBackend. The closure will
// be used during a scoped request to figure out the configuration depending on
// the scope. An option array will be returned by the closure.
func BackendOptions(be *Backend) ctxcors.ScopedOptionFunc {

	return func(sg config.ScopedGetter) (opts []ctxcors.Option) {

		// scp, id := sg.Scope()

		//exp, err := be.NetCtxjwtExpiration.Get(sg)
		//if err != nil {
		//	return append(opts, func(s *ctxcors.Service) {
		//		s.AddError(errors.Wrap(err, "[backendjwt] NetCtxjwtExpiration.Get"))
		//	})
		//}
		//opts = append(opts, ctxcors.WithExpiration(scp, id, exp))

		// WithSigningMethod must be added at the end of the slice to overwrite default signing methods
		// return append(opts, ctxcors.WithKey(scp, id, key), ctxcors.WithSigningMethod(scp, id, signingMethod))
		return nil
	}
}

//
//// WithBackendApplied allows to add the backend configuration struct and applying
//// all options. This option should only be used within the middleware while
//// creating a new Cors pointer for a specific scope.
//func WithBackendApplied(b *Backend, sg config.ScopedGetter) Option {
//	return func(c *Service) {
//		c.Backend = b
//
//		var opts [8]Option
//
//		headers, err := b.NetCtxcorsExposedHeaders.Get(sg)
//		if err != nil {
//			c.MultiErr = c.AppendErrors(err)
//		}
//		opts[0] = WithExposedHeaders(headers...)
//
//		ao, err := b.NetCtxcorsAllowedOrigins.Get(sg)
//		if err != nil {
//			c.MultiErr = c.AppendErrors(err)
//		}
//		opts[1] = WithAllowedOrigins(ao...)
//
//		am, err := b.NetCtxcorsAllowedMethods.Get(sg)
//		if err != nil {
//			c.MultiErr = c.AppendErrors(err)
//		}
//		opts[2] = WithAllowedMethods(am...)
//
//		ah, err := b.NetCtxcorsAllowedHeaders.Get(sg)
//		if err != nil {
//			c.MultiErr = c.AppendErrors(err)
//		}
//		opts[3] = WithAllowedHeaders(ah...)
//
//		ac, err := b.NetCtxcorsAllowCredentials.Get(sg)
//		if err != nil {
//			c.MultiErr = c.AppendErrors(err)
//		}
//		if ac {
//			opts[4] = WithAllowCredentials()
//		}
//
//		op, err := b.NetCtxcorsOptionsPassthrough.Get(sg)
//		if err != nil {
//			c.MultiErr = c.AppendErrors(err)
//		}
//		if op {
//			opts[5] = WithOptionsPassthrough()
//		}
//
//		ma, err := b.NetCtxcorsMaxAge.Get(sg)
//		if err != nil {
//			c.MultiErr = c.AppendErrors(err)
//		}
//		opts[6] = WithMaxAge(ma)
//
//		// inherit logger
//		if c.Log != nil {
//			opts[7] = WithLogger(c.Log)
//		}
//
//		_, _ = c.Options(opts[:]...) // ignore because already covered
//	}
//}
