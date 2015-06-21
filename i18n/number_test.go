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
	"math"
	"sync"
	"testing"

	"runtime"

	"github.com/corestoreio/csfw/i18n"
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
	dec     int64
	want    string
	wantErr error
}

func TestNumberFmtNumber1(t *testing.T) {

	// if Format empty default format displays at the moment: #,##0.### DefaultNumberFormat

	tests := []fmtNumberData{
		{"¤ #0.00", 1, 1234, 1, 9, "¤ 1234.90", nil},
		{"###0.###", 1, 1234, 1, 9, "1234.900", nil},

		{"¤ #0.00", -1, -1234, 2, 6, "¤ -1234.06", nil},
		{"#,##0.00 ¤", 1, 1234, 0, 0, "1,234.00 ¤", nil},
		{"¤\u00a0#,##0.00;¤\u00a0-#,##0.00", -1, -1234, 3, 615, "¤\u00a0—1,234.62", nil},
		{"¤\u00a0#,##0.00;¤\u00a0-#,##0.00", 1, 1234, 3, 454, "¤\u00a01,234.45", nil},
		{"¤\u00a0#,##0.00;¤\u00a0-#,##0.00", -1, -1234, 7, 1234567, "¤\u00a0—1,234.12", nil},

		{"+#,##0.###", 1, 1234, 3, 560, "+1,234.560", nil},
		{"", -1, -1234, 3, 56, "-1,234.056", nil},
		{"", -1, -1234, 3, 6, "-1,234.006", nil},
		{"", -1, -1234, 3, 76, "-1,234.076", nil},
		{"", 0, -1234, 3, 6, "-1,234.006", nil},
		{"", 0, -1234, 2, 6, "-1,234.060", nil},
		{"", 0, -1234, 1, 6, "-1,234.600", nil},

		{"#,##0.00;(#,##0.00)", 1, 1234, 2, 56, "1,234.56", nil},
		{"#,##0.00;(#,##0.00)", -1, -1234, 2, 56, "(1,234.56)", nil},
		{"#,###.;(#,###.)", 1, 1234, 2, 56, "1,235", nil},
		{"#,###.;(#,###.)", -1, -1234, 2, 56, "(1,235)", nil},
		{"#.;(#.)", 1, 1234, 2, 56, "1235", nil},
		{"#.;(#.)", -1, -1234, 2, 56, "(1235)", nil},

		{"#,###.##", 1, 1234, 2, 56, "1,234.56", nil},
		{"#,###.##", -1, -1234, 2, 56, "-1,234.56", nil},
		{"#,###.##", -1, -1234, 2, 6, "-1,234.06", nil},
		{"#,###.##", -1, -987651234, 3, 456, "-987,651,234.46", nil},
		{"#,###.##", -1, -9876512341, 3, 454, "-9,876,512,341.45", nil},

		{"#,##0.###", 1, 1234, 3, 56, "1,234.056", nil},
		{"#,##0.###", -1, -1234, 3, 56, "-1,234.056", nil},
		{"#,##0.###", -1, -1234, 3, 6, "-1,234.006", nil},
		{"#,##0.###", -1, -1234, 4, 7678, "-1,234.768", nil},
		{"#,##0.###", -1, -1234, 3, 6, "-1,234.006", nil},
		{"#,##0.###", -1, -987651234, 3, 456, "-987,651,234.456", nil},

		{"#0.###", 1, 1234, 3, 560, "1234.560", nil},
		{"#0.###", -1, -1234, 3, 560, "-1234.560", nil},
		{"#0.###", -1, -1234, 3, 60, "-1234.060", nil},
		{"#0.###", -1, -1234, 4, 60, "-1234.006", nil},
		{"#0.###", -1, -1234, 3, 76, "-1234.076", nil},
		{"#0.###", -1, -1234, 3, 6, "-1234.006", nil},
		{"#0.###", -1, -987651234, 3, 456, "-987651234.456", nil},

		{"#,###.", 1, 1234, 2, 56, "1,235", nil},
		{"#,###.", -1, -1234, 2, 56, "-1,235", nil},
		{"#,###.", -1, -1234, 3, 495, "-1,235", nil},
		{"#,###.", -1, 1234, 3, 495, "1,235", nil},
		{"#,###.", -1, -1234, 2, 76, "-1,235", nil},
		{"#,###.", -1, -1234, 2, 6, "-1,234", nil},
		{"#,###.", -1, -1234, 2, 45, "-1,235", nil},  // should we round down here?
		{"#,###.", -1, -1234, 3, 445, "-1,234", nil}, // should we round up here?
		{"#,###.", -1, -1234, 2, 44, "-1,234", nil},

		// invalid, because . is missing
		{"#,###", 1, 2234, 2, 56, "2,235", nil},
		{"#,###", 1, 22, 2, 56, "23", nil},
		{"+#,###", 1, 22, 2, 0, "+22", nil},

		// invalid because . and , switched
		{"#.###,######", 1, 1234, 6, 567891, "1,234.5678910000", nil},

		// invalid
		{"#\U0001f4b0###.##", 1, 1234, 2, 56, "1234.56\U0001f4b0", nil},
		{"#\U0001f4b0###.##", -1, -1234, 2, 56, "-1234.56\U0001f4b0", nil},

		{"+#,###.###", 1, 1234, 6, 567891, "+1,234.568", nil},

		//		{"#,###.###", math.NaN(), "NaN", nil},
		{"#,###0.###", 0, 0, 2, 6, "", i18n.ErrCannotDetectMinusSign},
		{"$%^", 1, 1, 2, 6, "$%^1", nil},
	}

	// all unique number formats and the last seen language
	// data from unicode CLDR
	// unicodeNumberFormats := map[string]string{"#0.###":"hy_AM", "#,##0.###":"zu_ZA", "#,##,##0.###":"te_IN", "#0.######":"en_US_POSIX"}
	// te_IN not tested as too rare

	for _, test := range tests {
		haveNumber := i18n.NewNumber(
			i18n.NumberFormat(test.format, testDefaultNumberSymbols),
		)
		var buf bytes.Buffer
		_, err := haveNumber.FmtNumber(&buf, test.sign, test.i, test.prec, test.dec)
		have := buf.String()
		if test.wantErr != nil {
			assert.Error(t, err, "%v", test)
			assert.EqualError(t, err, test.wantErr.Error(), "%v", test)
		} else {
			assert.NoError(t, err, "%v", test)
			assert.EqualValues(t, test.want, have, "%v", test)
		}
	}
}

func TestNumberFmtNumber2(t *testing.T) {
	// only to test the default format
	tests := []struct {
		opts    []i18n.NumberOptFunc
		sign    int
		i       int64
		prec    int
		dec     int64
		want    string
		wantErr error
	}{
		{
			[]i18n.NumberOptFunc{
				i18n.NumberFormat(""),
				i18n.NumberSymbols(testDefCurSym),
			},
			-1, -1234, 3, 6, "-1,234.006", nil, // euros with default Symbols
		},
		{
			[]i18n.NumberOptFunc{
				i18n.NumberFormat(""),
				i18n.NumberSymbols(testDefCurSym),
			},
			-1, -1234, 4, 6, "-1,234.001", nil, // euros with default Symbols
		},
	}

	var buf bytes.Buffer
	for _, test := range tests {
		haveNumber := i18n.NewNumber(test.opts...)

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

func genParallelTests(suffix string) []fmtNumberData {
	tests := []fmtNumberData{}

	for i := 0; i < 500; i++ {
		// format is: "#,##0.###"
		tests = append(tests, fmtNumberData{"", 1, 1234, 2, 56, "1,234.560" + suffix, nil})
		tests = append(tests, fmtNumberData{"", -1, -1234, 2, 56, "-1,234.560" + suffix, nil})
		tests = append(tests, fmtNumberData{"", -1, -1234, 2, 6, "-1,234.060" + suffix, nil})
		tests = append(tests, fmtNumberData{"", -1, -1234, 4, 7678, "-1,234.768" + suffix, nil})
		tests = append(tests, fmtNumberData{"", -1, -1234, 3, 9, "-1,234.009" + suffix, nil})
		tests = append(tests, fmtNumberData{"", -1, -1234, 1, 7, "-1,234.700" + suffix, nil})
		tests = append(tests, fmtNumberData{"", -1, -987651234, 3, 456, "-987,651,234.456" + suffix, nil})
		tests = append(tests, fmtNumberData{"", 0, 0, 2, 6, "", i18n.ErrCannotDetectMinusSign})
		tests = append(tests, fmtNumberData{"", -1, 0, 1, 61, "", i18n.ErrPrecIsTooShort})
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
		i18n.NumberFormat("#,##0.###", testDefaultNumberSymbols),
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
		_, err := nf.FmtNumber(&buf, test.sign, test.i, test.prec, test.dec)
		have := buf.String()
		if test.wantErr != nil {
			assert.Error(t, err, "Worker %d => %v", id, test)
			assert.EqualError(t, err, test.wantErr.Error(), "Worker %d => %v", id, test)
		} else {
			assert.NoError(t, err, "Worker %d => %v", id, test)
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
		wantErr error
	}{
		{"¤ #0.00", -1234, "¤ -1234.00", nil},
		{"#,##0.00 ¤", 1234, "1,234.00 ¤", nil},
		{"¤\u00a0#,##0.00;¤\u00a0-#,##0.00", -1234, "¤\u00a0—1,234.00", nil},
		{"¤\u00a0#,##0.00;¤\u00a0-#,##0.00", 1234, "¤\u00a01,234.00", nil},

		{"+#,##0.###", 1234, "+1,234.000", nil},
		{"", -1234, "-1,234.000", nil},

		{"#,##0.00;(#,##0.00)", 1234, "1,234.00", nil},
		{"#,##0.00;(#,##0.00)", -1234, "(1,234.00)", nil},
		{"#,###.;(#,###.)", 1234, "1,234", nil},
		{"#,###.;(#,###.)", -1234, "(1,234)", nil},
		{"#.;(#.)", 1234, "1234", nil},
		{"#.;(#.)", -1234, "(1234)", nil},

		{"#,##0.###", 123456, "123,456.000", nil},

		// invalid
		{"#\U0001f4b0###.##", 1234, "1234.00\U0001f4b0", nil},
		{"#\U0001f4b0###.##", -1234, "-1234.00\U0001f4b0", nil},

		{"$%^", 2, "$%^2", nil},
	}

	for _, test := range tests {
		haveNumber := i18n.NewNumber(
			i18n.NumberFormat(test.format, testDefaultNumberSymbols),
		)
		var buf bytes.Buffer
		_, err := haveNumber.FmtInt64(&buf, test.i)
		have := buf.String()
		if test.wantErr != nil {
			assert.Error(t, err, "%v", test)
			assert.EqualError(t, err, test.wantErr.Error(), "%v", test)
		} else {
			assert.NoError(t, err, "%v", test)
			assert.EqualValues(t, test.want, have, "%v", test)
		}
	}
}

func TestNumberFmtFloat64(t *testing.T) {

	tests := []struct {
		format  string
		f       float64
		want    string
		wantErr error
	}{
		{"¤ #0.00", -1234.456, "¤ -1234.46", nil},
		{"#,##0.00 ¤", 1234.4456, "1,234.45 ¤", nil},
		{"¤\u00a0#,##0.00;¤\u00a0-#,##0.00", -1234.1234567, "¤\u00a0—1,234.12", nil},
		{"¤\u00a0#,##0.00;¤\u00a0-#,##0.00", 1234.1234567, "¤\u00a01,234.12", nil},

		{"+#,##0.###", 1234.1 * 23.4, "+28,877.940", nil},
		{"", -12.34 * 11.22, "-138.455", nil},

		{"#,##0.00;(#,##0.00)", 1234, "1,234.00", nil},
		{"#,##0.00;(#,##0.00)", -1234, "(1,234.00)", nil},
		{"#,###.;(#,###.)", 1234.345 * 10, "12,343", nil},
		{"#,###.;(#,###.)", -1234.345 * 10, "(12,343)", nil},
		{"#.;(#.)", 1234 * 10, "12340", nil},
		{"#.;(#.)", -1234 * 10, "(12340)", nil},

		{"#,##0.###", 12345.6 * 10, "123,456.000", nil},
		//
		//		// invalid
		{"#\U0001f4b0###.##", 1234, "1234.00\U0001f4b0", nil},
		{"#\U0001f4b0###.##", -1234, "-1234.00\U0001f4b0", nil},

		{"$%^", 2, "$%^2", nil},
		{"$%^", math.NaN(), "NaN", nil},

		{"#,##0.###", math.MaxFloat64, "∞", nil},
		{"#,##0.###", -math.MaxFloat64, "—∞", nil},
	}

	for _, test := range tests {
		haveNumber := i18n.NewNumber(
			i18n.NumberFormat(test.format, testDefaultNumberSymbols),
		)
		var buf bytes.Buffer
		_, err := haveNumber.FmtFloat64(&buf, test.f)
		have := buf.String()
		if test.wantErr != nil {
			assert.Error(t, err, "%v", test)
			assert.EqualError(t, err, test.wantErr.Error(), "%v", test)
		} else {
			assert.NoError(t, err, "%v", test)
			assert.EqualValues(t, test.want, have, "%v", test)
		}
	}
}

func TestNumberGetFormat(t *testing.T) {
	haveNumber := i18n.NewNumber(
		i18n.NumberFormat(`€ #,##0.00 ;· (#,##0.00) °`, testDefaultNumberSymbols),
	)

	pf, err := haveNumber.GetFormat(false)
	assert.NoError(t, err)
	assert.EqualValues(t, "Parsed \ttrue\nPattern\t€ #,##0.00 \uf8ff\nPrec.  \t2\nPlus\t_\x00_\nMinus  \t_\x00_\nDecimal\t_._\nGroup \t_,_\nPrefix \t_€ _\nSuffix \t_ \uf8ff_\n", pf.String())

	nf, err := haveNumber.GetFormat(true)
	assert.NoError(t, err)
	assert.EqualValues(t, "Parsed \ttrue\nPattern\t· (#,##0.00) °\nPrec.  \t2\nPlus\t_\x00_\nMinus  \t_\x00_\nDecimal\t_._\nGroup \t_,_\nPrefix \t_· (_\nSuffix \t_) °_\n", nf.String())
}

var benchmarkFmtNumber string

// BenchmarkFmtNumber_UnCached_Pos	  300000	      5413 ns/op	    1760 B/op	      47 allocs/op
func BenchmarkFmtNumber_UnCached_Pos(b *testing.B) {
	bmFmtNumber_UnCached(b, "#,###.##", "1,234.57", 1, 1234, 3, 567)
}

// BenchmarkFmtNumber_UnCached_Neg	  200000	      7320 ns/op	    1888 B/op	      51 allocs/op
func BenchmarkFmtNumber_UnCached_Neg(b *testing.B) {
	bmFmtNumber_UnCached(b, "#,##0.00;(#,##0.00)", "(1,234.57)", -1, -1234, 3, 567)
}

func bmFmtNumber_UnCached(b *testing.B, format, want string, sign int, intgr int64, prec int, dec int64) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		haveNumber := i18n.NewNumber(
			i18n.NumberFormat(format, testDefaultNumberSymbols),
		)
		var buf bytes.Buffer
		if _, err := haveNumber.FmtNumber(&buf, sign, intgr, prec, dec); err != nil {
			b.Error(err)
		}
		have := buf.String()
		if have != want {
			b.Errorf("Missmatch %s vs %s", have, want)
		}
		benchmarkFmtNumber = have
	}
}

// BenchmarkFmtNumber___Cached_Pos	 2000000	       722 ns/op	      24 B/op	       5 allocs/op
func BenchmarkFmtNumber___Cached_Pos(b *testing.B) {
	bmFmtNumber_Cached(b, "#,###.##", "1,234.57", 1, 1234, 3, 567)
}

// BenchmarkFmtNumber___Cached_Int	 3000000	       593 ns/op	      21 B/op	       4 allocs/op
func BenchmarkFmtNumber___Cached_Int(b *testing.B) {
	bmFmtNumber_Cached(b, "#,###.", "1,234", 1, 1234, 2, 0)
}

// BenchmarkFmtNumber___Cached_Neg	 2000000	       722 ns/op	      32 B/op	       5 allocs/op
func BenchmarkFmtNumber___Cached_Neg(b *testing.B) {
	bmFmtNumber_Cached(b, "#,##0.00;(#,##0.00)", "(1,234.57)", -1, -1234, 3, 567)
}

func bmFmtNumber_Cached(b *testing.B, format, want string, sign int, intgr int64, prec int, dec int64) {
	b.ReportAllocs()
	haveNumber := i18n.NewNumber(
		i18n.NumberFormat(format, testDefaultNumberSymbols),
	)
	var buf bytes.Buffer
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		if _, err := haveNumber.FmtNumber(&buf, sign, intgr, prec, dec); err != nil {
			b.Error(err)
		}
		have := buf.String()
		if have != want {
			b.Errorf("Missmatch %s vs %s", have, want)
		}
		benchmarkFmtNumber = have
		buf.Reset()
	}
}
