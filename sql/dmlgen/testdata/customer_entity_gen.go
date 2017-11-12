// Auto generated via github.com/corestoreio/pkg/sql/dmlgen

package testdata

import (
	"encoding/json"
	"time"

	"github.com/corestoreio/errors"
	"github.com/corestoreio/pkg/sql/dml"
)

// CustomerEntity represents a single row for DB table `customer_entity`.
// Auto generated.
type CustomerEntity struct {
	EntityID               uint64         // entity_id int(10) unsigned NOT NULL PRI  auto_increment "Entity Id"
	WebsiteID              dml.NullInt64  // website_id smallint(5) unsigned NULL MUL DEFAULT 'NULL'  "Website Id"
	Email                  dml.NullString // email varchar(255) NULL MUL DEFAULT 'NULL'  "Email"
	GroupID                uint64         // group_id smallint(5) unsigned NOT NULL  DEFAULT '0'  "Group Id"
	IncrementID            dml.NullString // increment_id varchar(50) NULL  DEFAULT 'NULL'  "Increment Id"
	StoreID                dml.NullInt64  // store_id smallint(5) unsigned NULL MUL DEFAULT '0'  "Store Id"
	CreatedAt              time.Time      // created_at timestamp NOT NULL  DEFAULT 'current_timestamp()'  "Created At"
	UpdatedAt              time.Time      // updated_at timestamp NOT NULL  DEFAULT 'current_timestamp()' on update current_timestamp() "Updated At"
	IsActive               bool           // is_active smallint(5) unsigned NOT NULL  DEFAULT '1'  "Is Active"
	DisableAutoGroupChange uint64         // disable_auto_group_change smallint(5) unsigned NOT NULL  DEFAULT '0'  "Disable automatic group change based on VAT ID"
	CreatedIn              dml.NullString // created_in varchar(255) NULL  DEFAULT 'NULL'  "Created From"
	Prefix                 dml.NullString // prefix varchar(40) NULL  DEFAULT 'NULL'  "Name Prefix"
	Firstname              dml.NullString // firstname varchar(255) NULL MUL DEFAULT 'NULL'  "First Name"
	Middlename             dml.NullString // middlename varchar(255) NULL  DEFAULT 'NULL'  "Middle Name/Initial"
	Lastname               dml.NullString // lastname varchar(255) NULL MUL DEFAULT 'NULL'  "Last Name"
	Suffix                 dml.NullString // suffix varchar(40) NULL  DEFAULT 'NULL'  "Name Suffix"
	Dob                    dml.NullTime   // dob date NULL  DEFAULT 'NULL'  "Date of Birth"
	PasswordHash           dml.NullString // password_hash varchar(128) NULL  DEFAULT 'NULL'  "Password_hash"
	RpToken                dml.NullString // rp_token varchar(128) NULL  DEFAULT 'NULL'  "Reset password token"
	RpTokenCreatedAt       dml.NullTime   // rp_token_created_at datetime NULL  DEFAULT 'NULL'  "Reset password token creation time"
	DefaultBilling         dml.NullInt64  // default_billing int(10) unsigned NULL  DEFAULT 'NULL'  "Default Billing Address"
	DefaultShipping        dml.NullInt64  // default_shipping int(10) unsigned NULL  DEFAULT 'NULL'  "Default Shipping Address"
	Taxvat                 dml.NullString // taxvat varchar(50) NULL  DEFAULT 'NULL'  "Tax/VAT Number"
	Confirmation           dml.NullString // confirmation varchar(64) NULL  DEFAULT 'NULL'  "Is Confirmed"
	Gender                 dml.NullInt64  // gender smallint(5) unsigned NULL  DEFAULT 'NULL'  "Gender"
	FailuresNum            dml.NullInt64  // failures_num smallint(6) NULL  DEFAULT '0'  "Failure Number"
	FirstFailure           dml.NullTime   // first_failure timestamp NULL  DEFAULT 'NULL'  "First Failure"
	LockExpires            dml.NullTime   // lock_expires timestamp NULL  DEFAULT 'NULL'  "Lock Expiration Date"
}

// NewCustomerEntity creates a new pointer with pre-initialized fields. Auto
// generated.
func NewCustomerEntity() *CustomerEntity {
	return &CustomerEntity{}
}

// AssignLastInsertID updates the increment ID field with the last inserted ID
// from an INSERT operation. Implements dml.InsertIDAssigner. Auto generated.
func (e *CustomerEntity) AssignLastInsertID(id int64) {
	e.EntityID = uint64(id)
}

// MapColumns implements interface ColumnMapper only partially. Auto generated.
func (e *CustomerEntity) MapColumns(cm *dml.ColumnMap) error {
	if cm.Mode() == dml.ColumnMapEntityReadAll {
		return cm.Uint64(&e.EntityID).NullInt64(&e.WebsiteID).NullString(&e.Email).Uint64(&e.GroupID).NullString(&e.IncrementID).NullInt64(&e.StoreID).Time(&e.CreatedAt).Time(&e.UpdatedAt).Bool(&e.IsActive).Uint64(&e.DisableAutoGroupChange).NullString(&e.CreatedIn).NullString(&e.Prefix).NullString(&e.Firstname).NullString(&e.Middlename).NullString(&e.Lastname).NullString(&e.Suffix).NullTime(&e.Dob).NullString(&e.PasswordHash).NullString(&e.RpToken).NullTime(&e.RpTokenCreatedAt).NullInt64(&e.DefaultBilling).NullInt64(&e.DefaultShipping).NullString(&e.Taxvat).NullString(&e.Confirmation).NullInt64(&e.Gender).NullInt64(&e.FailuresNum).NullTime(&e.FirstFailure).NullTime(&e.LockExpires).Err()
	}
	for cm.Next() {
		switch c := cm.Column(); c {
		case "entity_id", "customer_id", "parent_id":
			cm.Uint64(&e.EntityID)
		case "website_id":
			cm.NullInt64(&e.WebsiteID)
		case "email":
			cm.NullString(&e.Email)
		case "group_id":
			cm.Uint64(&e.GroupID)
		case "increment_id":
			cm.NullString(&e.IncrementID)
		case "store_id":
			cm.NullInt64(&e.StoreID)
		case "created_at":
			cm.Time(&e.CreatedAt)
		case "updated_at":
			cm.Time(&e.UpdatedAt)
		case "is_active":
			cm.Bool(&e.IsActive)
		case "disable_auto_group_change":
			cm.Uint64(&e.DisableAutoGroupChange)
		case "created_in":
			cm.NullString(&e.CreatedIn)
		case "prefix":
			cm.NullString(&e.Prefix)
		case "firstname":
			cm.NullString(&e.Firstname)
		case "middlename":
			cm.NullString(&e.Middlename)
		case "lastname":
			cm.NullString(&e.Lastname)
		case "suffix":
			cm.NullString(&e.Suffix)
		case "dob":
			cm.NullTime(&e.Dob)
		case "password_hash":
			cm.NullString(&e.PasswordHash)
		case "rp_token":
			cm.NullString(&e.RpToken)
		case "rp_token_created_at":
			cm.NullTime(&e.RpTokenCreatedAt)
		case "default_billing":
			cm.NullInt64(&e.DefaultBilling)
		case "default_shipping":
			cm.NullInt64(&e.DefaultShipping)
		case "taxvat":
			cm.NullString(&e.Taxvat)
		case "confirmation":
			cm.NullString(&e.Confirmation)
		case "gender":
			cm.NullInt64(&e.Gender)
		case "failures_num":
			cm.NullInt64(&e.FailuresNum)
		case "first_failure":
			cm.NullTime(&e.FirstFailure)
		case "lock_expires":
			cm.NullTime(&e.LockExpires)
		default:
			return errors.NewNotFoundf("[testdata] CustomerEntity Column %q not found", c)
		}
	}
	return errors.WithStack(cm.Err())
}

// CustomerEntityCollection represents a collection type for DB table customer_entity
// Not thread safe. Auto generated.
type CustomerEntityCollection struct {
	Data             []*CustomerEntity
	BeforeMapColumns func(uint64, *CustomerEntity) error
	AfterMapColumns  func(uint64, *CustomerEntity) error
}

// MakeCustomerEntityCollection creates a new initialized collection. Auto generated.
func MakeCustomerEntityCollection() CustomerEntityCollection {
	return CustomerEntityCollection{
		Data: make([]*CustomerEntity, 0, 5),
	}
}

func (cc CustomerEntityCollection) scanColumns(cm *dml.ColumnMap, e *CustomerEntity, idx uint64) error {
	if err := cc.BeforeMapColumns(idx, e); err != nil {
		return errors.WithStack(err)
	}
	if err := e.MapColumns(cm); err != nil {
		return errors.WithStack(err)
	}
	if err := cc.AfterMapColumns(idx, e); err != nil {
		return errors.WithStack(err)
	}
	return nil
}

// MapColumns implements dml.ColumnMapper interface. Auto generated.
func (cc CustomerEntityCollection) MapColumns(cm *dml.ColumnMap) error {
	switch m := cm.Mode(); m {
	case dml.ColumnMapEntityReadAll, dml.ColumnMapEntityReadSet:
		for i, e := range cc.Data {
			if err := cc.scanColumns(cm, e, uint64(i)); err != nil {
				return errors.WithStack(err)
			}
		}
	case dml.ColumnMapScan:
		if cm.Count == 0 {
			cc.Data = cc.Data[:0]
		}
		e := NewCustomerEntity()
		if err := cc.scanColumns(cm, e, cm.Count); err != nil {
			return errors.WithStack(err)
		}
		cc.Data = append(cc.Data, e)
	case dml.ColumnMapCollectionReadSet:
		for cm.Next() {
			switch c := cm.Column(); c {
			case "entity_id", "customer_id", "parent_id":
				cm.Args = cm.Args.Uint64s(cc.EntityIDs()...)
			default:
				return errors.NewNotFoundf("[testdata] CustomerEntityCollection Column %q not found", c)
			}
		}
	default:
		return errors.NewNotSupportedf("[dml] Unknown Mode: %q", string(m))
	}
	return cm.Err()
}

// EntityIDs returns a slice or appends to a slice all values.
// Auto generated.
func (cc CustomerEntityCollection) EntityIDs(ret ...uint64) []uint64 {
	if ret == nil {
		ret = make([]uint64, 0, len(cc.Data))
	}
	for _, e := range cc.Data {
		ret = append(ret, e.EntityID)
	}
	return ret
}

// UnmarshalJSON implements interface json.Unmarshaler.
func (cc *CustomerEntityCollection) UnmarshalJSON(b []byte) (err error) {
	return json.Unmarshal(b, cc.Data)
}

// MarshalJSON implements interface json.Marshaler.
func (cc *CustomerEntityCollection) MarshalJSON() ([]byte, error) {
	return json.Marshal(cc.Data)
}

// TODO add MarshalText and UnmarshalText.
