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
	"strings"

	"github.com/corestoreio/csfw/config"
	"github.com/corestoreio/csfw/config/scope"
	"github.com/corestoreio/csfw/util/bufferpool"
	"github.com/juju/errgo"
)

// CSVSeparator separates CSV values
const CSVSeparator = ","

// StringCSV represents a path in config.Getter which will be saved as a
// CSV string and returned as a string slice. Separator is a comma.
type StringCSV struct{ Path }

// Get returns a slice from the 1. default field of a config.SectionSlice
// or 2. from the config.ScopedGetter.
func (p StringCSV) Get(pkgCfg config.SectionSlice, sg config.ScopedGetter) []string {
	return strings.Split(p.Path.LookupString(pkgCfg, sg), CSVSeparator)
}

// Set writes a slice with its scope and ID to the writer
func (p StringCSV) Set(w config.Writer, sl []string, s scope.Scope, id int64) error {
	return p.Path.Set(w, strings.Join(sl, CSVSeparator), s, id)
}

// Int64CSV represents a path in config.Getter which will be saved as a
// CSV string and returned as an int64 slice. Separator is a comma.
type Int64CSV struct{ Path }

func (p Int64CSV) Get(pkgCfg config.SectionSlice, sg config.ScopedGetter) []int64 {
	v := p.Path.LookupString(pkgCfg, sg)
	csv := strings.Split(v, CSVSeparator)
	ret := make([]int64, len(csv))

	for i, line := range csv {
		var err error
		if ret[i], err = strconv.ParseInt(line, 10, 64); err != nil && PkgLog.IsDebug() {
			PkgLog.Debug("model.Int64CSV.Get.strconv.ParseInt", "err", err, "position", i, "line", line)
		}
	}
	return ret
}

// Set writes int64 values as a CSV string
func (p Int64CSV) Set(w config.Writer, sl []int64, s scope.Scope, id int64) error {
	val := bufferpool.Get()
	defer bufferpool.Put(val)
	for i, v := range sl {
		if _, err := val.WriteString(strconv.FormatInt(v, 10)); err != nil {
			return errgo.Mask(err)
		}
		if i < len(sl)-1 {
			if _, err := val.WriteString(CSVSeparator); err != nil {
				return errgo.Mask(err)
			}
		}
	}
	return p.Path.Set(w, val.String(), s, id)
}
