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
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/corestoreio/csfw/sql/dml"
	"github.com/corestoreio/csfw/util/cstesting"
	"github.com/corestoreio/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

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

var (
	_ dml.ColumnMapper = (*baseTest)(nil)
	_ dml.ColumnMapper = (*baseTestCollection)(nil)
)

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
	Decimal     dml.Decimal
}

func (bt *baseTest) MapColumns(cm *dml.ColumnMap) error {
	if cm.Mode() == dml.ColumnMapEntityReadAll {
		return cm.Bool(&bt.Bool).NullBool(&bt.NullBool).Int(&bt.Int).Int64(&bt.Int64).NullInt64(&bt.NullInt64).Float64(&bt.Float64).NullFloat64(&bt.NullFloat64).Uint(&bt.Uint).Uint8(&bt.Uint8).Uint16(&bt.Uint16).Uint32(&bt.Uint32).Uint64(&bt.Uint64).Byte(&bt.Byte).String(&bt.Str).NullString(&bt.NullString).Time(&bt.Time).NullTime(&bt.NullTime).Err()
	}
	for cm.Next() {
		switch c := cm.Column(); c {
		case "bool":
			cm.Bool(&bt.Bool)
		case "null_bool":
			cm.NullBool(&bt.NullBool)
		case "int":
			cm.Int(&bt.Int)
		case "int64":
			cm.Int64(&bt.Int64)
		case "null_int64":
			cm.NullInt64(&bt.NullInt64)
		case "float64":
			cm.Float64(&bt.Float64)
		case "null_float64":
			cm.NullFloat64(&bt.NullFloat64)
		case "uint":
			cm.Uint(&bt.Uint)
		case "uint8":
			cm.Uint8(&bt.Uint8)
		case "uint16":
			cm.Uint16(&bt.Uint16)
		case "uint32":
			cm.Uint32(&bt.Uint32)
		case "uint64":
			cm.Uint64(&bt.Uint64)
		case "byte":
			cm.Byte(&bt.Byte)
		case "str":
			cm.String(&bt.Str)
		case "null_string":
			cm.NullString(&bt.NullString)
		case "time":
			cm.Time(&bt.Time)
		case "null_time":
			cm.NullTime(&bt.NullTime)
		case "decimal":
			cm.Decimal(&bt.Decimal)
		default:
			return errors.NewNotFoundf("[dml_test] dmlPerson Column %q not found", c)
		}
	}
	return cm.Err()
}

type baseTestCollection struct {
	Data           []*baseTest
	EventAfterScan func(*dml.ColumnMap, *baseTest)
	CheckValidUTF8 bool
}

func (vs *baseTestCollection) ToSQL() (string, []interface{}, error) {
	return "SELECT * FROM `test`", nil, nil
}

// RowScan implements dml.ColumnMapper interface and scans a single row from the
// database query result.
func (vs *baseTestCollection) MapColumns(cm *dml.ColumnMap) error {
	switch m := cm.Mode(); m {
	case dml.ColumnMapEntityReadAll: // INSERT STATEMENT requesting all columns aka arguments
		for _, p := range vs.Data {
			if err := p.MapColumns(cm); err != nil {
				return errors.WithStack(err)
			}
		}
	case dml.ColumnMapScan:
		// case for scanning when loading certain rows, hence we write data from
		// the DB into the struct in each for-loop.
		if cm.Count == 0 {
			vs.Data = vs.Data[:0]
			cm.CheckValidUTF8 = vs.CheckValidUTF8
		}
		p := new(baseTest)
		if err := p.MapColumns(cm); err != nil {
			return errors.WithStack(err)
		}
		if vs.EventAfterScan != nil {
			vs.EventAfterScan(cm, p)
		}
		vs.Data = append(vs.Data, p)
	case 'r':
		panic("not needed")
	default:
		return errors.NewNotSupportedf("[dml] Unknown Mode: %q", string(m))
	}
	return cm.Err()
}

func TestRowConvert(t *testing.T) {
	t.Parallel()

	dbc, dbMock := cstesting.MockDB(t)
	defer cstesting.MockClose(t, dbc, dbMock)

	// TODO(CyS) check that RowMap.Byte() returns a copy

	columns := []string{
		"bool", "null_bool",
		"int", "int64", "null_int64",
		"float64", "null_float64",
		"uint", "uint8", "uint16", "uint32", "uint64",
		"byte", "str", "null_string", "time", "null_time",
		"decimal",
	}

	t.Run("scan with error", func(t *testing.T) {
		r := sqlmock.NewRows(columns).AddRow(
			make(chan int), "Nope",
			1, 2, nil,
			0.1, nil,
			0, 8, 16, 32, 64,
			nil, "", nil, time.Time{}, nil,
			nil)
		dbMock.ExpectQuery("SELECT \\* FROM `test`").WillReturnRows(r)

		tbl := new(baseTestCollection)

		rc, err := dml.Load(context.TODO(), dbc.DB, tbl, tbl)
		assert.Exactly(t, uint64(0), rc)
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
			nil,
		)
		dbMock.ExpectQuery("SELECT \\* FROM `test`").WillReturnRows(r)

		tbl := new(baseTestCollection)
		tbl.EventAfterScan = func(b *dml.ColumnMap, _ *baseTest) {
			buf := new(bytes.Buffer)
			require.NoError(t, b.Debug(buf))
			assert.Exactly(t, "bool: \"1\"\nnull_bool: \"false\"\nint: \"-1\"\nint64: \"-64\"\nnull_int64: \"-128\"\nfloat64: \"0.1\"\nnull_float64: \"3.141\"\nuint: \"0\"\nuint8: \"8\"\nuint16: \"16\"\nuint32: \"32\"\nuint64: \"64\"\nbyte: \"byte data\"\nstr: \"I'm a string\"\nnull_string: <nil>\ntime: \"2006-01-02T15:04:05.000000002+00:00\"\nnull_time: <nil>\ndecimal: <nil>", buf.String())
		}

		rc, err := dml.Load(context.TODO(), dbc.DB, tbl, tbl)
		assert.Exactly(t, uint64(1), rc)
		require.NoError(t, err)
	})

	t.Run("all types non-nil", func(t *testing.T) {
		r := sqlmock.NewRows(columns).AddRow(
			"1", "false",
			-1, -64, -128,
			0.1, 3.141,
			0, 8, 16, 32, 64,
			"byte data", "I'm a string", "null_string",
			now(), now(), "2681.7000",
		)
		dbMock.ExpectQuery("SELECT \\* FROM `test`").WillReturnRows(r)

		tbl := new(baseTestCollection)

		rc, err := dml.Load(context.TODO(), dbc.DB, tbl, tbl)
		require.NoError(t, err)

		assert.Exactly(t, uint64(1), rc)
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
				Decimal: dml.Decimal{
					Precision: 26817000,
					Scale:     4,
					Valid:     true,
				},
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
			now(), nil, nil,
		)
		dbMock.ExpectQuery("SELECT \\* FROM `test`").WillReturnRows(r)

		tbl := new(baseTestCollection)

		rc, err := dml.Load(context.TODO(), dbc.DB, tbl, tbl)
		assert.Exactly(t, uint64(1), rc)
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
	})

	t.Run("invalid UTF8 Str", func(t *testing.T) {

		r := sqlmock.NewRows(columns).AddRow(
			"True", nil,
			-1, -64, nil,
			0.1, nil,
			0, 8, 16, 32, 64,
			nil, "aa\xe2", string([]byte{66, 250, 67}), // both are invalid
			now(), nil, nil)
		dbMock.ExpectQuery("SELECT \\* FROM `test`").WillReturnRows(r)

		tbl := new(baseTestCollection)
		tbl.CheckValidUTF8 = true

		rc, err := dml.Load(context.TODO(), dbc.DB, tbl, tbl)
		assert.Exactly(t, uint64(0), rc)
		assert.True(t, errors.IsNotValid(err), "%+v", err)
	})
	t.Run("invalid UTF8 NullStr", func(t *testing.T) {

		r := sqlmock.NewRows(columns).AddRow(
			"True", nil,
			-1, -64, nil,
			0.1, nil,
			0, 8, 16, 32, 64,
			nil, "aa", string([]byte{66, 250, 67}), // both are invalid
			now(), nil, nil)

		dbMock.ExpectQuery("SELECT \\* FROM `test`").WillReturnRows(r)

		tbl := new(baseTestCollection)
		tbl.CheckValidUTF8 = true

		rc, err := dml.Load(context.TODO(), dbc.DB, tbl, tbl)
		assert.Exactly(t, uint64(0), rc)
		assert.True(t, errors.IsNotValid(err), "%+v", err)

	})
}
