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

//go:generate go run ../internal/scopedservice/main_copy.go "$GOPACKAGE"

package ratelimit

// Service creates a middleware that facilitates using a Limiter to limit HTTP
// requests.
type Service struct {
	service
}

// New creates a new rate limit middleware.
//
// Default DeniedHandler returns http.StatusTooManyRequests.
//
// Default RateLimiterFactory is the NewGCRAMemStore(). If *PkgBackend has
// been provided the values from the configration will be taken otherwise
// GCRAMemStore() uses the Default* variables.
func New(opts ...Option) (*Service, error) {
	return newService(opts...)
}

// FlushCache clears the internal cache
func (s *Service) FlushCache() error {
	return s.flushCache()
}
