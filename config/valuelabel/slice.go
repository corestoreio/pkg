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
	"math"
	"sort"

	"github.com/corestoreio/csfw/util"
	"github.com/juju/errgo"
)

// Slice type is returned by the SourceModel.Options() interface
type Slice []Pair

// SortByLabel sorts by label in asc or desc direction
func (s Slice) SortByLabel(d util.SortDirection) Slice {
	var si sort.Interface
	si = vlSortByLabel{s}
	if d == util.SortDesc {
		si = sort.Reverse(si)
	}
	sort.Sort(si)
	return s
}

// SortByValue sorts by value in asc or desc direction. The underlying value
// will be converted to a string. You might expect strange results when sorting
// integers or other non-strings.
func (s Slice) SortByValue(d util.SortDirection) Slice {
	var si sort.Interface
	si = vlSortByValue{s}
	if d == util.SortDesc {
		si = sort.Reverse(si)
	}
	sort.Sort(si)
	return s
}

// SortByInt sorts by field Int in asc or desc direction
func (s Slice) SortByInt(d util.SortDirection) Slice {
	var si sort.Interface
	si = vlSortByInt{s}
	if d == util.SortDesc {
		si = sort.Reverse(si)
	}
	sort.Sort(si)
	return s
}

// SortByFloat64 sorts by field Float64 in asc or desc direction
func (s Slice) SortByFloat64(d util.SortDirection) Slice {
	var si sort.Interface
	si = vlSortByFloat64{s}
	if d == util.SortDesc {
		si = sort.Reverse(si)
	}
	sort.Sort(si)
	return s
}

// SortByBool sorts by field Bool in asc or desc direction
func (s Slice) SortByBool(d util.SortDirection) Slice {
	var si sort.Interface
	si = vlSortByBool{s}
	if d == util.SortDesc {
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

// ContainsKeyString checks if k has an entry as a key.
func (s Slice) ContainsKeyString(k string) bool {
	for _, p := range s {
		if p.String == k {
			return true
		}
	}
	return false
}

// ContainsKeyInt checks if k has an entry as a key.
func (s Slice) ContainsKeyInt(k int) bool {
	for _, p := range s {
		if p.Int == k {
			return true
		}
	}
	return false
}

// ContainsKeyFloat64 checks if k has an entry as a key.
func (s Slice) ContainsKeyFloat64(k float64) bool {
	for _, p := range s {
		abs := math.Abs(p.Float64 - k)
		if abs >= 0 && abs < 0.000001 {
			return true
		}
	}
	return false
}

// ContainsLabel checks if k has an entry as a label.
func (s Slice) ContainsLabel(l string) bool {
	for _, p := range s {
		if p.Label() == l {
			return true
		}
	}
	return false
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
