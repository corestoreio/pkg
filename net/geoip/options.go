// Copyright 2015, Cyrill @ Schumacher.fm and the CoreStore contributors
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

import (
	"fmt"

	"github.com/corestoreio/csfw/config/scope"
	"github.com/corestoreio/csfw/net/ctxhttp"
	"github.com/corestoreio/csfw/utils/os"
	"github.com/oschwald/geoip2-golang"
)

// Option can be used as an argument in NewService to configure a token service.
type Option func(*Service)

// WithAlternativeHandler sets for a scope.Scope and its ID the error handler
// on a Service. If the Handler h is nil falls back to the DefaultErrorHandler.
// This function can be called as many times as you have websites, groups
// or stores.
func WithAlternativeHandler(so scope.Scope, id int64, h ctxhttp.Handler) Option {
	if h == nil {
		h = DefaultAlternativeHandler
	}
	return func(s *Service) {
		switch so {
		case scope.StoreID:
			s.storeIDs.Append(id)
			s.storeAltH = append(s.storeAltH, h)
		case scope.GroupID:
			s.groupIDs.Append(id)
			s.groupAltH = append(s.groupAltH, h)
		case scope.WebsiteID:
			s.websiteIDs.Append(id)
			s.websiteAltH = append(s.websiteAltH, h)
		default:
			s.lastErrors = append(s.lastErrors, scope.ErrUnsupportedScope)
		}
	}
}

// WithCheckAllow sets your custom function which checks is the country of the IP
// address is allowed.
func WithCheckAllow(f IsAllowedFunc) Option {
	return func(s *Service) {
		s.IsAllowed = f
	}
}

// WithGeoIP2Reader creates a new GeoIP2.Reader. As long as there are no other
// readers this is a mandatory argument.
func WithGeoIP2Reader(file string) Option {
	return func(s *Service) {
		if false == os.FileExists(file) {
			s.lastErrors = append(s.lastErrors, fmt.Errorf("File %s not found", file))
			return
		}

		r, err := geoip2.Open(file) // that implementation is not nice for testing because no interface usages :(
		if err != nil {
			s.lastErrors = append(s.lastErrors, err)
			return
		}
		s.GeoIP = r
	}
}
