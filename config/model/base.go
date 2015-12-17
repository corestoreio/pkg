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
	"github.com/corestoreio/csfw/config"
	"github.com/corestoreio/csfw/config/scope"
	"github.com/corestoreio/csfw/config/valuelabel"
	"github.com/corestoreio/csfw/util/cast"
)

// Path defines the path in the "core_config_data" table like a/b/c. All other
// types in this package inherits from this path type.
type Path struct {
	string  // contains the path
	Options valuelabel.Slice
}

// NewPath creates a new path type optional value label Pair
func NewPath(path string, vlPairs ...valuelabel.Pair) Path {
	return Path{
		string:  path,
		Options: valuelabel.Slice(vlPairs),
	}
}

// Write writes a value v to the config.Writer without checking if the value
// has changed.
func (p Path) Write(w config.Writer, v interface{}, s scope.Scope, id int64) error {
	return w.Write(config.Path(p.string), config.Value(v), config.Scope(s, id))
}

// LookupString searches in default value in config.SectionSlice and overrides
// it with a value from ScopedGetter if ScopedGetter is not empty.
func (p Path) LookupString(pkgCfg config.SectionSlice, sg config.ScopedGetter) (v string) {

	if fields, err := pkgCfg.FindFieldByPath(p.string); err == nil {
		v, _ = cast.ToStringE(fields.Default)
	} else {
		if PkgLog.IsDebug() {
			PkgLog.Debug("model.StringSlice.SectionSlice.FindFieldByPath", "err", err, "path", p.string)
		}
	}

	if val, err := sg.String(p.string); err == nil {
		v = val
	}
	return
}

// String returns the path
func (p Path) String() string {
	return p.string
}

// InScope checks if a field from a path is allowed for current ScopedGetter
func (p Path) InScope(f *config.Field, sg config.ScopedGetter) bool {
	s, _ := sg.Scope()
	return f.Scope.Has(s)
}
