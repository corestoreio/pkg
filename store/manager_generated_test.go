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

var storeManager *store.Manager

func TestGeneratedNewManager(t *testing.T) {
	if storeManager == nil {
		t.Skip("storeManager variable is nil. Integration test skipped")
	}

	t.Log("@todo")

	tests := []struct {
		haveID   store.IDRetriever
		haveCode store.CodeRetriever
		wantErr  error
		wantCode string
	}{
		{nil, store.Code("de"), nil, "de"},
		{nil, store.Code("cz"), store.ErrStoreNotFound, ""},
		{nil, store.Code("de"), nil, "de"},
	}

	for _, test := range tests {
		s, err := storeManager.Store(test.haveID, test.haveCode)
		if test.wantErr == nil {
			assert.NoError(t, err)
			assert.NotNil(t, s)
			assert.EqualValues(t, test.wantCode, s.Data().Code.String)
		} else {
			assert.Error(t, err)
			assert.EqualError(t, test.wantErr, err.Error())
			assert.Nil(t, s)
		}
	}

}
