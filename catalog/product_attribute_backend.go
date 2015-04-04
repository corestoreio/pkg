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

package catalog

import "github.com/corestoreio/csfw/eav"

var (
	_ eav.AttributeBackendModeller = (*todoPABP)(nil)
	_ eav.AttributeBackendModeller = (*todoPABSD)(nil)
	_ eav.AttributeBackendModeller = (*todoPABB)(nil)
	_ eav.AttributeBackendModeller = (*todoPABGP)(nil)
	_ eav.AttributeBackendModeller = (*todoPABTP)(nil)
	_ eav.AttributeBackendModeller = (*todoPABM)(nil)
	_ eav.AttributeBackendModeller = (*todoPABR)(nil)
	_ eav.AttributeBackendModeller = (*todoPABSku)(nil)
	_ eav.AttributeBackendModeller = (*todoPABStock)(nil)
)

type (
	todoPABP struct {
		*eav.AttributeBackend
	}
	todoPABSD struct {
		*eav.AttributeBackend
	}
	todoPABB struct {
		*eav.AttributeBackend
	}
	todoPABGP struct {
		*eav.AttributeBackend
	}
	todoPABTP struct {
		*eav.AttributeBackend
	}
	todoPABM struct {
		*eav.AttributeBackend
	}
	todoPABR struct {
		*eav.AttributeBackend
	}
	todoPABSku struct {
		*eav.AttributeBackend
	}
	todoPABStock struct {
		*eav.AttributeBackend
	}
)

// ProductAttributeBackendPrice prices @todo
// @see magento2/site/app/code/Magento/Catalog/Model/Product/Attribute/Backend/Price.php
func ProductAttributeBackendPrice() *todoPABP {
	return &todoPABP{}
}

// ProductAttributeBackendStartDate @todo
// @see magento2/site/app/code/Magento/Catalog/Model/Product/Attribute/Backend/Startdate.php
func ProductAttributeBackendStartDate() *todoPABSD {
	return &todoPABSD{}
}

// ProductAttributeBackendBoolean @todo
// @see magento2/site/app/code/Magento/Catalog/Model/Product/Attribute/Backend/Boolean.php
func ProductAttributeBackendBoolean() *todoPABB {
	return &todoPABB{}
}

// ProductAttributeBackendGroupPrice @todo ... pretty complex
// @see magento2/site/app/code/Magento/Catalog/Model/Product/Attribute/Backend/GroupPrice.php
// @see magento2/site/app/code/Magento/Catalog/Model/Product/Attribute/Backend/GroupPrice/AbstractGroupPrice.php
func ProductAttributeBackendGroupPrice() *todoPABGP {
	return &todoPABGP{}
}

// IsScalar returns false because a group price is not a scalar type
func (gp *todoPABGP) IsScalar() bool {
	return false
}

// ProductAttributeBackendTierPrice @todo ... pretty complex
// @see magento2/site/app/code/Magento/Catalog/Model/Product/Attribute/Backend/TierPrice.php
// @see magento2/site/app/code/Magento/Catalog/Model/Product/Attribute/Backend/GroupPrice/AbstractGroupPrice.php
func ProductAttributeBackendTierPrice() *todoPABTP {
	return &todoPABTP{}
}

// IsScalar returns false because a tier price is not a scalar type
func (gp *todoPABTP) IsScalar() bool {
	return false
}

// ProductAttributeBackendMedia @todo ... pretty complex
// @see magento2/site/app/code/Magento/Catalog/Model/Product/Attribute/Backend/Media.php
func ProductAttributeBackendMedia() *todoPABM {
	return &todoPABM{}
}

// IsScalar returns false because a media is not a scalar type
func (gp *todoPABM) IsScalar() bool {
	return false
}

// ProductAttributeRecurring @todo
// @see Mage_Catalog_Model_Product_Attribute_Backend_Recurring
func ProductAttributeBackendRecurring() *todoPABR {
	return &todoPABR{}
}

// ProductAttributeBackendSku @todo
// @see magento2/site/app/code/Magento/Catalog/Model/Product/Attribute/Backend/Sku.php
func ProductAttributeBackendSku() *todoPABSku {
	return &todoPABSku{}
}

// ProductAttributeBackendStock @todo
// @see magento2/site/app/code/Magento/Catalog/Model/Product/Attribute/Backend/Stock.php
func ProductAttributeBackendStock() *todoPABStock {
	return &todoPABStock{}
}
