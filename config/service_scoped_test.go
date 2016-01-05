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
	"github.com/corestoreio/csfw/config"
	"github.com/corestoreio/csfw/store/scope"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestScopedServiceString(t *testing.T) {
	tests := []struct {
		desc                        string
		fqpath                      string
		path                        []string
		websiteID, groupID, storeID int64
		err                         error
	}{
		{
			"Default ScopedGetter should return default scope",
			scope.StrDefault.FQPath("0", "a/b/c"), []string{"a/b/c"}, 0, 0, 0, nil,
		},
		{
			"Website ID 1 ScopedGetter should fall back to default scope",
			scope.StrDefault.FQPath("0", "a/b/c"), []string{"a/b/c"}, 1, 0, 0, nil,
		},
		{
			"Website ID 10 + Group ID 12 ScopedGetter should fall back to website 10 scope",
			scope.StrWebsites.FQPath("10", "a/b/c"), []string{"a/b/c"}, 10, 12, 0, nil,
		},
		{
			"Path consists of only two elements which is incorrect",
			scope.StrDefault.FQPath("0", "a/b/c"), []string{"a", "b"}, 0, 0, 0, config.ErrPathEmpty,
		},
	}
	for _, test := range tests {

		cg := config.NewMockGetter(config.WithMockValues(config.MockPV{
			test.fqpath: "Gopher",
		}))

		sg := cg.NewScoped(test.websiteID, test.groupID, test.storeID)
		s, err := sg.String(test.path...)

		if test.err != nil {
			assert.Empty(t, s, test.desc)
			assert.EqualError(t, err, test.err.Error(), test.desc)
			continue
		}
		assert.NoError(t, err, test.desc)
		assert.Exactly(t, "Gopher", s, test.desc)
	}
}

func TestScopedServiceBool(t *testing.T) {
	t.Log("TODO")
}

func TestScopedServiceFloat64(t *testing.T) {
	t.Log("TODO")
}
func TestScopedServiceInt(t *testing.T) {
	t.Log("TODO")
}
func TestScopedServiceDateTime(t *testing.T) {
	t.Log("TODO")
}
