// Code generated by codegen. DO NOT EDIT.
package dmltestgenerated3

import (
	"github.com/corestoreio/pkg/storage/null"
	"time"
)

// CatalogCategoryEntity represents a single row for DB table
// catalog_category_entity. Auto generated.
type CatalogCategoryEntity struct {
	EntityID                uint32                   // entity_id int(10) unsigned NOT NULL MUL   "Entity Id"
	RowID                   uint32                   // row_id int(10) unsigned NOT NULL PRI  auto_increment "Version Id"
	SequenceCatalogCategory *SequenceCatalogCategory // 1:1 catalog_category_entity.entity_id => sequence_catalog_category.sequence_value
}

// CatalogCategoryEntityCollection represents a collection type for DB table
// catalog_category_entity
// Not thread safe. Auto generated.
type CatalogCategoryEntityCollection struct {
	Data []*CatalogCategoryEntity `json:"data,omitempty"`
}

// NewCatalogCategoryEntityCollection  creates a new initialized collection. Auto
// generated.
func NewCatalogCategoryEntityCollection() *CatalogCategoryEntityCollection {
	return &CatalogCategoryEntityCollection{
		Data: make([]*CatalogCategoryEntity, 0, 5),
	}
}

// CustomerAddressEntity represents a single row for DB table
// customer_address_entity. Auto generated.
//easyjson:json
type CustomerAddressEntity struct {
	EntityID          uint32      `max_len:"10"` // entity_id int(10) unsigned NOT NULL PRI  auto_increment "Entity ID"
	IncrementID       null.String `max_len:"50"` // increment_id varchar(50) NULL  DEFAULT 'NULL'  "Increment Id"
	ParentID          null.Uint32 `max_len:"10"` // parent_id int(10) unsigned NULL MUL DEFAULT 'NULL'  "Parent ID"
	CreatedAt         time.Time   // created_at timestamp NOT NULL  DEFAULT 'current_timestamp()'  "Created At"
	UpdatedAt         time.Time   // updated_at timestamp NOT NULL  DEFAULT 'current_timestamp()' on update current_timestamp() "Updated At"
	IsActive          bool        `max_len:"5"`     // is_active smallint(5) unsigned NOT NULL  DEFAULT '1'  "Is Active"
	City              string      `max_len:"255"`   // city varchar(255) NOT NULL    "City"
	Company           null.String `max_len:"255"`   // company varchar(255) NULL  DEFAULT 'NULL'  "Company"
	CountryID         string      `max_len:"255"`   // country_id varchar(255) NOT NULL    "Country"
	Fax               null.String `max_len:"255"`   // fax varchar(255) NULL  DEFAULT 'NULL'  "Fax"
	Firstname         string      `max_len:"255"`   // firstname varchar(255) NOT NULL    "First Name"
	Lastname          string      `max_len:"255"`   // lastname varchar(255) NOT NULL    "Last Name"
	Middlename        null.String `max_len:"255"`   // middlename varchar(255) NULL  DEFAULT 'NULL'  "Middle Name"
	Postcode          null.String `max_len:"255"`   // postcode varchar(255) NULL  DEFAULT 'NULL'  "Zip/Postal Code"
	Prefix            null.String `max_len:"40"`    // prefix varchar(40) NULL  DEFAULT 'NULL'  "Name Prefix"
	Region            null.String `max_len:"255"`   // region varchar(255) NULL  DEFAULT 'NULL'  "State/Province"
	RegionID          null.Uint32 `max_len:"10"`    // region_id int(10) unsigned NULL  DEFAULT 'NULL'  "State/Province"
	Street            string      `max_len:"65535"` // street text NOT NULL    "Street Address"
	Suffix            null.String `max_len:"40"`    // suffix varchar(40) NULL  DEFAULT 'NULL'  "Name Suffix"
	Telephone         string      `max_len:"255"`   // telephone varchar(255) NOT NULL    "Phone Number"
	VatID             null.String `max_len:"255"`   // vat_id varchar(255) NULL  DEFAULT 'NULL'  "VAT number"
	VatIsValid        null.Bool   `max_len:"10"`    // vat_is_valid int(10) unsigned NULL  DEFAULT 'NULL'  "VAT number validity"
	VatRequestDate    null.String `max_len:"255"`   // vat_request_date varchar(255) NULL  DEFAULT 'NULL'  "VAT number validation request date"
	VatRequestID      null.String `max_len:"255"`   // vat_request_id varchar(255) NULL  DEFAULT 'NULL'  "VAT number validation request ID"
	VatRequestSuccess null.Uint32 `max_len:"10"`    // vat_request_success int(10) unsigned NULL  DEFAULT 'NULL'  "VAT number validation request success"
}

// CustomerAddressEntityCollection represents a collection type for DB table
// customer_address_entity
// Not thread safe. Auto generated.
//easyjson:json
type CustomerAddressEntityCollection struct {
	Data []*CustomerAddressEntity `json:"data,omitempty"`
}

// NewCustomerAddressEntityCollection  creates a new initialized collection. Auto
// generated.
func NewCustomerAddressEntityCollection() *CustomerAddressEntityCollection {
	return &CustomerAddressEntityCollection{
		Data: make([]*CustomerAddressEntity, 0, 5),
	}
}

// CustomerEntity represents a single row for DB table customer_entity. Auto
// generated.
//easyjson:json
type CustomerEntity struct {
	EntityID               uint32                           `max_len:"10"`  // entity_id int(10) unsigned NOT NULL PRI  auto_increment "Entity ID"
	WebsiteID              null.Uint16                      `max_len:"5"`   // website_id smallint(5) unsigned NULL MUL DEFAULT 'NULL'  "Website ID"
	Email                  null.String                      `max_len:"255"` // email varchar(255) NULL MUL DEFAULT 'NULL'  "Email"
	GroupID                uint16                           `max_len:"5"`   // group_id smallint(5) unsigned NOT NULL  DEFAULT '0'  "Group ID"
	IncrementID            null.String                      `max_len:"50"`  // increment_id varchar(50) NULL  DEFAULT 'NULL'  "Increment Id"
	StoreID                null.Uint16                      `max_len:"5"`   // store_id smallint(5) unsigned NULL MUL DEFAULT '0'  "Store ID"
	CreatedAt              time.Time                        // created_at timestamp NOT NULL  DEFAULT 'current_timestamp()'  "Created At"
	UpdatedAt              time.Time                        // updated_at timestamp NOT NULL  DEFAULT 'current_timestamp()' on update current_timestamp() "Updated At"
	IsActive               bool                             `max_len:"5"`   // is_active smallint(5) unsigned NOT NULL  DEFAULT '1'  "Is Active"
	DisableAutoGroupChange uint16                           `max_len:"5"`   // disable_auto_group_change smallint(5) unsigned NOT NULL  DEFAULT '0'  "Disable automatic group change based on VAT ID"
	CreatedIn              null.String                      `max_len:"255"` // created_in varchar(255) NULL  DEFAULT 'NULL'  "Created From"
	Prefix                 null.String                      `max_len:"40"`  // prefix varchar(40) NULL  DEFAULT 'NULL'  "Name Prefix"
	Firstname              null.String                      `max_len:"255"` // firstname varchar(255) NULL MUL DEFAULT 'NULL'  "First Name"
	Middlename             null.String                      `max_len:"255"` // middlename varchar(255) NULL  DEFAULT 'NULL'  "Middle Name/Initial"
	Lastname               null.String                      `max_len:"255"` // lastname varchar(255) NULL MUL DEFAULT 'NULL'  "Last Name"
	Suffix                 null.String                      `max_len:"40"`  // suffix varchar(40) NULL  DEFAULT 'NULL'  "Name Suffix"
	Dob                    null.Time                        // dob date NULL  DEFAULT 'NULL'  "Date of Birth"
	PasswordHash           null.String                      `max_len:"128"` // password_hash varchar(128) NULL  DEFAULT 'NULL'  "Password_hash"
	RpToken                null.String                      `max_len:"128"` // rp_token varchar(128) NULL  DEFAULT 'NULL'  "Reset password token"
	RpTokenCreatedAt       null.Time                        // rp_token_created_at datetime NULL  DEFAULT 'NULL'  "Reset password token creation time"
	DefaultBilling         null.Uint32                      `max_len:"10"` // default_billing int(10) unsigned NULL  DEFAULT 'NULL'  "Default Billing Address"
	DefaultShipping        null.Uint32                      `max_len:"10"` // default_shipping int(10) unsigned NULL  DEFAULT 'NULL'  "Default Shipping Address"
	Taxvat                 null.String                      `max_len:"50"` // taxvat varchar(50) NULL  DEFAULT 'NULL'  "Tax/VAT Number"
	Confirmation           null.String                      `max_len:"64"` // confirmation varchar(64) NULL  DEFAULT 'NULL'  "Is Confirmed"
	Gender                 null.Uint16                      `max_len:"5"`  // gender smallint(5) unsigned NULL  DEFAULT 'NULL'  "Gender"
	FailuresNum            null.Int16                       `max_len:"5"`  // failures_num smallint(6) NULL  DEFAULT '0'  "Failure Number"
	FirstFailure           null.Time                        // first_failure timestamp NULL  DEFAULT 'NULL'  "First Failure"
	LockExpires            null.Time                        // lock_expires timestamp NULL  DEFAULT 'NULL'  "Lock Expiration Date"
	Address                *CustomerAddressEntityCollection // Reversed 1:M customer_entity.entity_id => customer_address_entity.parent_id
}

// CustomerEntityCollection represents a collection type for DB table
// customer_entity
// Not thread safe. Auto generated.
//easyjson:json
type CustomerEntityCollection struct {
	Data []*CustomerEntity `json:"data,omitempty"`
}

// NewCustomerEntityCollection  creates a new initialized collection. Auto
// generated.
func NewCustomerEntityCollection() *CustomerEntityCollection {
	return &CustomerEntityCollection{
		Data: make([]*CustomerEntity, 0, 5),
	}
}

// SequenceCatalogCategory represents a single row for DB table
// sequence_catalog_category. Auto generated.
type SequenceCatalogCategory struct {
	SequenceValue         uint32                 // sequence_value int(10) unsigned NOT NULL PRI  auto_increment ""
	CatalogCategoryEntity *CatalogCategoryEntity // Reversed 1:1 sequence_catalog_category.sequence_value => catalog_category_entity.entity_id
}

// SequenceCatalogCategoryCollection represents a collection type for DB table
// sequence_catalog_category
// Not thread safe. Auto generated.
type SequenceCatalogCategoryCollection struct {
	Data []*SequenceCatalogCategory `json:"data,omitempty"`
}

// NewSequenceCatalogCategoryCollection  creates a new initialized collection.
// Auto generated.
func NewSequenceCatalogCategoryCollection() *SequenceCatalogCategoryCollection {
	return &SequenceCatalogCategoryCollection{
		Data: make([]*SequenceCatalogCategory, 0, 5),
	}
}

// Store represents a single row for DB table store. Auto generated.
type Store struct {
	StoreID      uint16        // store_id smallint(5) unsigned NOT NULL PRI  auto_increment "Store Id"
	Code         null.String   // code varchar(32) NULL UNI DEFAULT 'NULL'  "Code"
	WebsiteID    uint16        // website_id smallint(5) unsigned NOT NULL MUL DEFAULT '0'  "Website Id"
	GroupID      uint16        // group_id smallint(5) unsigned NOT NULL MUL DEFAULT '0'  "Group Id"
	Name         string        // name varchar(255) NOT NULL    "Store Name"
	SortOrder    uint16        // sort_order smallint(5) unsigned NOT NULL  DEFAULT '0'  "Store Sort Order"
	IsActive     bool          // is_active smallint(5) unsigned NOT NULL MUL DEFAULT '0'  "Store Activity"
	StoreGroup   *StoreGroup   // 1:1 store.group_id => store_group.group_id
	StoreWebsite *StoreWebsite // 1:1 store.website_id => store_website.website_id
}

// StoreCollection represents a collection type for DB table store
// Not thread safe. Auto generated.
type StoreCollection struct {
	Data []*Store `json:"data,omitempty"`
}

// NewStoreCollection  creates a new initialized collection. Auto generated.
func NewStoreCollection() *StoreCollection {
	return &StoreCollection{
		Data: make([]*Store, 0, 5),
	}
}

// StoreGroup represents a single row for DB table store_group. Auto generated.
type StoreGroup struct {
	GroupID        uint16           // group_id smallint(5) unsigned NOT NULL PRI  auto_increment "Group Id"
	WebsiteID      uint16           // website_id smallint(5) unsigned NOT NULL MUL DEFAULT '0'  "Website Id"
	Name           string           // name varchar(255) NOT NULL    "Store Group Name"
	RootCategoryID uint32           // root_category_id int(10) unsigned NOT NULL  DEFAULT '0'  "Root Category Id"
	DefaultStoreID uint16           // default_store_id smallint(5) unsigned NOT NULL MUL DEFAULT '0'  "Default Store Id"
	Code           null.String      // code varchar(32) NULL UNI DEFAULT 'NULL'  "Store group unique code"
	StoreWebsite   *StoreWebsite    // 1:1 store_group.website_id => store_website.website_id
	Store          *StoreCollection // Reversed 1:M store_group.group_id => store.group_id
}

// StoreGroupCollection represents a collection type for DB table store_group
// Not thread safe. Auto generated.
type StoreGroupCollection struct {
	Data []*StoreGroup `json:"data,omitempty"`
}

// NewStoreGroupCollection  creates a new initialized collection. Auto generated.
func NewStoreGroupCollection() *StoreGroupCollection {
	return &StoreGroupCollection{
		Data: make([]*StoreGroup, 0, 5),
	}
}

// StoreWebsite represents a single row for DB table store_website. Auto
// generated.
type StoreWebsite struct {
	WebsiteID      uint16                // website_id smallint(5) unsigned NOT NULL PRI  auto_increment "Website Id"
	Code           null.String           // code varchar(32) NULL UNI DEFAULT 'NULL'  "Code"
	Name           null.String           // name varchar(64) NULL  DEFAULT 'NULL'  "Website Name"
	SortOrder      uint16                // sort_order smallint(5) unsigned NOT NULL MUL DEFAULT '0'  "Sort Order"
	DefaultGroupID uint16                // default_group_id smallint(5) unsigned NOT NULL MUL DEFAULT '0'  "Default Group Id"
	IsDefault      null.Bool             // is_default smallint(5) unsigned NULL  DEFAULT '0'  "Defines Is Website Default"
	Store          *StoreCollection      // Reversed 1:M store_website.website_id => store.website_id
	StoreGroup     *StoreGroupCollection // Reversed 1:M store_website.website_id => store_group.website_id
}

// StoreWebsiteCollection represents a collection type for DB table store_website
// Not thread safe. Auto generated.
type StoreWebsiteCollection struct {
	Data []*StoreWebsite `json:"data,omitempty"`
}

// NewStoreWebsiteCollection  creates a new initialized collection. Auto
// generated.
func NewStoreWebsiteCollection() *StoreWebsiteCollection {
	return &StoreWebsiteCollection{
		Data: make([]*StoreWebsite, 0, 5),
	}
}
