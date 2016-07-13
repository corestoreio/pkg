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

package backendgeoip

import (
	"encoding/gob"

	"github.com/corestoreio/csfw/config"
	"github.com/corestoreio/csfw/config/cfgmodel"
	"github.com/corestoreio/csfw/net/geoip"
	"github.com/corestoreio/csfw/storage/transcache"
	"github.com/corestoreio/csfw/storage/transcache/tcbigcache"
	"github.com/corestoreio/csfw/storage/transcache/tcredis"
	"github.com/corestoreio/csfw/util/errors"
)

func init() {
	gob.Register(geoip.Country{})
}

// PrepareOptions creates a closure around the type Backend. The closure will be
// used during a scoped request to figure out the configuration depending on the
// incoming scope. An option array will be returned by the closure.
func PrepareOptions(be *Backend) geoip.OptionFactoryFunc {

	return func(sg config.Scoped) []geoip.Option {
		var opts [6]geoip.Option
		var i int
		scp, id := sg.Scope()

		acc, err := be.NetGeoipAllowedCountries.Get(sg)
		if err != nil {
			return optError(errors.Wrap(err, "[backendgeoip] NetGeoipAllowedCountries.Get"))
		}
		opts[i] = geoip.WithAllowedCountryCodes(scp, id, acc...)
		i++

		// REDIRECT TO ALTERNATIVE URL
		ar, err := be.NetGeoipAlternativeRedirect.Get(sg)
		if err != nil {
			return optError(errors.Wrap(err, "[backendgeoip] NetGeoipAlternativeRedirect.Get"))
		}
		arc, err := be.NetGeoipAlternativeRedirectCode.Get(sg)
		if err != nil {
			return optError(errors.Wrap(err, "[backendgeoip] NetGeoipAlternativeRedirectCode.Get"))
		}
		if arc > 0 && ar != nil {
			opts[i] = geoip.WithAlternativeRedirect(scp, id, ar.String(), arc)
		}
		i++

		// LOCAL MAXMIND FILE
		mmlf, err := be.NetGeoipMaxmindLocalFile.Get(sg)
		if err != nil {
			return optError(errors.Wrap(err, "[backendgeoip] NetGeoipMaxmindLocalFile.Get"))
		}
		if mmlf != "" {
			opts[i] = geoip.WithGeoIP2File(mmlf)
			i++
			// we're done! skip the webservice part
			return opts[:]
		}

		// MAXMIND WEB SERVICE
		user, err := be.NetGeoipMaxmindWebserviceUserID.Get(sg)
		if err != nil {
			return optError(errors.Wrap(err, "[backendgeoip] NetGeoipMaxmindWebserviceUserID.Get"))
		}
		license, err := be.NetGeoipMaxmindWebserviceLicense.Get(sg)
		if err != nil {
			return optError(errors.Wrap(err, "[backendgeoip] NetGeoipMaxmindWebserviceLicense.Get"))
		}
		timeout, err := be.NetGeoipMaxmindWebserviceTimeout.Get(sg)
		if err != nil {
			return optError(errors.Wrap(err, "[backendgeoip] NetGeoipMaxmindWebserviceTimeout.Get"))
		}
		redisURL, err := be.NetGeoipMaxmindWebserviceRedisURL.Get(sg)
		if err != nil {
			return optError(errors.Wrap(err, "[backendgeoip] NetGeoipMaxmindWebserviceRedisURL.Get"))
		}

		var opt [2]transcache.Option
		switch {
		case redisURL != nil:
			opt[0] = tcredis.WithURL(redisURL.String(), nil, true)
		default:
			opt[0] = tcbigcache.With()
		}
		opt[1] = transcache.WithPooledEncoder(transcache.GobCodec{}, geoip.Country{}) // prime gob with the Country struct

		// for now only encoding/gob can be used, we might make it configurable
		// to choose the encoder/decoder.
		tc, err := transcache.NewProcessor(opt[:]...)
		if err != nil {
			return optError(errors.Wrap(err, "[backendgeoip] transcache.NewProcessor"))
		}

		if user != "" && license != "" && timeout > 0 {
			if be.WebServiceClient != nil {
				be.WebServiceClient.Timeout = timeout
				opts[i] = geoip.WithGeoIP2WebserviceHTTPClient(tc, user, license, be.WebServiceClient)
			} else {
				opts[i] = geoip.WithGeoIP2Webservice(tc, user, license, timeout)
			}
			i++
		}

		return opts[:]
	}
}
