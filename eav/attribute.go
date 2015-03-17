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
	"fmt"

	"github.com/corestoreio/csfw/storage/dbr"
	"github.com/juju/errgo"
)

func GetAttributeSelect(dbrSess dbr.SessionRunner, et *CSEntityType, websiteId int) error {

	tableStructs := et.AdditionalAttributeTable
	ta, err := GetTableStructure(TableAttribute)
	if err != nil {
		return errgo.Mask(err)
	}
	taa, err := tableStructs.TableAdditionalAttribute()
	if err != nil {
		return errgo.Mask(err)
	}

	selectSql := dbrSess.
		Select(ta.AllColumnAliasQuote("main_table")...).
		From(ta.Name+" AS `main_table`"). // @todo use a []string{"tablename","alias"}
		Join(
		taa.Name+" AS `additional_table`", // @todo use a []string{"tablename","alias"}
		taa.ColumnAliasQuote("additional_table"),
		dbr.JoinOn("`additional_table`.`attribute_id` = `main_table`.`attribute_id`"),
		dbr.JoinOn("`main_table`.`entity_type_id` = ?", et.EntityTypeID),
	)

	tew, err := tableStructs.TableEavWebsite()
	if err != nil {
		return errgo.Mask(err)
	}

	if tew != nil {
		selectSql.
			LeftJoin(
			tew.Name+" AS `scope_table`",
			tew.ColumnAliasQuote("scope_table"),
			dbr.JoinOn("scope_table.attribute_id = main_table.attribute_id"),
			dbr.JoinOn("scope_table.website_id = ?", websiteId),
		)
	}

	sql, args := selectSql.ToSql()

	fmt.Printf("\n%#v\n\n%s\n%#v\n", et, sql, args)

	return nil
}
