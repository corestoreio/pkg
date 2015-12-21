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

package model

import (
	"strconv"

	"github.com/corestoreio/csfw/config"
	"github.com/corestoreio/csfw/config/scope"
	"github.com/corestoreio/csfw/config/valuelabel"
)

// Bool represents a path in config.Getter which handles bool values.
type Bool struct {
	basePath
}

// NewBool creates a new Bool model with a given path.
func NewBool(path string, vlPairs ...valuelabel.Pair) Bool {
	return Bool{
		basePath: NewPath(path, vlPairs...),
	}
}

// Get returns a bool value.
func (p Bool) Get(pkgCfg config.SectionSlice, sg config.ScopedGetter) (v bool) {
	var err error
	if v, err = p.lookupBool(pkgCfg, sg); err != nil && PkgLog.IsDebug() {
		PkgLog.Debug("model.Bool.Get.lookupBool", "err", err, "path", p.string)
	}
	return
}

// Write writes a bool. Bool gets internally converted to type string.
func (p Bool) Write(w config.Writer, v bool, s scope.Scope, id int64) error {
	return p.basePath.Write(w, strconv.FormatBool(v), s, id)
}

// Str represents a path in config.Getter which handles string values.
// The name Str has been chosen to avoid conflict with the String() function
// in the Stringer interface.
type Str struct{ basePath }

// NewStr creates a new Str model with a given path.
func NewStr(path string, vlPairs ...valuelabel.Pair) Str {
	return Str{basePath: NewPath(path, vlPairs...)}
}

// Get returns a string value
func (p Str) Get(pkgCfg config.SectionSlice, sg config.ScopedGetter) (v string) {
	var err error
	if v, err = p.lookupString(pkgCfg, sg); err != nil && PkgLog.IsDebug() {
		PkgLog.Debug("model.Str.Get.lookupString", "err", err, "path", p.string)
	}
	return
}

// Write writes a string value
func (p Str) Write(w config.Writer, v string, s scope.Scope, id int64) error {
	return p.basePath.Write(w, v, s, id)
}

// Int represents a path in config.Getter which handles int values.
type Int struct{ basePath }

// NewInt creates a new Int model with a given path.
func NewInt(path string, vlPairs ...valuelabel.Pair) Int {
	return Int{basePath: NewPath(path, vlPairs...)}
}

// Get returns an int value.
func (p Int) Get(pkgCfg config.SectionSlice, sg config.ScopedGetter) (v int) {
	var err error
	if v, err = p.lookupInt(pkgCfg, sg); err != nil && PkgLog.IsDebug() {
		PkgLog.Debug("model.Int.Get.lookupInt", "err", err, "path", p.string)
	}
	return
}

// Write writes an int value as a string.
func (p Int) Write(w config.Writer, v int, s scope.Scope, id int64) error {
	return p.basePath.Write(w, strconv.Itoa(v), s, id)
}

// Float64 represents a path in config.Getter which handles int values.
type Float64 struct{ basePath }

// NewFloat64 creates a new Float64 model with a given path.
func NewFloat64(path string, vlPairs ...valuelabel.Pair) Float64 {
	return Float64{basePath: NewPath(path, vlPairs...)}
}

// Get returns a float64 value.
func (p Float64) Get(pkgCfg config.SectionSlice, sg config.ScopedGetter) (v float64) {
	var err error
	if v, err = p.lookupFloat64(pkgCfg, sg); err != nil && PkgLog.IsDebug() {
		PkgLog.Debug("model.Float64.Get.lookupFloat64", "err", err, "path", p.string)
	}
	return
}

// Write writes a float64 value as a string.
func (p Float64) Write(w config.Writer, v float64, s scope.Scope, id int64) error {
	return p.basePath.Write(w, strconv.FormatFloat(v, 'f', 14, 64), s, id)
}
