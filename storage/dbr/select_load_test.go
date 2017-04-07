// Copyright 2015-2017, Cyrill @ Schumacher.fm and the CoreStore contributors
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

package dbr_test

import (
	"context"
	"database/sql"
	"testing"

	"github.com/corestoreio/csfw/storage/dbr"
)

var _ dbr.Querier = (*benchMockQuerier)(nil)

type benchMockQuerier struct{}

func (benchMockQuerier) QueryContext(ctx context.Context, query string, args ...interface{}) (*sql.Rows, error) {
	return new(sql.Rows), nil
}

// BenchmarkSelect_Rows-4   	 1000000	      2188 ns/op	    1354 B/op	      19 allocs/op old
// BenchmarkSelect_Rows-4   	 1000000	      2223 ns/op	    1386 B/op	      20 allocs/op new
func BenchmarkSelect_Rows(b *testing.B) {

	tables := []string{"eav_attribute"}
	ctx := context.TODO()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {

		sel := dbr.NewSelect("information_schema.COLUMNS").AddColumns(
			"TABLE_NAME", "COLUMN_NAME", "ORDINAL_POSITION", "COLUMN_DEFAULT", "IS_NULLABLE",
			"DATA_TYPE", "CHARACTER_MAXIMUM_LENGTH", "NUMERIC_PRECISION", "NUMERIC_SCALE",
			"COLUMN_TYPE", "COLUMN_KEY", "EXTRA", "COLUMN_COMMENT").
			Where(dbr.Condition(`TABLE_SCHEMA=DATABASE()`))
		sel.DB.Querier = benchMockQuerier{}
		if len(tables) > 0 {
			sel.Where(dbr.Condition("TABLE_NAME IN ?", dbr.ArgString(tables...)))
		}

		rows, err := sel.Rows(ctx)
		if err != nil {
			b.Fatalf("%+v", err)
		}
		if rows == nil {
			b.Fatal("Rows should not be nil")
		}
	}
}
