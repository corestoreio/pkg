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

package dml_test

import (
	"bytes"
	"context"
	"database/sql/driver"
	"encoding"
	"encoding/json"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/corestoreio/errors"
	"github.com/corestoreio/pkg/sql/dml"
	"github.com/corestoreio/pkg/sql/dmltest"
	"github.com/corestoreio/pkg/storage/null"
	"github.com/corestoreio/pkg/util/assert"
)

var (
	_ dml.ColumnMapper           = (*baseTest)(nil)
	_ dml.ColumnMapper           = (*baseTestCollection)(nil)
	_ encoding.TextMarshaler     = (*textBinaryEncoder)(nil)
	_ encoding.TextUnmarshaler   = (*textBinaryEncoder)(nil)
	_ encoding.BinaryMarshaler   = (*textBinaryEncoder)(nil)
	_ encoding.BinaryUnmarshaler = (*textBinaryEncoder)(nil)
)

func TestColumnMap_BinaryText(t *testing.T) {
	cm := dml.NewColumnMap(1)

	assert.NoError(t, cm.Binary(&textBinaryEncoder{data: []byte(`BinaryTest`)}).Err())
	assert.Exactly(t, "dml.MakeArgs(1).Bytes([]byte{0x42, 0x69, 0x6e, 0x61, 0x72, 0x79, 0x54, 0x65, 0x73, 0x74})", cm.GoString())
	assert.NoError(t, cm.Text(&textBinaryEncoder{data: []byte(`TextTest`)}).Err())
	assert.Exactly(t, "dml.MakeArgs(2).Bytes([]byte{0x42, 0x69, 0x6e, 0x61, 0x72, 0x79, 0x54, 0x65, 0x73, 0x74}).Bytes([]byte{0x54, 0x65, 0x78, 0x74, 0x54, 0x65, 0x73, 0x74})", cm.GoString())

	cm.CheckValidUTF8 = true
	err := cm.Text(&textBinaryEncoder{data: []byte("\xc0\x80")}).Err()
	assert.True(t, errors.Is(err, errors.NotValid), "Want errors.NotValid; Got %s\n%+v", errors.UnwrapKind(err), err)

}

type textBinaryEncoder struct {
	data []byte
}

func (be textBinaryEncoder) MarshalBinary() (data []byte, err error) {
	if bytes.Equal(be.data, []byte(`error`)) {
		return nil, errors.DecryptionFailed.Newf("decryption failed test error")
	}
	return be.data, nil
}
func (be *textBinaryEncoder) UnmarshalBinary(data []byte) error {
	if bytes.Equal(data, []byte(`error`)) {
		return errors.Empty.Newf("test error empty")
	}
	be.data = append(be.data, data...)
	return nil
}
func (be textBinaryEncoder) MarshalText() (text []byte, err error) {
	if bytes.Equal(be.data, []byte(`error`)) {
		return nil, errors.DecryptionFailed.Newf("internal validation failed test error")
	}
	return be.data, nil
}
func (be *textBinaryEncoder) UnmarshalText(text []byte) error {
	if bytes.Equal(text, []byte(`error`)) {
		return errors.Empty.Newf("test error empty")
	}
	be.data = append(be.data, text...)
	return nil
}

type baseTest struct {
	Bool        bool
	NullBool    null.Bool
	Int         int
	Int8        int8
	Int16       int16
	Int32       int32
	Int64       int64
	NullInt64   null.Int64
	Float64     float64
	NullFloat64 null.Float64
	Uint        uint
	Uint8       uint8
	Uint16      uint16
	Uint32      uint32
	Uint64      uint64
	Byte        []byte
	Str         string
	NullString  null.String
	Time        time.Time
	NullTime    null.Time
	Decimal     null.Decimal
	Text        textBinaryEncoder
	Binary      textBinaryEncoder
}

// TODO add null types of the u/int
func (bt *baseTest) MapColumns(cm *dml.ColumnMap) error {
	if cm.Mode() == dml.ColumnMapEntityReadAll {
		return cm.Bool(&bt.Bool).NullBool(&bt.NullBool).
			Int(&bt.Int).Int8(&bt.Int8).Int16(&bt.Int16).Int32(&bt.Int32).Int64(&bt.Int64).NullInt64(&bt.NullInt64).
			Float64(&bt.Float64).NullFloat64(&bt.NullFloat64).
			Uint(&bt.Uint).Uint8(&bt.Uint8).Uint16(&bt.Uint16).Uint32(&bt.Uint32).Uint64(&bt.Uint64).
			Byte(&bt.Byte).String(&bt.Str).NullString(&bt.NullString).
			Time(&bt.Time).NullTime(&bt.NullTime).Err()
	}
	for cm.Next() {
		switch c := cm.Column(); c {
		case "bool":
			cm.Bool(&bt.Bool)
		case "null_bool":
			cm.NullBool(&bt.NullBool)
		case "int":
			cm.Int(&bt.Int)
		case "int8":
			cm.Int8(&bt.Int8)
		case "int16":
			cm.Int16(&bt.Int16)
		case "int32":
			cm.Int32(&bt.Int32)
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
		case "text":
			cm.Text(&bt.Text)
		case "binary":
			cm.Binary(&bt.Binary)
		default:
			return errors.NotFound.Newf("[dml_test] dmlPerson Column %q not found", c)
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
		return errors.NotSupported.Newf("[dml] Unknown Mode: %q", string(m))
	}
	return cm.Err()
}

func TestColumnMap_Query(t *testing.T) {
	t.Parallel()

	dbc, dbMock := dmltest.MockDB(t)
	defer dmltest.MockClose(t, dbc, dbMock)

	// TODO(CyS) check that RowMap.Byte() returns a copy

	columns := []string{
		"bool", "null_bool",
		"int", "int64", "null_int64",
		"float64", "null_float64",
		"uint", "uint8", "uint16", "uint32", "uint64",
		"byte", "str", "null_string", "time", "null_time",
		"decimal",
		"text", "binary",
	}

	t.Run("scan with error", func(t *testing.T) {
		r := sqlmock.NewRows(columns).AddRow(
			1234, "Nope",
			1, 2, nil,
			0.1, nil,
			0, 8, 16, 32, 64,
			nil, "", nil, time.Time{}, nil,
			nil, nil, nil)
		dbMock.ExpectQuery("SELECT \\* FROM `test`").WillReturnRows(r)

		tbl := new(baseTestCollection)

		rc, err := dbc.WithQueryBuilder(tbl).Load(context.TODO(), tbl)
		assert.Exactly(t, uint64(0), rc)
		assert.Contains(t, err.Error(), `[dml] Column "null_bool": strconv.ParseBool: parsing "Nope": invalid syntax`)
	})

	t.Run("fmt.Stringer", func(t *testing.T) {
		// TODO extend the rows to add all types for `baseTest`
		r := sqlmock.NewRows(columns).AddRow(
			"1", "false",
			-13, int64(-64), -128,
			0.1, 3.141,
			0, 8, 16, 32, 64,
			"byte data", "I'm a string", nil,
			now(), nil,
			nil, nil, nil,
		)
		dbMock.ExpectQuery("SELECT \\* FROM `test`").WillReturnRows(r)

		tbl := new(baseTestCollection)
		tbl.EventAfterScan = func(b *dml.ColumnMap, _ *baseTest) {
			buf := new(bytes.Buffer)
			assert.NoError(t, b.Debug(buf))
			assert.Exactly(t, "bool: \"1\"\nnull_bool: \"false\"\nint: \"-13\"\nint64: \"-64\"\nnull_int64: \"-128\"\nfloat64: \"0.1\"\nnull_float64: \"3.141\"\nuint: \"0\"\nuint8: \"8\"\nuint16: \"16\"\nuint32: \"32\"\nuint64: \"64\"\nbyte: \"byte data\"\nstr: \"I'm a string\"\nnull_string: <nil>\ntime: \"2006-01-02 15:04:05.000000002 +0000 hardcoded\"\nnull_time: <nil>\ndecimal: <nil>\ntext: <nil>\nbinary: <nil>",
				buf.String())
		}

		rc, err := dbc.WithQueryBuilder(tbl).Load(context.TODO(), tbl)
		assert.Exactly(t, uint64(1), rc)
		assert.NoError(t, err)
	})

	t.Run("all types non-nil", func(t *testing.T) {
		r := sqlmock.NewRows(columns).AddRow(
			"1", "false",
			-1, -64, -128,
			0.1, 3.141,
			0, 8, 16, 32, 64,
			"byte data", "I'm a string", "null_string",
			now(), now(), "2681.7000", []byte(`Hello Text`), []byte(`Hello Binary`),
		)
		dbMock.ExpectQuery("SELECT \\* FROM `test`").WillReturnRows(r)

		tbl := new(baseTestCollection)

		rc, err := dbc.WithQueryBuilder(tbl).Load(context.TODO(), tbl)
		assert.NoError(t, err)

		assert.Exactly(t, uint64(1), rc)
		assert.Len(t, tbl.Data, 1)

		tbl.Data[0].Time = now()
		tbl.Data[0].NullTime = null.MakeTime(now()) // otherwise test would fail ...

		bt := &baseTest{
			Bool:        true,
			NullBool:    null.MakeBool(false),
			Int:         -1,
			Int64:       -64,
			NullInt64:   null.MakeInt64(-128),
			Float64:     0.1,
			NullFloat64: null.MakeFloat64(3.141),
			Uint:        0x0,
			Uint8:       0x8,
			Uint16:      0x10,
			Uint32:      0x20,
			Uint64:      0x40,
			Byte:        []byte("byte data"),
			Str:         "I'm a string",
			NullString:  null.MakeString("null_string"),
			Time:        now(),
			NullTime:    null.MakeTime(now()),
			Decimal: null.Decimal{
				Precision: 26817000,
				Scale:     4,
				Valid:     true,
			},
			Text:   textBinaryEncoder{data: []byte(`Hello Text`)},
			Binary: textBinaryEncoder{data: []byte(`Hello Binary`)},
		}

		btj, err := json.Marshal(bt)
		assert.NoError(t, err)
		tdj, err := json.Marshal(tbl.Data[0])
		assert.NoError(t, err)

		assert.Exactly(t, btj, tdj, "\nWant: %q\nHave: %q", btj, tdj)
	})

	t.Run("all types nil", func(t *testing.T) {
		r := sqlmock.NewRows(columns).AddRow(
			"True", nil,
			-1, -64, nil,
			0.1, nil,
			0, 8, 16, 32, 64,
			nil, "I'm a string", nil,
			now(), nil, nil,
			nil, nil,
		)
		dbMock.ExpectQuery("SELECT \\* FROM `test`").WillReturnRows(r)

		tbl := new(baseTestCollection)

		rc, err := dbc.WithQueryBuilder(tbl).Load(context.TODO(), tbl)
		assert.Exactly(t, uint64(1), rc)
		assert.NoError(t, err)
		assert.Len(t, tbl.Data, 1)

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
			now(), nil, nil,
			nil, nil)
		dbMock.ExpectQuery("SELECT \\* FROM `test`").WillReturnRows(r)

		tbl := new(baseTestCollection)
		tbl.CheckValidUTF8 = true

		rc, err := dbc.WithQueryBuilder(tbl).Load(context.TODO(), tbl)
		assert.Exactly(t, uint64(0), rc)
		assert.True(t, errors.NotValid.Match(err), "%+v", err)
	})

	t.Run("invalid UTF8 NullStr", func(t *testing.T) {

		r := sqlmock.NewRows(columns).AddRow(
			"True", nil,
			-1, -64, nil,
			0.1, nil,
			0, 8, 16, 32, 64,
			nil, "aa", string([]byte{66, 250, 67}), // both are invalid
			now(), nil, nil,
			nil, nil)

		dbMock.ExpectQuery("SELECT \\* FROM `test`").WillReturnRows(r)

		tbl := new(baseTestCollection)
		tbl.CheckValidUTF8 = true

		rc, err := dbc.WithQueryBuilder(tbl).Load(context.TODO(), tbl)
		assert.Exactly(t, uint64(0), rc)
		assert.True(t, errors.NotValid.Match(err), "%+v", err)

	})
}

func TestColumnMap_Prepared(t *testing.T) {
	t.Parallel()

	dbc, dbMock := dmltest.MockDB(t)
	defer dmltest.MockClose(t, dbc, dbMock)

	columns := []string{
		"bool", "null_bool",
		"int", "int64", "null_int64",
		"float64", "null_float64",
		"uint", "uint8", "uint16", "uint32", "uint64",
		"byte", "str", "null_string", "time", "null_time",
		"decimal",
		"text", "binary",
	}

	runner := func(scanErrWantKind errors.Kind, want string, values ...driver.Value) func(*testing.T) {
		return func(t *testing.T) {
			r := sqlmock.NewRows(columns).AddRow(values...)
			dbMock.ExpectPrepare("SELECT \\* FROM `test`").ExpectQuery().WillReturnRows(r)
			tbl := new(baseTestCollection)
			tbl.CheckValidUTF8 = true
			tbl.EventAfterScan = func(b *dml.ColumnMap, _ *baseTest) {
				buf := new(bytes.Buffer)
				assert.NoError(t, b.Debug(buf))
				assert.Exactly(t, want, buf.String())
			}
			stmt, err := dbc.SelectFrom("test").Star().Prepare(context.TODO())
			assert.NoError(t, err)

			rc, err := stmt.WithArgs().Load(context.TODO(), tbl)
			if scanErrWantKind != errors.NoKind {
				assert.True(t, errors.Is(err, scanErrWantKind), "Should be Error Kind %s; Got: %s\n%+v", scanErrWantKind, errors.UnwrapKind(err), err)
			} else {
				assert.NoError(t, err)
				assert.Exactly(t, uint64(1), rc, "Should return one loaded row")
			}
		}
	}
	t.Run("native bool", runner(
		errors.NoKind,
		"bool: \"true\"\nnull_bool: \"false\"\nint: \"1\"\nint64: \"2\"\nnull_int64: <nil>\nfloat64: \"0.1\"\nnull_float64: <nil>\nuint: \"0\"\nuint8: \"8\"\nuint16: \"16\"\nuint32: \"32\"\nuint64: \"64\"\nbyte: <nil>\nstr: \"\"\nnull_string: <nil>\ntime: \"0001-01-01 00:00:00 +0000 UTC\"\nnull_time: <nil>\ndecimal: <nil>\ntext: <nil>\nbinary: <nil>",
		true, false, // "bool", "null_bool"
		1, 2, nil, // "int", "int64", "null_int64"
		0.1, nil, // "float64", "null_float64",
		0, 8, 16, 32, 64, // "uint", "uint8", "uint16", "uint32", "uint64"
		nil, "", nil, time.Time{}, nil, // "byte", "str", "null_string", "time", "null_time",
		nil,      //	"decimal",
		nil, nil, // "text", "binary",
	))
	t.Run("int bool", runner(
		errors.NoKind,
		"bool: \"1\"\nnull_bool: \"1\"\nint: \"1\"\nint64: \"2\"\nnull_int64: <nil>\nfloat64: \"0.1\"\nnull_float64: <nil>\nuint: \"0\"\nuint8: \"8\"\nuint16: \"16\"\nuint32: \"32\"\nuint64: \"64\"\nbyte: <nil>\nstr: \"\"\nnull_string: <nil>\ntime: \"0001-01-01 00:00:00 +0000 UTC\"\nnull_time: <nil>\ndecimal: <nil>\ntext: <nil>\nbinary: <nil>",
		1, 1,
		1, 2, nil,
		0.1, nil,
		0, 8, 16, 32, 64,
		nil, "", nil, time.Time{}, nil,
		nil,
		nil, nil, // "text", "binary",
	))
	t.Run("float64 bool error", runner(
		errors.NotSupported,
		"",
		float64(1), 1,
		1, 2, nil,
		0.1, nil,
		0, 8, 16, 32, 64,
		nil, "", nil, time.Time{}, nil,
		nil,
		nil, nil, // "text", "binary",
	))
	t.Run("string bool error", runner(
		errors.BadEncoding,
		"",
		"ok", 1,
		1, 2, nil,
		0.1, nil,
		0, 8, 16, 32, 64,
		nil, "", nil, time.Time{}, nil,
		nil,
		nil, nil, // "text", "binary",
	))
	t.Run("null bool empty string", runner(
		errors.NoKind,
		"bool: \"1\"\nnull_bool: \"\"\nint: \"1\"\nint64: \"2\"\nnull_int64: <nil>\nfloat64: \"0.1\"\nnull_float64: <nil>\nuint: \"0\"\nuint8: \"8\"\nuint16: \"16\"\nuint32: \"32\"\nuint64: \"64\"\nbyte: <nil>\nstr: \"\"\nnull_string: <nil>\ntime: \"0001-01-01 00:00:00 +0000 UTC\"\nnull_time: <nil>\ndecimal: <nil>\ntext: <nil>\nbinary: <nil>",
		1, "",
		1, 2, nil,
		0.1, nil,
		0, 8, 16, 32, 64,
		nil, "", nil, time.Time{}, nil,
		nil,
		nil, nil, // "text", "binary",
	))
	t.Run("null bool float64 error", runner(
		errors.NotSupported,
		"",
		1, float64(1),
		1, 2, nil,
		0.1, nil,
		0, 8, 16, 32, 64,
		nil, "", nil, time.Time{}, nil,
		nil,
		nil, nil, // "text", "binary",
	))
	t.Run("null bool byte error", runner(
		errors.BadEncoding,
		"",
		1, []byte(`ok`),
		1, 2, nil,
		0.1, nil,
		0, 8, 16, 32, 64,
		nil, "", nil, time.Time{}, nil,
		nil,
		nil, nil, // "text", "binary",
	))
	t.Run("null bool string error", runner(
		errors.BadEncoding,
		"",
		1, `ok`,
		1, 2, nil,
		0.1, nil,
		0, 8, 16, 32, 64,
		nil, "", nil, time.Time{}, nil,
		nil,
		nil, nil, // "text", "binary",
	))
	t.Run("int string error", runner(
		errors.NotSupported,
		"",
		true, false, // "bool", "null_bool"
		"1", 2, nil, // "int", "int64", "null_int64"
		0.1, nil, // "float64", "null_float64",
		0, 8, 16, 32, 64, // "uint", "uint8", "uint16", "uint32", "uint64"
		nil, "", nil, time.Time{}, nil, // "byte", "str", "null_string", "time", "null_time",
		nil,      //	"decimal",
		nil, nil, // "text", "binary",
	))
	t.Run("int byte error", runner(
		errors.BadEncoding,
		"",
		true, false, // "bool", "null_bool"
		[]byte("1.0"), 2, nil, // "int", "int64", "null_int64"
		0.1, nil, // "float64", "null_float64",
		0, 8, 16, 32, 64, // "uint", "uint8", "uint16", "uint32", "uint64"
		nil, "", nil, time.Time{}, nil, // "byte", "str", "null_string", "time", "null_time",
		nil,      //	"decimal",
		nil, nil, // "text", "binary",
	))

	t.Run("int64 string error", runner(
		errors.NotSupported,
		"",
		true, false, // "bool", "null_bool"
		1, "2", nil, // "int", "int64", "null_int64"
		0.1, nil, // "float64", "null_float64",
		0, 8, 16, 32, 64, // "uint", "uint8", "uint16", "uint32", "uint64"
		nil, "", nil, time.Time{}, nil, // "byte", "str", "null_string", "time", "null_time",
		nil,      //	"decimal",
		nil, nil, // "text", "binary",
	))
	t.Run("int64 byte error", runner(
		errors.BadEncoding,
		"",
		true, false, // "bool", "null_bool"
		1, []byte("2.0"), nil, // "int", "int64", "null_int64"
		0.1, nil, // "float64", "null_float64",
		0, 8, 16, 32, 64, // "uint", "uint8", "uint16", "uint32", "uint64"
		nil, "", nil, time.Time{}, nil, // "byte", "str", "null_string", "time", "null_time",
		nil,      //	"decimal",
		nil, nil, // "text", "binary",
	))
	t.Run("null int64 byte error", runner(
		errors.BadEncoding,
		"",
		true, false, // "bool", "null_bool"
		1, 2, []byte("2.0"), // "int", "int64", "null_int64"
		0.1, nil, // "float64", "null_float64",
		0, 8, 16, 32, 64, // "uint", "uint8", "uint16", "uint32", "uint64"
		nil, "", nil, time.Time{}, nil, // "byte", "str", "null_string", "time", "null_time",
		nil,      //	"decimal",
		nil, nil, // "text", "binary",
	))
	t.Run("null int64 float error", runner(
		errors.NotSupported,
		"",
		true, false, // "bool", "null_bool"
		1, 2, 2.0, // "int", "int64", "null_int64"
		0.1, nil, // "float64", "null_float64",
		0, 8, 16, 32, 64, // "uint", "uint8", "uint16", "uint32", "uint64"
		nil, "", nil, time.Time{}, nil, // "byte", "str", "null_string", "time", "null_time",
		nil,      //	"decimal",
		nil, nil, // "text", "binary",
	))
	t.Run("float64 as int error", runner(
		errors.NotSupported,
		"",
		true, false, // "bool", "null_bool"
		1, 2, 3, // "int", "int64", "null_int64"
		64, nil, // "float64", "null_float64",
		0, 8, 16, 32, 64, // "uint", "uint8", "uint16", "uint32", "uint64"
		nil, "", nil, time.Time{}, nil, // "byte", "str", "null_string", "time", "null_time",
		nil,      //	"decimal",
		nil, nil, // "text", "binary",
	))
	t.Run("float64 as byte error", runner(
		errors.BadEncoding,
		"",
		true, false, // "bool", "null_bool"
		1, 2, 3, // "int", "int64", "null_int64"
		[]byte(`6,4`), nil, // "float64", "null_float64",
		0, 8, 16, 32, 64, // "uint", "uint8", "uint16", "uint32", "uint64"
		nil, "", nil, time.Time{}, nil, // "byte", "str", "null_string", "time", "null_time",
		nil,      //	"decimal",
		nil, nil, // "text", "binary",
	))
	t.Run("null float64 as byte error", runner(
		errors.BadEncoding,
		"",
		true, false, // "bool", "null_bool"
		1, 2, 3, // "int", "int64", "null_int64"
		6.4, []byte(`6,4`), // "float64", "null_float64",
		0, 8, 16, 32, 64, // "uint", "uint8", "uint16", "uint32", "uint64"
		nil, "", nil, time.Time{}, nil, // "byte", "str", "null_string", "time", "null_time",
		nil,      //	"decimal",
		nil, nil, // "text", "binary",
	))
	t.Run("null float64 as int error", runner(
		errors.NotSupported,
		"",
		true, false, // "bool", "null_bool"
		1, 2, 3, // "int", "int64", "null_int64"
		6.4, 64, // "float64", "null_float64",
		0, 8, 16, 32, 64, // "uint", "uint8", "uint16", "uint32", "uint64"
		nil, "", nil, time.Time{}, nil, // "byte", "str", "null_string", "time", "null_time",
		nil,      //	"decimal",
		nil, nil, // "text", "binary",
	))
	t.Run("decimal no error", runner(
		errors.NoKind,
		"bool: \"true\"\nnull_bool: \"false\"\nint: \"1\"\nint64: \"2\"\nnull_int64: \"3\"\nfloat64: \"6.4\"\nnull_float64: <nil>\nuint: \"0\"\nuint8: \"8\"\nuint16: \"16\"\nuint32: \"32\"\nuint64: \"64\"\nbyte: <nil>\nstr: \"\"\nnull_string: <nil>\ntime: \"0001-01-01 00:00:00 +0000 UTC\"\nnull_time: <nil>\ndecimal: \"48.98\"\ntext: <nil>\nbinary: <nil>",
		true, false, // "bool", "null_bool"
		1, 2, 3, // "int", "int64", "null_int64"
		6.4, nil, // "float64", "null_float64",
		0, 8, 16, 32, 64, // "uint", "uint8", "uint16", "uint32", "uint64"
		nil, "", nil, time.Time{}, nil, // "byte", "str", "null_string", "time", "null_time",
		48.98,    //	"decimal",
		nil, nil, // "text", "binary",
	))
	t.Run("decimal int error", runner(
		errors.NotSupported,
		"",
		true, false, // "bool", "null_bool"
		1, 2, 3, // "int", "int64", "null_int64"
		6.4, nil, // "float64", "null_float64",
		0, 8, 16, 32, 64, // "uint", "uint8", "uint16", "uint32", "uint64"
		nil, "", nil, time.Time{}, nil, // "byte", "str", "null_string", "time", "null_time",
		4898,     //	"decimal",
		nil, nil, // "text", "binary",
	))
	t.Run("uint float error", runner(
		errors.NotSupported,
		"",
		true, false, // "bool", "null_bool"
		1, 2, 3, // "int", "int64", "null_int64"
		6.4, nil, // "float64", "null_float64",
		0.1, 8, 16, 32, 64, // "uint", "uint8", "uint16", "uint32", "uint64"
		nil, "", nil, time.Time{}, nil, // "byte", "str", "null_string", "time", "null_time",
		48.98,    //	"decimal",
		nil, nil, // "text", "binary",
	))
	t.Run("uint8 float error", runner(
		errors.NotSupported,
		"",
		true, false, // "bool", "null_bool"
		1, 2, 3, // "int", "int64", "null_int64"
		6.4, nil, // "float64", "null_float64",
		1, 8.1, 16, 32, 64, // "uint", "uint8", "uint16", "uint32", "uint64"
		nil, "", nil, time.Time{}, nil, // "byte", "str", "null_string", "time", "null_time",
		48.98,    //	"decimal",
		nil, nil, // "text", "binary",
	))
	t.Run("uint16 float error", runner(
		errors.NotSupported,
		"",
		true, false, // "bool", "null_bool"
		1, 2, 3, // "int", "int64", "null_int64"
		6.4, nil, // "float64", "null_float64",
		1, 8, 1.6, 32, 64, // "uint", "uint8", "uint16", "uint32", "uint64"
		nil, "", nil, time.Time{}, nil, // "byte", "str", "null_string", "time", "null_time",
		48.98,    //	"decimal",
		nil, nil, // "text", "binary",
	))
	t.Run("uint32 float error", runner(
		errors.NotSupported,
		"",
		true, false, // "bool", "null_bool"
		1, 2, 3, // "int", "int64", "null_int64"
		6.4, nil, // "float64", "null_float64",
		1, 8, 16, 3.2, 64, // "uint", "uint8", "uint16", "uint32", "uint64"
		nil, "", nil, time.Time{}, nil, // "byte", "str", "null_string", "time", "null_time",
		48.98,    //	"decimal",
		nil, nil, // "text", "binary",
	))
	t.Run("uint64 float error", runner(
		errors.NotSupported,
		"",
		true, false, // "bool", "null_bool"
		1, 2, 3, // "int", "int64", "null_int64"
		6.4, nil, // "float64", "null_float64",
		1, 8, 16, 32, 6.4, // "uint", "uint8", "uint16", "uint32", "uint64"
		nil, "", nil, time.Time{}, nil, // "byte", "str", "null_string", "time", "null_time",
		48.98,    //	"decimal",
		nil, nil, // "text", "binary",
	))
	t.Run("byte as float error", runner(
		errors.NotSupported,
		"",
		true, false, // "bool", "null_bool"
		1, 2, 3, // "int", "int64", "null_int64"
		6.4, nil, // "float64", "null_float64",
		1, 8, 16, 32, 64, // "uint", "uint8", "uint16", "uint32", "uint64"
		4.5678, "", nil, time.Time{}, nil, // "byte", "str", "null_string", "time", "null_time",
		48.98,    //	"decimal",
		nil, nil, // "text", "binary",
	))
	t.Run("string as float error", runner(
		errors.NotSupported,
		"",
		true, false, // "bool", "null_bool"
		1, 2, 3, // "int", "int64", "null_int64"
		6.4, nil, // "float64", "null_float64",
		1, 8, 16, 32, 64, // "uint", "uint8", "uint16", "uint32", "uint64"
		[]byte(`Hi`), 12.8, nil, time.Time{}, nil, // "byte", "str", "null_string", "time", "null_time",
		48.98,    //	"decimal",
		nil, nil, // "text", "binary",
	))
	t.Run("string as invalid utf8 string error", runner(
		errors.NotValid,
		"",
		true, false, // "bool", "null_bool"
		1, 2, 3, // "int", "int64", "null_int64"
		6.4, nil, // "float64", "null_float64",
		1, 8, 16, 32, 64, // "uint", "uint8", "uint16", "uint32", "uint64"
		[]byte(`Hi`), string("\xc0\x80"), nil, time.Time{}, nil, // "byte", "str", "null_string", "time", "null_time",
		48.98,    //	"decimal",
		nil, nil, // "text", "binary",
	))
	t.Run("string as invalid utf8 byte error", runner(
		errors.NotValid,
		"",
		true, false, // "bool", "null_bool"
		1, 2, 3, // "int", "int64", "null_int64"
		6.4, nil, // "float64", "null_float64",
		1, 8, 16, 32, 64, // "uint", "uint8", "uint16", "uint32", "uint64"
		[]byte(`Hi`), []byte("\xc0\x80"), nil, time.Time{}, nil, // "byte", "str", "null_string", "time", "null_time",
		48.98,    //	"decimal",
		nil, nil, // "text", "binary",
	))
	t.Run("null string as invalid utf8 byte error", runner(
		errors.NotValid,
		"",
		true, false, // "bool", "null_bool"
		1, 2, 3, // "int", "int64", "null_int64"
		6.4, nil, // "float64", "null_float64",
		1, 8, 16, 32, 64, // "uint", "uint8", "uint16", "uint32", "uint64"
		[]byte(`Hi`), "x", []byte("\xc0\x80"), time.Time{}, nil, // "byte", "str", "null_string", "time", "null_time",
		48.98,    //	"decimal",
		nil, nil, // "text", "binary",
	))
	t.Run("null string as invalid int error", runner(
		errors.NotSupported,
		"",
		true, false, // "bool", "null_bool"
		1, 2, 3, // "int", "int64", "null_int64"
		6.4, nil, // "float64", "null_float64",
		1, 8, 16, 32, 64, // "uint", "uint8", "uint16", "uint32", "uint64"
		[]byte(`Hi`), "x", 8767, time.Time{}, nil, // "byte", "str", "null_string", "time", "null_time",
		48.98,    //	"decimal",
		nil, nil, // "text", "binary",
	))
	t.Run("time as invalid int error", runner(
		errors.NotSupported,
		"",
		true, false, // "bool", "null_bool"
		1, 2, 3, // "int", "int64", "null_int64"
		6.4, nil, // "float64", "null_float64",
		1, 8, 16, 32, 64, // "uint", "uint8", "uint16", "uint32", "uint64"
		[]byte(`Hi`), "x", "y", 123456789, nil, // "byte", "str", "null_string", "time", "null_time",
		48.98,    //	"decimal",
		nil, nil, // "text", "binary",
	))
	t.Run("null time as invalid int error", runner(
		errors.NotSupported,
		"",
		true, false, // "bool", "null_bool"
		1, 2, 3, // "int", "int64", "null_int64"
		6.4, nil, // "float64", "null_float64",
		1, 8, 16, 32, 64, // "uint", "uint8", "uint16", "uint32", "uint64"
		[]byte(`Hi`), "x", "y", time.Time{}, 123456789, // "byte", "str", "null_string", "time", "null_time",
		48.98,    //	"decimal",
		nil, nil, // "text", "binary",
	))
	t.Run("null time as invalid string error", runner(
		errors.NotValid,
		"",
		true, false, // "bool", "null_bool"
		1, 2, 3, // "int", "int64", "null_int64"
		6.4, nil, // "float64", "null_float64",
		1, 8, 16, 32, 64, // "uint", "uint8", "uint16", "uint32", "uint64"
		[]byte(`Hi`), "x", "y", time.Time{}, "hello", // "byte", "str", "null_string", "time", "null_time",
		48.98,    //	"decimal",
		nil, nil, // "text", "binary",
	))
	t.Run("null time as invalid byte error", runner(
		errors.NotValid,
		"",
		true, false, // "bool", "null_bool"
		1, 2, 3, // "int", "int64", "null_int64"
		6.4, nil, // "float64", "null_float64",
		1, 8, 16, 32, 64, // "uint", "uint8", "uint16", "uint32", "uint64"
		[]byte(`Hi`), "x", "y", time.Time{}, []byte("hello"), // "byte", "str", "null_string", "time", "null_time",
		48.98,    //	"decimal",
		nil, nil, // "text", "binary",
	))

	t.Run("text/binary valid", runner(
		errors.NoKind,
		"bool: \"true\"\nnull_bool: \"false\"\nint: \"1\"\nint64: \"2\"\nnull_int64: \"3\"\nfloat64: \"6.4\"\nnull_float64: <nil>\nuint: \"1\"\nuint8: \"8\"\nuint16: \"16\"\nuint32: \"32\"\nuint64: \"64\"\nbyte: \"Hi\"\nstr: \"x\"\nnull_string: \"y\"\ntime: \"0001-01-01 00:00:00 +0000 UTC\"\nnull_time: <nil>\ndecimal: \"48.98\"\ntext: \"Hello World Text\"\nbinary: \"Hello Binary\"",
		true, false, // "bool", "null_bool"
		1, 2, 3, // "int", "int64", "null_int64"
		6.4, nil, // "float64", "null_float64",
		1, 8, 16, 32, 64, // "uint", "uint8", "uint16", "uint32", "uint64"
		[]byte(`Hi`), "x", "y", time.Time{}, nil, // "byte", "str", "null_string", "time", "null_time",
		48.98,                                              //	"decimal",
		[]byte(`Hello World Text`), []byte("Hello Binary"), // "text", "binary",
	))
	t.Run("text invalid UTF8", runner(
		errors.NotValid,
		"",
		true, false, // "bool", "null_bool"
		1, 2, 3, // "int", "int64", "null_int64"
		6.4, nil, // "float64", "null_float64",
		1, 8, 16, 32, 64, // "uint", "uint8", "uint16", "uint32", "uint64"
		[]byte(`Hi`), "x", "y", time.Time{}, nil, // "byte", "str", "null_string", "time", "null_time",
		48.98,                                            //	"decimal",
		[]byte("hello \xc0\x80"), []byte("Hello Binary"), // "text", "binary",
	))
	t.Run("text type error", runner(
		errors.NotSupported,
		"",
		true, false, // "bool", "null_bool"
		1, 2, 3, // "int", "int64", "null_int64"
		6.4, nil, // "float64", "null_float64",
		1, 8, 16, 32, 64, // "uint", "uint8", "uint16", "uint32", "uint64"
		[]byte(`Hi`), "x", "y", time.Time{}, nil, // "byte", "str", "null_string", "time", "null_time",
		48.98,                                    //	"decimal",
		"hello \xc0\x80", []byte("Hello Binary"), // "text", "binary",
	))
	t.Run("binary type error", runner(
		errors.NotSupported,
		"",
		true, false, // "bool", "null_bool"
		1, 2, 3, // "int", "int64", "null_int64"
		6.4, nil, // "float64", "null_float64",
		1, 8, 16, 32, 64, // "uint", "uint8", "uint16", "uint32", "uint64"
		[]byte(`Hi`), "x", "y", time.Time{}, nil, // "byte", "str", "null_string", "time", "null_time",
		48.98,               //	"decimal",
		nil, "Hello Binary", // "text", "binary",
	))

	t.Run("text marshal error", runner(
		errors.Empty,
		"",
		true, false, // "bool", "null_bool"
		1, 2, 3, // "int", "int64", "null_int64"
		6.4, nil, // "float64", "null_float64",
		1, 8, 16, 32, 64, // "uint", "uint8", "uint16", "uint32", "uint64"
		[]byte(`Hi`), "x", "y", time.Time{}, nil, // "byte", "str", "null_string", "time", "null_time",
		48.98,                //	"decimal",
		[]byte("error"), nil, // "text", "binary",
	))
	t.Run("binary marshal error", runner(
		errors.Empty,
		"",
		true, false, // "bool", "null_bool"
		1, 2, 3, // "int", "int64", "null_int64"
		6.4, nil, // "float64", "null_float64",
		1, 8, 16, 32, 64, // "uint", "uint8", "uint16", "uint32", "uint64"
		[]byte(`Hi`), "x", "y", time.Time{}, nil, // "byte", "str", "null_string", "time", "null_time",
		48.98,                //	"decimal",
		nil, []byte("error"), // "text", "binary",
	))

}
