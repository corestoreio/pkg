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

package runmode

import (
	"context"
	"net/http"

	"github.com/corestoreio/pkg/store/scope"
)

// TODO important info on how to rework/refactor this package and maybe other.
//
// https://twitter.com/peterbourgon/status/752022730812317696
//
// https://gist.github.com/SchumacherFM/c62783b57621f791271eedf122de8de9

// Default defines the default run mode which is zero. It indicates the
// fall back to the default website and its default store.
const Default scope.TypeID = 0

// NOT SURE ABOUT THIS ONE
// CalculateRunMode transforms the Hash into a runMode. On an invalid Hash (the
// Type is < Website or Type > Store) it falls back to the default run mode,
// which is a zero Hash. Implements interface Calculater.
//func (t TypeID) CalculateRunMode(_ *http.Request) TypeID {
//	if s := t.Type(); s < Website || s > Store {
//		// fall back to default because only Website, Group and Store are allowed.
//		t = Default
//	}
//	return t
//}

// Calculater core type to initialize the run mode of the current
// request. Allows you to create a multi-site / multi-tenant setup. An
// implementation of this lives in net.runmode.WithRunMode() middleware.
//
// Your custom function allows to initialize the runMode based on parameters in
// the http.Request.
type Calculater interface {
	CalculateRunMode(*http.Request) scope.TypeID
}

// RunModeFunc type is an adapter to allow the use of ordinary functions as
// Calculater. If f is a function with the appropriate signature,
// RunModeFunc(f) is a Handler that calls f.
type RunModeFunc func(*http.Request) scope.TypeID

// CalculateRunMode calls f(r).
func (f RunModeFunc) CalculateRunMode(r *http.Request) scope.TypeID {
	return f(r)
}

// TODO(CYS) remove context function as the scope.TypeID does not belong into a context. See also Merovious

// WithContextRunMode sets the main run mode for the current request. It panics
// when called multiple times for the current context. This function is used in
// net/runmode together with function RunMode.CalculateMode(r, w).
// Use case for the runMode: Cache Keys and app initialization.
func WithContextRunMode(ctx context.Context, runMode scope.TypeID) context.Context {
	if _, ok := ctx.Value(ctxRunModeKey{}).(scope.TypeID); ok {
		panic("[scope] You are not allowed to set the runMode more than once for the current context.")
	}
	return context.WithValue(ctx, ctxRunModeKey{}, runMode)
}

// FromContextRunMode returns the run mode scope.TypeID from a context. If no
// entry can be found in the context the returned TypeID has a default value
// (0). This default value indicates the fall back to the default website and
// its default store. Use case for the runMode: Cache Keys and app
// initialization.
func FromContextRunMode(ctx context.Context) scope.TypeID {
	h, ok := ctx.Value(ctxRunModeKey{}).(scope.TypeID)
	if !ok {
		return Default // indicates a fall back to a default store of the default website
	}
	return h
}

type ctxRunModeKey struct{}
