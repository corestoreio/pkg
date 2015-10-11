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

package cshttp

import "golang.org/x/net/context"

// ctxKey type is unexported to prevent collisions with context keys defined in
// other packages.
type ctxKey uint

const (
	CtxAdminUserKey ctxKey = iota
	CtxStoreKey
	// todo more keys
)

// NewRootContext returns a context with the database set. This serves as the root
// context for all other contexts. @todo implementation
func NewRootContext(todo1, todo2, todo3, todo4 int) context.Context {
	return context.WithValue(context.Background(), CtxStoreKey, todo1)
}
