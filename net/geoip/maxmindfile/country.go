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

package maxmindfile

import (
	"net"

	"github.com/corestoreio/cspkg/net/geoip"
	"github.com/corestoreio/errors"
	"github.com/oschwald/geoip2-golang"
)

// mmdb internal wrapper between geoip2 and our interface
type mmdb struct {
	r *geoip2.Reader
}

func newMMDBByFile(filename string) (*mmdb, error) {
	r, err := geoip2.Open(filename)
	return &mmdb{r}, errors.NewNotValid(err, "[geoip] Maxmind Open")
}

func (mm *mmdb) FindCountry(ipAddress net.IP) (*geoip.Country, error) {
	c, err := mm.r.Country(ipAddress)
	if err != nil {
		return nil, errors.NewNotValid(err, "[geoip] mmdb.Country")
	}
	c2 := &geoip.Country{
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
