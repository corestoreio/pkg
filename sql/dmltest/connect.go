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

package dmltest

import (
	"io"
	"os"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/corestoreio/errors"
	"github.com/corestoreio/pkg/sql/dml"
)

// EnvDSN is the name of the environment variable
const EnvDSN = "CS_DSN"

func getDSN(env string) (string, error) {
	dsn := os.Getenv(env)
	if dsn == "" {
		return "", errors.NotFound.Newf("DSN in environment variable %q not found.", EnvDSN)
	}
	return dsn, nil
}

// MustGetDSN returns the data source name from an environment variable or
// panics on error.
func MustGetDSN(t testing.TB) string {
	d, err := getDSN(EnvDSN)
	if err != nil {
		t.Skip(err)
	}
	return d
}

// MustConnectDB is a helper function that creates a new database connection
// using a DSN from an environment variable found in the constant csdb.EnvDSN.
// If the DSN environment variable has not been set it skips the test.
// Argument t specified usually the *testing.T/B struct.
func MustConnectDB(t testing.TB, opts ...dml.ConnPoolOption) *dml.ConnPool {
	t.Helper()
	if _, err := getDSN(EnvDSN); errors.NotFound.Match(err) {
		t.Skipf("%s", err)
	}
	if len(opts) == 0 {
		return dml.MustConnectAndVerify(dml.WithDSN(MustGetDSN(t)))
	}
	return dml.MustConnectAndVerify(opts...)
}

// Close for usage in conjunction with defer.
// 		defer cstesting.Close(t,db)
func Close(t testing.TB, c io.Closer) {
	t.Helper()
	if err := c.Close(); err != nil {
		t.Errorf("%+v", err)
	}
}

// MockDB creates a mocked database connection. Fatals on error.
func MockDB(t testing.TB) (*dml.ConnPool, sqlmock.Sqlmock) {
	if t != nil { // t can be nil in Example functions
		t.Helper()
	}
	db, sm, err := sqlmock.New()
	fatalIfError(t, err)

	dbc, err := dml.NewConnPool(dml.WithDB(db))
	fatalIfError(t, err)
	return dbc, sm
}

// MockClose for usage in conjunction with defer.
// 		defer cstesting.MockClose(t,db,dbMock)
func MockClose(t testing.TB, c io.Closer, m sqlmock.Sqlmock) {
	if t != nil { // t can be nil in Example functions
		t.Helper()
	}
	m.ExpectClose()
	if err := c.Close(); err != nil {
		if t != nil {
			t.Fatalf("%+v", err)
		} else {
			panic(err)
		}
	}
	if err := m.ExpectationsWereMet(); err != nil {
		if t != nil {
			t.Fatalf("There were unfulfilled expectations:\n%+v", err)
		} else {
			panic(err)
		}
	}
}

// fatalIfError fails the tests if an unexpected error occurred. If the error is
// gift wrapped prints the location.
func fatalIfError(t testing.TB, err error) {
	if err != nil {
		if t != nil {
			t.Fatalf("%+v", err)
		} else {
			panic(err)
		}
	}
}
