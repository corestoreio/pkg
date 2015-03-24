package tools

import (
	"database/sql"

	"github.com/corestoreio/csfw/storage/csdb"
	"github.com/juju/errgo"
)

// Depends on generated code from tableToStruct.
type AddAttrTables struct {
	EntityTypeCode string
	db             *sql.DB
}

// Implements interface eav.EntityTypeAdditionalAttributeTabler
func NewAddAttrTables(db *sql.DB, code string) *AddAttrTables {
	return &AddAttrTables{
		EntityTypeCode: code,
		db:             db,
	}
}

// Implements interface eav.EntityTypeAdditionalAttributeTabler
func (aa *AddAttrTables) TableAdditionalAttribute() (*csdb.TableStructure, error) {
	if t, ok := ConfigEntityType[aa.EntityTypeCode]; ok {
		if t.TempAdditionalAttributeTable != "" {
			return aa.newTableStructure(t.TempAdditionalAttributeTable)
		}
		return nil, nil
	}
	return nil, errgo.Newf("Table for %s not found", aa.EntityTypeCode)
}

// Implements interface eav.EntityTypeAdditionalAttributeTabler
func (aa *AddAttrTables) TableEavWebsite() (*csdb.TableStructure, error) {
	if t, ok := ConfigEntityType[aa.EntityTypeCode]; ok {
		if t.TempAdditionalAttributeTableWebsite != "" {
			return aa.newTableStructure(t.TempAdditionalAttributeTableWebsite)
		}
		return nil, nil
	}
	return nil, errgo.Newf("Table for %s not found", aa.EntityTypeCode)
}

func (aa *AddAttrTables) newTableStructure(tableName string) (*csdb.TableStructure, error) {
	cols, err := GetColumns(aa.db, tableName)
	if err != nil {
		return nil, errgo.Mask(err)
	}
	return csdb.NewTableStructure(tableName, cols.GetFieldNames(true), cols.GetFieldNames(false)), nil
}
