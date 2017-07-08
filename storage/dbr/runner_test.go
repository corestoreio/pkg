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

package dbr_test

import (
	"context"
	"database/sql"
	"fmt"
	"io"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/corestoreio/csfw/storage/dbr"
	"github.com/corestoreio/csfw/util/cstesting"
	"github.com/corestoreio/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var _ fmt.Stringer = (*dbr.Base)(nil)
var _ io.WriterTo = (*dbr.Base)(nil)

type myToSQL struct {
	sql  string
	args []interface{}
	error
}

func (m myToSQL) ToSQL() (string, []interface{}, error) {
	return m.sql, m.args, m.error
}

func TestExec(t *testing.T) {
	t.Parallel()
	haveErr := errors.NewAlreadyClosedf("Who closed myself?")

	t.Run("ToSQL error", func(t *testing.T) {
		stmt, err := dbr.Exec(context.TODO(), nil, myToSQL{error: haveErr})
		assert.Nil(t, stmt)
		assert.True(t, errors.IsAlreadyClosed(err), "%+v", err)
	})
}

func TestPrepare(t *testing.T) {
	t.Parallel()
	haveErr := errors.NewAlreadyClosedf("Who closed myself?")

	t.Run("ToSQL error", func(t *testing.T) {
		stmt, err := dbr.Prepare(context.TODO(), nil, myToSQL{error: haveErr})
		assert.Nil(t, stmt)
		assert.True(t, errors.IsAlreadyClosed(err), "%+v", err)
	})
	t.Run("ToSQL prepare error", func(t *testing.T) {
		dbc, dbMock := cstesting.MockDB(t)
		defer cstesting.MockClose(t, dbc, dbMock)

		dbMock.ExpectPrepare("SELECT `a` FROM `b`").WillReturnError(haveErr)

		stmt, err := dbr.Prepare(context.TODO(), dbc.DB, myToSQL{sql: "SELECT `a` FROM `b`"})
		assert.Nil(t, stmt)
		assert.True(t, errors.IsAlreadyClosed(err), "%+v", err)
	})
}

type baseTest struct {
	Bool        bool
	NullBool    sql.NullBool
	Int         int
	Int64       int64
	NullInt64   sql.NullInt64
	Float64     float64
	NullFloat64 sql.NullFloat64
	Uint        uint
	Uint8       uint8
	Uint16      uint16
	Uint32      uint32
	Uint64      uint64
	Byte        []byte
	Str         string
	NullString  sql.NullString
}

type baseTestCollection struct {
	Base           dbr.Base
	Data           []*baseTest
	EventAfterScan func(dbr.Base, *baseTest)
}

func (vs *baseTestCollection) ToSQL() (string, []interface{}, error) {
	return "SELECT * FROM `test`", nil, nil
}

// RowScan implements dbr.Scanner interface and scans a single row from the
// database query result.
func (vs *baseTestCollection) RowScan(r *sql.Rows) error {
	if err := vs.Base.Scan(r); err != nil {
		return err
	}

	o := new(baseTest)
	for i, col := range vs.Base.Columns {
		if vs.Base.Alias != nil {
			if orgCol, ok := vs.Base.Alias[col]; ok {
				col = orgCol
			}
		}
		b := vs.Base.Index(i)
		var err error

		switch col {
		case "bool":
			o.Bool, err = b.Bool()
		case "null_bool":
			o.NullBool, err = b.NullBool()
		case "int":
			o.Int, err = b.Int()
		case "int64":
			o.Int64, err = b.Int64()
		case "null_int64":
			o.NullInt64, err = b.NullInt64()
		case "float64":
			o.Float64, err = b.Float64()
		case "null_float64":
			o.NullFloat64, err = b.NullFloat64()
		case "uint":
			o.Uint, err = b.Uint()
		case "uint8":
			o.Uint8, err = b.Uint8()
		case "uint16":
			o.Uint16, err = b.Uint16()
		case "uint32":
			o.Uint32, err = b.Uint32()
		case "uint64":
			o.Uint64, err = b.Uint64()
		case "byte":
			o.Byte = b.Byte()
		case "str":
			o.Str, err = b.Str()
		case "null_string":
			o.NullString, err = b.NullString()
		}
		if err != nil {
			return errors.Wrapf(err, "[dbr_test] Failed to scan %q at row %d", col, b.Count)
		}
	}
	if vs.EventAfterScan != nil {
		vs.EventAfterScan(vs.Base, o)
	}
	vs.Data = append(vs.Data, o)
	return nil
}

func TestBase(t *testing.T) {
	t.Parallel()
	dbc, dbMock := cstesting.MockDB(t)
	defer cstesting.MockClose(t, dbc, dbMock)

	columns := []string{
		"bool", "null_bool",
		"int", "int64", "null_int64",
		"float64", "null_float64",
		"uint", "uint8", "uint16", "uint32", "uint64",
		"byte", "str", "null_string",
	}

	t.Run("scan with error", func(t *testing.T) {
		r := sqlmock.NewRows(columns).AddRow(
			make(chan int), "Nope",
			1, 2, nil,
			0.1, nil,
			0, 8, 16, 32, 64,
			nil, "", nil)
		dbMock.ExpectQuery("SELECT \\* FROM `test`").WillReturnRows(r)

		tbl := new(baseTestCollection)

		rc, err := dbr.Load(context.TODO(), dbc.DB, tbl, tbl)
		assert.Exactly(t, int64(0), rc)
		assert.EqualError(t, err, "sql: Scan error on column index 0: unsupported Scan, storing driver.Value type chan int into type *sql.RawBytes")
	})

	t.Run("fmt.Stringer", func(t *testing.T) {

		r := sqlmock.NewRows(columns).AddRow(
			"1", "false",
			-1, -64, -128,
			0.1, 3.141,
			0, 8, 16, 32, 64,
			"byte data", "I'm a string", nil)
		dbMock.ExpectQuery("SELECT \\* FROM `test`").WillReturnRows(r)

		tbl := new(baseTestCollection)
		tbl.EventAfterScan = func(b dbr.Base, _ *baseTest) {
			assert.Exactly(t, `bool: "1"
null_bool: "false"
int: "-1"
int64: "-64"
null_int64: "-128"
float64: "0.1"
null_float64: "3.141"
uint: "0"
uint8: "8"
uint16: "16"
uint32: "32"
uint64: "64"
byte: "byte data"
str: "I'm a string"
null_string: <nil>`, b.String())
		}

		rc, err := dbr.Load(context.TODO(), dbc.DB, tbl, tbl)
		assert.Exactly(t, int64(1), rc)
		require.NoError(t, err)
	})
}
