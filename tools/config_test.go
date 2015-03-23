package tools

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEntityTypeMapKeys(t *testing.T) {
	assert.Len(t, ConfigEntityType.Keys(), len(ConfigEntityType))
}

var benchmarkEntityTypeMapKeys []string

// BenchmarkEntityTypeMapKeys	 5000000	       328 ns/op	      64 B/op	       1 allocs/op
func BenchmarkEntityTypeMapKeys(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		benchmarkEntityTypeMapKeys = ConfigEntityType.Keys()
	}
}
