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

// StoreSlice a collection of pointers to the Store structs.
// StoreSlice has some nifty method receivers.
type StoreSlice []*Store

// Sort convenience helper
func (ss *StoreSlice) Sort() *StoreSlice {
	sort.Stable(ss)
	return ss
}

func (ss StoreSlice) Len() int { return len(ss) }

func (ss *StoreSlice) Swap(i, j int) { (*ss)[i], (*ss)[j] = (*ss)[j], (*ss)[i] }

// Less depends on the SortOrder
func (ss *StoreSlice) Less(i, j int) bool {
	return (*ss)[i].Data.SortOrder < (*ss)[j].Data.SortOrder
}

// Filter returns a new slice filtered by predicate f
func (ss StoreSlice) Filter(f func(*Store) bool) StoreSlice {
	var stores StoreSlice
	for _, v := range ss {
		if v != nil && f(v) {
			stores = append(stores, v)
		}
	}
	return stores
}

// Codes returns all store codes
func (ss StoreSlice) Codes() []string {
	if len(ss) == 0 {
		return nil
	}
	var c = make([]string, 0, len(ss))
	for _, st := range ss {
		if st != nil {
			c = append(c, st.Data.Code.String)
		}
	}
	return c
}

// IDs returns all store IDs
func (ss StoreSlice) IDs() []int64 {
	if len(ss) == 0 {
		return nil
	}
	var ids = make([]int64, 0, len(ss))
	for _, st := range ss {
		if st != nil {
			ids = append(ids, st.Data.StoreID)
		}
	}
	return ids
}

// LastItem returns the last item of this slice or nil
func (ss StoreSlice) LastItem() *Store {
	if ss.Len() > 0 {
		return ss[ss.Len()-1]
	}
	return nil
}
