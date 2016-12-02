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

package dbr

// These types are four callbacks to allow changes to the underlying SQL query
// builder by a 3rd party package.
type (
	SelectHook func(*Select)
	InsertHook func(*Insert)
	UpdateHook func(*Update)
	DeleteHook func(*Delete)

	SelectHooks []SelectHook
	InsertHooks []InsertHook
	UpdateHooks []UpdateHook
	DeleteHooks []DeleteHook
)

// Apply runs all SELECT hooks.
func (hs SelectHooks) Apply(b *Select) {
	for _, h := range hs {
		h(b)
	}
}

// Apply runs all INSERT hooks.
func (hs InsertHooks) Apply(b *Insert) {
	for _, h := range hs {
		h(b)
	}
}

// Apply runs all Update hooks.
func (hs UpdateHooks) Apply(b *Update) {
	for _, h := range hs {
		h(b)
	}
}

// Apply runs all DELETE hooks.
func (hs DeleteHooks) Apply(b *Delete) {
	for _, h := range hs {
		h(b)
	}
}

// Hook a type for embedding to define hooks for manipulating the SQL. DML
// stands for data manipulation language.
type Hook struct {
	SelectAfter SelectHooks
	InsertAfter InsertHooks
	UpdateAfter UpdateHooks
	DeleteAfter DeleteHooks
}

// NewHookDML creates a new set of hooks for data manipulation language
func NewHook() *Hook {
	return new(Hook)
}

// Merge merges one or more other hooks into the current hook.
func (h *Hook) Merge(hooks ...*Hook) *Hook {
	for _, hs := range hooks {
		h.AddSelectAfter(hs.SelectAfter...)
		h.AddInsertAfter(hs.InsertAfter...)
		h.AddUpdateAfter(hs.UpdateAfter...)
		h.AddDeleteAfter(hs.DeleteAfter...)
	}
	return h
}

func (h *Hook) AddSelectAfter(sh ...SelectHook) {
	h.SelectAfter = append(h.SelectAfter, sh...)
}

func (h *Hook) AddInsertAfter(sh ...InsertHook) {
	h.InsertAfter = append(h.InsertAfter, sh...)
}

func (h *Hook) AddUpdateAfter(sh ...UpdateHook) {
	h.UpdateAfter = append(h.UpdateAfter, sh...)
}

func (h *Hook) AddDeleteAfter(sh ...DeleteHook) {
	h.DeleteAfter = append(h.DeleteAfter, sh...)
}
