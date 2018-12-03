// Auto generated via github.com/corestoreio/pkg/sql/dmlgen

package testdata

import (
	"testing"
	"context"
	"sort"
	"github.com/corestoreio/pkg/sql/dmltest"
	"github.com/bxcodec/faker"
	"github.com/corestoreio/pkg/util/assert"

)
func TestNewTables(t *testing.T) {
	db := dmltest.MustConnectDB(t)
	defer dmltest.Close(t, db)

	defer dmltest.SQLDumpLoad(t, "test_*_tables.sql", &dmltest.SQLDumpOptions{
		SkipDBCleanup: true,
	})()

	ctx := context.TODO()
	tbls, err := NewTables(ctx, db.DB)
	assert.NoError(t, err)

	tblNames := tbls.Tables()
	sort.Strings(tblNames)
	assert.Exactly(t, []string{"core_config_data", "customer_entity", "dmlgen_types"}, tblNames)

	err = tbls.Validate(ctx)
	assert.NoError(t, err)
	t.Run("CoreConfigData_Entity", func(t *testing.T) {
		ccd := tbls.MustTable(TableNameCoreConfigData)

		inStmt, err := ccd.Insert().Ignore().BuildValues().Prepare(ctx)
		assert.NoError(t, err, "%+v", err)
		insArtisan := inStmt.WithArgs()
		defer dmltest.Close(t, inStmt)

		selArtisan := ccd.SelectByPK().WithArgs().ExpandPlaceHolders()

		for i := 0; i < 5; i++ {
			entityIn := new(CoreConfigData)
			assert.NoError(t, faker.FakeData(entityIn))

			lID := dmltest.CheckLastInsertID(t, "Error: TestNewTables.CoreConfigData_Entity")(insArtisan.Record("", entityIn).ExecContext(ctx))
			insArtisan.Reset()

			entityOut := new(CoreConfigData)
			rowCount, err := selArtisan.Int64s(lID).Load(ctx, entityOut)
			assert.NoError(t, err, "%+v", err)
			assert.Exactly(t, uint64(1), rowCount, "RowCount did not match")

			// assert.Exactly(t, entityIn.Scope, entityOut.Scope, "Scope did not match")
			// assert.Exactly(t, entityIn.ScopeID, entityOut.ScopeID, "ScopeID did not match")
			// assert.Exactly(t, entityIn.Expires, entityOut.Expires, "Expires did not match")
			// assert.Exactly(t, entityIn.Path, entityOut.Path, "Path did not match")
			// assert.Exactly(t, entityIn.Value, entityOut.Value, "Value did not match")
		}
	})
	t.Run("CustomerEntity_Entity", func(t *testing.T) {
		ccd := tbls.MustTable(TableNameCustomerEntity)

		inStmt, err := ccd.Insert().Ignore().BuildValues().Prepare(ctx)
		assert.NoError(t, err, "%+v", err)
		insArtisan := inStmt.WithArgs()
		defer dmltest.Close(t, inStmt)

		selArtisan := ccd.SelectByPK().WithArgs().ExpandPlaceHolders()

		for i := 0; i < 5; i++ {
			entityIn := new(CustomerEntity)
			assert.NoError(t, faker.FakeData(entityIn))

			lID := dmltest.CheckLastInsertID(t, "Error: TestNewTables.CustomerEntity_Entity")(insArtisan.Record("", entityIn).ExecContext(ctx))
			insArtisan.Reset()

			entityOut := new(CustomerEntity)
			rowCount, err := selArtisan.Int64s(lID).Load(ctx, entityOut)
			assert.NoError(t, err, "%+v", err)
			assert.Exactly(t, uint64(1), rowCount, "RowCount did not match")

			// assert.Exactly(t, entityIn.Scope, entityOut.Scope, "Scope did not match")
			// assert.Exactly(t, entityIn.ScopeID, entityOut.ScopeID, "ScopeID did not match")
			// assert.Exactly(t, entityIn.Expires, entityOut.Expires, "Expires did not match")
			// assert.Exactly(t, entityIn.Path, entityOut.Path, "Path did not match")
			// assert.Exactly(t, entityIn.Value, entityOut.Value, "Value did not match")
		}
	})
	t.Run("DmlgenTypes_Entity", func(t *testing.T) {
		ccd := tbls.MustTable(TableNameDmlgenTypes)

		inStmt, err := ccd.Insert().Ignore().BuildValues().Prepare(ctx)
		assert.NoError(t, err, "%+v", err)
		insArtisan := inStmt.WithArgs()
		defer dmltest.Close(t, inStmt)

		selArtisan := ccd.SelectByPK().WithArgs().ExpandPlaceHolders()

		for i := 0; i < 5; i++ {
			entityIn := new(DmlgenTypes)
			assert.NoError(t, faker.FakeData(entityIn))

			lID := dmltest.CheckLastInsertID(t, "Error: TestNewTables.DmlgenTypes_Entity")(insArtisan.Record("", entityIn).ExecContext(ctx))
			insArtisan.Reset()

			entityOut := new(DmlgenTypes)
			rowCount, err := selArtisan.Int64s(lID).Load(ctx, entityOut)
			assert.NoError(t, err, "%+v", err)
			assert.Exactly(t, uint64(1), rowCount, "RowCount did not match")

			// assert.Exactly(t, entityIn.Scope, entityOut.Scope, "Scope did not match")
			// assert.Exactly(t, entityIn.ScopeID, entityOut.ScopeID, "ScopeID did not match")
			// assert.Exactly(t, entityIn.Expires, entityOut.Expires, "Expires did not match")
			// assert.Exactly(t, entityIn.Path, entityOut.Path, "Path did not match")
			// assert.Exactly(t, entityIn.Value, entityOut.Value, "Value did not match")
		}
	})
}
