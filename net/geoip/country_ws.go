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
	"time"

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

// NewHttpClient creates a new HTTP client for the MaxMind webservice.
// You can provide here your own function or a mock for testing.
// This function will be used in the option function WithGeoIP2Webservice().
var NewHttpClient = func(timeout time.Duration) *http.Client {
	return &http.Client{Timeout: timeout}
}

// mmws resolves to MaxMind WebService
type mmws struct {
	userID     string
	licenseKey string
	// client instantiated once and used for all queries to MaxMind.
	client *http.Client
	TransCacher
}

func newMMWS(t TransCacher, userID, licenseKey string, httpTimeout time.Duration) *mmws {
	if httpTimeout < 1 {
		httpTimeout = time.Second * 20
	}
	return &mmws{
		userID:      userID,
		licenseKey:  licenseKey,
		client:      NewHttpClient(httpTimeout),
		TransCacher: t,
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
	c, err = mm.fetch("https://geoip.maxmind.com/geoip/v2.1/country/", ipAddress)
	if err == nil {
		if err2 := mm.TransCacher.Set(ipAddress, c); err2 != nil {
			return nil, errors.Wrap(err, "[geoip] mmws.Country.cacheSave")
		}
	}
	return c, errors.Wrap(err, "[geoip] mmws.Country.fetch")
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

func (a *mmws) fetch(prefix string, ipAddress net.IP) (*Country, error) {
	var response = new(Country)
	req, err := http.NewRequest("GET", prefix+ipAddress.String(), nil)
	if err != nil {
		return response, err
	}

	// authorize the request
	// http://dev.maxmind.com/geoip/geoip2/web-services/#Authorization
	req.SetBasicAuth(a.userID, a.licenseKey)

	// execute the request

	resp, err := a.client.Do(req)
	if err != nil {
		return response, err
	}
	defer resp.Body.Close()

	// handle errors that may occur
	// http://dev.maxmind.com/geoip/geoip2/web-services/#Response_Headers
	if resp.StatusCode >= 400 && resp.StatusCode < 600 {
		var v WebserviceError
		v.err = json.NewDecoder(resp.Body).Decode(&v)
		return nil, errors.NewNotValid(v, "[geoip] mmws.fetch URL: "+prefix)
	}

	// parse the response body
	// http://dev.maxmind.com/geoip/geoip2/web-services/#Response_Body

	err = json.NewDecoder(resp.Body).Decode(&response)
	return response, errors.NewNotValid(err, "[geoip] json.NewDecoder.Decode")
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
