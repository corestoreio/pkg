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
	"net"

	"github.com/corestoreio/csfw/util/errors"
	"github.com/oschwald/geoip2-golang"
)

// think about that ...
//type CountryTrimmed struct {
//	// IP contains the request IP address even if we run behind a proxy
//	IP      net.IP `json:"ip,omitempty"`
//	Country struct {
//		Confidence int               `json:"confidence,omitempty"`
//		GeoNameID  uint              `json:"geoname_id,omitempty"`
//		IsoCode    string            `json:"iso_code,omitempty"`
//	} `json:"country,omitempty"`
//	MaxMind struct {
//		QueriesRemaining int `json:"queries_remaining,omitempty"`
//	} `json:"maxmind,omitempty"`
//}

// The Country structure corresponds to the data in the GeoIP2/GeoLite2
// Country databases.
type Country struct {
	// IP contains the request IP address even if we run behind a proxy
	IP   net.IP `json:"ip,omitempty"`
	City struct {
		Confidence int               `json:"confidence,omitempty"`
		GeoNameID  uint              `json:"geoname_id,omitempty"`
		Names      map[string]string `json:"names,omitempty"`
	} `json:"city,omitempty"`
	Continent struct {
		Code      string            `json:"code,omitempty"`
		GeoNameID uint              `json:"geoname_id,omitempty"`
		Names     map[string]string `json:"names,omitempty"`
	} `json:"continent,omitempty"`
	Country struct {
		Confidence int               `json:"confidence,omitempty"`
		GeoNameID  uint              `json:"geoname_id,omitempty"`
		IsoCode    string            `json:"iso_code,omitempty"`
		Names      map[string]string `json:"names,omitempty"`
	} `json:"country,omitempty"`
	Location struct {
		AccuracyRadius    int     `json:"accuracy_radius,omitempty"`
		AverageIncome     int     `json:"average_income,omitempty"`
		Latitude          float64 `json:"latitude,omitempty"`
		Longitude         float64 `json:"longitude,omitempty"`
		MetroCode         int     `json:"metro_code,omitempty"`
		PopulationDensity int     `json:"population_density,omitempty"`
		TimeZone          string  `json:"time_zone,omitempty"`
	} `json:"location,omitempty"`
	Postal struct {
		Code       string `json:"code,omitempty"`
		Confidence int    `json:"confidence,omitempty"`
	} `json:"postal,omitempty"`
	RegisteredCountry struct {
		GeoNameID uint              `json:"geoname_id,omitempty"`
		IsoCode   string            `json:"iso_code,omitempty"`
		Names     map[string]string `json:"names,omitempty"`
	} `json:"registered_country,omitempty"`
	RepresentedCountry struct {
		GeoNameID uint              `json:"geoname_id,omitempty"`
		IsoCode   string            `json:"iso_code,omitempty"`
		Names     map[string]string `json:"names,omitempty"`
		Type      string            `json:"type,omitempty"`
	} `json:"represented_country,omitempty"`
	Subdivision []struct {
		Confidence int               `json:"confidence,omitempty"`
		GeoNameID  uint              `json:"geoname_id,omitempty"`
		IsoCode    string            `json:"iso_code,omitempty"`
		Names      map[string]string `json:"names,omitempty"`
	} `json:"subdivisions,omitempty"`
	Traits struct {
		AutonomousSystemNumber       int    `json:"autonomous_system_number,omitempty"`
		AutonomousSystemOrganization string `json:"autonomous_system_organization,omitempty"`
		Domain                       string `json:"domain,omitempty"`
		IsAnonymousProxy             bool   `json:"is_anonymous_proxy,omitempty"`
		IsSatelliteProvider          bool   `json:"is_satellite_provider,omitempty"`
		Isp                          string `json:"isp,omitempty"`
		IPAddress                    string `json:"ip_address,omitempty"`
		Organization                 string `json:"organization,omitempty"`
		UserType                     string `json:"user_type,omitempty"`
	} `json:"traits,omitempty"`
	MaxMind struct {
		QueriesRemaining int `json:"queries_remaining,omitempty"`
	} `json:"maxmind,omitempty"`
}

// CountryRetriever implements how to lookup the Country for an IP address.
// Supports IPv4 and IPv6 addresses.
type CountryRetriever interface {
	// Country todo add context for cancelling
	Country(net.IP) (*Country, error)
	// Close may be called on shutdown of the overall app and terminates
	// the underlying lookup service.
	Close() error
}

// mmdb internal wrapper between geoip2 and our interface
type mmdb struct {
	r *geoip2.Reader
}

func newMMDBByFile(filename string) (*mmdb, error) {
	r, err := geoip2.Open(filename)
	return &mmdb{r}, errors.NewNotValid(err, "[geoip] Maxmind Open")
}

func (mm *mmdb) Country(ipAddress net.IP) (*Country, error) {
	c, err := mm.r.Country(ipAddress)
	if err != nil {
		return nil, errors.NewNotValid(err, "[geoip] mmdb.Country")
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
