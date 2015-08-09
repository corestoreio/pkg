// Copyright 2015, Cyrill @ Schumacher.fm and the CoreStore contributors
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

package locale_test

import "testing"
import (
	"github.com/corestoreio/csfw/locale"
	"github.com/stretchr/testify/assert"
)

func TestExtractNumber(t *testing.T) {

	tests := []struct {
		have string
		want float64
	}{
		{"  2345.4356,1234", 23454356.1234},
		{" \t 2345 \n", 2345},
		{"'+23,3452.123'", 233452.123},
		{"'-23,3452.123'", -233452.123},
		{"'-323.452,123'", -323452.123},
		{"' 12343 '", 12343},
		{"' 123-43 '", 123},
		{"'-9456km'", -9456},
		{"0", 0},
		{"+0.2343", 0.2343},
		{"2 054,10", 2054.1},
		{"-", 0},
		{"2014.30", 2014.3},
		{"2'054.52", 2054.52},
		{"2,46 GB", 2.46},
		{"2,44 GB 34.56kb", 2.44},
		{"2,44 GB 34.56kb", 2.44},
		{"", 0},
		{`<IMG SRC=j&#X41vascript:alert('test2')>`, 41},
	}

	for _, test := range tests {
		out, err := locale.ExtractNumber(test.have)
		assert.NoError(t, err, "have %s want %f actual %f", test.have, test.want, out)
		assert.True(t, test.want == out, "have %s want %f actual %f", test.have, test.want, out) //   comparing floats ... >8-)
	}
}

// BenchmarkExtractNumberLong-4 	 3000000	       459 ns/op	      80 B/op	       2 allocs/op
func BenchmarkExtractNumberLong(b *testing.B) {
	benchmarkExtractNumber(b, "'-323.452,123'", -323452.123)
}

// BenchmarkExtractNumberInt-4  	10000000	       224 ns/op	      32 B/op	       2 allocs/op
func BenchmarkExtractNumberInt(b *testing.B) {
	benchmarkExtractNumber(b, " 452", 452)
}

// BenchmarkExtractNumberFloat-4	 5000000	       353 ns/op	      64 B/op	       2 allocs/op
func BenchmarkExtractNumberFloat(b *testing.B) {
	benchmarkExtractNumber(b, " 45,43232", 45.43232)
}

func benchmarkExtractNumber(b *testing.B, input string, want float64) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		have, err := locale.ExtractNumber(input)
		if err != nil {
			b.Error(err)
		}
		if have != want {
			b.Errorf("%f != %f", have, want)
		}
	}
}

func TestFormatPrice(t *testing.T) {
	t.Log("@todo")
}
