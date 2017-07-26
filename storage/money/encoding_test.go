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

package money_test

import (
	"bytes"
	"context"
	"database/sql"
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"testing"

	"github.com/corestoreio/csfw/storage/dbr"
	"github.com/corestoreio/csfw/storage/money"
	"github.com/corestoreio/csfw/util/cstesting"
	"github.com/corestoreio/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var (
	_ json.Unmarshaler = (*money.Money)(nil)
	_ json.Marshaler   = (*money.Money)(nil)
	_ sql.Scanner      = (*money.Money)(nil)
	_ driver.Valuer    = (*money.Money)(nil)
)

func TestJSONMarshal(t *testing.T) {

	// @todo these tests will fail once i18n has been fully implemented. so fix this.
	var prefix = `"` + string([]byte("$")) + " "
	tests := []struct {
		prec      int
		haveI     int64
		haveEnc   money.Encoder
		haveValid bool
		want      string
		wantErr   error
	}{
		{100, 123456, money.JSONNumber, true, `1234.56`, nil},
		{1000, 123456, money.JSONNumber, true, `123.456`, nil},
		{10000, 123456, money.JSONNumber, true, `12.3456`, nil},
		{10, 123456, money.JSONNumber, true, `12345.6`, nil},
		{100, 123456, money.JSONNumber, false, `null`, nil},
		{0, 123456, money.JSONNumber, true, `123456`, nil},

		{100, 123456, money.JSONLocale, true, prefix + `1,234.56"`, nil},
		{1000, 123456, money.JSONLocale, true, prefix + `123.46"`, nil},
		{10000, 123456, money.JSONLocale, true, prefix + `12.35"`, nil},
		{10, 123456, money.JSONLocale, true, prefix + `12,345.60"`, nil},
		{100, 123456, money.JSONLocale, false, `null`, nil},
		{0, 123456, money.JSONLocale, true, prefix + `123,456.00"`, nil},

		{100, 123456, money.JSONExtended, true, `[1234.56, "$", "$ 1,234.56"]`, nil},
		{1000, 123456, money.JSONExtended, true, `[123.456, "$", "$ 123.46"]`, nil},
		{10000, 123456, money.JSONExtended, true, `[12.3456, "$", "$ 12.35"]`, nil},
		{10, 123456, money.JSONExtended, true, `[12345.6, "$", "$ 12,345.60"]`, nil},
		{100, 123456, money.JSONExtended, false, `null`, nil},
		{0, 123456, money.JSONExtended, true, `[123456, "$", "$ 123,456.00"]`, nil},
	}

	for _, test := range tests {
		c := money.New(
			money.WithPrecision(test.prec),
		).Set(test.haveI)
		c.Valid = test.haveValid
		c.FmtCur = testFmtCur
		c.FmtNum = testFmtNum
		c.Encoder = test.haveEnc

		have, err := c.MarshalJSON()
		if test.wantErr != nil {
			assert.Error(t, err, "%v", test)
			assert.Nil(t, have, "%v", test)
		} else {
			haveS := string(have)
			assert.NoError(t, err, "%v", test)
			if haveS != test.want {
				// assert.Equal... is not useful in this case
				t.Errorf("\nHave: %s\nWant: %s\n", haveS, test.want)
			}
		}
	}
}
func TestJSONUnMarshalSingle(t *testing.T) {

	tests := []struct {
		haveEnc  money.Encoder
		jsonData []byte
		want     string
		wantErr  bool
	}{
		{money.JSONNumber, []byte{0xf1, 0x32, 0xd8, 0x8a, 0x12, 0x8a, 0x74, 0x2a, 0x5, 0x5d, 0x18, 0x39, 0xf9, 0xd7, 0x99, 0x8b}, `NaN`, true},

		{money.JSONNumber, []byte(`1999.0000`), `1999.0000`, false},
		{money.JSONNumber, []byte(`-0.01`), `-0.0100`, false},
		{money.JSONNumber, []byte(`null`), `NaN`, false},
		{money.JSONNumber, []byte(`1234.56789`), `1234.5679`, false},
		{money.JSONNumber, []byte(`-1234.56789`), `-1234.5679`, false},
		{money.JSONNumber, []byte(`2999x.0156`), `2999.0156`, false},
		{money.JSONNumber, []byte(`""`), `NaN`, false},

		{money.JSONLocale, []byte(`$ 999.00 `), `999.0000`, false},
		{money.JSONLocale, []byte(`EUR 999.00`), `999.0000`, false},
		{money.JSONLocale, []byte(`EUR 99x9.0'0`), `999.0000`, false},
		{money.JSONLocale, []byte("EUR \x00 99x9.0'0"), `999.0000`, false},
		{money.JSONLocale, []byte(`2 345 678,45 €`), `2345678.4500`, false},
		{money.JSONLocale, []byte(`2 345 367,834456 €`), `2345367.8345`, false},
		{money.JSONLocale, []byte(`null`), `NaN`, false},
		{money.JSONLocale, []byte(`1.705,99 €`), `1705.9900`, false},
		{money.JSONLocale, []byte(`705,99 €`), `705.9900`, false},
		{money.JSONLocale, []byte(`$ 5,123,705.94`), `5123705.9400`, false},
		{money.JSONLocale, []byte(`$ -6,705.99`), `-6705.9900`, false},
		{money.JSONLocale, []byte(`$ 705.99`), `705.9900`, false},
		{money.JSONLocale, []byte(`$ 70789`), `70789.0000`, false},

		{money.JSONExtended, []byte(`[999.0000,"$","$ 999.00"]`), `999.0000`, false},
		{money.JSONExtended, []byte(`[999.0000,null,null]`), `999.0000`, false},
		{money.JSONExtended, []byte(`[1,999.00236,null,null]`), `1.0000`, false},
		{money.JSONExtended, []byte(`[1999.00236,null,null]`), `1999.0024`, false},
		{money.JSONExtended, []byte(`null`), `NaN`, false},
		{money.JSONExtended, []byte(`[null,"$",null]`), `NaN`, false},
		{money.JSONExtended, []byte(`[null,null,null]`), `NaN`, false},
		{money.JSONExtended, []byte(`[ ]`), `NaN`, true},
	}
	for i, test := range tests {
		var c money.Money
		err := c.UnmarshalJSON(test.jsonData)

		if test.wantErr {
			assert.Error(t, err, "Index %d", i)
			assert.True(t, errors.IsNotValid(err), "Index %d => %s", err)
		} else {
			var buf []byte
			assert.NoError(t, err, "Index %d", i)
			buf = c.FtoaAppend(buf)
			have := string(buf)
			if test.want != have {
				t.Errorf("\nHave: %s\n\nWant: %s\nIndex %d", have, test.want, i)
			}
		}
	}
}

func TestJSONUnMarshalSlice(t *testing.T) {

	tests := []struct {
		haveEnc  money.Encoder
		jsonData []byte
		want     string
		wantErr  error
	}{
		{
			money.JSONNumber, []byte(`[{"ValueID":1,"AttributeID":73,"StoreID":0,"EntityID":1,"Value":1999.0000},{"ValueID":2,"AttributeID":78,"StoreID":0,"EntityID":1,"Value":null},{"ValueID":3,"AttributeID":74,"StoreID":0,"EntityID":1,"Value":-0.01},{"ValueID":4,"AttributeID":77,"StoreID":0,"EntityID":1,"Value":null},{"ValueID":5,"AttributeID":73,"StoreID":1,"EntityID":1,"Value":705.9933},{"ValueID":6,"AttributeID":73,"StoreID":4,"EntityID":1,"Value":705.9933},{"ValueID":7,"AttributeID":73,"StoreID":2,"EntityID":1,"Value":705.9933},{"ValueID":8,"AttributeID":73,"StoreID":3,"EntityID":1,"Value":705.9933}]`),
			`1999.000; NaN; -0.010; NaN; 705.993; 705.993; 705.993; 705.993; `,
			nil,
		},
		{
			money.JSONNumber, []byte(`[
			{"ValueID":1,"AttributeID":73,"StoreID":0,"EntityID":1,"Value": "2999.0156"},
			{"ValueID":2,"AttributeID":78,"StoreID":0,"EntityID":1,"Value":null},
			{"ValueID":3,"AttributeID":74,"StoreID":0,"EntityID":1,"Value":0.01},
			{"ValueID":4,"AttributeID":77,"StoreID":0,"EntityID":1,"Value":null},
			{"ValueID":5,"AttributeID":73,"StoreID":1,"EntityID":1,"Value":7059933},
			{"ValueID":6,"AttributeID":73,"StoreID":4,"EntityID":1,"Value":705.9933},
			{"ValueID":7,"AttributeID":73,"StoreID":2,"EntityID":1,"Value":705.9933},
			{"ValueID":8,"AttributeID":73,"StoreID":3,"EntityID":1,"Value":705.9933}
			]`),
			`2999.016; NaN; 0.010; NaN; 7059933.000; 705.993; 705.993; 705.993; `,
			nil,
		},
		{
			money.JSONNumber, []byte(`[{"ValueID":1,"AttributeID":73,"StoreID":0,"EntityID":1,"Value": "2999x.0156"}]`),
			`2999.016; `,
			nil,
		},
		{
			money.JSONLocale, []byte(`[
			{"ValueID":1,"AttributeID":73,"StoreID":0,"EntityID":1,"Value":"$ 999.00 "},
			{"ValueID":2,"AttributeID":78,"StoreID":0,"EntityID":1,"Value":"2 345 678,45 €"},
			{"ValueID":3,"AttributeID":74,"StoreID":0,"EntityID":1,"Value":null},
			{"ValueID":5,"AttributeID":73,"StoreID":1,"EntityID":1,"Value":"1.705,99 €"},
			{"ValueID":6,"AttributeID":73,"StoreID":4,"EntityID":1,"Value":"$ 5,123,705.94"},
			{"ValueID":7,"AttributeID":73,"StoreID":2,"EntityID":1,"Value":"$ -6,705.99"},
			{"ValueID":8,"AttributeID":73,"StoreID":3,"EntityID":1,"Value":"$ 705.99"},
			{"ValueID":8,"AttributeID":73,"StoreID":3,"EntityID":1,"Value":"$ 70789"}
			]`),
			`999.000; 2345678.450; NaN; 1705.990; 5123705.940; -6705.990; 705.990; 70789.000; `,
			nil,
		},
		{
			money.JSONExtended, []byte(`[{"ValueID":1,"AttributeID":73,"StoreID":0,"EntityID":1,"Value":[999.0000,"$","$ 999.00"]},
					{"ValueID":2,"AttributeID":78,"StoreID":0,"EntityID":1,"Value":null},
					{"ValueID":3,"AttributeID":74,"StoreID":0,"EntityID":1,"Value":null},
					{"ValueID":4,"AttributeID":77,"StoreID":0,"EntityID":1,"Value":null},
					{"ValueID":5,"AttributeID":73,"StoreID":1,"EntityID":1,"Value":[705.9933,"$","$ 705.99"]},
					{"ValueID":6,"AttributeID":73,"StoreID":4,"EntityID":1,"Value":[705.9933,"$","$ 705.99"]},
					{"ValueID":7,"AttributeID":73,"StoreID":2,"EntityID":1,"Value":[705.9933,"$","$ 705.99"]},
					{"ValueID":8,"AttributeID":73,"StoreID":3,"EntityID":1,"Value":[705.9933,"$","$ 705.99"]}
					]`),
			`999.000; NaN; NaN; NaN; 705.993; 705.993; 705.993; 705.993; `,
			nil,
		},
	}

	for _, test := range tests {
		var peds TableProductEntityDecimalSlice
		if err := json.Unmarshal(test.jsonData, &peds); err != nil {
			t.Error(err)
		}

		var buf bytes.Buffer
		for _, ped := range peds.Data {
			ped.Value.FmtNum = testFmtNum

			_, err := ped.Value.NumberWriter(&buf)
			if test.wantErr != nil {
				assert.Error(t, err)
				assert.EqualError(t, err, test.wantErr.Error())
			} else {
				assert.NoError(t, err)
			}
			buf.WriteString("; ")
		}
		have := buf.String()
		if test.want != have {
			t.Errorf("\nHave: %s\n\nWant: %s\n", have, test.want)
		}
	}
}

var _ dbr.ArgumentsAppender = (*TableProductEntityDecimal)(nil)
var _ dbr.Scanner = (*TableProductEntityDecimalSlice)(nil)

type TableProductEntityDecimal struct {
	ValueID     int64       // value_id int(11) NOT NULL PRI  auto_increment
	AttributeID int64       // attribute_id smallint(5) unsigned NOT NULL MUL DEFAULT '0'
	StoreID     int64       // store_id smallint(5) unsigned NOT NULL MUL DEFAULT '0'
	EntityID    int64       // entity_id int(10) unsigned NOT NULL MUL DEFAULT '0'
	Value       money.Money // value decimal(12,4) NULL
}

func (ped TableProductEntityDecimal) AppendArguments(stmtType int, args dbr.Arguments, columns []string) (dbr.Arguments, error) {
	for _, c := range columns {
		switch c {
		case "value_id":
			args = append(args, dbr.Int64(ped.ValueID))
		case "attribute_id":
			args = append(args, dbr.Int64(ped.AttributeID))
		case "store_id":
			args = append(args, dbr.Int64(ped.StoreID))
		case "entity_id":
			args = append(args, dbr.Int64(ped.EntityID))
		case "value":
			args = append(args, dbr.Float64(ped.Value.Getf()))
		default:
			panic("other statement types than insert are not yet supported")
		}
	}
	return args, nil
}

type TableProductEntityDecimalSlice struct {
	Convert dbr.RowConvert
	Data    []*TableProductEntityDecimal
}

func (peds *TableProductEntityDecimalSlice) RowScan(r *sql.Rows) error {
	if err := peds.Convert.Scan(r); err != nil {
		return err
	}
	o := new(TableProductEntityDecimal)
	for i, col := range peds.Convert.Columns {
		b := peds.Convert.Index(i)
		var err error
		switch col {
		case "value_id":
			o.ValueID, err = b.Int64()
		case "attribute_id":
			o.AttributeID, err = b.Int64()
		case "store_id":
			o.StoreID, err = b.Int64()
		case "entity_id":
			o.EntityID, err = b.Int64()
		case "value":
			var f sql.NullFloat64
			f, err = b.NullFloat64()
			o.Value = money.New().Setf(f.Float64)
			o.Value.Valid = f.Valid
		}
		if err != nil {
			return errors.Wrapf(err, "[dbr] Failed to convert value at row % with column index %d", peds.Convert.Count, i)
		}
	}
	peds.Data = append(peds.Data, o)
	return nil
}

func (peds *TableProductEntityDecimalSlice) AppendArguments(stmtType int, args dbr.Arguments, columns []string) (dbr.Arguments, error) {
	for _, ped := range peds.Data {
		for _, c := range columns {
			switch c {
			case "value_id":
				args = append(args, dbr.Int64(ped.ValueID))
			case "attribute_id":
				args = append(args, dbr.Int64(ped.AttributeID))
			case "store_id":
				args = append(args, dbr.Int64(ped.StoreID))
			case "entity_id":
				args = append(args, dbr.Int64(ped.EntityID))
			case "value":
				args = append(args, dbr.Float64(ped.Value.Getf()))
			default:
				panic("other statement types than insert are not yet supported")
			}
		}
	}
	return args, nil
}

func TestLoadFromDb(t *testing.T) {
	t.Skip("Only for hacking")

	conn := cstesting.MustConnectDB(t)
	defer cstesting.Close(t, conn.DB)

	sel := conn.Select("*").From(`catalog_product_entity_decimal`).Limit(10)
	peds := new(TableProductEntityDecimalSlice)
	if rows, err := sel.Load(context.TODO(), peds); err != nil {
		t.Error(err)
	} else if rows == 0 {
		t.Error("0 rows loaded")
	}

	for _, ped := range peds.Data {
		fmt.Printf("%#v\n", ped)
	}
}

func TestSaveToDb(t *testing.T) {
	//t.Skip("Only for hacking")

	conn := cstesting.MustConnectDB(t)
	defer cstesting.Close(t, conn.DB)

	peds := &TableProductEntityDecimalSlice{
		Data: []*TableProductEntityDecimal{
			{AttributeID: 73, StoreID: 0, EntityID: 1, Value: money.New(money.WithPrecision(4)).Set(9990000)},
			{AttributeID: 78, StoreID: 0, EntityID: 1, Value: money.New(money.WithPrecision(4))}, // null values
			{AttributeID: 74, StoreID: 0, EntityID: 1, Value: money.New(money.WithPrecision(4))}, // null values
			{AttributeID: 77, StoreID: 0, EntityID: 1, Value: money.New(money.WithPrecision(4))}, // null values
			{AttributeID: 73, StoreID: 1, EntityID: 1, Value: money.New(money.WithPrecision(4)).Set(7059933)},
			{AttributeID: 73, StoreID: 4, EntityID: 1, Value: money.New(money.WithPrecision(4)).Set(7059933)},
			{AttributeID: 73, StoreID: 2, EntityID: 1, Value: money.New(money.WithPrecision(4)).Set(7059933)},
			{AttributeID: 73, StoreID: 3, EntityID: 1, Value: money.New(money.WithPrecision(4)).Set(7059933)},
		},
	}

	//tuple := &TableProductEntityDecimal{ValueID: 0, AttributeID: 73, StoreID: 3, EntityID: 231, Value: money.New(money.WithPrecision(4)).Set(7779933)}
	//tuple2 := &TableProductEntityDecimal{ValueID: 0, AttributeID: 74, StoreID: 2, EntityID: 231, Value: money.New(money.WithPrecision(4)).Set(8889933)}
	ib := conn.InsertInto("catalog_product_entity_decimal2").
		AddColumns("attribute_id", "store_id", "entity_id", "value").
		AddRecords(peds).Interpolate()

	res, err := ib.Exec(context.TODO())
	require.NoError(t, err)
	t.Log(res.LastInsertId())
	t.Log(res.RowsAffected())
	//t.Logf("1: %#v", tuple)
	//t.Logf("2: %#v", tuple2)
}

func TestValue(t *testing.T) {

	tuple := &TableProductEntityDecimal{ValueID: 0, AttributeID: 73, StoreID: 3, EntityID: 231, Value: money.New(money.WithPrecision(4)).Set(7779933)}
	tuple2 := &TableProductEntityDecimal{ValueID: 0, AttributeID: 74, StoreID: 2, EntityID: 231, Value: money.New(money.WithPrecision(4)).Set(8889933)}
	ib := dbr.NewInsert("catalog_product_entity_decimal")

	ib.AddColumns("attribute_id", "store_id", "entity_id", "value")
	ib.AddRecords(tuple, tuple2)

	fullSql, _, err := ib.Interpolate().ToSQL()
	assert.NoError(t, err)
	assert.Contains(t, fullSql, `(73,3,231,777.9933),(74,2,231,888.9933)`)
}

func TestScan(t *testing.T) {

	tests := []struct {
		src     interface{}
		want    string
		wantErr error
	}{
		{nil, `NaN`, nil},
		{[]byte{0x39, 0x39, 0x39, 0x2e, 0x30, 0x30, 0x30, 0x30}, `999.0000`, nil},
		{[]byte{0x37, 0x30, 0x35, 0x2e, 0x39, 0x39, 0x33, 0x33}, `705.9933`, nil},
		{[]byte{0x37, 0x30, 0x35, 0x2e, 0x39, 0x39, 0x33, 0x33}, `705.9933`, nil},
		{[]byte{0x37, 0x30, 0x35, 0x2e, 0x39, 0x39, 0x33, 0x33}, `705.9933`, nil},
		{[]byte{0x37, 0x30, 0x35, 0x2e, 0x39, 0x39, 0x33, 0x33}, `705.9933`, nil},
		{[]byte{0x37, 0x30, 0x35, 0x2e, 0x19, 0x39, 0x33, 0x13}, `0.0000`, errors.New("strconv.ParseFloat: parsing \"705.\\x1993\\x13\": invalid syntax")},
		{[]byte{0x37, 0x33}, `73.0000`, nil},
		{[]byte{0x37, 0x38}, `78.0000`, nil},
		{[]byte{0x37, 0x34}, `74.0000`, nil},
		{[]byte{0x37, 0x37}, `77.0000`, nil},
		{[]byte{0xa7, 0x3e}, `0.0000`, errors.New("strconv.ParseFloat: parsing \"\\xa7>\": invalid syntax")},
		{int(33), `0.0000`, errors.New("Unsupported Type int for value '!'. Supported: []byte")},
	}

	var buf bytes.Buffer
	for i, test := range tests {
		var c money.Money
		err := c.Scan(test.src)
		c.FmtCur = testFmtCur
		c.FmtNum = testFmtNum

		if test.wantErr != nil {
			assert.Error(t, err, "%v", test, "Index %d", i)
			assert.EqualError(t, err, test.wantErr.Error(), "Index %d", i)
		} else {
			assert.NoError(t, err, "%v", test, "Index %d", i)
			assert.EqualValues(t, test.want, string(c.Ftoa()), "Index %d", i)

			if _, err := c.NumberWriter(&buf); err != nil {
				t.Error(err)
			}
			if _, err := buf.WriteString("; "); err != nil {
				t.Error(err)
			}
		}
	}

	want := `NaN; 999.000; 705.993; 705.993; 705.993; 705.993; 73.000; 78.000; 74.000; 77.000; `
	have := buf.String()
	if want != have {
		t.Errorf("\nHave: %s\n\nWant: %s\n", have, want)
	}
}

func TestJSONEncode(t *testing.T) {

	var peds = TableProductEntityDecimalSlice{
		Data: []*TableProductEntityDecimal{
			{ValueID: 1, AttributeID: 73, StoreID: 0, EntityID: 1, Value: money.New(money.WithPrecision(4)).Set(9990000)},
			{ValueID: 2, AttributeID: 78, StoreID: 0, EntityID: 1, Value: money.New(money.WithPrecision(4))}, // null values
			{ValueID: 3, AttributeID: 74, StoreID: 0, EntityID: 1, Value: money.New(money.WithPrecision(4))}, // null values
			{ValueID: 4, AttributeID: 77, StoreID: 0, EntityID: 1, Value: money.New(money.WithPrecision(4))}, // null values
			{ValueID: 5, AttributeID: 73, StoreID: 1, EntityID: 1, Value: money.New(money.WithPrecision(4)).Set(7059933)},
			{ValueID: 6, AttributeID: 73, StoreID: 4, EntityID: 1, Value: money.New(money.WithPrecision(4)).Set(7059933)},
			{ValueID: 7, AttributeID: 73, StoreID: 2, EntityID: 1, Value: money.New(money.WithPrecision(4)).Set(7059933)},
			{ValueID: 8, AttributeID: 73, StoreID: 3, EntityID: 1, Value: money.New(money.WithPrecision(4)).Set(7059933)},
		},
	}

	jb, err := json.Marshal(peds)
	assert.NoError(t, err)
	have := string(jb)
	want := `[{"ValueID":1,"AttributeID":73,"StoreID":0,"EntityID":1,"Value":"Cur:$ Sign:1 I:999 Prec:4 Frac:0"},{"ValueID":2,"AttributeID":78,"StoreID":0,"EntityID":1,"Value":null},{"ValueID":3,"AttributeID":74,"StoreID":0,"EntityID":1,"Value":null},{"ValueID":4,"AttributeID":77,"StoreID":0,"EntityID":1,"Value":null},{"ValueID":5,"AttributeID":73,"StoreID":1,"EntityID":1,"Value":"Cur:$ Sign:1 I:705 Prec:4 Frac:9933"},{"ValueID":6,"AttributeID":73,"StoreID":4,"EntityID":1,"Value":"Cur:$ Sign:1 I:705 Prec:4 Frac:9933"},{"ValueID":7,"AttributeID":73,"StoreID":2,"EntityID":1,"Value":"Cur:$ Sign:1 I:705 Prec:4 Frac:9933"},{"ValueID":8,"AttributeID":73,"StoreID":3,"EntityID":1,"Value":"Cur:$ Sign:1 I:705 Prec:4 Frac:9933"}]`
	if have != want {
		t.Errorf("\nHave: %s\n\nWant: %s\n", have, want)
	}
}
