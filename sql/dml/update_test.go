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

package dml

import (
	"context"
	"testing"

	"github.com/corestoreio/errors"
	"github.com/corestoreio/log"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestUpdate_Basics(t *testing.T) {
	t.Parallel()

	t.Run("all rows", func(t *testing.T) {
		qb := NewUpdate("a").Set(
			Column("b").Int64(1),
			Column("c").Int(2))
		compareToSQL(t, qb, nil, "UPDATE `a` SET `b`=?, `c`=?", "", int64(1), int64(2))
	})
	t.Run("single row", func(t *testing.T) {
		compareToSQL(t, NewUpdate("a").
			Set(
				Column("b").Int(1), Column("c").Int(2),
			).Where(Column("id").Int(1)),
			nil,
			"UPDATE `a` SET `b`=?, `c`=? WHERE (`id` = ?)",
			"UPDATE `a` SET `b`=1, `c`=2 WHERE (`id` = 1)",
			int64(1), int64(2), int64(1))
	})
	t.Run("order by", func(t *testing.T) {
		qb := NewUpdate("a").Set(Column("b").Int(1), Column("c").Int(2)).
			OrderBy("col1", "col2").OrderByDesc("col2", "col3").Unsafe().OrderBy("concat(1,2,3)")
		compareToSQL(t, qb, nil,
			"UPDATE `a` SET `b`=?, `c`=? ORDER BY `col1`, `col2`, `col2` DESC, `col3` DESC, concat(1,2,3)",
			"UPDATE `a` SET `b`=1, `c`=2 ORDER BY `col1`, `col2`, `col2` DESC, `col3` DESC, concat(1,2,3)",
			int64(1), int64(2))
	})
	t.Run("limit offset", func(t *testing.T) {
		compareToSQL(t, NewUpdate("a").Set(Column("b").Int(1)).Limit(10),
			nil,
			"UPDATE `a` SET `b`=? LIMIT 10",
			"UPDATE `a` SET `b`=1 LIMIT 10",
			int64(1))
	})
	t.Run("same column name in SET and WHERE", func(t *testing.T) {
		compareToSQL(t, NewUpdate("dml_people").Set(Column("key").Str("6-revoked")).Where(Column("key").Str("6")),
			nil,
			"UPDATE `dml_people` SET `key`=? WHERE (`key` = ?)",
			"UPDATE `dml_people` SET `key`='6-revoked' WHERE (`key` = '6')",
			"6-revoked", "6")
	})
}

func TestUpdateSetExprToSQL(t *testing.T) {
	t.Parallel()

	compareToSQL(t, NewUpdate("a").
		Set(
			Column("foo").Int(1),
			Column("bar").Expr("COALESCE(bar, 0) + 1"),
		).Where(Column("id").Int(9)),
		nil,
		"UPDATE `a` SET `foo`=?, `bar`=COALESCE(bar, 0) + 1 WHERE (`id` = ?)",
		"UPDATE `a` SET `foo`=1, `bar`=COALESCE(bar, 0) + 1 WHERE (`id` = 9)",
		int64(1), int64(9))

	compareToSQL(t, NewUpdate("a").
		Set(
			Column("foo").Int(1),
			Column("bar").Expr("COALESCE(bar, 0) + 1"),
		).Where(Column("id").In().Int64s(10, 11)),
		nil,
		"UPDATE `a` SET `foo`=?, `bar`=COALESCE(bar, 0) + 1 WHERE (`id` IN (?,?))",
		"UPDATE `a` SET `foo`=1, `bar`=COALESCE(bar, 0) + 1 WHERE (`id` IN (10,11))",
		int64(1), int64(10), int64(11))

	compareToSQL(t, NewUpdate("a").
		Set(
			Column("foo").Int(1),
			Column("bar").Expr("COALESCE(bar, 0) + ?").Int(2),
		).
		Where(Column("id").Int(9)),
		nil,
		"UPDATE `a` SET `foo`=?, `bar`=COALESCE(bar, 0) + ? WHERE (`id` = ?)",
		"UPDATE `a` SET `foo`=1, `bar`=COALESCE(bar, 0) + 2 WHERE (`id` = 9)",
		int64(1), int64(2), int64(9))
}

func TestUpdateKeywordColumnName(t *testing.T) {
	s := createRealSessionWithFixtures(t, nil)

	// Insert a user with a key
	_, err := s.InsertInto("dml_people").AddColumns("name", "email", "key").
		AddValues("Benjamin", "ben@whitehouse.gov", "6").Exec(context.TODO())
	assert.NoError(t, err)

	// Update the key
	res, err := s.Update("dml_people").Set(Column("key").Str("6-revoked")).Where(Column("key").Str("6")).Exec(context.TODO())
	assert.NoError(t, err)

	// Assert our record was updated (and only our record)
	rowsAff, err := res.RowsAffected()
	assert.NoError(t, err)
	assert.Equal(t, int64(1), rowsAff)

	var person dmlPerson
	_, err = s.SelectFrom("dml_people").Star().Where(Column("email").Str("ben@whitehouse.gov")).Load(context.TODO(), &person)
	assert.NoError(t, err)

	assert.Equal(t, "Benjamin", person.Name)
	assert.Equal(t, "6-revoked", person.Key.String)
}

func TestUpdateReal(t *testing.T) {
	s := createRealSessionWithFixtures(t, nil)

	// Insert a George
	res, err := s.InsertInto("dml_people").AddColumns("name", "email").
		AddValues("George", "george@whitehouse.gov").Exec(context.TODO())
	assert.NoError(t, err)

	// Get George'ab ID
	id, err := res.LastInsertId()
	assert.NoError(t, err)

	// Rename our George to Barack
	_, err = s.Update("dml_people").
		Set(Column("name").Str("Barack"), Column("email").Str("barack@whitehouse.gov")).
		Where(Column("id").In().Int64s(id, 8888)).Exec(context.TODO())
	// Meaning of 8888: Just to see if the SQL with place holders gets created correctly
	require.NoError(t, err)

	var person dmlPerson
	_, err = s.SelectFrom("dml_people").Star().Where(Column("id").Int64(id)).Load(context.TODO(), &person)
	assert.NoError(t, err)

	assert.Equal(t, id, int64(person.ID))
	assert.Equal(t, "Barack", person.Name)
	assert.Equal(t, true, person.Email.Valid)
	assert.Equal(t, "barack@whitehouse.gov", person.Email.String)
}

func TestUpdate_Prepare(t *testing.T) {
	t.Parallel()
	t.Run("ToSQL Error", func(t *testing.T) {
		in := &Update{}
		in.Set(Column("a").Int(1))
		stmt, err := in.Prepare(context.TODO())
		assert.Nil(t, stmt)
		assert.True(t, errors.IsEmpty(err))
	})

	t.Run("Prepare Error", func(t *testing.T) {
		u := &Update{}
		u.DB = dbMock{
			error: errors.NewAlreadyClosedf("Who closed myself?"),
		}
		u.Table.Name = "tableY"
		u.Set(Column("a").Int(1))

		stmt, err := u.Prepare(context.TODO())
		assert.Nil(t, stmt)
		assert.True(t, errors.IsAlreadyClosed(err), "%+v", err)
	})
}

func TestUpdate_ToSQL_Without_Column_Arguments(t *testing.T) {
	t.Parallel()
	t.Run("with condition values", func(t *testing.T) {
		u := NewUpdate("catalog_product_entity").AddColumns("sku", "updated_at")
		u.Where(Column("entity_id").In().Int64s(1, 2, 3))
		compareToSQL(t, u, nil,
			"UPDATE `catalog_product_entity` SET `sku`=?, `updated_at`=? WHERE (`entity_id` IN (?,?,?))",
			"",
			int64(1), int64(2), int64(3),
		)
	})
	t.Run("without condition values", func(t *testing.T) {
		u := NewUpdate("catalog_product_entity").AddColumns("sku", "updated_at")
		u.Where(Column("entity_id").In().PlaceHolder())

		args := []interface{}{}
		compareToSQL(t, u, nil,
			"UPDATE `catalog_product_entity` SET `sku`=?, `updated_at`=? WHERE (`entity_id` IN (?))",
			"",
			args...,
		)
	})
}

func TestUpdate_Events(t *testing.T) {
	t.Parallel()

	t.Run("Stop Propagation", func(t *testing.T) {
		d := NewUpdate("tableA")
		d.Set(Column("y").Int(25), Column("z").Int(26))

		d.Log = log.BlackHole{EnableInfo: true, EnableDebug: true}
		d.Listeners.Add(
			Listen{
				Name:      "listener1",
				EventType: OnBeforeToSQL,
				UpdateFunc: func(b *Update) {
					b.Set(Column("a").Int(1))
				},
			},
			Listen{
				Name:      "listener2",
				EventType: OnBeforeToSQL,
				UpdateFunc: func(b *Update) {
					b.Set(Column("b").Int(1))
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
		compareToSQL(t, d, nil,
			"UPDATE `tableA` SET `y`=?, `z`=?, `a`=?, `b`=?",
			"UPDATE `tableA` SET `y`=25, `z`=26, `a`=1, `b`=1, `a`=1, `b`=1", // each call ToSQL appends more columns
			int64(25), int64(26), int64(1), int64(1),
		)
	})

	t.Run("Missing EventType", func(t *testing.T) {
		up := NewUpdate("tableA")
		up.Set(Column("a").Int(1), Column("b").Bool(true))

		up.Listeners.Add(
			Listen{
				Name: "c=pi",
				Once: true,
				UpdateFunc: func(u *Update) {
					u.Set(Column("c").Float64(3.14159))
				},
			},
		)
		compareToSQL(t, up, errors.IsEmpty,
			"",
			"",
		)
	})

	t.Run("Should Dispatch", func(t *testing.T) {
		up := NewUpdate("tableA")
		up.Set(Column("a").Int(1), Column("b").Bool(true))
		up.Listeners.Add(
			Listen{
				Name:      "c=pi",
				Once:      true,
				EventType: OnBeforeToSQL,
				UpdateFunc: func(u *Update) {
					u.Set(Column("c").Float64(3.14159))
				},
			},
			Listen{
				Name:      "d=d",
				Once:      true,
				EventType: OnBeforeToSQL,
				UpdateFunc: func(u *Update) {
					u.Set(Column("d").Str("d"))
				},
			},
		)

		up.Listeners.Add(Listen{
			Name:      "e",
			EventType: OnBeforeToSQL,
			UpdateFunc: func(u *Update) {
				u.Set(Column("e").Str("e"))
			},
		})
		compareToSQL(t, up, nil,
			"UPDATE `tableA` SET `a`=?, `b`=?, `c`=?, `d`=?, `e`=?",
			"UPDATE `tableA` SET `a`=1, `b`=1, `c`=3.14159, `d`='d', `e`='e', `e`='e'", // each call ToSQL appends more columns
			int64(1), true, 3.14159, "d", "e",
		)
		assert.Exactly(t, `c=pi; d=d; e`, up.Listeners.String())
	})
}

func TestUpdate_SetRecord(t *testing.T) {
	t.Parallel()

	pRec := &dmlPerson{
		ID:    12345,
		Name:  "Gopher",
		Email: MakeNullString("gopher@g00gle.c0m"),
	}

	t.Run("without where", func(t *testing.T) {
		u := NewUpdate("dml_person").AddColumns("name", "email").BindRecord(Qualify("", pRec))
		compareToSQL(t, u, nil,
			"UPDATE `dml_person` SET `name`=?, `email`=?",
			"UPDATE `dml_person` SET `name`='Gopher', `email`='gopher@g00gle.c0m'",
			"Gopher", "gopher@g00gle.c0m",
		)
	})
	t.Run("with where", func(t *testing.T) {
		u := NewUpdate("dml_person").AddColumns("name", "email").BindRecord(Qualify("", pRec)).
			Where(Column("id").PlaceHolder())
		compareToSQL(t, u, nil,
			"UPDATE `dml_person` SET `name`=?, `email`=? WHERE (`id` = ?)",
			"UPDATE `dml_person` SET `name`='Gopher', `email`='gopher@g00gle.c0m' WHERE (`id` = 12345)",
			"Gopher", "gopher@g00gle.c0m", int64(12345),
		)
	})
	t.Run("fails column not in entity object", func(t *testing.T) {
		u := NewUpdate("dml_person").AddColumns("name", "email").BindRecord(Qualify("", pRec)).
			Set(Column("key").Str("JustAKey")).
			Where(Column("id").PlaceHolder())
		compareToSQL(t, u, errors.IsNotFound,
			"",
			"",
		)
	})
}

func TestUpdate_UseBuildCache(t *testing.T) {
	t.Parallel()

	up := NewUpdate("a").
		Set(
			Column("foo").Int(1),
			Column("bar").Expr("COALESCE(bar, 0) + ?").Int(2)).
		Where(Column("id").Int(9)).BuildCache()

	const cachedSQLPlaceHolder = "UPDATE `a` SET `foo`=?, `bar`=COALESCE(bar, 0) + ? WHERE (`id` = ?)"
	t.Run("without interpolate", func(t *testing.T) {
		for i := 0; i < 3; i++ {
			compareToSQL(t, up, nil,
				cachedSQLPlaceHolder,
				"",
				int64(1), int64(2), int64(9),
			)
			assert.Equal(t, cachedSQLPlaceHolder, string(up.cacheSQL))
		}
	})

	t.Run("with interpolate", func(t *testing.T) {
		up.cacheSQL = nil

		const cachedSQLInterpolated = "UPDATE `a` SET `foo`=1, `bar`=COALESCE(bar, 0) + 2 WHERE (`id` = 9)"
		for i := 0; i < 3; i++ {
			compareToSQL(t, up, nil,
				cachedSQLPlaceHolder,
				cachedSQLInterpolated,
				int64(1), int64(2), int64(9),
			)
			assert.Equal(t, cachedSQLPlaceHolder, string(up.cacheSQL))
		}
	})
}
