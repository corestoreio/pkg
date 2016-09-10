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
	"crypto/hmac"
	"crypto/sha256"
	"hash"

	"github.com/corestoreio/csfw/store/scope"
	"github.com/corestoreio/csfw/util/errors"
	"github.com/corestoreio/csfw/util/hashpool"
	"github.com/minio/blake2b-simd"
)

// WithDefaultConfig applies the default signed configuration settings based
// for a specific scope. This function overwrites any previous set options.
//
// Default values are:
//		- Denied Handler: http.Error(w, http.StatusText(http.StatusTooManyRequests), http.StatusTooManyRequests)
//		- VaryByer: returns an empty key
// Example:
//		s := MustNewService(WithDefaultConfig(scope.Store,1), WithVaryBy(scope.Store, 1, myVB))
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
		sc.hashPool = hashpool.New(func() hash.Hash {
			return hmac.New(hh, key)
		})
		sc.ScopeHash = h
		s.scopeCache[h] = sc
		return nil
	}
}

// WithHeaderHandler sets the writer and the parser. The writer knows how to
// write the hash value into the HTTP header. The parser knows how and where to
// extract the hash value from the header or even the trailer.
func WithHeaderHandler(scp scope.Scope, id int64, w HTTPWriter, p HTTPParser) Option {
	h := scope.NewHash(scp, id)
	return func(s *Service) error {
		s.rwmu.Lock()
		defer s.rwmu.Unlock()

		sc := s.scopeCache[h]
		if sc == nil {
			sc = optionInheritDefault(s)
		}
		sc.HTTPWriter = w
		sc.HTTPParser = p
		sc.ScopeHash = h
		s.scopeCache[h] = sc
		return nil
	}
}

// WithContentHMAC_SHA256 applies the SHA256 hash with your symmetric key.
func WithContentHMAC_SHA256(scp scope.Scope, id int64, key []byte) Option {
	return func(s *Service) error {
		if err := WithHash(scp, id, sha256.New, key)(s); err != nil {
			return errors.Wrap(err, "[signed] WithContentHMAC_SHA256.WithHash")
		}
		sig := NewHMAC("sha256")
		return WithHeaderHandler(scp, id, sig, sig)(s)
	}
}

// WithContentHMAC_Blake2b256 applies the very fast Blake2 hashing algorithm.
// The current package has been optimized with ASM with for x64 systems, hence
// Blake2 is faster than SHA.
func WithContentHMAC_Blake2b256(scp scope.Scope, id int64, key []byte) Option {
	return func(s *Service) error {
		if err := WithHash(scp, id, blake2b.New256, key)(s); err != nil {
			return errors.Wrap(err, "[signed] WithContentHMAC_Blake2b256.WithHash")
		}
		sig := NewHMAC("blk2b256")
		return WithHeaderHandler(scp, id, sig, sig)(s)
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
