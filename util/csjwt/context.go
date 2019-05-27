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

package csjwt

import (
	"context"
)

// keyCtxToken type is unexported to prevent collisions with context keys
// defined in other packages.
type keyCtxToken struct{}

// ctxTokenWrapper to prevent too much calls to runtime.convT2*
type ctxTokenWrapper struct {
	*Token
}

// WithContextToken creates a new context with a token attached.
func WithContextToken(ctx context.Context, t *Token) context.Context {
	return context.WithValue(ctx, keyCtxToken{}, ctxTokenWrapper{Token: t})
}

// FromContextToken returns the token in ctx if it exists or an error.
func FromContextToken(ctx context.Context) (*Token, bool) {
	wrp, ok := ctx.Value(keyCtxToken{}).(ctxTokenWrapper)
	if !ok {
		return nil, ok
	}
	return wrp.Token, ok
}
