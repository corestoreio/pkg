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
	"fmt"
	"net"

	"github.com/corestoreio/csfw/config/cfgmodel"
	"github.com/corestoreio/csfw/net/ctxhttp"
	"github.com/corestoreio/csfw/store/scope"
	"github.com/corestoreio/csfw/util/os"
	"github.com/oschwald/geoip2-golang"
)

// Option can be used as an argument in NewService to configure a token service.
type Option func(*Service)

// WithAlternativeHandler sets for a scope.Scope and its ID the error handler
// on a Service. If the Handler h is nil falls back to the DefaultErrorHandler.
// This function can be called as many times as you have websites or stores.
// Group scope is not suppored.
func WithAlternativeHandler(so scope.Scope, id int64, hf ctxhttp.Handler) Option {
	if hf == nil {
		hf = DefaultAlternativeHandler
	}
	return func(s *Service) {
		switch so {
		case scope.Store:
			s.storeIDs.Append(id)
			s.storeAltH = append(s.storeAltH, hf)
		case scope.Website:
			s.websiteIDs.Append(id)
			s.websiteAltH = append(s.websiteAltH, hf)
		default:
			s.MultiErr = s.AppendErrors(scope.ErrUnsupportedScopeID)
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

// WithGeoIP2Reader creates a new GeoIP2.Reader. As long as there are no other
// readers this is a mandatory argument.
func WithGeoIP2Reader(file string) Option {
	return func(s *Service) {
		if false == os.FileExists(file) {
			s.MultiErr = s.AppendErrors(fmt.Errorf("File %s not found", file))
			return
		}

		r, err := geoip2.Open(file) // that implementation is not nice for testing because no interface usages :(
		if err != nil {
			s.MultiErr = s.AppendErrors(err)
			return
		}
		s.GeoIP = &mmdb{
			r: r,
		}
	}
}

var _ Reader = (*mmdb)(nil)

// mmdb internal wrapper between geoip2 and our interface
type mmdb struct {
	r *geoip2.Reader
}

func (mm *mmdb) Country(ipAddress net.IP) (*Country, error) {
	c, err := mm.r.Country(ipAddress)
	if err != nil {
		return nil, err
	}
	c2 := &Country{
		IP: ipAddress,
	}
	c2.Continent.Code = c.Continent.Code
	c2.Continent.GeoNameID = c.Continent.GeoNameID
	c2.Continent.Names = c.Continent.Names // ! a map those names, should maybe copied away

	c2.Country.GeoNameID = c.Country.GeoNameID
	c2.Country.IsoCode = c.Country.IsoCode
	c2.Country.Names = c.Country.Names

	c2.RegisteredCountry.GeoNameID = c.RegisteredCountry.GeoNameID
	c2.RegisteredCountry.IsoCode = c.RegisteredCountry.IsoCode
	c2.RegisteredCountry.Names = c.RegisteredCountry.Names

	c2.RepresentedCountry.GeoNameID = c.RepresentedCountry.GeoNameID
	c2.RepresentedCountry.IsoCode = c.RepresentedCountry.IsoCode
	c2.RepresentedCountry.Names = c.RepresentedCountry.Names
	c2.RepresentedCountry.Type = c.RepresentedCountry.Type

	c2.Traits.IsAnonymousProxy = c.Traits.IsAnonymousProxy
	c2.Traits.IsSatelliteProvider = c.Traits.IsSatelliteProvider

	return c2, nil
}

func (mm *mmdb) Close() error {
	return mm.r.Close()
}
