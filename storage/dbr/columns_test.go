// Copyright 2015-2016, Cyrill @ Schumacher.fm and the CoreStore contributors
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

package dbr_test

import (
	"testing"

	"github.com/corestoreio/csfw/storage/dbr"
	"github.com/stretchr/testify/assert"
)

func TestTableColumnQuote(t *testing.T) {
	t.Parallel()
	tests := []struct {
		haveT string
		haveC []string
		want  []string
	}{
		{
			"t1",
			[]string{"col1", "col2"},
			[]string{"`t1`.`col1`", "`t1`.`col2`"},
		},
		{
			"t2",
			[]string{"col1", "col2", "`t2`.`col3`"},
			[]string{"`t2`.`col1`", "`t2`.`col2`", "`t2`.`col3`"},
		},
		{
			"t2a",
			[]string{"col1", "col2", "t2.col3"},
			[]string{"`t2a`.`col1`", "`t2a`.`col2`", "`t2`.`col3`"},
		},
		{
			"t3",
			[]string{"col1", "col2", "`col3`"},
			[]string{"`t3`.`col1`", "`t3`.`col2`", "`col3`"},
		},
	}

	for i, test := range tests {
		actC := dbr.Quoter.TableColumnAlias(test.haveT, test.haveC...)
		assert.Equal(t, test.want, actC, "Index %d", i)
	}
}

func TestIfNullAs(t *testing.T) {
	t.Parallel()
	runner := func(want string, have ...string) func(*testing.T) {
		return func(t *testing.T) {
			assert.Equal(t, want, dbr.IfNull(have...))
		}
	}
	t.Run("1 args expression", runner(
		"IFNULL((1/0),(NULL ))",
		"1/0",
	))
	t.Run("1 args no qualifier", runner(
		"IFNULL(`c1`,(NULL ))",
		"c1",
	))
	t.Run("1 args with qualifier", runner(
		"IFNULL(`t1`.`c1`,(NULL ))",
		"t1.c1",
	))

	t.Run("2 args expression left", runner(
		"IFNULL((1/0),`c2`)",
		"1/0", "c2",
	))
	t.Run("2 args expression right", runner(
		"IFNULL(`c2`,(1/0))",
		"c2", "1/0",
	))
	t.Run("2 args no qualifier", runner(
		"IFNULL(`c1`,`c2`)",
		"c1", "c2",
	))
	t.Run("2 args with qualifier", runner(
		"IFNULL(`t1`.`c1`,`t2`.`c2`)",
		"t1.c1", "t2.c2",
	))

	t.Run("3 args expression left", runner(
		"IFNULL((1/0),`c2`) AS `alias`",
		"1/0", "c2", "alias",
	))
	t.Run("3 args expression right", runner(
		"IFNULL(`c2`,(1/0)) AS `alias`",
		"c2", "1/0", "alias",
	))
	t.Run("3 args no qualifier", runner(
		"IFNULL(`c1`,`c2`) AS `alias`",
		"c1", "c2", "alias",
	))
	t.Run("3 args with qualifier", runner(
		"IFNULL(`t1`.`c1`,`t2`.`c2`) AS `alias`",
		"t1.c1", "t2.c2", "alias",
	))

	t.Run("4 args", runner(
		"IFNULL(`t1`.`c1`,`t2`.`c2`)",
		"t1", "c1", "t2", "c2",
	))
	t.Run("5 args", runner(
		"IFNULL(`t1`.`c1`,`t2`.`c2`) AS `alias`",
		"t1", "c1", "t2", "c2", "alias",
	))
	t.Run("6 args", runner(
		"IFNULL(`t1`.`c1`,`t2`.`c2`) AS `alias_x`",
		"t1", "c1", "t2", "c2", "alias", "x",
	))
	t.Run("7 args", runner(
		"IFNULL(`t1`.`c1`,`t2`.`c2`) AS `alias_x_y`",
		"t1", "c1", "t2", "c2", "alias", "x", "y",
	))
}

func BenchmarkIfNull(b *testing.B) {
	runner := func(want string, have ...string) func(*testing.B) {
		return func(b *testing.B) {
			var result string
			for i := 0; i < b.N; i++ {
				result = dbr.IfNull(have...)
			}
			if result != want {
				b.Fatalf("\nHave: %q\nWant: %q", result, want)
			}
		}
	}
	b.Run("3 args expression right", runner(
		"IFNULL(`c2`,(1/0)) AS `alias`",
		"c2", "1/0", "alias",
	))
	b.Run("3 args no qualifier", runner(
		"IFNULL(`c1`,`c2`) AS `alias`",
		"c1", "c2", "alias",
	))
	b.Run("3 args with qualifier", runner(
		"IFNULL(`t1`.`c1`,`t2`.`c2`) AS `alias`",
		"t1.c1", "t2.c2", "alias",
	))

	b.Run("4 args", runner(
		"IFNULL(`t1`.`c1`,`t2`.`c2`)",
		"t1", "c1", "t2", "c2",
	))
	b.Run("5 args", runner(
		"IFNULL(`t1`.`c1`,`t2`.`c2`) AS `alias`",
		"t1", "c1", "t2", "c2", "alias",
	))
	b.Run("6 args", runner(
		"IFNULL(`t1`.`c1`,`t2`.`c2`) AS `alias_x`",
		"t1", "c1", "t2", "c2", "alias", "x",
	))

}
