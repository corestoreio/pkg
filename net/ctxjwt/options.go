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
	"net/http"
	"time"

	"github.com/corestoreio/csfw/config"
	"github.com/corestoreio/csfw/net/ctxhttp"
	"github.com/corestoreio/csfw/store/scope"
	"github.com/corestoreio/csfw/util/csjwt"
	"github.com/corestoreio/csfw/util/csjwt/jwtclaim"
	"github.com/juju/errors"
)

// scopedConfig private internal scoped based configuration
type scopedConfig struct {
	scopeHash scope.Hash
	// Key contains the HMAC, RSA or ECDSA sensitive data. Yes the csjwt.Key
	// must not be embedded into this struct because otherwise when printing
	// or logging the sensitive data from csjwt.Key gets leaked into loggers
	// or where ever. If key would be lower case then %#v still prints
	// every field of the csjwt.Key.
	Key csjwt.Key
	// expire defines the duration when the token is about to expire
	expire time.Duration
	// signingMethod how to sign the JWT. For default value see the OptionFuncs
	signingMethod csjwt.Signer
	// jwtVerify token parser and verifier bound to ONE signing method. Setting
	// a new SigningMethod also overwrites the JWTVerify pointer.
	// TODO: add Option for setting custom Unmarshaler and HTTP FORM input name
	jwtVerify *csjwt.Verification
	// enableJTI activates the (JWT ID) Claim, a unique identifier. UUID.
	enableJTI bool
	// errorHandler specific for this scope. if nil, fallback to global one
	// stored in the Service
	errorHandler ctxhttp.Handler
	// keyFunc will receive the parsed token and should return the key for validating.
	keyFunc csjwt.Keyfunc
	// newClaims function to a create a new empty claim when parsing a token in a request.
	// Default value nil. Returned Claimer must be a pointer.
	newClaims func() csjwt.Claimer
}

// isValid a configuration for a scope is only then valid when the Key has been
// supplied and a applied non-nil signing method.
func (sc *scopedConfig) isValid() bool {
	return !sc.Key.IsEmpty() && sc.signingMethod != nil && sc.jwtVerify != nil
}

func (sc scopedConfig) parseFromRequest(r *http.Request) (csjwt.Token, error) {
	var claim csjwt.Claimer
	if sc.newClaims != nil {
		claim = sc.newClaims()
	} else {
		claim = &jwtclaim.Map{}
	}
	return sc.jwtVerify.ParseFromRequest(r, sc.keyFunc, claim)
}

func (sc scopedConfig) parseWithClaim(rawToken []byte) (csjwt.Token, error) {
	var claim csjwt.Claimer
	if sc.newClaims != nil {
		claim = sc.newClaims()
	} else {
		claim = &jwtclaim.Map{}
	}
	return sc.jwtVerify.ParseWithClaim(rawToken, sc.keyFunc, claim)
}

// getKeyFunc generates the key function for a specific scope and to used in caching
func getKeyFunc(scpCfg scopedConfig) csjwt.Keyfunc {
	return func(t csjwt.Token) (csjwt.Key, error) {

		if have, want := t.Alg(), scpCfg.signingMethod.Alg(); have != want {
			return csjwt.Key{}, errors.Errorf("[ctxjwt] Unknown signing method - Have: %q Want: %q", have, want)
		}
		if scpCfg.Key.Error != nil {
			return csjwt.Key{}, errors.Mask(scpCfg.Key.Error)
		}
		return scpCfg.Key, nil
	}
}

// Option can be used as an argument in NewService to configure a token service.
type Option func(*Service)

func defaultScopedConfig() (scopedConfig, error) {
	key := csjwt.WithPasswordRandom()
	hs256, err := csjwt.NewHMACFast256(key)
	return scopedConfig{
		scopeHash:     scope.DefaultHash,
		expire:        DefaultExpire,
		Key:           key,
		signingMethod: hs256,
		jwtVerify:     csjwt.NewVerification(hs256),
		enableJTI:     false,
	}, err
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
		s.mu.Lock()
		defer s.mu.Unlock()
		var err error
		s.scopeCache[h], err = defaultScopedConfig()
		s.MultiErr = s.AppendErrors(err)
	}
}

// WithBlacklist sets a new global black list service.
// Convenience helper function.
func WithBlacklist(blacklist Blacklister) Option {
	return func(s *Service) {
		s.Blacklist = blacklist
	}
}

// WithBackend applies the backend configuration to the service.
// Once this has been set all other option functions are not really
// needed.
// Convenience helper function.
func WithBackend(pb *PkgBackend) Option {
	return func(s *Service) {
		s.Backend = pb
	}
}

// WithNewClaims set a custom Claimer for each scope when parsing a token
// in a request. Function f will generate a new Claim for each request.
// This allows you to choose for using a map based claim or a struct based.
// The returned Claimer interface by predicate f must be a pointer.
func WithNewClaims(scp scope.Scope, id int64, f func() csjwt.Claimer) Option {
	h := scope.NewHash(scp, id)
	return func(s *Service) {
		s.mu.Lock()
		defer s.mu.Unlock()

		scNew := scopedConfig{
			newClaims: f,
		}
		if sc, ok := s.scopeCache[h]; ok {
			sc.newClaims = scNew.newClaims
			scNew = sc
		}
		scNew.scopeHash = h
		s.scopeCache[h] = scNew
	}
}

// WithSigningMethod this option function lets you overwrite the default 256 bit
// signing method for a specific scope. Used incorrectly token decryption can fail.
func WithSigningMethod(scp scope.Scope, id int64, sm csjwt.Signer) Option {
	h := scope.NewHash(scp, id)
	return func(s *Service) {
		s.mu.Lock()
		defer s.mu.Unlock()

		scNew := scopedConfig{
			signingMethod: sm,
			jwtVerify:     csjwt.NewVerification(sm),
		}
		if sc, ok := s.scopeCache[h]; ok {
			sc.signingMethod = scNew.signingMethod
			sc.jwtVerify = scNew.jwtVerify
			scNew = sc
		}
		scNew.scopeHash = h
		s.scopeCache[h] = scNew
	}
}

// WithErrorHandler sets the error handler for a scope and its ID. If the
// scope.DefaultID will be set the handler gets also applied to the global
// handler
func WithErrorHandler(scp scope.Scope, id int64, handler ctxhttp.Handler) Option {
	h := scope.NewHash(scp, id)
	return func(s *Service) {
		s.mu.Lock()
		defer s.mu.Unlock()
		scNew := scopedConfig{
			errorHandler: handler,
		}
		if sc, ok := s.scopeCache[h]; ok {
			sc.errorHandler = scNew.errorHandler
			scNew = sc
		}
		scNew.scopeHash = h
		s.scopeCache[h] = scNew
		if scp == scope.Default && id == 0 {
			s.DefaultErrorHandler = handler
		}
	}
}

// WithExpiration sets expiration duration depending on the scope
func WithExpiration(scp scope.Scope, id int64, d time.Duration) Option {
	h := scope.NewHash(scp, id)
	return func(s *Service) {
		s.mu.Lock()
		defer s.mu.Unlock()

		scNew := scopedConfig{
			expire: d,
		}

		if sc, ok := s.scopeCache[h]; ok {
			sc.expire = scNew.expire
			scNew = sc
		}
		scNew.scopeHash = h
		s.scopeCache[h] = scNew
	}
}

// WithTokenID enables JTI (JSON Web Token ID) for a specific scope
func WithTokenID(scp scope.Scope, id int64, enable bool) Option {
	h := scope.NewHash(scp, id)
	return func(s *Service) {
		s.mu.Lock()
		defer s.mu.Unlock()

		scNew := scopedConfig{
			enableJTI: enable,
		}

		if sc, ok := s.scopeCache[h]; ok {
			sc.enableJTI = scNew.enableJTI
			scNew = sc
		}
		scNew.scopeHash = h
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
			s.MultiErr = s.AppendErrors(key.Error)
		}
	}
	if key.IsEmpty() {
		return func(s *Service) {
			s.MultiErr = s.AppendErrors(errors.New("[ctxjwt] Provided key argument is empty"))
		}
	}
	return func(s *Service) {
		s.mu.Lock()
		defer s.mu.Unlock()

		scNew := scopedConfig{
			Key: key,
		}

		// if you are not satisfied with the bit size of 256 you can change it
		// by using WithSigningMethod
		switch key.Algorithm() {
		case csjwt.ES:
			scNew.signingMethod = csjwt.NewSigningMethodES256()
		case csjwt.HS:
			var err error
			scNew.signingMethod, err = csjwt.NewHMACFast256(key)
			if err != nil {
				s.MultiErr = s.AppendErrors(errors.Mask(err))
				return
			}
		case csjwt.RS:
			scNew.signingMethod = csjwt.NewSigningMethodRS256()
		default:
			s.MultiErr = s.AppendErrors(errors.Errorf("[ctxjwt] Unknown signing method - Have: %q Want: ES, HS or RS", key.Algorithm()))
			return
		}

		if sc, ok := s.scopeCache[h]; ok {
			sc.Key = scNew.Key
			sc.signingMethod = scNew.signingMethod
			sc.jwtVerify = csjwt.NewVerification(scNew.signingMethod)
			scNew = sc
		}
		scNew.scopeHash = h
		s.scopeCache[h] = scNew
	}
}

// optionsByBackend creates an option array containing the Options based
// on the configuration
func optionsByBackend(be *PkgBackend, sg config.ScopedGetter) (opts []Option) {
	scp, id := sg.Scope()

	exp, err := be.NetCtxjwtExpiration.Get(sg)
	if err != nil {
		return append(opts, func(s *Service) {
			s.MultiErr = s.AppendErrors(errors.Mask(err))
		})
	}
	opts = append(opts, WithExpiration(scp, id, exp))

	isJTI, err := be.NetCtxjwtEnableJTI.Get(sg)
	if err != nil {
		return append(opts, func(s *Service) {
			s.MultiErr = s.AppendErrors(errors.Mask(err))
		})
	}
	opts = append(opts, WithTokenID(scp, id, isJTI))

	signingMethod, err := be.NetCtxjwtSigningMethod.Get(sg)
	if err != nil {
		return append(opts, func(s *Service) {
			s.MultiErr = s.AppendErrors(errors.Mask(err))
		})
	}

	var key csjwt.Key

	switch signingMethod.Alg() {
	case csjwt.RS256, csjwt.RS384, csjwt.RS512:

		rsaKey, err := be.NetCtxjwtRSAKey.Get(sg)
		if err != nil {
			return append(opts, func(s *Service) {
				s.MultiErr = s.AppendErrors(errors.Mask(err))
			})
		}
		rsaPW, err := be.NetCtxjwtRSAKeyPassword.Get(sg)
		if err != nil {
			return append(opts, func(s *Service) {
				s.MultiErr = s.AppendErrors(errors.Mask(err))
			})
		}
		key = csjwt.WithRSAPrivateKeyFromPEM(rsaKey, rsaPW)

	case csjwt.ES256, csjwt.ES384, csjwt.ES512:

		ecdsaKey, err := be.NetCtxjwtECDSAKey.Get(sg)
		if err != nil {
			return append(opts, func(s *Service) {
				s.MultiErr = s.AppendErrors(errors.Mask(err))
			})
		}
		ecdsaPW, err := be.NetCtxjwtECDSAKeyPassword.Get(sg)
		if err != nil {
			return append(opts, func(s *Service) {
				s.MultiErr = s.AppendErrors(errors.Mask(err))
			})
		}
		key = csjwt.WithECPrivateKeyFromPEM(ecdsaKey, ecdsaPW)

	case csjwt.HS256, csjwt.HS384, csjwt.HS512:

		password, err := be.NetCtxjwtHmacPassword.Get(sg)
		if err != nil {
			return append(opts, func(s *Service) {
				s.MultiErr = s.AppendErrors(errors.Mask(err))
			})
		}
		key = csjwt.WithPassword(password)

	default:
		opts = append(opts, func(s *Service) {
			s.MultiErr = s.AppendErrors(errors.Errorf("[ctxjwt] Unknown signing method: %q", signingMethod.Alg()))
		})
	}

	// WithSigningMethod must be added at the end of the slice to overwrite default signing methods
	return append(opts, WithKey(scp, id, key), WithSigningMethod(scp, id, signingMethod))
}
