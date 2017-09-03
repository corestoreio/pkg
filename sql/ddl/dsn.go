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

package ddl

import (
	"fmt"
	"os"

	"github.com/corestoreio/errors"
	"github.com/go-sql-driver/mysql"
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
	return getDSN(EnvDSN, errors.NewNotFoundf("DSN in environment variable %q not found.", EnvDSN))
}

// MustGetDSN returns the data source name from an environment variable or
// panics on error.
func MustGetDSN() string {
	d, err := GetDSN()
	if err != nil {
		panic(fmt.Sprintf("%+v", err))
	}
	return d
}

// GetParsedDSN checks the environment variable EnvDSN and if a DSN can be found
// it gets parsed into an URL.
func GetParsedDSN() (*mysql.Config, error) {
	dsn, err := GetDSN()
	if err != nil {
		return nil, errors.Wrap(err, "[ddl] Cannot find DSN environment variable")
	}
	pd, err := mysql.ParseDSN(dsn)
	if err != nil {
		return nil, errors.Wrap(err, "[ddl] Cannot parse DSN into URL")
	}
	return pd, nil
}
