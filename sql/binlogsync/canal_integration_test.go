// Copyright 2015-2016, Cyrill @ Schumacher.fm and the CoreStore contributors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package binlogsync_test

import (
	"context"
	"flag"
	"os"
	"testing"
	"time"

	"github.com/corestoreio/errors"
	"github.com/corestoreio/pkg/sql/binlogsync"
	"github.com/corestoreio/pkg/sql/ddl"
	"github.com/corestoreio/pkg/sql/dml"
	"github.com/corestoreio/pkg/sql/dmltest"
	"github.com/corestoreio/pkg/util/assert"
)

// TODO(CyS): Add more tests

var (
	runIntegration = flag.Bool("integration", false, "Enables MySQL/MariaDB integration tests, env var CS_DSN must be set")
)

func TestIntegrationNewCanal_WithoutCfgSrv(t *testing.T) {
	if !*runIntegration {
		t.Skip("Skipping integration tests. You can enable them with via CLI option `-integration`")
	}

	dsn := os.Getenv(dml.EnvDSN)
	if dsn == "" {
		t.Skipf("Skipping integration test because environment variable %q not set.", dml.EnvDSN)
	}

	c, err := binlogsync.NewCanal(dsn, binlogsync.WithMySQL(), binlogsync.Options{
		IncludeTableRegex: []string{"catalog_product_entity"},
		OnClose: func(db *dml.ConnPool) (err error) {
			if _, err = db.DB.ExecContext(context.Background(), `DROP TABLE IF EXISTS 
			catalog_product_entity_datetime, catalog_product_entity_decimal, catalog_product_entity_int, 
			catalog_product_entity_text, catalog_product_entity_varchar, catalog_product_entity`); err != nil {
				return errors.WithStack(err)
			}
			_, err = db.DB.ExecContext(context.Background(), `DROP TABLE IF EXISTS
			catalog_category_entity_datetime,catalog_category_entity_decimal,catalog_category_entity_int,
			catalog_category_entity_text,catalog_category_entity_varchar,catalog_category_entity`)
			return errors.WithStack(err)
		},
	})
	assert.NoError(t, err, "%+v", err)

	cpe := &catalogProductEvent{idx: 1001, t: t, counter: make(map[string]int)}
	c.RegisterRowsEventHandler(cpe)
	// c.RegisterRowsEventHandler(catalogProductEvent{idx: 1002, t: t})

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()
	assert.NoError(t, c.Start(ctx), "Start error")

	err = dmltest.SQLDumpLoad(context.Background(), dsn, "testdata/*sql", dmltest.SQLDumpOptions{})
	assert.NoError(t, err)

	select {
	case <-ctx.Done():
		t.Logf("Err CTX: %+v", ctx.Err())
		err := c.Close()
		assert.NoError(t, err, "c.Close(): %+v", err)
	}
	assert.Exactly(t, 5, cpe.counter[binlogsync.InsertAction], "InsertActions")
	assert.Exactly(t, 3, cpe.counter[binlogsync.UpdateAction], "UpdateActions")
	assert.Exactly(t, 1, cpe.counter[binlogsync.DeleteAction], "DeleteActions")
}

type catalogProductEvent struct {
	idx     int
	t       *testing.T
	counter map[string]int
}

func (cpe *catalogProductEvent) Do(_ context.Context, action string, table ddl.Table, rows [][]interface{}) error {
	cpe.counter[action]++
	// Uncomment the following lines to see the data
	// cpe.t.Logf("%d: %q %q.%q", cpe.idx, action, table.Schema, table.Name)
	// for _, r := range rows {
	// 	cpe.t.Logf("%#v", r)
	// }

	if action == binlogsync.UpdateAction && table.Name == "catalog_product_entity" {
		assert.Exactly(cpe.t, []interface{}{int32(66), int16(9), "simple", "MH01-XL-Orange", int16(0), int16(0), "2018-04-17 21:42:21", "2018-04-17 21:42:21"}, rows[0], "A: Row0: %#v", rows[0])
		assert.Exactly(cpe.t, []interface{}{int32(66), int16(111), "simple", "MH01-XL-CS111", int16(1), int16(0), "2018-04-17 21:42:21", "2018-04-17 21:42:21"}, rows[1], "B: Row1: %#v", rows[1])
		assert.Len(cpe.t, rows, 2)
	}
	if action == binlogsync.DeleteAction && table.Name == "catalog_product_entity" {
		assert.Exactly(cpe.t, []interface{}{int32(65), int16(9), "simple", "MH01-XL-Gray", int16(0), int16(0), "2018-04-17 21:42:21", "2018-04-17 21:42:21"}, rows[0], "C: Row0: %#v", rows[0])
		assert.Len(cpe.t, rows, 1)
	}
	if action == binlogsync.UpdateAction && table.Name == "catalog_product_entity_varchar" {
		assert.Exactly(cpe.t, []interface{}{int32(436), int16(73), int16(0), int32(44), "Didi Sport Watch"}, rows[0], "D: Row0: %#v", rows[0])
		assert.Exactly(cpe.t, []interface{}{int32(436), int16(73), int16(0), int32(44), "Dodo Sport Watch and See"}, rows[1], "E: Row1: %#v", rows[1])
		assert.Len(cpe.t, rows, 2)
	}
	if action == binlogsync.UpdateAction && table.Name == "catalog_product_entity_int" {
		assert.Exactly(cpe.t, []interface{}{int32(212), int16(99), int16(0), int32(44), int32(4)}, rows[0], "F: Row0: %#v", rows[0])
		assert.Exactly(cpe.t, []interface{}{int32(212), int16(99), int16(0), int32(44), nil}, rows[1], "G: Row1: %#v", rows[1])
		assert.Len(cpe.t, rows, 2)
	}

	return nil
}
func (cpe *catalogProductEvent) Complete(_ context.Context) error {
	return nil // errors.NewFatalf("[test] What is incomplete?")
}
func (cpe *catalogProductEvent) String() string {
	return "WTF? catalogProductEvent"
}
