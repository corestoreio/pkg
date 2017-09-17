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

package dml_test

import (
	"bytes"
	"context"
	"database/sql"
	"io"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/corestoreio/csfw/sql/dml"
	"github.com/corestoreio/csfw/util/cstesting"
	"github.com/corestoreio/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var _ io.WriterTo = (*dml.RowConvert)(nil)

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
		stmt, err := dml.Exec(context.TODO(), nil, myToSQL{error: haveErr})
		assert.Nil(t, stmt)
		assert.True(t, errors.IsAlreadyClosed(err), "%+v", err)
	})
}

func TestPrepare(t *testing.T) {
	t.Parallel()
	haveErr := errors.NewAlreadyClosedf("Who closed myself?")

	t.Run("ToSQL error", func(t *testing.T) {
		stmt, err := dml.Prepare(context.TODO(), nil, myToSQL{error: haveErr})
		assert.Nil(t, stmt)
		assert.True(t, errors.IsAlreadyClosed(err), "%+v", err)
	})
	t.Run("ToSQL prepare error", func(t *testing.T) {
		dbc, dbMock := cstesting.MockDB(t)
		defer cstesting.MockClose(t, dbc, dbMock)

		dbMock.ExpectPrepare("SELECT `a` FROM `b`").WillReturnError(haveErr)

		stmt, err := dml.Prepare(context.TODO(), dbc.DB, myToSQL{sql: "SELECT `a` FROM `b`"})
		assert.Nil(t, stmt)
		assert.True(t, errors.IsAlreadyClosed(err), "%+v", err)
	})
}

type baseTest struct {
	Bool        bool
	NullBool    dml.NullBool
	Int         int
	Int64       int64
	NullInt64   dml.NullInt64
	Float64     float64
	NullFloat64 dml.NullFloat64
	Uint        uint
	Uint8       uint8
	Uint16      uint16
	Uint32      uint32
	Uint64      uint64
	Byte        []byte
	Str         string
	NullString  dml.NullString
	Time        time.Time
	NullTime    dml.NullTime
}

type baseTestCollection struct {
	Convert        dml.RowConvert
	Data           []*baseTest
	EventAfterScan func(dml.RowConvert, *baseTest)
}

func (vs *baseTestCollection) ToSQL() (string, []interface{}, error) {
	return "SELECT * FROM `test`", nil, nil
}

// RowScan implements dml.Scanner interface and scans a single row from the
// database query result.
func (vs *baseTestCollection) RowScan(r *sql.Rows) error {
	if err := vs.Convert.Scan(r); err != nil {
		return err
	}

	o := new(baseTest)
	for i, col := range vs.Convert.Columns {
		if vs.Convert.Alias != nil {
			if orgCol, ok := vs.Convert.Alias[col]; ok {
				col = orgCol
			}
		}
		b := vs.Convert.Index(i)
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
			o.Str, err = b.String()
		case "null_string":
			o.NullString, err = b.NullString()
		case "time":
			o.Time, err = b.Time()
		case "null_time":
			o.NullTime, err = b.NullTime()
		}
		if err != nil {
			return errors.Wrapf(err, "[dml_test] Failed to scan column %q at row %d", col, b.Count)
		}
	}
	if vs.EventAfterScan != nil {
		vs.EventAfterScan(vs.Convert, o)
	}
	vs.Data = append(vs.Data, o)
	return nil
}

func TestRowConvert(t *testing.T) {
	t.Parallel()

	dbc, dbMock := cstesting.MockDB(t)
	defer cstesting.MockClose(t, dbc, dbMock)

	// TODO(CyS) check that RowConvert.Byte() returns a copy

	columns := []string{
		"bool", "null_bool",
		"int", "int64", "null_int64",
		"float64", "null_float64",
		"uint", "uint8", "uint16", "uint32", "uint64",
		"byte", "str", "null_string", "time", "null_time",
	}

	t.Run("scan with error", func(t *testing.T) {
		r := sqlmock.NewRows(columns).AddRow(
			make(chan int), "Nope",
			1, 2, nil,
			0.1, nil,
			0, 8, 16, 32, 64,
			nil, "", nil, time.Time{}, nil)
		dbMock.ExpectQuery("SELECT \\* FROM `test`").WillReturnRows(r)

		tbl := new(baseTestCollection)

		rc, err := dml.Load(context.TODO(), dbc.DB, tbl, tbl)
		assert.Exactly(t, int64(0), rc)
		assert.Contains(t, err.Error(), "sql: Scan error on column index 0: unsupported Scan, storing driver.Value type chan int into type *sql.RawBytes")
	})

	t.Run("fmt.Stringer", func(t *testing.T) {

		r := sqlmock.NewRows(columns).AddRow(
			"1", "false",
			-1, -64, -128,
			0.1, 3.141,
			0, 8, 16, 32, 64,
			"byte data", "I'm a string", nil,
			now(), nil,
		)
		dbMock.ExpectQuery("SELECT \\* FROM `test`").WillReturnRows(r)

		tbl := new(baseTestCollection)
		tbl.EventAfterScan = func(b dml.RowConvert, _ *baseTest) {
			buf := new(bytes.Buffer)
			require.NoError(t, b.Debug(buf))
			assert.Exactly(t, "bool: \"1\"\nnull_bool: \"false\"\nint: \"-1\"\nint64: \"-64\"\nnull_int64: \"-128\"\nfloat64: \"0.1\"\nnull_float64: \"3.141\"\nuint: \"0\"\nuint8: \"8\"\nuint16: \"16\"\nuint32: \"32\"\nuint64: \"64\"\nbyte: \"byte data\"\nstr: \"I'm a string\"\nnull_string: <nil>\ntime: \"2006-01-02T15:04:05.000000002+00:00\"\nnull_time: <nil>", buf.String())
		}

		rc, err := dml.Load(context.TODO(), dbc.DB, tbl, tbl)
		assert.Exactly(t, int64(1), rc)
		require.NoError(t, err)
	})

	t.Run("all types non-nil", func(t *testing.T) {
		r := sqlmock.NewRows(columns).AddRow(
			"1", "false",
			-1, -64, -128,
			0.1, 3.141,
			0, 8, 16, 32, 64,
			"byte data", "I'm a string", "null_string",
			now(), now(),
		)
		dbMock.ExpectQuery("SELECT \\* FROM `test`").WillReturnRows(r)

		tbl := new(baseTestCollection)

		rc, err := dml.Load(context.TODO(), dbc.DB, tbl, tbl)
		require.NoError(t, err)

		assert.Exactly(t, int64(1), rc)
		require.Len(t, tbl.Data, 1)

		tbl.Data[0].Time = now()
		tbl.Data[0].NullTime = dml.MakeNullTime(now()) // otherwise test would fail ...
		assert.Exactly(t,
			&baseTest{
				Bool:        true,
				NullBool:    dml.MakeNullBool(false, true),
				Int:         -1,
				Int64:       -64,
				NullInt64:   dml.MakeNullInt64(-128),
				Float64:     0.1,
				NullFloat64: dml.MakeNullFloat64(3.141),
				Uint:        0x0,
				Uint8:       0x8,
				Uint16:      0x10,
				Uint32:      0x20,
				Uint64:      0x40,
				Byte:        []byte("byte data"),
				Str:         "I'm a string",
				NullString:  dml.MakeNullString("null_string"),
				Time:        now(),
				NullTime:    dml.MakeNullTime(now()),
			},
			tbl.Data[0])
	})

	t.Run("all types nil", func(t *testing.T) {
		r := sqlmock.NewRows(columns).AddRow(
			"True", nil,
			-1, -64, nil,
			0.1, nil,
			0, 8, 16, 32, 64,
			nil, "I'm a string", nil,
			now(), nil,
		)
		dbMock.ExpectQuery("SELECT \\* FROM `test`").WillReturnRows(r)

		tbl := new(baseTestCollection)

		rc, err := dml.Load(context.TODO(), dbc.DB, tbl, tbl)
		assert.Exactly(t, int64(1), rc)
		require.NoError(t, err)
		require.Len(t, tbl.Data, 1)

		want := &baseTest{
			Bool:    true,
			Int:     -1,
			Int64:   -64,
			Float64: 0.1,
			Uint:    0x0,
			Uint8:   0x8,
			Uint16:  0x10,
			Uint32:  0x20,
			Uint64:  0x40,
			Str:     "I'm a string",
			Time:    now(),
		}
		tbl.Data[0].Time = now() // useless to test this because location gets set new, hence a new pointer ...
		assert.Exactly(t, want, tbl.Data[0])

		buf := new(bytes.Buffer)
		require.NoError(t, tbl.Convert.Debug(buf))

		assert.Exactly(t, "bool: \"True\"\nnull_bool: <nil>\nint: \"-1\"\nint64: \"-64\"\nnull_int64: <nil>\nfloat64: \"0.1\"\nnull_float64: <nil>\nuint: \"0\"\nuint8: \"8\"\nuint16: \"16\"\nuint32: \"32\"\nuint64: \"64\"\nbyte: <nil>\nstr: \"I'm a string\"\nnull_string: <nil>\ntime: \"2006-01-02T15:04:05.000000002+00:00\"\nnull_time: <nil>",
			buf.String())
	})

	t.Run("invalid UTF8 Str", func(t *testing.T) {

		r := sqlmock.NewRows(columns).AddRow(
			"True", nil,
			-1, -64, nil,
			0.1, nil,
			0, 8, 16, 32, 64,
			nil, "aa\xe2", string([]byte{66, 250, 67}), // both are invalid
			now(), nil)
		dbMock.ExpectQuery("SELECT \\* FROM `test`").WillReturnRows(r)

		tbl := new(baseTestCollection)
		tbl.Convert.CheckValidUTF8 = true

		rc, err := dml.Load(context.TODO(), dbc.DB, tbl, tbl)
		assert.Exactly(t, int64(0), rc)
		assert.True(t, errors.IsNotValid(err), "%+v", err)
	})
	t.Run("invalid UTF8 NullStr", func(t *testing.T) {

		r := sqlmock.NewRows(columns).AddRow(
			"True", nil,
			-1, -64, nil,
			0.1, nil,
			0, 8, 16, 32, 64,
			nil, "aa", string([]byte{66, 250, 67}), // both are invalid
			now(), nil)

		dbMock.ExpectQuery("SELECT \\* FROM `test`").WillReturnRows(r)

		tbl := new(baseTestCollection)
		tbl.Convert.CheckValidUTF8 = true

		rc, err := dml.Load(context.TODO(), dbc.DB, tbl, tbl)
		assert.Exactly(t, int64(0), rc)
		assert.True(t, errors.IsNotValid(err), "%+v", err)

	})
	t.Run("WriteTo", func(t *testing.T) {

		r := sqlmock.NewRows(columns).AddRow(
			"True", nil,
			-1, -64, nil,
			0.1, nil,
			0, 8, 16, 32, 64,
			nil, "I'm writing to ...", nil,
			now(), now(),
		)
		dbMock.ExpectQuery("SELECT \\* FROM `test`").WillReturnRows(r)

		tbl := new(baseTestCollection)

		rc, err := dml.Load(context.TODO(), dbc.DB, tbl, tbl)
		assert.Exactly(t, int64(1), rc)
		require.NoError(t, err)

		// Does only work for one returned row OR when using a call back function
		buf := new(bytes.Buffer)
		l, err := tbl.Convert.Index(13).WriteTo(buf)
		require.NoError(t, err)
		assert.Exactly(t, int64(18), l)
		assert.Exactly(t, `I'm writing to ...`, buf.String())

	})
}
