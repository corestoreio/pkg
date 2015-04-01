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
	"errors"

	"github.com/corestoreio/csfw/storage/dbr"
	"github.com/juju/errgo"
)

var (
	ErrAttributeNotFound = errors.New("Attribute not found")
)

type (
	// AttributeIndex used for index in a slice with constants (iota)
	AttributeIndex uint

	// Attributer defines the minimal requirements for one attribute. The interface relates to the table
	// eav_attribute. Developers can even extend this table with additional columns. Additional columns
	// will result into more generated functions. Even other EAV entities can embed this interface to
	// extend an attribute. For example the customer attributes
	// extends this interface two times. The column of func AttributeModel() has been removed because
	// it is not used and misplaced in the eav_attribute table. This interface can be extended to more than
	// 40 functions which is of course not idiomatic but for generated code it provides the best
	// flexibility.
	Attributer interface {
		AttributeID() int64
		EntityTypeID() int64
		AttributeCode() string
		BackendModel() AttributeBackendModeller
		BackendType() string
		BackendTable() string
		FrontendModel() AttributeFrontendModeller
		FrontendInput() string
		FrontendLabel() string
		FrontendClass() string
		SourceModel() AttributeSourceModeller
		IsRequired() bool
		IsUserDefined() bool
		DefaultValue() string
		IsUnique() bool
		Note() string
	}

	// AttributeGetter implements functions on how to retrieve directly a certain attribute
	AttributeGetter interface {
		// ByID returns an index using the AttributeID. This index identifies an attribute within an AttributeSlice.
		ByID(id int64) (AttributeIndex, error)
		// ByCode returns an index using the AttributeCode. This index identifies an attribute within an AttributeSlice.
		ByCode(code string) (AttributeIndex, error)
	}

	// AttributeBackendModeller defines the attribute backend model @todo
	AttributeBackendModeller interface {
		GetTable() string
		IsStatic() bool
		GetType()
		GetEntityIdField()
		SetValueId(valueId int)
		GetValueId()
		//AfterLoad($object);
		//BeforeSave($object);
		//AfterSave($object);
		//BeforeDelete($object);
		//AfterDelete($object);

		GetEntityValueId(entity *CSEntityType)

		SetEntityValueId(entity *CSEntityType, valueId int)
	}

	// AttributeFrontendModeller defines the attribute frontend model @todo
	AttributeFrontendModeller interface {
		TBD()
	}

	// AttributeSourceModeller defines the source where an attribute can also be stored @todo
	AttributeSourceModeller interface {
		TBD()
	}
)

// GetAttributeSelectSql generates the select query to retrieve full attribute configuration
func GetAttributeSelectSql(dbrSess dbr.SessionRunner, aat EntityTypeAdditionalAttributeTabler, entityTypeID, websiteId int64) (*dbr.SelectBuilder, error) {

	ta, err := GetTableStructure(TableAttribute)
	if err != nil {
		return nil, errgo.Mask(err)
	}
	taa, err := aat.TableAdditionalAttribute()
	if err != nil {
		return nil, errgo.Mask(err)
	}

	selectSql := dbrSess.
		Select(ta.AllColumnAliasQuote("main_table")...).
		From(ta.Name, "main_table").
		Join(
		dbr.JoinTable(taa.Name, "additional_table"),
		taa.ColumnAliasQuote("additional_table"),
		dbr.JoinOn("`additional_table`.`attribute_id` = `main_table`.`attribute_id`"),
		dbr.JoinOn("`main_table`.`entity_type_id` = ?", entityTypeID),
	)

	tew, err := aat.TableEavWebsite()
	if err != nil {
		return nil, errgo.Mask(err)
	}

	if tew != nil {
		const scopeTable = "scope_table"
		l := len(tew.Columns) * 2
		cols := make([]string, l)
		j := 0
		for i := 0; i < l; i = i + 2 {
			cols[i] = scopeTable + "." + tew.Columns[j] // real column name
			cols[i+1] = "scope_" + tew.Columns[j]       // alias column name
			j++
		}

		selectSql.
			LeftJoin(
			dbr.JoinTable(tew.Name, "scope_table"),
			dbr.ColumnAlias(cols...),
			dbr.JoinOn("`scope_table`.`attribute_id` = `main_table`.`attribute_id`"),
			dbr.JoinOn("`scope_table`.`website_id` = ?", websiteId),
		)
	}
	return selectSql, nil
}
