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

package dml

import (
	"database/sql"
	"database/sql/driver"
	"fmt"
	"testing"
	"time"

	"github.com/corestoreio/errors"
	"github.com/corestoreio/pkg/storage/null"
	"github.com/corestoreio/pkg/util/assert"
)

var (
	_ fmt.Stringer = (*ip)(nil)
	_ QueryBuilder = (*ip)(nil)
)

func TestExpandPlaceHolders(t *testing.T) {
	t.Parallel()

	cp, err := NewConnPool()
	assert.NoError(t, err)

	t.Run("MisMatch", func(t *testing.T) {
		a := cp.WithRawSQL("SELECT * FROM `table` WHERE id IN (?)").ExpandPlaceHolders()
		compareToSQL2(t, a, errors.NoKind, "SELECT * FROM `table` WHERE id IN (?)")
	})
	t.Run("MisMatch length reps", func(t *testing.T) {
		a := cp.WithRawSQL("SELECT * FROM `table` WHERE id IN ?").ExpandPlaceHolders().TestWithArgs([]int64{1, 2}, []string{"d", "3"})
		compareToSQL2(t, a, errors.Mismatch, "")
	})
	t.Run("MisMatch qMarks", func(t *testing.T) {
		a := cp.WithRawSQL("SELECT * FROM `table` WHERE id IN(!)").ExpandPlaceHolders().TestWithArgs(3)
		compareToSQL2(t, a, errors.Mismatch, "")
	})
	t.Run("one arg with one value", func(t *testing.T) {
		a := cp.WithRawSQL("SELECT * FROM `table` WHERE id IN (?)").ExpandPlaceHolders().TestWithArgs(1)
		compareToSQL2(t, a, errors.NoKind, "SELECT * FROM `table` WHERE id IN (?)", int64(1))
	})
	t.Run("one arg with three values", func(t *testing.T) {
		a := cp.WithRawSQL("SELECT * FROM `table` WHERE id IN ?").ExpandPlaceHolders().TestWithArgs([]int64{11, 3, 5})
		compareToSQL2(t, a, errors.NoKind, "SELECT * FROM `table` WHERE id IN (?,?,?)", int64(11), int64(3), int64(5))
	})
	t.Run("multi 3,5 times replacement", func(t *testing.T) {
		a := cp.WithRawSQL("SELECT * FROM `table` WHERE id IN ? AND name IN ?").ExpandPlaceHolders().
			TestWithArgs([]int64{5, 7, 9}, []string{"a", "b", "c", "d", "e"})
		compareToSQL2(t, a, errors.NoKind, "SELECT * FROM `table` WHERE id IN (?,?,?) AND name IN (?,?,?,?,?)",
			int64(5), int64(7), int64(9), "a", "b", "c", "d", "e",
		)
	})
}

func TestInterpolate_Nil(t *testing.T) {
	t.Parallel()
	t.Run("one nil", func(t *testing.T) {
		ip := Interpolate("SELECT * FROM x WHERE a = ?").Null()
		assert.Exactly(t, "SELECT * FROM x WHERE a = NULL", ip.String())

		ip = Interpolate("SELECT * FROM x WHERE a = ?").Null()
		assert.Exactly(t, "SELECT * FROM x WHERE a = NULL", ip.String())
	})
	t.Run("two nil", func(t *testing.T) {
		ip := Interpolate("SELECT * FROM x WHERE a BETWEEN ? AND ?").Null().Null()
		assert.Exactly(t, "SELECT * FROM x WHERE a BETWEEN NULL AND NULL", ip.String())

		ip = Interpolate("SELECT * FROM x WHERE a BETWEEN ? AND ?").Null().Null()
		assert.Exactly(t, "SELECT * FROM x WHERE a BETWEEN NULL AND NULL", ip.String())
	})
	t.Run("one nil between two values", func(t *testing.T) {
		ip := Interpolate("SELECT * FROM x WHERE a BETWEEN ? AND ? OR Y=?").Int(1).Null().Str("X")
		assert.Exactly(t, "SELECT * FROM x WHERE a BETWEEN 1 AND NULL OR Y='X'", ip.String())
	})
}

func TestInterpolate_AllTypes(t *testing.T) {
	t.Parallel()

	sqlStr, args, err := Interpolate(`SELECT ?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?`).
		Null().
		Unsafe(`Unsafe`).
		Int(2).
		Ints(3, 4).
		Int64(5).
		Int64s(6, 7).
		Uint64(8).
		Uint64s(9, 10).
		Float64(11.1).
		Float64s(12.12, 13.13).
		Str("14").
		Strs("15", "16").
		Bool(true).
		Bools(false, true, false).
		Bytes([]byte(`17-18`)).
		BytesSlice([]byte(`19-20`), nil, []byte(`21`)).
		Time(now()).
		Times(now(), now()).
		NullString(null.MakeString("22")).
		NullStrings(null.MakeString("23"), null.String{}, null.MakeString("24")).
		NullFloat64(null.MakeFloat64(25.25)).
		NullFloat64s(null.MakeFloat64(26.26), null.Float64{}, null.MakeFloat64(27.27)).
		NullInt64(null.MakeInt64(28)).
		NullInt64s(null.MakeInt64(29), null.Int64{}, null.MakeInt64(30)).
		NullBool(null.MakeBool(true)).
		NullBools(null.MakeBool(true), null.Bool{}, null.MakeBool(false)).
		NullTime(null.MakeTime(now())).
		NullTimes(null.MakeTime(now()), null.Time{}, null.MakeTime(now())).
		ToSQL()
	assert.NoError(t, err)
	assert.Nil(t, args)
	assert.Exactly(t,
		"SELECT NULL,'Unsafe',2,(3,4),5,(6,7),8,(9,10),11.1,(12.12,13.13),'14',('15','16'),1,(0,1,0),'17-18',('19-20',NULL,'21'),'2006-01-02 15:04:05',('2006-01-02 15:04:05','2006-01-02 15:04:05'),'22',('23',NULL,'24'),25.25,(26.26,NULL,27.27),28,(29,NULL,30),1,(1,NULL,0),'2006-01-02 15:04:05',('2006-01-02 15:04:05',NULL,'2006-01-02 15:04:05')",
		sqlStr)
}

func TestInterpolate_Errors(t *testing.T) {
	t.Parallel()
	t.Run("non utf8", func(t *testing.T) {
		compareToSQL2(t,
			Interpolate("SELECT * FROM x WHERE a = ?").Str(string([]byte{0x34, 0xFF, 0xFE})),
			errors.NotValid,
			"",
		)
	})
	t.Run("too few qmarks", func(t *testing.T) {
		compareToSQL2(t,
			Interpolate("SELECT * FROM x WHERE a = ?").Int(3).Int(4),
			errors.Mismatch,
			"",
		)
	})
	t.Run("too many qmarks", func(t *testing.T) {
		compareToSQL2(t,
			Interpolate("SELECT * FROM x WHERE a = ? OR b = ? or c = ?").Ints(3, 4),
			errors.Mismatch,
			"",
		)
	})
	t.Run("way too many qmarks", func(t *testing.T) {
		_, _, err := Interpolate("SELECT * FROM x WHERE a IN ? OR b = ? OR c = ? AND d = ?").Ints(3, 4).Int64(2).ToSQL()
		assert.ErrorIsKind(t, errors.Mismatch, err)
	})
	t.Run("print error into String", func(t *testing.T) {
		ip := Interpolate("SELECT * FROM x WHERE a IN (?) AND b BETWEEN ? AND ? AND c = ? AND d IN (?,?)").
			Int64(3).Int64(-3).
			Uint64(7).
			Float64s(3.5, 4.4995)
		assert.Exactly(t, "[dml] Number of place holders (6) vs number of arguments (4) do not match.", ip.String())
	})
}

type argValUint16 uint16

func (u argValUint16) Value() (driver.Value, error) {
	if u > 0 {
		return nil, errors.Aborted.Newf("Not in the mood today")
	}
	return uint16(0), nil
}

func TestInterpolate_ArgValue(t *testing.T) {
	t.Parallel()

	aInt := null.MakeInt64(4711)
	aStr := null.MakeString("Goph'er")
	aFlo := null.MakeFloat64(2.7182818)
	aTim := null.MakeTime(Now.UTC())
	aBoo := null.MakeBool(true)
	aByt := driverValueBytes(`BytyGophe'r`)
	var aNil driverValueBytes

	t.Run("equal", func(t *testing.T) {
		compareToSQL2(t,
			Interpolate("SELECT * FROM x WHERE a = ? AND b = ? AND c = ? AND d = ? AND e = ? AND f = ? AND g = ? AND h = ?").
				DriverValues(aInt).
				DriverValues(aStr).
				DriverValues(aFlo).
				DriverValues(aTim).
				DriverValues(aBoo).
				DriverValues(aByt).
				DriverValues(aNil).
				DriverValues(aNil),
			errors.NoKind,
			"SELECT * FROM x WHERE a = 4711 AND b = 'Goph\\'er' AND c = 2.7182818 AND d = '2006-01-02 19:04:05' AND e = 1 AND f = 'BytyGophe\\'r' AND g = NULL AND h = NULL",
		)
	})
	t.Run("in", func(t *testing.T) {
		compareToSQL2(t,
			Interpolate("SELECT * FROM x WHERE a IN ? AND b IN ? AND c IN ? AND d IN ? AND e IN ? AND f IN ?").
				DriverValue(aInt, aInt).DriverValue(aStr, aStr).DriverValue(aFlo, aFlo).
				DriverValue(aTim, aTim).DriverValue(aBoo, aBoo).DriverValue(aByt, aByt),
			errors.NoKind,
			"SELECT * FROM x WHERE a IN (4711,4711) AND b IN ('Goph\\'er','Goph\\'er') AND c IN (2.7182818,2.7182818) AND d IN ('2006-01-02 19:04:05','2006-01-02 19:04:05') AND e IN (1,1) AND f IN ('BytyGophe\\'r','BytyGophe\\'r')",
		)
	})
	t.Run("type not supported", func(t *testing.T) {
		_, _, err := Interpolate("SELECT * FROM x WHERE a = ?").DriverValues(argValUint16(0)).ToSQL()
		assert.ErrorIsKind(t, errors.NotSupported, err)
	})
	t.Run("valuer error", func(t *testing.T) {
		_, _, err := Interpolate("SELECT * FROM x WHERE a = ?").DriverValues(argValUint16(1)).ToSQL()
		assert.ErrorIsKind(t, errors.Aborted, err)
	})
}

func TestInterpolate_Reset(t *testing.T) {
	t.Parallel()

	t.Run("call twice with different arguments", func(t *testing.T) {
		ip := Interpolate("SELECT * FROM x WHERE a IN ? AND b BETWEEN ? AND ? AND c = ? AND d IN ?").
			Ints(1, -2).
			Int64(3).Int64(-3).
			Uint64(7).
			Float64s(3.5, 4.4995)
		assert.Exactly(t, "SELECT * FROM x WHERE a IN (1,-2) AND b BETWEEN 3 AND -3 AND c = 7 AND d IN (3.5,4.4995)", ip.String())

		ip.Reset().
			Ints(10, -20).
			Int64(30).Int64(-30).
			Uint64(70).
			Float64s(30.5, 40.4995)
		assert.Exactly(t, "SELECT * FROM x WHERE a IN (10,-20) AND b BETWEEN 30 AND -30 AND c = 70 AND d IN (30.5,40.4995)", ip.String())
	})
}

func TestInterpolate_Int64(t *testing.T) {
	t.Parallel()

	t.Run("equal named params", func(t *testing.T) {
		compareToSQL2(t,
			Interpolate("SELECT * FROM x WHERE a = (:ArgX) AND b > @ArgY").
				Named(
					sql.Named(":ArgX", 3),
					sql.Named("@ArgY", 3.14159),
				),
			errors.NoKind,
			"SELECT * FROM x WHERE a = (3) AND b > 3.14159",
		)
	})
	t.Run("equal", func(t *testing.T) {
		compareToSQL2(t,
			Interpolate("SELECT * FROM x WHERE a = ? AND b = ? AND c = ?").
				Int64(1).Int64(-2).Int64(3),
			errors.NoKind,
			"SELECT * FROM x WHERE a = 1 AND b = -2 AND c = 3",
		)
	})
	t.Run("in", func(t *testing.T) {
		compareToSQL2(t,
			Interpolate("SELECT * FROM x WHERE a IN ?").Int64s(1, -2, 3, 4, 5, 6, 7, 8, 9, 10),
			errors.NoKind,
			"SELECT * FROM x WHERE a IN (1,-2,3,4,5,6,7,8,9,10)",
		)
	})
	t.Run("in and equal", func(t *testing.T) {
		compareToSQL2(t,
			Interpolate("SELECT * FROM x WHERE a = ? AND b = ? AND c = ? AND h = ? AND i = ? AND j = ? AND k = ? AND m IN ? OR n = ?").
				Int64(1).
				Int64(-2).
				Int64(3).
				Int64(4).
				Int64(5).
				Int64(6).
				Int64(11).
				Int64s(12, 13).
				Int64(-14),
			errors.NoKind,
			`SELECT * FROM x WHERE a = 1 AND b = -2 AND c = 3 AND h = 4 AND i = 5 AND j = 6 AND k = 11 AND m IN (12,13) OR n = -14`,
		)
	})
}

func TestInterpolate_Bools(t *testing.T) {
	t.Parallel()

	t.Run("single args", func(t *testing.T) {
		compareToSQL2(t,
			Interpolate("SELECT * FROM x WHERE a = ? AND b = ?").Bool(true).Bool(false),
			errors.NoKind,
			"SELECT * FROM x WHERE a = 1 AND b = 0",
		)
	})
	t.Run("IN args", func(t *testing.T) {
		compareToSQL2(t,
			Interpolate("SELECT * FROM x WHERE a IN ? AND b = ? OR c = ?").
				Bools(true, false).Bool(true).Bool(false),
			errors.NoKind,
			"SELECT * FROM x WHERE a IN (1,0) AND b = 1 OR c = 0",
		)
	})
}

func TestInterpolate_Bytes(t *testing.T) {
	t.Parallel()

	b1 := []byte(`Go`)
	b2 := []byte(`Further`)
	t.Run("UTF8 valid: single args", func(t *testing.T) {
		compareToSQL2(t,
			Interpolate("SELECT * FROM x WHERE a = ? AND b = ?").Bytes(b1).Bytes(b2),
			errors.NoKind,
			"SELECT * FROM x WHERE a = 'Go' AND b = 'Further'",
		)
	})
	t.Run("UTF8 valid: IN args", func(t *testing.T) {
		compareToSQL2(t,
			Interpolate("SELECT * FROM x WHERE a IN ?").BytesSlice(b1, b2),
			errors.NoKind,
			"SELECT * FROM x WHERE a IN ('Go','Further')",
		)
	})
	t.Run("empty arg triggers no error", func(t *testing.T) {
		compareToSQL2(t,
			Interpolate("SELECT * FROM x WHERE a IN ?").BytesSlice(),
			errors.NoKind,
			"SELECT * FROM x WHERE a IN ()",
		)
	})
	t.Run("Binary to hex", func(t *testing.T) {
		bin := []byte{66, 250, 67}
		compareToSQL2(t,
			Interpolate("SELECT * FROM x WHERE a = ?").Bytes(bin),
			errors.NoKind,
			"SELECT * FROM x WHERE a = 0x42fa43",
		)
	})
}

func TestInterpolate_Time(t *testing.T) {
	t.Parallel()

	t1 := now()
	t2 := now().Add(time.Minute)

	t.Run("single args", func(t *testing.T) {
		compareToSQL2(t,
			Interpolate("SELECT * FROM x WHERE a = ? AND b = ?").Time(t1).Time(t2),
			errors.NoKind,
			"SELECT * FROM x WHERE a = '2006-01-02 15:04:05' AND b = '2006-01-02 15:05:05'",
		)
	})
	t.Run("IN args", func(t *testing.T) {
		compareToSQL2(t,
			Interpolate("SELECT * FROM x WHERE a IN ?").Times(t1, t2),
			errors.NoKind,
			"SELECT * FROM x WHERE a IN ('2006-01-02 15:04:05','2006-01-02 15:05:05')",
		)
	})
	t.Run("empty arg", func(t *testing.T) {
		compareToSQL2(t,
			Interpolate("SELECT * FROM x WHERE a IN ? ?").Times(),
			errors.Mismatch,
			"",
		)
	})
}

func TestInterpolate_Floats(t *testing.T) {
	t.Parallel()
	t.Run("single args", func(t *testing.T) {
		compareToSQL2(t,
			Interpolate("SELECT * FROM x WHERE a = ? AND b = ?").Float64(3.14159).Float64(2.7182818),
			errors.NoKind,
			"SELECT * FROM x WHERE a = 3.14159 AND b = 2.7182818",
		)
	})
	t.Run("IN args", func(t *testing.T) {
		compareToSQL2(t,
			Interpolate("SELECT * FROM x WHERE a IN ? AND b = ?").Float64s(3.14159, 2.7182818).Float64(0.815),
			errors.NoKind,
			"SELECT * FROM x WHERE a IN (3.14159,2.7182818) AND b = 0.815",
		)
	})
	t.Run("IN args reverse", func(t *testing.T) {
		compareToSQL2(t,
			Interpolate("SELECT * FROM x WHERE b = ? AND a IN ?").Float64(0.815).Float64s(3.14159, 2.7182818),
			errors.NoKind,
			"SELECT * FROM x WHERE b = 0.815 AND a IN (3.14159,2.7182818)",
		)
	})
}

func TestInterpolate_Slices_Strings_Between(t *testing.T) {
	t.Parallel()

	t.Run("BETWEEN at the end", func(t *testing.T) {
		compareToSQL2(t,
			Interpolate("SELECT * FROM x WHERE a IN ? AND b IN ? AND c NOT IN ? AND d BETWEEN ? AND ?").
				Ints(1).Ints(1, 2, 3).Int64s(5, 6, 7).Str("wat").Str("ok"),
			errors.NoKind,
			"SELECT * FROM x WHERE a IN (1) AND b IN (1,2,3) AND c NOT IN (5,6,7) AND d BETWEEN 'wat' AND 'ok'",
		)
	})
	t.Run("BETWEEN in the middle", func(t *testing.T) {
		compareToSQL2(t,
			Interpolate("SELECT * FROM x WHERE a IN ? AND b IN ? AND d BETWEEN ? AND ? AND c NOT IN ?").
				Ints(1).Ints(1, 2, 3).Str("wat").Str("ok").Int64s(5, 6, 7),
			errors.NoKind,
			"SELECT * FROM x WHERE a IN (1) AND b IN (1,2,3) AND d BETWEEN 'wat' AND 'ok' AND c NOT IN (5,6,7)",
		)
	})

	t.Run("single args", func(t *testing.T) {
		compareToSQL2(t,
			Interpolate("SELECT * FROM x WHERE a = ? AND b = ? AND c = ?").Str("a'b").Str("c`d").Str("\"hello's \\ world\" \n\r\x00\x1a"),
			errors.NoKind,
			"SELECT * FROM x WHERE a = 'a\\'b' AND b = 'c`d' AND c = '\\\"hello\\'s \\\\ world\\\" \\n\\r\\x00\\x1a'",
		)
	})
	t.Run("IN args", func(t *testing.T) {
		compareToSQL2(t,
			Interpolate("SELECT * FROM x WHERE a IN ? AND b = ?").Strs("a'b", "c`d").Strs("1' or '1' = '1'))/*"),
			errors.NoKind,
			"SELECT * FROM x WHERE a IN ('a\\'b','c`d') AND b = ('1\\' or \\'1\\' = \\'1\\'))/*')",
		)
	})
	t.Run("empty args triggers incorrect interpolation of values", func(t *testing.T) {
		sl := make([]string, 0, 2)
		compareToSQL2(t,
			Interpolate("SELECT * FROM x WHERE a IN ? AND b = ? OR c = ?").Strs("a", "b").Str("c").Strs(sl...),
			errors.NoKind,
			"SELECT * FROM x WHERE a IN ('a','b') AND b = 'c' OR c = ()",
		)
	})
	t.Run("multiple slices", func(t *testing.T) {
		compareToSQL2(t,
			Interpolate("SELECT * FROM x WHERE a = ? AND b = ? AND c = ? AND d = ? AND e = ?").
				Ints(1).
				Ints(1, 2, 3).
				Int64s(5, 6, 7).
				Strs("wat", "ok").
				Int(8),
			errors.NoKind,
			"SELECT * FROM x WHERE a = (1) AND b = (1,2,3) AND c = (5,6,7) AND d = ('wat','ok') AND e = 8",
		)
	})
}

func TestInterpolate_Different(t *testing.T) {
	t.Parallel()
	ifs := func(vals ...interface{}) []interface{} { return vals }
	tests := []struct {
		sql     string
		argsIn  []interface{}
		expSQL  string
		errKind errors.Kind
	}{
		// NULL
		{
			"SELECT * FROM x WHERE a = ?", ifs(nil),
			"SELECT * FROM x WHERE a = NULL", errors.NoKind,
		},

		// integers
		{
			`SELECT * FROM x WHERE a = ? AND b = ?`,
			ifs(1, -2),
			`SELECT * FROM x WHERE a = 1 AND b = -2`, errors.NoKind,
		},
		{
			`SELECT * FROM x WHERE a IN ?`,
			ifs([]int{1, -2, 3, 4, 5, 6, 7, 8, 9, 10}),
			`SELECT * FROM x WHERE a IN (1,-2,3,4,5,6,7,8,9,10)`, errors.NoKind,
		},
		{
			`SELECT * FROM x WHERE a IN ?`,
			ifs([]int64{1, -2, 3, 4, 5, 6, 7, 8, 9, 10}),
			`SELECT * FROM x WHERE a IN (1,-2,3,4,5,6,7,8,9,10)`, errors.NoKind,
		},

		// boolean
		{
			"SELECT * FROM x WHERE a = ? AND b = ?", ifs(true, false),
			"SELECT * FROM x WHERE a = 1 AND b = 0", errors.NoKind,
		},
		{
			"SELECT * FROM x WHERE a IN ?", ifs([]bool{true, false}),
			"SELECT * FROM x WHERE a IN (1,0)", errors.NoKind,
		},

		// floats
		{
			"SELECT * FROM x WHERE a = ? AND b = ?", ifs(0.15625, 3.14159),
			"SELECT * FROM x WHERE a = 0.15625 AND b = 3.14159", errors.NoKind,
		},
		{
			"SELECT * FROM x WHERE a IN ?", ifs([]float64{0.15625, 3.14159}),
			"SELECT * FROM x WHERE a IN (0.15625,3.14159)", errors.NoKind,
		},
		{
			"SELECT * FROM x WHERE a = ? AND b = ? and C = ?", ifs([]float64{0.15625, 3.14159}),
			"", errors.Mismatch,
		},
		{
			`SELECT * FROM x WHERE a IN ?`,
			ifs([]float64{1.1, -2.2, 3.3}),
			`SELECT * FROM x WHERE a IN (1.1,-2.2,3.3)`, errors.NoKind,
		},

		// strings
		{
			`SELECT * FROM x WHERE a = ?
			AND b = ?`,
			ifs("hello", "\"hello's \\ world\" \n\r\x00\x1a"),
			`SELECT * FROM x WHERE a = 'hello'
			AND b = '\"hello\'s \\ world\" \n\r\x00\x1a'`, errors.NoKind,
		},
		{
			`SELECT * FROM x WHERE a IN ?`,
			ifs([]string{"a'a", "bb"}),
			`SELECT * FROM x WHERE a IN ('a\'a','bb')`, errors.NoKind,
		},
		{
			`SELECT * FROM x WHERE a IN ?`,
			ifs([]string{"a'a", "bb"}),
			`SELECT * FROM x WHERE a IN ('a\'a','bb')`, errors.NoKind,
		},

		// slices
		{
			"SELECT * FROM x WHERE a = ? AND b = ? AND c = ? AND d = ?",
			ifs(1, []int{1, 2, 3}, []int{5, 6, 7}, []string{"wat", "ok"}),
			"SELECT * FROM x WHERE a = 1 AND b = (1,2,3) AND c = (5,6,7) AND d = ('wat','ok')", errors.NoKind,
		},
		//
		////// TODO valuers
		////{"SELECT * FROM x WHERE a = ? AND b = ?",
		////	args{myString{true, "wat"}, myString{false, "fry"}},
		////	"SELECT * FROM x WHERE a = 'wat' AND b = NULL", nil},

		// errors
		{
			"SELECT * FROM x WHERE a = ? AND b = ?", ifs(int64(1)),
			"", errors.Mismatch,
		},

		{
			"SELECT * FROM x WHERE", ifs(1),
			"", errors.Mismatch,
		},

		{
			"SELECT * FROM x WHERE a = ?", ifs(string([]byte{0x34, 0xFF, 0xFE})),
			"", errors.NotValid,
		},

		// String() without arguments is equal to empty interface in the previous version.
		{"SELECT 'hello", ifs(""), "", errors.Mismatch},
		{`SELECT "hello`, ifs(""), "", errors.Mismatch},
		{`SELECT ? "hello`, ifs(""), "SELECT '' 'hello", errors.NoKind},

		// preprocessing
		{"SELECT '?'", nil, "SELECT '?'", errors.NoKind},
		{"SELECT `?`", nil, "SELECT `?`", errors.NoKind},
		{"SELECT [?]", nil, "SELECT `?`", errors.NoKind},
		{"SELECT [?]", nil, "SELECT `?`", errors.NoKind},
		{"SELECT [name] FROM [user]", nil, "SELECT `name` FROM `user`", errors.NoKind},
		{"SELECT [u.name] FROM [user] [u]", nil, "SELECT `u`.`name` FROM `user` `u`", errors.NoKind},
		{"SELECT [u.na`me] FROM [user] [u]", nil, "SELECT `u`.`na``me` FROM `user` `u`", errors.NoKind},
		{"SELECT * FROM [user] WHERE [name] = '[nick]'", nil, "SELECT * FROM `user` WHERE `name` = '[nick]'", errors.NoKind},
		{`SELECT * FROM [user] WHERE [name] = "nick[]"`, nil, "SELECT * FROM `user` WHERE `name` = 'nick[]'", errors.NoKind},
		{"SELECT * FROM [user] WHERE [name] = 'Hello`s World'", nil, "SELECT * FROM `user` WHERE `name` = 'Hello`s World'", errors.NoKind},
	}

	for _, test := range tests {
		str, _, err := Interpolate(test.sql).Unsafe(test.argsIn...).ToSQL()
		if !test.errKind.Empty() {
			if !test.errKind.Match(err) {
				isErr := test.errKind.Match(err)
				t.Errorf("IDX %q\ngot error: %v\nwant: %t", test.sql, err, isErr)
			}
		}
		// assert.NoError(t, err, "IDX %d", i)
		if str != test.expSQL {
			t.Errorf("IDX %q\ngot: %v\nwant: %v\nError: %+v", test.sql, str, test.expSQL, err)
		}
	}
}

func TestInterpolate_MultipleSingleQuotes(t *testing.T) {
	const rawSQL = "DELETE FROM `tableA` WHERE (`colA` >= 3.14159) AND (`colB` IN ?) AND (`colC` = 'He\\'llo') ORDER BY `id` LIMIT 10"
	str, args, err := Interpolate(rawSQL).Unsafe([]float64{3.1, 2.4}).ToSQL()
	assert.NoError(t, err)
	assert.Nil(t, args)
	assert.Exactly(t, "DELETE FROM `tableA` WHERE (`colA` >= 3.14159) AND (`colB` IN (3.1,2.4)) AND (`colC` = 'He\\'llo') ORDER BY `id` LIMIT 10", str)
}

func TestExtractNamedArgs(t *testing.T) {
	t.Parallel()

	runner := func(haveSQL, wantSQL string, wantQualifiedColumns ...string) func(*testing.T) {
		return func(t *testing.T) {
			gotSQL, qualifiedColumns, _ := extractReplaceNamedArgs(haveSQL, nil)
			assert.Exactly(t, wantSQL, string(gotSQL))
			assert.Exactly(t, wantQualifiedColumns, qualifiedColumns)
		}
	}
	t.Run("one", runner(
		"SELECT 1 AS `n`, CAST(:abc AS CHAR(20)) AS `str`",
		"SELECT 1 AS `n`, CAST(? AS CHAR(20)) AS `str`",
		namedArgStartStr+"abc",
	))
	t.Run("one in bracket", runner(
		"SELECT 1 AS `n`, CAST((:abc) AS CHAR(20)) AS `str`",
		"SELECT 1 AS `n`, CAST((?) AS CHAR(20)) AS `str`",
		namedArgStartStr+"abc",
	))
	t.Run("qualified one in bracket", runner(
		"SELECT 1 AS `n`, CAST((:xx.abc) AS CHAR(20)) AS `str`",
		"SELECT 1 AS `n`, CAST((?) AS CHAR(20)) AS `str`",
		namedArgStartStr+"xx.abc",
	))
	t.Run("none", runner(
		"SELECT 1 AS `n`, CAST(abc AS CHAR(20)) AS `str`",
		"SELECT 1 AS `n`, CAST(abc AS CHAR(20)) AS `str`",
	))
	t.Run("two same", runner(
		"SELECT 1 AS `n`, CAST(:abc AS CHAR(20)) AS `str`, CAST(:abc AS INT(20)) AS `intstr`",
		"SELECT 1 AS `n`, CAST(? AS CHAR(20)) AS `str`, CAST(? AS INT(20)) AS `intstr`",
		namedArgStartStr+"abc", namedArgStartStr+"abc",
	))
	t.Run("two different", runner(
		"SELECT 1 AS `n`, CAST(:abc2 AS CHAR(20)) AS `str`, CAST(:abc1 AS INT(20)) AS `intstr`",
		"SELECT 1 AS `n`, CAST(? AS CHAR(20)) AS `str`, CAST(? AS INT(20)) AS `intstr`",
		namedArgStartStr+"abc2", namedArgStartStr+"abc1",
	))
	t.Run("emoji and non-emoji with error", runner(
		"SELECT 1 AS `n`, CAST(:a😱bc AS CHAR(20)) AS `str`, CAST(:abc AS INT(20)) AS `intstr`",
		"SELECT 1 AS `n`, CAST(? AS CHAR(20)) AS `str`, CAST(?c AS INT(20)) AS `intstr`",
		namedArgStartStr+"a😱bc", namedArgStartStr+"ab",
	))
	t.Run("emoji only with incorrect SQL", runner(
		"SELECT 1 AS `n`, (:👨‍👨‍) AS `str`, (:🔜) AS `intstr`",
		"SELECT 1 AS `n`, (?\u200d👨\u200d) AS `str`, (?) AS `intstr`",
		namedArgStartStr+"👨", namedArgStartStr+"🔜",
	))
	t.Run("colon only", runner(
		"SELECT : AS `n`",
		"SELECT ? AS `n`",
	))
	t.Run("with number", runner(
		"SELECT (:x32)",
		"SELECT (?)",
		namedArgStartStr+"x32",
	))
	t.Run("date as argument short", runner(
		"CASE  WHEN date_start <= '2009-11-11 00:00:00'",
		"CASE  WHEN date_start <= '2009-11-11 00:00:00'",
	))
	t.Run("date as argument long", runner(
		"CASE  WHEN date_start <= '2009-11-11 00:00:00' AND date_end >= '2009-11-12 00:00:00' THEN `open` WHEN date_start > '2009-11-11 00:00:00' AND date_end > '2009-11-12 00:00:00' THEN `upcoming` ELSE `closed` END",
		"CASE  WHEN date_start <= '2009-11-11 00:00:00' AND date_end >= '2009-11-12 00:00:00' THEN `open` WHEN date_start > '2009-11-11 00:00:00' AND date_end > '2009-11-12 00:00:00' THEN `upcoming` ELSE `closed` END",
	))
	t.Run("single quote escaped", runner(
		"date_start = 'It\\'s xmas' ORDER BY X",
		"date_start = 'It\\'s xmas' ORDER BY X",
	))
}
