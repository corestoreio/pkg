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
	"github.com/stretchr/testify/assert"
)

var (
	_ fmt.Stringer = (*Value)(nil)
	_ io.WriterTo  = (*Value)(nil)
)

func TestValue(t *testing.T) {
	t.Parallel()

	t.Run("String", func(t *testing.T) {
		v := MakeValue([]byte(`Rothaus`))
		assert.Exactly(t, "\"Rothaus\"", v.String())
	})
	t.Run("String Convert Failed", func(t *testing.T) {
		v := MakeValue([]byte(`Rothaus`))

		assert.Contains(t, v.WithConvert(func(p Path, data []byte) ([]byte, error) {
			return nil, errors.AlreadyInUse.Newf("Convert already in use")
		}).String(), "[config] Value: Convert already in use")
	})

	t.Run("WriteTo", func(t *testing.T) {
		v := MakeValue([]byte(`Rothaus Beer`))
		var buf strings.Builder
		_, err := v.WriteTo(&buf)
		assert.NoError(t, err)
		assert.Exactly(t, "Rothaus Beer", buf.String())
	})

	t.Run("Str1", func(t *testing.T) {
		v := MakeValue([]byte(`Waldhaus Beer`))
		val, ok, err := v.Str()
		assert.True(t, ok)
		assert.NoError(t, err)
		assert.Exactly(t, "Waldhaus Beer", string(val))
	})
	t.Run("Str2", func(t *testing.T) {
		v := MakeValue(nil)
		val, ok, err := v.Str()
		assert.False(t, ok)
		assert.NoError(t, err)
		assert.Exactly(t, "", string(val))
	})

	t.Run("Strs1", func(t *testing.T) {
		v := MakeValue([]byte(`SitUps,AirSquats,PushUps`))
		val, err := v.Strs()
		assert.NoError(t, err)
		assert.Exactly(t, []string{"SitUps", "AirSquats", "PushUps"}, val)
	})
	t.Run("Strs2", func(t *testing.T) {
		v := MakeValue([]byte(`SitUps`))
		v.CSVComma = ''
		val, err := v.Strs()
		assert.NoError(t, err)
		assert.Exactly(t, []string{"SitUps"}, val)
	})
	t.Run("Strs3", func(t *testing.T) {
		v := MakeValue([]byte(`SitUps`))
		v.CSVComma = ''
		val := []string{"X"}
		val, err := v.Strs(val...)
		assert.NoError(t, err)
		assert.Exactly(t, []string{"X", "SitUps"}, val)
	})
	t.Run("Strs4", func(t *testing.T) {
		v := MakeValue(nil)
		val := []string{"X"}
		val, err := v.Strs(val...)
		assert.NoError(t, err)
		assert.Exactly(t, []string{"X"}, val)
	})
	t.Run("Strs5", func(t *testing.T) {
		v := MakeValue([]byte(`SitUps,,DU`))
		val, err := v.Strs()
		assert.NoError(t, err)
		assert.Exactly(t, []string{"SitUps", "DU"}, val)
	})
	t.Run("CSV1", func(t *testing.T) {
		v := MakeValue([]byte(`SitUps`))
		val := [][]string{{"X"}, {"Y"}}
		val, err := v.CSV(val...)
		assert.NoError(t, err)
		assert.Exactly(t, [][]string{{"X"}, {"Y"}, {"SitUps"}}, val)
	})
	t.Run("CSV2", func(t *testing.T) {
		v := MakeValue([]byte(`50xSitUps,21xHSPU`))
		val := [][]string{{"X"}, {"Y"}}
		val, err := v.CSV(val...)
		assert.NoError(t, err)
		assert.Exactly(t, [][]string{{"X"}, {"Y"}, {"50xSitUps", "21xHSPU"}}, val)
	})
	t.Run("CSV3", func(t *testing.T) {
		v := MakeValue([]byte("50xSitUps,21xHSPU\n18xBar MU,9xHPC"))
		var val [][]string
		val, err := v.CSV(val...)
		assert.NoError(t, err)
		assert.Exactly(t, [][]string{{"50xSitUps", "21xHSPU"}, {"18xBar MU", "9xHPC"}}, val)
	})
	t.Run("CSV4", func(t *testing.T) {
		v := MakeValue(nil)
		var val [][]string
		val, err := v.CSV(val...)
		assert.NoError(t, err)
		assert.Exactly(t, [][]string(nil), val)
	})

	t.Run("Unmarshal1", func(t *testing.T) {
		v := MakeValue([]byte(`{"X":1}`))
		val := map[string]int{}
		err := v.Unmarshal(json.Unmarshal, &val)
		assert.NoError(t, err)
		assert.Exactly(t, map[string]int{"X": 1}, val)
	})

	t.Run("Bool1", func(t *testing.T) {
		v := MakeValue([]byte(`true`))
		val, ok, err := v.Bool()
		assert.NoError(t, err)
		assert.True(t, ok)
		assert.Exactly(t, true, val)
	})
	t.Run("Bool2", func(t *testing.T) {
		v := MakeValue([]byte(`tru3`))
		val, ok, err := v.Bool()
		assert.Error(t, err)
		assert.False(t, ok)
		assert.Exactly(t, false, val)
	})
	t.Run("Bool3", func(t *testing.T) {
		v := MakeValue(nil)
		val, ok, err := v.Bool()
		assert.NoError(t, err)
		assert.False(t, ok)
		assert.Exactly(t, false, val)
	})

	t.Run("Float641", func(t *testing.T) {
		v := MakeValue([]byte(`-3.14159`))
		val, ok, err := v.Float64()
		assert.NoError(t, err)
		assert.True(t, ok)
		assert.Exactly(t, -3.14159, val)
	})
	t.Run("Float642", func(t *testing.T) {
		v := MakeValue([]byte(`tru3`))
		val, ok, err := v.Float64()
		assert.Error(t, err)
		assert.False(t, ok)
		assert.Exactly(t, float64(0), val)
	})
	t.Run("Float643", func(t *testing.T) {
		v := MakeValue(nil)
		val, ok, err := v.Float64()
		assert.NoError(t, err)
		assert.False(t, ok)
		assert.Exactly(t, float64(0), val)
	})

	t.Run("Float64s1", func(t *testing.T) {
		v := MakeValue([]byte(`3,2.7182,-21.15`))
		val, err := v.Float64s()
		assert.NoError(t, err)
		assert.Exactly(t, []float64{3, 2.7182, -21.15}, val)
	})
	t.Run("Float64s2", func(t *testing.T) {
		v := MakeValue([]byte(`3.33`))
		v.CSVComma = ''
		val, err := v.Float64s()
		assert.NoError(t, err)
		assert.Exactly(t, []float64{3.33}, val)
	})
	t.Run("Float64s3", func(t *testing.T) {
		v := MakeValue([]byte(`-0.01`))
		v.CSVComma = ''
		val := []float64{0.01}
		val, err := v.Float64s(val...)
		assert.NoError(t, err)
		assert.Exactly(t, []float64{0.01, -0.01}, val)
	})
	t.Run("Float64s4", func(t *testing.T) {
		v := MakeValue(nil)
		val := []float64{11}
		val, err := v.Float64s(val...)
		assert.NoError(t, err)
		assert.Exactly(t, []float64{11}, val)
	})
	t.Run("Float64s5", func(t *testing.T) {
		v := MakeValue([]byte(`3,X,-21.15`))
		val, err := v.Float64s()
		assert.Nil(t, val)
		assert.EqualError(t, err, "[config] Value.Float64s with index 1 and entry \"X\": strconv.ParseFloat: parsing \"X\": invalid syntax")
	})

	t.Run("Int1", func(t *testing.T) {
		v := MakeValue([]byte(`-314159`))
		val, ok, err := v.Int()
		assert.NoError(t, err)
		assert.True(t, ok)
		assert.Exactly(t, int(-314159), val)
	})
	t.Run("Int2", func(t *testing.T) {
		v := MakeValue([]byte(`tru3`))
		val, ok, err := v.Int()
		assert.Error(t, err)
		assert.False(t, ok)
		assert.Exactly(t, int(0), val)
	})
	t.Run("Int3", func(t *testing.T) {
		v := MakeValue(nil)
		val, ok, err := v.Int()
		assert.NoError(t, err)
		assert.False(t, ok)
		assert.Exactly(t, int(0), val)
	})

	t.Run("Ints1", func(t *testing.T) {
		v := MakeValue([]byte(`3,27182,-2115`))
		val, err := v.Ints()
		assert.NoError(t, err)
		assert.Exactly(t, []int{3, 27182, -2115}, val)
	})
	t.Run("Ints2", func(t *testing.T) {
		v := MakeValue([]byte(`333`))
		v.CSVComma = ''
		val, err := v.Ints()
		assert.NoError(t, err)
		assert.Exactly(t, []int{333}, val)
	})
	t.Run("Ints3", func(t *testing.T) {
		v := MakeValue([]byte(`-1`))
		v.CSVComma = ''
		val := []int{1}
		val, err := v.Ints(val...)
		assert.NoError(t, err)
		assert.Exactly(t, []int{1, -1}, val)
	})
	t.Run("Ints4", func(t *testing.T) {
		v := MakeValue(nil)
		val := []int{11}
		val, err := v.Ints(val...)
		assert.NoError(t, err)
		assert.Exactly(t, []int{11}, val)
	})
	t.Run("Ints5", func(t *testing.T) {
		v := MakeValue([]byte(`3,X,-2115`))
		val, err := v.Ints()
		assert.Nil(t, val)
		assert.EqualError(t, err, "[config] Value.Ints with index 1 and entry \"X\": strconv.ParseInt: parsing \"X\": invalid syntax")
	})

	t.Run("Int641", func(t *testing.T) {
		v := MakeValue([]byte(`-314159`))
		val, ok, err := v.Int64()
		assert.NoError(t, err)
		assert.True(t, ok)
		assert.Exactly(t, int64(-314159), val)
	})
	t.Run("Int642", func(t *testing.T) {
		v := MakeValue([]byte(`tru3`))
		val, ok, err := v.Int64()
		assert.Error(t, err)
		assert.False(t, ok)
		assert.Exactly(t, int64(0), val)
	})
	t.Run("Int643", func(t *testing.T) {
		v := MakeValue(nil)
		val, ok, err := v.Int64()
		assert.NoError(t, err)
		assert.False(t, ok)
		assert.Exactly(t, int64(0), val)
	})

	t.Run("Int64s1", func(t *testing.T) {
		v := MakeValue([]byte(`3,27182,-2115`))
		val, err := v.Int64s()
		assert.NoError(t, err)
		assert.Exactly(t, []int64{3, 27182, -2115}, val)
	})
	t.Run("Int64s2", func(t *testing.T) {
		v := MakeValue([]byte(`333`))
		v.CSVComma = ''
		val, err := v.Int64s()
		assert.NoError(t, err)
		assert.Exactly(t, []int64{333}, val)
	})
	t.Run("Int64s3", func(t *testing.T) {
		v := MakeValue([]byte(`-1`))
		v.CSVComma = ''
		val := []int64{1}
		val, err := v.Int64s(val...)
		assert.NoError(t, err)
		assert.Exactly(t, []int64{1, -1}, val)
	})
	t.Run("Int64s4", func(t *testing.T) {
		v := MakeValue(nil)
		val := []int64{11}
		val, err := v.Int64s(val...)
		assert.NoError(t, err)
		assert.Exactly(t, []int64{11}, val)
	})
	t.Run("Int64s5", func(t *testing.T) {
		v := MakeValue([]byte(`3,X,-2115`))
		val, err := v.Int64s()
		assert.Nil(t, val)
		assert.EqualError(t, err, "[config] Value.Int64s with index 1 and entry \"X\": strconv.ParseInt: parsing \"X\": invalid syntax")
	})

	t.Run("Uint641", func(t *testing.T) {
		v := MakeValue([]byte(`314159`))
		val, ok, err := v.Uint64()
		assert.NoError(t, err)
		assert.True(t, ok)
		assert.Exactly(t, uint64(314159), val)
	})
	t.Run("Uint642", func(t *testing.T) {
		v := MakeValue([]byte(`tru3`))
		val, ok, err := v.Uint64()
		assert.Error(t, err)
		assert.False(t, ok)
		assert.Exactly(t, uint64(0), val)
	})
	t.Run("Uint643", func(t *testing.T) {
		v := MakeValue(nil)
		val, ok, err := v.Uint64()
		assert.NoError(t, err)
		assert.False(t, ok)
		assert.Exactly(t, uint64(0), val)
	})

	t.Run("Uint64s1", func(t *testing.T) {
		v := MakeValue([]byte(`3,27182,2115`))
		val, err := v.Uint64s()
		assert.NoError(t, err)
		assert.Exactly(t, []uint64{3, 27182, 2115}, val)
	})
	t.Run("Uint64s2", func(t *testing.T) {
		v := MakeValue([]byte(`333`))
		v.CSVComma = ''
		val, err := v.Uint64s()
		assert.NoError(t, err)
		assert.Exactly(t, []uint64{333}, val)
	})
	t.Run("Uint64s3", func(t *testing.T) {
		v := MakeValue([]byte(`2`))
		v.CSVComma = ''
		val := []uint64{1}
		val, err := v.Uint64s(val...)
		assert.NoError(t, err)
		assert.Exactly(t, []uint64{1, 2}, val)
	})
	t.Run("Uint64s4", func(t *testing.T) {
		v := MakeValue(nil)
		val := []uint64{11}
		val, err := v.Uint64s(val...)
		assert.NoError(t, err)
		assert.Exactly(t, []uint64{11}, val)
	})
	t.Run("Uint64s5", func(t *testing.T) {
		v := MakeValue([]byte(`3,X,-2115`))
		val, err := v.Uint64s()
		assert.Nil(t, val)
		assert.EqualError(t, err, "[config] Value.Uint64s with index 1 and entry \"X\": strconv.ParseUint: parsing \"X\": invalid syntax")
	})

	t.Run("Time1", func(t *testing.T) {
		v := MakeValue([]byte(`2018-04-02`))
		val, err := v.Time()
		assert.NoError(t, err)
		assert.Exactly(t, "2018-04-02 00:00:00 +0000 UTC", val.String())
	})
	t.Run("Time2", func(t *testing.T) {
		ct := time.Now().Format("2006-01-02 15:04:05.999999999")
		v := MakeValue([]byte(ct))
		val, err := v.Time()
		assert.NoError(t, err)
		assert.Exactly(t, ct+" +0000 UTC", val.String())
	})
	t.Run("Time3", func(t *testing.T) {
		v := MakeValue([]byte(`X018-04-02`))
		val, err := v.Time()
		assert.EqualError(t, err, "parsing time \"X018-04-02\" as \"2006-01-02\": cannot parse \"X018-04-02\" as \"2006\"")
		assert.Exactly(t, time.Time{}, val)
	})

	t.Run("Times1", func(t *testing.T) {
		v := MakeValue([]byte(`2018-04-02,2018-04-02,`))
		val, err := v.Times()
		assert.NoError(t, err)
		assert.Exactly(t, "2018-04-02 00:00:00 +0000 UTC", val[0].String())
		assert.Exactly(t, "2018-04-02 00:00:00 +0000 UTC", val[1].String())
	})
	t.Run("Times2", func(t *testing.T) {
		v := MakeValue([]byte(`2018-04-02,2018-X4-02,`))
		val, err := v.Times()
		assert.Nil(t, val)
		assert.EqualError(t, err, "[config] Value.Times with index 1 and entry \"2018-X4-02\": parsing time \"2018-X4-02\": month out of range")
	})
	t.Run("Times3", func(t *testing.T) {
		v := MakeValue([]byte(`2018-04-02,,2018-X4-02`))
		val, err := v.Times()
		assert.Nil(t, val)
		assert.EqualError(t, err, "[config] Value.Times with index 2 and entry \"2018-X4-02\": parsing time \"2018-X4-02\": month out of range")
	})

	t.Run("Duration", func(t *testing.T) {
		v := MakeValue([]byte(`5m2s`))
		val, err := v.Duration()
		assert.NoError(t, err)
		assert.Exactly(t, "5m2s", val.String())
	})

	t.Run("IsEqual", func(t *testing.T) {
		d := []byte(`5m2s`)
		v := MakeValue(d)
		val, err := v.IsEqual(d)
		assert.NoError(t, err)
		assert.True(t, val)
	})
}
