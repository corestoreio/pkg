// Copyright 2015-present, Cyrill @ Schumacher.fm and the CoreStore contributors
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

	"github.com/corestoreio/errors"
	"github.com/corestoreio/pkg/store/scope"
	"github.com/corestoreio/pkg/util/bufferpool"
)

// Groups contains a set of Groups.
//  Thread safe for reading but not for modifying.
type Groups []*Group

// Group defines the layout of a group containing multiple Fields
//  Thread safe for reading but not for modifying.
type Group struct {
	// ID unique ID and merged with others. 2nd part of the path.
	ID      string
	Label   string `json:",omitempty"`
	Comment string `json:",omitempty"`
	// Scopes: bit value eg: showInDefault="1" showInWebsite="1" showInStore="1"
	Scopes    scope.Perm `json:",omitempty"`
	SortOrder int        `json:",omitempty"`

	HelpURL  string `json:",omitempty"`
	MoreURL  string `json:",omitempty"`
	DemoLink string `json:",omitempty"`

	Fields Fields `json:",omitempty"`
	// Groups     Groups @todo see recursive options <xs:element name="group"> in app/code/Magento/Config/etc/system_file.xsd
}

// MakeGroups wrapper function, for now.
func MakeGroups(gs ...*Group) Groups {
	return Groups(gs)
}

// Find returns a Group pointer or ErrGroupNotFound. Route must be a single
// part. E.g. if you have path "a/b/c" route would be in this case "b". For
// comparison the field Sum32 of a route will be used. Error behaviour: NotFound
func (gs Groups) Find(id string) (*Group, int, error) {
	for i, g := range gs {
		if g.ID != "" && g.ID == id {
			return g, i, nil
		}
	}
	return nil, 0, errors.NotFound.Newf("[element] Group %q not found", id)
}

// Merge copies the data from a groups into this slice and returns the new
// slice. Appends if ID is not found in this slice otherwise overrides struct
// fields if not empty. Not thread safe.
func (gs Groups) Merge(groups ...*Group) Groups {
	for _, g := range groups {
		gs = gs.merge(g)
	}
	return gs
}

func (gs Groups) merge(g *Group) Groups {
	cg, idx, err := gs.Find(g.ID) // cg current group
	if err != nil {
		cg = g
		gs = append(gs, cg)
		idx = len(gs) - 1
	}

	if "" != g.Label {
		cg.Label = g.Label
	}
	if "" != g.Comment {
		cg.Comment = g.Comment
	}
	if g.Scopes > 0 {
		cg.Scopes = g.Scopes
	}
	if g.SortOrder != 0 {
		cg.SortOrder = g.SortOrder
	}
	cg.Fields = cg.Fields.Merge(g.Fields...)

	gs[idx] = cg
	return gs
}

// ToJSON transforms the whole slice into JSON. If an error occurs the returned
// string starts with: "[element] Error:".
func (gs Groups) ToJSON() string {
	buf := bufferpool.Get()
	defer bufferpool.Put(buf)
	if err := json.NewEncoder(buf).Encode(gs); err != nil {
		return "[element] Error: " + err.Error()
	}
	return buf.String()
}

// Sort convenience helper. Not thread safe.
func (gs Groups) Sort() Groups {
	sort.Sort(gs)
	return gs
}

func (gs Groups) Len() int {
	return len(gs)
}

func (gs Groups) Swap(i, j int) {
	gs[i], gs[j] = gs[j], gs[i]
}

func (gs Groups) Less(i, j int) bool {
	return gs[i].SortOrder < gs[j].SortOrder
}
