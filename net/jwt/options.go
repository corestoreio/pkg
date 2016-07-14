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
	"time"

	"github.com/corestoreio/csfw/log"
	"github.com/corestoreio/csfw/store"
	"github.com/corestoreio/csfw/store/scope"
	"github.com/corestoreio/csfw/util/csjwt"
	"github.com/corestoreio/csfw/util/errors"
)

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
	return withDefaultConfig(scp, id)
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
		s.rwmu.Lock()
		defer s.rwmu.Unlock()

		sc := s.scopeCache[h]
		if sc == nil {
			sc = optionInheritDefault(s)
		}
		sc.templateTokenFunc = f
		sc.ScopeHash = h
		s.scopeCache[h] = sc
		return nil
	}
}

// WithSigningMethod this option function lets you overwrite the default 256 bit
// signing method for a specific scope. Used incorrectly token decryption can fail.
func WithSigningMethod(scp scope.Scope, id int64, sm csjwt.Signer) Option {
	h := scope.NewHash(scp, id)
	return func(s *Service) error {
		s.rwmu.Lock()
		defer s.rwmu.Unlock()

		sc := s.scopeCache[h]
		if sc == nil {
			sc = optionInheritDefault(s)
		}
		sc.SigningMethod = sm
		sc.Verifier = csjwt.NewVerification(sm)
		sc.initKeyFunc()
		sc.ScopeHash = h
		s.scopeCache[h] = sc
		return nil
	}
}

// WithExpiration sets expiration duration depending on the scope
func WithExpiration(scp scope.Scope, id int64, d time.Duration) Option {
	h := scope.NewHash(scp, id)
	return func(s *Service) error {
		s.rwmu.Lock()
		defer s.rwmu.Unlock()

		sc := s.scopeCache[h]
		if sc == nil {
			sc = optionInheritDefault(s)
		}
		sc.Expire = d
		sc.ScopeHash = h
		s.scopeCache[h] = sc
		return nil
	}
}

// WithSkew sets the duration of time skew we allow between signer and verifier.
// Must be a positive value.
func WithSkew(scp scope.Scope, id int64, d time.Duration) Option {
	h := scope.NewHash(scp, id)
	return func(s *Service) error {
		s.rwmu.Lock()
		defer s.rwmu.Unlock()

		sc := s.scopeCache[h]
		if sc == nil {
			sc = optionInheritDefault(s)
		}
		sc.Skew = d
		sc.ScopeHash = h
		s.scopeCache[h] = sc
		return nil
	}
}

// WithTokenID enables JTI (JSON Web Token ID) for a specific scope
func WithTokenID(scp scope.Scope, id int64, enable bool) Option {
	h := scope.NewHash(scp, id)
	return func(s *Service) error {
		s.rwmu.Lock()
		defer s.rwmu.Unlock()

		sc := s.scopeCache[h]
		if sc == nil {
			sc = optionInheritDefault(s)
		}
		sc.EnableJTI = enable
		sc.ScopeHash = h
		s.scopeCache[h] = sc
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

		sc := s.scopeCache[h]
		if sc == nil {
			sc = optionInheritDefault(s)
		}
		sc.ScopeHash = h

		// if you are not satisfied with the bit size of 256 you can change it
		// by using WithSigningMethod
		switch key.Algorithm() {
		case csjwt.ES:
			sc.SigningMethod = csjwt.NewSigningMethodES256()
		case csjwt.HS:
			sc.SigningMethod, err = csjwt.NewHMACFast256(key)
			if err != nil {
				return errors.Wrap(err, "[jwt] HMAC Fast 256 error")
			}
		case csjwt.RS:
			sc.SigningMethod = csjwt.NewSigningMethodRS256()
		default:
			return errors.NewNotImplementedf(errUnknownSigningMethodOptions, key.Algorithm())
		}

		sc.Key = key
		sc.Verifier = csjwt.NewVerification(sc.SigningMethod)
		sc.initKeyFunc()

		s.scopeCache[h] = sc
		return nil
	}
}

// WithDisable disables the whole JWT processing for a scope.
func WithDisable(scp scope.Scope, id int64, isDisabled bool) Option {
	h := scope.NewHash(scp, id)
	return func(s *Service) error {
		s.rwmu.Lock()
		defer s.rwmu.Unlock()

		sc := s.scopeCache[h]
		if sc == nil {
			sc = optionInheritDefault(s)
		}
		sc.Disabled = isDisabled
		sc.ScopeHash = h
		s.scopeCache[h] = sc
		return nil
	}
}
