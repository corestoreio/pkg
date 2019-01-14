// Auto generated via github.com/corestoreio/pkg/sql/dmlgen

package testdata

import (
	"context"
	"fmt"
	"github.com/alecthomas/repr"
	"github.com/corestoreio/pkg/sql/dmltest"
	"github.com/corestoreio/pkg/util/assert"
	"github.com/corestoreio/pkg/util/pseudo"
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
	assert.Exactly(t, []string{"core_config_data", "customer_address_entity", "customer_entity", "dmlgen_types"}, tblNames)

	err = tbls.Validate(ctx)
	assert.NoError(t, err)
	var ps *pseudo.Service
	ps = pseudo.MustNewService(0, &pseudo.Options{Lang: "de", FloatMaxDecimals: 6},
		pseudo.WithTagFakeFunc("website_id", func(maxLen int) (interface{}, error) {
			return 1, nil
		}),
		pseudo.WithTagFakeFunc("store_id", func(maxLen int) (interface{}, error) {
			return 1, nil
		}),
		pseudo.WithTagFakeFunc("CustomerAddressEntity.ParentID", func(maxLen int) (interface{}, error) {
			return nil, nil
		}),
		pseudo.WithTagFakeFunc("col_date1", func(maxLen int) (interface{}, error) {
			if ps.Intn(1000)%3 == 0 {
				return nil, nil
			}
			return ps.Dob18(), nil
		}),
		pseudo.WithTagFakeFunc("col_date2", func(maxLen int) (interface{}, error) {
			return ps.Dob18().MarshalText()
		}),
		pseudo.WithTagFakeFunc("col_decimal101", func(maxLen int) (interface{}, error) {
			return fmt.Sprintf("%.1f", ps.Price()), nil
		}),
		pseudo.WithTagFakeFunc("price124b", func(maxLen int) (interface{}, error) {
			return fmt.Sprintf("%.4f", ps.Price()), nil
		}),
		pseudo.WithTagFakeFunc("col_decimal123", func(maxLen int) (interface{}, error) {
			return fmt.Sprintf("%.3f", ps.Float64()), nil
		}),
		pseudo.WithTagFakeFunc("col_decimal206", func(maxLen int) (interface{}, error) {
			return fmt.Sprintf("%.6f", ps.Float64()), nil
		}),
		pseudo.WithTagFakeFunc("col_decimal2412", func(maxLen int) (interface{}, error) {
			return fmt.Sprintf("%.12f", ps.Float64()), nil
		}),
		pseudo.WithTagFakeFuncAlias(
			"col_decimal124", "price124b",
			"price124a", "price124b",
			"col_float", "col_decimal206",
		),
	)

	// TODO run those tests in parallel
	t.Run("CoreConfigData_Entity", func(t *testing.T) {
		ccd := tbls.MustTable(TableNameCoreConfigData)

		inStmt, err := ccd.Insert().BuildValues().Prepare(ctx) // Do not use Ignore() to suppress DB errors.
		assert.NoError(t, err, "%+v", err)
		insArtisan := inStmt.WithArgs()
		defer dmltest.Close(t, inStmt)

		selArtisan := ccd.SelectByPK().WithArgs().ExpandPlaceHolders()

		for i := 0; i < 9; i++ {
			entityIn := new(CoreConfigData)
			if err := ps.FakeData(entityIn); err != nil {
				t.Errorf("IDX[%d]: %+v", i, err)
				return
			}

			lID := dmltest.CheckLastInsertID(t, "Error: TestNewTables.CoreConfigData_Entity")(insArtisan.Record("", entityIn).ExecContext(ctx))
			insArtisan.Reset()

			entityOut := new(CoreConfigData)
			rowCount, err := selArtisan.Int64s(lID).Load(ctx, entityOut)
			assert.NoError(t, err, "%+v", err)
			assert.Exactly(t, uint64(1), rowCount, "IDX%d: RowCount did not match", i)
			assert.Exactly(t, entityIn.ConfigID, entityOut.ConfigID, "IDX%d: ConfigID should match", lID)
			assert.ExactlyLength(t, 8, &entityIn.Scope, &entityOut.Scope, "IDX%d: Scope should match", lID)
			assert.Exactly(t, entityIn.ScopeID, entityOut.ScopeID, "IDX%d: ScopeID should match", lID)
			assert.Exactly(t, entityIn.Expires, entityOut.Expires, "IDX%d: Expires should match", lID)
			assert.ExactlyLength(t, 255, &entityIn.Path, &entityOut.Path, "IDX%d: Path should match", lID)
			assert.ExactlyLength(t, 65535, &entityIn.Value, &entityOut.Value, "IDX%d: Value should match", lID)

		}
	})
	t.Run("CustomerAddressEntity_Entity", func(t *testing.T) {
		ccd := tbls.MustTable(TableNameCustomerAddressEntity)

		inStmt, err := ccd.Insert().BuildValues().Prepare(ctx) // Do not use Ignore() to suppress DB errors.
		assert.NoError(t, err, "%+v", err)
		insArtisan := inStmt.WithArgs()
		defer dmltest.Close(t, inStmt)

		selArtisan := ccd.SelectByPK().WithArgs().ExpandPlaceHolders()

		for i := 0; i < 9; i++ {
			entityIn := new(CustomerAddressEntity)
			if err := ps.FakeData(entityIn); err != nil {
				t.Errorf("IDX[%d]: %+v", i, err)
				return
			}

			lID := dmltest.CheckLastInsertID(t, "Error: TestNewTables.CustomerAddressEntity_Entity")(insArtisan.Record("", entityIn).ExecContext(ctx))
			insArtisan.Reset()

			entityOut := new(CustomerAddressEntity)
			rowCount, err := selArtisan.Int64s(lID).Load(ctx, entityOut)
			assert.NoError(t, err, "%+v", err)
			assert.Exactly(t, uint64(1), rowCount, "IDX%d: RowCount did not match", i)
			assert.Exactly(t, entityIn.EntityID, entityOut.EntityID, "IDX%d: EntityID should match", lID)
			assert.ExactlyLength(t, 50, &entityIn.IncrementID, &entityOut.IncrementID, "IDX%d: IncrementID should match", lID)
			assert.Exactly(t, entityIn.ParentID, entityOut.ParentID, "IDX%d: ParentID should match", lID)
			assert.Exactly(t, entityIn.CreatedAt, entityOut.CreatedAt, "IDX%d: CreatedAt should match", lID)
			assert.Exactly(t, entityIn.UpdatedAt, entityOut.UpdatedAt, "IDX%d: UpdatedAt should match", lID)
			assert.Exactly(t, entityIn.IsActive, entityOut.IsActive, "IDX%d: IsActive should match", lID)
			assert.ExactlyLength(t, 255, &entityIn.City, &entityOut.City, "IDX%d: City should match", lID)
			assert.ExactlyLength(t, 255, &entityIn.Company, &entityOut.Company, "IDX%d: Company should match", lID)
			assert.ExactlyLength(t, 255, &entityIn.CountryID, &entityOut.CountryID, "IDX%d: CountryID should match", lID)
			assert.ExactlyLength(t, 255, &entityIn.Fax, &entityOut.Fax, "IDX%d: Fax should match", lID)
			assert.ExactlyLength(t, 255, &entityIn.Firstname, &entityOut.Firstname, "IDX%d: Firstname should match", lID)
			assert.ExactlyLength(t, 255, &entityIn.Lastname, &entityOut.Lastname, "IDX%d: Lastname should match", lID)
			assert.ExactlyLength(t, 255, &entityIn.Middlename, &entityOut.Middlename, "IDX%d: Middlename should match", lID)
			assert.ExactlyLength(t, 255, &entityIn.Postcode, &entityOut.Postcode, "IDX%d: Postcode should match", lID)
			assert.ExactlyLength(t, 40, &entityIn.Prefix, &entityOut.Prefix, "IDX%d: Prefix should match", lID)
			assert.ExactlyLength(t, 255, &entityIn.Region, &entityOut.Region, "IDX%d: Region should match", lID)
			assert.Exactly(t, entityIn.RegionID, entityOut.RegionID, "IDX%d: RegionID should match", lID)
			assert.ExactlyLength(t, 65535, &entityIn.Street, &entityOut.Street, "IDX%d: Street should match", lID)
			assert.ExactlyLength(t, 40, &entityIn.Suffix, &entityOut.Suffix, "IDX%d: Suffix should match", lID)
			assert.ExactlyLength(t, 255, &entityIn.Telephone, &entityOut.Telephone, "IDX%d: Telephone should match", lID)
			assert.ExactlyLength(t, 255, &entityIn.VatID, &entityOut.VatID, "IDX%d: VatID should match", lID)
			assert.Exactly(t, entityIn.VatIsValid, entityOut.VatIsValid, "IDX%d: VatIsValid should match", lID)
			assert.ExactlyLength(t, 255, &entityIn.VatRequestDate, &entityOut.VatRequestDate, "IDX%d: VatRequestDate should match", lID)
			assert.ExactlyLength(t, 255, &entityIn.VatRequestID, &entityOut.VatRequestID, "IDX%d: VatRequestID should match", lID)
			assert.Exactly(t, entityIn.VatRequestSuccess, entityOut.VatRequestSuccess, "IDX%d: VatRequestSuccess should match", lID)
		}
	})
	t.Run("CustomerEntity_Entity", func(t *testing.T) {
		ccd := tbls.MustTable(TableNameCustomerEntity)

		inStmt, err := ccd.Insert().BuildValues().Prepare(ctx) // Do not use Ignore() to suppress DB errors.
		assert.NoError(t, err, "%+v", err)
		insArtisan := inStmt.WithArgs()
		defer dmltest.Close(t, inStmt)

		selArtisan := ccd.SelectByPK().WithArgs().ExpandPlaceHolders()

		for i := 0; i < 1; i++ {
			entityIn := new(CustomerEntity)
			if err := ps.FakeData(entityIn); err != nil {
				t.Errorf("IDX[%d]: %+v", i, err)
				return
			}

			repr.Println(entityIn)

			lID := dmltest.CheckLastInsertID(t, "Error: TestNewTables.CustomerEntity_Entity")(insArtisan.Record("", entityIn).ExecContext(ctx))
			insArtisan.Reset()

			entityOut := new(CustomerEntity)
			rowCount, err := selArtisan.Int64s(lID).Load(ctx, entityOut)
			assert.NoError(t, err, "%+v", err)
			assert.Exactly(t, uint64(1), rowCount, "IDX%d: RowCount did not match", i)
			assert.Exactly(t, entityIn.EntityID, entityOut.EntityID, "IDX%d: EntityID should match", lID)
			assert.Exactly(t, entityIn.WebsiteID, entityOut.WebsiteID, "IDX%d: WebsiteID should match", lID)
			assert.ExactlyLength(t, 255, &entityIn.Email, &entityOut.Email, "IDX%d: Email should match", lID)
			assert.Exactly(t, entityIn.GroupID, entityOut.GroupID, "IDX%d: GroupID should match", lID)
			assert.ExactlyLength(t, 50, &entityIn.IncrementID, &entityOut.IncrementID, "IDX%d: IncrementID should match", lID)
			assert.Exactly(t, entityIn.StoreID, entityOut.StoreID, "IDX%d: StoreID should match", lID)
			assert.Exactly(t, entityIn.CreatedAt, entityOut.CreatedAt, "IDX%d: CreatedAt should match", lID)
			assert.Exactly(t, entityIn.UpdatedAt, entityOut.UpdatedAt, "IDX%d: UpdatedAt should match", lID)
			assert.Exactly(t, entityIn.IsActive, entityOut.IsActive, "IDX%d: IsActive should match", lID)
			assert.Exactly(t, entityIn.DisableAutoGroupChange, entityOut.DisableAutoGroupChange, "IDX%d: DisableAutoGroupChange should match", lID)
			assert.ExactlyLength(t, 255, &entityIn.CreatedIn, &entityOut.CreatedIn, "IDX%d: CreatedIn should match", lID)
			assert.ExactlyLength(t, 40, &entityIn.Prefix, &entityOut.Prefix, "IDX%d: Prefix should match", lID)
			assert.ExactlyLength(t, 255, &entityIn.Firstname, &entityOut.Firstname, "IDX%d: Firstname should match", lID)
			assert.ExactlyLength(t, 255, &entityIn.Middlename, &entityOut.Middlename, "IDX%d: Middlename should match", lID)
			assert.ExactlyLength(t, 255, &entityIn.Lastname, &entityOut.Lastname, "IDX%d: Lastname should match", lID)
			assert.ExactlyLength(t, 40, &entityIn.Suffix, &entityOut.Suffix, "IDX%d: Suffix should match", lID)
			assert.Exactly(t, entityIn.Dob, entityOut.Dob, "IDX%d: Dob should match", lID)
			assert.ExactlyLength(t, 128, &entityIn.PasswordHash, &entityOut.PasswordHash, "IDX%d: PasswordHash should match", lID)
			assert.ExactlyLength(t, 128, &entityIn.RpToken, &entityOut.RpToken, "IDX%d: RpToken should match", lID)
			assert.Exactly(t, entityIn.RpTokenCreatedAt, entityOut.RpTokenCreatedAt, "IDX%d: RpTokenCreatedAt should match", lID)
			assert.Exactly(t, entityIn.DefaultBilling, entityOut.DefaultBilling, "IDX%d: DefaultBilling should match", lID)
			assert.Exactly(t, entityIn.DefaultShipping, entityOut.DefaultShipping, "IDX%d: DefaultShipping should match", lID)
			assert.ExactlyLength(t, 50, &entityIn.Taxvat, &entityOut.Taxvat, "IDX%d: Taxvat should match", lID)
			assert.ExactlyLength(t, 64, &entityIn.Confirmation, &entityOut.Confirmation, "IDX%d: Confirmation should match", lID)
			assert.Exactly(t, entityIn.Gender, entityOut.Gender, "IDX%d: Gender should match", lID)
			assert.Exactly(t, entityIn.FailuresNum, entityOut.FailuresNum, "IDX%d: FailuresNum should match", lID)
			assert.Exactly(t, entityIn.FirstFailure, entityOut.FirstFailure, "IDX%d: FirstFailure should match", lID)
			assert.Exactly(t, entityIn.LockExpires, entityOut.LockExpires, "IDX%d: LockExpires should match", lID)
		}
	})
	t.Run("DmlgenTypes_Entity", func(t *testing.T) {
		ccd := tbls.MustTable(TableNameDmlgenTypes)

		inStmt, err := ccd.Insert().BuildValues().Prepare(ctx) // Do not use Ignore() to suppress DB errors.
		assert.NoError(t, err, "%+v", err)
		insArtisan := inStmt.WithArgs()
		defer dmltest.Close(t, inStmt)

		selArtisan := ccd.SelectByPK().WithArgs().ExpandPlaceHolders()

		for i := 0; i < 9; i++ {
			entityIn := new(DmlgenTypes)
			if err := ps.FakeData(entityIn); err != nil {
				t.Errorf("IDX[%d]: %+v", i, err)
				return
			}

			lID := dmltest.CheckLastInsertID(t, "Error: TestNewTables.DmlgenTypes_Entity")(insArtisan.Record("", entityIn).ExecContext(ctx))
			insArtisan.Reset()

			entityOut := new(DmlgenTypes)
			rowCount, err := selArtisan.Int64s(lID).Load(ctx, entityOut)
			assert.NoError(t, err, "%+v", err)
			assert.Exactly(t, uint64(1), rowCount, "IDX%d: RowCount did not match", i)
			assert.Exactly(t, entityIn.ID, entityOut.ID, "IDX%d: ID should match", lID)
			assert.Exactly(t, entityIn.ColBigint1, entityOut.ColBigint1, "IDX%d: ColBigint1 should match", lID)
			assert.Exactly(t, entityIn.ColBigint2, entityOut.ColBigint2, "IDX%d: ColBigint2 should match", lID)
			assert.Exactly(t, entityIn.ColBigint3, entityOut.ColBigint3, "IDX%d: ColBigint3 should match", lID)
			assert.Exactly(t, entityIn.ColBigint4, entityOut.ColBigint4, "IDX%d: ColBigint4 should match", lID)
			assert.ExactlyLength(t, 65535, &entityIn.ColBlob, &entityOut.ColBlob, "IDX%d: ColBlob should match", lID)
			assert.Exactly(t, entityIn.ColDate1, entityOut.ColDate1, "IDX%d: ColDate1 should match", lID)
			assert.Exactly(t, entityIn.ColDate2, entityOut.ColDate2, "IDX%d: ColDate2 should match", lID)
			assert.Exactly(t, entityIn.ColDatetime1, entityOut.ColDatetime1, "IDX%d: ColDatetime1 should match", lID)
			assert.Exactly(t, entityIn.ColDatetime2, entityOut.ColDatetime2, "IDX%d: ColDatetime2 should match", lID)
			assert.Exactly(t, entityIn.ColDecimal101, entityOut.ColDecimal101, "IDX%d: ColDecimal101 should match", lID)
			assert.Exactly(t, entityIn.ColDecimal124, entityOut.ColDecimal124, "IDX%d: ColDecimal124 should match", lID)
			assert.Exactly(t, entityIn.Price124a, entityOut.Price124a, "IDX%d: Price124a should match", lID)
			assert.Exactly(t, entityIn.Price124b, entityOut.Price124b, "IDX%d: Price124b should match", lID)
			assert.Exactly(t, entityIn.ColDecimal123, entityOut.ColDecimal123, "IDX%d: ColDecimal123 should match", lID)
			assert.Exactly(t, entityIn.ColDecimal206, entityOut.ColDecimal206, "IDX%d: ColDecimal206 should match", lID)
			assert.Exactly(t, entityIn.ColDecimal2412, entityOut.ColDecimal2412, "IDX%d: ColDecimal2412 should match", lID)
			assert.Exactly(t, entityIn.ColInt1, entityOut.ColInt1, "IDX%d: ColInt1 should match", lID)
			assert.Exactly(t, entityIn.ColInt2, entityOut.ColInt2, "IDX%d: ColInt2 should match", lID)
			assert.Exactly(t, entityIn.ColInt3, entityOut.ColInt3, "IDX%d: ColInt3 should match", lID)
			assert.Exactly(t, entityIn.ColInt4, entityOut.ColInt4, "IDX%d: ColInt4 should match", lID)
			assert.ExactlyLength(t, 4294967295, &entityIn.ColLongtext1, &entityOut.ColLongtext1, "IDX%d: ColLongtext1 should match", lID)
			assert.ExactlyLength(t, 4294967295, &entityIn.ColLongtext2, &entityOut.ColLongtext2, "IDX%d: ColLongtext2 should match", lID)
			assert.ExactlyLength(t, 16777215, &entityIn.ColMediumblob, &entityOut.ColMediumblob, "IDX%d: ColMediumblob should match", lID)
			assert.ExactlyLength(t, 16777215, &entityIn.ColMediumtext1, &entityOut.ColMediumtext1, "IDX%d: ColMediumtext1 should match", lID)
			assert.ExactlyLength(t, 16777215, &entityIn.ColMediumtext2, &entityOut.ColMediumtext2, "IDX%d: ColMediumtext2 should match", lID)
			assert.Exactly(t, entityIn.ColSmallint1, entityOut.ColSmallint1, "IDX%d: ColSmallint1 should match", lID)
			assert.Exactly(t, entityIn.ColSmallint2, entityOut.ColSmallint2, "IDX%d: ColSmallint2 should match", lID)
			assert.Exactly(t, entityIn.ColSmallint3, entityOut.ColSmallint3, "IDX%d: ColSmallint3 should match", lID)
			assert.Exactly(t, entityIn.ColSmallint4, entityOut.ColSmallint4, "IDX%d: ColSmallint4 should match", lID)
			assert.Exactly(t, entityIn.HasSmallint5, entityOut.HasSmallint5, "IDX%d: HasSmallint5 should match", lID)
			assert.Exactly(t, entityIn.IsSmallint5, entityOut.IsSmallint5, "IDX%d: IsSmallint5 should match", lID)
			assert.ExactlyLength(t, 65535, &entityIn.ColText, &entityOut.ColText, "IDX%d: ColText should match", lID)
			assert.Exactly(t, entityIn.ColTimestamp1, entityOut.ColTimestamp1, "IDX%d: ColTimestamp1 should match", lID)
			assert.Exactly(t, entityIn.ColTimestamp2, entityOut.ColTimestamp2, "IDX%d: ColTimestamp2 should match", lID)
			assert.Exactly(t, entityIn.ColTinyint1, entityOut.ColTinyint1, "IDX%d: ColTinyint1 should match", lID)
			assert.ExactlyLength(t, 1, &entityIn.ColVarchar1, &entityOut.ColVarchar1, "IDX%d: ColVarchar1 should match", lID)
			assert.ExactlyLength(t, 100, &entityIn.ColVarchar100, &entityOut.ColVarchar100, "IDX%d: ColVarchar100 should match", lID)
			assert.ExactlyLength(t, 16, &entityIn.ColVarchar16, &entityOut.ColVarchar16, "IDX%d: ColVarchar16 should match", lID)
			assert.ExactlyLength(t, 21, &entityIn.ColChar1, &entityOut.ColChar1, "IDX%d: ColChar1 should match", lID)
			assert.ExactlyLength(t, 17, &entityIn.ColChar2, &entityOut.ColChar2, "IDX%d: ColChar2 should match", lID)
		}
	})
}
