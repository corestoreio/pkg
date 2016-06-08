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
	"hash/fnv"
	"net"
	"net/http"
	"runtime"

	"github.com/corestoreio/csfw/storage/suspend"
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

// MaxmindWebserviceBaseURL defines the used base url. The IP address will be
// added after the last slash.
const MaxmindWebserviceBaseURL = "https://geoip.maxmind.com/geoip/v2.1/country/"

// mmws resolves to MaxMind WebService
type mmws struct {
	suspend.State
	ipIN       chan net.IP
	cOUT       chan *Country
	errOUT     chan error
	userID     string
	licenseKey string
	// client instantiated once and used for all queries to MaxMind.
	client *http.Client
	TransCacher
}

func newMMWS(t TransCacher, userID, licenseKey string, hc *http.Client) *mmws {
	mm := &mmws{
		State:       suspend.NewStateWithHash(fnv.New64a()),
		ipIN:        make(chan net.IP),
		cOUT:        make(chan *Country),
		errOUT:      make(chan error),
		userID:      userID,
		licenseKey:  licenseKey,
		client:      hc,
		TransCacher: t,
	}

	for i := 0; i < runtime.NumCPU(); i++ {
		go workfetch(mm, mm.ipIN, mm.cOUT, mm.errOUT)
	}

	return mm
}

func workfetch(mm *mmws, ipIN <-chan net.IP, cOUT chan<- *Country, errOUT chan<- error) {
	for {
		select {
		case ip, ok := <-ipIN:
			if !ok {
				// channel closed, so quit
				return
			}
			c, err := fetch(mm.client, mm.userID, mm.licenseKey, MaxmindWebserviceBaseURL, ip)
			if err != nil {
				errOUT <- errors.Wrap(err, "[geoip] mmws.Country.fetch")
			} else {
				cOUT <- c
			}
		}
	}
}

// Country queries the MaxMind Webserver for one IP address. Implements the CountryRetriever interface.
// During concurrent requests with the same IP address it avoids querying the MaxMind
// database twice. It is guaranteed one request to MaxMind for an IP address. Those
// addresses gets cached in the Transcache along with the retrieved country.
func (mm *mmws) Country(ipAddress net.IP) (*Country, error) {

	var c = new(Country)
	err := mm.TransCacher.Get(ipAddress, c)
	if err != nil && !errors.IsNotFound(err) {
		return nil, errors.Wrap(err, "[geoip] mmws.Country.TransCacher.Get")
	}

	switch {
	case err == nil:
		// cache hit
		return c, nil
	case mm.ShouldStartBytes(ipAddress):
		defer mm.DoneBytes(ipAddress) // send Signal and release waiter
		mm.ipIN <- ipAddress
		select {
		case err := <-mm.errOUT:
			return nil, errors.Wrap(err, "[geoip] mmws.Country.fetch.errOUT")
		case c = <-mm.cOUT:
			if err := mm.TransCacher.Set(ipAddress, c); err != nil {
				return nil, errors.Wrap(err, "[geoip] mmws.Country.cacheSave")
			}
			if !c.IP.Equal(ipAddress) {
				// todo limit recursion
				// call itself as long until we get our ip. This is pretty rare and no idea
				// how to 100% test it
				return mm.Country(ipAddress)
			}
			return c, nil
		}
	case mm.ShouldWaitBytes(ipAddress):
		// try again ...
		err := mm.TransCacher.Get(ipAddress, c)
		if err != nil && !errors.IsNotFound(err) {
			return nil, errors.Wrap(err, "[geoip] mmws.Country.TransCacher.Get")
		}
		return c, err // can be a not-found error
	}

	return nil, errors.NewFatalf("[geoip] mmws.Country unreachable code and you reached it 8-)")
}

func (mm *mmws) Close() error {
	close(mm.ipIN)
	return nil
}

//func (a *mmws) City(ipAddress net.IP) (internal.Response, error) {
//	return a.fetch("https://geoip.maxmind.com/geoip/v2.1/city/", ipAddress)
//}
//
//func (a *mmws) Insights(ipAddress net.IP) (internal.Response, error) {
//	return a.fetch("https://geoip.maxmind.com/geoip/v2.1/insights/", ipAddress)
//}

func fetch(hc *http.Client, userID, licenseKey, baseURL string, ipAddress net.IP) (*Country, error) {
	var country = new(Country)
	req, err := http.NewRequest("GET", baseURL+ipAddress.String(), nil)
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
		return nil, errors.NewNotValid(v, "[geoip] mmws.fetch URL: "+baseURL)
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
