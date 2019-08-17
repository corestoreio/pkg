// Copyright 2015-present, Cyrill @ Schumacher.fm and the CoreStore contributors
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

package ddl

import (
	"context"
	"database/sql"
	"fmt"
	"io"
	"sort"

	"github.com/corestoreio/errors"
	"github.com/corestoreio/pkg/sql/dml"
	"github.com/corestoreio/pkg/storage/null"
)

// KeyColumnUsage represents a single row for DB table `KEY_COLUMN_USAGE`
type KeyColumnUsage struct {
	ConstraintCatalog          string      // CONSTRAINT_CATALOG varchar(512) NOT NULL  DEFAULT ''''  ""
	ConstraintSchema           string      // CONSTRAINT_SCHEMA varchar(64) NOT NULL  DEFAULT ''''  ""
	ConstraintName             string      // CONSTRAINT_NAME varchar(64) NOT NULL  DEFAULT ''''  ""
	TableCatalog               string      // TABLE_CATALOG varchar(512) NOT NULL  DEFAULT ''''  ""
	TableSchema                string      // TABLE_SCHEMA varchar(64) NOT NULL  DEFAULT ''''  ""
	TableName                  string      // TABLE_NAME varchar(64) NOT NULL  DEFAULT ''''  ""
	ColumnName                 string      // COLUMN_NAME varchar(64) NOT NULL  DEFAULT ''''  ""
	OrdinalPosition            int64       // ORDINAL_POSITION bigint(10) NOT NULL  DEFAULT '0'  ""
	PositionInUniqueConstraint null.Int64  // POSITION_IN_UNIQUE_CONSTRAINT bigint(10) NULL  DEFAULT 'NULL'  ""
	ReferencedTableSchema      null.String // REFERENCED_TABLE_SCHEMA varchar(64) NULL  DEFAULT 'NULL'  ""
	ReferencedTableName        null.String // REFERENCED_TABLE_NAME varchar(64) NULL  DEFAULT 'NULL'  ""
	ReferencedColumnName       null.String // REFERENCED_COLUMN_NAME varchar(64) NULL  DEFAULT 'NULL'  ""
}

// MapColumns implements interface ColumnMapper only partially.
func (e *KeyColumnUsage) MapColumns(cm *dml.ColumnMap) error {
	if cm.Mode() == dml.ColumnMapEntityReadAll {
		return cm.String(&e.ConstraintCatalog).String(&e.ConstraintSchema).String(&e.ConstraintName).String(&e.TableCatalog).String(&e.TableSchema).String(&e.TableName).String(&e.ColumnName).Int64(&e.OrdinalPosition).NullInt64(&e.PositionInUniqueConstraint).NullString(&e.ReferencedTableSchema).NullString(&e.ReferencedTableName).NullString(&e.ReferencedColumnName).Err()
	}
	for cm.Next() {
		switch c := cm.Column(); c {
		case "CONSTRAINT_CATALOG":
			cm.String(&e.ConstraintCatalog)
		case "CONSTRAINT_SCHEMA":
			cm.String(&e.ConstraintSchema)
		case "CONSTRAINT_NAME":
			cm.String(&e.ConstraintName)
		case "TABLE_CATALOG":
			cm.String(&e.TableCatalog)
		case "TABLE_SCHEMA":
			cm.String(&e.TableSchema)
		case "TABLE_NAME":
			cm.String(&e.TableName)
		case "COLUMN_NAME":
			cm.String(&e.ColumnName)
		case "ORDINAL_POSITION":
			cm.Int64(&e.OrdinalPosition)
		case "POSITION_IN_UNIQUE_CONSTRAINT":
			cm.NullInt64(&e.PositionInUniqueConstraint)
		case "REFERENCED_TABLE_SCHEMA":
			cm.NullString(&e.ReferencedTableSchema)
		case "REFERENCED_TABLE_NAME":
			cm.NullString(&e.ReferencedTableName)
		case "REFERENCED_COLUMN_NAME":
			cm.NullString(&e.ReferencedColumnName)
		default:
			return errors.NotFound.Newf("[testdata] KeyColumnUsage Column %q not found", c)
		}
	}
	return errors.WithStack(cm.Err())
}

// Reset resets the struct to its empty fields.
func (e *KeyColumnUsage) Reset() *KeyColumnUsage {
	*e = KeyColumnUsage{}
	return e
}

// KeyColumnUsageCollection represents a collection type for DB table KEY_COLUMN_USAGE
// Not thread safe. Generated via dmlgen.
type KeyColumnUsageCollection struct {
	Data             []*KeyColumnUsage
	BeforeMapColumns func(uint64, *KeyColumnUsage) error
	AfterMapColumns  func(uint64, *KeyColumnUsage) error
}

// Sort sorts the collection by constraint name.
func (cc KeyColumnUsageCollection) Sort() {
	sort.Slice(cc.Data, func(i, j int) bool {
		return cc.Data[i].ConstraintName < cc.Data[j].ConstraintName
	})
}

func (cc KeyColumnUsageCollection) scanColumns(cm *dml.ColumnMap, e *KeyColumnUsage, idx uint64) error {
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

// MapColumns implements dml.ColumnMapper interface
func (cc KeyColumnUsageCollection) MapColumns(cm *dml.ColumnMap) error {
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
		e := new(KeyColumnUsage)
		if err := cc.scanColumns(cm, e, cm.Count); err != nil {
			return errors.WithStack(err)
		}
		cc.Data = append(cc.Data, e)
	case dml.ColumnMapCollectionReadSet:
		for cm.Next() {
			switch c := cm.Column(); c {
			case "TABLE_NAME":
				cm.Strings(cc.TableNames()...)
			case "COLUMN_NAME":
				cm.Strings(cc.ColumnNames()...)
			default:
				return errors.NotFound.Newf("[testdata] KeyColumnUsageCollection Column %q not found", c)
			}
		}
	default:
		return errors.NotSupported.Newf("[dml] Unknown Mode: %q", string(m))
	}
	return cm.Err()
}

// TableNames belongs to the column `TABLE_NAME` and returns a
// slice or appends to a slice only unique values of that column. The values
// will be filtered internally in a Go map. No DB query gets executed.
func (cc KeyColumnUsageCollection) TableNames(ret ...string) []string {
	if ret == nil {
		ret = make([]string, 0, len(cc.Data))
	}

	dupCheck := make(map[string]struct{}, len(cc.Data))
	for _, e := range cc.Data {
		if _, ok := dupCheck[e.TableName]; !ok {
			ret = append(ret, e.TableName)
			dupCheck[e.TableName] = struct{}{}
		}
	}
	return ret
}

// ColumnNames belongs to the column `COLUMN_NAME` and returns a
// slice or appends to a slice only unique values of that column. The values
// will be filtered internally in a Go map. No DB query gets executed.
func (cc KeyColumnUsageCollection) ColumnNames(ret ...string) []string {
	if ret == nil {
		ret = make([]string, 0, len(cc.Data))
	}

	dupCheck := make(map[string]struct{}, len(cc.Data))
	for _, e := range cc.Data {
		if _, ok := dupCheck[e.ColumnName]; !ok {
			ret = append(ret, e.ColumnName)
			dupCheck[e.ColumnName] = struct{}{}
		}
	}
	return ret
}

// LoadKeyColumnUsage returns all foreign key columns from a list of table names
// in the current database. Map key contains TABLE_NAME and value
// contains all of the table foreign keys. All columns from all tables gets
// selected when you don't provide the argument `tables`.
func LoadKeyColumnUsage(ctx context.Context, db dml.Querier, tables ...string) (tc map[string]KeyColumnUsageCollection, err error) {
	const selFkWhere = ` AND REFERENCED_TABLE_NAME IN ?`
	const selFkOrderBy = ` ORDER BY TABLE_SCHEMA,TABLE_NAME,ORDINAL_POSITION, COLUMN_NAME`

	const selFkTablesColumns = `SELECT
	CONSTRAINT_CATALOG, CONSTRAINT_SCHEMA, CONSTRAINT_NAME, TABLE_CATALOG, TABLE_SCHEMA,
	TABLE_NAME, COLUMN_NAME, ORDINAL_POSITION, POSITION_IN_UNIQUE_CONSTRAINT,
	REFERENCED_TABLE_SCHEMA, REFERENCED_TABLE_NAME, REFERENCED_COLUMN_NAME
	 FROM information_schema.KEY_COLUMN_USAGE WHERE REFERENCED_TABLE_SCHEMA = DATABASE()` + selFkWhere + selFkOrderBy

	const selFkAllTablesColumns = `SELECT
	CONSTRAINT_CATALOG, CONSTRAINT_SCHEMA, CONSTRAINT_NAME, TABLE_CATALOG, TABLE_SCHEMA,
	TABLE_NAME, COLUMN_NAME, ORDINAL_POSITION, POSITION_IN_UNIQUE_CONSTRAINT,
	REFERENCED_TABLE_SCHEMA, REFERENCED_TABLE_NAME, REFERENCED_COLUMN_NAME
	 FROM information_schema.KEY_COLUMN_USAGE WHERE REFERENCED_TABLE_SCHEMA = DATABASE()` + selFkOrderBy

	var rows *sql.Rows
	if len(tables) == 0 {
		rows, err = db.QueryContext(ctx, selFkAllTablesColumns)
		if err != nil {
			return nil, errors.Wrapf(err, "[ddl] LoadKeyColumnUsage QueryContext for tables %v", tables)
		}
	} else {
		sqlStr, _, err := dml.Interpolate(selFkTablesColumns).Strs(tables...).ToSQL()
		if err != nil {
			return nil, errors.Wrapf(err, "[ddl] LoadKeyColumnUsage dml.ExpandPlaceHolders for tables %v", tables)
		}
		rows, err = db.QueryContext(ctx, sqlStr)
		if err != nil {
			return nil, errors.Wrapf(err, "[ddl] LoadKeyColumnUsage QueryContext for tables %v with WHERE clause", tables)
		}
	}

	defer func() {
		// Not testable with the sqlmock package :-(
		if err2 := rows.Close(); err2 != nil && err == nil {
			err = errors.Wrap(err2, "[ddl] LoadKeyColumnUsage.Rows.Close")
		}
	}()

	tc = make(map[string]KeyColumnUsageCollection)
	rc := new(dml.ColumnMap)

	for rows.Next() {
		if err = rc.Scan(rows); err != nil {
			return nil, errors.Wrapf(err, "[ddl] LoadKeyColumnUsage Scan Query for tables: %v", tables) // due to the defer
		}
		kcu := new(KeyColumnUsage)
		if err = kcu.MapColumns(rc); err != nil {
			return nil, errors.WithStack(err)
		}
		if !kcu.ReferencedTableName.Valid || !kcu.ReferencedColumnName.Valid {
			err = errors.Fatal.Newf("[ddl] LoadKeyColumnUsage: The columns ReferencedTableName or ReferencedColumnName cannot be null: %#v", kcu)
			return
		}

		kcuc := tc[kcu.TableName]
		kcuc.Data = append(kcuc.Data, kcu)
		tc[kcu.TableName] = kcuc
	}
	if err = rows.Err(); err != nil {
		err = errors.WithStack(err)
	}
	return
}

// ReverseKeyColumnUsage reverses the argument to a new key column usage
// collection. E.g. customer_entity, catalog_product_entity and other tables
// have a foreign key to table store.store_id which is a OneToOne relationship.
// When reversed the table store, as map key, points to customer_entity and
// catalog_product_entity which becomes then a OneToMany relationship. If that
// makes sense is another topic.
func ReverseKeyColumnUsage(kcu map[string]KeyColumnUsageCollection) (kcuRev map[string]KeyColumnUsageCollection) {
	kcuRev = make(map[string]KeyColumnUsageCollection, len(kcu))
	for _, kcuc := range kcu {
		for _, kcucd := range kcuc.Data {
			kcucRev := kcuRev[kcucd.ReferencedTableName.String]
			kcucRev.Data = append(kcucRev.Data, &KeyColumnUsage{
				ConstraintCatalog:          kcucd.ConstraintCatalog,
				ConstraintSchema:           kcucd.ConstraintSchema,
				ConstraintName:             kcucd.ConstraintName,
				TableCatalog:               kcucd.TableCatalog,
				TableSchema:                kcucd.ReferencedTableSchema.String,
				TableName:                  kcucd.ReferencedTableName.String,
				ColumnName:                 kcucd.ReferencedColumnName.String,
				OrdinalPosition:            kcucd.OrdinalPosition,
				PositionInUniqueConstraint: kcucd.PositionInUniqueConstraint,
				ReferencedTableSchema:      null.MakeString(kcucd.TableSchema),
				ReferencedTableName:        null.MakeString(kcucd.TableName),
				ReferencedColumnName:       null.MakeString(kcucd.ColumnName),
			})
			kcuRev[kcucd.ReferencedTableName.String] = kcucRev
		}
	}
	return kcuRev
}

type relationKeyType int

func (r relationKeyType) String() string {
	switch r {
	case fKeyTypeNone:
		return "relKey:none"
	case fKeyTypePRI:
		return "relKey:PRI"
	case fKeyTypeMUL:
		return "relKey:MUL"
	default:
		panic("relationKeyType unknown type")
	}
}

const (
	fKeyTypeNone relationKeyType = iota
	fKeyTypePRI
	fKeyTypeMUL
)

type relTarget struct {
	column           string
	referencedTable  string
	referencedColumn string
	relationKeyType
}

type relTargets []relTarget

func (rt relTargets) hasManyToMany() bool {
	if len(rt) != 2 {
		return false
	}
	// 1. both main columns must be different named.
	// 2. key relation type must be equal and PRI
	// 3. referenced tables must be different
	return rt[0].column != rt[1].column && rt[0].relationKeyType == fKeyTypePRI && rt[0].relationKeyType == rt[1].relationKeyType &&
		rt[0].referencedTable != rt[1].referencedTable && rt[0].referencedColumn != rt[1].referencedColumn
}

// KeyRelationShips contains an internal cache about the database foreign key
// structure. It can only be created via function GenerateKeyRelationships.
type KeyRelationShips struct {
	// making the map private makes this type race free as reading the map from
	// multiple goroutines is allowed without a lock.
	// map[mainTable][]relTarget
	relMap map[string]relTargets
}

// IsOneToOne
func (krs KeyRelationShips) IsOneToOne(mainTable, mainColumn, referencedTable, referencedColumn string) bool {
	for _, rel := range krs.relMap[mainTable] {
		if rel.column == mainColumn && rel.referencedTable == referencedTable && rel.referencedColumn == referencedColumn && rel.relationKeyType == fKeyTypePRI {
			return true
		}
	}
	return false
}

// IsOneToMany returns true for a oneToMany or switching the tables for a ManyToOne relationship
func (krs KeyRelationShips) IsOneToMany(referencedTable, referencedColumn, mainTable, mainColumn string) bool {
	for _, rel := range krs.relMap[referencedTable] {
		if rel.column == referencedColumn && rel.referencedTable == mainTable && rel.referencedColumn == mainColumn && rel.relationKeyType == fKeyTypeMUL {
			return true
		}
	}
	return false
}

// ManyToManyTarget figures out if a table has M:N relationships and returns the
// target table and its column or empty strings if not found.
func (krs KeyRelationShips) ManyToManyTarget(referencedTable, referencedColumn, mainTable, mainColumn string) (table string, column string) {
	if krs.relMap[mainTable].hasManyToMany() {
		targetRefs := krs.relMap[mainTable]
		if targetRefs[0].column == mainColumn {
			return targetRefs[1].referencedTable, targetRefs[1].referencedColumn
		}
		if targetRefs[1].column == mainColumn {
			return targetRefs[0].referencedTable, targetRefs[0].referencedColumn
		}
	}
	return "", ""
}

// Debug writes the internal map in a sorted list to w.
func (krs KeyRelationShips) Debug(w io.Writer) {
	keys := make([]string, 0, len(krs.relMap))

	for k := range krs.relMap {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	for _, k := range keys {
		for _, rel := range krs.relMap[k] {
			fmt.Fprintf(w, "main: %s.%s => ref: %s.%s => %s\n", k, rel.column, rel.referencedTable, rel.referencedColumn, rel.relationKeyType)
		}
	}
}

// GenerateKeyRelationships loads the foreign key relationships between a list of
// given tables or all tables in a database. Might not yet work across several
// databases on the same file system.
func GenerateKeyRelationships(ctx context.Context, db dml.Querier, foreignKeys map[string]KeyColumnUsageCollection) (KeyRelationShips, error) {
	krs := KeyRelationShips{
		relMap: map[string]relTargets{},
	}

	fieldCount, err := countFieldsForTables(ctx, db)
	if err != nil {
		return KeyRelationShips{}, errors.WithStack(err)
	}

	for _, kcuc := range foreignKeys {
		for _, kcu := range kcuc.Data {

			// OneToOne relationship
			krs.relMap[kcu.TableName] = append(krs.relMap[kcu.TableName], relTarget{
				column:           kcu.ColumnName,
				referencedTable:  kcu.ReferencedTableName.String,
				referencedColumn: kcu.ReferencedColumnName.String,
				relationKeyType:  fKeyTypePRI,
			})

			// if referenced table has only one PK, then the reversed relationship of OneToMany is not possible
			if tc, ok := fieldCount[kcu.ReferencedTableName.String]; ok && tc.Pri == 1 && tc.Empty == 0 && tc.Mul == 0 {
				// OneToOne reversed is also possible
				krs.relMap[kcu.ReferencedTableName.String] = append(krs.relMap[kcu.ReferencedTableName.String], relTarget{
					column:           kcu.ReferencedColumnName.String,
					referencedTable:  kcu.TableName,
					referencedColumn: kcu.ColumnName,
					relationKeyType:  fKeyTypePRI,
				})
			}
			if tc, ok := fieldCount[kcu.TableName]; ok && (tc.Empty > 0 || tc.Pri > 1) {
				krs.relMap[kcu.ReferencedTableName.String] = append(krs.relMap[kcu.ReferencedTableName.String], relTarget{
					column:           kcu.ReferencedColumnName.String,
					referencedTable:  kcu.TableName,
					referencedColumn: kcu.ColumnName,
					relationKeyType:  fKeyTypeMUL,
				})
			}
		}
	}

	return krs, nil
}

type columnKeyCount struct {
	Empty, Mul, Pri, Uni int
}

func countFieldsForTables(ctx context.Context, db dml.Querier) (_ map[string]*columnKeyCount, err error) {
	const sqlQry = `SELECT TABLE_NAME, COLUMN_KEY, COUNT(*) AS FIELD_COUNT
 	FROM information_schema.COLUMNS WHERE TABLE_SCHEMA=DATABASE() GROUP BY TABLE_NAME,COLUMN_KEY`
	// TODO limit query to referencedTables, if provided

	rows, err := db.QueryContext(ctx, sqlQry)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	defer func() {
		// Not testable with the sqlmock package :-(
		if err2 := rows.Close(); err2 != nil && err == nil {
			err = errors.WithStack(err2)
		}
	}()

	col3 := &struct {
		TableName string
		ColumnKey string
		Count     int
	}{
		"", "", 0,
	}

	ret := map[string]*columnKeyCount{}
	for rows.Next() {
		if err = rows.Scan(&col3.TableName, &col3.ColumnKey, &col3.Count); err != nil {
			return nil, errors.WithStack(err)
		}
		ckc := ret[col3.TableName]
		if ckc == nil {
			ret[col3.TableName] = new(columnKeyCount)
			ckc = ret[col3.TableName]
		}

		switch col3.ColumnKey {
		case "":
			ckc.Empty = col3.Count
		case "MUL":
			ckc.Mul = col3.Count
		case "PRI":
			ckc.Pri = col3.Count
		case "UNI":
			ckc.Uni = col3.Count
		default:
			return nil, errors.NotSupported.Newf("[ddl] ColumnKey %q not supported", col3.ColumnKey)
		}

		col3.TableName = ""
		col3.ColumnKey = ""
		col3.Count = 0
	}
	if err = rows.Err(); err != nil {
		err = errors.WithStack(err)
	}
	return ret, err
}

func DisableForeignKeys(ctx context.Context, db dml.Execer, callBack func() error) (err error) {
	if _, err = db.ExecContext(ctx, "SET foreign_key_checks = 0;"); err != nil {
		return errors.WithStack(err)
	}
	defer func() {
		if _, err2 := db.ExecContext(ctx, "SET foreign_key_checks = 1;"); err2 != nil && err == nil {
			err = errors.WithStack(err2)
		}
	}()
	if err = callBack(); err != nil {
		return errors.WithStack(err)
	}
	return nil
}
