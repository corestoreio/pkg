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
	"database/sql"
	"fmt"
	"testing"

	"github.com/corestoreio/errors"
	"github.com/corestoreio/log"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

//var _ UpdateArgProducer = (*someRecord)(nil)
var _ InsertArgProducer = (*someRecord)(nil)

type someRecord struct {
	SomethingID int
	UserID      int64
	Other       bool
}

func (sr someRecord) ProduceInsertArgs(args Arguments, columns []string) (Arguments, error) {
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

func TestInsert_Interpolate(t *testing.T) {
	sStr, args, err := NewInsert("a").AddColumns("b", "c").
		AddValues(
			argInt(1), argInt(2),
			argInt(3), argInt(4),
		).
		AddValues(
			argInt(5), argInt(6),
		).
		AddOnDuplicateKey("b", nil).
		AddOnDuplicateKey("c", nil).
		Interpolate().
		ToSQL()

	assert.NoError(t, err)
	assert.Equal(t, "INSERT INTO `a` (`b`,`c`) VALUES (1,2),(3,4),(5,6) ON DUPLICATE KEY UPDATE `b`=VALUES(`b`), `c`=VALUES(`c`)", sStr)
	assert.Nil(t, args)
}

func TestInsertRecordsToSQL(t *testing.T) {
	s := createFakeSession()

	objs := []someRecord{{1, 88, false}, {2, 99, true}, {3, 101, true}}
	sqlStr, args, err := s.InsertInto("a").
		AddColumns("something_id", "user_id", "other").
		AddRecords(objs[0]).AddRecords(objs[1], objs[2]).
		AddOnDuplicateKey("something_id", argInt64(99)).
		AddOnDuplicateKey("user_id", nil).
		ToSQL()
	require.NoError(t, err)
	assert.Equal(t, "INSERT INTO `a` (`something_id`,`user_id`,`other`) VALUES (?,?,?),(?,?,?),(?,?,?) ON DUPLICATE KEY UPDATE `something_id`=?, `user_id`=VALUES(`user_id`)", sqlStr)
	// without fmt.Sprint we have an error despite objects are equal ...
	assert.Equal(t, fmt.Sprint([]interface{}{1, 88, false, 2, 99, true, 3, 101, true, int64(99)}), fmt.Sprint(args.Interfaces()))
}

func TestInsertRecordsToSQLNotFoundMapping(t *testing.T) {
	s := createFakeSession()

	objs := []someRecord{{1, 88, false}, {2, 99, true}}
	sqlStr, args, err := s.InsertInto("a").AddColumns("something_it", "user_id", "other").AddRecords(objs[0]).AddRecords(objs[1]).ToSQL()
	assert.True(t, errors.IsNotFound(err), "%+v", err)
	assert.Nil(t, args)
	assert.Empty(t, sqlStr)
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
	_, err = c.Select("*").From("dbr_people").Where(Column("id", ArgInt64(id))).Load(context.TODO(), &person)
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
		_, err = s.Select("*").From("dbr_people").Where(Column("id", ArgInt64(inID))).Load(context.TODO(), &p)
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
		_, err = s.Select("*").From("dbr_people").Where(Column("id", ArgInt64(inID))).Load(context.TODO(), &p)
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

	iSQL, args, err := ins.FromSelect(NewSelect("something_id", "user_id", "other").
		From("some_table").
		Where(
			ParenthesisOpen(),
			Column("d", argInt64(1)),
			Column("e", ArgString("wat")).Or(),
			ParenthesisClose(),
		).
		Where(argEq).
		OrderByDesc("id").
		Paginate(1, 20)).ToSQL()
	if err != nil {
		t.Fatalf("%+v", err)
	}
	assert.Exactly(t, "INSERT INTO `tableA` SELECT `something_id`, `user_id`, `other` FROM `some_table` WHERE ((`d` = ?) OR (`e` = ?)) AND (`a` IN ?) ORDER BY `id` DESC LIMIT 20 OFFSET 0", iSQL)
	assert.Exactly(t, []interface{}{int64(1), "wat", int64(1), int64(2), int64(3)}, args.Interfaces())
}

func TestInsert_Replace_Ignore(t *testing.T) {

	// this generated statement does not comply the SQL standard
	sStr, args, err := NewInsert("a").
		Replace().Ignore().
		AddColumns("b", "c").
		AddValues(argInt(1), argInt(2)).
		AddValues(argInt64(3), argInt64(4)).
		ToSQL()
	assert.NoError(t, err)
	assert.Equal(t, "REPLACE IGNORE INTO `a` (`b`,`c`) VALUES (?,?),(?,?)", sStr)
	assert.Equal(t, []interface{}{int64(1), int64(2), int64(3), int64(4)}, args.Interfaces())
}

func TestInsert_WithoutColumns(t *testing.T) {

	t.Run("each column in its own Arg", func(t *testing.T) {
		ins := NewInsert("catalog_product_link").
			AddValues(ArgInt64(2046), ArgInt64(33), ArgInt64(3)).
			AddValues(ArgInt64(2046), ArgInt64(34), ArgInt64(3)).
			AddValues(ArgInt64(2046), ArgInt64(35), ArgInt64(3))

		sStr, args, err := ins.ToSQL()
		assert.NoError(t, err)
		assert.Exactly(t, []interface{}{int64(2046), int64(33), int64(3), int64(2046), int64(34), int64(3), int64(2046), int64(35), int64(3)}, args.Interfaces())
		assert.Exactly(t, "INSERT INTO `catalog_product_link` VALUES (?,?,?),(?,?,?),(?,?,?)", sStr)
	})
}

func TestInsert_Pair(t *testing.T) {
	t.Run("one row", func(t *testing.T) {
		ins := NewInsert("catalog_product_link").
			Pair("product_id", ArgInt64(2046)).
			Pair("linked_product_id", ArgInt64(33)).
			Pair("link_type_id", ArgInt64(3))
		sStr, args, err := ins.ToSQL()
		assert.NoError(t, err)
		assert.Exactly(t, []interface{}{int64(2046), int64(33), int64(3)}, args.Interfaces())
		assert.Exactly(t, "INSERT INTO `catalog_product_link` (`product_id`,`linked_product_id`,`link_type_id`) VALUES (?,?,?)", sStr)
	})
	t.Run("multiple rows triggers error", func(t *testing.T) {
		ins := NewInsert("catalog_product_link").
			Pair("product_id", ArgInt64(2046)).
			Pair("linked_product_id", ArgInt64(33)).
			Pair("link_type_id", ArgInt64(3)).
			// next row
			Pair("product_id", ArgInt64(2046)).
			Pair("linked_product_id", ArgInt64(34)).
			Pair("link_type_id", ArgInt64(3))

		sStr, args, err := ins.ToSQL()
		assert.Empty(t, sStr)
		assert.Nil(t, args)
		assert.True(t, errors.IsAlreadyExists(err), "%+v", err)
	})
}

func TestInsert_UseBuildCache(t *testing.T) {
	t.Parallel()

	ins := NewInsert("a").AddColumns("b", "c").
		AddValues(
			argInt(1), argInt(2),
			argInt(3), argInt(4),
		).
		AddValues(
			argInt(5), argInt(6),
		).
		AddOnDuplicateKey("b", nil).
		AddOnDuplicateKey("c", nil)

	ins.UseBuildCache = true

	const cachedSQLPlaceHolder = "INSERT INTO `a` (`b`,`c`) VALUES (?,?),(?,?),(?,?) ON DUPLICATE KEY UPDATE `b`=VALUES(`b`), `c`=VALUES(`c`)"
	t.Run("without interpolate", func(t *testing.T) {
		for i := 0; i < 3; i++ {
			sql, args, err := ins.ToSQL()
			require.NoError(t, err, "%+v", err)
			require.Equal(t, cachedSQLPlaceHolder, sql)
			assert.Equal(t, []interface{}{int64(1), int64(2), int64(3), int64(4), int64(5), int64(6)}, args.Interfaces())
			assert.Equal(t, cachedSQLPlaceHolder, string(ins.buildCache))
		}
	})

	t.Run("with interpolate", func(t *testing.T) {
		ins.Interpolate()
		ins.buildCache = nil
		ins.RawArguments = nil

		const cachedSQLInterpolated = "INSERT INTO `a` (`b`,`c`) VALUES (1,2),(3,4),(5,6) ON DUPLICATE KEY UPDATE `b`=VALUES(`b`), `c`=VALUES(`c`)"
		for i := 0; i < 3; i++ {
			sql, args, err := ins.ToSQL()
			assert.Equal(t, cachedSQLPlaceHolder, string(ins.buildCache))
			require.NoError(t, err, "%+v", err)
			require.Equal(t, cachedSQLInterpolated, sql)
			assert.Nil(t, args)
		}
	})
}
