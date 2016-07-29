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

package scope

import (
	"context"
	"net/http"
)

type ctxRunModeKey struct{}

// DefaultRunMode defines the default run mode if the programmer hasn't applied
// the field Mode or the function RunMode.WithContext() to specify a specific
// run mode. It indicates the fall back to the default website and its default
// store.
const DefaultRunMode Hash = 0

// RunMode core type to initialize the run mode of the current request. Allows
// you to create a multi-site / multi-tenant setup. An implementation of this
// lives in storenet.AppRunMode.WithRunMode() middleware.
type RunMode struct {
	Mode Hash
	// ModeFunc if not nil you can create your own function to set a run mode.
	ModeFunc func(http.ResponseWriter, *http.Request) Hash
}

// CalculateMode calls the user defined Mode field or ModeFunction. On an
// invalid mode it falls back to the default run mode, which is a zero Hash.
func (rm RunMode) CalculateMode(w http.ResponseWriter, r *http.Request) Hash {
	h := rm.Mode
	if rm.ModeFunc != nil {
		h = rm.ModeFunc(w, r)
	}
	if s := h.Scope(); s < Website || s > Store {
		// fall back to default because only Website, Group and Store are allowed.
		h = DefaultRunMode
	}
	return h
}

// WithContextRunMode sets the main run mode for the current request. Use the
// Hash value returned from the function RunMode.CalculateMode(r, w).
func WithContextRunMode(ctx context.Context, runMode Hash) context.Context {
	return context.WithValue(ctx, ctxRunModeKey{}, runMode)
}

// FromContextRunMode returns the run mode Hash from a context. If no entry can
// be found in the context the returned Hash has a default value. This default
// value indicates the fall back to the default website and its default store.
func FromContextRunMode(ctx context.Context) Hash {
	h, ok := ctx.Value(ctxRunModeKey{}).(Hash)
	if !ok {
		return DefaultRunMode // indicates a fall back to a default store of the default website
	}
	return h
}
