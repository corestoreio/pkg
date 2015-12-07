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

package ctxlog

import (
	"github.com/corestoreio/csfw/utils/log"
	"golang.org/x/net/context"
)

type keyLog struct{}

// WithContext creates a new context with jwt.Token attached.
func WithContext(ctx context.Context, l log.Logger) context.Context {
	return context.WithValue(ctx, keyLog{}, l)
}

// FromContext returns a log.Logger in ctx if it exists or an log.NullLogger.
func FromContext(ctx context.Context) log.Logger {
	if l, ok := ctx.Value(keyLog{}).(log.Logger); ok {
		return l
	}
	return log.BlackHole{}
}
