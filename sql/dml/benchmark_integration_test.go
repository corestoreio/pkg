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

// +build integration

package dml_test

import (
	"context"
	"math/rand"
	"testing"
	"time"

	"github.com/corestoreio/pkg/sql/dml"
	"github.com/corestoreio/pkg/sql/dmltest"
	"github.com/corestoreio/pkg/storage/null"
)

// table with 2007 rows and 5 columns

// MacBook Pro (13-inch, 2017, Two Thunderbolt 3 ports) 2.5 GHz Intel Core i7
// go version devel +5fae09b738 Tue Jan 15 23:30:58 2019 +0000 darwin/amd64
// BenchmarkSelectRows2007/Query-4         	    				1000	   2067316 ns/op	  742986 B/op	   17222 allocs/op
// BenchmarkSelectRows2007/Prepared,noSliceReuse-4         	    1000	   1912768 ns/op	  629676 B/op	   15738 allocs/op
// BenchmarkSelectRows2007/Prepared,SliceReuse-4           	    1000	   1921589 ns/op	  570973 B/op	   15723 allocs/op
func BenchmarkSelectRows2007(b *testing.B) {
	const coreConfigDataRowCount = 2007
	c := createRealSession(b)
	defer dmltest.Close(b, c)

	if err := c.RegisterByQueryBuilder(map[string]dml.QueryBuilder{
		"selectStar112": dml.NewSelect("*").From("core_config_data112"),
	}); err != nil {
		b.Fatal(err)
	}

	b.ResetTimer()
	b.Run("Query", func(b *testing.B) {
		s := c.WithCacheKey("selectStar112")
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
	})

	b.Run("Prepared,noSliceReuse", func(b *testing.B) {
		ctx := context.Background()
		stmt := c.WithPrepareCacheKey(ctx, "selectStar112")
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			var ccd TableCoreConfigDataSlice
			if _, err := stmt.Load(ctx, &ccd); err != nil {
				b.Fatalf("%+v", err)
			}
			if len(ccd.Data) != coreConfigDataRowCount {
				b.Fatal("Length mismatch")
			}
		}
	})
	b.Run("Prepared,SliceReuse", func(b *testing.B) {
		ctx := context.Background()
		stmt := c.WithPrepareCacheKey(ctx, "selectStar112")

		ccd := &TableCoreConfigDataSlice{
			Data: make([]*TableCoreConfigData, 0, coreConfigDataRowCount),
		}
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			if _, err := stmt.Load(ctx, ccd); err != nil {
				b.Fatalf("%+v", err)
			}
			if len(ccd.Data) != coreConfigDataRowCount {
				b.Fatal("Length mismatch")
			}
			ccd.Data = ccd.Data[:0]
		}
	})
}

// BenchmarkInsert_Prepared/ExecRecord-4       	    5000	    320371 ns/op	     512 B/op	      12 allocs/op
// BenchmarkInsert_Prepared/ExecRecord-4       	    2880	    558748 ns/op	    1000 B/op	      18 allocs/op

// BenchmarkInsert_Prepared/ExecArgs-4         	    5000	    310453 ns/op	     640 B/op	      14 allocs/op
// BenchmarkInsert_Prepared/ExecArgs-4           	  2374	    547961 ns/op	    1208 B/op	      20 allocs/op

// BenchmarkInsert_Prepared/ExecContext-4      	    5000	    312097 ns/op	     608 B/op	      13 allocs/op
// BenchmarkInsert_Prepared/ExecContext-4        	  2836	    547065 ns/op	    1233 B/op	      21 allocs/op
func BenchmarkInsert_Prepared(b *testing.B) {
	truncate := func(db dml.Execer) {
		if _, err := db.ExecContext(context.TODO(), "TRUNCATE TABLE `dml_people`"); err != nil {
			b.Fatal(err)
		}
	}

	c := createRealSession(b)
	defer func() {
		truncate(c.DB)
		dmltest.Close(b, c)
	}()
	truncate(c.DB)

	stmt := c.WithPrepare(context.TODO(), dml.NewInsert("dml_people").
		AddColumns("name", "email", "store_id", "created_at", "total_income").BuildValues())

	defer dmltest.Close(b, stmt)

	const totalIncome = 4.3215
	const storeID = 12345
	ctx := context.TODO()
	b.ResetTimer()

	b.Run("ExecRecord", func(b *testing.B) {
		truncate(c.DB)
		p := &dmlPerson{
			Name:        "Maria Gopher ExecRecord",
			Email:       null.MakeString("maria@gopherExecRecord.go"),
			StoreID:     storeID,
			CreatedAt:   now(),
			TotalIncome: totalIncome,
		}

		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			p.TotalIncome = totalIncome * float64(i)
			res, err := stmt.ExecContext(ctx, dml.Qualify("", p)) // TODO verify how the DB table looks like
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

		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			res, err := stmt.ExecContext(ctx, "Maria Gopher ExecArgs",
				null.MakeString("maria@gopherExecArgs.go"), storeID, now(), totalIncome*float64(i),
			)
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
			stmt.Reset()
		}
	})

	b.Run("ExecContext", func(b *testing.B) { // TODO rewrite this in many different ways.
		truncate(c.DB)
		name := "Maria Gopher ExecContext"
		email := null.String{Data: "maria@gopherExecContext.go", Valid: true}
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

// https://github.com/jackc/go_db_bench/blob/master/bench_test.go#L542
// https://gist.github.com/jackc/4996e8648a0c59839bff644f49d6e434#file-results-txt-L15
func BenchmarkJackC_GoDBBench(b *testing.B) {
	const maxSelectID = 24
	c := createRealSession(b)
	defer dmltest.Close(b, c)

	// prepared statement:
	// select id, first_name, last_name, sex, birth_date, weight, height, update_time
	// from dml_fake_person where id between ? and ? + 24
	stmt := c.WithPrepare(context.Background(), dml.NewSelect("id", "first_name", "last_name", "sex", "birth_date", "weight", "height", "update_time").
		From("dml_fake_person").
		Where(dml.Column("id").Between().PlaceHolder()))

	defer dmltest.Close(b, stmt)

	randPersonIDs := shuffledInts(10000)

	b.ResetTimer()

	b.Run("SelectMultipleRowsCollect DBR", func(b *testing.B) {
		ctx := context.Background()
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			id := randPersonIDs[i%len(randPersonIDs)]
			var fp fakePersons
			if _, err := stmt.Load(ctx, &fp, id, id+maxSelectID); err != nil {
				b.Fatalf("%+v", err)
			}
			for i := range fp.Data {
				checkPersonWasFilled(b, fp.Data[i])
			}
		}
	})
	b.Run("SelectMultipleRowsCollect expandInterfaces", func(b *testing.B) {
		ctx := context.Background()
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			id := randPersonIDs[i%len(randPersonIDs)]
			var fp fakePersons
			if _, err := stmt.Load(ctx, &fp, id, id+maxSelectID); err != nil {
				b.Fatalf("%+v", err)
			}
			for i := range fp.Data {
				checkPersonWasFilled(b, fp.Data[i])
			}
			stmt.Reset()
		}
	})

	b.Run("SelectMultipleRowsEntity DBR", func(b *testing.B) {
		ctx := context.Background()
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			id := randPersonIDs[i%len(randPersonIDs)]
			var fp fakePerson
			if _, err := stmt.Load(ctx, &fp, id, id+maxSelectID); err != nil {
				b.Fatalf("%+v", err)
			}
			checkPersonWasFilled(b, fp)
			stmt.Reset()
		}
	})
	b.Run("SelectMultipleRowsEntity Interface", func(b *testing.B) {
		ctx := context.Background()
		var args [2]any
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			id := randPersonIDs[i%len(randPersonIDs)]
			args[0] = id
			args[1] = id + maxSelectID
			argss := args[:]
			var fp fakePerson
			if _, err := stmt.Load(ctx, &fp, argss...); err != nil {
				b.Fatalf("%+v", err)
			}
			checkPersonWasFilled(b, fp)
		}
	})
}

func shuffledInts(size int) []int {
	randPersonIDs := make([]int, size)
	for i := range randPersonIDs {
		randPersonIDs[i] = i
	}

	vals := randPersonIDs
	r := rand.New(rand.NewSource(time.Now().Unix()))
	for len(vals) > 0 {
		n := len(vals)
		randIndex := r.Intn(n)
		vals[n-1], vals[randIndex] = vals[randIndex], vals[n-1]
		vals = vals[:n-1]
	}
	return randPersonIDs
}

func checkPersonWasFilled(b *testing.B, p fakePerson) {
	if p.ID == 0 {
		b.Fatal("id was 0")
	}
	if len(p.FirstName) == 0 {
		b.Fatal("FirstName was empty")
	}
	if len(p.LastName) == 0 {
		b.Fatal("LastName was empty")
	}
	if len(p.Sex) == 0 {
		b.Fatal("Sex was empty")
	}
	var zeroTime time.Time
	if p.BirthDate == zeroTime {
		b.Fatal("BirthDate was zero time")
	}
	if p.Weight == 0 {
		b.Fatal("Weight was 0")
	}
	if p.Height == 0 {
		b.Fatal("Height was 0")
	}
	if p.UpdateTime == zeroTime {
		b.Fatal("UpdateTime was zero time")
	}
}
