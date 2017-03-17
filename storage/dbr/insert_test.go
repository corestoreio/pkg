package dbr

import (
	"database/sql"
	"fmt"
	"testing"

	"github.com/corestoreio/errors"
	"github.com/corestoreio/log"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type someRecord struct {
	SomethingID int   `db:"something_id"`
	UserID      int64 `db:"user_id"`
	Other       bool
}

func BenchmarkInsertValuesSQL(b *testing.B) {
	s := createFakeSession()

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		s.InsertInto("alpha").Columns("something_id", "user_id", "other").Values(1, 2, true).ToSQL()
	}
}

func BenchmarkInsertRecordsSQL(b *testing.B) {
	s := createFakeSession()
	obj := someRecord{1, 99, false}

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		s.InsertInto("alpha").Columns("something_id", "user_id", "other").Record(obj).ToSQL()
	}
}

func TestInsertSingleToSQL(t *testing.T) {
	s := createFakeSession()

	sql, args, err := s.InsertInto("a").Columns("b", "c").Values(1, 2).ToSQL()
	assert.NoError(t, err)
	assert.Equal(t, sql, "INSERT INTO a (`b`,`c`) VALUES (?,?)")
	assert.Equal(t, args, []interface{}{1, 2})
}

func TestInsertMultipleToSQL(t *testing.T) {
	s := createFakeSession()

	sql, args, err := s.InsertInto("a").Columns("b", "c").Values(1, 2).Values(3, 4).ToSQL()
	assert.NoError(t, err)
	assert.Equal(t, sql, "INSERT INTO a (`b`,`c`) VALUES (?,?),(?,?)")
	assert.Equal(t, args, []interface{}{1, 2, 3, 4})
}

func TestInsertRecordsToSQL(t *testing.T) {
	s := createFakeSession()

	objs := []someRecord{{1, 88, false}, {2, 99, true}}
	sql, args, err := s.InsertInto("a").Columns("something_id", "user_id", "other").Record(objs[0]).Record(objs[1]).ToSQL()
	require.NoError(t, err)
	assert.Equal(t, "INSERT INTO a (`something_id`,`user_id`,`other`) VALUES (?,?,?),(?,?,?)", sql)
	// without fmt.Sprint we have an error despite objects are equal ...
	assert.Equal(t, fmt.Sprint([]interface{}{1, 88, false, 2, 99, true}), fmt.Sprint(args))
}

func TestInsertRecordsToSQLNotFoundMapping(t *testing.T) {
	s := createFakeSession()

	objs := []someRecord{{1, 88, false}, {2, 99, true}}
	sql, args, err := s.InsertInto("a").Columns("something_it", "user_id", "other").Record(objs[0]).Record(objs[1]).ToSQL()
	assert.True(t, errors.IsNotFound(err), "%+v", err)
	assert.Nil(t, args)
	assert.Empty(t, sql)
}

func TestInsertKeywordColumnName(t *testing.T) {
	// Insert a column whose name is reserved
	s := createRealSessionWithFixtures()
	res, err := s.InsertInto("dbr_people").Columns("name", "key").Values("Barack", "44").Exec()
	assert.NoError(t, err)

	rowsAff, err := res.RowsAffected()
	assert.NoError(t, err)
	assert.Equal(t, rowsAff, int64(1))
}

func TestInsertReal(t *testing.T) {
	// Insert by specifying values
	s := createRealSessionWithFixtures()
	res, err := s.InsertInto("dbr_people").Columns("name", "email").Values("Barack", "obama@whitehouse.gov").Exec()
	validateInsertingBarack(t, s, res, err)

	// Insert by specifying a record (ptr to struct)
	s = createRealSessionWithFixtures()
	person := dbrPerson{Name: "Barack"}
	person.Email.Valid = true
	person.Email.String = "obama@whitehouse.gov"
	res, err = s.InsertInto("dbr_people").Columns("name", "email").Record(&person).Exec()
	validateInsertingBarack(t, s, res, err)

	// Insert by specifying a record (struct)
	s = createRealSessionWithFixtures()
	res, err = s.InsertInto("dbr_people").Columns("name", "email").Record(person).Exec()
	validateInsertingBarack(t, s, res, err)
}

func validateInsertingBarack(t *testing.T, s *Session, res sql.Result, err error) {
	assert.NoError(t, err)
	id, err := res.LastInsertId()
	assert.NoError(t, err)
	rowsAff, err := res.RowsAffected()
	assert.NoError(t, err)

	assert.True(t, id > 0)
	assert.Equal(t, rowsAff, int64(1))

	var person dbrPerson
	err = s.Select("*").From("dbr_people").Where(ConditionRaw("id = ?", id)).LoadStruct(&person)
	assert.NoError(t, err)

	assert.Equal(t, person.ID, id)
	assert.Equal(t, person.Name, "Barack")
	assert.Equal(t, person.Email.Valid, true)
	assert.Equal(t, person.Email.String, "obama@whitehouse.gov")
}

// TODO: do a real test inserting multiple records

func TestInsert_Prepare(t *testing.T) {

	t.Run("ToSQL Error", func(t *testing.T) {
		in := &Insert{}
		in.Columns("a", "b")
		stmt, err := in.Prepare()
		assert.Nil(t, stmt)
		assert.True(t, errors.IsEmpty(err))
	})

	t.Run("Prepare Error", func(t *testing.T) {
		in := &Insert{
			Into: "table",
		}
		in.DB.Preparer = dbMock{
			error: errors.NewAlreadyClosedf("Who closed myself?"),
		}
		in.Columns("a", "b").Values(1, true)

		stmt, err := in.Prepare()
		assert.Nil(t, stmt)
		assert.True(t, errors.IsAlreadyClosed(err), "%+v", err)
	})
}

func TestInsert_Events(t *testing.T) {
	t.Parallel()

	t.Run("Stop Propagation", func(t *testing.T) {
		d := NewInsert("tableA")
		d.Columns("a", "b").Values(1, true)

		d.Log = log.BlackHole{EnableInfo: true, EnableDebug: true}
		d.Listeners.Add(
			Listen{
				Name:      "listener1",
				EventType: OnBeforeToSQL,
				InsertFunc: func(b *Insert) {
					b.Pair("col1", "X1")
				},
			},
			Listen{
				Name:      "listener2",
				EventType: OnBeforeToSQL,
				InsertFunc: func(b *Insert) {
					b.Pair("col2", "X2")
					b.PropagationStopped = true
				},
			},
			Listen{
				Name:      "listener3",
				EventType: OnBeforeToSQL,
				InsertFunc: func(b *Insert) {
					panic("Should not get called")
				},
			},
		)
		sql, _, err := d.ToSQL()
		assert.NoError(t, err, "%+v", err)
		assert.Exactly(t, "INSERT INTO tableA (`a`,`b`,`col1`,`col2`) VALUES (?,?,?,?)", sql)

		sql, _, err = d.ToSQL()
		assert.NoError(t, err, "%+v", err)
		assert.Exactly(t, "INSERT INTO tableA (`a`,`b`,`col1`,`col2`,`col1`,`col2`) VALUES (?,?,?,?,?,?)", sql)
	})

	t.Run("Missing EventType", func(t *testing.T) {
		ins := NewInsert("tableA")
		ins.Columns("a", "b").Values(1, true)

		ins.Listeners.Add(
			Listen{
				Name: "colC",
				InsertFunc: func(i *Insert) {
					i.Pair("colC", "X1")
				},
			},
		)
		sql, args, err := ins.ToSQL()
		assert.Empty(t, sql)
		assert.Nil(t, args)
		assert.True(t, errors.IsEmpty(err), "%+v", err)
	})

	t.Run("Should Dispatch", func(t *testing.T) {
		ins := NewInsert("tableA")

		ins.Columns("a", "b").Values(1, true)

		ins.Listeners.Add(
			Listen{
				EventType: OnBeforeToSQL,
				Name:      "colA",
				Once:      true,
				InsertFunc: func(i *Insert) {
					i.Pair("colA", 3.14159)
				},
			},
			Listen{
				EventType: OnBeforeToSQL,
				Name:      "colB",
				Once:      true,
				InsertFunc: func(i *Insert) {
					i.Pair("colB", 2.7182)
				},
			},
		)

		ins.Listeners.Add(
			Listen{
				EventType: OnBeforeToSQL,
				Name:      "colC",
				InsertFunc: func(i *Insert) {
					i.Pair("colC", "X1")
				},
			},
		)

		sql, args, err := ins.ToSQL()
		assert.NoError(t, err)
		assert.Exactly(t, []interface{}{1, true, 3.14159, 2.7182, "X1"}, args)
		assert.Exactly(t, "INSERT INTO tableA (`a`,`b`,`colA`,`colB`,`colC`) VALUES (?,?,?,?,?)", sql)

		sql, args, err = ins.ToSQL()
		assert.NoError(t, err)
		assert.Exactly(t, []interface{}{1, true, 3.14159, 2.7182, "X1", "X1"}, args)
		assert.Exactly(t, "INSERT INTO tableA (`a`,`b`,`colA`,`colB`,`colC`,`colC`) VALUES (?,?,?,?,?,?)", sql)

		assert.Exactly(t, `colA; colB; colC`, ins.Listeners.String())
	})

}
