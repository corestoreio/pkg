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

package scope_test

import (
	"context"
	"testing"

	"github.com/corestoreio/csfw/store/scope"
	"github.com/stretchr/testify/assert"
)

func TestFromContext(t *testing.T) {
	tests := []struct {
		c    scope.Hash
		p    scope.Hash
		want bool
	}{
		{scope.DefaultHash, scope.DefaultHash, true},
		{scope.NewHash(scope.Website, 1), scope.DefaultHash, true},
		{scope.NewHash(scope.Website, 0), scope.DefaultHash, true},
		{scope.NewHash(scope.Store, 1), scope.NewHash(scope.Website, 1), true},
		{scope.NewHash(scope.Store, -1), scope.NewHash(scope.Website, 1), false},
		{scope.NewHash(scope.Store, 1), scope.NewHash(scope.Website, -1), false},
		{scope.NewHash(scope.Store, 0), scope.NewHash(scope.Website, 0), true},
		{scope.DefaultHash, scope.NewHash(scope.Website, 1), false},
		{0, 0, false},
		{0, scope.DefaultHash, false},
		{scope.DefaultHash, 0, false},
	}
	for i, test := range tests {
		ctx := scope.WithContext(context.TODO(), test.c, test.p)
		haveC, haveP, haveOK := scope.FromContext(ctx)
		if have, want := haveOK, test.want; have != want {
			t.Errorf("(%d) Have: %v Want: %v", i, have, want)
		}
		if have, want := haveC, test.c; have != want {
			t.Errorf("Current Have: %v Want: %v", have, want)
		}
		if have, want := haveP, test.p; have != want {
			t.Errorf("Parent Have: %v Want: %v", have, want)
		}
	}
	c, h, ok := scope.FromContext(context.Background())
	assert.Exactly(t, scope.Hash(0), c)
	assert.Exactly(t, scope.Hash(0), h)
	assert.False(t, ok)
}
