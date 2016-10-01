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

package geoip

import "sync/atomic"

//go:generate go run ../internal/scopedservice/main_copy.go "$GOPACKAGE"

// Service represents a service manager for GeoIP detection and restriction.
// Please consider the law in your country if you would like to implement
// geo-blocking.
type Service struct {
	service

	// Finder finds a country by an IP address. If nil panics during execution
	// in the middleware. This field gets protected by a mutex to allowing
	// setting the field during requests.
	Finder

	// geoIPLoaded checks to only load the GeoIP CountryRetriever once because
	// we may set that within a request. It's defined in the backend
	// configuration but later we need to reset this value to zero to allow
	// reloading.
	geoIPLoaded uint32
}

// New creates a new GeoIP service to be used as a middleware or standalone.
func New(opts ...Option) (*Service, error) {
	s, err := newService(opts...)
	if err != nil {
		return nil, err
	}

	return s, nil
}

// Close closes the underlying GeoIP CountryRetriever service and resets the
// internal loading state of the GeoIP flag. It does not yet clear the internal
// cache.
func (s *Service) Close() error {
	atomic.StoreUint32(&s.geoIPLoaded, 0)
	return s.Finder.Close()
}

// isGeoIPLoaded checks if the geoip lookup interface has been set by an object.
// this can be adjusted dynamically with the scoped configuration.
func (s *Service) isGeoIPLoaded() bool {
	return atomic.LoadUint32(&s.geoIPLoaded) == 1
}
