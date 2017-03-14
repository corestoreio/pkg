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

package eav

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
	// @see magento2/site/app/code/Magento/Eav/Api/Data/AttributeInterface.php
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

	// WSASlice WebSiteAttributeSlice is a slice of pointers to attributes. Despite being a slice the websiteID
	// is used as an index to return the attribute for a specific websiteID.
	WSASlice []*Attribute

	// Attribute defines properties for an attribute. This struct must be embedded in other EAV attributes.
	Attribute struct {
		// wa contains website specific attribute properties which are pulled out from the table
		// e.g. customer_eav_attribute_website or any other entity eav website table.
		// This pointer slice is in most cases nil. All struct fields can have different values in a website.
		wa WSASlice
		// websiteID if 0 then default values. If greater 0 then we have different values for a website.
		websiteID     int64
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
	// is used in concrete attribute models by generated code.
	// The logic behind this interface is to provide a fast access to the AttributeIndex. We will use as
	// key the id int64 or code string which will then map to the value of an AttributeIndex.
	AttributeGetter interface {
		// ByID returns an index using the AttributeID. This index identifies an attribute within an AttributeSlice.
		ByID(id int64) (AttributeIndex, error)
		// ByCode returns an index using the AttributeCode. This index identifies an attribute within an AttributeSlice.
		ByCode(code string) (AttributeIndex, error)
	}
	AttributeSliceGetter interface {
		Index(i AttributeIndex) interface{}
		Len() int
		ByID(g AttributeGetter, id int64) (interface{}, bool)
		ByCode(g AttributeGetter, code string) (interface{}, bool)
	}
)

var _ Attributer = (*Attribute)(nil)
var _ AttributeGetter = (*AttributeMapGet)(nil)

// Index returns the attribute from the current index or nil
func (s WSASlice) Index(i int64) *Attribute {
	if i > int64(len(s)) || i < 0 {
		return nil
	}
	return s[i] // can also be nil
}

// NewAttributeMapGet returns a new pointer to an AttributeMapGet.
func NewAttributeMapGet(i map[int64]AttributeIndex, c map[string]AttributeIndex) *AttributeMapGet {
	return &AttributeMapGet{
		i: i,
		c: c,
	}
}

// AttributeMapGet contains two maps for faster retrieving of the attribute index.
// Only used in generated code. Implements interface AttributeGetter.
type AttributeMapGet struct {
	i map[int64]AttributeIndex
	c map[string]AttributeIndex
}

// ByID returns an attribute index by the id from the database table
func (si *AttributeMapGet) ByID(id int64) (_ AttributeIndex, ok bool) {
	if i, ok := si.i[id]; ok {
		return i, ok
	}
	return AttributeIndex(0), false
}

// ByCode returns an attribute index
func (si *AttributeMapGet) ByCode(code string) (_ AttributeIndex, ok bool) {
	if c, ok := si.c[code]; ok {
		return c, ok
	}
	return AttributeIndex(0), false
}

// NewAttribute only for use in auto generated code. Looks terrible 8-)
// TODO: Use functional options because only a few fields are required.
func NewAttribute(
	wa WSASlice,
	websiteID int64,
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
		wa:            wa,
		websiteID:     websiteID,
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

// hasWA checks if the attribute has a website specific attribute
func (a *Attribute) hasWA() bool {
	return a.websiteID > 0 && a.wa.Index(a.websiteID) != nil
}

// IsStatic checks if an attribute is a static one. Considers websiteID.
func (a *Attribute) IsStatic() bool {
	wa := a
	if a.hasWA() {
		wa = a.wa.Index(a.websiteID)
	}
	return wa.backendType == TypeStatic || wa.backendType == ""
}

// EntityType returns EntityType object or an error. Does not consider websiteID.
func (a *Attribute) EntityType() (*CSEntityType, error) {
	return GetEntityTypeCollection().GetByID(a.entityTypeID)
}

/*
	@todo implement wa check like in IsStatic()
*/

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

// Handler internal wrapper for attribute collection C, getter G and entity type id.
// must be embedded into a concrete attribute struct. Implements interface EntityTypeAttributeModeller
// and EntityTypeAttributeCollectioner.
type Handler struct {
	// EntityTyeID to load the entity type. @todo implementation
	EntityTyeID int64
	// C collection of a materialized slice
	C AttributeSliceGetter
	// G getter knows how to return an AttributeIndex by id or by code
	G AttributeGetter
}

var _ EntityTypeAttributeModeller = (*Handler)(nil)
var _ EntityTypeAttributeCollectioner = (*Handler)(nil)

// New creates a new attribute and returns interface custattr.Attributer
func (h *Handler) New() interface{} {
	panic("Please override this method")
	return nil
}

// Get uses an AttributeIndex to return an attribute or an error.
// Use type assertion to convert to Attributer.
func (h *Handler) Get(i AttributeIndex) (_ interface{}, ok bool) {
	if h.C != nil && int(i) < h.C.Len() {
		return h.C.Index(i), true
	}
	return nil, true
}

// MustGet returns an attribute by AttributeIndex. Panics if not found.
// Use type assertion to convert to Attributer.
func (h *Handler) MustGet(i AttributeIndex) interface{} {
	a, ok := h.Get(i)
	if !ok {
		panic("Not found: fix me AttributeIndex")
	}
	return a
}

// GetByID returns an address attribute by its id
// Use type assertion to convert to Attributer.
func (h *Handler) GetByID(id int64) (_ interface{}, ok bool) {
	if h.C != nil {
		return h.C.ByID(h.G, id)
	}
	return nil, false
}

// GetByCode returns an address attribute by its code
// Use type assertion to convert to Attributer.
func (h *Handler) GetByCode(code string) (_ interface{}, ok bool) {
	if h.C != nil {
		return h.C.ByCode(h.G, code)
	}
	return nil, false
}

// Collection returns the full attribute collection AttributeSlice.
// You must use type assertion to convert to custattr.AttributeSlice.
func (h *Handler) Collection() interface{} {
	return h.C
}
