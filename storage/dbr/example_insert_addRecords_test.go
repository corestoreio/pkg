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

// Make sure that type productEntity implements interface.
var _ dbr.ArgumentAssembler = (*productEntity)(nil)

// productEntity represents just a demo record.
type productEntity struct {
	EntityID       int64 // Auto Increment
	AttributeSetID int64
	TypeID         string
	SKU            dbr.NullString
	HasOptions     bool
}

func (pe productEntity) AssembleArguments(stmtType int, args dbr.Arguments, columns []string) (dbr.Arguments, error) {
	for _, c := range columns {
		switch c {
		case "attribute_set_id":
			args = append(args, dbr.ArgInt64(pe.AttributeSetID))
		case "type_id":
			args = append(args, dbr.ArgString(pe.TypeID))
		case "sku":
			args = append(args, pe.SKU)
		case "has_options":
			args = append(args, dbr.ArgBool(pe.HasOptions))
		default:
			return nil, errors.NewNotFoundf("[dbr_test] Column %q not found", c)
		}
	}
	if len(columns) == 0 && stmtType&(dbr.SQLPartValues) != 0 {
		args = append(args,
			dbr.ArgInt64(pe.EntityID),
			dbr.ArgInt64(pe.AttributeSetID),
			dbr.ArgString(pe.TypeID),
			pe.SKU,
			dbr.ArgBool(pe.HasOptions),
		)
	}
	return args, nil
}

func ExampleInsert_AddRecords() {

	objs := []productEntity{
		{1, 5, "simple", dbr.MakeNullString("SOA9"), false},
		{2, 5, "virtual", dbr.NullString{}, true},
	}

	i := dbr.NewInsert("catalog_product_entity").AddColumns("attribute_set_id", "type_id", "sku", "has_options").
		AddRecords(objs[0]).AddRecords(objs[1])
	writeToSqlAndPreprocess(i)

	fmt.Print("\n\n")
	i = dbr.NewInsert("catalog_product_entity").SetRecordValueCount(5).AddRecords(objs[0]).AddRecords(objs[1])
	writeToSqlAndPreprocess(i)

	// Output:
	//Prepared Statement:
	//INSERT INTO `catalog_product_entity`
	//(`attribute_set_id`,`type_id`,`sku`,`has_options`) VALUES (?,?,?,?),(?,?,?,?)
	//Arguments: [5 simple SOA9 false 5 virtual <nil> true]
	//
	//Preprocessed Statement:
	//INSERT INTO `catalog_product_entity`
	//(`attribute_set_id`,`type_id`,`sku`,`has_options`) VALUES
	//(5,'simple','SOA9',0),(5,'virtual',NULL,1)
	//
	//Prepared Statement:
	//INSERT INTO `catalog_product_entity` VALUES (?,?,?,?,?),(?,?,?,?,?)
	//Arguments: [1 5 simple SOA9 false 2 5 virtual <nil> true]
	//
	//Preprocessed Statement:
	//INSERT INTO `catalog_product_entity` VALUES
	//(1,5,'simple','SOA9',0),(2,5,'virtual',NULL,1)
}
