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
	"github.com/corestoreio/pkg/store/scope"
	"github.com/corestoreio/pkg/util/csjwt"
	"github.com/corestoreio/errors"
	"github.com/corestoreio/log"
)

//go:generate go run ../internal/scopedservice/main_copy.go "$GOPACKAGE"

const (
	claimExpiresAt = "exp"
	claimIssuedAt  = "iat"
	claimKeyID     = "jti"
)

// Service main type for handling JWT authentication, generation, blacklists and
// log outs depending on a scope.
type Service struct {
	service
	// JTI generates for each token a new ID. It panics on nil. Default ID
	// generator runs in package util/shortid.
	JTI IDGenerator
	// Blacklist concurrent safe black list service which handles blocked
	// tokens. Default black hole storage. Must be thread safe.
	Blacklist Blacklister
}

// New creates a new token service.
// Default values from option function WithDefaultConfig() will be
// applied if passing no functional option arguments.
func New(opts ...Option) (*Service, error) {
	s, err := newService(opts...)
	if err != nil {
		return nil, err
	}
	s.useWebsite = true
	s.optionAfterApply = func() error {
		s.rwmu.RLock()
		defer s.rwmu.RUnlock()
		for h := range s.scopeCache {
			// This one checks if the configuration contains only the default or
			// website scope. Store scope is neither allowed nor supported.
			if scp, _ := h.Unpack(); scp > scope.Website {
				return errors.NewNotSupportedf(errServiceUnsupportedScope, h)
			}
		}
		return nil
	}
	if s.JTI == nil {
		s.JTI = jti{}
	}
	if s.Blacklist == nil {
		s.Blacklist = nullBL{}
	}
	if err := s.optionAfterApply(); err != nil {
		return nil, err
	}
	return s, nil
}

// NewToken creates a new signed JSON web token based on the predefined scoped
// based template token function (WithTemplateToken) and merges the optional 3rd
// argument into the template token claim. The returned token is owned by the
// caller. The tokens Raw field contains the freshly signed byte slice.
// ExpiresAt, IssuedAt and ID are already set and cannot be overwritten, but you
// can access them. It panics if the provided template token has a nil Header or
// Claimer field.
func (s *Service) NewToken(scopeID scope.TypeID, claim ...csjwt.Claimer) (csjwt.Token, error) {
	var empty csjwt.Token
	now := csjwt.TimeFunc()

	sc, err := s.ConfigByScopeID(scopeID, 0)
	if err != nil {
		return empty, errors.Wrap(err, "[jwt] NewToken.ConfigByScopeID")
	}

	var tk = sc.TemplateToken()

	if len(claim) > 0 && claim[0] != nil {
		if err := csjwt.MergeClaims(tk.Claims, claim...); err != nil {
			return empty, errors.Wrap(err, "[jwt] NewToken.MergeClaims")
		}
	}

	if err := tk.Claims.Set(claimExpiresAt, now.Add(sc.Expire).Unix()); err != nil {
		return empty, errors.Wrap(err, "[jwt] NewToken.Claims.Set EXP")
	}
	if err := tk.Claims.Set(claimIssuedAt, now.Unix()); err != nil {
		return empty, errors.Wrap(err, "[jwt] NewToken.Claims.Set IAT")
	}

	jti, err := s.JTI.NewID()
	if err != nil {
		return empty, errors.Wrap(err, "[jwt] NewToken.Claims.Set JTI.NewID")
	}
	if err := tk.Claims.Set(claimKeyID, jti); err != nil {
		return empty, errors.Wrapf(err, "[jwt] NewToken.Claims.Set KID: %q", jti)
	}

	tk.Raw, err = tk.SignedString(sc.SigningMethod, sc.Key)
	return tk, errors.Wrap(err, "[jwt] NewToken.SignedString")
}

// Logout adds a token securely to a blacklist with the expiration duration. If
// the JTI or token ID is empty or missing, an error gets returned of behaviour
// Empty.
func (s *Service) Logout(token csjwt.Token) error {
	if len(token.Raw) == 0 || !token.Valid {
		return nil
	}

	kid, err := extractJTI(token)
	if err != nil {
		return errors.Wrap(err, "[jwt] Service.Logout extractJTI")
	}

	return errors.Wrap(s.Blacklist.Set(kid, token.Claims.Expires()), "[jwt] Service.Logout.Blacklist.Set")
}

// Parse parses a token string with the DefaultID scope and returns the
// valid token or an error.
func (s *Service) Parse(rawToken []byte) (csjwt.Token, error) {
	return s.ParseScoped(scope.DefaultTypeID, rawToken)
}

// ParseScoped parses a token based on the applied scope and the scope ID.
// Different configurations are passed to the token parsing function. The black
// list will be checked for containing entries.
func (s *Service) ParseScoped(scopeID scope.TypeID, rawToken []byte) (csjwt.Token, error) {
	var empty csjwt.Token

	sc, err := s.ConfigByScopeID(scopeID, 0)
	if err != nil {
		return empty, errors.Wrap(err, "[jwt] ParseScoped.ConfigByScopeID")
	}

	token, err := sc.Parse(rawToken)
	if err != nil {
		return empty, errors.Wrap(err, "[jwt] ParseScoped.Parse")
	}

	var inBL bool
	isValid := token.Valid && len(token.Raw) > 0
	if isValid {
		inBL = s.Blacklist.Has(token.Raw)
	}
	if isValid && !inBL {
		return token, nil
	}
	if s.Log.IsDebug() {
		s.Log.Debug("jwt.Service.ParseScoped", log.Err(err), log.Bool("inBlackList", inBL), log.String("rawToken", string(rawToken)), log.Marshal("token", token))
	}
	return empty, errors.NewNotValidf(errTokenParseNotValidOrBlackListed)
}
