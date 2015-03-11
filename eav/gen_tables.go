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

// auto generated via tableToStruct
package eav

import (
	"time"

	"github.com/corestoreio/csfw/storage/csdb"
	"github.com/gocraft/dbr"
)

const (
	TableNoop csdb.Index = iota
	TableAttribute
	TableAttributeGroup
	TableAttributeLabel
	TableAttributeOption
	TableAttributeOptionValue
	TableAttributeSet
	TableEntity
	TableEntityAttribute
	TableEntityDatetime
	TableEntityDecimal
	TableEntityInt
	TableEntityStore
	TableEntityText
	TableEntityType
	TableEntityVarchar
	TableFormElement
	TableFormFieldset
	TableFormFieldsetLabel
	TableFormType
	TableFormTypeEntity
	TableMax
)

var (
	// read only map
	tableMap = csdb.TableMap{
		TableAttribute: csdb.NewTableStructure(
			"eav_attribute",
			[]string{
				"attribute_id",
			},
			[]string{
				"attribute_id",
				"entity_type_id",
				"attribute_code",
				"attribute_model",
				"backend_model",
				"backend_type",
				"backend_table",
				"frontend_model",
				"frontend_input",
				"frontend_label",
				"frontend_class",
				"source_model",
				"is_required",
				"is_user_defined",
				"default_value",
				"is_unique",
				"note",
			},
		),
		TableAttributeGroup: csdb.NewTableStructure(
			"eav_attribute_group",
			[]string{
				"attribute_group_id",
			},
			[]string{
				"attribute_group_id",
				"attribute_set_id",
				"attribute_group_name",
				"sort_order",
				"default_id",
			},
		),
		TableAttributeLabel: csdb.NewTableStructure(
			"eav_attribute_label",
			[]string{
				"attribute_label_id",
			},
			[]string{
				"attribute_label_id",
				"attribute_id",
				"store_id",
				"value",
			},
		),
		TableAttributeOption: csdb.NewTableStructure(
			"eav_attribute_option",
			[]string{
				"option_id",
			},
			[]string{
				"option_id",
				"attribute_id",
				"sort_order",
			},
		),
		TableAttributeOptionValue: csdb.NewTableStructure(
			"eav_attribute_option_value",
			[]string{
				"value_id",
			},
			[]string{
				"value_id",
				"option_id",
				"store_id",
				"value",
			},
		),
		TableAttributeSet: csdb.NewTableStructure(
			"eav_attribute_set",
			[]string{
				"attribute_set_id",
			},
			[]string{
				"attribute_set_id",
				"entity_type_id",
				"attribute_set_name",
				"sort_order",
			},
		),
		TableEntity: csdb.NewTableStructure(
			"eav_entity",
			[]string{
				"entity_id",
			},
			[]string{
				"entity_id",
				"entity_type_id",
				"attribute_set_id",
				"increment_id",
				"parent_id",
				"store_id",
				"created_at",
				"updated_at",
				"is_active",
			},
		),
		TableEntityAttribute: csdb.NewTableStructure(
			"eav_entity_attribute",
			[]string{
				"entity_attribute_id",
			},
			[]string{
				"entity_attribute_id",
				"entity_type_id",
				"attribute_set_id",
				"attribute_group_id",
				"attribute_id",
				"sort_order",
			},
		),
		TableEntityDatetime: csdb.NewTableStructure(
			"eav_entity_datetime",
			[]string{
				"value_id",
			},
			[]string{
				"value_id",
				"entity_type_id",
				"attribute_id",
				"store_id",
				"entity_id",
				"value",
			},
		),
		TableEntityDecimal: csdb.NewTableStructure(
			"eav_entity_decimal",
			[]string{
				"value_id",
			},
			[]string{
				"value_id",
				"entity_type_id",
				"attribute_id",
				"store_id",
				"entity_id",
				"value",
			},
		),
		TableEntityInt: csdb.NewTableStructure(
			"eav_entity_int",
			[]string{
				"value_id",
			},
			[]string{
				"value_id",
				"entity_type_id",
				"attribute_id",
				"store_id",
				"entity_id",
				"value",
			},
		),
		TableEntityStore: csdb.NewTableStructure(
			"eav_entity_store",
			[]string{
				"entity_store_id",
			},
			[]string{
				"entity_store_id",
				"entity_type_id",
				"store_id",
				"increment_prefix",
				"increment_last_id",
			},
		),
		TableEntityText: csdb.NewTableStructure(
			"eav_entity_text",
			[]string{
				"value_id",
			},
			[]string{
				"value_id",
				"entity_type_id",
				"attribute_id",
				"store_id",
				"entity_id",
				"value",
			},
		),
		TableEntityType: csdb.NewTableStructure(
			"eav_entity_type",
			[]string{
				"entity_type_id",
			},
			[]string{
				"entity_type_id",
				"entity_type_code",
				"entity_model",
				"attribute_model",
				"entity_table",
				"value_table_prefix",
				"entity_id_field",
				"is_data_sharing",
				"data_sharing_key",
				"default_attribute_set_id",
				"increment_model",
				"increment_per_store",
				"increment_pad_length",
				"increment_pad_char",
				"additional_attribute_table",
				"entity_attribute_collection",
			},
		),
		TableEntityVarchar: csdb.NewTableStructure(
			"eav_entity_varchar",
			[]string{
				"value_id",
			},
			[]string{
				"value_id",
				"entity_type_id",
				"attribute_id",
				"store_id",
				"entity_id",
				"value",
			},
		),
		TableFormElement: csdb.NewTableStructure(
			"eav_form_element",
			[]string{
				"element_id",
			},
			[]string{
				"element_id",
				"type_id",
				"fieldset_id",
				"attribute_id",
				"sort_order",
			},
		),
		TableFormFieldset: csdb.NewTableStructure(
			"eav_form_fieldset",
			[]string{
				"fieldset_id",
			},
			[]string{
				"fieldset_id",
				"type_id",
				"code",
				"sort_order",
			},
		),
		TableFormFieldsetLabel: csdb.NewTableStructure(
			"eav_form_fieldset_label",
			[]string{
				"fieldset_id",
				"store_id",
			},
			[]string{
				"fieldset_id",
				"store_id",
				"label",
			},
		),
		TableFormType: csdb.NewTableStructure(
			"eav_form_type",
			[]string{
				"type_id",
			},
			[]string{
				"type_id",
				"code",
				"label",
				"is_system",
				"theme",
				"store_id",
			},
		),
		TableFormTypeEntity: csdb.NewTableStructure(
			"eav_form_type_entity",
			[]string{
				"type_id",
				"entity_type_id",
			},
			[]string{
				"type_id",
				"entity_type_id",
			},
		),
	}
)

func GetTableStructure(i csdb.Index) (*csdb.TableStructure, error) {
	return tableMap.Structure(i)
}

func GetTableName(i csdb.Index) string {
	return tableMap.Name(i)
}

type (
	AttributeSlice []*Attribute
	Attribute      struct {
		AttributeId    int64          `db:"attribute_id"`    // attribute_id smallint(5) unsigned NOT NULL PRI  auto_increment
		EntityTypeId   int64          `db:"entity_type_id"`  // entity_type_id smallint(5) unsigned NOT NULL MUL DEFAULT '0'
		AttributeCode  dbr.NullString `db:"attribute_code"`  // attribute_code varchar(255) NULL
		AttributeModel dbr.NullString `db:"attribute_model"` // attribute_model varchar(255) NULL
		BackendModel   dbr.NullString `db:"backend_model"`   // backend_model varchar(255) NULL
		BackendType    string         `db:"backend_type"`    // backend_type varchar(8) NOT NULL  DEFAULT 'static'
		BackendTable   dbr.NullString `db:"backend_table"`   // backend_table varchar(255) NULL
		FrontendModel  dbr.NullString `db:"frontend_model"`  // frontend_model varchar(255) NULL
		FrontendInput  dbr.NullString `db:"frontend_input"`  // frontend_input varchar(50) NULL
		FrontendLabel  dbr.NullString `db:"frontend_label"`  // frontend_label varchar(255) NULL
		FrontendClass  dbr.NullString `db:"frontend_class"`  // frontend_class varchar(255) NULL
		SourceModel    dbr.NullString `db:"source_model"`    // source_model varchar(255) NULL
		IsRequired     bool           `db:"is_required"`     // is_required smallint(5) unsigned NOT NULL  DEFAULT '0'
		IsUserDefined  bool           `db:"is_user_defined"` // is_user_defined smallint(5) unsigned NOT NULL  DEFAULT '0'
		DefaultValue   dbr.NullString `db:"default_value"`   // default_value text NULL
		IsUnique       bool           `db:"is_unique"`       // is_unique smallint(5) unsigned NOT NULL  DEFAULT '0'
		Note           dbr.NullString `db:"note"`            // note varchar(255) NULL
	}

	AttributeGroupSlice []*AttributeGroup
	AttributeGroup      struct {
		AttributeGroupId   int64          `db:"attribute_group_id"`   // attribute_group_id smallint(5) unsigned NOT NULL PRI  auto_increment
		AttributeSetId     int64          `db:"attribute_set_id"`     // attribute_set_id smallint(5) unsigned NOT NULL MUL DEFAULT '0'
		AttributeGroupName dbr.NullString `db:"attribute_group_name"` // attribute_group_name varchar(255) NULL
		SortOrder          int64          `db:"sort_order"`           // sort_order smallint(6) NOT NULL  DEFAULT '0'
		DefaultId          dbr.NullInt64  `db:"default_id"`           // default_id smallint(5) unsigned NULL  DEFAULT '0'
	}

	AttributeLabelSlice []*AttributeLabel
	AttributeLabel      struct {
		AttributeLabelId int64          `db:"attribute_label_id"` // attribute_label_id int(10) unsigned NOT NULL PRI  auto_increment
		AttributeId      int64          `db:"attribute_id"`       // attribute_id smallint(5) unsigned NOT NULL MUL DEFAULT '0'
		StoreId          int64          `db:"store_id"`           // store_id smallint(5) unsigned NOT NULL MUL DEFAULT '0'
		Value            dbr.NullString `db:"value"`              // value varchar(255) NULL
	}

	AttributeOptionSlice []*AttributeOption
	AttributeOption      struct {
		OptionId    int64 `db:"option_id"`    // option_id int(10) unsigned NOT NULL PRI  auto_increment
		AttributeId int64 `db:"attribute_id"` // attribute_id smallint(5) unsigned NOT NULL MUL DEFAULT '0'
		SortOrder   int64 `db:"sort_order"`   // sort_order smallint(5) unsigned NOT NULL  DEFAULT '0'
	}

	AttributeOptionValueSlice []*AttributeOptionValue
	AttributeOptionValue      struct {
		ValueId  int64          `db:"value_id"`  // value_id int(10) unsigned NOT NULL PRI  auto_increment
		OptionId int64          `db:"option_id"` // option_id int(10) unsigned NOT NULL MUL DEFAULT '0'
		StoreId  int64          `db:"store_id"`  // store_id smallint(5) unsigned NOT NULL MUL DEFAULT '0'
		Value    dbr.NullString `db:"value"`     // value varchar(255) NULL
	}

	AttributeSetSlice []*AttributeSet
	AttributeSet      struct {
		AttributeSetId   int64          `db:"attribute_set_id"`   // attribute_set_id smallint(5) unsigned NOT NULL PRI  auto_increment
		EntityTypeId     int64          `db:"entity_type_id"`     // entity_type_id smallint(5) unsigned NOT NULL MUL DEFAULT '0'
		AttributeSetName dbr.NullString `db:"attribute_set_name"` // attribute_set_name varchar(255) NULL
		SortOrder        int64          `db:"sort_order"`         // sort_order smallint(6) NOT NULL  DEFAULT '0'
	}

	EntitySlice []*Entity
	Entity      struct {
		EntityId       int64          `db:"entity_id"`        // entity_id int(10) unsigned NOT NULL PRI  auto_increment
		EntityTypeId   int64          `db:"entity_type_id"`   // entity_type_id smallint(5) unsigned NOT NULL MUL DEFAULT '0'
		AttributeSetId int64          `db:"attribute_set_id"` // attribute_set_id smallint(5) unsigned NOT NULL  DEFAULT '0'
		IncrementId    dbr.NullString `db:"increment_id"`     // increment_id varchar(50) NULL
		ParentId       int64          `db:"parent_id"`        // parent_id int(10) unsigned NOT NULL  DEFAULT '0'
		StoreId        int64          `db:"store_id"`         // store_id smallint(5) unsigned NOT NULL MUL DEFAULT '0'
		CreatedAt      time.Time      `db:"created_at"`       // created_at timestamp NOT NULL  DEFAULT 'CURRENT_TIMESTAMP' on update CURRENT_TIMESTAMP
		UpdatedAt      time.Time      `db:"updated_at"`       // updated_at timestamp NOT NULL  DEFAULT '0000-00-00 00:00:00'
		IsActive       bool           `db:"is_active"`        // is_active smallint(5) unsigned NOT NULL  DEFAULT '1'
	}

	EntityAttributeSlice []*EntityAttribute
	EntityAttribute      struct {
		EntityAttributeId int64 `db:"entity_attribute_id"` // entity_attribute_id int(10) unsigned NOT NULL PRI  auto_increment
		EntityTypeId      int64 `db:"entity_type_id"`      // entity_type_id smallint(5) unsigned NOT NULL  DEFAULT '0'
		AttributeSetId    int64 `db:"attribute_set_id"`    // attribute_set_id smallint(5) unsigned NOT NULL MUL DEFAULT '0'
		AttributeGroupId  int64 `db:"attribute_group_id"`  // attribute_group_id smallint(5) unsigned NOT NULL MUL DEFAULT '0'
		AttributeId       int64 `db:"attribute_id"`        // attribute_id smallint(5) unsigned NOT NULL MUL DEFAULT '0'
		SortOrder         int64 `db:"sort_order"`          // sort_order smallint(6) NOT NULL  DEFAULT '0'
	}

	EntityDatetimeSlice []*EntityDatetime
	EntityDatetime      struct {
		ValueId      int64     `db:"value_id"`       // value_id int(11) NOT NULL PRI  auto_increment
		EntityTypeId int64     `db:"entity_type_id"` // entity_type_id smallint(5) unsigned NOT NULL MUL DEFAULT '0'
		AttributeId  int64     `db:"attribute_id"`   // attribute_id smallint(5) unsigned NOT NULL MUL DEFAULT '0'
		StoreId      int64     `db:"store_id"`       // store_id smallint(5) unsigned NOT NULL MUL DEFAULT '0'
		EntityId     int64     `db:"entity_id"`      // entity_id int(10) unsigned NOT NULL MUL DEFAULT '0'
		Value        time.Time `db:"value"`          // value datetime NOT NULL  DEFAULT '0000-00-00 00:00:00'
	}

	EntityDecimalSlice []*EntityDecimal
	EntityDecimal      struct {
		ValueId      int64   `db:"value_id"`       // value_id int(11) NOT NULL PRI  auto_increment
		EntityTypeId int64   `db:"entity_type_id"` // entity_type_id smallint(5) unsigned NOT NULL MUL DEFAULT '0'
		AttributeId  int64   `db:"attribute_id"`   // attribute_id smallint(5) unsigned NOT NULL MUL DEFAULT '0'
		StoreId      int64   `db:"store_id"`       // store_id smallint(5) unsigned NOT NULL MUL DEFAULT '0'
		EntityId     int64   `db:"entity_id"`      // entity_id int(10) unsigned NOT NULL MUL DEFAULT '0'
		Value        float64 `db:"value"`          // value decimal(12,4) NOT NULL  DEFAULT '0.0000'
	}

	EntityIntSlice []*EntityInt
	EntityInt      struct {
		ValueId      int64 `db:"value_id"`       // value_id int(11) NOT NULL PRI  auto_increment
		EntityTypeId int64 `db:"entity_type_id"` // entity_type_id smallint(5) unsigned NOT NULL MUL DEFAULT '0'
		AttributeId  int64 `db:"attribute_id"`   // attribute_id smallint(5) unsigned NOT NULL MUL DEFAULT '0'
		StoreId      int64 `db:"store_id"`       // store_id smallint(5) unsigned NOT NULL MUL DEFAULT '0'
		EntityId     int64 `db:"entity_id"`      // entity_id int(10) unsigned NOT NULL MUL DEFAULT '0'
		Value        int64 `db:"value"`          // value int(11) NOT NULL  DEFAULT '0'
	}

	EntityStoreSlice []*EntityStore
	EntityStore      struct {
		EntityStoreId   int64          `db:"entity_store_id"`   // entity_store_id int(10) unsigned NOT NULL PRI  auto_increment
		EntityTypeId    int64          `db:"entity_type_id"`    // entity_type_id smallint(5) unsigned NOT NULL MUL DEFAULT '0'
		StoreId         int64          `db:"store_id"`          // store_id smallint(5) unsigned NOT NULL MUL DEFAULT '0'
		IncrementPrefix dbr.NullString `db:"increment_prefix"`  // increment_prefix varchar(20) NULL
		IncrementLastId dbr.NullString `db:"increment_last_id"` // increment_last_id varchar(50) NULL
	}

	EntityTextSlice []*EntityText
	EntityText      struct {
		ValueId      int64  `db:"value_id"`       // value_id int(11) NOT NULL PRI  auto_increment
		EntityTypeId int64  `db:"entity_type_id"` // entity_type_id smallint(5) unsigned NOT NULL MUL DEFAULT '0'
		AttributeId  int64  `db:"attribute_id"`   // attribute_id smallint(5) unsigned NOT NULL MUL DEFAULT '0'
		StoreId      int64  `db:"store_id"`       // store_id smallint(5) unsigned NOT NULL MUL DEFAULT '0'
		EntityId     int64  `db:"entity_id"`      // entity_id int(10) unsigned NOT NULL MUL DEFAULT '0'
		Value        string `db:"value"`          // value text NOT NULL
	}

	EntityTypeSlice []*EntityType
	EntityType      struct {
		EntityTypeId              int64          `db:"entity_type_id"`              // entity_type_id smallint(5) unsigned NOT NULL PRI  auto_increment
		EntityTypeCode            string         `db:"entity_type_code"`            // entity_type_code varchar(50) NOT NULL MUL
		EntityModel               string         `db:"entity_model"`                // entity_model varchar(255) NOT NULL
		AttributeModel            dbr.NullString `db:"attribute_model"`             // attribute_model varchar(255) NULL
		EntityTable               dbr.NullString `db:"entity_table"`                // entity_table varchar(255) NULL
		ValueTablePrefix          dbr.NullString `db:"value_table_prefix"`          // value_table_prefix varchar(255) NULL
		EntityIdField             dbr.NullString `db:"entity_id_field"`             // entity_id_field varchar(255) NULL
		IsDataSharing             bool           `db:"is_data_sharing"`             // is_data_sharing smallint(5) unsigned NOT NULL  DEFAULT '1'
		DataSharingKey            dbr.NullString `db:"data_sharing_key"`            // data_sharing_key varchar(100) NULL  DEFAULT 'default'
		DefaultAttributeSetId     int64          `db:"default_attribute_set_id"`    // default_attribute_set_id smallint(5) unsigned NOT NULL  DEFAULT '0'
		IncrementModel            dbr.NullString `db:"increment_model"`             // increment_model varchar(255) NULL
		IncrementPerStore         int64          `db:"increment_per_store"`         // increment_per_store smallint(5) unsigned NOT NULL  DEFAULT '0'
		IncrementPadLength        int64          `db:"increment_pad_length"`        // increment_pad_length smallint(5) unsigned NOT NULL  DEFAULT '8'
		IncrementPadChar          string         `db:"increment_pad_char"`          // increment_pad_char varchar(1) NOT NULL  DEFAULT '0'
		AdditionalAttributeTable  dbr.NullString `db:"additional_attribute_table"`  // additional_attribute_table varchar(255) NULL
		EntityAttributeCollection dbr.NullString `db:"entity_attribute_collection"` // entity_attribute_collection varchar(255) NULL
	}

	EntityVarcharSlice []*EntityVarchar
	EntityVarchar      struct {
		ValueId      int64          `db:"value_id"`       // value_id int(11) NOT NULL PRI  auto_increment
		EntityTypeId int64          `db:"entity_type_id"` // entity_type_id smallint(5) unsigned NOT NULL MUL DEFAULT '0'
		AttributeId  int64          `db:"attribute_id"`   // attribute_id smallint(5) unsigned NOT NULL MUL DEFAULT '0'
		StoreId      int64          `db:"store_id"`       // store_id smallint(5) unsigned NOT NULL MUL DEFAULT '0'
		EntityId     int64          `db:"entity_id"`      // entity_id int(10) unsigned NOT NULL MUL DEFAULT '0'
		Value        dbr.NullString `db:"value"`          // value varchar(255) NULL
	}

	FormElementSlice []*FormElement
	FormElement      struct {
		ElementId   int64         `db:"element_id"`   // element_id int(10) unsigned NOT NULL PRI  auto_increment
		TypeId      int64         `db:"type_id"`      // type_id smallint(5) unsigned NOT NULL MUL
		FieldsetId  dbr.NullInt64 `db:"fieldset_id"`  // fieldset_id smallint(5) unsigned NULL MUL
		AttributeId int64         `db:"attribute_id"` // attribute_id smallint(5) unsigned NOT NULL MUL
		SortOrder   int64         `db:"sort_order"`   // sort_order int(11) NOT NULL  DEFAULT '0'
	}

	FormFieldsetSlice []*FormFieldset
	FormFieldset      struct {
		FieldsetId int64  `db:"fieldset_id"` // fieldset_id smallint(5) unsigned NOT NULL PRI  auto_increment
		TypeId     int64  `db:"type_id"`     // type_id smallint(5) unsigned NOT NULL MUL
		Code       string `db:"code"`        // code varchar(64) NOT NULL
		SortOrder  int64  `db:"sort_order"`  // sort_order int(11) NOT NULL  DEFAULT '0'
	}

	FormFieldsetLabelSlice []*FormFieldsetLabel
	FormFieldsetLabel      struct {
		FieldsetId int64  `db:"fieldset_id"` // fieldset_id smallint(5) unsigned NOT NULL PRI
		StoreId    int64  `db:"store_id"`    // store_id smallint(5) unsigned NOT NULL PRI
		Label      string `db:"label"`       // label varchar(255) NOT NULL
	}

	FormTypeSlice []*FormType
	FormType      struct {
		TypeId   int64          `db:"type_id"`   // type_id smallint(5) unsigned NOT NULL PRI  auto_increment
		Code     string         `db:"code"`      // code varchar(64) NOT NULL MUL
		Label    string         `db:"label"`     // label varchar(255) NOT NULL
		IsSystem bool           `db:"is_system"` // is_system smallint(5) unsigned NOT NULL  DEFAULT '0'
		Theme    dbr.NullString `db:"theme"`     // theme varchar(64) NULL
		StoreId  int64          `db:"store_id"`  // store_id smallint(5) unsigned NOT NULL MUL
	}

	FormTypeEntitySlice []*FormTypeEntity
	FormTypeEntity      struct {
		TypeId       int64 `db:"type_id"`        // type_id smallint(5) unsigned NOT NULL PRI
		EntityTypeId int64 `db:"entity_type_id"` // entity_type_id smallint(5) unsigned NOT NULL PRI
	}
)
