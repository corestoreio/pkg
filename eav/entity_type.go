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
	// CSEntityTypeSlice Types starting with CS are the CoreStore mappings with the DB data
	CSEntityTypeSlice []*CSEntityType
	// CSEntityType Go Type of the Mage database models and types
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

func (et *EavEntityType) LoadByCode(dbrSess *dbr.Session, code string, cbs ...csdb.DbrSelectCb) error {
	s, err := GetTableStructure(TableEavEntityType)
	if err != nil {
		return errgo.Mask(err)
	}
	qry := dbrSess.Select(s.Columns...).From(s.Name).Where("entity_type_code = ?", code)
	for _, cb := range cbs {
		qry = cb(qry)
	}
	return errgo.Mask(qry.LoadStruct(et))
}

// IsRealEav checks if those types which have an attribute model and therefore are a real EAV.
// sales* tables are not real EAV tables as they are already flat tables.
func (et *EavEntityType) IsRealEav() bool {
	return et.EntityTypeID > 0 && et.AttributeModel.Valid == true && et.AttributeModel.String != ""
}

func (es EavEntityTypeSlice) GetByCode(code string) (*EavEntityType, error) {
	for _, e := range es {
		if e.EntityTypeCode == code {
			return e, nil
		}
	}
	return nil, errgo.Newf("Entity Code %s not found", code)
}

func (es CSEntityTypeSlice) GetByCode(code string) (*CSEntityType, error) {
	for _, e := range es {
		if e.EntityTypeCode == code {
			return e, nil
		}
	}
	return nil, errgo.Newf("Entity Code %s not found", code)
}
