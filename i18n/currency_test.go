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
	"testing"

	"github.com/corestoreio/csfw/i18n"
	"github.com/stretchr/testify/assert"
)

// all currency formats and the last seen in language
// data from unicode CLDR
var allFormats = map[string]string{
	"#,##0.00\u00a0¤": "xog_UG",
	// "#,##0.00 ¤":                        "de_DE", duplicate of xog_UG
	"¤\u00a0#,##0.00":                   "yi_001",
	"¤\u00a0#0.00":                      "en_US_POSIX",
	"¤\u00a0#,##0.00;¤\u00a0-#,##0.00":  "es_PY",
	"¤#,##0.00;¤-\u00a0#,##0.00":        "luy_KE",
	"¤\u00a0#,##0.00;(¤\u00a0#,##0.00)": "nl_SX",
	"\u200e¤#,##0.00;\u200e(¤#,##0.00)": "fa_IR",
	// "¤#,##,##0.00;(¤#,##,##0.00)":       "te_IN", also not tested in NumberFormatter
	"¤#0.00":                            "lv_LV",
	"#,##0.00¤;(#,##0.00¤)":             "uk_UA",
	"¤#,##0.00;¤-#,##0.00":              "sg_CF",
	"#0.00\u00a0¤":                      "hy_AM",
	"¤#,##0.00;(¤#,##0.00)":             "zu_ZA",
	"#,##0.00¤":                         "zgh_MA",
	"¤#,##0.00":                         "zh_Hans_SG",
	"¤\u00a0#,##,##0.00":                "ur_PK", // \u00a0 no breaking space
	"#,##,##0.00¤;(#,##,##0.00¤)":       "bn_IN",
	"#,##0.00\u00a0¤;(#,##0.00\u00a0¤)": "yav_CM",
	"¤\u00a0#,##0.00;¤-#,##0.00":        "it_CH",
	"¤#,##,##0.00":                      "hi_IN",
}

func TestFmtCurrency(t *testing.T) {

	tests := []struct {
		format  string
		n       float64
		want    string
		wantErr error
	}{
		//		{"", 1234.56, "+1234.560000000", nil},
		//		{"", -1234.56, "-1234.560000000", nil},
		//		{"", -1234.06, "-1234.060000000", nil},
		//		{"", -1234.076, "-1234.076000000", nil},
		//		{"", -1234.0006, "-1234.000600000", nil},

		{"#,##0.00 ¤", 1234.56, "1,234.56", nil},
		{"#,##0.00 ¤", -1234.56, "-1,234.56", nil},
		{"#,##0.00 ¤", -1234.06, "-1,234.06", nil},
		{"#,##0.00 ¤", -1234.076, "-1,234.08", nil},
		{"#,##0.00 ¤", -1234.0006, "-1,234.00", nil},
		{"#,##0.00 ¤", -987651234.456, "-987,651,234.46", nil},

		//		{"¤#,##0.00", 1234.56, "1,234.56", nil},
		//		{"¤#,##0.00", -1234.56, "-1,234.56", nil},
		//		{"¤#,##0.00", -1234.06, "-1,234.06", nil},
		//		{"¤#,##0.00", -1234.076, "-1,234.08", nil},
		//		{"¤#,##0.00", -1234.0006, "-1,234.00", nil},
		//		{"¤#,##0.00", -987651234.456, "-987,651,234.46", nil},
		//
		//		{"¤ #,##0.00;¤-#,##0.00", 1234.56, "1,234.56", nil},
		//		{"¤ #,##0.00;¤-#,##0.00", -1234.56, "-1,234.56", nil},
		//		{"¤ #,##0.00;¤-#,##0.00", -1234.06, "-1,234.06", nil},
		//		{"¤ #,##0.00;¤-#,##0.00", -1234.076, "-1,234.08", nil},
		//		{"¤ #,##0.00;¤-#,##0.00", -1234.0006, "-1,234.00", nil},
		//		{"¤ #,##0.00;¤-#,##0.00", -987651234.456, "-987,651,234.46", nil},
		//
		//		{"¤ #,##0.00;¤ -#,##0.00", 1234.56, "1,234.56", nil},
		//		{"¤ #,##0.00;¤ -#,##0.00", -1234.56, "-1,234.56", nil},
		//		{"¤ #,##0.00;¤ -#,##0.00", -1234.06, "-1,234.06", nil},
		//		{"¤ #,##0.00;¤ -#,##0.00", -1234.076, "-1,234.08", nil},
		//		{"¤ #,##0.00;¤ -#,##0.00", -1234.0006, "-1,234.00", nil},
		//		{"¤ #,##0.00;¤ -#,##0.00", -987651234.456, "-987,651,234.46", nil},

		//		{"#,###.", 1234.56, "1,235", nil},
		//		{"#,###.", -1234.56, "-1,235", nil},
		//		{"#,###.", -1234.495, "-1,234", nil},
		//		{"#,###.", 1234.495, "1,234", nil},
		//		{"#,###.", -1234.076, "-1,234", nil},
		//		{"#,###.", -1234.0006, "-1,234", nil},
		//
		//		{"#,###", 1234.56, "1234,560", nil},
		//		{"#,###", -1234.56, "-1234,560", nil},
		//		{"#,###", -1234.495, "-1234,495", nil},
		//		{"#,###", 1234.495, "1234,495", nil},
		//		{"#,###", -1234.076, "-1234,076", nil},
		//		{"#,###", -1234.0006, "-1234,001", nil},
		//
		//		{"#.###,######", 1234.567891, "1.234,567891", nil},
		//		{"#.###,######", -1234.56, "-1.234,560000", nil},
		//		{"#.###,######", -1234.495, "-1.234,495000", nil},
		//		{"#.###,######", 1234.495, "1.234,495000", nil},
		//		{"#.###,######", -1234.076, "-1.234,076000", nil},
		//		{"#.###,######", -1234.0006, "-1.234,000600", nil},
		//
		//		{"#\U0001f4b0###,##", 1234.56, "1\U0001f4b0234,56", nil},
		//		{"#\U0001f4b0###,##", -1234.56, "-1\U0001f4b0234,56", nil},

		// @todo more tests
	}

	for _, test := range tests {
		haveNumber := i18n.NewCurrency(
			i18n.CurrencyFormat(test.format),
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
