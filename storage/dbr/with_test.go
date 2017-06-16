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
	"testing"

	"github.com/corestoreio/errors"
	"github.com/stretchr/testify/assert"
)

func TestNewWith(t *testing.T) {
	t.Parallel()

	t.Run("one CTE", func(t *testing.T) {
		cte := NewWith(
			WithCTE{Name: "one", Select: NewSelect().AddColumnsExpr("1")},
		).Select(NewSelect().Star().From("one"))
		compareToSQL(t, cte, nil,
			"WITH\n`one` AS (SELECT 1)\nSELECT * FROM `one`",
			"WITH\n`one` AS (SELECT 1)\nSELECT * FROM `one`",
		)
	})
	t.Run("one CTE recursive", func(t *testing.T) {
		cte := NewWith(
			WithCTE{
				Name:    "cte",
				Columns: []string{"n"},
				Union: NewUnion(
					NewSelect().AddColumnsExpr("1"),
					NewSelect().AddColumnsExpr("n+1").From("cte").Where(Column("n", Less.Int(5))),
				).All(),
			},
		).Recursive().Select(NewSelect().Star().From("cte"))
		compareToSQL(t, cte, nil,
			"WITH RECURSIVE\n`cte` (`n`) AS ((SELECT 1)\nUNION ALL\n(SELECT n+1 FROM `cte` WHERE (`n` < ?)))\nSELECT * FROM `cte`",
			"WITH RECURSIVE\n`cte` (`n`) AS ((SELECT 1)\nUNION ALL\n(SELECT n+1 FROM `cte` WHERE (`n` < 5)))\nSELECT * FROM `cte`",
			int64(5),
		)
	})

	t.Run("two CTEs", func(t *testing.T) {
		cte := NewWith(
			WithCTE{Name: "intermed", Select: NewSelect().Star().From("test").Where(Column("x", GreaterOrEqual.Int(5)))},
			WithCTE{Name: "derived", Select: NewSelect().Star().From("intermed").Where(Column("x", Less.Int(10)))},
		).Select(NewSelect().Star().From("derived"))
		compareToSQL(t, cte, nil,
			"WITH\n`intermed` AS (SELECT * FROM `test` WHERE (`x` >= ?)),\n`derived` AS (SELECT * FROM `intermed` WHERE (`x` < ?))\nSELECT * FROM `derived`",
			"WITH\n`intermed` AS (SELECT * FROM `test` WHERE (`x` >= 5)),\n`derived` AS (SELECT * FROM `intermed` WHERE (`x` < 10))\nSELECT * FROM `derived`",
			int64(5), int64(10),
		)
	})
	t.Run("multi column", func(t *testing.T) {
		cte := NewWith(
			WithCTE{Name: "multi", Columns: []string{"x", "y"}, Select: NewSelect().AddColumnsExpr("1", "2")},
		).Select(NewSelect("x", "y").From("multi"))
		compareToSQL(t, cte, nil,
			"WITH\n`multi` (`x`,`y`) AS (SELECT 1, 2)\nSELECT `x`, `y` FROM `multi`",
			"",
		)
	})

	t.Run("Find best and worst month", func(t *testing.T) {
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
			WithCTE{Name: "sales_by_month", Columns: []string{"month", "total"},
				Select: NewSelect().AddColumnsExpr("Month(day_of_sale)", "Sum(amount)").From("sales_days").
					Where(Expression("Year(day_of_sale) = ?", ArgInt(2015))).
					GroupByExpr("Month(day_of_sale))"),
			},
			WithCTE{Name: "best_month", Columns: []string{"month", "total", "award"},
				Select: NewSelect().AddColumns("month", "total").AddColumnsExpr(`"best"`).From("sales_by_month").
					Where(SubSelect("total", Equal, NewSelect().AddColumnsExpr("Max(total)").From("sales_by_month"))),
			},
			WithCTE{Name: "worst_month", Columns: []string{"month", "total", "award"},
				Select: NewSelect().AddColumns("month", "total").AddColumnsExpr(`"worst"`).From("sales_by_month").
					Where(SubSelect("total", Equal, NewSelect().AddColumnsExpr("Min(total)").From("sales_by_month"))),
			},
		).Union(NewUnion(
			NewSelect().Star().From("best_month"),
			NewSelect().Star().From("worst_month"),
		).All())
		cte.UseBuildCache = true

		compareToSQL(t, cte, nil,
			"WITH\n`sales_by_month` (`month`,`total`) AS (SELECT Month(day_of_sale), Sum(amount) FROM `sales_days` WHERE (Year(day_of_sale) = ?) GROUP BY Month(day_of_sale))),\n`best_month` (`month`,`total`,`award`) AS (SELECT `month`, `total`, \"best\" FROM `sales_by_month` WHERE (`total` = (SELECT Max(total) FROM `sales_by_month`))),\n`worst_month` (`month`,`total`,`award`) AS (SELECT `month`, `total`, \"worst\" FROM `sales_by_month` WHERE (`total` = (SELECT Min(total) FROM `sales_by_month`)))\n(SELECT * FROM `best_month`)\nUNION ALL\n(SELECT * FROM `worst_month`)",
			"WITH\n`sales_by_month` (`month`,`total`) AS (SELECT Month(day_of_sale), Sum(amount) FROM `sales_days` WHERE (Year(day_of_sale) = 2015) GROUP BY Month(day_of_sale))),\n`best_month` (`month`,`total`,`award`) AS (SELECT `month`, `total`, 'best' FROM `sales_by_month` WHERE (`total` = (SELECT Max(total) FROM `sales_by_month`))),\n`worst_month` (`month`,`total`,`award`) AS (SELECT `month`, `total`, 'worst' FROM `sales_by_month` WHERE (`total` = (SELECT Min(total) FROM `sales_by_month`)))\n(SELECT * FROM `best_month`)\nUNION ALL\n(SELECT * FROM `worst_month`)",
			int64(2015),
		)
		// call it twice
		compareToSQL(t, cte, nil,
			"WITH\n`sales_by_month` (`month`,`total`) AS (SELECT Month(day_of_sale), Sum(amount) FROM `sales_days` WHERE (Year(day_of_sale) = ?) GROUP BY Month(day_of_sale))),\n`best_month` (`month`,`total`,`award`) AS (SELECT `month`, `total`, \"best\" FROM `sales_by_month` WHERE (`total` = (SELECT Max(total) FROM `sales_by_month`))),\n`worst_month` (`month`,`total`,`award`) AS (SELECT `month`, `total`, \"worst\" FROM `sales_by_month` WHERE (`total` = (SELECT Min(total) FROM `sales_by_month`)))\n(SELECT * FROM `best_month`)\nUNION ALL\n(SELECT * FROM `worst_month`)",
			"WITH\n`sales_by_month` (`month`,`total`) AS (SELECT Month(day_of_sale), Sum(amount) FROM `sales_days` WHERE (Year(day_of_sale) = 2015) GROUP BY Month(day_of_sale))),\n`best_month` (`month`,`total`,`award`) AS (SELECT `month`, `total`, 'best' FROM `sales_by_month` WHERE (`total` = (SELECT Max(total) FROM `sales_by_month`))),\n`worst_month` (`month`,`total`,`award`) AS (SELECT `month`, `total`, 'worst' FROM `sales_by_month` WHERE (`total` = (SELECT Min(total) FROM `sales_by_month`)))\n(SELECT * FROM `best_month`)\nUNION ALL\n(SELECT * FROM `worst_month`)",
			int64(2015),
		)
		assert.Equal(t, "WITH\n`sales_by_month` (`month`,`total`) AS (SELECT Month(day_of_sale), Sum(amount) FROM `sales_days` WHERE (Year(day_of_sale) = ?) GROUP BY Month(day_of_sale))),\n`best_month` (`month`,`total`,`award`) AS (SELECT `month`, `total`, \"best\" FROM `sales_by_month` WHERE (`total` = (SELECT Max(total) FROM `sales_by_month`))),\n`worst_month` (`month`,`total`,`award`) AS (SELECT `month`, `total`, \"worst\" FROM `sales_by_month` WHERE (`total` = (SELECT Min(total) FROM `sales_by_month`)))\n(SELECT * FROM `best_month`)\nUNION ALL\n(SELECT * FROM `worst_month`)",
			string(cte.cacheSQL))
	})

	t.Run("DELETE", func(t *testing.T) {
		cte := NewWith(
			WithCTE{Name: "check_vals", Columns: []string{"val"}, Select: NewSelect().AddColumnsExpr("123")},
		).Delete(NewDelete("test").Where(SubSelect("val", In, NewSelect("val").From("check_vals"))))

		compareToSQL(t, cte, nil,
			"WITH\n`check_vals` (`val`) AS (SELECT 123)\nDELETE FROM `test` WHERE (`val` IN (SELECT `val` FROM `check_vals`))",
			"WITH\n`check_vals` (`val`) AS (SELECT 123)\nDELETE FROM `test` WHERE (`val` IN (SELECT `val` FROM `check_vals`))",
		)
	})
	t.Run("UPDATE", func(t *testing.T) {
		cte := NewWith(
			WithCTE{Name: "my_cte", Columns: []string{"n"}, Union: NewUnion(
				NewSelect().AddColumnsExpr("1"),
				NewSelect().AddColumnsExpr("1+n").From("my_cte").Where(Column("n", Less.Int(6))),
			).All()},
			// UPDATE statement is wrong because we're missing a JOIN which is not yet implemented.
		).Update(NewUpdate("numbers").Set("n", ArgInt(0)).Where(Expression("n=my_cte.n*my_cte.n"))).
			Recursive()

		compareToSQL(t, cte, nil,
			"WITH RECURSIVE\n`my_cte` (`n`) AS ((SELECT 1)\nUNION ALL\n(SELECT 1+n FROM `my_cte` WHERE (`n` < ?)))\nUPDATE `numbers` SET `n`=? WHERE (n=my_cte.n*my_cte.n)",
			"WITH RECURSIVE\n`my_cte` (`n`) AS ((SELECT 1)\nUNION ALL\n(SELECT 1+n FROM `my_cte` WHERE (`n` < 6)))\nUPDATE `numbers` SET `n`=0 WHERE (n=my_cte.n*my_cte.n)",
			int64(6), int64(0),
		)
		//WITH RECURSIVE my_cte(n) AS
		//(
		//	SELECT 1
		//UNION ALL
		//SELECT 1+n FROM my_cte WHERE n<6
		//)
		//UPDATE numbers, my_cte
		//# Change to 0...
		//	SET numbers.n=0
		//# ... the numbers which are squares, i.e. 1 and 4
		//WHERE numbers.n=my_cte.n*my_cte.n;
	})

	t.Run("error EMPTY top clause", func(t *testing.T) {
		cte := NewWith(
			WithCTE{Name: "check_vals", Columns: []string{"val"}, Select: NewSelect().AddColumnsExpr("123")},
		)
		compareToSQL(t, cte, errors.IsEmpty,
			"",
			"",
		)
	})
}
