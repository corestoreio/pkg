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

package dml

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTransactionReal(t *testing.T) {
	s := createRealSessionWithFixtures(t, nil)

	tx, err := s.BeginTx(context.TODO(), nil)
	assert.NoError(t, err)

	res, err := tx.InsertInto("dml_people").AddColumns("name", "email").AddValues(
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

	var person dmlPerson
	_, err = tx.SelectFrom("dml_people").Star().Where(Column("id").Int64(id)).Load(context.TODO(), &person)
	assert.NoError(t, err)

	assert.Equal(t, id, int64(person.ID))
	assert.Equal(t, "Barack", person.Name)
	assert.Equal(t, true, person.Email.Valid)
	assert.Equal(t, "obama@whitehouse.gov", person.Email.String)

	err = tx.Commit()
	assert.NoError(t, err)
}

func TestTransactionRollbackReal(t *testing.T) {
	// Insert by specifying values
	s := createRealSessionWithFixtures(t, nil)

	tx, err := s.BeginTx(context.TODO(), nil)
	assert.NoError(t, err)

	var person dmlPerson
	_, err = tx.SelectFrom("dml_people").Star().Where(Column("email").Str("jonathan@uservoice.com")).Load(context.TODO(), &person)
	assert.NoError(t, err)
	assert.Equal(t, "Jonathan", person.Name)

	err = tx.Rollback()
	assert.NoError(t, err)
}
