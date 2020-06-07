// Copyright 2015-present, Cyrill @ Schumacher.fm and the CoreStore contributors
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

	"github.com/corestoreio/errors"
	"github.com/corestoreio/pkg/net/mw"
	"github.com/corestoreio/pkg/store/scope"
	"github.com/corestoreio/pkg/util/csjwt"
	"github.com/corestoreio/pkg/util/csjwt/jwtclaim"
)

// ScopedConfig contains the configuration for a scope
type ScopedConfig struct {
	scopedConfigGeneric
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
	// SingleTokenUsage if set to true for each request a token can be only used
	// once. The JTI (JSON Token Identifier) gets added to the blockList until it
	// expires.
	SingleTokenUsage bool
}

var defaultUnauthorizedHandler = mw.ErrorWithStatusCode(http.StatusUnauthorized)

// IsValid check if the scoped configuration is valid when:
//		- Key
//		- SigningMethod
//		- Verifier
// has been set and no other previous error has occurred.
func (sc *ScopedConfig) isValid() error {
	if err := sc.isValidPreCheck(); err != nil {
		return errors.Wrap(err, "[jwt] ScopedConfig.isValid as an lastErr")
	}
	if sc.Disabled {
		return nil
	}
	if sc.Key.IsEmpty() || sc.SigningMethod == nil || sc.Verifier == nil {
		return errors.NewNotValidf(errScopedConfigNotValid, sc.ScopeID)
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
func (sc ScopedConfig) ParseFromRequest(bl Blocklister, r *http.Request) (csjwt.Token, error) {
	dst := sc.TemplateToken()

	if err := sc.Verifier.ParseFromRequest(&dst, sc.KeyFunc, r); err != nil {
		return dst, errors.Wrap(err, "[jwt] ScopedConfig.Verifier.ParseFromRequest")
	}

	kid, err := extractJTI(dst)
	if err != nil {
		return dst, errors.Wrap(err, "[jwt] ScopedConfig.ParseFromRequest.extractJTI")
	}

	if bl.Has(kid) {
		return dst, errors.NewNotValidf(errTokenBlocklisted)
	}
	if sc.SingleTokenUsage {
		if err := bl.Set(kid, dst.Claims.Expires()); err != nil {
			return dst, errors.Wrap(err, "[jwt] ScopedConfig.ParseFromRequest.Blocklist.Set")
		}
	}
	return dst, nil
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
	// copy the data from sc pointer to avoid race conditions under high load
	// test in package backendjwt: $ go test -race -run=TestServiceWithBackend_WithRunMode_Valid_Request -count=8 .
	var alg string
	if sc.SigningMethod != nil {
		alg = sc.SigningMethod.Alg()
	}
	key := sc.Key
	keyErr := sc.Key.Error
	sc.KeyFunc = func(t *csjwt.Token) (csjwt.Key, error) {
		if have, want := t.Alg(), alg; have != want {
			return csjwt.Key{}, errors.NewNotImplementedf(errUnknownSigningMethod, have, want)
		}
		if keyErr != nil {
			return csjwt.Key{}, errors.Wrap(sc.Key.Error, "[jwt] ScopedConfig.initKeyFunc.Key.Error")
		}
		return key, nil
	}
}

func newScopedConfig(target, parent scope.TypeID) *ScopedConfig {
	key := csjwt.WithPasswordRandom()
	hs256, err := csjwt.NewSigningMethodHS256Fast(key)
	if err != nil {
		return &ScopedConfig{
			scopedConfigGeneric: scopedConfigGeneric{
				ScopeID:  target,
				ParentID: parent,
				lastErr:  errors.Wrap(err, "[jwt] defaultScopedConfig.NewHMACFast256"),
			},
		}
	}
	sc := &ScopedConfig{
		scopedConfigGeneric: newScopedConfigGeneric(target, parent),
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
