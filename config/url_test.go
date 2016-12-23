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

package config_test

import (
	"testing"

	"github.com/corestoreio/csfw/config"
	"github.com/corestoreio/errors"
	"github.com/stretchr/testify/assert"
)

func TestURLCache(t *testing.T) {

	tests := []struct {
		haveType  config.URLType
		url       string
		wantError errors.BehaviourFunc
	}{
		{config.URLTypeStatic, "", errors.IsEmpty},
		{config.URLTypeWeb, "http://corestore.io/", nil},
		{config.URLTypeStatic, "://corestore.io/", errors.IsNotValid},
		{config.URLType(254), "https://corestore.io/catalog", errors.IsNotFound},
	}
	for i, test := range tests {
		uc := config.NewURLCache()

		if test.wantError != nil {
			pu, err := uc.Set(test.haveType, test.url)
			assert.Nil(t, pu, "Index %d", i)
			assert.True(t, test.wantError(err), "Index %d => %s", i, err)
			assert.Nil(t, uc.Get(test.haveType))
			continue
		}

		pu, err := uc.Set(test.haveType, test.url) // pu = parsed URL
		assert.NoError(t, err, "Index %d", i)
		assert.Exactly(t, test.url, pu.String(), "Index %d", i)

		puCache := uc.Get(test.haveType)
		assert.Exactly(t, test.url, puCache.String(), "Index %d", i)

		assert.NoError(t, uc.Clear())
		assert.Nil(t, uc.Get(test.haveType), "Index %d", i)
	}
}
