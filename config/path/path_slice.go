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

package path

import (
	"bytes"
	"sort"
)

// PathSlice represents a collection of Paths
type PathSlice []Path

// add more functions if needed

// Contains return true if the Path p can be found within the slice.
// It must match ID, Scope and Route.
func (ps PathSlice) Contains(p Path) bool {
	for _, pps := range ps {
		if pps.ID == p.ID && pps.Scope == p.Scope && pps.Sum32 == pps.Sum32 {
			return true
		}
	}
	return false
}

func (ps PathSlice) Len() int { return len(ps) }
func (ps PathSlice) Less(i, j int) bool {
	return bytes.Compare(ps[i].Route.Chars, ps[j].Route.Chars) == -1
}
func (ps PathSlice) Swap(i, j int) { ps[i], ps[j] = ps[j], ps[i] }

// Sort is a convenience method.
func (ps PathSlice) Sort() { sort.Stable(ps) }
