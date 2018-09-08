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
	"testing"

	"github.com/corestoreio/errors"
	"github.com/corestoreio/pkg/util/assert"
)

func TestSQLCase_PlaceHolders(t *testing.T) {

	t.Run("placeholders in case", func(t *testing.T) {
		var start int64 = 20180101
		var end int64 = 20180201

		s := NewSelect().AddColumns("price", "sku", "name", "title", "description").
			AddColumnsConditions(
				SQLCase("", "`closed`",
					"date_start <= ? AND date_end >= ?", "`open`",
					"date_start > ? AND date_end > ?", "`upcoming`",
				).Alias("is_on_sale"),
			).
			From("catalog_promotions").
			Where(
				Column("promotion_id").NotIn().PlaceHolder(),
			)
		sa := s.WithArgs().Int64(start).Int64(end).Int64(start).Int64(end).Ints(4711, 815, 42)

		compareToSQL(t, sa, errors.NoKind,
			"SELECT `price`, `sku`, `name`, `title`, `description`, CASE  WHEN date_start <= ? AND date_end >= ? THEN `open` WHEN date_start > ? AND date_end > ? THEN `upcoming` ELSE `closed` END AS `is_on_sale` FROM `catalog_promotions` WHERE (`promotion_id` NOT IN ?)",
			"",
			start, end, start, end, int64(4711), int64(815), int64(42),
		)

		assert.Exactly(t, []string{"promotion_id"}, s.qualifiedColumns)
	})

}
