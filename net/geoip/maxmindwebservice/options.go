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

package maxmindwebservice

import (
	"encoding/gob"
	"net/http"
	"time"

	"github.com/corestoreio/csfw/config"
	"github.com/corestoreio/csfw/config/cfgmodel"
	"github.com/corestoreio/csfw/net/geoip"
	"github.com/corestoreio/csfw/storage/transcache"
	"github.com/corestoreio/csfw/storage/transcache/tcbigcache"
	"github.com/corestoreio/csfw/storage/transcache/tcredis"
	"github.com/corestoreio/errors"
)

// required by the transcache package
func init() {
	gob.Register(geoip.Country{})
}

// WithCountryFinder uses for each incoming a request a lookup request to the
// Maxmind Webservice http://dev.maxmind.com/geoip/geoip2/web-services/ and
// caches the result in Transcacher. Hint: use package storage/transcache. If
// the httpTimeout is lower 0 then the default 20s get applied.
func WithCountryFinder(t TransCacher, userID, licenseKey string, httpTimeout time.Duration) geoip.Option {
	if httpTimeout < 1 {
		httpTimeout = time.Second * 20
	}
	return WithCountryFinderHTTPClient(&http.Client{Timeout: httpTimeout}, t, userID, licenseKey)
}

// WithCountryFinderHTTPClient uses for each incoming a request a lookup
// request to the Maxmind Webservice
// http://dev.maxmind.com/geoip/geoip2/web-services/ and caches the result in
// Transcacher. Hint: use package storage/transcache.
func WithCountryFinderHTTPClient(hc *http.Client, t TransCacher, userID, licenseKey string) geoip.Option {
	return geoip.WithCountryFinder(newMMWS(t, userID, licenseKey, hc))
}

// OptionName identifies this package within the register of the
// backendgeoip.Configuration type.
const OptionName = `webservice`

// NewOptionFactory creates a new option factory function for the MaxMind web
// service in the backend package to be used for automatic scope based
// configuration initialization. Configuration values must be set from package
// backendgeoip.Configuration.
//
// First argument http.Client allows you to use a custom client when making
// requests to the MaxMind webservice. The timeout gets set by configuration
// path MaxmindWebserviceTimeout.
//
// gob.Register(geoip.Country{}) has already been called.
func NewOptionFactory(hc *http.Client, userID, license cfgmodel.Str, timeout cfgmodel.Duration, redisURL cfgmodel.URL) (optionName string, _ geoip.OptionFactoryFunc) {
	return OptionName, func(sg config.Scoped) []geoip.Option {

		vUserID, err := userID.Get(sg)
		if err != nil {
			return geoip.OptionsError(errors.Wrap(err, "[maxmindwebservice] MaxmindWebserviceUserID.Get"))
		}
		vLicense, err := license.Get(sg)
		if err != nil {
			return geoip.OptionsError(errors.Wrap(err, "[maxmindwebservice] MaxmindWebserviceLicense.Get"))
		}
		vTimeout, err := timeout.Get(sg)
		if err != nil {
			return geoip.OptionsError(errors.Wrap(err, "[maxmindwebservice] MaxmindWebserviceTimeout.Get"))
		}
		vRedisURL, err := redisURL.Get(sg)
		if err != nil {
			return geoip.OptionsError(errors.Wrap(err, "[maxmindwebservice] MaxmindWebserviceRedisURL.Get"))
		}

		var tco [2]transcache.Option
		switch {
		case vRedisURL != nil:
			tco[0] = tcredis.WithURL(vRedisURL.String(), nil, true)
		default:
			tco[0] = tcbigcache.With()
		}
		tco[1] = transcache.WithPooledEncoder(transcache.GobCodec{}, geoip.Country{}) // prime gob with the Country struct

		// for now only encoding/gob can be used, we might make it configurable
		// to choose the encoder/decoder.
		tc, err := transcache.NewProcessor(tco[:]...)
		if err != nil {
			return geoip.OptionsError(errors.Wrap(err, "[maxmindwebservice] transcache.NewProcessor"))
		}
		if vUserID == "" || vLicense == "" || vTimeout < 1 {
			return geoip.OptionsError(errors.NewNotValidf("[maxmindwebservice] Incomplete WebService configuration: User: %q License %q Timeout: %d (zero timeout not supported)", vUserID, vLicense, vTimeout))
		}
		hc.Timeout = vTimeout

		return []geoip.Option{
			WithCountryFinderHTTPClient(hc, tc, vUserID, vLicense),
		}
	}
}
