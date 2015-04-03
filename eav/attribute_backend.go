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
	// AttributeBackendModeller defines the attribute backend model @todo
	// @see magento2/site/app/code/Magento/Eav/Model/Entity/Attribute/Backend/BackendInterface.php
	AttributeBackendModeller interface {
		// GetTable @todo
		GetTable() string
		IsStatic() bool
		GetType() string
		// GetEntityIdField @todo
		GetEntityIdField() string
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
	}
	// AttributeBackend implements abstract functions @todo
	// @see magento2/site/app/code/Magento/Eav/Model/Entity/Attribute/Backend/AbstractBackend.php
	AttributeBackend struct {
		*Attribute
	}
)

var _ AttributeBackendModeller = (*AttributeBackend)(nil)

func (ab *AttributeBackend) IsStatic() bool           { return ab.Attribute.IsStatic() }
func (ab *AttributeBackend) GetTable() string         { return "" }
func (ab *AttributeBackend) GetType() string          { return ab.Attribute.BackendType() }
func (ab *AttributeBackend) GetEntityIdField() string { return "" }
func (ab *AttributeBackend) Validate() bool           { return true }
func (ab *AttributeBackend) IsScalar() bool           { return true }
