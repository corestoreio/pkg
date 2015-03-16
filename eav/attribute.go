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
	"strconv"

	"github.com/gocraft/dbr"
	"github.com/juju/errgo"
	sq "github.com/lann/squirrel"
)

// @todo dbr should implement handling of joins ... then we can remove squirrel

func GetAttributeSelect(dbrSess dbr.SessionRunner, et *CSEntityType) error {

	tabelStructs := et.AdditionalAttributeTable
	ta, err := GetTableStructure(TableAttribute)
	if err != nil {
		return errgo.Mask(err)
	}
	taa, err := tabelStructs.TableAdditionalAttribute()
	if err != nil {
		return errgo.Mask(err)
	}
	attr := sq.
		Select(append(ta.ColumnAliasQuote("main_table"), taa.ColumnAliasQuote("additional_table"))...).
		From(ta.Name).
		Join(taa.Name + " AS `additional_table` ON `additional_table`.`attribute_id` = `main_table`.`attribute_id`" +
		" AND `main_table`.`entity_type_id` = " + strconv.FormatInt(et.EntityTypeID, 10))

	tew, err := tabelStructs.TableEavWebsite()
	if err != nil {
		return errgo.Mask(err)
	}
	if tew != nil {
		attr.Select(taa.ColumnAliasQuote("scope_table"))
		attr.Join()
	}

	return nil
}
