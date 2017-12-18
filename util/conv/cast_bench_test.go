package conv

import (
	"errors"
	"fmt"
	"math"
	"testing"
)

var benchmarkToString string

func benchmarkToStringF(b *testing.B, s interface{}) {
	b.ResetTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		benchmarkToString = ToString(s)
	}
	if benchmarkToString == "" {
		b.Fatal("benchmarkToString is empty :-(")
	}
}

var benchmarkToFloat64 float64

func benchmarkToFloat64F(b *testing.B, s interface{}) {
	b.ResetTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		benchmarkToFloat64 = ToFloat64(s)
	}
	if benchmarkToFloat64 < 0.001 {
		b.Fatal("benchmarkToFloat64 is empty :-(")
	}
}

var benchmarkToBool bool

func benchmarkToBoolF(b *testing.B, s interface{}) {
	b.ResetTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		benchmarkToBool = ToBool(s)
	}
	if !benchmarkToBool {
		b.Fatal("benchmarkToBool is false :-(")
	}
}

func BenchmarkToString_String(b *testing.B) {
	var val = "Hello cute little Goph3rs out there! Eat an  :-)"
	benchmarkToStringF(b, val)
}
func BenchmarkToString_Bytes(b *testing.B) {
	var val = []byte("Hello cute little Goph3rs out there! Eat an  :-)")
	benchmarkToStringF(b, val)
}
func BenchmarkToString_Float64(b *testing.B) {
	var val = math.Pi * 33
	benchmarkToStringF(b, val)
}
func BenchmarkToString_Int(b *testing.B) {
	var val = int(math.MaxInt16)
	benchmarkToStringF(b, val)
}
func BenchmarkToString_Error(b *testing.B) {
	var val = errors.New("Luke, I'm not your father.")
	benchmarkToStringF(b, val)
}

func BenchmarkToFloat64_Float64(b *testing.B) {
	var val = math.Pi
	benchmarkToFloat64F(b, val)
}

func BenchmarkToFloat64_Int64(b *testing.B) {
	var val = math.MaxInt64 - int64(math.MaxInt16)
	benchmarkToFloat64F(b, val)
}

func BenchmarkToFloat64_String(b *testing.B) {
	var val = fmt.Sprintf("%.10f", math.E)
	benchmarkToFloat64F(b, val)
}

func BenchmarkToBool_String(b *testing.B) {
	benchmarkToBoolF(b, "true")
}

func BenchmarkToBool_Interface(b *testing.B) {
	benchmarkToBoolF(b, toBool{true})
}
