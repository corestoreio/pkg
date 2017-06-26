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
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTransactionReal(t *testing.T) {
	s := createRealSessionWithFixtures(t)

	tx, err := s.Begin()
	assert.NoError(t, err)

	res, err := tx.InsertInto("dbr_people").AddColumns("name", "email").AddValues(
		"Barack", "obama@whitehouse.gov",
		"Obama", "barack@whitehouse.gov",
	).Exec(context.TODO())

	assert.NoError(t, err)
	id, err := res.LastInsertId()
	assert.NoError(t, err)
	rowsAff, err := res.RowsAffected()
	assert.NoError(t, err)

	assert.True(t, id > 0)
	assert.Equal(t, int64(2), rowsAff)

	var person dbrPerson
	_, err = tx.Select("*").From("dbr_people").Where(Column("id", ArgInt64(id))).Load(context.TODO(), &person)
	assert.NoError(t, err)

	assert.Equal(t, id, person.ID)
	assert.Equal(t, "Barack", person.Name)
	assert.Equal(t, true, person.Email.Valid)
	assert.Equal(t, "obama@whitehouse.gov", person.Email.String)

	err = tx.Commit()
	assert.NoError(t, err)
}

func TestTransactionRollbackReal(t *testing.T) {
	// Insert by specifying values
	s := createRealSessionWithFixtures(t)

	tx, err := s.Begin()
	assert.NoError(t, err)

	var person dbrPerson
	_, err = tx.Select("*").From("dbr_people").Where(Column("email", ArgString("jonathan@uservoice.com"))).Load(context.TODO(), &person)
	assert.NoError(t, err)
	assert.Equal(t, "Jonathan", person.Name)

	err = tx.Rollback()
	assert.NoError(t, err)
}
