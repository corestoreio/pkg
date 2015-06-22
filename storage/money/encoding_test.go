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
	"strconv"
	"testing"

	"github.com/corestoreio/csfw/i18n"
	"github.com/corestoreio/csfw/storage/money"
	"github.com/stretchr/testify/assert"
)

func TestJSONMarshal(t *testing.T) {

	// @todo these tests will fail once i18n has been fully implemented. so fix this.
	var prefix = `"` + string(i18n.DefaultCurrencySign) + "\\u00A0" // because JSEscape
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

		{100, 123456, money.JSONExtended, true, `[1234.56, "$", "$\u00A01,234.56"]`, nil},
		{1000, 123456, money.JSONExtended, true, `[123.456, "$", "$\u00A0123.46"]`, nil},
		{10000, 123456, money.JSONExtended, true, `[12.3456, "$", "$\u00A012.35"]`, nil},
		{10, 123456, money.JSONExtended, true, `[12345.6, "$", "$\u00A012,345.60"]`, nil},
		{100, 123456, money.JSONExtended, false, `null`, nil},
		{0, 123456, money.JSONExtended, true, `[123456, "$", "$\u00A0123,456.00"]`, nil},
	}

	for _, test := range tests {
		c := money.New(
			money.Precision(test.prec),
			money.JSONMarshal(test.haveEnc),
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
	}

	for _, test := range tests {
		var c money.Currency
		err := c.Scan(test.src)

		if test.wantErr != nil {
			assert.Error(t, err, "%v", test)
			assert.Contains(t, err.Error(), test.wantErr.Error())
		} else {
			assert.NoError(t, err, "%v", test)
			assert.EqualValues(t, test.want, string(c.Ftoa()), "%v", test)
		}
	}

	// for hacking testing added :-)
	//	type TableProductEntityDecimal struct {
	//		ValueID     int64          `db:"value_id"`     // value_id int(11) NOT NULL PRI  auto_increment
	//		AttributeID int64          `db:"attribute_id"` // attribute_id smallint(5) unsigned NOT NULL MUL DEFAULT '0'
	//		StoreID     int64          `db:"store_id"`     // store_id smallint(5) unsigned NOT NULL MUL DEFAULT '0'
	//		EntityID    int64          `db:"entity_id"`    // entity_id int(10) unsigned NOT NULL MUL DEFAULT '0'
	//		Value       money.Currency `db:"value"`        // value decimal(12,4) NULL
	//	}
	//
	//	type TableProductEntityDecimalSlice []*TableProductEntityDecimal
	//
	//	db := csdb.MustConnectTest()
	//	defer db.Close()
	//	dbrSess := dbr.NewConnection(db, nil).NewSession(nil)
	//
	//	sel := dbrSess.SelectBySql("SELECT * FROM `catalog_product_entity_decimal`")
	//	var peds TableProductEntityDecimalSlice
	//
	//	if rows, err := sel.LoadStructs(&peds); err != nil {
	//		t.Error(err)
	//	} else if rows == 0 {
	//		t.Error("0 rows loaded")
	//	}
	//
	//	for _, ped := range peds {
	//		t.Logf("%s\n", ped.Value.Ftoa())
	//	}

}
