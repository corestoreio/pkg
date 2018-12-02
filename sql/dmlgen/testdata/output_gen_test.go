// Auto generated via github.com/corestoreio/pkg/sql/dmlgen

package testdata_test

import (
	"context"
	"sort"
	"testing"

	"github.com/bxcodec/faker"
	"github.com/corestoreio/pkg/sql/dmlgen/testdata"
	"github.com/corestoreio/pkg/sql/dmltest"
	"github.com/corestoreio/pkg/util/assert"
)

func TestNewTables(t *testing.T) {
	db := dmltest.MustConnectDB(t)
	defer dmltest.Close(t, db)

	defer dmltest.SQLDumpLoad(t, "test_*_tables.sql", &dmltest.SQLDumpOptions{
		SkipDBCleanup: true,
	})()

	ctx := context.TODO()
	tbls, err := testdata.NewTables(ctx, db.DB)
	assert.NoError(t, err)

	tblNames := tbls.Tables()
	sort.Strings(tblNames)
	assert.Exactly(t, []string{"core_config_data", "customer_entity", "dmlgen_types"}, tblNames)

	err = tbls.Validate(ctx)
	assert.NoError(t, err)

	t.Run("CoreConfigData_Entity", func(t *testing.T) {
		ccd := tbls.MustTable(testdata.TableNameCoreConfigData)

		inStmt, err := ccd.Insert().Ignore().BuildValues().Prepare(ctx)
		assert.NoError(t, err, "%+v", err)
		insArtisan := inStmt.WithArgs()
		defer dmltest.Close(t, inStmt)

		selArtisan := ccd.SelectByPK().WithArgs().ExpandPlaceHolders()

		for i := 0; i < 5; i++ {
			entity := new(testdata.CoreConfigData)
			assert.NoError(t, faker.FakeData(entity))

			lID := dmltest.CheckLastInsertID(t, "Error: TestNewTables.CoreConfigData_Entity")(insArtisan.Record("", entity).ExecContext(ctx))
			insArtisan.Reset()
			entity.Empty()

			rowCount, err := selArtisan.Int64s(lID).Load(ctx, entity)
			assert.NoError(t, err, "%+v", err)
			assert.Exactly(t, uint64(1), rowCount, "RowCount did not match")

			assert.Exactly(t, entity.Scope, entity.Scope, "Scope did not match")
			assert.Exactly(t, entity.ScopeID, entity.ScopeID, "ScopeID did not match")
			assert.Exactly(t, entity.Expires, entity.Expires, "Expires did not match")
			assert.Exactly(t, entity.Path, entity.Path, "Path did not match")
			assert.Exactly(t, entity.Value, entity.Value, "Value did not match")
		}
	})
}
