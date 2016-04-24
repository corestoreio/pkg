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

package csdb

import (
	"os"

	"github.com/corestoreio/csfw/storage/dbr"
	"github.com/corestoreio/csfw/util/errors"
	_ "github.com/go-sql-driver/mysql"
)

// EnvDSN is the name of the environment variable
const EnvDSN string = "CS_DSN"

func getDSN(env string, err error) (string, error) {
	dsn := os.Getenv(env)
	if dsn == "" {
		return "", err
	}
	return dsn, nil
}

// GetDSN returns the data source name from an environment variable or an error
func GetDSN() (string, error) {
	return getDSN(EnvDSN, errors.NewNotFoundf("Env var: %q not found", EnvDSN))
}

// Connect creates a new database connection from a DSN stored in an
// environment variable.
func Connect(opts ...dbr.ConnectionOption) (*dbr.Connection, error) {
	dsn, err := GetDSN()
	if err != nil {
		return nil, errors.Wrap(err, "[csdb] GetDSN")
	}
	c, err := dbr.NewConnection(dbr.WithDSN(dsn))
	return c.ApplyOpts(opts...), err
}

// MustConnectTest is a helper function that creates a
// new database connection using environment variables.
func MustConnectTest(opts ...dbr.ConnectionOption) *dbr.Connection {
	dsn, err := GetDSN()
	if err != nil {
		panic(err)
	}
	return dbr.MustConnectAndVerify(dbr.WithDSN(dsn)).ApplyOpts(opts...)
}
