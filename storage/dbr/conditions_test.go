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
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestColumn(t *testing.T) {
	t.Run("invalid column name", func(t *testing.T) {
		s := NewSelect("a", "b").From("c").Where(
			Column("a", ArgInt(111)),
			Column("b=c"),
		)
		sql, args, err := s.ToSQL()
		require.NoError(t, err)
		assert.Equal(t, "SELECT `a`, `b` FROM `c` WHERE (`a` = ?) AND (`b=c`)", sql)
		assert.Equal(t, []interface{}{int64(111)}, args)
	})

	t.Run("valid column name", func(t *testing.T) {
		s := NewSelect("a", "b").From("c").Where(
			Column("a", In.Int64(111, 222)),
			Column("b"),
			Column("d", Between.Float64(2.5, 2.7)),
		).Interpolate()
		sql, args, err := s.ToSQL()
		require.NoError(t, err)
		assert.Nil(t, args)
		assert.Equal(t, "SELECT `a`, `b` FROM `c` WHERE (`a` IN (111,222)) AND (`b`) AND (`d` BETWEEN 2.5 AND 2.7)", sql)
	})
}
