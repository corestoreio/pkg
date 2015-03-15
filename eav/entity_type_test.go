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
	"testing"

	"github.com/corestoreio/csfw/storage/csdb"
	"github.com/gocraft/dbr"
	"github.com/stretchr/testify/assert"
)

var (
	csEntityTypeCollection = CSEntityTypeSlice{
		&CSEntityType{
			EntityTypeID:          3,
			EntityTypeCode:        "catalog_category",
			EntityModel:           nil,
			AttributeModel:        nil,
			EntityTable:           nil,
			ValueTablePrefix:      "",
			IsDataSharing:         true,
			DataSharingKey:        "default",
			DefaultAttributeSetID: 3,

			IncrementPerStore:         false,
			IncrementPadLength:        8,
			IncrementPadChar:          "0",
			AdditionalAttributeTable:  nil,
			EntityAttributeCollection: nil,
		},
	}
)

func TestEntityType(t *testing.T) {
	db := csdb.MustConnectTest()
	defer db.Close()
	dbrSess := dbr.NewConnection(db, nil).NewSession(nil)
	var et EntityType
	et.LoadByCode(
		dbrSess,
		"catalog_product",
		func(sb *dbr.SelectBuilder) *dbr.SelectBuilder {
			sb.OrderBy("entity_type_id")
			return sb
		},
	)

	assert.NotEmpty(t, et.EntityModel)
	assert.NotEmpty(t, et.AttributeModel.String)
	assert.True(t, et.EntityTypeID > 0)
	assert.True(t, et.IsRealEav())
}

func TestEntityTypeSliceGetByCode(t *testing.T) {
	db := csdb.MustConnectTest()
	defer db.Close()
	dbrSess := dbr.NewConnection(db, nil).NewSession(nil)

	s, err := GetTableStructure(TableEntityType)
	if err != nil {
		t.Error(err)
	}

	var entityTypeCollection EntityTypeSlice
	_, err = dbrSess.
		Select(s.Columns...).
		From(s.Name).
		LoadStructs(&entityTypeCollection)
	if err != nil {
		t.Error(err)
	}

	etc, err := entityTypeCollection.GetByCode("catalog_categories")
	assert.Nil(t, etc)
	assert.Error(t, err)

	etc, err = entityTypeCollection.GetByCode("catalog_category")
	assert.NotNil(t, etc)
	assert.NoError(t, err)
}

func TestCSEntityTypeSliceGetByCode(t *testing.T) {
	etc, err := csEntityTypeCollection.GetByCode("catalog_category")
	assert.NotNil(t, etc)
	assert.NoError(t, err)

	etc, err = csEntityTypeCollection.GetByCode("catalog_categories")
	assert.Nil(t, etc)
	assert.Error(t, err)
}
