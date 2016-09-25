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

	"github.com/corestoreio/csfw/net/mw"
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
func WithDefaultConfig(id scope.TypeID) Option {
	return withDefaultConfig(id)
}

// WithBlacklist sets a new global black list service. Convenience helper
// function.
func WithBlacklist(bl Blacklister) Option {
	return func(s *Service) error {
		s.Blacklist = bl
		return nil
	}
}

// WithTemplateToken set a custom csjwt.Header and csjwt.Claimer for each scope
// when parsing a token in a request. Function f will generate a new base token
// for each request. This allows you to choose using a slow map as a claim or a
// fast struct based claim. Same goes with the header.
func WithTemplateToken(id scope.TypeID, f func() csjwt.Token) Option {
	return func(s *Service) error {
		s.rwmu.Lock()
		defer s.rwmu.Unlock()

		sc := s.scopeCache[id]
		if sc == nil {
			sc = optionInheritDefault(s)
		}
		sc.templateTokenFunc = f
		sc.ScopeID = id
		s.scopeCache[id] = sc
		return nil
	}
}

// WithSigningMethod this option function lets you overwrite the default 256 bit
// signing method for a specific scope. Used incorrectly token decryption can fail.
func WithSigningMethod(id scope.TypeID, sm csjwt.Signer) Option {
	return func(s *Service) error {
		s.rwmu.Lock()
		defer s.rwmu.Unlock()

		sc := s.scopeCache[id]
		if sc == nil {
			sc = optionInheritDefault(s)
		}
		sc.SigningMethod = sm
		sc.Verifier = csjwt.NewVerification(sm)
		sc.initKeyFunc()
		sc.ScopeID = id
		s.scopeCache[id] = sc
		return nil
	}
}

// WithExpiration sets expiration duration depending on the scope
func WithExpiration(id scope.TypeID, d time.Duration) Option {
	return func(s *Service) error {
		s.rwmu.Lock()
		defer s.rwmu.Unlock()

		sc := s.scopeCache[id]
		if sc == nil {
			sc = optionInheritDefault(s)
		}
		sc.Expire = d
		sc.ScopeID = id
		s.scopeCache[id] = sc
		return nil
	}
}

// WithSkew sets the duration of time skew we allow between signer and verifier.
// Must be a positive value.
func WithSkew(id scope.TypeID, d time.Duration) Option {
	return func(s *Service) error {
		s.rwmu.Lock()
		defer s.rwmu.Unlock()

		sc := s.scopeCache[id]
		if sc == nil {
			sc = optionInheritDefault(s)
		}
		sc.Skew = d
		sc.ScopeID = id
		s.scopeCache[id] = sc
		return nil
	}
}

// WithKey sets the key for the default signing method of 256 bits.
// You can also provide your own signing method by using additionally
// the function WithSigningMethod(), which must be called after this function :-/.
func WithKey(id scope.TypeID, key csjwt.Key) Option {
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

		sc := s.scopeCache[id]
		if sc == nil {
			sc = optionInheritDefault(s)
		}
		sc.ScopeID = id

		// if you are not satisfied with the bit size of 256 you can change it
		// by using WithSigningMethod
		switch key.Algorithm() {
		case csjwt.ES:
			sc.SigningMethod = csjwt.NewSigningMethodES256()
		case csjwt.HS:
			sc.SigningMethod, err = csjwt.NewSigningMethodHS256Fast(key)
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

		s.scopeCache[id] = sc
		return nil
	}
}

// WithStoreCodeFieldName sets the name of the key in the token claims section
// to extract the store code.
func WithStoreCodeFieldName(id scope.TypeID, name string) Option {
	return func(s *Service) error {
		s.rwmu.Lock()
		defer s.rwmu.Unlock()

		sc := s.scopeCache[id]
		if sc == nil {
			sc = optionInheritDefault(s)
		}
		sc.StoreCodeFieldName = name
		sc.ScopeID = id
		s.scopeCache[id] = sc
		return nil
	}
}

// WithUnauthorizedHandler adds a custom handler when a token cannot authorized to call the next handler in the chain.
// The default unauthorized handler prints the error to the user and
// returns a http.StatusUnauthorized.
func WithUnauthorizedHandler(id scope.TypeID, uh mw.ErrorHandler) Option {
	return func(s *Service) error {
		s.rwmu.Lock()
		defer s.rwmu.Unlock()

		sc := s.scopeCache[id]
		if sc == nil {
			sc = optionInheritDefault(s)
		}
		sc.UnauthorizedHandler = uh
		sc.ScopeID = id
		s.scopeCache[id] = sc
		return nil
	}
}

// WithSingleTokenUsage if set to true for each request a token can be only used
// once. The JTI (JSON Token Identifier) gets added to the blacklist until it
// expires.
func WithSingleTokenUsage(id scope.TypeID, enable bool) Option {
	return func(s *Service) error {
		s.rwmu.Lock()
		defer s.rwmu.Unlock()

		sc := s.scopeCache[id]
		if sc == nil {
			sc = optionInheritDefault(s)
		}
		sc.SingleTokenUsage = enable
		sc.ScopeID = id
		s.scopeCache[id] = sc
		return nil
	}
}
