package bgwork

import (
	"errors"
	"testing"

	"github.com/corestoreio/pkg/util/assert"
)

func TestAsync(t *testing.T) {
	t.Run("no error", func(t *testing.T) {
		fn1 := func() error {
			return nil
		}
		fn2 := func() error {
			return nil
		}
		fn3 := func() error {
			return nil
		}
		count := 0
		for err := range Async(fn1, fn2, fn3) {
			assert.NoError(t, err)
			count++
		}
		assert.Exactly(t, 3, count)
	})
	t.Run("all error", func(t *testing.T) {
		fn1 := func() error {
			return errors.New("fn1")
		}
		fn2 := func() error {
			return errors.New("fn2")
		}
		fn3 := func() error {
			return errors.New("fn3")
		}
		count := 0
		for err := range Async(fn1, fn2, fn3) {
			assert.Error(t, err)
			count++
		}
		assert.Exactly(t, 3, count)
	})
}
