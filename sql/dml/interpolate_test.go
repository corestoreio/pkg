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
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var _ fmt.Stringer = (*ip)(nil)
var _ QueryBuilder = (*ip)(nil)

func TestRepeat(t *testing.T) {
	t.Parallel()

	t.Run("MisMatch", func(t *testing.T) {
		s, err := ExpandPlaceHolders("SELECT * FROM `table` WHERE id IN (?)", nil)
		assert.Empty(t, s)
		assert.True(t, errors.IsMismatch(err), "%+v", err)
	})
	t.Run("MisMatch length reps", func(t *testing.T) {
		s, err := ExpandPlaceHolders("SELECT * FROM `table` WHERE id IN (?)", MakeArgs(2).Int64s(1, 2).Strings("d", "3"))
		assert.Empty(t, s)
		assert.True(t, errors.IsMismatch(err), "%+v", err)
	})
	t.Run("MisMatch qMarks", func(t *testing.T) {
		s, err := ExpandPlaceHolders("SELECT * FROM `table` WHERE id IN(!)", MakeArgs(1).Int64(3))
		assert.Empty(t, s)
		assert.True(t, errors.IsMismatch(err), "%+v", err)
	})
	t.Run("one arg with one value", func(t *testing.T) {
		args := MakeArgs(1).Int64(1)
		s, err := ExpandPlaceHolders("SELECT * FROM `table` WHERE id IN (?)", args)
		require.NoError(t, err)
		assert.Exactly(t, "SELECT * FROM `table` WHERE id IN (?)", s)
		assert.Exactly(t, []interface{}{int64(1)}, args.Interfaces())
	})
	t.Run("one arg with three values", func(t *testing.T) {
		args := MakeArgs(1).Int64s(11, 3, 5)
		s, err := ExpandPlaceHolders("SELECT * FROM `table` WHERE id IN ?", args)
		require.NoError(t, err)
		assert.Exactly(t, "SELECT * FROM `table` WHERE id IN (?,?,?)", s)
		assert.Exactly(t, []interface{}{int64(11), int64(3), int64(5)}, args.Interfaces())
	})
	t.Run("multi 3,5 times replacement", func(t *testing.T) {
		args := MakeArgs(3).Int64s(5, 7, 9).Strings("a", "b", "c", "d", "e")
		s, err := ExpandPlaceHolders("SELECT * FROM `table` WHERE id IN ? AND name IN ?", args)
		require.NoError(t, err)
		assert.Exactly(t, "SELECT * FROM `table` WHERE id IN (?,?,?) AND name IN (?,?,?,?,?)", s)
		assert.Exactly(t, []interface{}{int64(5), int64(7), int64(9), "a", "b", "c", "d", "e"}, args.Interfaces())
	})
}

func TestInterpolate_Nil(t *testing.T) {
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

func TestInterpolate_Errors(t *testing.T) {
	t.Parallel()
	t.Run("non utf8", func(t *testing.T) {
		compareToSQL2(t,
			Interpolate("SELECT * FROM x WHERE a = ?").Str(string([]byte{0x34, 0xFF, 0xFE})),
			errors.IsNotValid,
			"",
		)
	})
	t.Run("too few qmarks", func(t *testing.T) {
		compareToSQL2(t,
			Interpolate("SELECT * FROM x WHERE a = ?").Int(3).Int(4),
			errors.IsMismatch,
			"",
		)
	})
	t.Run("too many qmarks", func(t *testing.T) {
		compareToSQL2(t,
			Interpolate("SELECT * FROM x WHERE a = ? OR b = ? or c = ?").Ints(3, 4),
			errors.IsMismatch,
			"",
		)
	})
	t.Run("way too many qmarks", func(t *testing.T) {
		_, _, err := Interpolate("SELECT * FROM x WHERE a IN ? OR b = ? OR c = ? AND d = ?").Ints(3, 4).Int64(2).ToSQL()
		assert.True(t, errors.IsMismatch(err), "%+v", err)
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
		compareToSQL2(t,
			Interpolate("SELECT * FROM x WHERE a = ? AND b = ? AND c = ? AND d = ? AND e = ? AND f = ? AND g = ? AND h = ?").
				DriverValue(aInt).
				DriverValue(aStr).
				DriverValue(aFlo).
				DriverValue(aTim).
				DriverValue(aBoo).
				DriverValue(aByt).
				DriverValue(aNil).
				DriverValue(aNil),
			nil,
			"SELECT * FROM x WHERE a = (4711) AND b = ('Goph\\'er') AND c = (2.7182818) AND d = ('2006-01-02 19:04:05') AND e = (1) AND f = ('BytyGophe\\'r') AND g = (NULL) AND h = (NULL)",
		)
	})
	t.Run("in", func(t *testing.T) {
		compareToSQL2(t,
			Interpolate("SELECT * FROM x WHERE a IN ? AND b IN ? AND c IN ? AND d IN ? AND e IN ? AND f IN ?").
				DriverValue(aInt, aInt).DriverValue(aStr, aStr).DriverValue(aFlo, aFlo).
				DriverValue(aTim, aTim).DriverValue(aBoo, aBoo).DriverValue(aByt, aByt),
			nil,
			"SELECT * FROM x WHERE a IN (4711,4711) AND b IN ('Goph\\'er','Goph\\'er') AND c IN (2.7182818,2.7182818) AND d IN ('2006-01-02 19:04:05','2006-01-02 19:04:05') AND e IN (1,1) AND f IN ('BytyGophe\\'r','BytyGophe\\'r')",
		)
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
			nil,
			"SELECT * FROM x WHERE a = (3) AND b > 3.14159",
		)
	})
	t.Run("equal", func(t *testing.T) {
		compareToSQL2(t,
			Interpolate("SELECT * FROM x WHERE a = ? AND b = ? AND c = ?").
				Int64(1).Int64(-2).Int64(3),
			nil,
			"SELECT * FROM x WHERE a = 1 AND b = -2 AND c = 3",
		)
	})
	t.Run("in", func(t *testing.T) {
		compareToSQL2(t,
			Interpolate("SELECT * FROM x WHERE a IN ?").Int64s(1, -2, 3, 4, 5, 6, 7, 8, 9, 10),
			nil,
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
			nil,
			`SELECT * FROM x WHERE a = 1 AND b = -2 AND c = 3 AND h = 4 AND i = 5 AND j = 6 AND k = 11 AND m IN (12,13) OR n = -14`,
		)
	})
}

func TestInterpolate_Bools(t *testing.T) {
	t.Parallel()

	t.Run("single args", func(t *testing.T) {
		compareToSQL2(t,
			Interpolate("SELECT * FROM x WHERE a = ? AND b = ?").Bool(true).Bool(false),
			nil,
			"SELECT * FROM x WHERE a = 1 AND b = 0",
		)
	})
	t.Run("IN args", func(t *testing.T) {
		compareToSQL2(t,
			Interpolate("SELECT * FROM x WHERE a IN ? AND b = ? OR c = ?").
				Bools(true, false).Bool(true).Bool(false),
			nil,
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
			nil,
			"SELECT * FROM x WHERE a = 'Go' AND b = 'Further'",
		)
	})
	t.Run("UTF8 valid: IN args", func(t *testing.T) {
		compareToSQL2(t,
			Interpolate("SELECT * FROM x WHERE a IN ?").BytesSlice(b1, b2),
			nil,
			"SELECT * FROM x WHERE a IN ('Go','Further')",
		)
	})
	t.Run("empty arg triggers no error", func(t *testing.T) {
		compareToSQL2(t,
			Interpolate("SELECT * FROM x WHERE a IN ?").BytesSlice(),
			nil,
			"SELECT * FROM x WHERE a IN ()",
		)
	})
	t.Run("Binary to hex", func(t *testing.T) {
		bin := []byte{66, 250, 67}
		compareToSQL2(t,
			Interpolate("SELECT * FROM x WHERE a = ?").Bytes(bin),
			nil,
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
			nil,
			"SELECT * FROM x WHERE a = '2006-01-02 15:04:05' AND b = '2006-01-02 15:05:05'",
		)
	})
	t.Run("IN args", func(t *testing.T) {
		compareToSQL2(t,
			Interpolate("SELECT * FROM x WHERE a IN ?").Times(t1, t2),
			nil,
			"SELECT * FROM x WHERE a IN ('2006-01-02 15:04:05','2006-01-02 15:05:05')",
		)
	})
	t.Run("empty arg", func(t *testing.T) {
		compareToSQL2(t,
			Interpolate("SELECT * FROM x WHERE a IN ? ?").Times(),
			errors.IsMismatch,
			"",
		)
	})
}

func TestInterpolate_Floats(t *testing.T) {
	t.Parallel()
	t.Run("single args", func(t *testing.T) {
		compareToSQL2(t,
			Interpolate("SELECT * FROM x WHERE a = ? AND b = ?").Float64(3.14159).Float64(2.7182818),
			nil,
			"SELECT * FROM x WHERE a = 3.14159 AND b = 2.7182818",
		)
	})
	t.Run("IN args", func(t *testing.T) {
		compareToSQL2(t,
			Interpolate("SELECT * FROM x WHERE a IN ? AND b = ?").Float64s(3.14159, 2.7182818).Float64(0.815),
			nil,
			"SELECT * FROM x WHERE a IN (3.14159,2.7182818) AND b = 0.815",
		)
	})
	t.Run("IN args reverse", func(t *testing.T) {
		compareToSQL2(t,
			Interpolate("SELECT * FROM x WHERE b = ? AND a IN ?").Float64(0.815).Float64s(3.14159, 2.7182818),
			nil,
			"SELECT * FROM x WHERE b = 0.815 AND a IN (3.14159,2.7182818)",
		)
	})
}

func TestInterpolate_Slices_Strings_Between(t *testing.T) {
	t.Parallel()

	t.Run("BETWEEN at the end", func(t *testing.T) {
		compareToSQL2(t,
			Interpolate("SELECT * FROM x WHERE a IN ? AND b IN ? AND c NOT IN ? AND d BETWEEN ? AND ?").
				ArgUnions(MakeArgs(5).Ints(1).Ints(1, 2, 3).Int64s(5, 6, 7).String("wat").String("ok")),
			nil,
			"SELECT * FROM x WHERE a IN (1) AND b IN (1,2,3) AND c NOT IN (5,6,7) AND d BETWEEN 'wat' AND 'ok'",
		)
	})
	t.Run("BETWEEN in the middle", func(t *testing.T) {
		compareToSQL2(t,
			Interpolate("SELECT * FROM x WHERE a IN ? AND b IN ? AND d BETWEEN ? AND ? AND c NOT IN ?").
				ArgUnions(MakeArgs(5).Ints(1).Ints(1, 2, 3).String("wat").String("ok").Int64s(5, 6, 7)),
			nil,
			"SELECT * FROM x WHERE a IN (1) AND b IN (1,2,3) AND d BETWEEN 'wat' AND 'ok' AND c NOT IN (5,6,7)",
		)
	})

	t.Run("single args", func(t *testing.T) {
		compareToSQL2(t,
			Interpolate("SELECT * FROM x WHERE a = ? AND b = ? AND c = ?").Str("a'b").Str("c`d").Str("\"hello's \\ world\" \n\r\x00\x1a"),
			nil,
			"SELECT * FROM x WHERE a = 'a\\'b' AND b = 'c`d' AND c = '\\\"hello\\'s \\\\ world\\\" \\n\\r\\x00\\x1a'",
		)
	})
	t.Run("IN args", func(t *testing.T) {
		compareToSQL2(t,
			Interpolate("SELECT * FROM x WHERE a IN ? AND b = ?").Strs("a'b", "c`d").Strs("1' or '1' = '1'))/*"),
			nil,
			"SELECT * FROM x WHERE a IN ('a\\'b','c`d') AND b = ('1\\' or \\'1\\' = \\'1\\'))/*')",
		)
	})
	t.Run("empty args triggers incorrect interpolation of values", func(t *testing.T) {
		var sl = make([]string, 0, 2)
		compareToSQL2(t,
			Interpolate("SELECT * FROM x WHERE a IN ? AND b = ? OR c = ?").Strs("a", "b").Str("c").Strs(sl...),
			nil,
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
			nil,
			"SELECT * FROM x WHERE a = (1) AND b = (1,2,3) AND c = (5,6,7) AND d = ('wat','ok') AND e = 8",
		)
	})
}

func TestInterpolate_Different(t *testing.T) {
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
			`SELECT * FROM x WHERE a = ? AND b = ?`,
			MakeArgs(1).Int(1).Int(-2),
			`SELECT * FROM x WHERE a = 1 AND b = -2`, nil,
		},
		{
			`SELECT * FROM x WHERE a IN ?`,
			MakeArgs(1).Ints(1, -2, 3, 4, 5, 6, 7, 8, 9, 10),
			`SELECT * FROM x WHERE a IN (1,-2,3,4,5,6,7,8,9,10)`, nil,
		},
		{
			`SELECT * FROM x WHERE a IN ?`,
			MakeArgs(1).Int64s(1, -2, 3, 4, 5, 6, 7, 8, 9, 10),
			`SELECT * FROM x WHERE a IN (1,-2,3,4,5,6,7,8,9,10)`, nil,
		},

		// boolean
		{"SELECT * FROM x WHERE a = ? AND b = ?", MakeArgs(2).Bool(true).Bool(false),
			"SELECT * FROM x WHERE a = 1 AND b = 0", nil},
		{"SELECT * FROM x WHERE a IN ?", MakeArgs(1).Bools(true, false),
			"SELECT * FROM x WHERE a IN (1,0)", nil},

		// floats
		{"SELECT * FROM x WHERE a = ? AND b = ?", MakeArgs(2).Float64(0.15625).Float64(3.14159),
			"SELECT * FROM x WHERE a = 0.15625 AND b = 3.14159", nil},
		{"SELECT * FROM x WHERE a IN ?", MakeArgs(1).Float64s(0.15625, 3.14159),
			"SELECT * FROM x WHERE a IN (0.15625,3.14159)", nil},
		{"SELECT * FROM x WHERE a = ? AND b = ? and C = ?", MakeArgs(1).Float64s(0.15625, 3.14159),
			"", errors.IsMismatch},
		{
			`SELECT * FROM x WHERE a IN ?`,
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
			`SELECT * FROM x WHERE a IN ?`,
			MakeArgs(1).Strings("a'a", "bb"),
			`SELECT * FROM x WHERE a IN ('a\'a','bb')`, nil,
		},
		{
			`SELECT * FROM x WHERE a IN ?`,
			MakeArgs(1).Strings("a'a", "bb"),
			`SELECT * FROM x WHERE a IN ('a\'a','bb')`, nil,
		},

		// slices
		{"SELECT * FROM x WHERE a = ? AND b = ? AND c = ? AND d = ?",
			MakeArgs(4).Int(1).Ints(1, 2, 3).Ints(5, 6, 7).Strings("wat", "ok"),
			"SELECT * FROM x WHERE a = 1 AND b = (1,2,3) AND c = (5,6,7) AND d = ('wat','ok')", nil},
		//
		////// TODO valuers
		////{"SELECT * FROM x WHERE a = ? AND b = ?",
		////	args{myString{true, "wat"}, myString{false, "fry"}},
		////	"SELECT * FROM x WHERE a = 'wat' AND b = NULL", nil},

		// errors
		{"SELECT * FROM x WHERE a = ? AND b = ?", MakeArgs(1).Int64(1),
			"", errors.IsMismatch},

		{"SELECT * FROM x WHERE", MakeArgs(1).Int(1),
			"", errors.IsMismatch},

		{"SELECT * FROM x WHERE a = ?", MakeArgs(1).String(string([]byte{0x34, 0xFF, 0xFE})),
			"", errors.IsNotValid},

		// String() without arguments is equal to empty interface in the previous version.
		{"SELECT 'hello", MakeArgs(1).String(""), "", errors.IsMismatch},
		{`SELECT "hello`, MakeArgs(1).String(""), "", errors.IsMismatch},
		{`SELECT ? "hello`, MakeArgs(1).String(""), "SELECT '' 'hello", nil},

		// preprocessing
		{"SELECT '?'", nil, "SELECT '?'", nil},
		{"SELECT `?`", nil, "SELECT `?`", nil},
		{"SELECT [?]", nil, "SELECT `?`", nil},
		{"SELECT [?]", nil, "SELECT `?`", nil},
		{"SELECT [name] FROM [user]", nil, "SELECT `name` FROM `user`", nil},
		{"SELECT [u.name] FROM [user] [u]", nil, "SELECT `u`.`name` FROM `user` `u`", nil},
		{"SELECT [u.na`me] FROM [user] [u]", nil, "SELECT `u`.`na``me` FROM `user` `u`", nil},
		{"SELECT * FROM [user] WHERE [name] = '[nick]'", nil, "SELECT * FROM `user` WHERE `name` = '[nick]'", nil},
		{`SELECT * FROM [user] WHERE [name] = "nick[]"`, nil, "SELECT * FROM `user` WHERE `name` = 'nick[]'", nil},
		{"SELECT * FROM [user] WHERE [name] = 'Hello`s World'", nil, "SELECT * FROM `user` WHERE `name` = 'Hello`s World'", nil},
	}

	for _, test := range tests {
		str, _, err := Interpolate(test.sql).ArgUnions(test.Arguments).ToSQL()
		if test.errBhf != nil {
			if !test.errBhf(err) {
				isErr := test.errBhf(err)
				t.Errorf("IDX %q\ngot error: %v\nwant: %t", test.sql, err, isErr)
			}
		}
		//assert.NoError(t, err, "IDX %d", i)
		if str != test.expSQL {
			t.Errorf("IDX %q\ngot: %v\nwant: %v\nError: %+v", test.sql, str, test.expSQL, err)
		}
	}
}

func TestInterpolate_MultipleSingleQuotes(t *testing.T) {
	const rawSQL = "DELETE FROM `tableA` WHERE (`colA` >= 3.14159) AND (`colB` IN ?) AND (`colC` = 'He\\'llo') ORDER BY `id` LIMIT 10"
	str, args, err := Interpolate(rawSQL).ArgUnions(MakeArgs(1).Float64s(3.1, 2.4)).ToSQL()
	assert.NoError(t, err)
	assert.Nil(t, args)
	assert.Exactly(t, "DELETE FROM `tableA` WHERE (`colA` >= 3.14159) AND (`colB` IN (3.1,2.4)) AND (`colC` = 'He\\'llo') ORDER BY `id` LIMIT 10", str)
}

func TestExtractNamedArgs(t *testing.T) {
	t.Parallel()

	runner := func(haveSQL, wantSQL string, wantQualifiedColumns ...string) func(*testing.T) {
		return func(t *testing.T) {
			gotSQL, qualifiedColumns := extractReplaceNamedArgs(haveSQL, nil)
			assert.Exactly(t, wantSQL, gotSQL)
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
