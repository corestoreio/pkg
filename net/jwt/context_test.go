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

package jwt

import (
	"testing"

	"context"

	"github.com/corestoreio/pkg/util/csjwt"
	"github.com/corestoreio/pkg/util/assert"
)

func TestFromContext_Token(t *testing.T) {

	ctx := withContext(context.Background(), csjwt.Token{})
	assert.NotNil(t, ctx)

	haveToken, ok := FromContext(ctx)
	assert.NotNil(t, haveToken)
	assert.False(t, haveToken.Valid)
	assert.True(t, ok)
}

func TestFromContext_NoToken(t *testing.T) {
	haveToken, ok := FromContext(context.Background())
	assert.NotNil(t, haveToken)
	assert.False(t, haveToken.Valid)
	assert.False(t, ok)
}
