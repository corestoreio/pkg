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

package dmlgen

import (
	"bytes"
	"fmt"
	"strconv"
	"strings"
	"text/tabwriter"

	"github.com/corestoreio/pkg/sql/ddl"
	"github.com/corestoreio/pkg/util/bufferpool"
	"github.com/corestoreio/pkg/util/codegen"
	"github.com/corestoreio/pkg/util/strs"
)

// table writes one database table into Go source code.
type Table struct {
	Package              string // Name of the package
	Table                *ddl.Table
	Comment              string // Comment above the struct type declaration
	HasAutoIncrement     uint8  // 0=nil,1=false (has NO auto increment),2=true has auto increment
	HasEasyJSONMarshaler bool
	HasSerializer        bool // writes the .proto file if true

	// PrivateFields key=snake case name of the DB column, value=true, the field must be private
	debug                  bool // gets set via isDebug function
	privateFields          map[string]bool
	featuresInclude        FeatureToggle
	featuresExclude        FeatureToggle
	fieldMapFn             func(dbIdentifier string) (newName string)
	customStructTagFields  map[string]string
	relationshipSeen       map[string]bool // to not print twice a relationship
	availableRelationships []relationShipInfo
}

type relationShipInfo struct {
	isCollection          bool
	tableName             string
	structName            string
	mappedStructFieldName string // name of the struct in the relation struct
	columnName            string
}

func (t *Table) getFieldMapFn(g *Generator) func(dbIdentifier string) (newName string) {
	fieldMapFn := g.defaultTableConfig.FieldMapFn
	if fieldMapFn == nil {
		fieldMapFn = t.fieldMapFn
	}
	if fieldMapFn == nil {
		fieldMapFn = defaultFieldMapFn
	}
	return fieldMapFn
}

func (t *Table) IsFieldPublic(dbColumnName string) bool {
	return !t.privateFields[dbColumnName]
}

func (t *Table) IsFieldPrivate(dbColumnName string) bool {
	return t.privateFields[dbColumnName]
}

func (t *Table) GoCamelMaybePrivate(fieldName string) string {
	su := strs.ToGoCamelCase(fieldName)
	if t.IsFieldPublic(fieldName) {
		return su
	}
	return strs.LcFirst(su)
}

func (t *Table) CollectionName() string {
	return pluralize(t.Table.Name)
}

func (t *Table) EntityName() string {
	return strs.ToGoCamelCase(t.Table.Name)
}

func (t *Table) EntityNameLCFirst() string {
	return strs.LcFirst(strs.ToGoCamelCase(t.Table.Name))
}

func (t *Table) hasFeature(g *Generator, f FeatureToggle) bool {
	return g.hasFeature(t.featuresInclude, t.featuresExclude, f, 'a') // mode == AND
}

func (t *Table) fnCollectionStruct(mainGen *codegen.Go, g *Generator) {
	if !g.hasFeature(t.featuresInclude, t.featuresExclude, FeatureCollectionStruct) {
		return
	}

	mainGen.C(t.CollectionName(), `represents a collection type for DB table`, t.Table.Name)
	mainGen.C(`Not thread safe. Auto generated.`)

	mainGen.C(t.Comment != "", t.Comment)

	mainGen.Pln(t.HasEasyJSONMarshaler, `//easyjson:json`) // do not use C() because it adds a whitespace between "//" and "e"

	mainGen.Pln(`type `, t.CollectionName(), ` struct {`)
	{
		mainGen.In()
		mainGen.Pln(`Data []*`, t.EntityName(), codegen.EncloseBT(`json:"data,omitempty"`))

		if fn, ok := g.customCode["type_"+t.CollectionName()]; ok {
			fn(g, t, mainGen)
		}
		mainGen.Out()
	}
	mainGen.Pln(`}`)

	mainGen.C(`New`+t.CollectionName(), ` creates a new initialized collection. Auto generated.`)
	// TODO(idea): use a global pool which can register for each type the
	//  before/after mapcolumn function so that the dev does not need to
	//  assign each time. think if it's worth such a pattern.
	mainGen.Pln(`func New`+t.CollectionName(), `() *`, t.CollectionName(), ` {`)
	{
		mainGen.In()
		mainGen.Pln(`return &`, t.CollectionName(), `{`)
		{
			mainGen.In()
			mainGen.Pln(`Data: make([]*`, t.EntityName(), `, 0, 5),`)
			mainGen.Out()
		}
		mainGen.Pln(`}`)
		mainGen.Out()
	}
	mainGen.Pln(`}`)
}

func (t *Table) fnEntityStruct(mainGen *codegen.Go, g *Generator) {
	if !g.hasFeature(t.featuresInclude, t.featuresExclude, FeatureEntityStruct) {
		return
	}

	mainGen.C(t.EntityName(), `represents a single row for DB table`, t.Table.Name+`. Auto generated.`)
	if t.Comment != "" {
		mainGen.C(t.Comment)
	}
	if t.Table.TableComment != "" {
		mainGen.C("Table comment:", t.Table.TableComment)
	}
	if t.HasEasyJSONMarshaler {
		mainGen.Pln(`//easyjson:json`)
	}

	// Generate table structs
	mainGen.Pln(`type `, t.EntityName(), ` struct {`)
	{
		if fn, ok := g.customCode["type_"+t.EntityName()]; ok {
			fn(g, t, mainGen)
		} else {
			mainGen.In()
			for _, c := range t.Table.Columns {
				structTag := ""
				if c.StructTag != "" {
					structTag += "`" + c.StructTag + "`"
				}
				mainGen.Pln(t.GoCamelMaybePrivate(c.Field), g.goTypeNull(c), structTag, c.GoComment())
			}
			if len(t.availableRelationships) > 0 {
				mainGen.P(`Relations *`, t.relationStructName()) // TODO use customStructTagFields
				if fn, ok := g.customCode["type_relation_"+t.relationStructName()]; ok {
					fn(g, t, mainGen)
				} else {
					mainGen.Pln("")
				}
			}
			mainGen.Out()
		}
	}
	mainGen.Pln(`}`)
}

func (t *Table) relationStructName() string {
	return t.EntityNameLCFirst() + `Relations`
}

// this part is duplicated in the proto file generation function generateProto.
func (t *Table) fnEntityRelationStruct(mainGen *codegen.Go, g *Generator) {
	_, ok1 := g.kcu[t.Table.Name]
	_, ok2 := g.kcuRev[t.Table.Name]
	if !g.hasFeature(t.featuresInclude, t.featuresExclude, FeatureEntityRelationships) || (!ok1 && !ok2) {
		return
	}

	fieldMapFn := t.getFieldMapFn(g)
	// only write the relation struct if there are some relations, if non, revert to the old buffer.
	mainGenBuf := mainGen.Buffer
	mainGen.Buffer = new(bytes.Buffer)
	defer func() {
		data := mainGen.Buffer.Bytes()
		mainGen.Buffer = mainGenBuf // restore old buffer
		if len(t.availableRelationships) > 0 {
			mainGen.Buffer.Write(data)
		}
	}()

	mainGen.Pln(`type `, t.relationStructName(), ` struct {`)
	mainGen.In()
	if fn, ok := g.customCode["type_"+t.relationStructName()]; ok {
		fn(g, t, mainGen)
	}
	mainGen.Pln(`parent *`, t.EntityName())

	debugBuf := bufferpool.Get()
	defer bufferpool.Put(debugBuf)
	tabW := tabwriter.NewWriter(debugBuf, 6, 0, 2, ' ', 0)
	var hasAtLeastOneRelationShip int
	fmt.Fprintf(debugBuf, "RelationInfo for: %q\n", t.Table.Name)
	fmt.Fprintf(tabW, "Case\tis1:M\tis1:1\tseen?\tisRelAl\thasTable\tTarget Tbl M:N\tRelation\n")
	if kcuc, ok := g.kcu[t.Table.Name]; ok { // kcu = keyColumnUsage && kcuc = keyColumnUsageCollection
		for _, kcuce := range kcuc.Data {
			if !kcuce.ReferencedTableName.Valid {
				continue
			}
			hasAtLeastOneRelationShip++
			// case ONE-TO-MANY
			isOneToMany := g.krs.IsOneToMany(kcuce.TableName, kcuce.ColumnName, kcuce.ReferencedTableName.Data, kcuce.ReferencedColumnName.Data)
			isRelationAllowed := g.isAllowedRelationship(kcuce.TableName, kcuce.ColumnName, kcuce.ReferencedTableName.Data, kcuce.ReferencedColumnName.Data)
			hasTable := g.Tables[kcuce.ReferencedTableName.Data] != nil
			fmt.Fprintf(tabW, "A1_1:M\t%t\t%t\t%t\t%t\t%t\t-\t%s => %s\n", isOneToMany, false, false, isRelationAllowed, hasTable,
				kcuce.TableName+"."+kcuce.ColumnName, kcuce.ReferencedTableName.Data+"."+kcuce.ReferencedColumnName.Data)
			if isOneToMany && hasTable && isRelationAllowed {

				name := pluralize(kcuce.ReferencedTableName.Data)
				fieldName := fieldMapFn(name)
				t.availableRelationships = append(t.availableRelationships, relationShipInfo{
					isCollection:          true,
					tableName:             kcuce.ReferencedTableName.Data,
					structName:            name,
					mappedStructFieldName: fieldName,
					columnName:            kcuce.ReferencedColumnName.Data,
				})

				mainGen.Pln(fieldName, " *", name,
					t.customStructTagFields[kcuce.ReferencedTableName.Data],
					"// 1:M", kcuce.TableName+"."+kcuce.ColumnName, "=>", kcuce.ReferencedTableName.Data+"."+kcuce.ReferencedColumnName.Data)
			}

			// case ONE-TO-ONE
			isOneToOne := g.krs.IsOneToOne(kcuce.TableName, kcuce.ColumnName, kcuce.ReferencedTableName.Data, kcuce.ReferencedColumnName.Data)
			fmt.Fprintf(tabW, "B1_1:1\t%t\t%t\t%t\t%t\t%t\t-\t%s => %s\n", isOneToMany, isOneToOne, false, isRelationAllowed, hasTable,
				kcuce.TableName+"."+kcuce.ColumnName, kcuce.ReferencedTableName.Data+"."+kcuce.ReferencedColumnName.Data)
			if isOneToOne && hasTable && isRelationAllowed {

				name := strs.ToGoCamelCase(kcuce.ReferencedTableName.Data)
				fieldName := fieldMapFn(name)
				t.availableRelationships = append(t.availableRelationships, relationShipInfo{
					isCollection:          true,
					tableName:             kcuce.ReferencedTableName.Data,
					structName:            name,
					mappedStructFieldName: fieldName,
					columnName:            kcuce.ReferencedColumnName.Data,
				})

				mainGen.Pln(fieldName, " *", name, t.customStructTagFields[kcuce.ReferencedTableName.Data],
					"// 1:1", kcuce.TableName+"."+kcuce.ColumnName, "=>", kcuce.ReferencedTableName.Data+"."+kcuce.ReferencedColumnName.Data)
			}

			// case MANY-TO-MANY
			targetTbl, targetColumn := g.krs.ManyToManyTarget(kcuce.TableName, kcuce.ColumnName)
			fmt.Fprintf(tabW, "C1_M:N\t%t\t%t\t%t\t%t\t%t\t%s\t%s => %s\n", isOneToMany, isOneToOne, false, isRelationAllowed, hasTable,
				targetTbl+"."+targetColumn,
				kcuce.TableName+"."+kcuce.ColumnName, kcuce.ReferencedTableName.Data+"."+kcuce.ReferencedColumnName.Data)
			// hasTable variable shall not be added because usually the link table does not get loaded.
			if isRelationAllowed && targetTbl != "" && targetColumn != "" {

				name := pluralize(targetTbl)
				fieldName := fieldMapFn(name)
				t.availableRelationships = append(t.availableRelationships, relationShipInfo{
					isCollection:          true,
					tableName:             kcuce.ReferencedTableName.Data,
					structName:            name,
					mappedStructFieldName: fieldName,
					columnName:            kcuce.ReferencedColumnName.Data,
				})

				mainGen.Pln(fieldName, " *", name, t.customStructTagFields[targetTbl],
					"// M:N", kcuce.TableName+"."+kcuce.ColumnName, "via", kcuce.ReferencedTableName.Data+"."+kcuce.ReferencedColumnName.Data,
					"=>", targetTbl+"."+targetColumn,
				)
			}
		}
	}

	if kcuc, ok := g.kcuRev[t.Table.Name]; ok { // kcu = keyColumnUsage && kcuc = keyColumnUsageCollection
		for _, kcuce := range kcuc.Data {
			if !kcuce.ReferencedTableName.Valid {
				continue
			}
			hasAtLeastOneRelationShip++
			// case ONE-TO-MANY
			isOneToMany := g.krs.IsOneToMany(kcuce.TableName, kcuce.ColumnName, kcuce.ReferencedTableName.Data, kcuce.ReferencedColumnName.Data)
			isRelationAllowed := g.isAllowedRelationship(kcuce.TableName, kcuce.ColumnName, kcuce.ReferencedTableName.Data, kcuce.ReferencedColumnName.Data)
			hasTable := g.Tables[kcuce.ReferencedTableName.Data] != nil
			keySeen := fieldMapFn(pluralize(kcuce.ReferencedTableName.Data))
			relationShipSeenAlready := t.relationshipSeen[keySeen]
			// case ONE-TO-MANY
			fmt.Fprintf(tabW, "A2_1:M rev\t%t\t%t\t%t\t%t\t%t\t-\t%s => %s\n", isOneToMany, false, relationShipSeenAlready, isRelationAllowed, hasTable,
				kcuce.TableName+"."+kcuce.ColumnName, kcuce.ReferencedTableName.Data+"."+kcuce.ReferencedColumnName.Data)
			if isRelationAllowed && isOneToMany && hasTable && !relationShipSeenAlready {

				name := pluralize(kcuce.ReferencedTableName.Data)
				fieldName := fieldMapFn(name)
				t.availableRelationships = append(t.availableRelationships, relationShipInfo{
					isCollection:          true,
					tableName:             kcuce.ReferencedTableName.Data,
					structName:            name,
					mappedStructFieldName: fieldName,
					columnName:            kcuce.ReferencedColumnName.Data,
				})

				mainGen.Pln(fieldName, " *", name, t.customStructTagFields[kcuce.ReferencedTableName.Data],
					"// Reversed 1:M", kcuce.TableName+"."+kcuce.ColumnName, "=>", kcuce.ReferencedTableName.Data+"."+kcuce.ReferencedColumnName.Data)
				t.relationshipSeen[keySeen] = true
			}

			// case ONE-TO-ONE
			isOneToOne := g.krs.IsOneToOne(kcuce.TableName, kcuce.ColumnName, kcuce.ReferencedTableName.Data, kcuce.ReferencedColumnName.Data)
			fmt.Fprintf(tabW, "B2_1:1 rev\t%t\t%t\t%t\t%t\t%t\t-\t%s => %s\n", isOneToMany, isOneToOne, relationShipSeenAlready, isRelationAllowed, hasTable,
				kcuce.TableName+"."+kcuce.ColumnName, kcuce.ReferencedTableName.Data+"."+kcuce.ReferencedColumnName.Data)
			if isRelationAllowed && isOneToOne && hasTable {

				name := strs.ToGoCamelCase(kcuce.ReferencedTableName.Data)
				fieldName := fieldMapFn(name)
				t.availableRelationships = append(t.availableRelationships, relationShipInfo{
					tableName:             kcuce.ReferencedTableName.Data,
					structName:            name,
					mappedStructFieldName: fieldName,
					columnName:            kcuce.ReferencedColumnName.Data,
				})

				mainGen.Pln(fieldName, " *", name, t.customStructTagFields[kcuce.ReferencedTableName.Data],
					"// Reversed 1:1", kcuce.TableName+"."+kcuce.ColumnName, "=>", kcuce.ReferencedTableName.Data+"."+kcuce.ReferencedColumnName.Data)
			}

			// case MANY-TO-MANY
			targetTbl, targetColumn := g.krs.ManyToManyTarget(kcuce.ReferencedTableName.Data, kcuce.ReferencedColumnName.Data)
			if targetTbl != "" && targetColumn != "" {
				keySeen := fieldMapFn(pluralize(targetTbl))
				isRelationAllowed = g.isAllowedRelationship(kcuce.TableName, kcuce.ColumnName, targetTbl, targetColumn) && !t.relationshipSeen[keySeen]
				t.relationshipSeen[keySeen] = true
			}

			// case MANY-TO-MANY
			fmt.Fprintf(tabW, "C2_M:N rev\t%t\t%t\t%t\t%t\t%t\t%s\t%s => %s\n", isOneToMany, isOneToOne, relationShipSeenAlready, isRelationAllowed, hasTable,
				targetTbl+"."+targetColumn,
				kcuce.TableName+"."+kcuce.ColumnName, kcuce.ReferencedTableName.Data+"."+kcuce.ReferencedColumnName.Data)
			// hasTable shall not be added because usually the link table does not get loaded.
			if isRelationAllowed && targetTbl != "" && targetColumn != "" {

				name := pluralize(targetTbl)
				fieldName := fieldMapFn(name)
				t.availableRelationships = append(t.availableRelationships, relationShipInfo{
					isCollection:          true,
					tableName:             kcuce.ReferencedTableName.Data,
					structName:            name,
					mappedStructFieldName: fieldName,
					columnName:            kcuce.ReferencedColumnName.Data,
				})

				mainGen.Pln(fieldName, " *", name, t.customStructTagFields[targetTbl],
					"// Reversed M:N", kcuce.TableName+"."+kcuce.ColumnName, "via", kcuce.ReferencedTableName.Data+"."+kcuce.ReferencedColumnName.Data,
					"=>", targetTbl+"."+targetColumn,
				)
			}
		}
	}
	if t.debug && hasAtLeastOneRelationShip > 0 {
		_ = tabW.Flush()
		fmt.Fprintf(debugBuf, "Relationship count: %d\n", hasAtLeastOneRelationShip)
		fmt.Println(debugBuf.String())
	}

	mainGen.Out()
	mainGen.Pln(`}`) // end type struct

	mainGen.Pln(`func (e *`, t.EntityName(), `) setRelationParent() {`)
	mainGen.Pln(`if e.Relations != nil && e.Relations.parent == nil {
			e.Relations.parent = e
		}
	}`)

	mainGen.Pln(`func (e *`, t.EntityName(), `) NewRelations() *`, t.relationStructName(), ` {`)
	mainGen.Pln(`e.Relations = &`, t.relationStructName(), ` { parent: e }
			return e.Relations }`)
}

func (t *Table) fnEntityRelationMethods(mainGen *codegen.Go, g *Generator) {
	if !g.hasFeature(t.featuresInclude, t.featuresExclude, FeatureEntityRelationships) || len(t.availableRelationships) == 0 {
		return
	}

	parentPK := t.Table.Columns.PrimaryKeys().First() // TODO support multiple Primary Keys
	parentPKFieldName := strs.ToGoCamelCase(parentPK.Field)

	for _, rs := range t.availableRelationships {
		// <DELETE>
		mainGen.Pln(`func (r *`, t.relationStructName(), `) `, "Delete"+rs.mappedStructFieldName, `(ctx context.Context, dbm *DBM, opts ...dml.DBRFunc) error {`)
		mainGen.Pln(`dbr := dbm.ConnPool.WithCacheKey(`, codegen.SkipWS(`"`, rs.mappedStructFieldName, `DeleteByFK`, `"`), `, opts...)`)
		mainGen.Pln(`res, err := dbr.ExecContext(ctx, r.parent.`, parentPKFieldName, `)`)
		mainGen.Pln(`err = dbr.ResultCheckFn(`, constTableName(rs.tableName), `, len(r.`, rs.mappedStructFieldName, `.Data), res, err)`)
		mainGen.Pln(`if err == nil && r.`, rs.mappedStructFieldName, ` != nil { r.`, rs.mappedStructFieldName, `.Clear() }`)
		mainGen.Pln(`return errors.WithStack(err)
		}`)
		// </DELETE>

		// <INSERT>
		mainGen.Pln(`func (r *`, t.relationStructName(), `) `, codegen.SkipWS("Insert", rs.mappedStructFieldName), `(ctx context.Context, dbm *DBM, opts ...dml.DBRFunc) error {`)
		mainGen.Pln(`if r.`, rs.mappedStructFieldName, ` == nil || len(r.`, rs.mappedStructFieldName, `.Data) == 0 { return nil }
		for _, e2 := range r.`, rs.mappedStructFieldName, `.Data {
			e2.`, strs.ToGoCamelCase(rs.columnName), ` = `, g.convertType(
			parentPK,
			g.findColumn(rs.tableName, rs.columnName),
			`r.parent.`+parentPKFieldName,
		), `
		}
			return errors.WithStack(r.`, rs.mappedStructFieldName, `.DBInsert(ctx, dbm, opts...)) }`)
		// </INSERT>

		// <UPDATE>
		mainGen.Pln(`func (r *customerEntityRelations) `, codegen.SkipWS(`Update`, rs.mappedStructFieldName), `(ctx context.Context, dbm *DBM, opts ...dml.DBRFunc) (err error) {`)
		mainGen.Pln(`if r.`, rs.mappedStructFieldName, ` == nil || len(r.`, rs.mappedStructFieldName, `.Data) == 0 {
			dbr := dbm.ConnPool.WithCacheKey(`, codegen.SkipWS(`"`, rs.mappedStructFieldName, `DeleteByFK`, `"`), `, opts...)
			res, err := dbr.ExecContext(ctx, r.parent.`, parentPKFieldName, `)
			return dbr.ResultCheckFn(`, constTableName(rs.tableName), `, -1, res, errors.WithStack(err))
		}
		for _, e2 := range r.`, rs.mappedStructFieldName, `.Data {
				e2.`, strs.ToGoCamelCase(rs.columnName), ` = `, g.convertType(
			parentPK,
			g.findColumn(rs.tableName, rs.columnName),
			`r.parent.`+parentPKFieldName,
		), `
		}
		err = r.`, rs.mappedStructFieldName, `.DBUpdate(ctx, dbm, opts...)
		return errors.WithStack(err)
}`)
		// </UPDATE>

		// <SELECT>
		mainGen.Pln(`func (r *`, t.relationStructName(), `) `, "Load"+rs.mappedStructFieldName, `(ctx context.Context, dbm *DBM, opts ...dml.DBRFunc) (rowCount uint64, err error) {`)
		mainGen.Pln(`if r.`, rs.mappedStructFieldName, ` == nil { r.`, rs.mappedStructFieldName, ` = &`, rs.mappedStructFieldName, `{} }`)
		mainGen.Pln(`r.`, rs.mappedStructFieldName, `.Clear()
			  rowCount, err = dbm.ConnPool.WithCacheKey(`, codegen.SkipWS(`"`, rs.mappedStructFieldName, `SelectByFK"`), `, opts...).Load(ctx, r.`, rs.mappedStructFieldName, `, r.parent.EntityID)
				return rowCount, errors.WithStack(err) }`)
		// </SELECT>

	} // end for availableRelationships

	// TODO for all `All` functions add a possibility to load async.
	//g, ctx := errgroup.WithContext(ctx)
	//g.Go(func() error {
	//	_, err = r.LoadCustomerAddressEntities(ctx, dbm, opts...)
	//	return errors.WithStack(err)
	//})
	//return g.Wait()

	// <INSERT_ALL>
	mainGen.Pln(`func (r *`, t.relationStructName(), `) InsertAll(ctx context.Context, dbm *DBM, opts ...dml.DBRFunc) error {`)
	for _, rs := range t.availableRelationships {
		mainGen.Pln(`if err := r.`, codegen.SkipWS("Insert", rs.mappedStructFieldName), `(ctx, dbm, opts...); err != nil { return errors.WithStack(err) }`)
	}
	mainGen.Pln(`return nil }`)
	// </INSERT_ALL>

	// <SELECT_ALL>
	mainGen.Pln(`func (r *`, t.relationStructName(), `) LoadAll(ctx context.Context, dbm *DBM, opts ...dml.DBRFunc) (err error) {`)
	for _, rs := range t.availableRelationships {
		mainGen.Pln(`if _, err = r.`, codegen.SkipWS("Load", rs.mappedStructFieldName), `(ctx, dbm, opts...); err != nil { return errors.WithStack(err) }`)
	}
	mainGen.Pln(`return nil }`)
	// </SELECT_ALL>

	// <UPDATE_ALL>
	mainGen.Pln(`func (r *`, t.relationStructName(), `) UpdateAll(ctx context.Context, dbm *DBM, opts ...dml.DBRFunc) error {`)
	for _, rs := range t.availableRelationships {
		mainGen.Pln(`if err := r.`, codegen.SkipWS("Update", rs.mappedStructFieldName), `(ctx, dbm, opts...); err != nil { return errors.WithStack(err) }`)
	}
	mainGen.Pln(`return nil }`)
	// </UPDATE_ALL>

	// <DELETE_ALL>
	mainGen.Pln(`func (r *`, t.relationStructName(), `) DeleteAll(ctx context.Context, dbm *DBM, opts ...dml.DBRFunc) error {`)
	for _, rs := range t.availableRelationships {
		mainGen.Pln(`if err := r.`, codegen.SkipWS("Delete", rs.mappedStructFieldName), `(ctx, dbm, opts...); err != nil { return errors.WithStack(err) }`)
	}
	mainGen.Pln(`return nil }`)
	// </DELETE_ALL>
}

func (t *Table) fnEntityGetSetPrivateFields(mainGen *codegen.Go, g *Generator) {
	if !g.hasFeature(t.featuresInclude, t.featuresExclude, FeatureEntityGetSetPrivateFields) {
		return
	}
	// Generates the Getter/Setter for private fields
	for _, c := range t.Table.Columns {
		if !t.IsFieldPrivate(c.Field) {
			continue
		}
		mainGen.C(`Set`, strs.ToGoCamelCase(c.Field), ` sets the data for a private and security sensitive field.`)
		mainGen.Pln(`func (e *`, t.EntityName(), `) Set`+strs.ToGoCamelCase(c.Field), `(d `, g.goTypeNull(c), `) *`, t.EntityName(), ` {`)
		{
			mainGen.In()
			mainGen.Pln(`e.`, t.GoCamelMaybePrivate(c.Field), ` = d`)
			mainGen.Pln(`return e`)
			mainGen.Out()
		}
		mainGen.Pln(`}`)

		mainGen.C(`Get`, strs.ToGoCamelCase(c.Field), ` returns the data from a private and security sensitive field.`)
		mainGen.Pln(`func (e *`, t.EntityName(), `) Get`+strs.ToGoCamelCase(c.Field), `() `, g.goTypeNull(c), `{`)
		{
			mainGen.In()
			mainGen.Pln(`return e.`, t.GoCamelMaybePrivate(c.Field))
			mainGen.Out()
		}
		mainGen.Pln(`}`)

	}
}

func (t *Table) fnEntityEmpty(mainGen *codegen.Go, g *Generator) {
	if !g.hasFeature(t.featuresInclude, t.featuresExclude, FeatureEntityEmpty) {
		return
	}
	mainGen.C(`Empty empties all the fields of the current object. Also known as Reset.`)
	// no idea if pointer dereferencing is bad ...
	mainGen.Pln(`func (e *`, t.EntityName(), `) Empty() *`, t.EntityName(), ` { *e = `, t.EntityName(), `{}; return e }`)
}

func (t *Table) fnEntityIsSet(mainGen *codegen.Go, g *Generator) {
	if !g.hasFeature(t.featuresInclude, t.featuresExclude, FeatureEntityIsSet|FeatureDB|FeatureDBSelect) {
		return
	}
	tblPKs := t.Table.Columns.PrimaryKeys()
	if t.Table.IsView() {
		tblPKs = t.Table.Columns.ViewPrimaryKeys()
	}

	// TODO maybe unique keys should also be added.
	var buf strings.Builder
	i := 0
	tblPKs.Each(func(c *ddl.Column) {
		if i > 0 {
			buf.WriteString(" && ")
		}
		buf.WriteString("e.")
		buf.WriteString(strs.ToGoCamelCase(c.Field))
		buf.WriteString(mySQLType2GoComparisonOperator(c))
		i++
	})
	if i == 0 {
		return // no PK fields found
	}
	mainGen.C(`IsSet returns true if the entity has non-empty primary keys.`)
	mainGen.Pln(`func (e *`, t.EntityName(), `) IsSet() bool { return `, buf.String(), ` }`)
}

func (t *Table) fnEntityCopy(mainGen *codegen.Go, g *Generator) {
	if !g.hasFeature(t.featuresInclude, t.featuresExclude, FeatureEntityCopy) {
		return
	}
	mainGen.C(`Copy copies the struct and returns a new pointer. TODO use deepcopy tool to generate code afterwards`)
	mainGen.Pln(`func (e *`, t.EntityName(), `) Copy() *`, t.EntityName(), ` {
		if e == nil { return &`, t.EntityName(), `{} }
		e2 := *e // for now a shallow copy
		return &e2
}`)
}

func (t *Table) fnEntityWriteTo(mainGen *codegen.Go, g *Generator) {
	if !g.hasFeature(t.featuresInclude, t.featuresExclude, FeatureEntityWriteTo) {
		return
	}
	mainGen.C(`WriteTo implements io.WriterTo and writes the field names and their values to w.`,
		`This is especially useful for debugging or or generating a hash of the struct.`)

	mainGen.Pln(`func (e *`, t.EntityName(), `) WriteTo(w io.Writer) (n int64, err error) {
	// for now this printing is good enough. If you need better swap out with your code.`)

	if fn, ok := g.customCode["func_"+t.EntityName()+"_WriteTo"]; ok {
		fn(g, t, mainGen)
	} else {
		mainGen.Pln(`n2, err := fmt.Fprint(w,`)
		mainGen.In()
		t.Table.Columns.Each(func(c *ddl.Column) {
			if t.IsFieldPublic(c.Field) {
				mainGen.Pln(`"`+c.Field+`:"`, `, e.`, strs.ToGoCamelCase(c.Field), `,`, `"\n",`)
			}
		})
		mainGen.Pln(`)`)
		mainGen.Pln(`return int64(n2), err`)
		mainGen.Out()
	}
	mainGen.Pln(`}`)
}

func (t *Table) fnCollectionWriteTo(mainGen *codegen.Go, g *Generator) {
	if !g.hasFeature(t.featuresInclude, t.featuresExclude, FeatureEntityWriteTo) {
		return
	}

	mainGen.C(`WriteTo implements io.WriterTo and writes the field names and their values to w.`,
		`This is especially useful for debugging or or generating a hash of the struct.`)

	mainGen.Pln(`func (cc *`, t.CollectionName(), `) WriteTo(w io.Writer) (n int64, err error) {
		for i,d := range cc.Data {
			n2,err := d.WriteTo(w)
			if err != nil {
				return 0, errors.Wrapf(err,"[`+t.Package+`] WriteTo failed at index %d",i)
			}
			n+=n2
		}
		return n,nil
	}`)
}

func (t *Table) fnEntityDBMapColumns(mainGen *codegen.Go, g *Generator) {
	if !g.hasFeature(t.featuresInclude, t.featuresExclude, FeatureDBMapColumns|
		FeatureDB|FeatureDBSelect|FeatureDBDelete|
		FeatureDBInsert|FeatureDBUpdate|FeatureDBUpsert) {
		return
	}
	mainGen.C(`MapColumns implements interface ColumnMapper only partially. Auto generated.`)
	mainGen.Pln(`func (e *`, t.EntityName(), `) MapColumns(cm *dml.ColumnMap) error {`)
	{
		if fn, ok := g.customCode["func_"+t.EntityName()+"_MapColumns"]; ok {
			fn(g, t, mainGen)
		}

		mainGen.Pln(`for cm.Next(`, t.Table.Columns.Len(), `) {`)
		{
			mainGen.In()
			mainGen.Pln(`switch c := cm.Column(); c {`)
			{
				mainGen.In()
				t.Table.Columns.Each(func(c *ddl.Column) {
					mainGen.P(`case`, strconv.Quote(c.Field))
					for _, a := range c.Aliases {
						mainGen.P(`,`, strconv.Quote(a))
					}
					mainGen.Pln(codegen.SkipWS(`,"`, c.Pos-1, `"`), `:`)
					mainGen.Pln(`cm.`, g.goFuncNull(c), `(&e.`, t.GoCamelMaybePrivate(c.Field), `)`)
				})
				mainGen.Pln(`default:`)
				mainGen.Pln(`return errors.NotFound.Newf("[`+g.Package+`]`, t.EntityName(), `Column %q not found", c)`)
				mainGen.Out()
			}
			mainGen.Pln(`}`)
			mainGen.Out()
		}
		mainGen.Pln(`}`)
		mainGen.Pln(`return errors.WithStack(cm.Err())`)
		mainGen.Out()
	}
	mainGen.Pln(`}`)
}

func (t *Table) hasPKAutoInc() bool {
	var hasPKAutoInc bool
	t.Table.Columns.Each(func(c *ddl.Column) {
		if c.IsPK() && c.IsAutoIncrement() {
			hasPKAutoInc = true
		}
		if hasPKAutoInc {
			return
		}
	})
	return hasPKAutoInc
}

func (t *Table) fnEntityDBAssignLastInsertID(mainGen *codegen.Go, g *Generator) {
	if !g.hasFeature(t.featuresInclude, t.featuresExclude, FeatureDBAssignLastInsertID|
		FeatureDB|FeatureDBInsert|FeatureDBUpsert) {
		return
	}
	if !t.hasPKAutoInc() {
		return
	}

	mainGen.C(`AssignLastInsertID updates the increment ID field with the last inserted ID from an INSERT operation.`,
		`Implements dml.InsertIDAssigner. Auto generated.`)
	mainGen.Pln(`func (e *`, t.EntityName(), `) AssignLastInsertID(id int64) {`)
	{
		mainGen.In()
		t.Table.Columns.Each(func(c *ddl.Column) {
			if c.IsPK() && c.IsAutoIncrement() {
				mainGen.Pln(`e.`, t.GoCamelMaybePrivate(c.Field), ` = `, g.goType(c), `(id)`)
			}
		})
		mainGen.Out()
	}
	mainGen.Pln(`}`)
}

func (t *Table) fnCollectionDBAssignLastInsertID(mainGen *codegen.Go, g *Generator) {
	if !g.hasFeature(t.featuresInclude, t.featuresExclude, FeatureDBAssignLastInsertID|
		FeatureDB|FeatureDBInsert|FeatureDBUpsert) {
		return
	}
	if !t.hasPKAutoInc() {
		return
	}

	mainGen.C(`AssignLastInsertID traverses through the slice and sets an incrementing new ID to each entity.`)
	mainGen.Pln(`func (cc *`, t.CollectionName(), `) AssignLastInsertID(id int64) {`)
	{
		mainGen.In()
		mainGen.Pln(`for i:=0 ; i < len(cc.Data); i++ {`)
		{
			mainGen.In()
			mainGen.Pln(`cc.Data[i].AssignLastInsertID(id + int64(i))`)
			mainGen.Out()
		}
		mainGen.Pln(`}`)
		mainGen.Out()
	}
	mainGen.Pln(`}`)
}

func (t *Table) fnCollectionUniqueGetters(mainGen *codegen.Go, g *Generator) {
	if !g.hasFeature(t.featuresInclude, t.featuresExclude, FeatureCollectionUniqueGetters|
		FeatureDB|FeatureDBSelect|FeatureDBDelete|
		FeatureDBInsert|FeatureDBUpdate|FeatureDBUpsert) {
		return
	}

	// Generates functions to return all data as a slice from unique/primary
	// columns.
	for _, c := range t.Table.Columns.UniqueColumns() {
		gtn := g.goTypeNull(c)
		goCamel := strs.ToGoCamelCase(c.Field)
		mainGen.C(goCamel + `s returns a slice with the data or appends it to a slice.`)
		mainGen.C(`Auto generated.`)
		mainGen.Pln(`func (cc *`, t.CollectionName(), `) `, goCamel+`s(ret ...`+gtn, `) []`+gtn, ` {`)
		{
			mainGen.In()
			mainGen.Pln(`if cc == nil {	return nil }`)
			mainGen.Pln(`if ret == nil {`)
			{
				mainGen.In()
				mainGen.Pln(`ret = make([]`+gtn, `, 0, len(cc.Data))`)
				mainGen.Out()
			}
			mainGen.Pln(`}`)
			mainGen.Pln(`for _, e := range cc.Data {`)
			{
				mainGen.In()
				mainGen.Pln(`ret = append(ret, e.`+goCamel, `)`)
				mainGen.Out()
			}
			mainGen.Pln(`}`)
			mainGen.Pln(`return ret`)
			mainGen.Out()
		}
		mainGen.Pln(`}`)
	}
}

func (t *Table) fnCollectionUniquifiedGetters(mainGen *codegen.Go, g *Generator) {
	if !g.hasFeature(t.featuresInclude, t.featuresExclude, FeatureCollectionUniquifiedGetters) {
		return
	}
	// Generates functions to return data with removed duplicates from any
	// column which has set the flag Uniquified.
	for _, c := range t.Table.Columns.UniquifiedColumns() {
		goType := g.goType(c)
		goCamel := strs.ToGoCamelCase(c.Field)

		mainGen.C(goCamel+`s belongs to the column`, strconv.Quote(c.Field), `and returns a slice or appends to a slice only`,
			`unique values of that column. The values will be filtered internally in a Go map. No DB query gets`,
			`executed. Auto generated.`)
		mainGen.Pln(`func (cc *`, t.CollectionName(), `) Unique`+goCamel+`s(ret ...`, goType, `) []`, goType, ` {`)
		{
			mainGen.In()
			mainGen.Pln(`if cc == nil {	return nil }`)
			mainGen.Pln(`if ret == nil {
					ret = make([]`, goType, `, 0, len(cc.Data))
				}`)

			// TODO: a reusable map and use different algorithms depending on
			// the size of the cc.Data slice. Sometimes a for/for loop runs
			// faster than a map.
			goPrimNull := g.toGoPrimitiveFromNull(c)
			mainGen.Pln(`dupCheck := make(map[`, goType, `]bool, len(cc.Data))`)
			mainGen.Pln(`for _, e := range cc.Data {`)
			{
				mainGen.In()
				mainGen.Pln(`if !dupCheck[e.`+goPrimNull, `] {`)
				{
					mainGen.In()
					mainGen.Pln(`ret = append(ret, e.`, goPrimNull, `)`)
					mainGen.Pln(`dupCheck[e.`+goPrimNull, `] = true`)
					mainGen.Out()
				}
				mainGen.Pln(`}`)
				mainGen.Out()
			}
			mainGen.Pln(`}`)
			mainGen.Pln(`return ret`)
			mainGen.Out()
		}
		mainGen.Pln(`}`)
	}
}

func (t *Table) fnCollectionFilter(mainGen *codegen.Go, g *Generator) {
	if !g.hasFeature(t.featuresInclude, t.featuresExclude, FeatureCollectionFilter) {
		return
	}
	mainGen.C(`Filter filters the current slice by predicate f without memory allocation. Auto generated via dmlgen.`)
	mainGen.Pln(`func (cc *`, t.CollectionName(), `) Filter(f func(*`, t.EntityName(), `) bool) *`, t.CollectionName(), ` {`)
	{
		mainGen.In()
		mainGen.Pln(`if cc == nil {	return nil }`)
		mainGen.Pln(`b,i := cc.Data[:0],0`)
		mainGen.Pln(`for _, e := range cc.Data {`)
		{
			mainGen.In()
			mainGen.Pln(`if f(e) {`)
			{
				mainGen.Pln(`b = append(b, e)`)
			}
			mainGen.Pln(`}`) // endif
			mainGen.Pln(`i++`)
		}
		mainGen.Out()
		mainGen.Pln(`}`) // for loop
		mainGen.Pln(`for i := len(b); i < len(cc.Data); i++ {
				cc.Data[i] = nil // this should avoid the memory leak
			}`)

		mainGen.Pln(`cc.Data = b`)
		mainGen.Pln(`return cc`)
		mainGen.Out()
	}
	mainGen.Pln(`}`) // function
}

func (t *Table) fnCollectionEach(mainGen *codegen.Go, g *Generator) {
	if !g.hasFeature(t.featuresInclude, t.featuresExclude, FeatureCollectionEach) {
		return
	}
	mainGen.C(`Each will run function f on all items in []*`, t.EntityName(), `. Auto generated via dmlgen.`)
	mainGen.Pln(`func (cc *`, t.CollectionName(), `) Each(f func(*`, t.EntityName(), `)) *`, t.CollectionName(), ` {`)
	{
		mainGen.Pln(`if cc == nil {	return nil }`)
		mainGen.Pln(`for i := range cc.Data {`)
		{
			mainGen.Pln(`f(cc.Data[i])`)
		}
		mainGen.Pln(`}`)
		mainGen.Pln(`return cc`)
	}
	mainGen.Pln(`}`)
}

// Clear because Reset name is used by gogo protobuf
func (t *Table) fnCollectionClear(mainGen *codegen.Go, g *Generator) {
	if !g.hasFeature(t.featuresInclude, t.featuresExclude, FeatureCollectionClear) {
		return
	}
	mainGen.C(`Clear will reset the data slice or create a new type. Useful for reusing the underlying backing slice array. Auto generated via dmlgen.`)
	mainGen.Pln(`func (cc *`, t.CollectionName(), `) Clear() *`, t.CollectionName(), ` {
	if cc == nil {
		*cc = `, t.CollectionName(), `{}
		return cc
	}
	if c := cap(cc.Data); c > len(cc.Data) { cc.Data = cc.Data[:c] }
	for i := 0; i < len(cc.Data); i++ {
		cc.Data[i] = nil
	}
	cc.Data = cc.Data[:0]
	return cc
}`)
}

func (t *Table) fnCollectionCut(mainGen *codegen.Go, g *Generator) {
	if !g.hasFeature(t.featuresInclude, t.featuresExclude, FeatureCollectionCut) {
		return
	}

	mainGen.C(`Cut will remove items i through j-1. Auto generated via dmlgen.`)
	mainGen.Pln(`func (cc *`, t.CollectionName(), `) Cut(i, j int) *`, t.CollectionName(), ` {`)
	{
		mainGen.In()
		mainGen.Pln(`z := cc.Data // copy slice header`)
		mainGen.Pln(`copy(z[i:], z[j:])`)
		mainGen.Pln(`for k, n := len(z)-j+i, len(z); k < n; k++ {`)
		{
			mainGen.In()
			mainGen.Pln(`z[k] = nil // this avoids the memory leak`)
			mainGen.Out()
		}
		mainGen.Pln(`}`)
		mainGen.Pln(`z = z[:len(z)-j+i]`)
		mainGen.Pln(`cc.Data = z`)
		mainGen.Pln(`return cc`)
		mainGen.Out()
	}
	mainGen.Pln(`}`)
}

func (t *Table) fnCollectionSwap(mainGen *codegen.Go, g *Generator) {
	if !g.hasFeature(t.featuresInclude, t.featuresExclude, FeatureCollectionSwap) {
		return
	}
	mainGen.C(`Swap will satisfy the sort.Interface. Auto generated via dmlgen.`)
	mainGen.Pln(`func (cc *`, t.CollectionName(), `) Swap(i, j int) { cc.Data[i], cc.Data[j] = cc.Data[j], cc.Data[i] }`)

	mainGen.C(`Len will satisfy the sort.Interface. Auto generated via dmlgen.`)
	mainGen.Pln(`func (cc *`, t.CollectionName(), `) Len() int { if cc == nil { return 0; }; return len(cc.Data); }`)
}

func (t *Table) fnCollectionDelete(mainGen *codegen.Go, g *Generator) {
	if !g.hasFeature(t.featuresInclude, t.featuresExclude, FeatureCollectionDelete) {
		return
	}

	mainGen.C(`Delete will remove an item from the slice. Auto generated via dmlgen.`)
	mainGen.Pln(`func (cc *`, t.CollectionName(), `) Delete(i int) *`, t.CollectionName(), ` {`)
	{
		mainGen.Pln(`z := cc.Data // copy the slice header`)
		mainGen.Pln(`end := len(z) - 1`)
		mainGen.Pln(`cc.Swap(i, end)`)
		mainGen.Pln(`copy(z[i:], z[i+1:])`)
		mainGen.Pln(`z[end] = nil // this should avoid the memory leak`)
		mainGen.Pln(`z = z[:end]`)
		mainGen.Pln(`cc.Data = z`)
		mainGen.Pln(`return cc`)
	}
	mainGen.Pln(`}`)
}

func (t *Table) fnCollectionInsert(mainGen *codegen.Go, g *Generator) {
	if !g.hasFeature(t.featuresInclude, t.featuresExclude, FeatureCollectionInsert) {
		return
	}
	mainGen.C(`Insert will place a new item at position i. Auto generated via dmlgen.`)
	mainGen.Pln(`func (cc *`, t.CollectionName(), `) Insert(n *`, t.EntityName(), `, i int) *`, t.CollectionName(), ` {`)
	{
		mainGen.Pln(`z := cc.Data // copy the slice header`)
		mainGen.Pln(`z = append(z, &`+t.EntityName(), `{})`)
		mainGen.Pln(`copy(z[i+1:], z[i:])`)
		mainGen.Pln(`z[i] = n`)
		mainGen.Pln(`cc.Data = z`)
		mainGen.Pln(`return cc`)
	}
	mainGen.Pln(`}`)
}

func (t *Table) fnCollectionAppend(mainGen *codegen.Go, g *Generator) {
	if !g.hasFeature(t.featuresInclude, t.featuresExclude, FeatureCollectionAppend) {
		return
	}
	mainGen.C(`Append will add a new item at the end of *`, t.CollectionName(), `. Auto generated via dmlgen.`)
	mainGen.Pln(`func (cc *`, t.CollectionName(), `) Append(n ...*`, t.EntityName(), `) *`, t.CollectionName(), ` {`)
	{
		mainGen.Pln(`cc.Data = append(cc.Data, n...)`)
		mainGen.Pln(`return cc`)
	}
	mainGen.Pln(`}`)
}

func (t *Table) fnCollectionBinaryMarshaler(mainGen *codegen.Go, g *Generator) {
	if !t.HasSerializer || !g.hasFeature(t.featuresInclude, t.featuresExclude, FeatureCollectionBinaryMarshaler) {
		return
	}

	mainGen.C(`UnmarshalBinary implements encoding.BinaryUnmarshaler.`)
	mainGen.Pln(`func (cc *`, t.CollectionName(), `) UnmarshalBinary(data []byte) error {`)
	{
		mainGen.Pln(`return cc.Unmarshal(data) // Implemented via github.com/gogo/protobuf`)
	}
	mainGen.Pln(`}`)

	mainGen.C(`MarshalBinary implements encoding.BinaryMarshaler.`)
	mainGen.Pln(`func (cc *`, t.CollectionName(), `) MarshalBinary() (data []byte, err error) {`)
	{
		mainGen.Pln(`return cc.Marshal()  // Implemented via github.com/gogo/protobuf`)
	}
	mainGen.Pln(`}`)
}

func (t *Table) fnEntityValidate(mainGen *codegen.Go, g *Generator) {
	if !g.hasFeature(t.featuresInclude, t.featuresExclude, FeatureEntityValidate) {
		return
	}

	fn, ok := g.customCode[t.EntityName()+".Validate"]

	if !ok {
		mainGen.C(`This variable can be set in another file to provide a custom validator.`)
		mainGen.Pln(`var validate`+t.EntityName(), ` func(*`, t.EntityName(), `) error `)
	}
	mainGen.C(`Validate runs internal consistency tests.`)
	mainGen.Pln(`func (e *`, t.EntityName(), `) Validate() error {`)
	{
		mainGen.In()
		mainGen.Pln(`if e == nil { return errors.NotValid.Newf("Type %T cannot be nil", e) }`)
		if ok {
			fn(g, t, mainGen)
		} else {
			mainGen.Pln(`if validate`+t.EntityName(), ` != nil { return validate`+t.EntityName(), `(e) }`)
		}

		mainGen.Out()
	}
	mainGen.Pln(`return nil }`)
}

func (t *Table) fnCollectionValidate(mainGen *codegen.Go, g *Generator) {
	if !g.hasFeature(t.featuresInclude, t.featuresExclude, FeatureCollectionValidate) {
		return
	}
	mainGen.C(`Validate runs internal consistency tests on all items.`)
	mainGen.Pln(`func (cc *`, t.CollectionName(), `) Validate() (err error) {`)
	{
		mainGen.In()
		mainGen.Pln(`if len(cc.Data) == 0 { return nil }`)
		mainGen.Pln(`for i,ld := 0, len(cc.Data); i < ld && err == nil; i++ {`)
		{
			mainGen.Pln(`err = cc.Data[i].Validate()`)
		}
		mainGen.Pln(`}`)
		mainGen.Out()
	}
	mainGen.Pln(`return }`)
}

func (t *Table) fnCollectionDBMapColumns(mainGen *codegen.Go, g *Generator) {
	if !g.hasFeature(t.featuresInclude, t.featuresExclude, FeatureDBMapColumns|
		FeatureDB|FeatureDBSelect|FeatureDBDelete|FeatureDBInsert|FeatureDBUpdate|FeatureDBUpsert) {
		return
	}

	mainGen.Pln(`func (cc *`, t.CollectionName(), `) scanColumns(cm *dml.ColumnMap, e *`, t.EntityName(), `) error {
			if err := e.MapColumns(cm); err != nil {
				return errors.WithStack(err)
			}
			// this function might get extended.
			return nil
		}`)

	mainGen.C(`MapColumns implements dml.ColumnMapper interface. Auto generated.`)
	mainGen.Pln(`func (cc *`, t.CollectionName(), `) MapColumns(cm *dml.ColumnMap) error {`)
	{
		mainGen.Pln(`switch m := cm.Mode(); m {
						case dml.ColumnMapEntityReadAll, dml.ColumnMapEntityReadSet:
							for _, e := range cc.Data {
								if err := cc.scanColumns(cm, e); err != nil {
									return errors.WithStack(err)
								}
							}`)

		mainGen.Pln(`case dml.ColumnMapScan:
							if cm.Count == 0 { cc.Clear(); }
							var e `, t.EntityName(), `
							if err := cc.scanColumns(cm, &e); err != nil {
								return errors.WithStack(err)
							}
							cc.Data = append(cc.Data, &e)`)

		unqiueCols := t.Table.Columns.UniqueColumns()
		hasUniqueCols := unqiueCols.Len() > 0
		mainGen.Pln(hasUniqueCols, `case dml.ColumnMapCollectionReadSet:
							for cm.Next(0) {
								switch c := cm.Column(); c {`)
		unqiueCols.Each(func(c *ddl.Column) {
			if !c.IsFloat() {
				mainGen.P(`case`, strconv.Quote(c.Field))
				for _, a := range c.Aliases {
					mainGen.P(`,`, strconv.Quote(a))
				}
				mainGen.Pln(`:`)
				mainGen.Pln(`cm = cm.`, g.goFuncNull(c)+`s(cc.`, strs.ToGoCamelCase(c.Field)+`s()...)`)
			}
		})
		mainGen.Pln(hasUniqueCols, `default:
				return errors.NotFound.Newf("[`+t.Package+`]`, t.CollectionName(), `Column %q not found", c)
			}
		} // end for cm.Next`)

		mainGen.Pln(`default:
		return errors.NotSupported.Newf("[` + t.Package + `] Unknown Mode: %q", string(m))
	}
	return cm.Err()`)
	}
	mainGen.Pln(`}`) // end func MapColumns
}

func (t *Table) fnCollectionDBMHandler(mainGen *codegen.Go, g *Generator) {
	if !g.hasFeature(t.featuresInclude, t.featuresExclude, FeatureDB|FeatureDBSelect|FeatureDBDelete|
		FeatureDBInsert|FeatureDBUpdate|FeatureDBUpsert) {
		return
	}

	var bufPKStructArg strings.Builder
	var bufPKNames strings.Builder
	i := 0

	tblPkCols := t.Table.Columns.PrimaryKeys()
	if t.Table.IsView() {
		tblPkCols = t.Table.Columns.ViewPrimaryKeys()
	}
	dbLoadStructArgOrSliceName := t.CollectionName() + "DBLoadArgs"
	bufPKStructArg.WriteString("type " + dbLoadStructArgOrSliceName + " struct {\n")
	bufPKStructArg.WriteString("\t_Named_Fields_Required struct{}\n")
	tblPkCols.Each(func(c *ddl.Column) {
		if i > 0 {
			bufPKNames.WriteByte(',')
		}
		goNamedField := strs.ToGoCamelCase(c.Field)
		if tblPkCols.Len() == 1 {
			dbLoadStructArgOrSliceName = g.goTypeNull(c)
		} else {
			bufPKStructArg.WriteString(goNamedField)
			bufPKStructArg.WriteByte(' ')
			bufPKStructArg.WriteString(g.goTypeNull(c) + "\n")
		}
		bufPKNames.WriteString(goNamedField)

		i++
	})
	bufPKStructArg.WriteString("}\n")
	if i == 0 {
		mainGen.C("The table/view", t.CollectionName(), "does not have a primary key. Skipping to generate DML functions based on the PK.")
		mainGen.Pln("\n")
		return
	}
	collectionPTRName := codegen.SkipWS("*", t.CollectionName())
	entityEventName := codegen.SkipWS(`event`, t.EntityName(), `Func`)
	tracingEnabled := t.hasFeature(g, FeatureDBTracing)
	collectionFuncName := codegen.SkipWS(t.CollectionName(), "SelectAll")
	dmlEnabled := t.hasFeature(g, FeatureDBSelect)

	mainGen.Pln(dmlEnabled && tblPkCols.Len() > 1, bufPKStructArg.String())

	mainGen.Pln(dmlEnabled, `func (cc `, collectionPTRName, `) DBLoad(ctx context.Context,dbm *DBM, pkIDs []`, dbLoadStructArgOrSliceName, `, opts ...dml.DBRFunc) (err error) {`)
	mainGen.Pln(dmlEnabled && tracingEnabled, `	ctx, span := dbm.option.Trace.Start(ctx, `, codegen.SkipWS(`"`, t.CollectionName(), "DBLoad", `"`), `)
		defer func(){ cstrace.Status(span, err, ""); span.End(); }()`)
	mainGen.Pln(dmlEnabled, `cc.Clear()`)
	mainGen.Pln(dmlEnabled, `qo := dml.FromContextQueryOptions(ctx)`)

	mainGen.Pln(dmlEnabled, `// put the IDs`, bufPKNames.String(), `into the context as value to search for a cache entry in the event function.
	if err = dbm.`, entityEventName, `(ctx, dml.EventFlagBeforeSelect, qo.SkipEvents, cc, nil); err != nil {
		return errors.WithStack(err)
	}
	if cc.Data != nil {
		return nil // might return data from cache
	}`)

	if tblPkCols.Len() > 1 { // for tables with more than one PK
		mainGen.Pln(dmlEnabled, `	cacheKey := `, codegen.SkipWS(`"`, collectionFuncName, "", `"`), `
	var args []interface{}
	if len(pkIDs) > 0 {
		args = make([]interface{}, 0, len(pkIDs)*`, tblPkCols.Len(), `)
		for _, pk := range pkIDs {`)
		tblPkCols.Each(func(c *ddl.Column) {
			mainGen.Pln(dmlEnabled, `args = append(args, pk.`, strs.ToGoCamelCase(c.Field), `)`)
		})
		mainGen.Pln(dmlEnabled, `}
		cacheKey = `, codegen.SkipWS(`"`, t.CollectionName(), "SelectByPK", `"`), `
	}
	if _, err = dbm.ConnPool.WithCacheKey(cacheKey, opts...).Load(ctx, cc, args...); err != nil {
		return errors.WithStack(err)
	}`)
	} else {
		mainGen.Pln(dmlEnabled, `if len(pkIDs) > 0 {`)
		mainGen.In()
		{
			mainGen.Pln(dmlEnabled, `if _, err = dbm.ConnPool.WithCacheKey(`, codegen.SkipWS(`"`, t.CollectionName(), "SelectByPK", `"`), `, opts...).Load(ctx, cc, pkIDs); err != nil {
		return errors.WithStack(err); }`)
		}
		mainGen.Out()
		mainGen.Pln(dmlEnabled, `} else {`)
		mainGen.In()
		{
			mainGen.Pln(dmlEnabled, `if _, err = dbm.ConnPool.WithCacheKey(`, codegen.SkipWS(`"`, collectionFuncName, "", `"`), `, opts...).Load(ctx, cc); err != nil {
		return errors.WithStack(err); }`)
		}
		mainGen.Out()
		mainGen.Pln(dmlEnabled, `}`)
	}

	mainGen.Pln(dmlEnabled, `return errors.WithStack(dbm.`, entityEventName, `(ctx, dml.EventFlagAfterSelect, qo.SkipEvents,cc, nil))
}`)

	if t.Table.IsView() {
		// skip here the delete,insert,update and upsert functions.
		return
	}

	dmlEnabled = t.hasFeature(g, FeatureDBDelete)
	collectionFuncName = codegen.SkipWS(t.EntityName(), "DeleteByPK")
	mainGen.Pln(dmlEnabled, `func (cc `, collectionPTRName, `) DBDelete(ctx context.Context, dbm *DBM, opts ...dml.DBRFunc) (res sql.Result,err error) {`)
	mainGen.Pln(dmlEnabled && tracingEnabled, `	ctx, span := dbm.option.Trace.Start(ctx, `, codegen.SkipWS(`"`, t.CollectionName(), "DeleteByPK", `"`), `)
			defer func(){ cstrace.Status(span, err, ""); span.End(); }()`)
	mainGen.Pln(dmlEnabled, `if cc == nil {
		return nil, errors.NotValid.Newf(`, codegen.SkipWS(`"`, t.CollectionName()), `can't be nil")
	}`)
	mainGen.Pln(dmlEnabled, `qo := dml.FromContextQueryOptions(ctx)`)

	mainGen.Pln(dmlEnabled, `if err = dbm.`, entityEventName, `(ctx, dml.EventFlagBeforeDelete, qo.SkipEvents, cc, nil); err != nil {
			return nil, errors.WithStack(err)
		}
		if res, err = dbm.ConnPool.WithCacheKey(`, codegen.SkipWS(`"`, collectionFuncName, `"`), `, opts...).ExecContext(ctx, dml.Qualify("", cc)); err != nil {
			return nil, errors.WithStack(err)
		}
		if err = errors.WithStack(dbm.`, entityEventName, `(ctx, dml.EventFlagAfterDelete, qo.SkipEvents,cc, nil)); err != nil {
			return nil, errors.WithStack(err)
		}
		return res, nil
	}`)

	dmlEnabled = t.hasFeature(g, FeatureDBUpdate)
	collectionFuncName = codegen.SkipWS(t.EntityName(), "UpdateByPK")
	mainGen.Pln(dmlEnabled, `func (cc `, collectionPTRName, `) DBUpdate(ctx context.Context, dbm *DBM, opts ...dml.DBRFunc) (err error) {`)
	mainGen.Pln(dmlEnabled && tracingEnabled, `	ctx, span := dbm.option.Trace.Start(ctx, `, codegen.SkipWS(`"`, t.CollectionName(), "UpdateByPK", `"`), `);
			defer func(){ cstrace.Status(span, err, ""); span.End(); }()`)
	mainGen.Pln(dmlEnabled, `if cc == nil {
		return errors.NotValid.Newf(`, codegen.SkipWS(`"`, t.CollectionName()), `can't be nil")
	}`)
	mainGen.Pln(dmlEnabled, `qo := dml.FromContextQueryOptions(ctx)`)

	mainGen.Pln(dmlEnabled, `if err = dbm.`, entityEventName, `(ctx, dml.EventFlagBeforeUpdate, qo.SkipEvents, cc, nil); err != nil {
			return errors.WithStack(err)
		}`)

	mainGen.Pln(dmlEnabled, `dbr := dbm.ConnPool.WithCacheKey(`, codegen.SkipWS(`"`, collectionFuncName, `"`), `, opts...)`)
	mainGen.Pln(dmlEnabled, `dbrStmt, err := dbr.Prepare(ctx)
		if err != nil {	return errors.WithStack(err) }`)

	mainGen.Pln(dmlEnabled, `for _, c := range cc.Data {
		res, err := dbrStmt.ExecContext(ctx, c)
		if err := dbr.ResultCheckFn(`, constTableName(t.Table.Name), `, 1, res, err); err != nil {
			return errors.WithStack(err)
		}
	}`)

	mainGen.Pln(dmlEnabled, `return errors.WithStack(dbm.`, entityEventName, `(ctx, dml.EventFlagAfterUpdate, qo.SkipEvents,cc, nil))
	}`)

	dmlEnabled = t.hasFeature(g, FeatureDBInsert)
	collectionFuncName = codegen.SkipWS(t.EntityName(), "Insert")
	mainGen.Pln(dmlEnabled, `func (cc `, collectionPTRName, `) DBInsert(ctx context.Context, dbm *DBM, opts ...dml.DBRFunc) (err error) {`)
	mainGen.Pln(dmlEnabled && tracingEnabled, `	ctx, span := dbm.option.Trace.Start(ctx, `, codegen.SkipWS(`"`, t.CollectionName(), "Insert", `"`), `);
			defer func(){ cstrace.Status(span, err, ""); span.End(); }()`)
	mainGen.Pln(dmlEnabled, `if cc == nil {
		return errors.NotValid.Newf(`, codegen.SkipWS(`"`, t.CollectionName()), `can't be nil")
	}`)
	mainGen.Pln(dmlEnabled, `qo := dml.FromContextQueryOptions(ctx)`)

	mainGen.Pln(dmlEnabled, `if err := dbm.`, entityEventName, `(ctx, dml.EventFlagBeforeInsert, qo.SkipEvents, cc, nil); err != nil {
			return errors.WithStack(err)
		}
		dbr := dbm.ConnPool.WithCacheKey(`, codegen.SkipWS(`"`, collectionFuncName, `"`), `, opts...)
		res, err := dbr.ExecContext(ctx, cc)
		if err := dbr.ResultCheckFn(`, constTableName(t.Table.Name), `, len(cc.Data), res, err); err != nil {
			return errors.WithStack(err)
		}
		return errors.WithStack(dbm.`, entityEventName, `(ctx, dml.EventFlagAfterInsert, qo.SkipEvents,cc, nil))
	}`)

	dmlEnabled = t.hasFeature(g, FeatureDBUpsert)
	collectionFuncName = codegen.SkipWS(t.EntityName(), "UpsertByPK")
	mainGen.Pln(dmlEnabled, `func (cc `, collectionPTRName, `) DBUpsert(ctx context.Context, dbm *DBM, opts ...dml.DBRFunc)  (err error) {`)
	mainGen.Pln(dmlEnabled && tracingEnabled, `	ctx, span := dbm.option.Trace.Start(ctx, `, codegen.SkipWS(`"`, t.CollectionName(), "UpsertByPK", `"`), `);
			defer func(){ cstrace.Status(span, err, ""); span.End(); }()`)
	mainGen.Pln(dmlEnabled, `if cc == nil {
		return errors.NotValid.Newf(`, codegen.SkipWS(`"`, t.CollectionName()), `can't be nil")
	}`)
	mainGen.Pln(dmlEnabled, `qo := dml.FromContextQueryOptions(ctx)`)

	mainGen.Pln(dmlEnabled, `if err := dbm.`, entityEventName, `(ctx, dml.EventFlagBeforeUpsert, qo.SkipEvents, cc, nil); err != nil {
			return errors.WithStack(err)
		}
		dbr := dbm.ConnPool.WithCacheKey(`, codegen.SkipWS(`"`, collectionFuncName, `"`), `, opts...)
		res, err := dbr.ExecContext(ctx, dml.Qualify("", cc))
		if err := dbr.ResultCheckFn(`, constTableName(t.Table.Name), `, len(cc.Data), res, err); err != nil {
				return errors.WithStack(err)
		}
		return errors.WithStack(dbm.`, entityEventName, `(ctx, dml.EventFlagAfterUpsert, qo.SkipEvents,cc, nil))
	}`)
}

func (t *Table) fnEntityDBMHandler(mainGen *codegen.Go, g *Generator) {
	if !g.hasFeature(t.featuresInclude, t.featuresExclude, FeatureDB|FeatureDBSelect|FeatureDBDelete|
		FeatureDBInsert|FeatureDBUpdate|FeatureDBUpsert) {
		return
	}

	var bufPKNamesAsArgs strings.Builder
	var bufPKStructArg strings.Builder
	var bufPKNames strings.Builder
	i := 0

	tblPkCols := t.Table.Columns.PrimaryKeys()
	if t.Table.IsView() {
		tblPkCols = t.Table.Columns.ViewPrimaryKeys()
	}
	dbLoadStructArgOrSliceName := t.EntityName() + "LoadArgs"
	bufPKStructArg.WriteString("type " + dbLoadStructArgOrSliceName + " struct {\n")
	bufPKStructArg.WriteString("\t_Named_Fields_Required struct{}\n")
	loadArgName := "arg"
	if tblPkCols.Len() == 1 {
		loadArgName = "primaryKey"
	}
	tblPkCols.Each(func(c *ddl.Column) {
		if i > 0 {
			bufPKNames.WriteByte(',')
			bufPKNamesAsArgs.WriteByte(',')
		}
		goNamedField := strs.ToGoCamelCase(c.Field)
		if tblPkCols.Len() == 1 {
			dbLoadStructArgOrSliceName = g.goTypeNull(c)
			bufPKNames.WriteString(loadArgName)
		} else {
			bufPKStructArg.WriteString(goNamedField)
			bufPKStructArg.WriteByte(' ')
			bufPKStructArg.WriteString(g.goTypeNull(c) + "\n")
			bufPKNames.WriteString(loadArgName + "." + goNamedField)
		}

		bufPKNamesAsArgs.WriteString("e.")
		bufPKNamesAsArgs.WriteString(strs.ToGoCamelCase(c.Field))
		i++
	})
	bufPKStructArg.WriteString("}\n")

	if i == 0 {
		mainGen.C("The table/view", t.EntityName(), "does not have a primary key. SKipping to generate DML functions based on the PK.")
		mainGen.Pln("\n")
		return
	}
	entityPTRName := codegen.SkipWS("*", t.EntityName())
	entityEventName := codegen.SkipWS(`event`, t.EntityName(), `Func`)
	tracingEnabled := t.hasFeature(g, FeatureDBTracing)
	entityFuncName := codegen.SkipWS(t.EntityName(), "SelectByPK")

	dmlEnabled := t.hasFeature(g, FeatureDBSelect)
	mainGen.Pln(dmlEnabled && tblPkCols.Len() > 1, bufPKStructArg.String())
	mainGen.Pln(dmlEnabled, `func (e `, entityPTRName, `) Load(ctx context.Context,dbm *DBM, `, loadArgName, ` `, dbLoadStructArgOrSliceName, `, opts ...dml.DBRFunc) (err error) {`)
	mainGen.Pln(dmlEnabled && tracingEnabled, `	ctx, span := dbm.option.Trace.Start(ctx, `, codegen.SkipWS(`"`, entityFuncName, `"`), `)
		defer func(){ cstrace.Status(span, err, ""); span.End(); }()`)
	mainGen.Pln(dmlEnabled, `if e == nil {
		return errors.NotValid.Newf(`, codegen.SkipWS(`"`, t.EntityName()), `can't be nil")
	}`)
	mainGen.Pln(dmlEnabled && len(t.availableRelationships) > 0, `e.setRelationParent()`)
	mainGen.Pln(dmlEnabled, `qo := dml.FromContextQueryOptions(ctx)`)

	mainGen.Pln(dmlEnabled, `// put the IDs`, bufPKNames.String(), `into the context as value to search for a cache entry in the event function.
	if err = dbm.`, entityEventName, `(ctx, dml.EventFlagBeforeSelect, qo.SkipEvents, nil, e); err != nil {
		return errors.WithStack(err)
	}
	if e.IsSet() {
		return nil // might return data from cache
	}
	if _, err = dbm.ConnPool.WithCacheKey(`, codegen.SkipWS(`"`, entityFuncName, `"`), `, opts...).Load(ctx, e, `, &bufPKNames, `); err != nil {
		return errors.WithStack(err)
	}
	return errors.WithStack(dbm.`, entityEventName, `(ctx, dml.EventFlagAfterSelect, qo.SkipEvents,nil, e))
}`)

	if t.Table.IsView() {
		// skip here the delete,insert,update and upsert functions.
		return
	}

	dmlEnabled = t.hasFeature(g, FeatureDBDelete)
	entityFuncName = codegen.SkipWS(t.EntityName(), "DeleteByPK")
	mainGen.Pln(dmlEnabled, `func (e `, entityPTRName, `) Delete(ctx context.Context, dbm *DBM, opts ...dml.DBRFunc) (res sql.Result, err error) {`)
	mainGen.Pln(dmlEnabled && tracingEnabled, `	ctx, span := dbm.option.Trace.Start(ctx, `, codegen.SkipWS(`"`, entityFuncName, `"`), `)
			defer func(){ cstrace.Status(span, err, ""); span.End(); }()`)
	mainGen.Pln(dmlEnabled, `if e == nil {
		return nil, errors.NotValid.Newf(`, codegen.SkipWS(`"`, t.EntityName()), `can't be nil")
	}`)
	mainGen.Pln(dmlEnabled && len(t.availableRelationships) > 0, `e.setRelationParent()`)
	mainGen.Pln(dmlEnabled, `qo := dml.FromContextQueryOptions(ctx)`)

	mainGen.Pln(dmlEnabled, `if err = dbm.`, entityEventName, `(ctx, dml.EventFlagBeforeDelete, qo.SkipEvents, nil, e); err != nil {
			return nil, errors.WithStack(err)
		}
		if res, err = dbm.ConnPool.WithCacheKey(`, codegen.SkipWS(`"`, entityFuncName, `"`), `, opts...).ExecContext(ctx, `, bufPKNamesAsArgs.String(), `); err != nil {
			return nil, errors.WithStack(err)
		}
		if err = dbm.`, entityEventName, `(ctx, dml.EventFlagAfterDelete, qo.SkipEvents,nil, e); err != nil {
			return nil, errors.WithStack(err)
		}
		return res, nil
	}`)

	dmlEnabled = t.hasFeature(g, FeatureDBUpdate)
	entityFuncName = codegen.SkipWS(t.EntityName(), "UpdateByPK")
	mainGen.Pln(dmlEnabled, `func (e `, entityPTRName, `) Update(ctx context.Context, dbm *DBM, opts ...dml.DBRFunc) (res sql.Result, err error) {`)
	mainGen.Pln(dmlEnabled && tracingEnabled, `	ctx, span := dbm.option.Trace.Start(ctx, `, codegen.SkipWS(`"`, entityFuncName, `"`), `);
			defer func(){ cstrace.Status(span, err, ""); span.End(); }()`)
	mainGen.Pln(dmlEnabled, `if e == nil {
		return nil, errors.NotValid.Newf(`, codegen.SkipWS(`"`, t.EntityName()), `can't be nil")
	}`)
	mainGen.Pln(dmlEnabled && len(t.availableRelationships) > 0, `e.setRelationParent()`)
	mainGen.Pln(dmlEnabled, `qo := dml.FromContextQueryOptions(ctx)`)

	mainGen.Pln(dmlEnabled, `if err = dbm.`, entityEventName, `(ctx, dml.EventFlagBeforeUpdate, qo.SkipEvents, nil, e); err != nil {
			return nil, errors.WithStack(err)
		}
		if res, err = dbm.ConnPool.WithCacheKey(`, codegen.SkipWS(`"`, entityFuncName, `"`), `, opts...).ExecContext(ctx, e); err != nil {
			return nil, errors.WithStack(err)
		}
		if err = dbm.`, entityEventName, `(ctx, dml.EventFlagAfterUpdate, qo.SkipEvents,nil, e); err != nil {
			return nil, errors.WithStack(err)
		}
		return res, nil
	}`)

	dmlEnabled = t.hasFeature(g, FeatureDBInsert)
	entityFuncName = codegen.SkipWS(t.EntityName(), "Insert")
	mainGen.Pln(dmlEnabled, `func (e `, entityPTRName, `) Insert(ctx context.Context, dbm *DBM, opts ...dml.DBRFunc) (res sql.Result, err error) {`)
	mainGen.Pln(dmlEnabled && tracingEnabled, `	ctx, span := dbm.option.Trace.Start(ctx, `, codegen.SkipWS(`"`, entityFuncName, `"`), `);
			defer func(){ cstrace.Status(span, err, ""); span.End(); }()`)
	mainGen.Pln(dmlEnabled, `if e == nil {
		return nil, errors.NotValid.Newf(`, codegen.SkipWS(`"`, t.EntityName()), `can't be nil")
	}`)
	mainGen.Pln(dmlEnabled && len(t.availableRelationships) > 0, `e.setRelationParent()`)
	mainGen.Pln(dmlEnabled, `qo := dml.FromContextQueryOptions(ctx)`)

	mainGen.Pln(dmlEnabled, `if err = dbm.`, entityEventName, `(ctx, dml.EventFlagBeforeInsert, qo.SkipEvents, nil, e); err != nil {
			return nil, errors.WithStack(err)
		}
		if res, err = dbm.ConnPool.WithCacheKey(`, codegen.SkipWS(`"`, entityFuncName, `"`), `, opts...).ExecContext(ctx, e); err != nil {
			return nil, errors.WithStack(err)
		}
		if err = dbm.`, entityEventName, `(ctx, dml.EventFlagAfterInsert, qo.SkipEvents,nil, e); err != nil {
			return nil, errors.WithStack(err)
		}
		return res, nil
	}`)

	dmlEnabled = t.hasFeature(g, FeatureDBUpsert)
	entityFuncName = codegen.SkipWS(t.EntityName(), "UpsertByPK")
	mainGen.Pln(dmlEnabled, `func (e `, entityPTRName, `) Upsert(ctx context.Context, dbm *DBM, opts ...dml.DBRFunc) (res sql.Result, err error) {`)
	mainGen.Pln(dmlEnabled && tracingEnabled, `	ctx, span := dbm.option.Trace.Start(ctx, `, codegen.SkipWS(`"`, entityFuncName, `"`), `);
			defer func(){ cstrace.Status(span, err, ""); span.End(); }()`)
	mainGen.Pln(dmlEnabled, `if e == nil {
		return nil, errors.NotValid.Newf(`, codegen.SkipWS(`"`, t.EntityName()), `can't be nil")
	}`)
	mainGen.Pln(dmlEnabled && len(t.availableRelationships) > 0, `e.setRelationParent()`)
	mainGen.Pln(dmlEnabled, `qo := dml.FromContextQueryOptions(ctx)`)

	mainGen.Pln(dmlEnabled, `if err = dbm.`, entityEventName, `(ctx, dml.EventFlagBeforeUpsert, qo.SkipEvents, nil, e); err != nil {
			return nil, errors.WithStack(err)
		}
		if res, err = dbm.ConnPool.WithCacheKey(`, codegen.SkipWS(`"`, entityFuncName, `"`), `, opts...).ExecContext(ctx, dml.Qualify("", e)); err != nil {
			return nil, errors.WithStack(err)
		}
		if err = dbm.`, entityEventName, `(ctx, dml.EventFlagAfterUpsert, qo.SkipEvents,nil, e); err != nil {
			return nil, errors.WithStack(err)
		}
		return res, nil
	}`)
}

func (t *Table) fnDBMOptionsSQLBuildQueries(mainGen *codegen.Go, g *Generator) {
	tblPKLen := t.Table.Columns.PrimaryKeys().Len()
	tblPK := t.Table.Columns.PrimaryKeys()
	if t.Table.IsView() {
		tblPKLen = t.Table.Columns.ViewPrimaryKeys().Len()
		tblPK = t.Table.Columns.ViewPrimaryKeys()
	}

	var pkWhereIN strings.Builder
	var pkWhereEQ strings.Builder
	if tblPKLen == 1 {
		pkWhereIN.WriteString("\ndml.Column(`" + strings.Join(tblPK.FieldNames(), "`,`") + "`).In().")
		pkWhereIN.WriteString("PlaceHolder(),\n")
		pkWhereEQ.WriteString("\ndml.Column(`" + strings.Join(tblPK.FieldNames(), "`,`") + "`).Equal().")
		pkWhereEQ.WriteString("PlaceHolder(),\n")
	} else {
		pkWhereIN.WriteString("\ndml.Columns(`" + strings.Join(tblPK.FieldNames(), "`,`") + "`).In().")
		pkWhereIN.WriteString("Tuples(),\n")
		pkWhereEQ.WriteString("\ndml.Columns(`" + strings.Join(tblPK.FieldNames(), "`,`") + "`).Equal().")
		pkWhereEQ.WriteString("Tuples(),\n")
	}

	mainGen.Pln(tblPKLen > 0 && t.hasFeature(g, FeatureDBSelect|FeatureCollectionStruct),
		codegen.SkipWS(`"`, t.CollectionName(), `SelectAll"`),
		`: dbmo.InitSelectFn(tbls.MustTable(`, constTableName(t.Table.Name), `).Select("*")),`)

	mainGen.Pln(tblPKLen > 0 && t.hasFeature(g, FeatureDBSelect|FeatureEntityStruct|FeatureCollectionStruct),
		codegen.SkipWS(`"`, t.CollectionName(), `SelectByPK"`),
		`: dbmo.InitSelectFn(tbls.MustTable(`, constTableName(t.Table.Name), `).Select("*")).Where(`, pkWhereIN.String(), `),`)

	mainGen.Pln(tblPKLen > 0 && t.hasFeature(g, FeatureDBSelect|FeatureEntityStruct|FeatureCollectionStruct),
		codegen.SkipWS(`"`, t.EntityName(), `SelectByPK"`),
		`: dbmo.InitSelectFn(tbls.MustTable(`, constTableName(t.Table.Name), `).Select("*")).Where(`, pkWhereEQ.String(), `),`)

	if t.Table.IsView() {
		return
	}

	mainGen.Pln(t.hasFeature(g, FeatureDBUpdate|FeatureEntityStruct|FeatureCollectionStruct),
		codegen.SkipWS(`"`, t.EntityName(), `UpdateByPK"`),
		`: dbmo.InitUpdateFn(tbls.MustTable(`, constTableName(t.Table.Name), `).Update().Where(`, pkWhereEQ.String(), `)),`)
	mainGen.Pln(t.hasFeature(g, FeatureDBDelete|FeatureEntityStruct|FeatureCollectionStruct),
		codegen.SkipWS(`"`, t.EntityName(), `DeleteByPK"`),
		`: dbmo.InitDeleteFn(tbls.MustTable(`, constTableName(t.Table.Name), `).Delete().Where(`, pkWhereIN.String(), `)),`)
	mainGen.Pln(t.hasFeature(g, FeatureDBInsert|FeatureEntityStruct|FeatureCollectionStruct),
		codegen.SkipWS(`"`, t.EntityName(), `Insert"`),
		`: dbmo.InitInsertFn(tbls.MustTable(`, constTableName(t.Table.Name), `).Insert()),`)
	mainGen.Pln(t.hasFeature(g, FeatureDBUpsert|FeatureEntityStruct|FeatureCollectionStruct),
		codegen.SkipWS(`"`, t.EntityName(), `UpsertByPK"`),
		`: dbmo.InitInsertFn(tbls.MustTable(`, constTableName(t.Table.Name), `).Insert()).OnDuplicateKey(),`)

	// foreign keys
	if g.hasFeature(t.featuresInclude, t.featuresExclude, FeatureEntityRelationships) && len(t.availableRelationships) > 0 {
		mainGen.C(`<FOREIGN_KEY_QUERIES`, t.Table.Name, `>`)
		var fkWhereEQ bytes.Buffer
		for _, rs := range t.availableRelationships {
			fkWhereEQ.WriteString("\ndml.Column(`" + rs.columnName + "`).Equal().PlaceHolder(),\n")

			// DELETE FROM
			mainGen.Pln(codegen.SkipWS(`"`, rs.mappedStructFieldName, `DeleteByFK"`),
				`: dbmo.InitDeleteFn(tbls.MustTable(`, constTableName(rs.tableName), `).Delete().Where(`, fkWhereEQ.String(), `)),`)

			// SELECT FROM
			mainGen.Pln(codegen.SkipWS(`"`, rs.mappedStructFieldName, `SelectByFK"`),
				`: dbmo.InitSelectFn(tbls.MustTable(`, constTableName(rs.tableName), `).Select("*").Where(`, fkWhereEQ.String(), `)),`)

			// UPDATE not needed as it uses the default UPDATE

			fkWhereEQ.Reset()
		}
		mainGen.C(`</FOREIGN_KEY_QUERIES`, t.Table.Name, `>`)
	}
}

func (t *Table) generateTestOther(testGen *codegen.Go, g *Generator) (codeWritten int) {
	if g.hasFeature(t.featuresInclude, t.featuresExclude, FeatureEntityEmpty) {
		testGen.Pln(`t.Run("` + t.EntityName() + `_Empty", func(t *testing.T) {`)
		{
			testGen.Pln(`e:= new(`, t.EntityName(), `)`)
			testGen.Pln(`assert.NoError(t, ps.FakeData(e))`)
			testGen.Pln(`e.Empty()`)
			testGen.Pln(`assert.Exactly(t, *e, `, t.EntityName(), `{})`)
		}
		testGen.Pln(`})`) // end t.Run
		codeWritten++
	}
	if g.hasFeature(t.featuresInclude, t.featuresExclude, FeatureEntityCopy) {
		testGen.Pln(`t.Run("` + t.EntityName() + `_Copy", func(t *testing.T) {`)
		{
			testGen.Pln(`e:= new(`, t.EntityName(), `)`)
			testGen.Pln(`assert.NoError(t, ps.FakeData(e))`)
			testGen.Pln(`e2 := e.Copy()`)
			testGen.Pln(`assert.Exactly(t, e, e2)`)
			testGen.Pln(`assert.NoError(t, ps.FakeData(e))`)
			testGen.Pln(`assert.NotEqual(t, e, e2)`)
		}
		testGen.Pln(`})`) // end t.Run
		codeWritten++
	}
	if g.hasFeature(t.featuresInclude, t.featuresExclude, FeatureEntityValidate) {
		testGen.Pln(`t.Run("` + t.CollectionName() + `_Validate", func(t *testing.T) {`)
		{
			testGen.Pln(`c := `, t.CollectionName(), `{ Data: []*`, t.EntityName(), `{nil} }`)
			testGen.Pln(`assert.True(t, errors.NotValid.Match(c.Validate()))`)
		}
		testGen.Pln(`})`) // end t.Run
		codeWritten++
	}
	// more feature tests to follow
	return
}

func (t *Table) generateTestDB(testGen *codegen.Go) {
	testGen.Pln(`t.Run("` + strs.ToGoCamelCase(t.Table.Name) + `_Entity", func(t *testing.T) {`)
	testGen.Pln(`tbl := tbls.MustTable(TableName`+strs.ToGoCamelCase(t.Table.Name), `)`)

	testGen.Pln(`selOneRow := tbl.Select("*").Where(`)
	for _, c := range t.Table.Columns {
		if c.IsPK() && c.IsAutoIncrement() {
			testGen.Pln(`dml.Column(`, strconv.Quote(c.Field), `).Equal().PlaceHolder(),`)
		}
	}
	testGen.Pln(`)`)

	testGen.Pln(`selTenRows := tbl.Select("*").Where(`)
	for _, c := range t.Table.Columns {
		if c.IsPK() && c.IsAutoIncrement() {
			testGen.Pln(`dml.Column(`, strconv.Quote(c.Field), `).LessOrEqual().Int(10),`)
		}
	}
	testGen.Pln(`)`)

	testGen.Pln(`selOneRowDBR := tbls.ConnPool.WithPrepare(ctx, selOneRow)`)
	testGen.Pln(`defer selOneRowDBR.Close()`)
	testGen.Pln(`selTenRowsDBR := tbls.ConnPool.WithQueryBuilder(selTenRows)`)

	if t.HasAutoIncrement < 2 {
		testGen.C(`this table/view does not support auto_increment`)
		testGen.Pln(`entCol := New`+t.CollectionName(), `()`)
		testGen.Pln(`rowCount, err := selTenRowsDBR.Load(ctx, entCol)`)
		testGen.Pln(`assert.NoError(t, err)`)
		testGen.Pln(`t.Logf("Collection load rowCount: %d", rowCount)`)
	} else {
		testGen.Pln(`entINSERTStmtA := tbls.ConnPool.WithPrepare(ctx,tbl.Insert().BuildValues())`)

		testGen.Pln(`for i := 0; i < 9; i++ {`)
		{
			testGen.In()
			testGen.Pln(`entIn := new(`, strs.ToGoCamelCase(t.Table.Name), `)`)
			testGen.Pln(`assert.NoError(t, ps.FakeData(entIn), "Error at index %d", i)`)

			testGen.Pln(`lID := dmltest.CheckLastInsertID(t, "Error: TestNewTables.` + strs.ToGoCamelCase(t.Table.Name) + `_Entity")(entINSERTStmtA.ExecContext(ctx,dml.Qualify("", entIn)))`)
			testGen.Pln(`entINSERTStmtA.Reset()`)

			testGen.Pln(`entOut := new(`, strs.ToGoCamelCase(t.Table.Name), `)`)
			testGen.Pln(`rowCount, err := selOneRowDBR.Load(ctx, entOut, lID)`)
			testGen.Pln(`assert.NoError(t, err)`)
			testGen.Pln(`assert.Exactly(t, uint64(1), rowCount, "IDX%d: RowCount did not match", i)`)

			for _, c := range t.Table.Columns {
				fn := t.GoCamelMaybePrivate(c.Field)
				switch {
				case c.IsTime():
					// skip comparison as we can't mock time (yet) :-(
				case c.IsChar():
					testGen.Pln(`assert.ExactlyLength(t,`, c.CharMaxLength.Int64, `, `, `&entIn.`, fn, `,`, `&entOut.`, fn, `,`, `"IDX%d:`, fn, `should match", lID)`)
				case !c.IsSystemVersioned():
					testGen.Pln(`assert.Exactly(t, entIn.`, fn, `,`, `entOut.`, fn, `,`, `"IDX%d:`, fn, `should match", lID)`)
				default:
					testGen.C(`ignoring:`, c.Field)
				}
			}
			testGen.Out()
		}
		testGen.Pln(`}`) // endfor
		testGen.Pln(`dmltest.Close(t, entINSERTStmtA)`)
		testGen.Pln(`entCol := New`+t.CollectionName(), `()`)
		testGen.Pln(`rowCount, err := selTenRowsDBR.Load(ctx, entCol)`)
		testGen.Pln(`assert.NoError(t, err)`)
		testGen.Pln(`t.Logf("Collection load rowCount: %d", rowCount)`)

		testGen.Pln(`colInsertDBR := tbls.ConnPool.WithQueryBuilder(tbl.Insert().Replace().SetRowCount(len(entCol.Data)).BuildValues())`)
		testGen.Pln(`lID := dmltest.CheckLastInsertID(t, "Error: `, t.CollectionName(), `")(colInsertDBR.ExecContext(ctx, dml.Qualify("", entCol)))`)
		testGen.Pln(`t.Logf("Last insert ID into: %d", lID)`)
	}

	testGen.Pln(`})`)
}
