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

package dml

import (
	"bytes"
	"fmt"
	"testing"
	"time"

	"github.com/corestoreio/errors"
	"github.com/corestoreio/pkg/storage/null"
	"github.com/corestoreio/pkg/util/assert"
)

var _ fmt.Stringer = Op(0)

func TestConditionOpRune(t *testing.T) {
	t.Parallel()
	s := NewSelect().From("tableA").AddColumns("a", "b").
		Where(
			Column("a1").Like().Str("H_ll_"),
			Column("a1").Like().NullString(null.String{}),
			Column("a1").Like().NullString(null.MakeString("NullString")),
			Column("a1").Like().Float64(2.718281),
			Column("a1").Like().NullFloat64(null.Float64{}),
			Column("a1").Like().NullFloat64(null.MakeFloat64(-2.718281)),
			Column("a1").Like().Int64(2718281),
			Column("a1").Like().NullInt64(null.Int64{}),
			Column("a1").Like().NullInt64(null.MakeInt64(-987)),
			Column("a1").Like().Int(2718281),
			Column("a1").Like().Bool(true),
			Column("a1").Like().NullBool(null.Bool{}),
			Column("a1").Like().NullBool(null.MakeBool(false)),
			Column("a1").Like().Time(now()),
			Column("a1").Like().NullTime(null.MakeTime(now().Add(time.Minute))),
			Column("a1").Like().Null(),
			Column("a1").Like().Bytes([]byte(`H3llo`)),
			Column("a1").Like().DriverValue(null.MakeInt64(2345)),

			Column("a2").NotLike().Str("H_ll_"),
			Column("a2").NotLike().NullString(null.String{}),
			Column("a2").NotLike().NullString(null.MakeString("NullString")),
			Column("a2").NotLike().Float64(2.718281),
			Column("a2").NotLike().NullFloat64(null.Float64{}),
			Column("a2").NotLike().NullFloat64(null.MakeFloat64(-2.718281)),
			Column("a2").NotLike().Int64(2718281),
			Column("a2").NotLike().NullInt64(null.Int64{}),
			Column("a2").NotLike().NullInt64(null.MakeInt64(-987)),
			Column("a2").NotLike().Int(2718281),
			Column("a2").NotLike().Bool(true),
			Column("a2").NotLike().NullBool(null.Bool{}),
			Column("a2").NotLike().NullBool(null.MakeBool(false)),
			Column("a2").NotLike().Time(now()),
			Column("a2").NotLike().NullTime(null.MakeTime(now().Add(time.Minute))),
			Column("a2").NotLike().Null(),
			Column("a2").NotLike().Bytes([]byte(`H3llo`)),
			Column("a2").NotLike().DriverValue(null.MakeInt64(2345)),

			Column("a301").In().Strs("Go1", "Go2"),
			Column("a303").In().NullString(null.MakeString("NullXString")),
			Column("a302").In().NullStrings(null.String{}, null.String{}),
			Column("a304").In().Float64s(2.718281, 3.14159),
			Column("a305").In().NullFloat64(null.Float64{}),
			Column("a306").In().NullFloat64s(null.MakeFloat64(-2.718281), null.MakeFloat64(-3.14159)),
			Column("a307").In().Int64s(2718281, 314159),
			Column("a308").In().NullInt64(null.Int64{}),
			Column("a309").In().NullInt64s(null.MakeInt64(-987), null.MakeInt64(-654)),
			Column("a310").In().Ints(2718281, 314159),
			Column("a311").In().Bools(true, false),
			Column("a312").In().NullBool(null.Bool{}),
			Column("a313").In().NullBools(null.MakeBool(true)),
			Column("a314").In().Times(now(), now()),
			Column("a315a").In().NullTime(null.MakeTime(now().Add(time.Minute))),
			Column("a315b").In().NullTimes(null.MakeTime(now().Add(time.Minute)), null.MakeTime(now().Add(time.Minute*2))),
			Column("a316").In().Null(),
			Column("a317").In().Bytes([]byte(`H3llo1`)),
			Column("a320").In().DriverValue(null.MakeFloat64(674589), null.MakeFloat64(3.14159)),

			Column("a401").SpaceShip().Str("H_ll_"),
			Column("a402").SpaceShip().NullString(null.String{}),
			Column("a403").SpaceShip().NullString(null.MakeString("NullString")),
		)
	compareToSQL(t, s, errors.NoKind,
		"SELECT `a`, `b` FROM `tableA` WHERE (`a1` LIKE 'H_ll_') AND (`a1` LIKE NULL) AND (`a1` LIKE 'NullString') AND (`a1` LIKE 2.718281) AND (`a1` LIKE NULL) AND (`a1` LIKE -2.718281) AND (`a1` LIKE 2718281) AND (`a1` LIKE NULL) AND (`a1` LIKE -987) AND (`a1` LIKE 2718281) AND (`a1` LIKE 1) AND (`a1` LIKE NULL) AND (`a1` LIKE 0) AND (`a1` LIKE '2006-01-02 15:04:05') AND (`a1` LIKE '2006-01-02 15:05:05') AND (`a1` IS NULL) AND (`a1` LIKE 'H3llo') AND (`a1` LIKE (2345)) AND (`a2` NOT LIKE 'H_ll_') AND (`a2` NOT LIKE NULL) AND (`a2` NOT LIKE 'NullString') AND (`a2` NOT LIKE 2.718281) AND (`a2` NOT LIKE NULL) AND (`a2` NOT LIKE -2.718281) AND (`a2` NOT LIKE 2718281) AND (`a2` NOT LIKE NULL) AND (`a2` NOT LIKE -987) AND (`a2` NOT LIKE 2718281) AND (`a2` NOT LIKE 1) AND (`a2` NOT LIKE NULL) AND (`a2` NOT LIKE 0) AND (`a2` NOT LIKE '2006-01-02 15:04:05') AND (`a2` NOT LIKE '2006-01-02 15:05:05') AND (`a2` IS NULL) AND (`a2` NOT LIKE 'H3llo') AND (`a2` NOT LIKE (2345)) AND (`a301` IN ('Go1','Go2')) AND (`a303` IN 'NullXString') AND (`a302` IN (NULL,NULL)) AND (`a304` IN (2.718281,3.14159)) AND (`a305` IN NULL) AND (`a306` IN (-2.718281,-3.14159)) AND (`a307` IN (2718281,314159)) AND (`a308` IN NULL) AND (`a309` IN (-987,-654)) AND (`a310` IN (2718281,314159)) AND (`a311` IN (1,0)) AND (`a312` IN NULL) AND (`a313` IN (1)) AND (`a314` IN ('2006-01-02 15:04:05','2006-01-02 15:04:05')) AND (`a315a` IN '2006-01-02 15:05:05') AND (`a315b` IN ('2006-01-02 15:05:05','2006-01-02 15:06:05')) AND (`a316` IS NULL) AND (`a317` IN 'H3llo1') AND (`a320` IN (674589,3.14159)) AND (`a401` <=> 'H_ll_') AND (`a402` <=> NULL) AND (`a403` <=> 'NullString')",
		"SELECT `a`, `b` FROM `tableA` WHERE (`a1` LIKE 'H_ll_') AND (`a1` LIKE NULL) AND (`a1` LIKE 'NullString') AND (`a1` LIKE 2.718281) AND (`a1` LIKE NULL) AND (`a1` LIKE -2.718281) AND (`a1` LIKE 2718281) AND (`a1` LIKE NULL) AND (`a1` LIKE -987) AND (`a1` LIKE 2718281) AND (`a1` LIKE 1) AND (`a1` LIKE NULL) AND (`a1` LIKE 0) AND (`a1` LIKE '2006-01-02 15:04:05') AND (`a1` LIKE '2006-01-02 15:05:05') AND (`a1` IS NULL) AND (`a1` LIKE 'H3llo') AND (`a1` LIKE (2345)) AND (`a2` NOT LIKE 'H_ll_') AND (`a2` NOT LIKE NULL) AND (`a2` NOT LIKE 'NullString') AND (`a2` NOT LIKE 2.718281) AND (`a2` NOT LIKE NULL) AND (`a2` NOT LIKE -2.718281) AND (`a2` NOT LIKE 2718281) AND (`a2` NOT LIKE NULL) AND (`a2` NOT LIKE -987) AND (`a2` NOT LIKE 2718281) AND (`a2` NOT LIKE 1) AND (`a2` NOT LIKE NULL) AND (`a2` NOT LIKE 0) AND (`a2` NOT LIKE '2006-01-02 15:04:05') AND (`a2` NOT LIKE '2006-01-02 15:05:05') AND (`a2` IS NULL) AND (`a2` NOT LIKE 'H3llo') AND (`a2` NOT LIKE (2345)) AND (`a301` IN ('Go1','Go2')) AND (`a303` IN 'NullXString') AND (`a302` IN (NULL,NULL)) AND (`a304` IN (2.718281,3.14159)) AND (`a305` IN NULL) AND (`a306` IN (-2.718281,-3.14159)) AND (`a307` IN (2718281,314159)) AND (`a308` IN NULL) AND (`a309` IN (-987,-654)) AND (`a310` IN (2718281,314159)) AND (`a311` IN (1,0)) AND (`a312` IN NULL) AND (`a313` IN (1)) AND (`a314` IN ('2006-01-02 15:04:05','2006-01-02 15:04:05')) AND (`a315a` IN '2006-01-02 15:05:05') AND (`a315b` IN ('2006-01-02 15:05:05','2006-01-02 15:06:05')) AND (`a316` IS NULL) AND (`a317` IN 'H3llo1') AND (`a320` IN (674589,3.14159)) AND (`a401` <=> 'H_ll_') AND (`a402` <=> NULL) AND (`a403` <=> 'NullString')",
	)
}

func TestConditionOp_String(t *testing.T) {
	t.Parallel()
	var o Op
	assert.Exactly(t, "=", o.String())
	assert.Exactly(t, "🚀", SpaceShip.String())
}

func TestConditionOpArgs(t *testing.T) {
	t.Parallel()
	t.Run("Null with place holder IN,Regexp", func(t *testing.T) {
		compareToSQL(t,
			NewSelect("a", "b").From("t1").Where(
				Column("a315").In().Null(),
				Column("a316").In().PlaceHolder(),
				Column("a317").Regexp().PlaceHolder(),
				Column("a318").NotRegexp().PlaceHolder(),
			),
			errors.NoKind,
			"SELECT `a`, `b` FROM `t1` WHERE (`a315` IS NULL) AND (`a316` IN ?) AND (`a317` REGEXP ?) AND (`a318` NOT REGEXP ?)",
			"",
		)
	})
	t.Run("IN place holder and one value", func(t *testing.T) {
		compareToSQL(t,
			NewSelect("a", "b").From("t1").Where(
				Column("a315").In().Null(),
				Column("a316").In().PlaceHolder(),
			).WithDBR().TestWithArgs([]string{"aa"}),
			errors.NoKind,
			"SELECT `a`, `b` FROM `t1` WHERE (`a315` IS NULL) AND (`a316` IN ?)",
			"SELECT `a`, `b` FROM `t1` WHERE (`a315` IS NULL) AND (`a316` IN ('aa'))",
			"aa",
		)
	})
	t.Run("IN place holder and two values", func(t *testing.T) {
		compareToSQL(t,
			NewSelect("a", "b").From("t1").Where(
				Column("a315").In().Null(),
				Column("a316").In().PlaceHolder(),
			).WithDBR().TestWithArgs([]string{"aa", "bb"}),
			errors.NoKind,
			"SELECT `a`, `b` FROM `t1` WHERE (`a315` IS NULL) AND (`a316` IN ?)",
			"SELECT `a`, `b` FROM `t1` WHERE (`a315` IS NULL) AND (`a316` IN ('aa','bb'))",
			"aa", "bb",
		)
	})

	t.Run("Args In (single)", func(t *testing.T) {
		compareToSQL(t,
			NewSelect("a", "b").From("t1").Where(
				Column("a311").Xor().Int(9),
				Column("a313").In().Float64(3.3), // Wrong SQL
				Column("a314").In().Int64(33),    // Wrong SQL
				Column("a312").In().Int(44),      // Wrong SQL
				Column("a315").In().Str(`Go1`),   // Wrong SQL
				Column("a316").In().BytesSlice([]byte(`Go`), []byte(`Rust`)),
			),
			errors.NoKind,
			"SELECT `a`, `b` FROM `t1` WHERE (`a311` XOR 9) AND (`a313` IN 3.3) AND (`a314` IN 33) AND (`a312` IN 44) AND (`a315` IN 'Go1') AND (`a316` IN ('Go','Rust'))",
			"SELECT `a`, `b` FROM `t1` WHERE (`a311` XOR 9) AND (`a313` IN 3.3) AND (`a314` IN 33) AND (`a312` IN 44) AND (`a315` IN 'Go1') AND (`a316` IN ('Go','Rust'))",
		)
	})
	t.Run("Args In (plural)", func(t *testing.T) {
		compareToSQL(t,
			NewSelect("a", "b").From("t1").Where(
				Column("a313").In().Float64s(3.3),
				Column("a314").In().Int64s(33),
				Column("a312").In().Ints(44),
				Column("a315").In().Strs(`Go1`),
				Column("a316").In().BytesSlice([]byte(`Go`), []byte(`Rust`)),
			),
			errors.NoKind,
			"SELECT `a`, `b` FROM `t1` WHERE (`a313` IN (3.3)) AND (`a314` IN (33)) AND (`a312` IN (44)) AND (`a315` IN ('Go1')) AND (`a316` IN ('Go','Rust'))",
			"SELECT `a`, `b` FROM `t1` WHERE (`a313` IN (3.3)) AND (`a314` IN (33)) AND (`a312` IN (44)) AND (`a315` IN ('Go1')) AND (`a316` IN ('Go','Rust'))",
		)
	})

	t.Run("BytesSlice BETWEEN strings", func(t *testing.T) {
		compareToSQL(t,
			NewSelect("a", "b").From("t1").Where(
				Column("a316").Between().BytesSlice([]byte(`Go`), []byte(`Rust`)),
			),
			errors.NoKind,
			"SELECT `a`, `b` FROM `t1` WHERE (`a316` BETWEEN 'Go' AND 'Rust')",
			"SELECT `a`, `b` FROM `t1` WHERE (`a316` BETWEEN 'Go' AND 'Rust')",
		)
	})
	t.Run("BytesSlice IN binary", func(t *testing.T) {
		compareToSQL(t,
			NewSelect("a", "b").From("t1").Where(
				Column("a316").In().BytesSlice([]byte{66, 250, 67}, []byte(`Rust`), []byte("\xFB\xBF\xBF\xBF\xBF")),
			),
			errors.NoKind,
			"SELECT `a`, `b` FROM `t1` WHERE (`a316` IN (0x42fa43,'Rust',0xfbbfbfbfbf))",
			"SELECT `a`, `b` FROM `t1` WHERE (`a316` IN (0x42fa43,'Rust',0xfbbfbfbfbf))",
		)
	})
	t.Run("ArgValue IN", func(t *testing.T) {
		compareToSQL(t,
			NewSelect("a", "b").From("t1").Where(
				Column("a3419").In().DriverValues(
					null.MakeFloat64(3.141),
					null.MakeString("G'o"),
					driverValueBytes{66, 250, 67},
					null.MakeTime(now()),
					driverValueBytes([]byte("x\x00\xff")),
				),
			),
			errors.NoKind,
			"SELECT `a`, `b` FROM `t1` WHERE (`a3419` IN (3.141,'G\\'o',0x42fa43,'2006-01-02 15:04:05',0x7800ff))",
			"SELECT `a`, `b` FROM `t1` WHERE (`a3419` IN (3.141,'G\\'o',0x42fa43,'2006-01-02 15:04:05',0x7800ff))",
		)
	})
	t.Run("ArgValue BETWEEN error", func(t *testing.T) {
		// Between statements do not work with DriverValue
		compareToSQL(t,
			NewSelect("a", "b").From("t1").Where(
				Column("a319").Between().DriverValues(null.MakeFloat64(3.141), null.MakeString("G'o")),
			),
			errors.NoKind,
			"SELECT `a`, `b` FROM `t1` WHERE (`a319` BETWEEN ? AND )",
			"",
		)
	})
}

func TestConditionColumn(t *testing.T) {
	t.Parallel()

	t.Run("invalid column name", func(t *testing.T) {
		s := NewSelect("a", "b").From("c").Where(
			Column("a").Int(111),
			Expr("b=c"),
		)
		compareToSQL2(t, s, errors.NoKind, "SELECT `a`, `b` FROM `c` WHERE (`a` = 111) AND (b=c)")
	})

	t.Run("valid column name", func(t *testing.T) {
		s := NewSelect("a", "b").From("c").Where(
			Column("a").Ints(111, 222), // omitted In(). on purpose because default operator is IN for slices
			Column("b").Null(),
			Column("d").Between().Float64s(2.5, 2.7),
		)
		compareToSQL2(t, s, errors.NoKind, "SELECT `a`, `b` FROM `c` WHERE (`a` IN (111,222)) AND (`b` IS NULL) AND (`d` BETWEEN 2.5 AND 2.7)")
	})
}

func TestConditions_writeOnDuplicateKey(t *testing.T) {
	t.Parallel()

	runner := func(cnds Conditions, wantSQL string, wantArgs ...interface{}) func(*testing.T) {
		return func(t *testing.T) {
			buf := new(bytes.Buffer)

			ph, err := cnds.writeOnDuplicateKey(buf, nil)
			assert.Nil(t, ph, "TODO check me")
			assert.NoError(t, err)
			// args := MakeArgs(2)
			// args, _, err = cnds.appendArgs(args, appendArgsDUPKEY)
			// assert.NoError(t, err)
			// assert.Exactly(t, wantSQL, buf.String(), t.Name())
			// assert.Exactly(t, wantArgs, args.expandInterfaces(), t.Name())
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

func TestConditionExpr_Arguments(t *testing.T) {
	t.Parallel()

	t.Run("ints", func(t *testing.T) {
		sel := NewSelect("a").From("c").
			Where(
				Column("g").Int(3),
				Expr("i1 = ? AND i2 IN ? AND i64_1 = ? AND i64_2 IN ? AND ui64 > ? AND f64_1 = ? AND f64_2 IN ?").
					Int(1).Ints(2, 3).
					Int64(4).Int64s(5, 6).
					Uint64(7).
					Float64(4.51).Float64s(5.41, 6.66666),
			)

		compareToSQL(t, sel, errors.NoKind,
			"SELECT `a` FROM `c` WHERE (`g` = 3) AND (i1 = 1 AND i2 IN (2,3) AND i64_1 = 4 AND i64_2 IN (5,6) AND ui64 > 7 AND f64_1 = 4.51 AND f64_2 IN (5.41,6.66666))",
			"SELECT `a` FROM `c` WHERE (`g` = 3) AND (i1 = 1 AND i2 IN (2,3) AND i64_1 = 4 AND i64_2 IN (5,6) AND ui64 > 7 AND f64_1 = 4.51 AND f64_2 IN (5.41,6.66666))",
		)
	})

	t.Run("slice expression", func(t *testing.T) {
		sel := NewSelect("a").From("c").
			Where(
				Column("h").In().Int64s(1, 2, 3),
				Expr("l NOT IN ?").Strs("xx", "yy"),
			)
		compareToSQL(t, sel, errors.NoKind,
			"SELECT `a` FROM `c` WHERE (`h` IN (1,2,3)) AND (l NOT IN ('xx','yy'))",
			"SELECT `a` FROM `c` WHERE (`h` IN (1,2,3)) AND (l NOT IN ('xx','yy'))",
		)
	})

	t.Run("string bools", func(t *testing.T) {
		sel := NewSelect("a").From("c").
			Where(
				Column("h").In().Int64s(1, 2, 3),
				Expr("l = ? AND m IN ? AND n = ? AND o IN ? AND p = ? AND q IN ?").
					Str("xx").Strs("aa", "bb", "cc").
					Bool(true).Bools(true, false, true).
					Bytes([]byte(`Gopher`)).BytesSlice([]byte(`Go1`), []byte(`Go2`)),
			)

		compareToSQL(t, sel, errors.NoKind,
			"SELECT `a` FROM `c` WHERE (`h` IN (1,2,3)) AND (l = 'xx' AND m IN ('aa','bb','cc') AND n = 1 AND o IN (1,0,1) AND p = 'Gopher' AND q IN ('Go1','Go2'))",
			"SELECT `a` FROM `c` WHERE (`h` IN (1,2,3)) AND (l = 'xx' AND m IN ('aa','bb','cc') AND n = 1 AND o IN (1,0,1) AND p = 'Gopher' AND q IN ('Go1','Go2'))",
		)
	})

	t.Run("null types", func(t *testing.T) {
		sel := NewSelect("a").From("c").
			Where(
				Column("h").In().Int64s(1, 2, 3),
				Expr("t1 = ? AND t2 IN ? AND ns = ? OR nf = ? OR ni = ? OR nb = ? AND nt = ?").
					Time(now()).
					Times(now(), now()).
					NullString(null.MakeString("Goph3r")).
					NullFloat64(null.MakeFloat64(2.7182)).
					NullInt64(null.MakeInt64(27182)).
					NullBool(null.MakeBool(true)).
					NullTime(null.MakeTime(now())),
			)

		compareToSQL(t, sel, errors.NoKind,
			"SELECT `a` FROM `c` WHERE (`h` IN (1,2,3)) AND (t1 = '2006-01-02 15:04:05' AND t2 IN ('2006-01-02 15:04:05','2006-01-02 15:04:05') AND ns = 'Goph3r' OR nf = 2.7182 OR ni = 27182 OR nb = 1 AND nt = '2006-01-02 15:04:05')",
			"SELECT `a` FROM `c` WHERE (`h` IN (1,2,3)) AND (t1 = '2006-01-02 15:04:05' AND t2 IN ('2006-01-02 15:04:05','2006-01-02 15:04:05') AND ns = 'Goph3r' OR nf = 2.7182 OR ni = 27182 OR nb = 1 AND nt = '2006-01-02 15:04:05')",
		)
	})
}

func TestCondition_Column(t *testing.T) {
	t.Parallel()

	t.Run("complex", func(t *testing.T) {
		sel := NewSelect("t_d.attribute_id", "e.entity_id").
			AddColumnsAliases("t_d.value", "default_value").
			AddColumnsConditions(SQLIf("t_s.value_id IS NULL", "t_d.value", "t_s.value").Alias("value")).
			AddColumnsConditions(SQLIf("? IS NULL", "t_d.value", "t_s.value").NullFloat64(null.MakeFloat64(2.718281)).Alias("value")).
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

		compareToSQL(t, sel, errors.NoKind,
			"SELECT `t_d`.`attribute_id`, `e`.`entity_id`, `t_d`.`value` AS `default_value`, IF((t_s.value_id IS NULL), t_d.value, t_s.value) AS `value`, IF((2.718281 IS NULL), t_d.value, t_s.value) AS `value` FROM `catalog_category_entity` AS `e` INNER JOIN `catalog_category_entity_varchar` AS `t_d` ON (`e`.`entity_id` = `t_d`.`entity_id`) LEFT JOIN `catalog_category_entity_varchar` AS `t_s` ON (`t_s`.`attribute_id` >= `t_d`.`attribute_id`) WHERE (`e`.`entity_id` IN (28,16,25,17)) AND (`t_d`.`attribute_id` IN (45)) AND (`t_d`.`store_id` = IFNULL(`t_s`.`store_id`,0))",
			"SELECT `t_d`.`attribute_id`, `e`.`entity_id`, `t_d`.`value` AS `default_value`, IF((t_s.value_id IS NULL), t_d.value, t_s.value) AS `value`, IF((2.718281 IS NULL), t_d.value, t_s.value) AS `value` FROM `catalog_category_entity` AS `e` INNER JOIN `catalog_category_entity_varchar` AS `t_d` ON (`e`.`entity_id` = `t_d`.`entity_id`) LEFT JOIN `catalog_category_entity_varchar` AS `t_s` ON (`t_s`.`attribute_id` >= `t_d`.`attribute_id`) WHERE (`e`.`entity_id` IN (28,16,25,17)) AND (`t_d`.`attribute_id` IN (45)) AND (`t_d`.`store_id` = IFNULL(`t_s`.`store_id`,0))",
		)
	})

	t.Run("simple", func(t *testing.T) {
		sel := NewSelect("t_d.attribute_id", "e.entity_id").
			FromAlias("catalog_category_entity", "e").
			Where(
				Column("e.entity_id").In().Int64s(28, 16, 25, 17),
				Column("t_d.attribute_id").In().Int64s(45),
			)

		compareToSQL(t, sel, errors.NoKind,
			"SELECT `t_d`.`attribute_id`, `e`.`entity_id` FROM `catalog_category_entity` AS `e` WHERE (`e`.`entity_id` IN (28,16,25,17)) AND (`t_d`.`attribute_id` IN (45))",
			"SELECT `t_d`.`attribute_id`, `e`.`entity_id` FROM `catalog_category_entity` AS `e` WHERE (`e`.`entity_id` IN (28,16,25,17)) AND (`t_d`.`attribute_id` IN (45))",
		)
	})
}

func TestConditionExpr(t *testing.T) {
	t.Parallel()
	t.Run("quoted string", func(t *testing.T) {
		s := NewSelect().AddColumns("month", "total").AddColumnsConditions(Expr(`"best"`)).From("sales_by_month")
		compareToSQL(t, s, errors.NoKind,
			"SELECT `month`, `total`, \"best\" FROM `sales_by_month`",
			"",
		)
	})
}

func TestConditionSplitColumn(t *testing.T) {
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

func (ai appendInt) MapColumns(cm *ColumnMap) error {
	i := int(ai)
	return cm.Int(&i).Err()
}

func TestConditionAppendArgs(t *testing.T) {
	t.Parallel()
	t.Run("PH,val,expr,PH", func(t *testing.T) {
		s := NewSelect("sku").FromAlias("catalog", "e").
			// alias t_d ignored and not needed in this test case
			Where(
				Column("e.entity_id").PlaceHolder(),                           // 678
				Column("t_d.attribute_id").In().Int64s(45),                    // 45
				Column("t_d.store_id").Equal().SQLIfNull("t_s.store_id", "0"), // Does not make sense this WHERE condition ;-)
				Column("t_d.store_id").Equal().PlaceHolder(),                  // 17
			)

		compareToSQL(t, s.WithDBR().TestWithArgs(Qualify("e", appendInt(678)), Qualify("t_d", appendInt(17))),
			errors.NoKind,
			"SELECT `sku` FROM `catalog` AS `e` WHERE (`e`.`entity_id` = ?) AND (`t_d`.`attribute_id` IN (45)) AND (`t_d`.`store_id` = IFNULL(`t_s`.`store_id`,0)) AND (`t_d`.`store_id` = ?)",
			"SELECT `sku` FROM `catalog` AS `e` WHERE (`e`.`entity_id` = 678) AND (`t_d`.`attribute_id` IN (45)) AND (`t_d`.`store_id` = IFNULL(`t_s`.`store_id`,0)) AND (`t_d`.`store_id` = 17)",
			int64(678), int64(17),
		)
	})

	t.Run("PH,val,PH", func(t *testing.T) {
		s := NewSelect("sku").FromAlias("catalog", "e").
			// alias t_d ignored and not needed in this test case
			Where(
				Column("e.entity_id").PlaceHolder(),          // 678
				Column("t_d.attribute_id").In().Int64s(45),   // 45
				Column("t_d.store_id").Equal().PlaceHolder(), // 17
			)

		compareToSQL(t, s.WithDBR().TestWithArgs(Qualify("e", appendInt(678)), Qualify("t_d", appendInt(17))),
			errors.NoKind,
			"SELECT `sku` FROM `catalog` AS `e` WHERE (`e`.`entity_id` = ?) AND (`t_d`.`attribute_id` IN (45)) AND (`t_d`.`store_id` = ?)",
			"SELECT `sku` FROM `catalog` AS `e` WHERE (`e`.`entity_id` = 678) AND (`t_d`.`attribute_id` IN (45)) AND (`t_d`.`store_id` = 17)",
			int64(678), int64(17),
		)
	})
}

func TestCondition_Sub(t *testing.T) {
	t.Parallel()

	countSel := NewSelect().AddColumnsConditions(
		Expr("((? / COUNT(*)) * 10)"),
	).From("dml_fake_person")

	idSel := NewSelect("id", "first_name", "last_name").From("dml_fake_person").Where(
		Expr("RAND()").Less().Sub(countSel),
	).Limit(0, 40)
	idSel.IsOrderByRand = true

	compareToSQL2(t, idSel, errors.NoKind,
		"SELECT `id`, `first_name`, `last_name` FROM `dml_fake_person` WHERE (RAND() < (SELECT ((? / COUNT(*)) * 10) FROM `dml_fake_person`)) ORDER BY RAND() LIMIT 0,40",
	)
}

func TestConditions_Clone(t *testing.T) {
	t.Parallel()

	t.Run("non-nil", func(t *testing.T) {
		cnd := Conditions{
			Column("a").Equal().Float64(3.141),
			Column("b").Sub(NewSelect("x", "y").From("z")),
			Columns("g", "h", "i"),
			nil,
		}
		cnd2 := cnd.Clone()

		notEqualPointers(t, cnd, cnd2)
		assert.NotEqual(t, fmt.Sprintf("%#v", cnd), fmt.Sprintf("%#v", cnd2))
		// A weird case is that reflect.DeepEqual does execute ToSQL and hence fills
		// the internal cache field `cachedSQL`. Panicing before writing to field
		// `cachedSQL` does not work because of an internal recover in
		// reflect.DeepEqual.
		notEqualPointers(t, cnd[1].Right.Sub, cnd2[1].Right.Sub)

		assert.Exactly(t, cnd[2].Columns, cnd2[2].Columns)
		cnd2[2].Columns = cnd2[2].Columns[:0]
		cnd2[2].Columns = append(cnd2[2].Columns, "j", "k")

		assert.Exactly(t, []string{"g", "h", "i"}, cnd[2].Columns)
		assert.Exactly(t, []string{"j", "k"}, cnd2[2].Columns)
	})

	t.Run("nil", func(t *testing.T) {
		var cnd Conditions
		cnd2 := cnd.Clone()
		assert.Nil(t, cnd2)
	})
}

func TestConditionJoins_Clone(t *testing.T) {
	t.Parallel()

	jn := Joins{
		&join{
			JoinType: "LEFT",
			Table:    MakeIdentifier("tableX"),
			On: Conditions{
				Column("a").Equal().Float64(3.141),
				Columns("g", "h", "i"),
			},
		},
		nil,
	}
	jn2 := jn.Clone()
	assert.Exactly(t, jn, jn2)

	notEqualPointers(t, jn[0], jn2[0])
	notEqualPointers(t, jn[0].On[0], jn2[0].On[0])
	notEqualPointers(t, jn[0].On[1].Columns, jn2[0].On[1].Columns)
}

func TestConditions_Reset(t *testing.T) {
	t.Parallel()

	s := NewSelect("sku").FromAlias("catalog", "e").
		Where(
			Column("e.entity_id").PlaceHolder(),
			Column("t_d.attribute_id").In().Int64s(45),
		)
	sql, _, err := s.ToSQL()
	assert.NoError(t, err)
	assert.Exactly(t, "SELECT `sku` FROM `catalog` AS `e` WHERE (`e`.`entity_id` = ?) AND (`t_d`.`attribute_id` IN (45))", sql)

	s.WithCacheKey("resetted_where").Wheres.Reset()
	s.Where(
		Column("e.entity_id2").PlaceHolder(),
	)

	sql, _, err = s.ToSQL()
	assert.NoError(t, err)
	assert.Exactly(t, "SELECT `sku` FROM `catalog` AS `e` WHERE (`e`.`entity_id2` = ?)", sql)

	assert.Exactly(t, []string{
		"", "SELECT `sku` FROM `catalog` AS `e` WHERE (`e`.`entity_id` = ?) AND (`t_d`.`attribute_id` IN (45))",
		"resetted_where", "SELECT `sku` FROM `catalog` AS `e` WHERE (`e`.`entity_id2` = ?)",
	},
		s.CachedQueries(), "CachedQueries")
}
