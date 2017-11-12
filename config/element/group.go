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

package element

import (
	"encoding/json"
	"sort"

	"github.com/corestoreio/pkg/config/cfgpath"
	"github.com/corestoreio/pkg/storage/text"
	"github.com/corestoreio/pkg/store/scope"
	"github.com/corestoreio/pkg/util/bufferpool"
	"github.com/corestoreio/errors"
)

// GroupSlice contains a set of Groups.
//  Thread safe for reading but not for modifying.
type GroupSlice []Group

// Group defines the layout of a group containing multiple Fields
//  Thread safe for reading but not for modifying.
type Group struct {
	// ID unique ID and merged with others. 2nd part of the path.
	ID      cfgpath.Route
	Label   text.Chars `json:",omitempty"`
	Comment text.Chars `json:",omitempty"`
	// Scopes: bit value eg: showInDefault="1" showInWebsite="1" showInStore="1"
	Scopes    scope.Perm `json:",omitempty"`
	SortOrder int        `json:",omitempty"`

	HelpURL               text.Chars `json:",omitempty"`
	MoreURL               text.Chars `json:",omitempty"` // todo maybe a slice because we might have multiple URLs
	DemoLink              text.Chars `json:",omitempty"`
	HideInSingleStoreMode bool       `json:",omitempty"`

	Fields FieldSlice
	// Groups     GroupSlice @todo see recursive options <xs:element name="group"> in app/code/Magento/Config/etc/system_file.xsd
}

// NewGroupSlice wrapper function, for now.
func NewGroupSlice(gs ...Group) GroupSlice {
	return GroupSlice(gs)
}

// Find returns a Group pointer or ErrGroupNotFound. Route must be a single
// part. E.g. if you have path "a/b/c" route would be in this case "b". For
// comparison the field Sum32 of a route will be used. Error behaviour: NotFound
func (gs GroupSlice) Find(id cfgpath.Route) (Group, int, error) {
	for i, g := range gs {
		if g.ID.Sum32 > 0 && g.ID.Sum32 == id.Sum32 {
			return g, i, nil
		}
	}
	return Group{}, 0, errors.NewNotFoundf("[element] Group %q", id)
}

// Merge copies the data from a groups into this slice. Appends if ID is not
// found in this slice otherwise overrides struct fields if not empty. Not
// thread safe.
func (gs *GroupSlice) Merge(groups ...Group) error {
	for _, g := range groups {
		if err := gs.merge(g); err != nil {
			return errors.Wrap(err, "[element] GroupSlice.Merge")
		}
	}
	return nil
}

func (gs *GroupSlice) merge(g Group) error {
	cg, idx, err := (*gs).Find(g.ID) // cg current group
	if err != nil {
		cg = g
		*gs = append(*gs, cg)
		idx = len(*gs) - 1
	}

	if !g.Label.IsEmpty() {
		cg.Label = g.Label.Clone()
	}
	if !g.Comment.IsEmpty() {
		cg.Comment = g.Comment.Clone()
	}
	if g.Scopes > 0 {
		cg.Scopes = g.Scopes
	}
	if g.SortOrder != 0 {
		cg.SortOrder = g.SortOrder
	}
	if err := cg.Fields.Merge(g.Fields...); err != nil {
		return errors.Wrap(err, "[element] GroupSlice.merge.Fields.Merge")
	}

	(*gs)[idx] = cg
	return nil
}

// ToJSON transforms the whole slice into JSON. If an error occurs the returned
// string starts with: "[element] Error:".
func (gs GroupSlice) ToJSON() string {
	buf := bufferpool.Get()
	defer bufferpool.Put(buf)
	if err := json.NewEncoder(buf).Encode(gs); err != nil {
		return "[element] Error: " + err.Error()
	}
	return buf.String()
}

// Sort convenience helper. Not thread safe.
func (gs GroupSlice) Sort() GroupSlice {
	sort.Sort(gs)
	return gs
}

func (gs GroupSlice) Len() int {
	return len(gs)
}

func (gs GroupSlice) Swap(i, j int) {
	gs[i], gs[j] = gs[j], gs[i]
}

func (gs GroupSlice) Less(i, j int) bool {
	return gs[i].SortOrder < gs[j].SortOrder
}
