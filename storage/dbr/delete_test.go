package dbr

import (
	"context"
	"testing"

	"github.com/corestoreio/errors"
	"github.com/corestoreio/log"
	"github.com/stretchr/testify/assert"
)

var benchmarkDeleteSQL Arguments

func BenchmarkDeleteSQL(b *testing.B) {
	s := createFakeSession()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		var err error
		_, benchmarkDeleteSQL, err = s.DeleteFrom("alpha").Where(Condition("a", ArgString("b"))).Limit(1).OrderBy("id").ToSQL()
		if err != nil {
			b.Fatalf("%+v", err)
		}
	}
}

func TestDeleteAllToSQL(t *testing.T) {
	s := createFakeSession()

	sql, _, err := s.DeleteFrom("a").ToSQL()
	assert.NoError(t, err)
	assert.Equal(t, sql, "DELETE FROM `a`")

	sql, _, err = s.DeleteFrom("a", "b").ToSQL()
	assert.NoError(t, err)
	assert.Equal(t, sql, "DELETE FROM `a` AS `b`")
}

func TestDeleteSingleToSQL(t *testing.T) {
	s := createFakeSession()

	del := s.DeleteFrom("a").Where(Condition("id = ?", argInt(1)))
	sql, args, err := del.ToSQL()
	assert.NoError(t, err)
	assert.Equal(t, sql, "DELETE FROM `a` WHERE (id = ?)")
	assert.Equal(t, []interface{}{int64(1)}, args.Interfaces())

	// once where was a sync.Pool for the whereFragments with which it was
	// not possible to run ToSQL() twice.
	sql, args, err = del.ToSQL()
	assert.NoError(t, err)
	assert.Equal(t, sql, "DELETE FROM `a` WHERE (id = ?)")
	assert.Equal(t, []interface{}{int64(1)}, args.Interfaces())

}

func TestDeleteTenStaringFromTwentyToSQL(t *testing.T) {
	s := createFakeSession()

	sql, _, err := s.DeleteFrom("a").Limit(10).Offset(20).OrderBy("id").ToSQL()
	assert.NoError(t, err)
	assert.Equal(t, sql, "DELETE FROM `a` ORDER BY id LIMIT 10 OFFSET 20")
}

func TestDeleteReal(t *testing.T) {
	s := createRealSessionWithFixtures()

	// Insert a Barack
	res, err := s.InsertInto("dbr_people").AddColumns("name", "email").
		AddValues(ArgString("Barack"), ArgString("barack@whitehouse.gov")).Exec(context.TODO())
	assert.NoError(t, err)
	if res == nil {
		t.Fatal("result should not be nil. See previous error")
	}

	// Get Barack'ab ID
	id, err := res.LastInsertId()
	assert.NoError(t, err, "LastInsertId")

	// Delete Barack
	res, err = s.DeleteFrom("dbr_people").Where(Condition("id", ArgInt64(id))).Exec(context.TODO())
	assert.NoError(t, err, "DeleteFrom")

	// Ensure we only reflected one row and that the id no longer exists
	rowsAff, err := res.RowsAffected()
	assert.NoError(t, err)
	assert.Equal(t, rowsAff, int64(1), "RowsAffected")

	var count int64
	err = s.Select("count(*)").From("dbr_people").Where(Condition("id", ArgInt64(id))).LoadValue(context.TODO(), &count)
	assert.NoError(t, err)
	assert.Equal(t, count, int64(0), "count")
}

func TestDelete_Prepare(t *testing.T) {

	t.Run("ToSQL Error", func(t *testing.T) {
		d := &Delete{}
		d.Where(Condition("a", argInt64(1)))
		stmt, err := d.Prepare(context.TODO())
		assert.Nil(t, stmt)
		assert.True(t, errors.IsEmpty(err))
	})

	t.Run("Prepare Error", func(t *testing.T) {
		d := &Delete{
			From: MakeAlias("table"),
		}
		d.DB.Preparer = dbMock{
			error: errors.NewAlreadyClosedf("Who closed myself?"),
		}
		d.Where(Condition("a", argInt(1)))
		stmt, err := d.Prepare(context.TODO())
		assert.Nil(t, stmt)
		assert.True(t, errors.IsAlreadyClosed(err), "%+v", err)
	})
}

func TestDelete_Events(t *testing.T) {
	t.Parallel()

	t.Run("Stop Propagation", func(t *testing.T) {
		d := NewDelete("tableA", "main_table")
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
		sql, _, err := d.ToSQL()
		assert.NoError(t, err, "%+v", err)
		assert.Exactly(t, "DELETE FROM `tableA` AS `main_table` ORDER BY col1 DESC, col2 DESC", sql)

		sql, _, err = d.ToSQL()
		assert.NoError(t, err, "%+v", err)
		assert.Exactly(t, "DELETE FROM `tableA` AS `main_table` ORDER BY col1 DESC, col2 DESC, col1 DESC, col2 DESC", sql)

	})

	t.Run("Missing EventType", func(t *testing.T) {
		d := NewDelete("tableA", "main_table")

		d.OrderBy("col2")
		d.Listeners.Add(
			Listen{
				Name: "col1",
				DeleteFunc: func(b *Delete) {
					b.OrderByDesc("col1")
				},
			},
		)

		sql, args, err := d.ToSQL()
		assert.Empty(t, sql)
		assert.Nil(t, args)
		assert.True(t, errors.IsEmpty(err), "%+v", err)
	})

	t.Run("Should Dispatch", func(t *testing.T) {

		d := NewDelete("tableA", "main_table")

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
					b.Where(Condition("store_id=?", ArgInt64(1)))
				},
			},
		)

		d.Listeners.Add(
			Listen{
				Name:      "repetitive",
				EventType: OnBeforeToSQL,
				DeleteFunc: func(b *Delete) {
					b.Where(Condition("repetitive=?", argInt(3)))
				},
			},
		)

		sql, args, err := d.ToSQL()
		assert.NoError(t, err)
		assert.Exactly(t, []interface{}{int64(1), int64(3)}, args.Interfaces())
		assert.Exactly(t, "DELETE FROM `tableA` AS `main_table` WHERE (store_id=?) AND (repetitive=?) ORDER BY col2, col1 DESC", sql)

		sql, args, err = d.ToSQL()
		assert.NoError(t, err)
		assert.Exactly(t, []interface{}{int64(1), int64(3), int64(3)}, args.Interfaces())
		assert.Exactly(t, "DELETE FROM `tableA` AS `main_table` WHERE (store_id=?) AND (repetitive=?) AND (repetitive=?) ORDER BY col2, col1 DESC", sql)

		assert.Exactly(t, `col1; storeid; repetitive`, d.Listeners.String())
	})

}
