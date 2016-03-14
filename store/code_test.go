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

package store_test

import (
	"testing"

	"github.com/corestoreio/csfw/store"
	"github.com/stretchr/testify/assert"
)

func TestValidateStoreCode(t *testing.T) {
	t.Parallel()
	tests := []struct {
		have    string
		wantErr error
	}{
		{"@de", store.ErrStoreCodeInvalid},
		{" de", store.ErrStoreCodeInvalid},
		{"de", nil},
		{"DE", nil},
		{"deCH09_", nil},
		{"_de", store.ErrStoreCodeInvalid},
		{"", store.ErrStoreCodeInvalid},
		{"\U0001f41c", store.ErrStoreCodeInvalid},
		{"au_en", nil},
		{"au-fr", store.ErrStoreCodeInvalid},
		{"Hello GoLang", store.ErrStoreCodeInvalid},
		{"Hello€GoLang", store.ErrStoreCodeInvalid},
		{"HelloGoLdhashdfkjahdjfhaskjdfhuiwehfiawehfuahweldsnjkasfkjkwejqwehqang", store.ErrStoreCodeInvalid},
	}
	for _, test := range tests {
		haveErr := store.CodeIsValid(test.have)
		if test.wantErr != nil {
			assert.EqualError(t, haveErr, test.wantErr.Error(), "err codes switched: %#v", test)
		} else {
			assert.NoError(t, haveErr, "%#v", test)
		}
	}
}
