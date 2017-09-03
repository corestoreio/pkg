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

package dml_test

import (
	"context"
	"database/sql"
	"flag"
	"testing"

	"github.com/corestoreio/csfw/sql/dml"
	"github.com/corestoreio/csfw/util/cstesting"
)

var runIntegration bool

func init() {
	flag.BoolVar(&runIntegration, "integration", false, "Enables dml integration tests")
}

// table with 2007 rows and 5 columns
// BenchmarkSelect_Integration_LoadStructs-4   	     300	   3995130 ns/op	  839604 B/op	   23915 allocs/op <- Reflection with struct tags
// BenchmarkSelect_Integration_LoadX-4         	     500	   3190194 ns/op	  752296 B/op	   21883 allocs/op <- "No Reflection"
// BenchmarkSelect_Integration_LoadGoSQLDriver-4   	 500	   2975945 ns/op	  738824 B/op	   17859 allocs/op
// BenchmarkSelect_Integration_LoadPubNative-4       500	   2826601 ns/op	  669699 B/op	   11966 allocs/op <- no database/sql

// BenchmarkSelect_Integration_Load-4   	     500	   3393616 ns/op	  752254 B/op	   21882 allocs/op <- if/else orgie
// BenchmarkSelect_Integration_Load-4   	     500	   3461720 ns/op	  752234 B/op	   21882 allocs/op <- switch

// BenchmarkSelect_Integration_LScanner-4   	 500	   3425029 ns/op	  755206 B/op	   21878 allocs/op
// BenchmarkSelect_Integration_Scanner-4   	     500	   3288291 ns/op	  784423 B/op	   23890 allocs/op <- iFace with Scan function
// BenchmarkSelect_Integration_Scanner-4   	     500	   3001319 ns/op	  784290 B/op	   23888 allocs/op Go 1.9 with new Scanner iFace
// BenchmarkSelect_Integration_Scanner-4   	    1000	   1947410 ns/op	  743693 B/op	   17876 allocs/op Go 1.9 with RowConvert type and sql.RawBytes
func BenchmarkSelect_Integration_Scanner(b *testing.B) {
	if !runIntegration {
		b.Skip("Skipped. To enable use -integration=1")
	}

	const coreConfigDataRowCount = 2007

	c := createRealSession(b)
	defer cstesting.Close(b, c)

	s := c.SelectFrom("core_config_data112").Star()
	ctx := context.TODO()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		var ccd TableCoreConfigDataSlice
		if _, err := s.Load(ctx, &ccd); err != nil {
			b.Fatalf("%+v", err)
		}
		if len(ccd.Data) != coreConfigDataRowCount {
			b.Fatal("Length mismatch")
		}
	}
}

//BenchmarkInsert_Prepared/ExecRecord-4        	    5000	    320371 ns/op	     512 B/op	      12 allocs/op
//BenchmarkInsert_Prepared/ExecArgs-4         	    5000	    310453 ns/op	     640 B/op	      14 allocs/op
//BenchmarkInsert_Prepared/ExecContext-4      	    5000	    312097 ns/op	     608 B/op	      13 allocs/op
func BenchmarkInsert_Prepared(b *testing.B) {
	if !runIntegration {
		b.Skip("Skipped. To enable use -integration=1")
	}

	truncate := func(db dml.Execer) {
		if _, err := db.ExecContext(context.TODO(), "TRUNCATE TABLE `dml_people`"); err != nil {
			b.Fatal(err)
		}
	}

	c := createRealSession(b)
	defer func() {
		truncate(c.DB)
		cstesting.Close(b, c)
	}()
	truncate(c.DB)

	stmt, err := c.InsertInto("dml_people").
		AddColumns("name", "email", "store_id", "created_at", "total_income").
		Prepare(context.TODO())
	if err != nil {
		b.Fatal(err)
	}
	defer func() {
		if err := stmt.Close(); err != nil {
			b.Fatal(err)
		}
	}()
	const totalIncome = 4.3215
	const storeID = 12345
	ctx := context.TODO()
	b.ResetTimer()

	b.Run("ExecRecord", func(b *testing.B) {
		truncate(c.DB)
		p := &dmlPerson{
			Name:        "Maria Gopher ExecRecord",
			Email:       dml.MakeNullString("maria@gopherExecRecord.go"),
			StoreID:     storeID,
			CreatedAt:   now(),
			TotalIncome: totalIncome,
		}

		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			p.TotalIncome = totalIncome * float64(i)
			res, err := stmt.WithRecords(p).ExecContext(ctx)
			if err != nil {
				b.Fatal(err)
			}
			lid, err := res.LastInsertId()
			if err != nil {
				b.Fatal(err)
			}
			if lid < 1 {
				b.Fatalf("LastInsertID was %d", lid)
			}
		}
	})

	b.Run("ExecArgs", func(b *testing.B) {
		truncate(c.DB)
		args := dml.MakeArgs(5)

		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			args = args[:0]

			res, err := stmt.WithArguments(args.
				Str("Maria Gopher ExecArgs").NullString(dml.MakeNullString("maria@gopherExecArgs.go")).
				Int64(storeID).Time(now()).Float64(totalIncome * float64(i)),
			).ExecContext(ctx)
			if err != nil {
				b.Fatal(err)
			}
			lid, err := res.LastInsertId()
			if err != nil {
				b.Fatal(err)
			}
			if lid < 1 {
				b.Fatalf("LastInsertID was %d", lid)
			}
		}
	})

	b.Run("ExecContext", func(b *testing.B) {
		truncate(c.DB)
		name := "Maria Gopher ExecContext"
		email := sql.NullString{String: "maria@gopherExecContext.go", Valid: true}

		b.ResetTimer()
		for i := 0; i < b.N; i++ {

			res, err := stmt.ExecContext(ctx, name, email, storeID, now(), totalIncome*float64(i))
			if err != nil {
				b.Fatal(err)
			}
			lid, err := res.LastInsertId()
			if err != nil {
				b.Fatal(err)
			}
			if lid < 1 {
				b.Fatalf("LastInsertID was %d", lid)
			}
		}
	})

}
