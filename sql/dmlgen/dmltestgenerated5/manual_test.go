package dmltestgenerated5_test

import (
	"bytes"
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/teris-io/shortid"

	"github.com/alecthomas/repr"
	"github.com/corestoreio/pkg/sql/ddl"
	"github.com/corestoreio/pkg/sql/dml"
	"github.com/corestoreio/pkg/sql/dmlgen/dmltestgenerated5"
	"github.com/corestoreio/pkg/sql/dmltest"
	"github.com/corestoreio/pkg/storage/null"
	"github.com/corestoreio/pkg/util/assert"
	"github.com/corestoreio/pkg/util/pseudo"
)

func TestManualCustomerEntities_WithAddresses(t *testing.T) {
	var buf bytes.Buffer
	defer func() { println("\n", buf.String(), "\n") }()

	db := dmltest.MustConnectDB(t, dmltest.WithSQLLog(&buf, true), dml.WithLogger(nil, shortid.MustGenerate)) // , dml.WithLogger(nil, shortid.MustGenerate)
	defer dmltest.Close(t, db)
	defer dmltest.SQLDumpLoad(t, "../testdata/testCust_01_*_tables.sql", nil).Deferred()

	ctx := context.Background()

	opts := &dmltestgenerated5.DBMOption{
		TableOptions: []ddl.TableOption{
			ddl.WithConnPool(db),
		},
	}

	dbm, err := dmltestgenerated5.NewDBManager(ctx, opts)
	assert.NoError(t, err)

	err = dbm.ConnPool.RegisterByQueryBuilder(map[string]dml.QueryBuilder{
		"custUpdateCustomerAddressEntities_Company": dbm.MustTable(dmltestgenerated5.TableNameCustomerAddressEntity).
			Update().
			AddColumns(dmltestgenerated5.Columns.CustomerAddressEntity.Company).
			Where(
				dml.Column(dmltestgenerated5.Columns.CustomerAddressEntity.EntityID).Equal().PlaceHolder(),
			),
		// above custom query looks ugly but provides compile time safety ...
		// UPDATE customer_address_entity SET company=? WHERE entity_id=?
	})
	assert.NoError(t, err)

	ps := pseudo.MustNewService(uint64(time.Now().Unix()), &pseudo.Options{
		Lang:              "de",
		MaxFloatDecimals:  6,
		MaxLenStringLimit: 41,
	})

	var lastInsertID int64

	t.Run("01 insert and update fake customer entity with addresses", func(t *testing.T) {
		var ce dmltestgenerated5.CustomerEntity
		assert.NoError(t, ps.FakeData(&ce))

		res, err := ce.Insert(ctx, dbm)
		lastInsertID = dmltest.CheckLastInsertID(t)(res, err)

		i := uint16(10)
		ce.Relations.CustomerEntityInts.Each(func(e *dmltestgenerated5.CustomerEntityInt) {
			e.AttributeID = i
			i++
		})
		ce.Relations.CustomerEntityVarchars.Each(func(e *dmltestgenerated5.CustomerEntityVarchar) {
			e.AttributeID = i
			i++
		})

		err = ce.Relations.InsertAll(ctx, dbm)
		assert.NoError(t, err)

		ce.Relations.CustomerEntityInts.Each(func(e *dmltestgenerated5.CustomerEntityInt) {
			e.Value = int32(i * 10)
			i++
		})
		ce.Relations.CustomerEntityVarchars.Each(func(e *dmltestgenerated5.CustomerEntityVarchar) {
			e.Value = null.MakeString(fmt.Sprintf("Val:%d", i))
			i++
		})
		ce.Relations.CustomerAddressEntities.Each(func(e *dmltestgenerated5.CustomerAddressEntity) {
			e.City = "m" + e.City
			e.Firstname = "m" + e.Firstname
			e.Company.Valid = false // set to null
			i++
		})

		err = ce.Relations.UpdateAll(ctx, dbm)
		assert.NoError(t, err)

		// the following code just updates the field company in customer_address_entities

		ce.Relations.CustomerAddressEntities.Each(func(e *dmltestgenerated5.CustomerAddressEntity) {
			e.Company = null.MakeString(fmt.Sprintf("Company: %d", i))
			i++
		})

		err = ce.Relations.UpdateCustomerAddressEntities(ctx, dbm, dml.DBRValidateMinAffectedRow(1), func(dbr *dml.DBR) {
			dbr.WithCacheKey("custUpdateCustomerAddressEntities_Company")
		})
		assert.NoError(t, err)

		// TODO validate that data is correct
		// repr.Println(ce)
	})
	t.Run("02 Load all", func(t *testing.T) {
		var ce dmltestgenerated5.CustomerEntity
		assert.NoError(t, ce.Load(ctx, dbm, uint32(lastInsertID)))

		err := ce.NewRelations().LoadAll(ctx, dbm)
		assert.NoError(t, err)

		repr.Println(ce)

		err = ce.Relations.DeleteAll(ctx, dbm)
		assert.NoError(t, err)
	})
}
