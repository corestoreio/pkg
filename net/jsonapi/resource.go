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

package jsonapi

import "encoding/json"

type Resource struct {
	Data Data `json:"data"`
}

func NewResource(typ string, id string, v interface{}) *Resource {
	raw, _ := json.Marshal(v)
	return &Resource{
		Data: Data{
			Type:          typ,
			ID:            id,
			Attributes:    json.RawMessage(raw),
			Relationships: map[string]Relationship{},
			Links:         map[string]string{},
		},
	}
}

func (r *Resource) SetRelationship(name, typ, id, self, related string) {
	r.Data.Relationships[name] = Relationship{
		Data: RelationshipData{
			Type: typ,
			ID:   id,
		},
		Links: RelationshipLinks{
			Self:    self,
			Related: related,
		},
	}
}

func (r *Resource) SetLink(name, link string) {
	r.Data.Links[name] = link
}

type Data struct {
	Type          string                  `json:"type"`
	ID            string                  `json:"id"`
	Attributes    json.RawMessage         `json:"attributes"`
	Relationships map[string]Relationship `json:"relationships,omitempty"`
	Links         map[string]string       `json:"links,omitempty"`
}

type Relationship struct {
	Links RelationshipLinks `json:"links"`
	Data  RelationshipData  `json:"data"`
}

type RelationshipData struct {
	Type string `json:"type"`
	ID   string `json:"id"`
}

type RelationshipLinks struct {
	Self    string `json:"self"`
	Related string `json:"related,omitempty"`
}
