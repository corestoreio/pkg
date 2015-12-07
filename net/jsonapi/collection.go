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

type Collection struct {
	Data  []Data          `json:"data"`
	Links map[string]Link `json:"links,omitempty"`
}

func NewCollection(data []Data) *Collection {
	return &Collection{
		Data:  data,
		Links: map[string]Link{},
	}
}

func (c *Collection) SetLink(name, link string) {
	c.Links[name] = Link{HRef: link}
}

func (c *Collection) SetLinkWithMeta(name string, link string, m Meta) {
	c.Links[name] = Link{HRef: link, Meta: m}
}
