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

	"github.com/stretchr/testify/assert"
)

func TestWith_ToSQL(t *testing.T) {

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

}
