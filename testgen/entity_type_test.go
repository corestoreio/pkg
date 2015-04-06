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

package testgen

import (
	"testing"

	"github.com/corestoreio/csfw/eav"
	"github.com/stretchr/testify/assert"
)

// TestGetCSEntityType uses the generated CSEntityTypeSlice from the test DB
func TestGetCSEntityType(t *testing.T) {

	tests := []struct {
		id   int64
		code string
		err  bool
	}{
		{1, "customer", false},
		{2, "customer_address", false},
		{3, "catalog_category", false},
		{4, "catalog_product", false},
		{40, "catalog_products", true},
	}

	for _, test := range tests {

		etc, err := eav.GetEntityTypeByID(test.id)
		if test.err {
			assert.Nil(t, etc)
			assert.Error(t, err)
		} else {
			assert.NotNil(t, etc)
			assert.NoError(t, err)
			assert.EqualValues(t, test.id, etc.EntityTypeID)
		}

		etc, err = eav.GetEntityTypeByCode(test.code)
		if test.err {
			assert.Nil(t, etc)
			assert.Error(t, err)
		} else {
			assert.NotNil(t, etc)
			assert.NoError(t, err)
			assert.EqualValues(t, test.id, etc.EntityTypeID)
			assert.EqualValues(t, test.code, etc.EntityTypeCode)
		}

	}

}
