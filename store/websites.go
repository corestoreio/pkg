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

// WebsiteSlice contains pointer to Website struct and some nifty method receivers.
type WebsiteSlice []*Website

// Sort convenience helper
func (ws *WebsiteSlice) Sort() *WebsiteSlice {
	sort.Stable(ws)
	return ws
}

// Len returns the length of the slice
func (ws WebsiteSlice) Len() int { return len(ws) }

// Swap swaps positions within the slice
func (ws *WebsiteSlice) Swap(i, j int) { (*ws)[i], (*ws)[j] = (*ws)[j], (*ws)[i] }

// Less checks the Data field SortOrder if index i < index j.
func (ws *WebsiteSlice) Less(i, j int) bool {
	return (*ws)[i].Data.SortOrder < (*ws)[j].Data.SortOrder
}

// Filter returns a new slice filtered by predicate f
func (ws WebsiteSlice) Filter(f func(*Website) bool) WebsiteSlice {
	var nws = make(WebsiteSlice, 0, len(ws))
	for _, v := range ws {
		if f(v) {
			nws = append(nws, v)
		}
	}
	return nws
}

func (ws WebsiteSlice) Each(f func(*Website)) WebsiteSlice {
	for i := range ws {
		f(ws[i])
	}
	return ws
}

// Codes returns all website codes
func (ws WebsiteSlice) Codes() []string {
	if len(ws) == 0 {
		return nil
	}
	var c = make([]string, 0, len(ws))
	for _, w := range ws {
		c = append(c, w.Data.Code.String)
	}
	return c
}

// IDs returns an website IDs
func (ws WebsiteSlice) IDs() []int64 {
	if len(ws) == 0 {
		return nil
	}
	var ids = make([]int64, 0, len(ws))
	for _, w := range ws {
		ids = append(ids, w.Data.WebsiteID)
	}
	return ids
}
