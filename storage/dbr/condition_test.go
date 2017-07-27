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
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var _ fmt.Stringer = (*expressions)(nil)

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

	t.Run("col1=VALUES(val)+? and with arguments", runner(Conditions{
		Column("name").String("E0S 5D Mark II"),
		Column("stock").Expression("VALUES(`stock`)+?-?").Int64(13).Int(4),
	}, " ON DUPLICATE KEY UPDATE `name`=?, `stock`=VALUES(`stock`)+?-?",
		"E0S 5D Mark II", int64(13), int64(4)))

	t.Run("col2=VALUES(val) and with arguments and nil", runner(Conditions{
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

func TestExpression(t *testing.T) {

	t.Run("ints", func(t *testing.T) {
		sel := NewSelect("a").From("c").
			Where(
				Column("g").Int(3),
				Expression("i1 = ? AND i2 IN (?) AND i64_1 = ? AND i64_2 IN (?) AND ui64 > ? AND f64_1 = ? AND f64_2 IN (?)").
					Int(1).Ints(2, 3).
					Int64(4).Int64s(5, 6).
					Uint64(7).
					Float64(4.51).Float64s(5.41, 6.66666),
			)

		compareToSQL(t, sel, nil,
			"SELECT `a` FROM `c` WHERE (`g` = ?) AND (i1 = ? AND i2 IN (?,?) AND i64_1 = ? AND i64_2 IN (?,?) AND ui64 > ? AND f64_1 = ? AND f64_2 IN (?,?))",
			"SELECT `a` FROM `c` WHERE (`g` = 3) AND (i1 = 1 AND i2 IN (2,3) AND i64_1 = 4 AND i64_2 IN (5,6) AND ui64 > 7 AND f64_1 = 4.51 AND f64_2 IN (5.41,6.66666))",
			int64(3), int64(1), int64(2), int64(3), int64(4), int64(5), int64(6), []byte(`7`), 4.51, 5.41, 6.66666,
		)
	})

	t.Run("slice expression", func(t *testing.T) {
		sel := NewSelect("a").From("c").
			Where(
				Column("h").In().Int64s(1, 2, 3),
				Expression("l NOT IN (?)").Strings("xx", "yy"),
			)
		compareToSQL(t, sel, nil,
			"SELECT `a` FROM `c` WHERE (`h` IN (?,?,?)) AND (l NOT IN (?,?))",
			"SELECT `a` FROM `c` WHERE (`h` IN (1,2,3)) AND (l NOT IN ('xx','yy'))",
			int64(1), int64(2), int64(3), "xx", "yy",
		)
	})

	t.Run("string bools", func(t *testing.T) {
		sel := NewSelect("a").From("c").
			Where(
				Column("h").In().Int64s(1, 2, 3),
				Expression("l = ? AND m IN (?) AND n = ? AND o IN (?) AND p = ? AND q IN (?)").
					String("xx").Strings("aa", "bb", "cc").
					Bool(true).Bools(true, false, true).
					Bytes([]byte(`Gopher`)).BytesSlice([]byte(`Go1`), []byte(`Go2`)),
			)

		compareToSQL(t, sel, nil,
			"SELECT `a` FROM `c` WHERE (`h` IN (?,?,?)) AND (l = ? AND m IN (?,?,?) AND n = ? AND o IN (?,?,?) AND p = ? AND q IN (?,?))",
			"SELECT `a` FROM `c` WHERE (`h` IN (1,2,3)) AND (l = 'xx' AND m IN ('aa','bb','cc') AND n = 1 AND o IN (1,0,1) AND p = 'Gopher' AND q IN ('Go1','Go2'))",
			int64(1), int64(2), int64(3), "xx", "aa", "bb", "cc", true, true, false, true, []byte(`Gopher`), []byte(`Go1`), []byte(`Go2`),
		)
	})

	t.Run("null types", func(t *testing.T) {
		sel := NewSelect("a").From("c").
			Where(
				Column("h").In().Int64s(1, 2, 3),
				Expression("t1 = ? AND t2 IN (?) AND ns = ? OR nf = ? OR ni = ? OR nb = ? AND nt = ?").
					Time(now()).
					Times(now(), now()).
					NullString(MakeNullString("Goph3r")).
					NullFloat64(MakeNullFloat64(2.7182)).
					NullInt64(MakeNullInt64(27182)).
					NullBool(MakeNullBool(true)).
					NullTime(MakeNullTime(now())),
			)

		compareToSQL(t, sel, nil,
			"SELECT `a` FROM `c` WHERE (`h` IN (?,?,?)) AND (t1 = ? AND t2 IN (?,?) AND ns = ? OR nf = ? OR ni = ? OR nb = ? AND nt = ?)",
			"SELECT `a` FROM `c` WHERE (`h` IN (1,2,3)) AND (t1 = '2006-01-02 15:04:05' AND t2 IN ('2006-01-02 15:04:05','2006-01-02 15:04:05') AND ns = 'Goph3r' OR nf = 2.7182 OR ni = 27182 OR nb = 1 AND nt = '2006-01-02 15:04:05')",
			int64(1), int64(2), int64(3), now(), now(), now(), "Goph3r", 2.7182, int64(27182), true, now(),
		)
	})
}
