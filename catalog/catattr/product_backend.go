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

package catattr

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

// ProductBackendPrice prices @todo
// @see magento2/site/app/code/Magento/Catalog/Model/Product/Attribute/Backend/Price.php
func ProductBackendPrice() *todoPABP {
	return &todoPABP{
		AttributeBackend: eav.NewAttributeBackend(),
	}
}

// ProductBackendStartDate @todo
// @see magento2/site/app/code/Magento/Catalog/Model/Product/Attribute/Backend/Startdate.php
func ProductBackendStartDate() *todoPABSD {
	return &todoPABSD{
		AttributeBackend: eav.NewAttributeBackend(),
	}
}

// ProductBackendBoolean @todo
// @see magento2/site/app/code/Magento/Catalog/Model/Product/Attribute/Backend/Boolean.php
func ProductBackendBoolean() *todoPABB {
	return &todoPABB{
		AttributeBackend: eav.NewAttributeBackend(),
	}
}

// ProductBackendGroupPrice @todo ... pretty complex
// @see magento2/site/app/code/Magento/Catalog/Model/Product/Attribute/Backend/GroupPrice.php
// @see magento2/site/app/code/Magento/Catalog/Model/Product/Attribute/Backend/GroupPrice/AbstractGroupPrice.php
func ProductBackendGroupPrice() *todoPABGP {
	return &todoPABGP{
		AttributeBackend: eav.NewAttributeBackend(),
	}
}

// IsScalar returns false because a group price is not a scalar type
func (gp *todoPABGP) IsScalar() bool {
	return false
}

// ProductBackendTierPrice @todo ... pretty complex
// @see magento2/site/app/code/Magento/Catalog/Model/Product/Attribute/Backend/TierPrice.php
// @see magento2/site/app/code/Magento/Catalog/Model/Product/Attribute/Backend/GroupPrice/AbstractGroupPrice.php
func ProductBackendTierPrice() *todoPABTP {
	return &todoPABTP{
		AttributeBackend: eav.NewAttributeBackend(),
	}
}

// IsScalar returns false because a tier price is not a scalar type
func (gp *todoPABTP) IsScalar() bool {
	return false
}

// ProductBackendMedia @todo ... pretty complex
// @see magento2/site/app/code/Magento/Catalog/Model/Product/Attribute/Backend/Media.php
func ProductBackendMedia() *todoPABM {
	return &todoPABM{
		AttributeBackend: eav.NewAttributeBackend(),
	}
}

// IsScalar returns false because a media is not a scalar type
func (gp *todoPABM) IsScalar() bool {
	return false
}

// ProductRecurring @todo
// @see Mage_Catalog_Model_Product_Attribute_Backend_Recurring
func ProductBackendRecurring() *todoPABR {
	return &todoPABR{
		AttributeBackend: eav.NewAttributeBackend(),
	}
}

// ProductBackendSku @todo
// @see magento2/site/app/code/Magento/Catalog/Model/Product/Attribute/Backend/Sku.php
func ProductBackendSku() *todoPABSku {
	return &todoPABSku{
		AttributeBackend: eav.NewAttributeBackend(),
	}
}

// ProductBackendStock @todo
// @see magento2/site/app/code/Magento/Catalog/Model/Product/Attribute/Backend/Stock.php
func ProductBackendStock() *todoPABStock {
	return &todoPABStock{
		AttributeBackend: eav.NewAttributeBackend(),
	}
}
