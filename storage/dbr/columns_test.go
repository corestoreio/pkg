// Copyright 2015, Cyrill @ Schumacher.fm and the CoreStore contributors
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

func TestTableColumnQuote(t *testing.T) {
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
			"t3",
			[]string{"col1", "col2", "`col3`"},
			[]string{"`t3`.`col1`", "`t3`.`col2`", "`col3`"},
		},
	}

	for _, test := range tests {
		actC := TableColumnQuote(test.haveT, test.haveC...)
		assert.Equal(t, test.want, actC)
	}
}

func TestIfNullAs(t *testing.T) {
	s := IfNullAs("t1", "c1", "t2", "c2", "alias")
	assert.Equal(t, "IFNULL(`t1`.`c1`, `t2`.`c2`) AS `alias`", s)
}
