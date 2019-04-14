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
	"strconv"
	"unicode"

	"github.com/corestoreio/pkg/sql/ddl"
	"github.com/corestoreio/pkg/util/codegen"
	"github.com/corestoreio/pkg/util/strs"
)

// FeatureToggle allows certain generated code blocks to be switched off or on.
type FeatureToggle uint64

// List of available features
const (
	FeatureDB FeatureToggle = 1 << iota
	FeatureCollectionStruct
	FeatureEntityStruct
	FeatureEntityGetSetPrivateFields
	FeatureEntityEmpty
	FeatureEntityCopy
	FeatureEntityDBMapColumns
	FeatureEntityWriteTo
	FeatureEntityDBAssignLastInsertID
	FeatureEntityRelationships
	FeatureCollectionUniqueGetters
	FeatureCollectionUniquifiedGetters
	FeatureCollectionFilter
	FeatureCollectionEach
	FeatureCollectionCut
	FeatureCollectionSwap
	FeatureCollectionDelete
	FeatureCollectionInsert
	FeatureCollectionAppend
	FeatureCollectionBinaryMarshaler
	FeatureCollectionDBMapColumns
)

// table writes one database table into Go source code.
type Table struct {
	Package              string      // Name of the package
	TableName            string      // Name of the table
	Comment              string      // Comment above the struct type declaration
	Columns              ddl.Columns // all columns of a table
	HasAutoIncrement     uint8       // 0=nil,1=false (has NO auto increment),2=true has auto increment
	HasJSONMarshaler     bool
	HasEasyJSONMarshaler bool
	HasBinaryMarshaler   bool
	HasSerializer        bool // writes the .proto file if true

	// PrivateFields key=snake case name of the DB column, value=true, the field must be private
	privateFields   map[string]bool
	featuresInclude FeatureToggle
	featuresExclude FeatureToggle
	fieldMapFn      func(dbIdentifier string) (newName string)
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
	sr := []rune(su)
	sr[0] = unicode.ToLower(sr[0])
	return string(sr)
}

func (t *Table) CollectionName() string {
	return strs.ToGoCamelCase(t.TableName) + "Collection"
}

func (t *Table) EntityName() string {
	return strs.ToGoCamelCase(t.TableName)
}

func (t *Table) collectionStruct(mainGen *codegen.Go, ts *Generator) {
	if !ts.hasFeature(t.featuresInclude, t.featuresExclude, FeatureCollectionStruct) {
		return
	}

	mainGen.C(t.CollectionName(), `represents a collection type for DB table`, t.TableName)
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
		if ts.hasFeature(t.featuresInclude, t.featuresExclude, FeatureEntityDBMapColumns|FeatureDB) {
			mainGen.Pln(`BeforeMapColumns	func(uint64, *`, t.EntityName(), `) error`, codegen.EncloseBT(`json:"-"`))
			mainGen.Pln(`AfterMapColumns 	func(uint64, *`, t.EntityName(), `) error `, codegen.EncloseBT(`json:"-"`))
		}
		if fn, ok := ts.customCode["type_"+t.CollectionName()]; ok {
			fn(ts, t, mainGen)
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

func (t *Table) entityStruct(mainGen *codegen.Go, ts *Generator) {
	if !ts.hasFeature(t.featuresInclude, t.featuresExclude, FeatureEntityStruct) {
		return
	}

	fieldMapFn := ts.defaultTableConfig.FieldMapFn
	if fieldMapFn == nil {
		fieldMapFn = t.fieldMapFn
	}
	if fieldMapFn == nil {
		fieldMapFn = defaultFieldMapFn
	}

	mainGen.C(t.EntityName(), `represents a single row for DB table`, t.TableName+`. Auto generated.`)
	if t.Comment != "" {
		mainGen.C(t.Comment)
	}
	if t.HasEasyJSONMarshaler {
		mainGen.Pln(`//easyjson:json`)
	}

	// Generate table structs
	mainGen.Pln(`type `, t.EntityName(), ` struct {`)
	{
		if fn, ok := ts.customCode["type_"+t.EntityName()]; ok {
			fn(ts, t, mainGen)
		} else {
			mainGen.In()
			for _, c := range t.Columns {
				structTag := ""
				if c.StructTag != "" {
					structTag += "`" + c.StructTag + "`"
				}
				mainGen.Pln(t.GoCamelMaybePrivate(c.Field), ts.goTypeNull(c), structTag, c.GoComment())
			}

			if ts.hasFeature(t.featuresInclude, t.featuresExclude, FeatureEntityRelationships) {
				const debug = true
				if kcuc, ok := ts.kcu[t.TableName]; ok { // kcu = keyColumnUsage && kcuc = keyColumnUsageCollection
					for _, kcuce := range kcuc.Data {
						if !kcuce.ReferencedTableName.Valid {
							continue
						}

						// case ONE-TO-MANY
						isOneToMany := ts.krs.IsOneToMany(kcuce.TableName, kcuce.ColumnName, kcuce.ReferencedTableName.String, kcuce.ReferencedColumnName.String)
						isRelationAllowed := !ts.skipRelationship(kcuce.TableName, kcuce.ColumnName, kcuce.ReferencedTableName.String, kcuce.ReferencedColumnName.String)
						hasTable := ts.Tables[kcuce.ReferencedTableName.String] != nil
						if debug {
							println("A1: isOneToMany", isOneToMany, "\tisRelationAllowed", isRelationAllowed, "\thasTable", hasTable, "\t",
								t.TableName, "\t",
								kcuce.TableName+"."+kcuce.ColumnName, "=>", kcuce.ReferencedTableName.String+"."+kcuce.ReferencedColumnName.String)
						}
						if isOneToMany && hasTable && isRelationAllowed {
							mainGen.Pln(fieldMapFn(kcuce.ReferencedTableName.String), " *", strs.ToGoCamelCase(kcuce.ReferencedTableName.String)+"Collection",
								"// 1:M", kcuce.TableName+"."+kcuce.ColumnName, "=>", kcuce.ReferencedTableName.String+"."+kcuce.ReferencedColumnName.String)
						}

						// case ONE-TO-ONE
						isOneToOne := ts.krs.IsOneToOne(kcuce.TableName, kcuce.ColumnName, kcuce.ReferencedTableName.String, kcuce.ReferencedColumnName.String)
						if debug {
							println("B1: IsOneToOne", isOneToOne, "\tisRelationAllowed", isRelationAllowed, "\thasTable", hasTable, "\t",
								t.TableName, "\t",
								kcuce.TableName+"."+kcuce.ColumnName, "=>", kcuce.ReferencedTableName.String+"."+kcuce.ReferencedColumnName.String)
						}
						if isOneToOne && hasTable && isRelationAllowed {
							mainGen.Pln(fieldMapFn(kcuce.ReferencedTableName.String), " *", strs.ToGoCamelCase(kcuce.ReferencedTableName.String),
								"// 1:1", kcuce.TableName+"."+kcuce.ColumnName, "=>", kcuce.ReferencedTableName.String+"."+kcuce.ReferencedColumnName.String)
						}
					}
				}

				if kcuc, ok := ts.kcuRev[t.TableName]; ok { // kcu = keyColumnUsage && kcuc = keyColumnUsageCollection
					for _, kcuce := range kcuc.Data {
						if !kcuce.ReferencedTableName.Valid {
							continue
						}

						// case ONE-TO-MANY
						isOneToMany := ts.krs.IsOneToMany(kcuce.TableName, kcuce.ColumnName, kcuce.ReferencedTableName.String, kcuce.ReferencedColumnName.String)
						isRelationAllowed := !ts.skipRelationship(kcuce.TableName, kcuce.ColumnName, kcuce.ReferencedTableName.String, kcuce.ReferencedColumnName.String)
						hasTable := ts.Tables[kcuce.ReferencedTableName.String] != nil
						if debug {
							println("A2: isOneToMany", isOneToMany, "\tisRelationAllowed", isRelationAllowed, "\thasTable", hasTable, "\t",
								t.TableName, "\t",
								kcuce.TableName+"."+kcuce.ColumnName, "=>", kcuce.ReferencedTableName.String+"."+kcuce.ReferencedColumnName.String)
						}
						if isOneToMany && hasTable && isRelationAllowed {
							mainGen.Pln(fieldMapFn(kcuce.ReferencedTableName.String), " *", strs.ToGoCamelCase(kcuce.ReferencedTableName.String)+"Collection",
								"// Reversed 1:M", kcuce.TableName+"."+kcuce.ColumnName, "=>", kcuce.ReferencedTableName.String+"."+kcuce.ReferencedColumnName.String)
						}

						// case ONE-TO-ONE
						isOneToOne := ts.krs.IsOneToOne(kcuce.TableName, kcuce.ColumnName, kcuce.ReferencedTableName.String, kcuce.ReferencedColumnName.String)
						if debug {
							println("B2: IsOneToOne", isOneToOne, "\tisRelationAllowed", isRelationAllowed, "\thasTable", hasTable, "\t",
								t.TableName, "\t",
								kcuce.TableName+"."+kcuce.ColumnName, "=>", kcuce.ReferencedTableName.String+"."+kcuce.ReferencedColumnName.String)
						}
						if isOneToOne && hasTable && isRelationAllowed {
							mainGen.Pln(fieldMapFn(kcuce.ReferencedTableName.String), " *", strs.ToGoCamelCase(kcuce.ReferencedTableName.String),
								"// Reversed 1:1", kcuce.TableName+"."+kcuce.ColumnName, "=>", kcuce.ReferencedTableName.String+"."+kcuce.ReferencedColumnName.String)
						}
					}
				}
			}
			mainGen.Out()
		}
	}
	mainGen.Pln(`}`)
}

func (t *Table) fnEntityGetSetPrivateFields(mainGen *codegen.Go, ts *Generator) {
	if !ts.hasFeature(t.featuresInclude, t.featuresExclude, FeatureEntityGetSetPrivateFields) {
		return
	}
	// Generates the Getter/Setter for private fields
	for _, c := range t.Columns {
		if !t.IsFieldPrivate(c.Field) {
			continue
		}
		mainGen.C(`Set`, strs.ToGoCamelCase(c.Field), ` sets the data for a private and security sensitive field.`)
		mainGen.Pln(`func (e *`, t.EntityName(), `) Set`+strs.ToGoCamelCase(c.Field), `(d `, ts.goTypeNull(c), `) *`, t.EntityName(), ` {`)
		{
			mainGen.In()
			mainGen.Pln(`e.`, t.GoCamelMaybePrivate(c.Field), ` = d`)
			mainGen.Pln(`return e`)
			mainGen.Out()
		}
		mainGen.Pln(`}`)

		mainGen.C(`Get`, strs.ToGoCamelCase(c.Field), ` returns the data from a private and security sensitive field.`)
		mainGen.Pln(`func (e *`, t.EntityName(), `) Get`+strs.ToGoCamelCase(c.Field), `() `, ts.goTypeNull(c), `{`)
		{
			mainGen.In()
			mainGen.Pln(`return e.`, t.GoCamelMaybePrivate(c.Field))
			mainGen.Out()
		}
		mainGen.Pln(`}`)

	}
}

func (t *Table) fnEntityEmpty(mainGen *codegen.Go, ts *Generator) {
	if !ts.hasFeature(t.featuresInclude, t.featuresExclude, FeatureEntityEmpty) {
		return
	}
	mainGen.Pln(`// Empty empties all the fields of the current object. Also known as Reset.`)
	// no idea if pointer dereferencing is bad ...
	mainGen.Pln(`func (e *`, t.EntityName(), `) Empty() *`, t.EntityName(), ` { *e = `, t.EntityName(), `{}; return e }`)
}

func (t *Table) fnEntityCopy(mainGen *codegen.Go, ts *Generator) {
	if !ts.hasFeature(t.featuresInclude, t.featuresExclude, FeatureEntityCopy) {
		return
	}
	mainGen.Pln(`// Copy copies the struct and returns a new pointer`)
	mainGen.Pln(`func (e *`, t.EntityName(), `) Copy() *`, t.EntityName(), ` { 
		e2 := new(`, t.EntityName(), `)
		*e2 = *e // for now a shallow copy
		return e2 
}`)
}

func (t *Table) fnEntityWriteTo(mainGen *codegen.Go, ts *Generator) {
	if !ts.hasFeature(t.featuresInclude, t.featuresExclude, FeatureEntityWriteTo) {
		return
	}
	mainGen.C(`WriteTo implements io.WriterTo and writes the field names and their values to w.`,
		`This is especially useful for debugging or or generating a hash of the struct.`)

	mainGen.Pln(`func (e *`, t.EntityName(), `) WriteTo(w io.Writer) (n int64, err error) {
	// for now this printing is good enough. If you need better swap out with your code.`)

	if fn, ok := ts.customCode["func_"+t.EntityName()+"_WriteTo"]; ok {
		fn(ts, t, mainGen)
	} else {
		mainGen.Pln(`n2, err := fmt.Fprint(w,`)
		mainGen.In()
		t.Columns.Each(func(c *ddl.Column) {
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

func (t *Table) fnCollectionWriteTo(mainGen *codegen.Go, ts *Generator) {
	if !ts.hasFeature(t.featuresInclude, t.featuresExclude, FeatureEntityWriteTo) {
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

func (t *Table) fnEntityDBMapColumns(mainGen *codegen.Go, ts *Generator) {
	if !ts.hasFeature(t.featuresInclude, t.featuresExclude, FeatureEntityDBMapColumns|FeatureDB) {
		return
	}
	mainGen.C(`MapColumns implements interface ColumnMapper only partially. Auto generated.`)
	mainGen.Pln(`func (e *`, t.EntityName(), `) MapColumns(cm *dml.ColumnMap) error {`)
	{
		if fn, ok := ts.customCode["func_"+t.EntityName()+"_MapColumns"]; ok {
			fn(ts, t, mainGen)
		}

		mainGen.In()
		mainGen.Pln(`if cm.Mode() == dml.ColumnMapEntityReadAll {`)
		{
			mainGen.In()
			mainGen.P(`return cm`)
			t.Columns.Each(func(c *ddl.Column) {
				mainGen.P(`.`, ts.goFuncNull(c), `(&e.`, t.GoCamelMaybePrivate(c.Field), `)`)
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
				t.Columns.Each(func(c *ddl.Column) {
					mainGen.P(`case`, strconv.Quote(c.Field))
					for _, a := range c.Aliases {
						mainGen.P(`,`, strconv.Quote(a))
					}
					mainGen.Pln(`:`)
					mainGen.Pln(`cm.`, ts.goFuncNull(c), `(&e.`, t.GoCamelMaybePrivate(c.Field), `)`)
				})
				mainGen.Pln(`default:`)
				mainGen.Pln(`return errors.NotFound.Newf("[`+ts.Package+`]`, t.EntityName(), `Column %q not found", c)`)
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
	t.Columns.Each(func(c *ddl.Column) {
		if c.IsPK() && c.IsAutoIncrement() {
			hasPKAutoInc = true
		}
		if hasPKAutoInc {
			return
		}
	})
	return hasPKAutoInc
}

func (t *Table) fnEntityDBAssignLastInsertID(mainGen *codegen.Go, ts *Generator) {
	if !ts.hasFeature(t.featuresInclude, t.featuresExclude, FeatureEntityDBAssignLastInsertID|FeatureDB) {
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
		t.Columns.Each(func(c *ddl.Column) {
			if c.IsPK() && c.IsAutoIncrement() {
				mainGen.Pln(`e.`, t.GoCamelMaybePrivate(c.Field), ` = `, ts.goType(c), `(id)`)
			}
		})
		mainGen.Out()
	}
	mainGen.Pln(`}`)
}

func (t *Table) fnCollectionDBAssignLastInsertID(mainGen *codegen.Go, ts *Generator) {
	if !ts.hasFeature(t.featuresInclude, t.featuresExclude, FeatureEntityDBAssignLastInsertID|FeatureDB) {
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

func (t *Table) fnCollectionUniqueGetters(mainGen *codegen.Go, ts *Generator) {
	if !ts.hasFeature(t.featuresInclude, t.featuresExclude, FeatureCollectionUniqueGetters) {
		return
	}

	// Generates functions to return all data as a slice from unique/primary
	// columns.
	for _, c := range t.Columns.UniqueColumns() {
		gtn := ts.goTypeNull(c)
		goCamel := strs.ToGoCamelCase(c.Field)
		mainGen.C(goCamel + `s returns a slice with the data or appends it to a slice.`)
		mainGen.C(`Auto generated.`)
		mainGen.Pln(`func (cc *`, t.CollectionName(), `) `, goCamel+`s(ret ...`+gtn, `) []`+gtn, ` {`)
		{
			mainGen.In()
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

func (t *Table) fnCollectionUniquifiedGetters(mainGen *codegen.Go, ts *Generator) {
	if !ts.hasFeature(t.featuresInclude, t.featuresExclude, FeatureCollectionUniquifiedGetters) {
		return
	}
	// Generates functions to return data with removed duplicates from any
	// column which has set the flag Uniquified.
	for _, c := range t.Columns.UniquifiedColumns() {
		goType := ts.goType(c)
		goCamel := strs.ToGoCamelCase(c.Field)

		mainGen.C(goCamel+`s belongs to the column`, strconv.Quote(c.Field), `and returns a slice or appends to a slice only`,
			`unique values of that column. The values will be filtered internally in a Go map. No DB query gets`,
			`executed. Auto generated.`)
		mainGen.Pln(`func (cc *`, t.CollectionName(), `) Unique`+goCamel+`s(ret ...`, goType, `) []`, goType, ` {`)
		{
			mainGen.In()
			mainGen.Pln(`if ret == nil {
					ret = make([]`, goType, `, 0, len(cc.Data))
				}`)

			// TODO: a reusable map and use different algorithms depending on
			// the size of the cc.Data slice. Sometimes a for/for loop runs
			// faster than a map.
			goPrimNull := ts.toGoPrimitiveFromNull(c)
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

func (t *Table) fnCollectionFilter(mainGen *codegen.Go, ts *Generator) {
	if !ts.hasFeature(t.featuresInclude, t.featuresExclude, FeatureCollectionFilter) {
		return
	}
	mainGen.C(`Filter filters the current slice by predicate f without memory allocation. Auto generated via dmlgen.`)
	mainGen.Pln(`func (cc *`, t.CollectionName(), `) Filter(f func(*`, t.EntityName(), `) bool) *`, t.CollectionName(), ` {`)
	{
		mainGen.In()
		mainGen.Pln(`b,i := cc.Data[:0],0`)
		mainGen.Pln(`for _, e := range cc.Data {`)
		{
			mainGen.In()
			mainGen.Pln(`if f(e) {`)
			{
				mainGen.Pln(`b = append(b, e)`)
				mainGen.Pln(`cc.Data[i] = nil // this avoids the memory leak`)
			}
			mainGen.Pln(`}`) // endif
			mainGen.Pln(`i++`)
		}
		mainGen.Out()
		mainGen.Pln(`}`) // for loop
		mainGen.Pln(`cc.Data = b`)
		mainGen.Pln(`return cc`)
		mainGen.Out()
	}
	mainGen.Pln(`}`) // function
}

func (t *Table) fnCollectionEach(mainGen *codegen.Go, ts *Generator) {
	if !ts.hasFeature(t.featuresInclude, t.featuresExclude, FeatureCollectionEach) {
		return
	}
	mainGen.C(`Each will run function f on all items in []*`, t.EntityName(), `. Auto generated via dmlgen.`)
	mainGen.Pln(`func (cc *`, t.CollectionName(), `) Each(f func(*`, t.EntityName(), `)) *`, t.CollectionName(), ` {`)
	{
		mainGen.Pln(`for i := range cc.Data {`)
		{
			mainGen.Pln(`f(cc.Data[i])`)
		}
		mainGen.Pln(`}`)
		mainGen.Pln(`return cc`)
	}
	mainGen.Pln(`}`)
}

func (t *Table) fnCollectionCut(mainGen *codegen.Go, ts *Generator) {
	if !ts.hasFeature(t.featuresInclude, t.featuresExclude, FeatureCollectionCut) {
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

func (t *Table) fnCollectionSwap(mainGen *codegen.Go, ts *Generator) {
	if !ts.hasFeature(t.featuresInclude, t.featuresExclude, FeatureCollectionSwap) {
		return
	}
	mainGen.C(`Swap will satisfy the sort.Interface. Auto generated via dmlgen.`)
	mainGen.Pln(`func (cc *`, t.CollectionName(), `) Swap(i, j int) { cc.Data[i], cc.Data[j] = cc.Data[j], cc.Data[i] }`)
}

func (t *Table) fnCollectionDelete(mainGen *codegen.Go, ts *Generator) {
	if !ts.hasFeature(t.featuresInclude, t.featuresExclude, FeatureCollectionDelete) {
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

func (t *Table) fnCollectionInsert(mainGen *codegen.Go, ts *Generator) {
	if !ts.hasFeature(t.featuresInclude, t.featuresExclude, FeatureCollectionInsert) {
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

func (t *Table) fnCollectionAppend(mainGen *codegen.Go, ts *Generator) {
	if !ts.hasFeature(t.featuresInclude, t.featuresExclude, FeatureCollectionAppend) {
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

func (t *Table) fnCollectionBinaryMarshaler(mainGen *codegen.Go, ts *Generator) {
	if !ts.hasFeature(t.featuresInclude, t.featuresExclude, FeatureCollectionBinaryMarshaler) {
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

func (t *Table) fnCollectionDBMapColumns(mainGen *codegen.Go, ts *Generator) {
	if !ts.hasFeature(t.featuresInclude, t.featuresExclude, FeatureCollectionDBMapColumns|FeatureDB) {
		return
	}

	mainGen.Pln(`func (cc *`, t.CollectionName(), `) scanColumns(cm *dml.ColumnMap,e *`, t.EntityName(), `, idx uint64) error {
			if cc.BeforeMapColumns != nil {
				if err := cc.BeforeMapColumns(idx, e); err != nil {
					return errors.WithStack(err)
				}
			}
			if err := e.MapColumns(cm); err != nil {
				return errors.WithStack(err)
			}
			if cc.AfterMapColumns != nil {
				if err := cc.AfterMapColumns(idx, e); err != nil {
					return errors.WithStack(err)
				}
			}
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

		t.Columns.UniqueColumns().Each(func(c *ddl.Column) {
			if !c.IsFloat() {
				mainGen.P(`case`, strconv.Quote(c.Field))
				for _, a := range c.Aliases {
					mainGen.P(`,`, strconv.Quote(a))
				}
				mainGen.Pln(`:`)
				mainGen.Pln(`cm = cm.`, ts.goFuncNull(c)+`s(cc.`, strs.ToGoCamelCase(c.Field)+`s()...)`)
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

func (t *Table) generateTestOther(testGen *codegen.Go, ts *Generator) {

	if ts.hasFeature(t.featuresInclude, t.featuresExclude, FeatureEntityEmpty) {
		testGen.Pln(`t.Run("` + strs.ToGoCamelCase(t.TableName) + `_Empty", func(t *testing.T) {`)
		{
			testGen.Pln(`e:= new(`, t.EntityName(), `)`)
			testGen.Pln(`assert.NoError(t, ps.FakeData(e))`)
			testGen.Pln(`e.Empty()`)
			testGen.Pln(`assert.Exactly(t, *e, `, t.EntityName(), `{})`)
		}
		testGen.Pln(`})`) // end t.Run
	}
	if ts.hasFeature(t.featuresInclude, t.featuresExclude, FeatureEntityCopy) {
		testGen.Pln(`t.Run("` + strs.ToGoCamelCase(t.TableName) + `_Copy", func(t *testing.T) {`)
		{
			testGen.Pln(`e:= new(`, t.EntityName(), `)`)
			testGen.Pln(`assert.NoError(t, ps.FakeData(e))`)
			testGen.Pln(`e2 := e.Copy()`)
			testGen.Pln(`assert.Exactly(t, e, e2)`)
			testGen.Pln(`assert.NoError(t, ps.FakeData(e))`)
			testGen.Pln(`assert.NotEqual(t, e, e2)`)
		}
		testGen.Pln(`})`) // end t.Run
	}
	// more feature tests to follow
}

func (t *Table) generateTestDB(testGen *codegen.Go) {

	testGen.Pln(`t.Run("` + strs.ToGoCamelCase(t.TableName) + `_Entity", func(t *testing.T) {`)
	testGen.Pln(`tbl := tbls.MustTable(TableName`+strs.ToGoCamelCase(t.TableName), `)`)

	testGen.Pln(`entSELECT := tbl.SelectByPK("*")`)
	testGen.C(`WithArgs generates the cached SQL string with empty key "".`)
	testGen.Pln(`entSELECTStmtA := entSELECT.WithArgs().ExpandPlaceHolders()`)

	testGen.Pln(`entSELECT.WithCacheKey("select_10").Wheres.Reset()`)
	testGen.Pln(`_, _, err := entSELECT.Where(`)

	for _, c := range t.Columns {
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
		testGen.Pln(`entINSERTStmtA := entINSERT.PrepareWithArgs(ctx)`)

		testGen.Pln(`for i := 0; i < 9; i++ {`)
		{
			testGen.In()
			testGen.Pln(`entIn := new(`, strs.ToGoCamelCase(t.TableName), `)`)
			testGen.Pln(`if err := ps.FakeData(entIn); err != nil {`)
			{
				testGen.In()
				testGen.Pln(`t.Errorf("IDX[%d]: %+v", i, err)`)
				testGen.Pln(`return`)
				testGen.Out()
			}
			testGen.Pln(`}`)

			testGen.Pln(`lID := dmltest.CheckLastInsertID(t, "Error: TestNewTables.` + strs.ToGoCamelCase(t.TableName) + `_Entity")(entINSERTStmtA.Record("", entIn).ExecContext(ctx))`)
			testGen.Pln(`entINSERTStmtA.Reset()`)

			testGen.Pln(`entOut := new(`, strs.ToGoCamelCase(t.TableName), `)`)
			testGen.Pln(`rowCount, err := entSELECTStmtA.Int64s(lID).Load(ctx, entOut)`)
			testGen.Pln(`assert.NoError(t, err)`)
			testGen.Pln(`assert.Exactly(t, uint64(1), rowCount, "IDX%d: RowCount did not match", i)`)

			for _, c := range t.Columns {
				fn := t.GoCamelMaybePrivate(c.Field)
				switch {
				case c.IsString():
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

		testGen.Pln(`entINSERTStmtA = entINSERT.WithCacheKey("row_count_%d", len(entCol.Data)).Replace().SetRowCount(len(entCol.Data)).PrepareWithArgs(ctx)`)
		testGen.Pln(`lID := dmltest.CheckLastInsertID(t, "Error: `, t.CollectionName(), `")(entINSERTStmtA.Record("", entCol).ExecContext(ctx))`)
		testGen.Pln(`dmltest.Close(t, entINSERTStmtA)`)
		testGen.Pln(`t.Logf("Last insert ID into: %d", lID)`)
		testGen.Pln(`t.Logf("INSERT queries: %#v", entINSERT.CachedQueries())`)
		testGen.Pln(`t.Logf("SELECT queries: %#v", entSELECT.CachedQueries())`)
	}

	testGen.Pln(`})`)
}
