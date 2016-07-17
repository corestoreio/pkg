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

package store

import "sort"

// GroupSlice collection of Group. GroupSlice has some nice method receivers.
type GroupSlice []*Group

// Sort convenience helper
func (gs *GroupSlice) Sort() *GroupSlice {
	sort.Stable(gs)
	return gs
}

// Len returns the length of the slice
func (gs GroupSlice) Len() int { return len(gs) }

// Swap swaps positions within the slice
func (gs *GroupSlice) Swap(i, j int) { (*gs)[i], (*gs)[j] = (*gs)[j], (*gs)[i] }

// Less checks the Data field GroupID if index i < index j.
func (gs *GroupSlice) Less(i, j int) bool {
	return (*gs)[i].Data.GroupID < (*gs)[j].Data.GroupID
}

// Filter returns a new slice filtered by predicate f
func (gs GroupSlice) Filter(f func(*Group) bool) GroupSlice {
	var ret GroupSlice
	for _, v := range gs {
		if f(v) {
			ret = append(ret, v)
		}
	}
	return ret
}

func (gs GroupSlice) Each(f func(*Group)) GroupSlice {
	for i := range gs {
		f(gs[i])
	}
	return gs
}

// IDs returns all group IDs
func (gs GroupSlice) IDs() []int64 {
	if len(gs) == 0 {
		return nil
	}
	var ids = make([]int64, 0, len(gs))
	for _, g := range gs {
		ids = append(ids, g.Data.GroupID)
	}
	return ids
}
