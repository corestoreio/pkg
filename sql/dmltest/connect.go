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

package dmltest

import (
	"context"
	"database/sql"
	"database/sql/driver"
	"fmt"
	"io"
	"os"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/corestoreio/errors"
	"github.com/corestoreio/pkg/sql/dml"
	"github.com/fatih/color"
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
		t.Skip(color.MagentaString("%s", err))
	}
	return d
}

// MustConnectDB is a helper function that creates a new database connection
// using a DSN from an environment variable found in the constant csdb.EnvDSN.
// If the DSN environment variable has not been set it skips the test.
// It creates a random database if the DSN database name is the word "random".
func MustConnectDB(t testing.TB, opts ...dml.ConnPoolOption) *dml.ConnPool {
	t.Helper()
	if _, err := getDSN(EnvDSN); errors.NotFound.Match(err) {
		t.Skip(color.MagentaString("%s", err))
	}

	cfg := []dml.ConnPoolOption{
		dml.WithDSN(MustGetDSN(t)),
		dml.WithCreateDatabase(context.Background(), ""), // empty DB name gets derived from the DSN
	}
	return dml.MustConnectAndVerify(append(cfg, opts...)...)
}

// Close for usage in conjunction with defer.
// 		defer dmltest.Close(t, db)
func Close(t testing.TB, c io.Closer) {
	t.Helper()
	if err := c.Close(); err != nil {
		t.Errorf("%+v", err)
	}
}

// MockDB creates a mocked database connection. Fatals on error.
func MockDB(t testing.TB, opts ...dml.ConnPoolOption) (*dml.ConnPool, sqlmock.Sqlmock) {
	if t != nil { // t can be nil in Example functions
		t.Helper()
	}
	db, sm, err := sqlmock.New()
	FatalIfError(t, err)
	cfg := []dml.ConnPoolOption{dml.WithDB(db)}
	dbc, err := dml.NewConnPool(append(cfg, opts...)...)
	FatalIfError(t, err)
	return dbc, sm
}

// MockDBCallBack same as MockDB but allows to add expectations early to the
// mock.
func MockDBCallBack(t testing.TB, mockCB func(sqlmock.Sqlmock), opts ...dml.ConnPoolOption) (*dml.ConnPool, sqlmock.Sqlmock) {
	if t != nil { // t can be nil in Example functions
		t.Helper()
	}
	db, sm, err := sqlmock.New()
	FatalIfError(t, err)
	if mockCB != nil {
		mockCB(sm)
	}
	cfg := []dml.ConnPoolOption{dml.WithDB(db)}
	dbc, err := dml.NewConnPool(append(cfg, opts...)...)
	FatalIfError(t, err)
	return dbc, sm
}

// MockClose for usage in conjunction with defer.
// 		defer dmltest.MockClose(t, db, dbMock)
func MockClose(t testing.TB, c io.Closer, m sqlmock.Sqlmock) {
	if t != nil { // t can be nil in Example functions
		t.Helper()
	}
	m.ExpectClose()
	FatalIfError(t, c.Close())
	FatalIfError(t, m.ExpectationsWereMet())
}

// FatalIfError fails the tests if an unexpected error occurred. If the error is
// gift wrapped, it prints the location. If `t` is nil, this function panics.
func FatalIfError(t testing.TB, err error) {
	if err != nil {
		if t != nil {
			t.Fatalf("%+v", err)
		} else {
			panic(err)
		}
	}
}

// CheckLastInsertID returns a function which accepts the return result from
// Exec*() and returns itself the last_insert_id or emits an error.
func CheckLastInsertID(t interface {
	Errorf(format string, args ...any)
}, msg ...string) func(sql.Result, error) int64 {
	return func(res sql.Result, err error) int64 {
		if err != nil {
			if len(msg) == 1 {
				t.Errorf("%q: %+v", msg[0], err)
			} else {
				t.Errorf("%+v", err)
			}
			return 0
		}
		lid, err := res.LastInsertId()
		if err != nil {
			if len(msg) == 1 {
				t.Errorf("%q: %+v", msg[0], err)
			} else {
				t.Errorf("%+v", err)
			}
			return 0
		}
		if lid < 0 {
			t.Errorf("Expecting Last Insert ID to be greater than zero, got %d", lid)
		}
		return lid
	}
}

// WithSQLLog displays the SQL query and its named arguments on driver level.
func WithSQLLog(buf io.Writer, printNamedArgs bool) dml.ConnPoolOption {
	return dml.WithDriverCallBack(func(fnName string) func(error, string, []driver.NamedValue) error {
		start := time.Now()
		return func(err error, query string, namedArgs []driver.NamedValue) error {
			fmt.Fprintf(buf, "%q Took: %s\n", fnName, time.Now().Sub(start))
			if err != nil {
				fmt.Fprintf(buf, "Error: %s\n", err)
			}
			if query != "" {
				fmt.Fprintf(buf, "Query: %q\n", query)
			}
			if printNamedArgs && len(namedArgs) > 0 {
				printNamedValues(buf, namedArgs)
			}
			fmt.Fprint(buf, "\n")
			return err
		}
	})
}

func printNamedValues(w io.Writer, namedArgs []driver.NamedValue) {
	fmt.Fprint(w, "NamedArgs:")
	for _, na := range namedArgs {
		fmt.Fprint(w, "{")
		if na.Name != "" {
			fmt.Fprintf(w, "Name:%q, ", na.Name)
		}
		if na.Ordinal != 0 {
			fmt.Fprintf(w, "Ordinal:%d, ", na.Ordinal)
		}
		if na.Value != nil {
			fmt.Fprintf(w, "Value:%#v", na.Value)
		}
		fmt.Fprint(w, "}, ")
	}
	fmt.Fprint(w, "\n")
}
