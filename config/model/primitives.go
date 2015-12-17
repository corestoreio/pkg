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
	"github.com/corestoreio/csfw/util/cast"
)

// Bool represents a path in config.Getter which handles bool values.
type Bool struct {
	Path
}

// NewBool creates a new Bool model with a given path.
func NewBool(path string, vlPairs ...valuelabel.Pair) Bool {
	return Bool{
		Path: NewPath(path, vlPairs...),
	}
}

// Get returns a bool value.
func (p Bool) Get(pkgCfg config.SectionSlice, sg config.ScopedGetter) (v bool) {

	if fields, err := pkgCfg.FindFieldByPath(p.string); err == nil {
		v, _ = cast.ToBoolE(fields.Default)
	} else {
		if PkgLog.IsDebug() {
			PkgLog.Debug("model.Bool.SectionSlice.FindFieldByPath", "err", err, "path", p.string)
		}
	}

	if val, err := sg.Bool(p.string); err == nil {
		v = val
	}
	return v
}

// Set writes a bool. Bool gets internally converted to type string.
func (p Bool) Set(w config.Writer, v bool, s scope.Scope, id int64) error {
	return p.Path.Set(w, strconv.FormatBool(v), s, id)
}

// String represents a path in config.Getter which handles string values.
type String struct{ Path }

// NewString creates a new String model with a given path.
func NewString(path string, vlPairs ...valuelabel.Pair) String {
	return String{Path: NewPath(path, vlPairs...)}
}

// Get returns a string value
func (p String) Get(pkgCfg config.SectionSlice, sg config.ScopedGetter) (v string) {
	return p.LookupString(pkgCfg, sg)
}

// Set writes a string value
func (p String) Set(w config.Writer, v string, s scope.Scope, id int64) error {
	return p.Path.Set(w, v, s, id)
}

// Int represents a path in config.Getter which handles int values.
type Int struct{ Path }

// NewInt creates a new Int model with a given path.
func NewInt(path string, vlPairs ...valuelabel.Pair) Int { return Int{Path: NewPath(path, vlPairs...)} }

// Get returns an int value.
func (p Int) Get(pkgCfg config.SectionSlice, sg config.ScopedGetter) (v int) {
	if fields, err := pkgCfg.FindFieldByPath(p.string); err == nil {
		v, _ = cast.ToIntE(fields.Default)
	} else {
		if PkgLog.IsDebug() {
			PkgLog.Debug("model.Int.SectionSlice.FindFieldByPath", "err", err, "path", p.string)
		}
	}

	if val, err := sg.Int(p.string); err == nil {
		v = val
	}
	return v
}

// Set writes an int value as a string.
func (p Int) Set(w config.Writer, v int, s scope.Scope, id int64) error {
	return p.Path.Set(w, strconv.Itoa(v), s, id)
}

// Float64 represents a path in config.Getter which handles int values.
type Float64 struct{ Path }

// NewFloat64 creates a new Float64 model with a given path.
func NewFloat64(path string, vlPairs ...valuelabel.Pair) Float64 {
	return Float64{Path: NewPath(path, vlPairs...)}
}

// Get returns a float64 value.
func (p Float64) Get(pkgCfg config.SectionSlice, sg config.ScopedGetter) (v float64) {
	if fields, err := pkgCfg.FindFieldByPath(p.string); err == nil {
		v, _ = cast.ToFloat64E(fields.Default)
	} else {
		if PkgLog.IsDebug() {
			PkgLog.Debug("model.Float64.SectionSlice.FindFieldByPath", "err", err, "path", p.string)
		}
	}

	if val, err := sg.Float64(p.string); err == nil {
		v = val
	}
	return v
}

// Set writes a float64 value as a string.
func (p Float64) Set(w config.Writer, v float64, s scope.Scope, id int64) error {
	return p.Path.Set(w, strconv.FormatFloat(v, 'f', 14, 64), s, id)
}
