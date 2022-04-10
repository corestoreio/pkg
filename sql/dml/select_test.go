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
	"context"
	"database/sql"
	"fmt"
	"testing"

	"github.com/corestoreio/pkg/storage/null"
	"github.com/corestoreio/pkg/util/assert"
)

func TestSelect_BasicToSQL(t *testing.T) {
	t.Run("no table no args", func(t *testing.T) {
		sel := NewSelect().AddColumnsConditions(Expr("1").Alias("n")).AddColumnsAliases("abc", "str")
		compareToSQL2(t, sel, false,
			"SELECT 1 AS `n`, `abc` AS `str`",
		)
	})
	t.Run("no table with args", func(t *testing.T) {
		sel := NewSelect().
			AddColumnsConditions(
				Expr("?").Alias("n").Int64(1),
				Expr("CAST(? AS CHAR(20))").Alias("str").Str("a'bc"),
			)
		compareToSQL2(t, sel, false,
			"SELECT 1 AS `n`, CAST('a\\'bc' AS CHAR(20)) AS `str`",
		)
	})
	t.Run("no table with placeholders Args as Records", func(t *testing.T) {
		p := &dmlPerson{
			Name: "a'bc",
		}
		sel := NewSelect().
			AddColumnsConditions(
				Expr("?").Alias("n").Int64(1),
				Expr("CAST(:name AS CHAR(20))").Alias("str"),
			).WithDBR(dbMock{}).TestWithArgs(Qualify("", p))

		compareToSQL(t, sel, false,
			"SELECT 1 AS `n`, CAST(? AS CHAR(20)) AS `str`",
			"SELECT 1 AS `n`, CAST('a\\'bc' AS CHAR(20)) AS `str`",
			"a'bc",
		)
	})

	t.Run("two cols, one table, one condition", func(t *testing.T) {
		sel := NewSelect("a", "b").From("c").Where(Column("id").Equal().Int(1))
		compareToSQL2(t, sel, false,
			"SELECT `a`, `b` FROM `c` WHERE (`id` = 1)",
		)
	})

	t.Run("place holders", func(t *testing.T) {
		sel := NewSelect("a", "b").From("c").Where(
			Column("id").Greater().PlaceHolder(),
			Column("email").Like().NamedArg("ema1l"),
		)
		compareToSQL2(t, sel, false,
			"SELECT `a`, `b` FROM `c` WHERE (`id` > ?) AND (`email` LIKE ?)",
		)
		assert.Exactly(t, []string{"id", ":ema1l"}, sel.qualifiedColumns)
	})

	t.Run("column right expression without arguments", func(t *testing.T) {
		compareToSQL2(t,
			NewSelect("sku", "name").From("products").Where(
				Column("id").NotBetween().Ints(4, 7),
				Column("name").NotEqual().Expr("CONCAT('Canon','E0S 5D Mark III')"),
			),
			false,
			"SELECT `sku`, `name` FROM `products` WHERE (`id` NOT BETWEEN 4 AND 7) AND (`name` != CONCAT('Canon','E0S 5D Mark III'))",
		)
	})

	t.Run("column right expression with one argument", func(t *testing.T) {
		compareToSQL2(t,
			NewSelect("sku", "name").From("products").Where(
				Column("id").NotBetween().Ints(4, 7),
				Column("name").Like().Expr("CONCAT('Canon',?,'E0S 7D Mark VI')").Str("Camera"),
			),
			false,
			"SELECT `sku`, `name` FROM `products` WHERE (`id` NOT BETWEEN 4 AND 7) AND (`name` LIKE CONCAT('Canon','Camera','E0S 7D Mark VI'))",
		)
	})

	t.Run("column right expression with slice argument (wrong SQL code)", func(t *testing.T) {
		sel := NewSelect("sku", "name").From("products").Where(
			Column("id").NotBetween().Ints(4, 7),
			Column("name").NotLike().Expr("CONCAT('Canon',?,'E0S 8D Mark VII')").Strs("Camera", "Photo", "Flash"),
		)
		compareToSQL2(t, sel, false,
			"SELECT `sku`, `name` FROM `products` WHERE (`id` NOT BETWEEN 4 AND 7) AND (`name` NOT LIKE CONCAT('Canon',('Camera','Photo','Flash'),'E0S 8D Mark VII'))",
		)
	})
	t.Run("column right expression with slice argument (correct SQL code)", func(t *testing.T) {
		sel := NewSelect("sku", "name").From("products").Where(
			Column("id").NotBetween().Ints(4, 7),
			Column("name").NotLike().Expr("CONCAT('Canon',?,?,?,'E0S 8D Mark VII')").Str("Camera").Str("Photo").Str("Flash"),
		)
		compareToSQL2(t, sel, false,
			"SELECT `sku`, `name` FROM `products` WHERE (`id` NOT BETWEEN 4 AND 7) AND (`name` NOT LIKE CONCAT('Canon','Camera','Photo','Flash','E0S 8D Mark VII'))",
		)
	})

	t.Run("column left and right expression without arguments", func(t *testing.T) {
		sel := NewSelect("sku", "name").From("products").Where(
			Column("id").NotBetween().Ints(4, 7),
			Column("name").NotEqual().Expr("CONCAT('Canon','E0S 5D Mark III')"),
		)
		compareToSQL2(t, sel, false,
			"SELECT `sku`, `name` FROM `products` WHERE (`id` NOT BETWEEN 4 AND 7) AND (`name` != CONCAT('Canon','E0S 5D Mark III'))",
		)
	})

	t.Run("IN with expand", func(t *testing.T) {
		sel := NewSelect("sku", "name").From("products").Where(
			Column("id").In().PlaceHolder(),
			Column("name").NotIn().PlaceHolder(),
		)

		selA := sel.WithDBR(dbMock{}).ExpandPlaceHolders()

		compareToSQL(t, selA.TestWithArgs([]int{3, 4, 5}, []null.String{null.MakeString("A1"), {}, null.MakeString("A2")}), false,
			"SELECT `sku`, `name` FROM `products` WHERE (`id` IN (?,?,?)) AND (`name` NOT IN (?,?,?))",
			"",
			int64(3), int64(4), int64(5), "A1", nil, "A2",
		)

		compareToSQL(t, selA.TestWithArgs([]int{3, 4, 5, 6, 7}, []null.String{{}, null.MakeString("A2")}), false,
			"SELECT `sku`, `name` FROM `products` WHERE (`id` IN (?,?,?,?,?)) AND (`name` NOT IN (?,?))",
			"",
			int64(3), int64(4), int64(5), int64(6), int64(7), nil, "A2",
		)
	})

	t.Run("IN with PlaceHolders", func(t *testing.T) {
		sel := NewSelect("email").From("tableX").Where(Column("id").In().PlaceHolders(2))
		compareToSQL2(t, sel, false,
			"SELECT `email` FROM `tableX` WHERE (`id` IN (?,?))",
		)
		sel = NewSelect("email").From("tableX").Where(Column("id").In().PlaceHolders(1))
		compareToSQL2(t, sel, false,
			"SELECT `email` FROM `tableX` WHERE (`id` IN (?))",
		)
		sel = NewSelect("email").From("tableX").Where(Column("id").In().PlaceHolders(0))
		compareToSQL2(t, sel, false,
			"SELECT `email` FROM `tableX` WHERE (`id` IN ())",
		)
		sel = NewSelect("email").From("tableX").Where(Column("id").In().PlaceHolders(-10))
		compareToSQL2(t, sel, false,
			"SELECT `email` FROM `tableX` WHERE (`id` IN ())",
		)
	})
}

func TestSelect_FullToSQL(t *testing.T) {
	sel := NewSelect("a", "b").
		Distinct().
		FromAlias("c", "cc").
		Where(
			ParenthesisOpen(),
			Column("d").Int(1),
			Column("e").Str("wat").Or(),
			ParenthesisClose(),
			Column("f").Int(2),
			Column("g").Int(3),
			Column("h").In().Int64s(4, 5, 6),
		).
		GroupBy("ab").
		Having(
			ParenthesisOpen(),
			Column("m").Int(33),
			Column("n").Str("wh3r3").Or(),
			ParenthesisClose(),
			Expr("j = k"),
		).
		OrderBy("l").
		Limit(8, 7)

	compareToSQL2(t, sel, false,
		"SELECT DISTINCT `a`, `b` FROM `c` AS `cc` WHERE ((`d` = 1) OR (`e` = 'wat')) AND (`f` = 2) AND (`g` = 3) AND (`h` IN (4,5,6)) GROUP BY `ab` HAVING ((`m` = 33) OR (`n` = 'wh3r3')) AND (j = k) ORDER BY `l` LIMIT 8,7",
	)
}

func TestSelect_ComplexExpr(t *testing.T) {
	t.Run("two args in one condition", func(t *testing.T) {
		sel := NewSelect("a", "b", "z", "y", "x").From("c").
			Distinct().
			Where(Expr("`d` = ? OR `e` = ?").Int64(1).Str("wat")).
			Where(
				Column("g").Int(3),
				Column("h").In().Int64s(1, 2, 3),
			).
			GroupBy("ab").GroupBy("ii").GroupBy("iii").
			Having(Expr("j = k"), Column("jj").Int64(1)).
			Having(Column("jjj").Int64(2)).
			OrderBy("l1").OrderBy("l2").OrderBy("l3").
			Limit(8, 7)

		compareToSQL2(t, sel, false,
			"SELECT DISTINCT `a`, `b`, `z`, `y`, `x` FROM `c` WHERE (`d` = 1 OR `e` = 'wat') AND (`g` = 3) AND (`h` IN (1,2,3)) GROUP BY `ab`, `ii`, `iii` HAVING (j = k) AND (`jj` = 1) AND (`jjj` = 2) ORDER BY `l1`, `l2`, `l3` LIMIT 8,7",
			// int64(1), "wat", int64(3), int64(1), int64(2), int64(3), int64(1), int64(2),
		)
	})
}

func TestSelect_OrderByRandom_Strings(t *testing.T) {
	t.Run("simple select", func(t *testing.T) {
		compareToSQL2(t,
			NewSelect("id", "first_name", "last_name").
				From("dml_fake_person").
				OrderByRandom("id", 25),
			false,
			"SELECT `id`, `first_name`, `last_name` FROM `dml_fake_person`  JOIN (SELECT `id` FROM `dml_fake_person` WHERE (RAND() < (SELECT ((25 / COUNT(*)) * 10) FROM `dml_fake_person`)) ORDER BY RAND() LIMIT 0,25) AS `randdml_fake_person` USING (`id`)",
		)
	})

	t.Run("WHERE condition", func(t *testing.T) {
		compareToSQL2(t,
			NewSelect("id", "first_name", "last_name").
				From("dml_fake_person").
				Where(Column("id").LessOrEqual().Int(30)).
				OrderByRandom("id", 25),
			false,
			"SELECT `id`, `first_name`, `last_name` FROM `dml_fake_person`  JOIN (SELECT `id` FROM `dml_fake_person` WHERE (RAND() < (SELECT ((25 / COUNT(*)) * 10) FROM `dml_fake_person` WHERE (`id` <= 30))) AND (`id` <= 30) ORDER BY RAND() LIMIT 0,25) AS `randdml_fake_person` USING (`id`) WHERE (`id` <= 30)",
		)
	})

	t.Run("one join", func(t *testing.T) {
		sqlObj := NewSelect("p1.*", "p2.*").FromAlias("dml_people", "p1").
			Distinct().StraightJoin().SQLNoCache().
			Join(
				MakeIdentifier("dml_people").Alias("p2"),
				Expr("`p2`.`id` = `p1`.`id`"),
				Column("p1.id").Int(142),
			).OrderByRandom("id", 100)

		compareToSQL2(t, sqlObj, false,
			"SELECT DISTINCT STRAIGHT_JOIN SQL_NO_CACHE `p1`.*, `p2`.* FROM `dml_people` AS `p1` INNER JOIN `dml_people` AS `p2` ON (`p2`.`id` = `p1`.`id`) AND (`p1`.`id` = 142)  JOIN (SELECT `id` FROM `dml_people` WHERE (RAND() < (SELECT ((100 / COUNT(*)) * 10) FROM `dml_people`)) ORDER BY RAND() LIMIT 0,100) AS `randdml_people` USING (`id`)",
		)
	})
}

func TestSelect_OrderByRandom_Integration(t *testing.T) {
	s := createRealSessionWithFixtures(t, nil)
	defer testCloser(t, s)

	t.Run("Load IDs", func(t *testing.T) {
		ids, err := s.WithQueryBuilder(
			NewSelect("id").From("dml_people").OrderByRandom("id", 10),
		).LoadUint64s(context.TODO(), nil)
		assert.NoError(t, err)
		assert.Len(t, ids, 2)
	})

	// TODO test this when in a column is one row NULL ... then it should fail scanning
}

func TestSelect_Paginate(t *testing.T) {
	t.Run("asc", func(t *testing.T) {
		compareToSQL2(t,
			NewSelect("a", "b").
				From("c").
				Where(Column("d").Int(1)).
				Paginate(3, 30).
				OrderBy("id"),
			false,
			"SELECT `a`, `b` FROM `c` WHERE (`d` = 1) ORDER BY `id` LIMIT 60,30",
		)
	})
	t.Run("desc", func(t *testing.T) {
		compareToSQL2(t,
			NewSelect("a", "b").
				From("c").
				Where(Column("d").Int(1)).
				Paginate(1, 20).
				OrderByDesc("id"),
			false,
			"SELECT `a`, `b` FROM `c` WHERE (`d` = 1) ORDER BY `id` DESC LIMIT 0,20",
		)
	})
}

func TestSelect_WithoutWhere(t *testing.T) {
	compareToSQL2(t,
		NewSelect("a", "b").From("c"),
		false,
		"SELECT `a`, `b` FROM `c`",
	)
}

func TestSelect_MultiHavingSQL(t *testing.T) {
	compareToSQL2(t,
		NewSelect("a", "b").From("c").
			Where(Column("p").Int(1)).
			GroupBy("z").Having(Column("z`z").Int(2), Column("y").Int(3)),
		false,
		"SELECT `a`, `b` FROM `c` WHERE (`p` = 1) GROUP BY `z` HAVING (`zz` = 2) AND (`y` = 3)",
	)
}

func TestSelect_MultiOrderSQL(t *testing.T) {
	compareToSQL2(t,
		NewSelect("a", "b").From("c").OrderBy("name").OrderByDesc("id"),
		false,
		"SELECT `a`, `b` FROM `c` ORDER BY `name`, `id` DESC",
	)
}

func TestSelect_OrderByDeactivated(t *testing.T) {
	compareToSQL2(t,
		NewSelect("a", "b").From("c").OrderBy("name").OrderByDeactivated(),
		false,
		"SELECT `a`, `b` FROM `c` ORDER BY NULL",
	)
}

func TestSelect_ConditionColumn(t *testing.T) {
	// TODO rewrite test to use every type which implements interface Argument and every operator

	runner := func(wf *Condition, wantSQL string) func(*testing.T) {
		return func(t *testing.T) {
			compareToSQL2(t,
				NewSelect("a", "b").From("c").Where(wf),
				false,
				wantSQL,
			)
		}
	}
	t.Run("single int64", runner(
		Column("d").Int64(33),
		"SELECT `a`, `b` FROM `c` WHERE (`d` = 33)",
	))
	t.Run("IN int64", runner(
		Column("d").In().Int64s(33, 44),
		"SELECT `a`, `b` FROM `c` WHERE (`d` IN (33,44))",
	))
	t.Run("single float64", runner(
		Column("d").Float64(33.33),
		"SELECT `a`, `b` FROM `c` WHERE (`d` = 33.33)",
	))
	t.Run("IN float64", runner(
		Column("d").In().Float64s(33.44, 44.33),
		"SELECT `a`, `b` FROM `c` WHERE (`d` IN (33.44,44.33))",
	))
	t.Run("NOT IN float64", runner(
		Column("d").NotIn().Float64s(33.1, 44.2),
		"SELECT `a`, `b` FROM `c` WHERE (`d` NOT IN (33.1,44.2))",
	))
	t.Run("single int", runner(
		Column("d").Equal().Int(33),
		"SELECT `a`, `b` FROM `c` WHERE (`d` = 33)",
	))
	t.Run("IN int", runner(
		Column("d").In().Ints(33, 44),
		"SELECT `a`, `b` FROM `c` WHERE (`d` IN (33,44))",
	))
	t.Run("single string", runner(
		Column("d").Str("w"),
		"SELECT `a`, `b` FROM `c` WHERE (`d` = 'w')",
	))
	t.Run("IN string", runner(
		Column("d").In().Strs("x", "y"),
		"SELECT `a`, `b` FROM `c` WHERE (`d` IN ('x','y'))",
	))

	t.Run("BETWEEN int64", runner(
		Column("d").Between().Int64s(5, 6),
		"SELECT `a`, `b` FROM `c` WHERE (`d` BETWEEN 5 AND 6)",
	))
	t.Run("NOT BETWEEN int64", runner(
		Column("d").NotBetween().Int64s(5, 6),
		"SELECT `a`, `b` FROM `c` WHERE (`d` NOT BETWEEN 5 AND 6)",
	))

	t.Run("LIKE string", runner(
		Column("d").Like().Str("x%"),
		"SELECT `a`, `b` FROM `c` WHERE (`d` LIKE 'x%')",
	))
	t.Run("NOT LIKE string", runner(
		Column("d").NotLike().Str("x%"),
		"SELECT `a`, `b` FROM `c` WHERE (`d` NOT LIKE 'x%')",
	))

	t.Run("Less float64", runner(
		Column("d").Less().Float64(5.1),
		"SELECT `a`, `b` FROM `c` WHERE (`d` < 5.1)",
	))
	t.Run("Greater float64", runner(
		Column("d").Greater().Float64(5.1),
		"SELECT `a`, `b` FROM `c` WHERE (`d` > 5.1)",
	))
	t.Run("LessOrEqual float64", runner(
		Column("d").LessOrEqual().Float64(5.1),
		"SELECT `a`, `b` FROM `c` WHERE (`d` <= 5.1)",
	))
	t.Run("GreaterOrEqual float64", runner(
		Column("d").GreaterOrEqual().Float64(5.1),
		"SELECT `a`, `b` FROM `c` WHERE (`d` >= 5.1)",
	))
}

func TestSelect_Null(t *testing.T) {
	t.Run("col is null", func(t *testing.T) {
		compareToSQL2(t,
			NewSelect("a", "b").From("c").Where(Column("r").Null()),
			false,
			"SELECT `a`, `b` FROM `c` WHERE (`r` IS NULL)",
		)
	})

	t.Run("col is not null", func(t *testing.T) {
		compareToSQL2(t,
			NewSelect("a", "b").From("c").Where(Column("r").NotNull()),
			false,
			"SELECT `a`, `b` FROM `c` WHERE (`r` IS NOT NULL)",
		)
	})

	t.Run("complex", func(t *testing.T) {
		compareToSQL2(t,
			NewSelect("a", "b").From("c").
				Where(
					Column("r").Null(),
					Column("d").Int(3),
					Column("ab").Null(),
					Column("w").NotNull(),
				),
			false,
			"SELECT `a`, `b` FROM `c` WHERE (`r` IS NULL) AND (`d` = 3) AND (`ab` IS NULL) AND (`w` IS NOT NULL)",
		)
	})
}

func TestSelect_WhereNULL(t *testing.T) {
	t.Run("one nil", func(t *testing.T) {
		compareToSQL2(t,
			NewSelect("a").From("b").Where(Column("a")),
			false,
			"SELECT `a` FROM `b` WHERE (`a` IS NULL)",
		)
	})

	t.Run("no values", func(t *testing.T) {
		compareToSQL2(t,
			NewSelect("a").From("b").Where(Column("a").PlaceHolder()),
			false,
			"SELECT `a` FROM `b` WHERE (`a` = ?)",
		)
	})

	t.Run("empty Ints trigger invalid SQL", func(t *testing.T) {
		var iVal []int
		compareToSQL2(t,
			NewSelect("a").From("b").Where(Column("a").In().Ints(iVal...)),
			false,
			"SELECT `a` FROM `b` WHERE (`a` IN ())",
		)
	})

	t.Run("Map nil arg", func(t *testing.T) {
		s := NewSelect("a").From("b").
			Where(
				Column("a"),
				Column("b").Bool(false),
				Column("c").Null(),
				Column("d").NotNull(),
			)
		compareToSQL2(t, s, false,
			"SELECT `a` FROM `b` WHERE (`a` IS NULL) AND (`b` = 0) AND (`c` IS NULL) AND (`d` IS NOT NULL)",
		)
	})
}

func TestSelect_Varieties(t *testing.T) {
	// This would be incorrect SQL!
	compareToSQL2(t, NewSelect("id, name, email").From("users"), false,
		"SELECT `id, name, email` FROM `users`",
	)
	// With unsafe it still gets quoted because unsafe has been applied after
	// the column names has been added.
	compareToSQL2(t, NewSelect("id, name, email").Unsafe().From("users"), false,
		"SELECT `id, name, email` FROM `users`",
	)
	// correct way to handle it
	compareToSQL2(t, NewSelect("id", "name", "email").From("users"), false,
		"SELECT `id`, `name`, `email` FROM `users`",
	)
}

func TestSelect_Load_Slice_Scanner(t *testing.T) {
	s := createRealSessionWithFixtures(t, nil)
	defer testCloser(t, s)

	var people dmlPersons
	count, err := s.WithQueryBuilder(NewSelect("id", "name", "email").From("dml_people").OrderBy("id")).Load(context.TODO(), &people)

	assert.NoError(t, err)
	assert.Exactly(t, uint64(2), count)

	assert.Exactly(t, len(people.Data), 2)
	if len(people.Data) == 2 {
		// Make sure that the Ids are isSet. It'ab possible (maybe?) that different DBs isSet ids differently so
		// don't assume they're 1 and 2.
		assert.True(t, people.Data[0].ID > 0)
		assert.True(t, people.Data[1].ID > people.Data[0].ID)

		assert.Exactly(t, "Sir George", people.Data[0].Name)
		assert.True(t, people.Data[0].Email.Valid)
		assert.Exactly(t, "SirGeorge@GoIsland.com", people.Data[0].Email.Data)
		assert.Exactly(t, "Dmitri", people.Data[1].Name)
		assert.True(t, people.Data[1].Email.Valid)
		assert.Exactly(t, "userXYZZ@emailServerX.com", people.Data[1].Email.Data)
	}
}

func TestSelect_Load_Rows(t *testing.T) {
	s := createRealSessionWithFixtures(t, nil)
	defer testCloser(t, s)

	t.Run("found", func(t *testing.T) {
		var person dmlPerson
		_, err := s.WithQueryBuilder(NewSelect("id", "name", "email").From("dml_people").
			Where(Column("email").Str("SirGeorge@GoIsland.com"))).Load(context.TODO(), &person)
		assert.NoError(t, err)
		assert.True(t, person.ID > 0)
		assert.Exactly(t, "Sir George", person.Name)
		assert.True(t, person.Email.Valid)
		assert.Exactly(t, "SirGeorge@GoIsland.com", person.Email.Data)
	})

	t.Run("not found", func(t *testing.T) {
		var person2 dmlPerson
		count, err := s.WithQueryBuilder(NewSelect("id", "name", "email").From("dml_people").
			Where(Column("email").Str("dontexist@uservoice.com"))).Load(context.TODO(), &person2)

		assert.NoError(t, err)
		assert.Exactly(t, dmlPerson{}, person2)
		assert.Empty(t, count, "Should have no rows loaded")
	})
}

func TestSelectBySQL_Load_Slice(t *testing.T) {
	s := createRealSessionWithFixtures(t, nil)
	defer testCloser(t, s)

	t.Run("single slice item", func(t *testing.T) {
		var people dmlPersons
		count, err := s.WithQueryBuilder(QuerySQL("SELECT `name` FROM `dml_people` WHERE `email` = ?")).
			Load(context.TODO(), &people, "SirGeorge@GoIsland.com")

		assert.NoError(t, err)
		assert.Exactly(t, uint64(1), count)
		if len(people.Data) == 1 {
			assert.Exactly(t, "Sir George", people.Data[0].Name)
			assert.Exactly(t, uint64(0), people.Data[0].ID)      // not set
			assert.Exactly(t, false, people.Data[0].Email.Valid) // not set
			assert.Exactly(t, "", people.Data[0].Email.Data)     // not set
		}
	})

	t.Run("IN Clause (multiple args,interpolate)", func(t *testing.T) {
		ids, err := s.WithQueryBuilder(NewSelect("id").From("dml_people").
			Where(Column("id").In().Int64s(1, 2, 3))).LoadInt64s(context.TODO(), nil)
		assert.NoError(t, err)
		assert.Exactly(t, []int64{1, 2}, ids)
	})
	t.Run("IN Clause (single args,interpolate)", func(t *testing.T) {
		ids, err := s.WithQueryBuilder(NewSelect("id").From("dml_people").
			Where(Column("id").In().Int64s(2))).LoadInt64s(context.TODO(), nil)
		assert.NoError(t, err)
		assert.Exactly(t, []int64{2}, ids)
	})
	t.Run("NOT IN Clause (multiple args)", func(t *testing.T) {
		ids, err := NewSelect("id").From("dml_people").
			Where(Column("id").NotIn().Int64s(2, 3)).WithDBR(s.DB).LoadInt64s(context.TODO(), nil)
		assert.NoError(t, err)
		assert.Exactly(t, []int64{1}, ids)
	})
	t.Run("Scan string into arg UINT returns error", func(t *testing.T) {
		var people dmlPersons
		rc, err := NewSelect().AddColumnsAliases("email", "id", "name", "email").From("dml_people").
			WithDBR(s.DB).Load(context.TODO(), &people)
		assert.Error(t, err)
		assert.Error(t, err)
		assert.Empty(t, rc)
	})
}

func TestSelect_LoadType_Single(t *testing.T) {
	s := createRealSessionWithFixtures(t, nil)
	defer testCloser(t, s)

	t.Run("LoadNullString", func(t *testing.T) {
		name, found, err := NewSelect("name").From("dml_people").Where(Column("email").PlaceHolder()).
			WithDBR(s.DB).LoadNullString(context.TODO(), "SirGeorge@GoIsland.com")
		assert.NoError(t, err)
		assert.True(t, found)
		assert.Exactly(t, null.MakeString("Sir George"), name)
	})
	t.Run("LoadNullString too many columns", func(t *testing.T) {
		name, found, err := NewSelect("name", "email").From("dml_people").Where(Expr("email = 'SirGeorge@GoIsland.com'")).
			WithDBR(s.DB).LoadNullString(context.TODO())
		assert.Error(t, err)
		assert.False(t, found)
		assert.Empty(t, name.Data)
	})
	t.Run("LoadNullString not found", func(t *testing.T) {
		name, found, err := NewSelect("name").From("dml_people").Where(Expr("email = 'notfound@example.com'")).
			WithDBR(s.DB).LoadNullString(context.TODO())
		assert.NoError(t, err)
		assert.False(t, found)
		assert.Exactly(t, null.String{}, name)
	})

	t.Run("LoadNullInt64", func(t *testing.T) {
		id, found, err := NewSelect("id").From("dml_people").Limit(0, 1).WithDBR(s.DB).LoadNullInt64(context.TODO())
		assert.NoError(t, err)
		assert.True(t, found)
		assert.True(t, id.Int64 > 0)
	})
	t.Run("LoadNullInt64 too many columns", func(t *testing.T) {
		id, found, err := NewSelect("id", "email").From("dml_people").Limit(0, 1).WithDBR(s.DB).LoadNullInt64(context.TODO())
		assert.Error(t, err)
		assert.False(t, found)
		assert.Exactly(t, null.Int64{}, id)
	})
	t.Run("LoadNullInt64 not found", func(t *testing.T) {
		id, found, err := NewSelect("id").From("dml_people").Where(Expr("id=236478326")).WithDBR(s.DB).LoadNullInt64(context.TODO())
		assert.NoError(t, err)
		assert.False(t, found)
		assert.Exactly(t, null.Int64{}, id)
	})

	t.Run("LoadNullUint64", func(t *testing.T) {
		id, found, err := NewSelect("id").From("dml_people").Limit(0, 1).WithDBR(s.DB).LoadNullUint64(context.TODO())
		assert.NoError(t, err)
		assert.True(t, found)
		assert.True(t, id.Uint64 > 0)
	})
	t.Run("LoadNullUint64 too many columns", func(t *testing.T) {
		id, found, err := NewSelect("id", "email").From("dml_people").Limit(0, 1).WithDBR(s.DB).LoadNullUint64(context.TODO())
		assert.Error(t, err)
		assert.False(t, found)
		assert.Exactly(t, null.Uint64{}, id)
	})
	t.Run("LoadNullUint64 not found", func(t *testing.T) {
		id, found, err := NewSelect("id").From("dml_people").Where(Expr("id=236478326")).WithDBR(s.DB).LoadNullUint64(context.TODO())
		assert.NoError(t, err)
		assert.False(t, found)
		assert.Exactly(t, null.Uint64{}, id)
	})

	t.Run("LoadNullFloat64", func(t *testing.T) {
		id, found, err := NewSelect("id").From("dml_people").Limit(0, 1).WithDBR(s.DB).LoadNullFloat64(context.TODO())
		assert.NoError(t, err)
		assert.True(t, found)
		assert.True(t, id.Float64 > 0)
	})
	t.Run("LoadNullFloat64 too many columns", func(t *testing.T) {
		id, found, err := NewSelect("id", "email").From("dml_people").Limit(0, 1).WithDBR(s.DB).LoadNullFloat64(context.TODO())
		assert.Error(t, err)
		assert.False(t, found)
		assert.Exactly(t, null.Float64{}, id)
	})
	t.Run("LoadNullFloat64 not found", func(t *testing.T) {
		id, found, err := NewSelect("id").From("dml_people").Where(Expr("id=236478326")).WithDBR(s.DB).LoadNullFloat64(context.TODO())
		assert.NoError(t, err)
		assert.False(t, found)
		assert.Exactly(t, null.Float64{}, id)
	})

	t.Run("LoadDecimal", func(t *testing.T) {
		income, found, err := NewSelect("avg_income").From("dml_people").Where(Column("email").Like().Str(`SirGeorge@GoIsland.com`)).
			WithDBR(s.DB).LoadDecimal(context.TODO())
		assert.NoError(t, err)
		assert.True(t, found)
		assert.Exactly(t, null.MakeDecimalInt64(33366677, 5), income)
	})
}

func TestSelect_WithArgs_LoadUint64(t *testing.T) {
	s := createRealSessionWithFixtures(t, &installFixturesConfig{
		AddPeopleWithMaxUint64: true,
	})
	defer testCloser(t, s)

	// Despite it seems that Go can support large uint64 values ... the down side is that
	// the byte encoded uint64 gets transferred as a string and MySQL/MariaDB must convert that
	// string into a bigint.
	const bigID uint64 = 18446744073700551613 // see also file dml_test.go MaxUint64

	sel := NewSelect("id").From("dml_people").Where(Column("id").Uint64(bigID))

	t.Run("MaxUint64 prepared stmt o:equal", func(t *testing.T) {
		id, found, err := sel.WithDBR(s.DB).LoadNullUint64(context.TODO())
		assert.NoError(t, err)
		assert.True(t, found)
		assert.Exactly(t, null.MakeUint64(bigID), id)
	})
	t.Run("MaxUint64 interpolated o:equal", func(t *testing.T) {
		id, found, err := sel.WithDBR(s.DB).Interpolate().LoadNullUint64(context.TODO())
		assert.NoError(t, err)
		assert.True(t, found)
		assert.Exactly(t, null.MakeUint64(bigID), id)
	})
}

func TestSelect_WithArgs_LoadType_Slices(t *testing.T) {
	s := createRealSessionWithFixtures(t, nil)
	defer testCloser(t, s)
	t.Run("LoadStrings", func(t *testing.T) {
		names, err := NewSelect("name").From("dml_people").WithDBR(s.DB).LoadStrings(context.TODO(), nil)
		assert.NoError(t, err)
		assert.Exactly(t, []string{"Sir George", "Dmitri"}, names)
	})
	t.Run("LoadStrings too many columns", func(t *testing.T) {
		vals, err := NewSelect("name", "email").From("dml_people").WithDBR(s.DB).LoadStrings(context.TODO(), nil)
		assert.Error(t, err)
		assert.Exactly(t, []string(nil), vals)
	})
	t.Run("LoadStrings not found", func(t *testing.T) {
		names, err := NewSelect("name").From("dml_people").Where(Expr("name ='jdhsjdf'")).WithDBR(s.DB).LoadStrings(context.TODO(), nil)
		assert.NoError(t, err)
		assert.Nil(t, names)
	})

	t.Run("LoadInt64s", func(t *testing.T) {
		names, err := NewSelect("id").From("dml_people").WithDBR(s.DB).LoadInt64s(context.TODO(), nil)
		assert.NoError(t, err)
		assert.Exactly(t, []int64{1, 2}, names)
	})
	t.Run("LoadInt64s too many columns", func(t *testing.T) {
		vals, err := NewSelect("id", "email").From("dml_people").WithDBR(s.DB).LoadInt64s(context.TODO(), nil)
		assert.Error(t, err)
		assert.Exactly(t, []int64(nil), vals)
	})
	t.Run("LoadInt64s not found", func(t *testing.T) {
		names, err := NewSelect("id").From("dml_people").Where(Expr("name ='jdhsjdf'")).WithDBR(s.DB).LoadInt64s(context.TODO(), nil)
		assert.NoError(t, err)
		assert.Nil(t, names)
	})

	t.Run("LoadUint64s", func(t *testing.T) {
		names, err := NewSelect("id").From("dml_people").WithDBR(s.DB).LoadUint64s(context.TODO(), nil)
		assert.NoError(t, err)
		assert.Exactly(t, []uint64{1, 2}, names)
	})
	t.Run("LoadUint64s too many columns", func(t *testing.T) {
		vals, err := NewSelect("id", "email").From("dml_people").WithDBR(s.DB).LoadUint64s(context.TODO(), nil)
		assert.Error(t, err)
		assert.Exactly(t, []uint64(nil), vals)
	})
	t.Run("LoadUint64s not found", func(t *testing.T) {
		names, err := NewSelect("id").From("dml_people").Where(Expr("name ='jdhsjdf'")).WithDBR(s.DB).LoadUint64s(context.TODO(), nil)
		assert.NoError(t, err)
		assert.Nil(t, names)
	})

	t.Run("LoadFloat64s", func(t *testing.T) {
		names, err := NewSelect("id").From("dml_people").WithDBR(s.DB).LoadFloat64s(context.TODO(), nil)
		assert.NoError(t, err)
		assert.Exactly(t, []float64{1, 2}, names)
	})
	t.Run("LoadFloat64s too many columns", func(t *testing.T) {
		vals, err := NewSelect("id", "email").From("dml_people").WithDBR(s.DB).LoadFloat64s(context.TODO(), nil)
		assert.Error(t, err)
		assert.Exactly(t, []float64(nil), vals)
	})
	t.Run("LoadFloat64s not found", func(t *testing.T) {
		names, err := NewSelect("id").From("dml_people").Where(Expr("name ='jdhsjdf'")).WithDBR(s.DB).LoadFloat64s(context.TODO(), nil)
		assert.NoError(t, err)
		assert.Nil(t, names)
	})
}

func TestSelect_Join(t *testing.T) {
	s := createRealSessionWithFixtures(t, nil)
	defer testCloser(t, s)

	t.Run("inner, distinct, no cache, high prio", func(t *testing.T) {
		sqlObj := NewSelect("p1.*", "p2.*").FromAlias("dml_people", "p1").
			Distinct().StraightJoin().SQLNoCache().
			Join(
				MakeIdentifier("dml_people").Alias("p2"),
				Expr("`p2`.`id` = `p1`.`id`"),
				Column("p1.id").Int(42),
			)

		compareToSQL2(t, sqlObj, false,
			"SELECT DISTINCT STRAIGHT_JOIN SQL_NO_CACHE `p1`.*, `p2`.* FROM `dml_people` AS `p1` INNER JOIN `dml_people` AS `p2` ON (`p2`.`id` = `p1`.`id`) AND (`p1`.`id` = 42)",
		)
	})

	t.Run("inner", func(t *testing.T) {
		sqlObj := NewSelect("p1.*", "p2.*").FromAlias("dml_people", "p1").
			Join(
				MakeIdentifier("dml_people").Alias("p2"),
				Expr("`p2`.`id` = `p1`.`id`"),
				Column("p1.id").Int(42),
			)

		compareToSQL2(t, sqlObj, false,
			"SELECT `p1`.*, `p2`.* FROM `dml_people` AS `p1` INNER JOIN `dml_people` AS `p2` ON (`p2`.`id` = `p1`.`id`) AND (`p1`.`id` = 42)",
		)
	})

	t.Run("left", func(t *testing.T) {
		sqlObj := NewSelect("p1.*", "p2.name").FromAlias("dml_people", "p1").
			LeftJoin(
				MakeIdentifier("dml_people").Alias("p2"),
				Expr("`p2`.`id` = `p1`.`id`"),
				Column("p1.id").Int(42),
			)

		compareToSQL2(t, sqlObj, false,
			"SELECT `p1`.*, `p2`.`name` FROM `dml_people` AS `p1` LEFT JOIN `dml_people` AS `p2` ON (`p2`.`id` = `p1`.`id`) AND (`p1`.`id` = 42)",
		)
	})

	t.Run("right", func(t *testing.T) {
		sqlObj := NewSelect("p1.*").FromAlias("dml_people", "p1").
			AddColumnsAliases("p2.name", "p2Name", "p2.email", "p2Email", "id", "internalID").
			RightJoin(
				MakeIdentifier("dml_people").Alias("p2"),
				Expr("`p2`.`id` = `p1`.`id`"),
			)
		compareToSQL2(t, sqlObj, false,
			"SELECT `p1`.*, `p2`.`name` AS `p2Name`, `p2`.`email` AS `p2Email`, `id` AS `internalID` FROM `dml_people` AS `p1` RIGHT JOIN `dml_people` AS `p2` ON (`p2`.`id` = `p1`.`id`)",
		)
	})

	t.Run("using", func(t *testing.T) {
		sqlObj := NewSelect("p1.*").FromAlias("dml_people", "p1").
			AddColumnsAliases("p2.name", "p2Name", "p2.email", "p2Email").
			RightJoin(
				MakeIdentifier("dml_people").Alias("p2"),
				Columns("id", "email"),
			)
		compareToSQL2(t, sqlObj, false,
			"SELECT `p1`.*, `p2`.`name` AS `p2Name`, `p2`.`email` AS `p2Email` FROM `dml_people` AS `p1` RIGHT JOIN `dml_people` AS `p2` USING (`id`,`email`)",
		)
	})
}

func TestSelect_Locks(t *testing.T) {
	t.Run("LOCK IN SHARE MODE", func(t *testing.T) {
		s := NewSelect("p1.*").
			AddColumnsAliases("p2.name", "p2Name", "p2.email", "p2Email").
			FromAlias("dml_people", "p1").LockInShareMode()
		compareToSQL2(t, s, false,
			"SELECT `p1`.*, `p2`.`name` AS `p2Name`, `p2`.`email` AS `p2Email` FROM `dml_people` AS `p1` LOCK IN SHARE MODE",
		)
	})
	t.Run("FOR UPDATE", func(t *testing.T) {
		s := NewSelect("p1.*").
			AddColumnsAliases("p2.name", "p2Name", "p2.email", "p2Email").
			FromAlias("dml_people", "p1").ForUpdate()
		compareToSQL2(t, s, false,
			"SELECT `p1`.*, `p2`.`name` AS `p2Name`, `p2`.`email` AS `p2Email` FROM `dml_people` AS `p1` FOR UPDATE",
		)
	})
}

func TestSelect_Columns(t *testing.T) {
	t.Run("AddColumns, multiple args", func(t *testing.T) {
		s := NewSelect("a", "b")
		s.FromAlias("tableA", "tA")
		s.AddColumns("d,e, f", "g", "h", "i,j ,k")
		compareToSQL2(t, s, false,
			"SELECT `a`, `b`, `d,e, f`, `g`, `h`, `i,j ,k` FROM `tableA` AS `tA`",
		)
	})
	t.Run("AddColumns, each column itself", func(t *testing.T) {
		s := NewSelect("a", "b")
		s.FromAlias("tableA", "tA")
		s.AddColumns("d", "e", "f")
		compareToSQL2(t, s, false,
			"SELECT `a`, `b`, `d`, `e`, `f` FROM `tableA` AS `tA`",
		)
	})
	t.Run("AddColumnsAliases Expression Quoted", func(t *testing.T) {
		s := NewSelect().From("t3").
			AddColumnsAliases("x", "u", "y", "v").
			AddColumnsAliases("SUM(price)", "total_price")
		compareToSQL2(t, s, false,
			"SELECT `x` AS `u`, `y` AS `v`, `SUM(price)` AS `total_price` FROM `t3`",
		)
	})
	t.Run("AddColumns+AddColumnsConditions", func(t *testing.T) {
		s := NewSelect().From("t3").
			AddColumns("t3.name", "sku").
			AddColumnsConditions(Expr("SUM(price)").Alias("total_price"))
		compareToSQL2(t, s, false,
			"SELECT `t3`.`name`, `sku`, SUM(price) AS `total_price` FROM `t3`",
		)
	})

	t.Run("AddColumnsAliases multi", func(t *testing.T) {
		s := NewSelect().From("t3").
			AddColumnsAliases("t3.name", "t3Name", "t3.sku", "t3SKU")
		compareToSQL2(t, s, false,
			"SELECT `t3`.`name` AS `t3Name`, `t3`.`sku` AS `t3SKU` FROM `t3`",
		)
	})
	t.Run("AddColumnsAliases imbalanced", func(t *testing.T) {
		defer func() {
			if r := recover(); r != nil {
				if err, ok := r.(error); ok {
					assert.Error(t, err)
				} else {
					t.Errorf("Panic should contain an error but got:\n%+v", r)
				}
			} else {
				t.Error("Expecting a panic but got nothing")
			}
		}()
		NewSelect().From("t3").
			AddColumnsAliases("t3.name", "t3Name", "t3.sku")
	})

	t.Run("AddColumnsConditions", func(t *testing.T) {
		s := NewSelect().FromAlias("sales_bestsellers_aggregated_daily", "t3").
			AddColumnsConditions(Expr("DATE_FORMAT(t3.period, '%Y-%m-01')").Alias("period"))
		compareToSQL2(t, s, false,
			"SELECT DATE_FORMAT(t3.period, '%Y-%m-01') AS `period` FROM `sales_bestsellers_aggregated_daily` AS `t3`",
		)
	})
	t.Run("AddColumns with expression incorrect", func(t *testing.T) {
		s := NewSelect().AddColumns(" `t.value`", "`t`.`attribute_id`", "t.{column} AS `col_type`").FromAlias("catalog_product_entity_{type}", "t")
		compareToSQL2(t, s, false,
			"SELECT ` t`.`value`, `t`.`attribute_id`, `t`.`{column} AS col_type` FROM `catalog_product_entity_{type}` AS `t`",
		)
	})

	t.Run("AddColumnsConditions fails on interpolation", func(t *testing.T) {
		s := NewSelect().From("t3").
			AddColumns("t3.name", "sku").
			AddColumnsConditions(Expr("SUM(price)+?-?").Float64(3.14159).Alias("total_price"))
		compareToSQL2(t, s, true, "")
	})
}

func TestSelect_SubSelect(t *testing.T) {
	sub := NewSelect().From("catalog_category_product").
		AddColumns("entity_id").Where(Column("category_id").Int64(234))

	runner := func(op Op, wantSQL string) func(*testing.T) {
		c := Column("entity_id").Sub(sub)
		c.Operator = op
		return func(t *testing.T) {
			s := NewSelect("sku", "type_id").
				From("catalog_product_entity").
				Where(c)
			compareToSQL2(t, s, false, wantSQL)
		}
	}
	t.Run("IN", runner(In,
		"SELECT `sku`, `type_id` FROM `catalog_product_entity` WHERE (`entity_id` IN (SELECT `entity_id` FROM `catalog_category_product` WHERE (`category_id` = 234)))",
	))
	t.Run("EXISTS", runner(Exists,
		"SELECT `sku`, `type_id` FROM `catalog_product_entity` WHERE (`entity_id` EXISTS (SELECT `entity_id` FROM `catalog_category_product` WHERE (`category_id` = 234)))",
	))
	t.Run("NOT EXISTS", runner(NotExists,
		"SELECT `sku`, `type_id` FROM `catalog_product_entity` WHERE (`entity_id` NOT EXISTS (SELECT `entity_id` FROM `catalog_category_product` WHERE (`category_id` = 234)))",
	))
	t.Run("NOT EQUAL", runner(NotEqual,
		"SELECT `sku`, `type_id` FROM `catalog_product_entity` WHERE (`entity_id` != (SELECT `entity_id` FROM `catalog_category_product` WHERE (`category_id` = 234)))",
	))
	t.Run("NOT EQUAL", runner(Equal,
		"SELECT `sku`, `type_id` FROM `catalog_product_entity` WHERE (`entity_id` = (SELECT `entity_id` FROM `catalog_category_product` WHERE (`category_id` = 234)))",
	))
}

func TestSelect_Subselect_Complex(t *testing.T) {
	/* Something like:
	   SELECT
	     `t1`.`store_id`,
	     `t1`.`product_id`,
	     `t1`.`product_name`,
	     `t1`.`product_price`,
	     `t1`.`qty_ordered`
	   FROM (
	          SELECT
	            `t2`.`store_id`,
	            `t2`.`product_id`,
	            `t2`.`product_name`,
	            `t2`.`product_price`,
	            `t2`.`total_qty` AS `qty_ordered`
	          FROM (
	                 SELECT
	                   `t3`.`store_id`,
	                   `t3`.`product_id`,
	                   `t3`.`product_name`,
	                   AVG(`t3`.`product_price`) as `avg_price`,
	                   SUM(t3.qty_ordered) AS `total_qty`
	                 FROM `sales_bestsellers_aggregated_daily` AS `t3`
	                 GROUP BY `t3`.`store_id`,
	                   Date_format(t3.period, '%Y-%m-01'),
	                   `t3`.`product_id`
	                 ORDER BY `t3`.`store_id` ASC,
	                   Date_format(t3.period, '%Y-%m-01'),
	                   `total_qty` DESC
	               ) AS `t2`
	        ) AS `t1`
	*/

	t.Run("without args", func(t *testing.T) {
		sel3 := NewSelect().FromAlias("sales_bestsellers_aggregated_daily", "t3").
			Unsafe().
			AddColumnsConditions(Expr("DATE_FORMAT(t3.period, '%Y-%m-01')").Alias("period")).
			AddColumns("`t3`.`store_id`,`t3`.`product_id`,`t3`.`product_name`").
			AddColumnsConditions(
				Expr("AVG(`t3`.`product_price`)").Alias("avg_price"),
				Expr("SUM(t3.qty_ordered)").Alias("total_qty"),
			).
			GroupBy("t3.store_id").
			GroupBy("DATE_FORMAT(t3.period, '%Y-%m-01')").
			GroupBy("t3.product_id", "t3.product_name").
			OrderBy("t3.store_id").
			OrderBy("DATE_FORMAT(t3.period, '%Y-%m-01')").
			OrderByDesc("total_qty")

		sel2 := NewSelectWithDerivedTable(sel3, "t2").
			AddColumns("t2.period", "t2.store_id", "t2.product_id", "t2.product_name", "t2.avg_price").
			AddColumnsAliases("`t2`.`total_qty`", "`qty_ordered`")

		sel1 := NewSelectWithDerivedTable(sel2, "t1").
			AddColumns("t1.period", "t1.store_id", "t1.product_id", "t1.product_name", "t1.avg_price", "t1.qty_ordered").
			OrderBy("`t1`.period", "`t1`.product_id")

		compareToSQL2(t, sel1, false,
			"SELECT `t1`.`period`, `t1`.`store_id`, `t1`.`product_id`, `t1`.`product_name`, `t1`.`avg_price`, `t1`.`qty_ordered` FROM (SELECT `t2`.`period`, `t2`.`store_id`, `t2`.`product_id`, `t2`.`product_name`, `t2`.`avg_price`, `t2`.`total_qty` AS `qty_ordered` FROM (SELECT DATE_FORMAT(t3.period, '%Y-%m-01') AS `period`, `t3`.`store_id`,`t3`.`product_id`,`t3`.`product_name`, AVG(`t3`.`product_price`) AS `avg_price`, SUM(t3.qty_ordered) AS `total_qty` FROM `sales_bestsellers_aggregated_daily` AS `t3` GROUP BY `t3`.`store_id`, DATE_FORMAT(t3.period, '%Y-%m-01'), `t3`.`product_id`, `t3`.`product_name` ORDER BY `t3`.`store_id`, DATE_FORMAT(t3.period, '%Y-%m-01'), `total_qty` DESC) AS `t2`) AS `t1` ORDER BY `t1`.`period`, `t1`.`product_id`",
		)
	})

	t.Run("with args", func(t *testing.T) {
		// Full valid query which works in a M1 and M2 database.
		sel3 := NewSelect().FromAlias("sales_bestsellers_aggregated_daily", "t3").
			Unsafe().
			AddColumnsConditions(Expr("DATE_FORMAT(t3.period, '%Y-%m-01')").Alias("period")).
			AddColumns("`t3`.`store_id`,`t3`.`product_id`,`t3`.`product_name`").
			AddColumnsConditions(
				Expr("AVG(`t3`.`product_price`)").Alias("avg_price"),
				Expr("SUM(t3.qty_ordered)+?").Alias("total_qty").Float64(3.141),
			).
			GroupBy("t3.store_id").
			GroupBy("DATE_FORMAT(t3.period, '%Y-%m-01')").
			GroupBy("t3.product_id", "t3.product_name").
			Having(Expr("COUNT(*)>?").Int(3)).
			OrderBy("t3.store_id").
			OrderBy("DATE_FORMAT(t3.period, '%Y-%m-01')").
			OrderByDesc("total_qty").
			Where(Column("t3.store_id").In().Int64s(2, 3, 4))

		sel2 := NewSelectWithDerivedTable(sel3, "t2").
			AddColumns("t2.period", "t2.store_id", "t2.product_id", "t2.product_name", "t2.avg_price").
			AddColumnsAliases("t2.total_qty", "qty_ordered")

		sel1 := NewSelectWithDerivedTable(sel2, "t1").
			AddColumns("t1.period", "t1.store_id", "t1.product_id", "t1.product_name", "t1.avg_price", "t1.qty_ordered").
			OrderBy("`t1`.period", "`t1`.product_id")

		compareToSQL2(t, sel1, false,
			"SELECT `t1`.`period`, `t1`.`store_id`, `t1`.`product_id`, `t1`.`product_name`, `t1`.`avg_price`, `t1`.`qty_ordered` FROM (SELECT `t2`.`period`, `t2`.`store_id`, `t2`.`product_id`, `t2`.`product_name`, `t2`.`avg_price`, `t2`.`total_qty` AS `qty_ordered` FROM (SELECT DATE_FORMAT(t3.period, '%Y-%m-01') AS `period`, `t3`.`store_id`,`t3`.`product_id`,`t3`.`product_name`, AVG(`t3`.`product_price`) AS `avg_price`, SUM(t3.qty_ordered)+3.141 AS `total_qty` FROM `sales_bestsellers_aggregated_daily` AS `t3` WHERE (`t3`.`store_id` IN (2,3,4)) GROUP BY `t3`.`store_id`, DATE_FORMAT(t3.period, '%Y-%m-01'), `t3`.`product_id`, `t3`.`product_name` HAVING (COUNT(*)>3) ORDER BY `t3`.`store_id`, DATE_FORMAT(t3.period, '%Y-%m-01'), `total_qty` DESC) AS `t2`) AS `t1` ORDER BY `t1`.`period`, `t1`.`product_id`",
		)
	})
}

func TestSelect_Subselect_Compact(t *testing.T) {
	sel2 := NewSelect().FromAlias("sales_bestsellers_aggregated_daily", "t3").
		AddColumns("`t3`.`product_name`").
		Where(Column("t3.store_id").In().Int64s(2, 3, 4)).
		GroupBy("t3.store_id").
		Having(Expr("COUNT(*)>?").Int(5))

	sel := NewSelectWithDerivedTable(sel2, "t2").
		AddColumns("t2.product_name").
		Where(Column("t2.floatcol").Equal().Float64(3.14159))

	compareToSQL2(t, sel, false,
		"SELECT `t2`.`product_name` FROM (SELECT `t3`.`product_name` FROM `sales_bestsellers_aggregated_daily` AS `t3` WHERE (`t3`.`store_id` IN (2,3,4)) GROUP BY `t3`.`store_id` HAVING (COUNT(*)>5)) AS `t2` WHERE (`t2`.`floatcol` = 3.14159)",
	)
}

func TestSelect_ParenthesisOpen_Close(t *testing.T) {
	t.Run("beginning of WHERE", func(t *testing.T) {
		sel := NewSelect("a", "b").
			FromAlias("c", "cc").
			Where(
				ParenthesisOpen(),
				Column("d").Int(1),
				Column("e").Str("wat").Or(),
				ParenthesisClose(),
				Column("f").Float64(2.7182),
			).
			GroupBy("ab").
			Having(
				ParenthesisOpen(),
				Column("m").Int(33),
				Column("n").Str("wh3r3").Or(),
				ParenthesisClose(),
				Expr("j = k"),
			)
		compareToSQL2(t, sel, false,
			"SELECT `a`, `b` FROM `c` AS `cc` WHERE ((`d` = 1) OR (`e` = 'wat')) AND (`f` = 2.7182) GROUP BY `ab` HAVING ((`m` = 33) OR (`n` = 'wh3r3')) AND (j = k)",
		)
	})

	t.Run("end of WHERE", func(t *testing.T) {
		sel := NewSelect("a", "b").
			FromAlias("c", "cc").
			Where(
				Column("f").Float64(2.7182),
				ParenthesisOpen(),
				Column("d").Int(1),
				Column("e").Str("wat").Or(),
				ParenthesisClose(),
			).
			GroupBy("ab").
			Having(
				Expr("j = k"),
				ParenthesisOpen(),
				Column("m").Int(33),
				Column("n").Str("wh3r3").Or(),
				ParenthesisClose(),
			)
		compareToSQL2(t, sel, false,
			"SELECT `a`, `b` FROM `c` AS `cc` WHERE (`f` = 2.7182) AND ((`d` = 1) OR (`e` = 'wat')) GROUP BY `ab` HAVING (j = k) AND ((`m` = 33) OR (`n` = 'wh3r3'))",
		)
	})

	t.Run("middle of WHERE", func(t *testing.T) {
		sel := NewSelect("a", "b").
			FromAlias("c", "cc").
			Where(
				Column("f").Float64(2.7182),
				ParenthesisOpen(),
				Column("d").Int(1),
				Column("e").Str("wat").Or(),
				ParenthesisClose(),
				Column("p").Float64(3.141592),
			).
			GroupBy("ab").
			Having(
				Expr("j = k"),
				ParenthesisOpen(),
				Column("m").Int(33),
				Column("n").Str("wh3r3").Or(),
				ParenthesisClose(),
				Column("q").NotNull(),
			)
		compareToSQL2(t, sel, false,
			"SELECT `a`, `b` FROM `c` AS `cc` WHERE (`f` = 2.7182) AND ((`d` = 1) OR (`e` = 'wat')) AND (`p` = 3.141592) GROUP BY `ab` HAVING (j = k) AND ((`m` = 33) OR (`n` = 'wh3r3')) AND (`q` IS NOT NULL)",
		)
	})
}

func TestSelect_Count(t *testing.T) {
	t.Run("written count star gets quoted", func(t *testing.T) {
		compareToSQL2(t,
			NewSelect("count(*)").From("dml_people"),
			false,
			"SELECT `count(*)` FROM `dml_people`",
		)
	})
	t.Run("written count star gets not quoted Unsafe", func(t *testing.T) {
		compareToSQL2(t,
			NewSelect().Unsafe().AddColumns("count(*)").From("dml_people"),
			false,
			"SELECT count(*) FROM `dml_people`",
		)
	})
	t.Run("func count star", func(t *testing.T) {
		s := NewSelect("a", "b").Count().From("dml_people")
		compareToSQL2(t,
			s,
			false,
			"SELECT COUNT(*) AS `counted` FROM `dml_people`",
		)
	})
}

func TestSelect_DisableBuildCache(t *testing.T) {
	sel := NewSelect("a", "b").
		Distinct().
		FromAlias("c", "cc").
		Where(
			ParenthesisOpen(),
			Column("d").PlaceHolder(),
			Column("e").Str("wat").Or(),
			ParenthesisClose(),
			Column("f").Int(2),
			Column("g").Int(3),
			Column("h").In().Int64s(4, 5, 6),
		).
		GroupBy("ab").
		Having(
			ParenthesisOpen(),
			Column("m").Int(33),
			Column("n").Str("wh3r3").Or(),
			ParenthesisClose(),
			Expr("j = k"),
		).
		OrderBy("l").
		Limit(8, 7)

	const run1 = "SELECT DISTINCT `a`, `b` FROM `c` AS `cc` WHERE ((`d` = ?) OR (`e` = 'wat')) AND (`f` = 2) AND (`g` = 3) AND (`h` IN (4,5,6)) GROUP BY `ab` HAVING ((`m` = 33) OR (`n` = 'wh3r3')) AND (j = k) ORDER BY `l` LIMIT 8,7"
	const run2 = "SELECT DISTINCT `a`, `b` FROM `c` AS `cc` WHERE ((`d` = ?) OR (`e` = 'wat')) AND (`f` = 2) AND (`g` = 3) AND (`h` IN (4,5,6)) AND (`added_col` = 3.14159) GROUP BY `ab` HAVING ((`m` = 33) OR (`n` = 'wh3r3')) AND (j = k) ORDER BY `l` LIMIT 8,7"

	compareToSQL(t, sel.WithDBR(dbMock{}).TestWithArgs(87654), false,
		run1,
		"SELECT DISTINCT `a`, `b` FROM `c` AS `cc` WHERE ((`d` = 87654) OR (`e` = 'wat')) AND (`f` = 2) AND (`g` = 3) AND (`h` IN (4,5,6)) GROUP BY `ab` HAVING ((`m` = 33) OR (`n` = 'wh3r3')) AND (j = k) ORDER BY `l` LIMIT 8,7",
		int64(87654))

	sel.Where(
		Column("added_col").Float64(3.14159),
	)
	compareToSQL(t, sel.WithDBR(dbMock{}).TestWithArgs(87654), false, run2, "", int64(87654))
	// key2 still applies to the next 2 calls
	compareToSQL(t, sel.WithDBR(dbMock{}).TestWithArgs(87654), false, run2, "", int64(87654))
	compareToSQL(t, sel.WithDBR(dbMock{}).TestWithArgs(87654), false, run2,
		"SELECT DISTINCT `a`, `b` FROM `c` AS `cc` WHERE ((`d` = 87654) OR (`e` = 'wat')) AND (`f` = 2) AND (`g` = 3) AND (`h` IN (4,5,6)) AND (`added_col` = 3.14159) GROUP BY `ab` HAVING ((`m` = 33) OR (`n` = 'wh3r3')) AND (j = k) ORDER BY `l` LIMIT 8,7",
		int64(87654))
}

func TestSelect_NamedArguments(t *testing.T) {
	sel := NewSelect("config_id", "value").
		From("core_config_data").
		Where(
			Column("config_id1").Less().NamedArg(":configID"),
			Column("config_id2").Greater().NamedArg("configID"),
			Column("scope_id").Greater().Int(5),
			Column("value").Like().PlaceHolder(),
		)

	selDBR := sel.WithDBR(dbMock{})

	t.Run("With ID 3", func(t *testing.T) {
		compareToSQL2(t, selDBR.TestWithArgs(sql.Named("configID", 3), "GopherValue"), false,
			"SELECT `config_id`, `value` FROM `core_config_data` WHERE (`config_id1` < ?) AND (`config_id2` > ?) AND (`scope_id` > 5) AND (`value` LIKE ?)",
			int64(3), int64(3), "GopherValue",
		)
		assert.Exactly(t, []string{":configID", ":configID", "value"}, sel.qualifiedColumns, "qualifiedColumns should match")
	})
	t.Run("With ID 6", func(t *testing.T) {
		// Here positions are switched
		compareToSQL2(t, selDBR.TestWithArgs("G0pherValue", sql.Named("configID", 6)), false,
			"SELECT `config_id`, `value` FROM `core_config_data` WHERE (`config_id1` < ?) AND (`config_id2` > ?) AND (`scope_id` > 5) AND (`value` LIKE ?)",
			int64(6), int64(6), "G0pherValue",
		)
		assert.Exactly(t, []string{":configID", ":configID", "value"}, sel.qualifiedColumns, "qualifiedColumns should match")
	})
}

func TestSelect_SetRecord(t *testing.T) {
	p := &dmlPerson{
		ID:    6666,
		Name:  "Hans Wurst",
		Email: null.MakeString("hans@wurst.com"),
	}
	p2 := &dmlPerson{
		Dob: 1970,
	}

	t.Run("multiple args from record", func(t *testing.T) {
		sel := NewSelect("a", "b").
			FromAlias("dml_person", "dp").
			Join(MakeIdentifier("dml_group").Alias("dg"), Column("dp.id").PlaceHolder()).
			Where(
				Column("dg.dob").Greater().PlaceHolder(),
				Column("dg.size").Less().NamedArg("dbSIZE"),
				Column("age").Less().Int(56),
				ParenthesisOpen(),
				Column("dp.name").PlaceHolder(),
				Column("e").Str("wat").Or(),
				ParenthesisClose(),
				Column("f").LessOrEqual().Int(2),
				Column("g").Greater().Int(3),
				Column("h").In().Int64s(4, 5, 6),
			).
			GroupBy("ab").
			Having(
				Column("dp.email").PlaceHolder(),
				Column("n").Str("wh3r3"),
			).
			OrderBy("l").
			WithDBR(dbMock{}).TestWithArgs(Qualify("dp", p), Qualify("dg", p2), sql.Named("dbSIZE", uint(201801)))

		compareToSQL2(t, sel, false,
			"SELECT `a`, `b` FROM `dml_person` AS `dp` INNER JOIN `dml_group` AS `dg` ON (`dp`.`id` = ?) WHERE (`dg`.`dob` > ?) AND (`dg`.`size` < ?) AND (`age` < 56) AND ((`dp`.`name` = ?) OR (`e` = 'wat')) AND (`f` <= 2) AND (`g` > 3) AND (`h` IN (4,5,6)) GROUP BY `ab` HAVING (`dp`.`email` = ?) AND (`n` = 'wh3r3') ORDER BY `l`",
			int64(6666), int64(1970), int64(201801), "Hans Wurst", "hans@wurst.com",
		)
	})
	t.Run("single arg JOIN", func(t *testing.T) {
		sel := NewSelect("a").FromAlias("dml_people", "dp").
			Join(MakeIdentifier("dml_group").Alias("dg"), Column("dp.id").PlaceHolder(), Column("dg.name").Strs("XY%")).
			OrderBy("id").WithDBR(dbMock{}).TestWithArgs(Qualify("dp", p))

		compareToSQL2(t, sel, false,
			"SELECT `a` FROM `dml_people` AS `dp` INNER JOIN `dml_group` AS `dg` ON (`dp`.`id` = ?) AND (`dg`.`name` = ('XY%')) ORDER BY `id`",
			int64(6666),
		)
	})
	t.Run("single arg WHERE", func(t *testing.T) {
		sel := NewSelect("a").From("dml_people").
			Where(
				Column("id").PlaceHolder(),
			).
			OrderBy("id").WithDBR(dbMock{}).TestWithArgs(Qualify("", p))

		compareToSQL2(t, sel, false,
			"SELECT `a` FROM `dml_people` WHERE (`id` = ?) ORDER BY `id`",
			int64(6666),
		)
	})
	// t.Run("Warning when nothing got matched", func(t *testing.T) {
	//	sel := NewSelect("a").From("dml_people").
	//		Where(
	//			Column("id").PlaceHolder(),
	//		).
	//		WithRecords(Qualify("dml", p)).OrderBy("id")
	//
	//	compareToSQL2(t, sel, errors.Mismatch.Match, <-- TODO implement
	//		"SELECT `a` FROM `dml_people` WHERE (`id` = ?) ORDER BY `id`",
	//	)
	//})
	t.Run("HAVING", func(t *testing.T) {
		sel := NewSelect("a").From("dml_people").
			Having(
				Column("id").PlaceHolder(),
				Column("name").Like().PlaceHolder(),
			).
			OrderBy("id").WithDBR(dbMock{}).TestWithArgs(Qualify("", p))

		compareToSQL2(t, sel, false,
			"SELECT `a` FROM `dml_people` HAVING (`id` = ?) AND (`name` LIKE ?) ORDER BY `id`",
			int64(6666), "Hans Wurst",
		)
	})

	t.Run("slice as record", func(t *testing.T) {
		persons := &dmlPersons{
			Data: []*dmlPerson{
				{ID: 33, Name: "Muffin Hat", Email: null.MakeString("Muffin@Hat.head")},
				{ID: 44, Name: "Marianne Phyllis Finch", Email: null.MakeString("marianne@phyllis.finch")},
				{ID: 55, Name: "Daphne Augusta Perry", Email: null.MakeString("daphne@augusta.perry")},
			},
		}
		t.Run("one column in WHERE", func(t *testing.T) {
			compareToSQL2(t,
				NewSelect("name", "email").From("dml_person").
					Where(
						Column("id").In().PlaceHolder(),
					).
					WithDBR(dbMock{}).TestWithArgs(Qualify("", persons)),
				false,
				"SELECT `name`, `email` FROM `dml_person` WHERE (`id` IN ?)",
				int64(33), int64(44), int64(55),
			)
		})
		t.Run("two columns in WHERE", func(t *testing.T) {
			compareToSQL2(t,
				NewSelect("name", "email").From("dml_person").
					Where(
						Column("name").In().PlaceHolder(),
						Column("email").In().PlaceHolder(),
					).
					WithDBR(dbMock{}).TestWithArgs(Qualify("", persons)),
				false,
				"SELECT `name`, `email` FROM `dml_person` WHERE (`name` IN ?) AND (`email` IN ?)",
				// "SELECT `name`, `email` FROM `dml_person` WHERE (`name` IN ('Muffin Hat','Marianne Phyllis Finch','Daphne Augusta Perry')) AND (`email` IN ('Muffin@Hat.head','marianne@phyllis.finch','daphne@augusta.perry'))",
				"Muffin Hat", "Marianne Phyllis Finch", "Daphne Augusta Perry",
				"Muffin@Hat.head", "marianne@phyllis.finch", "daphne@augusta.perry",
			)
		})
		t.Run("three columns in WHERE", func(t *testing.T) {
			compareToSQL2(t,
				NewSelect("name", "email").From("dml_person").
					Where(
						Column("email").In().PlaceHolder(),
						Column("name").In().PlaceHolder(),
						Column("id").In().PlaceHolder(),
					).
					WithDBR(dbMock{}).TestWithArgs(Qualify("", persons)),
				false,
				"SELECT `name`, `email` FROM `dml_person` WHERE (`email` IN ?) AND (`name` IN ?) AND (`id` IN ?)",
				//"SELECT `name`, `email` FROM `dml_person` WHERE (`email` IN ('Muffin@Hat.head','marianne@phyllis.finch','daphne@augusta.perry')) AND (`name` IN ('Muffin Hat','Marianne Phyllis Finch','Daphne Augusta Perry')) AND (`id` IN (33,44,55))",
				"Muffin@Hat.head", "marianne@phyllis.finch", "daphne@augusta.perry",
				"Muffin Hat", "Marianne Phyllis Finch", "Daphne Augusta Perry",
				int64(33), int64(44), int64(55),
			)
		})
	})
}

func TestSelect_DBR_Load_Functions(t *testing.T) {
	s := createRealSessionWithFixtures(t, nil)
	defer testCloser(t, s)

	inA := NewInsert("dml_null_types").AddColumns("string_val", "int64_val", "float64_val", "time_val", "bool_val", "decimal_val").WithDBR(s.DB)
	res, err := inA.ExecContext(context.Background(),
		"A1", 11, 11.11, now(), true, 11.111,
		nil, nil, nil, nil, nil, nil,
		"A2", 22, 22.22, now(), false, 22.222,
		nil, nil, nil, nil, nil, nil,
		"-A3", -33, -33.33, now(), true, -33.333,
	)
	assert.NoError(t, err)
	lid, err := res.LastInsertId()
	assert.NoError(t, err)
	assert.Exactly(t, int64(1), lid)

	t.Run("LoadFloat64s all rows", func(t *testing.T) {
		var vals []float64
		vals, err := NewSelect("float64_val").From("dml_null_types").OrderBy("id").WithDBR(s.DB).LoadFloat64s(context.Background(), vals)
		assert.NoError(t, err)
		assert.Exactly(t, []float64{11.11, 22.22, -33.33}, vals)
	})

	t.Run("LoadFloat64s IN with one arg", func(t *testing.T) {
		var vals []float64
		vals, err := NewSelect("float64_val").From("dml_null_types").Where(
			Column("int64_val").In().PlaceHolders(1),
			// do not interpolate because we want a prepare statement
		).WithDBR(s.DB).LoadFloat64s(context.Background(), vals, []int{11})
		assert.NoError(t, err)
		// the mariaDB field type of float64_val is float and go-sql-driver/mysql
		// converts it to float32 as we use internally null.Float64 we must convert
		// float32 to float64 and so loose precision. If the mariaDB field type of
		// field float64_val would be double, then we have float64.
		assert.Exactly(t, "11.11", fmt.Sprintf("%.2f", vals[0]))
	})

	t.Run("LoadInt64s", func(t *testing.T) {
		var vals []int64
		vals, err := NewSelect("int64_val").From("dml_null_types").OrderBy("id").WithDBR(s.DB).LoadInt64s(context.Background(), vals)
		assert.NoError(t, err)
		assert.Exactly(t, []int64{11, 22, -33}, vals)
	})

	t.Run("LoadUint64s", func(t *testing.T) {
		var vals []uint64
		vals, err := NewSelect("int64_val").From("dml_null_types").Where(Column("int64_val").GreaterOrEqual().Int(0)).OrderBy("id").
			WithDBR(s.DB).LoadUint64s(context.Background(), vals)
		assert.NoError(t, err)
		assert.Exactly(t, []uint64{11, 22}, vals)
	})
	t.Run("LoadStrings found", func(t *testing.T) {
		var vals []string
		vals, err := NewSelect("string_val").From("dml_null_types").OrderBy("id").WithDBR(s.DB).LoadStrings(context.Background(), vals)
		assert.NoError(t, err)
		assert.Exactly(t, []string{"A1", "A2", "-A3"}, vals)
	})
	t.Run("LoadStrings not found", func(t *testing.T) {
		var vals []string
		vals, err := NewSelect("string_val").From("dml_null_types").Where(
			Column("int64_val").Equal().Int64(-34),
		).OrderBy("id").WithDBR(s.DB).LoadStrings(context.Background(), vals)
		assert.NoError(t, err)
		assert.Exactly(t, []string(nil), vals)
	})
}
