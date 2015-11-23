// Copyright 2015, Cyrill @ Schumacher.fm and the CoreStore contributors
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

package directory

import (
	"fmt"
	"strings"

	"github.com/corestoreio/csfw/config"
	"github.com/corestoreio/csfw/utils"
)

type (
	Country struct {
	}
)

// DefaultCountry returns the country code. Store argument is optional.
func DefaultCountry(cr config.ScopedGetter) string {
	return cr.String(PathDefaultCountry)
}

// AllowedCountries returns a list of all allowed countries per scope.
// This function might gets refactored into a SourceModel.
// May return nil,nil when it's unable to determine any list of countries.
func AllowedCountries(cr config.ScopedGetter) (utils.StringSlice, error) {
	cStr := cr.String(PathCountryAllowed)

	if cStr == "" {
		field, err := PackageConfiguration.FindFieldByPath(PathCountryAllowed) // get default value
		if err != nil {
			return nil, err
		}

		var ok bool
		cStr, ok = field.Default.(string)
		if cStr == "" || !ok {
			return nil, fmt.Errorf("Cannot type assert field.Default value to string: %#v", field)
		}
	}
	return utils.StringSlice(strings.Split(cStr, `,`)), nil
}
