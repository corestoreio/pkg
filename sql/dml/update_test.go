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
		compareToSQL(t, qb, errors.NoKind, "UPDATE `a` SET `b`=1, `c`=2", "")
	})
	t.Run("single row", func(t *testing.T) {
		compareToSQL(t, NewUpdate("a").
			Set(
				Column("b").Int(1), Column("c").Int(2),
			).Where(Column("id").Int(1)),
			errors.NoKind,
			"UPDATE `a` SET `b`=1, `c`=2 WHERE (`id` = 1)",
			"UPDATE `a` SET `b`=1, `c`=2 WHERE (`id` = 1)")
	})
	t.Run("order by", func(t *testing.T) {
		qb := NewUpdate("a").Set(Column("b").Int(1), Column("c").Int(2)).
			OrderBy("col1", "col2").OrderByDesc("col2", "col3").Unsafe().OrderBy("concat(1,2,3)")
		compareToSQL(t, qb, errors.NoKind,
			"UPDATE `a` SET `b`=1, `c`=2 ORDER BY `col1`, `col2`, `col2` DESC, `col3` DESC, concat(1,2,3)",
			"UPDATE `a` SET `b`=1, `c`=2 ORDER BY `col1`, `col2`, `col2` DESC, `col3` DESC, concat(1,2,3)")
	})
	t.Run("limit offset", func(t *testing.T) {
		compareToSQL(t, NewUpdate("a").Set(Column("b").Int(1)).Limit(10),
			errors.NoKind,
			"UPDATE `a` SET `b`=1 LIMIT 10",
			"UPDATE `a` SET `b`=1 LIMIT 10")
	})
	t.Run("same column name in SET and WHERE", func(t *testing.T) {
		compareToSQL(t, NewUpdate("dml_people").Set(Column("key").Str("6-revoked")).Where(Column("key").Str("6")),
			errors.NoKind,
			"UPDATE `dml_people` SET `key`='6-revoked' WHERE (`key` = '6')",
			"UPDATE `dml_people` SET `key`='6-revoked' WHERE (`key` = '6')")
	})

	t.Run("placeholder in columns", func(t *testing.T) {
		u := NewUpdate("dml_people").Set(
			Column("key").PlaceHolder(),
		).Where(Column("key").Str("6")).WithArgs().String("Ke' --yX")
		compareToSQL(t, u,
			errors.NoKind,
			"UPDATE `dml_people` SET `key`=? WHERE (`key` = '6')",
			"UPDATE `dml_people` SET `key`='Ke\\' --yX' WHERE (`key` = '6')",
			"Ke' --yX")
	})
}

func TestUpdate_SetExprToSQL(t *testing.T) {
	t.Parallel()

	t.Run("no placeholder", func(t *testing.T) {
		compareToSQL(t, NewUpdate("a").
			Set(
				Column("foo").Int(1),
				Column("bar").Expr("COALESCE(bar, 0) + 1"),
			).Where(Column("id").Int(9)),
			errors.NoKind,
			"UPDATE `a` SET `foo`=1, `bar`=COALESCE(bar, 0) + 1 WHERE (`id` = 9)",
			"UPDATE `a` SET `foo`=1, `bar`=COALESCE(bar, 0) + 1 WHERE (`id` = 9)",
		)
	})

	t.Run("with slice in WHERE clause", func(t *testing.T) {
		compareToSQL(t, NewUpdate("a").
			Set(
				Column("foo").Int(1),
				Column("bar").Expr("COALESCE(bar, 0) + 1"),
			).Where(Column("id").In().Int64s(10, 11)),
			errors.NoKind,
			"UPDATE `a` SET `foo`=1, `bar`=COALESCE(bar, 0) + 1 WHERE (`id` IN (10,11))",
			"UPDATE `a` SET `foo`=1, `bar`=COALESCE(bar, 0) + 1 WHERE (`id` IN (10,11))",
		)
	})

	t.Run("with placeholder", func(t *testing.T) {
		u := NewUpdate("a").
			Set(
				Column("fooNULL").PlaceHolder(),
				Column("bar99").Expr("COALESCE(bar, 0) + ?"),
			).
			Where(Column("id").Int(9)).
			WithArgs().NullString(NullString{}).Uint(99)
		compareToSQL(t, u, errors.NoKind,
			"UPDATE `a` SET `fooNULL`=?, `bar99`=COALESCE(bar, 0) + ? WHERE (`id` = 9)",
			"", //"UPDATE `a` SET `foo`=1, `bar`=COALESCE(bar, 0) + 2 WHERE (`id` = 9)",
			nil, int64(99))
		assert.Exactly(t, []string{"fooNULL", "bar99"}, u.base.qualifiedColumns)
	})
}

func TestUpdateKeywordColumnName(t *testing.T) {
	s := createRealSessionWithFixtures(t, nil)
	defer testCloser(t, s)

	// Insert a user with a key
	_, err := s.InsertInto("dml_people").AddColumns("name", "email", "key").
		WithArgs().Raw("Benjamin", "ben@whitehouse.gov", "6").ExecContext(context.TODO())
	assert.NoError(t, err)

	// Update the key
	res, err := s.Update("dml_people").Set(Column("key").Str("6-revoked")).Where(Column("key").Str("6")).WithArgs().ExecContext(context.TODO())
	assert.NoError(t, err)

	// Assert our record was updated (and only our record)
	rowsAff, err := res.RowsAffected()
	assert.NoError(t, err)
	assert.Equal(t, int64(1), rowsAff)

	var person dmlPerson
	_, err = s.SelectFrom("dml_people").AddColumns("id", "name", "key").
		Where(Column("email").Str("ben@whitehouse.gov")).WithArgs().Load(context.TODO(), &person)
	assert.NoError(t, err)

	assert.Equal(t, "Benjamin", person.Name)
	assert.Equal(t, "6-revoked", person.Key.String)
}

func TestUpdateReal(t *testing.T) {
	s := createRealSessionWithFixtures(t, nil)
	defer testCloser(t, s)

	// Insert a George
	res, err := s.InsertInto("dml_people").AddColumns("name", "email").
		WithArgs().ExecContext(context.TODO(), "George", "george@whitehouse.gov")
	assert.NoError(t, err)

	// Get George'ab ID
	id, err := res.LastInsertId()
	assert.NoError(t, err)

	// Rename our George to Barack
	_, err = s.Update("dml_people").
		Set(Column("name").Str("Barack"), Column("email").Str("barack@whitehouse.gov")).
		Where(Column("id").In().Int64s(id, 8888)).WithArgs().ExecContext(context.TODO())
	// Meaning of 8888: Just to see if the SQL with place holders gets created correctly
	require.NoError(t, err)

	var person dmlPerson
	_, err = s.SelectFrom("dml_people").Star().Where(Column("id").Int64(id)).WithArgs().Load(context.TODO(), &person)
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
		assert.True(t, errors.Empty.Match(err))
	})

	t.Run("Prepare Error", func(t *testing.T) {
		u := &Update{}
		u.DB = dbMock{
			error: errors.AlreadyClosed.Newf("Who closed myself?"),
		}
		u.Table.Name = "tableY"
		u.Set(Column("a").Int(1))

		stmt, err := u.Prepare(context.TODO())
		assert.Nil(t, stmt)
		assert.True(t, errors.AlreadyClosed.Match(err), "%+v", err)
	})
}

func TestUpdate_ToSQL_Without_Column_Arguments(t *testing.T) {
	t.Parallel()
	t.Run("with condition values", func(t *testing.T) {
		u := NewUpdate("catalog_product_entity").AddColumns("sku", "updated_at")
		u.Where(Column("entity_id").In().Int64s(1, 2, 3))
		compareToSQL(t, u, errors.NoKind,
			"UPDATE `catalog_product_entity` SET `sku`=?, `updated_at`=? WHERE (`entity_id` IN (1,2,3))",
			"",
		)
	})
	t.Run("without condition values", func(t *testing.T) {
		u := NewUpdate("catalog_product_entity").AddColumns("sku", "updated_at")
		u.Where(Column("entity_id").In().PlaceHolder())

		compareToSQL(t, u, errors.NoKind,
			"UPDATE `catalog_product_entity` SET `sku`=?, `updated_at`=? WHERE (`entity_id` IN ?)",
			"",
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
				ListenUpdateFn: func(b *Update) {
					b.Set(Column("a").Int(1))
				},
			},
			Listen{
				Name:      "listener2",
				EventType: OnBeforeToSQL,
				ListenUpdateFn: func(b *Update) {
					b.Set(Column("b").Int(1))
					b.PropagationStopped = true
				},
			},
			Listen{
				Name:      "listener3",
				EventType: OnBeforeToSQL,
				ListenUpdateFn: func(b *Update) {
					panic("Should not get called")
				},
			},
		)
		compareToSQL(t, d, errors.NoKind,
			"UPDATE `tableA` SET `y`=25, `z`=26, `a`=1, `b`=1",
			"UPDATE `tableA` SET `y`=25, `z`=26, `a`=1, `b`=1", // each call ToSQL appends more columns
		)
	})

	t.Run("Missing EventType", func(t *testing.T) {
		up := NewUpdate("tableA")
		up.Set(Column("a").Int(1), Column("b").Bool(true))

		up.Listeners.Add(
			Listen{
				Name: "c=pi",
				ListenUpdateFn: func(u *Update) {
					u.Set(Column("c").Float64(3.14159))
				},
			},
		)
		compareToSQL(t, up, errors.Empty,
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
				EventType: OnBeforeToSQL,
				ListenUpdateFn: func(u *Update) {
					u.Set(Column("c").Float64(3.14159))
				},
			},
			Listen{
				Name:      "d=d",
				EventType: OnBeforeToSQL,
				ListenUpdateFn: func(u *Update) {
					u.Set(Column("d").Str("d"))
				},
			},
		)

		up.Listeners.Add(Listen{
			Name:      "e",
			EventType: OnBeforeToSQL,
			ListenUpdateFn: func(u *Update) {
				u.Set(Column("e").Str("e"))
			},
		})
		compareToSQL(t, up, errors.NoKind,
			"UPDATE `tableA` SET `a`=1, `b`=1, `c`=3.14159, `d`='d', `e`='e'",
			"UPDATE `tableA` SET `a`=1, `b`=1, `c`=3.14159, `d`='d', `e`='e'", // each call ToSQL appends more columns
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
		u := NewUpdate("dml_person").AddColumns("name", "email").WithArgs().Record("", pRec)
		compareToSQL(t, u, errors.NoKind,
			"UPDATE `dml_person` SET `name`=?, `email`=?",
			"UPDATE `dml_person` SET `name`='Gopher', `email`='gopher@g00gle.c0m'",
			"Gopher", "gopher@g00gle.c0m",
		)
	})
	t.Run("with where", func(t *testing.T) {
		u := NewUpdate("dml_person").AddColumns("name", "email").
			Where(Column("id").PlaceHolder()).WithArgs().Record("", pRec)
		compareToSQL(t, u, errors.NoKind,
			"UPDATE `dml_person` SET `name`=?, `email`=? WHERE (`id` = ?)",
			"UPDATE `dml_person` SET `name`='Gopher', `email`='gopher@g00gle.c0m' WHERE (`id` = 12345)",
			"Gopher", "gopher@g00gle.c0m", int64(12345),
		)
		assert.Exactly(t, []string{"name", "email", "id"}, u.base.qualifiedColumns)
	})
	t.Run("fails column `key` not in entity object", func(t *testing.T) {
		u := NewUpdate("dml_person").AddColumns("name", "email").
			Set(Column("keyXXX").PlaceHolder()).
			Where(Column("id").PlaceHolder()).
			WithArgs().Record("", pRec)
		compareToSQL(t, u, errors.NotFound,
			"",
			"",
		)
	})
}

func TestUpdate_DisableBuildCache(t *testing.T) {
	t.Parallel()

	up := NewUpdate("a").
		Set(
			Column("foo").Int(1),
			Column("bar").Expr("COALESCE(bar, 0) + ?").Int(2)).
		Where(Column("id").PlaceHolder()).
		DisableBuildCache()

	const cachedSQLPlaceHolder = "UPDATE `a` SET `foo`=1, `bar`=COALESCE(bar, 0) + 2 WHERE (`id` = ?)"
	const cachedSQLInterpolated = "UPDATE `a` SET `foo`=1, `bar`=COALESCE(bar, 0) + 2 WHERE (`id` = 987654321)"

	for i := 0; i < 3; i++ {
		compareToSQL(t, up.WithArgs().Uint(987654321), errors.NoKind,
			cachedSQLPlaceHolder,
			cachedSQLInterpolated,
			int64(987654321),
		)
		assert.Empty(t, up.cachedSQL)
	}

}
