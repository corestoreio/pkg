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
	"github.com/corestoreio/csfw/storage/dbr"
	"github.com/juju/errgo"
)

type (
	// CSEntityTypeSlice Types starting with CS are the CoreStore mappings with the DB data
	CSEntityTypeSlice []*CSEntityType
	// CSEntityType Go Type of the Mage database models and types. The prefix CS indicates
	// that this EntityType is not generated because it contains special interfaces.
	CSEntityType struct {
		EntityTypeID              int64
		EntityTypeCode            string
		EntityModel               EntityTypeModeller
		AttributeModel            EntityTypeAttributeModeller
		EntityTable               EntityTypeTabler
		ValueTablePrefix          string
		EntityIDField             string
		IsDataSharing             bool
		DataSharingKey            string
		DefaultAttributeSetID     int64
		IncrementModel            EntityTypeIncrementModeller
		IncrementPerStore         bool
		IncrementPadLength        int64
		IncrementPadChar          string
		AdditionalAttributeTable  EntityTypeAdditionalAttributeTabler
		EntityAttributeCollection EntityAttributeCollectioner
	}
)

// csEntityTypeCollection contains all entity types mapped to their Go types/interfaces
var csEntityTypeCollection CSEntityTypeSlice

// GetEntityTypeCollection to avoid leaking global variable. Maybe returning a copy?
func GetEntityTypeCollection() CSEntityTypeSlice {
	return csEntityTypeCollection
}

// SetEntityTypeCollection sets the collection. Panics if slice is empty.
func SetEntityTypeCollection(sc CSEntityTypeSlice) {
	if len(sc) == 0 {
		panic("CSEntityTypeSlice is empty")
	}
	csEntityTypeCollection = sc
}

func (et *TableEntityType) LoadByCode(dbrSess *dbr.Session, code string, cbs ...csdb.DbrSelectCb) error {
	s, err := GetTableStructure(TableIndexEntityType)
	if err != nil {
		return errgo.Mask(err)
	}
	sb := dbrSess.Select(s.AllColumnAliasQuote(csdb.MainTable)...).From(s.Name, csdb.MainTable).Where("entity_type_code = ?", code)
	for _, cb := range cbs {
		sb = cb(sb)
	}
	return errgo.Mask(sb.LoadStruct(et))
}

// IsRealEav checks if those types which have an attribute model and therefore are a real EAV.
// sales* tables are not real EAV tables as they are already flat tables.
func (et *TableEntityType) IsRealEav() bool {
	return et.EntityTypeID > 0 && et.AttributeModel.Valid == true && et.AttributeModel.String != ""
}

// GetByCode returns a TableEntityType using the entity code
func (es TableEntityTypeSlice) GetByCode(code string) (*TableEntityType, error) {
	for _, e := range es {
		if e.EntityTypeCode == code {
			return e, nil
		}
	}
	return nil, errgo.Newf("Entity Code %s not found", code)
}

// GetByCode returns a CSEntityType using the entity code
func (es CSEntityTypeSlice) GetByCode(code string) (*CSEntityType, error) {
	for _, e := range es {
		if e.EntityTypeCode == code {
			return e, nil
		}
	}
	return nil, errgo.Newf("Entity Code %s not found", code)
}

// GetByID returns a CSEntityType using the entity id
func (es CSEntityTypeSlice) GetByID(id int64) (*CSEntityType, error) {
	for _, e := range es {
		if e.EntityTypeID == id {
			return e, nil
		}
	}
	return nil, errgo.Newf("Entity ID %d not found", id)
}

// EntityTablePrefix eav table name prefix
// @see magento2/site/app/code/Magento/Eav/Model/Entity/AbstractEntity.php::getEntityTablePrefix()
func (e *CSEntityType) GetEntityTablePrefix() string {
	return e.GetValueTablePrefix()
}

// ValueTablePrefix returns the table prefix for all value tables
// @see magento2/site/app/code/Magento/Eav/Model/Entity/AbstractEntity.php::getValueTablePrefix()
func (e *CSEntityType) GetValueTablePrefix() string {
	if e.ValueTablePrefix == "" {
		return e.EntityTable.TableNameBase()
	}
	return e.ValueTablePrefix
}
