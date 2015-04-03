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

	"github.com/corestoreio/csfw/storage/csdb"
	"github.com/corestoreio/csfw/storage/dbr"
	"github.com/juju/errgo"
)

var (
	ErrAttributeNotFound = errors.New("Attribute not found")
	// AttributeCoreColumns defnies the minimal required coulmns for table eav_attribute.
	// Developers can extend the table eav_attribute with additional columns but these additional
	// columns with its method receivers must get generated in the attribute materialize function.
	// These core columns are already defined below.
	AttributeCoreColumns = csdb.TableCoreColumns{
		"attribute_id",
		"entity_type_id",
		"attribute_code",
		"attribute_model", // this column is unused by Mage1+2
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
	}
)

const (
	// TypeStatic use to check if an attribute is static, means part of the eav prefix table
	TypeStatic string = "static"
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
	// flexibility to extend with other custom structs.
	// @see magento2/site/app/code/Magento/Eav/Model/Entity/Attribute/AbstractAttribute.php
	Attributer interface {
		// IsStatic checks if an attribute is a static one
		IsStatic() bool
		// EntityType returns EntityType object or an error
		EntityType() (*CSEntityType, error)
		// UsesSource checks whether possible attribute values are retrieved from a finite source
		UsesSource() bool
		// IsInSet checks if attribute in specified attribute set
		IsInSet(int64) bool
		// IsInGroup checks if attribute in specified attribute group
		IsInGroup(int64) bool

		//@see magento2/site/app/code/Magento/Eav/Model/Entity/Attribute/AbstractAttribute.php
		//FlatColumns() []fancyTableColumnType @todo
		//FlatIndexes() []fancyTableColumnIndex @todo
		//Options()

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

	// Attribute defines properties for an attribute. This struct must be embedded in other EAV attributes.
	Attribute struct {
		attributeID   int64
		entityTypeID  int64
		attributeCode string
		backendModel  AttributeBackendModeller
		backendType   string
		backendTable  string
		frontendModel AttributeFrontendModeller
		frontendInput string
		frontendLabel string
		frontendClass string
		sourceModel   AttributeSourceModeller
		isRequired    bool
		isUserDefined bool
		defaultValue  string
		isUnique      bool
		note          string
	}

	// AttributeGetter implements functions on how to retrieve directly a certain attribute. This interface
	// is used in concrete entity models by generated code.
	AttributeGetter interface {
		// ByID returns an index using the AttributeID. This index identifies an attribute within an AttributeSlice.
		ByID(id int64) (AttributeIndex, error)
		// ByCode returns an index using the AttributeCode. This index identifies an attribute within an AttributeSlice.
		ByCode(code string) (AttributeIndex, error)
	}
)

var _ Attributer = (*Attribute)(nil)

// NewAttribute only for use in auto generated code. Looks terrible 8-)
func NewAttribute(
	attributeCode string,
	attributeID int64,
	backendModel AttributeBackendModeller,
	backendTable string,
	backendType string,
	defaultValue string,
	entityTypeID int64,
	frontendClass string,
	frontendInput string,
	frontendLabel string,
	frontendModel AttributeFrontendModeller,
	isRequired bool,
	isUnique bool,
	isUserDefined bool,
	note string,
	sourceModel AttributeSourceModeller,
) *Attribute {
	return &Attribute{
		attributeID:   attributeID,
		entityTypeID:  entityTypeID,
		attributeCode: attributeCode,
		backendModel:  backendModel,
		backendType:   backendType,
		backendTable:  backendTable,
		frontendModel: frontendModel,
		frontendInput: frontendInput,
		frontendLabel: frontendLabel,
		frontendClass: frontendClass,
		sourceModel:   sourceModel,
		isRequired:    isRequired,
		isUserDefined: isUserDefined,
		defaultValue:  defaultValue,
		isUnique:      isUnique,
		note:          note,
	}
}

// IsStatic checks if an attribute is a static one
func (a *Attribute) IsStatic() bool {
	return a.backendType == TypeStatic || a.backendType == ""
}

// EntityType returns EntityType object or an error
func (a *Attribute) EntityType() (*CSEntityType, error) {
	return GetEntityTypeCollection().GetByID(a.entityTypeID)
}

// UsesSource checks whether possible attribute values are retrieved from a finite source
func (a *Attribute) UsesSource() bool {
	switch {
	case a.frontendInput == "select":
		return true
	case a.frontendInput == "multiselect":
		return true
	case a.sourceModel != nil:
		return true
	}
	return false
}

// IsInSet checks if attribute in specified attribute set @todo
func (a *Attribute) IsInSet(_ int64) bool {
	return false
}

// IsInGroup checks if attribute in specified attribute group @todo
func (a *Attribute) IsInGroup(_ int64) bool {
	return false
}

func (a *Attribute) AttributeID() int64 {
	return a.attributeID
}
func (a *Attribute) EntityTypeID() int64 {
	return a.entityTypeID
}
func (a *Attribute) AttributeCode() string {
	return a.attributeCode
}
func (a *Attribute) BackendModel() AttributeBackendModeller {
	return a.backendModel
}
func (a *Attribute) BackendType() string {
	return a.backendType
}

// BackendTable returns the attribute backend table name. This function panics.
// @see magento2/site/app/code/Magento/Eav/Model/Entity/Attribute/AbstractAttribute.php::getBackendTable
func (a *Attribute) BackendTable() string {
	et, err := a.EntityType()
	if err != nil {
		panic("EntityType not found for attribute " + a.attributeCode)
	}

	if a.IsStatic() {
		// this means that the attribute is directly a column of the base entity table like
		// catalog_product_entity or customer_entity or eav_entity
		return et.GetValueTablePrefix()
	}
	if a.backendTable != "" {
		return a.backendTable
	}
	return et.GetEntityTablePrefix() + "_" + a.backendType
}

func (a *Attribute) FrontendModel() AttributeFrontendModeller {
	return a.frontendModel
}
func (a *Attribute) FrontendInput() string {
	return a.frontendInput
}
func (a *Attribute) FrontendLabel() string {
	return a.frontendLabel
}
func (a *Attribute) FrontendClass() string {
	return a.frontendClass
}
func (a *Attribute) SourceModel() AttributeSourceModeller {
	return a.sourceModel
}
func (a *Attribute) IsRequired() bool {
	return a.isRequired
}
func (a *Attribute) IsUserDefined() bool {
	return a.isUserDefined
}
func (a *Attribute) DefaultValue() string {
	return a.defaultValue
}
func (a *Attribute) IsUnique() bool {
	return a.isUnique
}
func (a *Attribute) Note() string {
	return a.note
}

// GetAttributeSelectSql generates the select query to retrieve full attribute configuration
func GetAttributeSelectSql(dbrSess dbr.SessionRunner, aat EntityTypeAdditionalAttributeTabler, entityTypeID, websiteId int64) (*dbr.SelectBuilder, error) {

	ta, err := GetTableStructure(TableIndexAttribute)
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
