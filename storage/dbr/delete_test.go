package dbr

import (
	"testing"

	"github.com/corestoreio/csfw/util/errors"
	"github.com/stretchr/testify/assert"
)

func BenchmarkDeleteSql(b *testing.B) {
	s := createFakeSession()

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_, _, err := s.DeleteFrom("alpha").Where(ConditionRaw("a", "b")).Limit(1).OrderDir("id", true).ToSQL()
		if err != nil {
			b.Fatalf("%+v", err)
		}
	}
}

func TestDeleteAllToSql(t *testing.T) {
	s := createFakeSession()

	sql, _, err := s.DeleteFrom("a").ToSQL()
	assert.NoError(t, err)
	assert.Equal(t, sql, "DELETE FROM `a`")

	sql, _, err = s.DeleteFrom("a", "b").ToSQL()
	assert.NoError(t, err)
	assert.Equal(t, sql, "DELETE FROM `a` AS `b`")
}

func TestDeleteSingleToSql(t *testing.T) {
	s := createFakeSession()

	del := s.DeleteFrom("a").Where(ConditionRaw("id = ?", 1))
	sql, args, err := del.ToSQL()
	assert.NoError(t, err)
	assert.Equal(t, sql, "DELETE FROM `a` WHERE (id = ?)")
	assert.Equal(t, args, []interface{}{1})

	// once where was a sync.Pool for the whereFragments with which it was
	// not possible to run ToSQL() twice.
	sql, args, err = del.ToSQL()
	assert.NoError(t, err)
	assert.Equal(t, sql, "DELETE FROM `a` WHERE (id = ?)")
	assert.Equal(t, args, []interface{}{1})

}

func TestDeleteTenStaringFromTwentyToSql(t *testing.T) {
	s := createFakeSession()

	sql, _, err := s.DeleteFrom("a").Limit(10).Offset(20).OrderBy("id").ToSQL()
	assert.NoError(t, err)
	assert.Equal(t, sql, "DELETE FROM `a` ORDER BY id LIMIT 10 OFFSET 20")
}

func TestDeleteReal(t *testing.T) {
	s := createRealSessionWithFixtures()

	// Insert a Barack
	res, err := s.InsertInto("dbr_people").Columns("name", "email").Values("Barack", "barack@whitehouse.gov").Exec()
	assert.NoError(t, err)

	// Get Barack's ID
	id, err := res.LastInsertId()
	assert.NoError(t, err, "LastInsertId")

	// Delete Barack
	res, err = s.DeleteFrom("dbr_people").Where(ConditionRaw("id = ?", id)).Exec()
	assert.NoError(t, err, "DeleteFrom")

	// Ensure we only reflected one row and that the id no longer exists
	rowsAff, err := res.RowsAffected()
	assert.NoError(t, err)
	assert.Equal(t, rowsAff, int64(1), "RowsAffected")

	var count int64
	err = s.Select("count(*)").From("dbr_people").Where(ConditionRaw("id = ?", id)).LoadValue(&count)
	assert.NoError(t, err)
	assert.Equal(t, count, int64(0), "count")
}

func TestDelete_Prepare(t *testing.T) {

	t.Run("ToSQL Error", func(t *testing.T) {
		d := &Delete{}
		d.Where(ConditionRaw("a", 1))
		stmt, err := d.Prepare()
		assert.Nil(t, stmt)
		assert.True(t, errors.IsEmpty(err))
	})

	t.Run("Prepare Error", func(t *testing.T) {
		d := &Delete{
			From: MakeAlias("table"),
			Preparer: dbMock{
				error: errors.NewAlreadyClosedf("Who closed myself?"),
			},
		}
		d.Where(ConditionRaw("a", 1))
		stmt, err := d.Prepare()
		assert.Nil(t, stmt)
		assert.True(t, errors.IsAlreadyClosed(err), "%+v", err)
	})
}

func TestDelete_AddHookBeforeToSQLOnce(t *testing.T) {
	d := NewDelete("tableA", "main_table")

	d.OrderBy("col2")
	d.AddHookBeforeToSQLOnce(func(s2 *Delete) {
		s2.OrderDir("col1", false)
	})

	sql, args, err := d.ToSQL()
	assert.NoError(t, err)
	assert.Nil(t, args)
	assert.NotEmpty(t, sql)

	sql, args, err = d.ToSQL()
	assert.NoError(t, err)
	assert.Nil(t, args)
	assert.Exactly(t, "DELETE FROM `tableA` AS `main_table` ORDER BY col2, col1 DESC", sql)
}
