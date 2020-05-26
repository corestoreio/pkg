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
	"database/sql"
	"testing"

	"github.com/corestoreio/pkg/storage/null"

	"github.com/corestoreio/errors"
)

func TestWith_Placeholder(t *testing.T) {
	t.Run("placeholder DBR", func(t *testing.T) {
		cte := NewWith(
			WithCTE{
				Name:    "cte",
				Columns: []string{"n"},
				Union: NewUnion(
					NewSelect("a").AddColumnsAliases("d", "b").From("tableAD").Where(Column("b").PlaceHolder()),
					NewSelect("a", "b").From("tableAB").Where(Column("b").Like().NamedArg("nArg2")),
				).All(),
			},
		).
			Recursive().
			Select(NewSelect().Star().From("cte").Where(Column("a").GreaterOrEqual().PlaceHolder()))

		compareToSQL(t,
			cte.WithDBR(dbMock{}).TestWithArgs(sql.Named("nArg2", "hello%"), null.MakeString("arg1"), 2.7182),
			errors.NoKind,
			"WITH RECURSIVE `cte` (`n`) AS ((SELECT `a`, `d` AS `b` FROM `tableAD` WHERE (`b` = ?))\nUNION ALL\n(SELECT `a`, `b` FROM `tableAB` WHERE (`b` LIKE ?)))\nSELECT * FROM `cte` WHERE (`a` >= ?)",
			"WITH RECURSIVE `cte` (`n`) AS ((SELECT `a`, `d` AS `b` FROM `tableAD` WHERE (`b` = 'arg1'))\nUNION ALL\n(SELECT `a`, `b` FROM `tableAB` WHERE (`b` LIKE 'hello%')))\nSELECT * FROM `cte` WHERE (`a` >= 2.7182)",
			"arg1", "hello%", 2.7182,
		)
	})
}

func TestWith_ToSQL(t *testing.T) {
	t.Run("Find best and worst month With cache", func(t *testing.T) {
		/*
			WITH sales_by_month(month, total)
			     AS
			     -- first CTE: one row per month, with amount sold on all days of month
			     (SELECT Month(day_of_sale),Sum(amount)
			      FROM   sales_days
			      WHERE  Year(day_of_sale) = 2015
			      GROUP  BY Month(day_of_sale)),

			     best_month(month, total, award)
			     AS -- second CTE: best month
			     (SELECT month,
			             total,
			             "best"
			      FROM   sales_by_month
			      WHERE  total = (SELECT Max(total) FROM   sales_by_month)),
			     worst_month(month, total, award)
			     AS -- 3rd CTE: worst month
			     (SELECT month,
			             total,
			             "worst"
			      FROM   sales_by_month
			      WHERE  total = (SELECT Min(total)
			                      FROM   sales_by_month))
			-- Now show best and worst:
			SELECT *
			FROM   best_month
			UNION ALL
			SELECT *
			FROM   worst_month;
		*/
		cte := NewWith(
			WithCTE{
				Name: "sales_by_month", Columns: []string{"month", "total"},
				Select: NewSelect().Unsafe().AddColumns("Month(day_of_sale)", "Sum(amount)").
					From("sales_days").
					Where(Expr("Year(day_of_sale) = ?").Int(2015)).
					GroupBy("Month(day_of_sale))"),
			},
			WithCTE{
				Name: "best_month", Columns: []string{"month", "total", "award"},
				Select: NewSelect().AddColumns("month", "total").AddColumnsConditions(Expr(`"best"`)).From("sales_by_month").
					Where(
						Column("total").Equal().Sub(NewSelect().AddColumnsConditions(Expr("Max(total)")).From("sales_by_month"))),
			},
			WithCTE{
				Name: "worst_month", Columns: []string{"month", "total", "award"},
				Select: NewSelect().AddColumns("month", "total").AddColumnsConditions(Expr(`"worst"`)).From("sales_by_month").
					Where(
						Column("total").Equal().Sub(NewSelect().AddColumnsConditions(Expr("Min(total)")).From("sales_by_month"))),
			},
		).Union(NewUnion(
			NewSelect().Star().From("best_month"),
			NewSelect().Star().From("worst_month"),
		).All())

		compareToSQL(t, cte, errors.NoKind,
			"WITH `sales_by_month` (`month`,`total`) AS (SELECT Month(day_of_sale), Sum(amount) FROM `sales_days` WHERE (Year(day_of_sale) = 2015) GROUP BY Month(day_of_sale))),\n`best_month` (`month`,`total`,`award`) AS (SELECT `month`, `total`, \"best\" FROM `sales_by_month` WHERE (`total` = (SELECT Max(total) FROM `sales_by_month`))),\n`worst_month` (`month`,`total`,`award`) AS (SELECT `month`, `total`, \"worst\" FROM `sales_by_month` WHERE (`total` = (SELECT Min(total) FROM `sales_by_month`)))\n(SELECT * FROM `best_month`)\nUNION ALL\n(SELECT * FROM `worst_month`)",
			"WITH `sales_by_month` (`month`,`total`) AS (SELECT Month(day_of_sale), Sum(amount) FROM `sales_days` WHERE (Year(day_of_sale) = 2015) GROUP BY Month(day_of_sale))),\n`best_month` (`month`,`total`,`award`) AS (SELECT `month`, `total`, \"best\" FROM `sales_by_month` WHERE (`total` = (SELECT Max(total) FROM `sales_by_month`))),\n`worst_month` (`month`,`total`,`award`) AS (SELECT `month`, `total`, \"worst\" FROM `sales_by_month` WHERE (`total` = (SELECT Min(total) FROM `sales_by_month`)))\n(SELECT * FROM `best_month`)\nUNION ALL\n(SELECT * FROM `worst_month`)",
		)
		// call it twice
		compareToSQL(t, cte, errors.NoKind,
			"WITH `sales_by_month` (`month`,`total`) AS (SELECT Month(day_of_sale), Sum(amount) FROM `sales_days` WHERE (Year(day_of_sale) = 2015) GROUP BY Month(day_of_sale))),\n`best_month` (`month`,`total`,`award`) AS (SELECT `month`, `total`, \"best\" FROM `sales_by_month` WHERE (`total` = (SELECT Max(total) FROM `sales_by_month`))),\n`worst_month` (`month`,`total`,`award`) AS (SELECT `month`, `total`, \"worst\" FROM `sales_by_month` WHERE (`total` = (SELECT Min(total) FROM `sales_by_month`)))\n(SELECT * FROM `best_month`)\nUNION ALL\n(SELECT * FROM `worst_month`)",
			"WITH `sales_by_month` (`month`,`total`) AS (SELECT Month(day_of_sale), Sum(amount) FROM `sales_days` WHERE (Year(day_of_sale) = 2015) GROUP BY Month(day_of_sale))),\n`best_month` (`month`,`total`,`award`) AS (SELECT `month`, `total`, \"best\" FROM `sales_by_month` WHERE (`total` = (SELECT Max(total) FROM `sales_by_month`))),\n`worst_month` (`month`,`total`,`award`) AS (SELECT `month`, `total`, \"worst\" FROM `sales_by_month` WHERE (`total` = (SELECT Min(total) FROM `sales_by_month`)))\n(SELECT * FROM `best_month`)\nUNION ALL\n(SELECT * FROM `worst_month`)",
		)
	})
}
