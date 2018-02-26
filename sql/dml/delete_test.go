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
	"time"

	"github.com/corestoreio/errors"
	"github.com/corestoreio/log"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDeleteAllToSQL(t *testing.T) {
	t.Parallel()

	compareToSQL2(t, NewDelete("a"), errors.NoKind, "DELETE FROM `a`")
	compareToSQL2(t, NewDelete("a").Alias("b"), errors.NoKind, "DELETE FROM `a` AS `b`")
}

func TestDeleteSingleToSQL(t *testing.T) {
	t.Parallel()

	qb := NewDelete("a").Where(Column("id").Int(1))
	compareToSQL2(t, qb, errors.NoKind,
		"DELETE FROM `a` WHERE (`id` = 1)",
	)

	// test for being idempotent
	compareToSQL2(t, qb, errors.NoKind,
		"DELETE FROM `a` WHERE (`id` = 1)",
	)
}

func TestDelete_OrderBy(t *testing.T) {
	t.Parallel()
	t.Run("expr", func(t *testing.T) {
		compareToSQL2(t, NewDelete("a").Unsafe().OrderBy("b=c").OrderByDesc("d"), errors.NoKind,
			"DELETE FROM `a` ORDER BY b=c, `d` DESC",
		)
	})
	t.Run("asc", func(t *testing.T) {
		compareToSQL2(t, NewDelete("a").OrderBy("b").OrderBy("c"), errors.NoKind,
			"DELETE FROM `a` ORDER BY `b`, `c`",
		)
	})
	t.Run("desc", func(t *testing.T) {
		compareToSQL2(t, NewDelete("a").OrderBy("b").OrderByDesc("c").OrderBy("d").OrderByDesc("e", "f").OrderBy("g"), errors.NoKind,
			"DELETE FROM `a` ORDER BY `b`, `c` DESC, `d`, `e` DESC, `f` DESC, `g`",
		)
	})
}

func TestDelete_Limit_Offset(t *testing.T) {
	t.Parallel()
	compareToSQL2(t, NewDelete("a").Limit(10).OrderBy("id"), errors.NoKind,
		"DELETE FROM `a` ORDER BY `id` LIMIT 10",
	)
}

func TestDelete_Interpolate(t *testing.T) {
	t.Parallel()

	compareToSQL2(t, NewDelete("tableA").
		Where(
			Column("colA").GreaterOrEqual().Float64(3.14159),
			Column("colB").In().Ints(1, 2, 3, 45),
			Column("colC").Str("Hello"),
		).
		Limit(10).OrderBy("id"), errors.NoKind,
		"DELETE FROM `tableA` WHERE (`colA` >= 3.14159) AND (`colB` IN (1,2,3,45)) AND (`colC` = 'Hello') ORDER BY `id` LIMIT 10",
	)

	compareToSQL2(t, NewDelete("tableA").
		Where(
			Column("colA").GreaterOrEqual().Float64(3.14159),
			Column("colB").In().NamedArg("colB2"),
		).
		Limit(10).OrderBy("id").WithArgs().Name("colB2").Int64s(3, 4, 7, 8).Interpolate(), errors.NoKind,
		"DELETE FROM `tableA` WHERE (`colA` >= 3.14159) AND (`colB` IN (3,4,7,8)) ORDER BY `id` LIMIT 10",
	)

}

func TestDeleteReal(t *testing.T) {
	s := createRealSessionWithFixtures(t, nil)
	defer testCloser(t, s)
	// Insert a Barack
	res, err := s.InsertInto("dml_people").AddColumns("name", "email").
		WithArgs().ExecContext(context.TODO(), "Barack", "barack@whitehouse.gov")
	require.NoError(t, err)
	if res == nil {
		t.Fatal("result should not be nil. See previous error")
	}

	// Get Barack'ab ID
	id, err := res.LastInsertId()
	require.NoError(t, err, "LastInsertId")

	// Delete Barack
	res, err = s.DeleteFrom("dml_people").Where(Column("id").Int64(id)).WithArgs().ExecContext(context.TODO())
	require.NoError(t, err, "DeleteFrom")

	// Ensure we only reflected one row and that the id no longer exists
	rowsAff, err := res.RowsAffected()
	require.NoError(t, err)
	assert.Equal(t, int64(1), rowsAff, "RowsAffected")

	count, found, err := s.SelectFrom("dml_people").Count().Where(Column("id").PlaceHolder()).WithArgs().Int64(id).LoadNullInt64(context.TODO())
	require.NoError(t, err)
	require.True(t, found, "should have found a row")
	assert.Equal(t, int64(0), count.Int64, "count")
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
				ListenDeleteFn: func(b *Delete) {
					b.OrderByDesc("col1")
				},
			},
			Listen{
				Name:      "listener2",
				EventType: OnBeforeToSQL,
				ListenDeleteFn: func(b *Delete) {
					b.OrderByDesc("col2")
					b.PropagationStopped = true
				},
			},
			Listen{
				Name:      "listener3",
				EventType: OnBeforeToSQL,
				ListenDeleteFn: func(b *Delete) {
					panic("Should not get called")
				},
			},
		)
		compareToSQL2(t, d, errors.NoKind,
			"DELETE FROM `tableA` AS `main_table` ORDER BY `col1` DESC, `col2` DESC",
		)
	})

	t.Run("Missing EventType", func(t *testing.T) {
		d := NewDelete("tableA").Alias("main_table")

		d.OrderBy("col2")
		d.Listeners.Add(
			Listen{
				Name: "col1",
				ListenDeleteFn: func(b *Delete) {
					b.OrderByDesc("col1")
				},
			},
		)
		compareToSQL2(t, d, errors.Empty,
			"",
		)
	})

	t.Run("Should Dispatch", func(t *testing.T) {

		d := NewDelete("tableA").Alias("main_table")

		d.OrderBy("col2")
		d.Listeners.Add(
			Listen{
				Name:      "col1",
				EventType: OnBeforeToSQL,
				ListenDeleteFn: func(b *Delete) {
					b.OrderByDesc("col1")
				},
			},
			Listen{
				Name:      "storeid",
				EventType: OnBeforeToSQL,
				ListenDeleteFn: func(b *Delete) {
					b.Where(Column("store_id").Int64(1))
				},
			},
		)

		d.Listeners.Add(
			Listen{
				Name:      "repetitive",
				EventType: OnBeforeToSQL,
				ListenDeleteFn: func(b *Delete) {
					b.Where(Column("repetitive").Int(3))
				},
			},
		)
		compareToSQL2(t, d, errors.NoKind,
			"DELETE FROM `tableA` AS `main_table` WHERE (`store_id` = 1) AND (`repetitive` = 3) ORDER BY `col2`, `col1` DESC",
		)
		assert.Exactly(t, `col1; storeid; repetitive`, d.Listeners.String())
	})
}

func TestDelete_BuildCacheDisabled(t *testing.T) {
	t.Parallel()

	del := NewDelete("alpha").Where(
		Column("a").Str("b"),
		Column("b").PlaceHolder(),
	).Limit(1).OrderBy("id")

	del.IsBuildCacheDisabled = true

	const iterations = 3
	const cachedSQLPlaceHolder = "DELETE FROM `alpha` WHERE (`a` = 'b') AND (`b` = ?) ORDER BY `id` LIMIT 1"
	t.Run("without interpolate", func(t *testing.T) {
		for i := 0; i < iterations; i++ {
			sql, args, err := del.ToSQL()
			require.NoError(t, err, "%+v", err)
			require.Equal(t, cachedSQLPlaceHolder, sql)
			assert.Nil(t, args, "No arguments provided but got some")
			assert.Nil(t, del.cachedSQL, "cache []byte should be nil")
		}
	})

	t.Run("with interpolate", func(t *testing.T) {
		delA := del.WithArgs().Int64(3333).Interpolate()
		del.cachedSQL = nil
		const cachedSQLInterpolated = "DELETE FROM `alpha` WHERE (`a` = 'b') AND (`b` = 3333) ORDER BY `id` LIMIT 1"
		for i := 0; i < iterations; i++ {
			compareToSQL2(t, delA, errors.NoKind, cachedSQLInterpolated)
			assert.Nil(t, del.cachedSQL, "cache []byte should be nil")
		}
	})
}

func TestDelete_Bind(t *testing.T) {
	t.Parallel()
	p := &dmlPerson{
		ID:    5555,
		Email: MakeNullString("hans@wurst.com"),
	}
	t.Run("multiple args from Record", func(t *testing.T) {
		del := NewDelete("dml_people").
			Where(
				Column("idI64").Greater().Int64(4),
				Column("id").Equal().PlaceHolder(),
				Column("float64_pi").Float64(3.14159),
				Column("email").PlaceHolder(),
				Column("int_e").Int(2718281),
			).OrderBy("id").
			WithArgs().Records(Qualify("", p))

		compareToSQL2(t, del, errors.NoKind,
			"DELETE FROM `dml_people` WHERE (`idI64` > 4) AND (`id` = ?) AND (`float64_pi` = 3.14159) AND (`email` = ?) AND (`int_e` = 2718281) ORDER BY `id`",
			int64(5555), "hans@wurst.com",
		)
	})
	t.Run("single arg from Record unqualified", func(t *testing.T) {
		del := NewDelete("dml_people").
			Where(
				Column("id").PlaceHolder(),
			).OrderBy("id").
			WithArgs().Records(Qualify("", p))

		compareToSQL2(t, del, errors.NoKind,
			"DELETE FROM `dml_people` WHERE (`id` = ?) ORDER BY `id`",
			int64(5555),
		)
		assert.Exactly(t, []string{"id"}, del.base.qualifiedColumns)
	})
	t.Run("single arg from Record qualified", func(t *testing.T) {
		del := NewDelete("dml_people").Alias("dmlPpl").
			Where(
				Column("id").PlaceHolder(),
			).OrderBy("id").
			WithArgs().Records(Qualify("dmlPpl", p))

		compareToSQL(t, del, errors.NoKind,
			"DELETE FROM `dml_people` AS `dmlPpl` WHERE (`id` = ?) ORDER BY `id`",
			"DELETE FROM `dml_people` AS `dmlPpl` WHERE (`id` = 5555) ORDER BY `id`",
			int64(5555),
		)
		assert.Exactly(t, []string{"id"}, del.base.qualifiedColumns)
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
			).OrderBy("id").WithArgs().Record("", ntr)

		compareToSQL2(t, del, errors.NoKind,
			"DELETE FROM `null_type_table` WHERE (`string_val` = ?) AND (`int64_val` = ?) AND (`float64_val` = ?) AND (`random1` BETWEEN 1.2 AND 3.4) AND (`time_val` = ?) AND (`bool_val` = ?) ORDER BY `id`",
			"wow", int64(42), 1.618, time.Date(2009, 1, 3, 18, 15, 5, 0, time.UTC), true,
		)
	})
}

func TestDelete_Clone(t *testing.T) {
	t.Parallel()

	t.Run("nil", func(t *testing.T) {
		var d *Delete
		d2 := d.Clone()
		assert.Nil(t, d)
		assert.Nil(t, d2)
	})
	t.Run("non-nil", func(t *testing.T) {
		d := NewDelete("dml_people").Alias("dmlPpl").FromTables("a1", "b2").
			Where(
				Column("id").PlaceHolder(),
			).OrderBy("id")
		d2 := d.Clone()
		notEqualPointers(t, d, d2)
		notEqualPointers(t, d, d2)
		notEqualPointers(t, d.BuilderConditional.Wheres, d2.BuilderConditional.Wheres)
		notEqualPointers(t, d.MultiTables, d2.MultiTables)
	})
}
