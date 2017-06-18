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

func TestStmtChecker(t *testing.T) {
	tests := []struct {
		sel   string
		selok bool
		upd   string
		updok bool
		del   string
		delok bool
		ins   string
		insok bool
	}{
		{
			"SELECT ...",
			false,
			"UPDATE ...",
			false,
			"DELETE ...",
			false,
			"INSERT",
			false,
		},
		{
			"SELECT ... From ",
			true,
			"UPDATE ... From ",
			true,
			"DELETE ...From ",
			true,
			"INSERT ",
			true,
		},
	}

	for _, test := range tests {
		assert.Equal(t, test.selok, Stmt.IsSelect(test.sel), "%#v", test)
		assert.Equal(t, test.updok, Stmt.IsUpdate(test.upd), "%#v", test)
		assert.Equal(t, test.delok, Stmt.IsDelete(test.del), "%#v", test)
		assert.Equal(t, test.insok, Stmt.IsInsert(test.ins), "%#v", test)
	}
}
