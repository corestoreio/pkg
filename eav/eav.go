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

import (
	"errors"

	"github.com/corestoreio/csfw/storage/csdb"
)

const (
	EntityTypeDatetime ValueIndex = iota + 1
	EntityTypeDecimal
	EntityTypeInt
	EntityTypeText
	EntityTypeVarchar
)

var (
	ErrEntityTypeValueNotFound = errors.New("Unknown entity type value")
)

type (
	ValueIndex int

	// EntityTypeModeller defines an entity type model @todo
	EntityTypeModeller interface {
		TBD()
	}

	// EntityTypeTabler returns the table name
	EntityTypeTabler interface {
		// Base returns the base/prefix table name. E.g.: catalog_product_entity
		TableNameBase() string
		// Type returns for a type the table name. E.g.: catalog_product_entity_int
		TableNameValue(ValueIndex) string
	}

	// EntityTypeAttributeModeller defines an attribute model @todo
	EntityTypeAttributeModeller interface {
		TBD()
	}

	// EntityTypeAdditionalAttributeTabler implements methods for EAV table structures to retrieve attributes
	EntityTypeAdditionalAttributeTabler interface {
		TableAdditionalAttribute() (*csdb.TableStructure, error)
		// TableNameEavWebsite gets the table, where website-dependent attribute parameters are stored in.
		// If an EAV model doesn't demand this functionality, let this function just return an empty string
		TableEavWebsite() (*csdb.TableStructure, error)
	}

	// EntityTypeIncrementModeller defines who to increment a number @todo
	EntityTypeIncrementModeller interface {
		TBD()
	}

	// EntityAttributeCollectioner defines an attribute collection @todo
	EntityAttributeCollectioner interface {
		AttributeCollection()
	}
)
