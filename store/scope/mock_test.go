// Copyright 2015-2016, Cyrill @ Schumacher.fm and the CoreStore contributors
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

package scope_test

import (
	"testing"

	"github.com/corestoreio/csfw/store/scope"
	"github.com/stretchr/testify/assert"
)

func TestMockID(t *testing.T) {

	var e int64 = 29
	assert.Equal(t, e, scope.MockID(29).StoreID())
	assert.Equal(t, e, scope.MockID(29).GroupID())
	assert.Equal(t, e, scope.MockID(29).WebsiteID())
}

func TestMockCode(t *testing.T) {

	assert.Equal(t, "Waverly", scope.MockCode("Waverly").StoreCode())
	assert.Equal(t, "Waverly", scope.MockCode("Waverly").WebsiteCode())
	var i int64 = -1
	assert.Equal(t, i, scope.MockCode("Waverly").WebsiteID())
	assert.Equal(t, i, scope.MockCode("Waverly").GroupID())
	assert.Equal(t, i, scope.MockCode("Waverly").StoreID())
}

func TestMock(t *testing.T) {

	tests := []struct {
		s  scope.Scope
		id int64
	}{
		{scope.Default, 0},
		{scope.Website, 1},
		{scope.Group, 20},
		{scope.Store, 30},
	}
	for _, test := range tests {
		m := scope.Mock{
			Scp: test.s,
			ID:  test.id,
		}
		haveS, haveID := m.Scope()
		assert.Exactly(t, test.s, haveS)
		assert.Exactly(t, test.id, haveID)
	}
}
