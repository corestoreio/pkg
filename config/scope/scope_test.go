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

func TestScopeBits(t *testing.T) {
	const (
		scope1 scope.Group = iota + 1
		scope2
		scope3
		scope4
		scope5
	)

	tests := []struct {
		have    []scope.Group
		want    scope.Group
		notWant scope.Group
		human   []string
	}{
		{[]scope.Group{scope1, scope2}, scope2, scope3, []string{"ScopeDefault", "ScopeWebsite"}},
		{[]scope.Group{scope3, scope4}, scope3, scope2, []string{"ScopeGroup", "ScopeStore"}},
		{[]scope.Group{scope4, scope5}, scope4, scope2, []string{"ScopeStore", "ScopeGroup(5)"}},
	}

	for _, test := range tests {
		var b scope.Perm
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
