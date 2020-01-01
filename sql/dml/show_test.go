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

func TestShow(t *testing.T) {
	t.Parallel()

	t.Run("variables", func(t *testing.T) {
		s := NewShow().Variable()
		compareToSQL(t, s, errors.NoKind,
			"SHOW VARIABLES",
			"SHOW VARIABLES",
		)
	})
	t.Run("variables global", func(t *testing.T) {
		s := NewShow().Variable().Global()
		compareToSQL(t, s, errors.NoKind,
			"SHOW GLOBAL VARIABLES",
			"SHOW GLOBAL VARIABLES",
		)
	})
	t.Run("variables session", func(t *testing.T) {
		s := NewShow().Variable().Session()
		compareToSQL(t, s, errors.NoKind,
			"SHOW SESSION VARIABLES",
			"SHOW SESSION VARIABLES",
		)
	})
	t.Run("variables global session", func(t *testing.T) {
		s := NewShow().Variable().Session().Global()
		compareToSQL(t, s, errors.NoKind,
			"SHOW SESSION VARIABLES",
			"SHOW SESSION VARIABLES",
		)
	})
	t.Run("variables LIKE interpolated", func(t *testing.T) {
		s := NewShow().Variable().Like().WithDBR().TestWithArgs("aria%")
		compareToSQL(t, s, errors.NoKind,
			"SHOW VARIABLES LIKE ?",
			"SHOW VARIABLES LIKE 'aria%'",
			"aria%",
		)
	})
	t.Run("variables LIKE place holder", func(t *testing.T) {
		s := NewShow().Variable().
			Where(Column("Variable_name").PlaceHolder()).
			WithDBR()
		compareToSQL(t, s.TestWithArgs("aria%"), errors.NoKind,
			"SHOW VARIABLES WHERE (`Variable_name` = ?)",
			"SHOW VARIABLES WHERE (`Variable_name` = 'aria%')",
			"aria%",
		)
		assert.Exactly(t, []string{"Variable_name"}, s.base.qualifiedColumns)
	})
	t.Run("variables WHERE interpolate", func(t *testing.T) {
		s := NewShow().Variable().Where(Column("Variable_name").In().Strs("basedir", "back_log"))
		compareToSQL(t, s, errors.NoKind,
			"SHOW VARIABLES WHERE (`Variable_name` IN ('basedir','back_log'))",
			"SHOW VARIABLES WHERE (`Variable_name` IN ('basedir','back_log'))",
		)
	})
	t.Run("variables WHERE placeholder", func(t *testing.T) {
		s := NewShow().Variable().
			Where(Column("Variable_name").In().PlaceHolder()).
			WithDBR().TestWithArgs([]string{"basedir", "back_log"})
		compareToSQL(t, s, errors.NoKind,
			"SHOW VARIABLES WHERE (`Variable_name` IN ?)",
			"SHOW VARIABLES WHERE (`Variable_name` IN ('basedir','back_log'))",
			"basedir",
			"back_log",
		)
	})

	t.Run("master status", func(t *testing.T) {
		s := NewShow().MasterStatus()
		compareToSQL(t, s, errors.NoKind,
			"SHOW MASTER STATUS",
			"SHOW MASTER STATUS",
		)
	})

	t.Run("binary log", func(t *testing.T) {
		s := NewShow().BinaryLog()
		compareToSQL(t, s, errors.NoKind,
			"SHOW BINARY LOG",
			"SHOW BINARY LOG",
		)
	})

	t.Run("status WHERE", func(t *testing.T) {
		s := NewShow().Session().Status().Where(Column("Variable_name").Like().Str("%error%"))
		compareToSQL(t, s, errors.NoKind,
			"SHOW SESSION STATUS WHERE (`Variable_name` LIKE '%error%')",
			"SHOW SESSION STATUS WHERE (`Variable_name` LIKE '%error%')",
		)
	})

	t.Run("table status WHERE (build cache)", func(t *testing.T) {
		s := NewShow().TableStatus().Where(Column("Name").Regexp().Str(".*catalog[_]+"))
		compareToSQL(t, s, errors.NoKind,
			"SHOW TABLE STATUS WHERE (`Name` REGEXP '.*catalog[_]+')",
			"SHOW TABLE STATUS WHERE (`Name` REGEXP '.*catalog[_]+')",
		)
		assert.Exactly(t, []string{"", "SHOW TABLE STATUS WHERE (`Name` REGEXP '.*catalog[_]+')"}, s.CachedQueries())

		s.WithCacheKey("XsalesX").WhereFragments[0].Str("sales$") // set Equal on purpose ... because cache already written
		// twice to test the build cache
		compareToSQL(t, s, errors.NoKind,
			"SHOW TABLE STATUS WHERE (`Name` REGEXP 'sales$')",
			"SHOW TABLE STATUS WHERE (`Name` REGEXP 'sales$')",
		)
		assert.Exactly(t,
			[]string{"", "SHOW TABLE STATUS WHERE (`Name` REGEXP '.*catalog[_]+')", "XsalesX", "SHOW TABLE STATUS WHERE (`Name` REGEXP 'sales$')"},
			s.CachedQueries())
	})
}
