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

package ctxjwt

import (
	"errors"
	"testing"

	"github.com/corestoreio/csfw/util/csjwt"
	"github.com/stretchr/testify/assert"
	"golang.org/x/net/context"
)

func TestContextWithError(t *testing.T) {
	t.Parallel()

	var wantErr = errors.New("Contiki Context")
	ctx := withContextError(context.Background(), wantErr)
	assert.NotNil(t, ctx)

	haveToken, haveErr := FromContext(ctx)
	assert.NotNil(t, haveToken)
	assert.False(t, haveToken.Valid)
	assert.EqualError(t, haveErr, wantErr.Error())
}

func TestFromContext(t *testing.T) {
	t.Parallel()

	ctx := withContext(context.Background(), csjwt.Token{})
	assert.NotNil(t, ctx)

	haveToken, haveErr := FromContext(ctx)
	assert.NotNil(t, haveToken)
	assert.False(t, haveToken.Valid)
	assert.EqualError(t, haveErr, errContextJWTNotFound.Error())
}
