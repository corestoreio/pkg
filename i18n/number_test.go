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

package i18n_test

import (
	"bytes"
	"errors"
	"testing"

	"math"

	"github.com/corestoreio/csfw/i18n"
	"github.com/stretchr/testify/assert"
)

func TestFmtNumber(t *testing.T) {

	// all unique number formats and the last seen language
	// data from unicode CLDR
	// unicodeNumberFormats := map[string]string{"#0.###":"hy_AM", "#,##0.###":"zu_ZA", "#,##,##0.###":"te_IN", "#0.######":"en_US_POSIX"}
	// te_IN not tested as too rare

	tests := []struct {
		format  string
		n       float64
		want    string
		wantErr error
	}{
		{"", 1234.56, "+1234.560000000", nil},
		{"", -1234.56, "-1234.560000000", nil},
		{"", -1234.06, "-1234.060000000", nil},
		{"", -1234.076, "-1234.076000000", nil},
		{"", -1234.0006, "-1234.000600000", nil},

		{"#,###.##", 1234.56, "1,234.56", nil},
		{"#,###.##", -1234.56, "-1,234.56", nil},
		{"#,###.##", -1234.06, "-1,234.06", nil},
		{"#,###.##", -1234.076, "-1,234.08", nil},
		{"#,###.##", -1234.0006, "-1,234.00", nil},
		{"#,###.##", -987651234.456, "-987,651,234.46", nil},

		{"#,##0.###", 1234.56, "1,234.560", nil},
		{"#,##0.###", -1234.56, "-1,234.560", nil},
		{"#,##0.###", -1234.06, "-1,234.060", nil},
		{"#,##0.###", -1234.076, "-1,234.076", nil},
		{"#,##0.###", -1234.0006, "-1,234.001", nil},
		{"#,##0.###", -987651234.456, "-987,651,234.456", nil},

		{"#0.###", 1234.56, "1234.560", nil},
		{"#0.###", -1234.56, "-1234.560", nil},
		{"#0.###", -1234.06, "-1234.060", nil},
		{"#0.###", -1234.076, "-1234.076", nil},
		{"#0.###", -1234.0006, "-1234.001", nil},
		{"#0.###", -987651234.456, "-987651234.456", nil},

		{"#,###.", 1234.56, "1,235", nil},
		{"#,###.", -1234.56, "-1,235", nil},
		{"#,###.", -1234.495, "-1,234", nil},
		{"#,###.", 1234.495, "1,234", nil},
		{"#,###.", -1234.076, "-1,234", nil},
		{"#,###.", -1234.0006, "-1,234", nil},

		{"#,###", 1234.56, "1234,560", nil},
		{"#,###", -1234.56, "-1234,560", nil},
		{"#,###", -1234.495, "-1234,495", nil},
		{"#,###", 1234.495, "1234,495", nil},
		{"#,###", -1234.076, "-1234,076", nil},
		{"#,###", -1234.0006, "-1234,001", nil},

		{"#.###,######", 1234.567891, "1.234,567891", nil},
		{"#.###,######", -1234.56, "-1.234,560000", nil},
		{"#.###,######", -1234.495, "-1.234,495000", nil},
		{"#.###,######", 1234.495, "1.234,495000", nil},
		{"#.###,######", -1234.076, "-1.234,076000", nil},
		{"#.###,######", -1234.0006, "-1.234,000600", nil},

		{"#\U0001f4b0###,##", 1234.56, "1\U0001f4b0234,56", nil},
		{"#\U0001f4b0###,##", -1234.56, "-1\U0001f4b0234,56", nil},

		{"+#.###,###", 1234.567891, "+1.234,568", nil},

		{"#.###,###", math.NaN(), "NaN", nil},

		{"$#,##0.###", -1234.0006, "-1,234.001", errors.New("Invalid positive sign directive in format: $#,##0.###")},
		{"#,###0.###", -1234.0006, "-1,234.001", errors.New("Group separator directive must be followed by 3 digit-specifiers in format: #,###0.###")},
	}

	for _, test := range tests {
		haveNumber := i18n.NewNumber(
			i18n.NumberFormat(test.format),
		)
		var buf bytes.Buffer
		_, err := haveNumber.FmtNumber(&buf, test.n)
		have := buf.String()
		if test.wantErr != nil {
			assert.Error(t, err)
			assert.EqualError(t, err, test.wantErr.Error())
		} else {
			assert.NoError(t, err)

			assert.EqualValues(t, test.want, have, "%v", test)
		}
	}
}

var benchmarkFmtNumber string

// BenchmarkFmtNumber_UnCached	 1000000	      1875 ns/op	     520 B/op	      17 allocs/op
func BenchmarkFmtNumber_UnCached(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {

		haveNumber := i18n.NewNumber(
			i18n.NumberFormat("#,###.##"),
		)
		var buf bytes.Buffer
		if _, err := haveNumber.FmtNumber(&buf, 1234.567); err != nil {
			b.Error(err)
		}
		have := buf.String()
		if have != "1,234.57" {
			b.Errorf("Missmatch %s vs 1,234.56", have)
		}
		benchmarkFmtNumber = have
	}
}

// BenchmarkFmtNumber___Cached	 2000000	       814 ns/op	     160 B/op	       8 allocs/op
func BenchmarkFmtNumber___Cached(b *testing.B) {
	b.ReportAllocs()
	haveNumber := i18n.NewNumber(
		i18n.NumberFormat("#,###.##"),
	)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		var buf bytes.Buffer
		if _, err := haveNumber.FmtNumber(&buf, 1234.567); err != nil {
			b.Error(err)
		}
		have := buf.String()
		if have != "1,234.57" {
			b.Errorf("Missmatch %s vs 1,234.56", have)
		}
		benchmarkFmtNumber = have
	}
}
