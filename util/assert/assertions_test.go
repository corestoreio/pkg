package assert_test

import (
	"testing"

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

func TestMaxLen(t *testing.T) {
	sl := [10]int{}
	assert.LenBetween(t, sl, 10, 10)
	assert.LenBetween(t, sl, 9, 11)
}
