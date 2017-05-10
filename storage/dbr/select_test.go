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
	"context"
	"testing"

	"github.com/corestoreio/errors"
	"github.com/corestoreio/log"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSelectBasicToSQL(t *testing.T) {
	t.Parallel()
	s := createFakeSession()

	sel := s.Select("a", "b").From("c").Where(Column("id", argInt(1)))
	sql, args, err := sel.ToSQL()
	assert.NoError(t, err)
	assert.Equal(t, "SELECT `a`, `b` FROM `c` WHERE (`id` = ?)", sql)
	assert.Equal(t, []interface{}{int64(1)}, args.Interfaces())
}

func TestSelectFullToSQL(t *testing.T) {
	t.Parallel()

	sel := NewSelect("a", "b").
		Distinct().
		From("c", "cc").
		Where(
			ParenthesisOpen(),
			Column("d", argInt(1)),
			Column("e", ArgString("wat")).Or(),
			ParenthesisClose(),
			Eq{"f": argInt(2)}, Eq{"g": argInt(3)},
		).
		Where(Eq{"h": ArgInt64(4, 5, 6).Operator(In)}).
		GroupBy("ab").
		Having(
			ParenthesisOpen(),
			Column("m", argInt(33)),
			Column("n", ArgString("wh3r3")).Or(),
			ParenthesisClose(),
			Expression("j = k"),
		).
		OrderBy("l").
		Limit(7).
		Offset(8)

	sql, args, err := sel.ToSQL()
	assert.NoError(t, err)
	assert.Equal(t, "SELECT DISTINCT `a`, `b` FROM `c` AS `cc` WHERE ((`d` = ?) OR (`e` = ?)) AND (`f` = ?) AND (`g` = ?) AND (`h` IN ?) GROUP BY `ab` HAVING ((`m` = ?) OR (`n` = ?)) AND (j = k) ORDER BY `l` LIMIT 7 OFFSET 8", sql)
	assert.Equal(t, []interface{}{int64(1), "wat", int64(2), int64(3), int64(4), int64(5), int64(6), int64(33), "wh3r3"}, args.Interfaces())
}

func TestSelect_Interpolate(t *testing.T) {
	t.Parallel()

	sql, args, err := NewSelect("a", "b").
		Distinct().
		From("c", "cc").
		Where(
			ParenthesisOpen(),
			Column("d", argInt(1)),
			Column("e", ArgString("wat")).Or(),
			ParenthesisClose(),
			Eq{"f": argInt(2)}, Eq{"g": argInt(3)},
		).
		Where(Eq{"h": ArgInt64(4, 5, 6).Operator(In)}).
		GroupBy("ab").
		Having(
			ParenthesisOpen(),
			Column("m", argInt(33)),
			Column("n", ArgString("wh3r3")).Or(),
			ParenthesisClose(),
			Column("j = k"),
		).
		OrderBy("l").
		Limit(7).
		Offset(8).
		Interpolate().
		ToSQL()

	require.NoError(t, err)
	assert.Equal(t, "SELECT DISTINCT `a`, `b` FROM `c` AS `cc` WHERE ((`d` = 1) OR (`e` = 'wat')) AND (`f` = 2) AND (`g` = 3) AND (`h` IN (4,5,6)) GROUP BY `ab` HAVING ((`m` = 33) OR (`n` = 'wh3r3')) AND (`j = k`) ORDER BY `l` LIMIT 7 OFFSET 8", sql)
	assert.Nil(t, args)
}

func TestSelectPaginateOrderDirToSQL(t *testing.T) {
	t.Parallel()
	s := createFakeSession()

	sql, args, err := s.Select("a", "b").
		From("c").
		Where(Column("d", argInt(1))).
		Paginate(1, 20).
		OrderByDesc("id").
		ToSQL()
	assert.NoError(t, err)
	assert.Equal(t, "SELECT `a`, `b` FROM `c` WHERE (`d` = ?) ORDER BY `id` DESC LIMIT 20 OFFSET 0", sql)
	assert.Equal(t, []interface{}{int64(1)}, args.Interfaces())

	sql, args, err = s.Select("a", "b").
		From("c").
		Where(Column("d", argInt(1))).
		Paginate(3, 30).
		OrderBy("id").
		ToSQL()
	assert.NoError(t, err)
	assert.Equal(t, "SELECT `a`, `b` FROM `c` WHERE (`d` = ?) ORDER BY `id` LIMIT 30 OFFSET 60", sql)
	assert.Equal(t, []interface{}{int64(1)}, args.Interfaces())
}

func TestSelectNoWhereSQL(t *testing.T) {
	t.Parallel()
	s := createFakeSession()

	sql, args, err := s.Select("a", "b").From("c").ToSQL()
	assert.NoError(t, err)
	assert.Equal(t, "SELECT `a`, `b` FROM `c`", sql)
	assert.Equal(t, []interface{}(nil), args.Interfaces())
}

func TestSelectMultiHavingSQL(t *testing.T) {
	t.Parallel()
	s := createFakeSession()

	sql, args, err := s.Select("a", "b").From("c").
		Where(Column("p", argInt(1))).
		GroupBy("z").Having(Column("z`z", argInt(2)), Column("y", argInt(3))).ToSQL()
	assert.NoError(t, err)
	assert.Equal(t, "SELECT `a`, `b` FROM `c` WHERE (`p` = ?) GROUP BY `z` HAVING (`zz` = ?) AND (`y` = ?)", sql)
	assert.Equal(t, []interface{}{int64(1), int64(2), int64(3)}, args.Interfaces())
}

func TestSelectMultiOrderSQL(t *testing.T) {
	t.Parallel()
	s := createFakeSession()

	sql, args, err := s.Select("a", "b").From("c").OrderBy("name").OrderByDesc("id").ToSQL()
	assert.NoError(t, err)
	assert.Equal(t, "SELECT `a`, `b` FROM `c` ORDER BY `name`, `id` DESC", sql)
	assert.Equal(t, []interface{}(nil), args.Interfaces())
}

func TestSelect_OrderByDeactivated(t *testing.T) {
	t.Parallel()

	sql, args, err := NewSelect("a", "b").From("c").OrderBy("name").OrderByDeactivated().ToSQL()
	assert.NoError(t, err)
	assert.Equal(t, "SELECT `a`, `b` FROM `c` ORDER BY NULL", sql)
	assert.Exactly(t, Arguments{}, args)
}

func TestSelect_ConditionColumn(t *testing.T) {
	t.Parallel()
	// TODO rewrite test to use every type which implements interface Argument and every operator

	s := createFakeSession()
	runner := func(arg Argument, wantSQL string, wantVal []interface{}) func(*testing.T) {
		return func(t *testing.T) {
			sql, args, err := s.Select("a", "b").From("c").Where(Column("d", arg)).ToSQL()
			assert.NoError(t, err)
			assert.Exactly(t, wantSQL, sql)
			assert.Exactly(t, wantVal, args.Interfaces())

		}
	}
	t.Run("single int64", runner(
		argInt64(33),
		"SELECT `a`, `b` FROM `c` WHERE (`d` = ?)",
		[]interface{}{int64(33)},
	))
	t.Run("IN int64", runner(
		ArgInt64(33, 44).Operator(In),
		"SELECT `a`, `b` FROM `c` WHERE (`d` IN ?)",
		[]interface{}{int64(33), int64(44)},
	))
	t.Run("single float64", runner(
		ArgFloat64(33),
		"SELECT `a`, `b` FROM `c` WHERE (`d` = ?)",
		[]interface{}{float64(33)},
	))
	t.Run("IN float64", runner(
		ArgFloat64(33, 44).Operator(In),
		"SELECT `a`, `b` FROM `c` WHERE (`d` IN ?)",
		[]interface{}{float64(33), float64(44)},
	))
	t.Run("NOT IN float64", runner(
		ArgFloat64(33, 44).Operator(NotIn),
		"SELECT `a`, `b` FROM `c` WHERE (`d` NOT IN ?)",
		[]interface{}{float64(33), float64(44)},
	))
	t.Run("single int", runner(
		argInt(33),
		"SELECT `a`, `b` FROM `c` WHERE (`d` = ?)",
		[]interface{}{int64(33)},
	))
	t.Run("IN int", runner(
		ArgInt(33, 44).Operator(In),
		"SELECT `a`, `b` FROM `c` WHERE (`d` IN ?)",
		[]interface{}{int64(33), int64(44)},
	))
	t.Run("single string", runner(
		ArgString("w"),
		"SELECT `a`, `b` FROM `c` WHERE (`d` = ?)",
		[]interface{}{"w"},
	))
	t.Run("IN string", runner(
		ArgString("x", "y").Operator(In),
		"SELECT `a`, `b` FROM `c` WHERE (`d` IN ?)",
		[]interface{}{"x", "y"},
	))

	t.Run("BETWEEN int64", runner(
		ArgInt64(5, 6).Operator(Between),
		"SELECT `a`, `b` FROM `c` WHERE (`d` BETWEEN ? AND ?)",
		[]interface{}{int64(5), int64(6)},
	))
	t.Run("NOT BETWEEN int64", runner(
		ArgInt64(5, 6).Operator(NotBetween),
		"SELECT `a`, `b` FROM `c` WHERE (`d` NOT BETWEEN ? AND ?)",
		[]interface{}{int64(5), int64(6)},
	))

	t.Run("LIKE string", runner(
		ArgString("x%").Operator(Like),
		"SELECT `a`, `b` FROM `c` WHERE (`d` LIKE ?)",
		[]interface{}{"x%"},
	))
	t.Run("NOT LIKE string", runner(
		ArgString("x%").Operator(NotLike),
		"SELECT `a`, `b` FROM `c` WHERE (`d` NOT LIKE ?)",
		[]interface{}{"x%"},
	))

	t.Run("Less float64", runner(
		ArgFloat64(5.1).Operator(Less),
		"SELECT `a`, `b` FROM `c` WHERE (`d` < ?)",
		[]interface{}{float64(5.1)},
	))
	t.Run("Greater float64", runner(
		ArgFloat64(5.1).Operator(Greater),
		"SELECT `a`, `b` FROM `c` WHERE (`d` > ?)",
		[]interface{}{float64(5.1)},
	))
	t.Run("LessOrEqual float64", runner(
		ArgFloat64(5.1).Operator(LessOrEqual),
		"SELECT `a`, `b` FROM `c` WHERE (`d` <= ?)",
		[]interface{}{float64(5.1)},
	))
	t.Run("GreaterOrEqual float64", runner(
		ArgFloat64(5.1).Operator(GreaterOrEqual),
		"SELECT `a`, `b` FROM `c` WHERE (`d` >= ?)",
		[]interface{}{float64(5.1)},
	))

}

func TestSelect_Null(t *testing.T) {
	s := createFakeSession()

	t.Run("col is null", func(t *testing.T) {
		sql, args, err := s.Select("a", "b").From("c").Where(Column("r", ArgNull())).ToSQL()
		assert.NoError(t, err)
		assert.Exactly(t, "SELECT `a`, `b` FROM `c` WHERE (`r` IS NULL)", sql)
		assert.Exactly(t, []interface{}(nil), args.Interfaces())
	})

	t.Run("col is not null", func(t *testing.T) {
		sql, args, err := s.Select("a", "b").From("c").Where(Column("r", ArgNull().Operator(NotNull))).ToSQL()
		assert.NoError(t, err)
		assert.Exactly(t, "SELECT `a`, `b` FROM `c` WHERE (`r` IS NOT NULL)", sql)
		assert.Exactly(t, []interface{}(nil), args.Interfaces())
	})

	t.Run("complex", func(t *testing.T) {
		sql, args, err := s.Select("a", "b").From("c").
			Where(
				Column("r", ArgNull()),
				Column("d", argInt(3)),
				Column("ab", ArgNull()),
				Column("w", ArgNull().Operator(NotNull)),
			).ToSQL()
		assert.NoError(t, err)
		assert.Exactly(t, "SELECT `a`, `b` FROM `c` WHERE (`r` IS NULL) AND (`d` = ?) AND (`ab` IS NULL) AND (`w` IS NOT NULL)", sql)
		assert.Exactly(t, []interface{}{int64(3)}, args.Interfaces())
	})
}

func TestSelectWhereMapSQL(t *testing.T) {
	s := createFakeSession()

	t.Run("one", func(t *testing.T) {
		sql, args, err := s.Select("a").From("b").Where(Eq{"a": argInt(1)}).ToSQL()
		assert.NoError(t, err)
		assert.Equal(t, "SELECT `a` FROM `b` WHERE (`a` = ?)", sql)
		assert.Equal(t, []interface{}{int64(1)}, args.Interfaces())
	})

	t.Run("two", func(t *testing.T) {
		sql, args, err := s.Select("a").From("b").Where(Eq{"a": argInt(1), "b": ArgBool(true)}).ToSQL()
		assert.NoError(t, err)
		if sql == "SELECT `a` FROM `b` WHERE (`a` = ?) AND (`b` = ?)" {
			assert.Equal(t, []interface{}{int64(1), true}, args.Interfaces())
		} else {
			assert.Equal(t, "SELECT `a` FROM `b` WHERE (`b` = ?) AND (`a` = ?)", sql)
			assert.Equal(t, []interface{}{true, int64(1)}, args.Interfaces())
		}
	})

	t.Run("one nil", func(t *testing.T) {
		sql, args, err := s.Select("a").From("b").Where(Eq{"a": nil}).ToSQL()
		assert.NoError(t, err)
		assert.Equal(t, "SELECT `a` FROM `b` WHERE (`a` IS NULL)", sql)
		assert.Equal(t, []interface{}(nil), args.Interfaces())
	})

	t.Run("one IN", func(t *testing.T) {
		sql, args, err := s.Select("a").From("b").Where(Eq{"a": ArgInt(1, 2, 3).Operator(In)}).ToSQL()
		assert.NoError(t, err)
		assert.Equal(t, "SELECT `a` FROM `b` WHERE (`a` IN ?)", sql)
		assert.Equal(t, []interface{}{int64(1), int64(2), int64(3)}, args.Interfaces())
	})

	t.Run("no values", func(t *testing.T) {
		// NOTE: a has no valid values, we want a query that returns nothing
		// TODO(CyS): revise architecture and behaviour ... maybe
		sql, args, err := s.Select("a").From("b").Where(Eq{"a": ArgInt()}).ToSQL()
		assert.NoError(t, err)
		//assert.Equal(t, "SELECT a FROM `b` WHERE (1=0)", sql)
		assert.Equal(t, "SELECT `a` FROM `b` WHERE (`a` = ?)", sql)
		assert.Equal(t, []interface{}{}, args.Interfaces())
	})

	t.Run("empty ArgInt", func(t *testing.T) {
		// see subtest above "no values" and its TODO
		var iVal []int
		sql, args, err := s.Select("a").From("b").Where(Eq{"a": ArgInt(iVal...)}).ToSQL()
		assert.NoError(t, err)
		assert.Equal(t, "SELECT `a` FROM `b` WHERE (`a` = ?)", sql)
		assert.Equal(t, []interface{}{}, args.Interfaces())
	})

	t.Run("Map nil arg", func(t *testing.T) {
		sql, args, err := s.Select("a").From("b").
			Where(Eq{"a": nil}).
			Where(Eq{"b": ArgBool(false)}).
			Where(Eq{"c": ArgNull()}).
			Where(Eq{"d": ArgNull().Operator(NotNull)}).
			ToSQL()
		assert.NoError(t, err)
		assert.Equal(t, "SELECT `a` FROM `b` WHERE (`a` IS NULL) AND (`b` = ?) AND (`c` IS NULL) AND (`d` IS NOT NULL)", sql)
		assert.Equal(t, []interface{}{false}, args.Interfaces())
	})
}

func TestSelectWhereEqSQL(t *testing.T) {
	s := createFakeSession()

	sql, args, err := s.Select("a").From("b").Where(Eq{"a": argInt(1), "b": ArgInt64(1, 2, 3).Operator(In)}).ToSQL()
	assert.NoError(t, err)
	if sql == "SELECT `a` FROM `b` WHERE (`a` = ?) AND (`b` IN ?)" {
		assert.Equal(t, []interface{}{int64(1), int64(1), int64(2), int64(3)}, args.Interfaces())
	} else {
		assert.Equal(t, sql, "SELECT `a` FROM `b` WHERE (`b` IN ?) AND (`a` = ?)")
		assert.Equal(t, []interface{}{int64(1), int64(2), int64(3), int64(1)}, args.Interfaces())
	}
}

func TestSelectBySQL(t *testing.T) {
	s := createFakeSession()

	sql, args, err := s.SelectBySQL("SELECT * FROM users WHERE x = 1").ToSQL()
	assert.NoError(t, err)
	assert.Equal(t, sql, "SELECT * FROM users WHERE x = 1")
	assert.Equal(t, []interface{}(nil), args.Interfaces())

	sql, args, err = s.SelectBySQL("SELECT * FROM users WHERE x = ? AND y IN ?", argInt(9), ArgInt(5, 6, 7)).ToSQL()
	assert.NoError(t, err)
	assert.Equal(t, sql, "SELECT * FROM users WHERE x = ? AND y IN ?")
	assert.Equal(t, []interface{}{int64(9), int64(5), int64(6), int64(7)}, args.Interfaces())

	// Doesn't fix shit if it'ab broken:
	sql, args, err = s.SelectBySQL("wat", argInt(9), ArgInt(5, 6, 7)).ToSQL()
	assert.NoError(t, err)
	assert.Equal(t, sql, "wat")
	assert.Equal(t, []interface{}{int64(9), int64(5), int64(6), int64(7)}, args.Interfaces())
}

func TestSelectVarieties(t *testing.T) {
	s := createFakeSession()

	sql, _, err := s.Select("id, name, email").From("users").ToSQL()
	assert.NoError(t, err)
	assert.Equal(t, "SELECT `id, name, email` FROM `users`", sql)
	sql2, _, err2 := s.Select("id", "name", "email").From("users").ToSQL()
	assert.NoError(t, err2)
	assert.NotEqual(t, sql, sql2)
	assert.Equal(t, "SELECT `id`, `name`, `email` FROM `users`", sql2)
}

func TestSelect_Load_Slice_Scanner(t *testing.T) {
	s := createRealSessionWithFixtures()

	var people dbrPersons
	count, err := s.Select("id", "name", "email").From("dbr_people").OrderBy("id").Load(context.TODO(), &people)

	assert.NoError(t, err)
	assert.Equal(t, count, 2)

	assert.Equal(t, len(people.Data), 2)
	if len(people.Data) == 2 {
		// Make sure that the Ids are set. It'ab possible (maybe?) that different DBs set ids differently so
		// don't assume they're 1 and 2.
		assert.True(t, people.Data[0].ID > 0)
		assert.True(t, people.Data[1].ID > people.Data[0].ID)

		assert.Equal(t, "Jonathan", people.Data[0].Name)
		assert.True(t, people.Data[0].Email.Valid)
		assert.Equal(t, "jonathan@uservoice.com", people.Data[0].Email.String)
		assert.Equal(t, "Dmitri", people.Data[1].Name)
		assert.True(t, people.Data[1].Email.Valid)
		assert.Equal(t, "zavorotni@jadius.com", people.Data[1].Email.String)
	}
}

func TestSelect_Load_Single(t *testing.T) {
	s := createRealSessionWithFixtures()

	t.Run("found", func(t *testing.T) {
		var person dbrPerson
		_, err := s.Select("id", "name", "email").From("dbr_people").
			Where(Column("email", ArgString("jonathan@uservoice.com"))).Load(context.TODO(), &person)
		assert.NoError(t, err)
		assert.True(t, person.ID > 0)
		assert.Equal(t, "Jonathan", person.Name)
		assert.True(t, person.Email.Valid)
		assert.Equal(t, "jonathan@uservoice.com", person.Email.String)
	})

	t.Run("not found", func(t *testing.T) {
		var person2 dbrPerson
		count, err := s.Select("id", "name", "email").From("dbr_people").
			Where(Column("email", ArgString("dontexist@uservoice.com"))).Load(context.TODO(), &person2)

		require.NoError(t, err, "%+v", err)
		assert.Exactly(t, dbrPerson{}, person2)
		assert.Empty(t, count, "Should have no rows loaded")
	})

}

func TestSelectBySQL_Load_Slice(t *testing.T) {
	s := createRealSessionWithFixtures()

	var people dbrPersons
	count, err := s.SelectBySQL("SELECT name FROM dbr_people WHERE email = ?", ArgString("jonathan@uservoice.com")).Load(context.TODO(), &people)

	assert.NoError(t, err)
	assert.Equal(t, count, 1)
	if len(people.Data) == 1 {
		assert.Equal(t, "Jonathan", people.Data[0].Name)
		assert.Equal(t, int64(0), people.Data[0].ID)       // not set
		assert.Equal(t, false, people.Data[0].Email.Valid) // not set
		assert.Equal(t, "", people.Data[0].Email.String)   // not set
	}
}

func TestSelect_LoadType_Single(t *testing.T) {
	s := createRealSessionWithFixtures()

	t.Run("LoadString", func(t *testing.T) {
		name, err := s.Select("name").From("dbr_people").Where(Expression("email = 'jonathan@uservoice.com'")).LoadString(context.TODO())
		assert.NoError(t, err)
		assert.Equal(t, "Jonathan", name)
	})
	t.Run("LoadString too many columns", func(t *testing.T) {
		name, err := s.Select("name", "email").From("dbr_people").Where(Expression("email = 'jonathan@uservoice.com'")).LoadString(context.TODO())
		assert.Error(t, err, "%+v", err)
		assert.Empty(t, name)
	})
	t.Run("LoadString not found", func(t *testing.T) {
		name, err := s.Select("name").From("dbr_people").Where(Expression("email = 'notfound@example.com'")).LoadString(context.TODO())
		assert.True(t, errors.IsNotFound(err), "%+v", err)
		assert.Empty(t, name)
	})

	t.Run("LoadInt64", func(t *testing.T) {
		id, err := s.Select("id").From("dbr_people").Limit(1).LoadInt64(context.TODO())
		assert.NoError(t, err)
		assert.True(t, id > 0)
	})
	t.Run("LoadInt64 too many columns", func(t *testing.T) {
		id, err := s.Select("id", "email").From("dbr_people").Limit(1).LoadInt64(context.TODO())
		assert.Error(t, err, "%+v", err)
		assert.Empty(t, id)
	})
	t.Run("LoadInt64 not found", func(t *testing.T) {
		id, err := s.Select("id").From("dbr_people").Where(Expression("id=236478326")).LoadInt64(context.TODO())
		assert.True(t, errors.IsNotFound(err), "%+v", err)
		assert.Empty(t, id)
	})

	t.Run("LoadUint64", func(t *testing.T) {
		id, err := s.Select("id").From("dbr_people").Limit(1).LoadUint64(context.TODO())
		assert.NoError(t, err)
		assert.True(t, id > 0)
	})
	t.Run("LoadUint64 too many columns", func(t *testing.T) {
		id, err := s.Select("id", "email").From("dbr_people").Limit(1).LoadUint64(context.TODO())
		assert.Error(t, err, "%+v", err)
		assert.Empty(t, id)
	})
	t.Run("LoadUint64 not found", func(t *testing.T) {
		id, err := s.Select("id").From("dbr_people").Where(Expression("id=236478326")).LoadUint64(context.TODO())
		assert.True(t, errors.IsNotFound(err), "%+v", err)
		assert.Empty(t, id)
	})

	t.Run("LoadFloat64", func(t *testing.T) {
		id, err := s.Select("id").From("dbr_people").Limit(1).LoadFloat64(context.TODO())
		assert.NoError(t, err)
		assert.True(t, id > 0)
	})
	t.Run("LoadFloat64 too many columns", func(t *testing.T) {
		id, err := s.Select("id", "email").From("dbr_people").Limit(1).LoadFloat64(context.TODO())
		assert.Error(t, err, "%+v", err)
		assert.Empty(t, id)
	})
	t.Run("LoadFloat64 not found", func(t *testing.T) {
		id, err := s.Select("id").From("dbr_people").Where(Expression("id=236478326")).LoadFloat64(context.TODO())
		assert.True(t, errors.IsNotFound(err), "%+v", err)
		assert.Empty(t, id)
	})
}

func TestSelect_LoadType_Slices(t *testing.T) {
	s := createRealSessionWithFixtures()

	t.Run("LoadStrings", func(t *testing.T) {
		names, err := s.Select("name").From("dbr_people").LoadStrings(context.TODO())
		assert.NoError(t, err)
		assert.Equal(t, []string{"Jonathan", "Dmitri"}, names)
	})
	t.Run("LoadStrings too many columns", func(t *testing.T) {
		vals, err := s.Select("name", "email").From("dbr_people").LoadStrings(context.TODO())
		assert.Error(t, err, "%+v", err)
		assert.Exactly(t, []string(nil), vals)
	})
	t.Run("LoadStrings not found", func(t *testing.T) {
		names, err := s.Select("name").From("dbr_people").Where(Expression("name ='jdhsjdf'")).LoadStrings(context.TODO())
		assert.NoError(t, err)
		assert.Equal(t, []string{}, names)
	})

	t.Run("LoadInt64s", func(t *testing.T) {
		names, err := s.Select("id").From("dbr_people").LoadInt64s(context.TODO())
		assert.NoError(t, err)
		assert.Equal(t, []int64{1, 2}, names)
	})
	t.Run("LoadInt64s too many columns", func(t *testing.T) {
		vals, err := s.Select("id", "email").From("dbr_people").LoadInt64s(context.TODO())
		assert.Error(t, err, "%+v", err)
		assert.Exactly(t, []int64(nil), vals)
	})
	t.Run("LoadInt64s not found", func(t *testing.T) {
		names, err := s.Select("id").From("dbr_people").Where(Expression("name ='jdhsjdf'")).LoadInt64s(context.TODO())
		assert.NoError(t, err)
		assert.Equal(t, []int64{}, names)
	})

	t.Run("LoadUint64s", func(t *testing.T) {
		names, err := s.Select("id").From("dbr_people").LoadUint64s(context.TODO())
		assert.NoError(t, err)
		assert.Equal(t, []uint64{1, 2}, names)
	})
	t.Run("LoadUint64s too many columns", func(t *testing.T) {
		vals, err := s.Select("id", "email").From("dbr_people").LoadUint64s(context.TODO())
		assert.Error(t, err, "%+v", err)
		assert.Exactly(t, []uint64(nil), vals)
	})
	t.Run("LoadUint64s not found", func(t *testing.T) {
		names, err := s.Select("id").From("dbr_people").Where(Expression("name ='jdhsjdf'")).LoadUint64s(context.TODO())
		assert.NoError(t, err)
		assert.Equal(t, []uint64{}, names)
	})

	t.Run("LoadFloat64s", func(t *testing.T) {
		names, err := s.Select("id").From("dbr_people").LoadFloat64s(context.TODO())
		assert.NoError(t, err)
		assert.Equal(t, []float64{1, 2}, names)
	})
	t.Run("LoadFloat64s too many columns", func(t *testing.T) {
		vals, err := s.Select("id", "email").From("dbr_people").LoadFloat64s(context.TODO())
		assert.Error(t, err, "%+v", err)
		assert.Exactly(t, []float64(nil), vals)
	})
	t.Run("LoadFloat64s not found", func(t *testing.T) {
		names, err := s.Select("id").From("dbr_people").Where(Expression("name ='jdhsjdf'")).LoadFloat64s(context.TODO())
		assert.NoError(t, err)
		assert.Equal(t, []float64{}, names)
	})

}

func TestSelectJoin(t *testing.T) {
	s := createRealSessionWithFixtures()

	t.Run("inner, distinct, no cache, high proi", func(t *testing.T) {
		sqlObj := s.
			Select("p1.*", "p2.*").
			Distinct().StraightJoin().SQLNoCache().
			From("dbr_people", "p1").
			Join(
				MakeAlias("dbr_people", "p2"),
				Column("`p2`.`id` = `p1`.`id`"),
				Column("p1.id", argInt(42)),
			)

		sql, _, err := sqlObj.ToSQL()
		assert.NoError(t, err)
		assert.Equal(t,
			"SELECT DISTINCT STRAIGHT_JOIN SQL_NO_CACHE `p1`.*, `p2`.* FROM `dbr_people` AS `p1` INNER JOIN `dbr_people` AS `p2` ON (`p2`.`id` = `p1`.`id`) AND (`p1`.`id` = ?)",
			sql,
		)
	})

	t.Run("inner", func(t *testing.T) {
		sqlObj := s.
			Select("p1.*", "p2.*").
			From("dbr_people", "p1").
			Join(
				MakeAlias("dbr_people", "p2"),
				Column("`p2`.`id` = `p1`.`id`"),
				Column("p1.id", argInt(42)),
			)

		sql, _, err := sqlObj.ToSQL()
		assert.NoError(t, err)
		assert.Equal(t,
			"SELECT `p1`.*, `p2`.* FROM `dbr_people` AS `p1` INNER JOIN `dbr_people` AS `p2` ON (`p2`.`id` = `p1`.`id`) AND (`p1`.`id` = ?)",
			sql,
		)
	})

	t.Run("left", func(t *testing.T) {
		sqlObj := s.
			Select("p1.*", "p2.name").
			From("dbr_people", "p1").
			LeftJoin(
				MakeAlias("dbr_people", "p2"),
				Column("`p2`.`id` = `p1`.`id`"),
				Column("p1.id", argInt(42)),
			)
		sql, _, err := sqlObj.ToSQL()
		assert.NoError(t, err)
		assert.Equal(t,
			"SELECT `p1`.*, `p2`.`name` FROM `dbr_people` AS `p1` LEFT JOIN `dbr_people` AS `p2` ON (`p2`.`id` = `p1`.`id`) AND (`p1`.`id` = ?)",
			sql,
		)
	})

	t.Run("right", func(t *testing.T) {
		sqlObj := s.
			Select("p1.*").
			AddColumnsAlias("p2.name", "p2Name", "p2.email", "p2Email", "id", "internalID").
			From("dbr_people", "p1").
			RightJoin(
				MakeAlias("dbr_people", "p2"),
				Column("`p2`.`id` = `p1`.`id`"),
			)

		sql, _, err := sqlObj.ToSQL()
		assert.NoError(t, err)
		assert.Equal(t,
			"SELECT `p1`.*, `p2`.`name` AS `p2Name`, `p2`.`email` AS `p2Email`, `id` AS `internalID` FROM `dbr_people` AS `p1` RIGHT JOIN `dbr_people` AS `p2` ON (`p2`.`id` = `p1`.`id`)",
			sql,
		)
	})

	t.Run("using", func(t *testing.T) {
		sqlObj := s.
			Select("p1.*").
			AddColumnsAlias("p2.name", "p2Name", "p2.email", "p2Email").
			From("dbr_people", "p1").
			RightJoin(
				MakeAlias("dbr_people", "p2"),
				Using("id", "email"),
			)

		sql, _, err := sqlObj.ToSQL()
		assert.NoError(t, err)
		assert.Equal(t,
			"SELECT `p1`.*, `p2`.`name` AS `p2Name`, `p2`.`email` AS `p2Email` FROM `dbr_people` AS `p1` RIGHT JOIN `dbr_people` AS `p2` USING (`id`,`email`)",
			sql,
		)
	})
}
func TestSelect_Locks(t *testing.T) {
	t.Run("LOCK IN SHARE MODE", func(t *testing.T) {
		s := NewSelect("p1.*").
			AddColumnsAlias("p2.name", "p2Name", "p2.email", "p2Email").
			From("dbr_people", "p1").LockInShareMode()
		sql, _, err := s.ToSQL()
		assert.NoError(t, err)
		assert.Equal(t,
			"SELECT `p1`.*, `p2`.`name` AS `p2Name`, `p2`.`email` AS `p2Email` FROM `dbr_people` AS `p1` LOCK IN SHARE MODE",
			sql,
		)
	})
	t.Run("FOR UPDATE", func(t *testing.T) {
		s := NewSelect("p1.*").
			AddColumnsAlias("p2.name", "p2Name", "p2.email", "p2Email").
			From("dbr_people", "p1").ForUpdate()
		sql, _, err := s.ToSQL()
		assert.NoError(t, err)
		assert.Equal(t,
			"SELECT `p1`.*, `p2`.`name` AS `p2Name`, `p2`.`email` AS `p2Email` FROM `dbr_people` AS `p1` FOR UPDATE",
			sql,
		)
	})
}

func TestSelect_Events(t *testing.T) {
	t.Parallel()

	t.Run("Stop Propagation", func(t *testing.T) {
		d := NewSelect("a", "b").From("tableA", "tA")
		d.OrderBy("col3")

		d.Log = log.BlackHole{EnableInfo: true, EnableDebug: true}
		d.Listeners.Add(
			Listen{
				Name:      "listener1",
				EventType: OnBeforeToSQL,
				SelectFunc: func(b *Select) {
					b.OrderByDesc("col1")
				},
			},
			Listen{
				Name:      "listener2",
				EventType: OnBeforeToSQL,
				SelectFunc: func(b *Select) {
					b.OrderByDesc("col2")
					b.PropagationStopped = true
				},
			},
			Listen{
				Name:      "listener3",
				EventType: OnBeforeToSQL,
				SelectFunc: func(b *Select) {
					panic("Should not get called")
				},
			},
		)
		sql, _, err := d.ToSQL()
		assert.NoError(t, err, "%+v", err)
		assert.Exactly(t, "SELECT `a`, `b` FROM `tableA` AS `tA` ORDER BY `col3`, `col1` DESC, `col2` DESC", sql)

		sql, _, err = d.ToSQL()
		assert.NoError(t, err, "%+v", err)
		assert.Exactly(t, "SELECT `a`, `b` FROM `tableA` AS `tA` ORDER BY `col3`, `col1` DESC, `col2` DESC, `col1` DESC, `col2` DESC", sql)
	})

	t.Run("Missing EventType", func(t *testing.T) {
		s := NewSelect("a", "b").From("tableA", "tA")
		s.OrderBy("col3")
		s.Listeners.Add(Listen{
			Name: "a col1",
			SelectFunc: func(s2 *Select) {
				s2.Where(Column("a", ArgFloat64(3.14159)))
				s2.OrderByDesc("col1")
			},
		})

		sql, args, err := s.ToSQL()
		assert.Empty(t, sql)
		assert.Nil(t, args)
		assert.True(t, errors.IsEmpty(err), "%+v", err)
	})

	t.Run("Should Dispatch", func(t *testing.T) {
		s := NewSelect("a", "b").From("tableA", "tA")
		s.OrderBy("col3")
		s.Listeners.Add(Listen{
			Name:      "a col1",
			Once:      true,
			EventType: OnBeforeToSQL,
			SelectFunc: func(s2 *Select) {
				s2.Where(Column("a", ArgFloat64(3.14159)))
				s2.OrderByDesc("col1")
			},
		})
		s.Listeners.Add(Listen{
			Name:      "b col2",
			EventType: OnBeforeToSQL,
			SelectFunc: func(s2 *Select) {
				s2.OrderByDesc("col2").
					Where(Column("b", ArgString("a")))
			},
		})

		sql, args, err := s.ToSQL()
		assert.NoError(t, err)
		assert.Exactly(t, []interface{}{3.14159, "a"}, args.Interfaces())
		assert.Exactly(t, "SELECT `a`, `b` FROM `tableA` AS `tA` WHERE (`a` = ?) AND (`b` = ?) ORDER BY `col3`, `col1` DESC, `col2` DESC", sql)

		sql, args, err = s.ToSQL()
		assert.NoError(t, err)
		assert.Exactly(t, []interface{}{3.14159, "a", "a"}, args.Interfaces())
		assert.Exactly(t, "SELECT `a`, `b` FROM `tableA` AS `tA` WHERE (`a` = ?) AND (`b` = ?) AND (`b` = ?) ORDER BY `col3`, `col1` DESC, `col2` DESC, `col2` DESC", sql)

		assert.Exactly(t, `a col1; b col2`, s.Listeners.String())
	})
}

func TestSelect_Columns(t *testing.T) {
	t.Parallel()

	t.Run("AddColumns, multiple args", func(t *testing.T) {
		s := NewSelect("a", "b")
		s.From("tableA", "tA")
		s.AddColumns("d,e, f", "g", "h", "i,j ,k")
		sql, _, err := s.ToSQL()
		assert.NoError(t, err, "%+v", err)
		assert.Exactly(t, "SELECT `a`, `b`, `d,e, f`, `g`, `h`, `i,j ,k` FROM `tableA` AS `tA`", sql)
	})
	t.Run("AddColumns, each column itself", func(t *testing.T) {
		s := NewSelect("a", "b")
		s.From("tableA", "tA")
		s.AddColumns("d", "e", "f")
		sql, _, err := s.ToSQL()
		assert.NoError(t, err, "%+v", err)
		assert.Exactly(t, "SELECT `a`, `b`, `d`, `e`, `f` FROM `tableA` AS `tA`", sql)
	})
	t.Run("AddColumnsAlias Expression Quoted", func(t *testing.T) {
		s := NewSelect().From("t3").
			AddColumnsAlias("x", "u", "y", "v").
			AddColumnsAlias("SUM(price)", "total_price")
		sSQL, _, err := s.ToSQL()
		assert.NoError(t, err, "%+v", err)
		assert.Exactly(t, "SELECT `x` AS `u`, `y` AS `v`, `SUM(price)` AS `total_price` FROM `t3`", sSQL)
	})
	t.Run("AddColumns+AddColumnsExprAlias", func(t *testing.T) {
		s := NewSelect().From("t3").
			AddColumns("t3.name", "sku").
			AddColumnsExprAlias("SUM(price)", "total_price")
		sSQL, _, err := s.ToSQL()
		assert.NoError(t, err, "%+v", err)
		assert.Exactly(t, "SELECT `t3`.`name`, `sku`, SUM(price) AS `total_price` FROM `t3`", sSQL)
	})

	t.Run("AddColumnsAlias multi", func(t *testing.T) {
		s := NewSelect().From("t3").
			AddColumnsAlias("t3.name", "t3Name", "t3.sku", "t3SKU")
		sSQL, _, err := s.ToSQL()
		assert.NoError(t, err, "%+v", err)
		assert.Exactly(t, "SELECT `t3`.`name` AS `t3Name`, `t3`.`sku` AS `t3SKU` FROM `t3`", sSQL)
	})
	t.Run("AddColumnsAlias one", func(t *testing.T) {
		s := NewSelect().From("t3").
			AddColumnsAlias("t3.name", "t3Name", "t3.sku", "t3SKU")
		sSQL, _, err := s.ToSQL()
		assert.NoError(t, err, "%+v", err)
		assert.Exactly(t, "SELECT `t3`.`name` AS `t3Name`, `t3`.`sku` AS `t3SKU` FROM `t3`", sSQL)
	})
	t.Run("AddColumnsAlias imbalanced", func(t *testing.T) {
		s := NewSelect().From("t3").
			AddColumnsAlias("t3.name", "t3Name", "t3.sku")
		_, _, err := s.ToSQL()
		assert.True(t, errors.IsMismatch(err), "%+v", err)
	})
	t.Run("AddColumnsExprAlias imbalanced", func(t *testing.T) {
		s := NewSelect().From("t3").
			AddColumnsExprAlias("t3.name", "t3Name", "t3.sku")
		_, _, err := s.ToSQL()
		assert.True(t, errors.IsMismatch(err), "%+v", err)
	})
	t.Run("AddColumnsExprAlias", func(t *testing.T) {
		s := NewSelect().From("sales_bestsellers_aggregated_daily", "t3").
			AddColumnsExprAlias("DATE_FORMAT(t3.period, '%Y-%m-01')", "period")
		sSQL, _, err := s.ToSQL()
		assert.NoError(t, err, "%+v", err)
		assert.Exactly(t, "SELECT DATE_FORMAT(t3.period, '%Y-%m-01') AS `period` FROM `sales_bestsellers_aggregated_daily` AS `t3`", sSQL)
	})
	t.Run("AddColumns with expression incorrect", func(t *testing.T) {
		s := NewSelect().AddColumns(" `t.value`", "`t`.`attribute_id`", "t.{column} AS `col_type`").From("catalog_product_entity_{type}", "t")
		sSQL, _, err := s.ToSQL()
		assert.NoError(t, err, "%+v", err)
		assert.Exactly(t, "SELECT ` t`.`value`, `t`.`attribute_id`, `t`.`{column} AS col_type` FROM `catalog_product_entity_{type}` AS `t`", sSQL)
	})
}

func TestSubSelect(t *testing.T) {
	t.Parallel()
	sub := NewSelect().From("catalog_category_product").
		AddColumns("entity_id").Where(Column("category_id", ArgInt64(234)))

	runner := func(op rune, wantSQL string) func(*testing.T) {
		return func(t *testing.T) {
			s := NewSelect("sku", "type_id").
				From("catalog_product_entity").
				Where(SubSelect("entity_id", op, sub))

			sStr, args, err := s.ToSQL()
			assert.NoError(t, err)
			assert.Exactly(t, []interface{}{int64(234)}, args.Interfaces())
			assert.Exactly(t, wantSQL, sStr)
		}
	}
	t.Run("IN", runner(In,
		"SELECT `sku`, `type_id` FROM `catalog_product_entity` WHERE (`entity_id` IN (SELECT `entity_id` FROM `catalog_category_product` WHERE (`category_id` = ?)))",
	))
	t.Run("EXISTS", runner(Exists,
		"SELECT `sku`, `type_id` FROM `catalog_product_entity` WHERE (`entity_id` EXISTS (SELECT `entity_id` FROM `catalog_category_product` WHERE (`category_id` = ?)))",
	))
	t.Run("NOT EXISTS", runner(NotExists,
		"SELECT `sku`, `type_id` FROM `catalog_product_entity` WHERE (`entity_id` NOT EXISTS (SELECT `entity_id` FROM `catalog_category_product` WHERE (`category_id` = ?)))",
	))
	t.Run("NOT EQUAL", runner(NotEqual,
		"SELECT `sku`, `type_id` FROM `catalog_product_entity` WHERE (`entity_id` != (SELECT `entity_id` FROM `catalog_category_product` WHERE (`category_id` = ?)))",
	))
	t.Run("NOT EQUAL", runner(Equal,
		"SELECT `sku`, `type_id` FROM `catalog_product_entity` WHERE (`entity_id` = (SELECT `entity_id` FROM `catalog_category_product` WHERE (`category_id` = ?)))",
	))
}

func TestSelect_Subselect(t *testing.T) {
	t.Parallel()
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
		sel3 := NewSelect().From("sales_bestsellers_aggregated_daily", "t3").
			AddColumnsExprAlias("DATE_FORMAT(t3.period, '%Y-%m-01')", "period").
			AddColumns("`t3`.`store_id`,`t3`.`product_id`,`t3`.`product_name`").
			AddColumnsExprAlias("AVG(`t3`.`product_price`)", "avg_price", "SUM(t3.qty_ordered)", "total_qty").
			GroupBy("t3.store_id").
			GroupByExpr("DATE_FORMAT(t3.period, '%Y-%m-01')").
			GroupBy("t3.product_id", "t3.product_name").
			OrderBy("t3.store_id").
			OrderByExpr("DATE_FORMAT(t3.period, '%Y-%m-01')").
			OrderByDesc("total_qty")

		sel2 := NewSelectFromSub(sel3, "t2").
			AddColumns("t2.period", "t2.store_id", "t2.product_id", "t2.product_name", "t2.avg_price").
			AddColumnsAlias("`t2`.`total_qty`", "`qty_ordered`")

		sel1 := NewSelectFromSub(sel2, "t1").
			AddColumns("t1.period", "t1.store_id", "t1.product_id", "t1.product_name", "t1.avg_price", "t1.qty_ordered").
			OrderBy("`t1`.period", "`t1`.product_id")

		sSQL, args, err := sel1.ToSQL()
		if err != nil {
			t.Fatalf("%+v", err)
		}
		assert.Exactly(t, []interface{}(nil), args.Interfaces())
		//println(sSQL)
		const wantSQL = "SELECT `t1`.`period`, `t1`.`store_id`, `t1`.`product_id`, `t1`.`product_name`, `t1`.`avg_price`, `t1`.`qty_ordered` FROM (SELECT `t2`.`period`, `t2`.`store_id`, `t2`.`product_id`, `t2`.`product_name`, `t2`.`avg_price`, `t2`.`total_qty` AS `qty_ordered` FROM (SELECT DATE_FORMAT(t3.period, '%Y-%m-01') AS `period`, `t3`.`store_id`,`t3`.`product_id`,`t3`.`product_name`, AVG(`t3`.`product_price`) AS `avg_price`, SUM(t3.qty_ordered) AS `total_qty` FROM `sales_bestsellers_aggregated_daily` AS `t3` GROUP BY `t3`.`store_id`, DATE_FORMAT(t3.period, '%Y-%m-01'), `t3`.`product_id`, `t3`.`product_name` ORDER BY `t3`.`store_id`, DATE_FORMAT(t3.period, '%Y-%m-01'), `total_qty` DESC) AS `t2`) AS `t1` ORDER BY `t1`.`period`, `t1`.`product_id`"
		if sSQL != wantSQL {
			t.Errorf("\nHave: %q\nWant: %q", sSQL, wantSQL)
		}
	})

	t.Run("with args", func(t *testing.T) {
		sel3 := NewSelect().From("sales_bestsellers_aggregated_daily", "t3").
			AddColumnsExprAlias("DATE_FORMAT(t3.period, '%Y-%m-01')", "period").
			AddColumns("`t3`.`store_id`,`t3`.`product_id`,`t3`.`product_name`").
			AddColumnsExprAlias("AVG(`t3`.`product_price`)", "avg_price", "SUM(t3.qty_ordered)", "total_qty").
			GroupBy("t3.store_id").
			GroupByExpr("DATE_FORMAT(t3.period, '%Y-%m-01')").
			GroupBy("t3.product_id", "t3.product_name").
			Having(Expression("COUNT(*)>?", argInt(3))).
			OrderBy("t3.store_id").
			OrderByExpr("DATE_FORMAT(t3.period, '%Y-%m-01')").
			OrderByDesc("total_qty DESC").
			Where(Column("t3.store_id", ArgInt64(2, 3, 4).Operator(In)))

		sel2 := NewSelectFromSub(sel3, "t2").
			AddColumns("t2.period", "t2.store_id", "t2.product_id", "t2.product_name", "t2.avg_price").
			AddColumnsAlias("t2.total_qty", "qty_ordered")

		sel1 := NewSelectFromSub(sel2, "t1").
			AddColumns("t1.period", "t1.store_id", "t1.product_id", "t1.product_name", "t1.avg_price", "t1.qty_ordered").
			OrderBy("`t1`.period", "`t1`.product_id")

		sSQL, args, err := sel1.ToSQL()
		if err != nil {
			t.Fatalf("%+v", err)
		}
		assert.Exactly(t, []interface{}{int64(2), int64(3), int64(4), int64(3)}, args.Interfaces())
		//println(sSQL)
		const wantSQL = "SELECT `t1`.`period`, `t1`.`store_id`, `t1`.`product_id`, `t1`.`product_name`, `t1`.`avg_price`, `t1`.`qty_ordered` FROM (SELECT `t2`.`period`, `t2`.`store_id`, `t2`.`product_id`, `t2`.`product_name`, `t2`.`avg_price`, `t2`.`total_qty` AS `qty_ordered` FROM (SELECT DATE_FORMAT(t3.period, '%Y-%m-01') AS `period`, `t3`.`store_id`,`t3`.`product_id`,`t3`.`product_name`, AVG(`t3`.`product_price`) AS `avg_price`, SUM(t3.qty_ordered) AS `total_qty` FROM `sales_bestsellers_aggregated_daily` AS `t3` WHERE (`t3`.`store_id` IN ?) GROUP BY `t3`.`store_id`, DATE_FORMAT(t3.period, '%Y-%m-01'), `t3`.`product_id`, `t3`.`product_name` HAVING (COUNT(*)>?) ORDER BY `t3`.`store_id`, DATE_FORMAT(t3.period, '%Y-%m-01'), `total_qty DESC` DESC) AS `t2`) AS `t1` ORDER BY `t1`.`period`, `t1`.`product_id`"
		if sSQL != wantSQL {
			t.Errorf("\nHave: %q\nWant: %q", sSQL, wantSQL)
		}
	})
}

func TestParenthesisOpen_Close(t *testing.T) {
	t.Parallel()
	t.Run("beginning of WHERE", func(t *testing.T) {

		sel := NewSelect("a", "b").
			From("c", "cc").
			Where(
				ParenthesisOpen(),
				Column("d", argInt(1)),
				Column("e", ArgString("wat")).Or(),
				ParenthesisClose(),
				Column("f", ArgFloat64(2.7182)),
			).
			GroupBy("ab").
			Having(
				ParenthesisOpen(),
				Column("m", argInt(33)),
				Column("n", ArgString("wh3r3")).Or(),
				ParenthesisClose(),
				Column("j = k"),
			)

		sql, args, err := sel.ToSQL()
		assert.NoError(t, err)
		assert.Equal(t, "SELECT `a`, `b` FROM `c` AS `cc` WHERE ((`d` = ?) OR (`e` = ?)) AND (`f` = ?) GROUP BY `ab` HAVING ((`m` = ?) OR (`n` = ?)) AND (`j = k`)", sql)
		assert.Equal(t, []interface{}{int64(1), "wat", 2.7182, int64(33), "wh3r3"}, args.Interfaces())
	})

	t.Run("end of WHERE", func(t *testing.T) {
		sel := NewSelect("a", "b").
			From("c", "cc").
			Where(
				Column("f", ArgFloat64(2.7182)),
				ParenthesisOpen(),
				Column("d", argInt(1)),
				Column("e", ArgString("wat")).Or(),
				ParenthesisClose(),
			).
			GroupBy("ab").
			Having(
				Column("j = k"),
				ParenthesisOpen(),
				Column("m", argInt(33)),
				Column("n", ArgString("wh3r3")).Or(),
				ParenthesisClose(),
			)

		sql, _, err := sel.ToSQL()
		assert.NoError(t, err)
		assert.Equal(t, "SELECT `a`, `b` FROM `c` AS `cc` WHERE (`f` = ?) AND ((`d` = ?) OR (`e` = ?)) GROUP BY `ab` HAVING (`j = k`) AND ((`m` = ?) OR (`n` = ?))", sql)
	})

	t.Run("middle of WHERE", func(t *testing.T) {
		sel := NewSelect("a", "b").
			From("c", "cc").
			Where(
				Column("f", ArgFloat64(2.7182)),
				ParenthesisOpen(),
				Column("d", argInt(1)),
				Column("e", ArgString("wat")).Or(),
				ParenthesisClose(),
				Column("p", ArgFloat64(3.141592)),
			).
			GroupBy("ab").
			Having(
				Column("j = k"),
				ParenthesisOpen(),
				Column("m", argInt(33)),
				Column("n", ArgString("wh3r3")).Or(),
				ParenthesisClose(),
				Column("q", ArgNull().Operator(NotNull)),
			)

		sql, _, err := sel.ToSQL()
		assert.NoError(t, err)
		assert.Equal(t, "SELECT `a`, `b` FROM `c` AS `cc` WHERE (`f` = ?) AND ((`d` = ?) OR (`e` = ?)) AND (`p` = ?) GROUP BY `ab` HAVING (`j = k`) AND ((`m` = ?) OR (`n` = ?)) AND (`q` IS NOT NULL)", sql)
	})
}

func TestSelect_Count(t *testing.T) {
	t.Run("written count star gets quoted", func(t *testing.T) {
		sqlStr, args, err := NewSelect("count(*)").From("dbr_people").ToSQL()
		assert.NoError(t, err)
		assert.Exactly(t, Arguments{}, args)
		assert.Equal(t, "SELECT `count(*)` FROM `dbr_people`", sqlStr)
	})
	t.Run("func count star", func(t *testing.T) {
		sqlStr, args, err := NewSelect().Count().From("dbr_people").ToSQL()
		assert.NoError(t, err)
		assert.Exactly(t, Arguments{}, args)
		assert.Equal(t, "SELECT COUNT(*) AS `counted` FROM `dbr_people`", sqlStr)
	})
}

func TestSelect_UseBuildCache(t *testing.T) {
	t.Parallel()

	sel := NewSelect("a", "b").
		Distinct().
		From("c", "cc").
		Where(
			ParenthesisOpen(),
			Column("d", argInt(1)),
			Column("e", ArgString("wat")).Or(),
			ParenthesisClose(),
			Eq{"f": argInt(2)}, Eq{"g": argInt(3)},
		).
		Where(Eq{"h": ArgInt64(4, 5, 6).Operator(In)}).
		GroupBy("ab").
		Having(
			ParenthesisOpen(),
			Column("m", argInt(33)),
			Column("n", ArgString("wh3r3")).Or(),
			ParenthesisClose(),
			Expression("j = k"),
		).
		OrderBy("l").
		Limit(7).
		Offset(8)
	sel.UseBuildCache = true

	const cachedSQLPlaceHolder = "SELECT DISTINCT `a`, `b` FROM `c` AS `cc` WHERE ((`d` = ?) OR (`e` = ?)) AND (`f` = ?) AND (`g` = ?) AND (`h` IN ?) GROUP BY `ab` HAVING ((`m` = ?) OR (`n` = ?)) AND (j = k) ORDER BY `l` LIMIT 7 OFFSET 8"
	t.Run("without interpolate", func(t *testing.T) {
		for i := 0; i < 3; i++ {
			sql, args, err := sel.ToSQL()
			require.NoError(t, err, "%+v", err)
			require.Equal(t, cachedSQLPlaceHolder, sql)
			assert.Equal(t, []interface{}{int64(1), "wat", int64(2), int64(3), int64(4), int64(5), int64(6), int64(33), "wh3r3"}, args.Interfaces())
			assert.Equal(t, cachedSQLPlaceHolder, string(sel.buildCache))
		}
	})

	t.Run("with interpolate", func(t *testing.T) {
		sel.Interpolate()
		sel.buildCache = nil
		sel.RawArguments = nil

		const cachedSQLInterpolated = "SELECT DISTINCT `a`, `b` FROM `c` AS `cc` WHERE ((`d` = 1) OR (`e` = 'wat')) AND (`f` = 2) AND (`g` = 3) AND (`h` IN (4,5,6)) GROUP BY `ab` HAVING ((`m` = 33) OR (`n` = 'wh3r3')) AND (j = k) ORDER BY `l` LIMIT 7 OFFSET 8"
		for i := 0; i < 3; i++ {
			sql, args, err := sel.ToSQL()
			assert.Equal(t, cachedSQLPlaceHolder, string(sel.buildCache))
			require.NoError(t, err, "%+v", err)
			require.Equal(t, cachedSQLInterpolated, sql)
			assert.Nil(t, args)
		}
	})
}
