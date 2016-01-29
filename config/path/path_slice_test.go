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

package path_test

import (
	"github.com/corestoreio/csfw/config/path"
	"testing"

	"github.com/corestoreio/csfw/store/scope"
	"github.com/stretchr/testify/assert"
)

func TestPathSliceContains(t *testing.T) {
	tests := []struct {
		paths  path.PathSlice
		search path.Path
		want   bool
	}{
		{
			path.PathSlice{
				0: path.MustNewByParts("aa/bb/cc").Bind(scope.WebsiteID, 3),
				1: path.MustNewByParts("aa/bb/cc").Bind(scope.WebsiteID, 2),
			},
			path.MustNewByParts("aa/bb/cc").Bind(scope.WebsiteID, 2),
			true,
		},
		{
			path.PathSlice{
				0: path.MustNewByParts("aa/bb/cc").Bind(scope.WebsiteID, 3),
				1: path.MustNewByParts("aa/bb/cc").Bind(scope.WebsiteID, 2),
			},
			path.MustNewByParts("aa/bb/cc").Bind(scope.StoreID, 2),
			false,
		},
	}
	for i, test := range tests {
		assert.Exactly(t, test.want, test.paths.Contains(test.search), "Index %d", i)
	}
}
