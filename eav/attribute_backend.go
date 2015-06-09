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
	// AttributeBackendModeller defines the attribute backend model @todo
	// @see magento2/site/app/code/Magento/Eav/Model/Entity/Attribute/Backend/BackendInterface.php
	AttributeBackendModeller interface {
		// GetTable @todo
		GetTable() string
		IsStatic() bool
		GetType() string
		// GetEntityIDField @todo
		GetEntityIDField() string
		//SetValueId(valueId int) @todo
		//GetValueId()@todo

		//AfterLoad($object); object must be an interface @todo
		//BeforeSave($object);
		//AfterSave($object);
		//BeforeDelete($object);
		//AfterDelete($object);

		// Validate @todo
		Validate( /*data object*/ ) bool

		//GetEntityValueId(entity *CSEntityType) hmmmm
		//SetEntityValueId(entity *CSEntityType, valueId int) hmmmm

		// IsScalar By default attribute value is considered scalar that can be stored in a generic way
		IsScalar() bool
		// Config to configure the current instance
		Config(...AttributeBackendConfig) AttributeBackendModeller
	}
	// AttributeBackend implements abstract functions @todo
	// @see magento2/site/app/code/Magento/Eav/Model/Entity/Attribute/Backend/AbstractBackend.php
	AttributeBackend struct {
		// a is the reference to the already created attribute during init() call in a package.
		// Do not modify the attribute here
		a *Attribute
		// idx references to the generated constant and therefore references to itself. mainly used in
		// backend|source|frontend|etc_model
		idx AttributeIndex
	}
	AttributeBackendConfig func(*AttributeBackend)
)

var _ AttributeBackendModeller = (*AttributeBackend)(nil)

// NewAttributeBackend creates a pointer to a new attribute source
func NewAttributeBackend(cfgs ...AttributeBackendConfig) *AttributeBackend {
	as := &AttributeBackend{
		a: nil,
	}
	as.Config(cfgs...)
	return as
}

// AttributeBackendIdx only used in generated code to set the current index in the attribute slice
func AttributeBackendIdx(i AttributeIndex) AttributeBackendConfig {
	return func(as *AttributeBackend) {
		as.idx = i
	}
}

// Config runs the configuration functions
func (ab *AttributeBackend) Config(configs ...AttributeBackendConfig) AttributeBackendModeller {
	for _, cfg := range configs {
		cfg(ab)
	}
	return ab
}

func (ab *AttributeBackend) IsStatic() bool           { return ab.a.IsStatic() }
func (ab *AttributeBackend) GetTable() string         { return "" }
func (ab *AttributeBackend) GetType() string          { return ab.a.BackendType() }
func (ab *AttributeBackend) GetEntityIDField() string { return "" }
func (ab *AttributeBackend) Validate() bool           { return true }
func (ab *AttributeBackend) IsScalar() bool           { return true }
