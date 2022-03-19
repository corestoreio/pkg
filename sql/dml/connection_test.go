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

	s.RegisterByQueryBuilder(map[string]QueryBuilder{
		"insert01": NewInsert("dml_people").AddColumns("name", "email"),
		"selectID": NewSelect("*").From("dml_people").Where(Column("id").PlaceHolder()),
	})

	lastInsertID, _ := compareExecContext(t, tx.WithCacheKey("insert01"), []any{
		"Barack", "obama@whitehouse.gov",
		"Obama", "barack@whitehouse.gov",
	}, 3, 2)

	var person dmlPerson
	_, err = tx.WithCacheKey("selectID").Load(context.TODO(), &person, lastInsertID)
	assert.NoError(t, err)

	assert.Exactly(t, lastInsertID, int64(person.ID))
	assert.Exactly(t, "Barack", person.Name)
	assert.Exactly(t, true, person.Email.Valid)
	assert.Exactly(t, "obama@whitehouse.gov", person.Email.Data)

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
	_, err = tx.WithQueryBuilder(NewSelect("*").From("dml_people").Where(Column("email").PlaceHolder())).Load(context.TODO(), &person, "SirGeorge@GoIsland.com")
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

		cp, err := NewConnPool(WithDSNFromEnv("TEST_CS_DSN_WithDSNfromEnv"))
		assert.Nil(t, cp)
		assert.ErrorIsKind(t, errors.NotImplemented, err)
	})
	t.Run("env is missing", func(t *testing.T) {
		cp, err := NewConnPool(WithDSNFromEnv("TEST_CS_DSN_WithDSNFromEnv2"))
		assert.Nil(t, cp)
		assert.ErrorIsKind(t, errors.NotExists, err)
	})
}

func Test_hashSQL(t *testing.T) {
	tests := []struct {
		name string
		args string
		want string
	}{
		{"empty", "", "cbf29ce484222325"},
		{"one char", "a", "af63fc4c860222ec"},
		{"select01", "SELECT * FROM dual", "SELECT6928bed45f95652f"},
		{"select02", "SELECT*FROM\tdual", "SELECT*FROM6928bed45f95652f"},
		{"select03", sqlIDPrefix + "asdfasfasd*/SELECT\ncol1 FROM\tdual", "SELECTfb435c798e219c6e"},
		{"update04", sqlIDPrefix + "asdfasfasd*/UPDATE\r\ncol1 FROM\tdual", "UPDATE27ae3c89bd8e8d9f"},
		{"union", "(SELECT `a`, `d` AS `b` FROM `tableAD`) UNION (SELECT `a`, `b` FROM `tableAB` WHERE (`b` = 3.14159)) ORDER BY `a`, `b` DESC, concat(\"c\",b,\"d\")", "(SELECT13d6ec028c4394e5"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := hashSQL(tt.args); got != tt.want {
				t.Errorf("hashSQL() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_extractSQLIDPrefix(t *testing.T) {
	tests := []struct {
		name        string
		rawSQL      string
		wantPrefix  string
		wantLastPos int
	}{
		{
			name:        "001",
			rawSQL:      "/*$ID$7Q26wU5GR0*/UPDATE `customer_entity_int` SET `attribute_id`=?, `entity_id`=?, `value`=? WHERE (`value_id` = ?)",
			wantPrefix:  "7Q26wU5GR0",
			wantLastPos: 18,
		},
		{
			name:        "002",
			rawSQL:      "/*$ID$7*/UPDATE",
			wantPrefix:  "7",
			wantLastPos: 9,
		},
		{
			name:        "003",
			rawSQL:      "UPDATE",
			wantPrefix:  "",
			wantLastPos: 0,
		},
		{
			name:        "004",
			rawSQL:      "/*$ID$7 UPDATE",
			wantPrefix:  "",
			wantLastPos: 0,
		},
		{
			name:        "005",
			rawSQL:      "/*$ID$*/UPDATE",
			wantPrefix:  "",
			wantLastPos: 8,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotPrefix, gotLastPos := extractSQLIDPrefix(tt.rawSQL)
			assert.Exactly(t, tt.wantLastPos, gotLastPos)
			assert.Exactly(t, tt.wantPrefix, gotPrefix)
		})
	}
}
