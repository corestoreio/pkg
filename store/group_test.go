// Copyright 2015 CoreGroup Authors
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

func TestGroup(t *testing.T) {
	s1 := storeManager.Group().Collection()
	assert.True(t, len(s1) > 2, "There should be at least three groups in the slice")

	for i, group := range storeManager.Group().Collection() {
		if i == 0 {
			assert.Nil(t, group, "Expecting first index to be nil")
			continue
		}
		assert.True(t, len(group.Name) > 1, "group.Name should be longer than 1 char: %#v", group)
	}

}

func TestGetGroupByID(t *testing.T) {
	g, err := storeManager.Group().ByID(1)
	if err != nil {
		t.Error(err)
		assert.Nil(t, g)
	} else {
		assert.NoError(t, err)
		assert.Equal(t, "Madison Island", g.Name)
	}
}

func TestGroupGetStores(t *testing.T) {
	gInvalid, err := storeManager.Group().ByID(10000)
	assert.EqualError(t, err, store.ErrGroupNotFound.Error())
	assert.Nil(t, gInvalid)

	stores, err := storeManager.Group().Stores(1)
	assert.NoError(t, err)
	t.Logf("\n%#v\n", stores)

}
func TestGroupGetWebsite(t *testing.T) {}
