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

	"github.com/corestoreio/csfw/net/mw"
	"github.com/corestoreio/csfw/util/csjwt"
	"github.com/corestoreio/csfw/util/csjwt/jwtclaim"
	"github.com/corestoreio/csfw/util/errors"
)

// ScopedConfig contains the configuration for a scope
type ScopedConfig struct {
	scopedConfigGeneric

	// Disabled if true disables JWT completely
	Disabled bool
	// Key contains the HMAC, RSA or ECDSA sensitive data. The csjwt.Key must
	// not be embedded into this struct because otherwise when printing or
	// logging the sensitive data from csjwt.Key gets leaked into loggers or
	// where ever. If key would be lower case then %#v still prints every field
	// of the csjwt.Key.
	Key csjwt.Key
	// Expire defines the duration when the token is about to expire
	Expire time.Duration
	// Skew duration of time skew we allow between signer and verifier.
	Skew time.Duration
	// SigningMethod how to sign the JWT. For default value see the OptionFuncs
	SigningMethod csjwt.Signer
	// Verifier token parser and verifier bound to ONE signing method. Setting a
	// new SigningMethod also overwrites the JWTVerify pointer. TODO(newbies):
	// For Verification add Options for setting custom Unmarshaler, HTTP FORM
	// input name and cookie name.
	Verifier *csjwt.Verification
	// KeyFunc will receive the parsed token and should return the key for
	// validating.
	KeyFunc csjwt.Keyfunc
	// templateTokenFunc to a create a new template token when parsing a byte
	// token slice into the template token. Default value nil.
	templateTokenFunc func() csjwt.Token

	// UnauthorizedHandler gets called for invalid tokens. Returns the code
	// http.StatusUnauthorized
	UnauthorizedHandler mw.ErrorHandler

	// StoreCodeFieldName optional custom key name used to lookup the claims section
	// to find the store code, defaults to constant store.CodeFieldName.
	StoreCodeFieldName string
}

var defaultUnauthorizedHandler = mw.ErrorWithStatusCode(http.StatusUnauthorized)

// IsValid check if the scoped configuration is valid when:
//		- Key
//		- SigningMethod
//		- Verifier
// has been set and no other previous error has occurred.
func (sc *ScopedConfig) IsValid() error {
	if sc.lastErr != nil {
		return errors.Wrap(sc.lastErr, "[jwt] ScopedConfig.isValid as an lastErr")
	}

	if sc.ScopeHash == 0 || sc.Key.IsEmpty() || sc.SigningMethod == nil || sc.Verifier == nil {
		return errors.NewNotValidf(errScopedConfigNotValid, sc.ScopeHash)
	}
	return nil
}

// TemplateToken returns the template token. Default Claim is a map. You can
// provide your own by setting the template token function. WithTemplateToken()
func (sc ScopedConfig) TemplateToken() (tk csjwt.Token) {
	if sc.templateTokenFunc != nil {
		tk = sc.templateTokenFunc()
	} else {
		// must be a pointer because of the unmarshalling function
		// default claim defines a map[string]interface{}
		tk = csjwt.NewToken(&jwtclaim.Map{})
	}
	_ = tk.Claims.Set(jwtclaim.KeyTimeSkew, sc.Skew)
	return
}

// ParseFromRequest parses a request to find a token in either the header, a
// cookie or an HTML form.
func (sc ScopedConfig) ParseFromRequest(r *http.Request) (csjwt.Token, error) {
	dst := sc.TemplateToken()
	err := sc.Verifier.ParseFromRequest(&dst, sc.KeyFunc, r)
	return dst, errors.Wrap(err, "[jwt] ScopedConfig.Verifier.ParseFromRequest")
}

// Parse parses a raw token.
func (sc ScopedConfig) Parse(rawToken []byte) (csjwt.Token, error) {
	dst := sc.TemplateToken()
	err := sc.Verifier.Parse(&dst, rawToken, sc.KeyFunc)
	return dst, errors.Wrap(err, "[jwt] ScopedConfig.Verifier.Parse")
}

// initKeyFunc generates a closure for a specific scope to compare if the
// algorithm in the token matches with the current algorithm.
func (sc *ScopedConfig) initKeyFunc() {
	sc.KeyFunc = func(t *csjwt.Token) (csjwt.Key, error) {

		if have, want := t.Alg(), sc.SigningMethod.Alg(); have != want {
			return csjwt.Key{}, errors.NewNotImplementedf(errUnknownSigningMethod, have, want)
		}
		if sc.Key.Error != nil {
			return csjwt.Key{}, errors.Wrap(sc.Key.Error, "[jwt] ScopedConfig.initKeyFunc.Key.Error")
		}
		return sc.Key, nil
	}
}

func newScopedConfig() *ScopedConfig {
	key := csjwt.WithPasswordRandom()
	hs256, err := csjwt.NewSigningMethodHS256Fast(key)
	if err != nil {
		se := newScopedConfigError(errors.Wrap(err, "[jwt] defaultScopedConfig.NewHMACFast256"))
		return &se
	}
	sc := &ScopedConfig{
		scopedConfigGeneric: newScopedConfigGeneric(),
		Expire:              DefaultExpire,
		Skew:                DefaultSkew,
		Key:                 key,
		SigningMethod:       hs256,
		Verifier:            csjwt.NewVerification(hs256),
		UnauthorizedHandler: defaultUnauthorizedHandler,
	}
	sc.initKeyFunc()
	return sc
}
