package ccache

import (
	"testing"

	"github.com/corestoreio/pkg/util/assert"
)

func TestConfiguration_BucketsPowerOf2(t *testing.T) {
	for i := uint32(0); i < 31; i++ {
		c := Configure().Buckets(i)
		if i == 1 || i == 2 || i == 4 || i == 8 || i == 16 {
			assert.Exactly(t, int(i), c.buckets)
		} else {
			assert.Exactly(t, 16, c.buckets)
		}
	}
}
