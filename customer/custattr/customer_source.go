// Copyright 2015-2016, Cyrill @ Schumacher.fm and the CoreStore contributors
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

var (
	_ eav.AttributeSourceModeller = (*todoASG)(nil)
)

type (
	// example when you want to add custom struct fields ...
	todoASG struct {
		*eav.AttributeSource
	}
)

// CustomerSourceGroup customer group handling @todo
// @see magento2/site/app/code/Magento/Customer/Model/Customer/Attribute/Source/Group.php
func CustomerSourceGroup() todoASG {
	return todoASG{
		AttributeSource: eav.NewAttributeSource(),
	}
}

// CustomerSourceStore handle store source @todo
// @see magento2/site/app/code/Magento/Customer/Model/Customer/Attribute/Source/Store.php
func CustomerSourceStore() *eav.AttributeSource {
	return eav.NewAttributeSource()
}

// CustomerSourceWebsite handle store source @todo
// @see magento2/site/app/code/Magento/Customer/Model/Customer/Attribute/Source/Website.php
func CustomerSourceWebsite() *eav.AttributeSource {
	return eav.NewAttributeSource()
}
