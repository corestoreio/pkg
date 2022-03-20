// Copyright 2015-present, Cyrill @ Schumacher.fm and the CoreStore contributors
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

package mycanal_test

import (
	"context"
	"flag"
	"testing"
	"time"

	"github.com/corestoreio/errors"
	"github.com/corestoreio/pkg/sql/ddl"
	"github.com/corestoreio/pkg/sql/dml"
	"github.com/corestoreio/pkg/sql/dmltest"
	"github.com/corestoreio/pkg/sql/mycanal"
	"github.com/corestoreio/pkg/util/assert"
)

// TODO(CyS): Add more tests

var (
	runIntegration = flag.Bool("integration", false, "Enables MySQL/MariaDB integration tests, env var CS_DSN must be set")
)

func TestIntegrationNewCanal(t *testing.T) {
	if !*runIntegration {
		t.Skip("Skipping integration tests. You can enable them with via CLI option `-integration`")
	}

	dsn := dmltest.MustGetDSN(t)

	// var bufLog bytes.Buffer
	// myLog := logw.NewLog(logw.WithDebug(&bufLog, "INTG", log.LstdFlags))

	c, err := mycanal.NewCanal(dsn, mycanal.WithMySQL(), &mycanal.Options{
		// Log:               myLog,
		IncludeTableRegex: []string{"catalog_product_entity", "^sales_order$"},
		OnClose: func(db *dml.ConnPool) (err error) {
			// return nil
			if _, err = db.DB.ExecContext(context.Background(), `DROP TABLE IF EXISTS
			catalog_product_entity_datetime, catalog_product_entity_decimal, catalog_product_entity_int,
			catalog_product_entity_text, catalog_product_entity_varchar,catalog_product_entity`); err != nil {
				return errors.WithStack(err)
			}
			_, err = db.DB.ExecContext(context.Background(), `DROP TABLE IF EXISTS
			catalog_category_entity_datetime,catalog_category_entity_decimal,catalog_category_entity_int,
			catalog_category_entity_text,catalog_category_entity_varchar,catalog_category_entity,sales_order`)
			return errors.WithStack(err)
		},
	})
	assert.NoError(t, err, "%+v", err)

	cpe := &catalogProductEvent{idx: 1001, t: t, counter: make(map[string]int)}
	c.RegisterRowsEventHandler(nil, cpe)

	soe := &salesOrderEvent{idx: 1001, t: t, counter: make(map[string]int)}
	c.RegisterRowsEventHandler([]string{"sales_order"}, soe)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()
	assert.NoError(t, c.Start(ctx), "Start error")

	dmltest.SQLDumpLoad(t, "testdata/*sql", nil)

	select {
	case <-ctx.Done():
		t.Logf("Err CTX: %+v", ctx.Err())
		err := c.Close()
		assert.NoError(t, err, "c.Close(): %+v", err)
	}
	assert.Exactly(t, 6, cpe.counter[mycanal.InsertAction], "InsertActions:\n%#v", cpe.counter) // 5+1 (1=> sales_order)
	assert.Exactly(t, 3, cpe.counter[mycanal.UpdateAction], "UpdateActions:\n%#v", cpe.counter)
	assert.Exactly(t, 1, cpe.counter[mycanal.DeleteAction], "DeleteActions:\n%#v", cpe.counter)
	assert.Exactly(t, 1, soe.counter[mycanal.InsertAction], "InsertActions:\n%#v", cpe.counter)
}

var _ mycanal.RowsEventHandler = (*catalogProductEvent)(nil)

type catalogProductEvent struct {
	idx     int
	t       *testing.T
	counter map[string]int
}

func (cpe *catalogProductEvent) Do(_ context.Context, action string, table *ddl.Table, rows [][]any) error {
	cpe.counter[action]++
	// Uncomment the following lines to see the data
	// cpe.t.Logf("%d: %q %q.%q", cpe.idx, action, table.Schema, table.Name)
	// for _, r := range rows {
	// 	cpe.t.Logf("%#v", r)
	// }

	if action == mycanal.UpdateAction && table.Name == "catalog_product_entity" {
		assert.Exactly(cpe.t, []any{int32(66), int16(9), "simple", "MH01-XL-Orange", int16(0), int16(0), "2018-04-17 21:42:21", "2018-04-17 21:42:21"}, rows[0], "A: Row0: %#v", rows[0])
		assert.Exactly(cpe.t, []any{int32(66), int16(111), "simple", "MH01-XL-CS111", int16(1), int16(0), "2018-04-17 21:42:21", "2018-04-17 21:42:21"}, rows[1], "B: Row1: %#v", rows[1])
		assert.Len(cpe.t, rows, 2)
	}
	if action == mycanal.DeleteAction && table.Name == "catalog_product_entity" {
		assert.Exactly(cpe.t, []any{int32(65), int16(9), "simple", "MH01-XL-Gray", int16(0), int16(0), "2018-04-17 21:42:21", "2018-04-17 21:42:21"}, rows[0], "C: Row0: %#v", rows[0])
		assert.Len(cpe.t, rows, 1)
	}
	if action == mycanal.UpdateAction && table.Name == "catalog_product_entity_varchar" {
		assert.Exactly(cpe.t, []any{int32(436), int16(73), int16(0), int32(44), "Didi Sport Watch"}, rows[0], "D: Row0: %#v", rows[0])
		assert.Exactly(cpe.t, []any{int32(436), int16(73), int16(0), int32(44), "Dodo Sport Watch and See"}, rows[1], "E: Row1: %#v", rows[1])
		assert.Len(cpe.t, rows, 2)
	}
	if action == mycanal.UpdateAction && table.Name == "catalog_product_entity_int" {
		assert.Exactly(cpe.t, []any{int32(212), int16(99), int16(0), int32(44), int32(4)}, rows[0], "F: Row0: %#v", rows[0])
		assert.Exactly(cpe.t, []any{int32(212), int16(99), int16(0), int32(44), nil}, rows[1], "G: Row1: %#v", rows[1])
		assert.Len(cpe.t, rows, 2)
	}

	return nil
}
func (cpe *catalogProductEvent) Complete(_ context.Context) error {
	return nil // errors.NewFatalf("[test] What is incomplete?")
}

func (cpe *catalogProductEvent) String() string {
	return "catalogProductEvent"
}

var _ mycanal.RowsEventHandler = (*salesOrderEvent)(nil)

type salesOrderEvent struct {
	idx     int
	t       *testing.T
	counter map[string]int
}

func (cpe *salesOrderEvent) Do(_ context.Context, action string, table *ddl.Table, rows [][]any) error {
	if table.Name != "sales_order" {
		// should not happen due to the special registration of this handler
		return errors.Fatal.Newf("table name %q not allowed and not expected", table.Name)
	}
	cpe.counter[action]++
	// Uncomment the following lines to see the data
	// cpe.t.Logf("%d: %q %q.%q", cpe.idx, action, table.Schema, table.Name)
	// for _, r := range rows {
	// 	cpe.t.Logf("%#v", r)
	// }

	if action == mycanal.InsertAction {
		assert.Exactly(cpe.t, "89875168d4e71e08688d6a266413162a", rows[0][4], "A: Row4: %#v", rows[0][4])
		// 	assert.Exactly(cpe.t, []any{int32(66), int16(111), "simple", "MH01-XL-CS111", int16(1), int16(0), "2018-04-17 21:42:21", "2018-04-17 21:42:21"}, rows[1], "B: Row1: %#v", rows[1])
		assert.Len(cpe.t, rows, 2)
	}

	return nil
}
func (cpe *salesOrderEvent) Complete(_ context.Context) error {
	return nil // errors.NewFatalf("[test] What is incomplete?")
}
func (cpe *salesOrderEvent) String() string {
	return "salesOrderEvent"
}
