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
	"github.com/corestoreio/csfw/util/cast"
	"strconv"
)

// Bool represents a path in config.Getter which handles bool values.
type Bool path

func (p Bool) Get(pkgCfg config.SectionSlice, sg config.ScopedGetter) (v bool) {
	aPath := string(p)
	if fields, err := pkgCfg.FindFieldByPath(aPath); err == nil {
		v, _ = cast.ToBoolE(fields.Default)
	} else {
		if PkgLog.IsDebug() {
			PkgLog.Debug("model.Bool.SectionSlice.FindFieldByPath", "err", err, "path", aPath)
		}
	}

	if val, err := sg.Bool(aPath); err == nil {
		v = val
	}
	return v
}

func (p Bool) Set(w config.Writer, v bool, s scope.Scope, id int64) error {
	return path(p).set(w, strconv.FormatBool(v), s, id)
}

// String represents a path in config.Getter which handles string values.
type String path

func (p String) Get(pkgCfg config.SectionSlice, sg config.ScopedGetter) (v string) {
	path := string(p)
	if fields, err := pkgCfg.FindFieldByPath(path); err == nil {
		v, _ = cast.ToStringE(fields.Default)
	} else {
		if PkgLog.IsDebug() {
			PkgLog.Debug("model.String.SectionSlice.FindFieldByPath", "err", err, "path", path)
		}
	}

	if val, err := sg.String(path); err == nil {
		v = val
	}
	return v
}

func (p String) Set(w config.Writer, v string, s scope.Scope, id int64) error {
	return path(p).set(w, v, s, id)
}

// Int represents a path in config.Getter which handles int values.
type Int path

func (p Int) Get(pkgCfg config.SectionSlice, sg config.ScopedGetter) (v int) {
	path := string(p)
	if fields, err := pkgCfg.FindFieldByPath(path); err == nil {
		v, _ = cast.ToIntE(fields.Default)
	} else {
		if PkgLog.IsDebug() {
			PkgLog.Debug("model.Int.SectionSlice.FindFieldByPath", "err", err, "path", path)
		}
	}

	if val, err := sg.Int(path); err == nil {
		v = val
	}
	return v
}

func (p Int) Set(w config.Writer, v int, s scope.Scope, id int64) error {
	return path(p).set(w, strconv.Itoa(v), s, id)
}
