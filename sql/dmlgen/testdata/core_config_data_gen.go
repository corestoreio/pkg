// Auto generated via github.com/corestoreio/pkg/sql/dmlgen

package testdata

import (
	"github.com/corestoreio/pkg/sql/dml"
	"github.com/corestoreio/errors"

)
// CoreConfigData represents a single row for DB table `core_config_data`.
// Auto generated.
type CoreConfigData struct {
	ConfigID uint64         `json:"config_id,omitempty"`     // config_id int(10) unsigned NOT NULL PRI  auto_increment "Config Id"
	Scope    string         `json:"scope,omitempty"`         // scope varchar(8) NOT NULL MUL DEFAULT ''default''  "Config Scope"
	ScopeID  int64          `json:"scope_id" xml:"scope_id"` // scope_id int(11) NOT NULL  DEFAULT '0'  "Config Scope Id"
	Path     string         `json:"x_path" xml:"y_path"`     // path varchar(255) NOT NULL  DEFAULT ''general''  "Config Path"
	Value    dml.NullString `json:"value,omitempty"`         // value text NULL  DEFAULT 'NULL'  "Config Value"
}

// NewCoreConfigData creates a new pointer with pre-initialized fields. Auto
// generated.
func NewCoreConfigData() *CoreConfigData {
	return &CoreConfigData{}
}

// AssignLastInsertID updates the increment ID field with the last inserted ID
// from an INSERT operation. Implements dml.InsertIDAssigner. Auto generated.
func (e *CoreConfigData) AssignLastInsertID(id int64) {
	e.ConfigID = uint64(id)
}

// MapColumns implements interface ColumnMapper only partially. Auto generated.
func (e *CoreConfigData) MapColumns(cm *dml.ColumnMap) error {
	if cm.Mode() == dml.ColumnMapEntityReadAll {
		return cm.Uint64(&e.ConfigID).String(&e.Scope).Int64(&e.ScopeID).String(&e.Path).NullString(&e.Value).Err()
	}
	for cm.Next() {
		switch c := cm.Column(); c {
		case "config_id":
			cm.Uint64(&e.ConfigID)
		case "scope":
			cm.String(&e.Scope)
		case "scope_id":
			cm.Int64(&e.ScopeID)
		case "path", "storage_location", "config_directory":
			cm.String(&e.Path)
		case "value":
			cm.NullString(&e.Value)
		default:
			return errors.NewNotFoundf("[testdata] CoreConfigData Column %q not found", c)
		}
	}
	return errors.WithStack(cm.Err())
}

// CoreConfigDataCollection represents a collection type for DB table core_config_data
// Not thread safe. Auto generated.
type CoreConfigDataCollection struct {
	Data             []*CoreConfigData
	BeforeMapColumns func(uint64, *CoreConfigData) error
	AfterMapColumns  func(uint64, *CoreConfigData) error
}

// MakeCoreConfigDataCollection creates a new initialized collection. Auto generated.
func MakeCoreConfigDataCollection() CoreConfigDataCollection {
	return CoreConfigDataCollection{
		Data: make([]*CoreConfigData, 0, 5),
	}
}

func (cc CoreConfigDataCollection) scanColumns(cm *dml.ColumnMap, e *CoreConfigData, idx uint64) error {
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
func (cc CoreConfigDataCollection) MapColumns(cm *dml.ColumnMap) error {
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
		e := NewCoreConfigData()
		if err := cc.scanColumns(cm, e, cm.Count); err != nil {
			return errors.WithStack(err)
		}
		cc.Data = append(cc.Data, e)
	case dml.ColumnMapCollectionReadSet:
		for cm.Next() {
			switch c := cm.Column(); c {
			case "config_id":
				cm.Args = cm.Args.Uint64s(cc.ConfigIDs()...)
			case "path", "storage_location", "config_directory":
				cm.Args = cm.Args.Strings(cc.Paths()...)
			default:
				return errors.NewNotFoundf("[testdata] CoreConfigDataCollection Column %q not found", c)
			}
		}
	default:
		return errors.NewNotSupportedf("[dml] Unknown Mode: %q", string(m))
	}
	return cm.Err()
}

// ConfigIDs returns a slice or appends to a slice all values.
// Auto generated.
func (cc CoreConfigDataCollection) ConfigIDs(ret ...uint64) []uint64 {
	if ret == nil {
		ret = make([]uint64, 0, len(cc.Data))
	}
	for _, e := range cc.Data {
		ret = append(ret, e.ConfigID)
	}
	return ret
}

// Paths belongs to the column `path`
// and returns a slice or appends to a slice only unique values of that column.
// The values will be filtered internally in a Go map. No DB query gets
// executed. Auto generated.
func (cc CoreConfigDataCollection) Paths(ret ...string) []string {
	if ret == nil {
		ret = make([]string, 0, len(cc.Data))
	}

	dupCheck := make(map[string]struct{}, len(cc.Data))
	for _, e := range cc.Data {
		if _, ok := dupCheck[e.Path]; !ok {
			ret = append(ret, e.Path)
			dupCheck[e.Path] = struct{}{}
		}
	}
	return ret
} 
