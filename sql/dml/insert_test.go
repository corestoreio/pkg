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

package dml

import (
	"context"
	"sync"
	"testing"

	"github.com/corestoreio/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestInsert_NoArguments(t *testing.T) {
	t.Parallel()

	t.Run("ToSQL", func(t *testing.T) {
		ins := NewInsert("tableA").AddColumns("a", "b")
		compareToSQL2(t, ins, errors.NoKind,
			"INSERT INTO `tableA` (`a`,`b`) VALUES (?,?)")
	})
	t.Run("WithArgs.ToSQL", func(t *testing.T) {
		ins := NewInsert("tableA").AddColumns("a", "b").WithArgs()
		compareToSQL2(t, ins, errors.NoKind,
			"INSERT INTO `tableA` (`a`,`b`) VALUES (?,?)")
	})

}

func TestInsert_SetValuesCount(t *testing.T) {
	t.Parallel()

	t.Run("not set", func(t *testing.T) {
		ins := NewInsert("a").AddColumns("b", "c")
		compareToSQL2(t, ins, errors.NoKind,
			"INSERT INTO `a` (`b`,`c`) VALUES (?,?)",
		)
		assert.Exactly(t, []string{"b", "c"}, ins.qualifiedColumns)
	})
	t.Run("set to two", func(t *testing.T) {
		compareToSQL2(t,
			NewInsert("a").AddColumns("b", "c").SetRowCount(2),
			errors.NoKind,
			"INSERT INTO `a` (`b`,`c`) VALUES (?,?),(?,?)",
		)
	})
	t.Run("with values", func(t *testing.T) {
		ins := NewInsert("dml_people").AddColumns("name", "key")
		inA := ins.WithArgs("Barack", "44")
		compareToSQL2(t, inA, errors.NoKind, "INSERT INTO `dml_people` (`name`,`key`) VALUES (?,?)",
			"Barack", "44",
		)
		assert.Exactly(t, []string{"name", "key"}, ins.qualifiedColumns)
	})
	t.Run("with record", func(t *testing.T) {
		person := dmlPerson{Name: "Barack"}
		person.Email.Valid = true
		person.Email.String = "obama@whitehouse.gov"
		compareToSQL2(t,
			NewInsert("dml_people").AddColumns("name", "email").WithArgs().Record("", &person),
			errors.NoKind,
			"INSERT INTO `dml_people` (`name`,`email`) VALUES (?,?)",
			"Barack", "obama@whitehouse.gov",
		)
	})
}

func TestInsertKeywordColumnName(t *testing.T) {
	// Insert a column whose name is reserved
	s := createRealSessionWithFixtures(t, nil)
	defer testCloser(t, s)
	ins := s.InsertInto("dml_people").AddColumns("name", "key").WithArgs("Barack", "44")

	compareExecContext(t, ins, 0, 1)
}

func TestInsertReal(t *testing.T) {
	// Insert by specifying values
	s := createRealSessionWithFixtures(t, nil)
	defer testCloser(t, s)
	ins := s.InsertInto("dml_people").AddColumns("name", "email").WithArgs("Barack", "obama@whitehouse.gov")
	lastInsertID, _ := compareExecContext(t, ins, 3, 0)
	validateInsertingBarack(t, s, lastInsertID)

	// Insert by specifying a record (ptr to struct)

	person := dmlPerson{Name: "Barack"}
	person.Email.Valid = true
	person.Email.String = "obama@whitehouse.gov"
	ins = s.InsertInto("dml_people").AddColumns("name", "email").WithArgs().Record("", &person)
	lastInsertID, _ = compareExecContext(t, ins, 4, 0)

	validateInsertingBarack(t, s, lastInsertID)
}

func validateInsertingBarack(t *testing.T, c *ConnPool, lastInsertID int64) {

	var person dmlPerson
	_, err := c.SelectFrom("dml_people").Star().Where(Column("id").Int64(lastInsertID)).WithArgs().Load(context.TODO(), &person)
	require.NoError(t, err)

	require.Equal(t, lastInsertID, int64(person.ID))
	assert.Equal(t, "Barack", person.Name)
	assert.Equal(t, true, person.Email.Valid)
	assert.Equal(t, "obama@whitehouse.gov", person.Email.String)
}

func TestInsertReal_OnDuplicateKey(t *testing.T) {

	s := createRealSessionWithFixtures(t, nil)
	defer testCloser(t, s)

	p := &dmlPerson{
		Name:  "Pike",
		Email: MakeNullString("pikes@peak.co"),
	}

	res, err := s.InsertInto("dml_people").
		AddColumns("name", "email").
		WithArgs().Record("", p).ExecContext(context.TODO())
	if err != nil {
		t.Fatalf("%+v", err)
	}
	require.Exactly(t, uint64(3), p.ID, "Last Insert ID must be three")

	inID, err := res.LastInsertId()
	if err != nil {
		t.Fatalf("%+v", err)
	}
	{
		var p dmlPerson
		_, err = s.SelectFrom("dml_people").Star().Where(Column("id").Int64(inID)).WithArgs().Load(context.TODO(), &p)
		require.NoError(t, err)
		assert.Equal(t, "Pike", p.Name)
		assert.Equal(t, "pikes@peak.co", p.Email.String)
	}

	p.Name = "-"
	p.Email.String = "pikes@peak.com"
	res, err = s.InsertInto("dml_people").
		AddColumns("id", "name", "email").
		AddOnDuplicateKey(Column("name").Str("Pik3"), Column("email").Values()).
		WithArgs().Record("", p).
		ExecContext(context.TODO())
	if err != nil {
		t.Fatalf("%+v", err)
	}
	inID2, err := res.LastInsertId()
	if err != nil {
		t.Fatalf("%+v", err)
	}
	assert.Exactly(t, inID, inID2)

	{
		var p dmlPerson
		_, err = s.SelectFrom("dml_people").Star().Where(Column("id").Int64(inID)).WithArgs().Load(context.TODO(), &p)
		require.NoError(t, err)
		assert.Equal(t, "Pik3", p.Name)
		assert.Equal(t, "pikes@peak.com", p.Email.String)
	}
}

func TestInsert_Events(t *testing.T) {
	t.Parallel()

	t.Run("Stop Propagation", func(t *testing.T) {
		d := NewInsert("tableA").AddColumns("a", "b")

		d.Listeners.Add(
			Listen{
				Name:      "listener1",
				EventType: OnBeforeToSQL,
				ListenInsertFn: func(b *Insert) {
					b.AddColumns("col1")
				},
			},
			Listen{
				Name:      "listener2",
				EventType: OnBeforeToSQL,
				ListenInsertFn: func(b *Insert) {
					b.AddColumns("col2")
					b.PropagationStopped = true
				},
			},
			Listen{
				Name:      "listener3",
				EventType: OnBeforeToSQL,
				ListenInsertFn: func(b *Insert) {
					panic("Should not get called")
				},
			},
		)

		da := d.WithArgs().Int(1).Bool(true).String("X1").String("X2")
		sqlStr, args, err := da.Interpolate().ToSQL()
		require.NoError(t, err)
		assert.Nil(t, args)
		assert.Exactly(t, "INSERT INTO `tableA` (`a`,`b`,`col1`,`col2`) VALUES (1,1,'X1','X2')", sqlStr)

		da.Options = 0

		// call it twice (4x) to test for being NOT idempotent
		compareToSQL(t, da, errors.NoKind,
			"INSERT INTO `tableA` (`a`,`b`,`col1`,`col2`) VALUES (?,?,?,?)",
			"INSERT INTO `tableA` (`a`,`b`,`col1`,`col2`) VALUES (1,1,'X1','X2')",
			int64(1), true, "X1", "X2",
		)

	})

	t.Run("Missing EventType", func(t *testing.T) {
		ins := NewInsert("tableA").AddColumns("a", "b")

		ins.Listeners.Add(
			Listen{
				Name: "colC",
				ListenInsertFn: func(i *Insert) {
					i.WithPairs(Column("colC").Str("X1"))
				},
			},
		)
		compareToSQL(t, ins, errors.Empty,
			"",
			"",
		)
	})

	t.Run("Should Dispatch", func(t *testing.T) {
		ins := NewInsert("tableA").AddColumns("a", "b")

		ins.Listeners.Add(
			Listen{
				EventType: OnBeforeToSQL,
				Name:      "colA",
				ListenInsertFn: func(i *Insert) {
					i.AddColumns("colA")
				},
			},
			Listen{
				EventType: OnBeforeToSQL,
				Name:      "colB",
				ListenInsertFn: func(i *Insert) {
					i.AddColumns("colB")
				},
			},
		)

		ins.Listeners.Add(
			Listen{
				// Multiple calls and colC is getting ignored because the WithPairs
				// function only creates the next value slice when a column `a`
				// gets called with WithPairs.
				EventType: OnBeforeToSQL,
				Name:      "colC",
				ListenInsertFn: func(i *Insert) {
					i.AddColumns("colC")
				},
			},
		)
		compareToSQL2(t, ins.WithArgs(), errors.NoKind,
			"INSERT INTO `tableA` (`a`,`b`,`colA`,`colB`,`colC`) VALUES (?,?,?,?,?)")
		compareToSQL2(t, ins, errors.NoKind,
			"INSERT INTO `tableA` (`a`,`b`,`colA`,`colB`,`colC`) VALUES ",
		)

		assert.Exactly(t, `colA; colB; colC`, ins.Listeners.String())
	})
}

func TestInsert_FromSelect(t *testing.T) {
	t.Parallel()

	t.Run("One Placeholder, ON DUPLICATE KEY", func(t *testing.T) {
		ins := NewInsert("tableA").AddColumns("a", "b").OnDuplicateKey()

		compareToSQL(t, ins.FromSelect(NewSelect("something_id", "user_id").
			From("some_table").
			Where(
				Column("d").PlaceHolder(),
				Column("e").Str("wat"),
			).
			OrderByDesc("id"),
		).
			WithArgs().Int(897),
			errors.NoKind,
			"INSERT INTO `tableA` (`a`,`b`) SELECT `something_id`, `user_id` FROM `some_table` WHERE (`d` = ?) AND (`e` = 'wat') ORDER BY `id` DESC ON DUPLICATE KEY UPDATE `a`=VALUES(`a`), `b`=VALUES(`b`)",
			"INSERT INTO `tableA` (`a`,`b`) SELECT `something_id`, `user_id` FROM `some_table` WHERE (`d` = 897) AND (`e` = 'wat') ORDER BY `id` DESC ON DUPLICATE KEY UPDATE `a`=VALUES(`a`), `b`=VALUES(`b`)",
			int64(897),
		)
		assert.Exactly(t, []string{"d"}, ins.qualifiedColumns)

	})

	t.Run("one PH, complex SELECT", func(t *testing.T) {
		ins := NewInsert("tableA").AddColumns("a", "b", "c")

		compareToSQL(t, ins.FromSelect(NewSelect("something_id", "user_id", "other").
			From("some_table").
			Where(
				ParenthesisOpen(),
				Column("d").PlaceHolder(),
				Column("e").Str("wat").Or(),
				ParenthesisClose(),
				Column("a").In().Int64s(1, 2, 3),
			).
			OrderByDesc("id").
			Paginate(1, 20),
		).
			WithArgs().Int(4444),
			errors.NoKind,
			"INSERT INTO `tableA` (`a`,`b`,`c`) SELECT `something_id`, `user_id`, `other` FROM `some_table` WHERE ((`d` = ?) OR (`e` = 'wat')) AND (`a` IN (1,2,3)) ORDER BY `id` DESC LIMIT 20 OFFSET 0",
			"INSERT INTO `tableA` (`a`,`b`,`c`) SELECT `something_id`, `user_id`, `other` FROM `some_table` WHERE ((`d` = 4444) OR (`e` = 'wat')) AND (`a` IN (1,2,3)) ORDER BY `id` DESC LIMIT 20 OFFSET 0",
			int64(4444),
		)
		assert.Exactly(t, []string{"d"}, ins.qualifiedColumns)
	})

	t.Run("Two placeholders", func(t *testing.T) {
		ins := NewInsert("tableA").AddColumns("a", "b")

		compareToSQL(t, ins.FromSelect(NewSelect("something_id", "user_id").
			From("some_table").
			Where(
				Column("d").PlaceHolder(),
				Column("a").In().Int64s(1, 2, 3),
				Column("e").PlaceHolder(),
			),
		).WithArgs().String("Guys!").Int(4444),
			errors.NoKind,
			"INSERT INTO `tableA` (`a`,`b`) SELECT `something_id`, `user_id` FROM `some_table` WHERE (`d` = ?) AND (`a` IN (1,2,3)) AND (`e` = ?)",
			"INSERT INTO `tableA` (`a`,`b`) SELECT `something_id`, `user_id` FROM `some_table` WHERE (`d` = 'Guys!') AND (`a` IN (1,2,3)) AND (`e` = 4444)",
			"Guys!", int64(4444),
		)
		assert.Exactly(t, []string{"d", "e"}, ins.qualifiedColumns)
	})

	t.Run("Record Simple,no select", func(t *testing.T) {
		p := &dmlPerson{
			Name:  "Pike",
			Email: MakeNullString("pikes@peak.co"),
		}

		ins := NewInsert("dml_people").AddColumns("name", "email").
			WithArgs().Record("", p)
		compareToSQL(t, ins, errors.NoKind,
			"INSERT INTO `dml_people` (`name`,`email`) VALUES (?,?)",
			"INSERT INTO `dml_people` (`name`,`email`) VALUES ('Pike','pikes@peak.co')",
			"Pike", "pikes@peak.co",
		)
	})

	t.Run("Record Complex", func(t *testing.T) {
		p := &dmlPerson{
			ID:    20180128,
			Name:  "Hans Wurst",
			Email: MakeNullString("hans@wurst.com"),
		}

		sel := NewSelect("a", "b").
			FromAlias("dml_person", "dp").
			Join(MakeIdentifier("dml_group").Alias("dg"), Column("dp.id").PlaceHolder()).
			Where(
				Column("dg.dob").Greater().PlaceHolder(),
				Column("age").Less().Int(56),
				Column("size").Greater().NamedArg("xSize"),
				ParenthesisOpen(),
				Column("dp.name").PlaceHolder(),
				Column("e").Str("wat").Or(),
				ParenthesisClose(),
				Column("fPlaceholder").LessOrEqual().PlaceHolder(),
				Column("g").Greater().Int(3),
				Column("h").In().Int64s(4, 5, 6),
			).
			GroupBy("ab").
			Having(
				Column("dp.email").PlaceHolder(),
				Column("n").Str("wh3r3"),
			).
			OrderBy("l")

		ins := NewInsert("tableA").AddColumns("a", "b").FromSelect(sel).WithArgs().Records(
			Qualify("dp", p),
			Qualify("dg", MakeArgs(1).Name("dob").Int(1970)),
		).Name("xSize").Int(678).Float64( /*fPlaceholder*/ 3.14159)

		compareToSQL(t, ins, errors.NoKind,
			"INSERT INTO `tableA` (`a`,`b`) SELECT `a`, `b` FROM `dml_person` AS `dp` INNER JOIN `dml_group` AS `dg` ON (`dp`.`id` = ?) WHERE (`dg`.`dob` > ?) AND (`age` < 56) AND (`size` > ?) AND ((`dp`.`name` = ?) OR (`e` = 'wat')) AND (`fPlaceholder` <= ?) AND (`g` > 3) AND (`h` IN (4,5,6)) GROUP BY `ab` HAVING (`dp`.`email` = ?) AND (`n` = 'wh3r3') ORDER BY `l`",
			"INSERT INTO `tableA` (`a`,`b`) SELECT `a`, `b` FROM `dml_person` AS `dp` INNER JOIN `dml_group` AS `dg` ON (`dp`.`id` = 20180128) WHERE (`dg`.`dob` > 1970) AND (`age` < 56) AND (`size` > 678) AND ((`dp`.`name` = 'Hans Wurst') OR (`e` = 'wat')) AND (`fPlaceholder` <= 3.14159) AND (`g` > 3) AND (`h` IN (4,5,6)) GROUP BY `ab` HAVING (`dp`.`email` = 'hans@wurst.com') AND (`n` = 'wh3r3') ORDER BY `l`",
			int64(20180128), int64(1970), int64(678), "Hans Wurst", 3.14159, "hans@wurst.com",
		)
	})
}

func TestInsert_Replace_Ignore(t *testing.T) {
	t.Parallel()

	// this generated statement does not comply the SQL standard
	compareToSQL(t, NewInsert("a").
		Replace().Ignore().
		AddColumns("b", "c").
		WithArgs().Int(1).Int(2).Int(3).Int(4),
		errors.NoKind,
		"REPLACE IGNORE INTO `a` (`b`,`c`) VALUES (?,?),(?,?)",
		"REPLACE IGNORE INTO `a` (`b`,`c`) VALUES (1,2),(3,4)",
		int64(1), int64(2), int64(3), int64(4),
	)
}

func TestInsert_WithoutColumns(t *testing.T) {
	t.Parallel()

	t.Run("each column in its own Arg", func(t *testing.T) {
		compareToSQL(t, NewInsert("catalog_product_link").SetRowCount(3).
			WithArgs().Int(2046).Int(33).Int(3).Int(2046).Int(34).Int(3).Int(2046).Int(35).Int(3),
			errors.NoKind,
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
			WithPairs(
				Column("product_id").Int64(2046),
				Column("linked_product_id").Int64(33),
				Column("link_type_id").Int64(3),
			).WithArgs(),
			errors.NoKind,
			"INSERT INTO `catalog_product_link` (`product_id`,`linked_product_id`,`link_type_id`) VALUES (?,?,?)",
			"INSERT INTO `catalog_product_link` (`product_id`,`linked_product_id`,`link_type_id`) VALUES (2046,33,3)",
			int64(2046), int64(33), int64(3),
		)
	})
	// TODO implement expression handling, requires some refactorings
	//t.Run("expression no args", func(t *testing.T) {
	//	compareToSQL(t, NewInsert("catalog_product_link").
	//		WithPairs(
	//			Column("product_id").Int64(2046),
	//			Column("type_name").Expression("CONCAT(`product_id`,'Manufacturer')"),
	//			Column("link_type_id").Int64(3),
	//		),
	//		errors.NoKind,
	//		"INSERT INTO `catalog_product_link` (`product_id`,`linked_product_id`,`link_type_id`) VALUES (?,CONCAT(`product_id`,'Manufacturer'),?)",
	//		"INSERT INTO `catalog_product_link` (`product_id`,`linked_product_id`,`link_type_id`) VALUES (2046,CONCAT(`product_id`,'Manufacturer'),3)",
	//		int64(2046), int64(33), int64(3),
	//	)
	//})
	t.Run("multiple rows triggers NO error", func(t *testing.T) {
		compareToSQL(t, NewInsert("catalog_product_link").
			WithPairs(
				// First row
				Column("product_id").Int64(2046),
				Column("linked_product_id").Int64(33),
				Column("link_type_id").Int64(3),

				// second row
				Column("product_id").Int64(2046),
				Column("linked_product_id").Int64(34),
				Column("link_type_id").Int64(3),
			).WithArgs(),
			errors.NoKind,
			"INSERT INTO `catalog_product_link` (`product_id`,`linked_product_id`,`link_type_id`) VALUES (?,?,?),(?,?,?)",
			"INSERT INTO `catalog_product_link` (`product_id`,`linked_product_id`,`link_type_id`) VALUES (2046,33,3),(2046,34,3)",
			int64(2046), int64(33), int64(3), int64(2046), int64(34), int64(3),
		)
	})
}

func TestInsert_DisableBuildCache(t *testing.T) {
	t.Parallel()

	ins := NewInsert("a").AddColumns("b", "c").
		OnDuplicateKey().DisableBuildCache().WithArgs()

	const cachedSQLPlaceHolder = "INSERT INTO `a` (`b`,`c`) VALUES (?,?),(?,?),(?,?) ON DUPLICATE KEY UPDATE `b`=VALUES(`b`), `c`=VALUES(`c`)"
	t.Run("without interpolate", func(t *testing.T) {
		for i := 0; i < 3; i++ {
			sql, args, err := ins.Raw(1, 2, 3, 4, 5, 6).ToSQL()
			require.NoError(t, err, "%+v", err)
			require.Equal(t, cachedSQLPlaceHolder, sql)
			require.Equal(t, []interface{}{1, 2, 3, 4, 5, 6}, args)
			require.Equal(t, "INSERT INTO `a` (`b`,`c`) VALUES  ON DUPLICATE KEY UPDATE `b`=VALUES(`b`), `c`=VALUES(`c`)",
				string(ins.base.cachedSQL))
			ins.Reset()
		}
	})

	t.Run("with interpolate", func(t *testing.T) {
		insA := ins.Reset().Interpolate()
		insA.Int(1).Int(2).Int(3).Int(4).Int(5).Int(6)

		const cachedSQLInterpolated = "INSERT INTO `a` (`b`,`c`) VALUES (1,2),(3,4),(5,6) ON DUPLICATE KEY UPDATE `b`=VALUES(`b`), `c`=VALUES(`c`)"
		for i := 0; i < 3; i++ {
			sql, args, err := insA.ToSQL()
			require.NoError(t, err, "%+v", err)
			require.Equal(t, cachedSQLInterpolated, sql)
			require.Nil(t, args)
			assert.Equal(t, "INSERT INTO `a` (`b`,`c`) VALUES  ON DUPLICATE KEY UPDATE `b`=VALUES(`b`), `c`=VALUES(`c`)", string(ins.base.cachedSQL))
		}
	})
}

func TestInsert_AddArguments(t *testing.T) {
	t.Parallel()

	t.Run("single WithArgs", func(t *testing.T) {
		compareToSQL(t,
			NewInsert("a").AddColumns("b", "c").WithArgs().Int(1).Int(2),
			errors.NoKind,
			"INSERT INTO `a` (`b`,`c`) VALUES (?,?)",
			"INSERT INTO `a` (`b`,`c`) VALUES (1,2)",
			int64(1), int64(2),
		)
	})
	t.Run("multi WithArgs on duplicate key", func(t *testing.T) {
		compareToSQL(t,
			NewInsert("a").AddColumns("b", "c").
				OnDuplicateKey().
				WithArgs().Int(1).Int(2).Int(3).Int(4).Int(5).Int(6),
			errors.NoKind,
			"INSERT INTO `a` (`b`,`c`) VALUES (?,?),(?,?),(?,?) ON DUPLICATE KEY UPDATE `b`=VALUES(`b`), `c`=VALUES(`c`)",
			"INSERT INTO `a` (`b`,`c`) VALUES (1,2),(3,4),(5,6) ON DUPLICATE KEY UPDATE `b`=VALUES(`b`), `c`=VALUES(`c`)",
			int64(1), int64(2), int64(3), int64(4), int64(5), int64(6),
		)
	})
	t.Run("single AddValues", func(t *testing.T) {
		compareToSQL(t,
			NewInsert("a").AddColumns("b", "c").WithArgs().Int64(1).Int64(2),
			errors.NoKind,
			"INSERT INTO `a` (`b`,`c`) VALUES (?,?)",
			"INSERT INTO `a` (`b`,`c`) VALUES (1,2)",
			int64(1), int64(2),
		)
	})
	t.Run("multi AddValues on duplicate key", func(t *testing.T) {
		compareToSQL(t,
			NewInsert("a").AddColumns("b", "c").
				OnDuplicateKey().WithArgs().
				Int64(1).Int64(2).
				Int64(3).Int64(4).
				Int64(5).Int64(6),
			errors.NoKind,
			"INSERT INTO `a` (`b`,`c`) VALUES (?,?),(?,?),(?,?) ON DUPLICATE KEY UPDATE `b`=VALUES(`b`), `c`=VALUES(`c`)",
			"INSERT INTO `a` (`b`,`c`) VALUES (1,2),(3,4),(5,6) ON DUPLICATE KEY UPDATE `b`=VALUES(`b`), `c`=VALUES(`c`)",
			int64(1), int64(2), int64(3), int64(4), int64(5), int64(6),
		)
	})
}

func TestInsert_OnDuplicateKey(t *testing.T) {
	t.Parallel()

	t.Run("Exclude only", func(t *testing.T) {
		compareToSQL(t, NewInsert("customer_gr1d_flat").
			AddColumns("entity_id", "name", "email", "group_id", "created_at", "website_id").
			AddOnDuplicateKeyExclude("entity_id").
			WithArgs().Int(1).String("Martin").String("martin@go.go").Int(3).String("2019-01-01").Int(2),
			errors.NoKind,
			"INSERT INTO `customer_gr1d_flat` (`entity_id`,`name`,`email`,`group_id`,`created_at`,`website_id`) VALUES (?,?,?,?,?,?) ON DUPLICATE KEY UPDATE `name`=VALUES(`name`), `email`=VALUES(`email`), `group_id`=VALUES(`group_id`), `created_at`=VALUES(`created_at`), `website_id`=VALUES(`website_id`)",
			"INSERT INTO `customer_gr1d_flat` (`entity_id`,`name`,`email`,`group_id`,`created_at`,`website_id`) VALUES (1,'Martin','martin@go.go',3,'2019-01-01',2) ON DUPLICATE KEY UPDATE `name`=VALUES(`name`), `email`=VALUES(`email`), `group_id`=VALUES(`group_id`), `created_at`=VALUES(`created_at`), `website_id`=VALUES(`website_id`)",
			int64(1), "Martin", "martin@go.go", int64(3), "2019-01-01", int64(2),
		)
	})

	t.Run("Exclude plus custom field value", func(t *testing.T) {
		compareToSQL(t, NewInsert("customer_gr1d_flat").
			AddColumns("entity_id", "name", "email", "group_id", "created_at", "website_id").
			AddOnDuplicateKeyExclude("entity_id").
			AddOnDuplicateKey(Column("created_at").Time(now())).
			WithArgs().Int(1).String("Martin").String("martin@go.go").Int(3).String("2019-01-01").Int(2),
			errors.NoKind,
			"INSERT INTO `customer_gr1d_flat` (`entity_id`,`name`,`email`,`group_id`,`created_at`,`website_id`) VALUES (?,?,?,?,?,?) ON DUPLICATE KEY UPDATE `name`=VALUES(`name`), `email`=VALUES(`email`), `group_id`=VALUES(`group_id`), `website_id`=VALUES(`website_id`), `created_at`='2006-01-02 15:04:05'",
			"INSERT INTO `customer_gr1d_flat` (`entity_id`,`name`,`email`,`group_id`,`created_at`,`website_id`) VALUES (1,'Martin','martin@go.go',3,'2019-01-01',2) ON DUPLICATE KEY UPDATE `name`=VALUES(`name`), `email`=VALUES(`email`), `group_id`=VALUES(`group_id`), `website_id`=VALUES(`website_id`), `created_at`='2006-01-02 15:04:05'",
			int64(1), "Martin", "martin@go.go", int64(3), "2019-01-01", int64(2),
		)
	})

	t.Run("Exclude plus default place holder, Arguments", func(t *testing.T) {
		ins := NewInsert("customer_gr1d_flat").
			AddColumns("entity_id", "name", "email", "group_id", "created_at", "website_id").
			AddOnDuplicateKeyExclude("entity_id").
			AddOnDuplicateKey(Column("created_at").PlaceHolder())
		insA := ins.WithArgs().Int(1).String("Martin").String("martin@go.go").Int(3).String("2019-01-01").Int(2).Name("time").Time(now())
		compareToSQL(t, insA, errors.NoKind,
			"INSERT INTO `customer_gr1d_flat` (`entity_id`,`name`,`email`,`group_id`,`created_at`,`website_id`) VALUES (?,?,?,?,?,?) ON DUPLICATE KEY UPDATE `name`=VALUES(`name`), `email`=VALUES(`email`), `group_id`=VALUES(`group_id`), `website_id`=VALUES(`website_id`), `created_at`=?",
			"INSERT INTO `customer_gr1d_flat` (`entity_id`,`name`,`email`,`group_id`,`created_at`,`website_id`) VALUES (1,'Martin','martin@go.go',3,'2019-01-01',2) ON DUPLICATE KEY UPDATE `name`=VALUES(`name`), `email`=VALUES(`email`), `group_id`=VALUES(`group_id`), `website_id`=VALUES(`website_id`), `created_at`='2006-01-02 15:04:05'",
			int64(1), "Martin", "martin@go.go", int64(3), "2019-01-01", int64(2), now(),
		)
		assert.Exactly(t, []string{"entity_id", "name", "email", "group_id", "created_at", "website_id", "created_at"}, ins.qualifiedColumns)
	})

	t.Run("Exclude plus default place holder, iface", func(t *testing.T) {
		ins := NewInsert("customer_gr1d_flat").
			AddColumns("entity_id", "name", "email", "group_id", "created_at", "website_id").
			AddOnDuplicateKeyExclude("entity_id").
			AddOnDuplicateKey(Column("created_at").PlaceHolder())
		insA := ins.WithArgs().Int(1).String("Martin").String("martin@go.go").Int(3).String("2019-01-01").Int(2).Time(now())
		compareToSQL2(t, insA, errors.NoKind,
			"INSERT INTO `customer_gr1d_flat` (`entity_id`,`name`,`email`,`group_id`,`created_at`,`website_id`) VALUES (?,?,?,?,?,?) ON DUPLICATE KEY UPDATE `name`=VALUES(`name`), `email`=VALUES(`email`), `group_id`=VALUES(`group_id`), `website_id`=VALUES(`website_id`), `created_at`=?",
			int64(1), "Martin", "martin@go.go", int64(3), "2019-01-01", int64(2), now(),
		)
		assert.Exactly(t, []string{"entity_id", "name", "email", "group_id", "created_at", "website_id", "created_at"}, insA.base.qualifiedColumns)
	})

	t.Run("Exclude plus custom place holder", func(t *testing.T) {
		ins := NewInsert("customer_gr1d_flat").
			AddColumns("entity_id", "name", "email", "group_id", "created_at", "website_id").
			AddOnDuplicateKeyExclude("entity_id").
			AddOnDuplicateKey(Column("created_at").NamedArg("time")).
			WithArgs().Int(1).String("Martin").String("martin@go.go").Int(3).String("2019-01-01").Int(2).Name("time").Time(now())
		compareToSQL(t, ins, errors.NoKind,
			"INSERT INTO `customer_gr1d_flat` (`entity_id`,`name`,`email`,`group_id`,`created_at`,`website_id`) VALUES (?,?,?,?,?,?) ON DUPLICATE KEY UPDATE `name`=VALUES(`name`), `email`=VALUES(`email`), `group_id`=VALUES(`group_id`), `website_id`=VALUES(`website_id`), `created_at`=?",
			"INSERT INTO `customer_gr1d_flat` (`entity_id`,`name`,`email`,`group_id`,`created_at`,`website_id`) VALUES (1,'Martin','martin@go.go',3,'2019-01-01',2) ON DUPLICATE KEY UPDATE `name`=VALUES(`name`), `email`=VALUES(`email`), `group_id`=VALUES(`group_id`), `website_id`=VALUES(`website_id`), `created_at`='2006-01-02 15:04:05'",
			int64(1), "Martin", "martin@go.go", int64(3), "2019-01-01", int64(2), now(),
		)
		assert.Exactly(t, []string{"entity_id", "name", "email", "group_id", "created_at", "website_id", ":time"}, ins.base.qualifiedColumns)
	})

	t.Run("Enabled for all columns", func(t *testing.T) {
		ins := NewInsert("customer_gr1d_flat").
			AddColumns("name", "email", "group_id", "created_at", "website_id").
			OnDuplicateKey().WithArgs().String("Martin").String("martin@go.go").Int(3).String("2019-01-01").Int(2)
		compareToSQL(t, ins, errors.NoKind,
			"INSERT INTO `customer_gr1d_flat` (`name`,`email`,`group_id`,`created_at`,`website_id`) VALUES (?,?,?,?,?) ON DUPLICATE KEY UPDATE `name`=VALUES(`name`), `email`=VALUES(`email`), `group_id`=VALUES(`group_id`), `created_at`=VALUES(`created_at`), `website_id`=VALUES(`website_id`)",
			"INSERT INTO `customer_gr1d_flat` (`name`,`email`,`group_id`,`created_at`,`website_id`) VALUES ('Martin','martin@go.go',3,'2019-01-01',2) ON DUPLICATE KEY UPDATE `name`=VALUES(`name`), `email`=VALUES(`email`), `group_id`=VALUES(`group_id`), `created_at`=VALUES(`created_at`), `website_id`=VALUES(`website_id`)",
			"Martin", "martin@go.go", int64(3), "2019-01-01", int64(2),
		)
		// testing for being idempotent
		compareToSQL(t, ins, errors.NoKind,
			"INSERT INTO `customer_gr1d_flat` (`name`,`email`,`group_id`,`created_at`,`website_id`) VALUES (?,?,?,?,?) ON DUPLICATE KEY UPDATE `name`=VALUES(`name`), `email`=VALUES(`email`), `group_id`=VALUES(`group_id`), `created_at`=VALUES(`created_at`), `website_id`=VALUES(`website_id`)",
			"", //"INSERT INTO `customer_gr1d_flat` (`name`,`email`,`group_id`,`created_at`,`website_id`) VALUES ('Martin','martin@go.go',3,'2019-01-01',2) ON DUPLICATE KEY UPDATE `name`=VALUES(`name`), `email`=VALUES(`email`), `group_id`=VALUES(`group_id`), `created_at`=VALUES(`created_at`), `website_id`=VALUES(`website_id`)",
			"Martin", "martin@go.go", int64(3), "2019-01-01", int64(2),
		)
	})
}

func TestInsert_Bind_Slice(t *testing.T) {
	t.Parallel()

	wantArgs := []interface{}{
		"Muffin Hat", "Muffin@Hat.head",
		"Marianne Phyllis Finch", "marianne@phyllis.finch",
		"Daphne Augusta Perry", "daphne@augusta.perry",
	}
	persons := &dmlPersons{
		Data: []*dmlPerson{
			{Name: "Muffin Hat", Email: MakeNullString("Muffin@Hat.head")},
			{Name: "Marianne Phyllis Finch", Email: MakeNullString("marianne@phyllis.finch")},
			{Name: "Daphne Augusta Perry", Email: MakeNullString("daphne@augusta.perry")},
		},
	}

	ins := NewInsert("dml_personXXX").AddColumns("name", "email").
		SetRowCount(len(persons.Data))

	const (
		wantPH = "INSERT INTO `dml_personXXX` (`name`,`email`) VALUES (?,?),(?,?),(?,?)"
		wantIP = "INSERT INTO `dml_personXXX` (`name`,`email`) VALUES ('Muffin Hat','Muffin@Hat.head'),('Marianne Phyllis Finch','marianne@phyllis.finch'),('Daphne Augusta Perry','daphne@augusta.perry')"
	)

	const iterations = 5
	var wg sync.WaitGroup
	wg.Add(iterations)
	for i := 0; i < iterations; i++ {
		go func(wg *sync.WaitGroup) {
			defer wg.Done()
			insA := ins.WithArgs().Record("", persons)

			compareToSQL(t, insA, errors.NoKind, wantPH, wantIP, wantArgs...)
			insA.Reset().Record("", persons)
			compareToSQL(t, insA, errors.NoKind, wantPH, wantIP, wantArgs...)
		}(&wg)
	}
	wg.Wait()

}
