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

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/corestoreio/errors"
	"github.com/corestoreio/log"
	"github.com/corestoreio/pkg/sql/dml"
	"github.com/corestoreio/pkg/sql/dmltest"
	"github.com/corestoreio/pkg/storage/null"
	"github.com/corestoreio/pkg/util/assert"
)

func TestDelete_Prepare(t *testing.T) {
	t.Parallel()

	t.Run("ToSQL Error", func(t *testing.T) {
		compareToSQL(t, dml.NewDelete("").Where(dml.Column("a").Int64(1)), errors.Empty, "", "")
	})

	t.Run("Prepare Error", func(t *testing.T) {
		d := &dml.Delete{
			BuilderBase: dml.BuilderBase{
				Table: dml.MakeIdentifier("table"),
			},
		}
		d.WithDB(dbMock{
			error: errors.AlreadyClosed.Newf("Who closed myself?"),
		})

		d.Where(dml.Column("a").Int(1))
		stmt, err := d.Prepare(context.TODO())
		assert.Nil(t, stmt)
		assert.ErrorIsKind(t, errors.AlreadyClosed, err)
	})

	t.Run("ExecArgs One Row", func(t *testing.T) {
		dbc, dbMock := dmltest.MockDB(t)
		defer dmltest.MockClose(t, dbc, dbMock)

		prep := dbMock.ExpectPrepare(dmltest.SQLMockQuoteMeta("DELETE FROM `customer_entity` WHERE (`email` = ?) AND (`group_id` = ?)"))
		prep.ExpectExec().WithArgs("a@b.c", 33).WillReturnResult(sqlmock.NewResult(0, 1))
		prep.ExpectExec().WithArgs("x@y.z", 44).WillReturnResult(sqlmock.NewResult(0, 2))

		stmt, err := dml.NewDelete("customer_entity").
			Where(dml.Column("email").PlaceHolder(), dml.Column("group_id").PlaceHolder()).
			WithDB(dbc.DB).
			Prepare(context.TODO())
		assert.NoError(t, err, "failed creating a prepared statement")
		defer dmltest.Close(t, stmt)

		tests := []struct {
			email   string
			groupID int
			affRows int64
		}{
			{"a@b.c", 33, 1},
			{"x@y.z", 44, 2},
		}
		args := stmt.WithDBR()
		for i, test := range tests {
			res, err := args.String(test.email).Int(test.groupID).ExecContext(context.TODO())
			if err != nil {
				t.Fatalf("Index %d => %+v", i, err)
			}
			ra, err := res.RowsAffected()
			if err != nil {
				t.Fatalf("Result index %d with error: %s", i, err)
			}
			assert.Exactly(t, test.affRows, ra, "Index %d has different RowsAffected", i)
			args.Reset()
		}
	})

	t.Run("ExecRecord One Row", func(t *testing.T) {
		dbc, dbMock := dmltest.MockDB(t)
		defer dmltest.MockClose(t, dbc, dbMock)

		prep := dbMock.ExpectPrepare(dmltest.SQLMockQuoteMeta("DELETE FROM `dml_person` WHERE (`name` = ?) AND (`email` = ?)"))
		prep.ExpectExec().WithArgs("Peter Gopher", "peter@gopher.go").WillReturnResult(sqlmock.NewResult(0, 4))
		prep.ExpectExec().WithArgs("John Doe", "john@doe.go").WillReturnResult(sqlmock.NewResult(0, 5))

		stmt, err := dml.NewDelete("dml_person").
			Where(dml.Column("name").PlaceHolder(), dml.Column("email").PlaceHolder()).
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

			res, err := stmt.WithDBR().Raw(dml.Qualify("", p)).ExecContext(context.TODO())
			if err != nil {
				t.Fatalf("Index %d => %+v", i, err)
			}
			lid, err := res.RowsAffected()
			if err != nil {
				t.Fatalf("Result index %d with error: %s", i, err)
			}
			assert.Exactly(t, test.insertID, lid, "Index %d has different RowsAffected", i)
		}
	})

	t.Run("ExecContext", func(t *testing.T) {
		dbc, dbMock := dmltest.MockDB(t)
		defer dmltest.MockClose(t, dbc, dbMock)

		prep := dbMock.ExpectPrepare(dmltest.SQLMockQuoteMeta("DELETE FROM `dml_person` WHERE (`name` = ?) AND (`email` = ?)"))
		prep.ExpectExec().WithArgs("Peter Gopher", "peter@gopher.go").WillReturnResult(sqlmock.NewResult(0, 4))

		stmt, err := dml.NewDelete("dml_person").
			Where(dml.Column("name").PlaceHolder(), dml.Column("email").PlaceHolder()).
			WithDB(dbc.DB).
			Prepare(context.TODO())
		assert.NoError(t, err, "failed creating a prepared statement")
		defer func() {
			assert.NoError(t, stmt.Close(), "Close on a prepared statement")
		}()

		res, err := stmt.WithDBR().ExecContext(context.TODO(), "Peter Gopher", "peter@gopher.go")
		assert.NoError(t, err, "failed to execute ExecContext")

		lid, err := res.RowsAffected()
		if err != nil {
			t.Fatal(err)
		}
		assert.Exactly(t, int64(4), lid, "Different RowsAffected")
	})
}

func TestDelete_Join(t *testing.T) {
	t.Parallel()

	del1 := dml.NewDelete("customer_entity").Alias("ce").
		FromTables("customer_address", "customer_company").
		Join(
			dml.MakeIdentifier("customer_company").Alias("cc"),
			dml.Columns("ce.entity_id", "cc.customer_id"),
		).
		RightJoin(
			dml.MakeIdentifier("customer_address").Alias("ca"),
			dml.Columns("ce.entity_id", "ca.parent_id"),
		).
		Where(
			dml.Column("ce.created_at").Less().PlaceHolder(),
		)

	t.Run("JOIN USING with alias", func(t *testing.T) {
		compareToSQL(t, del1, errors.NoKind,
			"DELETE `ce`,`customer_address`,`customer_company` FROM `customer_entity` AS `ce` INNER JOIN `customer_company` AS `cc` USING (`ce.entity_id`,`cc.customer_id`) RIGHT JOIN `customer_address` AS `ca` USING (`ce.entity_id`,`ca.parent_id`) WHERE (`ce`.`created_at` < ?)",
			"",
		)
	})

	t.Run("JOIN USING with alias WithDBR", func(t *testing.T) {
		compareToSQL(t, del1.WithDBR().Time(now()), errors.NoKind,
			"DELETE `ce`,`customer_address`,`customer_company` FROM `customer_entity` AS `ce` INNER JOIN `customer_company` AS `cc` USING (`ce.entity_id`,`cc.customer_id`) RIGHT JOIN `customer_address` AS `ca` USING (`ce.entity_id`,`ca.parent_id`) WHERE (`ce`.`created_at` < ?)",
			"DELETE `ce`,`customer_address`,`customer_company` FROM `customer_entity` AS `ce` INNER JOIN `customer_company` AS `cc` USING (`ce.entity_id`,`cc.customer_id`) RIGHT JOIN `customer_address` AS `ca` USING (`ce.entity_id`,`ca.parent_id`) WHERE (`ce`.`created_at` < '2006-01-02 15:04:05')",
			now(),
		)
	})

	t.Run("LeftJoin USING without alias", func(t *testing.T) {
		del := dml.NewDelete("customer_entity").
			FromTables("customer_address").
			LeftJoin(
				dml.MakeIdentifier("customer_address").Alias("ca"),
				dml.Columns("ce.entity_id", "ca.parent_id"),
			).
			Where(
				dml.Column("ce.created_at").Less().PlaceHolder(),
			)

		compareToSQL(t, del, errors.NoKind,
			"DELETE `customer_entity`,`customer_address` FROM `customer_entity` LEFT JOIN `customer_address` AS `ca` USING (`ce.entity_id`,`ca.parent_id`) WHERE (`ce`.`created_at` < ?)",
			"",
		)
	})

	t.Run("OuterJoin USING without alias", func(t *testing.T) {
		del := dml.NewDelete("customer_entity").
			FromTables("customer_address").
			OuterJoin(
				dml.MakeIdentifier("customer_address").Alias("ca"),
				dml.Columns("ce.entity_id", "ca.parent_id"),
			)

		compareToSQL(t, del, errors.NoKind,
			"DELETE `customer_entity`,`customer_address` FROM `customer_entity` OUTER JOIN `customer_address` AS `ca` USING (`ce.entity_id`,`ca.parent_id`)",
			"",
		)
	})

	t.Run("JOIN USING without FromTables", func(t *testing.T) {
		del := dml.NewDelete("customer_entity").
			CrossJoin(
				dml.MakeIdentifier("customer_address").Alias("ca"),
				dml.Column("ce.entity_id").Equal().Column("ca.parent_id"),
			).
			Where(
				dml.Column("ce.created_at").Less().PlaceHolder(),
			)

		compareToSQL(t, del, errors.NoKind,
			"DELETE FROM `customer_entity` CROSS JOIN `customer_address` AS `ca` ON (`ce`.`entity_id` = `ca`.`parent_id`) WHERE (`ce`.`created_at` < ?)",
			"",
		)
	})
}

func TestDelete_Returning(t *testing.T) {
	t.Parallel()

	t.Run("not allowed", func(t *testing.T) {
		del := dml.NewDelete("customer_entity").
			FromTables("customer_address").
			OuterJoin(
				dml.MakeIdentifier("customer_address").Alias("ca"),
				dml.Columns("ce.entity_id", "ca.parent_id"),
			)
		del.Returning = dml.NewSelect()
		compareToSQL(t, del, errors.NotAllowed,
			"",
			"",
		)
	})

	t.Run("return delete rows", func(t *testing.T) {
		del := dml.NewDelete("customer_entity").
			Where(
				dml.Column("ce.entity_id").GreaterOrEqual().PlaceHolder(),
			)
		del.Returning = dml.NewSelect("entity_id", "created_at").From("customer_entity")
		compareToSQL(t, del, errors.NoKind,
			"DELETE FROM `customer_entity` WHERE (`ce`.`entity_id` >= ?) RETURNING SELECT `entity_id`, `created_at` FROM `customer_entity`",
			"",
		)
	})
}

func TestDelete_Clone(t *testing.T) {
	t.Parallel()

	dbc, dbMock := dmltest.MockDB(t, dml.WithLogger(log.BlackHole{}, func() string { return "uniqueID" }))
	defer dmltest.MockClose(t, dbc, dbMock)

	t.Run("nil", func(t *testing.T) {
		var d *dml.Delete
		d2 := d.Clone()
		assert.Nil(t, d)
		assert.Nil(t, d2)
	})

	t.Run("non-nil", func(t *testing.T) {
		d := dbc.DeleteFrom("dml_people").Alias("dmlPpl").FromTables("a1", "b2").
			Where(
				dml.Column("id").PlaceHolder(),
			).OrderBy("id")
		d2 := d.Clone()
		notEqualPointers(t, d, d2)
		notEqualPointers(t, d, d2)
		notEqualPointers(t, d.BuilderConditional.Wheres, d2.BuilderConditional.Wheres)
		notEqualPointers(t, d.BuilderConditional.OrderBys, d2.BuilderConditional.OrderBys)
		notEqualPointers(t, d.MultiTables, d2.MultiTables)
		assert.Exactly(t, d.DB, d2.DB)
		assert.Exactly(t, d.Log, d2.Log)
	})
}
