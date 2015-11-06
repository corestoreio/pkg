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

package valuelabel

import (
	"bytes"
	"encoding/json"
	"sort"

	"github.com/corestoreio/csfw/utils"
	"github.com/juju/errgo"
)

// Slice type is returned by the SourceModel.Options() interface
type Slice []Pair

// Options satisfies the SourceModeller interface
func (s Slice) Options() Slice { return s }

// SortByLabel sorts by label in asc or desc direction
func (s Slice) SortByLabel(d utils.SortDirection) Slice {
	var si sort.Interface
	si = vlSortByLabel{s}
	if d == utils.SortDesc {
		si = sort.Reverse(si)
	}
	sort.Sort(si)
	return s
}

// SortByValue sorts by value in asc or desc direction. The underlying value
// will be converted to a string. You might expect strange results when sorting
// integers or other non-strings.
func (s Slice) SortByValue(d utils.SortDirection) Slice {
	var si sort.Interface
	si = vlSortByValue{s}
	if d == utils.SortDesc {
		si = sort.Reverse(si)
	}
	sort.Sort(si)
	return s
}

// SortByInt sorts by field Int in asc or desc direction
func (s Slice) SortByInt(d utils.SortDirection) Slice {
	var si sort.Interface
	si = vlSortByInt{s}
	if d == utils.SortDesc {
		si = sort.Reverse(si)
	}
	sort.Sort(si)
	return s
}

// SortByFloat64 sorts by field Float64 in asc or desc direction
func (s Slice) SortByFloat64(d utils.SortDirection) Slice {
	var si sort.Interface
	si = vlSortByFloat64{s}
	if d == utils.SortDesc {
		si = sort.Reverse(si)
	}
	sort.Sort(si)
	return s
}

// SortByBool sorts by field Bool in asc or desc direction
func (s Slice) SortByBool(d utils.SortDirection) Slice {
	var si sort.Interface
	si = vlSortByBool{s}
	if d == utils.SortDesc {
		si = sort.Reverse(si)
	}
	sort.Sort(si)
	return s
}

// Len returns the length of the slice
func (s Slice) Len() int { return len(s) }

// Swap swaps elements. Will panic when slice index does not exists.
func (s Slice) Swap(i, j int) { s[i], s[j] = s[j], s[i] }

// ToJSON returns a JSON string, convenience function.
func (s Slice) ToJSON() (string, error) {
	var buf bytes.Buffer
	if err := json.NewEncoder(&buf).Encode(s); err != nil {
		if PkgLog.IsDebug() {
			PkgLog.Debug("config.ValueLabelSlice.ToJSON.Encode", "err", err, "slice", s)
		}
		return "", errgo.Mask(err)
	}
	return buf.String(), nil
}

type (
	vlSortByLabel struct {
		Slice
	}
	vlSortByValue struct {
		Slice
	}
	vlSortByInt struct {
		Slice
	}
	vlSortByFloat64 struct {
		Slice
	}
	vlSortByBool struct {
		Slice
	}
)

func (v vlSortByLabel) Less(i, j int) bool {
	return v.Slice[i].Label() < v.Slice[j].Label()
}

func (v vlSortByValue) Less(i, j int) bool {
	return v.Slice[i].Value() < v.Slice[j].Value()
}

func (v vlSortByInt) Less(i, j int) bool {
	return v.Slice[i].Int < v.Slice[j].Int
}

func (v vlSortByFloat64) Less(i, j int) bool {
	return v.Slice[i].Float64 < v.Slice[j].Float64
}

func (v vlSortByBool) Less(i, j int) bool {
	if !v.Slice[i].Bool {
		return v.Slice[i].Label() < v.Slice[j].Label()
	}
	return v.Slice[i].Bool && v.Slice[i].Label() < v.Slice[j].Label()
}
