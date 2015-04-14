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

	"github.com/stretchr/testify/assert"
)

// @todo add group and websites

func TestStore(t *testing.T) {
	s1 := storeManager.Store().Collection()
	assert.True(t, len(s1) > 1, "There should be at least two stores in the slice")
	//assert.Equal(t, storeCollection, s1)

	s2 := storeManager.Store().Collection()

	for _, store := range s2 {
		assert.NotNil(t, store, "Expecting first index to be nil")
		assert.True(t, len(store.Code.String) > 1, "store.Code.String should be longer than 1 char: %#v", store)
		//t.Logf("\n%d : %#v\n", i, store)
	}

}

func TestGetStoreByCode(t *testing.T) {
	s, err := storeManager.Store().ByCode("german")
	if err != nil {
		t.Error(err)
		assert.Nil(t, s)
	} else {
		assert.NoError(t, err)
		assert.Equal(t, "german", s.Code.String)
	}
}

func TestGetStoreByID(t *testing.T) {
	s, err := storeManager.Store().ByID(2)
	if err != nil {
		t.Error(err)
		assert.Nil(t, s)
	} else {
		assert.NoError(t, err)
		assert.Equal(t, "french", s.Code.String)
	}
}
