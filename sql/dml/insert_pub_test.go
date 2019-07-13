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

package dml_test

import (
	"context"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/corestoreio/errors"
	"github.com/corestoreio/log"
	"github.com/corestoreio/pkg/sql/dml"
	"github.com/corestoreio/pkg/sql/dmltest"
	"github.com/corestoreio/pkg/storage/null"
	"github.com/corestoreio/pkg/util/assert"
)

var _ dml.ColumnMapper = (*someRecord)(nil)

type someRecord struct {
	SomethingID int
	UserID      int64
	Other       bool
}

func (sr someRecord) MapColumns(cm *dml.ColumnMap) error {
	if cm.Mode() == dml.ColumnMapEntityReadAll {
		return cm.Int(&sr.SomethingID).Int64(&sr.UserID).Bool(&sr.Other).Err()
	}
	for cm.Next() {
		switch c := cm.Column(); c {
		case "something_id":
			cm.Int(&sr.SomethingID)
		case "user_id":
			cm.Int64(&sr.UserID)
		case "other":
			cm.Bool(&sr.Other)
		default:
			return errors.NotFound.Newf("[dml_test] Column %q not found", c)
		}
	}
	return cm.Err()
}

func TestInsert_Bind(t *testing.T) {
	t.Parallel()

	objs := []someRecord{{1, 88, false}, {2, 99, true}, {3, 101, true}}
	wantArgs := []interface{}{int64(1), int64(88), false, int64(2), int64(99), true, int64(3), int64(101), true}

	t.Run("valid with multiple records", func(t *testing.T) {
		compareToSQL(t,
			dml.NewInsert("a").
				AddColumns("something_id", "user_id", "other").
				AddOnDuplicateKey(
					dml.Column("something_id").Int64(99),
					dml.Column("user_id").Values(),
				).WithArgs().Record("", objs[0]).Record("", objs[1]).Record("", objs[2]),
			errors.NoKind,
			"INSERT INTO `a` (`something_id`,`user_id`,`other`) VALUES (?,?,?),(?,?,?),(?,?,?) ON DUPLICATE KEY UPDATE `something_id`=99, `user_id`=VALUES(`user_id`)",
			"INSERT INTO `a` (`something_id`,`user_id`,`other`) VALUES (1,88,0),(2,99,1),(3,101,1) ON DUPLICATE KEY UPDATE `something_id`=99, `user_id`=VALUES(`user_id`)",
			wantArgs...,
		)
	})
	t.Run("without columns, all columns requested, with AddOnDuplicateKey", func(t *testing.T) {
		compareToSQL(t,
			dml.NewInsert("a").
				SetRecordPlaceHolderCount(3).
				AddOnDuplicateKey(
					dml.Column("something_id").Int64(99),
					dml.Column("user_id").Values(),
				).WithArgs().Record("", objs[0]).Record("", objs[1]).Record("", objs[2]),
			errors.NoKind,
			"INSERT INTO `a` VALUES (?,?,?),(?,?,?),(?,?,?) ON DUPLICATE KEY UPDATE `something_id`=99, `user_id`=VALUES(`user_id`)",
			"INSERT INTO `a` VALUES (1,88,0),(2,99,1),(3,101,1) ON DUPLICATE KEY UPDATE `something_id`=99, `user_id`=VALUES(`user_id`)",
			wantArgs...,
		)
	})
	t.Run("without columns, all columns requested", func(t *testing.T) {
		customers := []*customerEntity{
			{EntityID: 11, Firstname: "Karl Gopher", StoreID: 0x7, LifetimeSales: null.MakeFloat64(47.11), VoucherCodes: exampleStringSlice{"1FE9983E", "28E76FBC"}},
			{EntityID: 12, Firstname: "Fung Go Roo", StoreID: 0x7, LifetimeSales: null.MakeFloat64(28.94), VoucherCodes: exampleStringSlice{"4FE7787E", "15E59FBB", "794EFDE8"}},
			{EntityID: 13, Firstname: "John Doe", StoreID: 0x6, LifetimeSales: null.MakeFloat64(138.54), VoucherCodes: exampleStringSlice{""}},
		}

		compareToSQL(t,
			dml.NewInsert("customer_entity").
				SetRecordPlaceHolderCount(5). // mandatory because no columns provided!
				WithArgs().Record("", customers[0]).Record("", customers[1]).Record("", customers[2]),
			errors.NoKind,
			"INSERT INTO `customer_entity` VALUES (?,?,?,?,?),(?,?,?,?,?),(?,?,?,?,?)",
			"INSERT INTO `customer_entity` VALUES (11,'Karl Gopher',7,47.11,'1FE9983E|28E76FBC'),(12,'Fung Go Roo',7,28.94,'4FE7787E|15E59FBB|794EFDE8'),(13,'John Doe',6,138.54,'')",
			int64(11), "Karl Gopher", int64(7), 47.11, "1FE9983E|28E76FBC", int64(12), "Fung Go Roo", int64(7), 28.94, "4FE7787E|15E59FBB|794EFDE8", int64(13), "John Doe", int64(6), 138.54, "",
		)
	})
	t.Run("column not found", func(t *testing.T) {
		objs := []someRecord{{1, 88, false}, {2, 99, true}}
		compareToSQL(t,
			dml.NewInsert("a").AddColumns("something_it", "user_id", "other").WithArgs().Record("", objs[0]).Record("", objs[1]),
			errors.NotFound,
			"",
			"",
		)
	})
}

func TestInsert_Prepare(t *testing.T) {
	t.Parallel()

	t.Run("BuildValues not set Error", func(t *testing.T) {
		in := &dml.Insert{}
		in.AddColumns("a", "b")
		stmt, err := in.Prepare(context.TODO())
		assert.Nil(t, stmt)
		assert.ErrorIsKind(t, errors.NotAcceptable, err)
	})

	t.Run("ToSQL Error", func(t *testing.T) {
		in := &dml.Insert{}
		in.AddColumns("a", "b")
		stmt, _, err := in.ToSQL()
		assert.Empty(t, stmt)
		assert.ErrorIsKind(t, errors.Empty, err)
	})

	t.Run("DB Error", func(t *testing.T) {
		in := &dml.Insert{
			Into: "table",
		}
		in.DB = dbMock{
			error: errors.AlreadyClosed.Newf("Who closed myself?"),
		}
		in.AddColumns("a", "b").BuildValues()

		stmt, err := in.Prepare(context.TODO())
		assert.Nil(t, stmt)
		assert.ErrorIsKind(t, errors.AlreadyClosed, err)
	})

	t.Run("ExecArgs One Row", func(t *testing.T) {
		dbc, dbMock := dmltest.MockDB(t)
		defer dmltest.MockClose(t, dbc, dbMock)

		prep := dbMock.ExpectPrepare(dmltest.SQLMockQuoteMeta("INSERT INTO `customer_entity` (`email`,`group_id`,`created_at`) VALUES (?,?,?)"))
		prep.ExpectExec().WithArgs("a@b.c", 33, now()).WillReturnResult(sqlmock.NewResult(4, 0))
		prep.ExpectExec().WithArgs("x@y.z", 44, now().Add(time.Minute)).WillReturnResult(sqlmock.NewResult(5, 0))

		stmt, err := dml.NewInsert("customer_entity").
			AddColumns("email", "group_id", "created_at").BuildValues().
			WithDB(dbc.DB).
			Prepare(context.TODO())
		assert.NoError(t, err, "failed creating a prepared statement")
		defer func() {
			assert.NoError(t, stmt.Close(), "Close on a prepared statement")
		}()

		tests := []struct {
			email      string
			groupID    int
			created_at time.Time
			insertID   int64
		}{
			{"a@b.c", 33, now(), 4},
			{"x@y.z", 44, now().Add(time.Minute), 5},
		}

		stmtA := stmt.WithArgs()
		for i, test := range tests {
			res, err := stmtA.String(test.email).Int(test.groupID).Time(test.created_at).ExecContext(context.TODO())
			if err != nil {
				t.Fatalf("Index %d => %+v", i, err)
			}
			lid, err := res.LastInsertId()
			if err != nil {
				t.Fatalf("Result index %d with error: %s", i, err)
			}
			assert.Exactly(t, test.insertID, lid, "Index %d has different LastInsertIDs", i)
			stmtA.Reset()
		}
	})

	t.Run("ExecArgs Multi Row", func(t *testing.T) {
		dbc, dbMock := dmltest.MockDB(t)
		defer dmltest.MockClose(t, dbc, dbMock)

		prep := dbMock.ExpectPrepare(dmltest.SQLMockQuoteMeta("INSERT INTO `customer_entity` (`email`,`group_id`) VALUES (?,?),(?,?)"))
		prep.ExpectExec().WithArgs("a@b.c", 33, "d@e.f", 33).WillReturnResult(sqlmock.NewResult(6, 0))
		prep.ExpectExec().WithArgs("x@y.z", 44, "u@v.w", 44).WillReturnResult(sqlmock.NewResult(7, 0))

		stmt, err := dml.NewInsert("customer_entity").
			AddColumns("email", "group_id").
			BuildValues().
			SetRowCount(2).
			WithDB(dbc.DB).
			Prepare(context.TODO())
		assert.NoError(t, err)
		defer func() {
			assert.NoError(t, stmt.Close(), "Close on a prepared statement")
		}()

		tests := []struct {
			email1   string
			groupID1 int
			email2   string
			groupID2 int
			insertID int64
		}{
			{"a@b.c", 33, "d@e.f", 33, 6},
			{"x@y.z", 44, "u@v.w", 44, 7},
		}

		stmtA := stmt.WithArgs()
		for i, test := range tests {

			res, err := stmtA.String(test.email1).Int(test.groupID1).String(test.email2).Int(test.groupID2).
				ExecContext(context.TODO())
			if err != nil {
				t.Fatalf("Index %d => %+v", i, err)
			}
			lid, err := res.LastInsertId()
			if err != nil {
				t.Fatalf("Result index %d with error: %s", i, err)
			}
			assert.Exactly(t, test.insertID, lid, "Index %d has different LastInsertIDs", i)
			stmtA.Reset()
		}
	})

	t.Run("ExecRecord One Row", func(t *testing.T) {
		dbc, dbMock := dmltest.MockDB(t)
		defer dmltest.MockClose(t, dbc, dbMock)

		prep := dbMock.ExpectPrepare(dmltest.SQLMockQuoteMeta("INSERT INTO `dml_person` (`name`,`email`) VALUES (?,?)"))
		prep.ExpectExec().WithArgs("Peter Gopher", "peter@gopher.go").WillReturnResult(sqlmock.NewResult(4, 0))
		prep.ExpectExec().WithArgs("John Doe", "john@doe.go").WillReturnResult(sqlmock.NewResult(5, 0))

		stmt, err := dml.NewInsert("dml_person").
			AddColumns("name", "email").
			BuildValues().
			WithDB(dbc.DB).
			Prepare(context.TODO())
		assert.NoError(t, err, "failed creating a prepared statement")
		defer func() {
			assert.NoError(t, stmt.Close(), "Close on a prepared statement")
		}()

		tests := []struct {
			name     string
			email    string
			insertID int64
		}{
			{"Peter Gopher", "peter@gopher.go", 4},
			{"John Doe", "john@doe.go", 5},
		}

		for i, test := range tests {

			p := &dmlPerson{
				Name:  test.name,
				Email: null.MakeString(test.email),
			}

			res, err := stmt.WithArgs().Record("", p).ExecContext(context.TODO())
			if err != nil {
				t.Fatalf("Index %d => %+v", i, err)
			}
			lid, err := res.LastInsertId()
			if err != nil {
				t.Fatalf("Result index %d with error: %s", i, err)
			}
			assert.Exactly(t, test.insertID, lid, "Index %d has different LastInsertIDs", i)
			assert.Exactly(t, test.insertID, p.ID, "Index %d and model p has different LastInsertIDs", i)
		}
	})

	t.Run("ExecContext", func(t *testing.T) {
		dbc, dbMock := dmltest.MockDB(t)
		defer dmltest.MockClose(t, dbc, dbMock)

		prep := dbMock.ExpectPrepare(dmltest.SQLMockQuoteMeta("INSERT INTO `dml_person` (`name`,`email`) VALUES (?,?)"))
		prep.ExpectExec().WithArgs("Peter Gopher", "peter@gopher.go").WillReturnResult(sqlmock.NewResult(4, 0))

		stmt, err := dml.NewInsert("dml_person").
			AddColumns("name", "email").
			BuildValues().
			WithDB(dbc.DB).
			Prepare(context.TODO())
		assert.NoError(t, err, "failed creating a prepared statement")
		defer func() {
			assert.NoError(t, stmt.Close(), "Close on a prepared statement")
		}()

		res, err := stmt.WithArgs().ExecContext(context.TODO(), "Peter Gopher", "peter@gopher.go")
		assert.NoError(t, err, "failed to execute ExecContext")

		lid, err := res.LastInsertId()
		if err != nil {
			t.Fatal(err)
		}
		assert.Exactly(t, int64(4), lid, "Different LastInsertIDs")
	})
}

func TestInsert_BuildValues(t *testing.T) {
	t.Parallel()

	t.Run("WithArgs", func(t *testing.T) {
		p := &dmlPerson{
			Name:  "Pike",
			Email: null.MakeString("pikes@peak.co"),
		}

		insA := dml.NewInsert("alpha").
			AddColumns("name", "email").BuildValues().
			WithArgs()

		// bug

		compareToSQL(t, insA.Record("", p), errors.NoKind,
			"INSERT INTO `alpha` (`name`,`email`) VALUES (?,?)",
			"",
			"Pike", "pikes@peak.co",
		)
	})

	t.Run("WithoutArgs", func(t *testing.T) {
		ins := dml.NewInsert("alpha").
			AddColumns("name", "email").BuildValues()

		compareToSQL(t, ins, errors.NoKind,
			"INSERT INTO `alpha` (`name`,`email`) VALUES (?,?)",
			"",
		)
	})
}

func TestInsert_Clone(t *testing.T) {
	t.Parallel()

	dbc, dbMock := dmltest.MockDB(t, dml.WithLogger(log.BlackHole{}, func() string { return "uniqueID" }))
	defer dmltest.MockClose(t, dbc, dbMock)

	t.Run("nil", func(t *testing.T) {
		var i *dml.Insert
		i2 := i.Clone()
		assert.Nil(t, i)
		assert.Nil(t, i2)
	})

	t.Run("non-nil AddColumns", func(t *testing.T) {
		i := dbc.InsertInto("dml_people").AddColumns("name", "email")

		i2 := i.Clone()
		notEqualPointers(t, i, i2)
		notEqualPointers(t, i.Columns, i2.Columns)
		assert.Exactly(t, i.Pairs, i2.Pairs)
		assert.Exactly(t, i.RecordPlaceHolderCount, i2.RecordPlaceHolderCount)
		assert.Exactly(t, i.DB, i2.DB)
		assert.Exactly(t, i.Log, i2.Log)
	})
	t.Run("non-nil AddColumns", func(t *testing.T) {
		i := dbc.InsertInto("dml_people").WithPairs(
			dml.Column("name").Str("Hans"),
			dml.Column("age").Int(79),
		)

		i2 := i.Clone()
		notEqualPointers(t, i, i2)
		assert.Exactly(t, i.Columns, i2.Columns)
		notEqualPointers(t, i.Pairs, i2.Pairs)
		assert.Exactly(t, i.RecordPlaceHolderCount, i2.RecordPlaceHolderCount)
		assert.Exactly(t, i.DB, i2.DB)
		assert.Exactly(t, i.Log, i2.Log)
	})
	t.Run("non-nil OnDulicateKey", func(t *testing.T) {
		i := dml.NewInsert("a").
			AddColumns("something_id", "user_id", "other").
			AddOnDuplicateKey(
				dml.Column("something_id").Int64(99),
				dml.Column("user_id").Values(),
			)
		i2 := i.Clone()
		notEqualPointers(t, i, i2)
		assert.Exactly(t, i.Columns, i2.Columns)
		assert.False(t, i2.IsOnDuplicateKey, "Should be false i2.IsOnDuplicateKey")
		assert.Exactly(t, i.IsOnDuplicateKey, i2.IsOnDuplicateKey)
		notEqualPointers(t, i.OnDuplicateKeys, i2.OnDuplicateKeys)
	})
}
