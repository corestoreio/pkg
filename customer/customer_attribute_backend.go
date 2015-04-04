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
	_ eav.AttributeBackendModeller = (*todoABB)(nil)
	_ eav.AttributeBackendModeller = (*todoABS)(nil)
	_ eav.AttributeBackendModeller = (*todoABDB)(nil)
	_ eav.AttributeBackendModeller = (*todoABP)(nil)
	_ eav.AttributeBackendModeller = (*todoABStore)(nil)
	_ eav.AttributeBackendModeller = (*todoABWebsite)(nil)
)

type (
	todoABB struct {
		*eav.AttributeBackend
	}
	todoABS struct {
		*eav.AttributeBackend
	}
	todoABDB struct {
		*eav.AttributeBackend
	}
	todoABP struct {
		*eav.AttributeBackend
	}
	todoABStore struct {
		*eav.AttributeBackend
	}
	todoABWebsite struct {
		*eav.AttributeBackend
	}
)

// AttributeBackendBilling handles billing address @todo
// @see magento2/site/app/code/Magento/Customer/Model/Customer/Attribute/Backend/Billing.php
func AttributeBackendBilling() *todoABB {
	return &todoABB{}
}

// AttributeBackendShipping handles shipping address @todo
// @see magento2/site/app/code/Magento/Customer/Model/Customer/Attribute/Backend/Shipping.php
func AttributeBackendShipping() *todoABS {
	return &todoABS{}
}

// AttributeBackendDataBoolean converts 1 or 0 to bool @todo
// @see magento2/site/app/code/Magento/Customer/Model/Attribute/Backend/Data/Boolean.php
func AttributeBackendDataBoolean() *todoABDB {
	return &todoABDB{}
}

// AttributeBackendPassword handles customer passwords @todo
// @see magento2/site/app/code/Magento/Customer/Model/Customer/Attribute/Backend/Password.php
func AttributeBackendPassword() *todoABP {
	return &todoABP{}
}

// AttributeBackendStore handles store
// @see magento2/site/app/code/Magento/Customer/Model/Customer/Attribute/Backend/Store.php
func AttributeBackendStore() *todoABStore {
	return &todoABStore{}
}

// AttributeBackendWebsite handles website
// @see magento2/site/app/code/Magento/Customer/Model/Customer/Attribute/Backend/Website.php
func AttributeBackendWebsite() *todoABWebsite {
	return &todoABWebsite{}
}
