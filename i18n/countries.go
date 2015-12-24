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

package i18n

import (
	"golang.org/x/text/language"
	"golang.org/x/text/language/display"
)

type (
	CountrySlice []language.Region
)

// Countries contains all supported countries
var Countries CountrySlice

// init countries
func init() {
	for _, r := range display.Values.Regions() {
		if r.IsCountry() {
			Countries = append(Countries, r)
		}
	}
}
