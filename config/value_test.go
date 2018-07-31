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
	"encoding/json"
	"fmt"
	"io"
	"strings"
	"testing"
	"time"

	"github.com/corestoreio/errors"
	"github.com/corestoreio/pkg/util/assert"
)

var (
	_ fmt.Stringer = (*Value)(nil)
	_ io.WriterTo  = (*Value)(nil)
	_ error        = (*Value)(nil)
)

type myTestUnmarshal struct {
	result string
}

func (m *myTestUnmarshal) UnmarshalText(text []byte) error {
	m.result = string(text)
	return nil
}
func (m *myTestUnmarshal) UnmarshalBinary(data []byte) error {
	m.result = string(data)
	return nil
}

func TestValue(t *testing.T) {
	t.Parallel()

	t.Run("String", func(t *testing.T) {
		v := NewValue([]byte(`Rothaus`))
		assert.Exactly(t, "\"Rothaus\"", v.String())

		v = NewValue([]byte(nil))
		assert.Exactly(t, "<nil>", v.String())

		v.found = valFoundNo
		assert.Exactly(t, "<notFound>", v.String())
	})

	t.Run("WriteTo", func(t *testing.T) {
		v := NewValue([]byte(`Rothaus Beer`))
		var buf strings.Builder
		_, err := v.WriteTo(&buf)
		assert.NoError(t, err)
		assert.Exactly(t, "Rothaus Beer", buf.String())
	})

	t.Run("Str1", func(t *testing.T) {
		v := NewValue([]byte(`Waldhaus Beer`))
		val, ok, err := v.Str()
		assert.True(t, ok)
		assert.NoError(t, err)
		assert.Exactly(t, "Waldhaus Beer", string(val))

		v.found = valFoundNo
		val, ok, err = v.Str()
		assert.False(t, ok)
		assert.NoError(t, err)
		assert.Exactly(t, "", string(val))
	})
	t.Run("Str2", func(t *testing.T) {
		v := NewValue(nil)
		val, ok, err := v.Str()
		assert.False(t, ok)
		assert.NoError(t, err)
		assert.Exactly(t, "", string(val))
	})
	t.Run("Error", func(t *testing.T) {
		v := NewValue(nil)
		v.lastErr = errors.New("Ups")
		assert.EqualError(t, v, "Ups")
		v.lastErr = nil
		assert.EqualError(t, v, "")
	})

	t.Run("Strs1", func(t *testing.T) {
		v := NewValue([]byte(`SitUps,AirSquats,PushUps`))
		val, err := v.Strs()
		assert.NoError(t, err)
		assert.Exactly(t, []string{"SitUps", "AirSquats", "PushUps"}, val)
	})
	t.Run("Strs2", func(t *testing.T) {
		v := NewValue([]byte(`SitUps`))
		v.CSVComma = ''
		val, err := v.Strs()
		assert.NoError(t, err)
		assert.Exactly(t, []string{"SitUps"}, val)
	})
	t.Run("Strs3", func(t *testing.T) {
		v := NewValue([]byte(`SitUps`))
		v.CSVComma = ''
		val := []string{"X"}
		val, err := v.Strs(val...)
		assert.NoError(t, err)
		assert.Exactly(t, []string{"X", "SitUps"}, val)
	})
	t.Run("Strs4", func(t *testing.T) {
		v := NewValue(nil)
		val := []string{"X"}
		val, err := v.Strs(val...)
		assert.NoError(t, err)
		assert.Exactly(t, []string{"X"}, val)
	})
	t.Run("Strs5", func(t *testing.T) {
		v := NewValue([]byte(`SitUps,,DU`))
		val, err := v.Strs()
		assert.NoError(t, err)
		assert.Exactly(t, []string{"SitUps", "DU"}, val)
	})
	t.Run("CSV1", func(t *testing.T) {
		v := NewValue([]byte(`SitUps`))
		val := [][]string{{"X"}, {"Y"}}
		val, err := v.CSV(val...)
		assert.NoError(t, err)
		assert.Exactly(t, [][]string{{"X"}, {"Y"}, {"SitUps"}}, val)
	})
	t.Run("CSV2", func(t *testing.T) {
		v := NewValue([]byte(`50xSitUps,21xHSPU`))
		val := [][]string{{"X"}, {"Y"}}
		val, err := v.CSV(val...)
		assert.NoError(t, err)
		assert.Exactly(t, [][]string{{"X"}, {"Y"}, {"50xSitUps", "21xHSPU"}}, val)
	})
	t.Run("CSV3", func(t *testing.T) {
		v := NewValue([]byte("50xSitUps,21xHSPU\n18xBar MU,9xHPC"))
		var val [][]string
		val, err := v.CSV(val...)
		assert.NoError(t, err)
		assert.Exactly(t, [][]string{{"50xSitUps", "21xHSPU"}, {"18xBar MU", "9xHPC"}}, val)
	})
	t.Run("CSV4", func(t *testing.T) {
		v := NewValue(nil)
		var val [][]string
		val, err := v.CSV(val...)
		assert.NoError(t, err)
		assert.Exactly(t, [][]string(nil), val)
	})

	t.Run("UnmarshalTo", func(t *testing.T) {
		v := NewValue([]byte(`{"X":1}`))
		val := map[string]int{}
		err := v.UnmarshalTo(json.Unmarshal, &val)
		assert.NoError(t, err)
		assert.Exactly(t, map[string]int{"X": 1}, val)
	})

	t.Run("UnmarshalTextTo", func(t *testing.T) {
		v := NewValue([]byte(`{zXz:1}`))
		val := &myTestUnmarshal{}
		err := v.UnmarshalTextTo(val)
		assert.NoError(t, err)
		assert.Exactly(t, `{zXz:1}`, val.result)
	})

	t.Run("UnmarshalBinaryTo", func(t *testing.T) {
		v := NewValue([]byte(`{eXe:1}`))
		val := &myTestUnmarshal{}
		err := v.UnmarshalBinaryTo(val)
		assert.NoError(t, err)
		assert.Exactly(t, `{eXe:1}`, val.result)
	})

	t.Run("Bool1", func(t *testing.T) {
		v := NewValue([]byte(`true`))
		val, ok, err := v.Bool()
		assert.NoError(t, err)
		assert.True(t, ok)
		assert.Exactly(t, true, val)
	})
	t.Run("Bool2", func(t *testing.T) {
		v := NewValue([]byte(`tru3`))
		val, ok, err := v.Bool()
		assert.Error(t, err)
		assert.False(t, ok)
		assert.Exactly(t, false, val)
	})
	t.Run("Bool3", func(t *testing.T) {
		v := NewValue(nil)
		val, ok, err := v.Bool()
		assert.NoError(t, err)
		assert.False(t, ok)
		assert.Exactly(t, false, val)
	})

	t.Run("Float641", func(t *testing.T) {
		v := NewValue([]byte(`-3.14159`))
		val, ok, err := v.Float64()
		assert.NoError(t, err)
		assert.True(t, ok)
		assert.Exactly(t, -3.14159, val)
	})
	t.Run("Float642", func(t *testing.T) {
		v := NewValue([]byte(`tru3`))
		val, ok, err := v.Float64()
		assert.Error(t, err)
		assert.False(t, ok)
		assert.Exactly(t, float64(0), val)
	})
	t.Run("Float643", func(t *testing.T) {
		v := NewValue(nil)
		val, ok, err := v.Float64()
		assert.NoError(t, err)
		assert.False(t, ok)
		assert.Exactly(t, float64(0), val)
	})

	t.Run("Float64s1", func(t *testing.T) {
		v := NewValue([]byte(`3,2.7182,-21.15`))
		val, err := v.Float64s()
		assert.NoError(t, err)
		assert.Exactly(t, []float64{3, 2.7182, -21.15}, val)
	})
	t.Run("Float64s2", func(t *testing.T) {
		v := NewValue([]byte(`3.33`))
		v.CSVComma = ''
		val, err := v.Float64s()
		assert.NoError(t, err)
		assert.Exactly(t, []float64{3.33}, val)
	})
	t.Run("Float64s3", func(t *testing.T) {
		v := NewValue([]byte(`-0.01`))
		v.CSVComma = ''
		val := []float64{0.01}
		val, err := v.Float64s(val...)
		assert.NoError(t, err)
		assert.Exactly(t, []float64{0.01, -0.01}, val)
	})
	t.Run("Float64s4", func(t *testing.T) {
		v := NewValue(nil)
		val := []float64{11}
		val, err := v.Float64s(val...)
		assert.NoError(t, err)
		assert.Exactly(t, []float64{11}, val)
	})
	t.Run("Float64s5", func(t *testing.T) {
		v := NewValue([]byte(`3,X,-21.15`))
		val, err := v.Float64s()
		assert.Nil(t, val)
		assert.EqualError(t, err, "[config] Value.Float64s with index 1 and entry \"X\": strconv.ParseFloat: parsing \"X\": invalid syntax")
	})

	t.Run("Int1", func(t *testing.T) {
		v := NewValue([]byte(`-314159`))
		val, ok, err := v.Int()
		assert.NoError(t, err)
		assert.True(t, ok)
		assert.Exactly(t, int(-314159), val)
	})
	t.Run("Int2", func(t *testing.T) {
		v := NewValue([]byte(`tru3`))
		val, ok, err := v.Int()
		assert.Error(t, err)
		assert.False(t, ok)
		assert.Exactly(t, int(0), val)
	})
	t.Run("Int3", func(t *testing.T) {
		v := NewValue(nil)
		val, ok, err := v.Int()
		assert.NoError(t, err)
		assert.False(t, ok)
		assert.Exactly(t, int(0), val)
	})

	t.Run("Ints1", func(t *testing.T) {
		v := NewValue([]byte(`3,27182,-2115`))
		val, err := v.Ints()
		assert.NoError(t, err)
		assert.Exactly(t, []int{3, 27182, -2115}, val)
	})
	t.Run("Ints2", func(t *testing.T) {
		v := NewValue([]byte(`333`))
		v.CSVComma = ''
		val, err := v.Ints()
		assert.NoError(t, err)
		assert.Exactly(t, []int{333}, val)
	})
	t.Run("Ints3", func(t *testing.T) {
		v := NewValue([]byte(`-1`))
		v.CSVComma = ''
		val := []int{1}
		val, err := v.Ints(val...)
		assert.NoError(t, err)
		assert.Exactly(t, []int{1, -1}, val)
	})
	t.Run("Ints4", func(t *testing.T) {
		v := NewValue(nil)
		val := []int{11}
		val, err := v.Ints(val...)
		assert.NoError(t, err)
		assert.Exactly(t, []int{11}, val)
	})
	t.Run("Ints5", func(t *testing.T) {
		v := NewValue([]byte(`3,X,-2115`))
		val, err := v.Ints()
		assert.Nil(t, val)
		assert.EqualError(t, err, "[config] Value.Ints with index 1 and entry \"X\": strconv.ParseInt: parsing \"X\": invalid syntax")
	})

	t.Run("Int641", func(t *testing.T) {
		v := NewValue([]byte(`-314159`))
		val, ok, err := v.Int64()
		assert.NoError(t, err)
		assert.True(t, ok)
		assert.Exactly(t, int64(-314159), val)
	})
	t.Run("Int642", func(t *testing.T) {
		v := NewValue([]byte(`tru3`))
		val, ok, err := v.Int64()
		assert.Error(t, err)
		assert.False(t, ok)
		assert.Exactly(t, int64(0), val)
	})
	t.Run("Int643", func(t *testing.T) {
		v := NewValue(nil)
		val, ok, err := v.Int64()
		assert.NoError(t, err)
		assert.False(t, ok)
		assert.Exactly(t, int64(0), val)
	})

	t.Run("Int64s1", func(t *testing.T) {
		v := NewValue([]byte(`3,27182,-2115`))
		val, err := v.Int64s()
		assert.NoError(t, err)
		assert.Exactly(t, []int64{3, 27182, -2115}, val)
	})
	t.Run("Int64s2", func(t *testing.T) {
		v := NewValue([]byte(`333`))
		v.CSVComma = ''
		val, err := v.Int64s()
		assert.NoError(t, err)
		assert.Exactly(t, []int64{333}, val)
	})
	t.Run("Int64s3", func(t *testing.T) {
		v := NewValue([]byte(`-1`))
		v.CSVComma = ''
		val := []int64{1}
		val, err := v.Int64s(val...)
		assert.NoError(t, err)
		assert.Exactly(t, []int64{1, -1}, val)
	})
	t.Run("Int64s4", func(t *testing.T) {
		v := NewValue(nil)
		val := []int64{11}
		val, err := v.Int64s(val...)
		assert.NoError(t, err)
		assert.Exactly(t, []int64{11}, val)
	})
	t.Run("Int64s5", func(t *testing.T) {
		v := NewValue([]byte(`3,X,-2115`))
		val, err := v.Int64s()
		assert.Nil(t, val)
		assert.EqualError(t, err, "[config] Value.Int64s with index 1 and entry \"X\": strconv.ParseInt: parsing \"X\": invalid syntax")
	})

	t.Run("Uint641", func(t *testing.T) {
		v := NewValue([]byte(`314159`))
		val, ok, err := v.Uint64()
		assert.NoError(t, err)
		assert.True(t, ok)
		assert.Exactly(t, uint64(314159), val)
	})
	t.Run("Uint642", func(t *testing.T) {
		v := NewValue([]byte(`tru3`))
		val, ok, err := v.Uint64()
		assert.Error(t, err)
		assert.False(t, ok)
		assert.Exactly(t, uint64(0), val)
	})
	t.Run("Uint643", func(t *testing.T) {
		v := NewValue(nil)
		val, ok, err := v.Uint64()
		assert.NoError(t, err)
		assert.False(t, ok)
		assert.Exactly(t, uint64(0), val)
	})

	t.Run("Uint64s1", func(t *testing.T) {
		v := NewValue([]byte(`3,27182,2115`))
		val, err := v.Uint64s()
		assert.NoError(t, err)
		assert.Exactly(t, []uint64{3, 27182, 2115}, val)
	})
	t.Run("Uint64s2", func(t *testing.T) {
		v := NewValue([]byte(`333`))
		v.CSVComma = ''
		val, err := v.Uint64s()
		assert.NoError(t, err)
		assert.Exactly(t, []uint64{333}, val)
	})
	t.Run("Uint64s3", func(t *testing.T) {
		v := NewValue([]byte(`2`))
		v.CSVComma = ''
		val := []uint64{1}
		val, err := v.Uint64s(val...)
		assert.NoError(t, err)
		assert.Exactly(t, []uint64{1, 2}, val)
	})
	t.Run("Uint64s4", func(t *testing.T) {
		v := NewValue(nil)
		val := []uint64{11}
		val, err := v.Uint64s(val...)
		assert.NoError(t, err)
		assert.Exactly(t, []uint64{11}, val)
	})
	t.Run("Uint64s5", func(t *testing.T) {
		v := NewValue([]byte(`3,X,-2115`))
		val, err := v.Uint64s()
		assert.Nil(t, val)
		assert.EqualError(t, err, "[config] Value.Uint64s with index 1 and entry \"X\": strconv.ParseUint: parsing \"X\": invalid syntax")
	})

	t.Run("Time1", func(t *testing.T) {
		v := NewValue([]byte(`2018-04-02`))
		val, ok, err := v.Time()
		assert.NoError(t, err)
		assert.True(t, ok, "Time should be set and not nil, so true.")
		assert.Exactly(t, "2018-04-02 00:00:00 +0000 UTC", val.String())
	})
	t.Run("Time2", func(t *testing.T) {
		ct := time.Now().Format("2006-01-02 15:04:05.999999999")
		v := NewValue([]byte(ct))
		val, ok, err := v.Time()
		assert.True(t, ok, "Time should be set and not nil, so true.")
		assert.NoError(t, err)
		assert.Exactly(t, ct+" +0000 UTC", val.String())
	})
	t.Run("Time3", func(t *testing.T) {
		v := NewValue([]byte(`X018-04-02`))
		val, ok, err := v.Time()
		assert.False(t, ok, "Time should NOT be set because invalid.")
		assert.EqualError(t, err, "parsing time \"X018-04-02\" as \"2006-01-02\": cannot parse \"X018-04-02\" as \"2006\"")
		assert.Exactly(t, time.Time{}, val)
	})

	t.Run("Times1", func(t *testing.T) {
		v := NewValue([]byte(`2018-04-02,2018-04-02,`))
		val, err := v.Times()
		assert.NoError(t, err)
		assert.Exactly(t, "2018-04-02 00:00:00 +0000 UTC", val[0].String())
		assert.Exactly(t, "2018-04-02 00:00:00 +0000 UTC", val[1].String())
	})
	t.Run("Times2", func(t *testing.T) {
		v := NewValue([]byte(`2018-04-02,2018-X4-02,`))
		val, err := v.Times()
		assert.Nil(t, val)
		assert.EqualError(t, err, "[config] Value.Times with index 1 and entry \"2018-X4-02\": parsing time \"2018-X4-02\": month out of range")
	})
	t.Run("Times3", func(t *testing.T) {
		v := NewValue([]byte(`2018-04-02,,2018-X4-02`))
		val, err := v.Times()
		assert.Nil(t, val)
		assert.EqualError(t, err, "[config] Value.Times with index 2 and entry \"2018-X4-02\": parsing time \"2018-X4-02\": month out of range")
	})

	t.Run("Duration", func(t *testing.T) {
		v := NewValue([]byte(`5m2s`))
		val, ok, err := v.Duration()
		assert.True(t, ok, "Duration should be set and not nil, so true.")
		assert.NoError(t, err)
		assert.Exactly(t, "5m2s", val.String())

		v.found = valFoundNo
		val, ok, err = v.Duration()
		assert.False(t, ok)
		assert.NoError(t, err)
		assert.Exactly(t, "0s", val.String())
	})

	t.Run("IsEqual", func(t *testing.T) {
		d := []byte(`5m2s`)
		v := NewValue(d)
		assert.True(t, v.IsEqual(d))
	})
	t.Run("ConstantTimeCompare", func(t *testing.T) {
		d := []byte(`5m2s`)
		v := NewValue(d)
		assert.True(t, v.ConstantTimeCompare(d))
	})
}

func TestValue_Unsafe(t *testing.T) {
	t.Parallel()
	t.Run("bool", func(t *testing.T) {
		assert.True(t, NewValue([]byte(`1`)).UnsafeBool())
		assert.False(t, NewValue([]byte(`0`)).UnsafeBool())
		assert.False(t, NewValue([]byte(``)).UnsafeBool())
		assert.False(t, NewValue(nil).UnsafeBool())
	})
	t.Run("string", func(t *testing.T) {
		assert.Exactly(t, `Ups`, NewValue([]byte(`Ups`)).UnsafeStr())
		assert.Exactly(t, ``, NewValue(nil).UnsafeStr())
	})
	t.Run("float64", func(t *testing.T) {
		assert.Exactly(t, 2.718281, NewValue([]byte(`2.718281`)).UnsafeFloat64())
		assert.Exactly(t, 0.0, NewValue([]byte(`=`)).UnsafeFloat64())
		assert.Exactly(t, 0.0, NewValue(nil).UnsafeFloat64())
	})
	t.Run("int", func(t *testing.T) {
		assert.Exactly(t, 2718281, NewValue([]byte(`2718281`)).UnsafeInt())
		assert.Exactly(t, 0, NewValue([]byte(`=`)).UnsafeInt())
		assert.Exactly(t, 0, NewValue(nil).UnsafeInt())
	})
	t.Run("int64", func(t *testing.T) {
		assert.Exactly(t, int64(2718281), NewValue([]byte(`2718281`)).UnsafeInt64())
		assert.Exactly(t, int64(0), NewValue([]byte(`=`)).UnsafeInt64())
		assert.Exactly(t, int64(0), NewValue(nil).UnsafeInt64())
	})
	t.Run("uint64", func(t *testing.T) {
		assert.Exactly(t, uint64(2718281), NewValue([]byte(`2718281`)).UnsafeUint64())
		assert.Exactly(t, uint64(0), NewValue([]byte(`=`)).UnsafeUint64())
		assert.Exactly(t, uint64(0), NewValue(nil).UnsafeUint64())
	})
	t.Run("time", func(t *testing.T) {
		assert.Exactly(t, time.Date(2018, 4, 2, 0, 0, 0, 0, time.UTC), NewValue([]byte(`2018-04-02`)).UnsafeTime())
		assert.Exactly(t, time.Time{}, NewValue([]byte(`12018-04-02`)).UnsafeTime())
		assert.Exactly(t, time.Time{}, NewValue(nil).UnsafeTime())
	})
	t.Run("duration", func(t *testing.T) {
		assert.Exactly(t, time.Duration(302000000000), NewValue([]byte(`5m2s`)).UnsafeDuration())
		assert.Exactly(t, time.Duration(0), NewValue(nil).UnsafeDuration())
		assert.Exactly(t, time.Duration(0), NewValue([]byte(`A`)).UnsafeDuration())
	})
}
