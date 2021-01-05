package dmlgen

import (
	"fmt"
	"strconv"

	"github.com/corestoreio/pkg/util/codegen"
)

func (g *Generator) fnTestMainOther(testGen *codegen.Go, tbls tables) {
	// Test Header
	lenBefore := testGen.Len()
	var codeWritten int

	testGen.Pln(`func TestNewDBManagerNonDB_` + tbls.nameID() + `(t *testing.T) {`)
	{
		testGen.Pln(`ps := pseudo.MustNewService(0, &pseudo.Options{Lang: "de",MaxFloatDecimals:6})`)
		// If some features haven't been enabled, then there are no tests so
		// assign ps to underscore to avoid the unused variable error.
		// Alternatively figure out how not to print the whole test function at
		// all.
		testGen.Pln(`_ = ps`)
		for _, t := range tbls {
			codeWritten += t.generateTestOther(testGen, g)
		}
	}
	testGen.Pln(`}`) // end TestNewDBManager

	if codeWritten == 0 {
		testGen.Truncate(lenBefore)
	}
}

func (g *Generator) fnTestMainDB(testGen *codegen.Go, tbls tables) {
	if !tbls.hasFeature(g, FeatureDB) {
		return
	}

	// Test Header
	testGen.Pln(`func TestNewDBManagerDB_` + tbls.nameID() + `(t *testing.T) {`)
	{
		testGen.Pln(`db := dmltest.MustConnectDB(t)`)
		testGen.Pln(`defer dmltest.Close(t, db)`)

		if g.TestSQLDumpGlobPath != "" {
			testGen.Pln(`defer dmltest.SQLDumpLoad(t,`, strconv.Quote(g.TestSQLDumpGlobPath), `, &dmltest.SQLDumpOptions{
					SkipDBCleanup: true,
				}).Deferred()`)
		}

		testGen.Pln(`ctx, cancel := context.WithTimeout(context.Background(), time.Minute*2)`)
		testGen.Pln(`defer cancel()`)
		testGen.Pln(`tbls, err := NewDBManager(ctx, &DBMOption{TableOptions: []ddl.TableOption{ddl.WithConnPool(db)}} )`)
		testGen.Pln(`assert.NoError(t, err)`)

		testGen.Pln(`tblNames := tbls.Tables.Tables()`)
		testGen.Pln(`sort.Strings(tblNames)`)
		testGen.Pln(`assert.Exactly(t, `, fmt.Sprintf("%#v", g.sortedTableNames()), `, tblNames)`)

		testGen.Pln(`err = tbls.Validate(ctx)`)
		testGen.Pln(`assert.NoError(t, err)`)
		testGen.Pln(`var ps *pseudo.Service`)
		testGen.Pln(`ps = pseudo.MustNewService(0, &pseudo.Options{Lang: "de",MaxFloatDecimals:6},`)
		testGen.In()
		testGen.Pln(`pseudo.WithTagFakeFunc("website_id", func(maxLen int) interface{} {`)
		testGen.Pln(`    return 1`)
		testGen.Pln(`}),`)

		testGen.Pln(`pseudo.WithTagFakeFunc("store_id", func(maxLen int) interface{} {`)
		testGen.Pln(`    return 1`)
		testGen.Pln(`}),`)
		if fn, ok := g.customCode["pseudo.MustNewService.Option"]; ok {
			fn(g, nil, testGen)
		}
		testGen.Out()
		testGen.Pln(`)`)

		for _, t := range tbls {
			t.generateTestDB(testGen)
		} // end for tables
	}
	testGen.C(`Uncomment the next line for debugging to see all the queries.`)
	testGen.Pln(`// t.Logf("queries: %#v", tbls.ConnPool.CachedQueries())`)
	testGen.Pln(`}`) // end TestNewDBManager
}
