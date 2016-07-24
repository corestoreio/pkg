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
	"github.com/corestoreio/csfw/store/scope"
)

// Tree represents a hierarchical structure of all available scopes.
type Tree struct {
	Scope  scope.Scope `json:"scope",xml:"scope"`
	ID     int64       `json:"id",xml:"id"`
	Scopes []Tree      `json:"scopes,omitempty",xml:"scopes,omitempty"`
}
