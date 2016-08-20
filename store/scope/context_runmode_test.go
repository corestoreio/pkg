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
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/corestoreio/csfw/store/scope"
	"github.com/stretchr/testify/assert"
)

var _ scope.RunModeCalculater = (*scope.RunModeFunc)(nil)

func TestRunModeFunc_CalculateRunMode(t *testing.T) {
	h := scope.NewHash(scope.Website, 33)
	rmf := scope.RunModeFunc(func(_ *http.Request) scope.Hash {
		return h
	})
	assert.Exactly(t, h, rmf.CalculateRunMode(nil))
}

func TestRunMode(t *testing.T) {

	tests := []struct {
		mode scope.RunModeCalculater
		want scope.Hash
	}{
		{scope.NewHash(scope.Website, 2), scope.NewHash(scope.Website, 2)},
		{scope.NewHash(scope.Store, 3), scope.NewHash(scope.Store, 3)},
		{scope.NewHash(scope.Group, 4), scope.NewHash(scope.Group, 4)},
		{scope.NewHash(scope.Store, 0), scope.NewHash(scope.Store, 0)},
		{scope.NewHash(scope.Default, 0), 0},
	}
	for i, test := range tests {
		req := httptest.NewRequest("GET", "http://corestore.io", nil)

		haveMode := test.mode.CalculateRunMode(nil)

		ctx := scope.WithContextRunMode(req.Context(), haveMode)

		assert.Exactly(t, test.want, haveMode, "Index %d", i)
		assert.Exactly(t, test.want, scope.FromContextRunMode(ctx), "Index %d", i)
	}
	assert.Exactly(t, scope.Hash(0), scope.FromContextRunMode(context.Background()))
}

func TestWithContextRunMode(t *testing.T) {
	defer func() {
		if r := recover(); r != nil {
			assert.NotEmpty(t, r.(string))
		} else {
			t.Fatal("Expecting a panic")
		}
	}()

	ctx := context.Background()
	ctx = scope.WithContextRunMode(ctx, scope.DefaultHash)
	_ = scope.WithContextRunMode(ctx, scope.DefaultHash)
}
