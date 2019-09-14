// Copyright 2015-present, Cyrill @ Schumacher.fm and the CoreStore contributors
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

package config

import (
	"testing"

	"github.com/corestoreio/pkg/store/scope"
)

var _ getter = (*Service)(nil)
var _ getter = (*FakeService)(nil)

func TestScopedService_Parent(t *testing.T) {
	t.Parallel()
	tests := []struct {
		sg               Scoped
		wantCurrentScope scope.Type
		wantCurrentId    uint32
		wantParentScope  scope.Type
		wantParentID     uint32
	}{
		{makeScoped(nil, 33, 1), scope.Store, 1, scope.Website, 33},
		{makeScoped(nil, 3, 0), scope.Website, 3, scope.Default, 0},
		{makeScoped(nil, 0, 0), scope.Default, 0, scope.Default, 0},
	}
	for _, test := range tests {
		haveScp, haveID := test.sg.ParentID().Unpack()
		if have, want := haveScp, test.wantParentScope; have != want {
			t.Errorf("ParentScope: Have: %v Want: %v", have, want)
		}
		if have, want := haveID, test.wantParentID; have != want {
			t.Errorf("ParentScopeID: Have: %v Want: %v", have, want)
		}

		haveScp, haveID = test.sg.ScopeID().Unpack()
		if have, want := haveScp, test.wantCurrentScope; have != want {
			t.Errorf("Scope: Have: %v Want: %v", have, want)
		}
		if have, want := haveID, test.wantCurrentId; have != want {
			t.Errorf("ScopeID: Have: %v Want: %v", have, want)
		}

	}
}
