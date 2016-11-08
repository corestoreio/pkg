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
	"context"
	"database/sql"
	"fmt"
	"os"
	"sync"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/corestoreio/csfw/storage/dbr"
	"github.com/corestoreio/csfw/util/errors"
	"github.com/corestoreio/csfw/util/magento"
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

type dsnVersion struct {
	sync.RWMutex
	cache map[string]magento.Version
}

func (dv *dsnVersion) v(dsn string) magento.Version {
	dv.RLock()
	defer dv.RUnlock()
	return dv.cache[dsn]
}

func (dv *dsnVersion) set(dsn string, v magento.Version) {
	dv.Lock()
	defer dv.Unlock()
	dv.cache[dsn] = v
}

// dsnVersion stores for a specific DSN the Magento version of the database.
// This struct avoids rescanning the database tables.
var dsnVersionCache = &dsnVersion{
	cache: make(map[string]magento.Version),
}

// MustConnectDB is a helper function that creates a new database connection
// using a DSN from an environment variable found in the constant csdb.EnvDSN.
// It queries the database to figure out the current version of Magento. If the
// DSN environment variable has not been set it returns nil,0.
func MustConnectDB(opts ...dbr.ConnectionOption) (*dbr.Connection, magento.Version) {
	dsn, err := getDSN(EnvDSN)
	if errors.IsNotFound(err) {
		return nil, 0
	}

	cos := make([]dbr.ConnectionOption, 0, 2)
	cos = append(cos, dbr.WithDSN(MustGetDSN()))
	dbc := dbr.MustConnectAndVerify(append(cos, opts...)...)

	if v := dsnVersionCache.v(dsn); v > 0 {
		return dbc, v
	}

	tables, err := showTables(context.Background(), dbc.DB)
	if err != nil {
		panic(fmt.Sprintf("%+v", err))
	}

	v := magento.DetectVersion("", tables) // hmm refactor later if we need a prefix
	dsnVersionCache.set(dsn, v)

	return dbc, v
}

// MockDB creates a mocked database connection. Fatals on error.
func MockDB(t fataler) (*dbr.Connection, sqlmock.Sqlmock) {
	db, sm, err := sqlmock.New()
	fatalIfError(t, err)

	dbc, err := dbr.NewConnection(dbr.WithDB(db))
	fatalIfError(t, err)
	return dbc, sm
}

// showTables executes the query SHOW TABLES and returns all tables within the
// current database.
func showTables(ctx context.Context, db *sql.DB) ([]string, error) {
	rows, err := db.QueryContext(ctx, "SHOW TABLES")
	if err != nil {
		return nil, errors.Wrap(err, "[csdb] ShowTables: SHOW TABLES failed")
	}
	defer rows.Close()

	var tables = make([]string, 0, 200)
	var table = new(string)
	for rows.Next() {
		if err := rows.Scan(table); err != nil {
			return nil, errors.Wrap(err, "[csdb] ShowTables: scan failed")
		}
		tables = append(tables, *table)
		*table = ""
	}
	if rows.Err() != nil {
		return nil, errors.Wrap(rows.Err(), "[csdb] ShowTables: row error")
	}
	return tables, nil
}
