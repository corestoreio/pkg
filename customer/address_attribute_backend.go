// Copyright 2015 CoreStore Authors
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

package customer

import "github.com/corestoreio/csfw/eav"

var (
	_ eav.AttributeBackendModeller = (*todoAABR)(nil)
	_ eav.AttributeBackendModeller = (*todoAABS)(nil)
)

type (
	todoAABR struct {
		*eav.AttributeBackend
	}
	todoAABS struct {
		*eav.AttributeBackend
	}
)

// AddressAttributeDataPostcode post code data model @todo
// @see magento2/site/app/code/Magento/Customer/Model/Attribute/Data/Postcode.php
func AddressAttributeBackendRegion() *todoAABR {
	return &todoAABR{}
}

// AddressAttributeBackendStreet handles multiline street address
// @see Mage_Customer_Model_Resource_Address_Attribute_Backend_Street
func AddressAttributeBackendStreet() *todoAABS {
	return &todoAABS{}
}
