// Copyright 2015, Cyrill @ Schumacher.fm and the CoreStore contributors
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

	"github.com/corestoreio/csfw/config/scope"
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
}
