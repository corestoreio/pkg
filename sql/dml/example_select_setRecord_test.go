// Copyright 2015-2017, Cyrill @ Schumacher.fm and the CoreStore contributors
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

package dml_test

import (
	"fmt"
	"time"

	"github.com/corestoreio/pkg/sql/dml"
	"github.com/corestoreio/errors"
)

// Make sure that type catalogCategoryEntity implements interface
var _ dml.ColumnMapper = (*catalogCategoryEntity)(nil)
var _ dml.ColumnMapper = (*tableStore)(nil)

// catalogCategoryEntity defined somewhere in a different package.
type catalogCategoryEntity struct {
	EntityID       int64  // Auto Increment
	AttributeSetID int64  // From the EAV model
	ParentID       int64  // Other EntityID
	Path           string // e.g.: 1/2/20/21/26
	Position       int    // Position within the category tree
	CreatedAt      time.Time
}

func (ce *catalogCategoryEntity) MapColumns(cm *dml.ColumnMap) error {
	if cm.Mode() == dml.ColumnMapEntityReadAll {
		// This case gets executed when an INSERT statement doesn't contain any
		// columns, hence it requests all columns.
		return cm.Int64(&ce.EntityID).Int64(&ce.AttributeSetID).Int64(&ce.ParentID).String(&ce.Path).Int(&ce.Position).Time(&ce.CreatedAt).Err()
	}
	for cm.Next() {
		switch c := cm.Column(); c {
		case "entity_id":
			cm.Int64(&ce.EntityID)
		case "attribute_set_id":
			cm.Int64(&ce.AttributeSetID)
		case "parent_id":
			cm.Int64(&ce.ParentID)
		case "path":
			cm.String(&ce.Path)
		case "position":
			cm.Int(&ce.Position)
		case "created_at":
			cm.Time(&ce.CreatedAt)
		default:
			return errors.NewNotFoundf("[dml_test] %T: Column %q not found", ce, c)
		}
	}
	return cm.Err()
}

// tableStore defined somewhere in a different package.
type tableStore struct {
	StoreID   int64  // store_id smallint(5) unsigned NOT NULL PRI  auto_increment
	Code      string // code varchar(32) NULL UNI
	WebsiteID int64  // website_id smallint(5) unsigned NOT NULL MUL DEFAULT '0'
	GroupID   int64  // group_id smallint(5) unsigned NOT NULL MUL DEFAULT '0'
	Name      string // name varchar(255) NOT NULL
}

func (ts *tableStore) MapColumns(cm *dml.ColumnMap) error {
	if cm.Mode() == dml.ColumnMapEntityReadAll {
		// This case gets executed when an INSERT statement doesn't contain any
		// columns, hence it requests all columns.
		return cm.Int64(&ts.StoreID).String(&ts.Code).Int64(&ts.WebsiteID).Int64(&ts.GroupID).String(&ts.Name).Err()
	}
	for cm.Next() {
		switch c := cm.Column(); c {
		case "store_id":
			cm.Int64(&ts.StoreID)
		case "code":
			cm.String(&ts.Code)
		case "website_id":
			cm.Int64(&ts.WebsiteID)
		case "group_id":
			cm.Int64(&ts.GroupID)
		case "name":
			cm.String(&ts.Name)
		default:
			return errors.NewNotFoundf("[dml_test] %T: Column %q not found", ts, c)
		}
	}
	return cm.Err()
}

func ExampleSelect_BindRecord() {

	ce := &catalogCategoryEntity{678, 6, 670, "2/13/670/678", 0, now()}
	st := &tableStore{17, "ch-en", 2, 4, "Swiss EN Store"}

	s := dml.NewSelect("t_d.attribute_id", "e.entity_id").
		AddColumnsAliases("t_d.value", "default_value").
		AddColumnsConditions(dml.SQLIf("t_s.value_id IS NULL", "t_d.value", "t_s.value").Alias("value")).
		FromAlias("catalog_category_entity", "e").
		Join(
			dml.MakeIdentifier("catalog_category_entity_varchar").Alias("t_d"), // t_d = table scope default
			dml.Column("e.entity_id").Equal().Column("t_d.entity_id"),
		).
		LeftJoin(
			dml.MakeIdentifier("catalog_category_entity_varchar").Alias("t_s"), // t_s = table scope store
			dml.Column("t_s.attribute_id").Equal().Column("t_d.attribute_id"),
		).
		Where(
			dml.Column("e.entity_id").In().PlaceHolder(),                      // 678
			dml.Column("t_d.attribute_id").In().Int64s(45),                    // 45
			dml.Column("t_d.store_id").Equal().SQLIfNull("t_s.store_id", "0"), // Just for testing
			dml.Column("t_d.store_id").Equal().PlaceHolder(),                  // 17
		).
		BindRecord(dml.Qualify("e", ce), dml.Qualify("t_d", st))

	writeToSQLAndInterpolate(s)
	fmt.Print("\n\n")

	// Output:
	//Prepared Statement:
	//SELECT `t_d`.`attribute_id`, `e`.`entity_id`, `t_d`.`value` AS `default_value`,
	//IF((t_s.value_id IS NULL), t_d.value, t_s.value) AS `value` FROM
	//`catalog_category_entity` AS `e` INNER JOIN `catalog_category_entity_varchar` AS
	//`t_d` ON (`e`.`entity_id` = `t_d`.`entity_id`) LEFT JOIN
	//`catalog_category_entity_varchar` AS `t_s` ON (`t_s`.`attribute_id` =
	//`t_d`.`attribute_id`) WHERE (`e`.`entity_id` IN (?)) AND (`t_d`.`attribute_id`
	//IN (?)) AND (`t_d`.`store_id` = IFNULL(`t_s`.`store_id`,0)) AND
	//(`t_d`.`store_id` = ?)
	//Arguments: [678 45 17]
	//
	//Interpolated Statement:
	//SELECT `t_d`.`attribute_id`, `e`.`entity_id`, `t_d`.`value` AS `default_value`,
	//IF((t_s.value_id IS NULL), t_d.value, t_s.value) AS `value` FROM
	//`catalog_category_entity` AS `e` INNER JOIN `catalog_category_entity_varchar` AS
	//`t_d` ON (`e`.`entity_id` = `t_d`.`entity_id`) LEFT JOIN
	//`catalog_category_entity_varchar` AS `t_s` ON (`t_s`.`attribute_id` =
	//`t_d`.`attribute_id`) WHERE (`e`.`entity_id` IN (678)) AND (`t_d`.`attribute_id`
	//IN (45)) AND (`t_d`.`store_id` = IFNULL(`t_s`.`store_id`,0)) AND
	//(`t_d`.`store_id` = 17)

}
