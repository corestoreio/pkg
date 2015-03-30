// Copyright 2015 CoreStore Authors
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

package store_test

import (
	"testing"

	"github.com/corestoreio/csfw/store"
	"github.com/stretchr/testify/assert"
)

func TestGetStores(t *testing.T) {
	s1 := store.GetStores()
	assert.True(t, len(s1) > 1, "There should be at least two stores in the slice")
	//assert.Equal(t, storeCollection, s1)

	s2 := store.GetStores()

	for i, store := range s2 {
		t.Logf("\n%d : %#v\n", i, store)
	}

	//assert.Len(t, s2, len(storeCollection)-1)
}

func TestGetStoreByID(t *testing.T) {
	s, err := store.GetStoreByCode("german")
	assert.NoError(t, err)
	assert.Equal(t, "german", s.Code.String)
}
