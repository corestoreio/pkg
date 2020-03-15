package dmltestgenerated

import (
	"context"
	"sort"
	"testing"

	"github.com/corestoreio/pkg/sql/ddl"
	"github.com/corestoreio/pkg/sql/dml"
	"github.com/corestoreio/pkg/sql/dmltest"
	"github.com/corestoreio/pkg/util/assert"
	"github.com/corestoreio/pkg/util/pseudo"
)

func TestNewDBManager_Manual_Tuples(t *testing.T) {
	// var logbuf bytes.Buffer
	// defer func() { println("\n", logbuf.String(), "\n") }()
	// l := logw.NewLog(logw.WithLevel(logw.LevelDebug), logw.WithWriter(&logbuf))
	// db := dmltest.MustConnectDB(t, dml.WithLogger(l, shortid.MustGenerate))

	db := dmltest.MustConnectDB(t)
	defer dmltest.Close(t, db)
	defer dmltest.SQLDumpLoad(t, "../testdata/test_*_tables.sql", nil).Deferred()

	availableEvents := []dml.EventFlag{
		dml.EventFlagBeforeInsert, dml.EventFlagAfterInsert,
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
		TableOptionsAfter: []ddl.TableOption{ddl.WithQueryDBRCallBack(func(key string, dbr *dml.DBR) {
			switch key {
			case "SalesOrderStatusStateSelectByPK":
				dbr.Options = 0
			}
		})},
		InitSelectFn: func(d *dml.Select) *dml.Select {
			d.Limit(0, 1000) // adds to every SELECT the LIMIT clause, for testing purposes
			return d
		},
	}

	for _, eventID := range availableEvents {
		eventID := eventID
		opts = opts.AddEventSalesOrderStatusState(eventID, func(_ context.Context, cc *SalesOrderStatusStates, c *SalesOrderStatusState) error {
			if cc != nil {
				calledEvents[eventIdxCollection][eventID]++ // set to 2 to verify that it has been called
			} else if c != nil {
				calledEvents[eventIdxEntity][eventID]++ // set to 2 to verify that it has been called
			}
			return nil
		})
	}

	// used for debugging or different query styles
	shouldInterpolateFn := func(dbr *dml.DBR) {
		// dbr.Interpolate()
	}

	dbm, err := NewDBManager(ctx, opts)
	assert.NoError(t, err)

	ps := pseudo.MustNewService(0, &pseudo.Options{Lang: "de", FloatMaxDecimals: 6, MaxLenStringLimit: 41})
	t.Run("Entity", func(t *testing.T) {
		var eFake SalesOrderStatusState // e=entity => entityFake or entityLoaded
		assert.NoError(t, ps.FakeData(&eFake))

		t.Run("Insert", func(t *testing.T) {
			res, err := eFake.Insert(ctx, dbm, shouldInterpolateFn)
			assert.NoError(t, err)
			assert.NoError(t, dml.ExecValidateOneAffectedRow(res, err))
			ra, _ := res.RowsAffected()
			assert.True(t, ra > 0, "RowsAffected should be greater than 0")
		})
		t.Run("Upsert", func(t *testing.T) {
			// this test, runs the ON DUPLICATE KEY clause as the table core_config_data has a unique key.
			res, err := eFake.Upsert(ctx, dbm, shouldInterpolateFn)
			assert.NoError(t, err)
			ra, _ := res.RowsAffected()
			assert.True(t, ra == 0, "RowsAffected should be zero")
		})
		t.Run("Load", func(t *testing.T) {
			eLoaded := &SalesOrderStatusState{}
			err = eLoaded.Load(ctx, dbm, eFake.Status, eFake.State)
			assert.NoError(t, err)
			assert.NotEmpty(t, eLoaded.Status)
			assert.NotEmpty(t, eLoaded.State)
		})
	})

	t.Run("Collection", func(t *testing.T) {
		var ec SalesOrderStatusStates
		assert.NoError(t, ps.FakeData(&ec))
		t.Run("DBInsert", func(t *testing.T) {
			res, err := ec.DBInsert(ctx, dbm, shouldInterpolateFn)
			assert.NoError(t, err)
			lid, _ := res.LastInsertId()
			ra, _ := res.RowsAffected()
			t.Logf("LastInsertId(%d) RowsAffected(%d) RowsIn:%d Len:%d", lid, ra, lid+ra, len(ec.Data))
			assert.True(t, lid == 0, "LastInsertID should be zero because no previous rows")
			assert.True(t, ra > 0, "RowsAffected should be greater than 0")
		})
		t.Run("DBUpsert", func(t *testing.T) {
			// this test, runs the ON DUPLICATE KEY clause as the table core_config_data has a unique key.
			res, err := ec.DBUpsert(ctx, dbm, shouldInterpolateFn)
			assert.NoError(t, err)
			lid, _ := res.LastInsertId()
			ra, _ := res.RowsAffected()
			t.Logf("LastInsertId(%d) RowsAffected(%d)", lid, ra)
			assert.True(t, lid == 0, "LastInsertID should be zero")
			assert.True(t, ra == 0, "RowsAffected should be zero")
		})

		t.Run("validate auto increment", func(t *testing.T) {
			calls := 0
			ec.Each(func(c *SalesOrderStatusState) {
				assert.NotEmpty(t, c.Status, "status shoyld not be empty")
				assert.NotEmpty(t, c.State, "state should not be empty")
				calls++
			})
			assert.Exactly(t, calls, len(ec.Data), "Length of ec must be equal")
			t.Logf("calls %d == %d len(ec.Data)", calls, len(ec.Data))
		})
		t.Run("DBLoad All", func(t *testing.T) {
			var eca SalesOrderStatusStates
			assert.NoError(t, eca.DBLoad(ctx, dbm, nil))
			// 10 = previous rows in the DB loaded from the testdata/test_01_....sql file
			assert.Exactly(t, len(ec.Data)+11, len(eca.Data), "former collection must have the same length as the loaded one")
		})
		t.Run("DBLoad partial IDs", func(t *testing.T) {
			args := make([]SalesOrderStatusStatesDBLoadArgs, 0, len(ec.Data))
			ec.Each(func(s *SalesOrderStatusState) {
				args = append(args, SalesOrderStatusStatesDBLoadArgs{
					Status: s.Status,
					State:  s.State,
				})
			})

			var eca SalesOrderStatusStates
			assert.NoError(t, eca.DBLoad(ctx, dbm, args))
			assert.Exactly(t, len(ec.Data), len(eca.Data), "former collection must have the same length as the loaded one")
		})
		t.Run("DBDelete", func(t *testing.T) {
			t.Skip("asdadasds")
			res, err := ec.DBDelete(ctx, dbm)
			assert.NoError(t, err)
			lid, _ := res.LastInsertId()
			ra, _ := res.RowsAffected()
			assert.True(t, lid == 0, "LastInsertID should be zero")
			assert.Exactly(t, int64(len(ec.Data)), ra, "RowsAffected should be same as ec.Data length")
			t.Logf("LastInsertId(%d) RowsAffected(%d)", lid, ra)
		})
	})

	t.Run("Events and cached queries", func(t *testing.T) {
		cq := dbm.CachedQueries()
		queries := make([]string, 0, len(cq)*2)
		for k, v := range cq {
			queries = append(queries, k+"::"+v)
		}
		sort.Strings(queries)

		wantSQLQueries := []string{
			"CatalogProductIndexEAVDecimalIDXDeleteByPK::DELETE FROM `catalog_product_index_eav_decimal_idx` WHERE ((`entity_id`, `attribute_id`, `store_id`, `source_id`) IN /*TUPLES=004*/)",
			"CatalogProductIndexEAVDecimalIDXInsert::INSERT INTO `catalog_product_index_eav_decimal_idx` (`entity_id`,`attribute_id`,`store_id`,`source_id`,`value`) VALUES (?,?,?,?,?)",
			"CatalogProductIndexEAVDecimalIDXSelectByPK::SELECT `entity_id`, `attribute_id`, `store_id`, `source_id`, `value` FROM `catalog_product_index_eav_decimal_idx` AS `main_table` WHERE ((`entity_id`, `attribute_id`, `store_id`, `source_id`) = /*TUPLES=004*/) LIMIT 0,1000",
			"CatalogProductIndexEAVDecimalIDXUpdateByPK::UPDATE `catalog_product_index_eav_decimal_idx` SET `entity_id`=?, `attribute_id`=?, `store_id`=?, `source_id`=?, `value`=? WHERE ((`entity_id`, `attribute_id`, `store_id`, `source_id`) = /*TUPLES=004*/)",
			"CatalogProductIndexEAVDecimalIDXUpsertByPK::INSERT INTO `catalog_product_index_eav_decimal_idx` (`entity_id`,`attribute_id`,`store_id`,`source_id`,`value`) VALUES (?,?,?,?,?) ON DUPLICATE KEY UPDATE `entity_id`=VALUES(`entity_id`), `attribute_id`=VALUES(`attribute_id`), `store_id`=VALUES(`store_id`), `source_id`=VALUES(`source_id`), `value`=VALUES(`value`)",
			"CatalogProductIndexEAVDecimalIDXesSelectAll::SELECT `entity_id`, `attribute_id`, `store_id`, `source_id`, `value` FROM `catalog_product_index_eav_decimal_idx` AS `main_table` LIMIT 0,1000",
			"CatalogProductIndexEAVDecimalIDXesSelectByPK::SELECT `entity_id`, `attribute_id`, `store_id`, `source_id`, `value` FROM `catalog_product_index_eav_decimal_idx` AS `main_table` WHERE ((`entity_id`, `attribute_id`, `store_id`, `source_id`) IN /*TUPLES=004*/) LIMIT 0,1000",
			"CoreConfigurationDeleteByPK::DELETE FROM `core_configuration` WHERE (`config_id` IN ?)",
			"CoreConfigurationInsert::INSERT INTO `core_configuration` (`scope`,`scope_id`,`expires`,`path`,`value`) VALUES (?,?,?,?,?)",
			"CoreConfigurationSelectByPK::SELECT `config_id`, `scope`, `scope_id`, `expires`, `path`, `value` FROM `core_configuration` AS `main_table` WHERE (`config_id` = ?) LIMIT 0,1000",
			"CoreConfigurationUpdateByPK::UPDATE `core_configuration` SET `scope`=?, `scope_id`=?, `expires`=?, `path`=?, `value`=? WHERE (`config_id` = ?)",
			"CoreConfigurationUpsertByPK::INSERT INTO `core_configuration` (`scope`,`scope_id`,`expires`,`path`,`value`) VALUES (?,?,?,?,?) ON DUPLICATE KEY UPDATE `scope`=VALUES(`scope`), `scope_id`=VALUES(`scope_id`), `expires`=VALUES(`expires`), `path`=VALUES(`path`), `value`=VALUES(`value`)",
			"CoreConfigurationsSelectAll::SELECT `config_id`, `scope`, `scope_id`, `expires`, `path`, `value` FROM `core_configuration` AS `main_table` LIMIT 0,1000",
			"CoreConfigurationsSelectByPK::SELECT `config_id`, `scope`, `scope_id`, `expires`, `path`, `value` FROM `core_configuration` AS `main_table` WHERE (`config_id` IN ?) LIMIT 0,1000",
			"CustomerAddressEntitiesSelectAll::SELECT `entity_id`, `increment_id`, `parent_id`, `created_at`, `updated_at`, `is_active`, `city`, `company`, `country_id`, `fax`, `firstname`, `lastname`, `middlename`, `postcode`, `prefix`, `region`, `region_id`, `street`, `suffix`, `telephone`, `vat_id`, `vat_is_valid`, `vat_request_date`, `vat_request_id`, `vat_request_success` FROM `customer_address_entity` AS `main_table` LIMIT 0,1000",
			"CustomerAddressEntitiesSelectByPK::SELECT `entity_id`, `increment_id`, `parent_id`, `created_at`, `updated_at`, `is_active`, `city`, `company`, `country_id`, `fax`, `firstname`, `lastname`, `middlename`, `postcode`, `prefix`, `region`, `region_id`, `street`, `suffix`, `telephone`, `vat_id`, `vat_is_valid`, `vat_request_date`, `vat_request_id`, `vat_request_success` FROM `customer_address_entity` AS `main_table` WHERE (`entity_id` IN ?) LIMIT 0,1000",
			"CustomerAddressEntityDeleteByPK::DELETE FROM `customer_address_entity` WHERE (`entity_id` IN ?)",
			"CustomerAddressEntityInsert::INSERT INTO `customer_address_entity` (`increment_id`,`parent_id`,`is_active`,`city`,`company`,`country_id`,`fax`,`firstname`,`lastname`,`middlename`,`postcode`,`prefix`,`region`,`region_id`,`street`,`suffix`,`telephone`,`vat_id`,`vat_is_valid`,`vat_request_date`,`vat_request_id`,`vat_request_success`) VALUES (?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?)",
			"CustomerAddressEntitySelectByPK::SELECT `entity_id`, `increment_id`, `parent_id`, `created_at`, `updated_at`, `is_active`, `city`, `company`, `country_id`, `fax`, `firstname`, `lastname`, `middlename`, `postcode`, `prefix`, `region`, `region_id`, `street`, `suffix`, `telephone`, `vat_id`, `vat_is_valid`, `vat_request_date`, `vat_request_id`, `vat_request_success` FROM `customer_address_entity` AS `main_table` WHERE (`entity_id` = ?) LIMIT 0,1000",
			"CustomerAddressEntityUpdateByPK::UPDATE `customer_address_entity` SET `increment_id`=?, `parent_id`=?, `is_active`=?, `city`=?, `company`=?, `country_id`=?, `fax`=?, `firstname`=?, `lastname`=?, `middlename`=?, `postcode`=?, `prefix`=?, `region`=?, `region_id`=?, `street`=?, `suffix`=?, `telephone`=?, `vat_id`=?, `vat_is_valid`=?, `vat_request_date`=?, `vat_request_id`=?, `vat_request_success`=? WHERE (`entity_id` = ?)",
			"CustomerAddressEntityUpsertByPK::INSERT INTO `customer_address_entity` (`increment_id`,`parent_id`,`is_active`,`city`,`company`,`country_id`,`fax`,`firstname`,`lastname`,`middlename`,`postcode`,`prefix`,`region`,`region_id`,`street`,`suffix`,`telephone`,`vat_id`,`vat_is_valid`,`vat_request_date`,`vat_request_id`,`vat_request_success`) VALUES (?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?) ON DUPLICATE KEY UPDATE `increment_id`=VALUES(`increment_id`), `parent_id`=VALUES(`parent_id`), `is_active`=VALUES(`is_active`), `city`=VALUES(`city`), `company`=VALUES(`company`), `country_id`=VALUES(`country_id`), `fax`=VALUES(`fax`), `firstname`=VALUES(`firstname`), `lastname`=VALUES(`lastname`), `middlename`=VALUES(`middlename`), `postcode`=VALUES(`postcode`), `prefix`=VALUES(`prefix`), `region`=VALUES(`region`), `region_id`=VALUES(`region_id`), `street`=VALUES(`street`), `suffix`=VALUES(`suffix`), `telephone`=VALUES(`telephone`), `vat_id`=VALUES(`vat_id`), `vat_is_valid`=VALUES(`vat_is_valid`), `vat_request_date`=VALUES(`vat_request_date`), `vat_request_id`=VALUES(`vat_request_id`), `vat_request_success`=VALUES(`vat_request_success`)",
			"CustomerEntitiesSelectAll::SELECT `entity_id`, `website_id`, `email`, `group_id`, `increment_id`, `store_id`, `created_at`, `updated_at`, `is_active`, `disable_auto_group_change`, `created_in`, `prefix`, `firstname`, `middlename`, `lastname`, `suffix`, `dob`, `password_hash`, `rp_token`, `rp_token_created_at`, `default_billing`, `default_shipping`, `taxvat`, `confirmation`, `gender`, `failures_num`, `first_failure`, `lock_expires` FROM `customer_entity` AS `main_table` LIMIT 0,1000",
			"CustomerEntitiesSelectByPK::SELECT `entity_id`, `website_id`, `email`, `group_id`, `increment_id`, `store_id`, `created_at`, `updated_at`, `is_active`, `disable_auto_group_change`, `created_in`, `prefix`, `firstname`, `middlename`, `lastname`, `suffix`, `dob`, `password_hash`, `rp_token`, `rp_token_created_at`, `default_billing`, `default_shipping`, `taxvat`, `confirmation`, `gender`, `failures_num`, `first_failure`, `lock_expires` FROM `customer_entity` AS `main_table` WHERE (`entity_id` IN ?) LIMIT 0,1000",
			"CustomerEntityDeleteByPK::DELETE FROM `customer_entity` WHERE (`entity_id` IN ?)",
			"CustomerEntityInsert::INSERT INTO `customer_entity` (`website_id`,`email`,`group_id`,`increment_id`,`store_id`,`is_active`,`disable_auto_group_change`,`created_in`,`prefix`,`firstname`,`middlename`,`lastname`,`suffix`,`dob`,`password_hash`,`rp_token`,`rp_token_created_at`,`default_billing`,`default_shipping`,`taxvat`,`confirmation`,`gender`,`failures_num`,`first_failure`,`lock_expires`) VALUES (?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?)",
			"CustomerEntitySelectByPK::SELECT `entity_id`, `website_id`, `email`, `group_id`, `increment_id`, `store_id`, `created_at`, `updated_at`, `is_active`, `disable_auto_group_change`, `created_in`, `prefix`, `firstname`, `middlename`, `lastname`, `suffix`, `dob`, `password_hash`, `rp_token`, `rp_token_created_at`, `default_billing`, `default_shipping`, `taxvat`, `confirmation`, `gender`, `failures_num`, `first_failure`, `lock_expires` FROM `customer_entity` AS `main_table` WHERE (`entity_id` = ?) LIMIT 0,1000",
			"CustomerEntityUpdateByPK::UPDATE `customer_entity` SET `website_id`=?, `email`=?, `group_id`=?, `increment_id`=?, `store_id`=?, `is_active`=?, `disable_auto_group_change`=?, `created_in`=?, `prefix`=?, `firstname`=?, `middlename`=?, `lastname`=?, `suffix`=?, `dob`=?, `password_hash`=?, `rp_token`=?, `rp_token_created_at`=?, `default_billing`=?, `default_shipping`=?, `taxvat`=?, `confirmation`=?, `gender`=?, `failures_num`=?, `first_failure`=?, `lock_expires`=? WHERE (`entity_id` = ?)",
			"CustomerEntityUpsertByPK::INSERT INTO `customer_entity` (`website_id`,`email`,`group_id`,`increment_id`,`store_id`,`is_active`,`disable_auto_group_change`,`created_in`,`prefix`,`firstname`,`middlename`,`lastname`,`suffix`,`dob`,`password_hash`,`rp_token`,`rp_token_created_at`,`default_billing`,`default_shipping`,`taxvat`,`confirmation`,`gender`,`failures_num`,`first_failure`,`lock_expires`) VALUES (?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?) ON DUPLICATE KEY UPDATE `website_id`=VALUES(`website_id`), `email`=VALUES(`email`), `group_id`=VALUES(`group_id`), `increment_id`=VALUES(`increment_id`), `store_id`=VALUES(`store_id`), `is_active`=VALUES(`is_active`), `disable_auto_group_change`=VALUES(`disable_auto_group_change`), `created_in`=VALUES(`created_in`), `prefix`=VALUES(`prefix`), `firstname`=VALUES(`firstname`), `middlename`=VALUES(`middlename`), `lastname`=VALUES(`lastname`), `suffix`=VALUES(`suffix`), `dob`=VALUES(`dob`), `password_hash`=VALUES(`password_hash`), `rp_token`=VALUES(`rp_token`), `rp_token_created_at`=VALUES(`rp_token_created_at`), `default_billing`=VALUES(`default_billing`), `default_shipping`=VALUES(`default_shipping`), `taxvat`=VALUES(`taxvat`), `confirmation`=VALUES(`confirmation`), `gender`=VALUES(`gender`), `failures_num`=VALUES(`failures_num`), `first_failure`=VALUES(`first_failure`), `lock_expires`=VALUES(`lock_expires`)",
			"DmlgenTypesCollectionSelectAll::SELECT `id`, `col_bigint_1`, `col_bigint_2`, `col_bigint_3`, `col_bigint_4`, `col_blob`, `col_date_1`, `col_date_2`, `col_datetime_1`, `col_datetime_2`, `col_decimal_10_1`, `col_decimal_12_4`, `price_a_12_4`, `price_b_12_4`, `col_decimal_12_3`, `col_decimal_20_6`, `col_decimal_24_12`, `col_int_1`, `col_int_2`, `col_int_3`, `col_int_4`, `col_longtext_1`, `col_longtext_2`, `col_mediumblob`, `col_mediumtext_1`, `col_mediumtext_2`, `col_smallint_1`, `col_smallint_2`, `col_smallint_3`, `col_smallint_4`, `has_smallint_5`, `is_smallint_5`, `col_text`, `col_timestamp_1`, `col_timestamp_2`, `col_tinyint_1`, `col_varchar_1`, `col_varchar_100`, `col_varchar_16`, `col_char_1`, `col_char_2` FROM `dmlgen_types` AS `main_table` LIMIT 0,1000",
			"DmlgenTypesCollectionSelectByPK::SELECT `id`, `col_bigint_1`, `col_bigint_2`, `col_bigint_3`, `col_bigint_4`, `col_blob`, `col_date_1`, `col_date_2`, `col_datetime_1`, `col_datetime_2`, `col_decimal_10_1`, `col_decimal_12_4`, `price_a_12_4`, `price_b_12_4`, `col_decimal_12_3`, `col_decimal_20_6`, `col_decimal_24_12`, `col_int_1`, `col_int_2`, `col_int_3`, `col_int_4`, `col_longtext_1`, `col_longtext_2`, `col_mediumblob`, `col_mediumtext_1`, `col_mediumtext_2`, `col_smallint_1`, `col_smallint_2`, `col_smallint_3`, `col_smallint_4`, `has_smallint_5`, `is_smallint_5`, `col_text`, `col_timestamp_1`, `col_timestamp_2`, `col_tinyint_1`, `col_varchar_1`, `col_varchar_100`, `col_varchar_16`, `col_char_1`, `col_char_2` FROM `dmlgen_types` AS `main_table` WHERE (`id` IN ?) LIMIT 0,1000",
			"DmlgenTypesDeleteByPK::DELETE FROM `dmlgen_types` WHERE (`id` IN ?)",
			"DmlgenTypesInsert::INSERT INTO `dmlgen_types` (`col_bigint_1`,`col_bigint_2`,`col_bigint_3`,`col_bigint_4`,`col_blob`,`col_date_1`,`col_date_2`,`col_datetime_1`,`col_datetime_2`,`col_decimal_10_1`,`col_decimal_12_4`,`price_a_12_4`,`price_b_12_4`,`col_decimal_12_3`,`col_decimal_20_6`,`col_decimal_24_12`,`col_int_1`,`col_int_2`,`col_int_3`,`col_int_4`,`col_longtext_1`,`col_longtext_2`,`col_mediumblob`,`col_mediumtext_1`,`col_mediumtext_2`,`col_smallint_1`,`col_smallint_2`,`col_smallint_3`,`col_smallint_4`,`has_smallint_5`,`is_smallint_5`,`col_text`,`col_timestamp_2`,`col_tinyint_1`,`col_varchar_1`,`col_varchar_100`,`col_varchar_16`,`col_char_1`,`col_char_2`) VALUES (?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?)",
			"DmlgenTypesSelectByPK::SELECT `id`, `col_bigint_1`, `col_bigint_2`, `col_bigint_3`, `col_bigint_4`, `col_blob`, `col_date_1`, `col_date_2`, `col_datetime_1`, `col_datetime_2`, `col_decimal_10_1`, `col_decimal_12_4`, `price_a_12_4`, `price_b_12_4`, `col_decimal_12_3`, `col_decimal_20_6`, `col_decimal_24_12`, `col_int_1`, `col_int_2`, `col_int_3`, `col_int_4`, `col_longtext_1`, `col_longtext_2`, `col_mediumblob`, `col_mediumtext_1`, `col_mediumtext_2`, `col_smallint_1`, `col_smallint_2`, `col_smallint_3`, `col_smallint_4`, `has_smallint_5`, `is_smallint_5`, `col_text`, `col_timestamp_1`, `col_timestamp_2`, `col_tinyint_1`, `col_varchar_1`, `col_varchar_100`, `col_varchar_16`, `col_char_1`, `col_char_2` FROM `dmlgen_types` AS `main_table` WHERE (`id` = ?) LIMIT 0,1000",
			"DmlgenTypesUpdateByPK::UPDATE `dmlgen_types` SET `col_bigint_1`=?, `col_bigint_2`=?, `col_bigint_3`=?, `col_bigint_4`=?, `col_blob`=?, `col_date_1`=?, `col_date_2`=?, `col_datetime_1`=?, `col_datetime_2`=?, `col_decimal_10_1`=?, `col_decimal_12_4`=?, `price_a_12_4`=?, `price_b_12_4`=?, `col_decimal_12_3`=?, `col_decimal_20_6`=?, `col_decimal_24_12`=?, `col_int_1`=?, `col_int_2`=?, `col_int_3`=?, `col_int_4`=?, `col_longtext_1`=?, `col_longtext_2`=?, `col_mediumblob`=?, `col_mediumtext_1`=?, `col_mediumtext_2`=?, `col_smallint_1`=?, `col_smallint_2`=?, `col_smallint_3`=?, `col_smallint_4`=?, `has_smallint_5`=?, `is_smallint_5`=?, `col_text`=?, `col_timestamp_2`=?, `col_tinyint_1`=?, `col_varchar_1`=?, `col_varchar_100`=?, `col_varchar_16`=?, `col_char_1`=?, `col_char_2`=? WHERE (`id` = ?)",
			"DmlgenTypesUpsertByPK::INSERT INTO `dmlgen_types` (`col_bigint_1`,`col_bigint_2`,`col_bigint_3`,`col_bigint_4`,`col_blob`,`col_date_1`,`col_date_2`,`col_datetime_1`,`col_datetime_2`,`col_decimal_10_1`,`col_decimal_12_4`,`price_a_12_4`,`price_b_12_4`,`col_decimal_12_3`,`col_decimal_20_6`,`col_decimal_24_12`,`col_int_1`,`col_int_2`,`col_int_3`,`col_int_4`,`col_longtext_1`,`col_longtext_2`,`col_mediumblob`,`col_mediumtext_1`,`col_mediumtext_2`,`col_smallint_1`,`col_smallint_2`,`col_smallint_3`,`col_smallint_4`,`has_smallint_5`,`is_smallint_5`,`col_text`,`col_timestamp_2`,`col_tinyint_1`,`col_varchar_1`,`col_varchar_100`,`col_varchar_16`,`col_char_1`,`col_char_2`) VALUES (?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?) ON DUPLICATE KEY UPDATE `col_bigint_1`=VALUES(`col_bigint_1`), `col_bigint_2`=VALUES(`col_bigint_2`), `col_bigint_3`=VALUES(`col_bigint_3`), `col_bigint_4`=VALUES(`col_bigint_4`), `col_blob`=VALUES(`col_blob`), `col_date_1`=VALUES(`col_date_1`), `col_date_2`=VALUES(`col_date_2`), `col_datetime_1`=VALUES(`col_datetime_1`), `col_datetime_2`=VALUES(`col_datetime_2`), `col_decimal_10_1`=VALUES(`col_decimal_10_1`), `col_decimal_12_4`=VALUES(`col_decimal_12_4`), `price_a_12_4`=VALUES(`price_a_12_4`), `price_b_12_4`=VALUES(`price_b_12_4`), `col_decimal_12_3`=VALUES(`col_decimal_12_3`), `col_decimal_20_6`=VALUES(`col_decimal_20_6`), `col_decimal_24_12`=VALUES(`col_decimal_24_12`), `col_int_1`=VALUES(`col_int_1`), `col_int_2`=VALUES(`col_int_2`), `col_int_3`=VALUES(`col_int_3`), `col_int_4`=VALUES(`col_int_4`), `col_longtext_1`=VALUES(`col_longtext_1`), `col_longtext_2`=VALUES(`col_longtext_2`), `col_mediumblob`=VALUES(`col_mediumblob`), `col_mediumtext_1`=VALUES(`col_mediumtext_1`), `col_mediumtext_2`=VALUES(`col_mediumtext_2`), `col_smallint_1`=VALUES(`col_smallint_1`), `col_smallint_2`=VALUES(`col_smallint_2`), `col_smallint_3`=VALUES(`col_smallint_3`), `col_smallint_4`=VALUES(`col_smallint_4`), `has_smallint_5`=VALUES(`has_smallint_5`), `is_smallint_5`=VALUES(`is_smallint_5`), `col_text`=VALUES(`col_text`), `col_timestamp_2`=VALUES(`col_timestamp_2`), `col_tinyint_1`=VALUES(`col_tinyint_1`), `col_varchar_1`=VALUES(`col_varchar_1`), `col_varchar_100`=VALUES(`col_varchar_100`), `col_varchar_16`=VALUES(`col_varchar_16`), `col_char_1`=VALUES(`col_char_1`), `col_char_2`=VALUES(`col_char_2`)",
			"SalesOrderStatusStateDeleteByPK::DELETE FROM `sales_order_status_state` WHERE ((`status`, `state`) IN /*TUPLES=002*/)",
			"SalesOrderStatusStateInsert::INSERT INTO `sales_order_status_state` (`status`,`state`,`is_default`,`visible_on_front`) VALUES (?,?,?,?)",
			"SalesOrderStatusStateSelectByPK::SELECT `status`, `state`, `is_default`, `visible_on_front` FROM `sales_order_status_state` AS `main_table` WHERE ((`status`, `state`) = /*TUPLES=002*/) LIMIT 0,1000",
			"SalesOrderStatusStateUpdateByPK::UPDATE `sales_order_status_state` SET `status`=?, `state`=?, `is_default`=?, `visible_on_front`=? WHERE ((`status`, `state`) = /*TUPLES=002*/)",
			"SalesOrderStatusStateUpsertByPK::INSERT INTO `sales_order_status_state` (`status`,`state`,`is_default`,`visible_on_front`) VALUES (?,?,?,?) ON DUPLICATE KEY UPDATE `status`=VALUES(`status`), `state`=VALUES(`state`), `is_default`=VALUES(`is_default`), `visible_on_front`=VALUES(`visible_on_front`)",
			"SalesOrderStatusStatesSelectAll::SELECT `status`, `state`, `is_default`, `visible_on_front` FROM `sales_order_status_state` AS `main_table` LIMIT 0,1000",
			"SalesOrderStatusStatesSelectByPK::SELECT `status`, `state`, `is_default`, `visible_on_front` FROM `sales_order_status_state` AS `main_table` WHERE ((`status`, `state`) IN /*TUPLES=002*/) LIMIT 0,1000",
			"ViewCustomerAutoIncrementSelectByPK::SELECT `ce_entity_id`, `email`, `firstname`, `lastname`, `city` FROM `view_customer_auto_increment` AS `main_table` WHERE (`ce_entity_id` = ?) LIMIT 0,1000",
			"ViewCustomerAutoIncrementsSelectAll::SELECT `ce_entity_id`, `email`, `firstname`, `lastname`, `city` FROM `view_customer_auto_increment` AS `main_table` LIMIT 0,1000",
			"ViewCustomerAutoIncrementsSelectByPK::SELECT `ce_entity_id`, `email`, `firstname`, `lastname`, `city` FROM `view_customer_auto_increment` AS `main_table` WHERE (`ce_entity_id` IN ?) LIMIT 0,1000",
		}

		assert.Exactly(t, wantSQLQueries, queries)

		for _, eventID := range availableEvents {
			assert.Exactly(t, 1, calledEvents[eventIdxEntity][eventID])
		}
	})
}
