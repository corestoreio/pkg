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

package jwtauth

import (
	"net/http"
	"time"

	"github.com/corestoreio/csfw/config"
	"github.com/corestoreio/csfw/store/scope"
	"github.com/corestoreio/csfw/util/csjwt"
	"github.com/corestoreio/csfw/util/csjwt/jwtclaim"
	"github.com/corestoreio/csfw/util/errors"
	"github.com/corestoreio/csfw/util/log"
)

// Option can be used as an argument in NewService to configure a token service.
type Option func(*Service)

// ScopedOptionFunc a closure around a scoped configuration to figure out which
// options should be returned depending on the scope brought to you during
// a request.
type ScopedOptionFunc func(config.ScopedGetter) []Option

// scopedConfig private internal scoped based configuration
type scopedConfig struct {
	// ScopeHash defines the scope bound to the configuration is.
	ScopeHash scope.Hash
	// Key contains the HMAC, RSA or ECDSA sensitive data. Yes the csjwt.Key
	// must not be embedded into this struct because otherwise when printing
	// or logging the sensitive data from csjwt.Key gets leaked into loggers
	// or where ever. If key would be lower case then %#v still prints
	// every field of the csjwt.Key.
	Key csjwt.Key
	// Expire defines the duration when the token is about to expire
	Expire time.Duration
	// SigningMethod how to sign the JWT. For default value see the OptionFuncs
	SigningMethod csjwt.Signer
	// Verifier token parser and verifier bound to ONE signing method. Setting
	// a new SigningMethod also overwrites the JWTVerify pointer.
	// TODO(newbies): For Verification add Options for setting custom Unmarshaler, HTTP FORM input name and cookie name.
	Verifier *csjwt.Verification
	// EnableJTI activates the (JWT ID) Claim, a unique identifier. UUID.
	EnableJTI bool
	// ErrorHandler specific for this scope. if nil, the the next handler in
	// the chain will be called.
	ErrorHandler http.Handler
	// KeyFunc will receive the parsed token and should return the key for validating.
	KeyFunc csjwt.Keyfunc
	// templateTokenFunc to a create a new template token when parsing
	// a byte token slice into the template token.
	// Default value nil.
	templateTokenFunc func() csjwt.Token
}

// TODO(cs) maybe we can replace csjwt.Token with our own interface definition but seems complex.

// IsValid a configuration for a scope is only then valid when the Key has been
// supplied, a non-nil signing method and a non-nil Verifier.
func (sc *scopedConfig) IsValid() bool {
	return !sc.Key.IsEmpty() && sc.SigningMethod != nil && sc.Verifier != nil
}

// TemplateToken returns the template token. Default Claim is a map. You can
// provide your own by setting the template token function. WithTemplateToken()
func (sc scopedConfig) TemplateToken() csjwt.Token {
	if sc.templateTokenFunc != nil {
		return sc.templateTokenFunc()
	}
	// must be a pointer because of the unmarshalling function
	// default claim defines a map[string]interface{}
	// TODO(cs) get rid of dependency on jwtclaim.Map
	return csjwt.NewToken(&jwtclaim.Map{})
}

// ParseFromRequest parses a request to find a token in either the header, a
// cookie or an HTML form.
func (sc scopedConfig) ParseFromRequest(r *http.Request) (csjwt.Token, error) {
	return sc.Verifier.ParseFromRequest(sc.TemplateToken(), sc.KeyFunc, r)
}

// Parse parses a raw token.
func (sc scopedConfig) Parse(rawToken []byte) (csjwt.Token, error) {
	return sc.Verifier.Parse(sc.TemplateToken(), rawToken, sc.KeyFunc)
}

// initKeyFunc generates a closure for a specific scope to compare if the
// algorithm in the token matches with the current algorithm.
func (sc *scopedConfig) initKeyFunc() {
	sc.KeyFunc = func(t csjwt.Token) (csjwt.Key, error) {

		if have, want := t.Alg(), sc.SigningMethod.Alg(); have != want {
			return csjwt.Key{}, errors.NewNotImplementedf(errUnknownSigningMethod, have, want)
		}
		if sc.Key.Error != nil {
			return csjwt.Key{}, errors.Wrap(sc.Key.Error, "[jwtauth] Key Error")
		}
		return sc.Key, nil
	}
}

func defaultScopedConfig() (scopedConfig, error) {
	key := csjwt.WithPasswordRandom()
	hs256, err := csjwt.NewHMACFast256(key)
	sc := scopedConfig{
		ScopeHash:     scope.DefaultHash,
		Expire:        DefaultExpire,
		Key:           key,
		SigningMethod: hs256,
		Verifier:      csjwt.NewVerification(hs256),
		EnableJTI:     false,
	}
	sc.initKeyFunc()
	return sc, err
}

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
	return func(s *Service) {
		if s.optionError != nil {
			return
		}

		if h == scope.DefaultHash {
			s.defaultScopeCache, s.optionError = defaultScopedConfig()
			s.optionError = errors.Wrap(s.optionError, "[jwtauth] Default Scope with Default Config")
			return
		}

		s.mu.Lock()
		s.scopeCache[h], s.optionError = defaultScopedConfig()
		s.optionError = errors.Wrapf(s.optionError, "[jwtauth] Scope %s with Default Config", h)
		s.mu.Unlock()
	}
}

// WithBlacklist sets a new global black list service.
// Convenience helper function.
func WithBlacklist(bl Blacklister) Option {
	return func(s *Service) {
		s.Blacklist = bl
	}
}

// WithLogger sets a new global logger.
// Convenience helper function.
func WithLogger(l log.Logger) Option {
	return func(s *Service) {
		s.Log = l
	}
}

// WithBackend applies the backend configuration to the service.
// Once this has been set all other option functions are not really
// needed.
//	cfgStruct, err := jwtauthbe.NewConfigStructure()
//	if err != nil {
//		panic(err)
//	}
//	pb := jwtauthbe.NewBackend(cfgStruct, cfgmodel.WithEncryptor(myEncryptor{}))
//
//	jwts := jwtauth.MustNewService(
//		jwtauth.WithBackend(jwtauthbe.BackendOptions(pb)),
//	)
func WithBackend(f ScopedOptionFunc) Option {
	return func(s *Service) {
		s.scpOptionFnc = f
	}
}

// WithTemplateToken set a custom csjwt.Header and csjwt.Claimer for each scope
// when parsing a token in a request. Function f will generate a new base token
// for each request. This allows you to choose using a slow map as a claim
// or a fast struct based claim. Same goes with the header.
func WithTemplateToken(scp scope.Scope, id int64, f func() csjwt.Token) Option {
	h := scope.NewHash(scp, id)
	return func(s *Service) {

		if h == scope.DefaultHash {
			s.defaultScopeCache.templateTokenFunc = f
			return
		}

		s.mu.Lock()
		defer s.mu.Unlock()

		// inherit default config
		scNew := s.defaultScopeCache
		scNew.templateTokenFunc = f

		if sc, ok := s.scopeCache[h]; ok {
			sc.templateTokenFunc = scNew.templateTokenFunc
			scNew = sc
		}
		scNew.ScopeHash = h
		s.scopeCache[h] = scNew
	}
}

// WithSigningMethod this option function lets you overwrite the default 256 bit
// signing method for a specific scope. Used incorrectly token decryption can fail.
func WithSigningMethod(scp scope.Scope, id int64, sm csjwt.Signer) Option {
	h := scope.NewHash(scp, id)
	return func(s *Service) {

		if h == scope.DefaultHash {
			s.defaultScopeCache.SigningMethod = sm
			s.defaultScopeCache.Verifier = csjwt.NewVerification(sm)
			return
		}

		s.mu.Lock()
		defer s.mu.Unlock()

		// inherit default config
		scNew := s.defaultScopeCache

		scNew.SigningMethod = sm
		scNew.Verifier = csjwt.NewVerification(sm)

		if sc, ok := s.scopeCache[h]; ok {
			sc.SigningMethod = scNew.SigningMethod
			sc.Verifier = scNew.Verifier
			scNew = sc
		}

		scNew.ScopeHash = h
		s.scopeCache[h] = scNew
	}
}

// WithErrorHandler sets the error handler for a scope and its ID. If the
// scope.DefaultID will be set the handler gets also applied to the global
// handler
func WithErrorHandler(scp scope.Scope, id int64, handler http.Handler) Option {
	h := scope.NewHash(scp, id)
	return func(s *Service) {

		if h == scope.DefaultHash {
			s.defaultScopeCache.ErrorHandler = handler
			return
		}

		s.mu.Lock()
		defer s.mu.Unlock()

		// inherit default config
		scNew := s.defaultScopeCache
		scNew.ErrorHandler = handler

		if sc, ok := s.scopeCache[h]; ok {
			sc.ErrorHandler = scNew.ErrorHandler
			scNew = sc
		}
		scNew.ScopeHash = h
		s.scopeCache[h] = scNew
	}
}

// WithExpiration sets expiration duration depending on the scope
func WithExpiration(scp scope.Scope, id int64, d time.Duration) Option {
	h := scope.NewHash(scp, id)
	return func(s *Service) {

		if h == scope.DefaultHash {
			s.defaultScopeCache.Expire = d
			return
		}

		s.mu.Lock()
		defer s.mu.Unlock()

		// inherit default config
		scNew := s.defaultScopeCache
		scNew.Expire = d

		if sc, ok := s.scopeCache[h]; ok {
			sc.Expire = scNew.Expire
			scNew = sc
		}
		scNew.ScopeHash = h
		s.scopeCache[h] = scNew
	}
}

// WithTokenID enables JTI (JSON Web Token ID) for a specific scope
func WithTokenID(scp scope.Scope, id int64, enable bool) Option {
	h := scope.NewHash(scp, id)
	return func(s *Service) {

		if h == scope.DefaultHash {
			s.defaultScopeCache.EnableJTI = enable
			return
		}

		s.mu.Lock()
		defer s.mu.Unlock()

		// inherit default config
		scNew := s.defaultScopeCache
		scNew.EnableJTI = enable

		if sc, ok := s.scopeCache[h]; ok {
			sc.EnableJTI = scNew.EnableJTI
			scNew = sc
		}
		scNew.ScopeHash = h
		s.scopeCache[h] = scNew
	}
}

// WithKey sets the key for the default signing method of 256 bits.
// You can also provide your own signing method by using additionally
// the function WithSigningMethod(), which must be called after this function :-/.
func WithKey(scp scope.Scope, id int64, key csjwt.Key) Option {
	h := scope.NewHash(scp, id)
	if key.Error != nil {
		return func(s *Service) {
			s.optionError = errors.Wrap(key.Error, "[jwtauth] Key Error")
		}
	}
	if key.IsEmpty() {
		return func(s *Service) {
			s.optionError = errors.NewEmptyf(errKeyEmpty)
		}
	}
	return func(s *Service) {
		if s.optionError != nil {
			return
		}

		s.mu.Lock()
		defer s.mu.Unlock()

		// inherit default config
		scNew := s.defaultScopeCache
		scNew.Key = key

		// if you are not satisfied with the bit size of 256 you can change it
		// by using WithSigningMethod
		switch key.Algorithm() {
		case csjwt.ES:
			scNew.SigningMethod = csjwt.NewSigningMethodES256()
		case csjwt.HS:
			scNew.SigningMethod, s.optionError = csjwt.NewHMACFast256(key)
			if s.optionError != nil {
				s.optionError = errors.Wrap(s.optionError, "[jwtauth] HMAC Fast 256 error")
				return
			}
		case csjwt.RS:
			scNew.SigningMethod = csjwt.NewSigningMethodRS256()
		default:
			s.optionError = errors.NewNotImplementedf(errUnknownSigningMethodOptions, key.Algorithm())
			return
		}

		scNew.Verifier = csjwt.NewVerification(scNew.SigningMethod)
		scNew.initKeyFunc()

		if h == scope.DefaultHash {
			s.defaultScopeCache.Key = scNew.Key
			s.defaultScopeCache.SigningMethod = scNew.SigningMethod
			s.defaultScopeCache.Verifier = scNew.Verifier
			s.defaultScopeCache.KeyFunc = scNew.KeyFunc
			return
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
	}
}
