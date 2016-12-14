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

package slices

import (
	"sort"

	"github.com/corestoreio/csfw/util/errors"
)

var errOutOfRange = errors.NewFatalf("[slices] Index out of range")

// Int64Slice contains Map/Filter/Reduce/Sort/Unique/etc method receivers for
// []int64.
// @todo think about the necessary gen functions
// +gen slice:"Where,Count,GroupBy[int64]"
type Int64 []int64

// ToInt64 converts to type int64 slice.
func (l Int64) ToInt64() []int64 { return []int64(l) }

// Len returns the length
func (l Int64) Len() int { return len(l) }

// Less compares two slice values
func (l Int64) Less(i, j int) bool { return l[i] < l[j] }

// Swap changes the position
func (l Int64) Swap(i, j int) { l[i], l[j] = l[j], l[i] }

// Sort is a convenience method.
func (l Int64) Sort() Int64 { sort.Stable(l); return l }

// Reverse is a convenience method.
func (l Int64) Reverse() Int64 { sort.Stable(sort.Reverse(l)); return l }

// Append adds s (variadic) to the Int64Slice
func (l *Int64) Append(s ...int64) Int64 {
	*l = append(*l, s...)
	return *l
}

// Update sets the int64 s on index i. If index is not found returns an fatal
// error behaviour.
func (l *Int64) Update(i int, s int64) error {
	if i > l.Len() || i < 0 {
		return errOutOfRange
	}
	(*l)[i] = s
	return nil
}

// Delete removes index i from slice
func (l *Int64) Delete(i int) error {
	if i > l.Len()-1 || i < 0 {
		return errOutOfRange
	}
	*l = append((*l)[:i], (*l)[i+1:]...)
	return nil
}

// Index returns -1 if not found or the current index for target t.
func (l Int64) Index(t int64) int {
	for i, v := range l {
		if v == t {
			return i
		}
	}
	return -1
}

// Contains returns true if the target int64 t is in the slice.
func (l Int64) Contains(t int64) bool {
	return l.Index(t) >= 0
}

// Any returns true if one of the int64s in the slice satisfies the predicate f.
func (l Int64) Any(f func(int64) bool) bool {
	for _, v := range l {
		if f(v) {
			return true
		}
	}
	return false
}

// All returns true if all of the int64s in the slice satisfy the predicate f.
func (l Int64) All(f func(int64) bool) bool {
	for _, v := range l {
		if !f(v) {
			return false
		}
	}
	return true
}

// Reduce reduces itself containing all int64s in the slice that satisfy the
// predicate f.
func (l *Int64) Reduce(f func(int64) bool) Int64 {
	vsf := (*l)[:0]
	for _, v := range *l {
		if f(v) {
			vsf = append(vsf, v)
		}
	}
	*l = vsf
	return *l
}

// Map changes itself containing the results of applying the function f to each
// int64 in itself.
func (l *Int64) Map(f func(int64) int64) Int64 {
	for i, v := range *l {
		(*l)[i] = f(v)
	}
	return *l
}

// Sum returns the sum
func (l Int64) Sum() int64 {
	var s int64
	for _, v := range l {
		s += v
	}
	return s
}

// Unique removes duplicate entries in O(n)
func (l *Int64) Unique() Int64 {
	seen := make(map[int64]struct{}, len(*l))
	unique := (*l)[:0]
	for _, p := range *l {
		if _, ok := seen[p]; !ok {
			unique = append(unique, p)
			seen[p] = struct{}{}
		}
	}
	*l = unique
	return *l
}
