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

package money

import (
	"database/sql/driver"

	"github.com/corestoreio/errors"
)

var (
	nullString = []byte(`null`)
)

// MarshalJSON generates JSON output depending on the Encoder.
func (m Money) MarshalJSON() ([]byte, error) {
	return nil, nil
}

// UnmarshalJSON reads JSON and fills the money struct depending on the Decoder.
func (m *Money) UnmarshalJSON(src []byte) error {

	if src == nil {
		m.m, m.Valid = 0, false
		return nil
	}
	return nil
}

// Value implements the SQL driver Valuer interface.
func (m Money) Value() (driver.Value, error) {
	if !m.Valid {
		return nil, nil
	}
	return m.Getf(), nil
}

// Scan scans a value into the Money struct. Returns an error on data loss.
// Errors will be logged. Initial default settings are the guard and precision
// value.
func (m *Money) Scan(src interface{}) error {
	// TODO type switch

	if src == nil {
		m.m, m.Valid = 0, false
		return nil
	}

	if b, ok := src.([]byte); ok {
		return m.ParseFloat(string(b))
	}
	return errors.Errorf("Unsupported Type %T for value %q. Supported: []byte", src, src)
}
