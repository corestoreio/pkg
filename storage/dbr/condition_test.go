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
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var _ fmt.Stringer = Op(0)

func TestOpRune(t *testing.T) {
	t.Parallel()
	s := NewSelect().From("tableA").AddColumns("a", "b").
		Where(
			Column("a1").Like().Str("H_ll_"),
			Column("a1").Like().NullString(NullString{}),
			Column("a1").Like().NullString(MakeNullString("NullString")),
			Column("a1").Like().Float64(2.718281),
			Column("a1").Like().NullFloat64(NullFloat64{}),
			Column("a1").Like().NullFloat64(MakeNullFloat64(-2.718281)),
			Column("a1").Like().Int64(2718281),
			Column("a1").Like().NullInt64(NullInt64{}),
			Column("a1").Like().NullInt64(MakeNullInt64(-987)),
			Column("a1").Like().Int(2718281),
			Column("a1").Like().Bool(true),
			Column("a1").Like().NullBool(NullBool{}),
			Column("a1").Like().NullBool(MakeNullBool(false)),
			Column("a1").Like().Time(now()),
			Column("a1").Like().NullTime(MakeNullTime(now().Add(time.Minute))),
			Column("a1").Like().Null(),
			Column("a1").Like().Bytes([]byte(`H3llo`)),
			Column("a1").Like().DriverValue(MakeNullInt64(2345)),

			Column("a2").NotLike().Str("H_ll_"),
			Column("a2").NotLike().NullString(NullString{}),
			Column("a2").NotLike().NullString(MakeNullString("NullString")),
			Column("a2").NotLike().Float64(2.718281),
			Column("a2").NotLike().NullFloat64(NullFloat64{}),
			Column("a2").NotLike().NullFloat64(MakeNullFloat64(-2.718281)),
			Column("a2").NotLike().Int64(2718281),
			Column("a2").NotLike().NullInt64(NullInt64{}),
			Column("a2").NotLike().NullInt64(MakeNullInt64(-987)),
			Column("a2").NotLike().Int(2718281),
			Column("a2").NotLike().Bool(true),
			Column("a2").NotLike().NullBool(NullBool{}),
			Column("a2").NotLike().NullBool(MakeNullBool(false)),
			Column("a2").NotLike().Time(now()),
			Column("a2").NotLike().NullTime(MakeNullTime(now().Add(time.Minute))),
			Column("a2").NotLike().Null(),
			Column("a2").NotLike().Bytes([]byte(`H3llo`)),
			Column("a2").NotLike().DriverValue(MakeNullInt64(2345)),

			Column("a301").In().Strs("Go1", "Go2"),
			Column("a302").In().NullString(NullString{}, NullString{}),
			Column("a303").In().NullString(MakeNullString("NullString")),
			Column("a304").In().Float64s(2.718281, 3.14159),
			Column("a305").In().NullFloat64(NullFloat64{}),
			Column("a306").In().NullFloat64(MakeNullFloat64(-2.718281), MakeNullFloat64(-3.14159)),
			Column("a307").In().Int64s(2718281, 314159),
			Column("a308").In().NullInt64(NullInt64{}),
			Column("a309").In().NullInt64(MakeNullInt64(-987), MakeNullInt64(-654)),
			Column("a310").In().Ints(2718281, 314159),
			Column("a311").In().Bools(true, false),
			Column("a312").In().NullBool(NullBool{}),
			Column("a313").In().NullBool(MakeNullBool(true)),
			Column("a314").In().Times(now(), now()),
			Column("a315").In().NullTime(MakeNullTime(now().Add(time.Minute))),
			Column("a316").In().Null(),
			Column("a317").In().Bytes([]byte(`H3llo1`)),
			Column("a320").In().DriverValue(MakeNullFloat64(674589), MakeNullFloat64(3.14159)),

			Column("a401").SpaceShip().Str("H_ll_"),
			Column("a402").SpaceShip().NullString(NullString{}),
			Column("a403").SpaceShip().NullString(MakeNullString("NullString")),
		)
	compareToSQL(t, s, nil,
		"SELECT `a`, `b` FROM `tableA` WHERE (`a1` LIKE ?) AND (`a1` LIKE ?) AND (`a1` LIKE ?) AND (`a1` LIKE ?) AND (`a1` LIKE ?) AND (`a1` LIKE ?) AND (`a1` LIKE ?) AND (`a1` LIKE ?) AND (`a1` LIKE ?) AND (`a1` LIKE ?) AND (`a1` LIKE ?) AND (`a1` LIKE ?) AND (`a1` LIKE ?) AND (`a1` LIKE ?) AND (`a1` LIKE ?) AND (`a1` IS NULL) AND (`a1` LIKE ?) AND (`a1` LIKE ?) AND (`a2` NOT LIKE ?) AND (`a2` NOT LIKE ?) AND (`a2` NOT LIKE ?) AND (`a2` NOT LIKE ?) AND (`a2` NOT LIKE ?) AND (`a2` NOT LIKE ?) AND (`a2` NOT LIKE ?) AND (`a2` NOT LIKE ?) AND (`a2` NOT LIKE ?) AND (`a2` NOT LIKE ?) AND (`a2` NOT LIKE ?) AND (`a2` NOT LIKE ?) AND (`a2` NOT LIKE ?) AND (`a2` NOT LIKE ?) AND (`a2` NOT LIKE ?) AND (`a2` IS NULL) AND (`a2` NOT LIKE ?) AND (`a2` NOT LIKE ?) AND (`a301` IN (?,?)) AND (`a302` IN (?,?)) AND (`a303` IN (?)) AND (`a304` IN (?,?)) AND (`a305` IN (?)) AND (`a306` IN (?,?)) AND (`a307` IN (?,?)) AND (`a308` IN (?)) AND (`a309` IN (?,?)) AND (`a310` IN (?,?)) AND (`a311` IN (?,?)) AND (`a312` IN (?)) AND (`a313` IN (?)) AND (`a314` IN (?,?)) AND (`a315` IN (?)) AND (`a316` IS NULL) AND (`a317` IN (?)) AND (`a320` IN (?,?)) AND (`a401` <=> ?) AND (`a402` <=> ?) AND (`a403` <=> ?)",
		"SELECT `a`, `b` FROM `tableA` WHERE (`a1` LIKE 'H_ll_') AND (`a1` LIKE NULL) AND (`a1` LIKE 'NullString') AND (`a1` LIKE 2.718281) AND (`a1` LIKE NULL) AND (`a1` LIKE -2.718281) AND (`a1` LIKE 2718281) AND (`a1` LIKE NULL) AND (`a1` LIKE -987) AND (`a1` LIKE 2718281) AND (`a1` LIKE 1) AND (`a1` LIKE NULL) AND (`a1` LIKE 0) AND (`a1` LIKE '2006-01-02 15:04:05') AND (`a1` LIKE '2006-01-02 15:05:05') AND (`a1` IS NULL) AND (`a1` LIKE 'H3llo') AND (`a1` LIKE 2345) AND (`a2` NOT LIKE 'H_ll_') AND (`a2` NOT LIKE NULL) AND (`a2` NOT LIKE 'NullString') AND (`a2` NOT LIKE 2.718281) AND (`a2` NOT LIKE NULL) AND (`a2` NOT LIKE -2.718281) AND (`a2` NOT LIKE 2718281) AND (`a2` NOT LIKE NULL) AND (`a2` NOT LIKE -987) AND (`a2` NOT LIKE 2718281) AND (`a2` NOT LIKE 1) AND (`a2` NOT LIKE NULL) AND (`a2` NOT LIKE 0) AND (`a2` NOT LIKE '2006-01-02 15:04:05') AND (`a2` NOT LIKE '2006-01-02 15:05:05') AND (`a2` IS NULL) AND (`a2` NOT LIKE 'H3llo') AND (`a2` NOT LIKE 2345) AND (`a301` IN ('Go1','Go2')) AND (`a302` IN (NULL,NULL)) AND (`a303` IN ('NullString')) AND (`a304` IN (2.718281,3.14159)) AND (`a305` IN (NULL)) AND (`a306` IN (-2.718281,-3.14159)) AND (`a307` IN (2718281,314159)) AND (`a308` IN (NULL)) AND (`a309` IN (-987,-654)) AND (`a310` IN (2718281,314159)) AND (`a311` IN (1,0)) AND (`a312` IN (NULL)) AND (`a313` IN (1)) AND (`a314` IN ('2006-01-02 15:04:05','2006-01-02 15:04:05')) AND (`a315` IN ('2006-01-02 15:05:05')) AND (`a316` IS NULL) AND (`a317` IN ('H3llo1')) AND (`a320` IN (674589,3.14159)) AND (`a401` <=> 'H_ll_') AND (`a402` <=> NULL) AND (`a403` <=> 'NullString')",
		"H_ll_", interface{}(nil), "NullString", 2.718281, interface{}(nil),
		-2.718281, int64(2718281), interface{}(nil), int64(-987), int64(2718281), true,
		interface{}(nil), false, now(), now().Add(time.Minute),
		[]uint8{0x48, 0x33, 0x6c, 0x6c, 0x6f}, int64(2345), "H_ll_",
		interface{}(nil), "NullString", 2.718281, interface{}(nil), -2.718281, int64(2718281),
		interface{}(nil), int64(-987), int64(2718281), true, interface{}(nil), false, now(), now().Add(time.Minute),
		[]uint8{0x48, 0x33, 0x6c, 0x6c, 0x6f}, int64(2345),
		"Go1", "Go2", interface{}(nil), interface{}(nil), "NullString", 2.718281, 3.14159,
		interface{}(nil), -2.718281, -3.14159, int64(2718281), int64(314159), interface{}(nil),
		int64(-987), int64(-654), int64(2718281), int64(314159), true, false, interface{}(nil), true,
		now(), now(), now().Add(time.Minute), []uint8{0x48, 0x33, 0x6c, 0x6c, 0x6f, 0x31},
		float64(674589), 3.14159, "H_ll_", interface{}(nil), "NullString",
	)
}

func TestOp_String(t *testing.T) {
	t.Parallel()
	var o Op
	assert.Exactly(t, "=", o.String())
	assert.Exactly(t, "🚀", SpaceShip.String())
}

func TestOpArgs(t *testing.T) {
	t.Parallel()
	t.Run("Null with place holder IN,Regexp", func(t *testing.T) {
		compareToSQL(t,
			NewSelect("a", "b").From("t1").Where(
				Column("a315").In().Null(),
				Column("a316").In().PlaceHolder(),
				Column("a317").Regexp().PlaceHolder(),
				Column("a317").NotRegexp().PlaceHolder(),
			),
			nil,
			"SELECT `a`, `b` FROM `t1` WHERE (`a315` IS NULL) AND (`a316` IN (?)) AND (`a317` REGEXP ?) AND (`a317` NOT REGEXP ?)",
			"SELECT `a`, `b` FROM `t1` WHERE (`a315` IS NULL) AND (`a316` IN (?)) AND (`a317` REGEXP ?) AND (`a317` NOT REGEXP ?)",
			[]interface{}{}...,
		)
	})

	t.Run("Args In", func(t *testing.T) {
		compareToSQL(t,
			NewSelect("a", "b").From("t1").Where(
				Column("a311").Xor().Int(9),
				Column("a313").In().Float64(3.3),
				Column("a314").In().Int64(33),
				Column("a312").In().Int(44),
				Column("a315").In().Str(`Go1`),
				Column("a316").In().BytesSlice([]byte(`Go`), []byte(`Rust`)),
			),
			nil,
			"SELECT `a`, `b` FROM `t1` WHERE (`a311` XOR ?) AND (`a313` IN (?)) AND (`a314` IN (?)) AND (`a312` IN (?)) AND (`a315` IN (?)) AND (`a316` IN (?,?))",
			"SELECT `a`, `b` FROM `t1` WHERE (`a311` XOR 9) AND (`a313` IN (3.3)) AND (`a314` IN (33)) AND (`a312` IN (44)) AND (`a315` IN ('Go1')) AND (`a316` IN ('Go','Rust'))",
			int64(9), 3.3, int64(33), int64(44), "Go1", []uint8{0x47, 0x6f}, []uint8{0x52, 0x75, 0x73, 0x74},
		)
	})

	t.Run("BytesSlice BETWEEN strings", func(t *testing.T) {
		compareToSQL(t,
			NewSelect("a", "b").From("t1").Where(
				Column("a316").Between().BytesSlice([]byte(`Go`), []byte(`Rust`)),
			),
			nil,
			"SELECT `a`, `b` FROM `t1` WHERE (`a316` BETWEEN ? AND ?)",
			"SELECT `a`, `b` FROM `t1` WHERE (`a316` BETWEEN 'Go' AND 'Rust')",
			[]uint8{0x47, 0x6f}, []uint8{0x52, 0x75, 0x73, 0x74},
		)
	})
	t.Run("BytesSlice IN binary", func(t *testing.T) {

		compareToSQL(t,
			NewSelect("a", "b").From("t1").Where(
				Column("a316").In().BytesSlice([]byte{66, 250, 67}, []byte(`Rust`), []byte("\xFB\xBF\xBF\xBF\xBF")),
			),
			nil,
			"SELECT `a`, `b` FROM `t1` WHERE (`a316` IN (?,?,?))",
			"SELECT `a`, `b` FROM `t1` WHERE (`a316` IN (0x42fa43,'Rust',0xfbbfbfbfbf))",
			[]uint8{0x42, 0xfa, 0x43}, []uint8{0x52, 0x75, 0x73, 0x74}, []uint8{0xfb, 0xbf, 0xbf, 0xbf, 0xbf},
		)
	})
	t.Run("ArgValue IN", func(t *testing.T) {

		compareToSQL(t,
			NewSelect("a", "b").From("t1").Where(
				Column("a3419").In().DriverValues(
					MakeNullFloat64(3.141),
					MakeNullString("G'o"),
					driverValueBytes{66, 250, 67},
					MakeNullTime(now()),
					driverValueBytes([]byte("x\x00\xff")),
				),
			),
			nil,
			"SELECT `a`, `b` FROM `t1` WHERE (`a3419` IN (?,?,?,?,?))",
			"SELECT `a`, `b` FROM `t1` WHERE (`a3419` IN (3.141,'G\\'o',0x42fa43,'2006-01-02 15:04:05',0x7800ff))",
			3.141, `G'o`, []uint8{0x42, 0xfa, 0x43}, now(), []uint8{0x78, 0x0, 0xff},
		)
	})
	t.Run("ArgValue BETWEEN values", func(t *testing.T) {
		compareToSQL(t,
			NewSelect("a", "b").From("t1").Where(
				Column("a319").Between().DriverValues(MakeNullFloat64(3.141), MakeNullString("G'o")),
			),
			nil,
			"SELECT `a`, `b` FROM `t1` WHERE (`a319` BETWEEN ? AND ?)",
			"SELECT `a`, `b` FROM `t1` WHERE (`a319` BETWEEN 3.141 AND 'G\\'o')",
			3.141, `G'o`,
		)
	})
}

func TestColumn(t *testing.T) {
	t.Run("invalid column name", func(t *testing.T) {
		s := NewSelect("a", "b").From("c").Where(
			Column("a").Int(111),
			Expr("b=c"),
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
			args := MakeArgs(2)
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
		Column("name").Str("E0S 5D Mark II"),
		Column("stock").Int64(12),
	}, " ON DUPLICATE KEY UPDATE `name`=?, `stock`=?",
		"E0S 5D Mark II", int64(12)))

	t.Run("col1=VALUES(val)+? and with arguments", runner(Conditions{
		Column("name").Str("E0S 5D Mark II"),
		Column("stock").Expr("VALUES(`stock`)+?-?").Int64(13).Int(4),
	}, " ON DUPLICATE KEY UPDATE `name`=?, `stock`=VALUES(`stock`)+?-?",
		"E0S 5D Mark II", int64(13), int64(4)))

	t.Run("col2=VALUES(val) and with arguments and nil", runner(Conditions{
		Column("name").Str("E0S 5D Mark III"),
		Column("sku").Values(),
		Column("stock").Int64(14),
	}, " ON DUPLICATE KEY UPDATE `name`=?, `sku`=VALUES(`sku`), `stock`=?",
		"E0S 5D Mark III", int64(14)))

	t.Run("col=expression without arguments", runner(Conditions{
		Column("name").Expr("CONCAT('Canon','E0S 5D Mark III')"),
	}, " ON DUPLICATE KEY UPDATE `name`=CONCAT('Canon','E0S 5D Mark III')",
	))
}

func TestExpr_Arguments(t *testing.T) {

	t.Run("ints", func(t *testing.T) {
		sel := NewSelect("a").From("c").
			Where(
				Column("g").Int(3),
				Expr("i1 = ? AND i2 IN (?) AND i64_1 = ? AND i64_2 IN (?) AND ui64 > ? AND f64_1 = ? AND f64_2 IN (?)").
					Int(1).Ints(2, 3).
					Int64(4).Int64s(5, 6).
					Uint64(7).
					Float64(4.51).Float64s(5.41, 6.66666),
			)

		compareToSQL(t, sel, nil,
			"SELECT `a` FROM `c` WHERE (`g` = ?) AND (i1 = ? AND i2 IN (?,?) AND i64_1 = ? AND i64_2 IN (?,?) AND ui64 > ? AND f64_1 = ? AND f64_2 IN (?,?))",
			"SELECT `a` FROM `c` WHERE (`g` = 3) AND (i1 = 1 AND i2 IN (2,3) AND i64_1 = 4 AND i64_2 IN (5,6) AND ui64 > 7 AND f64_1 = 4.51 AND f64_2 IN (5.41,6.66666))",
			int64(3), int64(1), int64(2), int64(3), int64(4), int64(5), int64(6), int64(7), 4.51, 5.41, 6.66666,
		)
	})

	t.Run("slice expression", func(t *testing.T) {
		sel := NewSelect("a").From("c").
			Where(
				Column("h").In().Int64s(1, 2, 3),
				Expr("l NOT IN (?)").Strs("xx", "yy"),
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
				Expr("l = ? AND m IN (?) AND n = ? AND o IN (?) AND p = ? AND q IN (?)").
					Str("xx").Strs("aa", "bb", "cc").
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
				Expr("t1 = ? AND t2 IN (?) AND ns = ? OR nf = ? OR ni = ? OR nb = ? AND nt = ?").
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

func TestCondition_Column(t *testing.T) {
	t.Parallel()
	sel := NewSelect("t_d.attribute_id", "e.entity_id").
		AddColumnsAliases("t_d.value", "default_value").
		AddColumnsConditions(SQLIf("t_s.value_id IS NULL", "t_d.value", "t_s.value").Alias("value")).
		AddColumnsConditions(SQLIf("? IS NULL", "t_d.value", "t_s.value").NullFloat64(MakeNullFloat64(2.718281)).Alias("value")).
		FromAlias("catalog_category_entity", "e").
		Join(
			MakeIdentifier("catalog_category_entity_varchar").Alias("t_d"), // t_d = table scope default
			Column("e.entity_id").Equal().Column("t_d.entity_id"),
		).
		LeftJoin(
			MakeIdentifier("catalog_category_entity_varchar").Alias("t_s"), // t_s = table scope store
			Column("t_s.attribute_id").GreaterOrEqual().Column("t_d.attribute_id"),
		).
		Where(
			Column("e.entity_id").In().Int64s(28, 16, 25, 17),
			Column("t_d.attribute_id").In().Int64s(45),
			Column("t_d.store_id").Equal().SQLIfNull("t_s.store_id", "0"),
		)

	compareToSQL(t, sel, nil,
		"SELECT `t_d`.`attribute_id`, `e`.`entity_id`, `t_d`.`value` AS `default_value`, IF((t_s.value_id IS NULL), t_d.value, t_s.value) AS `value`, IF((? IS NULL), t_d.value, t_s.value) AS `value` FROM `catalog_category_entity` AS `e` INNER JOIN `catalog_category_entity_varchar` AS `t_d` ON (`e`.`entity_id` = `t_d`.`entity_id`) LEFT JOIN `catalog_category_entity_varchar` AS `t_s` ON (`t_s`.`attribute_id` >= `t_d`.`attribute_id`) WHERE (`e`.`entity_id` IN (?,?,?,?)) AND (`t_d`.`attribute_id` IN (?)) AND (`t_d`.`store_id` = IFNULL(`t_s`.`store_id`,0))",
		"SELECT `t_d`.`attribute_id`, `e`.`entity_id`, `t_d`.`value` AS `default_value`, IF((t_s.value_id IS NULL), t_d.value, t_s.value) AS `value`, IF((2.718281 IS NULL), t_d.value, t_s.value) AS `value` FROM `catalog_category_entity` AS `e` INNER JOIN `catalog_category_entity_varchar` AS `t_d` ON (`e`.`entity_id` = `t_d`.`entity_id`) LEFT JOIN `catalog_category_entity_varchar` AS `t_s` ON (`t_s`.`attribute_id` >= `t_d`.`attribute_id`) WHERE (`e`.`entity_id` IN (28,16,25,17)) AND (`t_d`.`attribute_id` IN (45)) AND (`t_d`.`store_id` = IFNULL(`t_s`.`store_id`,0))",
		2.718281, int64(28), int64(16), int64(25), int64(17), int64(45),
	)
}

func TestExpr(t *testing.T) {
	t.Parallel()
	t.Run("quoted string", func(t *testing.T) {
		s := NewSelect().AddColumns("month", "total").AddColumnsConditions(Expr(`"best"`)).From("sales_by_month")
		compareToSQL(t, s, nil,
			"SELECT `month`, `total`, \"best\" FROM `sales_by_month`",
			"",
		)
	})
}

func TestSplitColumn(t *testing.T) {
	t.Parallel()

	tests := []struct {
		identifier string
		wantQuali  string
		wantCol    string
	}{
		{"id", "", "id"},
		{".id", "", ".id"},
		{".id.", "", ".id."},
		{"id.", "", "id."},
		{"cpe.entity_id", "cpe", "entity_id"},
		{"cpe.*", "cpe", "*"},
		{"database.cpe.entity_id", "database.cpe", "entity_id"},
	}
	for i, test := range tests {
		haveQ, haveC := splitColumn(test.identifier)
		assert.Exactly(t, test.wantQuali, haveQ, "Qualifier mismatch at index %d", i)
		assert.Exactly(t, test.wantCol, haveC, "Column mismatch at index %d", i)
	}
}

type appendInt int

func (ai appendInt) AppendBind(args Arguments, _ []string) (Arguments, error) {
	return args.Int(int(ai)), nil
}

func TestAppendArgs(t *testing.T) {
	t.Parallel()
	t.Run("PH,val,expr,PH", func(t *testing.T) {
		s := NewSelect("sku").FromAlias("catalog", "e").
			// alias t_d ignored and not needed in this test case
			Where(
				Column("e.entity_id").In().PlaceHolder(),                      // 678
				Column("t_d.attribute_id").In().Int64s(45),                    // 45
				Column("t_d.store_id").Equal().SQLIfNull("t_s.store_id", "0"), // Does not make sense this WHERE condition ;-)
				Column("t_d.store_id").Equal().PlaceHolder(),                  // 17
			).
			BindByQualifier("e", appendInt(678)).
			BindByQualifier("t_d", appendInt(17))

		compareToSQL(t, s, nil,
			"SELECT `sku` FROM `catalog` AS `e` WHERE (`e`.`entity_id` IN (?)) AND (`t_d`.`attribute_id` IN (?)) AND (`t_d`.`store_id` = IFNULL(`t_s`.`store_id`,0)) AND (`t_d`.`store_id` = ?)",
			"SELECT `sku` FROM `catalog` AS `e` WHERE (`e`.`entity_id` IN (678)) AND (`t_d`.`attribute_id` IN (45)) AND (`t_d`.`store_id` = IFNULL(`t_s`.`store_id`,0)) AND (`t_d`.`store_id` = 17)",
			int64(678), int64(45), int64(17),
		)
	})

	t.Run("PH,val,PH", func(t *testing.T) {
		s := NewSelect("sku").FromAlias("catalog", "e").
			// alias t_d ignored and not needed in this test case
			Where(
				Column("e.entity_id").In().PlaceHolder(),     // 678
				Column("t_d.attribute_id").In().Int64s(45),   // 45
				Column("t_d.store_id").Equal().PlaceHolder(), // 17
			).
			BindByQualifier("e", appendInt(678)).
			BindByQualifier("t_d", appendInt(17))

		compareToSQL(t, s, nil,
			"SELECT `sku` FROM `catalog` AS `e` WHERE (`e`.`entity_id` IN (?)) AND (`t_d`.`attribute_id` IN (?)) AND (`t_d`.`store_id` = ?)",
			"SELECT `sku` FROM `catalog` AS `e` WHERE (`e`.`entity_id` IN (678)) AND (`t_d`.`attribute_id` IN (45)) AND (`t_d`.`store_id` = 17)",
			int64(678), int64(45), int64(17),
		)
	})

}
