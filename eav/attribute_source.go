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

package eav

type (
	// AttributeSourceModeller interface implements the functions needed to retrieve data from a source model
	// @see magento2/site/app/code/Magento/Eav/Model/Entity/Attribute/Source/SourceInterface.php and
	// its abstract class plus other default implementations. Refers to tables eav_attribute_option and
	// eav_attribute_option_value OR other tables. @todo
	// tbd: Should return a []string where i%0 is the value and i%1 is the label.
	AttributeSourceModeller interface {
		GetAllOptions()
		GetOptionText()
	}
	// AttributeSource should implement all abstract ideas of
	// @see magento2/site/app/code/Magento/Eav/Model/Entity/Attribute/Source/AbstractSource.php
	// maybe extend also the interface
	AttributeSource struct {
		*Attribute
	}
)

var _ AttributeSourceModeller = (*AttributeSource)(nil)

func (as AttributeSource) GetAllOptions() {}
func (as AttributeSource) GetOptionText() {}
