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
		_, args, err := s.InsertInto("alpha").AddColumns("something_id", "user_id", "other").AddValues(
			argInt(1), argInt(2), ArgBool(true),
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
			AddColumns("something_id", "user_id", "other").
			AddRecords(obj).
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

	sStr, args, err := s.InsertInto("a").AddColumns("b", "c").AddValues(argInt(1), argInt(2)).ToSQL()
	assert.NoError(t, err)
	assert.Equal(t, "INSERT INTO `a` (`b`,`c`) VALUES (?,?)", sStr)
	assert.Equal(t, []interface{}{int64(1), int64(2)}, args.Interfaces())
}

func TestInsertMultipleToSQL(t *testing.T) {
	s := createFakeSession()

	sStr, args, err := s.InsertInto("a").AddColumns("b", "c").
		AddValues(
			argInt(1), argInt(2),
			argInt(3), argInt(4),
		).
		AddValues(
			argInt(5), argInt(6),
		).
		AddOnDuplicateKey("b", nil).
		AddOnDuplicateKey("c", nil).
		ToSQL()
	assert.NoError(t, err)
	assert.Equal(t, "INSERT INTO `a` (`b`,`c`) VALUES (?,?),(?,?),(?,?) ON DUPLICATE KEY UPDATE `b`=VALUES(`b`), `c`=VALUES(`c`)", sStr)
	assert.Equal(t, []interface{}{int64(1), int64(2), int64(3), int64(4), int64(5), int64(6)}, args.Interfaces())
}

func TestInsertRecordsToSQL(t *testing.T) {
	s := createFakeSession()

	objs := []someRecord{{1, 88, false}, {2, 99, true}, {3, 101, true}}
	sql, args, err := s.InsertInto("a").
		AddColumns("something_id", "user_id", "other").
		AddRecords(objs[0]).AddRecords(objs[1], objs[2]).
		AddOnDuplicateKey("something_id", argInt64(99)).
		AddOnDuplicateKey("user_id", nil).
		ToSQL()
	require.NoError(t, err)
	assert.Equal(t, "INSERT INTO `a` (`something_id`,`user_id`,`other`) VALUES (?,?,?),(?,?,?),(?,?,?) ON DUPLICATE KEY UPDATE `something_id`=?, `user_id`=VALUES(`user_id`)", sql)
	// without fmt.Sprint we have an error despite objects are equal ...
	assert.Equal(t, fmt.Sprint([]interface{}{1, 88, false, 2, 99, true, 3, 101, true, int64(99)}), fmt.Sprint(args.Interfaces()))
}

func TestInsertRecordsToSQLNotFoundMapping(t *testing.T) {
	s := createFakeSession()

	objs := []someRecord{{1, 88, false}, {2, 99, true}}
	sql, args, err := s.InsertInto("a").AddColumns("something_it", "user_id", "other").AddRecords(objs[0]).AddRecords(objs[1]).ToSQL()
	assert.True(t, errors.IsNotFound(err), "%+v", err)
	assert.Nil(t, args)
	assert.Empty(t, sql)
}

func TestInsertKeywordColumnName(t *testing.T) {
	// Insert a column whose name is reserved
	s := createRealSessionWithFixtures()
	res, err := s.InsertInto("dbr_people").AddColumns("name", "key").AddValues(ArgString("Barack"), ArgString("44")).Exec(context.TODO())
	assert.NoError(t, err)

	rowsAff, err := res.RowsAffected()
	assert.NoError(t, err)
	assert.Equal(t, rowsAff, int64(1))
}

func TestInsertReal(t *testing.T) {
	// Insert by specifying values
	s := createRealSessionWithFixtures()
	res, err := s.InsertInto("dbr_people").AddColumns("name", "email").AddValues(ArgString("Barack"), ArgString("obama@whitehouse.gov")).Exec(context.TODO())
	validateInsertingBarack(t, s, res, err)

	// Insert by specifying a record (ptr to struct)
	s = createRealSessionWithFixtures()
	person := dbrPerson{Name: "Barack"}
	person.Email.Valid = true
	person.Email.String = "obama@whitehouse.gov"
	ib := s.InsertInto("dbr_people").AddColumns("name", "email").AddRecords(&person)
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

func TestInsertReal_OnDuplicateKey(t *testing.T) {

	s := createRealSessionWithFixtures()
	res, err := s.InsertInto("dbr_people").
		AddColumns("id", "name", "email").
		AddValues(ArgInt64(678), ArgString("Pike"), ArgString("pikes@peak.co")).Exec(context.TODO())
	if err != nil {
		t.Fatalf("%+v", err)
	}
	inID, err := res.LastInsertId()
	if err != nil {
		t.Fatalf("%+v", err)
	}
	{
		var p dbrPerson
		err = s.Select("*").From("dbr_people").Where(Condition("id = ?", ArgInt64(inID))).LoadStruct(context.TODO(), &p)
		assert.NoError(t, err)
		assert.Equal(t, "Pike", p.Name)
		assert.Equal(t, "pikes@peak.co", p.Email.String)
	}
	res, err = s.InsertInto("dbr_people").
		AddColumns("id", "name", "email").
		AddValues(ArgInt64(inID), ArgString(""), ArgString("pikes@peak.com")).
		AddOnDuplicateKey("name", ArgString("Pik3")).
		AddOnDuplicateKey("email", nil).
		Exec(context.TODO())
	if err != nil {
		t.Fatalf("%+v", err)
	}
	inID2, err := res.LastInsertId()
	if err != nil {
		t.Fatalf("%+v", err)
	}
	assert.Exactly(t, inID, inID2)

	{
		var p dbrPerson
		err = s.Select("*").From("dbr_people").Where(Condition("id = ?", ArgInt64(inID))).LoadStruct(context.TODO(), &p)
		assert.NoError(t, err)
		assert.Equal(t, "Pik3", p.Name)
		assert.Equal(t, "pikes@peak.com", p.Email.String)
	}
}

// TODO: do a real test inserting multiple records

func TestInsert_Prepare(t *testing.T) {

	t.Run("ToSQL Error", func(t *testing.T) {
		in := &Insert{}
		in.AddColumns("a", "b")
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
		in.AddColumns("a", "b").AddValues(argInt(1), ArgBool(true))

		stmt, err := in.Prepare(context.TODO())
		assert.Nil(t, stmt)
		assert.True(t, errors.IsAlreadyClosed(err), "%+v", err)
	})
}

func TestInsert_Events(t *testing.T) {
	t.Parallel()

	t.Run("Stop Propagation", func(t *testing.T) {
		d := NewInsert("tableA")
		d.AddColumns("a", "b").AddValues(argInt(1), ArgBool(true))

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
		ins.AddColumns("a", "b").AddValues(argInt(1), ArgBool(true))

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

		ins.AddColumns("a", "b").AddValues(argInt(1), ArgBool(true))

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
	ins.AddColumns("a", "b").AddValues(argInt(1), ArgBool(true))

	argEq := Eq{"a": ArgInt64(1, 2, 3).Operator(In)}
	args := Arguments{argInt64(1), ArgString("wat")}

	iSQL, args, err := ins.FromSelect(NewSelect("something_id", "user_id", "other").
		From("some_table").
		Where(Condition("d = ? OR e = ?", args...)).
		Where(argEq).
		OrderByDesc("id").
		Paginate(1, 20))
	if err != nil {
		t.Fatalf("%+v", err)
	}
	assert.Exactly(t, "INSERT INTO `tableA` SELECT something_id, user_id, other FROM `some_table` WHERE (d = ? OR e = ?) AND (`a` IN ?) ORDER BY id DESC LIMIT 20 OFFSET 0", iSQL)
	assert.Exactly(t, []interface{}{int64(1), "wat", int64(1), int64(2), int64(3)}, args.Interfaces())

}
