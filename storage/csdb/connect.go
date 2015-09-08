// Copyright 2015, Cyrill @ Schumacher.fm and the CoreStore contributors
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
	"errors"
	"os"

	"github.com/corestoreio/csfw/storage/dbr"
	_ "github.com/go-sql-driver/mysql"
	"github.com/juju/errgo"
)

const (
	// EnvDSN is the name of the environment variable
	EnvDSN string = "CS_DSN"
	// EnvDSNTest test env DSN
	EnvDSNTest string = "CS_DSN_TEST"
)

var (
	ErrDSNNotFound     = errors.New("Env var: " + EnvDSN + " not found")
	ErrDSNTestNotFound = errors.New("Env var: " + EnvDSNTest + " not found")
)

func getDSN(env string, err error) (string, error) {
	dsn := os.Getenv(env)
	if dsn == "" {
		return "", err
	}
	return dsn, nil
}

// GetDSN returns the data source name from an environment variable or an error
func GetDSN() (string, error) {
	return getDSN(EnvDSN, ErrDSNNotFound)
}

// GetDSNTest returns the test data source name from an environment variable or an error
func GetDSNTest() (string, error) {
	return getDSN(EnvDSNTest, ErrDSNTestNotFound)
}

// Connect creates a new database connection from a DSN stored in an
// environment variable.
func Connect(opts ...dbr.ConnectionOption) (*dbr.Connection, error) {
	dsn, err := GetDSN()
	if err != nil {
		return nil, errgo.Mask(err)
	}
	c, err := dbr.NewConnection(dbr.SetDSN(dsn))
	return c.ApplyOpts(opts...), err
}

// MustConnectTest is a helper function that creates a
// new database connection using environment variables.
func MustConnectTest(opts ...dbr.ConnectionOption) *dbr.Connection {
	dsn, err := GetDSNTest()
	if err != nil {
		panic(err)
	}
	return dbr.MustConnectAndVerify(dbr.SetDSN(dsn)).ApplyOpts(opts...)
}
