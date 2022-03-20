// Code generated by corestoreio/pkg/util/codegen. DO NOT EDIT.
// Generated by sql/dmlgen. DO NOT EDIT.
package dmltestgeneratedMToM

import (
	"context"
	"sort"
	"testing"
	"time"

	"github.com/corestoreio/pkg/sql/ddl"
	"github.com/corestoreio/pkg/sql/dml"
	"github.com/corestoreio/pkg/sql/dmltest"
	"github.com/corestoreio/pkg/util/assert"
	"github.com/corestoreio/pkg/util/pseudo"
)

func TestNewDBManagerDB_b0014848206a53c73619e6569577340a(t *testing.T) {
	db := dmltest.MustConnectDB(t)
	defer dmltest.Close(t, db)
	defer dmltest.SQLDumpLoad(t, "../testdata/testAll_*_tables.sql", &dmltest.SQLDumpOptions{
		SkipDBCleanup: true,
	}).Deferred()
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute*2)
	defer cancel()
	tbls, err := NewDBManager(ctx, &DBMOption{TableOptions: []ddl.TableOption{ddl.WithConnPool(db)}})
	assert.NoError(t, err)
	tblNames := tbls.Tables.Tables()
	sort.Strings(tblNames)
	assert.Exactly(t, []string{"athlete", "athlete_team", "athlete_team_member", "customer_address_entity", "customer_entity"}, tblNames)
	err = tbls.Validate(ctx)
	assert.NoError(t, err)
	var ps *pseudo.Service
	ps = pseudo.MustNewService(0, &pseudo.Options{Lang: "de", MaxFloatDecimals: 6},
		pseudo.WithTagFakeFunc("website_id", func(maxLen int) any {
			return 1
		}),
		pseudo.WithTagFakeFunc("store_id", func(maxLen int) any {
			return 1
		}),
	)
	t.Run("Athlete_Entity", func(t *testing.T) {
		tbl := tbls.MustTable(TableNameAthlete)
		selOneRow := tbl.Select("*").Where(
			dml.Column("athlete_id").Equal().PlaceHolder(),
		)
		selTenRows := tbl.Select("*").Where(
			dml.Column("athlete_id").LessOrEqual().Int(10),
		)
		selOneRowDBR := tbls.ConnPool.WithPrepare(ctx, selOneRow)
		defer selOneRowDBR.Close()
		selTenRowsDBR := tbls.ConnPool.WithQueryBuilder(selTenRows)
		entINSERTStmtA := tbls.ConnPool.WithPrepare(ctx, tbl.Insert().BuildValues())
		for i := 0; i < 9; i++ {
			entIn := new(Athlete)
			assert.NoError(t, ps.FakeData(entIn), "Error at index %d", i)
			lID := dmltest.CheckLastInsertID(t, "Error: TestNewTables.Athlete_Entity")(entINSERTStmtA.ExecContext(ctx, dml.Qualify("", entIn)))
			entINSERTStmtA.Reset()
			entOut := new(Athlete)
			rowCount, err := selOneRowDBR.Load(ctx, entOut, lID)
			assert.NoError(t, err)
			assert.Exactly(t, uint64(1), rowCount, "IDX%d: RowCount did not match", i)
			assert.Exactly(t, entIn.AthleteID, entOut.AthleteID, "IDX%d: AthleteID should match", lID)
			assert.ExactlyLength(t, 340, &entIn.Firstname, &entOut.Firstname, "IDX%d: Firstname should match", lID)
			assert.ExactlyLength(t, 340, &entIn.Lastname, &entOut.Lastname, "IDX%d: Lastname should match", lID)
		}
		dmltest.Close(t, entINSERTStmtA)
		entCol := NewAthletes()
		rowCount, err := selTenRowsDBR.Load(ctx, entCol)
		assert.NoError(t, err)
		t.Logf("Collection load rowCount: %d", rowCount)
		colInsertDBR := tbls.ConnPool.WithQueryBuilder(tbl.Insert().Replace().SetRowCount(len(entCol.Data)).BuildValues())
		lID := dmltest.CheckLastInsertID(t, "Error:  Athletes ")(colInsertDBR.ExecContext(ctx, dml.Qualify("", entCol)))
		t.Logf("Last insert ID into: %d", lID)
	})
	t.Run("AthleteTeam_Entity", func(t *testing.T) {
		tbl := tbls.MustTable(TableNameAthleteTeam)
		selOneRow := tbl.Select("*").Where(
			dml.Column("team_id").Equal().PlaceHolder(),
		)
		selTenRows := tbl.Select("*").Where(
			dml.Column("team_id").LessOrEqual().Int(10),
		)
		selOneRowDBR := tbls.ConnPool.WithPrepare(ctx, selOneRow)
		defer selOneRowDBR.Close()
		selTenRowsDBR := tbls.ConnPool.WithQueryBuilder(selTenRows)
		entINSERTStmtA := tbls.ConnPool.WithPrepare(ctx, tbl.Insert().BuildValues())
		for i := 0; i < 9; i++ {
			entIn := new(AthleteTeam)
			assert.NoError(t, ps.FakeData(entIn), "Error at index %d", i)
			lID := dmltest.CheckLastInsertID(t, "Error: TestNewTables.AthleteTeam_Entity")(entINSERTStmtA.ExecContext(ctx, dml.Qualify("", entIn)))
			entINSERTStmtA.Reset()
			entOut := new(AthleteTeam)
			rowCount, err := selOneRowDBR.Load(ctx, entOut, lID)
			assert.NoError(t, err)
			assert.Exactly(t, uint64(1), rowCount, "IDX%d: RowCount did not match", i)
			assert.Exactly(t, entIn.TeamID, entOut.TeamID, "IDX%d: TeamID should match", lID)
			assert.ExactlyLength(t, 340, &entIn.Name, &entOut.Name, "IDX%d: Name should match", lID)
		}
		dmltest.Close(t, entINSERTStmtA)
		entCol := NewAthleteTeams()
		rowCount, err := selTenRowsDBR.Load(ctx, entCol)
		assert.NoError(t, err)
		t.Logf("Collection load rowCount: %d", rowCount)
		colInsertDBR := tbls.ConnPool.WithQueryBuilder(tbl.Insert().Replace().SetRowCount(len(entCol.Data)).BuildValues())
		lID := dmltest.CheckLastInsertID(t, "Error:  AthleteTeams ")(colInsertDBR.ExecContext(ctx, dml.Qualify("", entCol)))
		t.Logf("Last insert ID into: %d", lID)
	})
	t.Run("AthleteTeamMember_Entity", func(t *testing.T) {
		tbl := tbls.MustTable(TableNameAthleteTeamMember)
		selOneRow := tbl.Select("*").Where(
			dml.Column("id").Equal().PlaceHolder(),
		)
		selTenRows := tbl.Select("*").Where(
			dml.Column("id").LessOrEqual().Int(10),
		)
		selOneRowDBR := tbls.ConnPool.WithPrepare(ctx, selOneRow)
		defer selOneRowDBR.Close()
		selTenRowsDBR := tbls.ConnPool.WithQueryBuilder(selTenRows)
		entINSERTStmtA := tbls.ConnPool.WithPrepare(ctx, tbl.Insert().BuildValues())
		for i := 0; i < 9; i++ {
			entIn := new(AthleteTeamMember)
			assert.NoError(t, ps.FakeData(entIn), "Error at index %d", i)
			lID := dmltest.CheckLastInsertID(t, "Error: TestNewTables.AthleteTeamMember_Entity")(entINSERTStmtA.ExecContext(ctx, dml.Qualify("", entIn)))
			entINSERTStmtA.Reset()
			entOut := new(AthleteTeamMember)
			rowCount, err := selOneRowDBR.Load(ctx, entOut, lID)
			assert.NoError(t, err)
			assert.Exactly(t, uint64(1), rowCount, "IDX%d: RowCount did not match", i)
			assert.Exactly(t, entIn.ID, entOut.ID, "IDX%d: ID should match", lID)
			assert.Exactly(t, entIn.TeamID, entOut.TeamID, "IDX%d: TeamID should match", lID)
			assert.Exactly(t, entIn.AthleteID, entOut.AthleteID, "IDX%d: AthleteID should match", lID)
		}
		dmltest.Close(t, entINSERTStmtA)
		entCol := NewAthleteTeamMembers()
		rowCount, err := selTenRowsDBR.Load(ctx, entCol)
		assert.NoError(t, err)
		t.Logf("Collection load rowCount: %d", rowCount)
		colInsertDBR := tbls.ConnPool.WithQueryBuilder(tbl.Insert().Replace().SetRowCount(len(entCol.Data)).BuildValues())
		lID := dmltest.CheckLastInsertID(t, "Error:  AthleteTeamMembers ")(colInsertDBR.ExecContext(ctx, dml.Qualify("", entCol)))
		t.Logf("Last insert ID into: %d", lID)
	})
	t.Run("CustomerAddressEntity_Entity", func(t *testing.T) {
		tbl := tbls.MustTable(TableNameCustomerAddressEntity)
		selOneRow := tbl.Select("*").Where(
			dml.Column("entity_id").Equal().PlaceHolder(),
		)
		selTenRows := tbl.Select("*").Where(
			dml.Column("entity_id").LessOrEqual().Int(10),
		)
		selOneRowDBR := tbls.ConnPool.WithPrepare(ctx, selOneRow)
		defer selOneRowDBR.Close()
		selTenRowsDBR := tbls.ConnPool.WithQueryBuilder(selTenRows)
		entINSERTStmtA := tbls.ConnPool.WithPrepare(ctx, tbl.Insert().BuildValues())
		for i := 0; i < 9; i++ {
			entIn := new(CustomerAddressEntity)
			assert.NoError(t, ps.FakeData(entIn), "Error at index %d", i)
			lID := dmltest.CheckLastInsertID(t, "Error: TestNewTables.CustomerAddressEntity_Entity")(entINSERTStmtA.ExecContext(ctx, dml.Qualify("", entIn)))
			entINSERTStmtA.Reset()
			entOut := new(CustomerAddressEntity)
			rowCount, err := selOneRowDBR.Load(ctx, entOut, lID)
			assert.NoError(t, err)
			assert.Exactly(t, uint64(1), rowCount, "IDX%d: RowCount did not match", i)
			assert.Exactly(t, entIn.EntityID, entOut.EntityID, "IDX%d: EntityID should match", lID)
			assert.ExactlyLength(t, 50, &entIn.IncrementID, &entOut.IncrementID, "IDX%d: IncrementID should match", lID)
			assert.Exactly(t, entIn.ParentID, entOut.ParentID, "IDX%d: ParentID should match", lID)
			assert.Exactly(t, entIn.IsActive, entOut.IsActive, "IDX%d: IsActive should match", lID)
			assert.ExactlyLength(t, 255, &entIn.City, &entOut.City, "IDX%d: City should match", lID)
			assert.ExactlyLength(t, 255, &entIn.Company, &entOut.Company, "IDX%d: Company should match", lID)
			assert.ExactlyLength(t, 255, &entIn.CountryID, &entOut.CountryID, "IDX%d: CountryID should match", lID)
			assert.ExactlyLength(t, 255, &entIn.Firstname, &entOut.Firstname, "IDX%d: Firstname should match", lID)
			assert.ExactlyLength(t, 255, &entIn.Lastname, &entOut.Lastname, "IDX%d: Lastname should match", lID)
			assert.ExactlyLength(t, 255, &entIn.Postcode, &entOut.Postcode, "IDX%d: Postcode should match", lID)
			assert.ExactlyLength(t, 255, &entIn.Region, &entOut.Region, "IDX%d: Region should match", lID)
			assert.ExactlyLength(t, 65535, &entIn.Street, &entOut.Street, "IDX%d: Street should match", lID)
		}
		dmltest.Close(t, entINSERTStmtA)
		entCol := NewCustomerAddressEntities()
		rowCount, err := selTenRowsDBR.Load(ctx, entCol)
		assert.NoError(t, err)
		t.Logf("Collection load rowCount: %d", rowCount)
		colInsertDBR := tbls.ConnPool.WithQueryBuilder(tbl.Insert().Replace().SetRowCount(len(entCol.Data)).BuildValues())
		lID := dmltest.CheckLastInsertID(t, "Error:  CustomerAddressEntities ")(colInsertDBR.ExecContext(ctx, dml.Qualify("", entCol)))
		t.Logf("Last insert ID into: %d", lID)
	})
	t.Run("CustomerEntity_Entity", func(t *testing.T) {
		tbl := tbls.MustTable(TableNameCustomerEntity)
		selOneRow := tbl.Select("*").Where(
			dml.Column("entity_id").Equal().PlaceHolder(),
		)
		selTenRows := tbl.Select("*").Where(
			dml.Column("entity_id").LessOrEqual().Int(10),
		)
		selOneRowDBR := tbls.ConnPool.WithPrepare(ctx, selOneRow)
		defer selOneRowDBR.Close()
		selTenRowsDBR := tbls.ConnPool.WithQueryBuilder(selTenRows)
		entINSERTStmtA := tbls.ConnPool.WithPrepare(ctx, tbl.Insert().BuildValues())
		for i := 0; i < 9; i++ {
			entIn := new(CustomerEntity)
			assert.NoError(t, ps.FakeData(entIn), "Error at index %d", i)
			lID := dmltest.CheckLastInsertID(t, "Error: TestNewTables.CustomerEntity_Entity")(entINSERTStmtA.ExecContext(ctx, dml.Qualify("", entIn)))
			entINSERTStmtA.Reset()
			entOut := new(CustomerEntity)
			rowCount, err := selOneRowDBR.Load(ctx, entOut, lID)
			assert.NoError(t, err)
			assert.Exactly(t, uint64(1), rowCount, "IDX%d: RowCount did not match", i)
			assert.Exactly(t, entIn.EntityID, entOut.EntityID, "IDX%d: EntityID should match", lID)
			assert.Exactly(t, entIn.WebsiteID, entOut.WebsiteID, "IDX%d: WebsiteID should match", lID)
			assert.ExactlyLength(t, 255, &entIn.Email, &entOut.Email, "IDX%d: Email should match", lID)
			assert.Exactly(t, entIn.GroupID, entOut.GroupID, "IDX%d: GroupID should match", lID)
			assert.Exactly(t, entIn.StoreID, entOut.StoreID, "IDX%d: StoreID should match", lID)
			assert.Exactly(t, entIn.IsActive, entOut.IsActive, "IDX%d: IsActive should match", lID)
			assert.ExactlyLength(t, 255, &entIn.CreatedIn, &entOut.CreatedIn, "IDX%d: CreatedIn should match", lID)
			assert.ExactlyLength(t, 255, &entIn.Firstname, &entOut.Firstname, "IDX%d: Firstname should match", lID)
			assert.ExactlyLength(t, 255, &entIn.Lastname, &entOut.Lastname, "IDX%d: Lastname should match", lID)
			assert.ExactlyLength(t, 128, &entIn.PasswordHash, &entOut.PasswordHash, "IDX%d: PasswordHash should match", lID)
			assert.ExactlyLength(t, 128, &entIn.RpToken, &entOut.RpToken, "IDX%d: RpToken should match", lID)
			assert.Exactly(t, entIn.DefaultBilling, entOut.DefaultBilling, "IDX%d: DefaultBilling should match", lID)
			assert.Exactly(t, entIn.DefaultShipping, entOut.DefaultShipping, "IDX%d: DefaultShipping should match", lID)
			assert.Exactly(t, entIn.Gender, entOut.Gender, "IDX%d: Gender should match", lID)
		}
		dmltest.Close(t, entINSERTStmtA)
		entCol := NewCustomerEntities()
		rowCount, err := selTenRowsDBR.Load(ctx, entCol)
		assert.NoError(t, err)
		t.Logf("Collection load rowCount: %d", rowCount)
		colInsertDBR := tbls.ConnPool.WithQueryBuilder(tbl.Insert().Replace().SetRowCount(len(entCol.Data)).BuildValues())
		lID := dmltest.CheckLastInsertID(t, "Error:  CustomerEntities ")(colInsertDBR.ExecContext(ctx, dml.Qualify("", entCol)))
		t.Logf("Last insert ID into: %d", lID)
	})
	// Uncomment the next line for debugging to see all the queries.
	// t.Logf("queries: %#v", tbls.ConnPool.CachedQueries())
}
