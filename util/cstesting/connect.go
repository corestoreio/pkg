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

package cstesting

import (
	"fmt"
	"os"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/corestoreio/csfw/storage/dbr"
	"github.com/corestoreio/errors"
)

// EnvDSN is the name of the environment variable
const EnvDSN string = "CS_DSN"

func getDSN(env string) (string, error) {
	dsn := os.Getenv(env)
	if dsn == "" {
		return "", errors.NewNotFoundf("DSN in environment variable %q not found.", EnvDSN)
	}
	return dsn, nil
}

// MustGetDSN returns the data source name from an environment variable or
// panics on error.
func MustGetDSN() string {
	d, err := getDSN(EnvDSN)
	if err != nil {
		panic(fmt.Sprintf("%+v", err))
	}
	return d
}

// MustConnectDB is a helper function that creates a new database connection
// using a DSN from an environment variable found in the constant csdb.EnvDSN.
// It queries the database to figure out the current version of Magento. If the
// DSN environment variable has not been set it returns nil,0.
func MustConnectDB(opts ...dbr.ConnectionOption) *dbr.Connection {
	cos := append([]dbr.ConnectionOption{}, dbr.WithDSN(MustGetDSN()))
	dbc := dbr.MustConnectAndVerify(append(cos, opts...)...)
	return dbc
}

// MockDB creates a mocked database connection. Fatals on error.
func MockDB(t fataler) (*dbr.Connection, sqlmock.Sqlmock) {
	db, sm, err := sqlmock.New()
	fatalIfError(t, err)

	dbc, err := dbr.NewConnection(dbr.WithDB(db))
	fatalIfError(t, err)
	return dbc, sm
}
