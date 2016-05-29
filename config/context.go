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

package config

import (
	"github.com/corestoreio/csfw/config/internal/cfgctx"
	"golang.org/x/net/context"
)

// FromContextGetter returns a config.Getter from a context.
func FromContextGetter(ctx context.Context) (g Getter, ok bool) {
	g, ok = ctx.Value(cfgctx.KeyGetter{}).(Getter)
	return
}

// WithContextGetter adds a config.Getter to a context
func WithContextGetter(ctx context.Context, r Getter) context.Context {
	return context.WithValue(ctx, cfgctx.KeyGetter{}, r)
}

// FromContextScopedGetter returns a config.ScopedGetter from a context.
func FromContextScopedGetter(ctx context.Context) (r ScopedGetter, ok bool) {
	r, ok = ctx.Value(cfgctx.KeyScopedGetter{}).(ScopedGetter)
	return
}

// WithContextScopedGetter adds a config.ScopedGetter to a context
func WithContextScopedGetter(ctx context.Context, r ScopedGetter) context.Context {
	return context.WithValue(ctx, cfgctx.KeyScopedGetter{}, r)
}

// FromContextGetterPubSuber returns a config.GetterPubSuber from a context.
func FromContextGetterPubSuber(ctx context.Context) (r GetterPubSuber, ok bool) {
	r, ok = ctx.Value(cfgctx.KeyGetterPubSuber{}).(GetterPubSuber)
	return
}

// WithContextGetterPubSuber adds a GetterPubSuber to a context.
func WithContextGetterPubSuber(ctx context.Context, r GetterPubSuber) context.Context {
	return context.WithValue(ctx, cfgctx.KeyGetterPubSuber{}, r)
}
