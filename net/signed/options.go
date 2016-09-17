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

package signed

import (
	"hash"
	"time"

	"github.com/corestoreio/csfw/store/scope"
)

// WithDefaultConfig applies the default signed configuration settings based
// for a specific scope. This function overwrites any previous set options.
//
// Default settings: InTrailer activated, Content-HMAC header with sha256,
// allowed HTTP methods set to POST, PUT, PATCH and password for the HMAC SHA
// 256 from a cryptographically random source with a length of 64 bytes.
// Example:
//		s := MustNewService(WithDefaultConfig(scope.Store,1), WithOtherSettings(scope.Store, 1, ...))
func WithDefaultConfig(scp scope.Scope, id int64) Option {
	return withDefaultConfig(scp, id)
}

// WithHash sets the hashing algorithm to create a new hash and verify an
// incoming hash. Please use only cryptographically secure hash algorithms.
func WithHash(scp scope.Scope, id int64, hh func() hash.Hash, key []byte) Option {
	h := scope.NewHash(scp, id)
	return func(s *Service) error {
		s.rwmu.Lock()
		defer s.rwmu.Unlock()

		sc := s.scopeCache[h]
		if sc == nil {
			sc = optionInheritDefault(s)
		}
		sc.hashPoolInit(hh, key)
		sc.ScopeHash = h
		s.scopeCache[h] = sc
		return nil
	}
}

// WithHeaderHandler sets the writer and the parser. The writer knows how to
// write the hash value into the HTTP header. The parser knows how and where to
// extract the hash value from the header or even the trailer.
func WithHeaderHandler(scp scope.Scope, id int64, pw HeaderParseWriter) Option {
	h := scope.NewHash(scp, id)
	return func(s *Service) error {
		s.rwmu.Lock()
		defer s.rwmu.Unlock()

		sc := s.scopeCache[h]
		if sc == nil {
			sc = optionInheritDefault(s)
		}
		sc.HeaderParseWriter = pw
		sc.ScopeHash = h
		s.scopeCache[h] = sc
		return nil
	}
}

// WithDisable allows to disable a signing of the HTTP body or validation.
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

// WithAllowedMethods sets the allowed HTTP methods which can transport a
// signature hash.
func WithAllowedMethods(scp scope.Scope, id int64, methods ...string) Option {
	h := scope.NewHash(scp, id)
	return func(s *Service) error {
		s.rwmu.Lock()
		defer s.rwmu.Unlock()

		sc := s.scopeCache[h]
		if sc == nil {
			sc = optionInheritDefault(s)
		}
		sc.AllowedMethods = methods
		sc.ScopeHash = h
		s.scopeCache[h] = sc
		return nil
	}
}

// WithTrailer allows to write the hash sum into the trailer. The middleware switches
// to stream based hash calculation which results in faster processing instead of writing
// into a buffer. Make sure that your client can process HTTP trailers.
func WithTrailer(scp scope.Scope, id int64, inTrailer bool) Option {
	h := scope.NewHash(scp, id)
	return func(s *Service) error {
		s.rwmu.Lock()
		defer s.rwmu.Unlock()

		sc := s.scopeCache[h]
		if sc == nil {
			sc = optionInheritDefault(s)
		}
		sc.InTrailer = inTrailer
		sc.ScopeHash = h
		s.scopeCache[h] = sc
		return nil
	}
}

// WithTransparent allows to write the hashes into the Cacher with a
// time-to-live. Responses will not get a header key attached and requests won't
// get inspected for a header key which might contain the hash value.
func WithTransparent(scp scope.Scope, id int64, c Cacher, ttl time.Duration) Option {
	h := scope.NewHash(scp, id)
	return func(s *Service) error {
		s.rwmu.Lock()
		defer s.rwmu.Unlock()

		sc := s.scopeCache[h]
		if sc == nil {
			sc = optionInheritDefault(s)
		}
		sc.TransparentCacher = c
		sc.HeaderParseWriter = MakeTransparent(c, ttl)
		sc.TransparentTTL = ttl
		sc.InTrailer = true // enable streaming hash calculation
		sc.ScopeHash = h
		s.scopeCache[h] = sc
		return nil
	}
}
