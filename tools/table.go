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

package tools

import (
	"database/sql"
	"strings"

	"github.com/corestoreio/csfw/storage/dbr"
	"github.com/juju/errgo"
)

const (
	TableNameSeparator string = "_"
	TableEavEntityType string = "eav_entity_type"
)

var (
	// TableEntityTypeSuffix e.g. for catalog_product_entity, customer_entity
	TableEntityTypeSuffix = "entity"
	// TableEntityTypeValueSuffixes defines all possible value type tables which an EAV model can have.
	TableEntityTypeValueSuffixes = ValueSuffixes{
		"datetime",
		"decimal",
		"int",
		"text",
		"varchar",
	}
)

type (
	ValueSuffixes      []string
	TypeCodeValueTable map[string]map[string]string // 1. key entity_type_code 2. key table name => value ValueSuffix

	EntityTypeMap struct {
		ImportPath                string `json:"import_path"`
		EntityTypeID              int64  `db:"entity_type_id"`
		EntityTypeCode            string `db:"entity_type_code"`
		EntityModel               string `json:"entity_model"`
		AttributeModel            string `json:"attribute_model"`
		EntityTable               string `json:"entity_table"`
		ValueTablePrefix          string `db:"value_table_prefix"`
		EntityIDField             string
		IsDataSharing             bool
		DataSharingKey            string
		DefaultAttributeSetID     int64
		IncrementModel            string `json:"increment_model"`
		IncrementPerStore         bool
		IncrementPadLength        int64
		IncrementPadChar          string
		AdditionalAttributeTable  string `json:"additional_attribute_table"`
		EntityAttributeCollection string `json:"entity_attribute_collection"`
	}
)

func (vs ValueSuffixes) contains(suffix string) bool {
	for _, v := range vs {
		if v == suffix {
			return true
		}
	}
	return false
}

func (vs ValueSuffixes) String() string {
	return strings.Join(vs, ", ")
}

func (m TypeCodeValueTable) Empty() bool {
	_, ok := m[""]
	return len(m) < 1 || ok
}

func GetTables(db *sql.DB, prefix string) ([]string, error) {

	var tableNames = make([]string, 0, 200)
	qry := "SHOW TABLES like '" + prefix + "%'"

	rows, err := db.Query(qry)
	if err != nil {
		return nil, errgo.Mask(err)
	}
	defer rows.Close()

	for rows.Next() {
		var tableName string
		err := rows.Scan(&tableName)
		if err != nil {
			return nil, errgo.Mask(err)
		}
		tableNames = append(tableNames, tableName)
	}
	err = rows.Err()
	if err != nil {
		return nil, errgo.Mask(err)
	}
	return tableNames, nil
}

// GetEavValueTables returns a map of all custom and default value tables for entity type codes.
// Despite value_table_prefix can have in Magento a different table name we treat it here
// as the table name itself. Not thread safe.
func GetEavValueTables(dbrConn *dbr.Connection, prefix string, entityTypeCodes []string) (TypeCodeValueTable, error) {

	typeCodeTables := make(TypeCodeValueTable, len(entityTypeCodes))

	for _, typeCode := range entityTypeCodes {

		vtp, err := dbrConn.NewSession(nil).
			Select("`value_table_prefix`").
			From(prefix+TableEavEntityType).
			Where("`value_table_prefix` IS NOT NULL").
			Where("`entity_type_code` = ?", typeCode).
			ReturnString()

		if err != nil && err != dbr.ErrNotFound {
			return nil, errgo.Mask(err)
		}
		if vtp == "" {
			vtp = typeCode + TableNameSeparator + TableEntityTypeSuffix + TableNameSeparator // e.g. catalog_product_entity_
		} else {
			vtp = vtp + TableNameSeparator
		}

		tableNames, err := GetTables(dbrConn.Db, prefix+vtp)
		if err != nil {
			return nil, errgo.Mask(err)
		}

		if _, ok := typeCodeTables[typeCode]; !ok {
			typeCodeTables[typeCode] = make(map[string]string, len(tableNames))
		}
		for _, t := range tableNames {
			valueSuffix := t[len(prefix+vtp):]
			if TableEntityTypeValueSuffixes.contains(valueSuffix) {
				/*
				   other tables like catalog_product_entity_gallery, catalog_product_entity_group_price,
				   catalog_product_entity_tier_price, etc are the backend model tables for different storage systems.
				   they are not part of the default EAV model.
				*/
				typeCodeTables[typeCode][t] = valueSuffix
			}

		}

	}

	return typeCodeTables, nil
}
