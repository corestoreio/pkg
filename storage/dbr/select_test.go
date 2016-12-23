package dbr

import (
	"testing"

	"github.com/corestoreio/csfw/log"
	"github.com/corestoreio/errors"
	"github.com/stretchr/testify/assert"
)

func BenchmarkSelectBasicSQL(b *testing.B) {
	s := createFakeSession()

	// Do some allocations outside the loop so they don't affect the results
	argEq := ConditionMap(Eq{"a": []int{1, 2, 3}})

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		s.Select("something_id", "user_id", "other").
			From("some_table").
			Where(ConditionRaw("d = ? OR e = ?", 1, "wat")).
			Where(argEq).
			OrderDir("id", false).
			Paginate(1, 20).
			ToSQL()
	}
}

func BenchmarkSelectFullSQL(b *testing.B) {
	s := createFakeSession()

	// Do some allocations outside the loop so they don't affect the results
	argEq1 := ConditionMap(Eq{"f": 2, "x": "hi"})
	argEq2 := ConditionMap(Eq{"g": 3})
	argEq3 := ConditionMap(Eq{"h": []int{1, 2, 3}})

	b.ResetTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		s.Select("a", "b", "z", "y", "x").
			Distinct().
			From("c").
			Where(ConditionRaw("d = ? OR e = ?", 1, "wat")).
			Where(argEq1).
			Where(argEq2).
			Where(argEq3).
			GroupBy("i").
			GroupBy("ii").
			GroupBy("iii").
			Having(ConditionRaw("j = k"), ConditionRaw("jj = ?", 1)).
			Having(ConditionRaw("jjj = ?", 2)).
			OrderBy("l").
			OrderBy("l").
			OrderBy("l").
			Limit(7).
			Offset(8).
			ToSQL()
	}
}

func TestSelectBasicToSQL(t *testing.T) {
	s := createFakeSession()
	sel := s.Select("a", "b").From("c").Where(ConditionRaw("id = ?", 1))
	for i := 0; i < 3; i++ {
		sql, args, err := sel.ToSQL()
		assert.NoError(t, err)
		assert.Equal(t, "SELECT a, b FROM `c` WHERE (id = ?)", sql, "Loop %d", 0)
		assert.Equal(t, []interface{}{1}, args, "Loop %d", 0)
	}
}

func TestSelectFullToSQL(t *testing.T) {
	s := createFakeSession()

	sel := s.Select("a", "b").
		Distinct().
		From("c", "cc").
		Where(ConditionRaw("d = ? OR e = ?", 1, "wat"), ConditionMap(Eq{"f": 2}), ConditionMap(Eq{"g": 3})).
		Where(ConditionMap(Eq{"h": []int{4, 5, 6}})).
		GroupBy("i").
		Having(ConditionRaw("j = k")).
		OrderBy("l").
		Limit(7).
		Offset(8)

	sql, args, err := sel.ToSQL()
	assert.NoError(t, err)
	assert.Equal(t, "SELECT DISTINCT a, b FROM `c` AS `cc` WHERE (d = ? OR e = ?) AND (`f` = ?) AND (`g` = ?) AND (`h` IN ?) GROUP BY i HAVING (j = k) ORDER BY l LIMIT 7 OFFSET 8", sql)
	assert.Equal(t, []interface{}{1, "wat", 2, 3, []int{4, 5, 6}}, args)

}

func TestSelectPaginateOrderDirToSQL(t *testing.T) {
	s := createFakeSession()

	sql, args, err := s.Select("a", "b").
		From("c").
		Where(ConditionRaw("d = ?", 1)).
		Paginate(1, 20).
		OrderDir("id", false).
		ToSQL()
	assert.NoError(t, err)
	assert.Equal(t, "SELECT a, b FROM `c` WHERE (d = ?) ORDER BY id DESC LIMIT 20 OFFSET 0", sql)
	assert.Equal(t, []interface{}{1}, args)

	sql, args, err = s.Select("a", "b").
		From("c").
		Where(ConditionRaw("d = ?", 1)).
		Paginate(3, 30).
		OrderDir("id", true).
		ToSQL()
	assert.NoError(t, err)
	assert.Equal(t, "SELECT a, b FROM `c` WHERE (d = ?) ORDER BY id ASC LIMIT 30 OFFSET 60", sql)
	assert.Equal(t, []interface{}{1}, args)
}

func TestSelectNoWhereSQL(t *testing.T) {
	s := createFakeSession()

	sql, args, err := s.Select("a", "b").From("c").ToSQL()
	assert.NoError(t, err)
	assert.Equal(t, sql, "SELECT a, b FROM `c`")
	assert.Equal(t, args, []interface{}(nil))
}

func TestSelectMultiHavingSQL(t *testing.T) {
	s := createFakeSession()

	sql, args, err := s.Select("a", "b").From("c").Where(ConditionRaw("p = ?", 1)).GroupBy("z").Having(ConditionRaw("z = ?", 2), ConditionRaw("y = ?", 3)).ToSQL()
	assert.NoError(t, err)
	assert.Equal(t, sql, "SELECT a, b FROM `c` WHERE (p = ?) GROUP BY z HAVING (z = ?) AND (y = ?)")
	assert.Equal(t, args, []interface{}{1, 2, 3})
}

func TestSelectMultiOrderSQL(t *testing.T) {
	s := createFakeSession()

	sql, args, err := s.Select("a", "b").From("c").OrderBy("name ASC").OrderBy("id DESC").ToSQL()
	assert.NoError(t, err)
	assert.Equal(t, sql, "SELECT a, b FROM `c` ORDER BY name ASC, id DESC")
	assert.Equal(t, args, []interface{}(nil))
}

func TestSelectWhereMapSQL(t *testing.T) {
	s := createFakeSession()

	sql, args, err := s.Select("a").From("b").Where(ConditionMap(Eq{"a": 1})).ToSQL()
	assert.NoError(t, err)
	assert.Equal(t, sql, "SELECT a FROM `b` WHERE (`a` = ?)")
	assert.Equal(t, args, []interface{}{1})

	sql, args, err = s.Select("a").From("b").Where(ConditionMap(Eq{"a": 1, "b": true})).ToSQL()
	assert.NoError(t, err)
	if sql == "SELECT a FROM `b` WHERE (`a` = ?) AND (`b` = ?)" {
		assert.Equal(t, args, []interface{}{1, true})
	} else {
		assert.Equal(t, sql, "SELECT a FROM `b` WHERE (`b` = ?) AND (`a` = ?)")
		assert.Equal(t, args, []interface{}{true, 1})
	}

	sql, args, err = s.Select("a").From("b").Where(ConditionMap(Eq{"a": nil})).ToSQL()
	assert.NoError(t, err)
	assert.Equal(t, sql, "SELECT a FROM `b` WHERE (`a` IS NULL)")
	assert.Equal(t, args, []interface{}(nil))

	sql, args, err = s.Select("a").From("b").Where(ConditionMap(Eq{"a": []int{1, 2, 3}})).ToSQL()
	assert.NoError(t, err)
	assert.Equal(t, sql, "SELECT a FROM `b` WHERE (`a` IN ?)")
	assert.Equal(t, args, []interface{}{[]int{1, 2, 3}})

	sql, args, err = s.Select("a").From("b").Where(ConditionMap(Eq{"a": []int{1}})).ToSQL()
	assert.NoError(t, err)
	assert.Equal(t, sql, "SELECT a FROM `b` WHERE (`a` = ?)")
	assert.Equal(t, args, []interface{}{1})

	// NOTE: a has no valid values, we want a query that returns nothing
	sql, args, err = s.Select("a").From("b").Where(ConditionMap(Eq{"a": []int{}})).ToSQL()
	assert.NoError(t, err)
	assert.Equal(t, sql, "SELECT a FROM `b` WHERE (1=0)")
	assert.Equal(t, args, []interface{}(nil))

	var aval []int
	sql, args, err = s.Select("a").From("b").Where(ConditionMap(Eq{"a": aval})).ToSQL()
	assert.NoError(t, err)
	assert.Equal(t, sql, "SELECT a FROM `b` WHERE (`a` IS NULL)")
	assert.Equal(t, args, []interface{}(nil))

	sql, args, err = s.Select("a").From("b").
		Where(ConditionMap(Eq{"a": []int(nil)})).
		Where(ConditionMap(Eq{"b": false})).
		ToSQL()
	assert.NoError(t, err)
	assert.Equal(t, sql, "SELECT a FROM `b` WHERE (`a` IS NULL) AND (`b` = ?)")
	assert.Equal(t, args, []interface{}{false})
}

func TestSelectWhereEqSQL(t *testing.T) {
	s := createFakeSession()

	sql, args, err := s.Select("a").From("b").Where(ConditionMap(Eq{"a": 1, "b": []int64{1, 2, 3}})).ToSQL()
	assert.NoError(t, err)
	if sql == "SELECT a FROM `b` WHERE (`a` = ?) AND (`b` IN ?)" {
		assert.Equal(t, args, []interface{}{1, []int64{1, 2, 3}})
	} else {
		assert.Equal(t, sql, "SELECT a FROM `b` WHERE (`b` IN ?) AND (`a` = ?)")
		assert.Equal(t, args, []interface{}{[]int64{1, 2, 3}, 1})
	}
}

func TestSelectBySQL(t *testing.T) {
	s := createFakeSession()

	sql, args, err := s.SelectBySQL("SELECT * FROM users WHERE x = 1").ToSQL()
	assert.NoError(t, err)
	assert.Equal(t, sql, "SELECT * FROM users WHERE x = 1")
	assert.Equal(t, args, []interface{}(nil))

	sql, args, err = s.SelectBySQL("SELECT * FROM users WHERE x = ? AND y IN ?", 9, []int{5, 6, 7}).ToSQL()
	assert.NoError(t, err)
	assert.Equal(t, sql, "SELECT * FROM users WHERE x = ? AND y IN ?")
	assert.Equal(t, args, []interface{}{9, []int{5, 6, 7}})

	// Doesn't fix shit if it's broken:
	sql, args, err = s.SelectBySQL("wat", 9, []int{5, 6, 7}).ToSQL()
	assert.NoError(t, err)
	assert.Equal(t, sql, "wat")
	assert.Equal(t, args, []interface{}{9, []int{5, 6, 7}})
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
	count, err := s.Select("id", "name", "email").From("dbr_people").OrderBy("id ASC").LoadStructs(&people)

	assert.NoError(t, err)
	assert.Equal(t, count, 2)

	assert.Equal(t, len(people), 2)
	if len(people) == 2 {
		// Make sure that the Ids are set. It's possible (maybe?) that different DBs set ids differently so
		// don't assume they're 1 and 2.
		assert.True(t, people[0].ID > 0)
		assert.True(t, people[1].ID > people[0].ID)

		assert.Equal(t, people[0].Name, "Jonathan")
		assert.True(t, people[0].Email.Valid)
		assert.Equal(t, people[0].Email.String, "jonathan@uservoice.com")
		assert.Equal(t, people[1].Name, "Dmitri")
		assert.True(t, people[1].Email.Valid)
		assert.Equal(t, people[1].Email.String, "zavorotni@jadius.com")
	}

	// TODO: test map
}

func TestSelectLoadStruct(t *testing.T) {
	s := createRealSessionWithFixtures()

	// Found:
	var person dbrPerson
	err := s.Select("id", "name", "email").From("dbr_people").Where(ConditionRaw("email = ?", "jonathan@uservoice.com")).LoadStruct(&person)
	assert.NoError(t, err)
	assert.True(t, person.ID > 0)
	assert.Equal(t, person.Name, "Jonathan")
	assert.True(t, person.Email.Valid)
	assert.Equal(t, person.Email.String, "jonathan@uservoice.com")

	// Not found:
	var person2 dbrPerson
	err = s.Select("id", "name", "email").From("dbr_people").Where(ConditionRaw("email = ?", "dontexist@uservoice.com")).LoadStruct(&person2)
	assert.True(t, errors.IsNotFound(err), "%+v", err)
}

func TestSelectBySQLLoadStructs(t *testing.T) {
	s := createRealSessionWithFixtures()

	var people []*dbrPerson
	count, err := s.SelectBySQL("SELECT name FROM dbr_people WHERE email IN ?", []string{"jonathan@uservoice.com"}).LoadStructs(&people)

	assert.NoError(t, err)
	assert.Equal(t, count, 1)
	if len(people) == 1 {
		assert.Equal(t, people[0].Name, "Jonathan")
		assert.Equal(t, people[0].ID, int64(0))       // not set
		assert.Equal(t, people[0].Email.Valid, false) // not set
		assert.Equal(t, people[0].Email.String, "")   // not set
	}
}

func TestSelectLoadValue(t *testing.T) {
	s := createRealSessionWithFixtures()

	var name string
	err := s.Select("name").From("dbr_people").Where(ConditionRaw("email = 'jonathan@uservoice.com'")).LoadValue(&name)

	assert.NoError(t, err)
	assert.Equal(t, name, "Jonathan")

	var id int64
	err = s.Select("id").From("dbr_people").Limit(1).LoadValue(&id)

	assert.NoError(t, err)
	assert.True(t, id > 0)
}

func TestSelectLoadValues(t *testing.T) {
	s := createRealSessionWithFixtures()

	var names []string
	count, err := s.Select("name").From("dbr_people").LoadValues(&names)

	assert.NoError(t, err)
	assert.Equal(t, count, 2)
	assert.Equal(t, names, []string{"Jonathan", "Dmitri"})

	var ids []int64
	count, err = s.Select("id").From("dbr_people").Limit(1).LoadValues(&ids)

	assert.NoError(t, err)
	assert.Equal(t, count, 1)
	assert.Equal(t, ids, []int64{1})
}

//func TestSelectReturn(t *testing.T) {
//	s := createRealSessionWithFixtures()
//
//	name, err := s.Select("name").From("dbr_people").Where(ConditionRaw("email = 'jonathan@uservoice.com'")).ReturnString()
//	assert.NoError(t, err)
//	assert.Equal(t, name, "Jonathan")
//
//	count, err := s.Select("COUNT(*)").From("dbr_people").ReturnInt64()
//	assert.NoError(t, err)
//	assert.Equal(t, count, int64(2))
//
//	names, err := s.Select("name").From("dbr_people").Where(ConditionRaw("email = 'jonathan@uservoice.com'")).ReturnStrings()
//	assert.NoError(t, err)
//	assert.Equal(t, names, []string{"Jonathan"})
//
//	counts, err := s.Select("COUNT(*)").From("dbr_people").ReturnInt64s()
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
			ConditionRaw("`p2`.`id` = `p1`.`id`"),
			ConditionRaw("`p1`.`id` = ?", 42),
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
			ConditionRaw("`p2`.`id` = `p1`.`id`"),
			ConditionRaw("`p1`.`id` = ?", 42),
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
			ConditionRaw("`p2`.`id` = `p1`.`id`"),
		)

	sql, _, err = sqlObj.ToSQL()
	assert.NoError(t, err)
	assert.Equal(t,
		"SELECT p1.*, `p2`.`name` AS `p2Name`, `p2`.`email` AS `p2Email`, `id` AS `internalID` FROM `dbr_people` AS `p1` RIGHT JOIN `dbr_people` AS `p2` ON (`p2`.`id` = `p1`.`id`)",
		sql,
	)
}

func TestSelect_Join(t *testing.T) {

	const want = "SELECT IFNULL(`manufacturerStore`.`value`,IFNULL(`manufacturerGroup`.`value`,IFNULL(`manufacturerWebsite`.`value`,IFNULL(`manufacturerDefault`.`value`,'')))) AS `manufacturer`, cpe.* FROM `catalog_product_entity` AS `cpe` LEFT JOIN `catalog_product_entity_varchar` AS `manufacturerDefault` ON (manufacturerDefault.scope = 0) AND (manufacturerDefault.scope_id = 0) AND (manufacturerDefault.attribute_id = 83) AND (manufacturerDefault.value IS NOT NULL) LEFT JOIN `catalog_product_entity_varchar` AS `manufacturerWebsite` ON (manufacturerWebsite.scope = 1) AND (manufacturerWebsite.scope_id = 10) AND (manufacturerWebsite.attribute_id = 83) AND (manufacturerWebsite.value IS NOT NULL) LEFT JOIN `catalog_product_entity_varchar` AS `manufacturerGroup` ON (manufacturerGroup.scope = 2) AND (manufacturerGroup.scope_id = 20) AND (manufacturerGroup.attribute_id = 83) AND (manufacturerGroup.value IS NOT NULL) LEFT JOIN `catalog_product_entity_varchar` AS `manufacturerStore` ON (manufacturerStore.scope = 2) AND (manufacturerStore.scope_id = 20) AND (manufacturerStore.attribute_id = 83) AND (manufacturerStore.value IS NOT NULL)"

	s := NewSelect("catalog_product_entity", "cpe").
		LeftJoin(
			JoinTable("catalog_product_entity_varchar", "manufacturerDefault"),
			JoinColumns("cpe.*"),
			ConditionRaw("manufacturerDefault.scope = 0"),
			ConditionRaw("manufacturerDefault.scope_id = 0"),
			ConditionRaw("manufacturerDefault.attribute_id = 83"),
			ConditionRaw("manufacturerDefault.value IS NOT NULL"),
		).
		LeftJoin(
			JoinTable("catalog_product_entity_varchar", "manufacturerWebsite"),
			JoinColumns(),
			ConditionRaw("manufacturerWebsite.scope = 1"),
			ConditionRaw("manufacturerWebsite.scope_id = 10"),
			ConditionRaw("manufacturerWebsite.attribute_id = 83"),
			ConditionRaw("manufacturerWebsite.value IS NOT NULL"),
		).
		LeftJoin(
			JoinTable("catalog_product_entity_varchar", "manufacturerGroup"),
			JoinColumns(),
			ConditionRaw("manufacturerGroup.scope = 2"),
			ConditionRaw("manufacturerGroup.scope_id = 20"),
			ConditionRaw("manufacturerGroup.attribute_id = 83"),
			ConditionRaw("manufacturerGroup.value IS NOT NULL"),
		).
		LeftJoin(
			JoinTable("catalog_product_entity_varchar", "manufacturerStore"),
			JoinColumns(),
			ConditionRaw("manufacturerStore.scope = 2"),
			ConditionRaw("manufacturerStore.scope_id = 20"),
			ConditionRaw("manufacturerStore.attribute_id = 83"),
			ConditionRaw("manufacturerStore.value IS NOT NULL"),
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

		d.Logger = log.BlackHole{EnableInfo: true, EnableDebug: true}
		d.SelectListeners.Add(
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
		s.SelectListeners.Add(Listen{
			Name: "a col1",
			SelectFunc: func(s2 *Select) {
				s2.Where(ConditionRaw("a=?", 3.14159))
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
		s.SelectListeners.Add(Listen{
			Name:      "a col1",
			Once:      true,
			EventType: OnBeforeToSQL,
			SelectFunc: func(s2 *Select) {
				s2.Where(ConditionRaw("a=?", 3.14159))
				s2.OrderDir("col1", false)
			},
		})
		s.SelectListeners.Add(Listen{
			Name:      "b col2",
			EventType: OnBeforeToSQL,
			SelectFunc: func(s2 *Select) {
				s2.OrderDir("col2", false)
				s2.Where(ConditionRaw("b=?", "a"))
			},
		})

		sql, args, err := s.ToSQL()
		assert.NoError(t, err)
		assert.Exactly(t, []interface{}{3.14159, "a"}, args)
		assert.Exactly(t, "SELECT a, b FROM `tableA` AS `tA` WHERE (a=?) AND (b=?) ORDER BY col3, col1 DESC, col2 DESC", sql)

		sql, args, err = s.ToSQL()
		assert.NoError(t, err)
		assert.Exactly(t, []interface{}{3.14159, "a", "a"}, args)
		assert.Exactly(t, "SELECT a, b FROM `tableA` AS `tA` WHERE (a=?) AND (b=?) AND (b=?) ORDER BY col3, col1 DESC, col2 DESC, col2 DESC", sql)

		assert.Exactly(t, `a col1; b col2`, s.SelectListeners.String())
	})

}
