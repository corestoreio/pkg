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
	t.Run("variables", func(t *testing.T) {
		s := NewShow().Variable()
		compareToSQL(t, s, nil,
			"SHOW VARIABLES",
			"SHOW VARIABLES",
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

}
