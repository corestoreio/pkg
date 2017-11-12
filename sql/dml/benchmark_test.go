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

package dml

import (
	"bytes"
	"database/sql"
	"math"
	"testing"
	"time"

	"github.com/corestoreio/pkg/util/bufferpool"
)

var preprocessSink string

// BenchmarkInterpolate-4   	  500000	      4013 ns/o	     174 B/o	      11 allocs/o with reflection
// BenchmarkInterpolate-4   	  500000	      3591 ns/o	     174 B/o	      11 allocs/o string
// BenchmarkInterpolate-4   	  500000	      3599 ns/o	     174 B/o	      11 allocs/o []byte
// BenchmarkInterpolate-4   	  500000	      2684 ns/op	     160 B/op	       1 allocs/op
func BenchmarkInterpolate(b *testing.B) {
	ipBuf := bufferpool.Get()
	defer bufferpool.Put(ipBuf)

	const want = `SELECT * FROM x WHERE a = 1 AND b = -2 AND c = 3 AND d = 4 AND e = 5 AND f = 6 AND g = 7 AND h = 8 AND i = 9 AND j = 10 AND k = 'Hello' AND l = 1`
	var sqlBytes = []byte("SELECT * FROM x WHERE a = ? AND b = ? AND c = ? AND d = ? AND e = ? AND f = ? AND g = ? AND h = ? AND i = ? AND j = ? AND k = ? AND l = ?")
	args := MakeArgs(3).
		Int64s(1, -2, 3, 4, 5, 6, 7, 8, 9, 10).
		String("Hello").
		Bool(true)

	for i := 0; i < b.N; i++ {
		if err := writeInterpolate(ipBuf, sqlBytes, args); err != nil {
			b.Fatal(err)
		}
		preprocessSink = ipBuf.String()
		ipBuf.Reset()
	}
	if preprocessSink != want {
		b.Fatalf("Have: %v Want: %v", ipBuf.String(), want)
	}
}

var benchmarkIsValidIdentifier int8

// BenchmarkIsValidIdentifier-4   	20000000	        92.0 ns/o	       0 B/o	       0 allocs/o
// BenchmarkIsValidIdentifier-4   	 5000000	       280 ns/o	       0 B/o	       0 allocs/o
func BenchmarkIsValidIdentifier(b *testing.B) {
	for i := 0; i < b.N; i++ {
		benchmarkIsValidIdentifier = isValidIdentifier(`store_owner.catalog_product_entity_varchar`)
	}
	if benchmarkIsValidIdentifier != 0 {
		b.Fatalf("Should be zero but got %d", benchmarkIsValidIdentifier)
	}
}

func BenchmarkQuoteAlias(b *testing.B) {
	const want = "`e`.`price` AS `final_price`"

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		if have := Quoter.NameAlias("e.price", "final_price"); have != want {
			b.Fatalf("Have %s\nWant %s\n", have, want)
		}
	}
}

// BenchmarkConditions_writeOnDuplicateKey-4   	 5000000	       337 ns/o	       0 B/o	       0 allocs/o
func BenchmarkConditions_writeOnDuplicateKey(b *testing.B) {
	buf := new(bytes.Buffer)
	args := MakeArgs(3)
	dk := Conditions{
		Column("name").Str("E0S 5D Mark III"),
		Column("sku").Values(),
		Column("stock").Int64(14),
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		if err := dk.writeOnDuplicateKey(buf); err != nil {
			b.Fatalf("%+v", err)
		}
		var err error
		if args, _, err = dk.appendArgs(args, appendArgsDUPKEY); err != nil {
			b.Fatalf("%+v", err)
		}
		buf.Reset()
		args = args[:0]
	}
}

var benchmarkDialectEscapeTimeBuf = new(bytes.Buffer)

func BenchmarkDialectEscapeTime(b *testing.B) {
	date := now()
	for i := 0; i < b.N; i++ {
		dialect.EscapeTime(benchmarkDialectEscapeTimeBuf, date)
		benchmarkDialectEscapeTimeBuf.Reset()
	}
}

var benchmarkArgEnc argEncoded

func BenchmarkArgumentEncoding(b *testing.B) {

	b.Run("all types without warm up", func(b *testing.B) {
		t1 := now()
		t2 := now().Add(time.Minute * 2)

		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			benchmarkArgEnc = makeArgBytes().
				appendInt(3).
				appendInts(4, 5, 6).
				appendInt64(30).
				appendInt64s(40, 50, 60).
				appendUint64(math.MaxUint32).
				appendUint64s(800, 900).
				appendFloat64(math.MaxFloat32).
				appendFloat64s(80.5490, math.Pi).
				appendString("Finally, how will we ship and deliver Go 2?").
				appendStrings("Finally, how will we fly and deliver Go 1?", "Finally, how will we run and deliver Go 3?", "Finally, how will we walk and deliver Go 3?").
				appendBool(true).
				appendBool(false).
				appendBools(false, true, true, false, true).
				appendTime(t1).
				appendTimes(t1, t2, t1).
				appendNullString(sql.NullString{}, sql.NullString{Valid: true, String: "Hello"}).
				appendNullFloat64(sql.NullFloat64{Valid: true, Float64: math.E}, sql.NullFloat64{}).
				appendNullInt64(sql.NullInt64{Valid: true, Int64: 987654321}, sql.NullInt64{}).
				appendNullBool(sql.NullBool{}, sql.NullBool{Valid: true, Bool: true}, sql.NullBool{Valid: false, Bool: true}).
				appendNullTime(NullTime{Valid: true, Time: t1}, NullTime{})
		}
	})

	b.Run("all types with warm up", func(b *testing.B) {

		t1 := now()
		t2 := now().Add(time.Minute * 2)

		benchmarkArgEnc = makeArgBytes().
			appendInt(3).
			appendInts(4, 5, 6).
			appendInt64(30).
			appendInt64s(40, 50, 60).
			appendUint64(math.MaxUint32).
			appendUint64s(800, 900).
			appendFloat64(math.MaxFloat32).
			appendFloat64s(80.5490, math.Pi).
			appendString("Finally, how will we ship and deliver Go 2?").
			appendStrings("Finally, how will we fly and deliver Go 1?", "Finally, how will we run and deliver Go 3?", "Finally, how will we walk and deliver Go 3?").
			appendBool(true).
			appendBool(false).
			appendBools(false, true, true, false, true).
			appendTime(t1).
			appendTimes(t1, t2, t1).
			appendNullString(sql.NullString{}, sql.NullString{Valid: true, String: "Hello"}).
			appendNullFloat64(sql.NullFloat64{Valid: true, Float64: math.E}, sql.NullFloat64{}).
			appendNullInt64(sql.NullInt64{Valid: true, Int64: 987654321}, sql.NullInt64{}).
			appendNullBool(sql.NullBool{}, sql.NullBool{Valid: true, Bool: true}, sql.NullBool{Valid: false, Bool: true}).
			appendNullTime(NullTime{Valid: true, Time: t1}, NullTime{})

		ns := []sql.NullString{{}, {Valid: true, String: "Hello"}}

		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			benchmarkArgEnc.
				reset().
				appendInt(13).
				appendInts(14, 15, 16).
				appendInt64(130).
				appendInt64s(140, 150, 160).
				appendUint64(math.MaxUint16).
				appendUint64s(1800, 1900).
				appendFloat64(math.MaxFloat32).
				appendFloat64s(84.5490, math.Pi).
				appendString("F1nally, how will we ship and deliver Go 1?").
				appendStrings("F1nally, how will we fly and deliver Go 2?", "Finally, how will we run and deliver Go 3?", "Finally, how will we walk and deliver Go 4?").
				appendBool(false).
				appendBool(true).
				appendBools(false, true, true, false, true).
				appendTime(t1).
				appendTimes(t1, t2, t1).
				appendNullString(ns...).
				appendNullFloat64(sql.NullFloat64{Valid: true, Float64: math.E}, sql.NullFloat64{}).
				appendNullInt64(sql.NullInt64{Valid: true, Int64: 987654321}, sql.NullInt64{}).
				appendNullBool(sql.NullBool{}, sql.NullBool{Valid: true, Bool: true}, sql.NullBool{Valid: false, Bool: true}).
				appendNullTime(NullTime{Valid: true, Time: t1}, NullTime{})
			// b.Fatal(benchmarkArgEnc.DebugBytes())
		}
	})
	b.Run("number slices without warm up", func(b *testing.B) {
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			benchmarkArgEnc = makeArgBytes().
				appendInt(3).
				appendInts(4, 5, 6).
				appendInt64(30).
				appendInt64s(40, 50, 60).
				appendUint64(math.MaxUint32).
				appendUint64s(800, 900).
				appendFloat64(math.MaxFloat32).
				appendFloat64s(80.5490, math.Pi)
		}
	})
	b.Run("number slices with warm up", func(b *testing.B) {
		benchmarkArgEnc = makeArgBytes().
			appendInt(3).
			appendInts(4, 5, 6).
			appendInt64(30).
			appendInt64s(40, 50, 60).
			appendUint64(math.MaxUint32).
			appendUint64s(800, 900).
			appendFloat64(math.MaxFloat32).
			appendFloat64s(80.5490, math.Pi)

		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			benchmarkArgEnc = benchmarkArgEnc.
				reset().
				appendInt(13).
				appendInts(14, 15, 16).
				appendInt64(130).
				appendInt64s(140, 150, 160).
				appendUint64(math.MaxUint32).
				appendUint64s(1800, 1900).
				appendFloat64(math.MaxFloat32).
				appendFloat64s(180.5490, math.Pi)
		}
	})

	b.Run("numbers without warm up", func(b *testing.B) {
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			benchmarkArgEnc = makeArgBytes().
				appendInt(3).
				appendInt64(30).
				appendUint64(math.MaxUint32).
				appendFloat64(math.MaxFloat32)
		}
	})
	b.Run("numbers with warm up", func(b *testing.B) {
		benchmarkArgEnc = makeArgBytes().
			appendInt(3).
			appendInt64(30).
			appendUint64(math.MaxUint32).
			appendFloat64(math.MaxFloat32)

		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			benchmarkArgEnc = benchmarkArgEnc.
				reset().
				appendInt(9).
				appendInt64(130).
				appendUint64(math.MaxUint64).
				appendFloat64(math.MaxFloat64)
		}
	})
}
