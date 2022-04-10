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
	"database/sql"
	"testing"
	"time"

	"github.com/corestoreio/pkg/storage/null"
	"github.com/corestoreio/pkg/util/assert"
	"github.com/corestoreio/pkg/util/conv"
)

func TestDeleteAllToSQL(t *testing.T) {
	compareToSQL2(t, NewDelete("a"), false, "DELETE FROM `a`")
	compareToSQL2(t, NewDelete("a").Alias("b"), false, "DELETE FROM `a` AS `b`")
}

func TestDeleteSingleToSQL(t *testing.T) {
	qb := NewDelete("a").Where(Column("id").Int(1))
	compareToSQL2(t, qb, false,
		"DELETE FROM `a` WHERE (`id` = 1)",
	)

	// test for being idempotent
	compareToSQL2(t, qb, false,
		"DELETE FROM `a` WHERE (`id` = 1)",
	)
}

func TestDelete_OrderBy(t *testing.T) {
	t.Run("expr", func(t *testing.T) {
		compareToSQL2(t, NewDelete("a").Unsafe().OrderBy("b=c").OrderByDesc("d"), false,
			"DELETE FROM `a` ORDER BY b=c, `d` DESC",
		)
	})
	t.Run("asc", func(t *testing.T) {
		compareToSQL2(t, NewDelete("a").OrderBy("b").OrderBy("c"), false,
			"DELETE FROM `a` ORDER BY `b`, `c`",
		)
	})
	t.Run("desc", func(t *testing.T) {
		compareToSQL2(t, NewDelete("a").OrderBy("b").OrderByDesc("c").OrderBy("d").OrderByDesc("e", "f").OrderBy("g"), false,
			"DELETE FROM `a` ORDER BY `b`, `c` DESC, `d`, `e` DESC, `f` DESC, `g`",
		)
	})
}

func TestDelete_Limit_Offset(t *testing.T) {
	compareToSQL2(t, NewDelete("a").Limit(10).OrderBy("id"), false,
		"DELETE FROM `a` ORDER BY `id` LIMIT 10",
	)
}

func TestDelete_Interpolate(t *testing.T) {
	compareToSQL2(t, NewDelete("tableA").
		Where(
			Column("colA").GreaterOrEqual().Float64(3.14159),
			Column("colB").In().Ints(1, 2, 3, 45),
			Column("colC").Str("Hello"),
		).
		Limit(10).OrderBy("id"), false,
		"DELETE FROM `tableA` WHERE (`colA` >= 3.14159) AND (`colB` IN (1,2,3,45)) AND (`colC` = 'Hello') ORDER BY `id` LIMIT 10",
	)

	compareToSQL(t, NewDelete("tableA").
		Where(
			Column("colA").GreaterOrEqual().Float64(3.14159),
			Column("colB").In().NamedArg("colB2"),
		).
		Limit(10).OrderBy("id").WithDBR(nil).Interpolate().TestWithArgs(sql.Named("colB2", []int64{3, 4, 7, 8})), false,
		"DELETE FROM `tableA` WHERE (`colA` >= 3.14159) AND (`colB` IN ?) ORDER BY `id` LIMIT 10",
		"DELETE FROM `tableA` WHERE (`colA` >= 3.14159) AND (`colB` IN (3,4,7,8)) ORDER BY `id` LIMIT 10",
		int64(3), int64(4), int64(7), int64(8),
	)
}

func TestDeleteReal(t *testing.T) {
	s := createRealSessionWithFixtures(t, nil)
	defer testCloser(t, s)
	// Insert a Barack
	res, err := s.WithQueryBuilder(NewInsert("dml_people").AddColumns("name", "email")).
		ExecContext(context.TODO(), "Barack", "barack@whitehouse.gov")
	assert.NoError(t, err)
	assert.NotNil(t, res)

	// Get Barack'ab ID
	id, err := res.LastInsertId()
	assert.NoError(t, err, "LastInsertId")
	// Delete Barack
	res, err = s.WithQueryBuilder(NewDelete("dml_people").Where(Column("id").Int64(id))).ExecContext(context.TODO())
	assert.NoError(t, err, "DeleteFrom")

	// Ensure we only reflected one row and that the id no longer exists
	rowsAff, err := res.RowsAffected()
	assert.NoError(t, err)
	assert.Exactly(t, int64(1), rowsAff, "RowsAffected")
	count, found, err := s.WithQueryBuilder(NewSelect().Count().From("dml_people").Where(Column("id").PlaceHolder())).LoadNullInt64(context.TODO(), id)
	assert.NoError(t, err)
	assert.True(t, found, "should have found a row")
	assert.Exactly(t, int64(0), count.Int64, "count")

	assert.Exactly(t, []string{
		"DELETEf0d5f12cb21e7142", "DELETE FROM `dml_people` WHERE (`id` = 3)",
		"INSERT35b056eb359cf3ef", "INSERT INTO `dml_people` (`name`,`email`) VALUES ",
		"SELECT835edcaab6dec8d3", "SELECT COUNT(*) AS `counted` FROM `dml_people` WHERE (`id` = ?)",
	}, conv.ToStringSlice(s.CachedQueries()))
}

func TestDelete_BuildCacheDisabled(t *testing.T) {
	del := NewDelete("alpha").Where(
		Column("a").Str("b"),
		Column("b").PlaceHolder(),
	).Limit(1).OrderBy("id")

	const iterations = 3
	const cachedSQLPlaceHolder = "DELETE FROM `alpha` WHERE (`a` = 'b') AND (`b` = ?) ORDER BY `id` LIMIT 1"
	t.Run("without interpolate", func(t *testing.T) {
		for i := 0; i < iterations; i++ {
			sql, args, err := del.ToSQL()
			assert.NoError(t, err)
			assert.Exactly(t, cachedSQLPlaceHolder, sql)
			assert.Nil(t, args, "No arguments provided but got some")
		}
		// assert.Exactly(t, []string{"", "DELETE FROM `alpha` WHERE (`a` = 'b') AND (`b` = ?) ORDER BY `id` LIMIT 1"},
		//	del.CachedQueries())
	})

	t.Run("with interpolate", func(t *testing.T) {
		delA := del.WithDBR(nil)

		compareToSQL(t,
			delA.TestWithArgs(123),
			false,
			"DELETE FROM `alpha` WHERE (`a` = 'b') AND (`b` = ?) ORDER BY `id` LIMIT 1",
			"DELETE FROM `alpha` WHERE (`a` = 'b') AND (`b` = 123) ORDER BY `id` LIMIT 1",
			int64(123),
		)
		delA.Reset()
		compareToSQL(t,
			delA.TestWithArgs(124),
			false,
			"DELETE FROM `alpha` WHERE (`a` = 'b') AND (`b` = ?) ORDER BY `id` LIMIT 1",
			"DELETE FROM `alpha` WHERE (`a` = 'b') AND (`b` = 124) ORDER BY `id` LIMIT 1",
			int64(124),
		)

		// assert.Exactly(t, []string{"", "DELETE FROM `alpha` WHERE (`a` = 'b') AND (`b` = ?) ORDER BY `id` LIMIT 1"},
		//	del.CachedQueries())
	})
}

func TestDelete_Bind(t *testing.T) {
	p := &dmlPerson{
		ID:    5555,
		Email: null.MakeString("hans@wurst.com"),
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
			WithDBR(nil).TestWithArgs(Qualify("", p))

		compareToSQL2(t, del, false,
			"DELETE FROM `dml_people` WHERE (`idI64` > 4) AND (`id` = ?) AND (`float64_pi` = 3.14159) AND (`email` = ?) AND (`int_e` = 2718281) ORDER BY `id`",
			int64(5555), "hans@wurst.com",
		)
	})
	t.Run("single arg from Record unqualified", func(t *testing.T) {
		del := NewDelete("dml_people").
			Where(
				Column("id").PlaceHolder(),
			).OrderBy("id").
			WithDBR(nil)

		compareToSQL2(t, del.TestWithArgs(Qualify("", p)), false,
			"DELETE FROM `dml_people` WHERE (`id` = ?) ORDER BY `id`",
			int64(5555),
		)
		assert.Exactly(t, []string{"id"}, del.cachedSQL.qualifiedColumns)
	})
	t.Run("single arg from Record qualified", func(t *testing.T) {
		del := NewDelete("dml_people").Alias("dmlPpl").
			Where(
				Column("id").PlaceHolder(),
			).OrderBy("id").
			WithDBR(nil)

		compareToSQL(t, del.TestWithArgs(Qualify("dmlPpl", p)), false,
			"DELETE FROM `dml_people` AS `dmlPpl` WHERE (`id` = ?) ORDER BY `id`",
			"DELETE FROM `dml_people` AS `dmlPpl` WHERE (`id` = 5555) ORDER BY `id`",
			int64(5555),
		)
		assert.Exactly(t, []string{"id"}, del.cachedSQL.qualifiedColumns)
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
			).OrderBy("id").WithDBR(nil).TestWithArgs(Qualify("", ntr))

		compareToSQL2(t, del, false,
			"DELETE FROM `null_type_table` WHERE (`string_val` = ?) AND (`int64_val` = ?) AND (`float64_val` = ?) AND (`random1` BETWEEN 1.2 AND 3.4) AND (`time_val` = ?) AND (`bool_val` = ?) ORDER BY `id`",
			"wow", int64(42), 1.618, time.Date(2009, 1, 3, 18, 15, 5, 0, time.UTC), true,
		)
	})
}
