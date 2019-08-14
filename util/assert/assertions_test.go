package assert_test

import (
	"testing"

	"github.com/corestoreio/errors"
	"github.com/corestoreio/pkg/util/assert"
)

func TestObjectsAreEqual(t *testing.T) {
	assert.True(t, assert.ObjectsAreEqual(1, 1))
	assert.False(t, assert.ObjectsAreEqual(0, false))
}

type NullString struct {
	String string
	Valid  bool
}

func TestExactlyLength(t *testing.T) {
	t.Run("strings", func(t *testing.T) {
		a, b := "abcdef", "abcdefgh"
		assert.ExactlyLength(t, 6, &a, &b, "Strings must match")
	})

	t.Run("structs NullString", func(t *testing.T) {
		expected := NullString{
			String: "abcdef",
		}
		b := NullString{
			String: "abcdefgh",
		}
		assert.ExactlyLength(t, 6, &expected, &b)
	})
}

func TestLenBetween(t *testing.T) {
	t.Run("slice", func(t *testing.T) {
		sl := [10]int{}
		assert.LenBetween(t, sl, 10, 10)
		assert.LenBetween(t, sl, 9, 11)
	})

	t.Run("struct with numeric fields", func(t *testing.T) {
		dm := &DataMaxLen{
			U8:   8,
			U16:  16,
			U32:  32,
			U64:  64,
			Uint: 128,

			I8:  8,
			I16: 16,
			I32: 32,
			I64: 64,
			Int: 128,

			Float32: 128,
			Float64: 128,
		}
		assert.LenBetween(t, dm.U8, 0, 8)
		assert.LenBetween(t, dm.U16, 0, 16)
		assert.LenBetween(t, dm.U32, 0, 32)
		assert.LenBetween(t, dm.U64, 0, 64)
		assert.LenBetween(t, dm.Uint, 0, 128)
		assert.LenBetween(t, dm.I8, 0, 8)
		assert.LenBetween(t, dm.I16, 0, 16)
		assert.LenBetween(t, dm.I32, 0, 32)
		assert.LenBetween(t, dm.I64, 0, 64)
		assert.LenBetween(t, dm.Int, 0, 128)
		assert.LenBetween(t, dm.Float32, 0, 128)
		assert.LenBetween(t, dm.Float64, 0, 128)
	})
}

type DataMaxLen struct {
	U8   uint8
	U16  uint16
	U32  uint32
	U64  uint64
	Uint uint

	I8  int8
	I16 int16
	I32 int32
	I64 int64
	Int int

	Float32 float32
	Float64 float64
}

func TestErrorIsKind(t *testing.T) {
	err := errors.AlreadyCaptured.Newf("the already captured err")
	assert.ErrorIsKind(t, errors.AlreadyCaptured, err)
}
