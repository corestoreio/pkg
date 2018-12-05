// Auto generated via github.com/corestoreio/pkg/sql/dmlgen

package testdata

import (
	"context"
	"github.com/bxcodec/faker"
	"github.com/corestoreio/pkg/sql/dmltest"
	"github.com/corestoreio/pkg/util/assert"
	"sort"
	"testing"
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

	// TODO run those tests in parallel
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
			assert.Exactly(t, entityIn.ConfigID, entityOut.ConfigID, "ConfigID did not match")
			assert.ExactlyLength(t, 8, &entityIn.Scope, &entityOut.Scope, "Scope do not match")
			assert.Exactly(t, entityIn.ScopeID, entityOut.ScopeID, "ScopeID did not match")
			assert.Exactly(t, entityIn.Expires, entityOut.Expires, "Expires did not match")
			assert.ExactlyLength(t, 255, &entityIn.Path, &entityOut.Path, "Path do not match")
			assert.ExactlyLength(t, 65535, &entityIn.Value, &entityOut.Value, "Value do not match")
			assert.Exactly(t, entityIn.VersionTs, entityOut.VersionTs, "VersionTs did not match")
			assert.Exactly(t, entityIn.VersionTe, entityOut.VersionTe, "VersionTe did not match")
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
			assert.Exactly(t, entityIn.EntityID, entityOut.EntityID, "EntityID did not match")
			assert.Exactly(t, entityIn.WebsiteID, entityOut.WebsiteID, "WebsiteID did not match")
			assert.ExactlyLength(t, 255, &entityIn.Email, &entityOut.Email, "Email do not match")
			assert.Exactly(t, entityIn.GroupID, entityOut.GroupID, "GroupID did not match")
			assert.ExactlyLength(t, 50, &entityIn.IncrementID, &entityOut.IncrementID, "IncrementID do not match")
			assert.Exactly(t, entityIn.StoreID, entityOut.StoreID, "StoreID did not match")
			assert.Exactly(t, entityIn.CreatedAt, entityOut.CreatedAt, "CreatedAt did not match")
			assert.Exactly(t, entityIn.UpdatedAt, entityOut.UpdatedAt, "UpdatedAt did not match")
			assert.Exactly(t, entityIn.IsActive, entityOut.IsActive, "IsActive did not match")
			assert.Exactly(t, entityIn.DisableAutoGroupChange, entityOut.DisableAutoGroupChange, "DisableAutoGroupChange did not match")
			assert.ExactlyLength(t, 255, &entityIn.CreatedIn, &entityOut.CreatedIn, "CreatedIn do not match")
			assert.ExactlyLength(t, 40, &entityIn.Prefix, &entityOut.Prefix, "Prefix do not match")
			assert.ExactlyLength(t, 255, &entityIn.Firstname, &entityOut.Firstname, "Firstname do not match")
			assert.ExactlyLength(t, 255, &entityIn.Middlename, &entityOut.Middlename, "Middlename do not match")
			assert.ExactlyLength(t, 255, &entityIn.Lastname, &entityOut.Lastname, "Lastname do not match")
			assert.ExactlyLength(t, 40, &entityIn.Suffix, &entityOut.Suffix, "Suffix do not match")
			assert.Exactly(t, entityIn.Dob, entityOut.Dob, "Dob did not match")
			assert.ExactlyLength(t, 128, &entityIn.PasswordHash, &entityOut.PasswordHash, "PasswordHash do not match")
			assert.ExactlyLength(t, 128, &entityIn.RpToken, &entityOut.RpToken, "RpToken do not match")
			assert.Exactly(t, entityIn.RpTokenCreatedAt, entityOut.RpTokenCreatedAt, "RpTokenCreatedAt did not match")
			assert.Exactly(t, entityIn.DefaultBilling, entityOut.DefaultBilling, "DefaultBilling did not match")
			assert.Exactly(t, entityIn.DefaultShipping, entityOut.DefaultShipping, "DefaultShipping did not match")
			assert.ExactlyLength(t, 50, &entityIn.Taxvat, &entityOut.Taxvat, "Taxvat do not match")
			assert.ExactlyLength(t, 64, &entityIn.Confirmation, &entityOut.Confirmation, "Confirmation do not match")
			assert.Exactly(t, entityIn.Gender, entityOut.Gender, "Gender did not match")
			assert.Exactly(t, entityIn.FailuresNum, entityOut.FailuresNum, "FailuresNum did not match")
			assert.Exactly(t, entityIn.FirstFailure, entityOut.FirstFailure, "FirstFailure did not match")
			assert.Exactly(t, entityIn.LockExpires, entityOut.LockExpires, "LockExpires did not match")
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
			assert.Exactly(t, entityIn.ID, entityOut.ID, "ID did not match")
			assert.Exactly(t, entityIn.ColBigint1, entityOut.ColBigint1, "ColBigint1 did not match")
			assert.Exactly(t, entityIn.ColBigint2, entityOut.ColBigint2, "ColBigint2 did not match")
			assert.Exactly(t, entityIn.ColBigint3, entityOut.ColBigint3, "ColBigint3 did not match")
			assert.Exactly(t, entityIn.ColBigint4, entityOut.ColBigint4, "ColBigint4 did not match")
			assert.ExactlyLength(t, 65535, &entityIn.ColBlob, &entityOut.ColBlob, "ColBlob do not match")
			assert.Exactly(t, entityIn.ColDate1, entityOut.ColDate1, "ColDate1 did not match")
			assert.Exactly(t, entityIn.ColDate2, entityOut.ColDate2, "ColDate2 did not match")
			assert.Exactly(t, entityIn.ColDatetime1, entityOut.ColDatetime1, "ColDatetime1 did not match")
			assert.Exactly(t, entityIn.ColDatetime2, entityOut.ColDatetime2, "ColDatetime2 did not match")
			assert.Exactly(t, entityIn.ColDecimal100, entityOut.ColDecimal100, "ColDecimal100 did not match")
			assert.Exactly(t, entityIn.ColDecimal124, entityOut.ColDecimal124, "ColDecimal124 did not match")
			assert.Exactly(t, entityIn.Price124a, entityOut.Price124a, "Price124a did not match")
			assert.Exactly(t, entityIn.Price124b, entityOut.Price124b, "Price124b did not match")
			assert.Exactly(t, entityIn.ColDecimal123, entityOut.ColDecimal123, "ColDecimal123 did not match")
			assert.Exactly(t, entityIn.ColDecimal206, entityOut.ColDecimal206, "ColDecimal206 did not match")
			assert.Exactly(t, entityIn.ColDecimal2412, entityOut.ColDecimal2412, "ColDecimal2412 did not match")
			assert.Exactly(t, entityIn.ColFloat, entityOut.ColFloat, "ColFloat did not match")
			assert.Exactly(t, entityIn.ColInt1, entityOut.ColInt1, "ColInt1 did not match")
			assert.Exactly(t, entityIn.ColInt2, entityOut.ColInt2, "ColInt2 did not match")
			assert.Exactly(t, entityIn.ColInt3, entityOut.ColInt3, "ColInt3 did not match")
			assert.Exactly(t, entityIn.ColInt4, entityOut.ColInt4, "ColInt4 did not match")
			assert.ExactlyLength(t, 4294967295, &entityIn.ColLongtext1, &entityOut.ColLongtext1, "ColLongtext1 do not match")
			assert.ExactlyLength(t, 4294967295, &entityIn.ColLongtext2, &entityOut.ColLongtext2, "ColLongtext2 do not match")
			assert.ExactlyLength(t, 16777215, &entityIn.ColMediumblob, &entityOut.ColMediumblob, "ColMediumblob do not match")
			assert.ExactlyLength(t, 16777215, &entityIn.ColMediumtext1, &entityOut.ColMediumtext1, "ColMediumtext1 do not match")
			assert.ExactlyLength(t, 16777215, &entityIn.ColMediumtext2, &entityOut.ColMediumtext2, "ColMediumtext2 do not match")
			assert.Exactly(t, entityIn.ColSmallint1, entityOut.ColSmallint1, "ColSmallint1 did not match")
			assert.Exactly(t, entityIn.ColSmallint2, entityOut.ColSmallint2, "ColSmallint2 did not match")
			assert.Exactly(t, entityIn.ColSmallint3, entityOut.ColSmallint3, "ColSmallint3 did not match")
			assert.Exactly(t, entityIn.ColSmallint4, entityOut.ColSmallint4, "ColSmallint4 did not match")
			assert.Exactly(t, entityIn.HasSmallint5, entityOut.HasSmallint5, "HasSmallint5 did not match")
			assert.Exactly(t, entityIn.IsSmallint5, entityOut.IsSmallint5, "IsSmallint5 did not match")
			assert.ExactlyLength(t, 65535, &entityIn.ColText, &entityOut.ColText, "ColText do not match")
			assert.Exactly(t, entityIn.ColTimestamp1, entityOut.ColTimestamp1, "ColTimestamp1 did not match")
			assert.Exactly(t, entityIn.ColTimestamp2, entityOut.ColTimestamp2, "ColTimestamp2 did not match")
			assert.Exactly(t, entityIn.ColTinyint1, entityOut.ColTinyint1, "ColTinyint1 did not match")
			assert.ExactlyLength(t, 1, &entityIn.ColVarchar1, &entityOut.ColVarchar1, "ColVarchar1 do not match")
			assert.ExactlyLength(t, 100, &entityIn.ColVarchar100, &entityOut.ColVarchar100, "ColVarchar100 do not match")
			assert.ExactlyLength(t, 16, &entityIn.ColVarchar16, &entityOut.ColVarchar16, "ColVarchar16 do not match")
			assert.ExactlyLength(t, 21, &entityIn.ColChar1, &entityOut.ColChar1, "ColChar1 do not match")
			assert.ExactlyLength(t, 17, &entityIn.ColChar2, &entityOut.ColChar2, "ColChar2 do not match")
		}
	})
}
