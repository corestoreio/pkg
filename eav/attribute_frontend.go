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
	// AttributeFrontendModeller defines the attribute frontend model @todo
	// @see magento2/site/app/code/Magento/Eav/Model/Entity/Attribute/Frontend/AbstractFrontend.php
	AttributeFrontendModeller interface {
		InputRenderer() FrontendInputRendererIFace
		GetValue()
		GetInputType() string
		// @todo
	}
	// FrontendInputRendererIFace see table catalog_eav_attribute.frontend_input_renderer @todo
	// Stupid name :-( Fix later.
	FrontendInputRendererIFace interface {
		// TBD()
	}
	// AttributeFrontend implements abstract functions @todo
	// @see magento2/site/app/code/Magento/Eav/Model/Entity/Attribute/Backend/AbstractBackend.php
	AttributeFrontend struct {
		*Attribute
	}
)

var _ AttributeFrontendModeller = (*AttributeFrontend)(nil)

func (af *AttributeFrontend) InputRenderer() FrontendInputRendererIFace { return nil }
func (af *AttributeFrontend) GetValue()                                 {}
func (af *AttributeFrontend) GetInputType() string {
	return af.Attribute.FrontendInput()
}
