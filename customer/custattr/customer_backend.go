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

package custattr

import "github.com/corestoreio/csfw/eav"

var (
	_ eav.AttributeBackendModeller = (*todoABB)(nil)
)

type (
	// example for your own struct
	todoABB struct {
		*eav.AttributeBackend
	}
)

// CustomerBackendBilling handles billing address @todo
// @see magento2/site/app/code/Magento/Customer/Model/Customer/Attribute/Backend/Billing.php
func CustomerBackendBilling() *todoABB {
	return &todoABB{
		AttributeBackend: eav.NewAttributeBackend(),
	}
}

// CustomerBackendShipping handles shipping address @todo
// @see magento2/site/app/code/Magento/Customer/Model/Customer/Attribute/Backend/Shipping.php
func CustomerBackendShipping() *eav.AttributeBackend {
	return eav.NewAttributeBackend()
}

// CustomerBackendDataBoolean converts 1 or 0 to bool @todo
// @see magento2/site/app/code/Magento/Customer/Model/Attribute/Backend/Data/Boolean.php
func CustomerBackendDataBoolean() *eav.AttributeBackend {
	return eav.NewAttributeBackend()
}

// CustomerBackendPassword handles customer passwords @todo
// @see magento2/site/app/code/Magento/Customer/Model/Customer/Attribute/Backend/Password.php
func CustomerBackendPassword() *eav.AttributeBackend {
	return eav.NewAttributeBackend()
}

// CustomerBackendStore handles store
// @see magento2/site/app/code/Magento/Customer/Model/Customer/Attribute/Backend/Store.php
func CustomerBackendStore() *eav.AttributeBackend {
	return eav.NewAttributeBackend()
}

// CustomerBackendWebsite handles website
// @see magento2/site/app/code/Magento/Customer/Model/Customer/Attribute/Backend/Website.php
func CustomerBackendWebsite() *eav.AttributeBackend {
	return eav.NewAttributeBackend()
}
