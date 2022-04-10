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

	"github.com/corestoreio/pkg/util/assert"
)

func TestShow(t *testing.T) {
	t.Run("variables", func(t *testing.T) {
		s := NewShow().Variable()
		compareToSQL(t, s, false,
			"SHOW VARIABLES",
			"SHOW VARIABLES",
		)
	})
	t.Run("variables global", func(t *testing.T) {
		s := NewShow().Variable().Global()
		compareToSQL(t, s, false,
			"SHOW GLOBAL VARIABLES",
			"SHOW GLOBAL VARIABLES",
		)
	})
	t.Run("variables session", func(t *testing.T) {
		s := NewShow().Variable().Session()
		compareToSQL(t, s, false,
			"SHOW SESSION VARIABLES",
			"SHOW SESSION VARIABLES",
		)
	})
	t.Run("variables global session", func(t *testing.T) {
		s := NewShow().Variable().Session().Global()
		compareToSQL(t, s, false,
			"SHOW SESSION VARIABLES",
			"SHOW SESSION VARIABLES",
		)
	})
	t.Run("variables LIKE interpolated", func(t *testing.T) {
		s := NewShow().Variable().Like().WithDBR(dbMock{}).TestWithArgs("aria%")
		compareToSQL(t, s, false,
			"SHOW VARIABLES LIKE ?",
			"SHOW VARIABLES LIKE 'aria%'",
			"aria%",
		)
	})
	t.Run("variables LIKE place holder", func(t *testing.T) {
		s := NewShow().Variable().
			Where(Column("Variable_name").PlaceHolder()).
			WithDBR(dbMock{})
		compareToSQL(t, s.TestWithArgs("aria%"), false,
			"SHOW VARIABLES WHERE (`Variable_name` = ?)",
			"SHOW VARIABLES WHERE (`Variable_name` = 'aria%')",
			"aria%",
		)
		assert.Exactly(t, []string{"Variable_name"}, s.cachedSQL.qualifiedColumns)
	})
	t.Run("variables WHERE interpolate", func(t *testing.T) {
		s := NewShow().Variable().Where(Column("Variable_name").In().Strs("basedir", "back_log"))
		compareToSQL(t, s, false,
			"SHOW VARIABLES WHERE (`Variable_name` IN ('basedir','back_log'))",
			"SHOW VARIABLES WHERE (`Variable_name` IN ('basedir','back_log'))",
		)
	})
	t.Run("variables WHERE placeholder", func(t *testing.T) {
		s := NewShow().Variable().
			Where(Column("Variable_name").In().PlaceHolder()).
			WithDBR(dbMock{}).TestWithArgs([]string{"basedir", "back_log"})
		compareToSQL(t, s, false,
			"SHOW VARIABLES WHERE (`Variable_name` IN ?)",
			"SHOW VARIABLES WHERE (`Variable_name` IN ('basedir','back_log'))",
			"basedir",
			"back_log",
		)
	})

	t.Run("master status", func(t *testing.T) {
		s := NewShow().MasterStatus()
		compareToSQL(t, s, false,
			"SHOW MASTER STATUS",
			"SHOW MASTER STATUS",
		)
	})

	t.Run("binary log", func(t *testing.T) {
		s := NewShow().BinaryLog()
		compareToSQL(t, s, false,
			"SHOW BINARY LOG",
			"SHOW BINARY LOG",
		)
	})

	t.Run("status WHERE", func(t *testing.T) {
		s := NewShow().Session().Status().Where(Column("Variable_name").Like().Str("%error%"))
		compareToSQL(t, s, false,
			"SHOW SESSION STATUS WHERE (`Variable_name` LIKE '%error%')",
			"SHOW SESSION STATUS WHERE (`Variable_name` LIKE '%error%')",
		)
	})
}
