package dbr

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func BenchmarkDeleteSql(b *testing.B) {
	s := createFakeSession()

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_, _, err := s.DeleteFrom("alpha").Where(ConditionRaw("a", "b")).Limit(1).OrderDir("id", true).ToSql()
		if err != nil {
			b.Fatalf("%+v", err)
		}
	}
}

func TestDeleteAllToSql(t *testing.T) {
	s := createFakeSession()

	sql, _, err := s.DeleteFrom("a").ToSql()
	assert.NoError(t, err)
	assert.Equal(t, sql, "DELETE FROM `a`")

	sql, _, err = s.DeleteFrom("a", "b").ToSql()
	assert.NoError(t, err)
	assert.Equal(t, sql, "DELETE FROM `a` AS `b`")
}

func TestDeleteSingleToSql(t *testing.T) {
	s := createFakeSession()

	del := s.DeleteFrom("a").Where(ConditionRaw("id = ?", 1))
	sql, args, err := del.ToSql()
	assert.NoError(t, err)
	assert.Equal(t, sql, "DELETE FROM `a` WHERE (id = ?)")
	assert.Equal(t, args, []interface{}{1})

	// once where was a sync.Pool for the whereFragments with which it was
	// not possible to run ToSQL() twice.
	sql, args, err = del.ToSql()
	assert.NoError(t, err)
	assert.Equal(t, sql, "DELETE FROM `a` WHERE (id = ?)")
	assert.Equal(t, args, []interface{}{1})

}

func TestDeleteTenStaringFromTwentyToSql(t *testing.T) {
	s := createFakeSession()

	sql, _, err := s.DeleteFrom("a").Limit(10).Offset(20).OrderBy("id").ToSql()
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

func TestDeleteBuilder_Prepare(t *testing.T) {
	t.Skip("TODO")
}
