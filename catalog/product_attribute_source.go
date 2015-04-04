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
	_ eav.AttributeSourceModeller = (*todoPASCOM)(nil)
	_ eav.AttributeSourceModeller = (*todoPASDesignOptionsContainer)(nil)
	_ eav.AttributeSourceModeller = (*todoPASL)(nil)
	_ eav.AttributeSourceModeller = (*todoPASStatus)(nil)
	_ eav.AttributeSourceModeller = (*todoPASVisibility)(nil)
)

type (
	todoPASCOM struct {
		*eav.AttributeSource
	}
	todoPASDesignOptionsContainer struct {
		*eav.AttributeSource
	}
	todoPASL struct {
		*eav.AttributeSource
	}
	todoPASStatus struct {
		*eav.AttributeSource
	}
	todoPASVisibility struct {
		*eav.AttributeSource
	}
)

// ProductAttributeSourceCountryOfManufacture @todo hand out price for longest name ;-)
// @see magento2/site/app/code/Magento/Catalog/Model/Product/Attribute/Source/Countryofmanufacture.php
func ProductAttributeSourceCountryOfManufacture() *todoPASCOM {
	return &todoPASCOM{}
}

// ProductAttributeSourceDesignOptionsContainer @todo
// @see magento2/site/app/code/Magento/Catalog/Model/Entity/Product/Attribute/Design/Options/Container.php
// @see magento2/site/app/code/Magento/Catalog/etc/di.xml Line ~109
func ProductAttributeSourceDesignOptionsContainer() *todoPASDesignOptionsContainer {
	return &todoPASDesignOptionsContainer{}
}

// ProductAttributeSourceLayout @todo
// @see magento2/site/app/code/Magento/Catalog/Model/Product/Attribute/Source/Layout.php
func ProductAttributeSourceLayout() *todoPASL {
	return &todoPASL{}
}

// ProductAttributeSourceStatus @todo
// @see magento2/site/app/code/Magento/Catalog/Model/Product/Attribute/Source/Status.php
func ProductAttributeSourceStatus() *todoPASStatus {
	return &todoPASStatus{}
}

// ProductAttributeSourceVisibility @todo
// @see magento2/site/app/code/Magento/Catalog/Model/Product/Visibility.php
// This class is misplaced in Magento2 maybe they move it to the correct location ...
func ProductAttributeSourceVisibility() *todoPASVisibility {
	return &todoPASVisibility{}
}
