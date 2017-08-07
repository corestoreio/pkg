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

package store_test

import (
	"testing"

	"github.com/corestoreio/csfw/storage/dbr"
	"github.com/corestoreio/csfw/store"
	"github.com/corestoreio/csfw/util/cstesting"
	"github.com/corestoreio/csfw/util/magento"
	"github.com/corestoreio/log"
)

// BenchmarkIntegration_TableStoreSlice_Native-4   	   10000	    259183 ns/op	    8501 B/op	     306 allocs/op <= no prepare; Rows()
// BenchmarkIntegration_TableStoreSlice_Native-4   	   10000	    217399 ns/op	    4704 B/op	     207 allocs/op <= prepared statement
// with 18 rows
func BenchmarkIntegration_TableStoreSlice_Native(b *testing.B) {

	dbc, mageVersion := cstesting.MustConnectDB()
	if dbc == nil {
		b.Skip("Environment DB DSN not found")
	}
	if mageVersion != magento.Version1 {
		b.Skip("Expecting a Magento1 database structure")
	}
	defer dbc.Close()

	tsr := store.NewTableStoreResource(dbc.DB)
	tsr.Table.Name = "core_store"
	defer tsr.Close()

	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		tss, err := tsr.Select()
		if err != nil {
			b.Fatalf("%+v", err)
		}
		if have, want := len(tss), 18; have != want {
			b.Errorf("Have: %v Want: %v\n%#v", have, want, tss)
		}
	}
}

// BenchmarkIntegration_TableStoreSlice_Reflection-4   	    5000	    308089 ns/op	   11390 B/op	     434 allocs/op <= without prepare
func BenchmarkIntegration_TableStoreSlice_Reflection(b *testing.B) {
	dbc, mageVersion := cstesting.MustConnectDB()
	if dbc == nil {
		b.Skip("Environment DB DSN not found")
	}
	if mageVersion != magento.Version1 {
		b.Skip("Expecting a Magento1 database structure")
	}

	sb := &dbr.Select{
		Logger:    log.BlackHole{EnableDebug: false, EnableInfo: false},
		Querier:   dbc.DB,
		Columns:   []string{"`main_table`.`store_id`", "`main_table`.`code`", "`main_table`.`website_id`", "`main_table`.`group_id`", "`main_table`.`name`", "`main_table`.`sort_order`", "`main_table`.`is_active`"},
		FromTable: dbr.MakeIdentifier("core_store", "main_table"),
	}

	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {

		var tss store.TableStoreSlice
		rows, err := sb.LoadStructs(&tss)
		if err != nil {
			b.Fatalf("%+v", err)
		}
		if have, want := rows, 18; have != want {
			b.Errorf("Have: %v Want: %v", have, want)
		}
	}
}
