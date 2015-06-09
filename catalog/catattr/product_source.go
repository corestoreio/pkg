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

// ProductSourceCountryOfManufacture @todo hand out price for longest name ;-)
// @see magento2/site/app/code/Magento/Catalog/Model/Product/Attribute/Source/Countryofmanufacture.php
func ProductSourceCountryOfManufacture() *eav.AttributeSource {
	return eav.NewAttributeSource()
}

// ProductSourceDesignOptionsContainer @todo
// @see magento2/site/app/code/Magento/Catalog/Model/Entity/Product/Attribute/Design/Options/Container.php
// @see magento2/site/app/code/Magento/Catalog/etc/di.xml Line ~109
func ProductSourceDesignOptionsContainer() *eav.AttributeSource {
	return eav.NewAttributeSource()
}

// ProductSourceLayout @todo
// @see magento2/site/app/code/Magento/Catalog/Model/Product/Attribute/Source/Layout.php
func ProductSourceLayout() *eav.AttributeSource {
	return eav.NewAttributeSource()
}

// ProductSourceStatus @todo
// @see magento2/site/app/code/Magento/Catalog/Model/Product/Attribute/Source/Status.php
func ProductSourceStatus() *eav.AttributeSource {
	return eav.NewAttributeSource()
}

// ProductSourceVisibility @todo
// @see magento2/site/app/code/Magento/Catalog/Model/Product/Visibility.php
// This class is misplaced in Magento2 maybe they move it to the correct location ...
func ProductSourceVisibility() *eav.AttributeSource {
	return eav.NewAttributeSource()
}
