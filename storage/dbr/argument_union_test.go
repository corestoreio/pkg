// Copyright 2015-2017, Cyrill @ Schumacher.fm and the CoreStore contributors
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

package dbr

import (
	"database/sql/driver"
	"github.com/corestoreio/errors"
	"math"
	"reflect"
	"testing"
	"time"
)

func BenchmarkArgUnion(b *testing.B) {
	reflectIFaceContainer := make([]interface{}, 0, 25)
	var finalArgs = make([]interface{}, 0, 30)
	drvVal := []driver.Valuer{MakeNullString("I'm a valid null string: See the License for the specific language governing permissions and See the License for the specific language governing permissions and See the License for the specific language governing permissions and")}
	argUnion := makeArgUninons(30)
	now1 := now()
	b.ResetTimer()

	b.Run("reflection all types", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			reflectIFaceContainer = append(reflectIFaceContainer,
				int64(5), []int64{6, 7, 8},
				uint64(9), []uint64{10, 11, 12},
				float64(3.14159), []float64{33.44, 55.66, 77.88},
				true, []bool{true, false, true},
				`Licensed under the Apache License, Version 2.0 (the "License");`,
				[]string{`Unless required by applicable law or agreed to in writing, software`, `Licensed under the Apache License, Version 2.0 (the "License");`},
				drvVal[0],
				nil,
				now1,
			)
			var err error
			finalArgs, err = encodePlaceholder(finalArgs, reflectIFaceContainer)
			if err != nil {
				b.Fatal(err)
			}
			//b.Fatal("%#v", finalArgs)
			reflectIFaceContainer = reflectIFaceContainer[:0]
			finalArgs = finalArgs[:0]
		}
	})
	b.Run("ArgUninons all types", func(b *testing.B) {
		// two times faster than the reflection version

		finalArgs = finalArgs[:0]

		for i := 0; i < b.N; i++ {
			argUnion = argUnion.
				Int64(5).Int64s(6, 7, 8).
				Uint64(9).Uint64s(10, 11, 12).
				Float64(3.14159).Float64s(33.44, 55.66, 77.88).
				Bool(true).Bools(true, false, true).
				String(`Licensed under the Apache License, Version 2.0 (the "License");`).
				Strings(`Unless required by applicable law or agreed to in writing, software`, `Licensed under the Apache License, Version 2.0 (the "License");`).
				DriverValue(drvVal...).
				Null().
				Time(now1)

			finalArgs = argUnion.Interfaces(finalArgs...)
			//b.Fatal("%#v", finalArgs)
			argUnion = argUnion[:0]
			finalArgs = finalArgs[:0]
		}
	})

	b.Run("reflection numbers", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			reflectIFaceContainer = append(reflectIFaceContainer,
				int64(5), []int64{6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19},
				uint64(9), []uint64{10, 11, 12, 20, 21, 22, 23, 24, 25, 26, 27, 28, 29},
				float64(3.14159), []float64{33.44, 55.66, 77.88, 11.22, math.Pi, math.E, math.Sqrt2},
				nil,
			)
			var err error
			finalArgs, err = encodePlaceholder(finalArgs, reflectIFaceContainer)
			if err != nil {
				b.Fatal(err)
			}
			//b.Fatal("%#v", finalArgs)
			reflectIFaceContainer = reflectIFaceContainer[:0]
			finalArgs = finalArgs[:0]
		}
	})
	b.Run("ArgUninons numbers", func(b *testing.B) {
		// three times faster than the reflection version

		finalArgs = finalArgs[:0]
		for i := 0; i < b.N; i++ {
			argUnion = argUnion.
				Int64(5).Int64s(6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19).
				Uint64(9).Uint64s(10, 11, 12, 20, 21, 22, 23, 24, 25, 26, 27, 28, 29).
				Float64(3.14159).Float64s(33.44, 55.66, 77.88, 11.22, math.Pi, math.E, math.Sqrt2).
				Null()

			finalArgs = argUnion.Interfaces(finalArgs...)
			//b.Fatal("%#v", finalArgs)
			argUnion = argUnion[:0]
			finalArgs = finalArgs[:0]
		}
	})

}

func encodePlaceholder(args []interface{}, value interface{}) ([]interface{}, error) {

	if valuer, ok := value.(driver.Valuer); ok {
		// get driver.Valuer's data
		var err error
		value, err = valuer.Value()
		if err != nil {
			return args, err
		}
	}

	if value == nil {
		return append(args, nil), nil
	}
	v := reflect.ValueOf(value)
	switch v.Kind() {
	case reflect.String:
		return append(args, v.String()), nil
	case reflect.Bool:
		return append(args, v.Bool()), nil
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return append(args, v.Int()), nil
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return append(args, v.Uint()), nil
	case reflect.Float32, reflect.Float64:
		return append(args, v.Float()), nil
	case reflect.Struct:
		if v.Type() == reflect.TypeOf(time.Time{}) {
			return append(args, v.Interface().(time.Time)), nil
		}
	case reflect.Slice:
		if v.Type().Elem().Kind() == reflect.Uint8 {
			// []byte
			return append(args, v.Bytes()), nil
		}
		if v.Len() == 0 {
			// FIXME: support zero-length slice
			return args, errors.NewNotValidf("invalid slice length")
		}

		for n := 0; n < v.Len(); n++ {
			var err error
			// recursion
			args, err = encodePlaceholder(args, v.Index(n).Interface())
			if err != nil {
				return args, err
			}
		}
		return args, nil
	case reflect.Ptr:
		if v.IsNil() {
			return append(args, nil), nil
		}
		return encodePlaceholder(args, v.Elem().Interface())

	}
	return args, errors.NewNotSupportedf("Type %#v not supported", value)
}
