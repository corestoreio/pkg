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

package null_test

import (
	"encoding/json"
	"math"
	"testing"

	"github.com/corestoreio/pkg/storage/null"
	"github.com/corestoreio/pkg/util/assert"
	segjson "github.com/segmentio/encoding/json"
)

type AllTypes struct {
	Bool    null.Bool    `json:"bool"`
	Decimal null.Decimal `json:"decimal"`
	Float64 null.Float64 `json:"float64"`
	Int8    null.Int8    `json:"int8"`
	Int16   null.Int16   `json:"int16"`
	Int32   null.Int32   `json:"int32"`
	Int64   null.Int64   `json:"int64"`
	Uint8   null.Uint8   `json:"uint8"`
	Uint16  null.Uint16  `json:"uint16"`
	Uint32  null.Uint32  `json:"uint32"`
	Uint64  null.Uint64  `json:"uint64"`
	String  null.String  `json:"string"`
	// TODO time
}

func TestJSON_EncodingAll(t *testing.T) {
	at := AllTypes{
		Bool:    null.MakeBool(true),
		Decimal: null.Decimal{Precision: 987654321, Scale: 13, Negative: false, Valid: true, Quote: true},
		Float64: null.MakeFloat64(math.Pi),
		Int8:    null.MakeInt8(8),
		Int16:   null.MakeInt16(66),
		Int32:   null.MakeInt32(72),
		Int64:   null.MakeInt64(30),
		Uint8:   null.MakeUint8(62),
		Uint16:  null.MakeUint16(95),
		Uint32:  null.MakeUint32(8955412),
		Uint64:  null.MakeUint64(74),
		String:  null.MakeString(`OVRkYjbyamYIZiyBzBPQRGyLzyvEEbKjWMekMlSkbdNPFjIlVRAvjHMMlHwTeavthSaxjeoWuIoHNateBjTJhGranNSxPnezotCMnahKKQBqkOPqtMDeIZuvhRnHkWr`),
	}
	t.Run("packages", func(t *testing.T) {
		jsonEncDecEqual(t, json.Marshal, json.Unmarshal, &at, &AllTypes{})
		jsonEncDecEqual(t, segjson.Marshal, segjson.Unmarshal, &at, &AllTypes{})
		jsonEncDecEqual(t, json.Marshal, segjson.Unmarshal, &at, &AllTypes{})
		jsonEncDecEqual(t, segjson.Marshal, json.Unmarshal, &at, &AllTypes{})
	})

	t.Run("positive non-null", func(t *testing.T) {
		at := at // copy it
		at.Decimal.Quote = false
		jsonDecEqual(t, json.Unmarshal, []byte(`{
  "bool": true,
  "decimal": 0.0000987654321,
  "float64": 3.141592653589793,
  "int8": 8,
  "int16": 66,
  "int32": 72,
  "int64": 30,
  "uint8": 62,
  "uint16": 95,
  "uint32": 8955412,
  "uint64": 74,
  "string": "OVRkYjbyamYIZiyBzBPQRGyLzyvEEbKjWMekMlSkbdNPFjIlVRAvjHMMlHwTeavthSaxjeoWuIoHNateBjTJhGranNSxPnezotCMnahKKQBqkOPqtMDeIZuvhRnHkWr"
}`), &AllTypes{}, &at)
	})

	t.Run("negative non-null", func(t *testing.T) {
		at := AllTypes{
			Bool:    null.MakeBool(false),
			Decimal: null.Decimal{Precision: 987654321, Scale: 13, Negative: true, Valid: true, Quote: true},
			Float64: null.MakeFloat64(-math.Pi),
			Int8:    null.MakeInt8(-8),
			Int16:   null.MakeInt16(-66),
			Int32:   null.MakeInt32(-72),
			Int64:   null.MakeInt64(-30),
			Uint8:   null.MakeUint8(0),
			Uint16:  null.MakeUint16(0),
			Uint32:  null.MakeUint32(0),
			Uint64:  null.MakeUint64(0),
			String:  null.MakeString(`OVRkYjbyamYIZiyBzBPQRGyLzyvEEbKjWMekMlSkbdNPFjIlVRAvjHMMlHwTeavthSaxjeoWuIoHNateBjTJhGranNSxPnezotCMnahKKQBqkOPqtMDeIZuvhRnHkWr`),
		}

		jsonDecEqual(t, json.Unmarshal, []byte(`{
  "bool": false,
  "decimal": "-0.0000987654321",
  "float64": -3.141592653589793,
  "int8": -8,
  "int16": -66,
  "int32": -72,
  "int64": -30,
  "uint8": 0,
  "uint16": 0,
  "uint32": 0,
  "uint64": 0,
  "string": "OVRkYjbyamYIZiyBzBPQRGyLzyvEEbKjWMekMlSkbdNPFjIlVRAvjHMMlHwTeavthSaxjeoWuIoHNateBjTJhGranNSxPnezotCMnahKKQBqkOPqtMDeIZuvhRnHkWr"
}`), &AllTypes{}, &at)
	})

	t.Run("null", func(t *testing.T) {
		at := at // copy it
		at.Decimal.Quote = false
		jsonDecEqual(t, json.Unmarshal, []byte(`{
  "bool": null,
  "decimal": null,
  "float64": null,
  "int8": null,
  "int16": null,
  "int32": null,
  "int64": null,
  "uint8": null,
  "uint16": null,
  "uint32": null,
  "uint64": null,
  "string": null
}`), &AllTypes{}, &AllTypes{})
	})

	t.Run("empty", func(t *testing.T) {
		at := at // copy it
		at.Decimal.Quote = false
		jsonDecEqual(t, json.Unmarshal, []byte(`{}`), &AllTypes{}, &AllTypes{})
	})

	t.Run("positive sub structs non-null", func(t *testing.T) {
		at := at // copy it
		at.Decimal.Quote = false
		jsonDecEqual(t, json.Unmarshal, []byte(`{
  "bool": true,
  "decimal": 0.0000987654321,
  "float64": {"Float64":3.141592653589793,"Valid":true},
  "int8": {"Int8":8,"Valid":true},
  "int16": {"Int16":66,"Valid":true},
  "int32": {"Int32":72,"Valid":true},
  "int64": {"Int64":30,"Valid":true},
  "uint8": {"Uint8":62,"Valid":true},
  "uint16": {"Uint16":95,"Valid":true},
  "uint32": {"Uint32":8955412,"Valid":true},
  "uint64": {"Uint64":74,"Valid":true},
  "string": "OVRkYjbyamYIZiyBzBPQRGyLzyvEEbKjWMekMlSkbdNPFjIlVRAvjHMMlHwTeavthSaxjeoWuIoHNateBjTJhGranNSxPnezotCMnahKKQBqkOPqtMDeIZuvhRnHkWr"
}`), &AllTypes{}, &at)
	})
}

func jsonEncDecEqual(
	t *testing.T,
	jsonMarshalFn func(v interface{}) ([]byte, error),
	jsonUnMarshalFn func(data []byte, v interface{}) error,
	in interface{},
	out interface{},
) {
	data, err := jsonMarshalFn(in)
	assert.NoError(t, err)

	// println(string(data))

	assert.NoError(t, jsonUnMarshalFn(data, out))
	assert.Exactly(t, in, out)
}

func jsonDecEqual(
	t *testing.T,
	jsonUnMarshalFn func(data []byte, v interface{}) error,
	data []byte,
	out interface{},
	want interface{},
) {
	assert.NoError(t, jsonUnMarshalFn(data, out))
	assert.Exactly(t, want, out)
}
