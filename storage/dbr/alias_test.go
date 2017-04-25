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

package dbr

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMakeAlias(t *testing.T) {
	assert.Exactly(t, "`table1`", MakeAlias("table1").String())
	assert.Exactly(t, "`table0` AS `table1`", MakeAlias("table0", "table1").String())
	assert.Exactly(t, "(table1)", MakeAlias("(table1)").String())
	assert.Exactly(t, "(table1) AS `table2`", MakeAlias("(table1)", "table2").String())
	assert.Exactly(t, "(table1)", MakeAlias("(table1)", "").String())
	assert.Exactly(t, "`table1`", MakeAlias("table1", "").String())
}
