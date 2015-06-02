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

package config_test

import (
	"testing"

	"github.com/corestoreio/csfw/config"
	"github.com/stretchr/testify/assert"
)

func TestScopeBits(t *testing.T) {
	const (
		scope1 config.ScopeGroup = iota + 1
		scope2
		scope3
		scope4
		scope5
	)

	tests := []struct {
		have    []config.ScopeGroup
		want    config.ScopeGroup
		notWant config.ScopeGroup
		human   []string
	}{
		{[]config.ScopeGroup{scope1, scope2}, scope2, scope3, []string{"ScopeDefault", "ScopeWebsite"}},
		{[]config.ScopeGroup{scope3, scope4}, scope3, scope2, []string{"ScopeGroup", "ScopeStore"}},
		{[]config.ScopeGroup{scope4, scope5}, scope4, scope2, []string{"ScopeStore", "ScopeGroup(5)"}},
	}

	for _, test := range tests {
		var b config.ScopePerm
		b.Set(test.have...)
		if b.Has(test.want) == false {
			t.Errorf("%d should contain %d", b, test.want)
		}
		if b.Has(test.notWant) {
			t.Errorf("%d should not contain %d", b, test.notWant)
		}
		assert.EqualValues(t, test.human, b.Human())
	}
}
