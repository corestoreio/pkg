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

// auto generated via eavToStruct
package eav

import (
	"database/sql"

	"github.com/gocraft/dbr"
)

var (
	EntityTypeCollection = EntityTypeSlice{
		&EntityType{
			EntityTypeId:   1,
			EntityTypeCode: "customer",
			EntityModel:    "customer/customer",
			AttributeModel: dbr.NullString{
				NullString: sql.NullString{String: "customer/attribute", Valid: true},
			},
			EntityTable: dbr.NullString{
				NullString: sql.NullString{String: "customer/entity", Valid: true},
			},
			ValueTablePrefix: dbr.NullString{},
			EntityIdField:    dbr.NullString{},
			IsDataSharing:    true,
			DataSharingKey: dbr.NullString{
				NullString: sql.NullString{String: "default", Valid: true},
			},
			DefaultAttributeSetId: 1,
			IncrementModel: dbr.NullString{
				NullString: sql.NullString{String: "eav/entity_increment_numeric", Valid: true},
			},
			IncrementPerStore:  0,
			IncrementPadLength: 8,
			IncrementPadChar:   "0",
			AdditionalAttributeTable: dbr.NullString{
				NullString: sql.NullString{String: "customer/eav_attribute", Valid: true},
			},
			EntityAttributeCollection: dbr.NullString{
				NullString: sql.NullString{String: "customer/attribute_collection", Valid: true},
			},
		},
		&EntityType{
			EntityTypeId:   2,
			EntityTypeCode: "customer_address",
			EntityModel:    "customer/address",
			AttributeModel: dbr.NullString{
				NullString: sql.NullString{String: "customer/attribute", Valid: true},
			},
			EntityTable: dbr.NullString{
				NullString: sql.NullString{String: "customer/address_entity", Valid: true},
			},
			ValueTablePrefix: dbr.NullString{},
			EntityIdField:    dbr.NullString{},
			IsDataSharing:    true,
			DataSharingKey: dbr.NullString{
				NullString: sql.NullString{String: "default", Valid: true},
			},
			DefaultAttributeSetId: 2,
			IncrementModel:        dbr.NullString{},
			IncrementPerStore:     0,
			IncrementPadLength:    8,
			IncrementPadChar:      "0",
			AdditionalAttributeTable: dbr.NullString{
				NullString: sql.NullString{String: "customer/eav_attribute", Valid: true},
			},
			EntityAttributeCollection: dbr.NullString{
				NullString: sql.NullString{String: "customer/address_attribute_collection", Valid: true},
			},
		},
		&EntityType{
			EntityTypeId:   3,
			EntityTypeCode: "catalog_category",
			EntityModel:    "catalog/category",
			AttributeModel: dbr.NullString{
				NullString: sql.NullString{String: "catalog/resource_eav_attribute", Valid: true},
			},
			EntityTable: dbr.NullString{
				NullString: sql.NullString{String: "catalog/category", Valid: true},
			},
			ValueTablePrefix: dbr.NullString{},
			EntityIdField:    dbr.NullString{},
			IsDataSharing:    true,
			DataSharingKey: dbr.NullString{
				NullString: sql.NullString{String: "default", Valid: true},
			},
			DefaultAttributeSetId: 3,
			IncrementModel:        dbr.NullString{},
			IncrementPerStore:     0,
			IncrementPadLength:    8,
			IncrementPadChar:      "0",
			AdditionalAttributeTable: dbr.NullString{
				NullString: sql.NullString{String: "catalog/eav_attribute", Valid: true},
			},
			EntityAttributeCollection: dbr.NullString{
				NullString: sql.NullString{String: "catalog/category_attribute_collection", Valid: true},
			},
		},
		&EntityType{
			EntityTypeId:   4,
			EntityTypeCode: "catalog_product",
			EntityModel:    "catalog/product",
			AttributeModel: dbr.NullString{
				NullString: sql.NullString{String: "catalog/resource_eav_attribute", Valid: true},
			},
			EntityTable: dbr.NullString{
				NullString: sql.NullString{String: "catalog/product", Valid: true},
			},
			ValueTablePrefix: dbr.NullString{},
			EntityIdField:    dbr.NullString{},
			IsDataSharing:    true,
			DataSharingKey: dbr.NullString{
				NullString: sql.NullString{String: "default", Valid: true},
			},
			DefaultAttributeSetId: 4,
			IncrementModel:        dbr.NullString{},
			IncrementPerStore:     0,
			IncrementPadLength:    8,
			IncrementPadChar:      "0",
			AdditionalAttributeTable: dbr.NullString{
				NullString: sql.NullString{String: "catalog/eav_attribute", Valid: true},
			},
			EntityAttributeCollection: dbr.NullString{
				NullString: sql.NullString{String: "catalog/product_attribute_collection", Valid: true},
			},
		},
		&EntityType{
			EntityTypeId:   5,
			EntityTypeCode: "order",
			EntityModel:    "sales/order",
			AttributeModel: dbr.NullString{},
			EntityTable: dbr.NullString{
				NullString: sql.NullString{String: "sales/order", Valid: true},
			},
			ValueTablePrefix: dbr.NullString{},
			EntityIdField:    dbr.NullString{},
			IsDataSharing:    true,
			DataSharingKey: dbr.NullString{
				NullString: sql.NullString{String: "default", Valid: true},
			},
			DefaultAttributeSetId: 0,
			IncrementModel: dbr.NullString{
				NullString: sql.NullString{String: "eav/entity_increment_numeric", Valid: true},
			},
			IncrementPerStore:         1,
			IncrementPadLength:        8,
			IncrementPadChar:          "0",
			AdditionalAttributeTable:  dbr.NullString{},
			EntityAttributeCollection: dbr.NullString{},
		},
		&EntityType{
			EntityTypeId:   6,
			EntityTypeCode: "invoice",
			EntityModel:    "sales/order_invoice",
			AttributeModel: dbr.NullString{},
			EntityTable: dbr.NullString{
				NullString: sql.NullString{String: "sales/invoice", Valid: true},
			},
			ValueTablePrefix: dbr.NullString{},
			EntityIdField:    dbr.NullString{},
			IsDataSharing:    true,
			DataSharingKey: dbr.NullString{
				NullString: sql.NullString{String: "default", Valid: true},
			},
			DefaultAttributeSetId: 0,
			IncrementModel: dbr.NullString{
				NullString: sql.NullString{String: "eav/entity_increment_numeric", Valid: true},
			},
			IncrementPerStore:         1,
			IncrementPadLength:        8,
			IncrementPadChar:          "0",
			AdditionalAttributeTable:  dbr.NullString{},
			EntityAttributeCollection: dbr.NullString{},
		},
		&EntityType{
			EntityTypeId:   7,
			EntityTypeCode: "creditmemo",
			EntityModel:    "sales/order_creditmemo",
			AttributeModel: dbr.NullString{},
			EntityTable: dbr.NullString{
				NullString: sql.NullString{String: "sales/creditmemo", Valid: true},
			},
			ValueTablePrefix: dbr.NullString{},
			EntityIdField:    dbr.NullString{},
			IsDataSharing:    true,
			DataSharingKey: dbr.NullString{
				NullString: sql.NullString{String: "default", Valid: true},
			},
			DefaultAttributeSetId: 0,
			IncrementModel: dbr.NullString{
				NullString: sql.NullString{String: "eav/entity_increment_numeric", Valid: true},
			},
			IncrementPerStore:         1,
			IncrementPadLength:        8,
			IncrementPadChar:          "0",
			AdditionalAttributeTable:  dbr.NullString{},
			EntityAttributeCollection: dbr.NullString{},
		},
		&EntityType{
			EntityTypeId:   8,
			EntityTypeCode: "shipment",
			EntityModel:    "sales/order_shipment",
			AttributeModel: dbr.NullString{},
			EntityTable: dbr.NullString{
				NullString: sql.NullString{String: "sales/shipment", Valid: true},
			},
			ValueTablePrefix: dbr.NullString{},
			EntityIdField:    dbr.NullString{},
			IsDataSharing:    true,
			DataSharingKey: dbr.NullString{
				NullString: sql.NullString{String: "default", Valid: true},
			},
			DefaultAttributeSetId: 0,
			IncrementModel: dbr.NullString{
				NullString: sql.NullString{String: "eav/entity_increment_numeric", Valid: true},
			},
			IncrementPerStore:         1,
			IncrementPadLength:        8,
			IncrementPadChar:          "0",
			AdditionalAttributeTable:  dbr.NullString{},
			EntityAttributeCollection: dbr.NullString{},
		},
	}
)
