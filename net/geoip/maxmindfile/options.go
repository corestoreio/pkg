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
	"os"

	"github.com/corestoreio/cspkg/config"
	"github.com/corestoreio/cspkg/config/cfgmodel"
	"github.com/corestoreio/cspkg/net/geoip"
	"github.com/corestoreio/errors"
)

// WithCountryFinder creates a new GeoIP2.Reader which reads the geo information
// from a file stored on the server.
func WithCountryFinder(filename string) geoip.Option {
	return func(s *geoip.Service) error {
		if _, err := os.Stat(filename); os.IsNotExist(err) {
			return errors.NewNotFoundf("[maxmindfile] File %q not found", filename)
		}
		cr, err := newMMDBByFile(filename)
		if err != nil {
			return errors.NewNotValidf("[maxmindfile] Maxmind Open %s with file %q", err, filename)
		}
		return geoip.WithCountryFinder(cr)(s)
	}
}

// OptionName identifies this package within the register of the
// backendgeoip.Configuration type.
const OptionName = `file`

// NewOptionFactory specifies the file on the server to retrieve geo
// information. Alternatively you can choose the MaxMind web service via package
// maxmind.NewOptionFactoryWebservice(). This function will be triggered when
// you choose in backendgeoip.Configuration.DataSource the value `file`.
func NewOptionFactory(maxmindLocalFile cfgmodel.Str) (optionName string, _ geoip.OptionFactoryFunc) {
	return OptionName, func(sg config.Scoped) []geoip.Option {
		mmlf, err := maxmindLocalFile.Get(sg)
		if err != nil {
			return geoip.OptionsError(errors.Wrap(err, "[backendgeoip] NetGeoipMaxmindLocalFile.Get"))
		}
		if mmlf != "" {
			return []geoip.Option{
				WithCountryFinder(mmlf),
			}
		}
		return geoip.OptionsError(errors.NewEmptyf("[backendgeoip] Geo source as file specified but path to file name not provided"))
	}
}
