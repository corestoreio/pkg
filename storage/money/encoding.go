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

package money

import (
	"database/sql"
	"database/sql/driver"
	"encoding/json"

	"github.com/corestoreio/csfw/utils/log"
	"github.com/juju/errgo"
)

var (
	_          json.Unmarshaler = (*Money)(nil)
	_          json.Marshaler   = (*Money)(nil)
	_          sql.Scanner      = (*Money)(nil)
	_          driver.Valuer    = (*Money)(nil)
	nullString                  = []byte(`null`)
)

type (
	// Encoder interface to encode money into bytes
	Encoder interface {
		// Encode encodes the currency into bytes
		Encode(*Money) ([]byte, error)
	}
	// Decoder interface to decode money type
	Decoder interface {
		// Decode reads the bytes and decodes them into the currency
		Decode(*Money, []byte) error
	}
)

// MarshalJSON generates JSON output depending on the Encoder.
func (m Money) MarshalJSON() ([]byte, error) {
	return m.Encode(&m)
}

// UnmarshalJSON reads JSON and fills the money struct depending on the Decoder.
func (m *Money) UnmarshalJSON(src []byte) error {
	m.applyDefaults()
	if src == nil {
		m.m, m.Valid = 0, false
		return nil
	}
	return m.Decode(m, src)
}

// Value implements the SQL driver Valuer interface.
func (m Money) Value() (driver.Value, error) {
	if !m.Valid {
		return nil, nil
	}
	return m.Getf(), nil
}

// Scan scans a value into the Money struct. Returns an error on data loss.
// Errors will be logged. Initial default settings are the guard and precision value.
func (m *Money) Scan(src interface{}) error {
	m.applyDefaults()

	if src == nil {
		m.m, m.Valid = 0, false
		if log.IsDebug() {
			log.Debug("money.Currency.Scan", "case", 89, "c", m, "src", src)
		}
		return nil
	}

	if b, ok := src.([]byte); ok {
		return m.ParseFloat(string(b))
	}
	return log.Error("money.Currency.Scan.Assertion", "err", errgo.Newf("Unsupported Type %T for value. Supported: []byte", src), "src", src)
}
