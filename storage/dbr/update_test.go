package dbr

import (
	"bytes"
	"context"
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
		_, args, err := s.Update("alpha").Set("something_id", argInt64(1)).Where(Condition("id", argInt64(1))).ToSQL()
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
			Set("something_id", argInt64(1)).
			SetMap(map[string]Argument{
				"b": argInt64(2),
				"c": argInt64(3),
			}).
			Where(Condition("id", argInt(1))).
			ToSQL()
		if err != nil {
			b.Fatalf("%+v", err)
		}
		benchmarkUpdateValuesSQL = args
	}
}

func TestUpdateAllToSQL(t *testing.T) {
	s := createFakeSession()

	sql, args, err := s.Update("a").Set("b", argInt64(1)).Set("c", argInt(2)).ToSQL()
	assert.NoError(t, err)
	assert.Equal(t, "UPDATE `a` SET `b`=?, `c`=?", sql)
	assert.Equal(t, []interface{}{int64(1), int64(2)}, args.Interfaces())
}

func TestUpdateSingleToSQL(t *testing.T) {
	s := createFakeSession()

	sql, args, err := s.Update("a").Set("b", argInt(1)).Set("c", argInt(2)).Where(Condition("id = ?", argInt(1))).ToSQL()
	assert.NoError(t, err)
	assert.Equal(t, "UPDATE `a` SET `b`=?, `c`=? WHERE (id = ?)", sql)
	assert.Equal(t, []interface{}{int64(1), int64(2), int64(1)}, args.Interfaces())
}

func TestUpdateSetMapToSQL(t *testing.T) {
	s := createFakeSession()

	sql, args, err := s.Update("a").SetMap(map[string]Argument{"b": argInt64(1), "c": ArgInt64(2)}).Where(Condition("id = ?", argInt(1))).ToSQL()
	assert.NoError(t, err)
	if sql == "UPDATE `a` SET `b`=?, `c`=? WHERE (id = ?)" {
		assert.Equal(t, []interface{}{int64(1), int64(2), int64(1)}, args.Interfaces())
	} else {
		assert.Equal(t, "UPDATE `a` SET `c`=?, `b`=? WHERE (id = ?)", sql)
		assert.Equal(t, []interface{}{int64(2), int64(1), int64(1)}, args.Interfaces())
	}
}

func TestUpdateSetExprToSQL(t *testing.T) {
	s := createFakeSession()

	sql, args, err := s.Update("a").
		Set("foo", argInt(1)).
		Set("bar", ArgExpr("COALESCE(bar, 0) + 1")).Where(Condition("id = ?", argInt(9))).ToSQL()
	assert.NoError(t, err)
	assert.Equal(t, "UPDATE `a` SET `foo`=?, `bar`=COALESCE(bar, 0) + 1 WHERE (id = ?)", sql)
	assert.Equal(t, []interface{}{int64(1), int64(9)}, args.Interfaces())

	sql, args, err = s.Update("a").
		Set("foo", argInt(1)).
		Set("bar", ArgExpr("COALESCE(bar, 0) + ?", argInt(2))).Where(Condition("id = ?", argInt(9))).ToSQL()
	assert.NoError(t, err)
	assert.Equal(t, "UPDATE `a` SET `foo`=?, `bar`=COALESCE(bar, 0) + ? WHERE (id = ?)", sql)
	assert.Equal(t, []interface{}{int64(1), int64(2), int64(9)}, args.Interfaces())
}

func TestUpdateTenStaringFromTwentyToSQL(t *testing.T) {
	s := createFakeSession()

	sql, args, err := s.Update("a").Set("b", argInt(1)).Limit(10).Offset(20).ToSQL()
	assert.NoError(t, err)
	assert.Equal(t, "UPDATE `a` SET `b`=? LIMIT 10 OFFSET 20", sql)
	assert.Equal(t, []interface{}{int64(1)}, args.Interfaces())
}

func TestUpdateKeywordColumnName(t *testing.T) {
	s := createRealSessionWithFixtures()

	// Insert a user with a key
	_, err := s.InsertInto("dbr_people").AddColumns("name", "email", "key").
		AddValues(ArgString("Benjamin"), ArgString("ben@whitehouse.gov"), ArgString("6")).Exec(context.TODO())
	assert.NoError(t, err)

	// Update the key
	res, err := s.Update("dbr_people").Set("key", ArgString("6-revoked")).Where(Eq{"key": ArgString("6")}).Exec(context.TODO())
	assert.NoError(t, err)

	// Assert our record was updated (and only our record)
	rowsAff, err := res.RowsAffected()
	assert.NoError(t, err)
	assert.Equal(t, int64(1), rowsAff)

	var person dbrPerson
	err = s.Select("*").From("dbr_people").Where(Eq{"email": ArgString("ben@whitehouse.gov")}).LoadStruct(context.TODO(), &person)
	assert.NoError(t, err)

	assert.Equal(t, "Benjamin", person.Name)
	assert.Equal(t, "6-revoked", person.Key.String)
}

func TestUpdateReal(t *testing.T) {
	s := createRealSessionWithFixtures()

	// Insert a George
	res, err := s.InsertInto("dbr_people").AddColumns("name", "email").
		AddValues(ArgString("George"), ArgString("george@whitehouse.gov")).Exec(context.TODO())
	assert.NoError(t, err)

	// Get George'ab ID
	id, err := res.LastInsertId()
	assert.NoError(t, err)

	// Rename our George to Barack
	_, err = s.Update("dbr_people").
		SetMap(map[string]Argument{"name": ArgString("Barack"), "email": ArgString("barack@whitehouse.gov")}).
		Where(Condition("id = ?", ArgInt64(id))).Exec(context.TODO())

	assert.NoError(t, err)

	var person dbrPerson
	err = s.Select("*").From("dbr_people").Where(Condition("id = ?", ArgInt64(id))).LoadStruct(context.TODO(), &person)
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
		in.Set("a", argInt(1))
		stmt, err := in.Prepare(context.TODO())
		assert.Nil(t, stmt)
		assert.True(t, errors.IsEmpty(err))
	})

	t.Run("Prepare Error", func(t *testing.T) {
		u := &Update{}
		u.DB.Preparer = dbMock{
			error: errors.NewAlreadyClosedf("Who closed myself?"),
		}
		u.Table = MakeAlias("tableY")
		u.Set("a", argInt(1))

		stmt, err := u.Prepare(context.TODO())
		assert.Nil(t, stmt)
		assert.True(t, errors.IsAlreadyClosed(err), "%+v", err)
	})
}

func TestUpdate_ToSQL_Without_Column_Arguments(t *testing.T) {
	t.Parallel()
	t.Run("with condition values", func(t *testing.T) {
		u := NewUpdate("catalog_product_entity", "cpe")
		u.SetClauses.Columns = []string{"sku", "updated_at"}
		u.Where(Condition("entity_id", ArgInt64(1, 2, 3).Operator(In)))

		sqlStr, args, err := u.ToSQL()
		assert.NoError(t, err, "%+v", err)
		assert.Exactly(t, []interface{}{int64(1), int64(2), int64(3)}, args.Interfaces())
		assert.Exactly(t,
			"UPDATE `catalog_product_entity` AS `cpe` SET `sku`=?, `updated_at`=? WHERE (`entity_id` IN ?)",
			sqlStr)
	})
	t.Run("without condition values", func(t *testing.T) {
		u := NewUpdate("catalog_product_entity", "cpe")
		u.SetClauses.Columns = []string{"sku", "updated_at"}
		u.Where(Condition("entity_id", ArgInt64().Operator(In)))

		sqlStr, args, err := u.ToSQL()
		assert.NoError(t, err, "%+v", err)
		assert.Exactly(t, []interface{}{}, args.Interfaces())
		assert.Exactly(t,
			"UPDATE `catalog_product_entity` AS `cpe` SET `sku`=?, `updated_at`=? WHERE (`entity_id` IN ?)",
			sqlStr)
	})
}

func TestUpdate_Events(t *testing.T) {
	t.Parallel()

	t.Run("Stop Propagation", func(t *testing.T) {
		d := NewUpdate("tableA", "tA")
		d.Set("y", argInt(25)).Set("z", argInt(26))

		d.Log = log.BlackHole{EnableInfo: true, EnableDebug: true}
		d.Listeners.Add(
			Listen{
				Name:      "listener1",
				EventType: OnBeforeToSQL,
				UpdateFunc: func(b *Update) {
					b.Set("a", argInt(1))
				},
			},
			Listen{
				Name:      "listener2",
				EventType: OnBeforeToSQL,
				UpdateFunc: func(b *Update) {
					b.Set("b", argInt(1))
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
		assert.Exactly(t, "UPDATE `tableA` AS `tA` SET `y`=?, `z`=?, `a`=?, `b`=?", sql)

		sql, _, err = d.ToSQL()
		assert.NoError(t, err, "%+v", err)
		assert.Exactly(t, "UPDATE `tableA` AS `tA` SET `y`=?, `z`=?, `a`=?, `b`=?, `a`=?, `b`=?", sql)
	})

	t.Run("Missing EventType", func(t *testing.T) {
		up := NewUpdate("tableA", "tA")
		up.Set("a", argInt(1)).Set("b", ArgBool(true))

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
		up.Set("a", argInt(1)).Set("b", ArgBool(true))

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
		assert.Exactly(t, "UPDATE `tableA` AS `tA` SET `a`=?, `b`=?, `c`=?, `d`=?, `e`=?", sql)

		sql, args, err = up.ToSQL()
		assert.NoError(t, err)
		assert.Exactly(t, []interface{}{int64(1), true, 3.14159, "d", "e", "e"}, args.Interfaces())
		assert.Exactly(t, "UPDATE `tableA` AS `tA` SET `a`=?, `b`=?, `c`=?, `d`=?, `e`=?, `e`=?", sql)

		assert.Exactly(t, `c=pi; d=d; e`, up.Listeners.String())
	})
}

func TestUpdatedColumns_writeOnDuplicateKey(t *testing.T) {
	t.Run("empty columns does nothing", func(t *testing.T) {
		uc := UpdatedColumns{}
		buf := new(bytes.Buffer)
		args := make(Arguments, 0, 2)
		args, err := uc.writeOnDuplicateKey(buf, args)
		assert.NoError(t, err, "%+v", err)
		assert.Empty(t, buf.String())
		assert.Empty(t, args)
	})

	t.Run("col=VALUES(col) and no arguments", func(t *testing.T) {
		uc := UpdatedColumns{
			Columns: []string{"sku", "name", "stock"},
		}
		buf := new(bytes.Buffer)
		args := make(Arguments, 0, 2)
		args, err := uc.writeOnDuplicateKey(buf, args)
		assert.NoError(t, err, "%+v", err)
		assert.Exactly(t, " ON DUPLICATE KEY UPDATE `sku`=VALUES(`sku`), `name`=VALUES(`name`), `stock`=VALUES(`stock`)", buf.String())
		assert.Empty(t, args)
	})

	t.Run("col=? and with arguments", func(t *testing.T) {
		uc := UpdatedColumns{
			Columns:   []string{"name", "stock"},
			Arguments: Arguments{ArgString("E0S 5D Mark II"), argInt64(12)},
		}
		buf := new(bytes.Buffer)
		args := make(Arguments, 0, 2)
		args, err := uc.writeOnDuplicateKey(buf, args)
		assert.NoError(t, err, "%+v", err)
		assert.Exactly(t, " ON DUPLICATE KEY UPDATE `name`=?, `stock`=?", buf.String())
		assert.Exactly(t, []interface{}{"E0S 5D Mark II", int64(12)}, args.Interfaces())
	})

	t.Run("col=VALUES(val)+? and with arguments", func(t *testing.T) {
		uc := UpdatedColumns{
			Columns:   []string{"name", "stock"},
			Arguments: Arguments{ArgString("E0S 5D Mark II"), ArgExpr("VALUES(`stock`)+?", argInt64(13))},
		}
		buf := new(bytes.Buffer)
		args := make(Arguments, 0, 2)
		args, err := uc.writeOnDuplicateKey(buf, args)
		assert.NoError(t, err, "%+v", err)
		assert.Exactly(t, " ON DUPLICATE KEY UPDATE `name`=?, `stock`=VALUES(`stock`)+?", buf.String())
		assert.Exactly(t, []interface{}{"E0S 5D Mark II", int64(13)}, args.Interfaces())
	})

	t.Run("col=VALUES(val) and with arguments and nil", func(t *testing.T) {
		uc := UpdatedColumns{
			Columns:   []string{"name", "sku", "stock"},
			Arguments: Arguments{ArgString("E0S 5D Mark III"), nil, argInt64(14)},
		}
		buf := new(bytes.Buffer)
		args := make(Arguments, 0, 2)
		args, err := uc.writeOnDuplicateKey(buf, args)
		assert.NoError(t, err, "%+v", err)
		assert.Exactly(t, " ON DUPLICATE KEY UPDATE `name`=?, `sku`=VALUES(`sku`), `stock`=?", buf.String())
		assert.Exactly(t, []interface{}{"E0S 5D Mark III", int64(14)}, args.Interfaces())
	})
}

// BenchmarkUpdatedColumns_writeOnDuplicateKey-4   	 5000000	       337 ns/op	       0 B/op	       0 allocs/op
func BenchmarkUpdatedColumns_writeOnDuplicateKey(b *testing.B) {
	buf := new(bytes.Buffer)
	args := make(Arguments, 0, 2)
	uc := UpdatedColumns{
		Columns:   []string{"name", "sku", "stock"},
		Arguments: Arguments{ArgString("E0S 5D Mark III"), nil, argInt64(14)},
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		var err error
		if args, err = uc.writeOnDuplicateKey(buf, args); err != nil {
			b.Fatalf("%+v", err)
		}
		buf.Reset()
		args = args[:0]
	}
}

func TestUpdate_SetRecord(t *testing.T) {
	t.Parallel()

	pRec := &dbrPerson{
		ID:    12345,
		Name:  "Gopher",
		Email: MakeNullString("gopher@g00gle.c0m"),
	}

	t.Run("without where", func(t *testing.T) {

		u := NewUpdate("dbr_person").SetRecord([]string{"name", "email"}, pRec)

		sqlStr, args, err := u.ToSQL()
		assert.NoError(t, err, "%+v", err)
		assert.Exactly(t, []interface{}{"Gopher", "gopher@g00gle.c0m"}, args.Interfaces())
		assert.Exactly(t,
			"UPDATE `dbr_person` SET `name`=?, `email`=?",
			sqlStr)
	})
	t.Run("with where", func(t *testing.T) {

		u := NewUpdate("dbr_person").SetRecord([]string{"name", "email"}, pRec).
			Where(Condition("id", ArgInt().Operator(Equal)))

		sqlStr, args, err := u.ToSQL()
		assert.NoError(t, err, "%+v", err)
		assert.Exactly(t, []interface{}{"Gopher", "gopher@g00gle.c0m", int64(12345)}, args.Interfaces())
		assert.Exactly(t,
			"UPDATE `dbr_person` SET `name`=?, `email`=? WHERE (`id` = ?)",
			sqlStr)
	})

	t.Run("fails column not in entity object", func(t *testing.T) {
		u := NewUpdate("dbr_person").SetRecord([]string{"name", "email"}, pRec).
			Set("key", ArgString("JustAKey")).
			Where(Condition("id", ArgInt().Operator(Equal)))

		sqlStr, args, err := u.ToSQL()
		assert.Empty(t, sqlStr)
		assert.Nil(t, args)
		assert.True(t, errors.IsNotFound(err), "%+v", err)
	})

}
