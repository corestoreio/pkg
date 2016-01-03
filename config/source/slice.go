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

package source

import (
	"encoding/json"
	"math"
	"sort"

	"github.com/corestoreio/csfw/util/bufferpool"
	"github.com/juju/errgo"
)

// Slice type is returned by the SourceModel.Options() interface
type Slice []Pair

// SortByLabel sorts by label in asc = 0 or desc != 0 direction
func (s Slice) SortByLabel(direction int) Slice {
	var si sort.Interface
	si = vlSortByLabel{s}
	if 0 != direction {
		si = sort.Reverse(si)
	}
	sort.Sort(si)
	return s
}

// SortByValue sorts by value in asc = 0 or desc != 0 direction. The underlying value
// will be converted to a string. You might expect strange results when sorting
// integers or other non-strings.
func (s Slice) SortByValue(direction int) Slice {
	var si sort.Interface
	si = vlSortByValue{s}
	if 0 != direction {
		si = sort.Reverse(si)
	}
	sort.Sort(si)
	return s
}

// SortByInt sorts by field Int in asc = 0 or desc != 0 direction
func (s Slice) SortByInt(direction int) Slice {
	var si sort.Interface
	si = vlSortByInt{s}
	if 0 != direction {
		si = sort.Reverse(si)
	}
	sort.Sort(si)
	return s
}

// SortByFloat64 sorts by field Float64 in asc = 0 or desc != 0 direction
func (s Slice) SortByFloat64(direction int) Slice {
	var si sort.Interface
	si = vlSortByFloat64{s}
	if 0 != direction {
		si = sort.Reverse(si)
	}
	sort.Sort(si)
	return s
}

// SortByBool sorts by field Bool in asc = 0 or desc != 0 direction
func (s Slice) SortByBool(direction int) Slice {
	var si sort.Interface
	si = vlSortByBool{s}
	if 0 != direction {
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
	var buf = bufferpool.Get()
	defer bufferpool.Put(buf)
	if err := json.NewEncoder(buf).Encode(s); err != nil {
		if PkgLog.IsDebug() {
			PkgLog.Debug("config.ValueLabelSlice.ToJSON.Encode", "err", err, "slice", s)
		}
		return "", errgo.Mask(err)
	}
	return buf.String(), nil
}

// ContainsValString checks if value k exists.
func (s Slice) ContainsValString(k string) bool {
	return s.IndexValString(k) > -1
}

// IndexValString checks if value k exists and returns its position.
// Returns -1 when the value was not found.
func (s Slice) IndexValString(k string) int {
	for i, p := range s {
		if p.NotNull == NotNullString && p.String == k {
			return i
		}
	}
	return -1
}

// ContainsValInt checks if value k exists.
func (s Slice) ContainsValInt(k int) bool {
	return s.IndexValInt(k) > -1
}

// IndexValInt checks if value k exists and returns its position.
// Returns -1 when the value was not found.
func (s Slice) IndexValInt(k int) int {
	for i, p := range s {
		if p.NotNull == NotNullInt && p.Int == k {
			return i
		}
	}
	return -1
}

// ContainsValFloat64 checks if value k exists.
func (s Slice) ContainsValFloat64(k float64) bool {
	return s.IndexValFloat64(k) > -1
}

// IndexValFloat64 checks if value k exists and returns its position.
// Returns -1 when the value was not found.
func (s Slice) IndexValFloat64(k float64) int {
	for i, p := range s {
		abs := math.Abs(p.Float64 - k)
		if p.NotNull == NotNullFloat64 && abs >= 0 && abs < 0.0000001 { // hmmmmm better way?
			return i
		}
	}
	return -1
}

// ContainsValBool checks if value k exists.
func (s Slice) ContainsValBool(k bool) bool {
	return s.IndexValBool(k) > -1
}

// IndexValBool checks if value k exists and returns its position.
// Returns -1 when the value was not found.
func (s Slice) IndexValBool(k bool) int {
	for i, p := range s {
		if p.NotNull == NotNullBool && p.Bool == k {
			return i
		}
	}
	return -1
}

// ContainsLabel checks if k has an entry as a label.
func (s Slice) ContainsLabel(l string) bool {
	return s.IndexLabel(l) > -1
}

// IndexLabel checks if label l exists and returns its first position.
// Returns -1 when the label was not found.
func (s Slice) IndexLabel(l string) int {
	for i, p := range s {
		if p.Label() == l {
			return i
		}
	}
	return -1
}

// Merge integrates the argument Slice into the receiver slice and overwrites the
// existing values of the receiver slice.
func (s *Slice) Merge(sl Slice) Slice {

	for _, p := range sl {
		var idx = -1

		switch p.NotNull {
		case NotNullString:
			idx = s.IndexValString(p.String)
		case NotNullInt:
			idx = s.IndexValInt(p.Int)
		case NotNullFloat64:
			idx = s.IndexValFloat64(p.Float64)
		case NotNullBool:
			idx = s.IndexValBool(p.Bool)
		}

		if idx > -1 {
			(*s)[idx].label = p.Label()
			(*s)[idx].NotNull = p.NotNull
		} else {
			*s = append(*s, p)
		}

	}
	return *s
}

// Unique removes duplicate entries.
func (s *Slice) Unique() Slice {
	unique := (*s)[:0]
	for _, p := range *s {

		found := false

		switch p.NotNull {
		case NotNullString:
			found = unique.ContainsValString(p.String)
		case NotNullInt:
			found = unique.ContainsValInt(p.Int)
		case NotNullFloat64:
			found = unique.ContainsValFloat64(p.Float64)
		case NotNullBool:
			found = unique.ContainsValBool(p.Bool)
		}

		if false == found {
			unique = append(unique, p)
		}

	}
	*s = unique
	return *s
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
