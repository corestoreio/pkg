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

package ctxjwt

import (
	"errors"

	"github.com/dgrijalva/jwt-go"
	"golang.org/x/net/context"
)

type keyJSONWebToken struct{}
type keyctxErr struct{}

// ErrContextJWTNotFound gets returned when the jwt cannot be found.
var ErrContextJWTNotFound = errors.New("Cannot extract ctxjwt nor an error from context")

// WithContext creates a new context with jwt.Token attached.
func WithContext(ctx context.Context, t *jwt.Token) context.Context {
	return context.WithValue(ctx, keyJSONWebToken{}, t)
}

// FromContext returns the jwt.Token in ctx if it exists or an error.
// Check the ok bool value if an error or jwt.Token is within the
// context.Context
func FromContext(ctx context.Context) (*jwt.Token, error) {
	err, ok := ctx.Value(keyctxErr{}).(error)
	if ok {
		return nil, err
	}
	t, ok := ctx.Value(keyJSONWebToken{}).(*jwt.Token)
	if !ok || t == nil {
		return nil, ErrContextJWTNotFound
	}
	return t, nil
}

// WithContextError creates a new context with an error attached.
func WithContextError(ctx context.Context, err error) context.Context {
	return context.WithValue(ctx, keyctxErr{}, err)
}
