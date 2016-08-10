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
	var i int
	// pre-calculate the size because append is slow.
	for _, v := range ss {
		if f(v) {
			i++
		}
	}
	stores := make(StoreSlice, i)
	i = 0
	for _, s := range ss {
		if f(s) {
			stores[i] = s
			i++
		}
	}
	return stores
}

// Each applies predicate f on each item within the slice.
func (ss StoreSlice) Each(f func(Store)) StoreSlice {
	for _, s := range ss {
		f(s)
	}
	return ss
}

// Map applies predicate f on each item within the slice and allows changing it.
func (ss StoreSlice) Map(f func(*Store)) StoreSlice {
	for i, s := range ss {
		f(&s)
		ss[i] = s
	}
	return ss
}

// FindByID filters by Id, returns the website and true if found.
func (ss StoreSlice) FindByID(id int64) (Store, bool) {
	for _, s := range ss {
		if s.ID() == id {
			return s, true
		}
	}
	return Store{}, false
}

// FindOne filters by predicate f and returns the first hit.
func (ss StoreSlice) FindOne(f func(Store) bool) (Store, bool) {
	for _, s := range ss {
		if f(s) {
			return s, true
		}
	}
	return Store{}, false
}

// Codes returns all store codes
func (ss StoreSlice) Codes() []string {
	if len(ss) == 0 {
		return nil
	}
	var c = make([]string, len(ss))
	for i, st := range ss {
		c[i] = st.Code()
	}
	return c
}

func (ss StoreSlice) activeCount() (i int) {
	for _, st := range ss {
		if st.Data != nil && st.Data.IsActive {
			i++
		}
	}
	return
}

// ActiveCodes returns all active store codes
func (ss StoreSlice) ActiveCodes() []string {
	if len(ss) == 0 {
		return nil
	}

	i := ss.activeCount()
	c := make([]string, i)
	i = 0
	for _, st := range ss {
		if st.Data != nil && st.Data.IsActive {
			c[i] = st.Code()
			i++
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
	i := ss.activeCount()

	ids := make([]int64, i)
	i = 0
	for _, st := range ss {
		if st.Data != nil && st.Data.IsActive {
			ids[i] = st.Data.StoreID
			i++
		}
	}
	return ids
}
