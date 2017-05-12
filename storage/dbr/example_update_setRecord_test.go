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

package dbr_test

import (
	"fmt"

	"github.com/corestoreio/csfw/storage/dbr"
	"github.com/corestoreio/errors"
)

// Make sure that type categoryEntity implements interface
var _ dbr.ArgumentAssembler = (*categoryEntity)(nil)

// categoryEntity represents just a demo record.
type categoryEntity struct {
	EntityID       int64 // Auto Increment
	AttributeSetID int64
	ParentID       string
	Path           dbr.NullString
}

func (pe *categoryEntity) AssembleArguments(stmtType rune, args dbr.Arguments, columns, condition []string) (dbr.Arguments, error) {
	for _, c := range columns {
		switch c {
		case "attribute_set_id":
			args = append(args, dbr.ArgInt64(pe.AttributeSetID))
		case "parent_id":
			args = append(args, dbr.ArgString(pe.ParentID))
		case "path":
			args = append(args, pe.Path)
		default:
			return nil, errors.NewNotFoundf("[dbr_test] Column %q not found", c)
		}
	}
	for _, c := range condition {
		switch c {
		case "entity_id":
			args = append(args, dbr.ArgInt64(pe.EntityID))
		}
	}

	return args, nil
}

func ExampleUpdate_SetRecord() {

	ce := &categoryEntity{345, 6, "p123", dbr.MakeNullString("4/5/6/7")}

	// Updates all rows in the table
	u := dbr.NewUpdate("catalog_category_entity").
		SetRecord([]string{"attribute_set_id", "parent_id", "path"}, ce)
	writeToSqlAndPreprocess(u)

	fmt.Print("\n\n")

	ce = &categoryEntity{678, 6, "p456", dbr.NullString{}}

	// Updates only one row in the table
	u = dbr.NewUpdate("catalog_category_entity").
		SetRecord([]string{"attribute_set_id", "parent_id", "path"}, ce).
		Where(dbr.Column("entity_id", dbr.ArgInt64())) // dbr.Equal is default operator!
	writeToSqlAndPreprocess(u)

	// Output:
	//Prepared Statement:
	//UPDATE `catalog_category_entity` SET `attribute_set_id`=?, `parent_id`=?,
	//`path`=?
	//Arguments: [6 p123 4/5/6/7]
	//
	//Preprocessed Statement:
	//UPDATE `catalog_category_entity` SET `attribute_set_id`=6, `parent_id`='p123',
	//`path`='4/5/6/7'
	//
	//Prepared Statement:
	//UPDATE `catalog_category_entity` SET `attribute_set_id`=?, `parent_id`=?,
	//`path`=? WHERE (`entity_id` = ?)
	//Arguments: [6 p456 <nil> 678]
	//
	//Preprocessed Statement:
	//UPDATE `catalog_category_entity` SET `attribute_set_id`=6, `parent_id`='p456',
	//`path`=NULL WHERE (`entity_id` = 678)
}
