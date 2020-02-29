package dmltestgenerated2

import (
	"context"
	"testing"

	"github.com/corestoreio/pkg/sql/ddl"
	"github.com/corestoreio/pkg/sql/dml"
	"github.com/corestoreio/pkg/sql/dmltest"
	"github.com/corestoreio/pkg/util/assert"
	"github.com/corestoreio/pkg/util/pseudo"
)

func TestNewDBManager_Manual(t *testing.T) {
	db := dmltest.MustConnectDB(t)
	t.Cleanup(func() {
		dmltest.Close(t, db)
		dmltest.SQLDumpLoad(t, "../testdata/test_*_tables.sql", &dmltest.SQLDumpOptions{
			SkipDBCleanup: true,
		}).Deferred()
	})

	availableEvents := []dml.EventFlag{
		dml.EventFlagBeforeUpsert, dml.EventFlagAfterUpsert,
		dml.EventFlagBeforeSelect, dml.EventFlagAfterSelect,
	}
	const (
		eventIdxEntity = iota
		eventIdxCollection
		eventIdxMax
	)
	calledEvents := [eventIdxMax][dml.EventFlagMax]int{}
	ctx := context.Background()

	opts := &DBMOption{
		TableOptions: []ddl.TableOption{ddl.WithConnPool(db)},
		InitSelectFn: func(d *dml.Select) *dml.Select {
			d.Limit(0, 100) // adds to every SELECT the LIMIT clause
			return d
		},
		InitInsertFn: nil,
	}

	for _, eventID := range availableEvents {
		eventID := eventID
		opts = opts.AddEventCoreConfiguration(eventID, func(_ context.Context, cc *CoreConfigurations, c *CoreConfiguration) error {
			if cc != nil {
				calledEvents[eventIdxCollection][eventID]++ // set to 2 to verify that it has been called
			} else if c != nil {
				calledEvents[eventIdxEntity][eventID]++ // set to 2 to verify that it has been called
			}
			return nil
		})
	}

	dbm, err := NewDBManager(ctx, opts)
	assert.NoError(t, err)

	ps := pseudo.MustNewService(0, &pseudo.Options{Lang: "de", FloatMaxDecimals: 6})

	t.Run("Entity", func(t *testing.T) {
		var eFake CoreConfiguration // e=entity => entityFake or entityLoaded
		assert.NoError(t, ps.FakeData(&eFake))

		res, err := eFake.Upsert(ctx, dbm) // INSERT ON DUPLICATE KEY
		err = dml.ExecValidateOneAffectedRow(res, err)
		assert.NoError(t, err)

		eLoaded := &CoreConfiguration{}
		err = eLoaded.Load(ctx, dbm, eFake.ConfigID)
		assert.NoError(t, err)
		assert.NotEmpty(t, eLoaded.ConfigID)
		assert.NotEmpty(t, eLoaded.Scope)
		assert.NotEmpty(t, eLoaded.ScopeID)
		assert.NotEmpty(t, eLoaded.Path)

		cq := dbm.CachedQueries()
		assert.Exactly(t, "SELECT `config_id`, `scope`, `scope_id`, `expires`, `path`, `value` FROM `core_configuration` AS `main_table` WHERE (`config_id` = ?) LIMIT 0,100",
			cq["CoreConfigurationFindByPK"])
		assert.Exactly(t, "INSERT INTO `core_configuration` (`scope`,`scope_id`,`expires`,`path`,`value`) VALUES (?,?,?,?,?) ON DUPLICATE KEY UPDATE `scope`=VALUES(`scope`), `scope_id`=VALUES(`scope_id`), `expires`=VALUES(`expires`), `path`=VALUES(`path`), `value`=VALUES(`value`)",
			cq["CoreConfigurationUpsertByPK"])

		//"SalesOrderStatusStateFindByPK": "SELECT `status`, `state`, `is_default`, `visible_on_front` FROM `sales_order_status_state` AS `main_table` WHERE (`status` = ?) AND (`state` = ?) LIMIT 0,100",
		//"ViewCustomerAutoIncrementFindByPK": "SELECT `ce_entity_id`, `email`, `firstname`, `lastname`, `city` FROM `view_customer_auto_increment` AS `main_table` LIMIT 0,100",

		for _, eventID := range availableEvents {
			assert.Exactly(t, 1, calledEvents[eventIdxEntity][eventID])
		}
	})

	t.Run("Collection", func(t *testing.T) {
		var eFake CoreConfigurations // e=entity => entityFake or entityLoaded
		assert.NoError(t, ps.FakeData(&eFake))
		eFake.Each(func(c *CoreConfiguration) {
			c.ConfigID = 0
		}) // reset configIDs

		res, err := eFake.DBUpsert(ctx, dbm)
		assert.NoError(t, err)
		lid, _ := res.LastInsertId()
		ra, _ := res.RowsAffected()
		t.Logf("lastInsertID %d RowsAffected %d", lid, ra)
	})
}
