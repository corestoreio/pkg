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

import (
	"errors"

	"github.com/corestoreio/csfw/storage/csdb"
	"github.com/corestoreio/csfw/storage/dbr"
	"github.com/juju/errgo"
)

const (
	EntityTypeDatetime ValueIndex = iota + 1
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
		// Base returns the base/prefix table name. E.g.: catalog_product_entity
		TableNameBase() string
		// Type returns for a type the table name. E.g.: catalog_product_entity_int
		TableNameValue(ValueIndex) string
	}

	// @todo all the returning interfaces are all crap. even that structure is rarely used in Magento

	// EntityTypeAttributeModeller defines an attribute model @todo
	EntityTypeAttributeModeller interface {
		// Creates a new attribute to the corresponding entity. @todo options?
		// The return type must embed eav.Attributer interface and of course its custom attribute interface
		New() interface{}
		Get(i AttributeIndex) (interface{}, error)
		MustGet(i AttributeIndex) interface{}
		GetByID(id int64) (interface{}, error)
		GetByCode(code string) (interface{}, error)
	}

	// EntityTypeAttributeCollectioner defines an attribute collection @todo
	// it returns a slice so use type assertion.
	EntityTypeAttributeCollectioner interface {
		Collection() interface{}
	}

	// EntityTypeAdditionalAttributeTabler implements methods for EAV table structures to retrieve attributes
	EntityTypeAdditionalAttributeTabler interface {
		TableAdditionalAttribute() (*csdb.Table, error)
		// TableEavWebsite gets the table, where website-dependent attribute parameters are stored in.
		// If an EAV model doesn't demand this functionality, let this function just return nil,nil
		TableEavWebsite() (*csdb.Table, error)
	}

	// EntityTypeIncrementModeller defines who to increment a number @todo
	EntityTypeIncrementModeller interface {
		TBD()
	}

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
		EntityAttributeCollection EntityTypeAttributeCollectioner
	}
)

var (
	ErrEntityTypeValueNotFound = errors.New("Unknown entity type value")
	// csEntityTypeCollection contains all entity types mapped to their Go types/interfaces
	csEntityTypeCollection CSEntityTypeSlice
)

// GetEntityTypeCollection to avoid leaking global variable. Maybe returning a copy?
func GetEntityTypeCollection() CSEntityTypeSlice {
	return csEntityTypeCollection
}

// GetEntityTypeByID returns an entity type by its id
func GetEntityTypeByID(id int64) (*CSEntityType, error) {
	return csEntityTypeCollection.GetByID(id)
}

// GetEntityTypeByCode returns an entity type by its code
func GetEntityTypeByCode(code string) (*CSEntityType, error) {
	return csEntityTypeCollection.GetByCode(code)
}

// SetEntityTypeCollection sets the collection. Panics if slice is empty.
func SetEntityTypeCollection(sc CSEntityTypeSlice) {
	if len(sc) == 0 {
		panic("CSEntityTypeSlice is empty")
	}
	csEntityTypeCollection = sc
}

func (et *TableEntityType) LoadByCode(dbrSess *dbr.Session, code string, cbs ...csdb.DbrSelectCb) error {
	s, err := TableCollection.Structure(TableIndexEntityType)
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
