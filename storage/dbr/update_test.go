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
	"context"
	"testing"

	"github.com/corestoreio/errors"
	"github.com/corestoreio/log"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestUpdateAllToSQL(t *testing.T) {
	t.Parallel()
	qb := NewUpdate("a").Set("b", argInt64(1)).Set("c", argInt(2))
	compareToSQL(t, qb, nil, "UPDATE `a` SET `b`=?, `c`=?", "", int64(1), int64(2))
}

func TestUpdateSingleToSQL(t *testing.T) {
	t.Parallel()

	compareToSQL(t, NewUpdate("a").
		Set("b", argInt(1)).Set("c", argInt(2)).Where(Column("id", argInt(1))),
		nil,
		"UPDATE `a` SET `b`=?, `c`=? WHERE (`id` = ?)",
		"UPDATE `a` SET `b`=1, `c`=2 WHERE (`id` = 1)",
		int64(1), int64(2), int64(1))

}

func TestUpdateSetMapToSQL(t *testing.T) {
	t.Parallel()
	s := createFakeSession()

	sql, args, err := s.Update("a").SetMap(map[string]Argument{"b": argInt64(1), "c": ArgInt64(2)}).Where(Column("id", argInt(1))).ToSQL()
	assert.NoError(t, err)
	if sql == "UPDATE `a` SET `b`=?, `c`=? WHERE (`id` = ?)" {
		assert.Equal(t, []interface{}{int64(1), int64(2), int64(1)}, args.Interfaces())
	} else {
		assert.Equal(t, "UPDATE `a` SET `c`=?, `b`=? WHERE (`id` = ?)", sql)
		assert.Equal(t, []interface{}{int64(2), int64(1), int64(1)}, args.Interfaces())
	}
}

func TestUpdateSetExprToSQL(t *testing.T) {
	t.Parallel()

	compareToSQL(t, NewUpdate("a").
		Set("foo", argInt(1)).
		Set("bar", ArgExpr("COALESCE(bar, 0) + 1")).Where(Column("id", argInt(9))),
		nil,
		"UPDATE `a` SET `foo`=?, `bar`=COALESCE(bar, 0) + 1 WHERE (`id` = ?)",
		"UPDATE `a` SET `foo`=1, `bar`=COALESCE(bar, 0) + 1 WHERE (`id` = 9)",
		int64(1), int64(9))

	compareToSQL(t, NewUpdate("a").
		Set("foo", argInt(1)).
		Set("bar", ArgExpr("COALESCE(bar, 0) + ?", argInt(2))).Where(Column("id", argInt(9))),
		nil,
		"UPDATE `a` SET `foo`=?, `bar`=COALESCE(bar, 0) + ? WHERE (`id` = ?)",
		"UPDATE `a` SET `foo`=1, `bar`=COALESCE(bar, 0) + 2 WHERE (`id` = 9)",
		int64(1), int64(2), int64(9))
}

func TestUpdate_Limit_Offset(t *testing.T) {
	t.Parallel()

	compareToSQL(t, NewUpdate("a").
		Set("b", argInt(1)).Limit(10).Offset(20),
		nil,
		"UPDATE `a` SET `b`=? LIMIT 10 OFFSET 20",
		"UPDATE `a` SET `b`=1 LIMIT 10 OFFSET 20",
		int64(1))
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
	_, err = s.Select("*").From("dbr_people").Where(Eq{"email": ArgString("ben@whitehouse.gov")}).Load(context.TODO(), &person)
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
		Where(Column("id", ArgInt64(id))).Exec(context.TODO())

	assert.NoError(t, err)

	var person dbrPerson
	_, err = s.Select("*").From("dbr_people").Where(Column("id", ArgInt64(id))).Load(context.TODO(), &person)
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
		u.Where(Column("entity_id", In.Int64(1, 2, 3)))

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
		u.Where(Column("entity_id", In.Int64()))

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
			Where(Column("id", Equal.Int()))

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
			Where(Column("id", Equal.Int()))

		sqlStr, args, err := u.ToSQL()
		assert.Empty(t, sqlStr)
		assert.Nil(t, args)
		assert.True(t, errors.IsNotFound(err), "%+v", err)
	})
}

func TestUpdate_UseBuildCache(t *testing.T) {
	t.Parallel()

	up := NewUpdate("a").
		Set("foo", argInt(1)).
		Set("bar", ArgExpr("COALESCE(bar, 0) + ?", argInt(2))).Where(Column("id", argInt(9)))

	up.UseBuildCache = true

	const cachedSQLPlaceHolder = "UPDATE `a` SET `foo`=?, `bar`=COALESCE(bar, 0) + ? WHERE (`id` = ?)"
	t.Run("without interpolate", func(t *testing.T) {
		for i := 0; i < 3; i++ {
			sql, args, err := up.ToSQL()
			require.NoError(t, err, "%+v", err)
			require.Equal(t, cachedSQLPlaceHolder, sql)
			assert.Equal(t, []interface{}{int64(1), int64(2), int64(9)}, args.Interfaces())
			assert.Equal(t, cachedSQLPlaceHolder, string(up.buildCache))
		}
	})

	t.Run("with interpolate", func(t *testing.T) {
		up.Interpolate()
		up.buildCache = nil
		up.RawArguments = nil

		const cachedSQLInterpolated = "UPDATE `a` SET `foo`=1, `bar`=COALESCE(bar, 0) + 2 WHERE (`id` = 9)"
		for i := 0; i < 3; i++ {
			sql, args, err := up.ToSQL()
			assert.Equal(t, cachedSQLPlaceHolder, string(up.buildCache))
			require.NoError(t, err, "%+v", err)
			require.Equal(t, cachedSQLInterpolated, sql)
			assert.Nil(t, args)
		}
	})
}
