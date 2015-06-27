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

package money_test

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"testing"

	"github.com/corestoreio/csfw/i18n"
	"github.com/corestoreio/csfw/storage/csdb"
	"github.com/corestoreio/csfw/storage/dbr"
	"github.com/corestoreio/csfw/storage/money"
	"github.com/stretchr/testify/assert"
)

func TestJSONMarshal(t *testing.T) {

	// @todo these tests will fail once i18n has been fully implemented. so fix this.
	var prefix = `"` + string(i18n.DefaultCurrencySign) + " "
	tests := []struct {
		prec      int
		haveI     int64
		haveEnc   money.JSONMarshaller
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
			money.Precision(test.prec),
			money.JSONMarshal(test.haveEnc),
			money.FormatCurrency(testFmtCur),
			money.FormatNumber(testFmtNum),
		).Set(test.haveI)
		c.Valid = test.haveValid

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
		haveEnc  money.JSONMarshaller
		jsonData []byte
		want     string
		wantErr  error
	}{
		{money.JSONNumber, []byte{0xf1, 0x32, 0xd8, 0x8a, 0x12, 0x8a, 0x74, 0x2a, 0x5, 0x5d, 0x18, 0x39, 0xf9, 0xd7, 0x99, 0x8b}, `NaN`, nil},

		{money.JSONNumber, []byte(`1999.0000`), `1999.0000`, nil},
		{money.JSONNumber, []byte(`-0.01`), `-0.0100`, nil},
		{money.JSONNumber, []byte(`null`), `NaN`, nil},
		{money.JSONNumber, []byte(`1234.56789`), `1234.5679`, nil},
		{money.JSONNumber, []byte(`-1234.56789`), `-1234.5679`, nil},
		{money.JSONNumber, []byte(`2999x.0156`), `2999.0156`, nil},
		{money.JSONNumber, []byte(`""`), `NaN`, nil},

		{money.JSONLocale, []byte(`$ 999.00 `), `999.0000`, nil},
		{money.JSONLocale, []byte(`EUR 999.00`), `999.0000`, nil},
		{money.JSONLocale, []byte(`EUR 99x9.0'0`), `999.0000`, nil},
		{money.JSONLocale, []byte("EUR \x00 99x9.0'0"), `999.0000`, nil},
		{money.JSONLocale, []byte(`2 345 678,45 €`), `2345678.4500`, nil},
		{money.JSONLocale, []byte(`2 345 367,834456 €`), `2345367.8345`, nil},
		{money.JSONLocale, []byte(`null`), `NaN`, nil},
		{money.JSONLocale, []byte(`1.705,99 €`), `1705.9900`, nil},
		{money.JSONLocale, []byte(`705,99 €`), `705.9900`, nil},
		{money.JSONLocale, []byte(`$ 5,123,705.94`), `5123705.9400`, nil},
		{money.JSONLocale, []byte(`$ -6,705.99`), `-6705.9900`, nil},
		{money.JSONLocale, []byte(`$ 705.99`), `705.9900`, nil},
		{money.JSONLocale, []byte(`$ 70789`), `70789.0000`, nil},

		{money.JSONExtended, []byte(`[999.0000,"$","$ 999.00"]`), `999.0000`, nil},
		{money.JSONExtended, []byte(`[999.0000,null,null]`), `999.0000`, nil},
		{money.JSONExtended, []byte(`[1,999.00236,null,null]`), `1.0000`, nil},
		{money.JSONExtended, []byte(`[1999.00236,null,null]`), `1999.0024`, nil},
		{money.JSONExtended, []byte(`null`), `NaN`, nil},
		{money.JSONExtended, []byte(`[null,"$",null]`), `NaN`, nil},
		{money.JSONExtended, []byte(`[null,null,null]`), `NaN`, nil},
		{money.JSONExtended, []byte(`[ ]`), `NaN`, money.ErrDecodeMissingColon},
	}
	for _, test := range tests {
		var c money.Currency
		err := c.UnmarshalJSON(test.jsonData)

		if test.wantErr != nil {
			assert.Error(t, err)
			assert.EqualError(t, err, test.wantErr.Error())
		} else {
			var buf []byte
			assert.NoError(t, err)
			buf = c.FtoaAppend(buf)
			have := string(buf)
			if test.want != have {
				t.Errorf("\nHave: %s\n\nWant: %s\n", have, test.want)
			}

		}
	}
}

func TestJSONUnMarshalSlice(t *testing.T) {

	tests := []struct {
		haveEnc  money.JSONMarshaller
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
			`999.000; 2345678.450; NaN; 170599.000; 5123705.940; -6705.990; 705.990; 70789.000; `,
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
		for _, ped := range peds {
			ped.Value.Option(
				money.FormatNumber(testFmtNum),
			)
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

type TableProductEntityDecimal struct {
	ValueID     int64          `db:"value_id"`     // value_id int(11) NOT NULL PRI  auto_increment
	AttributeID int64          `db:"attribute_id"` // attribute_id smallint(5) unsigned NOT NULL MUL DEFAULT '0'
	StoreID     int64          `db:"store_id"`     // store_id smallint(5) unsigned NOT NULL MUL DEFAULT '0'
	EntityID    int64          `db:"entity_id"`    // entity_id int(10) unsigned NOT NULL MUL DEFAULT '0'
	Value       money.Currency `db:"value"`        // value decimal(12,4) NULL
}

type TableProductEntityDecimalSlice []*TableProductEntityDecimal

func off_TestLoadFromDb(t *testing.T) {
	//for hacking testing added :-)
	db := csdb.MustConnectTest()
	defer db.Close()
	dbrSess := dbr.NewConnection(db, nil).NewSession(nil)

	sel := dbrSess.SelectBySql("SELECT * FROM `catalog_product_entity_decimal`")
	var peds TableProductEntityDecimalSlice

	if rows, err := sel.LoadStructs(&peds); err != nil {
		t.Error(err)
	} else if rows == 0 {
		t.Error("0 rows loaded")
	}

	for _, ped := range peds {
		fmt.Printf("%#v\n", ped)
	}
}

func TestSaveToDb(t *testing.T) {
	//for hacking testing added :-)
	db := csdb.MustConnectTest()
	defer db.Close()
	dbrSess := dbr.NewConnection(db, nil).NewSession(nil)

	//	var peds = TableProductEntityDecimalSlice{
	//		&TableProductEntityDecimal{ValueID: 1, AttributeID: 73, StoreID: 0, EntityID: 1, Value: money.New(money.Precision(4)).Set(9990000)},
	//		&TableProductEntityDecimal{ValueID: 2, AttributeID: 78, StoreID: 0, EntityID: 1, Value: money.New(money.Precision(4))}, // null values
	//		&TableProductEntityDecimal{ValueID: 3, AttributeID: 74, StoreID: 0, EntityID: 1, Value: money.New(money.Precision(4))}, // null values
	//		&TableProductEntityDecimal{ValueID: 4, AttributeID: 77, StoreID: 0, EntityID: 1, Value: money.New(money.Precision(4))}, // null values
	//		&TableProductEntityDecimal{ValueID: 5, AttributeID: 73, StoreID: 1, EntityID: 1, Value: money.New(money.Precision(4)).Set(7059933)},
	//		&TableProductEntityDecimal{ValueID: 6, AttributeID: 73, StoreID: 4, EntityID: 1, Value: money.New(money.Precision(4)).Set(7059933)},
	//		&TableProductEntityDecimal{ValueID: 7, AttributeID: 73, StoreID: 2, EntityID: 1, Value: money.New(money.Precision(4)).Set(7059933)},
	//		&TableProductEntityDecimal{ValueID: 8, AttributeID: 73, StoreID: 3, EntityID: 1, Value: money.New(money.Precision(4)).Set(7059933)},
	//	}

	tuple := &TableProductEntityDecimal{ValueID: 0, AttributeID: 73, StoreID: 3, EntityID: 1, Value: money.New(money.Precision(4)).Set(7779933)}
	ib := dbrSess.InsertInto("catalog_product_entity_decimal")
	ib.Columns("attribute_id", "store_id", "entity_id", "value")

	ib.Record(tuple)
	// this is a bug in the dbr package ToSql because it ignores the Value() driver.Value interface
	t.Error(ib.ToSql())

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
		{[]byte{0x37, 0x30, 0x35, 0x2e, 0x19, 0x39, 0x33, 0x13}, `0.0000`, strconv.ErrSyntax},
		{[]byte{0x37, 0x33}, `73.0000`, nil},
		{[]byte{0x37, 0x38}, `78.0000`, nil},
		{[]byte{0x37, 0x34}, `74.0000`, nil},
		{[]byte{0x37, 0x37}, `77.0000`, nil},
		{[]byte{0xa7, 0x3e}, `0.0000`, strconv.ErrSyntax},
		{int(33), `0.0000`, errors.New("Unsupported Type for value")},
	}

	var buf bytes.Buffer
	for _, test := range tests {
		var c money.Currency
		err := c.Scan(test.src)
		c.Option(
			money.FormatCurrency(testFmtCur),
			money.FormatNumber(testFmtNum),
		)

		if test.wantErr != nil {
			assert.Error(t, err, "%v", test)
			assert.Contains(t, err.Error(), test.wantErr.Error())
		} else {
			assert.NoError(t, err, "%v", test)
			assert.EqualValues(t, test.want, string(c.Ftoa()), "%v", test)

			if _, err := c.NumberWriter(&buf); err != nil {
				t.Error(err)
			}
			buf.WriteString("; ")
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
		&TableProductEntityDecimal{ValueID: 1, AttributeID: 73, StoreID: 0, EntityID: 1, Value: money.New(money.Precision(4)).Set(9990000)},
		&TableProductEntityDecimal{ValueID: 2, AttributeID: 78, StoreID: 0, EntityID: 1, Value: money.New(money.Precision(4))}, // null values
		&TableProductEntityDecimal{ValueID: 3, AttributeID: 74, StoreID: 0, EntityID: 1, Value: money.New(money.Precision(4))}, // null values
		&TableProductEntityDecimal{ValueID: 4, AttributeID: 77, StoreID: 0, EntityID: 1, Value: money.New(money.Precision(4))}, // null values
		&TableProductEntityDecimal{ValueID: 5, AttributeID: 73, StoreID: 1, EntityID: 1, Value: money.New(money.Precision(4)).Set(7059933)},
		&TableProductEntityDecimal{ValueID: 6, AttributeID: 73, StoreID: 4, EntityID: 1, Value: money.New(money.Precision(4)).Set(7059933)},
		&TableProductEntityDecimal{ValueID: 7, AttributeID: 73, StoreID: 2, EntityID: 1, Value: money.New(money.Precision(4)).Set(7059933)},
		&TableProductEntityDecimal{ValueID: 8, AttributeID: 73, StoreID: 3, EntityID: 1, Value: money.New(money.Precision(4)).Set(7059933)},
	}

	jb, err := json.Marshal(peds)
	assert.NoError(t, err)
	have := string(jb)
	want := `[{"ValueID":1,"AttributeID":73,"StoreID":0,"EntityID":1,"Value":"$\u00A0999.00"},{"ValueID":2,"AttributeID":78,"StoreID":0,"EntityID":1,"Value":null},{"ValueID":3,"AttributeID":74,"StoreID":0,"EntityID":1,"Value":null},{"ValueID":4,"AttributeID":77,"StoreID":0,"EntityID":1,"Value":null},{"ValueID":5,"AttributeID":73,"StoreID":1,"EntityID":1,"Value":"$\u00A0705.99"},{"ValueID":6,"AttributeID":73,"StoreID":4,"EntityID":1,"Value":"$\u00A0705.99"},{"ValueID":7,"AttributeID":73,"StoreID":2,"EntityID":1,"Value":"$\u00A0705.99"},{"ValueID":8,"AttributeID":73,"StoreID":3,"EntityID":1,"Value":"$\u00A0705.99"}]`
	if have != want {
		t.Errorf("\nHave: %s\n\nWant: %s\n", have, want)
	}
}
