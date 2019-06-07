// Copyright 2015-present, Cyrill @ Schumacher.fm and the CoreStore contributors
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

package custattr

import "github.com/corestoreio/pkg/eav"

// AddressSourceCountry retrieves slice of countries @todo
// @see magento2/site/app/code/Magento/Customer/Model/Resource/Address/Attribute/Source/Country.php
func AddressSourceCountry() *eav.AttributeSource {
	return eav.NewAttributeSource(
		// temporary because later these values comes from another slice/container/database
		func(as *eav.AttributeSource) {
			as.Source = []string{
				"AU", "Straya",
				"NZ", "Kiwi land",
				"DE", "Autobahn",
				"SE", "Smørrebrød",
			}
		},
	)
}

// AddressSourceRegion
// @see magento2/site/app/code/Magento/Customer/Model/Resource/Address/Attribute/Source/Region.php
func AddressSourceRegion() *eav.AttributeSource {
	return eav.NewAttributeSource(
		// temporary because later these values comes from another slice/container/database
		func(as *eav.AttributeSource) {
			as.Source = []string{
				"BAY", "Bavaria",
				"BAW", "Baden-Würstchenberg",
				"HAM", "Hamburg",
				"BER", "Bärlin",
			}
		},
	)
}
