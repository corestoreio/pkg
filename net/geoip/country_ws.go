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
	"encoding/json"
	"fmt"
	"net"
	"net/http"

	"github.com/corestoreio/csfw/sys/singleflight"
	"github.com/corestoreio/csfw/util/errors"
)

// TransCacher transcodes Go objects. It knows how to encode and cache any Go
// object and knows how to retrieve from cache and decode into a new Go object.
// Hint: use package storage/transcache.
type TransCacher interface {
	Set(key []byte, src interface{}) error
	// Get must return an errors.NotFound if a key does not exists.
	Get(key []byte, dst interface{}) error
}

// MaxMindWebserviceBaseURL defines the used base url. The IP address will be
// added after the last slash.
const MaxMindWebserviceBaseURL = "https://geoip.maxmind.com/geoip/v2.1/country/"

// mmws resolves to MaxMind WebService
type mmws struct {
	inflight   *singleflight.Group
	userID     string
	licenseKey string
	// client instantiated once and used for all queries to MaxMind.
	client *http.Client
	TransCacher
}

func newMMWS(t TransCacher, userID, licenseKey string, hc *http.Client) *mmws {
	mm := &mmws{
		inflight:    new(singleflight.Group),
		userID:      userID,
		licenseKey:  licenseKey,
		client:      hc,
		TransCacher: t,
	}

	return mm
}

// Country queries the MaxMind Webserver for one IP address. Implements the
// CountryRetriever interface. During concurrent requests with the same IP
// address it avoids querying the MaxMind database twice. It is guaranteed one
// request to MaxMind for an IP address. Those addresses gets cached in the
// Transcache along with the retrieved country.
func (mm *mmws) Country(ipAddress net.IP) (*Country, error) {

	var c = new(Country)
	err := mm.TransCacher.Get(ipAddress, c)
	if err != nil && !errors.IsNotFound(err) {
		return nil, errors.Wrap(err, "[geoip] mmws.Country.TransCacher.Get")
	}
	if err == nil {
		return c, nil
	}

	// runs the fetching of the HTTP result in another goroutine provided by DoChan()
	chResult := mm.inflight.DoChan(ipAddress.String(), func() (interface{}, error) {
		cntry, err := fetch(mm.client, mm.userID, mm.licenseKey, ipAddress)
		if err != nil {
			return nil, errors.Wrap(err, "[geoip] mmws.Country.Inflight.DoChan fetch() error")
		}
		if err := mm.TransCacher.Set(ipAddress, cntry); err != nil {
			return nil, errors.Wrap(err, "[geoip] mmws.Country.TransCacher.Set")
		}
		return cntry, nil
	})

	select {
	case res, ok := <-chResult:
		if !ok {
			return nil, errors.NewFatalf("[geoip] mmws.Country.Inflight.DoChan returned a closed/unreadable channel")
		}
		if res.Err != nil {
			return nil, errors.Wrap(res.Err, "[geoip] mmws.Country.Inflight.DoChan.Error")
		}
		if c, ok = res.Val.(*Country); ok {
			return c, nil
		} else {
			return nil, errors.NewFatalf("[geoip] mmws.Country.InflightDoChan res.Val cannot be type asserted to *Country")
		}
	}

	return nil, errors.Wrapf(err, "[geoip] mmws.Country: Nothing to select for IP %q", ipAddress)
}

func (mm *mmws) Close() error {
	return nil
}

//func (a *mmws) City(ipAddress net.IP) (internal.Response, error) {
//	return a.fetch("https://geoip.maxmind.com/geoip/v2.1/city/", ipAddress)
//}
//
//func (a *mmws) Insights(ipAddress net.IP) (internal.Response, error) {
//	return a.fetch("https://geoip.maxmind.com/geoip/v2.1/insights/", ipAddress)
//}

func fetch(hc *http.Client, userID, licenseKey string, ipAddress net.IP) (*Country, error) {
	var country = new(Country)
	req, err := http.NewRequest("GET", MaxMindWebserviceBaseURL+ipAddress.String(), nil)
	if err != nil {
		return country, errors.Wrap(err, "[geoip] http.NewRequest")
	}

	// authorize the request
	// http://dev.maxmind.com/geoip/geoip2/web-services/#Authorization
	req.SetBasicAuth(userID, licenseKey)

	resp, err := hc.Do(req) // execute the request
	if err != nil {
		return country, errors.Wrap(err, "[geoip] http.Client.Do")
	}
	defer resp.Body.Close()

	// handle errors that may occur
	// http://dev.maxmind.com/geoip/geoip2/web-services/#Response_Headers
	if resp.StatusCode >= 400 && resp.StatusCode < 600 {
		var v WebserviceError
		v.err = json.NewDecoder(resp.Body).Decode(&v)
		return nil, errors.NewNotValid(v, "[geoip] mmws.fetch URL: "+MaxMindWebserviceBaseURL)
	}

	// parse the response body
	// http://dev.maxmind.com/geoip/geoip2/web-services/#Response_Body

	if err := json.NewDecoder(resp.Body).Decode(country); err != nil {
		return nil, errors.NewNotValid(err, "[geoip] json.NewDecoder.Decode")
	}
	country.IP = ipAddress
	return country, nil
}

// WebserviceError used in the Maxmind Webservice functional option.
type WebserviceError struct {
	err  error
	Code string `json:"code,omitempty"`
	Err  string `json:"error,omitempty"`
}

func (e WebserviceError) Error() string {
	if e.err != nil {
		return e.err.Error()
	}
	return fmt.Sprintf("%s: %s", e.Code, e.Err)
}
