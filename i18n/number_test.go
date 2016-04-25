// Copyright 2015-2016, Cyrill @ Schumacher.fm and the CoreStore contributors
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
	"math"
	"sync"
	"testing"

	"runtime"

	"github.com/corestoreio/csfw/i18n"
	"github.com/corestoreio/csfw/util/errors"
	"github.com/stretchr/testify/assert"
)

var testDefaultNumberSymbols = i18n.Symbols{
	// normally that all should come from golang.org/x/text package
	Decimal: '.',
	Group:   ',',
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

type fmtNumberData struct {
	format  string
	sign    int
	i       int64
	prec    int
	frac    int64
	want    string
	wantErr bool
}

func TestNumberFmtNumber1(t *testing.T) {

	// if Format empty default format displays at the moment: #,##0.### DefaultNumberFormat

	tests := []fmtNumberData{
		{"¤ #0.00", 1, 1234, 1, 9, "¤ 1234.90", false},
		{"###0.###", 1, 1234, 1, 9, "1234.900", false},
		{"###0.###", -1, 0, 2, 2, "-0.020", false},

		{"¤ #0.00", -1, -1234, 2, 6, "¤ -1234.06", false},
		{"#,##0.00 ¤", 1, 1234, 0, 0, "1,234.00 ¤", false},
		{"¤\u00a0#,##0.00;¤\u00a0-#,##0.00", -1, -1234, 3, 615, "¤\u00a0—1,234.62", false},
		{"¤\u00a0#,##0.00;¤\u00a0-#,##0.00", 1, 1234, 3, 454, "¤\u00a01,234.45", false},
		{"¤\u00a0#,##0.00;¤\u00a0-#,##0.00", -1, -1234, 7, 1234567, "¤\u00a0—1,234.12", false},

		{"+#,##0.###", 1, 1234, 3, 560, "+1,234.560", false},
		{"", -1, -1234, 3, 56, "-1,234.056", false},
		{"", -1, -1234, 3, 6, "-1,234.006", false},
		{"", -1, -1234, 3, 76, "-1,234.076", false},
		{"", 0, -1234, 3, 6, "-1,234.006", false},
		{"", 0, -1234, 2, 6, "-1,234.060", false},
		{"", 0, -1234, 1, 6, "-1,234.600", false},

		{"#,##0.00;(#,##0.00)", 1, 1234, 2, 56, "1,234.56", false},
		{"#,##0.00;(#,##0.00)", -1, -1234, 2, 56, "(1,234.56)", false},
		{"#,###.;(#,###.)", 1, 1234, 2, 56, "1,235", false},
		{"#,###.;(#,###.)", -1, -1234, 2, 56, "(1,235)", false},
		{"#.;(#.)", 1, 1234, 2, 56, "1235", false},
		{"#.;(#.)", -1, -1234, 2, 56, "(1235)", false},

		{"#,###.##", 1, 1234, 2, 56, "1,234.56", false},
		{"#,###.##", -1, -1234, 2, 56, "-1,234.56", false},
		{"#,###.##", -1, -1234, 2, 6, "-1,234.06", false},
		{"#,###.##", -1, -987651234, 3, 456, "-987,651,234.46", false},
		{"#,###.##", -1, -9876512341, 3, 454, "-9,876,512,341.45", false},

		{"#,##0.###", 1, 1234, 3, 56, "1,234.056", false},
		{"#,##0.###", -1, -1234, 3, 56, "-1,234.056", false},
		{"#,##0.###", -1, -1234, 3, 6, "-1,234.006", false},
		{"#,##0.###", -1, -1234, 4, 7678, "-1,234.768", false},
		{"#,##0.###", -1, -1234, 3, 6, "-1,234.006", false},
		{"#,##0.###", -1, -987651234, 3, 456, "-987,651,234.456", false},

		{"#0.###", 1, 1234, 3, 560, "1234.560", false},
		{"#0.###", -1, -1234, 3, 560, "-1234.560", false},
		{"#0.###", -1, -1234, 3, 60, "-1234.060", false},
		{"#0.###", -1, -1234, 4, 60, "-1234.006", false},
		{"#0.###", -1, -1234, 3, 76, "-1234.076", false},
		{"#0.###", -1, -1234, 3, 6, "-1234.006", false},
		{"#0.###", -1, -987651234, 3, 456, "-987651234.456", false},

		{"#,###.", 1, 1234, 2, 56, "1,235", false},
		{"#,###.", -1, -1234, 2, 56, "-1,235", false},
		{"#,###.", -1, -1234, 3, 495, "-1,235", false},
		{"#,###.", -1, 1234, 3, 495, "1,235", false},
		{"#,###.", -1, -1234, 2, 76, "-1,235", false},
		{"#,###.", -1, -1234, 2, 6, "-1,234", false},
		{"#,###.", -1, -1234, 2, 45, "-1,235", false},  // should we round down here?
		{"#,###.", -1, -1234, 3, 445, "-1,234", false}, // should we round up here?
		{"#,###.", -1, -1234, 2, 44, "-1,234", false},

		// invalid, because . is missing
		{"#,###", 1, 2234, 2, 56, "2,235", false},
		{"#,###", 1, 22, 2, 56, "23", false},
		{"+#,###", 1, 22, 2, 0, "+22", false},

		// invalid because . and , switched
		{"#.###,######", 1, 1234, 6, 567891, "1,234.5678910000", false},

		// invalid
		{"#\U0001f4b0###.##", 1, 1234, 2, 56, "1234.56\U0001f4b0", false},
		{"#\U0001f4b0###.##", -1, -1234, 2, 56, "-1234.56\U0001f4b0", false},

		{"+#,###.###", 1, 1234, 6, 567891, "+1,234.568", false},

		{"#,###0.###", 0, 0, 2, 6, "", true},
		{"$%^", 1, 1, 2, 6, "$%^1", false},
	}

	// all unique number formats and the last seen language
	// data from unicode CLDR
	// unicodeNumberFormats := map[string]string{"#0.###":"hy_AM", "#,##0.###":"zu_ZA", "#,##,##0.###":"te_IN", "#0.######":"en_US_POSIX"}
	// te_IN not tested as too rare

	for i, test := range tests {
		haveNumber := i18n.NewNumber(
			i18n.SetNumberFormat(test.format, testDefaultNumberSymbols),
		)
		var buf bytes.Buffer
		_, err := haveNumber.FmtNumber(&buf, test.sign, test.i, test.prec, test.frac)
		have := buf.String()
		if test.wantErr {
			assert.Error(t, err, "Index %d", i)
			assert.True(t, errors.IsNotValid(err), "Index %d => %s", i, err)
		} else {
			assert.NoError(t, err, "%v", test, "Index %d", i)
			assert.EqualValues(t, test.want, have, "Index %d", i)
		}
	}
}

func TestNumberFmtNumber2(t *testing.T) {
	// only to test the default format
	tests := []struct {
		opts    []i18n.NumberOptions
		sign    int
		i       int64
		prec    int
		frac    int64
		want    string
		wantErr bool
	}{
		{
			[]i18n.NumberOptions{
				i18n.SetNumberFormat(""),
				i18n.SetNumberSymbols(testDefCurSym),
			},
			-1, -1234, 3, 6, "-1,234.006", false, // euros with default Symbols
		},
		{
			[]i18n.NumberOptions{
				i18n.SetNumberFormat(""),
				i18n.SetNumberSymbols(testDefCurSym),
			},
			-1, -1234, 4, 6, "-1,234.001", false, // euros with default Symbols
		},
	}

	var buf bytes.Buffer
	for i, test := range tests {
		haveNumber := i18n.NewNumber(test.opts...)

		_, haveErr := haveNumber.FmtNumber(&buf, test.sign, test.i, test.prec, test.frac)
		have := buf.String()
		if test.wantErr {
			assert.Error(t, haveErr, "Index %d", i)
			assert.True(t, errors.IsNotValid(haveErr), "Index %d => %s", i, haveErr)
		} else {
			assert.NoError(t, haveErr, "Index %d", i)
			assert.EqualValues(t, test.want, have, "Index %d", i)
		}
		buf.Reset()
	}
}

func genParallelTests(suffix string) []fmtNumberData {
	tests := []fmtNumberData{}

	for i := 0; i < 500; i++ {
		// format is: "#,##0.###"
		tests = append(tests, fmtNumberData{"", 1, 1234, 2, 56, "1,234.560" + suffix, false})
		tests = append(tests, fmtNumberData{"", -1, -1234, 2, 56, "-1,234.560" + suffix, false})
		tests = append(tests, fmtNumberData{"", -1, -1234, 2, 6, "-1,234.060" + suffix, false})
		tests = append(tests, fmtNumberData{"", -1, -1234, 4, 7678, "-1,234.768" + suffix, false})
		tests = append(tests, fmtNumberData{"", -1, -1234, 3, 9, "-1,234.009" + suffix, false})
		tests = append(tests, fmtNumberData{"", -1, -1234, 1, 7, "-1,234.700" + suffix, false})
		tests = append(tests, fmtNumberData{"", -1, -987651234, 3, 456, "-987,651,234.456" + suffix, false})
		tests = append(tests, fmtNumberData{"", 0, 0, 2, 6, "", true})
		tests = append(tests, fmtNumberData{"", -1, 0, 1, 61, "", true})
	}
	return tests
}

func TestNumberFmtNumberParallel(t *testing.T) {
	queue := make(chan fmtNumberData)
	ncpu := runtime.NumCPU()
	prevCPU := runtime.GOMAXPROCS(ncpu)
	defer runtime.GOMAXPROCS(prevCPU)
	wg := new(sync.WaitGroup)

	haveNumber := i18n.NewNumber(
		i18n.SetNumberFormat("#,##0.###", testDefaultNumberSymbols),
	)

	// spawn workers
	for i := 0; i < ncpu; i++ {
		wg.Add(1)
		go testNumberWorker(t, haveNumber, i, queue, wg)
	}

	// master: give work
	for _, test := range genParallelTests("") {
		queue <- test
	}
	close(queue)
	wg.Wait()
}

func testNumberWorker(t *testing.T, nf i18n.NumberFormatter, id int, queue chan fmtNumberData, wg *sync.WaitGroup) {
	defer wg.Done()
	for {
		test, ok := <-queue
		if !ok {
			//t.Logf("Worker ID %d stopped", id)
			return
		}

		var buf bytes.Buffer
		_, haveErr := nf.FmtNumber(&buf, test.sign, test.i, test.prec, test.frac)
		have := buf.String()
		if test.wantErr {
			assert.Error(t, haveErr, "Worker %d => %v", id, test)
			assert.True(t, errors.IsNotValid(haveErr), "Worker %d => %s", id, haveErr)
		} else {
			assert.NoError(t, haveErr, "Worker %d => %v", id, test)
			assert.EqualValues(t, test.want, have, "Worker %d => %v", id, test)
		}
		//t.Logf("Worker %d run test: %v\n", id, test)
	}
}

func TestNumberFmtInt(t *testing.T) {

	tests := []struct {
		format  string
		i       int64
		want    string
		wantErr bool
	}{
		{"¤ #0.00", -1234, "¤ -1234.00", false},
		{"#,##0.00 ¤", 1234, "1,234.00 ¤", false},
		{"¤\u00a0#,##0.00;¤\u00a0-#,##0.00", -1234, "¤\u00a0—1,234.00", false},
		{"¤\u00a0#,##0.00;¤\u00a0-#,##0.00", 1234, "¤\u00a01,234.00", false},

		{"+#,##0.###", 1234, "+1,234.000", false},
		{"", -1234, "-1,234.000", false},

		{"#,##0.00;(#,##0.00)", 1234, "1,234.00", false},
		{"#,##0.00;(#,##0.00)", -1234, "(1,234.00)", false},
		{"#,###.;(#,###.)", 1234, "1,234", false},
		{"#,###.;(#,###.)", -1234, "(1,234)", false},
		{"#.;(#.)", 1234, "1234", false},
		{"#.;(#.)", -1234, "(1234)", false},

		{"#,##0.###", 123456, "123,456.000", false},

		// invalid
		{"#\U0001f4b0###.##", 1234, "1234.00\U0001f4b0", false},
		{"#\U0001f4b0###.##", -1234, "-1234.00\U0001f4b0", false},

		{"$%^", 2, "$%^2", false},
	}

	for i, test := range tests {
		haveNumber := i18n.NewNumber(
			i18n.SetNumberFormat(test.format, testDefaultNumberSymbols),
		)
		var buf bytes.Buffer
		_, haveErr := haveNumber.FmtInt64(&buf, test.i)
		have := buf.String()
		if test.wantErr {
			assert.Error(t, haveErr, "Index %d", i)
			assert.True(t, errors.IsNotValid(haveErr), "Index %d => %s", i, haveErr)
		} else {
			assert.NoError(t, haveErr, "Index %d", i)
			assert.EqualValues(t, test.want, have, "Index %d", i)
		}
	}
}

func TestNumberFmtFloat64(t *testing.T) {

	tests := []struct {
		format  string
		f       float64
		want    string
		wantErr bool
	}{
		{"¤ #0.00", -1234.456, "¤ -1234.46", false},
		{"#,##0.00 ¤", 1234.4456, "1,234.45 ¤", false},
		{"¤\u00a0#,##0.00;¤\u00a0-#,##0.00", -1234.1234567, "¤\u00a0—1,234.12", false},
		{"¤\u00a0#,##0.00;¤\u00a0-#,##0.00", 1234.1234567, "¤\u00a01,234.12", false},

		{"+#,##0.###", 1234.1 * 23.4, "+28,877.940", false},
		{"", -12.34 * 11.22, "-138.455", false},

		{"#,##0.00;(#,##0.00)", 1234, "1,234.00", false},
		{"#,##0.00;(#,##0.00)", -1234, "(1,234.00)", false},
		{"#,###.;(#,###.)", 1234.345 * 10, "12,343", false},
		{"#,###.;(#,###.)", -1234.345 * 10, "(12,343)", false},
		{"#.;(#.)", 1234 * 10, "12340", false},
		{"#.;(#.)", -1234 * 10, "(12340)", false},

		{"#,##0.###", 12345.6 * 10, "123,456.000", false},
		//
		//		// invalid
		{"#\U0001f4b0###.##", 1234, "1234.00\U0001f4b0", false},
		{"#\U0001f4b0###.##", -1234, "-1234.00\U0001f4b0", false},

		{"$%^", 2, "$%^2", false},
		{"$%^", math.NaN(), "NaN", false},

		{"#,##0.###", math.MaxFloat64, "∞", false},
		{"#,##0.###", -math.MaxFloat64, "—∞", false},
	}

	for i, test := range tests {
		haveNumber := i18n.NewNumber(
			i18n.SetNumberFormat(test.format, testDefaultNumberSymbols),
		)
		var buf bytes.Buffer
		_, haveErr := haveNumber.FmtFloat64(&buf, test.f)
		have := buf.String()
		if test.wantErr {
			assert.Error(t, haveErr, "Index %d", i)
			assert.True(t, errors.IsNotValid(haveErr), "Index %d => %s", i, haveErr)
		} else {
			assert.NoError(t, haveErr, "Index %d", i)
			assert.EqualValues(t, test.want, have, "Index %d", i)
		}
	}
}

func TestNumberGetFormat(t *testing.T) {
	haveNumber := i18n.NewNumber(
		i18n.SetNumberFormat(`€ #,##0.00 ;· (#,##0.00) °`, testDefaultNumberSymbols),
	)

	pf, err := haveNumber.GetFormat(false)
	assert.NoError(t, err)
	assert.EqualValues(t, "Parsed \ttrue\nPattern\t€ #,##0.00 \uf8ff\nPrec.  \t2\nPlus\t_\x00_\nMinus  \t_\x00_\nDecimal\t_._\nGroup \t_,_\nPrefix \t_€ _\nSuffix \t_ \uf8ff_\n", pf.String())

	nf, err := haveNumber.GetFormat(true)
	assert.NoError(t, err)
	assert.EqualValues(t, "Parsed \ttrue\nPattern\t· (#,##0.00) °\nPrec.  \t2\nPlus\t_\x00_\nMinus  \t_\x00_\nDecimal\t_._\nGroup \t_,_\nPrefix \t_· (_\nSuffix \t_) °_\n", nf.String())
}
