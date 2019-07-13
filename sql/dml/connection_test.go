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
	"context"
	"os"
	"testing"

	"github.com/corestoreio/errors"
	"github.com/corestoreio/pkg/util/assert"
)

func TestTransactionReal(t *testing.T) {
	s := createRealSessionWithFixtures(t, nil)
	defer testCloser(t, s)

	tx, err := s.BeginTx(context.TODO(), nil)
	assert.NoError(t, err)

	txIns := tx.InsertInto("dml_people").AddColumns("name", "email").WithArgs().Raw(
		"Barack", "obama@whitehouse.gov",
		"Obama", "barack@whitehouse.gov",
	)

	lastInsertID, _ := compareExecContext(t, txIns, 3, 2)

	var person dmlPerson
	_, err = tx.SelectFrom("dml_people").Star().Where(Column("id").Int64(lastInsertID)).WithArgs().Load(context.TODO(), &person)
	assert.NoError(t, err)

	assert.Exactly(t, lastInsertID, int64(person.ID))
	assert.Exactly(t, "Barack", person.Name)
	assert.Exactly(t, true, person.Email.Valid)
	assert.Exactly(t, "obama@whitehouse.gov", person.Email.String)

	err = tx.Commit()
	assert.NoError(t, err)
}

func TestTransactionRollbackReal(t *testing.T) {
	// Insert by specifying values
	s := createRealSessionWithFixtures(t, nil)
	defer testCloser(t, s)

	tx, err := s.BeginTx(context.TODO(), nil)
	assert.NoError(t, err)

	var person dmlPerson
	_, err = tx.SelectFrom("dml_people").Star().Where(Column("email").PlaceHolder()).WithArgs().Load(context.TODO(), &person, "SirGeorge@GoIsland.com")
	assert.NoError(t, err)
	assert.Exactly(t, "Sir George", person.Name)

	err = tx.Rollback()
	assert.NoError(t, err)
}

func TestWithDSNfromEnv(t *testing.T) {
	t.Run("incorrect env", func(t *testing.T) {
		os.Setenv("TEST_CS_DSN_WithDSNfromEnv", "errrrr")
		defer func() {
			os.Unsetenv("TEST_CS_DSN_WithDSNfromEnv")
		}()

		cp, err := NewConnPool(WithDSNfromEnv("TEST_CS_DSN_WithDSNfromEnv"))
		assert.Nil(t, cp)
		assert.ErrorIsKind(t, errors.NotImplemented, err)
	})
	t.Run("env is missing", func(t *testing.T) {
		cp, err := NewConnPool(WithDSNfromEnv("TEST_CS_DSN_WithDSNFromEnv2"))
		assert.Nil(t, cp)
		assert.ErrorIsKind(t, errors.NotExists, err)
	})
}
