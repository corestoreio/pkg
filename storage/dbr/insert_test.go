package dbr

import (
	"context"
	"database/sql"
	"fmt"
	"testing"

	"github.com/corestoreio/errors"
	"github.com/corestoreio/log"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var _ ArgumentGenerater = (*someRecord)(nil)

type someRecord struct {
	SomethingID int
	UserID      int64
	Other       bool
}

func (sr someRecord) GenerateArguments(statementType byte, columns, condition []string) (Arguments, error) {
	args := make(Arguments, 0, 3) // 3 == number of fields in the struct
	for _, c := range columns {
		switch c {
		case "something_id":
			args = append(args, ArgInt(sr.SomethingID))
		case "user_id":
			args = append(args, ArgInt64(sr.UserID))
		case "other":
			args = append(args, ArgBool(sr.Other))
		default:
			return nil, errors.NewNotFoundf("[dbr_test] Column %q not found", c)
		}
	}
	return args, nil
}

var benchmarkInsertValuesSQLArgs Arguments

func BenchmarkInsertValuesSQL(b *testing.B) {
	s := createFakeSession()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, args, err := s.InsertInto("alpha").Columns("something_id", "user_id", "other").Values(
			ArgInt(1), ArgInt(2), ArgBool(true),
		).ToSQL()
		if err != nil {
			b.Fatal(err)
		}
		benchmarkInsertValuesSQLArgs = args
	}
}

func BenchmarkInsertRecordsSQL(b *testing.B) {
	s := createFakeSession()
	obj := someRecord{SomethingID: 1, UserID: 99, Other: false}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, args, err := s.InsertInto("alpha").
			Columns("something_id", "user_id", "other").
			Record(obj).
			ToSQL()
		if err != nil {
			b.Fatal(err)
		}
		benchmarkInsertValuesSQLArgs = args
		// ifaces = args.Interfaces()
	}
}

func TestInsertSingleToSQL(t *testing.T) {
	s := createFakeSession()

	sStr, args, err := s.InsertInto("a").Columns("b", "c").Values(ArgInt(1), ArgInt(2)).ToSQL()
	assert.NoError(t, err)
	assert.Equal(t, "INSERT INTO `a` (`b`,`c`) VALUES (?,?)", sStr)
	assert.Equal(t, []interface{}{int64(1), int64(2)}, args.Interfaces())
}

func TestInsertMultipleToSQL(t *testing.T) {
	s := createFakeSession()

	sStr, args, err := s.InsertInto("a").Columns("b", "c").
		Values(
			ArgInt(1), ArgInt(2),
			ArgInt(3), ArgInt(4),
		).
		Values(
			ArgInt(5), ArgInt(6),
		).ToSQL()
	assert.NoError(t, err)
	assert.Equal(t, "INSERT INTO `a` (`b`,`c`) VALUES (?,?),(?,?),(?,?)", sStr)
	assert.Equal(t, []interface{}{int64(1), int64(2), int64(3), int64(4), int64(5), int64(6)}, args.Interfaces())
}

func TestInsertRecordsToSQL(t *testing.T) {
	s := createFakeSession()

	objs := []someRecord{{1, 88, false}, {2, 99, true}, {3, 101, true}}
	sql, args, err := s.InsertInto("a").Columns("something_id", "user_id", "other").Record(objs[0]).Record(objs[1], objs[2]).ToSQL()
	require.NoError(t, err)
	assert.Equal(t, "INSERT INTO `a` (`something_id`,`user_id`,`other`) VALUES (?,?,?),(?,?,?),(?,?,?)", sql)
	// without fmt.Sprint we have an error despite objects are equal ...
	assert.Equal(t, fmt.Sprint([]interface{}{1, 88, false, 2, 99, true, 3, 101, true}), fmt.Sprint(args.Interfaces()))
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
	res, err := s.InsertInto("dbr_people").Columns("name", "key").Values(ArgString("Barack"), ArgString("44")).Exec(context.TODO())
	assert.NoError(t, err)

	rowsAff, err := res.RowsAffected()
	assert.NoError(t, err)
	assert.Equal(t, rowsAff, int64(1))
}

func TestInsertReal(t *testing.T) {
	// Insert by specifying values
	s := createRealSessionWithFixtures()
	res, err := s.InsertInto("dbr_people").Columns("name", "email").Values(ArgString("Barack"), ArgString("obama@whitehouse.gov")).Exec(context.TODO())
	validateInsertingBarack(t, s, res, err)

	// Insert by specifying a record (ptr to struct)
	s = createRealSessionWithFixtures()
	person := dbrPerson{Name: "Barack"}
	person.Email.Valid = true
	person.Email.String = "obama@whitehouse.gov"
	ib := s.InsertInto("dbr_people").Columns("name", "email").Record(&person)
	res, err = ib.Exec(context.TODO())
	if err != nil {
		t.Errorf("%s: %s", err, ib.String())
	}
	validateInsertingBarack(t, s, res, err)
}

func validateInsertingBarack(t *testing.T, c *Connection, res sql.Result, err error) {
	assert.NoError(t, err)
	if res == nil {
		t.Fatal("result at nit but should not")
	}
	id, err := res.LastInsertId()
	assert.NoError(t, err)
	rowsAff, err := res.RowsAffected()
	assert.NoError(t, err)

	assert.True(t, id > 0)
	assert.Equal(t, rowsAff, int64(1))

	var person dbrPerson
	err = c.Select("*").From("dbr_people").Where(Condition("id = ?", ArgInt64(id))).LoadStruct(context.TODO(), &person)
	assert.NoError(t, err)

	assert.Equal(t, id, person.ID)
	assert.Equal(t, "Barack", person.Name)
	assert.Equal(t, true, person.Email.Valid)
	assert.Equal(t, "obama@whitehouse.gov", person.Email.String)
}

// TODO: do a real test inserting multiple records

func TestInsert_Prepare(t *testing.T) {

	t.Run("ToSQL Error", func(t *testing.T) {
		in := &Insert{}
		in.Columns("a", "b")
		stmt, err := in.Prepare(context.TODO())
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
		in.Columns("a", "b").Values(ArgInt(1), ArgBool(true))

		stmt, err := in.Prepare(context.TODO())
		assert.Nil(t, stmt)
		assert.True(t, errors.IsAlreadyClosed(err), "%+v", err)
	})
}

func TestInsert_Events(t *testing.T) {
	t.Parallel()

	t.Run("Stop Propagation", func(t *testing.T) {
		d := NewInsert("tableA")
		d.Columns("a", "b").Values(ArgInt(1), ArgBool(true))

		d.Log = log.BlackHole{EnableInfo: true, EnableDebug: true}
		d.Listeners.Add(
			Listen{
				Name:      "listener1",
				EventType: OnBeforeToSQL,
				InsertFunc: func(b *Insert) {
					b.Pair("col1", ArgString("X1"))
				},
			},
			Listen{
				Name:      "listener2",
				EventType: OnBeforeToSQL,
				InsertFunc: func(b *Insert) {
					b.Pair("col2", ArgString("X2"))
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
		assert.Exactly(t, "INSERT INTO `tableA` (`a`,`b`,`col1`,`col2`) VALUES (?,?,?,?)", sql)

		sql, _, err = d.ToSQL()
		assert.NoError(t, err, "%+v", err)
		assert.Exactly(t, "INSERT INTO `tableA` (`a`,`b`,`col1`,`col2`) VALUES (?,?,?,?)", sql)
	})

	t.Run("Missing EventType", func(t *testing.T) {
		ins := NewInsert("tableA")
		ins.Columns("a", "b").Values(ArgInt(1), ArgBool(true))

		ins.Listeners.Add(
			Listen{
				Name: "colC",
				InsertFunc: func(i *Insert) {
					i.Pair("colC", ArgString("X1"))
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

		ins.Columns("a", "b").Values(ArgInt(1), ArgBool(true))

		ins.Listeners.Add(
			Listen{
				EventType: OnBeforeToSQL,
				Name:      "colA",
				Once:      true,
				InsertFunc: func(i *Insert) {
					i.Pair("colA", ArgFloat64(3.14159))
				},
			},
			Listen{
				EventType: OnBeforeToSQL,
				Name:      "colB",
				Once:      true,
				InsertFunc: func(i *Insert) {
					i.Pair("colB", ArgFloat64(2.7182))
				},
			},
		)

		ins.Listeners.Add(
			Listen{
				EventType: OnBeforeToSQL,
				Name:      "colC",
				InsertFunc: func(i *Insert) {
					i.Pair("colC", ArgString("X1"))
				},
			},
		)

		sql, args, err := ins.ToSQL()
		assert.NoError(t, err)
		assert.Exactly(t, []interface{}{int64(1), true, 3.14159, 2.7182, "X1"}, args.Interfaces())
		assert.Exactly(t, "INSERT INTO `tableA` (`a`,`b`,`colA`,`colB`,`colC`) VALUES (?,?,?,?,?)", sql)

		sql, args, err = ins.ToSQL()
		assert.NoError(t, err)
		assert.Exactly(t, []interface{}{int64(1), true, 3.14159, 2.7182, "X1"}, args.Interfaces())
		assert.Exactly(t, "INSERT INTO `tableA` (`a`,`b`,`colA`,`colB`,`colC`) VALUES (?,?,?,?,?)", sql)

		assert.Exactly(t, `colA; colB; colC`, ins.Listeners.String())
	})
}

func TestInsert_FromSelect(t *testing.T) {
	ins := NewInsert("tableA")
	// columns and args just to check that they get ignored
	ins.Columns("a", "b").Values(ArgInt(1), ArgBool(true))

	argEq := Eq{"a": ArgInt64(1, 2, 3).Operator(OperatorIn)}
	args := Arguments{ArgInt64(1), ArgString("wat")}

	iSQL, args, err := ins.FromSelect(NewSelect("some_table").
		AddColumns("something_id", "user_id", "other").
		Where(Condition("d = ? OR e = ?", args...)).
		Where(argEq).
		OrderDir("id", false).
		Paginate(1, 20))
	if err != nil {
		t.Fatalf("%+v", err)
	}
	assert.Exactly(t, "INSERT INTO `tableA` SELECT something_id, user_id, other FROM `some_table` WHERE (d = ? OR e = ?) AND (`a` IN ?) ORDER BY id DESC LIMIT 20 OFFSET 0", iSQL)
	assert.Exactly(t, []interface{}{int64(1), "wat", int64(1), int64(2), int64(3)}, args.Interfaces())

}
