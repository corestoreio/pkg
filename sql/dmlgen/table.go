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
	"crypto/md5"
	"fmt"
	"strconv"
	"strings"
	"text/tabwriter"

	"github.com/corestoreio/pkg/sql/ddl"
	"github.com/corestoreio/pkg/util/codegen"
	"github.com/corestoreio/pkg/util/strs"
)

// FeatureToggle allows certain generated code blocks to be switched off or on.
type FeatureToggle uint64

// List of available features
const (
	FeatureCollectionAppend FeatureToggle = 1 << iota
	FeatureCollectionBinaryMarshaler
	FeatureCollectionCut
	FeatureCollectionDBMapColumns
	FeatureCollectionDelete
	FeatureCollectionEach
	FeatureCollectionFilter
	FeatureCollectionInsert
	FeatureCollectionStruct // creates the struct type
	FeatureCollectionSwap
	FeatureCollectionUniqueGetters
	FeatureCollectionUniquifiedGetters
	FeatureCollectionValidate
	FeatureDB
	FeatureEntityCopy
	FeatureEntityDBAssignLastInsertID
	FeatureEntityDBDelete
	FeatureEntityDBInsert
	FeatureEntityDBMapColumns
	FeatureEntityDBSelect
	FeatureEntityDBTracing // opentelemetry tracing
	FeatureEntityDBUpdate
	FeatureEntityDBUpsert
	FeatureEntityEmpty
	FeatureEntityGetSetPrivateFields
	FeatureEntityIsSet
	FeatureEntityRelationships
	FeatureEntityStruct // creates the struct type
	FeatureEntityValidate
	FeatureEntityWriteTo
	featureMax
)

var featureNames = map[FeatureToggle]string{
	FeatureCollectionAppend:            "FeatureCollectionAppend",
	FeatureCollectionBinaryMarshaler:   "FeatureCollectionBinaryMarshaler",
	FeatureCollectionCut:               "FeatureCollectionCut",
	FeatureCollectionDBMapColumns:      "FeatureCollectionDBMapColumns",
	FeatureCollectionDelete:            "FeatureCollectionDelete",
	FeatureCollectionEach:              "FeatureCollectionEach",
	FeatureCollectionFilter:            "FeatureCollectionFilter",
	FeatureCollectionInsert:            "FeatureCollectionInsert",
	FeatureCollectionStruct:            "FeatureCollectionStruct",
	FeatureCollectionSwap:              "FeatureCollectionSwap",
	FeatureCollectionUniqueGetters:     "FeatureCollectionUniqueGetters",
	FeatureCollectionUniquifiedGetters: "FeatureCollectionUniquifiedGetters",
	FeatureCollectionValidate:          "FeatureCollectionValidate",
	FeatureDB:                          "FeatureDB",
	FeatureEntityCopy:                  "FeatureEntityCopy",
	FeatureEntityDBAssignLastInsertID:  "FeatureEntityDBAssignLastInsertID",
	FeatureEntityDBDelete:              "FeatureEntityDBDelete",
	FeatureEntityDBInsert:              "FeatureEntityDBInsert",
	FeatureEntityDBMapColumns:          "FeatureEntityDBMapColumns",
	FeatureEntityDBSelect:              "FeatureEntityDBSelect",
	FeatureEntityDBTracing:             "FeatureEntityDBTracing",
	FeatureEntityDBUpdate:              "FeatureEntityDBUpdate",
	FeatureEntityDBUpsert:              "FeatureEntityDBUpsert",
	FeatureEntityEmpty:                 "FeatureEntityEmpty",
	FeatureEntityGetSetPrivateFields:   "FeatureEntityGetSetPrivateFields",
	FeatureEntityIsSet:                 "FeatureEntityIsSet",
	FeatureEntityRelationships:         "FeatureEntityRelationships",
	FeatureEntityStruct:                "FeatureEntityStruct",
	FeatureEntityValidate:              "FeatureEntityValidate",
	FeatureEntityWriteTo:               "FeatureEntityWriteTo",
}

func (f FeatureToggle) String() string {
	var buf strings.Builder
	j := 0
	for i, k := 0, 0; i <= int(featureMax); i, k = 1<<k, k+1 {
		if int(f)&i > 0 {
			if j > 0 {
				buf.WriteByte(',')
			}
			buf.WriteString(featureNames[FeatureToggle(i)])
			j++
		}
	}
	return buf.String()
}

type tables []*Table

// hasFeature returns false when any of the tables does not have the feature/s.
func (ts tables) hasFeature(g *Generator, feature FeatureToggle) bool {
	for _, tbl := range ts {
		if g.hasFeature(tbl.featuresInclude, tbl.featuresExclude, feature) {
			return true
		}
	}
	return false
}

func (ts tables) names() []string {
	names := make([]string, len(ts))
	for i, tbl := range ts {
		names[i] = tbl.Table.Name
	}
	return names
}

// nameID returns a consistent md5 hash of the table names.
func (ts tables) nameID() string {
	var buf bytes.Buffer
	for _, tbl := range ts {
		buf.WriteString(tbl.Table.Name)
	}
	return fmt.Sprintf("%x", md5.Sum(buf.Bytes()))
}

// table writes one database table into Go source code.
type Table struct {
	Package              string // Name of the package
	Table                *ddl.Table
	Comment              string // Comment above the struct type declaration
	HasAutoIncrement     uint8  // 0=nil,1=false (has NO auto increment),2=true has auto increment
	HasEasyJSONMarshaler bool
	HasSerializer        bool // writes the .proto file if true

	// PrivateFields key=snake case name of the DB column, value=true, the field must be private
	debug                 bool // gets set via isDebug function
	privateFields         map[string]bool
	featuresInclude       FeatureToggle
	featuresExclude       FeatureToggle
	fieldMapFn            func(dbIdentifier string) (newName string)
	customStructTagFields map[string]string
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
	return lcFirst(su)
}

func (t *Table) CollectionName() string {
	return collectionName(t.Table.Name)
}

func (t *Table) EntityName() string {
	return strs.ToGoCamelCase(t.Table.Name)
}

func (t *Table) hasFeature(g *Generator, f FeatureToggle) bool {
	return g.hasFeature(t.featuresInclude, t.featuresExclude, f)
}

func (t *Table) collectionStruct(mainGen *codegen.Go, g *Generator) {
	if !g.hasFeature(t.featuresInclude, t.featuresExclude, FeatureCollectionStruct) {
		return
	}

	mainGen.C(t.CollectionName(), `represents a collection type for DB table`, t.Table.Name)
	mainGen.C(`Not thread safe. Auto generated.`)
	if t.Comment != "" {
		mainGen.C(t.Comment)
	}
	if t.HasEasyJSONMarshaler {
		mainGen.Pln(`//easyjson:json`) // do not use C() because it adds a whitespace between "//" and "e"
	}
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

func (t *Table) entityStruct(mainGen *codegen.Go, g *Generator) {
	if !g.hasFeature(t.featuresInclude, t.featuresExclude, FeatureEntityStruct) {
		return
	}

	fieldMapFn := g.defaultTableConfig.FieldMapFn
	if fieldMapFn == nil {
		fieldMapFn = t.fieldMapFn
	}
	if fieldMapFn == nil {
		fieldMapFn = defaultFieldMapFn
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

			// this part is duplicated in the proto file generation function generateProto.
			if g.hasFeature(t.featuresInclude, t.featuresExclude, FeatureEntityRelationships) {
				var debugBuf bytes.Buffer
				tabW := tabwriter.NewWriter(&debugBuf, 6, 0, 2, ' ', 0)
				var hasAtLeastOneRelationShip int
				fmt.Fprintf(&debugBuf, "RelationInfo for: %q\n", t.Table.Name)
				fmt.Fprintf(tabW, "Case\tis1:M\tis1:1\tseen?\tisRelAl\thasTable\tTarget Tbl M:N\tRelation\n")
				relationShipSeen := map[string]bool{}
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
							mainGen.Pln(fieldMapFn(collectionName(kcuce.ReferencedTableName.Data)), " *", collectionName(kcuce.ReferencedTableName.Data),
								t.customStructTagFields[kcuce.ReferencedTableName.Data],
								"// 1:M", kcuce.TableName+"."+kcuce.ColumnName, "=>", kcuce.ReferencedTableName.Data+"."+kcuce.ReferencedColumnName.Data)
						}

						// case ONE-TO-ONE
						isOneToOne := g.krs.IsOneToOne(kcuce.TableName, kcuce.ColumnName, kcuce.ReferencedTableName.Data, kcuce.ReferencedColumnName.Data)
						fmt.Fprintf(tabW, "B1_1:1\t%t\t%t\t%t\t%t\t%t\t-\t%s => %s\n", isOneToMany, isOneToOne, false, isRelationAllowed, hasTable,
							kcuce.TableName+"."+kcuce.ColumnName, kcuce.ReferencedTableName.Data+"."+kcuce.ReferencedColumnName.Data)
						if isOneToOne && hasTable && isRelationAllowed {
							mainGen.Pln(fieldMapFn(strs.ToGoCamelCase(kcuce.ReferencedTableName.Data)), " *", strs.ToGoCamelCase(kcuce.ReferencedTableName.Data),
								t.customStructTagFields[kcuce.ReferencedTableName.Data],
								"// 1:1", kcuce.TableName+"."+kcuce.ColumnName, "=>", kcuce.ReferencedTableName.Data+"."+kcuce.ReferencedColumnName.Data)
						}

						// case MANY-TO-MANY
						targetTbl, targetColumn := g.krs.ManyToManyTarget(kcuce.TableName, kcuce.ColumnName)
						fmt.Fprintf(tabW, "C1_M:N\t%t\t%t\t%t\t%t\t%t\t%s\t%s => %s\n", isOneToMany, isOneToOne, false, isRelationAllowed, hasTable,
							targetTbl+"."+targetColumn,
							kcuce.TableName+"."+kcuce.ColumnName, kcuce.ReferencedTableName.Data+"."+kcuce.ReferencedColumnName.Data)
						// hasTable variable shall not be added because usually the link table does not get loaded.
						if isRelationAllowed && targetTbl != "" && targetColumn != "" {
							mainGen.Pln(fieldMapFn(collectionName(targetTbl)), " *", collectionName(targetTbl),
								t.customStructTagFields[targetTbl],
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
						keySeen := fieldMapFn(collectionName(kcuce.ReferencedTableName.Data))
						relationShipSeenAlready := relationShipSeen[keySeen]
						// case ONE-TO-MANY
						fmt.Fprintf(tabW, "A2_1:M rev\t%t\t%t\t%t\t%t\t%t\t-\t%s => %s\n", isOneToMany, false, relationShipSeenAlready, isRelationAllowed, hasTable,
							kcuce.TableName+"."+kcuce.ColumnName, kcuce.ReferencedTableName.Data+"."+kcuce.ReferencedColumnName.Data)
						if isRelationAllowed && isOneToMany && hasTable && !relationShipSeenAlready {
							mainGen.Pln(fieldMapFn(collectionName(kcuce.ReferencedTableName.Data)), " *", collectionName(kcuce.ReferencedTableName.Data),
								t.customStructTagFields[kcuce.ReferencedTableName.Data],
								"// Reversed 1:M", kcuce.TableName+"."+kcuce.ColumnName, "=>", kcuce.ReferencedTableName.Data+"."+kcuce.ReferencedColumnName.Data)
							relationShipSeen[keySeen] = true
						}

						// case ONE-TO-ONE
						isOneToOne := g.krs.IsOneToOne(kcuce.TableName, kcuce.ColumnName, kcuce.ReferencedTableName.Data, kcuce.ReferencedColumnName.Data)
						fmt.Fprintf(tabW, "B2_1:1 rev\t%t\t%t\t%t\t%t\t%t\t-\t%s => %s\n", isOneToMany, isOneToOne, relationShipSeenAlready, isRelationAllowed, hasTable,
							kcuce.TableName+"."+kcuce.ColumnName, kcuce.ReferencedTableName.Data+"."+kcuce.ReferencedColumnName.Data)
						if isRelationAllowed && isOneToOne && hasTable {
							mainGen.Pln(fieldMapFn(strs.ToGoCamelCase(kcuce.ReferencedTableName.Data)), " *", strs.ToGoCamelCase(kcuce.ReferencedTableName.Data),
								t.customStructTagFields[kcuce.ReferencedTableName.Data],
								"// Reversed 1:1", kcuce.TableName+"."+kcuce.ColumnName, "=>", kcuce.ReferencedTableName.Data+"."+kcuce.ReferencedColumnName.Data)
						}

						// case MANY-TO-MANY
						targetTbl, targetColumn := g.krs.ManyToManyTarget(kcuce.ReferencedTableName.Data, kcuce.ReferencedColumnName.Data)
						if targetTbl != "" && targetColumn != "" {
							keySeen := fieldMapFn(collectionName(targetTbl))
							isRelationAllowed = g.isAllowedRelationship(kcuce.TableName, kcuce.ColumnName, targetTbl, targetColumn) &&
								!relationShipSeen[keySeen]
							relationShipSeen[keySeen] = true
						}

						// case MANY-TO-MANY
						fmt.Fprintf(tabW, "C2_M:N rev\t%t\t%t\t%t\t%t\t%t\t%s\t%s => %s\n", isOneToMany, isOneToOne, relationShipSeenAlready, isRelationAllowed, hasTable,
							targetTbl+"."+targetColumn,
							kcuce.TableName+"."+kcuce.ColumnName, kcuce.ReferencedTableName.Data+"."+kcuce.ReferencedColumnName.Data)
						// hasTable shall not be added because usually the link table does not get loaded.
						if isRelationAllowed && targetTbl != "" && targetColumn != "" {
							mainGen.Pln(fieldMapFn(collectionName(targetTbl)), " *", collectionName(targetTbl),
								t.customStructTagFields[targetTbl],
								"// Reversed M:N", kcuce.TableName+"."+kcuce.ColumnName, "via", kcuce.ReferencedTableName.Data+"."+kcuce.ReferencedColumnName.Data,
								"=>", targetTbl+"."+targetColumn,
							)
						}
					}
				}
				if t.debug && hasAtLeastOneRelationShip > 0 {
					_ = tabW.Flush()
					fmt.Fprintf(&debugBuf, "Relationship count: %d\n", hasAtLeastOneRelationShip)
					fmt.Println(debugBuf.String())
				}
			}
			mainGen.Out()
		}
	}
	mainGen.Pln(`}`)
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
	if !g.hasFeature(t.featuresInclude, t.featuresExclude, FeatureEntityIsSet|FeatureDB|FeatureEntityDBSelect) {
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
		e2 := new(`, t.EntityName(), `)
		*e2 = *e // for now a shallow copy
		return e2
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
	if !g.hasFeature(t.featuresInclude, t.featuresExclude, FeatureEntityDBMapColumns|
		FeatureDB|FeatureEntityDBSelect|FeatureEntityDBDelete|
		FeatureEntityDBInsert|FeatureEntityDBUpdate|FeatureEntityDBUpsert) {
		return
	}
	mainGen.C(`MapColumns implements interface ColumnMapper only partially. Auto generated.`)
	mainGen.Pln(`func (e *`, t.EntityName(), `) MapColumns(cm *dml.ColumnMap) error {`)
	{
		if fn, ok := g.customCode["func_"+t.EntityName()+"_MapColumns"]; ok {
			fn(g, t, mainGen)
		}

		mainGen.In()
		mainGen.Pln(`if cm.Mode() == dml.ColumnMapEntityReadAll {`)
		{
			mainGen.In()
			mainGen.P(`return cm`)
			t.Table.Columns.Each(func(c *ddl.Column) {
				mainGen.P(`.`, g.goFuncNull(c), `(&e.`, t.GoCamelMaybePrivate(c.Field), `)`)
			})
			mainGen.Pln(`.Err()`)
			mainGen.Out()
		}
		mainGen.Pln(`}`)
		mainGen.Pln(`for cm.Next() {`)
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
					mainGen.Pln(`:`)
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
	if !g.hasFeature(t.featuresInclude, t.featuresExclude, FeatureEntityDBAssignLastInsertID|FeatureDB) {
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
	if !g.hasFeature(t.featuresInclude, t.featuresExclude, FeatureEntityDBAssignLastInsertID|FeatureDB) {
		return
	}
	if !t.hasPKAutoInc() {
		return
	}

	mainGen.C(`AssignLastInsertID traverses through the slice and sets a decrementing new ID to each entity.`)
	mainGen.Pln(`func (cc *`, t.CollectionName(), `) AssignLastInsertID(id int64) {`)
	{
		mainGen.In()
		mainGen.Pln(`var j int64`)
		mainGen.Pln(`for i := len(cc.Data) - 1; i >= 0; i-- {`)
		{
			mainGen.In()
			mainGen.Pln(`cc.Data[i].AssignLastInsertID(id - j)`)
			mainGen.Pln(`j++`)
			mainGen.Out()
		}
		mainGen.Pln(`}`)
		mainGen.Out()
	}
	mainGen.Pln(`}`)
}

func (t *Table) fnCollectionUniqueGetters(mainGen *codegen.Go, g *Generator) {
	if !g.hasFeature(t.featuresInclude, t.featuresExclude, FeatureCollectionUniqueGetters|
		FeatureDB|FeatureEntityDBSelect|FeatureEntityDBDelete|
		FeatureEntityDBInsert|FeatureEntityDBUpdate|FeatureEntityDBUpsert) {
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
	if !g.hasFeature(t.featuresInclude, t.featuresExclude, FeatureCollectionDBMapColumns|FeatureDB) {
		return
	}

	mainGen.Pln(`func (cc *`, t.CollectionName(), `) scanColumns(cm *dml.ColumnMap, e *`, t.EntityName(), `, idx uint64) error {
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
							for i, e := range cc.Data {
								if err := cc.scanColumns(cm, e, uint64(i)); err != nil {
									return errors.WithStack(err)
								}
							}`)

		mainGen.Pln(`case dml.ColumnMapScan:
							if cm.Count == 0 {
								cc.Data = cc.Data[:0]
							}
							e := new(`, t.EntityName(), `)
							if err := cc.scanColumns(cm, e, cm.Count); err != nil {
								return errors.WithStack(err)
							}
							cc.Data = append(cc.Data, e)`)

		mainGen.Pln(`case dml.ColumnMapCollectionReadSet:
							for cm.Next() {
								switch c := cm.Column(); c {`)

		t.Table.Columns.UniqueColumns().Each(func(c *ddl.Column) {
			if !c.IsFloat() {
				mainGen.P(`case`, strconv.Quote(c.Field))
				for _, a := range c.Aliases {
					mainGen.P(`,`, strconv.Quote(a))
				}
				mainGen.Pln(`:`)
				mainGen.Pln(`cm = cm.`, g.goFuncNull(c)+`s(cc.`, strs.ToGoCamelCase(c.Field)+`s()...)`)
			}
		})
		mainGen.Pln(`default:
				return errors.NotFound.Newf("[`+t.Package+`]`, t.CollectionName(), `Column %q not found", c)
			}
		} // end for cm.Next

	default:
		return errors.NotSupported.Newf("[`+t.Package+`] Unknown Mode: %q", string(m))
	}
	return cm.Err()`)
	}
	mainGen.Pln(`}`) // end func MapColumns
}

func (tbl *Table) fnCreateDBMHandler(mainGen *codegen.Go, g *Generator) {
	var bufPKNameTypes strings.Builder
	var bufPKNames strings.Builder
	var bufPKNamesAsArgs strings.Builder
	i := 0

	tblCols := tbl.Table.Columns.PrimaryKeys()
	if tbl.Table.IsView() {
		tblCols = tbl.Table.Columns.ViewPrimaryKeys()
	}

	tblCols.Each(func(c *ddl.Column) {
		if i > 0 {
			bufPKNameTypes.WriteByte(',')
			bufPKNames.WriteByte(',')
			bufPKNamesAsArgs.WriteByte(',')
		}
		goNamedField := lcFirst(strs.ToGoCamelCase(c.Field))

		bufPKNameTypes.WriteString(goNamedField)
		bufPKNameTypes.WriteByte(' ')
		bufPKNameTypes.WriteString(g.goTypeNull(c))

		bufPKNames.WriteString(goNamedField)

		bufPKNamesAsArgs.WriteString("e.")
		bufPKNamesAsArgs.WriteString(strs.ToGoCamelCase(c.Field))
		i++
	})
	if i == 0 {
		mainGen.C("The table/view", tbl.EntityName(), "does not have a primary key. SKipping to generate DML functions based on the PK.")
		mainGen.Pln("\n")
		return
	}
	entityPTRName := codegen.SkipWS("*", tbl.EntityName())
	entityEventName := codegen.SkipWS(`event`, tbl.EntityName(), `Func`)
	tracingEnabled := tbl.hasFeature(g, FeatureEntityDBTracing)
	entityFuncName := codegen.SkipWS(tbl.EntityName(), "FindByPK")

	dmlEnabled := tbl.hasFeature(g, FeatureEntityDBSelect)
	mainGen.Pln(dmlEnabled,
		`func (dbm DBM) `, entityFuncName, `(ctx context.Context,  `, &bufPKNameTypes, `, opts ...dml.DBRFunc) (_ `, entityPTRName, `,err error) {`)

	mainGen.Pln(dmlEnabled && tracingEnabled, `	ctx, span := dbm.option.Trace.Start(ctx, `, codegen.SkipWS(`"`, entityFuncName, `"`), `)
		defer func(){ cstrace.Status(span, err); span.End(); }()`)

	mainGen.Pln(dmlEnabled, `var e `, tbl.EntityName(), `
	// put the IDs`, bufPKNames.String(), `into the context as value to search for a cache entry in the event function.
	if err = dbm.`, entityEventName, `(ctx, dml.EventFlagBeforeSelect, &e); err != nil {
		return nil, errors.WithStack(err)
	}
	if e.IsSet() {
		return &e, nil // returns data from cache
	}
	if _, err = dbm.CachedQuery(`, codegen.SkipWS(`"`, entityFuncName, `"`), `).ApplyCallBacks(opts...).Load(ctx, &e, `, &bufPKNames, `); err != nil {
		return nil, errors.WithStack(err)
	}
	if err = dbm.`, entityEventName, `(ctx, dml.EventFlagAfterSelect, &e); err != nil {
		return nil, errors.WithStack(err)
	}
	return &e, nil
}`)

	if tbl.Table.IsView() {
		// skip here the delete,insert,update and upsert functions.
		return
	}

	dmlEnabled = tbl.hasFeature(g, FeatureEntityDBDelete)
	entityFuncName = codegen.SkipWS(tbl.EntityName(), "DeleteByPK")
	mainGen.Pln(dmlEnabled, `func (dbm DBM) `, entityFuncName, `(ctx context.Context, `, &bufPKNameTypes, `, opts ...dml.DBRFunc) (err error) {`)

	mainGen.Pln(dmlEnabled && tracingEnabled, `	ctx, span := dbm.option.Trace.Start(ctx, `, codegen.SkipWS(`"`, entityFuncName, `"`), `)
			defer func(){ cstrace.Status(span, err); span.End(); }()`)

	// TODO think about the nil pointer in the events functions.
	mainGen.Pln(dmlEnabled, `if err = dbm.`, entityEventName, `(ctx, dml.EventFlagBeforeDelete, nil); err != nil {
			return errors.WithStack(err)
		}
		if _, err = dbm.CachedQuery(`, codegen.SkipWS(`"`, entityFuncName, `"`), `).ApplyCallBacks(opts...).ExecContext(ctx, `, &bufPKNames, `); err != nil {
			return errors.WithStack(err)
		}
		if err = dbm.`, entityEventName, `(ctx, dml.EventFlagAfterDelete, nil); err != nil {
			return errors.WithStack(err)
		}
		return nil
	}`)

	dmlEnabled = tbl.hasFeature(g, FeatureEntityDBUpdate)
	entityFuncName = codegen.SkipWS(tbl.EntityName(), "UpdateByPK")
	mainGen.Pln(dmlEnabled, `func (dbm DBM) `, entityFuncName, `(ctx context.Context, e `, entityPTRName, `, opts ...dml.DBRFunc) (err error) {`)

	mainGen.Pln(dmlEnabled && tracingEnabled, `	ctx, span := dbm.option.Trace.Start(ctx, `, codegen.SkipWS(`"`, entityFuncName, `"`), `);
			defer func(){ cstrace.Status(span, err); span.End(); }()`)

	mainGen.Pln(dmlEnabled, `if err = dbm.`, entityEventName, `(ctx, dml.EventFlagBeforeUpdate, e); err != nil {
			return errors.WithStack(err)
		}
		if _, err = dbm.CachedQuery(`, codegen.SkipWS(`"`, entityFuncName, `"`), `).ApplyCallBacks(opts...).ExecContext(ctx, &e); err != nil {
			return errors.WithStack(err)
		}
		if err = dbm.`, entityEventName, `(ctx, dml.EventFlagAfterUpdate, e); err != nil {
			return errors.WithStack(err)
		}
		return nil
	}`)

	dmlEnabled = tbl.hasFeature(g, FeatureEntityDBInsert)
	entityFuncName = codegen.SkipWS(tbl.EntityName(), "Insert")
	mainGen.Pln(dmlEnabled, `func (dbm DBM) `, codegen.SkipWS(tbl.EntityName(), "Insert"), `(ctx context.Context, e `, entityPTRName, `, opts ...dml.DBRFunc) (err error) {`)

	mainGen.Pln(dmlEnabled && tracingEnabled, `	ctx, span := dbm.option.Trace.Start(ctx, `, codegen.SkipWS(`"`, entityFuncName, `"`), `);
			defer func(){ cstrace.Status(span, err); span.End(); }()`)

	mainGen.Pln(dmlEnabled, `if err = dbm.`, entityEventName, `(ctx, dml.EventFlagBeforeInsert, e); err != nil {
			return errors.WithStack(err)
		}
		if _, err = dbm.CachedQuery(`, codegen.SkipWS(`"`, entityFuncName, `"`), `).ApplyCallBacks(opts...).ExecContext(ctx, &e); err != nil {
			return errors.WithStack(err)
		}
		if err = dbm.`, entityEventName, `(ctx, dml.EventFlagAfterInsert, e); err != nil {
			return errors.WithStack(err)
		}
		return nil
	}`)

	dmlEnabled = tbl.hasFeature(g, FeatureEntityDBUpsert)
	entityFuncName = codegen.SkipWS(tbl.EntityName(), "Upsert")
	mainGen.Pln(dmlEnabled, `func (dbm DBM) `, codegen.SkipWS(tbl.EntityName(), "Upsert"), `(ctx context.Context, e `, entityPTRName, `, opts ...dml.DBRFunc) (err error) {`)

	mainGen.Pln(dmlEnabled && tracingEnabled, `	ctx, span := dbm.option.Trace.Start(ctx, `, codegen.SkipWS(`"`, entityFuncName, `"`), `);
			defer func(){ cstrace.Status(span, err); span.End(); }()`)

	mainGen.Pln(dmlEnabled, `if err = dbm.`, entityEventName, `(ctx, dml.EventFlagBeforeUpsert, e); err != nil {
			return errors.WithStack(err)
		}
		if _, err = dbm.CachedQuery(`, codegen.SkipWS(`"`, entityFuncName, `"`), `).ApplyCallBacks(opts...).ExecContext(ctx, &e); err != nil {
			return errors.WithStack(err)
		}
		var e2 `, tbl.EntityName(), `
		if _, err = dbm.CachedQuery(`, codegen.SkipWS(`"`, tbl.EntityName(), "FindByPK", `"`), `).ApplyCallBacks(opts...).Load(ctx, &e2, `, &bufPKNamesAsArgs, `); err != nil {
			return errors.WithStack(err)
		}
		if err = dbm.`, entityEventName, `(ctx, dml.EventFlagAfterUpsert, &e2); err != nil {
			return errors.WithStack(err)
		}
		*e = e2
		return nil
	}`)
}

func (t *Table) optionsSQLBuildQueries(mainGen *codegen.Go, g *Generator) {
	tblPKLen := t.Table.Columns.PrimaryKeys().Len()
	if t.Table.IsView() {
		tblPKLen = t.Table.Columns.ViewPrimaryKeys().Len()
	}

	mainGen.Pln(tblPKLen > 0 && t.hasFeature(g, FeatureEntityDBSelect), `ddl.WithQueryDBR( `,
		codegen.SkipWS(`"`, t.EntityName(), `FindByPK"`),
		`, dbmo.InitSelectFn(tbls.MustTable(`, codegen.SkipWS(`TableName`, t.EntityName()), `).SelectByPK("*")).WithDBR().Interpolate()),`)

	if t.Table.IsView() {
		return
	}

	mainGen.Pln(t.hasFeature(g, FeatureEntityDBUpdate), `ddl.WithQueryDBR( `,
		codegen.SkipWS(`"`, t.EntityName(), `UpdateByPK"`),
		`, dbmo.InitUpdateFn(tbls.MustTable(`, codegen.SkipWS(`TableName`, t.EntityName()), `).UpdateByPK()).WithDBR()),`)
	mainGen.Pln(t.hasFeature(g, FeatureEntityDBDelete), `ddl.WithQueryDBR( `,
		codegen.SkipWS(`"`, t.EntityName(), `DeleteByPK"`),
		`, dbmo.InitDeleteFn(tbls.MustTable(`, codegen.SkipWS(`TableName`, t.EntityName()), `).DeleteByPK()).WithDBR().Interpolate()),`)
	mainGen.Pln(t.hasFeature(g, FeatureEntityDBInsert), `ddl.WithQueryDBR( `,
		codegen.SkipWS(`"`, t.EntityName(), `Insert"`),
		`, dbmo.InitInsertFn(tbls.MustTable(`, codegen.SkipWS(`TableName`, t.EntityName()), `).Insert()).WithDBR()),`)
	mainGen.Pln(t.hasFeature(g, FeatureEntityDBUpsert), `ddl.WithQueryDBR( `,
		codegen.SkipWS(`"`, t.EntityName(), `UpsertByPK"`),
		`, dbmo.InitInsertFn(tbls.MustTable(`, codegen.SkipWS(`TableName`, t.EntityName()), `).Insert()).OnDuplicateKey().WithDBR()),`)
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

	testGen.Pln(`entSELECT := tbl.SelectByPK("*")`)
	testGen.C(`WithDBR generates the cached SQL string with empty key "".`)
	testGen.Pln(`entSELECTStmtA := entSELECT.WithDBR().ExpandPlaceHolders()`)

	testGen.Pln(`entSELECT.WithCacheKey("select_10").Wheres.Reset()`)
	testGen.Pln(`_, _, err := entSELECT.Where(`)

	for _, c := range t.Table.Columns {
		if c.IsPK() && c.IsAutoIncrement() {
			testGen.Pln(`dml.Column(`, strconv.Quote(c.Field), `).LessOrEqual().Int(10),`)
		}
	}

	testGen.Pln(`).ToSQL() // ToSQL generates the new cached SQL string with key select_10`)
	testGen.Pln(`assert.NoError(t, err)`)
	testGen.Pln(`entCol := New`+t.CollectionName(), `()`)

	if t.HasAutoIncrement < 2 {
		testGen.C(`this table/view does not support auto_increment`)
		testGen.Pln(`rowCount, err := entSELECTStmtA.WithCacheKey("select_10").Load(ctx, entCol)`)
		testGen.Pln(`assert.NoError(t, err)`)
		testGen.Pln(`t.Logf("SELECT queries: %#v", entSELECT.CachedQueries())`)
		testGen.Pln(`t.Logf("Collection load rowCount: %d", rowCount)`)
	} else {
		testGen.Pln(`entINSERT := tbl.Insert().BuildValues()`)
		testGen.Pln(`entINSERTStmtA := entINSERT.PrepareWithDBR(ctx)`)

		testGen.Pln(`for i := 0; i < 9; i++ {`)
		{
			testGen.In()
			testGen.Pln(`entIn := new(`, strs.ToGoCamelCase(t.Table.Name), `)`)
			testGen.Pln(`if err := ps.FakeData(entIn); err != nil {`)
			{
				testGen.In()
				testGen.Pln(`t.Errorf("IDX[%d]: %+v", i, err)`)
				testGen.Pln(`return`)
				testGen.Out()
			}
			testGen.Pln(`}`)

			testGen.Pln(`lID := dmltest.CheckLastInsertID(t, "Error: TestNewTables.` + strs.ToGoCamelCase(t.Table.Name) + `_Entity")(entINSERTStmtA.ExecContext(ctx,dml.Qualify("", entIn)))`)
			testGen.Pln(`entINSERTStmtA.Reset()`)

			testGen.Pln(`entOut := new(`, strs.ToGoCamelCase(t.Table.Name), `)`)
			testGen.Pln(`rowCount, err := entSELECTStmtA.Load(ctx, entOut, lID)`)
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

		testGen.Pln(`rowCount, err := entSELECTStmtA.WithCacheKey("select_10").Load(ctx, entCol)`)
		testGen.Pln(`assert.NoError(t, err)`)
		testGen.Pln(`t.Logf("Collection load rowCount: %d", rowCount)`)

		testGen.Pln(`entINSERTStmtA = entINSERT.WithCacheKey("row_count_%d", len(entCol.Data)).Replace().SetRowCount(len(entCol.Data)).PrepareWithDBR(ctx)`)
		testGen.Pln(`lID := dmltest.CheckLastInsertID(t, "Error: `, t.CollectionName(), `")(entINSERTStmtA.ExecContext(ctx, dml.Qualify("", entCol)))`)
		testGen.Pln(`dmltest.Close(t, entINSERTStmtA)`)
		testGen.Pln(`t.Logf("Last insert ID into: %d", lID)`)
		testGen.Pln(`t.Logf("INSERT queries: %#v", entINSERT.CachedQueries())`)
		testGen.Pln(`t.Logf("SELECT queries: %#v", entSELECT.CachedQueries())`)
	}

	testGen.Pln(`})`)
}
