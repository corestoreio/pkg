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

package dml

import (
	"database/sql"
	"database/sql/driver"
	"fmt"
	"testing"
	"time"

	"github.com/corestoreio/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var _ fmt.Stringer = (*iPolate)(nil)
var _ QueryBuilder = (*iPolate)(nil)

func TestRepeat(t *testing.T) {
	t.Parallel()

	t.Run("MisMatch", func(t *testing.T) {
		s, err := Repeat("SELECT * FROM `table` WHERE id IN (?)", nil)
		assert.Empty(t, s)
		assert.True(t, errors.IsMismatch(err), "%+v", err)
	})
	t.Run("MisMatch length reps", func(t *testing.T) {
		s, err := Repeat("SELECT * FROM `table` WHERE id IN (?)", MakeArgs(2).Int64s(1, 2).Strings("d", "3"))
		assert.Empty(t, s)
		assert.True(t, errors.IsMismatch(err), "%+v", err)
	})
	t.Run("MisMatch qMarks", func(t *testing.T) {
		s, err := Repeat("SELECT * FROM `table` WHERE id IN(!)", MakeArgs(1).Int64(3))
		assert.Empty(t, s)
		assert.True(t, errors.IsMismatch(err), "%+v", err)
	})
	t.Run("one arg with one value", func(t *testing.T) {
		args := MakeArgs(1).Int64(1)
		s, err := Repeat("SELECT * FROM `table` WHERE id IN (?)", args)
		require.NoError(t, err)
		assert.Exactly(t, "SELECT * FROM `table` WHERE id IN (?)", s)
		assert.Exactly(t, []interface{}{int64(1)}, args.Interfaces())
	})
	t.Run("one arg with three values", func(t *testing.T) {
		args := MakeArgs(1).Int64s(11, 3, 5)
		s, err := Repeat("SELECT * FROM `table` WHERE id IN (?)", args)
		require.NoError(t, err)
		assert.Exactly(t, "SELECT * FROM `table` WHERE id IN (?,?,?)", s)
		assert.Exactly(t, []interface{}{int64(11), int64(3), int64(5)}, args.Interfaces())
	})
	t.Run("multi 3,5 times replacement", func(t *testing.T) {
		args := MakeArgs(3).Int64s(5, 7, 9).Strings("a", "b", "c", "d", "e")
		s, err := Repeat("SELECT * FROM `table` WHERE id IN (?) AND name IN (?)", args)
		require.NoError(t, err)
		assert.Exactly(t, "SELECT * FROM `table` WHERE id IN (?,?,?) AND name IN (?,?,?,?,?)", s)
		assert.Exactly(t, []interface{}{int64(5), int64(7), int64(9), "a", "b", "c", "d", "e"}, args.Interfaces())
	})
}

func TestInterpolateNil(t *testing.T) {
	t.Parallel()
	t.Run("one nil", func(t *testing.T) {
		ip := Interpolate("SELECT * FROM x WHERE a = ?").Null()
		assert.Equal(t, "SELECT * FROM x WHERE a = NULL", ip.String())

		ip = Interpolate("SELECT * FROM x WHERE a = ?").Null()
		assert.Equal(t, "SELECT * FROM x WHERE a = NULL", ip.String())
	})
	t.Run("two nil", func(t *testing.T) {
		ip := Interpolate("SELECT * FROM x WHERE a BETWEEN ? AND ?").Null().Null()
		assert.Equal(t, "SELECT * FROM x WHERE a BETWEEN NULL AND NULL", ip.String())

		ip = Interpolate("SELECT * FROM x WHERE a BETWEEN ? AND ?").Null().Null()
		assert.Equal(t, "SELECT * FROM x WHERE a BETWEEN NULL AND NULL", ip.String())
	})
	t.Run("one nil between two values", func(t *testing.T) {
		ip := Interpolate("SELECT * FROM x WHERE a BETWEEN ? AND ? OR Y=?").Int(1).Null().Str("X")
		assert.Equal(t, "SELECT * FROM x WHERE a BETWEEN 1 AND NULL OR Y='X'", ip.String())
	})
}

func TestInterpolateErrors(t *testing.T) {
	t.Parallel()
	t.Run("non utf8", func(t *testing.T) {
		_, _, err := Interpolate("SELECT * FROM x WHERE a = ?").Str(string([]byte{0x34, 0xFF, 0xFE})).ToSQL()
		assert.True(t, errors.IsNotValid(err), "%+v", err)
	})
	t.Run("too few qmarks", func(t *testing.T) {
		_, _, err := Interpolate("SELECT * FROM x WHERE a = ?").Int(3).Int(4).ToSQL()
		assert.True(t, errors.IsNotValid(err), "%+v", err)
	})
	t.Run("too many qmarks", func(t *testing.T) {
		_, _, err := Interpolate("SELECT * FROM x WHERE a = ? OR b = ? or c = ?").Ints(3, 4).ToSQL()
		assert.True(t, errors.IsNotValid(err), "%+v", err)
	})
	t.Run("way too many qmarks", func(t *testing.T) {
		_, _, err := Interpolate("SELECT * FROM x WHERE a IN ? OR b = ? OR c = ? AND d = ?").Ints(3, 4).Int64(2).ToSQL()
		assert.True(t, errors.IsNotValid(err), "%+v", err)
	})
	t.Run("print error into String", func(t *testing.T) {
		ip := Interpolate("SELECT * FROM x WHERE a IN (?) AND b BETWEEN ? AND ? AND c = ? AND d IN (?,?)").
			Int64(3).Int64(-3).
			Uint64(7).
			Float64s(3.5, 4.4995)
		assert.Exactly(t, "[dml] Interpolate: Arguments are imbalanced. Argument Index 4 is greater than argument count 3", ip.String())
	})
}

type argValUint16 uint16

func (u argValUint16) Value() (driver.Value, error) {
	if u > 0 {
		return nil, errors.NewAbortedf("Not in the mood today")
	}
	return uint16(0), nil
}

func TestInterpolate_ArgValue(t *testing.T) {
	t.Parallel()

	aInt := MakeNullInt64(4711)
	aStr := MakeNullString("Goph'er")
	aFlo := MakeNullFloat64(2.7182818)
	aTim := MakeNullTime(Now.UTC())
	aBoo := MakeNullBool(true)
	aByt := driverValueBytes(`BytyGophe'r`)
	var aNil driverValueBytes

	t.Run("equal", func(t *testing.T) {
		str := Interpolate("SELECT * FROM x WHERE a = ? AND b = ? AND c = ? AND d = ? AND e = ? AND f = ? AND g = ? AND h = ?").
			DriverValue(aInt).
			DriverValue(aStr).
			DriverValue(aFlo).
			DriverValue(aTim).
			DriverValue(aBoo).
			DriverValue(aByt).
			DriverValue(aNil).
			DriverValue(aNil).
			String()
		assert.Equal(t, "SELECT * FROM x WHERE a = 4711 AND b = 'Goph\\'er' AND c = 2.7182818 AND d = '2006-01-02 15:04:12' AND e = 1 AND f = 'BytyGophe\\'r' AND g = NULL AND h = NULL",
			str)
	})
	t.Run("in", func(t *testing.T) {
		str := Interpolate("SELECT * FROM x WHERE a IN (?) AND b IN (?) AND c IN (?) AND d IN (?) AND e IN (?) AND f IN (?)").
			DriverValue(aInt, aInt).DriverValue(aStr, aStr).DriverValue(aFlo, aFlo).
			DriverValue(aTim, aTim).DriverValue(aBoo, aBoo).DriverValue(aByt, aByt).String()
		assert.Equal(t,
			"SELECT * FROM x WHERE a IN (4711,4711) AND b IN ('Goph\\'er','Goph\\'er') AND c IN (2.7182818,2.7182818) AND d IN ('2006-01-02 15:04:12','2006-01-02 15:04:12') AND e IN (1,1) AND f IN ('BytyGophe\\'r','BytyGophe\\'r')",
			str)
	})
	t.Run("type not supported", func(t *testing.T) {
		defer func() {
			if r := recover(); r != nil {
				if err, ok := r.(error); ok {
					assert.True(t, errors.IsNotSupported(err), "%+v", err)
				} else {
					t.Errorf("Panic should contain an error but got:\n%+v", r)
				}
			} else {
				t.Error("Expecting a panic but got nothing")
			}
		}()
		_, _, _ = Interpolate("SELECT * FROM x WHERE a = ?").DriverValue(argValUint16(0)).ToSQL()
	})
	t.Run("valuer error", func(t *testing.T) {
		defer func() {
			if r := recover(); r != nil {
				if err, ok := r.(error); ok {
					err = errors.Cause(err)
					assert.True(t, errors.IsAborted(err), "%+v", err)
				} else {
					t.Errorf("Panic should contain an error but got:\n%+v", r)
				}
			} else {
				t.Error("Expecting a panic but got nothing")
			}
		}()
		_, _, _ = Interpolate("SELECT * FROM x WHERE a = ?").DriverValue(argValUint16(1)).ToSQL()
	})
}

func TestInterpolate_Reset(t *testing.T) {
	t.Parallel()

	t.Run("call twice with different arguments", func(t *testing.T) {
		ip := Interpolate("SELECT * FROM x WHERE a IN (?) AND b BETWEEN ? AND ? AND c = ? AND d IN (?)").
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

func TestInterpolateInt64(t *testing.T) {
	t.Parallel()

	t.Run("equal named params", func(t *testing.T) {
		str := Interpolate("SELECT * FROM x WHERE a = (:ArgX) AND b > @ArgY").
			Named(
				sql.Named(":ArgX", 3),
				sql.Named("@ArgY", 3.14159),
			).String()

		assert.Equal(t, "SELECT * FROM x WHERE a = (3) AND b > 3.14159", str)
	})

	t.Run("equal", func(t *testing.T) {
		str := Interpolate("SELECT * FROM x WHERE a = ? AND b = ? AND c = ? AND d = ? AND e = ? AND f = ? AND g = ? AND h = ? AND ab = ? AND j = ?").
			Int64s(1, -2, 3, 4, 5, 6, 7, 8, 9, 10).String()
		assert.Equal(t, "SELECT * FROM x WHERE a = 1 AND b = -2 AND c = 3 AND d = 4 AND e = 5 AND f = 6 AND g = 7 AND h = 8 AND ab = 9 AND j = 10", str)
	})
	t.Run("in", func(t *testing.T) {
		str := Interpolate("SELECT * FROM x WHERE a IN (?)").Int64s(1, -2, 3, 4, 5, 6, 7, 8, 9, 10).String()
		assert.Exactly(t, "SELECT * FROM x WHERE a IN (1,-2,3,4,5,6,7,8,9,10)", str)
	})
	t.Run("in and equal", func(t *testing.T) {
		str := Interpolate("SELECT * FROM x WHERE a = ? AND b = ? AND c = ? AND h = ? AND i = ? AND j = ? AND k = ? AND m IN (?) OR n = ?").
			Int64(1).
			Int64(-2).
			Int64(3).
			Int64(4).
			Int64(5).
			Int64(6).
			Int64(11).
			Int64s(12, 13).
			Int64(-14).String()
		assert.Exactly(t,
			`SELECT * FROM x WHERE a = 1 AND b = -2 AND c = 3 AND h = 4 AND i = 5 AND j = 6 AND k = 11 AND m IN (12,13) OR n = -14`,
			str)
	})
}

func TestInterpolateBools(t *testing.T) {
	t.Parallel()

	t.Run("single args", func(t *testing.T) {
		str := Interpolate("SELECT * FROM x WHERE a = ? AND b = ?").Bools(true, false).String()
		assert.Equal(t, "SELECT * FROM x WHERE a = 1 AND b = 0", str)
	})
	t.Run("IN args", func(t *testing.T) {
		str := Interpolate("SELECT * FROM x WHERE a IN (?) AND b = ? OR c = ?").
			Bools(true, false).Bool(true).Bool(false).String()
		assert.Equal(t, "SELECT * FROM x WHERE a IN (1,0) AND b = 1 OR c = 0", str)
	})
}

func TestInterpolate_Bytes(t *testing.T) {
	t.Parallel()

	b1 := []byte(`Go`)
	b2 := []byte(`Further`)
	t.Run("single args", func(t *testing.T) {
		str := Interpolate("SELECT * FROM x WHERE a = ? AND b = ?").BytesSlice(b1, b2).String()
		assert.Equal(t, "SELECT * FROM x WHERE a = 'Go' AND b = 'Further'", str)
	})
	t.Run("IN args", func(t *testing.T) {
		str := Interpolate("SELECT * FROM x WHERE a IN (?)").BytesSlice(b1, b2).String()
		assert.Equal(t, "SELECT * FROM x WHERE a IN ('Go','Further')", str)
	})
	t.Run("empty arg triggers not valid error", func(t *testing.T) {
		str, _, err := Interpolate("SELECT * FROM x WHERE a IN (?)").BytesSlice().ToSQL()
		assert.True(t, errors.IsNotValid(err), "%+v", err)
		assert.Equal(t, "", str)
	})
}

func TestInterpolate_Time(t *testing.T) {
	t.Parallel()

	t1 := now()
	t2 := now().Add(time.Minute)

	t.Run("single args", func(t *testing.T) {
		str := Interpolate("SELECT * FROM x WHERE a = ? AND b = ?").Times(t1, t2).String()
		assert.Equal(t, "SELECT * FROM x WHERE a = '2006-01-02 15:04:05' AND b = '2006-01-02 15:05:05'", str)
	})
	t.Run("IN args", func(t *testing.T) {
		str := Interpolate("SELECT * FROM x WHERE a IN (?)").Times(t1, t2).String()
		assert.Equal(t, "SELECT * FROM x WHERE a IN ('2006-01-02 15:04:05','2006-01-02 15:05:05')", str)
	})
	t.Run("empty arg", func(t *testing.T) {
		str, _, err := Interpolate("SELECT * FROM x WHERE a IN (?)").Times().ToSQL()
		assert.True(t, errors.IsNotValid(err), "%+v", err)
		assert.Empty(t, str)
	})
}

func TestInterpolateFloats(t *testing.T) {
	t.Parallel()
	t.Run("single args", func(t *testing.T) {
		str := Interpolate("SELECT * FROM x WHERE a = ? AND b = ?").Float64s(3.14159, 2.7182818).String()
		assert.Equal(t, "SELECT * FROM x WHERE a = 3.14159 AND b = 2.7182818", str)
	})
	t.Run("IN args", func(t *testing.T) {
		str := Interpolate("SELECT * FROM x WHERE a IN (?) AND b = ?").Float64s(3.14159, 2.7182818).Float64(0.815).String()
		assert.Equal(t, "SELECT * FROM x WHERE a IN (3.14159,2.7182818) AND b = 0.815", str)
	})
	t.Run("IN args reverse", func(t *testing.T) {
		str := Interpolate("SELECT * FROM x WHERE b = ? AND a IN (?)").Float64(0.815).Float64s(3.14159, 2.7182818).String()
		assert.Equal(t, "SELECT * FROM x WHERE b = 0.815 AND a IN (3.14159,2.7182818)", str)
	})
}

func TestInterpolateBetween(t *testing.T) {
	t.Parallel()

	runner := func(placeHolderStr, interpolatedStr string, wantErr errors.BehaviourFunc, args Arguments) func(*testing.T) {
		return func(t *testing.T) {
			have, _, err := Interpolate(placeHolderStr).ArgUnions(args).ToSQL()
			if wantErr != nil {
				assert.True(t, wantErr(err))
				return
			}
			require.NoError(t, err, t.Name())
			assert.Exactly(t, interpolatedStr, have, t.Name())
		}
	}
	t.Run("BETWEEN at the end", runner(
		"SELECT * FROM x WHERE a IN (?) AND b IN (?) AND c NOT IN (?) AND d BETWEEN ? AND ?",
		"SELECT * FROM x WHERE a IN (1) AND b IN (1,2,3) AND c NOT IN (5,6,7) AND d BETWEEN 'wat' AND 'ok'",
		nil,
		MakeArgs(5).Ints(1).Ints(1, 2, 3).Int64s(5, 6, 7).String("wat").String("ok"),
	))
	t.Run("BETWEEN in the middle", runner(
		"SELECT * FROM x WHERE a IN (?) AND b IN (?) AND d BETWEEN ? AND ? AND c NOT IN (?)",
		"SELECT * FROM x WHERE a IN (1) AND b IN (1,2,3) AND d BETWEEN 'wat' AND 'ok' AND c NOT IN (5,6,7)",
		nil,
		MakeArgs(5).Ints(1).Ints(1, 2, 3).String("wat").String("ok").Int64s(5, 6, 7),
	))
}

func TestInterpolateStrings(t *testing.T) {
	t.Parallel()
	t.Run("single args", func(t *testing.T) {
		str := Interpolate("SELECT * FROM x WHERE a = ? AND b = ? AND c = ?").Str("a'b").Str("c`d").Str("\"hello's \\ world\" \n\r\x00\x1a").String()
		assert.Equal(t, "SELECT * FROM x WHERE a = 'a\\'b' AND b = 'c`d' AND c = '\\\"hello\\'s \\\\ world\\\" \\n\\r\\x00\\x1a'", str)
	})
	t.Run("IN args", func(t *testing.T) {
		str := Interpolate("SELECT * FROM x WHERE a IN (?) AND b = ?").Strs("a'b", "c`d").Strs("1' or '1' = '1'))/*").String()
		assert.Equal(t, "SELECT * FROM x WHERE a IN ('a\\'b','c`d') AND b = '1\\' or \\'1\\' = \\'1\\'))/*'", str)
	})
	t.Run("empty args triggers incorrect interpolation of values", func(t *testing.T) {
		var sl = make([]string, 0, 2)
		str := Interpolate("SELECT * FROM x WHERE a IN (?) AND b = ? OR c = ?").Strs("a", "b").Str("c").Strs(sl...).String()
		assert.Exactly(t, "SELECT * FROM x WHERE a IN ('a') AND b = 'b' OR c = 'c'", str)
	})
}

func TestInterpolateSlices(t *testing.T) {
	t.Parallel()
	str := Interpolate("SELECT * FROM x WHERE a = (?) AND b = (?) AND c = (?) AND d = (?) AND e = ?").
		Ints(1).
		Ints(1, 2, 3).
		Int64s(5, 6, 7).
		Strs("wat", "ok").
		Int(8).String()
	assert.Equal(t, "SELECT * FROM x WHERE a = (1) AND b = (1,2,3) AND c = (5,6,7) AND d = ('wat','ok') AND e = 8", str)
}

func TestInterpolate(t *testing.T) {
	t.Parallel()
	tests := []struct {
		sql string
		Arguments
		expSQL string
		errBhf errors.BehaviourFunc
	}{
		// NULL
		{"SELECT * FROM x WHERE a = ?", MakeArgs(1).Null(),
			"SELECT * FROM x WHERE a = NULL", nil},

		// integers
		{
			`SELECT * FROM x WHERE a = ? AND b = ? AND c = ? AND d = ? AND e = ? AND f = ?
			AND g = ? AND h = ? AND i = ? AND j = ?`,
			MakeArgs(1).Ints(1, -2, 3, 4, 5, 6, 7, 8, 9, 10),
			`SELECT * FROM x WHERE a = 1 AND b = -2 AND c = 3 AND d = 4 AND e = 5 AND f = 6
			AND g = 7 AND h = 8 AND i = 9 AND j = 10`, nil,
		},
		{
			`SELECT * FROM x WHERE a IN (?)`,
			MakeArgs(1).Ints(1, -2, 3, 4, 5, 6, 7, 8, 9, 10),
			`SELECT * FROM x WHERE a IN (1,-2,3,4,5,6,7,8,9,10)`, nil,
		},
		{
			`SELECT * FROM x WHERE a IN (?)`,
			MakeArgs(1).Int64s(1, -2, 3, 4, 5, 6, 7, 8, 9, 10),
			`SELECT * FROM x WHERE a IN (1,-2,3,4,5,6,7,8,9,10)`, nil,
		},

		// boolean
		{"SELECT * FROM x WHERE a = ? AND b = ?", MakeArgs(2).Bool(true).Bool(false),
			"SELECT * FROM x WHERE a = 1 AND b = 0", nil},
		{"SELECT * FROM x WHERE a = ? AND b = ?", MakeArgs(1).Bools(true, false),
			"SELECT * FROM x WHERE a = 1 AND b = 0", nil},

		// floats
		{"SELECT * FROM x WHERE a = ? AND b = ?", MakeArgs(2).Float64(0.15625).Float64(3.14159),
			"SELECT * FROM x WHERE a = 0.15625 AND b = 3.14159", nil},
		{"SELECT * FROM x WHERE a = ? AND b = ?", MakeArgs(1).Float64s(0.15625, 3.14159),
			"SELECT * FROM x WHERE a = 0.15625 AND b = 3.14159", nil},
		{"SELECT * FROM x WHERE a = ? AND b = ? and C = ?", MakeArgs(1).Float64s(0.15625, 3.14159),
			"", errors.IsNotValid},
		{
			`SELECT * FROM x WHERE a IN (?)`,
			MakeArgs(1).Float64s(1.1, -2.2, 3.3),
			`SELECT * FROM x WHERE a IN (1.1,-2.2,3.3)`, nil,
		},

		// strings
		{
			`SELECT * FROM x WHERE a = ?
			AND b = ?`,
			MakeArgs(2).String("hello").String("\"hello's \\ world\" \n\r\x00\x1a"),
			`SELECT * FROM x WHERE a = 'hello'
			AND b = '\"hello\'s \\ world\" \n\r\x00\x1a'`, nil,
		},
		{
			`SELECT * FROM x WHERE a IN (?)`,
			MakeArgs(1).Strings("a'a", "bb"),
			`SELECT * FROM x WHERE a IN ('a\'a','bb')`, nil,
		},
		{
			`SELECT * FROM x WHERE a IN (?,?)`,
			MakeArgs(1).Strings("a'a", "bb"),
			`SELECT * FROM x WHERE a IN ('a\'a','bb')`, nil,
		},

		// slices
		{"SELECT * FROM x WHERE a = ? AND b = (?) AND c = (?) AND d = (?)",
			MakeArgs(4).Int(1).Ints(1, 2, 3).Ints(5, 6, 7).Strings("wat", "ok"),
			"SELECT * FROM x WHERE a = 1 AND b = (1,2,3) AND c = (5,6,7) AND d = ('wat','ok')", nil},
		//
		////// TODO valuers
		////{"SELECT * FROM x WHERE a = ? AND b = ?",
		////	args{myString{true, "wat"}, myString{false, "fry"}},
		////	"SELECT * FROM x WHERE a = 'wat' AND b = NULL", nil},

		// errors
		{"SELECT * FROM x WHERE a = ? AND b = ?", MakeArgs(1).Int64(1),
			"", errors.IsNotValid},

		{"SELECT * FROM x WHERE", MakeArgs(1).Int(1),
			"", errors.IsNotValid},

		{"SELECT * FROM x WHERE a = ?", MakeArgs(1).String(string([]byte{0x34, 0xFF, 0xFE})),
			"", errors.IsNotValid},

		// String() without arguments is equal to empty interface in the previous version.
		{"SELECT 'hello", MakeArgs(1).String(""), "", errors.IsNotValid},
		{`SELECT "hello`, MakeArgs(1).String(""), "", errors.IsNotValid},

		// preprocessing
		{"SELECT '?'", nil, "SELECT '?'", nil},
		{"SELECT `?`", nil, "SELECT `?`", nil},
		{"SELECT [?]", nil, "SELECT `?`", nil},
		{"SELECT [?]", nil, "SELECT `?`", nil},
		{"SELECT [name] FROM [user]", nil, "SELECT `name` FROM `user`", nil},
		{"SELECT [u.name] FROM [user] [u]", nil, "SELECT `u`.`name` FROM `user` `u`", nil},
		{"SELECT [u.na`me] FROM [user] [u]", nil, "SELECT `u`.`na``me` FROM `user` `u`", nil},
		{"SELECT * FROM [user] WHERE [name] = '[nick]'", nil,
			"SELECT * FROM `user` WHERE `name` = '[nick]'", nil},
		{`SELECT * FROM [user] WHERE [name] = "nick[]"`, nil,
			"SELECT * FROM `user` WHERE `name` = 'nick[]'", nil},
	}

	for i, test := range tests {
		str, _, err := Interpolate(test.sql).ArgUnions(test.Arguments).ToSQL()
		if test.errBhf != nil {
			if !test.errBhf(err) {
				t.Errorf("IDX %d\ngot error: %v\nwant: %s", i, err, test.errBhf(err))
			}
		}
		//assert.NoError(t, err, "IDX %d", i)
		if str != test.expSQL {
			t.Errorf("IDX %d\ngot: %v\nwant: %v\nError: %+v", i, str, test.expSQL, err)
		}
	}
}
