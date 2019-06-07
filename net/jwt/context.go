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
	"context"

	"github.com/corestoreio/pkg/util/csjwt"
)

type keyCtxToken struct{}

type ctxTokenWrapper struct {
	t csjwt.Token
}

// withContext creates a new context with csjwt.Token attached.
func withContext(ctx context.Context, t csjwt.Token) context.Context {
	return context.WithValue(ctx, keyCtxToken{}, ctxTokenWrapper{t: t})
}

// FromContext returns the csjwt.Token in ctx if it exists or an error. If there
// is no token in the context then the error ErrContextJWTNotFound gets
// returned. Error behaviour: NotFound.
func FromContext(ctx context.Context) (csjwt.Token, bool) {
	wrp, ok := ctx.Value(keyCtxToken{}).(ctxTokenWrapper)
	return wrp.t, ok
}
