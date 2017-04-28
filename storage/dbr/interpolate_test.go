package dbr

import (
	"testing"

	"github.com/corestoreio/errors"
	"github.com/stretchr/testify/assert"
)

func TestRepeat(t *testing.T) {
	t.Parallel()
	t.Run("MisMatch", func(t *testing.T) {
		s, args, err := Repeat("SELECT * FROM `table` WHERE id IN (?)")
		assert.Empty(t, s)
		assert.Nil(t, args)
		assert.True(t, errors.IsMismatch(err), "%+v", err)
	})
	t.Run("MisMatch length reps", func(t *testing.T) {
		s, args, err := Repeat("SELECT * FROM `table` WHERE id IN (?)", ArgInt(1, 2), ArgString("d", "3"))
		assert.Empty(t, s)
		assert.Nil(t, args)
		assert.True(t, errors.IsMismatch(err), "%+v", err)
	})
	t.Run("MisMatch qMarks", func(t *testing.T) {
		s, args, err := Repeat("SELECT * FROM `table` WHERE id IN(!)", argInt(3))
		assert.Empty(t, s)
		assert.Nil(t, args)
		assert.True(t, errors.IsMismatch(err), "%+v", err)
	})
	t.Run("one arg with one value", func(t *testing.T) {
		s, args, err := Repeat("SELECT * FROM `table` WHERE id IN (?)", argInt(1))
		assert.Exactly(t, "SELECT * FROM `table` WHERE id IN (?)", s)
		assert.Exactly(t, []interface{}{int64(1)}, args)
		assert.NoError(t, err, "%+v", err)
	})
	t.Run("one arg with three values", func(t *testing.T) {
		s, args, err := Repeat("SELECT * FROM `table` WHERE id IN (?)", ArgInt(11, 3, 5))
		assert.Exactly(t, "SELECT * FROM `table` WHERE id IN (?,?,?)", s)
		assert.Exactly(t, []interface{}{int64(11), int64(3), int64(5)}, args)
		assert.NoError(t, err, "%+v", err)
	})
	t.Run("multi 3,5 times replacement", func(t *testing.T) {
		sl := []string{"a", "b", "c", "d", "e"}
		s, args, err := Repeat("SELECT * FROM `table` WHERE id IN (?) AND name IN (?)",
			ArgInt(5, 7, 9), ArgString(sl...))
		assert.Exactly(t, "SELECT * FROM `table` WHERE id IN (?,?,?) AND name IN (?,?,?,?,?)", s)
		assert.Exactly(t, []interface{}{int64(5), int64(7), int64(9), "a", "b", "c", "d", "e"}, args)
		assert.NoError(t, err, "%+v", err)
	})
}

//BenchmarkRepeat/multi-4         	 3000000	       492 ns/op	      96 B/op	       1 allocs/op no iFace wrapping
//BenchmarkRepeat/single-4        	 5000000	       311 ns/op	      48 B/op	       1 allocs/op no iFace wrapping

//BenchmarkRepeat/multi-4         	 1000000	      1753 ns/op	    1192 B/op	      19 allocs/op
//BenchmarkRepeat/single-4        	 2000000	       899 ns/op	     448 B/op	      11 allocs/op

func BenchmarkRepeat(b *testing.B) {

	b.Run("multi", func(b *testing.B) {
		sl := []string{"a", "b", "c", "d", "e"}
		const want = "SELECT * FROM `table` WHERE id IN (?,?,?,?) AND name IN (?,?,?,?,?) AND status IN (?)"
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			s, args, err := Repeat("SELECT * FROM `table` WHERE id IN (?) AND name IN (?) AND status IN (?)",
				ArgInt(5, 7, 9, 11), ArgString(sl...), argInt(22))
			if err != nil {
				b.Fatalf("%+v", err)
			}
			if s != want {
				b.Fatalf("\nHave: %q\nWant: %q", s, want)
			}
			if len(args) == 0 {
				b.Fatal("Args cannot be empty")
			}
		}
	})

	b.Run("single", func(b *testing.B) {
		const want = "SELECT * FROM `table` WHERE id IN (?,?,?,?)"
		for i := 0; i < b.N; i++ {
			s, args, err := Repeat("SELECT * FROM `table` WHERE id IN (?)", ArgInt(9, 8, 7, 6))
			if err != nil {
				b.Fatalf("%+v", err)
			}
			if s != want {
				b.Fatalf("\nHave: %q\nWant: %q", s, want)
			}
			if len(args) == 0 {
				b.Fatal("Args cannot be empty")
			}
		}
	})
}

func TestInterpolateNil(t *testing.T) {
	t.Parallel()
	str, err := Interpolate("SELECT * FROM x WHERE a = ?", ArgNull())
	assert.NoError(t, err)
	assert.Equal(t, "SELECT * FROM x WHERE a = NULL", str)
}

func TestInterpolateErrors(t *testing.T) {
	t.Parallel()
	t.Run("non utf8", func(t *testing.T) {
		_, err := Interpolate("SELECT * FROM x WHERE a = ?", ArgString(string([]byte{0x34, 0xFF, 0xFE})))
		assert.True(t, errors.IsNotValid(err), "%+v", err)
	})
	t.Run("too few qmarks", func(t *testing.T) {
		str, err := Interpolate("SELECT * FROM x WHERE a = ?",
			ArgInt(3, 4),
		)
		assert.Empty(t, str)
		assert.True(t, errors.IsNotValid(err), "%+v", err)
	})
	t.Run("too many qmarks", func(t *testing.T) {
		str, err := Interpolate("SELECT * FROM x WHERE a = ? OR b = ? or c = ?",
			ArgInt(3, 4),
		)
		assert.Empty(t, str)
		assert.True(t, errors.IsNotValid(err), "%+v", err)
	})
	t.Run("way too many qmarks", func(t *testing.T) {
		str, err := Interpolate("SELECT * FROM x WHERE a IN ? OR b = ? OR c = ? AND d = ?",
			ArgInt(3, 4).Operator(In),
			argInt64(2),
		)
		assert.Empty(t, str)
		assert.True(t, errors.IsNotValid(err), "%+v", err)
	})
}

func TestInterpolateInt64(t *testing.T) {
	t.Parallel()
	t.Run("each", func(t *testing.T) {
		str, err := Interpolate("SELECT * FROM x WHERE a = ? AND b = ? AND c = ? AND d = ? AND e = ? AND f = ? AND g = ? AND h = ? AND ab = ? AND j = ?",
			ArgInt64(1, -2, 3, 4, 5, 6, 7, 8, 9, 10),
		)
		assert.NoError(t, err)
		assert.Equal(t, "SELECT * FROM x WHERE a = 1 AND b = -2 AND c = 3 AND d = 4 AND e = 5 AND f = 6 AND g = 7 AND h = 8 AND ab = 9 AND j = 10", str)
	})
	t.Run("in", func(t *testing.T) {
		str, err := Interpolate("SELECT * FROM x WHERE a IN ?",
			ArgInt64(1, -2, 3, 4, 5, 6, 7, 8, 9, 10).Operator(In),
		)
		assert.NoError(t, err)
		assert.Exactly(t,
			"SELECT * FROM x WHERE a IN (1,-2,3,4,5,6,7,8,9,10)",
			str)
	})
	t.Run("in and each", func(t *testing.T) {
		str, err := Interpolate("SELECT * FROM x WHERE a = ? AND b = ? AND c = ? AND h = ? AND i = ? AND j = ? AND k = ? AND m IN ? OR n = ?",
			ArgInt64(1, -2, 3, 4, 5, 6),
			argInt64(11),
			ArgInt64(12, 13).Operator(In),
			ArgInt64(-14),
		)
		assert.NoError(t, err)
		assert.Exactly(t,
			`SELECT * FROM x WHERE a = 1 AND b = -2 AND c = 3 AND h = 4 AND i = 5 AND j = 6 AND k = 11 AND m IN (12,13) OR n = -14`,
			str)
	})
	t.Run("empty arg", func(t *testing.T) {
		str, err := Interpolate("SELECT * FROM x WHERE a = ? AND b = ? AND c = ? AND h = ? AND i = ? AND j = ? AND k = ? AND m IN ? OR n = ?",
			ArgInt64(1, -2, 3, 4, 5, 6),
			argInt64(11),
			ArgInt64(12, 13).Operator(In),
			ArgInt64(),
		)
		assert.Empty(t, str)
		assert.True(t, errors.IsEmpty(err), "%+v", err)
	})
}

var preprocessSink string

// BenchmarkPreprocess-4   	  500000	      4013 ns/op	     174 B/op	      11 allocs/op with reflection
// BenchmarkPreprocess-4   	  500000	      3591 ns/op	     174 B/op	      11 allocs/op
func BenchmarkPreprocess(b *testing.B) {
	const want = `SELECT * FROM x WHERE a = 1 AND b = -2 AND c = 3 AND d = 4 AND e = 5 AND f = 6 AND g = 7 AND h = 8 AND i = 9 AND j = 10 AND k = 'Hello' AND l = 1`
	args := Arguments{
		ArgInt64(1, -2, 3, 4, 5, 6, 7, 8, 9, 10),
		ArgString("Hello"),
		ArgBool(true),
	}
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		var err error
		preprocessSink, err = Interpolate("SELECT * FROM x WHERE a = ? AND b = ? AND c = ? AND d = ? AND e = ? AND f = ? AND g = ? AND h = ? AND i = ? AND j = ? AND k = ? AND l = ?",
			args...,
		)
		if err != nil {
			b.Fatal(err)
		}
	}
	if preprocessSink != want {
		b.Fatalf("Have: %v Want: %v", preprocessSink, want)
	}
}

func TestInterpolateBools(t *testing.T) {
	t.Parallel()
	t.Run("single args", func(t *testing.T) {
		str, err := Interpolate("SELECT * FROM x WHERE a = ? AND b = ?", ArgBool(true, false))
		assert.NoError(t, err)
		assert.Equal(t, "SELECT * FROM x WHERE a = 1 AND b = 0", str)
	})
	t.Run("IN args", func(t *testing.T) {
		str, err := Interpolate("SELECT * FROM x WHERE a IN ? AND b = ? OR c = ?",
			ArgBool(true, false).Operator(In), ArgBool(true), ArgBool(false))
		assert.NoError(t, err)
		assert.Equal(t, "SELECT * FROM x WHERE a IN (1,0) AND b = 1 OR c = 0", str)
	})
	t.Run("empty arg", func(t *testing.T) {
		str, err := Interpolate("SELECT * FROM x WHERE a IN ? AND b = ? OR c = ?",
			ArgBool(true, false).Operator(In), ArgBool(true), ArgBool())
		assert.Empty(t, str)
		assert.True(t, errors.IsEmpty(err), "%+v", err)
	})
}

func TestInterpolateFloats(t *testing.T) {
	t.Parallel()
	t.Run("single args", func(t *testing.T) {
		str, err := Interpolate("SELECT * FROM x WHERE a = ? AND b = ?", ArgFloat64(3.14159, 2.7182818))
		assert.NoError(t, err)
		assert.Equal(t, "SELECT * FROM x WHERE a = 3.14159 AND b = 2.7182818", str)
	})
	t.Run("IN args", func(t *testing.T) {
		str, err := Interpolate("SELECT * FROM x WHERE a IN ? AND b = ?",
			ArgFloat64(3.14159, 2.7182818).Operator(In), ArgFloat64(0.815))
		assert.NoError(t, err)
		assert.Equal(t, "SELECT * FROM x WHERE a IN (3.14159,2.7182818) AND b = 0.815", str)
	})
	t.Run("empty args", func(t *testing.T) {
		var fl = make([]float64, 0, 2)
		str, err := Interpolate("SELECT * FROM x WHERE a IN ? AND b = ? OR c = ?",
			ArgFloat64(3.14159, 2.7182818).Operator(In), ArgFloat64(0.815), ArgFloat64(fl...))
		assert.True(t, errors.IsEmpty(err), "%+v", err)
		assert.Empty(t, str)
	})
}

func TestInterpolateStrings(t *testing.T) {
	t.Parallel()
	t.Run("single args", func(t *testing.T) {
		str, err := Interpolate("SELECT * FROM x WHERE a = ? AND b = ? AND c = ?", ArgString("a'b", "c`d"), ArgString("\"hello's \\ world\" \n\r\x00\x1a"))
		assert.NoError(t, err)
		assert.Equal(t, "SELECT * FROM x WHERE a = 'a\\'b' AND b = 'c`d' AND c = '\\\"hello\\'s \\\\ world\\\" \\n\\r\\x00\\x1a'", str)
	})
	t.Run("IN args", func(t *testing.T) {
		str, err := Interpolate("SELECT * FROM x WHERE a IN ? AND b = ?",
			ArgString("a'b", "c`d").Operator(In), ArgString("1' or '1' = '1'))/*"))
		assert.NoError(t, err)
		assert.Equal(t, "SELECT * FROM x WHERE a IN ('a\\'b','c`d') AND b = '1\\' or \\'1\\' = \\'1\\'))/*'", str)
	})
	t.Run("empty args", func(t *testing.T) {
		var fl = make([]string, 0, 2)
		str, err := Interpolate("SELECT * FROM x WHERE a IN ? AND b = ? OR c = ?",
			ArgString("a", "b").Operator(In), ArgString("c"), ArgString(fl...))
		assert.True(t, errors.IsEmpty(err), "%+v", err)
		assert.Empty(t, str)
	})

}

func TestInterpolateSlices(t *testing.T) {
	t.Parallel()
	str, err := Interpolate("SELECT * FROM x WHERE a = ? AND b = ? AND c = ? AND d = ? AND e = ?",
		argInt(1).Operator(In),
		ArgInt(1, 2, 3).Operator(In),
		ArgInt64(5, 6, 7).Operator(In),
		ArgString("wat", "ok").Operator(In),
		argInt(8),
	)
	assert.NoError(t, err)
	assert.Equal(t, "SELECT * FROM x WHERE a = (1) AND b = (1,2,3) AND c = (5,6,7) AND d = ('wat','ok') AND e = 8", str)
}

// TODO driver.Valuer
//type myString struct {
//	Present bool
//	Val     string
//}
//
//func (m myString) Value() (driver.Value, error) {
//	if m.Present {
//		return m.Val, nil
//	}
//	return nil, nil
//}
//
//func TestIntepolatingValuers(t *testing.T) {
//	args := []interface{}{myString{true, "wat"}, myString{false, "fry"}}
//
//	str, err := Interpolate("SELECT * FROM x WHERE a = ? AND b = ?", args)
//	assert.NoError(t, err)
//	assert.Equal(t, str, "SELECT * FROM x WHERE a = 'wat' AND b = NULL")
//}

func TestPreprocess(t *testing.T) {

	tests := []struct {
		sql string
		Arguments
		expSQL string
		errBhf errors.BehaviourFunc
	}{
		// NULL
		{"SELECT * FROM x WHERE a = ?", Arguments{ArgNull()},
			"SELECT * FROM x WHERE a = NULL", nil},

		// integers
		{
			`SELECT * FROM x WHERE a = ? AND b = ? AND c = ? AND d = ? AND e = ? AND f = ?
			AND g = ? AND h = ? AND i = ? AND j = ?`,
			Arguments{ArgInt(1, -2, 3, 4, 5, 6, 7, 8, 9, 10)},
			`SELECT * FROM x WHERE a = 1 AND b = -2 AND c = 3 AND d = 4 AND e = 5 AND f = 6
			AND g = 7 AND h = 8 AND i = 9 AND j = 10`, nil,
		},

		// boolean
		{"SELECT * FROM x WHERE a = ? AND b = ?", Arguments{ArgBool(true), ArgBool(false)},
			"SELECT * FROM x WHERE a = 1 AND b = 0", nil},

		// floats
		{"SELECT * FROM x WHERE a = ? AND b = ?", Arguments{ArgFloat64(0.15625), ArgFloat64(3.14159)},
			"SELECT * FROM x WHERE a = 0.15625 AND b = 3.14159", nil},

		// strings
		{
			`SELECT * FROM x WHERE a = ?
			AND b = ?`,
			Arguments{ArgString("hello", "\"hello's \\ world\" \n\r\x00\x1a")},
			`SELECT * FROM x WHERE a = 'hello'
			AND b = '\"hello\'s \\ world\" \n\r\x00\x1a'`, nil,
		},

		// slices
		{"SELECT * FROM x WHERE a = ? AND b = ? AND c = ? AND d = ?",
			Arguments{argInt(1), ArgInt(1, 2, 3).Operator(In), ArgInt(5, 6, 7).Operator(In), ArgString("wat", "ok").Operator(In)},
			"SELECT * FROM x WHERE a = 1 AND b = (1,2,3) AND c = (5,6,7) AND d = ('wat','ok')", nil},

		//// TODO valuers
		//{"SELECT * FROM x WHERE a = ? AND b = ?",
		//	Arguments{myString{true, "wat"}, myString{false, "fry"}},
		//	"SELECT * FROM x WHERE a = 'wat' AND b = NULL", nil},

		// errors
		{"SELECT * FROM x WHERE a = ? AND b = ?", Arguments{argInt64(1)},
			"", errors.IsNotValid},

		{"SELECT * FROM x WHERE", Arguments{argInt(1)},
			"", errors.IsNotValid},

		{"SELECT * FROM x WHERE a = ?", Arguments{ArgString(string([]byte{0x34, 0xFF, 0xFE}))},
			"", errors.IsNotValid},

		// ArgString() without arguments is equal to empty interface in the previous version.
		{"SELECT 'hello", Arguments{ArgString()}, "", errors.IsNotValid},
		{`SELECT "hello`, Arguments{ArgString()}, "", errors.IsNotValid},

		// preprocessing
		{"SELECT '?'", Arguments{ArgString()}, "SELECT '?'", nil},
		{"SELECT `?`", Arguments{ArgString()}, "SELECT `?`", nil},
		{"SELECT [?]", Arguments{ArgString()}, "SELECT `?`", nil},
		{"SELECT [name] FROM [user]", Arguments{ArgString()}, "SELECT `name` FROM `user`", nil},
		{"SELECT [u.name] FROM [user] [u]", Arguments{ArgString()}, "SELECT `u`.`name` FROM `user` `u`", nil},
		{"SELECT [u.na`me] FROM [user] [u]", Arguments{ArgString()}, "SELECT `u`.`na``me` FROM `user` `u`", nil},
		{"SELECT * FROM [user] WHERE [name] = '[nick]'", Arguments{ArgString()},
			"SELECT * FROM `user` WHERE `name` = '[nick]'", nil},
		{`SELECT * FROM [user] WHERE [name] = "nick[]"`, Arguments{ArgString()},
			"SELECT * FROM `user` WHERE `name` = 'nick[]'", nil},
	}

	for i, test := range tests {
		str, err := Interpolate(test.sql, test.Arguments...)
		if test.errBhf != nil {
			if !test.errBhf(err) {
				t.Errorf("IDX %d\ngot error: %v\nwant: %s", i, err, test.errBhf(err))
			}
		}
		if str != test.expSQL {
			t.Errorf("IDX %d\ngot: %v\nwant: %v", i, str, test.expSQL)
		}
	}
}
