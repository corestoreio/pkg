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

package cstesting_test

import (
	"errors"
	"os"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/corestoreio/csfw/storage/csdb"
	"github.com/corestoreio/csfw/storage/dbr"
	"github.com/corestoreio/csfw/util/cstesting"
	"github.com/corestoreio/csfw/util/magento"
	"github.com/stretchr/testify/assert"
)

func TestMockDB(t *testing.T) {
	dbc, mockDB := cstesting.MockDB(t)
	assert.NotNil(t, dbc)
	assert.NotNil(t, mockDB)
}

func TestMustConnectDB(t *testing.T) {
	t.Parallel()

	defer func() {
		if r := recover(); r != nil {
			assert.NotNil(t, r, "There should be no panic")
		}
	}()

	db, v := cstesting.MustConnectDB()
	if v == 0 {
		assert.Nil(t, db)
	} else {
		assert.NotNil(t, db)
	}
}

func TestMustConnectDB_Mock(t *testing.T) {
	t.Parallel()

	db, sm, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}

	oldEnv := os.Getenv(csdb.EnvDSN)
	defer os.Setenv(csdb.EnvDSN, oldEnv)

	t.Run("Detect V1", func(t *testing.T) {
		// t.Parallel() not possible because sqlmock is not thread safe AND ENV vars ;-)

		os.Setenv(csdb.EnvDSN, "us3r:passw0rd@tcp(localhost:3306)/database1")

		sm.ExpectQuery("SHOW TABLES").WillReturnRows(sqlmock.NewRows([]string{"Tables_in_Database"}).
			FromCSVString("core_store\ncore_website\ncore_store_group\napi_user"),
		)

		dbc, version := cstesting.MustConnectDB(dbr.WithDB(db))

		assert.Exactly(t, magento.Version1, version)
		assert.NotNil(t, dbc)

		// 2nd call is idempotent
		dbc, version = cstesting.MustConnectDB(dbr.WithDB(db))
		assert.Exactly(t, magento.Version1, version)
		assert.NotNil(t, dbc)
	})

	t.Run("Detect V2", func(t *testing.T) {

		os.Setenv(csdb.EnvDSN, "us3r:passw0rd@tcp(localhost:3306)/database2")

		sm.ExpectQuery("SHOW TABLES").WillReturnRows(sqlmock.NewRows([]string{"Tables_in_Database"}).
			FromCSVString("integration\nstore_website\nstore_group\nauthorization_role"),
		)

		dbc, version := cstesting.MustConnectDB(dbr.WithDB(db))

		assert.Exactly(t, magento.Version2, version)
		assert.NotNil(t, dbc)

		// 2nd call is idempotent
		dbc, version = cstesting.MustConnectDB(dbr.WithDB(db))
		assert.Exactly(t, magento.Version2, version)
		assert.NotNil(t, dbc)
	})

	t.Run("Show tables error", func(t *testing.T) {
		defer func() {
			if r := recover(); r != nil {
				assert.Contains(t, r.(string), `Shop ware not supported ;-)`)
			} else {
				t.Fatal("Expecting a panic")
			}
		}()

		os.Setenv(csdb.EnvDSN, "us3r:passw0rd@tcp(localhost:3306)/database3")

		rowErr := errors.New("Shop ware not supported ;-)")
		sm.ExpectQuery("SHOW TABLES").WillReturnError(rowErr)

		dbc, version := cstesting.MustConnectDB(dbr.WithDB(db))

		assert.Exactly(t, magento.Version(0), version)
		assert.Nil(t, dbc)

	})

	if err := sm.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expections: %s", err)
	}
}
