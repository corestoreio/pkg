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
	_ sql.Scanner              = (*CSN[int])(nil)
	_ driver.Valuer            = (*CSN[int])(nil)
	_ encoding.TextMarshaler   = (*CSN[int])(nil)
	_ encoding.TextUnmarshaler = (*CSN[int])(nil)
)

func numbers_runner[K numbers](data any, wantErr error, wantCSV []K, wantStringified string) func(*testing.T) {
	return func(t *testing.T) {
		var csv CSN[K]
		haveErr := csv.Scan(data)
		if wantErr != nil {
			require.EqualError(t, haveErr, wantErr.Error())
			return
		}
		require.NoError(t, haveErr)
		assert.Exactly(t, wantCSV, []K(csv))

		haveBytes, err := csv.Bytes()
		require.NoError(t, err)
		assert.Exactly(t, wantStringified, string(haveBytes))
	}
}

func TestCSN_Scan_Bytes_Int(t *testing.T) {
	t.Run("simple int", numbers_runner[int](`1,2`, nil, []int{1, 2}, "1,2"))
	t.Run("simple float", numbers_runner[float32](`1.2,-2.3`, nil, []float32{1.2, -2.3}, "1.2,-2.3"))
	t.Run("minus", numbers_runner[int]([]byte(`-1,-3`), nil, []int{-1, -3}, "-1,-3"))
	t.Run("ws quoted", numbers_runner(`;44;    -1;0`, nil, []int{44, -1, 0}, "44,-1,0"))
	// t.Run("weird char as first rune", numbers_runner(`-3,1`, nil, []int{-3, 1}, "-3,1"))
}
