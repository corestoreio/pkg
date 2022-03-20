// Code generated by corestoreio/pkg/util/codegen. DO NOT EDIT.
// Generated by sql/dmlgen. DO NOT EDIT.
package dmltestgenerated4

import (
	"context"
	"database/sql"
	"time"

	"github.com/corestoreio/errors"
	"github.com/corestoreio/pkg/sql/ddl"
	"github.com/corestoreio/pkg/sql/dml"
	"github.com/corestoreio/pkg/storage/null"
)

// TableName constants define the names of all tables.
const (
	TableNameCoreConfiguration         = "core_configuration"
	TableNameSalesOrderStatusState     = "sales_order_status_state"
	TableNameViewCustomerAutoIncrement = "view_customer_auto_increment"
)

// Columns struct provides for all tables the name of the columns. Allows type
// safety.
var Columns = struct {
	CoreConfiguration struct {
		ConfigID  string
		Scope     string
		ScopeID   string
		Expires   string
		Path      string
		Value     string
		VersionTs string
		VersionTe string
	}
	SalesOrderStatusState struct {
		Status         string
		State          string
		IsDefault      string
		VisibleOnFront string
	}
	ViewCustomerAutoIncrement struct {
		CeEntityID string
		Email      string
		Firstname  string
		Lastname   string
		City       string
	}
}{
	CoreConfiguration: struct {
		ConfigID  string
		Scope     string
		ScopeID   string
		Expires   string
		Path      string
		Value     string
		VersionTs string
		VersionTe string
	}{
		ConfigID:  "config_id",
		Scope:     "scope",
		ScopeID:   "scope_id",
		Expires:   "expires",
		Path:      "path",
		Value:     "value",
		VersionTs: "version_ts",
		VersionTe: "version_te",
	},
	SalesOrderStatusState: struct {
		Status         string
		State          string
		IsDefault      string
		VisibleOnFront string
	}{
		Status:         "status",
		State:          "state",
		IsDefault:      "is_default",
		VisibleOnFront: "visible_on_front",
	},
	ViewCustomerAutoIncrement: struct {
		CeEntityID string
		Email      string
		Firstname  string
		Lastname   string
		City       string
	}{
		CeEntityID: "ce_entity_id",
		Email:      "email",
		Firstname:  "firstname",
		Lastname:   "lastname",
		City:       "city",
	},
}

var dbmEmptyOpts = []dml.DBRFunc{func(dbr *dml.DBR) {
	// do nothing because Clone gets called automatically
}}

func dbmNoopResultCheckFn(_ sql.Result, err error) error { return err }

// Event functions are getting dispatched during before or after handling a
// collection or an entity.
// Context is always non-nil but either collection or entity pointer will be set.
type (
	EventCoreConfigurationFn         func(context.Context, *CoreConfigurations, *CoreConfiguration) error
	EventSalesOrderStatusStateFn     func(context.Context, *SalesOrderStatusStates, *SalesOrderStatusState) error
	EventViewCustomerAutoIncrementFn func(context.Context, *ViewCustomerAutoIncrements, *ViewCustomerAutoIncrement) error
)

// DBMOption provides various options to the DBM object.
type DBMOption struct {
	TableOptions                       []ddl.TableOption // gets applied at the beginning
	TableOptionsAfter                  []ddl.TableOption // gets applied at the end
	InitSelectFn                       func(*dml.Select) *dml.Select
	InitUpdateFn                       func(*dml.Update) *dml.Update
	InitDeleteFn                       func(*dml.Delete) *dml.Delete
	InitInsertFn                       func(*dml.Insert) *dml.Insert
	eventCoreConfigurationFunc         [dml.EventFlagMax][]EventCoreConfigurationFn
	eventSalesOrderStatusStateFunc     [dml.EventFlagMax][]EventSalesOrderStatusStateFn
	eventViewCustomerAutoIncrementFunc [dml.EventFlagMax][]EventViewCustomerAutoIncrementFn
}

// AddEventCoreConfiguration adds a specific defined event call back to the DBM.
// It panics if the event argument is larger than dml.EventFlagMax.
func (o *DBMOption) AddEventCoreConfiguration(event dml.EventFlag, fn EventCoreConfigurationFn) *DBMOption {
	o.eventCoreConfigurationFunc[event] = append(o.eventCoreConfigurationFunc[event], fn)
	return o
}

// AddEventSalesOrderStatusState adds a specific defined event call back to the
// DBM.
// It panics if the event argument is larger than dml.EventFlagMax.
func (o *DBMOption) AddEventSalesOrderStatusState(event dml.EventFlag, fn EventSalesOrderStatusStateFn) *DBMOption {
	o.eventSalesOrderStatusStateFunc[event] = append(o.eventSalesOrderStatusStateFunc[event], fn)
	return o
}

// AddEventViewCustomerAutoIncrement adds a specific defined event call back to
// the DBM.
// It panics if the event argument is larger than dml.EventFlagMax.
func (o *DBMOption) AddEventViewCustomerAutoIncrement(event dml.EventFlag, fn EventViewCustomerAutoIncrementFn) *DBMOption {
	o.eventViewCustomerAutoIncrementFunc[event] = append(o.eventViewCustomerAutoIncrementFunc[event], fn)
	return o
}

// DBM defines the DataBaseManagement object for the tables  core_configuration,
// sales_order_status_state, view_customer_auto_increment
type DBM struct {
	*ddl.Tables
	option DBMOption
}

func (dbm DBM) eventCoreConfigurationFunc(ctx context.Context, ef dml.EventFlag, skipEvents bool, ec *CoreConfigurations, e *CoreConfiguration) error {
	if len(dbm.option.eventCoreConfigurationFunc[ef]) == 0 || skipEvents {
		return nil
	}
	for _, fn := range dbm.option.eventCoreConfigurationFunc[ef] {
		if err := fn(ctx, ec, e); err != nil {
			return errors.WithStack(err)
		}
	}
	return nil
}

func (dbm DBM) eventSalesOrderStatusStateFunc(ctx context.Context, ef dml.EventFlag, skipEvents bool, ec *SalesOrderStatusStates, e *SalesOrderStatusState) error {
	if len(dbm.option.eventSalesOrderStatusStateFunc[ef]) == 0 || skipEvents {
		return nil
	}
	for _, fn := range dbm.option.eventSalesOrderStatusStateFunc[ef] {
		if err := fn(ctx, ec, e); err != nil {
			return errors.WithStack(err)
		}
	}
	return nil
}

func (dbm DBM) eventViewCustomerAutoIncrementFunc(ctx context.Context, ef dml.EventFlag, skipEvents bool, ec *ViewCustomerAutoIncrements, e *ViewCustomerAutoIncrement) error {
	if len(dbm.option.eventViewCustomerAutoIncrementFunc[ef]) == 0 || skipEvents {
		return nil
	}
	for _, fn := range dbm.option.eventViewCustomerAutoIncrementFunc[ef] {
		if err := fn(ctx, ec, e); err != nil {
			return errors.WithStack(err)
		}
	}
	return nil
}

// NewDBManager returns a goified version of the MySQL/MariaDB table schema for
// the tables:  core_configuration, sales_order_status_state,
// view_customer_auto_increment Auto generated by dmlgen.
func NewDBManager(ctx context.Context, dbmo *DBMOption) (*DBM, error) {
	tbls, err := ddl.NewTables(append([]ddl.TableOption{ddl.WithCreateTable(ctx, TableNameCoreConfiguration, "", TableNameSalesOrderStatusState, "", TableNameViewCustomerAutoIncrement, "")}, dbmo.TableOptions...)...)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	if dbmo.InitSelectFn == nil {
		dbmo.InitSelectFn = func(s *dml.Select) *dml.Select { return s }
	}
	if dbmo.InitUpdateFn == nil {
		dbmo.InitUpdateFn = func(s *dml.Update) *dml.Update { return s }
	}
	if dbmo.InitDeleteFn == nil {
		dbmo.InitDeleteFn = func(s *dml.Delete) *dml.Delete { return s }
	}
	if dbmo.InitInsertFn == nil {
		dbmo.InitInsertFn = func(s *dml.Insert) *dml.Insert { return s }
	}
	err = tbls.Options(
		ddl.WithQueryDBR(map[string]dml.QueryBuilder{
			"CoreConfigurationsSelectAll": dbmo.InitSelectFn(tbls.MustTable(TableNameCoreConfiguration).Select("*")),
			"CoreConfigurationsSelectByPK": dbmo.InitSelectFn(tbls.MustTable(TableNameCoreConfiguration).Select("*")).Where(
				dml.Column(`config_id`).In().PlaceHolder(),
			),
			"CoreConfigurationSelectByPK": dbmo.InitSelectFn(tbls.MustTable(TableNameCoreConfiguration).Select("*")).Where(
				dml.Column(`config_id`).Equal().PlaceHolder(),
			),
			"CoreConfigurationUpdateByPK": dbmo.InitUpdateFn(tbls.MustTable(TableNameCoreConfiguration).Update().Where(
				dml.Column(`config_id`).Equal().PlaceHolder(),
			)),
			"CoreConfigurationDeleteByPK": dbmo.InitDeleteFn(tbls.MustTable(TableNameCoreConfiguration).Delete().Where(
				dml.Column(`config_id`).In().PlaceHolder(),
			)),
			"CoreConfigurationInsert":         dbmo.InitInsertFn(tbls.MustTable(TableNameCoreConfiguration).Insert()),
			"CoreConfigurationUpsertByPK":     dbmo.InitInsertFn(tbls.MustTable(TableNameCoreConfiguration).Insert()).OnDuplicateKey(),
			"SalesOrderStatusStatesSelectAll": dbmo.InitSelectFn(tbls.MustTable(TableNameSalesOrderStatusState).Select("*")),
			"SalesOrderStatusStatesSelectByPK": dbmo.InitSelectFn(tbls.MustTable(TableNameSalesOrderStatusState).Select("*")).Where(
				dml.Columns(`status`, `state`).In().Tuples(),
			),
			"SalesOrderStatusStateSelectByPK": dbmo.InitSelectFn(tbls.MustTable(TableNameSalesOrderStatusState).Select("*")).Where(
				dml.Columns(`status`, `state`).Equal().Tuples(),
			),
			"ViewCustomerAutoIncrementsSelectAll": dbmo.InitSelectFn(tbls.MustTable(TableNameViewCustomerAutoIncrement).Select("*")),
			"ViewCustomerAutoIncrementsSelectByPK": dbmo.InitSelectFn(tbls.MustTable(TableNameViewCustomerAutoIncrement).Select("*")).Where(
				dml.Column(`ce_entity_id`).In().PlaceHolder(),
			),
			"ViewCustomerAutoIncrementSelectByPK": dbmo.InitSelectFn(tbls.MustTable(TableNameViewCustomerAutoIncrement).Select("*")).Where(
				dml.Column(`ce_entity_id`).Equal().PlaceHolder(),
			),
		}),
	)
	if err != nil {
		return nil, err
	}
	if err := tbls.Options(dbmo.TableOptionsAfter...); err != nil {
		return nil, err
	}
	return &DBM{Tables: tbls, option: *dbmo}, nil
}

// CoreConfiguration represents a single row for DB table core_configuration.
// Auto generated.
// Table comment: Config Data
type CoreConfiguration struct {
	ConfigID  uint32      `max_len:"10"` // config_id int(10) unsigned NOT NULL PRI  auto_increment "Id"
	Scope     string      `max_len:"8"`  // scope varchar(8) NOT NULL MUL DEFAULT ''default''  "Scope"
	ScopeID   int32       `max_len:"10"` // scope_id int(11) NOT NULL  DEFAULT '0'  "Scope Id"
	Expires   null.Time   // expires datetime NULL  DEFAULT 'NULL'  "Value expiration time"
	Path      string      `max_len:"255"`   // path varchar(255) NOT NULL    "Path"
	Value     null.String `max_len:"65535"` // value text NULL  DEFAULT 'NULL'  "Value"
	VersionTs time.Time   // version_ts timestamp(6) NOT NULL   STORED GENERATED "Timestamp Start Versioning"
	VersionTe time.Time   // version_te timestamp(6) NOT NULL PRI  STORED GENERATED "Timestamp End Versioning"
}

// AssignLastInsertID updates the increment ID field with the last inserted ID
// from an INSERT operation. Implements dml.InsertIDAssigner. Auto generated.
func (e *CoreConfiguration) AssignLastInsertID(id int64) {
	e.ConfigID = uint32(id)
}

// MapColumns implements interface ColumnMapper only partially. Auto generated.
func (e *CoreConfiguration) MapColumns(cm *dml.ColumnMap) error {
	for cm.Next(8) {
		switch c := cm.Column(); c {
		case "config_id", "0":
			cm.Uint32(&e.ConfigID)
		case "scope", "1":
			cm.String(&e.Scope)
		case "scope_id", "2":
			cm.Int32(&e.ScopeID)
		case "expires", "3":
			cm.NullTime(&e.Expires)
		case "path", "4":
			cm.String(&e.Path)
		case "value", "5":
			cm.NullString(&e.Value)
		case "version_ts", "6":
			cm.Time(&e.VersionTs)
		case "version_te", "7":
			cm.Time(&e.VersionTe)
		default:
			return errors.NotFound.Newf("[dmltestgenerated4] CoreConfiguration Column %q not found", c)
		}
	}
	return errors.WithStack(cm.Err())
}

func (e *CoreConfiguration) Load(ctx context.Context, dbm *DBM, primaryKey uint32, opts ...dml.DBRFunc) (err error) {
	if e == nil {
		return errors.NotValid.Newf("CoreConfiguration can't be nil")
	}
	qo := dml.FromContextQueryOptions(ctx)
	// put the IDs primaryKey into the context as value to search for a cache entry in the event function.
	if err = dbm.eventCoreConfigurationFunc(ctx, dml.EventFlagBeforeSelect, qo.SkipEvents, nil, e); err != nil {
		return errors.WithStack(err)
	}
	if e.IsSet() {
		return nil // might return data from cache
	}
	if _, err = dbm.ConnPool.WithCacheKey("CoreConfigurationSelectByPK", opts...).Load(ctx, e, primaryKey); err != nil {
		return errors.WithStack(err)
	}
	return errors.WithStack(dbm.eventCoreConfigurationFunc(ctx, dml.EventFlagAfterSelect, qo.SkipEvents, nil, e))
}

func (e *CoreConfiguration) Delete(ctx context.Context, dbm *DBM, opts ...dml.DBRFunc) (res sql.Result, err error) {
	if e == nil {
		return nil, errors.NotValid.Newf("CoreConfiguration can't be nil")
	}
	qo := dml.FromContextQueryOptions(ctx)
	if err = dbm.eventCoreConfigurationFunc(ctx, dml.EventFlagBeforeDelete, qo.SkipEvents, nil, e); err != nil {
		return nil, errors.WithStack(err)
	}
	if res, err = dbm.ConnPool.WithCacheKey("CoreConfigurationDeleteByPK", opts...).ExecContext(ctx, e.ConfigID); err != nil {
		return nil, errors.WithStack(err)
	}
	if err = dbm.eventCoreConfigurationFunc(ctx, dml.EventFlagAfterDelete, qo.SkipEvents, nil, e); err != nil {
		return nil, errors.WithStack(err)
	}
	return res, nil
}

func (e *CoreConfiguration) Update(ctx context.Context, dbm *DBM, opts ...dml.DBRFunc) (res sql.Result, err error) {
	if e == nil {
		return nil, errors.NotValid.Newf("CoreConfiguration can't be nil")
	}
	qo := dml.FromContextQueryOptions(ctx)
	if err = dbm.eventCoreConfigurationFunc(ctx, dml.EventFlagBeforeUpdate, qo.SkipEvents, nil, e); err != nil {
		return nil, errors.WithStack(err)
	}
	if res, err = dbm.ConnPool.WithCacheKey("CoreConfigurationUpdateByPK", opts...).ExecContext(ctx, e); err != nil {
		return nil, errors.WithStack(err)
	}
	if err = dbm.eventCoreConfigurationFunc(ctx, dml.EventFlagAfterUpdate, qo.SkipEvents, nil, e); err != nil {
		return nil, errors.WithStack(err)
	}
	return res, nil
}

func (e *CoreConfiguration) Insert(ctx context.Context, dbm *DBM, opts ...dml.DBRFunc) (res sql.Result, err error) {
	if e == nil {
		return nil, errors.NotValid.Newf("CoreConfiguration can't be nil")
	}
	qo := dml.FromContextQueryOptions(ctx)
	if err = dbm.eventCoreConfigurationFunc(ctx, dml.EventFlagBeforeInsert, qo.SkipEvents, nil, e); err != nil {
		return nil, errors.WithStack(err)
	}
	if res, err = dbm.ConnPool.WithCacheKey("CoreConfigurationInsert", opts...).ExecContext(ctx, e); err != nil {
		return nil, errors.WithStack(err)
	}
	if err = dbm.eventCoreConfigurationFunc(ctx, dml.EventFlagAfterInsert, qo.SkipEvents, nil, e); err != nil {
		return nil, errors.WithStack(err)
	}
	return res, nil
}

func (e *CoreConfiguration) Upsert(ctx context.Context, dbm *DBM, opts ...dml.DBRFunc) (res sql.Result, err error) {
	if e == nil {
		return nil, errors.NotValid.Newf("CoreConfiguration can't be nil")
	}
	qo := dml.FromContextQueryOptions(ctx)
	if err = dbm.eventCoreConfigurationFunc(ctx, dml.EventFlagBeforeUpsert, qo.SkipEvents, nil, e); err != nil {
		return nil, errors.WithStack(err)
	}
	if res, err = dbm.ConnPool.WithCacheKey("CoreConfigurationUpsertByPK", opts...).ExecContext(ctx, dml.Qualify("", e)); err != nil {
		return nil, errors.WithStack(err)
	}
	if err = dbm.eventCoreConfigurationFunc(ctx, dml.EventFlagAfterUpsert, qo.SkipEvents, nil, e); err != nil {
		return nil, errors.WithStack(err)
	}
	return res, nil
}

// IsSet returns true if the entity has non-empty primary keys.
func (e *CoreConfiguration) IsSet() bool { return e.ConfigID > 0 }

// CoreConfigurations represents a collection type for DB table
// core_configuration
// Not thread safe. Auto generated.
type CoreConfigurations struct {
	Data []*CoreConfiguration `json:"data,omitempty"`
}

// NewCoreConfigurations  creates a new initialized collection. Auto generated.
func NewCoreConfigurations() *CoreConfigurations {
	return &CoreConfigurations{
		Data: make([]*CoreConfiguration, 0, 5),
	}
}

// AssignLastInsertID traverses through the slice and sets an incrementing new ID
// to each entity.
func (cc *CoreConfigurations) AssignLastInsertID(id int64) {
	for i := 0; i < len(cc.Data); i++ {
		cc.Data[i].AssignLastInsertID(id + int64(i))
	}
}

func (cc *CoreConfigurations) scanColumns(cm *dml.ColumnMap, e *CoreConfiguration) error {
	if err := e.MapColumns(cm); err != nil {
		return errors.WithStack(err)
	}
	// this function might get extended.
	return nil
}

// MapColumns implements dml.ColumnMapper interface. Auto generated.
func (cc *CoreConfigurations) MapColumns(cm *dml.ColumnMap) error {
	switch m := cm.Mode(); m {
	case dml.ColumnMapEntityReadAll, dml.ColumnMapEntityReadSet:
		for _, e := range cc.Data {
			if err := cc.scanColumns(cm, e); err != nil {
				return errors.WithStack(err)
			}
		}
	case dml.ColumnMapScan:
		if cm.Count == 0 {
			cc.Data = cc.Data[:0]
		}
		var e CoreConfiguration
		if err := cc.scanColumns(cm, &e); err != nil {
			return errors.WithStack(err)
		}
		cc.Data = append(cc.Data, &e)
	case dml.ColumnMapCollectionReadSet:
		for cm.Next(0) {
			switch c := cm.Column(); c {
			case "config_id":
				cm = cm.Uint32s(cc.ConfigIDs()...)
			default:
				return errors.NotFound.Newf("[dmltestgenerated4] CoreConfigurations Column %q not found", c)
			}
		} // end for cm.Next
	default:
		return errors.NotSupported.Newf("[dmltestgenerated4] Unknown Mode: %q", string(m))
	}
	return cm.Err()
}

func (cc *CoreConfigurations) DBLoad(ctx context.Context, dbm *DBM, pkIDs []uint32, opts ...dml.DBRFunc) (err error) {
	if cc == nil {
		return errors.NotValid.Newf("CoreConfiguration can't be nil")
	}
	qo := dml.FromContextQueryOptions(ctx)
	// put the IDs ConfigID into the context as value to search for a cache entry in the event function.
	if err = dbm.eventCoreConfigurationFunc(ctx, dml.EventFlagBeforeSelect, qo.SkipEvents, cc, nil); err != nil {
		return errors.WithStack(err)
	}
	if cc.Data != nil {
		return nil // might return data from cache
	}
	if len(pkIDs) > 0 {
		if _, err = dbm.ConnPool.WithCacheKey("CoreConfigurationsSelectByPK", opts...).Load(ctx, cc, pkIDs); err != nil {
			return errors.WithStack(err)
		}
	} else {
		if _, err = dbm.ConnPool.WithCacheKey("CoreConfigurationsSelectAll", opts...).Load(ctx, cc); err != nil {
			return errors.WithStack(err)
		}
	}
	return errors.WithStack(dbm.eventCoreConfigurationFunc(ctx, dml.EventFlagAfterSelect, qo.SkipEvents, cc, nil))
}

func (cc *CoreConfigurations) DBDelete(ctx context.Context, dbm *DBM, opts ...dml.DBRFunc) (res sql.Result, err error) {
	if cc == nil {
		return nil, errors.NotValid.Newf("CoreConfigurations can't be nil")
	}
	qo := dml.FromContextQueryOptions(ctx)
	if err = dbm.eventCoreConfigurationFunc(ctx, dml.EventFlagBeforeDelete, qo.SkipEvents, cc, nil); err != nil {
		return nil, errors.WithStack(err)
	}
	if res, err = dbm.ConnPool.WithCacheKey("CoreConfigurationDeleteByPK", opts...).ExecContext(ctx, dml.Qualify("", cc)); err != nil {
		return nil, errors.WithStack(err)
	}
	if err = errors.WithStack(dbm.eventCoreConfigurationFunc(ctx, dml.EventFlagAfterDelete, qo.SkipEvents, cc, nil)); err != nil {
		return nil, errors.WithStack(err)
	}
	return res, nil
}

func (cc *CoreConfigurations) DBUpdate(ctx context.Context, dbm *DBM, resCheckFn func(sql.Result, error) error, opts ...dml.DBRFunc) (err error) {
	if cc == nil {
		return errors.NotValid.Newf("CoreConfigurations can't be nil")
	}
	qo := dml.FromContextQueryOptions(ctx)
	if err = dbm.eventCoreConfigurationFunc(ctx, dml.EventFlagBeforeUpdate, qo.SkipEvents, cc, nil); err != nil {
		return errors.WithStack(err)
	}
	if len(opts) == 0 {
		opts = dbmEmptyOpts
	}
	if resCheckFn == nil {
		resCheckFn = dbmNoopResultCheckFn
	}
	dbrStmt, err := dbm.ConnPool.WithCacheKey("CoreConfigurationUpdateByPK", opts...).Prepare(ctx)
	if err != nil {
		return errors.WithStack(err)
	}
	for _, c := range cc.Data {
		if err := resCheckFn(dbrStmt.ExecContext(ctx, c)); err != nil {
			return errors.WithStack(err)
		}
	}
	return errors.WithStack(dbm.eventCoreConfigurationFunc(ctx, dml.EventFlagAfterUpdate, qo.SkipEvents, cc, nil))
}

func (cc *CoreConfigurations) DBInsert(ctx context.Context, dbm *DBM, opts ...dml.DBRFunc) (res sql.Result, err error) {
	if cc == nil {
		return nil, errors.NotValid.Newf("CoreConfigurations can't be nil")
	}
	qo := dml.FromContextQueryOptions(ctx)
	if err = dbm.eventCoreConfigurationFunc(ctx, dml.EventFlagBeforeInsert, qo.SkipEvents, cc, nil); err != nil {
		return nil, errors.WithStack(err)
	}
	if res, err = dbm.ConnPool.WithCacheKey("CoreConfigurationInsert", opts...).ExecContext(ctx, cc); err != nil {
		return nil, errors.WithStack(err)
	}
	if err = errors.WithStack(dbm.eventCoreConfigurationFunc(ctx, dml.EventFlagAfterInsert, qo.SkipEvents, cc, nil)); err != nil {
		return nil, errors.WithStack(err)
	}
	return res, nil
}

func (cc *CoreConfigurations) DBUpsert(ctx context.Context, dbm *DBM, opts ...dml.DBRFunc) (res sql.Result, err error) {
	if cc == nil {
		return nil, errors.NotValid.Newf("CoreConfigurations can't be nil")
	}
	qo := dml.FromContextQueryOptions(ctx)
	if err = dbm.eventCoreConfigurationFunc(ctx, dml.EventFlagBeforeUpsert, qo.SkipEvents, cc, nil); err != nil {
		return nil, errors.WithStack(err)
	}
	if res, err = dbm.ConnPool.WithCacheKey("CoreConfigurationUpsertByPK", opts...).ExecContext(ctx, dml.Qualify("", cc)); err != nil {
		return nil, errors.WithStack(err)
	}
	if err = dbm.eventCoreConfigurationFunc(ctx, dml.EventFlagAfterUpsert, qo.SkipEvents, cc, nil); err != nil {
		return nil, errors.WithStack(err)
	}
	return res, nil
}

// Each will run function f on all items in []* CoreConfiguration . Auto
// generated via dmlgen.
func (cc *CoreConfigurations) Each(f func(*CoreConfiguration)) *CoreConfigurations {
	if cc == nil {
		return nil
	}
	for i := range cc.Data {
		f(cc.Data[i])
	}
	return cc
}

// ConfigIDs returns a slice with the data or appends it to a slice.
// Auto generated.
func (cc *CoreConfigurations) ConfigIDs(ret ...uint32) []uint32 {
	if cc == nil {
		return nil
	}
	if ret == nil {
		ret = make([]uint32, 0, len(cc.Data))
	}
	for _, e := range cc.Data {
		ret = append(ret, e.ConfigID)
	}
	return ret
}

// SalesOrderStatusState represents a single row for DB table
// sales_order_status_state. Auto generated.
// Table comment: Sales Order Status Table
type SalesOrderStatusState struct {
	Status         string // status varchar(32) NOT NULL PRI   "Status"
	State          string // state varchar(32) NOT NULL PRI   "Label"
	IsDefault      bool   // is_default smallint(5) unsigned NOT NULL  DEFAULT '0'  "Is Default"
	VisibleOnFront uint16 // visible_on_front smallint(5) unsigned NOT NULL  DEFAULT '0'  "Visible on front"
}

// MapColumns implements interface ColumnMapper only partially. Auto generated.
func (e *SalesOrderStatusState) MapColumns(cm *dml.ColumnMap) error {
	for cm.Next(4) {
		switch c := cm.Column(); c {
		case "status", "0":
			cm.String(&e.Status)
		case "state", "1":
			cm.String(&e.State)
		case "is_default", "2":
			cm.Bool(&e.IsDefault)
		case "visible_on_front", "3":
			cm.Uint16(&e.VisibleOnFront)
		default:
			return errors.NotFound.Newf("[dmltestgenerated4] SalesOrderStatusState Column %q not found", c)
		}
	}
	return errors.WithStack(cm.Err())
}

type SalesOrderStatusStateLoadArgs struct {
	_Named_Fields_Required struct{}
	Status                 string
	State                  string
}

func (e *SalesOrderStatusState) Load(ctx context.Context, dbm *DBM, arg SalesOrderStatusStateLoadArgs, opts ...dml.DBRFunc) (err error) {
	if e == nil {
		return errors.NotValid.Newf("SalesOrderStatusState can't be nil")
	}
	qo := dml.FromContextQueryOptions(ctx)
	// put the IDs arg.Status,arg.State into the context as value to search for a cache entry in the event function.
	if err = dbm.eventSalesOrderStatusStateFunc(ctx, dml.EventFlagBeforeSelect, qo.SkipEvents, nil, e); err != nil {
		return errors.WithStack(err)
	}
	if e.IsSet() {
		return nil // might return data from cache
	}
	if _, err = dbm.ConnPool.WithCacheKey("SalesOrderStatusStateSelectByPK", opts...).Load(ctx, e, arg.Status, arg.State); err != nil {
		return errors.WithStack(err)
	}
	return errors.WithStack(dbm.eventSalesOrderStatusStateFunc(ctx, dml.EventFlagAfterSelect, qo.SkipEvents, nil, e))
}

// IsSet returns true if the entity has non-empty primary keys.
func (e *SalesOrderStatusState) IsSet() bool { return e.Status != "" && e.State != "" }

// SalesOrderStatusStates represents a collection type for DB table
// sales_order_status_state
// Not thread safe. Auto generated.
type SalesOrderStatusStates struct {
	Data []*SalesOrderStatusState `json:"data,omitempty"`
}

// NewSalesOrderStatusStates  creates a new initialized collection. Auto
// generated.
func NewSalesOrderStatusStates() *SalesOrderStatusStates {
	return &SalesOrderStatusStates{
		Data: make([]*SalesOrderStatusState, 0, 5),
	}
}

func (cc *SalesOrderStatusStates) scanColumns(cm *dml.ColumnMap, e *SalesOrderStatusState) error {
	if err := e.MapColumns(cm); err != nil {
		return errors.WithStack(err)
	}
	// this function might get extended.
	return nil
}

// MapColumns implements dml.ColumnMapper interface. Auto generated.
func (cc *SalesOrderStatusStates) MapColumns(cm *dml.ColumnMap) error {
	switch m := cm.Mode(); m {
	case dml.ColumnMapEntityReadAll, dml.ColumnMapEntityReadSet:
		for _, e := range cc.Data {
			if err := cc.scanColumns(cm, e); err != nil {
				return errors.WithStack(err)
			}
		}
	case dml.ColumnMapScan:
		if cm.Count == 0 {
			cc.Data = cc.Data[:0]
		}
		var e SalesOrderStatusState
		if err := cc.scanColumns(cm, &e); err != nil {
			return errors.WithStack(err)
		}
		cc.Data = append(cc.Data, &e)
	case dml.ColumnMapCollectionReadSet:
		for cm.Next(0) {
			switch c := cm.Column(); c {
			case "status":
				cm = cm.Strings(cc.Statuss()...)
			case "state":
				cm = cm.Strings(cc.States()...)
			default:
				return errors.NotFound.Newf("[dmltestgenerated4] SalesOrderStatusStates Column %q not found", c)
			}
		} // end for cm.Next
	default:
		return errors.NotSupported.Newf("[dmltestgenerated4] Unknown Mode: %q", string(m))
	}
	return cm.Err()
}

type SalesOrderStatusStatesDBLoadArgs struct {
	_Named_Fields_Required struct{}
	Status                 string
	State                  string
}

func (cc *SalesOrderStatusStates) DBLoad(ctx context.Context, dbm *DBM, pkIDs []SalesOrderStatusStatesDBLoadArgs, opts ...dml.DBRFunc) (err error) {
	if cc == nil {
		return errors.NotValid.Newf("SalesOrderStatusState can't be nil")
	}
	qo := dml.FromContextQueryOptions(ctx)
	// put the IDs Status,State into the context as value to search for a cache entry in the event function.
	if err = dbm.eventSalesOrderStatusStateFunc(ctx, dml.EventFlagBeforeSelect, qo.SkipEvents, cc, nil); err != nil {
		return errors.WithStack(err)
	}
	if cc.Data != nil {
		return nil // might return data from cache
	}
	cacheKey := "SalesOrderStatusStatesSelectAll"
	var args []any
	if len(pkIDs) > 0 {
		args = make([]any, 0, len(pkIDs)*2)
		for _, pk := range pkIDs {
			args = append(args, pk.Status)
			args = append(args, pk.State)
		}
		cacheKey = "SalesOrderStatusStatesSelectByPK"
	}
	if _, err = dbm.ConnPool.WithCacheKey(cacheKey, opts...).Load(ctx, cc, args...); err != nil {
		return errors.WithStack(err)
	}
	return errors.WithStack(dbm.eventSalesOrderStatusStateFunc(ctx, dml.EventFlagAfterSelect, qo.SkipEvents, cc, nil))
}

// Each will run function f on all items in []* SalesOrderStatusState . Auto
// generated via dmlgen.
func (cc *SalesOrderStatusStates) Each(f func(*SalesOrderStatusState)) *SalesOrderStatusStates {
	if cc == nil {
		return nil
	}
	for i := range cc.Data {
		f(cc.Data[i])
	}
	return cc
}

// Statuss returns a slice with the data or appends it to a slice.
// Auto generated.
func (cc *SalesOrderStatusStates) Statuss(ret ...string) []string {
	if cc == nil {
		return nil
	}
	if ret == nil {
		ret = make([]string, 0, len(cc.Data))
	}
	for _, e := range cc.Data {
		ret = append(ret, e.Status)
	}
	return ret
}

// States returns a slice with the data or appends it to a slice.
// Auto generated.
func (cc *SalesOrderStatusStates) States(ret ...string) []string {
	if cc == nil {
		return nil
	}
	if ret == nil {
		ret = make([]string, 0, len(cc.Data))
	}
	for _, e := range cc.Data {
		ret = append(ret, e.State)
	}
	return ret
}

// ViewCustomerAutoIncrement represents a single row for DB table
// view_customer_auto_increment. Auto generated.
// Table comment: VIEW
type ViewCustomerAutoIncrement struct {
	CeEntityID uint32      // ce_entity_id int(10) unsigned NOT NULL  DEFAULT '0'  "Entity ID"
	Email      null.String // email varchar(255) NULL  DEFAULT 'NULL'  "Email"
	Firstname  string      // firstname varchar(255) NOT NULL    "First Name"
	Lastname   string      // lastname varchar(255) NOT NULL    "Last Name"
	City       string      // city varchar(255) NOT NULL    "City"
}

// MapColumns implements interface ColumnMapper only partially. Auto generated.
func (e *ViewCustomerAutoIncrement) MapColumns(cm *dml.ColumnMap) error {
	for cm.Next(5) {
		switch c := cm.Column(); c {
		case "ce_entity_id", "0":
			cm.Uint32(&e.CeEntityID)
		case "email", "1":
			cm.NullString(&e.Email)
		case "firstname", "2":
			cm.String(&e.Firstname)
		case "lastname", "3":
			cm.String(&e.Lastname)
		case "city", "4":
			cm.String(&e.City)
		default:
			return errors.NotFound.Newf("[dmltestgenerated4] ViewCustomerAutoIncrement Column %q not found", c)
		}
	}
	return errors.WithStack(cm.Err())
}

func (e *ViewCustomerAutoIncrement) Load(ctx context.Context, dbm *DBM, primaryKey uint32, opts ...dml.DBRFunc) (err error) {
	if e == nil {
		return errors.NotValid.Newf("ViewCustomerAutoIncrement can't be nil")
	}
	qo := dml.FromContextQueryOptions(ctx)
	// put the IDs primaryKey into the context as value to search for a cache entry in the event function.
	if err = dbm.eventViewCustomerAutoIncrementFunc(ctx, dml.EventFlagBeforeSelect, qo.SkipEvents, nil, e); err != nil {
		return errors.WithStack(err)
	}
	if e.IsSet() {
		return nil // might return data from cache
	}
	if _, err = dbm.ConnPool.WithCacheKey("ViewCustomerAutoIncrementSelectByPK", opts...).Load(ctx, e, primaryKey); err != nil {
		return errors.WithStack(err)
	}
	return errors.WithStack(dbm.eventViewCustomerAutoIncrementFunc(ctx, dml.EventFlagAfterSelect, qo.SkipEvents, nil, e))
}

// IsSet returns true if the entity has non-empty primary keys.
func (e *ViewCustomerAutoIncrement) IsSet() bool { return e.CeEntityID > 0 }

// ViewCustomerAutoIncrements represents a collection type for DB table
// view_customer_auto_increment
// Not thread safe. Auto generated.
type ViewCustomerAutoIncrements struct {
	Data []*ViewCustomerAutoIncrement `json:"data,omitempty"`
}

// NewViewCustomerAutoIncrements  creates a new initialized collection. Auto
// generated.
func NewViewCustomerAutoIncrements() *ViewCustomerAutoIncrements {
	return &ViewCustomerAutoIncrements{
		Data: make([]*ViewCustomerAutoIncrement, 0, 5),
	}
}

func (cc *ViewCustomerAutoIncrements) scanColumns(cm *dml.ColumnMap, e *ViewCustomerAutoIncrement) error {
	if err := e.MapColumns(cm); err != nil {
		return errors.WithStack(err)
	}
	// this function might get extended.
	return nil
}

// MapColumns implements dml.ColumnMapper interface. Auto generated.
func (cc *ViewCustomerAutoIncrements) MapColumns(cm *dml.ColumnMap) error {
	switch m := cm.Mode(); m {
	case dml.ColumnMapEntityReadAll, dml.ColumnMapEntityReadSet:
		for _, e := range cc.Data {
			if err := cc.scanColumns(cm, e); err != nil {
				return errors.WithStack(err)
			}
		}
	case dml.ColumnMapScan:
		if cm.Count == 0 {
			cc.Data = cc.Data[:0]
		}
		var e ViewCustomerAutoIncrement
		if err := cc.scanColumns(cm, &e); err != nil {
			return errors.WithStack(err)
		}
		cc.Data = append(cc.Data, &e)
	default:
		return errors.NotSupported.Newf("[dmltestgenerated4] Unknown Mode: %q", string(m))
	}
	return cm.Err()
}

func (cc *ViewCustomerAutoIncrements) DBLoad(ctx context.Context, dbm *DBM, pkIDs []uint32, opts ...dml.DBRFunc) (err error) {
	if cc == nil {
		return errors.NotValid.Newf("ViewCustomerAutoIncrement can't be nil")
	}
	qo := dml.FromContextQueryOptions(ctx)
	// put the IDs CeEntityID into the context as value to search for a cache entry in the event function.
	if err = dbm.eventViewCustomerAutoIncrementFunc(ctx, dml.EventFlagBeforeSelect, qo.SkipEvents, cc, nil); err != nil {
		return errors.WithStack(err)
	}
	if cc.Data != nil {
		return nil // might return data from cache
	}
	if len(pkIDs) > 0 {
		if _, err = dbm.ConnPool.WithCacheKey("ViewCustomerAutoIncrementsSelectByPK", opts...).Load(ctx, cc, pkIDs); err != nil {
			return errors.WithStack(err)
		}
	} else {
		if _, err = dbm.ConnPool.WithCacheKey("ViewCustomerAutoIncrementsSelectAll", opts...).Load(ctx, cc); err != nil {
			return errors.WithStack(err)
		}
	}
	return errors.WithStack(dbm.eventViewCustomerAutoIncrementFunc(ctx, dml.EventFlagAfterSelect, qo.SkipEvents, cc, nil))
}

// Each will run function f on all items in []* ViewCustomerAutoIncrement . Auto
// generated via dmlgen.
func (cc *ViewCustomerAutoIncrements) Each(f func(*ViewCustomerAutoIncrement)) *ViewCustomerAutoIncrements {
	if cc == nil {
		return nil
	}
	for i := range cc.Data {
		f(cc.Data[i])
	}
	return cc
}
