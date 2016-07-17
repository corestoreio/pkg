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
type StoreSlice []Store

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
func (ss StoreSlice) Filter(f func(Store) bool) StoreSlice {
	var stores StoreSlice
	for _, v := range ss {
		if f(v) {
			stores = append(stores, v)
		}
	}
	return stores
}

func (ss StoreSlice) Each(f func(Store)) StoreSlice {
	for i := range ss {
		f(ss[i])
	}
	return ss
}

// Codes returns all store codes
func (ss StoreSlice) Codes() []string {
	if len(ss) == 0 {
		return nil
	}
	var c = make([]string, len(ss))
	for i, st := range ss {
		c[i] = st.StoreCode()
	}
	return c
}

// ActiveCodes returns all active store codes
func (ss StoreSlice) ActiveCodes() []string {
	if len(ss) == 0 {
		return nil
	}
	var c = make([]string, 0, len(ss))
	for _, st := range ss {
		if st.Data.IsActive {
			c = append(c, st.StoreCode())
		}
	}
	return c
}

// IDs returns all store IDs
func (ss StoreSlice) IDs() []int64 {
	if len(ss) == 0 {
		return nil
	}
	ids := make([]int64, len(ss))
	for i, st := range ss {
		ids[i] = st.Data.StoreID
	}
	return ids
}

// ActiveIDs returns all active store IDs
func (ss StoreSlice) ActiveIDs() []int64 {
	if len(ss) == 0 {
		return nil
	}
	var ids = make([]int64, 0, len(ss))
	for _, st := range ss {
		if st.Data.IsActive {
			ids = append(ids, st.Data.StoreID)
		}
	}
	return ids
}
