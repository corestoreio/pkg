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
//	// "¤#,##,##0.00;(¤#,##,##0.00)":       "te_IN", not implemented
//	"¤#0.00":                            "lv_LV",
//	"#,##0.00¤;(#,##0.00¤)":             "uk_UA",
//	"¤#,##0.00;¤-#,##0.00":              "sg_CF",
//	"#0.00\u00a0¤":                      "hy_AM",
//	"¤#,##0.00;(¤#,##0.00)":             "zu_ZA",
//	"#,##0.00¤":                         "zgh_MA",
//	"¤#,##0.00":                         "zh_Hans_SG",
//	"¤\u00a0#,##,##0.00":                "ur_PK", // \u00a0 no breaking space
//	"#,##,##0.00¤;(#,##,##0.00¤)":       "bn_IN", // 5,61,23,000.00 not implemented
//	"#,##0.00\u00a0¤;(#,##0.00\u00a0¤)": "yav_CM",
//	"¤\u00a0#,##0.00;¤-#,##0.00":        "it_CH",
//	"¤#,##,##0.00":                      "hi_IN", not implemented
//}

var testDefCurSym = i18n.Symbols{
	// normally that all should come from golang.org/x/text package
	Decimal: ',',
	Group:   '.',
	//	List:                   ';',
	//	PercentSign:            '%',
	//	CurrencySign:           '¤',
	PlusSign:  '+',
	MinusSign: '—', //  em dash \u2014 ;-)
	//	Exponential:            'E',
	//	SuperscriptingExponent: '×',
	//	PerMille:               '‰',
	//	Infinity:               '∞',
	Nan: []byte(`NaN`),
}

func TestFmtCurrency2(t *testing.T) {

	tests := []struct {
		opts    []i18n.CurrencyOptions
		sign    int
		i       int64
		prec    int
		dec     int64
		want    string
		wantErr error
	}{
		{
			[]i18n.CurrencyOptions{
				i18n.SetCurrencyFormat("#,##0.00 ¤"),
				i18n.SetCurrencyFraction(2, 0, 2, 0), // euro, 2 digits, no rounding
				i18n.SetCurrencySign([]byte("€")),
			},
			-1, -1234, 2, 6, "-1.234,06 €", nil, // euros with default Symbols
		},
		{
			[]i18n.CurrencyOptions{
				i18n.SetCurrencyFormat("#,##0.00 ¤"),
				i18n.SetCurrencyFraction(2, 0, 2, 0), // euro, 2 digits, no rounding
				i18n.SetCurrencySign([]byte("€")),
			},
			-1, -1234, 1, 6, "-1.234,60 €", nil, // euros with default Symbols
		},
		{
			[]i18n.CurrencyOptions{i18n.SetCurrencyFormat("+#,##0.00 ¤"), i18n.SetCurrencyFraction(2, 0, 2, 0), i18n.SetCurrencySign([]byte("€"))},
			0, 1234, 2, 6, "+1.234,06 €", nil, // number is 1234.06
		},
		{
			[]i18n.CurrencyOptions{i18n.SetCurrencyFormat("#,##0. ¤"), i18n.SetCurrencyFraction(2, 0, 2, 0), i18n.SetCurrencySign([]byte("€"))},
			1, 1234, 2, 6, "1.234,06 €", nil, // number is 1234.06
		},
		{
			[]i18n.CurrencyOptions{i18n.SetCurrencyFormat("#,##0. ¤"), i18n.SetCurrencyFraction(2, 0, 2, 0), i18n.SetCurrencySign([]byte("€"))},
			1, 1234, 3, 6, "1.234,01 €", nil, // number is 1234.006
		},
		{
			[]i18n.CurrencyOptions{i18n.SetCurrencyFormat("#,##0. ¤"), i18n.SetCurrencyFraction(2, 0, 2, 0), i18n.SetCurrencySign([]byte("€"))},
			1, 1234, 3, 345, "1.234,35 €", nil, // number is 1234.345
		},
		{
			[]i18n.CurrencyOptions{i18n.SetCurrencyFormat("#,##0. ¤"), i18n.SetCurrencyFraction(2, 0, 2, 0), i18n.SetCurrencySign([]byte("€"))},
			1, 1234, 3, 45, "1.234,05 €", nil, // number is 1234.045
		},
		{
			[]i18n.CurrencyOptions{i18n.SetCurrencyFormat("#,##0. ¤"), i18n.SetCurrencyFraction(0, 0, 2, 0), i18n.SetCurrencySign([]byte("€"))},
			1, 123456789, 1, 6, "123.456.790 €", nil, // 123456789.6
		},
		{
			[]i18n.CurrencyOptions{i18n.SetCurrencyFormat("#,##0. ¤"), i18n.SetCurrencyFraction(0, 0, 2, 0), i18n.SetCurrencySign([]byte("€"))},
			1, 123456789, 2, 6, "123.456.789 €", nil, // 123456789.06
		},
		{
			[]i18n.CurrencyOptions{i18n.SetCurrencyFormat("+0.00 ¤"), i18n.SetCurrencyFraction(2, 0, 2, 0), i18n.SetCurrencySign([]byte("€"))},
			-1, 4, 3, 245, "+4,25 €", nil, // 4.245
		},
		{
			[]i18n.CurrencyOptions{i18n.SetCurrencyFormat("¤\u00a0#0.00;¤\u00a0-#0.00"), i18n.SetCurrencyFraction(2, 0, 2, 0), i18n.SetCurrencySign([]byte("€"))},
			-1, 4, 2, 245, "€\u00a04,25", i18n.ErrPrecIsTooShort, // 4.245
		},
		{
			[]i18n.CurrencyOptions{i18n.SetCurrencyFormat("¤\u00a0#0.00;¤\u00a0-#0.00"), i18n.SetCurrencyFraction(2, 0, 2, 0), i18n.SetCurrencySign([]byte("€"))},
			-1, -12345678, 3, 245, "€\u00a0—12345678,25", nil, // 12345678.245
		},
		{
			[]i18n.CurrencyOptions{i18n.SetCurrencyFormat("¤\u00a0#0.00;¤\u00a0-#0.00"), i18n.SetCurrencyFraction(-1, -1, -1, -1), i18n.SetCurrencySign([]byte(""))},
			-1, -12345678, 3, 245, "\uf8ff\u00a0-12345678", nil, // 12345678.245
		},
		{
			[]i18n.CurrencyOptions{i18n.SetCurrencyFormat("¤\u00a0#0.00;¤\u00a0-#0.00"), i18n.SetCurrencyFraction(0, 0, 2, 0), i18n.SetCurrencySign([]byte(""))},
			-1, -12345678, 3, 495, "\uf8ff\u00a0-12345679", nil, // 12345678.495
		},
		{
			[]i18n.CurrencyOptions{i18n.SetCurrencyFormat("#,##0.00¤;(#,##0.00¤)"), i18n.SetCurrencyFraction(2, 0, 2, 0), i18n.SetCurrencySign([]byte("C"))},
			-1, -12345678, 3, 495, "(12.345.678,50C)", nil, // 12345678.495
		},
		{
			[]i18n.CurrencyOptions{i18n.SetCurrencyFormat("#,##0.00¤;(#,##0.00¤)"), i18n.SetCurrencyFraction(2, 0, 2, 0), i18n.SetCurrencySign([]byte("C"))},
			0, 0, 3, 495, "", i18n.ErrCannotDetectMinusSign, // 0.495
		},
		{
			[]i18n.CurrencyOptions{i18n.SetCurrencyFormat("#,##0.00¤;(#,##0.00¤)"), i18n.SetCurrencyFraction(2, 0, 2, 0), i18n.SetCurrencySign([]byte("C"))},
			-1, 0, 3, 495, "(0,50C)", nil, // 0.495
		},
		{
			[]i18n.CurrencyOptions{
				i18n.SetCurrencyFormat("#,##0.00 ¤"),
				i18n.SetCurrencyFraction(0, 0, 0, 0), // japanese yen, no digits no rounding
				i18n.SetCurrencySign([]byte("¥JP")),
			},
			1, 1234, 3, 456, "1.235 ¥JP", nil, // yen with default symbols
		},
	}
	var buf bytes.Buffer
	for _, test := range tests {
		haveNumber := i18n.NewCurrency(i18n.SetCurrencySymbols(testDefCurSym))
		haveNumber.CSetOptions(test.opts...)

		_, err := haveNumber.FmtNumber(&buf, test.sign, test.i, test.prec, test.dec)
		have := buf.String()
		if test.wantErr != nil {
			assert.Error(t, err)
			assert.EqualError(t, err, test.wantErr.Error())
		} else {
			assert.NoError(t, err)

			assert.EqualValues(t, test.want, have, "%v", test)
		}
		buf.Reset()
	}
}

func TestFmtCurrency3(t *testing.T) {
	// only to test the default format
	tests := []struct {
		opts    []i18n.CurrencyOptions
		sign    int
		i       int64
		prec    int
		dec     int64
		want    string
		wantErr error
	}{
		{
			[]i18n.CurrencyOptions{
				i18n.SetCurrencyFormat("", testDefCurSym),
				i18n.SetCurrencyFraction(2, 0, 2, 0), // euro, 2 digits, no rounding
				i18n.SetCurrencySign([]byte("€")),
			},
			-1, -1234, 2, 6, "€\u00a0-1.234,06", nil, // euros with default Symbols, -1234.06
		},
		{
			[]i18n.CurrencyOptions{
				i18n.SetCurrencyFormat("", testDefCurSym),
				i18n.SetCurrencyFraction(2, 0, 2, 0), // euro, 2 digits, no rounding
				i18n.SetCurrencySign([]byte("€")),
			},
			-1, -1234, 3, 6, "€\u00a0-1.234,01", nil, // euros with default Symbols, -1234.006
		},
		{
			[]i18n.CurrencyOptions{
				i18n.SetCurrencyFormat("", testDefCurSym),
				i18n.SetCurrencyFraction(0, 0, 0, 0), // euro, 2 digits, no rounding
				i18n.SetCurrencySign([]byte("€")),
			},
			1, 1234, 2, 495, "€\u00a0-1.235", i18n.ErrPrecIsTooShort, // euros with default Symbols, -1234.495
		},
		{
			[]i18n.CurrencyOptions{
				i18n.SetCurrencyFormat("", testDefCurSym),
				i18n.SetCurrencyFraction(0, 0, 0, 0), // euro, 2 digits, no rounding
				i18n.SetCurrencySign([]byte("€")),
			},
			1, 1234, 2, 44, "€\u00a01.234", nil, // euros with default Symbols, -1234.495
		},
	}

	var buf bytes.Buffer
	for _, test := range tests {
		haveNumber := i18n.NewCurrency(test.opts...)

		_, err := haveNumber.FmtNumber(&buf, test.sign, test.i, test.prec, test.dec)
		have := buf.String()
		if test.wantErr != nil {
			assert.Error(t, err)
			assert.EqualError(t, err, test.wantErr.Error())
		} else {
			assert.NoError(t, err)

			assert.EqualValues(t, test.want, have, "%v", test)
		}
		buf.Reset()
	}
}

func TestFmtCurrencyParallel(t *testing.T) {
	queue := make(chan fmtNumberData)
	ncpu := runtime.NumCPU()
	prevCPU := runtime.GOMAXPROCS(ncpu)
	defer runtime.GOMAXPROCS(prevCPU)
	wg := new(sync.WaitGroup)

	haveNumber := i18n.NewCurrency(
		i18n.SetCurrencyFormat("#,##0.000 ¤", testDefaultNumberSymbols),
		i18n.SetCurrencyFraction(3, 0, 0, 0),
	)

	// spawn workers
	for i := 0; i < ncpu; i++ {
		wg.Add(1)
		go testCurrencyWorker(t, haveNumber, i, queue, wg)
	}

	// master: give work
	for _, test := range genParallelTests(" $") {
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

		_, err := cf.FmtNumber(&buf, test.sign, test.i, test.prec, test.frac)
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
