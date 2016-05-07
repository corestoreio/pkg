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

import (
	"net/http"

	"github.com/corestoreio/csfw/config/cfgmodel"
	"github.com/corestoreio/csfw/store/scope"
	"github.com/corestoreio/csfw/util/errors"
	"github.com/corestoreio/csfw/util/log"
	"github.com/corestoreio/csfw/util/os"
	"github.com/oschwald/geoip2-golang"
)

// Option can be used as an argument in NewService to configure a token service.
type Option func(*Service)

// WithAlternativeHandler sets for a scope.Scope and its ID the error handler
// on a Service. If the Handler h is nil falls back to the DefaultErrorHandler.
// This function can be called as many times as you have websites or stores.
// Group scope is not suppored.
func WithAlternativeHandler(so scope.Scope, id int64, hf http.Handler) Option {
	return func(s *Service) {
		if s.optionError != nil {
			return
		}
		switch so {
		case scope.Store:
			s.storeIDs.Append(id)
			s.storeAltH = append(s.storeAltH, hf)
		case scope.Website:
			s.websiteIDs.Append(id)
			s.websiteAltH = append(s.websiteAltH, hf)
		default:
			s.optionError = errors.NewNotSupportedf(errUnsupportedScope, so, id)
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

func WithAllowedCountryConfigModel(m cfgmodel.StringCSV) Option {
	return func(s *Service) {
		s.AllowedCountries = m
	}
}

func WithLogger(l log.Logger) Option {
	return func(s *Service) {
		s.Log = l
	}
}

// WithGeoIP2Reader creates a new GeoIP2.Reader. As long as there are no other
// readers this is a mandatory argument.
// Error behaviour: NotFound, NotValid
func WithGeoIP2Reader(file string) Option {
	return func(s *Service) {
		if s.optionError != nil {
			return
		}
		if false == os.FileExists(file) {
			s.optionError = errors.NewNotFoundf("[geoip] File %s not found", file)
			return
		}

		r, err := geoip2.Open(file) // that implementation is not nice for testing because no interface usages :(
		if err != nil {
			s.optionError = errors.NewNotValid(err, "[geoip] Maxmind Open")
			return
		}
		s.GeoIP = &mmdb{
			r: r,
		}
	}
}
