package conv_test

import (
	"errors"
	"github.com/corestoreio/csfw/config/cfgpath"
	"github.com/corestoreio/csfw/store/scope"
	"github.com/corestoreio/csfw/util/conv"
	"math"
	"testing"
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
	val := cfgpath.MustNewByParts("aa/bb/cc").Bind(scope.StoreID, 33)
	benchmarkToStringF(b, val)
}
func BenchmarkToString_Error(b *testing.B) {
	var val = errors.New("Luke, I'm not your father.")
	benchmarkToStringF(b, val)
}
