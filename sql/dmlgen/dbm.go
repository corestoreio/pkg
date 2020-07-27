package dmlgen

import (
	"fmt"

	"github.com/corestoreio/pkg/sql/ddl"
	"github.com/corestoreio/pkg/util/codegen"
	"github.com/corestoreio/pkg/util/strs"
)

func (g *Generator) fnCreateDBM(mainGen *codegen.Go, tbls tables) {
	if !tbls.hasFeature(g, FeatureDB|FeatureDBTracing|FeatureDBSelect|FeatureDBDelete|
		FeatureDBInsert|FeatureDBUpdate|FeatureDBUpsert|FeatureDBTableColumnNames) {
		return
	}

	// TODO add some feature switches to include/excluded parts.
	var tableNames []string
	var tableCreateStmt []string
	var tableConstants []string
	for _, tblname := range g.sortedTableNames() {
		constName := `TableName` + strs.ToGoCamelCase(tblname)
		tableConstants = append(tableConstants, fmt.Sprintf("%s = %q", constName, tblname))
		tableNames = append(tableNames, tblname)
		tableCreateStmt = append(tableCreateStmt, constName, `""`)
	}
	mainGen.C(`TableName constants define the names of all tables.`)
	mainGen.WriteConstants(tableConstants...)

	if tbls.hasFeature(g, FeatureDBTableColumnNames) {
		mainGen.C(`Columns struct provides for all tables the name of the columns. Allows type safety.`)
		mainGen.Pln(`var Columns = struct {`)
		{
			for _, tbl := range tbls {
				mainGen.Pln(tbl.EntityName(), `struct {`)
				{
					tbl.Table.Columns.Each(func(c *ddl.Column) {
						mainGen.Pln(strs.ToGoCamelCase(c.Field), `string`)
					})
				}
				mainGen.Pln(`}`)
			}
		}
		mainGen.Pln(`}{`)
		{
			for _, tbl := range tbls {
				mainGen.Pln(tbl.EntityName(), `: struct {`)
				{
					tbl.Table.Columns.Each(func(c *ddl.Column) {
						mainGen.Pln(strs.ToGoCamelCase(c.Field), `string`)
					})
				}
				mainGen.Pln(`}{`)
				{
					tbl.Table.Columns.Each(func(c *ddl.Column) {
						mainGen.Pln(strs.ToGoCamelCase(c.Field), `:`, fmt.Sprintf("%q", c.Field), ",")
					})
				}
				mainGen.Pln(`},`)
			}
		}
		mainGen.Pln(`}`) // end main struct
	}

	mainGen.Pln(`var dbmEmptyOpts = []dml.DBRFunc{func(dbr *dml.DBR) {
			// do nothing because Clone gets called automatically
		}}
		func dbmNoopResultCheckFn(_ sql.Result, err error) error { return err }
`)

	// <event functions>
	mainGen.C(`Event functions are getting dispatched during before or after handling a collection or an entity.
Context is always non-nil but either collection or entity pointer will be set.`)
	mainGen.Pln(`type (`)
	for _, tbl := range tbls {
		mainGen.Pln(`Event` + tbl.EntityName() + `Fn func(context.Context, *` + tbl.CollectionName() + `, *` + tbl.EntityName() + `) error`)
	}
	mainGen.Pln(`)`)
	// </event functions>

	// <DBM option struct>
	mainGen.C(`DBMOption provides various options to the DBM object.`)
	mainGen.Pln(`type DBMOption struct {`)
	{
		mainGen.Pln(tbls.hasFeature(g, FeatureDBTracing), `Trace                trace.Tracer`)
		mainGen.Pln(`TableOptions         []ddl.TableOption // gets applied at the beginning`)
		mainGen.Pln(`TableOptionsAfter    []ddl.TableOption // gets applied at the end`)
		mainGen.Pln(tbls.hasFeature(g, FeatureDBSelect), `InitSelectFn         func(*dml.Select) *dml.Select`)
		mainGen.Pln(tbls.hasFeature(g, FeatureDBUpdate), `InitUpdateFn         func(*dml.Update) *dml.Update`)
		mainGen.Pln(tbls.hasFeature(g, FeatureDBDelete), `InitDeleteFn         func(*dml.Delete) *dml.Delete`)
		mainGen.Pln(tbls.hasFeature(g, FeatureDBInsert|FeatureDBUpsert), `InitInsertFn         func(*dml.Insert) *dml.Insert`)
		for _, tbl := range tbls {
			mainGen.Pln(`event` + tbl.EntityName() + `Func [dml.EventFlagMax][]Event` + tbl.EntityName() + `Fn`)
		}
	}
	mainGen.Pln(`}`)
	// </DBM option struct>

	// <event adder>
	for _, tbl := range tbls {
		mainGen.C(codegen.SkipWS(`AddEvent`, tbl.EntityName()), `adds a specific defined event call back to the DBM.
It panics if the event argument is larger than dml.EventFlagMax.`)
		mainGen.Pln(`func (o *DBMOption) `, codegen.SkipWS(`AddEvent`, tbl.EntityName(), `(`), `event dml.EventFlag,`,
			`fn Event`+tbl.EntityName()+`Fn) *DBMOption {`)
		{
			mainGen.In()
			mainGen.Pln(`o.event` + tbl.EntityName() + `Func[event] = append(o.event` + tbl.EntityName() + `Func[event], fn)`)
			mainGen.Pln(`return o`)
			mainGen.Out()
		}
		mainGen.Pln(`}`)
	}
	// </event adder>

	mainGen.C(`DBM defines the DataBaseManagement object for the tables `, tableNames)
	mainGen.Pln(`type DBM struct { *ddl.Tables; option DBMOption }`)

	// <event dispatcher>
	for _, tbl := range tbls {
		mainGen.Pln(codegen.SkipWS(`func (dbm DBM) event`, tbl.EntityName(), `Func(ctx context.Context, ef dml.EventFlag, skipEvents bool, ec `, codegen.SkipWS(`*`, tbl.CollectionName()), `, e `, codegen.SkipWS(`*`, tbl.EntityName()), `) error`), ` {`)
		{
			mainGen.In()
			mainGen.Pln(`if len(dbm.option.`, codegen.SkipWS(`event`, tbl.EntityName(), `Func`), `[ef]) == 0 || skipEvents {`)
			mainGen.In()
			{
				mainGen.Pln(`return nil`)
			}
			mainGen.Out()
			mainGen.Pln(`}`)

			mainGen.Pln(`for _, fn := range dbm.option.`, codegen.SkipWS(`event`, tbl.EntityName(), `Func`), `[ef] {`)
			{
				mainGen.In()
				mainGen.Pln(`if err := fn(ctx, ec, e); err != nil {
				return errors.WithStack(err)
			}`)
				mainGen.Out()
			}
			mainGen.Pln(`}`)
			mainGen.Pln(`return nil`)
			mainGen.Out()
		}
		mainGen.Pln(`}`)
	} // </event dispatcher>

	mainGen.C(`NewDBManager returns a goified version of the MySQL/MariaDB table schema for the tables: `, tableNames, `Auto generated by dmlgen.`)
	mainGen.Pln(`func NewDBManager(ctx context.Context, dbmo *DBMOption) (*DBM, error) {`)
	{
		mainGen.In()
		mainGen.Pln(`tbls, err := ddl.NewTables(append([]ddl.TableOption{ddl.WithCreateTable(ctx, `, tableCreateStmt, `)},dbmo.TableOptions...)...)`)
		mainGen.Pln(`if err != nil { return nil, errors.WithStack(err); }`)

		mainGen.Pln(tbls.hasFeature(g, FeatureDBSelect),
			`	if dbmo.InitSelectFn == nil { dbmo.InitSelectFn = func(s *dml.Select) *dml.Select { return s; }; } `)
		mainGen.Pln(tbls.hasFeature(g, FeatureDBUpdate),
			`	if dbmo.InitUpdateFn == nil { dbmo.InitUpdateFn = func(s *dml.Update) *dml.Update { return s; }; } `)
		mainGen.Pln(tbls.hasFeature(g, FeatureDBDelete),
			`	if dbmo.InitDeleteFn == nil { dbmo.InitDeleteFn = func(s *dml.Delete) *dml.Delete { return s; }; } `)
		mainGen.Pln(tbls.hasFeature(g, FeatureDBInsert|FeatureDBUpsert),
			`	if dbmo.InitInsertFn == nil { dbmo.InitInsertFn = func(s *dml.Insert) *dml.Insert { return s; }; } `)

		{
			mainGen.Pln(`err = tbls.Options(`)
			mainGen.Pln(`ddl.WithQueryDBR(map[string]dml.QueryBuilder{`)
			for _, tbl := range tbls {
				tbl.fnDBMOptionsSQLBuildQueries(mainGen, g)
			}
			mainGen.Pln("}),")
			mainGen.Pln(`)`) // end options
			mainGen.Pln(`if err != nil { return nil, err }`)
			mainGen.Pln(`if err := tbls.Options(dbmo.TableOptionsAfter...); err != nil { return nil, err }`)
		}
		mainGen.Out()
	}

	mainGen.Pln(tbls.hasFeature(g, FeatureDBTracing), `	if dbmo.Trace == nil { dbmo.Trace = trace.NoopTracer{}; }`)
	mainGen.Pln(`return &DBM{	Tables: tbls, option: *dbmo, }, nil }`)
}
