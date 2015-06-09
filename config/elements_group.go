// Copyright 2015, Cyrill @ Schumacher.fm and the CoreStore contributors
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

package config

import (
	"bytes"
	"encoding/json"
	"errors"
	"sort"

	"github.com/corestoreio/csfw/utils/log"
	"github.com/juju/errgo"
)

// ErrGroupNotFound error when a group cannot be found
var ErrGroupNotFound = errors.New("Group not found")

type (

	// GroupSlice contains a set of Groups
	GroupSlice []*Group
	// Group defines the layout of a group containing multiple Fields
	Group struct {
		// ID unique ID and merged with others. 2nd part of the path.
		ID      string
		Label   string `json:",omitempty"`
		Comment string `json:",omitempty"`
		// Scope: bit value eg: showInDefault="1" showInWebsite="1" showInStore="1"
		Scope     ScopePerm `json:",omitempty"`
		SortOrder int       `json:",omitempty"`
		Fields    FieldSlice
	}
)

// FindByID returns a Group pointer or nil if not found
func (gs GroupSlice) FindByID(id string) (*Group, error) {
	for _, g := range gs {
		if g != nil && g.ID == id {
			return g, nil
		}
	}
	return nil, ErrGroupNotFound
}

// Append adds *Group (variadic) to the GroupSlice
func (gs *GroupSlice) Append(g ...*Group) *GroupSlice {
	*gs = append(*gs, g...)
	return gs
}

// Merge copies the data from a groups into this slice. Appends if ID is not found
// in this slice otherwise overrides struct fields if not empty.
func (gs *GroupSlice) Merge(groups ...*Group) error {
	for _, g := range groups {
		if err := (*gs).merge(g); err != nil {
			return errgo.Mask(err)
		}
	}
	return nil
}

func (gs *GroupSlice) merge(g *Group) error {
	if g == nil {
		return nil
	}
	cg, err := (*gs).FindByID(g.ID) // cg current group
	if cg == nil || err != nil {
		cg = g
		(*gs).Append(cg)
	}

	if g.Label != "" {
		cg.Label = g.Label
	}
	if g.Comment != "" {
		cg.Comment = g.Comment
	}
	if g.Scope > 0 {
		cg.Scope = g.Scope
	}
	if g.SortOrder != 0 {
		cg.SortOrder = g.SortOrder
	}
	cg.Fields.Merge(g.Fields...)
	return nil
}

// ToJSON transforms the whole slice into JSON
func (gs GroupSlice) ToJSON() string {
	var buf bytes.Buffer
	if err := json.NewEncoder(&buf).Encode(gs); err != nil {
		log.Error("GroupSlice=ToJSON", "err", err)
		return ""
	}
	return buf.String()
}

// Sort convenience helper
func (gs *GroupSlice) Sort() *GroupSlice {
	sort.Sort(gs)
	return gs
}

func (gs *GroupSlice) Len() int {
	return len(*gs)
}

func (gs *GroupSlice) Swap(i, j int) {
	(*gs)[i], (*gs)[j] = (*gs)[j], (*gs)[i]
}

func (gs *GroupSlice) Less(i, j int) bool {
	return (*gs)[i].SortOrder < (*gs)[j].SortOrder
}
