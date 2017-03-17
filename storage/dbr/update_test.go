package dbr

import (
	"testing"

	"github.com/corestoreio/errors"
	"github.com/corestoreio/log"
	"github.com/stretchr/testify/assert"
)

func BenchmarkUpdateValuesSQL(b *testing.B) {
	s := createFakeSession()

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		s.Update("alpha").Set("something_id", 1).Where(ConditionRaw("id", 1)).ToSQL()
	}
}

func BenchmarkUpdateValueMapSQL(b *testing.B) {
	s := createFakeSession()

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		s.Update("alpha").Set("something_id", 1).SetMap(map[string]interface{}{"b": 1, "c": 2}).Where(ConditionRaw("id", 1)).ToSQL()
	}
}

func TestUpdateAllToSQL(t *testing.T) {
	s := createFakeSession()

	sql, args, err := s.Update("a").Set("b", 1).Set("c", 2).ToSQL()
	assert.NoError(t, err)
	assert.Equal(t, sql, "UPDATE `a` SET `b` = ?, `c` = ?")
	assert.Equal(t, args, []interface{}{1, 2})
}

func TestUpdateSingleToSQL(t *testing.T) {
	s := createFakeSession()

	sql, args, err := s.Update("a").Set("b", 1).Set("c", 2).Where(ConditionRaw("id = ?", 1)).ToSQL()
	assert.NoError(t, err)
	assert.Equal(t, sql, "UPDATE `a` SET `b` = ?, `c` = ? WHERE (id = ?)")
	assert.Equal(t, args, []interface{}{1, 2, 1})
}

func TestUpdateSetMapToSQL(t *testing.T) {
	s := createFakeSession()

	sql, args, err := s.Update("a").SetMap(map[string]interface{}{"b": 1, "c": 2}).Where(ConditionRaw("id = ?", 1)).ToSQL()
	assert.NoError(t, err)
	if sql == "UPDATE `a` SET `b` = ?, `c` = ? WHERE (id = ?)" {
		assert.Equal(t, args, []interface{}{1, 2, 1})
	} else {
		assert.Equal(t, sql, "UPDATE `a` SET `c` = ?, `b` = ? WHERE (id = ?)")
		assert.Equal(t, args, []interface{}{2, 1, 1})
	}
}

func TestUpdateSetExprToSQL(t *testing.T) {
	s := createFakeSession()

	sql, args, err := s.Update("a").Set("foo", 1).Set("bar", Expr("COALESCE(bar, 0) + 1")).Where(ConditionRaw("id = ?", 9)).ToSQL()
	assert.NoError(t, err)
	assert.Equal(t, sql, "UPDATE `a` SET `foo` = ?, `bar` = COALESCE(bar, 0) + 1 WHERE (id = ?)")
	assert.Equal(t, args, []interface{}{1, 9})

	sql, args, err = s.Update("a").Set("foo", 1).Set("bar", Expr("COALESCE(bar, 0) + ?", 2)).Where(ConditionRaw("id = ?", 9)).ToSQL()
	assert.NoError(t, err)
	assert.Equal(t, sql, "UPDATE `a` SET `foo` = ?, `bar` = COALESCE(bar, 0) + ? WHERE (id = ?)")
	assert.Equal(t, args, []interface{}{1, 2, 9})
}

func TestUpdateTenStaringFromTwentyToSQL(t *testing.T) {
	s := createFakeSession()

	sql, args, err := s.Update("a").Set("b", 1).Limit(10).Offset(20).ToSQL()
	assert.NoError(t, err)
	assert.Equal(t, sql, "UPDATE `a` SET `b` = ? LIMIT 10 OFFSET 20")
	assert.Equal(t, args, []interface{}{1})
}

func TestUpdateKeywordColumnName(t *testing.T) {
	s := createRealSessionWithFixtures()

	// Insert a user with a key
	res, err := s.InsertInto("dbr_people").Columns("name", "email", "key").Values("Benjamin", "ben@whitehouse.gov", "6").Exec()
	assert.NoError(t, err)

	// Update the key
	res, err = s.Update("dbr_people").Set("key", "6-revoked").Where(Eq{"key": "6"}).Exec()
	assert.NoError(t, err)

	// Assert our record was updated (and only our record)
	rowsAff, err := res.RowsAffected()
	assert.NoError(t, err)
	assert.Equal(t, rowsAff, int64(1))

	var person dbrPerson
	err = s.Select("*").From("dbr_people").Where(Eq{"email": "ben@whitehouse.gov"}).LoadStruct(&person)
	assert.NoError(t, err)

	assert.Equal(t, person.Name, "Benjamin")
	assert.Equal(t, person.Key.String, "6-revoked")
}

func TestUpdateReal(t *testing.T) {
	s := createRealSessionWithFixtures()

	// Insert a George
	res, err := s.InsertInto("dbr_people").Columns("name", "email").Values("George", "george@whitehouse.gov").Exec()
	assert.NoError(t, err)

	// Get George's ID
	id, err := res.LastInsertId()
	assert.NoError(t, err)

	// Rename our George to Barack
	res, err = s.Update("dbr_people").SetMap(map[string]interface{}{"name": "Barack", "email": "barack@whitehouse.gov"}).Where(ConditionRaw("id = ?", id)).Exec()

	assert.NoError(t, err)

	var person dbrPerson
	err = s.Select("*").From("dbr_people").Where(ConditionRaw("id = ?", id)).LoadStruct(&person)
	assert.NoError(t, err)

	assert.Equal(t, person.ID, id)
	assert.Equal(t, person.Name, "Barack")
	assert.Equal(t, person.Email.Valid, true)
	assert.Equal(t, person.Email.String, "barack@whitehouse.gov")
}

func TestUpdate_Prepare(t *testing.T) {

	t.Run("ToSQL Error", func(t *testing.T) {
		in := &Update{}
		in.Set("a", 1)
		stmt, err := in.Prepare()
		assert.Nil(t, stmt)
		assert.True(t, errors.IsEmpty(err))
	})

	t.Run("Prepare Error", func(t *testing.T) {
		in := &Update{
			Preparer: dbMock{
				error: errors.NewAlreadyClosedf("Who closed myself?"),
			},
		}
		in.Table = MakeAlias("tableY")
		in.Set("a", 1)

		stmt, err := in.Prepare()
		assert.Nil(t, stmt)
		assert.True(t, errors.IsAlreadyClosed(err), "%+v", err)
	})
}

func TestUpdate_Events(t *testing.T) {
	t.Parallel()

	t.Run("Stop Propagation", func(t *testing.T) {
		d := NewUpdate("tableA", "tA")
		d.Set("y", 25).Set("z", 26)

		d.Log = log.BlackHole{EnableInfo: true, EnableDebug: true}
		d.Listeners.Add(
			Listen{
				Name:      "listener1",
				EventType: OnBeforeToSQL,
				UpdateFunc: func(b *Update) {
					b.Set("a", 1)
				},
			},
			Listen{
				Name:      "listener2",
				EventType: OnBeforeToSQL,
				UpdateFunc: func(b *Update) {
					b.Set("b", 2)
					b.PropagationStopped = true
				},
			},
			Listen{
				Name:      "listener3",
				EventType: OnBeforeToSQL,
				UpdateFunc: func(b *Update) {
					panic("Should not get called")
				},
			},
		)
		sql, _, err := d.ToSQL()
		assert.NoError(t, err, "%+v", err)
		assert.Exactly(t, "UPDATE `tableA` AS `tA` SET `y` = ?, `z` = ?, `a` = ?, `b` = ?", sql)

		sql, _, err = d.ToSQL()
		assert.NoError(t, err, "%+v", err)
		assert.Exactly(t, "UPDATE `tableA` AS `tA` SET `y` = ?, `z` = ?, `a` = ?, `b` = ?, `a` = ?, `b` = ?", sql)
	})

	t.Run("Missing EventType", func(t *testing.T) {
		up := NewUpdate("tableA", "tA")
		up.Set("a", 1).Set("b", true)

		up.Listeners.Add(
			Listen{
				Name: "c=pi",
				Once: true,
				UpdateFunc: func(u *Update) {
					u.Set("c", 3.14159)
				},
			},
		)
		sql, args, err := up.ToSQL()
		assert.Empty(t, sql)
		assert.Nil(t, args)
		assert.True(t, errors.IsEmpty(err), "%+v", err)
	})

	t.Run("Should Dispatch", func(t *testing.T) {
		up := NewUpdate("tableA", "tA")
		up.Set("a", 1).Set("b", true)

		up.Listeners.Add(
			Listen{
				Name:      "c=pi",
				Once:      true,
				EventType: OnBeforeToSQL,
				UpdateFunc: func(u *Update) {
					u.Set("c", 3.14159)
				},
			},
			Listen{
				Name:      "d=d",
				Once:      true,
				EventType: OnBeforeToSQL,
				UpdateFunc: func(u *Update) {
					u.Set("d", "d")
				},
			},
		)

		up.Listeners.Add(Listen{
			Name:      "e",
			EventType: OnBeforeToSQL,
			UpdateFunc: func(u *Update) {
				u.Set("e", "e")
			},
		})

		sql, args, err := up.ToSQL()
		assert.NoError(t, err)
		assert.Exactly(t, []interface{}{1, true, 3.14159, "d", "e"}, args)
		assert.Exactly(t, "UPDATE `tableA` AS `tA` SET `a` = ?, `b` = ?, `c` = ?, `d` = ?, `e` = ?", sql)

		sql, args, err = up.ToSQL()
		assert.NoError(t, err)
		assert.Exactly(t, []interface{}{1, true, 3.14159, "d", "e", "e"}, args)
		assert.Exactly(t, "UPDATE `tableA` AS `tA` SET `a` = ?, `b` = ?, `c` = ?, `d` = ?, `e` = ?, `e` = ?", sql)

		assert.Exactly(t, `c=pi; d=d; e`, up.Listeners.String())
	})

}
