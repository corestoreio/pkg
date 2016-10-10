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
func WithDefaultConfig(h scope.TypeID) Option {
	return withDefaultConfig(h)
}

// WithHash sets the hashing algorithm to create a new hash and verify an
// incoming hash. Please use only cryptographically secure hash algorithms.
func WithHash(name string, key []byte, scopeIDs ...scope.TypeID) Option {
	return func(s *Service) error {
		sc := s.findScopedConfig(scopeIDs...)
		sc.hashPoolInit(name, key)
		return s.updateScopedConfig(sc)
	}
}

// WithHeaderHandler sets the writer and the parser. The writer knows how to
// write the hash value into the HTTP header. The parser knows how and where to
// extract the hash value from the header or even the trailer. Compatible types
// in this package are ContentHMAC, ContentSignature and Transparent.
func WithHeaderHandler(pw HeaderParseWriter, scopeIDs ...scope.TypeID) Option {
	return func(s *Service) error {
		sc := s.findScopedConfig(scopeIDs...)
		sc.HeaderParseWriter = pw
		return s.updateScopedConfig(sc)
	}
}

// WithAllowedMethods sets the allowed HTTP methods which can transport a
// signature hash.
func WithAllowedMethods(methods []string, scopeIDs ...scope.TypeID) Option {
	return func(s *Service) error {
		sc := s.findScopedConfig(scopeIDs...)
		sc.AllowedMethods = methods
		return s.updateScopedConfig(sc)
	}
}

// WithTrailer allows to write the hash sum into the trailer. The middleware switches
// to stream based hash calculation which results in faster processing instead of writing
// into a buffer. Make sure that your client can process HTTP trailers.
func WithTrailer(inTrailer bool, scopeIDs ...scope.TypeID) Option {
	return func(s *Service) error {
		sc := s.findScopedConfig(scopeIDs...)
		sc.InTrailer = inTrailer
		return s.updateScopedConfig(sc)
	}
}

// WithTransparent allows to write the hashes into the Cacher with a
// time-to-live. Responses will not get a header key attached and requests won't
// get inspected for a header key which might contain the hash value.
func WithTransparent(c Cacher, ttl time.Duration, scopeIDs ...scope.TypeID) Option {
	return func(s *Service) error {
		sc := s.findScopedConfig(scopeIDs...)
		sc.TransparentCacher = c
		sc.HeaderParseWriter = MakeTransparent(c, ttl)
		sc.TransparentTTL = ttl
		sc.InTrailer = true // enable streaming hash calculation
		return s.updateScopedConfig(sc)
	}
}
