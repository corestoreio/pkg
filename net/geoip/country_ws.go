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
	"io"
	"io/ioutil"
	"net"
	"net/http"
	"runtime"

	"github.com/corestoreio/csfw/util/errors"
)

// TransCacher transcodes Go objects. It knows how to encode and cache any
// Go object and knows how to retrieve from cache and decode into a new Go object.
// Hint: use package storage/transcache.
type TransCacher interface {
	Set(key []byte, src interface{}) error
	// Get must return an errors.NotFound if a key does not exists.
	Get(key []byte, dst interface{}) error
}

// mmws resolves to MaxMind WebService
type mmws struct {
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
				return
			}
			c, err := fetch(mm.client, mm.userID, mm.licenseKey, "https://geoip.maxmind.com/geoip/v2.1/country/", ip)
			if err != nil {
				errOUT <- errors.Wrap(err, "[geoip] mmws.Country.fetch")
			}

			if err := mm.TransCacher.Set(ip, c); err != nil {
				errOUT <- errors.Wrap(err, "[geoip] mmws.Country.cacheSave")
			}
			cOUT <- c
		}
	}
}

func (mm *mmws) Country(ipAddress net.IP) (*Country, error) {
	var c = new(Country)
	err := mm.TransCacher.Get(ipAddress, c)
	if err != nil && !errors.IsNotFound(err) {
		return nil, errors.Wrap(err, "[geoip] mmws.Country.TransCacher.Get")
	}
	if err == nil {
		return c, nil
	}
	mm.ipIN <- ipAddress
	select {
	case ctry := <-mm.cOUT:

		if !ctry.IP.Equal(ipAddress) {
			// can this happen?
			panic(fmt.Sprintf("Dude, a bug! IPs are not equal Have %s Want %s", ctry.IP, ipAddress))
		}

		return ctry, nil
	case err := <-mm.errOUT:
		return nil, errors.Wrap(err, "[geoip] mmws.Country.fetch.errOUT")
	}
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
		return country, err
	}

	// authorize the request
	// http://dev.maxmind.com/geoip/geoip2/web-services/#Authorization
	req.SetBasicAuth(userID, licenseKey)

	// execute the request

	resp, err := hc.Do(req)
	if err != nil {
		return country, err
	}
	defer resp.Body.Close()
	defer func() {
		// read until the response is complete
		if _, err := io.Copy(ioutil.Discard, resp.Body); err != nil {
			panic(err) // todo remove panic or find another better way to avoid this
		}
	}()

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
