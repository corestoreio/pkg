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

package eav

type (
	// AttributeFrontendModeller defines the attribute frontend model @todo
	// @see magento2/site/app/code/Magento/Eav/Model/Entity/Attribute/Frontend/AbstractFrontend.php
	AttributeFrontendModeller interface {
		InputRenderer() FrontendInputRendererIFace
		GetValue()
		GetInputType() string
		// @todo

		// Config to configure the current instance
		Config(...AttributeFrontendConfig) AttributeFrontendModeller
	}
	// FrontendInputRendererIFace see table catalog_eav_attribute.frontend_input_renderer @todo
	// Stupid name :-( Fix later.
	FrontendInputRendererIFace interface {
		// TBD()
	}
	// AttributeFrontend implements abstract functions @todo
	// @see magento2/site/app/code/Magento/Eav/Model/Entity/Attribute/Backend/AbstractBackend.php
	AttributeFrontend struct {
		// a is the reference to the already created attribute during init() call in a package.
		// Do not modify a here
		a *Attribute
		// idx references to the generated constant and therefore references to itself. mainly used in
		// backend|source|frontend|etc_model
		idx AttributeIndex
	}
	AttributeFrontendConfig func(*AttributeFrontend)
)

var _ AttributeFrontendModeller = (*AttributeFrontend)(nil)

// NewAttributeFrontend creates a pointer to a new attribute source
func NewAttributeFrontend(cfgs ...AttributeFrontendConfig) *AttributeFrontend {
	as := &AttributeFrontend{
		a: nil,
	}
	as.Config(cfgs...)
	return as
}

// AttributeFrontendIdx only used in generated code to set the current index in the attribute slice
func AttributeFrontendIdx(i AttributeIndex) AttributeFrontendConfig {
	return func(as *AttributeFrontend) {
		as.idx = i
	}
}

// Config runs the configuration functions
func (af *AttributeFrontend) Config(configs ...AttributeFrontendConfig) AttributeFrontendModeller {
	for _, cfg := range configs {
		cfg(af)
	}
	return af
}

func (af *AttributeFrontend) InputRenderer() FrontendInputRendererIFace { return nil }
func (af *AttributeFrontend) GetValue()                                 {}
func (af *AttributeFrontend) GetInputType() string {
	return af.a.FrontendInput()
}
