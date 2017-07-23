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
	"fmt"
	"testing"

	"github.com/corestoreio/errors"
	"github.com/stretchr/testify/assert"
)

// check if the types implement the interfaces

var _ fmt.Stringer = (*Delete)(nil)
var _ fmt.Stringer = (*Insert)(nil)
var _ fmt.Stringer = (*Update)(nil)
var _ fmt.Stringer = (*Select)(nil)

var _ QueryBuilder = (*Select)(nil)
var _ QueryBuilder = (*Delete)(nil)
var _ QueryBuilder = (*Update)(nil)
var _ QueryBuilder = (*Insert)(nil)
var _ queryBuilder = (*buildQueryMock)(nil)

type buildQueryMock struct{ error }

func (m buildQueryMock) toSQL(queryWriter) error { return m.error }

func (m buildQueryMock) appendArgs(Arguments) (Arguments, error) { return nil, m.error }
func (m buildQueryMock) hasBuildCache() bool                     { return false }
func (m buildQueryMock) writeBuildCache(sql []byte)              {}
func (m buildQueryMock) readBuildCache() (sql []byte, args Arguments, err error) {
	return nil, nil, m.error
}

func TestMakeSQL(t *testing.T) {
	t.Parallel()
	t.Run("error", func(t *testing.T) {
		s := makeSQL(buildQueryMock{errors.NewAbortedf("Canceled")}, false)
		assert.Contains(t, s, "[dbr] ToSQL Error: Canceled\n")
	})
	t.Run("DELETE", func(t *testing.T) {
		b := NewDelete("tableX").Where(Column("columnA").Greater().Int64(2))
		assert.Exactly(t, "DELETE FROM `tableX` WHERE (`columnA` > ?)", b.String())
	})
	t.Run("INSERT", func(t *testing.T) {
		b := NewInsert("tableX").AddColumns("columnA", "columnB").AddValues(2, "Go")
		assert.Exactly(t, "INSERT INTO `tableX` (`columnA`,`columnB`) VALUES (?,?)", b.String())
	})
	t.Run("SELECT", func(t *testing.T) {
		b := NewSelect("columnA").FromAlias("tableX", "X").Where(Column("columnA").LessOrEqual().Float64(2.4))
		assert.Exactly(t, "SELECT `columnA` FROM `tableX` AS `X` WHERE (`columnA` <= ?)", b.String())
	})
	t.Run("UPDATE", func(t *testing.T) {
		b := NewUpdate("tableX").Set(
			Column("columnA").Int64(4),
		).Where(Column("columnB").Between().Ints(5, 7))
		assert.Exactly(t, "UPDATE `tableX` SET `columnA`=? WHERE (`columnB` BETWEEN ? AND ?)", b.String())
	})
}
