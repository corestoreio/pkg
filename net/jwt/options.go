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

package jwt

import (
	"net/http"
	"time"

	"github.com/corestoreio/csfw/config"
	"github.com/corestoreio/csfw/log"
	"github.com/corestoreio/csfw/store"
	"github.com/corestoreio/csfw/store/scope"
	"github.com/corestoreio/csfw/sync/singleflight"
	"github.com/corestoreio/csfw/util/csjwt"
	"github.com/corestoreio/csfw/util/errors"
)

// Option can be used as an argument in NewService to configure a token service.
type Option func(*Service) error

// OptionFactoryFunc a closure around a scoped configuration to figure out which
// options should be returned depending on the scope brought to you during
// a request.
type OptionFactoryFunc func(config.ScopedGetter) []Option

// WithDefaultConfig applies the default JWT configuration settings based for
// a specific scope.
//
// Default values are:
//		- constant DefaultExpire
//		- HMAC Password: random bytes, for each scope different.
//		- Signing Method HMAC SHA 256 (fast version from pkg csjwt)
//		- HTTP error handler returns http.StatusUnauthorized
//		- JTI disabled
func WithDefaultConfig(scp scope.Scope, id int64) Option {
	h := scope.NewHash(scp, id)
	return func(s *Service) (err error) {
		if h == scope.DefaultHash {
			s.defaultScopeCache = defaultScopedConfig()
			return nil
		}

		s.rwmu.Lock()
		defer s.rwmu.Unlock()
		s.scopeCache[h] = defaultScopedConfig()
		return nil
	}
}

// WithBlacklist sets a new global black list service. Convenience helper
// function.
func WithBlacklist(bl Blacklister) Option {
	return func(s *Service) error {
		s.Blacklist = bl
		return nil
	}
}

// WithLogger sets a new global logger. Convenience helper function.
func WithLogger(l log.Logger) Option {
	return func(s *Service) error {
		s.Log = l
		return nil
	}
}

// WithStoreService apply a store service aka. requested store to the middleware
// to allow a store change if requested via token. Convenience helper function.
func WithStoreService(sr store.Requester) Option {
	return func(s *Service) error {
		s.StoreService = sr
		return nil
	}
}

// WithTemplateToken set a custom csjwt.Header and csjwt.Claimer for each scope
// when parsing a token in a request. Function f will generate a new base token
// for each request. This allows you to choose using a slow map as a claim or a
// fast struct based claim. Same goes with the header.
func WithTemplateToken(scp scope.Scope, id int64, f func() csjwt.Token) Option {
	h := scope.NewHash(scp, id)
	return func(s *Service) error {

		if h == scope.DefaultHash {
			s.defaultScopeCache.templateTokenFunc = f
			return nil
		}

		s.rwmu.Lock()
		defer s.rwmu.Unlock()

		// inherit default config
		scNew := s.defaultScopeCache
		scNew.templateTokenFunc = f

		if sc, ok := s.scopeCache[h]; ok {
			sc.templateTokenFunc = scNew.templateTokenFunc
			scNew = sc
		}
		scNew.ScopeHash = h
		s.scopeCache[h] = scNew
		return nil
	}
}

// WithSigningMethod this option function lets you overwrite the default 256 bit
// signing method for a specific scope. Used incorrectly token decryption can fail.
func WithSigningMethod(scp scope.Scope, id int64, sm csjwt.Signer) Option {
	h := scope.NewHash(scp, id)
	return func(s *Service) error {

		if h == scope.DefaultHash {
			s.defaultScopeCache.SigningMethod = sm
			s.defaultScopeCache.Verifier = csjwt.NewVerification(sm)
			return nil
		}

		s.rwmu.Lock()
		defer s.rwmu.Unlock()

		// inherit default config
		scNew := s.defaultScopeCache

		scNew.SigningMethod = sm
		scNew.Verifier = csjwt.NewVerification(sm)
		scNew.initKeyFunc()

		if sc, ok := s.scopeCache[h]; ok {
			sc.SigningMethod = scNew.SigningMethod
			sc.Verifier = scNew.Verifier
			sc.KeyFunc = scNew.KeyFunc
			scNew = sc
		}

		scNew.ScopeHash = h
		s.scopeCache[h] = scNew
		return nil
	}
}

// WithErrorHandler sets the error handler for a scope and its ID. If the
// scope.DefaultID will be set the handler gets also applied to the global
// handler
func WithErrorHandler(scp scope.Scope, id int64, handler http.Handler) Option {
	h := scope.NewHash(scp, id)
	return func(s *Service) error {

		if h == scope.DefaultHash {
			s.defaultScopeCache.ErrorHandler = handler
			return nil
		}

		s.rwmu.Lock()
		defer s.rwmu.Unlock()

		// inherit default config
		scNew := s.defaultScopeCache
		scNew.ErrorHandler = handler

		if sc, ok := s.scopeCache[h]; ok {
			sc.ErrorHandler = scNew.ErrorHandler
			scNew = sc
		}
		scNew.ScopeHash = h
		s.scopeCache[h] = scNew
		return nil
	}
}

// WithExpiration sets expiration duration depending on the scope
func WithExpiration(scp scope.Scope, id int64, d time.Duration) Option {
	h := scope.NewHash(scp, id)
	return func(s *Service) error {

		if h == scope.DefaultHash {
			s.defaultScopeCache.Expire = d
			return nil
		}

		s.rwmu.Lock()
		defer s.rwmu.Unlock()

		// inherit default config
		scNew := s.defaultScopeCache
		scNew.Expire = d

		if sc, ok := s.scopeCache[h]; ok {
			sc.Expire = scNew.Expire
			scNew = sc
		}
		scNew.ScopeHash = h
		s.scopeCache[h] = scNew
		return nil
	}
}

// WithSkew sets the duration of time skew we allow between signer and verifier.
// Must be a positive value.
func WithSkew(scp scope.Scope, id int64, d time.Duration) Option {
	h := scope.NewHash(scp, id)
	return func(s *Service) error {

		if h == scope.DefaultHash {
			s.defaultScopeCache.Skew = d
			return nil
		}

		s.rwmu.Lock()
		defer s.rwmu.Unlock()

		// inherit default config
		scNew := s.defaultScopeCache
		scNew.Skew = d

		if sc, ok := s.scopeCache[h]; ok {
			sc.Skew = scNew.Skew
			scNew = sc
		}
		scNew.ScopeHash = h
		s.scopeCache[h] = scNew
		return nil
	}
}

// WithTokenID enables JTI (JSON Web Token ID) for a specific scope
func WithTokenID(scp scope.Scope, id int64, enable bool) Option {
	h := scope.NewHash(scp, id)
	return func(s *Service) error {

		if h == scope.DefaultHash {
			s.defaultScopeCache.EnableJTI = enable
			return nil
		}

		s.rwmu.Lock()
		defer s.rwmu.Unlock()

		// inherit default config
		scNew := s.defaultScopeCache
		scNew.EnableJTI = enable

		if sc, ok := s.scopeCache[h]; ok {
			sc.EnableJTI = scNew.EnableJTI
			scNew = sc
		}
		scNew.ScopeHash = h
		s.scopeCache[h] = scNew
		return nil
	}
}

// WithKey sets the key for the default signing method of 256 bits.
// You can also provide your own signing method by using additionally
// the function WithSigningMethod(), which must be called after this function :-/.
func WithKey(scp scope.Scope, id int64, key csjwt.Key) Option {
	h := scope.NewHash(scp, id)
	if key.Error != nil {
		return func(s *Service) error {
			return errors.Wrap(key.Error, "[jwt] Key Error")
		}
	}
	if key.IsEmpty() {
		return func(s *Service) error {
			return errors.NewEmptyf(errKeyEmpty)
		}
	}
	return func(s *Service) (err error) {
		s.rwmu.Lock()
		defer s.rwmu.Unlock()

		// inherit default config
		scNew := s.defaultScopeCache
		scNew.Key = key

		// if you are not satisfied with the bit size of 256 you can change it
		// by using WithSigningMethod
		switch key.Algorithm() {
		case csjwt.ES:
			scNew.SigningMethod = csjwt.NewSigningMethodES256()
		case csjwt.HS:
			scNew.SigningMethod, err = csjwt.NewHMACFast256(key)
			if err != nil {
				return errors.Wrap(err, "[jwt] HMAC Fast 256 error")
			}
		case csjwt.RS:
			scNew.SigningMethod = csjwt.NewSigningMethodRS256()
		default:
			return errors.NewNotImplementedf(errUnknownSigningMethodOptions, key.Algorithm())
		}

		scNew.Verifier = csjwt.NewVerification(scNew.SigningMethod)
		scNew.initKeyFunc()

		if h == scope.DefaultHash {
			s.defaultScopeCache.Key = scNew.Key
			s.defaultScopeCache.SigningMethod = scNew.SigningMethod
			s.defaultScopeCache.Verifier = scNew.Verifier
			s.defaultScopeCache.KeyFunc = scNew.KeyFunc
			return nil
		}

		if sc, ok := s.scopeCache[h]; ok {
			sc.Key = scNew.Key
			sc.SigningMethod = scNew.SigningMethod
			sc.Verifier = scNew.Verifier
			sc.KeyFunc = scNew.KeyFunc
			scNew = sc
		}

		scNew.ScopeHash = h
		s.scopeCache[h] = scNew
		return nil
	}
}

// WithDisable disables the whole JWT processing for a scope.
func WithDisable(scp scope.Scope, id int64, ok bool) Option {
	h := scope.NewHash(scp, id)
	return func(s *Service) error {

		if h == scope.DefaultHash {
			s.defaultScopeCache.Disabled = ok
			return nil
		}

		s.rwmu.Lock()
		defer s.rwmu.Unlock()

		// inherit default config
		scNew := s.defaultScopeCache
		scNew.Disabled = ok

		if sc, ok := s.scopeCache[h]; ok {
			sc.Disabled = scNew.Disabled
			scNew = sc
		}
		scNew.ScopeHash = h
		s.scopeCache[h] = scNew
		return nil
	}
}

// WithOptionFactory applies a function which lazily loads the options depending
// on the incoming scope within an HTTP request. For example applies the backend
// configuration to the service.
//
// In the case of the jwt package the configuration will also be used when
// calling the functions ConfigByScopeID(), NewToken(), Parse(), ParseScoped().
//
// Once this option function has been set, all other manually set option
// functions, which accept a scope and a scope ID as an argument, will be
// overwritten by the new values retrieved from the configuration service.
//
//	cfgStruct, err := backendjwt.NewConfigStructure()
//	if err != nil {
//		panic(err)
//	}
//	pb := backendjwt.New(cfgStruct)
//
//	jwts := jwt.MustNewService(
//		jwt.WithOptionFactory(backendjwt.PrepareOptions(pb), configService),
//	)
func WithOptionFactory(f OptionFactoryFunc, rootConfig config.Getter) Option {
	return func(s *Service) error {
		s.rootConfig = rootConfig
		s.optionFactoryFunc = f
		s.optionInflight = new(singleflight.Group)
		return nil
	}
}
