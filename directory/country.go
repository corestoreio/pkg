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
	"github.com/corestoreio/csfw/config/valuelabel"
	"github.com/corestoreio/csfw/storage/dbr"
)

var CountryCollection valuelabel.Slice

func InitCountryCollection(dbrsess dbr.SessionRunner) error {
	// CountryCollection
	// load from database the iso code and as value the names

	// apply the list of country codes to

	PathDefaultCountry.ValueLabel = CountryCollection
	PathCountryAllowed.ValueLabel = CountryCollection

	return nil
}
