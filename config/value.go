// Copyright 2015-present, Cyrill @ Schumacher.fm and the CoreStore contributors
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

package config

import (
	"bytes"
	"crypto/subtle"
	"encoding/csv"
	"fmt"
	"io"
	"strings"
	"time"
	"unicode/utf8"

	"github.com/corestoreio/errors"
	"github.com/corestoreio/pkg/util/byteconv"
)

// CSVColumnSeparator separates CSV values. Default value.
const CSVColumnSeparator = ','

// ValueConverter can convert the raw underlying value to the final data.
type ValueConverter interface {
	// Convert takes an optionally Path and the raw bytes to convert the raw
	// bytes into their final representation. An example would be base64
	// decoding.
	Convert(Path, []byte) ([]byte, error)
}

type convertFn func(Path, []byte) ([]byte, error)

func (cfn convertFn) Convert(p Path, data []byte) ([]byte, error) {
	return cfn(p, data)
}

const (
	valFoundNo = iota
	valFoundYes
	valFoundLRU
)

func valFoundStringer(found uint8) string {
	switch found {
	case valFoundNo:
		return "NO"
	case valFoundYes:
		return "YES"
	case valFoundLRU:
		return "LRU"
	}
	return "CONFIG:FOUND_UNDEFINED"
}

// Value represents an immutable value returned from the configuration service.
// A value is meant to be only for reading.
type Value struct {
	// Path optionally assigned to, to know to which path a value belongs to and to
	// provide different converter behaviour.
	Path Path
	data []byte
	// Converter converts on any method receiver call the byte value.
	Converter ValueConverter
	// CSVColumnSep defines the CSV column separator, default a comma.
	CSVReader *csv.Reader
	CSVComma  rune
	// found gets set to greater zero if any value can be found under the given
	// path. Even a NULL value can be valid. found gets also used as a
	// statistical flag to identify where a value comes from, e.g. from backend
	// or from LRU.
	found   uint8
	lastErr error
}

// MakeValue makes a new non-pointer value type.
func MakeValue(data []byte) Value {
	return Value{
		data:  data,
		found: valFoundYes,
	}
}

// WithConvert applies a function as a ValueConverter.
func (v Value) WithConvert(fn func(Path, []byte) ([]byte, error)) Value {
	v.Converter = convertFn(fn)
	return v
}

func (v Value) init() (_ Value, found bool, err error) {
	if v.lastErr != nil || valFoundNo == v.found {
		return Value{}, false, v.lastErr
	}

	if v.CSVComma == 0 {
		v.CSVComma = CSVColumnSeparator
	}
	if v.Converter != nil {
		v.data, err = v.Converter.Convert(v.Path, v.data)
	}
	return v, v.found > valFoundNo && err == nil, err
}

// String implements fmt.Stringer and returns the textual representation and Go
// syntax escaped of the underlying data.
func (v Value) String() string {
	v, found, err := v.init()
	if err != nil {
		return fmt.Sprintf("[config] Value: %+v", err)
	}
	if !found {
		return "<notFound>"
	}
	if v.data == nil {
		return "<nil>"
	}
	return fmt.Sprintf("%q", v.data)
}

// WriteTo writes the converted data to w.
func (v Value) WriteTo(w io.Writer) (n int64, err error) {
	if v, _, err = v.init(); err != nil {
		return 0, errors.WithStack(err)
	}
	nw, err := w.Write(v.data)
	return int64(nw), err
}

// Str returns the underlying value as a string. Ok is true when data slice
// bytes are not nil. Whereas function String implements fmt.Stringer and
// returns something different.
func (v Value) Str() (_ string, ok bool, err error) {
	v, ok, err = v.init()
	return string(v.data), ok && v.data != nil, err
}

// Strs splits the converted data using the CSVComma and appends it to `ret`.
func (v Value) Strs(ret ...string) (_ []string, err error) {
	if v, _, err = v.init(); err != nil {
		return nil, errors.WithStack(err)
	}
	if v.data == nil {
		return ret, nil
	}
	sep := string(v.CSVComma)
	s := string(v.data)
	n := strings.Count(s, sep)

	if ret == nil {
		ret = make([]string, 0, n+1)
	}
	i := 0
	for i < n {
		m := strings.Index(s, sep)
		if m < 0 {
			break
		}
		if s[:m] != "" {
			ret = append(ret, s[:m])
		}
		s = s[m+len(sep):]
		i++
	}
	if s != "" {
		ret = append(ret, s)
	}
	return ret, nil
}

// CSV reads a multiline CSV data value and appends it to ret.
func (v Value) CSV(ret ...[]string) (_ [][]string, err error) {
	if v, _, err = v.init(); err != nil {
		return nil, errors.WithStack(err)
	}
	if v.data == nil {
		return ret, nil
	}
	if err == nil && v.CSVReader == nil {
		v.CSVReader = csv.NewReader(bytes.NewReader(v.data))
	}
	v.CSVReader.Comma = v.CSVComma
	// Implement reuse record correctly ... bit complex and not sure if worth.
	for err == nil {
		var r []string
		if r, err = v.CSVReader.Read(); err == nil {
			ret = append(ret, r)
		}
	}
	if err == io.EOF {
		err = nil
	}
	return ret, err
}

// Unmarshal decodes the value into the final type. vPtr must be a pointer. The
// function signature for `fn` matches e.g. json.Unmarshal, xml.Unmarshal and
// many others.
func (v Value) Unmarshal(fn func(data []byte, vPtr interface{}) error, vPtr interface{}) (err error) {
	if v, _, err = v.init(); err != nil {
		return errors.WithStack(err)
	}
	if v.data == nil {
		return nil
	}
	return fn(v.data, vPtr)
}

// Bool converts the underlying converted data into the final type.
func (v Value) Bool() (val bool, ok bool, err error) {
	if v, ok, err = v.init(); err != nil || !ok {
		return false, false, errors.WithStack(err)
	}
	return byteconv.ParseBool(v.data)
}

// Float64 converts the underlying converted data into the final type.
func (v Value) Float64() (val float64, ok bool, err error) {
	if v, ok, err = v.init(); err != nil || !ok {
		return 0, false, errors.WithStack(err)
	}
	return byteconv.ParseFloat(v.data)
}

func (v Value) countSep() (n int, sep []byte) {
	var sep2 [4]byte
	rl := utf8.EncodeRune(sep2[:], v.CSVComma)
	sep = sep2[:rl]
	n = bytes.Count(v.data, sep)
	return
}

// Float64s splits the converted data using the CSVComma and appends it to `ret`.
func (v Value) Float64s(ret ...float64) (_ []float64, err error) {
	if v, _, err = v.init(); err != nil {
		return nil, errors.WithStack(err)
	}
	if v.data == nil {
		return ret, nil
	}
	n, sep := v.countSep()
	if ret == nil {
		ret = make([]float64, 0, n+1)
	}
	i := 0
	s := v.data
	for i < n {
		m := bytes.Index(s, sep)
		if m < 0 {
			break
		}

		if f64, ok, err := byteconv.ParseFloat(s[:m]); err != nil {
			return nil, errors.Wrapf(err, "[config] Value.Float64s with index %d and entry %q", i, s[:m])
		} else if ok {
			ret = append(ret, f64)
		}
		s = s[m+len(sep):]
		i++
	}

	if f64, ok, err := byteconv.ParseFloat(s); err != nil {
		return nil, errors.Wrapf(err, "[config] Value.Float64s with index %d and entry %q", i, s)
	} else if ok {
		ret = append(ret, f64)
	}

	return ret, nil
}

// Int converts the underlying converted data into the final type.
func (v Value) Int() (val int, ok bool, err error) {
	if v, ok, err = v.init(); err != nil || !ok {
		return 0, false, errors.WithStack(err)
	}
	i64, ok, err := byteconv.ParseInt(v.data)
	return int(i64), ok, err
}

// Ints converts the underlying byte slice into an int64 slice using
// v.CSVComma as a separator. The result gets append to argument `ret`.
func (v Value) Ints(ret ...int) (_ []int, err error) {
	if v, _, err = v.init(); err != nil {
		return nil, errors.WithStack(err)
	}
	if v.data == nil {
		return ret, nil
	}
	n, sep := v.countSep()
	if ret == nil {
		ret = make([]int, 0, n+1)
	}
	i := 0
	s := v.data
	for i < n {
		m := bytes.Index(s, sep)
		if m < 0 {
			break
		}

		if i64, ok, err := byteconv.ParseInt(s[:m]); err != nil {
			return nil, errors.Wrapf(err, "[config] Value.Ints with index %d and entry %q", i, s[:m])
		} else if ok {
			ret = append(ret, int(i64))
		}
		s = s[m+len(sep):]
		i++
	}

	if i64, ok, err := byteconv.ParseInt(s); err != nil {
		return nil, errors.Wrapf(err, "[config] Value.Ints with index %d and entry %q", i, s)
	} else if ok {
		ret = append(ret, int(i64))
	}

	return ret, nil
}

// Int64 converts the underlying converted data into the final type.
func (v Value) Int64() (_ int64, ok bool, err error) {
	if v, ok, err = v.init(); err != nil || !ok {
		return 0, false, errors.WithStack(err)
	}
	return byteconv.ParseInt(v.data)
}

// Int64s converts the underlying byte slice into an int64 slice using
// v.CSVComma as a separator. The result gets append to argument `ret`.
func (v Value) Int64s(ret ...int64) (_ []int64, err error) {
	if v, _, err = v.init(); err != nil {
		return nil, errors.WithStack(err)
	}
	if v.data == nil {
		return ret, nil
	}
	n, sep := v.countSep()
	if ret == nil {
		ret = make([]int64, 0, n+1)
	}
	i := 0
	s := v.data
	for i < n {
		m := bytes.Index(s, sep)
		if m < 0 {
			break
		}

		if i64, ok, err := byteconv.ParseInt(s[:m]); err != nil {
			return nil, errors.Wrapf(err, "[config] Value.Int64s with index %d and entry %q", i, s[:m])
		} else if ok {
			ret = append(ret, i64)
		}
		s = s[m+len(sep):]
		i++
	}

	if i64, ok, err := byteconv.ParseInt(s); err != nil {
		return nil, errors.Wrapf(err, "[config] Value.Int64s with index %d and entry %q", i, s)
	} else if ok {
		ret = append(ret, i64)
	}

	return ret, nil
}

// Uint64 converts the underlying converted data into the final type.
func (v Value) Uint64() (_ uint64, ok bool, err error) {
	if v, ok, err = v.init(); err != nil || !ok {
		return 0, false, errors.WithStack(err)
	}
	return byteconv.ParseUint(v.data, 10, 64)
}

// Uint64s converts the underlying byte slice into an int64 slice using
// v.CSVComma as a separator. The result gets append to argument `ret`.
func (v Value) Uint64s(ret ...uint64) (_ []uint64, err error) {
	if v, _, err = v.init(); err != nil {
		return nil, errors.WithStack(err)
	}
	if v.data == nil {
		return ret, nil
	}
	n, sep := v.countSep()
	if ret == nil {
		ret = make([]uint64, 0, n+1)
	}
	i := 0
	s := v.data
	for i < n {
		m := bytes.Index(s, sep)
		if m < 0 {
			break
		}
		if i64, ok, err := byteconv.ParseUint(s[:m], 10, 64); err != nil {
			return nil, errors.Wrapf(err, "[config] Value.Uint64s with index %d and entry %q", i, s[:m])
		} else if ok {
			ret = append(ret, i64)
		}
		s = s[m+len(sep):]
		i++
	}
	if len(s) > 0 {
		if i64, ok, err := byteconv.ParseUint(s, 10, 64); err != nil {
			return nil, errors.Wrapf(err, "[config] Value.Uint64s with index %d and entry %q", i, s)
		} else if ok {
			ret = append(ret, i64)
		}
	}
	return ret, nil
}

// Time parses a MySQL/MariaDB like time format:
// "2006-01-02 15:04:05.999999999 07:00" up to 35
// places supports. Minimal format must be a year, e.g. 2006.
// time.UTC location gets used.
func (v Value) Time() (t time.Time, ok bool, err error) {
	if v, ok, err = v.init(); err != nil || !ok {
		err = errors.WithStack(err)
		return
	}
	if v.IsEmpty() {
		return
	}
	t, err = parseDateTime(string(v.data), time.UTC)
	ok = err == nil
	return
}

// Times same as Time but parses the CSVComma separated list.
func (v Value) Times(ret ...time.Time) (t []time.Time, err error) {
	if v, _, err = v.init(); err != nil {
		return nil, errors.WithStack(err)
	}
	if v.data == nil {
		return ret, nil
	}
	sep := string(v.CSVComma)
	s := string(v.data)
	n := strings.Count(s, sep)

	if ret == nil {
		ret = make([]time.Time, 0, n+1)
	}
	i := 0
	for i < n {
		m := strings.Index(s, sep)
		if m < 0 {
			break
		}
		if s[:m] != "" {
			t, err := parseDateTime(s[:m], time.UTC)
			if err != nil {
				return nil, errors.Wrapf(err, "[config] Value.Times with index %d and entry %q", i, s[:m])
			}
			ret = append(ret, t)
		}
		s = s[m+len(sep):]
		i++
	}
	if s != "" {
		t, err := parseDateTime(s, time.UTC)
		if err != nil {
			return nil, errors.Wrapf(err, "[config] Value.Times with index %d and entry %q", i, s)
		}
		ret = append(ret, t)
	}
	return ret, nil
}

// Duration converts the underlying converted data into the final type.
// Uses internally time.ParseDuration
func (v Value) Duration() (d time.Duration, ok bool, err error) {
	if v, ok, err = v.init(); err != nil || !ok {
		return 0, false, errors.WithStack(err)
	}
	d, err = time.ParseDuration(string(v.data))
	ok = err == nil
	return
}

// IsEqual does a constant time comparison of the underlying data with the input
// data `d`. Useful for passwords. Only equal when no error occurs.
func (v Value) IsEqual(d []byte) bool {
	v, ok, err := v.init()
	return ok && subtle.ConstantTimeCompare(v.data, d) == 1 && err == nil
}

// IsEmpty returns true when its length is equal to zero.
func (v Value) IsEmpty() bool {
	v, _, err := v.init()
	return len(v.data) == 0 && err == nil
}

// IsValid returns the last error or in case when a Path is really not valid,
// returns an error with kind NotValid.
func (v Value) IsValid() bool {
	return v.lastErr == nil && v.found > valFoundNo
}

// equal compares if current is fully equal to v2.
func (v Value) equal(v2 Value) bool {
	return v.Path.Equal(v2.Path) && v.equalData(v2)
}

// equalData compares if the data part of the current value is equal to v2.
func (v Value) equalData(v2 Value) bool {
	return v.lastErr == nil && v2.lastErr == nil && v.found > valFoundNo && bytes.Equal(v.data, v2.data)
}

// Error implements error interface and returns the last error.
func (v Value) Error() string {
	if v.lastErr == nil {
		return ""
	}
	return v.lastErr.Error()
}

func parseDateTime(str string, loc *time.Location) (t time.Time, err error) {
	zeroBase := "0000-00-00 00:00:00.000000000+00:00"
	base := "2006-01-02 15:04:05.999999999 07:00"
	if strings.IndexByte(str, 'T') > 0 {
		base = time.RFC3339Nano
	}

	switch lStr := len(str); lStr {
	case 10, 19, 21, 22, 23, 24, 25, 26, 27, 28, 29, 35: // up to "YYYY-MM-DD HH:MM:SS.MMMMMMM+HH:II"
		if str == zeroBase[:lStr] {
			return
		}
		t, err = time.Parse(base[:lStr], str) // time.RFC3339Nano cannot be used due to the T
		if err != nil {
			err = errors.WithStack(err)
		}
	default:
		err = errors.NotValid.Newf("[config] invalid length %d in time string: %q", lStr, str)
		return
	}

	// Adjust location
	if err == nil && loc != time.UTC {
		y, mo, d := t.Date()
		h, mi, s := t.Clock()
		t, err = time.Date(y, mo, d, h, mi, s, t.Nanosecond(), loc), nil
	}
	return
}
