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
	"strings"
)

func lookupString(path string, pkgCfg config.SectionSlice, sg config.ScopedGetter) (v string) {
	if fields, err := pkgCfg.FindFieldByPath(path); err == nil {
		v, _ = cast.ToStringE(fields.Default)
	} else {
		if PkgLog.IsDebug() {
			PkgLog.Debug("model.StringSlice.SectionSlice.FindFieldByPath", "err", err, "path", path)
		}
	}

	if val, err := sg.String(path); err == nil {
		v = val
	}
	return
}

const csvSep = ","

type StringCSV path

func (p StringCSV) Get(pkgCfg config.SectionSlice, sg config.ScopedGetter) []string {
	path := string(p)
	v := lookupString(path, pkgCfg, sg)
	return strings.Split(v, csvSep)
}

func (p StringCSV) Set(w config.Writer, sl []string, s scope.Scope, id int64) error {
	path := string(p)
	return w.Write(config.Path(path), config.Value(strings.Join(sl, csvSep)), config.Scope(s, id))
}

type Int64CSV path

func (p Int64CSV) Get(pkgCfg config.SectionSlice, sg config.ScopedGetter) []int64 {
	path := string(p)
	v := lookupString(path, pkgCfg, sg)
	csv := strings.Split(v, csvSep)
	ret := make([]int64, len(csv))

	for i, line := range csv {
		ret[i], _ = strconv.ParseInt(line, 10, 64)
	}

	return ret
}

func (p Int64CSV) Set(w config.Writer, sl []int64, s scope.Scope, id int64) error {
	path := string(p)

	var final string // todo use bufferpool
	for i, v := range sl {
		final = final + strconv.FormatInt(v, 10)
		if i < len(sl)-1 {
			final = final + csvSep
		}
	}

	return w.Write(config.Path(path), config.Value(final), config.Scope(s, id))
}
