package dbr

import (
	"context"
	"testing"

	"github.com/corestoreio/errors"
	"github.com/corestoreio/log"
	"github.com/stretchr/testify/assert"
)

func TestSelectBasicToSQL(t *testing.T) {
	s := createFakeSession()

	sel := s.Select("a", "b").From("c").Where(Condition("id = ?", ArgInt(1)))
	sql, args, err := sel.ToSQL()
	assert.NoError(t, err)
	assert.Equal(t, "SELECT a, b FROM `c` WHERE (id = ?)", sql)
	assert.Equal(t, []interface{}{int64(1)}, args.Interfaces())
}

func TestSelectFullToSQL(t *testing.T) {
	s := createFakeSession()

	sel := s.Select("a", "b").
		Distinct().
		From("c", "cc").
		Where(Condition("d = ? OR e = ?",
			ArgInt(1), ArgString("wat")),
			Eq{"f": ArgInt(2)}, Eq{"g": ArgInt(3)},
		).
		Where(Eq{"h": ArgInt64(4, 5, 6).Operator(OperatorIn)}).
		GroupBy("ab").
		Having(Condition("j = k")).
		OrderBy("l").
		Limit(7).
		Offset(8)

	sql, args, err := sel.ToSQL()
	assert.NoError(t, err)
	assert.Equal(t, "SELECT DISTINCT a, b FROM `c` AS `cc` WHERE (d = ? OR e = ?) AND (`f` = ?) AND (`g` = ?) AND (`h` IN ?) GROUP BY ab HAVING (j = k) ORDER BY l LIMIT 7 OFFSET 8", sql)
	assert.Equal(t, []interface{}{int64(1), "wat", int64(2), int64(3), int64(4), int64(5), int64(6)}, args.Interfaces())
}

func TestSelectPaginateOrderDirToSQL(t *testing.T) {
	s := createFakeSession()

	sql, args, err := s.Select("a", "b").
		From("c").
		Where(Condition("d = ?", ArgInt(1))).
		Paginate(1, 20).
		OrderDir("id", false).
		ToSQL()
	assert.NoError(t, err)
	assert.Equal(t, "SELECT a, b FROM `c` WHERE (d = ?) ORDER BY id DESC LIMIT 20 OFFSET 0", sql)
	assert.Equal(t, []interface{}{int64(1)}, args.Interfaces())

	sql, args, err = s.Select("a", "b").
		From("c").
		Where(Condition("d = ?", ArgInt(1))).
		Paginate(3, 30).
		OrderDir("id", true).
		ToSQL()
	assert.NoError(t, err)
	assert.Equal(t, "SELECT a, b FROM `c` WHERE (d = ?) ORDER BY id ASC LIMIT 30 OFFSET 60", sql)
	assert.Equal(t, []interface{}{int64(1)}, args.Interfaces())
}

func TestSelectNoWhereSQL(t *testing.T) {
	s := createFakeSession()

	sql, args, err := s.Select("a", "b").From("c").ToSQL()
	assert.NoError(t, err)
	assert.Equal(t, "SELECT a, b FROM `c`", sql)
	assert.Equal(t, []interface{}(nil), args.Interfaces())
}

func TestSelectMultiHavingSQL(t *testing.T) {
	s := createFakeSession()

	sql, args, err := s.Select("a", "b").From("c").
		Where(Condition("p = ?", ArgInt(1))).
		GroupBy("z").Having(Condition("z = ?", ArgInt(2)), Condition("y = ?", ArgInt(3))).ToSQL()
	assert.NoError(t, err)
	assert.Equal(t, "SELECT a, b FROM `c` WHERE (p = ?) GROUP BY z HAVING (z = ?) AND (y = ?)", sql)
	assert.Equal(t, []interface{}{int64(1), int64(2), int64(3)}, args.Interfaces())
}

func TestSelectMultiOrderSQL(t *testing.T) {
	s := createFakeSession()

	sql, args, err := s.Select("a", "b").From("c").OrderBy("name ASC").OrderBy("id DESC").ToSQL()
	assert.NoError(t, err)
	assert.Equal(t, "SELECT a, b FROM `c` ORDER BY name ASC, id DESC", sql)
	assert.Equal(t, []interface{}(nil), args.Interfaces())
}

func TestSelect_ConditionColumn(t *testing.T) {
	// TODO rewrite test to use every type which implements interface Argument and every operator

	s := createFakeSession()
	runner := func(arg Argument, wantSQL string, wantVal []interface{}) func(*testing.T) {
		return func(t *testing.T) {
			sql, args, err := s.Select("a", "b").From("c").Where(Condition("d", arg)).ToSQL()
			assert.NoError(t, err)
			assert.Exactly(t, wantSQL, sql)
			assert.Exactly(t, wantVal, args.Interfaces())

		}
	}
	t.Run("single int64", runner(
		ArgInt64(33),
		"SELECT a, b FROM `c` WHERE (`d` = ?)",
		[]interface{}{int64(33)},
	))
	t.Run("IN int64", runner(
		ArgInt64(33, 44).Operator(OperatorIn),
		"SELECT a, b FROM `c` WHERE (`d` IN ?)",
		[]interface{}{int64(33), int64(44)},
	))
	t.Run("single float64", runner(
		ArgFloat64(33),
		"SELECT a, b FROM `c` WHERE (`d` = ?)",
		[]interface{}{float64(33)},
	))
	t.Run("IN float64", runner(
		ArgFloat64(33, 44).Operator('i'),
		"SELECT a, b FROM `c` WHERE (`d` IN ?)",
		[]interface{}{float64(33), float64(44)},
	))
	t.Run("NOT IN float64", runner(
		ArgFloat64(33, 44).Operator('I'),
		"SELECT a, b FROM `c` WHERE (`d` NOT IN ?)",
		[]interface{}{float64(33), float64(44)},
	))
	t.Run("single int", runner(
		ArgInt(33),
		"SELECT a, b FROM `c` WHERE (`d` = ?)",
		[]interface{}{int64(33)},
	))
	t.Run("IN int", runner(
		ArgInt(33, 44).Operator(OperatorIn),
		"SELECT a, b FROM `c` WHERE (`d` IN ?)",
		[]interface{}{int64(33), int64(44)},
	))
	t.Run("single string", runner(
		ArgString("w"),
		"SELECT a, b FROM `c` WHERE (`d` = ?)",
		[]interface{}{"w"},
	))
	t.Run("IN string", runner(
		ArgString("x", "y").Operator(OperatorIn),
		"SELECT a, b FROM `c` WHERE (`d` IN ?)",
		[]interface{}{"x", "y"},
	))

	t.Run("BETWEEN int64", runner(
		ArgInt64(5, 6).Operator(OperatorBetween),
		"SELECT a, b FROM `c` WHERE (`d` BETWEEN ? AND ?)",
		[]interface{}{int64(5), int64(6)},
	))
	t.Run("NOT BETWEEN int64", runner(
		ArgInt64(5, 6).Operator(OperatorNotBetween),
		"SELECT a, b FROM `c` WHERE (`d` NOT BETWEEN ? AND ?)",
		[]interface{}{int64(5), int64(6)},
	))

	t.Run("LIKE string", runner(
		ArgString("x%").Operator(OperatorLike),
		"SELECT a, b FROM `c` WHERE (`d` LIKE ?)",
		[]interface{}{"x%"},
	))
	t.Run("NOT LIKE string", runner(
		ArgString("x%").Operator(OperatorNotLike),
		"SELECT a, b FROM `c` WHERE (`d` NOT LIKE ?)",
		[]interface{}{"x%"},
	))

}

func TestSelect_Null(t *testing.T) {
	s := createFakeSession()

	t.Run("col is null", func(t *testing.T) {
		sql, args, err := s.Select("a", "b").From("c").Where(Condition("r", ArgNull())).ToSQL()
		assert.NoError(t, err)
		assert.Exactly(t, "SELECT a, b FROM `c` WHERE (`r` IS NULL)", sql)
		assert.Exactly(t, []interface{}(nil), args.Interfaces())
	})

	t.Run("col is not null", func(t *testing.T) {
		sql, args, err := s.Select("a", "b").From("c").Where(Condition("r", ArgNotNull())).ToSQL()
		assert.NoError(t, err)
		assert.Exactly(t, "SELECT a, b FROM `c` WHERE (`r` IS NOT NULL)", sql)
		assert.Exactly(t, []interface{}(nil), args.Interfaces())
	})

	t.Run("complex", func(t *testing.T) {
		sql, args, err := s.Select("a", "b").From("c").
			Where(
				Condition("r", ArgNull()),
				Condition("d = ?", ArgInt(3)),
				Condition("ab", ArgNull()),
				Condition("w", ArgNotNull()),
			).ToSQL()
		assert.NoError(t, err)
		assert.Exactly(t, "SELECT a, b FROM `c` WHERE (`r` IS NULL) AND (d = ?) AND (`ab` IS NULL) AND (`w` IS NOT NULL)", sql)
		assert.Exactly(t, []interface{}{int64(3)}, args.Interfaces())
	})
}

func TestSelectWhereMapSQL(t *testing.T) {
	s := createFakeSession()

	t.Run("one", func(t *testing.T) {
		sql, args, err := s.Select("a").From("b").Where(Eq{"a": ArgInt(1)}).ToSQL()
		assert.NoError(t, err)
		assert.Equal(t, "SELECT a FROM `b` WHERE (`a` = ?)", sql)
		assert.Equal(t, []interface{}{int64(1)}, args.Interfaces())
	})

	t.Run("two", func(t *testing.T) {
		sql, args, err := s.Select("a").From("b").Where(Eq{"a": ArgInt(1), "b": ArgBool(true)}).ToSQL()
		assert.NoError(t, err)
		if sql == "SELECT a FROM `b` WHERE (`a` = ?) AND (`b` = ?)" {
			assert.Equal(t, []interface{}{int64(1), true}, args.Interfaces())
		} else {
			assert.Equal(t, "SELECT a FROM `b` WHERE (`b` = ?) AND (`a` = ?)", sql)
			assert.Equal(t, []interface{}{true, int64(1)}, args.Interfaces())
		}
	})

	t.Run("one nil", func(t *testing.T) {
		sql, args, err := s.Select("a").From("b").Where(Eq{"a": nil}).ToSQL()
		assert.NoError(t, err)
		assert.Equal(t, "SELECT a FROM `b` WHERE (`a` IS NULL)", sql)
		assert.Equal(t, []interface{}(nil), args.Interfaces())
	})

	t.Run("one IN", func(t *testing.T) {
		sql, args, err := s.Select("a").From("b").Where(Eq{"a": ArgInt(1, 2, 3).Operator(OperatorIn)}).ToSQL()
		assert.NoError(t, err)
		assert.Equal(t, "SELECT a FROM `b` WHERE (`a` IN ?)", sql)
		assert.Equal(t, []interface{}{int64(1), int64(2), int64(3)}, args.Interfaces())
	})

	t.Run("no values", func(t *testing.T) {
		// NOTE: a has no valid values, we want a query that returns nothing
		// TODO(CyS): revise architecture and behaviour ... maybe
		sql, args, err := s.Select("a").From("b").Where(Eq{"a": ArgInt()}).ToSQL()
		assert.NoError(t, err)
		//assert.Equal(t, "SELECT a FROM `b` WHERE (1=0)", sql)
		assert.Equal(t, "SELECT a FROM `b` WHERE (`a` = ?)", sql)
		assert.Equal(t, []interface{}{}, args.Interfaces())
	})

	t.Run("empty ArgInt", func(t *testing.T) {
		// see subtest above "no values" and its TODO
		var iVal []int
		sql, args, err := s.Select("a").From("b").Where(Eq{"a": ArgInt(iVal...)}).ToSQL()
		assert.NoError(t, err)
		assert.Equal(t, "SELECT a FROM `b` WHERE (`a` = ?)", sql)
		assert.Equal(t, []interface{}{}, args.Interfaces())
	})

	t.Run("Map nil arg", func(t *testing.T) {
		sql, args, err := s.Select("a").From("b").
			Where(Eq{"a": nil}).
			Where(Eq{"b": ArgBool(false)}).
			Where(Eq{"c": ArgNull()}).
			Where(Eq{"d": ArgNotNull()}).
			ToSQL()
		assert.NoError(t, err)
		assert.Equal(t, "SELECT a FROM `b` WHERE (`a` IS NULL) AND (`b` = ?) AND (`c` IS NULL) AND (`d` IS NOT NULL)", sql)
		assert.Equal(t, []interface{}{false}, args.Interfaces())
	})
}

func TestSelectWhereEqSQL(t *testing.T) {
	s := createFakeSession()

	sql, args, err := s.Select("a").From("b").Where(Eq{"a": ArgInt(1), "b": ArgInt64(1, 2, 3).Operator(OperatorIn)}).ToSQL()
	assert.NoError(t, err)
	if sql == "SELECT a FROM `b` WHERE (`a` = ?) AND (`b` IN ?)" {
		assert.Equal(t, []interface{}{int64(1), int64(1), int64(2), int64(3)}, args.Interfaces())
	} else {
		assert.Equal(t, sql, "SELECT a FROM `b` WHERE (`b` IN ?) AND (`a` = ?)")
		assert.Equal(t, []interface{}{int64(1), int64(2), int64(3), int64(1)}, args.Interfaces())
	}
}

func TestSelectBySQL(t *testing.T) {
	s := createFakeSession()

	sql, args, err := s.SelectBySQL("SELECT * FROM users WHERE x = 1").ToSQL()
	assert.NoError(t, err)
	assert.Equal(t, sql, "SELECT * FROM users WHERE x = 1")
	assert.Equal(t, []interface{}(nil), args.Interfaces())

	sql, args, err = s.SelectBySQL("SELECT * FROM users WHERE x = ? AND y IN ?", ArgInt(9), ArgInt(5, 6, 7)).ToSQL()
	assert.NoError(t, err)
	assert.Equal(t, sql, "SELECT * FROM users WHERE x = ? AND y IN ?")
	assert.Equal(t, []interface{}{int64(9), int64(5), int64(6), int64(7)}, args.Interfaces())

	// Doesn't fix shit if it'ab broken:
	sql, args, err = s.SelectBySQL("wat", ArgInt(9), ArgInt(5, 6, 7)).ToSQL()
	assert.NoError(t, err)
	assert.Equal(t, sql, "wat")
	assert.Equal(t, []interface{}{int64(9), int64(5), int64(6), int64(7)}, args.Interfaces())
}

func TestSelectVarieties(t *testing.T) {
	s := createFakeSession()

	sql, _, err := s.Select("id, name, email").From("users").ToSQL()
	assert.NoError(t, err)
	sql2, _, err2 := s.Select("id", "name", "email").From("users").ToSQL()
	assert.NoError(t, err2)
	assert.Equal(t, sql, sql2)
}

func TestSelectLoadStructs(t *testing.T) {
	s := createRealSessionWithFixtures()

	var people []*dbrPerson
	count, err := s.Select("id", "name", "email").From("dbr_people").OrderBy("id ASC").LoadStructs(context.TODO(), &people)

	assert.NoError(t, err)
	assert.Equal(t, count, 2)

	assert.Equal(t, len(people), 2)
	if len(people) == 2 {
		// Make sure that the Ids are set. It'ab possible (maybe?) that different DBs set ids differently so
		// don't assume they're 1 and 2.
		assert.True(t, people[0].ID > 0)
		assert.True(t, people[1].ID > people[0].ID)

		assert.Equal(t, "Jonathan", people[0].Name)
		assert.True(t, people[0].Email.Valid)
		assert.Equal(t, "jonathan@uservoice.com", people[0].Email.String)
		assert.Equal(t, "Dmitri", people[1].Name)
		assert.True(t, people[1].Email.Valid)
		assert.Equal(t, "zavorotni@jadius.com", people[1].Email.String)
	}

	// TODO: test map
}

func TestSelectLoadStruct(t *testing.T) {
	s := createRealSessionWithFixtures()

	// Found:
	var person dbrPerson
	err := s.Select("id", "name", "email").From("dbr_people").Where(Condition("email = ?", ArgString("jonathan@uservoice.com"))).LoadStruct(context.TODO(), &person)
	assert.NoError(t, err)
	assert.True(t, person.ID > 0)
	assert.Equal(t, "Jonathan", person.Name)
	assert.True(t, person.Email.Valid)
	assert.Equal(t, "jonathan@uservoice.com", person.Email.String)

	// Not found:
	var person2 dbrPerson
	err = s.Select("id", "name", "email").From("dbr_people").Where(Condition("email = ?", ArgString("dontexist@uservoice.com"))).LoadStruct(context.TODO(), &person2)
	assert.True(t, errors.IsNotFound(err), "%+v", err)
}

func TestSelectBySQLLoadStructs(t *testing.T) {
	s := createRealSessionWithFixtures()

	var people []*dbrPerson
	count, err := s.SelectBySQL("SELECT name FROM dbr_people WHERE email = ?", ArgString("jonathan@uservoice.com")).LoadStructs(context.TODO(), &people)

	assert.NoError(t, err)
	assert.Equal(t, count, 1)
	if len(people) == 1 {
		assert.Equal(t, "Jonathan", people[0].Name)
		assert.Equal(t, int64(0), people[0].ID)       // not set
		assert.Equal(t, false, people[0].Email.Valid) // not set
		assert.Equal(t, "", people[0].Email.String)   // not set
	}
}

func TestSelectLoadValue(t *testing.T) {
	s := createRealSessionWithFixtures()

	var name string
	err := s.Select("name").From("dbr_people").Where(Condition("email = 'jonathan@uservoice.com'")).LoadValue(context.TODO(), &name)

	assert.NoError(t, err)
	assert.Equal(t, "Jonathan", name)

	var id int64
	err = s.Select("id").From("dbr_people").Limit(1).LoadValue(context.TODO(), &id)

	assert.NoError(t, err)
	assert.True(t, id > 0)
}

func TestSelectLoadValues(t *testing.T) {
	s := createRealSessionWithFixtures()

	var names []string
	count, err := s.Select("name").From("dbr_people").LoadValues(context.TODO(), &names)

	assert.NoError(t, err)
	assert.Equal(t, 2, count)
	assert.Equal(t, []string{"Jonathan", "Dmitri"}, names)

	var ids []int64
	count, err = s.Select("id").From("dbr_people").Limit(1).LoadValues(context.TODO(), &ids)

	assert.NoError(t, err)
	assert.Equal(t, count, 1)
	assert.Equal(t, ids, []int64{1})
}

//func TestSelectReturn(t *testing.T) {
//	ab := createRealSessionWithFixtures()
//
//	name, err := ab.Select("name").From("dbr_people").Where(Condition("email = 'jonathan@uservoice.com'")).ReturnString()
//	assert.NoError(t, err)
//	assert.Equal(t, name, "Jonathan")
//
//	count, err := ab.Select("COUNT(*)").From("dbr_people").ReturnInt64()
//	assert.NoError(t, err)
//	assert.Equal(t, count, int64(2))
//
//	names, err := ab.Select("name").From("dbr_people").Where(Condition("email = 'jonathan@uservoice.com'")).ReturnStrings()
//	assert.NoError(t, err)
//	assert.Equal(t, names, []string{"Jonathan"})
//
//	counts, err := ab.Select("COUNT(*)").From("dbr_people").ReturnInt64s()
//	assert.NoError(t, err)
//	assert.Equal(t, counts, []int64{2})
//}

func TestSelectJoin(t *testing.T) {
	s := createRealSessionWithFixtures()

	sqlObj := s.
		Select("p1.*", "p2.*").
		From("dbr_people", "p1").
		Join(
			JoinTable("dbr_people", "p2"),
			JoinColumns(),
			Condition("`p2`.`id` = `p1`.`id`"),
			Condition("p1.id", ArgInt(42)),
		)

	sql, _, err := sqlObj.ToSQL()
	assert.NoError(t, err)
	assert.Equal(t,
		"SELECT p1.*, p2.* FROM `dbr_people` AS `p1` INNER JOIN `dbr_people` AS `p2` ON (`p2`.`id` = `p1`.`id`) AND (`p1`.`id` = ?)",
		sql,
	)

	sqlObj = s.
		Select("p1.*").
		From("dbr_people", "p1").
		LeftJoin(
			JoinTable("dbr_people", "p2"),
			JoinColumns("p2.name"),
			Condition("`p2`.`id` = `p1`.`id`"),
			Condition("p1.id", ArgInt(42)),
		)

	sql, _, err = sqlObj.ToSQL()
	assert.NoError(t, err)
	assert.Equal(t,
		"SELECT p1.*, p2.name FROM `dbr_people` AS `p1` LEFT JOIN `dbr_people` AS `p2` ON (`p2`.`id` = `p1`.`id`) AND (`p1`.`id` = ?)",
		sql,
	)

	sqlObj = s.
		Select("p1.*").
		From("dbr_people", "p1").
		RightJoin(
			JoinTable("dbr_people", "p2"),
			Quoter.ColumnAlias("p2.name", "p2Name", "p2.email", "p2Email", "id", "internalID"),
			Condition("`p2`.`id` = `p1`.`id`"),
		)

	sql, _, err = sqlObj.ToSQL()
	assert.NoError(t, err)
	assert.Equal(t,
		"SELECT p1.*, `p2`.`name` AS `p2Name`, `p2`.`email` AS `p2Email`, `id` AS `internalID` FROM `dbr_people` AS `p1` RIGHT JOIN `dbr_people` AS `p2` ON (`p2`.`id` = `p1`.`id`)",
		sql,
	)
}

func TestSelect_Join(t *testing.T) {
	t.Parallel()
	const want = "SELECT IFNULL(`manufacturerStore`.`value`,IFNULL(`manufacturerGroup`.`value`,IFNULL(`manufacturerWebsite`.`value`,IFNULL(`manufacturerDefault`.`value`,'')))) AS `manufacturer`, cpe.* FROM `catalog_product_entity` AS `cpe` LEFT JOIN `catalog_product_entity_varchar` AS `manufacturerDefault` ON (manufacturerDefault.scope = 0) AND (manufacturerDefault.scope_id = 0) AND (manufacturerDefault.attribute_id = 83) AND (manufacturerDefault.value IS NOT NULL) LEFT JOIN `catalog_product_entity_varchar` AS `manufacturerWebsite` ON (manufacturerWebsite.scope = 1) AND (manufacturerWebsite.scope_id = 10) AND (manufacturerWebsite.attribute_id = 83) AND (manufacturerWebsite.value IS NOT NULL) LEFT JOIN `catalog_product_entity_varchar` AS `manufacturerGroup` ON (manufacturerGroup.scope = 2) AND (manufacturerGroup.scope_id = 20) AND (manufacturerGroup.attribute_id = 83) AND (manufacturerGroup.value IS NOT NULL) LEFT JOIN `catalog_product_entity_varchar` AS `manufacturerStore` ON (manufacturerStore.scope = 2) AND (manufacturerStore.scope_id = 20) AND (manufacturerStore.attribute_id = 83) AND (manufacturerStore.value IS NOT NULL)"

	s := NewSelect("catalog_product_entity", "cpe").
		LeftJoin(
			JoinTable("catalog_product_entity_varchar", "manufacturerDefault"),
			JoinColumns("cpe.*"),
			Condition("manufacturerDefault.scope = 0"),
			Condition("manufacturerDefault.scope_id = 0"),
			Condition("manufacturerDefault.attribute_id = 83"),
			Condition("manufacturerDefault.value IS NOT NULL"),
		).
		LeftJoin(
			JoinTable("catalog_product_entity_varchar", "manufacturerWebsite"),
			JoinColumns(),
			Condition("manufacturerWebsite.scope = 1"),
			Condition("manufacturerWebsite.scope_id = 10"),
			Condition("manufacturerWebsite.attribute_id = 83"),
			Condition("manufacturerWebsite.value IS NOT NULL"),
		).
		LeftJoin(
			JoinTable("catalog_product_entity_varchar", "manufacturerGroup"),
			JoinColumns(),
			Condition("manufacturerGroup.scope = 2"),
			Condition("manufacturerGroup.scope_id = 20"),
			Condition("manufacturerGroup.attribute_id = 83"),
			Condition("manufacturerGroup.value IS NOT NULL"),
		).
		LeftJoin(
			JoinTable("catalog_product_entity_varchar", "manufacturerStore"),
			JoinColumns(),
			Condition("manufacturerStore.scope = 2"),
			Condition("manufacturerStore.scope_id = 20"),
			Condition("manufacturerStore.attribute_id = 83"),
			Condition("manufacturerStore.value IS NOT NULL"),
		)

	s.Columns = []string{EAVIfNull("manufacturer", "value", "''")}

	sql, _, err := s.ToSQL()
	assert.NoError(t, err)
	assert.Equal(t,
		want,
		sql,
	)
}

func TestSelect_Events(t *testing.T) {
	t.Parallel()

	t.Run("Stop Propagation", func(t *testing.T) {
		d := NewSelect("tableA", "tA")
		d.Columns = []string{"a", "b"}
		d.OrderBy("col3")

		d.Log = log.BlackHole{EnableInfo: true, EnableDebug: true}
		d.Listeners.Add(
			Listen{
				Name:      "listener1",
				EventType: OnBeforeToSQL,
				SelectFunc: func(b *Select) {
					b.OrderDir("col1", false)
				},
			},
			Listen{
				Name:      "listener2",
				EventType: OnBeforeToSQL,
				SelectFunc: func(b *Select) {
					b.OrderDir("col2", false)
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
		assert.Exactly(t, "SELECT a, b FROM `tableA` AS `tA` ORDER BY col3, col1 DESC, col2 DESC", sql)

		sql, _, err = d.ToSQL()
		assert.NoError(t, err, "%+v", err)
		assert.Exactly(t, "SELECT a, b FROM `tableA` AS `tA` ORDER BY col3, col1 DESC, col2 DESC, col1 DESC, col2 DESC", sql)
	})

	t.Run("Missing EventType", func(t *testing.T) {
		s := NewSelect("tableA", "tA")

		s.Columns = []string{"a", "b"}
		s.OrderBy("col3")
		s.Listeners.Add(Listen{
			Name: "a col1",
			SelectFunc: func(s2 *Select) {
				s2.Where(Condition("a=?", ArgFloat64(3.14159)))
				s2.OrderDir("col1", false)
			},
		})

		sql, args, err := s.ToSQL()
		assert.Empty(t, sql)
		assert.Nil(t, args)
		assert.True(t, errors.IsEmpty(err), "%+v", err)
	})

	t.Run("Should Dispatch", func(t *testing.T) {
		s := NewSelect("tableA", "tA")

		s.Columns = []string{"a", "b"}
		s.OrderBy("col3")
		s.Listeners.Add(Listen{
			Name:      "a col1",
			Once:      true,
			EventType: OnBeforeToSQL,
			SelectFunc: func(s2 *Select) {
				s2.Where(Condition("a=?", ArgFloat64(3.14159)))
				s2.OrderDir("col1", false)
			},
		})
		s.Listeners.Add(Listen{
			Name:      "b col2",
			EventType: OnBeforeToSQL,
			SelectFunc: func(s2 *Select) {
				s2.OrderDir("col2", false)
				s2.Where(Condition("b=?", ArgString("a")))
			},
		})

		sql, args, err := s.ToSQL()
		assert.NoError(t, err)
		assert.Exactly(t, []interface{}{3.14159, "a"}, args.Interfaces())
		assert.Exactly(t, "SELECT a, b FROM `tableA` AS `tA` WHERE (a=?) AND (b=?) ORDER BY col3, col1 DESC, col2 DESC", sql)

		sql, args, err = s.ToSQL()
		assert.NoError(t, err)
		assert.Exactly(t, []interface{}{3.14159, "a", "a"}, args.Interfaces())
		assert.Exactly(t, "SELECT a, b FROM `tableA` AS `tA` WHERE (a=?) AND (b=?) AND (b=?) ORDER BY col3, col1 DESC, col2 DESC, col2 DESC", sql)

		assert.Exactly(t, `a col1; b col2`, s.Listeners.String())
	})
}

func TestSelect_AddColumns(t *testing.T) {
	t.Parallel()
	s := NewSelect("tableA", "tA")
	s.AddColumns("a", "b", "c")
	s.AddColumns("d,e,f", "should not get added!")
	s.AddColumnsAliases("x", "u", "y", "v")
	sql, _, err := s.ToSQL()
	assert.NoError(t, err, "%+v", err)
	assert.Exactly(t, "SELECT a, b, c, d, e, f, x AS `u`, y AS `v` FROM `tableA` AS `tA`", sql)
}

func TestSelect_AddColumnsQuoted(t *testing.T) {
	s := NewSelect("t3").
		AddColumnsQuoted("t3.name", "sku").
		AddColumnsAliases("SUM(price)", "total_price")

	sSQL, _, err := s.ToSQL()
	assert.NoError(t, err, "%+v", err)
	assert.Exactly(t, "SELECT `t3`.`name`, `sku`, SUM(price) AS `total_price` FROM `t3`", sSQL)
}

func TestSelect_Subselect(t *testing.T) {
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
		sel3 := NewSelect("sales_bestsellers_aggregated_daily", "t3").
			AddColumnsAliases("DATE_FORMAT(t3.period, '%Y-%m-01')", "period").
			AddColumns("`t3`.`store_id`,`t3`.`product_id`,`t3`.`product_name`").
			AddColumnsAliases("AVG(`t3`.`product_price`)", "avg_price", "SUM(t3.qty_ordered)", "total_qty").
			GroupBy("`t3`.`store_id`", "DATE_FORMAT(t3.period, '%Y-%m-01')", "`t3`.`product_id`", "`t3`.`product_name`").
			OrderBy("`t3`.`store_id`", "DATE_FORMAT(t3.period, '%Y-%m-01')", "`total_qty` DESC")

		sel2 := NewSelectFromSub(sel3, "t2").
			AddColumns("`t2`.`period`,`t2`.`store_id`,`t2`.`product_id`,`t2`.`product_name`,`t2`.`avg_price`").
			AddColumnsAliases("`t2`.`total_qty`", "`qty_ordered`")

		sel1 := NewSelectFromSub(sel2, "t1").
			AddColumns("`t1`.`period`,`t1`.`store_id`,`t1`.`product_id`,`t1`.`product_name`,`t1`.`avg_price`,`t1`.`qty_ordered`").
			OrderBy("`t1`.period", "`t1`.product_id")

		sSQL, args, err := sel1.ToSQL()
		if err != nil {
			t.Fatalf("%+v", err)
		}
		assert.Exactly(t, []interface{}(nil), args.Interfaces())
		//println(sSQL)
		const wantSQL = "SELECT `t1`.`period`, `t1`.`store_id`, `t1`.`product_id`, `t1`.`product_name`, `t1`.`avg_price`, `t1`.`qty_ordered` FROM (SELECT `t2`.`period`, `t2`.`store_id`, `t2`.`product_id`, `t2`.`product_name`, `t2`.`avg_price`, `t2`.`total_qty` AS `qty_ordered` FROM (SELECT DATE_FORMAT(t3.period, '%Y-%m-01') AS `period`, `t3`.`store_id`, `t3`.`product_id`, `t3`.`product_name`, AVG(`t3`.`product_price`) AS `avg_price`, SUM(t3.qty_ordered) AS `total_qty` FROM `sales_bestsellers_aggregated_daily` AS `t3` GROUP BY `t3`.`store_id`, DATE_FORMAT(t3.period, '%Y-%m-01'), `t3`.`product_id`, `t3`.`product_name` ORDER BY `t3`.`store_id`, DATE_FORMAT(t3.period, '%Y-%m-01'), `total_qty` DESC) AS `t2`) AS `t1` ORDER BY `t1`.period, `t1`.product_id"
		if sSQL != wantSQL {
			t.Errorf("\nHave: %q\nWant: %q", sSQL, wantSQL)
		}
	})

	t.Run("with args", func(t *testing.T) {
		sel3 := NewSelect("sales_bestsellers_aggregated_daily", "t3").
			AddColumnsAliases("DATE_FORMAT(t3.period, '%Y-%m-01')", "period").
			AddColumns("`t3`.`store_id`,`t3`.`product_id`,`t3`.`product_name`").
			AddColumnsAliases("AVG(`t3`.`product_price`)", "avg_price", "SUM(t3.qty_ordered)", "total_qty").
			GroupBy("`t3`.`store_id`", "DATE_FORMAT(t3.period, '%Y-%m-01')", "`t3`.`product_id`", "`t3`.`product_name`").
			Having(Condition("COUNT(*)>?", ArgInt(3))).
			OrderBy("`t3`.`store_id`", "DATE_FORMAT(t3.period, '%Y-%m-01')", "`total_qty` DESC").
			Where(Condition("t3.store_id", ArgInt64(2, 3, 4).Operator(OperatorIn)))

		sel2 := NewSelectFromSub(sel3, "t2").
			AddColumns("`t2`.`period`,`t2`.`store_id`,`t2`.`product_id`,`t2`.`product_name`,`t2`.`avg_price`").
			AddColumnsAliases("`t2`.`total_qty`", "`qty_ordered`")

		sel1 := NewSelectFromSub(sel2, "t1").
			AddColumns("`t1`.`period`,`t1`.`store_id`,`t1`.`product_id`,`t1`.`product_name`,`t1`.`avg_price`,`t1`.`qty_ordered`").
			OrderBy("`t1`.period", "`t1`.product_id")

		sSQL, args, err := sel1.ToSQL()
		if err != nil {
			t.Fatalf("%+v", err)
		}
		assert.Exactly(t, []interface{}(nil), args.Interfaces())
		//println(sSQL)
		const wantSQL = "SELECT `t1`.`period`, `t1`.`store_id`, `t1`.`product_id`, `t1`.`product_name`, `t1`.`avg_price`, `t1`.`qty_ordered` FROM (SELECT `t2`.`period`, `t2`.`store_id`, `t2`.`product_id`, `t2`.`product_name`, `t2`.`avg_price`, `t2`.`total_qty` AS `qty_ordered` FROM (SELECT DATE_FORMAT(t3.period, '%Y-%m-01') AS `period`, `t3`.`store_id`, `t3`.`product_id`, `t3`.`product_name`, AVG(`t3`.`product_price`) AS `avg_price`, SUM(t3.qty_ordered) AS `total_qty` FROM `sales_bestsellers_aggregated_daily` AS `t3` GROUP BY `t3`.`store_id`, DATE_FORMAT(t3.period, '%Y-%m-01'), `t3`.`product_id`, `t3`.`product_name` ORDER BY `t3`.`store_id`, DATE_FORMAT(t3.period, '%Y-%m-01'), `total_qty` DESC) AS `t2`) AS `t1` ORDER BY `t1`.period, `t1`.product_id"
		if sSQL != wantSQL {
			t.Errorf("\nHave: %q\nWant: %q", sSQL, wantSQL)
		}
	})

}
