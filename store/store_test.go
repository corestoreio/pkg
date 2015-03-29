package store

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestGetStores(t *testing.T) {
	s1 := GetStores()
	assert.True(t, len(s1) > 1, "There should be at least two stores in the slice")
	assert.Equal(t, storeCollection, s1)

	s2 := GetStores(false)

	for i, store := range s2 {
		t.Logf("\n%d : %#v\n", i, store)
	}

	//assert.Len(t, s2, len(storeCollection)-1)
}
