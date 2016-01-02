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

import "golang.org/x/net/context"

// ctx types are unexported to avoid collisions in context.Context with other packages
type (
	ctxKeyGetter         struct{}
	ctxKeyGetterPubSuber struct{}
	ctxKeyScopedGetter   struct{}
	ctxKeyWriter         struct{}
)

// FromContextGetter returns a config.Getter from a context. If not found returns
// the config.DefaultService
func FromContextGetter(ctx context.Context) Getter {
	if r, ok := ctx.Value(ctxKeyGetter{}).(Getter); ok {
		return r
	}
	return DefaultService
}

// WithContextGetter adds a config.Getter to a context
func WithContextGetter(ctx context.Context, r Getter) context.Context {
	return context.WithValue(ctx, ctxKeyGetter{}, r)
}

// FromContextScopedGetter returns a config.ScopedGetter from a context.
func FromContextScopedGetter(ctx context.Context) (r ScopedGetter, ok bool) {
	r, ok = ctx.Value(ctxKeyScopedGetter{}).(ScopedGetter)
	return
}

// WithContextScopedGetter adds a config.ScopedGetter to a context
func WithContextScopedGetter(ctx context.Context, r ScopedGetter) context.Context {
	return context.WithValue(ctx, ctxKeyScopedGetter{}, r)
}

// FromContextGetterPubSuber returns a config.GetterPubSuber from a context.
func FromContextGetterPubSuber(ctx context.Context) (r GetterPubSuber, ok bool) {
	r, ok = ctx.Value(ctxKeyGetterPubSuber{}).(GetterPubSuber)
	return
}

// WithContextGetterPubSuber adds a GetterPubSuber to a context.
func WithContextGetterPubSuber(ctx context.Context, r GetterPubSuber) context.Context {
	return context.WithValue(ctx, ctxKeyGetterPubSuber{}, r)
}

// FromContextWriter returns a config.Writer from a context.
func FromContextWriter(ctx context.Context) (w Writer, ok bool) {
	w, ok = ctx.Value(ctxKeyWriter{}).(Writer)
	return
}

// WithContextWriter adds a writer to a context
func WithContextWriter(ctx context.Context, w Writer) context.Context {
	return context.WithValue(ctx, ctxKeyWriter{}, w)
}
