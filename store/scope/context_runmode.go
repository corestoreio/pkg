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

// defaultRunMode defines the default run mode if the programmer hasn't applied
// the function RunMode.WithContext() to specify a specific run mode.
// It indicates the fall back to the default website and its default store.
const defaultRunMode Hash = 0

// RunMode core type to initialize the run mode of the current request. Allows
// you to create a multi-site / multi-tenant setup.
type RunMode struct {
	Mode Hash
	// ModeFunc if not nil you can create your own function to set a run mode.
	ModeFunc func(http.ResponseWriter, *http.Request) Hash
}

// WithContext sets the main run mode for a request. Precedence for applying the
// mode are: First field Mode and if not nil field ModeFunc. Returns a shallow
// copy of the http.Request and the applied run mode Hash.
func (rm RunMode) WithContext(w http.ResponseWriter, r *http.Request) (*http.Request, Hash) {
	ctx := r.Context()

	h := rm.Mode
	if rm.ModeFunc != nil {
		h = rm.ModeFunc(w, r)
	}
	if s := h.Scope(); s < Website || s > Store {
		// fall back to default because only Website, Group and Store are allowed.
		h = defaultRunMode
	}
	ctx = context.WithValue(ctx, ctxRunModeKey{}, h)
	return r.WithContext(ctx), h
}

// FromContextRunMode returns the run mode Hash from a context. If no entry can
// be found in the context the returned Hash has a default value. This default
// value indicates the fall back to the default website and its default store.
func FromContextRunMode(ctx context.Context) Hash {
	h, ok := ctx.Value(ctxRunModeKey{}).(Hash)
	if !ok {
		return defaultRunMode // indicates a fall back to a default store of the default website
	}
	return h
}
