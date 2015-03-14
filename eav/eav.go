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

const (
	EntityTypeDateTime ValueIndex = iota + 1
	EntityTypeDecimal
	EntityTypeInt
	EntityTypeText
	EntityTypeVarchar
)

type (
	ValueIndex int

	// EntityTypeModeller defines an entity type model @todo
	EntityTypeModeller interface {
		TBD()
	}

	// EntityTypeTabler returns the table name
	EntityTypeTabler interface {
		// Base returns the base table name. E.g.: catalog_product_entity
		TableNameBase() string
		// Type returns for a type the table name. E.g.: catalog_product_entity_int
		TableNameValue(ValueIndex) string
	}

	// EntityTypeAttributeModeller defines an attribute model @todo
	EntityTypeAttributeModeller interface {
		TBD()
	}

	// EntityTypeAdditionalAttributeTabler returns the table name
	EntityTypeAdditionalAttributeTabler interface {
		TableNameAdditionalAttribute() string
	}

	// EntityTypeIncrementModeller defines who to increment a number @todo
	EntityTypeIncrementModeller interface {
		TBD()
	}

	// EntityAttributeCollectioner defines an attribute collection @todo
	EntityAttributeCollectioner interface {
		AttributeCollection()
	}

	// AttributeBackendModeller defines the attribute backend model @todo
	AttributeBackendModeller interface {
		TBD()
	}

	// AttributeFrontendModeller defines the attribute frontend model @todo
	AttributeFrontendModeller interface {
		TBD()
	}

	// AttributeSourceModeller defines the source where an attribute can also be stored @todo
	AttributeSourceModeller interface {
		TBD()
	}
)
