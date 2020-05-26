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

package dml

import (
	"bytes"
	"database/sql/driver"
	"math"
	"reflect"
	"testing"
	"time"

	"github.com/corestoreio/errors"

	"github.com/corestoreio/pkg/storage/null"
	"github.com/corestoreio/pkg/util/bufferpool"
)

func BenchmarkArgsToIFace(b *testing.B) {
	reflectIFaceContainer := make([]interface{}, 0, 25)
	finalArgs := make([]interface{}, 0, 40)
	drvVal := []driver.Valuer{null.MakeString("I'm a valid null string: See the License for the specific language governing permissions and See the License for the specific language governing permissions and See the License for the specific language governing permissions and")}
	now1 := Now.UTC()
	argUnion := make([]interface{}, 0, 30)
	b.ResetTimer()
	b.Run("reflection all types", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			reflectIFaceContainer = append(reflectIFaceContainer,
				int64(5), []int64{6, 7, 8},
				uint64(9), []uint64{10, 11, 12},
				3.14159, []float64{33.44, 55.66, 77.88},
				true, []bool{true, false, true},
				`Licensed under the Apache License, Version 2.0 (the "License");`,
				[]string{`Unless required by applicable law or agreed to in writing, software`, `Licensed under the Apache License, Version 2.0 (the "License");`},
				drvVal[0],
				nil,
				now1,
			)
			var err error
			finalArgs, err = encodePlaceholder(finalArgs, reflectIFaceContainer)
			if err != nil {
				b.Fatal(err)
			}
			// b.Fatal("%#v", finalArgs)
			reflectIFaceContainer = reflectIFaceContainer[:0]
			finalArgs = finalArgs[:0]
		}
	})
	b.Run("args all types", func(b *testing.B) {
		// two times faster than the reflection version
		finalArgs = finalArgs[:0]
		for i := 0; i < b.N; i++ {
			argUnion = append(argUnion,
				5, []int64{6, 7, 8}, uint64(9), []uint64{10, 11, 12},
				3.14159, []float64{33.44, 55.66, 77.88}, true, []bool{true, false, true},
				`Licensed under the Apache License, Version 2.0 (the "License");`,
				[]string{`Unless required by applicable law or agreed to in writing, software`, `Licensed under the Apache License, Version 2.0 (the "License");`},
				drvVal[0], nil, now1)

			finalArgs = expandInterfaces(argUnion)
			// b.Fatal("%#v", finalArgs)
			argUnion = argUnion[:0]
			finalArgs = finalArgs[:0]
		}
	})

	b.Run("reflection numbers", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			reflectIFaceContainer = append(reflectIFaceContainer,
				int64(5), []int64{6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19},
				uint64(9), []uint64{10, 11, 12, 20, 21, 22, 23, 24, 25, 26, 27, 28, 29},
				float64(3.14159), []float64{33.44, 55.66, 77.88, 11.22, math.Pi, math.E, math.Sqrt2},
				nil,
			)
			var err error
			finalArgs, err = encodePlaceholder(finalArgs, reflectIFaceContainer)
			if err != nil {
				b.Fatal(err)
			}
			// b.Fatal("%#v", finalArgs)
			reflectIFaceContainer = reflectIFaceContainer[:0]
			finalArgs = finalArgs[:0]
		}
	})
	b.Run("args numbers", func(b *testing.B) {
		// three times faster than the reflection version

		finalArgs = finalArgs[:0]
		for i := 0; i < b.N; i++ {
			argUnion = append(argUnion,
				int64(5), []int64{6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19},
				uint64(9), []uint64{10, 11, 12, 20, 21, 22, 23, 24, 25, 26, 27, 28, 29},
				3.14159, []float64{33.44, 55.66, 77.88, 11.22, math.Pi, math.E, math.Sqrt2},
				nil)

			finalArgs = expandInterfaces(argUnion)
			// b.Fatal("%#v", finalArgs)
			argUnion = argUnion[:0]
			finalArgs = finalArgs[:0]
		}
	})
}

func encodePlaceholder(args []interface{}, value interface{}) ([]interface{}, error) {
	if valuer, ok := value.(driver.Valuer); ok {
		// get driver.Valuer's data
		var err error
		value, err = valuer.Value()
		if err != nil {
			return args, err
		}
	}

	if value == nil {
		return append(args, nil), nil
	}
	v := reflect.ValueOf(value)
	switch v.Kind() {
	case reflect.String:
		return append(args, v.String()), nil
	case reflect.Bool:
		return append(args, v.Bool()), nil
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return append(args, v.Int()), nil
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return append(args, v.Uint()), nil
	case reflect.Float32, reflect.Float64:
		return append(args, v.Float()), nil
	case reflect.Struct:
		if v.Type() == reflect.TypeOf(time.Time{}) {
			return append(args, v.Interface().(time.Time)), nil
		}
	case reflect.Slice:
		if v.Type().Elem().Kind() == reflect.Uint8 {
			// []byte
			return append(args, v.Bytes()), nil
		}
		if v.Len() == 0 {
			// FIXME: support zero-length slice
			return args, errors.NotValid.Newf("invalid slice length")
		}

		for n := 0; n < v.Len(); n++ {
			var err error
			// recursion
			args, err = encodePlaceholder(args, v.Index(n).Interface())
			if err != nil {
				return args, err
			}
		}
		return args, nil
	case reflect.Ptr:
		if v.IsNil() {
			return append(args, nil), nil
		}
		return encodePlaceholder(args, v.Elem().Interface())

	}
	return args, errors.NotSupported.Newf("Type %#v not supported", value)
}

var preprocessSink string

// BenchmarkInterpolate-4   	  500000	      4013 ns/o	     174 B/o	      11 allocs/o with reflection
// BenchmarkInterpolate-4   	  500000	      3591 ns/o	     174 B/o	      11 allocs/o string
// BenchmarkInterpolate-4   	  500000	      3599 ns/o	     174 B/o	      11 allocs/o []byte
// BenchmarkInterpolate-4   	  500000	      2684 ns/op	     160 B/op	       1 allocs/op
func BenchmarkInterpolate(b *testing.B) {
	ipBuf := bufferpool.Get()
	defer bufferpool.Put(ipBuf)

	const want = `SELECT * FROM x WHERE a = 1 AND b = -2 AND c = 3 AND d = 4 AND e = 5 AND f = 6 AND g = 7 AND h = 8 AND i = 9 AND j = 10 AND k = 'Hello' AND l = 1`
	const sqlBytes = `SELECT * FROM x WHERE a = ? AND b = ? AND c = ? AND d = ? AND e = ? AND f = ? AND g = ? AND h = ? AND i = ? AND j = ? AND k = ? AND l = ?`

	args := []interface{}{1, -2, 3, 4, 5, 6, 7, 8, 9, 10, "Hello", true}
	b.ResetTimer()
	b.ReportAllocs()
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
	dk := Conditions{
		Column("name").Str("E0S 5D Mark III"),
		Column("sku").Values(),
		Column("stock").Int64(14),
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		if _, err := dk.writeOnDuplicateKey(buf, nil); err != nil {
			b.Fatalf("%+v", err)
		}
		buf.Reset()
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
				appendNullString(null.String{}, null.String{Valid: true, Data: "Hello"}).
				appendNullFloat64(null.Float64{Valid: true, Float64: math.E}, null.Float64{}).
				appendNullInt64(null.Int64{Valid: true, Int64: 987654321}, null.Int64{}).
				appendNullBool(null.Bool{}, null.Bool{Valid: true, Bool: true}, null.Bool{Valid: false, Bool: true}).
				appendNullTime(null.MakeTime(t1), null.Time{})
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
			appendNullString(null.String{}, null.String{Valid: true, Data: "Hello"}).
			appendNullFloat64(null.Float64{Valid: true, Float64: math.E}, null.Float64{}).
			appendNullInt64(null.Int64{Valid: true, Int64: 987654321}, null.Int64{}).
			appendNullBool(null.Bool{}, null.Bool{Valid: true, Bool: true}, null.Bool{Valid: false, Bool: true}).
			appendNullTime(null.MakeTime(t1), null.Time{})

		ns := []null.String{{}, {Valid: true, Data: "Hello"}}

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
				appendNullFloat64(null.Float64{Valid: true, Float64: math.E}, null.Float64{}).
				appendNullInt64(null.Int64{Valid: true, Int64: 987654321}, null.Int64{}).
				appendNullBool(null.Bool{}, null.Bool{Valid: true, Bool: true}, null.Bool{Valid: false, Bool: true}).
				appendNullTime(null.MakeTime(t1), null.Time{})
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

// BenchmarkHashSQL-4   	  617323	      1890 ns/op	      64 B/op	       4 allocs/op
func BenchmarkHashSQL(b *testing.B) {
	const sql = "WITH RECURSIVE `cte` (`n`) AS ((SELECT `name`, `d` AS `email` FROM `dml_person`) UNION ALL (SELECT `name`, `email` FROM `dml_person2` WHERE (`id` = ?))) SELECT * FROM `cte`"
	for i := 0; i < b.N; i++ {
		preprocessSink = hashSQL(sql)
	}
}
