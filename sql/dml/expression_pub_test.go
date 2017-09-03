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

package dml_test

import (
	"testing"

	"github.com/corestoreio/csfw/sql/dml"
	"github.com/corestoreio/errors"
	"github.com/stretchr/testify/assert"
)

func TestTableColumnQuote(t *testing.T) {
	t.Parallel()
	tests := []struct {
		haveT string
		haveC []string
		want  []string
	}{
		{
			"t1",
			[]string{"col1", "col2"},
			[]string{"`t1`.`col1`", "`t1`.`col2`"},
		},
		{
			"t2",
			[]string{"col1", "col2", "`t2`.`col3`"},
			[]string{"`t2`.`col1`", "`t2`.`col2`", "`t2`.`col3`"},
		},
		{
			"t2a",
			[]string{"col1", "col2", "t2.col3"},
			[]string{"`t2a`.`col1`", "`t2a`.`col2`", "`t2`.`col3`"},
		},
		{
			"t3",
			[]string{"col1", "col2", "`col3`"},
			[]string{"`t3`.`col1`", "`t3`.`col2`", "`col3`"},
		},
	}

	for i, test := range tests {
		actC := dml.Quoter.ColumnsWithQualifier(test.haveT, test.haveC...)
		assert.Equal(t, test.want, actC, "Index %d", i)
	}
}

func TestSQLIfNull(t *testing.T) {
	t.Parallel()
	runner := func(want string, have ...string) func(*testing.T) {
		return func(t *testing.T) {
			var alias string
			if lh := len(have); lh%2 == 1 && lh > 1 {
				alias = have[lh-1]
				have = have[:lh-1]
			}
			ifn := dml.SQLIfNull(have...)
			if alias != "" {
				ifn = ifn.Alias(alias)
			}
			assert.Equal(t, want, ifn.Left)
			assert.True(t, ifn.IsLeftExpression, "IsLeftExpression should be true")
		}
	}
	t.Run("1 args expression", runner(
		"IFNULL(1/0,NULL)",
		"1/0",
	))
	t.Run("1 args no qualifier", runner(
		"IFNULL(`c1`,NULL)",
		"c1",
	))
	t.Run("1 args with qualifier", runner(
		"IFNULL(`t1`.`c1`,NULL)",
		"t1.c1",
	))

	t.Run("2 args expression left", runner(
		"IFNULL(1/0,`c2`)",
		"1/0", "c2",
	))
	t.Run("2 args expression right", runner(
		"IFNULL(`c2`,1/0)",
		"c2", "1/0",
	))
	t.Run("2 args no qualifier", runner(
		"IFNULL(`c1`,`c2`)",
		"c1", "c2",
	))
	t.Run("2 args with qualifier", runner(
		"IFNULL(`t1`.`c1`,`t2`.`c2`)",
		"t1.c1", "t2.c2",
	))

	t.Run("3 args expression left", runner(
		"IFNULL(1/0,`c2`)",
		"1/0", "c2",
	))
	t.Run("3 args expression right", runner(
		"IFNULL(`c2`,1/0)",
		"c2", "1/0",
	))
	t.Run("3 args no qualifier", runner(
		"IFNULL(`c1`,`c2`)",
		"c1", "c2",
	))
	t.Run("3 args with qualifier", runner(
		"IFNULL(`t1`.`c1`,`t2`.`c2`)",
		"t1.c1", "t2.c2",
	))

	t.Run("4 args", runner(
		"IFNULL(`t1`.`c1`,`t2`.`c2`)",
		"t1", "c1", "t2", "c2",
	))
	t.Run("5 args", runner(
		"IFNULL(`t1`.`c1`,`t2`.`c2`)",
		"t1", "c1", "t2", "c2",
	))

	// its own tests
	t.Run("6 args", func(t *testing.T) {
		defer func() {
			if r := recover(); r != nil {
				if err, ok := r.(error); ok {
					assert.True(t, errors.IsNotValid(err), "%+v", err)
				} else {
					t.Errorf("Panic should contain an error but got:\n%+v", r)
				}
			} else {
				t.Error("Expecting a panic but got nothing")
			}
		}()

		runner(
			"IFNULL(`t1`.`c1`,`t2`.`c2`) AS `alias_x`",
			"t1", "c1", "t2", "c2", "alias", "x",
		)(t)
	})
}

func TestSQLIf_Expression(t *testing.T) {
	t.Parallel()

	t.Run("single call", func(t *testing.T) {
		assert.Exactly(t, "IF((c.value_id > 0), c.value, d.value)", dml.SQLIf("c.value_id > 0", "c.value", "d.value").Left)
	})

	t.Run("just EXPRESSION", func(t *testing.T) {
		s := dml.NewSelect().AddColumns("a", "b", "c").
			From("table1").Where(
			dml.Expr(
				"IF((a > 0), b, c) > ?",
			).Greater().Int(4711))

		compareToSQL(t, s, nil,
			"SELECT `a`, `b`, `c` FROM `table1` WHERE (IF((a > 0), b, c) > ?)",
			"SELECT `a`, `b`, `c` FROM `table1` WHERE (IF((a > 0), b, c) > 4711)",
			int64(4711),
		)
	})

	t.Run("IF in EXPRESSION", func(t *testing.T) {
		s := dml.NewSelect().AddColumns("a", "b", "c").
			From("table1").Where(
			dml.SQLIf("a > 0", "b", "c").Greater().Int(4711))

		compareToSQL(t, s, nil,
			"SELECT `a`, `b`, `c` FROM `table1` WHERE (IF((a > 0), b, c) > ?)",
			"SELECT `a`, `b`, `c` FROM `table1` WHERE (IF((a > 0), b, c) > 4711)",
			int64(4711),
		)
	})
}

func TestSQLCase(t *testing.T) {
	t.Parallel()

	t.Run("UPDATE in columns with args", func(t *testing.T) {
		/*
					UPDATE `cataloginventory_stock_item`
					SET    `qty` = CASE product_id
								 WHEN 23434 THEN qty + 3
								 WHEN 23435 THEN qty + 2
								 WHEN 23436 THEN qty + 4
							 ELSE qty
						   end
					WHERE  ( product_id IN ( 23434, 23435, 23436 ) )
			       AND ( website_id = 4 )
		*/

		u := dml.NewUpdate("cataloginventory_stock_item").
			Set(dml.Column("qty").SQLCase("`product_id`", "qty",
				"3456", "qty+?",
				"3457", "qty+?",
				"3458", "qty+?",
			).Ints(3, 4, 5)).
			Where(
				dml.Column("product_id").In().Int64s(345, 567, 897),
				dml.Column("website_id").Int64(6),
			)

		sqlStr, args, err := u.ToSQL()
		if err != nil {
			t.Fatalf("%+v", err)
		}
		assert.Exactly(t, []interface{}{int64(3), int64(4), int64(5), int64(345), int64(567), int64(897), int64(6)}, args)
		assert.Exactly(t, "UPDATE `cataloginventory_stock_item` SET `qty`=CASE `product_id` WHEN 3456 THEN qty+? WHEN 3457 THEN qty+? WHEN 3458 THEN qty+? ELSE qty END WHERE (`product_id` IN (?,?,?)) AND (`website_id` = ?)", sqlStr)

		sqlStr = dml.Interpolate(sqlStr).ArgUnions(iFaceToArgs(args...)).String()
		assert.Exactly(t, "UPDATE `cataloginventory_stock_item` SET `qty`=CASE `product_id` WHEN 3456 THEN qty+3 WHEN 3457 THEN qty+4 WHEN 3458 THEN qty+5 ELSE qty END WHERE (`product_id` IN (345,567,897)) AND (`website_id` = 6)", sqlStr)
	})

	t.Run("cases", func(t *testing.T) {
		assert.Exactly(t,
			"CASE `product_id` WHEN 3456 THEN qty+1 WHEN 3457 THEN qty+4 WHEN 3458 THEN qty-3 ELSE qty END",
			dml.SQLCase("`product_id`", "qty",
				"3456", "qty+1",
				"3457", "qty+4",
				"3458", "qty-3",
			).Left,
		)
		assert.Exactly(t,
			"(CASE `product_id` WHEN 3456 THEN qty WHEN 3457 THEN qty ELSE qty END) AS `product_qty`",
			dml.SQLCase("`product_id`", "qty",
				"3456", "qty",
				"3457", "qty",
				"product_qty",
			).Left,
		)
		assert.Exactly(t,
			"CASE `product_id` WHEN 3456 THEN qty+1 WHEN 3457 THEN qty+4 WHEN 3458 THEN qty-3 END",
			dml.SQLCase("`product_id`", "",
				"3456", "qty+1",
				"3457", "qty+4",
				"3458", "qty-3",
			).Left,
		)
		ce := dml.SQLCase("", "", "1=1", "2", "3=2", "4")
		assert.Exactly(t, "CASE  WHEN 1=1 THEN 2 WHEN 3=2 THEN 4 END", ce.Left)
		assert.True(t, ce.IsLeftExpression, "IsLeftExpression should be true")
	})
	t.Run("case panics because of invalid length of comparions", func(t *testing.T) {
		defer func() {
			if r := recover(); r != nil {
				if err, ok := r.(error); ok {
					assert.True(t, errors.IsFatal(err))
				} else {
					t.Errorf("Panic should contain an error but got:\n%+v", r)
				}
			} else {
				t.Error("Expecting a panic but got nothing")
			}
		}()
		dml.SQLCase("", "", "1=1")
	})
}
