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

	"runtime"
	"sync"

	"github.com/corestoreio/csfw/i18n"
	"github.com/stretchr/testify/assert"
)

// all currency formats and the last seen in language
// data from unicode CLDR
//var allFormats = map[string]string{
//	"#,##0.00\u00a0¤": "xog_UG",
//	// "#,##0.00 ¤":                        "de_DE", duplicate of xog_UG
//	"¤\u00a0#,##0.00":                   "yi_001",
//	"¤\u00a0#0.00":                      "en_US_POSIX",
//	"¤\u00a0#,##0.00;¤\u00a0-#,##0.00":  "es_PY",
//	"¤#,##0.00;¤-\u00a0#,##0.00":        "luy_KE",
//	"¤\u00a0#,##0.00;(¤\u00a0#,##0.00)": "nl_SX",
//	"\u200e¤#,##0.00;\u200e(¤#,##0.00)": "fa_IR",
//	// "¤#,##,##0.00;(¤#,##,##0.00)":       "te_IN", also not tested in NumberFormatter
//	"¤#0.00":                            "lv_LV",
//	"#,##0.00¤;(#,##0.00¤)":             "uk_UA",
//	"¤#,##0.00;¤-#,##0.00":              "sg_CF",
//	"#0.00\u00a0¤":                      "hy_AM",
//	"¤#,##0.00;(¤#,##0.00)":             "zu_ZA",
//	"#,##0.00¤":                         "zgh_MA",
//	"¤#,##0.00":                         "zh_Hans_SG",
//	"¤\u00a0#,##,##0.00":                "ur_PK", // \u00a0 no breaking space
//	"#,##,##0.00¤;(#,##,##0.00¤)":       "bn_IN", // 5,61,23,000.00
//	"#,##0.00\u00a0¤;(#,##0.00\u00a0¤)": "yav_CM",
//	"¤\u00a0#,##0.00;¤-#,##0.00":        "it_CH",
//	"¤#,##,##0.00":                      "hi_IN",
//}

func TestFmtCurrency(t *testing.T) {

	tests := []struct {
		opts    []i18n.CurrencyOptFunc
		sign    int
		i       int64
		dec     int64
		want    string
		wantErr error
	}{
		{
			[]i18n.CurrencyOptFunc{
				i18n.CurrencyFormat("#,##0.00 ¤"),
				i18n.CurrencyFraction(2, 0, 2, 0), // euro, 2 digits, no rounding
				i18n.CurrencySign([]byte("€")),
			},
			-1, -1234, 6, "-1,234.06 €", nil, // euros with default Symbols
		},
		{
			[]i18n.CurrencyOptFunc{
				i18n.CurrencyFormat("#,##0.00 ¤"),
				i18n.CurrencyFraction(0, 0, 0, 0), // japanese yen, no digits no rounding
				i18n.CurrencySign([]byte("¥JP")),
			},
			1, 1234, 456, "1,235 ¥JP", nil, // yen with default symbols
		},
		//		{"¤\u00a0#,##0.00;¤\u00a0-#,##0.00", -1, -1234, 615, "¤\u00a0—1,234.62", nil},
		//		{"¤\u00a0#,##0.00;¤\u00a0-#,##0.00", 1, 1234, 454, "¤\u00a01,234.45", nil},
	}

	for _, test := range tests {
		haveNumber := i18n.NewCurrency(test.opts...)
		var buf bytes.Buffer
		_, err := haveNumber.FmtCurrency(&buf, test.sign, test.i, test.dec)
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

func TestFmtCurrencyParallel(t *testing.T) {
	queue := make(chan fmtNumberData)
	ncpu := runtime.NumCPU()
	prevCPU := runtime.GOMAXPROCS(ncpu)
	defer runtime.GOMAXPROCS(prevCPU)
	wg := new(sync.WaitGroup)

	haveNumber := i18n.NewCurrency(
		i18n.CurrencyFormat("#,##0.000 ¤", testDefaultNumberSymbols),
		i18n.CurrencyFraction(3, 0, 0, 0),
	)

	// spawn workers
	for i := 0; i < ncpu; i++ {
		wg.Add(1)
		go testCurrencyWorker(t, haveNumber, i, queue, wg)
	}

	// master: give work
	for _, test := range genParallelTests(" \U0001f4b0") {
		queue <- test
	}
	close(queue)
	wg.Wait()
}

func testCurrencyWorker(t *testing.T, cf i18n.CurrencyFormatter, id int, queue chan fmtNumberData, wg *sync.WaitGroup) {
	defer wg.Done()
	var buf bytes.Buffer
	for {
		test, ok := <-queue
		if !ok {
			//t.Logf("Worker ID %d stopped", id)
			return
		}

		_, err := cf.FmtCurrency(&buf, test.sign, test.i, test.dec)
		have := buf.String()
		if test.wantErr != nil {
			assert.Error(t, err, "Worker %d => %v", id, test)
			assert.EqualError(t, err, test.wantErr.Error(), "Worker %d => %v", id, test)
		} else {
			assert.NoError(t, err, "Worker %d => %v", id, test)
			assert.EqualValues(t, test.want, have, "Worker %d => %v", id, test)
		}
		buf.Reset()
		//t.Logf("Worker %d run test: %v\n", id, test)
	}
}
