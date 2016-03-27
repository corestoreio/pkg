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

package freecache_test

import (
	"github.com/corestoreio/csfw/config/cfgpath"
	"github.com/corestoreio/csfw/config/storage"
	"github.com/corestoreio/csfw/config/storage/freecache"
	"github.com/stretchr/testify/assert"
	"testing"
)

var _ storage.Storager = (*freecache.Storage)(nil)

func TestCacheGet(t *testing.T) {

	sc := freecache.New(0)

	tests := []struct {
		key        cfgpath.Path
		val        interface{}
		wantSetErr error
		wantGetErr error
	}{
		{cfgpath.MustNewByParts("aa/bb/cc"), 1, nil, nil},
	}
	for idx, test := range tests {

		haveSetErr := sc.Set(test.key, test.val)
		if test.wantSetErr != nil {
			assert.EqualError(t, haveSetErr, test.wantSetErr.Error(), "Index %d", idx)
		} else {
			assert.NoError(t, haveSetErr, "Index %d", idx)
		}

		haveVal, haveGetErr := sc.Get(test.key)
		if test.wantGetErr != nil {
			assert.EqualError(t, haveGetErr, test.wantGetErr.Error(), "Index %d", idx)
		} else {
			assert.NoError(t, haveGetErr, "Index %d", idx)
		}

		assert.Exactly(t, test.val, haveVal, "Index %d", idx)
	}
}
