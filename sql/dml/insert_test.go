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
	"database/sql"
	"testing"

	"github.com/corestoreio/errors"
	"github.com/corestoreio/pkg/storage/null"
	"github.com/corestoreio/pkg/sync/bgwork"
	"github.com/corestoreio/pkg/util/assert"
)

var _ ColumnMapper = (*someRecord)(nil)

type someRecord struct {
	SomethingID int
	UserID      int64
	Other       bool
}

func (sr someRecord) MapColumns(cm *ColumnMap) error {
	if cm.Mode() == ColumnMapEntityReadAll {
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

func TestInsert_NoArguments(t *testing.T) {
	t.Parallel()

	t.Run("ToSQL", func(t *testing.T) {
		ins := NewInsert("tableA").AddColumns("a", "b").BuildValues()
		compareToSQL2(t, ins, errors.NoKind,
			"INSERT INTO `tableA` (`a`,`b`) VALUES (?,?)")
	})
	t.Run("WithDBR.ToSQL", func(t *testing.T) {
		ins := NewInsert("tableA").AddColumns("a", "b").WithDBR()
		compareToSQL2(t, ins, errors.NoKind,
			"INSERT INTO `tableA` (`a`,`b`) VALUES (?,?)")
	})
}

func TestInsert_SetValuesCount(t *testing.T) {
	t.Parallel()

	t.Run("No BuildValues", func(t *testing.T) {
		ins := NewInsert("a").AddColumns("b", "c")
		compareToSQL2(t, ins, errors.NoKind,
			"INSERT INTO `a` (`b`,`c`) VALUES ",
		)
		assert.Exactly(t, []string{"b", "c"}, ins.qualifiedColumns)
	})
	t.Run("BuildValues", func(t *testing.T) {
		ins := NewInsert("a").AddColumns("b", "c").BuildValues()
		compareToSQL2(t, ins, errors.NoKind,
			"INSERT INTO `a` (`b`,`c`) VALUES (?,?)",
		)
		assert.Exactly(t, []string{"b", "c"}, ins.qualifiedColumns)
	})
	t.Run("set to two", func(t *testing.T) {
		compareToSQL2(t,
			NewInsert("a").AddColumns("b", "c").SetRowCount(2).BuildValues(),
			errors.NoKind,
			"INSERT INTO `a` (`b`,`c`) VALUES (?,?),(?,?)",
		)
	})
	t.Run("with values", func(t *testing.T) {
		ins := NewInsert("dml_people").AddColumns("name", "key")
		inA := ins.WithDBR()
		compareToSQL2(t, inA.TestWithArgs("Barack", "44"), errors.NoKind, "INSERT INTO `dml_people` (`name`,`key`) VALUES (?,?)",
			"Barack", "44",
		)
		assert.Exactly(t, []string{"name", "key"}, ins.qualifiedColumns)
	})
	t.Run("with record", func(t *testing.T) {
		person := dmlPerson{Name: "Barack"}
		person.Email.Valid = true
		person.Email.Data = "obama@whitehouse.gov"
		compareToSQL2(t,
			NewInsert("dml_people").AddColumns("name", "email").WithDBR().TestWithArgs(Qualify("", &person)),
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
	ins := s.InsertInto("dml_people").AddColumns("name", "key").WithDBR()

	compareExecContext(t, ins, []interface{}{"Barack", "44"}, 0, 1)
}

func TestInsertReal(t *testing.T) {
	// Insert by specifying values
	s := createRealSessionWithFixtures(t, nil)
	defer testCloser(t, s)
	ins := s.InsertInto("dml_people").AddColumns("name", "email").WithDBR()
	lastInsertID, _ := compareExecContext(t, ins, []interface{}{"Barack", "obama@whitehouse.gov"}, 3, 0)
	validateInsertingBarack(t, s, lastInsertID)

	// Insert by specifying a record (ptr to struct)

	person := dmlPerson{Name: "Barack"}
	person.Email.Valid = true
	person.Email.Data = "obama@whitehouse.gov"
	ins = s.InsertInto("dml_people").AddColumns("name", "email").WithDBR()
	lastInsertID, _ = compareExecContext(t, ins, []interface{}{Qualify("", &person)}, 4, 0)

	validateInsertingBarack(t, s, lastInsertID)
}

func validateInsertingBarack(t *testing.T, c *ConnPool, lastInsertID int64) {
	var person dmlPerson
	_, err := c.SelectFrom("dml_people").Star().Where(Column("id").Int64(lastInsertID)).WithDBR().Load(context.TODO(), &person)
	assert.NoError(t, err)

	assert.Exactly(t, lastInsertID, int64(person.ID))
	assert.Exactly(t, "Barack", person.Name)
	assert.Exactly(t, true, person.Email.Valid)
	assert.Exactly(t, "obama@whitehouse.gov", person.Email.Data)
}

func TestInsertReal_OnDuplicateKey(t *testing.T) {
	s := createRealSessionWithFixtures(t, nil)
	defer testCloser(t, s)

	p := &dmlPerson{
		Name:  "Pike",
		Email: null.MakeString("pikes@peak.co"),
	}

	res, err := s.InsertInto("dml_people").
		AddColumns("name", "email").
		WithDBR().ExecContext(context.TODO(), Qualify("", p))
	if err != nil {
		t.Fatalf("%+v", err)
	}
	assert.Exactly(t, uint64(3), p.ID, "Last Insert ID must be three")

	inID, err := res.LastInsertId()
	if err != nil {
		t.Fatalf("%+v", err)
	}
	{
		var p dmlPerson
		_, err = s.SelectFrom("dml_people").Star().Where(Column("id").Int64(inID)).WithDBR().Load(context.TODO(), &p)
		assert.NoError(t, err)
		assert.Exactly(t, "Pike", p.Name)
		assert.Exactly(t, "pikes@peak.co", p.Email.Data)
	}

	p.Name = "-"
	p.Email.Data = "pikes@peak.com"
	res, err = s.InsertInto("dml_people").
		AddColumns("id", "name", "email").
		AddOnDuplicateKey(Column("name").Str("Pik3"), Column("email").Values()).
		WithDBR().
		ExecContext(context.TODO(), Qualify("", p))
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
		_, err = s.SelectFrom("dml_people").Star().Where(Column("id").Int64(inID)).WithDBR().Load(context.TODO(), &p)
		assert.NoError(t, err)
		assert.Exactly(t, "Pik3", p.Name)
		assert.Exactly(t, "pikes@peak.com", p.Email.Data)
	}
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
			WithDBR().TestWithArgs(897),
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
			WithDBR().TestWithArgs(4444),
			errors.NoKind,
			"INSERT INTO `tableA` (`a`,`b`,`c`) SELECT `something_id`, `user_id`, `other` FROM `some_table` WHERE ((`d` = ?) OR (`e` = 'wat')) AND (`a` IN (1,2,3)) ORDER BY `id` DESC LIMIT 0,20",
			"INSERT INTO `tableA` (`a`,`b`,`c`) SELECT `something_id`, `user_id`, `other` FROM `some_table` WHERE ((`d` = 4444) OR (`e` = 'wat')) AND (`a` IN (1,2,3)) ORDER BY `id` DESC LIMIT 0,20",
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
		).WithDBR().TestWithArgs("Guys!", 4444),
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
			Email: null.MakeString("pikes@peak.co"),
		}

		ins := NewInsert("dml_people").AddColumns("name", "email").
			WithDBR().TestWithArgs(Qualify("", p))
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
			Email: null.MakeString("hans@wurst.com"),
		}
		p2 := &dmlPerson{
			Dob: 1970,
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

		ins := NewInsert("tableA").AddColumns("a", "b").FromSelect(sel).WithDBR()

		compareToSQL(t, ins.TestWithArgs(Qualify("dp", p), Qualify("dg", p2), sql.Named("xSize", 678) /*fPlaceholder*/, 3.14159), errors.NoKind,
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
		WithDBR().TestWithArgs(1, 2, 3, 4),
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
			WithDBR().TestWithArgs(2046, 33, 3, 2046, 34, 3, 2046, 35, 3),
			errors.NoKind,
			"INSERT INTO `catalog_product_link` VALUES (?,?,?),(?,?,?),(?,?,?)",
			"INSERT INTO `catalog_product_link` VALUES (2046,33,3),(2046,34,3),(2046,35,3)",
			int64(2046), int64(33), int64(3), int64(2046), int64(34), int64(3), int64(2046), int64(35), int64(3),
		)
	})
}

func TestInsert_sqlNamedArg(t *testing.T) {
	t.Parallel()

	t.Run("one row", func(t *testing.T) {
		compareToSQL(t, NewInsert("catalog_product_link").AddColumns("product_id", "linked_product_id", "link_type_id").
			WithDBR().TestWithArgs([]sql.NamedArg{
			{Name: "product_id", Value: 2046},
			{Name: "linked_product_id", Value: 33},
			{Name: "link_type_id", Value: 3},
		}),
			errors.NoKind,
			"INSERT INTO `catalog_product_link` (`product_id`,`linked_product_id`,`link_type_id`) VALUES (?,?,?)",
			"INSERT INTO `catalog_product_link` (`product_id`,`linked_product_id`,`link_type_id`) VALUES (2046,33,3)",
			int64(2046), int64(33), int64(3),
		)
	})
	// TODO implement expression handling, requires some refactorings
	// t.Run("expression no args", func(t *testing.T) {
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
		compareToSQL(t, NewInsert("catalog_product_link").AddColumns("product_id", "linked_product_id", "link_type_id").
			WithDBR().TestWithArgs([]sql.NamedArg{
			{Name: "product_id", Value: 2046},
			{Name: "linked_product_id", Value: 33},
			{Name: "link_type_id", Value: 3},

			{Name: "product_id", Value: 2046},
			{Name: "linked_product_id", Value: 34},
			{Name: "link_type_id", Value: 3},
		}),
			errors.NoKind,
			"INSERT INTO `catalog_product_link` (`product_id`,`linked_product_id`,`link_type_id`) VALUES (?,?,?),(?,?,?)",
			"INSERT INTO `catalog_product_link` (`product_id`,`linked_product_id`,`link_type_id`) VALUES (2046,33,3),(2046,34,3)",
			int64(2046), int64(33), int64(3), int64(2046), int64(34), int64(3),
		)
	})
}

func TestInsert_DisableBuildCache(t *testing.T) {
	t.Parallel()

	insA := NewInsert("a").AddColumns("b", "c").
		OnDuplicateKey().WithDBR()

	const cachedSQLPlaceHolder = "INSERT INTO `a` (`b`,`c`) VALUES (?,?),(?,?),(?,?) ON DUPLICATE KEY UPDATE `b`=VALUES(`b`), `c`=VALUES(`c`)"
	t.Run("without interpolate", func(t *testing.T) {
		for i := 0; i < 3; i++ {
			sql, args, err := insA.TestWithArgs(1, 2, 3, 4, 5, 6).ToSQL()
			assert.NoError(t, err)
			assert.Exactly(t, cachedSQLPlaceHolder, sql)
			assert.Exactly(t, []interface{}{int64(1), int64(2), int64(3), int64(4), int64(5), int64(6)}, args)
			insA.Reset()
		}
		assert.Exactly(t, []string{"", "INSERT INTO `a` (`b`,`c`) VALUES  ON DUPLICATE KEY UPDATE `b`=VALUES(`b`), `c`=VALUES(`c`)"},
			insA.CachedQueries())
	})

	t.Run("with interpolate", func(t *testing.T) {
		insA := insA.Reset().Interpolate()
		qb := insA.testWithArgs(1, 2, 3, 4, 5, 6)

		const cachedSQLInterpolated = "INSERT INTO `a` (`b`,`c`) VALUES (1,2),(3,4),(5,6) ON DUPLICATE KEY UPDATE `b`=VALUES(`b`), `c`=VALUES(`c`)"
		for i := 0; i < 3; i++ {
			sql, args, err := qb.ToSQL()
			assert.NoError(t, err)
			assert.Exactly(t, cachedSQLInterpolated, sql)
			assert.Nil(t, args)
		}
		assert.Exactly(t, []string{"", "INSERT INTO `a` (`b`,`c`) VALUES  ON DUPLICATE KEY UPDATE `b`=VALUES(`b`), `c`=VALUES(`c`)"},
			insA.CachedQueries())
	})
}

func TestInsert_AddArguments(t *testing.T) {
	t.Parallel()

	t.Run("single WithDBR", func(t *testing.T) {
		compareToSQL(t,
			NewInsert("a").AddColumns("b", "c").WithDBR().TestWithArgs(1, 2),
			errors.NoKind,
			"INSERT INTO `a` (`b`,`c`) VALUES (?,?)",
			"INSERT INTO `a` (`b`,`c`) VALUES (1,2)",
			int64(1), int64(2),
		)
	})
	t.Run("multi WithDBR on duplicate key", func(t *testing.T) {
		compareToSQL(t,
			NewInsert("a").AddColumns("b", "c").
				OnDuplicateKey().
				WithDBR().TestWithArgs(1, 2, 3, 4, 5, 6),
			errors.NoKind,
			"INSERT INTO `a` (`b`,`c`) VALUES (?,?),(?,?),(?,?) ON DUPLICATE KEY UPDATE `b`=VALUES(`b`), `c`=VALUES(`c`)",
			"INSERT INTO `a` (`b`,`c`) VALUES (1,2),(3,4),(5,6) ON DUPLICATE KEY UPDATE `b`=VALUES(`b`), `c`=VALUES(`c`)",
			int64(1), int64(2), int64(3), int64(4), int64(5), int64(6),
		)
	})
	t.Run("single AddValues", func(t *testing.T) {
		compareToSQL(t,
			NewInsert("a").AddColumns("b", "c").WithDBR().TestWithArgs(1, 2),
			errors.NoKind,
			"INSERT INTO `a` (`b`,`c`) VALUES (?,?)",
			"INSERT INTO `a` (`b`,`c`) VALUES (1,2)",
			int64(1), int64(2),
		)
	})
	t.Run("multi AddValues on duplicate key", func(t *testing.T) {
		compareToSQL(t,
			NewInsert("a").AddColumns("b", "c").
				OnDuplicateKey().WithDBR().TestWithArgs(1, 2, 3, 4, 5, 6),
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
			WithDBR().TestWithArgs(1, "Martin", "martin@go.go", 3, "2019-01-01", 2),
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
			WithDBR().TestWithArgs(1, "Martin", "martin@go.go", 3, "2019-01-01", 2),
			errors.NoKind,
			"INSERT INTO `customer_gr1d_flat` (`entity_id`,`name`,`email`,`group_id`,`created_at`,`website_id`) VALUES (?,?,?,?,?,?) ON DUPLICATE KEY UPDATE `name`=VALUES(`name`), `email`=VALUES(`email`), `group_id`=VALUES(`group_id`), `website_id`=VALUES(`website_id`), `created_at`='2006-01-02 15:04:05'",
			"INSERT INTO `customer_gr1d_flat` (`entity_id`,`name`,`email`,`group_id`,`created_at`,`website_id`) VALUES (1,'Martin','martin@go.go',3,'2019-01-01',2) ON DUPLICATE KEY UPDATE `name`=VALUES(`name`), `email`=VALUES(`email`), `group_id`=VALUES(`group_id`), `website_id`=VALUES(`website_id`), `created_at`='2006-01-02 15:04:05'",
			int64(1), "Martin", "martin@go.go", int64(3), "2019-01-01", int64(2),
		)
	})

	t.Run("Exclude plus default place holder, DBR", func(t *testing.T) {
		ins := NewInsert("customer_gr1d_flat").
			AddColumns("entity_id", "name", "email", "group_id", "created_at", "website_id").
			AddOnDuplicateKeyExclude("entity_id").
			AddOnDuplicateKey(Column("created_at").PlaceHolder())
		insA := ins.WithDBR().TestWithArgs(
			1, "Martin", "martin@go.go", 3, "2019-01-01", 2, sql.Named("time", now()),
		)
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
		insA := ins.WithDBR()
		compareToSQL2(t, insA.TestWithArgs(1, "Martin", "martin@go.go", 3, "2019-01-01", 2, now()), errors.NoKind,
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
			WithDBR()
		compareToSQL(t, ins.TestWithArgs(1, "Martin", "martin@go.go", 3, "2019-01-01", 2, sql.Named("time", now())), errors.NoKind,
			"INSERT INTO `customer_gr1d_flat` (`entity_id`,`name`,`email`,`group_id`,`created_at`,`website_id`) VALUES (?,?,?,?,?,?) ON DUPLICATE KEY UPDATE `name`=VALUES(`name`), `email`=VALUES(`email`), `group_id`=VALUES(`group_id`), `website_id`=VALUES(`website_id`), `created_at`=?",
			"INSERT INTO `customer_gr1d_flat` (`entity_id`,`name`,`email`,`group_id`,`created_at`,`website_id`) VALUES (1,'Martin','martin@go.go',3,'2019-01-01',2) ON DUPLICATE KEY UPDATE `name`=VALUES(`name`), `email`=VALUES(`email`), `group_id`=VALUES(`group_id`), `website_id`=VALUES(`website_id`), `created_at`='2006-01-02 15:04:05'",
			int64(1), "Martin", "martin@go.go", int64(3), "2019-01-01", int64(2), now(),
		)
		assert.Exactly(t, []string{"entity_id", "name", "email", "group_id", "created_at", "website_id", ":time"}, ins.base.qualifiedColumns)
	})

	t.Run("Enabled for all columns", func(t *testing.T) {
		ins := NewInsert("customer_gr1d_flat").
			AddColumns("name", "email", "group_id", "created_at", "website_id").
			OnDuplicateKey().WithDBR().TestWithArgs("Martin", "martin@go.go", 3, "2019-01-01", 2)
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

// TestInsert_Parallel_Bind_Slice is a tough test because first a complex SQL
// statement from a collection and second it runs in parallel.
func TestInsert_Parallel_Bind_Slice(t *testing.T) {
	t.Parallel()

	wantArgs := []interface{}{
		"Muffin Hat", "Muffin@Hat.head",
		"Marianne Phyllis Finch", "marianne@phyllis.finch",
		"Daphne Augusta Perry", "daphne@augusta.perry",
	}
	persons := &dmlPersons{
		Data: []*dmlPerson{
			{Name: "Muffin Hat", Email: null.MakeString("Muffin@Hat.head")},
			{Name: "Marianne Phyllis Finch", Email: null.MakeString("marianne@phyllis.finch")},
			{Name: "Daphne Augusta Perry", Email: null.MakeString("daphne@augusta.perry")},
		},
	}

	ins := NewInsert("dml_personXXX").AddColumns("name", "email").
		SetRowCount(len(persons.Data))

	const (
		wantPH = "INSERT INTO `dml_personXXX` (`name`,`email`) VALUES (?,?),(?,?),(?,?)"
		wantIP = "INSERT INTO `dml_personXXX` (`name`,`email`) VALUES ('Muffin Hat','Muffin@Hat.head'),('Marianne Phyllis Finch','marianne@phyllis.finch'),('Daphne Augusta Perry','daphne@augusta.perry')"
	)

	const concurrencyLevel = 10
	bgwork.Wait(concurrencyLevel, func(index int) {
		// Don't use such a construct in production code!
		// TODO try to move this above to see if it's race free.
		insA := ins.WithDBR()

		compareToSQL(t, insA.TestWithArgs(Qualify("", persons)), errors.NoKind, wantPH, wantIP, wantArgs...)

		compareToSQL(t, insA.TestWithArgs(Qualify("", persons)), errors.NoKind, wantPH, wantIP, wantArgs...)
	})
}

func TestInsert_Expressions_In_Values(t *testing.T) {
	t.Parallel()

	t.Run("1 string expression one row", func(t *testing.T) {
		ins := NewInsert("catalog_product_customer_relation").
			AddColumns("product_id", "sort_order").
			WithPairs(
				Column("customer_id").Expr("IFNULL(SELECT entity_id FROM customer_entity WHERE email like ?,0)"),
			).BuildValues()

		compareToSQL2(t, ins, errors.NoKind,
			"INSERT INTO `catalog_product_customer_relation` (`product_id`,`sort_order`,`customer_id`) VALUES (?,?,IFNULL(SELECT entity_id FROM customer_entity WHERE email like ?,0))",
		)
	})
	// TODO Not yet supported. some calculations necessary in Insert.toSQL
	// t.Run("2 string expression multiple rows", func(t *testing.T) {
	//	ins := NewInsert("catalog_product_customer_relation").
	//		AddColumns("product_id", "sort_order").
	//		WithPairs(
	//			Column("customer_id").Expr("IFNULL(SELECT entity_id FROM customer_entity WHERE email like ?,0)"),
	//			Column("customer_id").Expr("IFNULL(SELECT entity_id FROM customer_entity WHERE email like ?,0)"),
	//		).BuildValues()
	//
	//	compareToSQL2(t, ins, errors.NoKind,
	//		"INSERT INTO `catalog_product_customer_relation` (`product_id`,`sort_order`,`customer_id`) VALUES (?,?,IFNULL(SELECT entity_id FROM customer_entity WHERE email like ?,0)),(?,?,IFNULL(SELECT entity_id FROM customer_entity WHERE email like ?,0))",
	//	)
	//})

	t.Run("sub select", func(t *testing.T) {
		// do not use such a construct like the test query. use such a construct:
		/*
			INSERT INTO catalog_product_customer_relation (product_id, sort_order, group_id)
			  SELECT
					? AS product_id,
					? AS sort_order,
					group_id
					FROM customer_group
					WHERE name = ?;
		*/

		ins := NewInsert("catalog_product_customer_relation").
			AddColumns("product_id", "sort_order").
			WithPairs(
				Column("group_id").Sub(
					NewSelect("group_id").From("customer_group").Where(
						Column("name").Equal().PlaceHolder(),
					),
				),
			).BuildValues()
		compareToSQL(t, ins, errors.NoKind,
			"INSERT INTO `catalog_product_customer_relation` (`product_id`,`sort_order`,`group_id`) VALUES (?,?,(SELECT `group_id` FROM `customer_group` WHERE (`name` = ?)))",
			"",
		)
	})

	t.Run("all possibilities", func(t *testing.T) {
		ins := NewInsert("catalog_product_customer_relation").
			AddColumns("product_id", "sort_order").
			WithPairs(
				Column("customer_id").Expr("IFNULL(SELECT entity_id FROM customer_entity WHERE email like ?,0)"),
				Column("group_id").Sub(
					NewSelect("group_id").From("customer_group").Where(
						Column("name").Equal().PlaceHolder(),
					),
				),
			).BuildValues()
		compareToSQL(t, ins, errors.NoKind,
			"INSERT INTO `catalog_product_customer_relation` (`product_id`,`sort_order`,`customer_id`,`group_id`) VALUES (?,?,IFNULL(SELECT entity_id FROM customer_entity WHERE email like ?,0),(SELECT `group_id` FROM `customer_group` WHERE (`name` = ?)))",
			"",
		)
	})
}
