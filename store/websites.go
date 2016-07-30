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

import (
	"sort"

	"github.com/corestoreio/csfw/store/scope"
	"github.com/corestoreio/csfw/util/errors"
)

// WebsiteSlice contains pointer to Website struct and some nifty method receivers.
type WebsiteSlice []Website

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
func (ws WebsiteSlice) Less(i, j int) bool {
	return ws[i].Data.SortOrder < ws[j].Data.SortOrder
}

// Filter returns a new slice filtered by predicate f
func (ws WebsiteSlice) Filter(f func(Website) bool) WebsiteSlice {
	var nws = make(WebsiteSlice, 0, len(ws))
	for _, v := range ws {
		if f(v) {
			nws = append(nws, v)
		}
	}
	return nws
}

func (ws WebsiteSlice) Each(f func(Website)) WebsiteSlice {
	for _, w := range ws {
		f(w)
	}
	return ws
}

// Map applies predicate f on each item within the slice and allows changing it.
func (ws WebsiteSlice) Map(f func(*Website)) WebsiteSlice {
	for i, w := range ws {
		f(&w)
		ws[i] = w
	}
	return ws
}

// FindByID filters by Id, returns the website and true if found.
func (ws WebsiteSlice) FindByID(id int64) (Website, bool) {
	for _, w := range ws {
		if w.ID() == id {
			return w, true
		}
	}
	return Website{}, false
}

// Codes returns all website codes
func (ws WebsiteSlice) Codes() []string {
	if len(ws) == 0 {
		return nil
	}
	var c = make([]string, len(ws))
	for i, w := range ws {
		c[i] = w.Data.Code.String
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

// Default returns the default website or a not-found error.
func (ws WebsiteSlice) Default() (Website, error) {
	for _, w := range ws {
		if w.Data.IsDefault.Valid && w.Data.IsDefault.Bool {
			return w, nil
		}
	}
	return Website{}, errors.NewNotFoundf("[store] WebsiteSlice Default Website not found")
}

// Tree represents a hierarchical structure of all available scopes.
type Tree struct {
	Scope  scope.Scope `json:"scope",xml:"scope"`
	ID     int64       `json:"id",xml:"id"`
	Scopes []Tree      `json:"scopes,omitempty",xml:"scopes,omitempty"`
}

// Tree returns the hierarchical overview of the scopes: default -> website
// -> group -> store represented in a Tree.
func (ws WebsiteSlice) Tree() Tree {
	t := Tree{
		Scope: scope.Default,
	}

	t.Scopes = make([]Tree, 0, ws.Len())
	ws.Each(func(w Website) {
		tw := Tree{
			Scope: scope.Website,
			ID:    w.ID(),
		}

		tw.Scopes = make([]Tree, 0, w.Groups.Len())
		w.Groups.Each(func(g Group) {
			tg := Tree{
				Scope: scope.Group,
				ID:    g.ID(),
			}

			tg.Scopes = make([]Tree, 0, g.Stores.Len())
			g.Stores.Each(func(s Store) {
				ts := Tree{
					Scope: scope.Store,
					ID:    s.ID(),
				}
				tg.Scopes = append(tg.Scopes, ts)
			})

			tw.Scopes = append(tw.Scopes, tg)
		})
		t.Scopes = append(t.Scopes, tw)
	})
	return t
}
