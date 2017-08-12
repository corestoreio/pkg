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
	"testing"

	"github.com/corestoreio/errors"
	"github.com/corestoreio/log"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var _ ArgumentsAppender = (*someRecord)(nil)

type someRecord struct {
	SomethingID int
	UserID      int64
	Other       bool
}

func (sr someRecord) AppendArguments(stmtType SQLStmt, args Arguments, columns []string) (Arguments, error) {
	for _, c := range columns {
		switch c {
		case "something_id":
			args = args.Int(sr.SomethingID)
		case "user_id":
			args = args.Int64(sr.UserID)
		case "other":
			args = args.Bool(sr.Other)
		default:
			return nil, errors.NewNotFoundf("[dbr_test] Column %q not found", c)
		}
	}
	if len(columns) == 0 && stmtType.IsValues() {
		args = args.Int(sr.SomethingID).Int64(sr.UserID).Bool(sr.Other)

	}
	return args, nil
}

func TestInsert_SetValuesCount(t *testing.T) {
	t.Parallel()

	t.Run("not set", func(t *testing.T) {
		compareToSQL(t,
			NewInsert("a").AddColumns("b", "c"),
			nil,
			"INSERT INTO `a` (`b`,`c`) VALUES (?,?)",
			"",
		)
	})
	t.Run("set to two", func(t *testing.T) {
		compareToSQL(t,
			NewInsert("a").AddColumns("b", "c").SetRowCount(2),
			nil,
			"INSERT INTO `a` (`b`,`c`) VALUES (?,?),(?,?)",
			"",
		)
	})
	t.Run("with values", func(t *testing.T) {
		compareToSQL(t,
			NewInsert("dbr_people").AddColumns("name", "key").AddValues("Barack", "44"),
			nil,
			"INSERT INTO `dbr_people` (`name`,`key`) VALUES (?,?)",
			"INSERT INTO `dbr_people` (`name`,`key`) VALUES ('Barack','44')",
			"Barack", "44",
		)
	})
	t.Run("with record", func(t *testing.T) {
		person := dbrPerson{Name: "Barack"}
		person.Email.Valid = true
		person.Email.String = "obama@whitehouse.gov"
		compareToSQL(t,
			NewInsert("dbr_people").AddColumns("name", "email").AddRecords(&person),
			nil,
			"INSERT INTO `dbr_people` (`name`,`email`) VALUES (?,?)",
			"INSERT INTO `dbr_people` (`name`,`email`) VALUES ('Barack','obama@whitehouse.gov')",
			"Barack", "obama@whitehouse.gov",
		)
	})
}

func TestInsert_Add(t *testing.T) {
	t.Parallel()
	t.Run("AddValues error", func(t *testing.T) {
		defer func() {
			if r := recover(); r != nil {
				if err, ok := r.(error); ok {
					assert.True(t, errors.IsNotSupported(err), "%+v", err)
				} else {
					t.Errorf("Panic should contain an error but got:\n%+v", r)
				}
			} else {
				t.Error("Expecting a panic but got nothing")
			}
		}()
		NewInsert("a").AddColumns("b").AddValues(make(chan int))
	})
	t.Run("single AddValues", func(t *testing.T) {
		compareToSQL(t,
			NewInsert("a").AddColumns("b", "c").AddValues(1, 2),
			nil,
			"INSERT INTO `a` (`b`,`c`) VALUES (?,?)",
			"INSERT INTO `a` (`b`,`c`) VALUES (1,2)",
			int64(1), int64(2),
		)
	})
	t.Run("multi AddValues on duplicate key", func(t *testing.T) {
		compareToSQL(t,
			NewInsert("a").AddColumns("b", "c").
				AddValues(
					1, 2,
					3, 4,
				).
				AddValues(
					5, 6,
				).
				OnDuplicateKey(),
			nil,
			"INSERT INTO `a` (`b`,`c`) VALUES (?,?),(?,?),(?,?) ON DUPLICATE KEY UPDATE `b`=VALUES(`b`), `c`=VALUES(`c`)",
			"INSERT INTO `a` (`b`,`c`) VALUES (1,2),(3,4),(5,6) ON DUPLICATE KEY UPDATE `b`=VALUES(`b`), `c`=VALUES(`c`)",
			int64(1), int64(2), int64(3), int64(4), int64(5), int64(6),
		)
	})
	t.Run("single AddArguments", func(t *testing.T) {
		compareToSQL(t,
			NewInsert("a").AddColumns("b", "c").AddArguments(MakeArgs(2).Int64(1).Int64(2)),
			nil,
			"INSERT INTO `a` (`b`,`c`) VALUES (?,?)",
			"INSERT INTO `a` (`b`,`c`) VALUES (1,2)",
			int64(1), int64(2),
		)
	})
	t.Run("multi AddArguments on duplicate key", func(t *testing.T) {
		compareToSQL(t,
			NewInsert("a").AddColumns("b", "c").
				AddArguments(MakeArgs(4).
					Int64(1).Int64(2).
					Int64(3).Int64(4),
				).
				AddArguments(MakeArgs(2).
					Int64(5).Int64(6),
				).
				OnDuplicateKey(),
			nil,
			"INSERT INTO `a` (`b`,`c`) VALUES (?,?),(?,?),(?,?) ON DUPLICATE KEY UPDATE `b`=VALUES(`b`), `c`=VALUES(`c`)",
			"INSERT INTO `a` (`b`,`c`) VALUES (1,2),(3,4),(5,6) ON DUPLICATE KEY UPDATE `b`=VALUES(`b`), `c`=VALUES(`c`)",
			int64(1), int64(2), int64(3), int64(4), int64(5), int64(6),
		)
	})
}

func TestInsert_AddRecords(t *testing.T) {
	t.Parallel()
	objs := []someRecord{{1, 88, false}, {2, 99, true}, {3, 101, true}}
	wantArgs := []interface{}{int64(1), int64(88), false, int64(2), int64(99), true, int64(3), int64(101), true, int64(99)}

	t.Run("valid with multiple records", func(t *testing.T) {
		compareToSQL(t,
			NewInsert("a").
				AddColumns("something_id", "user_id", "other").
				AddRecords(objs[0]).AddRecords(objs[1], objs[2]).
				AddOnDuplicateKey(
					Column("something_id").Int64(99),
					Column("user_id").Values(),
				),
			nil,
			"INSERT INTO `a` (`something_id`,`user_id`,`other`) VALUES (?,?,?),(?,?,?),(?,?,?) ON DUPLICATE KEY UPDATE `something_id`=?, `user_id`=VALUES(`user_id`)",
			"INSERT INTO `a` (`something_id`,`user_id`,`other`) VALUES (1,88,0),(2,99,1),(3,101,1) ON DUPLICATE KEY UPDATE `something_id`=99, `user_id`=VALUES(`user_id`)",
			wantArgs...,
		)
	})
	t.Run("without columns, all columns requested", func(t *testing.T) {
		compareToSQL(t,
			NewInsert("a").
				SetRecordValueCount(3).
				AddRecords(objs[0]).AddRecords(objs[1], objs[2]).
				AddOnDuplicateKey(
					Column("something_id").Int64(99),
					Column("user_id").Values(),
				),
			nil,
			"INSERT INTO `a` VALUES (?,?,?),(?,?,?),(?,?,?) ON DUPLICATE KEY UPDATE `something_id`=?, `user_id`=VALUES(`user_id`)",
			"INSERT INTO `a` VALUES (1,88,0),(2,99,1),(3,101,1) ON DUPLICATE KEY UPDATE `something_id`=99, `user_id`=VALUES(`user_id`)",
			wantArgs...,
		)
	})
	t.Run("column not found", func(t *testing.T) {
		objs := []someRecord{{1, 88, false}, {2, 99, true}}
		compareToSQL(t,
			NewInsert("a").AddColumns("something_it", "user_id", "other").AddRecords(objs[0]).AddRecords(objs[1]),
			errors.IsNotFound,
			"",
			"",
		)
	})
	t.Run("slice as record - nice feature", func(t *testing.T) {
		wantArgs := []interface{}{"Muffin Hat", "Muffin@Hat.head", "Marianne Phyllis Finch", "marianne@phyllis.finch", "Daphne Augusta Perry", "daphne@augusta.perry"}
		persons := &dbrPersons{
			Data: []*dbrPerson{
				{Name: "Muffin Hat", Email: MakeNullString("Muffin@Hat.head")},
				{Name: "Marianne Phyllis Finch", Email: MakeNullString("marianne@phyllis.finch")},
				{Name: "Daphne Augusta Perry", Email: MakeNullString("daphne@augusta.perry")},
			},
		}

		compareToSQL(t,
			NewInsert("dbr_person").
				AddColumns("name", "email").
				AddRecords(persons).
				SetRowCount(len(persons.Data)),
			nil,
			"INSERT INTO `dbr_person` (`name`,`email`) VALUES (?,?),(?,?),(?,?)",
			"INSERT INTO `dbr_person` (`name`,`email`) VALUES ('Muffin Hat','Muffin@Hat.head'),('Marianne Phyllis Finch','marianne@phyllis.finch'),('Daphne Augusta Perry','daphne@augusta.perry')",
			wantArgs...,
		)
	})

}

func TestInsertKeywordColumnName(t *testing.T) {
	// Insert a column whose name is reserved
	s := createRealSessionWithFixtures(t, nil)
	res, err := s.InsertInto("dbr_people").AddColumns("name", "key").AddValues("Barack", "44").Exec(context.TODO())
	require.NoError(t, err)

	rowsAff, err := res.RowsAffected()
	require.NoError(t, err)
	assert.Equal(t, rowsAff, int64(1))
}

func TestInsertReal(t *testing.T) {
	// Insert by specifying values
	s := createRealSessionWithFixtures(t, nil)
	res, err := s.InsertInto("dbr_people").AddColumns("name", "email").AddValues("Barack", "obama@whitehouse.gov").Exec(context.TODO())
	validateInsertingBarack(t, s, res, err)

	// Insert by specifying a record (ptr to struct)
	s = createRealSessionWithFixtures(t, nil)
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
	require.NoError(t, err)
	if res == nil {
		t.Fatal("result at nit but should not")
	}
	id, err := res.LastInsertId()
	require.NoError(t, err)
	rowsAff, err := res.RowsAffected()
	require.NoError(t, err)

	assert.True(t, id > 0)
	assert.Equal(t, int64(1), rowsAff)

	var person dbrPerson
	_, err = c.Select("*").From("dbr_people").Where(Column("id").Int64(id)).Load(context.TODO(), &person)
	require.NoError(t, err)

	assert.Equal(t, id, int64(person.ID))
	assert.Equal(t, "Barack", person.Name)
	assert.Equal(t, true, person.Email.Valid)
	assert.Equal(t, "obama@whitehouse.gov", person.Email.String)
}

func TestInsertReal_OnDuplicateKey(t *testing.T) {

	s := createRealSessionWithFixtures(t, nil)
	res, err := s.InsertInto("dbr_people").
		AddColumns("id", "name", "email").
		AddValues(678, "Pike", "pikes@peak.co").Exec(context.TODO())
	if err != nil {
		t.Fatalf("%+v", err)
	}
	inID, err := res.LastInsertId()
	if err != nil {
		t.Fatalf("%+v", err)
	}
	{
		var p dbrPerson
		_, err = s.Select("*").From("dbr_people").Where(Column("id").Int64(inID)).Load(context.TODO(), &p)
		require.NoError(t, err)
		assert.Equal(t, "Pike", p.Name)
		assert.Equal(t, "pikes@peak.co", p.Email.String)
	}
	res, err = s.InsertInto("dbr_people").
		AddColumns("id", "name", "email").
		AddValues(inID, "", "pikes@peak.com").
		AddOnDuplicateKey(Column("name").Str("Pik3"), Column("email").Values()).
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
		_, err = s.Select("*").From("dbr_people").Where(Column("id").Int64(inID)).Load(context.TODO(), &p)
		require.NoError(t, err)
		assert.Equal(t, "Pik3", p.Name)
		assert.Equal(t, "pikes@peak.com", p.Email.String)
	}
}

// TODO: do a real test inserting multiple records

func TestInsert_Prepare(t *testing.T) {
	t.Parallel()
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
		in.DB = dbMock{
			error: errors.NewAlreadyClosedf("Who closed myself?"),
		}
		in.AddColumns("a", "b").AddValues(1, true)

		stmt, err := in.Prepare(context.TODO())
		assert.Nil(t, stmt)
		assert.True(t, errors.IsAlreadyClosed(err), "%+v", err)
	})
}

func TestInsert_Events(t *testing.T) {
	t.Parallel()

	t.Run("Stop Propagation", func(t *testing.T) {
		d := NewInsert("tableA")
		d.AddColumns("a", "b").AddValues(1, true)

		d.Log = log.BlackHole{EnableInfo: true, EnableDebug: true}
		d.Listeners.Add(
			Listen{
				Name:      "listener1",
				EventType: OnBeforeToSQL,
				InsertFunc: func(b *Insert) {
					b.Pair(Column("col1").Str("X1"))
				},
			},
			Listen{
				Name:      "listener2",
				EventType: OnBeforeToSQL,
				InsertFunc: func(b *Insert) {
					b.Pair(Column("col2").Str("X2"))
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

		sqlStr, args, err := d.Interpolate().ToSQL()
		require.NoError(t, err)
		assert.Nil(t, args)
		assert.Exactly(t, "INSERT INTO `tableA` (`a`,`b`,`col1`,`col2`) VALUES (1,1,'X1','X2')", sqlStr)

		// call it twice (4x) to test for being NOT idempotent
		compareToSQL(t, d, nil,
			"INSERT INTO `tableA` (`a`,`b`,`col1`,`col2`) VALUES (1,1,'X1','X2')",
			"",
		)

	})

	t.Run("Missing EventType", func(t *testing.T) {
		ins := NewInsert("tableA")
		ins.AddColumns("a", "b").AddValues(1, true)

		ins.Listeners.Add(
			Listen{
				Name: "colC",
				InsertFunc: func(i *Insert) {
					i.Pair(Column("colC").Str("X1"))
				},
			},
		)
		compareToSQL(t, ins, errors.IsEmpty,
			"",
			"",
		)
	})

	t.Run("Should Dispatch", func(t *testing.T) {
		ins := NewInsert("tableA")

		ins.AddColumns("a", "b").AddValues(1, true)

		ins.Listeners.Add(
			Listen{
				EventType: OnBeforeToSQL,
				Name:      "colA",
				Once:      true,
				InsertFunc: func(i *Insert) {
					i.Pair(Column("colA").Float64(3.14159))
				},
			},
			Listen{
				EventType: OnBeforeToSQL,
				Name:      "colB",
				Once:      true,
				InsertFunc: func(i *Insert) {
					i.Pair(Column("colB").Float64(2.7182))
				},
			},
		)

		ins.Listeners.Add(
			Listen{
				// Multiple calls and colC is getting ignored because the Pair
				// function only creates the next value slice when a column `a`
				// gets called with Pair.
				EventType: OnBeforeToSQL,
				Name:      "colC",
				InsertFunc: func(i *Insert) {
					i.Pair(Column("colC").Str("X1"))
				},
			},
		)
		sqlStr, args, err := ins.Interpolate().ToSQL()
		require.NoError(t, err)
		assert.Nil(t, args)
		assert.Exactly(t, "INSERT INTO `tableA` (`a`,`b`,`colA`,`colB`,`colC`) VALUES (1,1,3.14159,2.7182,'X1')", sqlStr)

		compareToSQL(t, ins, nil,
			"INSERT INTO `tableA` (`a`,`b`,`colA`,`colB`,`colC`) VALUES (1,1,3.14159,2.7182,'X1')",
			"",
		)

		assert.Exactly(t, `colA; colB; colC`, ins.Listeners.String())
	})
}

func TestInsert_FromSelect(t *testing.T) {
	t.Parallel()

	ins := NewInsert("tableA")
	// columns and args just to check that they get ignored
	ins.AddColumns("a", "b").AddValues(1, true)

	compareToSQL(t, ins.FromSelect(NewSelect("something_id", "user_id", "other").
		From("some_table").
		Where(
			ParenthesisOpen(),
			Column("d").Int64(1),
			Column("e").Str("wat").Or(),
			ParenthesisClose(),
			Column("a").In().Int64s(1, 2, 3),
		).
		OrderByDesc("id").
		Paginate(1, 20)),
		nil,
		"INSERT INTO `tableA` SELECT `something_id`, `user_id`, `other` FROM `some_table` WHERE ((`d` = ?) OR (`e` = ?)) AND (`a` IN (?,?,?)) ORDER BY `id` DESC LIMIT 20 OFFSET 0",
		"INSERT INTO `tableA` SELECT `something_id`, `user_id`, `other` FROM `some_table` WHERE ((`d` = 1) OR (`e` = 'wat')) AND (`a` IN (1,2,3)) ORDER BY `id` DESC LIMIT 20 OFFSET 0",
		int64(1), "wat", int64(1), int64(2), int64(3),
	)
}

func TestInsert_Replace_Ignore(t *testing.T) {
	t.Parallel()

	// this generated statement does not comply the SQL standard
	compareToSQL(t, NewInsert("a").
		Replace().Ignore().
		AddColumns("b", "c").
		AddValues(1, 2).
		AddValues(3, 4),
		nil,
		"REPLACE IGNORE INTO `a` (`b`,`c`) VALUES (?,?),(?,?)",
		"REPLACE IGNORE INTO `a` (`b`,`c`) VALUES (1,2),(3,4)",
		int64(1), int64(2), int64(3), int64(4),
	)
}

func TestInsert_WithoutColumns(t *testing.T) {
	t.Parallel()

	t.Run("each column in its own Arg", func(t *testing.T) {
		compareToSQL(t, NewInsert("catalog_product_link").
			AddValues(2046, 33, 3).
			AddValues(2046, 34, 3).
			AddValues(2046, 35, 3),
			nil,
			"INSERT INTO `catalog_product_link` VALUES (?,?,?),(?,?,?),(?,?,?)",
			"INSERT INTO `catalog_product_link` VALUES (2046,33,3),(2046,34,3),(2046,35,3)",
			int64(2046), int64(33), int64(3), int64(2046), int64(34), int64(3), int64(2046), int64(35), int64(3),
		)
	})
}

func TestInsert_Pair(t *testing.T) {
	t.Parallel()

	t.Run("one row", func(t *testing.T) {
		compareToSQL(t, NewInsert("catalog_product_link").
			Pair(
				Column("product_id").Int64(2046),
				Column("linked_product_id").Int64(33),
				Column("link_type_id").Int64(3),
			),
			nil,
			"INSERT INTO `catalog_product_link` (`product_id`,`linked_product_id`,`link_type_id`) VALUES (?,?,?)",
			"INSERT INTO `catalog_product_link` (`product_id`,`linked_product_id`,`link_type_id`) VALUES (2046,33,3)",
			int64(2046), int64(33), int64(3),
		)
	})
	// TODO implement expression handling, requires some refactorings
	//t.Run("expression no args", func(t *testing.T) {
	//	compareToSQL(t, NewInsert("catalog_product_link").
	//		Pair(
	//			Column("product_id").Int64(2046),
	//			Column("type_name").Expression("CONCAT(`product_id`,'Manufacturer')"),
	//			Column("link_type_id").Int64(3),
	//		),
	//		nil,
	//		"INSERT INTO `catalog_product_link` (`product_id`,`linked_product_id`,`link_type_id`) VALUES (?,CONCAT(`product_id`,'Manufacturer'),?)",
	//		"INSERT INTO `catalog_product_link` (`product_id`,`linked_product_id`,`link_type_id`) VALUES (2046,CONCAT(`product_id`,'Manufacturer'),3)",
	//		int64(2046), int64(33), int64(3),
	//	)
	//})
	t.Run("multiple rows triggers NO error", func(t *testing.T) {
		compareToSQL(t, NewInsert("catalog_product_link").
			Pair(
				// First row
				Column("product_id").Int64(2046),
				Column("linked_product_id").Int64(33),
				Column("link_type_id").Int64(3),

				// second row
				Column("product_id").Int64(2046),
				Column("linked_product_id").Int64(34),
				Column("link_type_id").Int64(3),
			),
			nil,
			"INSERT INTO `catalog_product_link` (`product_id`,`linked_product_id`,`link_type_id`) VALUES (?,?,?),(?,?,?)",
			"INSERT INTO `catalog_product_link` (`product_id`,`linked_product_id`,`link_type_id`) VALUES (2046,33,3),(2046,34,3)",
			int64(2046), int64(33), int64(3), int64(2046), int64(34), int64(3),
		)
	})
}

func TestInsert_UseBuildCache(t *testing.T) {
	t.Parallel()

	ins := NewInsert("a").AddColumns("b", "c").
		AddValues(
			1, 2,
			3, 4,
		).
		AddValues(
			5, 6,
		).
		OnDuplicateKey().BuildCache()

	const cachedSQLPlaceHolder = "INSERT INTO `a` (`b`,`c`) VALUES (?,?),(?,?),(?,?) ON DUPLICATE KEY UPDATE `b`=VALUES(`b`), `c`=VALUES(`c`)"
	t.Run("without interpolate", func(t *testing.T) {
		for i := 0; i < 3; i++ {
			sql, args, err := ins.ToSQL()
			require.NoError(t, err, "%+v", err)
			require.Equal(t, cachedSQLPlaceHolder, sql)
			assert.Equal(t, []interface{}{int64(1), int64(2), int64(3), int64(4), int64(5), int64(6)}, args)
			assert.Equal(t, cachedSQLPlaceHolder, string(ins.cacheSQL))
		}
	})

	t.Run("with interpolate", func(t *testing.T) {
		ins.Interpolate()
		ins.cacheSQL = nil

		const cachedSQLInterpolated = "INSERT INTO `a` (`b`,`c`) VALUES (1,2),(3,4),(5,6) ON DUPLICATE KEY UPDATE `b`=VALUES(`b`), `c`=VALUES(`c`)"
		for i := 0; i < 3; i++ {
			sql, args, err := ins.ToSQL()
			assert.Equal(t, cachedSQLPlaceHolder, string(ins.cacheSQL))
			require.NoError(t, err, "%+v", err)
			require.Equal(t, cachedSQLInterpolated, sql)
			assert.Nil(t, args)
		}
	})
}

func TestInsert_AddUpdateAllNonPrimary(t *testing.T) {
	t.Parallel()

	t.Run("AddOnDuplicateKeyExclude only", func(t *testing.T) {
		compareToSQL(t, NewInsert("customer_gr1d_flat").
			AddColumns("entity_id", "name", "email", "group_id", "created_at", "website_id").
			AddValues(1, "Martin", "martin@go.go", 3, "2019-01-01", 2).
			AddOnDuplicateKeyExclude("entity_id"),
			nil,
			"INSERT INTO `customer_gr1d_flat` (`entity_id`,`name`,`email`,`group_id`,`created_at`,`website_id`) VALUES (?,?,?,?,?,?) ON DUPLICATE KEY UPDATE `name`=VALUES(`name`), `email`=VALUES(`email`), `group_id`=VALUES(`group_id`), `created_at`=VALUES(`created_at`), `website_id`=VALUES(`website_id`)",
			"INSERT INTO `customer_gr1d_flat` (`entity_id`,`name`,`email`,`group_id`,`created_at`,`website_id`) VALUES (1,'Martin','martin@go.go',3,'2019-01-01',2) ON DUPLICATE KEY UPDATE `name`=VALUES(`name`), `email`=VALUES(`email`), `group_id`=VALUES(`group_id`), `created_at`=VALUES(`created_at`), `website_id`=VALUES(`website_id`)",
			int64(1), "Martin", "martin@go.go", int64(3), "2019-01-01", int64(2),
		)
	})

	t.Run("AddOnDuplicateKeyExclude plus custom field value", func(t *testing.T) {
		compareToSQL(t, NewInsert("customer_gr1d_flat").
			AddColumns("entity_id", "name", "email", "group_id", "created_at", "website_id").
			AddValues(1, "Martin", "martin@go.go", 3, "2019-01-01", 2).
			AddOnDuplicateKeyExclude("entity_id").
			AddOnDuplicateKey(Column("created_at").Time(now())),
			nil,
			"INSERT INTO `customer_gr1d_flat` (`entity_id`,`name`,`email`,`group_id`,`created_at`,`website_id`) VALUES (?,?,?,?,?,?) ON DUPLICATE KEY UPDATE `name`=VALUES(`name`), `email`=VALUES(`email`), `group_id`=VALUES(`group_id`), `website_id`=VALUES(`website_id`), `created_at`=?",
			"INSERT INTO `customer_gr1d_flat` (`entity_id`,`name`,`email`,`group_id`,`created_at`,`website_id`) VALUES (1,'Martin','martin@go.go',3,'2019-01-01',2) ON DUPLICATE KEY UPDATE `name`=VALUES(`name`), `email`=VALUES(`email`), `group_id`=VALUES(`group_id`), `website_id`=VALUES(`website_id`), `created_at`='2006-01-02 15:04:05'",
			int64(1), "Martin", "martin@go.go", int64(3), "2019-01-01", int64(2), now(),
		)
	})

	t.Run("OnDuplicateKey enabled for all columns", func(t *testing.T) {
		ins := NewInsert("customer_gr1d_flat").
			AddColumns("name", "email", "group_id", "created_at", "website_id").
			AddValues("Martin", "martin@go.go", 3, "2019-01-01", 2).
			OnDuplicateKey()
		compareToSQL(t, ins, nil,
			"INSERT INTO `customer_gr1d_flat` (`name`,`email`,`group_id`,`created_at`,`website_id`) VALUES (?,?,?,?,?) ON DUPLICATE KEY UPDATE `name`=VALUES(`name`), `email`=VALUES(`email`), `group_id`=VALUES(`group_id`), `created_at`=VALUES(`created_at`), `website_id`=VALUES(`website_id`)",
			"INSERT INTO `customer_gr1d_flat` (`name`,`email`,`group_id`,`created_at`,`website_id`) VALUES ('Martin','martin@go.go',3,'2019-01-01',2) ON DUPLICATE KEY UPDATE `name`=VALUES(`name`), `email`=VALUES(`email`), `group_id`=VALUES(`group_id`), `created_at`=VALUES(`created_at`), `website_id`=VALUES(`website_id`)",
			"Martin", "martin@go.go", int64(3), "2019-01-01", int64(2),
		)
		// testing for being idempotent
		compareToSQL(t, ins, nil,
			"INSERT INTO `customer_gr1d_flat` (`name`,`email`,`group_id`,`created_at`,`website_id`) VALUES (?,?,?,?,?) ON DUPLICATE KEY UPDATE `name`=VALUES(`name`), `email`=VALUES(`email`), `group_id`=VALUES(`group_id`), `created_at`=VALUES(`created_at`), `website_id`=VALUES(`website_id`)",
			"INSERT INTO `customer_gr1d_flat` (`name`,`email`,`group_id`,`created_at`,`website_id`) VALUES ('Martin','martin@go.go',3,'2019-01-01',2) ON DUPLICATE KEY UPDATE `name`=VALUES(`name`), `email`=VALUES(`email`), `group_id`=VALUES(`group_id`), `created_at`=VALUES(`created_at`), `website_id`=VALUES(`website_id`)",
			"Martin", "martin@go.go", int64(3), "2019-01-01", int64(2),
		)
	})

}
