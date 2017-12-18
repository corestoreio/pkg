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
	"math/rand"
	"sort"
	"strings"

	"github.com/corestoreio/errors"
)

// String contains Map/Filter/Reduce/Sort/Unique/etc method receivers for
// []string.
type String []string

// ToString converts to string slice.
func (l String) ToString() []string { return []string(l) }

// Len returns the length
func (l String) Len() int { return len(l) }

// Less compares two slice values
func (l String) Less(i, j int) bool { return l[i] < l[j] }

// Swap changes the position
func (l String) Swap(i, j int) { l[i], l[j] = l[j], l[i] }

// Sort is a convenience method.
func (l String) Sort() String { sort.Stable(l); return l }

// Append adds s (variadic) to the StringSlice
func (l *String) Append(s ...string) String {
	*l = append(*l, s...)
	return *l
}

// Update sets the string s on index i. If index is not found returns an fatal
// error behaviour.
func (l *String) Update(i int, s string) error {
	if i > l.Len() || i < 0 {
		return errors.OutofRange.Newf("[slices] Index %d is out of range", i)
	}
	(*l)[i] = s
	return nil
}

// Delete removes index i from slice
func (l *String) Delete(i int) error {
	if i > l.Len()-1 || i < 0 {
		return errors.OutofRange.Newf("[slices] Index %d is out of range", i)
	}
	*l = append((*l)[:i], (*l)[i+1:]...)
	return nil
}

// Index returns -1 if not found or the current index for target t.
func (l String) Index(t string) int {
	for i, v := range l {
		if v == t {
			return i
		}
	}
	return -1
}

// Contains returns true if the target string t is in the slice.
func (l String) Contains(t string) bool {
	return l.Index(t) >= 0
}

// Any returns true if one of the strings in the slice satisfies the predicate
// f.
func (l String) Any(f func(string) bool) bool {
	for _, v := range l {
		if f(v) {
			return true
		}
	}
	return false
}

// All returns true if all of the strings in the slice satisfy the predicate f.
func (l String) All(f func(string) bool) bool {
	for _, v := range l {
		if !f(v) {
			return false
		}
	}
	return true
}

// Filter filters all strings in the slice that satisfy the predicate f and
// returns a new slice.
func (l String) Filter(f func(string) bool) String {
	vsf := l[:0]
	for _, v := range l {
		if f(v) {
			vsf = append(vsf, v)
		}
	}
	return vsf
}

// Reduce reduces itself containing all strings in the slice that satisfy the
// predicate f.
func (l *String) Reduce(f func(string) bool) String {
	vsf := (*l)[:0]
	for _, v := range *l {
		if f(v) {
			vsf = append(vsf, v)
		}
	}
	*l = vsf
	return *l
}

// ReduceContains reduces itself if the parts of the in slice are contained
// within itself.
func (l *String) ReduceContains(in ...string) String {
	r := (*l)[:0]
	for _, s := range *l {
		isInScope := false
		for _, sin := range in {
			if strings.Contains(s, sin) {
				isInScope = true
				break
			}
		}
		if isInScope == false {
			r = append(r, s)
		}
	}
	*l = r
	return *l
}

// ContainsReverse checks if the StringSlice has at least one occurrence in the
// string s.
func (l String) ContainsReverse(s string) bool {
	for _, substr := range l {
		if strings.Contains(s, substr) {
			return true
		}
	}
	return false
}

// StartsWithReverse checks if the StringSlice can be found at least once in the
// provided string s.
func (l String) StartsWithReverse(s string) bool {
	for _, substr := range l {
		if strings.Index(s, substr) == 0 {
			return true
		}
	}
	return false
}

// Map changes itself containing the results of applying the function f to each
// string in itself.
func (l *String) Map(f func(string) string) String {
	for i, v := range *l {
		(*l)[i] = f(v)
	}
	return *l
}

// Unique removes duplicate entries and discards "" empty strings.
func (l *String) Unique() String {
	unique := (*l)[:0]
	for _, p := range *l {
		found := false
		for _, u := range unique {
			if u == p {
				found = true
				break
			}
		}
		if false == found && p != "" {
			unique = append(unique, p)
		}
	}
	*l = unique
	return *l
}

// Join joins the slice using a separator.
func (l String) Join(sep string) string { return strings.Join(l, sep) }

// Split uses string s and a separator and appends the parts to the slice.
func (l *String) Split(s, sep string) String { return l.Append(strings.Split(s, sep)...) }

// SplitStringer8 uses a name and position indexes to split the name and appends
// the parts to the slice. Cracking the names and indexes which the stringer
// command generates.
func (l *String) SplitStringer8(n string, ps ...uint8) String {
	var next uint8
	ln := uint8(len(n))
	lu := len(ps) - 1
	for i := 0; i < lu; i++ {
		if i+1 < lu {
			next = ps[i+1]
		} else {
			next = ln
		}
		(*l).Append(n[ps[i]:next])
	}
	return *l
}

// Shuffle destroys the order
func (l *String) Shuffle() String {
	for i := range *l {
		j := rand.Intn(i + 1)
		(*l)[i], (*l)[j] = (*l)[j], (*l)[i]
	}
	return *l
}
