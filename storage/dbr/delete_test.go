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
	"time"

	"github.com/corestoreio/errors"
	"github.com/corestoreio/log"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDeleteAllToSQL(t *testing.T) {
	t.Parallel()

	compareToSQL(t, NewDelete("a"), nil, "DELETE FROM `a`", "DELETE FROM `a`")
	compareToSQL(t, NewDelete("a").Alias("b"), nil, "DELETE FROM `a` AS `b`", "DELETE FROM `a` AS `b`")
}

func TestDeleteSingleToSQL(t *testing.T) {
	t.Parallel()

	qb := NewDelete("a").Where(Column("id").Int(1))
	compareToSQL(t, qb, nil,
		"DELETE FROM `a` WHERE (`id` = ?)",
		"DELETE FROM `a` WHERE (`id` = 1)",
		int64(1),
	)

	// test for being idempotent
	compareToSQL(t, qb, nil,
		"DELETE FROM `a` WHERE (`id` = ?)",
		"DELETE FROM `a` WHERE (`id` = 1)",
		int64(1),
	)
}

func TestDelete_OrderBy(t *testing.T) {
	t.Parallel()
	t.Run("expr", func(t *testing.T) {
		compareToSQL(t, NewDelete("a").Unsafe().OrderBy("b=c").OrderByDesc("d"), nil,
			"DELETE FROM `a` ORDER BY b=c, `d` DESC",
			"DELETE FROM `a` ORDER BY b=c, `d` DESC",
		)
	})
	t.Run("asc", func(t *testing.T) {
		compareToSQL(t, NewDelete("a").OrderBy("b").OrderBy("c"), nil,
			"DELETE FROM `a` ORDER BY `b`, `c`",
			"DELETE FROM `a` ORDER BY `b`, `c`",
		)
	})
	t.Run("desc", func(t *testing.T) {
		compareToSQL(t, NewDelete("a").OrderBy("b").OrderByDesc("c").OrderBy("d").OrderByDesc("e", "f").OrderBy("g"), nil,
			"DELETE FROM `a` ORDER BY `b`, `c` DESC, `d`, `e` DESC, `f` DESC, `g`",
			"DELETE FROM `a` ORDER BY `b`, `c` DESC, `d`, `e` DESC, `f` DESC, `g`",
		)
	})
}

func TestDelete_Limit_Offset(t *testing.T) {
	t.Parallel()
	compareToSQL(t, NewDelete("a").Limit(10).OrderBy("id"), nil,
		"DELETE FROM `a` ORDER BY `id` LIMIT 10",
		"DELETE FROM `a` ORDER BY `id` LIMIT 10",
	)
}

func TestDelete_Interpolate(t *testing.T) {
	t.Parallel()
	compareToSQL(t, NewDelete("tableA").
		Where(
			Column("colA").GreaterOrEqual().Float64(3.14159),
			Column("colB").In().Ints(1, 2, 3, 45),
			Column("colC").Str("He'l`lo"),
		).
		Limit(10).OrderBy("id"), nil,
		"DELETE FROM `tableA` WHERE (`colA` >= ?) AND (`colB` IN (?,?,?,?)) AND (`colC` = ?) ORDER BY `id` LIMIT 10",
		"DELETE FROM `tableA` WHERE (`colA` >= 3.14159) AND (`colB` IN (1,2,3,45)) AND (`colC` = 'He\\'l`lo') ORDER BY `id` LIMIT 10",
		3.14159, int64(1), int64(2), int64(3), int64(45), "He'l`lo",
	)

}

func TestDeleteReal(t *testing.T) {
	s := createRealSessionWithFixtures(t, nil)

	// Insert a Barack
	res, err := s.InsertInto("dbr_people").AddColumns("name", "email").
		AddValues("Barack", "barack@whitehouse.gov").Exec(context.TODO())
	require.NoError(t, err)
	if res == nil {
		t.Fatal("result should not be nil. See previous error")
	}

	// Get Barack'ab ID
	id, err := res.LastInsertId()
	require.NoError(t, err, "LastInsertId")

	// Delete Barack
	res, err = s.DeleteFrom("dbr_people").Where(Column("id").Int64(id)).Exec(context.TODO())
	require.NoError(t, err, "DeleteFrom")

	// Ensure we only reflected one row and that the id no longer exists
	rowsAff, err := res.RowsAffected()
	require.NoError(t, err)
	assert.Equal(t, int64(1), rowsAff, "RowsAffected")

	count, err := s.SelectFrom("dbr_people").Count().Where(Column("id").Int64(id)).LoadInt64(context.TODO())
	require.NoError(t, err)
	assert.Equal(t, int64(0), count, "count")
}

func TestDelete_Events(t *testing.T) {
	t.Parallel()

	t.Run("Stop Propagation", func(t *testing.T) {
		d := NewDelete("tableA").Alias("main_table")
		d.Log = log.BlackHole{EnableInfo: true, EnableDebug: true}
		d.Listeners.Add(
			Listen{
				Name:      "listener1",
				EventType: OnBeforeToSQL,
				DeleteFunc: func(b *Delete) {
					b.OrderByDesc("col1")
				},
			},
			Listen{
				Name:      "listener2",
				EventType: OnBeforeToSQL,
				DeleteFunc: func(b *Delete) {
					b.OrderByDesc("col2")
					b.PropagationStopped = true
				},
			},
			Listen{
				Name:      "listener3",
				EventType: OnBeforeToSQL,
				DeleteFunc: func(b *Delete) {
					panic("Should not get called")
				},
			},
		)
		compareToSQL(t, d, nil,
			"DELETE FROM `tableA` AS `main_table` ORDER BY `col1` DESC, `col2` DESC",
			"DELETE FROM `tableA` AS `main_table` ORDER BY `col1` DESC, `col2` DESC, `col1` DESC, `col2` DESC",
		)
	})

	t.Run("Missing EventType", func(t *testing.T) {
		d := NewDelete("tableA").Alias("main_table")

		d.OrderBy("col2")
		d.Listeners.Add(
			Listen{
				Name: "col1",
				DeleteFunc: func(b *Delete) {
					b.OrderByDesc("col1")
				},
			},
		)
		compareToSQL(t, d, errors.IsEmpty,
			"",
			"",
		)
	})

	t.Run("Should Dispatch", func(t *testing.T) {

		d := NewDelete("tableA").Alias("main_table")

		d.OrderBy("col2")
		d.Listeners.Add(
			Listen{
				Name:      "col1",
				Once:      true,
				EventType: OnBeforeToSQL,
				DeleteFunc: func(b *Delete) {
					b.OrderByDesc("col1")
				},
			},
			Listen{
				Name:      "storeid",
				Once:      true,
				EventType: OnBeforeToSQL,
				DeleteFunc: func(b *Delete) {
					b.Where(Column("store_id").Int64(1))
				},
			},
		)

		d.Listeners.Add(
			Listen{
				Name:      "repetitive",
				EventType: OnBeforeToSQL,
				DeleteFunc: func(b *Delete) {
					b.Where(Column("repetitive").Int(3))
				},
			},
		)
		compareToSQL(t, d, nil,
			"DELETE FROM `tableA` AS `main_table` WHERE (`store_id` = ?) AND (`repetitive` = ?) ORDER BY `col2`, `col1` DESC",
			"DELETE FROM `tableA` AS `main_table` WHERE (`store_id` = 1) AND (`repetitive` = 3) AND (`repetitive` = 3) ORDER BY `col2`, `col1` DESC",
			int64(1), int64(3),
		)
		assert.Exactly(t, `col1; storeid; repetitive`, d.Listeners.String())
	})
}

func TestDelete_UseBuildCache(t *testing.T) {
	t.Parallel()

	del := NewDelete("alpha").Where(Column("a").Str("b")).Limit(1).OrderBy("id")
	del.IsBuildCache = true

	const cachedSQLPlaceHolder = "DELETE FROM `alpha` WHERE (`a` = ?) ORDER BY `id` LIMIT 1"
	t.Run("without interpolate", func(t *testing.T) {
		for i := 0; i < 3; i++ {
			sql, args, err := del.ToSQL()
			require.NoError(t, err, "%+v", err)
			require.Equal(t, cachedSQLPlaceHolder, sql)
			assert.Equal(t, []interface{}{"b"}, args)
			assert.Equal(t, cachedSQLPlaceHolder, string(del.cacheSQL))
		}
	})

	t.Run("with interpolate", func(t *testing.T) {
		del.Interpolate()
		del.cacheSQL = nil

		const cachedSQLInterpolated = "DELETE FROM `alpha` WHERE (`a` = 'b') ORDER BY `id` LIMIT 1"
		for i := 0; i < 3; i++ {
			sql, args, err := del.ToSQL()
			assert.Equal(t, cachedSQLPlaceHolder, string(del.cacheSQL))
			require.NoError(t, err, "%+v", err)
			require.Equal(t, cachedSQLInterpolated, sql)
			assert.Nil(t, args)
		}
	})
}

func TestDelete_Bind(t *testing.T) {
	t.Parallel()
	p := &dbrPerson{
		ID:    5555,
		Email: MakeNullString("hans@wurst.com"),
	}
	t.Run("multiple args from Record", func(t *testing.T) {
		del := NewDelete("dbr_people").
			Where(
				Column("idI64").Greater().Int64(4),
				Column("id").Equal().PlaceHolder(),
				Column("float64_pi").Float64(3.14159),
				Column("email").PlaceHolder(),
				Column("int_e").Int(2718281),
			).
			BindRecord(Qualify("", p)).OrderBy("id")

		compareToSQL(t, del, nil,
			"DELETE FROM `dbr_people` WHERE (`idI64` > ?) AND (`id` = ?) AND (`float64_pi` = ?) AND (`email` = ?) AND (`int_e` = ?) ORDER BY `id`",
			"DELETE FROM `dbr_people` WHERE (`idI64` > 4) AND (`id` = 5555) AND (`float64_pi` = 3.14159) AND (`email` = 'hans@wurst.com') AND (`int_e` = 2718281) ORDER BY `id`",
			int64(4), int64(5555), 3.14159, "hans@wurst.com", int64(2718281),
		)
	})
	t.Run("single arg from Record", func(t *testing.T) {
		del := NewDelete("dbr_people").
			Where(
				Column("id").PlaceHolder(),
			).
			BindRecord(Qualify("dbr_people", p)).OrderBy("id")

		compareToSQL(t, del, nil,
			"DELETE FROM `dbr_people` WHERE (`id` = ?) ORDER BY `id`",
			"DELETE FROM `dbr_people` WHERE (`id` = 5555) ORDER BY `id`",
			int64(5555),
		)
	})
	t.Run("null type records", func(t *testing.T) {
		ntr := newNullTypedRecordWithData()

		del := NewDelete("null_type_table").
			Where(
				Column("string_val").PlaceHolder(),
				Column("int64_val").PlaceHolder(),
				Column("float64_val").PlaceHolder(),
				Column("random1").Between().Float64s(1.2, 3.4),
				Column("time_val").PlaceHolder(),
				Column("bool_val").PlaceHolder(),
			).
			BindRecord(Qualify("", ntr)).OrderBy("id")

		compareToSQL(t, del, nil,
			"DELETE FROM `null_type_table` WHERE (`string_val` = ?) AND (`int64_val` = ?) AND (`float64_val` = ?) AND (`random1` BETWEEN ? AND ?) AND (`time_val` = ?) AND (`bool_val` = ?) ORDER BY `id`",
			"DELETE FROM `null_type_table` WHERE (`string_val` = 'wow') AND (`int64_val` = 42) AND (`float64_val` = 1.618) AND (`random1` BETWEEN 1.2 AND 3.4) AND (`time_val` = '2009-01-03 18:15:05') AND (`bool_val` = 1) ORDER BY `id`",
			"wow", int64(42), 1.618, 1.2, 3.4, time.Date(2009, 1, 3, 18, 15, 5, 0, time.UTC), true,
		)
	})
}
