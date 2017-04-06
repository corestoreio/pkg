package dbr

import (
	"testing"

	"github.com/corestoreio/errors"
	"github.com/corestoreio/log"
	"github.com/stretchr/testify/assert"
)

var benchmarkUpdateValuesSQL Arguments

func BenchmarkUpdateValuesSQL(b *testing.B) {
	s := createFakeSession()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, args, err := s.Update("alpha").Set("something_id", ArgInt64(1)).Where(Condition("id", ArgInt64(1))).ToSQL()
		if err != nil {
			b.Fatalf("%+v", err)
		}
		benchmarkUpdateValuesSQL = args
	}
}

func BenchmarkUpdateValueMapSQL(b *testing.B) {
	s := createFakeSession()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, args, err := s.Update("alpha").
			Set("something_id", ArgInt64(1)).
			SetMap(map[string]Argument{
				"b": ArgInt64(2),
				"c": ArgInt64(3),
			}).
			Where(Condition("id", ArgInt(1))).
			ToSQL()
		if err != nil {
			b.Fatalf("%+v", err)
		}
		benchmarkUpdateValuesSQL = args
	}
}

func TestUpdateAllToSQL(t *testing.T) {
	s := createFakeSession()

	sql, args, err := s.Update("a").Set("b", ArgInt64(1)).Set("c", ArgInt(2)).ToSQL()
	assert.NoError(t, err)
	assert.Equal(t, "UPDATE `a` SET `b` = ?, `c` = ?", sql)
	assert.Equal(t, []interface{}{int64(1), int64(2)}, args.Interfaces())
}

func TestUpdateSingleToSQL(t *testing.T) {
	s := createFakeSession()

	sql, args, err := s.Update("a").Set("b", ArgInt(1)).Set("c", ArgInt(2)).Where(Condition("id = ?", ArgInt(1))).ToSQL()
	assert.NoError(t, err)
	assert.Equal(t, "UPDATE `a` SET `b` = ?, `c` = ? WHERE (id = ?)", sql)
	assert.Equal(t, []interface{}{int64(1), int64(2), int64(1)}, args.Interfaces())
}

func TestUpdateSetMapToSQL(t *testing.T) {
	s := createFakeSession()

	sql, args, err := s.Update("a").SetMap(map[string]Argument{"b": ArgInt64(1), "c": ArgInt64(2)}).Where(Condition("id = ?", ArgInt(1))).ToSQL()
	assert.NoError(t, err)
	if sql == "UPDATE `a` SET `b` = ?, `c` = ? WHERE (id = ?)" {
		assert.Equal(t, []interface{}{int64(1), int64(2), int64(1)}, args.Interfaces())
	} else {
		assert.Equal(t, "UPDATE `a` SET `c` = ?, `b` = ? WHERE (id = ?)", sql)
		assert.Equal(t, []interface{}{int64(2), int64(1), int64(1)}, args.Interfaces())
	}
}

func TestUpdateSetExprToSQL(t *testing.T) {
	s := createFakeSession()

	sql, args, err := s.Update("a").
		Set("foo", ArgInt(1)).
		Set("bar", Expr("COALESCE(bar, 0) + 1")).Where(Condition("id = ?", ArgInt(9))).ToSQL()
	assert.NoError(t, err)
	assert.Equal(t, "UPDATE `a` SET `foo` = ?, `bar` = COALESCE(bar, 0) + 1 WHERE (id = ?)", sql)
	assert.Equal(t, []interface{}{int64(1), int64(9)}, args.Interfaces())

	sql, args, err = s.Update("a").
		Set("foo", ArgInt(1)).
		Set("bar", Expr("COALESCE(bar, 0) + ?", ArgInt(2))).Where(Condition("id = ?", ArgInt(9))).ToSQL()
	assert.NoError(t, err)
	assert.Equal(t, "UPDATE `a` SET `foo` = ?, `bar` = COALESCE(bar, 0) + ? WHERE (id = ?)", sql)
	assert.Equal(t, []interface{}{int64(1), int64(2), int64(9)}, args.Interfaces())
}

func TestUpdateTenStaringFromTwentyToSQL(t *testing.T) {
	s := createFakeSession()

	sql, args, err := s.Update("a").Set("b", ArgInt(1)).Limit(10).Offset(20).ToSQL()
	assert.NoError(t, err)
	assert.Equal(t, "UPDATE `a` SET `b` = ? LIMIT 10 OFFSET 20", sql)
	assert.Equal(t, []interface{}{int64(1)}, args.Interfaces())
}

func TestUpdateKeywordColumnName(t *testing.T) {
	s := createRealSessionWithFixtures()

	// Insert a user with a key
	res, err := s.InsertInto("dbr_people").Columns("name", "email", "key").
		Values(ArgString("Benjamin"), ArgString("ben@whitehouse.gov"), ArgString("6")).Exec()
	assert.NoError(t, err)

	// Update the key
	res, err = s.Update("dbr_people").Set("key", ArgString("6-revoked")).Where(Eq{"key": ArgString("6")}).Exec()
	assert.NoError(t, err)

	// Assert our record was updated (and only our record)
	rowsAff, err := res.RowsAffected()
	assert.NoError(t, err)
	assert.Equal(t, int64(1), rowsAff)

	var person dbrPerson
	err = s.Select("*").From("dbr_people").Where(Eq{"email": ArgString("ben@whitehouse.gov")}).LoadStruct(&person)
	assert.NoError(t, err)

	assert.Equal(t, "Benjamin", person.Name)
	assert.Equal(t, "6-revoked", person.Key.String)
}

func TestUpdateReal(t *testing.T) {
	s := createRealSessionWithFixtures()

	// Insert a George
	res, err := s.InsertInto("dbr_people").Columns("name", "email").
		Values(ArgString("George"), ArgString("george@whitehouse.gov")).Exec()
	assert.NoError(t, err)

	// Get George'ab ID
	id, err := res.LastInsertId()
	assert.NoError(t, err)

	// Rename our George to Barack
	res, err = s.Update("dbr_people").
		SetMap(map[string]Argument{"name": ArgString("Barack"), "email": ArgString("barack@whitehouse.gov")}).
		Where(Condition("id = ?", ArgInt64(id))).Exec()

	assert.NoError(t, err)

	var person dbrPerson
	err = s.Select("*").From("dbr_people").Where(Condition("id = ?", ArgInt64(id))).LoadStruct(&person)
	assert.NoError(t, err)

	assert.Equal(t, id, person.ID)
	assert.Equal(t, "Barack", person.Name)
	assert.Equal(t, true, person.Email.Valid)
	assert.Equal(t, "barack@whitehouse.gov", person.Email.String)
}

func TestUpdate_Prepare(t *testing.T) {
	t.Parallel()
	t.Run("ToSQL Error", func(t *testing.T) {
		in := &Update{}
		in.Set("a", ArgInt(1))
		stmt, err := in.Prepare()
		assert.Nil(t, stmt)
		assert.True(t, errors.IsEmpty(err))
	})

	t.Run("Prepare Error", func(t *testing.T) {
		u := &Update{}
		u.DB.Preparer = dbMock{
			error: errors.NewAlreadyClosedf("Who closed myself?"),
		}
		u.Table = MakeAlias("tableY")
		u.Set("a", ArgInt(1))

		stmt, err := u.Prepare()
		assert.Nil(t, stmt)
		assert.True(t, errors.IsAlreadyClosed(err), "%+v", err)
	})
}

func TestUpdate_ToSQL_Without_Column_Arguments(t *testing.T) {
	t.Parallel()
	t.Run("with condition values", func(t *testing.T) {
		u := NewUpdate("catalog_product_entity", "cpe")
		u.SetClauses.Columns = []string{"sku", "updated_at"}
		u.Where(Condition("entity_id", ArgInt64(1, 2, 3).Operator('i')))

		sqlStr, args, err := u.ToSQL()
		assert.NoError(t, err, "%+v", err)
		assert.Exactly(t, []interface{}{int64(1), int64(2), int64(3)}, args.Interfaces())
		assert.Exactly(t,
			"UPDATE `catalog_product_entity` AS `cpe` SET `sku` = ?, `updated_at` = ? WHERE (`entity_id` IN ?)",
			sqlStr)
	})
	t.Run("without condition values", func(t *testing.T) {
		u := NewUpdate("catalog_product_entity", "cpe")
		u.SetClauses.Columns = []string{"sku", "updated_at"}
		u.Where(Condition("entity_id", ArgInt64().Operator('i')))

		sqlStr, args, err := u.ToSQL()
		assert.NoError(t, err, "%+v", err)
		assert.Exactly(t, []interface{}{}, args.Interfaces())
		assert.Exactly(t,
			"UPDATE `catalog_product_entity` AS `cpe` SET `sku` = ?, `updated_at` = ? WHERE (`entity_id` IN ?)",
			sqlStr)
	})
}

func TestUpdateMulti_Exec(t *testing.T) {

	//mu := NewUpdateMulti("catalog_product_entity", "cpe")
	//mu.Stmt.SetClauses.Columns = []string{"sku", "updated_at"}
	//mu.Stmt.Where(Condition("entity_id", ArgInt64().Operator('i'))) // ArgInt64 must be without arguments

	t.Run("no columns provided", func(t *testing.T) {
		mu := NewUpdateMulti("catalog_product_entity", "cpe")
		mu.Stmt.Where(Condition("entity_id", ArgInt64().Operator('i'))) // ArgInt64 must be without arguments
		res, err := mu.Exec()
		assert.Nil(t, res)
		assert.True(t, errors.IsEmpty(err), "%+v", err)
	})
}

func TestUpdate_Events(t *testing.T) {
	t.Parallel()

	t.Run("Stop Propagation", func(t *testing.T) {
		d := NewUpdate("tableA", "tA")
		d.Set("y", ArgInt(25)).Set("z", ArgInt(26))

		d.Log = log.BlackHole{EnableInfo: true, EnableDebug: true}
		d.Listeners.Add(
			Listen{
				Name:      "listener1",
				EventType: OnBeforeToSQL,
				UpdateFunc: func(b *Update) {
					b.Set("a", ArgInt(1))
				},
			},
			Listen{
				Name:      "listener2",
				EventType: OnBeforeToSQL,
				UpdateFunc: func(b *Update) {
					b.Set("b", ArgInt(1))
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
		up.Set("a", ArgInt(1)).Set("b", ArgBool(true))

		up.Listeners.Add(
			Listen{
				Name: "c=pi",
				Once: true,
				UpdateFunc: func(u *Update) {
					u.Set("c", ArgFloat64(3.14159))
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
		up.Set("a", ArgInt(1)).Set("b", ArgBool(true))

		up.Listeners.Add(
			Listen{
				Name:      "c=pi",
				Once:      true,
				EventType: OnBeforeToSQL,
				UpdateFunc: func(u *Update) {
					u.Set("c", ArgFloat64(3.14159))
				},
			},
			Listen{
				Name:      "d=d",
				Once:      true,
				EventType: OnBeforeToSQL,
				UpdateFunc: func(u *Update) {
					u.Set("d", ArgString("d"))
				},
			},
		)

		up.Listeners.Add(Listen{
			Name:      "e",
			EventType: OnBeforeToSQL,
			UpdateFunc: func(u *Update) {
				u.Set("e", ArgString("e"))
			},
		})

		sql, args, err := up.ToSQL()
		assert.NoError(t, err)
		assert.Exactly(t, []interface{}{int64(1), true, 3.14159, "d", "e"}, args.Interfaces())
		assert.Exactly(t, "UPDATE `tableA` AS `tA` SET `a` = ?, `b` = ?, `c` = ?, `d` = ?, `e` = ?", sql)

		sql, args, err = up.ToSQL()
		assert.NoError(t, err)
		assert.Exactly(t, []interface{}{int64(1), true, 3.14159, "d", "e", "e"}, args.Interfaces())
		assert.Exactly(t, "UPDATE `tableA` AS `tA` SET `a` = ?, `b` = ?, `c` = ?, `d` = ?, `e` = ?, `e` = ?", sql)

		assert.Exactly(t, `c=pi; d=d; e`, up.Listeners.String())
	})

}
