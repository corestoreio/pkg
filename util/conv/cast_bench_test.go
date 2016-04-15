package conv_test

import (
	"errors"
	"fmt"
	"math"
	"testing"

	"github.com/corestoreio/csfw/config/cfgpath"
	"github.com/corestoreio/csfw/store/scope"
	"github.com/corestoreio/csfw/util/conv"
)

var benchmarkToString string

func benchmarkToStringF(b *testing.B, s interface{}) {
	b.ResetTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		benchmarkToString = conv.ToString(s)
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
		benchmarkToFloat64 = conv.ToFloat64(s)
	}
	if benchmarkToFloat64 < 0.001 {
		b.Fatal("benchmarkToFloat64 is empty :-(")
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
func BenchmarkToString_CfgPathPath(b *testing.B) {
	val := cfgpath.MustNewByParts("aa/bb/cc").Bind(scope.Store, 33)
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
