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

package null_test

import (
	"math"
	"testing"
	"time"

	"github.com/corestoreio/pkg/storage/null"
)

func BenchmarkSQLScanner(b *testing.B) {
	var valInt64 int64
	var valFloat64 float64
	var valUint64 uint64
	var valString string
	var valTime time.Time
	b.Run("NullInt64", func(b *testing.B) {
		val := []byte(`12345678`)
		for i := 0; i < b.N; i++ {
			var nv null.Int64
			if err := nv.Scan(val); err != nil {
				b.Fatal(err)
			}
			if nv.Int64 != 12345678 {
				b.Fatalf("Have %d Want %d", nv.Int64, 12345678)
			}
			valInt64 = nv.Int64
		}
	})
	b.Run("NullFloat64", func(b *testing.B) {
		val := []byte(`-1234.5678`)
		for i := 0; i < b.N; i++ {
			var nv null.Float64
			if err := nv.Scan(val); err != nil {
				b.Fatal(err)
			}
			if nv.Float64 != -1234.5678 {
				b.Fatalf("Have %f Want %f", nv.Float64, -1234.5678)
			}
			valFloat64 = nv.Float64
		}
	})
	b.Run("NullUint64", func(b *testing.B) {
		val := []byte(`12345678910`)
		for i := 0; i < b.N; i++ {
			var nv null.Uint64
			if err := nv.Scan(val); err != nil {
				b.Fatal(err)
			}
			if nv.Uint64 != 12345678910 {
				b.Fatalf("Have %d Want %d", nv.Uint64, 12345678910)
			}
			valUint64 = nv.Uint64
		}
	})
	b.Run("NullString", func(b *testing.B) {
		const want = `12345678910`
		val := []byte(want)
		for i := 0; i < b.N; i++ {
			var nv null.String
			if err := nv.Scan(val); err != nil {
				b.Fatal(err)
			}
			if nv.Data != want {
				b.Fatalf("Have %q Want %q", nv.Data, want)
			}
			valString = nv.Data
		}
	})
	b.Run("NullTime", func(b *testing.B) {
		const want = `2006-01-02 19:04:05`
		val := []byte(want)
		for i := 0; i < b.N; i++ {
			var nv null.Time
			if err := nv.Scan(val); err != nil {
				b.Fatal(err)
			}
			if nv.Time.IsZero() {
				b.Fatalf("Time cannot be zero %s", nv.String())
			}
			valTime = nv.Time
		}
	})
	_ = valInt64
	_ = valFloat64
	_ = valUint64
	_ = valString
	_ = valTime
}

var benchmarkDecimal_String string

func BenchmarkDecimal_String(b *testing.B) {
	b.Run("123456789", func(b *testing.B) {
		d := null.Decimal{
			Precision: 123456789,
			Valid:     true,
		}
		for i := 0; i < b.N; i++ {
			benchmarkDecimal_String = d.String()
		}
	})
	b.Run("-123456789", func(b *testing.B) {
		d := null.Decimal{
			Precision: 123456789,
			Valid:     true,
			Negative:  true,
		}
		for i := 0; i < b.N; i++ {
			benchmarkDecimal_String = d.String()
		}
	})
	b.Run("12345.6789", func(b *testing.B) {
		d := null.Decimal{
			Precision: 123456789,
			Scale:     4,
			Valid:     true,
		}
		for i := 0; i < b.N; i++ {
			benchmarkDecimal_String = d.String()
		}
	})
	b.Run("-12345.6789", func(b *testing.B) {
		d := null.Decimal{
			Precision: 123456789,
			Valid:     true,
			Scale:     4,
			Negative:  true,
		}
		for i := 0; i < b.N; i++ {
			benchmarkDecimal_String = d.String()
		}
	})
	b.Run("-Scale140", func(b *testing.B) {
		d := null.Decimal{
			Valid:     true,
			Precision: math.MaxUint64,
			Scale:     140,
			Negative:  true,
		}
		for i := 0; i < b.N; i++ {
			benchmarkDecimal_String = d.String()
		}
	})
	b.Run("PrecStr Scale30", func(b *testing.B) {
		d := null.Decimal{
			Valid:        true,
			PrecisionStr: "123456789012345678912345",
			Scale:        30,
			Negative:     true,
		}
		for i := 0; i < b.N; i++ {
			benchmarkDecimal_String = d.String()
		}
	})
}

var benchmarkDecimal_MarshalBinary []byte

func BenchmarkDecimal_Binary(b *testing.B) {
	b.Run("Marshal", func(b *testing.B) {
		d := null.Decimal{
			Precision: 123456789,
			Valid:     true,
			Scale:     4,
			Negative:  true,
		}
		for i := 0; i < b.N; i++ {
			var err error
			benchmarkDecimal_MarshalBinary, err = d.MarshalBinary()
			if err != nil {
				b.Fatal(err)
			}
		}
	})

	var dUn null.Decimal
	b.Run("Unmarshal", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			if err := dUn.UnmarshalBinary(benchmarkDecimal_MarshalBinary); err != nil {
				b.Fatal(err)
			}
		}
	})
}

var benchmarkMakeDecimalBytes null.Decimal

func benchMakeDecimalBytes(data []byte) func(b *testing.B) {
	return func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			var err error
			if benchmarkMakeDecimalBytes, err = null.MakeDecimalBytes(data); err != nil {
				b.Fatal(err)
			}
		}
	}
}

func BenchmarkMakeDecimalBytes(b *testing.B) {
	tests := []string{
		`-10.550000000000000000001`,
		`-0010.651234560000000000000000`,
		`123.45`,
		`0.1234567890123456789`,
		`0`,
		`6789012345678912345678901234.123`,
	}
	for _, data := range tests {
		b.Run(string(data), benchMakeDecimalBytes([]byte(data)))
	}
}
