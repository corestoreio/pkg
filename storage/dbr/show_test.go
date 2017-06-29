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

import "testing"

func TestShow(t *testing.T) {
	t.Parallel()

	t.Run("variables", func(t *testing.T) {
		s := NewShow().Variable()
		compareToSQL(t, s, nil,
			"SHOW VARIABLES",
			"SHOW VARIABLES",
		)
	})
	t.Run("variables global", func(t *testing.T) {
		s := NewShow().Variable().Global()
		compareToSQL(t, s, nil,
			"SHOW GLOBAL VARIABLES",
			"SHOW GLOBAL VARIABLES",
		)
	})
	t.Run("variables session", func(t *testing.T) {
		s := NewShow().Variable().Session()
		compareToSQL(t, s, nil,
			"SHOW SESSION VARIABLES",
			"SHOW SESSION VARIABLES",
		)
	})
	t.Run("variables global session", func(t *testing.T) {
		s := NewShow().Variable().Session().Global()
		compareToSQL(t, s, nil,
			"SHOW SESSION VARIABLES",
			"SHOW SESSION VARIABLES",
		)
	})
	t.Run("variables LIKE", func(t *testing.T) {
		s := NewShow().Variable().Like(ArgString("aria%"))
		compareToSQL(t, s, nil,
			"SHOW VARIABLES LIKE ?",
			"SHOW VARIABLES LIKE 'aria%'",
			"aria%",
		)
	})
	t.Run("variables WHERE", func(t *testing.T) {
		s := NewShow().Variable().Where(Column("Variable_name", In.Str("basedir", "back_log")))
		compareToSQL(t, s, nil,
			"SHOW VARIABLES WHERE (`Variable_name` IN (?,?))",
			"SHOW VARIABLES WHERE (`Variable_name` IN ('basedir','back_log'))",
			"basedir",
			"back_log",
		)
	})

	t.Run("master status", func(t *testing.T) {
		s := NewShow().MasterStatus()
		compareToSQL(t, s, nil,
			"SHOW MASTER STATUS",
			"SHOW MASTER STATUS",
		)
	})

	t.Run("binary log", func(t *testing.T) {
		s := NewShow().BinaryLog()
		compareToSQL(t, s, nil,
			"SHOW BINARY LOG",
			"SHOW BINARY LOG",
		)
	})

	t.Run("status WHERE", func(t *testing.T) {
		s := NewShow().Session().Status().Where(Column("Variable_name", Like.Str("%error%")))
		compareToSQL(t, s, nil,
			"SHOW SESSION STATUS WHERE (`Variable_name` LIKE ?)",
			"SHOW SESSION STATUS WHERE (`Variable_name` LIKE '%error%')",
			"%error%",
		)
	})

	t.Run("table status WHERE", func(t *testing.T) {
		s := NewShow().TableStatus().Where(Column("Name", Like.Str("%catalog%")))
		s.UseBuildCache = true
		compareToSQL(t, s, nil,
			"SHOW TABLE STATUS WHERE (`Name` LIKE ?)",
			"SHOW TABLE STATUS WHERE (`Name` LIKE '%catalog%')",
			"%catalog%",
		)
		// twice to test the build cache
		compareToSQL(t, s, nil,
			"SHOW TABLE STATUS WHERE (`Name` LIKE ?)",
			"SHOW TABLE STATUS WHERE (`Name` LIKE '%catalog%')",
			"%catalog%",
		)
	})

}
