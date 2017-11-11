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

	"github.com/corestoreio/cspkg/store"
	"github.com/corestoreio/errors"
	"github.com/stretchr/testify/assert"
)

func TestValidateStoreCode(t *testing.T) {

	tests := []struct {
		have       string
		wantErrBhf errors.BehaviourFunc
	}{
		{"@de", errors.IsNotValid},
		{" de", errors.IsNotValid},
		{"de", nil},
		{"DE", nil},
		{"deCH09_", nil},
		{"_de", errors.IsNotValid},
		{"", errors.IsNotValid},
		{"\U0001f41c", errors.IsNotValid},
		{"au_en", nil},
		{"au-fr", errors.IsNotValid},
		{"Hello GoLang", errors.IsNotValid},
		{"Hello€GoLang", errors.IsNotValid},
		{"HelloGoLdhashdfkjahdjfhaskjdfhuiwehfiawehfuahweldsnjkasfkjkwejqwehqang", errors.IsNotValid},
	}
	for i, test := range tests {
		haveErr := store.CodeIsValid(test.have)
		if test.wantErrBhf != nil {
			assert.True(t, test.wantErrBhf(haveErr), "Index %d => %s", i, haveErr)
		} else {
			assert.NoError(t, haveErr, "Index %d", i)
		}
	}
}
