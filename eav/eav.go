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
	"github.com/corestoreio/csfw/storage/csdb"
	"github.com/gocraft/dbr"
	"github.com/juju/errgo"
)

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
	// not added entity_attribute_collection

	AttributeBackendModeller interface {
		TBD()
	}
	AttributeFrontendModeller interface {
		TBD()
	}
	AttributeSourceModeller interface {
		TBD()
	}
)

func (et *EntityType) LoadByCode(dbrSess *dbr.Session, code string, cb ...csdb.DbrSessionCallback) error {

	// @todo add check entry in EntityTypeCollection

	s, err := GetTableStructure(TableEntityType)
	if err != nil {
		return errgo.Mask(err)
	}

	return errgo.Mask(
		dbrSess.
			Select(s.Columns...).
			From(s.Name).
			Where("entity_type_code = ?", code).
			LoadStruct(et),
	)
}

// IsRealEav checks if those types which have an attribute model and therefore are a real EAV.
// sales* tables are not real EAV tables as they are already flat tables.
func (et *EntityType) IsRealEav() bool {
	return et.EntityTypeId > 0 && et.AttributeModel.Valid == true && et.AttributeModel.String != ""
}
