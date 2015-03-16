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

// Generates code for all EAV attribute types
package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/corestoreio/csfw/storage/csdb"
	"github.com/corestoreio/csfw/tools"
)

var (
	pkg        = flag.String("p", "", "Package name in template")
	run        = flag.Bool("run", false, "If true program runs")
	outputFile = flag.String("o", "", "Output file name")
)

type (
	dataContainer struct {
		Package, Tick string
	}
)

func main() {
	flag.Parse()

	if false == *run || *outputFile == "" || *pkg == "" {
		flag.Usage()
		os.Exit(1)
	}

	db, dbrConn, err := csdb.Connect()
	tools.LogFatal(err)
	defer db.Close()

	tplData := &dataContainer{
		Package: *pkg,
		Tick:    "`",
	}

	formatted, err := tools.GenerateCode(*pkg, tplEav, tplData)
	if err != nil {
		fmt.Printf("\n%s\n", formatted)
		tools.LogFatal(err)
	}

	ioutil.WriteFile(*outputFile, formatted, 0600)
}

/*
to retrieve the attributes. The eav library must implement:

EAV -> handle additional_attribute_table
    -> implements _getEavWebsiteTable -> see Mage_Eav_Model_Resource_Attribute

SELECT
  `main_table`.*,
  `additional_table`.*
FROM `eav_attribute` AS `main_table`
  INNER JOIN `catalog_eav_attribute` AS `additional_table`
    ON additional_table.attribute_id = main_table.attribute_id AND main_table.entity_type_id = 4

see Mage_Customer_Model_Resource_Attribute_Collection
SELECT
  `main_table`.`attribute_id`,
  `main_table`.`entity_type_id`,
  `main_table`.`attribute_code`,
  `main_table`.`attribute_model`,
  `main_table`.`backend_model`,
  `main_table`.`backend_type`,
  `main_table`.`backend_table`,
  `main_table`.`frontend_model`,
  `main_table`.`frontend_input`,
  `main_table`.`frontend_label`,
  `main_table`.`frontend_class`,
  `main_table`.`source_model`,
  `main_table`.`is_required`,
  `main_table`.`is_user_defined`,
  `main_table`.`default_value`,
  `main_table`.`is_unique`,
  `main_table`.`note`,
  `additional_table`.`is_visible`,
  `additional_table`.`input_filter`,
  `additional_table`.`multiline_count`,
  `additional_table`.`validate_rules`,
  `additional_table`.`is_system`,
  `additional_table`.`sort_order`,
  `additional_table`.`data_model`,
  `scope_table`.`website_id`      AS `scope_website_id`,
  `scope_table`.`is_visible`      AS `scope_is_visible`,
  `scope_table`.`is_required`     AS `scope_is_required`,
  `scope_table`.`default_value`   AS `scope_default_value`,
  `scope_table`.`multiline_count` AS `scope_multiline_count`
FROM `eav_attribute` AS `main_table`
  INNER JOIN `customer_eav_attribute` AS `additional_table` ON additional_table.attribute_id = main_table.attribute_id
  LEFT JOIN `customer_eav_attribute_website` AS `scope_table`
    ON scope_table.attribute_id = main_table.attribute_id AND scope_table.website_id = :scope_website_id
WHERE (main_table.entity_type_id = :mt_entity_type_id)

*/

//func getEntityTypeData(dbrSess *dbr.Session) (JsonEntityTypeMap, error) {
//
//	s, err := eav.GetTableStructure(eav.TableEntityType)
//	if err != nil {
//		return nil, errgo.Mask(err)
//	}
//
//	var entityTypeCollection eav.EntityTypeSlice
//	_, err = dbrSess.
//		Select(s.Columns...).
//		From(s.Name).
//		LoadStructs(&entityTypeCollection)
//	if err != nil {
//		return nil, errgo.Mask(err)
//	}
//
//	return mapCollection, nil
//}
