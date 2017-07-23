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
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestColumn(t *testing.T) {
	t.Run("invalid column name", func(t *testing.T) {
		s := NewSelect("a", "b").From("c").Where(
			Column("a").Int(111),
			Expression("b=c"),
		)
		sql, args, err := s.ToSQL()
		require.NoError(t, err)
		assert.Equal(t, "SELECT `a`, `b` FROM `c` WHERE (`a` = ?) AND (b=c)", sql)
		assert.Equal(t, []interface{}{int64(111)}, args)
	})

	t.Run("valid column name", func(t *testing.T) {
		s := NewSelect("a", "b").From("c").Where(
			Column("a").Ints(111, 222), // omitted In(). on purpose because default operator is IN for slices
			Column("b").Null(),
			Column("d").Between().Float64s(2.5, 2.7),
		).Interpolate()
		sql, args, err := s.ToSQL()
		require.NoError(t, err)
		assert.Nil(t, args)
		assert.Equal(t, "SELECT `a`, `b` FROM `c` WHERE (`a` IN (111,222)) AND (`b` IS NULL) AND (`d` BETWEEN 2.5 AND 2.7)", sql)
	})
}

func TestConditions_writeOnDuplicateKey(t *testing.T) {

	runner := func(cnds Conditions, wantSQL string, wantArgs ...interface{}) func(*testing.T) {
		return func(t *testing.T) {
			buf := new(bytes.Buffer)
			args := make(Arguments, 0, 2)
			err := cnds.writeOnDuplicateKey(buf)
			require.NoError(t, err)
			args, _, err = cnds.appendArgs(args, appendArgsDUPKEY)
			require.NoError(t, err)
			assert.Exactly(t, wantSQL, buf.String(), t.Name())
			assert.Exactly(t, wantArgs, args.Interfaces(), t.Name())
		}
	}
	t.Run("empty columns does nothing", runner(Conditions{}, ""))

	t.Run("col=VALUES(col) and no arguments", runner(Conditions{
		Columns("sku", "name", "stock"),
	}, " ON DUPLICATE KEY UPDATE `sku`=VALUES(`sku`), `name`=VALUES(`name`), `stock`=VALUES(`stock`)"))

	t.Run("col=? and with arguments", runner(Conditions{
		Column("name").String("E0S 5D Mark II"),
		Column("stock").Int64(12),
	}, " ON DUPLICATE KEY UPDATE `name`=?, `stock`=?",
		"E0S 5D Mark II", int64(12)))

	t.Run("col=VALUES(val)+? and with arguments", runner(Conditions{
		Column("name").String("E0S 5D Mark II"),
		Column("stock").Expression("VALUES(`stock`)+?").Int64(13), // TODO add more args
	}, " ON DUPLICATE KEY UPDATE `name`=?, `stock`=VALUES(`stock`)+?",
		"E0S 5D Mark II", int64(13)))

	t.Run("col=VALUES(val) and with arguments and nil", runner(Conditions{
		Column("name").String("E0S 5D Mark III"),
		Column("sku").Values(),
		Column("stock").Int64(14),
	}, " ON DUPLICATE KEY UPDATE `name`=?, `sku`=VALUES(`sku`), `stock`=?",
		"E0S 5D Mark III", int64(14)))

	t.Run("col=expression without arguments", runner(Conditions{
		Column("name").Expression("CONCAT('Canon','E0S 5D Mark III')"),
	}, " ON DUPLICATE KEY UPDATE `name`=CONCAT('Canon','E0S 5D Mark III')",
	))
}
