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

package directory_test

import (
	"testing"

	"github.com/corestoreio/csfw/config"
	"github.com/corestoreio/csfw/config/model"
	"github.com/corestoreio/csfw/directory"
	"github.com/corestoreio/csfw/store/scope"
	"github.com/stretchr/testify/assert"
)

func init() {

	directory.InitSources(nil)

}

func TestPathCountryAllowedCustom(t *testing.T) {
	t.Parallel()
	defer debugLogBuf.Reset()

	previous := directory.Backend.GeneralCountryAllow.Option(model.WithSourceByString(
		"DE", "Germany", "AU", "'Straya", "CH", "Switzerland",
	))
	defer directory.Backend.GeneralCountryAllow.Option(previous)

	cr := config.NewMockGetter(
		config.WithMockValues(config.MockPV{
			directory.Backend.GeneralCountryAllow.MustFQPathInt64(scope.StrStores, 1): "DE,AU,CH,AT",
		}),
	)

	haveCountries := directory.Backend.GeneralCountryAllow.Get(cr.NewScoped(1, 1, 1))

	assert.Exactly(t, []string{"DE", "AU", "CH", "AT"}, haveCountries)
	// todo validation
}

func TestPathGeneralCountryAllowDefault(t *testing.T) {
	t.Parallel()
	defer debugLogBuf.Reset()

	cr := config.NewMockGetter(
		config.WithMockValues(config.MockPV{}),
	)

	haveCountries := directory.Backend.GeneralCountryAllow.Get(cr.NewScoped(1, 1, 1))

	assert.Exactly(t,
		[]string{"AF", "AL", "DZ", "AS", "AD", "AO", "AI", "AQ", "AG", "AR", "AM", "AW", "AU", "AT", "AX", "AZ", "BS", "BH", "BD", "BB", "BY", "BE", "BZ", "BJ", "BM", "BL", "BT", "BO", "BA", "BW", "BV", "BR", "IO", "VG", "BN", "BG", "BF", "BI", "KH", "CM", "CA", "CD", "CV", "KY", "CF", "TD", "CL", "CN", "CX", "CC", "CO", "KM", "CG", "CK", "CR", "HR", "CU", "CY", "CZ", "DK", "DJ", "DM", "DO", "EC", "EG", "SV", "GQ", "ER", "EE", "ET", "FK", "FO", "FJ", "FI", "FR", "GF", "PF", "TF", "GA", "GM", "GE", "DE", "GG", "GH", "GI", "GR", "GL", "GD", "GP", "GU", "GT", "GN", "GW", "GY", "HT", "HM", "HN", "HK", "HU", "IS", "IM", "IN", "ID", "IR", "IQ", "IE", "IL", "IT", "CI", "JE", "JM", "JP", "JO", "KZ", "KE", "KI", "KW", "KG", "LA", "LV", "LB", "LS", "LR", "LY", "LI", "LT", "LU", "ME", "MF", "MO", "MK", "MG", "MW", "MY", "MV", "ML", "MT", "MH", "MQ", "MR", "MU", "YT", "FX", "MX", "FM", "MD", "MC", "MN", "MS", "MA", "MZ", "MM", "NA", "NR", "NP", "NL", "AN", "NC", "NZ", "NI", "NE", "NG", "NU", "NF", "KP", "MP", "NO", "OM", "PK", "PW", "PA", "PG", "PY", "PE", "PH", "PN", "PL", "PS", "PT", "PR", "QA", "RE", "RO", "RS", "RU", "RW", "SH", "KN", "LC", "PM", "VC", "WS", "SM", "ST", "SA", "SN", "SC", "SL", "SG", "SK", "SI", "SB", "SO", "ZA", "GS", "KR", "ES", "LK", "SD", "SR", "SJ", "SZ", "SE", "CH", "SY", "TL", "TW", "TJ", "TZ", "TH", "TG", "TK", "TO", "TT", "TN", "TR", "TM", "TC", "TV", "VI", "UG", "UA", "AE", "GB", "US", "UM", "UY", "UZ", "VU", "VA", "VE", "VN", "WF", "EH", "YE", "ZM", "ZW"},
		haveCountries,
	)
}
