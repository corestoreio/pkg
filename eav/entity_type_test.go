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

package eav_test

import (
	"testing"

	"github.com/corestoreio/csfw/eav"
	"github.com/corestoreio/csfw/storage/csdb"
	"github.com/corestoreio/csfw/storage/dbr"
	"github.com/stretchr/testify/assert"
)

var (
	csEntityTypeCollection = eav.CSEntityTypeSlice{
		&eav.CSEntityType{
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

func init() {
	dbc := csdb.MustConnectTest()
	defer dbc.Close()
	if err := eav.TableCollection.Init(dbc.NewSession()); err != nil {
		panic(err)
	}
}

func TestEntityType(t *testing.T) {
	dbc := csdb.MustConnectTest()
	defer dbc.Close()
	dbrSess := dbc.NewSession()

	var et eav.TableEntityType
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
	assert.True(t, et.EntityTypeID > 0, "EntityTypeID should be greater 0 but is: %#v\n", et)
	assert.True(t, et.IsRealEav())
}

func TestEntityTypeSliceGetByCode(t *testing.T) {
	dbc := csdb.MustConnectTest()
	defer dbc.Close()
	dbrSess := dbc.NewSession()

	s, err := eav.TableCollection.Structure(eav.TableIndexEntityType)
	if err != nil {
		t.Error(err)
	}

	var entityTypeCollection eav.TableEntityTypeSlice
	_, err = dbrSess.
		Select(s.Columns.FieldNames()...).
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
