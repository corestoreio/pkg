// Copyright 2015 CoreStore Authors
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
	"database/sql"
	"errors"
	"os"

	"github.com/gocraft/dbr"
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

// GetDSN returns the DSN from env or an error
func GetDSN() (string, error) {
	return getDSN(EnvDSN, ErrDSNNotFound)
}

// GetDSNTest returns the DSN from env or an error
func GetDSNTest() (string, error) {
	return getDSN(EnvDSNTest, ErrDSNTestNotFound)
}

func Connect() (*sql.DB, *dbr.Session, error) {
	dsn, err := GetDSN()
	if err != nil {
		return nil, nil, errgo.Mask(err)
	}
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, nil, errgo.Mask(err)
	}
	dbrSess := dbr.NewConnection(db, nil).NewSession(nil)

	return db, dbrSess, nil
}

// mustConnectTest is a helper function that creates a
// new database connection using environment variables.
func MustConnectTest() *sql.DB {
	dsn, err := GetDSNTest()
	if err != nil {
		panic(err)
	}
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		panic(err)
	}
	return db
}
