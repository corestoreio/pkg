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

package byteconv

import (
	"database/sql"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseNullBoolSQL_ParseBoolSQL(t *testing.T) {

	runner := func(have string, want sql.NullBool) func(*testing.T) {
		return func(t *testing.T) {
			b := sql.RawBytes(have)
			if have == "NULL" {
				b = nil
			}
			assert.Exactly(t, want, ParseNullBoolSQL(&b), t.Name())
			assert.Exactly(t, want.Bool, ParseBoolSQL(&b), t.Name())
		}
	}
	t.Run("NULL is false and invalid", runner("NULL", sql.NullBool{}))
	t.Run("empty is false and invalid", runner("", sql.NullBool{}))
	t.Run(" is false and invalid", runner("", sql.NullBool{}))
	t.Run("£ is false and invalid", runner("£", sql.NullBool{}))
	t.Run("0 is false and valid", runner("0", sql.NullBool{Valid: true}))
	t.Run("1 is true and valid", runner("1", sql.NullBool{Valid: true, Bool: true}))
	t.Run("10 is false and invalid", runner("10", sql.NullBool{}))
	t.Run("01 is false and invalid", runner("01", sql.NullBool{}))
}
