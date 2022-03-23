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

package dmltype

import (
	"database/sql"
	"database/sql/driver"
	"encoding"
	"testing"

	"github.com/corestoreio/pkg/util/assert"
	"github.com/stretchr/testify/require"
)

var (
	_ sql.Scanner              = (*CSV)(nil)
	_ driver.Valuer            = (*CSV)(nil)
	_ encoding.TextMarshaler   = (*CSV)(nil)
	_ encoding.TextUnmarshaler = (*CSV)(nil)
)

func TestStrings_Scan_Bytes(t *testing.T) {
	runner := func(data any, wantErr error, wantCSV []string, wantStringified string) func(*testing.T) {
		return func(t *testing.T) {
			var csv CSV
			haveErr := csv.Scan(data)
			if wantErr != nil {
				require.EqualError(t, haveErr, wantErr.Error())
				return
			}
			require.NoError(t, haveErr)
			assert.Exactly(t, wantCSV, []string(csv))

			haveBytes, err := csv.Bytes()
			require.NoError(t, err)
			assert.Exactly(t, wantStringified, string(haveBytes))
		}
	}
	t.Run("simple", runner(`a,b`, nil, []string{"a", "b"}, "a,b"))
	t.Run("quoted", runner(`a"a,b'b`, nil, []string{"a\"a", "b'b"}, "a\"a,b'b"))
	t.Run("ws quoted", runner(`a"a,    b'b`, nil, []string{"a\"a", "b'b"}, "a\"a,b'b"))
	t.Run("comma quoted1", runner(`a,"a,,,    b'b`, nil, []string{"a", "\"a", "b'b"}, "a,\"a,b'b"))
	t.Run("comma quoted2", runner(`a,"a,,,b,c`, nil, []string{"a", "\"a", "b", "c"}, "a,\"a,b,c"))
	t.Run("semi colon as first rune", runner(`;a;a;;b;c`, nil, []string{"a", "a", "b", "c"}, "a,a,b,c"))
	t.Run("weird char as first rune", runner(`a,c`, nil, []string{"a", "c"}, "a,c"))
}
