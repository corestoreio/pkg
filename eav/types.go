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
	// Defines EntityTypeModel @todo this guy handles everything with the model
	EntityTypeModeller interface {
		TBD()
	}

	// used in EntityType to map Mage1+2 data to Go packages
	// Interface name is the name of the column +er
	EntityTypeTabler interface {
		TableName() string
	}

	// @todo this dude handles everything related to attributes
	EntityTypeAttributeModeller interface {
		TBD()
	}

	// @todo
	EntityTypeAdditionalAttributeTabler interface {
		TableName() string
	}

	// @todo How increment a number e.g. customer, invoice, order ...
	EntityTypeIncrementModeller interface {
		TBD()
	}

	EntityAttributeCollectioner interface {
		TBD()
	}

	AttributeBackendModeller interface {
		TBD()
	}
	AttributeFrontendModeller interface {
		TBD()
	}
	AttributeSourceModeller interface {
		TBD()
	}

	// Types starting with CS are the CoreStore mappings with the DB data
	CSEntityTypeSlice []*CSEntityType
	CSEntityType      struct {
		EntityTypeId              int64
		EntityTypeCode            string
		EntityModel               EntityTypeModeller
		AttributeModel            EntityTypeAttributeModeller
		EntityTable               EntityTypeTabler
		ValueTablePrefix          string
		EntityIdField             string
		IsDataSharing             bool
		DataSharingKey            string
		DefaultAttributeSetId     int64
		IncrementModel            EntityTypeIncrementModeller
		IncrementPerStore         int64
		IncrementPadLength        int64
		IncrementPadChar          string
		AdditionalAttributeTable  EntityTypeAdditionalAttributeTabler
		EntityAttributeCollection EntityAttributeCollectioner
	}
)
